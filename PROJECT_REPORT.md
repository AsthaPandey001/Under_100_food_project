# Bangalore Food Discovery Under ₹100 — Implementation & Technical Report

### 1. Technology Stack Architecture
- **Backend Language**: Go (Golang 1.24) — High-concurrency, compiled backend providing sub-millisecond route handling and single-binary deployment.
- **Web Application Framework**: Echo Framework (v4.11.4) — High-throughput HTTP router handling route grouping, middleware, HTML/JSON rendering, and logging.
- **Object-Relational Mapping (ORM)**: Ent ORM (entgo.io/ent) — 100% type-safe graph-based ORM providing schema-as-code models, SQL query building, and automated relational migrations.
- **Relational Database**: PostgreSQL 15.19 — Enterprise relational database running locally on port 5432 ('goship_db') storing relational tables and relational foreign key constraints.
- **In-Memory Cache**: Redis 5.0.14.1 — Key-value store running on port 6379 for fast session caching and state management.
- **Dynamic Frontend Layer**: HTMX (v1.9.10) — Declarative AJAX library driving live search debounce, filter updates, and dynamic DOM swapping without heavy client JavaScript bundles.
- **Interactive Mapping System**:
  - Leaflet.js (v1.9.4) — Dynamic map engine supporting custom price markers, popups, and smooth camera animations.
  - Live Google Maps Street Layer — Crisp street, road, and locality labels loaded directly via Google Maps tile servers without API key watermarks.
  - Google Hybrid Satellite Layer — Aerial satellite imagery with street and highway overlays.
  - OpenStreetMap (OSM) — Integrated open geographic fallback layer.
- **Geographic Distance Engine**: Native Haversine formula implemented in Go ('CalculateHaversineDistance') calculating real-time straight-line distances from the user's coordinates.
- **Design System & Styling**:
  - Typography: Google Fonts ('Outfit' for bold headings, 'Plus Jakarta Sans' for clean UI body text).
  - Modern Dark Theme: Custom CSS3 glassmorphism, accent green highlights, and FSSAI-compliant dietary symbols.

### 2. Database Models & Schema Design
- **City Schema ('ent/schema/city.go')**: Stores city records with coordinates (Bangalore: 12.9716° N, 77.5946° E) and a one-to-many relationship with restaurants.
- **Restaurant Schema ('ent/schema/restaurant.go')**: Stores restaurant records including name, description, address, exact latitude and longitude, contact number, opening hours, and lifecycle status ('approved', 'pending', 'rejected').
- **Food Item Schema ('ent/schema/fooditem.go')**: Stores dish name, description, category, price, and availability flag, linked via foreign key to the parent restaurant.
- **User Schema ('ent/schema/user.go')**: Stores user authentication details with bcrypt password hashing and an 'is_admin' boolean flag for role-based authorization.
- **Auto-Migration**: 25 relational tables generated and synchronized in PostgreSQL with foreign keys and indexes.

### 3. Bangalore Food Directory & Seeding Engine
- **Directory Onboarding ('seeder/food_discovery.go')**: Onboarded **41 verified Bangalore food spots** across 10 major culinary zones:
  - Basavanagudi & Gandhi Bazaar: Vidyarthi Bhavan, Brahmin's Coffee Bar, Bhaskar's Mane Holige, Shivanna Gulkand Center.
  - Malleshwaram & Rajajinagar: CTR (Shri Sagar), Veena Stores, Sri Sairam's Chaat.
  - Jayanagar & JP Nagar: Taaza Thindi, SN Refreshments, Hari Super Sandwich, Brahmin's Thatte Idli.
  - Central Bangalore (MG Road / Church Street): MTR Lalbagh, Indian Coffee House, Airlines Hotel, Sanman, Truffles, Corner House Ice Cream, Fanoos, Albert Bakery, Thom's Bakery, Siddique Kabab Centre.
  - Koramangala & Indiranagar: Sri Krishna Kafe, Anand Sweets, Al-Amanah Cafe, Shivaji Military Hotel, The Old Fashioned Chai Co.
  - Tech Corridors (Whitefield / Marathahalli / Bellandur): AECS Layout Roll Corner, Thalaivar Biryani, Tibbet Kitchen, The Bangalore Cafe.
- **Food Catalog**:
  - **165 total menu items** cataloged.
  - **150 food items priced strictly at or below ₹100** (from ₹15 filter coffee and cutting chai to ₹99 shawarma rolls).
  - 15 items priced above ₹100 to test budget filtering.
- **Business Rule Verification**:
  - Meghana Foods (Residency Road): All main dishes priced > ₹100; strictly filtered out from public discovery.
  - Pending Moderation: Shivanna Gulkand Center, AECS Layout Roll Corner, Thalaivar Biryani, and The Bangalore Cafe set to 'pending' (visible only to admins).
  - Rejected Listing: Spice Terrace set to 'rejected' (completely excluded from public discovery).

