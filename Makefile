PYTHON ?= python3
PLUGIN_DIR := plugins/yarstack
OPENAI_CODEX_REF := 315195492c80fdade38e917c18f9584efd599304
OPENAI_VALIDATOR_URL := https://raw.githubusercontent.com/openai/codex/$(OPENAI_CODEX_REF)/codex-rs/skills/src/assets/samples/plugin-creator/scripts/validate_plugin.py

.DEFAULT_GOAL := help

.PHONY: help validate validate-metadata validate-codex validate-claude smoke-codex-install

help:
	@echo "Yarstack marketplace targets"
	@echo ""
	@echo "  make validate             Run metadata and official vendor validation"
	@echo "  make validate-metadata    Check cross-platform metadata and package integrity"
	@echo "  make validate-codex       Run OpenAI's pinned plugin validator"
	@echo "  make validate-claude      Run Claude Code strict validation"
	@echo "  make smoke-codex-install  Register and install the local Codex marketplace"

validate: validate-metadata validate-codex validate-claude

validate-metadata:
	@command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }
	@jq --exit-status --slurp '.[0].name == .[1].name and .[0].version == .[1].version and .[0].author.name == .[1].author.name and .[0].repository == .[1].repository' \
		$(PLUGIN_DIR)/.codex-plugin/plugin.json \
		$(PLUGIN_DIR)/.claude-plugin/plugin.json >/dev/null
	@jq --exit-status '.name == "yarstack" and .plugins == [{"name": "yarstack", "source": {"source": "local", "path": "./plugins/yarstack"}, "policy": {"installation": "AVAILABLE", "authentication": "ON_INSTALL"}, "category": "Productivity"}]' .agents/plugins/marketplace.json >/dev/null
	@jq --exit-status '.name == "yarstack" and (.plugins | length) == 1 and .plugins[0].name == "yarstack" and .plugins[0].source == "./plugins/yarstack" and .plugins[0].strict == true' .claude-plugin/marketplace.json >/dev/null
	@cmp LICENSE $(PLUGIN_DIR)/LICENSE
	@echo "Metadata validation passed"

validate-codex:
	@command -v curl >/dev/null || { echo "curl is required" >&2; exit 1; }
	@command -v $(PYTHON) >/dev/null || { echo "$(PYTHON) is required" >&2; exit 1; }
	@set -eu; \
	validator_dir="$$(mktemp -d)"; \
	trap 'rm -rf "$$validator_dir"' EXIT HUP INT TERM; \
	$(PYTHON) -m venv "$$validator_dir/venv"; \
	"$$validator_dir/venv/bin/pip" install --quiet --disable-pip-version-check PyYAML==6.0.3; \
	curl --fail --location --silent --show-error \
		"$(OPENAI_VALIDATOR_URL)" \
		--output "$$validator_dir/validate_plugin.py"; \
	"$$validator_dir/venv/bin/python" "$$validator_dir/validate_plugin.py" $(PLUGIN_DIR)

validate-claude:
	@command -v claude >/dev/null || { echo "claude is required" >&2; exit 1; }
	@claude plugin validate $(PLUGIN_DIR) --strict
	@claude plugin validate . --strict

smoke-codex-install:
	@command -v codex >/dev/null || { echo "codex is required" >&2; exit 1; }
	@echo "Registering this checkout as the Yarstack Codex marketplace"
	@codex plugin marketplace add "$(CURDIR)"
	@codex plugin list --marketplace yarstack
	@codex plugin add --json yarstack@yarstack
