package model

import (
	"context"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/irvankadhafi/articles-go/search-service/event"
	"time"
)

type (
	// ArticleSearchUsecase :nodoc:
	ArticleSearchUsecase interface {
		IndexByID(ctx context.Context, id int) error
		HandleEvent(ctx context.Context, msg *event.PubSubPayload) error
		Search(ctx context.Context, req SearchArticleRequest) (spaces []*Article, count int, err error)
	}
	// ArticleRepository :nodoc:
	ArticleRepository interface {
		Search(ctx context.Context, req SearchArticleRequest) (ids []int, count int, err error)
		IndexES(ctx context.Context, article *Article) error
	}

	// ArticleSortType :nodoc:
	ArticleSortType int
)

// ArticleSortType enums
const (
	SortArticleByCreatedAtDesc = ArticleSortType(0)
	SortEArticleByCreatedAtAsc = ArticleSortType(1)
)

// ArticleFilter :nodoc:
type ArticleFilter struct {
	Author string
}

// SearchArticleRequest :nodoc:
type SearchArticleRequest struct {
	Query  string
	Page   int
	Size   int
	Filter ArticleFilter
	Sorter ArticleSortType
}

var (
	//ESArticleIndex :nodoc:
	ESArticleIndex = "articles-v1"

	//ESArticleIndexAlias :nodoc:
	ESArticleIndexAlias = "articles"

	//ESArticleMapping :nodoc:
	ESArticleMapping = `
		{
			"mappings": {
				"properties": {
					"id": {
						"type": "integer"
					},
					"author": {
						"type": "text"
					},
					"title": {
						"type": "text",
						"analyzer": "autocomplete"
					},
					"body": {
						"type": "text",
						"analyzer": "autocomplete"
					},
					"created_at": {
						"type": "date"
					}
				}
			},
			"settings": {
				"analysis": {
					"analyzer": {
						"autocomplete": {
							"type": "custom",
							"tokenizer": "standard",
							"filter": [
								"lowercase",
								"autocomplete_filter"
							]
						}
					},
					"filter": {
						"autocomplete_filter": {
							"type": "edge_ngram",
							"min_gram": 1,
							"max_gram": 20
						}
					}
				}
			}
		}
`
)

type Article struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// ToFieldAndDirection converts sortType into respected sort query based on olivere's library. True means sort ascendingly, False means sort descendingly.
func (st ArticleSortType) ToFieldAndDirection() (string, bool) {
	switch st {
	case SortEArticleByCreatedAtAsc:
		return "created_at", true
	default:
		return "created_at", false
	}
}

// NewArticleFromProto :nodoc:
func NewArticleFromProto(p *pb.Article) *Article {
	createdAt, _ := time.Parse(time.RFC3339Nano, p.GetCreatedAt())
	return &Article{
		ID:        int(p.GetId()),
		Author:    p.GetAuthor(),
		Title:     p.GetTitle(),
		Body:      p.GetBody(),
		CreatedAt: createdAt,
	}
}
