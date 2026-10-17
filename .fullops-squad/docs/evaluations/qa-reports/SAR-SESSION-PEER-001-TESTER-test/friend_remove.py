"""Owner <label> removes the friendship with the first active friend on the isolated relay (web form). usage: python3 friend_remove.py <relay> <mail-dir> <email>"""
import re
import sys

relay, mail_dir, email = sys.argv[1:4]
src = open(__file__.replace("friend_remove.py", "fixture.py"), encoding="utf-8").read().split("RESUME =")[0]
sys.argv = ["fixture.py", relay, mail_dir, "x", "x", "x"]
ns = {"__file__": "fixture.py"}
exec(compile(src, "fixture.py", "exec"), ns)
b = ns["login"](email)
page = ns["expect"](b.do("GET", "/home/friends"), 200)
fid = re.search(r'href="/home/friends/([^"]+)"', page).group(1)
detail = ns["expect"](b.do("GET", f"/home/friends/{fid}"), 200)
gen = re.search(r'name="generation" value="(\d+)"><button class="danger-button" name="decision" value="remove"', detail).group(1)
ns["expect"](b.do("POST", "/home/friends/decision", {"friend": fid, "generation": gen, "decision": "remove"}), 303)
print("removed friend", fid, "generation", gen)
