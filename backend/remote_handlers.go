package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *App) controlForwardTarget() (string, error) {
	if a.cfg.Mode != "control" {
		return "", errors.New("control forwarding only available in control mode")
	}
	servers := a.listServers()
	if len(servers) == 0 {
		return "", errors.New("no active agent available")
	}
	return servers[0], nil
}

func (a *App) controlForward(cmd RemoteCommand) (RemoteResult, error) {
	agentID, err := a.controlForwardTarget()
	if err != nil {
		return RemoteResult{}, err
	}
	cmd.ID = newCommandID()
	return a.queueAndAwait(agentID, cmd)
}

func (a *App) dispatchRemote(w http.ResponseWriter, r *http.Request, cmd RemoteCommand) (RemoteResult, bool) {
	if a.cfg.Mode != "control" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Remote control only available in control mode"})
		return RemoteResult{}, false
	}
	agentID := chi.URLParam(r, "server")
	if agentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Missing server id"})
		return RemoteResult{}, false
	}

	cmd.ID = newCommandID()
	result, err := a.queueAndAwait(agentID, cmd)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"message": err.Error()})
		return RemoteResult{}, false
	}
	if !result.OK {
		writeJSON(w, http.StatusBadGateway, map[string]string{"message": result.Message})
		return RemoteResult{}, false
	}
	return result, true
}

func (a *App) remoteStateGet(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "get_state", Account: account})
	if !ok {
		return
	}
	lockState := false
	if result.LockState != nil {
		lockState = *result.LockState
	}
	duration := -1
	if result.Duration != nil {
		duration = *result.Duration
	}
	remaining := -1
	if result.Remaining != nil {
		remaining = *result.Remaining
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message":   result.Message,
		"lockstate": lockState,
		"duration":  duration,
		"remaining": remaining,
	})
}

func (a *App) remoteConfigGet(w http.ResponseWriter, r *http.Request) {
	if a.cfg.Mode != "control" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Remote config only available in control mode"})
		return
	}
	agentID := chi.URLParam(r, "server")
	if agentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Missing server id"})
		return
	}

	a.controller.mu.Lock()
	state, ok := a.controller.agents[agentID]
	a.controller.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "agent unknown"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"account": state.Account,
	})
}

func (a *App) remoteLockGet(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "get_lock", Account: account})
	if !ok {
		return
	}
	lockState := false
	if result.LockState != nil {
		lockState = *result.LockState
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message, "lockstate": lockState})
}

func (a *App) remoteLockPost(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	var in struct {
		LockState bool `json:"lockstate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON"})
		return
	}
	lock := in.LockState
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "set_lock", Account: account, LockState: &lock})
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message})
}

func (a *App) remoteTimeGet(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "get_time", Account: account})
	if !ok {
		return
	}
	duration := -1
	if result.Duration != nil {
		duration = *result.Duration
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message, "duration": duration})
}

func (a *App) remoteKillallPost(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "killall", Account: account})
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message})
}

func (a *App) remoteTimeoutGet(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "get_timeout", Account: account})
	if !ok {
		return
	}
	remaining := -1
	if result.Remaining != nil {
		remaining = *result.Remaining
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message, "remaining": remaining})
}

func (a *App) remoteTimeoutPost(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	var in struct {
		Duration int `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid JSON"})
		return
	}
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "set_timeout", Account: account, Duration: in.Duration})
	if !ok {
		return
	}
	remaining := in.Duration
	if result.Remaining != nil {
		remaining = *result.Remaining
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message, "remaining": remaining})
}

func (a *App) remoteTimeoutDelete(w http.ResponseWriter, r *http.Request) {
	account := chi.URLParam(r, "account")
	result, ok := a.dispatchRemote(w, r, RemoteCommand{Op: "del_timeout", Account: account})
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": result.Message})
}
