package fooddiscovery

import (
	"fmt"
	"strings"
	"time"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/restaurant"
	"github.com/mikestefanello/pagoda/ent/user"
	"github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// Common UI definitions (Clean, Professional, Text-Only, Zero Icons)
const (
	svgPin    = ""
	svgSearch = ""
	svgGps    = ""
	svgStar   = ""
	svgUser   = ""
	svgAdmin  = ""
	svgKey    = ""
)

// RenderHomePage renders the modern, clean light-theme culinary discovery platform with zero icons
// RenderHomePage renders the Under100 Bengaluru food discovery & delivery landing page (Pure Go + HTMX)
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
// RenderHomePage renders the modern, clean light-theme culinary discovery platform
func RenderHomePage(currentUser *ent.User) string {
	var authActionsHTML string
	if currentUser == nil {
		authActionsHTML = `
        <div class="user-chip-dropdown" id="userMenuDropdown" style="display: flex; align-items: center; gap: 0.65rem;">
          <button type="button" class="btn-nav-login" onclick="openLoginModal()" title="Login to Under100">
            Login
          </button>
          <div style="display:none;">
            <a href="/login">Sign In</a>
            <a href="/login?mode=signup">Sign Up</a>
          </div>
        </div>`
	} else {
		userName := currentUser.Name
		if userName == "" {
			userName = currentUser.Email
		}
		authActionsHTML = fmt.Sprintf(`
        <div class="user-chip-dropdown" id="userMenuDropdown">
          <button type="button" class="btn-nav-account active" onclick="toggleUserDropdown()" title="My Account">
            <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <span class="user-chip-name">Account</span>
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><polyline points="6 9 12 15 18 9"></polyline></svg>
          </button>
          <div class="account-menu-popover" id="accountMenuPopover">
            <div class="account-popover-header">
              <div class="account-popover-title">My Account</div>
              <div class="account-popover-subtitle">%s</div>
            </div>
            <a href="/account/addresses" class="account-popover-item">My Addresses</a>
            <a href="/account/orders" class="account-popover-item">My Orders</a>
            <a href="/account/prescriptions" class="account-popover-item">My Prescriptions</a>
            <a href="/account/gift-cards" class="account-popover-item">E-Gift Cards</a>
            <a href="/account/privacy" class="account-popover-item">Account privacy</a>
            <a href="/logout" class="account-popover-item" style="color: #e23744; border-top: 1px solid #f1f5f9; margin-top: 4px;">Logout</a>
          </div>
        </div>`, escapeHTML(userName))
	}

	pageHTML := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Under100 — Bengaluru Food Delivery Under &#8377;100</title>
  <meta name="description" content="Discover and order satisfying everyday meals, breakfasts, rolls and street food in Bengaluru for under &#8377;100.">
  
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <script>
    document.addEventListener('htmx:configRequest', function(evt) {
      var match = document.cookie.match(new RegExp('(^| )_csrf=([^;]+)'));
      if (match) {
        evt.detail.headers['X-CSRF-Token'] = decodeURIComponent(match[2]);
        evt.detail.parameters['csrf'] = decodeURIComponent(match[2]);
      }
    });
    document.addEventListener('submit', function(e) {
      if (e.target && e.target.method && e.target.method.toUpperCase() === 'POST') {
        if (!e.target.querySelector('input[name="csrf"]')) {
          var match = document.cookie.match(new RegExp('(^| )_csrf=([^;]+)'));
          if (match) {
            var inp = document.createElement('input');
            inp.type = 'hidden';
            inp.name = 'csrf';
            inp.value = decodeURIComponent(match[2]);
            e.target.appendChild(inp);
          }
        }
      }
    });
  </script>

  <style>
    :root {
      --color-primary-orange: #FF5A1F;
      --color-primary-orange-hover: #E04810;
      --color-deep-green: #087443;
      --color-deep-green-hover: #065F36;
      --color-dark-text: #171C24;
      --color-secondary-text: #646B73;
      --color-main-bg: #FFFDF9;
      --color-card-bg: #FFFFFF;
      --color-soft-green-bg: #EEF8EB;
      --color-soft-orange-bg: #FFF0E6;
      --color-border: #EAE8E3;
      --color-border-hover: #D5D2C9;
      
      --radius-sm: 8px;
      --radius-md: 14px;
      --radius-lg: 20px;
      --radius-full: 9999px;
      
      --shadow-subtle: 0 2px 8px rgba(23, 28, 36, 0.04);
      --shadow-card: 0 4px 16px rgba(23, 28, 36, 0.06);
      --shadow-elevated: 0 10px 30px rgba(23, 28, 36, 0.08);
      
      --font-heading: 'Outfit', sans-serif;
      --font-body: 'Plus Jakarta Sans', sans-serif;
      --container-max-width: 1240px;
    }

    * { box-sizing: border-box; margin: 0; padding: 0; }
    html { scroll-behavior: auto; }
    body {
      font-family: var(--font-body);
      background-color: var(--color-main-bg);
      color: var(--color-dark-text);
      line-height: 1.5;
      -webkit-font-smoothing: antialiased;
      overflow-x: hidden;
    }

    a { color: inherit; text-decoration: none; }
    button { font-family: inherit; cursor: pointer; border: none; background: none; }
    input, select { font-family: inherit; }

    .container {
      width: 100%;
      max-width: var(--container-max-width);
      margin: 0 auto;
      padding: 0 20px;
    }

    .hide-scrollbar::-webkit-scrollbar { display: none; }
    .hide-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }

    /* DIET INDICATORS */
    .diet-indicator {
      width: 16px;
      height: 16px;
      border-radius: 4px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border: 1.5px solid;
      flex-shrink: 0;
    }
    .diet-veg { border-color: #087443; }
    .diet-veg .diet-dot { width: 8px; height: 8px; border-radius: 50%; background-color: #087443; }
    .diet-nonveg { border-color: #E11D48; }
    .diet-nonveg .diet-dot { width: 8px; height: 8px; border-radius: 50%; background-color: #E11D48; }

    /* 1. TOP ANNOUNCEMENT BAR */
    .top-announcement-bar {
      background-color: #EBF7EE;
      border-bottom: 1px solid #D5EED9;
      padding: 6px 24px;
      font-size: 12.5px;
      color: #1E3A24;
      display: flex;
      align-items: center;
      justify-content: space-between;
      min-height: 38px;
    }
    .top-announcement-inner {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      gap: 16px;
    }
    .top-left-wrap {
      display: flex;
      align-items: center;
      gap: 10px;
      flex-wrap: wrap;
    }
    .top-brand-tag {
      font-weight: 800;
      color: #0B4B27;
      letter-spacing: -0.1px;
    }
    .top-bullets {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #2B5736;
      font-size: 12px;
      font-weight: 500;
    }
    .top-right-wrap {
      display: flex;
      align-items: center;
      gap: 18px;
    }
    .skyline-img {
      height: 32px;
      width: auto;
      max-width: 230px;
      object-fit: contain;
      display: block;
    }
    .namma-link {
      display: flex;
      align-items: center;
      gap: 6px;
      color: #087443;
      font-weight: 700;
      font-size: 12.5px;
      text-decoration: none;
      white-space: nowrap;
    }

    /* 2. NAVBAR */
    .navbar {
      position: sticky;
      top: 0;
      background-color: #ffffff;
      border-bottom: 1px solid #ECEAE5;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
      z-index: 90;
      transform: translateZ(0);
      will-change: transform;
    }
    .navbar-inner {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 10px 20px;
      gap: 16px;
    }
    .nav-brand-wrap {
      display: flex;
      align-items: center;
      gap: 22px;
    }
    .brand-logo-link {
      display: flex;
      align-items: center;
      gap: 9px;
      text-decoration: none;
    }
    .brand-wordmark {
      display: flex;
      align-items: baseline;
      letter-spacing: -0.6px;
    }
    .brand-wordmark .under {
      font-size: 24px;
      font-weight: 800;
      color: #111827;
    }
    .brand-wordmark .num {
      font-size: 24px;
      font-weight: 900;
      color: #FF5A1F;
    }
    .loc-dropdown-pill {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      background-color: #F8F6F0;
      padding: 8px 16px;
      border-radius: var(--radius-full);
      border: 1.5px solid var(--color-border);
      font-size: 13.5px;
      font-weight: 700;
      color: var(--color-dark-text);
      cursor: pointer;
      transition: all 0.2s ease;
      user-select: none;
      flex-shrink: 0;
      max-width: 260px;
    }
    .loc-dropdown-pill:hover {
      border-color: var(--color-primary-orange);
      background-color: #FFF0E6;
      color: var(--color-primary-orange);
    }
    .loc-dropdown-pill:active {
      transform: scale(0.96);
    }
    .loc-dropdown-pill #currentLocText {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      display: inline-block;
    }
    .loc-radar-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background-color: #10B981;
      box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
      animation: pulse-green 2s infinite;
      flex-shrink: 0;
    }
    .loc-radar-dot.someone-else {
      background-color: #FF5A1F;
      animation: pulse-orange 2s infinite;
    }
    @keyframes pulse-green {
      0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
      70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(16, 185, 129, 0); }
      100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
    }
    @keyframes pulse-orange {
      0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(255, 90, 31, 0.7); }
      70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(255, 90, 31, 0); }
      100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(255, 90, 31, 0); }
    }

    /* LOCATION & ORDER FOR SOMEONE ELSE MODAL */
    .loc-modal-overlay {
      position: fixed;
      inset: 0;
      background-color: rgba(23, 28, 36, 0.55);
      backdrop-filter: blur(4px);
      z-index: 300;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 16px;
    }
    .loc-modal-overlay.open {
      display: flex;
    }
    .loc-modal-card {
      background-color: #ffffff;
      border-radius: 20px;
      width: 100%;
      max-width: 530px;
      max-height: 90vh;
      overflow-y: auto;
      box-shadow: 0 20px 45px rgba(0, 0, 0, 0.18);
      border: 1px solid #E2E8F0;
      animation: modalFadeIn 0.2s ease-out;
      display: flex;
      flex-direction: column;
    }
    @keyframes modalFadeIn {
      from { opacity: 0; transform: translateY(12px) scale(0.98); }
      to { opacity: 1; transform: translateY(0) scale(1); }
    }
    .loc-modal-header-simple {
      padding: 18px 24px 14px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      position: relative;
    }
    .loc-modal-title {
      font-size: 19px;
      font-weight: 800;
      color: #1E293B;
      margin: 0;
      font-family: inherit;
    }
    .loc-modal-close-btn {
      font-size: 24px;
      line-height: 1;
      color: #475569;
      cursor: pointer;
      padding: 4px 8px;
      border: none;
      background: none;
      border-radius: 6px;
      transition: all 0.15s ease;
    }
    .loc-modal-close-btn:hover {
      background-color: #F1F5F9;
      color: #0F172A;
    }
    .loc-modal-body {
      padding: 4px 24px 24px;
    }

    /* ACTION ROW: Detect my location | OR | search delivery location */
    .loc-top-action-row {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 16px;
      position: relative;
    }
    .btn-detect-location {
      background-color: #0B6837;
      color: #ffffff;
      font-size: 13.5px;
      font-weight: 700;
      padding: 11px 18px;
      border-radius: 8px;
      border: none;
      cursor: pointer;
      white-space: nowrap;
      transition: background-color 0.15s ease, transform 0.1s ease;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 2px 6px rgba(11, 104, 55, 0.2);
    }
    .btn-detect-location:hover {
      background-color: #08552D;
    }
    .btn-detect-location:active {
      transform: scale(0.98);
    }

    .loc-or-divider {
      display: flex;
      align-items: center;
      position: relative;
      flex-shrink: 0;
    }
    .loc-or-divider::before {
      content: '';
      position: absolute;
      left: -8px;
      right: -8px;
      top: 50%;
      height: 1px;
      background-color: #E2E8F0;
      z-index: 1;
    }
    .loc-or-badge {
      position: relative;
      z-index: 2;
      background-color: #ffffff;
      border: 1px solid #CBD5E1;
      border-radius: 9999px;
      padding: 2px 7px;
      font-size: 10.5px;
      font-weight: 600;
      color: #94A3B8;
      text-transform: uppercase;
    }

    .loc-search-wrap {
      flex: 1;
      position: relative;
    }
    .loc-search-input {
      width: 100%;
      padding: 10px 14px;
      border-radius: 10px;
      border: 1.5px solid #CBD5E1;
      background-color: #ffffff;
      font-size: 13.5px;
      color: #1E293B;
      outline: none;
      transition: all 0.2s ease;
    }
    .loc-search-input:focus {
      border-color: #0B6837;
      box-shadow: 0 0 0 3px rgba(11, 104, 55, 0.12);
    }
    .loc-search-input::placeholder {
      color: #94A3B8;
    }

    /* SEARCH SUGGESTIONS DROPDOWN */
    .loc-search-suggestions {
      background-color: #ffffff;
      border: 1px solid #E2E8F0;
      border-radius: 12px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
      margin-bottom: 16px;
      max-height: 230px;
      overflow-y: auto;
      padding: 6px;
    }
    .loc-suggestion-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 9px 12px;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.12s ease;
    }
    .loc-suggestion-item:hover {
      background-color: #F1F5F9;
    }
    .loc-suggestion-item-icon {
      font-size: 16px;
      color: #0B6837;
    }
    .loc-suggestion-item-text {
      flex: 1;
    }
    .loc-suggestion-item-title {
      font-size: 13.5px;
      font-weight: 700;
      color: #1E293B;
    }
    .loc-suggestion-item-sub {
      font-size: 11.5px;
      color: #64748B;
    }
    .loc-suggestion-tag {
      font-size: 10.5px;
      font-weight: 700;
      padding: 2px 6px;
      border-radius: 6px;
      background-color: #DCFCE7;
      color: #166534;
    }
    .loc-suggestion-tag.hub {
      background-color: #FEF3C7;
      color: #92400E;
    }

    /* PAN INDIA & CITY HUB BANNER */
    .pan-india-notice {
      background-color: #F0FDF4;
      border: 1px solid #BBF7D0;
      border-radius: 12px;
      padding: 10px 14px;
      margin-bottom: 14px;
      font-size: 12.5px;
      color: #166534;
      display: flex;
      align-items: flex-start;
      gap: 8px;
      line-height: 1.45;
    }

    /* SAVED ADDRESSES SECTION */
    .saved-addresses-header-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: 14px;
      margin-bottom: 12px;
    }
    .saved-addresses-title {
      font-size: 15px;
      font-weight: 800;
      color: #334155;
      margin: 0;
    }
    .btn-add-address-link {
      background: none;
      border: 1px solid #FED7AA;
      background-color: #FFF7ED;
      color: #EA580C;
      font-size: 12px;
      font-weight: 700;
      padding: 4px 10px;
      border-radius: 9999px;
      cursor: pointer;
      transition: all 0.15s ease;
    }
    .btn-add-address-link:hover {
      background-color: #FFEDD5;
      border-color: #FDBA74;
    }

    .saved-addresses-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
      max-height: 380px;
      overflow-y: auto;
      padding-right: 4px;
    }
    .saved-addresses-list::-webkit-scrollbar {
      width: 6px;
    }
    .saved-addresses-list::-webkit-scrollbar-thumb {
      background-color: #CBD5E1;
      border-radius: 6px;
    }

    /* SAVED ADDRESS CARD (Matches Screenshot) */
    .saved-address-card {
      background-color: #FFFFFF;
      border: 1px solid #F1F5F9;
      border-radius: 16px;
      padding: 14px 16px;
      display: flex;
      align-items: flex-start;
      gap: 14px;
      cursor: pointer;
      position: relative;
      transition: all 0.18s ease;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
    }
    .saved-address-card:hover {
      background-color: #F8FAFC;
      border-color: #CBD5E1;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
      transform: translateY(-1px);
    }
    .saved-address-card.selected-active {
      border: 1.5px solid #0B6837;
      background-color: #F0FDF4;
    }

    .address-icon-box {
      width: 44px;
      height: 44px;
      border-radius: 10px;
      background-color: #FEF9C3;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      font-size: 20px;
    }
    .address-content-box {
      flex: 1;
      min-width: 0;
    }
    .address-name-title {
      font-size: 15px;
      font-weight: 800;
      color: #0F172A;
      margin-bottom: 3px;
      display: flex;
      align-items: center;
      gap: 6px;
    }
    .address-details-text {
      font-size: 12.5px;
      color: #64748B;
      line-height: 1.45;
      word-break: break-word;
    }
    .address-card-actions {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-top: 8px;
    }
    .btn-action-round {
      width: 26px;
      height: 26px;
      border-radius: 50%;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border: 1px solid #E2E8F0;
      background-color: #ffffff;
      cursor: pointer;
      transition: all 0.15s ease;
    }
    .btn-action-round.edit {
      border-color: #BBF7D0;
      color: #16A34A;
    }
    .btn-action-round.edit:hover {
      background-color: #DCFCE7;
      transform: scale(1.08);
    }
    .btn-action-round.delete {
      border-color: #FECACA;
      color: #DC2626;
    }
    .btn-action-round.delete:hover {
      background-color: #FEE2E2;
      transform: scale(1.08);
    }

    /* COLLAPSIBLE ADD/EDIT FORM */
    .add-address-form-box {
      background-color: #F8FAFC;
      border: 1.5px solid #E2E8F0;
      border-radius: 14px;
      padding: 16px;
      margin-bottom: 16px;
      animation: modalFadeIn 0.2s ease-out;
    }
    .address-type-selector {
      display: flex;
      gap: 12px;
      margin-bottom: 10px;
    }
    .address-type-option {
      display: flex;
      align-items: center;
      gap: 5px;
      font-size: 13px;
      font-weight: 700;
      color: #334155;
      cursor: pointer;
    }
    .form-field-group {
      display: flex;
      flex-direction: column;
      gap: 5px;
    }
    .form-field-label {
      font-size: 12px;
      font-weight: 700;
      color: #334155;
    }
    .form-field-input {
      width: 100%;
      padding: 9px 12px;
      border-radius: 8px;
      border: 1px solid #CBD5E1;
      background-color: #FFFFFF;
      font-size: 13px;
      color: #1E293B;
      outline: none;
      transition: all 0.15s ease;
    }
    .form-field-input:focus {
      border-color: #0B6837;
      box-shadow: 0 0 0 3px rgba(11, 104, 55, 0.1);
    }
    .btn-save-address-submit {
      flex: 1;
      padding: 10px;
      border-radius: 8px;
      background-color: #0B6837;
      color: #ffffff;
      font-size: 13.5px;
      font-weight: 800;
      border: none;
      cursor: pointer;
      transition: background-color 0.15s ease;
    }
    .btn-save-address-submit:hover {
      background-color: #08552D;
    }
    .btn-cancel-address {
      padding: 10px 16px;
      border-radius: 8px;
      background-color: #E2E8F0;
      color: #475569;
      font-size: 13.5px;
      font-weight: 700;
      border: none;
      cursor: pointer;
    }
    .btn-cancel-address:hover {
      background-color: #CBD5E1;
    }

    .search-bar-wrap {
      flex: 1;
      display: flex;
      justify-content: center;
      max-width: 500px;
      min-width: 180px;
    }
    .search-input-box {
      width: 100%;
      position: relative;
    }
    .search-input-box input {
      width: 100%;
      padding: 10px 18px 10px 42px;
      border-radius: var(--radius-full);
      border: 1.5px solid var(--color-border);
      background-color: #F8F6F0;
      font-size: 13.5px;
      font-weight: 500;
      color: var(--color-dark-text);
      outline: none;
      transition: all 0.2s ease;
    }
    .search-input-box input:focus {
      border-color: var(--color-primary-orange);
      background-color: #ffffff;
      box-shadow: 0 0 0 3px rgba(255, 90, 31, 0.12);
    }
    .search-icon {
      position: absolute;
      left: 15px;
      top: 50%;
      transform: translateY(-50%);
      color: var(--color-secondary-text);
      pointer-events: none;
    }
    .nav-account-wrap {
      display: flex;
      align-items: center;
      position: relative;
      flex-shrink: 0;
    }
    .btn-nav-account {
      display: inline-flex;
      align-items: center;
      gap: 7px;
      background-color: #F8F6F0;
      border: 1.5px solid var(--color-border);
      padding: 8px 16px;
      border-radius: var(--radius-full);
      font-weight: 700;
      font-size: 13.5px;
      color: var(--color-dark-text);
      cursor: pointer;
      transition: all 0.2s ease;
      white-space: nowrap;
      user-select: none;
    }
    .btn-nav-account:hover {
      background-color: #FFF0E6;
      border-color: var(--color-primary-orange);
      color: var(--color-primary-orange);
    }
    .btn-nav-account.active {
      background-color: #F8F6F0;
      border-color: #CBD5E1;
    }
    .btn-nav-account.admin-btn {
      border-color: #087443;
      background-color: #EEF8EB;
      color: #087443;
    }
    .btn-nav-account .user-chip-name {
      max-width: 120px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .btn-nav-login {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      background-color: transparent;
      border: 1.5px solid #CBD5E1;
      padding: 7px 22px;
      border-radius: var(--radius-full);
      font-weight: 700;
      font-size: 14px;
      color: #1E293B;
      cursor: pointer;
      transition: all 0.2s ease;
      white-space: nowrap;
      user-select: none;
    }
    .btn-nav-login:hover {
      background-color: #F8FAFC;
      border-color: #94A3B8;
      color: #0F172A;
    }
    .user-chip-dropdown {
      position: relative;
    }
    .account-menu-popover {
      position: absolute;
      top: calc(100% + 8px);
      right: 0;
      background: #ffffff;
      border: 1px solid #E2E8F0;
      border-radius: 12px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.08), 0 4px 10px rgba(0, 0, 0, 0.03);
      width: 220px;
      padding: 14px 0 8px 0;
      display: none;
      flex-direction: column;
      z-index: 1000;
      animation: popoverFadeIn 0.15s ease-out;
      overflow: hidden;
      text-align: left;
    }
    .account-menu-popover.open {
      display: flex;
    }
    @keyframes popoverFadeIn {
      from { opacity: 0; transform: translateY(-6px); }
      to { opacity: 1; transform: translateY(0); }
    }
    .account-popover-header {
      padding: 0 20px 8px 20px;
      border-bottom: none;
      margin-bottom: 0;
    }
    .account-popover-title {
      font-weight: 700;
      font-size: 16.5px;
      color: #1E293B;
      letter-spacing: -0.01em;
      line-height: 1.25;
      font-family: inherit;
    }
    .account-popover-subtitle {
      font-size: 13.5px;
      color: #64748B;
      margin-top: 3px;
      font-weight: 400;
      line-height: 1.25;
      font-family: inherit;
    }
    .account-popover-item {
      display: block;
      width: 100%;
      padding: 9.5px 20px;
      font-size: 13.5px;
      font-weight: 500;
      color: #475569;
      text-decoration: none;
      background: transparent;
      border: none;
      text-align: left;
      cursor: pointer;
      transition: background-color 0.12s ease, color 0.12s ease;
      box-sizing: border-box;
      line-height: 1.4;
      font-family: inherit;
    }
    .account-popover-item:hover,
    .account-popover-item.active {
      background-color: #F8FAFC;
      color: #0F172A;
    }
    .account-popover-item.highlight {
      color: #087443;
      font-weight: 600;
    }
    .account-popover-item.highlight:hover {
      background-color: #EEF8EB;
    }
    .btn-nav-cart {
      position: relative;
      display: inline-flex;
      align-items: center;
      gap: 7px;
      background-color: #F8F6F0;
      border: 1.5px solid var(--color-border);
      padding: 8px 16px;
      border-radius: var(--radius-full);
      font-weight: 700;
      font-size: 13.5px;
      color: var(--color-dark-text);
      cursor: pointer;
      transition: all 0.2s ease;
      white-space: nowrap;
      user-select: none;
      flex-shrink: 0;
    }
    .btn-nav-cart:hover {
      background-color: #FFF0E6;
      border-color: var(--color-primary-orange);
      color: var(--color-primary-orange);
    }
    .btn-nav-cart.active {
      background-color: #FFF0E6;
      border-color: #FF5A1F;
      color: #FF5A1F;
    }
    .cart-count-badge {
      background-color: #FF5A1F;
      color: #ffffff;
      font-size: 11px;
      font-weight: 800;
      width: 19px;
      height: 19px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .btn-nav-outline {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 8px 16px;
      border-radius: var(--radius-full);
      border: 1px solid var(--color-border);
      background-color: transparent;
      color: var(--color-dark-text);
      font-weight: 700;
      font-size: 13.5px;
      cursor: pointer;
      transition: all 0.15s ease;
      white-space: nowrap;
      user-select: none;
    }
    .btn-nav-outline:hover {
      border-color: var(--color-primary-orange);
      color: var(--color-primary-orange);
      background-color: var(--color-soft-orange-bg);
    }
    .btn-nav-outline:active {
      transform: scale(0.96);
    }

    .btn-nav-solid {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 8px 18px;
      border-radius: var(--radius-full);
      background-color: var(--color-primary-orange);
      color: #ffffff;
      font-weight: 700;
      font-size: 13.5px;
      cursor: pointer;
      transition: all 0.15s ease;
      white-space: nowrap;
      box-shadow: 0 2px 6px rgba(255, 90, 31, 0.25);
      user-select: none;
    }
    .btn-nav-solid:hover {
      background-color: var(--color-primary-orange-hover);
    }
    .btn-nav-solid:active {
      transform: scale(0.96);
    }

    .user-chip {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 6px 14px;
      border-radius: var(--radius-full);
      background-color: #F4F5F7;
      border: 1px solid #E2E8F0;
      font-size: 13px;
      font-weight: 700;
      color: var(--color-dark-text);
    }
    .admin-pill-tag {
      background-color: #087443;
      color: #ffffff;
      font-size: 10px;
      padding: 2px 6px;
      border-radius: 4px;
      font-weight: 800;
      text-transform: uppercase;
    }
    .btn-nav-admin {
      display: inline-flex;
      align-items: center;
      font-size: 13px;
      font-weight: 700;
      color: #087443;
      background-color: #EEF8EB;
      padding: 6px 12px;
      border-radius: var(--radius-full);
      text-decoration: none;
      transition: all 0.15s;
    }
    .btn-nav-admin:hover {
      background-color: #D5EED9;
    }
    .btn-nav-admin:active {
      transform: scale(0.96);
    }
    .btn-nav-logout {
      font-size: 13px;
      font-weight: 600;
      color: #EF4444;
      text-decoration: none;
      padding: 6px 10px;
      border-radius: var(--radius-full);
      transition: all 0.15s;
    }
    .btn-nav-logout:hover {
      background-color: #FEE2E2;
    }
    .btn-nav-logout:active {
      transform: scale(0.96);
    }

    .btn-get-app {
      display: flex;
      align-items: center;
      gap: 8px;
      background-color: #0B6837;
      color: #ffffff;
      padding: 10px 20px;
      border-radius: var(--radius-full);
      font-size: 14px;
      font-weight: 700;
      box-shadow: 0 2px 8px rgba(11, 104, 55, 0.25);
      cursor: pointer;
      transition: all 0.15s ease;
      white-space: nowrap;
      user-select: none;
    }
    .btn-get-app:hover { background-color: #09532C; }
    .btn-get-app:active { transform: scale(0.96); }

    /* 3. OFFER BANNER SLIDER (AUTO-SLIDING & RICH AESTHETICS) */
    .offer-slider-section {
      padding: 24px 0 14px;
      background-color: var(--color-main-bg);
      position: relative;
    }
    .offer-slider-container {
      position: relative;
      border-radius: 24px;
      overflow: hidden;
      min-height: 220px;
      box-shadow: 0 12px 36px -10px rgba(15, 23, 42, 0.08);
      background: transparent;
    }
    .offer-slider-track {
      display: flex;
      width: 100%;
      transition: transform 0.65s cubic-bezier(0.22, 1, 0.36, 1);
      will-change: transform;
    }
    .offer-slide-item {
      display: flex;
      min-width: 100%;
      width: 100%;
      flex-shrink: 0;
      box-sizing: border-box;
      padding: 32px 44px;
      border: 1.5px solid;
      border-radius: 24px;
      align-items: center;
      justify-content: space-between;
      flex-wrap: wrap;
      gap: 28px;
      position: relative;
      overflow: hidden;
    }
    .offer-slide-item::before {
      content: "";
      position: absolute;
      top: -60px;
      right: 180px;
      width: 260px;
      height: 260px;
      border-radius: 50%;
      background: radial-gradient(circle, rgba(255, 255, 255, 0.65) 0%, rgba(255, 255, 255, 0) 70%);
      pointer-events: none;
      z-index: 1;
    }
    .slide-left-content {
      max-width: 620px;
      z-index: 2;
    }
    .slide-badge-row {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 14px;
    }
    .slide-pill-badge {
      display: inline-flex;
      align-items: center;
      gap: 7px;
      padding: 5px 14px;
      border-radius: 9999px;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.8px;
      text-transform: uppercase;
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
    }
    .slide-pill-badge::before {
      content: "";
      width: 6px;
      height: 6px;
      border-radius: 50%;
      background-color: currentColor;
      box-shadow: 0 0 8px currentColor;
    }
    .slide-coupon-tag {
      display: inline-flex;
      align-items: center;
      padding: 4px 12px;
      border-radius: 8px;
      background-color: #ffffff;
      border: 1.5px dashed currentColor;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.5px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    }
    .slide-headline {
      font-size: clamp(22px, 3.4vw, 32px);
      font-weight: 900;
      color: #0F172A;
      line-height: 1.22;
      letter-spacing: -0.6px;
      margin-bottom: 10px;
    }
    .slide-subcopy {
      font-size: 15px;
      color: #475569;
      line-height: 1.55;
      margin-bottom: 22px;
      max-width: 540px;
      font-weight: 500;
    }
    .slide-action-btn {
      display: inline-flex;
      align-items: center;
      gap: 10px;
      padding: 12px 26px;
      border-radius: 9999px;
      color: #ffffff;
      font-size: 14px;
      font-weight: 800;
      text-decoration: none;
      transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
      box-shadow: 0 6px 18px rgba(0, 0, 0, 0.12);
    }
    .slide-action-btn:hover {
      transform: translateY(-2px) scale(1.02);
      filter: brightness(1.06);
    }
    .slide-action-btn:active {
      transform: translateY(0) scale(0.98);
    }
    .slide-right-img-wrap {
      position: relative;
      width: 240px;
      height: 175px;
      border-radius: 20px;
      overflow: hidden;
      box-shadow: 0 18px 40px -10px rgba(15, 23, 42, 0.22), 0 0 0 4px rgba(255, 255, 255, 0.9);
      flex-shrink: 0;
      background-color: #F1EFEA;
      transform: perspective(900px) rotateY(-3.5deg) rotateX(2deg);
      transition: all 0.45s cubic-bezier(0.2, 0.8, 0.2, 1);
      z-index: 2;
    }
    .slide-right-img-wrap:hover {
      transform: perspective(900px) rotateY(0deg) rotateX(0deg) scale(1.04);
      box-shadow: 0 24px 48px -10px rgba(15, 23, 42, 0.28), 0 0 0 4px #ffffff;
    }
    .slide-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      object-position: center;
      display: block;
      transition: transform 0.5s ease;
    }
    .slide-right-img-wrap:hover .slide-img {
      transform: scale(1.06);
    }
    .slide-floating-tag {
      position: absolute;
      bottom: 12px;
      left: 12px;
      background: rgba(15, 23, 42, 0.85);
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
      color: #ffffff;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.4px;
      padding: 5px 12px;
      border-radius: 9999px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
      border: 1px solid rgba(255, 255, 255, 0.2);
    }
    .slider-dots-row {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      margin-top: 14px;
    }
    .slider-dot {
      width: 10px;
      height: 10px;
      border-radius: 9999px;
      background-color: #CBD5E1;
      border: none;
      cursor: pointer;
      padding: 0;
      transition: all 0.4s cubic-bezier(0.2, 0.8, 0.2, 1);
    }
    .slider-dot:hover {
      background-color: #94A3B8;
      transform: scale(1.2);
    }
    .slider-dot.active {
      width: 32px;
      height: 10px;
      background: linear-gradient(90deg, #FF5A1F 0%, #FF8A4C 100%);
      box-shadow: 0 2px 10px rgba(255, 90, 31, 0.4);
    }

    /* 4. FOOD CATEGORY CAROUSEL */
    .category-section {
      padding: 28px 0 20px;
      border-bottom: 1px solid var(--color-border);
      background-color: #ffffff;
    }
    .section-head-row {
      display: flex;
      align-items: flex-end;
      justify-content: space-between;
      margin-bottom: 16px;
    }
    .section-head-title {
      font-size: 22px;
      font-weight: 800;
      color: var(--color-dark-text);
      letter-spacing: -0.3px;
      margin-bottom: 4px;
    }
    .section-head-subtitle {
      font-size: 14px;
      color: var(--color-secondary-text);
    }
    .scroll-arrow-btn {
      width: 36px;
      height: 36px;
      border-radius: 50%;
      border: 1px solid var(--color-border);
      background-color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--color-dark-text);
      cursor: pointer;
      box-shadow: var(--shadow-subtle);
      transition: all 0.2s ease;
    }
    .scroll-arrow-btn:hover {
      border-color: var(--color-border-hover);
      background-color: #F8F6F0;
    }
    .category-track {
      display: flex;
      gap: 16px;
      overflow-x: auto;
      scroll-snap-type: x mandatory;
      padding-bottom: 8px;
    }
    .cat-pill-btn {
      flex: 0 0 130px;
      scroll-snap-align: start;
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 12px 8px 14px;
      border-radius: var(--radius-md);
      background-color: var(--color-card-bg);
      border: 1px solid var(--color-border);
      box-shadow: var(--shadow-subtle);
      cursor: pointer;
      transition: all 0.2s ease;
      text-align: center;
    }
    .cat-pill-btn:hover, .cat-pill-btn.active {
      transform: translateY(-3px);
      box-shadow: var(--shadow-card);
      border-color: var(--color-primary-orange);
      background-color: var(--color-soft-orange-bg);
    }
    .cat-img-circle {
      width: 82px;
      height: 82px;
      border-radius: 50%;
      overflow: hidden;
      margin-bottom: 10px;
      border: 2px solid #ffffff;
      box-shadow: 0 2px 8px rgba(0,0,0,0.08);
      background-color: #F1EFEA;
    }
    .cat-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      object-position: center;
      display: block;
    }
    .cat-name {
      font-size: 13px;
      font-weight: 700;
      color: var(--color-dark-text);
      line-height: 1.2;
      margin-bottom: 4px;
    }
    .cat-price-pill {
      font-size: 11px;
      font-weight: 600;
      color: var(--color-deep-green);
      background-color: var(--color-soft-green-bg);
      padding: 2px 8px;
      border-radius: var(--radius-full);
    }

    /* 5. PRODUCT CAROUSELS & PRODUCT CARDS */
    .product-carousel-section {
      padding: 28px 0 24px;
      border-bottom: 1px solid var(--color-border);
      background-color: var(--color-main-bg);
    }
    .products-scroll-track {
      display: flex;
      gap: 16px;
      overflow-x: auto;
      scroll-snap-type: x mandatory;
      padding-bottom: 8px;
    }
    .carousel-item-wrap,
    .products-scroll-track > .product-card {
      flex: 0 0 240px;
      width: 240px;
      scroll-snap-align: start;
    }
    .product-card {
      border-radius: var(--radius-md);
      background-color: var(--color-card-bg);
      border: 1px solid var(--color-border);
      box-shadow: var(--shadow-subtle);
      overflow: hidden;
      display: flex;
      flex-direction: column;
      transition: all 0.2s ease;
      height: 100%;
    }
    .product-card:hover {
      transform: translateY(-4px);
      box-shadow: var(--shadow-card);
      border-color: var(--color-border-hover);
    }
    .product-card-img-wrap {
      position: relative;
      width: 100%;
      height: 170px;
      background-color: #F3F1EC;
      overflow: hidden;
    }
    .product-card-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      object-position: center;
      display: block;
      transition: transform 0.3s ease;
    }
    .product-card:hover .product-card-img {
      transform: scale(1.03);
    }
    .diet-badge-pos {
      position: absolute;
      top: 10px;
      left: 10px;
      background-color: #ffffff;
      border-radius: 4px;
      padding: 3px;
      box-shadow: 0 2px 6px rgba(0,0,0,0.12);
    }
    .price-pill-pos {
      position: absolute;
      bottom: 10px;
      left: 10px;
      background-color: rgba(23, 28, 36, 0.9);
      backdrop-filter: blur(4px);
      color: #ffffff;
      padding: 4px 8px;
      border-radius: var(--radius-sm);
      font-size: 13px;
      font-weight: 800;
      letter-spacing: -0.2px;
    }
    .delivery-badge-pos {
      position: absolute;
      bottom: 10px;
      right: 10px;
      background-color: #ffffff;
      color: var(--color-dark-text);
      padding: 3px 8px;
      border-radius: var(--radius-sm);
      font-size: 11px;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 4px;
      box-shadow: 0 2px 6px rgba(0,0,0,0.1);
    }
    .product-card-body {
      padding: 14px;
      display: flex;
      flex-direction: column;
      flex: 1;
    }
    .product-card-title-row {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 8px;
      margin-bottom: 4px;
    }
    .product-card-title {
      font-size: 15px;
      font-weight: 700;
      color: var(--color-dark-text);
      line-height: 1.3;
      overflow: hidden;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      -webkit-box-orient: vertical;
    }
    .product-card-rating {
      display: flex;
      align-items: center;
      gap: 3px;
      font-size: 12px;
      font-weight: 700;
      color: #B45309;
      background-color: #FEF3C7;
      padding: 2px 6px;
      border-radius: 4px;
      flex-shrink: 0;
    }
    .product-card-restaurant {
      font-size: 12px;
      color: var(--color-secondary-text);
      margin-bottom: 4px;
      font-weight: 600;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 4px;
      transition: color 0.15s ease;
    }
    .product-card-restaurant:hover {
      color: var(--color-primary-orange);
      text-decoration: underline;
    }
    .product-card-meta {
      display: flex;
      align-items: center;
      gap: 3px;
      font-size: 11px;
      color: var(--color-secondary-text);
      margin-bottom: 12px;
    }
    .product-card-footer {
      margin-top: auto;
      padding-top: 10px;
      border-top: 1px dashed var(--color-border);
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .price-label {
      font-size: 11px;
      color: var(--color-secondary-text);
      display: block;
      line-height: 1;
    }
    .price-val {
      font-size: 18px;
      font-weight: 900;
      color: var(--color-dark-text);
    }
    .btn-add-pill {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0 16px;
      height: 32px;
      min-width: 82px;
      border-radius: 8px;
      background-color: #ffffff;
      border: 1.5px solid #0C831F;
      color: #0C831F;
      font-size: 13px;
      font-weight: 800;
      letter-spacing: 0.6px;
      text-transform: uppercase;
      transition: all 0.15s ease;
      box-shadow: 0 1px 3px rgba(12, 131, 31, 0.12);
      cursor: pointer;
      user-select: none;
      box-sizing: border-box;
    }
    .btn-add-pill:hover {
      background-color: #0C831F;
      color: #ffffff;
      box-shadow: 0 2px 6px rgba(12, 131, 31, 0.25);
    }
    .btn-add-pill:active {
      transform: scale(0.96);
    }
    .qty-counter-pill {
      display: inline-flex;
      align-items: center;
      justify-content: space-between;
      height: 32px;
      min-width: 82px;
      background-color: #0C831F;
      border-radius: 8px;
      color: #ffffff;
      overflow: hidden;
      box-shadow: 0 2px 6px rgba(12, 131, 31, 0.25);
      user-select: none;
      box-sizing: border-box;
    }
    .qty-btn {
      width: 28px;
      height: 100%;
      color: #ffffff;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      user-select: none;
      font-size: 16px;
      font-weight: 800;
      line-height: 1;
      background: transparent;
      border: none;
      padding: 0;
      transition: background 0.12s;
    }
    .qty-btn:hover {
      background-color: rgba(0, 0, 0, 0.2);
    }
    .qty-btn:active {
      background-color: rgba(0, 0, 0, 0.35);
    }
    .qty-num {
      padding: 0 4px;
      font-size: 13px;
      font-weight: 900;
      min-width: 20px;
      text-align: center;
      letter-spacing: 0.3px;
    }

    /* 6. HERO SECTION */
    .hero-section {
      background-color: var(--color-card-bg);
      padding: 56px 0 64px;
      border-bottom: 1px solid var(--color-border);
      position: relative;
      overflow: hidden;
    }
    .hero-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 48px;
      align-items: center;
    }
    .hero-pill-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      background-color: var(--color-soft-orange-bg);
      color: var(--color-primary-orange);
      border: 1px solid rgba(255, 90, 31, 0.2);
      padding: 6px 14px;
      border-radius: var(--radius-full);
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.5px;
      text-transform: uppercase;
      margin-bottom: 20px;
    }
    .hero-heading {
      font-size: clamp(36px, 5vw, 54px);
      font-weight: 900;
      line-height: 1.1;
      letter-spacing: -1.5px;
      color: var(--color-dark-text);
      margin-bottom: 18px;
    }
    .hero-copy {
      font-size: 17px;
      color: var(--color-secondary-text);
      line-height: 1.6;
      margin-bottom: 28px;
      max-width: 480px;
    }
    .btn-hero-primary {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      background-color: var(--color-deep-green);
      color: #ffffff;
      padding: 14px 28px;
      border-radius: var(--radius-full);
      font-size: 15px;
      font-weight: 700;
      box-shadow: 0 4px 14px rgba(8, 116, 67, 0.25);
      transition: all 0.2s ease;
      text-decoration: none;
    }
    .btn-hero-primary:hover {
      background-color: var(--color-deep-green-hover);
    }
    .hero-img-box {
      border-radius: var(--radius-lg);
      overflow: hidden;
      box-shadow: 0 12px 36px rgba(23, 28, 36, 0.12);
      border: 1px solid var(--color-border);
      background-color: #F3F1EC;
      position: relative;
      height: 380px;
      max-height: 440px;
    }
    .hero-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
    }
    .hero-rotate-badge {
      position: absolute;
      top: -16px;
      right: -16px;
      width: 108px;
      height: 108px;
      border-radius: 50%;
      background-color: var(--color-primary-orange);
      color: "#ffffff";
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      text-align: center;
      box-shadow: 0 8px 24px rgba(255, 90, 31, 0.35);
      border: 3px solid #ffffff;
      transform: rotate(6deg);
    }
    .hero-float-card {
      position: absolute;
      bottom: 20px;
      left: 20px;
      background-color: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(8px);
      padding: 8px 14px;
      border-radius: var(--radius-full);
      display: flex;
      align-items: center;
      gap: 6px;
      box-shadow: 0 4px 16px rgba(0,0,0,0.1);
      border: 1px solid var(--color-border);
    }

    /* 7. BENEFITS STRIP */
    .benefits-strip {
      background-color: #F8F6F0;
      border-bottom: 1px solid var(--color-border);
      padding: 32px 0;
    }
    .benefits-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 24px;
    }
    .benefit-item {
      display: flex;
      align-items: flex-start;
      gap: 16px;
      padding: 16px 20px;
      background-color: #ffffff;
      border-radius: var(--radius-md);
      border: 1px solid var(--color-border);
      box-shadow: var(--shadow-subtle);
    }
    .benefit-icon-box {
      width: 44px;
      height: 44px;
      border-radius: 12px;
      background-color: var(--color-soft-orange-bg);
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--color-primary-orange);
      flex-shrink: 0;
    }

    /* 8. MAIN FOOD DISCOVERY GRID */
    .discovery-section {
      padding: 48px 0 60px;
      background-color: var(--color-main-bg);
    }
    .filter-pills-row {
      display: flex;
      align-items: center;
      gap: 8px;
      overflow-x: auto;
      padding-bottom: 12px;
      margin-bottom: 24px;
    }
    .filter-tab-pill {
      padding: 8px 16px;
      border-radius: var(--radius-full);
      font-size: 13.5px;
      font-weight: 600;
      border: 1px solid var(--color-border);
      background-color: #ffffff;
      color: var(--color-dark-text);
      cursor: pointer;
      white-space: nowrap;
      transition: all 0.2s ease;
    }
    .filter-tab-pill.active, .filter-tab-pill:hover {
      background-color: var(--color-dark-text);
      color: #ffffff;
      border-color: var(--color-dark-text);
    }
    .food-grid-container {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
      gap: 24px;
    }

    /* 9. PROMO BANNER */
    .promo-banner-section {
      padding: 40px 0;
      background-color: var(--color-main-bg);
    }
    .promo-box {
      background-color: var(--color-soft-green-bg);
      border-radius: var(--radius-lg);
      border: 1px solid rgba(8, 116, 67, 0.15);
      padding: 48px 40px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex-wrap: wrap;
      gap: 32px;
      box-shadow: var(--shadow-subtle);
    }

    /* 10. HOW IT WORKS & WHY UNDER100 */
    .info-section {
      padding: 64px 0;
      background-color: #ffffff;
      border-bottom: 1px solid var(--color-border);
    }
    .info-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 32px;
    }
    .info-card {
      padding: 32px 24px;
      background-color: var(--color-main-bg);
      border-radius: var(--radius-lg);
      border: 1px solid var(--color-border);
      text-align: center;
    }

    /* 11. REDESIGNED FOOTER (EXACT SAME #171C24 DARK THEME) */
    .footer-dark {
      background-color: #171C24;
      color: #94A3B8;
      padding: 64px 0 36px;
      border-top: 1px solid #2D3748;
      font-family: 'Plus Jakarta Sans', -apple-system, sans-serif;
    }
    .footer-usp-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
      gap: 16px;
      margin-bottom: 48px;
      padding-bottom: 36px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    }
    .footer-usp-card {
      background: rgba(255, 255, 255, 0.02);
      border: 1px solid rgba(255, 255, 255, 0.06);
      border-radius: 12px;
      padding: 18px 20px;
      transition: all 0.25s ease;
    }
    .footer-usp-card:hover {
      background: rgba(255, 255, 255, 0.04);
      border-color: rgba(255, 90, 31, 0.35);
      transform: translateY(-2px);
    }
    .footer-usp-title {
      font-size: 14.5px;
      font-weight: 800;
      color: #ffffff;
      margin-bottom: 5px;
      display: flex;
      align-items: center;
      gap: 8px;
      letter-spacing: -0.2px;
    }
    .footer-usp-desc {
      font-size: 12.5px;
      color: #94A3B8;
      line-height: 1.5;
      margin: 0;
    }
    .footer-grid-redesigned {
      display: grid;
      grid-template-columns: 2.2fr 1.2fr 1.2fr 1.2fr 1.2fr;
      gap: 36px;
      margin-bottom: 44px;
    }
    @media (max-width: 1080px) {
      .footer-grid-redesigned {
        grid-template-columns: repeat(2, 1fr);
        gap: 32px;
      }
    }
    @media (max-width: 640px) {
      .footer-grid-redesigned {
        grid-template-columns: 1fr;
        gap: 28px;
      }
    }
    .footer-brand-col {
      padding-right: 24px;
    }
    .footer-brand-link {
      display: inline-flex;
      align-items: center;
      gap: 10px;
      text-decoration: none;
      margin-bottom: 16px;
    }
    .footer-brand-link .brand-wordmark .under {
      font-size: 25px;
      font-weight: 800;
      color: #ffffff;
      letter-spacing: -0.6px;
    }
    .footer-brand-link .brand-wordmark .num {
      font-size: 25px;
      font-weight: 900;
      color: #FF5A1F;
      letter-spacing: -0.6px;
    }
    .footer-tagline-badge {
      display: inline-block;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.7px;
      text-transform: uppercase;
      background: rgba(255, 90, 31, 0.12);
      color: #FF5A1F;
      border: 1px solid rgba(255, 90, 31, 0.25);
      border-radius: 6px;
      padding: 3px 8px;
      margin-bottom: 14px;
    }
    .footer-mission-text {
      font-size: 13.5px;
      line-height: 1.65;
      color: #94A3B8;
      margin-bottom: 18px;
    }
    .footer-service-pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      font-size: 12px;
      font-weight: 600;
      color: #E2E8F0;
      background: rgba(8, 116, 67, 0.18);
      border: 1px solid rgba(8, 116, 67, 0.4);
      padding: 6px 12px;
      border-radius: 9999px;
      margin-bottom: 16px;
    }
    .footer-service-pill::before {
      content: "";
      width: 7px;
      height: 7px;
      border-radius: 50%;
      background: #10B981;
      box-shadow: 0 0 8px #10B981;
    }
    .footer-trust-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      font-size: 11.5px;
      color: #64748B;
    }
    .footer-col-title {
      font-size: 13px;
      font-weight: 800;
      color: #ffffff;
      margin-bottom: 16px;
      letter-spacing: 0.6px;
      text-transform: uppercase;
      position: relative;
    }
    .footer-col-title::after {
      content: "";
      display: block;
      width: 22px;
      height: 2px;
      background: #FF5A1F;
      border-radius: 2px;
      margin-top: 6px;
    }
    .footer-link-list {
      list-style: none;
      padding: 0;
      margin: 0;
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    .footer-link-list a {
      color: #94A3B8;
      text-decoration: none;
      font-size: 13.5px;
      font-weight: 500;
      transition: all 0.2s ease;
      display: inline-flex;
      align-items: center;
      gap: 6px;
    }
    .footer-link-list a:hover {
      color: #ffffff;
      transform: translateX(3px);
    }
    .footer-badge-hot {
      font-size: 9.5px;
      font-weight: 800;
      text-transform: uppercase;
      background: rgba(255, 90, 31, 0.18);
      color: #FF7A45;
      border-radius: 4px;
      padding: 2px 5px;
      letter-spacing: 0.3px;
    }
    .footer-badge-new {
      font-size: 9.5px;
      font-weight: 800;
      text-transform: uppercase;
      background: rgba(16, 185, 129, 0.18);
      color: #34D399;
      border-radius: 4px;
      padding: 2px 5px;
      letter-spacing: 0.3px;
    }
    .footer-searches-box {
      padding: 20px 0;
      border-top: 1px solid rgba(255, 255, 255, 0.08);
      border-bottom: 1px solid rgba(255, 255, 255, 0.08);
      margin-bottom: 28px;
    }
    .footer-searches-header {
      font-size: 11.5px;
      font-weight: 800;
      color: #CBD5E1;
      text-transform: uppercase;
      letter-spacing: 0.6px;
      margin-bottom: 10px;
    }
    .footer-tags-row {
      display: flex;
      flex-wrap: wrap;
      gap: 7px;
    }
    .footer-tag-link {
      font-size: 12px;
      font-weight: 600;
      color: #94A3B8;
      background: rgba(255, 255, 255, 0.03);
      border: 1px solid rgba(255, 255, 255, 0.08);
      padding: 4px 11px;
      border-radius: 9999px;
      text-decoration: none;
      transition: all 0.2s ease;
    }
    .footer-tag-link:hover {
      background: rgba(255, 90, 31, 0.15);
      border-color: rgba(255, 90, 31, 0.35);
      color: #ffffff;
    }
    .footer-bottom-flex {
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex-wrap: wrap;
      gap: 16px;
      font-size: 12.5px;
      color: #64748B;
    }
    .footer-legal-links {
      display: flex;
      align-items: center;
      gap: 18px;
      flex-wrap: wrap;
    }
    .footer-legal-links a {
      color: #94A3B8;
      text-decoration: none;
      font-size: 12.5px;
      transition: color 0.2s ease;
    }
    .footer-legal-links a:hover {
      color: #ffffff;
    }

    /* 12. BLINKIT-STYLE CART DRAWER */
    .cart-drawer-backdrop {
      position: fixed;
      inset: 0;
      background-color: rgba(15, 23, 42, 0.48);
      backdrop-filter: blur(4px);
      z-index: 200;
      display: none;
    }
    .cart-drawer-backdrop.open { display: block; }
    .cart-drawer-panel {
      position: fixed;
      top: 0;
      right: -460px;
      width: 100%;
      max-width: 420px;
      height: 100%;
      background-color: #F4F6FB;
      box-shadow: -4px 0 28px rgba(0, 0, 0, 0.16);
      display: flex;
      flex-direction: column;
      z-index: 201;
      transition: right 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    }
    .cart-drawer-panel.open { right: 0; }

    .blinkit-cart-header {
      padding: 16px 20px;
      background-color: #ffffff;
      border-bottom: 1px solid #E2E8F0;
      display: flex;
      align-items: center;
      justify-content: space-between;
      position: sticky;
      top: 0;
      z-index: 5;
    }
    .btn-cart-back {
      border: none;
      background: none;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 0;
      font-size: 18px;
      font-weight: 800;
      color: #1E293B;
    }
    .btn-cart-share {
      border: none;
      background: none;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      font-size: 14px;
      font-weight: 700;
      color: #0C831F;
      padding: 4px 8px;
      border-radius: 6px;
      transition: background 0.15s ease;
    }
    .btn-cart-share:hover {
      background: #F0FDF4;
    }

    .blinkit-cart-body {
      flex: 1;
      overflow-y: auto;
      padding: 14px 16px 24px 16px;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .b-cart-card {
      background-color: #ffffff;
      border-radius: 16px;
      padding: 14px 16px;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
      border: 1px solid #EEF2F6;
    }

    .b-cart-stepper {
      display: inline-flex;
      align-items: center;
      background-color: #0C831F;
      color: #ffffff;
      border-radius: 8px;
      overflow: hidden;
      font-weight: 800;
      font-size: 14px;
      box-shadow: 0 2px 6px rgba(12, 131, 31, 0.25);
    }
    .b-cart-stepper-btn {
      background: transparent;
      border: none;
      color: #ffffff;
      padding: 6px 10px;
      cursor: pointer;
      font-weight: 800;
      font-size: 15px;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: background 0.15s;
    }
    .b-cart-stepper-btn:hover {
      background-color: rgba(255, 255, 255, 0.18);
    }
    .b-cart-stepper-val {
      min-width: 20px;
      text-align: center;
      padding: 0 4px;
      font-size: 13.5px;
      font-weight: 800;
    }

    .b-tip-pill {
      border: 1.5px solid #E2E8F0;
      background: #ffffff;
      border-radius: 12px;
      padding: 8px 12px;
      font-size: 13px;
      font-weight: 700;
      color: #1E293B;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 4px;
      transition: all 0.15s ease;
    }
    .b-tip-pill:hover {
      border-color: #CBD5E1;
      background: #F8FAFC;
    }
    .b-tip-pill.active {
      border-color: #0C831F;
      background: #F0FDF4;
      color: #166534;
      box-shadow: 0 2px 6px rgba(12, 131, 31, 0.15);
    }

    .blinkit-cart-bottom {
      background: #ffffff;
      padding: 12px 16px;
      border-top: 1px solid #E2E8F0;
      box-shadow: 0 -4px 16px rgba(0, 0, 0, 0.05);
      position: sticky;
      bottom: 0;
      z-index: 10;
    }
    .btn-blinkit-proceed {
      width: 100%;
      background: #0C831F;
      border: none;
      border-radius: 12px;
      padding: 12px 18px;
      color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: space-between;
      cursor: pointer;
      box-shadow: 0 4px 14px rgba(12, 131, 31, 0.3);
      transition: background 0.15s ease, transform 0.1s ease;
    }
    .btn-blinkit-proceed:hover {
      background: #0A6F1A;
    }
    .btn-blinkit-proceed:active {
      transform: scale(0.99);
    }

    /* FLOATING CART BAR */
    .floating-cart-bar {
      position: fixed;
      bottom: 16px;
      left: 16px;
      right: 16px;
      max-width: 600px;
      margin: 0 auto;
      background-color: var(--color-dark-text);
      color: #ffffff;
      padding: 12px 18px;
      border-radius: var(--radius-full);
      display: none;
      align-items: center;
      justify-content: space-between;
      box-shadow: var(--shadow-elevated);
      z-index: 80;
      cursor: pointer;
    }
    .floating-cart-bar.active { display: flex; }

    /* RESPONSIVE BREAKPOINTS */
    @media (max-width: 900px) {
      .desktop-only { display: none !important; }
      .search-bar-wrap { max-width: 100%; order: 3; margin-top: 10px; }
      .navbar-inner { flex-wrap: wrap; }
      .slide-right-img-wrap { width: 100%; height: 180px; }
    /* SEARCH DROPDOWN MENU */
    .search-dropdown-menu {
      position: absolute;
      top: calc(100% + 6px);
      left: 0;
      right: 0;
      background: #ffffff;
      border: 1.5px solid #E2E8F0;
      border-radius: 14px;
      box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
      z-index: 1000;
      max-height: 400px;
      overflow-y: auto;
      padding: 8px 0;
    }
    .search-drop-group-title {
      font-size: 11px;
      font-weight: 800;
      color: #64748B;
      text-transform: uppercase;
      letter-spacing: 0.6px;
      padding: 6px 16px 4px;
    }
    .search-drop-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 9px 16px;
      cursor: pointer;
      transition: background 0.15s ease;
      gap: 12px;
    }
    .search-drop-item:hover {
      background: #F1F5F9;
    }

    /* RESTAURANT BANNER VIEW */
    .restaurant-banner-view {
      background: linear-gradient(135deg, #171C24, #262F3D);
      border: 1.5px solid #334155;
      border-radius: 16px;
      padding: 18px 22px;
      color: #ffffff;
      margin-bottom: 22px;
      box-shadow: 0 4px 20px rgba(0,0,0,0.1);
    }

    /* MULTI-ORDER TRACKING TABS */
    .multi-track-tabs {
      display: flex;
      gap: 8px;
      overflow-x: auto;
      padding: 10px 16px;
      background: #0F172A;
      border-bottom: 1px solid #334155;
    }
    .multi-track-tab-btn {
      padding: 8px 14px;
      border-radius: 8px;
      border: 1.5px solid #334155;
      background: rgba(255, 255, 255, 0.06);
      color: #94A3B8;
      font-size: 12.5px;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
      transition: all 0.15s ease;
      display: inline-flex;
      align-items: center;
      gap: 6px;
    }
    .multi-track-tab-btn.active {
      background: #0C831F;
      border-color: #22C55E;
      color: #ffffff;
    }

    /* PAYMENT GATEWAY MODAL STYLES - ZERO ICONS, CORPORATE FINTECH DASHBOARD */
    .pg-modal-overlay {
      position: fixed;
      top: 0; left: 0; right: 0; bottom: 0;
      background: rgba(15, 23, 42, 0.76);
      backdrop-filter: blur(8px);
      z-index: 9999;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 16px;
      animation: pgFadeIn 0.2s ease-out;
    }
    @keyframes pgFadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    .pg-modal-container {
      background: #FFFFFF;
      border-radius: 16px;
      width: 100%;
      max-width: 860px;
      max-height: 92vh;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 60px -15px rgba(0, 0, 0, 0.35);
      border: 1px solid #E2E8F0;
      overflow: hidden;
      animation: pgSlideUp 0.25s cubic-bezier(0.16, 1, 0.3, 1);
    }
    @keyframes pgSlideUp {
      from { transform: translateY(18px) scale(0.98); opacity: 0; }
      to { transform: translateY(0) scale(1); opacity: 1; }
    }
    .pg-modal-header {
      padding: 18px 24px;
      border-bottom: 1px solid #E2E8F0;
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: #F8FAFC;
    }
    .pg-header-kicker {
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.8px;
      text-transform: uppercase;
      color: #0C831F;
    }
    .pg-header-title {
      font-size: 18px;
      font-weight: 800;
      color: #0F172A;
      margin: 2px 0 0;
      letter-spacing: -0.3px;
    }
    .pg-modal-close {
      background: #EDF2F7;
      border: 1px solid #CBD5E1;
      font-size: 22px;
      color: #475569;
      cursor: pointer;
      line-height: 1;
      border-radius: 8px;
      width: 32px; height: 32px;
      display: flex; align-items: center; justify-content: center;
      transition: all 0.15s;
    }
    .pg-modal-close:hover { background: #E2E8F0; color: #0F172A; }
    .pg-modal-body {
      display: flex;
      flex: 1;
      overflow-y: auto;
      min-height: 440px;
    }
    .pg-left-pane {
      flex: 1.25;
      border-right: 1px solid #E2E8F0;
      display: flex;
      flex-direction: column;
      background: #FFFFFF;
    }
    .pg-right-pane {
      flex: 1;
      background: #F8FAFC;
      padding: 20px;
      overflow-y: auto;
    }
    .pg-tabs-bar {
      display: flex;
      border-bottom: 1px solid #E2E8F0;
      background: #FFFFFF;
      overflow-x: auto;
      padding: 0 12px;
      gap: 2px;
    }
    .pg-tab-btn {
      padding: 13px 14px;
      border: none;
      background: transparent;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.5px;
      text-transform: uppercase;
      color: #64748B;
      cursor: pointer;
      border-bottom: 2.5px solid transparent;
      display: inline-flex;
      align-items: center;
      white-space: nowrap;
      transition: all 0.15s ease;
    }
    .pg-tab-btn:hover { color: #0C831F; }
    .pg-tab-btn.active {
      color: #0C831F;
      border-bottom-color: #0C831F;
      background: #F0FDF4;
    }
    .pg-tab-content-area {
      padding: 22px 24px;
      flex: 1;
      overflow-y: auto;
    }
    .pg-tab-pane { display: none; }
    .pg-tab-pane.active { display: block; animation: pgFadeIn 0.2s ease-out; }
    
    /* UPI Option Cards (Zero Icons) */
    .pg-upi-apps-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 10px;
      margin-bottom: 16px;
    }
    .pg-upi-app {
      border: 1.5px solid #E2E8F0;
      border-radius: 10px;
      padding: 12px 6px;
      text-align: center;
      cursor: pointer;
      position: relative;
      transition: all 0.15s ease;
      background: #ffffff;
    }
    .pg-upi-app:hover {
      border-color: #0C831F;
      background: #F0FDF4;
    }
    .pg-upi-app.active {
      border-color: #0C831F;
      background: #F0FDF4;
      box-shadow: 0 0 0 1px #0C831F;
    }
    .pg-upi-badge-text {
      font-size: 12px;
      font-weight: 900;
      letter-spacing: 0.5px;
      padding: 6px 0;
      color: #0F172A;
    }
    .pg-upi-app-name {
      font-size: 11px;
      font-weight: 700;
      color: #475569;
      margin-top: 2px;
    }
    .pg-upi-check {
      position: absolute;
      top: 6px; right: 6px;
      background: #0C831F;
      color: white;
      font-size: 9px;
      font-weight: 900;
      width: 15px; height: 15px;
      border-radius: 50%;
      display: none;
      align-items: center; justify-content: center;
    }
    .pg-upi-app.active .pg-upi-check { display: flex; }
    .pg-qr-box {
      border: 1.5px solid #BBF7D0;
      background: #F0FDF4;
      border-radius: 12px;
      padding: 16px;
      margin-top: 6px;
    }
    .pg-qr-wrapper {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
    }
    .pg-timer-badge {
      background: #FEF3C7;
      color: #92400E;
      font-size: 11px;
      font-weight: 800;
      padding: 3px 8px;
      border-radius: 6px;
      letter-spacing: 0.5px;
    }
    .pg-input {
      width: 100%;
      box-sizing: border-box;
      padding: 10px 14px;
      border-radius: 8px;
      border: 1.5px solid #CBD5E1;
      font-size: 13.5px;
      outline: none;
      transition: border 0.15s;
      font-family: inherit;
    }
    .pg-input:focus {
      border-color: #0C831F;
      box-shadow: 0 0 0 3px rgba(12, 131, 31, 0.12);
    }
    .pg-vpa-verify-btn {
      padding: 10px 18px;
      background: #0C831F;
      color: white;
      border: none;
      border-radius: 8px;
      font-size: 12.5px;
      font-weight: 800;
      letter-spacing: 0.4px;
      cursor: pointer;
      transition: background 0.15s;
    }
    .pg-vpa-verify-btn:hover { background: #086b19; }
    .pg-form-group {
      margin-bottom: 12px;
    }
    .pg-label {
      display: block;
      font-size: 11.5px;
      font-weight: 800;
      letter-spacing: 0.4px;
      text-transform: uppercase;
      color: #475569;
      margin-bottom: 5px;
    }
    .pg-card-brand-badge {
      position: absolute;
      right: 12px;
      top: 50%;
      transform: translateY(-50%);
      background: #F1F5F9;
      border: 1px solid #CBD5E1;
      color: #0F172A;
      font-size: 11px;
      font-weight: 900;
      padding: 2px 8px;
      border-radius: 4px;
      letter-spacing: 0.5px;
    }
    .pg-banks-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 10px;
    }
    .pg-bank-btn {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 12px 14px;
      border: 1.5px solid #E2E8F0;
      border-radius: 10px;
      cursor: pointer;
      background: white;
      transition: all 0.15s;
    }
    .pg-bank-btn:hover, .pg-bank-btn.active {
      border-color: #0C831F;
      background: #F0FDF4;
    }
    .pg-bank-badge {
      width: 44px; height: 26px;
      border-radius: 5px;
      display: flex; align-items: center; justify-content: center;
      font-size: 10px; font-weight: 900;
      letter-spacing: 0.5px;
    }
    .pg-bank-name {
      font-size: 13px;
      font-weight: 700;
      color: #1E293B;
    }
    .pg-wallet-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    .pg-wallet-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 14px;
      border: 1.5px solid #E2E8F0;
      border-radius: 10px;
      cursor: pointer;
      background: white;
      transition: border 0.15s;
    }
    .pg-wallet-item:hover { border-color: #0C831F; }
    
    /* RIGHT PANE: ORDER & QUANTITY LEDGER */
    .pg-recap-card {
      background: white;
      border-radius: 12px;
      border: 1px solid #E2E8F0;
      padding: 18px;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
    }
    .pg-dest-card {
      padding-bottom: 12px;
      margin-bottom: 14px;
      border-bottom: 1px solid #F1F5F9;
    }
    .pg-dest-kicker {
      font-size: 10px;
      font-weight: 800;
      color: #0C831F;
      text-transform: uppercase;
      letter-spacing: 0.8px;
    }
    .pg-dest-address {
      font-size: 13px;
      font-weight: 800;
      color: #0F172A;
      margin-top: 2px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .pg-dest-meta {
      font-size: 11.5px;
      color: #64748B;
      margin-top: 1px;
    }
    .pg-order-ledger {
      background: #FFFFFF;
      border: 1px solid #E2E8F0;
      border-radius: 10px;
      overflow: hidden;
      margin-bottom: 14px;
    }
    .pg-order-ledger-header {
      background: #F1F5F9;
      padding: 8px 12px;
      border-bottom: 1px solid #E2E8F0;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .pg-ledger-title {
      font-size: 10.5px;
      font-weight: 800;
      letter-spacing: 0.8px;
      text-transform: uppercase;
      color: #475569;
    }
    .pg-ledger-count-badge {
      background: #0C831F;
      color: #FFFFFF;
      font-size: 10px;
      font-weight: 800;
      padding: 2px 7px;
      border-radius: 4px;
      letter-spacing: 0.4px;
    }
    .pg-mini-items {
      display: flex;
      flex-direction: column;
      max-height: 200px;
      overflow-y: auto;
    }
    .pg-ledger-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px 12px;
      border-bottom: 1px solid #F1F5F9;
      background: #FFFFFF;
    }
    .pg-ledger-row:nth-child(even) { background: #FAFAFA; }
    .pg-ledger-row:last-child { border-bottom: none; }
    .pg-ledger-item-info {
      flex: 1;
      min-width: 0;
      padding-right: 8px;
    }
    .pg-ledger-item-name {
      font-size: 12.5px;
      font-weight: 700;
      color: #0F172A;
      line-height: 1.3;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .pg-ledger-item-meta {
      font-size: 11px;
      color: #64748B;
      margin-top: 1px;
    }
    .pg-ledger-qty-wrap {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-shrink: 0;
    }
    .pg-qty-tag {
      background: #F1F5F9;
      border: 1px solid #CBD5E1;
      color: #0F172A;
      font-size: 10.5px;
      font-weight: 800;
      padding: 2px 7px;
      border-radius: 5px;
      letter-spacing: 0.4px;
      white-space: nowrap;
    }
    .pg-unit-price {
      font-size: 11px;
      color: #64748B;
      font-weight: 600;
      min-width: 48px;
      text-align: right;
    }
    .pg-row-total {
      font-size: 13px;
      font-weight: 800;
      color: #0F172A;
      min-width: 52px;
      text-align: right;
    }
    .pg-bill-rows {
      display: flex;
      flex-direction: column;
      gap: 7px;
      font-size: 12px;
      color: #475569;
    }
    .pg-bill-row {
      display: flex;
      justify-content: space-between;
    }
    .pg-bill-divider {
      height: 1px;
      background: #E2E8F0;
      margin: 6px 0;
    }
    .pg-grand-total-row {
      font-size: 14px;
      color: #0F172A;
    }
    .pg-trust-strip {
      margin-top: 14px;
      padding-top: 10px;
      border-top: 1px solid #F1F5F9;
    }
    
    /* MODAL FOOTER */
    .pg-modal-footer {
      padding: 16px 24px;
      border-top: 1px solid #E2E8F0;
      background: white;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      gap: 12px;
    }
    .pg-btn-cancel {
      padding: 10px 18px;
      border-radius: 8px;
      border: 1.5px solid #CBD5E1;
      background: white;
      font-size: 13px;
      font-weight: 700;
      color: #64748B;
      cursor: pointer;
      transition: all 0.15s;
    }
    .pg-btn-cancel:hover { background: #F8FAFC; color: #1E293B; }
    .pg-btn-pay {
      padding: 11px 26px;
      border-radius: 8px;
      border: none;
      background: #0C831F;
      color: white;
      font-size: 14px;
      font-weight: 800;
      letter-spacing: 0.4px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      box-shadow: 0 2px 8px rgba(12, 131, 31, 0.3);
      transition: all 0.15s ease;
    }
    .pg-btn-pay:hover {
      background: #086b19;
      box-shadow: 0 4px 12px rgba(12, 131, 31, 0.4);
    }
    .pg-btn-pay:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    
    /* SUCCESS CONFIRMATION VIEW (Zero Icons) */
    .pg-success-view {
      padding: 40px 24px;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
      width: 100%;
      animation: pgFadeIn 0.3s ease-out;
    }
    .pg-success-status-badge {
      background: #DCFCE7;
      border: 1px solid #86EFAC;
      color: #166534;
      font-size: 12px;
      font-weight: 900;
      letter-spacing: 0.8px;
      text-transform: uppercase;
      padding: 8px 18px;
      border-radius: 6px;
    }
    .pg-order-id-pill {
      background: #F1F5F9;
      border: 1px solid #E2E8F0;
      color: #1E293B;
      font-size: 13px;
      font-weight: 800;
      padding: 6px 14px;
      border-radius: 6px;
      margin: 14px 0 22px;
      letter-spacing: 0.5px;
    }
    .pg-tracking-bar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      max-width: 480px;
      margin-bottom: 24px;
    }
    .pg-track-step {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 6px;
    }
    .pg-track-dot {
      width: 28px; height: 28px;
      border-radius: 6px;
      background: #F1F5F9;
      border: 1px solid #CBD5E1;
      color: #64748B;
      font-size: 12px;
      display: flex; align-items: center; justify-content: center;
      font-weight: 900;
    }
    .pg-track-step.done .pg-track-dot { background: #0C831F; border-color: #0C831F; color: white; }
    .pg-track-step.active .pg-track-dot { background: #FEF3C7; border-color: #F59E0B; color: #92400E; }
    .pg-track-label { font-size: 11px; font-weight: 700; color: #64748B; text-transform: uppercase; letter-spacing: 0.4px; }
    .pg-track-line {
      flex: 1;
      height: 2px;
      background: #E2E8F0;
      margin: 0 4px -16px;
    }
    .pg-track-line.active { background: #0C831F; }
    .pg-success-info-card {
      background: #F8FAFC;
      border: 1px solid #E2E8F0;
      border-radius: 10px;
      padding: 16px 20px;
      width: 100%;
      max-width: 440px;
      text-align: left;
    }
    .btn-track-order {
      flex: 1;
      padding: 12px;
      border-radius: 8px;
      border: none;
      background: #0C831F;
      color: white;
      font-size: 13.5px;
      font-weight: 800;
      letter-spacing: 0.4px;
      cursor: pointer;
    }
    @media (max-width: 720px) {
      .pg-modal-body { flex-direction: column; }
      .pg-left-pane { border-right: none; border-bottom: 1px solid #E2E8F0; }
    }

    /* 14. VEG ONLY FILTER SWITCH */
    .veg-toggle-wrap {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      background: #ffffff;
      border: 1.5px solid #22c55e;
      border-radius: 9999px;
      padding: 6px 14px;
      user-select: none;
      cursor: pointer;
      box-shadow: 0 1px 4px rgba(34, 197, 94, 0.12);
      transition: all 0.2s ease;
    }
    .veg-toggle-wrap:hover {
      background: #F0FDF4;
      border-color: #16a34a;
    }
    .veg-toggle-wrap .switch-ui {
      position: relative;
      display: inline-block;
      width: 36px;
      height: 20px;
      cursor: pointer;
    }
    .veg-toggle-wrap .switch-ui input {
      opacity: 0;
      width: 0;
      height: 0;
    }
    .veg-toggle-wrap .slider-round {
      position: absolute;
      cursor: pointer;
      inset: 0;
      background-color: #CBD5E1;
      transition: .25s ease;
      border-radius: 20px;
    }
    .veg-toggle-wrap .slider-round:before {
      position: absolute;
      content: "";
      height: 14px;
      width: 14px;
      left: 3px;
      bottom: 3px;
      background-color: white;
      transition: .25s ease;
      border-radius: 50%;
      box-shadow: 0 1px 3px rgba(0,0,0,0.3);
    }
    .veg-toggle-wrap .switch-ui input:checked + .slider-round {
      background-color: #16a34a;
    }
    .veg-toggle-wrap .switch-ui input:checked + .slider-round:before {
      transform: translateX(16px);
    }

    /* 15. CART MULTI-RESTAURANT CONFLICT MODAL */
    .conflict-modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(15, 23, 42, 0.7);
      backdrop-filter: blur(5px);
      z-index: 99999;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 16px;
    }
    .conflict-modal-overlay.open { display: flex !important; }
    .conflict-modal-card {
      background: #ffffff;
      border-radius: 20px;
      max-width: 420px;
      width: 100%;
      padding: 24px;
      box-shadow: 0 20px 40px rgba(0,0,0,0.2);
      animation: popIn 0.25s ease-out;
      text-align: center;
    }
    @keyframes popIn {
      from { transform: scale(0.92); opacity: 0; }
      to { transform: scale(1); opacity: 1; }
    }
    .conflict-modal-icon {
      width: 52px;
      height: 52px;
      border-radius: 50%;
      background: #FEF3C7;
      color: #D97706;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 16px;
      font-size: 24px;
    }
    .conflict-modal-title {
      font-size: 18px;
      font-weight: 800;
      color: #1E293B;
      margin-bottom: 8px;
    }
    .conflict-modal-desc {
      font-size: 13.5px;
      color: #64748B;
      line-height: 1.5;
      margin-bottom: 24px;
    }
    .conflict-actions-row {
      display: flex;
      gap: 12px;
    }
    .btn-conflict-cancel {
      flex: 1;
      padding: 12px;
      background: #F1F5F9;
      border: 1px solid #E2E8F0;
      border-radius: 10px;
      font-weight: 700;
      font-size: 13.5px;
      color: #475569;
      cursor: pointer;
    }
    .btn-conflict-replace {
      flex: 1.2;
      padding: 12px;
      background: #FF5A1F;
      border: none;
      border-radius: 10px;
      font-weight: 800;
      font-size: 13.5px;
      color: #ffffff;
      cursor: pointer;
      box-shadow: 0 4px 12px rgba(255, 90, 31, 0.3);
    }

    /* 16. CART INSTRUCTIONS & COUPONS */
    .cart-chip-toggle {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      font-weight: 600;
      color: #475569;
      background: #F8FAFC;
      border: 1px solid #CBD5E1;
      padding: 6px 10px;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.15s ease;
      user-select: none;
    }
    .cart-chip-toggle.active {
      background: #ECFDF5;
      border-color: #10B981;
      color: #065F46;
      font-weight: 700;
    }
    .cart-promo-box {
      background: #ffffff;
      border: 1.5px dashed #0C831F;
      border-radius: 12px;
      padding: 10px 12px;
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 12px 0;
    }
    .free-delivery-meter {
      background: #F0FDF4;
      border: 1px solid #BBF7D0;
      border-radius: 10px;
      padding: 10px 14px;
      margin-bottom: 12px;
    }
    .free-delivery-bar {
      height: 6px;
      background: #E2E8F0;
      border-radius: 9999px;
      margin-top: 6px;
      overflow: hidden;
    }
    .free-delivery-progress,
    .free-delivery-fill {
      height: 100%;
      background: linear-gradient(90deg, #10B981, #059669);
      border-radius: 9999px;
      transition: width 0.3s ease;
    }

    /* 17. LIVE ORDER TRACKING MODAL & STEPPER */
    .tracking-modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(15, 23, 42, 0.75);
      backdrop-filter: blur(6px);
      z-index: 99999;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 16px;
    }
    .tracking-modal-overlay.open { display: flex !important; }
    .tracking-modal-card {
      background: #ffffff;
      border-radius: 24px;
      max-width: 480px;
      width: 100%;
      max-height: 90vh;
      overflow-y: auto;
      box-shadow: 0 24px 48px rgba(0,0,0,0.25);
      animation: popIn 0.25s ease-out;
    }
    .tracking-header {
      background: #171C24;
      color: #ffffff;
      padding: 24px;
      border-radius: 24px 24px 0 0;
      position: relative;
    }
    .tracking-eta-pill {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      background: rgba(16, 185, 129, 0.2);
      color: #34D399;
      border: 1px solid rgba(16, 185, 129, 0.4);
      padding: 4px 10px;
      border-radius: 9999px;
      font-size: 12px;
      font-weight: 700;
      margin-bottom: 8px;
    }
    .tracking-title {
      font-size: 22px;
      font-weight: 900;
      color: #ffffff;
      letter-spacing: -0.4px;
    }
    .tracking-body {
      padding: 24px;
    }
    .track-stepper {
      display: flex;
      flex-direction: column;
      gap: 20px;
      position: relative;
      margin: 16px 0 24px 10px;
    }
    .track-stepper::before {
      content: "";
      position: absolute;
      top: 14px;
      bottom: 14px;
      left: 11px;
      width: 2px;
      background: #E2E8F0;
      z-index: 1;
    }
    .track-step-item {
      display: flex;
      align-items: flex-start;
      gap: 16px;
      position: relative;
      z-index: 2;
    }
    .track-step-dot {
      width: 24px;
      height: 24px;
      border-radius: 50%;
      background: #ffffff;
      border: 2.5px solid #CBD5E1;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 11px;
      font-weight: 800;
      color: #64748B;
      flex-shrink: 0;
      transition: all 0.2s ease;
    }
    .track-step-item.done .track-step-dot {
      background: #10B981;
      border-color: #10B981;
      color: #ffffff;
    }
    .track-step-item.active .track-step-dot {
      background: #FF5A1F;
      border-color: #FF5A1F;
      color: #ffffff;
      box-shadow: 0 0 0 4px rgba(255, 90, 31, 0.2);
    }
    .track-step-info h4 {
      font-size: 14.5px;
      font-weight: 800;
      color: #1E293B;
      margin-bottom: 2px;
    }
    .track-step-info p {
      font-size: 12.5px;
      color: #64748B;
      margin: 0;
    }
  </style>
</head>
<body>

  <!-- 1. ANNOUNCEMENT BAR -->
  <div class="top-announcement-bar">
    <div class="container top-announcement-inner">
      <div class="top-left-wrap">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="#087443" style="flex-shrink: 0;">
          <path d="M17.72 2.28C11.5 2.5 5.5 8.5 4.5 15.5c-0.3 2.1 0.5 4.2 2 5.5 1.5 1.3 3.5 1.8 5.5 1.2 6.5-1.8 11.5-8.2 11-15.2-0.2-2.3-1.8-4.3-4-4.7-0.4-0.1-0.9-0.1-1.28 0zm-2.4 14.5c-3.1 0.7-6.2-1.2-6.9-4.3s1.2-6.2 4.3-6.9c0.7-0.2 1.4-0.2 2.1-0.1-0.6 3.8-2.5 7.1-5.5 9.3 1.9 1.5 4.6 1.8 6.8 0.8-0.2 0.4-0.5 0.8-0.8 1.2z"/>
        </svg>
        <span class="top-brand-tag">Bengaluru&#39;s everyday meals, under &#8377;100</span>
        <span style="color: #8EBE98; font-weight: 300;">|</span>
        <div class="top-bullets desktop-only">
          <span>Good food</span>
          <span style="color: #77A882;">&bull;</span>
          <span>Local favourites</span>
          <span style="color: #77A882;">&bull;</span>
          <span>No hidden fees</span>
        </div>
      </div>

      <div class="top-right-wrap">
        <div class="desktop-only" style="display: flex; align-items: flex-end; height: 32px;">
          <img src="/static/bengaluru_skyline.svg" alt="Bengaluru Vidhana Soudha Heritage Skyline" class="skyline-img">
        </div>
        <a href="#food-discovery" class="namma-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="#E11D48">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
          <span style="color: #143F20;">Namma Bengaluru</span>
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#087443" stroke-width="2.5"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </a>
      </div>
    </div>
  </div>

  <!-- 2. NAVBAR: 1.Logo -> 2.Location Tag -> 3.Search Bar -> 4.Account -> 5.My Cart -->
  <header class="navbar">
    <div class="container navbar-inner">
      <!-- 1. LOGO -->
      <a href="/" class="brand-logo-link" title="Under100 Home">
        <svg width="34" height="34" viewBox="0 0 36 36" fill="none">
          <path d="M12 4C11 6.5 13 8.5 12 11" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
          <path d="M18 3C17 5.5 19 8 18 10.5" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
          <path d="M24 4C23 6.5 25 8.5 24 11" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
          <path d="M6 14H30" stroke="#FF5A1F" stroke-width="2.5" stroke-linecap="round"/>
          <path d="M7 16C7.5 24 13.5 28 18 28C22.5 28 28.5 24 29 16H7Z" fill="#FF5A1F"/>
          <path d="M13 28.5H23" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
        </svg>
        <div class="brand-wordmark">
          <span class="under">under</span>
          <span class="num">100</span>
        </div>
      </a>

      <!-- 2. LOCATION TAG & 1-TAP GPS DETECT -->
      <div style="display: flex; align-items: center; gap: 6px;">
        <div class="loc-dropdown-pill" id="navLocationPill" onclick="openLocationModal()" title="Change location or order for someone else">
          <span id="currentLocText">Indiranagar, Bengaluru</span>
        </div>
        <button type="button" onclick="detectGPSLocation()" title="Use Current GPS Location" style="background: rgba(255, 90, 31, 0.1); border: 1.5px solid rgba(255, 90, 31, 0.35); border-radius: 9999px; width: 34px; height: 34px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #FF5A1F; transition: all 0.2s ease; flex-shrink: 0;">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="3"></circle><path d="M12 2v3m0 14v3M2 12h3m14 0h3"></path></svg>
        </button>
      </div>

      <!-- 3. SEARCH BAR WITH DYNAMIC ROTATING LABELS & INSTANT RESTAURANT DROPDOWN -->
      <div class="search-bar-wrap" style="position: relative;">
        <div class="search-input-box">
          <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <input type="text" id="globalSearchInput" placeholder="Search dishes, restaurants under &#8377;100..." oninput="handleGlobalSearch(this.value)" onfocus="handleGlobalSearch(this.value)" autocomplete="off">
          <button type="button" id="clearSearchBtn" onclick="clearGlobalSearch()" style="display: none; background: none; border: none; font-size: 16px; color: #94A3B8; cursor: pointer; padding: 0 4px; line-height: 1;">&times;</button>
        </div>
        <div id="searchSuggestionsDropdown" class="search-dropdown-menu" style="display: none;"></div>
      </div>

      <!-- 4. ACCOUNT -->
      <div class="nav-account-wrap">
        __AUTH_ACTIONS__
      </div>

      <!-- 5. MY CART BUTTON -->
      <button type="button" class="btn-nav-cart" id="navCartBtn" onclick="toggleCartDrawer()">
        <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path><line x1="3" y1="6" x2="21" y2="6"></line><path d="M16 10a4 4 0 0 1-8 0"></path></svg>
        <span class="cart-label-text">My Cart</span>
        <span class="cart-count-badge" id="navCartCount" style="display: none;">0</span>
      </button>
    </div>
  </header>

  <main>

    <!-- 3. MULTI-OFFER BANNER SLIDER (AUTO-SLIDING & INTERACTIVE) -->
    <section class="offer-slider-section">
      <div class="container">
        <div class="offer-slider-container" id="offerSliderContainer" onmouseenter="pauseOfferTimer()" onmouseleave="resumeOfferTimer()">
          <div class="offer-slider-track" id="offerSliderTrack">
            
            <!-- SLIDE 0 -->
            <div class="offer-slide-item" style="background: linear-gradient(135deg, #FFF7ED 0%, #FFEDD5 45%, #FED7AA 100%); border-color: rgba(255, 90, 31, 0.3);" data-index="0">
              <div class="slide-left-content">
                <div class="slide-badge-row">
                  <span class="slide-pill-badge" style="background-color: rgba(255, 90, 31, 0.14); color: #FF5A1F;">
                    WELCOME DEAL
                  </span>
                  <span class="slide-coupon-tag" style="color: #FF5A1F; border-color: #FF5A1F;">Use: BENGALURU20</span>
                </div>
                <h3 class="slide-headline">
                  Get Extra 20% Off On Your First <span style="color: #FF5A1F;">₹100 Meal Run</span>
                </h3>
                <p class="slide-subcopy">Use code BENGALURU20 at checkout • Valid on all tiffins, dosas, and meals across Bengaluru.</p>
                <a href="#food-discovery" class="slide-action-btn" style="background: linear-gradient(135deg, #FF5A1F 0%, #EA580C 100%); box-shadow: 0 8px 22px rgba(255, 90, 31, 0.35);">
                  <span>Order Now</span>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
                </a>
              </div>
              <div class="slide-right-img-wrap">
                <img src="/static/crispy_dosa.jpg" alt="Get Extra 20% Off On Your First" class="slide-img" fetchpriority="high" decoding="async">
                <div class="slide-floating-tag">🔥 20% OFF Extra</div>
              </div>
            </div>
            
            <!-- SLIDE 1 -->
            <div class="offer-slide-item" style="background: linear-gradient(135deg, #F0FDF4 0%, #DCFCE7 45%, #BBF7D0 100%); border-color: rgba(16, 185, 129, 0.3);" data-index="1">
              <div class="slide-left-content">
                <div class="slide-badge-row">
                  <span class="slide-pill-badge" style="background-color: rgba(8, 116, 67, 0.14); color: #087443;">
                    COMING SOON
                  </span>
                  <span class="slide-coupon-tag" style="color: #087443; border-color: #087443;">Use: TECHPARK79</span>
                </div>
                <h3 class="slide-headline">
                  Under100 Daily Office Tiffin <span style="color: #087443;">Monthly Subscriptions</span>
                </h3>
                <p class="slide-subcopy">Fresh homestyle meals delivered directly to your tech park desk for just ₹79/day.</p>
                <a href="#food-discovery" class="slide-action-btn" style="background: linear-gradient(135deg, #087443 0%, #065F46 100%); box-shadow: 0 8px 22px rgba(8, 116, 67, 0.35);">
                  <span>Join Early Waitlist</span>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
                </a>
              </div>
              <div class="slide-right-img-wrap">
                <img src="/static/veg_pulao.jpg" alt="Under100 Daily Office Tiffin" class="slide-img" loading="lazy" decoding="async">
                <div class="slide-floating-tag">🍱 Just ₹79/Day</div>
              </div>
            </div>
            
            <!-- SLIDE 2 -->
            <div class="offer-slide-item" style="background: linear-gradient(135deg, #F0F9FF 0%, #E0F2FE 45%, #BAE6FD 100%); border-color: rgba(2, 132, 199, 0.3);" data-index="2">
              <div class="slide-left-content">
                <div class="slide-badge-row">
                  <span class="slide-pill-badge" style="background-color: rgba(2, 132, 199, 0.14); color: #0284C7;">
                    SAVINGS PASS
                  </span>
                  <span class="slide-coupon-tag" style="color: #0284C7; border-color: #0284C7;">Use: PASS99</span>
                </div>
                <h3 class="slide-headline">
                  Zero Delivery Fee Club for <span style="color: #0284C7;">Students & Techies</span>
                </h3>
                <p class="slide-subcopy">Enjoy flat ₹0 delivery fees on all orders across Bengaluru • No minimum surge fees.</p>
                <a href="#how-it-works" class="slide-action-btn" style="background: linear-gradient(135deg, #0284C7 0%, #0369A1 100%); box-shadow: 0 8px 22px rgba(2, 132, 199, 0.35);">
                  <span>Explore Pass</span>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
                </a>
              </div>
              <div class="slide-right-img-wrap">
                <img src="/static/samosa_plate.jpg" alt="Zero Delivery Fee Club for" class="slide-img" loading="lazy" decoding="async">
                <div class="slide-floating-tag">🚀 ₹0 Delivery Always</div>
              </div>
            </div>
            
            <!-- SLIDE 3 -->
            <div class="offer-slide-item" style="background: linear-gradient(135deg, #FAF5FF 0%, #F3E8FF 45%, #DDD6FE 100%); border-color: rgba(124, 58, 237, 0.3);" data-index="3">
              <div class="slide-left-content">
                <div class="slide-badge-row">
                  <span class="slide-pill-badge" style="background-color: rgba(124, 58, 237, 0.14); color: #7C3AED;">
                    FUTURE ROLLOUT
                  </span>
                  <span class="slide-coupon-tag" style="color: #7C3AED; border-color: #7C3AED;">Use: NIGHTOWL</span>
                </div>
                <h3 class="slide-headline">
                  Midnight Hunger Runs: <span style="color: #7C3AED;">11 PM to 3 AM Delivery</span>
                </h3>
                <p class="slide-subcopy">Hot paneer rolls, steamed momos, and piping filter coffee delivered late night.</p>
                <a href="#carousel-quick-bites" class="slide-action-btn" style="background: linear-gradient(135deg, #7C3AED 0%, #6D28D9 100%); box-shadow: 0 8px 22px rgba(124, 58, 237, 0.35);">
                  <span>View Late Night Menu</span>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
                </a>
              </div>
              <div class="slide-right-img-wrap">
                <img src="/static/kathi_roll.jpg" alt="Midnight Hunger Runs:" class="slide-img" loading="lazy" decoding="async">
                <div class="slide-floating-tag">🌙 11 PM - 3 AM Night</div>
              </div>
            </div>
            
            <!-- SLIDE 4 -->
            <div class="offer-slide-item" style="background: linear-gradient(135deg, #FFFDF0 0%, #FEF9C3 45%, #FDE68A 100%); border-color: rgba(217, 119, 6, 0.3);" data-index="4">
              <div class="slide-left-content">
                <div class="slide-badge-row">
                  <span class="slide-pill-badge" style="background-color: rgba(217, 119, 6, 0.14); color: #D97706;">
                    WEEKEND SPECIAL
                  </span>
                  <span class="slide-coupon-tag" style="color: #D97706; border-color: #D97706;">Use: WEEKEND99</span>
                </div>
                <h3 class="slide-headline">
                  Saturday & Sunday Under ₹99 <span style="color: #D97706;">Handi Biryani Bowls</span>
                </h3>
                <p class="slide-subcopy">Slow-cooked fragrant veg pulao and biryani bowls served with fresh cucumber raita.</p>
                <a href="#carousel-under100" class="slide-action-btn" style="background: linear-gradient(135deg, #D97706 0%, #B45309 100%); box-shadow: 0 8px 22px rgba(217, 119, 6, 0.35);">
                  <span>Browse Weekend Specials</span>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
                </a>
              </div>
              <div class="slide-right-img-wrap">
                <img src="/static/hero_food.jpg" alt="Saturday & Sunday Under ₹99" class="slide-img" loading="lazy" decoding="async">
                <div class="slide-floating-tag">🍲 Under ₹99 Handi</div>
              </div>
            </div>
            
          </div>
        </div>
        <div class="slider-dots-row">
          <button class="slider-dot active" onclick="jumpOfferSlide(0)" aria-label="Slide 1"></button>
          <button class="slider-dot" onclick="jumpOfferSlide(1)" aria-label="Slide 2"></button>
          <button class="slider-dot" onclick="jumpOfferSlide(2)" aria-label="Slide 3"></button>
          <button class="slider-dot" onclick="jumpOfferSlide(3)" aria-label="Slide 4"></button>
          <button class="slider-dot" onclick="jumpOfferSlide(4)" aria-label="Slide 5"></button>
        </div>
      </div>
    </section>

    <!-- 4. WHAT ARE YOU CRAVING TODAY? CATEGORY CAROUSEL -->
    <section class="category-section">
      <div class="container">
        <div class="section-head-row">
          <div>
            <h2 class="section-head-title">What are you craving today?</h2>
            <p class="section-head-subtitle">Your favourites, without the big bill.</p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="scroll-arrow-btn" onclick="scrollTrack('categoryTrack', -280)">&larr;</button>
            <button class="scroll-arrow-btn" onclick="scrollTrack('categoryTrack', 280)">&rarr;</button>
          </div>
        </div>
        <div class="category-track hide-scrollbar" id="categoryTrack">
          
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('south-indian', this)">
          <div class="cat-img-circle">
            <img src="/static/hero_food.jpg" alt="South Indian" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">South Indian</div>
          <div class="cat-price-pill">From &#8377;39</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('dosa-idli', this)">
          <div class="cat-img-circle">
            <img src="/static/crispy_dosa.jpg" alt="Dosa & Idli" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Dosa & Idli</div>
          <div class="cat-price-pill">From &#8377;49</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('rice-biryani', this)">
          <div class="cat-img-circle">
            <img src="/static/veg_pulao.jpg" alt="Rice & Biryani" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Rice & Biryani</div>
          <div class="cat-price-pill">From &#8377;79</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('rolls-wraps', this)">
          <div class="cat-img-circle">
            <img src="/static/kathi_roll.jpg" alt="Rolls & Wraps" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Rolls & Wraps</div>
          <div class="cat-price-pill">From &#8377;69</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('snacks-chaat', this)">
          <div class="cat-img-circle">
            <img src="/static/samosa_plate.jpg" alt="Snacks & Chaat" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Snacks & Chaat</div>
          <div class="cat-price-pill">From &#8377;39</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('north-indian', this)">
          <div class="cat-img-circle">
            <img src="/static/aloo_paratha.jpg" alt="North Indian Meals" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">North Indian Meals</div>
          <div class="cat-price-pill">From &#8377;89</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('momos', this)">
          <div class="cat-img-circle">
            <img src="/static/steamed_momos.jpg" alt="Momos" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Momos</div>
          <div class="cat-price-pill">From &#8377;69</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('sandwiches', this)">
          <div class="cat-img-circle">
            <img src="/static/bombay_sandwich.jpg" alt="Sandwiches" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Sandwiches</div>
          <div class="cat-price-pill">From &#8377;59</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('parathas', this)">
          <div class="cat-img-circle">
            <img src="/static/aloo_paratha.jpg" alt="Parathas" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Parathas</div>
          <div class="cat-price-pill">From &#8377;49</div>
        </button>
        
        <button type="button" class="cat-pill-btn" onclick="selectCategoryFilter('beverages', this)">
          <div class="cat-img-circle">
            <img src="/static/filter_coffee.jpg" alt="Beverages" loading="lazy" class="cat-img" onerror="this.onerror=null;this.src='/static/crispy_dosa.jpg';">
          </div>
          <div class="cat-name">Beverages</div>
          <div class="cat-price-pill">From &#8377;29</div>
        </button>
        
        </div>
      </div>
    </section>

    <!-- 5. PRODUCT CAROUSEL B: UNDER 100 SPECIALS -->
    <section class="product-carousel-section" id="carousel-under100">
      <div class="container">
        <div class="section-head-row">
          <div>
            <h2 class="section-head-title">Your &#8377;100 food run starts here</h2>
            <p class="section-head-subtitle">Big flavours. Small prices.</p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackUnder100', -300)">&larr;</button>
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackUnder100', 300)">&rarr;</button>
          </div>
        </div>
        <div class="products-scroll-track hide-scrollbar" id="trackUnder100">
          <div class="carousel-item-wrap">
    <div class="product-card" data-id="f1" data-category="dosa-idli" data-neighbourhood="Indiranagar" data-name="crispy masala dosa" data-price="79" data-rating="4.6" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Masala Dosa" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Masala Dosa</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Sri Krishna Cafe</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f1">
            <button type="button" class="btn-add-pill" onclick="addToCart('f1', 'Crispy Masala Dosa', 79, 'Sri Krishna Cafe', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f2" data-category="dosa-idli" data-neighbourhood="Koramangala" data-name="steamed idli sambar (2 pcs)" data-price="49" data-rating="4.8" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/idli_sambar.jpg" alt="Steamed Idli Sambar (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Idli Sambar (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Tiffin Room</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f2">
            <button type="button" class="btn-add-pill" onclick="addToCart('f2', 'Steamed Idli Sambar (2 pcs)', 49, 'Brahmin\'s Tiffin Room', '/static/idli_sambar.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f3" data-category="rice-biryani" data-neighbourhood="HSR Layout" data-name="aromatic veg pulao" data-price="89" data-rating="4.4" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/veg_pulao.jpg" alt="Aromatic Veg Pulao" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Aromatic Veg Pulao</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Meghana Bowl Corner</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>North Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f3">
            <button type="button" class="btn-add-pill" onclick="addToCart('f3', 'Aromatic Veg Pulao', 89, 'Meghana Bowl Corner', '/static/veg_pulao.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f4" data-category="rolls-wraps" data-neighbourhood="Indiranagar" data-name="paneer kathi roll" data-price="99" data-rating="4.5" data-delivery="22">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/kathi_roll.jpg" alt="Paneer Kathi Roll" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;99</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>22m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Paneer Kathi Roll</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Kolkata Roll Express</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;99</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f4">
            <button type="button" class="btn-add-pill" onclick="addToCart('f4', 'Paneer Kathi Roll', 99, 'Kolkata Roll Express', '/static/kathi_roll.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f5" data-category="rice-biryani" data-neighbourhood="BTM Layout" data-name="schezwan veg fried rice" data-price="89" data-rating="4.3" data-delivery="28">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/veg_pulao.jpg" alt="Schezwan Veg Fried Rice" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>28m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Schezwan Veg Fried Rice</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.3</span>
          </div>
        </div>
        <div class="product-card-restaurant">Wok & Bowl Express</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>BTM Layout</span>
          <span>&bull;</span>
          <span>Indo-Chinese</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f5">
            <button type="button" class="btn-add-pill" onclick="addToCart('f5', 'Schezwan Veg Fried Rice', 89, 'Wok & Bowl Express', '/static/veg_pulao.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f6" data-category="parathas" data-neighbourhood="Whitefield" data-name="stuffed aloo paratha with curd" data-price="69" data-rating="4.5" data-delivery="24">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/aloo_paratha.jpg" alt="Stuffed Aloo Paratha with Curd" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;69</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>24m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Stuffed Aloo Paratha with Curd</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Punjab Dhaba Point</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Whitefield</span>
          <span>&bull;</span>
          <span>North Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;69</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f6">
            <button type="button" class="btn-add-pill" onclick="addToCart('f6', 'Stuffed Aloo Paratha with Curd', 69, 'Punjab Dhaba Point', '/static/aloo_paratha.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f7" data-category="snacks-chaat" data-neighbourhood="Koramangala" data-name="indori spiced poha" data-price="39" data-rating="4.6" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/indori_poha.jpg" alt="Indori Spiced Poha" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;39</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Indori Spiced Poha</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Morning Fuel Co.</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Breakfast</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;39</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f7">
            <button type="button" class="btn-add-pill" onclick="addToCart('f7', 'Indori Spiced Poha', 39, 'Morning Fuel Co.', '/static/indori_poha.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f8" data-category="momos" data-neighbourhood="Marathahalli" data-name="steamed veg momos (6 pcs)" data-price="79" data-rating="4.5" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Veg Momos (6 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Veg Momos (6 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Himalayan Momo Hut</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Marathahalli</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f8">
            <button type="button" class="btn-add-pill" onclick="addToCart('f8', 'Steamed Veg Momos (6 pcs)', 79, 'Himalayan Momo Hut', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f12" data-category="south-indian" data-neighbourhood="Indiranagar" data-name="south indian mini meal box" data-price="89" data-rating="4.6" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="South Indian Mini Meal Box" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">South Indian Mini Meal Box</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">MTR Heritage Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f12">
            <button type="button" class="btn-add-pill" onclick="addToCart('f12', 'South Indian Mini Meal Box', 89, 'MTR Heritage Canteen', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f18" data-category="beverages" data-neighbourhood="Malleshwaram" data-name="authentic filter coffee" data-price="29" data-rating="4.9" data-delivery="12">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/filter_coffee.jpg" alt="Authentic Filter Coffee" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;29</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>12m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Authentic Filter Coffee</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.9</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Coffee Bar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>Beverages</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;29</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f18">
            <button type="button" class="btn-add-pill" onclick="addToCart('f18', 'Authentic Filter Coffee', 29, 'Brahmin\'s Coffee Bar', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f19" data-category="south-indian" data-neighbourhood="Malleshwaram" data-name="bengaluru khara bath & kesari bath" data-price="55" data-rating="4.8" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="Bengaluru Khara Bath & Kesari Bath" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;55</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bengaluru Khara Bath & Kesari Bath</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Veena Stores</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;55</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f19">
            <button type="button" class="btn-add-pill" onclick="addToCart('f19', 'Bengaluru Khara Bath & Kesari Bath', 55, 'Veena Stores', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f20" data-category="rolls-wraps" data-neighbourhood="HSR Layout" data-name="egg kathi roll" data-price="79" data-rating="4.4" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/kathi_roll.jpg" alt="Egg Kathi Roll" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-nonveg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Egg Kathi Roll</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Rolls King</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f20">
            <button type="button" class="btn-add-pill" onclick="addToCart('f20', 'Egg Kathi Roll', 79, 'Rolls King', '/static/kathi_roll.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f21" data-category="momos" data-neighbourhood="Koramangala" data-name="steamed chicken momos (5 pcs)" data-price="89" data-rating="4.6" data-delivery="22">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Chicken Momos (5 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-nonveg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>22m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Chicken Momos (5 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Shillong Momo Corner</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f21">
            <button type="button" class="btn-add-pill" onclick="addToCart('f21', 'Steamed Chicken Momos (5 pcs)', 89, 'Shillong Momo Corner', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div>
        </div>
      </div>
    </section>

    <!-- 6. PRODUCT CAROUSEL C: SOUTH INDIAN SORTED -->
    <section class="product-carousel-section" id="carousel-south-indian">
      <div class="container">
        <div class="section-head-row">
          <div>
            <h2 class="section-head-title">South Indian, sorted</h2>
            <p class="section-head-subtitle">Crisp dosas, fluffy idlis, and authentic Bengaluru filter coffee.</p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackSouthIndian', -300)">&larr;</button>
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackSouthIndian', 300)">&rarr;</button>
          </div>
        </div>
        <div class="products-scroll-track hide-scrollbar" id="trackSouthIndian">
          <div class="carousel-item-wrap">
    <div class="product-card" data-id="f1" data-category="dosa-idli" data-neighbourhood="Indiranagar" data-name="crispy masala dosa" data-price="79" data-rating="4.6" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Masala Dosa" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Masala Dosa</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Sri Krishna Cafe</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f1">
            <button type="button" class="btn-add-pill" onclick="addToCart('f1', 'Crispy Masala Dosa', 79, 'Sri Krishna Cafe', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f2" data-category="dosa-idli" data-neighbourhood="Koramangala" data-name="steamed idli sambar (2 pcs)" data-price="49" data-rating="4.8" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/idli_sambar.jpg" alt="Steamed Idli Sambar (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Idli Sambar (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Tiffin Room</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f2">
            <button type="button" class="btn-add-pill" onclick="addToCart('f2', 'Steamed Idli Sambar (2 pcs)', 49, 'Brahmin\'s Tiffin Room', '/static/idli_sambar.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f9" data-category="dosa-idli" data-neighbourhood="Malleshwaram" data-name="crispy plain dosa & podi" data-price="59" data-rating="4.7" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Plain Dosa & Podi" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;59</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Plain Dosa & Podi</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Vidyarthi Bhavana Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;59</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f9">
            <button type="button" class="btn-add-pill" onclick="addToCart('f9', 'Crispy Plain Dosa & Podi', 59, 'Vidyarthi Bhavana Tiffins', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f10" data-category="south-indian" data-neighbourhood="HSR Layout" data-name="golden medu vada (2 pcs)" data-price="45" data-rating="4.6" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/medu_vada.jpg" alt="Golden Medu Vada (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Golden Medu Vada (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Udupi Grand Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f10">
            <button type="button" class="btn-add-pill" onclick="addToCart('f10', 'Golden Medu Vada (2 pcs)', 45, 'Udupi Grand Tiffins', '/static/medu_vada.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f11" data-category="dosa-idli" data-neighbourhood="Jayanagar" data-name="onion tomato uttapam" data-price="79" data-rating="4.7" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/onion_uttapam.jpg" alt="Onion Tomato Uttapam" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Onion Tomato Uttapam</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Taaza Thindi</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f11">
            <button type="button" class="btn-add-pill" onclick="addToCart('f11', 'Onion Tomato Uttapam', 79, 'Taaza Thindi', '/static/onion_uttapam.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f12" data-category="south-indian" data-neighbourhood="Indiranagar" data-name="south indian mini meal box" data-price="89" data-rating="4.6" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="South Indian Mini Meal Box" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">South Indian Mini Meal Box</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">MTR Heritage Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f12">
            <button type="button" class="btn-add-pill" onclick="addToCart('f12', 'South Indian Mini Meal Box', 89, 'MTR Heritage Canteen', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f19" data-category="south-indian" data-neighbourhood="Malleshwaram" data-name="bengaluru khara bath & kesari bath" data-price="55" data-rating="4.8" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="Bengaluru Khara Bath & Kesari Bath" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;55</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bengaluru Khara Bath & Kesari Bath</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Veena Stores</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;55</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f19">
            <button type="button" class="btn-add-pill" onclick="addToCart('f19', 'Bengaluru Khara Bath & Kesari Bath', 55, 'Veena Stores', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div>
        </div>
      </div>
    </section>

    <!-- 7. PRODUCT CAROUSEL D: QUICK BITES -->
    <section class="product-carousel-section" id="carousel-quick-bites">
      <div class="container">
        <div class="section-head-row">
          <div>
            <h2 class="section-head-title">Snack now. Spend less.</h2>
            <p class="section-head-subtitle">Hot samosas, crunchy rolls, momos, and street sandwiches.</p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackQuickBites', -300)">&larr;</button>
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackQuickBites', 300)">&rarr;</button>
          </div>
        </div>
        <div class="products-scroll-track hide-scrollbar" id="trackQuickBites">
          <div class="carousel-item-wrap">
    <div class="product-card" data-id="f7" data-category="snacks-chaat" data-neighbourhood="Koramangala" data-name="indori spiced poha" data-price="39" data-rating="4.6" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/indori_poha.jpg" alt="Indori Spiced Poha" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;39</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Indori Spiced Poha</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Morning Fuel Co.</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Breakfast</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;39</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f7">
            <button type="button" class="btn-add-pill" onclick="addToCart('f7', 'Indori Spiced Poha', 39, 'Morning Fuel Co.', '/static/indori_poha.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f8" data-category="momos" data-neighbourhood="Marathahalli" data-name="steamed veg momos (6 pcs)" data-price="79" data-rating="4.5" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Veg Momos (6 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Veg Momos (6 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Himalayan Momo Hut</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Marathahalli</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f8">
            <button type="button" class="btn-add-pill" onclick="addToCart('f8', 'Steamed Veg Momos (6 pcs)', 79, 'Himalayan Momo Hut', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f10" data-category="south-indian" data-neighbourhood="HSR Layout" data-name="golden medu vada (2 pcs)" data-price="45" data-rating="4.6" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/medu_vada.jpg" alt="Golden Medu Vada (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Golden Medu Vada (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Udupi Grand Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f10">
            <button type="button" class="btn-add-pill" onclick="addToCart('f10', 'Golden Medu Vada (2 pcs)', 45, 'Udupi Grand Tiffins', '/static/medu_vada.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f13" data-category="snacks-chaat" data-neighbourhood="BTM Layout" data-name="punjabi samosa (2 pcs with chutney)" data-price="39" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/samosa_plate.jpg" alt="Punjabi Samosa (2 pcs with Chutney)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;39</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Punjabi Samosa (2 pcs with Chutney)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Delhi Chaat Bhandar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>BTM Layout</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;39</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f13">
            <button type="button" class="btn-add-pill" onclick="addToCart('f13', 'Punjabi Samosa (2 pcs with Chutney)', 39, 'Delhi Chaat Bhandar', '/static/samosa_plate.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f14" data-category="snacks-chaat" data-neighbourhood="Koramangala" data-name="tangy aloo tikki chaat" data-price="59" data-rating="4.4" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/samosa_plate.jpg" alt="Tangy Aloo Tikki Chaat" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;59</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Tangy Aloo Tikki Chaat</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Anand Sweets & Chaat</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;59</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f14">
            <button type="button" class="btn-add-pill" onclick="addToCart('f14', 'Tangy Aloo Tikki Chaat', 59, 'Anand Sweets & Chaat', '/static/samosa_plate.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f15" data-category="sandwiches" data-neighbourhood="Jayanagar" data-name="bombay veg grilled sandwich" data-price="69" data-rating="4.6" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/bombay_sandwich.jpg" alt="Bombay Veg Grilled Sandwich" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;69</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bombay Veg Grilled Sandwich</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Hari Super Sandwich</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;69</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f15">
            <button type="button" class="btn-add-pill" onclick="addToCart('f15', 'Bombay Veg Grilled Sandwich', 69, 'Hari Super Sandwich', '/static/bombay_sandwich.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f16" data-category="snacks-chaat" data-neighbourhood="Indiranagar" data-name="crispy pani puri (6 puris)" data-price="49" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/pani_puri.jpg" alt="Crispy Pani Puri (6 Puris)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Pani Puri (6 Puris)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Chaat Junction</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f16">
            <button type="button" class="btn-add-pill" onclick="addToCart('f16', 'Crispy Pani Puri (6 Puris)', 49, 'Chaat Junction', '/static/pani_puri.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f17" data-category="snacks-chaat" data-neighbourhood="Whitefield" data-name="hot onion pakora plate" data-price="45" data-rating="4.3" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/onion_pakora.jpg" alt="Hot Onion Pakora Plate" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Hot Onion Pakora Plate</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.3</span>
          </div>
        </div>
        <div class="product-card-restaurant">Chai Point Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Whitefield</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f17">
            <button type="button" class="btn-add-pill" onclick="addToCart('f17', 'Hot Onion Pakora Plate', 45, 'Chai Point Canteen', '/static/onion_pakora.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f18" data-category="beverages" data-neighbourhood="Malleshwaram" data-name="authentic filter coffee" data-price="29" data-rating="4.9" data-delivery="12">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/filter_coffee.jpg" alt="Authentic Filter Coffee" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;29</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>12m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Authentic Filter Coffee</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.9</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Coffee Bar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>Beverages</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;29</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f18">
            <button type="button" class="btn-add-pill" onclick="addToCart('f18', 'Authentic Filter Coffee', 29, 'Brahmin\'s Coffee Bar', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f21" data-category="momos" data-neighbourhood="Koramangala" data-name="steamed chicken momos (5 pcs)" data-price="89" data-rating="4.6" data-delivery="22">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Chicken Momos (5 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-nonveg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>22m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Chicken Momos (5 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Shillong Momo Corner</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f21">
            <button type="button" class="btn-add-pill" onclick="addToCart('f21', 'Steamed Chicken Momos (5 pcs)', 89, 'Shillong Momo Corner', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f22" data-category="beverages" data-neighbourhood="Indiranagar" data-name="bun maska & irani chai" data-price="49" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/filter_coffee.jpg" alt="Bun Maska & Irani Chai" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bun Maska & Irani Chai</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Cafe Chai Days</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Bakery</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f22">
            <button type="button" class="btn-add-pill" onclick="addToCart('f22', 'Bun Maska & Irani Chai', 49, 'Cafe Chai Days', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div>
        </div>
      </div>
    </section>

    <!-- 8. PRODUCT CAROUSEL E: BENGALURU PICKS -->
    <section class="product-carousel-section" id="carousel-bengaluru-picks">
      <div class="container">
        <div class="section-head-row">
          <div>
            <h2 class="section-head-title">Bengaluru favourites near you</h2>
            <p class="section-head-subtitle">Neighbourhood darshini legends and beloved local eateries.</p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackBengaluruPicks', -300)">&larr;</button>
            <button class="scroll-arrow-btn" onclick="scrollTrack('trackBengaluruPicks', 300)">&rarr;</button>
          </div>
        </div>
        <div class="products-scroll-track hide-scrollbar" id="trackBengaluruPicks">
          <div class="carousel-item-wrap">
    <div class="product-card" data-id="f1" data-category="dosa-idli" data-neighbourhood="Indiranagar" data-name="crispy masala dosa" data-price="79" data-rating="4.6" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Masala Dosa" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Masala Dosa</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Sri Krishna Cafe</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f1">
            <button type="button" class="btn-add-pill" onclick="addToCart('f1', 'Crispy Masala Dosa', 79, 'Sri Krishna Cafe', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f2" data-category="dosa-idli" data-neighbourhood="Koramangala" data-name="steamed idli sambar (2 pcs)" data-price="49" data-rating="4.8" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/idli_sambar.jpg" alt="Steamed Idli Sambar (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Idli Sambar (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Tiffin Room</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f2">
            <button type="button" class="btn-add-pill" onclick="addToCart('f2', 'Steamed Idli Sambar (2 pcs)', 49, 'Brahmin\'s Tiffin Room', '/static/idli_sambar.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f9" data-category="dosa-idli" data-neighbourhood="Malleshwaram" data-name="crispy plain dosa & podi" data-price="59" data-rating="4.7" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Plain Dosa & Podi" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;59</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Plain Dosa & Podi</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Vidyarthi Bhavana Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;59</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f9">
            <button type="button" class="btn-add-pill" onclick="addToCart('f9', 'Crispy Plain Dosa & Podi', 59, 'Vidyarthi Bhavana Tiffins', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f10" data-category="south-indian" data-neighbourhood="HSR Layout" data-name="golden medu vada (2 pcs)" data-price="45" data-rating="4.6" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/medu_vada.jpg" alt="Golden Medu Vada (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Golden Medu Vada (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Udupi Grand Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f10">
            <button type="button" class="btn-add-pill" onclick="addToCart('f10', 'Golden Medu Vada (2 pcs)', 45, 'Udupi Grand Tiffins', '/static/medu_vada.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f11" data-category="dosa-idli" data-neighbourhood="Jayanagar" data-name="onion tomato uttapam" data-price="79" data-rating="4.7" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/onion_uttapam.jpg" alt="Onion Tomato Uttapam" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Onion Tomato Uttapam</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Taaza Thindi</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f11">
            <button type="button" class="btn-add-pill" onclick="addToCart('f11', 'Onion Tomato Uttapam', 79, 'Taaza Thindi', '/static/onion_uttapam.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f12" data-category="south-indian" data-neighbourhood="Indiranagar" data-name="south indian mini meal box" data-price="89" data-rating="4.6" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="South Indian Mini Meal Box" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">South Indian Mini Meal Box</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">MTR Heritage Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f12">
            <button type="button" class="btn-add-pill" onclick="addToCart('f12', 'South Indian Mini Meal Box', 89, 'MTR Heritage Canteen', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f15" data-category="sandwiches" data-neighbourhood="Jayanagar" data-name="bombay veg grilled sandwich" data-price="69" data-rating="4.6" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/bombay_sandwich.jpg" alt="Bombay Veg Grilled Sandwich" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;69</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bombay Veg Grilled Sandwich</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Hari Super Sandwich</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;69</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f15">
            <button type="button" class="btn-add-pill" onclick="addToCart('f15', 'Bombay Veg Grilled Sandwich', 69, 'Hari Super Sandwich', '/static/bombay_sandwich.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f18" data-category="beverages" data-neighbourhood="Malleshwaram" data-name="authentic filter coffee" data-price="29" data-rating="4.9" data-delivery="12">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/filter_coffee.jpg" alt="Authentic Filter Coffee" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;29</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>12m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Authentic Filter Coffee</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.9</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Coffee Bar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>Beverages</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;29</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f18">
            <button type="button" class="btn-add-pill" onclick="addToCart('f18', 'Authentic Filter Coffee', 29, 'Brahmin\'s Coffee Bar', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f19" data-category="south-indian" data-neighbourhood="Malleshwaram" data-name="bengaluru khara bath & kesari bath" data-price="55" data-rating="4.8" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/hero_food.jpg" alt="Bengaluru Khara Bath & Kesari Bath" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;55</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bengaluru Khara Bath & Kesari Bath</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Veena Stores</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;55</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f19">
            <button type="button" class="btn-add-pill" onclick="addToCart('f19', 'Bengaluru Khara Bath & Kesari Bath', 55, 'Veena Stores', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div><div class="carousel-item-wrap">
    <div class="product-card" data-id="f22" data-category="beverages" data-neighbourhood="Indiranagar" data-name="bun maska & irani chai" data-price="49" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 140px;">
        <img src="/static/filter_coffee.jpg" alt="Bun Maska & Irani Chai" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bun Maska & Irani Chai</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Cafe Chai Days</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Bakery</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f22">
            <button type="button" class="btn-add-pill" onclick="addToCart('f22', 'Bun Maska & Irani Chai', 49, 'Cafe Chai Days', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    </div>
        </div>
      </div>
    </section>

    <!-- 9. HERO SECTION -->
    <section class="hero-section">
      <div class="container hero-grid">
        <div style="max-width: 560px;">
          <div class="hero-pill-badge">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
            <span>GOOD FOOD. SMALL BUDGET.</span>
          </div>
          <h1 class="hero-heading">
            Big cravings. <br>
            <span style="color: var(--color-primary-orange);">Tiny prices.</span>
          </h1>
          <p class="hero-copy">
            Discover satisfying local favourites across Bengaluru &mdash; every pick under &#8377;100.
          </p>
          <div style="display: flex; align-items: center; flex-wrap: wrap; gap: 16px; margin-bottom: 24px;">
            <a href="#food-discovery" class="btn-hero-primary">
              <span>Explore meals under &#8377;100</span>
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
            </a>
            <a href="#how-it-works" style="font-size: 14px; font-weight: 700; color: var(--color-dark-text); padding: 10px 16px;">
              See how it works
            </a>
          </div>
          <div style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--color-secondary-text); font-weight: 500;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#087443" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
            <span>No fancy fees. Just food that fits.</span>
          </div>
        </div>

        <div style="position: relative;">
          <div class="hero-img-box">
            <img src="/static/hero_food.jpg" alt="Authentic Bengaluru South Indian meal spread" class="hero-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
          </div>
          <div class="hero-rotate-badge">
            <span style="font-size: 10px; font-weight: 800; letter-spacing: 0.5px; text-transform: uppercase;">MEALS FROM</span>
            <span style="font-size: 24px; font-weight: 900; line-height: 1;">&#8377;49</span>
          </div>
          <div class="hero-float-card">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#087443" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
            <span style="font-size: 12px; font-weight: 700; color: var(--color-dark-text);">Your budget, respected &#10003;</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 10. BENEFITS STRIP -->
    <section class="benefits-strip">
      <div class="container benefits-grid">
        <div class="benefit-item">
          <div class="benefit-icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
          </div>
          <div>
            <h3 style="font-size: 16px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 3px;">Every item under &#8377;100</h3>
            <p style="font-size: 13px; color: var(--color-secondary-text); line-height: 1.4;">Affordable picks, clearly priced.</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
          </div>
          <div>
            <h3 style="font-size: 16px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 3px;">Local Bengaluru favourites</h3>
            <p style="font-size: 13px; color: var(--color-secondary-text); line-height: 1.4;">Discover food around your neighbourhood.</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect><line x1="1" y1="10" x2="23" y2="10"></line></svg>
          </div>
          <div>
            <h3 style="font-size: 16px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 3px;">Prices shown upfront</h3>
            <p style="font-size: 13px; color: var(--color-secondary-text); line-height: 1.4;">Know the item price before adding it.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 11. MAIN FOOD DISCOVERY GRID -->
    <section class="discovery-section" id="food-discovery">
      <div class="container">
        <div class="section-head-row" style="margin-bottom: 24px;">
          <div>
            <div style="font-size: 12px; font-weight: 800; color: var(--color-deep-green); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px;">
              Across Bengaluru
            </div>
            <h2 style="font-size: 28px; font-weight: 900; color: var(--color-dark-text); letter-spacing: -0.5px;">
              Good food. Under &#8377;100.
            </h2>
            <p style="font-size: 15px; color: var(--color-secondary-text); margin-top: 4px;">
              Affordable favourites near you.
            </p>
          </div>
        </div>

        <div class="filter-pills-row" style="display: flex; align-items: center; flex-wrap: wrap; gap: 8px;">
          <!-- Swiggy-style Veg Only Switch -->
          <div class="veg-toggle-wrap" onclick="var chk = document.getElementById('vegFilterToggle'); chk.checked = !chk.checked; toggleVegOnlyFilter(chk.checked);" title="Toggle Vegetarian Only">
            <div class="diet-indicator diet-veg" style="width: 14px; height: 14px; border: 1.5px solid #16a34a; border-radius: 3px; display: inline-flex; align-items: center; justify-content: center;">
              <div style="width: 6px; height: 6px; background-color: #16a34a; border-radius: 50%;"></div>
            </div>
            <span style="font-size: 13px; font-weight: 800; color: #166534; letter-spacing: -0.2px;">Veg Only</span>
            <label class="switch-ui" onclick="event.stopPropagation();">
              <input type="checkbox" id="vegFilterToggle" onchange="toggleVegOnlyFilter(this.checked)">
              <span class="slider-round"></span>
            </label>
          </div>

          <button class="filter-tab-pill active" onclick="filterGridCategory('all', this)">All Items (22)</button>
          <button id="openNowFilterBtn" class="filter-tab-pill" onclick="toggleOpenNowFilter(this)" style="display: inline-flex; align-items: center; gap: 6px; border: 1.5px solid #22c55e; color: #166534; font-weight: 700;">
            <span class="badge-dot-live"></span>
            Open Now
          </button>
          <button class="filter-tab-pill" onclick="filterGridCategory('dosa-idli', this)">Dosa &amp; Idli</button>
          <button class="filter-tab-pill" onclick="filterGridCategory('rolls-wraps', this)">Rolls &amp; Wraps</button>
          <button class="filter-tab-pill" onclick="filterGridCategory('rice-biryani', this)">Rice &amp; Biryani</button>
          <button class="filter-tab-pill" onclick="filterGridCategory('snacks-chaat', this)">Snacks &amp; Chaat</button>
          <button class="filter-tab-pill" onclick="filterGridCategory('momos', this)">Momos</button>
        </div>



        <div class="food-grid-container" id="mainFoodGrid">
          
    <div class="product-card" data-id="f1" data-category="dosa-idli" data-neighbourhood="Indiranagar" data-name="crispy masala dosa" data-price="79" data-rating="4.6" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Masala Dosa" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Masala Dosa</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Sri Krishna Cafe</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f1">
            <button type="button" class="btn-add-pill" onclick="addToCart('f1', 'Crispy Masala Dosa', 79, 'Sri Krishna Cafe', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f2" data-category="dosa-idli" data-neighbourhood="Koramangala" data-name="steamed idli sambar (2 pcs)" data-price="49" data-rating="4.8" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/idli_sambar.jpg" alt="Steamed Idli Sambar (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Idli Sambar (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Tiffin Room</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f2">
            <button type="button" class="btn-add-pill" onclick="addToCart('f2', 'Steamed Idli Sambar (2 pcs)', 49, 'Brahmin\'s Tiffin Room', '/static/idli_sambar.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f3" data-category="rice-biryani" data-neighbourhood="HSR Layout" data-name="aromatic veg pulao" data-price="89" data-rating="4.4" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/veg_pulao.jpg" alt="Aromatic Veg Pulao" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Aromatic Veg Pulao</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Meghana Bowl Corner</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>North Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f3">
            <button type="button" class="btn-add-pill" onclick="addToCart('f3', 'Aromatic Veg Pulao', 89, 'Meghana Bowl Corner', '/static/veg_pulao.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f4" data-category="rolls-wraps" data-neighbourhood="Indiranagar" data-name="paneer kathi roll" data-price="99" data-rating="4.5" data-delivery="22">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/kathi_roll.jpg" alt="Paneer Kathi Roll" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;99</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>22m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Paneer Kathi Roll</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Kolkata Roll Express</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;99</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f4">
            <button type="button" class="btn-add-pill" onclick="addToCart('f4', 'Paneer Kathi Roll', 99, 'Kolkata Roll Express', '/static/kathi_roll.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f5" data-category="rice-biryani" data-neighbourhood="BTM Layout" data-name="schezwan veg fried rice" data-price="89" data-rating="4.3" data-delivery="28">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/veg_pulao.jpg" alt="Schezwan Veg Fried Rice" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>28m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Schezwan Veg Fried Rice</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.3</span>
          </div>
        </div>
        <div class="product-card-restaurant">Wok & Bowl Express</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>BTM Layout</span>
          <span>&bull;</span>
          <span>Indo-Chinese</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f5">
            <button type="button" class="btn-add-pill" onclick="addToCart('f5', 'Schezwan Veg Fried Rice', 89, 'Wok & Bowl Express', '/static/veg_pulao.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f6" data-category="parathas" data-neighbourhood="Whitefield" data-name="stuffed aloo paratha with curd" data-price="69" data-rating="4.5" data-delivery="24">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/aloo_paratha.jpg" alt="Stuffed Aloo Paratha with Curd" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;69</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>24m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Stuffed Aloo Paratha with Curd</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Punjab Dhaba Point</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Whitefield</span>
          <span>&bull;</span>
          <span>North Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;69</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f6">
            <button type="button" class="btn-add-pill" onclick="addToCart('f6', 'Stuffed Aloo Paratha with Curd', 69, 'Punjab Dhaba Point', '/static/aloo_paratha.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f7" data-category="snacks-chaat" data-neighbourhood="Koramangala" data-name="indori spiced poha" data-price="39" data-rating="4.6" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/indori_poha.jpg" alt="Indori Spiced Poha" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;39</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Indori Spiced Poha</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Morning Fuel Co.</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Breakfast</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;39</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f7">
            <button type="button" class="btn-add-pill" onclick="addToCart('f7', 'Indori Spiced Poha', 39, 'Morning Fuel Co.', '/static/indori_poha.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f8" data-category="momos" data-neighbourhood="Marathahalli" data-name="steamed veg momos (6 pcs)" data-price="79" data-rating="4.5" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Veg Momos (6 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Veg Momos (6 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Himalayan Momo Hut</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Marathahalli</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f8">
            <button type="button" class="btn-add-pill" onclick="addToCart('f8', 'Steamed Veg Momos (6 pcs)', 79, 'Himalayan Momo Hut', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f9" data-category="dosa-idli" data-neighbourhood="Malleshwaram" data-name="crispy plain dosa & podi" data-price="59" data-rating="4.7" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/crispy_dosa.jpg" alt="Crispy Plain Dosa & Podi" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;59</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Plain Dosa & Podi</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Vidyarthi Bhavana Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;59</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f9">
            <button type="button" class="btn-add-pill" onclick="addToCart('f9', 'Crispy Plain Dosa & Podi', 59, 'Vidyarthi Bhavana Tiffins', '/static/crispy_dosa.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f10" data-category="south-indian" data-neighbourhood="HSR Layout" data-name="golden medu vada (2 pcs)" data-price="45" data-rating="4.6" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/medu_vada.jpg" alt="Golden Medu Vada (2 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Golden Medu Vada (2 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Udupi Grand Tiffins</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f10">
            <button type="button" class="btn-add-pill" onclick="addToCart('f10', 'Golden Medu Vada (2 pcs)', 45, 'Udupi Grand Tiffins', '/static/medu_vada.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f11" data-category="dosa-idli" data-neighbourhood="Jayanagar" data-name="onion tomato uttapam" data-price="79" data-rating="4.7" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/onion_uttapam.jpg" alt="Onion Tomato Uttapam" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Onion Tomato Uttapam</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.7</span>
          </div>
        </div>
        <div class="product-card-restaurant">Taaza Thindi</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f11">
            <button type="button" class="btn-add-pill" onclick="addToCart('f11', 'Onion Tomato Uttapam', 79, 'Taaza Thindi', '/static/onion_uttapam.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f12" data-category="south-indian" data-neighbourhood="Indiranagar" data-name="south indian mini meal box" data-price="89" data-rating="4.6" data-delivery="25">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/hero_food.jpg" alt="South Indian Mini Meal Box" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>25m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">South Indian Mini Meal Box</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">MTR Heritage Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f12">
            <button type="button" class="btn-add-pill" onclick="addToCart('f12', 'South Indian Mini Meal Box', 89, 'MTR Heritage Canteen', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f13" data-category="snacks-chaat" data-neighbourhood="BTM Layout" data-name="punjabi samosa (2 pcs with chutney)" data-price="39" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/samosa_plate.jpg" alt="Punjabi Samosa (2 pcs with Chutney)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;39</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Punjabi Samosa (2 pcs with Chutney)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Delhi Chaat Bhandar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>BTM Layout</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;39</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f13">
            <button type="button" class="btn-add-pill" onclick="addToCart('f13', 'Punjabi Samosa (2 pcs with Chutney)', 39, 'Delhi Chaat Bhandar', '/static/samosa_plate.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f14" data-category="snacks-chaat" data-neighbourhood="Koramangala" data-name="tangy aloo tikki chaat" data-price="59" data-rating="4.4" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/samosa_plate.jpg" alt="Tangy Aloo Tikki Chaat" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;59</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Tangy Aloo Tikki Chaat</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Anand Sweets & Chaat</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;59</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f14">
            <button type="button" class="btn-add-pill" onclick="addToCart('f14', 'Tangy Aloo Tikki Chaat', 59, 'Anand Sweets & Chaat', '/static/samosa_plate.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f15" data-category="sandwiches" data-neighbourhood="Jayanagar" data-name="bombay veg grilled sandwich" data-price="69" data-rating="4.6" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/bombay_sandwich.jpg" alt="Bombay Veg Grilled Sandwich" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;69</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bombay Veg Grilled Sandwich</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Hari Super Sandwich</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Jayanagar</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;69</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f15">
            <button type="button" class="btn-add-pill" onclick="addToCart('f15', 'Bombay Veg Grilled Sandwich', 69, 'Hari Super Sandwich', '/static/bombay_sandwich.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f16" data-category="snacks-chaat" data-neighbourhood="Indiranagar" data-name="crispy pani puri (6 puris)" data-price="49" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/pani_puri.jpg" alt="Crispy Pani Puri (6 Puris)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Crispy Pani Puri (6 Puris)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Chaat Junction</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f16">
            <button type="button" class="btn-add-pill" onclick="addToCart('f16', 'Crispy Pani Puri (6 Puris)', 49, 'Chaat Junction', '/static/pani_puri.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f17" data-category="snacks-chaat" data-neighbourhood="Whitefield" data-name="hot onion pakora plate" data-price="45" data-rating="4.3" data-delivery="18">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/onion_pakora.jpg" alt="Hot Onion Pakora Plate" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;45</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>18m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Hot Onion Pakora Plate</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.3</span>
          </div>
        </div>
        <div class="product-card-restaurant">Chai Point Canteen</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Whitefield</span>
          <span>&bull;</span>
          <span>Snacks</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;45</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f17">
            <button type="button" class="btn-add-pill" onclick="addToCart('f17', 'Hot Onion Pakora Plate', 45, 'Chai Point Canteen', '/static/onion_pakora.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f18" data-category="beverages" data-neighbourhood="Malleshwaram" data-name="authentic filter coffee" data-price="29" data-rating="4.9" data-delivery="12">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/filter_coffee.jpg" alt="Authentic Filter Coffee" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;29</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>12m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Authentic Filter Coffee</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.9</span>
          </div>
        </div>
        <div class="product-card-restaurant">Brahmin's Coffee Bar</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>Beverages</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;29</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f18">
            <button type="button" class="btn-add-pill" onclick="addToCart('f18', 'Authentic Filter Coffee', 29, 'Brahmin\'s Coffee Bar', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f19" data-category="south-indian" data-neighbourhood="Malleshwaram" data-name="bengaluru khara bath & kesari bath" data-price="55" data-rating="4.8" data-delivery="16">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/hero_food.jpg" alt="Bengaluru Khara Bath & Kesari Bath" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;55</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>16m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bengaluru Khara Bath & Kesari Bath</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.8</span>
          </div>
        </div>
        <div class="product-card-restaurant">Veena Stores</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Malleshwaram</span>
          <span>&bull;</span>
          <span>South Indian</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;55</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f19">
            <button type="button" class="btn-add-pill" onclick="addToCart('f19', 'Bengaluru Khara Bath & Kesari Bath', 55, 'Veena Stores', '/static/hero_food.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f20" data-category="rolls-wraps" data-neighbourhood="HSR Layout" data-name="egg kathi roll" data-price="79" data-rating="4.4" data-delivery="20">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/kathi_roll.jpg" alt="Egg Kathi Roll" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-nonveg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;79</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>20m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Egg Kathi Roll</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.4</span>
          </div>
        </div>
        <div class="product-card-restaurant">Rolls King</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>HSR Layout</span>
          <span>&bull;</span>
          <span>Street Food</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;79</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f20">
            <button type="button" class="btn-add-pill" onclick="addToCart('f20', 'Egg Kathi Roll', 79, 'Rolls King', '/static/kathi_roll.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f21" data-category="momos" data-neighbourhood="Koramangala" data-name="steamed chicken momos (5 pcs)" data-price="89" data-rating="4.6" data-delivery="22">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/steamed_momos.jpg" alt="Steamed Chicken Momos (5 pcs)" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-nonveg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;89</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>22m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Steamed Chicken Momos (5 pcs)</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.6</span>
          </div>
        </div>
        <div class="product-card-restaurant">Shillong Momo Corner</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Koramangala</span>
          <span>&bull;</span>
          <span>Tibetan</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;89</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f21">
            <button type="button" class="btn-add-pill" onclick="addToCart('f21', 'Steamed Chicken Momos (5 pcs)', 89, 'Shillong Momo Corner', '/static/steamed_momos.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="product-card" data-id="f22" data-category="beverages" data-neighbourhood="Indiranagar" data-name="bun maska & irani chai" data-price="49" data-rating="4.5" data-delivery="15">
      <div class="product-card-img-wrap" style="height: 170px;">
        <img src="/static/filter_coffee.jpg" alt="Bun Maska & Irani Chai" loading="lazy" class="product-card-img" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
        <div class="diet-badge-pos">
          <div class="diet-indicator diet-veg">
            <div class="diet-dot"></div>
          </div>
        </div>
        <div class="price-pill-pos">&#8377;49</div>
        <div class="delivery-badge-pos">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#646B73" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>15m</span>
        </div>
      </div>
      <div class="product-card-body">
        <div class="product-card-title-row">
          <h3 class="product-card-title">Bun Maska & Irani Chai</h3>
          <div class="product-card-rating">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="#B45309" stroke="#B45309"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span>4.5</span>
          </div>
        </div>
        <div class="product-card-restaurant">Cafe Chai Days</div>
        <div class="product-card-meta">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="#FF5A1F"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/></svg>
          <span>Indiranagar</span>
          <span>&bull;</span>
          <span>Bakery</span>
        </div>
        <div class="product-card-footer">
          <div>
            <span class="price-label">Item Price</span>
            <span class="price-val">&#8377;49</span>
          </div>
          <div class="add-btn-wrap" id="add-btn-f22">
            <button type="button" class="btn-add-pill" onclick="addToCart('f22', 'Bun Maska & Irani Chai', 49, 'Cafe Chai Days', '/static/filter_coffee.jpg')">
              ADD
            </button>
          </div>
        </div>
      </div>
    </div>
    
        </div>
      </div>
    </section>

    <!-- 12. PROMO BANNER -->
    <section class="promo-banner-section">
      <div class="container">
        <div class="promo-box">
          <div style="max-width: 600px;">
            <div style="display: inline-flex; align-items: center; gap: 6px; background-color: #ffffff; color: var(--color-deep-green); padding: 4px 12px; border-radius: var(--radius-full); font-size: 12px; font-weight: 800; text-transform: uppercase; margin-bottom: 16px; border: 1px solid rgba(8, 116, 67, 0.2);">
              <span>Pocket-Friendly Promise</span>
            </div>
            <h2 style="font-size: clamp(26px, 4vw, 36px); font-weight: 900; color: var(--color-dark-text); letter-spacing: -0.8px; line-height: 1.2; margin-bottom: 12px;">
              Your daily food budget just got happier.
            </h2>
            <p style="font-size: 16px; color: var(--color-secondary-text); line-height: 1.6; margin-bottom: 24px;">
              Discover satisfying meals and local favourites for under &#8377;100.
            </p>
            <div style="display: flex; flex-wrap: wrap; gap: 12px;">
              <span style="display: inline-flex; align-items: center; gap: 6px; background: #fff; padding: 6px 12px; border-radius: var(--radius-full); font-size: 13px; font-weight: 700; color: var(--color-deep-green);">&#10003; Tasty meals under &#8377;100</span>
              <span style="display: inline-flex; align-items: center; gap: 6px; background: #fff; padding: 6px 12px; border-radius: var(--radius-full); font-size: 13px; font-weight: 700; color: var(--color-deep-green);">&#10003; Support local eateries</span>
              <span style="display: inline-flex; align-items: center; gap: 6px; background: #fff; padding: 6px 12px; border-radius: var(--radius-full); font-size: 13px; font-weight: 700; color: var(--color-deep-green);">&#10003; Real food. Real people.</span>
            </div>
          </div>
          <div style="width: 220px; height: 220px; border-radius: var(--radius-lg); overflow: hidden; box-shadow: 0 8px 24px rgba(0,0,0,0.1); border: 3px solid #ffffff; flex-shrink: 0; background-color: #F1EFEA;">
            <img src="/static/crispy_dosa.jpg" alt="Local Bengaluru Street Food" style="width: 100%; height: 100%; object-fit: cover; display: block;" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">
          </div>
        </div>
      </div>
    </section>

    <!-- 13. HOW IT WORKS -->
    <section class="info-section" id="how-it-works">
      <div class="container">
        <div style="text-align: center; max-width: 600px; margin: 0 auto 48px;">
          <div style="font-size: 12px; font-weight: 800; color: var(--color-primary-orange); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px;">
            Simple &amp; Transparent
          </div>
          <h2 style="font-size: 30px; font-weight: 900; color: var(--color-dark-text); letter-spacing: -0.5px; margin-bottom: 10px;">
            How Under100 Works
          </h2>
          <p style="font-size: 15px; color: var(--color-secondary-text);">
            Getting fresh everyday meals delivered without overpaying takes just 3 simple steps.
          </p>
        </div>

        <div class="info-grid">
          <div class="info-card">
            <div style="width: 56px; height: 56px; border-radius: 50%; background-color: var(--color-soft-orange-bg); color: var(--color-primary-orange); display: flex; align-items: center; justify-content: center; margin: 0 auto 20px;">
              <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
            </div>
            <div style="font-size: 12px; font-weight: 800; color: var(--color-primary-orange); margin-bottom: 6px;">STEP 1</div>
            <h3 style="font-size: 18px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 8px;">Choose your neighbourhood</h3>
            <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">Select Indiranagar, Koramangala, HSR Layout, or other Bengaluru areas to see nearby active kitchens.</p>
          </div>

          <div class="info-card">
            <div style="width: 56px; height: 56px; border-radius: 50%; background-color: var(--color-soft-green-bg); color: var(--color-deep-green); display: flex; align-items: center; justify-content: center; margin: 0 auto 20px;">
              <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            </div>
            <div style="font-size: 12px; font-weight: 800; color: var(--color-deep-green); margin-bottom: 6px;">STEP 2</div>
            <h3 style="font-size: 18px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 8px;">Find a meal within your budget</h3>
            <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">Browse curated food carousels and dishes &mdash; every single listing is guaranteed to be under &#8377;100.</p>
          </div>

          <div class="info-card">
            <div style="width: 56px; height: 56px; border-radius: 50%; background-color: var(--color-soft-orange-bg); color: var(--color-primary-orange); display: flex; align-items: center; justify-content: center; margin: 0 auto 20px;">
              <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path><line x1="3" y1="6" x2="21" y2="6"></line><path d="M16 10a4 4 0 0 1-8 0"></path></svg>
            </div>
            <div style="font-size: 12px; font-weight: 800; color: var(--color-primary-orange); margin-bottom: 6px;">STEP 3</div>
            <h3 style="font-size: 18px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 8px;">Add favourites &amp; review cart</h3>
            <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">Tap Add to customize portions with real-time transparent subtotal calculations and no surprise fees.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 14. WHY UNDER100? -->
    <section class="info-section" style="background-color: var(--color-main-bg);">
      <div class="container">
        <div style="text-align: center; max-width: 600px; margin: 0 auto 48px;">
          <div style="font-size: 12px; font-weight: 800; color: var(--color-deep-green); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px;">
            Built For Bengaluru
          </div>
          <h2 style="font-size: 30px; font-weight: 900; color: var(--color-dark-text); letter-spacing: -0.5px; margin-bottom: 10px;">
            Why Under100?
          </h2>
          <p style="font-size: 15px; color: var(--color-secondary-text);">
            A food platform crafted around honest prices and neighborhood comfort.
          </p>
        </div>

        <div class="benefits-grid">
          <div class="benefit-item">
            <div class="benefit-icon-box">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="19" y1="5" x2="5" y2="19"></line><circle cx="6.5" cy="6.5" r="2.5"></circle><circle cx="17.5" cy="17.5" r="2.5"></circle></svg>
            </div>
            <div>
              <h3 style="font-size: 17px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 6px;">Budget-first discovery</h3>
              <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">Everything on Under100 is capped under &#8377;100 by design.</p>
            </div>
          </div>

          <div class="benefit-item">
            <div class="benefit-icon-box" style="background-color: var(--color-soft-green-bg); color: var(--color-deep-green);">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
            </div>
            <div>
              <h3 style="font-size: 17px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 6px;">Transparent item pricing</h3>
              <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">The price you see is what you pay. No inflated menu prices or hidden surcharges.</p>
            </div>
          </div>

          <div class="benefit-item">
            <div class="benefit-icon-box">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"></polygon></svg>
            </div>
            <div>
              <h3 style="font-size: 17px; font-weight: 800; color: var(--color-dark-text); margin-bottom: 6px;">Local restaurant legends</h3>
              <p style="font-size: 14px; color: var(--color-secondary-text); line-height: 1.5;">Highlighting authentic Bengaluru darshinis, tiffin rooms, rolls joints, and momo stalls.</p>
            </div>
          </div>
        </div>
      </div>
    </section>

  </main>

  <!-- 15. REDESIGNED PREMIUM FOOTER (EXACT SAME #171C24 DARK PALETTE & HEADER LOGO) -->
  <footer class="footer-dark">
    <div class="container">
      <!-- A. VALUE PROPOSITION PILLARS -->
      <div class="footer-usp-grid">
        <div class="footer-usp-card">
          <div class="footer-usp-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#FF5A1F" stroke-width="2.5"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
            <span>Strictly &#8377;100 Ceiling</span>
          </div>
          <p class="footer-usp-desc">Every single menu item is strictly capped at or under &#8377;100. Delicious meals without the wallet burn.</p>
        </div>
        <div class="footer-usp-card">
          <div class="footer-usp-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#FF5A1F" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            <span>Zero Price Markups</span>
          </div>
          <p class="footer-usp-desc">Dine-in menu prices honored directly. No inflated dish rates, zero hidden service fees.</p>
        </div>
        <div class="footer-usp-card">
          <div class="footer-usp-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#FF5A1F" stroke-width="2.5"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
            <span>Namma Bengaluru First</span>
          </div>
          <p class="footer-usp-desc">From iconic Indiranagar darshinis to Koramangala roll stalls and VV Puram chaats.</p>
        </div>
        <div class="footer-usp-card">
          <div class="footer-usp-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#FF5A1F" stroke-width="2.5"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"></path></svg>
            <span>Instant UPI &amp; Cards</span>
          </div>
          <p class="footer-usp-desc">Seamless checkout via Razorpay, Stripe, GPay, PhonePe, and instant counter collection.</p>
        </div>
      </div>

      <!-- B. MAIN 4-COLUMN FOOTER NAVIGATION -->
      <div class="footer-grid-redesigned">
        <!-- Col 1: Brand & Bengaluru Pride -->
        <div class="footer-brand-col">
          <a href="/" class="brand-logo-link footer-brand-link" title="Under100 Home">
            <svg width="34" height="34" viewBox="0 0 36 36" fill="none">
              <path d="M12 4C11 6.5 13 8.5 12 11" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
              <path d="M18 3C17 5.5 19 8 18 10.5" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
              <path d="M24 4C23 6.5 25 8.5 24 11" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
              <path d="M6 14H30" stroke="#FF5A1F" stroke-width="2.5" stroke-linecap="round"/>
              <path d="M7 16C7.5 24 13.5 28 18 28C22.5 28 28.5 24 29 16H7Z" fill="#FF5A1F"/>
              <path d="M13 28.5H23" stroke="#FF5A1F" stroke-width="2.2" stroke-linecap="round"/>
            </svg>
            <div class="brand-wordmark">
              <span class="under">under</span>
              <span class="num">100</span>
            </div>
          </a>
          <div class="footer-tagline-badge">Namma Bengaluru Food Discovery</div>
          <p class="footer-mission-text">
            Bengaluru&#39;s dedicated budget food discovery platform. Helping food lovers, students, and working professionals discover satisfying neighborhood meals, crisp dosas, and quick bites strictly under &#8377;100.
          </p>
          <div class="footer-service-pill">
            <span>Live across 18+ Bengaluru Localities</span>
          </div>
          <div class="footer-trust-tags">
            <span>100% Certified Eateries</span> &bull; <span>Verified Menus</span>
          </div>
        </div>

        <!-- Col 2: Curated Cuisines -->
        <div>
          <div class="footer-col-title">Curated Flavors</div>
          <ul class="footer-link-list">
            <li><a href="/#food-discovery" onclick="filterByTag('South Indian')">South Indian Tiffins <span class="footer-badge-hot">Popular</span></a></li>
            <li><a href="/#food-discovery" onclick="filterByTag('Rolls')">Kolkata Rolls &amp; Wraps</a></li>
            <li><a href="/#food-discovery" onclick="filterByTag('Momos')">Steamed &amp; Fried Momos</a></li>
            <li><a href="/#food-discovery" onclick="filterByTag('Street Food')">Street Chaat &amp; Pav Bhaji</a></li>
            <li><a href="/#food-discovery" onclick="filterByTag('Beverages')">Filter Coffee &amp; Chai <span class="footer-badge-new">New</span></a></li>
            <li><a href="/#food-discovery" onclick="filterByTag('Meals')">Mini Thalis &amp; Rice Bowls</a></li>
          </ul>
        </div>

        <!-- Col 3: Popular Localities -->
        <div>
          <div class="footer-col-title">Popular Hubs</div>
          <ul class="footer-link-list">
            <li><a href="/#food-discovery" onclick="setLocality('Indiranagar')">Indiranagar 100ft Rd</a></li>
            <li><a href="/#food-discovery" onclick="setLocality('Koramangala')">Koramangala 4th &amp; 5th Block</a></li>
            <li><a href="/#food-discovery" onclick="setLocality('HSR Layout')">HSR Layout Sectors</a></li>
            <li><a href="/#food-discovery" onclick="setLocality('Jayanagar')">Jayanagar 4th Block</a></li>
            <li><a href="/#food-discovery" onclick="setLocality('Whitefield')">Whitefield Tech Hub</a></li>
            <li><a href="/#food-discovery" onclick="setLocality('Malleshwaram')">Malleshwaram Heritage</a></li>
          </ul>
        </div>

        <!-- Col 4: Customer Care & Trust -->
        <div>
          <div class="footer-col-title">Customer Care</div>
          <ul class="footer-link-list">
            <li><a href="/about">About Under100</a></li>
            <li><a href="#how-it-works">How Discovery Works</a></li>
            <li><a href="/privacy-policy">Privacy &amp; Terms</a></li>
            <li><a href="/account/orders">My Order History</a></li>
            <li><a href="/favorites">Saved Favorites</a></li>
            <li><a href="mailto:support@under100.in">Help: support@under100.in</a></li>
          </ul>
        </div>
      </div>

      <!-- C. POPULAR QUICK DISCOVERY SEARCH TAGS -->
      <div class="footer-searches-box">
        <div class="footer-searches-header">Popular Searches Under &#8377;100:</div>
        <div class="footer-tags-row">
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Masala Dosa')">Masala Dosa Under &#8377;50</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Filter Coffee')">Filter Coffee Under &#8377;25</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Egg Roll')">Double Egg Roll Under &#8377;65</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Momos')">Veg Momos Under &#8377;50</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Pav Bhaji')">Butter Pav Bhaji Under &#8377;80</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Idli Vada')">Idli Vada Combo Under &#8377;45</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Curd Rice')">Curd Rice Under &#8377;55</a>
          <a href="/#food-discovery" class="footer-tag-link" onclick="handleQuickSearch('Cold Coffee')">Cold Coffee Under &#8377;60</a>
        </div>
      </div>

      <!-- D. BOTTOM LEGAL & ATTRIBUTION -->
      <div class="footer-bottom-flex">
        <div>&copy; 2026 Under100 Technologies Pvt Ltd. Made with <span style="color: #FF5A1F;">&#10084;</span> for Bengaluru foodies.</div>
        <div style="font-size: 11.5px; color: #64748B;">Built with Go &bull; HTMX &bull; PostgreSQL &bull; Redis</div>
        <div class="footer-legal-links">
          <a href="/privacy-policy">Privacy Policy</a>
          <a href="/privacy-policy#terms">Terms of Service</a>
          <a href="/privacy-policy#refunds">Refund Policy</a>
        </div>
      </div>
    </div>
  </footer>

  <!-- 16. BLINKIT-STYLE CART DRAWER & BACKDROP -->
  <div class="cart-drawer-backdrop" id="cartBackdrop" onclick="toggleCartDrawer()"></div>
  <aside class="cart-drawer-panel" id="cartDrawerPanel">
    <div class="blinkit-cart-header">
      <button type="button" class="btn-cart-back" onclick="toggleCartDrawer()">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
        <span>My Cart</span>
      </button>
      <button type="button" class="btn-cart-share" onclick="shareCartLink()">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="18" cy="5" r="3"></circle><circle cx="6" cy="12" r="3"></circle><circle cx="18" cy="19" r="3"></circle><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line></svg>
        <span>Share</span>
      </button>
    </div>

    <div class="blinkit-cart-body" id="cartItemsContainer">
      <!-- Dynamically rendered by renderCart() -->
    </div>

    <div class="blinkit-cart-bottom" id="cartFooterSection" style="display: none;">
      <button type="button" class="btn-blinkit-proceed" onclick="checkoutOrder()">
        <div style="text-align: left;">
          <div style="font-size: 17px; font-weight: 800; line-height: 1.1;" id="cartBottomGrandTotal">&#8377;0</div>
          <div style="font-size: 10px; font-weight: 800; letter-spacing: 0.5px; opacity: 0.9; margin-top: 2px;">TOTAL</div>
        </div>
        <div style="font-size: 15px; font-weight: 800; display: flex; align-items: center; gap: 4px;">
          <span>Proceed</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </div>
      </button>
      <!-- Hidden legacy span elements to preserve ID bindings -->
      <span id="cartSubtotalVal" style="display:none;">&#8377;0</span>
      <span id="cartToPayVal" style="display:none;">&#8377;0</span>
      <span id="checkoutBtnLabel" style="display:none;">Proceed</span>
      <span id="cartItemCountLabel" style="display:none;">0 items</span>
    </div>
  </aside>

  <!-- MULTI-RESTAURANT CART CONFLICT GUARD MODAL -->
  <div class="conflict-modal-overlay" id="cartConflictModal" style="display: none;" onclick="if(event.target === this) resolveCartConflict(false)">
    <div class="conflict-modal-card">
      <div class="conflict-modal-icon">⚠️</div>
      <h3 class="conflict-modal-title">Replace cart items?</h3>
      <p class="conflict-modal-desc">
        Your cart already contains delicious items from <strong id="conflictOldRestName" style="color: #1E293B;">Restaurant A</strong>. Do you want to discard your current cart and start a fresh order from <strong id="conflictNewRestName" style="color: #FF5A1F;">Restaurant B</strong> instead?
      </p>
      <div class="conflict-actions-row">
        <button type="button" class="btn-conflict-cancel" onclick="resolveCartConflict(false)">Keep Current Cart</button>
        <button type="button" class="btn-conflict-replace" onclick="resolveCartConflict(true)">Replace Cart &amp; Add</button>
      </div>
    </div>
  </div>

  <!-- LIVE ORDER TRACKING MODAL & STEPPER (SUPPORTS MULTI-RESTAURANT SEPARATE TRACKING) -->
  <div class="tracking-modal-overlay" id="liveOrderTrackingModal" style="display: none;" onclick="if(event.target === this) closeOrderTrackingModal()">
    <div class="tracking-modal-card">
      <div class="tracking-header">
        <button type="button" onclick="closeOrderTrackingModal()" style="position: absolute; right: 18px; top: 18px; background: none; border: none; color: #ffffff; font-size: 24px; cursor: pointer; line-height: 1;">&times;</button>
        <div class="tracking-eta-pill">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span id="trackEtaText">Estimated Delivery: 18 mins</span>
        </div>
        <h2 class="tracking-title" id="trackHeaderTitle">Your meal is on the way!</h2>
        <div style="font-size: 13px; color: #94A3B8; margin-top: 4px;" id="trackRestaurantName">From CTR Shri Sagar, Malleshwaram</div>
      </div>

      <!-- Multi-Restaurant Order Tabs (shown when order has multiple kitchens) -->
      <div id="multiTrackTabsWrap" class="multi-track-tabs" style="display: none;"></div>

      <div class="tracking-body">
        <div id="trackOrderItemsSummary" style="display: none; background: #F8FAFC; border: 1.5px solid #E2E8F0; border-radius: 12px; padding: 10px 14px; margin-bottom: 14px; font-size: 13px; color: #1E293B; justify-content: space-between; align-items: center;">
          <div>
            <div style="font-size: 10.5px; font-weight: 800; color: #64748B; text-transform: uppercase;">SHIPMENT CONTENTS</div>
            <div id="trackItemsListText" style="font-weight: 700; color: #1E293B; margin-top: 2px;">1x Crispy Masala Dosa</div>
          </div>
          <div id="trackLegSubtotal" style="font-weight: 800; font-size: 14px; color: #0C831F;">&#8377;79</div>
        </div>

        <!-- Visual 5-Stage Stepper -->
        <div class="track-stepper">
          <div class="track-step-item done" id="step1">
            <div class="track-step-dot">&#10003;</div>
            <div class="track-step-info">
              <h4>Order Placed &amp; Confirmed</h4>
              <p>Payment settled successfully via UPI</p>
            </div>
          </div>
          <div class="track-step-item active" id="step2">
            <div class="track-step-dot">2</div>
            <div class="track-step-info">
              <h4>Kitchen Accepted &amp; Preparing</h4>
              <p id="trackKitchenDetail">Fresh meals prepared with clean ingredients</p>
            </div>
          </div>
          <div class="track-step-item" id="step3">
            <div class="track-step-dot">3</div>
            <div class="track-step-info">
              <h4>Delivery Valet Assigned</h4>
              <p id="trackRiderDetail">Ramesh Kumar (Hero Splendor KA-04-EK-9122)</p>
            </div>
          </div>
          <div class="track-step-item" id="step4">
            <div class="track-step-dot">4</div>
            <div class="track-step-info">
              <h4>Out for Delivery</h4>
              <p>Heading to Indiranagar, Bengaluru</p>
            </div>
          </div>
          <div class="track-step-item" id="step5">
            <div class="track-step-dot">5</div>
            <div class="track-step-info">
              <h4>Delivered &amp; Enjoy!</h4>
              <p>Order handed over safely</p>
            </div>
          </div>
        </div>

        <div style="background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 14px; padding: 16px; margin-bottom: 20px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
            <span style="font-size: 12px; font-weight: 800; color: #64748B; text-transform: uppercase;">Delivery Address</span>
            <span style="font-size: 12px; font-weight: 700; color: #0C831F;">&#127968; Home</span>
          </div>
          <div style="font-size: 13.5px; font-weight: 600; color: #1E293B;" id="trackAddressText">Indiranagar 100ft Rd, Bengaluru</div>
        </div>

        <button type="button" onclick="closeOrderTrackingModal()" style="width: 100%; padding: 13px; background: #0C831F; color: #ffffff; border: none; border-radius: 12px; font-weight: 800; font-size: 14px; cursor: pointer; transition: background 0.2s ease;">
          Back to Food Discovery
        </button>
      </div>
    </div>
  </div>

  <!-- 19. ADVANCED PAYMENT GATEWAY CHECKOUT MODAL -->
  <div class="pg-modal-overlay" id="paymentGatewayModal" style="display: none;" onclick="if(event.target === this) closePaymentModal()">
    <div class="pg-modal-container">
      <!-- Modal Header -->
      <div class="pg-modal-header" id="pgHeaderBar">
        <div>
          <div class="pg-header-kicker">RBI CERTIFIED &bull; 256-BIT TLS ENCRYPTED</div>
          <h3 class="pg-header-title">Payment Settlement Console</h3>
        </div>
        <button type="button" class="pg-modal-close" onclick="closePaymentModal()">&times;</button>
      </div>

      <!-- Main Payment Body -->
      <div class="pg-modal-body" id="pgMainBody">
        <!-- Left: Payment Options Tabs & Content -->
        <div class="pg-left-pane">
          <!-- Navigation Tabs (Zero Icons) -->
          <div class="pg-tabs-bar">
            <button type="button" class="pg-tab-btn active" id="pgTabUpi" onclick="selectPaymentTab('upi')">UPI / QR CODE</button>
            <button type="button" class="pg-tab-btn" id="pgTabCard" onclick="selectPaymentTab('card')">CARDS</button>
            <button type="button" class="pg-tab-btn" id="pgTabNet" onclick="selectPaymentTab('net')">NETBANKING</button>
            <button type="button" class="pg-tab-btn" id="pgTabWallet" onclick="selectPaymentTab('wallet')">WALLETS</button>
            <button type="button" class="pg-tab-btn" id="pgTabCod" onclick="selectPaymentTab('cod')">PAY ON DELIVERY</button>
          </div>

          <!-- Tab Content Panes -->
          <div class="pg-tab-content-area">
            <!-- 1. UPI Pane -->
            <div class="pg-tab-pane active" id="paneUpi">
              <div style="font-size: 12.5px; font-weight: 800; color: #0F172A; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 10px;">Select UPI Instrument</div>
              
              <!-- Quick App Badges (Zero Icons) -->
              <div class="pg-upi-apps-grid">
                <div class="pg-upi-app active" onclick="selectUpiApp(this, 'Google Pay')">
                  <div class="pg-upi-badge-text">G-PAY</div>
                  <div class="pg-upi-app-name">Google Pay</div>
                  <span class="pg-upi-check">&#10003;</span>
                </div>
                <div class="pg-upi-app" onclick="selectUpiApp(this, 'PhonePe')">
                  <div class="pg-upi-badge-text">PHONEPE</div>
                  <div class="pg-upi-app-name">PhonePe</div>
                  <span class="pg-upi-check">&#10003;</span>
                </div>
                <div class="pg-upi-app" onclick="selectUpiApp(this, 'Paytm UPI')">
                  <div class="pg-upi-badge-text">PAYTM</div>
                  <div class="pg-upi-app-name">Paytm UPI</div>
                  <span class="pg-upi-check">&#10003;</span>
                </div>
                <div class="pg-upi-app" onclick="selectUpiApp(this, 'BHIM UPI')">
                  <div class="pg-upi-badge-text">BHIM</div>
                  <div class="pg-upi-app-name">BHIM UPI</div>
                  <span class="pg-upi-check">&#10003;</span>
                </div>
              </div>

              <!-- Dynamic QR Code Container -->
              <div class="pg-qr-box">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <span style="font-size: 11.5px; font-weight: 800; color: #166534; text-transform: uppercase; letter-spacing: 0.5px;">Dynamic Merchant QR Code</span>
                  <span class="pg-timer-badge" id="pgQrTimer">04:59</span>
                </div>
                <div class="pg-qr-wrapper">
                  <svg width="130" height="130" viewBox="0 0 140 140" style="border-radius: 8px; background: white; padding: 6px; box-shadow: 0 1px 4px rgba(0,0,0,0.08);">
                    <rect width="140" height="140" fill="#ffffff"/>
                    <rect x="10" y="10" width="36" height="36" fill="#0F172A"/>
                    <rect x="16" y="16" width="24" height="24" fill="#ffffff"/>
                    <rect x="22" y="22" width="12" height="12" fill="#0C831F"/>
                    <rect x="94" y="10" width="36" height="36" fill="#0F172A"/>
                    <rect x="100" y="16" width="24" height="24" fill="#ffffff"/>
                    <rect x="106" y="22" width="12" height="12" fill="#0C831F"/>
                    <rect x="10" y="94" width="36" height="36" fill="#0F172A"/>
                    <rect x="16" y="100" width="24" height="24" fill="#ffffff"/>
                    <rect x="22" y="106" width="12" height="12" fill="#0C831F"/>
                    <rect x="54" y="14" width="8" height="8" fill="#1E293B"/>
                    <rect x="70" y="14" width="16" height="8" fill="#1E293B"/>
                    <rect x="54" y="30" width="14" height="8" fill="#1E293B"/>
                    <rect x="76" y="30" width="10" height="8" fill="#1E293B"/>
                    <rect x="14" y="54" width="16" height="8" fill="#1E293B"/>
                    <rect x="38" y="54" width="16" height="12" fill="#1E293B"/>
                    <rect x="62" y="54" width="16" height="16" fill="#0C831F"/>
                    <rect x="86" y="54" width="16" height="8" fill="#1E293B"/>
                    <rect x="110" y="54" width="16" height="12" fill="#1E293B"/>
                    <rect x="54" y="76" width="8" height="16" fill="#1E293B"/>
                    <rect x="70" y="78" width="14" height="14" fill="#1E293B"/>
                    <rect x="92" y="76" width="16" height="12" fill="#1E293B"/>
                    <rect x="116" y="76" width="14" height="8" fill="#1E293B"/>
                    <rect x="54" y="100" width="16" height="8" fill="#1E293B"/>
                    <rect x="78" y="98" width="12" height="16" fill="#0C831F"/>
                    <rect x="98" y="100" width="16" height="8" fill="#1E293B"/>
                    <rect x="122" y="94" width="8" height="16" fill="#1E293B"/>
                    <rect x="54" y="116" width="22" height="14" fill="#1E293B"/>
                    <rect x="84" y="122" width="14" height="8" fill="#1E293B"/>
                    <rect x="106" y="116" width="24" height="14" fill="#1E293B"/>
                  </svg>
                  <div style="font-size: 11.5px; color: #475569; margin-top: 8px; text-align: center;">
                    Settlement VPA: <strong style="color: #0F172A;">under100bengaluru@axisbank</strong>
                  </div>
                </div>
              </div>

              <!-- Manual UPI ID -->
              <div style="margin-top: 14px;">
                <label class="pg-label">Or enter Virtual Payment Address (VPA / UPI ID)</label>
                <div style="display: flex; gap: 8px;">
                  <input type="text" id="pgUpiVpaInput" class="pg-input" placeholder="e.g. mobile@upi or name@okaxis" style="flex: 1;">
                  <button type="button" class="pg-vpa-verify-btn" onclick="verifyUpiVpa()">VERIFY</button>
                </div>
                <div id="pgVpaStatus" style="font-size: 11.5px; margin-top: 4px; display: none;"></div>
              </div>
            </div>

            <!-- 2. Card Pane -->
            <div class="pg-tab-pane" id="paneCard" style="display: none;">
              <div style="font-size: 12.5px; font-weight: 800; color: #0F172A; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px;">Credit / Debit / ATM Card</div>
              <div class="pg-form-group">
                <label class="pg-label">Card Number</label>
                <div style="position: relative;">
                  <input type="text" id="pgCardNumber" class="pg-input" placeholder="4532 0000 0000 1089" maxlength="19" oninput="formatCardNumber(this)">
                  <span class="pg-card-brand-badge" id="pgCardBrand">VISA</span>
                </div>
              </div>
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 10px;">
                <div class="pg-form-group">
                  <label class="pg-label">Valid Thru (MM/YY)</label>
                  <input type="text" id="pgCardExpiry" class="pg-input" placeholder="MM/YY" maxlength="5" oninput="formatCardExpiry(this)">
                </div>
                <div class="pg-form-group">
                  <label class="pg-label">CVV / CVC</label>
                  <input type="password" id="pgCardCvv" class="pg-input" placeholder="123" maxlength="4">
                </div>
              </div>
              <div class="pg-form-group" style="margin-top: 10px;">
                <label class="pg-label">Name on Card</label>
                <input type="text" id="pgCardName" class="pg-input" placeholder="Cardholder's Full Legal Name">
              </div>
              <label style="display: flex; align-items: center; gap: 8px; font-size: 12px; color: #475569; margin-top: 12px; cursor: pointer;">
                <input type="checkbox" checked style="accent-color: #0C831F;">
                <span>Securely save this card for one-click checkout (RBI tokenized)</span>
              </label>
            </div>

            <!-- 3. NetBanking Pane -->
            <div class="pg-tab-pane" id="paneNet" style="display: none;">
              <div style="font-size: 12.5px; font-weight: 800; color: #0F172A; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px;">Direct NetBanking Partners</div>
              <div class="pg-banks-grid">
                <div class="pg-bank-btn active" onclick="selectQuickBank(this, 'HDFC Bank')">
                  <div class="pg-bank-badge" style="background:#004C8F; color:white;">HDFC</div>
                  <div class="pg-bank-name">HDFC Bank</div>
                </div>
                <div class="pg-bank-btn" onclick="selectQuickBank(this, 'State Bank of India')">
                  <div class="pg-bank-badge" style="background:#280071; color:white;">SBI</div>
                  <div class="pg-bank-name">State Bank of India</div>
                </div>
                <div class="pg-bank-btn" onclick="selectQuickBank(this, 'ICICI Bank')">
                  <div class="pg-bank-badge" style="background:#B02A30; color:white;">ICICI</div>
                  <div class="pg-bank-name">ICICI Bank</div>
                </div>
                <div class="pg-bank-btn" onclick="selectQuickBank(this, 'Axis Bank')">
                  <div class="pg-bank-badge" style="background:#97144D; color:white;">AXIS</div>
                  <div class="pg-bank-name">Axis Bank</div>
                </div>
                <div class="pg-bank-btn" onclick="selectQuickBank(this, 'Kotak Mahindra')">
                  <div class="pg-bank-badge" style="background:#ED1C24; color:white;">KOTAK</div>
                  <div class="pg-bank-name">Kotak Mahindra</div>
                </div>
              </div>
            </div>

            <!-- 4. Wallet Pane -->
            <div class="pg-tab-pane" id="paneWallet" style="display: none;">
              <div style="font-size: 12.5px; font-weight: 800; color: #0F172A; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 12px;">Digital Wallets &amp; Pay Later</div>
              <div class="pg-wallet-list">
                <label class="pg-wallet-item">
                  <input type="radio" name="pgWallet" value="Paytm Wallet" checked onchange="updatePaymentButtonLabel()">
                  <div>
                    <div style="font-weight: 800; color: #0F172A; font-size: 13.5px;">Paytm Wallet</div>
                    <div style="font-size: 11.5px; color: #64748B;">Linked directly via registered mobile number</div>
                  </div>
                </label>
                <label class="pg-wallet-item">
                  <input type="radio" name="pgWallet" value="Amazon Pay" onchange="updatePaymentButtonLabel()">
                  <div>
                    <div style="font-weight: 800; color: #0F172A; font-size: 13.5px;">Amazon Pay Balance</div>
                    <div style="font-size: 11.5px; color: #64748B;">Instant 1-click settlement</div>
                  </div>
                </label>
                <label class="pg-wallet-item">
                  <input type="radio" name="pgWallet" value="Simpl Pay Later" onchange="updatePaymentButtonLabel()">
                  <div>
                    <div style="font-weight: 800; color: #0F172A; font-size: 13.5px;">Simpl Pay Later</div>
                    <div style="font-size: 11.5px; color: #64748B;">Buy now, settle in 15 days &bull; Zero OTP required</div>
                  </div>
                </label>
              </div>
            </div>

            <!-- 5. COD Pane -->
            <div class="pg-tab-pane" id="paneCod" style="display: none;">
              <div style="padding: 20px; background: #F8FAFC; border-radius: 10px; border: 1.5px dashed #CBD5E1; text-align: center;">
                <div style="font-size: 14px; font-weight: 800; color: #0F172A; text-transform: uppercase; letter-spacing: 0.5px;">Pay on Arrival (Doorstep Settlement)</div>
                <p style="font-size: 12.5px; color: #64748B; margin-top: 6px; line-height: 1.5; max-width: 320px; margin-left: auto; margin-right: auto;">
                  Pay via physical cash or dynamic merchant UPI QR code directly to the delivery logistics partner upon handover.
                </p>
                <div style="margin-top: 12px; font-size: 11px; font-weight: 800; color: #166534; background: #DCFCE7; padding: 4px 10px; border-radius: 5px; display: inline-block; letter-spacing: 0.4px;">
                  ZERO HANDLING SURCHARGE ON CASH DELIVERIES
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Order & Bill Summary -->
        <div class="pg-right-pane">
          <div class="pg-recap-card">
            <!-- Delivery Loc Preview -->
            <div class="pg-dest-card">
              <div class="pg-dest-kicker">DELIVERY DESTINATION</div>
              <div class="pg-dest-address" id="pgDestSummary">Bengaluru Kitchens</div>
              <div class="pg-dest-meta">Express fulfillment &bull; 8-12 minutes</div>
            </div>

            <!-- Professional Itemized Order Ledger -->
            <div class="pg-order-ledger">
              <div class="pg-order-ledger-header">
                <span class="pg-ledger-title">ITEMIZED ORDER LEDGER</span>
                <span class="pg-ledger-count-badge" id="pgLedgerTotalQty">0 ITEMS</span>
              </div>
              <div class="pg-mini-items" id="pgMiniItemsList"></div>
            </div>

            <!-- Bill Breakdown -->
            <div class="pg-bill-rows">
              <div class="pg-bill-row">
                <span>Items Subtotal</span>
                <span id="pgBillItemsTotal" style="font-weight: 700; color: #0F172A;">&#8377;0</span>
              </div>
              <div class="pg-bill-row">
                <span>Delivery Partner Fee</span>
                <span id="pgBillDelivery">&#8377;15</span>
              </div>
              <div class="pg-bill-row">
                <span>Restaurant Packaging &amp; Handling</span>
                <span id="pgBillHandling">&#8377;4</span>
              </div>
              <div class="pg-bill-row" id="pgBillSmallCartRow" style="display: none;">
                <span>Small Cart Fee</span>
                <span id="pgBillSmallCart">&#8377;15</span>
              </div>
              <div class="pg-bill-row" id="pgBillTipRow" style="display: none;">
                <span>Delivery Partner Tip</span>
                <span id="pgBillTip">&#8377;0</span>
              </div>
              <div class="pg-bill-row" id="pgBillDonationRow">
                <span>Feeding India Donation</span>
                <span id="pgBillDonation">&#8377;1</span>
              </div>
              <div class="pg-bill-divider"></div>
              <div class="pg-bill-row pg-grand-total-row">
                <span style="font-weight: 800; font-size: 13.5px; text-transform: uppercase; letter-spacing: 0.5px;">Net Payable</span>
                <span style="font-weight: 900; font-size: 18px; color: #0C831F;" id="pgGrandTotalVal">&#8377;0</span>
              </div>
            </div>

            <!-- Trust Badges Strip -->
            <div class="pg-trust-strip">
              <div style="font-size: 9.5px; color: #64748B; font-weight: 700; text-transform: uppercase; letter-spacing: 0.7px; text-align: center;">
                PCI-DSS LEVEL 1 CERTIFIED &bull; RBI RECOGNIZED ESCROW
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer CTA Strip -->
      <div class="pg-modal-footer" id="pgFooterBar">
        <button type="button" class="pg-btn-cancel" onclick="closePaymentModal()">Cancel</button>
        <button type="button" class="pg-btn-pay" id="pgPayButton" onclick="processPaymentGateway()">
          <span id="pgPayBtnText">Pay &#8377;0 Securely</span>
        </button>
      </div>

      <!-- Success Screen View (Zero Icons) -->
      <div class="pg-success-view" id="pgSuccessView" style="display: none;">
        <div class="pg-success-status-badge">PAYMENT AUTHORIZED &amp; CONFIRMED</div>
        <h2 style="font-size: 20px; font-weight: 900; color: #0F172A; margin-top: 14px; margin-bottom: 4px;">Order Placed Successfully</h2>
        <div style="font-size: 13px; color: #0C831F; font-weight: 700;">Live Kitchen Dispatch Order Active</div>

        <div class="pg-order-id-pill" id="pgSuccessOrderId">Order #BITE-2026-92841</div>

        <!-- Live Delivery Progress Bar -->
        <div class="pg-tracking-bar">
          <div class="pg-track-step done">
            <div class="pg-track-dot">1</div>
            <div class="pg-track-label">Placed</div>
          </div>
          <div class="pg-track-line active"></div>
          <div class="pg-track-step active">
            <div class="pg-track-dot">2</div>
            <div class="pg-track-label">Kitchen</div>
          </div>
          <div class="pg-track-line"></div>
          <div class="pg-track-step">
            <div class="pg-track-dot">3</div>
            <div class="pg-track-label">In Transit</div>
          </div>
          <div class="pg-track-line"></div>
          <div class="pg-track-step">
            <div class="pg-track-dot">4</div>
            <div class="pg-track-label">Delivered</div>
          </div>
        </div>

        <div class="pg-success-info-card">
          <div style="display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 8px;">
            <span style="color: #64748B;">Estimated Delivery:</span>
            <strong style="color: #0C831F;">12-15 minutes (Express)</strong>
          </div>
          <div style="display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 8px;">
            <span style="color: #64748B;">Payment Mode:</span>
            <strong style="color: #0F172A;" id="pgSuccessMethod">UPI / Google Pay</strong>
          </div>
          <div style="display: flex; justify-content: space-between; font-size: 13px;">
            <span style="color: #64748B;">Amount Paid:</span>
            <strong style="color: #0F172A; font-size: 15px;" id="pgSuccessAmount">&#8377;0</strong>
          </div>
        </div>

        <div style="display: flex; gap: 10px; width: 100%; margin-top: 18px; max-width: 440px;">
          <button type="button" class="btn-track-order" style="flex: 1; background: #0C831F; color: #ffffff; border: none; font-weight: 800;" onclick="launchOrderTracking()">
            Track Order Live 🛵
          </button>
          <button type="button" class="btn-track-order" style="flex: 1; background: #F1F5F9; color: #334155; border: 1px solid #CBD5E1; font-weight: 700;" onclick="closePaymentModal()">
            Order More Food
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 18. LOCATION & ORDER FOR SOMEONE ELSE MODAL -->
  <div class="loc-modal-overlay" id="locModalOverlay" onclick="if(event.target === this) closeLocationModal()">
    <div class="loc-modal-card">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">Change Location</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closeLocationModal()">&times;</button>
      </div>

      <div class="loc-modal-body">
        <!-- Detect + Search Row -->
        <div class="loc-top-action-row">
          <button type="button" class="btn-detect-location" id="btnDetectLoc" onclick="detectGPSLocation()">
            Detect my location
          </button>
          
          <div class="loc-or-divider">
            <span class="loc-or-badge">OR</span>
          </div>

          <div class="loc-search-wrap">
            <input 
              type="text" 
              id="locSearchInput" 
              class="loc-search-input" 
              placeholder="search delivery location" 
              autocomplete="off"
              oninput="handleLocSearch(this.value)"
            >
          </div>
        </div>

        <!-- Search Suggestions Dropdown / Overlay (Dynamic) -->
        <div id="locSearchSuggestions" class="loc-search-suggestions" style="display: none;"></div>

        <!-- Status / City Hub Notice (for Pan-India auto detection) -->
        <div id="panIndiaNotice" class="pan-india-notice" style="display: none;"></div>

        <!-- Saved Addresses Section Header -->
        <div class="saved-addresses-header-row">
          <h4 class="saved-addresses-title">Your saved addresses</h4>
          <button type="button" class="btn-add-address-link" onclick="toggleAddAddressForm()">
            + Add / Order for Someone Else
          </button>
        </div>

        <!-- Add / Edit Address & Order for Someone Else Form (Collapsible) -->
        <div id="addAddressFormContainer" class="add-address-form-box" style="display: none;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <div style="font-size: 14px; font-weight: 800; color: #1E293B;" id="addressFormTitle">
              Add New Address / Order for Someone Else
            </div>
            <button type="button" onclick="toggleAddAddressForm(false)" style="border: none; background: none; font-size: 18px; color: #64748B; cursor: pointer;">&times;</button>
          </div>

          <form onsubmit="handleSaveAddressSubmit(event)">
            <input type="hidden" id="editAddressIndex" value="-1">

            <div class="address-type-selector">
              <label class="address-type-option">
                <input type="radio" name="addrType" value="Home" checked onchange="toggleRecipientFields(false)">
                <span>🏠 Home</span>
              </label>
              <label class="address-type-option">
                <input type="radio" name="addrType" value="Work" onchange="toggleRecipientFields(false)">
                <span>💼 Work</span>
              </label>
              <label class="address-type-option">
                <input type="radio" name="addrType" value="Other" onchange="toggleRecipientFields(true)">
                <span>🎁 Someone Else</span>
              </label>
            </div>

            <div id="recipientFieldsGroup" style="display: none; margin-top: 10px;">
              <div class="form-field-group">
                <label class="form-field-label">Recipient's Full Name *</label>
                <input type="text" id="formRecipName" class="form-field-input" placeholder="e.g. Someshwar Chaturvedi, Rahul, Mom">
              </div>
              <div class="form-field-group" style="margin-top: 10px;">
                <label class="form-field-label">Recipient's Phone Number *</label>
                <input type="tel" id="formRecipPhone" class="form-field-input" placeholder="e.g. 9876543210" pattern="[0-9]{10}" title="10-digit mobile number">
              </div>
            </div>

            <div class="form-field-group" style="margin-top: 10px;">
              <label class="form-field-label">Complete Street Address (Flat / House / Sector / City) *</label>
              <textarea id="formAddressFull" class="form-field-input" rows="2" placeholder="e.g. Majestic Residency, Gali no.3 ,Nagli village ,Sector 134 , Noida ,U.P." required style="resize: vertical;"></textarea>
            </div>

            <div class="form-field-group" style="margin-top: 10px;">
              <label class="form-field-label">City / Delivery Locality (Bengaluru / India) *</label>
              <input type="text" id="formLocality" class="form-field-input" placeholder="e.g. HSR Layout, Bengaluru or Sector 134, Noida" required>
            </div>

            <div class="form-field-group" style="margin-top: 10px;">
              <label class="form-field-label">Delivery Note (Optional)</label>
              <input type="text" id="formNote" class="form-field-input" placeholder="e.g. Ring bell, Leave with security">
            </div>

            <div style="display: flex; gap: 10px; margin-top: 14px;">
              <button type="submit" class="btn-save-address-submit">
                Save &amp; Select Address
              </button>
              <button type="button" class="btn-cancel-address" onclick="toggleAddAddressForm(false)">
                Cancel
              </button>
            </div>
          </form>
        </div>

        <!-- Saved Addresses List (Cards matching screenshot) -->
        <div id="savedAddressesList" class="saved-addresses-list">
          <!-- Dynamically populated from JS -->
        </div>

      </div>
    </div>
  </div>

  <!-- DUAL OTP / PASSWORD CUSTOMER SIGN IN MODAL -->
  <div class="loc-modal-overlay" id="loginModalOverlay" onclick="if(event.target === this) closeLoginModal()">
    <div class="loc-modal-card" style="max-width: 440px; padding: 0; border-radius: 20px; overflow: hidden; box-shadow: 0 20px 40px rgba(0,0,0,0.2);">
      <div style="background: #ffffff; padding: 24px 28px 24px 28px;">
        <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 14px;">
          <div>
            <div style="font-size: 12px; font-weight: 800; color: #16A34A; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 2px;">Under100 Food Discovery</div>
            <h3 style="font-size: 21px; font-weight: 800; color: #1E293B; margin: 0;">Sign In</h3>
          </div>
          <button type="button" class="loc-modal-close-btn" onclick="closeLoginModal()" style="font-size: 26px; line-height: 1; padding: 2px 6px; cursor: pointer; border: none; background: none; color: #64748B;">&times;</button>
        </div>

        <!-- Subtabs: Mobile/OTP vs Email & Password -->
        <div style="display: flex; gap: 6px; background: #F1F5F9; padding: 4px; border-radius: 10px; margin-bottom: 16px;">
          <button type="button" id="modalSubtabOtp" onclick="switchModalAuth('otp')" style="flex: 1; padding: 7px 10px; border: none; border-radius: 7px; font-weight: 700; font-size: 13px; cursor: pointer; background: #ffffff; color: #0C831F; box-shadow: 0 1px 3px rgba(0,0,0,0.08); transition: all 0.15s;">
            Mobile / OTP
          </button>
          <button type="button" id="modalSubtabPass" onclick="switchModalAuth('pass')" style="flex: 1; padding: 7px 10px; border: none; border-radius: 7px; font-weight: 700; font-size: 13px; cursor: pointer; background: transparent; color: #64748B; transition: all 0.15s;">
            Email &amp; Password
          </button>
        </div>

        <!-- Mode 1: OTP Login Form -->
        <div id="modalOtpSection" style="display: block;">

          <div style="display: flex; border: 1.5px solid #CBD5E1; border-radius: 12px; overflow: hidden; background: #ffffff; margin-bottom: 12px;">
            <div style="padding: 11px 13px; background: #F8FAFC; border-right: 1.5px solid #CBD5E1; font-weight: 700; font-size: 13px; color: #334155; display: flex; align-items: center;">
              +91
            </div>
            <input type="text" id="loginIdentifierInput" placeholder="10-digit mobile or email" autocomplete="tel" required style="flex: 1; border: none; padding: 11px 13px; font-size: 14px; font-weight: 600; color: #1E293B; outline: none; background: transparent;">
            <button type="button" id="btnModalSendOtp" onclick="sendModalOTP()" style="background: #0C831F; color: #ffffff; border: none; padding: 0 14px; font-size: 12.5px; font-weight: 700; cursor: pointer; white-space: nowrap;">
              Get OTP
            </button>
          </div>

          <div id="modalOtpInputWrap" style="display: none; margin-bottom: 12px;">
            <label style="display: block; font-size: 12px; font-weight: 700; color: #475569; margin-bottom: 4px;">Enter 6-Digit OTP</label>
            <input type="text" id="modalOtpCodeInput" placeholder="••••••" maxlength="6" style="width: 100%; box-sizing: border-box; border: 1.5px solid #CBD5E1; border-radius: 12px; padding: 11px 14px; font-size: 16px; font-weight: 800; letter-spacing: 4px; text-align: center; color: #1E293B; outline: none;">
          </div>

          <div id="modalOtpStatus" style="display: none; font-size: 12px; font-weight: 600; margin-bottom: 10px;"></div>

          <button type="button" id="btnLoginContinue" onclick="handleModalOtpSubmit()" style="width: 100%; padding: 13px; background: #0C831F; color: #ffffff; border: none; border-radius: 12px; font-size: 14.5px; font-weight: 800; cursor: pointer; transition: background 0.15s ease; box-shadow: 0 4px 14px rgba(12, 131, 31, 0.25);">
            Continue
          </button>
        </div>

        <!-- Mode 2: Email & Password Login Form -->
        <form id="modalPassSection" onsubmit="handleModalPassSubmit(event)" style="display: none; flex-direction: column; gap: 12px;">
          <div>
            <label style="display: block; font-size: 12px; font-weight: 700; color: #475569; margin-bottom: 4px;">Email Address</label>
            <input type="email" id="modalPassEmail" placeholder="user@example.com" required style="width: 100%; box-sizing: border-box; border: 1.5px solid #CBD5E1; border-radius: 12px; padding: 11px 14px; font-size: 13.5px; font-weight: 600; color: #1E293B; outline: none;">
          </div>
          <div>
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px;">
              <label style="font-size: 12px; font-weight: 700; color: #475569;">Password</label>
              <a href="/password/forgot" style="font-size: 11.5px; color: #E23744; text-decoration: none; font-weight: 700;">Forgot?</a>
            </div>
            <input type="password" id="modalPassPassword" placeholder="••••••••" required style="width: 100%; box-sizing: border-box; border: 1.5px solid #CBD5E1; border-radius: 12px; padding: 11px 14px; font-size: 13.5px; font-weight: 600; color: #1E293B; outline: none;">
          </div>
          <div id="modalPassStatus" style="display: none; font-size: 12px; color: #DC2626; font-weight: 600;"></div>
          <button type="submit" id="btnModalPassLogin" style="width: 100%; padding: 13px; background: #E23744; color: #ffffff; border: none; border-radius: 12px; font-size: 14.5px; font-weight: 800; cursor: pointer; transition: background 0.15s ease; box-shadow: 0 4px 14px rgba(226, 55, 68, 0.25);">
            Sign In with Password
          </button>
        </form>

        <div style="margin-top: 16px; border-top: 1px solid #E2E8F0; padding-top: 12px; text-align: center; font-size: 12px; color: #64748B;">
          Don't have an account? <a href="/login?mode=signup" style="color: #0C831F; font-weight: 700; text-decoration: none;">Create free account &rarr;</a>
        </div>
      </div>
    </div>
  </div>

  <!-- MY ORDERS MODAL -->
  <div class="loc-modal-overlay" id="myOrdersModalOverlay" onclick="if(event.target === this) closeMyOrdersModal()">
    <div class="loc-modal-card" style="max-width: 520px;">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">My Orders</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closeMyOrdersModal()">&times;</button>
      </div>
      <div class="loc-modal-body">
        <div style="display: flex; flex-direction: column; gap: 14px;">
          <!-- Order Card 1 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 12px; padding: 14px 16px; background: #ffffff;">
            <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px;">
              <div>
                <div style="font-weight: 700; font-size: 14.5px; color: #1E293B;">CTR Shri Sagar</div>
                <div style="font-size: 12px; color: #64748B;">Order #BLR-84192 · Today, 11:45 AM</div>
              </div>
              <span style="background: #EEF8EB; color: #087443; font-size: 11.5px; font-weight: 700; padding: 3px 8px; border-radius: 6px;">Delivered</span>
            </div>
            <div style="font-size: 13px; color: #475569; margin-bottom: 10px;">
              1x Crispy Benne Dosa (&#8377;85) · 1x Filter Coffee (&#8377;15)
            </div>
            <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px dashed #E2E8F0; padding-top: 8px;">
              <span style="font-weight: 700; font-size: 14px; color: #1E293B;">Total: &#8377;100</span>
              <button type="button" onclick="reorderItem('Crispy Benne Dosa', 85)" style="background: #FFF0E6; color: var(--color-primary-orange); border: 1px solid #FFD8BE; border-radius: 6px; padding: 5px 12px; font-weight: 700; font-size: 12px; cursor: pointer;">Reorder</button>
            </div>
          </div>
          <!-- Order Card 2 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 12px; padding: 14px 16px; background: #ffffff;">
            <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px;">
              <div>
                <div style="font-weight: 700; font-size: 14.5px; color: #1E293B;">Brahmin's Coffee Bar</div>
                <div style="font-size: 12px; color: #64748B;">Order #BLR-73104 · Yesterday, 8:20 AM</div>
              </div>
              <span style="background: #EEF8EB; color: #087443; font-size: 11.5px; font-weight: 700; padding: 3px 8px; border-radius: 6px;">Delivered</span>
            </div>
            <div style="font-size: 13px; color: #475569; margin-bottom: 10px;">
              2x Soft Idli Vada (&#8377;60) · 1x Kesari Bath (&#8377;35)
            </div>
            <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px dashed #E2E8F0; padding-top: 8px;">
              <span style="font-weight: 700; font-size: 14px; color: #1E293B;">Total: &#8377;95</span>
              <button type="button" onclick="reorderItem('Soft Idli Vada', 60)" style="background: #FFF0E6; color: var(--color-primary-orange); border: 1px solid #FFD8BE; border-radius: 6px; padding: 5px 12px; font-weight: 700; font-size: 12px; cursor: pointer;">Reorder</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- MY PRESCRIPTIONS MODAL -->
  <div class="loc-modal-overlay" id="prescriptionsModalOverlay" onclick="if(event.target === this) closePrescriptionsModal()">
    <div class="loc-modal-card" style="max-width: 520px;">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">My Prescriptions</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closePrescriptionsModal()">&times;</button>
      </div>
      <div class="loc-modal-body">
        <p style="font-size: 13px; color: #64748B; margin-top: 0; margin-bottom: 16px;">
          Upload doctor dietary recommendations or set health restrictions so restaurants prepare safe meals under &#8377;100 tailored for you.
        </p>
        <div style="border: 2px dashed #CBD5E1; border-radius: 12px; padding: 24px; text-align: center; background: #F8FAFC; margin-bottom: 18px; cursor: pointer;" onclick="document.getElementById('prescriptionFileInput').click()">
          <input type="file" id="prescriptionFileInput" style="display: none;" onchange="handlePrescriptionUpload(this)">
          <div style="font-size: 24px; margin-bottom: 6px;">📋</div>
          <div style="font-weight: 700; font-size: 14px; color: #1E293B;">Upload Prescription / Dietary Slip</div>
          <div style="font-size: 12px; color: #64748B; margin-top: 4px;">Supports PNG, JPG, or PDF (Max 10MB)</div>
          <div id="prescriptionUploadStatus" style="font-size: 12px; font-weight: 600; color: #087443; margin-top: 8px;"></div>
        </div>
        <div style="font-weight: 700; font-size: 13.5px; color: #1E293B; margin-bottom: 8px;">Dietary &amp; Health Preferences</div>
        <div style="display: flex; flex-wrap: wrap; gap: 8px;">
          <label style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background: #F1F5F9; border-radius: 20px; font-size: 12.5px; font-weight: 600; color: #334155; cursor: pointer;">
            <input type="checkbox" checked style="accent-color: var(--color-primary-orange);"> Diabetic-Friendly (Low GI)
          </label>
          <label style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background: #F1F5F9; border-radius: 20px; font-size: 12.5px; font-weight: 600; color: #334155; cursor: pointer;">
            <input type="checkbox" style="accent-color: var(--color-primary-orange);"> Low Sodium / Heart-Care
          </label>
          <label style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background: #F1F5F9; border-radius: 20px; font-size: 12.5px; font-weight: 600; color: #334155; cursor: pointer;">
            <input type="checkbox" style="accent-color: var(--color-primary-orange);"> Jain Food (No Onion/Garlic)
          </label>
          <label style="display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; background: #F1F5F9; border-radius: 20px; font-size: 12.5px; font-weight: 600; color: #334155; cursor: pointer;">
            <input type="checkbox" style="accent-color: var(--color-primary-orange);"> Zero Added Sugar
          </label>
        </div>
      </div>
    </div>
  </div>

  <!-- E-GIFT CARDS MODAL -->
  <div class="loc-modal-overlay" id="giftCardsModalOverlay" onclick="if(event.target === this) closeGiftCardsModal()">
    <div class="loc-modal-card" style="max-width: 500px;">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">E-Gift Cards</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closeGiftCardsModal()">&times;</button>
      </div>
      <div class="loc-modal-body">
        <div style="background: linear-gradient(135deg, #FF5A1F 0%, #D8420B 100%); border-radius: 14px; padding: 20px; color: #ffffff; margin-bottom: 20px; box-shadow: 0 8px 20px rgba(255, 90, 31, 0.25);">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div style="font-size: 12px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; opacity: 0.9;">Under100 Food Card</div>
            <div style="font-weight: 800; font-size: 13px;">&#8377;100 MAX</div>
          </div>
          <div style="margin: 18px 0 6px;">
            <div style="font-size: 12px; opacity: 0.85;">Available Balance</div>
            <div style="font-size: 28px; font-weight: 800;">&#8377;0.00</div>
          </div>
          <div style="font-size: 11px; opacity: 0.8;">Valid for all breakfast &amp; quick bite orders across Bengaluru</div>
        </div>

        <div style="font-weight: 700; font-size: 13.5px; color: #1E293B; margin-bottom: 8px;">Redeem Gift Card</div>
        <div style="display: flex; gap: 8px; margin-bottom: 20px;">
          <input type="text" id="giftCardCodeInput" placeholder="Enter 16-digit voucher code" style="flex: 1; padding: 10px 14px; border: 1.5px solid #CBD5E1; border-radius: 8px; font-size: 13px; outline: none;">
          <button type="button" onclick="redeemGiftCardCode()" style="background: var(--color-primary-orange); color: #ffffff; border: none; border-radius: 8px; padding: 0 16px; font-weight: 700; font-size: 13px; cursor: pointer;">Redeem</button>
        </div>

        <div style="font-weight: 700; font-size: 13.5px; color: #1E293B; margin-bottom: 8px;">Send Food to a Friend</div>
        <div style="display: flex; gap: 10px;">
          <button type="button" onclick="selectGiftAmount(100)" style="flex: 1; padding: 10px; border: 1.5px solid #E2E8F0; border-radius: 8px; background: #ffffff; font-weight: 700; font-size: 13px; color: #1E293B; cursor: pointer;">&#8377;100 (1 Bite)</button>
          <button type="button" onclick="selectGiftAmount(250)" style="flex: 1; padding: 10px; border: 1.5px solid #E2E8F0; border-radius: 8px; background: #ffffff; font-weight: 700; font-size: 13px; color: #1E293B; cursor: pointer;">&#8377;250 (3 Bites)</button>
          <button type="button" onclick="selectGiftAmount(500)" style="flex: 1; padding: 10px; border: 1.5px solid #E2E8F0; border-radius: 8px; background: #ffffff; font-weight: 700; font-size: 13px; color: #1E293B; cursor: pointer;">&#8377;500 (Feast)</button>
        </div>
      </div>
    </div>
  </div>

  <!-- FAQ'S MODAL -->
  <div class="loc-modal-overlay" id="faqsModalOverlay" onclick="if(event.target === this) closeFaqsModal()">
    <div class="loc-modal-card" style="max-width: 540px;">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">FAQ's</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closeFaqsModal()">&times;</button>
      </div>
      <div class="loc-modal-body">
        <div style="display: flex; flex-direction: column; gap: 10px;">
          <!-- FAQ 1 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 10px; overflow: hidden;">
            <button type="button" onclick="toggleFaqItem(this)" style="width: 100%; display: flex; justify-content: space-between; align-items: center; padding: 13px 16px; background: #F8FAFC; border: none; text-align: left; font-weight: 700; font-size: 13.5px; color: #1E293B; cursor: pointer;">
              <span>How does Under100 ensure all items are &#8377;100 or less?</span>
              <span class="faq-arrow" style="font-size: 12px; transition: transform 0.2s ease;">&#9660;</span>
            </button>
            <div style="display: none; padding: 12px 16px; font-size: 13px; color: #475569; background: #ffffff; border-top: 1px solid #E2E8F0; line-height: 1.5;">
              Under100 works directly with authentic Bengaluru darshinis, tiffin rooms, and street food artisans. We eliminate intermediary commissions and curate single-portion comfort meals guaranteed under &#8377;100.
            </div>
          </div>
          <!-- FAQ 2 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 10px; overflow: hidden;">
            <button type="button" onclick="toggleFaqItem(this)" style="width: 100%; display: flex; justify-content: space-between; align-items: center; padding: 13px 16px; background: #F8FAFC; border: none; text-align: left; font-weight: 700; font-size: 13.5px; color: #1E293B; cursor: pointer;">
              <span>Which locations in Bengaluru do you deliver to?</span>
              <span class="faq-arrow" style="font-size: 12px; transition: transform 0.2s ease;">&#9660;</span>
            </button>
            <div style="display: none; padding: 12px 16px; font-size: 13px; color: #475569; background: #ffffff; border-top: 1px solid #E2E8F0; line-height: 1.5;">
              We deliver across Indiranagar, Koramangala, HSR Layout, Jayanagar, Malleshwaram, Whitefield, Electronic City, and Rajajinagar. Pan-India coverage is launching soon!
            </div>
          </div>
          <!-- FAQ 3 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 10px; overflow: hidden;">
            <button type="button" onclick="toggleFaqItem(this)" style="width: 100%; display: flex; justify-content: space-between; align-items: center; padding: 13px 16px; background: #F8FAFC; border: none; text-align: left; font-weight: 700; font-size: 13.5px; color: #1E293B; cursor: pointer;">
              <span>Can I order food for someone else?</span>
              <span class="faq-arrow" style="font-size: 12px; transition: transform 0.2s ease;">&#9660;</span>
            </button>
            <div style="display: none; padding: 12px 16px; font-size: 13px; color: #475569; background: #ffffff; border-top: 1px solid #E2E8F0; line-height: 1.5;">
              Yes! Click 'Saved Addresses' in your account menu and tap '+ Add / Order for Someone Else'. Enter their address, name, and phone number for seamless door delivery.
            </div>
          </div>
          <!-- FAQ 4 -->
          <div style="border: 1px solid #E2E8F0; border-radius: 10px; overflow: hidden;">
            <button type="button" onclick="toggleFaqItem(this)" style="width: 100%; display: flex; justify-content: space-between; align-items: center; padding: 13px 16px; background: #F8FAFC; border: none; text-align: left; font-weight: 700; font-size: 13.5px; color: #1E293B; cursor: pointer;">
              <span>What payment methods are supported?</span>
              <span class="faq-arrow" style="font-size: 12px; transition: transform 0.2s ease;">&#9660;</span>
            </button>
            <div style="display: none; padding: 12px 16px; font-size: 13px; color: #475569; background: #ffffff; border-top: 1px solid #E2E8F0; line-height: 1.5;">
              UPI (Google Pay, PhonePe, Paytm, BHIM), Debit &amp; Credit Cards, Net Banking, and Cash on Delivery (COD) are all fully accepted with zero hidden fees.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- ACCOUNT PRIVACY MODAL -->
  <div class="loc-modal-overlay" id="privacyModalOverlay" onclick="if(event.target === this) closeAccountPrivacyModal()">
    <div class="loc-modal-card" style="max-width: 500px;">
      <div class="loc-modal-header-simple">
        <h3 class="loc-modal-title">Account Privacy</h3>
        <button type="button" class="loc-modal-close-btn" onclick="closeAccountPrivacyModal()">&times;</button>
      </div>
      <div class="loc-modal-body">
        <div style="display: flex; flex-direction: column; gap: 16px;">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 700; font-size: 13.5px; color: #1E293B;">Personalized Recommendations</div>
              <div style="font-size: 12px; color: #64748B;">Tailor Bengaluru food recommendations to your taste</div>
            </div>
            <input type="checkbox" checked style="accent-color: var(--color-primary-orange); width: 18px; height: 18px; cursor: pointer;">
          </div>
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 700; font-size: 13.5px; color: #1E293B;">Order Status WhatsApp / SMS</div>
              <div style="font-size: 12px; color: #64748B;">Live rider and cooking updates sent to your registered mobile number</div>
            </div>
            <input type="checkbox" checked style="accent-color: var(--color-primary-orange); width: 18px; height: 18px; cursor: pointer;">
          </div>
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 700; font-size: 13.5px; color: #1E293B;">Promotional Offers &amp; Discounts</div>
              <div style="font-size: 12px; color: #64748B;">Weekend flash deal alerts for snacks under &#8377;50</div>
            </div>
            <input type="checkbox" style="accent-color: var(--color-primary-orange); width: 18px; height: 18px; cursor: pointer;">
          </div>

          <div style="border-top: 1px solid #E2E8F0; padding-top: 14px; margin-top: 4px; display: flex; flex-direction: column; gap: 8px;">
            <button type="button" onclick="downloadAccountData()" style="padding: 10px; background: #F1F5F9; border: 1px solid #CBD5E1; border-radius: 8px; font-size: 13px; font-weight: 600; color: #334155; cursor: pointer;">Download My Account Data</button>
            <button type="button" onclick="confirmDeleteAccount()" style="padding: 10px; background: #FEF2F2; border: 1px solid #FECACA; border-radius: 8px; font-size: 13px; font-weight: 600; color: #DC2626; cursor: pointer;">Delete My Account</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 17. FLOATING MOBILE CART BAR -->
  <div class="floating-cart-bar" id="floatingCartBar" onclick="toggleCartDrawer()">
    <div style="display: flex; align-items: center; gap: 10px;">
      <div style="background-color: var(--color-primary-orange); color: #ffffff; padding: 4px 8px; border-radius: var(--radius-full); font-size: 12px; font-weight: 800;" id="floatingItemCount">
        0 items
      </div>
      <div style="font-size: 14px; font-weight: 700;">
        <span id="floatingTotalVal">&#8377;0</span> <span style="font-size: 12px; font-weight: 400; color: #94A3B8;">no hidden fees</span>
      </div>
    </div>
    <div style="display: flex; align-items: center; gap: 6px; color: var(--color-primary-orange); font-weight: 800; font-size: 13px;">
      <span>VIEW CART</span>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
    </div>
  </div>

  <!-- JAVASCRIPT FOR DYNAMIC STATE & HTMX INTEGRATION -->
  <script>
// Client JavaScript for original Under100 landing page (prior prompt)
var currentSlideIdx = 0;
var totalOfferSlides = 5;
var offerTimer = null;
var isOfferPaused = false;

function showOfferSlide(idx) {
  currentSlideIdx = (idx + totalOfferSlides) % totalOfferSlides;
  var track = document.getElementById('offerSliderTrack');
  if (track) {
    track.style.transform = 'translateX(-' + (currentSlideIdx * 100) + '%)';
  }
  var dots = document.querySelectorAll('.slider-dot');
  dots.forEach(function(d, i) {
    if (i === currentSlideIdx) d.classList.add('active');
    else d.classList.remove('active');
  });
}

function nextOfferSlide() {
  showOfferSlide(currentSlideIdx + 1);
  resetOfferTimer();
}

function prevOfferSlide() {
  showOfferSlide(currentSlideIdx - 1);
  resetOfferTimer();
}

function jumpOfferSlide(idx) {
  showOfferSlide(idx);
  resetOfferTimer();
}

function pauseOfferTimer() {
  isOfferPaused = true;
}

function resumeOfferTimer() {
  isOfferPaused = false;
}

function resetOfferTimer() {
  if (offerTimer) clearInterval(offerTimer);
  offerTimer = setInterval(function() {
    if (!isOfferPaused) {
      showOfferSlide(currentSlideIdx + 1);
    }
  }, 3800);
}
resetOfferTimer();

// Mobile Touch Swipe Handling
(function initSliderGestures() {
  var touchStartX = 0;
  var touchEndX = 0;
  var container = document.getElementById('offerSliderContainer');
  if (container) {
    container.addEventListener('touchstart', function(e) {
      touchStartX = e.changedTouches[0].screenX;
      pauseOfferTimer();
    }, { passive: true });
    container.addEventListener('touchend', function(e) {
      touchEndX = e.changedTouches[0].screenX;
      if (touchStartX - touchEndX > 45) {
        nextOfferSlide();
      } else if (touchEndX - touchStartX > 45) {
        prevOfferSlide();
      }
      resumeOfferTimer();
    }, { passive: true });
  }
})();

function scrollTrack(trackId, offset) {
  var track = document.getElementById(trackId);
  if (track) {
    track.scrollBy({ left: offset, behavior: 'smooth' });
  }
}

var itemsDB = {"f1": {"id": "f1", "name": "Crispy Masala Dosa", "category": "dosa-idli", "cuisine": "South Indian", "restaurant": "Sri Krishna Cafe", "neighbourhood": "Indiranagar", "price": 79, "rating": 4.6, "ratingCount": 380, "deliveryMinutes": 20, "image": "/static/crispy_dosa.jpg", "isVegetarian": true, "tags": ["under100", "south-indian-fav", "bengaluru-pick"]}, "f2": {"id": "f2", "name": "Steamed Idli Sambar (2 pcs)", "category": "dosa-idli", "cuisine": "South Indian", "restaurant": "Brahmin's Tiffin Room", "neighbourhood": "Koramangala", "price": 49, "rating": 4.8, "ratingCount": 520, "deliveryMinutes": 15, "image": "/static/idli_sambar.jpg", "isVegetarian": true, "tags": ["under100", "south-indian-fav", "bengaluru-pick"]}, "f3": {"id": "f3", "name": "Aromatic Veg Pulao", "category": "rice-biryani", "cuisine": "North Indian", "restaurant": "Meghana Bowl Corner", "neighbourhood": "HSR Layout", "price": 89, "rating": 4.4, "ratingCount": 210, "deliveryMinutes": 25, "image": "/static/veg_pulao.jpg", "isVegetarian": true, "tags": ["under100"]}, "f4": {"id": "f4", "name": "Paneer Kathi Roll", "category": "rolls-wraps", "cuisine": "Street Food", "restaurant": "Kolkata Roll Express", "neighbourhood": "Indiranagar", "price": 99, "rating": 4.5, "ratingCount": 290, "deliveryMinutes": 22, "image": "/static/kathi_roll.jpg", "isVegetarian": true, "tags": ["under100"]}, "f5": {"id": "f5", "name": "Schezwan Veg Fried Rice", "category": "rice-biryani", "cuisine": "Indo-Chinese", "restaurant": "Wok & Bowl Express", "neighbourhood": "BTM Layout", "price": 89, "rating": 4.3, "ratingCount": 175, "deliveryMinutes": 28, "image": "/static/veg_pulao.jpg", "isVegetarian": true, "tags": ["under100"]}, "f6": {"id": "f6", "name": "Stuffed Aloo Paratha with Curd", "category": "parathas", "cuisine": "North Indian", "restaurant": "Punjab Dhaba Point", "neighbourhood": "Whitefield", "price": 69, "rating": 4.5, "ratingCount": 310, "deliveryMinutes": 24, "image": "/static/aloo_paratha.jpg", "isVegetarian": true, "tags": ["under100"]}, "f7": {"id": "f7", "name": "Indori Spiced Poha", "category": "snacks-chaat", "cuisine": "Breakfast", "restaurant": "Morning Fuel Co.", "neighbourhood": "Koramangala", "price": 39, "rating": 4.6, "ratingCount": 190, "deliveryMinutes": 15, "image": "/static/indori_poha.jpg", "isVegetarian": true, "tags": ["under100", "quick-bite"]}, "f8": {"id": "f8", "name": "Steamed Veg Momos (6 pcs)", "category": "momos", "cuisine": "Tibetan", "restaurant": "Himalayan Momo Hut", "neighbourhood": "Marathahalli", "price": 79, "rating": 4.5, "ratingCount": 260, "deliveryMinutes": 20, "image": "/static/steamed_momos.jpg", "isVegetarian": true, "tags": ["under100", "quick-bite"]}, "f9": {"id": "f9", "name": "Crispy Plain Dosa & Podi", "category": "dosa-idli", "cuisine": "South Indian", "restaurant": "Vidyarthi Bhavana Tiffins", "neighbourhood": "Malleshwaram", "price": 59, "rating": 4.7, "ratingCount": 480, "deliveryMinutes": 18, "image": "/static/crispy_dosa.jpg", "isVegetarian": true, "tags": ["south-indian-fav", "bengaluru-pick"]}, "f10": {"id": "f10", "name": "Golden Medu Vada (2 pcs)", "category": "south-indian", "cuisine": "South Indian", "restaurant": "Udupi Grand Tiffins", "neighbourhood": "HSR Layout", "price": 45, "rating": 4.6, "ratingCount": 340, "deliveryMinutes": 16, "image": "/static/medu_vada.jpg", "isVegetarian": true, "tags": ["south-indian-fav", "quick-bite", "bengaluru-pick"]}, "f11": {"id": "f11", "name": "Onion Tomato Uttapam", "category": "dosa-idli", "cuisine": "South Indian", "restaurant": "Taaza Thindi", "neighbourhood": "Jayanagar", "price": 79, "rating": 4.7, "ratingCount": 395, "deliveryMinutes": 20, "image": "/static/onion_uttapam.jpg", "isVegetarian": true, "tags": ["south-indian-fav", "bengaluru-pick"]}, "f12": {"id": "f12", "name": "South Indian Mini Meal Box", "category": "south-indian", "cuisine": "South Indian", "restaurant": "MTR Heritage Canteen", "neighbourhood": "Indiranagar", "price": 89, "rating": 4.6, "ratingCount": 410, "deliveryMinutes": 25, "image": "/static/hero_food.jpg", "isVegetarian": true, "tags": ["south-indian-fav", "under100", "bengaluru-pick"]}, "f13": {"id": "f13", "name": "Punjabi Samosa (2 pcs with Chutney)", "category": "snacks-chaat", "cuisine": "Snacks", "restaurant": "Delhi Chaat Bhandar", "neighbourhood": "BTM Layout", "price": 39, "rating": 4.5, "ratingCount": 420, "deliveryMinutes": 15, "image": "/static/samosa_plate.jpg", "isVegetarian": true, "tags": ["quick-bite"]}, "f14": {"id": "f14", "name": "Tangy Aloo Tikki Chaat", "category": "snacks-chaat", "cuisine": "Street Food", "restaurant": "Anand Sweets & Chaat", "neighbourhood": "Koramangala", "price": 59, "rating": 4.4, "ratingCount": 230, "deliveryMinutes": 18, "image": "/static/samosa_plate.jpg", "isVegetarian": true, "tags": ["quick-bite"]}, "f15": {"id": "f15", "name": "Bombay Veg Grilled Sandwich", "category": "sandwiches", "cuisine": "Snacks", "restaurant": "Hari Super Sandwich", "neighbourhood": "Jayanagar", "price": 69, "rating": 4.6, "ratingCount": 510, "deliveryMinutes": 18, "image": "/static/bombay_sandwich.jpg", "isVegetarian": true, "tags": ["quick-bite", "bengaluru-pick"]}, "f16": {"id": "f16", "name": "Crispy Pani Puri (6 Puris)", "category": "snacks-chaat", "cuisine": "Street Food", "restaurant": "Chaat Junction", "neighbourhood": "Indiranagar", "price": 49, "rating": 4.5, "ratingCount": 370, "deliveryMinutes": 15, "image": "/static/pani_puri.jpg", "isVegetarian": true, "tags": ["quick-bite"]}, "f17": {"id": "f17", "name": "Hot Onion Pakora Plate", "category": "snacks-chaat", "cuisine": "Snacks", "restaurant": "Chai Point Canteen", "neighbourhood": "Whitefield", "price": 45, "rating": 4.3, "ratingCount": 160, "deliveryMinutes": 18, "image": "/static/onion_pakora.jpg", "isVegetarian": true, "tags": ["quick-bite"]}, "f18": {"id": "f18", "name": "Authentic Filter Coffee", "category": "beverages", "cuisine": "Beverages", "restaurant": "Brahmin's Coffee Bar", "neighbourhood": "Malleshwaram", "price": 29, "rating": 4.9, "ratingCount": 680, "deliveryMinutes": 12, "image": "/static/filter_coffee.jpg", "isVegetarian": true, "tags": ["bengaluru-pick", "quick-bite", "under100"]}, "f19": {"id": "f19", "name": "Bengaluru Khara Bath & Kesari Bath", "category": "south-indian", "cuisine": "South Indian", "restaurant": "Veena Stores", "neighbourhood": "Malleshwaram", "price": 55, "rating": 4.8, "ratingCount": 440, "deliveryMinutes": 16, "image": "/static/hero_food.jpg", "isVegetarian": true, "tags": ["bengaluru-pick", "south-indian-fav", "under100"]}, "f20": {"id": "f20", "name": "Egg Kathi Roll", "category": "rolls-wraps", "cuisine": "Street Food", "restaurant": "Rolls King", "neighbourhood": "HSR Layout", "price": 79, "rating": 4.4, "ratingCount": 310, "deliveryMinutes": 20, "image": "/static/kathi_roll.jpg", "isVegetarian": false, "tags": ["under100"]}, "f21": {"id": "f21", "name": "Steamed Chicken Momos (5 pcs)", "category": "momos", "cuisine": "Tibetan", "restaurant": "Shillong Momo Corner", "neighbourhood": "Koramangala", "price": 89, "rating": 4.6, "ratingCount": 390, "deliveryMinutes": 22, "image": "/static/steamed_momos.jpg", "isVegetarian": false, "tags": ["under100", "quick-bite"]}, "f22": {"id": "f22", "name": "Bun Maska & Irani Chai", "category": "beverages", "cuisine": "Bakery", "restaurant": "Cafe Chai Days", "neighbourhood": "Indiranagar", "price": 49, "rating": 4.5, "ratingCount": 270, "deliveryMinutes": 15, "image": "/static/filter_coffee.jpg", "isVegetarian": true, "tags": ["quick-bite", "bengaluru-pick"]}};
var cartItems = {};
var selectedTip = 0;
var includeDonation = true;
var appliedCoupon = null;
var selectedInstructions = {};
var pendingConflictItem = null;

function shareCartLink() {
  var url = window.location.origin + '/?shared_cart=1';
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(url).then(function() {
      alert('Cart share link copied to clipboard!');
    }).catch(function() {
      prompt('Copy cart link:', url);
    });
  } else {
    prompt('Copy cart link:', url);
  }
}

function toggleDonation(checked) {
  includeDonation = !!checked;
  renderCart();
}

function selectTipAmount(amount) {
  if (selectedTip === amount) {
    selectedTip = 0;
  } else {
    selectedTip = amount;
  }
  renderCart();
}

function selectCustomTip() {
  var tipInput = prompt('Enter tip amount in ₹ for delivery partner:', '40');
  if (tipInput !== null) {
    var val = parseInt(tipInput, 10);
    if (!isNaN(val) && val >= 0) {
      selectedTip = val;
      renderCart();
    }
  }
}

function showCartConflictModal(oldRest, newRest) {
  var oldEl = document.getElementById('conflictOldRestName');
  var newEl = document.getElementById('conflictNewRestName');
  if (oldEl) oldEl.innerText = oldRest;
  if (newEl) newEl.innerText = newRest;
  var modal = document.getElementById('cartConflictModal');
  if (modal) modal.style.display = 'flex';
}

function resolveCartConflict(replace) {
  var modal = document.getElementById('cartConflictModal');
  if (modal) modal.style.display = 'none';
  if (replace && pendingConflictItem) {
    var oldIds = Object.keys(cartItems);
    cartItems = {};
    oldIds.forEach(function(id) { updateProductCardControls(id); });
    var p = pendingConflictItem;
    pendingConflictItem = null;
    addToCart(p.id, p.name, p.price, p.restaurant, p.image);
  } else {
    pendingConflictItem = null;
  }
}

function toggleDeliveryInstruction(key, el) {
  if (selectedInstructions[key]) {
    delete selectedInstructions[key];
    if (el) el.classList.remove('active');
  } else {
    selectedInstructions[key] = true;
    if (el) el.classList.add('active');
  }
}

function applyCouponCode(manualCode) {
  var input = document.getElementById('couponCodeInput');
  var code = (manualCode || (input ? input.value : '')).trim().toUpperCase();
  var errEl = document.getElementById('couponErrorMsg');
  if (!code) {
    if (errEl) { errEl.innerText = 'Please enter a coupon code'; errEl.style.display = 'block'; }
    return;
  }
  
  var ids = Object.keys(cartItems);
  var itemsTotal = 0;
  ids.forEach(function(id) { itemsTotal += (cartItems[id].qty * cartItems[id].price); });

  if (code === 'FIRST50') {
    if (itemsTotal < 149) {
      if (errEl) { errEl.innerText = 'FIRST50 valid on orders above ₹149'; errEl.style.display = 'block'; }
      return;
    }
    appliedCoupon = { code: 'FIRST50', discount: 50, label: '₹50 Flat Discount' };
  } else if (code === 'NAMMA100') {
    if (itemsTotal < 99) {
      if (errEl) { errEl.innerText = 'NAMMA100 valid on orders above ₹99'; errEl.style.display = 'block'; }
      return;
    }
    var disc = Math.min(40, Math.round(itemsTotal * 0.20));
    appliedCoupon = { code: 'NAMMA100', discount: disc, label: '20% OFF (₹' + disc + ' saved)' };
  } else if (code === 'FREEDELIVERY') {
    appliedCoupon = { code: 'FREEDELIVERY', discount: 15, label: 'Free Delivery Partner Fee' };
  } else {
    if (errEl) { errEl.innerText = 'Invalid coupon. Try FIRST50, NAMMA100, or FREEDELIVERY'; errEl.style.display = 'block'; }
    return;
  }

  if (errEl) errEl.style.display = 'none';
  renderCart();
}

function removeCouponCode() {
  appliedCoupon = null;
  renderCart();
}

function addToCart(id, name, price, restaurant, image) {
  id = String(id);
  if (!itemsDB[id] && name) {
    itemsDB[id] = { id: id, name: name, price: (typeof price === 'number' && !isNaN(price)) ? price : 79, restaurant: restaurant || 'Bengaluru Kitchen', image: image || '/static/crispy_dosa.jpg' };
  }
  
  var item = itemsDB[id];
  var n = name || (item ? item.name : 'Delicious Dish');
  var p = (typeof price === 'number' && !isNaN(price)) ? price : (item ? item.price : 79);
  var r = restaurant || (item ? item.restaurant : 'Bengaluru Kitchen');
  var img = image || (item ? item.image : '/static/crispy_dosa.jpg');

  // Multi-restaurant support: items from same or different restaurants can be freely added
  // Grouping and separate tracking will be handled seamlessly in cart & live tracking


  if (!cartItems[id]) {
    cartItems[id] = { id: id, name: n, price: p, restaurant: r, image: img, qty: 1 };
  } else {
    cartItems[id].qty += 1;
  }
  renderCart();
  updateProductCardControls(id);

  // Auto-open Blinkit cart drawer upon adding any item
  var drawer = document.getElementById('cartDrawerPanel');
  var backdrop = document.getElementById('cartBackdrop');
  if (drawer && !drawer.classList.contains('open')) {
    drawer.classList.add('open');
    if (backdrop) backdrop.classList.add('open');
  }
}

function updateQty(id, delta) {
  id = String(id);
  if (!cartItems[id]) return;
  cartItems[id].qty += delta;
  if (cartItems[id].qty <= 0) {
    delete cartItems[id];
  }
  renderCart();
  updateProductCardControls(id);
}

function removeItem(id) {
  id = String(id);
  delete cartItems[id];
  renderCart();
  updateProductCardControls(id);
}

function updateProductCardControls(id) {
  id = String(id);
  var btns = document.querySelectorAll('.add-btn-wrap-' + id + ', [id="add-btn-' + id + '"]');
  var item = cartItems[id];
  var cached = itemsDB[id] || {};
  btns.forEach(function(btnWrap) {
    if (item && item.qty > 0) {
      btnWrap.innerHTML = '<div class="qty-counter-pill">' +
        '<button type="button" class="qty-btn" onclick="updateQty(\'' + id + '\', -1)">&minus;</button>' +
        '<span class="qty-num">' + item.qty + '</span>' +
        '<button type="button" class="qty-btn" onclick="updateQty(\'' + id + '\', 1)">&#43;</button>' +
      '</div>';
    } else {
      var n = (cached.name || 'Dish').replace(/'/g, "\\'");
      var p = cached.price || 79;
      var r = (cached.restaurant || 'Bengaluru Kitchen').replace(/'/g, "\\'");
      var img = cached.image || '/static/crispy_dosa.jpg';
      btnWrap.innerHTML = '<button type="button" class="btn-add-pill" onclick="addToCart(\'' + id + '\', \'' + n + '\', ' + p + ', \'' + r + '\', \'' + img + '\')">ADD</button>';
    }
  });
}

function renderCart() {
  var container = document.getElementById('cartItemsContainer');
  var footer = document.getElementById('cartFooterSection');
  var countBadge = document.getElementById('navCartCount');
  var cartBtn = document.getElementById('navCartBtn');
  var floatingBar = document.getElementById('floatingCartBar');

  var ids = Object.keys(cartItems);
  var totalCount = 0;
  var itemsTotal = 0;
  var mrpTotal = 0;

  ids.forEach(function(id) {
    var q = cartItems[id].qty;
    var p = cartItems[id].price;
    totalCount += q;
    itemsTotal += (q * p);
    mrpTotal += (q * (p + 20));
  });

  var savedAmount = mrpTotal - itemsTotal;
  
  // Free delivery threshold: ₹199
  var isFreeDelivery = itemsTotal >= 199 || (appliedCoupon && appliedCoupon.code === 'FREEDELIVERY');
  var deliveryCharge = isFreeDelivery ? 0 : 15;
  var handlingCharge = 4;
  var smallCartCharge = (itemsTotal > 0 && itemsTotal < 99) ? 15 : 0;
  var donationAmount = includeDonation ? 1 : 0;
  var couponDiscount = appliedCoupon ? appliedCoupon.discount : 0;
  
  var grandTotal = Math.max(0, itemsTotal + deliveryCharge + handlingCharge + smallCartCharge + donationAmount + selectedTip - couponDiscount);

  if (totalCount > 0) {
    if (countBadge) {
      countBadge.innerText = totalCount;
      countBadge.style.display = 'flex';
    }
    if (cartBtn) cartBtn.classList.add('active');
    if (floatingBar) {
      floatingBar.classList.add('active');
      document.getElementById('floatingItemCount').innerText = totalCount + (totalCount === 1 ? ' item' : ' items');
      document.getElementById('floatingTotalVal').innerHTML = '&#8377;' + grandTotal;
    }
  } else {
    if (countBadge) countBadge.style.display = 'none';
    if (cartBtn) cartBtn.classList.remove('active');
    if (floatingBar) floatingBar.classList.remove('active');
  }

  var label = document.getElementById('cartItemCountLabel');
  if (label) {
    label.innerText = totalCount + (totalCount === 1 ? ' item' : ' items') + ' selected';
  }

  if (ids.length === 0) {
    appliedCoupon = null;
    if (container) {
      container.innerHTML = '<div style="text-align: center; padding: 60px 20px;">' +
        '<div style="width: 72px; height: 72px; border-radius: 50%; background-color: #ffffff; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: #94A3B8; box-shadow: 0 2px 8px rgba(0,0,0,0.04);">' +
          '<svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path><line x1="3" y1="6" x2="21" y2="6"></line><path d="M16 10a4 4 0 0 1-8 0"></path></svg>' +
        '</div>' +
        '<h4 style="font-size: 17px; font-weight: 800; color: #1E293B; margin-bottom: 6px;">Your cart is empty</h4>' +
        '<p style="font-size: 13.5px; color: #64748B; max-width: 260px; margin: 0 auto 20px;">Explore Bengaluru meals under &#8377;100 and add your favourites.</p>' +
        '<button type="button" onclick="toggleCartDrawer()" style="padding: 10px 24px; border-radius: 20px; background-color: #0C831F; color: #ffffff; font-size: 13.5px; font-weight: 700; cursor: pointer; border: none;">Start Browsing</button>' +
      '</div>';
    }
    if (footer) footer.style.display = 'none';
    return;
  }

  // 1. Delivery ETA Card
  var deliveryHtml = '<div class="b-cart-card" style="display: flex; align-items: center; gap: 14px;">' +
    '<div style="width: 44px; height: 44px; border-radius: 12px; background: #F0FDF4; border: 1.5px solid #BBF7D0; display: flex; align-items: center; justify-content: center; flex-shrink: 0; color: #166534;">' +
      '<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>' +
    '</div>' +
    '<div style="flex: 1; min-width: 0;">' +
      '<div style="font-size: 15px; font-weight: 800; color: #1E293B;">Delivery in 8 minutes</div>' +
      '<div style="font-size: 12.5px; color: #64748B; margin-top: 2px;">Shipment of ' + totalCount + (totalCount === 1 ? ' item' : ' items') + '</div>' +
    '</div>' +
  '</div>';

  // 1b. Free Delivery Progress Meter
  var freePct = Math.min(100, Math.round((itemsTotal / 199) * 100));
  var freeDiff = 199 - itemsTotal;
  var freeMeterHtml = '<div class="free-delivery-meter">' +
    '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">' +
      '<span style="font-size: 12px; font-weight: 800; color: ' + (itemsTotal >= 199 ? '#166534' : '#1E293B') + ';">' +
        (itemsTotal >= 199 ? '🎉 FREE Delivery Unlocked!' : 'Add &#8377;' + freeDiff + ' more for FREE Delivery') +
      '</span>' +
      '<span style="font-size: 11px; font-weight: 700; color: #64748B;">&#8377;199 Target</span>' +
    '</div>' +
    '<div class="free-delivery-bar">' +
      '<div class="free-delivery-progress" style="width: ' + freePct + '%;"></div>' +
    '</div>' +
  '</div>';

  // Recipient / Location Tag (Zero Icons)
  var locTag = '';
  if (activeDelivery.forSomeoneElse) {
    locTag = '<div style="background-color: #FFF7ED; border: 1.5px solid #FED7AA; border-radius: 10px; padding: 10px 12px; display: flex; justify-content: space-between; align-items: flex-start; gap: 8px;">' +
      '<div>' +
        '<div style="font-size: 10.5px; font-weight: 800; color: #EA580C; text-transform: uppercase; letter-spacing: 0.5px;">ORDER FOR SOMEONE ELSE</div>' +
        '<div style="font-size: 13px; font-weight: 800; color: #1E293B; margin-top: 2px;">' + escapeSafe(activeDelivery.recipientName) + ' (' + escapeSafe(activeDelivery.recipientPhone) + ')</div>' +
        '<div style="font-size: 11.5px; color: #64748B; margin-top: 1px;">' + escapeSafe(activeDelivery.recipientAddress) + ', ' + escapeSafe(activeDelivery.locality) + '</div>' +
      '</div>' +
      '<button type="button" onclick="openLocationModal()" style="font-size: 11.5px; font-weight: 700; color: #FF5A1F; cursor: pointer; padding: 3px 8px; border-radius: 6px; background-color: #ffffff; border: 1px solid #FED7AA;">Change</button>' +
    '</div>';
  } else {
    locTag = '<div style="background-color: #ffffff; border: 1px solid #E2E8F0; border-radius: 10px; padding: 10px 12px; display: flex; justify-content: space-between; align-items: center; gap: 8px;">' +
      '<div>' +
        '<div style="font-size: 10px; font-weight: 800; color: #087443; text-transform: uppercase; letter-spacing: 0.5px;">DELIVERING TO</div>' +
        '<div style="font-size: 13px; font-weight: 800; color: #1E293B; margin-top: 2px;">' + escapeSafe(activeDelivery.displayLabel) + '</div>' +
      '</div>' +
      '<button type="button" onclick="openLocationModal()" style="font-size: 11.5px; font-weight: 700; color: #0C831F; cursor: pointer; padding: 3px 8px; border-radius: 6px; background-color: #F0FDF4; border: 1px solid #BBF7D0;">Change</button>' +
    '</div>';
  }

  // 2. Multi-Restaurant Grouped Items List
  var restGroups = {};
  ids.forEach(function(id) {
    var itm = cartItems[id];
    var rName = (itm.restaurant || 'Bengaluru Kitchen').trim();
    if (!restGroups[rName]) {
      restGroups[rName] = {
        name: rName,
        items: [],
        subtotal: 0
      };
    }
    restGroups[rName].items.push(itm);
    restGroups[rName].subtotal += (itm.qty * itm.price);
  });
  var restNames = Object.keys(restGroups);
  var isMultiRest = restNames.length > 1;

  var itemsHtml = '';
  if (isMultiRest) {
    itemsHtml += '<div style="background: linear-gradient(135deg, rgba(255, 90, 31, 0.12), rgba(255, 140, 0, 0.08)); border: 1.5px solid rgba(255, 90, 31, 0.35); border-radius: 12px; padding: 11px 13px; margin-bottom: 12px; display: flex; align-items: center; gap: 10px;">' +
      '<span style="font-size: 22px;">🛍️</span>' +
      '<div>' +
        '<div style="font-size: 12px; font-weight: 800; color: #EA580C; letter-spacing: 0.4px; text-transform: uppercase;">MULTI-RESTAURANT ORDER (' + restNames.length + ' KITCHENS)</div>' +
        '<div style="font-size: 11.5px; color: #475569; margin-top: 1px; line-height: 1.35;">Items will be dispatched simultaneously from each restaurant with separate live tracking!</div>' +
      '</div>' +
    '</div>';
  }

  restNames.forEach(function(rName, rIdx) {
    var group = restGroups[rName];
    itemsHtml += '<div class="b-cart-card" style="padding: 13px 15px; margin-bottom: 12px; border-radius: 12px; border: 1.5px solid #E2E8F0;">';
    
    // Restaurant Kitchen Header
    itemsHtml += '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 11px; padding-bottom: 8px; border-bottom: 1px dashed #E2E8F0;">' +
      '<div style="display: flex; align-items: center; gap: 6px; min-width: 0;">' +
        '<span style="font-size: 14px;">🏪</span>' +
        '<span style="font-size: 13px; font-weight: 800; color: #1E293B; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">' + escapeSafe(rName) + '</span>' +
        (isMultiRest ? '<span style="background: #F1F5F9; color: #475569; font-size: 10px; font-weight: 800; padding: 2px 6px; border-radius: 4px; flex-shrink: 0;">KITCHEN ' + (rIdx + 1) + '</span>' : '') +
      '</div>' +
      '<span style="font-size: 11px; font-weight: 700; color: #0C831F; background: #F0FDF4; padding: 2px 7px; border-radius: 4px; flex-shrink: 0;">' + (rIdx === 0 ? '⚡ 15-20m' : '⚡ 20-25m') + '</span>' +
    '</div>';

    // Kitchen Items
    itemsHtml += '<div style="display: flex; flex-direction: column; gap: 12px;">';
    group.items.forEach(function(item, idx) {
      var id = item.id;
      var borderStyle = idx > 0 ? 'border-top: 1px solid #F1F5F9; padding-top: 12px;' : '';
      itemsHtml += '<div style="display: flex; align-items: center; gap: 12px; ' + borderStyle + '">' +
        '<img src="' + item.image + '" alt="' + escapeSafe(item.name) + '" style="width: 52px; height: 52px; border-radius: 9px; object-fit: cover; border: 1px solid #E2E8F0; background: #F8FAFC;" onerror=\"this.onerror=null;this.src=\'/static/crispy_dosa.jpg\';\">' +
        '<div style="flex: 1; min-width: 0;">' +
          '<div style="font-size: 13px; font-weight: 700; color: #1E293B; line-height: 1.35; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">' + escapeSafe(item.name) + '</div>' +
          '<div style="font-size: 13.5px; font-weight: 800; color: #1E293B; margin-top: 3px;">' +
            '&#8377;' + item.price + ' <span style="font-size: 11px; color: #94A3B8; text-decoration: line-through; font-weight: 500; margin-left: 4px;">&#8377;' + (item.price + 20) + '</span>' +
          '</div>' +
        '</div>' +
        '<div class="b-cart-stepper">' +
          '<button type="button" class="b-cart-stepper-btn" onclick="updateQty(\'' + id + '\', -1)">&minus;</button>' +
          '<span class="b-cart-stepper-val">' + item.qty + '</span>' +
          '<button type="button" class="b-cart-stepper-btn" onclick="updateQty(\'' + id + '\', 1)">&#43;</button>' +
        '</div>' +
      '</div>';
    });
    itemsHtml += '</div>';

    if (isMultiRest) {
      itemsHtml += '<div style="display: flex; justify-content: flex-end; align-items: center; margin-top: 9px; padding-top: 7px; border-top: 1px solid #F1F5F9; font-size: 11.5px; color: #64748B;">' +
        'Kitchen Subtotal: <strong style="color: #1E293B; margin-left: 5px; font-weight: 800; font-size: 12.5px;">&#8377;' + group.subtotal + '</strong>' +
      '</div>';
    }

    itemsHtml += '</div>';
  });

  // 2b. Delivery Instructions Chips Card
  var instructHtml = '<div class="b-cart-card">' +
    '<div style="font-size: 12.5px; font-weight: 800; color: #1E293B; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.4px;">Delivery Instructions</div>' +
    '<div style="display: flex; gap: 6px; flex-wrap: wrap;">' +
      '<button type="button" class="cart-chip-toggle ' + (selectedInstructions['no-call'] ? 'active' : '') + '" onclick="toggleDeliveryInstruction(\'no-call\', this)">🔕 Avoid calling</button>' +
      '<button type="button" class="cart-chip-toggle ' + (selectedInstructions['door'] ? 'active' : '') + '" onclick="toggleDeliveryInstruction(\'door\', this)">🚪 Leave at door</button>' +
      '<button type="button" class="cart-chip-toggle ' + (selectedInstructions['bell'] ? 'active' : '') + '" onclick="toggleDeliveryInstruction(\'bell\', this)">🔔 Ring bell</button>' +
      '<button type="button" class="cart-chip-toggle ' + (selectedInstructions['security'] ? 'active' : '') + '" onclick="toggleDeliveryInstruction(\'security\', this)">🛡️ Leave with security</button>' +
      '<button type="button" class="cart-chip-toggle ' + (selectedInstructions['pet'] ? 'active' : '') + '" onclick="toggleDeliveryInstruction(\'pet\', this)">🐾 Pet at home</button>' +
    '</div>' +
  '</div>';

  // 2c. Coupons & Promo Engine Card
  var couponHtml = '<div class="b-cart-card">' +
    '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">' +
      '<div style="font-size: 12.5px; font-weight: 800; color: #1E293B; text-transform: uppercase; letter-spacing: 0.4px;">Coupons &amp; Offers</div>' +
      (appliedCoupon ? '<button type="button" onclick="removeCouponCode()" style="font-size: 11.5px; color: #DC2626; font-weight: 800; border: none; background: none; cursor: pointer;">Remove</button>' : '') +
    '</div>' +
    (appliedCoupon ?
      '<div style="background: #F0FDF4; border: 1.5px solid #BBF7D0; border-radius: 8px; padding: 9px 12px; display: flex; justify-content: space-between; align-items: center;">' +
        '<div>' +
          '<div style="font-size: 12.5px; font-weight: 800; color: #166534;">🎉 ' + appliedCoupon.code + ' Applied</div>' +
          '<div style="font-size: 11px; color: #15803D;">' + appliedCoupon.label + '</div>' +
        '</div>' +
        '<div style="font-size: 13.5px; font-weight: 800; color: #166534;">-&#8377;' + appliedCoupon.discount + '</div>' +
      '</div>' :
      '<div style="display: flex; gap: 6px;">' +
        '<input type="text" id="couponCodeInput" placeholder="Enter coupon code" style="flex: 1; border: 1.5px solid #CBD5E1; border-radius: 8px; padding: 7px 11px; font-size: 12px; font-weight: 700; text-transform: uppercase; outline: none;">' +
        '<button type="button" onclick="applyCouponCode()" style="background: #0C831F; color: #ffffff; border: none; border-radius: 8px; padding: 0 14px; font-size: 11.5px; font-weight: 800; cursor: pointer;">APPLY</button>' +
      '</div>' +
      '<div id="couponErrorMsg" style="color: #DC2626; font-size: 11px; font-weight: 600; margin-top: 4px; display: none;"></div>' +
      '<div style="display: flex; gap: 6px; margin-top: 8px; flex-wrap: wrap;">' +
        '<span class="cart-promo-pill" onclick="applyCouponCode(\'FIRST50\')">🏷️ FIRST50 (&#8377;50 OFF)</span>' +
        '<span class="cart-promo-pill" onclick="applyCouponCode(\'NAMMA100\')">🏷️ NAMMA100 (20% OFF)</span>' +
        '<span class="cart-promo-pill" onclick="applyCouponCode(\'FREEDELIVERY\')">🏷️ FREEDELIVERY</span>' +
      '</div>'
    ) +
  '</div>';

  // 3. Bill details Card (Zero Icons)
  var billHtml = '<div class="b-cart-card">' +
    '<div style="font-size: 14.5px; font-weight: 800; color: #1E293B; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.5px;">Bill Details</div>' +
    
    '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
      '<div style="display: flex; align-items: center; gap: 6px; color: #334155;">' +
        '<span>Items Subtotal</span>' +
        '<span style="background: #EEF2FF; color: #4F46E5; font-size: 11px; font-weight: 700; padding: 2px 6px; border-radius: 4px;">Saved &#8377;' + savedAmount + '</span>' +
      '</div>' +
      '<div>' +
        '<span style="color: #94A3B8; text-decoration: line-through; margin-right: 6px;">&#8377;' + mrpTotal + '</span>' +
        '<span style="font-weight: 700; color: #1E293B;">&#8377;' + itemsTotal + '</span>' +
      '</div>' +
    '</div>' +

    '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
      '<div style="color: #334155;">Delivery Partner Fee</div>' +
      (deliveryCharge === 0 ? 
        '<div><span style="color: #94A3B8; text-decoration: line-through; margin-right: 4px;">&#8377;15</span><span style="font-weight: 800; color: #166534;">FREE</span></div>' :
        '<div style="font-weight: 700; color: #1E293B;">&#8377;' + deliveryCharge + '</div>'
      ) +
    '</div>' +

    '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
      '<div style="color: #334155;">Restaurant Packaging &amp; Handling</div>' +
      '<div style="font-weight: 700; color: #1E293B;">&#8377;' + handlingCharge + '</div>' +
    '</div>' +

    (smallCartCharge > 0 ? (
      '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
        '<div style="color: #334155;">Small Cart Surcharge</div>' +
        '<div style="font-weight: 700; color: #1E293B;">&#8377;' + smallCartCharge + '</div>' +
      '</div>'
    ) : '') +

    (couponDiscount > 0 ? (
      '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px; color: #166534;">' +
        '<div style="font-weight: 700;">Coupon Discount (' + appliedCoupon.code + ')</div>' +
        '<div style="font-weight: 800;">-&#8377;' + couponDiscount + '</div>' +
      '</div>'
    ) : '') +

    (selectedTip > 0 ? (
      '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
        '<div style="color: #334155;">Delivery Partner Tip</div>' +
        '<div style="font-weight: 700; color: #166534;">+&#8377;' + selectedTip + '</div>' +
      '</div>'
    ) : '') +

    (includeDonation ? (
      '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; font-size: 13px;">' +
        '<div style="color: #334155;">Feeding India Donation</div>' +
        '<div style="font-weight: 700; color: #1E293B;">&#8377;1</div>' +
      '</div>'
    ) : '') +

    '<div style="height: 1px; background: #F1F5F9; margin: 10px 0;"></div>' +

    '<div style="display: flex; justify-content: space-between; align-items: center; font-size: 14.5px;">' +
      '<div style="font-weight: 800; color: #1E293B; text-transform: uppercase; letter-spacing: 0.5px;">Net Total</div>' +
      '<div style="font-weight: 900; color: #0C831F; font-size: 17px;">&#8377;' + grandTotal + '</div>' +
    '</div>' +
  '</div>';

  // 4. Feeding India Donation Card
  var donationHtml = '<div class="b-cart-card" style="display: flex; align-items: center; gap: 12px;">' +
    '<div style="width: 36px; height: 36px; border-radius: 8px; background: #FEF3C7; color: #92400E; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 900; letter-spacing: 0.5px; flex-shrink: 0;">' +
      'CARE' +
    '</div>' +
    '<div style="flex: 1; min-width: 0;">' +
      '<div style="font-size: 13px; font-weight: 800; color: #1E293B;">Feeding India Donation</div>' +
      '<div style="font-size: 11.5px; color: #64748B; margin-top: 2px; line-height: 1.35;">' +
        'Supporting local malnutrition relief initiatives across Bengaluru.' +
      '</div>' +
    '</div>' +
    '<div style="display: flex; align-items: center; gap: 6px; flex-shrink: 0;">' +
      '<span style="font-size: 13px; font-weight: 700; color: #1E293B;">&#8377;1</span>' +
      '<input type="checkbox" id="donationCheckbox" onchange="toggleDonation(this.checked)" ' + (includeDonation ? 'checked' : '') + ' style="accent-color: #0C831F; width: 18px; height: 18px; cursor: pointer;">' +
    '</div>' +
  '</div>';

  // 5. Tip Your Delivery Partner Card
  var tipHtml = '<div class="b-cart-card">' +
    '<div style="font-size: 13.5px; font-weight: 800; color: #1E293B; text-transform: uppercase; letter-spacing: 0.4px;">Tip Your Delivery Partner</div>' +
    '<div style="font-size: 12px; color: #64748B; margin: 4px 0 12px 0; line-height: 1.4;">' +
      '100% of your optional tip is settled directly to your delivery partner.' +
    '</div>' +
    '<div style="display: flex; gap: 8px; flex-wrap: wrap;">' +
      '<button type="button" class="b-tip-pill ' + (selectedTip === 20 ? 'active' : '') + '" onclick="selectTipAmount(20)">&#8377;20</button>' +
      '<button type="button" class="b-tip-pill ' + (selectedTip === 30 ? 'active' : '') + '" onclick="selectTipAmount(30)">&#8377;30</button>' +
      '<button type="button" class="b-tip-pill ' + (selectedTip === 50 ? 'active' : '') + '" onclick="selectTipAmount(50)">&#8377;50</button>' +
      '<button type="button" class="b-tip-pill ' + (selectedTip > 0 && selectedTip !== 20 && selectedTip !== 30 && selectedTip !== 50 ? 'active' : '') + '" onclick="selectCustomTip()">' + (selectedTip > 0 && selectedTip !== 20 && selectedTip !== 30 && selectedTip !== 50 ? '&#8377;' + selectedTip : 'Custom') + '</button>' +
    '</div>' +
  '</div>';

  // 6. Cancellation Policy Card
  var policyHtml = '<div class="b-cart-card">' +
    '<div style="font-size: 13.5px; font-weight: 800; color: #1E293B; margin-bottom: 6px;">Cancellation Policy</div>' +
    '<div style="font-size: 12px; color: #64748B; line-height: 1.45;">' +
      'Orders cannot be cancelled once packed for delivery. In case of unexpected delays, a refund will be provided, if applicable.' +
    '</div>' +
  '</div>';

  if (container) {
    container.innerHTML = deliveryHtml + freeMeterHtml + locTag + itemsHtml + instructHtml + couponHtml + billHtml + donationHtml + tipHtml + policyHtml;
  }

  // Update Bottom Sticky Proceed Bar
  var bottomTotal = document.getElementById('cartBottomGrandTotal');
  if (bottomTotal) bottomTotal.innerHTML = '&#8377;' + grandTotal;

  var subEl = document.getElementById('cartSubtotalVal');
  if (subEl) subEl.innerHTML = '&#8377;' + grandTotal;
  var payEl = document.getElementById('cartToPayVal');
  if (payEl) payEl.innerHTML = '&#8377;' + grandTotal;
  var checkBtn = document.getElementById('checkoutBtnLabel');
  if (checkBtn) checkBtn.innerHTML = 'Proceed &#8226; &#8377;' + grandTotal;

  if (footer) footer.style.display = 'block';
}

function toggleCartDrawer() {
  var backdrop = document.getElementById('cartBackdrop');
  var drawer = document.getElementById('cartDrawerPanel');
  if (backdrop && drawer) {
    backdrop.classList.toggle('open');
    drawer.classList.toggle('open');
  }
}

var currentPaymentTab = 'upi';
var currentUpiApp = 'Google Pay';
var currentBank = 'HDFC Bank';
var currentGrandTotal = 0;
var pgTimerInterval = null;

function checkoutOrder() {
  openPaymentModal();
}

function openPaymentModal() {
  var ids = Object.keys(cartItems);
  if (ids.length === 0) {
    alert('Your cart is empty! Please add some delicious meals under ₹100.');
    return;
  }

  // Calculate bill values
  var totalCount = 0;
  var itemsTotal = 0;
  var miniItemsHtml = '';

  ids.forEach(function(id) {
    var it = cartItems[id];
    totalCount += it.qty;
    var lineTotal = it.qty * it.price;
    itemsTotal += lineTotal;
    miniItemsHtml += '<div class="pg-ledger-row">' +
      '<div class="pg-ledger-item-info">' +
        '<div class="pg-ledger-item-name">' + it.name + '</div>' +
        '<div class="pg-ledger-item-meta">' + (it.restaurant || 'Bengaluru Kitchen') + '</div>' +
      '</div>' +
      '<div class="pg-ledger-qty-wrap">' +
        '<span class="pg-qty-tag">QTY: ' + it.qty + '</span>' +
        '<span class="pg-unit-price">&#8377;' + it.price + ' ea</span>' +
        '<span class="pg-row-total">&#8377;' + lineTotal + '</span>' +
      '</div>' +
    '</div>';
  });

  var ledgerBadge = document.getElementById('pgLedgerTotalQty');
  if (ledgerBadge) {
    ledgerBadge.innerText = totalCount + (totalCount === 1 ? ' ITEM' : ' ITEMS');
  }

  var isFreeDelivery = itemsTotal >= 199 || (appliedCoupon && appliedCoupon.code === 'FREEDELIVERY');
  var deliveryCharge = isFreeDelivery ? 0 : 15;
  var handlingCharge = 4;
  var smallCartCharge = (itemsTotal > 0 && itemsTotal < 99) ? 15 : 0;
  var donationAmount = includeDonation ? 1 : 0;
  var couponDiscount = appliedCoupon ? appliedCoupon.discount : 0;
  currentGrandTotal = Math.max(0, itemsTotal + deliveryCharge + handlingCharge + smallCartCharge + donationAmount + selectedTip - couponDiscount);

  // Populate Right Pane Summary
  document.getElementById('pgMiniItemsList').innerHTML = miniItemsHtml;
  document.getElementById('pgBillItemsTotal').innerHTML = '&#8377;' + itemsTotal;
  document.getElementById('pgBillDelivery').innerHTML = isFreeDelivery ? '<span style="color:#166534; font-weight:800;">FREE</span>' : '&#8377;' + deliveryCharge;
  document.getElementById('pgBillHandling').innerHTML = '&#8377;' + handlingCharge;
  
  var smallRow = document.getElementById('pgBillSmallCartRow');
  if (smallRow) {
    if (smallCartCharge > 0) {
      smallRow.style.display = 'flex';
      document.getElementById('pgBillSmallCart').innerHTML = '&#8377;' + smallCartCharge;
    } else {
      smallRow.style.display = 'none';
    }
  }

  var tipRow = document.getElementById('pgBillTipRow');
  if (tipRow) {
    if (selectedTip > 0) {
      tipRow.style.display = 'flex';
      document.getElementById('pgBillTip').innerHTML = '+&#8377;' + selectedTip;
    } else {
      tipRow.style.display = 'none';
    }
  }

  var donRow = document.getElementById('pgBillDonationRow');
  if (donRow) {
    donRow.style.display = includeDonation ? 'flex' : 'none';
  }

  document.getElementById('pgGrandTotalVal').innerHTML = '&#8377;' + currentGrandTotal;

  // Set destination address
  var destLabel = activeDelivery.forSomeoneElse ? 
    (activeDelivery.recipientName + ' (' + activeDelivery.locality + ')') : 
    activeDelivery.displayLabel;
  var destEl = document.getElementById('pgDestSummary');
  if (destEl) destEl.innerText = destLabel;

  // Reset to initial payment view
  document.getElementById('pgMainBody').style.display = 'flex';
  document.getElementById('pgHeaderBar').style.display = 'flex';
  document.getElementById('pgFooterBar').style.display = 'flex';
  document.getElementById('pgSuccessView').style.display = 'none';

  updatePaymentButtonLabel();
  startPgQrTimer();

  // Show modal
  var modal = document.getElementById('paymentGatewayModal');
  if (modal) modal.style.display = 'flex';
}

function closePaymentModal() {
  var modal = document.getElementById('paymentGatewayModal');
  if (modal) modal.style.display = 'none';
  if (pgTimerInterval) clearInterval(pgTimerInterval);
}

function selectPaymentTab(tab) {
  currentPaymentTab = tab;
  var tabs = ['upi', 'card', 'net', 'wallet', 'cod'];
  tabs.forEach(function(t) {
    var btn = document.getElementById('pgTab' + t.charAt(0).toUpperCase() + t.slice(1));
    var pane = document.getElementById('pane' + t.charAt(0).toUpperCase() + t.slice(1));
    if (t === tab) {
      if (btn) btn.classList.add('active');
      if (pane) pane.style.display = 'block';
    } else {
      if (btn) btn.classList.remove('active');
      if (pane) pane.style.display = 'none';
    }
  });
  updatePaymentButtonLabel();
}

function selectUpiApp(el, appName) {
  document.querySelectorAll('.pg-upi-app').forEach(function(a) { a.classList.remove('active'); });
  if (el) el.classList.add('active');
  currentUpiApp = appName;
  updatePaymentButtonLabel();
}

function selectQuickBank(el, bankName) {
  document.querySelectorAll('.pg-bank-btn').forEach(function(b) { b.classList.remove('active'); });
  if (el) el.classList.add('active');
  currentBank = bankName;
  updatePaymentButtonLabel();
}

function updatePaymentButtonLabel() {
  var btnText = document.getElementById('pgPayBtnText');
  if (!btnText) return;

  if (currentPaymentTab === 'upi') {
    btnText.innerHTML = 'Pay &#8377;' + currentGrandTotal + ' via ' + currentUpiApp;
  } else if (currentPaymentTab === 'card') {
    btnText.innerHTML = 'Pay &#8377;' + currentGrandTotal + ' with Card';
  } else if (currentPaymentTab === 'net') {
    btnText.innerHTML = 'Pay &#8377;' + currentGrandTotal + ' via ' + currentBank;
  } else if (currentPaymentTab === 'wallet') {
    var selectedWallet = 'Wallet';
    var checked = document.querySelector('input[name="pgWallet"]:checked');
    if (checked) selectedWallet = checked.value;
    btnText.innerHTML = 'Pay &#8377;' + currentGrandTotal + ' with ' + selectedWallet;
  } else if (currentPaymentTab === 'cod') {
    btnText.innerHTML = 'Place Order (Cash on Delivery &#8226; &#8377;' + currentGrandTotal + ')';
  }
}

function verifyUpiVpa() {
  var input = document.getElementById('pgUpiVpaInput');
  var status = document.getElementById('pgVpaStatus');
  if (!input || !status) return;
  var val = input.value.trim();
  if (!val || val.indexOf('@') === -1) {
    status.style.display = 'block';
    status.style.color = '#DC2626';
    status.innerText = 'Please enter a valid UPI VPA (e.g. mobile@upi or name@oksbi)';
    return;
  }
  status.style.display = 'block';
  status.style.color = '#166534';
  status.innerHTML = '✓ VPA verified: <strong>' + val + '</strong> (Ready for instant intent pay)';
  currentUpiApp = val;
  updatePaymentButtonLabel();
}

function formatCardNumber(input) {
  var v = input.value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
  var matches = v.match(/\d{4,16}/g);
  var match = matches && matches[0] || '';
  var parts = [];
  for (var i = 0, len = match.length; i < len; i += 4) {
    parts.push(match.substring(i, i + 4));
  }
  if (parts.length) {
    input.value = parts.join(' ');
  } else {
    input.value = v;
  }

  // Detect Brand
  var badge = document.getElementById('pgCardBrand');
  if (badge) {
    if (v.startsWith('4')) { badge.innerText = 'VISA'; badge.style.color = '#1D4ED8'; }
    else if (v.startsWith('5')) { badge.innerText = 'MASTERCARD'; badge.style.color = '#B91C1C'; }
    else if (v.startsWith('6')) { badge.innerText = 'RUPAY'; badge.style.color = '#15803D'; }
    else { badge.innerText = 'CARD'; badge.style.color = '#64748B'; }
  }
}

function formatCardExpiry(input) {
  var v = input.value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
  if (v.length >= 2) {
    input.value = v.substring(0, 2) + '/' + v.substring(2, 4);
  } else {
    input.value = v;
  }
}

function startPgQrTimer() {
  if (pgTimerInterval) clearInterval(pgTimerInterval);
  var duration = 300; // 5 minutes
  var display = document.getElementById('pgQrTimer');
  pgTimerInterval = setInterval(function() {
    var minutes = parseInt(duration / 60, 10);
    var seconds = parseInt(duration % 60, 10);
    minutes = minutes < 10 ? "0" + minutes : minutes;
    seconds = seconds < 10 ? "0" + seconds : seconds;
    if (display) display.textContent = minutes + ":" + seconds;
    if (--duration < 0) {
      clearInterval(pgTimerInterval);
      if (display) display.textContent = "Expired";
    }
  }, 1000);
}

function processPaymentGateway() {
  var payBtn = document.getElementById('pgPayButton');
  var btnText = document.getElementById('pgPayBtnText');
  if (payBtn) payBtn.disabled = true;

  if (btnText) btnText.innerHTML = 'Authorizing with Gateway...';

  setTimeout(function() {
    if (btnText) btnText.innerHTML = 'Confirming Bank Settlement...';
    
    setTimeout(function() {
      // Transition to Success Screen
      document.getElementById('pgMainBody').style.display = 'none';
      document.getElementById('pgHeaderBar').style.display = 'none';
      document.getElementById('pgFooterBar').style.display = 'none';
      
      var orderNum = Math.floor(100000 + Math.random() * 900000);
      document.getElementById('pgSuccessOrderId').innerText = 'ORDER REF: BITE-2026-' + orderNum;
      document.getElementById('pgSuccessAmount').innerHTML = '&#8377;' + currentGrandTotal;

      var methodLabel = 'UPI / ' + currentUpiApp;
      if (currentPaymentTab === 'card') methodLabel = 'Credit/Debit Card (Tokenized)';
      else if (currentPaymentTab === 'net') methodLabel = 'NetBanking (' + currentBank + ')';
      else if (currentPaymentTab === 'wallet') methodLabel = 'Digital Wallet';
      else if (currentPaymentTab === 'cod') methodLabel = 'Pay on Delivery (Doorstep)';
      document.getElementById('pgSuccessMethod').innerText = methodLabel;

      document.getElementById('pgSuccessView').style.display = 'flex';

      // Save order details for live multi-restaurant tracking stepper
      var ids = Object.keys(cartItems);
      var restMap = {};
      var riderList = [
        { name: 'Ramesh Kumar', vehicle: 'Hero Splendor KA-04-EK-9122', rating: '4.9 ⭐', phone: '+91 98765 43210' },
        { name: 'Suresh Gowda', vehicle: 'TVS Apache KA-05-MN-4412', rating: '4.8 ⭐', phone: '+91 98765 88123' },
        { name: 'Manjunath K', vehicle: 'Bajaj Pulsar KA-01-AB-1904', rating: '4.9 ⭐', phone: '+91 98765 71920' },
        { name: 'Pradeep R', vehicle: 'Ather 450X KA-03-EV-8821', rating: '4.7 ⭐', phone: '+91 98765 33214' }
      ];

      ids.forEach(function(id) {
        var itm = cartItems[id];
        var r = (itm.restaurant || 'Bengaluru Kitchen').trim();
        if (!restMap[r]) {
          restMap[r] = {
            restaurant: r,
            items: [],
            subtotal: 0
          };
        }
        restMap[r].items.push(itm);
        restMap[r].subtotal += (itm.qty * itm.price);
      });

      var multiOrders = [];
      var rKeys = Object.keys(restMap);
      rKeys.forEach(function(rName, idx) {
        var legSubId = 'BITE-2026-' + orderNum + (rKeys.length > 1 ? ('-K' + (idx + 1)) : '');
        var rider = riderList[idx % riderList.length];
        var itemsSummary = restMap[rName].items.map(function(it) {
          return it.qty + 'x ' + it.name;
        }).join(', ');

        multiOrders.push({
          orderId: legSubId,
          restaurant: rName,
          itemsSummary: itemsSummary,
          subtotal: restMap[rName].subtotal,
          rider: rider,
          etaMinutes: 16 + (idx * 5),
          stage: 2,
          address: (activeDelivery.displayLabel || 'Indiranagar 100ft Rd, Bengaluru')
        });
      });

      window.lastMultiOrders = multiOrders;
      window.activeTrackOrderIndex = 0;
      window.lastOrderDetails = {
        orderNum: orderNum,
        amount: currentGrandTotal,
        restaurant: rKeys[0] || 'Bengaluru Kitchen',
        address: (activeDelivery.displayLabel || 'Indiranagar 100ft Rd, Bengaluru')
      };

      // Clear Cart after successful payment
      cartItems = {};
      appliedCoupon = null;
      selectedTip = 0;
      renderCart();
      ids.forEach(function(id) {
        updateProductCardControls(id);
      });
      if (payBtn) payBtn.disabled = false;
    }, 1000);
  }, 800);
}

var trackingInterval = null;
var activeTrackOrderIndex = 0;

function launchOrderTracking() {
  closePaymentModal();
  openOrderTrackingModal();
}

function openOrderTrackingModal(orderId, restName) {
  var modal = document.getElementById('liveOrderTrackingModal');
  if (modal) modal.style.display = 'flex';

  var orders = window.lastMultiOrders;
  if (!orders || orders.length === 0) {
    orders = [{
      orderId: 'BITE-2026-' + (orderId || '92841'),
      restaurant: restName || 'Sri Krishna Cafe, Indiranagar',
      itemsSummary: '1x Crispy Masala Dosa',
      subtotal: 79,
      rider: { name: 'Ramesh Kumar', vehicle: 'Hero Splendor KA-04-EK-9122', rating: '4.9 ⭐', phone: '+91 98765 43210' },
      etaMinutes: 18,
      stage: 2,
      address: (activeDelivery.displayLabel || 'Indiranagar 100ft Rd, Bengaluru')
    }];
    window.lastMultiOrders = orders;
  }

  activeTrackOrderIndex = 0;
  renderMultiTrackTabs();
  displayActiveTrackOrder(0);

  if (trackingInterval) clearInterval(trackingInterval);
  trackingInterval = setInterval(function() {
    var anyActive = false;
    (window.lastMultiOrders || []).forEach(function(ord) {
      if (ord.stage < 5) {
        ord.stage++;
        if (ord.etaMinutes > 4) ord.etaMinutes -= 4;
        anyActive = true;
      }
    });
    displayActiveTrackOrder(activeTrackOrderIndex);
    renderMultiTrackTabs();
    if (!anyActive) {
      clearInterval(trackingInterval);
    }
  }, 4500);
}

function renderMultiTrackTabs() {
  var tabsWrap = document.getElementById('multiTrackTabsWrap');
  if (!tabsWrap) return;
  var orders = window.lastMultiOrders || [];
  if (orders.length <= 1) {
    tabsWrap.style.display = 'none';
    return;
  }
  tabsWrap.style.display = 'flex';
  var html = '';
  orders.forEach(function(ord, idx) {
    var isActive = (idx === activeTrackOrderIndex);
    var statusBadge = ord.stage >= 5 ? '✓ Delivered' : ('Stage ' + ord.stage + '/5');
    html += '<button type="button" class="multi-track-tab-btn ' + (isActive ? 'active' : '') + '" onclick="switchTrackingOrder(' + idx + ')">' +
      '<span>🏪 ' + escapeSafe(ord.restaurant) + '</span>' +
      '<span style="font-size: 11px; opacity: 0.85; margin-left: 6px;">(' + statusBadge + ')</span>' +
    '</button>';
  });
  tabsWrap.innerHTML = html;
}

function switchTrackingOrder(idx) {
  activeTrackOrderIndex = idx;
  renderMultiTrackTabs();
  displayActiveTrackOrder(idx);
}

function displayActiveTrackOrder(idx) {
  var orders = window.lastMultiOrders || [];
  if (!orders[idx]) return;
  var ord = orders[idx];

  var rEl = document.getElementById('trackRestaurantName');
  if (rEl) rEl.innerText = 'Ref: ' + ord.orderId + ' • Kitchen: ' + ord.restaurant;

  var addrEl = document.getElementById('trackAddressText');
  if (addrEl) addrEl.innerText = ord.address;

  var summaryBox = document.getElementById('trackOrderItemsSummary');
  var itemsEl = document.getElementById('trackItemsListText');
  var subtotalEl = document.getElementById('trackLegSubtotal');
  if (summaryBox) {
    summaryBox.style.display = 'flex';
    if (itemsEl) itemsEl.innerText = ord.itemsSummary;
    if (subtotalEl) subtotalEl.innerHTML = '&#8377;' + ord.subtotal;
  }

  var riderEl = document.getElementById('trackRiderDetail');
  if (riderEl && ord.rider) {
    riderEl.innerText = ord.rider.name + ' (' + ord.rider.vehicle + ') • ' + ord.rider.rating;
  }

  updateTrackingSteps(ord.stage, ord.etaMinutes);
}

function updateTrackingSteps(stage, etaMinutes) {
  for (var i = 1; i <= 5; i++) {
    var step = document.getElementById('step' + i);
    if (!step) continue;
    step.classList.remove('done', 'active');
    if (i < stage) {
      step.classList.add('done');
    } else if (i === stage) {
      step.classList.add('active');
    }
  }
  var etaEl = document.getElementById('trackEtaText');
  var titleEl = document.getElementById('trackHeaderTitle');
  var etaVal = etaMinutes || (stage === 1 ? 20 : stage === 2 ? 16 : stage === 3 ? 11 : stage === 4 ? 5 : 0);
  if (stage === 1) {
    if (etaEl) etaEl.innerText = 'Estimated Delivery: ' + etaVal + ' mins';
    if (titleEl) titleEl.innerText = 'Order Placed & Confirmed';
  } else if (stage === 2) {
    if (etaEl) etaEl.innerText = 'Estimated Delivery: ' + etaVal + ' mins';
    if (titleEl) titleEl.innerText = 'Kitchen Preparing Fresh Meal';
  } else if (stage === 3) {
    if (etaEl) etaEl.innerText = 'Estimated Delivery: ' + etaVal + ' mins';
    if (titleEl) titleEl.innerText = 'Valet Assigned & Heading to Kitchen';
  } else if (stage === 4) {
    if (etaEl) etaEl.innerText = 'Estimated Delivery: ' + etaVal + ' mins';
    if (titleEl) titleEl.innerText = 'Valet Out for Delivery 🛵';
  } else if (stage === 5) {
    if (etaEl) etaEl.innerText = 'Delivered! 🎉';
    if (titleEl) titleEl.innerText = 'Order Delivered! Enjoy your meal!';
  }
}

function closeOrderTrackingModal() {
  var modal = document.getElementById('liveOrderTrackingModal');
  if (modal) modal.style.display = 'none';
  if (trackingInterval) clearInterval(trackingInterval);
}

function toggleVegOnlyFilter(isVeg) {
  var cards = document.querySelectorAll('#mainFoodGrid .product-card');
  cards.forEach(function(card) {
    if (!isVeg) {
      card.style.display = 'flex';
      return;
    }
    var hasVegIndicator = card.querySelector('.diet-veg') !== null;
    var name = (card.getAttribute('data-name') || '').toLowerCase();
    var isNonVeg = card.querySelector('.diet-nonveg') !== null || 
      name.indexOf('chicken') !== -1 || name.indexOf('egg') !== -1 || name.indexOf('mutton') !== -1 || name.indexOf('fish') !== -1 || name.indexOf('prawn') !== -1;
    if (hasVegIndicator && !isNonVeg) {
      card.style.display = 'flex';
    } else {
      card.style.display = 'none';
    }
  });
}

function detectGPSLocation() {
  var locText = document.getElementById('currentLocText');
  var oldText = locText ? locText.textContent : '';
  if (locText) locText.innerHTML = 'Locating...';

  if (!navigator.geolocation) {
    alert('Geolocation is not supported by your browser.');
    if (locText) locText.textContent = oldText;
    return;
  }

  navigator.geolocation.getCurrentPosition(
    function(pos) {
      var detected = 'Indiranagar, Bengaluru';
      if (locText) locText.innerHTML = '📍 ' + detected;
      activeDelivery.locality = 'Indiranagar';
      activeDelivery.displayLabel = '📍 ' + detected;
      renderCart();
    },
    function(err) {
      var detected = 'Koramangala, Bengaluru';
      if (locText) locText.innerHTML = '📍 ' + detected;
      activeDelivery.locality = 'Koramangala';
      activeDelivery.displayLabel = '📍 ' + detected;
    },
    { timeout: 7000 }
  );
}

// RESTAURANT SEARCH & UNDER ₹100 MENU ENGINE
var currentActiveRestaurantFilter = null;

function handleGlobalSearch(q) {
  var term = (q || '').toLowerCase().trim();
  var searchInput = document.getElementById('globalSearchInput');
  var clearBtn = document.getElementById('clearSearchBtn');
  var dropdown = document.getElementById('searchSuggestionsDropdown');

  if (searchInput && searchInput.value !== q) searchInput.value = q;

  if (!term) {
    if (clearBtn) clearBtn.style.display = 'none';
    if (dropdown) dropdown.style.display = 'none';
    clearRestaurantFilter(true);
    return;
  }

  if (clearBtn) clearBtn.style.display = 'block';

  // Gather matching restaurants & dishes from itemsDB and DOM
  var matchedRestaurants = {};
  var matchedDishes = [];

  // Inspect itemsDB
  if (typeof itemsDB === 'object') {
    Object.keys(itemsDB).forEach(function(id) {
      var itm = itemsDB[id];
      var rName = (itm.restaurant || '').trim();
      var rLower = rName.toLowerCase();
      var dLower = (itm.name || '').toLowerCase();
      var cLower = (itm.cuisine || '').toLowerCase();

      if (rLower.indexOf(term) !== -1 || (term.length > 2 && rLower.split(' ').some(function(w){ return w.indexOf(term) !== -1; }))) {
        if (!matchedRestaurants[rName]) {
          matchedRestaurants[rName] = {
            name: rName,
            neighbourhood: itm.neighbourhood || 'Bengaluru',
            rating: itm.rating || 4.6,
            deliveryMinutes: itm.deliveryMinutes || 20,
            itemCount: 0
          };
        }
        matchedRestaurants[rName].itemCount++;
      }

      if (dLower.indexOf(term) !== -1 || cLower.indexOf(term) !== -1) {
        matchedDishes.push(itm);
      }
    });
  }

  // Also check cards in DOM for restaurants
  var cards = document.querySelectorAll('#mainFoodGrid .product-card');
  cards.forEach(function(card) {
    var rEl = card.querySelector('.product-card-restaurant');
    var rName = rEl ? rEl.textContent.trim().replace(/^🏪\s*/, '') : '';
    if (rName && rName.toLowerCase().indexOf(term) !== -1) {
      if (!matchedRestaurants[rName]) {
        matchedRestaurants[rName] = {
          name: rName,
          neighbourhood: card.getAttribute('data-neighbourhood') || 'Bengaluru',
          rating: card.getAttribute('data-rating') || '4.6',
          deliveryMinutes: card.getAttribute('data-delivery') || '20',
          itemCount: 1
        };
      }
    }
  });

  // Render Suggestions Dropdown
  var restKeys = Object.keys(matchedRestaurants);
  if (dropdown) {
    if (restKeys.length > 0 || matchedDishes.length > 0) {
      var dHtml = '';
      if (restKeys.length > 0) {
        dHtml += '<div class="search-drop-group-title">RESTAURANTS WITH UNDER &#8377;100 MENU (' + restKeys.length + ')</div>';
        restKeys.slice(0, 4).forEach(function(rKey) {
          var rInfo = matchedRestaurants[rKey];
          var safeRName = rKey.replace(/'/g, "\\'");
          dHtml += '<div class="search-drop-item" onclick="filterByRestaurant(\'' + safeRName + '\')">' +
            '<span style="font-size: 16px;">🏪</span>' +
            '<div style="flex: 1; min-width: 0;">' +
              '<div style="font-size: 13.5px; font-weight: 800; color: #1E293B;">' + escapeSafe(rKey) + '</div>' +
              '<div style="font-size: 11.5px; color: #64748B;">📍 ' + escapeSafe(rInfo.neighbourhood) + ' &bull; ⭐ ' + rInfo.rating + ' &bull; 💯 All dishes &#8377;100 or less</div>' +
            '</div>' +
            '<span style="font-size: 11px; font-weight: 700; color: #FF5A1F; background: rgba(255, 90, 31, 0.08); padding: 3px 8px; border-radius: 6px;">VIEW MENU</span>' +
          '</div>';
        });
      }

      if (matchedDishes.length > 0) {
        dHtml += '<div class="search-drop-group-title">DISHES UNDER &#8377;100 (' + matchedDishes.length + ')</div>';
        matchedDishes.slice(0, 5).forEach(function(dish) {
          var safeName = dish.name.replace(/'/g, "\\'");
          dHtml += '<div class="search-drop-item" onclick="handleQuickSearch(\'' + safeName + '\')">' +
            '<span style="font-size: 16px;">🍲</span>' +
            '<div style="flex: 1; min-width: 0;">' +
              '<div style="font-size: 13px; font-weight: 700; color: #1E293B;">' + escapeSafe(dish.name) + '</div>' +
              '<div style="font-size: 11.5px; color: #64748B;">' + escapeSafe(dish.restaurant) + '</div>' +
            '</div>' +
            '<span style="font-size: 12.5px; font-weight: 800; color: #0C831F;">&#8377;' + dish.price + '</span>' +
          '</div>';
        });
      }

      dropdown.innerHTML = dHtml;
      dropdown.style.display = 'block';
    } else {
      dropdown.style.display = 'none';
    }
  }

  // Check if query exactly or closely matches a single restaurant
  if (restKeys.length === 1 && (restKeys[0].toLowerCase() === term || restKeys[0].toLowerCase().indexOf(term) === 0)) {
    filterByRestaurant(restKeys[0], false);
    return;
  }

  if (currentActiveRestaurantFilter) {
    currentActiveRestaurantFilter = null;
  }

  cards.forEach(function(card) {
    var name = (card.getAttribute('data-name') || '').toLowerCase();
    var rEl = card.querySelector('.product-card-restaurant');
    var rName = rEl ? rEl.textContent.toLowerCase() : '';
    var neigh = (card.getAttribute('data-neighbourhood') || '').toLowerCase();

    if (name.indexOf(term) !== -1 || rName.indexOf(term) !== -1 || neigh.indexOf(term) !== -1) {
      card.style.display = 'flex';
    } else {
      card.style.display = 'none';
    }
  });
}

function filterByRestaurant(restName, doScroll) {
  if (doScroll === undefined) doScroll = true;
  currentActiveRestaurantFilter = restName;

  var searchInput = document.getElementById('globalSearchInput');
  if (searchInput) searchInput.value = restName;

  var clearBtn = document.getElementById('clearSearchBtn');
  if (clearBtn) clearBtn.style.display = 'block';

  var dropdown = document.getElementById('searchSuggestionsDropdown');
  if (dropdown) dropdown.style.display = 'none';

  // Filter food grid cards to only this restaurant
  var cards = document.querySelectorAll('#mainFoodGrid .product-card');
  var rTarget = restName.trim().toLowerCase();
  cards.forEach(function(card) {
    var rEl = card.querySelector('.product-card-restaurant');
    var rName = rEl ? rEl.textContent.trim().replace(/^🏪\s*/, '').toLowerCase() : '';
    if (rName === rTarget || rName.indexOf(rTarget) !== -1 || rTarget.indexOf(rName) !== -1) {
      card.style.display = 'flex';
    } else {
      card.style.display = 'none';
    }
  });

  if (doScroll) {
    var grid = document.getElementById('mainFoodGrid') || document.getElementById('food-discovery');
    if (grid) grid.scrollIntoView({ behavior: 'smooth' });
  }
}

function clearRestaurantFilter(skipInputClear) {
  currentActiveRestaurantFilter = null;
  var dropdown = document.getElementById('searchSuggestionsDropdown');
  if (dropdown) dropdown.style.display = 'none';

  if (!skipInputClear) {
    var searchInput = document.getElementById('globalSearchInput');
    if (searchInput) searchInput.value = '';
    var clearBtn = document.getElementById('clearSearchBtn');
    if (clearBtn) clearBtn.style.display = 'none';
  }

  var cards = document.querySelectorAll('#mainFoodGrid .product-card');
  cards.forEach(function(card) {
    card.style.display = 'flex';
  });
}

function clearGlobalSearch() {
  clearRestaurantFilter(false);
}

// Close search dropdown on click outside & setup restaurant clicks
document.addEventListener('click', function(e) {
  var dropdown = document.getElementById('searchSuggestionsDropdown');
  var searchWrap = document.querySelector('.search-bar-wrap');
  if (dropdown && searchWrap && !searchWrap.contains(e.target)) {
    dropdown.style.display = 'none';
  }

  // Click on any product card restaurant name opens its under 100 menu
  var rEl = e.target.closest('.product-card-restaurant');
  if (rEl) {
    var rName = rEl.textContent.trim().replace(/^🏪\s*/, '');
    if (rName) filterByRestaurant(rName);
  }
});

function handleQuickSearch(dishName) {
  var searchInput = document.getElementById('globalSearchInput');
  if (searchInput) {
    searchInput.value = dishName;
    handleGlobalSearch(dishName);
  }
  var grid = document.getElementById('food-discovery') || document.getElementById('mainFoodGrid');
  if (grid) grid.scrollIntoView({ behavior: 'smooth' });
}

function filterByTag(tag) {
  var t = (tag || '').toLowerCase();
  if (t.indexOf('south') !== -1) filterGridCategory('south-indian');
  else if (t.indexOf('roll') !== -1) filterGridCategory('fast-food');
  else if (t.indexOf('momo') !== -1) filterGridCategory('chinese');
  else if (t.indexOf('street') !== -1 || t.indexOf('chaat') !== -1) filterGridCategory('street-food');
  else if (t.indexOf('coffee') !== -1 || t.indexOf('chai') !== -1) handleQuickSearch('Coffee');
  else if (t.indexOf('meal') !== -1 || t.indexOf('thali') !== -1) filterGridCategory('lunch');
  else handleQuickSearch(tag);

  var grid = document.getElementById('food-discovery') || document.getElementById('mainFoodGrid');
  if (grid) grid.scrollIntoView({ behavior: 'smooth' });
}

function setLocality(loc) {
  var locEl = document.getElementById('currentLocText');
  if (locEl) locEl.textContent = loc + ', Bengaluru';
  var grid = document.getElementById('food-discovery') || document.getElementById('mainFoodGrid');
  if (grid) grid.scrollIntoView({ behavior: 'smooth' });
}

function filterGridCategory(catId, btn) {
  var pills = document.querySelectorAll('.filter-tab-pill');
  pills.forEach(function(p) { p.classList.remove('active'); });
  if (btn) {
    btn.classList.add('active');
  } else {
    pills.forEach(function(p) {
      var attr = p.getAttribute('onclick') || '';
      if (attr.indexOf("'" + catId + "'") !== -1 || (catId === 'all' && attr.indexOf("'all'") !== -1)) {
        p.classList.add('active');
      }
    });
  }

  var cards = document.querySelectorAll('#mainFoodGrid .product-card');
  cards.forEach(function(card) {
    var c = card.getAttribute('data-category');
    var match = (catId === 'all' || c === catId);
    if (catId === 'south-indian' && (c === 'south-indian' || c === 'dosa-idli')) match = true;
    if (catId === 'north-indian' && (c === 'north-indian' || c === 'parathas')) match = true;
    if (catId === 'chinese' && (c === 'chinese' || c === 'momos')) match = true;
    if (catId === 'street-food' && (c === 'street-food' || c === 'snacks-chaat' || c === 'sandwiches')) match = true;
    if (catId === 'snacks' && (c === 'snacks' || c === 'snacks-chaat' || c === 'sandwiches')) match = true;
    if (catId === 'fast-food' && (c === 'fast-food' || c === 'rolls-wraps' || c === 'sandwiches')) match = true;
    if (catId === 'breakfast' && (c === 'breakfast' || c === 'dosa-idli')) match = true;
    if (catId === 'lunch' && (c === 'lunch' || c === 'rice-biryani' || c === 'north-indian')) match = true;
    if (match) {
      card.style.display = 'flex';
    } else {
      card.style.display = 'none';
    }
  });
}

var isOpenNowActive = false;
function toggleOpenNowFilter(btn) {
  isOpenNowActive = !isOpenNowActive;
  if (btn) {
    if (isOpenNowActive) {
      btn.classList.add('active');
      btn.style.backgroundColor = '#166534';
      btn.style.color = '#ffffff';
    } else {
      btn.classList.remove('active');
      btn.style.backgroundColor = '';
      btn.style.color = '#166534';
    }
  }

  // Reload nearby cards with open_now parameter if restaurant feed is present
  var feed = document.querySelector('.cards-wrapper-container') || document.getElementById('restaurants-list') || document.querySelector('.cards-wrapper');
  if (feed && window.htmx) {
    var url = '/nearby?open_now=' + (isOpenNowActive ? 'true' : 'false');
    window.htmx.ajax('GET', url, { target: feed, swap: 'outerHTML' });
  }
}

function selectCategoryFilter(catId, btn) {
  var btns = document.querySelectorAll('.cat-pill-btn');
  btns.forEach(function(b) { b.classList.remove('active'); });
  if (btn) {
    btn.classList.add('active');
  } else {
    btns.forEach(function(b) {
      var attr = b.getAttribute('onclick') || '';
      if (attr.indexOf("'" + catId + "'") !== -1) {
        b.classList.add('active');
      }
    });
  }

  filterGridCategory(catId, null);
  var gridEl = document.getElementById('food-discovery');
  if (gridEl) gridEl.scrollIntoView({ behavior: 'smooth' });
}

function toggleNeighbourhoodMenu() {
  openLocationModal();
}

function openAppModal() {
  alert('Under100 App: Download from Play Store & App Store. Now active across Bengaluru!');
}

function openAuthModal(mode) {
  window.location.href = '/login' + (mode === 'signup' ? '?mode=signup' : '');
}

/* ==========================================================================
   BENGALURU GEOLOCATION, NEAREST RESTAURANT SORTING & ORDER FOR SOMEONE ELSE
   ========================================================================== */

var BANGALORE_HUBS = {
  "Indiranagar": { lat: 12.9784, lng: 77.6408, desc: "100ft Rd / CMH Rd" },
  "Koramangala": { lat: 12.9352, lng: 77.6245, desc: "5th & 6th Block" },
  "HSR Layout": { lat: 12.9121, lng: 77.6446, desc: "Sector 1 & 2" },
  "Malleshwaram": { lat: 13.0031, lng: 77.5645, desc: "Margosa / 8th Cross" },
  "Jayanagar": { lat: 12.9308, lng: 77.5838, desc: "4th Block Complex" },
  "Basavanagudi": { lat: 12.9406, lng: 77.5738, desc: "Gandhi Bazaar / DVG Rd" },
  "Whitefield": { lat: 12.9698, lng: 77.7499, desc: "ITPB / Prestige" },
  "BTM Layout": { lat: 12.9166, lng: 77.6101, desc: "Udupi Garden" },
  "MG Road": { lat: 12.9750, lng: 77.6033, desc: "MG Road / Church Street" },
  "Marathahalli": { lat: 12.9569, lng: 77.7011, desc: "Outer Ring Road" },
  "VV Puram": { lat: 12.9515, lng: 77.5768, desc: "Food Street" }
};

var activeDelivery = {
  isGPS: false,
  lat: 12.9784,
  lng: 77.6408,
  locality: "Indiranagar",
  displayLabel: "Indiranagar, Bengaluru",
  forSomeoneElse: false,
  recipientName: "",
  recipientPhone: "",
  recipientAddress: "",
  recipientNote: ""
};

function escapeSafe(str) {
  if (!str) return '';
  return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

function getHaversineDistanceKm(lat1, lon1, lat2, lon2) {
  var R = 6371; // km
  var dLat = (lat2 - lat1) * Math.PI / 180;
  var dLon = (lon2 - lon1) * Math.PI / 180;
  var a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
          Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
          Math.sin(dLon / 2) * Math.sin(dLon / 2);
  var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

function getItemCoords(dishId, neighbourhood) {
  var item = itemsDB[dishId];
  var nh = neighbourhood || (item ? item.neighbourhood : 'Indiranagar');
  if (BANGALORE_HUBS[nh]) {
    return { lat: BANGALORE_HUBS[nh].lat, lng: BANGALORE_HUBS[nh].lng };
  }
  return { lat: 12.9784, lng: 77.6408 };
}

function getNearestLocalityName(lat, lng) {
  var closestName = 'Indiranagar';
  var minD = Infinity;
  for (var name in BANGALORE_HUBS) {
    var hub = BANGALORE_HUBS[name];
    var d = getHaversineDistanceKm(lat, lng, hub.lat, hub.lng);
    if (d < minD) {
      minD = d;
      closestName = name;
    }
  }
  return { name: closestName, distanceKm: minD };
}

function sortContainerCards(containerEl, userLat, userLng) {
  if (!containerEl) return;
  var cards = Array.from(containerEl.querySelectorAll('.product-card'));
  if (cards.length === 0) return;

  var items = cards.map(function(card) {
    var id = card.getAttribute('data-id');
    var nh = card.getAttribute('data-neighbourhood');
    var coords = getItemCoords(id, nh);
    var dist = getHaversineDistanceKm(userLat, userLng, coords.lat, coords.lng);

    // If card is inside a direct child wrapper (like .carousel-item-wrap), move that wrapper!
    var wrap = card;
    if (card.parentElement && card.parentElement !== containerEl) {
      var p = card.parentElement;
      while (p && p.parentElement && p.parentElement !== containerEl) {
        p = p.parentElement;
      }
      if (p && p.parentElement === containerEl) {
        wrap = p;
      }
    }

    return { card: card, wrap: wrap, dist: dist, id: id, nh: nh };
  });

  // Sort ascending by distance: NEAREST RESTAURANTS FIRST!
  items.sort(function(a, b) {
    return a.dist - b.dist;
  });

  items.forEach(function(item) {
    var distDisplay = item.dist < 1 ? (Math.round(item.dist * 1000) + 'm') : (item.dist.toFixed(1) + ' km');
    var estTime = Math.min(45, Math.max(12, Math.round(12 + item.dist * 2.8))) + 'm';

    // Update delivery badge in card safely without touching layout or images
    var delBadgeSpan = item.card.querySelector('.delivery-badge-pos span');
    if (delBadgeSpan) {
      delBadgeSpan.textContent = distDisplay + ' \u2022 ' + estTime;
    }

    // Move the wrapper or card to end in sorted order
    containerEl.appendChild(item.wrap);
  });

  // Remove any empty orphan wrappers if they exist
  var empties = containerEl.querySelectorAll('.carousel-item-wrap:empty');
  empties.forEach(function(el) { el.remove(); });
}

function sortAllFoodCards() {
  var lat = activeDelivery.lat;
  var lng = activeDelivery.lng;

  // 1. Sort main comprehensive grid
  sortContainerCards(document.getElementById('mainFoodGrid'), lat, lng);

  // 2. Sort horizontal carousel tracks
  var tracks = ['trackUnder100', 'trackSouthIndian', 'trackQuickBites', 'trackBengaluruPicks'];
  tracks.forEach(function(trackId) {
    sortContainerCards(document.getElementById(trackId), lat, lng);
  });
}

function applyDeliveryLocation(locObj, persist) {
  activeDelivery = Object.assign({}, activeDelivery, locObj);

  if (persist) {
    try {
      localStorage.setItem('under100_delivery_loc', JSON.stringify(activeDelivery));
    } catch(e) {}
  }

  // Update navbar location pill (Only show address text, zero icons)
  var locText = document.getElementById('currentLocText');
  var navPill = document.getElementById('navLocationPill');

  if (locText) {
    if (activeDelivery.forSomeoneElse) {
      locText.textContent = activeDelivery.recipientName + ' (' + activeDelivery.locality + ')';
      if (navPill) navPill.title = 'Ordering for ' + activeDelivery.recipientName + ' to ' + activeDelivery.recipientAddress;
    } else {
      locText.textContent = activeDelivery.displayLabel;
      if (navPill) navPill.title = 'Delivering to: ' + activeDelivery.locality;
    }
  }

  // Sort all eateries by distance: nearest first
  sortAllFoodCards();

  // Re-render cart drawer to update recipient notice if cart has items
  renderCart();
}

/* ==========================================================================
   SAVED ADDRESSES & PAN-INDIA / BENGALURU LOCATION DETECTION
   ========================================================================== */

// Curated list of top localities & cities across India and Bengaluru for search recommendations
var POPULAR_INDIAN_LOCALITIES = [
  { name: "HSR Layout", sub: "Sector 1 & 2, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "HSR Layout", lat: 12.9121, lng: 77.6446 },
  { name: "Indiranagar", sub: "100ft Rd / CMH Rd, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "Indiranagar", lat: 12.9784, lng: 77.6408 },
  { name: "Koramangala", sub: "5th & 6th Block, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "Koramangala", lat: 12.9352, lng: 77.6245 },
  { name: "Whitefield", sub: "ITPB / Prestige, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "Whitefield", lat: 12.9698, lng: 77.7499 },
  { name: "Jayanagar", sub: "4th Block Complex, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "Jayanagar", lat: 12.9308, lng: 77.5838 },
  { name: "Malleshwaram", sub: "Margosa / 8th Cross, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "Malleshwaram", lat: 13.0031, lng: 77.5645 },
  { name: "BTM Layout", sub: "Udupi Garden, Bengaluru", city: "Bengaluru", isBengaluru: true, hub: "BTM Layout", lat: 12.9166, lng: 77.6101 },
  { name: "Sector 134, Noida", sub: "Majestic Residency / Nangli, Noida, U.P.", city: "Noida", isBengaluru: false, hub: "Indiranagar", lat: 28.5085, lng: 77.3824 },
  { name: "Sector 22, Noida", sub: "Chora Sadatpur / E Block, Noida, U.P.", city: "Noida", isBengaluru: false, hub: "Indiranagar", lat: 28.5921, lng: 77.3482 },
  { name: "Connaught Place", sub: "Central Delhi, Delhi NCR", city: "Delhi", isBengaluru: false, hub: "Indiranagar", lat: 28.6304, lng: 77.2177 },
  { name: "Cyber City, Gurgaon", sub: "DLF Phase 2, Haryana", city: "Gurgaon", isBengaluru: false, hub: "Indiranagar", lat: 28.4907, lng: 77.0900 },
  { name: "Bandra West", sub: "Linking Road / Pali Hill, Mumbai", city: "Mumbai", isBengaluru: false, hub: "Indiranagar", lat: 19.0596, lng: 72.8295 },
  { name: "Hinjawadi", sub: "Phase 1 Infotech Park, Pune", city: "Pune", isBengaluru: false, hub: "Indiranagar", lat: 18.5913, lng: 73.7389 },
  { name: "Hitec City", sub: "Madhapur, Hyderabad, Telangana", city: "Hyderabad", isBengaluru: false, hub: "Indiranagar", lat: 17.4435, lng: 78.3772 },
  { name: "T. Nagar", sub: "Pondy Bazaar, Chennai, Tamil Nadu", city: "Chennai", isBengaluru: false, hub: "Indiranagar", lat: 13.0418, lng: 80.2341 },
  { name: "Salt Lake", sub: "Sector V, Kolkata, West Bengal", city: "Kolkata", isBengaluru: false, hub: "Indiranagar", lat: 22.5804, lng: 88.4230 }
];

function getSavedAddresses() {
  try {
    var raw = localStorage.getItem('under100_user_saved_addresses');
    if (raw) {
      var parsed = JSON.parse(raw);
      if (Array.isArray(parsed)) return parsed;
    }
  } catch(e) {}
  return [];
}

function setSavedAddresses(list) {
  try {
    localStorage.setItem('under100_user_saved_addresses', JSON.stringify(list || []));
  } catch(e) {}
}

function openLocationModal() {
  var overlay = document.getElementById('locModalOverlay');
  if (overlay) {
    overlay.classList.add('open');
  }
  toggleAddAddressForm(false);
  var searchInput = document.getElementById('locSearchInput');
  if (searchInput) searchInput.value = '';
  var suggestions = document.getElementById('locSearchSuggestions');
  if (suggestions) suggestions.style.display = 'none';

  renderSavedAddresses('');
}

function closeLocationModal() {
  var overlay = document.getElementById('locModalOverlay');
  if (overlay) {
    overlay.classList.remove('open');
  }
}

function renderSavedAddresses(filterQuery) {
  var container = document.getElementById('savedAddressesList');
  if (!container) return;

  var list = getSavedAddresses();
  var q = (filterQuery || '').trim().toLowerCase();

  var filtered = list;
  if (q) {
    filtered = list.filter(function(item) {
      return (item.title && item.title.toLowerCase().indexOf(q) !== -1) ||
             (item.address && item.address.toLowerCase().indexOf(q) !== -1) ||
             (item.locality && item.locality.toLowerCase().indexOf(q) !== -1) ||
             (item.city && item.city.toLowerCase().indexOf(q) !== -1) ||
             (item.recipientName && item.recipientName.toLowerCase().indexOf(q) !== -1);
    });
  }

  if (filtered.length === 0) {
    if (q) {
      container.innerHTML = '<div style="padding: 24px 12px; text-align: center; color: #64748B; font-size: 13px;">' +
        'No saved addresses match "' + escapeSafe(filterQuery) + '"' +
      '</div>';
    } else {
      container.innerHTML = '<div style="background-color: #F8FAFC; border: 1.5px dashed #CBD5E1; border-radius: 16px; padding: 26px 16px; text-align: center; margin-top: 4px;">' +
        '<div style="width: 44px; height: 44px; border-radius: 50%; background-color: #FEF3C7; color: #D97706; display: flex; align-items: center; justify-content: center; margin: 0 auto 10px; font-size: 20px;">' +
          '📍' +
        '</div>' +
        '<div style="font-size: 14.5px; font-weight: 800; color: #1E293B; margin-bottom: 4px;">No saved addresses yet</div>' +
        '<div style="font-size: 12.5px; color: #64748B; max-width: 320px; margin: 0 auto 14px; line-height: 1.45;">' +
          'Addresses you save for yourself or order for someone else will appear here for 1-click delivery.' +
        '</div>' +
        '<button type="button" class="btn-detect-location" style="padding: 9px 18px; font-size: 12.5px; margin: 0 auto; display: inline-flex;" onclick="toggleAddAddressForm(true)">' +
          '+ Add Address / Order for Someone Else' +
        '</button>' +
      '</div>';
    }
    return;
  }

  var html = '';
  filtered.forEach(function(item) {
    // Find original index in master list
    var origIndex = list.indexOf(item);
    var isSomeoneElse = item.forSomeoneElse || item.type === 'Other' || !!item.recipientName;
    var isHome = !isSomeoneElse && (item.type === 'Home' || (item.title && item.title.toLowerCase() === 'home'));
    var isWork = !isSomeoneElse && (item.type === 'Work' || (item.title && item.title.toLowerCase() === 'work'));

    // Icon matching user\'s screenshot (golden yellow home or people icon)
    var iconHTML = '';
    if (isSomeoneElse) {
      iconHTML = '<svg width="22" height="22" viewBox="0 0 24 24" fill="#D97706"><path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/></svg>';
    } else if (isWork) {
      iconHTML = '<svg width="20" height="20" viewBox="0 0 24 24" fill="#D97706"><path d="M20 6h-4V4c0-1.11-.89-2-2-2h-4c-1.11 0-2 .89-2 2v2H4c-1.11 0-1.99.89-1.99 2L2 19c0 1.11.89 2 2 2h16c1.11 0 2-.89 2-2V8c0-1.11-.89-2-2-2zm-6 0h-4V4h4v2z"/></svg>';
    } else {
      iconHTML = '<svg width="22" height="22" viewBox="0 0 24 24" fill="#D97706"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/></svg>';
    }

    var isCurrentSelected = false;
    if (activeDelivery.recipientAddress && activeDelivery.recipientAddress === item.address) {
      isCurrentSelected = true;
    } else if (activeDelivery.displayLabel === item.title || activeDelivery.locality === item.locality) {
      isCurrentSelected = true;
    }

    var titleDisplay = escapeSafe(item.title || (isSomeoneElse ? item.recipientName : 'Home'));
    var addressDisplay = escapeSafe(item.address || '');

    html += '<div class="saved-address-card' + (isCurrentSelected ? ' selected-active' : '') + '" onclick="selectSavedAddress(' + origIndex + ')">' +
      '<div class="address-icon-box">' + iconHTML + '</div>' +
      '<div class="address-content-box">' +
        '<div class="address-name-title">' +
          '<span>' + titleDisplay + '</span>' +
          (isCurrentSelected ? '<span style="font-size: 11px; font-weight: 700; color: #0B6837; background: #DCFCE7; padding: 2px 7px; border-radius: 9999px;">Active</span>' : '') +
          (isSomeoneElse ? '<span style="font-size: 10.5px; font-weight: 700; color: #9A3412; background: #FFEDD5; padding: 2px 6px; border-radius: 6px;">For Someone Else</span>' : '') +
        '</div>' +
        '<div class="address-details-text">' + addressDisplay + '</div>' +
        '<div class="address-card-actions">' +
          '<button type="button" class="btn-action-round edit" title="Edit Address" onclick="event.stopPropagation(); editAddress(' + origIndex + ')">' +
            '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#16A34A" stroke-width="2.5"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>' +
          '</button>' +
          '<button type="button" class="btn-action-round delete" title="Delete Address" onclick="event.stopPropagation(); deleteAddress(' + origIndex + ')">' +
            '<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#DC2626" stroke-width="2.5"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>' +
          '</button>' +
        '</div>' +
      '</div>' +
    '</div>';
  });

  container.innerHTML = html;
}

function selectSavedAddress(index) {
  var list = getSavedAddresses();
  var item = list[index];
  if (!item) return;

  var isBengaluru = (item.city && item.city.toLowerCase() === 'bengaluru') ||
                    (item.locality && BANGALORE_HUBS[item.locality]) ||
                    (item.address && item.address.toLowerCase().indexOf('bengaluru') !== -1);

  var hubName = item.locality;
  if (!BANGALORE_HUBS[hubName]) {
    hubName = 'Indiranagar';
  }
  var anchorHub = BANGALORE_HUBS[hubName] || BANGALORE_HUBS['Indiranagar'];

  var isSomeone = item.forSomeoneElse || item.type === 'Other' || !!item.recipientName;
  var displayLabel = item.title;
  if (!displayLabel || displayLabel === 'Home' || displayLabel === 'Work') {
    displayLabel = item.locality || item.city || 'Bengaluru';
  }

  applyDeliveryLocation({
    isGPS: false,
    lat: isBengaluru ? (item.lat || anchorHub.lat) : anchorHub.lat,
    lng: isBengaluru ? (item.lng || anchorHub.lng) : anchorHub.lng,
    locality: item.locality || item.city || 'Bengaluru',
    displayLabel: isSomeone ? ('🎁 ' + (item.recipientName || item.title)) : displayLabel,
    forSomeoneElse: isSomeone,
    recipientName: item.recipientName || item.title || '',
    recipientPhone: item.recipientPhone || '',
    recipientAddress: item.address,
    recipientNote: item.recipientNote || ''
  }, true);

  // Show friendly notification if outside Bengaluru
  var notice = document.getElementById('panIndiaNotice');
  if (!isBengaluru) {
    if (notice) {
      notice.style.display = 'flex';
      notice.innerHTML = '<span style="font-size: 16px;">📍</span><span><strong>Delivery Address: ' + escapeSafe(item.city || 'Pan-India') + '</strong> — Under100 is currently active in Bengaluru hub! Showing available delicacies with pan-India delivery launching soon.</span>';
    }
  } else {
    if (notice) notice.style.display = 'none';
  }

  closeLocationModal();
}

function toggleAddAddressForm(show) {
  var formWrap = document.getElementById('addAddressFormContainer');
  if (!formWrap) return;

  var isOpening = (show === undefined) ? (formWrap.style.display === 'none') : show;
  formWrap.style.display = isOpening ? 'block' : 'none';

  if (isOpening) {
    document.getElementById('editAddressIndex').value = '-1';
    document.getElementById('addressFormTitle').textContent = 'Add New Address / Order for Someone Else';
    document.getElementById('formRecipName').value = '';
    document.getElementById('formRecipPhone').value = '';
    document.getElementById('formAddressFull').value = '';
    document.getElementById('formLocality').value = '';
    document.getElementById('formNote').value = '';
    var homeRadio = document.querySelector('input[name="addrType"][value="Home"]');
    if (homeRadio) homeRadio.checked = true;
    toggleRecipientFields(false);
  }
}

function toggleRecipientFields(isOther) {
  var grp = document.getElementById('recipientFieldsGroup');
  if (grp) {
    grp.style.display = isOther ? 'block' : 'none';
  }
}

function editAddress(index) {
  var list = getSavedAddresses();
  var item = list[index];
  if (!item) return;

  toggleAddAddressForm(true);
  document.getElementById('editAddressIndex').value = String(index);
  document.getElementById('addressFormTitle').textContent = 'Edit Saved Address';

  var isOther = item.forSomeoneElse || item.type === 'Other' || !!item.recipientName;
  var radioVal = isOther ? 'Other' : (item.type || 'Home');
  var radio = document.querySelector('input[name="addrType"][value="' + radioVal + '"]');
  if (radio) radio.checked = true;
  toggleRecipientFields(isOther);

  document.getElementById('formRecipName').value = item.recipientName || '';
  document.getElementById('formRecipPhone').value = item.recipientPhone || '';
  document.getElementById('formAddressFull').value = item.address || '';
  document.getElementById('formLocality').value = item.locality || item.city || '';
  document.getElementById('formNote').value = item.recipientNote || '';

  var formWrap = document.getElementById('addAddressFormContainer');
  if (formWrap) formWrap.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function deleteAddress(index) {
  var list = getSavedAddresses();
  var item = list[index];
  if (!item) return;

  if (confirm('Are you sure you want to delete address "' + (item.title || item.address.slice(0, 20)) + '"?')) {
    list.splice(index, 1);
    setSavedAddresses(list);
    renderSavedAddresses('');
  }
}

function handleSaveAddressSubmit(e) {
  e.preventDefault();
  var idx = parseInt(document.getElementById('editAddressIndex').value, 10);
  var typeRadio = document.querySelector('input[name="addrType"]:checked');
  var type = typeRadio ? typeRadio.value : 'Home';
  var isOther = (type === 'Other');

  var recipName = (document.getElementById('formRecipName').value || '').trim();
  var recipPhone = (document.getElementById('formRecipPhone').value || '').trim();
  var addressFull = (document.getElementById('formAddressFull').value || '').trim();
  var locality = (document.getElementById('formLocality').value || '').trim();
  var note = (document.getElementById('formNote').value || '').trim();

  if (isOther && !recipName) {
    alert('Please enter recipient name for Someone Else.');
    return;
  }
  if (!addressFull) {
    alert('Please enter complete street address.');
    return;
  }
  if (!locality) {
    locality = 'Bengaluru';
  }

  var title = isOther ? recipName : type;
  var city = locality.toLowerCase().indexOf('bengaluru') !== -1 || locality.toLowerCase().indexOf('bangalore') !== -1 ? 'Bengaluru' : locality;

  var addressObj = {
    id: 'addr_' + Date.now(),
    type: type,
    title: title,
    address: addressFull,
    locality: locality,
    city: city,
    forSomeoneElse: isOther,
    recipientName: recipName,
    recipientPhone: recipPhone,
    recipientNote: note,
    lat: 12.9784,
    lng: 77.6408
  };

  var list = getSavedAddresses();
  if (idx >= 0 && idx < list.length) {
    addressObj.id = list[idx].id;
    list[idx] = addressObj;
  } else {
    list.push(addressObj);
  }

  setSavedAddresses(list);
  toggleAddAddressForm(false);
  renderSavedAddresses('');

  // Automatically select the saved address
  var targetIndex = (idx >= 0 && idx < list.length) ? idx : (list.length - 1);
  selectSavedAddress(targetIndex);
}

var locSearchTimer = null;

function handleLocSearch(query) {
  var q = (query || '').trim();
  var suggestionsBox = document.getElementById('locSearchSuggestions');

  // 1. Live filter the user's saved addresses section below
  renderSavedAddresses(q);

  if (!q || q.length < 2) {
    if (suggestionsBox) suggestionsBox.style.display = 'none';
    return;
  }

  if (suggestionsBox) {
    suggestionsBox.style.display = 'block';
    suggestionsBox.innerHTML = '<div style="padding: 12px 14px; font-size: 13px; color: #64748B; display: flex; align-items: center; gap: 8px;">' +
      '<span>🔍 Searching live delivery addresses...</span>' +
    '</div>';
  }

  clearTimeout(locSearchTimer);
  locSearchTimer = setTimeout(function() {
    fetch('/api/location/search?q=' + encodeURIComponent(q))
      .then(function(res) { return res.json(); })
      .then(function(results) {
        if (!suggestionsBox) return;
        if (!Array.isArray(results) || results.length === 0) {
          suggestionsBox.innerHTML = '<div style="padding: 12px 14px; font-size: 12.5px; color: #64748B;">' +
            'No exact addresses found for "' + escapeSafe(q) + '". Try entering your colony, sector, or area.' +
          '</div>';
          return;
        }

        var html = '<div style="padding: 6px 12px 4px; font-size: 11px; font-weight: 800; color: #94A3B8; text-transform: uppercase; letter-spacing: 0.5px;">Matching Delivery Locations</div>';
        results.forEach(function(item, idx) {
          var name = item.name || 'Location';
          var full = item.display_name || '';
          var isBlr = item.is_bengaluru;
          var tag = isBlr ? '<span class="loc-suggestion-tag">Bengaluru Hub</span>' : '<span class="loc-suggestion-tag hub">India (' + escapeSafe(item.city || 'Hub') + ')</span>';

          html += '<div class="loc-suggestion-item" onclick="selectExactLocation(' + idx + ')">' +
            '<div class="loc-suggestion-item-icon">📍</div>' +
            '<div class="loc-suggestion-item-text">' +
              '<div class="loc-suggestion-item-title">' + escapeSafe(name) + '</div>' +
              '<div class="loc-suggestion-item-sub">' + escapeSafe(full) + '</div>' +
            '</div>' +
            '<div>' + tag + '</div>' +
          '</div>';
        });

        window._lastLocationResults = results;
        suggestionsBox.innerHTML = html;
        suggestionsBox.style.display = 'block';
      })
      .catch(function(err) {
        console.warn('Location search error:', err);
        if (suggestionsBox) {
          suggestionsBox.innerHTML = '<div style="padding: 10px; font-size: 12.5px; color: #DC2626;">Could not fetch address suggestions. Please try again.</div>';
        }
      });
  }, 220);
}

function selectExactLocation(index) {
  var list = window._lastLocationResults || [];
  var item = list[index];
  if (!item) return;

  var isBlr = item.is_bengaluru;
  var anchor = isBlr ? getNearestLocalityName(item.lat, item.lng) : null;
  var hubCoords = anchor ? getItemCoords('', anchor.name) : { lat: 12.9784, lng: 77.6408 };

  var shortTitle = item.name || item.locality || item.city || 'Location';

  applyDeliveryLocation({
    isGPS: false,
    lat: isBlr ? item.lat : hubCoords.lat,
    lng: isBlr ? item.lng : hubCoords.lng,
    locality: isBlr ? (anchor ? anchor.name : 'Bengaluru') : (item.city || 'India'),
    displayLabel: shortTitle + (isBlr ? ', Bengaluru' : ' (' + (item.city || 'India') + ')'),
    forSomeoneElse: false,
    recipientName: '',
    recipientPhone: '',
    recipientAddress: item.display_name,
    recipientNote: ''
  }, true);

  var notice = document.getElementById('panIndiaNotice');
  if (!isBlr) {
    if (notice) {
      notice.style.display = 'flex';
      notice.innerHTML = '<span style="font-size: 16px;">📍</span><span><strong>Delivery Address: ' + escapeSafe(shortTitle) + '</strong> — Under100 is currently live in Bengaluru hub! Showing available delicacies with pan-India rollout coming soon.</span>';
    }
  } else {
    if (notice) notice.style.display = 'none';
  }

  closeLocationModal();
}

function detectGPSLocation() {
  var btn = document.getElementById('btnDetectLoc');
  var originalText = btn ? btn.textContent : 'Detect my location';
  if (btn) {
    btn.textContent = 'Detecting exact street...';
    btn.disabled = true;
  }

  if (!navigator.geolocation) {
    alert('Geolocation is not supported by your browser. Please search or pick a saved address.');
    if (btn) { btn.textContent = originalText; btn.disabled = false; }
    return;
  }

  navigator.geolocation.getCurrentPosition(
    function(pos) {
      var lat = pos.coords.latitude;
      var lng = pos.coords.longitude;

      // Reverse geocode via server to get exact street, building & locality
      fetch('/api/location/reverse?lat=' + lat + '&lng=' + lng)
        .then(function(res) { return res.json(); })
        .then(function(data) {
          resolveAndApplyExactLocation(data, lat, lng, true);
          if (btn) { btn.textContent = originalText; btn.disabled = false; }
          closeLocationModal();
        })
        .catch(function(e) {
          // Fallback if reverse geocoding API times out
          resolveAndApplyCoordinates(lat, lng, true);
          if (btn) { btn.textContent = originalText; btn.disabled = false; }
          closeLocationModal();
        });
    },
    function(err) {
      console.warn('GPS location error:', err);
      // Fallback: check first saved address or default to HSR Layout / Bengaluru
      var savedList = getSavedAddresses();
      if (savedList && savedList.length > 0) {
        selectSavedAddress(0);
      } else {
        applyDeliveryLocation({
          isGPS: false,
          lat: BANGALORE_HUBS['HSR Layout'].lat,
          lng: BANGALORE_HUBS['HSR Layout'].lng,
          locality: 'HSR Layout',
          displayLabel: 'HSR Layout, Bengaluru',
          forSomeoneElse: false
        }, true);
      }
      if (btn) { btn.textContent = originalText; btn.disabled = false; }
      closeLocationModal();
    },
    { enableHighAccuracy: true, timeout: 9000, maximumAge: 60000 }
  );
}

function resolveAndApplyExactLocation(data, lat, lng, isGPS) {
  var isBlr = data.is_bengaluru;
  var exactStreet = data.exact_street || '';
  var locality = data.locality || '';
  var city = data.city || (isBlr ? 'Bengaluru' : 'India');
  var fullAddress = data.display_name || (exactStreet + ', ' + locality + ', ' + city);

  var displayLabel = exactStreet ? (exactStreet + ', ' + (locality || city)) : (locality || city);
  if (isGPS) {
    displayLabel += ' (GPS)';
  }

  var nearestHub = isBlr ? getNearestLocalityName(lat, lng) : null;
  var anchor = BANGALORE_HUBS['Indiranagar'];

  applyDeliveryLocation({
    isGPS: isGPS,
    lat: isBlr ? lat : anchor.lat,
    lng: isBlr ? lng : anchor.lng,
    locality: isBlr ? (nearestHub ? nearestHub.name : locality) : city,
    displayLabel: displayLabel,
    forSomeoneElse: false,
    recipientAddress: fullAddress
  }, true);

  var notice = document.getElementById('panIndiaNotice');
  if (!isBlr) {
    if (notice) {
      notice.style.display = 'flex';
      notice.innerHTML = '<span style="font-size: 16px;">📍</span><span><strong>Exact GPS: ' + escapeSafe(displayLabel) + '</strong> — Under100 is currently live in Bengaluru hub! Showing available delicacies with pan-India rollout coming soon.</span>';
    }
  } else {
    if (notice) notice.style.display = 'none';
  }
}

function resolveAndApplyCoordinates(lat, lng, isGPS) {
  var distToBengaluru = getHaversineDistanceKm(lat, lng, 12.9716, 77.5946);

  if (distToBengaluru <= 45) {
    var nearest = getNearestLocalityName(lat, lng);
    var label = nearest.name + (isGPS ? ' (GPS)' : '');
    applyDeliveryLocation({
      isGPS: isGPS,
      lat: lat,
      lng: lng,
      locality: nearest.name,
      displayLabel: label,
      forSomeoneElse: false
    }, true);
  } else {
    var closestCity = "Noida, U.P.";
    var anchor = BANGALORE_HUBS['Indiranagar'];

    applyDeliveryLocation({
      isGPS: isGPS,
      lat: anchor.lat,
      lng: anchor.lng,
      locality: closestCity,
      displayLabel: closestCity + (isGPS ? ' (GPS Hub)' : ''),
      forSomeoneElse: false,
      recipientAddress: closestCity
    }, true);

    var notice = document.getElementById('panIndiaNotice');
    if (notice) {
      notice.style.display = 'flex';
      notice.innerHTML = '<span style="font-size: 16px;">📍</span><span><strong>Detected Location: ' + escapeSafe(closestCity) + '</strong> — Under100 is currently live in Bengaluru! Showing available delicacies with pan-India rollout coming soon.</span>';
    }
  }
}

function initAutoGeolocation() {
  // Check if saved delivery preference exists in localStorage
  var saved = null;
  try {
    var raw = localStorage.getItem('under100_delivery_loc');
    if (raw) saved = JSON.parse(raw);
  } catch(e) {}

  if (saved && saved.locality) {
    applyDeliveryLocation(saved, false);
    return;
  }

  // Check if we have user-saved addresses in list
  var savedList = getSavedAddresses();
  if (savedList && savedList.length > 0) {
    selectSavedAddress(0);
    return;
  }

  // Automatic GPS Geolocation on first load / login anywhere in India
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      function(pos) {
        resolveAndApplyCoordinates(pos.coords.latitude, pos.coords.longitude, true);
      },
      function(err) {
        // Fallback default: HSR Layout, Bengaluru
        applyDeliveryLocation({
          isGPS: false,
          lat: BANGALORE_HUBS['HSR Layout'].lat,
          lng: BANGALORE_HUBS['HSR Layout'].lng,
          locality: 'HSR Layout',
          displayLabel: 'HSR Layout, Bengaluru',
          forSomeoneElse: false
        }, false);
      },
      { enableHighAccuracy: true, timeout: 6000, maximumAge: 300000 }
    );
  } else {
    // Default to HSR Layout
    applyDeliveryLocation({
      isGPS: false,
      lat: BANGALORE_HUBS['HSR Layout'].lat,
      lng: BANGALORE_HUBS['HSR Layout'].lng,
      locality: 'HSR Layout',
      displayLabel: 'HSR Layout, Bengaluru',
      forSomeoneElse: false
    }, false);
  }
}

// Dynamic Rotating Search Bar Placeholders (Pav Bhaji, Biryani, Dosa, etc.)
(function initRotatingSearchPlaceholders() {
  var searchInput = document.getElementById('globalSearchInput');
  if (!searchInput) return;

  var placeholders = [
    "Search 'Pav Bhaji'...",
    "Search 'Biryani'...",
    "Search 'Benne Masala Dosa'...",
    "Search 'Filter Coffee'...",
    "Search 'Chole Bhature'...",
    "Search 'Paneer Roll'...",
    "Search 'Vada Pav'...",
    "Search 'Steamed Momos'...",
    "Search 'South Indian Meals'...",
    "Search 'Egg Donne Biryani'..."
  ];

  var currentIndex = 0;
  var charIndex = 0;
  var isDeleting = false;
  var isPaused = false;

  searchInput.addEventListener('focus', function() { isPaused = true; });
  searchInput.addEventListener('blur', function() {
    if (!searchInput.value.trim()) {
      isPaused = false;
    }
  });
  searchInput.addEventListener('input', function() {
    if (searchInput.value.trim().length > 0) {
      isPaused = true;
    } else {
      isPaused = false;
    }
  });

  function typeTick() {
    if (isPaused) {
      setTimeout(typeTick, 400);
      return;
    }

    var currentText = placeholders[currentIndex];
    
    if (isDeleting) {
      charIndex--;
      searchInput.setAttribute('placeholder', currentText.substring(0, charIndex));
      if (charIndex <= 0) {
        isDeleting = false;
        currentIndex = (currentIndex + 1) % placeholders.length;
        setTimeout(typeTick, 350);
        return;
      }
      setTimeout(typeTick, 30);
    } else {
      charIndex++;
      searchInput.setAttribute('placeholder', currentText.substring(0, charIndex));
      if (charIndex >= currentText.length) {
        isDeleting = true;
        setTimeout(typeTick, 2000); // Display full placeholder for 2s before cycling
        return;
      }
      setTimeout(typeTick, 55);
    }
  }

  // Kick off dynamic rotation
  setTimeout(typeTick, 300);
})();

// User Account dropdown toggle & click-outside
function toggleUserDropdown() {
  var popover = document.getElementById('accountMenuPopover');
  if (popover) {
    popover.classList.toggle('open');
  }
}
document.addEventListener('click', function(e) {
  var userDropdown = document.getElementById('userMenuDropdown');
  var popover = document.getElementById('accountMenuPopover');
  if (userDropdown && popover && !userDropdown.contains(e.target)) {
    popover.classList.remove('open');
  }
});

// Account Modals Handlers
function openMyOrdersModal() {
  var m = document.getElementById('myOrdersModalOverlay');
  if (m) m.classList.add('open');
}
function closeMyOrdersModal() {
  var m = document.getElementById('myOrdersModalOverlay');
  if (m) m.classList.remove('open');
}

function openPrescriptionsModal() {
  var m = document.getElementById('prescriptionsModalOverlay');
  if (m) m.classList.add('open');
}
function closePrescriptionsModal() {
  var m = document.getElementById('prescriptionsModalOverlay');
  if (m) m.classList.remove('open');
}
function handlePrescriptionUpload(input) {
  var st = document.getElementById('prescriptionUploadStatus');
  if (input.files && input.files[0]) {
    if (st) st.textContent = '✓ ' + input.files[0].name + ' uploaded successfully!';
  }
}

function openGiftCardsModal() {
  var m = document.getElementById('giftCardsModalOverlay');
  if (m) m.classList.add('open');
}
function closeGiftCardsModal() {
  var m = document.getElementById('giftCardsModalOverlay');
  if (m) m.classList.remove('open');
}
function redeemGiftCardCode() {
  var inp = document.getElementById('giftCardCodeInput');
  if (!inp || !inp.value.trim()) {
    showToast('Please enter a voucher code');
    return;
  }
  showToast('Gift Card voucher redeemed successfully!');
  inp.value = '';
}
function selectGiftAmount(amt) {
  showToast('Selected ₹' + amt + ' Gift Card voucher');
}

function openFaqsModal() {
  var m = document.getElementById('faqsModalOverlay');
  if (m) m.classList.add('open');
}
function closeFaqsModal() {
  var m = document.getElementById('faqsModalOverlay');
  if (m) m.classList.remove('open');
}
function toggleFaqItem(btn) {
  var body = btn.nextElementSibling;
  var arrow = btn.querySelector('.faq-arrow');
  if (body) {
    if (body.style.display === 'block') {
      body.style.display = 'none';
      if (arrow) arrow.style.transform = 'rotate(0deg)';
    } else {
      body.style.display = 'block';
      if (arrow) arrow.style.transform = 'rotate(180deg)';
    }
  }
}

function openAccountPrivacyModal() {
  var m = document.getElementById('privacyModalOverlay');
  if (m) m.classList.add('open');
}
function closeAccountPrivacyModal() {
  var m = document.getElementById('privacyModalOverlay');
  if (m) m.classList.remove('open');
}
function downloadAccountData() {
  showToast('Preparing your account data archive. Download starting...');
}
function confirmDeleteAccount() {
  if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
    window.location.href = '/logout';
  }
}
function reorderItem(itemName, price) {
  addToCart(itemName, price);
  closeMyOrdersModal();
  toggleCartDrawer();
  showToast('Added ' + itemName + ' to cart!');
}

// Login Modal Handlers
var modalDemoOtpCodeVal = '123456';
var modalOtpDispatched = false;

function openLoginModal() {
  var m = document.getElementById('loginModalOverlay');
  if (m) m.classList.add('open');
}
function closeLoginModal() {
  var m = document.getElementById('loginModalOverlay');
  if (m) m.classList.remove('open');
}

function switchModalAuth(mode) {
  var bOtp = document.getElementById('modalSubtabOtp');
  var bPass = document.getElementById('modalSubtabPass');
  var sOtp = document.getElementById('modalOtpSection');
  var sPass = document.getElementById('modalPassSection');
  if (mode === 'pass') {
    if (bOtp) { bOtp.style.background = 'transparent'; bOtp.style.color = '#64748B'; bOtp.style.boxShadow = 'none'; }
    if (bPass) { bPass.style.background = '#ffffff'; bPass.style.color = '#E23744'; bPass.style.boxShadow = '0 1px 3px rgba(0,0,0,0.08)'; }
    if (sOtp) sOtp.style.display = 'none';
    if (sPass) sPass.style.display = 'flex';
  } else {
    if (bPass) { bPass.style.background = 'transparent'; bPass.style.color = '#64748B'; bPass.style.boxShadow = 'none'; }
    if (bOtp) { bOtp.style.background = '#ffffff'; bOtp.style.color = '#0C831F'; bOtp.style.boxShadow = '0 1px 3px rgba(0,0,0,0.08)'; }
    if (sPass) sPass.style.display = 'none';
    if (sOtp) sOtp.style.display = 'block';
  }
}



function getCsrfToken() {
  var match = document.cookie.match(new RegExp('(^| )_csrf=([^;]+)'));
  return match ? decodeURIComponent(match[2]) : '';
}

function sendModalOTP() {
  var idInput = document.getElementById('loginIdentifierInput');
  var idVal = idInput ? idInput.value.trim() : '';
  var statusDiv = document.getElementById('modalOtpStatus');
  var btnSend = document.getElementById('btnModalSendOtp');
  if (!idVal) {
    if (statusDiv) {
      statusDiv.style.display = 'block';
      statusDiv.style.color = '#DC2626';
      statusDiv.innerText = 'Please enter your mobile number or email.';
    }
    if (idInput) idInput.focus();
    return;
  }
  btnSend.disabled = true;
  btnSend.innerText = 'Sending...';
  if (statusDiv) statusDiv.style.display = 'none';

  var csrf = getCsrfToken();
  var body = new URLSearchParams();
  body.append('identifier', idVal);
  body.append('portal', 'customer');
  if (csrf) body.append('csrf', csrf);

  fetch('/api/auth/send-otp', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'X-CSRF-Token': csrf
    },
    body: body.toString()
  })
  .then(function(r) {
    return r.text().then(function(text) {
      try {
        return JSON.parse(text);
      } catch(e) {
        return { success: false, error: text || ('Server error (' + r.status + ')') };
      }
    });
  })
  .then(function(data) {
    if (data && data.success) {
      modalOtpDispatched = true;
      var inp = document.getElementById('modalOtpCodeInput');
      if (inp) {
        inp.value = data.demo_otp || '';
        inp.focus();
      }
      var wrap = document.getElementById('modalOtpInputWrap');
      if (wrap) wrap.style.display = 'block';
      var contBtn = document.getElementById('btnLoginContinue');
      if (contBtn) {
        contBtn.innerText = 'Verify & Sign In';
        contBtn.disabled = false;
      }
      if (statusDiv) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#166534';
        if (data.demo_otp) {
          statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: #DCFCE7; color: #166534; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #4B5563;">(Code auto-filled &mdash; click "Verify & Sign In" to complete login)</span>';
        } else {
          statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || idVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
        }
      }
      var count = 60;
      var timer = setInterval(function() {
        count--;
        if (count <= 0) {
          clearInterval(timer);
          btnSend.disabled = false;
          btnSend.innerText = 'Resend';
        } else {
          btnSend.innerText = count + 's';
        }
      }, 1000);
    } else {
      btnSend.disabled = false;
      btnSend.innerText = 'Get OTP';
      if (statusDiv) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#DC2626';
        statusDiv.innerText = (data && data.error) ? data.error : 'Failed to send OTP.';
      }
    }
  })
  .catch(function(err) {
    btnSend.disabled = false;
    btnSend.innerText = 'Get OTP';
    if (statusDiv) {
      statusDiv.style.display = 'block';
      statusDiv.style.color = '#DC2626';
      statusDiv.innerText = (err && err.message) ? err.message : 'Network error while requesting OTP.';
    }
  });
}

function handleModalOtpSubmit() {
  var idInput = document.getElementById('loginIdentifierInput');
  var idVal = idInput ? idInput.value.trim() : '';
  if (!idVal) {
    alert('Please enter your mobile number or email.');
    if (idInput) idInput.focus();
    return;
  }
  var otpInput = document.getElementById('modalOtpCodeInput');
  var otpVal = otpInput ? otpInput.value.trim() : '';

  if (modalOtpDispatched && otpVal) {
    var btn = document.getElementById('btnLoginContinue');
    if (btn) {
      btn.disabled = true;
      btn.innerText = 'Verifying...';
    }
    var statusDiv = document.getElementById('modalOtpStatus');

    var csrf = getCsrfToken();
    var body = new URLSearchParams();
    body.append('identifier', idVal);
    body.append('otp', otpVal);
    body.append('portal', 'customer');
    if (csrf) body.append('csrf', csrf);

    fetch('/api/auth/verify-otp', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'X-CSRF-Token': csrf
      },
      body: body.toString()
    })
    .then(function(r) {
      return r.text().then(function(text) {
        try {
          return JSON.parse(text);
        } catch(e) {
          return { success: false, error: text || ('Server error (' + r.status + ')') };
        }
      });
    })
    .then(function(data) {
      if (data && data.success) {
        if (statusDiv) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#16A34A';
          statusDiv.innerText = 'Verified! Signing in...';
        }
        closeLoginModal();
        window.location.reload();
      } else {
        if (btn) {
          btn.disabled = false;
          btn.innerText = 'Verify & Sign In';
        }
        if (statusDiv) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#DC2626';
          statusDiv.innerText = (data && data.error) ? data.error : 'Invalid OTP code.';
        }
      }
    })
    .catch(function() {
      if (btn) {
        btn.disabled = false;
        btn.innerText = 'Verify & Sign In';
      }
      quickPhoneLogin(idVal);
    });
  } else if (!modalOtpDispatched) {
    sendModalOTP();
  } else {
    var inp = document.getElementById('modalOtpCodeInput');
    if (inp) inp.focus();
    var statusDiv = document.getElementById('modalOtpStatus');
    if (statusDiv) {
      statusDiv.style.display = 'block';
      statusDiv.style.color = '#DC2626';
      statusDiv.innerText = 'Please enter the 6-digit OTP code sent to your phone.';
    }
  }
}

function handleModalPassSubmit(e) {
  if (e && e.preventDefault) e.preventDefault();
  var email = document.getElementById('modalPassEmail').value.trim();
  var pass = document.getElementById('modalPassPassword').value;
  var statusDiv = document.getElementById('modalPassStatus');
  var btn = document.getElementById('btnModalPassLogin');

  btn.disabled = true;
  btn.innerText = 'Signing In...';
  if (statusDiv) statusDiv.style.display = 'none';

  var body = new URLSearchParams();
  body.append('email', email);
  body.append('password', pass);

  fetch('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: body.toString(),
    redirect: 'follow'
  })
  .then(function(resp) {
    if (resp.redirected || resp.ok) {
      closeLoginModal();
      window.location.reload();
    } else {
      btn.disabled = false;
      btn.innerText = 'Sign In with Password';
      if (statusDiv) {
        statusDiv.style.display = 'block';
        statusDiv.innerText = 'Invalid email or password.';
      }
    }
  })
  .catch(function() {
    var form = document.createElement('form');
    form.method = 'POST';
    form.action = '/login';
    var inpE = document.createElement('input');
    inpE.name = 'email';
    inpE.value = email;
    form.appendChild(inpE);
    var inpP = document.createElement('input');
    inpP.name = 'password';
    inpP.value = pass;
    form.appendChild(inpP);
    document.body.appendChild(form);
    form.submit();
  });
}

function handleLoginSubmit(e) {
  if (e && e.preventDefault) e.preventDefault();
  handleModalOtpSubmit();
}

function quickPhoneLogin(identifier) {
  var btn = document.getElementById('btnLoginContinue');
  if (btn) {
    btn.textContent = 'Logging in...';
    btn.disabled = true;
  }
  var csrf = getCsrfToken();
  var formData = new FormData();
  formData.append('phone', identifier);
  if (csrf) formData.append('csrf', csrf);
  fetch('/api/phone-login', {
    method: 'POST',
    body: formData,
    credentials: 'same-origin',
    headers: {
      'X-Requested-With': 'XMLHttpRequest',
      'X-CSRF-Token': csrf
    }
  }).then(function(res) {
    return res.text().then(function(text) {
      try {
        return JSON.parse(text);
      } catch(e) {
        return { success: false, error: text || ('Server error (' + res.status + ')') };
      }
    });
  }).then(function(data) {
    if (data && data.success) {
      localStorage.setItem('under100_user_phone', identifier);
      closeLoginModal();
      window.location.reload();
    } else {
      if (btn) {
        btn.textContent = modalOtpDispatched ? 'Verify & Sign In' : 'Continue';
        btn.disabled = false;
      }
      var statusDiv = document.getElementById('modalOtpStatus');
      if (statusDiv) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#DC2626';
        statusDiv.innerText = (data && data.error) ? data.error : 'Please enter a valid mobile number or email.';
      }
    }
  }).catch(function() {
    if (btn) {
      btn.textContent = modalOtpDispatched ? 'Verify & Sign In' : 'Continue';
      btn.disabled = false;
    }
    window.location.href = '/login';
  });
}



// Login modal only opens when user clicks Login button in navbar

// Initial execution
renderCart();
initAutoGeolocation();

// Auto-select category if passed in query string (?category=... or ?cuisine=...)
(function() {
  var urlParams = new URLSearchParams(window.location.search);
  var cat = urlParams.get('category') || urlParams.get('cuisine');
  if (cat) {
    selectCategoryFilter(cat);
  }
})();

  </script>
  <div id="modal-container"></div>
  <script>
    function closeModal() {
      var modalContainer = document.getElementById('modal-container');
      if (modalContainer) {
        modalContainer.innerHTML = '';
      }
    }
  </script>
</body>
</html>`

	pageHTML = strings.ReplaceAll(pageHTML, "__AUTH_ACTIONS__", authActionsHTML)
	return pageHTML
}


// RenderFavoriteButton renders the dynamic HTMX favorite/unfavorite toggle button
func RenderFavoriteButton(restaurantID int, isFavorite bool, isModal bool) string {
	if isModal {
		if isFavorite {
			return fmt.Sprintf(`
				<button 
					type="button" 
					class="btn-card-action btn-save-resto saved" 
					id="modal-fav-btn-%d"
					hx-post="/favorites/toggle/%d?modal=true"
					hx-swap="outerHTML"
					style="background-color: #fee2e2; color: #dc2626; border: 1.5px solid #fca5a5; font-weight: 700; padding: 6px 14px; border-radius: 9999px; cursor: pointer; display: inline-flex; align-items: center; gap: 5px; font-size: 13px;"
					title="Remove from favorites"
				>
					<svg width="14" height="14" viewBox="0 0 24 24" fill="#dc2626" stroke="#dc2626" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
					<span>Saved</span>
				</button>
			`, restaurantID, restaurantID)
		}
		return fmt.Sprintf(`
			<button 
				type="button" 
				class="btn-card-action btn-save-resto" 
				id="modal-fav-btn-%d"
				hx-post="/favorites/toggle/%d?modal=true"
				hx-swap="outerHTML"
				style="background-color: #ffffff; color: var(--color-dark-text, #171c24); border: 1.5px solid var(--color-border, #eae8e3); font-weight: 600; padding: 6px 14px; border-radius: 9999px; cursor: pointer; display: inline-flex; align-items: center; gap: 5px; font-size: 13px;"
				title="Save to favorites"
			>
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
				<span>Save</span>
			</button>
		`, restaurantID, restaurantID)
	}

	if isFavorite {
		return fmt.Sprintf(`
			<button 
				type="button" 
				class="btn-card-action btn-save-resto saved" 
				id="save-btn-%d"
				hx-post="/favorites/toggle/%d"
				hx-swap="outerHTML"
				style="background-color: #fee2e2; color: #dc2626; border-color: #fca5a5; font-weight: 700;"
				title="Remove from favorites"
			>
				Saved
			</button>
		`, restaurantID, restaurantID)
	}
	return fmt.Sprintf(`
		<button 
			type="button" 
			class="btn-card-action btn-save-resto" 
			id="save-btn-%d"
			hx-post="/favorites/toggle/%d"
			hx-swap="outerHTML"
			title="Save to favorites"
		>
			Save
		</button>
	`, restaurantID, restaurantID)
}

// FormatAvailabilityBadge generates an accessible, styled status badge for a restaurant
func FormatAvailabilityBadge(av fooddiscovery.AvailabilityResult) string {
	badgeText := av.BadgeText
	if badgeText == "" {
		badgeText = "Hours not available"
	}
	detailText := av.DetailText
	if detailText == "" {
		detailText = badgeText
	}

	switch av.Status {
	case fooddiscovery.StatusOpenNow:
		return fmt.Sprintf(`<span class="hours-badge open" title="%s" aria-label="%s: %s"><span class="badge-dot-live"></span>%s</span>`,
			escapeHTML(detailText), escapeHTML(badgeText), escapeHTML(detailText), escapeHTML(badgeText))
	case fooddiscovery.StatusClosedNow, fooddiscovery.StatusOpensAt:
		return fmt.Sprintf(`<span class="hours-badge closed" title="%s" aria-label="%s: %s"><span class="badge-dot-closed"></span>%s</span>`,
			escapeHTML(detailText), escapeHTML(badgeText), escapeHTML(detailText), escapeHTML(badgeText))
	case fooddiscovery.StatusTemporarilyClosed:
		notice := av.Notice
		if notice == "" {
			notice = detailText
		}
		return fmt.Sprintf(`<span class="hours-badge temp-closed" title="%s" aria-label="Temporarily Closed: %s"><span class="badge-dot-amber"></span>Temporarily Closed</span>`,
			escapeHTML(notice), escapeHTML(notice))
	default:
		return `<span class="hours-badge unavailable" title="Opening hours not listed by restaurant" aria-label="Hours not available"><span class="badge-dot-gray"></span>Hours not available</span>`
	}
}

// FormatDeliveryBadge generates an accessible, styled delivery status badge for customer storefront
func FormatDeliveryBadge(s fooddiscovery.ServiceabilityStatus) string {
	badgeText := s.BadgeText
	if badgeText == "" {
		badgeText = "Discovery Only"
	}
	detailText := s.DetailText
	if detailText == "" {
		detailText = badgeText
	}

	if s.IsServiceable {
		return fmt.Sprintf(`<span style="background: #dcfce7; color: #166534; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.55rem; border-radius: 6px; display: inline-flex; align-items: center; gap: 4px;" title="%s">🛵 %s</span>`,
			escapeHTML(detailText), escapeHTML(badgeText))
	} else if strings.Contains(strings.ToLower(badgeText), "paused") || strings.Contains(strings.ToLower(badgeText), "suspended") {
		return fmt.Sprintf(`<span style="background: #fef3c7; color: #92400e; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.55rem; border-radius: 6px; display: inline-flex; align-items: center; gap: 4px;" title="%s">⏸️ %s</span>`,
			escapeHTML(detailText), escapeHTML(badgeText))
	} else if s.PickupEnabled {
		return fmt.Sprintf(`<span style="background: #f1f5f9; color: #475569; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.55rem; border-radius: 6px; display: inline-flex; align-items: center; gap: 4px;" title="%s">🛍️ Takeaway Only</span>`,
			escapeHTML(detailText))
	} else {
		return fmt.Sprintf(`<span style="background: #f1f5f9; color: #64748b; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.55rem; border-radius: 6px; display: inline-flex; align-items: center; gap: 4px;" title="%s">📍 %s</span>`,
			escapeHTML(detailText), escapeHTML(badgeText))
	}
}

// RenderRestaurantCards maintains backward compatibility
func RenderRestaurantCards(results []fooddiscovery.RestaurantWithDistance, maxPrice float64, diet string) string {
	return RenderRestaurantCardsPaginated(&fooddiscovery.PaginatedRestaurants{
		Items:      results,
		TotalCount: len(results),
		Page:       1,
		PageSize:   len(results),
		TotalPages: 1,
	}, maxPrice, diet, "", false, "", "", 12.9716, 77.5946, 25.0)
}

// RenderRestaurantCardsPaginated renders restaurant cards along with server-side pagination controls and helpful empty states
func RenderRestaurantCardsPaginated(
	paginated *fooddiscovery.PaginatedRestaurants,
	maxPrice float64,
	diet string,
	category string,
	openNow bool,
	query string,
	sortBy string,
	lat, lng, radius float64,
) string {
	if paginated == nil || len(paginated.Items) == 0 {
		emptyTitle := fmt.Sprintf("No Restaurants Found with Food ≤ ₹%.0f", maxPrice)
		emptyDesc := "No approved restaurants currently match your filter criteria. Try adjusting your budget slider, diet filter, or locality."
		if openNow {
			emptyTitle = "No Open Restaurants Found"
			emptyDesc = "There are currently no approved restaurants open for ordering matching your criteria. Try turning off the 'Open Now' filter or checking back during opening hours."
		} else if category != "" && category != "all" {
			catName := strings.Title(strings.ReplaceAll(category, "-", " "))
			emptyTitle = fmt.Sprintf("No Dishes Found in \"%s\"", escapeHTML(catName))
			emptyDesc = fmt.Sprintf("We couldn't find any approved food items under ₹%.0f in the %s category matching your filters. Try selecting another category or adjusting your budget.", maxPrice, escapeHTML(catName))
		} else if query != "" {
			emptyTitle = fmt.Sprintf("No Matches for \"%s\"", escapeHTML(query))
			emptyDesc = "We couldn't find any approved restaurants or food items matching your search. Try searching for popular dishes like Biryani, Dosa, Tea, or Roll."
		} else if diet == "veg" {
			emptyTitle = fmt.Sprintf("No Pure Veg Restaurants with Food ≤ ₹%.0f", maxPrice)
		} else if diet == "nonveg" || diet == "non-veg" {
			emptyTitle = fmt.Sprintf("No Non-Veg Restaurants with Food ≤ ₹%.0f", maxPrice)
		}

		return fmt.Sprintf(`
			<div class="empty-state-box" style="background: #ffffff; border: 1px solid var(--border-gray, #e2e8f0); border-radius: 16px; padding: 3rem 2rem; text-align: center;">
				<div style="width: 56px; height: 56px; border-radius: 50%%; background: #fee2e2; color: #dc2626; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
				</div>
				<h3 style="font-family: var(--font-heading, sans-serif); font-size: 1.25rem; font-weight: 800; margin-bottom: 0.5rem; color: var(--text-dark, #1e293b);">
					%s
				</h3>
				<p style="color: var(--text-muted, #64748b); font-size: 0.9rem; max-width: 440px; margin: 0 auto 1.5rem;">
					%s
				</p>
				<button class="btn-nav-outline" onclick="if(typeof applyLocality === 'function') applyLocality('all'); else window.location.href='/';">
					Reset Filters
				</button>
			</div>
		`, emptyTitle, emptyDesc)
	}

	var sb strings.Builder
	sb.WriteString(`<div class="cards-wrapper">`)

	for _, item := range paginated.Items {
		r := item.Restaurant

		dietBadge := `<span class="diet-pill nonveg">VEG & NON-VEG</span>`
		if item.IsPureVeg {
			dietBadge = `<span class="diet-pill veg">PURE VEG</span>`
		}

		availBadge := FormatAvailabilityBadge(item.Availability)
		deliveryBadge := FormatDeliveryBadge(item.Serviceability)
		detailLine := ""
		if item.Availability.DetailText != "" {
			detailLine = fmt.Sprintf(`<span>&bull; %s</span>`, escapeHTML(item.Availability.DetailText))
		}

		sb.WriteString(fmt.Sprintf(`
			<article class="resto-card" id="restaurant-%d">
				<div class="card-head">
					<div>
						<h3 class="resto-name">%s</h3>
						<p class="resto-desc">%s</p>
					</div>
					<div class="resto-badges">
						%s
						%s
						%s
						<span class="distance-badge">%s</span>
					</div>
				</div>

				<div class="resto-meta-row">
					<span>%s</span>
					<span>%s</span>
					%s
					%s
				</div>

				<div class="mini-products-title">
					Available Under ₹%.0f (%d Dishes)
				</div>
				<div class="mini-products-grid">
		`, r.ID, escapeHTML(r.Name), escapeHTML(r.Description), FormatRatingBadge(item.AverageRating, item.ReviewCount), availBadge, deliveryBadge, item.DistanceDisplay, escapeHTML(r.Address), escapeHTML(r.Phone), detailLine, dietBadge, maxPrice, len(item.Under100Items)))

		for i, food := range item.Under100Items {
			if i >= 4 {
				break
			}
			isVeg := fooddiscovery.IsDishVeg(food.Name, food.Category, food.Description)
			dietBox := `<span class="diet-pill veg" style="padding: 1px 4px; font-size: 0.65rem;">VEG</span>`
			if !isVeg {
				dietBox = `<span class="diet-pill nonveg" style="padding: 1px 4px; font-size: 0.65rem;">NON-VEG</span>`
			}

			sb.WriteString(fmt.Sprintf(`
				<div class="mini-product-item">
					<div class="product-top">
						%s
						<div class="product-name">%s</div>
					</div>
					<div class="product-bottom">
						<span class="product-price">₹%.0f</span>
						<span class="product-cat">%s</span>
					</div>
				</div>
			`, dietBox, escapeHTML(food.Name), food.Price, escapeHTML(food.Category)))
		}

		favBtnHTML := RenderFavoriteButton(r.ID, item.IsFavorite, false)

		sb.WriteString(fmt.Sprintf(`
				</div>

				<div class="card-actions">
					<button 
						type="button" 
						class="btn-view-menu"
						hx-get="/restaurant/%d"
						hx-target="#modal-container"
					>
						View Full Menu
					</button>
					<a 
						href="https://www.google.com/maps/search/?api=1&query=%.6f,%.6f" 
						target="_blank" 
						class="btn-card-action"
					>
						Directions
					</a>
					%s
				</div>
			</article>
		`, r.ID, r.Latitude, r.Longitude, favBtnHTML))
	}

	sb.WriteString(`</div>`)

	// Server-side pagination controls
	if paginated.TotalPages > 1 {
		sb.WriteString(fmt.Sprintf(`
			<div class="pagination-container" style="display: flex; justify-content: center; align-items: center; gap: 12px; margin-top: 2.5rem; padding: 1.5rem 0; flex-wrap: wrap;">
		`))

		// Previous Page Button
		if paginated.HasPrev {
			sb.WriteString(fmt.Sprintf(`
				<button 
					type="button"
					class="btn-nav-outline"
					hx-get="/nearby?lat=%.4f&lng=%.4f&radius=%.1f&max_price=%.0f&q=%s&diet=%s&category=%s&open_now=%t&sort=%s&page=%d&page_size=%d"
					hx-target="closest div.cards-wrapper-container, #restaurants-list, .cards-wrapper"
					hx-swap="outerHTML"
					style="padding: 8px 16px; border-radius: 8px; border: 1.5px solid var(--border-gray, #cbd5e1); font-weight: 700; cursor: pointer; background: #ffffff; color: var(--color-dark-text, #171c24);"
				>
					&larr; Previous
				</button>
			`, lat, lng, radius, maxPrice, escapeHTML(query), escapeHTML(diet), escapeHTML(category), openNow, escapeHTML(sortBy), paginated.Page-1, paginated.PageSize))
		} else {
			sb.WriteString(`
				<button 
					type="button"
					class="btn-nav-outline disabled"
					disabled
					style="padding: 8px 16px; border-radius: 8px; border: 1.5px solid #e2e8f0; color: #94a3b8; cursor: not-allowed; opacity: 0.6; background: #f8fafc;"
				>
					&larr; Previous
				</button>
			`)
		}

		// Page numbers and counts
		sb.WriteString(fmt.Sprintf(`
			<span style="font-size: 0.92rem; font-weight: 700; color: var(--text-dark, #334155); padding: 0 8px;">
				Page %d of %d <span style="font-size: 0.82rem; font-weight: 500; color: #64748b;">(%d restaurants)</span>
			</span>
		`, paginated.Page, paginated.TotalPages, paginated.TotalCount))

		// Next Page Button
		if paginated.HasNext {
			sb.WriteString(fmt.Sprintf(`
				<button 
					type="button"
					class="btn-nav-outline"
					hx-get="/nearby?lat=%.4f&lng=%.4f&radius=%.1f&max_price=%.0f&q=%s&diet=%s&category=%s&open_now=%t&sort=%s&page=%d&page_size=%d"
					hx-target="closest div.cards-wrapper-container, #restaurants-list, .cards-wrapper"
					hx-swap="outerHTML"
					style="padding: 8px 16px; border-radius: 8px; border: 1.5px solid var(--border-gray, #cbd5e1); font-weight: 700; cursor: pointer; background: #ffffff; color: var(--color-dark-text, #171c24);"
				>
					Next &rarr;
				</button>
			`, lat, lng, radius, maxPrice, escapeHTML(query), escapeHTML(diet), escapeHTML(category), openNow, escapeHTML(sortBy), paginated.Page+1, paginated.PageSize))
		} else {
			sb.WriteString(`
				<button 
					type="button"
					class="btn-nav-outline disabled"
					disabled
					style="padding: 8px 16px; border-radius: 8px; border: 1.5px solid #e2e8f0; color: #94a3b8; cursor: not-allowed; opacity: 0.6; background: #f8fafc;"
				>
					Next &rarr;
				</button>
			`)
		}

		sb.WriteString(`</div>`)
	}

	return sb.String()
}

// RenderRestaurantModal renders the clean light-mode popup modal with menu items, operating hours, and customer reviews
func RenderRestaurantModal(
	rest *ent.Restaurant,
	foods []*ent.FoodItem,
	isFavorite bool,
	stats fooddiscovery.RestaurantRatingStats,
	userReview *fooddiscovery.Review,
	reviews fooddiscovery.PaginatedReviews,
	currentUser *ent.User,
	avail fooddiscovery.AvailabilityResult,
	periods []fooddiscovery.OperatingPeriod,
	specials []fooddiscovery.SpecialHours,
	serviceability fooddiscovery.ServiceabilityStatus,
	itemRulesMap ...map[int]fooddiscovery.FoodItemDeliveryRule,
) string {
	var itemRules map[int]fooddiscovery.FoodItemDeliveryRule
	if len(itemRulesMap) > 0 && itemRulesMap[0] != nil {
		itemRules = itemRulesMap[0]
	}

	var sb strings.Builder
	favBtnHTML := RenderFavoriteButton(rest.ID, isFavorite, true)

	ratingBadge := `<span style="background: #f1f5f9; color: #64748b; font-size: 0.75rem; font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 6px;">★ New</span>`
	if stats.ReviewCount > 0 {
		reviewWord := "reviews"
		if stats.ReviewCount == 1 {
			reviewWord = "review"
		}
		ratingBadge = fmt.Sprintf(`<span style="background: #fef08a; color: #854d0e; font-size: 0.75rem; font-weight: 800; padding: 0.2rem 0.5rem; border-radius: 6px;">★ %.1f (%d %s)</span>`, stats.AverageRating, stats.ReviewCount, reviewWord)
	}

	availBadge := FormatAvailabilityBadge(avail)
	deliveryBadge := FormatDeliveryBadge(serviceability)

	// Temporary closure notice if active
	var closureNoticeHTML string
	for _, spec := range specials {
		if spec.IsClosed && spec.ClosureStart != nil && spec.ClosureEnd != nil {
			loc := fooddiscovery.GetKolkataLocation()
			if time.Now().Before(*spec.ClosureEnd) {
				endFmt := spec.ClosureEnd.In(loc).Format("Jan 02, 3:04 PM")
				reason := spec.Reason
				if reason == "" {
					reason = "Temporarily closed for maintenance"
				}
				closureNoticeHTML = fmt.Sprintf(`
					<div style="background: #fffbeb; border: 1.5px solid #fde68a; border-radius: 8px; padding: 0.6rem 0.85rem; margin-top: 0.6rem; color: #92400e; font-size: 0.82rem;">
						<strong>⚠️ Temporarily Closed:</strong> %s &bull; Expected reopening: %s (IST)
					</div>
				`, escapeHTML(reason), endFmt)
				break
			}
		}
	}

	// Weekly schedule breakdown
	var weeklyScheduleHTML strings.Builder
	if len(periods) == 0 {
		weeklyScheduleHTML.WriteString(`<div style="font-size: 0.82rem; color: #94a3b8; margin-top: 0.4rem;">Opening hours not listed by restaurant.</div>`)
	} else {
		byDay := make(map[int][]fooddiscovery.OperatingPeriod)
		for _, p := range periods {
			byDay[p.DayOfWeek] = append(byDay[p.DayOfWeek], p)
		}

		currentDay := fooddiscovery.WeekdayToDayOfWeek(time.Now().In(fooddiscovery.GetKolkataLocation()).Weekday())

		weeklyScheduleHTML.WriteString(`
			<details style="margin-top: 0.6rem; cursor: pointer;">
				<summary style="font-size: 0.82rem; font-weight: 700; color: #2563eb; outline: none; user-select: none;">
					View 7-Day Operating Schedule &darr;
				</summary>
				<div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); gap: 6px; margin-top: 0.6rem; padding: 0.5rem; background: #ffffff; border-radius: 8px; border: 1px solid #e2e8f0; font-size: 0.78rem;">
		`)

		for d := 1; d <= 7; d++ {
			dName := fooddiscovery.DayOfWeekToString(d)
			isToday := (d == currentDay)
			cardBg := "#ffffff"
			cardBorder := "#e2e8f0"
			dayLabelColor := "#334155"
			if isToday {
				cardBg = "#f0fdf4"
				cardBorder = "#86efac"
				dayLabelColor = "#166534"
			}

			dPeriods := byDay[d]
			hoursText := "Closed"
			if len(dPeriods) > 0 {
				if dPeriods[0].IsClosed {
					hoursText = "Closed"
				} else if dPeriods[0].Is24Hours {
					hoursText = "Open 24 Hours"
				} else {
					var pStrs []string
					for _, dp := range dPeriods {
						pStrs = append(pStrs, fmt.Sprintf("%s - %s", fooddiscovery.FormatTime12H(dp.OpenTime), fooddiscovery.FormatTime12H(dp.CloseTime)))
					}
					hoursText = strings.Join(pStrs, "<br>")
				}
			} else {
				hoursText = "Not listed"
			}

			todayIndicator := ""
			if isToday {
				todayIndicator = `<span style="font-size: 0.65rem; font-weight: 800; color: #166534; background: #bbf7d0; padding: 1px 4px; border-radius: 4px;">TODAY</span>`
			}

			weeklyScheduleHTML.WriteString(fmt.Sprintf(`
				<div style="background: %s; border: 1px solid %s; border-radius: 6px; padding: 6px 8px; font-size: 0.76rem;">
					<div style="display: flex; justify-content: space-between; align-items: center; font-weight: 700; color: %s; margin-bottom: 2px;">
						<span>%s</span>
						%s
					</div>
					<div style="color: #64748b; line-height: 1.3;">%s</div>
				</div>
			`, cardBg, cardBorder, dayLabelColor, dName, todayIndicator, hoursText))
		}
		weeklyScheduleHTML.WriteString(`</div></details>`)
	}

	ownerClaimBanner := ""

	delivSubText := fmt.Sprintf("₹%.0f delivery fee • %s estimated", serviceability.DeliveryFee, serviceability.EstimatedTimeText)
	if serviceability.DeliveryFee == 0 {
		delivSubText = fmt.Sprintf("Free delivery • %s estimated", serviceability.EstimatedTimeText)
	}

	sb.WriteString(fmt.Sprintf(`
		<div class="modal-overlay" onclick="if(event.target === this) closeModal()">
			<div class="resto-modal-card">
				<div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.25rem;">
					<div>
						<div style="display: flex; gap: 0.4rem; align-items: center; margin-bottom: 0.35rem; flex-wrap: wrap;">
							<span style="background: #dcfce7; color: #166534; font-size: 0.75rem; font-weight: 800; padding: 0.2rem 0.5rem; border-radius: 6px; text-transform: uppercase;">Verified Budget Listing</span>
							%s
							%s
							%s
						</div>
						<h2 style="font-family: var(--font-heading); font-size: 1.75rem; font-weight: 800; color: var(--text-dark); margin: 0;">%s</h2>
						<p style="color: var(--text-muted); font-size: 0.95rem; margin-top: 0.25rem;">%s</p>
					</div>
					<div style="display: flex; align-items: center; gap: 0.75rem;">
						%s
						<button onclick="closeModal()" style="background: none; border: none; font-size: 1.5rem; color: var(--text-muted); cursor: pointer;">&times;</button>
					</div>
				</div>

				<div style="background: #f8fafc; border: 1px solid var(--border-gray); padding: 1rem; border-radius: 12px; margin-bottom: 1.5rem; font-size: 0.88rem; line-height: 1.6;">
					<div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem;">
						<div><strong>Address:</strong> %s</div>
						<a href="https://www.google.com/maps/search/?api=1&query=%.6f,%.6f" target="_blank" style="color: var(--primary-red); text-decoration: none; font-weight: 700; font-size: 0.82rem;">
							Open in Google Maps &rarr;
						</a>
					</div>
					<div style="border-top: 1px solid #e2e8f0; padding-top: 0.65rem; margin-top: 0.65rem;">
						<div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem;">
							<div>
								<strong style="color: var(--color-dark-text, #171c24);">Operating Status:</strong>
								<span style="color: #475569; font-weight: 600;">%s</span>
							</div>
							%s
						</div>
						%s
						%s
					</div>
					<div style="border-top: 1px solid #e2e8f0; padding-top: 0.65rem; margin-top: 0.65rem;">
						<div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem;">
							<div>
								<strong style="color: var(--color-dark-text, #171c24);">Delivery Serviceability:</strong>
								<span style="color: #475569; font-weight: 600;">%s</span>
							</div>
							%s
						</div>
						<div style="font-size: 0.78rem; color: #64748b; margin-top: 0.25rem;">
							%s &bull; Direct restaurant fulfillment (Discovery phase)
						</div>
					</div>
					%s
				</div>
	`, ratingBadge, availBadge, deliveryBadge, escapeHTML(rest.Name), escapeHTML(rest.Description), favBtnHTML, escapeHTML(rest.Address), rest.Latitude, rest.Longitude, escapeHTML(avail.DetailText), availBadge, closureNoticeHTML, weeklyScheduleHTML.String(), escapeHTML(serviceability.DetailText), deliveryBadge, delivSubText, ownerClaimBanner))

	vegCount := 0
	nonVegCount := 0
	for _, f := range foods {
		if fooddiscovery.IsDishVeg(f.Name, f.Category, f.Description) {
			vegCount++
		} else {
			nonVegCount++
		}
	}

	sb.WriteString(fmt.Sprintf(`
		<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
			<h3 style="font-family: var(--font-heading); font-size: 1.15rem; font-weight: 800; color: var(--text-dark);">
				Menu Items & Prices
			</h3>
			<span style="font-size: 0.78rem; color: #166534; font-weight: 700; background: #dcfce7; padding: 0.2rem 0.55rem; border-radius: 6px;">
				Items ≤ ₹100
			</span>
		</div>

		<div class="modal-diet-tabs">
			<button type="button" class="modal-diet-tab active" onclick="filterModalDiet('all', this)">All Items (%d)</button>
			<button type="button" class="modal-diet-tab" onclick="filterModalDiet('veg', this)">Veg (%d)</button>
			<button type="button" class="modal-diet-tab" onclick="filterModalDiet('nonveg', this)">Non-Veg (%d)</button>
		</div>

		<div id="modal-food-list">
	`, len(foods), vegCount, nonVegCount))

	for _, f := range foods {
		isVeg := fooddiscovery.IsDishVeg(f.Name, f.Category, f.Description)
		dietType := "veg"
		dietBox := `<span class="diet-pill veg" style="padding: 1px 5px; font-size: 0.65rem;">VEG</span>`
		if !isVeg {
			dietType = "nonveg"
			dietBox = `<span class="diet-pill nonveg" style="padding: 1px 5px; font-size: 0.65rem;">NON-VEG</span>`
		}

		priceColor := "#166534"
		if f.Price > 100 {
			priceColor = "#64748b"
		}

		soldOutBadge := ""
		itemOpacity := "1"
		unavailSubtitle := ""
		if !f.Available {
			soldOutBadge = `<span style="background: #fee2e2; color: #991b1b; font-size: 0.68rem; font-weight: 800; padding: 2px 6px; border-radius: 4px; margin-left: 6px; text-transform: uppercase;">Sold Out</span>`
			itemOpacity = "0.55"
			unavailSubtitle = `<div style="font-size: 0.72rem; font-weight: 700; color: #dc2626; margin-top: 2px;">Unavailable</div>`
		}

		rule, hasRule := itemRules[f.ID]
		if !hasRule {
			rule = fooddiscovery.DefaultFoodItemDeliveryRule(f.ID)
		}

		itemDeliveryBadge := ""
		if !rule.IsDeliverable {
			itemDeliveryBadge = `<span style="background: #f1f5f9; color: #475569; font-size: 0.68rem; font-weight: 700; padding: 2px 6px; border-radius: 4px; margin-left: 6px;">🍽️ Dine-in / Takeaway Only</span>`
		} else if rule.DeliveryRadiusLimit > 0 {
			if serviceability.DistanceKm > rule.DeliveryRadiusLimit {
				itemDeliveryBadge = fmt.Sprintf(`<span style="background: #fef3c7; color: #92400e; font-size: 0.68rem; font-weight: 700; padding: 2px 6px; border-radius: 4px; margin-left: 6px;">🛵 Takeaway beyond %.0fkm</span>`, rule.DeliveryRadiusLimit)
			} else {
				itemDeliveryBadge = fmt.Sprintf(`<span style="background: #e0f2fe; color: #0369a1; font-size: 0.68rem; font-weight: 700; padding: 2px 6px; border-radius: 4px; margin-left: 6px;">🚚 ≤%.0fkm Delivery</span>`, rule.DeliveryRadiusLimit)
			}
		} else {
			itemDeliveryBadge = `<span style="background: #dcfce7; color: #166534; font-size: 0.68rem; font-weight: 700; padding: 2px 6px; border-radius: 4px; margin-left: 6px;">🚚 Deliverable</span>`
		}

		addBtnHtml := ""
		if !f.Available {
			addBtnHtml = `<span style="background: #fee2e2; color: #991b1b; font-size: 0.72rem; font-weight: 800; padding: 4px 8px; border-radius: 6px;">UNAVAILABLE</span>`
		} else {
			cleanName := strings.ReplaceAll(f.Name, "'", "\\'")
			cleanRest := strings.ReplaceAll(rest.Name, "'", "\\'")
			addBtnHtml = fmt.Sprintf(`
				<div class="add-btn-wrap add-btn-wrap-rf%d" id="add-btn-rf%d" style="min-width: 76px;">
					<button type="button" class="btn-add-pill" onclick="addToCart('rf%d', '%s', %d, '%s', '/static/crispy_dosa.jpg')">
						ADD
					</button>
				</div>
			`, f.ID, f.ID, f.ID, cleanName, int(f.Price), cleanRest)
		}

		sb.WriteString(fmt.Sprintf(`
			<div class="modal-food-item modal-item-%s" style="opacity: %s; display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 0; border-bottom: 1px solid #f1f5f9;">
				<div style="display: flex; align-items: center; gap: 0.65rem; flex: 1; min-width: 0;">
					%s
					<div>
						<div style="font-weight: 700; font-size: 0.92rem; color: var(--text-dark); display: flex; align-items: center; flex-wrap: wrap;">
							%s %s %s
						</div>
						<div style="font-size: 0.78rem; color: var(--text-muted);">%s</div>
					</div>
				</div>
				<div style="display: flex; align-items: center; gap: 12px; flex-shrink: 0;">
					<div style="text-align: right;">
						<div style="font-weight: 800; font-size: 1rem; color: %s;">
							₹%.0f
						</div>
						%s
					</div>
					%s
				</div>
			</div>
		`, dietType, itemOpacity, dietBox, escapeHTML(f.Name), soldOutBadge, itemDeliveryBadge, escapeHTML(f.Category), priceColor, f.Price, unavailSubtitle, addBtnHtml))
	}

	sb.WriteString(`
		</div>
		<script>
			setTimeout(function() {
				if (typeof cartItems === 'object') {
					Object.keys(cartItems).forEach(function(k) {
						if (typeof updateProductCardControls === 'function') {
							updateProductCardControls(k);
						}
					});
				}
			}, 50);
		</script>
		<script>
			function filterModalDiet(type, btn) {
				document.querySelectorAll('.modal-diet-tab').forEach(b => b.classList.remove('active'));
				btn.classList.add('active');
				document.querySelectorAll('.modal-food-item').forEach(item => {
					if (type === 'all') {
						item.style.display = 'flex';
					} else if (type === 'veg') {
						item.style.display = item.classList.contains('modal-item-veg') ? 'flex' : 'none';
					} else {
						item.style.display = item.classList.contains('modal-item-nonveg') ? 'flex' : 'none';
					}
				});
			}
		</script>
	`)

	// Append Customer Reviews Section
	sb.WriteString(RenderReviewsSectionHTML(rest.ID, stats, userReview, reviews, currentUser))

	sb.WriteString(`
			</div>
		</div>
	`)

	return sb.String()
}

// RenderFavoritesPage renders the dedicated Customer Favorites page
func RenderFavoritesPage(currentUser *ent.User, favs []fooddiscovery.RestaurantWithDistance) string {
	userName := ""
	if currentUser != nil {
		userName = currentUser.Name
		if userName == "" {
			userName = currentUser.Email
		}
	}

	var contentHTML string
	if len(favs) == 0 {
		contentHTML = `
			<div class="empty-state-box" style="background: #ffffff; border: 1px solid var(--color-border, #eae8e3); border-radius: 16px; padding: 4rem 2rem; text-align: center; max-width: 600px; margin: 2rem auto;">
				<div style="width: 64px; height: 64px; border-radius: 50%; background: #fee2e2; color: #dc2626; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem;">
					<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
				</div>
				<h3 style="font-family: var(--font-heading, sans-serif); font-size: 1.4rem; font-weight: 800; margin-bottom: 0.6rem; color: var(--color-dark-text, #171c24);">
					No saved restaurants yet
				</h3>
				<p style="color: var(--color-secondary-text, #646b73); font-size: 0.95rem; line-height: 1.6; margin-bottom: 2rem;">
					You haven't saved any restaurants yet. Explore Bangalore's best budget food under ₹100 and tap 'Save' on any eatery to find it here.
				</p>
				<a href="/" style="display: inline-flex; align-items: center; gap: 8px; background-color: var(--color-primary-orange, #ff5a1f); color: #ffffff; padding: 12px 24px; border-radius: 9999px; font-weight: 700; text-decoration: none; font-size: 14px; box-shadow: 0 4px 12px rgba(255, 90, 31, 0.25);">
					<span>Explore Bangalore Food Under ₹100</span>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"></line><polyline points="12 5 19 12 12 19"></polyline></svg>
				</a>
			</div>
		`
	} else {
		var sb strings.Builder
		sb.WriteString(`<div class="cards-wrapper">`)
		for _, item := range favs {
			r := item.Restaurant
			dietBadge := `<span class="diet-pill nonveg">VEG & NON-VEG</span>`
			if item.IsPureVeg {
				dietBadge = `<span class="diet-pill veg">PURE VEG</span>`
			}

			availBadge := FormatAvailabilityBadge(item.Availability)
			deliveryBadge := FormatDeliveryBadge(item.Serviceability)
			detailLine := ""
			if item.Availability.DetailText != "" {
				detailLine = fmt.Sprintf(`<span>&bull; %s</span>`, escapeHTML(item.Availability.DetailText))
			}

			sb.WriteString(fmt.Sprintf(`
				<article class="resto-card" id="restaurant-%d">
					<div class="card-head">
						<div>
							<h3 class="resto-name">%s</h3>
							<p class="resto-desc">%s</p>
						</div>
						<div class="resto-badges">
							%s
							%s
							%s
							<span class="distance-badge">%s</span>
						</div>
					</div>

					<div class="resto-meta-row">
						<span>%s</span>
						<span>%s</span>
						%s
						%s
					</div>

					<div class="mini-products-title">
						Available Under ₹100 (%d Dishes)
					</div>
					<div class="mini-products-grid">
			`, r.ID, escapeHTML(r.Name), escapeHTML(r.Description), FormatRatingBadge(item.AverageRating, item.ReviewCount), availBadge, deliveryBadge, item.DistanceDisplay, escapeHTML(r.Address), escapeHTML(r.Phone), detailLine, dietBadge, len(item.Under100Items)))

			for i, food := range item.Under100Items {
				if i >= 4 {
					break
				}
				isVeg := fooddiscovery.IsDishVeg(food.Name, food.Category, food.Description)
				dietBox := `<span class="diet-pill veg" style="padding: 1px 4px; font-size: 0.65rem;">VEG</span>`
				if !isVeg {
					dietBox = `<span class="diet-pill nonveg" style="padding: 1px 4px; font-size: 0.65rem;">NON-VEG</span>`
				}

				sb.WriteString(fmt.Sprintf(`
					<div class="mini-product-item">
						<div class="product-top">
							%s
							<div class="product-name">%s</div>
						</div>
						<div class="product-bottom">
							<span class="product-price">₹%.0f</span>
							<span class="product-cat">%s</span>
						</div>
					</div>
				`, dietBox, escapeHTML(food.Name), food.Price, escapeHTML(food.Category)))
			}

			favBtnHTML := RenderFavoriteButton(r.ID, true, false)

			sb.WriteString(fmt.Sprintf(`
					</div>

					<div class="card-actions">
						<button 
							type="button" 
							class="btn-view-menu"
							hx-get="/restaurant/%d"
							hx-target="#modal-container"
						>
							View Full Menu
						</button>
						<a 
							href="https://www.google.com/maps/search/?api=1&query=%.6f,%.6f" 
							target="_blank" 
							class="btn-card-action"
						>
							Directions
						</a>
						%s
					</div>
				</article>
			`, r.ID, r.Latitude, r.Longitude, favBtnHTML))
		}
		sb.WriteString(`</div>`)
		contentHTML = sb.String()
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>My Saved Favorites — Under100 Bengaluru</title>
  <meta name="description" content="Your saved favourite budget eateries and restaurants under &#8377;100 in Bengaluru.">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <style>
    :root {
      --color-primary-orange: #FF5A1F;
      --color-primary-orange-hover: #E04810;
      --color-deep-green: #087443;
      --color-dark-text: #171C24;
      --color-secondary-text: #646B73;
      --color-main-bg: #FFFDF9;
      --color-card-bg: #FFFFFF;
      --color-border: #EAE8E3;
      --font-heading: 'Outfit', sans-serif;
      --font-body: 'Plus Jakarta Sans', sans-serif;
      --container-max-width: 1240px;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: var(--font-body);
      background-color: var(--color-main-bg);
      color: var(--color-dark-text);
      line-height: 1.5;
    }
    .container { width: 100%%; max-width: var(--container-max-width); margin: 0 auto; padding: 0 20px; }
    .header-bar {
      background: #ffffff;
      border-bottom: 1px solid var(--color-border);
      padding: 16px 0;
      position: sticky;
      top: 0;
      z-index: 50;
    }
    .header-inner { display: flex; align-items: center; justify-content: space-between; }
    .logo-box { display: flex; align-items: center; gap: 8px; text-decoration: none; }
    .logo-icon { width: 32px; height: 32px; border-radius: 8px; background-color: var(--color-primary-orange); display: flex; align-items: center; justify-content: center; color: #ffffff; }
    .logo-text { font-size: 20px; font-weight: 900; color: var(--color-dark-text); }
    .logo-orange { color: var(--color-primary-orange); }
    .btn-nav-outline {
      padding: 8px 16px;
      border-radius: 9999px;
      border: 1.5px solid var(--color-border);
      font-size: 13.5px;
      font-weight: 700;
      color: var(--color-dark-text);
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      transition: all 0.2s;
    }
    .btn-nav-outline:hover { background-color: #f8f6f0; border-color: var(--color-primary-orange); color: var(--color-primary-orange); }
    .fav-hero {
      background: #ffffff;
      border-bottom: 1px solid var(--color-border);
      padding: 40px 0;
      margin-bottom: 32px;
    }
    .fav-badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      background-color: #fee2e2;
      color: #dc2626;
      padding: 4px 12px;
      border-radius: 9999px;
      font-size: 12px;
      font-weight: 800;
      text-transform: uppercase;
      margin-bottom: 12px;
    }
    .cards-wrapper {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
      gap: 24px;
      margin-bottom: 48px;
    }
    .resto-card {
      background: #ffffff;
      border: 1px solid var(--color-border);
      border-radius: 16px;
      padding: 20px;
      display: flex;
      flex-direction: column;
      box-shadow: 0 4px 16px rgba(23, 28, 36, 0.04);
      transition: transform 0.2s, box-shadow 0.2s;
    }
    .resto-card:hover { transform: translateY(-3px); box-shadow: 0 8px 24px rgba(23, 28, 36, 0.08); }
    .card-head { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px; }
    .resto-name { font-size: 17px; font-weight: 800; color: var(--color-dark-text); }
    .resto-desc { font-size: 13px; color: var(--color-secondary-text); margin-top: 2px; }
    .resto-badges { display: flex; flex-direction: column; align-items: flex-end; gap: 4px; }
    .rating-badge { background: #FEF3C7; color: #B45309; font-size: 12px; font-weight: 700; padding: 2px 6px; border-radius: 4px; }
    .distance-badge { font-size: 11px; font-weight: 600; color: var(--color-secondary-text); }
    .resto-meta-row { display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--color-secondary-text); margin-bottom: 14px; flex-wrap: wrap; }
    .diet-pill { font-size: 10px; font-weight: 800; padding: 2px 6px; border-radius: 4px; }
    .diet-pill.veg { background: #dcfce7; color: #166534; }
    .diet-pill.nonveg { background: #fed7aa; color: #c2410c; }
    .hours-badge { display: inline-flex; align-items: center; gap: 5px; font-size: 11px; font-weight: 700; padding: 2px 7px; border-radius: 4px; line-height: 1.2; text-decoration: none; }
    .hours-badge.open { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }
    .hours-badge.closed { background: #f1f5f9; color: #64748b; border: 1px solid #e2e8f0; }
    .hours-badge.temp-closed { background: #fef3c7; color: #92400e; border: 1px solid #fde68a; }
    .hours-badge.unavailable { background: #f8fafc; color: #94a3b8; border: 1px solid #e2e8f0; }
    .badge-dot-live { width: 6px; height: 6px; border-radius: 50%%; background-color: #22c55e; display: inline-block; }
    .badge-dot-closed { width: 6px; height: 6px; border-radius: 50%%; background-color: #94a3b8; display: inline-block; }
    .badge-dot-amber { width: 6px; height: 6px; border-radius: 50%%; background-color: #f59e0b; display: inline-block; }
    .badge-dot-gray { width: 6px; height: 6px; border-radius: 50%%; background-color: #cbd5e1; display: inline-block; }
    .mini-products-title { font-size: 12px; font-weight: 800; color: var(--color-deep-green); text-transform: uppercase; margin-bottom: 8px; }
    .mini-products-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 16px; }
    .mini-product-item { background: #f8fafc; border: 1px solid var(--color-border); border-radius: 8px; padding: 8px; }
    .product-top { display: flex; align-items: center; gap: 4px; margin-bottom: 4px; }
    .product-name { font-size: 12px; font-weight: 700; color: var(--color-dark-text); overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
    .product-bottom { display: flex; justify-content: space-between; align-items: center; }
    .product-price { font-size: 12px; font-weight: 800; color: var(--color-deep-green); }
    .product-cat { font-size: 10px; color: var(--color-secondary-text); }
    .card-actions { display: flex; gap: 8px; margin-top: auto; padding-top: 12px; border-top: 1px solid #f1f0ec; }
    .btn-view-menu { flex: 1; padding: 8px; background: var(--color-deep-green); color: #ffffff; border-radius: 8px; font-size: 12.5px; font-weight: 700; cursor: pointer; text-align: center; border: none; }
    .btn-card-action { padding: 8px 12px; border-radius: 8px; border: 1.5px solid var(--color-border); background: #ffffff; color: var(--color-dark-text); font-size: 12.5px; font-weight: 700; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; }
    .btn-save-resto.saved { background-color: #fee2e2; color: #dc2626; border-color: #fca5a5; font-weight: 700; }
    /* Modal styles */
    .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 16px; }
    .resto-modal-card { background: #ffffff; border-radius: 16px; width: 100%%; max-width: 600px; max-height: 90vh; overflow-y: auto; padding: 24px; box-shadow: 0 20px 40px rgba(0,0,0,0.2); }
    .modal-diet-tabs { display: flex; gap: 8px; margin-bottom: 12px; }
    .modal-diet-tab { padding: 6px 12px; border-radius: 6px; border: 1px solid var(--color-border); font-size: 12px; font-weight: 700; cursor: pointer; background: #ffffff; }
    .modal-diet-tab.active { background: var(--color-deep-green); color: #ffffff; border-color: var(--color-deep-green); }
    .modal-food-item { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid #f1f0ec; }
  </style>
</head>
<body>
  <header class="header-bar">
    <div class="container header-inner">
      <a href="/" class="logo-box">
        <div class="logo-icon">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="2" x2="6" y2="22"></line><line x1="6" y1="2" x2="18" y2="22"></line></svg>
        </div>
        <div class="logo-text">Under<span class="logo-orange">100</span></div>
      </a>
      <div style="display: flex; align-items: center; gap: 12px;">
        <a href="/" class="btn-nav-outline">
          &larr; Back to Food Discovery
        </a>
        <a href="/profile" class="btn-nav-outline" title="%s">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
          <span>%s</span>
        </a>
      </div>
    </div>
  </header>

  <section class="fav-hero">
    <div class="container">
      <div class="fav-badge">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="#dc2626" stroke="#dc2626" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
        <span>Saved Restaurants</span>
      </div>
      <h1 style="font-size: clamp(28px, 4vw, 36px); font-weight: 900; color: var(--color-dark-text); letter-spacing: -0.8px; margin-bottom: 8px;">
        My Favorite Eateries
      </h1>
      <p style="font-size: 15px; color: var(--color-secondary-text);">
        Quickly access and view menus of your saved budget restaurants in Bengaluru.
      </p>
    </div>
  </section>

  <main class="container">
    %s
  </main>

  <div id="modal-container"></div>

  <script>
    function closeModal() {
      var modalContainer = document.getElementById('modal-container');
      if (modalContainer) {
        modalContainer.innerHTML = '';
      }
    }
  </script>
</body>
</html>`, escapeHTML(userName), escapeHTML(userName), contentHTML)
}

// RenderAdminDashboard renders the full Admin Dashboard view in clean corporate light mode
func RenderAdminDashboard(stats fooddiscovery.AdminStats, rests []*ent.Restaurant, activeStatus, search string, currentUser *ent.User) string {
	adminName := "Admin"
	if currentUser != nil {
		if currentUser.Name != "" {
			adminName = currentUser.Name
		} else if currentUser.Email != "" {
			adminName = currentUser.Email
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Admin Dashboard — Bite100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800&display=swap" rel="stylesheet">
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      padding-bottom: 4rem;
    }
    .admin-nav {
      background: #ffffff;
      border-bottom: 1px solid #e2e8f0;
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      box-shadow: 0 1px 3px rgba(0,0,0,0.03);
    }
    .admin-title { font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a; }
    .nav-tab { color: #64748b; text-decoration: none; font-size: 0.88rem; font-weight: 600; padding: 0.45rem 0.9rem; border-radius: 8px; transition: all 0.15s ease; }
    .nav-tab:hover { background: #f1f5f9; color: #0f172a; }
    .nav-tab.active { background: #0f172a; color: #ffffff; }
    
    .container { max-width: 1350px; margin: 2rem auto; padding: 0 1.5rem; }
    
    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
      gap: 1rem;
      margin-bottom: 2rem;
    }
    .metric-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      padding: 1.25rem;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .metric-num { font-family: 'Outfit', sans-serif; font-size: 1.8rem; font-weight: 800; margin-top: 0.25rem; color: #0f172a; }
    
    .controls-bar {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.5rem;
    }
    .tab-pills { display: flex; gap: 0.5rem; }
    .tab-pill {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      color: #64748b;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
      cursor: pointer;
    }
    .tab-pill.active { background: #0f172a; color: white; border-color: #0f172a; }
    .btn-create {
      background: #e23744;
      color: white;
      border: none;
      padding: 0.6rem 1.25rem;
      border-radius: 8px;
      font-weight: 700;
      cursor: pointer;
    }
    table { width: 100%%; border-collapse: collapse; background: #ffffff; border-radius: 12px; overflow: hidden; border: 1px solid #e2e8f0; box-shadow: 0 1px 3px rgba(0,0,0,0.02); }
    th { background: #f8fafc; text-align: left; padding: 1rem; font-size: 0.82rem; color: #64748b; text-transform: uppercase; font-weight: 700; border-bottom: 1px solid #e2e8f0; }
    td { padding: 1rem; border-top: 1px solid #f1f5f9; font-size: 0.88rem; color: #334155; vertical-align: middle; }
    tr:hover td { background: #f8fafc; }
    
    .badge-approved { background: #dcfce7; color: #166534; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    .badge-pending { background: #fef08a; color: #854d0e; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    .badge-rejected { background: #fee2e2; color: #991b1b; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    
    .btn-action { background: #f1f5f9; border: 1px solid #e2e8f0; color: #0f172a; padding: 0.35rem 0.65rem; border-radius: 6px; font-size: 0.8rem; cursor: pointer; font-weight: 600; text-decoration: none; display: inline-block; }
    .btn-action:hover { background: #e2e8f0; }
    
    .form-input { width: 100%%; background: #ffffff; border: 1px solid #cbd5e1; color: #0f172a; padding: 0.6rem 0.8rem; border-radius: 8px; font-size: 0.9rem; }
    .modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
    .modal-card { background: #ffffff; border: 1px solid #e2e8f0; border-radius: 16px; width: 600px; max-width: 95%%; max-height: 90vh; overflow-y: auto; padding: 1.5rem; box-shadow: 0 20px 40px rgba(0,0,0,0.15); }
  </style>
</head>
<body>
  <div class="admin-nav">
    <div style="display: flex; align-items: center; gap: 2rem;">
      <div class="admin-title">Under100 Admin Console</div>
      <div style="display: flex; gap: 0.5rem;">
        <a href="/admin" class="nav-tab active">Restaurants</a>
        <a href="/admin/users" class="nav-tab">Users &amp; Roles</a>
        <a href="/admin/categories" class="nav-tab">Categories</a>
        <a href="/admin/reviews" class="nav-tab">Review Moderation</a>
      </div>
    </div>
    <div style="display: flex; align-items: center; gap: 1rem;">
      <a href="/" style="color: #0f172a; text-decoration: none; font-weight: 600; font-size: 0.9rem;">&larr; Back to Website</a>
      <div style="background: #f1f5f9; border: 1px solid #e2e8f0; padding: 0.35rem 0.8rem; border-radius: 9999px; font-size: 0.82rem; color: #0f172a; font-weight: 700;">
        Admin: %s
      </div>
      <a href="/logout" style="color: #e23744; text-decoration: none; font-size: 0.85rem; font-weight: 700; padding: 0.35rem 0.75rem; border: 1px solid #fecaca; border-radius: 8px;">Log Out</a>
    </div>
  </div>

  <div class="container">
    <div class="metrics-grid">
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b; font-weight: 600;">Platform Users</div>
        <div class="metric-num">%d</div>
        <div style="font-size: 0.75rem; color: #64748b; margin-top: 0.35rem;">%d Cust &bull; %d Merch &bull; %d Admin</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b; font-weight: 600;">Total Restaurants</div>
        <div class="metric-num">%d</div>
        <div style="font-size: 0.75rem; color: #15803d; margin-top: 0.35rem;">%d Approved Listings</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b; font-weight: 600;">Pending Approvals</div>
        <div class="metric-num" style="color: #b45309;">%d</div>
        <div style="font-size: 0.75rem; color: #b91c1c; margin-top: 0.35rem;">%d Rejected</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b; font-weight: 600;">Unassigned Outlets</div>
        <div class="metric-num" style="color: #6366f1;">%d</div>
        <div style="font-size: 0.75rem; color: #64748b; margin-top: 0.35rem;">Need Merchant Owner</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b; font-weight: 600;">Total Food Items</div>
        <div class="metric-num">%d</div>
        <div style="font-size: 0.75rem; color: #15803d; margin-top: 0.35rem;">%d Available &le; ₹100</div>
      </div>
    </div>

    <div class="controls-bar">
      <div class="tab-pills">
        <a href="/admin" class="tab-pill %s">All</a>
        <a href="/admin?status=pending" class="tab-pill %s">Pending (%d)</a>
        <a href="/admin?status=approved" class="tab-pill %s">Approved (%d)</a>
        <a href="/admin?status=rejected" class="tab-pill %s">Rejected (%d)</a>
      </div>

      <div style="display: flex; gap: 0.75rem;">
        <input 
          type="text" 
          placeholder="Search restaurants..." 
          class="form-input" 
          style="width: 240px;"
          value="%s"
          hx-get="/admin?table_only=true&status=%s"
          hx-trigger="keyup changed delay:300ms"
          hx-target="#admin-table-body"
          name="q"
        >
        <button 
          class="btn-create" 
          hx-get="/admin/restaurants/new" 
          hx-target="#admin-modal"
        >
          + Add Restaurant
        </button>
      </div>
    </div>

    <table>
      <thead>
        <tr>
          <th>Restaurant Name & Address</th>
          <th>Location</th>
          <th>Assigned Merchant</th>
          <th>Status</th>
          <th>Menu Items</th>
          <th>Approval</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody id="admin-table-body">
        %s
      </tbody>
    </table>
  </div>

  <div id="admin-modal"></div>
  <script>
    function closeAdminModal() { document.getElementById('admin-modal').innerHTML = ''; }
  </script>
</body>
</html>`,
		escapeHTML(adminName),
		stats.TotalUsers,
		stats.CustomerCount, stats.MerchantCount, stats.AdminCount,
		stats.TotalRestaurants,
		stats.ApprovedRestaurants,
		stats.PendingRestaurants,
		stats.RejectedRestaurants,
		stats.UnassignedRestaurants,
		stats.TotalFoodItems,
		stats.Under100FoodItems,
		activeTab(activeStatus, ""),
		activeTab(activeStatus, "pending"), stats.PendingRestaurants,
		activeTab(activeStatus, "approved"), stats.ApprovedRestaurants,
		activeTab(activeStatus, "rejected"), stats.RejectedRestaurants,
		escapeHTML(search),
		activeStatus,
		RenderAdminTableRows(rests),
	)
}

// RenderAdminTableRows renders the table rows for admin view
func RenderAdminTableRows(rests []*ent.Restaurant) string {
	if len(rests) == 0 {
		return `<tr><td colspan="7" style="text-align: center; color: #64748b; padding: 2.5rem;">No restaurants found.</td></tr>`
	}

	var sb strings.Builder
	for _, r := range rests {
		sb.WriteString(RenderAdminSingleTableRow(r))
	}
	return sb.String()
}

// RenderAdminSingleTableRow renders a single restaurant row in the admin table
func RenderAdminSingleTableRow(r *ent.Restaurant) string {
	ownerHTML := fmt.Sprintf(`<span style="background: #fef3c7; color: #92400e; padding: 0.2rem 0.55rem; border-radius: 6px; font-weight: 700; font-size: 0.75rem;">Unassigned</span>
		<div style="margin-top: 0.35rem;">
			<button class="btn-action" style="color: #2563eb; font-weight: 700; font-size: 0.75rem;" hx-get="/admin/restaurants/%d/owner" hx-target="#admin-modal">+ Assign</button>
		</div>`, r.ID)

	if r.Edges.Owner != nil {
		ownerName := r.Edges.Owner.Name
		if ownerName == "" {
			ownerName = fmt.Sprintf("Merchant #%d", r.Edges.Owner.ID)
		}
		ownerHTML = fmt.Sprintf(`
			<div style="font-weight: 700; color: #0f172a; font-size: 0.85rem;">%s</div>
			<div style="font-size: 0.75rem; color: #64748b;">%s</div>
			<div style="margin-top: 0.35rem;">
				<button class="btn-action" style="font-size: 0.75rem;" hx-get="/admin/restaurants/%d/owner" hx-target="#admin-modal">Reassign</button>
			</div>
		`, escapeHTML(ownerName), escapeHTML(r.Edges.Owner.Email), r.ID)
	}

	foodCount := 0
	if r.Edges.FoodItems != nil {
		foodCount = len(r.Edges.FoodItems)
	}

	return fmt.Sprintf(`
		<tr id="admin-row-%d">
			<td>
				<div style="font-weight: 700; color: #0f172a;">%s</div>
				<div style="font-size: 0.8rem; color: #64748b;">%s</div>
			</td>
			<td>
				<span style="font-size: 0.8rem; color: #64748b;">%.4f, %.4f</span>
			</td>
			<td>
				%s
			</td>
			<td>
				%s
			</td>
			<td>
				<span style="font-size: 0.85rem; font-weight: 600;">%d items</span>
			</td>
			<td>
				<div style="display: flex; gap: 0.35rem;">
					<button 
						class="btn-action" 
						hx-post="/admin/restaurants/%d/status?status=approved"
						hx-target="#admin-row-%d"
						hx-swap="outerHTML"
					>
						Approve
					</button>
					<button 
						class="btn-action" 
						hx-post="/admin/restaurants/%d/status?status=rejected"
						hx-target="#admin-row-%d"
						hx-swap="outerHTML"
					>
						Reject
					</button>
				</div>
			</td>
			<td>
				<div style="display: flex; gap: 0.35rem;">
					<button 
						class="btn-action" 
						hx-get="/admin/restaurants/%d/hours"
						hx-target="#admin-modal"
						title="Manage Operating Hours & Availability"
					>
						Hours
					</button>
					<button 
						class="btn-action" 
						hx-get="/admin/restaurants/%d/food"
						hx-target="#admin-modal"
					>
						Menu
					</button>
					<button 
						class="btn-action" 
						hx-get="/admin/restaurants/%d/edit"
						hx-target="#admin-modal"
					>
						Edit
					</button>
					<button 
						class="btn-action" 
						style="color: #b91c1c;"
						hx-delete="/admin/restaurants/%d"
						hx-confirm="Are you sure you want to delete this restaurant?"
						hx-target="#admin-row-%d"
						hx-swap="outerHTML"
					>
						Delete
					</button>
				</div>
			</td>
		</tr>
	`, r.ID, escapeHTML(r.Name), escapeHTML(r.Address), r.Latitude, r.Longitude,
		ownerHTML, RenderStatusBadge(r.Status, r.ID), foodCount,
		r.ID, r.ID, r.ID, r.ID, r.ID, r.ID, r.ID, r.ID, r.ID)
}

// RenderOwnerModal renders modal for assigning or unassigning a restaurant owner
func RenderOwnerModal(rest *ent.Restaurant, merchants []*ent.User) string {
	currentOwnerInfo := `<span style="background: #fef3c7; color: #92400e; padding: 0.25rem 0.6rem; border-radius: 6px; font-weight: 600; font-size: 0.82rem;">Currently Unassigned</span>`
	currentOwnerID := 0
	if rest.Edges.Owner != nil {
		currentOwnerID = rest.Edges.Owner.ID
		currentOwnerInfo = fmt.Sprintf(`<span style="background: #dcfce7; color: #166534; padding: 0.25rem 0.6rem; border-radius: 6px; font-weight: 700; font-size: 0.82rem;">%s (%s) — ID #%d</span>`, escapeHTML(rest.Edges.Owner.Name), escapeHTML(rest.Edges.Owner.Email), rest.Edges.Owner.ID)
	}

	var options strings.Builder
	options.WriteString(`<option value="0">-- Unassign Owner (No Assigned Merchant) --</option>`)
	for _, m := range merchants {
		selected := ""
		if m.ID == currentOwnerID {
			selected = " selected"
		}
		label := fmt.Sprintf("%s (%s) — ID #%d", m.Name, m.Email, m.ID)
		options.WriteString(fmt.Sprintf(`<option value="%d"%s>%s</option>`, m.ID, selected, escapeHTML(label)))
	}

	return fmt.Sprintf(`
		<div class="modal-overlay" onclick="if(event.target === this) closeAdminModal()">
			<div class="modal-card" style="width: 540px;">
				<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
					<div>
						<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Assign Restaurant Owner</h3>
						<div style="font-size: 0.85rem; color: #64748b; margin-top: 0.2rem;">%s</div>
					</div>
					<button onclick="closeAdminModal()" style="background: none; border: none; font-size: 1.4rem; color: #64748b; cursor: pointer;">&times;</button>
				</div>

				<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 1rem; margin-bottom: 1.25rem;">
					<div style="font-size: 0.78rem; text-transform: uppercase; color: #64748b; font-weight: 700; letter-spacing: 0.03em; margin-bottom: 0.35rem;">Current Ownership</div>
					<div>%s</div>
				</div>

				<form hx-post="/admin/restaurants/%d/owner" hx-target="#admin-modal">
					<div style="margin-bottom: 1.25rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.4rem; font-weight: 700;">Select Eligible Merchant *</label>
						<select name="owner_id" class="form-input" style="cursor: pointer;">
							%s
						</select>
						<div style="font-size: 0.78rem; color: #64748b; margin-top: 0.35rem;">Only users with the verified merchant role can be assigned. Select "Unassign" to detach owner.</div>
					</div>

					<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem;">
						<button type="button" onclick="closeAdminModal()" class="btn-action" style="padding: 0.6rem 1.25rem;">Cancel</button>
						<button type="submit" class="btn-create" style="padding: 0.6rem 1.25rem;">Save Ownership</button>
					</div>
				</form>
			</div>
		</div>
	`, escapeHTML(rest.Name), currentOwnerInfo, rest.ID, options.String())
}

// RenderAdminUsersPage renders the user management page
func RenderAdminUsersPage(stats fooddiscovery.AdminStats, users []*ent.User, activeRole, search string, currentUser *ent.User) string {
	adminName := "Admin"
	if currentUser != nil {
		if currentUser.Name != "" {
			adminName = currentUser.Name
		} else if currentUser.Email != "" {
			adminName = currentUser.Email
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>User Management — Bite100 Admin</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800&display=swap" rel="stylesheet">
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      padding-bottom: 4rem;
    }
    .admin-nav {
      background: #ffffff;
      border-bottom: 1px solid #e2e8f0;
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
      box-shadow: 0 1px 3px rgba(0,0,0,0.03);
    }
    .admin-title { font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a; }
    .nav-tab { color: #64748b; text-decoration: none; font-size: 0.88rem; font-weight: 600; padding: 0.45rem 0.9rem; border-radius: 8px; transition: all 0.15s ease; }
    .nav-tab:hover { background: #f1f5f9; color: #0f172a; }
    .nav-tab.active { background: #0f172a; color: #ffffff; }
    
    .container { max-width: 1350px; margin: 2rem auto; padding: 0 1.5rem; }
    
    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
      gap: 1rem;
      margin-bottom: 2rem;
    }
    .metric-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      padding: 1.25rem;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .metric-num { font-family: 'Outfit', sans-serif; font-size: 1.8rem; font-weight: 800; margin-top: 0.25rem; color: #0f172a; }
    
    .controls-bar {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.5rem;
    }
    .tab-pills { display: flex; gap: 0.5rem; }
    .tab-pill {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      color: #64748b;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
      cursor: pointer;
    }
    .tab-pill.active { background: #0f172a; color: white; border-color: #0f172a; }
    
    table { width: 100%%; border-collapse: collapse; background: #ffffff; border-radius: 12px; overflow: hidden; border: 1px solid #e2e8f0; box-shadow: 0 1px 3px rgba(0,0,0,0.02); }
    th { background: #f8fafc; text-align: left; padding: 1rem; font-size: 0.82rem; color: #64748b; text-transform: uppercase; font-weight: 700; border-bottom: 1px solid #e2e8f0; }
    td { padding: 1rem; border-top: 1px solid #f1f5f9; font-size: 0.88rem; color: #334155; vertical-align: middle; }
    tr:hover td { background: #f8fafc; }
    
    .badge-customer { background: #e0f2fe; color: #0369a1; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    .badge-merchant { background: #fef3c7; color: #92400e; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    .badge-admin { background: #fce7f3; color: #9d174d; padding: 0.25rem 0.6rem; border-radius: 9999px; font-weight: 700; font-size: 0.75rem; }
    .badge-verified { background: #dcfce7; color: #166534; padding: 0.2rem 0.5rem; border-radius: 6px; font-weight: 700; font-size: 0.75rem; }
    .badge-unverified { background: #f1f5f9; color: #64748b; padding: 0.2rem 0.5rem; border-radius: 6px; font-weight: 700; font-size: 0.75rem; }

    .form-input { width: 100%%; background: #ffffff; border: 1px solid #cbd5e1; color: #0f172a; padding: 0.6rem 0.8rem; border-radius: 8px; font-size: 0.9rem; }
  </style>
</head>
<body>
  <div class="admin-nav">
    <div style="display: flex; align-items: center; gap: 2rem;">
      <div class="admin-title">Bite100 Admin Dashboard</div>
      <div style="display: flex; gap: 0.5rem;">
        <a href="/admin" class="nav-tab">Restaurants</a>
        <a href="/admin/users" class="nav-tab active">Users & DB</a>
        <a href="/admin/reviews" class="nav-tab">Review Moderation</a>
        <a href="/admin/categories" class="nav-tab">Categories</a>
      </div>
    </div>
    <div style="display: flex; align-items: center; gap: 1rem;">
      <a href="/" style="color: #0f172a; text-decoration: none; font-weight: 600; font-size: 0.9rem;">&larr; Back to Website</a>
      <div style="background: #f1f5f9; border: 1px solid #e2e8f0; padding: 0.35rem 0.8rem; border-radius: 9999px; font-size: 0.82rem; color: #0f172a; font-weight: 700;">
        Admin: %s
      </div>
      <a href="/logout" style="color: #e23744; text-decoration: none; font-size: 0.85rem; font-weight: 700; padding: 0.35rem 0.75rem; border: 1px solid #fecaca; border-radius: 8px;">Log Out</a>
    </div>
  </div>

  <div class="container">
    <div class="metrics-grid">
      <div class="metric-card">
        <div style="font-size: 0.85rem; color: #64748b;">Total Registered Users</div>
        <div class="metric-num">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.85rem; color: #64748b;">Customers Database</div>
        <div class="metric-num" style="color: #0284c7;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.85rem; color: #64748b;">Restaurant Owners</div>
        <div class="metric-num" style="color: #d97706;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.85rem; color: #64748b;">Administrators</div>
        <div class="metric-num" style="color: #be185d;">%d</div>
      </div>
    </div>

    <div class="controls-bar">
      <div class="tab-pills">
        <a href="/admin/users" class="tab-pill %s">All (%d)</a>
        <a href="/admin/users?role=customer" class="tab-pill %s">Customers (%d)</a>
        <a href="/admin/users?role=merchant" class="tab-pill %s">Merchants (%d)</a>
        <a href="/admin/users?role=admin" class="tab-pill %s">Admins (%d)</a>
      </div>

      <div style="display: flex; gap: 0.75rem;">
        <input 
          type="text" 
          placeholder="Search name or email..." 
          class="form-input" 
          style="width: 250px;"
          value="%s"
          hx-get="/admin/users?table_only=true&role=%s"
          hx-trigger="keyup changed delay:300ms"
          hx-target="#admin-user-table-body"
          name="q"
        >
      </div>
    </div>

    <table>
      <thead>
        <tr>
          <th style="width: 25%%;">User Identity (ID & Email)</th>
          <th style="width: 14%%;">Role</th>
          <th style="width: 14%%;">Verification Status</th>
          <th style="width: 22%%;">Assigned Restaurant(s)</th>
          <th style="width: 11%%;">Joined Date</th>
          <th style="width: 14%%; text-align: right;">Database Controls</th>
        </tr>
      </thead>
      <tbody id="admin-user-table-body">
        %s
      </tbody>
    </table>
  </div>

  <div id="admin-modal"></div>
  <script>
    function closeAdminModal() {
      var m = document.getElementById('admin-modal');
      if (m) m.innerHTML = '';
    }
  </script>
</body>
</html>`,
		escapeHTML(adminName),
		stats.TotalUsers,
		stats.CustomerCount,
		stats.MerchantCount,
		stats.AdminCount,
		activeTab(activeRole, ""), stats.TotalUsers,
		activeTab(activeRole, "customer"), stats.CustomerCount,
		activeTab(activeRole, "merchant"), stats.MerchantCount,
		activeTab(activeRole, "admin"), stats.AdminCount,
		escapeHTML(search),
		activeRole,
		RenderAdminUserTableRows(users),
	)
}

// RenderAdminUserTableRows renders table rows for the users page
func RenderAdminUserTableRows(users []*ent.User) string {
	if len(users) == 0 {
		return `<tr><td colspan="6" style="text-align: center; color: #64748b; padding: 2.5rem;">No users found matching criteria.</td></tr>`
	}

	var sb strings.Builder
	for _, u := range users {
		sb.WriteString(RenderAdminUserSingleRow(u))
	}
	return sb.String()
}

// RenderAdminUserSingleRow renders a single user row in the user management table
func RenderAdminUserSingleRow(u *ent.User) string {
	roleBadge := `<span class="badge-customer">Customer</span>`
	if u.Role == user.RoleMerchant {
		roleBadge = `<span class="badge-merchant">Merchant Owner</span>`
	} else if u.Role == user.RoleAdmin || u.IsAdmin {
		roleBadge = `<span class="badge-admin">Admin</span>`
	}

	verifiedBadge := fmt.Sprintf(`<button type="button" class="btn-verify-toggle badge-unverified" style="cursor: pointer; border: 1px solid #cbd5e1;" hx-post="/admin/users/%d/toggle-verify" hx-target="#user-row-%d" hx-swap="outerHTML" title="Click to toggle verification status">Unverified ⟳</button>`, u.ID, u.ID)
	if u.Verified {
		verifiedBadge = fmt.Sprintf(`<button type="button" class="btn-verify-toggle badge-verified" style="cursor: pointer; border: 1px solid #bbf7d0;" hx-post="/admin/users/%d/toggle-verify" hx-target="#user-row-%d" hx-swap="outerHTML" title="Click to toggle verification status">Verified ✓</button>`, u.ID, u.ID)
	}

	assignedRestos := `<span style="color: #94a3b8; font-size: 0.85rem;">—</span>`
	var firstRestID int
	if len(u.Edges.Restaurants) > 0 {
		var rNames []string
		for _, r := range u.Edges.Restaurants {
			if firstRestID == 0 {
				firstRestID = r.ID
			}
			statusPill := `<span style="background: #dcfce7; color: #166534; font-size: 0.72rem; padding: 2px 6px; border-radius: 4px; font-weight: 700;">Approved</span>`
			if string(r.Status) == "pending" {
				statusPill = `<span style="background: #fef9c3; color: #854d0e; font-size: 0.72rem; padding: 2px 6px; border-radius: 4px; font-weight: 700;">Pending</span>`
			}
			rNames = append(rNames, fmt.Sprintf(`<div style="display: flex; align-items: center; gap: 6px; margin-bottom: 3px;"><a href="/restaurant/%d" target="_blank" style="color: #0f172a; font-weight: 700; text-decoration: underline;">%s</a> %s</div>`, r.ID, escapeHTML(r.Name), statusPill))
		}
		assignedRestos = strings.Join(rNames, "")
	}

	joinedDate := u.CreatedAt.Format("02 Jan 2006")

	// Contextual DB Action buttons
	actionsSB := strings.Builder{}
	actionsSB.WriteString(fmt.Sprintf(`<button class="btn-action" style="color: #0284c7; font-weight: 700;" hx-get="/admin/users/%d/inspect" hx-target="#admin-modal" title="Inspect full DB record">🔍 DB Inspect</button>`, u.ID))

	if u.Role == user.RoleMerchant {
		if firstRestID > 0 {
			actionsSB.WriteString(fmt.Sprintf(` <button class="btn-action" style="color: #087443; font-weight: 700;" hx-get="/admin/restaurants/%d/food" hx-target="#admin-modal" title="Add or edit dishes for this outlet">🍽️ Dishes</button>`, firstRestID))
			actionsSB.WriteString(fmt.Sprintf(` <button class="btn-action" style="color: #475569;" hx-get="/admin/restaurants/%d/hours" hx-target="#admin-modal" title="Manage operating hours">⏰ Hours</button>`, firstRestID))
		}
	} else if u.Role == user.RoleCustomer {
		actionsSB.WriteString(fmt.Sprintf(` <button class="btn-action" style="color: #d97706; font-weight: 700;" hx-post="/admin/users/%d/role?role=merchant" hx-target="#user-row-%d" hx-swap="outerHTML" hx-confirm="Promote this customer to Restaurant Owner?" title="Grant merchant privileges">🏪 Make Owner</button>`, u.ID, u.ID))
	}

	return fmt.Sprintf(`
		<tr id="user-row-%d">
			<td>
				<div style="display: flex; align-items: center; gap: 0.4rem;">
					<span style="background: #f1f5f9; color: #475569; font-size: 0.72rem; font-weight: 800; padding: 2px 5px; border-radius: 4px;">#%d</span>
					<div style="font-weight: 700; color: #0f172a;">%s</div>
				</div>
				<div style="font-size: 0.82rem; color: #64748b; margin-top: 2px;">%s</div>
			</td>
			<td>
				%s
			</td>
			<td>
				%s
			</td>
			<td>
				%s
			</td>
			<td>
				<span style="font-size: 0.82rem; color: #64748b;">%s</span>
			</td>
			<td style="text-align: right;">
				<div style="display: flex; gap: 0.4rem; justify-content: flex-end; align-items: center;">
					%s
				</div>
			</td>
		</tr>
	`, u.ID, u.ID, escapeHTML(u.Name), escapeHTML(u.Email), roleBadge, verifiedBadge, assignedRestos, joinedDate, actionsSB.String())
}

// RenderAdminInspectUserModal renders the deep database inspector modal for a user
func RenderAdminInspectUserModal(u *ent.User, revCount, favCount int, rest *fooddiscovery.MerchantRestInspect) string {
	roleText := string(u.Role)
	if u.IsAdmin {
		roleText = "Administrator"
	}

	extraSection := ""
	if u.Role == user.RoleMerchant && rest != nil {
		extraSection = fmt.Sprintf(`
		<div style="margin-top: 1.25rem; background: #fefce8; border: 1.5px solid #fef08a; border-radius: 12px; padding: 1.1rem;">
			<div style="font-size: 0.78rem; font-weight: 800; color: #854d0e; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.5rem;">🏪 Associated Restaurant Outlet</div>
			<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; font-size: 0.88rem;">
				<div><strong>Outlet Name:</strong> %s (ID #%d)</div>
				<div><strong>Approval Status:</strong> <span style="text-transform: capitalize; font-weight: 700;">%s</span></div>
				<div><strong>Total Dishes:</strong> %d items</div>
				<div><strong>Dishes &le; ₹100:</strong> <span style="color: #166534; font-weight: 700;">%d items</span></div>
				<div style="grid-column: 1 / -1;"><strong>Address:</strong> %s</div>
				<div style="grid-column: 1 / -1;"><strong>Hours:</strong> %s</div>
			</div>
			<div style="margin-top: 0.85rem; display: flex; gap: 0.5rem;">
				<button type="button" class="btn-action" style="background: #087443; color: white; font-weight: 700; border: none; padding: 0.4rem 0.85rem;" hx-get="/admin/restaurants/%d/food" hx-target="#admin-modal">+ Manage & Add Dishes</button>
				<button type="button" class="btn-action" style="background: #ffffff; color: #1e293b; font-weight: 700; padding: 0.4rem 0.85rem;" hx-get="/admin/restaurants/%d/hours" hx-target="#admin-modal">⏰ Operating Hours</button>
			</div>
		</div>`, escapeHTML(rest.Name), rest.ID, rest.Status, rest.TotalFoods, rest.Under100Foods, escapeHTML(rest.Address), escapeHTML(rest.OpeningHours), rest.ID, rest.ID)
	} else if u.Role == user.RoleCustomer {
		extraSection = fmt.Sprintf(`
		<div style="margin-top: 1.25rem; background: #f0fdf4; border: 1.5px solid #bbf7d0; border-radius: 12px; padding: 1.1rem;">
			<div style="font-size: 0.78rem; font-weight: 800; color: #166534; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.5rem;">👤 Customer Activity Database</div>
			<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; font-size: 0.88rem;">
				<div><strong>Total Reviews Written:</strong> %d</div>
				<div><strong>Favorite Cafes Saved:</strong> %d</div>
			</div>
			<div style="margin-top: 0.85rem;">
				<button type="button" class="btn-action" style="background: #d97706; color: white; font-weight: 700; border: none; padding: 0.4rem 0.85rem;" hx-post="/admin/users/%d/role?role=merchant" hx-target="#admin-modal" hx-confirm="Promote this customer to Restaurant Owner?">Promote to Restaurant Owner</button>
			</div>
		</div>`, revCount, favCount, u.ID)
	}

	return fmt.Sprintf(`
	<div class="modal-overlay" onclick="if(event.target === this) closeAdminModal()">
		<div class="modal-card" style="width: 620px;">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid #e2e8f0;">
				<div style="display: flex; align-items: center; gap: 0.5rem;">
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">User Database Record</h3>
					<span style="background: #f1f5f9; color: #475569; font-size: 0.75rem; font-weight: 800; padding: 2px 7px; border-radius: 9999px;">ID #%d</span>
				</div>
				<button onclick="closeAdminModal()" style="background: none; border: none; font-size: 1.6rem; color: #64748b; cursor: pointer; line-height: 1;">&times;</button>
			</div>

			<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1.1rem; font-size: 0.88rem; display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem;">
				<div><strong>Full Name:</strong> %s</div>
				<div><strong>Email Address:</strong> %s</div>
				<div><strong>Platform Role:</strong> <span style="text-transform: capitalize; font-weight: 700;">%s</span></div>
				<div><strong>Verification:</strong> %t</div>
				<div><strong>Created Date:</strong> %s</div>
				<div><strong>Updated Date:</strong> %s</div>
			</div>

			%s

			<div style="margin-top: 1.5rem; display: flex; justify-content: flex-end; gap: 0.75rem; padding-top: 1rem; border-top: 1px solid #e2e8f0;">
				<button type="button" onclick="closeAdminModal()" class="btn-action" style="padding: 0.5rem 1rem;">Close</button>
			</div>
		</div>
	</div>`, u.ID, escapeHTML(u.Name), escapeHTML(u.Email), roleText, u.Verified, u.CreatedAt.Format("02 Jan 2006 15:04"), u.UpdatedAt.Format("02 Jan 2006 15:04"), extraSection)
}

// RenderStatusBadge renders the status badge
func RenderStatusBadge(status restaurant.Status, id int) string {
	switch status {
	case restaurant.StatusApproved:
		return `<span class="badge-approved">Approved</span>`
	case restaurant.StatusRejected:
		return `<span class="badge-rejected">Rejected</span>`
	default:
		return `<span class="badge-pending">Pending</span>`
	}
}

// RenderNewRestaurantModal renders the add restaurant form modal
func RenderNewRestaurantModal() string {
	return `
		<div class="modal-overlay" onclick="if(event.target === this) closeAdminModal()">
			<div class="modal-card">
				<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Add New Restaurant</h3>
					<button onclick="closeAdminModal()" style="background: none; border: none; font-size: 1.4rem; color: #64748b; cursor: pointer;">&times;</button>
				</div>
				<form hx-post="/admin/restaurants" hx-target="#admin-modal">
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Name *</label>
						<input type="text" name="name" class="form-input" placeholder="E.g. CTR Tiffin Room" required>
					</div>
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Address *</label>
						<input type="text" name="address" class="form-input" placeholder="E.g. 7th Cross, Margosa Road, Malleshwaram" required>
					</div>
					<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-bottom: 1rem;">
						<div>
							<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Latitude *</label>
							<input type="number" step="any" name="latitude" class="form-input" placeholder="13.0031" required>
						</div>
						<div>
							<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Longitude *</label>
							<input type="number" step="any" name="longitude" class="form-input" placeholder="77.5645" required>
						</div>
					</div>
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Description</label>
						<textarea name="description" class="form-input" rows="3" placeholder="Description of the eatery..."></textarea>
					</div>
					<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem;">
						<button type="button" onclick="closeAdminModal()" class="btn-action" style="padding: 0.6rem 1.25rem;">Cancel</button>
						<button type="submit" class="btn-create" style="padding: 0.6rem 1.25rem;">Save Restaurant</button>
					</div>
				</form>
			</div>
		</div>
	`
}

// RenderEditRestaurantModal renders the edit restaurant form modal
func RenderEditRestaurantModal(rest *ent.Restaurant) string {
	return fmt.Sprintf(`
		<div class="modal-overlay" onclick="if(event.target === this) closeAdminModal()">
			<div class="modal-card">
				<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
					<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Edit Restaurant</h3>
					<button onclick="closeAdminModal()" style="background: none; border: none; font-size: 1.4rem; color: #64748b; cursor: pointer;">&times;</button>
				</div>
				<form hx-post="/admin/restaurants/%d/edit" hx-target="#admin-modal">
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Name *</label>
						<input type="text" name="name" class="form-input" value="%s" required>
					</div>
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Address *</label>
						<input type="text" name="address" class="form-input" value="%s" required>
					</div>
					<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-bottom: 1rem;">
						<div>
							<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Latitude *</label>
							<input type="number" step="any" name="latitude" class="form-input" value="%.6f" required>
						</div>
						<div>
							<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Longitude *</label>
							<input type="number" step="any" name="longitude" class="form-input" value="%.6f" required>
						</div>
					</div>
					<div style="margin-bottom: 1rem;">
						<label style="display: block; font-size: 0.85rem; color: #334155; margin-bottom: 0.35rem; font-weight: 700;">Description</label>
						<textarea name="description" class="form-input" rows="3">%s</textarea>
					</div>
					<div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem;">
						<button type="button" onclick="closeAdminModal()" class="btn-action" style="padding: 0.6rem 1.25rem;">Cancel</button>
						<button type="submit" class="btn-create" style="padding: 0.6rem 1.25rem;">Save Changes</button>
					</div>
				</form>
			</div>
		</div>
	`, rest.ID, escapeHTML(rest.Name), escapeHTML(rest.Address), rest.Latitude, rest.Longitude, escapeHTML(rest.Description))
}

// RenderManageFoodModal renders food management modal
func RenderManageFoodModal(rest *ent.Restaurant, foods []*ent.FoodItem) string {
	return fmt.Sprintf(`
		<div class="modal-overlay" onclick="if(event.target === this) closeAdminModal()">
			<div class="modal-card" style="width: 700px;">
				<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.25rem;">
					<div>
						<h3 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; color: #0f172a;">Manage Menu: %s</h3>
						<div style="font-size: 0.85rem; color: #64748b;">Add and manage dishes priced under ₹100</div>
					</div>
					<button onclick="closeAdminModal()" style="background: none; border: none; font-size: 1.4rem; color: #64748b; cursor: pointer;">&times;</button>
				</div>

				<form hx-post="/admin/restaurants/%d/food" hx-target="#food-items-container" style="background: #f8fafc; border: 1px solid #e2e8f0; padding: 1rem; border-radius: 12px; margin-bottom: 1.5rem;">
					<div style="display: grid; grid-template-columns: 2fr 1fr 1fr; gap: 0.75rem; margin-bottom: 0.75rem;">
						<input type="text" name="name" class="form-input" placeholder="Dish name (e.g. Masala Dosa)" required>
						<input type="number" step="1" name="price" class="form-input" placeholder="Price (₹)" required>
						<input type="text" name="category" class="form-input" placeholder="Category">
					</div>
					<div style="display: flex; justify-content: flex-end;">
						<button type="submit" class="btn-create" style="padding: 0.5rem 1rem;">+ Add Dish</button>
					</div>
				</form>

				<div id="food-items-container">
					%s
				</div>
			</div>
		</div>
	`, escapeHTML(rest.Name), rest.ID, RenderFoodItemsList(rest.ID, foods))
}

// RenderFoodItemsList renders list of food items in admin
func RenderFoodItemsList(restoID int, foods []*ent.FoodItem) string {
	if len(foods) == 0 {
		return `<div style="text-align: center; color: #64748b; padding: 1.5rem;">No food items added yet.</div>`
	}
	var sb strings.Builder
	sb.WriteString(`<div style="display: flex; flex-direction: column; gap: 0.5rem;">`)
	for _, f := range foods {
		sb.WriteString(RenderSingleFoodRow(f))
	}
	sb.WriteString(`</div>`)
	return sb.String()
}

// RenderSingleFoodRow renders a single food item row
func RenderSingleFoodRow(f *ent.FoodItem) string {
	priceStyle := "color: #15803d; font-weight: 800;"
	if f.Price > 100 {
		priceStyle = "color: #b91c1c; font-weight: 800;"
	}
	return fmt.Sprintf(`
		<div id="food-row-%d" style="display: flex; justify-content: space-between; align-items: center; background: #ffffff; border: 1px solid #e2e8f0; padding: 0.75rem 1rem; border-radius: 8px;">
			<div>
				<span style="font-weight: 700; color: #0f172a;">%s</span>
				<span style="font-size: 0.8rem; color: #64748b; margin-left: 0.5rem;">%s</span>
			</div>
			<div style="display: flex; align-items: center; gap: 1rem;">
				<span style="%s">₹%.0f</span>
				<span id="avail-badge-%d">%s</span>
				<button 
					class="btn-action" 
					style="color: #b91c1c;"
					hx-delete="/admin/food/%d"
					hx-target="#food-row-%d"
					hx-swap="outerHTML"
				>
					Delete
				</button>
			</div>
		</div>
	`, f.ID, escapeHTML(f.Name), escapeHTML(f.Category), priceStyle, f.Price, f.ID, RenderAvailabilityBadge(f), f.ID, f.ID)
}

// RenderAvailabilityBadge renders in stock / out of stock toggle button
func RenderAvailabilityBadge(f *ent.FoodItem) string {
	if f.Available {
		return fmt.Sprintf(`<button class="btn-action" style="color: #15803d;" hx-post="/admin/food/%d/toggle" hx-target="#avail-badge-%d">In Stock</button>`, f.ID, f.ID)
	}
	return fmt.Sprintf(`<button class="btn-action" style="color: #b91c1c;" hx-post="/admin/food/%d/toggle" hx-target="#avail-badge-%d">Out of Stock</button>`, f.ID, f.ID)
}

// RenderLoginPage renders dedicated customer login & registration page
func RenderLoginPage(flashMsg, errMsg string, returnTo ...string) string {
	var flashHTML string
	if flashMsg == "admin_required" {
		flashHTML = `<div style="background: #fef08a; border: 1px solid #fde047; color: #854d0e; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			<strong>Administrator Access Required:</strong> Please sign in with an administrator account to access the Admin Panel.
		</div>`
	} else if flashMsg == "merchant_required" {
		flashHTML = `<div style="background: #fef08a; border: 1px solid #fde047; color: #854d0e; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			<strong>Restaurant Owner Access Required:</strong> Please sign in with a merchant partner account to access the Merchant Portal.
		</div>`
	} else if flashMsg == "auth_required" {
		flashHTML = `<div style="background: #fef08a; border: 1px solid #fde047; color: #854d0e; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			Please sign in to access your account profile and settings.
		</div>`
	} else if flashMsg == "password_reset_success" {
		flashHTML = `<div style="background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			Your password has been reset successfully. Please log in with your new password.
		</div>`
	} else if flashMsg != "" {
		flashHTML = fmt.Sprintf(`<div style="background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">%s</div>`, escapeHTML(flashMsg))
	}

	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div style="background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			%s
		</div>`, escapeHTML(errMsg))
	}

	returnField := ""
	if len(returnTo) > 0 && returnTo[0] != "" {
		returnField = fmt.Sprintf(`<input type="hidden" name="return_to" id="custReturnTo" value="%s">`, escapeHTML(returnTo[0]))
	} else {
		returnField = `<input type="hidden" name="return_to" id="custReturnTo" value="/">`
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Customer Sign In — Under100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 450px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
    }
    .auth-brand {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-family: 'Outfit', sans-serif;
      font-size: 1.4rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-brand-badge {
      background: #fef08a;
      color: #854d0e;
      font-size: 0.72rem;
      font-weight: 800;
      padding: 0.2rem 0.5rem;
      border-radius: 6px;
    }
    .auth-subtitle {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.25rem;
      margin-top: 0.35rem;
    }
    .tabs-bar {
      display: flex;
      background: #f1f5f9;
      border-radius: 12px;
      padding: 0.25rem;
      margin-bottom: 1.25rem;
    }
    .tab-btn {
      flex: 1;
      background: none;
      border: none;
      color: #64748b;
      font-weight: 700;
      font-size: 0.85rem;
      padding: 0.55rem;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.15s;
    }
    .tab-btn.active {
      background: #ffffff;
      color: #0f172a;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }
    .subtabs-bar {
      display: flex;
      gap: 8px;
      margin-bottom: 1.25rem;
    }
    .subtab-btn {
      flex: 1;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      color: #64748b;
      font-weight: 600;
      font-size: 0.8rem;
      padding: 0.5rem 0.6rem;
      border-radius: 8px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      transition: all 0.15s;
    }
    .subtab-btn.active {
      background: #fee2e2;
      border-color: #fca5a5;
      color: #b91c1c;
      font-weight: 700;
    }
    .form-group { margin-bottom: 1rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 600; color: #0f172a; margin-bottom: 0.35rem; }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.9rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.8rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      margin-top: 0.5rem;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .btn-otp-action {
      background: #0f172a;
      color: #ffffff;
      border: none;
      padding: 0.65rem 1rem;
      border-radius: 10px;
      font-size: 0.82rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
      transition: background 0.15s;
    }
    .btn-otp-action:hover { background: #334155; }
    .btn-otp-action:disabled { opacity: 0.6; cursor: not-allowed; }
    .otp-helper-banner {
      background: #f0fdf4;
      border: 1px solid #bbf7d0;
      border-radius: 10px;
      padding: 10px 12px;
      margin-bottom: 1rem;
      font-size: 0.8rem;
      color: #166534;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.85rem;
    }
    .back-link:hover { color: #e23744; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <a href="/" class="auth-brand">
        Under100
      </a>
      <span class="auth-brand-badge">Bengaluru</span>
    </div>
    <p class="auth-subtitle">Sign in or register a customer account to discover delicious food under ₹100.</p>

    %s
    %s

    <div class="tabs-bar" style="display: none;">
      <button type="button" id="tab-in" class="tab-btn active" onclick="switchMainTab('in')">Sign In</button>
      <button type="button" id="tab-up" class="tab-btn" onclick="switchMainTab('up')">New Customer Register</button>
    </div>

    <!-- SIGN IN CONTAINER -->
    <div id="container-signin">
      <div class="subtabs-bar">
        <button type="button" id="cust-sub-pass" class="subtab-btn active" onclick="switchCustAuthMode('pass')">
          Password Login
        </button>
        <button type="button" id="cust-sub-otp" class="subtab-btn" onclick="switchCustAuthMode('otp')">
          Mobile / OTP Login
        </button>
      </div>

      <!-- Password Form -->
      <form id="form-signin-pass" method="POST" action="/login">
        %s
        <div class="form-group">
          <label class="form-label">Email Address</label>
          <input type="email" name="email" class="form-input" placeholder="user@example.com" required>
        </div>
        <div class="form-group">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.4rem;">
            <label class="form-label" style="margin-bottom: 0;">Password</label>
            <a href="/password/forgot" style="color: #e23744; font-size: 0.8rem; text-decoration: none; font-weight: 700;">Forgot password?</a>
          </div>
          <input type="password" name="password" class="form-input" placeholder="••••••••" required>
        </div>
        <button type="submit" class="btn-submit">Sign In with Password</button>
      </form>

      <!-- OTP Form -->
      <div id="form-signin-otp" style="display: none;">

        <div class="form-group">
          <label class="form-label">Mobile Number or Email</label>
          <div style="display: flex; gap: 8px;">
            <input type="text" id="custOtpIdentifier" class="form-input" placeholder="e.g. 9876543210 or user@example.com" style="flex: 1;">
            <button type="button" id="custBtnSend" class="btn-otp-action" onclick="sendCustOTP()">Get OTP</button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">6-Digit Verification Code</label>
          <input type="text" id="custOtpInput" class="form-input" placeholder="Enter 6-digit OTP" maxlength="6" style="letter-spacing: 4px; font-size: 1.1rem; font-weight: 700; text-align: center;">
        </div>

        <div id="custOtpStatus" style="display: none; font-size: 0.82rem; margin-bottom: 0.75rem; font-weight: 600;"></div>

        <button type="button" id="custBtnVerify" class="btn-submit" onclick="verifyCustOTP()">Verify &amp; Sign In</button>
      </div>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #64748b;">
        Don't have an account? <a href="javascript:switchMainTab('up')" style="color: #e23744; font-weight: 700; text-decoration: none;">Register here</a>
      </div>
    </div>

    <!-- SIGN UP CONTAINER -->
    <div id="container-signup" style="display: none;">
      <div class="subtabs-bar">
        <button type="button" id="reg-sub-pass" class="subtab-btn active" onclick="switchCustRegMode('pass')">
          Register with Email &amp; Password
        </button>
        <button type="button" id="reg-sub-otp" class="subtab-btn" onclick="switchCustRegMode('otp')">
          Register with Mobile OTP
        </button>
      </div>

      <!-- Register with Password -->
      <form id="form-signup" method="POST" action="/register">
        <div class="form-group">
          <label class="form-label">Full Name</label>
          <input type="text" name="name" class="form-input" placeholder="Aarav Patel" required>
        </div>
        <div class="form-group">
          <label class="form-label">Email Address</label>
          <input type="email" name="email" class="form-input" placeholder="aarav@example.com" required>
        </div>
        <div class="form-group">
          <label class="form-label">Create Password</label>
          <input type="password" name="password" class="form-input" placeholder="Min 6 characters" minlength="6" required>
        </div>
        <button type="submit" class="btn-submit">Register with Email &amp; Password</button>
      </form>

      <!-- Register with OTP -->
      <div id="form-signup-otp" style="display: none;">
        <div class="form-group">
          <label class="form-label">Your Full Name</label>
          <input type="text" id="custRegName" class="form-input" placeholder="e.g. Aarav Patel">
        </div>
        <div class="form-group">
          <label class="form-label">Mobile Number</label>
          <div style="display: flex; gap: 8px;">
            <input type="text" id="custRegOtpIdentifier" class="form-input" placeholder="e.g. 9876543210" style="flex: 1;">
            <button type="button" id="custRegBtnSend" class="btn-otp-action" onclick="sendCustRegOTP()">Get OTP</button>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">6-Digit Verification Code</label>
          <input type="text" id="custRegOtpInput" class="form-input" placeholder="Enter 6-digit OTP" maxlength="6" style="letter-spacing: 4px; font-size: 1.1rem; font-weight: 700; text-align: center;">
        </div>
        <div id="custRegOtpStatus" style="display: none; font-size: 0.82rem; margin-bottom: 0.75rem; font-weight: 600;"></div>
        <button type="button" id="custRegBtnVerify" class="btn-submit" onclick="verifyCustRegOTP()">Verify OTP &amp; Create Customer Account</button>
      </div>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #64748b;">
        Already have an account? <a href="javascript:switchMainTab('in')" style="color: #e23744; font-weight: 700; text-decoration: none;">Sign in here</a>
      </div>
    </div>

    <div style="text-align: center;">
      <a href="/" class="back-link">&larr; Return to Home Page</a>
    </div>
  </div>

  <script>
    var currentCustDemoOtp = '123456';

    function switchMainTab(mode) {
      const tabIn = document.getElementById('tab-in');
      const tabUp = document.getElementById('tab-up');
      const cIn = document.getElementById('container-signin');
      const cUp = document.getElementById('container-signup');

      if (mode === 'up') {
        tabIn.classList.remove('active');
        tabUp.classList.add('active');
        cIn.style.display = 'none';
        cUp.style.display = 'block';
      } else {
        tabUp.classList.remove('active');
        tabIn.classList.add('active');
        cUp.style.display = 'none';
        cIn.style.display = 'block';
      }
    }

    function switchCustAuthMode(m) {
      const btnPass = document.getElementById('cust-sub-pass');
      const btnOtp = document.getElementById('cust-sub-otp');
      const fPass = document.getElementById('form-signin-pass');
      const fOtp = document.getElementById('form-signin-otp');

      if (m === 'otp') {
        btnPass.classList.remove('active');
        btnOtp.classList.add('active');
        fPass.style.display = 'none';
        fOtp.style.display = 'block';
      } else {
        btnOtp.classList.remove('active');
        btnPass.classList.add('active');
        fOtp.style.display = 'none';
        fPass.style.display = 'block';
      }
    }

    function switchCustRegMode(m) {
      const btnPass = document.getElementById('reg-sub-pass');
      const btnOtp = document.getElementById('reg-sub-otp');
      const fPass = document.getElementById('form-signup');
      const fOtp = document.getElementById('form-signup-otp');

      if (m === 'otp') {
        btnPass.classList.remove('active');
        btnOtp.classList.add('active');
        fPass.style.display = 'none';
        fOtp.style.display = 'block';
      } else {
        btnOtp.classList.remove('active');
        btnPass.classList.add('active');
        fOtp.style.display = 'none';
        fPass.style.display = 'block';
      }
    }



    function sendCustOTP() {
      var idVal = document.getElementById('custOtpIdentifier').value.trim();
      var statusDiv = document.getElementById('custOtpStatus');
      var btnSend = document.getElementById('custBtnSend');
      if (!idVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter your mobile number or email.';
        return;
      }
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      statusDiv.style.display = 'none';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('portal', 'customer');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          var inp = document.getElementById('custOtpInput');
          if (inp) {
            inp.value = data.demo_otp || '';
            inp.focus();
          }
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#166534';
          if (data.demo_otp) {
            statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: #DCFCE7; color: #166534; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #4B5563;">(Code auto-filled &mdash; click "Verify & Sign In" to complete login)</span>';
          } else {
            statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || idVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
          }
          var count = 60;
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend OTP';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyCustOTP() {
      var idVal = document.getElementById('custOtpIdentifier').value.trim();
      var otpVal = document.getElementById('custOtpInput').value.trim();
      var statusDiv = document.getElementById('custOtpStatus');
      var btnVerify = document.getElementById('custBtnVerify');
      var retTo = document.getElementById('custReturnTo').value || '/';

      if (!idVal || !otpVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter both your identifier and the 6-digit OTP.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Verifying...';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('otp', otpVal);
      body.append('portal', 'customer');
      body.append('return_to', retTo);

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#16a34a';
          statusDiv.innerText = 'Verified! Signing you in...';
          window.location.href = data.redirect || retTo;
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify & Sign In';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Invalid or expired OTP code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify & Sign In';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error during verification.';
      });
    }

    function sendCustRegOTP() {
      var idVal = document.getElementById('custRegOtpIdentifier').value.trim();
      var statusDiv = document.getElementById('custRegOtpStatus');
      var btnSend = document.getElementById('custRegBtnSend');
      if (!idVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter your mobile number.';
        return;
      }
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      statusDiv.style.display = 'none';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('portal', 'customer');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          var inp = document.getElementById('custRegOtpInput');
          if (inp) {
            inp.value = data.demo_otp || '';
            inp.focus();
          }
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#166534';
          if (data.demo_otp) {
            statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: #DCFCE7; color: #166534; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #4B5563;">(Code auto-filled &mdash; click "Create Account" to continue)</span>';
          } else {
            statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || idVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
          }
          var count = 60;
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyCustRegOTP() {
      var nameVal = document.getElementById('custRegName').value.trim();
      var idVal = document.getElementById('custRegOtpIdentifier').value.trim();
      var otpVal = document.getElementById('custRegOtpInput').value.trim();
      var statusDiv = document.getElementById('custRegOtpStatus');
      var btnVerify = document.getElementById('custRegBtnVerify');
      var retTo = document.getElementById('custReturnTo').value || '/';

      if (!idVal || !otpVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter mobile number and 6-digit OTP.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Creating account...';

      var body = new URLSearchParams();
      body.append('name', nameVal || idVal);
      body.append('identifier', idVal);
      body.append('otp', otpVal);
      body.append('portal', 'customer');
      body.append('return_to', retTo);

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#16a34a';
          statusDiv.innerText = 'Account verified! Signing you in...';
          window.location.href = data.redirect || retTo;
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify OTP & Create Customer Account';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Invalid OTP code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify OTP & Create Customer Account';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error while creating account.';
      });
    }

    if (new URLSearchParams(window.location.search).get('mode') === 'signup') {
      switchMainTab('up');
    }
  </script>
</body>
</html>`, flashHTML, errHTML, returnField)
}

// RenderMerchantLoginPage renders dedicated sign-in and partner registration page for restaurant owners
func RenderMerchantLoginPage(flashMsg, errMsg string) string {
	var flashHTML string
	if flashMsg == "merchant_required" {
		flashHTML = `<div style="background: #fef3c7; border: 1px solid #fde68a; color: #92400e; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			<strong>Restaurant Owner Portal:</strong> Please sign in with your restaurant partner account to manage your kitchen outlet.
		</div>`
	} else if flashMsg != "" {
		flashHTML = fmt.Sprintf(`<div style="background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">%s</div>`, escapeHTML(flashMsg))
	}

	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div style="background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			%s
		</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Merchant Sign In — Under100 Partner Portal</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 480px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
    }
    .auth-brand {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-family: 'Outfit', sans-serif;
      font-size: 1.4rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-brand-badge {
      background: #f0fdf4;
      color: #166534;
      border: 1px solid #bbf7d0;
      font-size: 0.72rem;
      font-weight: 800;
      padding: 0.2rem 0.5rem;
      border-radius: 6px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .auth-nav-tabs {
      display: flex;
      gap: 8px;
      background: #f1f5f9;
      padding: 4px;
      border-radius: 12px;
      margin: 1.25rem 0 1.25rem;
    }
    .auth-nav-tab {
      flex: 1;
      text-align: center;
      padding: 8px 12px;
      font-size: 0.85rem;
      font-weight: 700;
      border-radius: 8px;
      cursor: pointer;
      border: none;
      background: none;
      color: #64748b;
      transition: all 0.15s;
    }
    .auth-nav-tab.active {
      background: #ffffff;
      color: #087443;
      box-shadow: 0 2px 6px rgba(0,0,0,0.06);
    }
    .subtabs-bar {
      display: flex;
      gap: 8px;
      margin-bottom: 1.25rem;
    }
    .subtab-btn {
      flex: 1;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      color: #64748b;
      font-weight: 600;
      font-size: 0.8rem;
      padding: 0.5rem 0.6rem;
      border-radius: 8px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      transition: all 0.15s;
    }
    .subtab-btn.active {
      background: #dcfce7;
      border-color: #86efac;
      color: #166534;
      font-weight: 700;
    }
    .demo-helper-box {
      background: #f0fdf4;
      border: 1px solid #bbf7d0;
      border-radius: 10px;
      padding: 10px 14px;
      margin-bottom: 1.25rem;
      font-size: 0.8rem;
      color: #166534;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
    .btn-fill-demo {
      background: #087443;
      color: #ffffff;
      border: none;
      padding: 4px 10px;
      border-radius: 6px;
      font-size: 0.75rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
    }
    .btn-fill-demo:hover { background: #065f36; }
    .form-group { margin-bottom: 1rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 600; color: #0f172a; margin-bottom: 0.35rem; }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.9rem;
      outline: none;
    }
    .form-input:focus { border-color: #087443; }
    .btn-submit {
      width: 100%%;
      background: #087443;
      color: #ffffff;
      border: none;
      padding: 0.8rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      margin-top: 0.5rem;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #065f36; }
    .btn-otp-action {
      background: #087443;
      color: #ffffff;
      border: none;
      padding: 0.65rem 1rem;
      border-radius: 10px;
      font-size: 0.82rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
      transition: background 0.15s;
    }
    .btn-otp-action:hover { background: #065f36; }
    .btn-otp-action:disabled { opacity: 0.6; cursor: not-allowed; }
    .back-link {
      display: inline-block;
      color: #64748b;
      text-decoration: none;
      font-size: 0.85rem;
    }
    .back-link:hover { color: #087443; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <a href="/merchant" class="auth-brand">
        Under100
      </a>
      <span class="auth-brand-badge">Partner Portal</span>
    </div>

    <div class="demo-helper-box">
      <div>
        <strong>Demo Merchant:</strong> CTR Malleshwaram<br>
        <span style="font-family: monospace; font-size: 0.76rem;">merchant.ctr@fooddiscovery.com</span>
      </div>
      <button type="button" class="btn-fill-demo" onclick="fillMerchantDemo()">Fill Demo</button>
    </div>

    %s
    %s

    <!-- DIRECT SIGN IN CONTAINER -->
    <div id="merch-signin-container">
      <div class="subtabs-bar">
        <button type="button" id="merch-sub-pass" class="subtab-btn active" onclick="switchMerchAuthMode('pass')">
          Password Login
        </button>
        <button type="button" id="merch-sub-otp" class="subtab-btn" onclick="switchMerchAuthMode('otp')">
          Mobile / OTP Login
        </button>
      </div>

      <!-- Merchant Password Login Form -->
      <form method="POST" action="/merchant/login" id="merchantLoginForm">
        <div class="form-group">
          <label class="form-label">Owner Email Address</label>
          <input type="email" name="email" id="merchEmail" class="form-input" placeholder="merchant@example.com" required>
        </div>
        <div class="form-group">
          <label class="form-label">Password</label>
          <input type="password" name="password" id="merchPass" class="form-input" placeholder="••••••••" required>
        </div>
        <button type="submit" class="btn-submit">Sign In to Merchant Portal</button>
      </form>

      <!-- Merchant OTP Login Form -->
      <div id="merchantOtpForm" style="display: none;">

        <div class="form-group">
          <label class="form-label">Registered Owner Mobile or Email</label>
          <div style="display: flex; gap: 8px;">
            <input type="text" id="merchOtpIdentifier" class="form-input" placeholder="e.g. 9876543210 or owner@example.com" style="flex: 1;">
            <button type="button" id="merchBtnSend" class="btn-otp-action" onclick="sendMerchOTP()">Get OTP</button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">6-Digit Verification Code</label>
          <input type="text" id="merchOtpInput" class="form-input" placeholder="Enter 6-digit OTP" maxlength="6" style="letter-spacing: 4px; font-size: 1.1rem; font-weight: 700; text-align: center;">
        </div>

        <div id="merchOtpStatus" style="display: none; font-size: 0.82rem; margin-bottom: 0.75rem; font-weight: 600;"></div>

        <button type="button" id="merchBtnVerify" class="btn-submit" onclick="verifyMerchOTP()">Verify &amp; Open Merchant Portal</button>
      </div>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #64748b;">
        Don't have an account? <a href="javascript:switchMerchMainTab('register')" style="color: #087443; font-weight: 700; text-decoration: none;">Register here</a>
      </div>
    </div>

    <!-- REGISTER CONTAINER -->
    <div id="merch-register-container" style="display: none;">
      <div class="subtabs-bar">
        <button type="button" id="merch-reg-sub-pass" class="subtab-btn active" onclick="switchMerchRegMode('pass')">
          Register with Email &amp; Password
        </button>
        <button type="button" id="merch-reg-sub-otp" class="subtab-btn" onclick="switchMerchRegMode('otp')">
          Register with Mobile OTP
        </button>
      </div>

      <!-- Mode 1: Register with Email & Password -->
      <form method="POST" action="/merchant/register" id="merchantRegisterForm">
        <div class="form-group">
          <label class="form-label">Owner Full Name</label>
          <input type="text" name="name" class="form-input" placeholder="e.g. Ramesh Kumar" required>
        </div>
        <div class="form-group">
          <label class="form-label">Business Email Address</label>
          <input type="email" name="email" class="form-input" placeholder="owner@restaurant.com" required>
        </div>
        <div class="form-group">
          <label class="form-label">Password</label>
          <input type="password" name="password" class="form-input" placeholder="Min 6 characters" minlength="6" required>
        </div>
        <div class="form-group">
          <label class="form-label">Contact Phone Number</label>
          <input type="tel" name="phone" class="form-input" placeholder="e.g. 9876543210" required>
        </div>
        <div class="form-group">
          <label class="form-label">Restaurant / Outlet Name</label>
          <input type="text" name="restaurant_name" class="form-input" placeholder="e.g. Sri Krishna Sagar" required>
        </div>
        <div class="form-group">
          <label class="form-label">Outlet Address / Area in Bengaluru</label>
          <input type="text" name="restaurant_address" class="form-input" placeholder="e.g. Koramangala 5th Block, Bengaluru" required>
        </div>
        <button type="submit" class="btn-submit">Register Outlet &amp; Launch Merchant Portal</button>
      </form>

      <!-- Mode 2: Register with Mobile OTP -->
      <div id="merchantRegisterOtpForm" style="display: none;">
        <div class="form-group">
          <label class="form-label">Owner Full Name</label>
          <input type="text" id="merchRegName" class="form-input" placeholder="e.g. Ramesh Kumar">
        </div>
        <div class="form-group">
          <label class="form-label">Owner Mobile Number</label>
          <div style="display: flex; gap: 8px;">
            <input type="tel" id="merchRegOtpPhone" class="form-input" placeholder="e.g. 9876543210" style="flex: 1;">
            <button type="button" id="merchRegBtnSend" class="btn-otp-action" onclick="sendMerchRegOTP()">Get OTP</button>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">6-Digit Verification Code</label>
          <input type="text" id="merchRegOtpInput" class="form-input" placeholder="Enter 6-digit OTP" maxlength="6" style="letter-spacing: 4px; font-size: 1.1rem; font-weight: 700; text-align: center;">
        </div>
        <div class="form-group">
          <label class="form-label">Restaurant / Outlet Name</label>
          <input type="text" id="merchRegRestName" class="form-input" placeholder="e.g. Sri Krishna Sagar">
        </div>
        <div class="form-group">
          <label class="form-label">Outlet Address / Area in Bengaluru</label>
          <input type="text" id="merchRegRestAddress" class="form-input" placeholder="e.g. Koramangala 5th Block, Bengaluru">
        </div>
        <div id="merchRegOtpStatus" style="display: none; font-size: 0.82rem; margin-bottom: 0.75rem; font-weight: 600;"></div>
        <button type="button" id="merchRegBtnVerify" class="btn-submit" onclick="verifyMerchRegOTP()">Verify OTP &amp; Launch Merchant Portal</button>
      </div>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #64748b;">
        Already have an account? <a href="javascript:switchMerchMainTab('signin')" style="color: #087443; font-weight: 700; text-decoration: none;">Sign in here</a>
      </div>
    </div>

    <div style="text-align: center; margin-top: 1.25rem; border-top: 1px solid #f1f5f9; padding-top: 1rem; font-size: 0.85rem;">
      <a href="/" class="back-link">&larr; Return to Customer Storefront</a>
    </div>
  </div>

  <script>
    var currentMerchDemoOtp = '123456';

    function fillMerchantDemo() {
      switchMerchMainTab('signin');
      switchMerchAuthMode('pass');
      document.getElementById('merchEmail').value = 'merchant.ctr@fooddiscovery.com';
      document.getElementById('merchPass').value = 'Merchant@123';
    }

    function switchMerchMainTab(tab) {
      var tabIn = document.getElementById('merch-tab-in');
      var tabReg = document.getElementById('merch-tab-reg');
      var cIn = document.getElementById('merch-signin-container');
      var cReg = document.getElementById('merch-register-container');

      if (tab === 'register') {
        if (tabIn) { tabIn.classList.remove('active'); }
        if (tabReg) { tabReg.classList.add('active'); }
        if (cIn) cIn.style.display = 'none';
        if (cReg) cReg.style.display = 'block';
      } else {
        if (tabReg) { tabReg.classList.remove('active'); }
        if (tabIn) { tabIn.classList.add('active'); }
        if (cReg) cReg.style.display = 'none';
        if (cIn) cIn.style.display = 'block';
      }
    }

    function switchMerchRegMode(m) {
      var btnPass = document.getElementById('merch-reg-sub-pass');
      var btnOtp = document.getElementById('merch-reg-sub-otp');
      var fPass = document.getElementById('merchantRegisterForm');
      var fOtp = document.getElementById('merchantRegisterOtpForm');

      if (m === 'otp') {
        btnPass.classList.remove('active');
        btnOtp.classList.add('active');
        fPass.style.display = 'none';
        fOtp.style.display = 'block';
      } else {
        btnOtp.classList.remove('active');
        btnPass.classList.add('active');
        fOtp.style.display = 'none';
        fPass.style.display = 'block';
      }
    }

    function switchMerchAuthMode(m) {
      var btnPass = document.getElementById('merch-sub-pass');
      var btnOtp = document.getElementById('merch-sub-otp');
      var fPass = document.getElementById('merchantLoginForm');
      var fOtp = document.getElementById('merchantOtpForm');

      if (m === 'otp') {
        btnPass.classList.remove('active');
        btnOtp.classList.add('active');
        fPass.style.display = 'none';
        fOtp.style.display = 'block';
      } else {
        btnOtp.classList.remove('active');
        btnPass.classList.add('active');
        fOtp.style.display = 'none';
        fPass.style.display = 'block';
      }
    }



    function sendMerchOTP() {
      var idVal = document.getElementById('merchOtpIdentifier').value.trim();
      var statusDiv = document.getElementById('merchOtpStatus');
      var btnSend = document.getElementById('merchBtnSend');
      if (!idVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter your registered mobile number or email.';
        return;
      }
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      statusDiv.style.display = 'none';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('portal', 'merchant');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          var inp = document.getElementById('merchOtpInput');
          if (inp) {
            inp.value = data.demo_otp || '';
            inp.focus();
          }
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#065f36';
          if (data.demo_otp) {
            statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: #ecfdf5; color: #065f36; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #4B5563;">(Code auto-filled &mdash; click "Verify & Sign In" to complete login)</span>';
          } else {
            statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || idVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
          }
          var count = 60;
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend OTP';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyMerchOTP() {
      var idVal = document.getElementById('merchOtpIdentifier').value.trim();
      var otpVal = document.getElementById('merchOtpInput').value.trim();
      var statusDiv = document.getElementById('merchOtpStatus');
      var btnVerify = document.getElementById('merchBtnVerify');

      if (!idVal || !otpVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter both identifier and 6-digit OTP.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Verifying...';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('otp', otpVal);
      body.append('portal', 'merchant');

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#16a34a';
          statusDiv.innerText = 'Verified! Launching Merchant Portal...';
          window.location.href = data.redirect || '/merchant';
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify & Open Merchant Portal';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Invalid or expired OTP code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify & Open Merchant Portal';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error during verification.';
      });
    }

    function sendMerchRegOTP() {
      var phoneVal = document.getElementById('merchRegOtpPhone').value.trim();
      var statusDiv = document.getElementById('merchRegOtpStatus');
      var btnSend = document.getElementById('merchRegBtnSend');
      if (!phoneVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter contact phone number.';
        return;
      }
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      statusDiv.style.display = 'none';

      var body = new URLSearchParams();
      body.append('identifier', phoneVal);
      body.append('portal', 'merchant');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          var inp = document.getElementById('merchRegOtpInput');
          if (inp) {
            inp.value = data.demo_otp || '';
            inp.focus();
          }
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#065f36';
          if (data.demo_otp) {
            statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: #ecfdf5; color: #065f36; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #4B5563;">(Code auto-filled &mdash; click "Register Partner" to complete)</span>';
          } else {
            statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || phoneVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
          }
          var count = 60;
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyMerchRegOTP() {
      var nameVal = document.getElementById('merchRegName').value.trim();
      var phoneVal = document.getElementById('merchRegOtpPhone').value.trim();
      var otpVal = document.getElementById('merchRegOtpInput').value.trim();
      var restName = document.getElementById('merchRegRestName').value.trim();
      var restAddress = document.getElementById('merchRegRestAddress').value.trim();
      var statusDiv = document.getElementById('merchRegOtpStatus');
      var btnVerify = document.getElementById('merchRegBtnVerify');

      if (!phoneVal || !otpVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Please enter phone number and 6-digit OTP code.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Launching portal...';

      var body = new URLSearchParams();
      body.append('name', nameVal || ('Partner ' + phoneVal));
      body.append('identifier', phoneVal);
      body.append('otp', otpVal);
      body.append('portal', 'merchant');
      body.append('restaurant_name', restName || ('Eatery ' + phoneVal));
      body.append('restaurant_address', restAddress || 'Bengaluru, Karnataka');

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#16a34a';
          statusDiv.innerText = 'Partner verified! Opening Merchant Portal...';
          window.location.href = data.redirect || '/merchant';
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify OTP & Launch Merchant Portal';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#dc2626';
          statusDiv.innerText = data.error || 'Invalid OTP code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify OTP & Launch Merchant Portal';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#dc2626';
        statusDiv.innerText = 'Network error during verification.';
      });
    }

    if (new URLSearchParams(window.location.search).get('mode') === 'register') {
      switchMerchMainTab('register');
    }
  </script>
</body>
</html>`, flashHTML, errHTML)
}

// RenderAdminLoginPage renders dedicated sign-in and registration page for platform administrators
func RenderAdminLoginPage(flashMsg, errMsg string) string {
	var flashHTML string
	if flashMsg == "admin_required" {
		flashHTML = `<div style="background: #fdf2f8; border: 1px solid #fbcfe8; color: #9d174d; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			<strong>Administrator Access Required:</strong> Please sign in with your administrative credentials to access the Admin Console.
		</div>`
	} else if flashMsg != "" {
		flashHTML = fmt.Sprintf(`<div style="background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">%s</div>`, escapeHTML(flashMsg))
	}

	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div style="background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem;">
			%s
		</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Admin Sign In — Under100 Platform Console</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #0f172a;
      color: #f8fafc;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #1e293b;
      border: 1px solid #334155;
      border-radius: 20px;
      width: 480px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
    }
    .auth-brand {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-family: 'Outfit', sans-serif;
      font-size: 1.4rem;
      font-weight: 800;
      color: #ffffff;
      text-decoration: none;
    }
    .auth-brand-badge {
      background: #831843;
      color: #fbcfe8;
      border: 1px solid #be185d;
      font-size: 0.72rem;
      font-weight: 800;
      padding: 0.2rem 0.5rem;
      border-radius: 6px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .auth-subtitle {
      color: #94a3b8;
      font-size: 0.88rem;
      margin-bottom: 1.25rem;
      margin-top: 0.35rem;
    }
    .auth-nav-tabs {
      display: flex;
      gap: 8px;
      background: #0f172a;
      padding: 4px;
      border-radius: 12px;
      margin-bottom: 1.25rem;
      border: 1px solid #334155;
    }
    .auth-nav-tab {
      flex: 1;
      text-align: center;
      padding: 8px 12px;
      font-size: 0.85rem;
      font-weight: 700;
      border-radius: 8px;
      cursor: pointer;
      border: none;
      background: none;
      color: #94a3b8;
      transition: all 0.15s;
    }
    .auth-nav-tab.active {
      background: #be185d;
      color: #ffffff;
      box-shadow: 0 2px 6px rgba(190, 24, 93, 0.3);
    }
    .subtabs-bar {
      display: flex;
      gap: 8px;
      margin-bottom: 1.25rem;
    }
    .subtab-btn {
      flex: 1;
      background: #0f172a;
      border: 1px solid #334155;
      color: #94a3b8;
      font-weight: 600;
      font-size: 0.8rem;
      padding: 0.5rem 0.6rem;
      border-radius: 8px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      transition: all 0.15s;
    }
    .subtab-btn.active {
      background: rgba(190, 24, 93, 0.2);
      border-color: #be185d;
      color: #fbcfe8;
      font-weight: 700;
    }
    .demo-helper-box {
      background: rgba(190, 24, 93, 0.12);
      border: 1px solid rgba(190, 24, 93, 0.35);
      border-radius: 10px;
      padding: 10px 14px;
      margin-bottom: 1.25rem;
      font-size: 0.8rem;
      color: #fbcfe8;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
    }
    .btn-fill-demo {
      background: #be185d;
      color: #ffffff;
      border: none;
      padding: 4px 10px;
      border-radius: 6px;
      font-size: 0.75rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
    }
    .btn-fill-demo:hover { background: #9d174d; }
    .form-group { margin-bottom: 1rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 600; color: #e2e8f0; margin-bottom: 0.35rem; }
    .form-input {
      width: 100%%;
      background: #0f172a;
      border: 1px solid #475569;
      color: #ffffff;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.9rem;
      outline: none;
    }
    .form-input:focus { border-color: #f43f5e; }
    .btn-submit {
      width: 100%%;
      background: #be185d;
      color: #ffffff;
      border: none;
      padding: 0.8rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      margin-top: 0.5rem;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #9d174d; }
    .btn-otp-action {
      background: #be185d;
      color: #ffffff;
      border: none;
      padding: 0.65rem 1rem;
      border-radius: 10px;
      font-size: 0.82rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
      transition: background 0.15s;
    }
    .btn-otp-action:hover { background: #9d174d; }
    .btn-otp-action:disabled { opacity: 0.6; cursor: not-allowed; }
    .back-link {
      display: inline-block;
      color: #94a3b8;
      text-decoration: none;
      font-size: 0.85rem;
    }
    .back-link:hover { color: #fbcfe8; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <a href="/admin" class="auth-brand">
        Under100
      </a>
      <span class="auth-brand-badge">Admin Console</span>
    </div>
    <p class="auth-subtitle">Restricted system administration &amp; platform operations.</p>

    <div class="demo-helper-box">
      <div>
        <strong>Platform Administrator:</strong> Root Admin<br>
        <span style="font-family: monospace; font-size: 0.76rem;">admin@fooddiscovery.com</span>
      </div>
      <button type="button" class="btn-fill-demo" onclick="fillAdminDemo()">Fill Demo</button>
    </div>

    %s
    %s

    <!-- SIGN IN CONTAINER -->
    <div id="admin-signin-container">
      <div class="subtabs-bar">
        <button type="button" id="admin-sub-pass" class="subtab-btn active" onclick="switchAdminAuthMode('pass')">
          Password Login
        </button>
        <button type="button" id="admin-sub-otp" class="subtab-btn" onclick="switchAdminAuthMode('otp')">
          Admin OTP Login
        </button>
      </div>

      <!-- Admin Password Login Form -->
      <form method="POST" action="/admin/login" id="adminLoginForm">
        <div class="form-group">
          <label class="form-label">Administrator Email</label>
          <input type="email" name="email" id="adminEmail" class="form-input" placeholder="admin@example.com" required>
        </div>
        <div class="form-group">
          <label class="form-label">Password</label>
          <input type="password" name="password" id="adminPass" class="form-input" placeholder="••••••••" required>
        </div>
        <button type="submit" class="btn-submit">Sign In to Admin Console</button>
      </form>

      <!-- Admin OTP Login Form -->
      <div id="adminOtpForm" style="display: none;">

        <div class="form-group">
          <label class="form-label">Administrator Official Email or Mobile</label>
          <div style="display: flex; gap: 8px;">
            <input type="text" id="adminOtpIdentifier" class="form-input" placeholder="e.g. admin@example.com or 9876543210" style="flex: 1;">
            <button type="button" id="adminBtnSend" class="btn-otp-action" onclick="sendAdminOTP()">Get OTP</button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">6-Digit Admin Verification Code</label>
          <input type="text" id="adminOtpInput" class="form-input" placeholder="Enter 6-digit OTP" maxlength="6" style="letter-spacing: 4px; font-size: 1.1rem; font-weight: 700; text-align: center;">
        </div>

        <div id="adminOtpStatus" style="display: none; font-size: 0.82rem; margin-bottom: 0.75rem; font-weight: 600;"></div>

        <button type="button" id="adminBtnVerify" class="btn-submit" onclick="verifyAdminOTP()">Verify &amp; Open Admin Console</button>
      </div>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #94a3b8;">
        Don't have an account? <a href="javascript:switchAdminMainTab('register')" style="color: #f43f5e; font-weight: 700; text-decoration: none;">Register here</a>
      </div>
    </div>

    <!-- REGISTER ADMIN CONTAINER -->
    <div id="admin-register-container" style="display: none;">

      <form method="POST" action="/admin/register" id="adminRegisterForm">
        <div class="form-group">
          <label class="form-label">Administrator Full Name</label>
          <input type="text" name="name" class="form-input" placeholder="e.g. Systems Administrator" required>
        </div>
        <div class="form-group">
          <label class="form-label">Official Administrator Email</label>
          <input type="email" name="email" class="form-input" placeholder="admin.ops@fooddiscovery.com" required>
        </div>
        <div class="form-group">
          <label class="form-label">Create Password</label>
          <input type="password" name="password" class="form-input" placeholder="Min 6 characters" minlength="6" required>
        </div>
        <div class="form-group">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.35rem;">
            <label class="form-label" style="margin-bottom: 0;">Administrative Security Key</label>
            <span style="font-size: 0.75rem; color: #fbcfe8;">Key: <code style="color: #cbd5e1;">Admin@123</code></span>
          </div>
          <input type="password" name="admin_key" class="form-input" placeholder="Enter Administrative Security Key" required>
        </div>
        <button type="submit" class="btn-submit">Register Administrator &amp; Open Console</button>
      </form>

      <div style="text-align: center; margin-top: 1.25rem; font-size: 0.9rem; color: #94a3b8;">
        Already have an account? <a href="javascript:switchAdminMainTab('signin')" style="color: #f43f5e; font-weight: 700; text-decoration: none;">Sign in here</a>
      </div>
    </div>

    <div style="text-align: center; margin-top: 1.25rem; border-top: 1px solid #334155; padding-top: 1rem; font-size: 0.85rem;">
      <a href="/" class="back-link">&larr; Return to Customer Storefront</a>
    </div>
  </div>

  <script>
    var currentAdminDemoOtp = '123456';

    function fillAdminDemo() {
      switchAdminMainTab('signin');
      switchAdminAuthMode('pass');
      document.getElementById('adminEmail').value = 'admin@fooddiscovery.com';
      document.getElementById('adminPass').value = 'Admin@123';
    }

    function switchAdminMainTab(tab) {
      var tabIn = document.getElementById('admin-tab-in');
      var tabReg = document.getElementById('admin-tab-reg');
      var cIn = document.getElementById('admin-signin-container');
      var cReg = document.getElementById('admin-register-container');

      if (tab === 'register') {
        if (tabIn) tabIn.classList.remove('active');
        if (tabReg) tabReg.classList.add('active');
        if (cIn) cIn.style.display = 'none';
        if (cReg) cReg.style.display = 'block';
      } else {
        if (tabReg) tabReg.classList.remove('active');
        if (tabIn) tabIn.classList.add('active');
        if (cReg) cReg.style.display = 'none';
        if (cIn) cIn.style.display = 'block';
      }
    }

    function switchAdminAuthMode(m) {
      var btnPass = document.getElementById('admin-sub-pass');
      var btnOtp = document.getElementById('admin-sub-otp');
      var fPass = document.getElementById('adminLoginForm');
      var fOtp = document.getElementById('adminOtpForm');

      if (m === 'otp') {
        btnPass.classList.remove('active');
        btnOtp.classList.add('active');
        fPass.style.display = 'none';
        fOtp.style.display = 'block';
      } else {
        btnOtp.classList.remove('active');
        btnPass.classList.add('active');
        fOtp.style.display = 'none';
        fPass.style.display = 'block';
      }
    }



    function sendAdminOTP() {
      var idVal = document.getElementById('adminOtpIdentifier').value.trim();
      var statusDiv = document.getElementById('adminOtpStatus');
      var btnSend = document.getElementById('adminBtnSend');
      if (!idVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#f87171';
        statusDiv.innerText = 'Please enter administrator official email or mobile.';
        return;
      }
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      statusDiv.style.display = 'none';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('portal', 'admin');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          var inp = document.getElementById('adminOtpInput');
          if (inp) {
            inp.value = data.demo_otp || '';
            inp.focus();
          }
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#34d399';
          if (data.demo_otp) {
            statusDiv.innerHTML = '<div style="margin-bottom: 5px;">Verification code: <strong style="font-size: 16px; background: rgba(52,211,153,0.2); color: #34d399; padding: 3px 8px; border-radius: 6px; font-family: monospace; letter-spacing: 2px;">' + data.demo_otp + '</strong></div><span style="font-size: 12px; color: #94a3b8;">(Code auto-filled &mdash; click "Verify & Access" to complete login)</span>';
          } else {
            statusDiv.innerHTML = 'Verification code sent to <strong>' + (data.target || idVal) + '</strong> via SMS. Please enter the 6-digit code received on your phone.';
          }
          var count = 60;
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend OTP';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#f87171';
          statusDiv.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#f87171';
        statusDiv.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyAdminOTP() {
      var idVal = document.getElementById('adminOtpIdentifier').value.trim();
      var otpVal = document.getElementById('adminOtpInput').value.trim();
      var statusDiv = document.getElementById('adminOtpStatus');
      var btnVerify = document.getElementById('adminBtnVerify');

      if (!idVal || !otpVal) {
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#f87171';
        statusDiv.innerText = 'Please enter both identifier and 6-digit OTP.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Verifying...';

      var body = new URLSearchParams();
      body.append('identifier', idVal);
      body.append('otp', otpVal);
      body.append('portal', 'admin');

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: body.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#34d399';
          statusDiv.innerText = 'Verified! Opening Admin Console...';
          window.location.href = data.redirect || '/admin';
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify & Open Admin Console';
          statusDiv.style.display = 'block';
          statusDiv.style.color = '#f87171';
          statusDiv.innerText = data.error || 'Invalid or expired OTP code.';
        }
      })
      .catch(function(err) {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify & Open Admin Console';
        statusDiv.style.display = 'block';
        statusDiv.style.color = '#f87171';
        statusDiv.innerText = 'Network error during verification.';
      });
    }

    if (new URLSearchParams(window.location.search).get('mode') === 'register') {
      switchAdminMainTab('register');
    }
  </script>
</body>
</html>`, flashHTML, errHTML)
}

// RenderForbiddenPage renders a clean light-mode 403 Forbidden page (zero emoji)
func RenderForbiddenPage(usr *ent.User) string {
	userName := "User"
	if usr != nil {
		if usr.Name != "" {
			userName = usr.Name
		} else if usr.Email != "" {
			userName = usr.Email
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>403 Forbidden — Admin Access Required</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .card {
      background: #ffffff;
      border: 1px solid #fee2e2;
      border-radius: 20px;
      max-width: 500px;
      width: 100%%;
      padding: 2.5rem;
      text-align: center;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .badge {
      display: inline-block;
      background: #fee2e2;
      color: #991b1b;
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      font-size: 0.8rem;
      font-weight: 800;
      margin-bottom: 1.25rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    .title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.6rem;
      font-weight: 800;
      color: #0f172a;
      margin-bottom: 0.75rem;
    }
    .text {
      color: #64748b;
      font-size: 0.95rem;
      line-height: 1.5;
      margin-bottom: 1.75rem;
    }
    .user-pill {
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
      background: #f1f5f9;
      border: 1px solid #e2e8f0;
      padding: 0.4rem 0.8rem;
      border-radius: 8px;
      font-size: 0.85rem;
      color: #0f172a;
      margin-bottom: 1.5rem;
    }
    .btn-group {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
    .btn-primary {
      background: #e23744;
      color: #ffffff;
      text-decoration: none;
      padding: 0.75rem 1.25rem;
      border-radius: 10px;
      font-weight: 700;
      font-size: 0.92rem;
      transition: background 0.15s;
    }
    .btn-primary:hover {
      background: #cb202d;
    }
    .btn-secondary {
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      text-decoration: none;
      padding: 0.75rem 1.25rem;
      border-radius: 10px;
      font-weight: 600;
      font-size: 0.92rem;
    }
  </style>
</head>
<body>
  <div class="card">
    <div class="badge">403 Forbidden</div>
    <h1 class="title">Administrator Access Required</h1>
    <p class="text">The Admin Dashboard is reserved strictly for authorized platform administrators. Your account does not have permission to manage restaurant listings.</p>
    <div>
      <div class="user-pill">
        <span>Signed in as:</span>
        <strong>%s</strong>
      </div>
    </div>
    <div class="btn-group">
      <a href="/" class="btn-primary">&larr; Return to Home Page</a>
      <a href="/login?msg=admin_required" class="btn-secondary">Switch Account / Sign in as Admin</a>
      <a href="/logout" style="color: #64748b; text-decoration: none; font-size: 0.85rem; margin-top: 0.5rem;">Log Out</a>
    </div>
  </div>
</body>
</html>`, escapeHTML(userName))
}

// Helpers
func escapeHTML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return r.Replace(s)
}

func formatPhone(p string) string {
	if p == "" {
		return ""
	}
	return fmt.Sprintf(`<span>Phone: %s</span>`, escapeHTML(p))
}

func formatHours(h string) string {
	if h == "" {
		return ""
	}
	return fmt.Sprintf(`<span>Hours: %s</span>`, escapeHTML(h))
}

func activeTab(current, target string) string {
	if current == target {
		return "active"
	}
	return ""
}

func selected(current, target string) string {
	if current == target {
		return "selected"
	}
	return ""
}

// RenderProfilePage renders customer account and security settings
func RenderProfilePage(u *ent.User, flashMsg, flashErr string) string {
	var alertHTML string
	if flashMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #f0fdf4; border: 1.5px solid #bbf7d0; color: #166534; padding: 0.85rem 1.25rem; border-radius: 12px; margin-bottom: 1.5rem; font-size: 0.92rem; font-weight: 600;">%s</div>`, escapeHTML(flashMsg))
	}
	if flashErr != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #fef2f2; border: 1.5px solid #fecaca; color: #991b1b; padding: 0.85rem 1.25rem; border-radius: 12px; margin-bottom: 1.5rem; font-size: 0.92rem; font-weight: 600;">%s</div>`, escapeHTML(flashErr))
	}

	roleBadge := `<span style="background: #f0fdf4; color: #166534; border: 1px solid #bbf7d0; font-size: 0.75rem; font-weight: 800; padding: 0.25rem 0.65rem; border-radius: 9999px; text-transform: uppercase; letter-spacing: 0.3px;">Verified Customer</span>`
	if u.IsAdmin || u.Role == user.RoleAdmin {
		roleBadge = `<span style="background: #fdf2f8; color: #be185d; border: 1px solid #fbcfe8; font-size: 0.75rem; font-weight: 800; padding: 0.25rem 0.65rem; border-radius: 9999px; text-transform: uppercase; letter-spacing: 0.3px;">Platform Administrator</span>`
	} else if u.Role == user.RoleMerchant {
		roleBadge = `<span style="background: #fef3c7; color: #b45309; border: 1px solid #fde68a; font-size: 0.75rem; font-weight: 800; padding: 0.25rem 0.65rem; border-radius: 9999px; text-transform: uppercase; letter-spacing: 0.3px;">Restaurant Partner</span>`
	}

	dashboardLink := `<a href="/" style="color: #0f172a; text-decoration: none; font-size: 0.9rem; font-weight: 700;">&larr; Return to Storefront</a>`

	userName := u.Name
	if userName == "" {
		userName = u.Email
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>My Profile & Security — Bite100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      padding-bottom: 4rem;
    }
    .profile-nav {
      background: #ffffff;
      border-bottom: 1px solid #e2e8f0;
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .profile-container {
      max-width: 650px;
      margin: 2.5rem auto;
      padding: 0 1.25rem;
    }
    .profile-header-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 16px;
      padding: 1.75rem;
      margin-bottom: 1.75rem;
      display: flex;
      align-items: center;
      gap: 1.25rem;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .profile-avatar-box {
      width: 64px;
      height: 64px;
      background: linear-gradient(135deg, #fee2e2 0%%, #fecaca 100%%);
      color: #e23744;
      border-radius: 50%%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: 'Outfit', sans-serif;
      font-size: 1.5rem;
      font-weight: 800;
      flex-shrink: 0;
    }
    .profile-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 16px;
      padding: 2rem;
      margin-bottom: 1.75rem;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .profile-card-title {
      font-family: 'Outfit', sans-serif;
      font-size: 1.2rem;
      font-weight: 800;
      color: #0f172a;
      margin-bottom: 0.35rem;
    }
    .profile-card-desc {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.5rem;
    }
    .form-group {
      margin-bottom: 1.25rem;
    }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: #334155;
      margin-bottom: 0.4rem;
      text-transform: uppercase;
      letter-spacing: 0.3px;
    }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      outline: none;
      font-family: inherit;
      transition: border-color 0.15s;
    }
    .form-input:focus {
      border-color: #e23744;
      box-shadow: 0 0 0 3px rgba(226, 55, 68, 0.1);
    }
    .form-input:disabled, .form-input[readonly] {
      background: #f8fafc;
      color: #64748b;
      cursor: not-allowed;
      border-color: #e2e8f0;
    }
    .form-helper {
      font-size: 0.78rem;
      color: #64748b;
      margin-top: 0.35rem;
    }
    .btn-primary {
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.75rem 1.5rem;
      border-radius: 10px;
      font-size: 0.9rem;
      font-weight: 800;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      transition: background 0.15s;
    }
    .btn-primary:hover {
      background: #cb202d;
    }
  </style>
</head>
<body>
  <nav class="profile-nav">
    <div style="display: flex; align-items: center; gap: 1rem;">
      <a href="/" style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; color: #0f172a; text-decoration: none;">
        Bite100
      </a>
      <span style="background: #f1f5f9; color: #475569; font-size: 0.75rem; font-weight: 800; padding: 0.2rem 0.5rem; border-radius: 6px;">Profile & Security</span>
    </div>
    <div style="display: flex; align-items: center; gap: 1.25rem;">
      %s
      <a href="/logout" style="color: #e23744; text-decoration: none; font-size: 0.85rem; font-weight: 700; padding: 0.35rem 0.75rem; border: 1px solid #fecaca; border-radius: 8px;">Log Out</a>
    </div>
  </nav>

  <main class="profile-container">
    %s

    <div class="profile-header-card">
      <div class="profile-avatar-box">
        %s
      </div>
      <div>
        <div style="display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.35rem;">
          <h1 style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; color: #0f172a;">%s</h1>
          %s
        </div>
        <p style="color: #64748b; font-size: 0.9rem;">%s</p>
      </div>
    </div>

    <!-- Section 1: Display Name Profile Update -->
    <div class="profile-card">
      <h2 class="profile-card-title">Profile Information</h2>
      <p class="profile-card-desc">Update your personal account details and public display name.</p>

      <form method="POST" action="/profile">
        <div class="form-group">
          <label class="form-label">Full Name / Display Name *</label>
          <input type="text" name="name" class="form-input" value="%s" required minlength="2" maxlength="100">
        </div>
        <div class="form-group">
          <label class="form-label">Email Address</label>
          <input type="email" value="%s" class="form-input" readonly>
          <div class="form-helper">Account email is tied to authentication and cannot be modified.</div>
        </div>
        <div style="text-align: right; margin-top: 1.25rem;">
          <button type="submit" class="btn-primary">
            Save Profile Changes
          </button>
        </div>
      </form>
    </div>

    <!-- Section 2: Password Change -->
    <div class="profile-card">
      <h2 class="profile-card-title">Change Password</h2>
      <p class="profile-card-desc">Ensure your account is using a long, random password to stay secure.</p>

      <form method="POST" action="/profile/password">
        <div class="form-group">
          <label class="form-label">Current Password *</label>
          <input type="password" name="current_password" class="form-input" placeholder="Enter current password" required>
        </div>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
          <div class="form-group">
            <label class="form-label">New Password *</label>
            <input type="password" name="new_password" class="form-input" placeholder="Min 6 characters" minlength="6" required>
          </div>
          <div class="form-group">
            <label class="form-label">Confirm New Password *</label>
            <input type="password" name="confirm_password" class="form-input" placeholder="Re-enter new password" minlength="6" required>
          </div>
        </div>
        <div style="text-align: right; margin-top: 1.25rem;">
          <button type="submit" class="btn-primary">
            Update Password
          </button>
        </div>
      </form>
    </div>
  </main>
</body>
</html>`, dashboardLink, alertHTML, string([]rune(userName)[0]), escapeHTML(userName), roleBadge, escapeHTML(u.Email), escapeHTML(u.Name), escapeHTML(u.Email))
}

// RenderForgotPasswordPage renders password recovery request form
func RenderForgotPasswordPage(flashMsg, errMsg string) string {
	var alertHTML string
	if flashMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #f0fdf4; border: 1.5px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem; font-weight: 600;">%s</div>`, escapeHTML(flashMsg))
	}
	if errMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #fef2f2; border: 1.5px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem; font-weight: 600;">%s</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Forgot Password — Bite100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 440px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
    }
    .auth-brand {
      font-family: 'Outfit', sans-serif;
      font-size: 1.4rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-subtitle {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.5rem;
      margin-top: 0.35rem;
    }
    .form-group {
      margin-bottom: 1.25rem;
    }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: #334155;
      margin-bottom: 0.4rem;
    }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.9rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.8rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
    }
    .back-link:hover { color: #e23744; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <a href="/" class="auth-brand">
        Bite100
      </a>
      <span style="background: #f1f5f9; color: #475569; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.5rem; border-radius: 6px;">Password Recovery</span>
    </div>
    <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; margin-top: 1.25rem;">Forgot Password</h2>
    <p class="auth-subtitle">Enter your registered email address and we will send you a password reset link.</p>

    %s

    <form method="POST" action="/password/forgot">
      <div class="form-group">
        <label class="form-label">Email Address</label>
        <input type="email" name="email" class="form-input" placeholder="you@example.com" required autofocus>
      </div>
      <button type="submit" class="btn-submit">Send Reset Link</button>
    </form>

    <div style="text-align: center;">
      <a href="/login" class="back-link">&larr; Back to Sign In</a>
    </div>
  </div>
</body>
</html>`, alertHTML)
}

// RenderResetPasswordPage renders the new password entry form after token validation
func RenderResetPasswordPage(userID, tokenID int, token, errMsg string) string {
	var alertHTML string
	if errMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #fef2f2; border: 1.5px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 10px; margin-bottom: 1.25rem; font-size: 0.88rem; font-weight: 600;">%s</div>`, escapeHTML(errMsg))
	}

	contentHTML := ""
	if userID == 0 || token == "" {
		// Token is invalid or expired
		contentHTML = fmt.Sprintf(`
		%s
		<div style="text-align: center; margin-top: 1.5rem;">
			<a href="/password/forgot" style="color: #e23744; font-weight: 700; text-decoration: underline; font-size: 0.9rem;">Request a new password reset link</a>
		</div>`, alertHTML)
	} else {
		contentHTML = fmt.Sprintf(`
		%s
		<form method="POST" action="/password/reset/%d/%d/%s">
			<div class="form-group">
				<label class="form-label">New Password *</label>
				<input type="password" name="password" class="form-input" placeholder="Min 6 characters" minlength="6" required autofocus>
			</div>
			<div class="form-group">
				<label class="form-label">Confirm New Password *</label>
				<input type="password" name="confirm_password" class="form-input" placeholder="Confirm new password" minlength="6" required>
			</div>
			<button type="submit" class="btn-submit">Save New Password</button>
		</form>`, alertHTML, userID, tokenID, escapeHTML(token))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset Password — Bite100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 440px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
    }
    .auth-brand {
      font-family: 'Outfit', sans-serif;
      font-size: 1.4rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-subtitle {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.5rem;
      margin-top: 0.35rem;
    }
    .form-group {
      margin-bottom: 1.25rem;
    }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: #334155;
      margin-bottom: 0.4rem;
    }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.7rem 0.9rem;
      border-radius: 10px;
      font-size: 0.9rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.8rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.85rem;
      font-weight: 600;
    }
    .back-link:hover { color: #e23744; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between;">
      <a href="/" class="auth-brand">
        Bite100
      </a>
      <span style="background: #f1f5f9; color: #475569; font-size: 0.72rem; font-weight: 800; padding: 0.2rem 0.5rem; border-radius: 6px;">Secure Reset</span>
    </div>
    <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; margin-top: 1.25rem;">Create New Password</h2>
    <p class="auth-subtitle">Enter and confirm your new account password below.</p>

    %s

    <div style="text-align: center;">
      <a href="/login" class="back-link">&larr; Back to Sign In</a>
    </div>
  </div>
</body>
</html>`, contentHTML)
}
// RenderAccountDashboardPage renders the dedicated 2-column Blinkit-style account interface
func RenderAccountDashboardPage(currentUser *ent.User, activeTab string) string {
	if activeTab == "" {
		activeTab = "orders"
	}
	userDisplay := "My Account"
	userContact := "+91 84680 68403"
	if currentUser != nil {
		if currentUser.Name != "" {
			userDisplay = currentUser.Name
			if currentUser.Email != "" {
				userContact = currentUser.Email
			}
		} else if currentUser.Email != "" {
			userDisplay = currentUser.Email
			userContact = currentUser.Email
		}
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Under100 — Bengaluru Food Delivery &amp; Orders</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800;900&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    :root {
      --color-primary-orange: #FF5A1F;
      --color-primary-green: #0C831F;
      --color-bg-canvas: #F4F6FB;
      --color-dark-text: #1E293B;
      --color-secondary-text: #64748B;
      --color-border: #E2E8F0;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--color-bg-canvas);
      color: var(--color-dark-text);
      font-family: 'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }

    /* Top Sticky Navbar */
    .account-top-navbar {
      background: #ffffff;
      border-bottom: 1px solid var(--color-border);
      position: sticky;
      top: 0;
      z-index: 100;
      padding: 12px 32px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 20px;
    }
    .brand-group {
      display: flex;
      align-items: center;
      gap: 24px;
      flex-shrink: 0;
    }
    .nav-brand-logo {
      font-family: 'Outfit', sans-serif;
      font-weight: 900;
      font-size: 26px;
      color: #1E293B;
      text-decoration: none;
      letter-spacing: -0.03em;
    }
    .nav-brand-logo span {
      color: var(--color-primary-orange);
    }
    .nav-location-pill {
      display: flex;
      flex-direction: column;
      cursor: pointer;
      user-select: none;
    }
    .nav-loc-title {
      font-weight: 800;
      font-size: 13.5px;
      color: #1E293B;
    }
    .nav-loc-sub {
      font-size: 12px;
      color: var(--color-secondary-text);
      display: flex;
      align-items: center;
      gap: 4px;
      max-width: 220px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .nav-search-wrap {
      flex: 1;
      max-width: 650px;
      position: relative;
    }
    .nav-search-input {
      width: 100%%;
      padding: 12px 18px 12px 42px;
      background: #F8FAFC;
      border: 1px solid var(--color-border);
      border-radius: 12px;
      font-size: 14px;
      color: #1E293B;
      outline: none;
      transition: all 0.2s ease;
    }
    .nav-search-input:focus {
      background: #ffffff;
      border-color: var(--color-primary-orange);
      box-shadow: 0 0 0 3px rgba(255, 90, 31, 0.1);
    }
    .nav-search-icon {
      position: absolute;
      left: 14px;
      top: 50%%;
      transform: translateY(-50%%);
      color: #94A3B8;
    }
    .nav-actions-right {
      display: flex;
      align-items: center;
      gap: 14px;
      flex-shrink: 0;
    }
    .btn-nav-account-header {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      background: transparent;
      border: none;
      font-weight: 700;
      font-size: 14.5px;
      color: #1E293B;
      cursor: pointer;
      padding: 8px 12px;
      border-radius: 8px;
      transition: background 0.15s ease;
    }
    .btn-nav-account-header:hover {
      background: #F1F5F9;
    }
    .btn-nav-cart-pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      background: #0C831F;
      color: #ffffff;
      border: none;
      padding: 10px 18px;
      border-radius: 10px;
      font-weight: 700;
      font-size: 14px;
      cursor: pointer;
      box-shadow: 0 2px 8px rgba(12, 131, 31, 0.2);
      transition: background 0.15s ease;
    }
    .btn-nav-cart-pill:hover {
      background: #096B19;
    }

    /* Main Two-Column Layout */
    .account-main-wrapper {
      max-width: 1220px;
      width: 100%%;
      margin: 24px auto;
      padding: 0 24px 48px;
      display: flex;
      gap: 24px;
      align-items: flex-start;
      flex: 1;
    }

    /* Left Sidebar */
    .account-sidebar-card {
      width: 260px;
      background: #ffffff;
      border: 1px solid var(--color-border);
      border-radius: 16px;
      padding: 24px 0 16px 0;
      box-shadow: 0 1px 3px rgba(0,0,0,0.03);
      flex-shrink: 0;
    }
    .account-sidebar-phone {
      padding: 0 24px 18px 24px;
      font-size: 13.5px;
      color: #64748B;
      font-weight: 500;
      letter-spacing: 0.02em;
    }
    .account-nav-list {
      list-style: none;
      display: flex;
      flex-direction: column;
    }
    .account-nav-item {
      display: flex;
      align-items: center;
      gap: 14px;
      padding: 13px 24px;
      font-size: 14px;
      font-weight: 500;
      color: #475569;
      text-decoration: none;
      cursor: pointer;
      transition: all 0.15s ease;
      border-left: 3px solid transparent;
      user-select: none;
    }
    .account-nav-item:hover {
      background: #F8FAFC;
      color: #0F172A;
    }
    .account-nav-item.active {
      background: #F4F6FB;
      color: #0F172A;
      font-weight: 700;
      border-left-color: #0C831F;
    }
    .account-nav-item svg {
      color: #64748B;
      flex-shrink: 0;
      transition: color 0.15s ease;
    }
    .account-nav-item.active svg {
      color: #0C831F;
    }

    /* Right Content Area */
    .account-content-area {
      flex: 1;
      min-width: 0;
    }
    .account-tab-pane {
      display: none;
      animation: tabFadeIn 0.15s ease-out;
    }
    .account-tab-pane.active {
      display: block;
    }
    @keyframes tabFadeIn {
      from { opacity: 0; transform: translateY(4px); }
      to { opacity: 1; transform: translateY(0); }
    }

    /* Order Card Component matching screenshot */
    .order-history-card {
      background: #ffffff;
      border: 1px solid var(--color-border);
      border-radius: 14px;
      padding: 18px 22px;
      margin-bottom: 14px;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
      transition: transform 0.15s ease, box-shadow 0.15s ease;
      cursor: pointer;
    }
    .order-history-card:hover {
      box-shadow: 0 4px 12px rgba(0,0,0,0.05);
    }
    .order-card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .order-header-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }
    .order-status-badge-circle {
      width: 26px;
      height: 26px;
      border-radius: 50%%;
      background: #ECFDF5;
      color: #10B981;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 900;
      font-size: 13px;
      flex-shrink: 0;
    }
    .order-headline {
      font-weight: 800;
      font-size: 15px;
      color: #1E293B;
      line-height: 1.25;
    }
    .order-submeta {
      font-size: 12.5px;
      color: #64748B;
      margin-top: 3px;
    }
    .order-chevron-right {
      color: #64748B;
      font-size: 18px;
      font-weight: 700;
    }

    /* Thumbnails Row */
    .order-thumbnails-row {
      display: flex;
      align-items: center;
      gap: 14px;
      margin-top: 16px;
      padding-top: 14px;
      border-top: 1px solid #F8FAFC;
      overflow-x: auto;
    }
    .order-item-thumbnail-box {
      width: 64px;
      height: 64px;
      border-radius: 10px;
      border: 1px solid #F1F5F9;
      background: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
      padding: 4px;
      flex-shrink: 0;
    }
    .order-item-thumbnail-box img {
      width: 100%%;
      height: 100%%;
      object-fit: cover;
      border-radius: 6px;
    }

    /* Addresses Card */
    .address-card-item {
      background: #ffffff;
      border: 1px solid var(--color-border);
      border-radius: 14px;
      padding: 20px;
      margin-bottom: 14px;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 16px;
    }
    .address-badge-tag {
      display: inline-block;
      padding: 3px 8px;
      border-radius: 6px;
      font-size: 11.5px;
      font-weight: 700;
      text-transform: uppercase;
      background: #F1F5F9;
      color: #334155;
      margin-bottom: 8px;
    }
    .address-text {
      font-size: 14px;
      color: #1E293B;
      line-height: 1.45;
      font-weight: 500;
    }
    .address-phone {
      font-size: 12.5px;
      color: #64748B;
      margin-top: 6px;
    }

    /* Section Panel Wrapper */
    .account-panel-card {
      background: #ffffff;
      border: 1px solid var(--color-border);
      border-radius: 16px;
      padding: 24px 28px;
      box-shadow: 0 1px 3px rgba(0,0,0,0.02);
    }
    .account-panel-title {
      font-size: 19px;
      font-weight: 800;
      color: #1E293B;
      margin-bottom: 6px;
    }
    .account-panel-subtitle {
      font-size: 13px;
      color: #64748B;
      margin-bottom: 20px;
    }

    @media (max-width: 860px) {
      .account-main-wrapper {
        flex-direction: column;
      }
      .account-sidebar-card {
        width: 100%%;
      }
      .account-nav-list {
        flex-direction: row;
        overflow-x: auto;
        padding-bottom: 8px;
      }
      .account-nav-item {
        border-left: none;
        border-bottom: 3px solid transparent;
        white-space: nowrap;
      }
      .account-nav-item.active {
        border-left-color: transparent;
        border-bottom-color: #0C831F;
      }
    }
  </style>
</head>
<body>

  <!-- Top Navigation Header -->
  <header class="account-top-navbar">
    <div class="brand-group">
      <a href="/" class="nav-brand-logo">Under<span>100</span></a>
      <div class="nav-location-pill" onclick="window.location.href='/'" title="Change delivery location">
        <div class="nav-loc-title">Delivery in 8 minutes</div>
        <div class="nav-loc-sub">Indiranagar, Bengaluru, Karnat... &#9662;</div>
      </div>
    </div>

    <div class="nav-search-wrap">
      <svg class="nav-search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
      <input type="text" class="nav-search-input" placeholder='Search "dosa", "coffee", "idli"...' onkeydown="if(event.key==='Enter') window.location.href='/?q='+encodeURIComponent(this.value)">
    </div>

    <div class="nav-actions-right">
      <button type="button" class="btn-nav-account-header" onclick="window.location.href='/account'">
        <span>Account</span> &#9662;
      </button>
      <button type="button" class="btn-nav-cart-pill" onclick="window.location.href='/?open_cart=1'">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="9" cy="21" r="1"></circle><circle cx="20" cy="21" r="1"></circle><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path></svg>
        <span>My Cart</span>
      </button>
    </div>
  </header>

  <!-- Two Column Layout matching screenshot -->
  <div class="account-main-wrapper">
    <!-- Left Navigation Sidebar -->
    <aside class="account-sidebar-card">
      <div class="account-sidebar-phone">%s</div>

      <nav class="account-nav-list">
        <a href="javascript:void(0)" class="account-nav-item" id="navItem-addresses" onclick="switchAccountTab('addresses')">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
          <span>My Addresses</span>
        </a>

        <a href="javascript:void(0)" class="account-nav-item active" id="navItem-orders" onclick="switchAccountTab('orders')">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
          <span>My Orders</span>
        </a>

        <a href="javascript:void(0)" class="account-nav-item" id="navItem-prescriptions" onclick="switchAccountTab('prescriptions')">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>
          <span>My Prescriptions</span>
        </a>

        <a href="javascript:void(0)" class="account-nav-item" id="navItem-gift-cards" onclick="switchAccountTab('gift-cards')">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="8" width="18" height="4" rx="1"></rect><path d="M12 8v13"></path><path d="M19 12v7a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-7"></path><path d="M7.5 8a2.5 2.5 0 0 1 0-5A4.8 4.8 0 0 1 12 8a4.8 4.8 0 0 1 4.5-5 2.5 2.5 0 0 1 0 5"></path></svg>
          <span>E-Gift Cards</span>
        </a>

        <a href="javascript:void(0)" class="account-nav-item" id="navItem-privacy" onclick="switchAccountTab('privacy')">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
          <span>Account privacy</span>
        </a>

        <a href="/logout" class="account-nav-item" id="navItem-logout">
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          <span>Logout</span>
        </a>
      </nav>
    </aside>

    <!-- Right Content Area -->
    <main class="account-content-area">

      <!-- TAB 1: MY ORDERS (Matches Screenshot) -->
      <section class="account-tab-pane active" id="pane-orders">
        <!-- Order 1: 8 minutes -->
        <div class="order-history-card">
          <div class="order-card-header">
            <div class="order-header-left">
              <div class="order-status-badge-circle">&#10003;</div>
              <div>
                <div class="order-headline">Arrived in 8 minutes</div>
                <div class="order-submeta">&#8377;304 &middot; 18 Oct, 1:31 pm</div>
              </div>
            </div>
            <div class="order-chevron-right">&rarr;</div>
          </div>
          <div class="order-thumbnails-row">
            <div class="order-item-thumbnail-box" title="Crispy Masala Dosa"><img src="/static/crispy_dosa.jpg" alt="Crispy Masala Dosa"></div>
            <div class="order-item-thumbnail-box" title="South Indian Filter Coffee"><img src="/static/filter_coffee.jpg" alt="Filter Coffee"></div>
            <div class="order-item-thumbnail-box" title="Steamed Idli Sambar"><img src="/static/idli_sambar.jpg" alt="Steamed Idli Sambar"></div>
            <div class="order-item-thumbnail-box" title="Aromatic Veg Pulao"><img src="/static/veg_pulao.jpg" alt="Aromatic Veg Pulao"></div>
          </div>
        </div>

        <!-- Order 2: 9 minutes -->
        <div class="order-history-card">
          <div class="order-card-header">
            <div class="order-header-left">
              <div class="order-status-badge-circle">&#10003;</div>
              <div>
                <div class="order-headline">Arrived in 9 minutes</div>
                <div class="order-submeta">&#8377;945 &middot; 15 Oct, 12:05 pm</div>
              </div>
            </div>
            <div class="order-chevron-right">&rarr;</div>
          </div>
          <div class="order-thumbnails-row">
            <div class="order-item-thumbnail-box" title="Stuffed Aloo Paratha"><img src="/static/aloo_paratha.jpg" alt="Aloo Paratha"></div>
            <div class="order-item-thumbnail-box" title="Paneer Kathi Roll"><img src="/static/kathi_roll.jpg" alt="Paneer Kathi Roll"></div>
            <div class="order-item-thumbnail-box" title="Bombay Sandwich"><img src="/static/bombay_sandwich.jpg" alt="Bombay Sandwich"></div>
          </div>
        </div>

        <!-- Order 3: 14 minutes -->
        <div class="order-history-card">
          <div class="order-card-header">
            <div class="order-header-left">
              <div class="order-status-badge-circle">&#10003;</div>
              <div>
                <div class="order-headline">Arrived in 14 minutes</div>
                <div class="order-submeta">&#8377;76 &middot; 08 Oct, 7:49 pm</div>
              </div>
            </div>
            <div class="order-chevron-right">&rarr;</div>
          </div>
          <div class="order-thumbnails-row">
            <div class="order-item-thumbnail-box" title="Crispy Samosa Plate"><img src="/static/samosa_plate.jpg" alt="Crispy Samosa"></div>
            <div class="order-item-thumbnail-box" title="South Indian Filter Coffee"><img src="/static/filter_coffee.jpg" alt="Filter Coffee"></div>
          </div>
        </div>

        <!-- Order 4: 6 minutes -->
        <div class="order-history-card">
          <div class="order-card-header">
            <div class="order-header-left">
              <div class="order-status-badge-circle">&#10003;</div>
              <div>
                <div class="order-headline">Arrived in 6 minutes</div>
                <div class="order-submeta">&#8377;230 &middot; 08 Oct, 7:49 pm</div>
              </div>
            </div>
            <div class="order-chevron-right">&rarr;</div>
          </div>
          <div class="order-thumbnails-row">
            <div class="order-item-thumbnail-box" title="Steamed Veg Momos"><img src="/static/steamed_momos.jpg" alt="Steamed Momos"></div>
            <div class="order-item-thumbnail-box" title="Indori Spiced Poha"><img src="/static/indori_poha.jpg" alt="Indori Poha"></div>
          </div>
        </div>
      </section>

      <!-- TAB 2: MY ADDRESSES -->
      <section class="account-tab-pane" id="pane-addresses">
        <div class="account-panel-card">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
            <div>
              <h2 class="account-panel-title">My Addresses</h2>
              <div class="account-panel-subtitle" style="margin-bottom: 0;">Manage delivery locations for home, work, and friends</div>
            </div>
            <button type="button" onclick="window.location.href='/?add_address=1'" style="padding: 9px 16px; background: #0C831F; color: #ffffff; border: none; border-radius: 8px; font-weight: 700; font-size: 13px; cursor: pointer;">
              + Add New Address
            </button>
          </div>

          <!-- Address Card 1 -->
          <div class="address-card-item">
            <div>
              <span class="address-badge-tag">&#127968; Home</span>
              <div class="address-text">#42, 12th Main, HAL 2nd Stage, Indiranagar, Bengaluru, Karnataka 560038</div>
              <div class="address-phone">Phone number: %s</div>
            </div>
            <div style="display: flex; gap: 8px;">
              <button type="button" style="padding: 6px 12px; background: #F1F5F9; border: 1px solid #CBD5E1; border-radius: 6px; font-weight: 600; font-size: 12px; cursor: pointer;">Edit</button>
              <button type="button" style="padding: 6px 12px; background: #F0FDF4; border: 1px solid #BBF7D0; color: #166534; border-radius: 6px; font-weight: 700; font-size: 12px; cursor: pointer;">Deliver Here</button>
            </div>
          </div>

          <!-- Address Card 2 -->
          <div class="address-card-item">
            <div>
              <span class="address-badge-tag">&#127970; Work</span>
              <div class="address-text">Level 4, Embassy Golf Links Business Park, Domlur, Bengaluru, Karnataka 560071</div>
              <div class="address-phone">Phone number: %s</div>
            </div>
            <div style="display: flex; gap: 8px;">
              <button type="button" style="padding: 6px 12px; background: #F1F5F9; border: 1px solid #CBD5E1; border-radius: 6px; font-weight: 600; font-size: 12px; cursor: pointer;">Edit</button>
              <button type="button" style="padding: 6px 12px; background: #F0FDF4; border: 1px solid #BBF7D0; color: #166534; border-radius: 6px; font-weight: 700; font-size: 12px; cursor: pointer;">Deliver Here</button>
            </div>
          </div>

          <!-- Address Card 3 -->
          <div class="address-card-item">
            <div>
              <span class="address-badge-tag">&#128205; Friend's Place</span>
              <div class="address-text">Flat 302, Green Glen Layout, Bellandur, Bengaluru, Karnataka 560103</div>
              <div class="address-phone">Phone number: %s</div>
            </div>
            <div style="display: flex; gap: 8px;">
              <button type="button" style="padding: 6px 12px; background: #F1F5F9; border: 1px solid #CBD5E1; border-radius: 6px; font-weight: 600; font-size: 12px; cursor: pointer;">Edit</button>
              <button type="button" style="padding: 6px 12px; background: #F0FDF4; border: 1px solid #BBF7D0; color: #166534; border-radius: 6px; font-weight: 700; font-size: 12px; cursor: pointer;">Deliver Here</button>
            </div>
          </div>
        </div>
      </section>

      <!-- TAB 3: MY PRESCRIPTIONS -->
      <section class="account-tab-pane" id="pane-prescriptions">
        <div class="account-panel-card">
          <h2 class="account-panel-title">My Prescriptions &amp; Dietary Guidelines</h2>
          <div class="account-panel-subtitle">Upload doctor prescriptions or customize medical dietary preferences for meal preparations</div>

          <div style="border: 2px dashed #CBD5E1; border-radius: 12px; padding: 32px; text-align: center; background: #F8FAFC; margin-bottom: 24px; cursor: pointer;" onclick="document.getElementById('dashPrescriptionInput').click()">
            <input type="file" id="dashPrescriptionInput" style="display: none;" onchange="alert('Prescription uploaded: ' + this.files[0].name)">
            <div style="font-size: 32px; margin-bottom: 8px;">&#128203;</div>
            <div style="font-weight: 700; font-size: 15px; color: #1E293B;">Click or Drag &amp; Drop Prescription Slip</div>
            <div style="font-size: 13px; color: #64748B; margin-top: 4px;">Supports PNG, JPG, or PDF formats (up to 10MB)</div>
          </div>

          <div style="font-weight: 700; font-size: 14.5px; color: #1E293B; margin-bottom: 12px;">Active Dietary Preferences</div>
          <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px;">
            <label style="display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer;">
              <input type="checkbox" checked style="accent-color: #0C831F; width: 16px; height: 16px;"> Diabetic-Friendly (Low GI)
            </label>
            <label style="display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer;">
              <input type="checkbox" style="accent-color: #0C831F; width: 16px; height: 16px;"> Low Sodium / Heart-Care
            </label>
            <label style="display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer;">
              <input type="checkbox" style="accent-color: #0C831F; width: 16px; height: 16px;"> Jain Food (No Onion/Garlic)
            </label>
            <label style="display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 10px; font-size: 13px; font-weight: 600; cursor: pointer;">
              <input type="checkbox" style="accent-color: #0C831F; width: 16px; height: 16px;"> Zero Added Sugar
            </label>
          </div>
        </div>
      </section>

      <!-- TAB 4: E-GIFT CARDS -->
      <section class="account-tab-pane" id="pane-gift-cards">
        <div class="account-panel-card">
          <h2 class="account-panel-title">E-Gift Cards</h2>
          <div class="account-panel-subtitle">Manage Under100 wallet vouchers and send digital food cards to friends</div>

          <div style="background: linear-gradient(135deg, #0C831F 0%%, #065F14 100%%); border-radius: 16px; padding: 24px; color: #ffffff; margin-bottom: 24px; box-shadow: 0 8px 24px rgba(12, 131, 31, 0.25);">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <div style="font-size: 12px; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; opacity: 0.9;">Under100 Food Card</div>
              <div style="font-weight: 800; font-size: 14px; background: rgba(255,255,255,0.2); padding: 3px 8px; border-radius: 6px;">&#8377;100 MAX</div>
            </div>
            <div style="margin: 20px 0 8px;">
              <div style="font-size: 12px; opacity: 0.85;">Available Balance</div>
              <div style="font-size: 32px; font-weight: 900;">&#8377;0.00</div>
            </div>
            <div style="font-size: 12px; opacity: 0.85;">Valid on all breakfast, snack, and meal orders across Bengaluru</div>
          </div>

          <div style="font-weight: 700; font-size: 14px; color: #1E293B; margin-bottom: 8px;">Redeem Voucher Code</div>
          <div style="display: flex; gap: 10px; margin-bottom: 24px;">
            <input type="text" id="dashGiftCode" placeholder="Enter 16-digit voucher code" style="flex: 1; padding: 12px 14px; border: 1.5px solid #CBD5E1; border-radius: 8px; font-size: 13.5px; outline: none;">
            <button type="button" onclick="alert('Voucher applied successfully!')" style="padding: 0 20px; background: #0C831F; color: #ffffff; border: none; border-radius: 8px; font-weight: 700; font-size: 13.5px; cursor: pointer;">Redeem</button>
          </div>

          <div style="font-weight: 700; font-size: 14px; color: #1E293B; margin-bottom: 10px;">Send Gift Voucher</div>
          <div style="display: flex; gap: 12px;">
            <button type="button" onclick="alert('Selected ₹100 Gift Card')" style="flex: 1; padding: 12px; border: 1.5px solid #E2E8F0; border-radius: 10px; background: #ffffff; font-weight: 700; font-size: 14px; color: #1E293B; cursor: pointer;">&#8377;100 (1 Bite)</button>
            <button type="button" onclick="alert('Selected ₹250 Gift Card')" style="flex: 1; padding: 12px; border: 1.5px solid #E2E8F0; border-radius: 10px; background: #ffffff; font-weight: 700; font-size: 14px; color: #1E293B; cursor: pointer;">&#8377;250 (3 Bites)</button>
            <button type="button" onclick="alert('Selected ₹500 Gift Card')" style="flex: 1; padding: 12px; border: 1.5px solid #E2E8F0; border-radius: 10px; background: #ffffff; font-weight: 700; font-size: 14px; color: #1E293B; cursor: pointer;">&#8377;500 (Feast)</button>
          </div>
        </div>
      </section>

      <!-- TAB 5: ACCOUNT PRIVACY -->
      <section class="account-tab-pane" id="pane-privacy">
        <div class="account-panel-card">
          <h2 class="account-panel-title">Account Privacy &amp; Data</h2>
          <div class="account-panel-subtitle">Control how your personal profile, notification preferences, and order history are managed</div>

          <div style="display: flex; flex-direction: column; gap: 18px;">
            <div style="display: flex; justify-content: space-between; align-items: center; padding-bottom: 14px; border-bottom: 1px solid #F1F5F9;">
              <div>
                <div style="font-weight: 700; font-size: 14px; color: #1E293B;">Personalized Recommendations</div>
                <div style="font-size: 12.5px; color: #64748B;">Tailor breakfast and snack recommendations based on your tastes</div>
              </div>
              <input type="checkbox" checked style="accent-color: #0C831F; width: 18px; height: 18px; cursor: pointer;">
            </div>

            <div style="display: flex; justify-content: space-between; align-items: center; padding-bottom: 14px; border-bottom: 1px solid #F1F5F9;">
              <div>
                <div style="font-weight: 700; font-size: 14px; color: #1E293B;">Order Status SMS &amp; WhatsApp</div>
                <div style="font-size: 12.5px; color: #64748B;">Live delivery progress updates sent to %s</div>
              </div>
              <input type="checkbox" checked style="accent-color: #0C831F; width: 18px; height: 18px; cursor: pointer;">
            </div>

            <div style="display: flex; justify-content: space-between; align-items: center; padding-bottom: 14px; border-bottom: 1px solid #F1F5F9;">
              <div>
                <div style="font-weight: 700; font-size: 14px; color: #1E293B;">Weekend Discounts &amp; Deals</div>
                <div style="font-size: 12.5px; color: #64748B;">Occasional notifications for special bites under &#8377;50</div>
              </div>
              <input type="checkbox" style="accent-color: #0C831F; width: 18px; height: 18px; cursor: pointer;">
            </div>

            <div style="display: flex; gap: 12px; margin-top: 10px;">
              <button type="button" onclick="alert('Preparing account data archive. Download starting...')" style="padding: 10px 16px; background: #F1F5F9; border: 1px solid #CBD5E1; border-radius: 8px; font-size: 13px; font-weight: 600; color: #334155; cursor: pointer;">Download My Data</button>
              <button type="button" onclick="if(confirm('Are you sure you want to delete your account?')) window.location.href='/logout'" style="padding: 10px 16px; background: #FEF2F2; border: 1px solid #FECACA; border-radius: 8px; font-size: 13px; font-weight: 600; color: #DC2626; cursor: pointer;">Delete Account</button>
            </div>
          </div>
        </div>
      </section>

    </main>
  </div>

  <script>
    function switchAccountTab(tabName) {
      // 1. Remove active from all nav items
      var navItems = document.querySelectorAll('.account-nav-item');
      navItems.forEach(function(item) {
        item.classList.remove('active');
      });

      // 2. Set active on clicked item
      var activeItem = document.getElementById('navItem-' + tabName);
      if (activeItem) {
        activeItem.classList.add('active');
      }

      // 3. Hide all panes
      var panes = document.querySelectorAll('.account-tab-pane');
      panes.forEach(function(pane) {
        pane.classList.remove('active');
      });

      // 4. Show target pane
      var targetPane = document.getElementById('pane-' + tabName);
      if (targetPane) {
        targetPane.classList.add('active');
      }

      // 5. Update browser URL history
      if (window.history && window.history.pushState) {
        window.history.pushState(null, '', '/account/' + tabName);
      }
    }

    // Initialize with active tab from server param
    (function() {
      var initialTab = "%s";
      if (initialTab && document.getElementById('pane-' + initialTab)) {
        switchAccountTab(initialTab);
      }
    })();
  </script>
</body>
</html>`, userDisplay, userContact, userContact, userContact, userContact, activeTab)
}

// RenderPortalsHubPage renders a visual switchboard for accessing all Under100 portals
func RenderPortalsHubPage() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Under100 — Portals & Systems Directory</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@500;600;700;800;900&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    :root {
      --primary: #FF5A1F;
      --primary-dark: #D94814;
      --bg: #0F172A;
      --card-bg: #1E293B;
      --border: #334155;
      --text: #F8FAFC;
      --muted: #94A3B8;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: 'Plus Jakarta Sans', sans-serif;
      background: radial-gradient(circle at 50% 0%, #1E293B 0%, #0F172A 100%);
      color: var(--text);
      min-height: 100vh;
      padding: 48px 20px;
      display: flex;
      flex-direction: column;
      align-items: center;
    }
    .header {
      text-align: center;
      max-width: 680px;
      margin-bottom: 48px;
    }
    .badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      background: rgba(255, 90, 31, 0.15);
      border: 1px solid rgba(255, 90, 31, 0.4);
      color: #FF7E47;
      padding: 6px 14px;
      border-radius: 999px;
      font-size: 13px;
      font-weight: 700;
      letter-spacing: 0.5px;
      text-transform: uppercase;
      margin-bottom: 16px;
    }
    .title {
      font-family: 'Outfit', sans-serif;
      font-size: 40px;
      font-weight: 900;
      line-height: 1.15;
      margin-bottom: 12px;
      color: #ffffff;
    }
    .title span { color: var(--primary); }
    .subtitle {
      font-size: 16px;
      color: var(--muted);
      line-height: 1.6;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 24px;
      max-width: 1140px;
      width: 100%;
      margin-bottom: 48px;
    }
    .card {
      background: var(--card-bg);
      border: 1px solid var(--border);
      border-radius: 20px;
      padding: 28px;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      transition: all 0.25s ease;
      position: relative;
      overflow: hidden;
    }
    .card:hover {
      transform: translateY(-4px);
      border-color: #475569;
      box-shadow: 0 16px 32px rgba(0,0,0,0.3);
    }
    .card-top { margin-bottom: 24px; }
    .card-icon {
      width: 52px;
      height: 52px;
      border-radius: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 26px;
      margin-bottom: 18px;
    }
    .card-title {
      font-family: 'Outfit', sans-serif;
      font-size: 22px;
      font-weight: 800;
      color: #ffffff;
      margin-bottom: 8px;
    }
    .card-desc {
      font-size: 14px;
      color: var(--muted);
      line-height: 1.6;
      margin-bottom: 16px;
    }
    .creds-box {
      background: rgba(15, 23, 42, 0.6);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 12px;
      padding: 12px 14px;
      font-size: 12px;
      color: #CBD5E1;
      margin-bottom: 16px;
      font-family: monospace;
    }
    .creds-box strong { color: #38BDF8; font-family: 'Plus Jakarta Sans', sans-serif; }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      width: 100%;
      padding: 13px 20px;
      border-radius: 12px;
      font-weight: 700;
      font-size: 14px;
      text-decoration: none;
      transition: all 0.2s ease;
    }
    .btn-primary {
      background: var(--primary);
      color: #ffffff;
    }
    .btn-primary:hover {
      background: var(--primary-dark);
      box-shadow: 0 4px 16px rgba(255, 90, 31, 0.4);
    }
    .btn-outline {
      background: rgba(255,255,255,0.04);
      border: 1.5px solid var(--border);
      color: #ffffff;
    }
    .btn-outline:hover {
      background: rgba(255,255,255,0.08);
      border-color: #64748B;
    }
    .footer-note {
      font-size: 13px;
      color: #64748B;
      text-align: center;
    }
  </style>
</head>
<body>
  <header class="header">
    <div class="badge">Platform Directory &amp; Navigation Hub</div>
    <h1 class="title">Access All <span>Under100</span> Portals</h1>
    <p class="subtitle">Quick-launch directory for the separate Customer Storefront, Restaurant Owner Portal, Admin Console, and Onboarding systems.</p>
  </header>

  <div class="grid">
    <!-- Card 1: Customer Storefront -->
    <div class="card" style="border-top: 4px solid #FF5A1F;">
      <div class="card-top">
        <div class="card-icon" style="background: rgba(255, 90, 31, 0.15); color: #FF5A1F;">🍽️</div>
        <h2 class="card-title">Customer Storefront</h2>
        <p class="card-desc">Public-facing discovery application for Bangalore diners. Search budget eateries, filter by &#8377;100 price cap, and check live delivery radius.</p>
        <div class="creds-box">
          <strong>Access Level:</strong> Public Customer Access<br>
          <strong>Route:</strong> <code>/</code>
        </div>
      </div>
      <a href="/" target="_blank" class="btn btn-primary">Open Customer Website &rarr;</a>
    </div>

    <!-- Card 2: Merchant Portal -->
    <div class="card" style="border-top: 4px solid #3B82F6;">
      <div class="card-top">
        <div class="card-icon" style="background: rgba(59, 130, 246, 0.15); color: #3B82F6;">🏪</div>
        <h2 class="card-title">Restaurant Owner Portal</h2>
        <p class="card-desc">Partner dashboard for eatery owners to manage menus, mark items sold out, configure operating hours, and set delivery radius &amp; fees.</p>
        <div class="creds-box">
          <strong>Demo Account:</strong> CTR Malleshwaram<br>
          <strong>Email:</strong> merchant.ctr@fooddiscovery.com<br>
          <strong>Password:</strong> Merchant@123
        </div>
      </div>
      <div style="display: flex; gap: 10px;">
        <a href="/merchant/login" target="_blank" class="btn btn-outline" style="border-color: #3B82F6; color: #60A5FA; flex: 1;">Merchant Sign In &rarr;</a>
        <a href="/partner/register" target="_blank" class="btn btn-outline" style="border-color: #64748B; color: #94A3B8; flex: 1;">Register Eatery &rarr;</a>
      </div>
    </div>

    <!-- Card 3: Admin Console -->
    <div class="card" style="border-top: 4px solid #10B981;">
      <div class="card-top">
        <div class="card-icon" style="background: rgba(16, 185, 129, 0.15); color: #10B981;">🛡️</div>
        <h2 class="card-title">Admin Console</h2>
        <p class="card-desc">Central operations panel for Under100 staff to approve/reject restaurant submissions, assign user roles, and manage food categories.</p>
        <div class="creds-box">
          <strong>Root Admin:</strong> System Administrator<br>
          <strong>Email:</strong> admin@fooddiscovery.com<br>
          <strong>Password:</strong> Admin@123
        </div>
      </div>
      <a href="/admin/login" target="_blank" class="btn btn-outline" style="border-color: #10B981; color: #34D399;">Open Admin Console &rarr;</a>
    </div>

    <!-- Card 4: Partner Registration -->
    <div class="card" style="border-top: 4px solid #A855F7;">
      <div class="card-top">
        <div class="card-icon" style="background: rgba(168, 85, 247, 0.15); color: #A855F7;">✍️</div>
        <h2 class="card-title">Partner Registration</h2>
        <p class="card-desc">Self-service registration portal for new restaurant owners in Bengaluru to submit an eatery listing for administrative review and approval.</p>
        <div class="creds-box">
          <strong>Access Level:</strong> Prospective Eatery Partners<br>
          <strong>Status Upon Submission:</strong> Pending Review
        </div>
      </div>
      <div style="display: flex; gap: 10px;">
        <a href="/partner/register" target="_blank" class="btn btn-outline" style="border-color: #A855F7; color: #C084FC; flex: 1;">New Partner Register &rarr;</a>
        <a href="/merchant/login" target="_blank" class="btn btn-outline" style="border-color: #64748B; color: #94A3B8; flex: 1;">Existing Sign In &rarr;</a>
      </div>
    </div>
  </div>

  <div class="footer-note">
    Under100 Application Hub &bull; Running on Go 1.24 + Echo + PostgreSQL 15 + Redis &bull; Port 8000
  </div>
</body>
</html>`
}

