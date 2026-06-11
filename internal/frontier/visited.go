package frontier

import "sync"

type Visited struct {
	mu   sync.RWMutex
	urls map[string]struct{}
}

func NewVisited() *Visited {
	return &Visited{
		urls: make(map[string]struct{}),
	}
}

func (v *Visited) Seen(url string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, exists := v.urls[url]
	return exists
}

func (v *Visited) Add(url string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.urls[url] = struct{}{}
}

func (v *Visited) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.urls)
}
