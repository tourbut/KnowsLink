#!/usr/bin/env bash
# Owner-only synthetic beta lifecycle for the separate "knowslink" Compose project; never touches other projects.
set -euo pipefail

SRC=${KNOWSLINK_SRC:-$(git -C "$(dirname "$0")" rev-parse --show-toplevel)}
DEPLOY=${KNOWSLINK_DEPLOY:-/home/shin/deploy/knowslink}
STATE=${KNOWSLINK_STATE_DIR:-/home/shin/deploy/knowslink-state}
HOST=link.knowslog.com
TEAM=${KNOWSLINK_TEAM:-scshin88}
POSTGRES_IMAGE=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24
export KNOWSLINK_STATE_DIR=$STATE KNOWSLINK_UID=$(id -u) KNOWSLINK_GID=$(id -g)

dc() { (cd "$DEPLOY" && docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file "$STATE/.env" "$@"); }
die() { echo "ERROR: $*" >&2; exit 1; }

prepare() { # prepare <sha>: detached checkout, 0700 state dir, new 0600 secrets
  sha=${1:?usage: beta.sh prepare <sha>}
  ! ss -ltn | grep -q '127.0.0.1:8080 \|0.0.0.0:8080 ' || die "port 8080 is busy"
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
  (umask 077; dc exec -T postgres pg_dump -U knowslink -Fc knowslink >"$out")
  pg_list=$(docker run --rm -i "$POSTGRES_IMAGE" pg_restore --list <"$out" | wc -l)
  echo "backup $(basename "$out") bytes=$(stat -c %s "$out") toc_lines=$pg_list"
}

restore_verify() { # restore a dump into a throwaway container with no network; live DB untouched
  dump=${1:?usage: beta.sh restore-verify <dump>}
  name=knowslink-restore-check
  docker rm -f $name >/dev/null 2>&1 || true
  docker run -d --name $name --network none -e POSTGRES_PASSWORD="$(openssl rand -hex 12)" "$POSTGRES_IMAGE" >/dev/null
  trap 'docker rm -f '$name' >/dev/null 2>&1' RETURN
  for _ in $(seq 30); do docker exec $name pg_isready -U postgres >/dev/null 2>&1 && break; sleep 1; done
  docker exec $name psql -U postgres -qc 'CREATE ROLE knowslink' -c 'CREATE DATABASE restored OWNER knowslink'
  docker exec -i $name pg_restore -U postgres -d restored --no-owner <"$dump"
  docker exec $name psql -U postgres -d restored -Atc \
    "select 'tables='||count(*) from information_schema.tables where table_schema='public'"
}

seed() { # synthetic fixture on the server; adapters only talk to loopback. Gate expires 180 s after this runs
  cd "$DEPLOY"
  [ -d adapters/dist ] || { npm ci --prefix adapters --silent && npm run build --prefix adapters --silent; }
  node adapters/dist/synthetic.js http://127.0.0.1:8080 --seed
}

owner_login() { # browser HTTP Basic prompt values for the synthetic owner B; for the beta owner's own terminal only
  python3 - "$DEPLOY/build/qa-fixture.json" <<'PY'
import json, sys
owner = json.load(open(sys.argv[1]))["b"]["owner"]
print("username:", owner["owner"])
print("password:", owner["credential"])
PY
}

tunnel_create() {
  [ ! -e "$STATE/tunnel.uuid" ] || die "tunnel already recorded"
  out=$(cloudflared tunnel create --credentials-file "$STATE/tunnel/new.json" knowslink 2>&1) || die "tunnel create failed"
  uuid=$(grep -Eo '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' <<<"$out" | head -1)
  [ -n "$uuid" ] || die "no tunnel id in output"
  mv "$STATE/tunnel/new.json" "$STATE/tunnel/$uuid.json"; chmod 0600 "$STATE/tunnel/$uuid.json"
  echo "$uuid" >"$STATE/tunnel.uuid"; echo "tunnel $uuid"
}

render_config() { # needs access.aud from access_apply.py
  uuid=$(<"$STATE/tunnel.uuid"); aud=$(<"$STATE/access.aud")
  sed -e "s/__TUNNEL_UUID__/$uuid/g" -e "s/__TEAM_NAME__/$TEAM/" -e "s/__AUD_TAG__/$aud/" \
    "$DEPLOY/deploy/knowslink/tunnel/config.yml.tmpl" >"$STATE/tunnel/config.yml"
  chmod 0600 "$STATE/tunnel/config.yml"
  cloudflared tunnel --config "$STATE/tunnel/config.yml" ingress validate
}

expose() { # gate: protected Access app recorded, config requires the JWT, DNS record absent
  [ -s "$STATE/access.aud" ] || die "Access app not recorded; run access_apply.py first"
  grep -q 'required: true' "$STATE/tunnel/config.yml" || die "origin JWT check missing in tunnel config"
  curl -fsS -o /dev/null "http://127.0.0.1:8080/healthz" || die "relay not healthy"
  ! dig +short "$HOST" | grep -q . || die "$HOST already resolves; not overwriting"
  cloudflared tunnel route dns "$(<"$STATE/tunnel.uuid")" "$HOST"
  dc --profile tunnel up -d cloudflared
}

case "${1:-}" in
  prepare|up|stop|unexpose|backup|expose|seed) "$1" "${@:2}" ;;
  owner-login) owner_login ;;
  restore-verify) restore_verify "${@:2}" ;;
  tunnel-create) tunnel_create ;;
  render-config) render_config ;;
  *) die "usage: beta.sh prepare <sha>|up|stop|unexpose|backup|restore-verify <dump>|seed|owner-login|tunnel-create|render-config|expose" ;;
esac
