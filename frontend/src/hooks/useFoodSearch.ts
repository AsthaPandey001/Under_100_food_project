import { useState, useMemo } from "react";
import { FOOD_ITEMS } from "../data/foodItems";

export function useFoodSearch() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("all");
  const [selectedNeighbourhood, setSelectedNeighbourhood] = useState("all");
  const [sortBy, setSortBy] = useState<"popular" | "price-asc" | "rating" | "delivery">("popular");

  const filteredItems = useMemo(() => {
    return FOOD_ITEMS.filter(item => {
      // Category filter
      if (selectedCategory !== "all" && item.category !== selectedCategory) {
        return false;
      }

      // Neighbourhood filter
      if (selectedNeighbourhood !== "all" && item.neighbourhood !== selectedNeighbourhood) {
        return false;
      }

      // Search query filter
      if (searchQuery.trim() !== "") {
        const q = searchQuery.toLowerCase().trim();
        const matchesName = item.name.toLowerCase().includes(q);
        const matchesRestaurant = item.restaurant.toLowerCase().includes(q);
        const matchesCuisine = item.cuisine.toLowerCase().includes(q);
        const matchesNeighbourhood = item.neighbourhood.toLowerCase().includes(q);
        if (!matchesName && !matchesRestaurant && !matchesCuisine && !matchesNeighbourhood) {
          return false;
        }
      }

      return true;
    }).sort((a, b) => {
      if (sortBy === "price-asc") return a.price - b.price;
      if (sortBy === "rating") return b.rating - a.rating;
      if (sortBy === "delivery") return a.deliveryMinutes - b.deliveryMinutes;
      // Default: popularity / ratingCount
      return b.ratingCount - a.ratingCount;
    });
  }, [searchQuery, selectedCategory, selectedNeighbourhood, sortBy]);

  const resetFilters = () => {
    setSearchQuery("");
    setSelectedCategory("all");
    setSelectedNeighbourhood("all");
    setSortBy("popular");
  };

  return {
    searchQuery,
    setSearchQuery,
    selectedCategory,
    setSelectedCategory,
    selectedNeighbourhood,
    setSelectedNeighbourhood,
    sortBy,
    setSortBy,
    filteredItems,
    resetFilters
  };
}
