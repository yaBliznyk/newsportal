package public

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

type Controller struct {
	log *slog.Logger
	svc domain.Service
}

func NewController(log *slog.Logger, svc domain.Service) *Controller {
	return &Controller{
		log: log,
		svc: svc,
	}
}

// Init регистрирует роуты на стандартном http.ServeMux
func (c *Controller) Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/listNews", c.listNews)
	mux.HandleFunc("GET /v1/countNews", c.countNews)
	mux.HandleFunc("GET /v1/getNews", c.getNews)
	mux.HandleFunc("GET /v1/getCategories", c.getCategories)
	mux.HandleFunc("GET /v1/getTags", c.getTags)
}

// writeJSON отправляет JSON-ответ
func (c *Controller) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeError отправляет ошибку
func (c *Controller) writeError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, svcerrs.ErrDataNotFound):
		code = http.StatusNotFound
		msg = "data not found"
	case errors.Is(err, svcerrs.ErrInvalidData):
		code = http.StatusBadRequest
		msg = "invalid data"
	case errors.Is(err, svcerrs.ErrAlreadyExist):
		code = http.StatusConflict
		msg = "already exists"
	case errors.Is(err, svcerrs.ErrAccessDenied):
		code = http.StatusForbidden
		msg = "access denied"
	case errors.Is(err, svcerrs.ErrUnauthorized):
		code = http.StatusUnauthorized
		msg = "unauthorized"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
