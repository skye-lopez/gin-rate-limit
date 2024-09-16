package limiter

// Parent struct for limiter with config options.
type Limiter struct{}

// Returns a new *Limiter with the given options.
func NewLimiter() *Limiter {
	return &Limiter{}
}
