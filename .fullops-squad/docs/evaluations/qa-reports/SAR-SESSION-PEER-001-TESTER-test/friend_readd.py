"""Re-request and re-accept the A-B friendship after removal. usage: python3 friend_readd.py <relay> <mail-dir> <emailA> <emailB>"""
import sys

relay, mail_dir, ea, eb = sys.argv[1:5]
src = open(__file__.replace("friend_readd.py", "fixture.py"), encoding="utf-8").read().split("RESUME =")[0]
sys.argv = ["fixture.py", relay, mail_dir, "x", "x", "x"]
ns = {"__file__": "fixture.py"}
exec(compile(src, "fixture.py", "exec"), ns)
a, b = ns["login"](ea), ns["login"](eb)
ns["request_friend"](a, ns["friend_code"](b))
ns["accept"](b, "합성 소유자 A")
print("re-added")
