package fooddiscovery

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Category represents a food category or cuisine classification
type Category struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	CategoryType string    `json:"category_type"` // "cuisine", "meal_type", "food_type"
	Description  string    `json:"description"`
	Icon         string    `json:"icon"`
	ImageURL     string    `json:"image_url"`
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	DishCount    int       `json:"dish_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Slugify converts a title into a clean URL-friendly slug
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// EnsureCategoriesTable idempotently creates the categories and food_item_categories tables and seeds initial categories
func (r *Repo) EnsureCategoriesTable(ctx context.Context) {
	if r.db == nil {
		return
	}

	ddl := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(100) NOT NULL UNIQUE,
		category_type VARCHAR(50) NOT NULL DEFAULT 'cuisine',
		description VARCHAR(255) DEFAULT '',
		icon VARCHAR(50) DEFAULT '',
		image_url VARCHAR(255) DEFAULT '',
		is_active BOOLEAN NOT NULL DEFAULT true,
		display_order INT NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);
	CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);

	CREATE TABLE IF NOT EXISTS food_item_categories (
		food_item_id BIGINT NOT NULL REFERENCES food_items(id) ON DELETE CASCADE,
		category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY (food_item_id, category_id)
	);

	CREATE INDEX IF NOT EXISTS idx_food_item_categories_food_id ON food_item_categories(food_item_id);
	CREATE INDEX IF NOT EXISTS idx_food_item_categories_cat_id ON food_item_categories(category_id);
	`
	_, _ = r.db.ExecContext(ctx, ddl)

	// Idempotently seed standard categories
	seedSQL := `
	INSERT INTO categories (name, slug, category_type, description, icon, image_url, display_order)
	VALUES
	  ('South Indian', 'south-indian', 'cuisine', 'Authentic Bangalore & Karnataka breakfast and tiffin delicacies', '🍛', '/static/hero_food.jpg', 1),
	  ('North Indian', 'north-indian', 'cuisine', 'Rich curries, fresh rotis, parathas, and Punjabi flavors', '🥘', '/static/aloo_paratha.jpg', 2),
	  ('Chinese', 'chinese', 'cuisine', 'Desi Chinese noodles, fried rice, and savory starters', '🥢', '/static/steamed_momos.jpg', 3),
	  ('Street Food', 'street-food', 'meal_type', 'Bangalore chaat, pani puri, bhel, and spicy street eats', '🥟', '/static/pani_puri.jpg', 4),
	  ('Snacks', 'snacks', 'meal_type', 'Quick crispy bites, samosas, pakodas, and evening snacks', '🥪', '/static/samosa_plate.jpg', 5),
	  ('Breakfast', 'breakfast', 'meal_type', 'Steaming idlis, crispy vadas, and hot morning filter coffee', '☕', '/static/idli_sambar.jpg', 6),
	  ('Lunch', 'lunch', 'meal_type', 'Wholesome budget meals, pulao, thalis, and rice bowls', '🍱', '/static/veg_pulao.jpg', 7),
	  ('Desserts', 'desserts', 'food_type', 'Sweet treats, halwa, kesari bath, and Indian mithai', '🍧', '/static/filter_coffee.jpg', 8),
	  ('Beverages', 'beverages', 'food_type', 'Filter coffee, badam milk, fresh juices, and hot chai', '🥤', '/static/filter_coffee.jpg', 9),
	  ('Fast Food', 'fast-food', 'food_type', 'Burgers, crispy rolls, sandwiches, and quick bites', '🍔', '/static/kathi_roll.jpg', 10),
	  ('Bakery', 'bakery', 'food_type', 'Fresh puffs, buns, tea cakes, and local bakery treats', '🥐', '/static/samosa_plate.jpg', 11),
	  ('Healthy Food', 'healthy-food', 'food_type', 'Wholesome salads, sprouts, fresh fruit bowls, and low-oil eats', '🥗', '/static/indori_poha.jpg', 12)
	ON CONFLICT (slug) DO NOTHING;
	`
	_, _ = r.db.ExecContext(ctx, seedSQL)

	// Idempotently seed mapping from existing food items
	mappingSQL := `
	INSERT INTO food_item_categories (food_item_id, category_id)
	SELECT f.id, c.id
	FROM food_items f
	JOIN categories c ON (
		(c.slug = 'south-indian' AND (LOWER(f.category) LIKE '%south indian%' OR LOWER(f.name) LIKE '%dosa%' OR LOWER(f.name) LIKE '%idli%' OR LOWER(f.name) LIKE '%vada%' OR LOWER(f.name) LIKE '%uttapam%')) OR
		(c.slug = 'north-indian' AND (LOWER(f.category) LIKE '%north indian%' OR LOWER(f.name) LIKE '%paratha%' OR LOWER(f.name) LIKE '%paneer%' OR LOWER(f.category) LIKE '%curries%')) OR
		(c.slug = 'chinese' AND (LOWER(f.category) LIKE '%chinese%' OR LOWER(f.name) LIKE '%momo%' OR LOWER(f.name) LIKE '%noodle%' OR LOWER(f.name) LIKE '%fried rice%')) OR
		(c.slug = 'street-food' AND (LOWER(f.category) LIKE '%street%' OR LOWER(f.category) LIKE '%chaat%' OR LOWER(f.name) LIKE '%puri%' OR LOWER(f.name) LIKE '%bhel%')) OR
		(c.slug = 'snacks' AND (LOWER(f.category) LIKE '%snack%' OR LOWER(f.name) LIKE '%samosa%' OR LOWER(f.name) LIKE '%pakora%' OR LOWER(f.name) LIKE '%roll%')) OR
		(c.slug = 'breakfast' AND (LOWER(f.category) LIKE '%breakfast%' OR LOWER(f.name) LIKE '%dosa%' OR LOWER(f.name) LIKE '%idli%' OR LOWER(f.name) LIKE '%poha%')) OR
		(c.slug = 'lunch' AND (LOWER(f.category) LIKE '%main course%' OR LOWER(f.category) LIKE '%biryani%' OR LOWER(f.name) LIKE '%rice%' OR LOWER(f.name) LIKE '%pulao%' OR LOWER(f.name) LIKE '%thali%')) OR
		(c.slug = 'desserts' AND (LOWER(f.category) LIKE '%dessert%' OR LOWER(f.name) LIKE '%bath%' OR LOWER(f.name) LIKE '%halwa%' OR LOWER(f.name) LIKE '%sweet%' OR LOWER(f.category) LIKE '%ice cream%' OR LOWER(f.category) LIKE '%sundae%')) OR
		(c.slug = 'beverages' AND (LOWER(f.category) LIKE '%beverage%' OR LOWER(f.name) LIKE '%coffee%' OR LOWER(f.name) LIKE '%tea%' OR LOWER(f.name) LIKE '%juice%')) OR
		(c.slug = 'fast-food' AND (LOWER(f.category) LIKE '%fast food%' OR LOWER(f.category) LIKE '%rolls%' OR LOWER(f.category) LIKE '%burgers%' OR LOWER(f.category) LIKE '%sandwiches%'))
	)
	ON CONFLICT (food_item_id, category_id) DO NOTHING;
	`
	_, _ = r.db.ExecContext(ctx, mappingSQL)
}

