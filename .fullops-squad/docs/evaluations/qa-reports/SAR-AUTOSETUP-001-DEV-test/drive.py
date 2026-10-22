"""SAR-AUTOSETUP-001-DEV installer driver on Windows against the isolated autosetup fixture (127.0.0.1:18411).

Each "box" isolates LOCALAPPDATA, CLAUDE_CONFIG_DIR, CODEX_HOME and TEMP; the user's real ~/.claude and ~/.codex are
never used. The fixture's /fixture/* endpoints stand in for the person's browser approval or cancel.
usage: python drive.py <work dir> <scenario> [...]   (scenarios: see SCENARIOS at the end)
"""
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import threading
import time
import urllib.parse
import urllib.request
from pathlib import Path

RELAY = "http://127.0.0.1:18411"
BASE = f"{RELAY}/install/0.1.0"
WORK = Path(sys.argv[1])
NODE_CACHE = WORK / "node-runtime"


def fixture(path, **form):
    data = urllib.parse.urlencode(form).encode()
    with urllib.request.urlopen(urllib.request.Request(RELAY + path, data=data, method="POST")) as r:
        return r.read().decode()


def state():
    with urllib.request.urlopen(RELAY + "/fixture/state") as r:
        return json.loads(r.read())


class Box:
    def __init__(self, name, fresh_node=False):
        if name != "main" and not name.endswith("-main"):
            fixture("/fixture/revoke-all")
        self.root = WORK / name
        for sub in ("local", "claude", "codex", "tmp"):
            (self.root / sub).mkdir(parents=True, exist_ok=True)
        self.kl = self.root / "local" / "KnowsLink"
        runtime = self.kl / "runtime"
        if not fresh_node and NODE_CACHE.exists() and not runtime.exists():
            shutil.copytree(NODE_CACHE, runtime)
        self.env = dict(os.environ, LOCALAPPDATA=str(self.root / "local"), CLAUDE_CONFIG_DIR=str(self.root / "claude"),
                        CODEX_HOME=str(self.root / "codex"), TEMP=str(self.root / "tmp"), TMP=str(self.root / "tmp"))

    def run(self, app, lang, answers, browser=None, extra="", env=None, label=""):
        """browser(connection_id) is called once per printed connection link, like the person in the browser."""
        command = (f"[Console]::OutputEncoding=[Text.Encoding]::UTF8; & ([scriptblock]::Create((irm {BASE}/knowslink-setup.ps1)))"
                   f" -App {app} -Lang {lang} {extra}; exit $LASTEXITCODE")
        started = time.time()
        p = subprocess.Popen(["powershell.exe", "-NoProfile", "-Command", command], stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                             stderr=subprocess.STDOUT, env=env or self.env)
        p.stdin.write(("\n".join(answers) + "\n").encode())
        p.stdin.close()
        lines = []
        for raw in p.stdout:
            line = raw.decode("utf-8", "replace").rstrip()
            lines.append(line)
            m = re.search(r"/connect/([A-Za-z0-9_-]{43})", line)
            if m and browser:
                threading.Thread(target=browser, args=(m.group(1),), daemon=True).start()
        code = p.wait()
        out = "\n".join(lines)
        log = WORK / "logs" / (f"{label}.log" if label else f"{app}-{lang}.log")
        log.parent.mkdir(exist_ok=True)
        log.write_text(f"exit={code} seconds={time.time() - started:.1f}\n{out}\n", encoding="utf-8")
        return code, out


def approve(name):
    def act(connection):
        time.sleep(1.5)
        fixture("/fixture/approve", connection=connection, name=name)
    return act


def cancel(connection):
    time.sleep(1.5)
    fixture("/fixture/cancel", connection=connection)


def tree_hash(path):
    """Content hash of every file under path (relative names and bytes); missing path hashes as 'absent'."""
    path = Path(path)
    if not path.exists():
        return "absent"
    h = hashlib.sha256()
    if path.is_file():
        return hashlib.sha256(path.read_bytes()).hexdigest()
    for f in sorted(p for p in path.rglob("*") if p.is_file()):
        h.update(str(f.relative_to(path)).encode() + b"\0" + f.read_bytes())
    return h.hexdigest()


def check(cond, message):
    print(("PASS " if cond else "FAIL ") + message, flush=True)
    if not cond:
        raise SystemExit(1)


def agent(private):
    data = json.loads((Path(private) / "agent.json").read_text())
    pem = (Path(private) / "private.pem").read_text()
    return data, pem


def secrets_absent(out, private, install):
    data, pem = agent(private)
    record = (Path(install) / "knowslink-install.json").read_text(encoding="utf-8")
    body = pem.split("-----")[2].strip().splitlines()[0]
    return all(s not in out and s not in record for s in (data["credential"], body))
