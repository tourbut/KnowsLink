import hashlib, os, sys
from pathlib import Path
h = Path.home(); prof = Path(os.environ.get("CLAUDE_CONFIG_DIR", h / ".claude"))
files = [h/".claude/settings.json", h/".claude.json", h/".claude/plugins/installed_plugins.json", h/".claude/plugins/known_marketplaces.json",
         h/".codex/config.toml", h/".codex/auth.json", h/".codex/hooks.json", h/".codex/AGENTS.md",
         prof/"settings.json", prof/".claude.json", prof/"plugins/installed_plugins.json", prof/"plugins/known_marketplaces.json"]
for f in files:
    if f.is_file():
        b = f.read_bytes()
        print(f"{hashlib.sha256(b).hexdigest()[:16]} {len(b):>8} knowslink-mentions={b.lower().count(b'knowslink')} {str(f).replace(str(h),'~')}")
    else:
        print(f"absent {f}")
