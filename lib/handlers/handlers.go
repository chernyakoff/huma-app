package handlers

import (
	"huma-app/lib/middleware"
	"huma-app/lib/security"
	"huma-app/store"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Login  *middleware.IPRateLimiter
	Verify *middleware.IPRateLimiter
}

type ApiHandlers struct {
	repo     *store.Queries
	security *security.Security
	limiter  *RateLimiter
}

func NewApiHandlers(repo *store.Queries, security *security.Security) *ApiHandlers {

	return &ApiHandlers{repo, security, &RateLimiter{
		Login:  middleware.NewIPRateLimiter(rate.Every(time.Minute), 2),
		Verify: middleware.NewIPRateLimiter(rate.Every(time.Minute), 1),
	}}
}
