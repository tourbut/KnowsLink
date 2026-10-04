"""Render trial-only Access request bodies and a protected Tunnel candidate; never mutate live Cloudflare or deployment state."""
import argparse
import json
from pathlib import Path
import re

HOST = "link.knowslog.com"
UUID = re.compile(r"^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$")


def identifier(value):
    if not UUID.fullmatch(value):
        raise argparse.ArgumentTypeError("expected a real UUID from the create response")
    return value


def policy(codex, grok):
    if codex == grok:
        raise ValueError("agents require distinct Access service tokens")
    return {"name": "knowslink-trial-agents", "decision": "non_identity", "include": [
        {"service_token": {"token_id": codex}}, {"service_token": {"token_id": grok}}], "exclude": [], "require": []}


def app(policy_id):
    return {"name": "KnowsLink trial messages", "type": "self_hosted", "domain": HOST + "/v1/test/*",
            "destinations": [{"type": "public", "uri": HOST + "/v1/test/*"}],
            "session_duration": "24h", "app_launcher_visible": False,
            "policies": [{"id": policy_id, "precedence": 1}]}


def tunnel(uuid, team, owner_aud, trial_aud):
    if not re.fullmatch(r"[a-z0-9-]+", team) or not all(re.fullmatch(r"[a-f0-9]{64}", x) for x in (owner_aud, trial_aud)) or owner_aud == trial_aud:
        raise ValueError("invalid team or distinct AUD tags required")
    source = (Path(__file__).parent / "tunnel/config.yml.tmpl").read_text()
    source = source.replace("__TUNNEL_UUID__", uuid).replace("__TEAM_NAME__", team).replace("__AUD_TAG__", owner_aud)
    entry = f"""  - hostname: {HOST}
    path: /v1/test/.*
    service: http://relay:8080
    originRequest:
      access:
        required: true
        teamName: {team}
        audTag:
          - {trial_aud}
"""
    return source.replace("ingress:\n", "ingress:\n" + entry, 1)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="action", required=True)
    sub.add_parser("tokens")
    p = sub.add_parser("policy")
    p.add_argument("--codex-token", required=True, type=identifier)
    p.add_argument("--grok-token", required=True, type=identifier)
    p = sub.add_parser("app")
    p.add_argument("--policy", required=True, type=identifier)
    p = sub.add_parser("tunnel")
    p.add_argument("--uuid", required=True, type=identifier)
    p.add_argument("--team", required=True)
    p.add_argument("--owner-aud", required=True)
    p.add_argument("--trial-aud", required=True)
    sub.add_parser("selftest")
    args = parser.parse_args()
    if args.action == "tunnel":
        print(tunnel(args.uuid, args.team, args.owner_aud, args.trial_aud), end="")
        return
    if args.action == "tokens":
        result = [{"name": "knowslink-trial-" + agent, "duration": "24h"} for agent in ("codex", "grok")]
    elif args.action == "policy":
        result = policy(args.codex_token, args.grok_token)
    elif args.action == "app":
        result = app(args.policy)
    else:
        body = app("test-policy")
        rules = policy("test-codex", "test-grok")
        candidate = tunnel("00000000-0000-4000-8000-000000000001", "test-team", "a" * 64, "b" * 64)
        if body["domain"] != HOST + "/v1/test/*" or rules["decision"] != "non_identity" or len(rules["include"]) != 2 or candidate.count("required: true") != 2 or candidate.index("path: /v1/test/.*") > candidate.rindex("hostname:"):
            raise ValueError("trial plan boundary failed")
        try:
            policy("same", "same")
        except ValueError:
            pass
        else:
            raise ValueError("shared token accepted")
        result = {"state": "pass", "mode": "render-only", "liveMutation": False}
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
