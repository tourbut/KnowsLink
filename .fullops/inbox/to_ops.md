---
id: handover-ops
title: SAR-ORCA-DIAG-001-OPS — 쿠어→Grok decision 01a220ef…이 Grok 받은편지함에 안 보이는 원인을 운영 기록으로 찾는다(읽기 전용)
status: approved
updated: 2026-11-28
owner: ops
summary: 쿠어→Grok decision 01a220ef…이 Grok 받은편지함에 안 보이는 원인을 운영 기록으로 찾는다(읽기 전용)
key: SAR-ORCA-DIAG-001-OPS
role: ops
goal: 쿠어→Grok decision 01a220ef…이 Grok 받은편지함에 안 보이는 원인을 운영 기록으로 찾는다(읽기 전용)
base: 4d3998e52dcb7e31045b83c3e524cac2a6c21dab
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-DIAG-001-OPS — 쿠어→Grok decision 01a220ef…이 Grok 받은편지함에 안 보이는 원인을 운영 기록으로 찾는다(읽기 전용)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 운영 읽기 전용 진단. 대상은 relay 로그, Redis와 Postgres 읽기, 이 PC 설치본 로그 읽기다. ops 브랜치 push(`fullops done`)는 지난 과제와 같은 기준으로 승인된 것으로 본다.
- 추가 승인이 필요한 행위: 운영 데이터 쓰기와 삭제, 배포, 재시작, `.env`·Tunnel 변경, 메시지 재전송. 고칠 것이 보이면 고치지 말고 원인과 추천을 보고한다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)

## 현재 상황과 확인 근거

- 쿠어 보고(msg_fd3bfd0579b5):
  - 2026-11-28 15:57 KST쯤 쿠어가 로컬 MCP `knowslink_send`로 Grok Bot(`agent_2b73608560c8ba6e8c3d06`)에 `orca.decision.v1`을 보냈다(item 18).
  - 결과는 `{"state":"sent","id":"01a220ef-c401-790c-a6bc-cb05472fc5e1","transport":"queued"}`였다.
  - 그런데 Grok Bot의 `knowslink_receive`가 `messages: []`를 돌려줬다.
- `transport: queued`는 relay의 발송 응답 값이다(`adapters/src/mcp.ts` 375, `public-check.ts` 50). relay가 받은 것으로 보이지만 확인이 필요하다.
- 이전 운영 왕복(item 13)에서는 Grok이 원격 `/mcp` receive로 decision을 받았다.
- 의심 후보는 다음과 같다. 추정일 뿐이며 직접 확인한다.
  1. 보존 1시간(Redis)이 지나 만료됐다.
  2. Grok의 receive가 `kinds` 없이 불려 orca 형식이 걸러졌다.
  3. 원격 호스트 agent의 서버 보관 Inbox로 배달되지 않았다.
  4. 이미 다른 receive가 소비하고 ACK했다.
  5. 로컬 MCP 발송이 다른 receiver 경로나 grant 판정을 탔다.

## 적용 기준과 예외

- D12 운영 절차를 따른다. 키·credential·본문을 로그와 보고에 넣지 않는다. 형식·ID·크기·시각만 적는다.
- SSH 스크립트는 파일로 올려 실행한다. 서버 명령은 Bash 도구로 실행한다.

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`
- `.fullops/tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`의 4차(item 13 왕복)와 `state-check.sh`
- `internal/relay`의 send·receive·원격 `/mcp` receive(kinds 필터)·보존 만료 코드

## 해야 할 일과 파일 소유권

- [ ] relay 로그와 저장소에서 id `01a220ef-c401-790c-a6bc-cb05472fc5e1`을 추적한다. 다음을 확인한다:
  - 접수 시각
  - receiver
  - kind
  - grant 사용 수 변화
  - 배달 상태: 대기, ACK, 만료
- [ ] Grok의 receive 호출 기록을 확인한다(시각, 결과 건수, kinds 인자 유무 등 남는 범위만).
- [ ] 원인을 하나로 좁힌다. 재현하려면 쓰기가 필요하면 하지 말고 추천만 적는다.
- [ ] 원인별 추천을 적는다:
  - 사용 안내만 바꾸면 되는지(예: kinds 지정, 보존 시간 안에 받기)
  - 코드 수정이 필요한지. 필요하면 어느 파일, 어느 과제인지

## 완료 기준과 검증

- `tasks/SAR-ORCA-DIAG-001-OPS/qa/report.md`에 다음을 적는다:
  - 추적 결과 발췌(본문 없이)
  - 원인
  - 근거
  - 추천

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- 없음(진단). 원인이 운영 절차면 D12에 한 줄.

## 기대 산출물

- `tasks/SAR-ORCA-DIAG-001-OPS/qa/report.md`

## 제약·협업·후속

- 코드 수정은 범위 밖이다. dev는 SAR-MCP-UPGRADE-001(원격 `/mcp` 수정)을 진행 중이다. `internal/relay`의 MCP 코드와 겹치는 수정 추천이면 그 사실을 적는다.

## 완료 보고

브랜치 / SHA / 원인 / 근거 / 추천 / 못 한 것.
`fullops done`을 실행하고 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
