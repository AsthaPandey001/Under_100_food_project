import React from "react";
import { BadgePercent, Eye, Compass, Smartphone } from "lucide-react";

export const WhyUnder100: React.FC = () => {
  return (
    <section style={{
      padding: "64px 0",
      backgroundColor: "var(--color-main-bg)",
      borderBottom: "1px solid var(--color-border)"
    }}>
      <div className="container">
        <div style={{ textAlign: "center", maxWidth: "600px", margin: "0 auto 48px" }}>
          <div style={{
            fontSize: "12px",
            fontWeight: 800,
            color: "var(--color-deep-green)",
            textTransform: "uppercase",
            letterSpacing: "0.5px",
            marginBottom: "6px"
          }}>
            Built For Bengaluru
          </div>
          <h2 style={{
            fontSize: "30px",
            fontWeight: 900,
            color: "var(--color-dark-text)",
            letterSpacing: "-0.5px",
            marginBottom: "10px"
          }}>
            Why Under100?
          </h2>
          <p style={{ fontSize: "15px", color: "var(--color-secondary-text)" }}>
            A food platform crafted around honest prices and neighborhood comfort.
          </p>
        </div>

        <div style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))",
          gap: "24px"
        }}>
          {/* Card 1 */}
          <div style={{
            padding: "24px",
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
              color: "var(--color-primary-orange)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginBottom: "16px"
            }}>
              <BadgePercent size={22} />
            </div>
            <h3 style={{ fontSize: "17px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Budget-first food discovery
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              Instead of filtering out expensive dishes, everything on Under100 is capped at ₹99 by design.
            </p>
          </div>

          {/* Card 2 */}
          <div style={{
            padding: "24px",
            backgroundColor: "#ffffff",
            borderRadius: "var(--radius-md)",
            border: "1px solid var(--color-border)",
            boxShadow: "var(--shadow-subtle)"
          }}>
            <div style={{
              width: "44px",
              height: "44px",
              borderRadius: "12px",
              backgroundColor: "var(--color-soft-green-bg)",
              color: "var(--color-deep-green)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginBottom: "16px"
            }}>
              <Eye size={22} />
            </div>
            <h3 style={{ fontSize: "17px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Transparent item pricing
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              The price you see on the food card is the item price you pay. No inflated menu prices or hidden surcharges.
            </p>
          </div>

          {/* Card 3 */}
          <div style={{
            padding: "24px",
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
              color: "var(--color-primary-orange)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginBottom: "16px"
            }}>
              <Compass size={22} />
            </div>
            <h3 style={{ fontSize: "17px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Local restaurant discovery
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              Highlighting authentic Bengaluru darshinis, tiffin rooms, rolls joints, and local momo stalls.
            </p>
          </div>

          {/* Card 4 */}
          <div style={{
            padding: "24px",
            backgroundColor: "#ffffff",
            borderRadius: "var(--radius-md)",
            border: "1px solid var(--color-border)",
            boxShadow: "var(--shadow-subtle)"
          }}>
            <div style={{
              width: "44px",
              height: "44px",
              borderRadius: "12px",
              backgroundColor: "var(--color-soft-green-bg)",
              color: "var(--color-deep-green)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginBottom: "16px"
            }}>
              <Smartphone size={22} />
            </div>
            <h3 style={{ fontSize: "17px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Easy browsing on mobile
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              One-hand friendly horizontal carousels, responsive touch swipe, and lightning-fast dish search.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
