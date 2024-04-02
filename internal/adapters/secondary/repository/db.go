package repository

import (
	"todo-hexagonal/internal/adapters/secondary/cache"

	"gorm.io/gorm"
)

type DB struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

// new database
func NewDB(db *gorm.DB, cache *cache.RedisCache) *DB {
	return &DB{
		db:    db,
		cache: cache,
	}
}
