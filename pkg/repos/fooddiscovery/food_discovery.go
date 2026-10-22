package fooddiscovery

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/city"
	"github.com/mikestefanello/pagoda/ent/fooditem"
	"github.com/mikestefanello/pagoda/ent/restaurant"
	"github.com/mikestefanello/pagoda/ent/user"
)

type Repo struct {
	client *ent.Client
	db     *sql.DB
}

func NewRepo(client *ent.Client, db *sql.DB) *Repo {
	r := &Repo{client: client, db: db}
	r.EnsureFavoritesTable(context.Background())
	r.EnsureReviewsTable(context.Background())
	r.EnsureCategoriesTable(context.Background())
	r.EnsureOpeningHoursTables(context.Background())
	r.EnsureDeliverySettingsTable(context.Background())
	return r
}

type RestaurantWithDistance struct {
	Restaurant      *ent.Restaurant      `json:"restaurant"`
	DistanceKm      float64              `json:"distance_km"`
	DistanceDisplay string               `json:"distance_display"`
	Under100Items   []*ent.FoodItem      `json:"under_100_items"`
	AllItems        []*ent.FoodItem      `json:"all_items,omitempty"`
	MinPrice        float64              `json:"min_price"`
	IsPureVeg       bool                 `json:"is_pure_veg"`
	HasNonVeg       bool                 `json:"has_non_veg"`
	IsFavorite      bool                 `json:"is_favorite"`
	AverageRating   float64              `json:"average_rating"`
	ReviewCount     int                  `json:"review_count"`
	Availability    AvailabilityResult   `json:"availability"`
	Serviceability  ServiceabilityStatus `json:"serviceability"`
}

type Review struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	RestaurantID     int       `json:"restaurant_id"`
	Rating           int       `json:"rating"`
	Comment          string    `json:"comment"`
	Status           string    `json:"status"` // "approved", "hidden"
	ModerationReason string    `json:"moderation_reason"`
	UserName         string    `json:"user_name"`
	UserEmail        string    `json:"user_email,omitempty"`
	RestaurantName   string    `json:"restaurant_name,omitempty"`
	ReportCount      int       `json:"report_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ReviewReport struct {
	ID        int       `json:"id"`
	ReviewID  int       `json:"review_id"`
	UserID    int       `json:"user_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type RestaurantRatingStats struct {
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`
}

type PaginatedReviews struct {
	Items      []Review `json:"items"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
	HasNext    bool     `json:"has_next"`
	HasPrev    bool     `json:"has_prev"`
}

type PaginatedRestaurants struct {
	Items          []RestaurantWithDistance `json:"items"`
	TotalCount     int                      `json:"total_count"`
	Page           int                      `json:"page"`
	PageSize       int                      `json:"page_size"`
	TotalPages     int                      `json:"total_pages"`
	HasNext        bool                     `json:"has_next"`
	HasPrev        bool                     `json:"has_prev"`
	SearchQuery    string                   `json:"search_query,omitempty"`
	DietFilter     string                   `json:"diet_filter,omitempty"`
	CategoryFilter string                   `json:"category_filter,omitempty"`
	OpenNowFilter  bool                     `json:"open_now_filter,omitempty"`
	MaxPrice       float64                  `json:"max_price,omitempty"`
	SortBy         string                   `json:"sort_by,omitempty"`
}

// IsDishVeg accurately classifies whether a dish is Vegetarian or Non-Vegetarian
func IsDishVeg(name, category, description string) bool {
	catLower := strings.ToLower(category)
	if strings.Contains(catLower, "non-veg") || strings.Contains(catLower, "non veg") {
		return false
	}
	if strings.Contains(catLower, "pure veg") || catLower == "veg" {
		return true
	}

	text := " " + strings.ToLower(name+" "+description) + " "

	// Explicit veg exceptions first
	if strings.Contains(text, "eggless") || strings.Contains(text, "veg biryani") || strings.Contains(text, "paneer biryani") || strings.Contains(text, "mushroom biryani") {
		return true
	}

	nonVegKeywords := []string{
		"chicken", "mutton", "egg", "fish", "prawn", "seafood", "beef", "pork",
		"meat", "kabab", "kebab", "keema", "shawarma", "kothu", "sheekh", "wings",
		"omelette", "bhurji", "brain",
	}

	for _, kw := range nonVegKeywords {
		if strings.Contains(text, " "+kw+" ") ||
			strings.Contains(text, " "+kw+"s ") ||
			strings.Contains(text, " "+kw+",") ||
			strings.Contains(text, " "+kw+".") ||
			strings.Contains(text, "("+kw) ||
			strings.Contains(text, kw+")") {
			return false
		}
	}

	// Biryani check (non-veg by default unless explicitly tagged veg/paneer/mushroom above)
	if strings.Contains(text, "biryani") {
		return false
	}

	return true
}

type AdminStats struct {
	TotalRestaurants       int `json:"total_restaurants"`
	PendingRestaurants     int `json:"pending_restaurants"`
	ApprovedRestaurants    int `json:"approved_restaurants"`
	RejectedRestaurants    int `json:"rejected_restaurants"`
	UnassignedRestaurants  int `json:"unassigned_restaurants"`
	TotalFoodItems         int `json:"total_food_items"`
	Under100FoodItems      int `json:"under_100_food_items"`
	TotalUsers             int `json:"total_users"`
	CustomerCount          int `json:"customer_count"`
	MerchantCount          int `json:"merchant_count"`
	AdminCount             int `json:"admin_count"`
}

type RestaurantInput struct {
	Name         string  `json:"name" form:"name"`
	Description  string  `json:"description" form:"description"`
	Address      string  `json:"address" form:"address"`
	CityName     string  `json:"city_name" form:"city_name"`
	Latitude     float64 `json:"latitude" form:"latitude"`
	Longitude    float64 `json:"longitude" form:"longitude"`
	Phone        string  `json:"phone" form:"phone"`
	Website      string  `json:"website" form:"website"`
	OpeningHours string  `json:"opening_hours" form:"opening_hours"`
	Status       string  `json:"status" form:"status"`
}

type FoodInput struct {
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	Category    string  `json:"category" form:"category"`
	Price       float64 `json:"price" form:"price"`
	Available   bool    `json:"available" form:"available"`
}

// CalculateHaversineDistance calculates spherical distance between two coordinates in kilometers
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	rLat1 := lat1 * (math.Pi / 180.0)
	rLat2 := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// FormatDistance returns a human-friendly string for distance
