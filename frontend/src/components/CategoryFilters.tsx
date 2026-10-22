import React from "react";
import { SlidersHorizontal } from "lucide-react";
import { CATEGORIES } from "../data/categories";

interface CategoryFiltersProps {
  selectedCategory: string;
  onSelectCategory: (cat: string) => void;
  sortBy: "popular" | "price-asc" | "rating" | "delivery";
  onSortChange: (sort: "popular" | "price-asc" | "rating" | "delivery") => void;
  resultCount: number;
}

export const CategoryFilters: React.FC<CategoryFiltersProps> = ({
  selectedCategory,
  onSelectCategory,
  sortBy,
  onSortChange,
  resultCount
}) => {
  return (
    <div style={{
      display: "flex",
      alignItems: "center",
      justifyContent: "space-between",
      flexWrap: "wrap",
      gap: "16px",
      marginBottom: "24px"
    }}>
      {/* Category Chips */}
      <div style={{
        display: "flex",
        alignItems: "center",
        gap: "8px",
        overflowX: "auto",
        paddingBottom: "4px",
        maxWidth: "100%",
        scrollbarWidth: "none"
      }} className="hide-scrollbar">
        <button
          type="button"
          onClick={() => onSelectCategory("all")}
          style={{
            padding: "8px 16px",
            borderRadius: "var(--radius-full)",
            fontSize: "13px",
            fontWeight: 700,
            cursor: "pointer",
            border: "1px solid",
            borderColor: selectedCategory === "all" ? "var(--color-primary-orange)" : "var(--color-border)",
            backgroundColor: selectedCategory === "all" ? "var(--color-primary-orange)" : "#ffffff",
            color: selectedCategory === "all" ? "#ffffff" : "var(--color-dark-text)",
            transition: "all 0.2s ease",
            whiteSpace: "nowrap"
          }}
        >
          All Items ({resultCount})
        </button>

        {CATEGORIES.map((cat) => {
          const active = selectedCategory === cat.id;
          return (
            <button
              key={cat.id}
              type="button"
              onClick={() => onSelectCategory(cat.id)}
              style={{
                padding: "8px 14px",
                borderRadius: "var(--radius-full)",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
                border: "1px solid",
                borderColor: active ? "var(--color-primary-orange)" : "var(--color-border)",
                backgroundColor: active ? "var(--color-soft-orange-bg)" : "#ffffff",
                color: active ? "var(--color-primary-orange)" : "var(--color-dark-text)",
                transition: "all 0.2s ease",
                whiteSpace: "nowrap"
              }}
            >
              {cat.name}
            </button>
          );
        })}
      </div>

      {/* Sort Controls */}
      <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
        <SlidersHorizontal size={16} color="var(--color-secondary-text)" />
        <label htmlFor="sort-select" style={{ fontSize: "13px", color: "var(--color-secondary-text)", fontWeight: 500 }}>
          Sort by:
        </label>
        <select
          id="sort-select"
          value={sortBy}
          onChange={(e) => onSortChange(e.target.value as any)}
          style={{
            padding: "8px 12px",
            borderRadius: "var(--radius-md)",
            border: "1px solid var(--color-border)",
            backgroundColor: "#ffffff",
            fontSize: "13px",
            fontWeight: 600,
            color: "var(--color-dark-text)",
            outline: "none",
            cursor: "pointer"
          }}
        >
          <option value="popular">Most Popular</option>
          <option value="price-asc">Price: Low to High</option>
          <option value="rating">Top Rated</option>
          <option value="delivery">Fastest Delivery</option>
        </select>
      </div>
    </div>
  );
};
