package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/clshu/srv-go/utils"
)

// The middleware to get token in the http header
// and sav it in the context

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"tokenId"}
var errCtxKey = &contextKey{"tokenIdError"}

type contextKey struct {
	name string
}

// A stand-in for token
// type User struct {
// 	token utils.TokenId
// }

func extractBearerToken(auths []string) (string, error) {
	var tokens []string
	for _, auth := range auths {
		// log.Println(auth)
		strs := strings.Split(auth, " ")
		bearer := strings.Trim(strs[0], " ")
		btoken := strings.Trim(strs[1], " ")
		// log.Println(bearer)
		// log.Println(btoken)
		// It's supposed to be Bearer, some are not.
		// Normalize it to compare
		if strings.ToLower(bearer) == "bearer" {
			tokens = append(tokens, btoken)
		}
	}
	// log.Println(tokens)
	size := len(tokens)
	// log.Println(size)
	switch {
	case size == 0:
		return "", fmt.Errorf("%s", "No Bearer Token")
	case size > 1:
		return "", fmt.Errorf("%s", "Too Many Bearer Token")
	default:
		return tokens[0], nil
	}
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auths := r.Header.Values("Authorization")
			// Allow unauthenticated users in
			if len(auths) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			token, err := extractBearerToken(auths)
			if err != nil {
				// No Bearer token or too many bearer token
				log.Printf("auth.middleware: %s\n", err.Error())
				// Let resolvers return the error
				ctx := context.WithValue(r.Context(), errCtxKey, err)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
				return
			}

			tokenId, err := utils.ParseToken(token)
			if err != nil {
				// Token is bad, don't put it in the context
				// The resolvers requiring token authentication will fail

				log.Printf("auth.middleware: %s\n", err.Error())
				// Let resolvers return the error
				ctx := context.WithValue(r.Context(), errCtxKey, err)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
				return
			}
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, tokenId)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the tokenId from the context. REQUIRES Middleware to have run.
func forContext(ctx context.Context) *utils.TokenId {
	raw, _ := ctx.Value(userCtxKey).(*utils.TokenId)
	return raw
}
func forContextError(ctx context.Context) error {
	raw, _ := ctx.Value(errCtxKey).(error)
	return raw
}
func GetTokenIdFromCtx(ctx context.Context) (*utils.TokenId, error) {
	err := forContextError(ctx)
	if err != nil {
		return nil, err
	}
	return forContext(ctx), nil
}
