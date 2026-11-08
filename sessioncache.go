package main

import (
	"context"
	"sync"
	"time"

	kmf1 "github.com/KohanMathers/KmF1/go"
)

const sessionCacheTTL = 48 * time.Hour

type sessionCache struct {
	mu      sync.Mutex
	entries map[sessionCacheKey]*sessionCacheEntry
}

type sessionCacheKey struct {
	year        int
	round       int
	event       string
	eventStart  time.Time
	sessionType string
}

type sessionCacheEntry struct {
	ready     chan struct{}
	session   *kmf1.Session
	err       error
	expiresAt time.Time
}

func newSessionCache() *sessionCache {
	return &sessionCache{entries: make(map[sessionCacheKey]*sessionCacheEntry)}
}

func (c *sessionCache) getOrLoad(ctx context.Context, key sessionCacheKey, load func() (*kmf1.Session, error)) (*kmf1.Session, error) {
	for {
		c.mu.Lock()
		entry, ok := c.entries[key]
		if !ok {
			entry = &sessionCacheEntry{ready: make(chan struct{})}
			c.entries[key] = entry
			c.mu.Unlock()

			entry.session, entry.err = load()
			entry.expiresAt = time.Now().Add(sessionCacheTTL)
			close(entry.ready)

			if entry.err != nil {
				c.mu.Lock()
				if c.entries[key] == entry {
					delete(c.entries, key)
				}
				c.mu.Unlock()
			}
			return entry.session, entry.err
		}
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
		if c.entries[key] == entry {
			delete(c.entries, key)
		}
		c.mu.Unlock()
	}
}
