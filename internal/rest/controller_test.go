package rest_test

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/domain/mocks"
	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// fixedTime используется для предсказуемых результатов в тестах
var fixedTime = time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)

// setupController создает контроллер с моком сервиса для тестов
func setupController(t *testing.T) (*mocks.Service, *http.ServeMux) {
	mockSvc := mocks.NewService(t)
	log := slog.New(slog.NewTextHandler(io.Discard, nil)) // discard logs
	ctrl := rest.NewNewsHandler(log, mockSvc)

	mux := http.NewServeMux()
	ctrl.Init(mux)

	return mockSvc, mux
}

// doRequest выполняет HTTP запрос и возвращает ответ
func doRequest(mux *http.ServeMux, method, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

// parseErrorResp парсит JSON ответ с ошибкой
func parseErrorResp(t *testing.T, body []byte) map[string]string {
	var resp map[string]string
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

// ========================================
// Тесты ошибок (используются во всех обработчиках)
// ========================================

// errorTestCases возвращает тест-кейсы для всех типов ошибок
func errorTestCases() []struct {
	name       string
	err        error
	wantStatus int
	wantMsg    string
} {
	return []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"ErrDataNotFound", svcerrs.ErrDataNotFound, http.StatusNotFound, "data not found"},
		{"ErrInvalidData", svcerrs.ErrInvalidData, http.StatusBadRequest, "invalid data"},
		{"ErrAlreadyExist", svcerrs.ErrAlreadyExist, http.StatusConflict, "already exists"},
		{"ErrAccessDenied", svcerrs.ErrAccessDenied, http.StatusForbidden, "access denied"},
		{"ErrUnauthorized", svcerrs.ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{"UnknownError", errors.New("some unexpected error"), http.StatusInternalServerError, "internal server error"},
	}
}

// ========================================
// ListNews
// ========================================

func TestListNews_Success(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.ListNewsResp{
		News: []domain.ListNewsItem{
			{
				ID:          1,
				Title:       "Test News",
				Category:    domain.Category{ID: 1, Name: "Tech"},
				Tags:        []domain.Tag{{ID: 1, Name: "Go"}},
				Author:      "Test Author",
				CreatedAt:   fixedTime,
				PublishedAt: fixedTime,
			},
		},
	}

	mockSvc.EXPECT().ListNews(mock.Anything, domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/listNews")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp domain.ListNewsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.News, 1)
	assert.Equal(t, 1, resp.News[0].ID)
	assert.Equal(t, "Test News", resp.News[0].Title)
	assert.Equal(t, "Test Author", resp.News[0].Author)
	assert.Equal(t, 1, resp.News[0].Category.ID)
	assert.Equal(t, "Tech", resp.News[0].Category.Name)
	require.Len(t, resp.News[0].Tags, 1)
	assert.Equal(t, 1, resp.News[0].Tags[0].ID)
	assert.Equal(t, "Go", resp.News[0].Tags[0].Name)
}

func TestListNews_WithFilters(t *testing.T) {
	mockSvc, mux := setupController(t)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	expectedResp := &domain.ListNewsResp{News: []domain.ListNewsItem{}}

	mockSvc.EXPECT().ListNews(mock.Anything, domain.ListNewsReq{
		CategoryID: 5,
		TagID:      10,
		From:       from,
		To:         to,
		Page:       2,
		Limit:      50,
	}).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet,
		"/v1/listNews?category=5&tag=10&from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z&page=2&limit=50")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListNews_DefaultPagination(t *testing.T) {
	mockSvc, mux := setupController(t)

	mockSvc.EXPECT().ListNews(mock.Anything, domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}).Return(&domain.ListNewsResp{}, nil)

	w := doRequest(mux, http.MethodGet, "/v1/listNews")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListNews_InvalidParamsIgnored(t *testing.T) {
	mockSvc, mux := setupController(t)

	// Невалидные параметры должны игнорироваться, используются дефолты
	mockSvc.EXPECT().ListNews(mock.Anything, domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}).Return(&domain.ListNewsResp{}, nil)

	w := doRequest(mux, http.MethodGet,
		"/v1/listNews?category=invalid&tag=abc&page=-1&limit=0")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListNews_Errors(t *testing.T) {
	for _, tc := range errorTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc, mux := setupController(t)

			mockSvc.EXPECT().ListNews(mock.Anything, domain.ListNewsReq{
				Page:  1,
				Limit: 20,
			}).Return(nil, tc.err)

			w := doRequest(mux, http.MethodGet, "/v1/listNews")

			assert.Equal(t, tc.wantStatus, w.Code)
			resp := parseErrorResp(t, w.Body.Bytes())
			assert.Equal(t, tc.wantMsg, resp["error"])
		})
	}
}

