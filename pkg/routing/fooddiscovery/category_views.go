package fooddiscovery

import (
	"fmt"
	"strings"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// RenderCategoryChips renders responsive category chips for discovery filtering
func RenderCategoryChips(categories []fooddiscovery.Category, activeSlug string) string {
	var sb strings.Builder
	sb.WriteString(`<div class="category-chips-bar" style="display: flex; gap: 8px; overflow-x: auto; padding: 10px 0; margin-bottom: 20px; -webkit-overflow-scrolling: touch;" class="hide-scrollbar">`)

	allActive := ""
	if activeSlug == "" || activeSlug == "all" {
		allActive = "active"
	}

	sb.WriteString(fmt.Sprintf(`
		<button 
			type="button" 
			class="cat-chip-pill %s" 
			onclick="applyCategoryFilter('all', this)"
			style="white-space: nowrap; padding: 7px 16px; border-radius: 9999px; font-weight: 700; font-size: 0.85rem; border: 1.5px solid #cbd5e1; background: #ffffff; color: #1e293b; cursor: pointer; transition: all 0.15s ease;"
		>
			All Dishes
		</button>
	`, allActive))

	for _, cat := range categories {
		if !cat.IsActive {
			continue
		}
		activeClass := ""
		if strings.ToLower(cat.Slug) == strings.ToLower(activeSlug) {
			activeClass = "active"
		}
		sb.WriteString(fmt.Sprintf(`
			<button 
				type="button" 
				class="cat-chip-pill %s" 
				onclick="applyCategoryFilter('%s', this)"
				data-slug="%s"
				style="white-space: nowrap; padding: 7px 16px; border-radius: 9999px; font-weight: 700; font-size: 0.85rem; border: 1.5px solid #cbd5e1; background: #ffffff; color: #1e293b; cursor: pointer; transition: all 0.15s ease; display: inline-flex; align-items: center; gap: 6px;"
			>
				<span>%s</span>
				<span>%s</span>
			</button>
		`, activeClass, cat.Slug, cat.Slug, cat.Icon, escapeHTML(cat.Name)))
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// RenderAdminCategoriesPage renders the admin categories management dashboard
func RenderAdminCategoriesPage(
	categories []fooddiscovery.Category,
	currentUser *ent.User,
	msg, errMsg string,
) string {
	var sb strings.Builder

	totalCount := len(categories)
	activeCount := 0
	inactiveCount := 0
	totalDishes := 0

	for _, c := range categories {
		if c.IsActive {
			activeCount++
		} else {
			inactiveCount++
		}
		totalDishes += c.DishCount
	}

	sb.WriteString(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Food Categories & Cuisines — Bite100 Admin</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <link rel="stylesheet" href="/static/css/main.css">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8fafc; color: #0f172a; margin: 0; padding: 0; }
    .admin-nav { background: #ffffff; border-bottom: 1px solid #e2e8f0; padding: 0.75rem 2rem; display: flex; justify-content: space-between; align-items: center; }
    .admin-title { font-size: 1.25rem; font-weight: 800; color: #e23744; }
    .nav-tab { color: #64748b; text-decoration: none; font-weight: 600; font-size: 0.95rem; padding: 0.5rem 1rem; border-radius: 8px; }
    .nav-tab.active { color: #0f172a; background: #f1f5f9; }
    .container { max-width: 1300px; margin: 2rem auto; padding: 0 1.5rem; }
    .metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; margin-bottom: 1.5rem; }
    .metric-card { background: #ffffff; border: 1px solid #e2e8f0; padding: 1.25rem; border-radius: 12px; }
    .metric-num { font-size: 1.75rem; font-weight: 800; color: #0f172a; margin-top: 0.25rem; }
    table { width: 100%%; border-collapse: collapse; background: #ffffff; border-radius: 12px; overflow: hidden; border: 1px solid #e2e8f0; }
    th { background: #f8fafc; padding: 0.75rem 1rem; text-align: left; font-size: 0.8rem; font-weight: 700; color: #64748b; border-bottom: 1px solid #e2e8f0; }
    td { padding: 0.85rem 1rem; border-bottom: 1px solid #f1f5f9; font-size: 0.88rem; vertical-align: middle; }
    .badge { display: inline-block; padding: 0.2rem 0.55rem; border-radius: 6px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; }
    .badge-active { background: #dcfce7; color: #166534; }
    .badge-inactive { background: #fee2e2; color: #991b1b; }
    .badge-type { background: #f1f5f9; color: #475569; }
    .btn-action { padding: 0.35rem 0.75rem; border-radius: 6px; font-size: 0.8rem; font-weight: 700; cursor: pointer; border: 1px solid transparent; }
    .create-card { background: #ffffff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1.5rem; margin-bottom: 2rem; }
  </style>
</head>
<body>
  <div class="admin-nav">
    <div style="display: flex; align-items: center; gap: 2rem;">
      <div class="admin-title">Bite100 Admin Dashboard</div>
      <div style="display: flex; gap: 0.5rem;">
        <a href="/admin" class="nav-tab">Restaurants</a>
        <a href="/admin/users" class="nav-tab">Users</a>
        <a href="/admin/reviews" class="nav-tab">Review Moderation</a>
        <a href="/admin/categories" class="nav-tab active">Categories</a>
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
	`, escapeHTML(currentUser.Name)))

	if msg != "" {
		sb.WriteString(fmt.Sprintf(`
			<div style="background: #dcfce7; border: 1px solid #86efac; color: #166534; padding: 0.85rem 1.25rem; border-radius: 8px; margin-bottom: 1.5rem; font-weight: 600;">
				%s
			</div>
		`, escapeHTML(msg)))
	}
	if errMsg != "" {
		sb.WriteString(fmt.Sprintf(`
			<div style="background: #fee2e2; border: 1px solid #fca5a5; color: #991b1b; padding: 0.85rem 1.25rem; border-radius: 8px; margin-bottom: 1.5rem; font-weight: 600;">
				%s
			</div>
		`, escapeHTML(errMsg)))
	}

	sb.WriteString(fmt.Sprintf(`
    <div class="metrics-grid">
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Total Categories</div>
        <div class="metric-num">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Active Categories</div>
        <div class="metric-num" style="color: #166534;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Inactive Categories</div>
        <div class="metric-num" style="color: #991b1b;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Total Dish Mappings</div>
        <div class="metric-num" style="color: #2563eb;">%d</div>
      </div>
    </div>

    <!-- Create New Category Card -->
    <div class="create-card">
      <h3 style="margin-top: 0; margin-bottom: 1rem; font-size: 1.1rem; font-weight: 800;">Create Food Category / Cuisine</h3>
      <form method="POST" action="/admin/categories" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)) auto; gap: 0.75rem; align-items: flex-end;">
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; margin-bottom: 0.3rem;">Category Name *</label>
          <input type="text" name="name" placeholder="e.g. Continental" required style="width: 100%%; box-sizing: border-box; padding: 0.5rem 0.75rem; border: 1px solid #cbd5e1; border-radius: 6px;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; margin-bottom: 0.3rem;">Type *</label>
          <select name="category_type" style="width: 100%%; box-sizing: border-box; padding: 0.5rem 0.75rem; border: 1px solid #cbd5e1; border-radius: 6px;">
            <option value="cuisine">Cuisine</option>
            <option value="meal_type">Meal Type</option>
            <option value="food_type">Food Type</option>
          </select>
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; margin-bottom: 0.3rem;">Icon (Emoji)</label>
          <input type="text" name="icon" value="🍽️" placeholder="🍽️" style="width: 100%%; box-sizing: border-box; padding: 0.5rem 0.75rem; border: 1px solid #cbd5e1; border-radius: 6px;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; margin-bottom: 0.3rem;">Description</label>
          <input type="text" name="description" placeholder="Brief description..." style="width: 100%%; box-sizing: border-box; padding: 0.5rem 0.75rem; border: 1px solid #cbd5e1; border-radius: 6px;">
        </div>
        <div>
          <button type="submit" style="background: #e23744; color: white; border: none; padding: 0.55rem 1.25rem; border-radius: 6px; font-weight: 700; cursor: pointer;">
            + Add Category
          </button>
        </div>
      </form>
    </div>

    <table>
      <thead>
        <tr>
          <th style="width: 50px;">ID</th>
          <th style="width: 180px;">Category</th>
          <th style="width: 140px;">Slug</th>
          <th style="width: 120px;">Type</th>
          <th>Description</th>
          <th style="width: 90px;">Dishes</th>
          <th style="width: 100px;">Status</th>
          <th style="width: 140px;">Action</th>
        </tr>
      </thead>
      <tbody>
	`, totalCount, activeCount, inactiveCount, totalDishes))

	for _, cat := range categories {
		sb.WriteString(RenderAdminCategoryRow(cat))
	}

	sb.WriteString(`
      </tbody>
    </table>
  </div>
</body>
</html>
	`)

	return sb.String()
}

// RenderAdminCategoryRow renders a single category row
func RenderAdminCategoryRow(cat fooddiscovery.Category) string {
	statusBadge := `<span class="badge badge-active">Active</span>`
	toggleActionText := "Deactivate"
	toggleStyle := "background: #fee2e2; color: #991b1b; border: 1px solid #fca5a5;"
	if !cat.IsActive {
		statusBadge = `<span class="badge badge-inactive">Inactive</span>`
		toggleActionText = "Activate"
		toggleStyle = "background: #dcfce7; color: #166534; border: 1px solid #86efac;"
	}

	typeDisplay := strings.Title(strings.ReplaceAll(cat.CategoryType, "_", " "))

	return fmt.Sprintf(`
		<tr id="admin-cat-row-%d">
			<td style="font-weight: 700; color: #64748b;">#%d</td>
			<td>
				<div style="font-weight: 700; display: flex; align-items: center; gap: 0.4rem;">
					<span style="font-size: 1.2rem;">%s</span>
					<span>%s</span>
				</div>
			</td>
			<td style="font-family: monospace; font-size: 0.82rem; color: #475569;">%s</td>
			<td><span class="badge badge-type">%s</span></td>
			<td style="color: #64748b; font-size: 0.85rem;">%s</td>
			<td style="font-weight: 800; color: #0f172a;">%d</td>
			<td>%s</td>
			<td>
				<form hx-post="/admin/categories/%d/toggle" hx-target="#admin-cat-row-%d" hx-swap="outerHTML" style="display: inline;">
					<button type="submit" class="btn-action" style="%s">%s</button>
				</form>
			</td>
		</tr>
	`, cat.ID, cat.ID, cat.Icon, escapeHTML(cat.Name), cat.Slug, typeDisplay, escapeHTML(cat.Description), cat.DishCount, statusBadge, cat.ID, cat.ID, toggleStyle, toggleActionText)
}
