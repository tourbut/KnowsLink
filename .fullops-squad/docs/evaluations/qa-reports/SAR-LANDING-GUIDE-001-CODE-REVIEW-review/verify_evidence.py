"""Read fixed review inputs and preserve integrity evidence outside the snapshot."""
import hashlib
import json
from pathlib import Path
import struct
import subprocess

HERE = Path(__file__).parent
SNAP = Path('C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-6599ed38950f468c916585fcd9f1f417')
RAW = Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-dev')
BASE = 'eae9a8bffae706410f684e4c837e396e7b0bcf26'
HEAD = 'fe2387a7676f7997aea76a2606f4cc2dbcd4b56b'
PRODUCT = '4795d619513fde4fbcd4a55b94ec91b056ebbaf0'
QA = '.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test'

def digest(data):
    return hashlib.sha256(data).hexdigest()

def git(*args):
    result = subprocess.run(['git', '--no-pager', '-C', str(SNAP), *args], capture_output=True, text=True, encoding='utf-8')
    return {'args': list(args), 'exit_code': result.returncode, 'stdout': result.stdout, 'stderr': result.stderr}

preview = json.loads((HERE / 'preview.json').read_text(encoding='utf-8'))
files = []
for item in preview['reviewable_files'] + preview['excluded_files']:
    path = SNAP / item['path']
    data = path.read_bytes()
    record = {'path': item['path'], 'bytes': len(data), 'sha256': digest(data)}
    if path.suffix == '.json':
        value = json.loads(data)
        record['json_valid'] = True
        record['top_keys'] = list(value) if isinstance(value, dict) else ['array']
    files.append(record)

captures = json.loads((SNAP / QA / 'browser-audit/captures.json').read_text(encoding='utf-8'))
for record in captures:
    data = (SNAP / QA / record['path']).read_bytes()
    assert digest(data) == record['sha256'], record['path']
    assert struct.unpack('>II', data[16:24]) == (record['png_width'], record['png_height'])

checks = [git('rev-parse', 'HEAD'), git('symbolic-ref', '-q', 'HEAD'), git('status', '--porcelain', '--untracked-files=all'),
          git('diff', '--no-ext-diff', '--no-textconv', '--check', BASE, HEAD),
          git('diff', '--no-ext-diff', '--no-textconv', '--exit-code', PRODUCT, HEAD, '--', 'cmd', 'internal', 'adapters', 'db', 'scripts', 'deploy', 'Makefile', 'go.mod', 'go.sum', 'Dockerfile', 'compose.yaml', 'sqlc.yaml', QA + '/gate'),
          git('diff', '--no-ext-diff', '--no-textconv', '--exit-code', BASE, HEAD, '--', 'AGENTS.md', '.fullops-squad/FULLOPS.md', '.fullops-squad/rules/common', '.fullops-squad/project.md', '.fullops-squad/docs/agents/document-writing.md', '.fullops-squad/review/rule.json', '.fullops-squad/lint')]
assert checks[0]['stdout'].strip() == HEAD
assert checks[1]['exit_code'] == 1 and not checks[1]['stdout']
assert checks[2]['exit_code'] == 0 and not checks[2]['stdout']
assert all(item['exit_code'] == 0 for item in checks[3:])

lint_bytes = (RAW / 'final-lint.json').read_bytes()
lint = json.loads(lint_bytes)
assert lint['head'] == HEAD and lint['base'] == BASE
assert lint['summary']['errors'] == 0 and lint['summary']['unavailable'] == 0
assert all(item['status'] == 'passed' and item['exit_code'] == 0 for item in lint['commands'])
for filename in ['lint-original.json', 'lint.json']:
    (HERE / filename).write_bytes(lint_bytes)
logs = []
for label in ['install', 'lint', 'test']:
    path = RAW / (HEAD + '-' + label + '.log')
    data = path.read_bytes()
    assert data.rstrip().endswith(b'EXIT=0'), path.name
    (HERE / path.name).write_bytes(data)
    logs.append({'path': str(path), 'sha256': digest(data), 'exit_code': 0})

result = {'base': BASE, 'head': HEAD, 'product_source': PRODUCT, 'snapshot': str(SNAP), 'files': files, 'captures_verified': len(captures),
          'checks': checks, 'final_lint_sha256': digest(lint_bytes), 'final_lint_summary': lint['summary'], 'final_logs': logs,
          'limitations': ['No native keyboard or browser zoom 200% acceptance.', 'No independent visual acceptance.', 'Historical generated extracts are not independently re-audited in full.']}
(HERE / 'verification.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
print(json.dumps({'files': len(files), 'captures_verified': len(captures), 'final_lint': lint['summary'], 'fixed_snapshot': True}))
