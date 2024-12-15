package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisCache provides a Redis-based caching mechanism
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// CacheConfig holds configuration for Redis cache
type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(config CacheConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set stores a value in the cache with a specified expiration
func (rc *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	// Convert value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal cache value")
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	// Store in Redis
	err = rc.client.Set(rc.ctx, key, jsonData, expiration).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":        key,
			"expiration": expiration,
			"error":      err,
		}).Error("Failed to set cache value")
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"key":        key,
		"expiration": expiration,
	}).Debug("Cached value set successfully")

	return nil
}

// Get retrieves a value from the cache
func (rc *RedisCache) Get(key string, dest interface{}) error {
	// Retrieve from Redis
	result, err := rc.client.Get(rc.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Cache miss is not an error, just return nil
			return nil
		}
		logrus.WithFields(logrus.Fields{
			"key":   key,
			"error": err,
		}).Error("Failed to retrieve cache value")
		return fmt.Errorf("failed to retrieve cache value: %w", err)
	}

	// Unmarshal JSON
	err = json.Unmarshal(result, dest)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal cache value")
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Delete removes a key from the cache
func (rc *RedisCache) Delete(key string) error {
	err := rc.client.Del(rc.ctx, key).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":   key,
			"error": err,
		}).Error("Failed to delete cache key")
		return fmt.Errorf("failed to delete cache key: %w", err)
	}

	logrus.WithField("key", key).Debug("Cache key deleted successfully")
	return nil
}

// Exists checks if a key exists in the cache
func (rc *RedisCache) Exists(key string) (bool, error) {
	count, err := rc.client.Exists(rc.ctx, key).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":   key,
			"error": err,
		}).Error("Failed to check cache key existence")
		return false, fmt.Errorf("failed to check cache key existence: %w", err)
	}

	return count > 0, nil
}

// Invalidate handles cache invalidation for specific patterns
func (rc *RedisCache) Invalidate(pattern string) (int64, error) {
	// Find and delete keys matching the pattern
	keys, err := rc.client.Keys(rc.ctx, pattern).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pattern": pattern,
			"error":   err,
		}).Error("Failed to find keys for invalidation")
		return 0, fmt.Errorf("failed to find keys for invalidation: %w", err)
	}

	// Delete found keys
	if len(keys) > 0 {
		deletedCount, err := rc.client.Del(rc.ctx, keys...).Result()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pattern": pattern,
				"error":   err,
			}).Error("Failed to delete keys")
			return 0, fmt.Errorf("failed to delete keys: %w", err)
		}

		logrus.WithFields(logrus.Fields{
			"pattern":       pattern,
			"deleted_count": deletedCount,
		}).Debug("Cache keys invalidated")

		return deletedCount, nil
	}

	return 0, nil
}

// Close terminates the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}