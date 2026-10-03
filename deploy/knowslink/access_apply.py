"""Create or remove the owner-only Cloudflare Access app for link.knowslog.com; never prints the email or token."""
import json
import os
import re
import stat
import time
import sys
import urllib.error
import urllib.request

HOST = "link.knowslog.com"
EMAIL_FILE = os.environ.get("KNOWSLINK_OWNER_EMAIL_FILE", "/tmp/knowslink-beta-owner-email")
STATE = os.environ.get("KNOWSLINK_STATE_DIR", "/home/shin/deploy/knowslink-state")
POLICY_NAME = "knowslink-beta-owner-only"
APP_NAME = "KnowsLink beta (owner-only)"


def private_text(path):
    info = os.stat(path)
    if stat.S_IMODE(info.st_mode) & 0o077 or info.st_uid != os.getuid():
        raise SystemExit(f"{path} must be owned by this user with mode 0600")
    return open(path).read().strip()


def policy_body(email):
    return {"name": POLICY_NAME, "decision": "allow", "session_duration": "24h",
            "include": [{"email": {"email": email}}], "exclude": [], "require": []}


def app_body(policy_id, idp_id):
    return {"name": APP_NAME, "type": "self_hosted", "domain": HOST,
            "destinations": [{"type": "public", "uri": HOST}], "session_duration": "24h",
            "allowed_idps": [idp_id], "auto_redirect_to_identity": True,
            "app_launcher_visible": False, "policies": [{"id": policy_id, "precedence": 1}]}


def only_owner(policy, email):
    return policy["decision"] == "allow" and policy["include"] == [{"email": {"email": email}}] \
        and not policy["exclude"] and not policy["require"]


def verify_live(app, policy, email, aud, team, config_text):
    """Raise AssertionError unless the live Access app, policy and tunnel config agree and admit only the owner."""
    assert app["type"] == "self_hosted" and app.get("domain") in (None, HOST), "app type or domain"
    assert [(d["type"], d["uri"]) for d in app["destinations"]] == [("public", HOST)], "destinations"
    assert [p["id"] for p in app["policies"]] == [policy["id"]], "app must carry only the owner policy"
    assert only_owner(policy, email), "policy must allow only the owner email"
    assert len(app.get("allowed_idps") or []) == 1, "exactly one identity provider"
    assert not app.get("options_preflight_bypass") and not app.get("custom_pages"), "bypass options"
    assert app["aud"] == aud, "app aud differs from the recorded aud"
    assert re.search(rf"^\s+teamName: {re.escape(team)}$", config_text, re.M), "tunnel teamName"
    assert re.findall(r"^\s+- ([0-9a-f]{64})$", config_text, re.M) == [aud], "tunnel audTag must be exactly the app aud"


def selftest():
    body = app_body("p", "i")
    assert body["destinations"] == [{"type": "public", "uri": HOST}] and len(body["policies"]) == 1
    assert only_owner(policy_body("x@example.test") | {"id": "p"}, "x@example.test")
    assert not only_owner(policy_body("x@example.test"), "y@example.test")
    email, aud = "x@example.test", "a" * 64
    policy = policy_body(email) | {"id": "p"}
    app = app_body("p", "i") | {"aud": aud}
    config = f"        teamName: team\n        audTag:\n          - {aud}\n"
    verify_live(app, policy, email, aud, "team", config)
    for broken in (app | {"aud": "b" * 64}, app | {"options_preflight_bypass": True},
                   app | {"policies": [{"id": "q"}]}, app | {"destinations": []}):
        try:
            verify_live(broken, policy, email, aud, "team", config)
        except AssertionError:
            continue
        raise SystemExit("verify_live accepted a broken app")
    for broken_config in (config.replace(aud, ""), config.replace("team", "other")):
        try:
            verify_live(app, policy, email, aud, "team", broken_config)
        except AssertionError:
            continue
        raise SystemExit("verify_live accepted a broken tunnel config")
    print("access_apply selftest ok")


