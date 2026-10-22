package fooddiscovery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// AdminCategoriesList displays the categories and cuisines management table
func (h *Handler) AdminCategoriesList(ctx echo.Context) error {
	u, _ := h.c.Auth.GetAuthenticatedUser(ctx)
	msg := ctx.QueryParam("msg")
	errMsg := ctx.QueryParam("err")

	categories, err := h.c.FoodDiscovery.GetAllCategoriesAdmin(ctx.Request().Context())
	if err != nil {
		return ctx.HTML(http.StatusOK, fmt.Sprintf(`<div class="alert alert-error">Failed to load categories: %s</div>`, err.Error()))
	}

	html := RenderAdminCategoriesPage(categories, u, msg, errMsg)
	return ctx.HTML(http.StatusOK, html)
}

// AdminCreateCategory handles adding a new category
func (h *Handler) AdminCreateCategory(ctx echo.Context) error {
	name := ctx.FormValue("name")
	catType := ctx.FormValue("category_type")
	description := ctx.FormValue("description")
	icon := ctx.FormValue("icon")

	if strings.TrimSpace(name) == "" {
		return ctx.Redirect(http.StatusSeeOther, "/admin/categories?err=Category+name+is+required")
	}

	_, err := h.c.FoodDiscovery.CreateCategory(ctx.Request().Context(), name, catType, description, icon)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/categories?err=%s", url.QueryEscape(err.Error())))
	}

	return ctx.Redirect(http.StatusSeeOther, "/admin/categories?msg=Category+created+successfully")
}

// AdminToggleCategory toggles active status for a category
func (h *Handler) AdminToggleCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.String(http.StatusBadRequest, "Invalid category ID")
	}

	_, err = h.c.FoodDiscovery.ToggleCategoryActive(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	// Fetch updated category
	cat, err := h.c.FoodDiscovery.GetCategoryByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	// If HTMX request, render updated row
	if ctx.Request().Header.Get("HX-Request") == "true" {
		return ctx.HTML(http.StatusOK, RenderAdminCategoryRow(*cat))
	}

	return ctx.Redirect(http.StatusSeeOther, "/admin/categories?msg=Category+status+updated")
}

// AdminUpdateCategory updates an existing category
func (h *Handler) AdminUpdateCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return ctx.Redirect(http.StatusSeeOther, "/admin/categories?err=Invalid+category+ID")
	}

	name := ctx.FormValue("name")
	catType := ctx.FormValue("category_type")
	description := ctx.FormValue("description")
	icon := ctx.FormValue("icon")

	err = h.c.FoodDiscovery.UpdateCategory(ctx.Request().Context(), id, name, catType, description, icon)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/categories?err=%s", url.QueryEscape(err.Error())))
	}

	return ctx.Redirect(http.StatusSeeOther, "/admin/categories?msg=Category+updated+successfully")
}

// CategoryDiscoveryPage handles direct links to a category results page (/category/:slug)
func (h *Handler) CategoryDiscoveryPage(ctx echo.Context) error {
	slug := strings.TrimSpace(ctx.Param("slug"))
	if slug == "" {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	cat, err := h.c.FoodDiscovery.GetCategoryBySlug(ctx.Request().Context(), slug)
	if err != nil || cat == nil {
		// Try to see if it's a known slug or fallback
		return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/?category=%s", url.QueryEscape(slug)))
	}

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/?category=%s", url.QueryEscape(cat.Slug)))
}
