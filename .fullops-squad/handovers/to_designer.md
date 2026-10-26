---
title: SAR-CONNECT-SIMPLE-001 — 연동 구조 재점검과 UI/UX·연결 처리 단순화(제품 기획)
status: draft
updated: 2026-10-26
owner: coor
tasks: [SAR-CONNECT-SIMPLE-001]
summary: 에이전트 한도 제거·연동 절차 단순화의 제품 기준과 역할별 지시서 작성
---

# SAR-CONNECT-SIMPLE-001 — 연동 구조 재점검과 UI/UX·연결 처리 단순화

- 작성일: 2026-10-26
- From / To: coor → designer(제품 기획)
- 상태: ready
- 승인된 범위: 제품 기획 문서·설계 문서·역할별 지시서(to_dev / to_tester) 작성과 커밋. 코드 구현은 하지 않는다.
- 추가 승인이 필요한 행위: 운영 배포, 운영 데이터 변경, 기존 agent·키 철회, 실제 메시지 발송, 보안·권한 수락 기준 완화.
- 담당 워크트리 / 브랜치: C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀: coor 워크트리 C:/Users/shin/orca/workspaces/KnowsLink/coor (run id·task id·dispatch id는 Orca preamble 값)

## 사용자 요청 원문 (2026-10-26)

1. "일단 연결 갯수 한도를 풀어놔야겟다"
2. "그리고 기본적으로 연동하는게 너무 복잡해 전체적으로 손좀 봐라"
3. "전체적인 구조를 다시한번 살피고 최대한 UI/UX를 편하게 수정하고 연동 처리도 깔끔하게 되도록해"
4. "에이전트 연결한도도 없애고"
5. "테스트 중지" (링크·Claude 실연결 테스트는 중지했다)

확정 규칙(사용자 결정): **회원당 에이전트 연결 한도를 없앤다.** 서버 전체 보호 상한(`total` 200)·철회 기록 보존(`records`)·키 기록 보호 상한을 함께 없앨지, 남긴다면 어떤 남용 방지로 대신할지는 제품 기준으로 정하고 근거를 적는다.

## 현재 상황과 확인 근거 (coor 실측, 운영 relay link.knowslog.com)

- 한도: `internal/relay/connections.go` `agentLimit` — 회원당 활성 5개(`active`), 전체 200(`total`), 회원별 기록 상한(`records`, `ownerAgentRecords`). 문구는 `internal/relay/member_agents.go` `agentLimits`.
- 테스트 중 회원 화면에 이름 없는 에이전트가 3~4개 쌓였다(연결 실패·중단 시도). 이것 때문에 Claude Code 세 번째 연결이 "연결 에이전트 초과"로 거절됐다.
- 운영에 Claude Code·Codex 자동 설정이 아직 게시되지 않았다. 회원 화면은 "연결 문서를 따르라"고만 안내하고, 실제로는 레포 빌드 + `scripts/package_plugin.py --host claude --install-dir … --connect` + `claude --plugin-dir …`가 필요했다.
- 연결 실패 사유가 사용자에게 보이지 않는다. `knowslink_connect`가 다음 경우 모두 `{"state":"failed"}`만 돌려줬다.
  - 이전 시도의 `login.json`·`private.pem`이 남은 agent 폴더를 재사용한 경우
  - 빈 agent 폴더를 미리 만들어 둔 경우(`private-files.ts` `privateDirectory`가 기존 폴더를 거부)
  - Windows의 비시스템 드라이브(D:)에 private 폴더를 둔 경우(`windowsACL`의 `SetOwner`가 소유자 본인이어도 UnauthorizedAccess)
