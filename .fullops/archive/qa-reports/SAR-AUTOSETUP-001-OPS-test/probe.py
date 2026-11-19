"""Read-only public probe for SAR-AUTOSETUP-001-OPS: no credentials, no account state change. Usage: probe.py <label>"""
import datetime, hashlib, http.client, json, pathlib, re, sys, time
label = sys.argv[1]
root = pathlib.Path(__file__).parent
dev = pathlib.Path('D:/workspace/KnowsLink/.devtools')
rel = {'0.1.1': dev / 'autosetup-001-ops/release/0.1.1', '0.1.0': dev / 'connect-simple-001-ops/release/0.1.0'}
rep = {'label': label, 'utc': datetime.datetime.now(datetime.timezone.utc).isoformat(), 'checks': []}
def call(method, path, body=None, host='link.knowslog.com'):
    c = http.client.HTTPSConnection(host, timeout=20)
    c.request(method, path, body=body, headers={'Content-Type': 'application/json'} if body else {})
    r = c.getresponse(); data = r.read(); hdr = {k.lower(): v for k, v in r.getheaders()}; c.close(); time.sleep(1.5)
    return r.status, data, hdr
def add(name, ok, **kw):
    rep['checks'].append({'name': name, 'pass': bool(ok), **kw})
inst = {'knowslink-setup.ps1': 'text/plain; charset=utf-8', 'knowslink-claude-plugin.zip': 'application/zip', 'knowslink-codex-plugin.zip': 'application/zip'}
got = {}
for version in ('0.1.1', '0.1.0'):
    for name, ctype in inst.items():
        s, data, h = call('GET', f'/install/{version}/{name}')
        local = (rel[version] / name).read_bytes()
        got[(version, name)] = data
        add(f'install {version} {name}', s == 200 and h.get('content-type') == ctype and data == local, status=s, content_type=h.get('content-type'), bytes=len(data), sha256=hashlib.sha256(data).hexdigest(), matches_release=data == local)
# the installer pins the ZIP hashes of its own build
for version in ('0.1.1', '0.1.0'):
    script = got[(version, 'knowslink-setup.ps1')].decode('utf-8')
    sums = dict(re.findall(r"(claude|codex) = '([0-9a-f]{64})'", script))
    add(f'installer {version} pins zip sha256 + version/base url', sums.get('claude') == hashlib.sha256(got[(version, 'knowslink-claude-plugin.zip')]).hexdigest() and sums.get('codex') == hashlib.sha256(got[(version, 'knowslink-codex-plugin.zip')]).hexdigest() and f"$Version = '{version}'" in script and f"/install/{version}'" in script, claude=sums.get('claude'), codex=sums.get('codex'))
for path in ['/install/0.1.1/other.txt', '/install/0.1.2/knowslink-setup.ps1', '/install/0.1/knowslink-setup.ps1', '/install/0.1.1/', '/install/', '/install/0.1.1/a/knowslink-setup.ps1', '/install/0.1.1/knowslink-setup.PS1', '/install/0.1.1/..%2fknowslink-setup.ps1', '/install/0.1.1/knowslink-evil.zip']:
    s, _, _ = call('GET', path); add('install neighbour ' + path, s == 404, status=s)
s, data, h = call('GET', '/')
add('landing', s == 200, status=s)
s, _, h = call('GET', '/home/agents')
add('/home/agents anonymous redirect', s in (302, 303), status=s, location=h.get('location'))
s, _, _ = call('GET', '/privacy'); add('privacy', s == 200, status=s)
s, _, _ = call('GET', '/terms'); add('terms', s == 200, status=s)
# a well-formed start request with a bad proof is refused as invalid_signature; nothing is stored
body = json.dumps({'token': 'A' * 43, 'client': 'node-local', 'kid': 'probe_kid', 'public': 'AAAA', 'proof': 'AAAA', 'host': 'claude'})
s, data, _ = call('POST', '/v1/connect/start', body)
try:
    code = json.loads(data).get('error')
except Exception:
    code = None
add('connect/start host field', code == 'invalid_signature', status=s, error=code)
s, _, _ = call('GET', '/install/../v1/owners'); add('owner api stays blocked', s in (400, 404), status=s)
for path in ['/owner', '/v1/owners', '/healthz']:
    s, _, _ = call('GET', path); add('blocked ' + path, s == 404, status=s)
rep['all_pass'] = all(c['pass'] for c in rep['checks'])
(root / f'probe-{label}.json').write_text(json.dumps(rep, ensure_ascii=False, indent=1), encoding='utf-8')
print(json.dumps({'label': label, 'all_pass': rep['all_pass'], 'failed': [c['name'] for c in rep['checks'] if not c['pass']], 'n': len(rep['checks'])}))
