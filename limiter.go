package limiter

import (
	"sync"
	"time"
)

// Parent struct for limiter with config options.
type Limiter struct {
	IPMap map[string]*IPCtx
	mu    sync.Mutex
	Limit int
	Burst int
	Rate  time.Duration // time.Duration(1) * time.Minute
}

// The context for a specific IP adress
type IPCtx struct {
	LastTime time.Time
	NextTime time.Time
	Amount   int
	Burst    int
}

// Returns a new *Limiter with the given options.
func NewLimiter(rate time.Duration, limit int, burst int) *Limiter {
	return &Limiter{
		Rate:  rate,
		Limit: limit,
		Burst: burst,
		IPMap: make(map[string]*IPCtx, 0),
	}
}

// Returns an HTTP status code letting you know if a current IP is allowed for a given request.
func (l *Limiter) Allowed(ip string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()

	ctx, ok := l.IPMap[ip]
	// If ip doesnt exist create and allow
	if !ok {
		nextTime := now.Add(l.Rate)
		l.IPMap[ip] = &IPCtx{
			LastTime: time.Now(),
			NextTime: nextTime,
			Amount:   l.Limit,
			Burst:    l.Burst - 1,
		}
		return 200
	}

	// If burst hasnt been fulfilled, allow
	if ctx.Burst > 0 {
		ctx.Burst -= 1
		return 200
	}

	// If the nextTime interval has passed, reset
	if ctx.NextTime.Unix() < now.Unix() {
		nextTime := now.Add(l.Rate)
		ctx.LastTime = now
		ctx.NextTime = nextTime
		ctx.Amount = l.Limit
		ctx.Burst = l.Burst - 1

		return 200
	}

	// If we are within the interval span and still have remaining requests, allow
	if ctx.NextTime.Unix() > now.Unix() && ctx.Amount > 0 {
		ctx.Amount -= 1
		return 200
	}

	// Otherwise reject.
	if ctx.NextTime.Unix() > now.Unix() && ctx.Amount <= 0 {
		return 429
	}

	return 500
}
