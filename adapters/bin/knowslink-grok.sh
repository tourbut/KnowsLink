#!/bin/sh
# Grok Bot plugin launcher: start the bundled MCP with the pinned Node, preparing that Node on first use.
# stdout belongs to MCP stdio, so every message goes to stderr. It never touches ~/.grok or app credentials.
# It owns only $PREFIX/node and its own $PREFIX/.stage.* folder; other files in PREFIX are kept.
# `--prepare-node` only prepares Node and prints its path (scripts/install_bot_mcp.sh reuses this).
set -eu

PREFIX="${KNOWSLINK_PREFIX:-/workspace/.knowslink}"
case "$PREFIX" in
/*) ;;
*)
	echo "KNOWSLINK_PREFIX must be an absolute path: $PREFIX" >&2
	exit 1
	;;
esac
VERSION=v22.22.2
node="$PREFIX/node/bin/node"

if [ "$("$node" --version 2>/dev/null || true)" != "$VERSION" ]; then
	case "$(uname -s)-$(uname -m)" in
	Linux-x86_64) arch=x64 sum=978978a635eef872fa68beae09f0aad0bbbae6757e444da80b570964a97e62a3 ;;
	Linux-aarch64 | Linux-arm64) arch=arm64 sum=b2f3a96f31486bfc365192ad65ced14833ad2a3c2e1bcefec4846902f264fa28 ;;
	*)
		echo "unsupported platform: $(uname -s) $(uname -m)" >&2
		exit 1
		;;
	esac
	# Replace only what an earlier run made; anything else at this path is the user's, so stop before changing it.
	if [ -L "$PREFIX/node" ] || { [ -e "$PREFIX/node" ] && [ ! -f "$node" ]; }; then
		echo "refusing to replace $PREFIX/node: not made by this installer (node/bin/node missing); move it away first" >&2
		exit 1
	fi
	mkdir -p "$PREFIX"
	stage="$(mktemp -d "$PREFIX/.stage.XXXXXX")"
	trap 'rm -rf "$stage"' EXIT
	trap 'exit 1' HUP INT TERM
	echo "KnowsLink: preparing Node $VERSION in $PREFIX/node" >&2
	# .tar.gz needs only gzip; the Bot image may lack xz.
	archive="$stage/node.tar.gz"
	curl -fsSL "https://nodejs.org/dist/$VERSION/node-$VERSION-linux-$arch.tar.gz" -o "$archive" >&2
	if command -v sha256sum >/dev/null; then got="$(sha256sum "$archive")"; else got="$(shasum -a 256 "$archive")"; fi
	if [ "${got%% *}" != "$sum" ]; then
		echo "Node checksum mismatch; the existing $PREFIX/node is kept" >&2
		exit 1
	fi
	mkdir "$stage/node"
	tar -xzf "$archive" -C "$stage/node" --strip-components=1
	[ "$("$stage/node/bin/node" --version)" = "$VERSION" ]
	# shortcut: two first starts at once can both download; the loser re-checks and keeps the winner's Node. Add a lock if Bots start many MCPs together.
	if [ "$("$node" --version 2>/dev/null || true)" != "$VERSION" ]; then
		# The old tree leaves only after the verified one is ready; it is removed with the stage.
		if [ -d "$PREFIX/node" ]; then mv "$PREFIX/node" "$stage/old-node"; fi
		mv "$stage/node" "$PREFIX/node"
	fi
	# exec below skips the EXIT trap, so clean up here.
	rm -rf "$stage"
	trap - EXIT
fi

if [ "${1:-}" = --prepare-node ]; then
	echo "$node"
	exit 0
fi

# Release defaults (T-03): a variable set in the app wins. No request is made before knowslink_connect with confirmed:true.
: "${KNOWSLINK_MODE:=public-node}" "${RELAY_URL:=https://link.knowslog.com}" "${KNOWSLINK_HOST:=grok}" "${KNOWSLINK_AGENT_ROOT:=$PREFIX/connections}"
export KNOWSLINK_MODE RELAY_URL KNOWSLINK_HOST KNOWSLINK_AGENT_ROOT
exec "$node" "$(dirname "$0")/../dist/plugin.js" "$@"