func FormatDistance(km float64) string {
	if km < 1.0 {
		return fmt.Sprintf("%d m away", int(km*1000))
	}
	return fmt.Sprintf("%.1f km away", km)
}

// GetNearbyRestaurants searches approved restaurants nearby that have at least one available food item <= maxPrice (default ₹100), filtered by diet (all, veg, nonveg)
func (r *Repo) GetNearbyRestaurants(
	ctx context.Context,
	userLat, userLng float64,
	radiusKm float64,
	maxPrice float64,
	searchQuery string,
	dietFilter string,
	sortBy string,
) ([]RestaurantWithDistance, error) {
	paginated, err := r.GetNearbyRestaurantsPaginated(ctx, userLat, userLng, radiusKm, maxPrice, searchQuery, dietFilter, "", false, sortBy, 1, 1000, nil)
	if err != nil {
		return nil, err
	}
	return paginated.Items, nil
}

// GetNearbyRestaurantsPaginated searches and paginates approved restaurants with food <= maxPrice matching criteria
func (r *Repo) GetNearbyRestaurantsPaginated(
	ctx context.Context,
	userLat, userLng float64,
	radiusKm float64,
	maxPrice float64,
	searchQuery string,
	dietFilter string,
	categoryFilter string,
	openNowFilter bool,
	sortBy string,
	page int,
	pageSize int,
	userFavIDs map[int]bool,
) (*PaginatedRestaurants, error) {
	if maxPrice <= 0 || maxPrice > 100 {
		maxPrice = 100.0
	}
	if radiusKm <= 0 {
		radiusKm = 25.0 // 25km radius default for metropolitan Bangalore
	}

	// Pagination parameters normalization
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 50 {
		pageSize = 50
	}

	// 1. Fetch only APPROVED restaurants with their food items
	query := r.client.Restaurant.Query().
		Where(restaurant.StatusEQ(restaurant.StatusApproved)).
		WithFoodItems().
		WithCity()

	restaurants, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query restaurants: %w", err)
	}

	cleanSearch := strings.TrimSpace(strings.ToLower(searchQuery))
	cleanDiet := strings.TrimSpace(strings.ToLower(dietFilter))
	cleanCategory := strings.TrimSpace(strings.ToLower(categoryFilter))

	var foodCatMap map[int]map[string]bool
	if cleanCategory != "" && cleanCategory != "all" {
		foodCatMap, _ = r.GetFoodCategorySlugsMap(ctx)
	}

	var allMatching []RestaurantWithDistance

	for _, rest := range restaurants {
		// Calculate distance
		dist := CalculateHaversineDistance(userLat, userLng, rest.Latitude, rest.Longitude)

		// Filter by radius if provided
		if radiusKm > 0 && dist > radiusKm {
			continue
		}

		// Filter food items: available AND price <= maxPrice AND matching diet filter AND matching category filter
		var under100 []*ent.FoodItem
		var allAvailable []*ent.FoodItem
		minPrice := math.MaxFloat64
		hasMatchingFoodSearch := false
		restoIsPureVeg := true
		restoHasNonVeg := false

		for _, food := range rest.Edges.FoodItems {
			if !food.Available {
				continue
			}
			allAvailable = append(allAvailable, food)

			isVeg := IsDishVeg(food.Name, food.Category, food.Description)
			if !isVeg {
				restoHasNonVeg = true
				restoIsPureVeg = false
			}

			// Apply diet filter to under100 items
			if cleanDiet == "veg" && !isVeg {
				continue
			}
			if (cleanDiet == "nonveg" || cleanDiet == "non-veg") && isVeg {
				continue
			}

			// Apply category filter to food items
			if cleanCategory != "" && cleanCategory != "all" {
				isCatMatch := false
				if foodCatMap != nil && foodCatMap[food.ID] != nil && foodCatMap[food.ID][cleanCategory] {
					isCatMatch = true
				} else if MatchCategoryKeywords(cleanCategory, food.Category, food.Name, food.Description) {
					isCatMatch = true
				}
				if !isCatMatch {
					continue
				}
			}

			if food.Price < minPrice {
				minPrice = food.Price
			}

			if food.Price <= maxPrice {
				under100 = append(under100, food)
			}

			if cleanSearch != "" && (strings.Contains(strings.ToLower(food.Name), cleanSearch) || strings.Contains(strings.ToLower(food.Category), cleanSearch) || strings.Contains(strings.ToLower(food.Description), cleanSearch)) {
				hasMatchingFoodSearch = true
			}
		}

		// CRITICAL RULE:
		// "Show only restaurants having at least one available food item priced at ₹100 or less matching the criteria."
		if len(under100) == 0 {
			continue
		}

		// Search filter: check restaurant name, cuisine/description, address, or food items
		if cleanSearch != "" {
			nameMatch := strings.Contains(strings.ToLower(rest.Name), cleanSearch)
			descMatch := strings.Contains(strings.ToLower(rest.Description), cleanSearch)
			addrMatch := strings.Contains(strings.ToLower(rest.Address), cleanSearch)

			if !nameMatch && !descMatch && !addrMatch && !hasMatchingFoodSearch {
				continue
			}
		}

		// Sort under 100 food items by price ascending
		sort.Slice(under100, func(i, j int) bool {
			return under100[i].Price < under100[j].Price
		})

		isFav := false
		if userFavIDs != nil && userFavIDs[rest.ID] {
			isFav = true
		}

		allMatching = append(allMatching, RestaurantWithDistance{
			Restaurant:      rest,
			DistanceKm:      dist,
			DistanceDisplay: FormatDistance(dist),
			Under100Items:   under100,
			AllItems:        allAvailable,
			MinPrice:        minPrice,
			IsPureVeg:       restoIsPureVeg,
			HasNonVeg:       restoHasNonVeg,
			IsFavorite:      isFav,
		})
	}

	// 2. Batch load availability for all matching restaurants
	if len(allMatching) > 0 {
		candIDs := make([]int, len(allMatching))
		for i, itm := range allMatching {
			candIDs[i] = itm.Restaurant.ID
		}
		availMap, _ := r.GetBatchRestaurantAvailability(ctx, candIDs, time.Now())

		if openNowFilter {
			var openOnly []RestaurantWithDistance
			for _, item := range allMatching {
				av, ok := availMap[item.Restaurant.ID]
				if ok && av.IsOpen {
					item.Availability = av
					openOnly = append(openOnly, item)
				}
			}
			allMatching = openOnly
		} else {
			for i := range allMatching {
				if av, ok := availMap[allMatching[i].Restaurant.ID]; ok {
					allMatching[i].Availability = av
				}
			}
		}
	}

	// Dynamic Sorting: distance, price, name, or alphabetical
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "price":
		sort.Slice(allMatching, func(i, j int) bool {
			return allMatching[i].MinPrice < allMatching[j].MinPrice
		})
	case "name", "alphabetical":
		sort.Slice(allMatching, func(i, j int) bool {
			return strings.ToLower(allMatching[i].Restaurant.Name) < strings.ToLower(allMatching[j].Restaurant.Name)
		})
	default: // "distance"
		sort.Slice(allMatching, func(i, j int) bool {
			return allMatching[i].DistanceKm < allMatching[j].DistanceKm
		})
	}

	totalCount := len(allMatching)
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Clamp page
	if page > totalPages && totalCount > 0 {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	var pagedItems []RestaurantWithDistance
	if start < totalCount {
		if end > totalCount {
			end = totalCount
		}
		pagedItems = allMatching[start:end]
	}

	if len(pagedItems) > 0 {
		rIDs := make([]int, len(pagedItems))
		for i, itm := range pagedItems {
			rIDs[i] = itm.Restaurant.ID
		}
		ratingsMap, _ := r.GetRestaurantRatingsBatch(ctx, rIDs)
		delivMap, _ := r.GetDeliverySettingsBatch(ctx, rIDs)
		for i := range pagedItems {
			if st, ok := ratingsMap[pagedItems[i].Restaurant.ID]; ok {
				pagedItems[i].AverageRating = st.AverageRating
				pagedItems[i].ReviewCount = st.ReviewCount
			}
			ds := DefaultDeliverySettings(pagedItems[i].Restaurant.ID)
			if set, ok := delivMap[pagedItems[i].Restaurant.ID]; ok {
				ds = set
			}
			pagedItems[i].Serviceability = CalculateServiceability(ds, pagedItems[i].DistanceKm)
		}
	}

	return &PaginatedRestaurants{
		Items:          pagedItems,
		TotalCount:     totalCount,
		Page:           page,
		PageSize:       pageSize,
		TotalPages:     totalPages,
		HasNext:        page < totalPages,
		HasPrev:        page > 1,
		SearchQuery:    searchQuery,
		DietFilter:     dietFilter,
		CategoryFilter: categoryFilter,
		OpenNowFilter:  openNowFilter,
		MaxPrice:       maxPrice,
		SortBy:         sortBy,
	}, nil
}

