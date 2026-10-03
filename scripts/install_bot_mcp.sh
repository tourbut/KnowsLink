#!/bin/sh
# Stage the KnowsLink MCP bundle and pinned Node in a durable Grok Bot folder for a Custom MCP "Command" server.
# Run from the repo root on the Bot computer. It never touches ~/.grok or app credentials and keeps KNOWSLINK_MODE held.
set -eu

PREFIX="${KNOWSLINK_PREFIX:-/workspace/.knowslink}"
VERSION=v22.22.2
case "$(uname -s)-$(uname -m)" in
Linux-x86_64) arch=x64 sum=978978a635eef872fa68beae09f0aad0bbbae6757e444da80b570964a97e62a3 ;;
Linux-aarch64 | Linux-arm64) arch=arm64 sum=b2f3a96f31486bfc365192ad65ced14833ad2a3c2e1bcefec4846902f264fa28 ;;
*)
	echo "unsupported platform: $(uname -s) $(uname -m)" >&2
	exit 1
	;;
esac

node="$PREFIX/node/bin/node"
if [ "$("$node" --version 2>/dev/null || true)" != "$VERSION" ]; then
	# .tar.gz needs only gzip; the Bot image may lack xz.
	archive="$PREFIX/node-$VERSION-linux-$arch.tar.gz"
	mkdir -p "$PREFIX"
	curl -fsSL "https://nodejs.org/dist/$VERSION/node-$VERSION-linux-$arch.tar.gz" -o "$archive"
	echo "$sum  $archive" | sha256sum -c -
	rm -rf "$PREFIX/node" && mkdir "$PREFIX/node"
	tar -xzf "$archive" -C "$PREFIX/node" --strip-components=1
	rm "$archive"
fi

PATH="$PREFIX/node/bin:$PATH"
export PATH
npm ci --prefix adapters
npm run build --prefix adapters
python3 scripts/package_plugin.py --verify

rm -rf "$PREFIX/knowslink"
python3 -m zipfile -e build/knowslink-grok-bot-plugin.zip "$PREFIX/package"
mv "$PREFIX/package/knowslink" "$PREFIX/knowslink"
rm -rf "$PREFIX/package"
# Same boundary test, run as the app would: staged node, staged bundle, no inherited PATH or env.
env -i "$node" adapters/dist/mcp.test.js "$PREFIX/knowslink/dist/plugin.js"

cat <<EOF
Ready. Register in the Grok Bot app as a Custom MCP server:
  Name:        knowslink
  Type:        Command
  Command:     $node
  Arguments:   $PREFIX/knowslink/dist/plugin.js
  Environment: none (held is the default)
EOF
