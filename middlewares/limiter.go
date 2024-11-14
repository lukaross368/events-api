package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var Clients = make(map[string]*Client)
var Mu sync.Mutex

func RateLimit(context *gin.Context) {

	ip := context.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = context.ClientIP()
	}

	if client := getClientFromIp(ip); !client.Limiter.Allow() {
		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
		return
	}

	context.Next()
}

func getClientFromIp(ip string) *Client {
	Mu.Lock()
	client, exists := Clients[ip]

	if !exists {

		client = &Client{
			Limiter:  rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
			LastSeen: time.Now(),
		}

		Clients[ip] = client
	}

	client.LastSeen = time.Now()
	Mu.Unlock()
	return client
}
