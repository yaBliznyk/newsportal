package public

import (
	"net/http"
)

// getTags обрабатывает GET /v1/getTags
func (c *Controller) getTags(w http.ResponseWriter, r *http.Request) {
	resp, err := c.svc.GetTags(r.Context())
	if err != nil {
		c.log.Error("GetTags failed", "error", err)
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, resp)
}
