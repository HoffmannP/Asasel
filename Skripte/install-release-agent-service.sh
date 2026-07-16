#!/usr/bin/env bash
set -euo pipefail

REPO_DEFAULT="HoffmannP/Asasel"
ASSET_NAME="Asasel"
BIN_PATH="/usr/local/bin/Asasel"
WRAPPER_PATH="/usr/local/bin/asasel-agent-wrapper"
SERVICE_NAME="asasel"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
CONF_DIR="/etc/asasel"
ENV_FILE="${CONF_DIR}/asasel.env"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Fehlt: $1" >&2
    exit 1
  fi
}

prompt_nonempty() {
  local label="$1"
  local value=""
  while [[ -z "$value" ]]; do
    read -r -p "$label" value
  done
  printf '%s' "$value"
}

prompt_default() {
  local label="$1"
  local def="$2"
  local value=""
  read -r -p "$label [$def]: " value
  if [[ -z "$value" ]]; then
    value="$def"
  fi
  printf '%s' "$value"
}

require_cmd curl
require_cmd sudo
require_cmd systemctl
require_cmd install

REPO=$(prompt_default "GitHub Repo (owner/name)" "$REPO_DEFAULT")
CONTROLLER_URL=$(prompt_nonempty "Controller URL (z.B. https://controller.example): ")
ACCOUNT=$(prompt_nonempty "Linux Account fuer diesen Agent (z.B. linus): ")
LISTEN=$(prompt_default "Listen Address" "127.0.0.1:2828")
AGENT_ID=$(prompt_default "Agent ID" "$(hostname)")

read -r -p "Shared Secret (optional, Enter = leer): " SHARED_SECRET
read -r -p "Basic Auth User (optional, Enter = leer): " AUTH_USER
if [[ -n "$AUTH_USER" ]]; then
  read -r -s -p "Basic Auth Pass: " AUTH_PASS
  echo
else
  AUTH_PASS=""
fi

TMP_BIN=$(mktemp)
trap 'rm -f "$TMP_BIN"' EXIT

echo "Lade neuestes Release von $REPO ..."
curl -fL "https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}" -o "$TMP_BIN"

sudo install -d -m 755 "$CONF_DIR"
sudo install -m 755 "$TMP_BIN" "$BIN_PATH"

# Environment-Datei fuer Service-Konfiguration.
sudo tee "$ENV_FILE" >/dev/null <<EOF
ASASEL_LISTEN=${LISTEN}
ASASEL_CONTROLLER_URL=${CONTROLLER_URL}
ASASEL_AGENT_ID=${AGENT_ID}
ASASEL_ACCOUNT=${ACCOUNT}
ASASEL_SHARED_SECRET=${SHARED_SECRET}
ASASEL_AUTH_USER=${AUTH_USER}
ASASEL_AUTH_PASS=${AUTH_PASS}
EOF
sudo chmod 600 "$ENV_FILE"

# Wrapper baut optionale Argumente sauber zusammen.
sudo tee "$WRAPPER_PATH" >/dev/null <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

args=(
  -mode local
  -listen "${ASASEL_LISTEN:-127.0.0.1:2828}"
  -controller-url "${ASASEL_CONTROLLER_URL:?ASASEL_CONTROLLER_URL fehlt}"
  -agent-id "${ASASEL_AGENT_ID:-$(hostname)}"
  -account "${ASASEL_ACCOUNT:?ASASEL_ACCOUNT fehlt}"
)

if [[ -n "${ASASEL_SHARED_SECRET:-}" ]]; then
  args+=( -shared-secret "${ASASEL_SHARED_SECRET}" )
fi

if [[ -n "${ASASEL_AUTH_USER:-}" || -n "${ASASEL_AUTH_PASS:-}" ]]; then
  if [[ -z "${ASASEL_AUTH_USER:-}" || -z "${ASASEL_AUTH_PASS:-}" ]]; then
    echo "ASASEL_AUTH_USER und ASASEL_AUTH_PASS muessen zusammen gesetzt sein" >&2
    exit 1
  fi
  args+=( -auth-user "${ASASEL_AUTH_USER}" -auth-pass "${ASASEL_AUTH_PASS}" )
fi

exec /usr/local/bin/Asasel "${args[@]}"
EOF
sudo chmod 755 "$WRAPPER_PATH"

sudo tee "$SERVICE_FILE" >/dev/null <<EOF
[Unit]
Description=Asasel Agent Service
Documentation=https://github.com/${REPO}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
EnvironmentFile=${ENV_FILE}
ExecStart=${WRAPPER_PATH}
Restart=always
RestartSec=2
User=root

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now "$SERVICE_NAME"

echo
echo "Fertig. Status pruefen mit:"
echo "  sudo systemctl status ${SERVICE_NAME}"
echo "Logs live mit:"
echo "  sudo journalctl -u ${SERVICE_NAME} -f"
