package fooddiscovery

import (
	"fmt"
	"strings"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/repos/fooddiscovery"
)

// renderStars outputs visual unicode star rating
func renderStars(rating int) string {
	var sb strings.Builder
	for i := 1; i <= 5; i++ {
		if i <= rating {
			sb.WriteString("★")
		} else {
			sb.WriteString("☆")
		}
	}
	return sb.String()
}

// RenderReviewsSectionHTML renders the complete reviews and rating section for the restaurant modal
func RenderReviewsSectionHTML(
	restaurantID int,
	stats fooddiscovery.RestaurantRatingStats,
	userReview *fooddiscovery.Review,
	reviews fooddiscovery.PaginatedReviews,
	currentUser *ent.User,
) string {
	var sb strings.Builder

	sb.WriteString(`<div id="reviews-section" style="margin-top: 2rem; border-top: 1.5px solid #e2e8f0; padding-top: 1.5rem;">`)

	// Section Header
	sb.WriteString(`
		<div style="display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 1.25rem;">
			<div>
				<h3 style="font-family: var(--font-heading); font-size: 1.25rem; font-weight: 800; color: var(--text-dark); margin: 0 0 0.25rem 0;">
					Customer Ratings & Reviews
				</h3>
				<p style="color: var(--text-muted); font-size: 0.85rem; margin: 0;">Verified community feedback for budget dining</p>
			</div>
	`)

	// Overall Rating Summary Box
	if stats.ReviewCount > 0 {
		reviewWord := "reviews"
		if stats.ReviewCount == 1 {
			reviewWord = "review"
		}
		sb.WriteString(fmt.Sprintf(`
			<div style="display: flex; align-items: center; gap: 0.75rem; background: #fefce8; border: 1.5px solid #fef08a; padding: 0.5rem 0.85rem; border-radius: 10px;">
				<div style="font-size: 1.75rem; font-weight: 800; color: #854d0e; line-height: 1;">%.1f</div>
				<div>
					<div style="color: #f59e0b; font-size: 0.95rem; line-height: 1.1;">%s</div>
					<div style="font-size: 0.75rem; color: #713f12; font-weight: 600; margin-top: 2px;">%d %s</div>
				</div>
			</div>
		`, stats.AverageRating, renderStars(int(stats.AverageRating+0.5)), stats.ReviewCount, reviewWord))
	} else {
		sb.WriteString(`
			<div style="background: #f8fafc; border: 1px solid #e2e8f0; padding: 0.45rem 0.8rem; border-radius: 8px; text-align: right;">
				<div style="font-size: 0.85rem; font-weight: 700; color: #64748b;">★ New</div>
				<div style="font-size: 0.72rem; color: #94a3b8;">No reviews yet</div>
			</div>
		`)
	}
	sb.WriteString(`</div>`)

	// Customer Review Form or "Your Review" Box
	if currentUser != nil {
		if userReview != nil {
			// Display the authenticated customer's existing review with Edit / Delete actions
			sb.WriteString(fmt.Sprintf(`
				<div id="user-review-box" style="background: #f0fdf4; border: 1.5px solid #86efac; border-radius: 12px; padding: 1.1rem 1.25rem; margin-bottom: 1.5rem;">
					<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
						<div style="display: flex; align-items: center; gap: 0.5rem;">
							<span style="background: #166534; color: white; font-size: 0.7rem; font-weight: 800; padding: 0.15rem 0.45rem; border-radius: 4px; text-transform: uppercase;">Your Review</span>
							<span style="color: #f59e0b; font-weight: 800; font-size: 0.95rem;">%s (%d/5)</span>
						</div>
						<div style="display: flex; gap: 0.5rem;">
							<button 
								type="button" 
								class="btn-action" 
								style="color: #2563eb; font-weight: 700; font-size: 0.8rem; background: none; border: none; cursor: pointer; text-decoration: underline;" 
								hx-get="/reviews/%d/edit" 
								hx-target="#user-review-box"
							>
								Edit
							</button>
							<button 
								type="button" 
								class="btn-action" 
								style="color: #b91c1c; font-weight: 700; font-size: 0.8rem; background: none; border: none; cursor: pointer; text-decoration: underline;" 
								hx-delete="/reviews/%d" 
								hx-target="#reviews-section" 
								hx-confirm="Are you sure you want to delete your review?"
							>
								Delete
							</button>
						</div>
					</div>
					<p style="color: #1e293b; font-size: 0.9rem; line-height: 1.5; margin: 0;">%s</p>
					<div style="font-size: 0.72rem; color: #166534; font-weight: 600; margin-top: 0.4rem;">
						Reviewed on %s
					</div>
				</div>
			`, renderStars(userReview.Rating), userReview.Rating, userReview.ID, userReview.ID, escapeHTML(userReview.Comment), userReview.CreatedAt.Format("02 Jan 2006")))
		} else {
			// Display the interactive Review Submission Form
			sb.WriteString(fmt.Sprintf(`
				<form 
					id="review-submit-form" 
					hx-post="/restaurant/%d/reviews" 
					hx-target="#reviews-section" 
					style="background: #f8fafc; border: 1.5px solid #e2e8f0; border-radius: 12px; padding: 1.1rem 1.25rem; margin-bottom: 1.5rem;"
				>
					<div style="font-weight: 700; font-size: 0.95rem; color: var(--text-dark); margin-bottom: 0.75rem;">
						Rate & Review this Restaurant
					</div>
					<div id="review-form-error"></div>
					<div style="margin-bottom: 0.75rem;">
						<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.35rem;">
							Your Rating (1 to 5 Stars)
						</label>
						<div style="display: flex; gap: 0.65rem; align-items: center; flex-wrap: wrap;">
							<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
								<input type="radio" name="rating" value="5" checked> ★★★★★ (5)
							</label>
							<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
								<input type="radio" name="rating" value="4"> ★★★★ (4)
							</label>
							<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
								<input type="radio" name="rating" value="3"> ★★★ (3)
							</label>
							<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
								<input type="radio" name="rating" value="2"> ★★ (2)
							</label>
							<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
								<input type="radio" name="rating" value="1"> ★ (1)
							</label>
						</div>
					</div>
					<div style="margin-bottom: 0.75rem;">
						<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.35rem;">
							Your Review (Maximum 1000 characters)
						</label>
						<textarea 
							name="comment" 
							rows="3" 
							maxlength="1000" 
							placeholder="What food items were best under ₹100? How was the value, cleanliness, and taste?" 
							style="width: 100%%; box-sizing: border-box; border: 1.5px solid #cbd5e1; border-radius: 8px; padding: 0.6rem 0.75rem; font-family: inherit; font-size: 0.88rem; resize: vertical;" 
							required
						></textarea>
					</div>
					<div style="display: flex; justify-content: flex-end;">
						<button 
							type="submit" 
							style="background: #e23744; color: white; border: none; padding: 0.55rem 1.35rem; border-radius: 8px; font-weight: 700; font-size: 0.88rem; cursor: pointer;"
						>
							Submit Review
						</button>
					</div>
				</form>
			`, restaurantID))
		}
	} else {
		// Logged out prompt
		sb.WriteString(`
			<div style="background: #fef2f2; border: 1px solid #fecaca; border-radius: 12px; padding: 0.85rem 1.25rem; margin-bottom: 1.5rem; display: flex; justify-content: space-between; align-items: center;">
				<div>
					<div style="font-weight: 700; color: #991b1b; font-size: 0.9rem;">Dined here recently?</div>
					<div style="font-size: 0.8rem; color: #7f1d1d;">Share your food review and help other Bengaluru foodies!</div>
				</div>
				<a href="/login?msg=auth_required" style="background: #e23744; color: white; padding: 0.45rem 1rem; border-radius: 8px; text-decoration: none; font-weight: 700; font-size: 0.82rem;">
					Log in to Review
				</a>
			</div>
		`)
	}

	// Reviews List Container
	sb.WriteString(`<div id="reviews-list">`)
	sb.WriteString(RenderReviewsListHTML(restaurantID, reviews, currentUser))
	sb.WriteString(`</div>`)

	sb.WriteString(`</div>`) // end #reviews-section
	return sb.String()
}

