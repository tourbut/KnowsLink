"""Check the render-only Google member ingress fragment without changing Cloudflare or deployment state."""
from pathlib import Path
import re

source = (Path(__file__).resolve().parents[1] / "deploy/knowslink/tunnel/public-ingress.yml").read_text()
pattern = re.search(r"^  path: '(.+)'$", source, re.M)
assert pattern, "public path rule required"
public = re.compile(pattern[1])
for path in ("/", "/auth/google", "/auth/google/callback", "/auth/reauth", "/home", "/home/device-confirm",
             "/connect/" + "a" * 43, "/v1/connect/start", "/v1/connect/poll", "/v1/connect/complete",
             "/v1/text/send", "/v1/keys/agent_a/key_a", "/v1/receipts/00000000-0000-7000-8000-000000000001"):
    assert public.fullmatch(path), f"member path blocked: {path}"
for path in ("/owner", "/owner/gates/x", "/v1/owners", "/v1/agents", "/v1/keys", "/v1/key-revoke",
             "/v1/owner-revoke", "/v1/authorize", "/v1/test/pull", "/healthz", "/_probe", "/homepage",
             "/connect/short", "/v1/connect/unknown", "/v1/keys/agent/key/extra"):
    assert not public.fullmatch(path), f"protected path exposed: {path}"
assert "hostname: link.knowslog.com" in source and "service: http://relay:8080" in source
assert "required: false" in source
print("Public ingress: Google/member paths only; owner, administrative APIs and trial remain outside this rule")
