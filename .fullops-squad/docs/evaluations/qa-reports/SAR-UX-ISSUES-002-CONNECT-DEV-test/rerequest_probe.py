"""Loopback check of the running fixture's re-request paths (blocked, receive limit, unknown) as synthetic account a.

Reuses remote.py's verified SSH session and task clone; counts A's '입력한 코드 …abcd' rows before and after one
same-code re-request per path, and checks D/E pages do not change. Codes stay on the remote host.
"""
import pathlib
import sys

sys.path.insert(0, str(pathlib.Path(__file__).parent))
import remote

SCRIPT = r'''
import json, re, sys, urllib.parse, urllib.request
m = json.load(open(sys.argv[1])); o = m["origin"]; assert o.startswith("http://127.0.0.1:")
def session(name):
    # The __Host- session cookie is Secure, so a cookie jar drops it over loopback http; carry it by hand.
    op = urllib.request.build_opener()
    op.addheaders = [("Origin", o), ("Sec-Fetch-Site", "same-origin")]
    r = op.open(urllib.request.Request(o + "/fixture/as/" + name, method="POST"))
    cookie = [h.split(";")[0] for h in r.headers.get_all("Set-Cookie") if h.startswith("__Host-")]
    assert cookie, "no session cookie"
    op.addheaders.append(("Cookie", cookie[0])); return op
def page(op): return op.open(o + "/home/friends").read().decode()
def plain(html): return re.sub(r'name="(attempt|csrf)" value="[^"]*"', "", html)
a, d, e = session("a"), session("d"), session("e")
failures = 0
for key, code in m["rerequest_codes"].items():
    tail = "입력한 코드 …" + code[-4:]
    before, db, eb = page(a), plain(page(d)), plain(page(e))
    a.open(o + "/home/friends/request", urllib.parse.urlencode({"code": code, "consent": "yes"}).encode())
    after = page(a)
    rows = (before.count(tail), after.count(tail))
    same = plain(page(d)) == db and plain(page(e)) == eb
    ok = rows[0] == 1 and rows[1] == 2 and same  # the pending row also names the tail in its cancel button
    failures += not ok
    print(key, "rows(label mentions) before/after", rows, "targets unchanged", same, "OK" if ok else "FAIL")
print("FAILURES", failures); sys.exit(1 if failures else 0)
'''


def main():
    sha = remote.sha_head()
    base = (remote.HERE / f'{sha}.base').read_text().strip()
    client, values = remote.connect()
    try:
        with client.open_sftp() as sftp, sftp.open(f'{base}/rerequest_probe_remote.py', 'w') as f:
            f.write(SCRIPT)
        stdin, stdout, stderr = client.exec_command(f'python3 {base}/rerequest_probe_remote.py {base}/fixture.json', timeout=120)
        stdin.close()
        out = stdout.read().decode() + stderr.read().decode()
        code = stdout.channel.recv_exit_status()
        (remote.HERE / f'{sha}-fixture-rerequest.log').write_text(out + f'\nEXIT={code}\n', encoding='utf-8')
        print(out + f'EXIT={code}')
        return code
    finally:
        client.close()


if __name__ == '__main__':
    raise SystemExit(main())
