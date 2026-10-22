"""usage: turn.py B|C '<prompt>' '<answers json list>' [timeout]  -- append a turn for the App Server driver and wait for its turnEnd."""
import json, sys, time
n, prompt, answers = sys.argv[1], sys.argv[2], json.loads(sys.argv[3])
timeout = int(sys.argv[4]) if len(sys.argv) > 4 else 240
cmd, logf = f"cmd-{n}.jsonl", f"log-{n}.jsonl"
before = sum(1 for l in open(logf, encoding="utf-8") if '"turnEnd"' in l)
with open(cmd, "a", encoding="utf-8") as f:
    f.write(json.dumps({"turn": prompt, "answers": answers}, ensure_ascii=False) + "\n")
t0 = time.time()
while time.time() - t0 < timeout:
    time.sleep(2)
    ls = [l for l in open(logf, encoding="utf-8")]
    if sum(1 for l in ls if '"turnEnd"' in l) > before:
        break
else:
    print("TIMEOUT")
# print events of this turn
idx = max(i for i, l in enumerate(ls) if '"turn":' in l and '"turnEnd"' not in l and '"event"' not in l)
for l in ls[idx:]:
    d = json.loads(l)
    if "elicitation" in d: print("DIALOG", d.get("knowslinkDialog"), d["answer"], "|", d["elicitation"][:400].replace("\n", " "))
    elif d.get("event") == "item/completed":
        it = d["params"]["item"]
        if it["type"] == "mcpToolCall": print("TOOL", it.get("tool"), json.dumps(it.get("arguments"), ensure_ascii=False)[:200], "->", json.dumps(it.get("result"), ensure_ascii=False)[:700])
        elif it["type"] == "agentMessage": print("AGENT", it.get("text", "")[:500])
    elif d.get("serverRequest"): print("SRVREQ", d["serverRequest"])
    elif "turnEnd" in d: print("TURNEND")
