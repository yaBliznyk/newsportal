package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yaBliznyk/newsportal/internal/domain"
	"github.com/yaBliznyk/newsportal/internal/service"
	"github.com/yaBliznyk/newsportal/internal/service/mocks"
	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

// fixedTime используется для предсказуемых результатов в тестах
var fixedTime = time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)

// setupService создает сервис с моком репозитория для тестов
func setupService(t *testing.T) (*mocks.NewsRepository, *service.NewsService) {
	mockRepo := mocks.NewNewsRepository(t)
	svc := service.New(mockRepo)
	return mockRepo, svc
}

// ========================================
// ListNews
// ========================================

func TestListNews_Success(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedNews := []domain.ListNewsItem{
		{
			ID:          1,
			Title:       "Test News",
			Category:    domain.Category{ID: 1, Name: "Tech"},
			Tags:        []domain.Tag{{ID: 1, Name: "Go"}},
			Author:      "Test Author",
			CreatedAt:   fixedTime,
			PublishedAt: fixedTime,
		},
	}

	req := domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}

	mockRepo.EXPECT().ListNews(mock.Anything, req).Return(expectedNews, nil)

	resp, err := svc.ListNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
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
	mockRepo, svc := setupService(t)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	req := domain.ListNewsReq{
		CategoryID: 5,
		TagID:      10,
		From:       from,
		To:         to,
		Page:       2,
		Limit:      50,
	}

	expectedNews := []domain.ListNewsItem{}

	mockRepo.EXPECT().ListNews(mock.Anything, req).Return(expectedNews, nil)

	resp, err := svc.ListNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.News)
}

func TestListNews_EmptyResult(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}

	mockRepo.EXPECT().ListNews(mock.Anything, req).Return([]domain.ListNewsItem{}, nil)

	resp, err := svc.ListNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.News)
}

func TestListNews_RepositoryError(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.ListNewsReq{
		Page:  1,
		Limit: 20,
	}

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().ListNews(mock.Anything, req).Return(nil, expectedErr)

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "repo.ListNews")
	assert.Contains(t, err.Error(), "database error")
}

func TestListNews_InvalidCategoryID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.ListNewsReq{CategoryID: -1}

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "category_id")
}

func TestListNews_InvalidTagID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.ListNewsReq{TagID: -1}

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "tag_id")
}

func TestListNews_InvalidDateRange(t *testing.T) {
	_, svc := setupService(t)

	req := domain.ListNewsReq{
		From: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "from")
}

func TestListNews_InvalidPage(t *testing.T) {
	_, svc := setupService(t)

	req := domain.ListNewsReq{Page: -1}

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "page")
}

func TestListNews_InvalidLimit(t *testing.T) {
	_, svc := setupService(t)

	req := domain.ListNewsReq{Limit: -1}

	resp, err := svc.ListNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "limit")
}

// ========================================
// CountNews
// ========================================

func TestCountNews_Success(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.CountNewsReq{}
	mockRepo.EXPECT().CountNews(mock.Anything, req).Return(42, nil)

	resp, err := svc.CountNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 42, resp.Count)
}

func TestCountNews_WithFilters(t *testing.T) {
	mockRepo, svc := setupService(t)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	req := domain.CountNewsReq{
		CategoryID: 5,
		TagID:      10,
		From:       from,
		To:         to,
	}

	mockRepo.EXPECT().CountNews(mock.Anything, req).Return(10, nil)

	resp, err := svc.CountNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 10, resp.Count)
}

func TestCountNews_ZeroResult(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.CountNewsReq{}
	mockRepo.EXPECT().CountNews(mock.Anything, req).Return(0, nil)

	resp, err := svc.CountNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 0, resp.Count)
}

func TestCountNews_RepositoryError(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.CountNewsReq{}
	expectedErr := errors.New("database error")
	mockRepo.EXPECT().CountNews(mock.Anything, req).Return(0, expectedErr)

	resp, err := svc.CountNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "repo.CountNews")
	assert.Contains(t, err.Error(), "database error")
}

func TestCountNews_InvalidCategoryID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.CountNewsReq{CategoryID: -1}

	resp, err := svc.CountNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "category_id")
}

func TestCountNews_InvalidTagID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.CountNewsReq{TagID: -1}

	resp, err := svc.CountNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "tag_id")
}

