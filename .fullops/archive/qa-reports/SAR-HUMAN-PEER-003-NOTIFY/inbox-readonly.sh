set -eu
cd /home/shin/deploy/knowslink
date -u +%FT%TZ
docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file /home/shin/deploy/knowslink-state/.env exec -T postgres psql -U knowslink -d knowslink -At <<'SQL'
BEGIN READ ONLY;
SELECT jsonb_build_object('agent',i.value->>'Agent','inbox_ids',ARRAY(SELECT jsonb_object_keys(COALESCE(i.value->'Inbox','{}'::jsonb)))) FROM relay_state CROSS JOIN LATERAL jsonb_each(data->'Remote'->'Installs') i WHERE i.value->>'Agent'='agent_a3afe886b9fed6dac5c2b5';
COMMIT;
BEGIN READ ONLY;
SELECT jsonb_build_object('id',m.key,'receipt',m.value->'Receipt','persisted',m.value->'Persisted','claimed',m.value->'Claimed','completion',m.value->'Completion','attempts',m.value->'Attempts') FROM relay_state CROSS JOIN LATERAL jsonb_each(data->'Messages') m WHERE m.key IN ('01a13097-8352-7f8a-bf8f-d61f516b5d1e','01a1309e-da56-712a-998c-6019b240888f');
COMMIT;
SQL
