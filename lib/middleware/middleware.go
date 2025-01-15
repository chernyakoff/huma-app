package middleware

import (
	"huma-app/lib/security"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func JwtAuthMiddleware(api huma.API, security *security.Security) func(ctx huma.Context, next func(huma.Context)) {
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

		claims, err := security.ParseToken(token.Value)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx = huma.WithValue(ctx, "user_id", claims.UserID)
		ctx = huma.WithValue(ctx, "is_admin", claims.IsAdmin)

		if len(anyOfNeededRoles) == 0 {
			next(ctx)
			return
		}

		for _, role := range anyOfNeededRoles {
			if role == "admin" && claims.IsAdmin {
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
