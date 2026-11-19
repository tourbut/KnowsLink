set -eu
cd /home/shin/deploy/knowslink
date -u +%FT%TZ
docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file /home/shin/deploy/knowslink-state/.env exec -T postgres psql -U knowslink -d knowslink -At <<'SQL'
BEGIN READ ONLY;
SELECT jsonb_build_object('id', s.key, 'agent', data->'Remote'->'Installs'->(s.value->>'Install')->>'Agent', 'expiry',s.value->>'Exp','verified',s.value->>'Verified','event_at',s.value->>'EventAt','attempts',s.value->'Attempts','next',s.value->>'Next','seen_set',COALESCE(s.value->>'Seen','')<>'','pending_set',COALESCE(s.value->>'Pending','')<>'') FROM relay_state CROSS JOIN LATERAL jsonb_each(COALESCE(data->'Remote'->'Subscriptions','{}'::jsonb)) s;
COMMIT;
SQL
