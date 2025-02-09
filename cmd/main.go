package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	repo "github.com/Hao1995/short-url/internal/adapter/repository/mysql"
	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/internal/router/handler"
	"github.com/Hao1995/short-url/internal/usecase"
	"github.com/Hao1995/short-url/pkg/migrationkit"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/viney-shih/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MIGRATION_DIR = "database/migration"
)

func main() {
	// Migration
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DB,
	)
	if cfg.App.Env == "dev" {
		if err := migrationkit.GooseMigrate(dsn, MIGRATION_DIR); err != nil {
			log.Fatalf("failed to connect to migrate database: %s", map[string]interface{}{"error": err.Error(), "dsn": dsn})
		}
		log.Print("Migrate the DB successfully")
	}

	// Init DB connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("failed to connect to DB: %s", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("failed to get sqlDB from gorm: %s", err)
		} else {
			sqlDB.Close()
		}
	}()
	log.Print("Connect to the DB successfully")

	// Init Cache
	tinyLfu := cache.NewTinyLFU(10000)
	rds := cache.NewRedis(redis.NewRing(&redis.RingOptions{Addrs: cfg.Redis.Addrs}))
	cacheFactory := cache.NewFactory(rds, tinyLfu)

	c := cacheFactory.NewCache([]cache.Setting{
		{
			Prefix: domain.CACHE_PREFIX_SHORT_URL,
			CacheAttributes: map[cache.Type]cache.Attribute{
				cache.SharedCacheType: {TTL: time.Hour},
				cache.LocalCacheType:  {TTL: 10 * time.Second},
			},
			MarshalFunc:   json.Marshal,
			UnmarshalFunc: json.Unmarshal,
		},
	})

	// DI
	repoImpl := repo.NewShortUrlRepository(db)
	ucImpl := usecase.NewShortUrlUseCase(repoImpl, c)
	hlrImpl := handler.NewShortUrlHandler(ucImpl)

	// Run server
	log.Print("Start API server ...")
	if err := endless.ListenAndServe(":"+cfg.App.Port, RegisterGinRouter(hlrImpl)); err != nil {
		log.Fatalf("failed to connect to DB: %s", err)
	}
}

func RegisterGinRouter(hlrImpl *handler.ShortUrlHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/urls", hlrImpl.Create)
	r.GET("/:id", hlrImpl.Get)
	return r
}
