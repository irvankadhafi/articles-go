package query

import (
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/olivere/elastic/v7"
)

// BuildSearchArticleQuery :nodoc:
func BuildSearchArticleQuery(req model.SearchArticleRequest) (q *elastic.BoolQuery, sorter elastic.Sorter) {
	var mustQueries []elastic.Query
	var filters []elastic.Query

	if req.Query == "" {
		req.Sorter = model.SortArticleByCreatedAtDesc
	}

	sorter = SortByField(req.Sorter.ToFieldAndDirection())

	if req.Query != "" {
		fields := []string{
			"body^4",
			"title^3",
			"title.autocomplete^2",
			"body.autocomplete",
		}

		mustQueries = append(mustQueries, elastic.NewMultiMatchQuery(req.Query, fields...).
			Operator("or").Type("best_fields"))
	}

	if req.Filter.Author != "" {
		mustQueries = append(mustQueries, elastic.NewMatchQuery("author", req.Filter.Author))
	}

	return elastic.NewBoolQuery().Must(mustQueries...).Filter(filters...), sorter
}

// SortByField :nodoc:
func SortByField(field string, ascending bool) *elastic.FieldSort {
	fieldSort := elastic.NewFieldSort(field)
	fieldSort.Order(ascending)
	return fieldSort
}
