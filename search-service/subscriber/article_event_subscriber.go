package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/irvankadhafi/articles-go/search-service/event"
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ArticleEventSubscriber is an interface for subscribing to Redis channel
type ArticleEventSubscriber interface {
	Subscribe(channel string) error
	RegisterArticleUsecase(articleUsecase model.ArticleSearchUsecase)
	Consume() error
}

// articleEventSubscriber is a concrete implementation of ArticleEventSubscriber using Redis
type articleEventSubscriber struct {
	conn           redis.Conn
	psc            redis.PubSubConn
	articleUsecase model.ArticleSearchUsecase
}

// NewRedisSubscriber creates a new articleEventSubscriber
func NewRedisSubscriber(conn redis.Conn) ArticleEventSubscriber {
	return &articleEventSubscriber{conn: conn}
}

// RegisterArticleUsecase :nodoc:
func (r *articleEventSubscriber) RegisterArticleUsecase(articleUsecase model.ArticleSearchUsecase) {
	r.articleUsecase = articleUsecase
}

// Subscribe subscribes to a Redis channel
func (r *articleEventSubscriber) Subscribe(channel string) error {
	r.psc = redis.PubSubConn{Conn: r.conn}
	if err := r.psc.Subscribe(channel); err != nil {
		return errors.Wrap(err, "failed to subscribe to channel")
	}
	return nil
}

// Consume consumes messages from the Redis channel
func (r *articleEventSubscriber) Consume() error {
	for {
		switch v := r.psc.Receive().(type) {
		case redis.Message:
			var article event.PubSubPayload
			if err := json.Unmarshal(v.Data, &article); err != nil {
				logrus.Error(err)
				return err
			}
			// Handle Reindex with usecase
			err := r.articleUsecase.HandleEvent(context.Background(), &article)
			if err != nil {
				logrus.Error(err)
				return err
			}

		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return errors.Wrap(v, "failed to consume message")
		}
	}
}
