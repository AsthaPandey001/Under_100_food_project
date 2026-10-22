package fooddiscovery

import (
	"fmt"
	"strings"

	"github.com/mikestefanello/pagoda/ent"
)

// RenderPartnerRegisterPage renders the dedicated onboarding page for restaurant owners
func RenderPartnerRegisterPage(
	currentUser *ent.User,
	unassignedRests []*ent.Restaurant,
	preselectedClaimID int,
	flashErr, flashMsg string,
) string {
	errBanner := ""
	if flashErr != "" {
		errBanner = fmt.Sprintf(`
		<div class="partner-alert error" style="background: #fef2f2; border: 1.5px solid #fecaca; color: #991b1b; padding: 1rem 1.25rem; border-radius: 12px; margin-bottom: 1.5rem; font-weight: 600; display: flex; align-items: center; gap: 0.65rem;">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
			<span>%s</span>
		</div>`, escapeHTML(flashErr))
	}

	msgBanner := ""
	if flashMsg != "" {
		msgBanner = fmt.Sprintf(`
		<div class="partner-alert success" style="background: #f0fdf4; border: 1.5px solid #bbf7d0; color: #166534; padding: 1rem 1.25rem; border-radius: 12px; margin-bottom: 1.5rem; font-weight: 600; display: flex; align-items: center; gap: 0.65rem;">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
			<span>%s</span>
		</div>`, escapeHTML(flashMsg))
	}

	// Build dropdown options and cards for unassigned eateries
	var claimSelectOptions strings.Builder
	var claimCardsHTML strings.Builder

	claimSelectOptions.WriteString(`<option value="">-- Choose an eatery to claim from the directory --</option>`)
	for _, r := range unassignedRests {
		selectedAttr := ""
		if r.ID == preselectedClaimID {
			selectedAttr = "selected"
		}
		claimSelectOptions.WriteString(fmt.Sprintf(
			`<option value="%d" %s>%s — %s</option>`,
			r.ID, selectedAttr, escapeHTML(r.Name), escapeHTML(r.Address),
		))

		cardBorder := "#e2e8f0"
		cardBg := "#ffffff"
		checkedAttr := ""
		if r.ID == preselectedClaimID {
			cardBorder = "#087443"
			cardBg = "#f0fdf4"
			checkedAttr = "checked"
		}

		phoneText := r.Phone
		if phoneText == "" {
			phoneText = "Contact available on listing"
		}

		claimCardsHTML.WriteString(fmt.Sprintf(`
			<label class="claim-resto-card" id="claim-card-%d" style="display: flex; gap: 0.85rem; align-items: flex-start; padding: 0.85rem 1rem; border: 1.5px solid %s; background: %s; border-radius: 10px; cursor: pointer; transition: all 0.2s ease;">
				<input type="radio" name="claim_restaurant_id" value="%d" %s onchange="selectClaimRadio(%d)" style="margin-top: 0.25rem; accent-color: #087443;">
				<div style="flex: 1;">
					<div style="font-weight: 800; font-size: 0.95rem; color: #0f172a;">%s</div>
					<div style="font-size: 0.82rem; color: #64748b; margin-top: 0.15rem;">📍 %s</div>
					<div style="font-size: 0.78rem; color: #087443; margin-top: 0.25rem; font-weight: 600;">📞 %s</div>
				</div>
			</label>
		`, r.ID, cardBorder, cardBg, r.ID, checkedAttr, r.ID, escapeHTML(r.Name), escapeHTML(r.Address), escapeHTML(phoneText)))
	}

	// Owner credentials section
	ownerFieldsHTML := ""
	if currentUser == nil {
		ownerFieldsHTML = `
		<div style="background: #ffffff; border: 1.5px solid #e2e8f0; border-radius: 14px; padding: 1.5rem; margin-top: 1.5rem;">
			<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.15rem; font-weight: 800; color: #0f172a; margin-bottom: 0.35rem; display: flex; align-items: center; gap: 0.45rem;">
				<span>👤</span> Step 2: Restaurant Owner Account Credentials
			</h3>
			<p style="font-size: 0.85rem; color: #64748b; margin-bottom: 1.25rem;">Create your secure merchant owner credentials to access your dashboard, add menu dishes, and manage hours.</p>

			<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
				<div class="form-group">
					<label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Owner Full Name *</label>
					<input type="text" name="owner_name" class="form-input" placeholder="e.g. Ramesh Kumar" required style="width: 100%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
				</div>
				<div class="form-group">
					<label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Business Email Address *</label>
					<input type="email" name="owner_email" class="form-input" placeholder="e.g. owner@bengalurucafe.com" required style="width: 100%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
				</div>
			</div>

			<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-top: 0.75rem;">
				<div class="form-group">
					<label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Owner Contact Phone *</label>
					<input type="tel" name="owner_phone" class="form-input" placeholder="e.g. +91 98450 12345" required style="width: 100%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
				</div>
				<div class="form-group">
					<label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Password * (min 6 chars)</label>
					<input type="password" name="owner_password" class="form-input" placeholder="••••••••" required minlength="6" style="width: 100%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
				</div>
			</div>
		</div>`
	} else {
		ownerName := currentUser.Name
		if ownerName == "" {
			ownerName = currentUser.Email
		}
		ownerFieldsHTML = fmt.Sprintf(`
		<div style="background: #f0fdf4; border: 1.5px solid #bbf7d0; border-radius: 14px; padding: 1.25rem 1.5rem; margin-top: 1.5rem; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 1rem;">
			<div style="display: flex; align-items: center; gap: 0.75rem;">
				<div style="width: 40px; height: 40px; border-radius: 50%%; background: #087443; color: white; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 1.1rem;">
					✓
				</div>
				<div>
					<div style="font-weight: 800; font-size: 0.95rem; color: #166534;">Connected Account: %s</div>
					<div style="font-size: 0.82rem; color: #15803d;">%s &bull; Submitting will grant you Merchant Partner role and link your restaurant immediately.</div>
				</div>
			</div>
			<input type="hidden" name="existing_user" value="true">
		</div>`, escapeHTML(ownerName), escapeHTML(currentUser.Email))
	}

	activeMode := "claim"
	if preselectedClaimID <= 0 && len(unassignedRests) == 0 {
		activeMode = "new"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Partner With Us — Add Your Restaurant to FoodDiscovery</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    :root {
      --primary-green: #087443;
      --primary-green-hover: #065f37;
      --primary-red: #e23744;
      --text-dark: #0f172a;
      --text-muted: #64748b;
      --bg-page: #f8fafc;
      --border-gray: #e2e8f0;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-page);
      color: var(--text-dark);
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      line-height: 1.5;
    }
    .partner-nav {
      background: #ffffff;
      border-bottom: 1px solid var(--border-gray);
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      box-shadow: 0 1px 3px rgba(0,0,0,0.03);
    }
    .partner-brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      text-decoration: none;
    }
    .partner-logo-pill {
      background: #087443;
      color: white;
      font-family: 'Outfit', sans-serif;
      font-weight: 900;
      font-size: 0.95rem;
      padding: 0.35rem 0.65rem;
      border-radius: 8px;
    }
    .partner-portal-title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.25rem;
      font-weight: 800;
      color: var(--text-dark);
    }
    .partner-container {
      max-width: 900px;
      margin: 2.5rem auto 4rem auto;
      padding: 0 1.5rem;
    }
    .hero-card {
      text-align: center;
      margin-bottom: 2.5rem;
    }
    .hero-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
      background: #dcfce7;
      color: #166534;
      font-size: 0.78rem;
      font-weight: 800;
      padding: 0.3rem 0.8rem;
      border-radius: 9999px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      margin-bottom: 0.85rem;
    }
    .hero-title {
      font-family: 'Outfit', sans-serif;
      font-size: 2.35rem;
      font-weight: 900;
      color: #0f172a;
      letter-spacing: -0.02em;
      line-height: 1.2;
    }
    .hero-subtitle {
      font-size: 1.05rem;
      color: #64748b;
      max-width: 680px;
      margin: 0.75rem auto 0 auto;
    }
    .perks-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 1rem;
      margin: 1.75rem 0 2.5rem 0;
    }
    @media (max-width: 640px) {
      .perks-grid { grid-template-columns: 1fr; }
    }
    .perk-card {
      background: #ffffff;
      border: 1px solid var(--border-gray);
      border-radius: 12px;
      padding: 1.15rem;
      text-align: left;
    }
    .perk-icon {
      font-size: 1.5rem;
      margin-bottom: 0.4rem;
    }
    .perk-title {
      font-weight: 800;
      font-size: 0.95rem;
      color: #0f172a;
    }
    .perk-desc {
      font-size: 0.82rem;
      color: #64748b;
      margin-top: 0.25rem;
    }
    .mode-tab-buttons {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 0.5rem;
      background: #f1f5f9;
      padding: 0.4rem;
      border-radius: 12px;
      margin-bottom: 1.5rem;
    }
    .mode-tab-btn {
      background: none;
      border: none;
      padding: 0.75rem 1rem;
      border-radius: 8px;
      font-size: 0.92rem;
      font-weight: 700;
      color: #475569;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.4rem;
      transition: all 0.2s ease;
    }
    .mode-tab-btn.active {
      background: #ffffff;
      color: #087443;
      box-shadow: 0 2px 4px rgba(0,0,0,0.06);
    }
    .form-group {
      margin-bottom: 1rem;
    }
    .claim-resto-list {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 0.75rem;
      max-height: 380px;
      overflow-y: auto;
      padding: 0.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 10px;
      background: #fafaf9;
    }
    @media (max-width: 640px) {
      .claim-resto-list { grid-template-columns: 1fr; }
    }
    .btn-submit-partner {
      width: 100%%;
      background: #087443;
      color: #ffffff;
      border: none;
      border-radius: 12px;
      padding: 1rem 1.5rem;
      font-size: 1.05rem;
      font-weight: 800;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      box-shadow: 0 4px 12px rgba(8,116,67,0.25);
      transition: all 0.2s ease;
      margin-top: 1.5rem;
    }
    .btn-submit-partner:hover {
      background: #065f37;
      transform: translateY(-1px);
    }
  </style>
