package interceptors

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"
)

type contextKey struct {
	name string
}

var userJWTCtxKey = &contextKey{"userJWT"}

// User struct
type User struct {
	Username string
	IsAuth   bool
	Error    error
}

// JWTCheck decodes the share session cookie and packs the session into context
func JWTCheck(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userJWT := r.Header.Get("Authorization")
			// Check if header has bearer jwt
			if userJWT == "" {
				user := User{Error: fmt.Errorf("user is not authenticated")}
				ctx := context.WithValue(r.Context(), userJWTCtxKey, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			pl := jwt.Payload{}
			expValidator := jwt.ExpirationTimeValidator(time.Now())
			audValidator := jwt.AudienceValidator(jwt.Audience{"https://ecrpe.fr"})
			validatePayload := jwt.ValidatePayload(&pl, audValidator, expValidator)
			signature := jwt.NewHS256([]byte(secretKey))

			// Split bearer from jwt
			if _, err := jwt.Verify([]byte(strings.Split(userJWT, " ")[1]), signature, &pl, validatePayload); err != nil {
				user := User{Error: fmt.Errorf("JWT is wrong")}
				ctx := context.WithValue(r.Context(), userJWTCtxKey, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			user := User{Username: pl.Subject, IsAuth: true}
			ctx := context.WithValue(r.Context(), userJWTCtxKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userJWTCtxKey).(*User)
	return raw
}
