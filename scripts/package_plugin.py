"""Package a standalone Cursor/Grok Bot plugin with a strict file allowlist and no credentials."""

from hashlib import sha256
import argparse
import json
from pathlib import Path
import subprocess
from tempfile import TemporaryDirectory
from zipfile import ZipFile, ZipInfo, ZIP_DEFLATED

ROOT = Path(__file__).resolve().parents[1]
FILES = [".cursor-plugin/plugin.json", "mcp.json", "README.md", "skills/knowslink/SKILL.md", "dist/plugin.js"]


def add(archive, name, content):
    entry = ZipInfo(name, date_time=(1980, 1, 1, 0, 0, 0))
    entry.compress_type = ZIP_DEFLATED
    entry.external_attr = 0o100644 << 16
    archive.writestr(entry, content)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--verify", action="store_true", help="Run MCP checks against the extracted standalone bundle")
    options = parser.parse_args()
    sources = [(name, ROOT / "adapters" / name) for name in FILES]
    for name, path in sources:
        if not path.is_file() or path.is_symlink():
            raise RuntimeError(f"Missing regular plugin file: {name}; run make build first")
    output = ROOT / "build/knowslink-grok-bot-plugin.zip"
    output.parent.mkdir(exist_ok=True)
    with ZipFile(output, "w", ZIP_DEFLATED) as archive:
        marketplace = {"name": "knowslink-plugins", "owner": {"name": "KnowsLink"}, "plugins": [{"name": "knowslink", "source": "knowslink", "description": "Human-gated pull relay connector; actual connection held"}]}
        add(archive, ".cursor-plugin/marketplace.json", json.dumps(marketplace, indent=2) + "\n")
        for name, path in sources:
            add(archive, "knowslink/" + name, path.read_bytes())
        add(archive, "knowslink/package.json", json.dumps({"name": "knowslink", "version": "0.1.0", "private": True, "type": "module", "engines": {"node": ">=22.22.2 <23"}}) + "\n")
        notices = []
        packages = json.loads((ROOT / "adapters/package-lock.json").read_text())["packages"]
        for name, metadata in sorted(packages.items()):
            if not name or metadata.get("dev"):
                continue
            directory = ROOT / "adapters" / name
            licenses = sorted(path for path in directory.iterdir() if path.is_file() and path.name.lower().startswith(("license", "copying")))
            if not licenses:
                raise RuntimeError(f"Missing third-party license: {name}")
            notices.append(f"{name} {metadata['version']}\n" + "\n".join(path.read_text() for path in licenses))
        add(archive, "knowslink/THIRD_PARTY_NOTICES.txt", "\n\n".join(notices))
    print(f"Plugin: {output.relative_to(ROOT)}; sha256: {sha256(output.read_bytes()).hexdigest()}")
    if options.verify:
        with TemporaryDirectory(prefix="knowslink-plugin-test-") as temporary:
            with ZipFile(output) as archive:
                archive.extractall(temporary)
            command = ["node", "adapters/dist/mcp.test.js", str(Path(temporary) / "knowslink/dist/plugin.js")]
            result = subprocess.run(command, cwd=ROOT)
            print(f"extracted standalone MCP test exit={result.returncode}", flush=True)
            if result.returncode:
                raise SystemExit(result.returncode)


if __name__ == "__main__":
    main()
