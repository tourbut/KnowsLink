"""Summarize raw-walk.out (walk.py output via explore.py); split pages into raw/ (gitignored)."""
import json, pathlib
t = open('raw-walk.out', encoding='utf-8').read()
j = json.loads(t[:t.rindex('\nEXIT=')].strip().splitlines()[-1])
raw = pathlib.Path('raw'); raw.mkdir(exist_ok=True)
for k, v in j['pages'].items(): (raw / f'{k}.html').write_text(v, encoding='utf-8')
for c in j['checks']: print('PASS' if c['pass'] else 'FAIL', c['id'], '|', c['detail'][:200] if not c['pass'] else '')
print(sum(c['pass'] for c in j['checks']), '/', len(j['checks']))
