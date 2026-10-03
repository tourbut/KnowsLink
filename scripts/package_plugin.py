"""Package a standalone Cursor/Grok Bot plugin with a strict file allowlist and no credentials."""

from hashlib import sha256
import json
from pathlib import Path
from zipfile import ZipFile, ZipInfo, ZIP_DEFLATED

ROOT = Path(__file__).resolve().parents[1]
FILES = [".cursor-plugin/plugin.json", "mcp.json", "README.md", "skills/knowslink/SKILL.md", "dist/plugin.js"]


def add(archive, name, content):
    entry = ZipInfo(name, date_time=(1980, 1, 1, 0, 0, 0))
    entry.compress_type = ZIP_DEFLATED
    entry.external_attr = 0o100644 << 16
    archive.writestr(entry, content)


def main():
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
    print(f"Plugin: {output.relative_to(ROOT)}; sha256: {sha256(output.read_bytes()).hexdigest()}")


if __name__ == "__main__":
    main()
