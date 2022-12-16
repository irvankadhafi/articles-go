package console

import (
	"context"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"time"
)

//func serverInterceptor(ctx context.Context,
//	req interface{},
//	_ *grpc.UnaryServerInfo,
//	handler grpc.UnaryHandler,
//) (interface{}, error) {
//	ctx, cancel := context.WithTimeout(ctx, config.RPCServerTimeout())
//	defer cancel()
//	return handler(ctx, req)
//}

func gracefulShutdown(graphqlSvr *echo.Echo) {
	if graphqlSvr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := graphqlSvr.Shutdown(ctx); err != nil {
			graphqlSvr.Logger.Fatal(err)
		}
	}
}

func continueOrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
