import React, { useMemo } from "react";
import { AnnouncementBar } from "../components/AnnouncementBar";
import { Navbar } from "../components/Navbar";
import { OfferBannerSlider } from "../components/OfferBannerSlider";
import { FoodCategoryCarousel } from "../components/FoodCategoryCarousel";
import { FoodProductCarousel } from "../components/FoodProductCarousel";
import { HeroSection } from "../components/HeroSection";
import { BenefitsStrip } from "../components/BenefitsStrip";
import { FoodGrid } from "../components/FoodGrid";
import { PromoBanner } from "../components/PromoBanner";
import { HowItWorks } from "../components/HowItWorks";
import { WhyUnder100 } from "../components/WhyUnder100";
import { CartDrawer } from "../components/CartDrawer";
import { Footer } from "../components/Footer";
import { useCart } from "../hooks/useCart";
import { useFoodSearch } from "../hooks/useFoodSearch";
import { FOOD_ITEMS } from "../data/foodItems";
import { ArrowRight } from "lucide-react";

export const Home: React.FC = () => {
  const {
    cartItems,
    addToCart,
    updateQuantity,
    removeFromCart,
    clearCart,
    getItemQuantity,
    totalCount,
    subtotal,
    isDrawerOpen,
    setIsDrawerOpen
  } = useCart();

  const {
    searchQuery,
    setSearchQuery,
    selectedCategory,
    setSelectedCategory,
    selectedNeighbourhood,
    setSelectedNeighbourhood,
    sortBy,
    setSortBy,
    filteredItems,
    resetFilters
  } = useFoodSearch();

  // Curate items for Carousels B, C, D, E
  const under100Items = useMemo(() => {
    return FOOD_ITEMS.filter(f => f.tags.includes("under100"));
  }, []);

  const southIndianItems = useMemo(() => {
    return FOOD_ITEMS.filter(f => f.tags.includes("south-indian-fav"));
  }, []);

  const quickBiteItems = useMemo(() => {
    return FOOD_ITEMS.filter(f => f.tags.includes("quick-bite"));
  }, []);

  const bengaluruPickItems = useMemo(() => {
    return FOOD_ITEMS.filter(f => f.tags.includes("bengaluru-pick"));
  }, []);

  return (
    <div style={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      {/* 1. Announcement Strip */}
      <AnnouncementBar />

      {/* 2. Main Navigation */}
      <Navbar
        selectedNeighbourhood={selectedNeighbourhood}
        onSelectNeighbourhood={setSelectedNeighbourhood}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        onClearSearch={() => setSearchQuery("")}
        cartCount={totalCount}
        onOpenCart={() => setIsDrawerOpen(true)}
      />

      {/* 3. MULTI-OFFER BANNER SLIDER (ABOVE CAROUSELS) */}
      <OfferBannerSlider />

      {/* 4. FOOD CAROUSELS DIRECTLY BELOW (BEFORE HERO) */}
      <div id="carousels-area">
        {/* Carousel A: Popular cravings */}
        <FoodCategoryCarousel
          selectedCategory={selectedCategory}
          onSelectCategory={(catId: string) => {
            setSelectedCategory(catId);
            // Smoothly scroll down to discovery section
            const el = document.getElementById("food-discovery");
            if (el) el.scrollIntoView({ behavior: "smooth" });
          }}
        />

        {/* Carousel B: Meals under ₹100 */}
        <FoodProductCarousel
          id="carousel-under100"
          heading="Your ₹100 food run starts here"
          supportingText="Big flavours. Small prices."
          items={under100Items}
          getItemQuantity={getItemQuantity}
          onAdd={addToCart}
          onUpdateQuantity={updateQuantity}
        />

        {/* Carousel C: South Indian favourites */}
        <FoodProductCarousel
          id="carousel-south-indian"
          heading="South Indian, sorted"
          supportingText="Crisp dosas, fluffy idlis, and authentic Bengaluru filter coffee."
          items={southIndianItems}
          getItemQuantity={getItemQuantity}
          onAdd={addToCart}
          onUpdateQuantity={updateQuantity}
        />

        {/* Carousel D: Quick bites */}
        <FoodProductCarousel
          id="carousel-quick-bites"
          heading="Snack now. Spend less."
          supportingText="Hot samosas, crunchy rolls, momos, and street sandwiches."
          items={quickBiteItems}
          getItemQuantity={getItemQuantity}
          onAdd={addToCart}
          onUpdateQuantity={updateQuantity}
        />

        {/* Carousel E: Local Bengaluru picks */}
        <FoodProductCarousel
          id="carousel-bengaluru-picks"
          heading="Bengaluru favourites near you"
          supportingText="Neighbourhood darshini legends and beloved local eateries."
          items={bengaluruPickItems}
          getItemQuantity={getItemQuantity}
          onAdd={addToCart}
          onUpdateQuantity={updateQuantity}
        />
      </div>

      {/* 5. HERO SECTION — DIRECTLY AFTER THE FOOD CAROUSELS */}
      <HeroSection />

      {/* 6. TRUST AND VALUE STRIP */}
      <BenefitsStrip />

      {/* 7. MAIN FOOD DISCOVERY SECTION */}
      <FoodGrid
        items={filteredItems}
        selectedCategory={selectedCategory}
        onSelectCategory={setSelectedCategory}
        sortBy={sortBy}
        onSortChange={setSortBy}
        selectedNeighbourhood={selectedNeighbourhood}
        searchQuery={searchQuery}
        onResetFilters={resetFilters}
        getItemQuantity={getItemQuantity}
        onAdd={addToCart}
        onUpdateQuantity={updateQuantity}
      />

      {/* 8. AFFORDABILITY PROMOTIONAL BANNER */}
      <PromoBanner />

      {/* 9. HOW IT WORKS */}
      <HowItWorks />

      {/* 10. WHY UNDER100? */}
      <WhyUnder100 />

      {/* 11. FOOTER */}
      <Footer />

      {/* Slide-out Cart Drawer */}
      <CartDrawer
        isOpen={isDrawerOpen}
        onClose={() => setIsDrawerOpen(false)}
        cartItems={cartItems}
        subtotal={subtotal}
        totalCount={totalCount}
        onUpdateQuantity={updateQuantity}
        onRemoveItem={removeFromCart}
        onClearCart={clearCart}
      />

      {/* Floating Sticky Mobile Cart Bar */}
      {totalCount > 0 && !isDrawerOpen && (
        <div style={{
          position: "fixed",
          bottom: "16px",
          left: "16px",
          right: "16px",
          backgroundColor: "var(--color-dark-text)",
          color: "#ffffff",
          padding: "12px 18px",
          borderRadius: "var(--radius-full)",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          boxShadow: "var(--shadow-elevated)",
          zIndex: 80,
          cursor: "pointer"
        }}
        onClick={() => setIsDrawerOpen(true)}
        >
          <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
            <div style={{
              backgroundColor: "var(--color-primary-orange)",
              color: "#ffffff",
              padding: "4px 8px",
              borderRadius: "var(--radius-full)",
              fontSize: "12px",
              fontWeight: 800
            }}>
              {totalCount} {totalCount === 1 ? "item" : "items"}
            </div>
            <div style={{ fontSize: "14px", fontWeight: 700 }}>
              ₹{subtotal} <span style={{ fontSize: "12px", fontWeight: 400, color: "#94A3B8" }}>plus zero surprise fees</span>
            </div>
          </div>

          <div style={{
            display: "flex",
            alignItems: "center",
            gap: "6px",
            color: "var(--color-primary-orange)",
            fontWeight: 800,
            fontSize: "13px"
          }}>
            <span>VIEW CART</span>
            <ArrowRight size={16} />
          </div>
        </div>
      )}
    </div>
  );
};
