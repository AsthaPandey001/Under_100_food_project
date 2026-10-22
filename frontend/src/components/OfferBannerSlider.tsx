import React, { useState, useEffect, useRef } from "react";
import { ChevronLeft, ChevronRight, Tag, Sparkles, ArrowRight } from "lucide-react";

export interface OfferSlide {
  id: string;
  badge: string;
  badgeColor: string;
  badgeBg: string;
  title: string;
  highlightText: string;
  subtitle: string;
  code?: string;
  ctaText: string;
  ctaLink: string;
  image: string;
  bgColor: string;
  borderColor: string;
}

export const OFFER_SLIDES: OfferSlide[] = [
  {
    id: "offer-1",
    badge: "WELCOME DEAL",
    badgeColor: "#FF5A1F",
    badgeBg: "rgba(255, 90, 31, 0.12)",
    title: "Get Extra 20% Off On Your First",
    highlightText: "₹100 Meal Run",
    subtitle: "Use code BENGALURU20 at checkout • Valid on all tiffins, dosas, and meals.",
    code: "BENGALURU20",
    ctaText: "Order Now",
    ctaLink: "#food-discovery",
    image: "/crispy_dosa.jpg",
    bgColor: "#FFF6F0",
    borderColor: "#FED7AA"
  },
  {
    id: "offer-2",
    badge: "COMING SOON",
    badgeColor: "#087443",
    badgeBg: "rgba(8, 116, 67, 0.12)",
    title: "Under100 Daily Office Tiffin",
    highlightText: "Monthly Subscriptions",
    subtitle: "Fresh homestyle meals delivered directly to your tech park desk for just ₹79/day.",
    ctaText: "Join Early Waitlist",
    ctaLink: "#food-discovery",
    image: "/veg_pulao.jpg",
    bgColor: "#F0FDF4",
    borderColor: "#BBF7D0"
  },
  {
    id: "offer-3",
    badge: "SAVINGS PASS",
    badgeColor: "#0284C7",
    badgeBg: "rgba(2, 132, 199, 0.12)",
    title: "Zero Delivery Fee Club for",
    highlightText: "Students & Techies",
    subtitle: "Enjoy flat ₹0 delivery fees on all orders across Bengaluru • No minimum surge fees.",
    code: "PASS99",
    ctaText: "Explore Pass",
    ctaLink: "#how-it-works",
    image: "/samosa_plate.jpg",
    bgColor: "#F0F9FF",
    borderColor: "#BAE6FD"
  },
  {
    id: "offer-4",
    badge: "FUTURE ROLLOUT",
    badgeColor: "#7C3AED",
    badgeBg: "rgba(124, 58, 237, 0.12)",
    title: "Midnight Hunger Runs:",
    highlightText: "11 PM to 3 AM Delivery",
    subtitle: "Hot paneer rolls, steamed momos, and piping filter coffee delivered late night.",
    ctaText: "View Late Night Menu",
    ctaLink: "#carousel-quick-bites",
    image: "/kathi_roll.jpg",
    bgColor: "#FAF5FF",
    borderColor: "#E9D5FF"
  },
  {
    id: "offer-5",
    badge: "WEEKEND SPECIAL",
    badgeColor: "#D97706",
    badgeBg: "rgba(217, 119, 6, 0.12)",
    title: "Saturday & Sunday Under ₹99",
    highlightText: "Handi Biryani Bowls",
    subtitle: "Slow-cooked fragrant veg pulao and biryani bowls served with fresh cucumber raita.",
    ctaText: "Browse Weekend Specials",
    ctaLink: "#carousel-under100",
    image: "/hero_food.jpg",
    bgColor: "#FFFBEB",
    borderColor: "#FDE68A"
  }
];

