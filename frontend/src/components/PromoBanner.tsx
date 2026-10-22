import React from "react";
import { CheckCircle2, HeartHandshake, Sparkles, Store } from "lucide-react";

export const PromoBanner: React.FC = () => {
  return (
    <section style={{
      padding: "40px 0",
      backgroundColor: "var(--color-main-bg)"
    }}>
      <div className="container">
        <div style={{
          backgroundColor: "var(--color-soft-green-bg)",
          borderRadius: "var(--radius-lg)",
          border: "1px solid rgba(8, 116, 67, 0.15)",
          padding: "48px 40px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          flexWrap: "wrap",
          gap: "32px",
          boxShadow: "var(--shadow-subtle)",
          position: "relative",
          overflow: "hidden"
        }}>
          {/* Left Text */}
          <div style={{ maxWidth: "600px", zIndex: 1 }}>
            <div style={{
              display: "inline-flex",
              alignItems: "center",
              gap: "6px",
              backgroundColor: "#ffffff",
              color: "var(--color-deep-green)",
              padding: "4px 12px",
              borderRadius: "var(--radius-full)",
              fontSize: "12px",
              fontWeight: 800,
              textTransform: "uppercase",
              letterSpacing: "0.5px",
              marginBottom: "16px",
              border: "1px solid rgba(8, 116, 67, 0.2)"
            }}>
              <Sparkles size={14} />
              <span>Pocket-Friendly Promise</span>
            </div>

            <h2 style={{
              fontSize: "clamp(26px, 4vw, 36px)",
              fontWeight: 900,
              color: "var(--color-dark-text)",
              letterSpacing: "-0.8px",
              lineHeight: 1.2,
              marginBottom: "12px"
            }}>
              Your daily food budget just got happier.
            </h2>

            <p style={{
              fontSize: "16px",
              color: "var(--color-secondary-text)",
              lineHeight: 1.6,
              marginBottom: "28px"
            }}>
              Discover satisfying meals and local favourites for under ₹100.
            </p>

            {/* Three Compact Benefit Labels */}
            <div style={{
              display: "flex",
              flexWrap: "wrap",
              gap: "16px"
            }}>
              <div style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                backgroundColor: "#ffffff",
                padding: "8px 14px",
                borderRadius: "var(--radius-full)",
                fontSize: "13px",
                fontWeight: 700,
                color: "var(--color-deep-green)",
                boxShadow: "0 2px 6px rgba(8,116,67,0.08)"
              }}>
                <CheckCircle2 size={16} />
                <span>Tasty meals under ₹100</span>
              </div>

              <div style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                backgroundColor: "#ffffff",
                padding: "8px 14px",
                borderRadius: "var(--radius-full)",
                fontSize: "13px",
                fontWeight: 700,
                color: "var(--color-deep-green)",
                boxShadow: "0 2px 6px rgba(8,116,67,0.08)"
              }}>
                <Store size={16} />
                <span>Support local eateries</span>
              </div>

              <div style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                backgroundColor: "#ffffff",
                padding: "8px 14px",
                borderRadius: "var(--radius-full)",
                fontSize: "13px",
                fontWeight: 700,
                color: "var(--color-deep-green)",
                boxShadow: "0 2px 6px rgba(8,116,67,0.08)"
              }}>
                <HeartHandshake size={16} />
                <span>Real food. Real people.</span>
              </div>
            </div>
          </div>

          {/* Right Visual Image */}
          <div style={{
            width: "220px",
            height: "220px",
            borderRadius: "var(--radius-lg)",
            overflow: "hidden",
            boxShadow: "0 8px 24px rgba(0,0,0,0.1)",
            border: "3px solid #ffffff",
            flexShrink: 0
          }}>
            <img
              src="/crispy_dosa.jpg"
              alt="Local Bengaluru Street Food"
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          </div>
        </div>
      </div>
    </section>
  );
};
