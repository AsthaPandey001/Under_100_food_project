import React from "react";
import { Tag, MapPin, Receipt } from "lucide-react";

export const BenefitsStrip: React.FC = () => {
  return (
    <section style={{
      backgroundColor: "#F8F6F0",
      borderBottom: "1px solid var(--color-border)",
      padding: "32px 0"
    }}>
      <div className="container" style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))",
        gap: "24px"
      }}>
        {/* Benefit 1 */}
        <div style={{
          display: "flex",
          alignItems: "flex-start",
          gap: "16px",
          padding: "16px 20px",
          backgroundColor: "#ffffff",
          borderRadius: "var(--radius-md)",
          border: "1px solid var(--color-border)",
          boxShadow: "var(--shadow-subtle)"
        }}>
          <div style={{
            width: "44px",
            height: "44px",
            borderRadius: "12px",
            backgroundColor: "var(--color-soft-orange-bg)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "var(--color-primary-orange)",
            flexShrink: 0
          }}>
            <Tag size={22} />
          </div>
          <div>
            <h3 style={{ fontSize: "16px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "3px" }}>
              Every item under ₹100
            </h3>
            <p style={{ fontSize: "13px", color: "var(--color-secondary-text)", lineHeight: 1.4 }}>
              Affordable picks, clearly priced.
            </p>
          </div>
        </div>

        {/* Benefit 2 */}
        <div style={{
          display: "flex",
          alignItems: "flex-start",
          gap: "16px",
          padding: "16px 20px",
          backgroundColor: "#ffffff",
          borderRadius: "var(--radius-md)",
          border: "1px solid var(--color-border)",
          boxShadow: "var(--shadow-subtle)"
        }}>
          <div style={{
            width: "44px",
            height: "44px",
            borderRadius: "12px",
            backgroundColor: "var(--color-soft-orange-bg)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "var(--color-primary-orange)",
            flexShrink: 0
          }}>
            <MapPin size={22} />
          </div>
          <div>
            <h3 style={{ fontSize: "16px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "3px" }}>
              Local Bengaluru favourites
            </h3>
            <p style={{ fontSize: "13px", color: "var(--color-secondary-text)", lineHeight: 1.4 }}>
              Discover food around your neighbourhood.
            </p>
          </div>
        </div>

        {/* Benefit 3 */}
        <div style={{
          display: "flex",
          alignItems: "flex-start",
          gap: "16px",
          padding: "16px 20px",
          backgroundColor: "#ffffff",
          borderRadius: "var(--radius-md)",
          border: "1px solid var(--color-border)",
          boxShadow: "var(--shadow-subtle)"
        }}>
          <div style={{
            width: "44px",
            height: "44px",
            borderRadius: "12px",
            backgroundColor: "var(--color-soft-orange-bg)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "var(--color-primary-orange)",
            flexShrink: 0
          }}>
            <Receipt size={22} />
          </div>
          <div>
            <h3 style={{ fontSize: "16px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "3px" }}>
              Prices shown upfront
            </h3>
            <p style={{ fontSize: "13px", color: "var(--color-secondary-text)", lineHeight: 1.4 }}>
              Know the item price before adding it.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
