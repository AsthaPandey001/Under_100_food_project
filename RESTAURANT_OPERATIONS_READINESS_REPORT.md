# UNDER100: RESTAURANT OPERATIONS & CUSTOMER READINESS REPORT

**Project:** Under100 — Bangalore Food Discovery Application  
**Execution Date:** October 21-22, 2026  
**Auditor & Lead Architect:** Senior Go Backend Engineer, Product Architect & Application Security Specialist  
**Verification Status:** 30 / 30 Operations Tests PASSED (100%) | 29 / 29 Separation Tests PASSED (100%)  
**Build Status:** `go build -o under100_server.exe ./cmd/web` (Exit Code: 0)  
**Static Analysis:** `go vet ./...` (Exit Code: 0, Clean)  
**Database Integrity:** 100% Preserved (Zero baseline records added or destroyed)  

---

## 1. Executive Summary & Audit Findings

Following the completion of the three-way interface separation (Customer Storefront, Restaurant Owner Portal, Admin Console), a deep-dive operational and security audit was conducted on the Under100 codebase (`C:\Users\astha\Desktop\food-discovery`).

### Key Findings & Reality Checks:
1. **Honest Application Stage (Discovery vs. Delivery/Fulfillment):**
   - The Under100 database contains *no* tables for carts, checkout sessions, payment intents, dispatch queues, or rider tracking.
   - The merchant portal dashboard previously rendered a mock "Kitchen Order Queue" with hard-coded dummy data (`#BLR-84192`, `Rava Masala Dosa x 2`, "Accepted - Food being prepared").
   - **Correction:** The mock kitchen queue was removed and replaced with an authentic **Operations & Delivery Configuration** card. This prevents misleading restaurant partners and customers into believing the platform dispatches riders or processes in-flight orders.
2. **Missing Under-₹100 Price Cap on Merchant Menu Mutations:**
   - While customer discovery enforced the Under-₹100 limit, the merchant endpoint `/merchant/menu` permitted dishes with arbitrary prices (e.g., ₹150) to be submitted and updated without server-side validation.
   - **Correction:** Strict server-side validation (`price > 0 && price <= 100.0`) was implemented for dish creation and editing in both merchant and admin handlers.
3. **Cross-Tenant Fallback Bug in Merchant Restaurant Resolution:**
   - In `getMerchantRestaurant(ctx)`, when a merchant specified a query parameter `?restaurant_id=X` that they did *not* own, the function silently fell back to returning the merchant's first owned restaurant instead of failing closed (`nil`).
   - **Correction:** Fixed to fail closed. Attempting to mutate an unowned restaurant yields a clean `403 Forbidden` / `Restaurant not found`.
4. **Operating Hours (Phase 10) Multi-Outlet Scope:**
   - Phase 10 operating hours and temporary closures were functional in `pkg/repos/fooddiscovery/availability.go`, but the merchant HTTP handlers did not propagate the specific `restaurant_id` across multi-outlet accounts.
   - **Correction:** Multi-outlet restaurant selection was wired into all operating hours and temporary closure handlers.
5. **Delivery Serviceability Architecture:**
   - Missing persistent delivery settings for restaurants (radius, delivery fee, enabled/paused toggles, temporary monsoon/kitchen suspensions).
   - **Correction:** Implemented `restaurant_delivery_settings` with automatic schema initialization, merchant configuration modal, and customer-facing serviceability calculation (checking radius, distance, and temporary delivery suspensions).

---

## 2. Core Architectural Principles & Operational Distinctions

To ensure clarity for developers, restaurant partners, and customers, the operational state model is formalized into four distinct tiers:

| Operational Dimension | Status Source | Meaning to Customer | Meaning to Merchant |
| :--- | :--- | :--- | :--- |
| **1. Restaurant Opening Status** | Operating periods & temporary closures in `restaurant_operating_hours` / `restaurant_special_hours` | Indicates if the physical eatery is currently open (Asia/Kolkata timezone). Displays "Open Now", "Closed", or "Temporarily Closed". | Configured via `/merchant/hours`. Used for customer expectation and dine-in/pickup timing. |
| **2. Dish Availability** | `available` boolean on `food_items` | An individual dish is either active or "Sold Out". Sold-out dishes cannot be added to order plans. | Toggled via `/merchant/menu/:id/availability`. Immediate sync to customer modal. |
| **3. Delivery Serviceability** | `restaurant_delivery_settings` evaluated against customer coordinates | Evaluates whether delivery is offered, paused, or whether customer is within `delivery_radius_km`. | Configured via `/merchant/delivery`. Sets radius, fee, and temporary weather pause. |
| **4. Order Acceptance** | *Not applicable in discovery stage* | Transparently communicated that Under100 is an authentic discovery platform; direct ordering/payments belong to a future phase. | Prevents fake kitchen tickets and false rider dispatch promises. |

---

## 3. Features Implemented & Fixed

### A. Merchant Profile & Multi-Tenant Management
- **Authentic Operations Card:** Replaced fake kitchen order card with live operational controls for Delivery Settings, Operating Hours, and Edit Outlet Profile.
- **Outlet Editing:** Enabled merchants to edit Name, Description, Address, Phone, Website, and GPS coordinates for their specific outlets.
- **Cross-Tenant Isolation:** Validated that a merchant cannot view, edit, or submit changes to restaurants owned by other partners.

### B. Menu Management & Price Cap Enforcement
- **Server-Side ₹100 Price Cap:** Hard-enforced in `MerchantCreateFood`, `MerchantUpdateFood`, `AdminAddFood`, and `AdminEditFood`. Rejects any dish with price <= 0 or > 100 with HTTP 400 Bad Request.
- **Dietary Classification & Categorization:** Automatically associates dishes with category slugs (`category_repo`) and tags vegetarian / non-vegetarian designations.
- **Sold Out Handling:**
  - Merchants can toggle dish availability (`POST /merchant/menu/:id/availability`).
  - Unavailable dishes are prominently badged as **"Sold Out"** with reduced opacity (0.6) and disabled interaction on the customer modal.

### C. Operating Hours & Temporary Closures (Phase 10 Integrity)
- **Asia/Kolkata Timezone Engine:** Evaluates standard weekly periods (days 1-7) against IST (`time.Now().In(GetKolkataLocation())`).
- **Temporary Maintenance Closure:** Merchants can set temporary closure reasons (e.g., "Deep kitchen cleaning and maintenance") and expected reopening timestamps via `POST /merchant/hours/temporary-closure`.
- **Customer Storefront Reflection:** Customer modal displays an alert notice with the closure reason and reopening schedule, accurately overriding weekly operating hours.

### D. Delivery Settings & Customer Serviceability
- **Delivery Repository (`delivery_repo.go`):**
  - Schema: `restaurant_delivery_settings (restaurant_id PK/FK, delivery_enabled, pickup_enabled, delivery_radius_km, delivery_fee, min_order_amount, estimated_delivery_min, is_delivery_suspended, suspension_reason)`.
  - Idempotent table creation via `EnsureDeliverySettingsTable`.
  - Batch retrieval optimization (`GetDeliverySettingsBatch`) to prevent N+1 queries during feed rendering.
- **Merchant Delivery Modal & Handlers:**
  - Route: `GET /merchant/delivery` & `POST /merchant/delivery`.
  - Configurable delivery radius (km), flat fee (₹), estimated delivery time, and temporary suspension toggle with reason.
- **Customer-Facing Serviceability:**
  - Evaluated via Haversine distance from customer location to restaurant.
  - If within radius: Displays `🚚 Delivery: ₹X (~Y mins)` badge.
  - If beyond radius: Displays `🚫 Delivery Not Available (Beyond Z km radius)`.
  - If suspended: Displays `⏸️ Delivery Temporarily Paused (Reason)`.

