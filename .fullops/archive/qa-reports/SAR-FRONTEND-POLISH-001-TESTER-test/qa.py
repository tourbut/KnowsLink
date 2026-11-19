"""SAR-FRONTEND-POLISH-001-TESTER narrow independent QA over synthetic fixtures (Orca embedded browser).

Fixtures are rendered by `KNOWSLINK_UI_FIXTURE_DIR=<dir> go test ./internal/relay -run TestMemberUIFixtures` from the
review SHA's product code and served only from http://127.0.0.1:<port>. No production origin, account or key is used.
Usage: python3 -I qa.py <origin> <out_json>
"""
import json, subprocess, sys

sys.stdout.reconfigure(encoding="utf-8")

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
ORIGIN, OUT = sys.argv[1].rstrip("/"), sys.argv[2]
FIXTURES = ["connection-approved", "connection-cancelled", "connection-consumed", "connection-expired", "connection-prepared",
            "connection-token", "connection-waiting", "device-approved", "device-login", "device-login-off", "device-prepared",
            "gate-active", "gate-closed", "google-complete-connect", "google-complete-home", "home", "home-empty", "receipt",
            "refusal-create", "refusal-limited", "refusal-login", "refusal-reauth", "remote-consent", "remote-login",
            "start-both-problem", "start-google", "start-logout", "start-none", "verify"]
VPS = {"1440": (1440, 900), "390": (390, 844), "320": (320, 640)}
PROBE = r"""(()=>{
const vis=e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none'};
const nm=e=>(e.getAttribute('aria-label')||e.textContent||e.value||'').trim().replace(/\s+/g,' ');
document.querySelectorAll('details').forEach(d=>d.open=true);
const de=document.documentElement,iw=innerWidth;
const ctl=[...document.querySelectorAll('a[href],button,input:not([type=hidden]),select,summary')].filter(vis);
const small=ctl.filter(e=>e.getBoundingClientRect().height<43.5&&!e.classList.contains('skip')).map(e=>e.tagName+':'+nm(e).slice(0,30));
const over=[...document.querySelectorAll('body *')].filter(vis).filter(e=>{const r=e.getBoundingClientRect();return r.right>iw+1||r.left<-1}).map(e=>e.tagName.toLowerCase()+'.'+e.className);
const btn={};document.querySelectorAll('button').forEach(b=>{if(vis(b)){const k=nm(b);btn[k]=(btn[k]||0)+1}});
const emptySel=[...document.querySelectorAll('select')].filter(s=>s.options.length===0).map(s=>s.id||s.name);
const forms=[...document.querySelectorAll('form')].map(f=>(f.method||'get')+' '+f.getAttribute('action'));
return JSON.stringify({iw,sw:de.scrollWidth,over,small,mains:document.querySelectorAll('main').length,h1:document.querySelectorAll('h1').length,
 skip:!!document.querySelector('a.skip[href="#main"]'),dupBtn:Object.keys(btn).filter(k=>btn[k]>1),emptySel,forms,
 primaries:[...document.querySelectorAll('button.primary,.button.primary')].filter(vis).map(nm),
 buttons:[...document.querySelectorAll('button')].filter(vis).map(nm)});
})()"""
ACTIVE = r"""(()=>{const e=document.activeElement,s=getComputedStyle(e);return JSON.stringify({tag:e.tagName.toLowerCase(),
name:(e.getAttribute('aria-label')||e.textContent||'').trim().replace(/\s+/g,' ').slice(0,40),outline:s.outlineStyle+' '+s.outlineWidth,
top:Math.round(e.getBoundingClientRect().top)})})()"""


def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + " " + (p.stdout or p.stderr)[:300])
    return json.loads(p.stdout)["result"]


def ev(js):
    return json.loads(run("eval", "--expression", js)["result"])


def open_vp(name, vp):
    w, h = VPS[vp]
    run("goto", "--url", f"{ORIGIN}/{name}.html")
    for _ in range(5):  # viewport emulation applies asynchronously
        run("exec", "--command", f"set viewport {w} {h}")
        m = ev(PROBE)
        if m["iw"] == w:
            return m
    raise SystemExit(f"viewport not applied {name} {vp}: {m['iw']}")


res, fails = {}, []


def check(cond, msg):
    if not cond:
        fails.append(msg)


for n in FIXTURES:
    for vp in VPS:
        m = open_vp(n, vp)
        res[f"{n}@{vp}"] = m
        check(m["sw"] <= m["iw"] + 1 and not m["over"], f"{n}@{vp} horizontal overflow sw={m['sw']} iw={m['iw']} {m['over'][:2]}")
        check(not m["small"], f"{n}@{vp} control under 44px {m['small']}")
        check(m["mains"] == 1 and m["h1"] == 1 and m["skip"], f"{n}@{vp} landmarks main={m['mains']} h1={m['h1']} skip={m['skip']}")
        check(not m["dupBtn"], f"{n}@{vp} duplicate button names {m['dupBtn']}")
        check(not m["emptySel"], f"{n}@{vp} empty select {m['emptySel']}")

