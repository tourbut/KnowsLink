#!/usr/bin/env bash
# Owner-only synthetic beta lifecycle for the separate "knowslink" Compose project; never touches other projects.
set -euo pipefail

SRC=${KNOWSLINK_SRC:-$(git -C "$(dirname "$0")" rev-parse --show-toplevel)}
DEPLOY=${KNOWSLINK_DEPLOY:-/home/shin/deploy/knowslink}
STATE=${KNOWSLINK_STATE_DIR:-/home/shin/deploy/knowslink-state}
POSTGRES_IMAGE=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24
export KNOWSLINK_STATE_DIR=$STATE KNOWSLINK_UID=$(id -u) KNOWSLINK_GID=$(id -g)

dc() { (cd "$DEPLOY" && docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file "$STATE/.env" "$@"); }
die() { echo "ERROR: $*" >&2; exit 1; }

prepare() { # prepare <sha>: detached checkout, 0700 state dir, new 0600 secrets
  sha=${1:?usage: beta.sh prepare <sha>}
  ! ss -ltnH | grep -Eq '(127\.0\.0\.1|0\.0\.0\.0|\*|\[::\]):8080 ' || die "port 8080 is busy"
  [ ! -e "$DEPLOY" ] || die "$DEPLOY exists"
  [ ! -e "$STATE/.env" ] || die "$STATE/.env exists; refusing to rotate secrets"
  git clone --quiet "$SRC" "$DEPLOY" && git -C "$DEPLOY" checkout --quiet --detach "$sha"
  install -d -m 0700 "$STATE" "$STATE/tunnel" "$STATE/backups"
  umask 077
  password=$(openssl rand -hex 24)
  cat >"$STATE/.env" <<ENV
POSTGRES_PASSWORD=$password
DATABASE_URL=postgres://knowslink:$password@postgres:5432/knowslink?sslmode=disable
RELAY_PORT=8080
ENV
  echo "prepared $(git -C "$DEPLOY" rev-parse HEAD); state $(stat -c %a "$STATE") env $(stat -c %a "$STATE/.env")"
}

up() { dc up -d --build --wait relay; dc ps -a; }
stop() { dc --profile tunnel stop; }  # volumes and other projects untouched
unexpose() { dc --profile tunnel stop cloudflared; }

backup() { # pg_dump -Fc into 0600 file, then list it
  out="$STATE/backups/$(git -C "$DEPLOY" rev-parse --short HEAD)-$(date -u +%Y%m%dT%H%M%SZ).dump"
  (umask 077; dc exec -T postgres pg_dump -U knowslink -Fc knowslink >"$out") || { rm -f "$out"; die "pg_dump failed"; }
  pg_list=$(docker run --rm -i "$POSTGRES_IMAGE" pg_restore --list <"$out" | wc -l) || { rm -f "$out"; die "dump unreadable"; }
  rows=$(dc exec -T postgres psql -U knowslink -Atc 'select count(*) from relay_state')
  echo "backup $(basename "$out") bytes=$(stat -c %s "$out") toc_lines=$pg_list relay_state_rows=$rows"
}

restore_verify() { # restore a dump into a throwaway container with no network; live DB untouched
  dump=${1:?usage: beta.sh restore-verify <dump>}
  name=knowslink-restore-check
  docker rm -f $name >/dev/null 2>&1 || true
  docker run -d --name $name --network none -e POSTGRES_HOST_AUTH_METHOD=trust "$POSTGRES_IMAGE" >/dev/null
  trap 'docker rm -f '$name' >/dev/null 2>&1 || true' EXIT  # RETURN would not fire when set -e aborts the function
  # the image starts a temporary init server first; ready is the second "ready to accept connections"
  for _ in $(seq 60); do [ "$(docker logs $name 2>&1 | grep -c 'ready to accept connections')" -ge 2 ] && break; sleep 1; done
  docker exec $name psql -U postgres -qc 'CREATE ROLE knowslink' -c 'CREATE DATABASE restored OWNER knowslink'
  docker exec -i $name pg_restore -U postgres -d restored --no-owner <"$dump"
  docker exec $name psql -U postgres -d restored -Atc \
    "select 'tables='||(select count(*) from information_schema.tables where table_schema='public')||' relay_state_rows='||(select count(*) from relay_state)||' state_type='||(select data_type from information_schema.columns where table_name='relay_state' and column_name in ('state','data'))||' human_model='||coalesce((select to_jsonb(r)->>'human_model' from relay_state r),'none')"
}

tunnel_create() {
  [ ! -e "$STATE/tunnel.uuid" ] || die "tunnel already recorded"
  out=$(cloudflared tunnel create --credentials-file "$STATE/tunnel/new.json" knowslink 2>&1) || die "tunnel create failed"
  uuid=$(grep -Eo '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' <<<"$out" | head -1)
  [ -n "$uuid" ] || die "no tunnel id in output"
  mv "$STATE/tunnel/new.json" "$STATE/tunnel/$uuid.json"; chmod 0600 "$STATE/tunnel/$uuid.json"
  echo "$uuid" >"$STATE/tunnel.uuid"; echo "tunnel $uuid"
}

render_public_config() { # render-only, no Access subscription or live config changes
  uuid=$(<"$STATE/tunnel.uuid")
  [[ "$uuid" =~ ^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$ ]] || die "invalid tunnel UUID"
  candidate="$STATE/tunnel/config.public.yml"
  (umask 077; set -C
    {
      printf 'tunnel: %s\ncredentials-file: /etc/cloudflared/%s.json\ningress:\n' "$uuid" "$uuid"
      sed 's/^/  /' "$DEPLOY/deploy/knowslink/tunnel/public-ingress.yml" || exit 1
      printf '  - service: http_status:404\n'
    } >"$candidate"
  ) || die "candidate exists or cannot be written; live config unchanged"
  cloudflared tunnel --config "$candidate" ingress validate
  echo "rendered public candidate; D12 review and edge app removal required before activation"
}

deploy() { # deploy <sha>: backup, move the detached checkout, rebuild. Also the rollback path (see D12 11.3)
  sha=${1:?usage: beta.sh deploy <sha>}
  current=$(git -C "$DEPLOY" rev-parse HEAD)
  git -C "$DEPLOY" fetch -q origin || echo "origin fetch failed; using commits already in $DEPLOY" >&2
  git -C "$DEPLOY" cat-file -e "$sha^{commit}" || die "unknown commit $sha"
  git -C "$DEPLOY" diff --quiet "$current" "$sha" -- db/migrations \
    || die "db/migrations differ between $current and $sha; schema compatibility needs D12 section 5 review"
  backup
  echo "$(date -u +%FT%TZ) $current -> $sha" >>"$STATE/deploy-history.log"
  git -C "$DEPLOY" checkout -q --detach "$sha"
  up
  python3 "$DEPLOY/deploy/knowslink/verify.py" local
}

{ # braces: bash parses the whole dispatch first, so "deploy" may replace this file mid-run
case "${1:-}" in
  prepare|up|stop|unexpose|backup|deploy) "$1" "${@:2}" ;;
  restore-verify) restore_verify "${@:2}" ;;
  tunnel-create) tunnel_create ;;
  render-public-config) render_public_config ;;
  *) die "usage: beta.sh prepare <sha>|up|stop|unexpose|backup|deploy <sha>|restore-verify <dump>|tunnel-create|render-public-config" ;;
esac
  exit
}
