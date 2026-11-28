#!/bin/sh
# 3b. Download failure (no network) keeps the existing install left by qa.sh (fake v22.0.0).
set -u
P=/workspace/.knowslink
env -i PATH=/usr/bin:/bin sh /plugin/knowslink/bin/knowslink-grok.sh --prepare-node </dev/null >/dev/null 2>/tmp/e; rc=$?
[ $rc -ne 0 ] || { echo "FAIL 3b: exit 0"; exit 1; }
[ "$("$P/node/bin/node" --version)" = v22.0.0 ] && [ -z "$(ls -A "$P" | grep '^\.stage' || true)" ] || { echo "FAIL 3b: changed/stage"; exit 1; }
echo "PASS 3b download failure (network none): exit $rc, '$(tail -1 /tmp/e)', existing node kept, no stage left"
