package limiter

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	counts   map[string]int // method -> active count
	lastSeen time.Time      // когда последний раз был запрос
}

type Limiter struct {
	mu       sync.Mutex
	limits   map[string]int    // method -> limit
	clients  map[string]*entry // clientID -> entry
	ttl      time.Duration     // через сколько чистим "мертвых"
	stopChan chan struct{}
}

func NewLimiter(limits map[string]int, ttl time.Duration) *Limiter {
	l := &Limiter{
		limits:   limits,
		clients:  make(map[string]*entry),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}

	// запускаем чистильщик в фоне
	go l.cleanupLoop()
	return l
}

func (l *Limiter) Inc(clientID, method string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	limit, ok := l.limits[method]
	if !ok {
		return true
	}

	e, exists := l.clients[clientID]
	if !exists {
		e = &entry{counts: make(map[string]int)}
		l.clients[clientID] = e
	}

	fmt.Printf("client=%s method=%s active=%d limit=%d\n",
		clientID, method, e.counts[method], limit)

	if e.counts[method] >= limit {
		return false
	}

	e.counts[method]++
	e.lastSeen = time.Now()
	return true
}

func (l *Limiter) Dec(clientID, method string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, exists := l.clients[clientID]
	if !exists {
		return
	}

	e.counts[method]--
	if e.counts[method] < 0 {
		e.counts[method] = 0
	}
	e.lastSeen = time.Now()
}

// cleanupLoop фоновая чистка неактивных клиентов
func (l *Limiter) cleanupLoop() {
	ticker := time.NewTicker(l.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.cleanup()
		case <-l.stopChan:
			return
		}
	}
}

func (l *Limiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for clientID, e := range l.clients {
		if now.Sub(e.lastSeen) > l.ttl {
			delete(l.clients, clientID)
		}
	}
}

func (l *Limiter) Stop() {
	close(l.stopChan)
}