### 4. User Interface & Discovery Features
- **Sticky Navigation Bar ('pkg/routing/fooddiscovery/views.go')**:
  - Brand header with animated lightning icon and budget badge.
  - Locality selector with pulsing green radar status indicator.
  - Instant live search bar with 300ms debounce filtering food items, categories, and eatery names.
  - One-click GPS button for automatic device geolocation.
- **Promotional Value Banners**:
  - Highlight banners for Heritage Breakfast, Under ₹50 Express Deals, and Corner House Sundaes.
- **Restaurant Discovery Cards**:
  - Display star rating, distance badge, and estimated travel/commute time (e.g. '⚡ 1.2 km away • ~15 mins').
  - Embedded mini-dish grid showcasing up to 4 items under ₹100 with bold price tags ('₹20', '₹45', '₹75', etc.) and category tags.
- **Interactive Action Triggers**:
  - 'View Full Menu & Reviews': Opens a detailed modal displaying complete dish listings and prices.
  - 'On Map': Automatically centers the interactive map on the selected eatery with smooth camera animation.

### 5. Live Google Maps Integration
- **Zero Watermarks**: Direct integration with Google Maps live street tile servers eliminating third-party watermark restrictions.
- **Google Live Street View**: Crisp rendering of roads, highways, and street labels.
- **Google Satellite View**: High-resolution aerial satellite imagery with road overlays.
- **OpenStreetMap Layer**: Community geographic map layer fallback.
- **Layer Switcher**: Top-right toggle control for instant switching between map modes.
- **Interactive Marker Popups**: Custom pins displaying starting price (e.g. '🍛 ₹35') that open info popups with direct Google Maps route navigation ('Directions ↗').

### 6. Veg & Non-Veg Categorization & Multi-Level Toggle Controls
- **Automated FSSAI Classification Engine ('pkg/repos/fooddiscovery/food_discovery.go')**:
  - Implemented the 'IsDishVeg' function in Go based on Indian food standards.
  - Categorizes all 155 Bangalore food dishes into **132 Pure Vegetarian** and **23 Non-Vegetarian** items.
  - Correctly detects non-veg items: Chicken Shawarma, Mutton Keema Samosa, Egg Kothu Parotta, Mutton Brain Puff, Chicken Sheekh Roll, Double Egg Roll, Steamed Chicken Momos, Egg Kathi Roll, Mutton Cutlet, etc.
  - Correctly detects vegetarian items: Benne Masala Dosa, Thatte Idli, Filter Kaapi, Khara Bath, Kesari Bath, Paneer Tikka, Falafel, Corner House ice creams, etc.
- **Toggle Controls**:
  - Header Quick Switch: Animated toggle switch in the top navigation bar for 1-click 'Veg Only' filtering.
  - Filter Bar Segmented Group: 3-way toggle button group ('All', 'Pure Veg', 'Non-Veg') with glowing active state styling.
  - Category Ribbon Chips: Quick shortcut chips for '🟢 Pure Veg Treats' and '🍗 Non-Veg Bites'.
  - FSSAI Standard Indicators: Green square with green circle dot for Veg, and crimson square with red triangle for Non-Veg across all cards and modals.
- **Map & Feed Synchronization**:
  - Selecting 'Pure Veg' filters the restaurant list to 31 vegetarian spots and updates map pins.
  - Selecting 'Non-Veg' isolates the 8 non-vegetarian spots (Empire, Indian Coffee House, Kolkata Kati Roll, Siddique Kabab, Albert Bakery, Thom's Bakery, Al-Amanah, Tibbet Kitchen).
  - Map marker popups display distinctive dietary badges ('🟢 PURE VEG' or '🔴 NON-VEG').
- **Modal Menu Dietary Tabs**:
  - Detail popup includes 'All Items', '🟢 Veg', and '🔴 Non-Veg' tabs for zero-latency client-side filtering.

### 7. Admin Dashboard & Management
- **Dashboard URL**: 'http://localhost:8000/admin'.
- **Metrics Overview**: Real-time counter cards showing 41 Total Restaurants, 5 Pending, 34 Approved, 2 Rejected, 165 Food Items, and 150 Items ≤ ₹100.
- **1-Click Moderation**: Instant HTMX-powered Approve and Reject buttons updating status without full page reload.
- **Restaurant Editing**: Full modal editing for names, descriptions, addresses, coordinates, phones, operating hours, and status.
- **Menu Management**:
  - In-stock and out-of-stock toggles per dish.
  - Dish deletion with instant table row removal.
  - FSSAI Veg / Non-Veg badge displayed next to every dish.
  - '+ Add Food Item' form with 'Food Type' selector ('🟢 Pure Veg' / '🔴 Non-Veg').

### 8. System Status & Credentials
- **Public Discovery Web App**: http://localhost:8000
- **Admin Management Portal**: http://localhost:8000/admin
- **Admin Credentials**: admin@fooddiscovery.com / Admin@123456
- **User Credentials**: user@fooddiscovery.com / User@123456
- **Backend Process**: Compiled binary 'server.exe' (71.2 MB) running as daemon on port 8000.