# changed behaviour: empty home hides unusable forms, connected home keeps them, OAuth approve is the one primary
for vp in VPS:
    he = res[f"home-empty@{vp}"]
    check(not any(f.endswith("/home/invites") or f.endswith("/home/receipts") for f in he["forms"]), f"home-empty@{vp} unusable form {he['forms']}")
    h = res[f"home@{vp}"]
    check("post /home/invites" in h["forms"] and "get /home/receipts" in h["forms"], f"home@{vp} invite/receipt form missing {h['forms']}")
    rc = res[f"remote-consent@{vp}"]
    check(rc["primaries"] == ["이 설치 연결 승인"] and "취소" in rc["buttons"], f"remote-consent@{vp} primaries {rc['primaries']}")
    rl = res[f"remote-login@{vp}"]
    check(rl["primaries"] == ["Google 계정 확인"], f"remote-login@{vp} primaries {rl['primaries']}")

# accessible names: visible text must be contained in the accessible name (WCAG 2.5.3)
open_vp("home", "390")
rb = ev("""JSON.stringify([...document.querySelectorAll('form[action="/auth/reauth"] button')].map(b=>({aria:b.getAttribute('aria-label'),text:b.textContent.trim()})))""")
res["home-reauth-buttons"] = rb
check(all((b["aria"] or b["text"]).startswith(b["text"]) for b in rb), f"reauth accessible name lacks visible text {rb}")

for vp in ("390", "1440"):
    open_vp("home-empty", vp)
    t = ev("JSON.stringify(document.body.innerText)")
    res[f"home-empty-text@{vp}"] = [s for s in ("에이전트를 연결한 뒤", "에이전트를 만든 뒤 요청 상태를 확인할 수 있습니다.", "관계가 없습니다") if s in t]
    check("에이전트를 만든 뒤 요청 상태를 확인할 수 있습니다." in t, f"home-empty@{vp} receipt guidance missing")

# verify: six digits fit in the enlarged code input without clipping
for vp in VPS:
    open_vp("verify", vp)
    r = ev("""(()=>{const i=document.querySelector('input.code');i.value='000000';const r=i.getBoundingClientRect();
return JSON.stringify({w:r.width,h:r.height,sw:i.scrollWidth,cw:i.clientWidth,fs:getComputedStyle(i).fontSize})})()""")
    res[f"verify-code@{vp}"] = r
    check(r["sw"] <= r["cw"] + 1 and r["h"] >= 56, f"verify@{vp} code input clipped or short {r}")

# connection-token: the deadline note stays on one line
for vp in VPS:
    open_vp("connection-token", vp)
    r = ev("""(()=>{const s=document.querySelector('.nowrap');const g=new Set([...s.getClientRects()].map(x=>Math.round(x.top)));
return JSON.stringify({text:s.textContent,lines:g.size})})()""")
    res[f"token-note@{vp}"] = r
    check(r["lines"] == 1, f"connection-token@{vp} deadline note wraps {r}")

# keyboard (DOM substitute): Orca's CLI keypress did not move focus in this session (activeElement stayed body after
# click+Tab, also with --page), so real key delivery is NOT exercised. Instead: DOM tab order, programmatic focus ring.
TABS = r"""(()=>{const vis=e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'};
const out=[];for(const e of document.querySelectorAll('a[href],button,input:not([type=hidden]),select,summary,textarea,[tabindex]')){
if(e.tabIndex<0||e.disabled)continue;e.focus();const s=getComputedStyle(e);
out.push({name:(e.getAttribute('aria-label')||e.textContent||'').trim().replace(/\s+/g,' ').slice(0,40),focused:document.activeElement===e,
outline:s.outlineStyle+' '+s.outlineWidth,top:Math.round(e.getBoundingClientRect().top),vis:vis(e)})}
document.activeElement.blur();return JSON.stringify(out)})()"""
kb = {}
for n in ("remote-consent", "remote-login", "home-empty", "home"):
    open_vp(n, "390")
    seq = ev(TABS)
    kb[n] = seq
    # The embedded browser view is not OS-focused here, so :focus/:focus-visible styles never apply (skip link stays at
    # top:-80px, outline none). Ring and skip-link reveal are therefore NOT verified; only order and focusability are.
    check(seq[0]["name"] == "본문으로 건너뛰기", f"{n} first stop is not the skip link {seq[0]}")
    check(all(st["focused"] for st in seq), f"{n} element not focusable {[st for st in seq if not st['focused']]}")
res["keyboard-dom"] = kb
exp = {"remote-consent": ["본문으로 건너뛰기", "이 설치 연결 승인", "취소", "내 agent 관리"],
       "remote-login": ["본문으로 건너뛰기", "Google 계정 확인", "내 agent 관리"]}
for n, names in exp.items():
    got = [st["name"] for st in kb[n]]
    check(got == names, f"{n} tab order {got} != {names}")
open_vp("remote-consent", "390")
res["skip-target"] = ev("""JSON.stringify({href:document.querySelector('a.skip').getAttribute('href'),main:!!document.querySelector('main#main')})""")
check(res["skip-target"] == {"href": "#main", "main": True}, f"skip link target {res['skip-target']}")

json.dump({"results": res, "failures": fails}, open(OUT, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
print("fixtures", len(FIXTURES), "viewports", len(VPS), "failures", len(fails))
for f in fails:
    print(" ", f)
sys.exit(1 if fails else 0)
