"""Read-only evidence checks; no product tests, browser actions, or remote commands."""
import hashlib
import json
import pathlib
import subprocess

S = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-bb360033216c4b3785b0b7e8e2b33383')
OUT = pathlib.Path(__file__).parent
EXTERNAL = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-css-dev')
QA = S / '.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume'
JEV = S / '.fullops-squad/docs/evaluations/jev'
HEAD = 'fc4ad58b65471e2de298cb9f1aba5d2253ea8135'
BASE = 'fe2387a7676f7997aea76a2606f4cc2dbcd4b56b'
START = 'dc5ec35b2c0fa297ba91387b9794b0be8398872c'
SOURCE = 'f2b8cb622dc397034423582978e77560e0355771'

def git(*args):
    return subprocess.check_output(['git', '--no-pager', '-C', str(S), *args])

def digest(raw):
    return hashlib.sha256(raw).hexdigest()

def read(path):
    return json.loads(path.read_text(encoding='utf-8'))

data = {'base': BASE, 'head': HEAD, 'incremental_base': START, 'source': SOURCE, 'files': []}
for item in read(OUT / 'result.json')['files']:
    name = item['path']
    raw = (S / name).read_bytes()
    record = {'path': name, 'status': item['status'], 'sha256': digest(raw), 'bytes': len(raw),
              'changed_since_resume': bool(git('diff', '--name-only', START, HEAD, '--', name).strip())}
    if name.endswith('.json'):
        json.loads(raw)
        record['json_parse'] = 'passed'
    data['files'].append(record)
manifest = read(QA / 'fixture-manifest.json')
data['fixture_sources'] = [{'file': item['file'], 'matches': digest((S / item['file']).read_bytes()) == item['sha256']}
                           for item in manifest['sources']]
data['fixtures'] = [{'file': item['file'], 'matches': digest((EXTERNAL / 'fixtures' / item['file']).read_bytes()) == item['sha256']}
                    for item in manifest['fixtures']]
results = read(QA / 'results.json')
data['browser'] = {'source': results['source'], 'checks': len(results['checks']),
                   'failed': [item for item in results['checks'] if not item['passed']]}
data['captures'] = [{'file': item['file'], 'matches': digest((QA / item['file']).read_bytes()) == item['sha256']}
                    for item in results['captures']]
data['request_scope'] = all(item['url'].startswith('http://127.0.0.1:18896/') and item['method'] == 'GET'
                            for item in results['requests'])
product = ['internal', 'cmd', 'adapters', 'db', 'scripts', 'deploy', 'go.mod', 'go.sum', 'Makefile',
           'Dockerfile', 'compose.yaml', 'sqlc.yaml', '.github']
data['product_diff_source_to_final'] = git('diff', '--name-only', SOURCE, HEAD, '--', *product).decode().splitlines()
data['product_diff_base_to_final'] = git('diff', '--name-only', BASE, HEAD, '--', *product).decode().splitlines()
data['lint_original_sha256'] = digest((EXTERNAL / 'final-lint.json').read_bytes())
data['lint_copy_sha256'] = digest((OUT / 'lint.json').read_bytes())
data['final_logs'] = [{'name': f'{HEAD}-{target}.log',
                       'exit0': (EXTERNAL / f'{HEAD}-{target}.log').read_text().strip().endswith('EXIT=0'),
                       'sha256': digest((EXTERNAL / f'{HEAD}-{target}.log').read_bytes())}
                      for target in ['install', 'lint', 'test']]
name = '.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md'
data['author_body_unchanged'] = git('show', HEAD + ':' + name).decode().split('---\n', 2)[2].lstrip('\n') == git('show', START + ':' + name).decode()
packet = read(JEV / 'SAR-LANDING-GUIDE-001-DEV-packet.json')
outcomes = read(JEV / 'SAR-LANDING-GUIDE-001-DEV-packet-outcomes.json')
data['packet'] = {'items': len(packet['items']), 'outcomes': len(outcomes['items']),
                  'unknown': len(packet['unknown']), 'remaining': len(packet['remaining_context_paths']),
                  'partial': packet['partial'], 'attempt_matches': packet['attempt'] == outcomes['attempt'],
                  'input_matches': packet['input_sha256'] == outcomes['packet_input_sha256'],
                  'paths_match': {(v['path'], c) for v in packet['items'] for c in v['categories']} ==
                                 {(v['path'], v['category']) for v in outcomes['items']}}
data['prior_outcomes_unchanged'] = git('show', START + ':.fullops-squad/docs/evaluations/jev/SAR-LANDING-GUIDE-001-DEV-packet-outcomes.json') == git('show', HEAD + ':.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume/prior-packet-outcomes.json')
data['empty_inbox'] = (S / '.fullops-squad/handovers/to_dev.md').read_bytes().strip() == b''
data['snapshot_head'] = git('rev-parse', 'HEAD').decode().strip()
data['snapshot_clean'] = not git('status', '--porcelain').strip()
data['patch_check_exit'] = subprocess.run(['git', '--no-pager', '-C', str(S), 'diff', '--check', BASE, HEAD], capture_output=True).returncode
data['checks_all'] = (all(v['matches'] for group in ['fixtures', 'fixture_sources', 'captures'] for v in data[group])
                      and all(v['exit0'] for v in data['final_logs']) and data['author_body_unchanged']
                      and data['prior_outcomes_unchanged'] and data['empty_inbox'] and data['snapshot_clean']
                      and data['snapshot_head'] == HEAD and data['patch_check_exit'] == 0
                      and data['lint_original_sha256'] == data['lint_copy_sha256']
                      and not data['product_diff_source_to_final'] and not data['browser']['failed']
                      and data['request_scope'] and all(data['packet'][key] for key in ['attempt_matches', 'input_matches', 'paths_match']))
(OUT / 'verification.json').write_text(json.dumps(data, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
print(json.dumps({key: value for key, value in data.items() if key not in ['files', 'fixtures', 'captures']}, ensure_ascii=False, indent=2))
assert data['checks_all'], 'Evidence mismatch; inspect verification.json'