// ========================================
// CountNews
// ========================================

func TestCountNews_Success(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.CountNewsResp{Count: 42}

	mockSvc.EXPECT().CountNews(mock.Anything, domain.CountNewsReq{}).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/countNews")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp domain.CountNewsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 42, resp.Count)
}

func TestCountNews_WithFilters(t *testing.T) {
	mockSvc, mux := setupController(t)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	expectedResp := &domain.CountNewsResp{Count: 10}

	mockSvc.EXPECT().CountNews(mock.Anything, domain.CountNewsReq{
		CategoryID: 5,
		TagID:      10,
		From:       from,
		To:         to,
	}).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet,
		"/v1/countNews?category=5&tag=10&from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.CountNewsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 10, resp.Count)
}

func TestCountNews_Errors(t *testing.T) {
	for _, tc := range errorTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc, mux := setupController(t)

			mockSvc.EXPECT().CountNews(mock.Anything, domain.CountNewsReq{}).Return(nil, tc.err)

			w := doRequest(mux, http.MethodGet, "/v1/countNews")

			assert.Equal(t, tc.wantStatus, w.Code)
			resp := parseErrorResp(t, w.Body.Bytes())
			assert.Equal(t, tc.wantMsg, resp["error"])
		})
	}
}

// ========================================
// GetNews
// ========================================

func TestGetNews_Success(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.GetNewsResp{
		News: domain.News{
			ID:          123,
			Title:       "Full News Title",
			Preamble:    "Short description",
			Content:     "Full content of the news article",
			Category:    domain.Category{ID: 1, Name: "Technology"},
			Tags:        []domain.Tag{{ID: 1, Name: "Go"}, {ID: 2, Name: "Backend"}},
			Author:      "John Doe",
			CreatedAt:   fixedTime,
			PublishedAt: fixedTime,
		},
	}

	mockSvc.EXPECT().GetNews(mock.Anything, domain.GetNewsReq{ID: 123}).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/getNews?id=123")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp domain.GetNewsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 123, resp.News.ID)
	assert.Equal(t, "Full News Title", resp.News.Title)
	assert.Equal(t, "Short description", resp.News.Preamble)
	assert.Equal(t, "Full content of the news article", resp.News.Content)
	assert.Equal(t, "John Doe", resp.News.Author)
	assert.Equal(t, 1, resp.News.Category.ID)
	assert.Equal(t, "Technology", resp.News.Category.Name)
	require.Len(t, resp.News.Tags, 2)
	assert.Equal(t, 1, resp.News.Tags[0].ID)
	assert.Equal(t, "Go", resp.News.Tags[0].Name)
	assert.Equal(t, 2, resp.News.Tags[1].ID)
	assert.Equal(t, "Backend", resp.News.Tags[1].Name)
}

func TestGetNews_InvalidID(t *testing.T) {
	_, mux := setupController(t)

	w := doRequest(mux, http.MethodGet, "/v1/getNews?id=invalid")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseErrorResp(t, w.Body.Bytes())
	assert.Equal(t, "invalid data", resp["error"])
}

func TestGetNews_MissingID(t *testing.T) {
	_, mux := setupController(t)

	w := doRequest(mux, http.MethodGet, "/v1/getNews")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseErrorResp(t, w.Body.Bytes())
	assert.Equal(t, "invalid data", resp["error"])
}

