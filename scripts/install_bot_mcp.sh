#!/bin/sh
# Stage the KnowsLink MCP bundle and pinned Node in a durable Grok Bot folder for a Custom MCP "Command" server.
# Run from the repo root on the Bot computer. It never touches ~/.grok or app credentials; its release defaults start no connection.
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
# Node pin, its checksum and the node/ ownership rule live in the plugin launcher; this script adds the bundle.
# Replace only what an earlier run made; anything else at this path is the user's, so stop before changing it.
if [ -L "$PREFIX/knowslink" ] || { [ -e "$PREFIX/knowslink" ] && [ ! -f "$PREFIX/knowslink/dist/plugin.js" ]; }; then
	echo "refusing to replace $PREFIX/knowslink: not made by this installer (knowslink/dist/plugin.js missing); move it away first" >&2
	exit 1
fi

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

node="$(KNOWSLINK_PREFIX="$PREFIX" sh adapters/bin/knowslink-grok.sh --prepare-node)"
PATH="$PREFIX/node/bin:$PATH"
export PATH
npm ci --prefix adapters
npm run build --prefix adapters
python3 scripts/package_plugin.py --host grok --verify

python3 -m zipfile -e build/knowslink-grok-bot-plugin.zip "$stage"
# Same boundary test, run as the app would: staged node, new bundle, no inherited PATH or env.
env -i "$node" adapters/dist/mcp.test.js "$stage/knowslink/dist/plugin.js"
# Release defaults (T-03): with no environment the Command server is public-node on the service with an automatic agent folder.
# Variables set in the app still win. No request is made before knowslink_connect with confirmed:true.
python3 -c 'import json, sys; json.dump({"KNOWSLINK_MODE": "public-node", "RELAY_URL": "https://link.knowslog.com", "KNOWSLINK_HOST": "grok", "KNOWSLINK_AGENT_ROOT": sys.argv[1] + "/connections"}, open(sys.argv[2], "w"))' "$PREFIX" "$stage/knowslink/defaults.json"
swap knowslink

cat <<EOF
Ready. Register in the Grok Bot app as a Custom MCP server:
  Name:        knowslink
  Type:        Command
  Command:     $node
  Arguments:   $PREFIX/knowslink/dist/plugin.js
  Environment: none

Defaults ($PREFIX/knowslink/defaults.json): KNOWSLINK_MODE=public-node,
RELAY_URL=https://link.knowslog.com and the agent folder
$PREFIX/connections/grok-agent. A variable set in the app wins.
Nothing contacts the service until the user agrees: call knowslink_connect with
confirmed:true and have the user open its URL, compare the fingerprint, and
approve Google connection. Check knowslink_connect_status for connected; a
failure shows its reason and next action. A retry reuses or sets the folder
aside by itself; never delete it. Grok Bots on this computer that use these
defaults share one agent (knowslink_status shows it). Person friendship
consent is separate.
Default receive remains polling. KNOWSLINK_RECEIVE_MODE=sse opts into arrival
signals plus existing API pull and periodic recovery. Neither starts a host turn.
EOF
