package client

import (
	"context"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"google.golang.org/grpc"
)

func (c *Client) FindArticleByID(ctx context.Context, in *pb.FindByIDRequest, opts ...grpc.CallOption) (*pb.Article, error) {
	conn, err := c.Conn.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewArticleServiceClient(conn.ClientConn)
	result, err := client.FindArticleByID(ctx, in, opts...)

	return result, err
}
