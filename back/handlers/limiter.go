package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type entry struct {
	prevTime  time.Time
	prevCount float64
	currCount float64
}

func (e *entry) slide(windowSize time.Duration) {
	now := time.Now()
	dt := now.Sub(e.prevTime)

	if dt > windowSize*3 {
		e.prevTime = now.Truncate(windowSize).Add(-windowSize)
		e.prevCount = 0.0
		e.currCount = 0.0

		return
	}

	if dt > windowSize*2 {
		e.prevTime = e.prevTime.Add(windowSize)
		e.prevCount = e.currCount
		e.currCount = 0.0
	}
}

func (e *entry) currRatePerWindow(windowSize time.Duration) float64 {
	e.slide(windowSize)

	now := time.Now()
	tt := now.Truncate(windowSize)
	fc := now.Sub(tt).Hours() / windowSize.Hours()
	fp := tt.Sub(now.Add(-windowSize)).Hours() / windowSize.Hours()

	return fp*e.prevCount + fc*e.currCount
}

// Limiter sliding window request limiter based on ip
type Limiter struct {
	rwlock     sync.RWMutex
	m          map[string]entry
	windowSize time.Duration
	maxRate    float64
}

func (l *Limiter) incGetRate(k string) float64 {
	l.rwlock.Lock()
	defer l.rwlock.Unlock()

	now := time.Now()
	e, ok := l.m[k]

	if !ok {
		l.m[k] = entry{
			prevTime:  now.Truncate(l.windowSize).Add(-l.windowSize),
			prevCount: 0.0,
			currCount: 1.0,
		}

		return 1.0
	}

	e.slide(l.windowSize)
	e.currCount += 1.0
	l.m[k] = e

	return e.currRatePerWindow(l.windowSize)
}

// NewLimiter create a sliding window limiter
func NewLimiter(cleanEvery, windowSize time.Duration, maxRate float64) *Limiter {
	l := &Limiter{
		rwlock:     sync.RWMutex{},
		m:          make(map[string]entry),
		windowSize: windowSize,
		maxRate:    maxRate,
	}

	go doEvery(cleanEvery, func() {
		// iterate over map and clean it. Minimal retention
		var keysToDelete []string
		now := time.Now()
		l.rwlock.RLock()
		for k, e := range l.m {
			if now.Sub(e.prevTime) > 3*windowSize {
				keysToDelete = append(keysToDelete, k)
			}
		}
		l.rwlock.RUnlock()
		if len(keysToDelete) > 0 {
			l.rwlock.Lock()
			for _, k := range keysToDelete {
				delete(l.m, k)
			}
			l.rwlock.Unlock()
		}

	})

	return l
}

// Handler rate limitting middleware
func (l *Limiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ise(w, err)
			return
		}

		rate := l.incGetRate(ip)
		if rate > l.maxRate {
			limited(w)
			return
		}
		// not limited
		next.ServeHTTP(w, r)
	})
}

func doEvery(d time.Duration, fn func()) {
	fn()

	for range time.Tick(d) {
		fn()
	}
}
