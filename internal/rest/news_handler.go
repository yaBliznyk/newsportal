package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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
func (c *NewsHandler) Handle() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/listNews", c.listNews)
	mux.HandleFunc("GET /v1/countNews", c.countNews)
	mux.HandleFunc("GET /v1/getNews", c.getNews)
	mux.HandleFunc("GET /v1/getCategories", c.getCategories)
	mux.HandleFunc("GET /v1/getTags", c.getTags)
	return mux
}

// listNews обрабатывает GET /v1/listNews
func (c *NewsHandler) listNews(w http.ResponseWriter, r *http.Request) {
	var filter portal.ListNewsFilter
	pagination := portal.Pagination{
		Page:  1,
		Limit: 20,
	}

	if v := r.URL.Query().Get("category"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.CategoryID = id
		}
	}
	if v := r.URL.Query().Get("tag"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.TagID = id
		}
	}
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = t
		}
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			pagination.Page = p
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			pagination.Limit = l
		}
	}

	news, err := c.newsManager.ListNews(r.Context(), filter, pagination)
	if err != nil {
		c.log.Error("ListNews failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, ListNewsResp{News: NewNewsList(news)})
}

// countNews обрабатывает GET /v1/countNews
func (c *NewsHandler) countNews(w http.ResponseWriter, r *http.Request) {
	var filter portal.ListNewsFilter

	if v := r.URL.Query().Get("category"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.CategoryID = id
		}
	}
	if v := r.URL.Query().Get("tag"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.TagID = id
		}
	}
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = t
		}
	}

	count, err := c.newsManager.CountNews(r.Context(), filter)
	if err != nil {
		c.log.Error("CountNews failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, CountNewsResp{Count: count})
}

// getNews обрабатывает GET /v1/getNews
func (c *NewsHandler) getNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		c.writeError(w, portal.ErrInvalidData)
		return
	}

	news, err := c.newsManager.GetNews(r.Context(), id)
	if err != nil {
		c.log.Error("GetNews failed", "error", err, "id", id)
		c.writeError(w, err)
		return
	}

	if news == nil {
		c.writeError(w, portal.ErrNewsNotFound)
		return
	}

	c.writeJSON(w, GetNewsResp{News: NewNews(news)})
}

// getCategories обрабатывает GET /v1/getCategories
func (c *NewsHandler) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.newsManager.ListCategories(r.Context())
	if err != nil {
		c.log.Error("ListCategories failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, GetCategoriesResp{Categories: NewCategories(categories)})
}

// getTags обрабатывает GET /v1/getTags
func (c *NewsHandler) getTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.newsManager.ListTags(r.Context())
	if err != nil {
		c.log.Error("ListTags failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, GetTagsResp{Tags: NewTags(tags)})
}

// writeJSON отправляет JSON-ответ
func (c *NewsHandler) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeError отправляет ошибку
func (c *NewsHandler) writeError(w http.ResponseWriter, err error) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
