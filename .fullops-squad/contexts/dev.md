---
title: dev 컨텍스트
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-SETUP-001-DEV, SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-AUTO-RECEIVE-001-DEV]
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

- 2026-10-10 SAR-AUTO-RECEIVE-001-DEV: 자동 pull→로컬 inbox→ACK→MCP 알림 후보 뒤, 선택형 loopback wake를 `grok-wake.ts`로 분리 구현했다. 미문서화 vendor 경로는 기본 off이고 지원 상태(공식 없음/미문서화/실제 미검증)를 문서에 구분한다.
- 사용자 권한 입력(`sendPrompt`)에는 받은 text를 넣지 않는다. 고정 doorbell+검증 ID만 보내고 본문은 untrusted 도구 결과로 읽힌다. 결과는 retry(확실한 미전달)와 uncertain(재전송 금지)으로 나눈다.
- 이전 보고의 금지 문구가 사용자 지시인지 확인한다. 아니면 정정한다. 실제 Bot 컴퓨터 확인은 coor/owner 후속이다. [실행 기록](../docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md).
