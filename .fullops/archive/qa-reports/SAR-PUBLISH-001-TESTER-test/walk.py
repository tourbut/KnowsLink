"""Independent walk of the SAR-PUBLISH-001 synthetic fixture (runs on the remote loopback via explore.py).

Prints one JSON object {checks:[{id,pass,detail}], pages:{...}}. No cookie values, no real Google account.
"""
import json, os, re, urllib.request, urllib.error, urllib.parse

B = os.environ['BASE']; m = json.load(open(B + '/fixture.json')); U = m['url']
G_PEND, G_RES = m['gate_pending'], m['gate_with_result']
checks, pages = [], {}


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *a, **k): return None


op = urllib.request.build_opener(NoRedirect)


def req(method, path, headers=None, form=None):
    data = urllib.parse.urlencode(form).encode() if form is not None else (b'' if method == 'POST' else None)
    h = dict(headers or {})
    if form is not None: h['Content-Type'] = 'application/x-www-form-urlencoded'
    try:
        x = op.open(urllib.request.Request(U + path, method=method, headers=h, data=data))
        return x.status, x.read().decode(), x.headers
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode(), e.headers


def chk(i, ok, detail=''): checks.append({'id': i, 'pass': bool(ok), 'detail': str(detail)[:300]})
def state(): return json.loads(req('GET', '/fixture/state')[1])
def post(path): return json.loads(req('POST', path)[1])
def text(h): return re.sub(r'\s+', ' ', re.sub(r'<[^>]+>', ' ', re.sub(r'<(script|style)[^>]*>.*?</\1>', ' ', h, flags=re.S)))
def login(lang):
    s, _, h = req('POST', f'/fixture/as/{lang}')
    c = [v for k, v in h.items() if k.lower() == 'set-cookie' and 'kl_session' in v]
    return {'Cookie': c[0].split(';')[0], 'Sec-Fetch-Site': 'same-origin'}


RECORD = re.compile(r'^\d{4}/\d\d/\d\d \d\d:\d\d:\d\d calendar cleanup: connections=(\d+) revoked=(\d+) failed=(\d+)$')
# 1. cleanup: three cases from a freshly seeded pre-cleanup state, each then run a second time
req('POST', '/fixture/seed')
s0 = state()
chk('CL-00 seed state', s0['connections'] == 2 and s0['lookup_in_state'] and s0['disclosure_in_state'] and s0['calendar_key_in_state'] and s0['gate_with_result'] == 'approved' and s0['gate_pending'] == 'pending', s0)
for case, q, want in (('key=present', '?key=present', (2, 2, 0)), ('key=present revoke=fail', '?key=present&revoke=fail', (2, 0, 2)), ('key=absent', '?key=absent', (2, 0, 2))):
    if case != 'key=present': req('POST', '/fixture/seed')
    base_req = state()['revoke_requests']
    r1 = post('/fixture/retire' + q)
    mt = RECORD.match(r1['cleanup_record'])
    chk(f'CL {case} record', mt and tuple(map(int, mt.groups())) == want, r1['cleanup_record'])
    chk(f'CL {case} state', r1['connections'] == 0 and not r1['lookup_in_state'] and not r1['disclosure_in_state'] and not r1['calendar_key_in_state'] and r1['gate_with_result'] == 'denied', {k: v for k, v in r1.items() if k != 'cleanup_record'})
    chk(f'CL {case} revoke requests', r1['revoke_requests'] - base_req == (0 if case == 'key=absent' else 2), r1['revoke_requests'] - base_req)
    chk(f'CL {case} no secrets in record', '@' not in r1['cleanup_record'] and not re.search(r'[A-Za-z0-9_\-]{24,}', r1['cleanup_record']), r1['cleanup_record'])
    pend_after = r1['gate_pending']
    r2 = post('/fixture/retire' + q)
    chk(f'CL {case} 2nd run idempotent', r2['cleanup_record'] == '' and r2['revoke_requests'] == r1['revoke_requests'] and r2['connections'] == 0 and r2['gate_with_result'] == 'denied' and r2['gate_pending'] == pend_after, r2)
chk('CL pending gate without lookup untouched (spec: only gates with lookup are denied)', pend_after == 'pending', pend_after)

