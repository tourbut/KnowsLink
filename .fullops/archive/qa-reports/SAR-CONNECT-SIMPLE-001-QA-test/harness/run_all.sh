#!/bin/sh
cd C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA
export QA_RESULTS=results-run3.jsonl
for s in s_limits s_attempts s_cleanup s_folders s_names "s_browser main" s_grok s_reasons s_security "s_browser signup"; do
  echo "=== $s $(date -u +%H:%M:%S)"
  python $s.py 2>&1
  echo "=== exit $? $s"
done
echo ALLDONE