// EnsureFavoritesTable idempotently creates the favorites table and indexes
func (r *Repo) EnsureFavoritesTable(ctx context.Context) {
	if r.db == nil {
		return
	}
	ddl := `
	CREATE TABLE IF NOT EXISTS favorites (
		id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CONSTRAINT favorites_user_restaurant_unique UNIQUE (user_id, restaurant_id)
	);
	CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
	CREATE INDEX IF NOT EXISTS idx_favorites_restaurant_id ON favorites(restaurant_id);
	`
	_, _ = r.db.ExecContext(ctx, ddl)
}

// IsFavorite returns true if the specified restaurant is favorited by the user
func (r *Repo) IsFavorite(ctx context.Context, userID, restaurantID int) (bool, error) {
	if r.db == nil || userID <= 0 || restaurantID <= 0 {
		return false, nil
	}
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND restaurant_id = $2)", userID, restaurantID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ToggleFavorite toggles a restaurant favorite for the user, returning the new status
func (r *Repo) ToggleFavorite(ctx context.Context, userID, restaurantID int) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("database connection not available")
	}
	if userID <= 0 || restaurantID <= 0 {
		return false, fmt.Errorf("invalid user ID or restaurant ID")
	}

	isFav, err := r.IsFavorite(ctx, userID, restaurantID)
	if err != nil {
		return false, err
	}

	if isFav {
		_, err := r.db.ExecContext(ctx, "DELETE FROM favorites WHERE user_id = $1 AND restaurant_id = $2", userID, restaurantID)
		if err != nil {
			return false, fmt.Errorf("failed to remove favorite: %w", err)
		}
		return false, nil
	}

	// Insert favorite with duplicate prevention
	_, err = r.db.ExecContext(ctx, "INSERT INTO favorites (user_id, restaurant_id, created_at) VALUES ($1, $2, NOW()) ON CONFLICT (user_id, restaurant_id) DO NOTHING", userID, restaurantID)
	if err != nil {
		return false, fmt.Errorf("failed to save favorite: %w", err)
	}
	return true, nil
}

// GetUserFavoriteRestaurantIDs returns a set of restaurant IDs favorited by the user
func (r *Repo) GetUserFavoriteRestaurantIDs(ctx context.Context, userID int) (map[int]bool, error) {
	favMap := make(map[int]bool)
	if r.db == nil || userID <= 0 {
		return favMap, nil
	}

	rows, err := r.db.QueryContext(ctx, "SELECT restaurant_id FROM favorites WHERE user_id = $1", userID)
	if err != nil {
		return favMap, fmt.Errorf("failed to query user favorites: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rID int
		if err := rows.Scan(&rID); err == nil {
			favMap[rID] = true
		}
	}
	return favMap, nil
}

// GetUserFavoriteRestaurants returns the full list of favorited restaurants for the user
func (r *Repo) GetUserFavoriteRestaurants(ctx context.Context, userID int, userLat, userLng float64) ([]RestaurantWithDistance, error) {
	if r.db == nil || userID <= 0 {
		return nil, nil
	}

	if userLat == 0 && userLng == 0 {
		userLat = 12.9716
		userLng = 77.5946
	}

	rows, err := r.db.QueryContext(ctx, "SELECT restaurant_id FROM favorites WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorite restaurant IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return []RestaurantWithDistance{}, nil
	}

	// Query only APPROVED restaurants (strictly enforcing visibility rules)
	restaurants, err := r.client.Restaurant.Query().
		Where(
			restaurant.IDIn(ids...),
			restaurant.StatusEQ(restaurant.StatusApproved),
		).
		WithFoodItems().
		WithCity().
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to load favorite restaurants: %w", err)
	}

	// Preserve order of favorites
	restMap := make(map[int]*ent.Restaurant)
	for _, rest := range restaurants {
		restMap[rest.ID] = rest
	}

	var results []RestaurantWithDistance
	for _, id := range ids {
		rest, ok := restMap[id]
		if !ok {
			continue // Deleted, pending, or rejected restaurant filtered out safely
		}

		dist := CalculateHaversineDistance(userLat, userLng, rest.Latitude, rest.Longitude)

		var under100 []*ent.FoodItem
		var allAvailable []*ent.FoodItem
		minPrice := math.MaxFloat64
		restoIsPureVeg := true
		restoHasNonVeg := false

		for _, food := range rest.Edges.FoodItems {
			if !food.Available {
				continue
			}
			allAvailable = append(allAvailable, food)

			isVeg := IsDishVeg(food.Name, food.Category, food.Description)
			if !isVeg {
				restoHasNonVeg = true
				restoIsPureVeg = false
			}

			if food.Price < minPrice {
				minPrice = food.Price
			}
			if food.Price <= 100.0 {
				under100 = append(under100, food)
			}
		}

		sort.Slice(under100, func(i, j int) bool {
			return under100[i].Price < under100[j].Price
		})

		results = append(results, RestaurantWithDistance{
			Restaurant:      rest,
			DistanceKm:      dist,
			DistanceDisplay: FormatDistance(dist),
			Under100Items:   under100,
			AllItems:        allAvailable,
			MinPrice:        minPrice,
			IsPureVeg:       restoIsPureVeg,
			HasNonVeg:       restoHasNonVeg,
			IsFavorite:      true,
		})
	}

	if len(results) > 0 {
		rIDs := make([]int, len(results))
		for i, itm := range results {
			rIDs[i] = itm.Restaurant.ID
		}
		ratingsMap, _ := r.GetRestaurantRatingsBatch(ctx, rIDs)
		availMap, _ := r.GetBatchRestaurantAvailability(ctx, rIDs, time.Now())
		delivMap, _ := r.GetDeliverySettingsBatch(ctx, rIDs)
		for i := range results {
			if st, ok := ratingsMap[results[i].Restaurant.ID]; ok {
				results[i].AverageRating = st.AverageRating
				results[i].ReviewCount = st.ReviewCount
			}
			if av, ok := availMap[results[i].Restaurant.ID]; ok {
				results[i].Availability = av
			}
			ds := DefaultDeliverySettings(results[i].Restaurant.ID)
			if set, ok := delivMap[results[i].Restaurant.ID]; ok {
				ds = set
			}
			results[i].Serviceability = CalculateServiceability(ds, results[i].DistanceKm)
		}
	}

	return results, nil
}

