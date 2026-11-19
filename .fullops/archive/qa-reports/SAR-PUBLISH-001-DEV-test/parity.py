"""PC-02: compare rendered /privacy and /terms text with the SAR-PUBLISH-001 KO/EN drafts (P-01..04 values filled)."""
import re, sys, html, difflib, pathlib
spec = pathlib.Path(sys.argv[1]).read_text(encoding='utf-8')
out = pathlib.Path(sys.argv[2])
heads = {'privacy.ko':'## 개인정보처리방침 — 한국어','privacy.en':'## Privacy Policy — English','terms.ko':'## 서비스 약관 — 한국어','terms.en':'## Terms of Service — English'}
fill = {'{시행일}':'2026-10-22','{effective date}':'2026-10-22','{보호책임자 표시}':'KnowsLink 운영자','{privacy officer}':'KnowsLink operator',
 '백업은 {백업 보존 기간}이 지나면 삭제합니다.':'백업은 정기적으로 교체하며 교체할 때 함께 삭제합니다.',
 'We delete backups after {backup retention period}.':'Backups are replaced regularly, and the information in them is deleted when they are replaced.'}
def words(s): return re.findall(r"\S+", s)
ok = True
for key, head in heads.items():
    a = spec.index(head) + len(head); b = spec.index('\n## ', a)
    sec = spec[a:b]
    for k, v in fill.items(): sec = sec.replace(k, v)
    lines = [l for l in sec.splitlines() if not re.match(r'^\|[-| ]+\|$', l.strip()) and '{이메일 발송 업체}' not in l and '{email provider}' not in l]
    sec = '\n'.join(lines)
    sec = re.sub(r'\*\*h1\*\* ', '', sec); sec = sec.replace('**', '')
    sec = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', sec); sec = re.sub(r'\[([^\]]+)\]', r'\1', sec)
    sec = sec.replace('`', '').replace('|', ' ')
    sec = re.sub(r'^- ', '', sec, flags=re.M)
    page = (out / ('policy-' + key.replace('.', '-') + '.html')).read_text(encoding='utf-8')
    body = page[page.index('<h1>'):page.index('</main>')]
    body = re.sub(r'<br>', '\n', body); body = html.unescape(re.sub(r'<[^>]+>', ' ', body))
    body = body.replace('시행일: ', '시행일: ').replace(' .', '.').replace(' ,', ',')
    sw, bw = words(sec), words(body)
    # rendered text splits a word around a link (e.g. "약관 도"); join spec/body by removing spaces for comparison
    sj, bj = ''.join(sw), ''.join(bw)
    if sj == bj: print(key, 'MATCH', len(sw), 'words')
    else:
        ok = False
        for op in difflib.SequenceMatcher(None, sw, bw).get_opcodes():
            if op[0] != 'equal': print(key, op[0], sw[op[1]:op[2]], '->', bw[op[3]:op[4]])
    print(key, 'placeholders', '{' in body or '}' in body)
sys.exit(0 if ok else 1)
