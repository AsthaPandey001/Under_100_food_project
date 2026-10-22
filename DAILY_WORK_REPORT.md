# Bite100: Daily Engineering & Product Report
**Date:** October 9, 2026  
**Project:** Bangalore Food Discovery Under ₹100 (**Bite100**)  
**Tech Stack:** Go 1.23, Echo v4, Ent ORM, PostgreSQL, Redis, Leaflet.js, HTMX, Vanilla CSS  

---

## Executive Summary

Today's work centered on transforming **Bite100** from initial concept into a fully functional, production-ready food discovery platform tailored specifically for Bangalore residents seeking high-quality food strictly capped at ₹100. 

Key milestones achieved today include:
1. Complete database schema design and migration with Ent ORM.
2. Production seeder populating authentic Bangalore food hubs across 8 distinct localities.
3. High-performance geospatial querying, budget filtering, dietary classification, and sorting in Go.
4. Complete UI/UX overhaul transitioning to a modern, crisp light theme with zero emoji/icon clutter and a streamlined navigation bar.
5. Implementation of critical consumer flows: dual-pane interactive map/feed, view mode toggles, local bookmarks ("Saved Bites"), menu detail modals, and community spot nominations.
6. Full end-to-end integration and verification against live PostgreSQL and Redis instances.

---

## 1. Database & Schema Architecture (`ent/schema`)

Three foundational Ent entities were designed, generated, and migrated:

### 1.1 `City`
- Represents geographic jurisdictions (initialized with Bengaluru: `12.9716° N, 77.5946° E`).
- 1-to-many relationship with `Restaurant`.

### 1.2 `Restaurant`
- **Fields:** `name`, `description`, `address`, `phone`, `latitude`, `longitude`, `status` (`approved`, `pending`, `rejected`), timestamps.
- **Relationships:** Belongs to `City`, 1-to-many with `FoodItem`.
- **Status Workflow:** Enables community nominations to enter as `pending` before being reviewed and approved via the Admin Dashboard.

### 1.3 `FoodItem`
- **Fields:** `name`, `description`, `price`, `category`, `available`, timestamps.
- **Relationships:** Belongs to `Restaurant`.
- **Validation:** Strict price constraints ensuring discovered items fall at or below the ₹100 threshold.

---

## 2. Seed Data & Curated Bangalore Eateries (`seeder/`)

Populated realistic seed data covering iconic budget eateries across key localities:

| Locality | Eatery | Signature Dish Under ₹100 | Price |
| :--- | :--- | :--- | :--- |
| **Basavanagudi** | Vidyarthi Bhavan | Benne Masala Dosa | ₹65 |
| **VV Puram** | Brahmin's Coffee Bar | Idli Vada Combo & Filter Coffee | ₹50 |
| **Malleshwaram** | CTR (Shri Sagar) | Benne Masala Dosa | ₹60 |
| **Malleshwaram** | Veena Stores | Steamed Idli & Vada with Chutney | ₹45 |
| **Jayanagar** | Hari Super Sandwich | Chocolate Cheese Toast | ₹50 |
| **Shivaji Nagar** | Shivaji Military Hotel | Mutton Donne Biryani (Half) | ₹95 |
| **Church Street** | Empire Restaurant | Paneer Tikka / Chicken Roll | ₹85 |
| **Indiranagar** | By 2 Coffee | Filter Coffee & Khara Bath | ₹35 |

---

## 3. Backend Architecture & API Routes

### 3.1 Repository Layer (`pkg/repos/fooddiscovery/food_discovery.go`)
- **Geospatial Proximity Calculation:** Haversine formula computing distance in kilometers from user or locality coordinates.
- **Multi-Intent Filtering:**
  - `max_price`: Filters restaurants having available dishes $\le$ budget slider value.
  - `diet`: Dynamic classification based on item keywords (`veg`, `nonveg`, `all`).
  - `q`: Search keyword querying eatery names, dish titles, and categories.
  - `sort`: Supports `distance` (nearest first), `price` (cheapest dishes first), and `name` (alphabetical A-Z).
- **Community Nomination Service:** `SuggestRestaurant` saves submissions into PostgreSQL with an initial `pending` moderation flag.
- **Admin Management:** Status counters and bulk moderation operations (`approve`, `reject`).