func TestGetNews_Errors(t *testing.T) {
	for _, tc := range errorTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc, mux := setupController(t)

			mockSvc.EXPECT().GetNews(mock.Anything, domain.GetNewsReq{ID: 1}).Return(nil, tc.err)

			w := doRequest(mux, http.MethodGet, "/v1/getNews?id=1")

			assert.Equal(t, tc.wantStatus, w.Code)
			resp := parseErrorResp(t, w.Body.Bytes())
			assert.Equal(t, tc.wantMsg, resp["error"])
		})
	}
}

// ========================================
// GetCategories
// ========================================

func TestGetCategories_Success(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.GetCategoriesResp{
		Categories: []domain.Category{
			{ID: 1, Name: "Technology"},
			{ID: 2, Name: "Science"},
			{ID: 3, Name: "Sports"},
		},
	}

	mockSvc.EXPECT().GetCategories(mock.Anything).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/getCategories")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp domain.GetCategoriesResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Categories, 3)
	assert.Equal(t, 1, resp.Categories[0].ID)
	assert.Equal(t, "Technology", resp.Categories[0].Name)
	assert.Equal(t, 2, resp.Categories[1].ID)
	assert.Equal(t, "Science", resp.Categories[1].Name)
	assert.Equal(t, 3, resp.Categories[2].ID)
	assert.Equal(t, "Sports", resp.Categories[2].Name)
}

func TestGetCategories_Empty(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.GetCategoriesResp{Categories: []domain.Category{}}

	mockSvc.EXPECT().GetCategories(mock.Anything).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/getCategories")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.GetCategoriesResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp.Categories)
}

func TestGetCategories_Errors(t *testing.T) {
	for _, tc := range errorTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc, mux := setupController(t)

			mockSvc.EXPECT().GetCategories(mock.Anything).Return(nil, tc.err)

			w := doRequest(mux, http.MethodGet, "/v1/getCategories")

			assert.Equal(t, tc.wantStatus, w.Code)
			resp := parseErrorResp(t, w.Body.Bytes())
			assert.Equal(t, tc.wantMsg, resp["error"])
		})
	}
}

// ========================================
// GetTags
// ========================================

func TestGetTags_Success(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.GetTagsResp{
		Tags: []domain.Tag{
			{ID: 1, Name: "Go"},
			{ID: 2, Name: "Python"},
			{ID: 3, Name: "JavaScript"},
		},
	}

	mockSvc.EXPECT().GetTags(mock.Anything).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/getTags")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp domain.GetTagsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Tags, 3)
	assert.Equal(t, 1, resp.Tags[0].ID)
	assert.Equal(t, "Go", resp.Tags[0].Name)
	assert.Equal(t, 2, resp.Tags[1].ID)
	assert.Equal(t, "Python", resp.Tags[1].Name)
	assert.Equal(t, 3, resp.Tags[2].ID)
	assert.Equal(t, "JavaScript", resp.Tags[2].Name)
}

func TestGetTags_Empty(t *testing.T) {
	mockSvc, mux := setupController(t)

	expectedResp := &domain.GetTagsResp{Tags: []domain.Tag{}}

	mockSvc.EXPECT().GetTags(mock.Anything).Return(expectedResp, nil)

	w := doRequest(mux, http.MethodGet, "/v1/getTags")

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.GetTagsResp
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp.Tags)
}

func TestGetTags_Errors(t *testing.T) {
	for _, tc := range errorTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc, mux := setupController(t)

			mockSvc.EXPECT().GetTags(mock.Anything).Return(nil, tc.err)

			w := doRequest(mux, http.MethodGet, "/v1/getTags")

			assert.Equal(t, tc.wantStatus, w.Code)
			resp := parseErrorResp(t, w.Body.Bytes())
			assert.Equal(t, tc.wantMsg, resp["error"])
		})
	}
}