// GetActiveCategories returns all active categories ordered by display_order
func (r *Repo) GetActiveCategories(ctx context.Context) ([]Category, error) {
	if r.db == nil {
		return []Category{}, nil
	}

	query := `
		SELECT 
			c.id, c.name, c.slug, c.category_type, c.description, c.icon, c.image_url, 
			c.is_active, c.display_order, c.created_at, c.updated_at,
			COUNT(fc.food_item_id) as dish_count
		FROM categories c
		LEFT JOIN food_item_categories fc ON c.id = fc.category_id
		WHERE c.is_active = true
		GROUP BY c.id
		ORDER BY c.display_order ASC, c.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.CategoryType, &cat.Description,
			&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.DisplayOrder,
			&cat.CreatedAt, &cat.UpdatedAt, &cat.DishCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetAllCategoriesAdmin returns all categories (both active and inactive) for admin management
func (r *Repo) GetAllCategoriesAdmin(ctx context.Context) ([]Category, error) {
	if r.db == nil {
		return []Category{}, nil
	}

	query := `
		SELECT 
			c.id, c.name, c.slug, c.category_type, c.description, c.icon, c.image_url, 
			c.is_active, c.display_order, c.created_at, c.updated_at,
			COUNT(fc.food_item_id) as dish_count
		FROM categories c
		LEFT JOIN food_item_categories fc ON c.id = fc.category_id
		GROUP BY c.id
		ORDER BY c.display_order ASC, c.id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query admin categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.CategoryType, &cat.Description,
			&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.DisplayOrder,
			&cat.CreatedAt, &cat.UpdatedAt, &cat.DishCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetCategoryBySlug returns a single category by its slug
func (r *Repo) GetCategoryBySlug(ctx context.Context, slug string) (*Category, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}

	slug = strings.ToLower(strings.TrimSpace(slug))
	query := `
		SELECT 
			c.id, c.name, c.slug, c.category_type, c.description, c.icon, c.image_url, 
			c.is_active, c.display_order, c.created_at, c.updated_at,
			COUNT(fc.food_item_id) as dish_count
		FROM categories c
		LEFT JOIN food_item_categories fc ON c.id = fc.category_id
		WHERE c.slug = $1
		GROUP BY c.id
	`
	var cat Category
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.CategoryType, &cat.Description,
		&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.DisplayOrder,
		&cat.CreatedAt, &cat.UpdatedAt, &cat.DishCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &cat, nil
}

