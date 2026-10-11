"""Build inert consent fixtures from the fixed candidate's exact template and shared CSS.
Only synthetic values and the template's two existing branches are expanded.
"""
import hashlib
import html
import json
import os
from pathlib import Path
import re
import subprocess

SHA = "febd1282df2e84a29fa0a53834569e7b9cc028e2"
ROOT = Path(__file__).resolve().parents[5]
OUT = Path(__file__).resolve().parent


def source(path):
    return subprocess.check_output(["git", "-C", str(ROOT), "show", f"{SHA}:{path}"])


remote = source("internal/relay/remote_auth.go")
member = source("internal/relay/member.go")
style = re.search(r"const memberStyle = `([^`]+)`", member.decode()).group(1)
parts = re.search(r'var remotePage = .*?Parse\(`([^`]+)` \+ memberStyle \+ `([^`]+)`\)\)', remote.decode())
template = parts.group(1) + style + parts.group(2)
branch = re.search(r"{{if .Login}}(.*?){{else}}(.*?){{end}}", template)
assert branch is not None and template.count("{{if ") == 1
assert "overflow-wrap:anywhere" in style
fixtures = []
for long_uri in (False, True):
    redirect = "https://client.example/oauth/callback"
    if long_uri:
        redirect += "/" + "syntheticcallbacksegment" * 20 + "?fixture=long&review=visual"
    for login in (True, False):
        name = ("login" if login else "consent") + ("-long" if long_uri else "")
        body = template[:branch.start()] + branch.group(1 if login else 2) + template[branch.end():]
        values = {"Name": "Synthetic Work client", "Redirect": redirect,
                  "ID": "synthetic-visual-attempt", "CSRF": "synthetic-not-a-credential"}
        for key, value in values.items():
            body = body.replace("{{." + key + "}}", html.escape(value, quote=True))
        assert "{{" not in body and style in body
        assert body.count("<button") == (1 if login else 2)
        assert ("Google 계정 확인" in body) == login
        assert ("이 설치 연결 승인" in body and ">취소</button>" in body) == (not login)
        target = OUT / (name + ".html")
        target.write_text(body, encoding="utf-8")
        fixtures.append({"path": target.name, "login": login, "long_uri": long_uri,
                         "redirect_length": len(redirect),
                         "sha256": hashlib.sha256(target.read_bytes()).hexdigest()})
manifest = {
    "source_sha": SHA, "session_id": os.environ.get("CODEX_SESSION_ID"),
    "terminal": os.environ.get("ORCA_TERMINAL_HANDLE"),
    "task_id": "task_415d9cd925c6", "dispatch_id": "ctx_4283f224b494",
    "source": {p: {"blob": subprocess.check_output(["git", "-C", str(ROOT), "rev-parse", f"{SHA}:{p}"]).decode().strip(),
                   "sha256": hashlib.sha256(data).hexdigest()}
               for p, data in (("internal/relay/remote_auth.go", remote), ("internal/relay/member.go", member))},
    "member_style_sha256": hashlib.sha256(style.encode()).hexdigest(),
    "fixtures": fixtures,
    "method": "Exact fixed git blobs; bounded template expansion; HTML-escaped synthetic values; no server or account",
}
(OUT / "source-manifest.json").write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
