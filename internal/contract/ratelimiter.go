package contract

import (
	"context"

	"github.com/canonflow/canonflow-go-ddd/pkg/response"
)

type RateLimiterContract interface {
	Check(ctx context.Context, identifier string) response.RateLimiterResponse
}
