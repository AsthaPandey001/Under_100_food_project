package fooddiscovery

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AvailabilityStatus represents the calculated operational state of a restaurant
type AvailabilityStatus string

const (
	StatusOpenNow           AvailabilityStatus = "open_now"
	StatusClosedNow         AvailabilityStatus = "closed_now"
	StatusOpensAt           AvailabilityStatus = "opens_at"
	StatusClosesAt          AvailabilityStatus = "closes_at"
	StatusHoursNotAvailable AvailabilityStatus = "hours_not_available"
	StatusTemporarilyClosed AvailabilityStatus = "temporarily_closed"
)

// OperatingPeriod represents a scheduled operating time block for a restaurant
type OperatingPeriod struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurant_id"`
	DayOfWeek    int       `json:"day_of_week"` // 1 = Monday, 2 = Tuesday, ..., 7 = Sunday
	OpenTime     string    `json:"open_time"`  // 24-hr format "HH:MM", e.g. "07:00", "11:30"
	CloseTime    string    `json:"close_time"` // 24-hr format "HH:MM", e.g. "12:00", "23:00", "01:00"
	IsClosed     bool      `json:"is_closed"`
	Is24Hours    bool      `json:"is_24_hours"`
	PeriodIndex  int       `json:"period_index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SpecialHours represents special overrides or temporary closures
type SpecialHours struct {
	ID           int        `json:"id"`
	RestaurantID int        `json:"restaurant_id"`
	SpecialDate  *time.Time `json:"special_date,omitempty"`  // YYYY-MM-DD
	ClosureStart *time.Time `json:"closure_start,omitempty"` // UTC / Local timestamp
	ClosureEnd   *time.Time `json:"closure_end,omitempty"`   // UTC / Local timestamp
	IsClosed     bool       `json:"is_closed"`
	Reason       string     `json:"reason"`
	OpenTime     string     `json:"open_time,omitempty"`
	CloseTime    string     `json:"close_time,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AvailabilityResult represents the evaluated status for customer presentation
type AvailabilityResult struct {
	Status         AvailabilityStatus `json:"status"`
	IsOpen         bool               `json:"is_open"`
	BadgeText      string             `json:"badge_text"`       // "Open Now", "Closed", "Hours not available", "Temporarily Closed"
	DetailText     string             `json:"detail_text"`      // "Closes at 10:30 PM", "Opens at 7:00 AM", etc.
	NextChangeTime *time.Time         `json:"next_change_time"` // optional next state transition time
	Notice         string             `json:"notice,omitempty"` // temporary closure notice/reason
}

// EnsureOpeningHoursTables creates operating hours tables idempotently
func (r *Repo) EnsureOpeningHoursTables(ctx context.Context) {
	if r.db == nil {
		return
	}

	ddl := `
	CREATE TABLE IF NOT EXISTS restaurant_operating_hours (
		id SERIAL PRIMARY KEY,
		restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
		open_time VARCHAR(5) NOT NULL DEFAULT '09:00',
		close_time VARCHAR(5) NOT NULL DEFAULT '22:00',
		is_closed BOOLEAN NOT NULL DEFAULT false,
		is_24_hours BOOLEAN NOT NULL DEFAULT false,
		period_index INT NOT NULL DEFAULT 1,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE(restaurant_id, day_of_week, period_index)
	);

	CREATE INDEX IF NOT EXISTS idx_resto_hours_resto_id ON restaurant_operating_hours(restaurant_id);
	CREATE INDEX IF NOT EXISTS idx_resto_hours_day ON restaurant_operating_hours(day_of_week);

	CREATE TABLE IF NOT EXISTS restaurant_special_hours (
		id SERIAL PRIMARY KEY,
		restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		special_date DATE,
		closure_start TIMESTAMPTZ,
		closure_end TIMESTAMPTZ,
		is_closed BOOLEAN NOT NULL DEFAULT true,
		reason VARCHAR(255) DEFAULT '',
		open_time VARCHAR(5) DEFAULT '',
		close_time VARCHAR(5) DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_resto_special_resto_id ON restaurant_special_hours(restaurant_id);
	CREATE INDEX IF NOT EXISTS idx_resto_special_range ON restaurant_special_hours(restaurant_id, closure_start, closure_end);
	`
	_, _ = r.db.ExecContext(ctx, ddl)
}

