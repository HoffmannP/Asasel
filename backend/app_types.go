package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type AppConfig struct {
	Mode           string
	ListenAddr     string
	ControllerURL  string
	AgentID        string
	SharedSecret   string
	DefaultAccount string
	AuthUser       string
	AuthPass       string
	StaticServers  []string
}

type RemoteCommand struct {
	ID        string `json:"id"`
	Op        string `json:"op"`
	Account   string `json:"account,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	LockState *bool  `json:"lockstate,omitempty"`
}

type RemoteResult struct {
	ID        string `json:"id"`
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	LockState *bool  `json:"lockstate,omitempty"`
	Duration  *int   `json:"duration,omitempty"`
	Remaining *int   `json:"remaining,omitempty"`
	Error     string `json:"error,omitempty"`
}

type PollRequest struct {
	Name   string        `json:"name,omitempty"`
	Result *RemoteResult `json:"result,omitempty"`
}

type PollResponse struct {
	Command *RemoteCommand `json:"command,omitempty"`
}

type AgentState struct {
	ID       string
	Name     string
	LastSeen time.Time
	Pending  chan RemoteCommand
	Waiters  map[string]chan RemoteResult
}

type ControllerState struct {
	mu     sync.Mutex
	agents map[string]*AgentState
}

type App struct {
	cfg        AppConfig
	controller ControllerState
	client     *http.Client
}

func newApp(cfg AppConfig) *App {
	return &App{
		cfg: cfg,
		controller: ControllerState{
			agents: make(map[string]*AgentState),
		},
		client: &http.Client{Timeout: 40 * time.Second},
	}
}

func parseServers(input string) []string {
	parts := strings.Split(input, ",")
	servers := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			servers = append(servers, trimmed)
		}
	}
	if len(servers) == 0 {
		return []string{"yoga", "nura", "hackmack", "bereisen"}
	}
	return servers
}

func (a *App) listServers() []string {
	if a.cfg.Mode != "control" {
		return a.cfg.StaticServers
	}
	now := time.Now()
	a.controller.mu.Lock()
	defer a.controller.mu.Unlock()
	servers := make([]string, 0)
	for id, state := range a.controller.agents {
		if now.Sub(state.LastSeen) <= 2*time.Minute {
			servers = append(servers, id)
		}
	}
	sort.Strings(servers)
	return servers
}

func (a *App) ensureAgent(id string) *AgentState {
	state, ok := a.controller.agents[id]
	if !ok {
		state = &AgentState{
			ID:       id,
			Name:     id,
			LastSeen: time.Now(),
			Pending:  make(chan RemoteCommand, 32),
			Waiters:  make(map[string]chan RemoteResult),
		}
		a.controller.agents[id] = state
	}
	return state
}

func newCommandID() string {
	return fmt.Sprintf("%d-%x", time.Now().UnixNano(), rand.Uint64())
}

func (a *App) queueAndAwait(agentID string, cmd RemoteCommand) (RemoteResult, error) {
	a.controller.mu.Lock()
	state, ok := a.controller.agents[agentID]
	if !ok {
		a.controller.mu.Unlock()
		return RemoteResult{}, fmt.Errorf("agent %s unknown", agentID)
	}
	if time.Since(state.LastSeen) > 2*time.Minute {
		a.controller.mu.Unlock()
		return RemoteResult{}, fmt.Errorf("agent %s offline", agentID)
	}
	respCh := make(chan RemoteResult, 1)
	state.Waiters[cmd.ID] = respCh
	a.controller.mu.Unlock()

	select {
	case state.Pending <- cmd:
	case <-time.After(2 * time.Second):
		a.controller.mu.Lock()
		delete(state.Waiters, cmd.ID)
		a.controller.mu.Unlock()
		return RemoteResult{}, errors.New("agent command queue full")
	}

	select {
	case result := <-respCh:
		return result, nil
	case <-time.After(35 * time.Second):
		a.controller.mu.Lock()
		delete(state.Waiters, cmd.ID)
		a.controller.mu.Unlock()
		return RemoteResult{}, errors.New("agent did not answer in time")
	}
}
