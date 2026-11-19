#!/bin/sh
cd C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA
export QA_RESULTS=results-run3.jsonl
sed -i 's#work/run3#work/run4#' lib.py; sed -i 's#work-c3#work-c4#g' s_folders.py
for s in "s_folders.py" "s_names.py" "s_nohost.py" "s_browser.py main" "s_reasons.py" "s_security.py" "s_browser.py signup"; do
  echo "=== $s $(date -u +%H:%M:%S)"
  python $s 2>&1
  echo "=== exit $? $s"
done
echo ALLDONE
