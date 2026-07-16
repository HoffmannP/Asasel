package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (a *App) runAgentLoop(ctx context.Context) {
	if a.cfg.ControllerURL == "" {
		log.Println("agent mode disabled: no controller URL configured")
		return
	}

	endpoint := strings.TrimRight(a.cfg.ControllerURL, "/") + "/api/control/poll/" + url.PathEscape(a.cfg.AgentID)
	var pending *RemoteResult

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		payload := PollRequest{Result: pending}
		if pending == nil {
			payload.Name = a.cfg.AgentID
		}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			log.Printf("agent poll request error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if a.cfg.SharedSecret != "" {
			req.Header.Set("X-Asasel-Secret", a.cfg.SharedSecret)
		}

		resp, err := a.client.Do(req)
		if err != nil {
			log.Printf("agent poll failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		data, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			log.Printf("agent poll read failed: %v", readErr)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("agent poll status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
			time.Sleep(3 * time.Second)
			continue
		}

		var out PollResponse
		if err := json.Unmarshal(data, &out); err != nil {
			log.Printf("agent poll decode failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		pending = nil
		if out.Command != nil {
			log.Printf("agent poll received command op=%s account=%s id=%s", out.Command.Op, out.Command.Account, out.Command.ID)
			result := executeRemoteCommand(*out.Command)
			pending = &result
		}
	}
}
