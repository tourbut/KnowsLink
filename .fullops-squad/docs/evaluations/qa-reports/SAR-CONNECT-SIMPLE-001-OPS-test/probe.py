"""Read-only public probe for SAR-CONNECT-SIMPLE-001-OPS: no credentials, no account state change. Usage: probe.py <label> [expect-new]"""
import datetime, hashlib, http.client, json, pathlib, sys, time
label = sys.argv[1]
new = len(sys.argv) > 2 and sys.argv[2] == 'expect-new'
root = pathlib.Path(__file__).parent
rel = root / 'release' / '0.1.0'
rep = {'label': label, 'utc': datetime.datetime.now(datetime.timezone.utc).isoformat(), 'checks': []}
def call(method, path, body=None, host='link.knowslog.com'):
    c = http.client.HTTPSConnection(host, timeout=20)
    c.request(method, path, body=body, headers={'Content-Type': 'application/json'} if body else {})
    r = c.getresponse(); data = r.read(); hdr = {k.lower(): v for k, v in r.getheaders()}; c.close(); time.sleep(1.5)
    return r.status, data, hdr
def add(name, ok, **kw):
    rep['checks'].append({'name': name, 'pass': bool(ok), **kw})
inst = {'knowslink-setup.ps1': 'text/plain; charset=utf-8', 'knowslink-claude-plugin.zip': 'application/zip', 'knowslink-codex-plugin.zip': 'application/zip'}
for name, ctype in inst.items():
    s, data, h = call('GET', '/install/0.1.0/' + name)
    if new:
        add('install ' + name, s == 200 and h.get('content-type') == ctype and hashlib.sha256(data).digest() == hashlib.sha256((rel / name).read_bytes()).digest(), status=s, content_type=h.get('content-type'), bytes=len(data), sha256_matches_release=hashlib.sha256(data).hexdigest() == hashlib.sha256((rel / name).read_bytes()).hexdigest())
    else:
        add('install ' + name + ' (not yet public)', s == 404, status=s)
if new:
    for path in ['/install/0.1.0/other.txt', '/install/0.1.1/knowslink-setup.ps1', '/install/0.1/knowslink-setup.ps1', '/install/0.1.0/', '/install/', '/install/0.1.0/a/knowslink-setup.ps1', '/install/0.1.0/knowslink-setup.PS1', '/install/0.1.0/..%2fknowslink-setup.ps1', '/install/0.1.0/knowslink-evil.zip']:
        s, _, _ = call('GET', path); add('install neighbour ' + path, s == 404, status=s)
s, data, h = call('GET', '/')
add('landing', s == 200, status=s)
s, _, h = call('GET', '/home/agents')
add('/home/agents anonymous redirect', s in (302, 303), status=s, location=h.get('location'))
s, _, _ = call('GET', '/privacy'); add('privacy', s == 200, status=s)
s, _, _ = call('GET', '/terms'); add('terms', s == 200, status=s)
# host field: the old relay refuses the unknown field (invalid_schema, 422); the new one parses it and refuses only the bad proof (invalid_signature). Nothing is stored.
body = json.dumps({'token': 'A' * 43, 'client': 'node-local', 'kid': 'probe_kid', 'public': 'AAAA', 'proof': 'AAAA', 'host': 'grok'})
s, data, _ = call('POST', '/v1/connect/start', body)
try:
    code = json.loads(data).get('error')
except Exception:
    code = None
add('connect/start host field', (code == 'invalid_signature') if new else (code == 'invalid_schema'), status=s, error=code)
s, _, _ = call('GET', '/install/../v1/owners'); add('owner api stays blocked', s in (400, 404), status=s)
for path in ['/owner', '/v1/owners', '/healthz']:
    s, _, _ = call('GET', path); add('blocked ' + path, s == 404, status=s)
rep['all_pass'] = all(c['pass'] for c in rep['checks'])
(root / f'probe-{label}.json').write_text(json.dumps(rep, ensure_ascii=False, indent=1), encoding='utf-8')
print(json.dumps(rep, ensure_ascii=False, indent=1))
