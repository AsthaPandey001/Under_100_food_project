package fooddiscovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/restaurant"
	foodrepo "github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// MerchantHoursView renders the hours management modal for the authenticated merchant's restaurant
func (h *Handler) MerchantHoursView(ctx echo.Context) error {
	u, err := h.c.Auth.GetAuthenticatedUser(ctx)
	if err != nil || u == nil {
		return ctx.HTML(http.StatusUnauthorized, `<div class="p-4 text-red-600">Unauthorized</div>`)
	}

	reqID, _ := strconv.Atoi(ctx.QueryParam("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.FormValue("restaurant_id"))
	}

	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">You do not own or manage an assigned restaurant.</div>`)
	}

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())

	msg := ctx.QueryParam("msg")
	errMsg := ctx.QueryParam("err")
	html := RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, errMsg, msg, false)
	return ctx.HTML(http.StatusOK, html)
}

// MerchantSaveHours updates the weekly schedule with strict ownership validation
func (h *Handler) MerchantSaveHours(ctx echo.Context) error {
	u, err := h.c.Auth.GetAuthenticatedUser(ctx)
	if err != nil || u == nil {
		return ctx.HTML(http.StatusUnauthorized, `<div class="p-4 text-red-600">Unauthorized</div>`)
	}

	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}

	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">You do not own or manage an assigned restaurant.</div>`)
	}

	schedule, err := parseScheduleFromForm(ctx)
	if err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusBadRequest, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, err.Error(), "", false))
	}

	// Validate server-side
	if err := foodrepo.ValidateOperatingPeriods(schedule); err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusBadRequest, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, err.Error(), "", false))
	}

	// Atomic transaction save
	if err := h.c.FoodDiscovery.SaveRestaurantWeeklySchedule(ctx.Request().Context(), rest.ID, schedule); err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusInternalServerError, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, "Failed to save schedule: "+err.Error(), "", false))
	}

	// Re-fetch saved schedule and calculated availability
	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())

	return ctx.HTML(http.StatusOK, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, "", "Weekly operating schedule updated successfully!", false))
}

// MerchantSetTemporaryClosure sets temporary closure range for merchant's restaurant
func (h *Handler) MerchantSetTemporaryClosure(ctx echo.Context) error {
	u, err := h.c.Auth.GetAuthenticatedUser(ctx)
	if err != nil || u == nil {
		return ctx.HTML(http.StatusUnauthorized, `<div class="p-4 text-red-600">Unauthorized</div>`)
	}

	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}

	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">You do not own or manage an assigned restaurant.</div>`)
	}

	startStr := strings.TrimSpace(ctx.FormValue("closure_start"))
	endStr := strings.TrimSpace(ctx.FormValue("closure_end"))
	reason := strings.TrimSpace(ctx.FormValue("reason"))

	if startStr == "" || endStr == "" {
		return h.reRenderHoursModalWithErr(ctx, "/merchant/hours", rest, "Both closure start and end date/time are required.", false)
	}

	loc := foodrepo.GetKolkataLocation()
	layout := "2006-01-02T15:04"
	startTime, errS := time.ParseInLocation(layout, startStr, loc)
	endTime, errE := time.ParseInLocation(layout, endStr, loc)
	if errS != nil || errE != nil {
		return h.reRenderHoursModalWithErr(ctx, "/merchant/hours", rest, "Invalid date/time format. Please use the date picker.", false)
	}

	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return h.reRenderHoursModalWithErr(ctx, "/merchant/hours", rest, "Closure end time must be after the start time.", false)
	}

	if reason == "" {
		reason = "Temporary closure"
	}

	if err := h.c.FoodDiscovery.SetRestaurantTemporaryClosure(ctx.Request().Context(), rest.ID, startTime, endTime, reason); err != nil {
		return h.reRenderHoursModalWithErr(ctx, "/merchant/hours", rest, "Failed to set closure: "+err.Error(), false)
	}

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusOK, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, "", "Temporary closure scheduled successfully!", false))
}

// MerchantClearTemporaryClosure clears active temporary closures
func (h *Handler) MerchantClearTemporaryClosure(ctx echo.Context) error {
	u, err := h.c.Auth.GetAuthenticatedUser(ctx)
	if err != nil || u == nil {
		return ctx.HTML(http.StatusUnauthorized, `<div class="p-4 text-red-600">Unauthorized</div>`)
	}

	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}

	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">You do not own or manage an assigned restaurant.</div>`)
	}

	_ = h.c.FoodDiscovery.ClearRestaurantTemporaryClosure(ctx.Request().Context(), rest.ID)

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusOK, RenderHoursEditorModal("/merchant/hours", rest, periods, specials, avail, "", "Temporary closure cleared.", false))
}

