package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

// listNews обрабатывает GET /v1/listNews
func (c *NewsHandler) listNews(w http.ResponseWriter, r *http.Request) {
	req := domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}

	if v := r.URL.Query().Get("category"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			req.CategoryID = id
		}
	}
	if v := r.URL.Query().Get("tag"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			req.TagID = id
		}
	}
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			req.From = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			req.To = t
		}
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			req.Page = p
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			req.Limit = l
		}
	}

	resp, err := c.newsManager.ListNews(r.Context(), req)
	if err != nil {
		c.log.Error("ListNews failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
