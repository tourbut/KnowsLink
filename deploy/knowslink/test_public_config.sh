#!/usr/bin/env bash
# Check render-only public configuration, private file mode and failure propagation without live deployment or Cloudflare calls.
set -euo pipefail
root=$(git rev-parse --show-toplevel)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT
mkdir -p "$test_dir/state/tunnel" "$test_dir/bin"
export KNOWSLINK_DEPLOY="$root" KNOWSLINK_STATE_DIR="$test_dir/state"
export PATH="$test_dir/bin:$PATH"
cat >"$test_dir/bin/cloudflared" <<'SH'
#!/bin/sh
[ "$1 $2 $4 $5" = 'tunnel --config ingress validate' ] || exit 90
exit "${VALIDATE_EXIT:-0}"
SH
chmod +x "$test_dir/bin/cloudflared"
printf '%s\n' '00000000-0000-4000-8000-000000000001' >"$KNOWSLINK_STATE_DIR/tunnel.uuid"
printf 'unchanged-live-config\n' >"$KNOWSLINK_STATE_DIR/tunnel/config.yml"
bash "$root/deploy/knowslink/beta.sh" render-public-config
candidate="$KNOWSLINK_STATE_DIR/tunnel/config.public.yml"
[ "$(stat -c %a "$candidate")" = 600 ]
python3 - "$candidate" "$root" <<'PY'
from pathlib import Path
import sys
candidate = Path(sys.argv[1]).read_text()
fragment = (Path(sys.argv[2]) / 'deploy/knowslink/tunnel/public-ingress.yml').read_text()
assert candidate.endswith('  - service: http_status:404\n')
assert ''.join('  ' + line + '\n' for line in fragment.splitlines()) in candidate
assert candidate.count('service: http://relay:8080') == 1
assert 'required: true' not in candidate and 'audTag:' not in candidate
PY
before=$(sha256sum "$candidate")
if bash "$root/deploy/knowslink/beta.sh" render-public-config; then exit 1; fi
[ "$(sha256sum "$candidate")" = "$before" ]
rm "$candidate"
if VALIDATE_EXIT=7 bash "$root/deploy/knowslink/beta.sh" render-public-config; then exit 1; fi
rm "$candidate"
printf 'invalid\n' >"$KNOWSLINK_STATE_DIR/tunnel.uuid"
if bash "$root/deploy/knowslink/beta.sh" render-public-config; then exit 1; fi
[ ! -e "$candidate" ]
[ "$(cat "$KNOWSLINK_STATE_DIR/tunnel/config.yml")" = unchanged-live-config ]
echo 'Public config: render-only, private, default deny, no overwrite, validation failures propagated'
