#!/bin/sh
# SAR-GROK-INSTALL-001-DEV QA in a throwaway Linux container. /plugin = release tree (ro), /r = repo worktree (ro).
# Prints PASS/FAIL lines; exit 1 on the first failure.
set -u
L=/plugin/knowslink/bin/knowslink-grok.sh
P=/workspace/.knowslink
export PATH=/usr/bin:/bin # hide the image's /usr/local/bin/node: only the launcher's pinned Node may run
ok() { echo "PASS $*"; }
no() { echo "FAIL $*"; exit 1; }
init='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"qa","version":"0"}}}'
call() { # MCP over stdio through the launcher, no environment set (as the app runs it)
	{ echo "$init"; echo '{"jsonrpc":"2.0","method":"notifications/initialized"}'; echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"knowslink_status","arguments":{}}}'; sleep 3; } |
		env -i PATH=/usr/bin:/bin sh "$L" 2>"$1"
}
fake_old() { mkdir -p "$P/node/bin"; printf '#!/bin/sh\necho v22.0.0\n' >"$P/node/bin/node"; chmod +x "$P/node/bin/node"; }
same_fake() { [ "$("$P/node/bin/node" --version)" = v22.0.0 ] && [ -z "$(ls -A "$P" | grep '^\.stage' || true)" ]; }

echo "arch: $(uname -m)"
command -v node && no "node visible on PATH"

# 1. Fresh install: first MCP start prepares Node, answers MCP, release defaults active.
out="$(call /tmp/e1)"
echo "$out" | grep -q '"serverInfo"' || { cat /tmp/e1; no "1 fresh: no initialize result"; }
echo "$out" | grep -q 'public_node_configured_unverified' || { echo "$out"; no "1 fresh: defaults not public-node"; }
grep -q "preparing Node v22.22.2" /tmp/e1 || no "1 fresh: Node not prepared"
[ "$(env -i "$P/node/bin/node" --version)" = v22.22.2 ] || no "1 fresh: node version"
[ -f /plugin/knowslink/dist/plugin.js ] || no "1 fresh: plugin.js"
[ -z "$(ls -A "$P" | grep '^\.stage' || true)" ] || no "1 fresh: stage left"
(cd /r && env -i "$P/node/bin/node" adapters/dist/mcp.test.js /plugin/knowslink/dist/plugin.js) >/tmp/b 2>&1 || { tail -5 /tmp/b; no "1 fresh: env -i boundary test"; }
ok "1 fresh install: env -i node --version=v22.22.2, plugin.js present, initialize ok, knowslink_status public_node_configured_unverified, env -i MCP boundary test exit 0, no stage left"

# 2. Same command again reuses Node (no download); an older Node made by us is upgraded in place.
before="$(stat -c %Y "$P/node/bin/node")"
call /tmp/e2 | grep -q '"serverInfo"' || no "2 rerun: no initialize"
grep -q "preparing Node" /tmp/e2 && no "2 rerun: downloaded again"
[ "$(stat -c %Y "$P/node/bin/node")" = "$before" ] || no "2 rerun: node changed"
rm -rf "$P/node" && fake_old
call /tmp/e3 | grep -q '"serverInfo"' || no "2 upgrade: no initialize"
[ "$(env -i "$P/node/bin/node" --version)" = v22.22.2 ] || no "2 upgrade: not replaced"
ok "2 rerun reuses Node v22.22.2 without download; old v22.0.0 made by the launcher is replaced by verified v22.22.2"

# 3a. Checksum mismatch keeps the existing install.
rm -rf "$P/node" && fake_old
sed 's/sum=978978a6/sum=000000a6/; s/sum=b2f3a96f/sum=000000ff/' "$L" >/tmp/bad.sh
env -i PATH=/usr/bin:/bin sh /tmp/bad.sh --prepare-node </dev/null >/tmp/o 2>/tmp/e4 && no "3a mismatch: exit 0"
grep -q "checksum mismatch" /tmp/e4 || no "3a mismatch: message"
same_fake || no "3a mismatch: existing node changed or stage left"
ok "3a checksum mismatch: exit $(env -i PATH=/usr/bin:/bin sh /tmp/bad.sh --prepare-node </dev/null >/dev/null 2>&1; echo $?), '$(cat /tmp/e4 | tail -1)', existing node kept, no stage left"

# 3c. A node/ folder the launcher did not make is refused untouched.
rm -rf "$P/node" && mkdir -p "$P/node" && echo mine >"$P/node/keep"
env -i PATH=/usr/bin:/bin sh "$L" --prepare-node </dev/null >/dev/null 2>/tmp/e5 && no "3c foreign: exit 0"
[ "$(cat "$P/node/keep")" = mine ] || no "3c foreign: changed"
ok "3c foreign node/: refused ('$(tail -1 /tmp/e5)'), kept"
rm -rf "$P/node" && fake_old
echo "READY-FOR-OFFLINE"