// AdminRestaurantHoursModal renders hours modal for an administrator
func (h *Handler) AdminRestaurantHoursModal(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid restaurant ID")
	}

	rest, err := h.c.ORM.Restaurant.Query().Where(restaurant.IDEQ(id)).Only(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusNotFound, "Restaurant not found")
	}

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())

	actionUrl := fmt.Sprintf("/admin/restaurants/%d/hours", rest.ID)
	html := RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, "", "", true)
	return ctx.HTML(http.StatusOK, html)
}

// AdminSaveRestaurantHours saves hours for any restaurant by platform admin
func (h *Handler) AdminSaveRestaurantHours(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid restaurant ID")
	}

	rest, err := h.c.ORM.Restaurant.Query().Where(restaurant.IDEQ(id)).Only(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusNotFound, "Restaurant not found")
	}

	actionUrl := fmt.Sprintf("/admin/restaurants/%d/hours", rest.ID)
	schedule, err := parseScheduleFromForm(ctx)
	if err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusBadRequest, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, err.Error(), "", true))
	}

	if err := foodrepo.ValidateOperatingPeriods(schedule); err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusBadRequest, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, err.Error(), "", true))
	}

	if err := h.c.FoodDiscovery.SaveRestaurantWeeklySchedule(ctx.Request().Context(), rest.ID, schedule); err != nil {
		periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
		specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
		avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
		return ctx.HTML(http.StatusInternalServerError, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, "Failed to save: "+err.Error(), "", true))
	}

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusOK, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, "", "Operating hours saved successfully!", true))
}

// AdminSetRestaurantTemporaryClosure sets temporary closure by admin
func (h *Handler) AdminSetRestaurantTemporaryClosure(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid restaurant ID")
	}

	rest, err := h.c.ORM.Restaurant.Query().Where(restaurant.IDEQ(id)).Only(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusNotFound, "Restaurant not found")
	}

	actionUrl := fmt.Sprintf("/admin/restaurants/%d/hours", rest.ID)
	startStr := strings.TrimSpace(ctx.FormValue("closure_start"))
	endStr := strings.TrimSpace(ctx.FormValue("closure_end"))
	reason := strings.TrimSpace(ctx.FormValue("reason"))

	if startStr == "" || endStr == "" {
		return h.reRenderHoursModalWithErr(ctx, actionUrl, rest, "Start and end dates are required.", true)
	}

	loc := foodrepo.GetKolkataLocation()
	layout := "2006-01-02T15:04"
	startTime, errS := time.ParseInLocation(layout, startStr, loc)
	endTime, errE := time.ParseInLocation(layout, endStr, loc)
	if errS != nil || errE != nil {
		return h.reRenderHoursModalWithErr(ctx, actionUrl, rest, "Invalid date format.", true)
	}

	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return h.reRenderHoursModalWithErr(ctx, actionUrl, rest, "Closure end time must be after start time.", true)
	}

	if reason == "" {
		reason = "Administrative temporary closure"
	}

	_ = h.c.FoodDiscovery.SetRestaurantTemporaryClosure(ctx.Request().Context(), rest.ID, startTime, endTime, reason)

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusOK, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, "", "Temporary closure scheduled successfully!", true))
}

// AdminClearRestaurantTemporaryClosure clears closure by admin
func (h *Handler) AdminClearRestaurantTemporaryClosure(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid restaurant ID")
	}

	rest, err := h.c.ORM.Restaurant.Query().Where(restaurant.IDEQ(id)).Only(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusNotFound, "Restaurant not found")
	}

	actionUrl := fmt.Sprintf("/admin/restaurants/%d/hours", rest.ID)
	_ = h.c.FoodDiscovery.ClearRestaurantTemporaryClosure(ctx.Request().Context(), rest.ID)

	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusOK, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, "", "Temporary closure cleared.", true))
}

func (h *Handler) reRenderHoursModalWithErr(ctx echo.Context, actionUrl string, rest *ent.Restaurant, errMsg string, isAdmin bool) error {
	periods, _ := h.c.FoodDiscovery.GetRestaurantOperatingHours(ctx.Request().Context(), rest.ID)
	specials, _ := h.c.FoodDiscovery.GetRestaurantSpecialHours(ctx.Request().Context(), rest.ID)
	avail := foodrepo.CalculateAvailability(periods, specials, time.Now())
	return ctx.HTML(http.StatusBadRequest, RenderHoursEditorModal(actionUrl, rest, periods, specials, avail, errMsg, "", isAdmin))
}

