"""SAR-AUTOSETUP-001-DEV scenarios: N1-N6/U1-U4 per host, AC-05-AC-09, F-01-F-09 in KO and EN, AS-10 choice.
usage: python scenarios.py <work dir> main|failures-ko|failures-en|pick ...   (fixture on 127.0.0.1:18411)
"""
import json
import os
import re
import shutil
import subprocess
import sys
import tomllib
from pathlib import Path

from drive import Box, agent, approve, cancel, check, fixture, secrets_absent, state, tree_hash

SERVED = Path(os.environ.get("AUTOSETUP_SERVED", r"C:\Users\shin\AppData\Local\KnowsLinkDevTools\autosetup-001-dev\rel2\0.1.0"))
REAL_CLAUDE = shutil.which("claude.exe")
YES = {"claude": ["y", "y"], "codex": ["y"]}  # Claude Code 2.1.308 is not the verified 2.1.302, so it asks first (AS-09)
WORDS = {"ko": {"done": "설치와 키 저장 완료", "updated": "갱신 완료", "declined": "바뀐 것이 없습니다", "unsure": "[확인 필요]"},
         "en": {"done": "Install and key storage complete", "updated": "Update complete", "declined": "Nothing was changed", "unsure": "[Check needed]"}}


def seed(box, app):
    """Things the installer must not read or change (AC-05/06): other agent folders, Grok/Dots, host auth, other MCP entries."""
    (box.kl / f"{app}-01-private").mkdir(parents=True)  # an orphan agent folder: the install must pick another name
    (box.kl / f"{app}-01-private" / "agent.json").write_text('{"agent":"agent_orphan"}')
    (box.kl / "grok-other-private").mkdir()
    (box.kl / "grok-other-private" / "agent.json").write_text('{"agent":"agent_grok"}')
    (box.root / "home" / ".grok").mkdir(parents=True)
    (box.root / "home" / ".grok" / "config.json").write_text('{"mcp":{"knowslink":{"env":{"KNOWSLINK_AGENT_FOLDER":"x"}}}}')
    (box.root / "claude" / ".credentials.json").write_text('{"claudeAiOauth":{"accessToken":"fixture-not-real"}}')
    (box.root / "claude" / "settings.json").write_text('{"model":"opus","permissions":{"allow":["Bash(ls)"]}}')
    (box.root / "claude" / ".claude.json").write_text('{"mcpServers":{"other":{"command":"other-mcp"}}}')
    (box.root / "codex" / "auth.json").write_text('{"OPENAI_API_KEY":"fixture-not-real"}')
    (box.root / "codex" / "config.toml").write_text('model = "gpt-fixture"\napproval_policy = "on-request"\n\n[mcp_servers.other]\ncommand = "other-mcp"\n')
    # A used host: Claude Code 2.1.308 migrates the model alias once on its first plugin command in a new config dir
    # (opus -> opus[1m]); a real user's settings already went through it, so the fixture starts after that.
    subprocess.run([REAL_CLAUDE, "plugin", "list", "--json"], env=box.env, capture_output=True)
    subprocess.run("codex plugin list --json", shell=True, env=box.env, capture_output=True)


def protected(box):
    """Hashes of everything outside the KnowsLink entries; host files compared without the two KnowsLink entries."""
    out = {name: tree_hash(box.kl / name) for name in ("claude-01-private", "codex-01-private", "grok-other-private")}
    out["grok"] = tree_hash(box.root / "home")
    out["claude-auth"] = tree_hash(box.root / "claude" / ".credentials.json")
    out["codex-auth"] = tree_hash(box.root / "codex" / "auth.json")
    settings = json.loads((box.root / "claude" / "settings.json").read_text())
    settings.get("extraKnownMarketplaces", {}).pop("knowslink-plugins", None)
    settings.get("enabledPlugins", {}).pop("knowslink@knowslink-plugins", None)
    out["claude-settings"] = json.dumps({k: v for k, v in settings.items() if v != {}}, sort_keys=True)
    out["claude-mcp"] = json.dumps(json.loads((box.root / "claude" / ".claude.json").read_text()).get("mcpServers"), sort_keys=True)
    toml = tomllib.loads((box.root / "codex" / "config.toml").read_text())
    toml.get("marketplaces", {}).pop("knowslink-plugins", None)
    toml.get("plugins", {}).pop("knowslink@knowslink-plugins", None)
    out["codex-config"] = json.dumps({k: v for k, v in toml.items() if v != {}}, sort_keys=True)
    return out


def registered(box, app):
    if app == "claude":
        r = subprocess.run([REAL_CLAUDE, "plugin", "list", "--json"], env=box.env, capture_output=True, text=True, encoding="utf-8")
        entry = next(e for e in json.loads(r.stdout) if e["id"] == "knowslink@knowslink-plugins")
        return entry["enabled"], entry["mcpServers"]["knowslink"]
    r = subprocess.run("codex mcp list --json", shell=True, env=box.env, capture_output=True, text=True, encoding="utf-8")
    entry = next(e for e in json.loads(r.stdout) if e["name"] == "knowslink")
    return entry["enabled"], entry["transport"]


