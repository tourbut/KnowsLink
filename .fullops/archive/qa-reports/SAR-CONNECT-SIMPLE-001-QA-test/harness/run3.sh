#!/bin/sh
cd C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA
export QA_RESULTS=results-run3.jsonl
sed -i 's#work/run4#work/run5#' lib.py
for s in "s_names.py" "s_browser.py main" "s_security.py" "s_browser.py signup"; do
  echo "=== $s $(date -u +%H:%M:%S)"
  python $s 2>&1
  echo "=== exit $? $s"
done
echo ALLDONE
