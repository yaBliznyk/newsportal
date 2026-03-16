package newsportal_test

import (
	"context"
	"os"
	"testing"

	"github.com/yaBliznyk/newsportal/tests/newsportal"
)

func getRPCURL() string {
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}
	return "http://localhost" + httpAddr + "/rpc/"
}

func TestCategories(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	categories, err := client.News.Categories(ctx)
	if err != nil {
		t.Fatalf("Categories failed: %v", err)
	}

	if len(categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(categories))
	}
}

func TestTags(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	tags, err := client.News.Tags(ctx)
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if len(tags) != 6 {
		t.Errorf("expected 6 tags, got %d", len(tags))
	}
}

func TestCount(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{})
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != 5 {
		t.Errorf("expected 5, got %d", count)
	}
}

func TestCountWithCategoryFilter(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{Category_id: 5})
	if err != nil {
		t.Fatalf("Count with category filter failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestCountWithTagFilter(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{Tag_id: 5})
	if err != nil {
		t.Fatalf("Count with tag filter failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestCountWithCombinedFilters(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{Category_id: 1, Tag_id: 1})
	if err != nil {
		t.Fatalf("Count with combined filters failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestList(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(news) != 5 {
		t.Errorf("expected 5 news, got %d", len(news))
	}
}

func TestListWithCategoryFilter(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Category_id: 1}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with category filter failed: %v", err)
	}

	if len(news) != 3 {
		t.Errorf("expected 3 news, got %d", len(news))
	}

	if news[0].Title != "Прорыв в области искусственного интеллекта" {
		t.Errorf("unexpected title: %s", news[0].Title)
	}
}

func TestListWithTagFilter(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Tag_id: 5}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with tag filter failed: %v", err)
	}

	if len(news) != 1 {
		t.Errorf("expected 1 news, got %d", len(news))
	}

	if news[0].Title != "Открыта новая экзопланета" {
		t.Errorf("unexpected title: %s", news[0].Title)
	}
}

func TestListWithPagination(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{}, newsportal.Pager{Page: 1, Limit: 2})
	if err != nil {
		t.Fatalf("List with pagination failed: %v", err)
	}

	if len(news) != 2 {
		t.Errorf("expected 2 news, got %d", len(news))
	}
}

func TestListWithDateFrom(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.List(ctx, newsportal.NewsFilter{From: "2026-03-01T00:00:00Z"}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with From filter failed: %v", err)
	}
}

func TestListWithDateTo(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.List(ctx, newsportal.NewsFilter{To: "2026-02-28T23:59:59Z"}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with To filter failed: %v", err)
	}
}

func TestListWithDateRange(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.List(ctx, newsportal.NewsFilter{From: "2026-02-01T00:00:00Z", To: "2026-02-28T23:59:59Z"}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with From+To filters failed: %v", err)
	}
}

func TestListWithCombinedFilters(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Category_id: 5, Tag_id: 5}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with category+tag filters failed: %v", err)
	}

	if len(news) != 1 {
		t.Errorf("expected 1 news, got %d", len(news))
	}
}

func TestListWithInvalidCategory(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.List(ctx, newsportal.NewsFilter{Category_id: -1}, newsportal.Pager{})
	if err == nil {
		t.Error("expected error for invalid category")
	}
}

func TestListWithInvalidDateRange(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.List(ctx, newsportal.NewsFilter{From: "2026-12-31T00:00:00Z", To: "2026-01-01T00:00:00Z"}, newsportal.Pager{})
	if err == nil {
		t.Error("expected error for invalid date range")
	}
}

func TestCountWithInvalidCategory(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.Count(ctx, newsportal.NewsFilter{Category_id: -5})
	if err == nil {
		t.Error("expected error for invalid category")
	}
}

func TestCountWithInvalidDateRange(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	_, err := client.News.Count(ctx, newsportal.NewsFilter{From: "2026-06-01T00:00:00Z", To: "2026-01-01T00:00:00Z"})
	if err == nil {
		t.Error("expected error for invalid date range")
	}
}

func TestListNoResults(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Category_id: 999}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with non-existent category failed: %v", err)
	}

	if len(news) != 0 {
		t.Errorf("expected empty array, got %d news", len(news))
	}
}

