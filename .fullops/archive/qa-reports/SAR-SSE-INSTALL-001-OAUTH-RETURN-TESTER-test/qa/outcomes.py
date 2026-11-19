"""Tester packet outcomes for SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER from the packet's non-optional path/category pairs."""
import json, pathlib, re
ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/tester')
KEY = 'SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER'
JEV = ROOT / '.fullops-squad/docs/evaluations/jev'
packet = json.loads((JEV / f'{KEY}-packet.json').read_text(encoding='utf-8'))
F = '.fullops-squad/'
READ = {
    F + 'FULLOPS.md': '직접 읽었다. 독립 리뷰(구현자와 다른 세션, 고정 SHA detached snapshot)와 tester 책임, 미해결 C/H 차단 규약을 적용했다.',
    F + 'handovers/to_tester.md': '직접 읽고 완료 기준 1~6을 따랐다. DEV 최종 SHA eac603b와 review base c399a5e를 구분했다.',
    F + 'project.md': '직접 읽었다. 검증 명령(make lint·make test), 독립 QA와 designer 시각 검수 구분, 보안 기준을 적용했다.',
    F + 'rules/common/README.md': '직접 읽었다. fullops-common-0.3.3 적용 기준을 report.md에 기록했다.',
    F + 'rules/common/coding-style.md': '직접 읽었다. 단순 설계·오류 경계·UI 정본 기준으로 변경을 검토했다.',
    F + 'rules/common/security.md': '직접 읽었다. 출력 이스케이프(html/template), CSRF, 오류 응답 비노출, 비밀값 비기록을 검증했다.',
    F + 'rules/common/testing.md': '직접 읽었다. lite 레벨과 담당·반복 범위, 재사용 증거와 새 실행의 구분을 적용했다.',
    F + 'docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md': '직접 전문을 읽고 원인·수정·검증·한계 주장을 독립 실행(Chromium RED/GREEN, 경계, lint/test)으로 대조했다.',
    F + 'docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/dev-eac603b6-final-gate.json': '직접 읽었다. head eac603b·product-lint/test 종료코드 0·ERROR 0/WARNING 3 확인. 같은 SHA에서 lint.py를 다시 실행해 head를 맞췄다(WARNING 2).',
    F + 'contexts/dev.md': 'DEV 변경 diff를 읽었다. 결정 기록이 구현과 일치한다.',
    F + 'docs/design-docs/interface-design.md': 'DEV 변경 diff를 읽었다. 200 반환 화면·401/403 HTML 설명이 독립 실행 결과와 일치한다.',
    F + 'docs/design-docs/module-design.md': 'DEV 변경 diff를 읽었다. remoteReturn·form-action 미확장·TEST_BROWSER 설명이 실제 코드와 일치한다.',
    F + 'docs/operations/ops-guide.md': 'DEV 변경 diff와 관련 줄(42~44)을 읽었다. 지시서의 conflict 분류에 따라 과거 배포 수락 기록을 이번 성공 근거로 쓰지 않았다.',
    F + 'docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md': '관련 줄(consent 경로·TestRemoteOAuthHTTP·미실행 항목)을 읽었다. 실제 Dots 미실행 상태를 성공으로 바꾸지 않았다.',
    F + 'PLANS.md': '관련 줄(1466~1477)만 읽었다. 운영 관찰(install1/token0/code 만료)이 독립 RED 재현과 일치한다.',
    'internal/relay/remote_auth.go': '고정 SHA snapshot에서 전문 중 consent·authorize·exchange 구간을 읽고 diff 검토했다.',
    'internal/relay/remote_browser_integration_test.go': '전문을 읽고 같은 SHA의 Linux 격리 PG·Chromium에서 독립 실행했다(PASS).',
    'internal/relay/remote_integration_test.go': '변경 diff와 TestRemoteOAuthHTTP 전체 흐름을 읽고 독립 실행했다(PASS).',
    'internal/relay/testdata/consent_browser.mjs': '전문을 읽었다. 출력에서 URL query를 지우는 것을 확인했다. SLOP-004는 결과 JSON 출력 채널이다.',
    'internal/relay/http.go': 'CSP·Cache-Control·Referrer-Policy 설정과 statusFor를 읽었다. 변경 없음을 diff로 확인했고 consent POST 응답 헤더를 Chromium에서 직접 비교했다.',
    'internal/relay/capacity.go': 'boundedHTTP 설정(no-store·CSP)과 consent budget 경로 줄을 읽었다. 변경 없음.',
    'internal/relay/google.go': 'Google 시작·remote attempt 소유 binding 줄을 읽었다. 변경 없음. TestGoogleHTTP·TestGoogleDeviceHTTP를 같은 SHA에서 독립 재실행했다(PASS).',
    'internal/relay/remote_test.go': 'unit 테스트의 attempt/exchange 사용 줄을 읽었다. make test(race)에서 같은 SHA 통과.',
    'internal/relay/device.go': '관련 줄만 확인했다. consent 반환 화면과 무관하며 변경 없음.',
    'deploy/knowslink/tunnel/public-ingress.yml': 'consent 경로 규칙(/remote/consent/<43자>)을 읽었다. POST/GET이 같은 경로라 반환 화면에 ingress 변경이 필요 없다.',
    'scripts/check_public_ingress.py': 'consent 경로 검사 줄을 읽었다. 경로 집합 변경 없음.',
}
OTHER_DOC = '전문을 재독하지 않았다. consent·form-action·authorization_failed·OAUTH-RETURN 키워드 검색 결과 이번 callback 복귀 판단과 무관한 과거 기록(0건 또는 승인 화면 시각 증거 한정)으로 판단했다. 과거 증거를 새 SHA 성공 근거로 쓰지 않았다.'
OTHER_CODE = '경로와 consent 키워드 검색 결과 이번 변경(remote consent handler)의 호출·응답 형식에 의존하지 않아 변경·검증 대상에서 제외했다.'
NOCHANGE_UPDATE = {
    F + 'PLANS.md': 'PLANS.md는 coor 소유다. tester는 갱신하지 않고 발견을 worker_done으로 보고한다.',
    F + 'docs/design-docs/module-design.md': '제품 설계서는 dev 소유다. tester는 수정하지 않는다. 내용은 실제 코드와 일치해 갱신 필요가 없다.',
    F + 'docs/operations/ops-guide.md': '운영 가이드는 dev/ops 소유다. tester는 수정하지 않는다. 운영 변경은 직접 적용하지 않고 coor에 보고한다.',
}
SEARCHED = {'document_read', 'impact_check', 'direct_edit'}
items = []
for item in packet['items']:
    if item.get('optional'):
        continue
    path = item['path']
    for cat in item['categories']:
        if cat == 'document_update':
            status, reason = 'no_change', NOCHANGE_UPDATE.get(path, '소유 밖 문서이며 변경이 필요하지 않다.')
        elif cat == 'direct_edit':
            status, reason = 'no_change', f'tester는 제품 소스를 수정하지 않는다(source 수정 금지). {READ.get(path, "")}'
        elif path in READ:
            status, reason = 'completed', READ[path]
        elif path.endswith(('.go', '.py', '.ts', '.mjs', '.yml', '.sh')) or path in ('Makefile',):
            status, reason = 'no_change', OTHER_CODE
        else:
            status, reason = 'no_change', OTHER_DOC
        items.append({'path': path, 'category': cat, 'status': status, 'reason': reason})
out = {
    'packet_input_sha256': packet['input_sha256'],
    'attempt': packet['attempt'],
    'items': items,
    'uncertainty_review': 'partial 패킷의 unknown 15건(Google 구현·일부 테스트·OPS-LIVE 본문 미전송 포함)은 로컬 원문을 직접 읽었다: remote_auth.go·remote_browser_integration_test.go·remote_integration_test.go·consent_browser.mjs·http.go·capacity.go·google.go 관련 구간을 고정 SHA snapshot에서 확인했고, 미전송을 검증 완료로 해석하지 않았다. 정적 검색의 한계(동적 참조·언어 서버 의미)는 consent·authorization_failed·remote/consent grep과 실제 Chromium 실행으로 보강했다. 실제 Dots·실계정 경로는 미검증으로 유지한다.',
}
(JEV / f'{KEY}-packet-outcomes.json').write_text(json.dumps(out, ensure_ascii=False, indent=1) + '\n', encoding='utf-8')
print(len(items), sum(1 for i in items if i['status'] == 'completed'))
