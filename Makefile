YARSTACK_PLUGIN_DIR := plugins/yarstack
YARBRAIN_PLUGIN_DIR := plugins/yarbrain
PLUGIN_SCANNER := pipx run plugin-scanner==2.0.1112

.PHONY: validate validate-yarbrain validate-codex validate-claude install-system-prompt

validate: validate-yarbrain validate-codex validate-claude

validate-yarbrain:
	python3 -m unittest discover -s $(YARBRAIN_PLUGIN_DIR)/tests -p 'test_*.py'

validate-codex:
	$(PLUGIN_SCANNER) lint $(YARSTACK_PLUGIN_DIR)
	$(PLUGIN_SCANNER) verify $(YARSTACK_PLUGIN_DIR)
	$(PLUGIN_SCANNER) lint $(YARBRAIN_PLUGIN_DIR)
	$(PLUGIN_SCANNER) verify $(YARBRAIN_PLUGIN_DIR)

validate-claude:
	claude plugin validate $(YARSTACK_PLUGIN_DIR) --strict
	claude plugin validate $(YARBRAIN_PLUGIN_DIR) --strict
	claude plugin validate . --strict

install-system-prompt:
	$(YARSTACK_PLUGIN_DIR)/scripts/install-engineering-standards.sh
