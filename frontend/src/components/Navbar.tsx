import React, { useState } from "react";
import { Smartphone, ShoppingBag, Menu, X } from "lucide-react";
import { LocationSelector } from "./LocationSelector";
import { SearchBar } from "./SearchBar";

interface NavbarProps {
  selectedNeighbourhood: string;
  onSelectNeighbourhood: (area: string) => void;
  searchQuery: string;
  onSearchChange: (q: string) => void;
  onClearSearch: () => void;
  cartCount: number;
  onOpenCart: () => void;
}

export const Navbar: React.FC<NavbarProps> = ({
  selectedNeighbourhood,
  onSelectNeighbourhood,
  searchQuery,
  onSearchChange,
  onClearSearch,
  cartCount,
  onOpenCart
}) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [appModalOpen, setAppModalOpen] = useState(false);

  return (
    <>
      <header style={{
        position: "sticky",
        top: 0,
        backgroundColor: "#ffffff",
        borderBottom: "1px solid #ECEAE5",
        boxShadow: "0 2px 10px rgba(0, 0, 0, 0.03)",
        zIndex: 90
      }}>
        <div className="container" style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "10px 24px",
          gap: "16px"
        }}>
          {/* Left: Brand Logo + Location Selector */}
          <div style={{ display: "flex", alignItems: "center", gap: "22px" }}>
            {/* Steaming Bowl Brand Logo as in Screenshot */}
            <a href="#" style={{ display: "flex", alignItems: "center", gap: "9px", textDecoration: "none" }}>
              {/* Steaming Orange Bowl Icon */}
              <div style={{ display: "flex", alignItems: "center", justifyContent: "center" }}>
                <svg width="34" height="34" viewBox="0 0 36 36" fill="none" xmlns="http://www.w3.org/2000/svg">
                  {/* 3 Wavy Steam Lines */}
                  <path d="M12 4C11 6.5 13 8.5 12 11" stroke="#FF5A1F" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M18 3C17 5.5 19 8 18 10.5" stroke="#FF5A1F" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M24 4C23 6.5 25 8.5 24 11" stroke="#FF5A1F" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round"/>
                  {/* Bowl Rim */}
                  <path d="M6 14H30" stroke="#FF5A1F" strokeWidth="2.5" strokeLinecap="round"/>
                  {/* Bowl Body */}
                  <path d="M7 16C7.5 24 13.5 28 18 28C22.5 28 28.5 24 29 16H7Z" fill="#FF5A1F"/>
                  {/* Bowl Base */}
                  <path d="M13 28.5H23" stroke="#FF5A1F" strokeWidth="2.2" strokeLinecap="round"/>
                </svg>
              </div>

              {/* under100 Wordmark */}
              <div style={{ display: "flex", alignItems: "baseline", letterSpacing: "-0.6px" }}>
                <span style={{ fontSize: "24px", fontWeight: 800, color: "#111827" }}>
                  under
                </span>
                <span style={{ fontSize: "24px", fontWeight: 900, color: "#FF5A1F" }}>
                  100
                </span>
              </div>
            </a>

            {/* Location Selector Pill */}
            <div className="desktop-only" style={{ display: "flex" }}>
              <LocationSelector
                selected={selectedNeighbourhood}
                onSelect={onSelectNeighbourhood}
              />
            </div>
          </div>

          {/* Centre: Search Bar */}
          <div className="desktop-search" style={{ flex: 1, display: "flex", justifyContent: "center", maxWidth: "480px" }}>
            <SearchBar
              value={searchQuery}
              onChange={onSearchChange}
              onClear={onClearSearch}
            />
          </div>

          {/* Right Navigation & CTAs */}
          <div style={{ display: "flex", alignItems: "center", gap: "20px" }}>
            {/* 'How it works' Link */}
            <a
              href="#how-it-works"
              className="desktop-only"
              style={{
                fontSize: "14px",
                fontWeight: 600,
                color: "#1E293B",
                textDecoration: "none",
                whiteSpace: "nowrap"
              }}
            >
              How it works
            </a>

            {/* 'Sign in' Link */}
            <button
              type="button"
              onClick={() => alert("Sign in modal: Demo feature for registered users.")}
              className="desktop-only"
              style={{
                fontSize: "14px",
                fontWeight: 600,
                color: "#1E293B",
                cursor: "pointer",
                whiteSpace: "nowrap"
              }}
            >
              Sign in
            </button>

            {/* Cart Icon / Counter */}
            <button
              type="button"
              onClick={onOpenCart}
              style={{
                position: "relative",
                display: "flex",
                alignItems: "center",
                gap: "6px",
                backgroundColor: cartCount > 0 ? "#FFF0E6" : "#F4F5F7",
                border: "1px solid",
                borderColor: cartCount > 0 ? "#FF5A1F" : "#E2E8F0",
                padding: "8px 14px",
                borderRadius: "9999px",
                fontWeight: 600,
                fontSize: "13px",
                color: cartCount > 0 ? "#FF5A1F" : "#1E293B",
                cursor: "pointer",
                transition: "all 0.2s ease"
              }}
              aria-label="View shopping cart"
            >
              <ShoppingBag size={17} />
              <span className="cart-text">Cart</span>
              {cartCount > 0 && (
                <span style={{
                  backgroundColor: "#FF5A1F",
                  color: "#ffffff",
                  fontSize: "11px",
                  fontWeight: 800,
                  width: "19px",
                  height: "19px",
                  borderRadius: "50%",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center"
                }}>
                  {cartCount}
                </span>
              )}
            </button>

            {/* Deep Green 'Get the app' Button */}
            <button
              type="button"
              onClick={() => setAppModalOpen(true)}
              className="desktop-only"
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                backgroundColor: "#0B6837",
                color: "#ffffff",
                padding: "10px 20px",
                borderRadius: "9999px",
                fontSize: "14px",
                fontWeight: 700,
                boxShadow: "0 2px 8px rgba(11, 104, 55, 0.25)",
                cursor: "pointer",
                transition: "all 0.2s ease",
                whiteSpace: "nowrap"
              }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#09532C")}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#0B6837")}
            >
              <Smartphone size={17} strokeWidth={2.3} />
              <span>Get the app</span>
            </button>

            {/* Mobile Menu Toggle */}
            <button
              type="button"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="mobile-only"
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                padding: "8px",
                color: "#1E293B",
                cursor: "pointer"
              }}
              aria-label="Toggle navigation menu"
            >
              {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>

        {/* Mobile Search & Location Bar */}
        <div className="mobile-search-bar" style={{
          padding: "10px 20px 14px",
          display: "flex",
          flexDirection: "column",
          gap: "10px",
          borderTop: "1px solid #ECEAE5"
        }}>
          <SearchBar
            value={searchQuery}
            onChange={onSearchChange}
            onClear={onClearSearch}
          />
          <div style={{ display: "flex", justifyContent: "flex-start" }}>
            <LocationSelector
              selected={selectedNeighbourhood}
              onSelect={onSelectNeighbourhood}
            />
          </div>
        </div>

        {/* Mobile Slide-down Menu */}
        {mobileMenuOpen && (
          <div style={{
            backgroundColor: "#ffffff",
            borderBottom: "1px solid #ECEAE5",
            padding: "16px 20px",
            display: "flex",
            flexDirection: "column",
            gap: "12px"
          }}>
            <a
              href="#carousels-area"
              onClick={() => setMobileMenuOpen(false)}
              style={{ fontSize: "15px", fontWeight: 600, color: "#1E293B", padding: "6px 0" }}
            >
              Explore Cravings
            </a>
            <a
              href="#food-discovery"
              onClick={() => setMobileMenuOpen(false)}
              style={{ fontSize: "15px", fontWeight: 600, color: "#1E293B", padding: "6px 0" }}
            >
              Meals Under ₹100
            </a>
            <a
              href="#how-it-works"
              onClick={() => setMobileMenuOpen(false)}
              style={{ fontSize: "15px", fontWeight: 600, color: "#1E293B", padding: "6px 0" }}
            >
              How It Works
            </a>
            <button
              type="button"
              onClick={() => { setMobileMenuOpen(false); setAppModalOpen(true); }}
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                gap: "8px",
                backgroundColor: "#0B6837",
                color: "#ffffff",
                padding: "12px",
                borderRadius: "12px",
                fontSize: "14px",
                fontWeight: 700,
                marginTop: "8px",
                cursor: "pointer"
              }}
            >
              <Smartphone size={17} />
              <span>Get the Under100 App</span>
            </button>
          </div>
        )}
      </header>

      {/* App Download Modal */}
      {appModalOpen && (
        <div style={{
          position: "fixed",
          inset: 0,
          backgroundColor: "rgba(0, 0, 0, 0.5)",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          zIndex: 200,
          padding: "20px"
        }}>
          <div style={{
            backgroundColor: "#ffffff",
            borderRadius: "20px",
            padding: "28px",
            maxWidth: "420px",
            width: "100%",
            boxShadow: "0 20px 40px rgba(0,0,0,0.2)",
            position: "relative"
          }}>
            <button
              type="button"
              onClick={() => setAppModalOpen(false)}
              style={{ position: "absolute", top: "16px", right: "16px", color: "#64748B", cursor: "pointer" }}
              aria-label="Close modal"
            >
              <X size={20} />
            </button>
            <div style={{ display: "inline-flex", padding: "10px", borderRadius: "12px", backgroundColor: "#EBF7EE", color: "#0B6837", marginBottom: "14px" }}>
              <Smartphone size={24} />
            </div>
            <h3 style={{ fontSize: "20px", fontWeight: 800, marginBottom: "8px", color: "#1E293B" }}>
              Get Under100 on Mobile
            </h3>
            <p style={{ fontSize: "14px", color: "#64748B", lineHeight: 1.5, marginBottom: "20px" }}>
              Discover fast, satisfying meals under ₹100 across Bengaluru with instant delivery and zero hidden charges.
            </p>
            <div style={{
              padding: "16px",
              backgroundColor: "#F8F6F0",
              borderRadius: "14px",
              border: "1px dashed #CBD5E1",
              textAlign: "center",
              marginBottom: "20px"
            }}>
              <div style={{ fontWeight: 700, fontSize: "13px", color: "#0B6837", marginBottom: "4px" }}>
                Active in Bengaluru Hubs
              </div>
              <div style={{ fontSize: "12px", color: "#64748B" }}>
                Indiranagar • Koramangala • HSR Layout • Whitefield • Malleshwaram
              </div>
            </div>
            <button
              type="button"
              onClick={() => setAppModalOpen(false)}
              style={{
                width: "100%",
                padding: "12px",
                borderRadius: "14px",
                backgroundColor: "#0B6837",
                color: "#ffffff",
                fontWeight: 700,
                fontSize: "14px",
                cursor: "pointer"
              }}
            >
              Continue Browsing
            </button>
          </div>
        </div>
      )}

      <style>{`
        @media (min-width: 901px) {
          .mobile-only { display: none !important; }
          .mobile-search-bar { display: none !important; }
        }
        @media (max-width: 900px) {
          .desktop-only { display: none !important; }
          .desktop-search { display: none !important; }
        }
      `}</style>
    </>
  );
};
