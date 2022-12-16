package client

import (
	"context"
	"github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client :nodoc:
type Client struct {
	Conn *grpcpool.Pool
}

// NewClient :nodoc:
func NewClient(target string, timeout time.Duration, idleConnPool, maxConnPool int) (article.ArticleServiceClient, error) {
	factory := newFactory(target, timeout)

	pool, err := grpcpool.New(factory, idleConnPool, maxConnPool, time.Second)
	if err != nil {
		log.Errorf("Error : %v", err)
		return nil, err
	}

	return &Client{
		Conn: pool,
	}, nil
}

func newFactory(target string, timeout time.Duration) grpcpool.Factory {
	return func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(target, grpc.WithInsecure(), withClientUnaryInterceptor(timeout))
		if err != nil {
			log.Errorf("Error : %v", err)
			return nil, err
		}

		return conn, err
	}
}

func withClientUnaryInterceptor(timeout time.Duration) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		ctx = metadata.AppendToOutgoingContext(ctx, "caller", utils.MyCaller(5))
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	})
}
