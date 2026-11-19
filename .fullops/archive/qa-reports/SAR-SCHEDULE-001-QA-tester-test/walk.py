"""Independent HTTP walk of the synthetic calendar fixture. No real Google account.

Reads the local forward origin from the task manifest. Prints verdicts and writes
walk-result.json. Cookie values and OAuth state stay out of that file.
"""
import html
import json
import re
import subprocess
import time
import urllib.parse
from datetime import datetime, timezone
from pathlib import Path

HERE = Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/schedule001-qa')
OUT = Path(__file__).resolve().parent
SECRET_TEXT = ('secret-title', 'secret-attendee', 'secret-location', '비밀 회의', 'secret-attendee@example.test')
KO = {
    'help': '친구가 요청하고 내가 승인할 때만, 해당 기간의 바쁜 시간대를 알려 줍니다. 일정 제목·참석자·장소는 공유하지 않습니다.',
    'connected': '연결됨 — 바쁜 시간만 확인합니다',
    'scope': '승인하면 이 요청의 기간에 대해 바쁜 시간대만 한 번 공유합니다. 일정을 만들거나 바꾸지 않습니다.',
    'disclosed': '바쁜 시간대만',
    'v1': '이 요청은 승인할 수 없습니다. 거절하거나 기한이 지나면 조회 없이 끝납니다.',
    'unconnected': '캘린더가 연결되지 않아 이 요청에 답할 수 없습니다.',
    'd2': '해제하면 KnowsLink에 보관한 토큰을 즉시 폐기하고 Google에도 권한 철회를 요청합니다. 아직 상대에게 전달되지 않은 조회 요청과 결과는 거절로 끝나고 공유되지 않습니다. 이미 전달된 결과는 되돌릴 수 없습니다.',
    'off': '캘린더 연결을 해제했습니다. KnowsLink에 보관한 토큰을 폐기했고 아직 전달되지 않은 조회 요청과 결과는 거절로 끝났습니다.',
    'partial': '캘린더 권한이 허용되지 않아 연결하지 않았습니다',
    'done_connect': '캘린더 연결 완료',
    'reconnect': '다시 연결 필요',
    'consent_note': 'Google 동의 화면에서는 빈 시간 조회 권한만 요청합니다. 로그인 동의 범위는 바뀌지 않습니다.',
    'no_decision': '지금은 승인하거나 거절할 수 없습니다',
    'old_scope': '정보 공개나 일정 실행을 허용하지 않습니다',
    'old_policy': '가능 시간 공개 정책이 없어',
}
EN = {
    'help': 'Only when a friend asks and you approve, your busy time slots for that period are shared. Event titles, attendees and locations are never shared.',
    'connected': 'Connected — only busy times are checked',
    'scope': "Approving shares only your busy time slots for this request's period, once. No event is created or changed.",
    'disclosed': 'Busy time slots only',
    'v1': 'This request cannot be approved. If you decline it or it expires, it ends without a lookup.',
    'unconnected': 'This request cannot be answered because no calendar is connected.',
    'd2': 'Disconnecting discards the token KnowsLink stores at once and asks Google to revoke access. Availability requests and results not yet delivered end as declined and are not shared. Results already delivered cannot be taken back.',
    'off': 'Calendar disconnected. The token KnowsLink stored was discarded, and availability requests and results not yet delivered ended as declined.',
    'reconnect': 'Reconnection needed',
    'consent_note': 'The Google consent screen asks only for free/busy access. The sign-in consent scope does not change.',
}


