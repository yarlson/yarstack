# Security

Report vulnerabilities privately. Do not open a public issue or include
credentials or sensitive data:

https://github.com/yarlson/yarstack/security/advisories/new

Yarstack can write global agent guidance when the user runs its explicit
installer. Yarbrain writes only to the vault that the user explicitly
configures, and its lifecycle hooks store session locators rather than
transcript contents.

Changes that add executable code, network access, authentication, sensitive
data capture, or new write targets must make that behavior explicit.