// GetRestaurantByID fetches a restaurant and all its food items
func (r *Repo) GetRestaurantByID(ctx context.Context, id int) (*ent.Restaurant, []*ent.FoodItem, error) {
	rest, err := r.client.Restaurant.Query().
		Where(restaurant.IDEQ(id)).
		WithCity().
		WithFoodItems().
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Sort food items by price ascending
	foods := rest.Edges.FoodItems
	sort.Slice(foods, func(i, j int) bool {
		return foods[i].Price < foods[j].Price
	})

	return rest, foods, nil
}

// GetAllRestaurantsAdmin returns restaurants for the admin view with optional status and search filter, eager-loading owner
func (r *Repo) GetAllRestaurantsAdmin(ctx context.Context, statusFilter, search string) ([]*ent.Restaurant, error) {
	query := r.client.Restaurant.Query().WithCity().WithFoodItems().WithOwner()

	if statusFilter != "" && statusFilter != "all" {
		query = query.Where(restaurant.StatusEQ(restaurant.Status(statusFilter)))
	}

	rests, err := query.Order(ent.Desc(restaurant.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	cleanSearch := strings.TrimSpace(strings.ToLower(search))
	if cleanSearch == "" {
		return rests, nil
	}

	var filtered []*ent.Restaurant
	for _, item := range rests {
		if strings.Contains(strings.ToLower(item.Name), cleanSearch) ||
			strings.Contains(strings.ToLower(item.Address), cleanSearch) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

// GetAdminStats returns comprehensive dashboard metrics calculated from active database records
func (r *Repo) GetAdminStats(ctx context.Context) (AdminStats, error) {
	var stats AdminStats

	total, err := r.client.Restaurant.Query().Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.TotalRestaurants = total

	pending, _ := r.client.Restaurant.Query().Where(restaurant.StatusEQ(restaurant.StatusPending)).Count(ctx)
	stats.PendingRestaurants = pending

	approved, _ := r.client.Restaurant.Query().Where(restaurant.StatusEQ(restaurant.StatusApproved)).Count(ctx)
	stats.ApprovedRestaurants = approved

	rejected, _ := r.client.Restaurant.Query().Where(restaurant.StatusEQ(restaurant.StatusRejected)).Count(ctx)
	stats.RejectedRestaurants = rejected

	unassigned, _ := r.client.Restaurant.Query().Where(restaurant.Not(restaurant.HasOwner())).Count(ctx)
	stats.UnassignedRestaurants = unassigned

	totalFood, _ := r.client.FoodItem.Query().Count(ctx)
	stats.TotalFoodItems = totalFood

	under100, _ := r.client.FoodItem.Query().
		Where(
			fooditem.PriceLTE(100.0),
			fooditem.AvailableEQ(true),
		).Count(ctx)
	stats.Under100FoodItems = under100

	totalUsers, _ := r.client.User.Query().Count(ctx)
	stats.TotalUsers = totalUsers

	customers, _ := r.client.User.Query().Where(user.RoleEQ(user.RoleCustomer)).Count(ctx)
	stats.CustomerCount = customers

	merchants, _ := r.client.User.Query().Where(user.RoleEQ(user.RoleMerchant)).Count(ctx)
	stats.MerchantCount = merchants

	admins, _ := r.client.User.Query().Where(user.Or(user.RoleEQ(user.RoleAdmin), user.IsAdminEQ(true))).Count(ctx)
	stats.AdminCount = admins

	return stats, nil
}

// CreateRestaurant creates a new restaurant
func (r *Repo) CreateRestaurant(ctx context.Context, input RestaurantInput) (*ent.Restaurant, error) {
	// Default city to Bangalore if not provided
	cityName := strings.TrimSpace(input.CityName)
	if cityName == "" {
		cityName = "Bangalore"
	}

	c, err := r.client.City.Query().Where(city.NameEQ(cityName)).Only(ctx)
	if err != nil {
		// Fallback to first city
		c, err = r.client.City.Query().First(ctx)
		if err != nil {
			return nil, fmt.Errorf("city not found: %w", err)
		}
	}

	stat := restaurant.Status(input.Status)
	if stat != restaurant.StatusApproved && stat != restaurant.StatusRejected {
		stat = restaurant.StatusPending
	}

	return r.client.Restaurant.Create().
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetAddress(strings.TrimSpace(input.Address)).
		SetLatitude(input.Latitude).
		SetLongitude(input.Longitude).
		SetPhone(strings.TrimSpace(input.Phone)).
		SetWebsite(strings.TrimSpace(input.Website)).
		SetOpeningHours(strings.TrimSpace(input.OpeningHours)).
		SetStatus(stat).
		SetCity(c).
		Save(ctx)
}

// UpdateRestaurant updates restaurant details
func (r *Repo) UpdateRestaurant(ctx context.Context, id int, input RestaurantInput) (*ent.Restaurant, error) {
	builder := r.client.Restaurant.UpdateOneID(id).
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetAddress(strings.TrimSpace(input.Address)).
		SetLatitude(input.Latitude).
		SetLongitude(input.Longitude).
		SetPhone(strings.TrimSpace(input.Phone)).
		SetWebsite(strings.TrimSpace(input.Website)).
		SetOpeningHours(strings.TrimSpace(input.OpeningHours))

	if input.Status != "" {
		builder.SetStatus(restaurant.Status(input.Status))
	}

	return builder.Save(ctx)
}

// UpdateRestaurantStatus updates approval status
func (r *Repo) UpdateRestaurantStatus(ctx context.Context, id int, status restaurant.Status) (*ent.Restaurant, error) {
	return r.client.Restaurant.UpdateOneID(id).
		SetStatus(status).
		Save(ctx)
}

// DeleteRestaurant deletes a restaurant and cascaded food items
func (r *Repo) DeleteRestaurant(ctx context.Context, id int) error {
	return r.client.Restaurant.DeleteOneID(id).Exec(ctx)
}

// AddFoodItem adds a food item to a restaurant
func (r *Repo) AddFoodItem(ctx context.Context, restaurantID int, input FoodInput) (*ent.FoodItem, error) {
	rest, err := r.client.Restaurant.Get(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("restaurant not found: %w", err)
	}

	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = "General"
	}

	return r.client.FoodItem.Create().
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetCategory(category).
		SetPrice(input.Price).
		SetAvailable(input.Available).
		SetRestaurant(rest).
		Save(ctx)
}

// UpdateFoodItem updates food item fields
func (r *Repo) UpdateFoodItem(ctx context.Context, foodID int, input FoodInput) (*ent.FoodItem, error) {
	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = "General"
	}

	return r.client.FoodItem.UpdateOneID(foodID).
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetCategory(category).
		SetPrice(input.Price).
		SetAvailable(input.Available).
		Save(ctx)
}

// ToggleFoodAvailability toggles availability status
func (r *Repo) ToggleFoodAvailability(ctx context.Context, foodID int) (*ent.FoodItem, error) {
	item, err := r.client.FoodItem.Get(ctx, foodID)
	if err != nil {
		return nil, err
	}
	return r.client.FoodItem.UpdateOneID(foodID).
		SetAvailable(!item.Available).
		Save(ctx)
}

// DeleteFoodItem deletes a food item
func (r *Repo) DeleteFoodItem(ctx context.Context, foodID int) error {
	return r.client.FoodItem.DeleteOneID(foodID).Exec(ctx)
}

// CommunitySuggestInput represents a user nomination for an eatery under ₹100
type CommunitySuggestInput struct {
	Name      string  `json:"name" form:"name"`
	Locality  string  `json:"locality" form:"locality"`
	DishName  string  `json:"dish_name" form:"dish_name"`
	DishPrice float64 `json:"dish_price" form:"dish_price"`
	DietType  string  `json:"diet_type" form:"diet_type"`
	Notes     string  `json:"notes" form:"notes"`
	Submitter string  `json:"submitter" form:"submitter"`
}

// SuggestRestaurant saves a community suggestion as a pending restaurant and food item
func (r *Repo) SuggestRestaurant(ctx context.Context, input CommunitySuggestInput) (*ent.Restaurant, error) {
	cityName := "Bangalore"
	c, err := r.client.City.Query().Where(city.NameEQ(cityName)).First(ctx)
	if err != nil {
		c, _ = r.client.City.Query().First(ctx)
	}

	lat := 12.9716
	lng := 77.5946

	desc := fmt.Sprintf("Community nomination. Must-try dish: %s (₹%.0f). %s", input.DishName, input.DishPrice, input.Notes)
	if strings.TrimSpace(input.Submitter) != "" {
		desc = fmt.Sprintf("Suggested by %s. Must-try dish: %s (₹%.0f). %s", input.Submitter, input.DishName, input.DishPrice, input.Notes)
	}

	restQuery := r.client.Restaurant.Create().
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(desc).
		SetAddress(strings.TrimSpace(input.Locality)).
		SetLatitude(lat).
		SetLongitude(lng).
		SetStatus(restaurant.StatusPending)

	if c != nil {
		restQuery = restQuery.SetCity(c)
	}

	newRest, err := restQuery.Save(ctx)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.DishName) != "" && input.DishPrice > 0 {
		cat := "Specialty"
		if strings.ToLower(input.DietType) == "veg" {
			cat = "Pure Veg Specialty"
		}
		_, _ = r.client.FoodItem.Create().
			SetName(strings.TrimSpace(input.DishName)).
			SetPrice(input.DishPrice).
			SetCategory(cat).
			SetDescription(fmt.Sprintf("%s - ₹%.0f", input.DietType, input.DishPrice)).
			SetAvailable(true).
			SetRestaurant(newRest).
			Save(ctx)
	}

	return newRest, nil
}

// CanManageRestaurant checks if a user is an admin or the assigned owner of the restaurant
func (r *Repo) CanManageRestaurant(ctx context.Context, u *ent.User, restaurantID int) (bool, error) {
	if u == nil {
		return false, nil
	}
	// Admins can manage any restaurant (checking both is_admin and role)
	if u.IsAdmin || u.Role == user.RoleAdmin {
		return true, nil
	}
	// Merchants can only manage restaurants they own
	if u.Role != user.RoleMerchant {
		return false, nil
	}
	exists, err := r.client.Restaurant.Query().
		Where(restaurant.IDEQ(restaurantID), restaurant.HasOwnerWith(user.IDEQ(u.ID))).
		Exist(ctx)
	return exists, err
}

// GetRestaurantsByOwner returns all restaurants owned by the specified user
func (r *Repo) GetRestaurantsByOwner(ctx context.Context, userID int) ([]*ent.Restaurant, error) {
	return r.client.Restaurant.Query().
		Where(restaurant.HasOwnerWith(user.IDEQ(userID))).
		WithCity().
		All(ctx)
}

// AssignRestaurantOwner assigns an owner to a restaurant
func (r *Repo) AssignRestaurantOwner(ctx context.Context, restaurantID int, ownerID int) (*ent.Restaurant, error) {
	return r.client.Restaurant.UpdateOneID(restaurantID).
		SetOwnerID(ownerID).
		Save(ctx)
}

// SubmitMerchantRestaurant creates a new pending restaurant owned by the authenticated merchant
func (r *Repo) SubmitMerchantRestaurant(ctx context.Context, merchantID int, input RestaurantInput) (*ent.Restaurant, error) {
	cityName := strings.TrimSpace(input.CityName)
	if cityName == "" {
		cityName = "Bangalore"
	}

	c, err := r.client.City.Query().Where(city.NameEQ(cityName)).Only(ctx)
	if err != nil {
		c, err = r.client.City.Query().First(ctx)
		if err != nil {
			return nil, fmt.Errorf("city not found: %w", err)
		}
	}

	lat := input.Latitude
	lng := input.Longitude
	if lat == 0 && lng == 0 {
		lat = 12.9716
		lng = 77.5946
	}

	return r.client.Restaurant.Create().
		SetName(strings.TrimSpace(input.Name)).
		SetDescription(strings.TrimSpace(input.Description)).
		SetAddress(strings.TrimSpace(input.Address)).
		SetLatitude(lat).
		SetLongitude(lng).
		SetPhone(strings.TrimSpace(input.Phone)).
		SetWebsite(strings.TrimSpace(input.Website)).
		SetOpeningHours(strings.TrimSpace(input.OpeningHours)).
		SetStatus(restaurant.StatusPending).
		SetOwnerID(merchantID).
		SetCity(c).
		Save(ctx)
}

// CanManageFoodItem checks if a user is an admin or the owner of the restaurant serving this food item
func (r *Repo) CanManageFoodItem(ctx context.Context, u *ent.User, foodID int) (bool, *ent.Restaurant, error) {
	if u == nil {
		return false, nil, nil
	}
	food, err := r.client.FoodItem.Query().
		Where(fooditem.IDEQ(foodID)).
		WithRestaurant().
		Only(ctx)
	if err != nil {
		return false, nil, err
	}
	if food.Edges.Restaurant == nil {
		return false, nil, nil
	}
	// Admins can manage any food item
	if u.IsAdmin || u.Role == user.RoleAdmin {
		return true, food.Edges.Restaurant, nil
	}
	// Merchants must own the restaurant
	if u.Role != user.RoleMerchant {
		return false, nil, nil
	}
	canManage, err := r.CanManageRestaurant(ctx, u, food.Edges.Restaurant.ID)
	if err != nil || !canManage {
		return false, nil, err
	}
	return true, food.Edges.Restaurant, nil
}

// UnassignRestaurantOwner clears the owner edge from a restaurant
func (r *Repo) UnassignRestaurantOwner(ctx context.Context, restaurantID int) (*ent.Restaurant, error) {
	return r.client.Restaurant.UpdateOneID(restaurantID).
		ClearOwner().
		Save(ctx)
}

// AssignRestaurantOwnerValidated validates that the target user exists and is a merchant before assigning
func (r *Repo) AssignRestaurantOwnerValidated(ctx context.Context, restaurantID int, ownerID int) (*ent.Restaurant, error) {
	if ownerID <= 0 {
		return r.UnassignRestaurantOwner(ctx, restaurantID)
	}

	u, err := r.client.User.Query().
		Where(user.IDEQ(ownerID)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("user with ID %d not found: %w", ownerID, err)
	}
	if u.Role != user.RoleMerchant {
		return nil, fmt.Errorf("user %s is not an eligible merchant (role is %s)", u.Email, u.Role)
	}

	return r.client.Restaurant.UpdateOneID(restaurantID).
		SetOwner(u).
		Save(ctx)
}

// GetEligibleMerchants returns all users with the merchant role
func (r *Repo) GetEligibleMerchants(ctx context.Context) ([]*ent.User, error) {
	return r.client.User.Query().
		Where(user.RoleEQ(user.RoleMerchant)).
		Order(ent.Asc(user.FieldName)).
		All(ctx)
}

// GetAllUsersAdmin queries users with role and search filtering, including their assigned restaurants
func (r *Repo) GetAllUsersAdmin(ctx context.Context, roleFilter, search string) ([]*ent.User, error) {
	query := r.client.User.Query().WithRestaurants()

	if roleFilter != "" && roleFilter != "all" {
		query = query.Where(user.RoleEQ(user.Role(roleFilter)))
	}

	users, err := query.Order(ent.Desc(user.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}

	cleanSearch := strings.TrimSpace(strings.ToLower(search))
	if cleanSearch == "" {
		return users, nil
	}

	var filtered []*ent.User
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), cleanSearch) ||
			strings.Contains(strings.ToLower(u.Email), cleanSearch) {
			filtered = append(filtered, u)
		}
	}
	return filtered, nil
}