export const OfferBannerSlider: React.FC = () => {
  const [currentSlide, setCurrentSlide] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const timeoutRef = useRef<number | null>(null);

  const totalSlides = OFFER_SLIDES.length;

  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };

  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };

  useEffect(() => {
    if (!isPaused) {
      timeoutRef.current = window.setTimeout(nextSlide, 5000);
    }
    return () => {
      if (timeoutRef.current !== null) {
        window.clearTimeout(timeoutRef.current);
      }
    };
  }, [currentSlide, isPaused]);

  const slide = OFFER_SLIDES[currentSlide];

  return (
    <section
      style={{
        padding: "20px 0 10px",
        backgroundColor: "var(--color-main-bg)",
        position: "relative"
      }}
      onMouseEnter={() => setIsPaused(true)}
      onMouseLeave={() => setIsPaused(false)}
    >
      <div className="container">
        {/* Main Banner Card */}
        <div
          style={{
            position: "relative",
            borderRadius: "20px",
            backgroundColor: slide.bgColor,
            border: `1.5px solid ${slide.borderColor}`,
            padding: "26px 36px",
            overflow: "hidden",
            boxShadow: "0 6px 20px rgba(0, 0, 0, 0.04)",
            transition: "background-color 0.4s ease, border-color 0.4s ease",
            minHeight: "200px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            flexWrap: "wrap",
            gap: "24px"
          }}
        >
          {/* Left Text & Offer Content */}
          <div style={{ maxWidth: "620px", zIndex: 2 }}>
            {/* Pill Badge */}
            <div style={{ display: "flex", alignItems: "center", gap: "10px", marginBottom: "12px" }}>
              <span
                style={{
                  display: "inline-flex",
                  alignItems: "center",
                  gap: "6px",
                  padding: "4px 12px",
                  borderRadius: "9999px",
                  backgroundColor: slide.badgeBg,
                  color: slide.badgeColor,
                  fontSize: "11px",
                  fontWeight: 800,
                  letterSpacing: "0.5px",
                  textTransform: "uppercase"
                }}
              >
                <Sparkles size={13} />
                <span>{slide.badge}</span>
              </span>

              {slide.code && (
                <span
                  style={{
                    display: "inline-flex",
                    alignItems: "center",
                    gap: "4px",
                    padding: "3px 10px",
                    borderRadius: "6px",
                    backgroundColor: "#ffffff",
                    border: `1px dashed ${slide.badgeColor}`,
                    color: slide.badgeColor,
                    fontSize: "12px",
                    fontWeight: 700
                  }}
                >
                  <Tag size={12} />
                  <span>Use: {slide.code}</span>
                </span>
              )}
            </div>

            {/* Headline */}
            <h3
              style={{
                fontSize: "clamp(22px, 3.2vw, 30px)",
                fontWeight: 900,
                color: "#171C24",
                lineHeight: 1.25,
                letterSpacing: "-0.5px",
                marginBottom: "8px"
              }}
            >
              {slide.title}{" "}
              <span style={{ color: slide.badgeColor }}>{slide.highlightText}</span>
            </h3>

            {/* Subtitle */}
            <p
              style={{
                fontSize: "14.5px",
                color: "#64748B",
                lineHeight: 1.5,
                marginBottom: "20px",
                maxWidth: "520px"
              }}
            >
              {slide.subtitle}
            </p>

            {/* CTA Button */}
            <a
              href={slide.ctaLink}
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "8px",
                padding: "10px 22px",
                borderRadius: "9999px",
                backgroundColor: slide.badgeColor,
                color: "#ffffff",
                fontSize: "13.5px",
                fontWeight: 700,
                textDecoration: "none",
                boxShadow: `0 4px 14px ${slide.badgeColor}33`,
                transition: "all 0.2s ease"
              }}
              onMouseEnter={(e) => (e.currentTarget.style.transform = "translateY(-1px)")}
              onMouseLeave={(e) => (e.currentTarget.style.transform = "none")}
            >
              <span>{slide.ctaText}</span>
              <ArrowRight size={15} />
            </a>
          </div>

          {/* Right Visual Image */}
          <div
            style={{
              position: "relative",
              width: "220px",
              height: "165px",
              borderRadius: "16px",
              overflow: "hidden",
              boxShadow: "0 8px 24px rgba(0,0,0,0.1)",
              border: "3px solid #ffffff",
              flexShrink: 0
            }}
          >
            <img
              src={slide.image}
              alt={slide.title}
              style={{
                width: "100%",
                height: "100%",
                objectFit: "cover",
                transition: "transform 0.4s ease"
              }}
            />
          </div>

          {/* Desktop Left/Right Controls on the Banner */}
          <button
            type="button"
            onClick={prevSlide}
            style={{
              position: "absolute",
              left: "10px",
              top: "50%",
              transform: "translateY(-50%)",
              width: "36px",
              height: "36px",
              borderRadius: "50%",
              backgroundColor: "rgba(255, 255, 255, 0.9)",
              border: "1px solid rgba(0, 0, 0, 0.08)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "#171C24",
              cursor: "pointer",
              boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
              zIndex: 3
            }}
            aria-label="Previous offer slide"
          >
            <ChevronLeft size={18} />
          </button>

          <button
            type="button"
            onClick={nextSlide}
            style={{
              position: "absolute",
              right: "10px",
              top: "50%",
              transform: "translateY(-50%)",
              width: "36px",
              height: "36px",
              borderRadius: "50%",
              backgroundColor: "rgba(255, 255, 255, 0.9)",
              border: "1px solid rgba(0, 0, 0, 0.08)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "#171C24",
              cursor: "pointer",
              boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
              zIndex: 3
            }}
            aria-label="Next offer slide"
          >
            <ChevronRight size={18} />
          </button>
        </div>

        {/* Slide Indicator Dots */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: "8px",
            marginTop: "12px"
          }}
        >
          {OFFER_SLIDES.map((_, idx) => (
            <button
              key={idx}
              type="button"
              onClick={() => setCurrentSlide(idx)}
              style={{
                width: currentSlide === idx ? "24px" : "8px",
                height: "8px",
                borderRadius: "4px",
                backgroundColor: currentSlide === idx ? "#FF5A1F" : "#CBD5E1",
                border: "none",
                cursor: "pointer",
                padding: 0,
                transition: "all 0.3s ease"
              }}
              aria-label={`Go to slide ${idx + 1}`}
            />
          ))}
        </div>
      </div>
    </section>
  );
};
