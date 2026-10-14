"""Record the independent review coverage; no product or snapshot writes."""
import json
from pathlib import Path

HERE = Path(__file__).parent
result = json.loads((HERE / 'result.json').read_text(encoding='utf-8'))
code_reasons = {
    'google.go': 'googleFailure 공유 view 전환과 전체 googleStart/callback/finish 흐름을 읽었다. provider 설정, 재확인, connection/remote Next와 기존 오류 상태를 보존한다.',
    'human_pages.go': 'human-account의 서비스 안내 링크만 추가한다. 회원 메뉴·폼·개인정보의 기존 권한 경계는 유지한다.',
    'identity_integration_test.go': '공개 h1의 의도된 카피 변경만 기대값에 반영했다. 기존 인증·실패·CSRF·세션 회귀는 유지한다.',
    'landing_fixture_test.go': 'opt-in build tags, 합성 DB 전제, loopback 서버, 경로 allowlist, 로컬 Google 취소/메일 mock, 종료 채널을 확인했다. 운영 인증/송신은 허용하지 않는다.',
    'landing_integration_test.go': '실제 DB 익명/authed Cache-Control/CSP, 회원 인가·이름·만료·개인정보·provider 오류 회귀와 원 실행 PASS를 확인했다.',
    'landing_pages.go': 'D02 LG-01–08/D04 카피·README/installer·기존 memberStyle과 대조했다. 단일 인증 영역, 정적 합성 예시, .landing CSS, 텍스트 escaping·앵커/details 구조를 확인했다.',
    'landing_test.go': '4 provider 상태, 정확히 한 번의 활성 폼, 내용·오류 escape 검사를 확인했다. nil DB 전체 Handler 대신 startPage 단위 호출이며 헤더 assertion은 integration에 보존한다.',
    'member.go': 'startView 모든 호출자와 startPage의 세션 존재 여부만 전달하는 변경을 읽었다. 실제 설정·오류 반환·no identity rendering·기존 POST 계약을 확인했다.',
    'member_pages.go': 'start 정의를 landingTemplates로 옮기고 기존 memberStyle/form golden/공용 template 계약을 보존한다.'
}
for item in result['files']:
    path = item['path']
    name = Path(path).name
    status, reason = 'reviewed', '고정 SHA 변경과 관련 정본·원 실행 결과를 직접 대조했다. report.md에 상세 범위와 미실행을 기록했다.'
    if name.endswith('.png'):
        status, reason = 'skipped', '시각 판정은 별도 designer 담당이다. 17개 캡처의 존재·PNG 크기·SHA256은 captures.json과 전부 대조했다. native 입력/확대 수락을 추정하지 않는다.'
    elif '/jev/' in path and name not in ['SAR-LANDING-GUIDE-001-DEV-model-dev.json', 'SAR-LANDING-GUIDE-001-DEV-route.json', 'SAR-LANDING-GUIDE-001-DEV-packet-outcomes.json']:
        status, reason = 'skipped', '생성 탐색/과거 history 전체 발췌를 재감사하지 않았다. JSON 파싱·hash·현재 packet의 HEAD/attempt/partial/unknown22와 outcomes 결속을 확인하고 핵심 원문·호출자를 직접 읽었다.'
    elif 'SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/' in path and name != 'report.md':
        status, reason = 'skipped', '이미 수락된 별도 설계 리뷰의 역사 증거다. JSON 파싱·hash와 report의 세션/기준/19 reviewed·13 skipped/한계를 확인했다. 이번 제품 gate나 독립 UI 수락으로 재사용하지 않는다.'
    elif name == 'last-snapshot.json':
        status, reason = 'skipped', '생성 DOM snapshot 전체의 시각 의미는 재감사하지 않았다. JSON 파싱·hash, checks/실패/fixture/report와 실행 helper는 확인했다. 독립 UI QA를 대신하지 않는다.'
    elif name == 'source-audit.json':
        status, reason = 'skipped', '생성된 61개 과거 원문 발췌 전체는 재감사하지 않았다. 파싱·hash·원문 식별 구조를 확인하고 현재 D02/D04/README·인증 호출자는 직접 대조했다.'
    elif path.startswith('internal/relay/'):
        reason = code_reasons[name]
    elif name == 'browser.py':
        reason = 'DOM 실제 submit/click/focus와 native 입력·CSS 확대를 구분하는 실행 코드를 전부 읽었다. 합성 loopback만 사용하며 실패를 assertion으로 기록한다. 재사용 시 고정 page와 cookie 상태 제약은 QA에 있다.'
    elif name in ['remote.py', 'forward.py']:
        reason = '전체 helper의 exact HEAD/clean 검사, 신규 고유 clone·합성 DB, RejectPolicy·비밀값 비출력, loopback 포워드, stop 보존을 확인했다. 삭제/이동과 운영 자원 명령이 없다.'
    elif '/logs/' in path and name.endswith('.log'):
        reason = '원 로그를 직접 읽었다. 4795 실행의 실제 명령·PASS/음수 시나리오·EXIT=0을 확인했다. final fe의 필수 lint/test는 별도 원 로그로 다시 확인했다.'
    elif name == 'captures.json':
        reason = '17개 PNG의 파일 존재·실제 크기·SHA256을 독립 검증했다. verification.json에 보존했다. 직접 시각 수락은 별도다.'
    elif name == 'to_dev.md':
        reason = 'work.py finish 뒤 빈 인박스이며 동일 과제 attempt 지시와 완료 전문이 logs/2026-10-14_to_dev.md에 보존되어 있다.'
    item['review_status'], item['reason'] = status, reason
result['reviewer'] = 'Codex gpt-6-astra high; actual session 01a13aaa-968a-7d81-b8ca-f2d71a16a0d4'
result['conclusion'] = '코드·기술 문서 리뷰 수락 가능. C/H/M/L 0. final fe2387a 필수 lint/test exit0 확인. 독립 tester·designer 및 native keyboard/native200% 미실행은 제품 최종 수락 전 별도 게이트로 유지한다.'
result['findings'] = []
(HERE / 'result.json').write_text(json.dumps(result, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
raw = Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-dev/final-handoff.md').read_bytes()
(HERE / 'author-final-handoff.md').write_bytes(raw)
counts = {'total_files': len(result['files']), 'reviewed_files': sum(f['review_status'] == 'reviewed' for f in result['files'])}
counts['skipped_files'] = counts['total_files'] - counts['reviewed_files']
counts['coverage_rate'] = counts['reviewed_files'] / counts['total_files'] * 100
(HERE / 'coverage.json').write_text(json.dumps(counts, indent=2) + '\n', encoding='utf-8')
print(json.dumps(counts))
