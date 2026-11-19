"""Read-only public check of the 14-day backup sentence on /privacy (KO/EN), /terms and /; no credentials, bodies and cookie values not kept."""
import datetime, http.client, json, pathlib, re, time
KO = '백업은 만든 날부터 14일이 지나면 삭제합니다.'
EN = 'We delete backups 14 days after they are made.'
OLD = ['정기적으로 교체', 'replaced regularly']
rep = {'utc': datetime.datetime.now(datetime.timezone.utc).isoformat(), 'checks': []}
def call(path, headers=None):
    c = http.client.HTTPSConnection('link.knowslog.com', timeout=20)
    c.request('GET', path, headers=headers or {})
    r = c.getresponse(); body = r.read().decode('utf-8', 'replace'); hdr = r.getheaders(); c.close(); time.sleep(1.5)
    return r.status, body, hdr
def add(name, ok, **kw):
    rep['checks'].append({'name': name, 'pass': bool(ok), **kw})
s, t, h = call('/privacy')
add('privacy ko', s == 200 and KO in t and not any(o in t for o in OLD), status=s, has_new=KO in t, has_old=any(o in t for o in OLD))
# English via the language switch the page offers: /lang/en?back=/privacy sets a cookie and redirects back
s, _, h = call('/lang/en?back=/privacy')
loc = next((v for k, v in h if k.lower() == 'location'), None)
ck = next((v.split(';')[0] for k, v in h if k.lower() == 'set-cookie'), None)
add('lang/en redirect', s in (302, 303) and loc in ('/privacy', '/privacy?lang=en') and bool(ck), status=s, location=loc, cookie_set=bool(ck))
s, t, h = call('/privacy', {'Cookie': ck} if ck else {})
add('privacy en', s == 200 and EN in t and not any(o in t for o in OLD), status=s, has_new=EN in t, has_old=any(o in t for o in OLD))
s, t, h = call('/terms'); add('terms', s == 200, status=s)
s, t, h = call('/'); add('landing', s == 200, status=s)
for p in ['/privacy/', '/privacyx']:
    s, _, _ = call(p); add('neighbour ' + p, s == 404, status=s)
rep['all_pass'] = all(c['pass'] for c in rep['checks'])
pathlib.Path('smoke.json').write_text(json.dumps(rep, ensure_ascii=False, indent=1), encoding='utf-8')
print(json.dumps(rep, ensure_ascii=False, indent=1))
