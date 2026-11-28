#!/usr/bin/env bash
# SAR-ORCA-DIAG-001-OPS: read-only trace. Prints type, id, size and time only; never text, keys, tokens or secrets.
set -o errexit -o nounset -o pipefail
psql_ro() { docker exec -i -e PGOPTIONS='-c default_transaction_read_only=on' knowslink-postgres-1 psql -U knowslink -At; }
echo "now=$(date -u +%FT%TZ)"
state=$(psql_ro <<'SQL'
select coalesce(to_jsonb(r)->'state', to_jsonb(r)->'data') from relay_state r
SQL
)
name=$(printf %s "$state" | python3 -I -c 'import json,sys; print(json.load(sys.stdin).get("Ephemeral",""))')
copy='{}'
if [ -n "$name" ]; then
  copy=$(docker exec knowslink-redis-1 redis-cli --raw GET "knowslink:messages:$name")
  echo "redis_ttl_s=$(docker exec knowslink-redis-1 redis-cli TTL "knowslink:messages:$name")"
fi
printf '%s\n%s\n' "$state" "${copy:-{\}}" | python3 -I -c '
import json, sys
st = json.loads(sys.stdin.readline())
eph = json.loads(sys.stdin.readline() or "{}")
who = ("agent_5513ec51367fed3748903d", "agent_2b73608560c8ba6e8c3d06")
for a, v in (st.get("Agents") or {}).items():
    if a in who or not v.get("Revoked"):
        print("agent", a, "revoked=%s" % v.get("Revoked"), "client=%s" % v.get("Client"), "created=%s" % v.get("Created"), "same_owner=%s" % (v.get("Owner") == st["Agents"][who[0]]["Owner"]))
for g in (st.get("Grants") or {}).values():
    print("grant", g["From"][:14], ">", g["To"][:14], g["Kind"], "sent_last_hour=%s" % (g.get("Sent") or []))
r = st.get("Remote") or {}
installs = r.get("Installs") or {}
for i, v in installs.items():
    tok = sorted((t["Kind"], t["Exp"]) for t in (r.get("Tokens") or {}).values() if t["Install"] == i)
    access = [e for k, e in tok if k == "access"]
    print("install", i[:6], v["Agent"], "exp=%s" % v["Exp"], "access_exp=%s" % access, "inbox_ids=%s" % sorted((eph.get("Inboxes") or {}).get(i, {})))
for s in (r.get("Subscriptions") or {}).values():
    print("subscription", s["Install"][:6], "exp=%s" % s["Exp"], "verified=%s" % s.get("Verified"), "event_at=%s" % s.get("EventAt"), "attempts=%s" % s.get("Attempts"), "pending=%s" % s.get("Pending"), "seen=%s" % s.get("Seen"))
msgs = eph.get("Messages") or {}
print("messages=%d" % len(msgs))
for i, m in sorted(msgs.items(), key=lambda kv: kv[1]["Receipt"]["accepted_at"]):
    c, env = m["Receipt"], m.get("Envelope")
    print("message", i, c["from"][:14], ">", c["to"][:14], c["intent"], "transport=%s" % c["transport"], "accepted=%s" % c["accepted_at"], "exp=%s" % c["exp"],
          "attempts=%s" % m.get("Attempts"), "persisted=%s" % m.get("Persisted"), "lease_until=%s" % m.get("LeaseUntil"), "session_seen=%s" % m.get("SessionSeen"),
          "envelope=%s" % (len(json.dumps(env)) if env else None), "kind=%s" % (env.get("kind") if isinstance(env, dict) else None), "web_copy=%s" % bool(m.get("Conversation")), "reply_id=%s" % m.get("ReplyID"))
'
echo '--- relay log since 3h (all lines)'
docker logs --since 3h knowslink-relay-1 2>&1 | tail -n 30
docker ps --filter name=knowslink --format '{{.Names}} {{.Status}}'
echo "head=$(git -C /home/shin/deploy/knowslink rev-parse HEAD)"
