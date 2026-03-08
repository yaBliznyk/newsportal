package rest

import (
	"net/http"
	"strconv"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// getNews обрабатывает GET /v1/getNews
func (c *NewsHandler) getNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		c.writeError(w, svcerrs.ErrInvalidData)
		return
	}

	resp, err := c.newsManager.GetNews(r.Context(), domain.GetNewsReq{ID: id})
	if err != nil {
		c.log.Error("GetNews failed", "error", err, "id", id)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
