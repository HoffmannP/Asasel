package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (a *App) registerControlRoutes(router chi.Router) {
	router.Route("/api/control", func(r chi.Router) {
		r.Post("/poll/{agent}", func(w http.ResponseWriter, r *http.Request) {
			if !a.requireSecret(w, r) {
				return
			}

			agentID := chi.URLParam(r, "agent")
			if agentID == "" {
				http.Error(w, "missing agent id", http.StatusBadRequest)
				return
			}

			var in PollRequest
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}

			a.controller.mu.Lock()
			state := a.ensureAgent(agentID)
			state.LastSeen = time.Now()
			if in.Name != "" {
				state.Name = in.Name
			}
			if in.Result != nil {
				if waiter, ok := state.Waiters[in.Result.ID]; ok {
					waiter <- *in.Result
					delete(state.Waiters, in.Result.ID)
				}
			}
			a.controller.mu.Unlock()

			select {
			case cmd := <-state.Pending:
				writeJSON(w, http.StatusOK, PollResponse{Command: &cmd})
			case <-time.After(25 * time.Second):
				writeJSON(w, http.StatusOK, PollResponse{})
			}
		})
	})

	router.Route("/api/remote/{server}", func(r chi.Router) {
		r.Get("/accounts/lock/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteLockGet(w, r)
		})
		r.Post("/accounts/lock/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteLockPost(w, r)
		})
		r.Get("/accounts/time/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteTimeGet(w, r)
		})
		r.Get("/accounts/killall/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteKillallGet(w, r)
		})
		r.Get("/timeouts/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteTimeoutGet(w, r)
		})
		r.Post("/timeouts/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteTimeoutPost(w, r)
		})
		r.Delete("/timeouts/{account}", func(w http.ResponseWriter, r *http.Request) {
			a.remoteTimeoutDelete(w, r)
		})
	})
}
