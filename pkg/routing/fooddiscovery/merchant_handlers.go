package fooddiscovery

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/restaurant"
	"github.com/mikestefanello/pagoda/ent/user"
	foodrepo "github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// MerchantAuthMiddleware ensures only merchants and administrators can access merchant routes
func (h *Handler) MerchantAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := h.c.Auth.GetAuthenticatedUser(c)
		if err != nil || u == nil {
			if c.Request().Header.Get("HX-Request") == "true" {
				c.Response().Header().Set("HX-Redirect", "/merchant/login?msg=merchant_required")
				return c.NoContent(http.StatusOK)
			}
			return c.Redirect(http.StatusSeeOther, "/merchant/login?msg=merchant_required")
		}

		// Only users with role "merchant" or administrators are authorized
		if u.Role != user.RoleMerchant && !u.IsAdmin && u.Role != user.RoleAdmin {
			return c.HTML(http.StatusForbidden, RenderForbiddenPage(u))
		}

		return next(c)
	}
}

// getMerchantRestaurant finds the requested or primary restaurant owned by the merchant user
func (h *Handler) getMerchantRestaurant(ctx context.Context, u *ent.User, requestedID ...int) (*ent.Restaurant, error) {
	if u == nil {
		return nil, nil
	}

	// 1. If a specific restaurant ID was requested (>0):
	if len(requestedID) > 0 && requestedID[0] > 0 {
		reqID := requestedID[0]
		// Platform administrator can inspect any restaurant
		if u.IsAdmin || u.Role == user.RoleAdmin {
			rest, err := h.c.ORM.Restaurant.Query().
				Where(restaurant.IDEQ(reqID)).
				WithCity().
				WithOwner().
				First(ctx)
			if err == nil && rest != nil {
				return rest, nil
			}
		} else {
			// Merchant must own this restaurant
			rest, err := h.c.ORM.Restaurant.Query().
				Where(
					restaurant.IDEQ(reqID),
					restaurant.HasOwnerWith(user.IDEQ(u.ID)),
				).
				WithCity().
				WithOwner().
				First(ctx)
			if err == nil && rest != nil {
				return rest, nil
			}
			return nil, nil
		}
	}

	// 2. Fallback: First restaurant owned by this user
	rest, err := h.c.ORM.Restaurant.Query().
		Where(restaurant.HasOwnerWith(user.IDEQ(u.ID))).
		WithCity().
		WithOwner().
		Order(ent.Asc(restaurant.FieldID)).
		First(ctx)
	if err == nil && rest != nil {
		return rest, nil
	}

	return nil, nil
}

// getAllMerchantRestaurants returns all restaurants owned by the merchant (or all restaurants if administrator)
func (h *Handler) getAllMerchantRestaurants(ctx context.Context, u *ent.User) ([]*ent.Restaurant, error) {
	if u == nil {
		return nil, nil
	}
	if u.IsAdmin || u.Role == user.RoleAdmin {
		return h.c.ORM.Restaurant.Query().
			WithCity().
			WithOwner().
			Order(ent.Asc(restaurant.FieldID)).
			All(ctx)
	}
	return h.c.ORM.Restaurant.Query().
		Where(restaurant.HasOwnerWith(user.IDEQ(u.ID))).
		WithCity().
		WithOwner().
		Order(ent.Asc(restaurant.FieldID)).
		All(ctx)
}

