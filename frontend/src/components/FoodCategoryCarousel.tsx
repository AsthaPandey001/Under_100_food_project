import React, { useRef, useState, useEffect } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { CATEGORIES, Category } from "../data/categories";

interface FoodCategoryCarouselProps {
  selectedCategory: string;
  onSelectCategory: (catId: string) => void;
}

export const FoodCategoryCarousel: React.FC<FoodCategoryCarouselProps> = ({
  selectedCategory,
  onSelectCategory
}) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(true);

  const checkScroll = () => {
    if (scrollRef.current) {
      const { scrollLeft, scrollWidth, clientWidth } = scrollRef.current;
      setCanScrollLeft(scrollLeft > 10);
      setCanScrollRight(scrollLeft < scrollWidth - clientWidth - 10);
    }
  };

  useEffect(() => {
    checkScroll();
    window.addEventListener("resize", checkScroll);
    return () => window.removeEventListener("resize", checkScroll);
  }, []);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const offset = direction === "left" ? -320 : 320;
      scrollRef.current.scrollBy({ left: offset, behavior: "smooth" });
      setTimeout(checkScroll, 350);
    }
  };

  return (
    <section style={{ padding: "28px 0 20px", borderBottom: "1px solid var(--color-border)", backgroundColor: "#ffffff" }}>
      <div className="container">
        {/* Header with Title and Arrows */}
        <div style={{
          display: "flex",
          alignItems: "flex-end",
          justifyContent: "space-between",
          marginBottom: "16px"
        }}>
          <div>
            <h2 style={{
              fontSize: "22px",
              fontWeight: 800,
              color: "var(--color-dark-text)",
              letterSpacing: "-0.3px",
              marginBottom: "4px"
            }}>
              What are you craving today?
            </h2>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)" }}>
              Your favourites, without the big bill.
            </p>
          </div>

          {/* Desktop Arrow Controls */}
          <div style={{ display: "flex", gap: "8px" }}>
            <button
              type="button"
              onClick={() => scroll("left")}
              disabled={!canScrollLeft}
              style={{
                width: "36px",
                height: "36px",
                borderRadius: "50%",
                border: "1px solid var(--color-border)",
                backgroundColor: "#ffffff",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: canScrollLeft ? "var(--color-dark-text)" : "#CBD5E1",
                cursor: canScrollLeft ? "pointer" : "not-allowed",
                boxShadow: "var(--shadow-subtle)",
                transition: "all 0.2s ease"
              }}
              aria-label="Scroll cravings left"
            >
              <ChevronLeft size={18} />
            </button>
            <button
              type="button"
              onClick={() => scroll("right")}
              disabled={!canScrollRight}
              style={{
                width: "36px",
                height: "36px",
                borderRadius: "50%",
                border: "1px solid var(--color-border)",
                backgroundColor: "#ffffff",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: canScrollRight ? "var(--color-dark-text)" : "#CBD5E1",
                cursor: canScrollRight ? "pointer" : "not-allowed",
                boxShadow: "var(--shadow-subtle)",
                transition: "all 0.2s ease"
              }}
              aria-label="Scroll cravings right"
            >
              <ChevronRight size={18} />
            </button>
          </div>
        </div>

        {/* Scrollable Track */}
        <div
          ref={scrollRef}
          onScroll={checkScroll}
          className="hide-scrollbar"
          style={{
            display: "flex",
            gap: "16px",
            overflowX: "auto",
            scrollSnapType: "x mandatory",
            paddingBottom: "8px"
          }}
        >
          {CATEGORIES.map((cat: Category) => {
            const isSelected = selectedCategory === cat.id;
            return (
              <button
                key={cat.id}
                type="button"
                onClick={() => onSelectCategory(isSelected ? "all" : cat.id)}
                style={{
                  flex: "0 0 130px",
                  scrollSnapAlign: "start",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  padding: "12px 8px 14px",
                  borderRadius: "var(--radius-md)",
                  backgroundColor: isSelected ? "var(--color-soft-orange-bg)" : "var(--color-card-bg)",
                  border: isSelected ? "2px solid var(--color-primary-orange)" : "1px solid var(--color-border)",
                  boxShadow: "var(--shadow-subtle)",
                  cursor: "pointer",
                  transition: "all 0.2s ease",
                  textAlign: "center"
                }}
                onMouseEnter={(e) => {
                  if (!isSelected) {
                    e.currentTarget.style.transform = "translateY(-3px)";
                    e.currentTarget.style.boxShadow = "var(--shadow-card)";
                  }
                }}
                onMouseLeave={(e) => {
                  if (!isSelected) {
                    e.currentTarget.style.transform = "none";
                    e.currentTarget.style.boxShadow = "var(--shadow-subtle)";
                  }
                }}
              >
                {/* Category Image Circle */}
                <div style={{
                  width: "82px",
                  height: "82px",
                  borderRadius: "50%",
                  overflow: "hidden",
                  marginBottom: "10px",
                  border: "2px solid #ffffff",
                  boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
                  backgroundColor: "#F1EFEA"
                }}>
                  <img
                    src={cat.image}
                    alt={cat.name}
                    loading="lazy"
                    style={{
                      width: "100%",
                      height: "100%",
                      objectFit: "cover"
                    }}
                    onError={(e) => {
                      (e.target as HTMLImageElement).src = "/crispy_dosa.jpg";
                    }}
                  />
                </div>

                {/* Title */}
                <div style={{
                  fontSize: "13px",
                  fontWeight: 700,
                  color: isSelected ? "var(--color-primary-orange)" : "var(--color-dark-text)",
                  lineHeight: 1.2,
                  marginBottom: "4px"
                }}>
                  {cat.name}
                </div>

                {/* Starting Price Pill */}
                <div style={{
                  fontSize: "11px",
                  fontWeight: 600,
                  color: "var(--color-deep-green)",
                  backgroundColor: "var(--color-soft-green-bg)",
                  padding: "2px 8px",
                  borderRadius: "var(--radius-full)"
                }}>
                  From ₹{cat.startingPrice}
                </div>
              </button>
            );
          })}
        </div>
      </div>
    </section>
  );
};
