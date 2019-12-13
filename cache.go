package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Cache service
type Cache struct {
	db         *sql.DB
	expiration int
}

// Find finds from cache
func (cache *Cache) Find(c context.Context, q string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(2*time.Second))
	defer cancel()

	normalizedQuery := NormalizeQuery(q)

	var id int
	var query string
	var createdAt time.Time
	var entry []byte

	row := cache.db.QueryRowContext(ctx, "SELECT id, query, result, created_at FROM cache WHERE query = $1", normalizedQuery)
	switch err := row.Scan(&id, &query, &entry, &createdAt); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		now := time.Now()
		if now.After(createdAt.Add(time.Duration(cache.expiration) * time.Second)) {
			log.Debug(fmt.Sprintf("found expired result for query \"%s\"", normalizedQuery))
			return nil, nil
		}

		var result []map[string]interface{}
		err := json.Unmarshal(entry, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New("failed to fetch data from cache")
	}
}

// Create does something
func (cache *Cache) Create(c context.Context, q string, entry []map[string]interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(c, time.Duration(2*time.Second))
	defer cancel()

	raw, err := json.Marshal(entry)
	if err != nil {
		return false, err
	}

	normalizedQuery := NormalizeQuery(q)
	_, err = cache.db.ExecContext(ctx, "INSERT INTO cache (query, result, created_at) VALUES ($1, $2, now())", normalizedQuery, raw)
	if err != nil {
		log.Error(err)
		return false, err
	}

	log.Debug(fmt.Sprintf("created entry in cache for query \"%s\"", q))
	return true, nil
}

// NewCache constructor
func NewCache(db *sql.DB) *Cache {
	seconds := viper.GetInt("database.cache_expiration")
	return &Cache{db: db, expiration: seconds}
}