</head>
<body>

  <!-- Navigation Bar -->
  <header class="partner-nav">
    <a href="/" class="partner-brand">
      <span class="partner-logo-pill">FOOD</span>
      <span class="partner-portal-title">Partner Program</span>
    </a>
    <div style="display: flex; align-items: center; gap: 1rem;">
      <a href="/" style="color: #64748b; font-size: 0.88rem; font-weight: 700; text-decoration: none;">&larr; Back to Food Discovery</a>
      <a href="/merchant/login" style="background: #087443; border: 1px solid #087443; color: #ffffff; padding: 0.45rem 0.95rem; border-radius: 8px; font-size: 0.85rem; font-weight: 700; text-decoration: none;">Merchant Sign In</a>
    </div>
  </header>

  <div class="partner-container">
    <!-- Hero Header -->
    <div class="hero-card">
      <div class="hero-badge">
        <span>🏪</span> Bengaluru Restaurant Network
      </div>
      <h1 class="hero-title">Bring Your Restaurant to FoodDiscovery</h1>
      <p class="hero-subtitle">
        Whether you are expanding your famous brand or opening a new neighbourhood eatery, join Bengaluru's premiere food discovery network with over 40+ verified restaurant partners.
      </p>
    </div>

    <!-- Value Perks -->
    <div class="perks-grid">
      <div class="perk-card">
        <div class="perk-icon">⚡</div>
        <div class="perk-title">Dedicated Merchant Portal</div>
        <div class="perk-desc">Live Kitchen POS queue, real-time menu dish updates, and item availability controls.</div>
      </div>
      <div class="perk-card">
        <div class="perk-icon">🏷️</div>
        <div class="perk-title">Under &#8377;100 Discovery Feed</div>
        <div class="perk-desc">Showcase signature dishes priced &le; &#8377;100 to reach tens of thousands of hungry diners.</div>
      </div>
      <div class="perk-card">
        <div class="perk-icon">🕒</div>
        <div class="perk-title">Live Operating Hours</div>
        <div class="perk-desc">Full 7-day schedule control plus temporary maintenance closure broadcast.</div>
      </div>
    </div>

    %s
    %s

    <!-- Partner Onboarding Form -->
    <form method="POST" action="/partner/register" id="partnerForm" style="background: #ffffff; border: 1.5px solid #e2e8f0; border-radius: 16px; padding: 2rem; box-shadow: 0 4px 16px rgba(0,0,0,0.03);">
      
      <input type="hidden" name="mode" id="modeInput" value="%s">

      <div style="margin-bottom: 1.25rem;">
        <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a; margin-bottom: 0.35rem;">
          Step 1: Choose Your Restaurant Onboarding Method
        </h2>
        <p style="font-size: 0.88rem; color: #64748b;">
          If your outlet is already listed in Bengaluru, you can claim ownership instantly. If it's a new eatery, register your details below.
        </p>
      </div>

      <!-- Mode Selector Tabs -->
      <div class="mode-tab-buttons">
        <button type="button" class="mode-tab-btn active" id="tabClaimBtn" onclick="switchPartnerMode('claim')">
          <span>📍</span> Claim Existing Bengaluru Eatery (%d available)
        </button>
        <button type="button" class="mode-tab-btn" id="tabNewBtn" onclick="switchPartnerMode('new')">
          <span>+</span> Register a Brand New Restaurant
        </button>
      </div>

      <!-- SECTION A: Claim Existing Bengaluru Eatery -->
      <div id="sectionClaim" style="display: block;">
        <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1.25rem; margin-bottom: 1rem;">
          <label style="display: block; font-size: 0.85rem; font-weight: 700; color: #334155; margin-bottom: 0.5rem;">
            Search & Select Your Bengaluru Eatery:
          </label>
          <input 
            type="text" 
            id="partnerClaimSearch" 
            placeholder="🔍 Type restaurant name or area (e.g. Vidyarthi Bhavan, Brahmin's, Taaza Thindi)..." 
            oninput="filterPartnerEateries(this.value)"
            style="width: 100%%; padding: 0.75rem 1rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.95rem; margin-bottom: 1rem; outline: none;"
          >
          
          <div class="claim-resto-list" id="partnerClaimList">
            %s
          </div>

          <div style="margin-top: 0.85rem; font-size: 0.82rem; color: #64748b; text-align: center;">
            Don't see your restaurant in the list? <button type="button" onclick="switchPartnerMode('new')" style="background: none; border: none; color: #087443; font-weight: 800; cursor: pointer; text-decoration: underline;">Click here to register a brand new outlet &rarr;</button>
          </div>
        </div>
      </div>

      <!-- SECTION B: Register Brand New Restaurant -->
      <div id="sectionNew" style="display: none;">
        <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1.5rem; margin-bottom: 1rem;">
          <div class="form-group">
            <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Restaurant Name *</label>
            <input type="text" name="restaurant_name" id="newRestName" class="form-input" placeholder="e.g. Sri Krishna Tiffin Center" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
          </div>

          <div class="form-group">
            <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Street Address & Locality (Bengaluru) *</label>
            <input type="text" name="restaurant_address" id="newRestAddress" class="form-input" placeholder="e.g. 52, 1st Cross, Koramangala 4th Block, Bengaluru" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
          </div>

          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
            <div class="form-group">
              <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Outlet Phone Number</label>
              <input type="tel" name="restaurant_phone" class="form-input" placeholder="e.g. +91 80 2345 6789" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
            </div>
            <div class="form-group">
              <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Opening Hours Summary</label>
              <input type="text" name="restaurant_hours" class="form-input" placeholder="e.g. 7:00 AM – 10:30 PM" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
            </div>
          </div>

          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
            <div class="form-group">
              <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Latitude</label>
              <input type="number" step="0.000001" name="restaurant_lat" value="12.9716" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
            </div>
            <div class="form-group">
              <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Longitude</label>
              <input type="number" step="0.000001" name="restaurant_lng" value="77.5946" style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
            </div>
          </div>

          <div class="form-group">
            <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Website or Social Profile (Optional)</label>
            <input type="url" name="restaurant_website" class="form-input" placeholder="https://..." style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;">
          </div>

          <div class="form-group">
            <label class="form-label" style="display: block; font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 0.35rem;">Kitchen Story & Specialties</label>
            <textarea name="restaurant_description" rows="2" placeholder="Tell diners about your kitchen's famous signature recipes..." style="width: 100%%; padding: 0.65rem 0.85rem; border: 1.5px solid #cbd5e1; border-radius: 8px; font-size: 0.9rem;"></textarea>
          </div>
        </div>
      </div>

      <!-- Step 2: Owner Account Fields -->
      %s

      <!-- Submit Action -->
      <button type="submit" class="btn-submit-partner">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
        Complete Registration & Launch Merchant Portal
      </button>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.85rem; color: #64748b;">
        By joining, you agree to feature at least 1 dish priced &le; &#8377;100 for the Bengaluru discovery feed.
      </div>
    </form>
  </div>

  <script>
    function switchPartnerMode(mode) {
      document.getElementById('modeInput').value = mode;
      var tabClaim = document.getElementById('tabClaimBtn');
      var tabNew = document.getElementById('tabNewBtn');
      var secClaim = document.getElementById('sectionClaim');
      var secNew = document.getElementById('sectionNew');
      var restName = document.getElementById('newRestName');
      var restAddr = document.getElementById('newRestAddress');

      if (mode === 'claim') {
        tabClaim.classList.add('active');
        tabNew.classList.remove('active');
        secClaim.style.display = 'block';
        secNew.style.display = 'none';
        if (restName) restName.required = false;
        if (restAddr) restAddr.required = false;
      } else {
        tabNew.classList.add('active');
        tabClaim.classList.remove('active');
        secNew.style.display = 'block';
        secClaim.style.display = 'none';
        if (restName) restName.required = true;
        if (restAddr) restAddr.required = true;
      }
    }

    function selectClaimRadio(id) {
      var allCards = document.querySelectorAll('.claim-resto-card');
      allCards.forEach(function(c) {
        c.style.borderColor = '#e2e8f0';
        c.style.background = '#ffffff';
      });
      var activeCard = document.getElementById('claim-card-' + id);
      if (activeCard) {
        activeCard.style.borderColor = '#087443';
        activeCard.style.background = '#f0fdf4';
      }
    }

    function filterPartnerEateries(term) {
      term = (term || '').toLowerCase().trim();
      var cards = document.querySelectorAll('.claim-resto-card');
      cards.forEach(function(card) {
        var text = (card.textContent || '').toLowerCase();
        if (!term || text.indexOf(term) !== -1) {
          card.style.display = 'flex';
        } else {
          card.style.display = 'none';
        }
      });
    }

    // Set initial tab mode
    switchPartnerMode('%s');
  </script>
</body>
</html>`, errBanner, msgBanner, activeMode, len(unassignedRests), claimCardsHTML.String(), ownerFieldsHTML, activeMode)
}
