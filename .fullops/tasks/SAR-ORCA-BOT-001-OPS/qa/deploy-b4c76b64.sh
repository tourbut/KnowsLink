#!/usr/bin/env bash
# SAR-ORCA-BOT-001-OPS: deploy b4c76b64 on knowslink-relay-01 as shin.
# Upload and run as a file: piping into `ssh ... bash -s` lets `docker compose exec` eat the rest of the script.
set -o errexit -o nounset -o pipefail
NEW=b4c76b646f315178852ce56650685ca64a4a343e
OLD=40aafe6ff396db2681a9423c3b344732934756e0
D=/home/shin/deploy/knowslink
STATE=/home/shin/deploy/knowslink-state

[ "$(git -C "$D" rev-parse HEAD)" = "$OLD" ] || { echo "ERROR: HEAD is not $OLD" >&2; exit 1; }
[ -z "$(git -C "$D" status --porcelain --untracked-files=no)" ] || { echo "ERROR: checkout is dirty" >&2; exit 1; }

# Rollback tags, taken before the build replaces the images.
for name in relay migrate; do
  tag="knowslink-rollback:40aafe6f-$name-orcabot001"
  docker tag "knowslink-$name" "$tag"
  docker image inspect -f "$tag {{.Id}}" "$tag"
done

# beta.sh deploy: refuses migration changes, takes the backup, checks out, builds, runs verify.py local.
bash "$D/deploy/knowslink/beta.sh" deploy "$NEW" </dev/null

dump=$(ls -t "$STATE"/backups/*.dump | sed -n 1p)
stat -c 'backup mode=%a bytes=%s %n' "$dump"
sha256sum "$dump"
python3 "$D/deploy/knowslink/verify.py" public
echo "relay_error_lines=$(docker logs knowslink-relay-1 2>&1 | grep -ci error || true)"
echo "head=$(git -C "$D" rev-parse HEAD)"
docker ps --filter label=com.docker.compose.project=knowslink --format '{{.Names}} {{.Status}}'
echo "rollback if needed: bash $D/deploy/knowslink/beta.sh deploy $OLD"
