#!/usr/bin/env bash
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATUS=0

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  STATUS=1
}

has_text() {
  grep -Fq "$1" "$2"
}

first_line() {
  sed -n '1p' "$1"
}

skill_dirs() {
  find "$ROOT" -mindepth 1 -maxdepth 1 -type d \
    ! -name '.*' \
    ! -name 'plugins' \
    ! -name 'scripts' \
    -exec basename {} \; | sort
}

frontmatter_name() {
  sed -n 's/^name:[[:space:]]*//p' "$1" | head -n 1 | sed 's/^"//; s/"$//'
}

validate_skill_dir() {
  local skill="$1"
  local dir="$ROOT/$skill"
  local skill_md="$dir/SKILL.md"
  local readme="$dir/README.md"
  local install_url="https://github.com/yarlson/skills/tree/main/$skill"

  [ -f "$skill_md" ] || fail "$skill missing SKILL.md"
  [ -f "$readme" ] || fail "$skill missing README.md"

  [ -f "$skill_md" ] || return

  local name
  name="$(frontmatter_name "$skill_md")"
  [ "$name" = "$skill" ] || fail "$skill SKILL.md name is '$name'"

  [ -f "$readme" ] || return

  [ "$(first_line "$readme")" = "# $skill" ] || fail "$skill README title mismatch"
  has_text "$install_url" "$readme" || fail "$skill README install link mismatch"
  has_text "($skill/)" "$ROOT/README.md" || fail "root README missing $skill table link"
  has_text "$install_url" "$ROOT/README.md" || fail "root README missing $skill install link"
}

validate_plugin_metadata() {
  local marketplace="$ROOT/.claude-plugin/marketplace.json"
  local plugin="$ROOT/plugins/core/.claude-plugin/plugin.json"
  local core_skills="frame plan implement review infra-review ship journey-docs repo-context-docs"

  [ -f "$marketplace" ] || fail "missing .claude-plugin/marketplace.json"
  [ -f "$plugin" ] || fail "missing plugins/core/.claude-plugin/plugin.json"

  if [ -f "$marketplace" ]; then
    has_text '"name": "yarstack"' "$marketplace" || fail "marketplace name mismatch"
    has_text '"source": "./plugins/core"' "$marketplace" || fail "marketplace core plugin source mismatch"
  fi

  if [ -f "$plugin" ]; then
    has_text '"name": "yarstack-core"' "$plugin" || fail "core plugin name mismatch"
  fi

  for skill in $core_skills; do
    local plugin_skill="$ROOT/plugins/core/skills/$skill"
    [ -e "$plugin_skill" ] || fail "core plugin missing skill path $skill"
    [ "$(frontmatter_name "$plugin_skill/SKILL.md" 2>/dev/null)" = "$skill" ] ||
      fail "core plugin skill $skill does not resolve to matching SKILL.md"
  done

  for path in "$ROOT"/plugins/core/skills/*; do
    [ -e "$path" ] || continue
    local name
    name="$(basename "$path")"
    case " $core_skills " in
      *" $name "*) ;;
      *) fail "core plugin has non-core skill path $name" ;;
    esac
  done
}

validate_readme_docs() {
  has_text 'scripts/validate-skills.sh' "$ROOT/README.md" ||
    fail "root README missing validation command"
}

for skill in $(skill_dirs); do
  validate_skill_dir "$skill"
done

validate_plugin_metadata
validate_readme_docs

if [ "$STATUS" -eq 0 ]; then
  printf 'OK: skill validation passed\n'
fi

exit "$STATUS"
