package rest

import (
	"net/http"
)

// getCategories обрабатывает GET /v1/getCategories
func (c *NewsHandler) getCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := c.newsManager.GetCategories(r.Context())
	if err != nil {
		c.log.Error("GetCategories failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