class Client:
    def __init__(self, origin):
        self.origin = origin
        parsed = urllib.parse.urlparse(origin)
        self.host, self.port = parsed.hostname, parsed.port
        self.cookies = {}

    def req(self, method, path, form=None, lang=None):
        last = None
        for _ in range(8):
            last = self._req(method, path, form, lang)
            if last[0] != 429 or path.startswith('/fixture'):
                return last
            time.sleep(8)
        return last

    def _req(self, method, path, form, lang):
        import http.client
        conn = http.client.HTTPConnection(self.host, self.port, timeout=30)
        headers = {'Origin': self.origin, 'Cache-Control': 'no-store'}
        cookies = dict(self.cookies)
        if lang:
            cookies['__Host-kl_lang'] = lang
        if cookies:
            headers['Cookie'] = '; '.join(f'{k}={v}' for k, v in cookies.items())
        body = None
        if form is not None:
            body = urllib.parse.urlencode(form).encode()
            headers['Content-Type'] = 'application/x-www-form-urlencoded'
        conn.request(method, path, body=body, headers=headers)
        response = conn.getresponse()
        data = response.read().decode('utf-8', 'replace')
        for header in response.getheaders():
            if header[0].lower() != 'set-cookie':
                continue
            name, _, rest = header[1].partition('=')
            value = rest.split(';', 1)[0]
            if 'Max-Age=0' in header[1] or value == '':
                self.cookies.pop(name, None)
            else:
                self.cookies[name] = value
        location = ''
        for key, value in response.getheaders():
            if key.lower() == 'location':
                location = value
        conn.close()
        return response.status, location, data


def browser(origin, stage):
    script = Path(__file__).resolve().parent / 'browser.mjs'
    proc = subprocess.run(['node', str(script), origin, stage], cwd=script.parent, capture_output=True, text=True, timeout=120)
    tail = (proc.stdout + proc.stderr)[-500:].replace('\n', ' ')
    return proc.returncode == 0, tail


def now():
    return datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')


def csrf(page):
    match = re.search(r'name="csrf" value="([^"]+)"', page)
    return match.group(1) if match else ''


def shown(page):
    """Templates escape apostrophes. Compare the text a browser would show."""
    return html.unescape(page)


def visible(page):
    text = re.sub(r'<script[\s\S]*?</script>|<style[\s\S]*?</style>', ' ', page)
    text = re.sub(r'<[^>]+>', ' ', text)
    return re.sub(r'\s+', ' ', text)


def query_of(location):
    parsed = urllib.parse.urlparse(location)
    return urllib.parse.parse_qs(parsed.query)


def same_path(origin, location):
    """Follow a fixture redirect on the local forward, even if the server named its own port."""
    if not location:
        return '/'
    parsed = urllib.parse.urlparse(location)
    if parsed.scheme in ('http', 'https') and parsed.hostname not in (None, '127.0.0.1'):
        raise RuntimeError('refusing non-loopback redirect')
    path = parsed.path or '/'
    if parsed.query:
        path += '?' + parsed.query
    return path


def disclosure_problems(view, gran):
    raw = json.dumps(view.get('disclosure'), ensure_ascii=False)
    problems = []
    for word in SECRET_TEXT:
        if word in raw or word in json.dumps(view, ensure_ascii=False):
            problems.append('secret:' + word)
    data = view.get('disclosure')
    if not isinstance(data, dict):
        return problems + ['missing']
    if set(data) != {'window', 'granularity_min', 'busy'}:
        problems.append('keys')
    if data.get('granularity_min') != gran:
        problems.append('gran')
    window = data.get('window') or {}
    if set(window) != {'start', 'end'}:
        problems.append('window')
    for item in data.get('busy') or []:
        if set(item) != {'start', 'end'}:
            problems.append('busy-keys')
    return problems