// GetKolkataLocation returns the Asia/Kolkata timezone location
func GetKolkataLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err == nil {
		return loc
	}
	// Fallback to Indian Standard Time (UTC+05:30)
	return time.FixedZone("IST", 5*3600+30*60)
}

// WeekdayToDayOfWeek maps Go's time.Weekday (0=Sunday..6=Saturday) to standard 1=Monday..7=Sunday
func WeekdayToDayOfWeek(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

// DayOfWeekToString returns English name for day of week (1=Monday..7=Sunday)
func DayOfWeekToString(d int) string {
	switch d {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	case 7:
		return "Sunday"
	default:
		return "Unknown"
	}
}

// ParseMinutes parses a 24-hr time string "HH:MM" into minutes from midnight
func ParseMinutes(timeStr string) (int, error) {
	timeStr = strings.TrimSpace(timeStr)
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s (expected HH:MM)", timeStr)
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 24 {
		return 0, fmt.Errorf("invalid hour: %s", parts[0])
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid minute: %s", parts[1])
	}
	if h == 24 && m > 0 {
		return 0, fmt.Errorf("invalid time: %s", timeStr)
	}
	return h*60 + m, nil
}

// FormatTime12H converts "HH:MM" or minutes into human-readable 12-hour format e.g. "7:00 AM", "10:30 PM"
func FormatTime12H(timeStr string) string {
	mins, err := ParseMinutes(timeStr)
	if err != nil {
		return timeStr
	}
	h := (mins / 60) % 24
	m := mins % 60
	period := "AM"
	if h >= 12 {
		period = "PM"
	}
	displayH := h % 12
	if displayH == 0 {
		displayH = 12
	}
	return fmt.Sprintf("%d:%02d %s", displayH, m, period)
}

// CalculateAvailability evaluates restaurant operating status given configured periods, special hours, and reference time
func CalculateAvailability(periods []OperatingPeriod, specials []SpecialHours, refTime time.Time) AvailabilityResult {
	if len(periods) == 0 {
		return AvailabilityResult{
			Status:     StatusHoursNotAvailable,
			IsOpen:     false,
			BadgeText:  "Hours not available",
			DetailText: "Opening hours not listed",
		}
	}

	loc := GetKolkataLocation()
	t := refTime.In(loc)

	// 1. Check Temporary Closures
	for _, spec := range specials {
		if spec.IsClosed && spec.ClosureStart != nil && spec.ClosureEnd != nil {
			start := spec.ClosureStart.In(loc)
			end := spec.ClosureEnd.In(loc)
			if (t.Equal(start) || t.After(start)) && t.Before(end) {
				notice := spec.Reason
				if notice == "" {
					notice = "Temporarily closed for maintenance"
				}
				reopenStr := end.Format("Jan 02, 3:04 PM")
				return AvailabilityResult{
					Status:         StatusTemporarilyClosed,
					IsOpen:         false,
					BadgeText:      "Temporarily Closed",
					DetailText:     fmt.Sprintf("Reopens %s", reopenStr),
					Notice:         notice,
					NextChangeTime: &end,
				}
			}
		}
	}

	todayDayOfWeek := WeekdayToDayOfWeek(t.Weekday())
	currentMinutes := t.Hour()*60 + t.Minute()

	// Check if today has a special date override
	var todayOverride *SpecialHours
	todayDateStr := t.Format("2006-01-02")
	for i, spec := range specials {
		if spec.SpecialDate != nil {
			specDateStr := spec.SpecialDate.In(loc).Format("2006-01-02")
			if specDateStr == todayDateStr {
				todayOverride = &specials[i]
				break
			}
		}
	}

	if todayOverride != nil && todayOverride.IsClosed {
		reason := todayOverride.Reason
		if reason == "" {
			reason = "Closed for holiday"
		}
		return AvailabilityResult{
			Status:     StatusClosedNow,
			IsOpen:     false,
			BadgeText:  "Closed Today",
			DetailText: reason,
		}
	}

	// 2. Check Yesterday's Overnight Periods (e.g. Yesterday started 18:00 and closes at 02:00 today)
	yesterdayDayOfWeek := todayDayOfWeek - 1
	if yesterdayDayOfWeek < 1 {
		yesterdayDayOfWeek = 7
	}

	for _, p := range periods {
		if p.DayOfWeek == yesterdayDayOfWeek && !p.IsClosed && !p.Is24Hours {
			openM, errO := ParseMinutes(p.OpenTime)
			closeM, errC := ParseMinutes(p.CloseTime)
			if errO == nil && errC == nil && closeM <= openM {
				// Overnight shift spanning past midnight
				if currentMinutes < closeM {
					closeTimeToday := time.Date(t.Year(), t.Month(), t.Day(), closeM/60, closeM%60, 0, 0, loc)
					return AvailabilityResult{
						Status:         StatusOpenNow,
						IsOpen:         true,
						BadgeText:      "Open Now",
						DetailText:     fmt.Sprintf("Closes at %s", FormatTime12H(p.CloseTime)),
						NextChangeTime: &closeTimeToday,
					}
				}
			}
		}
	}

	// 3. Evaluate Today's Configured Periods
	var todayPeriods []OperatingPeriod
	for _, p := range periods {
		if p.DayOfWeek == todayDayOfWeek {
			todayPeriods = append(todayPeriods, p)
		}
	}

	// Check if today is completely marked closed
	if len(todayPeriods) > 0 {
		allClosed := true
		for _, p := range todayPeriods {
			if !p.IsClosed {
				allClosed = false
				break
			}
		}
		if allClosed {
			// Find next opening day
			nextDayText := findNextOpeningDay(periods, todayDayOfWeek, loc, t)
			return AvailabilityResult{
				Status:     StatusClosedNow,
				IsOpen:     false,
				BadgeText:  "Closed Today",
				DetailText: nextDayText,
			}
		}
	}

	// Check 24 Hours operation today
	for _, p := range todayPeriods {
		if !p.IsClosed && p.Is24Hours {
			return AvailabilityResult{
				Status:     StatusOpenNow,
				IsOpen:     true,
				BadgeText:  "Open 24 Hours",
				DetailText: "Open all day & night",
			}
		}
	}

	// Check if open right now in any of today's periods
	for _, p := range todayPeriods {
		if p.IsClosed {
			continue
		}
		openM, errO := ParseMinutes(p.OpenTime)
		closeM, errC := ParseMinutes(p.CloseTime)
		if errO != nil || errC != nil {
			continue
		}

		if closeM > openM {
			// Standard same-day period: [openM, closeM)
			if currentMinutes >= openM && currentMinutes < closeM {
				closeTime := time.Date(t.Year(), t.Month(), t.Day(), closeM/60, closeM%60, 0, 0, loc)
				return AvailabilityResult{
					Status:         StatusOpenNow,
					IsOpen:         true,
					BadgeText:      "Open Now",
					DetailText:     fmt.Sprintf("Closes at %s", FormatTime12H(p.CloseTime)),
					NextChangeTime: &closeTime,
				}
			}
		} else {
			// Overnight shift: [openM, 24:00)
			if currentMinutes >= openM {
				tomorrow := t.AddDate(0, 0, 1)
				closeTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), closeM/60, closeM%60, 0, 0, loc)
				return AvailabilityResult{
					Status:         StatusOpenNow,
					IsOpen:         true,
					BadgeText:      "Open Now",
					DetailText:     fmt.Sprintf("Closes at %s", FormatTime12H(p.CloseTime)),
					NextChangeTime: &closeTime,
				}
			}
		}
	}

	// 4. Restaurant is NOT currently open. Find Next Opening Time.
	// First check remaining periods today
	var nextOpenToday *OperatingPeriod
	var minDiff = 999999
	for i, p := range todayPeriods {
		if p.IsClosed {
			continue
		}
		openM, err := ParseMinutes(p.OpenTime)
		if err == nil && openM > currentMinutes {
			diff := openM - currentMinutes
			if diff < minDiff {
				minDiff = diff
				nextOpenToday = &todayPeriods[i]
			}
		}
	}

	if nextOpenToday != nil {
		openM, _ := ParseMinutes(nextOpenToday.OpenTime)
		nextOpenTime := time.Date(t.Year(), t.Month(), t.Day(), openM/60, openM%60, 0, 0, loc)
		return AvailabilityResult{
			Status:         StatusClosedNow,
			IsOpen:         false,
			BadgeText:      "Closed Now",
			DetailText:     fmt.Sprintf("Opens at %s", FormatTime12H(nextOpenToday.OpenTime)),
			NextChangeTime: &nextOpenTime,
		}
	}

	// Check subsequent days (tomorrow up to 7 days ahead)
	nextDayText := findNextOpeningDay(periods, todayDayOfWeek, loc, t)
	return AvailabilityResult{
		Status:     StatusClosedNow,
		IsOpen:     false,
		BadgeText:  "Closed Now",
		DetailText: nextDayText,
	}
}

