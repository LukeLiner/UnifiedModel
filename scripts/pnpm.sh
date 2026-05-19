#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEB_PACKAGE_JSON="${ROOT_DIR}/web/package.json"

if ! command -v node >/dev/null 2>&1; then
  echo "Node.js 22 or newer is required for the web UI, but node was not found." >&2
  exit 1
fi

NODE_MAJOR="$(node -p 'Number(process.versions.node.split(".")[0])')"
if [[ "${NODE_MAJOR}" -lt 22 ]]; then
  echo "Node.js 22 or newer is required for the web UI; current version is $(node --version)." >&2
  exit 1
fi

PNPM_VERSION="$(
  node -e "
    const fs = require('fs');
    const pkg = JSON.parse(fs.readFileSync(process.argv[1], 'utf8'));
    const manager = pkg.packageManager || 'pnpm@9.15.9';
    const match = manager.match(/^pnpm@(.+)$/);
    process.stdout.write(match ? match[1] : '9.15.9');
  " "${WEB_PACKAGE_JSON}"
)"

if command -v pnpm >/dev/null 2>&1; then
  PNPM_ACTUAL_VERSION="$(pnpm --version 2>/dev/null || true)"
  PNPM_MAJOR="${PNPM_ACTUAL_VERSION%%.*}"
  if [[ "${PNPM_MAJOR}" =~ ^[0-9]+$ && "${PNPM_MAJOR}" -ge 9 ]]; then
    exec pnpm "$@"
  fi
fi

if command -v corepack >/dev/null 2>&1; then
  exec corepack "pnpm@${PNPM_VERSION}" "$@"
fi

cat >&2 <<EOF
pnpm 9 or newer is required for the web UI.
Install pnpm, enable Corepack, or set PNPM=/path/to/pnpm when running make.
EOF
exit 1
