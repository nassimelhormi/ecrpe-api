package interceptors

import (
	"context"
	"net/http"
)

type RefreshTokenContextKey struct {
	name string
}

var refreshTokenCtxKey = &RefreshTokenContextKey{"refreshToken"}

type RefreshToken struct {
	string
}

// GetRefreshToken decodes the share session cookie and packs the session into context
func GetRefreshToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			refreshToken := r.Header.Get("X-Refresh-Token")
			ctx := context.WithValue(r.Context(), refreshTokenCtxKey, refreshToken)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForRefreshToken finds the user from the context. REQUIRES Middleware to have run.
func ForRefreshToken(ctx context.Context) string {
	refreshToken := ctx.Value(refreshTokenCtxKey).(*RefreshToken)
	return refreshToken.string
}
