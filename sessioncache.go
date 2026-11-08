package main

import (
	"container/list"
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	kmf1 "github.com/KohanMathers/KmF1/go"
)

const (
	sessionCacheTTL         = 48 * time.Hour
	defaultSessionCacheSize = 4
	sessionCacheSizeEnvVar  = "KMF1API_SESSION_CACHE_SIZE"
)

type sessionCache struct {
	mu       sync.Mutex
	capacity int
	entries  map[sessionCacheKey]*list.Element
	order    *list.List
}

type sessionCacheKey struct {
	year        int
	round       int
	event       string
	eventStart  time.Time
	sessionType string
}

type sessionCacheEntry struct {
	key       sessionCacheKey
	ready     chan struct{}
	session   *kmf1.Session
	err       error
	expiresAt time.Time
}

func newSessionCache() *sessionCache {
	capacity := defaultSessionCacheSize
	if raw := os.Getenv(sessionCacheSizeEnvVar); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			capacity = n
		}
	}
	return &sessionCache{
		capacity: capacity,
		entries:  make(map[sessionCacheKey]*list.Element),
		order:    list.New(),
	}
}

func isReady(entry *sessionCacheEntry) bool {
	select {
	case <-entry.ready:
		return true
	default:
		return false
	}
}

func (c *sessionCache) evictLocked() {
	for len(c.entries) > c.capacity {
		var victim *list.Element
		for e := c.order.Back(); e != nil; e = e.Prev() {
			if isReady(e.Value.(*sessionCacheEntry)) {
				victim = e
				break
			}
		}
		if victim == nil {
			return
		}
		entry := victim.Value.(*sessionCacheEntry)
		c.order.Remove(victim)
		delete(c.entries, entry.key)
	}
}

func (c *sessionCache) removeLocked(key sessionCacheKey, elem *list.Element) {
	if e, ok := c.entries[key]; ok && e == elem {
		c.order.Remove(elem)
		delete(c.entries, key)
	}
}

func (c *sessionCache) getOrLoad(ctx context.Context, key sessionCacheKey, load func() (*kmf1.Session, error)) (*kmf1.Session, error) {
	for {
		c.mu.Lock()
		elem, ok := c.entries[key]
		if !ok {
			entry := &sessionCacheEntry{key: key, ready: make(chan struct{})}
			elem = c.order.PushFront(entry)
			c.entries[key] = elem
			c.evictLocked()
			c.mu.Unlock()

			entry.session, entry.err = load()
			entry.expiresAt = time.Now().Add(sessionCacheTTL)
			close(entry.ready)

			if entry.err != nil {
				c.mu.Lock()
				c.removeLocked(key, elem)
				c.mu.Unlock()
			}
			return entry.session, entry.err
		}
		entry := elem.Value.(*sessionCacheEntry)
		c.order.MoveToFront(elem)
		c.mu.Unlock()

		select {
		case <-entry.ready:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if entry.err != nil || time.Now().Before(entry.expiresAt) {
			return entry.session, entry.err
		}

		c.mu.Lock()
		c.removeLocked(key, elem)
		c.mu.Unlock()
	}
}
