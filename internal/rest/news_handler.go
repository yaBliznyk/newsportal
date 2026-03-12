package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"github.com/yaBliznyk/newsportal/internal/portal"
)

type NewsHandler struct {
	log         *slog.Logger
	newsManager *portal.NewsManager
}

func NewNewsHandler(log *slog.Logger, svc *portal.NewsManager) *NewsHandler {
	return &NewsHandler{
		log:         log,
		newsManager: svc,
	}
}

// Handle возвращает http.Handler с зарегистрированными роутами
func (h *NewsHandler) Handle() http.Handler {
	mux := echo.New()
	mux.GET("/v1/listNews", h.listNews)
	mux.GET("/v1/countNews", h.countNews)
	mux.GET("/v1/getNews/:id", h.getNews)
	mux.GET("/v1/getCategories", h.getCategories)
	mux.GET("/v1/getTags", h.getTags)
	return mux
}

// listNews обрабатывает GET /v1/listNews
func (h *NewsHandler) listNews(c *echo.Context) error {
	var req ListNewsReq
	err := c.Bind(&req)
	if err != nil {
		h.log.Error("bind request failed", "error", err)
		return h.sendJsonError(c, portal.ErrInvalidData)
	}

	news, err := h.newsManager.ListNews(c.Request().Context(), req.ToPortalFilter(), req.ToPortalPagination())
	if err != nil {
		h.log.Error("ListNews failed", "error", err)
		return h.sendJsonError(c, err)
	}

	return c.JSON(http.StatusOK, ListNewsResp{News: NewNewsList(news)})
}

// countNews обрабатывает GET /v1/countNews
func (h *NewsHandler) countNews(c *echo.Context) error {
	var req ListNewsReq
	err := c.Bind(&req)
	if err != nil {
		h.log.Error("bind request failed", "error", err)
		return h.sendJsonError(c, portal.ErrInvalidData)
	}

	count, err := h.newsManager.CountNews(c.Request().Context(), req.ToPortalFilter())
	if err != nil {
		h.log.Error("CountNews failed", "error", err)
		return h.sendJsonError(c, err)
	}

	return c.JSON(http.StatusOK, CountNewsResp{Count: count})
}

// getNews обрабатывает GET /v1/getNews
func (h *NewsHandler) getNews(c *echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return h.sendJsonError(c, portal.ErrInvalidData)
	}

	news, err := h.newsManager.GetNews(c.Request().Context(), id)
	if err != nil {
		h.log.Error("GetNews failed", "error", err, "id", id)
		return h.sendJsonError(c, err)
	}

	if news == nil {
		return h.sendJsonError(c, portal.ErrNewsNotFound)
	}

	return c.JSON(http.StatusOK, GetNewsResp{News: NewNews(news)})
}

// getCategories обрабатывает GET /v1/getCategories
func (h *NewsHandler) getCategories(c *echo.Context) error {
	categories, err := h.newsManager.ListCategories(c.Request().Context())
	if err != nil {
		h.log.Error("ListCategories failed", "error", err)
		return h.sendJsonError(c, err)
	}

	return c.JSON(http.StatusOK, GetCategoriesResp{Categories: NewCategories(categories)})
}

// getTags обрабатывает GET /v1/getTags
func (h *NewsHandler) getTags(c *echo.Context) error {
	tags, err := h.newsManager.ListTags(c.Request().Context())
	if err != nil {
		h.log.Error("ListTags failed", "error", err)
		return h.sendJsonError(c, err)
	}

	return c.JSON(http.StatusOK, GetTagsResp{Tags: NewTags(tags)})
}

// sendJsonError отправляет ошибку в json формате
func (h *NewsHandler) sendJsonError(c *echo.Context, err error) error {
	code := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, portal.ErrNewsNotFound):
		code = http.StatusNotFound
		msg = "not found"
	case errors.Is(err, portal.ErrInvalidData),
		errors.Is(err, portal.ErrInvalidCategoryID),
		errors.Is(err, portal.ErrInvalidPage),
		errors.Is(err, portal.ErrInvalidLimit),
		errors.Is(err, portal.ErrInvalidDateRange):
		code = http.StatusBadRequest
		msg = "invalid data"
	}

	return c.JSON(code, map[string]string{"error": msg})
}