// parseScheduleFromForm extracts OperatingPeriod slices for days 1 to 7 from form fields or JSON
func parseScheduleFromForm(ctx echo.Context) ([]foodrepo.OperatingPeriod, error) {
	jsonStr := strings.TrimSpace(ctx.FormValue("periods_json"))
	if jsonStr != "" {
		var jsonPeriods []foodrepo.OperatingPeriod
		if err := json.Unmarshal([]byte(jsonStr), &jsonPeriods); err != nil {
			return nil, fmt.Errorf("invalid periods JSON: %w", err)
		}
		return jsonPeriods, nil
	}

	var schedule []foodrepo.OperatingPeriod

	for day := 1; day <= 7; day++ {
		prefix := fmt.Sprintf("day_%d", day)
		closedVal := ctx.FormValue(prefix + "_closed")
		isClosed := closedVal == "on" || closedVal == "true" || closedVal == "1"

		is24hVal := ctx.FormValue(prefix + "_24h")
		is24h := is24hVal == "on" || is24hVal == "true" || is24hVal == "1"

		openTime := strings.TrimSpace(ctx.FormValue(prefix + "_open"))
		closeTime := strings.TrimSpace(ctx.FormValue(prefix + "_close"))

		if isClosed {
			schedule = append(schedule, foodrepo.OperatingPeriod{
				DayOfWeek:   day,
				OpenTime:    "00:00",
				CloseTime:   "00:00",
				IsClosed:    true,
				Is24Hours:   false,
				PeriodIndex: 1,
			})
			continue
		}

		if is24h {
			schedule = append(schedule, foodrepo.OperatingPeriod{
				DayOfWeek:   day,
				OpenTime:    "00:00",
				CloseTime:   "23:59",
				IsClosed:    false,
				Is24Hours:   true,
				PeriodIndex: 1,
			})
			continue
		}

		if openTime == "" {
			openTime = "09:00"
		}
		if closeTime == "" {
			closeTime = "22:00"
		}

		schedule = append(schedule, foodrepo.OperatingPeriod{
			DayOfWeek:   day,
			OpenTime:    openTime,
			CloseTime:   closeTime,
			IsClosed:    false,
			Is24Hours:   false,
			PeriodIndex: 1,
		})

		// Check optional second period
		openTime2 := strings.TrimSpace(ctx.FormValue(prefix + "_open_2"))
		closeTime2 := strings.TrimSpace(ctx.FormValue(prefix + "_close_2"))
		hasPeriod2Val := ctx.FormValue(prefix + "_has_period_2")
		if (hasPeriod2Val == "on" || hasPeriod2Val == "true" || hasPeriod2Val == "1") && openTime2 != "" && closeTime2 != "" {
			schedule = append(schedule, foodrepo.OperatingPeriod{
				DayOfWeek:   day,
				OpenTime:    openTime2,
				CloseTime:   closeTime2,
				IsClosed:    false,
				Is24Hours:   false,
				PeriodIndex: 2,
			})
		}
	}

	return schedule, nil
}

