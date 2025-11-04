#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "${SCRIPT_DIR}/.." && pwd)
cd "${REPO_ROOT}"

DB_FILE="var/duplynx.db"
ASSETS_DIR="backend/web/dist"
ADDR="127.0.0.1:8080"
OUTPUT="var/duplynx_bench.json"
BIN_DIR="${REPO_ROOT}/var/bin"
CLI_BIN="${BIN_DIR}/duplynx-bench"
SERVER_PID=""
ADDR_OVERRIDDEN=0
PYTHON_BIN="${PYTHON_BIN:-python3}"

if ! command -v "$PYTHON_BIN" >/dev/null 2>&1; then
  PYTHON_BIN="python"
fi

if ! command -v "$PYTHON_BIN" >/dev/null 2>&1; then
  echo "Required Python interpreter not found (looked for python3/python)." >&2
  exit 1
fi

usage() {
  cat <<USAGE
Usage: $(basename "$0") [options]

Options:
  --db-file PATH        Path to SQLite database file (default: ${DB_FILE})
  --assets-dir PATH     Path to Tailwind asset directory (default: ${ASSETS_DIR})
  --addr HOST:PORT      Address for temporary serve check (default: ${ADDR})
  --output PATH         Path for benchmark JSON output (default: ${OUTPUT})
  -h, --help            Show this help message
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --db-file)
      DB_FILE="$2"
      shift 2
      ;;
    --assets-dir)
      ASSETS_DIR="$2"
      shift 2
      ;;
    --addr)
      ADDR="$2"
      ADDR_OVERRIDDEN=1
      shift 2
      ;;
    --output)
      OUTPUT="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

case "$DB_FILE" in
  /*) ;;
  *) DB_FILE="$REPO_ROOT/$DB_FILE" ;;
esac

case "$ASSETS_DIR" in
  /*) ;;
  *) ASSETS_DIR="$REPO_ROOT/$ASSETS_DIR" ;;
esac

case "$OUTPUT" in
  /*) ;;
  *) OUTPUT="$REPO_ROOT/$OUTPUT" ;;
esac

mkdir -p "$(dirname "$DB_FILE")"
mkdir -p "$BIN_DIR"

if [[ ! -f "$ASSETS_DIR/tailwind.css" ]]; then
  echo "Tailwind bundle missing at $ASSETS_DIR/tailwind.css" >&2
  echo "Run 'npm run build:tailwind' before measuring quickstart." >&2
  exit 1
fi

# Ensure any background duplynx serve instance is cleaned up on exit or signal.
cleanup() {
  if [[ -z "${SERVER_PID:-}" ]]; then
    return
  fi

  if kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
    for _ in 1 2 3 4 5 6 7 8 9 10; do
      if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
        break
      fi
      sleep 0.2
    done
    if kill -0 "$SERVER_PID" >/dev/null 2>&1; then
      kill -9 "$SERVER_PID" >/dev/null 2>&1 || true
    fi
  fi

  wait "$SERVER_PID" >/dev/null 2>&1 || true
  SERVER_PID=""
}

handle_signal_exit() {
  local status="${1:-130}"
  cleanup
  exit "$status"
}

trap cleanup EXIT
trap 'handle_signal_exit 130' INT
trap 'handle_signal_exit 143' TERM

# Build the CLI once so follow-up commands use a direct binary (no lingering go run wrapper).
build_cli() {
  (cd backend && go build -o "$CLI_BIN" ./cmd/duplynx)
}

now_ns() {
  date +%s%N
}

exec_seed() {
  "$CLI_BIN" seed --db-file "$DB_FILE" --assets-dir "$ASSETS_DIR"
}

start_serve() {
  "$CLI_BIN" serve --db-file "$DB_FILE" --assets-dir "$ASSETS_DIR" --addr "$ADDR" &
  SERVER_PID=$!
}

# Check whether the chosen address is available; if the default port is busy, pick a free one.
ensure_addr_available() {
  local addr="$1"
  local host="${addr%:*}"
  local port="${addr##*:}"

  if [[ "$port" == "$addr" ]]; then
    echo "Invalid address format: $addr (expected host:port)" >&2
    exit 1
  fi

  if [[ "$port" == "0" ]]; then
    return
  fi

  "$PYTHON_BIN" - "$host" "$port" <<'PY'
import socket, sys
host, port = sys.argv[1], int(sys.argv[2])
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    try:
        sock.bind((host, port))
    except OSError:
        sys.exit(1)
sys.exit(0)
PY
}

pick_free_port() {
  "$PYTHON_BIN" - <<'PY'
import socket
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.bind(("127.0.0.1", 0))
    print(sock.getsockname()[1])
PY
}

build_cli
if ! ensure_addr_available "$ADDR"; then
  if (( ADDR_OVERRIDDEN == 1 )); then
    echo "Address $ADDR is already in use; quickstart cannot proceed" >&2
    exit 1
  fi
  free_port=$(pick_free_port)
  ADDR="127.0.0.1:${free_port}"
  echo "Address 127.0.0.1:8080 busy, retrying on ${ADDR}" >&2
fi
seed_start=$(now_ns)
exec_seed
seed_end=$(now_ns)
seed_ms=$(( (seed_end - seed_start) / 1000000 ))

serve_start=$(now_ns)
start_serve
if [[ -z "${SERVER_PID:-}" ]]; then
  echo "Failed to start serve process" >&2
  cleanup
  exit 1
fi
server_started=0
declare -i serve_timeout_deadline=$((SECONDS + 120))

declare -i wait_deadline=$((SECONDS + 60))
until (( SECONDS > serve_timeout_deadline )); do
  if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    break
  fi
  if curl -fsS "http://$ADDR/healthz" >/dev/null 2>&1; then
    server_started=1
    break
  fi
  if (( SECONDS > wait_deadline )); then
    break
  fi
  sleep 1
done

if (( server_started == 0 )); then
  if kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    echo "Server did not become healthy within 60 seconds; killing serve process" >&2
  else
    echo "Serve process exited before becoming healthy" >&2
  fi
  cleanup
  exit 1
fi

root_html=$(curl -fsS "http://$ADDR/")
if [[ "$root_html" != *"Orion Analytics"* ]]; then
  echo "Dashboard response missing expected tenant marker" >&2
  cleanup
  exit 1
fi

if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
  echo "Serve process exited unexpectedly after responding; quickstart failed" >&2
  cleanup
  exit 1
fi

serve_end=$(now_ns)
serve_duration_ms=$(( (serve_end - serve_start) / 1000000 ))
cleanup

timestamp=$(date -Iseconds)
mkdir -p "$(dirname "$OUTPUT")"
cat <<JSON >"$OUTPUT"
{
  "timestamp": "${timestamp}",
  "db_file": "${DB_FILE}",
  "assets_dir": "${ASSETS_DIR}",
  "addr": "${ADDR}",
  "seed_ms": ${seed_ms},
  "serve_ms": ${serve_duration_ms},
  "total_ms": $(( seed_ms + serve_duration_ms ))
}
JSON

echo "Seed completed in ${seed_ms} ms"
echo "Serve verified in ${serve_duration_ms} ms"
echo "Benchmark written to ${OUTPUT}"