func TestCountNews_InvalidDateRange(t *testing.T) {
	_, svc := setupService(t)

	req := domain.CountNewsReq{
		From: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := svc.CountNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "from")
}

// ========================================
// GetNews
// ========================================

func TestGetNews_Success(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedNews := &domain.News{
		ID:          123,
		Title:       "Full News Title",
		Preamble:    "Short description",
		Content:     "Full content of the news article",
		Category:    domain.Category{ID: 1, Name: "Technology"},
		Tags:        []domain.Tag{{ID: 1, Name: "Go"}, {ID: 2, Name: "Backend"}},
		Author:      "John Doe",
		CreatedAt:   fixedTime,
		PublishedAt: fixedTime,
	}

	req := domain.GetNewsReq{ID: 123}

	mockRepo.EXPECT().GetNews(mock.Anything, 123).Return(expectedNews, nil)

	resp, err := svc.GetNews(t.Context(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
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

func TestGetNews_RepositoryError(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := domain.GetNewsReq{ID: 1}
	expectedErr := errors.New("database error")
	mockRepo.EXPECT().GetNews(mock.Anything, 1).Return(nil, expectedErr)

	resp, err := svc.GetNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "repo.GetNews(1)")
	assert.Contains(t, err.Error(), "database error")
}

func TestGetNews_InvalidID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.GetNewsReq{ID: 0}

	resp, err := svc.GetNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "id")
}

func TestGetNews_NegativeID(t *testing.T) {
	_, svc := setupService(t)

	req := domain.GetNewsReq{ID: -1}

	resp, err := svc.GetNews(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.True(t, errors.Is(err, svcerrs.ErrInvalidData))
	assert.Contains(t, err.Error(), "id")
}

// ========================================
// GetCategories
// ========================================

func TestGetCategories_Success(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedCategories := []domain.Category{
		{ID: 1, Name: "Technology"},
		{ID: 2, Name: "Science"},
		{ID: 3, Name: "Sports"},
	}

	mockRepo.EXPECT().GetCategories(mock.Anything).Return(expectedCategories, nil)

	resp, err := svc.GetCategories(t.Context())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Categories, 3)
	assert.Equal(t, 1, resp.Categories[0].ID)
	assert.Equal(t, "Technology", resp.Categories[0].Name)
	assert.Equal(t, 2, resp.Categories[1].ID)
	assert.Equal(t, "Science", resp.Categories[1].Name)
	assert.Equal(t, 3, resp.Categories[2].ID)
	assert.Equal(t, "Sports", resp.Categories[2].Name)
}

func TestGetCategories_Empty(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.EXPECT().GetCategories(mock.Anything).Return([]domain.Category{}, nil)

	resp, err := svc.GetCategories(t.Context())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.Categories)
}

func TestGetCategories_RepositoryError(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().GetCategories(mock.Anything).Return(nil, expectedErr)

	resp, err := svc.GetCategories(t.Context())

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "repo.GetCategories")
	assert.Contains(t, err.Error(), "database error")
}

// ========================================
// GetTags
// ========================================

func TestGetTags_Success(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedTags := []domain.Tag{
		{ID: 1, Name: "Go"},
		{ID: 2, Name: "Python"},
		{ID: 3, Name: "JavaScript"},
	}

	mockRepo.EXPECT().GetTags(mock.Anything).Return(expectedTags, nil)

	resp, err := svc.GetTags(t.Context())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Tags, 3)
	assert.Equal(t, 1, resp.Tags[0].ID)
	assert.Equal(t, "Go", resp.Tags[0].Name)
	assert.Equal(t, 2, resp.Tags[1].ID)
	assert.Equal(t, "Python", resp.Tags[1].Name)
	assert.Equal(t, 3, resp.Tags[2].ID)
	assert.Equal(t, "JavaScript", resp.Tags[2].Name)
}

func TestGetTags_Empty(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.EXPECT().GetTags(mock.Anything).Return([]domain.Tag{}, nil)

	resp, err := svc.GetTags(t.Context())

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.Tags)
}

func TestGetTags_RepositoryError(t *testing.T) {
	mockRepo, svc := setupService(t)

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().GetTags(mock.Anything).Return(nil, expectedErr)

	resp, err := svc.GetTags(t.Context())

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "repo.GetTags")
	assert.Contains(t, err.Error(), "database error")
}
