#!/usr/bin/env bash
# SAR-MCP-DEPLOY-001-OPS: read-only pre/post check. No keys, bodies, tokens or env values.
set -o errexit -o nounset -o pipefail
D=/home/shin/deploy/knowslink
echo "now=$(date -u +%FT%TZ)"
echo "head=$(git -C "$D" rev-parse HEAD)"
echo "dirty_lines=$(git -C "$D" status --porcelain --untracked-files=no | wc -l)"
docker ps --filter label=com.docker.compose.project=knowslink --format '{{.Names}} {{.Status}}'
for k in maxmemory maxmemory-policy; do docker exec knowslink-redis-1 redis-cli config get "$k" | paste -sd' ' -; done
docker exec knowslink-redis-1 redis-cli info memory | grep -E '^(used_memory|used_memory_human|used_memory_peak_human|maxmemory_human):' | tr -d '\r'
echo "dbsize=$(docker exec knowslink-redis-1 redis-cli dbsize)"
echo "redis_mem_limit=$(docker inspect -f '{{.HostConfig.Memory}}' knowslink-redis-1)"
docker exec -i -e PGOPTIONS='-c default_transaction_read_only=on' knowslink-postgres-1 psql -U knowslink -At -F ' ' <<'SQL'
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'agent', a.key, 'revoked=' || (a.value->>'Revoked') from s, jsonb_each(s.j->'Agents') a
 where a.key in ('agent_5513ec51367fed3748903d','agent_2b73608560c8ba6e8c3d06');
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'grant', g.value->>'From', g.value->>'To', g.value->>'Kind', g.value->>'Exp' from s, jsonb_each(s.j->'Grants') g order by 4;
SQL
