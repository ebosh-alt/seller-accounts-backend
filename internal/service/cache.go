package cache

import (
	"context"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type Repository interface {
	BotLink(ctx context.Context) (string, error)
}

type Storage struct {
	url string
	exp time.Time
	mu  sync.RWMutex
}

func (st *Storage) Get() (string, time.Time) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	return st.url, st.exp
}

func (st *Storage) Set(url string, exp time.Time) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.url = url
	st.exp = exp
}

type Cache struct {
	st    *Storage
	repo  Repository
	ttl   time.Duration
	group singleflight.Group
}

func (c *Cache) GetLink(ctx context.Context) (string, error) {
	var url string
	var exp time.Time
	url, exp = c.st.Get()
	if url != "" && time.Now().Before(exp) {
		return url, nil
	}
	v, err, _ := c.group.Do("bot-link", func() (interface{}, error) {
		url, exp = c.st.Get()
		if url != "" || time.Now().Before(exp) {
			return url, nil
		}
		url, err := c.repo.BotLink(ctx)
		if err != nil {
			return "", err
		}
		c.st.Set(url, time.Now().Add(c.ttl))

		return url, nil
	})
	if err != nil {
		return "", err
	}
	url = v.(string)
	return url, nil
}

func NewCache(repo Repository, ttl time.Duration) *Cache {
	return &Cache{
		st:   &Storage{},
		repo: repo,
		ttl:  ttl,
	}
}