// findNextOpeningDay searches the upcoming days for the next opening period
func findNextOpeningDay(periods []OperatingPeriod, currentDayOfWeek int, loc *time.Location, now time.Time) string {
	for offset := 1; offset <= 7; offset++ {
		targetDay := ((currentDayOfWeek - 1 + offset) % 7) + 1
		for _, p := range periods {
			if p.DayOfWeek == targetDay && !p.IsClosed {
				dayName := DayOfWeekToString(targetDay)
				if offset == 1 {
					dayName = "tomorrow"
				}
				if p.Is24Hours {
					return fmt.Sprintf("Opens %s (24 hours)", dayName)
				}
				return fmt.Sprintf("Opens at %s %s", FormatTime12H(p.OpenTime), dayName)
			}
		}
	}
	return "Currently closed"
}

// GetRestaurantOperatingHours returns configured periods for a restaurant
func (r *Repo) GetRestaurantOperatingHours(ctx context.Context, restaurantID int) ([]OperatingPeriod, error) {
	if r.db == nil {
		return []OperatingPeriod{}, nil
	}

	query := `
		SELECT id, restaurant_id, day_of_week, open_time, close_time, is_closed, is_24_hours, period_index, created_at, updated_at
		FROM restaurant_operating_hours
		WHERE restaurant_id = $1
		ORDER BY day_of_week ASC, period_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query operating hours: %w", err)
	}
	defer rows.Close()

	var periods []OperatingPeriod
	for rows.Next() {
		var p OperatingPeriod
		err := rows.Scan(
			&p.ID, &p.RestaurantID, &p.DayOfWeek, &p.OpenTime, &p.CloseTime,
			&p.IsClosed, &p.Is24Hours, &p.PeriodIndex, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		periods = append(periods, p)
	}
	return periods, nil
}

// GetRestaurantSpecialHours returns active special hours / temporary closures for a restaurant
func (r *Repo) GetRestaurantSpecialHours(ctx context.Context, restaurantID int) ([]SpecialHours, error) {
	if r.db == nil {
		return []SpecialHours{}, nil
	}

	query := `
		SELECT id, restaurant_id, special_date, closure_start, closure_end, is_closed, reason, open_time, close_time, created_at, updated_at
		FROM restaurant_special_hours
		WHERE restaurant_id = $1
		ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query special hours: %w", err)
	}
	defer rows.Close()

	var specials []SpecialHours
	for rows.Next() {
		var s SpecialHours
		var specDate, cStart, cEnd sql.NullTime
		var reason, openT, closeT sql.NullString
		err := rows.Scan(
			&s.ID, &s.RestaurantID, &specDate, &cStart, &cEnd, &s.IsClosed,
			&reason, &openT, &closeT, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if specDate.Valid {
			s.SpecialDate = &specDate.Time
		}
		if cStart.Valid {
			s.ClosureStart = &cStart.Time
		}
		if cEnd.Valid {
			s.ClosureEnd = &cEnd.Time
		}
		s.Reason = reason.String
		s.OpenTime = openT.String
		s.CloseTime = closeT.String
		specials = append(specials, s)
	}
	return specials, nil
}