// MerchantDashboard renders the merchant portal overview with multi-outlet switcher support
func (h *Handler) MerchantDashboard(ctx echo.Context) error {
	u, err := h.c.Auth.GetAuthenticatedUser(ctx)
	if err != nil || u == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login?msg=merchant_required")
	}

	reqID, _ := strconv.Atoi(ctx.QueryParam("restaurant_id"))
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}

	rest, _ := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	allRests, _ := h.getAllMerchantRestaurants(ctx.Request().Context(), u)

	if rest != nil {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			sess.Values["merchant_active_restaurant_id"] = rest.ID
			_ = sess.Save(ctx.Request(), ctx.Response())
		}
	}

	var foods []*ent.FoodItem
	var delivSettings foodrepo.DeliverySettings
	itemRules := make(map[int]foodrepo.FoodItemDeliveryRule)
	if rest != nil {
		_, foods, _ = h.c.FoodDiscovery.GetRestaurantByID(ctx.Request().Context(), rest.ID)
		delivSettings, _ = h.c.FoodDiscovery.GetRestaurantDeliverySettings(ctx.Request().Context(), rest.ID)
		if len(foods) > 0 {
			foodIDs := make([]int, len(foods))
			for i, f := range foods {
				foodIDs[i] = f.ID
			}
			if rules, err := h.c.FoodDiscovery.GetFoodItemDeliveryRulesBatch(ctx.Request().Context(), foodIDs); err == nil {
				itemRules = rules
			}
		}
	}

	msg := ctx.QueryParam("msg")
	errMsg := ctx.QueryParam("err")
	html := RenderMerchantDashboard(u, rest, allRests, foods, delivSettings, msg, errMsg, itemRules)
	return ctx.HTML(http.StatusOK, html)
}

// MerchantRestaurantView redirects to main dashboard
func (h *Handler) MerchantRestaurantView(ctx echo.Context) error {
	return ctx.Redirect(http.StatusSeeOther, "/merchant")
}

// MerchantMenuView redirects to main dashboard
func (h *Handler) MerchantMenuView(ctx echo.Context) error {
	return ctx.Redirect(http.StatusSeeOther, "/merchant")
}

// MerchantNewFoodForm returns modal to add a food item
func (h *Handler) MerchantNewFoodForm(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.QueryParam("restaurant_id"))
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">You must have an assigned restaurant before adding menu items.</div>`)
	}

	html := RenderMerchantFoodModal(rest, nil, "/merchant/menu", false)
	return ctx.HTML(http.StatusOK, html)
}

// MerchantCreateFood handles creation of a new food item with strict ownership validation
func (h *Handler) MerchantCreateFood(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="text-red-500">Access denied: no assigned restaurant.</div>`)
	}

	name := strings.TrimSpace(ctx.FormValue("name"))
	price, err := strconv.ParseFloat(strings.TrimSpace(ctx.FormValue("price")), 64)
	if err != nil || price <= 0 {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Please provide a valid positive price. Negative prices are prohibited.</div>`)
	}
	if price > 100.0 {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Under100 Policy: Food items must be priced at or below ₹100.00.</div>`)
	}

	if name == "" {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Dish name is required.</div>`)
	}

	avail := ctx.FormValue("available") == "true" || ctx.FormValue("available") == "on" || ctx.FormValue("available") == "1"

	dietType := ctx.FormValue("diet_type")
	category := strings.TrimSpace(ctx.FormValue("category"))
	if category == "" {
		category = "General"
	}
	if dietType == "nonveg" && !strings.Contains(strings.ToLower(category), "non-veg") {
		category = "Non-Veg " + category
	} else if dietType == "veg" && !strings.Contains(strings.ToLower(category), "veg") {
		category = "Pure Veg " + category
	}

	input := foodrepo.FoodInput{
		Name:        name,
		Description: strings.TrimSpace(ctx.FormValue("description")),
		Category:    category,
		Price:       price,
		Available:   avail,
	}

	createdFood, err := h.c.FoodDiscovery.AddFoodItem(ctx.Request().Context(), rest.ID, input)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, `<div class="p-4 text-red-600">Failed to create dish: `+err.Error()+`</div>`)
	}

	// Assign category slug to new food item (Phase 9)
	rawCategory := strings.TrimSpace(ctx.FormValue("category"))
	if rawCategory != "" && createdFood != nil {
		slug := foodrepo.Slugify(rawCategory)
		_ = h.c.FoodDiscovery.AssignFoodItemCategories(ctx.Request().Context(), createdFood.ID, []string{slug})
	}

	// Persist per-item delivery restrictions (deliverable vs takeaway only, radial limit)
	if createdFood != nil {
		isDeliverable := ctx.FormValue("is_deliverable") == "true" || ctx.FormValue("is_deliverable") == "on" || ctx.FormValue("is_deliverable") == "1"
		radiusLimit, _ := strconv.ParseFloat(ctx.FormValue("delivery_radius_limit"), 64)
		if radiusLimit < 0 {
			radiusLimit = 0
		}
		_ = h.c.FoodDiscovery.SaveFoodItemDeliveryRule(ctx.Request().Context(), foodrepo.FoodItemDeliveryRule{
			FoodItemID:          createdFood.ID,
			IsDeliverable:       isDeliverable,
			DeliveryRadiusLimit: radiusLimit,
		})
	}

	targetURL := fmt.Sprintf("/merchant?restaurant_id=%d&msg=dish_added", rest.ID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", targetURL)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, targetURL)
}

