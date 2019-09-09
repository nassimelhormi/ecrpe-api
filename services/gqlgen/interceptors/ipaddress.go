package interceptors

import (
	"context"
	"net/http"
)

type UserIPAddressContextKey struct {
	name string
}

var userIPAddressCtxKey = &UserIPAddressContextKey{"userIPAddress"}

type IPAddress struct {
	string
}

// GetIPAddress decodes the share session cookie and packs the session into context
func GetIPAddress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			IPAddress := IPAddress{r.RemoteAddr}
			ctx := context.WithValue(r.Context(), userIPAddressCtxKey, IPAddress)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForIPAddress finds the user from the context. REQUIRES Middleware to have run.
func ForIPAddress(ctx context.Context) string {
	IP := ctx.Value(userIPAddressCtxKey).(*IPAddress)
	return IP.string
}
