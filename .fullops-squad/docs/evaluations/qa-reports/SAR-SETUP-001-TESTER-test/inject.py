#!/usr/bin/env python3
# usage: inject.py DIR KIND  — applies one temporary violation to the throwaway copy
import sys
d, kind = sys.argv[1:]
def app(p, s):
    open(f"{d}/{p}", "a").write(s)
if kind == "go-format":
    app("internal/config/config.go", "\nfunc   badFormat( ){ }\n")
elif kind == "ts-lint":
    app("adapters/src/index.ts", '\nconst unusedValue = 1;\n')
elif kind == "ts-type":
    app("adapters/src/index.ts", '\nconst typeMismatch: number = "text";\nconsole.info(typeMismatch);\n')
elif kind == "ts-format":
    app("adapters/src/index.ts", "\nconsole.info(   'x'  )\n")
else:
    sys.exit("unknown kind")