// MerchantEditFoodForm returns modal to edit a dish, strictly verifying ownership
func (h *Handler) MerchantEditFoodForm(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	foodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || foodID <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid food item ID")
	}

	// Strict server-side cross-tenant ownership check
	canManage, rest, err := h.c.FoodDiscovery.CanManageFoodItem(ctx.Request().Context(), u, foodID)
	if err != nil || !canManage || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: You do not have permission to edit this dish.</div>`)
	}

	food, err := h.c.ORM.FoodItem.Get(ctx.Request().Context(), foodID)
	if err != nil {
		return ctx.HTML(http.StatusNotFound, "Dish not found")
	}

	actionURL := fmt.Sprintf("/merchant/menu/%d", foodID)
	rule := h.c.FoodDiscovery.GetFoodItemDeliveryRule(ctx.Request().Context(), foodID)
	html := RenderMerchantFoodModal(rest, food, actionURL, true, rule)
	return ctx.HTML(http.StatusOK, html)
}

// MerchantUpdateFood updates a food item, strictly verifying ownership
func (h *Handler) MerchantUpdateFood(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	foodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || foodID <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid food item ID")
	}

	// Strict server-side cross-tenant ownership check
	canManage, _, err := h.c.FoodDiscovery.CanManageFoodItem(ctx.Request().Context(), u, foodID)
	if err != nil || !canManage {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: You cannot edit a dish belonging to another restaurant.</div>`)
	}

	name := strings.TrimSpace(ctx.FormValue("name"))
	price, err := strconv.ParseFloat(strings.TrimSpace(ctx.FormValue("price")), 64)
	if err != nil || price <= 0 {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Please provide a valid positive price. Negative prices are prohibited.</div>`)
	}
	if price > 100.0 {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Under100 Policy: Food items must be priced at or below ₹100.00.</div>`)
	}

	if name == "" {
		return ctx.HTML(http.StatusBadRequest, `<div class="p-4 text-red-600">Dish name is required.</div>`)
	}

	avail := ctx.FormValue("available") == "true" || ctx.FormValue("available") == "on" || ctx.FormValue("available") == "1"

	dietType := ctx.FormValue("diet_type")
	category := strings.TrimSpace(ctx.FormValue("category"))
	if category == "" {
		category = "General"
	}
	if dietType == "nonveg" && !strings.Contains(strings.ToLower(category), "non-veg") {
		category = "Non-Veg " + category
	} else if dietType == "veg" && !strings.Contains(strings.ToLower(category), "veg") {
		category = "Pure Veg " + category
	}

	input := foodrepo.FoodInput{
		Name:        name,
		Description: strings.TrimSpace(ctx.FormValue("description")),
		Category:    category,
		Price:       price,
		Available:   avail,
	}

	_, err = h.c.FoodDiscovery.UpdateFoodItem(ctx.Request().Context(), foodID, input)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, `<div class="p-4 text-red-600">Failed to update dish: `+err.Error()+`</div>`)
	}

	// Assign category slug to updated food item (Phase 9)
	rawCategory := strings.TrimSpace(ctx.FormValue("category"))
	if rawCategory != "" {
		slug := foodrepo.Slugify(rawCategory)
		_ = h.c.FoodDiscovery.AssignFoodItemCategories(ctx.Request().Context(), foodID, []string{slug})
	}

	// Persist per-item delivery restrictions
	isDeliverable := ctx.FormValue("is_deliverable") == "true" || ctx.FormValue("is_deliverable") == "on" || ctx.FormValue("is_deliverable") == "1"
	radiusLimit, _ := strconv.ParseFloat(ctx.FormValue("delivery_radius_limit"), 64)
	if radiusLimit < 0 {
		radiusLimit = 0
	}
	_ = h.c.FoodDiscovery.SaveFoodItemDeliveryRule(ctx.Request().Context(), foodrepo.FoodItemDeliveryRule{
		FoodItemID:          foodID,
		IsDeliverable:       isDeliverable,
		DeliveryRadiusLimit: radiusLimit,
	})

	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", "/merchant?msg=dish_updated")
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, "/merchant?msg=dish_updated")
}

