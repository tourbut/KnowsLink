#!/usr/bin/env bash
# SAR-ORCA-HOST-001-OPS: read-only. Prints only agent id/revoked and grant from/to/kind/expiry; no keys, bodies or tokens.
set -o errexit -o nounset -o pipefail
docker exec -i -e PGOPTIONS='-c default_transaction_read_only=on' knowslink-postgres-1 psql -U knowslink -At -F ' ' <<'SQL'
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'agent', a.key, 'revoked=' || (a.value->>'Revoked') from s, jsonb_each(s.j->'Agents') a
 where a.key in ('agent_5513ec51367fed3748903d','agent_2b73608560c8ba6e8c3d06');
with s as (select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') j from relay_state r)
select 'grant', g.value->>'From', g.value->>'To', g.value->>'Kind', g.value->>'Exp' from s, jsonb_each(s.j->'Grants') g order by 4;
SQL
