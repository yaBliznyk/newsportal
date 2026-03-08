package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/yaBliznyk/newsportal/internal/domain"
)

// countNews обрабатывает GET /v1/countNews
func (c *NewsHandler) countNews(w http.ResponseWriter, r *http.Request) {
	req := domain.CountNewsReq{}

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

	resp, err := c.newsManager.CountNews(r.Context(), req)
	if err != nil {
		c.log.Error("CountNews failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
