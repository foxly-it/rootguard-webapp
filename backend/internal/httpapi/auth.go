package httpapi

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

func RequireBasicAuth(next http.Handler, expectedUser, expectedPassword string) http.Handler {
	expectedUserHash := sha256.Sum256([]byte(expectedUser))
	expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		user, password, ok := r.BasicAuth()
		userHash := sha256.Sum256([]byte(user))
		passwordHash := sha256.Sum256([]byte(password))
		valid := ok &&
			subtle.ConstantTimeCompare(userHash[:], expectedUserHash[:]) == 1 &&
			subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
		if !valid {
			w.Header().Set("WWW-Authenticate", `Basic realm="RootGuard", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
