---
name: security-review
description: Perform a focused security review of changed code and configuration. Use when work touches authentication, authorization, secrets, command execution, file paths, archives, network calls, redirects, user input, databases, dependencies, plugins, MCP tools, logs, permissions, sandboxes, or agent boundaries.
---

# Security Review

Review the security-sensitive surface changed by the current task.

## Workflow

1. Identify touched trust boundaries, untrusted inputs, and dangerous sinks.
2. Trace external input through parsing, validation, authorization, and side effects.
3. Check for committed or exposed credentials and sensitive log data.
4. Check command construction, quoting, argument boundaries, and injection paths.
5. Check file paths, archives, uploads, and downloads for traversal, unintended overwrite, unsafe deletion, and resource exhaustion.
6. Check network behavior for SSRF, open redirects, insecure transport, credential leakage, and unbounded requests.
7. Check authentication and authorization for missing enforcement, privilege escalation, tenant leakage, insecure defaults, unsafe session or token lifecycle, and unauthenticated persistent connections such as WebSockets.
8. When cryptographic material or protected backups are touched, check key ownership, rotation, recovery copies, and exposure during backup or disaster recovery.
9. Check dependency and tool changes for unexpected sources, unsafe permissions, and inappropriate version pinning.
10. Check containers, plugins, MCP, sandboxes, and agent configuration for hardcoded secrets, broad mounts or egress, metadata-service access, token exfiltration, escape paths, and hidden side effects.
11. Fix confirmed issues within scope and report unrelated risks separately.

## Finding Standard

For each confirmed finding, identify:

- affected location and behavior;
- exploit or failure path;
- severity and impact;
- smallest safe fix;
- verification performed.

Do not report speculative concerns as confirmed vulnerabilities, rotate credentials, modify live services, add security frameworks, or broaden the task into a general security rewrite without authorization.
