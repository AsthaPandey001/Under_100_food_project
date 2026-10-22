package fooddiscovery

import (
	"fmt"
	"strings"

	"github.com/mikestefanello/pagoda/ent"
	foodrepo "github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// RenderMerchantDashboard renders the dedicated, responsive Merchant Owner Portal with multi-restaurant management
func RenderMerchantDashboard(u *ent.User, rest *ent.Restaurant, allRests []*ent.Restaurant, foods []*ent.FoodItem, delivSettings foodrepo.DeliverySettings, flashMsg, flashErr string, itemRules ...map[int]foodrepo.FoodItemDeliveryRule) string {
	userName := u.Name
	if userName == "" {
		userName = u.Email
	}

	flashHTML := ""
	if flashMsg == "restaurant_registered" {
		flashHTML = `<div class="merchant-alert success" style="background: #f0fdf4; border-color: #bbf7d0; color: #166534;">🎉 Your restaurant outlet has been registered successfully! You can now manage your menu dishes and operating hours below.</div>`
	} else if flashMsg == "outlet_claimed" {
		flashHTML = `<div class="merchant-alert success" style="background: #f0fdf4; border-color: #bbf7d0; color: #166534;">🎉 Eatery listing claimed successfully! You now have full management access to this outlet.</div>`
	} else if flashMsg == "welcome_merchant" {
		flashHTML = `<div class="merchant-alert success" style="background: #f0fdf4; border-color: #bbf7d0; color: #166534;">🎉 Welcome to FoodDiscovery! Your merchant partner account is active. Start adding dishes to join the Under ₹100 feed.</div>`
	} else if flashMsg != "" {
		flashHTML = fmt.Sprintf(`<div class="merchant-alert success">%s</div>`, escapeHTML(flashMsg))
	}
	if flashErr != "" {
		flashHTML = fmt.Sprintf(`<div class="merchant-alert error">%s</div>`, escapeHTML(flashErr))
	}

	mainContent := ""
	if rest == nil {
		// Empty state: Merchant has no assigned restaurant
		mainContent = fmt.Sprintf(`
		<div class="empty-outlet-card">
			<div class="empty-icon-box">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#e23744" stroke-width="1.8">
					<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
					<polyline points="9 22 9 12 15 12 15 22"></polyline>
				</svg>
			</div>
			<h2>Welcome, %s!</h2>
			<p>You have not registered or claimed an eatery on FoodDiscovery yet. Choose an option below to get started:</p>

			<div style="margin: 1.75rem 0; display: flex; gap: 0.85rem; justify-content: center; flex-wrap: wrap;">
				<button 
					type="button" 
					class="btn-primary" 
					style="font-size: 0.95rem; padding: 0.85rem 1.6rem; display: inline-flex; align-items: center; gap: 0.5rem; background: #087443; color: white; border: none; border-radius: 10px; font-weight: 700; cursor: pointer;"
					hx-get="/merchant/restaurant/claim" 
					hx-target="#merchant-modal"
				>
					<span>📍</span>
					Claim an Existing Bengaluru Eatery
				</button>
				<button 
					type="button" 
					class="btn-secondary" 
					style="font-size: 0.95rem; padding: 0.85rem 1.6rem; display: inline-flex; align-items: center; gap: 0.5rem; background: #ffffff; border: 1.5px solid #cbd5e1; border-radius: 10px; font-weight: 700; cursor: pointer; color: #1e293b;"
					hx-get="/merchant/restaurant/new" 
					hx-target="#merchant-modal"
				>
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
					Register a New Restaurant
				</button>
			</div>

			<div class="empty-outlet-help">
				<strong>Need platform assistance?</strong><br>
				Please contact platform support at <a href="mailto:admin@fooddiscovery.com">admin@fooddiscovery.com</a> or use the platform admin panel.
			</div>
			<div style="margin-top: 1.5rem;">
				<a href="/" class="btn-secondary">&larr; Return to Storefront</a>
			</div>
		</div>`, escapeHTML(userName))
	} else {
		// Calculate merchant metrics
		totalFoods := len(foods)
		availCount := 0
		unavailCount := 0
		under100Count := 0

		for _, f := range foods {
			if f.Available {
				availCount++
				if f.Price <= 100.0 {
					under100Count++
				}
			} else {
				unavailCount++
			}
		}

		// Under 100 Qualification Banner
		qualificationBanner := ""
		if string(rest.Status) == "pending" {
			qualificationBanner = `
			<div class="qualify-banner" style="background: #fefce8; border: 1.5px solid #fef08a;">
				<div class="qualify-icon">
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#854d0e" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
				</div>
				<div>
					<strong style="color: #854d0e;">Restaurant Listing Under Administrative Review</strong>
					<p style="color: #a16207;">Your restaurant has been submitted and is currently awaiting administrator approval. You can manage your menu below so your Under ₹100 items are ready upon approval.</p>
				</div>
			</div>`
		} else if string(rest.Status) == "rejected" {
			qualificationBanner = `
			<div class="qualify-banner" style="background: #fef2f2; border: 1.5px solid #fecaca;">
				<div class="qualify-icon">
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#991b1b" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
				</div>
				<div>
					<strong style="color: #991b1b;">Restaurant Listing Rejected</strong>
					<p style="color: #b91c1c;">Your listing was not approved for publication. Please contact platform support at admin@fooddiscovery.com for more details.</p>
				</div>
			</div>`
		} else if under100Count > 0 {
			qualificationBanner = fmt.Sprintf(`
			<div class="qualify-banner qualified">
				<div class="qualify-icon">
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#166534" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
				</div>
				<div>
					<strong>Discoverable in Bengaluru Under ₹100 Feed</strong>
					<p>Your restaurant has %d available dish(es) priced at or below ₹100, meeting the platform discovery requirement.</p>
				</div>
			</div>`, under100Count)
		} else {
			qualificationBanner = `
			<div class="qualify-banner unqualified">
				<div class="qualify-icon">
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#b45309" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
				</div>
				<div>
					<strong>Action Required: Add an Under ₹100 Dish</strong>
					<p>The platform features restaurants with at least 1 item priced ≤ ₹100. Add an affordable signature dish so customers can discover you.</p>
				</div>
			</div>`
		}

		statusBadge := `<span class="badge-status approved">Approved Listing</span>`
		if string(rest.Status) == "pending" {
			statusBadge = `<span class="badge-status pending">Pending Review</span>`
		} else if string(rest.Status) == "rejected" {
			statusBadge = `<span class="badge-status rejected">Rejected</span>`
		}

		// Delivery and serviceability display
		delivStateHTML := `<span style="color: #64748b; font-weight: 700;">❌ Delivery Disabled</span>`
		delivStateSub := "Dine-in / Takeaway only"
		if delivSettings.DeliveryEnabled {
			if delivSettings.IsDeliverySuspended {
				delivStateHTML = `<span style="color: #b45309; font-weight: 700;">⏸️ Temporarily Suspended</span>`
				delivStateSub = "Merchant paused active deliveries"
			} else {
				delivStateHTML = fmt.Sprintf(`<span style="color: #166534; font-weight: 700;">✅ Active (Within %.1f km)</span>`, delivSettings.DeliveryRadiusKm)
				delivStateSub = fmt.Sprintf("₹%.0f delivery fee per order", delivSettings.DeliveryFee)
			}
		}

		pickupStateHTML := `<span style="color: #64748b; font-weight: 700;">❌ Disabled</span>`
		if delivSettings.PickupEnabled {
			pickupStateHTML = `<span style="color: #166534; font-weight: 700;">✅ Supported</span>`
		}

		// Build menu items table rows
		rowsHTML := ""
		if len(foods) == 0 {
			rowsHTML = `<tr><td colspan="6" class="empty-table-cell">No food items added yet. Click "+ Add Dish" below to start your menu!</td></tr>`
		} else {
			var rulesMap map[int]foodrepo.FoodItemDeliveryRule
			if len(itemRules) > 0 {
				rulesMap = itemRules[0]
			}
			var sb strings.Builder
			for _, f := range foods {
				var rule foodrepo.FoodItemDeliveryRule
				if rulesMap != nil {
					rule = rulesMap[f.ID]
				} else {
					rule = foodrepo.DefaultFoodItemDeliveryRule(f.ID)
				}
				sb.WriteString(RenderMerchantSingleFoodRow(f, rule))
			}
			rowsHTML = sb.String()
		}

		outletSwitcherHTML := ""
		if len(allRests) > 0 {
			var optSb strings.Builder
			for _, r := range allRests {
				sel := ""
				if rest != nil && r.ID == rest.ID {
					sel = "selected"
				}
				optSb.WriteString(fmt.Sprintf(`<option value="%d" %s>%s — %s</option>`, r.ID, sel, escapeHTML(r.Name), escapeHTML(r.Address)))
			}

			outletCountBadge := fmt.Sprintf("%d Outlets Active", len(allRests))
			if len(allRests) == 1 {
				outletCountBadge = "1 Outlet Active"
			}

			outletSwitcherHTML = fmt.Sprintf(`
			<div class="outlet-switcher-bar" style="background: #ffffff; border: 1.5px solid #e2e8f0; border-radius: 14px; padding: 0.9rem 1.35rem; margin-bottom: 1.5rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem; box-shadow: 0 1px 4px rgba(0,0,0,0.03);">
				<div style="display: flex; align-items: center; gap: 0.85rem; flex-wrap: wrap;">
					<div style="width: 38px; height: 38px; border-radius: 10px; background: #f0fdf4; border: 1px solid #bbf7d0; display: flex; align-items: center; justify-content: center; font-size: 1.25rem;">
						🏪
					</div>
					<div>
						<div style="font-size: 0.72rem; font-weight: 800; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px;">Active Managing Outlet</div>
						<div style="display: flex; align-items: center; gap: 0.6rem; margin-top: 0.2rem;">
							<select class="form-input" style="padding: 0.35rem 0.75rem; font-size: 0.95rem; font-weight: 800; border-radius: 8px; border: 1.5px solid #cbd5e1; cursor: pointer; background: #ffffff; color: #0f172a; max-width: 420px;" onchange="window.location.href='/merchant?restaurant_id='+this.value">
								%s
							</select>
							<span style="background: #dcfce7; color: #166534; font-size: 0.75rem; font-weight: 800; padding: 0.2rem 0.6rem; border-radius: 9999px; border: 1px solid #bbf7d0;">%s</span>
						</div>
					</div>
				</div>
				<div style="display: flex; gap: 0.5rem; flex-wrap: wrap; align-items: center;">
					<button 
						type="button" 
						class="btn-secondary" 
						style="font-size: 0.82rem; padding: 0.5rem 0.85rem; border-radius: 8px; font-weight: 700; display: inline-flex; align-items: center; gap: 0.35rem; background: #ffffff; border: 1.5px solid #cbd5e1; cursor: pointer; color: #1e293b;"
						hx-get="/merchant/restaurant/claim" 
						hx-target="#merchant-modal"
					>
						<span>📍</span> <span>Claim Eatery</span>
					</button>
					<button 
						type="button" 
						class="btn-primary" 
						style="font-size: 0.82rem; padding: 0.5rem 0.85rem; border-radius: 8px; font-weight: 700; display: inline-flex; align-items: center; gap: 0.35rem; background: #087443; color: white; border: none; cursor: pointer;"
						hx-get="/merchant/restaurant/new" 
						hx-target="#merchant-modal"
					>
						<span>+</span> <span>Add Another Outlet</span>
					</button>
				</div>
			</div>`, optSb.String(), outletCountBadge)
		}

		mainContent = fmt.Sprintf(`
		%s
		%s

		<div class="outlet-header-card">
			<div class="outlet-info-col">
				<div class="outlet-meta-top">
					<span class="outlet-tag">Assigned Outlet</span>
					%s
				</div>
				<h1 class="outlet-name">%s</h1>
				<p class="outlet-desc">%s</p>
				<div class="outlet-contact-row">
					<span><strong>Address:</strong> %s</span>
					<span><strong>Phone:</strong> %s</span>
					<span><strong>Hours:</strong> %s</span>
				</div>
			</div>
			<div class="outlet-action-col" style="display: flex; gap: 0.6rem; flex-wrap: wrap; align-items: center;">
				<a 
					href="/restaurant/%d" 
					target="_blank" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #1e293b; text-decoration: none;"
				>
					<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
					View Storefront
				</a>
				<button 
					type="button" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #1e293b;"
					hx-get="/merchant/restaurant/edit?restaurant_id=%d" 
					hx-target="#merchant-modal"
				>
					<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
					Edit Outlet Info
				</button>
				<button 
					type="button" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #1e293b;"
					hx-get="/merchant/hours?restaurant_id=%d" 
					hx-target="#merchant-modal"
				>
					<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
					Manage Operating Hours
				</button>
				<button 
					type="button" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #1e293b;"
					hx-get="/merchant/delivery?restaurant_id=%d" 
					hx-target="#merchant-modal"
				>
					<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="3" width="15" height="13"></rect><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"></polygon><circle cx="5.5" cy="18.5" r="2.5"></circle><circle cx="18.5" cy="18.5" r="2.5"></circle></svg>
					Delivery Settings
				</button>
				<button 
					type="button" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #087443;"
					hx-get="/merchant/restaurant/claim" 
					hx-target="#merchant-modal"
				>
					<span>📍</span>
					Claim Eatery
				</button>
				<button 
					type="button" 
					class="btn-secondary" 
					style="display: inline-flex; align-items: center; gap: 0.45rem; background: #ffffff; border: 1.5px solid #cbd5e1; padding: 0.65rem 1.15rem; border-radius: 8px; font-weight: 700; cursor: pointer; color: #087443;"
					hx-get="/merchant/restaurant/new" 
					hx-target="#merchant-modal"
				>
					<span>+</span>
					Add Outlet
				</button>
				<button 
					type="button" 
					class="btn-primary" 
					hx-get="/merchant/menu/new?restaurant_id=%d" 
					hx-target="#merchant-modal"
				>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
					Add New Dish
				</button>
			</div>
		</div>

		<div class="metrics-grid">
			<div class="metric-card">
				<div class="metric-label">Total Menu Items</div>
				<div class="metric-val">%d</div>
			</div>
			<div class="metric-card">
				<div class="metric-label">Available for Ordering</div>
				<div class="metric-val" style="color: #15803d;">%d</div>
			</div>
			<div class="metric-card">
				<div class="metric-label">Unavailable / Sold Out</div>
				<div class="metric-val" style="color: #64748b;">%d</div>
			</div>
			<div class="metric-card highlight">
				<div class="metric-label">Items Priced ≤ ₹100</div>
				<div class="metric-val" style="color: #e23744;">%d</div>
				<div class="metric-sub">Affordable Discovery Tier</div>
			</div>
		</div>

		<div class="menu-section-card">
			<div class="menu-section-header">
				<div>
					<h2 class="menu-section-title">Menu & Food Items Management</h2>
					<p class="menu-section-subtitle">Manage item pricing, dietary classifications, and instant live availability toggles.</p>
				</div>
				<button 
					type="button" 
					class="btn-primary-sm" 
					hx-get="/merchant/menu/new" 
					hx-target="#merchant-modal"
				>
					+ Add Dish
				</button>
			</div>

			<div class="table-responsive">
				<table class="merchant-table">
					<thead>
						<tr>
							<th style="width: 28%%;">Dish Name & Description</th>
							<th style="width: 16%%;">Category</th>
							<th style="width: 14%%;">Diet Type</th>
							<th style="width: 14%%;">Price</th>
							<th style="width: 14%%;">Status</th>
							<th style="width: 14%%; text-align: right;">Actions</th>
						</tr>
					</thead>
					<tbody id="merchant-menu-body">
						%s
					</tbody>
				</table>
			</div>
		</div>

		<div class="menu-section-card" style="margin-top: 2rem;">
			<div class="menu-section-header">
				<div>
					<h2 class="menu-section-title">Operations & Delivery Configuration</h2>
					<p class="menu-section-subtitle">Real-time delivery serviceability and dispatch parameters for this outlet.</p>
				</div>
				<button 
					type="button" 
					class="btn-secondary" 
					style="font-size: 0.85rem; padding: 0.45rem 0.9rem; border-radius: 8px; font-weight: 700; background: #ffffff; border: 1.5px solid #cbd5e1; cursor: pointer; color: #1e293b; display: inline-flex; align-items: center; gap: 0.35rem;"
					hx-get="/merchant/delivery?restaurant_id=%d" 
					hx-target="#merchant-modal"
				>
					⚙️ Configure Delivery
				</button>
			</div>

			<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; margin-top: 1rem;">
				<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1.25rem;">
					<div style="font-size: 0.78rem; font-weight: 800; color: #64748b; text-transform: uppercase;">Delivery Service</div>
					<div style="margin-top: 0.5rem; font-size: 1.05rem; font-weight: 800;">%s</div>
					<div style="font-size: 0.82rem; color: #64748b; margin-top: 0.35rem;">%s</div>
				</div>
				<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1.25rem;">
					<div style="font-size: 0.78rem; font-weight: 800; color: #64748b; text-transform: uppercase;">Takeaway / Pickup</div>
					<div style="margin-top: 0.5rem; font-size: 1.05rem; font-weight: 800;">%s</div>
					<div style="font-size: 0.82rem; color: #64748b; margin-top: 0.35rem;">Customer counter takeaway</div>
				</div>
				<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1.25rem;">
					<div style="font-size: 0.78rem; font-weight: 800; color: #64748b; text-transform: uppercase;">Service Radius & Fee</div>
					<div style="margin-top: 0.5rem; font-size: 1.05rem; font-weight: 800;">%.1f km &bull; &#8377;%.0f</div>
					<div style="font-size: 0.82rem; color: #64748b; margin-top: 0.35rem;">Configured delivery radius & fee</div>
				</div>
				<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1.25rem;">
					<div style="font-size: 0.78rem; font-weight: 800; color: #64748b; text-transform: uppercase;">Platform Fulfillment State</div>
					<div style="margin-top: 0.5rem; font-size: 0.95rem; font-weight: 800; color: #0284c7;">Discovery Stage Mode</div>
					<div style="font-size: 0.82rem; color: #64748b; margin-top: 0.35rem;">Direct merchant dispatch; no automatic rider assignment</div>
				</div>
			</div>
		</div>`, outletSwitcherHTML, qualificationBanner, statusBadge, escapeHTML(rest.Name), escapeHTML(rest.Description), escapeHTML(rest.Address), escapeHTML(rest.Phone), escapeHTML(rest.OpeningHours), rest.ID, rest.ID, rest.ID, rest.ID, rest.ID, totalFoods, availCount, unavailCount, under100Count, rowsHTML, rest.ID, delivStateHTML, delivStateSub, pickupStateHTML, delivSettings.DeliveryRadiusKm, delivSettings.DeliveryFee)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Merchant Portal — Under100 Food Discovery</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <style>
    :root {
      --primary-red: #e23744;
      --primary-red-hover: #c82333;
      --primary-green: #087443;
      --text-dark: #0f172a;
      --text-muted: #64748b;
      --bg-light: #f8fafc;
      --border-gray: #e2e8f0;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-light);
      color: var(--text-dark);
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      padding-bottom: 4rem;
    }
    .merchant-nav {
      background: #ffffff;
      border-bottom: 1px solid var(--border-gray);
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      box-shadow: 0 1px 3px rgba(0,0,0,0.03);
    }
    .merchant-brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }
    .merchant-logo-pill {
      background: #e23744;
      color: white;
      font-family: 'Outfit', sans-serif;
      font-weight: 900;
      font-size: 0.95rem;
      padding: 0.35rem 0.65rem;
      border-radius: 8px;
      letter-spacing: -0.5px;
    }
    .merchant-portal-title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.25rem;
      font-weight: 800;
      color: var(--text-dark);
    }
    .merchant-nav-right {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    .merchant-user-badge {
      background: #f1f5f9;
      border: 1px solid var(--border-gray);
      padding: 0.35rem 0.8rem;
      border-radius: 9999px;
      font-size: 0.82rem;
      font-weight: 700;
      color: var(--text-dark);
      display: flex;
      align-items: center;
      gap: 0.4rem;
    }
    .container {
      max-width: 1200px;
      margin: 2rem auto;
      padding: 0 1.5rem;
    }
    .merchant-alert {
      padding: 0.85rem 1.25rem;
      border-radius: 10px;
      margin-bottom: 1.5rem;
      font-size: 0.92rem;
      font-weight: 600;
    }
    .merchant-alert.success { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }
    .merchant-alert.error { background: #fee2e2; color: #991b1b; border: 1px solid #fecaca; }

    .qualify-banner {
      display: flex;
      align-items: flex-start;
      gap: 1rem;
      padding: 1.25rem 1.5rem;
      border-radius: 14px;
      margin-bottom: 1.5rem;
    }
    .qualify-banner.qualified { background: #f0fdf4; border: 1px solid #bbf7d0; color: #166534; }
    .qualify-banner.unqualified { background: #fffbeb; border: 1px solid #fde68a; color: #92400e; }
    .qualify-banner p { font-size: 0.88rem; margin-top: 0.2rem; opacity: 0.9; }

    .outlet-header-card {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      border-radius: 16px;
      padding: 1.75rem 2rem;
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1.5rem;
      margin-bottom: 1.5rem;
      box-shadow: 0 2px 6px rgba(0,0,0,0.02);
    }
    .outlet-meta-top {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-bottom: 0.5rem;
    }
    .outlet-tag {
      background: #f1f5f9;
      color: #475569;
      font-size: 0.72rem;
      font-weight: 800;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      padding: 0.2rem 0.55rem;
      border-radius: 6px;
    }
    .badge-status {
      font-size: 0.75rem;
      font-weight: 800;
      text-transform: uppercase;
      padding: 0.2rem 0.6rem;
      border-radius: 6px;
    }
    .badge-status.approved { background: #dcfce7; color: #166534; }
    .badge-status.pending { background: #fef08a; color: #854d0e; }
    .badge-status.rejected { background: #fee2e2; color: #991b1b; }
    .outlet-name {
      font-family: 'Outfit', sans-serif;
      font-size: 1.75rem;
      font-weight: 800;
      color: var(--text-dark);
      margin-bottom: 0.35rem;
    }
    .outlet-desc {
      color: var(--text-muted);
      font-size: 0.95rem;
      max-width: 650px;
      line-height: 1.5;
      margin-bottom: 0.85rem;
    }
    .outlet-contact-row {
      display: flex;
      flex-wrap: wrap;
      gap: 1.25rem;
      font-size: 0.85rem;
      color: #475569;
    }

    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 1rem;
      margin-bottom: 2rem;
    }
    .metric-card {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      border-radius: 14px;
      padding: 1.35rem;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .metric-card.highlight {
      background: linear-gradient(135deg, #fff5f5 0%%, #ffffff 100%%);
      border-color: #fecaca;
    }
    .metric-label {
      font-size: 0.82rem;
      font-weight: 700;
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.4px;
    }
    .metric-val {
      font-family: 'Outfit', sans-serif;
      font-size: 2.2rem;
      font-weight: 800;
      margin-top: 0.35rem;
      color: var(--text-dark);
    }
    .metric-sub {
      font-size: 0.78rem;
      font-weight: 700;
      color: #e23744;
      margin-top: 0.25rem;
    }

    .menu-section-card {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      border-radius: 16px;
      padding: 1.75rem;
      box-shadow: 0 2px 6px rgba(0,0,0,0.02);
    }
    .menu-section-header {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.5rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid var(--border-gray);
    }
    .menu-section-title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.35rem;
      font-weight: 800;
      color: var(--text-dark);
    }
    .menu-section-subtitle {
      font-size: 0.88rem;
      color: var(--text-muted);
      margin-top: 0.2rem;
    }

    .table-responsive {
      width: 100%%;
      overflow-x: auto;
    }
    .merchant-table {
      width: 100%%;
      border-collapse: collapse;
      font-size: 0.9rem;
    }
    .merchant-table th {
      text-align: left;
      padding: 0.85rem 1rem;
      background: #f8fafc;
      color: #64748b;
      font-size: 0.78rem;
      font-weight: 800;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      border-bottom: 1px solid var(--border-gray);
    }
    .merchant-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      vertical-align: middle;
    }
    .merchant-table tr:hover td {
      background: #fbfcfe;
    }
    .food-name-cell {
      font-weight: 700;
      color: var(--text-dark);
    }
    .food-desc-cell {
      font-size: 0.8rem;
      color: var(--text-muted);
      margin-top: 0.2rem;
    }
    .diet-badge {
      display: inline-block;
      font-size: 0.7rem;
      font-weight: 800;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      text-transform: uppercase;
    }
    .diet-badge.veg { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }
    .diet-badge.nonveg { background: #fee2e2; color: #991b1b; border: 1px solid #fecaca; }

    .price-text {
      font-family: 'Outfit', sans-serif;
      font-weight: 800;
      font-size: 1.05rem;
      color: var(--text-dark);
    }
    .price-badge-under100 {
      display: inline-block;
      background: #f0fdf4;
      color: #166534;
      border: 1px solid #bbf7d0;
      font-size: 0.7rem;
      font-weight: 800;
      padding: 0.15rem 0.45rem;
      border-radius: 4px;
      margin-left: 0.35rem;
    }

    .btn-toggle-avail {
      background: #f1f5f9;
      border: 1px solid var(--border-gray);
      border-radius: 9999px;
      padding: 0.35rem 0.75rem;
      font-size: 0.78rem;
      font-weight: 700;
      cursor: pointer;
      transition: all 0.15s ease;
    }
    .btn-toggle-avail.avail {
      background: #dcfce7;
      color: #166534;
      border-color: #bbf7d0;
    }
    .btn-toggle-avail.soldout {
      background: #f1f5f9;
      color: #64748b;
      border-color: #cbd5e1;
    }

    .btn-action-sm {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      color: var(--text-dark);
      padding: 0.4rem 0.75rem;
      border-radius: 8px;
      font-size: 0.8rem;
      font-weight: 700;
      cursor: pointer;
      text-decoration: none;
      display: inline-block;
      transition: background 0.15s ease;
    }
    .btn-action-sm:hover { background: #f1f5f9; }
    .btn-action-sm.delete {
      color: #e23744;
      border-color: #fecaca;
    }
    .btn-action-sm.delete:hover {
      background: #fee2e2;
    }

    .btn-primary {
      background: var(--primary-red);
      color: white;
      border: none;
      border-radius: 10px;
      padding: 0.7rem 1.4rem;
      font-weight: 700;
      font-size: 0.95rem;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
      text-decoration: none;
      transition: background 0.15s ease;
    }
    .btn-primary:hover { background: var(--primary-red-hover); }

    .btn-primary-sm {
      background: var(--primary-red);
      color: white;
      border: none;
      border-radius: 8px;
      padding: 0.45rem 0.9rem;
      font-weight: 700;
      font-size: 0.85rem;
      cursor: pointer;
      text-decoration: none;
    }
    .btn-secondary {
      background: #f1f5f9;
      border: 1px solid var(--border-gray);
      color: var(--text-dark);
      border-radius: 8px;
      padding: 0.6rem 1.2rem;
      font-weight: 700;
      font-size: 0.9rem;
      text-decoration: none;
      display: inline-block;
    }

    .empty-outlet-card {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      border-radius: 16px;
      padding: 4rem 2rem;
      text-align: center;
      max-width: 600px;
      margin: 2rem auto;
      box-shadow: 0 4px 12px rgba(0,0,0,0.03);
    }
    .empty-icon-box {
      width: 72px;
      height: 72px;
      background: #fee2e2;
      border-radius: 50%%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem;
    }
    .empty-outlet-card h2 {
      font-family: 'Outfit', sans-serif;
      font-size: 1.5rem;
      font-weight: 800;
      margin-bottom: 0.5rem;
    }
    .empty-outlet-card p {
      color: var(--text-muted);
      font-size: 0.95rem;
      line-height: 1.6;
      margin-bottom: 1.25rem;
    }
    .empty-outlet-help {
      background: #f8fafc;
      border: 1px dashed var(--border-gray);
      padding: 1rem;
      border-radius: 10px;
      font-size: 0.85rem;
      color: #475569;
    }

    /* Modal styling */
    .merchant-modal-overlay {
      position: fixed;
      top: 0; left: 0; right: 0; bottom: 0;
      background: rgba(15, 23, 42, 0.6);
      backdrop-filter: blur(2px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9999;
    }
    .merchant-modal-card {
      background: #ffffff;
      border-radius: 18px;
      width: 540px;
      max-width: 95%%;
      max-height: 90vh;
      overflow-y: auto;
      padding: 2rem;
      box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    }
    .form-group {
      margin-bottom: 1.25rem;
    }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: var(--text-dark);
      margin-bottom: 0.4rem;
      text-transform: uppercase;
      letter-spacing: 0.3px;
    }
    .form-input {
      width: 100%%;
      border: 1px solid var(--border-gray);
      background: #ffffff;
      padding: 0.65rem 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-family: inherit;
      color: var(--text-dark);
    }
    .form-input:focus {
      outline: none;
      border-color: var(--primary-red);
      box-shadow: 0 0 0 3px rgba(226, 55, 68, 0.1);
    }
    .form-helper {
      font-size: 0.78rem;
      color: var(--text-muted);
      margin-top: 0.3rem;
    }
  </style>
</head>
<body>
  <nav class="merchant-nav">
    <div class="merchant-brand">
      <div class="merchant-logo-pill">Bite100</div>
      <div class="merchant-portal-title">Merchant Owner Portal</div>
    </div>
    <div class="merchant-nav-right">
      <a href="/" style="color: var(--text-dark); text-decoration: none; font-size: 0.9rem; font-weight: 700;">&larr; Customer Storefront</a>
      <div class="merchant-user-badge">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
        %s (Owner)
      </div>
      <a href="/logout" style="color: #e23744; text-decoration: none; font-size: 0.85rem; font-weight: 700; padding: 0.35rem 0.75rem; border: 1px solid #fecaca; border-radius: 8px;">Log Out</a>
    </div>
  </nav>

  <main class="container">
    %s
    %s
  </main>

  <div id="merchant-modal"></div>

  <script>
    function closeMerchantModal() {
      var m = document.getElementById('merchant-modal');
      if (m) m.innerHTML = '';
    }
  </script>
</body>
</html>`, escapeHTML(userName), flashHTML, mainContent)
}

// RenderMerchantSingleFoodRow renders an individual food item table row
func RenderMerchantSingleFoodRow(f *ent.FoodItem, itemRule ...foodrepo.FoodItemDeliveryRule) string {
	isVeg := foodrepo.IsDishVeg(f.Name, f.Category, f.Description)
	dietBadge := `<span class="diet-badge veg">Pure Veg</span>`
	if !isVeg {
		dietBadge = `<span class="diet-badge nonveg">Non-Veg</span>`
	}

	under100Tag := ""
	if f.Price <= 100.0 {
		under100Tag = `<span class="price-badge-under100">≤ ₹100</span>`
	}

	availClass := "avail"
	availText := "Available"
	if !f.Available {
		availClass = "soldout"
		availText = "Sold Out"
	}

	desc := f.Description
	if desc == "" {
		desc = "No description provided"
	}

	delivBadge := `<span style="display: inline-flex; align-items: center; gap: 3px; font-size: 0.72rem; font-weight: 700; background: #dcfce7; color: #166534; padding: 2px 7px; border-radius: 4px; border: 1px solid #bbf7d0;">🚚 Deliverable</span>`
	if len(itemRule) > 0 {
		r := itemRule[0]
		if !r.IsDeliverable {
			delivBadge = `<span style="display: inline-flex; align-items: center; gap: 3px; font-size: 0.72rem; font-weight: 700; background: #f1f5f9; color: #475569; padding: 2px 7px; border-radius: 4px; border: 1px solid #cbd5e1;">🍽️ Takeaway Only</span>`
		} else if r.DeliveryRadiusLimit > 0 {
			delivBadge = fmt.Sprintf(`<span style="display: inline-flex; align-items: center; gap: 3px; font-size: 0.72rem; font-weight: 700; background: #fef3c7; color: #92400e; padding: 2px 7px; border-radius: 4px; border: 1px solid #fde68a;">🛵 ≤%.0fkm Range</span>`, r.DeliveryRadiusLimit)
		}
	}

	return fmt.Sprintf(`
	<tr id="food-row-%d">
		<td>
			<div class="food-name-cell" style="display: flex; align-items: center; gap: 8px;">
				<span style="font-weight: 700;">%s</span>
				%s
			</div>
			<div class="food-desc-cell">%s</div>
		</td>
		<td><span style="font-size: 0.85rem; color: #475569; font-weight: 600;">%s</span></td>
		<td>%s</td>
		<td>
			<span class="price-text">₹%.0f</span>
			%s
		</td>
		<td>
			<button 
				type="button" 
				class="btn-toggle-avail %s"
				hx-post="/merchant/menu/%d/availability"
				hx-target="#food-row-%d"
				hx-swap="outerHTML"
				title="Click to toggle Available vs Sold Out"
			>
				%s
			</button>
		</td>
		<td style="text-align: right;">
			<div style="display: flex; gap: 0.4rem; justify-content: flex-end;">
				<button 
					type="button" 
					class="btn-action-sm" 
					hx-get="/merchant/menu/%d/edit" 
					hx-target="#merchant-modal"
				>
					Edit
				</button>
				<button 
					type="button" 
					class="btn-action-sm delete" 
					hx-post="/merchant/menu/%d/delete" 
					hx-target="#food-row-%d"
					hx-swap="outerHTML"
					hx-confirm="Are you sure you want to delete this dish?"
				>
					Delete
				</button>
			</div>
		</td>
	</tr>`, f.ID, escapeHTML(f.Name), delivBadge, escapeHTML(desc), escapeHTML(f.Category), dietBadge, f.Price, under100Tag, availClass, f.ID, f.ID, availText, f.ID, f.ID, f.ID)
}

// RenderMerchantFoodModal renders the Add/Edit dish modal
func RenderMerchantFoodModal(rest *ent.Restaurant, food *ent.FoodItem, actionURL string, isEdit bool, itemRule ...foodrepo.FoodItemDeliveryRule) string {
	modalTitle := "Add New Dish to Menu"
	btnText := "Add Dish"
	nameVal := ""
	priceVal := ""
	catVal := "Breakfast"
	descVal := ""
	isVegChecked := "checked"
	isNonVegChecked := ""
	availChecked := "checked"
	delivChecked := "checked"
	radiusLimitVal := "0"
	radiusBoxDisplay := "block"

	if isEdit && food != nil {
		modalTitle = fmt.Sprintf("Edit Dish: %s", escapeHTML(food.Name))
		btnText = "Save Changes"
		nameVal = escapeHTML(food.Name)
		priceVal = fmt.Sprintf("%.0f", food.Price)
		catVal = escapeHTML(food.Category)
		descVal = escapeHTML(food.Description)
		if !foodrepo.IsDishVeg(food.Name, food.Category, food.Description) {
			isVegChecked = ""
			isNonVegChecked = "checked"
		}
		if !food.Available {
			availChecked = ""
		}
	}

	if len(itemRule) > 0 {
		r := itemRule[0]
		if !r.IsDeliverable {
			delivChecked = ""
			radiusBoxDisplay = "none"
		}
		if r.DeliveryRadiusLimit > 0 {
			radiusLimitVal = fmt.Sprintf("%.0f", r.DeliveryRadiusLimit)
		}
	}

	return fmt.Sprintf(`
	<div class="merchant-modal-overlay" onclick="if(event.target === this) closeMerchantModal()">
		<div class="merchant-modal-card">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">%s</h3>
				<button type="button" onclick="closeMerchantModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			<form 
				method="POST" 
				action="%s" 
				hx-post="%s" 
				hx-target="#merchant-modal"
			>
				<input type="hidden" name="restaurant_id" value="%d">
				<div style="margin-bottom: 1.25rem; background: #f8fafc; border: 1px dashed #cbd5e1; border-radius: 10px; padding: 0.75rem;">
					<div style="font-size: 0.75rem; font-weight: 800; color: #64748b; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.5rem;">⚡ Quick-Fill Popular Dishes (&le; &#8377;100):</div>
					<div style="display: flex; gap: 0.4rem; flex-wrap: wrap;">
						<button type="button" onclick="autofillDish('Crispy Benne Masala Dosa', 80, 'South Indian', 'veg', 'Classic butter roast dosa with potato palya and coconut chutney')" style="background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 0.35rem 0.6rem; font-size: 0.78rem; font-weight: 700; cursor: pointer; color: #1e293b;">🥞 Benne Dosa (&#8377;80)</button>
						<button type="button" onclick="autofillDish('Degree Filter Coffee', 15, 'Beverages', 'veg', 'Traditional South Indian chicory filter brew')" style="background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 0.35rem 0.6rem; font-size: 0.78rem; font-weight: 700; cursor: pointer; color: #1e293b;">☕ Filter Coffee (&#8377;15)</button>
						<button type="button" onclick="autofillDish('Steamed Idli Vada Combo', 45, 'Breakfast', 'veg', 'Two hot steamed idlis and one crispy medu vada with sambar')" style="background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 0.35rem 0.6rem; font-size: 0.78rem; font-weight: 700; cursor: pointer; color: #1e293b;">🥣 Idli Vada (&#8377;45)</button>
						<button type="button" onclick="autofillDish('Pav Bhaji with Extra Butter', 95, 'Street Food', 'veg', 'Spicy mashed vegetable curry with 2 butter-toasted pavs')" style="background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 0.35rem 0.6rem; font-size: 0.78rem; font-weight: 700; cursor: pointer; color: #1e293b;">🥘 Pav Bhaji (&#8377;95)</button>
						<button type="button" onclick="autofillDish('Egg Kathi Roll', 75, 'Rolls & Wraps', 'nonveg', 'Fresh paratha wrap rolled with spiced scrambled eggs and onions')" style="background: #ffffff; border: 1px solid #cbd5e1; border-radius: 6px; padding: 0.35rem 0.6rem; font-size: 0.78rem; font-weight: 700; cursor: pointer; color: #1e293b;">🌯 Egg Roll (&#8377;75)</button>
					</div>
				</div>

				<div class="form-group">
					<label class="form-label">Dish Name *</label>
					<input type="text" id="dishNameInput" name="name" class="form-input" value="%s" placeholder="e.g. Benne Masala Dosa" required>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
					<div class="form-group">
						<label class="form-label">Price (₹ INR) *</label>
						<input type="number" step="1" min="1" id="dishPriceInput" name="price" class="form-input" value="%s" placeholder="e.g. 75" required>
						<div class="form-helper">Dishes ≤ ₹100 qualify for budget discovery.</div>
					</div>

					<div class="form-group">
						<label class="form-label">Category</label>
						<select id="dishCategorySelect" name="category" class="form-input">
							<option value="South Indian" %s>South Indian</option>
							<option value="North Indian" %s>North Indian</option>
							<option value="Chinese" %s>Chinese</option>
							<option value="Street Food" %s>Street Food</option>
							<option value="Snacks" %s>Snacks</option>
							<option value="Breakfast" %s>Breakfast</option>
							<option value="Lunch" %s>Lunch</option>
							<option value="Desserts" %s>Desserts</option>
							<option value="Beverages" %s>Beverages</option>
							<option value="Fast Food" %s>Fast Food</option>
							<option value="Bakery" %s>Bakery</option>
							<option value="Healthy Food" %s>Healthy Food</option>
							<option value="Snacks & Chaat" %s>Snacks & Chaat</option>
							<option value="Rolls & Wraps" %s>Rolls & Wraps</option>
							<option value="Combos" %s>Combos</option>
						</select>
					</div>
				</div>

				<div class="form-group">
					<label class="form-label">Dietary Classification</label>
					<div style="display: flex; gap: 1.5rem; margin-top: 0.35rem;">
						<label style="display: flex; align-items: center; gap: 0.4rem; font-size: 0.9rem; font-weight: 600; cursor: pointer;">
							<input type="radio" name="diet_type" value="veg" %s>
							<span class="diet-badge veg">Pure Veg</span>
						</label>
						<label style="display: flex; align-items: center; gap: 0.4rem; font-size: 0.9rem; font-weight: 600; cursor: pointer;">
							<input type="radio" name="diet_type" value="nonveg" %s>
							<span class="diet-badge nonveg">Non-Veg</span>
						</label>
					</div>
				</div>

				<div class="form-group">
					<label class="form-label">Description (Optional)</label>
					<textarea id="dishDescInput" name="description" class="form-input" rows="2" placeholder="Briefly describe the ingredients or special preparation...">%s</textarea>
				</div>

				<div class="form-group">
					<label style="display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; font-weight: 700; cursor: pointer;">
						<input type="checkbox" name="available" value="true" %s>
						<span>Item is Currently Available in Kitchen (In-Stock)</span>
					</label>
				</div>

				<!-- Item Delivery Availability & Location Range Restriction -->
				<div class="form-group" style="background: #f8fafc; border: 1.5px solid #e2e8f0; border-radius: 12px; padding: 1rem; margin-top: 1rem;">
					<div style="font-size: 0.8rem; font-weight: 800; color: #0f172a; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 0.5rem;">
						🚚 Item Delivery &amp; Area Availability
					</div>
					<label style="display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; font-weight: 700; cursor: pointer;">
						<input type="checkbox" name="is_deliverable" value="true" %s onchange="document.getElementById('radiusLimitBox').style.display = this.checked ? 'block' : 'none'">
						<span>Item is Deliverable (Uncheck for Dine-in / Takeaway Only)</span>
					</label>

					<div id="radiusLimitBox" style="margin-top: 0.75rem; padding-top: 0.75rem; border-top: 1px dashed #cbd5e1; display: %s;">
						<label class="form-label" style="font-size: 0.82rem; margin-bottom: 0.25rem;">Deliverable Area / Location Constraint</label>
						<select name="delivery_radius_limit" class="form-input" style="font-size: 0.88rem;">
							<option value="0" %s>Full Delivery Zone (Standard Restaurant Radius)</option>
							<option value="3" %s>Short-Range Only (&le; 3 km) — For crispy dosas, hot tea/coffee, ice creams</option>
							<option value="5" %s>Medium Range (&le; 5 km)</option>
						</select>
						<div class="form-helper" style="font-size: 0.76rem; color: #64748b; margin-top: 0.35rem;">
							Distant customers beyond this radius will see this item marked as Takeaway Only.
						</div>
					</div>
				</div>

				<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.75rem; padding-top: 1rem; border-top: 1px solid #e2e8f0;">
					<button type="button" onclick="closeMerchantModal()" class="btn-secondary">Cancel</button>
					<button type="submit" class="btn-primary">%s</button>
				</div>
			</form>
			<script>
				function autofillDish(name, price, cat, diet, desc) {
					var n = document.getElementById('dishNameInput'); if(n) n.value = name;
					var p = document.getElementById('dishPriceInput'); if(p) p.value = price;
					var d = document.getElementById('dishDescInput'); if(d) d.value = desc;
					var c = document.getElementById('dishCategorySelect'); if(c) c.value = cat;
					var r = document.querySelector('input[name="diet_type"][value="' + diet + '"]'); if(r) r.checked = true;
				}
			</script>
		</div>
	</div>`, modalTitle, actionURL, actionURL, rest.ID, nameVal, priceVal,
		selectedAttr(catVal, "South Indian"),
		selectedAttr(catVal, "North Indian"),
		selectedAttr(catVal, "Chinese"),
		selectedAttr(catVal, "Street Food"),
		selectedAttr(catVal, "Snacks"),
		selectedAttr(catVal, "Breakfast"),
		selectedAttr(catVal, "Lunch"),
		selectedAttr(catVal, "Desserts"),
		selectedAttr(catVal, "Beverages"),
		selectedAttr(catVal, "Fast Food"),
		selectedAttr(catVal, "Bakery"),
		selectedAttr(catVal, "Healthy Food"),
		selectedAttr(catVal, "Snacks & Chaat"),
		selectedAttr(catVal, "Rolls & Wraps"),
		selectedAttr(catVal, "Combos"),
		isVegChecked, isNonVegChecked, descVal, availChecked,
		delivChecked, radiusBoxDisplay,
		selectedAttr(radiusLimitVal, "0"),
		selectedAttr(radiusLimitVal, "3"),
		selectedAttr(radiusLimitVal, "5"),
		btnText)
}

func selectedAttr(current, target string) string {
	if strings.EqualFold(current, target) {
		return "selected"
	}
	return ""
}

// RenderMerchantRestaurantModal renders modal form for self-service restaurant submission
func RenderMerchantRestaurantModal(errMsg string, input foodrepo.RestaurantInput) string {
	errBanner := ""
	if errMsg != "" {
		errBanner = fmt.Sprintf(`<div class="merchant-alert error" style="margin-bottom: 1.25rem;">%s</div>`, escapeHTML(errMsg))
	}

	latVal := "12.9716"
	if input.Latitude != 0 {
		latVal = fmt.Sprintf("%.6f", input.Latitude)
	}
	lngVal := "77.5946"
	if input.Longitude != 0 {
		lngVal = fmt.Sprintf("%.6f", input.Longitude)
	}

	return fmt.Sprintf(`
	<div class="merchant-modal-overlay" onclick="if(event.target === this) closeMerchantModal()">
		<div class="merchant-modal-card" style="max-width: 650px;">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<div>
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; color: #0f172a;">Register My Restaurant</h3>
					<p style="font-size: 0.85rem; color: #64748b; margin-top: 0.2rem;">Submit your eatery for Under ₹100 platform review and verification.</p>
				</div>
				<button type="button" onclick="closeMerchantModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			%s

			<form 
				method="POST" 
				action="/merchant/restaurant/submit" 
				hx-post="/merchant/restaurant/submit" 
				hx-target="#merchant-modal"
			>
				<div class="form-group">
					<label class="form-label">Restaurant Name *</label>
					<input type="text" name="name" class="form-input" value="%s" placeholder="e.g. Nagarjuna Andhra Meals" required autofocus>
				</div>

				<div class="form-group">
					<label class="form-label">Street Address & Locality *</label>
					<input type="text" name="address" class="form-input" value="%s" placeholder="e.g. 44/1, Residency Road, Shanthala Nagar, Ashok Nagar, Bengaluru" required>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
					<div class="form-group">
						<label class="form-label">Contact Phone</label>
						<input type="text" name="phone" class="form-input" value="%s" placeholder="e.g. +91 80 2555 1234">
					</div>
					<div class="form-group">
						<label class="form-label">Opening Hours</label>
						<input type="text" name="opening_hours" class="form-input" value="%s" placeholder="e.g. 7:00 AM – 10:30 PM">
					</div>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
					<div class="form-group">
						<label class="form-label">Latitude</label>
						<input type="number" step="0.000001" name="latitude" class="form-input" value="%s">
						<div class="form-helper">Defaulted to Bangalore center (12.9716)</div>
					</div>
					<div class="form-group">
						<label class="form-label">Longitude</label>
						<input type="number" step="0.000001" name="longitude" class="form-input" value="%s">
						<div class="form-helper">Defaulted to Bangalore center (77.5946)</div>
					</div>
				</div>

				<div class="form-group">
					<label class="form-label">Website or Online Menu (Optional)</label>
					<input type="url" name="website" class="form-input" value="%s" placeholder="https://example.com">
				</div>

				<div class="form-group">
					<label class="form-label">Description (Optional)</label>
					<textarea name="description" class="form-input" rows="2" placeholder="Briefly describe your cuisine, signature dishes, or dining style...">%s</textarea>
				</div>

				<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid #e2e8f0;">
					<button type="button" onclick="closeMerchantModal()" class="btn-secondary">Cancel</button>
					<button type="submit" class="btn-primary" style="display: inline-flex; align-items: center; gap: 0.4rem;">
						<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
						Submit for Review
					</button>
				</div>
			</form>
		</div>
	</div>`, errBanner, escapeHTML(input.Name), escapeHTML(input.Address), escapeHTML(input.Phone), escapeHTML(input.OpeningHours), latVal, lngVal, escapeHTML(input.Website), escapeHTML(input.Description))
}

// RenderMerchantEditRestaurantModal renders the outlet info editor modal
func RenderMerchantEditRestaurantModal(rest *ent.Restaurant, errMsg string) string {
	errHTML := ""
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="merchant-alert error" style="margin-bottom: 1rem;">%s</div>`, escapeHTML(errMsg))
	}

	websiteVal := ""
	if rest.Website != "" {
		websiteVal = escapeHTML(rest.Website)
	}

	return fmt.Sprintf(`
	<div class="merchant-modal-overlay" onclick="if(event.target === this) closeMerchantModal()">
		<div class="merchant-modal-card">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Edit Outlet Information</h3>
				<button type="button" onclick="closeMerchantModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			%s

			<form method="POST" action="/merchant/restaurant/edit" hx-post="/merchant/restaurant/edit" hx-target="#merchant-modal">
				<input type="hidden" name="restaurant_id" value="%d">

				<div class="form-group">
					<label class="form-label">Restaurant Name *</label>
					<input type="text" name="name" class="form-input" value="%s" required>
				</div>

				<div class="form-group">
					<label class="form-label">Street Address & Locality *</label>
					<input type="text" name="address" class="form-input" value="%s" required>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
					<div class="form-group">
						<label class="form-label">Contact Phone</label>
						<input type="text" name="phone" class="form-input" value="%s" placeholder="+91 80 2331 7536">
					</div>
					<div class="form-group">
						<label class="form-label">Website / Menu Link</label>
						<input type="text" name="website" class="form-input" value="%s" placeholder="https://example.com">
					</div>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1rem;">
					<div class="form-group">
						<label class="form-label">Operating Hours Summary</label>
						<input type="text" name="opening_hours" class="form-input" value="%s" placeholder="7:30 AM - 10:00 PM">
					</div>
					<div class="form-group">
						<label class="form-label">Latitude</label>
						<input type="number" step="0.0001" name="latitude" class="form-input" value="%.6f">
					</div>
					<div class="form-group">
						<label class="form-label">Longitude</label>
						<input type="number" step="0.0001" name="longitude" class="form-input" value="%.6f">
					</div>
				</div>

				<div class="form-group">
					<label class="form-label">Bio / Story</label>
					<textarea name="description" class="form-input" rows="3" placeholder="Tell customers about your kitchen, heritage, or specialties...">%s</textarea>
				</div>

				<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid #e2e8f0;">
					<button type="button" onclick="closeMerchantModal()" class="btn-secondary">Cancel</button>
					<button type="submit" class="btn-primary">Save Changes</button>
				</div>
			</form>
		</div>
	</div>`, errHTML, rest.ID, escapeHTML(rest.Name), escapeHTML(rest.Address), escapeHTML(rest.Phone), websiteVal, escapeHTML(rest.OpeningHours), rest.Latitude, rest.Longitude, escapeHTML(rest.Description))
}

// RenderMerchantDeliveryModal renders modal for configuring restaurant delivery and serviceability
func RenderMerchantDeliveryModal(rest *ent.Restaurant, settings foodrepo.DeliverySettings, successMsg, errMsg string) string {
	errHTML := ""
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="merchant-alert error" style="margin-bottom: 1rem;">%s</div>`, escapeHTML(errMsg))
	} else if successMsg != "" {
		errHTML = fmt.Sprintf(`<div class="merchant-alert" style="background: #f0fdf4; border: 1.5px solid #bbf7d0; color: #166534; margin-bottom: 1rem;">%s</div>`, escapeHTML(successMsg))
	}

	delivChecked := ""
	if settings.DeliveryEnabled {
		delivChecked = "checked"
	}
	pickupChecked := ""
	if settings.PickupEnabled {
		pickupChecked = "checked"
	}
	suspendChecked := ""
	if settings.IsDeliverySuspended {
		suspendChecked = "checked"
	}

	radiusVal := settings.DeliveryRadiusKm
	if radiusVal <= 0 {
		radiusVal = 5.0
	}

	estMinVal := settings.EstimatedDeliveryMin
	if estMinVal <= 0 {
		estMinVal = 25
	}

	return fmt.Sprintf(`
	<div class="merchant-modal-overlay" onclick="if(event.target === this) closeMerchantModal()">
		<div class="merchant-modal-card" style="max-width: 620px;">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<div>
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Delivery & Serviceability Settings</h3>
					<p style="font-size: 0.85rem; color: #64748b; margin-top: 0.2rem;">%s &bull; Configure delivery radius, fee, and temporary pauses.</p>
				</div>
				<button type="button" onclick="closeMerchantModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			%s

			<form method="POST" action="/merchant/delivery" hx-post="/merchant/delivery" hx-target="#merchant-modal">
				<input type="hidden" name="restaurant_id" value="%d">

				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-bottom: 1.25rem;">
					<div style="background: #f8fafc; border: 1.5px solid #e2e8f0; border-radius: 10px; padding: 1rem;">
						<label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: 800; font-size: 0.95rem; color: #0f172a;">
							<input type="checkbox" name="delivery_enabled" value="true" %s style="width: 18px; height: 18px; accent-color: #087443;">
							Offer Direct Delivery
						</label>
						<p style="font-size: 0.78rem; color: #64748b; margin-top: 4px; margin-left: 26px;">
							Allows customers within your delivery radius to see this eatery as serviceable.
						</p>
					</div>
					<div style="background: #f8fafc; border: 1.5px solid #e2e8f0; border-radius: 10px; padding: 1rem;">
						<label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: 800; font-size: 0.95rem; color: #0f172a;">
							<input type="checkbox" name="pickup_enabled" value="true" %s style="width: 18px; height: 18px; accent-color: #087443;">
							Offer Counter Takeaway
						</label>
						<p style="font-size: 0.78rem; color: #64748b; margin-top: 4px; margin-left: 26px;">
							Allows customers to pick up prepared food directly from your counter.
						</p>
					</div>
				</div>

				<div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1rem; margin-bottom: 1.25rem;">
					<div class="form-group">
						<label class="form-label">Delivery Radius (km) *</label>
						<input type="number" step="0.5" min="1" max="30" name="delivery_radius_km" class="form-input" value="%.1f" required>
						<span style="font-size: 0.72rem; color: #64748b;">Max service distance</span>
					</div>
					<div class="form-group">
						<label class="form-label">Delivery Fee (₹) *</label>
						<input type="number" step="1" min="0" max="200" name="delivery_fee" class="form-input" value="%.0f" required>
						<span style="font-size: 0.72rem; color: #64748b;">₹0 = Free delivery</span>
					</div>
					<div class="form-group">
						<label class="form-label">Est. Time (mins) *</label>
						<input type="number" step="5" min="10" max="120" name="estimated_delivery_min" class="form-input" value="%d" required>
						<span style="font-size: 0.72rem; color: #64748b;">Prep + transit</span>
					</div>
				</div>

				<div style="background: #fffbeb; border: 1.5px solid #fde68a; border-radius: 10px; padding: 1rem; margin-bottom: 1.25rem;">
					<div style="font-weight: 800; font-size: 0.95rem; color: #92400e; margin-bottom: 0.5rem; display: flex; align-items: center; gap: 6px;">
						<span>⚠️</span> Temporary Delivery Suspension
					</div>
					<label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: 700; font-size: 0.88rem; color: #78350f; margin-bottom: 0.6rem;">
						<input type="checkbox" name="is_delivery_suspended" value="true" %s style="width: 16px; height: 16px; accent-color: #d97706;">
						Temporarily pause incoming delivery orders (Kitchen rush / heavy rain)
					</label>
					<div class="form-group" style="margin-bottom: 0;">
						<label class="form-label" style="color: #78350f;">Reason displayed to customers</label>
						<input type="text" name="suspension_reason" class="form-input" value="%s" placeholder="e.g. Delivery temporarily paused due to torrential rain">
					</div>
				</div>

				<div style="display: flex; justify-content: flex-end; gap: 0.75rem; padding-top: 1rem; border-top: 1px solid #e2e8f0;">
					<button type="button" onclick="closeMerchantModal()" class="btn-secondary">Cancel</button>
					<button type="submit" class="btn-primary" style="background: #087443;">Save Settings</button>
				</div>
			</form>
		</div>
	</div>`, escapeHTML(rest.Name), errHTML, rest.ID, delivChecked, pickupChecked, radiusVal, settings.DeliveryFee, estMinVal, suspendChecked, escapeHTML(settings.SuspensionReason))
}

// RenderMerchantClaimModal renders modal for merchant to claim an unassigned restaurant
func RenderMerchantClaimModal(unassignedRests []*ent.Restaurant, preselectedID int, errMsg string) string {
	errHTML := ""
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="merchant-alert error" style="margin-bottom: 1.25rem;">%s</div>`, escapeHTML(errMsg))
	}

	var cardsSB strings.Builder
	if len(unassignedRests) == 0 {
		cardsSB.WriteString(`<div style="padding: 1.5rem; text-align: center; color: #64748b;">No unassigned restaurants currently found in the directory. You can register a new restaurant instead!</div>`)
	} else {
		for _, r := range unassignedRests {
			phone := r.Phone
			if phone == "" {
				phone = "Listing verified"
			}
			cardBorder := "#e2e8f0"
			cardBg := "#ffffff"
			if r.ID == preselectedID {
				cardBorder = "#087443"
				cardBg = "#f0fdf4"
			}
			cardsSB.WriteString(fmt.Sprintf(`
			<div class="claim-modal-item" id="claim-modal-item-%d" style="display: flex; justify-content: space-between; align-items: center; padding: 0.85rem 1rem; border: 1.5px solid %s; background: %s; border-radius: 10px; margin-bottom: 0.6rem; gap: 0.75rem;">
				<div style="flex: 1;">
					<div style="font-weight: 800; font-size: 0.95rem; color: #0f172a;">%s</div>
					<div style="font-size: 0.8rem; color: #64748b; margin-top: 0.15rem;">📍 %s</div>
					<div style="font-size: 0.75rem; color: #087443; margin-top: 0.2rem; font-weight: 600;">📞 %s</div>
				</div>
				<form method="POST" action="/merchant/restaurant/claim" hx-post="/merchant/restaurant/claim" hx-target="#merchant-modal" style="margin: 0;">
					<input type="hidden" name="restaurant_id" value="%d">
					<button type="submit" class="btn-primary" style="font-size: 0.82rem; padding: 0.45rem 0.85rem; border-radius: 8px; white-space: nowrap; background: #087443;">
						Claim Eatery &rarr;
					</button>
				</form>
			</div>`, r.ID, cardBorder, cardBg, escapeHTML(r.Name), escapeHTML(r.Address), escapeHTML(phone), r.ID))
		}
	}

	return fmt.Sprintf(`
	<div class="merchant-modal-overlay" onclick="if(event.target === this) closeMerchantModal()">
		<div class="merchant-modal-card" style="max-width: 650px;">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<div>
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Claim an Existing Bengaluru Eatery</h3>
					<p style="font-size: 0.85rem; color: #64748b; margin-top: 0.2rem;">Take ownership of a verified outlet to update menu dishes, prices, and operating hours.</p>
				</div>
				<button type="button" onclick="closeMerchantModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			%s

			<div style="margin-bottom: 1rem;">
				<input 
					type="text" 
					id="claimModalSearch" 
					placeholder="🔍 Search eateries by name or locality..." 
					oninput="filterClaimModalItems(this.value)"
					class="form-input" 
					style="width: 100%%; padding: 0.7rem 0.9rem; font-size: 0.92rem; border-radius: 8px; border: 1.5px solid #cbd5e1;"
				>
			</div>

			<div id="claimModalList" style="max-height: 360px; overflow-y: auto; padding-right: 0.25rem;">
				%s
			</div>

			<div style="display: flex; justify-content: space-between; align-items: center; margin-top: 1.25rem; padding-top: 0.85rem; border-top: 1px solid #e2e8f0; font-size: 0.85rem;">
				<span style="color: #64748b;">Not seeing your restaurant?</span>
				<button 
					type="button" 
					hx-get="/merchant/restaurant/new" 
					hx-target="#merchant-modal" 
					style="background: none; border: none; color: #087443; font-weight: 800; cursor: pointer; text-decoration: underline;"
				>
					Register a New Restaurant &rarr;
				</button>
			</div>
		</div>
	</div>
	<script>
		function filterClaimModalItems(term) {
			term = (term || '').toLowerCase().trim();
			var items = document.querySelectorAll('.claim-modal-item');
			items.forEach(function(item) {
				var text = (item.textContent || '').toLowerCase();
				if (!term || text.indexOf(term) !== -1) {
					item.style.display = 'flex';
				} else {
					item.style.display = 'none';
				}
			});
		}
	</script>`, errHTML, cardsSB.String())
}