def new_and_update(app, lang_new, lang_update):
    box = Box(f"{app}-main", fresh_node=app == "claude")  # the first install downloads and checks the real Node zip
    seed(box, app)
    before, counts = protected(box), state()
    code, out = box.run(app, lang_new, YES[app], approve(f"{app} 자동 설정"), label=f"{app}-new")
    install, private = box.kl / f"{app}-02", box.kl / f"{app}-02-private"
    check(code == 0 and WORDS[lang_new]["done"] in out, f"{app} N1-N6 new install exit 0 ({lang_new})")
    enabled, mcp = registered(box, app)
    check(enabled and mcp["args"][0] == str(install / "knowslink" / "dist" / "plugin.js") and mcp["env"]["KNOWSLINK_AGENT_FOLDER"] == str(private),
          f"{app} host entry points at the new install and agent folder")
    check(Path(mcp["command"]).name == "node.exe" and "runtime" in mcp["command"], f"{app} MCP runs the pinned runtime Node")
    check(protected(box) == before, f"{app} AC-05 other agent folders, Grok/Dots, host auth and other MCP/settings entries unchanged")
    check(not (box.kl / f"{app}-01").exists(), f"{app} AC-06 orphan {app}-01-private left alone; install used {app}-02")
    check(secrets_absent(out, private, install), f"{app} AC-06 credential and private key absent from output and install record")
    data, _ = agent(private)
    after_new = state()
    check(after_new.get("agents", 0) == counts.get("agents", 0) + 1, f"{app} one new agent on the web ({after_new})")
    private_hash = tree_hash(private)
    code, out2 = box.run(app, lang_update, YES[app], label=f"{app}-update")
    check(code == 0 and WORDS[lang_update]["updated"] in out2 and "/connect/" not in out2, f"{app} U1-U4 same command updates without Google ({lang_update})")
    check(tree_hash(private) == private_hash and agent(private)[0]["agent"] == data["agent"], f"{app} AC-07 agent folder (ID, key) byte-identical after update")
    check(state() == after_new, f"{app} AC-07 web agent and friend counts unchanged {state()}")
    check(any(p.name.startswith(f"{app}-02.previous-") for p in box.kl.iterdir()), f"{app} AS-03 previous install kept")
    check(protected(box) == before, f"{app} AC-05 still unchanged after update")
    check(secrets_absent(out2, private, install), f"{app} AC-06 no credential in update output or record")