- Google 동의 단계에서 다른 Google 계정(authuser=1)을 고르면 403 "이 계정으로 로그인하거나 재확인할 수 없습니다…"라는 일반 문구가 나온다. 원인(새 가입 한도/다른 소유자)을 구별할 수 없다.
- 동의 완료 뒤 연결 링크를 다시 열면 "만료·취소·사용되었거나 유효하지 않습니다" 오류가 보여, 성공했는지 실패했는지 헷갈린다.
- 동의 화면에서 이름을 정하지 않으면 '이름 없는 에이전트'가 생기고, Grok 쪽 도구 결과에 표시 이름이 나오지 않는다.
- Grok Bot은 계정 단위 플러그인 하나만 env 편집 없이 쓸 수 있다. 두 번째 Bot은 별도 custom MCP를 손으로 등록해야 했다.
- `package_plugin.py --verify`의 `adapters/dist/mcp.test.js` 실패 경계 검사가 main b413464에서 실패한다(`/v1/registry` 뒤 `/v1/pull` 미호출). 원인은 아직 모른다. DEV 지시서에 포함한다.
- 열린 결정 D040(tester 실세션 QA)과 AUTOSETUP 운영 게시 승인은 그대로 열려 있다.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`와 연결된 코딩·테스트·보안 규칙. 기준 커밋 main `b413464adb86fba56488cb715d9bae58f247e96e`.
- 제품 정본: `.fullops-squad/project.md`, `docs/planning/product-specs/`, `docs/design-docs/` (기존 SAR-AUTOSETUP-001, SAR-UX-ISSUES-002 문서 포함).
- 산출물 ID `D01~D13` → `docs/deliverables/README.md`. 결정 `D040` 등은 fullops-decisions.
- 보안 경계(키 private 폴더, 지문 확인, 발송 승인, 친구 권한)는 유지한다. 단순화는 이 경계를 없애는 방식이 아니어야 한다.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/jev/SAR-CONNECT-SIMPLE-001-route.json`
- `README.md`의 자동 설정 절, `adapters/README.md`
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`, `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md`
- `.fullops-squad/docs/operations/user-guide.md`
- `internal/relay/connections.go`, `internal/relay/member_agents.go`, `adapters/src/login.ts`, `adapters/src/private-files.ts`, `adapters/src/setup-connect.ts` (사실 확인용, 수정은 DEV)

## 해야 할 일과 파일 소유권

- [ ] 현재 연동 흐름(Grok Bot·Dots·Claude Code·Codex, 처음 연결·다시 연결·실패 복구)을 사용자 행동 단위로 정리하고 복잡한 지점을 목록화한다.
- [ ] 목표 흐름을 정한다: 사용자 행동 수, 실패 시 보이는 사유와 다음 행동, 재시도 시 폴더·키 처리, 이름 정하기, 완료 화면.
- [ ] 회원당 에이전트 한도 제거의 제품 규칙과, 남길 보호 장치를 정한다.
- [ ] 회원 화면(내 에이전트·연결 안내·동의 화면·완료/오류 화면)의 UI/UX 개선 기준을 정한다. 미사용·이름 없는 에이전트 정리 방식을 포함한다.
- [ ] 제품 문서를 갱신한다: D10, D02, D04 (+ 필요 시 D03, D05, D06, D09, D12, D13).
- [ ] `handovers/to_dev.md`에 구현 지시서를 쓴다(한도 제거, 연결 실패 사유 노출, 폴더 재사용·선생성·비시스템 드라이브 처리, mcp.test 실패 원인 확인, 화면 수정). 구현 파일·함수 계획은 DEV가 정한다.
- [ ] 동작이 바뀌므로 `handovers/to_tester.md`에 검증 지시서를 쓴다(DEV 완료 SHA 대기, 다시 돌릴/새로 만들 시나리오 경로 `docs/evaluations/scenarios/`).
- [ ] UI 변경이 있으면 ART 시각 검수 담당과 인계 조건을 지시서에 적는다.

## 완료 기준과 검증

- 제품 문서와 역할별 지시서가 커밋돼 있고, 지시서마다 완료 조건·검증 명령·파일 소유권이 정해져 있다.
- lint: `python3 <fullops scripts>/lint.py --repo <designer 워크트리> --from b413464adb86fba56488cb715d9bae58f247e96e` 통과.
- worker_done 본문: `[설계] SAR-CONNECT-SIMPLE-001 | 지시서: <경로들> | 역할: … | SHA …`

## 갱신할 산출물

D10, D02, D04 (route 추천). 추가 확인: D03, D05, D06, D09, D12, D13.

## 제약·협업·후속

- 범위 밖: 운영 배포·AUTOSETUP 게시(사용자 승인 대기), D040 실세션 QA.
- 사용자 판단이 필요한 제품 결정만 `ask`로 묻는다. 기술 방법은 DEV에 맡긴다.
- 테스트 중 생성된 운영 agent(링크 `agent_c18976751d5a96d8adfbb0` 등)와 이름 없는 agent는 건드리지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물 링크 / 남은 일. worker_done 본문 첫 줄은 `[완료] <과제 키> | 브랜치 <이름> | SHA <브랜치 HEAD>` 형식을 따른다(설계 보고는 위 `[설계]` 줄).
