package core

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const CoreKey = "CORE_KEY"

type Core struct {
	db *bolt.DB
	c  *cache.Cache
	entityId string
	accessToken string
	apiBaseUrl string
}

func Initialize(db *bolt.DB, c *cache.Cache, entityId, accessToken string, apiBaseUrl string) *Core {
	return &Core{db: db, c: c, entityId: entityId, accessToken: accessToken, apiBaseUrl: apiBaseUrl}
}

func (c *Core) SetCore() gin.HandlerFunc {
	return func(gc *gin.Context) {
		gc.Set(CoreKey, c)
		gc.Next()
	}
}