def failures(lang):
    other = "en" if lang == "ko" else "ko"
    # F-01: host missing from PATH.
    box = Box(f"f01-{lang}")
    env = dict(box.env, PATH=r"C:\Windows\System32;C:\Windows\System32\WindowsPowerShell\v1.0")
    code, out = box.run("codex", lang, YES["codex"], env=env, label=f"f01-{lang}")
    check(code == 1 and "[F-01]" in out and not (box.kl / "codex-01").exists(), f"F-01 host missing: message, nothing changed ({lang})")
    # F-02: unverified host version, user declines; then the plan itself declined.
    box = Box(f"f02-{lang}")
    code, out = box.run("claude", lang, ["n"], label=f"f02-{lang}")
    check(code == 1 and "[F-02]" in out and not (box.kl / "claude-01").exists(), f"F-02 declined unverified version: nothing changed ({lang})")
    code, out = box.run("claude", lang, ["y", "n"], label=f"decline-{lang}")
    check(WORDS[lang]["declined"] in out and not (box.kl / "claude-01").exists(), f"AS-02 plan declined: nothing written ({lang})")
    # F-03 / F-04: missing and altered download, then the same command completes (AC-09).
    served = SERVED / "knowslink-claude-plugin.zip"
    hold = served.with_suffix(".hold")
    for fid, tamper in (("F-03", None), ("F-04", SERVED / "knowslink-codex-plugin.zip")):
        box = Box(f"{fid.lower()}-{lang}")
        served.rename(hold)
        try:
            if tamper:
                shutil.copyfile(tamper, served)
            code, out = box.run("claude", lang, YES["claude"], label=f"{fid.lower()}-{lang}")
        finally:
            if served.exists():
                served.unlink()
            hold.rename(served)
        leftovers = list((box.root / "tmp").iterdir())
        check(code == 1 and f"[{fid}]" in out and not (box.kl / "claude-01").exists() and not leftovers, f"{fid} download failure: nothing changed, temp removed ({lang})")
        code, out = box.run("claude", other, YES["claude"], approve("재실행"), label=f"{fid.lower()}-rerun-{other}")
        check(code == 0 and WORDS[other]["done"] in out, f"AC-09 {fid} rerun completes without cleanup ({other})")
    # F-05: host rejects the plugin install; rerun resumes from registration.
    box = Box(f"f05-{lang}")
    shim = box.root / "shim"
    shim.mkdir(exist_ok=True)
    (shim / "claude.cmd").write_text(f'@echo off\r\nif "%1 %2"=="plugin install" (echo simulated install failure 1>&2 & exit /b 1)\r\n"{REAL_CLAUDE}" %*\r\n')
    env = dict(box.env, PATH=str(shim) + ";" + box.env["PATH"])
    code, out = box.run("claude", lang, YES["claude"], env=env, label=f"f05-{lang}")
    check(code == 1 and "[F-05]" in out and "simulated install failure" in out and (box.kl / "claude-01").exists(), f"F-05 registration failure: summary shown, install folder kept ({lang})")
    code, out = box.run("claude", other, YES["claude"], approve("재실행"), label=f"f05-rerun-{other}")
    check(code == 0 and WORDS[other]["done"] in out and not any(p.name.startswith("claude-01.previous") for p in box.kl.iterdir()),
          f"AC-09 F-05 rerun resumes at registration, files kept ({other})")
    # F-06: browser cancel; rerun starts a new connection.
    box = Box(f"f06-{lang}")
    code, out = box.run("claude", lang, YES["claude"], cancel, label=f"f06-{lang}")
    check(code == 1 and "[F-06]" in out and not (box.kl / "claude-01-private" / "agent.json").exists(), f"F-06 cancelled: install and registration kept, no key ({lang})")
    first = re.search(r"/connect/(\S+)", out).group(1)
    code, out = box.run("claude", other, YES["claude"], approve("재실행"), label=f"f06-rerun-{other}")
    check(code == 0 and re.search(r"/connect/(\S+)", out).group(1) != first, f"AC-09 F-06 rerun uses a new connection ({other})")
    # F-07: approval reached the relay, the credential answer was lost.
    box = Box(f"f07-{lang}")
    fixture("/fixture/lose-complete", on="1")
    try:
        code, out = box.run("claude", lang, YES["claude"] + [""], approve("응답 유실"), label=f"f07-{lang}")
    finally:
        fixture("/fixture/lose-complete", on="0")
    check(code == 1 and "[F-07]" in out and out.count("/connect/") == 1 and "응답 유실" in out, f"F-07 lost answer: check-needed with the agent name, no new attempt ({lang})")
    code, out = box.run("claude", other, YES["claude"] + ["new"], approve("재실행"), label=f"f07-rerun-{other}")
    check(code == 0 and "[F-07]" in out and out.count("/connect/") == 1 and WORDS[other]["done"] in out, f"AC-09 F-07 rerun re-checks first, new connection only after 'new' ({other})")
    # F-08: the install folder is in use during an update; the previous install stays.
    box = Box(f"f08-{lang}")
    box.run("claude", lang, YES["claude"], approve("갱신 실패"), label=f"f08-setup-{lang}")
    install = box.kl / "claude-01"
    before = tree_hash(install)
    held = open(install / "knowslink" / "README.md", "rb")
    try:
        code, out = box.run("claude", lang, YES["claude"], label=f"f08-{lang}")
    finally:
        held.close()
    enabled, mcp = registered(box, "claude")
    check(code == 1 and "[F-08]" in out and tree_hash(install) == before and enabled and Path(mcp["args"][0]).is_file(),
          f"F-08 swap failure: previous install byte-identical and still registered ({lang})")
    code, out = box.run("claude", other, YES["claude"], label=f"f08-rerun-{other}")
    check(code == 0 and WORDS[other]["updated"] in out, f"AC-09 F-08 rerun updates ({other})")
    # F-09: Codex YOLO approval mode is reported, never changed.
    box = Box(f"f09-{lang}")
    config = box.root / "codex" / "config.toml"
    config.write_text('approval_policy = "never"\nsandbox_mode = "danger-full-access"\n')
    code, out = box.run("codex", lang, YES["codex"], approve("Codex YOLO"), label=f"f09-{lang}")
    check(code == 0 and "[F-09]" in out and WORDS[lang]["done"] in out and 'approval_policy = "never"' in config.read_text(),
          f"F-09 YOLO: install complete, mode reported and unchanged ({lang})")


def pick():
    box = Box("pick")
    box.run("claude", "ko", YES["claude"], approve("첫째"), label="pick-first")
    code, out = box.run("claude", "ko", YES["claude"], approve("둘째"), extra="-NewAgent", label="pick-new-agent")
    check(code == 0 and (box.kl / "claude-02").exists(), "AS-10 -NewAgent installs a separate agent")
    code, out = box.run("claude", "en", ["y", "x"], label="pick-unsure")
    check(code == 1 and WORDS["en"]["unsure"] in out, "AS-10 several installs: list shown, unclear choice stops")


SCENARIOS = {"main": lambda: (new_and_update("claude", "ko", "en"), new_and_update("codex", "en", "ko")),
             "failures-ko": lambda: failures("ko"), "failures-en": lambda: failures("en"), "pick": pick}

if __name__ == "__main__":
    for name in sys.argv[2:]:
        print(f"== {name}", flush=True)
        SCENARIOS[name]()
