package fooddiscovery

import (
	"testing"
	"time"
)

// Helper to create a specific time in Asia/Kolkata (IST = UTC+05:30)
func istTime(year int, month time.Month, day, hour, min int) time.Time {
	loc := GetKolkataLocation()
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func TestAvailabilityService(t *testing.T) {
	// Standard weekly schedule helper: Mon-Sun open 09:00 - 22:00
	makeStandardWeek := func() []OperatingPeriod {
		var list []OperatingPeriod
		for d := 1; d <= 7; d++ {
			list = append(list, OperatingPeriod{
				DayOfWeek:   d,
				OpenTime:    "09:00",
				CloseTime:   "22:00",
				IsClosed:    false,
				Is24Hours:   false,
				PeriodIndex: 1,
			})
		}
		return list
	}

	// 1. Restaurant currently open
	t.Run("1_CurrentlyOpen", func(t *testing.T) {
		periods := makeStandardWeek()
		// Wednesday 14:30 IST
		ref := istTime(2026, time.October, 14, 14, 30)
		res := CalculateAvailability(periods, nil, ref)

		if !res.IsOpen {
			t.Errorf("expected IsOpen=true, got %v", res.IsOpen)
		}
		if res.Status != StatusOpenNow {
			t.Errorf("expected Status=%v, got %v", StatusOpenNow, res.Status)
		}
		if res.DetailText != "Closes at 10:00 PM" {
			t.Errorf("expected detail 'Closes at 10:00 PM', got '%s'", res.DetailText)
		}
	})

	// 2. Restaurant currently closed (late night)
	t.Run("2_CurrentlyClosed", func(t *testing.T) {
		periods := makeStandardWeek()
		// Wednesday 23:15 IST (after 22:00 close)
		ref := istTime(2026, time.October, 14, 23, 15)
		res := CalculateAvailability(periods, nil, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false, got %v", res.IsOpen)
		}
		if res.Status != StatusClosedNow {
			t.Errorf("expected Status=%v, got %v", StatusClosedNow, res.Status)
		}
	})

	// 3. Opens later today
	t.Run("3_OpensLaterToday", func(t *testing.T) {
		periods := makeStandardWeek()
		// Wednesday 07:45 IST (opens at 09:00)
		ref := istTime(2026, time.October, 14, 7, 45)
		res := CalculateAvailability(periods, nil, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false, got %v", res.IsOpen)
		}
		if res.DetailText != "Opens at 9:00 AM" {
			t.Errorf("expected 'Opens at 9:00 AM', got '%s'", res.DetailText)
		}
	})

	// 4. Closes later today
	t.Run("4_ClosesLaterToday", func(t *testing.T) {
		periods := makeStandardWeek()
		// Wednesday 21:30 IST (closes at 22:00)
		ref := istTime(2026, time.October, 14, 21, 30)
		res := CalculateAvailability(periods, nil, ref)

		if !res.IsOpen {
			t.Errorf("expected IsOpen=true, got %v", res.IsOpen)
		}
		if res.DetailText != "Closes at 10:00 PM" {
			t.Errorf("expected 'Closes at 10:00 PM', got '%s'", res.DetailText)
		}
	})

	// 5. Entire day marked closed
	t.Run("5_EntireDayMarkedClosed", func(t *testing.T) {
		periods := []OperatingPeriod{
			{DayOfWeek: 1, OpenTime: "09:00", CloseTime: "22:00", IsClosed: false, PeriodIndex: 1}, // Mon
			{DayOfWeek: 2, OpenTime: "00:00", CloseTime: "00:00", IsClosed: true, PeriodIndex: 1},  // Tue Closed
			{DayOfWeek: 3, OpenTime: "09:00", CloseTime: "22:00", IsClosed: false, PeriodIndex: 1}, // Wed
		}
		// Tuesday 12:00 IST
		ref := istTime(2026, time.October, 13, 12, 0) // Oct 13, 2026 is Tuesday
		res := CalculateAvailability(periods, nil, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false on closed day, got %v", res.IsOpen)
		}
		if res.BadgeText != "Closed Today" {
			t.Errorf("expected BadgeText 'Closed Today', got '%s'", res.BadgeText)
		}
	})

	// 6. Overnight opening period (18:00 to 02:00 next day)
	t.Run("6_OvernightOpeningPeriod", func(t *testing.T) {
		periods := []OperatingPeriod{
			// Friday 18:00 to 02:00 (Saturday morning)
			{DayOfWeek: 5, OpenTime: "18:00", CloseTime: "02:00", IsClosed: false, PeriodIndex: 1},
			// Saturday 18:00 to 02:00 (Sunday morning)
			{DayOfWeek: 6, OpenTime: "18:00", CloseTime: "02:00", IsClosed: false, PeriodIndex: 1},
		}

		// Test A: Friday night at 23:30 (during Friday shift)
		refFriNight := istTime(2026, time.October, 16, 23, 30) // Oct 16 is Friday
		resA := CalculateAvailability(periods, nil, refFriNight)
		if !resA.IsOpen {
			t.Errorf("expected Friday 23:30 to be Open, got %v", resA.IsOpen)
		}
		if resA.DetailText != "Closes at 2:00 AM" {
			t.Errorf("expected Friday 23:30 detail 'Closes at 2:00 AM', got '%s'", resA.DetailText)
		}

		// Test B: Saturday early morning at 01:15 (spillover from Friday night)
		refSatEarly := istTime(2026, time.October, 17, 1, 15) // Oct 17 is Saturday
		resB := CalculateAvailability(periods, nil, refSatEarly)
		if !resB.IsOpen {
			t.Errorf("expected Saturday 01:15 to be Open from Friday overnight, got %v", resB.IsOpen)
		}
		if resB.DetailText != "Closes at 2:00 AM" {
			t.Errorf("expected Saturday 01:15 detail 'Closes at 2:00 AM', got '%s'", resB.DetailText)
		}

		// Test C: Saturday morning at 03:00 (after 02:00 close, before Saturday 18:00 open)
		refSatMorning := istTime(2026, time.October, 17, 3, 0)
		resC := CalculateAvailability(periods, nil, refSatMorning)
		if resC.IsOpen {
			t.Errorf("expected Saturday 03:00 to be Closed, got %v", resC.IsOpen)
		}
		if resC.DetailText != "Opens at 6:00 PM" {
			t.Errorf("expected Saturday 03:00 detail 'Opens at 6:00 PM', got '%s'", resC.DetailText)
		}
	})

	// 7. 24-hour restaurant
	t.Run("7_24HourRestaurant", func(t *testing.T) {
		periods := []OperatingPeriod{
			{DayOfWeek: 1, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 2, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 3, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 4, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 5, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 6, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
			{DayOfWeek: 7, OpenTime: "00:00", CloseTime: "23:59", IsClosed: false, Is24Hours: true, PeriodIndex: 1},
		}
		// Any random time: 03:42 AM
		ref := istTime(2026, time.October, 15, 3, 42)
		res := CalculateAvailability(periods, nil, ref)

		if !res.IsOpen {
			t.Errorf("expected 24h restaurant to be Open, got %v", res.IsOpen)
		}
		if res.BadgeText != "Open 24 Hours" {
			t.Errorf("expected BadgeText 'Open 24 Hours', got '%s'", res.BadgeText)
		}
	})

	// 8. Multiple operating periods (Lunch: 11:30 - 15:00, Dinner: 19:00 - 23:00)
	t.Run("8_MultipleOperatingPeriods", func(t *testing.T) {
		periods := []OperatingPeriod{
			{DayOfWeek: 3, OpenTime: "11:30", CloseTime: "15:00", IsClosed: false, PeriodIndex: 1},
			{DayOfWeek: 3, OpenTime: "19:00", CloseTime: "23:00", IsClosed: false, PeriodIndex: 2},
		}

		// A: During Lunch at 13:00 -> Open, Closes at 3:00 PM
		refLunch := istTime(2026, time.October, 14, 13, 0)
		resLunch := CalculateAvailability(periods, nil, refLunch)
		if !resLunch.IsOpen || resLunch.DetailText != "Closes at 3:00 PM" {
			t.Errorf("lunch period failed: IsOpen=%v, Detail=%s", resLunch.IsOpen, resLunch.DetailText)
		}

		// B: Between shifts at 16:30 -> Closed, Opens at 7:00 PM
		refBetween := istTime(2026, time.October, 14, 16, 30)
		resBetween := CalculateAvailability(periods, nil, refBetween)
		if resBetween.IsOpen || resBetween.DetailText != "Opens at 7:00 PM" {
			t.Errorf("between shifts failed: IsOpen=%v, Detail=%s", resBetween.IsOpen, resBetween.DetailText)
		}

		// C: During Dinner at 20:30 -> Open, Closes at 11:00 PM
		refDinner := istTime(2026, time.October, 14, 20, 30)
		resDinner := CalculateAvailability(periods, nil, refDinner)
		if !resDinner.IsOpen || resDinner.DetailText != "Closes at 11:00 PM" {
			t.Errorf("dinner period failed: IsOpen=%v, Detail=%s", resDinner.IsOpen, resDinner.DetailText)
		}
	})

	// 9. Temporary closure overrides normal hours
	t.Run("9_TemporaryClosureOverride", func(t *testing.T) {
		periods := makeStandardWeek()
		closureStart := istTime(2026, time.October, 14, 10, 0)
		closureEnd := istTime(2026, time.October, 14, 18, 0)
		specials := []SpecialHours{
			{
				ClosureStart: &closureStart,
				ClosureEnd:   &closureEnd,
				IsClosed:     true,
				Reason:       "Renovation work in progress",
			},
		}

		// Current time: 14:00 (normally open, but within temporary closure)
		ref := istTime(2026, time.October, 14, 14, 0)
		res := CalculateAvailability(periods, specials, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false during temp closure, got %v", res.IsOpen)
		}
		if res.Status != StatusTemporarilyClosed {
			t.Errorf("expected Status=%v, got %v", StatusTemporarilyClosed, res.Status)
		}
		if res.Notice != "Renovation work in progress" {
			t.Errorf("expected notice 'Renovation work in progress', got '%s'", res.Notice)
		}

		// After temporary closure at 19:00 (back to normal open hours until 22:00)
		refAfter := istTime(2026, time.October, 14, 19, 0)
		resAfter := CalculateAvailability(periods, specials, refAfter)
		if !resAfter.IsOpen {
			t.Errorf("expected IsOpen=true after temp closure expired, got %v", resAfter.IsOpen)
		}
	})

	// 10. Special-date hours override weekly hours
	t.Run("10_SpecialDateOverride", func(t *testing.T) {
		periods := makeStandardWeek()
		specDate := istTime(2026, time.October, 14, 0, 0)
		specials := []SpecialHours{
			{
				SpecialDate: &specDate,
				IsClosed:    true,
				Reason:      "Closed for Deepavali",
			},
		}

		// Oct 14 at 12:00 IST
		ref := istTime(2026, time.October, 14, 12, 0)
		res := CalculateAvailability(periods, specials, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false on special holiday, got %v", res.IsOpen)
		}
		if res.DetailText != "Closed for Deepavali" {
			t.Errorf("expected DetailText 'Closed for Deepavali', got '%s'", res.DetailText)
		}
	})

	// 11. Missing operating hours
	t.Run("11_MissingOperatingHours", func(t *testing.T) {
		ref := istTime(2026, time.October, 14, 12, 0)
		res := CalculateAvailability(nil, nil, ref)

		if res.IsOpen {
			t.Errorf("expected missing hours IsOpen=false, got %v", res.IsOpen)
		}
		if res.Status != StatusHoursNotAvailable {
			t.Errorf("expected Status=%v, got %v", StatusHoursNotAvailable, res.Status)
		}
		if res.BadgeText != "Hours not available" {
			t.Errorf("expected BadgeText 'Hours not available', got '%s'", res.BadgeText)
		}
	})

	// 12. Exact opening boundary (e.g. exactly at 09:00:00)
	t.Run("12_ExactOpeningBoundary", func(t *testing.T) {
		periods := makeStandardWeek()
		// Exact open boundary: 09:00:00
		ref := istTime(2026, time.October, 14, 9, 0)
		res := CalculateAvailability(periods, nil, ref)

		if !res.IsOpen {
			t.Errorf("expected IsOpen=true at exact open boundary 09:00, got %v", res.IsOpen)
		}
	})

	// 13. Exact closing boundary (e.g. exactly at 22:00:00)
	t.Run("13_ExactClosingBoundary", func(t *testing.T) {
		periods := makeStandardWeek()
		// Exact close boundary: 22:00:00 (customer is past or at closing minute -> closed)
		ref := istTime(2026, time.October, 14, 22, 0)
		res := CalculateAvailability(periods, nil, ref)

		if res.IsOpen {
			t.Errorf("expected IsOpen=false at exact close boundary 22:00, got %v", res.IsOpen)
		}
	})

	// 14. Asia/Kolkata timezone handling
	t.Run("14_AsiaKolkataTimezoneHandling", func(t *testing.T) {
		periods := makeStandardWeek()
		// 04:00 UTC on Wednesday is 09:30 AM IST (UTC+05:30)
		utcTime := time.Date(2026, time.October, 14, 4, 0, 0, 0, time.UTC)
		res := CalculateAvailability(periods, nil, utcTime)

		if !res.IsOpen {
			t.Errorf("expected 04:00 UTC (09:30 IST) to be Open, got %v", res.IsOpen)
		}
		if res.DetailText != "Closes at 10:00 PM" {
			t.Errorf("expected 'Closes at 10:00 PM', got '%s'", res.DetailText)
		}
	})

	// 15. Invalid or contradictory schedules
	t.Run("15_InvalidOrContradictorySchedule", func(t *testing.T) {
		// Period with unparseable hours or malformed times
		periods := []OperatingPeriod{
			{DayOfWeek: 3, OpenTime: "invalid", CloseTime: "garbage", IsClosed: false, PeriodIndex: 1},
		}
		ref := istTime(2026, time.October, 14, 12, 0)
		res := CalculateAvailability(periods, nil, ref)

		// Must safely handle without panicking and return Closed
		if res.IsOpen {
			t.Errorf("expected invalid schedule to not be open, got %v", res.IsOpen)
		}

		// Also verify schedule validator catches overlapping periods
		overlappingPeriods := []OperatingPeriod{
			{DayOfWeek: 1, OpenTime: "09:00", CloseTime: "14:00", PeriodIndex: 1},
			{DayOfWeek: 1, OpenTime: "13:00", CloseTime: "18:00", PeriodIndex: 2}, // Overlaps 13:00-14:00
		}
		err := ValidateOperatingPeriods(overlappingPeriods)
		if err == nil {
			t.Errorf("expected error for overlapping periods, got nil")
		}
	})

	// 16. Weekday and weekend transitions (Sunday night into Monday morning)
	t.Run("16_WeekdayAndWeekendTransition", func(t *testing.T) {
		// Sunday (Day 7) open 20:00 - 02:00 (Monday morning)
		// Monday (Day 1) open 11:00 - 22:00
		periods := []OperatingPeriod{
			{DayOfWeek: 7, OpenTime: "20:00", CloseTime: "02:00", IsClosed: false, PeriodIndex: 1}, // Sunday overnight
			{DayOfWeek: 1, OpenTime: "11:00", CloseTime: "22:00", IsClosed: false, PeriodIndex: 1}, // Monday day
		}

		// Monday morning at 01:15 IST (spillover from Sunday overnight)
		// Oct 18, 2026 is Sunday, Oct 19, 2026 is Monday
		refMonEarly := istTime(2026, time.October, 19, 1, 15)
		resMonEarly := CalculateAvailability(periods, nil, refMonEarly)

		if !resMonEarly.IsOpen {
			t.Errorf("expected Sunday night shift to keep restaurant open on Monday 01:15, got %v", resMonEarly.IsOpen)
		}
		if resMonEarly.DetailText != "Closes at 2:00 AM" {
			t.Errorf("expected 'Closes at 2:00 AM', got '%s'", resMonEarly.DetailText)
		}

		// Monday morning at 04:00 IST (after 02:00 close, before Monday 11:00 open)
		refMonBreak := istTime(2026, time.October, 19, 4, 0)
		resMonBreak := CalculateAvailability(periods, nil, refMonBreak)

		if resMonBreak.IsOpen {
			t.Errorf("expected Monday 04:00 to be Closed, got %v", resMonBreak.IsOpen)
		}
		if resMonBreak.DetailText != "Opens at 11:00 AM" {
			t.Errorf("expected 'Opens at 11:00 AM', got '%s'", resMonBreak.DetailText)
		}
	})
}
