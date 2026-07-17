PLUGIN_DIR := plugins/yarstack
PLUGIN_SCANNER := pipx run plugin-scanner==2.0.1112

.PHONY: validate validate-codex validate-claude install-system-prompt

validate: validate-codex validate-claude

validate-codex:
	$(PLUGIN_SCANNER) lint $(PLUGIN_DIR)
	$(PLUGIN_SCANNER) verify $(PLUGIN_DIR)

validate-claude:
	claude plugin validate $(PLUGIN_DIR) --strict
	claude plugin validate . --strict

install-system-prompt:
	$(PLUGIN_DIR)/scripts/install-engineering-standards.sh