// ValidateOperatingPeriods checks periods for valid formats, valid days, and time overlaps
func ValidateOperatingPeriods(periods []OperatingPeriod) error {
	byDay := make(map[int][]OperatingPeriod)
	for _, p := range periods {
		if p.DayOfWeek < 1 || p.DayOfWeek > 7 {
			return fmt.Errorf("invalid day of week: %d (expected 1-7)", p.DayOfWeek)
		}
		if !p.IsClosed && !p.Is24Hours {
			openM, errO := ParseMinutes(p.OpenTime)
			if errO != nil {
				return fmt.Errorf("invalid open time for %s: %s", DayOfWeekToString(p.DayOfWeek), p.OpenTime)
			}
			closeM, errC := ParseMinutes(p.CloseTime)
			if errC != nil {
				return fmt.Errorf("invalid close time for %s: %s", DayOfWeekToString(p.DayOfWeek), p.CloseTime)
			}
			if openM == closeM {
				return fmt.Errorf("open and close times cannot be equal for %s (%s)", DayOfWeekToString(p.DayOfWeek), p.OpenTime)
			}
		}
		byDay[p.DayOfWeek] = append(byDay[p.DayOfWeek], p)
	}

	// Check for overlapping periods on the same day
	for day, dayPeriods := range byDay {
		if len(dayPeriods) <= 1 {
			continue
		}
		for i := 0; i < len(dayPeriods); i++ {
			p1 := dayPeriods[i]
			if p1.IsClosed || p1.Is24Hours {
				continue
			}
			o1, _ := ParseMinutes(p1.OpenTime)
			c1, _ := ParseMinutes(p1.CloseTime)
			if c1 <= o1 {
				c1 += 24 * 60 // Spans past midnight
			}

			for j := i + 1; j < len(dayPeriods); j++ {
				p2 := dayPeriods[j]
				if p2.IsClosed || p2.Is24Hours {
					continue
				}
				o2, _ := ParseMinutes(p2.OpenTime)
				c2, _ := ParseMinutes(p2.CloseTime)
				if c2 <= o2 {
					c2 += 24 * 60
				}

				// Check overlap: max(o1, o2) < min(c1, c2)
				startMax := o1
				if o2 > startMax {
					startMax = o2
				}
				endMin := c1
				if c2 < endMin {
					endMin = c2
				}

				if startMax < endMin {
					return fmt.Errorf("overlapping operating periods on %s (%s-%s and %s-%s)",
						DayOfWeekToString(day), p1.OpenTime, p1.CloseTime, p2.OpenTime, p2.CloseTime)
				}
			}
		}
	}
	return nil
}

