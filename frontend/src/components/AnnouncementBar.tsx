import React from "react";
import { ChevronRight } from "lucide-react";

export const AnnouncementBar: React.FC = () => {
  return (
    <div style={{
      backgroundColor: "#EBF7EE",
      borderBottom: "1px solid #D5EED9",
      padding: "4px 24px",
      fontSize: "12.5px",
      color: "#1E3A24",
      display: "flex",
      alignItems: "center",
      justifyContent: "space-between",
      minHeight: "38px",
      position: "relative",
      overflow: "hidden"
    }}>
      <div className="container" style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        width: "100%",
        gap: "16px"
      }}>
        {/* Left: Leaf + Primary Tagline + Bullets */}
        <div style={{ display: "flex", alignItems: "center", gap: "10px", flexWrap: "wrap" }}>
          {/* Green Leaf Icon */}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="#087443" style={{ flexShrink: 0 }}>
            <path d="M17.72 2.28C11.5 2.5 5.5 8.5 4.5 15.5c-0.3 2.1 0.5 4.2 2 5.5 1.5 1.3 3.5 1.8 5.5 1.2 6.5-1.8 11.5-8.2 11-15.2-0.2-2.3-1.8-4.3-4-4.7-0.4-0.1-0.9-0.1-1.28 0zm-2.4 14.5c-3.1 0.7-6.2-1.2-6.9-4.3s1.2-6.2 4.3-6.9c0.7-0.2 1.4-0.2 2.1-0.1-0.6 3.8-2.5 7.1-5.5 9.3 1.9 1.5 4.6 1.8 6.8 0.8-0.2 0.4-0.5 0.8-0.8 1.2z"/>
          </svg>

          <span style={{ fontWeight: 800, color: "#0B4B27", letterSpacing: "-0.1px" }}>
            Bengaluru's everyday meals, under ₹100
          </span>

          <span style={{ color: "#8EBE98", fontWeight: 300, margin: "0 2px" }}>|</span>

          <div style={{ display: "flex", alignItems: "center", gap: "8px", color: "#2B5736", fontSize: "12px", fontWeight: 500 }}>
            <span>Good food</span>
            <span style={{ color: "#77A882" }}>•</span>
            <span>Local favourites</span>
            <span style={{ color: "#77A882" }}>•</span>
            <span>No hidden fees</span>
          </div>
        </div>

        {/* Right: Bengaluru Skyline Silhouette + Namma Bengaluru Link */}
        <div style={{ display: "flex", alignItems: "center", gap: "18px" }}>
          {/* Enhanced Bengaluru Heritage Skyline Silhouette */}
          <div className="desktop-only" style={{ display: "flex", alignItems: "flex-end", height: "32px" }}>
            <img
              src="/bengaluru_skyline.svg"
              alt="Bengaluru Vidhana Soudha, High Court & Heritage Skyline"
              style={{
                height: "32px",
                width: "auto",
                maxWidth: "230px",
                objectFit: "contain",
                display: "block"
              }}
            />
          </div>

          {/* Red Heart + Namma Bengaluru > */}
          <a
            href="#food-discovery"
            style={{
              display: "flex",
              alignItems: "center",
              gap: "6px",
              color: "#087443",
              fontWeight: 700,
              fontSize: "12.5px",
              textDecoration: "none",
              whiteSpace: "nowrap"
            }}
          >
            {/* Red Heart */}
            <svg width="14" height="14" viewBox="0 0 24 24" fill="#E11D48">
              <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
            </svg>
            <span style={{ color: "#143F20" }}>Namma Bengaluru</span>
            <ChevronRight size={13} color="#087443" strokeWidth={2.5} />
          </a>
        </div>
      </div>
    </div>
  );
};
