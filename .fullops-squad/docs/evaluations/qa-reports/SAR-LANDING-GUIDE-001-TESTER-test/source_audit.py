"""Record bounded source inspection and provenance for the tester packet, including historical/partial candidates."""
import hashlib
import json
import pathlib
import re
import subprocess

ROOT = pathlib.Path(__file__).resolve().parents[5]
OUT = pathlib.Path(__file__).parent
KEY = 'SAR-LANDING-GUIDE-001-TESTER'
packet = json.loads((ROOT / f'.fullops-squad/docs/evaluations/jev/{KEY}-packet.json').read_text(encoding='utf-8'))
paths = set(packet['context_paths']) | set(packet.get('remaining_context_paths', [])) | set(packet['required_paths'])
paths.update(item['path'] for item in packet['items'])
rows = []
for name in sorted(paths):
    data = subprocess.check_output(['git', '-C', str(ROOT), 'show', packet['head'] + ':' + name])
    lines = data.decode('utf-8-sig', 'replace').splitlines()
    headings = [{'line': i + 1, 'text': line} for i, line in enumerate(lines) if line.startswith(('#', 'title:', 'summary:'))][:8]
    hits = [i for i, line in enumerate(lines) if re.search('LANDING|LG-0|startPage|startView|공개 첫페이지|native|200%|not_run|FAIL|미검증', line)]
    excerpts = [{'line': i + 1, 'text': '\n'.join(lines[max(0, i-1):i+2])} for i in hits[:4]]
    rows.append({'path': name, 'sha256': hashlib.sha256(data).hexdigest(), 'lines': len(lines), 'headings': headings, 'excerpts': excerpts})
(OUT / 'source-audit.json').write_text(json.dumps(rows, ensure_ascii=False, indent=2), encoding='utf-8')
product = ['cmd', 'internal', 'adapters', 'db', 'deploy', 'scripts', 'go.mod', 'go.sum', 'Makefile', 'Dockerfile', 'compose.yaml', 'sqlc.yaml', '.dockerignore']
identity = []
for first, second in [('4795d619513fde4fbcd4a55b94ec91b056ebbaf0', 'fe2387a7676f7997aea76a2606f4cc2dbcd4b56b'), ('fe2387a7676f7997aea76a2606f4cc2dbcd4b56b', packet['head'])]:
    result = subprocess.run(['git', '-C', str(ROOT), 'diff', '--exit-code', first, second, '--', *product], capture_output=True, text=True)
    identity.append({'from': first, 'to': second, 'paths': product, 'exit_code': result.returncode, 'diff': result.stdout})
before = subprocess.check_output(['git', '-C', str(ROOT), 'show', 'eae9a8bffae706410f684e4c837e396e7b0bcf26:internal/relay/member_pages.go'], text=True)
after = (ROOT / 'internal/relay/member_pages.go').read_text(encoding='utf-8')
css = lambda text: text.split('const memberStyle = `', 1)[1].split('`', 1)[0]
identity.append({'check': 'preexisting-memberStyle-identical', 'passed': css(before) == css(after)})
(OUT / 'identity.json').write_text(json.dumps(identity, indent=2), encoding='utf-8')
assert all(r.get('exit_code', 0) == 0 and r.get('passed', True) for r in identity)
print(f'Audited bounded headings/excerpts for {len(rows)} paths; runtime comparisons and unchanged memberStyle passed')