// MerchantToggleAvailability toggles item availability, strictly verifying ownership
func (h *Handler) MerchantToggleAvailability(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	foodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || foodID <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid food item ID")
	}

	// Strict server-side cross-tenant ownership check
	canManage, _, err := h.c.FoodDiscovery.CanManageFoodItem(ctx.Request().Context(), u, foodID)
	if err != nil || !canManage {
		return ctx.HTML(http.StatusForbidden, "Forbidden")
	}

	updated, err := h.c.FoodDiscovery.ToggleFoodAvailability(ctx.Request().Context(), foodID)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}

	rule := h.c.FoodDiscovery.GetFoodItemDeliveryRule(ctx.Request().Context(), foodID)
	return ctx.HTML(http.StatusOK, RenderMerchantSingleFoodRow(updated, rule))
}

// MerchantDeleteFood deletes a food item, strictly verifying ownership
func (h *Handler) MerchantDeleteFood(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	foodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || foodID <= 0 {
		return ctx.HTML(http.StatusBadRequest, "Invalid food item ID")
	}

	// Strict server-side cross-tenant ownership check
	canManage, _, err := h.c.FoodDiscovery.CanManageFoodItem(ctx.Request().Context(), u, foodID)
	if err != nil || !canManage {
		return ctx.HTML(http.StatusForbidden, "Forbidden")
	}

	err = h.c.FoodDiscovery.DeleteFoodItem(ctx.Request().Context(), foodID)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, err.Error())
	}

	// Return empty string so HTMX cleanly removes the row from the DOM
	return ctx.HTML(http.StatusOK, "")
}

