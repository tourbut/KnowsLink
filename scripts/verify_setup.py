"""Prove product lint detects violations in temporary copies and propagates failure."""

import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]


def run(command, directory, expected=0, diagnostic=None):
    result = subprocess.run(command, cwd=directory, text=True, capture_output=True)
    print(f"command: {command}; exit: {result.returncode}", flush=True)
    print(result.stdout + result.stderr, end="", flush=True)
    if (expected == 0 and result.returncode != 0) or (expected != 0 and result.returncode == 0):
        raise AssertionError("unexpected command exit code")
    if diagnostic and diagnostic not in result.stdout + result.stderr:
        raise AssertionError(f"expected diagnostic missing: {diagnostic}")


def main():
    run(["make", "lint"], ROOT)
    registered = json.loads((ROOT / ".fullops-squad/lint/lint.json").read_text())["commands"]
    product = next(item for item in registered if item["name"] == "product-lint")
    assert product["run"] == ["make", "lint"] and product["cwd"] == "."
    with tempfile.TemporaryDirectory(prefix="knowslink-lint-") as temporary:
        copy = Path(temporary) / "repo"
        shutil.copytree(ROOT, copy, ignore=shutil.ignore_patterns(
            ".git", ".env", ".env.*", "node_modules", "dist", "build", "__pycache__",
        ))
        shutil.copyfile(ROOT / ".env.example", copy / ".env.example")
        os.symlink(ROOT / "adapters/node_modules", copy / "adapters/node_modules")
        cases = [
            ("cmd/relay/lint_probe.go", 'package main\nfunc lintProbe( ){ }\n', "Go formatting violations"),
            ("adapters/src/lint_probe.ts", 'const unusedProbe = 1;\n', "@typescript-eslint/no-unused-vars"),
            ("adapters/src/lint_probe.ts", 'export const probe: string = 1;\n', "TS2322"),
            ("adapters/src/lint_probe.ts", 'export const probe="bad format"\n', "Code style issues"),
        ]
        for relative, content, diagnostic in cases:
            path = copy / relative
            path.write_text(content)
            try:
                run(product["run"], copy, expected=1, diagnostic=diagnostic)
            finally:
                path.unlink()
        run(product["run"], copy)
    print("PASS: format, lint, type errors fail registered integrated command; restored copy passes")


if __name__ == "__main__":
    main()
