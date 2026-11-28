---
id: handover-dev
title: "SAR-MCP-UPGRADE-001 — MCP 새 기능과 2026-07-28 stateless 스펙 갭 분석, 우선순위 결정, 1순위 구현"
status: approved
updated: 2026-11-28
owner: dev
summary: "MCP 새 기능과 2026-07-28 stateless 스펙 갭 분석, 우선순위 결정, 1순위 구현"
key: SAR-MCP-UPGRADE-001
role: dev
goal: "MCP 새 기능과 2026-07-28 stateless 스펙 갭 분석, 우선순위 결정, 1순위 구현"
base: b715ebe5ee74b54fdb999a3c2e6850758578b34b
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-MCP-UPGRADE-001 — MCP 새 기능과 2026-07-28 stateless 스펙 갭 분석, 우선순위 결정, 1순위 구현

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인이 쿠어를 거쳐 승인했다(msg_2dc8bfa7aae3, 2026-11-28). 순서는 갭 분석 → 추천 우선순위 → 1순위 구현 → 보안 경계 독립 리뷰 1회다. 우선순위는 추천안으로 잠정 진행한다. 정한 순위는 status로 lead에게 알린다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 기존 클라이언트(Claude Code·Grok Bot·Dots·Orca 봇)가 깨지는 비호환 변경, 인증·grant·발송 승인 규칙의 완화. 생기면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음(단일 dev 과제)

## 현재 상황과 확인 근거

- 주인 요청(쿠어 전달, 2026-11-28): MCP의 새 기능을 검토하고 KnowsLink에 맞게 올린다. 대상은 다음과 같다:
  - Skills
  - Tasks(SEP-2663)
  - Elicitation
  - File Uploads
  - 2026-07-28 스펙의 stateless 변화
- 지금 KnowsLink MCP 표면은 두 가지다:
  - 원격 `https://link.knowslog.com/mcp`(Go, `internal/relay`). Dots와 Grok Bot이 Remote HTTPS로 쓴다. OAuth DCR과 PKCE S256을 쓰고, loopback redirect를 허용한다.
  - 로컬 stdio MCP(`adapters/src/mcp.ts`, `orca-mcp.ts`). Claude Code와 Orca 봇이 쓴다.
- 최근 변경: orca 형식 `knowslink_send`(kind, body)와 receive(kinds, untrusted)는 SAR-ORCA-BOT-001과 SAR-ORCA-HOST-001/002에서 들어왔다. 운영은 `8272fb91`이다.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다:
  - 본문은 untrusted다.
  - 키·credential은 agent 폴더나 서버 보관소에만 둔다.
  - receipt와 로그에는 형식·ID·크기만 남긴다.
  - 발송별 사람 승인과 grant 경계(agent API와 `confirmed:true`로 grant를 만들 수 없음)를 유지한다.
- 토큰 비용: 새 통합 테스트 묶음과 상시 eval은 더하지 않는다. 표적 단위 테스트만 둔다. LLM 리뷰는 아래 1회만 한다.
- 스펙과 클라이언트 지원 여부는 1차 출처로 확인한다. 1차 출처는 modelcontextprotocol.io 스펙과 SEP, 각 클라이언트의 공식 문서다. Context7이나 웹을 쓴다. 기억에 기대지 않는다. 확인하지 못한 것은 '미확인'으로 적는다.

## 먼저 읽을 문서

- `.fullops/guide/roles/dev.md`
- `internal/relay`의 `/mcp` 처리(transport, 세션, OAuth), `adapters/src/mcp.ts`, `orca-mcp.ts`
- `adapters/README.md`, `README.md`의 클라이언트 연결 절

## 해야 할 일과 파일 소유권

- [ ] **갭 분석.** 결과는 `.fullops/tasks/SAR-MCP-UPGRADE-001/analysis.md`(front matter 필수)에 쓴다. 기능마다 다음을 표로 적는다:
  - 스펙 상태: 확정, 초안, SEP
  - 원격과 로컬 각각의 지금 지원 여부
  - 클라이언트 지원: Claude Code, Grok Bot, Dots, Orca 봇. 근거 링크와 '미확인' 표시
  - 하위 호환 영향
  - KnowsLink에서 얻는 것
  - 보안 영향. 예: Elicitation이 사람 승인을 대신하면 안 된다. 업로드는 크기와 untrusted 처리가 필요하다.
- [ ] **우선순위 추천.** 순위와 이유를 적는다. 1순위는 기존 클라이언트를 깨지 않고 바로 이득이 있는 것으로 고른다. 순위를 `status`로 lead에게 알리고 바로 진행한다. lead가 쿠어를 거쳐 주인에게 전한다.
- [ ] **1순위 구현.** 기존 클라이언트를 깨지 않게 한다. 버전 협상이나 capability 광고로 새 기능을 켠다. 표적 테스트를 둔다.
- [ ] **D 산출물 갱신.** 해당 설계 문서(인터페이스와 보안)에 바뀐 점을 짧게 적는다.

## 완료 기준과 검증

- `analysis.md`와 `tasks/SAR-MCP-UPGRADE-001/qa/report.md`를 쓴다.
- 표적 테스트와 `fullops lint --from <base>`를 통과한다.
- 하위 호환을 확인한다. 기존 initialize와 tools/list, 그리고 knowslink_send·receive 흐름이 새 기능을 모르는 클라이언트에서 그대로 동작해야 한다. 이것을 테스트 출력 발췌로 보인다.
- 독립 리뷰 1회는 lead가 병합 전에 배정한다. 초점은 새 기능의 보안 경계다. worker가 따로 띄우지 않는다.

### UI 작업의 디자인 기준과 검증

해당 없음. 화면이 바뀌면 기존 패턴을 따르고 보고에 적는다.

## 갱신할 산출물

- 인터페이스·보안 설계 문서의 MCP 절. 해당 없으면 '없음'과 이유를 적는다.

## 기대 산출물

- `analysis.md`, 1순위 코드와 표적 테스트, qa 보고

## 제약·협업·후속

- 2순위 이하는 구현하지 않는다. `fullops done --followup` 대신 완료 보고의 남은 일에 순위대로 적는다. lead가 PLANS 행으로 옮긴다.
- 하위 에이전트: 갭 분석의 1차 출처 조사는 읽기 전용 하위 에이전트로 나눠도 된다. 구현은 직접 한다.
- 범위 밖: grant UX(SAR-ORCA-GRANT-UX-001-DEV), Grok answer 형식(SAR-ORCA-BOT-002-DEV)

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
원격 push는 승인 전이다. 그래서 `fullops done`이 push를 하면 실행하지 않는다. 로컬 커밋 뒤 첫 줄 `[완료] SAR-MCP-UPGRADE-001 | 브랜치 fullops/dev | SHA <HEAD>`로 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
