package basicauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// BasicAuthSingleCredentials wraps an http.HandlerFunc with HTTP Basic Authentication
// using a single set of hardcoded credentials.
//
// If the request Authorization header is missing or the provided credentials do not
// match the expected `username` and `password`, the function responds with a
// 401 Unauthorized status and a "WWW-Authenticate: Basic" header.
//
// Parameters:
//   - handlerFunction: The function to be called after successful authentication.
//   - expectedUsername: The expected username for authentication.
//   - expectedPassword: The expected password for authentication.
func BasicAuthSingleCredentials(handlerFunction http.HandlerFunc, expectedUsername string, expectedPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				// Only if username and password is correct the handler function is called:
				handlerFunction.ServeHTTP(w, r)
				return
			}
		}

		// In any other case a unauthorized is returned.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
