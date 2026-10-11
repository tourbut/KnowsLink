"""Keyboard/focus/contrast probe for SAR-SERVICE-UI-001-VISUAL (Orca embedded browser, synthetic fixtures).
Usage: python3 -I keyboard.py <fixture_dir> <out_dir>. Writes keyboard.json and focus PNG files; never prints base64."""
import base64, json, os, subprocess, sys

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
FX, OUT = sys.argv[1], sys.argv[2]


def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + p.stderr[:200])
    return json.loads(p.stdout)["result"]


def ev(js):
    return json.loads(run("eval", "--expression", js)["result"])


def kb_start():
    """Give the embedded page input focus (keypress is ignored while document.hasFocus() is false),
    then park sequential focus before the skip link with a tabindex=-1 sentinel."""
    run("exec", "--command", "click h1")
    ev("(()=>{const s=document.createElement('span');s.tabIndex=-1;document.body.prepend(s);s.focus();return JSON.stringify(document.hasFocus())})()")


def open_page(name, w, h):
    run("goto", "--url", "file:///" + os.path.abspath(os.path.join(FX, name + ".html")).replace("\\", "/"))
    run("exec", "--command", f"set viewport {w} {h}")


ACTIVE = """JSON.stringify((()=>{const e=document.activeElement;const s=getComputedStyle(e);const r=e.getBoundingClientRect();
return {tag:e.tagName.toLowerCase(),text:(e.getAttribute('aria-label')||e.textContent||'').trim().replace(/\\s+/g,' ').slice(0,36),
outline:s.outlineWidth+' '+s.outlineStyle+' '+s.outlineColor,offset:s.outlineOffset,y:Math.round(r.top+scrollY),h:Math.round(r.height),
inDanger:!!e.closest('.danger,.danger-zone,details.danger')}})())"""

res = {}
# 1) tab order on home at 390 and 1440 (first 40 stops), focus ring on each
for w, h in ((1440, 900), (390, 844)):
    open_page("home", w, h)
    kb_start()
    stops = []
    for i in range(40):
        run("keypress", "--key", "Tab")
        s = ev(ACTIVE)
        stops.append(s)
        if i in (0,):
            r = run("screenshot")
            open(os.path.join(OUT, f"png/focus-home-{w}-skiplink.png"), "wb").write(base64.b64decode(r["data"]))
    res[f"tab-order-home@{w}"] = stops
    ys = [s["y"] for s in stops if s["tag"] != "body"]
    res[f"tab-order-home@{w}-ysorted"] = ys == sorted(ys) if w == 390 else "n/a (two-column layout)"

# 2) skip link -> main, summary toggle via Enter and Space, danger summary focus colour
open_page("home", 390, 844)
kb_start()
run("keypress", "--key", "Tab")
run("keypress", "--key", "Enter")
res["skip-enter-hash"] = ev("JSON.stringify({hash:location.hash,active:document.activeElement.tagName+'#'+document.activeElement.id})")
open_page("home", 390, 844)
kb_start()
ev("(()=>{const s=document.querySelector('#agents details>summary');s.focus();return JSON.stringify(1)})()")
before = ev("JSON.stringify(document.querySelector('#agents details').open)")
run("keypress", "--key", "Enter")
after_enter = ev("JSON.stringify(document.querySelector('#agents details').open)")
run("keypress", "--key", "Space")
after_space = ev("JSON.stringify(document.querySelector('#agents details').open)")
res["summary-toggle"] = {"before": before, "after_enter": after_enter, "after_space": after_space}
run("keypress", "--key", "Enter")  # reopen
ev("(()=>{const d=document.querySelector('#agents details details>summary');d.focus();return JSON.stringify(1)})()")
res["danger-summary-focus"] = ev(ACTIVE)
run("eval", "--expression", "document.activeElement.scrollIntoView({block:'center'});1")
r = run("screenshot")
open(os.path.join(OUT, "png/focus-danger-summary-390.png"), "wb").write(base64.b64decode(r["data"]))
run("keypress", "--key", "Enter")  # open danger zone
run("keypress", "--key", "Tab")
res["danger-first-button-focus"] = ev(ACTIVE)
r = run("screenshot")
open(os.path.join(OUT, "png/focus-danger-button-390.png"), "wb").write(base64.b64decode(r["data"]))

# 3) contrast of rendered colours (WCAG relative luminance on computed styles)
CONTRAST = """JSON.stringify((()=>{
const lum=c=>{const m=c.match(/[\\d.]+/g).map(Number);const a=m.slice(0,3).map(v=>{v/=255;return v<=.03928?v/12.92:Math.pow((v+.055)/1.055,2.4)});return .2126*a[0]+.7152*a[1]+.0722*a[2]};
const bgOf=e=>{while(e){const c=getComputedStyle(e).backgroundColor;const m=c.match(/[\\d.]+/g);if(m&&(m.length<4||+m[3]>0))return c;e=e.parentElement}return 'rgb(255,255,255)'};
const cr=(a,b)=>{const x=lum(a),y=lum(b);return +((Math.max(x,y)+.05)/(Math.min(x,y)+.05)).toFixed(2)};
const out={};
const pick=(label,sel)=>{const e=document.querySelector(sel);if(!e)return;const s=getComputedStyle(e);out[label]={fg:s.color,bg:bgOf(e),ratio:cr(s.color,bgOf(e)),size:s.fontSize}};
pick('body','main p, main');pick('muted','.muted');pick('primary-button','button.primary');pick('danger-button','.danger-button');pick('badge-ok','.badge.ok, .ok');pick('badge-warn','.warn');pick('badge-off','.off');pick('link','main a');pick('note','.note');pick('problem','.problem');pick('code','code');
const inp=document.querySelector('input:not([type=hidden]),select');if(inp){const s=getComputedStyle(inp);out['input-border']={c:s.borderTopColor,bg:bgOf(inp),ratio:cr(s.borderTopColor,bgOf(inp))}}
return out})())"""
for n in ("home", "start-both-problem", "device-prepared", "connection-token"):
    open_page(n, 1440, 900)
    res[f"contrast-{n}"] = ev(CONTRAST)

json.dump(res, open(os.path.join(OUT, "keyboard.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
print("ok", len(res))
