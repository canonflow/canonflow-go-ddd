package factory

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/canonflow/canonflow-go-ddd/pkg/utils"
	"github.com/redis/go-redis/v9"
)

var (
	ALGORITHMS               = []string{"token_bucket"}
	TOKEN_BUCKET_REFILL_TIME = 1 * time.Minute
)

type tokenBucketRateLimiter struct {
	Rdb        *redis.Client
	Limit      int
	RefillTime time.Duration
	Mu         sync.Mutex
}

func newTokenBucketRateLimiter(rdb *redis.Client, limit int) *tokenBucketRateLimiter {
	return &tokenBucketRateLimiter{
		Rdb:        rdb,
		Limit:      limit,
		RefillTime: TOKEN_BUCKET_REFILL_TIME,
	}
}

func (limiter *tokenBucketRateLimiter) Check(ctx context.Context, identifier string) response.RateLimiterResponse {
	limiter.Mu.Lock()
	defer limiter.Mu.Unlock()

	key := "rate_limiter:" + identifier
	now := time.Now().Unix()
	refillSec := int64(limiter.RefillTime.Seconds())

	//* Read current state from Redis Hash
	data, err := limiter.Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	//* Parse Tokens
	tokens := limiter.Limit
	if val, ok := data["tokens"]; ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			tokens = parsed
		}
	}

	//* Parse last_refill
	lastRefill := now
	if val, ok := data["last_refill"]; ok {
		if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
			lastRefill = parsed
		}
	}

	//* Refill tokens based on elapsed refill intervals
	elapsed := now - lastRefill
	intervals := elapsed / refillSec
	if intervals > 0 {
		tokens = min(limiter.Limit, tokens+int(intervals))
		lastRefill = lastRefill + (intervals + refillSec)
	}

	//* Check and consume token
	allowed := false
	if tokens > 0 {
		tokens--
		allowed = true
	}

	//* Persist update state
	if err := limiter.Rdb.HSet(ctx, key, map[string]interface{}{
		"tokens":      tokens,
		"last_refill": lastRefill,
	}).Err(); err != nil {
		panic(err)
	}

	//* Auto cleanup: TTL = 2 x refill interval
	limiter.Rdb.Expire(ctx, key, time.Duration(refillSec*2)*time.Second)

	return response.RateLimiterResponse{
		Allow: allowed,
		Metadata: map[string]interface{}{
			"tokens":      tokens,
			"last_refill": lastRefill,
		},
	}
}

func NewRateLimiterFactory(algorithm string, rdb *redis.Client, limit int) contract.RateLimiterContract {
	if !utils.SliceContains(ALGORITHMS, strings.ToLower(algorithm)) {
		return nil
	}

	if algorithm == "token_bucket" {
		return newTokenBucketRateLimiter(rdb, limit)
	}

	return nil
}
