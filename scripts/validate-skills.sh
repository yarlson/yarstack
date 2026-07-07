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
    -exec test -f '{}/SKILL.md' \; \
    -exec basename {} \; | sort
}

frontmatter_name() {
  sed -n 's/^name:[[:space:]]*//p' "$1" | head -n 1 | sed 's/^"//; s/"$//'
}

validate_frontmatter_yaml() {
  ruby -ryaml -e '
    content = File.read(ARGV[0])
    parts = content.split(/^---\s*$/, 3)
    exit 1 unless parts.length >= 3
    data = YAML.load(parts[1])
    exit 1 unless data.is_a?(Hash)
    exit 1 unless data["name"].is_a?(String)
    exit 1 unless data["description"].is_a?(String)
  ' "$1"
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
  validate_frontmatter_yaml "$skill_md" || fail "$skill SKILL.md frontmatter is invalid"

  name="$(frontmatter_name "$skill_md")"
  [ "$name" = "$skill" ] || fail "$skill SKILL.md name is '$name'"

  [ -f "$readme" ] || return

  [ "$(first_line "$readme")" = "# $skill" ] || fail "$skill README title mismatch"
  has_text "$install_url" "$readme" || fail "$skill README install link mismatch"
  has_text "($skill/)" "$ROOT/README.md" || fail "root README missing $skill table link"
  has_text "$install_url" "$ROOT/README.md" || fail "root README missing $skill install link"
}

validate_plugin_metadata() {
  local claude_marketplace="$ROOT/.claude-plugin/marketplace.json"
  local claude_plugin="$ROOT/plugins/core/.claude-plugin/plugin.json"
  local codex_marketplace="$ROOT/.agents/plugins/marketplace.json"
  local codex_plugin="$ROOT/plugins/core/.codex-plugin/plugin.json"
  local core_skills="yarstack-architecture-sparring yarstack-implementation-planning yarstack-maintainable-implementation yarstack-code-review yarstack-infra-review yarstack-draft-pr-shipping yarstack-critical-journey-docs yarstack-repo-context-docs"
  local manual_only_skills="yarstack-architecture-sparring yarstack-draft-pr-shipping yarstack-repo-context-docs"

  [ -f "$claude_marketplace" ] || fail "missing .claude-plugin/marketplace.json"
  [ -f "$claude_plugin" ] || fail "missing plugins/core/.claude-plugin/plugin.json"
  [ -f "$codex_marketplace" ] || fail "missing .agents/plugins/marketplace.json"
  [ -f "$codex_plugin" ] || fail "missing plugins/core/.codex-plugin/plugin.json"

  if [ -f "$claude_marketplace" ]; then
    has_text '"name": "yarstack"' "$claude_marketplace" || fail "Claude marketplace name mismatch"
    has_text '"source": "./plugins/core"' "$claude_marketplace" || fail "Claude marketplace core plugin source mismatch"
  fi

  if [ -f "$claude_plugin" ]; then
    has_text '"name": "yarstack-core"' "$claude_plugin" || fail "Claude core plugin name mismatch"
  fi

  if [ -f "$codex_marketplace" ]; then
    has_text '"name": "yarstack"' "$codex_marketplace" || fail "Codex marketplace name mismatch"
    has_text '"displayName": "Yarstack"' "$codex_marketplace" || fail "Codex marketplace displayName mismatch"
    has_text '"name": "yarstack-core"' "$codex_marketplace" || fail "Codex marketplace missing core plugin entry"
    has_text '"source": "local"' "$codex_marketplace" || fail "Codex marketplace core plugin source type mismatch"
    has_text '"path": "./plugins/core"' "$codex_marketplace" || fail "Codex marketplace core plugin source path mismatch"
    has_text '"installation": "AVAILABLE"' "$codex_marketplace" || fail "Codex marketplace installation policy mismatch"
    has_text '"authentication": "ON_INSTALL"' "$codex_marketplace" || fail "Codex marketplace authentication policy mismatch"
    has_text '"category": "Productivity"' "$codex_marketplace" || fail "Codex marketplace category mismatch"
  fi

  if [ -f "$codex_plugin" ]; then
    has_text '"name": "yarstack-core"' "$codex_plugin" || fail "Codex core plugin name mismatch"
    has_text '"version": "0.1.0"' "$codex_plugin" || fail "Codex core plugin version mismatch"
    has_text '"description": "Core workflow skills for disciplined coding agents."' "$codex_plugin" || fail "Codex core plugin description mismatch"
    has_text '"name": "Yarstack"' "$codex_plugin" || fail "Codex core plugin author mismatch"
    has_text '"skills": "./skills/"' "$codex_plugin" || fail "Codex core plugin skills path mismatch"
    has_text '"displayName": "Yarstack Core"' "$codex_plugin" || fail "Codex core plugin displayName mismatch"
    has_text '"shortDescription": "Core workflow skills for disciplined coding agents."' "$codex_plugin" || fail "Codex core plugin shortDescription mismatch"
    has_text '"longDescription": "Core workflow skills for framing, planning, implementing, reviewing, shipping, and documenting engineering work."' "$codex_plugin" || fail "Codex core plugin longDescription mismatch"
    has_text '"developerName": "Yarstack"' "$codex_plugin" || fail "Codex core plugin developerName mismatch"
    has_text '"capabilities": [' "$codex_plugin" || fail "Codex core plugin capabilities missing"
    has_text '"defaultPrompt": [' "$codex_plugin" || fail "Codex core plugin defaultPrompt missing"
  fi

  for skill in $core_skills; do
    local plugin_skill="$ROOT/plugins/core/skills/$skill"
    [ -e "$plugin_skill" ] || fail "core plugin missing skill path $skill"
    [ -L "$plugin_skill" ] || fail "core plugin skill $skill must be a symlink to the source skill"
    [ "$(readlink "$plugin_skill" 2>/dev/null)" = "../../../$skill" ] ||
      fail "core plugin skill $skill symlink target mismatch"
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

  for skill in $manual_only_skills; do
    local skill_md="$ROOT/$skill/SKILL.md"
    local openai_yaml="$ROOT/$skill/agents/openai.yaml"

    has_text 'disable-model-invocation: true' "$skill_md" ||
      fail "$skill missing Claude manual-only disable-model-invocation"
    [ -f "$openai_yaml" ] || fail "$skill missing Codex agents/openai.yaml"
    [ -f "$openai_yaml" ] || continue
    has_text 'display_name:' "$openai_yaml" || fail "$skill Codex openai.yaml missing display_name"
    has_text 'short_description:' "$openai_yaml" || fail "$skill Codex openai.yaml missing short_description"
    has_text 'allow_implicit_invocation: false' "$openai_yaml" ||
      fail "$skill missing Codex manual-only allow_implicit_invocation"
  done
}

validate_readme_docs() {
  has_text 'scripts/validate-skills.sh' "$ROOT/README.md" ||
    fail "root README missing validation command"
  has_text 'codex plugin marketplace add yarlson/skills' "$ROOT/README.md" ||
    fail "root README missing Codex marketplace install command"
  has_text 'codex plugin add yarstack-core@yarstack' "$ROOT/README.md" ||
    fail "root README missing Codex plugin install command"
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
