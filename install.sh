#!/usr/bin/env sh
# local-agent installer: downloads the latest release binary for the current
# OS/arch and drops it into /usr/local/bin (or ~/.local/bin if not writable).
set -eu

REPO="${LOCAL_AGENT_REPO:-hyuck0221/local-agent}"
VERSION="${LOCAL_AGENT_VERSION:-latest}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

case "$os" in
  darwin|linux) ;;
  *) echo "unsupported os: $os (use install.ps1 on Windows)" >&2; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep -E '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
fi
[ -n "$VERSION" ] || { echo "could not resolve latest version" >&2; exit 1; }

asset="local-agent_${VERSION#v}_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$VERSION/$asset"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $asset..."
curl -fsSL "$url" -o "$tmp/a.tar.gz"
tar -xzf "$tmp/a.tar.gz" -C "$tmp"

dest="/usr/local/bin"
if [ ! -w "$dest" ]; then
  dest="$HOME/.local/bin"
  mkdir -p "$dest"
fi

install -m 0755 "$tmp/local-agent" "$dest/local-agent"
echo "Installed local-agent $VERSION to $dest/local-agent"

case ":$PATH:" in
  *":$dest:"*) ;;
  *) echo "Note: add $dest to PATH to use local-agent from any shell." ;;
esac

echo
echo "Next: local-agent start"
