package usecase

import (
	"context"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/irvankadhafi/articles-go/article-service/event"
	"github.com/irvankadhafi/articles-go/article-service/internal/model"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"github.com/sirupsen/logrus"
)

type articleUsecase struct {
	articleRepo model.ArticleRepository
	redisPubSub redigo.Conn
}

// NewArticleUsecase articleUsecase constructor
func NewArticleUsecase(articleRepo model.ArticleRepository, redisPubSub redigo.Conn) model.ArticleUsecase {
	return &articleUsecase{articleRepo: articleRepo, redisPubSub: redisPubSub}
}

func (a *articleUsecase) Create(ctx context.Context, input model.CreateArticleInput) (*model.Article, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	err := input.Validate()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	article := &model.Article{
		Author: input.Author,
		Title:  input.Title,
		Body:   input.Body,
	}

	err = a.articleRepo.Create(ctx, article)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	go func() {
		err := a.publishRedisPubSub("ARTICLE", event.ArticleSubjectAdd, article.ID)
		if err != nil {
			logger.Error(err)
		}
	}()

	return a.FindByID(ctx, article.ID)
}

func (a *articleUsecase) FindByID(ctx context.Context, id int) (*model.Article, error) {
	article, err := a.articleRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithField("id", id).Error(err)
		return nil, err
	}
	if article == nil {
		return nil, ErrNotFound
	}

	return article, nil
}

func (a *articleUsecase) publishRedisPubSub(streamName, eventName string, id int) error {
	e := &event.PubSubPayload{
		Event:     eventName,
		ArticleID: id,
	}
	_, err := a.redisPubSub.Do("PUBLISH", streamName, utils.Dump(e))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"streamName": streamName,
			"e":          utils.Dump(e),
		}).Error(err)

		return err
	}

	return nil
}
