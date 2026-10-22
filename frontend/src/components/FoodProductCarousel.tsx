import React, { useRef, useState, useEffect } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { FoodItem } from "../data/foodItems";
import { ProductCard } from "./ProductCard";

interface FoodProductCarouselProps {
  id: string;
  heading: string;
  supportingText?: string;
  items: FoodItem[];
  getItemQuantity: (id: string) => number;
  onAdd: (item: FoodItem) => void;
  onUpdateQuantity: (id: string, delta: number) => void;
}

export const FoodProductCarousel: React.FC<FoodProductCarouselProps> = ({
  id,
  heading,
  supportingText,
  items,
  getItemQuantity,
  onAdd,
  onUpdateQuantity
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
  }, [items]);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const offset = direction === "left" ? -340 : 340;
      scrollRef.current.scrollBy({ left: offset, behavior: "smooth" });
      setTimeout(checkScroll, 350);
    }
  };

  if (items.length === 0) return null;

  return (
    <section id={id} style={{ padding: "28px 0 24px", borderBottom: "1px solid var(--color-border)", backgroundColor: "var(--color-main-bg)" }}>
      <div className="container">
        {/* Header with Title & Arrow Controls */}
        <div style={{
          display: "flex",
          alignItems: "flex-end",
          justifyContent: "space-between",
          marginBottom: "16px"
        }}>
          <div>
            <h2 style={{
              fontSize: "20px",
              fontWeight: 800,
              color: "var(--color-dark-text)",
              letterSpacing: "-0.3px",
              marginBottom: supportingText ? "4px" : "0"
            }}>
              {heading}
            </h2>
            {supportingText && (
              <p style={{ fontSize: "14px", color: "var(--color-secondary-text)" }}>
                {supportingText}
              </p>
            )}
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
              aria-label={`Scroll ${heading} left`}
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
              aria-label={`Scroll ${heading} right`}
            >
              <ChevronRight size={18} />
            </button>
          </div>
        </div>

        {/* Scrollable Products Track */}
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
          {items.map((item) => (
            <div
              key={item.id}
              style={{
                flex: "0 0 240px",
                scrollSnapAlign: "start"
              }}
            >
              <ProductCard
                item={item}
                quantity={getItemQuantity(item.id)}
                onAdd={onAdd}
                onUpdateQuantity={onUpdateQuantity}
                isCompact={true}
              />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
