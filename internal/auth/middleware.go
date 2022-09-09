package auth

import (
	"net/http"
	"server/internal/util/response"
)

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())

		if err != nil {
			response.Unauthorized(w, r)
			return
		}

		if token == nil || !token.Valid {
			response.Unauthorized(w, r)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
