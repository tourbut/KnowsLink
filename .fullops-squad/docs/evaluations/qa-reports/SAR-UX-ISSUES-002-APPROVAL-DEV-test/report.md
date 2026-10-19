---
title: SAR-UX-ISSUES-002-APPROVAL-DEV DEV 보고
status: review
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-APPROVAL-DEV]
summary: "검증된 승인 요약(U5)·승인 화면 위치(U7 해당)·O-NM-04·O-NM-08의 기술 계획, 대응표, 검사 결과와 미실행 항목"
---

# SAR-UX-ISSUES-002-APPROVAL-DEV DEV 보고

## 대상과 결론

대상 제품 후보는 fullops/dev `f43d05cc153fc125274f7c6188aa3dc92cbbb8fa`다(첫 구현 `47117a4`, 테스트 보정 `74bf643`, 리뷰 수정 `f43d05c`). 기준은 main `6883050`(지시서 base `688305054dc4`)이다. 원격 Linux 격리 clone의 make lint·make test·integration과 fixture 컴파일이 모두 EXIT=0이다. DEV 완료이며 제품 수락이 아니다. tester 독립 QA·designer 직접 시각 검수·고정 SHA 리뷰·통합은 coor가 진행한다.

## 기술 계획

- 표시 전용 helper `internal/relay/gate_summary.go`를 추가한다. `gateState`는 저장 상태와 현재 검증 결과로 표시 상태를 정한다. `gateSummary`는 검증된 부모 envelope의 schema 필드만 사용자 문구로 바꾼다. `gateRows`는 승인 목록 행을 만든다.
- 결정 권한은 기존 `gateDecision`의 owner·CSRF·pending·기한·`verifiedGate` 재검사만 판단한다. 이번 변경은 이 재검사를 줄이지 않는다. `verifiedGate`는 지원하지 않는 부모 intent를 검증 실패로 닫는다(현재 gate 부모는 `schedule.query`·`schedule.commit`만 가능하므로 동작 변화 없음).
- 이름은 기존 `visibleAgent` 권한을 그대로 쓴다. 이름 권한이 없으면 ‘확인할 수 없는 에이전트’와 이미 화면에 있던 ID의 짧은 구분값만 보인다.
- 새 DB schema·migration·공개 API·route·JS·색 토큰은 없다. CSS는 O-NM-08의 `.deny-help` 간격 한 규칙만 추가했다.

## 요청 유형별 검증 필드와 사용자 문구 대응표

| 부모 intent | 검증된 필드 | 사용자 문구(요청한 내용) | 승인하면 허용되는 범위 | 기술 상세 |
|---|---|---|---|---|
| `schedule.query` | `body.window.start/end`(RFC3339 Z, end>start), `body.granularity_min`(양의 정수) | 요청 유형 ‘일정 가능 시간 확인’, 확인할 기간 ‘시작 KST – 끝 KST’, 시간 단위 ‘N분’ | 고정 문구 + ‘가능 시간 공개 정책이 없어 승인해도 상대에게 가능 시간이 공개되지 않습니다.’(policy: disclosure absent) | 원문 JSON, 원요청 ID, 승인 요청 ID, 두 agent ID, intent, policy |
| `schedule.commit` | `body.slot.start/end`, `body.timezone`(로드 가능한 IANA), `body.commitment`(비지 않은 문자열) | 요청 유형 ‘일정 확정’, 일정 ‘시작 KST – 끝 KST’, 시간대 원값. `commitment`는 ‘요청자가 적은 내용’으로 분리 | 고정 문구 + ‘승인해도 일정이 확정되거나 실행되지 않습니다. 이 서비스에는 일정 실행 기능이 없습니다.’(policy: non-executable stub) | 같음 |
| (행 이름 보정) | | `schedule.commit`의 시간대 행 이름은 ‘요청에 적힌 시간대’다. 일정 시각은 KST로 보이므로 원값이 표시 기준으로 읽히지 않게 했다. | | |
| 그 밖(검증 실패·원문 부재·digest/서명/endpoint/세대 불일치·미지원 intent) | 없음 | 요약 없음. ‘승인할 내용을 확인할 수 없습니다…’ alert | 표시하지 않음 | 표시하지 않음 |