def call(token, method, path, body=None):
    request = urllib.request.Request(
        "https://api.cloudflare.com/client/v4" + path, method=method,
        data=None if body is None else json.dumps(body).encode(),
        headers={"Authorization": "Bearer " + token, "Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            data = json.load(response)
    except urllib.error.HTTPError as error:
        raise SystemExit(f"{method} {path.split('?')[0]} failed: HTTP {error.code}")
    return data["result"]


def apply(token, email, account):
    base = f"/accounts/{account}/access"
    if any(HOST in (app.get("domain") or "") for app in call(token, "GET", base + "/apps?per_page=100")):
        raise SystemExit(f"an Access app for {HOST} already exists; not overwriting")
    idps = [i for i in call(token, "GET", base + "/identity_providers") if i["type"] == "onetimepin"]
    assert len(idps) == 1, "expected exactly one One-time PIN identity provider"
    named = [p for p in call(token, "GET", base + "/policies") if p["name"] == POLICY_NAME]
    policy = named[0] if named else call(token, "POST", base + "/policies", policy_body(email))
    assert only_owner(policy, email), "existing policy does not allow only the owner email"
    app = call(token, "POST", base + "/apps", app_body(policy["id"], idps[0]["id"]))
    check = call(token, "GET", f"{base}/apps/{app['id']}")
    assert check["aud"] == app["aud"]
    assert [p["id"] for p in check["policies"]] == [policy["id"]], "app must carry only the owner policy"
    with open(os.path.join(STATE, "access.json"), "w") as out:
        json.dump({"app": app["id"], "policy": policy["id"], "policy_created": not named}, out)
    with open(os.path.join(STATE, "access.aud"), "w") as out:
        out.write(app["aud"])
    print("access app created; aud and ids stored in the state directory")


def live_snapshot(account):
    """Live app and policy: fetched with a token file, else a fresh (10 min) snapshot saved from a read-only MCP call."""
    saved = json.load(open(os.path.join(STATE, "access.json")))
    if os.environ.get("CF_API_TOKEN_FILE"):
        token = private_text(os.environ["CF_API_TOKEN_FILE"])
        base = f"/accounts/{account or call(token, 'GET', '/accounts')[0]['id']}/access"
        return call(token, "GET", f"{base}/apps/{saved['app']}"), call(token, "GET", f"{base}/policies/{saved['policy']}")
    path = os.path.join(STATE, "access.live.json")
    if time.time() - os.stat(path).st_mtime > 600:
        raise SystemExit("access.live.json is older than 10 minutes; refresh it with a read-only call")
    snapshot = json.load(open(path))
    return snapshot["app"], snapshot["policy"]


def check():
    app, policy = live_snapshot(None)
    config_text = open(os.path.join(STATE, "tunnel", "config.yml")).read()
    verify_live(app, policy, private_text(EMAIL_FILE), open(os.path.join(STATE, "access.aud")).read().strip(),
                os.environ.get("KNOWSLINK_TEAM", "scshin88"), config_text)
    print("live Access app, owner-only policy and tunnel config agree")


def remove(token, account):
    saved = json.load(open(os.path.join(STATE, "access.json")))
    base = f"/accounts/{account}/access"
    call(token, "DELETE", f"{base}/apps/{saved['app']}")
    if saved["policy_created"]:
        call(token, "DELETE", f"{base}/policies/{saved['policy']}")
    os.remove(os.path.join(STATE, "access.json"))
    os.remove(os.path.join(STATE, "access.aud"))
    print("access app removed")


if __name__ == "__main__":
    command = sys.argv[1] if len(sys.argv) > 1 else ""
    if command == "selftest":
        selftest()
    elif command == "check":
        check()
    elif command in ("apply", "remove"):
        token = private_text(os.environ["CF_API_TOKEN_FILE"])
        accounts = call(token, "GET", "/accounts")
        assert len(accounts) == 1, "the token must see exactly one account"
        account = accounts[0]["id"]
        if command == "apply":
            apply(token, private_text(EMAIL_FILE), account)
        else:
            remove(token, account)
    else:
        raise SystemExit("usage: CF_API_TOKEN_FILE=<0600 file> access_apply.py apply|remove|check|selftest")
