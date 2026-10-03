#!/usr/bin/env python3
"""External negative HTTP for link.knowslog.com. Saves status and location host only."""

import json
import socket
import subprocess
import tempfile
import urllib.parse
from pathlib import Path

HOST = "link.knowslog.com"
EVIDENCE = Path(__file__).resolve().parent
PATHS = ["/", "/owner", "/v1/registry", "/healthz", "/_probe"]
HEADER_CASES = [
    ("none", []),
    ("fake-jwt", ["-H", "Cf-Access-Jwt-Assertion: x.y.z"]),
    ("service-token", ["-H", "CF-Access-Client-Id: x.access", "-H", "CF-Access-Client-Secret: x"]),
    ("bearer", ["-H", "Authorization: Bearer x"]),
]


def dig(server):
    command = ["dig", "+time=3", "+tries=1", "+short", HOST, "A"]
    if server:
        command[1:1] = [f"@{server}"]
    completed = subprocess.run(command, text=True, capture_output=True)
    ips = [line.strip() for line in completed.stdout.splitlines() if line.strip() and line[0].isdigit()]
    return completed.returncode, ips


def system_resolve():
    try:
        socket.getaddrinfo(HOST, 443)
    except socket.gaierror as exc:
        return {"ok": False, "error": exc.__class__.__name__}
    return {"ok": True, "error": ""}


def curl(ip, port, scheme, path, extra):
    url = f"{scheme}://{HOST}{path}"
    with tempfile.TemporaryDirectory(prefix="knowslink-public-body-") as temporary:
        body_path = Path(temporary) / "body"
        command = [
            "curl", "-sS", "--http1.1", "--max-redirs", "0", "--connect-timeout", "10", "--max-time", "20",
            "--resolve", f"{HOST}:{port}:{ip}", "-D", "-", "-o", str(body_path), *extra, url,
        ]
        completed = subprocess.run(command, text=True, capture_output=True)
        body = body_path.read_bytes() if body_path.exists() else b""
    status = None
    location = ""
    for line in completed.stdout.splitlines():
        if line.startswith("HTTP/") and status is None:
            parts = line.split()
            status = int(parts[1]) if len(parts) > 1 and parts[1].isdigit() else None
        if line.lower().startswith("location:"):
            location = line.split(":", 1)[1].strip()
    parsed = urllib.parse.urlparse(location) if location else None
    location_host = parsed.hostname or "" if parsed else ""
    return {
        "exit": completed.returncode,
        "status": status,
        "location_host": location_host,
        "location_scheme": parsed.scheme if parsed else "",
        "body_has_relay": b"relay" in body.lower(),
        "body_len": len(body),
        "curl_error": completed.returncode != 0,
    }


def main():
    public = {}
    for server in ("1.1.1.1", "8.8.8.8"):
        code, ips = dig(server)
        public[server] = {"exit": code, "ips": ips}
    system_code, system_ips = dig(None)
    ips = public["1.1.1.1"]["ips"]
    same = public["1.1.1.1"]["ips"] == public["8.8.8.8"]["ips"] and bool(ips)
    ip = ips[0]
    probes = []
    for path in PATHS:
        item = curl(ip, 443, "https", path, [])
        item.update({"path": path, "headers": "none", "port": 443})
        probes.append(item)
    for label, extra in HEADER_CASES:
        if label == "none":
            continue
        item = curl(ip, 443, "https", "/healthz", extra)
        item.update({"path": "/healthz", "headers": label, "port": 443})
        probes.append(item)
    http80 = curl(ip, 80, "http", "/", [])
    http80.update({"path": "/", "headers": "none", "port": 80})
    access_ok = []
    for item in probes:
        ok = (
            item["exit"] == 0 and item["status"] in (302, 401, 403) and item["status"] != 200
            and not item["body_has_relay"]
            and (item["status"] != 302 or item["location_host"].endswith("cloudflareaccess.com"))
        )
        item["pass"] = ok
        access_ok.append(ok)
    http80["pass"] = (
        http80["exit"] == 0 and http80["status"] == 301 and http80["location_scheme"] == "https"
        and http80["location_host"] == HOST and not http80["body_has_relay"]
    )
    result = {
        "system_resolver": system_resolve(),
        "system_dig_exit": system_code,
        "system_dig_empty": system_ips == [],
        "public_ips_agree": same,
        "edge_ip_count": len(ips),
        "probes": probes,
        "http80": http80,
    }
    (EVIDENCE / "public.json").write_text(json.dumps(result, indent=2) + "\n")
    failed = [f"{item['path']}:{item['headers']}" for item in probes if not item["pass"]]
    if not http80["pass"]:
        failed.append("http80")
    if not same or result["system_resolver"]["ok"] or not result["system_dig_empty"]:
        failed.append("dns-distinction")
    print(json.dumps({"failed": failed, "statuses": [
        {"path": item["path"], "headers": item["headers"], "status": item["status"], "location_host": item["location_host"]}
        for item in probes
    ], "http80": {"status": http80["status"], "location_host": http80["location_host"], "scheme": http80["location_scheme"]}}, indent=2))
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
