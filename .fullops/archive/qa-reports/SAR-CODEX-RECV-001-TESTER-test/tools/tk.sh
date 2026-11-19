# tk.sh <prompt> : send a prompt to the TUI terminal, auto-pick "Allow" on Codex's own tool dialog, then print the screen
ORCA="C:/Users/shin/AppData/Local/Programs/orca/resources/bin/orca.EXE"; D=/c/Users/shin/AppData/Local/KnowsLinkDevTools/codex-recv-001-tester; T=$(cat $D/term-Y.txt)
export MSYS_NO_PATHCONV=1
"$ORCA" terminal send --terminal $T --text "$1" --enter --json >/dev/null
for i in $(seq 1 12); do sleep 3; S=$("$ORCA" terminal read --terminal $T 2>&1 | tail -n 16); if echo "$S" | grep -q "enter to submit | esc to cancel" && echo "$S" | grep -q "Allow the knowslink"; then "$ORCA" terminal send --terminal $T --text '' --enter --json >/dev/null; sleep 4; break; fi; done
"$ORCA" terminal read --terminal $T 2>&1 | tail -n 16
