package middleware

import (
	"huma-app/lib/security"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/time/rate"
)

func getRealIP(ctx huma.Context) string {
	var ip string
	if tcip := ctx.Header(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := ctx.Header(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := ctx.Header(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func RealIpMiddleware(ctx huma.Context, next func(huma.Context)) {
	ip := getRealIP(ctx)
	if ip == "" {
		ip = ctx.RemoteAddr()
	}
	ctx = huma.WithValue(ctx, "real-ip", ip)
	next(ctx)
}

// Структура для хранения лимитеров по IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter создает новый IPRateLimiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter возвращает лимитер для конкретного IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Если лимитер для IP уже существует, возвращаем его
	if limiter, exists := i.ips[ip]; exists {
		return limiter
	}

	// Создаем новый лимитер для IP
	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func RateLimitMiddleware(api huma.API, limiter *IPRateLimiter) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ip := ctx.Context().Value("real-ip").(string)
		limiter := limiter.GetLimiter(ip)
		if !limiter.Allow() {
			huma.WriteErr(api, ctx, http.StatusTooManyRequests, "To many requests")
			return
		}
		next(ctx)
	}
}

func JwtAuthMiddleware(api huma.API, sec *security.Security) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		var anyOfNeededRoles []string
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if anyOfNeededRoles, ok = opScheme["Bearer"]; ok {
				isAuthorizationRequired = true
			}
		}

		if !isAuthorizationRequired {
			next(ctx)
			return
		}

		token, err := huma.ReadCookie(ctx, "jwt")
		if err != nil || len(token.Value) == 0 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, err := sec.VerifyToken(token.Value, security.AccessToken)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx = huma.WithValue(ctx, "user_id", claims.UserID)
		ctx = huma.WithValue(ctx, "user_role", claims.UserRole)

		if len(anyOfNeededRoles) == 0 {
			next(ctx)
			return
		}

		for _, role := range anyOfNeededRoles {
			if role == string(claims.UserRole) {
				next(ctx)
				return
			}
		}

		huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden")
	}
}

func RequestLoggerMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

	}
}
