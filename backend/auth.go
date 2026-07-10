package main

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func (a *App) requireSecret(w http.ResponseWriter, r *http.Request) bool {
	if a.cfg.SharedSecret == "" {
		return true
	}
	if r.Header.Get("X-Asasel-Secret") != a.cfg.SharedSecret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func (a *App) requireBasicAuth(next http.Handler) http.Handler {
	if a.cfg.AuthUser == "" && a.cfg.AuthPass == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Agent polling has dedicated shared-secret auth and must not require browser auth.
		if strings.HasPrefix(r.URL.Path, "/api/control/poll/") {
			next.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if ok {
			userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(a.cfg.AuthUser)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(a.cfg.AuthPass)) == 1
			if userMatch && passMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Asasel"`)
		http.Error(w, "authentication required", http.StatusUnauthorized)
	})
}
