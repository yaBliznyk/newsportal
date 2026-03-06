package public

import (
	"net/http"
)

// getCategories обрабатывает GET /v1/getCategories
func (c *Controller) getCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := c.svc.GetCategories(r.Context())
	if err != nil {
		c.log.Error("GetCategories failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
