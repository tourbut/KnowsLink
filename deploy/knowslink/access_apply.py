"""Create or remove the owner-only Cloudflare Access app for link.knowslog.com; never prints the email or token."""
import json
import os
import stat
import sys
import urllib.error
import urllib.request

HOST = "link.knowslog.com"
EMAIL_FILE = "/tmp/knowslink-beta-owner-email"
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


def selftest():
    body = app_body("p", "i")
    assert body["destinations"] == [{"type": "public", "uri": HOST}] and len(body["policies"]) == 1
    assert only_owner(policy_body("x@example.test") | {"id": "p"}, "x@example.test")
    assert not only_owner(policy_body("x@example.test"), "y@example.test")
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
    if any(HOST in (app.get("domain") or "") for app in call(token, "GET", base + "/apps")):
        raise SystemExit(f"an Access app for {HOST} already exists; not overwriting")
    idps = [i for i in call(token, "GET", base + "/identity_providers") if i["type"] == "onetimepin"]
    assert len(idps) == 1, "expected exactly one One-time PIN identity provider"
    named = [p for p in call(token, "GET", base + "/policies") if p["name"] == POLICY_NAME]
    policy = named[0] if named else call(token, "POST", base + "/policies", policy_body(email))
    assert only_owner(policy, email), "existing policy does not allow only the owner email"
    app = call(token, "POST", base + "/apps", app_body(policy["id"], idps[0]["id"]))
    check = call(token, "GET", f"{base}/apps/{app['id']}")
    assert [p["id"] for p in check["policies"]] == [policy["id"]], "app must carry only the owner policy"
    with open(os.path.join(STATE, "access.json"), "w") as out:
        json.dump({"app": app["id"], "policy": policy["id"], "policy_created": not named}, out)
    with open(os.path.join(STATE, "access.aud"), "w") as out:
        out.write(app["aud"])
    print("access app created; aud and ids stored in the state directory")


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
    elif command in ("apply", "remove"):
        token = private_text(os.environ["CF_API_TOKEN_FILE"])
        account = call(token, "GET", "/accounts")[0]["id"]
        if command == "apply":
            apply(token, private_text(EMAIL_FILE), account)
        else:
            remove(token, account)
    else:
        raise SystemExit("usage: CF_API_TOKEN_FILE=<0600 file> access_apply.py apply|remove|selftest")
