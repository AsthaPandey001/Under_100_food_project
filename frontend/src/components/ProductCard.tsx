import React from "react";
import { Star, Clock, Plus, Minus, MapPin } from "lucide-react";
import { FoodItem } from "../data/foodItems";

interface ProductCardProps {
  item: FoodItem;
  quantity: number;
  onAdd: (item: FoodItem) => void;
  onUpdateQuantity: (id: string, delta: number) => void;
  isCompact?: boolean;
}

export const ProductCard: React.FC<ProductCardProps> = ({
  item,
  quantity,
  onAdd,
  onUpdateQuantity,
  isCompact = false
}) => {
  return (
    <div
      style={{
        borderRadius: "var(--radius-md)",
        backgroundColor: "var(--color-card-bg)",
        border: "1px solid var(--color-border)",
        boxShadow: "var(--shadow-subtle)",
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
        transition: "all 0.2s ease",
        height: "100%"
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = "translateY(-4px)";
        e.currentTarget.style.boxShadow = "var(--shadow-card)";
        e.currentTarget.style.borderColor = "var(--color-border-hover)";
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = "none";
        e.currentTarget.style.boxShadow = "var(--shadow-subtle)";
        e.currentTarget.style.borderColor = "var(--color-border)";
      }}
    >
      {/* Image Container */}
      <div style={{ position: "relative", width: "100%", height: isCompact ? "140px" : "170px", backgroundColor: "#F3F1EC", overflow: "hidden" }}>
        <img
          src={item.image}
          alt={item.name}
          loading="lazy"
          style={{ width: "100%", height: "100%", objectFit: "cover", transition: "transform 0.3s ease" }}
          onError={(e) => {
            (e.target as HTMLImageElement).src = "/crispy_dosa.jpg";
          }}
        />

        {/* Veg / Non-Veg Indicator */}
        <div style={{
          position: "absolute",
          top: "10px",
          left: "10px",
          backgroundColor: "#ffffff",
          borderRadius: "4px",
          padding: "3px",
          boxShadow: "0 2px 6px rgba(0,0,0,0.12)"
        }}>
          <div className={`diet-indicator ${item.isVegetarian ? "diet-veg" : "diet-nonveg"}`}>
            <div className="diet-dot" />
          </div>
        </div>

        {/* Price Pill */}
        <div style={{
          position: "absolute",
          bottom: "10px",
          left: "10px",
          backgroundColor: "rgba(23, 28, 36, 0.9)",
          backdropFilter: "blur(4px)",
          color: "#ffffff",
          padding: "4px 8px",
          borderRadius: "var(--radius-sm)",
          fontSize: "13px",
          fontWeight: 800,
          letterSpacing: "-0.2px"
        }}>
          ₹{item.price}
        </div>

        {/* Delivery Time Badge */}
        <div style={{
          position: "absolute",
          bottom: "10px",
          right: "10px",
          backgroundColor: "#ffffff",
          color: "var(--color-dark-text)",
          padding: "3px 8px",
          borderRadius: "var(--radius-sm)",
          fontSize: "11px",
          fontWeight: 600,
          display: "flex",
          alignItems: "center",
          gap: "4px",
          boxShadow: "0 2px 6px rgba(0,0,0,0.1)"
        }}>
          <Clock size={12} color="var(--color-secondary-text)" />
          <span>{item.deliveryMinutes}m</span>
        </div>
      </div>

      {/* Content */}
      <div style={{ padding: "14px", display: "flex", flexDirection: "column", flex: 1 }}>
        <div style={{ display: "flex", alignItems: "flex-start", justifyContent: "space-between", gap: "8px", marginBottom: "4px" }}>
          <h3 style={{
            fontSize: "15px",
            fontWeight: 700,
            color: "var(--color-dark-text)",
            lineHeight: 1.3,
            overflow: "hidden",
            display: "-webkit-box",
            WebkitLineClamp: 1,
            WebkitBoxOrient: "vertical"
          }}>
            {item.name}
          </h3>
          <div style={{
            display: "flex",
            alignItems: "center",
            gap: "3px",
            fontSize: "12px",
            fontWeight: 700,
            color: "#B45309",
            backgroundColor: "#FEF3C7",
            padding: "2px 6px",
            borderRadius: "4px",
            flexShrink: 0
          }}>
            <Star size={11} fill="#B45309" stroke="#B45309" />
            <span>{item.rating}</span>
          </div>
        </div>

        {/* Restaurant & Neighbourhood */}
        <div style={{ fontSize: "12px", color: "var(--color-secondary-text)", marginBottom: "4px", fontWeight: 500 }}>
          {item.restaurant}
        </div>
        <div style={{
          display: "flex",
          alignItems: "center",
          gap: "3px",
          fontSize: "11px",
          color: "var(--color-secondary-text)",
          marginBottom: "12px"
        }}>
          <MapPin size={11} color="var(--color-primary-orange)" />
          <span>{item.neighbourhood}</span>
          <span>•</span>
          <span>{item.cuisine}</span>
        </div>

        {/* Price & Add to Cart Controls */}
        <div style={{
          marginTop: "auto",
          paddingTop: "10px",
          borderTop: "1px dashed var(--color-border)",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between"
        }}>
          <div>
            <span style={{ fontSize: "11px", color: "var(--color-secondary-text)", display: "block", lineHeight: 1 }}>Item Price</span>
            <span style={{ fontSize: "18px", fontWeight: 900, color: "var(--color-dark-text)" }}>₹{item.price}</span>
          </div>

          {quantity === 0 ? (
            <button
              type="button"
              onClick={() => onAdd(item)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "4px",
                padding: "8px 16px",
                borderRadius: "var(--radius-full)",
                backgroundColor: "#ffffff",
                border: "1.5px solid var(--color-deep-green)",
                color: "var(--color-deep-green)",
                fontSize: "13px",
                fontWeight: 700,
                transition: "all 0.2s ease",
                boxShadow: "0 1px 4px rgba(8, 116, 67, 0.08)"
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = "var(--color-deep-green)";
                e.currentTarget.style.color = "#ffffff";
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = "#ffffff";
                e.currentTarget.style.color = "var(--color-deep-green)";
              }}
            >
              <Plus size={14} strokeWidth={2.5} />
              <span>ADD</span>
            </button>
          ) : (
            <div style={{
              display: "flex",
              alignItems: "center",
              backgroundColor: "var(--color-deep-green)",
              borderRadius: "var(--radius-full)",
              color: "#ffffff",
              overflow: "hidden",
              boxShadow: "0 2px 8px rgba(8, 116, 67, 0.2)"
            }}>
              <button
                type="button"
                onClick={() => onUpdateQuantity(item.id, -1)}
                style={{
                  padding: "6px 10px",
                  color: "#ffffff",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center"
                }}
                aria-label={`Decrease ${item.name}`}
              >
                <Minus size={14} strokeWidth={2.5} />
              </button>
              <span style={{
                padding: "0 6px",
                fontSize: "13px",
                fontWeight: 800,
                minWidth: "20px",
                textAlign: "center"
              }}>
                {quantity}
              </span>
              <button
                type="button"
                onClick={() => onUpdateQuantity(item.id, 1)}
                style={{
                  padding: "6px 10px",
                  color: "#ffffff",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center"
                }}
                aria-label={`Increase ${item.name}`}
              >
                <Plus size={14} strokeWidth={2.5} />
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
