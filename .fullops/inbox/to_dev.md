---
id: handover-dev
title: "SAR-GROK-WAKE-001-DEV — Grok이 매번 receive하지 않아도 도착을 알게 하는 방법(MCP Events, relay webhook)을 비교해 추천안과 순서를 낸다(구현 없음)"
status: approved
updated: 2026-11-28
owner: dev
summary: "Grok이 매번 receive하지 않아도 도착을 알게 하는 방법(MCP Events, relay webhook)을 비교해 추천안과 순서를 낸다(구현 없음)"
key: SAR-GROK-WAKE-001-DEV
role: dev
goal: "Grok이 매번 receive하지 않아도 도착을 알게 하는 방법(MCP Events, relay webhook)을 비교해 추천안과 순서를 낸다(구현 없음)"
base: fddd34c3e2c7aa4a3befd533894114ff199ced81
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-GROK-WAKE-001-DEV — Grok이 매번 receive하지 않아도 도착을 알게 하는 방법(MCP Events, relay webhook)을 비교해 추천안과 순서를 낸다(구현 없음)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 비교 분석과 추천만 한다. 코드는 바꾸지 않는다. 주인 요청을 쿠어가 전달했다(msg_56471dc13fae). 구현은 주인이 추천안을 확인한 뒤 별도 과제로 한다.
- 추가 승인이 필요한 행위: 코드 수정, push, 배포.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)

## 현재 상황과 확인 근거

- 지금은 메시지가 와도 Grok Bot(원격 Remote HTTPS 커넥터, 서버 보관 agent)이 모른다. 주인이 매번 Grok에게 `knowslink_receive`를 시켜야 한다. 주인은 메시지가 오면 Grok이 바로 알게 하고 싶어 한다.
- README에 '선택형 Grok Bot webhook routine wake'(`KNOWSLINK_GROK_WEBHOOK_FILE`)가 이미 있다. 다만 이것은 로컬 watch가 Grok 루틴 URL을 부르는 방식이다. 원격 커넥터에는 로컬 watch가 없다.
- README 한계: 'SSE·MCP log notice만으로 Grok 채팅에 새 턴이 생기지 않는다'.
- 후보는 두 가지다:
  - A. **relay 직접 webhook.** 원격 agent에 도착하면 relay가 주인이 웹에 등록한 Grok 루틴 webhook URL을 부른다. 비밀값은 서버에 암호화해 보관한다. 본문 없이 '새 메시지 N건' 같은 고정 문장만 보낸다.
  - B. **MCP Events 확장.** https://modelcontextprotocol.io/extensions/events/overview
- 주인이 정한 보안 조건(A):
  - URL 등록과 해제는 웹에서 주인만 한다. 5분 재인증을 거친다.
  - 호출 빈도에 상한을 둔다.
  - 본문과 발신자 이름은 넣지 않는다.
  - SSRF를 막는다: https만 허용하고, 사설·loopback·link-local IP와 DNS rebinding을 차단한다.
- PLANS 행 SAR-MCP-SKILLS-001-DEV에 걸린 Events 재검토 항목을 이 분석에서 같이 처리한다(msg_3f6af5fcff02). 대상은 다음과 같다:
  - 상태, `subscriptions/listen`과의 관계
  - 클라이언트 지원: Claude Code v2, Grok Bot, Dots, Orca 봇
  - Grok 채팅에 새 턴이 생기는지
  - 어느 클라이언트에 이득인지

## 적용 기준과 예외

- 1차 출처만 쓴다. 대상은 MCP 공식 문서와 확장 문서, Grok·xAI 공식 문서(루틴, webhook, 커넥터 알림), 각 클라이언트 공식 문서다. 확인하지 못한 것은 '미확인'으로 적는다.

## 먼저 읽을 문서

- `README.md`와 `adapters/README.md`의 Grok webhook routine wake, SSE·log notice 한계 부분
- `internal/relay`의 원격 agent 배달, 원격 `/mcp` events/subscribe(이미 있는 부분), 웹 '내 에이전트'·재인증 패턴(grant 화면)
- `.fullops/tasks/SAR-MCP-UPGRADE-001/analysis.md`

## 해야 할 일과 파일 소유권

- [ ] A와 B를 표로 비교한다. 비교 항목은 다음과 같다:
  - Grok 채팅에 새 턴이 생기는지
  - 지연
  - 클라이언트별 지원
  - 구현 범위(파일·대략 줄 수)
  - 보안 위험과 대책
  - 운영 비용(외부 호출, 재시도)
  - 되돌리기
- [ ] Grok 루틴 webhook을 확인한다. 실제 존재하는지, 호출 형식(인증, 본문), 새 턴을 만드는지를 공식 문서로 본다.
- [ ] 추천안을 1개 고른다. 단계적으로 진행하자면 그 순서를 적는다.
- [ ] 순서를 제안한다. 진행 중이거나 대기 중인 과제와 비교한다:
  - SAR-REDIS-CAP-001-DEV(Redis 용량 상한, 12-28 전 필수)
  - SAR-MCP-SKILLS-001-DEV(Skills, Events 재검토 포함)
  - SAR-MCP-STDIO-SDK2-001-DEV
  - SAR-ORCA-BOT-002-DEV(Grok answer가 번호로 옴)
  - SAR-ORCA-GRANT-UX-001-DEV
- [ ] 결과를 `.fullops/tasks/SAR-GROK-WAKE-001-DEV/analysis.md`(front matter)에 쓴다. 요약 5줄을 `worker_done` 본문에 넣는다.

## 완료 기준과 검증

- `analysis.md`에 비교 표, 근거 링크, 추천, 순서 제안이 있다.
- 문서만 바꿨으므로 `fullops lint --from <base>`를 통과한다.

### UI 작업의 디자인 기준과 검증

해당 없음. 구현할 때 웹 등록 화면은 grant 화면 패턴을 따른다는 점만 적는다.

## 갱신할 산출물

- 없음. 분석 문서만 쓴다.

## 기대 산출물

- `analysis.md`

## 제약·협업·후속

- 하위 에이전트: 공식 문서 조사는 읽기 전용 하위 에이전트로 나눠도 된다.

## 완료 보고

첫 줄은 `[완료] SAR-GROK-WAKE-001-DEV | 브랜치 fullops/dev | SHA <HEAD> | 추천 <A|B|단계>`다. 문서만 바꿨고 push 승인은 받지 않았다. 그래서 `fullops done`을 실행하지 않는다. 로컬 커밋 뒤 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
