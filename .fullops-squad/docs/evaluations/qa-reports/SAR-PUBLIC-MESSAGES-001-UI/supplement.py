"""Verify deny navigation and normalize the expired receipt fixture without repeating unaffected captures."""

import json
import os
import signal
import socket
import subprocess
import time
from datetime import datetime, timedelta, timezone

import run as fixture


def main():
    with socket.socket() as listener:
        listener.bind(('127.0.0.1', 0))
        fixture.CFG['smtp_port'] = listener.getsockname()[1]
    sink = subprocess.Popen(['python3',str(fixture.ROOT/'snapshot/scripts/mail_sink.py'),str(fixture.ROOT/'mail'),
                '--port',str(fixture.CFG['smtp_port'])],stdout=(fixture.ROOT/'smtp.log').open('a'),stderr=subprocess.STDOUT)
    fixture.CFG['smtp_pid']=sink.pid
    fixture.sql("UPDATE relay_state SET data='{}',epoch=epoch+1,clock=clock_timestamp();")
    fixture.mutate(lambda st: None,'supplement fixture startup')
    browser=fixture.Browser()
    try:
        member=browser.new()
        fixture.login(browser,member,'a')
        keys,tokens=fixture.seed()
        fixture.send_text(keys,tokens)
        fixture.mutate(lambda st: st['Messages'][fixture.REQUEST]['Receipt'].update(
            exp=fixture.instant(datetime.now(timezone.utc)-timedelta(seconds=1)),
            accepted_at=fixture.instant(datetime.now(timezone.utc)-timedelta(seconds=181))),
            'normalize expired acceptance metadata to TTL180; initial negative fixture kept')
        path='/home/receipts?agent='+fixture.AGENT+'&id='+fixture.REQUEST
        fixture.capture(browser,member,'receipt-expired',path)
        keys,tokens=fixture.seed()
        fixture.gate(keys,tokens)
        path='/home/gates/'+fixture.GATE
        browser.navigate(member,fixture.BASE+path)
        browser.evaluate(member,"document.querySelector('button[value=deny]').focus()")
        for kind in ['keyDown','keyUp']:
            browser.call('Input.dispatchKeyEvent',{'type':kind,'key':'Enter','code':'Enter','windowsVirtualKeyCode':13,
                         'text':'\r' if kind=='keyDown' else ''},member)
        time.sleep(.3)
        navigation=browser.evaluate(member,"({path:location.pathname,status:performance.getEntriesByType('navigation')[0].responseStatus,text:document.body.innerText})")
        navigation['storedGateState']=fixture.state()['Gates'][fixture.GATE]['State']
        navigation['networkResponses']=browser.responses[-6:]
        (fixture.EVIDENCE/'supplement-navigation.json').write_text(json.dumps(navigation,ensure_ascii=False,indent=2)+'\n')
        assert navigation['text'].strip()=='Method Not Allowed' and navigation['storedGateState']=='denied'
        assert any(r['status']==405 and r['url'].endswith('/deny') for r in browser.responses)
        fixture.events.append({'denyKeyboardEnterResult':navigation})
        # Canonical gate is reachable when the user supplies its URL; browser error itself offers no recovery link.
        browser.navigate(member,fixture.BASE+path)
        fixture.events.append({'canonicalDeniedRecovery':browser.evaluate(member,"document.body.innerText.includes('상태: denied')")})
        browser.evaluate(member,"document.querySelector('a').click()")
        time.sleep(.2)
        fixture.events.append({'homeRecoveryPath':browser.evaluate(member,'location.pathname')})
        original=json.loads((fixture.EVIDENCE/'observations.json').read_text())
        for capture in original['captures']:
            if capture['name'] in ['receipt-expired-desktop','receipt-expired-mobile']:
                capture['name']+='-initial-fixture'
                capture['fixtureLimitation']='accepted_at was not shifted; negative TTL is fixture error, not product defect'
        original['captures']+=fixture.captures
        original['events']+=fixture.events
        (fixture.EVIDENCE/'observations.json').write_text(json.dumps(original,ensure_ascii=False,indent=2)+'\n')
        (fixture.EVIDENCE/'supplement-result.json').write_text(json.dumps(navigation,ensure_ascii=False,indent=2)+'\n')
    finally:
        browser.close()
        os.kill(fixture.CFG['relay_pid'],signal.SIGTERM)
        sink.terminate()
        sink.wait(timeout=10)


if __name__=='__main__':
    main()
