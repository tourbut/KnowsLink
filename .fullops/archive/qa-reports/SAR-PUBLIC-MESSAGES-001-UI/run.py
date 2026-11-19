"""Exercise real member/receipt/gate pages on unchanged relay code and an isolated SQL fixture."""

import base64
import email
import email.policy
import hashlib
import json
import os
import re
import signal
import socket
import subprocess
import time
import urllib.error
import urllib.parse
import urllib.request
from datetime import datetime, timedelta, timezone
from pathlib import Path

from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey

from browser import Browser

EVIDENCE = Path(__file__).parent
IMAGES = EVIDENCE.parents[2] / "design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI"
ROOT = Path(Path('/tmp/sar-messages-ui-location').read_text())
CFG = json.loads((ROOT / 'private.json').read_text())
BASE = 'http://localhost:' + str(CFG['relay_port'])
AGENT = 'agent_a_' + 'a' * 120
PEER = 'agent_b_' + 'b' * 120
REQUEST = '0199a3f2-4c10-7a11-8b22-334455667788'
REPLY = '0199a3f2-4c10-7a11-8b22-334455667799'
GATE = '0199a3f2-4c10-7a11-8b22-3344556677aa'
events, captures = [], []


def instant(value=None):
    return (value or datetime.now(timezone.utc)).isoformat(timespec='seconds').replace('+00:00', 'Z')


def encoded(value):
    return base64.urlsafe_b64encode(value).decode().rstrip('=')


def canonical(value):
    return json.dumps(value, ensure_ascii=False, separators=(',', ':'), sort_keys=True).encode()


def sql(statement):
    result = subprocess.run(['docker', 'exec', '-i', CFG['container'], 'psql', '-U', 'knowslink', '-d', 'knowslink', '-At', '-v', 'ON_ERROR_STOP=1'],
                            input=statement, text=True, capture_output=True, check=True)
    return result.stdout


def state():
    return json.loads(sql('SELECT data FROM relay_state;'))


def mutate(operation, label):
    # Stop only this fixture relay before changing this fixture database.
    try:
        os.kill(CFG['relay_pid'], signal.SIGTERM)
    except ProcessLookupError:
        events.append({'fixtureProcess': 'previous relay already exited'})
    for _ in range(100):
        try:
            urllib.request.urlopen(BASE + '/healthz', timeout=.2)
        except (urllib.error.URLError, ConnectionError):
            break
        time.sleep(.03)
    current = state()
    operation(current)
    raw = json.dumps(current, ensure_ascii=False).replace("'", "''")
    sql("UPDATE relay_state SET data='" + raw + "'::jsonb,epoch=epoch+1,clock=clock_timestamp();")
    env = dict(os.environ, DATABASE_URL=f"postgres://knowslink:{CFG['password']}@127.0.0.1:{CFG['db_port']}/knowslink?sslmode=disable",
               RELAY_ADDR=f"127.0.0.1:{CFG['relay_port']}", KNOWSLINK_SMTP_URL=f"smtp://127.0.0.1:{CFG['smtp_port']}",
               KNOWSLINK_MAIL_FROM='noreply@example.invalid', KNOWSLINK_SYNTHETIC_SIGNUP='', KNOWSLINK_TEST_AGENTS='')
    CFG['relay_pid'] = subprocess.Popen([str(ROOT / 'relay')], cwd=ROOT / 'snapshot', env=env,
                         stdout=(ROOT / 'relay.log').open('a'), stderr=subprocess.STDOUT).pid
    (ROOT / 'private.json').write_text(json.dumps(CFG))
    for _ in range(100):
        try:
            urllib.request.urlopen(BASE + '/healthz', timeout=.3)
            break
        except (urllib.error.URLError, ConnectionError):
            time.sleep(.03)
    events.append({'fixtureMutation': label, 'privateDatabaseOnly': True})


