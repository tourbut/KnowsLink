import sys, json
sys.path.insert(0, ".")
from client import *
SCREENS = {
 "public": "/", "verify": "/auth/verify", "home": "/home", "name": "/home/name", "conversation": M["conversation"],
 "friends": "/home/friends", "agents": "/home/agents", "account": "/home/account", "legacy": "/home/legacy", "new": "/home/new", "receipts": "/home/receipts",
 "conn-waiting": "/home/connections/" + ID["connection_waiting"], "conn-prepared": "/home/connections/" + ID["connection_prepared"], "conn-approved": "/home/connections/" + ID["connection_approved"],
 "connect-device": "/connect/" + ID["device"], "gate": "/home/gates/" + ID["gate"], "remote": "/remote/consent/" + ID["remote"],
}
if __name__ == "__main__":
    for who in ["en", "ko"]:
        b = as_(who)
        for k, p in SCREENS.items():
            s, h, body = b.get(p)
            open(f"pages/{who}-{k}.html", "w", encoding="utf-8").write(body)
            t = re.search(r"<title>(.*?)</title>", body); l = re.search(r"<html lang=\"(\w*)\"", body)
            print(who, k, s, h.get("content-language"), l and l.group(1), t and t.group(1), len(body))
