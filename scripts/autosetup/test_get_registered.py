"""Run test-get-registered.ps1 under every PowerShell found (Windows PowerShell 5.1 and PowerShell 7); skip when none exists."""
from pathlib import Path
import shutil
import subprocess
import sys

script = Path(__file__).resolve().with_name("test-get-registered.ps1")
shells = [shell for shell in ("powershell.exe", "pwsh") if shutil.which(shell)]
if not shells:
    print("SKIP: no PowerShell on PATH")
for shell in shells:
    done = subprocess.run([shell, "-NoProfile", "-File", str(script)], capture_output=True, text=True, encoding="utf-8", errors="replace")
    print(f"== {shell}\n{done.stdout}{done.stderr}", end="")
    if done.returncode:
        sys.exit(f"{shell} exited {done.returncode}")
