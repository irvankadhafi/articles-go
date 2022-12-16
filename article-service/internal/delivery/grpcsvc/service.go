package grpcsvc

import (
	"github.com/irvankadhafi/articles-go/article-service/cacher"
	"github.com/irvankadhafi/articles-go/article-service/internal/model"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
)

// Service :nodoc:
type Service struct {
	pb.UnimplementedArticleServiceServer
	cacheManager   cacher.CacheManager
	articleUsecase model.ArticleUsecase
}

// NewService :nodoc:
func NewService() *Service {
	return new(Service)
}

// RegisterCacheManager :nodoc:
func (s *Service) RegisterCacheManager(k cacher.CacheManager) {
	s.cacheManager = k
}

// RegisterArticleUsecase :nodoc:
func (s *Service) RegisterArticleUsecase(au model.ArticleUsecase) {
	s.articleUsecase = au
}
