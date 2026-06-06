package processes

import (
	"log"
	"time"

	"github.com/lukaross368/events-api/middlewares"
)

const staleThreshold = 3 * time.Minute

func CleanUpOldLimiters() {
	middlewares.Mu.Lock()
	defer middlewares.Mu.Unlock()

	now := time.Now()

	for ip, client := range middlewares.Clients {
		if now.Sub(client.LastSeen) > staleThreshold {
			delete(middlewares.Clients, ip)
			log.Printf("Removed stale rate limiter for IP: %s\n", ip)
		}
	}
}