// GetUserStats returns customer activity stats (reviews count, favorites count)
func (r *Repo) GetUserStats(ctx context.Context, userID int) (int, int, error) {
	revCount := 0
	favCount := 0
	if r.db != nil && userID > 0 {
		_ = r.db.QueryRowContext(ctx, "SELECT count(*) FROM reviews WHERE user_id = $1", userID).Scan(&revCount)
		_ = r.db.QueryRowContext(ctx, "SELECT count(*) FROM favorites WHERE user_id = $1", userID).Scan(&favCount)
	}
	return revCount, favCount, nil
}

// MerchantRestInspect holds restaurant metrics for the admin inspector
type MerchantRestInspect struct {
	ID            int
	Name          string
	Address       string
	Phone         string
	Status        string
	TotalFoods    int
	Under100Foods int
	OpeningHours  string
}

// EnsureReviewsTable idempotently creates the reviews and review_reports tables
func (r *Repo) EnsureReviewsTable(ctx context.Context) {
	if r.db == nil {
		return
	}
	ddl := `
	CREATE TABLE IF NOT EXISTS reviews (
		id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
		comment TEXT NOT NULL,
		status VARCHAR(32) NOT NULL DEFAULT 'approved',
		moderation_reason TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CONSTRAINT reviews_user_restaurant_unique UNIQUE (user_id, restaurant_id)
	);
	CREATE INDEX IF NOT EXISTS idx_reviews_restaurant_id ON reviews(restaurant_id);
	CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
	CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);

	CREATE TABLE IF NOT EXISTS review_reports (
		id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		review_id BIGINT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		reason TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		CONSTRAINT review_reports_review_user_unique UNIQUE (review_id, user_id)
	);
	CREATE INDEX IF NOT EXISTS idx_review_reports_review_id ON review_reports(review_id);
	CREATE INDEX IF NOT EXISTS idx_review_reports_user_id ON review_reports(user_id);
	`
	_, _ = r.db.ExecContext(ctx, ddl)
}

