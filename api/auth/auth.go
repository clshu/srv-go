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
				log.Println(err.Error())
				next.ServeHTTP(w, r)
				return
			}

			tokenId, err := utils.ParseToken(token)
			if err != nil {
				// Token is bad, don't put it in the context
				// The resolvers requiring token authentication will fail
				log.Println(err.Error())
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
func ForContext(ctx context.Context) *utils.TokenId {
	raw, _ := ctx.Value(userCtxKey).(*utils.TokenId)
	return raw
}
