# UNDER100: CUSTOMER, MERCHANT & ADMIN INTERFACE SEPARATION AUDIT & IMPLEMENTATION REPORT

**Application:** Under100 Food Discovery Platform  
**Repository Directory:** `C:\Users\astha\Desktop\food-discovery`  
**Date of Audit & Execution:** October 21, 2026  
**Engineering Disciplines:** Full-Stack Architecture, UI/UX Design & Application Security  
**Test Suite Status:** **29 passed, 0 failed (100% pass rate)**  
**Database Record Preservation:** **100% Preserved (Zero data loss or schema regression)**  

---

## 1. Executive Summary

The Under100 application has been audited, redesigned, and refactored to achieve **strict three-way separation** across its user interfaces and backend access controls:

1. **Customer-Facing Storefront (`/`)**:
   - Modern, high-conversion food discovery experience inspired by leading platforms (Blinkit, Swiggy, Zomato).
   - **Zero operational clutter**: Completely eradicated restaurant onboarding banners, "Partner With Us" buttons, merchant portal links, admin console triggers, 1-click role switches, and eatery claiming forms from the public browsing flow.
   - Clean customer account navigation focused exclusively on buyer needs (Addresses, Orders, Prescriptions, E-Gift Cards, Privacy, Logout).

2. **Restaurant Owner / Merchant Portal (`/merchant`, `/merchant/login`)**:
   - Distinct, partner-branded portal with dedicated green theme (`#059669`) and enterprise navigation.
   - Isolated multi-tenant dashboard allowing restaurant owners to manage only their own outlets, menus, item availability, operating hours, and temporary closures.
   - Strict server-side authorization: unauthenticated requests are redirected with HTTP 303 to `/merchant/login?msg=merchant_required`; non-merchant customers attempting to access `/merchant` or its mutation endpoints receive **HTTP 403 Forbidden**.

3. **Platform Administrator Dashboard (`/admin`, `/admin/login`)**:
   - Dedicated administrative console with slate-navy branding (`#0f172a`), platform health monitoring, user management, eatery approval workflows, category taxonomies, and review moderation.
   - Strict server-side enforcement: unauthenticated users redirect to `/admin/login?msg=admin_required`; any non-admin attempting access receives **HTTP 403 Forbidden**.