공통: ‘누가 누구에게’는 gate의 From/To에 `visibleAgent` 권한 이름·종류·구분값을 쓴다. `render.hint`는 ‘참고 문구 — 승인 근거 아님’ 접힘 영역에 escape된 텍스트로만 둔다. 자유문과 hint는 링크로 만들지 않는다. 승인 요청 envelope의 `reason`은 이번 요약에 쓰지 않았다(지시 범위 밖, 기술 판단으로 생략).

표시 상태: 저장 상태가 pending이 아니면 그 값, 기한이 지나면 `expired`(만료), 검증 실패면 `unavailable`(확인할 수 없음), 검증은 되지만 현재 관계·키·owner가 current가 아니면 `revoked`(철회), 그 밖은 `pending`(대기)이다. 기존 housekeeping도 같은 조건을 저장 상태로 기록한다.

## 화면 변경

- 승인 상세 `/home/gates/{id}`·`/owner/gates/{id}`: 위치(`승인 목록 / 요청 승인`, aria-current) → h1 → 누가 누구에게 → 요청한 내용(요청자가 적은 내용·참고 문구) → 승인하면 허용되는 범위 → 기한·상태 → ‘이 요청 승인’·‘거절’ 또는 status 안내 → 기술 상세 → ‘승인 목록으로 돌아가기’. 검증 불가면 alert와 기한·상태·복귀만 보인다.
- 승인 목록(`/home/legacy#gates`): 행마다 요청한 에이전트·받은 내 에이전트·요청 유형·기한·상태와 ‘요청 검토’ 링크. 링크 접근성 이름은 ‘요청 검토: 유형 / 요청 에이전트 이름·종류·구분값 / 기한’이다. 원요청 ID는 상세의 기술 상세로 옮겼다. 빈 목록은 설계 문구를 쓴다. 대화 화면의 ‘승인 N건’은 표시 상태 pending만 센다.
- 회원 gate 거부 화면은 ‘승인 목록으로 돌아가기’(`/home/legacy#gates`)로 복귀한다. POST 거부는 ‘이번 제출은 적용되지 않았습니다…’, GET 거부는 ‘이 승인 요청을 볼 수 없습니다…’다. owner Basic 경로는 기존 응답을 유지한다.
- O-NM-04: 익명 429 조회 화면은 선택기·조회 폼·하단 링크 없이 재시도 시각과 ‘로그인 화면으로 이동’(`/`)만 보인다. ‘저장된 입력… 복원’은 세션 쿠키가 있는 브라우저에만 붙는다. 빈 `href`를 없앴다(`/home/receipts`로 대체). 렌더 fixture 전체에 `href=""` 금지 검사를 더했다.
- O-NM-08: 받은 친구 요청의 거절 안내에 `deny-help`(위 16px, 아래 4px)를 붙여 ‘거절’과 묶었다. 문구·`aria-describedby`는 그대로다.

## U5-A 자동 검사와 실패 경계

`internal/relay/gate_summary_integration_test.go` `TestApprovalSummaryBoundaries`(실제 Postgres, f43d05c PASS):

| 하위 검사 | 확인 내용 |
|---|---|
| `query_summary_matches_verified_body` | 읽기 순서, 요약 값이 서명된 body(window·granularity)와 일치, 같은 값이 기술 상세 JSON에 있음. 악성 hint(`<script>`, URL, ‘모든 권한 허용’)가 escape되어 참고 영역에만 1회. 목록에 이름·유형·‘요청 검토’가 있고 원요청 ID는 없음. agent 이름 변경 뒤 같은 원요청·gate ID 유지 |
| `commit_free_text_is_separate_and_inert` | HTML 링크가 든 commitment가 escape된 ‘요청자가 적은 내용’에만 있고 허용 범위 절에 없음, 링크 미생성 |
| `unverifiable_states_offer_no_decision` | 바뀐 digest·원문 부재·관계 철회·만료 각각 결정 버튼 없음, 상태 표시(확인할 수 없음/철회/만료), 승인 목록 절(`section#gates`)도 같은 상태이고 대기 표시 없음, 낡은 화면의 approve POST 409와 승인 목록 복귀, 결정 미기록 |
| `other_owner_stale_form_and_double_click` | 다른 회원 GET 403(이름·body 비노출), 다른 회원 POST 403, 다른 세션 CSRF 403, 한 gate의 token을 다른 gate에 제출하면 403이고 두 gate 모두 대기 유지, 승인 303 뒤 중복 클릭(deny) 409·‘이번 제출은 적용되지 않았습니다’, 첫 결정 유지·새 버튼 없음 |

