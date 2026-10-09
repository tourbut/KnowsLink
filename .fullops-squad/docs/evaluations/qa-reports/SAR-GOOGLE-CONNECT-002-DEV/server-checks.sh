#!/usr/bin/env bash
# Reproduce this task's isolated server checks; never apply production configuration or access credentials.
set -eu
base=/tmp/knowslink-google002-ctx9dc5e970
cd "$base/work"
tar -xzf ../fullops-assets.tar.gz -C ../fullops
docker ps --format '{{.Names}} {{.ID}}' | sort >../services-before.txt
pg=knowslink-dev-google002-ctx9dc5e970
cf=knowslink-cf-google002-ctx9dc5e970
trap 'docker rm -f "$pg" "$cf" >/dev/null 2>&1 || true' EXIT
image=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24
docker run -d --name "$pg" --memory 256m --cpus 1 -e POSTGRES_PASSWORD=example-local-only -e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 "$image" >/dev/null
for i in $(seq 60); do docker exec "$pg" pg_isready -U knowslink >/dev/null 2>&1 && break; sleep 1; done
port=$(docker port "$pg" 5432/tcp | cut -d: -f2)
export DATABASE_URL="postgres://knowslink:example-local-only@127.0.0.1:$port/knowslink?sslmode=disable"
export TEST_DATABASE_URL="$DATABASE_URL" TEST_SYNTHETIC_DATABASE=1
go run ./cmd/migrate >../migration.log 2>&1
go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP|TestGoogleHTTP|TestEmailIdentity' -count=1 >../integration.log 2>&1
mkdir -p ../tools ../render-state/tunnel
docker create --name "$cf" cloudflare/cloudflared:2026.9.1 >/dev/null
docker cp "$cf":/usr/local/bin/cloudflared ../tools/cloudflared
docker rm "$cf" >/dev/null
chmod +x ../tools/cloudflared
export PATH="$base/tools:$PATH" KNOWSLINK_DEPLOY="$base/work" KNOWSLINK_STATE_DIR="$base/render-state"
printf '%s\n' 00000000-0000-4000-8000-000000000001 >../render-state/tunnel.uuid
bash deploy/knowslink/beta.sh render-public-config >../cloudflared.log 2>&1
python3 - <<'PY' >>../cloudflared.log 2>&1
import os, subprocess
config = os.environ['KNOWSLINK_STATE_DIR'] + '/tunnel/config.public.yml'
paths = [('/',0),('/auth/google/callback',0),('/home',0),('/connect/'+'a'*43,0),('/v1/connect/start',0),('/v1/text/send',0),('/v1/keys/agent/key',0),('/v1/receipts/00000000-0000-7000-8000-000000000001',0),('/owner',1),('/v1/owners',1),('/v1/keys',1),('/v1/authorize',1),('/v1/test/pull',1),('/healthz',1),('/unlisted',1)]
for path, rule in paths:
    result = subprocess.run(['cloudflared','tunnel','--config',config,'ingress','rule','https://link.knowslog.com'+path],capture_output=True,text=True)
    assert result.returncode == 0, result.stderr
    assert f'Matched rule #{rule}' in result.stdout, result.stdout
    print(path, 'rule', rule, 'exit', result.returncode)
PY
KNOWSLINK_PREFIX="$base/bot-install" sh scripts/install_bot_mcp.sh >../bot-install.log 2>&1
docker rm -f "$pg" >/dev/null
trap - EXIT
docker ps --format '{{.Names}} {{.ID}}' | sort >../services-after.txt
cmp ../services-before.txt ../services-after.txt
printf 'migration=0 integration=0 cloudflared=0 bot-install=0 shared-services-unchanged=0\n'
