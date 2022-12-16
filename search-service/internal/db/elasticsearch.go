package db

import (
	"github.com/irvankadhafi/articles-go/search-service/internal/config"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ESInfoLogger :nodoc:
type ESInfoLogger struct{}

// ESErrorLogger :nodoc:
type ESErrorLogger struct{}

// ESTraceLogger :nodoc:
type ESTraceLogger struct{}

// NewElasticsearchClient :nodoc:
func NewElasticsearchClient(url string, httpClient *http.Client) (*elastic.Client, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetScheme("https"),
		elastic.SetSniff(config.ElasticsearchSetSniff()),
		elastic.SetHealthcheck(config.ElasticsearchSetHealthcheck()),
		elastic.SetErrorLog(&ESErrorLogger{}),
		elastic.SetInfoLog(&ESInfoLogger{}),
		elastic.SetHttpClient(httpClient),
	}
	if config.LogLevel() == "debug" {
		opts = append(opts, elastic.SetTraceLog(&ESTraceLogger{}))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

// Printf :nodoc:
func (*ESTraceLogger) Printf(format string, values ...interface{}) {
	log.WithFields(log.Fields{"type": "elasticsearch-log"}).Debugf(format, values...)
}

// Printf :nodoc:
func (*ESInfoLogger) Printf(format string, values ...interface{}) {
	log.WithFields(log.Fields{"type": "elasticsearch-log"}).Infof(format, values...)
}

// Printf :nodoc:
func (*ESErrorLogger) Printf(format string, values ...interface{}) {
	log.WithFields(log.Fields{"type": "elasticsearch-log"}).Errorf(format, values...)
}
