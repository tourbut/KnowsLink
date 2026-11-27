"""Run a make target for FullOps lint; on Windows add the repo-local tools to this process only."""

import os
import subprocess
import sys
from pathlib import Path

command = ["make", *sys.argv[1:]]
environment = dict(os.environ)
if os.name == "nt":
    def git(*args):
        return Path(subprocess.check_output(["git", *args], cwd=Path(__file__).parent, text=True).strip())

    tools = git("rev-parse", "--path-format=absolute", "--git-common-dir").parent / ".devtools"
    # Git for Windows supplies sh, bash and coreutils: <root>/mingw64/libexec/git-core -> <root>/usr/bin.
    shell = git("--exec-path").parents[2] / "usr/bin"
    node, w64 = tools / "node-v22.22.2-win-x64", tools / "w64devkit/bin"
    for needed in (node / "node.exe", shell / "sh.exe", w64 / "make.exe", w64 / "gcc.exe"):
        if not needed.is_file():
            raise SystemExit(f"missing Windows build tool: {needed}")
    # w64devkit goes last (make, gcc for -race) so its busybox does not shadow Git's shell tools.
    environment["PATH"] = os.pathsep.join([str(node), str(shell), environment["PATH"], str(w64)])
    command[0] = str(w64 / "make.exe")
sys.exit(subprocess.run(command, env=environment).returncode)
