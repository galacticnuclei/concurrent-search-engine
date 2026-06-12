package robots

import (
	"net/url"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	rules map[string]*Rules
}

func NewCache() *Cache {
	return &Cache{
		rules: make(map[string]*Rules),
	}
}

func (c *Cache) Get(host string) (*Rules, error) {
	c.mu.RLock()
	rules, exists := c.rules[host]
	c.mu.RUnlock()

	if exists {
		return rules, nil
	}

	rules, err := Fetch(host)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.rules[host] = rules
	c.mu.Unlock()

	return rules, nil
}

func (c *Cache) Allowed(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	rules, err := c.Get(parsed.Host)
	if err != nil {
		// If robots.txt can't be fetched,
		// allow crawling for now.
		return true
	}
	return rules.Allowed(parsed.Path)
}
