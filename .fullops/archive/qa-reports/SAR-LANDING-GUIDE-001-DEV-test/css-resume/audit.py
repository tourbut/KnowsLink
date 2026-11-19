"""Preserve bounded source excerpts and hashes for the resumed CSS task's partial packet."""
import hashlib
import json
import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[6]
PACKET = ROOT / '.fullops-squad/docs/evaluations/jev/SAR-LANDING-GUIDE-001-DEV-packet.json'
packet = json.loads(PACKET.read_text(encoding='utf-8'))
paths = {item['path'] for item in packet['items']}
paths.update(item['path'] for item in packet['unknown'] if 'path' in item)
paths.update(packet['remaining_context_paths'])
records = []
for name in sorted(paths):
    raw = (ROOT / name).read_bytes()
    lines = raw.decode('utf-8-sig').splitlines()
    matches = [(i + 1, line[:420]) for i, line in enumerate(lines)
               if re.search('SAR-LANDING|conversation-list|간격|대비|공개 첫|200%|배포|송신', line)]
    record = {'path': name, 'sha256': hashlib.sha256(raw).hexdigest(), 'lines': len(lines),
              'excerpts': matches[:4], 'scope': 'bounded relevant passages, not a full historical audit'}
    records.append(record)
    print(name)
    for number, text in record['excerpts'][:2]:
        print(f'  {number}: {text[:220]}')
pathlib.Path(__file__).with_name('source-audit.json').write_text(json.dumps(records, ensure_ascii=False, indent=2), encoding='utf-8')
