"""Read fixed review evidence only; never execute product or browser tests."""
import hashlib
import json
import re
import subprocess
import sys
from pathlib import Path

OUT = Path(__file__).parent
SNAP = Path('C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-9b3748e62cc04148b188a1f0860fe532')
HEAD = 'e939cfbb1939adb4126ad379317c26ec39b692b3'
BASE = 'a3d85ba85f00b9f5d65829aca4195d111e8fd906'
SOURCE = 'f2b8cb622dc397034423582978e77560e0355771'
DEV = 'fc4ad58b65471e2de298cb9f1aba5d2253ea8135'
QA = '.fullops-squad/docs/evaluations/qa-reports/'
TP = QA + 'SAR-LANDING-GUIDE-001-TESTER-test/'
VP = QA + 'SAR-LANDING-GUIDE-001-VISUAL/'
DP = QA + 'SAR-LANDING-GUIDE-001-DEV-test/css-resume/'
checks = []


def git(*args):
    return subprocess.check_output(['git', '--no-pager', '-C', str(SNAP), *args])


def blob(path, ref=HEAD):
    return git('show', ref + ':' + path)


def read(path):
    return json.loads(blob(path))


def digest(data):
    return hashlib.sha256(data).hexdigest()


def check(name, passed, **detail):
    checks.append(dict(name=name, passed=bool(passed), **detail))


paths = ['cmd', 'internal', 'adapters', 'db', 'deploy', 'scripts', 'go.mod', 'go.sum', 'Makefile', 'Dockerfile', 'compose.yaml', 'sqlc.yaml', '.dockerignore', '.github', '.fullops-squad/lint/lint.json']
for ref in [BASE, DEV, SOURCE]:
    changed = git('diff', '--name-only', ref, HEAD, '--', *paths).decode().splitlines()
    check('product-inputs:' + ref, not changed, changed=changed)
for ref in ['7fd58fcdbf731b0079350369ef012b03b1ee8911', 'c6d5a28374a7cdd05c812dcd9fe89efbe5479a58', '5149ac1']:
    exact = git('rev-parse', ref).decode().strip()
    rc = subprocess.run(['git', '--no-pager', '-C', str(SNAP), 'merge-base', '--is-ancestor', exact, HEAD]).returncode
    check('integrated-ancestor:' + ref, rc == 0, exact_sha=exact)

for key in ['SAR-LANDING-GUIDE-001-CSS-REVIEW', 'SAR-UX-ISSUES-002-DESIGN-REVIEW']:
    prefix = QA + key + '-review/'
    r = read(prefix + 'result.json')
    hashes = {n: digest(blob(prefix + n)) for n in ['result.json', 'report.md', 'check-execution.json']}
    execution = read(prefix + 'check-execution.json')
    check('prior-review-preserved:' + key, all(blob(prefix + n) == blob(prefix + n, BASE) for n in hashes), hashes=hashes, head=r['head'], independence=r['independence'], check_execution=execution)

manifest = read(DP + 'fixture-manifest.json')
for item in manifest['sources']:
    check('source-hash:' + item['file'], digest(blob(item['file'])) == item['sha256'] and blob(item['file']) == blob(item['file'], SOURCE))
private = Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-css-dev')
member = blob('internal/relay/member_pages.go').decode().split('const memberStyle = `', 1)[1].split('`', 1)[0]
landing = blob('internal/relay/landing_pages.go').decode().split('<style>', 1)[1].split('</style>', 1)[0]
for name in ['start-google.html', 'human-home-empty.html', 'human-home.html', 'human-account.html']:
    data = (private / 'fixtures' / name).read_bytes()
    expected = next(x['sha256'] for x in manifest['fixtures'] if x['file'] == name)
    check('preserved-fixture:' + name, digest(data) == expected and member in data.decode() and (name != 'start-google.html' or landing in data.decode()), sha256=digest(data))

for rel, count in [(TP + 'css-followup/results.json', 44), (TP + 'css-followup/identity.json', 18), (TP + 'browser-results.json', 22), (TP + 'http-results.json', 34)]:
    d = read(rel)
    check('recorded-checks:' + rel, len(d['checks']) == count and all(x['passed'] for x in d['checks']), count=len(d['checks']))
for rel in [TP + 'css-followup/results.json', TP + 'browser-results.json']:
    d = read(rel)
    for capture in d['captures']:
        p = str(Path(rel).parent).replace('\\', '/') + '/' + capture['file']
        check('capture-hash:' + p, digest(blob(p)) == capture['sha256'])
