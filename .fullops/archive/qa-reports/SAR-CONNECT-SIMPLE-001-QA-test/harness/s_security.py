# CS-AC-16: no key / token / credential text in any captured output, log or page; other members' state, names and e-mails not exposed.
from lib import *
import glob
secrets = {}
for root, _, files in os.walk(W):
    for f in files:
        p = os.path.join(root, f)
        if f in ("agent.json", "pending.json", "login.json") or f.endswith(".pem"):
            try: txt = open(p, encoding="utf-8").read()
            except Exception: continue
            if f.endswith(".pem"):
                body = "".join(l for l in txt.splitlines() if "-----" not in l)
                if len(body) > 20: secrets[f"pem:{os.path.relpath(p, W)}"] = body
            else:
                try: d = json.loads(txt)
                except Exception: continue
                for k, v in d.items():
                    if k in ("credential", "token", "private", "secret", "lease_token", "proof") and isinstance(v, str) and len(v) >= 16: secrets[f"{f}:{k}:{os.path.relpath(p, W)}"] = v
rec("CS-AC-16", "secret material found on disk to search for (private.pem bodies, credentials, login tokens)", len(secrets) > 0, f"{len(secrets)} secret strings")
targets = [p for p in glob.glob(os.path.join(OUT, "*")) if os.path.isfile(p) and p.endswith((".log", ".json", ".jsonl", ".txt")) and not p.endswith("pg.secret")]
hits = []
for p in targets:
    try: t = open(p, encoding="utf-8", errors="replace").read()
    except Exception: continue
    for k, v in secrets.items():
        if v in t: hits.append((os.path.basename(p), k))
rec("CS-AC-16", f"none of the secret strings appear in {len(targets)} captured outputs/logs (CLI JSON lines, MCP replies, fixture log, evidence)", not hits, f"hits={hits[:5]}")
pat = re.compile(r"BEGIN (?:OPENSSH |EC |RSA )?PRIVATE KEY|\"credential\"\s*:\s*\"[A-Za-z0-9_-]{30,}|lease_token|Bearer [A-Za-z0-9._-]{20,}")
pat_hits = [(os.path.basename(p), m.group(0)[:30]) for p in targets for m in [pat.search(open(p, encoding="utf-8", errors="replace").read())] if m]
rec("CS-AC-16", "pattern search (PEM header, credential value, lease token, bearer) over the same files: 0 hits", not pat_hits, f"hits={pat_hits[:5]}")
# member pages: other members' names, e-mails, agent ids not exposed
pages = ["/home", "/home/agents", "/home/friends", "/home/account", "/home/new"]
leaks = []
for x, others in (("a", ["b", "c", "d", "e"]), ("b", ["a", "c", "d"]), ("d", ["a", "b", "c", "e"])):
    b = Browser().login(x)
    other_agents = [a["id"] for o in others for a in (state(o)["agents"] or [])]
    for p in pages:
        s, raw, _ = b.get(p)
        for o in others:
            if f"csqa-{o}@" in raw: leaks.append((x, p, "email", o))
        for aid in other_agents:
            if aid in raw: leaks.append((x, p, "agent", aid[:14]))
rec("CS-AC-16", "member pages of A, B, D show no other non-friend member's e-mail or agent id", not leaks, f"leaks={leaks[:5]}")
