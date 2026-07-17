package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (a *App) registerPublicRoutes(router interface {
	Get(pattern string, handlerFn http.HandlerFunc)
}) {
	router.Get("/quick", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/quick.html", http.StatusTemporaryRedirect)
	})
	router.Get("/api/public/quick", a.quickStatus)
}

func (a *App) quickStatus(w http.ResponseWriter, r *http.Request) {
	account := strings.TrimSpace(r.URL.Query().Get("account"))
	if account == "" {
		account = strings.TrimSpace(a.cfg.DefaultAccount)
	}

	if account == "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"account":   "",
			"remaining": -1,
			"message":   "No account configured",
		})
		return
	}

	if a.cfg.Mode != LocalMode {
		writeJSON(w, http.StatusOK, map[string]any{
			"account":   account,
			"remaining": -1,
			"message":   "Quick status is only available in local mode",
		})
		return
	}

	remaining, err := getTimeout(account)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"account":   account,
			"remaining": -1,
			"message":   "Timeout unreadable",
		})
		return
	}

	msg := "Timeout not found"
	if remaining >= 0 {
		msg = fmt.Sprintf("Timeout is %d", remaining)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"account":   account,
		"remaining": remaining,
		"message":   msg,
	})
}
