import hashlib, json, sys, urllib.request
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent)); 
RELAY="http://127.0.0.1:18441"
def tree(p):
    p=Path(p); h=hashlib.sha256()
    for f in sorted(x for x in p.rglob("*") if x.is_file()): h.update(str(f.relative_to(p)).encode()+b"\0"+f.read_bytes())
    return h.hexdigest()[:16]
box, app, member = sys.argv[1:4]
root=Path(__file__).parent/"work"/box; priv=root/"local/KnowsLink"/f"{app}-01-private"
a=json.loads((priv/"agent.json").read_text())
st=json.loads(urllib.request.urlopen(f"{RELAY}/fixture/state?member={member}").read())
print(json.dumps({"agent":a.get("agent"),"private_tree":tree(priv),"keys":sorted(p.name for p in priv.iterdir()),"web_agents":st["agents"],"friends":st["friends"]},ensure_ascii=False))
