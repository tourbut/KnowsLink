"""Beta checks: baseline/regression of shared services, local relay hardening, and public Access negatives."""
import http.client
import json
import os
import re
import subprocess
import sys
import urllib.parse

STATE = os.environ.get("KNOWSLINK_STATE_DIR", "/home/shin/deploy/knowslink-state")
HOST = "link.knowslog.com"
SHARED = ["orca.knowslog.com", "s8.knowslog.com", "mcp.knowslog.com"]


def run(*command):
    return subprocess.run(command, capture_output=True, text=True, check=True).stdout


def fetch(scheme, host, path, headers=None, port=None):
    connection = (http.client.HTTPSConnection if scheme == "https" else http.client.HTTPConnection)(host, port, timeout=15)
    connection.request("GET", path, headers=headers or {})
    response = connection.getresponse()
    return response.status, response.getheader("Location") or "", response.read(200)


def shared_state():
    pids = re.findall(r"^[0-9]+(?= /home/shin/.local/bin/cloudflared tunnel --config /home/shin/\.cloudflared)",
                      run("pgrep", "-af", "cloudflared"), re.M)
    containers = sorted(run("docker", "ps", "-a", "--filter", "name=myportfolio", "--format", "{{.Names}} {{.Status}}").splitlines())
    return {"codes": {h: fetch("https", h, "/")[0] for h in SHARED}, "host_cloudflared_pids": pids,
            "containers": [" ".join(c.split()[:2]) for c in containers],
            "projects": sorted(l.split()[0] for l in run("docker", "compose", "ls").splitlines()[1:])}


def baseline():
    with open(os.path.join(STATE, "shared-baseline.json"), "w") as out:
        json.dump(shared_state(), out)
    print("baseline stored")


def regression():
    saved = json.load(open(os.path.join(STATE, "shared-baseline.json")))
    now = shared_state()
    now["projects"] = [p for p in now["projects"] if p != "knowslink"]
    assert now == saved, f"shared services changed: {now} != {saved}"
    print("shared services unchanged:", saved["codes"], "pids", saved["host_cloudflared_pids"])


def local():
    assert fetch("http", "127.0.0.1", "/healthz", port=8080)[0] == 200
    assert fetch("http", "127.0.0.1", "/_probe", port=8080)[0] == 404
    for path in ("/owner", "/v1/contacts", "/v1/receipts/x"):
        assert fetch("http", "127.0.0.1", path, port=8080)[0] == 401, path
    listeners = run("ss", "-ltnH")
    assert re.search(r"127\.0\.0\.1:8080\b", listeners) and not re.search(r"(0\.0\.0\.0|\[::\]|\*):8080\b", listeners)
    ports = run("docker", "ps", "--filter", "label=com.docker.compose.project=knowslink", "--format", "{{.Names}} {{.Ports}}")
    assert "5432->" not in ports and ports.count("->") == 1, ports
    for name in ("postgres", "relay"):
        inspect = json.loads(run("docker", "inspect", f"knowslink-{name}-1"))[0]
        assert inspect["HostConfig"]["RestartPolicy"]["Name"] == "unless-stopped" and inspect["HostConfig"]["Memory"] > 0
    for path in (".env", "tunnel"):
        assert os.stat(os.path.join(STATE, path)).st_mode & 0o077 == 0, path
    print("local checks ok: health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700")


def public():
    probes = [("/healthz", {}), ("/owner", {}), ("/v1/registry", {}), ("/healthz", {"Cf-Access-Jwt-Assertion": "x.y.z"}),
              ("/healthz", {"CF-Access-Client-Id": "x.access", "CF-Access-Client-Secret": "x"})]
    for path, headers in probes:
        status, location, body = fetch("https", HOST, path, headers)
        assert status != 200 and b"relay" not in body.lower(), (path, status)
        assert status in (302, 401, 403), (path, status)
        if status == 302:
            assert urllib.parse.urlparse(location).hostname.endswith("cloudflareaccess.com"), location
        print("unauthenticated", path, bool(headers), "->", status)


if __name__ == "__main__":
    {"baseline": baseline, "regression": regression, "local": local, "public": public}[sys.argv[1] if len(sys.argv) > 1 else ""]()
