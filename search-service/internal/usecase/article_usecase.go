package usecase

import (
	"context"
	articleEvent "github.com/irvankadhafi/articles-go/article-service/event"
	articleSvc "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/irvankadhafi/articles-go/search-service/event"
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/irvankadhafi/articles-go/search-service/utils"
	log "github.com/sirupsen/logrus"
	"sync"
)

type articleUsecase struct {
	articleRepo   model.ArticleRepository
	articleClient articleSvc.ArticleServiceClient
}

// NewArticleUsecase articleUsecase constructor
func NewArticleUsecase(articleRepo model.ArticleRepository, articleClient articleSvc.ArticleServiceClient) model.ArticleSearchUsecase {
	return &articleUsecase{articleRepo: articleRepo, articleClient: articleClient}
}

func (a *articleUsecase) IndexByID(ctx context.Context, id int) error {
	article, err := a.articleClient.FindArticleByID(ctx, &articleSvc.FindByIDRequest{Id: int64(id)})
	if err != nil {
		log.Error(err)
		return err
	}

	articleFromProto := model.NewArticleFromProto(article)
	return a.articleRepo.IndexES(ctx, articleFromProto)
}

func (a *articleUsecase) HandleEvent(ctx context.Context, msg *event.PubSubPayload) error {
	switch msg.Event {
	case articleEvent.ArticleSubjectAdd:
		return a.IndexByID(ctx, msg.ArticleID)
	}
	return nil
}

func (a *articleUsecase) Search(ctx context.Context, req model.SearchArticleRequest) (spaces []*model.Article, count int, err error) {
	logger := log.WithFields(log.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"searchCriteria": utils.Dump(req),
	})

	ids, count, err := a.articleRepo.Search(ctx, req)
	if err != nil {
		logger.Error(err)
		return
	}

	if len(ids) == 0 || count == 0 {
		return
	}

	spaces = a.FindAllByIDs(ctx, ids)
	if len(spaces) <= 0 {
		err = ErrNotFound
		return
	}

	return
}

// FindAllByIDs :nodoc:
func (a *articleUsecase) FindAllByIDs(ctx context.Context, ids []int) (articles []*model.Article) {
	logger := log.WithFields(log.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"ids": ids,
	})

	var wg sync.WaitGroup
	c := make(chan *model.Article, len(ids))
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			article, err := a.articleClient.FindArticleByID(ctx, &articleSvc.FindByIDRequest{Id: int64(id)})
			if err != nil {
				logger.Error(err)
				return
			}

			c <- model.NewArticleFromProto(article)
		}(id)
	}
	wg.Wait()
	close(c)

	if len(c) <= 0 {
		return
	}

	rs := map[int]*model.Article{}
	for article := range c {
		if article != nil {
			rs[article.ID] = article
		}
	}

	for _, id := range ids {
		if article, ok := rs[id]; ok {
			articles = append(articles, article)
		}
	}

	return
}