// SaveRestaurantWeeklySchedule saves the full 7-day schedule inside a database transaction
func (r *Repo) SaveRestaurantWeeklySchedule(ctx context.Context, restaurantID int, schedule []OperatingPeriod) error {
	if r.db == nil {
		return fmt.Errorf("database unavailable")
	}

	// Validate schedule items before writing to DB
	if err := ValidateOperatingPeriods(schedule); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete old schedule for this restaurant
	_, err = tx.ExecContext(ctx, "DELETE FROM restaurant_operating_hours WHERE restaurant_id = $1", restaurantID)
	if err != nil {
		return fmt.Errorf("failed to clear old schedule: %w", err)
	}

	// Insert new schedule entries
	insertSQL := `
		INSERT INTO restaurant_operating_hours (
			restaurant_id, day_of_week, open_time, close_time, is_closed, is_24_hours, period_index, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	for _, p := range schedule {
		periodIdx := p.PeriodIndex
		if periodIdx <= 0 {
			periodIdx = 1
		}
		openT := p.OpenTime
		if openT == "" {
			openT = "09:00"
		}
		closeT := p.CloseTime
		if closeT == "" {
			closeT = "22:00"
		}
		_, err = tx.ExecContext(ctx, insertSQL, restaurantID, p.DayOfWeek, openT, closeT, p.IsClosed, p.Is24Hours, periodIdx)
		if err != nil {
			return fmt.Errorf("failed to insert period: %w", err)
		}
	}

	return tx.Commit()
}

// SetRestaurantTemporaryClosure sets a temporary closure range for a restaurant
func (r *Repo) SetRestaurantTemporaryClosure(ctx context.Context, restaurantID int, start, end time.Time, reason string) error {
	if r.db == nil {
		return fmt.Errorf("database unavailable")
	}
	if end.Before(start) {
		return fmt.Errorf("closure end time cannot be before start time")
	}

	query := `
		INSERT INTO restaurant_special_hours (restaurant_id, closure_start, closure_end, is_closed, reason, created_at, updated_at)
		VALUES ($1, $2, $3, true, $4, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query, restaurantID, start, end, reason)
	if err != nil {
		return fmt.Errorf("failed to set temporary closure: %w", err)
	}
	return nil
}

