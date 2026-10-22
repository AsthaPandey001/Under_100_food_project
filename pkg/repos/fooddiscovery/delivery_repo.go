package fooddiscovery

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DeliverySettings stores restaurant-managed delivery configurations and operational status
type DeliverySettings struct {
	RestaurantID         int       `json:"restaurant_id"`
	DeliveryEnabled      bool      `json:"delivery_enabled"`       // Does eatery support direct or courier delivery?
	PickupEnabled        bool      `json:"pickup_enabled"`         // Does eatery support counter takeaway/pickup?
	DeliveryRadiusKm     float64   `json:"delivery_radius_km"`     // Maximum delivery radius in kilometers
	DeliveryFee          float64   `json:"delivery_fee"`           // Delivery fee in INR (0.0 for free delivery)
	IsDeliverySuspended  bool      `json:"is_delivery_suspended"`  // Temporary delivery pause (e.g., rain, surge)
	SuspensionReason     string    `json:"suspension_reason"`      // Reason shown to customer when paused
	EstimatedDeliveryMin int       `json:"estimated_delivery_min"` // Estimated delivery duration in minutes
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ServiceabilityStatus represents customer-facing evaluated serviceability for their specific location
type ServiceabilityStatus struct {
	IsServiceable     bool    `json:"is_serviceable"`
	BadgeText         string  `json:"badge_text"`          // "Delivery Available", "Outside Delivery Area", "Delivery Paused", "Takeaway Only"
	DetailText        string  `json:"detail_text"`         // Detailed contextual explanation
	DeliveryFee       float64 `json:"delivery_fee"`
	EstimatedTimeText string  `json:"estimated_time_text"` // "25-35 mins"
	DistanceKm        float64 `json:"distance_km"`
	RadiusKm          float64 `json:"radius_km"`
	PickupEnabled     bool    `json:"pickup_enabled"`
}

// EnsureDeliverySettingsTable creates the delivery settings table and item delivery rules table idempotently
func (r *Repo) EnsureDeliverySettingsTable(ctx context.Context) {
	if r.db == nil {
		return
	}

	ddl := `
	CREATE TABLE IF NOT EXISTS restaurant_delivery_settings (
		restaurant_id INT PRIMARY KEY REFERENCES restaurants(id) ON DELETE CASCADE,
		delivery_enabled BOOLEAN NOT NULL DEFAULT true,
		pickup_enabled BOOLEAN NOT NULL DEFAULT true,
		delivery_radius_km DOUBLE PRECISION NOT NULL DEFAULT 5.0,
		delivery_fee DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		is_delivery_suspended BOOLEAN NOT NULL DEFAULT false,
		suspension_reason VARCHAR(255) DEFAULT '',
		estimated_delivery_min INT NOT NULL DEFAULT 25,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_resto_deliv_enabled ON restaurant_delivery_settings(delivery_enabled);

	CREATE TABLE IF NOT EXISTS food_item_delivery_rules (
		food_item_id INT PRIMARY KEY REFERENCES food_items(id) ON DELETE CASCADE,
		is_deliverable BOOLEAN NOT NULL DEFAULT true,
		delivery_radius_limit DOUBLE PRECISION NOT NULL DEFAULT 0.0,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	_, _ = r.db.ExecContext(ctx, ddl)
}

// DefaultDeliverySettings returns default delivery configurations for a restaurant
func DefaultDeliverySettings(restaurantID int) DeliverySettings {
	return DeliverySettings{
		RestaurantID:         restaurantID,
		DeliveryEnabled:      true,
		PickupEnabled:        true,
		DeliveryRadiusKm:     5.0,
		DeliveryFee:          0.0,
		IsDeliverySuspended:  false,
		SuspensionReason:     "",
		EstimatedDeliveryMin: 25,
	}
}

// CalculateServiceability evaluates whether a customer at distanceKm is serviceable
func CalculateServiceability(settings DeliverySettings, distanceKm float64) ServiceabilityStatus {
	estMin := settings.EstimatedDeliveryMin
	if estMin <= 0 {
		estMin = 25
	}
	timeText := fmt.Sprintf("%d-%d mins", estMin, estMin+10)

	// 1. Delivery completely disabled by merchant (e.g. Dine-in / Takeaway only)
	if !settings.DeliveryEnabled {
		detail := "This eatery offers in-store dining and takeaway only."
		if settings.PickupEnabled {
			detail = "Takeaway available • Delivery not offered by this eatery"
		}
		return ServiceabilityStatus{
			IsServiceable:     false,
			BadgeText:         "Takeaway Only",
			DetailText:        detail,
			DeliveryFee:       0,
			EstimatedTimeText: "",
			DistanceKm:        distanceKm,
			RadiusKm:          settings.DeliveryRadiusKm,
			PickupEnabled:     settings.PickupEnabled,
		}
	}

	// 2. Delivery temporarily suspended by merchant
	if settings.IsDeliverySuspended {
		reason := strings.TrimSpace(settings.SuspensionReason)
		if reason == "" {
			reason = "Delivery temporarily paused due to kitchen rush"
		}
		detail := reason
		if settings.PickupEnabled {
			detail = fmt.Sprintf("%s (Takeaway still available)", reason)
		}
		return ServiceabilityStatus{
			IsServiceable:     false,
			BadgeText:         "Delivery Paused",
			DetailText:        detail,
			DeliveryFee:       settings.DeliveryFee,
			EstimatedTimeText: "",
			DistanceKm:        distanceKm,
			RadiusKm:          settings.DeliveryRadiusKm,
			PickupEnabled:     settings.PickupEnabled,
		}
	}

	// 3. Customer is outside the restaurant's configured delivery radius
	radius := settings.DeliveryRadiusKm
	if radius <= 0 {
		radius = 5.0
	}
	if distanceKm > radius {
		detail := fmt.Sprintf("Delivers up to %.1f km (you are %.1f km away)", radius, distanceKm)
		if settings.PickupEnabled {
			detail += " • Takeaway available"
		}
		return ServiceabilityStatus{
			IsServiceable:     false,
			BadgeText:         "Outside Delivery Area",
			DetailText:        detail,
			DeliveryFee:       settings.DeliveryFee,
			EstimatedTimeText: "",
			DistanceKm:        distanceKm,
			RadiusKm:          radius,
			PickupEnabled:     settings.PickupEnabled,
		}
	}

	// 4. Fully serviceable
	feeText := "Free Delivery"
	if settings.DeliveryFee > 0 {
		feeText = fmt.Sprintf("₹%.0f Delivery Fee", settings.DeliveryFee)
	}
	detail := fmt.Sprintf("Delivers in %s • %s", timeText, feeText)

	return ServiceabilityStatus{
		IsServiceable:     true,
		BadgeText:         "Delivery Available",
		DetailText:        detail,
		DeliveryFee:       settings.DeliveryFee,
		EstimatedTimeText: timeText,
		DistanceKm:        distanceKm,
		RadiusKm:          radius,
		PickupEnabled:     settings.PickupEnabled,
	}
}

// GetRestaurantDeliverySettings retrieves configured settings or returns defaults
func (r *Repo) GetRestaurantDeliverySettings(ctx context.Context, restaurantID int) (DeliverySettings, error) {
	defaults := DefaultDeliverySettings(restaurantID)
	if r.db == nil {
		return defaults, nil
	}

	row := r.db.QueryRowContext(ctx, `
		SELECT restaurant_id, delivery_enabled, pickup_enabled, delivery_radius_km, delivery_fee, 
		       is_delivery_suspended, suspension_reason, estimated_delivery_min, created_at, updated_at
		FROM restaurant_delivery_settings
		WHERE restaurant_id = $1
	`, restaurantID)

	var s DeliverySettings
	err := row.Scan(
		&s.RestaurantID,
		&s.DeliveryEnabled,
		&s.PickupEnabled,
		&s.DeliveryRadiusKm,
		&s.DeliveryFee,
		&s.IsDeliverySuspended,
		&s.SuspensionReason,
		&s.EstimatedDeliveryMin,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return defaults, nil
	}
	if err != nil {
		return defaults, err
	}
	return s, nil
}

// SaveRestaurantDeliverySettings updates or inserts delivery settings for a restaurant
func (r *Repo) SaveRestaurantDeliverySettings(ctx context.Context, restaurantID int, settings DeliverySettings) error {
	if r.db == nil {
		return fmt.Errorf("database connection not available")
	}

	if settings.DeliveryRadiusKm <= 0 {
		settings.DeliveryRadiusKm = 5.0
	}
	if settings.EstimatedDeliveryMin <= 0 {
		settings.EstimatedDeliveryMin = 25
	}
	if settings.DeliveryFee < 0 {
		settings.DeliveryFee = 0.0
	}

	query := `
		INSERT INTO restaurant_delivery_settings 
			(restaurant_id, delivery_enabled, pickup_enabled, delivery_radius_km, delivery_fee,
			 is_delivery_suspended, suspension_reason, estimated_delivery_min, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (restaurant_id) DO UPDATE SET
			delivery_enabled = EXCLUDED.delivery_enabled,
			pickup_enabled = EXCLUDED.pickup_enabled,
			delivery_radius_km = EXCLUDED.delivery_radius_km,
			delivery_fee = EXCLUDED.delivery_fee,
			is_delivery_suspended = EXCLUDED.is_delivery_suspended,
			suspension_reason = EXCLUDED.suspension_reason,
			estimated_delivery_min = EXCLUDED.estimated_delivery_min,
			updated_at = NOW();
	`

	_, err := r.db.ExecContext(ctx, query,
		restaurantID,
		settings.DeliveryEnabled,
		settings.PickupEnabled,
		settings.DeliveryRadiusKm,
		settings.DeliveryFee,
		settings.IsDeliverySuspended,
		strings.TrimSpace(settings.SuspensionReason),
		settings.EstimatedDeliveryMin,
	)
	return err
}

// GetDeliverySettingsBatch loads delivery settings for multiple restaurants in one query
func (r *Repo) GetDeliverySettingsBatch(ctx context.Context, restaurantIDs []int) (map[int]DeliverySettings, error) {
	res := make(map[int]DeliverySettings)
	if len(restaurantIDs) == 0 || r.db == nil {
		return res, nil
	}

	// Prepare placeholder list
	placeholders := make([]string, len(restaurantIDs))
	args := make([]interface{}, len(restaurantIDs))
	for i, id := range restaurantIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
		res[id] = DefaultDeliverySettings(id)
	}

	query := fmt.Sprintf(`
		SELECT restaurant_id, delivery_enabled, pickup_enabled, delivery_radius_km, delivery_fee, 
		       is_delivery_suspended, suspension_reason, estimated_delivery_min, created_at, updated_at
		FROM restaurant_delivery_settings
		WHERE restaurant_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var s DeliverySettings
		if err := rows.Scan(
			&s.RestaurantID,
			&s.DeliveryEnabled,
			&s.PickupEnabled,
			&s.DeliveryRadiusKm,
			&s.DeliveryFee,
			&s.IsDeliverySuspended,
			&s.SuspensionReason,
			&s.EstimatedDeliveryMin,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err == nil {
			res[s.RestaurantID] = s
		}
	}

	return res, nil
}

// FoodItemDeliveryRule defines per-dish delivery availability and radius restriction
type FoodItemDeliveryRule struct {
	FoodItemID          int     `json:"food_item_id"`
	IsDeliverable       bool    `json:"is_deliverable"`        // true = Deliverable, false = Dine-in / Takeaway Only
	DeliveryRadiusLimit float64 `json:"delivery_radius_limit"` // 0.0 = full restaurant radius, >0 = restricted radius (e.g. 3.0 km)
}

// DefaultFoodItemDeliveryRule returns default rule (deliverable across full radius)
func DefaultFoodItemDeliveryRule(foodItemID int) FoodItemDeliveryRule {
	return FoodItemDeliveryRule{
		FoodItemID:          foodItemID,
		IsDeliverable:       true,
		DeliveryRadiusLimit: 0.0,
	}
}

// Evaluate evaluates if a specific dish can be delivered to a customer at distanceKm
func (r FoodItemDeliveryRule) Evaluate(distanceKm, restoRadiusKm float64) (canDeliver bool, badgeText, badgeClass string) {
	if !r.IsDeliverable {
		return false, "Dine-in / Takeaway Only", "badge-takeaway-only"
	}
	if r.DeliveryRadiusLimit > 0.0 {
		if distanceKm > r.DeliveryRadiusLimit {
			return false, fmt.Sprintf("Takeaway only beyond %.0fkm", r.DeliveryRadiusLimit), "badge-radius-limited"
		}
		return true, fmt.Sprintf("Deliverable (≤%.0fkm range)", r.DeliveryRadiusLimit), "badge-deliverable-short"
	}
	if restoRadiusKm > 0.0 && distanceKm > restoRadiusKm {
		return false, "Outside delivery area", "badge-radius-limited"
	}
	return true, "Deliverable", "badge-deliverable"
}

// GetFoodItemDeliveryRule gets delivery rule for a specific food item
func (r *Repo) GetFoodItemDeliveryRule(ctx context.Context, foodItemID int) FoodItemDeliveryRule {
	rule := DefaultFoodItemDeliveryRule(foodItemID)
	if r.db == nil {
		return rule
	}

	query := `SELECT is_deliverable, delivery_radius_limit FROM food_item_delivery_rules WHERE food_item_id = $1`
	_ = r.db.QueryRowContext(ctx, query, foodItemID).Scan(&rule.IsDeliverable, &rule.DeliveryRadiusLimit)
	return rule
}

// SaveFoodItemDeliveryRule saves delivery rule for a specific food item
func (r *Repo) SaveFoodItemDeliveryRule(ctx context.Context, rule FoodItemDeliveryRule) error {
	if r.db == nil {
		return nil
	}

	query := `
		INSERT INTO food_item_delivery_rules (food_item_id, is_deliverable, delivery_radius_limit, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (food_item_id)
		DO UPDATE SET
			is_deliverable = EXCLUDED.is_deliverable,
			delivery_radius_limit = EXCLUDED.delivery_radius_limit,
			updated_at = NOW();
	`
	_, err := r.db.ExecContext(ctx, query, rule.FoodItemID, rule.IsDeliverable, rule.DeliveryRadiusLimit)
	return err
}

// GetFoodItemDeliveryRulesBatch loads delivery rules for multiple food items in one query
func (r *Repo) GetFoodItemDeliveryRulesBatch(ctx context.Context, foodItemIDs []int) (map[int]FoodItemDeliveryRule, error) {
	res := make(map[int]FoodItemDeliveryRule)
	if len(foodItemIDs) == 0 || r.db == nil {
		return res, nil
	}

	placeholders := make([]string, len(foodItemIDs))
	args := make([]interface{}, len(foodItemIDs))
	for i, id := range foodItemIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
		res[id] = DefaultFoodItemDeliveryRule(id)
	}

	query := fmt.Sprintf(`
		SELECT food_item_id, is_deliverable, delivery_radius_limit
		FROM food_item_delivery_rules
		WHERE food_item_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var rule FoodItemDeliveryRule
		if err := rows.Scan(&rule.FoodItemID, &rule.IsDeliverable, &rule.DeliveryRadiusLimit); err == nil {
			res[rule.FoodItemID] = rule
		}
	}
	return res, nil
}