### E. Restaurant Approval Lifecycle Verification
- **Partner Registration (`/partner/register`):** New restaurant partners register an account and submit outlet details.
- **Pending Review State:** New restaurant enters `status = 'pending'`.
- **Storefront Suppression:** Pending restaurant is completely invisible in customer discovery search (`/nearby`) and direct customer URLs (`/restaurant/:id` returns 404).
- **Merchant Pre-Launch Setup:** Partner can log into `/merchant`, sees an "Under Administrative Review" notification banner, and can safely configure their menu and operating hours.
- **Admin Review & Approval (`/admin/restaurants/:id/status`):** Platform admin reviews and transitions status to `approved`.
- **Public Discovery:** Once approved and containing dishes <= ₹100, the restaurant appears immediately in Bangalore customer discovery.
- **Rejection Suppression:** If admin transitions status to `rejected`, the restaurant is instantly suppressed from public customer discovery.

---

## 4. Exact Files Modified & Created

| File | Change Type | Description |
| :--- | :--- | :--- |
| `pkg/repos/fooddiscovery/delivery_repo.go` | **CREATED** | Complete delivery settings repository, serviceability evaluator, and schema management. |
| `pkg/repos/fooddiscovery/food_discovery.go` | **MODIFIED** | Added delivery settings table initialization, extended `RestaurantWithDistance` with `ServiceabilityStatus`, and integrated batch delivery calculation in paginated feeds. |
| `pkg/routing/fooddiscovery/merchant_views.go` | **MODIFIED** | Removed mock kitchen order card; added authentic operations card, delivery settings modal, and multi-outlet profile edit fields. |
| `pkg/routing/fooddiscovery/merchant_handlers.go` | **MODIFIED** | Added Under-₹100 validation on create/update food, fixed cross-tenant fallback in `getMerchantRestaurant`, and implemented `MerchantDeliveryView` / `MerchantSaveDeliverySettings`. |
| `pkg/routing/fooddiscovery/hours_handlers.go` | **MODIFIED** | Added `restaurant_id` propagation for multi-outlet operating hours and closure management. |
| `pkg/routing/fooddiscovery/views.go` | **MODIFIED** | Formatted delivery serviceability badges on restaurant cards, added delivery details to restaurant modal, and styled Sold Out dish badges. |
| `pkg/routing/fooddiscovery/food_handlers.go` | **MODIFIED** | Registered delivery routes, calculated delivery serviceability in `RestaurantDetails`, and enforced ₹100 cap on admin food handlers. |
| `verify_restaurant_operations.ps1` | **CREATED** | 30-step end-to-end operational, security, and lifecycle verification suite. |

---

## 5. Verification & Test Execution Results

### 1. Operations Verification Suite (`verify_restaurant_operations.ps1`)
**Result: 30 / 30 PASSED (100%)**

- `[PASS] Test 1.1`: Merchant Dashboard Displays Operations & Fulfillment Card
- `[PASS] Test 1.2`: Merchant Successfully Updates Outlet Profile
- `[PASS] Test 1.3`: Cross-Tenant Outlet Profile Editing Blocked
- `[PASS] Test 2.1`: Menu Creation Rejects Dishes > Rs. 100 (Price Cap Enforced)
- `[PASS] Test 2.2`: Menu Creation Accepts Valid Dish <= Rs. 100
- `[PASS] Test 2.3`: Menu Update Rejects Price Escalation > Rs. 100
- `[PASS] Test 2.4`: Cross-Tenant Dish Mutation Blocked (403)
- `[PASS] Test 3.1`: Merchant Marks Dish as Sold Out
- `[PASS] Test 3.2`: Customer Storefront Displays Sold Out Badge
- `[PASS] Test 3.3`: Merchant Restores Dish Availability
- `[PASS] Test 4.1`: Merchant Operating Hours Modal with Asia/Kolkata Timezone
- `[PASS] Test 4.2`: Merchant Activates Temporary Closure
- `[PASS] Test 4.3`: Customer Storefront Reflects Temporary Closure Status
- `[PASS] Test 4.4`: Merchant Clears Temporary Closure
- `[PASS] Test 5.1`: Merchant Delivery Configuration Modal Accessible
- `[PASS] Test 5.2`: Merchant Configures Delivery Radius & Service Fee
- `[PASS] Test 5.3`: Customer Within Radius Evaluated as Serviceable
- `[PASS] Test 5.4`: Customer Outside Radius Correctly Shown Unserviceable
- `[PASS] Test 5.5`: Temporary Delivery Suspension Displayed to Customers
- `[PASS] Test 6.1`: Partner Self-Service Registration
- `[PASS] Test 6.2`: Submitted Restaurant Enters Pending Review State
- `[PASS] Test 6.3`: Pending Restaurant Hidden from Public Discovery
- `[PASS] Test 6.4`: Owner Has Portal Access with Pending Review Notice
- `[PASS] Test 6.5`: Admin Approves Restaurant Submission
- `[PASS] Test 6.6`: Approved Restaurant Appears in Customer Discovery
- `[PASS] Test 6.7`: Admin Rejects Restaurant Listing
- `[PASS] Test 6.8`: Rejected Restaurant Removed from Discovery
- `[PASS] Test 7.1`: Customer Favorites Page Accessible
- `[PASS] Test 7.2`: Search and Filter Workflow Operational
- `[PASS] Test 7.3`: Database Baseline Completely Preserved

