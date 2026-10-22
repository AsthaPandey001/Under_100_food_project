import React from "react";
import { MapPin, Search, ShoppingBag } from "lucide-react";

export const HowItWorks: React.FC = () => {
  return (
    <section id="how-it-works" style={{
      padding: "64px 0",
      backgroundColor: "#ffffff",
      borderBottom: "1px solid var(--color-border)"
    }}>
      <div className="container">
        <div style={{ textAlign: "center", maxWidth: "600px", margin: "0 auto 48px" }}>
          <div style={{
            fontSize: "12px",
            fontWeight: 800,
            color: "var(--color-primary-orange)",
            textTransform: "uppercase",
            letterSpacing: "0.5px",
            marginBottom: "6px"
          }}>
            Simple & Transparent
          </div>
          <h2 style={{
            fontSize: "30px",
            fontWeight: 900,
            color: "var(--color-dark-text)",
            letterSpacing: "-0.5px",
            marginBottom: "10px"
          }}>
            How Under100 Works
          </h2>
          <p style={{ fontSize: "15px", color: "var(--color-secondary-text)" }}>
            Getting fresh everyday meals delivered without overpaying takes just 3 simple steps.
          </p>
        </div>

        <div style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(280px, 1fr))",
          gap: "32px"
        }}>
          {/* Step 1 */}
          <div style={{
            padding: "32px 24px",
            backgroundColor: "var(--color-main-bg)",
            borderRadius: "var(--radius-lg)",
            border: "1px solid var(--color-border)",
            textAlign: "center",
            position: "relative"
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
              margin: "0 auto 20px"
            }}>
              <MapPin size={26} />
            </div>
            <div style={{
              fontSize: "12px",
              fontWeight: 800,
              color: "var(--color-primary-orange)",
              marginBottom: "6px"
            }}>
              STEP 1
            </div>
            <h3 style={{ fontSize: "18px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Choose your neighbourhood
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              Select Indiranagar, Koramangala, HSR Layout, or other Bengaluru areas to see nearby active kitchens.
            </p>
          </div>

          {/* Step 2 */}
          <div style={{
            padding: "32px 24px",
            backgroundColor: "var(--color-main-bg)",
            borderRadius: "var(--radius-lg)",
            border: "1px solid var(--color-border)",
            textAlign: "center",
            position: "relative"
          }}>
            <div style={{
              width: "56px",
              height: "56px",
              borderRadius: "50%",
              backgroundColor: "var(--color-soft-green-bg)",
              color: "var(--color-deep-green)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              margin: "0 auto 20px"
            }}>
              <Search size={26} />
            </div>
            <div style={{
              fontSize: "12px",
              fontWeight: 800,
              color: "var(--color-deep-green)",
              marginBottom: "6px"
            }}>
              STEP 2
            </div>
            <h3 style={{ fontSize: "18px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Find a meal within your budget
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              Browse curated food carousels and dishes — every single listing is guaranteed to be under ₹100.
            </p>
          </div>

          {/* Step 3 */}
          <div style={{
            padding: "32px 24px",
            backgroundColor: "var(--color-main-bg)",
            borderRadius: "var(--radius-lg)",
            border: "1px solid var(--color-border)",
            textAlign: "center",
            position: "relative"
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
              margin: "0 auto 20px"
            }}>
              <ShoppingBag size={26} />
            </div>
            <div style={{
              fontSize: "12px",
              fontWeight: 800,
              color: "var(--color-primary-orange)",
              marginBottom: "6px"
            }}>
              STEP 3
            </div>
            <h3 style={{ fontSize: "18px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "8px" }}>
              Add favourites & review cart
            </h3>
            <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", lineHeight: 1.5 }}>
              Tap Add to customize portions with real-time transparent subtotal calculations and no surprise fees.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
