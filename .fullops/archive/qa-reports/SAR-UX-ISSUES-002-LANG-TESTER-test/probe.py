import sys, json
sys.path.insert(0, ".")
from client import *
for who in ["en", "ko", "none"]:
    b = as_(who); print(who, b.cookies.keys())
    for p in ["/home/account", "/home", "/home/legacy"]:
        s, h, body = b.get(p)
        print(" ", p, s, h.get("content-language"), re.search(r"<html lang=\"(\w+)\"", body).group(1) if s == 200 else "", re.search(r"<title>(.*?)</title>", body).group(1) if s == 200 else "")
