package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const sessionCookieName = "rootguard_session"

type SessionAuth struct {
	expectedUserHash     [32]byte
	expectedPasswordHash [32]byte
	ttl                  time.Duration
	mu                   sync.Mutex
	sessions             map[string]session
	persistencePath      string
}

type session struct {
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSessionAuth(expectedUser, expectedPassword string, ttl time.Duration, persistencePath string) *SessionAuth {
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}
	auth := &SessionAuth{
		expectedUserHash:     sha256.Sum256([]byte(expectedUser)),
		expectedPasswordHash: sha256.Sum256([]byte(expectedPassword)),
		ttl:                  ttl,
		sessions:             make(map[string]session),
		persistencePath:      persistencePath,
	}
	auth.loadSessions()
	return auth
}

func (a *SessionAuth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth/login":
			a.handleLogin(w, r)
			return
		case "/api/auth/logout":
			a.handleLogout(w, r)
			return
		case "/api/auth/session":
			a.handleSession(w, r)
			return
		case "/health":
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/adguard-ui/") {
			if _, ok := a.authenticatedUser(r); !ok {
				writeJSON(w, http.StatusUnauthorized, map[string]any{"authenticated": false})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *SessionAuth) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input credentials
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request"})
		return
	}

	userHash := sha256.Sum256([]byte(input.Username))
	passwordHash := sha256.Sum256([]byte(input.Password))
	valid := subtle.ConstantTimeCompare(userHash[:], a.expectedUserHash[:]) == 1 &&
		subtle.ConstantTimeCompare(passwordHash[:], a.expectedPasswordHash[:]) == 1
	if !valid {
		time.Sleep(250 * time.Millisecond)
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid_credentials"})
		return
	}

	token, err := randomSessionToken()
	if err != nil {
		http.Error(w, "Unable to create session", http.StatusInternalServerError)
		return
	}
	expiresAt := time.Now().Add(a.ttl)
	a.mu.Lock()
	a.deleteExpiredLocked(time.Now())
	a.sessions[token] = session{Username: input.Username, ExpiresAt: expiresAt}
	if err := a.persistLocked(); err != nil {
		delete(a.sessions, token)
		a.mu.Unlock()
		http.Error(w, "Unable to persist session", http.StatusInternalServerError)
		return
	}
	a.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
		MaxAge:   int(a.ttl.Seconds()),
	})
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "username": input.Username})
}

func (a *SessionAuth) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		a.mu.Lock()
		delete(a.sessions, cookie.Value)
		_ = a.persistLocked()
		a.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": false})
}

func (a *SessionAuth) handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	username, ok := a.authenticatedUser(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"authenticated": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "username": username})
}

func (a *SessionAuth) authenticatedUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}
	now := time.Now()
	a.mu.Lock()
	defer a.mu.Unlock()
	entry, ok := a.sessions[cookie.Value]
	if !ok || !entry.ExpiresAt.After(now) {
		delete(a.sessions, cookie.Value)
		_ = a.persistLocked()
		return "", false
	}
	return entry.Username, true
}

func (a *SessionAuth) deleteExpiredLocked(now time.Time) {
	for token, entry := range a.sessions {
		if !entry.ExpiresAt.After(now) {
			delete(a.sessions, token)
		}
	}
}

func (a *SessionAuth) loadSessions() {
	if a.persistencePath == "" {
		return
	}
	data, err := os.ReadFile(a.persistencePath)
	if err != nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_ = json.Unmarshal(data, &a.sessions)
	a.deleteExpiredLocked(time.Now())
}

func (a *SessionAuth) persistLocked() error {
	if a.persistencePath == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(a.persistencePath), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(a.sessions)
	if err != nil {
		return err
	}
	temp := a.persistencePath + ".tmp"
	if err := os.WriteFile(temp, data, 0600); err != nil {
		return err
	}
	return os.Rename(temp, a.persistencePath)
}

func randomSessionToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(buffer)
	if token == "" {
		return "", errors.New("empty session token")
	}
	return token, nil
}

func requestIsHTTPS(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