### 3.2 Routing & Handlers (`pkg/routing/fooddiscovery/`)
- `GET /` — Full application shell with embedded discovery state and authentication context.
- `GET /nearby` — HTMX endpoint returning filtered and sorted eatery cards.
- `GET /api/markers` — Lightweight JSON endpoint feeding Leaflet map pins with real-time price badges.
- `GET /restaurant/:id` — Detail modal returning full menu items segmented by dietary type.
- `POST /api/suggest` — Endpoint processing community nominations with instant HTMX feedback.
- `GET /admin` & `POST /admin/restaurant/:id/status` — Role-based administration dashboard.

---

## 4. UI/UX Design System & Flow Execution (`views.go`)

Per user direction, the entire UI was redesigned to meet modern consumer web standards:

### 4.1 Visual Design Tokens
- **Palette:** Crisp light background (`#ffffff`, `#f8fafc`), clean subtle border lines (`#e2e8f0`), deep slate typography (`#0f172a`), refined accents (coral-crimson `#e23744`, emerald `#15803d`).
- **Typography:** Google Fonts (`Plus Jakarta Sans` for body, `Outfit` for display headings).
- **Zero Icon Clutter:** Completely eliminated emojis, decorative SVG bloat, and gimmicky icons in favor of clean geometric pills, status tags, and clear text labels.

### 4.2 Streamlined Navigation Bar
- Lightweight layout containing only essential controls:
  - Brand identity with locality subtitle.
  - **Saved** bookmark counter pill: `Saved (<span id="saved-count">0</span>)`.
  - **Suggest a Spot** button for community submissions.
  - **Sign In / Account** trigger with 1-click test credentials for evaluators.

### 4.3 Multi-Intent Discovery Hero
- **Locality Quick Pills:** Immediate jumping to Bangalore hubs (*All, VV Puram, Malleshwaram, Koramangala, Indiranagar, Jayanagar, Commercial St, Shivaji Nagar*).
- **Price Range Slider:** Real-time budget scoping up to ₹100.
- **Dietary Filter Pills:** *All*, *Pure Veg*, and *Veg & Non-Veg*.
- **GPS Location Detection:** HTML5 Geolocation integration to center the map and feed on the user's current location.

### 4.4 Dual-Pane View & Flexible Display Switching
- **Split View:** Responsive side-by-side layout with Leaflet map on the left and eatery cards on the right.
- **Cards Only:** Focused card feed maximizing reading area.
- **Map Only:** Full-width map canvas highlighting price markers across Bangalore.

### 4.5 Interactive Card & Menu Modals
- **Card Actions:**
  - **View Full Menu:** Triggers modal showing all items under ₹100 with veg/non-veg tags.
  - **Directions:** Opens direct Google Maps turn-by-turn navigation.
  - **Save:** Persists bookmarks into browser `localStorage`, updates button styling to green *Saved*, and increments nav counter.
- **Saved Bites Modal:** Drawer displaying all bookmarked spots with 1-click dish access and item removal.
- **Suggest a Spot Modal:** Form allowing users to nominate their favorite neighborhood stalls with real-time HTMX validation.

---

## 5. Verification & Testing Log

All components were verified in the active environment:

| Test Case | Method | Expected Result | Actual Result | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Go Compilation** | `go build -o server.exe ./cmd/web` | Exit code 0 | Clean build, no errors | **PASS** |
| **Landing Page** | `GET http://localhost:8000/` | HTTP 200, full HTML | Status 200, modals and script present | **PASS** |
| **Nearby API with Sort** | `GET /nearby?sort=price` | HTTP 200, sorted cards | Status 200, cards sorted by min price | **PASS** |
| **Map Markers JSON** | `GET /api/markers` | JSON array of pins | Array of verified restaurants returned | **PASS** |
| **Restaurant Modal** | `GET /restaurant/1` | HTMX menu markup | Vidyarthi Bhavan menu rendered | **PASS** |
| **Community Nomination** | `POST /api/suggest` | HTTP 200, success message | Status 200, saved as pending in PostgreSQL | **PASS** |

---

## 6. Architecture Status & Next Opportunities

The core discovery engine, light-mode interface, and interaction loops are completely operational. Next logical expansions include:
1. **Photo Uploads:** Allow community members to submit photos with their eatery nominations.
2. **Review Aggregation:** Add user ratings and short text reviews for individual dishes.
3. **PWA Offline Support:** Cache bookmarked eateries and maps for offline street food hunting.
