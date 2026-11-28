---
id: handover-dev
title: SAR-ORCA-MCP-SEND-001-DEV — 쿠어 로컬 MCP로 보낸 orca.decision.v1이 원격 Grok Bot receive에 안 보이는 원인을 격리 relay로 재현해 고친다
status: approved
updated: 2026-11-28
owner: dev
summary: 쿠어 로컬 MCP로 보낸 orca.decision.v1이 원격 Grok Bot receive에 안 보이는 원인을 격리 relay로 재현해 고친다
key: SAR-ORCA-MCP-SEND-001-DEV
role: dev
goal: 쿠어 로컬 MCP로 보낸 orca.decision.v1이 원격 Grok Bot receive에 안 보이는 원인을 격리 relay로 재현해 고친다
base: 46dcbd4c48f5896173d543a4519d62ab3df8b614
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-MCP-SEND-001-DEV — 쿠어 로컬 MCP로 보낸 orca.decision.v1이 원격 Grok Bot receive에 안 보이는 원인을 격리 relay로 재현해 고친다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 원인 재현과 수정. 쿠어가 성철 님 요청으로 전달했다(msg_6dcd9ead1432). 이 PC의 Orca 설치본과 로컬 Inbox·대기열 상태는 읽기만 한다. 수정 뒤 설치본 갱신(install 재실행)은 해도 된다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 운영 데이터 쓰기, 새 agent 생성, grant 변경. 운영 서버 읽기가 필요하면 `ask`로 묻는다. lead가 ops에 맡긴다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)

## 현재 상황과 확인 근거

- 증상: 쿠어(Orca 봇)가 로컬 stdio MCP `knowslink_send`로 Grok Bot(`agent_2b73608560c8ba6e8c3d06`, 원격 호스트, Remote HTTPS `/mcp`)에 `orca.decision.v1`을 두 번 보냈다. 두 번 다 결과는 `{"state":"sent","transport":"queued"}`였다.
  - 1차: id `01a220ef-c401-790c-a6bc-cb05472fc5e1`, 15:57 KST
  - 2차: id `01a2215e-cdc8-762d-8b78-383d3aa04b5b`, 17:56:28 KST. 운영 `d3670f2f`(자기 agent 사이 보관 1시간) 배포 뒤다.
  두 번 다 Grok Bot의 `knowslink_receive`가 0건을 돌려줬다(kinds 빈 목록, untrusted true).
