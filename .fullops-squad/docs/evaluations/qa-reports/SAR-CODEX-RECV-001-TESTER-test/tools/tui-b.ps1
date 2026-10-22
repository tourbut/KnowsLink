# SAR-CODEX-RECV-001-TESTER: independent Codex TUI for agent B with the prepared install's MCP settings as per-invocation overrides (no ~/.codex change).
$i = 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/codex-recv-001-tester'
Set-Location "$i/work-tui"
codex -c "mcp_servers.knowslink.command='C:/Users/shin/AppData/Local/KnowsLinkDevTools/node-v22.22.2-win-x64/node.EXE'" -c "mcp_servers.knowslink.args=['$i/inst-b/knowslink/dist/plugin.js']" -c "mcp_servers.knowslink.env={KNOWSLINK_MODE='public-node', RELAY_URL='http://127.0.0.1:18321', KNOWSLINK_AGENT_FOLDER='$i/qa/m-b-private', KNOWSLINK_SESSION_HOST='codex'}" -c "model_reasoning_effort='low'" @args
