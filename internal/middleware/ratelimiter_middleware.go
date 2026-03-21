package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(limiter contract.RateLimiterContract, rds *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIp := ctx.ClientIP()

		result := limiter.Check(ctx, userIp)

		//* Attaches all metadata from rate limiter
		for key, value := range result.Metadata {
			headerKey := "X-" + strings.ReplaceAll(strings.Title(key), "_", "-")
			ctx.Header(headerKey, fmt.Sprintf("%v", value))
		}

		if !result.Allow {
			ctx.JSON(http.StatusTooManyRequests, response.BaseErrorResponse{
				Code:   http.StatusTooManyRequests,
				Status: "Too Many Requests",
				Error:  "Too Many Request",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
