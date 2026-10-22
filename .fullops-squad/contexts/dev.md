---
title: dev 컨텍스트
status: draft
updated: 2026-10-21
owner: dev
tasks: [SAR-PUBLISH-001-DEV, SAR-SCHEDULE-001-DEV, SAR-CODEX-HOOK-QUOTE-001-DEV, SAR-SERVICE-UI-001-FIX-DEV, SAR-SETUP-001-DEV, SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-AUTO-RECEIVE-001-DEV, SAR-SSE-INSTALL-001-DEV, SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV, SAR-GROK-WEBHOOK-001-DEV, SAR-SERVICE-UI-001-DEV, SAR-FRONTEND-POLISH-001-DEV, SAR-REPO-DOCS-001-DEV, SAR-SYSTEM-DIAGRAM-001-DEV, SAR-REFACTOR-002-DEV, SAR-UX-ISSUES-002-RECEIPT-DEV, SAR-UX-ISSUES-002-APPROVAL-DEV, SAR-UX-ISSUES-002-CONNECT-DEV, SAR-UX-ISSUES-002-LOWS-DEV]
summary: dev의 유효 설계 원칙 요약과 최근 결정. 2026-10-10까지의 전문은 archive에 보관
---

# dev 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다. 2026-10-10까지의 전문은 [보관본](archive/dev-2026-10-10.md)에 있다. 과제별 근거는 각 `docs/exec-plans/phases/<과제키>.md`를 따른다.

## 유효 원칙 (보관본 요약)

- 검증: `make lint`·`make test`는 `cmd/`·`internal/`과 adapter만 대상이다. Windows는 KnowsLinkDevTools의 Node22·w64devkit을 쓴다. Windows product-lint는 deploy 스크립트 한계가 있어 최종 gate는 서버 고유 `/tmp` clean clone에서 실행한다.
- 실행권: claim은 한 번만 발급한다. 재시작·철회·세대 교체 뒤 과거 receipt·gate로 복구하지 않는다. 처리 경로(`deliver`)는 저장 상태에 남기고 lease·persist·ACK·claim 확정 지점마다 검사한다.
- 보안 게이트를 Python `assert`로 쓰지 않는다(`-O`에서 사라진다). 회귀 검사는 원 결함 대조군·변이로 검출력을 먼저 증명한다.
- timeout: inline `AbortSignal.timeout()`과 reader만 남은 `Response`는 GC로 회수될 수 있다. 실제 기본값과 강제 GC로 검사한다.
- rate·예약 budget: principal bucket을 공유 bucket보다 앞에 둔다. 예약 판정은 transaction과 같은 자기 기록 대조로 한다. 공정성 단위는 "누가 누구를 막을 수 있는가"로 나눈다. snapshot 순서는 (시각, epoch)다.
- 회원·키: grant·새키 PoP·owner 지문·client 1회 완료를 분리한다. 철회는 항상 허용하고 증가 상한은 신규 기록만 거부한다. 살아 있는 agent의 철회 키는 지우지 않는다.
- bundle: 새 CLI를 import하면 `import.meta.url` 조건이 bundle URL이 된다. CLI 파일명을 확인하고 공통 helper는 core에 둔다.
- Windows 비공개 파일은 소유자·ACL로 검사한다. Go 문자열 CSS는 DESIGN lint 대상이 아니므로 390×844 렌더링으로 폭을 확인한다.
- 증거 구분: 로컬 Node/MCP 왕복, 실제 외부 계정·Bot 왕복, 독립 QA·직접 UI 수락은 서로 다른 증거다.

## 최근 결정

- 2026-10-15 SAR-LANDING-GUIDE-001-DEV: 공용 목록 링크의 색상 규칙은 button variant를 제외한다. source f2b8cb6에서 회원 CTA 7.958:1·공개 예시 간격24px와 세 폭의 Chrome 검사를 확인했다. [재개 QA](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume/report.md).
- 기존 native keyboard/200% 미검증 기록은 tester5149ac1의 실제 PASS로 해소됐다. 이번 CSS 검증은 그 원 실행을 재사용하며 새 SHA의 키보드/200% 실행으로 바꾸어 기록하지 않는다.

