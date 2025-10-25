package main

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found in cache")
	ErrInvalidKey  = errors.New("invalid key: cannot be empty")
	ErrInvalidTTL  = errors.New("invalid TTL: must be positive duration")
	ErrCacheClosed = errors.New("cache is closed")
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, bool, error)
	Delete(key string) error
}

type Closer interface {
	Close() error
}

type SimpleCache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string]interface{}),
	}
}

func validateKey(key string) error {
	if key == "" {
		return ErrInvalidKey
	}
	return nil
}

func (c *SimpleCache) Set(key string, value interface{}) error {
	if err := validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *SimpleCache) Get(key string) (interface{}, bool, error) {
	if err := validateKey(key); err != nil {
		return nil, false, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.data[key]
	return value, exists, nil
}

func (c *SimpleCache) Delete(key string) error {
	if err := validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

type cacheItem struct {
	value   interface{}
	expires time.Time
}

type TTLCache struct {
	data            map[string]cacheItem
	mu              sync.RWMutex
	ttl             time.Duration
	cleanupInterval time.Duration
	ticker          *time.Ticker
	done            chan struct{}
	closed          bool
}

type TTLCacheConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
}

func NewTTLCache(ttl time.Duration) (*TTLCache, error) {
	return NewTTLCacheWithConfig(TTLCacheConfig{
		TTL:             ttl,
		CleanupInterval: 0,
	})
}

func NewTTLCacheWithConfig(config TTLCacheConfig) (*TTLCache, error) {
	if config.TTL <= 0 {
		return nil, ErrInvalidTTL
	}

	cleanupInterval := config.CleanupInterval
	if cleanupInterval <= 0 {
		cleanupInterval = config.TTL / 2
		if cleanupInterval < 100*time.Millisecond {
			cleanupInterval = 100 * time.Millisecond
		}
	}

	cache := &TTLCache{
		data:            make(map[string]cacheItem),
		ttl:             config.TTL,
		cleanupInterval: cleanupInterval,
		ticker:          time.NewTicker(cleanupInterval),
		done:            make(chan struct{}),
		closed:          false,
	}

	go cache.cleanup()

	return cache, nil
}

func (c *TTLCache) Set(key string, value interface{}) error {
	if c.closed {
		return ErrCacheClosed
	}

	if err := validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}
	return nil
}

func (c *TTLCache) Get(key string) (interface{}, bool, error) {
	if c.closed {
		return nil, false, ErrCacheClosed
	}

	if err := validateKey(key); err != nil {
		return nil, false, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false, nil
	}

	if time.Now().After(item.expires) {
		return nil, false, nil
	}

	return item.value, true, nil
}

func (c *TTLCache) Delete(key string) error {
	if c.closed {
		return ErrCacheClosed
	}

	if err := validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *TTLCache) cleanup() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	for {
		select {
		case <-c.ticker.C:
			c.performCleanup()
		case <-c.done:
			c.ticker.Stop()
			return
		}
	}
}

func (c *TTLCache) performCleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expires) {
			delete(c.data, key)
		}
	}
}

func (c *TTLCache) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true
	close(c.done)
	return nil
}

func (c *TTLCache) IsClosed() bool {
	return c.closed
}