// GetRestaurantRatingStats returns the calculated average rating and count of approved reviews
func (r *Repo) GetRestaurantRatingStats(ctx context.Context, restaurantID int) (RestaurantRatingStats, error) {
	stats := RestaurantRatingStats{}
	if r.db == nil || restaurantID <= 0 {
		return stats, nil
	}

	query := `
		SELECT 
			COALESCE(ROUND(AVG(rating)::numeric, 1), 0.0),
			COUNT(*)
		FROM reviews 
		WHERE restaurant_id = $1 AND status = 'approved'
	`
	err := r.db.QueryRowContext(ctx, query, restaurantID).Scan(&stats.AverageRating, &stats.ReviewCount)
	if err != nil {
		return stats, fmt.Errorf("failed to get restaurant rating stats: %w", err)
	}
	return stats, nil
}

// GetRestaurantRatingsBatch fetches ratings in batch for a list of restaurant IDs
func (r *Repo) GetRestaurantRatingsBatch(ctx context.Context, restaurantIDs []int) (map[int]RestaurantRatingStats, error) {
	result := make(map[int]RestaurantRatingStats)
	if r.db == nil || len(restaurantIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(restaurantIDs))
	args := make([]any, len(restaurantIDs))
	for i, id := range restaurantIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			restaurant_id,
			COALESCE(ROUND(AVG(rating)::numeric, 1), 0.0),
			COUNT(*)
		FROM reviews 
		WHERE restaurant_id IN (%s) AND status = 'approved'
		GROUP BY restaurant_id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, fmt.Errorf("failed to batch get ratings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rID int
		var avg float64
		var count int
		if err := rows.Scan(&rID, &avg, &count); err == nil {
			result[rID] = RestaurantRatingStats{
				AverageRating: avg,
				ReviewCount:   count,
			}
		}
	}
	return result, nil
}

