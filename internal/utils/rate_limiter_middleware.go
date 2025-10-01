package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterUtil struct {
	Redis      *redis.Client
	MaxRequest int64
	Duration   time.Duration
}

func NewRateLimiterUtil(redis *redis.Client) *RateLimiterUtil {
	return &RateLimiterUtil{
		Redis:      redis,
		MaxRequest: 1,
		Duration:   time.Second * 10,
	}
}

func (u RateLimiterUtil) IsAllowed(ctx context.Context, token string) bool {

	increment, err := u.Redis.Incr(ctx, token).Result()
	if err != nil {
		fmt.Println("Error incrementing:", err)
		return false
	}

	if increment == 1 {
		err := u.Redis.Expire(ctx, token, u.Duration).Err()
		if err != nil {
			fmt.Println("Error setting expiration:", err)
			return false
		}
	}

	return increment <= u.MaxRequest
}
