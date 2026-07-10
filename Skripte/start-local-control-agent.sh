#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BACKEND_DIR="${REPO_ROOT}/backend"

BIN="${ASASEL_BIN:-${BACKEND_DIR}/Asasel}"
CONTROLLER_HOST="${CONTROLLER_HOST:-127.0.0.1}"
CONTROLLER_PORT="${CONTROLLER_PORT:-2727}"
AGENT_HOST="${AGENT_HOST:-127.0.0.1}"
AGENT_PORT="${AGENT_PORT:-2828}"
AGENT_ID="${AGENT_ID:-local-agent}"
ACCOUNT="${ACCOUNT:-linus}"
AUTH_USER="${AUTH_USER:-test}"
AUTH_PASS="${AUTH_PASS:-test}"
SHARED_SECRET="${SHARED_SECRET:-local-dev-secret}"

LOG_DIR="${LOG_DIR:-/tmp/asasel-local-test}"
CONTROLLER_LOG="${LOG_DIR}/controller.log"
AGENT_LOG="${LOG_DIR}/agent.log"

mkdir -p "${LOG_DIR}"

if [[ ! -x "${BIN}" ]]; then
  echo "Binary not found or not executable: ${BIN}" >&2
  echo "Trying to build backend binary..." >&2
  (
    cd "${BACKEND_DIR}"
    if command -v go >/dev/null 2>&1; then
      go build .
    elif [[ -x "${HOME}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.linux-amd64/bin/go" ]]; then
      "${HOME}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.1.linux-amd64/bin/go" build .
    else
      echo "No Go toolchain found. Please build backend manually first." >&2
      exit 1
    fi
  )
fi

if [[ ! -x "${BIN}" ]]; then
  echo "Binary still missing after build: ${BIN}" >&2
  exit 1
fi

cleanup() {
  if [[ -n "${AGENT_PID:-}" ]] && kill -0 "${AGENT_PID}" 2>/dev/null; then
    kill "${AGENT_PID}" 2>/dev/null || true
  fi
  if [[ -n "${CONTROLLER_PID:-}" ]] && kill -0 "${CONTROLLER_PID}" 2>/dev/null; then
    kill "${CONTROLLER_PID}" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

echo "Starting controller on ${CONTROLLER_HOST}:${CONTROLLER_PORT} ..."
"${BIN}" \
  -mode control \
  -listen "${CONTROLLER_HOST}:${CONTROLLER_PORT}" \
  -account "${ACCOUNT}" \
  -auth-user "${AUTH_USER}" \
  -auth-pass "${AUTH_PASS}" \
  -shared-secret "${SHARED_SECRET}" \
  >"${CONTROLLER_LOG}" 2>&1 &
CONTROLLER_PID=$!

echo "Starting agent on ${AGENT_HOST}:${AGENT_PORT} ..."
"${BIN}" \
  -mode local \
  -listen "${AGENT_HOST}:${AGENT_PORT}" \
  -agent-id "${AGENT_ID}" \
  -controller-url "http://${CONTROLLER_HOST}:${CONTROLLER_PORT}" \
  -account "${ACCOUNT}" \
  -auth-user "${AUTH_USER}" \
  -auth-pass "${AUTH_PASS}" \
  -shared-secret "${SHARED_SECRET}" \
  >"${AGENT_LOG}" 2>&1 &
AGENT_PID=$!

sleep 1

if ! kill -0 "${CONTROLLER_PID}" 2>/dev/null; then
  echo "Controller exited unexpectedly. See: ${CONTROLLER_LOG}" >&2
  exit 1
fi

if ! kill -0 "${AGENT_PID}" 2>/dev/null; then
  echo "Agent exited unexpectedly. See: ${AGENT_LOG}" >&2
  exit 1
fi

echo
echo "Controller and Agent are running."
echo "Controller URL: http://${CONTROLLER_HOST}:${CONTROLLER_PORT}/"
echo "Agent ID: ${AGENT_ID}"
echo "Auth: ${AUTH_USER} / ${AUTH_PASS}"
echo "Logs:"
echo "  - ${CONTROLLER_LOG}"
echo "  - ${AGENT_LOG}"
echo
echo "Press Ctrl+C to stop both processes."

wait "${CONTROLLER_PID}" "${AGENT_PID}"