// ClearRestaurantTemporaryClosure clears any active temporary closures for a restaurant
func (r *Repo) ClearRestaurantTemporaryClosure(ctx context.Context, restaurantID int) error {
	if r.db == nil {
		return fmt.Errorf("database unavailable")
	}
	_, err := r.db.ExecContext(ctx, "DELETE FROM restaurant_special_hours WHERE restaurant_id = $1 AND closure_start IS NOT NULL", restaurantID)
	return err
}

// GetBatchRestaurantAvailability loads hours and calculates availability for multiple restaurants in a single query
func (r *Repo) GetBatchRestaurantAvailability(ctx context.Context, restaurantIDs []int, refTime time.Time) (map[int]AvailabilityResult, error) {
	results := make(map[int]AvailabilityResult)
	if len(restaurantIDs) == 0 {
		return results, nil
	}

	// Default unconfigured restaurants to StatusHoursNotAvailable
	for _, id := range restaurantIDs {
		results[id] = AvailabilityResult{
			Status:     StatusHoursNotAvailable,
			IsOpen:     false,
			BadgeText:  "Hours not available",
			DetailText: "Opening hours not listed",
		}
	}

	if r.db == nil {
		return results, nil
	}

	// Format array parameter for IN query
	idStrings := make([]string, len(restaurantIDs))
	for i, id := range restaurantIDs {
		idStrings[i] = strconv.Itoa(id)
	}
	inClause := strings.Join(idStrings, ",")

	// 1. Fetch operating hours for all restaurants
	hoursQuery := fmt.Sprintf(`
		SELECT id, restaurant_id, day_of_week, open_time, close_time, is_closed, is_24_hours, period_index, created_at, updated_at
		FROM restaurant_operating_hours
		WHERE restaurant_id IN (%s)
		ORDER BY restaurant_id ASC, day_of_week ASC, period_index ASC
	`, inClause)

	hRows, err := r.db.QueryContext(ctx, hoursQuery)
	if err != nil {
		return results, fmt.Errorf("failed to query batch operating hours: %w", err)
	}
	defer hRows.Close()

	periodsMap := make(map[int][]OperatingPeriod)
	for hRows.Next() {
		var p OperatingPeriod
		err := hRows.Scan(
			&p.ID, &p.RestaurantID, &p.DayOfWeek, &p.OpenTime, &p.CloseTime,
			&p.IsClosed, &p.Is24Hours, &p.PeriodIndex, &p.CreatedAt, &p.UpdatedAt,
		)
		if err == nil {
			periodsMap[p.RestaurantID] = append(periodsMap[p.RestaurantID], p)
		}
	}

	// 2. Fetch special hours / closures for all restaurants
	specialsQuery := fmt.Sprintf(`
		SELECT id, restaurant_id, special_date, closure_start, closure_end, is_closed, reason, open_time, close_time, created_at, updated_at
		FROM restaurant_special_hours
		WHERE restaurant_id IN (%s)
		ORDER BY id DESC
	`, inClause)

	sRows, err := r.db.QueryContext(ctx, specialsQuery)
	if err == nil {
		defer sRows.Close()
		specialsMap := make(map[int][]SpecialHours)
		for sRows.Next() {
			var s SpecialHours
			var specDate, cStart, cEnd sql.NullTime
			var reason, openT, closeT sql.NullString
			err := sRows.Scan(
				&s.ID, &s.RestaurantID, &specDate, &cStart, &cEnd, &s.IsClosed,
				&reason, &openT, &closeT, &s.CreatedAt, &s.UpdatedAt,
			)
			if err == nil {
				if specDate.Valid {
					s.SpecialDate = &specDate.Time
				}
				if cStart.Valid {
					s.ClosureStart = &cStart.Time
				}
				if cEnd.Valid {
					s.ClosureEnd = &cEnd.Time
				}
				s.Reason = reason.String
				s.OpenTime = openT.String
				s.CloseTime = closeT.String
				specialsMap[s.RestaurantID] = append(specialsMap[s.RestaurantID], s)
			}
		}

		// Calculate availability for each restaurant in memory
		for _, id := range restaurantIDs {
			periods := periodsMap[id]
			specials := specialsMap[id]
			if len(periods) > 0 {
				results[id] = CalculateAvailability(periods, specials, refTime)
			}
		}
	} else {
		for _, id := range restaurantIDs {
			periods := periodsMap[id]
			if len(periods) > 0 {
				results[id] = CalculateAvailability(periods, nil, refTime)
			}
		}
	}

	return results, nil
}
