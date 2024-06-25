package ratelimiter

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

type RedisRateLimiter struct {
	Client  *redis.Client
	Limiter *redis_rate.Limiter
}

func NewRedisRateLimiter() *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis", err)
		return nil
	}
	fmt.Println("Redis is connected", pong)

	// Initialize Redis rate limiter
	limiter := redis_rate.NewLimiter(client)

	return &RedisRateLimiter{
		Client:  client,
		Limiter: limiter,
	}
}

// middleware

func RunLimit(rl *RedisRateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limit := redis_rate.PerMinute(3)
			res, err := rl.Limiter.Allow(context.Background(), "global", limit)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if res.Allowed == 0 {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				log.Println("Too many requests - rate limit exceeded.")
				return
			}
			fmt.Println("middleware: request sent , requests left:", res.Remaining)
			next.ServeHTTP(w, r)
		})
	}
}
