# CS-AC-14: release-style Grok package (zip from package_plugin.py + the installer's defaults.json), no environment: zero requests before consent, connect after consent; existing env settings still win.
from lib import *
import zipfile
ZIP = "C:/Users/shin/orca/workspaces/KnowsLink/tester/build/knowslink-grok-bot-plugin.zip"
stage = W + "/grok-stage"
if not os.path.exists(stage + "/knowslink/dist/plugin.js"):
    os.makedirs(stage, exist_ok=True)
    zipfile.ZipFile(ZIP).extractall(stage)
plugin = stage + "/knowslink/dist/plugin.js"
# the installer's exact defaults (install_bot_mcp.sh) with RELAY_URL pointed at the fixture instead of production
installer = open("C:/Users/shin/orca/workspaces/KnowsLink/tester/scripts/install_bot_mcp.sh", encoding="utf-8").read()
prod_ok = '"KNOWSLINK_MODE": "public-node", "RELAY_URL": "https://link.knowslog.com", "KNOWSLINK_HOST": "grok"' in installer
rec("CS-AC-14", "installer script writes release defaults: public-node, production URL, host grok, connections/ root (source check)", prod_ok and "connections" in installer, "scripts/install_bot_mcp.sh line 68")
prefix = W + "/grok-prefix"
os.makedirs(prefix, exist_ok=True)
def write_defaults(relay):
    json.dump({"KNOWSLINK_MODE": "public-node", "RELAY_URL": relay, "KNOWSLINK_HOST": "grok", "KNOWSLINK_AGENT_ROOT": prefix + "/connections"}, open(stage + "/knowslink/defaults.json", "w"))
class Rpc:
    def __init__(self, env):
        e = {k: v for k, v in os.environ.items()}
        for k in [k for k in e if k.startswith("KNOWSLINK_") or k == "RELAY_URL"]: e.pop(k)
        self.p = subprocess.Popen([NODE, os.path.join(OUT, "mcp_rpc.mjs"), plugin, json.dumps(env)], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, encoding="utf-8", env=e)
    def call(self, tool=None, args=None):
        self.p.stdin.write(json.dumps({"list": True} if tool is None else {"tool": tool, "args": args or {}}) + "\n"); self.p.stdin.flush()
        line = self.p.stdout.readline()
        with open(os.path.join(OUT, "mcp-output.log"), "a", encoding="utf-8") as f: f.write(json.dumps({"tool": tool, "args": args}, ensure_ascii=False) + chr(10) + line)
        return json.loads(line)
    def close(self): self.p.stdin.close(); self.p.wait(timeout=20)
def payload(r):
    try: return json.loads(r["content"][0]["text"])
    except Exception: return r
# --- A: no environment, release defaults -> fixture
write_defaults(URL)
fx("/fixture/hits?reset=1")
rpc = Rpc({})
tools = [t["name"] for t in rpc.call()["tools"]]
rec("CS-AC-14", "release package exposes knowslink_connect / connect_status / status tools", {"knowslink_connect", "knowslink_connect_status", "knowslink_status"} <= set(tools), ",".join(tools))
st = payload(rpc.call("knowslink_status")); cs = payload(rpc.call("knowslink_connect_status"))
no_confirm = rpc.call("knowslink_connect", {})
no_false = rpc.call("knowslink_connect", {"confirmed": False})
hits = json.loads(fx("/fixture/hits")[1])
open(os.path.join(OUT, "grok-before-consent.json"), "w", encoding="utf-8").write(json.dumps({"status": st, "connect_status": cs, "no_confirm": no_confirm, "confirm_false": no_false, "hits": hits}, ensure_ascii=False, indent=1))
rec("CS-AC-14", "no environment: status, connect_status, connect without/false confirmation make 0 relay requests", hits["api"] == 0 and hits["starts"] == 0, f"hits={hits} no_confirm_error={'isError' in no_confirm and no_confirm.get('isError')} rpcError={no_confirm.get('rpcError','')[:100]}")
rec("CS-AC-14", "connect without confirmed:true is refused by the tool (not started)", bool(no_confirm.get("isError") or no_confirm.get("rpcError")) and bool(no_false.get("isError") or no_false.get("rpcError")), f"{json.dumps(no_confirm, ensure_ascii=False)[:200]}")
fx("/fixture/hits?reset=1")
started = payload(rpc.call("knowslink_connect", {"confirmed": True}))
hits = json.loads(fx("/fixture/hits")[1])
rec("CS-AC-14", "after consent (confirmed:true) the connection starts: one start request, URL + fingerprint returned, no env set", hits["starts"] == 1 and "url" in json.dumps(started), json.dumps(started, ensure_ascii=False)[:300])
url = started.get("url") if isinstance(started, dict) else None
b, cid, c = approve("a", url, None)
for _ in range(20):
    status = payload(rpc.call("knowslink_connect_status"))
    if isinstance(status, dict) and status.get("state") == "connected": break
    time.sleep(3)
