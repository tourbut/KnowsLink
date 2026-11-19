"""Run make lint and make test on the clean final-SHA Linux clone. No secret output."""
import subprocess
import sys
from pathlib import Path

base = Path(sys.argv[1]).resolve()
sha = sys.argv[2]
repo = base / "repo"
assert subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=repo, text=True).strip() == sha
assert not subprocess.check_output(["git", "status", "--porcelain"], cwd=repo, text=True).strip()
results = []
for args in (["make", "lint"], ["make", "test"]):
    done = subprocess.run(args, cwd=repo, capture_output=True, text=True)
    tail = "\n".join((done.stdout + done.stderr).splitlines()[-30:])
    results.append({"run": args, "exit": done.returncode, "tail": tail})
    print(f"COMMAND={' '.join(args)} EXIT={done.returncode}")
    print(tail)
sys.exit(0 if all(item["exit"] == 0 for item in results) else 1)
