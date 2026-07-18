#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
readonly GUIDANCE_DIR="$SCRIPT_DIR/../agent-guidance/engineering-standards"
readonly AGENT_HOME="${YARSTACK_AGENT_HOME:-${HOME:?HOME must be set}}"
COMBINED_FILE="$(mktemp "${TMPDIR:-/tmp}/yarstack-engineering-standards.XXXXXX")"
readonly SCRIPT_DIR COMBINED_FILE

cleanup() {
  rm -f -- "$COMBINED_FILE"
}
trap cleanup EXIT

{
  cat "$GUIDANCE_DIR/engineering-quality-gate.md"
  printf '\n'
  cat "$GUIDANCE_DIR/existing-codebase-first.md"
  printf '\n'
  cat "$GUIDANCE_DIR/small-design-before-code.md"
  printf '\n'
  cat "$GUIDANCE_DIR/new-code-complexity-budget.md"
  printf '\n'
  cat "$GUIDANCE_DIR/tests-as-product-contracts.md"
  printf '\n'
  cat "$GUIDANCE_DIR/plain-english-writing.md"
  printf '\n'
  cat "$GUIDANCE_DIR/code-comment-policy.md"
  printf '\n'
  cat "$GUIDANCE_DIR/implementation-change-report.md"
} > "$COMBINED_FILE"

if [[ "${1:-}" == "--print" ]]; then
  cat "$COMBINED_FILE"
  exit 0
fi
if [[ $# -ne 0 ]]; then
  echo "usage: $(basename "$0") [--print]" >&2
  exit 2
fi

destinations=(
  "$AGENT_HOME/.claude/CLAUDE.md"
  "$AGENT_HOME/.agents/AGENTS.md"
  "$AGENT_HOME/.codex/AGENTS.md"
)
backup_timestamp="$(date -u +%Y%m%dT%H%M%SZ)"

for destination in "${destinations[@]}"; do
  if [[ -L "$destination" && ! -e "$destination" ]]; then
    echo "refusing to write through dangling symlink: $destination" >&2
    exit 1
  fi
done

for destination in "${destinations[@]}"; do
  mkdir -p -- "$(dirname -- "$destination")"
  if [[ -e "$destination" ]] && cmp -s "$COMBINED_FILE" "$destination"; then
    echo "unchanged: $destination"
    continue
  fi
  if [[ -e "$destination" ]]; then
    backup="$(mktemp "$destination.backup.$backup_timestamp.XXXXXX")"
    cp -p -- "$destination" "$backup"
    echo "backed up: $backup"
  fi
  if [[ -L "$destination" ]]; then
    cat "$COMBINED_FILE" > "$destination"
    chmod 0644 "$destination"
  else
    install -m 0644 "$COMBINED_FILE" "$destination"
  fi
  echo "installed: $destination"
done
