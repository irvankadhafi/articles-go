package console

import (
	"crypto/tls"
	"fmt"
	articleSvc "github.com/irvankadhafi/articles-go/article-service/client"
	"github.com/irvankadhafi/articles-go/search-service/internal/config"
	"github.com/irvankadhafi/articles-go/search-service/internal/db"
	"github.com/irvankadhafi/articles-go/search-service/internal/delivery/httpsvc"
	"github.com/irvankadhafi/articles-go/search-service/internal/helper"
	"github.com/irvankadhafi/articles-go/search-service/internal/repository"
	"github.com/irvankadhafi/articles-go/search-service/internal/usecase"
	"github.com/irvankadhafi/articles-go/search-service/pubsub"
	"github.com/irvankadhafi/articles-go/search-service/subscriber"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	esClient, err := db.NewElasticsearchClient(config.ElasticsearchHost(), &http.Client{
		Timeout: config.ElasticsearchHTTPTimeout(),
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
			MaxIdleConnsPerHost: config.ElasticsearchMaxIdleConnections(),
			MaxConnsPerHost:     config.ElasticsearchMaxConnsPerHost(),
		},
	})
	continueOrFatal(err)

	articleSvcClient, err := articleSvc.NewClient(config.ArticleServiceTarget(), config.RPCClientTimeout(), config.ServiceIdleConnPool(), config.ServiceMaxConnPool())
	continueOrFatal(err)

	redisOpts := &pubsub.RedisConnectionPoolOptions{
		DialTimeout:     config.RedisDialTimeout(),
		ReadTimeout:     config.RedisReadTimeout(),
		WriteTimeout:    config.RedisWriteTimeout(),
		IdleCount:       config.RedisMaxIdleConn(),
		PoolSize:        config.RedisMaxActiveConn(),
		IdleTimeout:     240 * time.Second,
		MaxConnLifetime: 1 * time.Minute,
	}

	redisSubscriber, err := pubsub.NewRedisPubSub(config.RedisPubSubHost(), redisOpts)
	continueOrFatal(err)
	defer helper.WrapCloser(redisSubscriber.Conn.Close)

	articleRepo := repository.NewArticleRepository(esClient)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, articleSvcClient)

	// Create articleEventSubscriber
	articleEventSubscriber := subscriber.NewRedisSubscriber(redisSubscriber.Conn)
	articleEventSubscriber.RegisterArticleUsecase(articleUsecase)

	go func() {
		// Subscribe to channel
		if err := articleEventSubscriber.Subscribe("ARTICLE"); err != nil {
			log.Fatal(err)
		}

		// Consume messages from channel
		if err := articleEventSubscriber.Consume(); err != nil {
			log.Fatal(err)
		}
	}()

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
				gracefulShutdown(httpServer)
				quitCh <- true
			case e := <-errCh:
				log.Error(e)
				gracefulShutdown(httpServer)
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

	<-quitCh
	log.Info("exiting")
}
