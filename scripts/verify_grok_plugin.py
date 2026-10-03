"""Install the built plugin ZIP with the real Grok CLI in a throwaway HOME and check MCP discovery; ~/.grok is untouched."""

import json
import os
from pathlib import Path
import subprocess
from tempfile import TemporaryDirectory
from zipfile import ZipFile

ROOT = Path(__file__).resolve().parents[1]


def main():
    with TemporaryDirectory(prefix="knowslink-grok-") as temporary:
        home, package = Path(temporary) / "home", Path(temporary) / "package"
        home.mkdir()
        with ZipFile(ROOT / "build/knowslink-grok-bot-plugin.zip") as archive:
            archive.extractall(package)
        environment = {**os.environ, "HOME": str(home)}
        environment.pop("GROK_HOME", None)

        def grok(*arguments):
            command = ["grok", *arguments]
            return subprocess.run(command, cwd=temporary, env=environment, check=True, capture_output=True, text=True, timeout=120).stdout

        grok("plugin", "validate", str(package / "knowslink"))
        grok("plugin", "install", str(package / "knowslink"), "--trust")
        plugins = json.loads(grok("plugin", "list", "--json"))
        assert [(plugin["name"], plugin["version"]) for plugin in plugins] == [("knowslink", "0.1.0")], plugins
        doctor = json.loads(grok("mcp", "doctor", "knowslink", "--json"))
        checks = {check["label"]: check["passed"] for server in doctor["servers"] for check in server["checks"]}
        assert doctor["healthy_count"] == 1 and checks.get("2 tools discovered"), doctor
        installed = Path(plugins[0]["path"]) / "dist/plugin.js"
        subprocess.run(["node", "adapters/dist/mcp.test.js", str(installed)], cwd=ROOT, check=True)
    print("PASS: grok plugin validate/install, knowslink 0.1.0 listed, mcp doctor healthy with 2 tools, installed copy held", flush=True)


if __name__ == "__main__":
    main()
