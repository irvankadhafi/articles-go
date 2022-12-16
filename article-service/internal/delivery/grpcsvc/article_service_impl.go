package grpcsvc

import (
	"context"
	"github.com/irvankadhafi/articles-go/article-service/internal/usecase"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FindArticleByID :nodoc:
func (s *Service) FindArticleByID(ctx context.Context, in *pb.FindByIDRequest) (out *pb.Article, err error) {
	article, err := s.articleUsecase.FindByID(ctx, int(in.GetId()))
	switch err {
	case nil:
		return article.ToProto(), nil
	case usecase.ErrNotFound:
		return nil, status.Error(codes.NotFound, "not found")
	default:
		logrus.WithFields(logrus.Fields{
			"ctx": utils.DumpIncomingContext(ctx),
			"req": utils.Dump(in),
		}).Error(err)
		return nil, status.Error(codes.Internal, "something wrong")
	}
}