// CreateReview creates a new customer rating and review with validation and duplicate prevention
func (r *Repo) CreateReview(ctx context.Context, userID, restaurantID, rating int, comment string) (*Review, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	if userID <= 0 {
		return nil, fmt.Errorf("unauthenticated customer")
	}
	if restaurantID <= 0 {
		return nil, fmt.Errorf("invalid restaurant ID")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5 stars")
	}
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return nil, fmt.Errorf("review comment cannot be empty")
	}
	if len(comment) > 1000 {
		return nil, fmt.Errorf("review comment exceeds maximum length of 1000 characters")
	}

	// Strictly verify that the restaurant exists and is APPROVED
	rest, err := r.client.Restaurant.Query().Where(restaurant.IDEQ(restaurantID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("restaurant not found")
	}
	if rest.Status != restaurant.StatusApproved {
		return nil, fmt.Errorf("reviews are only permitted for approved restaurants")
	}

	var newID int
	var createdAt, updatedAt time.Time
	query := `
		INSERT INTO reviews (user_id, restaurant_id, rating, comment, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'approved', NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = r.db.QueryRowContext(ctx, query, userID, restaurantID, rating, comment).Scan(&newID, &createdAt, &updatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "reviews_user_restaurant_unique") {
			return nil, fmt.Errorf("you have already reviewed this restaurant")
		}
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	userName := "Customer"
	u, err := r.client.User.Get(ctx, userID)
	if err == nil && u != nil {
		userName = u.Name
	}

	return &Review{
		ID:           newID,
		UserID:       userID,
		RestaurantID: restaurantID,
		Rating:       rating,
		Comment:      comment,
		Status:       "approved",
		UserName:     userName,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

// GetReviewByID fetches a review by its ID
func (r *Repo) GetReviewByID(ctx context.Context, reviewID int) (*Review, error) {
	if r.db == nil || reviewID <= 0 {
		return nil, fmt.Errorf("invalid review ID")
	}
	query := `
		SELECT 
			r.id, r.user_id, r.restaurant_id, r.rating, r.comment, r.status, r.moderation_reason,
			r.created_at, r.updated_at, u.name, u.email
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.id = $1
	`
	var rev Review
	err := r.db.QueryRowContext(ctx, query, reviewID).Scan(
		&rev.ID, &rev.UserID, &rev.RestaurantID, &rev.Rating, &rev.Comment, &rev.Status, &rev.ModerationReason,
		&rev.CreatedAt, &rev.UpdatedAt, &rev.UserName, &rev.UserEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, err
	}
	return &rev, nil
}

// GetUserReviewForRestaurant returns the existing review for a specific user and restaurant, if any
func (r *Repo) GetUserReviewForRestaurant(ctx context.Context, userID, restaurantID int) (*Review, error) {
	if r.db == nil || userID <= 0 || restaurantID <= 0 {
		return nil, nil
	}
	query := `
		SELECT 
			r.id, r.user_id, r.restaurant_id, r.rating, r.comment, r.status, r.moderation_reason,
			r.created_at, r.updated_at, u.name, u.email
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1 AND r.restaurant_id = $2
	`
	var rev Review
	err := r.db.QueryRowContext(ctx, query, userID, restaurantID).Scan(
		&rev.ID, &rev.UserID, &rev.RestaurantID, &rev.Rating, &rev.Comment, &rev.Status, &rev.ModerationReason,
		&rev.CreatedAt, &rev.UpdatedAt, &rev.UserName, &rev.UserEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rev, nil
}

// UpdateReview updates an existing review belonging to the user
func (r *Repo) UpdateReview(ctx context.Context, reviewID, userID, rating int, comment string) (*Review, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection not available")
	}
	if reviewID <= 0 || userID <= 0 {
		return nil, fmt.Errorf("invalid review or user ID")
	}
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5 stars")
	}
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return nil, fmt.Errorf("review comment cannot be empty")
	}
	if len(comment) > 1000 {
		return nil, fmt.Errorf("review comment exceeds maximum length of 1000 characters")
	}

	var restaurantID int
	var status string
	err := r.db.QueryRowContext(ctx, "SELECT restaurant_id, status FROM reviews WHERE id = $1 AND user_id = $2", reviewID, userID).Scan(&restaurantID, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found or unauthorized")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	rest, err := r.client.Restaurant.Query().Where(restaurant.IDEQ(restaurantID)).Only(ctx)
	if err != nil || rest.Status != restaurant.StatusApproved {
		return nil, fmt.Errorf("restaurant is not eligible for reviews")
	}

	query := `
		UPDATE reviews
		SET rating = $1, comment = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
	`
	_, err = r.db.ExecContext(ctx, query, rating, comment, reviewID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return r.GetReviewByID(ctx, reviewID)
}

// DeleteReview deletes a review belonging to the user
func (r *Repo) DeleteReview(ctx context.Context, reviewID, userID int) error {
	if r.db == nil {
		return fmt.Errorf("database connection not available")
	}
	if reviewID <= 0 || userID <= 0 {
		return fmt.Errorf("invalid review or user ID")
	}

	res, err := r.db.ExecContext(ctx, "DELETE FROM reviews WHERE id = $1 AND user_id = $2", reviewID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("review not found or unauthorized")
	}
	return nil
}

// ReportReview submits an abuse report against a review
func (r *Repo) ReportReview(ctx context.Context, reviewID, userID int, reason string) error {
	if r.db == nil {
		return fmt.Errorf("database connection not available")
	}
	if reviewID <= 0 || userID <= 0 {
		return fmt.Errorf("invalid parameters")
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "Inappropriate content"
	}
	if len(reason) > 500 {
		reason = reason[:500]
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM reviews WHERE id = $1)", reviewID).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("review not found")
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO review_reports (review_id, user_id, reason, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (review_id, user_id) DO NOTHING
	`, reviewID, userID, reason)
	if err != nil {
		return fmt.Errorf("failed to report review: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("already reported")
	}
	return nil
}

