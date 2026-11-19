"""Run candidate regression, standalone tool probe, Grok CLI isolation, and document alignment."""
import hashlib
import json
import shutil
import subprocess
from pathlib import Path

CLONE = Path("/tmp/sar-mvp-003-qa")
OUT = Path(__file__).resolve().parent / "logs"
NODE = shutil.which("node")
GROK_PATHS = [
    Path("/home/shin/.grok/config.toml"),
    Path("/home/shin/.grok/trusted_folders.toml"),
]


def run(name, command, cwd):
    result = subprocess.run(command, cwd=cwd, text=True, capture_output=True)
    (OUT / f"{name}.out").write_text(result.stdout)
    (OUT / f"{name}.err").write_text(result.stderr)
    (OUT / f"{name}.exit").write_text(str(result.returncode) + "\n")
    return result.returncode


def digest(path):
    if not path.exists():
        return "absent"
    if path.is_file():
        return hashlib.sha256(path.read_bytes()).hexdigest()
    rows = []
    for child in sorted(path.rglob("*")):
        if child.is_file():
            rows.append(hashlib.sha256(child.read_bytes()).hexdigest() + " " + str(child.relative_to(path)))
    return hashlib.sha256("\n".join(rows).encode()).hexdigest()


def grok_fingerprint():
    plugins = Path("/home/shin/.grok/installed-plugins")
    return {str(path): digest(path) for path in [*GROK_PATHS, plugins]}


def bare_tools():
    plugin = CLONE / "adapters/dist/plugin.js"
    proc = subprocess.run(["env", "-i", NODE, str(plugin)], input=handshake(), text=True, capture_output=True, cwd="/tmp", timeout=20)
    (OUT / "env-i.out").write_text(proc.stdout)
    (OUT / "env-i.err").write_text(proc.stderr)
    (OUT / "env-i.exit").write_text(str(proc.returncode) + "\n")
    names = []
    status = ""
    for line in proc.stdout.splitlines():
        if not line.strip():
            continue
        message = json.loads(line)
        tools = message.get("result", {}).get("tools")
        if tools is not None:
            names = sorted(tool["name"] for tool in tools)
        content = message.get("result", {}).get("content")
        if content:
            status = content[0].get("text", "")
    expected = ["knowslink_pull_once", "knowslink_status", "knowslink_test_receive", "knowslink_test_send"]
    held = '"state":"held"' in status and '"actualConnection":"held"' in status
    ok = names == expected and held and proc.stderr == ""
    return {"ok": ok, "names": names, "status": status, "stderr_len": len(proc.stderr), "exit": proc.returncode}


def handshake():
    calls = [
        {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-11-25", "capabilities": {}, "clientInfo": {"name": "env-i", "version": "0"}}},
        {"jsonrpc": "2.0", "method": "notifications/initialized"},
        {"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}},
        {"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "knowslink_status", "arguments": {}}},
    ]
    return "".join(json.dumps(call) + "\n" for call in calls)


def documents():
    needed = ["scripts/run_trial.py", "--config", "--node", "--bundle", "receive", "knowslink_test_send", "knowslink_test_receive", "knowslink_status", "knowslink_pull_once", "idempotency_key", "test-remote", "https://link.knowslog.com", "0600", "untrusted"]
    files = {
        "exec": CLONE / ".fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md",
        "ops": CLONE / ".fullops-squad/docs/operations/ops-guide.md",
        "readme": CLONE / "adapters/README.md",
    }
    missing = {name: [item for item in needed if item not in path.read_text()] for name, path in files.items()}
    return {"ok": not any(missing.values()), "missing": missing}


def main():
    OUT.mkdir(exist_ok=True)
    before = grok_fingerprint()
    codes = {
        "make-test": run("make-test", ["make", "test"], CLONE),
        "verify-mvp": run("verify-mvp", ["make", "verify-mvp"], CLONE),
        "verify-grok-plugin": run("verify-grok-plugin", ["make", "verify-grok-plugin"], CLONE),
    }
    after = grok_fingerprint()
    bare = bare_tools()
    docs = documents()
    clone_status = subprocess.run(["git", "-C", str(CLONE), "status", "--porcelain"], text=True, capture_output=True)
    summary = {
        "codes": codes,
        "bare": bare,
        "documents": docs,
        "grok_unchanged": before == after,
        "grok_before": before,
        "grok_after": after,
        "clone_head": subprocess.check_output(["git", "-C", str(CLONE), "rev-parse", "HEAD"], text=True).strip(),
        "clone_status": clone_status.stdout,
    }
    (OUT / "suite.json").write_text(json.dumps(summary, indent=2) + "\n")
    ok = all(code == 0 for code in codes.values()) and bare["ok"] and docs["ok"] and before == after and clone_status.stdout == "" and summary["clone_head"] == "cd60e7f87eb5ce137eca887980f232b3f67a18d0"
    print(json.dumps({"ok": ok, "codes": codes, "bare": bare, "documents_ok": docs["ok"], "grok_unchanged": before == after, "clone_status": clone_status.stdout}, indent=2))
    raise SystemExit(0 if ok else 1)


if __name__ == "__main__":
    main()
