# Security policy

## Supported versions

Only the latest published Yarstack version receives security fixes.

| Version        | Supported |
| -------------- | --------- |
| 0.1.x          | Yes       |
| Older versions | No        |

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use GitHub's private
vulnerability reporting for `yarlson/yarstack`:

https://github.com/yarlson/yarstack/security/advisories/new

Include the affected version, platform, impact, reproduction steps, and any
suggested mitigation. Avoid including real credentials or sensitive user data.
If private vulnerability reporting is unavailable, contact the maintainer
through https://github.com/yarlson before sharing technical details publicly.

## Security expectations

Plugin changes must make executable code, network access, external services,
write actions, authentication, and persistent state explicit. Secrets must not
be stored in manifests or logs. Dependencies and downloaded artifacts must be
pinned or integrity-checked where practical.
