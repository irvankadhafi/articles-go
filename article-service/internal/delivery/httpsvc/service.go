package httpsvc

import (
	"github.com/irvankadhafi/articles-go/article-service/internal/model"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Service http service
type Service struct {
	route          *echo.Echo
	articleUsecase model.ArticleUsecase
}

// RouteService add dependencies and use route for routing
func RouteService(
	route *echo.Echo,
	usecase model.ArticleUsecase,
) {
	srv := &Service{
		route:          route,
		articleUsecase: usecase,
	}
	srv.initRoutes()
}

func (s *Service) initRoutes() {
	// create episode
	s.route.POST("/articles", s.handleCreateArticle())
	s.route.GET("/woy", s.handleWoy())

}

func (s *Service) handleWoy() echo.HandlerFunc {
	logrus.Warn("MASUK KE HANDLER")
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "HALOO TEMAN")
	}
}

func (s *Service) handleCreateArticle() echo.HandlerFunc {
	type request struct {
		Author string `json:"author"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	return func(c echo.Context) error {
		req := request{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		newArticle, err := s.articleUsecase.Create(c.Request().Context(), model.CreateArticleInput{
			Author: req.Author,
			Title:  req.Title,
			Body:   req.Body,
		})
		if err != nil {
			logrus.Error(err)
			return err
		}

		res := model.Article{
			ID:        newArticle.ID,
			Author:    newArticle.Author,
			Title:     newArticle.Title,
			Body:      newArticle.Body,
			CreatedAt: newArticle.CreatedAt,
		}

		return c.JSON(http.StatusOK, res)
	}
}
