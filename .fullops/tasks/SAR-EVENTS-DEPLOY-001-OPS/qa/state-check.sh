#!/usr/bin/env bash
# SAR-EVENTS-DEPLOY-001-OPS: read-only check. Prints counts and names only: no keys, secrets, tokens or env values.
set -o errexit -o nounset -o pipefail
D=/home/shin/deploy/knowslink
STATE=/home/shin/deploy/knowslink-state
NEW=65a752fe2e88e0da1cbf1a844cb0000c51595295
echo "now=$(date -u +%FT%TZ)"
echo "head=$(git -C "$D" rev-parse HEAD)"
echo "dirty_lines=$(git -C "$D" status --porcelain --untracked-files=no | wc -l)"
echo "target_commit_present=$(git -C "$D" cat-file -e "$NEW^{commit}" 2>/dev/null && echo yes || echo no)"
docker ps --filter label=com.docker.compose.project=knowslink --format '{{.Names}} {{.Status}}'
stat -c 'env mode=%a bytes=%s' "$STATE/.env"
echo "seal_key_lines=$(grep -c '^KNOWSLINK_SEAL_KEY=' "$STATE/.env" || true)"
awk -F= '$1=="KNOWSLINK_SEAL_KEY"{print "seal_key_value_length=" length(substr($0, index($0,"=")+1))}' "$STATE/.env"
for k in maxmemory maxmemory-policy; do docker exec knowslink-redis-1 redis-cli config get "$k" </dev/null | paste -sd' ' -; done
docker exec knowslink-redis-1 redis-cli info memory </dev/null | grep -E '^(used_memory|used_memory_human|used_memory_peak_human):' | tr -d '\r'
echo "dbsize=$(docker exec knowslink-redis-1 redis-cli dbsize </dev/null)"
echo "relay_error_lines=$(docker logs knowslink-relay-1 2>&1 | grep -ci error || true)"
docker logs knowslink-relay-1 2>&1 | grep 'relay starting' | tail -n 1 || true
docker exec -i -e PGOPTIONS='-c default_transaction_read_only=on' knowslink-postgres-1 psql -U knowslink -At -F ' ' -v ON_ERROR_STOP=1 <<'SQL'
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'agent', a.key, 'revoked=' || (a.value->>'Revoked') from s, jsonb_each(s.j->'Agents') a
 where a.key in ('agent_5513ec51367fed3748903d','agent_2b73608560c8ba6e8c3d06');
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'grant', g.value->>'From', g.value->>'To', g.value->>'Kind', g.value->>'Exp' from s, jsonb_each(s.j->'Grants') g order by 4;
-- Events subscriptions: stored Secret/OldSecret by form. Values are never selected.
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r),
     subs as (select e.value v from s, jsonb_each(case when jsonb_typeof(s.j->'Remote'->'Subscriptions')='object'
                                                      then s.j->'Remote'->'Subscriptions' else '{}'::jsonb end) e)
select 'subscriptions total=' || count(*),
       'plain=' || count(*) filter (where v->>'Secret' like 'whsec\_%' or v->>'OldSecret' like 'whsec\_%'),
       'sealed=' || count(*) filter (where v->>'Secret' like 'sealed:v1:%'),
       'empty=' || count(*) filter (where coalesce(v->>'Secret','')='') from subs;
SQL