// RenderReviewsListHTML renders the inner reviews list and pagination
func RenderReviewsListHTML(
	restaurantID int,
	reviews fooddiscovery.PaginatedReviews,
	currentUser *ent.User,
) string {
	var sb strings.Builder

	if len(reviews.Items) == 0 {
		return `
			<div style="text-align: center; padding: 2rem 1rem; background: #f8fafc; border: 1px dashed #cbd5e1; border-radius: 12px; color: #64748b; font-size: 0.9rem;">
				No customer reviews yet. Be the first to share your experience!
			</div>
		`
	}

	for _, rev := range reviews.Items {
		isAuthor := currentUser != nil && currentUser.ID == rev.UserID

		reportAction := ""
		if currentUser != nil && !isAuthor {
			reportAction = fmt.Sprintf(`
				<div id="report-status-%d">
					<button 
						type="button" 
						class="btn-report" 
						style="background: none; border: none; color: #94a3b8; font-size: 0.75rem; cursor: pointer; text-decoration: underline;" 
						hx-post="/reviews/%d/report" 
						hx-target="#report-status-%d" 
						hx-confirm="Report this review as inappropriate?"
						title="Report this review"
					>
						Report
					</button>
				</div>
			`, rev.ID, rev.ID, rev.ID)
		}

		authorBadge := ""
		if isAuthor {
			authorBadge = `<span style="background: #dcfce7; color: #166534; font-size: 0.68rem; font-weight: 800; padding: 0.1rem 0.35rem; border-radius: 3px; margin-left: 0.35rem;">YOU</span>`
		}

		sb.WriteString(fmt.Sprintf(`
			<div class="review-item" style="border-bottom: 1px solid #f1f5f9; padding: 0.9rem 0;">
				<div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 0.35rem;">
					<div>
						<div style="display: flex; align-items: center; gap: 0.4rem;">
							<span style="font-weight: 700; font-size: 0.92rem; color: #1e293b;">%s</span>
							%s
						</div>
						<div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 2px;">
							<span style="color: #f59e0b; font-size: 0.85rem; letter-spacing: 1px;">%s</span>
							<span style="color: #94a3b8; font-size: 0.75rem;">%s</span>
						</div>
					</div>
					<div>
						%s
					</div>
				</div>
				<p style="color: #334155; font-size: 0.88rem; line-height: 1.55; margin: 0.35rem 0 0 0;">
					%s
				</p>
			</div>
		`, escapeHTML(rev.UserName), authorBadge, renderStars(rev.Rating), rev.CreatedAt.Format("02 Jan 2006"), reportAction, escapeHTML(rev.Comment)))
	}

	// Pagination
	if reviews.TotalPages > 1 {
		sb.WriteString(`
			<div style="display: flex; justify-content: space-between; align-items: center; margin-top: 1.25rem; padding-top: 0.75rem; border-top: 1px solid #e2e8f0; font-size: 0.82rem; color: #64748b;">
		`)
		sb.WriteString(fmt.Sprintf(`<div>Page %d of %d (%d reviews)</div>`, reviews.Page, reviews.TotalPages, reviews.TotalCount))
		sb.WriteString(`<div style="display: flex; gap: 0.5rem;">`)

		if reviews.HasPrev {
			sb.WriteString(fmt.Sprintf(`
				<button 
					type="button" 
					class="btn-page" 
					style="background: #f1f5f9; border: 1px solid #cbd5e1; padding: 0.3rem 0.65rem; border-radius: 6px; font-weight: 700; cursor: pointer; color: #0f172a;" 
					hx-get="/restaurant/%d/reviews?page=%d" 
					hx-target="#reviews-list"
				>
					&larr; Prev
				</button>
			`, restaurantID, reviews.Page-1))
		}
		if reviews.HasNext {
			sb.WriteString(fmt.Sprintf(`
				<button 
					type="button" 
					class="btn-page" 
					style="background: #f1f5f9; border: 1px solid #cbd5e1; padding: 0.3rem 0.65rem; border-radius: 6px; font-weight: 700; cursor: pointer; color: #0f172a;" 
					hx-get="/restaurant/%d/reviews?page=%d" 
					hx-target="#reviews-list"
				>
					Next &rarr;
				</button>
			`, restaurantID, reviews.Page+1))
		}

		sb.WriteString(`</div></div>`)
	}

	return sb.String()
}

