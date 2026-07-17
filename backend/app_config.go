package main

import (
	"flag"
	"os"
	"strings"
)

func parseCerts(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.SplitN(input, ",", 2)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func parseAuth(input string) (string, string) {
	if input == "" {
		return "", ""
	}
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		panic("invalid auth format, expected user:pass")
	}
	return parts[0], parts[1]
}

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
	auth := flag.String("auth", "", "HTTP basic auth username:password for web/API access")
	certs := flag.String("certs", "", "Comma separated key/certificate for SSL")
	flag.Parse()

	authUser, authPass := parseAuth(strings.TrimSpace(*auth))

	return AppConfig{
		Mode:           strings.ToLower(strings.TrimSpace(*mode)),
		ListenAddr:     *listen,
		ControllerURL:  strings.TrimSpace(*controllerURL),
		AgentID:        strings.TrimSpace(*agentID),
		SharedSecret:   strings.TrimSpace(*sharedSecret),
		DefaultAccount: strings.TrimSpace(*defaultAccount),
		AuthUser:       authUser,
		AuthPass:       authPass,
		Certs:          parseCerts(*certs),
	}
}
