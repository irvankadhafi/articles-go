package model

import (
	"context"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"time"
)

type Article struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:'now()':::STRING::TIMESTAMP" gorm:"->;<-:create"`
}

type ArticleUsecase interface {
	Create(ctx context.Context, input CreateArticleInput) (*Article, error)
	FindByID(ctx context.Context, id int) (*Article, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, article *Article) error
	FindByID(ctx context.Context, id int) (*Article, error)
}

// ToProto convert article to pb.Article
func (a *Article) ToProto() *pb.Article {
	return &pb.Article{
		Id:        int32(a.ID),
		Author:    a.Author,
		Title:     a.Title,
		Body:      a.Body,
		CreatedAt: utils.FormatTimeRFC3339(&a.CreatedAt),
	}
}

// CreateArticleInput create article input
type CreateArticleInput struct {
	Author string `json:"author" validate:"required,min=3,max=60"`
	Title  string `json:"title" validate:"required,min=3,max=60"`
	Body   string `json:"body" validate:"required"`
}

// Validate validate article input body
func (c *CreateArticleInput) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	return nil
}
