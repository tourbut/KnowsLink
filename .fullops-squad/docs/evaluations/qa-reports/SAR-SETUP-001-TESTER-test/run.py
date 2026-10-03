#!/usr/bin/env python3
# usage: run.py LOG CWD cmd...  runs cmd (no shell), appends output and real exit code to LOG
import subprocess, sys, time
log, cwd, *cmd = sys.argv[1:]
t = time.strftime("%Y-%m-%dT%H:%M:%S%z")
p = subprocess.run(cmd, cwd=cwd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
out = "\n".join(l.rstrip() for l in p.stdout.splitlines())
with open(log, "a") as f:
    f.write(f"## {t} cwd={cwd}\n$ {' '.join(cmd)}\n{out}\n[exit {p.returncode}]\n\n")
print(f"exit={p.returncode} :: {' '.join(cmd)}")