// MerchantEditRestaurantForm renders modal for merchant to edit outlet information
func (h *Handler) MerchantEditRestaurantForm(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.QueryParam("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.FormValue("restaurant_id"))
	}
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: No linked restaurant found.</div>`)
	}

	html := RenderMerchantEditRestaurantModal(rest, "")
	return ctx.HTML(http.StatusOK, html)
}

// MerchantUpdateRestaurant updates restaurant details for the authenticated merchant
func (h *Handler) MerchantUpdateRestaurant(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: No linked restaurant found.</div>`)
	}

	name := strings.TrimSpace(ctx.FormValue("name"))
	desc := strings.TrimSpace(ctx.FormValue("description"))
	phone := strings.TrimSpace(ctx.FormValue("phone"))
	addr := strings.TrimSpace(ctx.FormValue("address"))
	hours := strings.TrimSpace(ctx.FormValue("opening_hours"))
	website := strings.TrimSpace(ctx.FormValue("website"))
	latStr := strings.TrimSpace(ctx.FormValue("latitude"))
	lngStr := strings.TrimSpace(ctx.FormValue("longitude"))

	if name == "" || addr == "" {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantEditRestaurantModal(rest, "Restaurant name and street address are required."))
	}

	updater := h.c.ORM.Restaurant.UpdateOneID(rest.ID).
		SetName(name).
		SetDescription(desc).
		SetAddress(addr).
		SetOpeningHours(hours).
		SetWebsite(website)

	if phone != "" {
		updater = updater.SetPhone(phone)
	}
	if latStr != "" && lngStr != "" {
		lat, errLat := strconv.ParseFloat(latStr, 64)
		lng, errLng := strconv.ParseFloat(lngStr, 64)
		if errLat == nil && errLng == nil && lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180 {
			updater = updater.SetLatitude(lat).SetLongitude(lng)
		}
	}

	_, err = updater.Save(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, RenderMerchantEditRestaurantModal(rest, "Failed to update restaurant: "+err.Error()))
	}

	redirectURL := fmt.Sprintf("/merchant?restaurant_id=%d&msg=restaurant_updated", rest.ID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", redirectURL)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, redirectURL)
}

