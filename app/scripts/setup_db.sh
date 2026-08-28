#!/bin/bash
# =============================================================================
# setup_db.sh — Runs postgres_setup.sql using credentials from .envrc
# =============================================================================
# USAGE:
#   chmod +x setup_db.sh
#   ./setup_db.sh
#
# REQUIRES:
#   - direnv installed and .envrc loaded (direnv allow)
#   - OR manually export variables before running:
#       export DB_NAME=myapp
#       export DB_SCHEMA=app_auth
#       export DB_APP_USER=app_user
#       export DB_APP_PASSWORD=secret
#       export PGHOST=localhost      # optional, defaults to localhost
#       export PGPORT=5432           # optional, defaults to 5432
# =============================================================================

set -euo pipefail

# ── Validate required env vars are set ────────────────────────────────────────
REQUIRED_VARS=(
  DB_NAME
  DB_APP_USER
  DB_APP_PASSWORD
)

for var in "${REQUIRED_VARS[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    echo "❌ ERROR: Required environment variable '$var' is not set."
    echo "   Make sure your .envrc is loaded (run: direnv allow)"
    exit 1
  fi
done

# ── Optional overrides with defaults ──────────────────────────────────────────
PGHOST="${PGHOST:-localhost}"
PGPORT="${PGPORT:-5432}"
PGUSER="${PGUSER:-postgres}"

echo "🚀 Setting up PostgreSQL database..."
echo "   Host     : $PGHOST:$PGPORT"
echo "   Database : $DB_NAME"
echo "   App user : $DB_APP_USER"
echo ""

# ── Run the SQL script ────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PGPASSWORD="${PGPASSWORD:-}" psql \
  -h "$PGHOST" \
  -p "$PGPORT" \
  -U "$PGUSER" \
  -v db_name="$DB_NAME" \
  -v app_user="$DB_APP_USER" \
  -v "app_password=$DB_APP_PASSWORD" \
  -f "$SCRIPT_DIR/postgres_setup.sql"

echo ""
echo "✅ Database setup complete."