def api(path, body, token, expected=200):
    request = urllib.request.Request(BASE + path, data=json.dumps(body).encode(),
                headers={'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token})
    try:
        response = urllib.request.urlopen(request)
    except urllib.error.HTTPError as error:
        response = error
    value = json.load(response)
    events.append({'api': path, 'status': response.status, 'error': value.get('error') if isinstance(value, dict) else None})
    assert response.status == expected, (path, response.status, value.get('error'))
    return value


def seed():
    keys = {agent: Ed25519PrivateKey.generate() for agent in (AGENT, PEER)}
    tokens = {agent: encoded(os.urandom(32)) for agent in keys}
    def apply(current):
        owner = next(m['Owner'] for m in current['Members'].values() if m['Email'].startswith('ui-member'))
        for name in ('Agents', 'Pairs', 'Messages', 'Gates', 'Idempotency', 'Rates', 'HTTP', 'Connections'):
            current[name] = {}
        for agent, key in keys.items():
            public = key.public_key().public_bytes(serialization.Encoding.Raw, serialization.PublicFormat.Raw)
            current['Agents'][agent] = {'Owner': owner, 'Credential': '', 'Keys': {'key1': {
                'Public': base64.b64encode(public).decode(), 'Revoked': False, 'Changed': instant(),
                'Credential': encoded(hashlib.sha256(tokens[agent].encode()).digest())}}, 'Revoked': False, 'Changed': instant()}
        current['Pairs'][AGENT + '/' + PEER] = {'A': AGENT, 'B': PEER, 'Inviter': AGENT, 'Recipient': PEER,
                'State': 'active', 'Generation': 1, 'Exp': instant(datetime.now(timezone.utc) + timedelta(hours=24))}
    mutate(apply, 'new member-owned keys and accepted pair; preserve real browser sessions')
    return keys, tokens


def send_text(keys, tokens, reply=False, content='비민감 연결 확인', expected=200):
    sender, target = (PEER, AGENT) if reply else (AGENT, PEER)
    envelope = {'v': 'knowslink.text.v1', 'id': REPLY if reply else REQUEST, 'from': sender, 'to': target,
        'text': content, 'exp': instant(datetime.now(timezone.utc) + timedelta(seconds=170 if reply else 180)),
        'idempotency_key': 'designer-text-reply-key' if reply else 'designer-text-request-key'}
    if reply:
        envelope['reply_to'] = REQUEST
    unsigned = dict(envelope, reply_to=envelope.get('reply_to', ''))
    signature = keys[sender].sign(b'KNOWSLINK-TEXT\x00knowslink.text.v1\x00Ed25519\x00key1\x00' + canonical(unsigned))
    envelope['sig'] = {'alg': 'Ed25519', 'kid': 'key1', 'value': encoded(signature)}
    return api('/v1/text/send', envelope, tokens[sender], expected)


def receive(tokens, recipient, prefix='/v1/text', ack=True):
    lease = api(prefix + '/pull', {}, tokens[recipient])
    delivery = {'id': lease['envelope']['id'], 'token': lease['lease_token']}
    if ack:
        api(prefix + '/persist', delivery, tokens[recipient])
        api(prefix + '/ack', delivery, tokens[recipient])
    return delivery['id']


def gate(keys, tokens):
    exp = instant(datetime.now(timezone.utc) + timedelta(seconds=180))
    body = {'window': {'start': '2026-10-03T10:00:00Z', 'end': '2026-10-03T11:00:00Z'}, 'granularity_min': 30}
    parent = {'v': 'relay.v1', 'id': REQUEST, 'from': AGENT, 'to': PEER, 'intent': 'schedule.query',
       'body': body, 'deliver': 'agent', 'exp': exp, 'idempotency_key': 'designer-gate-parent-key',
       'render': {'hint': '이 hint는 실제 본문과 다릅니다. <script>window.uiInjected=true</script> $(touch /tmp/should-not-run)'}}
    digest = encoded(hashlib.sha256(canonical({'v': 'relay.v1', 'from': AGENT, 'to': PEER, 'intent': 'schedule.query',
              'body': body, 'deliver': 'agent', 'priority': 'normal', 'evidence': [], 'reply_to': None, 'ext': {}})).digest())
    def send(envelope, sender, claim=None):
        signature = keys[sender].sign(b'SILENT-AGENT-RELAY\x00relay.v1\x00Ed25519\x00key1\x00' + canonical(envelope))
        envelope['sig'] = {'alg': 'Ed25519', 'kid': 'key1', 'value': encoded(signature)}
        headers = {'Content-Type': 'application/json', 'Authorization': 'Bearer ' + tokens[sender]}
        if claim:
            headers['X-Execution-Claim'] = claim
        request = urllib.request.Request(BASE + '/v1/send', data=json.dumps(envelope).encode(), headers=headers)
        with urllib.request.urlopen(request) as response:
            json.load(response)
            events.append({'api': '/v1/send', 'status': response.status, 'intent': envelope['intent']})
    send(parent, AGENT)
    receive(tokens, PEER, '/v1')
    claim = api('/v1/claim', {'id': REQUEST}, tokens[PEER])['claim']
    send({'v': 'relay.v1', 'id': GATE, 'from': PEER, 'to': PEER, 'intent': 'relay.approval.request',
        'body': {'reason': 'judgment_required', 'request_digest': digest}, 'deliver': 'human', 'exp': exp,
        'idempotency_key': 'designer-gate-request-key', 'reply_to': REQUEST}, PEER, claim)


def login(browser, session, who):
    browser.navigate(session, BASE)
    address = ('ui-member' if who == 'a' else 'other-member') + '@example.invalid'
    browser.evaluate(session, f"document.querySelector('#email').value={json.dumps(address)};document.querySelector('form').requestSubmit()")
    time.sleep(.3)
    messages = sorted((ROOT / 'mail').glob('*.eml'), key=lambda p: p.stat().st_mtime)
    message = email.message_from_bytes(messages[-1].read_bytes(), policy=email.policy.default)
    code = re.search(r'코드: (\d{6})', message.get_content()).group(1)
    browser.evaluate(session, f"document.querySelector('#code').value={json.dumps(code)};document.querySelector('form').requestSubmit()")
    time.sleep(.3)
    assert browser.evaluate(session, 'location.pathname') == '/home'
    events.append({'browserLogin': who, 'realStartVerify': True, 'smtp': 'private local sink'})


def capture(browser, session, name, path, mobile_only=False):
    browser.navigate(session, BASE + path)
    for width in ([390] if mobile_only else [1280, 390]):
        label = name + ('-mobile' if width == 390 else '-desktop')
        metrics = browser.capture(session, IMAGES, label, width)
        captures.append({'name': label, 'path': path, 'sha': '09c523da8a3407288d9f5d711e1834af12bc7808', **metrics})


def main():
    IMAGES.mkdir(parents=True, exist_ok=True)
    with socket.socket() as listener:
        listener.bind(('127.0.0.1', 0))
        CFG['smtp_port'] = listener.getsockname()[1]
    sink = subprocess.Popen(['python3', str(ROOT/'snapshot/scripts/mail_sink.py'), str(ROOT/'mail'),
               '--port', str(CFG['smtp_port'])], stdout=(ROOT/'smtp.log').open('a'), stderr=subprocess.STDOUT)
    CFG['smtp_pid'] = sink.pid
    sql("UPDATE relay_state SET data='{}',epoch=epoch+1,clock=clock_timestamp();")
    mutate(lambda st: None, 'initial unchanged host relay startup')
    browser = Browser()
    try:
        member, other = browser.new(), browser.new()
        login(browser, member, 'a')
        login(browser, other, 'b')
        receipt = '/home/receipts?' + urllib.parse.urlencode({'agent': AGENT, 'id': REQUEST})
        for name in ['queued', 'leased', 'delivered', 'reply-queued', 'reply-delivered', 'expired', 'revoked', 'max-attempts', 'rate']:
            keys, tokens = seed()
            send_text(keys, tokens)
            if name == 'leased':
                receive(tokens, PEER, ack=False)
            if name in ['delivered', 'reply-queued', 'reply-delivered']:
                receive(tokens, PEER)
            if name.startswith('reply-'):
                send_text(keys, tokens, reply=True)
            if name == 'reply-delivered':
                receive(tokens, AGENT)
            if name == 'expired':
                mutate(lambda st: st['Messages'][REQUEST]['Receipt'].update(
                    exp=instant(datetime.now(timezone.utc)-timedelta(seconds=1)),
                    accepted_at=instant(datetime.now(timezone.utc)-timedelta(seconds=181))),
                    'controlled text expiration with consistent 180-second acceptance; actual sweep')
            if name == 'revoked':
                mutate(lambda st: st['Agents'][AGENT]['Keys']['key1'].update(Revoked=True), 'controlled key revocation; actual sweep')
            if name == 'max-attempts':
                receive(tokens, PEER, ack=False)
                mutate(lambda st: st['Messages'][REQUEST].update(Attempts=3, LeaseUntil=instant(datetime.now(timezone.utc)-timedelta(seconds=1))), 'three leases elapsed; actual sweep')
            if name == 'rate':
                mutate(lambda st: st['Rates'].update({'http:new': [instant()] * 200}), 'controlled new-work rate saturation')
            capture(browser, member, 'receipt-' + name, receipt)
            if name == 'reply-delivered':
                capture(browser, member, 'reply-parent-id', receipt.replace(REQUEST, REPLY))
            if name == 'queued':
                capture(browser, member, 'home-tool-guidance', '/home')
                capture(browser, other, 'foreign-receipt-refusal', receipt)
                focus = []
                browser.navigate(member, BASE + receipt)
                for _ in range(3):
                    browser.call('Input.dispatchKeyEvent', {'type': 'keyDown', 'key': 'Tab', 'code': 'Tab', 'windowsVirtualKeyCode': 9}, member)
                    browser.call('Input.dispatchKeyEvent', {'type': 'keyUp', 'key': 'Tab', 'code': 'Tab', 'windowsVirtualKeyCode': 9}, member)
                    focus.append(browser.evaluate(member, "({tag:document.activeElement.tagName,text:document.activeElement.innerText})"))
                events.append({'receiptKeyboardTab': focus})
                # The browser form performs real GET refresh, preserving agent/request selection.
                browser.evaluate(member, "document.querySelector('button').click()")
                time.sleep(.2)
                capture(browser, member, 'receipt-manual-refresh', receipt, True)
                send_text(keys, tokens, content='다른 내용', expected=409)
                send_text(keys, {**tokens, AGENT: encoded(os.urandom(32))}, expected=401)
        for name in ['pending', 'approved', 'denied', 'expired', 'revoked', 'unavailable', 'invalid-signature', 'csrf-error', 'capacity-error']:
            keys, tokens = seed()
            gate(keys, tokens)
            path = '/home/gates/' + GATE
            if name == 'expired':
                mutate(lambda st: st['Gates'][GATE].update(Exp=instant(datetime.now(timezone.utc)-timedelta(seconds=1))), 'controlled gate expiration; actual sweep')
            if name == 'revoked':
                mutate(lambda st: st['Agents'][AGENT]['Keys']['key1'].update(Revoked=True), 'controlled gate key revocation; actual sweep')
            if name == 'unavailable':
                mutate(lambda st: st['Messages'][REQUEST].pop('Envelope', None), 'remove parent raw body; actual unavailable gate')
            if name == 'invalid-signature':
                mutate(lambda st: st['Messages'][REQUEST].update(Envelope={'invalid': 'signature'}), 'invalid signed parent fixture; actual verification failure')
            browser.navigate(member, BASE + path)
            if name == 'pending':
                capture(browser, other, 'foreign-gate-refusal', path)
                focus=[]
                for _ in range(3):
                    browser.call('Input.dispatchKeyEvent', {'type':'keyDown','key':'Tab','code':'Tab','windowsVirtualKeyCode':9}, member)
                    browser.call('Input.dispatchKeyEvent', {'type':'keyUp','key':'Tab','code':'Tab','windowsVirtualKeyCode':9}, member)
                    focus.append(browser.evaluate(member,"({tag:document.activeElement.tagName,text:document.activeElement.innerText})"))
                events.append({'gateKeyboardTab':focus})
            if name in ['approved', 'denied']:
                browser.evaluate(member, "document.querySelector('button[value=\"" + ('approve' if name == 'approved' else 'deny') + "\"]').click()")
                time.sleep(.2)
                if name == 'denied':
                    events.append({'denyRedirectPath':browser.evaluate(member,'location.pathname')})
                    capture(browser, member, 'gate-deny-redirect', browser.evaluate(member,'location.pathname'))
            if name == 'csrf-error':
                browser.evaluate(member,"document.querySelector('input[name=csrf]').value='invalid';document.querySelector('button[value=approve]').click()")
                time.sleep(.2)
            if name == 'capacity-error':
                mutate(lambda st: st['HTTP'].update({str(i):{'Clean':False,'Exp':instant(datetime.now(timezone.utc)+timedelta(seconds=20))} for i in range(16)}), 'controlled 16 shared new-work slots')
                browser.evaluate(member,"document.querySelector('button[value=approve]').click()")
                time.sleep(.2)
            # Error page is captured before a fresh GET can replace it.
            if name in ['csrf-error', 'capacity-error']:
                for width in [1280,390]:
                    label='gate-'+name+('-mobile' if width==390 else '-desktop')
                    metrics=browser.capture(member,IMAGES,label,width)
                    captures.append({'name':label,'path':path,'sha':'09c523da8a3407288d9f5d711e1834af12bc7808',**metrics})
            else:
                capture(browser,member,'gate-'+name,path)
        (EVIDENCE/'observations.json').write_text(json.dumps({'events':events,'captures':captures},ensure_ascii=False,indent=2)+'\n')
    finally:
        browser.close()
        os.kill(CFG['relay_pid'], signal.SIGTERM)
        sink.terminate()
        sink.wait(timeout=10)


if __name__ == '__main__':
    main()
