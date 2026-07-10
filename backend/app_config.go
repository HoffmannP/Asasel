package main

import (
	"flag"
	"os"
	"strings"
)

func parseConfig() AppConfig {
	defaultID, _ := os.Hostname()
	if defaultID == "" {
		defaultID = "node"
	}

	mode := flag.String("mode", "local", "Run mode: local|control")
	listen := flag.String("listen", ":2727", "Listen address for HTTP server")
	controllerURL := flag.String("controller-url", "", "Controller base URL for outbound agent polling from local mode, e.g. https://controller.example")
	agentID := flag.String("agent-id", defaultID, "Unique node id used for registration at the controller")
	sharedSecret := flag.String("shared-secret", "", "Shared secret for agent-control polling, sent as X-Asasel-Secret")
	defaultAccount := flag.String("account", "", "Default account used by the web UI (required)")
	authUser := flag.String("auth-user", "", "HTTP basic auth username for web/API access")
	authPass := flag.String("auth-pass", "", "HTTP basic auth password for web/API access")
	servers := flag.String("servers", "", "Comma separated server list for local mode")
	flag.Parse()

	return AppConfig{
		Mode:           strings.ToLower(strings.TrimSpace(*mode)),
		ListenAddr:     *listen,
		ControllerURL:  strings.TrimSpace(*controllerURL),
		AgentID:        strings.TrimSpace(*agentID),
		SharedSecret:   strings.TrimSpace(*sharedSecret),
		DefaultAccount: strings.TrimSpace(*defaultAccount),
		AuthUser:       strings.TrimSpace(*authUser),
		AuthPass:       strings.TrimSpace(*authPass),
		StaticServers:  parseServers(*servers),
	}
}