// GetCategoryByID returns a category by its ID
func (r *Repo) GetCategoryByID(ctx context.Context, id int) (*Category, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}

	query := `
		SELECT 
			c.id, c.name, c.slug, c.category_type, c.description, c.icon, c.image_url, 
			c.is_active, c.display_order, c.created_at, c.updated_at,
			COUNT(fc.food_item_id) as dish_count
		FROM categories c
		LEFT JOIN food_item_categories fc ON c.id = fc.category_id
		WHERE c.id = $1
		GROUP BY c.id
	`
	var cat Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.CategoryType, &cat.Description,
		&cat.Icon, &cat.ImageURL, &cat.IsActive, &cat.DisplayOrder,
		&cat.CreatedAt, &cat.UpdatedAt, &cat.DishCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &cat, nil
}

// CreateCategory adds a new category with server-side validation and duplicate prevention
func (r *Repo) CreateCategory(ctx context.Context, name, catType, desc, icon string) (*Category, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection unavailable")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	slug := Slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("invalid category name")
	}

	catType = strings.ToLower(strings.TrimSpace(catType))
	if catType != "cuisine" && catType != "meal_type" && catType != "food_type" {
		catType = "cuisine"
	}

	if icon == "" {
		icon = "🍽️"
	}

	// Check for duplicate slug
	var existingID int
	err := r.db.QueryRowContext(ctx, "SELECT id FROM categories WHERE slug = $1", slug).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("a category with slug '%s' already exists", slug)
	}

	var newID int
	query := `
		INSERT INTO categories (name, slug, category_type, description, icon, is_active, display_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, 50, NOW(), NOW())
		RETURNING id
	`
	err = r.db.QueryRowContext(ctx, query, name, slug, catType, desc, icon).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return r.GetCategoryByID(ctx, newID)
}

// UpdateCategory edits an existing category
func (r *Repo) UpdateCategory(ctx context.Context, id int, name, catType, desc, icon string) error {
	if r.db == nil {
		return fmt.Errorf("database connection unavailable")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("category name is required")
	}

	catType = strings.ToLower(strings.TrimSpace(catType))
	if catType != "cuisine" && catType != "meal_type" && catType != "food_type" {
		catType = "cuisine"
	}

	query := `
		UPDATE categories
		SET name = $1, category_type = $2, description = $3, icon = $4, updated_at = NOW()
		WHERE id = $5
	`
	res, err := r.db.ExecContext(ctx, query, name, catType, desc, icon, id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// ToggleCategoryActive activates or deactivates a category without deleting food records
func (r *Repo) ToggleCategoryActive(ctx context.Context, id int) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("database connection unavailable")
	}

	query := `
		UPDATE categories
		SET is_active = NOT is_active, updated_at = NOW()
		WHERE id = $1
		RETURNING is_active
	`
	var newStatus bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&newStatus)
	if err != nil {
		return false, fmt.Errorf("failed to toggle category status: %w", err)
	}
	return newStatus, nil
}

// GetFoodItemCategorySlugs returns all category slugs associated with a food item
func (r *Repo) GetFoodItemCategorySlugs(ctx context.Context, foodID int) ([]string, error) {
	if r.db == nil || foodID <= 0 {
		return []string{}, nil
	}

	query := `
		SELECT c.slug
		FROM categories c
		JOIN food_item_categories fc ON c.id = fc.category_id
		WHERE fc.food_item_id = $1 AND c.is_active = true
	`
	rows, err := r.db.QueryContext(ctx, query, foodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err == nil {
			slugs = append(slugs, s)
		}
	}
	return slugs, nil
}

