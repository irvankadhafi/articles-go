package repository

import (
	"context"
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/irvankadhafi/articles-go/search-service/internal/repository/query"
	"github.com/irvankadhafi/articles-go/search-service/utils"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strconv"
)

type articleRepo struct {
	es *elastic.Client
}

// NewArticleRepository :nodoc:
func NewArticleRepository(es *elastic.Client) model.ArticleRepository {
	return &articleRepo{
		es: es,
	}
}

func (a *articleRepo) Search(ctx context.Context, req model.SearchArticleRequest) (ids []int, count int, err error) {
	query, sorter := query.BuildSearchArticleQuery(req)
	result, err := a.es.Search().
		Index(model.ESArticleIndexAlias).
		Query(query).
		FetchSource(false).
		TrackTotalHits(true).
		From(int(utils.Offset(req.Page, req.Size))).
		Size(int(req.Size)).
		SortBy(sorter).
		Do(ctx)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": utils.DumpIncomingContext(ctx),
			"req":     utils.Dump(req)}).
			Error(err)
		return
	}

	for _, hit := range result.Hits.Hits {
		ids = append(ids, utils.StringToInt(hit.Id))
	}

	return ids, int(result.TotalHits()), nil
}

func (a *articleRepo) IndexES(ctx context.Context, article *model.Article) error {
	_, err := a.es.Index().
		Index(model.ESArticleIndex).
		Id(strconv.Itoa(article.ID)).
		BodyJson(article).
		Do(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context":   utils.DumpIncomingContext(ctx),
			"indexData": utils.Dump(article),
		}).Error(err)
	}

	return nil
}
