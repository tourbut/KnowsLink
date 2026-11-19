"""Verify CSS fixture provenance and unchanged product-gate inputs; preserve bounded packet evidence."""
import hashlib
import json
import pathlib
import re
import subprocess
import urllib.request

ROOT = pathlib.Path(__file__).resolve().parents[6]
OUT = pathlib.Path(__file__).parent
DEV = ROOT / '.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume'
PRIVATE = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-css-dev')
KEY = 'SAR-LANDING-GUIDE-001-TESTER'
SOURCE = 'f2b8cb622dc397034423582978e77560e0355771'
BASE = 'fc4ad58b65471e2de298cb9f1aba5d2253ea8135'
HEAD = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()


def digest(data):
    return hashlib.sha256(data).hexdigest()


def git(*args):
    return subprocess.check_output(['git', '-C', str(ROOT), *args])


checks = []


def check(name, passed, **evidence):
    checks.append(dict(name=name, passed=passed, **evidence))


manifest = json.loads((DEV / 'fixture-manifest.json').read_text(encoding='utf-8'))
for item in manifest['sources']:
    source = git('show', SOURCE + ':' + item['file'])
    current = git('show', HEAD + ':' + item['file'])
    check('source:' + item['file'], source == current and digest(source) == item['sha256'], sha256=digest(source))
member = git('show', SOURCE + ':internal/relay/member_pages.go').decode().split('const memberStyle = `', 1)[1].split('`', 1)[0]
landing = git('show', SOURCE + ':internal/relay/landing_pages.go').decode()
landing_css = '<style>' + landing.split('<style>', 1)[1].split('</style>', 1)[0] + '</style>'
for name in ['start-google.html', 'human-home-empty.html', 'human-home.html', 'human-account.html']:
    local = (PRIVATE / 'fixtures' / name).read_bytes()
    with urllib.request.urlopen('http://127.0.0.1:18896/' + name, timeout=10) as response:
        served = response.read()
    expected = next(x['sha256'] for x in manifest['fixtures'] if x['file'] == name)
    html = served.decode()
    check('fixture:' + name, local == served and digest(served) == expected, sha256=digest(served), bytes=len(served))
    check('source-style:' + name, member in html and (name != 'start-google.html' or landing_css in html))
paths = ['cmd', 'internal', 'adapters', 'db', 'deploy', 'scripts', 'go.mod', 'go.sum', 'Makefile', 'Dockerfile', 'compose.yaml', 'sqlc.yaml', '.dockerignore', '.fullops-squad/lint/lint.json']
for first in [SOURCE, BASE]:
    result = subprocess.run(['git', '-C', str(ROOT), 'diff', '--exit-code', first, HEAD, '--', *paths], capture_output=True, text=True)
    check('gate-inputs:' + first, result.returncode == 0, to=HEAD, paths=paths, exit_code=result.returncode, diff=result.stdout)
gate = json.loads((PRIVATE / 'final-lint.json').read_text(encoding='utf-8'))
check('prior-gate-sha-and-exits', gate['head'] == BASE and all(c['status'] == 'passed' and c['exit_code'] == 0 for c in gate['commands']), receipt=gate)
for target in ['lint', 'test', 'install']:
    data = (PRIVATE / f'{BASE}-{target}.log').read_bytes()
    check('prior-linux-log:' + target, data.rstrip().endswith(b'EXIT=0'), sha256=digest(data), tail=data.decode('utf-8')[-650:])
old = 'fe2387a7676f7997aea76a2606f4cc2dbcd4b56b'
changed = git('diff', '--name-only', old, SOURCE, '--', *paths).decode().splitlines()
check('old-QA-impact-boundary', set(changed) == {'internal/relay/member_pages.go', 'internal/relay/landing_pages.go', 'internal/relay/human_pages_test.go'}, changed=changed)
packet_path = ROOT / f'.fullops-squad/docs/evaluations/jev/{KEY}-packet.json'
packet = json.loads(packet_path.read_text(encoding='utf-8'))
names = set(packet['context_paths']) | set(packet.get('remaining_context_paths', [])) | set(packet['required_paths'])
names.update(x['path'] for x in packet['items'])
names.update(x['path'] for x in packet['unknown'] if 'path' in x)
rows = []
for name in sorted(names):
    path = ROOT / name
    data = path.read_bytes()
    lines = data.decode('utf-8-sig', 'replace').splitlines()
    hits = [i for i, line in enumerate(lines) if re.search('LANDING|LG-0|공개 첫페이지|native|200%|not_run|FAIL|미검증|CSS|fixture', line)]
    rows.append(dict(path=name, sha256=digest(data), lines=len(lines), scope='bounded relevance inspection, not full historical re-audit', headings=[line for line in lines if line.startswith(('#', 'title:', 'summary:'))][:5], excerpts=[{'line': i + 1, 'text': '\n'.join(lines[max(0, i-1):i+2])} for i in hits[:3]]))
(OUT / 'source-audit.json').write_text(json.dumps(rows, ensure_ascii=False, indent=2), encoding='utf-8')
(OUT / 'identity.json').write_text(json.dumps(dict(head=HEAD, source=SOURCE, checks=checks, packet_sha256=digest(packet_path.read_bytes()), packet_paths=len(rows), scope='Reuse Linux gate at original fc4 SHA; no new Linux/auth/DB execution'), ensure_ascii=False, indent=2), encoding='utf-8')
assert all(c['passed'] for c in checks), [c['name'] for c in checks if not c['passed']]
print(f'PASS {len(checks)} provenance checks; {len(rows)} bounded packet paths')
