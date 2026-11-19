#!/usr/bin/env python3
"""Run narrow timeout probes against a detached candidate checkout and keep each exit code."""
import json
import subprocess
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
LOGS = HERE / "logs"
BASE = "cd60e7f87eb5ce137eca887980f232b3f67a18d0"
FIXED = "711f2532be423d1ca7707463a20fdc168f50bece"
UNCHANGED = [
    "cmd",
    "internal",
    "db",
    "scripts",
    "Makefile",
    "adapters/package.json",
    "adapters/package-lock.json",
    "adapters/src/mcp.ts",
    "adapters/src/test-transport.ts",
]


def run(name, args, cwd):
    LOGS.mkdir(parents=True, exist_ok=True)
    out_path = LOGS / f"{name}.out"
    err_path = LOGS / f"{name}.err"
    with out_path.open("w", encoding="utf-8") as out, err_path.open("w", encoding="utf-8") as err:
        completed = subprocess.run(args, cwd=cwd, stdout=out, stderr=err)
    (LOGS / f"{name}.exit").write_text(f"{completed.returncode}\n", encoding="utf-8")
    return completed.returncode


def main():
    if len(sys.argv) != 2:
        sys.exit("usage: run_narrow.py <detached-clone>")
    clone = Path(sys.argv[1]).resolve()
    node = run("node-version", ["node", "-p", "process.version"], clone)
    head = run("clone-head", ["git", "rev-parse", "HEAD"], clone)
    status = run("clone-status", ["git", "status", "--porcelain"], clone)
    names = run(
        "product-names",
        ["git", "diff", "--name-status", BASE, "HEAD", "--", ".", ":(exclude).fullops-squad"],
        clone,
    )
    unchanged = run("unchanged-diff", ["git", "diff", "--exit-code", BASE, "HEAD", "--", *UNCHANGED], clone)
    dist = str(clone / "adapters" / "dist")
    root = str(clone / "adapters")
    probes = [
        ("short-100ms", ["node", str(HERE / "short-100ms.mjs"), dist]),
        ("held-network", ["node", str(HERE / "held-network.mjs"), root]),
        ("default-1", ["node", str(HERE / "default-stall.mjs"), dist, "10500", "1"]),
        ("default-2", ["node", str(HERE / "default-stall.mjs"), dist, "10500", "2"]),
        ("default-3", ["node", str(HERE / "default-stall.mjs"), dist, "10500", "3"]),
        ("gc-default", ["node", "--expose-gc", str(HERE / "gc-stall.mjs"), dist, "10500"]),
        ("cleanup", ["node", str(HERE / "cleanup.mjs"), dist, "10500"]),
        ("mcp-busy", ["node", str(HERE / "mcp-busy.mjs"), root]),
    ]
    results = {
        "node-version": node,
        "clone-head": head,
        "clone-status": status,
        "product-names": names,
        "unchanged-diff": unchanged,
    }
    for name, args in probes:
        results[name] = run(name, args, clone)
    head_text = (LOGS / "clone-head.out").read_text(encoding="utf-8").strip()
    status_text = (LOGS / "clone-status.out").read_text(encoding="utf-8").strip()
    node_text = (LOGS / "node-version.out").read_text(encoding="utf-8").strip()
    identity_ok = head_text == FIXED and status_text == "" and node_text == "v22.22.2"
    summary = {
        "clone": str(clone),
        "fixed": FIXED,
        "head": head_text,
        "node": node_text,
        "clean": status_text == "",
        "results": results,
        "ok": identity_ok and all(code == 0 for code in results.values()),
    }
    (LOGS / "summary.json").write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(summary))
    return 0 if summary["ok"] else 1


if __name__ == "__main__":
    sys.exit(main())
