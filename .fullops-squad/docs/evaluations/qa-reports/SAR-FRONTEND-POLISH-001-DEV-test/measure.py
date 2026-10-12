"""SAR-FRONTEND-POLISH-001-DEV before/after inspection (adapted from SAR-SERVICE-UI-001-FIX-VISUAL measure.py).

Drives Orca's embedded browser via the orca CLI over synthetic fixtures only. Measures every visible
a[href]/button/input/select/summary at 1440x900, 390x844 and 320x640 with all details open: height,
line count, overflow, overlap, plus landmarks, skip link, h1 count, duplicate unlabeled button names,
and agent summary rows whose disclosure marker sits alone on a line. PNGs are saved as files only.
Usage: python3 -I measure.py <fixture_dir> <out_dir>
"""
import base64, hashlib, json, os, subprocess, sys

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
FX, OUT = sys.argv[1], sys.argv[2]
VPS = {"1440": (1440, 900), "390": (390, 844), "320": (320, 640)}
# representative PNGs only (V-01 link types + neighbours); name -> viewports
SHOTS = {"home": ["1440", "390"], "home-empty": ["390"], "verify": ["390"], "remote-consent": ["390"],
         "connection-token": ["390"], "refusal-limited": ["390"]}
MEASURE = r"""(()=>{
const vis=e=>{const r=e.getBoundingClientRect(),s=getComputedStyle(e);return r.width>0&&r.height>0&&s.visibility!=='hidden'&&s.display!=='none'};
const name=e=>(e.getAttribute('aria-label')||e.textContent||e.value||'').trim().replace(/\s+/g,' ').slice(0,40);
document.querySelectorAll('details').forEach(d=>d.open=true);
const iw=innerWidth,de=document.documentElement;
const sel='a[href],button,input:not([type=hidden]),select,textarea,summary';
const els=[...document.querySelectorAll(sel)].filter(vis);
const rect=e=>{const r=e.getBoundingClientRect();return {t:r.top+scrollY,b:r.bottom+scrollY,l:r.left,r:r.right,h:r.height,w:r.width}};
const info=els.map(e=>{
  const r=rect(e),lh=parseFloat(getComputedStyle(e).lineHeight)||0;
  const lines=e.tagName==='A'?new Set([...e.getClientRects()].map(x=>Math.round(x.top))).size:0;
  const range=document.createRange();range.selectNodeContents(e);
  const tl=new Set([...range.getClientRects()].map(x=>Math.round(x.top))).size;
  return {tag:e.tagName.toLowerCase(),cls:(e.className||'').toString(),name:name(e),h:+r.h.toFixed(1),w:+r.w.toFixed(1),textLines:tl,_r:r};
});
// overlap + min vertical gap to nearest other interactive box
info.forEach((a,i)=>{a.overlap=[];a.gap=null;
  info.forEach((b,j)=>{if(i===j)return;
    const A=a._r,B=b._r,ox=Math.min(A.r,B.r)-Math.max(A.l,B.l),oy=Math.min(A.b,B.b)-Math.max(A.t,B.t);
    if(ox>0.5&&oy>0.5&&!(a.cls.includes('skip')||b.cls.includes('skip')))a.overlap.push(b.tag+':'+b.name);
    if(ox>0.5&&oy<=0.5){const g=Math.max(A.t-B.b,B.t-A.b);if(g>=-0.5&&(a.gap===null||g<a.gap))a.gap=+g.toFixed(1)}
  })});
info.forEach(a=>delete a._r);
const over=[...document.querySelectorAll('body *')].filter(vis).filter(e=>{const r=e.getBoundingClientRect();return r.right>iw+1||r.left<-1}).map(e=>e.tagName.toLowerCase()+'.'+e.className);
const skip=document.querySelector('a[href="#main"]');
const bn={};document.querySelectorAll('button').forEach(b=>{if(!vis(b))return;const k=name(b);bn[k]=(bn[k]||0)+1});
const dup=Object.keys(bn).filter(k=>bn[k]>1);
const lone=[...document.querySelectorAll('details.agent>summary')].filter(vis).filter(s=>{const h=s.querySelector('h3');return h&&h.getBoundingClientRect().top-s.getBoundingClientRect().top-parseFloat(getComputedStyle(s).paddingTop)>8}).length;
return JSON.stringify({iw,sw:de.scrollWidth,docH:de.scrollHeight,info,overflow:over,skipPos:skip?getComputedStyle(skip).position:null,
 mains:document.querySelectorAll('main').length,h1:document.querySelectorAll('h1').length,dupButtons:dup,loneMarker:lone});
})()"""

def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + " " + p.stderr[:200])
    return json.loads(p.stdout)["result"]

def shot(name, vp, h, docH):
    n = max(1, -(-docH // h)); out = []
    for i in range(n):
        run("eval", "--expression", f"scrollTo(0,{min(i*h, max(0, docH-h))});scrollY")
        d = base64.b64decode(run("screenshot")["data"])
        rel = f"png/{name}-{vp}" + ("" if n == 1 else f"-s{i+1}of{n}") + ".png"
        open(os.path.join(OUT, rel), "wb").write(d)
        out.append({"path": rel, "viewport": vp, "sha256": hashlib.sha256(d).hexdigest(), "bytes": len(d)})
    run("eval", "--expression", "scrollTo(0,0);1")
    return out

res, manifest, viol = {}, [], []
os.makedirs(os.path.join(OUT, "png"), exist_ok=True)
for fn in sorted(os.listdir(FX)):
    if not fn.endswith(".html"):
        continue
    n = fn[:-5]
    url = "file:///" + os.path.abspath(os.path.join(FX, fn)).replace("\\", "/")
    for vp, (w, h) in VPS.items():
        run("goto", "--url", url)
        for _ in range(4):  # viewport emulation applies asynchronously after goto; never accept a wrong width
            run("exec", "--command", f"set viewport {w} {h}")
            m = json.loads(run("eval", "--expression", MEASURE)["result"])
            if m["iw"] == w:
                break
        assert m["iw"] == w, (fn, vp, m["iw"])
        res[f"{n}@{vp}"] = m
        for e in m["info"]:
            if e["h"] < 43.5 and "skip" not in e["cls"]:
                viol.append(f"{n}@{vp} small {e['tag']}:{e['name']} h={e['h']}")
            if e["tag"] == "a" and e["textLines"] > 1 and "skip" not in e["cls"]:
                viol.append(f"{n}@{vp} wrap {e['name']} lines={e['textLines']}")
            if e["overlap"]:
                viol.append(f"{n}@{vp} overlap {e['name']} x {e['overlap']}")
        if m["mains"] != 1 or m["h1"] != 1 or not m["skipPos"]:
            viol.append(f"{n}@{vp} landmarks main={m['mains']} h1={m['h1']} skip={m['skipPos']}")
        if m["dupButtons"]:
            viol.append(f"{n}@{vp} duplicate button names {m['dupButtons']}")
        if m["loneMarker"]:
            viol.append(f"{n}@{vp} summary marker alone x{m['loneMarker']}")
        if m["overflow"]:
            viol.append(f"{n}@{vp} overflow {m['overflow'][:3]}")
        if n in SHOTS and vp in SHOTS[n]:
            manifest += shot(n, vp, h, m["docH"])
json.dump(res, open(os.path.join(OUT, "measurements.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
json.dump(manifest, open(os.path.join(OUT, "manifest.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
links = sum(1 for m in res.values() for e in m["info"] if e["tag"] == "a")
print("keys", len(res), "links", links, "png", len(manifest), "violations", len(viol))
for v in viol[:40]:
    print(" ", v)
