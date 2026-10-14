"""Deterministic Orca checks and captures of the isolated landing fixture; never use operational accounts."""
import base64
import hashlib
import json
import pathlib
import subprocess
import sys
import struct
import urllib.request

ORCA = 'C:/Users/shin/AppData/Local/Programs/orca/resources/bin/orca.EXE'
PAGE = '7c54ca56-3195-4698-8d36-1d120a8d6997'
HERE = pathlib.Path(__file__).parent
OUT = HERE / ('browser-' + sys.argv[1])
OUT.mkdir(exist_ok=True)
records = []


def run(*args):
    result = subprocess.run([ORCA, *args, '--page', PAGE, '--json'], capture_output=True, text=True, encoding='utf-8')
    value = json.loads(result.stdout)
    if result.returncode or not value.get('ok'):
        raise RuntimeError(value)
    return value['result']


def js(expression):
    result = run('eval', '--expression', expression)
    result = result.get('result', result) if isinstance(result, dict) else result
    if result in ('true', 'false', 'null'):
        return json.loads(result)
    return result


def check(name, value):
    records.append({'check': name, 'value': value})
    (OUT / 'checks.json').write_text(json.dumps(records, ensure_ascii=False, indent=2), encoding='utf-8')
    if value is not True:
        raise AssertionError((name, value))


def shot(name):
    data = run('screenshot')
    (OUT / (name + '.png')).write_bytes(base64.b64decode(data['data']))


def snapshot():
    result = run('snapshot')
    (OUT / 'last-snapshot.json').write_text(json.dumps(result, ensure_ascii=False, indent=2), encoding='utf-8')
    return result


def ref(name, role):
    matches = [key for key, value in snapshot()['refs'].items() if value['name'] == name and value['role'] == role]
    if len(matches) != 1:
        raise ValueError((name, role, matches))
    return matches[0]


def layout():
    run('exec', '--command', 'set viewport 390 844 1')
    zoom = float(js('devicePixelRatio'))
    for width, height in [(1440, 900), (390, 844), (320, 844)]:
        run('goto', '--url', 'http://127.0.0.1:18782/')
        snapshot()
        scale = 0.65 if width == 1440 else 1
        for _ in range(4):
            run('exec', '--command', f'set viewport {round(width * zoom)} {round(height * zoom)} {scale}')
            if js(f'innerWidth === {width}'):
                break
        check(f'viewport-{width}', js(f'innerWidth === {width}'))
        check(f'hero-cta-{width}', js("document.querySelector('.hero .primary').getBoundingClientRect().bottom <= innerHeight"))
        shot(f'hero-{width}')
        js("document.querySelectorAll('details').forEach(d=>d.open=true)")
        check(f'overflow-{width}', js('document.documentElement.scrollWidth <= innerWidth'))
        check(f'hit-area-{width}', js("[...document.querySelectorAll('a,button,input,summary')].filter(e=>e.getClientRects().length).every(e=>parseFloat(getComputedStyle(e).minHeight)>=44 && e.getBoundingClientRect().height>=43.98)"))
        check(f'headings-{width}', js("document.querySelectorAll('h1').length===1 && document.documentElement.lang==='ko'"))
        for section in ['install', 'examples', 'start']:
            js(f"document.getElementById('{section}').scrollIntoView()")
            shot(f'{section}-{width}')
    # CSS zoom tests text/layout enlargement; browser chrome zoom is recorded separately.
    js("document.body.style.zoom='2'")
    check('css-zoom-200-overflow', js('document.documentElement.scrollWidth <= innerWidth'))
    shot('css-zoom-200')
    js("document.body.style.zoom='';scrollTo(0,0)")


def keyboard():
    run('goto', '--url', 'http://127.0.0.1:18782/')
    snapshot()
    run('exec', '--command', 'set viewport 390 844')
    run('click', '--element', ref('시작 순서 보기', 'link'))
    snapshot()
    js("document.querySelector('summary').focus()")
    run('keypress', '--key', 'Enter')
    check('summary-enter-opens', js("document.querySelector('details').open"))
    run('keypress', '--key', 'Space')
    check('summary-space-closes', js("!document.querySelector('details').open"))
    js("document.querySelector('header a[href=\"#start\"]').focus()")
    run('keypress', '--key', 'Enter')
    snapshot()
    check('anchor-enter', js("location.hash==='#start' && document.activeElement.id==='start'"))
    run('keypress', '--key', 'Tab')
    check('tab-focus', js("document.activeElement.tagName==='BUTTON' && getComputedStyle(document.activeElement).outlineWidth==='3px'"))
    run('keypress', '--key', 'Shift+Tab')
    check('shift-tab', js("document.activeElement.id==='start' || document.activeElement.tagName==='A'"))
    shot('keyboard-focus')


