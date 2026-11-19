"""Independent synthetic HTTP assertions; requires a fresh fixture with an unnamed landing account.

Run before browser.mjs; preserves the named synthetic member for visual QA and never records cookies or codes.
"""
import http.client
import json
import re
from http.cookies import SimpleCookie
from pathlib import Path
from urllib.parse import urlencode, urlsplit

OUT = Path(__file__).parent
records = []
requests = []


class Client:
    def __init__(self, port):
        self.port, self.cookies = port, {}

    def request(self, method, path, values=None, origin=None):
        assert method in ('GET', 'POST')
        assert path.split('?')[0] in ('/', '/home', '/home/name', '/home/account', '/auth/start', '/auth/verify', '/auth/google', '/auth/google/callback', '/auth/logout', '/home/profile', '/fixture/code', '/fixture/google-cancel', '/fixture/expire')
        headers = {'Cookie': '; '.join(k + '=' + v for k, v in self.cookies.items())}
        if method == 'POST':
            headers.update({'Content-Type': 'application/x-www-form-urlencoded', 'Origin': origin or f'http://127.0.0.1:{self.port}'})
        connection = http.client.HTTPConnection('127.0.0.1', self.port, timeout=20)
        connection.request(method, path, urlencode(values or {}) if method == 'POST' else None, headers)
        response = connection.getresponse()
        text = response.read().decode()
        for key, value in response.getheaders():
            if key.lower() == 'set-cookie':
                cookie = SimpleCookie(value)
                for name, item in cookie.items():
                    self.cookies[name] = item.value
        result = response.status, dict((k.lower(), v) for k, v in response.getheaders()), text
        requests.append({'method': method, 'port': self.port, 'path': urlsplit(path).path, 'status': response.status})
        connection.close()
        return result


def check(name, condition):
    records.append({'name': name, 'passed': bool(condition)})
    (OUT / 'http-results.json').write_text(json.dumps({'checks': records, 'requests': requests}, indent=2), encoding='utf-8')
    if not condition:
        raise AssertionError(name)


for port, google, email in [(18781, True, False), (18782, True, True), (18783, False, True), (18784, False, False)]:
    c = Client(port)
    status, headers, body = c.request('GET', '/')
    check(f'provider-{port}', status == 200 and body.count('action="/auth/google"') == int(google) and body.count('action="/auth/start"') == int(email))
    check(f'public-boundary-{port}', headers.get('cache-control') == 'no-store' and "form-action 'self' https://accounts.google.com" in headers.get('content-security-policy', '') and 'landing@example.test' not in body and '.fullops-squad' not in body)
    check(f'private-denied-{port}', c.request('GET', '/home')[0] == 303)
    if not google and not email:
        check('none-safe-reload', '안내 다시 불러오기' in body and '로그인 상태 확인' in body and '<form' not in body and not re.search(r'<(?:button|input)[^>]*\bdisabled\b', body))

c = Client(18783)
status, headers, body = c.request('POST', '/auth/start', {'email': '<invalid>'})
check('invalid-email-422-associated-escaped', status == 422 and 'aria-describedby="email-help problem"' in body and '<invalid>' not in body and '로그인 문제를 확인하세요' in body)
status, _, _ = c.request('POST', '/auth/start', {'email': 'landing@example.test'}, 'https://foreign.test')
check('cross-origin-login-403', status == 403)
status, headers, _ = c.request('POST', '/auth/start', {'email': 'landing@example.test'})
check('email-start-303', status == 303 and headers.get('location') == '/auth/verify')
status, _, body = c.request('GET', '/auth/verify')
check('verify-pending', status == 200 and '아직 로그인되지 않았습니다' in body)
_, _, code = c.request('GET', '/fixture/code')
assert code.isdigit() and len(code) == 6
wrong = '000000' if code != '000000' else '111111'
status, _, body = c.request('POST', '/auth/verify', {'code': wrong})
check('wrong-code-401-associated', status == 401 and '남은 시도: 4회' in body and 'aria-describedby="code-help problem"' in body)
status, headers, _ = c.request('POST', '/auth/verify', {'code': code})
check('valid-code-303-home', status == 303 and headers.get('location') == '/home')
status, _, body = c.request('GET', '/')
check('authed-public-200-no-private-data', status == 200 and '로그인되어 있습니다.' in body and '<form' not in body and 'landing@example.test' not in body and '비공개 QA 이름' not in body)
status, headers, body = c.request('GET', '/home')
check('unnamed-home-redirect', status == 303 and headers.get('location') == '/home/name')
status, _, body = c.request('GET', '/home/name')
check('name-form', status == 200 and 'action="/home/profile"' in body)
status, _, _ = c.request('POST', '/home/profile', {'name': '비공개 QA 이름'}, 'https://foreign.test')
check('profile-csrf-403', status == 403)
status, _, _ = c.request('POST', '/home/profile', {'name': '비공개 QA 이름'})
check('name-save-303', status == 303)
status, _, body = c.request('GET', '/home')
check('home-and-primary-navigation', status == 200 and all(f'href="{path}"' in body for path in ['/home', '/home/friends', '/home/agents']) and 'class="landing"' not in body and '.landing .hero' not in body)
status, _, body = c.request('GET', '/home/account')
check('account-service-guide-link', status == 200 and 'href="/">서비스 안내</a>' in body)
status, _, body = c.request('GET', '/')
check('named-public-no-name-email', status == 200 and '비공개 QA 이름' not in body and 'landing@example.test' not in body)
status, _, _ = c.request('POST', '/fixture/expire')
check('fixture-expire-303', status == 303)
status, _, body = c.request('GET', '/')
check('expired-session-restores-form', status == 200 and '세션이 끝났습니다. 다시 로그인하세요.' in body and 'action="/auth/start"' in body and not any(v for k, v in c.cookies.items() if 'session' in k))

g = Client(18781)
status, headers, _ = g.request('POST', '/auth/google')
check('google-start-local-mock', status == 303 and headers.get('location', '').startswith('/fixture/google-cancel?'))
status, headers, _ = g.request('GET', headers['location'])
check('google-local-cancel-redirect', status == 303 and headers.get('location', '').startswith('/auth/google/callback?'))
status, _, body = g.request('GET', headers['location'])
check('google-cancel-safe-recovery', status == 401 and '로그인 문제를 확인하세요' in body and '로그인 다시 시도' in body and 'action="/auth/start"' not in body)
status, _, body = g.request('GET', '/auth/verify')
check('missing-challenge-410-provider-safe', status == 410 and 'action="/auth/google"' in body and 'action="/auth/start"' not in body)
status, _, body = g.request('POST', '/auth/start', {'email': 'mail-failure@example.test'})
check('mail-unavailable-503-no-false-success', status == 503 and '로그인되지 않았습니다' in body and '1분 뒤' in body and 'action="/auth/start"' not in body)
print(f'HTTP checks passed: {len(records)}; requests: {len(requests)}')
