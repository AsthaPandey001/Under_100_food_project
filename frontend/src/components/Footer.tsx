import React from "react";
import { UtensilsCrossed, MapPin, Heart } from "lucide-react";

export const Footer: React.FC = () => {
  return (
    <footer style={{
      backgroundColor: "#171C24",
      color: "#94A3B8",
      padding: "64px 0 32px",
      borderTop: "1px solid #2D3748"
    }}>
      <div className="container">
        <div style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))",
          gap: "40px",
          marginBottom: "48px"
        }}>
          {/* Brand Col */}
          <div>
            <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "16px" }}>
              <div style={{
                width: "32px",
                height: "32px",
                borderRadius: "8px",
                backgroundColor: "var(--color-primary-orange)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#ffffff"
              }}>
                <UtensilsCrossed size={18} strokeWidth={2.5} />
              </div>
              <div style={{ display: "flex", alignItems: "baseline" }}>
                <span style={{ fontSize: "20px", fontWeight: 900, color: "#ffffff" }}>Under</span>
                <span style={{ fontSize: "20px", fontWeight: 900, color: "var(--color-primary-orange)" }}>100</span>
              </div>
            </div>
            <p style={{ fontSize: "14px", lineHeight: 1.6, color: "#94A3B8", marginBottom: "20px" }}>
              Bengaluru's dedicated food platform where customers discover satisfying local meals for under ₹100. Transparent prices, zero hidden fees.
            </p>
            <div style={{ display: "flex", alignItems: "center", gap: "6px", fontSize: "13px", color: "#E2E8F0" }}>
              <MapPin size={16} color="var(--color-primary-orange)" />
              <span>Service Area: Bengaluru, Karnataka</span>
            </div>
          </div>

          {/* About & Mission */}
          <div>
            <h4 style={{ fontSize: "15px", fontWeight: 800, color: "#ffffff", marginBottom: "16px", letterSpacing: "0.5px" }}>
              ABOUT UNDER100
            </h4>
            <ul style={{ listStyle: "none", display: "flex", flexDirection: "column", gap: "10px", fontSize: "14px" }}>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Our Mission</a></li>
              <li><a href="#how-it-works" style={{ color: "#94A3B8", transition: "color 0.2s" }}>How It Works</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Affordable Partner Kitchens</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Bengaluru Neighbourhoods</a></li>
            </ul>
          </div>

          {/* Help & Support */}
          <div>
            <h4 style={{ fontSize: "15px", fontWeight: 800, color: "#ffffff", marginBottom: "16px", letterSpacing: "0.5px" }}>
              HELP & SUPPORT
            </h4>
            <ul style={{ listStyle: "none", display: "flex", flexDirection: "column", gap: "10px", fontSize: "14px" }}>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Customer FAQs</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Restaurant Partner Support</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Order Issues & Feedback</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Contact Bengaluru Team</a></li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h4 style={{ fontSize: "15px", fontWeight: 800, color: "#ffffff", marginBottom: "16px", letterSpacing: "0.5px" }}>
              LEGAL & POLICIES
            </h4>
            <ul style={{ listStyle: "none", display: "flex", flexDirection: "column", gap: "10px", fontSize: "14px" }}>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Privacy Policy</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Terms of Service</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>Pricing Transparency Charter</a></li>
              <li><a href="#" style={{ color: "#94A3B8", transition: "color 0.2s" }}>FSSAI Standards</a></li>
            </ul>
          </div>
        </div>

        {/* Bottom Bar */}
        <div style={{
          borderTop: "1px solid #2D3748",
          paddingTop: "24px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          flexWrap: "wrap",
          gap: "16px",
          fontSize: "13px"
        }}>
          <div>
            © {new Date().getFullYear()} Under100 Technologies. All rights reserved.
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            <span>Crafted with</span>
            <Heart size={14} color="var(--color-primary-orange)" fill="var(--color-primary-orange)" />
            <span>for Bengaluru's everyday food lovers.</span>
          </div>
        </div>
      </div>
    </footer>
  );
};