// MerchantDeliveryView returns the delivery configuration modal
func (h *Handler) MerchantDeliveryView(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.QueryParam("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.FormValue("restaurant_id"))
	}
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: No linked restaurant found.</div>`)
	}

	settings, err := h.c.FoodDiscovery.GetRestaurantDeliverySettings(ctx.Request().Context(), rest.ID)
	if err != nil {
		settings = foodrepo.DeliverySettings{
			RestaurantID:         rest.ID,
			DeliveryEnabled:      true,
			PickupEnabled:        true,
			DeliveryRadiusKm:     5.0,
			DeliveryFee:          30.0,
			EstimatedDeliveryMin: 25,
		}
	}

	html := RenderMerchantDeliveryModal(rest, settings, "", "")
	return ctx.HTML(http.StatusOK, html)
}

// MerchantSaveDeliverySettings updates delivery and serviceability parameters
func (h *Handler) MerchantSaveDeliverySettings(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	reqID, _ := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if reqID <= 0 {
		reqID, _ = strconv.Atoi(ctx.QueryParam("restaurant_id"))
	}
	if reqID <= 0 {
		if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
			if sessRestID, ok := sess.Values["merchant_active_restaurant_id"].(int); ok && sessRestID > 0 {
				reqID = sessRestID
			}
		}
	}
	rest, err := h.getMerchantRestaurant(ctx.Request().Context(), u, reqID)
	if err != nil || rest == nil {
		return ctx.HTML(http.StatusForbidden, `<div class="p-4 text-red-600">Access denied: No linked restaurant found.</div>`)
	}

	delivEnabled := ctx.FormValue("delivery_enabled") == "true" || ctx.FormValue("delivery_enabled") == "on" || ctx.FormValue("delivery_enabled") == "1"
	pickupEnabled := ctx.FormValue("pickup_enabled") == "true" || ctx.FormValue("pickup_enabled") == "on" || ctx.FormValue("pickup_enabled") == "1"
	isSuspended := ctx.FormValue("is_delivery_suspended") == "true" || ctx.FormValue("is_delivery_suspended") == "on" || ctx.FormValue("is_delivery_suspended") == "1"

	radiusStr := strings.TrimSpace(ctx.FormValue("delivery_radius_km"))
	radius, rErr := strconv.ParseFloat(radiusStr, 64)
	if rErr != nil || radius <= 0 || radius > 50 {
		radius = 5.0
	}

	feeStr := strings.TrimSpace(ctx.FormValue("delivery_fee"))
	fee, fErr := strconv.ParseFloat(feeStr, 64)
	if fErr != nil || fee < 0 || fee > 200 {
		fee = 30.0
	}

	estMinStr := strings.TrimSpace(ctx.FormValue("estimated_delivery_min"))
	estMin, mErr := strconv.Atoi(estMinStr)
	if mErr != nil || estMin < 5 || estMin > 180 {
		estMin = 25
	}

	suspReason := strings.TrimSpace(ctx.FormValue("suspension_reason"))

	settings := foodrepo.DeliverySettings{
		RestaurantID:         rest.ID,
		DeliveryEnabled:      delivEnabled,
		PickupEnabled:        pickupEnabled,
		DeliveryRadiusKm:     radius,
		DeliveryFee:          fee,
		IsDeliverySuspended:  isSuspended,
		SuspensionReason:     suspReason,
		EstimatedDeliveryMin: estMin,
	}

	err = h.c.FoodDiscovery.SaveRestaurantDeliverySettings(ctx.Request().Context(), rest.ID, settings)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, RenderMerchantDeliveryModal(rest, settings, "", "Failed to save delivery parameters: "+err.Error()))
	}

	redirectURL := fmt.Sprintf("/merchant?restaurant_id=%d&msg=delivery_updated", rest.ID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", redirectURL)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, redirectURL)
}

// MerchantNewRestaurantForm renders modal for merchants to register a new outlet or branch
func (h *Handler) MerchantNewRestaurantForm(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	if u == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login?msg=merchant_required")
	}

	html := RenderMerchantRestaurantModal("", foodrepo.RestaurantInput{
		CityName:  "Bangalore",
		Latitude:  12.9716,
		Longitude: 77.5946,
	})
	return ctx.HTML(http.StatusOK, html)
}

// MerchantSubmitRestaurant processes self-service restaurant registration by an authenticated merchant
func (h *Handler) MerchantSubmitRestaurant(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	if u == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login?msg=merchant_required")
	}

	// 1. Extract and sanitize input
	name := strings.TrimSpace(ctx.FormValue("name"))
	address := strings.TrimSpace(ctx.FormValue("address"))
	phone := strings.TrimSpace(ctx.FormValue("phone"))
	latStr := strings.TrimSpace(ctx.FormValue("latitude"))
	lngStr := strings.TrimSpace(ctx.FormValue("longitude"))
	hours := strings.TrimSpace(ctx.FormValue("opening_hours"))
	website := strings.TrimSpace(ctx.FormValue("website"))
	desc := strings.TrimSpace(ctx.FormValue("description"))

	var lat, lng float64 = 12.9716, 77.5946
	if latStr != "" || lngStr != "" {
		var errLat, errLng error
		lat, errLat = strconv.ParseFloat(latStr, 64)
		lng, errLng = strconv.ParseFloat(lngStr, 64)
		if errLat != nil || errLng != nil || lat < -90 || lat > 90 || lng < -180 || lng > 180 {
			input := foodrepo.RestaurantInput{
				Name: name, Address: address, Phone: phone, Description: desc,
				OpeningHours: hours, Website: website, Latitude: lat, Longitude: lng,
			}
			return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("Please enter valid geographic coordinates (Latitude: -90 to 90, Longitude: -180 to 180).", input))
		}
	}

	input := foodrepo.RestaurantInput{
		Name:         name,
		Address:      address,
		Phone:        phone,
		Description:  desc,
		OpeningHours: hours,
		Website:      website,
		Latitude:     lat,
		Longitude:    lng,
		CityName:     "Bangalore",
	}

	// 2. Server-side validations
	if len(name) < 2 {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("Restaurant name is required (minimum 2 characters).", input))
	}
	if len(address) < 5 {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("A valid street address and locality is required (minimum 5 characters).", input))
	}
	if phone != "" {
		if len(phone) < 7 || len(phone) > 20 {
			return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("Please enter a valid contact phone number (7 to 20 digits).", input))
		}
		for _, ch := range phone {
			if (ch < '0' || ch > '9') && ch != '+' && ch != '-' && ch != ' ' && ch != '(' && ch != ')' {
				return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("Phone number contains invalid characters.", input))
			}
		}
	}

	// 3. Duplicate prevention: check for existing restaurant with same name and address
	dup, err := h.c.ORM.Restaurant.Query().
		Where(
			restaurant.NameEqualFold(name),
			restaurant.AddressEqualFold(address),
		).Exist(ctx.Request().Context())
	if err == nil && dup {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantRestaurantModal("A restaurant with this name and address has already been submitted or exists on the platform.", input))
	}

	// 4. Create the restaurant with status=pending and owner=u.ID
	newRest, err := h.c.FoodDiscovery.SubmitMerchantRestaurant(ctx.Request().Context(), u.ID, input)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, RenderMerchantRestaurantModal("Failed to register restaurant: "+err.Error(), input))
	}

	// 5. Store active restaurant in session
	if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
		sess.Values["merchant_active_restaurant_id"] = newRest.ID
		_ = sess.Save(ctx.Request(), ctx.Response())
	}

	// 6. Respond with HX-Redirect or HTTP 303 Redirect to merchant dashboard
	targetURL := fmt.Sprintf("/merchant?restaurant_id=%d&msg=restaurant_registered", newRest.ID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", targetURL)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, targetURL)
}

// MerchantClaimRestaurantModal renders modal with unassigned restaurants that can be claimed
func (h *Handler) MerchantClaimRestaurantModal(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	if u == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login?msg=merchant_required")
	}

	claimID, _ := strconv.Atoi(ctx.QueryParam("claim_id"))

	// Query unassigned restaurants in Bangalore
	unassigned, err := h.c.ORM.Restaurant.Query().
		Where(restaurant.Not(restaurant.HasOwner())).
		WithCity().
		Order(ent.Asc(restaurant.FieldName)).
		All(ctx.Request().Context())
	if err != nil {
		unassigned = nil
	}

	html := RenderMerchantClaimModal(unassigned, claimID, "")
	return ctx.HTML(http.StatusOK, html)
}

// MerchantClaimRestaurantSubmit handles claiming an unassigned restaurant
func (h *Handler) MerchantClaimRestaurantSubmit(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	if u == nil {
		return ctx.Redirect(http.StatusSeeOther, "/login?msg=merchant_required")
	}

	restID, err := strconv.Atoi(ctx.FormValue("restaurant_id"))
	if err != nil || restID <= 0 {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantClaimModal(nil, 0, "Please select a valid restaurant to claim."))
	}

	// Verify restaurant exists and has no owner (or current user is admin)
	rest, err := h.c.ORM.Restaurant.Query().
		Where(restaurant.IDEQ(restID)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusNotFound, RenderMerchantClaimModal(nil, restID, "Restaurant not found."))
	}

	if rest.Edges.Owner != nil && (!u.IsAdmin && u.Role != user.RoleAdmin) {
		return ctx.HTML(http.StatusBadRequest, RenderMerchantClaimModal(nil, restID, "This restaurant already has a registered owner. If you represent this eatery, please contact support."))
	}

	// Assign owner
	_, err = h.c.FoodDiscovery.AssignRestaurantOwner(ctx.Request().Context(), restID, u.ID)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, RenderMerchantClaimModal(nil, restID, "Failed to claim restaurant: "+err.Error()))
	}

	// Set active restaurant in session
	if sess, sErr := session.Get("ua", ctx); sErr == nil && sess != nil {
		sess.Values["merchant_active_restaurant_id"] = restID
		_ = sess.Save(ctx.Request(), ctx.Response())
	}

	targetURL := fmt.Sprintf("/merchant?restaurant_id=%d&msg=outlet_claimed", restID)
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", targetURL)
		return ctx.NoContent(http.StatusOK)
	}
	return ctx.Redirect(http.StatusSeeOther, targetURL)
}