- 2026-10-14 SAR-REFACTOR-002-DEV: 로컬 Inbox는 실패할 수 있는 housekeeping을 unread claim 전에 끝낸다. 만료 이력의 권한 순회를 생략해도 delivered 사본·24h metadata 계약은 유지한다. CPU-only 개선율과 실제 DB transaction 개선율은 별도 측정하며 [QA](../docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/report.md)에 실패한 실행도 보존한다.
- 2026-10-13 SAR-REFACTOR-001-DEV: 실제 template·외부 host·동적 등록을 확인한 뒤 공통 agent 선택/연결 응답/text 검증/MCP 포장만 공유했다. 제품 12줄 감소이며 legacy·held/test·보안·UI 계약은 유지한다. exact cd4dae8 Linux gate 0을 [QA 보고](../docs/evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md)에 보존했다. 중단 뒤 재개는 원 작성자 세션과 현재 세션을 분리하고 제품 변경 없이 기록을 마무리한다.
- 2026-10-12 SAR-FRONTEND-POLISH-001-DEV: 빈 상태는 쓸 수 없는 form을 숨기고 다음 동작 한 줄로 바꾼다(Go template `{{$x = true}}` 재할당). 측정 스크립트의 위치 판정은 padding을 보정한다. 테스트가 고정한 동작(받은 초대의 unpair)은 polish로 바꾸지 않고 제품 판단으로 남긴다. [실행 기록](../docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md).
- 2026-10-12 SAR-CODEX-HOOK-QUOTE-001-DEV: Codex Windows hook은 기본 `%COMSPEC% /C "<command>"` raw 인자이고 셸을 설정하면 PowerShell일 수 있다. hook 명령은 따옴표를 쓰고 `; exit $LASTEXITCODE`로 PowerShell의 종료코드 2→1 변환을 막는다. 테스트는 argv 목록이 아니라 실제 실행 방식(raw 문자열)으로 검증한다. [실행 기록](../docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-DEV.md).
- 2026-10-11 SAR-SERVICE-UI-001-DEV: UI 재설계 전에 재설계 전 템플릿으로 form 계약(method·action·hidden·필드·button 라우팅) golden을 먼저 커밋한다. 이후 diff가 의도한 변경만인지 기계적으로 보인다. `memberStyle`은 OAuth 동의·gate 화면도 공유하므로 토큰 변경 시 함께 렌더한다. Orca 내장 브라우저 키 입력은 페이지를 한 번 클릭해 포커스를 준 뒤에만 전달된다.
- 2026-10-11 SAR-GROK-WEBHOOK-001-DEV: webhook routine wake는 `KNOWSLINK_GROK_WEBHOOK_FILE`로만 켜고 그 process에서 `sendPrompt`와 배타다. 실패 fallback으로 이중 doorbell을 만들지 않는다. secret POST는 고정 host 정규식·`redirect: "manual"`로 다른 host 유출을 막는다. `AbortSignal.timeout`은 unref timer라 stub 검사에는 ref timer가 필요하다.
- 2026-10-11 SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV: Chromium은 form 제출의 redirect에도 form-action을 적용한다. 외부 OAuth callback은 POST 303 대신 동일 출처 반환 화면(meta refresh+링크)으로 이동한다. 이 경계는 HTTP recorder가 아닌 실제 Chromium(`TEST_BROWSER`)으로 검사한다. SSH 세션의 Linux Chrome에는 `--password-store=basic`이 필요하다.
- 2026-10-11 SAR-SSE-INSTALL-001-DEV: SSE heartbeat는 DB polling 없이 LISTEN으로 대기한다. boundedHTTP의 read/write deadline은 handshake 뒤 명시적으로 해제하며 별도 stream 수명·write deadline으로 제한한다.
- Dots는 공식 MCP2 Events·remote OAuth 경로로 연결한다. local stdio ZIP·Google local 등록과 구분하고 Cloud plugin은 실제 등록 app ID를 참조한다. 새 설치 키는 명시적 동의 뒤 private relay DB에 보관한다.
- 임시 Postgres readiness는 TCP로 검사한다. initdb의 잠깐 열린 Unix socket 성공을 최종 서비스 준비로 오인하지 않는다. 원격 검사는 고유 /tmp clone·자기 DB만 사용하고 기존 서비스 ID를 대조한다.

