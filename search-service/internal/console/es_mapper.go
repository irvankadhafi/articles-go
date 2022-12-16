package console

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/irvankadhafi/articles-go/search-service/internal/config"
	"github.com/irvankadhafi/articles-go/search-service/internal/db"
	"github.com/irvankadhafi/articles-go/search-service/internal/model"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

var (
	mapperCmd = &cobra.Command{
		Use:   "map",
		Short: "run mapper for elasticsearch",
		Long:  `run mapper for elasticsearch`,
		Run:   putAllMapping,
	}

	indexes = []string{
		model.ESArticleIndex,
	}

	indexMappingMap = map[string]string{
		model.ESArticleIndex: model.ESArticleMapping,
	}

	indexAliasMap = map[string]string{
		model.ESArticleIndex: model.ESArticleIndexAlias,
	}
)

func init() {
	RootCmd.AddCommand(mapperCmd)
}

func putMapping(client *elastic.Client, index string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Error(err)
	}

	if exists {
		return
	}

	// Create a new index.
	createIndex, err := client.CreateIndex(index).BodyString(indexMappingMap[index]).Do(ctx)
	if err != nil {
		log.WithField("index", index).Error(err)
		return
	}
	if !createIndex.Acknowledged {
		log.Error(fmt.Errorf("create index %s not acknowledged", index))
		return
	}

	// add alias for the new index
	createAlias, err := client.Alias().
		Add(index, indexAliasMap[index]).
		Do(ctx)
	if err != nil {
		log.WithFields(log.Fields{"alias": indexAliasMap[index], "index": index}).Error(err)
		return
	}
	if !createAlias.Acknowledged {
		log.Error(fmt.Errorf("add alias %s for index %s not acknowledged", indexAliasMap[index], index))
	}
}

func putAllMapping(_ *cobra.Command, _ []string) {
	log.Info("Creating connection to elasticsearch...")
	esClient, err := db.NewElasticsearchClient(config.ElasticsearchHost(), &http.Client{
		Timeout: 200 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		},
	})
	if err != nil {
		log.Error(err)
		return
	}

	for _, index := range indexes {
		log.Infof("Putting mapping for %s", index)
		putMapping(esClient, index)

		log.Infof("Done putting mapping for %s", index)
	}
}
