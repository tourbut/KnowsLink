"""Render fixture pages in fixed-width iframes (1440/390/320) of headless Chrome; record page and table overflow.

file:// pages share an origin only with --allow-file-access-from-files, so the wrapper can read each frame's layout.
"""
import json, pathlib, re, subprocess, sys, tempfile
CHROME = r'C:\Program Files\Google\Chrome\Application\chrome.exe'
src = pathlib.Path(sys.argv[1]); out = pathlib.Path(sys.argv[2])
pages = ['policy-privacy-ko', 'policy-privacy-en', 'policy-terms-ko', 'policy-terms-en', 'policy-start-ko', 'policy-start-en', 'human-account', 'en-human-account', 'gate-active', 'en-gate-active']
frames = ''.join(f'<iframe data-p="{p}" data-w="{w}" style="width:{w}px;height:900px;border:0;display:block" src="{(src / (p + ".html")).as_uri()}"></iframe>' for p in pages for w in (1440, 390, 320))
script = '''<script>addEventListener("load",()=>{const r=[...document.querySelectorAll("iframe")].map(f=>{const d=f.contentDocument;return{page:f.dataset.p,width:+f.dataset.w,iw:f.contentWindow.innerWidth,sw:d.documentElement.scrollWidth,tables:[...d.querySelectorAll(".table-scroll")].map(e=>[e.scrollWidth,e.clientWidth])}});document.body.setAttribute("data-m",JSON.stringify(r))})</script>'''
with tempfile.TemporaryDirectory() as tmp:
    wrap = pathlib.Path(tmp) / 'wrap.html'
    wrap.write_text('<!doctype html><body>' + frames + script + '</body>', encoding='utf-8')
    dom = subprocess.run([CHROME, '--headless=new', '--disable-gpu', '--allow-file-access-from-files', f'--user-data-dir={tmp}/p', '--window-size=1500,900', '--virtual-time-budget=5000', '--dump-dom', wrap.as_uri()], capture_output=True, text=True, encoding='utf-8', timeout=120).stdout
m = re.search(r'data-m="([^"]*)"', dom)
import html as h
results = json.loads(h.unescape(m.group(1)))
for r in results:
    r['overflow'] = r['sw'] > r['iw']
(out / 'widths.json').write_text(json.dumps(results, indent=1) + '\n', encoding='utf-8')
bad = [r for r in results if r['overflow'] or r['iw'] != r['width']]
for r in results: print(r['page'], r['width'], 'iw', r['iw'], 'sw', r['sw'], 'tables', r['tables'])
print('BAD', bad)
sys.exit(1 if bad else 0)