- 2026-10-10 SAR-AUTO-RECEIVE-001-DEV: 자동 pull→로컬 inbox→ACK→MCP 알림 후보 뒤, 선택형 loopback wake를 `grok-wake.ts`로 분리 구현했다. 미문서화 vendor 경로는 기본 off이고 지원 상태(공식 없음/미문서화/실제 미검증)를 문서에 구분한다.
- 사용자 권한 입력(`sendPrompt`)에는 받은 text를 넣지 않는다. 고정 doorbell+검증 ID만 보내고 본문은 untrusted 도구 결과로 읽힌다. 결과는 retry(확실한 미전달)와 uncertain(재전송 금지)으로 나눈다.
- 이전 보고의 금지 문구가 사용자 지시인지 확인한다. 아니면 정정한다. 실제 Bot 컴퓨터 확인은 coor/owner 후속이다. [실행 기록](../docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md).
- 2026-10-11 SAR-SERVICE-UI-001-FIX-DEV: 44px 터치 기준은 버튼·입력만이 아니라 `a[href]` 전체를 DOM으로 측정한다. 모든 회원·gate·OAuth 화면이 `memberStyle` 하나를 쓰므로 공통 `a` 규칙 한 곳에서 고친다. [실행 기록](../docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md).

- 2026-10-12 SAR-HUMAN-PEER-001-DEV: 친구·자기·legacy 인가는 출처와 세대를 함께 비교하고, 과거 사본 열람은 전송 권한과 분리한다. 구버전 writer 차단은 전체 과거 backup 복원을 보호하지 못한다. [실행 기록](../docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md).

- 2026-10-12 SAR-HUMAN-PEER-001-DEV 시각 수정: 두 줄 clamp는 뒤의 display:block에 무효화될 수 있으므로 실제 computed geometry로 확인한다. 브라우저 세션 설치 전에 isolated profile을 고정한다. [수정 기록](../docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/v01-v02-report.md).

- 2026-10-13 SAR-REPO-DOCS-001-DEV: 설치와 기존 업데이트는 키/agent/Inbox 보존 여부로 구분한다. Grok 로컬 재시작과 Dots 서버 요청 반영·metadata Refresh를 분리하고 수신/ACK·화면/답장을 구분한다. [결과](../docs/evaluations/qa-reports/SAR-REPO-DOCS-001-DEV/report.md)의 기존 heading alias·inbound 검사로 과거 문서 링크를 보존했다.

- 2026-10-14 SAR-SYSTEM-DIAGRAM-001-DEV: Dots receive의 transport ACK와 knowslink_ack의 설치 Inbox 읽음 제거를 구분한다. 로컬 미확인 Inbox와 서버 설치 Inbox의 만료도 다르다. [상세 문서 QA](../docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md). Mermaid 구조 추출과 실제 렌더 수락은 별도다.

- 2026-10-14 SAR-LANDING-GUIDE-001-DEV: 공개 `/`는 유효 세션에도200 안내이며 개인정보 대신 `/home` 링크만 표시한다. 설정별 로그인은 Service.startView 한 곳에서 결정한다. [DEV QA](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md).
- Orca 입력 성공 응답만으로 키보드 PASS를 선언하지 않는다. 이번 native 입력은 hasFocus=false로 미실행이며 CSS200%와 native200%도 분리했다. 같은 host의 포트는 cookie를 공유하므로 fixture 상태별 새 isolated profile을 준비한다.

- 2026-10-15 SAR-UX-ISSUES-002-RECEIPT-DEV: 조회 오류 폼은 route뿐 아니라 상위 admission의429/503에서도 보존해야 한다. 쿠키는 입력이며 현재 회원/소유권을 재확인한다. DB 인증 불가 때는 캐시로 사적 입력을 노출하지 않는다. [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md).

- 2026-10-19 SAR-UX-ISSUES-002-NAMES-DEV: 긴 이름 선택은 닫힌 select 대신 줄바꿈되는 radio로 해결한다. 상대 agent 이름은 현재 친구 권한에서만 보이고 그 밖은 ID만 둔다. CDP 측정은 모바일 emulation이 넘친 폭만큼 viewport를 넓히므로 비모바일 viewport로 가로 넘침을 잰다. [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md).

