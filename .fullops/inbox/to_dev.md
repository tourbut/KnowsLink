---
id: handover-dev
title: SAR-EVENTS-SECRET-001-DEV — MCP Events 구독 secret을 암호화 보관하고 Events 광고·배달을 Beta 스펙에 맞춘다(Dots 회귀 없이)
status: approved
updated: 2026-11-28
owner: dev
summary: MCP Events 구독 secret을 암호화 보관하고 Events 광고·배달을 Beta 스펙에 맞춘다(Dots 회귀 없이)
key: SAR-EVENTS-SECRET-001-DEV
role: dev
goal: MCP Events 구독 secret을 암호화 보관하고 Events 광고·배달을 Beta 스펙에 맞춘다(Dots 회귀 없이)
base: ec9dded8e49bec1b3717f6dd928ef1e75154b2b2
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-EVENTS-SECRET-001-DEV — MCP Events 구독 secret을 암호화 보관하고 Events 광고·배달을 Beta 스펙에 맞춘다(Dots 회귀 없이)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인이 승인한 순서 ③이다(msg_2e4feb01e049, msg_17037c7513a0). 되돌릴 수 있는 결정은 추천안으로 정한다. 검사를 통과한 병합은 lead가 바로 push한다. 그래서 `fullops done`을 실행해도 된다.
- 추가 승인이 필요한 행위: 운영 배포, 운영 `.env`에 새 키 넣기(ops 단계), DB 마이그레이션, Dots가 깨지는 비호환 변경. 해당하면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음

## 현재 상황과 확인 근거

- 근거 문서는 `.fullops/tasks/SAR-GROK-WAKE-001-DEV/analysis.md`다. 'MCP Events (B)' 절과 '기존 갭' 절을 본다.
  - MCP Events 구독의 `whsec_` secret이 상태(jsonb)에 평문으로 있다(`remote_events.go` `remoteSubscription.Secret`).
  - relay는 `server/discover`에서 `capabilities.events`로 광고한다(`remote_mcp.go:99`).
  - Beta 스펙(`io.modelcontextprotocol/events`, SEP-3415)은 `capabilities.extensions["io.modelcontextprotocol/events"]` 광고, 요청별 `_meta` 선언, `events/poll`·`events/stream`을 요구한다.
  - 지금 형식은 OpenAI 문서를 따른 것이고, Dots가 이 형식으로 동작한다.
- 도착 이벤트 `knowslink.text.arrived`는 이미 있다. 본문은 `{"arrival":true}`뿐이다.
- 다음 과제 A(SAR-GROK-WAKE-002-DEV)도 서버 암호화 키(AES-256-GCM)를 쓴다. 이번에 키 도입과 암호화 helper를 공용으로 만들어 A가 재사용하게 한다.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다.
  - 키는 env로만 받는다. 로그와 상태에 남기지 않는다.
  - 키가 없으면 secret을 평문으로 저장하지 않는다. 이때는 Events 구독을 거부하거나 기능을 끈다. 추천안을 정하고 근거를 적는다.
- 기존 평문 secret 처리: 배포 때 읽기는 평문 fallback, 쓰기는 암호화로 하는 이행, 또는 재구독 유도 중 추천안을 정한다. DB 마이그레이션이 필요하면 묻는다.
- Dots 회귀: Dots가 쓰는 지금 형식(`capabilities.events` 광고, webhook 구독)은 그대로 둔다. Beta 형식은 추가로 지원한다(dual). 1차 출처는 Events 공식 문서와 OpenAI mcp-events 문서다.
- 표적 테스트만 둔다.

## 먼저 읽을 문서

- 위 analysis.md
- `internal/relay/remote_events.go`, `remote_mcp.go`, `remote_era.go`
- https://modelcontextprotocol.io/extensions/events/overview
- SEP-3415
- https://developers.openai.com/plugins/build/mcp-events

## 해야 할 일과 파일 소유권

- [ ] 공용 서버 암호화를 만든다. env 키(32바이트 base64), AES-256-GCM, AAD(구독 ID 또는 agent ID)를 쓴다. A가 재사용할 수 있게 둔다.
- [ ] `whsec` secret을 암호화해 저장한다. 기존 평문은 이행 방식으로 처리한다. 화면, 로그, 상태 응답에 secret이 다시 나오지 않게 한다.
- [ ] Beta 스펙 광고와 `_meta` 선언을 추가한다(dual). 지금 형식은 유지한다.
- [ ] `events/poll`(상태 없음)을 추가한다. `events/stream`은 범위와 비용을 판단해 추천안을 정한다. 하지 않으면 이유를 적는다.
- [ ] 이벤트 본문에 본문과 발신자가 없는 것을 테스트로 고정한다.
- [ ] 문서를 고친다. 대상은 D05(인터페이스), D07(보안)이고, D12에는 새 env 키와 이행 절차를 적는다.

## 완료 기준과 검증

- 표적 테스트로 다음을 확인한다:
  - 암호화 왕복
  - 키가 없을 때의 동작
  - 평문 이행
  - dual 광고
  - poll 결과
  - 지금 형식의 구독이 그대로 동작함(Dots 회귀 단언)
- `fullops lint --from <base>`를 통과한다.
- 결과는 `tasks/SAR-EVENTS-SECRET-001-DEV/qa/report.md`에 적는다.

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- D05, D07, D12

## 기대 산출물

- 코드, 표적 테스트, 문서, qa 보고

## 제약·협업·후속

- 범위 밖: routine webhook(A)과 등록 화면. A는 다음 과제다.
- 하위 에이전트: 공식 문서 조사는 읽기 전용 하위 에이전트로 나눠도 된다.

## 완료 보고

`fullops done`을 실행하고, 그 출력 한 줄을 worker_done 본문으로 쓴다. 운영 배포 때 ops가 할 일(새 env 키 생성·주입, 이행 확인)은 남은 일에 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
