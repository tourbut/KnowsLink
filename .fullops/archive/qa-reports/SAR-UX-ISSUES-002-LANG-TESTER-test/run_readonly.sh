#!/bin/bash
# Read-only browser stages. Each stage keeps its own log and exit code (no pipes that hide the exit status).
cd "$(dirname "$0")"
run() { local name="$1"; shift; node browser.mjs "$@" > "$name.log" 2>&1; echo "exit=$?" >> "$name.log"; }
run chrome-layout-en 9333 chrome layout en
run chrome-keys-en 9333 chrome keys en
run chrome-keys-ko 9333 chrome keys ko
run chrome-langflow 9333 chrome langflow
run chrome-zoom 9333 chrome zoom
run edge-layout-en 9334 edge layout en
run edge-layout-ko 9334 edge layout ko
run edge-keys-en 9334 edge keys en
run edge-keys-ko 9334 edge keys ko
run edge-langflow 9334 edge langflow
echo done > run_readonly.done
