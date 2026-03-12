// nolint
//
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"
)

var Columns = struct {
	Category struct {
		ID, Name, SortOrder, StatusID string

		Status string
	}
	News struct {
		ID, Title, CategoryID, TagIDs, Author, Preamble, Content, CreatedAt, PublishedAt, StatusID string

		Category, Status string
	}
	Status struct {
		ID, Name string
	}
	Tag struct {
		ID, Name, StatusID string

		Status string
	}
}{
	Category: struct {
		ID, Name, SortOrder, StatusID string

		Status string
	}{
		ID:        "categoryId",
		Name:      "name",
		SortOrder: "sortOrder",
		StatusID:  "statusId",

		Status: "Status",
	},
	News: struct {
		ID, Title, CategoryID, TagIDs, Author, Preamble, Content, CreatedAt, PublishedAt, StatusID string

		Category, Status string
	}{
		ID:          "newsId",
		Title:       "title",
		CategoryID:  "categoryId",
		TagIDs:      "tagIds",
		Author:      "author",
		Preamble:    "preamble",
		Content:     "content",
		CreatedAt:   "createdAt",
		PublishedAt: "publishedAt",
		StatusID:    "statusId",

		Category: "Category",
		Status:   "Status",
	},
	Status: struct {
		ID, Name string
	}{
		ID:   "statusId",
		Name: "name",
	},
	Tag: struct {
		ID, Name, StatusID string

		Status string
	}{
		ID:       "tagId",
		Name:     "name",
		StatusID: "statusId",

		Status: "Status",
	},
}

var Tables = struct {
	Category struct {
		Name, Alias string
	}
	News struct {
		Name, Alias string
	}
	Status struct {
		Name, Alias string
	}
	Tag struct {
		Name, Alias string
	}
}{
	Category: struct {
		Name, Alias string
	}{
		Name:  "categories",
		Alias: "t",
	},
	News: struct {
		Name, Alias string
	}{
		Name:  "news",
		Alias: "t",
	},
	Status: struct {
		Name, Alias string
	}{
		Name:  "statuses",
		Alias: "t",
	},
	Tag: struct {
		Name, Alias string
	}{
		Name:  "tags",
		Alias: "t",
	},
}

type Category struct {
	tableName struct{} `pg:"categories,alias:t,discard_unknown_columns"`

	ID        int    `pg:"categoryId,pk"`
	Name      string `pg:"name,use_zero"`
	SortOrder int    `pg:"sortOrder,use_zero"`
	StatusID  int    `pg:"statusId,use_zero"`

	Status *Status `pg:"fk:statusId,rel:has-one"`
}

type News struct {
	tableName struct{} `pg:"news,alias:t,discard_unknown_columns"`

	ID          int       `pg:"newsId,pk"`
	Title       string    `pg:"title,use_zero"`
	CategoryID  int       `pg:"categoryId,use_zero"`
	TagIDs      []int     `pg:"tagIds,array"`
	Author      string    `pg:"author,use_zero"`
	Preamble    string    `pg:"preamble,use_zero"`
	Content     *string   `pg:"content"`
	CreatedAt   time.Time `pg:"createdAt,use_zero"`
	PublishedAt time.Time `pg:"publishedAt,use_zero"`
	StatusID    int       `pg:"statusId,use_zero"`

	Category *Category `pg:"fk:categoryId,rel:has-one"`
	Status   *Status   `pg:"fk:statusId,rel:has-one"`
}

type Status struct {
	tableName struct{} `pg:"statuses,alias:t,discard_unknown_columns"`

	ID   int     `pg:"statusId,pk"`
	Name *string `pg:"name"`
}

type Tag struct {
	tableName struct{} `pg:"tags,alias:t,discard_unknown_columns"`

	ID       int    `pg:"tagId,pk"`
	Name     string `pg:"name,use_zero"`
	StatusID int    `pg:"statusId,use_zero"`

	Status *Status `pg:"fk:statusId,rel:has-one"`
}