- 2026-10-19 SAR-UX-ISSUES-002-APPROVAL-DEV: 승인 요약은 표시 전용 helper(`gate_summary.go`)에 두고 결정 재검사는 `gateDecision`에만 남긴다. 기존 housekeeping이 만료·철회·원문 부재를 저장 상태로 바꾸므로 낡은 제출 검사는 ‘pending 유지’가 아니라 ‘결정 미기록’으로 단언한다. 새 원격 clone의 integration은 make test가 `adapters/dist`를 만든 뒤 실행한다. [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md). fixture가 claim된 부모를 여러 개 만들면 전역 claim 4 한도(`claimCapacity`)에 걸리므로 gate 생성 뒤 claim token을 잠시 비우고 끝에서 복원한다. fixture는 컴파일이 아니라 실제 기동으로 확인한다.

- 2026-10-20 SAR-UX-ISSUES-002-CONNECT-DEV: 연결 재개 상태는 표시 전용 helper(`connect_guide.go`)가 저장 상태만으로 계산하고, 관측할 수 없는 앱 설치·세션 연결·온라인은 ‘앱에서 확인 필요’로 둔다. 존재 비노출 표시는 저장 모양뿐 아니라 상대의 차단·반복 요청 뒤 화면까지 같은 모양인지 확인한다. Python `write_text`는 Windows에서 CRLF로 쓰므로 바이트로 쓴다. [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md).

- 2026-10-20 SAR-UX-ISSUES-002-LOWS-DEV: 응답 내부 표식은 render 한 곳이 아니라 진입 미들웨어의 ResponseWriter 래퍼(`Unwrap` 포함)에서 지워야 redirect·JSON까지 막는다. 보류 행처럼 상태를 바꿔 숨긴 행은 표시뿐 아니라 방향·결정 판정(block)도 이전 모습과 같게 맞춘다. 수정 전 실패는 원격 clone 사본에서 제품 파일만 기준 SHA로 되돌려 확인한다. [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LOWS-DEV-test/report.md).

- 2026-10-21 SAR-SCHEDULE-001-DEV: 외부 조회는 lock 안에서 lease(`Gate.Lookup.Started`)만 잡고 lock 밖에서 한 번 실행한 뒤 두 번째 transaction이 결과를 1회 기록한다. 재시작은 lease 만료로 `failed`에 닫고 재조회하지 않는다. agent에는 lookup이 끝나기 전까지 `pending`으로 보여 중간 상태를 노출하지 않는다. 결과 공개는 서버 canonical JSON과 바이트 동일한 echo만 수락해 agent가 필드를 더하지 못한다. 같은 Google project의 refresh token을 revoke하면 그 project의 모든 grant가 끊기므로 재연결(덮어쓰기) 때는 revoke하지 않는다. Makefile의 gofmt는 `$(go env GOROOT)/bin/gofmt`이며 PATH의 다른 Go 버전 gofmt와 정렬 결과가 다를 수 있다. [DEV 보고](../docs/evaluations/qa-reports/SAR-SCHEDULE-001-DEV-test/report.md).

- 2026-10-21 SAR-PUBLISH-001-DEV: 기능 제거는 그 과제 이전 커밋 이후의 변경이 그 과제뿐인지 `git log -- internal cmd`로 확인한 뒤 파일 단위로 되돌린다. 저장 상태 JSON의 필드를 바로 지우면 첫 transaction이 철회 없이 버리므로, legacy 필드를 남기고 HTTP 시작 전 정리 한 곳에서 철회·삭제한다. 언어별로 따로 쓴 문안은 run 번역 표 대신 embed 정적 HTML로 둔다. 한국어 조사가 링크 바로 뒤에 붙으면 영어 값에 앞 공백을 넣는다. headless Chrome 창은 최소 500px이라 390·320 측정은 `--allow-file-access-from-files`의 고정 폭 iframe으로 한다. [DEV 보고](../docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md).

- 2026-10-22 SAR-CODEX-RECV-001-DEV: Codex 세션 수신은 세션별 MCP 프로세스 결속 + MCP elicitation 승인으로 구현했다(`codex queue`·`turn/steer`·`inject_items`는 본문을 소유자 입력·대화 기록으로 넣어 제외). Codex YOLO는 elicitation을 보이지 않고 자동 수락하므로 단순 accept를 승인으로 쓰지 않고 대화상자의 일회 코드 입력을 요구한다. 호스트 승인 UI는 실제 기본 설정(YOLO 포함)으로도 실측한다. [DEV 보고](../docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md).
