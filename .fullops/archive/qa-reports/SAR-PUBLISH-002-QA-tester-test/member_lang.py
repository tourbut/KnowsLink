"""Remote loopback fixture: a signed-in member with stored locale en/ko sees /, /privacy, /terms in that language. Run via explore.py."""
import json, os, re, urllib.request, urllib.error
B = os.environ['BASE']; U = json.load(open(B + '/fixture.json'))['url']
class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *a, **k): return None
op = urllib.request.build_opener(NoRedirect)
def req(method, path, cookie=None):
    h = {'Cookie': cookie} if cookie else {}
    try:
        x = op.open(urllib.request.Request(U + path, method=method, headers=h, data=b'' if method == 'POST' else None)); return x.status, x.read().decode(), x.headers
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode(), e.headers
out = []
for loc, want in (('en', 'Privacy Policy'), ('ko', '개인정보처리방침')):
    st, _, h = req('POST', '/fixture/as/' + loc)
    ck = '; '.join(c.split(';')[0] for c in h.get_all('Set-Cookie') or [])
    out.append({'as': loc, 'login_status': st, 'cookie_set': bool(ck)})
    for path, h1 in (('/privacy', want), ('/terms', 'Terms of Service' if loc == 'en' else '서비스 약관'), ('/', None)):
        st, body, hd = req('GET', path, ck)
        m = re.search(r'<h1>(.*?)</h1>', body)
        out.append({'as': loc, 'path': path, 'status': st, 'html_lang': (re.search(r'<html lang="(\w+)"', body) or [None, None])[1], 'content_language': hd.get('Content-Language'), 'h1': m and m.group(1), 'ok': st == 200 and (m and h1 in m.group(1) if h1 else True) and hd.get('Content-Language') == loc})
print(json.dumps(out, ensure_ascii=False, indent=1))