def states():
    for port, google, email in [(18781, True, False), (18782, True, True), (18783, False, True), (18784, False, False)]:
        run('goto', '--url', f'http://127.0.0.1:{port}/')
        snapshot()
        check(f'providers-{port}', js(f"!!document.querySelector('form[action=\"/auth/google\"]')==={str(google).lower()} && !!document.querySelector('form[action=\"/auth/start\"]')==={str(email).lower()}"))
    run('goto', '--url', 'http://127.0.0.1:18782/')
    snapshot()
    js("document.querySelector('header a[href=\"#start\"]').click()")
    snapshot()
    js("document.querySelector('form[action=\"/auth/google\"] button').click()")
    run('wait', '--text', '로그인 문제를 확인하세요')
    snapshot()
    check('google-cancel-recovery', js("document.querySelector('[role=alert]')!==null && document.querySelector('header+p a').getAttribute('href')==='#start'"))
    shot('google-error-top')
    js("document.querySelector('header+p a').click()")
    snapshot()
    shot('google-error-form')
    run('goto', '--url', 'http://127.0.0.1:18782/?n=expired')
    snapshot()
    check('expired-notice', js("document.body.innerText.includes('세션이 끝났습니다. 다시 로그인하세요.')"))
    run('fill', '--element', ref('이메일 주소', 'textbox'), '--value', 'landing@example.test')
    js("document.querySelector('form[action=\"/auth/start\"] button').click()")
    run('wait', '--url', '/auth/verify')
    snapshot()
    check('email-post-verification', js("location.pathname==='/auth/verify'"))
    # Synthetic-only fixture endpoint, never a real mail inbox.
    # Production CSP forbids browser fetch; this reads only the synthetic local mail fixture.
    with urllib.request.urlopen('http://127.0.0.1:18782/fixture/code') as response:
        code = response.read().decode()
    run('fill', '--element', ref('6자리 확인 코드', 'textbox'), '--value', code)
    js("document.querySelector('form button').click()")
    run('wait', '--url', '/home/name')
    snapshot()
    check('unnamed-member', js("location.pathname==='/home/name'"))
    run('goto', '--url', 'http://127.0.0.1:18782/')
    snapshot()
    check('authenticated-public', js("location.pathname==='/' && document.body.innerText.includes('로그인되어 있습니다.') && document.querySelectorAll('form').length===0 && !document.body.innerText.includes('landing@example.test')"))
    shot('authed')


def audit():
    # Host-only cookies span ports; localhost keeps this read-only audit anonymous.
    run('goto', '--url', 'http://localhost:18781/')
    snapshot()
    tokens = json.loads(js("JSON.stringify(Object.fromEntries(['canvas','ink','muted','accent','edge','danger','selected'].map(k=>[k,getComputedStyle(document.documentElement).getPropertyValue('--'+k).trim()])))"))

    def luminance(color):
        channels = [int(color[n:n+2], 16) / 255 for n in (1, 3, 5)]
        linear = [c / 12.92 if c <= .04045 else ((c+.055)/1.055)**2.4 for c in channels]
        return sum(c*w for c, w in zip(linear, (.2126, .7152, .0722)))

    contrasts = []
    for foreground, background, minimum in [('ink','canvas',4.5), ('muted','canvas',4.5), ('accent','canvas',4.5), ('canvas','accent',4.5), ('danger','canvas',4.5), ('ink','selected',4.5), ('edge','canvas',3)]:
        values = sorted([luminance(tokens[foreground]), luminance(tokens[background])])
        ratio = (values[1]+.05)/(values[0]+.05)
        contrasts.append({'foreground': foreground, 'background': background, 'ratio': ratio, 'minimum': minimum})
        check(f'contrast-{foreground}-{background}', ratio >= minimum)
    (OUT / 'contrast.json').write_text(json.dumps({'tokens':tokens, 'pairs':contrasts}, indent=2), encoding='utf-8')
    js("document.querySelector('summary').click()")
    check('dom-summary-click-opens', js("document.querySelector('details').open"))
    js("document.querySelector('summary').click()")
    check('dom-summary-click-closes', js("!document.querySelector('details').open"))
    js("document.querySelector('header a[href=\"#start\"]').click()")
    check('dom-anchor-focus', js("location.hash==='#start' && document.activeElement.id==='start'"))
    check('focus-outline-rule', js("[...document.styleSheets].some(s=>[...s.cssRules].some(r=>r.selectorText===':focus-visible' && r.style.outline.includes('3px')))"))
    manifest = []
    for path in sorted(HERE.glob('browser-*/*.png')):
        data = path.read_bytes()
        width, height = struct.unpack('>II', data[16:24])
        manifest.append({'path':path.relative_to(HERE).as_posix(), 'png_width':width, 'png_height':height, 'sha256':hashlib.sha256(data).hexdigest()})
    (OUT / 'captures.json').write_text(json.dumps(manifest, indent=2), encoding='utf-8')


{'layout': layout, 'keyboard': keyboard, 'states': states, 'audit':audit}[sys.argv[1]]()
print(json.dumps({'phase': sys.argv[1], 'checks': len(records), 'passed': True}))