visual = read(VP + 'css-followup/identity.json')
for capture in visual['captures']:
    check('designer-capture:' + capture['file'], digest(blob(DP + capture['file'])) == capture['sha256'] and capture['directly_viewed'])
results = read(TP + 'css-followup/results.json')
check('fixture-requests', bool(results['requests']) and all(x['method'] == 'GET' and x['url'].startswith('http://127.0.0.1:18896/') for x in results['requests']))
zoom = read(TP + 'css-followup/zoom-results.json')
check('recorded-actual-zoom', zoom['passed'] and zoom['setting'] == '200%' and zoom['dpr'] == 2 and zoom['scale'] == 1 and zoom['cssZoom'] == '1' and zoom['gap'] == 24 and zoom['overflow'] <= zoom['width'], evidence=zoom)
check('original-contrast-failure', any(x['passed'] is False and x['ratio'] < 4.5 for x in read(TP + 'contrast-results.json')))
check('original-tester-report-preserved', blob(TP + 'report.md') == blob(TP + 'report.md', '5149ac1'))
check('original-visual-prefix-preserved', blob(VP + 'report.md').startswith(blob(VP + 'report.md', BASE)))
check('startup-failure-preserved', b'ECONNREFUSED' in blob(TP + 'css-followup/browser-startup-failure.log'))
check('finish-failure-preserved', read(TP + 'css-followup/finish-first-failure.json')['exit_code'] == 1)

receipt = read(VP + 'css-followup/dev-final-lint.json')
check('linux-gate-original', receipt['head'] == DEV and len(receipt['commands']) == 2 and all(c['status'] == 'passed' and c['exit_code'] == 0 for c in receipt['commands']))
for target in ['lint', 'test']:
    original = (private / f'{DEV}-{target}.log').read_bytes()
    copy = (SNAP / (VP + f'css-followup/dev-{DEV}-{target}.log')).read_bytes()
    expected = next(c['sha256'] for c in read(TP + 'css-followup/identity.json')['checks'] if c['name'] == 'prior-linux-log:' + target)
    check('linux-log:' + target, digest(original) == expected and original.replace(b'\r\n', b'\n') == copy.replace(b'\r\n', b'\n') and original.rstrip().endswith(b'EXIT=0'), original_sha256=digest(original), committed_copy_sha256=digest(copy), comparison='Original CRLF bytes match tester receipt hash; Git text copy differs only by CRLF to LF normalization.')

sys.path.insert(0, 'C:/Users/shin/.codex/plugins/cache/fullops-squad/fullops-squad/1.4.4/scripts')
import jev_packet
for role, name in [('tester', 'TESTER'), ('designer', 'VISUAL')]:
    key = 'SAR-LANDING-GUIDE-001-' + name
    packet = read('.fullops-squad/docs/evaluations/jev/' + key + '-packet.json')
    log = blob('.fullops-squad/handovers/logs/2026-10-15_to_' + role + '.md').decode()
    tail = log.split('<!-- fullops-attempt: ' + key + ' ' + packet['attempt'] + ' -->', 1)[1].lstrip()
    try:
        jev_packet.check(SNAP, role, key, completion=True, required=True, head=HEAD, text=tail)
        check('official-packet-completion:' + role, True, attempt=packet['attempt'], partial=packet.get('partial'), unknown_count=len(packet['unknown']))
    except ValueError as error:
        check('official-packet-completion:' + role, False, error=str(error))
    check('inbox-empty:' + role, not blob('.fullops-squad/handovers/to_' + role + '.md').strip())

review = json.loads((OUT / 'result.json').read_text(encoding='utf-8'))
files = []
for row in review['files']:
    data = blob(row['path'])
    if row['path'].endswith('.json'):
        json.loads(data)
    files.append({'path': row['path'], 'status': row['status'], 'sha256': digest(data)})
check('all-diff-paths-accounted', len(files) == len(git('diff', '--name-only', BASE, HEAD).decode().splitlines()), count=len(files))
check('snapshot-clean', not git('status', '--porcelain').strip())
check('snapshot-detached', subprocess.run(['git', '--no-pager', '-C', str(SNAP), 'symbolic-ref', '-q', 'HEAD'], capture_output=True).returncode == 1)
report = {'scope': 'Static evidence checks only; no product/browser rerun', 'base': BASE, 'head': HEAD, 'checks': checks, 'file_hashes': files}
(OUT / 'verification.json').write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
failed = [x for x in checks if not x['passed']]
print(json.dumps({'checks': len(checks), 'failures': failed}, ensure_ascii=False, indent=2))
assert not failed, failed
