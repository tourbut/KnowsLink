#!/bin/bash
# usage: launch.sh <title> <workdir-name> <plugin-letter>  -> prints terminal handle, accepts folder trust
export MSYS_NO_PATHCONV=1
cd /c/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-002-qa
BS='\'
Q="C:${BS}Users${BS}shin${BS}AppData${BS}Local${BS}KnowsLinkDevTools${BS}session-peer-002-qa"
orca terminal create --worktree path:C:/Users/shin/orca/workspaces/KnowsLink/tester --title "$1" --shell pwsh.exe --command "Set-Location '${Q}${BS}$2'; claude --plugin-dir '${Q}${BS}inst-$3${BS}knowslink'" --json > term-$1.json
H=$(python3 -c "import json;print(json.load(open('term-$1.json'))['result']['terminal']['handle'])")
sleep 6
orca terminal read --terminal $H --screen | grep -q "trust this folder" && { orca terminal send --terminal $H --text $'\e[B' >/dev/null; sleep 1; orca terminal send --terminal $H --text "" --enter >/dev/null; }
sleep 6
echo $H
