package rest

import (
	"net/http"
)

// getTags обрабатывает GET /v1/getTags
func (c *NewsHandler) getTags(w http.ResponseWriter, r *http.Request) {
	resp, err := c.newsManager.ListTags(r.Context())
	if err != nil {
		c.log.Error("ListTags failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