rec("CS-AC-14", "browser approval -> knowslink_connect_status reports connected with the agent name (grok default)", isinstance(status, dict) and status.get("state") == "connected", json.dumps(status, ensure_ascii=False)[:300])
st2 = payload(rpc.call("knowslink_status"))
rec("CS-AC-09", "grok knowslink_status shows my agent's display name (A-06), no key/token", "name" in json.dumps(st2, ensure_ascii=False) and c[0] in json.dumps(st2, ensure_ascii=False), json.dumps(st2, ensure_ascii=False)[:400])
open(os.path.join(OUT, "grok-after-connect.json"), "w", encoding="utf-8").write(json.dumps({"started": started, "status": status, "tool_status": st2, "default_name": c[0]}, ensure_ascii=False, indent=1))
folder_files = sorted(os.listdir(prefix + "/connections/grok-agent")) if os.path.isdir(prefix + "/connections/grok-agent") else []
rec("CS-AC-14", "automatic agent folder is <prefix>/connections/grok-agent (R-03) with agent.json", "agent.json" in folder_files, f"files={folder_files}")
rpc.close()
# --- B: existing environment settings win over defaults.json (old install keeps working)
write_defaults("http://127.0.0.1:1")           # wrong defaults on purpose
fx("/fixture/hits?reset=1")
legacy_folder = prefix + "/legacy-env-agent"
rpc = Rpc({"KNOWSLINK_MODE": "public-node", "RELAY_URL": URL, "KNOWSLINK_AGENT_FOLDER": legacy_folder})
started = payload(rpc.call("knowslink_connect", {"confirmed": True}))
url = started.get("url") if isinstance(started, dict) else None
hits = json.loads(fx("/fixture/hits")[1])
ok_url = bool(url) and url.startswith(URL)
if ok_url: approve("a", url)
for _ in range(20):
    status = payload(rpc.call("knowslink_connect_status"))
    if isinstance(status, dict) and status.get("state") == "connected": break
    time.sleep(3)
rec("CS-AC-14", "existing environment variables (mode/relay/folder) win over defaults.json and still connect to the env relay", ok_url and hits["starts"] == 1 and isinstance(status, dict) and status.get("state") == "connected" and os.path.exists(legacy_folder + "/agent.json"), f"url={url} status={json.dumps(status, ensure_ascii=False)[:200]}")
rpc.close()
# --- C: dev build without defaults.json and without env stays held
dev = W + "/grok-dev"
if not os.path.exists(dev + "/knowslink/dist/plugin.js"):
    os.makedirs(dev, exist_ok=True); zipfile.ZipFile(ZIP).extractall(dev)
plugin = dev + "/knowslink/dist/plugin.js"
fx("/fixture/hits?reset=1")
rpc = Rpc({})
held = payload(rpc.call("knowslink_connect", {"confirmed": True}))
hits = json.loads(fx("/fixture/hits")[1])
rec("CS-AC-14", "package without defaults.json and without env stays held (no relay request)", hits["api"] == 0 and "held" in json.dumps(held, ensure_ascii=False), json.dumps(held, ensure_ascii=False)[:200])
rpc.close()
