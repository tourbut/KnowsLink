---
title: dev 컨텍스트
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-CODEX-HOOK-QUOTE-001-DEV, SAR-SERVICE-UI-001-FIX-DEV, SAR-SETUP-001-DEV, SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-AUTO-RECEIVE-001-DEV, SAR-SSE-INSTALL-001-DEV, SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV, SAR-GROK-WEBHOOK-001-DEV, SAR-SERVICE-UI-001-DEV, SAR-FRONTEND-POLISH-001-DEV]
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
