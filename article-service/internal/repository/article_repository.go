package repository

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/irvankadhafi/articles-go/article-service/cacher"
	"github.com/irvankadhafi/articles-go/article-service/internal/model"
	"github.com/irvankadhafi/articles-go/article-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type articleRepository struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

// NewArticleRepository articleRepository constructor
func NewArticleRepository(
	db *gorm.DB,
	cacheManager cacher.CacheManager,
) model.ArticleRepository {
	return &articleRepository{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (a *articleRepository) Create(ctx context.Context, article *model.Article) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"article": utils.Dump(article),
	})

	err := a.db.WithContext(ctx).Create(article).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	err = a.cacheManager.DeleteByKeys([]string{a.newCacheByID(article.ID)})
	if err != nil {
		logger.WithField("id", article.ID).Error(err)
	}

	return nil
}

func (a *articleRepository) FindByID(ctx context.Context, id int) (*model.Article, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := a.newCacheByID(id)
	reply, mu, err := a.findFromCacheByKey(cacheKey)
	defer cacher.SafeUnlock(mu)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if mu == nil {
		return reply, nil
	}

	article := model.Article{}
	err = a.db.WithContext(ctx).Take(&article, "id = ?", id).Error
	switch err {
	case nil:
		err = a.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.Dump(article)))
		if err != nil {
			logger.Error(err)
		}
		return &article, nil
	case gorm.ErrRecordNotFound:
		storeNil(a.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

func (a *articleRepository) newCacheByID(id int) string {
	return fmt.Sprintf("cache:object:article:id:%d", id)
}

func (a *articleRepository) findFromCacheByKey(key string) (episode *model.Article, mu *redsync.Mutex, err error) {
	reply, mu, err := a.cacheManager.GetOrLock(key)
	if err != nil || reply == nil {
		return
	}

	episode = utils.InterfaceBytesToType[*model.Article](reply)
	return
}

func storeNil(ck cacher.CacheManager, key string) {
	err := ck.StoreNil(key)
	if err != nil {
		logrus.Error(err)
	}
}
