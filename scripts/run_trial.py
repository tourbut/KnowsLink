"""Load a private trial environment without exposing credentials in shell arguments; launch Codex CLI or Grok MCP."""
import argparse
import json
import os
from pathlib import Path
import stat


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--config", required=True)
    parser.add_argument("--node", default="node")
    parser.add_argument("--bundle", help="Grok standalone MCP bundle; omit for Codex trial CLI")
    parser.add_argument("arguments", nargs="*")
    args = parser.parse_args()
    path = Path(args.config)
    info = path.lstat()
    if not stat.S_ISREG(info.st_mode) or info.st_mode & 0o077 or info.st_uid != os.getuid():
        raise ValueError("trial config must be an owned private regular file")
    config = json.loads(path.read_text())
    allowed = {"KNOWSLINK_MODE", "RELAY_URL", "AGENT_ID", "AGENT_CREDENTIAL", "AGENT_KID", "AGENT_KEY_FILE", "KNOWSLINK_TEST_PEER", "CF_ACCESS_CLIENT_ID", "CF_ACCESS_CLIENT_SECRET"}
    if not isinstance(config, dict) or not set(config) <= allowed or not all(isinstance(value, str) for value in config.values()):
        raise ValueError("invalid trial configuration")
    # Do not inherit another agent's relay credentials or another mode's leftover settings.
    environment = {key: value for key, value in os.environ.items() if key not in allowed}
    environment.update(config)
    bundle = args.bundle or str(Path(__file__).resolve().parents[1] / "adapters/dist/trial-cli.js")
    os.execvpe(args.node, [args.node, bundle, *args.arguments], environment)


if __name__ == "__main__":
    try:
        main()
    except (OSError, ValueError):
        raise SystemExit("trial launch failed: check private config and executable paths")