기존 회귀: `TestGateFailureStates` 상태별 UI(거절·만료·철회·확인할 수 없음)와 gate-consume 409, `TestPublicHTTPAdmissionAndGateSafety`(hint XSS·잘못된 서명 GET/POST·admission), `TestEmailIdentity/member_gate_*`(회원 격리·결정·deny 복귀)가 통과했다. 바꾼 기대 문구는 ‘참고 hint’→‘참고 문구 — 승인 근거 아님’, ‘원문 부재’→‘승인할 내용을 확인할 수 없습니다’, 버튼 ‘승인’→‘이 요청 승인’, 복귀 `href="/home"`→`/home/legacy#gates`다. 실제 외부 업무 실행은 켜지 않았다(`executable:false`·`disclosure:false` 유지).

첫 integration 실행(47117a4)의 실패 두 건은 테스트 기대 문제였다. 기존 housekeeping이 원문 부재 gate를 저장 상태 `unavailable`로 바꾸므로 ‘pending 유지’ 대신 ‘결정 미기록’으로 고쳤다. 기존 테스트의 버튼 문구 기대를 새 문구로 고쳤다. 74bf643 첫 실행의 `TestPublicNodeProcesses` 실패는 새 clone에 `adapters/dist`가 없어서였다(make test가 생성). 순서대로 다시 실행해 통과했다.

## 검사 로그

원격 Linux 격리 clone(`/tmp/knowslink-approval002-<sha12>-<rand>`, 자기 Postgres 컨테이너, [remote.py](remote.py))에서 실행했다. 로그는 이 폴더에 있고 SHA256은 [raw-log-sha256.txt](raw-log-sha256.txt)다.

| SHA | make lint | make test | integration | fixture-compile |
|---|---|---|---|---|
| f43d05c | EXIT=0 | EXIT=0 | EXIT=0 (`-skip TestRemoteConsentBrowserReturn`, 기존과 같음) | EXIT=0 (`go vet`·`go test -c -tags=integration,approvalfixture`) |
| 74bf643 | EXIT=0 | EXIT=0 | EXIT=0 | EXIT=0 |
| 47117a4 | EXIT=0 | EXIT=0 | EXIT=1(위 기대 문제) | EXIT=0 |

로컬 Windows: `go vet`(태그 없음·integration·approvalfixture), `go test ./internal/relay` PASS.

## 렌더 확인

`KNOWSLINK_UI_FIXTURE_DIR`로 만든 합성 HTML을 headless Chrome 155(비모바일 viewport)로 1440·390·320에서 렌더했다. 스크립트 [render/probe.mjs](render/probe.mjs), 결과 [render/probe.json](render/probe.json), 캡처 `render/*.png` 24장이다. 대상은 gate-active(이름 있는 요청자), gate-commit-hint(39자 이름·긴 자유문·악성 hint·긴 JSON), gate-expired, gate-unavailable, home(승인 목록), home-no-gates(빈 목록), human-friends(O-NM-08), receipt-anonymous-limited(O-NM-04)다. 24개 모두 페이지 가로 넘침 0이다. DEV 확인이며 designer 시각 수락이 아니다.

## tester/designer용 격리 fixture

`internal/relay/approval_browser_fixture_test.go` `TestApprovalBrowserFixture`(태그 `integration && approvalfixture`). 기동은 coor가 `remote.py fixture`로 한다(자기 Postgres 컨테이너, `APPROVAL_FIXTURE_MANIFEST`).