// RenderEditReviewForm renders the inline edit form for a customer's review
func RenderEditReviewForm(rev *fooddiscovery.Review) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
		<form 
			id="review-edit-form" 
			hx-post="/reviews/%d/edit" 
			hx-target="#reviews-section" 
			style="background: #eff6ff; border: 1.5px solid #93c5fd; border-radius: 12px; padding: 1.1rem 1.25rem; margin-bottom: 1.5rem;"
		>
			<div style="font-weight: 800; font-size: 0.95rem; color: #1e3a8a; margin-bottom: 0.75rem;">
				Edit Your Rating & Review
			</div>
			<div style="margin-bottom: 0.75rem;">
				<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #1e3a8a; margin-bottom: 0.35rem;">
					Rating
				</label>
				<div style="display: flex; gap: 0.65rem; align-items: center; flex-wrap: wrap;">
	`, rev.ID))

	for i := 5; i >= 1; i-- {
		checked := ""
		if rev.Rating == i {
			checked = "checked"
		}
		stars := strings.Repeat("★", i)
		sb.WriteString(fmt.Sprintf(`
			<label style="cursor: pointer; display: flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600;">
				<input type="radio" name="rating" value="%d" %s> %s (%d)
			</label>
		`, i, checked, stars, i))
	}

	sb.WriteString(fmt.Sprintf(`
				</div>
			</div>
			<div style="margin-bottom: 0.75rem;">
				<label style="display: block; font-size: 0.8rem; font-weight: 700; color: #1e3a8a; margin-bottom: 0.35rem;">
					Review (Maximum 1000 characters)
				</label>
				<textarea 
					name="comment" 
					rows="3" 
					maxlength="1000" 
					style="width: 100%%; box-sizing: border-box; border: 1.5px solid #93c5fd; border-radius: 8px; padding: 0.6rem 0.75rem; font-family: inherit; font-size: 0.88rem; resize: vertical;" 
					required
				>%s</textarea>
			</div>
			<div style="display: flex; justify-content: flex-end; gap: 0.5rem;">
				<button 
					type="button" 
					style="background: #e2e8f0; color: #475569; border: none; padding: 0.5rem 1rem; border-radius: 8px; font-weight: 700; font-size: 0.85rem; cursor: pointer;" 
					hx-get="/restaurant/%d/reviews?modal=true" 
					hx-target="#reviews-section"
				>
					Cancel
				</button>
				<button 
					type="submit" 
					style="background: #2563eb; color: white; border: none; padding: 0.5rem 1.25rem; border-radius: 8px; font-weight: 700; font-size: 0.85rem; cursor: pointer;"
				>
					Save Changes
				</button>
			</div>
		</form>
	`, escapeHTML(rev.Comment), rev.RestaurantID))

	return sb.String()
}

// RenderAdminReviewsPage renders the administrator moderation page for customer reviews
func RenderAdminReviewsPage(
	reviews fooddiscovery.PaginatedReviews,
	activeFilter string,
	currentUser *ent.User,
	stats map[string]int,
) string {
	var sb strings.Builder

	if activeFilter == "" {
		activeFilter = "all"
	}

	sb.WriteString(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Review Moderation — Bite100 Admin</title>
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
    .tab-pills { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; }
    .tab-pill { padding: 0.4rem 0.9rem; border-radius: 9999px; text-decoration: none; font-weight: 700; font-size: 0.85rem; color: #64748b; background: #e2e8f0; }
    .tab-pill.active { background: #0f172a; color: #ffffff; }
    table { width: 100%%; border-collapse: collapse; background: #ffffff; border-radius: 12px; overflow: hidden; border: 1px solid #e2e8f0; }
    th { background: #f8fafc; padding: 0.75rem 1rem; text-align: left; font-size: 0.8rem; font-weight: 700; color: #64748b; border-bottom: 1px solid #e2e8f0; }
    td { padding: 0.85rem 1rem; border-bottom: 1px solid #f1f5f9; font-size: 0.88rem; vertical-align: top; }
    .badge { display: inline-block; padding: 0.2rem 0.55rem; border-radius: 6px; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; }
    .badge-approved { background: #dcfce7; color: #166534; }
    .badge-hidden { background: #fee2e2; color: #991b1b; }
    .badge-reports { background: #fef08a; color: #854d0e; }
    .btn-action { padding: 0.35rem 0.75rem; border-radius: 6px; font-size: 0.8rem; font-weight: 700; cursor: pointer; border: 1px solid transparent; }
  </style>
</head>
<body>
  <div class="admin-nav">
    <div style="display: flex; align-items: center; gap: 2rem;">
      <div class="admin-title">Bite100 Admin Dashboard</div>
      <div style="display: flex; gap: 0.5rem;">
        <a href="/admin" class="nav-tab">Restaurants</a>
        <a href="/admin/users" class="nav-tab">Users</a>
        <a href="/admin/reviews" class="nav-tab active">Review Moderation</a>
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
        <div style="font-size: 0.82rem; color: #64748b;">Total Reviews</div>
        <div class="metric-num">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Reported by Users</div>
        <div class="metric-num" style="color: #d97706;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Hidden / Moderated</div>
        <div class="metric-num" style="color: #b91c1c;">%d</div>
      </div>
      <div class="metric-card">
        <div style="font-size: 0.82rem; color: #64748b;">Active / Approved</div>
        <div class="metric-num" style="color: #166534;">%d</div>
      </div>
    </div>

    <div class="tab-pills">
      <a href="/admin/reviews" class="tab-pill %s">All Reviews (%d)</a>
      <a href="/admin/reviews?filter=reported" class="tab-pill %s">Reported (%d)</a>
      <a href="/admin/reviews?filter=hidden" class="tab-pill %s">Hidden (%d)</a>
      <a href="/admin/reviews?filter=approved" class="tab-pill %s">Approved (%d)</a>
    </div>

    <table>
      <thead>
        <tr>
          <th style="width: 50px;">ID</th>
          <th style="width: 180px;">Restaurant</th>
          <th style="width: 160px;">Customer</th>
          <th style="width: 100px;">Rating</th>
          <th>Comment</th>
          <th style="width: 100px;">Status</th>
          <th style="width: 90px;">Reports</th>
          <th style="width: 200px;">Moderation Action</th>
        </tr>
      </thead>
      <tbody>
	`,
		escapeHTML(currentUser.Name),
		stats["all"],
		stats["reported"],
		stats["hidden"],
		stats["approved"],
		filterClass(activeFilter, "all"), stats["all"],
		filterClass(activeFilter, "reported"), stats["reported"],
		filterClass(activeFilter, "hidden"), stats["hidden"],
		filterClass(activeFilter, "approved"), stats["approved"],
	))

	if len(reviews.Items) == 0 {
		sb.WriteString(`
			<tr>
				<td colspan="8" style="text-align: center; padding: 2.5rem; color: #64748b;">
					No reviews found under the selected filter.
				</td>
			</tr>
		`)
	} else {
		for _, rev := range reviews.Items {
			sb.WriteString(RenderAdminReviewRow(rev))
		}
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

func filterClass(active, current string) string {
	if strings.ToLower(active) == strings.ToLower(current) {
		return "active"
	}
	return ""
}

// RenderAdminReviewRow renders a single review row in the admin moderation table
func RenderAdminReviewRow(rev fooddiscovery.Review) string {
	statusBadge := `<span class="badge badge-approved">Approved</span>`
	if rev.Status == "hidden" {
		statusBadge = `<span class="badge badge-hidden">Hidden</span>`
	}

	reportBadge := `<span style="color: #64748b; font-size: 0.8rem;">0</span>`
	if rev.ReportCount > 0 {
		reportBadge = fmt.Sprintf(`<span class="badge badge-reports">⚠️ %d</span>`, rev.ReportCount)
	}

	actionControls := ""
	if rev.Status == "approved" {
		actionControls = fmt.Sprintf(`
			<form hx-post="/admin/reviews/%d/status?status=hidden" hx-target="#admin-review-row-%d" hx-swap="outerHTML" style="display: flex; gap: 0.35rem; align-items: center;">
				<input type="text" name="reason" placeholder="Moderation reason..." style="padding: 0.25rem 0.45rem; font-size: 0.78rem; border: 1px solid #cbd5e1; border-radius: 4px; width: 120px;" required>
				<button type="submit" class="btn-action" style="background: #fee2e2; color: #991b1b; border: 1px solid #fca5a5;">Hide</button>
			</form>
		`, rev.ID, rev.ID)
	} else {
		reasonNote := ""
		if rev.ModerationReason != "" {
			reasonNote = fmt.Sprintf(`<div style="font-size: 0.75rem; color: #b91c1c; margin-bottom: 0.25rem;">Reason: %s</div>`, escapeHTML(rev.ModerationReason))
		}
		actionControls = fmt.Sprintf(`
			<div>
				%s
				<form hx-post="/admin/reviews/%d/status?status=approved" hx-target="#admin-review-row-%d" hx-swap="outerHTML">
					<button type="submit" class="btn-action" style="background: #dcfce7; color: #166534; border: 1px solid #86efac;">Restore</button>
				</form>
			</div>
		`, reasonNote, rev.ID, rev.ID)
	}

	return fmt.Sprintf(`
		<tr id="admin-review-row-%d">
			<td style="font-weight: 700; color: #64748b;">#%d</td>
			<td>
				<a href="/restaurant/%d" target="_blank" style="color: #0f172a; font-weight: 700; text-decoration: underline;">
					%s
				</a>
			</td>
			<td>
				<div style="font-weight: 700;">%s</div>
				<div style="font-size: 0.75rem; color: #64748b;">%s</div>
			</td>
			<td>
				<span style="color: #f59e0b; font-weight: 800;">%s (%d/5)</span>
			</td>
			<td>
				<div style="max-width: 320px; line-height: 1.45; word-break: break-word;">%s</div>
				<div style="font-size: 0.72rem; color: #94a3b8; margin-top: 0.25rem;">%s</div>
			</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
	`, rev.ID, rev.ID, rev.RestaurantID, escapeHTML(rev.RestaurantName), escapeHTML(rev.UserName), escapeHTML(rev.UserEmail), renderStars(rev.Rating), rev.Rating, escapeHTML(rev.Comment), rev.CreatedAt.Format("02 Jan 2006 15:04"), statusBadge, reportBadge, actionControls)
}

// FormatRatingBadge formats the rating badge HTML for restaurant cards
func FormatRatingBadge(avg float64, count int) string {
	if count > 0 {
		return fmt.Sprintf(`<span class="rating-badge">★ %.1f (%d)</span>`, avg, count)
	}
	return `<span class="rating-badge">★ New</span>`
}