// AssignFoodItemCategories associates a food item with categories by slug
func (r *Repo) AssignFoodItemCategories(ctx context.Context, foodID int, categorySlugs []string) error {
	if r.db == nil || foodID <= 0 {
		return nil
	}

	// Remove old associations
	_, _ = r.db.ExecContext(ctx, "DELETE FROM food_item_categories WHERE food_item_id = $1", foodID)

	for _, slug := range categorySlugs {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if slug == "" {
			continue
		}
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO food_item_categories (food_item_id, category_id)
			SELECT $1, id FROM categories WHERE slug = $2
			ON CONFLICT DO NOTHING
		`, foodID, slug)
	}
	return nil
}

// GetFoodCategorySlugsMap builds a fast in-memory map of food_item_id -> set of active category slugs
func (r *Repo) GetFoodCategorySlugsMap(ctx context.Context) (map[int]map[string]bool, error) {
	result := make(map[int]map[string]bool)
	if r.db == nil {
		return result, nil
	}

	query := `
		SELECT fc.food_item_id, c.slug
		FROM food_item_categories fc
		JOIN categories c ON fc.category_id = c.id
		WHERE c.is_active = true
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var foodID int
		var slug string
		if err := rows.Scan(&foodID, &slug); err == nil {
			if result[foodID] == nil {
				result[foodID] = make(map[string]bool)
			}
			result[foodID][slug] = true
		}
	}
	return result, nil
}

// MatchCategoryKeywords checks if a food item matches a category slug by keywords
func MatchCategoryKeywords(slug, cat, name, desc string) bool {
	slug = strings.ToLower(strings.TrimSpace(slug))
	text := strings.ToLower(cat + " " + name + " " + desc)

	switch slug {
	case "south-indian":
		return strings.Contains(text, "south indian") || strings.Contains(text, "dosa") || strings.Contains(text, "idli") || strings.Contains(text, "vada") || strings.Contains(text, "uttapam") || strings.Contains(text, "upma")
	case "north-indian":
		return strings.Contains(text, "north indian") || strings.Contains(text, "paratha") || strings.Contains(text, "paneer") || strings.Contains(text, "roti") || strings.Contains(text, "curry") || strings.Contains(text, "curries") || strings.Contains(text, "dal")
	case "chinese":
		return strings.Contains(text, "chinese") || strings.Contains(text, "momo") || strings.Contains(text, "noodle") || strings.Contains(text, "fried rice") || strings.Contains(text, "manchurian")
	case "street-food":
		return strings.Contains(text, "street") || strings.Contains(text, "chaat") || strings.Contains(text, "puri") || strings.Contains(text, "bhel") || strings.Contains(text, "sev") || strings.Contains(text, "kachori")
	case "snacks":
		return strings.Contains(text, "snack") || strings.Contains(text, "samosa") || strings.Contains(text, "pakora") || strings.Contains(text, "pakoda") || strings.Contains(text, "bajji") || strings.Contains(text, "cutlet")
	case "breakfast":
		return strings.Contains(text, "breakfast") || strings.Contains(text, "idli") || strings.Contains(text, "dosa") || strings.Contains(text, "poha") || strings.Contains(text, "upma") || strings.Contains(text, "khara bath")
	case "lunch":
		return strings.Contains(text, "lunch") || strings.Contains(text, "meals") || strings.Contains(text, "thali") || strings.Contains(text, "rice") || strings.Contains(text, "biryani") || strings.Contains(text, "pulao") || strings.Contains(text, "main course")
	case "desserts":
		return strings.Contains(text, "dessert") || strings.Contains(text, "sweet") || strings.Contains(text, "bath") || strings.Contains(text, "halwa") || strings.Contains(text, "ice cream") || strings.Contains(text, "sundae") || strings.Contains(text, "gulab") || strings.Contains(text, "kesari")
	case "beverages":
		return strings.Contains(text, "beverage") || strings.Contains(text, "coffee") || strings.Contains(text, "tea") || strings.Contains(text, "chai") || strings.Contains(text, "juice") || strings.Contains(text, "shake") || strings.Contains(text, "badam") || strings.Contains(text, "lassi")
	case "fast-food":
		return strings.Contains(text, "fast food") || strings.Contains(text, "roll") || strings.Contains(text, "burger") || strings.Contains(text, "sandwich") || strings.Contains(text, "pizza") || strings.Contains(text, "wrap")
	case "bakery":
		return strings.Contains(text, "bakery") || strings.Contains(text, "puff") || strings.Contains(text, "bun") || strings.Contains(text, "cake") || strings.Contains(text, "pastry") || strings.Contains(text, "biscuit")
	case "healthy-food":
		return strings.Contains(text, "healthy") || strings.Contains(text, "salad") || strings.Contains(text, "sprout") || strings.Contains(text, "fruit") || strings.Contains(text, "oats")
	default:
		slugNorm := strings.ReplaceAll(slug, "-", " ")
		return strings.Contains(text, slugNorm)
	}
}
