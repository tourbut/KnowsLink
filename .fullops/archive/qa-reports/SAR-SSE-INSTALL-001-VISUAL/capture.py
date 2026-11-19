"""Capture the synthetic fixed-SHA fixtures through Orca's embedded browser only.
Record actual viewport and PNG dimensions; never click or submit consent.
"""
import base64
import hashlib
import json
from pathlib import Path
import struct
import subprocess
import sys

OUT = Path(__file__).resolve().parent
PAGE = "c1de83f5-2c16-4f6c-ad74-0e508af70166"
ORCA = "C:/Users/shin/AppData/Local/Programs/orca/resources/bin/orca.EXE"
METRICS = """JSON.stringify({url:location.href,innerWidth,innerHeight,scrollY,dpr:devicePixelRatio,
scrollWidth:document.documentElement.scrollWidth,clientWidth:document.documentElement.clientWidth,
scrollHeight:document.documentElement.scrollHeight,bodyStyle:{font:getComputedStyle(document.body).font,
color:getComputedStyle(document.body).color,background:getComputedStyle(document.body).backgroundColor,
overflowWrap:getComputedStyle(document.body).overflowWrap},
buttons:[...document.querySelectorAll('button')].map(e=>({text:e.innerText,
x:e.getBoundingClientRect().x,y:e.getBoundingClientRect().y,width:e.getBoundingClientRect().width,
height:e.getBoundingClientRect().height,visible:!!e.getBoundingClientRect().width})),
links:[...document.links].map(e=>({text:e.innerText,y:e.getBoundingClientRect().y})),
text:document.body.innerText})"""


def call(*args):
    run = subprocess.run([ORCA, *args, "--page", PAGE, "--json"], capture_output=True, text=True, encoding="utf-8")
    response = json.loads(run.stdout)
    if run.returncode != 0 or not response.get("ok"):
        raise RuntimeError((args, run.returncode, response, run.stderr))
    return response["result"]


def screenshot(command, name):
    result = call(command)
    data = base64.b64decode(result["data"])
    assert data[:8] == b"\x89PNG\r\n\x1a\n"
    size = struct.unpack(">II", data[16:24])
    target = OUT / (name + ".png")
    target.write_bytes(data)
    return {"path": target.name, "width": size[0], "height": size[1],
            "sha256": hashlib.sha256(data).hexdigest()}


viewport = sys.argv[1]
bottom = len(sys.argv) > 2 and sys.argv[2] == "bottom"
width, height = {"desktop": (1280, 800), "mobile": (390, 844)}[viewport]
records = []
for fixture in ("login", "consent", "login-long", "consent-long"):
    if bottom and not fixture.endswith("-long"):
        continue
    name = viewport + "-" + fixture + ("-bottom" if bottom else "")
    call("goto", "--url", (OUT / (fixture + ".html")).as_uri())
    call("exec", "--command", f"set viewport {width} {height}")
    if bottom:
        call("scroll", "--direction", "down", "--amount", "2000")
    snapshot = call("snapshot")
    metrics = json.loads(call("eval", "--expression", METRICS)["result"])
    assert (metrics["innerWidth"], metrics["innerHeight"]) == (width, height)
    shot = screenshot("screenshot", name)
    assert (shot["width"], shot["height"]) == (width, height)
    record = {"case": name, "fixture": fixture + ".html", "metrics": metrics,
              "snapshot": snapshot, "capture": shot}
    if not bottom and metrics["scrollHeight"] > height:
        record["full_capture"] = screenshot("full-screenshot", name + "-full")
        # Orca full capture excludes the vertical scrollbar; viewport capture retains it.
        assert record["full_capture"]["width"] == metrics["clientWidth"]
    records.append(record)
(OUT / (viewport + ("-bottom" if bottom else "") + "-observations.json")).write_text(json.dumps(records, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
