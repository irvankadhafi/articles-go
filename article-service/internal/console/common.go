package console

import (
	"context"
	"errors"
	goredis "github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/imdario/mergo"
	"github.com/irvankadhafi/articles-go/article-service/internal/config"
	"github.com/irvankadhafi/articles-go/article-service/internal/db"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// RedisConnectionPoolOptions options for the redis connection
type RedisConnectionPoolOptions struct {
	// Dial timeout for establishing new connections.
	// Default is 5 seconds. Only for go-redis.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds. Only for go-redis.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout. Only for go-redis.
	WriteTimeout time.Duration

	// Number of idle connections in the pool.
	IdleCount int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	PoolSize int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
}

var defaultRedisConnectionPoolOptions = &RedisConnectionPoolOptions{
	IdleCount:       20,
	PoolSize:        100,
	IdleTimeout:     60 * time.Second,
	MaxConnLifetime: 0,
	DialTimeout:     5 * time.Second,
	WriteTimeout:    2 * time.Second,
	ReadTimeout:     2 * time.Second,
}

// NewRedigoRedisConnectionPool uses redigo library to establish the redis connection pool
func NewRedigoRedisConnectionPool(url string, opt *RedisConnectionPoolOptions) (*redigo.Pool, error) {
	if !isValidRedisStandaloneURL(url) {
		return nil, errors.New("invalid redis URL: " + url)
	}

	options := applyRedisConnectionPoolOptions(opt)

	return &redigo.Pool{
		MaxIdle:     options.IdleCount,
		MaxActive:   options.PoolSize,
		IdleTimeout: options.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		MaxConnLifetime: options.MaxConnLifetime,
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Wait: true, // wait for connection available when maxActive is reached
	}, nil
}

func applyRedisConnectionPoolOptions(opt *RedisConnectionPoolOptions) *RedisConnectionPoolOptions {
	if opt == nil {
		return defaultRedisConnectionPoolOptions
	}
	// if error occurs, also return options from input
	_ = mergo.Merge(opt, *defaultRedisConnectionPoolOptions)
	return opt
}

func isValidRedisStandaloneURL(url string) bool {
	_, err := goredis.ParseURL(url)
	return err == nil
}

func serverInterceptor(ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, config.RPCServerTimeout())
	defer cancel()
	return handler(ctx, req)
}

func gracefulShutdown(grpcSvr *grpc.Server, graphqlSvr *echo.Echo) {
	db.StopTickerCh <- true

	if grpcSvr != nil {
		grpcSvr.GracefulStop()
	}

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