// ModerateReview allows administrators to hide or restore a review with an auditable reason
func (r *Repo) ModerateReview(ctx context.Context, reviewID int, status string, reason string) error {
	if r.db == nil {
		return fmt.Errorf("database connection not available")
	}
	if reviewID <= 0 {
		return fmt.Errorf("invalid review ID")
	}
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "approved" && status != "hidden" {
		return fmt.Errorf("invalid moderation status (must be approved or hidden)")
	}
	reason = strings.TrimSpace(reason)
	if len(reason) > 500 {
		reason = reason[:500]
	}

	res, err := r.db.ExecContext(ctx, `
		UPDATE reviews 
		SET status = $1, moderation_reason = $2, updated_at = NOW() 
		WHERE id = $3
	`, status, reason, reviewID)
	if err != nil {
		return fmt.Errorf("failed to moderate review: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

// GetRestaurantReviews returns paginated approved reviews for a restaurant
func (r *Repo) GetRestaurantReviews(ctx context.Context, restaurantID int, page, pageSize int) (PaginatedReviews, error) {
	result := PaginatedReviews{
		Items:      []Review{},
		Page:       page,
		PageSize:   pageSize,
		TotalPages: 1,
	}
	if r.db == nil || restaurantID <= 0 {
		return result, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 5
	}
	result.Page = page
	result.PageSize = pageSize

	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews WHERE restaurant_id = $1 AND status = 'approved'", restaurantID).Scan(&total)
	if err != nil {
		return result, fmt.Errorf("failed to count reviews: %w", err)
	}
	result.TotalCount = total

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}
	result.TotalPages = totalPages
	result.HasNext = page < totalPages
	result.HasPrev = page > 1

	if total == 0 {
		return result, nil
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT 
			r.id, r.user_id, r.restaurant_id, r.rating, r.comment, r.status, r.moderation_reason,
			r.created_at, r.updated_at, u.name, u.email
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.restaurant_id = $1 AND r.status = 'approved'
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, restaurantID, pageSize, offset)
	if err != nil {
		return result, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rev Review
		if err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.RestaurantID, &rev.Rating, &rev.Comment, &rev.Status, &rev.ModerationReason,
			&rev.CreatedAt, &rev.UpdatedAt, &rev.UserName, &rev.UserEmail,
		); err == nil {
			result.Items = append(result.Items, rev)
		}
	}
	return result, nil
}

// GetAdminReviews returns reviews for admin moderation filtered by status
func (r *Repo) GetAdminReviews(ctx context.Context, filter string, page, pageSize int) (PaginatedReviews, error) {
	result := PaginatedReviews{
		Items:      []Review{},
		Page:       page,
		PageSize:   pageSize,
		TotalPages: 1,
	}
	if r.db == nil {
		return result, nil
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 15
	}
	result.Page = page
	result.PageSize = pageSize

	whereClause := "1=1"
	switch strings.ToLower(filter) {
	case "reported":
		whereClause = "(SELECT COUNT(*) FROM review_reports rr WHERE rr.review_id = r.id) > 0"
	case "hidden":
		whereClause = "r.status = 'hidden'"
	case "approved":
		whereClause = "r.status = 'approved'"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM reviews r WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return result, fmt.Errorf("failed to count admin reviews: %w", err)
	}
	result.TotalCount = total

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}
	result.TotalPages = totalPages
	result.HasNext = page < totalPages
	result.HasPrev = page > 1

	if total == 0 {
		return result, nil
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.user_id, r.restaurant_id, r.rating, r.comment, r.status, r.moderation_reason,
			r.created_at, r.updated_at, u.name, u.email, rest.name,
			(SELECT COUNT(*) FROM review_reports rr WHERE rr.review_id = r.id) AS report_count
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN restaurants rest ON r.restaurant_id = rest.id
		WHERE %s
		ORDER BY 
			(SELECT COUNT(*) FROM review_reports rr WHERE rr.review_id = r.id) DESC,
			r.created_at DESC
		LIMIT $1 OFFSET $2
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return result, fmt.Errorf("failed to query admin reviews: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rev Review
		if err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.RestaurantID, &rev.Rating, &rev.Comment, &rev.Status, &rev.ModerationReason,
			&rev.CreatedAt, &rev.UpdatedAt, &rev.UserName, &rev.UserEmail, &rev.RestaurantName, &rev.ReportCount,
		); err == nil {
			result.Items = append(result.Items, rev)
		}
	}
	return result, nil
}

// GetAdminReviewStats returns counts for admin review tabs
func (r *Repo) GetAdminReviewStats(ctx context.Context) (map[string]int, error) {
	counts := map[string]int{"all": 0, "reported": 0, "hidden": 0, "approved": 0}
	if r.db == nil {
		return counts, nil
	}
	var all, approved, hidden, reported int
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews").Scan(&all)
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews WHERE status = 'approved'").Scan(&approved)
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews WHERE status = 'hidden'").Scan(&hidden)
	_ = r.db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT review_id) FROM review_reports").Scan(&reported)
	counts["all"] = all
	counts["approved"] = approved
	counts["hidden"] = hidden
	counts["reported"] = reported
	return counts, nil
}