func TestGet(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if news == nil {
		t.Fatal("expected news, got nil")
	}

	if news.ID != 1 {
		t.Errorf("expected ID 1, got %d", news.ID)
	}

	if news.Title != "Прорыв в области искусственного интеллекта" {
		t.Errorf("unexpected title: %s", news.Title)
	}

	if news.Content == nil || *news.Content == "" {
		t.Error("expected content")
	}

	if news.Preamble == "" {
		t.Error("expected preamble")
	}

	if news.Category == nil {
		t.Error("expected category")
	}

	if len(news.Tags) == 0 {
		t.Error("expected tags")
	}
}

func TestGetNotFound(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 999)
	if err != nil {
		t.Fatalf("Get with invalid id failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil, got %+v", news)
	}
}

func TestGetMissingID(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 0)
	if err != nil {
		t.Fatalf("Get without id failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil, got %+v", news)
	}
}

func TestGetDraftNews(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 4)
	if err != nil {
		t.Fatalf("Get draft news failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil for draft news, got %+v", news)
	}
}

func TestGetNewsInUnpublishedCategory(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 5)
	if err != nil {
		t.Fatalf("Get news in unpublished category failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil for news in unpublished category, got %+v", news)
	}
}

func TestGetDeletedNews(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 6)
	if err != nil {
		t.Fatalf("Get deleted news failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil for deleted news, got %+v", news)
	}
}

func TestGetNewsInDeletedCategory(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 7)
	if err != nil {
		t.Fatalf("Get news in deleted category failed: %v", err)
	}

	if news != nil {
		t.Errorf("expected nil for news in deleted category, got %+v", news)
	}
}

func TestGetNewsWithUnpublishedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 8)
	if err != nil {
		t.Fatalf("Get news with unpublished tag failed: %v", err)
	}

	if news == nil {
		t.Fatal("expected news, got nil")
	}

	if news.ID != 8 {
		t.Errorf("expected ID 8, got %d", news.ID)
	}

	if news.Title != "Новость с неопубликованным тегом" {
		t.Errorf("unexpected title: %s", news.Title)
	}

	if len(news.Tags) != 0 {
		t.Errorf("expected empty tags array, got %d tags", len(news.Tags))
	}
}

func TestGetNewsWithDeletedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.Get(ctx, 9)
	if err != nil {
		t.Fatalf("Get news with deleted tag failed: %v", err)
	}

	if news == nil {
		t.Fatal("expected news, got nil")
	}

	if news.ID != 9 {
		t.Errorf("expected ID 9, got %d", news.ID)
	}

	if news.Title != "Новость с удаленным тегом" {
		t.Errorf("unexpected title: %s", news.Title)
	}

	if len(news.Tags) != 0 {
		t.Errorf("expected empty tags array, got %d tags", len(news.Tags))
	}
}

func TestCountWithUnpublishedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{Tag_id: 7})
	if err != nil {
		t.Fatalf("Count with unpublished tag failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestListWithUnpublishedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Tag_id: 7}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with unpublished tag failed: %v", err)
	}

	if len(news) != 1 {
		t.Errorf("expected 1 news, got %d", len(news))
	}

	if len(news[0].Tags) != 0 {
		t.Errorf("expected empty tags array, got %d tags", len(news[0].Tags))
	}
}

func TestCountWithDeletedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{Tag_id: 8})
	if err != nil {
		t.Fatalf("Count with deleted tag failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestListWithDeletedTag(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	news, err := client.News.List(ctx, newsportal.NewsFilter{Tag_id: 8}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List with deleted tag failed: %v", err)
	}

	if len(news) != 1 {
		t.Errorf("expected 1 news, got %d", len(news))
	}

	if len(news[0].Tags) != 0 {
		t.Errorf("expected empty tags array, got %d tags", len(news[0].Tags))
	}
}

func TestCountVsListConsistency(t *testing.T) {
	client := newsportal.NewClient(getRPCURL(), nil)
	ctx := context.Background()

	count, err := client.News.Count(ctx, newsportal.NewsFilter{})
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	news, err := client.News.List(ctx, newsportal.NewsFilter{}, newsportal.Pager{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if count != len(news) {
		t.Errorf("count (%d) != list length (%d)", count, len(news))
	}
}
