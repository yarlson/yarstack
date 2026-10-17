#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
readonly GUIDANCE_FILE="$SCRIPT_DIR/../agent-guidance/engineering-standards/standards.md"
readonly AGENT_HOME="${YARSTACK_AGENT_HOME:-${HOME:?HOME must be set}}"
readonly SCRIPT_DIR

if [[ ! -r "$GUIDANCE_FILE" ]]; then
  echo "cannot read guidance file: $GUIDANCE_FILE" >&2
  exit 1
fi

if [[ "${1:-}" == "--print" ]]; then
  cat "$GUIDANCE_FILE"
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
  if [[ -e "$destination" ]] && cmp -s "$GUIDANCE_FILE" "$destination"; then
    echo "unchanged: $destination"
    continue
  fi
  if [[ -e "$destination" ]]; then
    backup="$(mktemp "$destination.backup.$backup_timestamp.XXXXXX")"
    cp -p -- "$destination" "$backup"
    echo "backed up: $backup"
  fi
  if [[ -L "$destination" ]]; then
    cat "$GUIDANCE_FILE" > "$destination"
    chmod 0644 "$destination"
  else
    install -m 0644 "$GUIDANCE_FILE" "$destination"
  fi
  echo "installed: $destination"
done