def main():
    public = json.loads((HERE / 'browser-fixture.json').read_text(encoding='utf-8'))
    origin = public['local_origin']
    gates = public['gates']
    results = []
    mutations = []

    def record(name, ok, detail):
        results.append({'name': name, 'verdict': 'PASS' if ok else 'FAIL', 'detail': detail})
        print(('PASS' if ok else 'FAIL'), name, detail if not ok else '', flush=True)

    def sign(name):
        client = Client(origin)
        status, _, _ = client.req('POST', f'/fixture/as/{name}')
        record(f'sign-in-{name}', status == 204, str(status))
        return client

    root = Client(origin)
    status, _, _ = root.req('GET', '/')
    record('anon-root', status == 200, str(status))
    status, _, body = root.req('GET', '/v1/pull')
    record('agent-api-blocked', status == 403 and 'outside calendar fixture scope' in body, str(status))

    b = sign('b')
    status, _, page = b.req('GET', '/home/account')
    record('sc07-b-connected-ko', status == 200 and all(shown(page).find(KO[k]) >= 0 for k in ('help', 'connected', 'd2')), str(status))
    record('sc07-b-no-connect-button', 'action="/home/calendar/connect"' not in page and 'action="/home/calendar/disconnect"' in page, '')
    status, _, en = b.req('GET', '/home/account', lang='en')
    record('sc07-b-connected-en', status == 200 and all(shown(en).find(EN[k]) >= 0 for k in ('help', 'connected', 'd2')), str(status))

    status, location, _ = b.req('POST', '/auth/reauth')
    scope = ' '.join(query_of(location).get('scope', []))
    record('sc01-login-scope-unchanged', status == 303 and scope == 'openid email' and 'access_type' not in query_of(location), scope)

    gate = gates['connected_pending_2h_30']
    status, _, page = b.req('GET', f'/home/gates/{gate}')
    record('sc07-gate-ko', status == 200 and KO['scope'] in shown(page) and KO['disclosed'] in shown(page) and 'value="approve"' in page, str(status))
    record('sc07-gate-no-old-copy', KO['old_scope'] not in page and KO['old_policy'] not in page, '')
    status, _, en = b.req('GET', f'/home/gates/{gate}', lang='en')
    record('sc07-gate-en', status == 200 and EN['scope'] in shown(en) and EN['disclosed'] in shown(en) and 'value="approve"' in en, str(status))

    c = sign('c')
    ugate = gates['unconnected_pending']
    status, _, page = c.req('GET', f'/home/gates/{ugate}')
    record('sc07-v1-ko', status == 200 and KO['v1'] in shown(page) and KO['unconnected'] in shown(page) and 'value="approve"' not in page and 'value="deny"' in page, str(status))
    record('sc07-v1-no-old', KO['old_scope'] not in page and KO['scope'] not in page, '')
    status, _, en = c.req('GET', f'/home/gates/{ugate}', lang='en')
    record('sc07-v1-en', status == 200 and EN['v1'] in shown(en) and EN['unconnected'] in shown(en) and 'value="approve"' not in en, str(status))
    token = csrf(page)
    status, _, body = c.req('POST', f'/home/gates/{ugate}', {'csrf': token, 'decision': 'approve'})
    view = json.loads(c.req('GET', f'/fixture/gate/{ugate}')[2])
    record('sc05-unconnected-approve-refused', status == 409 and view.get('state') == 'pending' and 'lookup' not in view, f'{status} {view.get("state")}')

    a = sign('a')
    status, _, page = a.req('GET', '/home/account')
    record('sc01-a-unconnected', status == 200 and KO['help'] in page and KO['consent_note'] in page and 'action="/home/calendar/connect"' in page and KO['connected'] not in page, str(status))
    status, location, _ = a.req('POST', '/home/calendar/connect')
    consent = query_of(location)
    record(
        'sc01-consent-scope',
        status == 303 and ' '.join(consent.get('scope', [])) == 'openid email https://www.googleapis.com/auth/calendar.freebusy'
        and consent.get('access_type') == ['offline'] and consent.get('include_granted_scopes') == ['true'] and consent.get('code_challenge_method') == ['S256'],
        'scope-keys-only',
    )
    status, _, auth = a.req('GET', same_path(origin, location))
    record('sc01-consent-page', status == 200 and '합성 Google 동의 화면' in auth and 'calendar.freebusy' in auth and '일부만 허용' in auth and '거절' in auth, str(status))
    partial = re.search(r'href="([^"]+)">일부만 허용', auth)
    status, _, body = a.req('GET', same_path(origin, partial.group(1).replace('&amp;', '&')))
    record('sc01-partial-not-stored', status == 403 and KO['partial'] in body, str(status))
    status, _, page = a.req('GET', '/home/account')
    record('sc01-partial-still-unconnected', KO['connected'] not in page and 'action="/home/calendar/connect"' in page, '')
    status, location, _ = a.req('POST', '/home/calendar/connect')
    status, _, auth = a.req('GET', same_path(origin, location))
    allow = re.search(r'href="([^"]+)">허용 \(빈 시간 조회 포함\)', auth)
    mutations.append({'at': now(), 'who': 'a', 'action': 'allow-connect'})
    status, _, body = a.req('GET', same_path(origin, allow.group(1).replace('&amp;', '&')))
    record('sc01-a-connect', status == 200 and KO['done_connect'] in body, str(status))
    status, _, page = a.req('GET', '/home/account')
    record('sc01-a-connected-copy', KO['connected'] in page and KO['d2'] in page, str(status))
    mutations.append({'at': now(), 'who': 'a', 'action': 'disconnect'})
    status, location, _ = a.req('POST', '/home/calendar/disconnect')
    record('sc01-a-disconnect-redirect', status == 303 and 'calendar-off' in location, str(status))
    status, _, page = a.req('GET', '/home/account?n=calendar-off')
    record('sc01-a-disconnect-copy', status == 200 and KO['off'] in page and 'action="/home/calendar/connect"' in page and 'action="/home/calendar/disconnect"' not in page, str(status))
    status, _, en = a.req('GET', '/home/account?n=calendar-off', lang='en')
    record('sc07-a-off-en', EN['off'] in en, str(status))

    b = sign('b')
    success = gates['connected_pending_4h_15']
    status, _, page = b.req('GET', f'/home/gates/{success}')
    before = json.loads(b.req('GET', f'/fixture/gate/{success}')[2])
    token = csrf(page)
    mutations.append({'at': now(), 'who': 'b', 'action': 'approve-4h-15'})
    status, _, _ = b.req('POST', f'/home/gates/{success}', {'csrf': token, 'decision': 'approve'})
    record('sc02-approve-redirect', status == 303, str(status))

    def settled(view):
        return view.get('lookup') in ('done', 'failed')

    view = None
    deadline = time.time() + 25
    while time.time() < deadline:
        view = json.loads(b.req('GET', f'/fixture/gate/{success}')[2])
        if settled(view):
            break
        time.sleep(0.4)
    problems = disclosure_problems(view, 15)
    record('sc02-one-lookup-done', before.get('lookup') is None and view.get('lookup') == 'done' and view.get('state') == 'approved', view.get('lookup'))
    record('sc03-disclosure-minimised', not problems, ','.join(problems))
    status, _, page = b.req('GET', f'/home/gates/{success}')
    record('sc02-revisit-no-form', 'value="approve"' not in page and KO['no_decision'] in page, '')
    record('sc03-page-has-no-secrets', all(word not in page for word in SECRET_TEXT), '')
    status, _, body = b.req('POST', f'/home/gates/{success}', {'csrf': token, 'decision': 'approve'})
    again = json.loads(b.req('GET', f'/fixture/gate/{success}')[2])
    record('sc02-duplicate-409', status == 409 and '이번 제출은 적용되지 않았습니다' in body and again.get('lookup') == 'done' and again.get('disclosure') == view.get('disclosure'), str(status))

    mutations.append({'at': now(), 'who': 'b', 'action': 'toggle-outage-on'})
    status, _, body = b.req('POST', '/fixture/google/toggle/outage')
    compact = body.replace(' ', '')
    record('sc05-outage-on', status == 200 and '"outage":true' in compact, compact[:80])
    outage = gates['connected_pending_8h_60']
    status, _, page = b.req('GET', f'/home/gates/{outage}')
    mutations.append({'at': now(), 'who': 'b', 'action': 'approve-8h-during-outage'})
    b.req('POST', f'/home/gates/{outage}', {'csrf': csrf(page), 'decision': 'approve'})
    view = {'lookup': ''}
    deadline = time.time() + 25
    while time.time() < deadline:
        view = json.loads(b.req('GET', f'/fixture/gate/{outage}')[2])
        if view.get('lookup') in ('done', 'failed'):
            break
        time.sleep(0.4)
    status, _, account = b.req('GET', '/home/account')
    record('sc05-outage-failed-not-done', view.get('lookup') == 'failed' and 'disclosure' not in view and KO['connected'] in account and KO['reconnect'] not in account, view.get('lookup'))
    b.req('POST', '/fixture/google/toggle/outage')
    mutations.append({'at': now(), 'who': 'b', 'action': 'toggle-outage-off'})

    mutations.append({'at': now(), 'who': 'b', 'action': 'toggle-revoke-grant'})
    b.req('POST', '/fixture/google/toggle/revoke-grant')
    revoked = gates['connected_pending_2h_30']
    status, _, page = b.req('GET', f'/home/gates/{revoked}')
    mutations.append({'at': now(), 'who': 'b', 'action': 'approve-2h-after-revoke'})
    b.req('POST', f'/home/gates/{revoked}', {'csrf': csrf(page), 'decision': 'approve'})
    view = {'lookup': ''}
    deadline = time.time() + 25
    while time.time() < deadline:
        view = json.loads(b.req('GET', f'/fixture/gate/{revoked}')[2])
        if view.get('lookup') in ('done', 'failed'):
            break
        time.sleep(0.4)
    status, _, account = b.req('GET', '/home/account')
    record('sc05-revoke-failed', view.get('lookup') == 'failed' and 'disclosure' not in view, view.get('lookup'))
    record('sc07-v3-d2-reconnect', KO['reconnect'] in account and KO['d2'] in account and 'action="/home/calendar/connect"' in account and 'action="/home/calendar/disconnect"' in account, '')
    status, _, en = b.req('GET', '/home/account', lang='en')
    record('sc07-v3-en', EN['reconnect'] in en and EN['d2'] in en, '')
    ok, detail = browser(origin, 'reconnect')
    record('browser-reconnect', ok, detail)

    held = json.loads(b.req('GET', f'/fixture/gate/{success}')[2])
    mutations.append({'at': now(), 'who': 'b', 'action': 'disconnect-with-undelivered'})
    status, location, _ = b.req('POST', '/home/calendar/disconnect')
    record('sc06-disconnect-redirect', status == 303 and 'calendar-off' in location, str(status))
    status, _, page = b.req('GET', '/home/account?n=calendar-off')
    dropped = json.loads(b.req('GET', f'/fixture/gate/{success}')[2])
    pending_left = json.loads(b.req('GET', f'/fixture/gate/{revoked}')[2])
    record('sc06-undelivered-hidden', held.get('lookup') == 'done' and 'disclosure' in held and dropped.get('state') == 'denied' and 'disclosure' not in dropped and dropped.get('agent_sees') == 'denied', dropped.get('state'))
    record('sc06-other-pending-denied', pending_left.get('agent_sees') == 'denied' and 'disclosure' not in pending_left, pending_left.get('state'))
    record('sc06-off-copy-and-no-token-ui', KO['off'] in page and 'action="/home/calendar/disconnect"' not in page and 'action="/home/calendar/connect"' in page and KO['connected'] not in page, '')
    status, _, en = b.req('GET', '/home/account?n=calendar-off', lang='en')
    record('sc07-off-en', EN['off'] in en, '')
    ok, detail = browser(origin, 'off')
    record('browser-off', ok, detail)

    expired = gates['connected_expired']
    status, _, page = b.req('GET', f'/home/gates/{expired}')
    record('sc02-expired-no-approve', status == 200 and 'value="approve"' not in page, str(status))
    denied = gates['connected_processed_denied']
    status, _, page = b.req('GET', f'/home/gates/{denied}')
    record('sc02-denied-no-approve', status == 200 and 'value="approve"' not in page, str(status))

    payload = {'origin': origin, 'sha': public['sha'], 'finished_at': now(), 'mutations': mutations, 'results': results, 'fail': sum(1 for item in results if item['verdict'] != 'PASS')}
    (OUT / 'walk-result.json').write_text(json.dumps(payload, ensure_ascii=False, indent=1), encoding='utf-8')
    print(f'FAILS={payload["fail"]} TOTAL={len(results)}', flush=True)
    return 1 if payload['fail'] else 0


if __name__ == '__main__':
    raise SystemExit(main())
