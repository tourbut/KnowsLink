#!/bin/sh
# Stage the KnowsLink MCP bundle and pinned Node in a durable Grok Bot folder for a Custom MCP "Command" server.
# Run from the repo root on the Bot computer. It never touches ~/.grok or app credentials and keeps KNOWSLINK_MODE held.
# It owns only $PREFIX/node, $PREFIX/knowslink and its own $PREFIX/.stage.* folder; other files in PREFIX are kept.
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
case "$(uname -s)-$(uname -m)" in
Linux-x86_64) arch=x64 sum=978978a635eef872fa68beae09f0aad0bbbae6757e444da80b570964a97e62a3 ;;
Linux-aarch64 | Linux-arm64) arch=arm64 sum=b2f3a96f31486bfc365192ad65ced14833ad2a3c2e1bcefec4846902f264fa28 ;;
*)
	echo "unsupported platform: $(uname -s) $(uname -m)" >&2
	exit 1
	;;
esac

# Replace only what an earlier run made; anything else at these paths is the user's, so stop before changing it.
for owned in node/bin/node knowslink/dist/plugin.js; do
	dir="$PREFIX/${owned%%/*}"
	if [ -L "$dir" ] || { [ -e "$dir" ] && [ ! -f "$PREFIX/$owned" ]; }; then
		echo "refusing to replace $dir: not made by this installer ($owned missing); move it away first" >&2
		exit 1
	fi
done

mkdir -p "$PREFIX"
# A fresh area per run: leftovers of an interrupted run never mix into the new bundle.
stage="$(mktemp -d "$PREFIX/.stage.XXXXXX")"
trap 'rm -rf "$stage"' EXIT
trap 'exit 1' HUP INT TERM

# Swap a verified new tree in; the old one is removed with the stage only after the swap.
swap() {
	if [ -d "$PREFIX/$1" ]; then mv "$PREFIX/$1" "$stage/old-$1"; fi
	mv "$stage/$1" "$PREFIX/$1"
}

node="$PREFIX/node/bin/node"
if [ "$("$node" --version 2>/dev/null || true)" != "$VERSION" ]; then
	# .tar.gz needs only gzip; the Bot image may lack xz.
	archive="$stage/node-$VERSION-linux-$arch.tar.gz"
	curl -fsSL "https://nodejs.org/dist/$VERSION/node-$VERSION-linux-$arch.tar.gz" -o "$archive"
	echo "$sum  $archive" | sha256sum -c -
	mkdir "$stage/node"
	tar -xzf "$archive" -C "$stage/node" --strip-components=1
	swap node
fi

PATH="$PREFIX/node/bin:$PATH"
export PATH
npm ci --prefix adapters
npm run build --prefix adapters
python3 scripts/package_plugin.py --verify

python3 -m zipfile -e build/knowslink-grok-bot-plugin.zip "$stage"
# Same boundary test, run as the app would: staged node, new bundle, no inherited PATH or env.
env -i "$node" adapters/dist/mcp.test.js "$stage/knowslink/dist/plugin.js"
swap knowslink

cat <<EOF
Ready. Register in the Grok Bot app as a Custom MCP server:
  Name:        knowslink
  Type:        Command
  Command:     $node
  Arguments:   $PREFIX/knowslink/dist/plugin.js
  Environment: none (held is the default)
EOF
