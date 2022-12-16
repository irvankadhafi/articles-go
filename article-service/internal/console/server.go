package console

import (
	"fmt"
	"github.com/irvankadhafi/articles-go/article-service/cacher"
	"github.com/irvankadhafi/articles-go/article-service/internal/config"
	"github.com/irvankadhafi/articles-go/article-service/internal/db"
	"github.com/irvankadhafi/articles-go/article-service/internal/delivery/grpcsvc"
	"github.com/irvankadhafi/articles-go/article-service/internal/delivery/httpsvc"
	"github.com/irvankadhafi/articles-go/article-service/internal/helper"
	"github.com/irvankadhafi/articles-go/article-service/internal/repository"
	"github.com/irvankadhafi/articles-go/article-service/internal/usecase"
	pb "github.com/irvankadhafi/articles-go/article-service/pb/article"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// Initiate all connection like db, redis, etc
	db.InitializePostgresConn()

	pgDB, err := db.PostgreSQL.DB()
	continueOrFatal(err)
	defer helper.WrapCloser(pgDB.Close)

	redisOpts := &RedisConnectionPoolOptions{
		DialTimeout:     config.RedisDialTimeout(),
		ReadTimeout:     config.RedisReadTimeout(),
		WriteTimeout:    config.RedisWriteTimeout(),
		IdleCount:       config.RedisMaxIdleConn(),
		PoolSize:        config.RedisMaxActiveConn(),
		IdleTimeout:     240 * time.Second,
		MaxConnLifetime: 1 * time.Minute,
	}

	redisConn, err := NewRedigoRedisConnectionPool(config.RedisCacheHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisConn.Close)

	redisLockConn, err := NewRedigoRedisConnectionPool(config.RedisLockHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisLockConn.Close)

	redisPubSub, err := NewRedigoRedisConnectionPool(config.RedisPubSubHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisPubSub.Close)

	cacheManager := cacher.NewCacheManager()
	cacheManager.SetConnectionPool(redisConn)
	cacheManager.SetLockConnectionPool(redisLockConn)
	cacheManager.SetDefaultTTL(config.CacheTTL())

	articleRepository := repository.NewArticleRepository(db.PostgreSQL, cacheManager)
	articleUsecase := usecase.NewArticleUsecase(articleRepository, redisPubSub.Get())

	grpcSvc := grpc.NewServer(grpc.ChainUnaryInterceptor(
		serverInterceptor,
	))

	httpServer := echo.New()
	httpServer.Pre(middleware.RemoveTrailingSlash())
	httpServer.Use(middleware.Logger())
	httpServer.Use(middleware.Recover())
	httpServer.Use(middleware.CORS())

	httpsvc.RouteService(httpServer, articleUsecase)

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	quitCh := make(chan bool, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for {
			select {
			case <-sigCh:
				gracefulShutdown(grpcSvc, httpServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(grpcSvc, httpServer)
				quitCh <- true
			}
		}
	}()

	go func() {
		// Start HTTP server
		if err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort())); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go func() {
		// Service definition
		svc := grpcsvc.NewService()
		svc.RegisterCacheManager(cacheManager)
		svc.RegisterArticleUsecase(articleUsecase)

		pb.RegisterArticleServiceServer(grpcSvc, svc)

		// Start gRPC server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCPort()))
		if err != nil {
			log.WithField("port", config.GRPCPort()).Fatalf("failed to listen: %v", err)
		}

		log.Info("Listening on ", config.GRPCPort())

		errCh <- grpcSvc.Serve(lis)
	}()

	<-quitCh
	log.Info("exiting")
}
