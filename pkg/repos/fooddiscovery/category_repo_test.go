package fooddiscovery

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"South Indian", "south-indian"},
		{"North Indian", "north-indian"},
		{"Snacks & Chaat", "snacks-chaat"},
		{"Special Dosa / Idli Combo!", "special-dosa-idli-combo"},
		{"   Street Food 100   ", "street-food-100"},
		{"Chinese", "chinese"},
		{"Desserts", "desserts"},
		{"Beverages", "beverages"},
		{"Fast Food", "fast-food"},
		{"Bakery", "bakery"},
		{"Healthy Food", "healthy-food"},
		{"", ""},
		{"---", ""},
	}

	for _, tt := range tests {
		actual := Slugify(tt.input)
		if actual != tt.expected {
			t.Errorf("Slugify(%q) = %q; expected %q", tt.input, actual, tt.expected)
		}
	}
}

func TestMatchCategoryKeywords(t *testing.T) {
	tests := []struct {
		slug     string
		foodCat  string
		foodName string
		foodDesc string
		expected bool
	}{
		// South Indian matching
		{"south-indian", "South Indian", "Crispy Masala Dosa", "Served with chutney", true},
		{"south-indian", "Breakfast", "Steamed Idli Sambar", "Traditional breakfast", true},
		{"south-indian", "Chinese", "Fried Rice", "Chinese meal", false},

		// North Indian matching
		{"north-indian", "North Indian", "Aloo Paratha with Curd", "Stuffed potato flatbread", true},
		{"north-indian", "Curries", "Paneer Butter Masala", "Rich tomato curry", true},
		{"north-indian", "Breakfast", "Filter Coffee", "Hot coffee", false},

		// Chinese matching
		{"chinese", "Chinese", "Veg Hakka Noodles", "Wok tossed noodles", true},
		{"chinese", "Snacks", "Steamed Veg Momos", "Tibetan momos with chilli sauce", true},
		{"chinese", "Breakfast", "Idli Vada", "South Indian", false},

		// Street food matching
		{"street-food", "Snacks", "Special Pani Puri", "Crispy puris with spicy mint water", true},
		{"street-food", "Fast Food", "Sev Puri", "Chaat snack", true},

		// Beverages matching
		{"beverages", "Beverages", "Degree Filter Coffee", "Bangalore filter coffee", true},
		{"beverages", "Drinks", "Badam Milk", "Almond milk", true},
		{"beverages", "Main Course", "Biryani Bowl", "Rice meal", false},

		// Desserts matching
		{"desserts", "Dessert", "Gulab Jamun", "Sweet dessert balls", true},
		{"desserts", "Sweets", "Kesari Bath", "Semolina halwa", true},

		// Healthy food matching
		{"healthy-food", "Salads", "Sprout Salad", "Fresh mung sprouts with lemon", true},
	}

	for _, tt := range tests {
		actual := MatchCategoryKeywords(tt.slug, tt.foodCat, tt.foodName, tt.foodDesc)
		if actual != tt.expected {
			t.Errorf("MatchCategoryKeywords(%q, %q, %q, %q) = %v; expected %v",
				tt.slug, tt.foodCat, tt.foodName, tt.foodDesc, actual, tt.expected)
		}
	}
}