- 계정: 합성 A(agent_b 소유, 승인자), B(agent_a 소유, 요청자), C(gate 없음). `POST /fixture/as/{a|b|c}`로 세션을 받는다(동일 origin만).
- gate 7개(manifest `gates`): `pending_query_hint`(악성 hint), `pending_commit_long`(긴 자유문·링크), `pending_decide_1`, `pending_decide_2`(승인·거절·중복 클릭용), `expired`, `processed`(승인됨), `original_missing`(원문 부재).
- A와 B는 친구가 아니라 legacy pair다. 친구로 바꾸면 부모 메시지의 scope가 바뀌어 철회가 되므로 요청자는 ‘확인할 수 없는 에이전트 / 구분값’으로 보인다. 이름 있는 요청자 화면은 렌더 fixture `gate-active`로 확인한다.
- C→A 친구 요청 1건이 대기 중이다(O-NM-08 확인용).
- origin: loopback `httptest` 주소(manifest `origin`). 허용 GET: `/`, `/home`, `/home/legacy`, `/home/friends`, `/home/account`, `/home/gates/*`. 허용 POST: `/home/gates/*`(승인·거절)와 `/fixture/as/*`. 그 밖은 403이다.
- 수명: 12시간. 저장된 기한만 늘리며 서명·digest·세대는 보낸 그대로다. 외부 업무 실행·메일·실계정은 없다.

## 하위 리뷰와 반영

같은 워크트리의 읽기 전용 리뷰 하위 에이전트 1건(Purpose: review)이 278001a..74bf643을 검토했다. critical·high는 없었다. medium 1건(철회 사례의 목록 단언이 관계 절의 ‘철회’ 배지로도 통과)과 low 4건(교차 gate token 미검사, commit 시간대 행 오독, 목록 유형을 policy에서 읽어 확인 불가 fallback 불가, 결정된 gate까지 서명 검사)을 모두 `f43d05c`에서 고쳤다. 목록 행은 대기 gate만 검증하고 유형을 검증된 envelope에서 읽는다. 검증 불가면 ‘확인할 수 없는 요청’이다. 이 리뷰는 coor의 독립 고정 SHA 리뷰를 대신하지 않는다.

## 미실행(not_run)

- 실제 브라우저 키보드(Tab/Shift+Tab/Enter/Space)·실제 200% 확대·스크린리더 낭독·보조 브라우저 확인(U7-A): tester/designer 몫이며 DEV는 실행하지 않았다.
- fixture 실제 기동과 브라우저 왕복: coor 기동 뒤 tester가 실행한다. DEV는 컴파일만 확인했다.
- `TestRemoteConsentBrowserReturn`: 기존 원격 실행과 같이 제외했다(이번 변경과 무관한 remote consent 브라우저 검사).
- 헤더 ‘내 계정’→‘설정’ 이름 변경(U7 일반)과 언어(U3)·연결(U8)·O-NM-03: 지시서 범위 밖이다.
- FullOps lint 결과는 완료 보고에 기록한다.

## FullOps lint

`lint.py --repo . --from 688305054dc4`를 커밋된 HEAD에서 실행했다. 미추적 사용자 폴더(`videos/service-intro/`, `SAR-MOTION-INTRO-001-DEV-test/` 하위 다섯 폴더)는 지시대로 건드리지 않고, 그 실행 프로세스에만 `core.excludesFile`(환경 변수 `GIT_CONFIG_*`)로 제외해 작업 트리 깨끗함 검사를 통과시켰다. 저장소 설정과 폴더는 바꾸지 않았다.

- ERROR 0. 첫 실행의 DOC-003(이 보고의 front matter 표기)은 `deliverables.py --stamp`로 고쳤다.
- WARNING SIZE-001 `internal/relay/http.go` 688줄(이전 674): gateView·gateError 보강 14줄이다. 요약 로직은 새 파일 `gate_summary.go`로 분리했고 기존 파일 분할은 이번 범위 밖이다.
- WARNING SIZE-002 변경 추가 약 500줄(상한 400): 지시서 예상 300–700줄 범위다. 제품 코드는 약 150줄이고 나머지는 U5-A integration 검사·격리 fixture·렌더 fixture다.
- WARNING SIZE-001 `PLANS.md`: 이번 과제에서 수정하지 않은 coor 소유 파일이다.
- 실행 불가 product-lint·product-test: Windows에 make가 없다. 같은 명령을 원격 Linux 격리 clone에서 f43d05c 기준 EXIT=0으로 실행했다(위 검사 로그). 이후 커밋은 문서·보고뿐이다.
- 의존성 변경 없음(go.mod·package.json 미수정). DEP-001 해당 없음.
