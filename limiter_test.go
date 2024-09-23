package limiter

import (
	"math/rand"
	"testing"
	"time"
)

// TODO: Better testing.
func TestLimiter(t *testing.T) {
	l := NewLimiter(time.Duration(30)*time.Second, 2, 10)

	// Mock incoming requests
	ips := [3]string{"ip1", "ip2", "ip3"}
	ipTracker := make(map[string]int, 3)
	ipTracker["ip1"] = 0
	ipTracker["ip2"] = 0
	ipTracker["ip3"] = 0

	endTime := time.Now().Add(time.Duration(30) * time.Second)
	for time.Now().Unix() < endTime.Unix() {
		ip := ips[rand.Intn(len(ips))]
		existing := ipTracker[ip]
		res := l.Allowed(ip)

		if existing <= 12 && res == 200 {
			ipTracker[ip] += 1
			continue
		} else if res != 429 {
			t.Fatal("Failed")
		}
	}
}
