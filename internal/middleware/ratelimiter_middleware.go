package middleware

import (
	"fmt"
	"net/http"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(limiter contract.RateLimiterContract, rds *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIp := ctx.ClientIP()

		result := limiter.Check(ctx, userIp)
		fmt.Println(result)

		if !result.Allow {
			ctx.JSON(http.StatusTooManyRequests, response.BaseErrorResponse{
				Code:   http.StatusTooManyRequests,
				Status: "Too Many Requests",
				Error:  "Too Many Request",
			})
			ctx.Abort()
			return
		}
	}
}
