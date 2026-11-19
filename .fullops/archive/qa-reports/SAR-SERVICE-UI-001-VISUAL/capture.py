"""Synthetic member-UI visual capture for SAR-SERVICE-UI-001-VISUAL.

Drives Orca's embedded browser through the orca CLI, saves PNG files (never prints
base64) and writes manifest.json + measurements.json. Usage:
  python3 -I capture.py <fixture_dir> <out_dir>
"""
import base64, hashlib, json, os, subprocess, sys

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
FX, OUT = sys.argv[1], sys.argv[2]
VPS = {"1440": (1440, 900), "390": (390, 844), "320": (320, 640)}
OPEN_SET = {"home", "home-empty", "connection-prepared", "connection-token"}
KEY320 = {"home", "home-empty", "start-both-problem", "connection-token", "connection-prepared",
          "device-prepared", "receipt", "refusal-limited", "verify"}

MEASURE = r"""(()=>{
const iw=innerWidth, de=document.documentElement;
const vis=e=>{const r=e.getBoundingClientRect();const s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none'};
const name=e=>(e.getAttribute('aria-label')||e.textContent||e.value||'').trim().replace(/\s+/g,' ').slice(0,40);
const over=[...document.querySelectorAll('body *')].filter(e=>vis(e)).filter(e=>{const r=e.getBoundingClientRect();return r.right>iw+1||r.left<-1}).slice(0,6).map(e=>e.tagName.toLowerCase()+'.'+e.className+':'+name(e));
const tg='a[href],button,input:not([type=hidden]),select,textarea,summary';
const small=[...document.querySelectorAll(tg)].filter(vis).filter(e=>{const r=e.getBoundingClientRect();return r.height<43.5&&!(e.className||'').toString().includes('skip')}).slice(0,8).map(e=>e.tagName.toLowerCase()+':'+name(e)+' h='+Math.round(e.getBoundingClientRect().height));
const btnNames={};[...document.querySelectorAll('button,input[type=submit]')].filter(vis).forEach(e=>{const n=name(e);btnNames[n]=(btnNames[n]||0)+1});
const dupBtn=Object.entries(btnNames).filter(([k,v])=>v>1).map(([k,v])=>k+' x'+v);
const nolabel=[...document.querySelectorAll('input:not([type=hidden]),select,textarea')].filter(vis).filter(e=>!(e.labels&&e.labels.length)&&!e.getAttribute('aria-label')&&!e.getAttribute('aria-labelledby')).map(e=>e.name||e.tagName);
const bs=getComputedStyle(document.body);
const codes=[...document.querySelectorAll('code,pre')].filter(vis);
const fs=[...document.querySelectorAll('body *')].filter(vis).filter(e=>[...e.childNodes].some(n=>n.nodeType===3&&n.textContent.trim())).map(e=>parseFloat(getComputedStyle(e).fontSize));
const hs=[...document.querySelectorAll('h1,h2,h3,h4')].map(e=>+e.tagName[1]);
let skipped=false;for(let i=1;i<hs.length;i++)if(hs[i]-hs[i-1]>1)skipped=true;
const html=de.outerHTML;
return JSON.stringify({sw:de.scrollWidth,iw,docH:de.scrollHeight,title:document.title,lang:de.lang,
 overflow:over,smallTargets:small,dupButtonNames:dupBtn,unlabeledInputs:nolabel,
 bodyWordBreak:bs.wordBreak,bodyOverflowWrap:bs.overflowWrap,bodyFont:bs.fontSize+'/'+bs.lineHeight,
 minFont:Math.min(...fs),codeOverflowWrap:[...new Set(codes.map(e=>getComputedStyle(e).overflowWrap))],
 h1:document.querySelectorAll('h1').length,headingSkip:skipped,
 landmarks:{header:document.querySelectorAll('header').length,nav:document.querySelectorAll('nav').length,main:document.querySelectorAll('main').length},
 skipLink:!!document.querySelector('a[href="#main"]'),
 detailsTotal:document.querySelectorAll('details').length,detailsOpen:document.querySelectorAll('details[open]').length,
 external:(html.match(/\ssrc=|<link|url\(|<script|@import/g)||[]).length,
 roleAlert:document.querySelectorAll('[role=alert]').length,roleStatus:document.querySelectorAll('[role=status]').length,
 textLen:document.body.innerText.length});
})()"""


def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + " " + p.stderr[:200])
    return json.loads(p.stdout)["result"]


def shot(path, w, h, docH, fold=True):
    """Viewport PNGs scrolled in h-high steps (-sKofN for tall pages). full-screenshot and
    tall-viewport capture only tile the first viewport under emulation, so neither is used."""
    n = max(1, -(-docH // h))
    out = []
    for i in range(n):
        y = min(i * h, max(0, docH - h))
        run("eval", "--expression", f"scrollTo(0,{y});scrollY")
        r = run("screenshot")
        p = path if n == 1 else path.replace(".png", f"-s{i + 1}of{n}.png")
        with open(p, "wb") as f:
            f.write(base64.b64decode(r["data"]))
        b = open(p, "rb").read()
        out.append((p, hashlib.sha256(b).hexdigest(), len(b), h))
    run("eval", "--expression", "scrollTo(0,0);1")
    return out


def measure():
    return json.loads(run("eval", "--expression", MEASURE)["result"])


manifest, meas = [], {}
os.makedirs(os.path.join(OUT, "png"), exist_ok=True)
for fn in sorted(os.listdir(FX)):
    if not fn.endswith(".html"):
        continue
    n = fn[:-5]
    url = "file:///" + os.path.abspath(os.path.join(FX, fn)).replace("\\", "/")
    for vp, (w, h) in VPS.items():
        run("goto", "--url", url)  # goto resets emulation, so set the viewport after it
        run("exec", "--command", f"set viewport {w} {h}")
        m = measure()
        assert m["iw"] == w, (n, vp, m["iw"])  # fail loudly if emulation did not apply
        meas[f"{n}@{vp}"] = m
        if vp != "320" or n in KEY320:
            p = f"png/{n}-{vp}.png"
            for fp, sha, size, vh in shot(os.path.join(OUT, p), w, h, m["docH"]):
                manifest.append({"screen": n, "viewport": f"{w}x{vh}", "state": "default", "file": os.path.relpath(fp, OUT).replace("\\", "/"), "sha256": sha, "bytes": size, "docH": m["docH"]})
        if n in OPEN_SET:
            run("eval", "--expression", "document.querySelectorAll('details').forEach(d=>d.open=true);1")
            m2 = measure()
            meas[f"{n}@{vp}+open"] = m2
            if vp != "320":
                p = f"png/{n}-{vp}-open.png"
                for fp, sha, size, vh in shot(os.path.join(OUT, p), w, h, m2["docH"], fold=False):
                    manifest.append({"screen": n, "viewport": f"{w}x{vh}", "state": "all-details-open", "file": os.path.relpath(fp, OUT).replace("\\", "/"), "sha256": sha, "bytes": size, "docH": m2["docH"]})
    print("done", n, flush=True)

json.dump(manifest, open(os.path.join(OUT, "manifest.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
json.dump(meas, open(os.path.join(OUT, "measurements.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
print("PNG", len(manifest), "measure keys", len(meas))
