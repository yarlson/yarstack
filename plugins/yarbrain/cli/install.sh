#!/bin/sh
# Install the yarbrain CLI from a GitHub release on macOS or Linux, after
# verifying its SHA-256 against the release's checksums.txt.
#
#   YARBRAIN_VERSION      version to install, such as 0.3.0 (default: latest)
#   YARBRAIN_INSTALL_DIR  target directory (default: ~/.local/bin)
set -eu

repo="yarlson/yarstack"
tag_prefix="plugins/yarbrain/cli/v"
# YARBRAIN_RELEASE_BASE and YARBRAIN_RELEASE_API exist for tests only.
base="${YARBRAIN_RELEASE_BASE:-https://github.com/$repo/releases/download}"
api="${YARBRAIN_RELEASE_API:-https://api.github.com/repos/$repo/releases?per_page=50}"
install_dir="${YARBRAIN_INSTALL_DIR:-$HOME/.local/bin}"

fail() {
	printf 'yarbrain install: %s\n' "$*" >&2
	exit 1
}

command -v curl >/dev/null 2>&1 || fail "curl is required; install curl and run this script again"

case "$(uname -s)" in
Darwin) os=darwin ;;
Linux) os=linux ;;
*) fail "unsupported OS $(uname -s); on Windows, run install.ps1" ;;
esac

case "$(uname -m)" in
x86_64 | amd64) arch=amd64 ;;
arm64 | aarch64) arch=arm64 ;;
*) fail "unsupported architecture $(uname -m)" ;;
esac

# A shell under Rosetta reports x86_64 on Apple silicon; install the native build.
if [ "$os" = darwin ] && [ "$arch" = amd64 ] && [ "$(sysctl -n hw.optional.arm64 2>/dev/null || true)" = 1 ]; then
	arch=arm64
fi

version="${YARBRAIN_VERSION:-}"
if [ -z "$version" ]; then
	releases="$(curl -fsSL --max-time 30 "$api")" || fail "could not read the release list from GitHub"
	tag="$(printf '%s\n' "$releases" | grep -o "\"tag_name\": *\"${tag_prefix}[^\"]*\"" | head -n 1 | sed 's/.*"\(.*\)"$/\1/')"
	[ -n "$tag" ] || fail "no yarbrain release is published yet"
	version="${tag#"$tag_prefix"}"
fi

tag="$tag_prefix$version"
asset="yarbrain_${version}_${os}_${arch}"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT INT TERM

curl -fsSL --max-time 60 -o "$work/$asset" "$base/$tag/$asset" || fail "could not download $asset from release $tag"
curl -fsSL --max-time 60 -o "$work/checksums.txt" "$base/$tag/checksums.txt" || fail "could not download checksums.txt from release $tag"

expected="$(awk -v name="$asset" '$2 == name { print $1 }' "$work/checksums.txt")"
[ -n "$expected" ] || fail "checksums.txt of release $tag lists no $asset"

sha256() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | awk '{ print $1 }'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | awk '{ print $1 }'
	else
		fail "sha256sum or shasum is required to verify the download"
	fi
}

actual="$(sha256 "$work/$asset")"
[ "$expected" = "$actual" ] || fail "checksum mismatch for $asset; nothing was installed"

mkdir -p "$install_dir"
target="$install_dir/yarbrain"

if [ -f "$target" ] && [ "$(sha256 "$target")" = "$expected" ]; then
	printf 'yarbrain %s is already installed at %s\n' "$version" "$target"
else
	staged="$install_dir/.yarbrain.$$"
	cp "$work/$asset" "$staged"
	chmod 0755 "$staged"
	mv -f "$staged" "$target"
	printf 'Installed yarbrain %s at %s\n' "$version" "$target"
fi

# Codex runs hooks in a login shell, so the login PATH must find yarbrain.
if ! "${SHELL:-/bin/sh}" -lc 'command -v yarbrain' >/dev/null 2>&1; then
	printf 'Add %s to PATH in your login profile (for example ~/.zprofile or ~/.profile):\n' "$install_dir"
	printf '  export PATH="%s:$PATH"\n' "$install_dir"
fi

printf 'Next: ask your agent to set up Yarbrain, or run: yarbrain init <vault-path>\n'
