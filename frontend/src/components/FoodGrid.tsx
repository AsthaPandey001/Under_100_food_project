import React from "react";
import { Utensils, RotateCcw } from "lucide-react";
import { FoodItem } from "../data/foodItems";
import { ProductCard } from "./ProductCard";
import { CategoryFilters } from "./CategoryFilters";

interface FoodGridProps {
  items: FoodItem[];
  selectedCategory: string;
  onSelectCategory: (cat: string) => void;
  sortBy: "popular" | "price-asc" | "rating" | "delivery";
  onSortChange: (sort: "popular" | "price-asc" | "rating" | "delivery") => void;
  selectedNeighbourhood: string;
  searchQuery: string;
  onResetFilters: () => void;
  getItemQuantity: (id: string) => number;
  onAdd: (item: FoodItem) => void;
  onUpdateQuantity: (id: string, delta: number) => void;
}

export const FoodGrid: React.FC<FoodGridProps> = ({
  items,
  selectedCategory,
  onSelectCategory,
  sortBy,
  onSortChange,
  selectedNeighbourhood,
  searchQuery,
  onResetFilters,
  getItemQuantity,
  onAdd,
  onUpdateQuantity
}) => {
  return (
    <section id="food-discovery" style={{ padding: "48px 0 60px", backgroundColor: "var(--color-main-bg)" }}>
      <div className="container">
        {/* Section Heading */}
        <div style={{
          display: "flex",
          alignItems: "flex-end",
          justifyContent: "space-between",
          flexWrap: "wrap",
          gap: "12px",
          marginBottom: "24px"
        }}>
          <div>
            <div style={{
              fontSize: "12px",
              fontWeight: 800,
              color: "var(--color-deep-green)",
              textTransform: "uppercase",
              letterSpacing: "0.5px",
              marginBottom: "4px"
            }}>
              {selectedNeighbourhood === "all" ? "Across Bengaluru" : `In ${selectedNeighbourhood}`}
            </div>
            <h2 style={{
              fontSize: "28px",
              fontWeight: 900,
              color: "var(--color-dark-text)",
              letterSpacing: "-0.5px"
            }}>
              Good food. Under ₹100.
            </h2>
            <p style={{ fontSize: "15px", color: "var(--color-secondary-text)", marginTop: "4px" }}>
              Affordable favourites near you.
            </p>
          </div>

          {(selectedCategory !== "all" || searchQuery !== "" || selectedNeighbourhood !== "all") && (
            <button
              type="button"
              onClick={onResetFilters}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "6px",
                fontSize: "13px",
                fontWeight: 600,
                color: "var(--color-primary-orange)",
                padding: "6px 12px",
                borderRadius: "var(--radius-sm)",
                backgroundColor: "var(--color-soft-orange-bg)"
              }}
            >
              <RotateCcw size={14} />
              <span>Reset filters</span>
            </button>
          )}
        </div>

        {/* Filter Chips & Sorting */}
        <CategoryFilters
          selectedCategory={selectedCategory}
          onSelectCategory={onSelectCategory}
          sortBy={sortBy}
          onSortChange={onSortChange}
          resultCount={items.length}
        />

        {/* Product Grid */}
        {items.length > 0 ? (
          <div style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))",
            gap: "24px"
          }}>
            {items.map((item) => (
              <ProductCard
                key={item.id}
                item={item}
                quantity={getItemQuantity(item.id)}
                onAdd={onAdd}
                onUpdateQuantity={onUpdateQuantity}
              />
            ))}
          </div>
        ) : (
          /* Empty State */
          <div style={{
            textAlign: "center",
            padding: "64px 20px",
            backgroundColor: "#ffffff",
            borderRadius: "var(--radius-lg)",
            border: "1px dashed var(--color-border)"
          }}>
            <div style={{
              width: "56px",
              height: "56px",
              borderRadius: "50%",
              backgroundColor: "var(--color-soft-orange-bg)",
              color: "var(--color-primary-orange)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              margin: "0 auto 16px"
            }}>
              <Utensils size={28} />
            </div>
            <h3 style={{ fontSize: "18px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              No dishes found matching your criteria
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", maxWidth: "420px", margin: "0 auto 20px" }}>
              Try searching for something else like "dosa", "momos", or "roll", or reset your current filters.
            </p>
            <button
              type="button"
              onClick={onResetFilters}
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "8px",
                padding: "10px 20px",
                backgroundColor: "var(--color-deep-green)",
                color: "#ffffff",
                borderRadius: "var(--radius-full)",
                fontSize: "14px",
                fontWeight: 700
              }}
            >
              <RotateCcw size={16} />
              <span>Show All Meals</span>
            </button>
          </div>
        )}
      </div>
    </section>
  );
};
