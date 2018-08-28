package render

import "sync"

type toggle struct {
	active bool
	mu     sync.RWMutex
}

func (r *toggle) Set(a bool) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.active == a {
		return false
	}
	r.active = a
	return true
}

func (r *toggle) State() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.active
}
