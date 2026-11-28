#!/usr/bin/env bash
# SAR-EVENTS-DEPLOY-001-OPS: backup, rollback tags, add KNOWSLINK_SEAL_KEY, deploy 65a752fe on knowslink-relay-01 as shin.
# Upload and run as a file: piping into `ssh ... bash -s` lets `docker compose exec` eat the rest of the script.
# Never add `set -x`: the key must not reach output, logs or argv.
set -o errexit -o nounset -o pipefail
NEW=65a752fe2e88e0da1cbf1a844cb0000c51595295
OLD=d3670f2fa3c85a8d320542055188296948aa783b
D=/home/shin/deploy/knowslink
STATE=/home/shin/deploy/knowslink-state
ENVFILE=$STATE/.env

[ "$(git -C "$D" rev-parse HEAD)" = "$OLD" ] || { echo "ERROR: HEAD is not $OLD" >&2; exit 1; }
[ -z "$(git -C "$D" status --porcelain --untracked-files=no)" ] || { echo "ERROR: checkout is dirty" >&2; exit 1; }
if grep -q 'KNOWSLINK_SEAL_KEY' "$ENVFILE"; then echo "ERROR: KNOWSLINK_SEAL_KEY already in env; not overwriting" >&2; exit 1; fi
git -C "$D" fetch -q origin || echo "origin fetch failed; using commits already in $D" >&2
git -C "$D" cat-file -e "$NEW^{commit}" || { echo "ERROR: $NEW not in $D" >&2; exit 1; }
git -C "$D" diff --quiet "$OLD" "$NEW" -- db Dockerfile tunnel || { echo "ERROR: db/Dockerfile/tunnel changed" >&2; exit 1; }

# 1. DB backup, rollback tags (before the build replaces the images), env backup.
bash "$D/deploy/knowslink/beta.sh" backup </dev/null
dump=$(ls -t "$STATE"/backups/*.dump | sed -n 1p)
stat -c 'backup mode=%a bytes=%s %n' "$dump"
sha256sum "$dump"
for name in relay migrate; do
  tag="knowslink-rollback:d3670f2f-$name-eventsdeploy001"
  docker tag "knowslink-$name" "$tag"
  docker image inspect -f "$tag {{.Id}}" "$tag"
done
umask 077
keep="$STATE/eventsdeploy001-$(date -u +%Y%m%d)"
install -d -m 0700 "$keep"
[ ! -e "$keep/env.pre-eventsdeploy001" ] || { echo "ERROR: env backup exists" >&2; exit 1; }
cp -p "$ENVFILE" "$keep/env.pre-eventsdeploy001"
stat -c 'env backup mode=%a bytes=%s %n' "$keep/env.pre-eventsdeploy001"
sha256sum "$keep/env.pre-eventsdeploy001"

# 2. Key. Held in a variable and appended with the printf builtin, so it never reaches argv or output.
key=$(openssl rand -base64 32)
[ "${#key}" -eq 44 ] || { echo "ERROR: unexpected key length ${#key}" >&2; exit 1; }
[ -z "$(tail -c1 "$ENVFILE")" ] || printf '\n' >>"$ENVFILE"
printf 'KNOWSLINK_SEAL_KEY=%s\n' "$key" >>"$ENVFILE"
unset key
echo "seal_key_lines=$(grep -c '^KNOWSLINK_SEAL_KEY=' "$ENVFILE")"
awk -F= '$1=="KNOWSLINK_SEAL_KEY"{print "seal_key_value_length=" length(substr($0, index($0,"=")+1))}' "$ENVFILE"
# Everything above the new line is byte for byte the backup.
head -n -1 "$ENVFILE" | cmp - "$keep/env.pre-eventsdeploy001" && echo "env_otherwise_unchanged=yes"
stat -c 'env mode=%a bytes=%s' "$ENVFILE"

# 3. beta.sh deploy: refuses migration changes, backs up again, checks out, builds, runs verify.py local (includes redis).
bash "$D/deploy/knowslink/beta.sh" deploy "$NEW" </dev/null

# 4. Verify.
python3 "$D/deploy/knowslink/verify.py" public
python3 "$D/deploy/knowslink/verify.py" redis
docker logs knowslink-relay-1 2>&1 | grep 'relay starting' | tail -n 1
echo "relay_error_lines=$(docker logs knowslink-relay-1 2>&1 | grep -ci error || true)"
echo "relay_restarts=$(docker inspect -f '{{.RestartCount}}' knowslink-relay-1)"
echo "head=$(git -C "$D" rev-parse HEAD)"
docker ps --filter label=com.docker.compose.project=knowslink --format '{{.Names}} {{.Status}}'
echo "rollback if needed: bash $D/deploy/knowslink/beta.sh deploy $OLD"