### 2. Interface Separation Regression Suite (`verify_separation.ps1`)
**Result: 29 / 29 PASSED (100%)**
- Zero regressions across customer storefront cleanliness, dedicated merchant portal security, admin access controls, and open redirect protections.

### 3. Static Analysis & Unit Tests
- `go vet ./...`: 0 errors / clean exit.
- `go test ./pkg/repos/fooddiscovery`: **PASS** (`availability_test.go`, `category_repo_test.go`).
- `go test ./pkg/context`, `./pkg/funcmap`, `./pkg/htmx`, `./pkg/middleware`, `./pkg/controller`: **PASS**.
- *(Note: `emailsmanager` test expects a local Docker daemon for `testcontainers`, an environmental limitation on Windows).*

---

## 6. Database Baseline & Zero-Data-Loss Proof

No database reset, truncation, or reseed occurred. Test fixtures were isolated and cleaned up automatically.

| Entity / Table | Baseline Pre-Audit | Final Post-Verification | Differential | Integrity Status |
| :--- | :--- | :--- | :--- | :--- |
| `users` | 9 | 9 | +0 | **100% Preserved** |
| `restaurants` | 43 | 43 | +0 | **100% Preserved** |
| `food_items` | 171 | 171 | +0 | **100% Preserved** |
| `categories` | 12 | 12 | +0 | **100% Preserved** |
| `food_item_categories` | 228 | 228 | +0 | **100% Preserved** |
| `restaurant_operating_hours` | 7 | 7 | +0 | **100% Preserved** |
| `reviews` | 0 | 0 | +0 | **100% Preserved** |

---

## 7. Known Limitations & Recommendations for Future Phases

1. **Ordering & Fulfillment Phase (Future Scope):**
   - Under100 is strictly a food discovery and menu serviceability platform.
   - When the business decides to introduce transactional ordering, a dedicated phase must implement:
     - Cart session storage (in Redis or PostgreSQL).
     - Payment gateway integration (Razorpay / UPI / Stripe).
     - Kitchen order acceptance state machine (`placed` -> `accepted` -> `preparing` -> `ready_for_pickup`).
     - Rider dispatch and GPS tracking integrations.
2. **Photo Management:**
   - Dish and restaurant schemas currently store static image paths or URLs (`/static/...`). If image file uploads are desired for merchants, an S3 / cloud storage upload pipeline should be attached to the merchant edit modals.
3. **Multi-Location Geofencing:**
   - Serviceability currently uses Haversine radial distance from the restaurant's stored coordinates. If polygonal delivery zones (e.g., Ward / Postcode boundaries) are required, PostGIS polygon queries can replace radial math.

---

## 8. Conclusion

The Under100 restaurant operations layer is now robust, secure, and production-ready for the Bangalore food discovery ecosystem. All mock fulfillment elements have been replaced with authentic operations controls, tenant isolation is strictly enforced, prices are validated against the ₹100 threshold on the server, operating hours and temporary closures respect Asia/Kolkata timezone rules, and the restaurant approval lifecycle operates with full integrity.
