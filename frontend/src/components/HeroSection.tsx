import React from "react";
import { ArrowRight, CheckCircle2, ShieldCheck, Sparkles } from "lucide-react";

export const HeroSection: React.FC = () => {
  return (
    <section style={{
      backgroundColor: "var(--color-card-bg)",
      padding: "56px 0 64px",
      borderBottom: "1px solid var(--color-border)",
      position: "relative",
      overflow: "hidden"
    }}>
      <div className="container" style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(320px, 1fr))",
        gap: "48px",
        alignItems: "center"
      }}>
        {/* Left Column: Brand Statement */}
        <div style={{ maxWidth: "560px" }}>
          {/* Small Orange Pill */}
          <div style={{
            display: "inline-flex",
            alignItems: "center",
            gap: "6px",
            backgroundColor: "var(--color-soft-orange-bg)",
            color: "var(--color-primary-orange)",
            border: "1px solid rgba(255, 90, 31, 0.2)",
            padding: "6px 14px",
            borderRadius: "var(--radius-full)",
            fontSize: "12px",
            fontWeight: 800,
            letterSpacing: "0.5px",
            textTransform: "uppercase",
            marginBottom: "20px"
          }}>
            <Sparkles size={14} />
            <span>GOOD FOOD. SMALL BUDGET.</span>
          </div>

          {/* Main Headline */}
          <h1 style={{
            fontSize: "clamp(36px, 5vw, 54px)",
            fontWeight: 900,
            lineHeight: 1.1,
            letterSpacing: "-1.5px",
            color: "var(--color-dark-text)",
            marginBottom: "18px"
          }}>
            Big cravings. <br />
            <span style={{ color: "var(--color-primary-orange)" }}>Tiny prices.</span>
          </h1>

          {/* Supporting Copy */}
          <p style={{
            fontSize: "17px",
            color: "var(--color-secondary-text)",
            lineHeight: 1.6,
            marginBottom: "28px",
            maxWidth: "480px"
          }}>
            Discover satisfying local favourites across Bengaluru — every pick under ₹100.
          </p>

          {/* Actions */}
          <div style={{
            display: "flex",
            alignItems: "center",
            flexWrap: "wrap",
            gap: "16px",
            marginBottom: "24px"
          }}>
            <a
              href="#food-discovery"
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "8px",
                backgroundColor: "var(--color-deep-green)",
                color: "#ffffff",
                padding: "14px 28px",
                borderRadius: "var(--radius-full)",
                fontSize: "15px",
                fontWeight: 700,
                boxShadow: "0 4px 14px rgba(8, 116, 67, 0.25)",
                transition: "all 0.2s ease"
              }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "var(--color-deep-green-hover)")}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "var(--color-deep-green)")}
            >
              <span>Explore meals under ₹100</span>
              <ArrowRight size={18} />
            </a>

            <a
              href="#how-it-works"
              style={{
                fontSize: "14px",
                fontWeight: 700,
                color: "var(--color-dark-text)",
                padding: "10px 16px",
                borderRadius: "var(--radius-sm)",
                display: "inline-flex",
                alignItems: "center",
                gap: "4px"
              }}
            >
              See how it works
            </a>
          </div>

          {/* Small Trust Message */}
          <div style={{
            display: "flex",
            alignItems: "center",
            gap: "8px",
            fontSize: "13px",
            color: "var(--color-secondary-text)",
            fontWeight: 500
          }}>
            <ShieldCheck size={16} color="var(--color-deep-green)" />
            <span>No fancy fees. Just food that fits.</span>
          </div>
        </div>

        {/* Right Column: Hero Food Spread Photography */}
        <div style={{ position: "relative" }}>
          <div style={{
            borderRadius: "var(--radius-lg)",
            overflow: "hidden",
            boxShadow: "0 12px 36px rgba(23, 28, 36, 0.12)",
            border: "1px solid var(--color-border)",
            backgroundColor: "#F3F1EC",
            position: "relative",
            aspectRatio: "4 / 3",
            maxHeight: "440px"
          }}>
            <img
              src="/hero_food.jpg"
              alt="Authentic Bengaluru South Indian meal spread with Masala Dosa, Idli, Sambar, and Chutneys"
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

          {/* Circular Orange Price Badge */}
          <div style={{
            position: "absolute",
            top: "-16px",
            right: "-16px",
            width: "108px",
            height: "108px",
            borderRadius: "50%",
            backgroundColor: "var(--color-primary-orange)",
            color: "#ffffff",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            textAlign: "center",
            boxShadow: "0 8px 24px rgba(255, 90, 31, 0.35)",
            border: "3px solid #ffffff",
            transform: "rotate(6deg)"
          }}>
            <span style={{ fontSize: "10px", fontWeight: 800, letterSpacing: "0.5px", textTransform: "uppercase" }}>
              MEALS FROM
            </span>
            <span style={{ fontSize: "24px", fontWeight: 900, lineHeight: 1 }}>
              ₹49
            </span>
          </div>

          {/* Small Supporting Label Card */}
          <div style={{
            position: "absolute",
            bottom: "20px",
            left: "20px",
            backgroundColor: "rgba(255, 255, 255, 0.95)",
            backdropFilter: "blur(8px)",
            padding: "8px 14px",
            borderRadius: "var(--radius-full)",
            display: "flex",
            alignItems: "center",
            gap: "6px",
            boxShadow: "0 4px 16px rgba(0,0,0,0.1)",
            border: "1px solid var(--color-border)"
          }}>
            <CheckCircle2 size={16} color="var(--color-deep-green)" />
            <span style={{ fontSize: "12px", fontWeight: 700, color: "var(--color-dark-text)" }}>
              Your budget, respected ✓
            </span>
          </div>
        </div>
      </div>
    </section>
  );
};