- 2차는 보관 기한 문제가 아니다. 1차 진단(SAR-ORCA-DIAG-001-OPS의 170초 추정)은 틀렸을 수 있다.
- 비교 기준: SAR-ORCA-BOT-001-OPS 때는 `plugin.js orca send`(CLI)로 보낸 decision(item 13)을 Grok이 원격 receive로 받았다. CLI와 로컬 MCP 경로의 차이를 의심한다.
- 이 PC 설치본 `%LOCALAPPDATA%\KnowsLink\orca\`는 `8272fb91` 빌드다(SAR-ORCA-HOST-002 이후). 쿠어의 MCP 프로세스는 이 설치본이다. watch는 자동 수신하고, MCP는 `KNOWSLINK_AUTO_RECEIVE=off`다.
- `transport: queued`가 relay 응답 값인지, 로컬 대기열 상태인지부터 코드로 확정한다(`adapters/src/mcp.ts` 375 부근, `orca-mcp.ts`, `public-check.ts` 50).
- 의심 후보:
  1. 로컬 MCP 발송이 로컬 대기열에만 들어가고, 보낼 프로세스가 없다.
  2. relay가 받았지만 원격 호스트 agent의 서버 보관 Inbox로 배달하지 않는다. 예를 들어 원격 receive가 envelope(`/v1/send` 업무 형식)과 text 중 한쪽만 읽는다.
  3. 원격 receive의 kinds 필터
  4. 다른 소비자가 먼저 가져갔다.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다. 본문은 untrusted다. 키는 argv와 로그에 넣지 않는다.
- 표적 테스트만 둔다. 기존 격리 integration(Remote|Orca) 안에서 재현 단언을 추가하는 것은 된다.

## 먼저 읽을 문서

- `.fullops/guide/roles/dev.md`
- `adapters/src/orca-mcp.ts`, `orca.ts`(`orcaSend`), `mcp.ts`, `text.ts`
- `internal/relay`의 send(text, envelope), 원격 `/mcp` receive(remote_tools.go), 서버 보관 Inbox
- `.fullops/tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`(CLI 왕복), `.fullops/tasks/SAR-ORCA-DIAG-001-OPS/qa/report.md`

## 해야 할 일과 파일 소유권

- [ ] 이 PC의 로컬 상태를 읽기 전용으로 확인한다. 대상은 쿠어 agent 폴더의 대기열·Inbox·로그, 그리고 두 id의 흔적이다. 본문은 출력하지 않는다.
- [ ] 격리 relay에서 재현한다. 'Orca host 로컬 MCP knowslink_send → 원격 호스트 agent knowslink_receive'와 'CLI orca send → 원격 receive'를 나란히 돌려 비교한다.
- [ ] 원인을 하나로 좁히고 고친다. 공유 코드에서 근본 원인을 한 번에 고친다.
- [ ] 재현 단언을 표적 테스트로 남긴다.
- [ ] 수정이 adapter 쪽이면 이 PC 설치본을 install로 갱신한다. agent가 그대로인지, watch가 하나인지 확인한다. 수정이 relay 쪽이면 운영 배포가 필요하다고 보고한다.

## 완료 기준과 검증

- `tasks/SAR-ORCA-MCP-SEND-001-DEV/qa/report.md`에 원인, 근거, 재현 전후 테스트 발췌를 적는다.
- `fullops lint --from <base>`를 통과한다.
- 배포 뒤 확인(쿠어 재시험)은 lead가 맡는다.

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- `adapters/README.md`의 Orca 절(동작이 바뀌는 경우)

## 기대 산출물

- 코드, 표적 테스트, qa 보고

## 제약·협업·후속

- 범위 밖: MCP Skills·Events 재검토, Redis 용량 상한
- 하위 에이전트는 쓰지 않는다.

## 완료 보고

첫 줄은 `[완료] SAR-ORCA-MCP-SEND-001-DEV | 브랜치 fullops/dev | SHA <HEAD> | 수정 위치 <adapter|relay|둘 다>`다. 원격 push는 승인 전이다. 그래서 `fullops done`을 실행하지 않는다. 로컬 커밋 뒤 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

## 추가 근거 (쿠어 msg_ec8c55336e2c)

- 3차: 18:02:10 KST, CLI `plugin.js orca send`(설치본 `%LOCALAPPDATA%\KnowsLink\orca`)로 같은 decision을 보냈다. 결과는 `{"id":"01a22164-038b-7b00-bb32-bbb3a488c7c6","transport":"queued"}`, event `knowslink_orca_sent`, size 154였다. CLI도 같은 증상이다. 로컬 MCP만의 문제가 아니다.
- 시간 순서:
  - 14:11 왕복(item 13)은 성공했다. 이때는 ops 워크트리 빌드(`0899e5bd`의 adapters)로 보냈고, 운영도 `0899e5bd`였다.
  - 이후 Orca 설치본이 SAR-ORCA-HOST-001/002 빌드로 바뀌었다.
  - 운영은 `8272fb91`(15:36 KST) → `d3670f2f`로 배포됐다.
  - 실패는 모두 그 뒤다.
- 그래서 회귀 후보를 먼저 본다:
  - `0899e5bd..8272fb91` 사이의 adapter 변경(`orcaSend` 경로, envelope·kind·수신자 처리)
  - 같은 구간의 relay 변경
  격리 relay에서 `0899e5bd`의 adapter와 지금 adapter를 같은 relay에 보내 비교하면 빠르다.