// RenderHoursEditorModal renders the HTML modal for managing restaurant weekly hours and temporary closures
func RenderHoursEditorModal(
	actionUrl string,
	rest *ent.Restaurant,
	periods []foodrepo.OperatingPeriod,
	specials []foodrepo.SpecialHours,
	avail foodrepo.AvailabilityResult,
	errMsg, successMsg string,
	isAdmin bool,
) string {
	targetModal := "#merchant-modal"
	if isAdmin {
		targetModal = "#admin-modal"
	}

	var alertHTML string
	if errMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #fef2f2; border: 1.5px solid #fecaca; color: #991b1b; padding: 0.75rem 1rem; border-radius: 8px; margin-bottom: 1rem; font-size: 0.9rem; font-weight: 600;">%s</div>`, escapeHTML(errMsg))
	} else if successMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #f0fdf4; border: 1.5px solid #bbf7d0; color: #166534; padding: 0.75rem 1rem; border-radius: 8px; margin-bottom: 1rem; font-size: 0.9rem; font-weight: 600;">%s</div>`, escapeHTML(successMsg))
	}

	// Map configured periods by DayOfWeek and PeriodIndex
	pMap := make(map[int]map[int]foodrepo.OperatingPeriod)
	for _, p := range periods {
		if pMap[p.DayOfWeek] == nil {
			pMap[p.DayOfWeek] = make(map[int]foodrepo.OperatingPeriod)
		}
		pMap[p.DayOfWeek][p.PeriodIndex] = p
	}

	statusBadgeClass := "unavailable"
	switch avail.Status {
	case foodrepo.StatusOpenNow:
		statusBadgeClass = "open"
	case foodrepo.StatusClosedNow, foodrepo.StatusOpensAt:
		statusBadgeClass = "closed"
	case foodrepo.StatusTemporarilyClosed:
		statusBadgeClass = "temp-closed"
	}

	currentStatusHTML := fmt.Sprintf(`
		<div style="background: #f8fafc; border: 1.5px solid var(--border-gray, #e2e8f0); border-radius: 10px; padding: 1rem; margin-bottom: 1.25rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
			<div>
				<div style="font-size: 0.8rem; font-weight: 700; color: var(--color-secondary-text, #64748b); text-transform: uppercase;">Current Live Status (Asia/Kolkata)</div>
				<div style="font-size: 1.05rem; font-weight: 800; color: var(--color-dark-text, #171c24); margin-top: 2px;">
					%s
				</div>
			</div>
			<div class="hours-badge %s" style="font-size: 13px; padding: 6px 12px;">
				<span class="badge-dot-live"></span>
				%s
			</div>
		</div>
	`, escapeHTML(avail.DetailText), statusBadgeClass, escapeHTML(avail.BadgeText))

	// Active Temporary Closure Banner
	var tempClosureHTML strings.Builder
	var activeClosure *foodrepo.SpecialHours
	now := time.Now()
	for i, spec := range specials {
		if spec.IsClosed && spec.ClosureStart != nil && spec.ClosureEnd != nil {
			if now.Before(*spec.ClosureEnd) {
				activeClosure = &specials[i]
				break
			}
		}
	}

	clearClosureUrl := actionUrl + "/temporary-closure/clear"
	setClosureUrl := actionUrl + "/temporary-closure"

	if activeClosure != nil {
		loc := foodrepo.GetKolkataLocation()
		startFmt := activeClosure.ClosureStart.In(loc).Format("Jan 02, 3:04 PM")
		endFmt := activeClosure.ClosureEnd.In(loc).Format("Jan 02, 3:04 PM")
		tempClosureHTML.WriteString(fmt.Sprintf(`
			<div style="background: #fffbeb; border: 1.5px solid #fde68a; border-radius: 10px; padding: 1rem; margin-bottom: 1.25rem;">
				<div style="display: flex; justify-content: space-between; align-items: flex-start; gap: 0.5rem;">
					<div>
						<div style="font-size: 0.82rem; font-weight: 800; color: #92400e; text-transform: uppercase;">⚠️ Active Temporary Closure</div>
						<div style="font-size: 0.95rem; font-weight: 700; color: #78350f; margin-top: 3px;">%s</div>
						<div style="font-size: 0.82rem; color: #92400e; margin-top: 2px;">%s &mdash; %s (IST)</div>
					</div>
					<button 
						type="button" 
						hx-post="%s" 
						hx-target="%s" 
						style="background: #ffffff; border: 1.5px solid #d97706; color: #b45309; padding: 4px 10px; border-radius: 6px; font-size: 0.8rem; font-weight: 700; cursor: pointer;"
					>
						Clear Closure
					</button>
				</div>
			</div>
		`, escapeHTML(activeClosure.Reason), startFmt, endFmt, clearClosureUrl, targetModal))
	}

	// Build Weekly Schedule rows
	var daysRows strings.Builder
	for day := 1; day <= 7; day++ {
		dayName := foodrepo.DayOfWeekToString(day)
		p1, hasP1 := pMap[day][1]
		p2, hasP2 := pMap[day][2]

		isClosed := hasP1 && p1.IsClosed
		is24h := hasP1 && p1.Is24Hours
		open1 := "09:00"
		close1 := "22:00"
		if hasP1 && !p1.IsClosed && !p1.Is24Hours {
			open1 = p1.OpenTime
			close1 = p1.CloseTime
		}

		open2 := "18:00"
		close2 := "23:00"
		if hasP2 {
			open2 = p2.OpenTime
			close2 = p2.CloseTime
		}

		closedChecked := ""
		if isClosed {
			closedChecked = "checked"
		}
		is24hChecked := ""
		if is24h {
			is24hChecked = "checked"
		}
		hasP2Checked := ""
		if hasP2 {
			hasP2Checked = "checked"
		}

		timesDisplay := ""
		if isClosed || is24h {
			timesDisplay = "display: none;"
		}
		period2Display := "none"
		if hasP2 {
			period2Display = "flex"
		}

		daysRows.WriteString(fmt.Sprintf(`
			<div class="schedule-day-row" id="row_day_%d" style="background: #ffffff; border: 1.5px solid var(--border-gray, #e2e8f0); border-radius: 10px; padding: 0.9rem; margin-bottom: 0.75rem;">
				<div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 0.6rem;">
					<div style="font-weight: 800; font-size: 0.98rem; color: var(--color-dark-text, #171c24); min-width: 100px;">
						%s
					</div>
					<div style="display: flex; align-items: center; gap: 1rem; font-size: 0.85rem;">
						<label style="display: inline-flex; align-items: center; gap: 4px; cursor: pointer; color: #dc2626; font-weight: 700;">
							<input type="checkbox" name="day_%d_closed" value="true" %s onchange="toggleDayControls(%d)">
							Closed
						</label>
						<label style="display: inline-flex; align-items: center; gap: 4px; cursor: pointer; color: #087443; font-weight: 700;">
							<input type="checkbox" name="day_%d_24h" value="true" %s onchange="toggleDayControls(%d)">
							24 Hours
						</label>
					</div>
				</div>

				<div id="day_%d_times" style="display: flex; flex-direction: column; gap: 0.5rem; %s">
					<div style="display: flex; align-items: center; gap: 0.6rem; flex-wrap: wrap;">
						<span style="font-size: 0.8rem; font-weight: 700; color: #64748b; width: 60px;">Period 1:</span>
						<label style="font-size: 0.8rem; color: #64748b;">Open:</label>
						<input type="time" name="day_%d_open" value="%s" style="padding: 4px 8px; border: 1.5px solid #cbd5e1; border-radius: 6px; font-weight: 600; font-size: 0.88rem;">
						<label style="font-size: 0.8rem; color: #64748b;">Close:</label>
						<input type="time" name="day_%d_close" value="%s" style="padding: 4px 8px; border: 1.5px solid #cbd5e1; border-radius: 6px; font-weight: 600; font-size: 0.88rem;">
						<span style="font-size: 0.75rem; color: #94a3b8;">(Past midnight supported e.g. 02:00)</span>
					</div>

					<div style="display: flex; align-items: center; gap: 0.6rem; flex-wrap: wrap; margin-top: 4px;">
						<label style="font-size: 0.8rem; font-weight: 700; color: #475569; display: inline-flex; align-items: center; gap: 4px; cursor: pointer;">
							<input type="checkbox" name="day_%d_has_period_2" value="true" %s onchange="document.getElementById('day_%d_period2_inputs').style.display = this.checked ? 'flex' : 'none';">
							Add 2nd Shift (e.g. Dinner)
						</label>
						<div id="day_%d_period2_inputs" style="display: %s; align-items: center; gap: 0.5rem; flex-wrap: wrap;">
							<label style="font-size: 0.8rem; color: #64748b;">Open:</label>
							<input type="time" name="day_%d_open_2" value="%s" style="padding: 4px 8px; border: 1.5px solid #cbd5e1; border-radius: 6px; font-weight: 600; font-size: 0.88rem;">
							<label style="font-size: 0.8rem; color: #64748b;">Close:</label>
							<input type="time" name="day_%d_close_2" value="%s" style="padding: 4px 8px; border: 1.5px solid #cbd5e1; border-radius: 6px; font-weight: 600; font-size: 0.88rem;">
						</div>
					</div>
				</div>
			</div>
		`, day, dayName, day, closedChecked, day, day, is24hChecked, day, day,
			timesDisplay,
			day, open1, day, close1,
			day, hasP2Checked, day,
			day, period2Display,
			day, open2, day, close2,
		))
	}

	return fmt.Sprintf(`
		<div class="modal-overlay" onclick="if(event.target === this) { document.querySelector('%s').innerHTML = ''; }">
			<div class="modal-card" style="max-width: 660px; max-height: 90vh; overflow-y: auto;">
				<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
					<div>
						<h2 style="font-family: var(--font-heading, sans-serif); font-size: 1.45rem; font-weight: 800; color: var(--color-dark-text, #171c24); margin: 0;">
							Operating Hours & Availability
						</h2>
						<p style="color: var(--color-secondary-text, #64748b); font-size: 0.88rem; margin-top: 3px;">
							%s &bull; Timezone: Asia/Kolkata (IST)
						</p>
					</div>
					<button type="button" onclick="document.querySelector('%s').innerHTML = '';" style="background: none; border: none; font-size: 1.5rem; color: #94a3b8; cursor: pointer;">&times;</button>
				</div>

				%s
				%s
				%s

				<!-- Weekly Schedule Form -->
				<form hx-post="%s" hx-target="%s" style="margin-bottom: 1.75rem;">
					<input type="hidden" name="restaurant_id" value="%d">
					<div style="font-size: 0.95rem; font-weight: 800; color: var(--color-dark-text, #171c24); margin-bottom: 0.75rem;">
						Weekly Operating Hours
					</div>

					%s

					<div style="margin-top: 1.25rem; display: flex; justify-content: flex-end; gap: 0.75rem;">
						<button type="button" class="btn-secondary" onclick="document.querySelector('%s').innerHTML = '';" style="padding: 8px 18px; border-radius: 8px; border: 1.5px solid #cbd5e1; font-weight: 700; cursor: pointer;">
							Cancel
						</button>
						<button type="submit" class="btn-primary" style="background: var(--color-primary-orange, #ff5a1f); color: #ffffff; padding: 8px 22px; border-radius: 8px; font-weight: 800; border: none; cursor: pointer;">
							Save Weekly Schedule
						</button>
					</div>
				</form>

				<!-- Temporary Closure Form -->
				<div style="background: #fafaf9; border: 1.5px solid #e7e5e4; border-radius: 12px; padding: 1.25rem; margin-top: 1.5rem;">
					<div style="font-weight: 800; font-size: 1rem; color: #44403c; margin-bottom: 0.4rem;">
						Schedule Temporary Closure
					</div>
					<p style="font-size: 0.83rem; color: #78716c; margin-bottom: 1rem;">
						Need to close for maintenance, private events, or holidays? Set a temporary closure window. Normal hours resume automatically after the end time.
					</p>

					<form hx-post="%s" hx-target="%s" style="display: flex; flex-direction: column; gap: 0.75rem;">
						<input type="hidden" name="restaurant_id" value="%d">
						<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem;">
							<div>
								<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #57534e; margin-bottom: 4px;">Closure Starts</label>
								<input type="datetime-local" name="closure_start" required style="width: 100%%; padding: 6px 10px; border: 1.5px solid #d6d3d1; border-radius: 6px; font-size: 0.85rem;">
							</div>
							<div>
								<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #57534e; margin-bottom: 4px;">Closure Ends / Reopens</label>
								<input type="datetime-local" name="closure_end" required style="width: 100%%; padding: 6px 10px; border: 1.5px solid #d6d3d1; border-radius: 6px; font-size: 0.85rem;">
							</div>
						</div>
						<div>
							<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #57534e; margin-bottom: 4px;">Customer Notice / Reason</label>
							<input type="text" name="reason" placeholder="e.g. Closed for kitchen renovation" style="width: 100%%; padding: 6px 10px; border: 1.5px solid #d6d3d1; border-radius: 6px; font-size: 0.85rem;">
						</div>
						<div style="display: flex; justify-content: flex-end; margin-top: 0.25rem;">
							<button type="submit" style="background: #78716c; color: #ffffff; border: none; padding: 7px 16px; border-radius: 6px; font-size: 0.85rem; font-weight: 700; cursor: pointer;">
								Set Temporary Closure
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>

		<script>
			function toggleDayControls(day) {
				var row = document.getElementById('row_day_' + day);
				var timesDiv = document.getElementById('day_' + day + '_times');
				if (!row || !timesDiv) return;
				var closedCheck = row.querySelector('input[name="day_' + day + '_closed"]');
				var is24hCheck = row.querySelector('input[name="day_' + day + '_24h"]');
				if ((closedCheck && closedCheck.checked) || (is24hCheck && is24hCheck.checked)) {
					timesDiv.style.display = 'none';
				} else {
					timesDiv.style.display = 'flex';
				}
			}
		</script>
	`, targetModal, escapeHTML(rest.Name), targetModal, alertHTML, currentStatusHTML, tempClosureHTML.String(),
		actionUrl, targetModal, rest.ID, daysRows.String(), targetModal, setClosureUrl, targetModal, rest.ID)
}
