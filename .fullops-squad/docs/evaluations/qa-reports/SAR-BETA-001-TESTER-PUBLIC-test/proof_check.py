#!/usr/bin/env python3
"""Compare the preserved Access GET with the live tunnel config. Prints no secret values."""

import json
import os
import re
import time
from pathlib import Path

STATE = Path(os.environ.get("KNOWSLINK_STATE_DIR", "/home/shin/deploy/knowslink-state"))
HOST = "link.knowslog.com"
EMAIL_FILE = Path(os.environ.get("KNOWSLINK_OWNER_EMAIL_FILE", "/tmp/knowslink-beta-owner-email"))


def mode(path):
    return path.stat().st_mode & 0o777


def main():
    config_text = (STATE / "tunnel" / "config.yml").read_text()
    live = json.loads((STATE / "access.live.json").read_text())
    saved = json.loads((STATE / "access.json").read_text())
    aud = (STATE / "access.aud").read_text().strip()
    app = live["app"]
    policy = live["policy"]
    tags = re.findall(r"^\s+- ([0-9a-f]{64})$", config_text, re.M)
    team = re.search(r"^\s+teamName: (\S+)$", config_text, re.M)
    include = policy.get("include") or []
    email_mode_ok = EMAIL_FILE.is_file() and mode(EMAIL_FILE) == 0o600 and EMAIL_FILE.stat().st_uid == os.getuid()
    email_match = False
    if email_mode_ok and len(include) == 1 and set(include[0]) == {"email"}:
        owner = EMAIL_FILE.read_text().strip()
        value = include[0]["email"].get("email") if isinstance(include[0].get("email"), dict) else None
        email_match = bool(owner) and value == owner
    idps = app.get("allowed_idps") or []
    creds = [p for p in (STATE / "tunnel").iterdir() if p.suffix == ".json"]
    age = time.time() - (STATE / "access.live.json").stat().st_mtime
    domains = app.get("self_hosted_domains")
    bypass_names = sorted(k for k in app if "bypass" in k.lower() or "service_token" in k.lower())
    result = {
        "live_0600": mode(STATE / "access.live.json") == 0o600,
        "access_json_0600": mode(STATE / "access.json") == 0o600,
        "access_aud_0600": mode(STATE / "access.aud") == 0o600,
        "config_0600": mode(STATE / "tunnel" / "config.yml") == 0o600,
        "tunnel_dir_0700": mode(STATE / "tunnel") == 0o700,
        "state_0700": mode(STATE) == 0o700,
        "env_0600": mode(STATE / ".env") == 0o600,
        "credential_json_count": len(creds),
        "credential_json_0600": bool(creds) and all(mode(p) == 0o600 for p in creds),
        "email_file_0600": email_mode_ok,
        "live_age_seconds": int(age),
        "live_within_10min": 0 <= age <= 600,
        "app_id_match": app.get("id") == saved.get("app") and bool(saved.get("app")),
        "policy_id_match": policy.get("id") == saved.get("policy") and bool(saved.get("policy")),
        "policy_name_match": policy.get("name") == "knowslink-beta-owner-only",
        "policy_reusable": policy.get("reusable") is True,
        "app_name_match": app.get("name") == "KnowsLink beta (owner-only)",
        "type_self_hosted": app.get("type") == "self_hosted",
        "domain_ok": app.get("domain") in (None, HOST),
        "self_hosted_domains_host_only": domains == [HOST],
        "destinations_exact": [(d.get("type"), d.get("uri")) for d in app.get("destinations") or []] == [("public", HOST)],
        "single_policy": [p.get("id") for p in app.get("policies") or []] == [policy.get("id")],
        "decision_allow": policy.get("decision") == "allow",
        "include_count": len(include),
        "exclude_empty": policy.get("exclude") == [],
        "require_empty": policy.get("require") == [],
        "owner_email_match": email_match,
        "idp_count": len(idps),
        "idp_id_only": len(idps) == 1 and isinstance(idps[0], str) and bool(idps[0]),
        "idp_type_in_proof": isinstance(idps[0], dict) and "type" in idps[0] if idps else False,
        "no_preflight_bypass": app.get("options_preflight_bypass") is False,
        "no_custom_pages": not app.get("custom_pages"),
        "launcher_hidden": app.get("app_launcher_visible") is False,
        "auto_redirect": app.get("auto_redirect_to_identity") is True,
        "bypass_flag_names": bypass_names,
        "bypass_flags_false": all(app.get(name) is False for name in bypass_names),
        "aud_match": app.get("aud") == aud and tags == [aud],
        "aud_tag_count": len(tags),
        "aud_64hex": bool(re.fullmatch(r"[0-9a-f]{64}", aud)),
        "team_match_default": bool(team) and team.group(1) == os.environ.get("KNOWSLINK_TEAM", "scshin88"),
        "required_true": bool(re.search(r"^\s+required:\s*true\s*$", config_text, re.M)),
        "hostname_present": f"hostname: {HOST}" in config_text,
        "relay_service": "service: http://relay:8080" in config_text,
        "fallback_404": "service: http_status:404" in config_text,
        "no_template_placeholders": not any(token in config_text for token in ("__AUD_TAG__", "__TEAM_NAME__", "__TUNNEL_UUID__")),
        "credentials_file_key_present": "credentials-file:" in config_text,
    }
    print(json.dumps(result, indent=2))
    required = [k for k, v in result.items() if k not in ("live_age_seconds", "idp_type_in_proof", "bypass_flag_names")]
    failed = [k for k in required if result[k] is False or result[k] == 0 and k.endswith("_count")]
    # idp_count and aud_tag_count and include_count and credential count must be 1
    for key, expected in (("idp_count", 1), ("aud_tag_count", 1), ("include_count", 1), ("credential_json_count", 1)):
        if result[key] != expected:
            failed.append(key)
    if result["idp_type_in_proof"]:
        failed.append("idp_type_unexpected_object")
    failed = sorted(set(failed))
    print(f"failed={failed}")
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