4. **Security Hardening & Role-Based Routing**:
   - Dedicated entry routes (`/merchant/login` and `/admin/login`) with role validation.
   - Unified `/login` role-based auto-routing (merchants automatically directed to `/merchant`, admins to `/admin`, customers to `/`).
   - Open redirect defense: robust sanitization on `return_to` parameters prevents external URLs (`https://`, `//`, `/\`) and blocks privilege escalation attempts (`return_to=/admin`).
   - Registration privilege escalation blocked: public `/register` form strictly assigns `role = "customer"`, ignoring any user-submitted role parameters.

---

## 2. Initial Audit Findings

Prior to refactoring, an exhaustive code and UI audit revealed several architectural and security vulnerabilities:

| # | Vulnerability / Issue | Location | Risk Severity |
|---|---|---|---|
| **1** | **Partner Navigation in Header** | `RenderHomePage` desktop navbar | High UI Clutter |
| | Links `btn-partner-nav` (`/partner/register`), `btn-nav-merchant-direct` (`/merchant`), and `btn-nav-admin-direct` (`/admin`) were rendered in the main customer header. | |
| **2** | **Account Popover Role Switchers** | `RenderHomePage` user dropdown menu | Critical Security & UX |
| | Rendered "Quick Switch (Dev Tools)" buttons and direct links to `/admin` and `/merchant`, allowing any logged-in user to navigate to administrative pages. | |
| **3** | **Footer Restaurant Column** | `RenderHomePage` footer | Moderate UI Clutter |
| | Contained a prominent "FOR RESTAURANTS" column with "Partner With Us", "Merchant Owner Portal", and "Admin Database Console" links on every public page. | |
| **4** | **Homepage Login Modal Role Shortcuts** | `#loginModalOverlay` | High Security Risk |
| | Contained "Fast Role Access (Demo)" buttons allowing 1-click login as Admin or Merchant directly from the customer homepage modal, plus an owner banner. | |
| **5** | **Customer Restaurant Detail Claim Banner** | `RenderRestaurantDetailsModal` | Moderate UI Clutter |
| | Showed a prominent green "Claim This Eatery" banner inside the customer food ordering modal. | |
| **6** | **Role Escalation in Registration** | `POST /register` | High Security Risk |
| | The registration form presented a `<select name="role">` dropdown that could be tampered with to self-assign merchant or admin roles. | |
| **7** | **Shared Login Demo Box** | `RenderLoginPage` | Moderate Security Risk |
| | Showed 1-click test credentials and account-type selector on public `/login`. | |
| **8** | **Missing Dedicated Portal Login Pages** | Routes `/merchant/login` & `/admin/login` | Usability & Brand Separation |
| | Unauthenticated access to `/merchant` or `/admin` dumped users onto the generic consumer `/login` without role-specific context. | |

---

## 3. Changes Implemented

### 3.1 Customer Website Cleanliness (`views.go`)
- **Desktop Navbar**:
  - Removed `btn-partner-nav` (`/partner/register`).
  - Removed `btn-nav-merchant-direct` (`/merchant`).
  - Removed `btn-nav-admin-direct` (`/admin`).
  - Replaced with customer-focused navigation: Location selector, Search bar, Favorites shortcut, and User Account menu.
- **Account Dropdown Popover**:
  - Removed all "Quick Switch" role buttons.
  - Removed direct links to `/merchant` and `/admin`.
  - Replaced with clean customer menu:
    - *My Account* (`/account`)
    - *My Addresses* (`/account/addresses`)
    - *My Orders* (`/account/orders`)
    - *My Prescriptions* (`/account/prescriptions`)
    - *E-Gift Cards* (`/account/giftcards`)
    - *Account Privacy* (`/account/privacy`)
    - *Logout* (`/logout`)
- **Public Footer**:
  - Replaced the "FOR RESTAURANTS" column with a customer-centric **"CUSTOMER CARE"** column:
    - *About Under100*
    - *Privacy & Terms*
    - *Order Help & Support*
    - *Saved Favorites* (`/favorites`)
- **Customer Sign-In Modal (`#loginModalOverlay`)**:
  - Removed "Fast Role Access" demo buttons and client-side `loginDirectRole` function.
  - Removed the bottom partner banner (`Partner with us & add your restaurant`).
- **Restaurant Details Modal**:
  - Removed the `ownerClaimBanner` so customers only see menus, pricing, dietary tags, operating hours, and reviews.
- **Public Login Page (`/login`)**:
  - Removed the demo credentials box.
  - Removed the `<select name="role">` dropdown from the sign-up tab.
  - Added hidden `return_to` support to safely preserve navigation intent.

### 3.2 Dedicated Merchant Portal & Security (`merchant_handlers.go`, `views.go`, `food_handlers.go`)
- **Dedicated Merchant Login (`RenderMerchantLoginPage`)**:
  - Available at `/merchant/login`.
  - Styled with partner-centric green branding (`#059669`), emerald badges, and clear restaurant owner guidance.
  - Displays context messages (e.g. `?msg=merchant_required`).
- **Dedicated Merchant Auth Handler (`MerchantLoginSubmit`)**:
  - Validates credentials against `users` table.
  - Server-side check: requires `user.Role == user.RoleMerchant` or `user.IsAdmin == true`.
  - If a customer account attempts login at `/merchant/login`, it returns **HTTP 403 Forbidden** with an explicit error message.
- **Merchant Auth Middleware (`MerchantAuthMiddleware`)**:
  - Unauthenticated requests to any `/merchant` endpoint are redirected with **HTTP 303** to `/merchant/login?msg=merchant_required` (with `HX-Redirect` support for HTMX requests).
  - Authenticated customers attempting to access `/merchant` or perform mutation endpoints (e.g., `POST /merchant/menu`) receive **HTTP 403 Forbidden**.
- **Multi-Tenant Ownership Isolation**:
  - Merchant dashboard queries only outlets where `user_restaurants = user.ID`.
  - Mutation endpoints verify ownership of the target eatery and dish; unauthorized cross-tenant mutations return **HTTP 403 Forbidden**.

### 3.3 Platform Administrator Panel & Access Controls (`food_handlers.go`, `views.go`)
- **Dedicated Admin Login (`RenderAdminLoginPage`)**:
  - Available at `/admin/login`.
  - Styled with slate-navy administrative branding (`#0f172a`), shield icons, and platform console notices.
- **Dedicated Admin Auth Handler (`AdminLoginSubmit`)**:
  - Validates credentials against `users` table.
  - Strictly requires `u.IsAdmin == true || u.Role == user.RoleAdmin`. Non-admins receive **HTTP 403 Forbidden**.
- **Admin Auth Middleware (`AdminAuthMiddleware`)**:
  - Unauthenticated requests to `/admin` redirect with **HTTP 303** to `/admin/login?msg=admin_required`.
  - Non-admin users attempting to access `/admin` or execute administrative POST endpoints (e.g., `POST /admin/restaurants/:id/status`) receive **HTTP 403 Forbidden**.

### 3.4 Role-Based Routing & Open Redirect Prevention (`food_handlers.go`)
- **Sanitized Return URLs (`sanitizeReturnURL`)**:
  - Rejects URLs containing scheme delimiters (`:`), protocol-relative prefixes (`//`), backslash bypasses (`/\`), or control characters.
  - Only allows safe, relative paths starting with a single `/`.
- **Intelligent Role Routing (`LoginSubmit`)**:
  - When a user logs in via the shared `/login` route:
    - Admins are automatically routed to `/admin`.
    - Restaurant owners are automatically routed to `/merchant`.
    - Customers are routed to the sanitized `return_to` path (or `/` if empty/external/privileged like `/admin`).
- **Privilege Escalation Defense (`RegisterSubmit`)**:
  - Overrides any client-submitted `role` field; all public self-registrations are hardcoded to `role = user.RoleCustomer` and `is_admin = false`.

---

## 4. Modified & Created Files Summary

| File Path | Action | Description |
|---|---|---|
| `pkg/routing/fooddiscovery/views.go` | Modified | Eradicated merchant/admin links from header, account dropdown, footer, modals, and login view. Created `RenderMerchantLoginPage` and `RenderAdminLoginPage`. |
| `pkg/routing/fooddiscovery/food_handlers.go` | Modified | Registered `/merchant/login` and `/admin/login` GET & POST routes. Implemented `MerchantLoginPage`, `MerchantLoginSubmit`, `AdminLoginPage`, `AdminLoginSubmit`, and `sanitizeReturnURL`. Hardcoded customer role in `RegisterSubmit`. |
| `pkg/routing/fooddiscovery/merchant_handlers.go` | Modified | Updated `MerchantAuthMiddleware` to redirect unauthenticated requests to `/merchant/login?msg=merchant_required`. |
| `verify_separation.ps1` | Created | Comprehensive 29-step automated verification suite testing storefront cleanliness, merchant isolation, admin security, redirects, and database preservation. |
| `CUSTOMER_MERCHANT_ADMIN_SEPARATION_REPORT.md` | Created | This architectural audit and compliance report. |

---

## 5. Verification Test Suite Results

The automated test script `verify_separation.ps1` was executed against the running application stack (PostgreSQL + Redis + Echo Go Server on port 8000).

**Summary: 29 Passed, 0 Failed (100% Pass Rate)**

```
=================================================================
 UNDER100: THREE-WAY INTERFACE SEPARATION & SECURITY TEST SUITE   
=================================================================

-- PART 1: Customer Website Cleanliness & Storefront Functionality --
 [PASS] Test 1.1: Customer Homepage Loads
        Homepage returned valid HTML markup with Under100 branding
 [PASS] Test 1.2: Customer Header Free of Management Links
        Desktop header contains zero merchant/admin entry points
 [PASS] Test 1.3: Customer Footer Free of Management Controls
        Public footer replaced with customer-facing CUSTOMER CARE links
 [PASS] Test 1.4: Customer Login Modal Clean of Role Shortcuts
        Login modal contains only standard phone/email entry
 [PASS] Test 1.5: Restaurant Detail Modal Free of Claim Banner
        Restaurant detail page shows food items and prices without management actions
 [PASS] Test 1.6: Customer Login & Sign-up Page Free of Role Escalation
        /login presents focused customer sign-in/sign-up without role dropdown

-- PART 2: Restaurant Owner Portal & Multi-Tenant Security --
 [PASS] Test 2.1: Dedicated Merchant Login Page Accessible
        /merchant/login presents dedicated partner portal sign-in
 [PASS] Test 2.2: Unauthenticated Access to /merchant Protected
        Unauthenticated user redirected to /merchant/login
 [PASS] Test 2.3: Customer Blocked from Merchant Portal (403)
        Customer account receives 403 Forbidden when requesting /merchant
 [PASS] Test 2.4: Customer Blocked from Merchant Mutations (403)
        Direct POST to /merchant/menu blocked with 403 Forbidden
 [PASS] Test 2.5: Authorized Merchant Signs In & Accesses Dashboard
        Merchant authenticated via /merchant/login and dashboard loaded
 [PASS] Test 2.6: Merchant Blocked from Admin Panel (403)
        Merchant account receives 403 Forbidden when attempting to access /admin
 [PASS] Test 2.7: Cross-Tenant Dish Mutation Blocked (403)
        Merchant cannot edit dishes belonging to another restaurant

-- PART 3: Administrator Panel & Access Controls --
 [PASS] Test 3.1: Dedicated Admin Login Page Accessible
        /admin/login presents dedicated administrative sign-in
 [PASS] Test 3.2: Unauthenticated Access to /admin Protected
        Unauthenticated user redirected to /admin/login
 [PASS] Test 3.3: Customer Blocked from Admin Panel (403)
        Customer account receives 403 Forbidden when requesting /admin
 [PASS] Test 3.4: Customer Blocked from Admin Mutations (403)
        Direct POST to /admin/restaurants/1/status blocked with 403 Forbidden
 [PASS] Test 3.5: Authorized Admin Signs In & Accesses Dashboard
        Admin authenticated via /admin/login and admin panel loaded
 [PASS] Test 3.6: Customer Rejected at /admin/login Entry Point (403)
        Non-admin credentials rejected with 403 Forbidden on /admin/login

-- PART 4: Role-Based Auto-Redirects & Open Redirect Prevention --
 [PASS] Test 4.1: Customer Login Redirects to Storefront (/)
        Customer authenticated and redirected to /
 [PASS] Test 4.2: Merchant Login via Shared /login Redirects to /merchant
        Merchant role detected and auto-routed to /merchant
 [PASS] Test 4.3: Admin Login via Shared /login Redirects to /admin
        Admin role detected and auto-routed to /admin
 [PASS] Test 4.4: Safe Internal Return URL Honored
        User redirected to requested customer return URL (/favorites)
 [PASS] Test 4.5: Open Redirect (https://evil.com) Blocked
        External absolute URL rejected, safely redirected to /
 [PASS] Test 4.6: Open Redirect (//evil.com) Blocked
        Protocol-relative URL rejected, safely redirected to /
 [PASS] Test 4.7: Open Redirect (/\\evil.com) Blocked
        Backslash bypass URL rejected, safely redirected to /
 [PASS] Test 4.8: Customer Cannot Use return_to to Escalate to /admin
        Return URL pointing to /admin ignored for customer, redirected to /
 [PASS] Test 4.9: Registration Role Parameter Cannot Escalate Privileges
        Submitted role=merchant ignored: user created as role=customer, is_admin=false

-- PART 5: Database Integrity Check --
 [PASS] Test 5.1: 100% Database Records Preserved
        All database tables match exact pre-test baseline counts

=================================================================
 TEST SUITE COMPLETE: 29 passed, 0 failed out of 29 tests
=================================================================
```

---

## 6. Database Preservation Verification

Per architectural instructions, zero database reset, truncation, or data loss occurred. Exact record counts were queried via `psql` before and after the verification test suite:

| Database Table | Pre-Execution Baseline | Post-Execution Final | Integrity Status |
|---|---|---|---|
| `users` | 9 | 9 | **PRESERVED (100%)** |
| `restaurants` | 42 | 42 | **PRESERVED (100%)** |
| `food_items` | 166 | 166 | **PRESERVED (100%)** |
| `categories` | 12 | 12 | **PRESERVED (100%)** |
| `food_item_categories` | 228 | 228 | **PRESERVED (100%)** |
| `restaurant_operating_hours` | 7 | 7 | **PRESERVED (100%)** |
| `reviews` | 0 | 0 | **PRESERVED (100%)** |

---

## 7. Operational Runbook

To run and verify the Under100 system in local development:

### 7.1 Start Supporting Daemons
```powershell
# 1. PostgreSQL 15 Daemon
& "C:\Users\astha\pgsql\bin\postgres.exe" -D "C:\Users\astha\pgdata"

# 2. Redis Cache Daemon
& "C:\Users\astha\redis-server.exe"
```

### 7.2 Compile and Run the Server
```powershell
cd C:\Users\astha\Desktop\food-discovery
$env:GOTMPDIR = "C:\Users\astha\Desktop\food-discovery\scratch\gotmp"
go build -o under100_server.exe ./cmd/web
.\under100_server.exe
```

### 7.3 Execute Verification Suite
```powershell
powershell.exe -ExecutionPolicy Bypass -File .\verify_separation.ps1
```

### 7.4 Standard Seeded Credentials for Testing
- **Customer Storefront (`/login`)**:
  - Email: `user@fooddiscovery.com`
  - Password: `User@123456`
- **Restaurant Owner Portal (`/merchant/login`)**:
  - Email: `merchant.ctr@fooddiscovery.com`
  - Password: `Merchant@123456`
- **Platform Administrator Console (`/admin/login`)**:
  - Email: `admin@fooddiscovery.com`
  - Password: `Admin@123456`

---

## 8. Conclusion

The Under100 food discovery platform now features three strictly isolated, access-controlled interfaces. The customer experience is clean, modern, and uncluttered like Blinkit, Swiggy, or Zomato, while restaurant owners and platform administrators operate through dedicated, security-hardened portals. Zero database records were lost, and all 29 automated test cases pass with 100% compliance.