# 2. approval screens, calendar words, routes (state is cleaned now; reseed first so the approved gate is also seen post-cleanup)
req('POST', '/fixture/seed'); post('/fixture/retire?key=present')
WORDS = re.compile(r'calendar|캘린더|busy time|바쁜 시간|freebusy|Google Calendar', re.I)
for lang in ('ko', 'en'):
    ck = login(lang)
    acct = req('GET', '/home/account', ck)[1]; pages[f'account-{lang}'] = acct
    chk(f'ST {lang} account: no calendar words', not WORDS.search(text(acct)), WORDS.findall(text(acct)))
    chk(f'ST {lang} account: policy links', 'href="/privacy"' in acct and 'href="/terms"' in acct)
    pend = req('GET', f'/home/gates/{G_PEND}', ck)[1]; pages[f'gate-pending-{lang}'] = pend
    chk(f'GT {lang} pending: approve button + no calendar words', 'value="approve"' in pend and not WORDS.search(text(pend)), WORDS.findall(text(pend)))
    den = req('GET', f'/home/gates/{G_RES}', ck)[1]; pages[f'gate-cleaned-{lang}'] = den
    chk(f'GT {lang} cleaned gate: no approve, no calendar words', 'value="approve"' not in den and not WORDS.search(text(den)), WORDS.findall(text(den)))
    for p in ('connect', 'disconnect'):
        chk(f'RT {lang} POST /home/calendar/{p} = 404', req('POST', f'/home/calendar/{p}', ck)[0] == 404, req('POST', f'/home/calendar/{p}', ck)[0])
        chk(f'RT {lang} GET /home/calendar/{p} not 2xx', req('GET', f'/home/calendar/{p}', ck)[0] >= 400)
    chk(f'RT {lang} logged-out POST connect 404', req('POST', '/home/calendar/connect', {'Sec-Fetch-Site': 'same-origin'})[0] in (404, 401, 403))
cons = req('POST', f'/fixture/consume?gate={G_RES}')
chk('GT cleaned gate consume -> invalid_gate (no disclosure)', 'invalid_gate' in cons[1] and 'disclosure":true' not in cons[1], cons[1][:200])

# 3. schedule.query approval: approve the pending gate, then consume; disclosure must be false
ck = login('ko')
pend = req('GET', f'/home/gates/{G_PEND}', ck)[1]
form = {k: v for k, v in re.findall(r'<input[^>]*type="hidden"[^>]*name="([^"]+)"[^>]*value="([^"]*)"', pend)}
form.update({k: v for v, k in re.findall(r'<input[^>]*type="hidden"[^>]*value="([^"]*)"[^>]*name="([^"]+)"', pend)})
form['decision'] = 'approve'
a = req('POST', f'/home/gates/{G_PEND}', ck, form)
chk('GT approve schedule.query accepted (2xx/3xx)', a[0] < 400, (a[0], form.keys()))
after = req('GET', f'/home/gates/{G_PEND}', ck)[1]
chk('GT after approval: no approve button, no calendar words', 'value="approve"' not in after and not WORDS.search(text(after)), WORDS.findall(text(after)))
cons = req('POST', f'/fixture/consume?gate={G_PEND}')
chk('GT approved schedule.query consume: disclosure false', '"disclosure":false' in cons[1].replace(' ', '') or 'disclosure' not in cons[1], cons[1][:300])
pages['consume-approved'] = cons[1]

# 4. policy pages
for path in ('/privacy', '/terms'):
    for lang in ('ko', 'en'):
        hd = {'Cookie': '__Host-kl_lang=en'} if lang == 'en' else {}
        s_, h, _ = req('GET', path, hd); pages[f'policy-{path[1:]}-{lang}'] = h
        t = text(h)
        chk(f'PC {path} {lang} 200 public', s_ == 200, s_)
        chk(f'PC {path} {lang} no braces', '{' not in t and '}' not in t)
        if lang == 'en': chk(f'PC {path} en no Hangul', not re.search('[가-힣]', t.replace('한국어', '')), re.findall('[가-힣]+', t)[:5])
        other = '/terms' if path == '/privacy' else '/privacy'
        chk(f'PC {path} {lang} cross link', f'href="{other}"' in h)
        chk(f'PC {path} {lang} one h1 + title', h.count('<h1>') == 1 and '<title>' in h)
for lang, hd in (('ko', {}), ('en', {'Cookie': '__Host-kl_lang=en'})):
    land = req('GET', '/', hd)[1]; pages[f'landing-{lang}'] = land
    chk(f'PC landing {lang} footer links', 'href="/privacy"' in land and 'href="/terms"' in land)
    # the #start policy line renders only when Google or email login is enabled; the fixture has neither, so unit TestPolicyLinks covers it
    chk(f'PC landing {lang} #start fallback without login (line hidden by design)', 'start-h' in land and land.count('href="/privacy"') == 1)
ck = login('ko')
chk('PC logged-in /privacy 200', req('GET', '/privacy', ck)[0] == 200)
print(json.dumps({'checks': checks, 'pages': pages}, ensure_ascii=False))

