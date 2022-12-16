package httpsvc

import (
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/irvankadhafi/articles-go/search-service/utils"
	"github.com/labstack/echo"
	"math"
	"net/http"
	"strconv"
)

// Service http service
type Service struct {
	route          *echo.Echo
	articleUsecase model.ArticleSearchUsecase
}

// RouteService add dependencies and use route for routing
func RouteService(
	route *echo.Echo,
	usecase model.ArticleSearchUsecase,
) {
	srv := &Service{
		route:          route,
		articleUsecase: usecase,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	s.route.GET("/articles", s.handleSearchArticle())
}

type CursorInfo struct {
	Size      int    `json:"size"`
	Count     int    `json:"count"`
	CountPage int    `json:"countPage"`
	HasMore   bool   `json:"hasMore"`
	Cursor    string `json:"cursor"`
}
type searchResponse struct {
	Articles   []*model.Article `json:"articles"`
	CursorInfo *CursorInfo      `json:"cursor_info"`
}

func (s *Service) handleSearchArticle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		query := c.QueryParam("query")
		page := c.QueryParam("page")
		size := c.QueryParam("size")
		author := c.QueryParam("author")

		criterias := model.SearchArticleRequest{
			Query: query,
			Page:  utils.StringToInt(page),
			Size:  utils.StringToInt(size),
		}

		if author != "" {
			criterias.Filter.Author = author
		}

		articles, count, err := s.articleUsecase.Search(ctx, criterias)
		if err != nil {
			return err
		}

		hasMore := count-(criterias.Page*criterias.Size) > 0
		return c.JSON(http.StatusOK, &searchResponse{
			Articles: articles,
			CursorInfo: &CursorInfo{
				Size:      criterias.Size,
				Count:     count,
				CountPage: int(int64(math.Ceil(float64(count) / float64(criterias.Size)))),
				HasMore:   hasMore,
				Cursor:    strconv.Itoa(criterias.Page),
			},
		})
	}
}
