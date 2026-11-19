"""Keyboard/focus check for SAR-FRONTEND-POLISH-001-DEV over synthetic fixtures in Orca's embedded browser.

Tabs through home (390x844) and remote-consent, records each focused element and its outline, and opens the
first agent row and its danger zone with Enter. Usage: python3 -I keyboard.py <fixture_dir> <out_json>
"""
import json, os, subprocess, sys

ORCA = r"C:\Users\shin\AppData\Local\Programs\orca\resources\bin\orca.EXE"
FX, OUT = sys.argv[1], sys.argv[2]
ACTIVE = r"""(()=>{const e=document.activeElement,s=getComputedStyle(e);
return JSON.stringify({tag:e.tagName.toLowerCase(),name:(e.getAttribute('aria-label')||e.textContent||'').trim().replace(/\s+/g,' ').slice(0,40),
outline:s.outlineStyle+' '+s.outlineWidth,open:e.tagName==='SUMMARY'?e.parentElement.open:null})})()"""

def run(*a):
    p = subprocess.run([ORCA, *a, "--json"], capture_output=True, text=True, encoding="utf-8")
    if p.returncode:
        raise SystemExit("orca failed: " + " ".join(a[:2]) + " " + p.stdout[:300])
    return json.loads(p.stdout)["result"]

def active():
    return json.loads(run("eval", "--expression", ACTIVE)["result"])

res, fails = {}, []
for name, tabs in (("home", 14), ("remote-consent", 4)):
    run("goto", "--url", "file:///" + os.path.abspath(os.path.join(FX, name + ".html")).replace("\\", "/"))
    run("exec", "--command", "set viewport 390 844")
    run("exec", "--command", "click h1")  # Orca delivers keys only after the page has been clicked once
    run("eval", "--expression", "document.activeElement&&document.activeElement.blur();scrollTo(0,0);1")
    seq = []
    for _ in range(tabs):
        run("keypress", "--key", "Tab")
        a = active()
        seq.append(a)
        if a["tag"] != "body" and not a["outline"].startswith("solid 3px"):
            fails.append(f"{name}: no 3px focus ring on {a}")
        if name == "home" and a["tag"] == "summary" and a["open"] is False:
            run("keypress", "--key", "Enter")
            opened = active()
            seq.append({"enter": opened})
            if not opened["open"]:
                fails.append(f"{name}: Enter did not open {a['name']}")
            break
    if all(x.get("tag") == "body" for x in seq if "tag" in x):
        fails.append(f"{name}: Tab never moved focus off body")
    res[name] = seq
json.dump({"results": res, "failures": fails}, open(OUT, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
print("failures", len(fails))
for f in fails:
    print(" ", f)
