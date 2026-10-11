"""V-01 regression: measure every visible a[href] (plus buttons/inputs/summary) at 320/390/1440.

Drives Orca's embedded browser via the orca CLI. Usage: python3 -I measure_links.py <fixture_dir> <out.json>
Fixtures: KNOWSLINK_UI_FIXTURE_DIR=<dir> go test ./internal/relay -run TestMemberUIFixtures -count=1
"""
import json, os, subprocess, sys

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
FX, OUT = sys.argv[1], sys.argv[2]
VPS = {"1440": (1440, 900), "390": (390, 844), "320": (320, 640)}
MEASURE = r"""(()=>{
const vis=e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none'};
const name=e=>(e.getAttribute('aria-label')||e.textContent||'').trim().replace(/\s+/g,' ').slice(0,40);
document.querySelectorAll('details').forEach(d=>d.open=true);
const iw=innerWidth,de=document.documentElement;
const links=[...document.querySelectorAll('a[href]')].filter(vis);
const other=[...document.querySelectorAll('button,input:not([type=hidden]),select,summary')].filter(vis);
const small=l=>l.filter(e=>e.getBoundingClientRect().height<43.5).map(e=>e.tagName.toLowerCase()+':'+name(e)+' h='+e.getBoundingClientRect().height);
const over=[...document.querySelectorAll('body *')].filter(vis).filter(e=>{const r=e.getBoundingClientRect();return r.right>iw+1||r.left<-1}).map(e=>e.tagName.toLowerCase()+'.'+e.className);
return JSON.stringify({iw,sw:de.scrollWidth,links:links.map(name),others:other.length,smallLinks:small(links),smallOthers:small(other),overflow:over,scrollOver:de.scrollWidth>iw});
})()"""


def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + " " + p.stderr[:200])
    return json.loads(p.stdout)["result"]


res, bad = {}, 0
for fn in sorted(os.listdir(FX)):
    if not fn.endswith(".html"):
        continue
    url = "file:///" + os.path.abspath(os.path.join(FX, fn)).replace("\\", "/")
    for vp, (w, h) in VPS.items():
        run("goto", "--url", url)
        run("exec", "--command", f"set viewport {w} {h}")
        m = json.loads(run("eval", "--expression", MEASURE)["result"])
        assert m["iw"] == w, (fn, vp, m["iw"])
        res[f"{fn[:-5]}@{vp}"] = m
        # the 15px (scrollbar) gap seen in older reports is emulation; real overflow = overflow elements
        bad += len(m["smallLinks"]) + len(m["smallOthers"]) + len(m["overflow"])
json.dump(res, open(OUT, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
links = sum(len(m["links"]) for m in res.values())
print("keys", len(res), "link measurements", links, "violations", bad)
sys.exit(1 if bad else 0)
