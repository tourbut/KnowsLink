---
id: handover-dev
title: SAR-MSG-TTL-1H-001-DEV — 모든 KnowsLink 메시지의 relay 보관 기한을 180초에서 1시간으로 늘리고 발송 기한·도구 설명·보관 안내를 맞춘다
status: approved
updated: 2026-11-28
owner: dev
summary: 모든 KnowsLink 메시지의 relay 보관 기한을 180초에서 1시간으로 늘리고 발송 기한·도구 설명·보관 안내를 맞춘다
key: SAR-MSG-TTL-1H-001-DEV
role: dev
goal: 모든 KnowsLink 메시지의 relay 보관 기한을 180초에서 1시간으로 늘리고 발송 기한·도구 설명·보관 안내를 맞춘다
base: 112e8810ac7a3c025b6a7815e2fd850871b4d889
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-MSG-TTL-1H-001-DEV — 모든 KnowsLink 메시지의 relay 보관 기한을 180초에서 1시간으로 늘리고 발송 기한·도구 설명·보관 안내를 맞춘다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인 결정(쿠어 msg_01c39e2e0d7e, 2026-11-28). KnowsLink 메시지 전부의 relay 보관 기한을 180초에서 1시간으로 늘린다. orca.*만이 아니다. 로컬 발송 기한과 원격 도구 설명도 맞춘다. 저장 용량, 정리 작업, 보관 규칙 문서(개인정보, 보관 기한 안내 포함)에 영향이 있으면 같이 고친다. 되돌릴 수 있는 결정은 추천안으로 잠정 진행한다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, DB 마이그레이션, 사람 발송 승인 규칙의 변경(아래 '구분' 참고). 생기면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음(단일 dev 과제)

## 현재 상황과 확인 근거

- 진단 SAR-ORCA-DIAG-001-OPS(`.fullops/tasks/SAR-ORCA-DIAG-001-OPS/qa/report.md`)의 추정은 다음과 같다:
  - 로컬 발송 기한이 170초다(`adapters/src/text.ts:130`, "under the relay's strict 180 s ceiling").
  - 그래서 Grok이 3분 넘게 지나 receive하면 `messages: []`가 나온다.
  - 원격 도구 설명의 "최대 1시간 보관"(`internal/relay/remote_mcp.go:145` 부근)은 실제와 다르다.
- 랜딩과 정책에는 "최대 180초"와 "최대 1시간 서버 메모리(Redis)" 같은 문구가 있다(`landing_pages.go`, `landing_test.go`, policy pages, i18n). `180`이 나오는 relay 파일은 23개다.
- SAR-MCP-UPGRADE-001(원격 /mcp 2026-07-28)은 main에 들어갔다. 그 리뷰의 low 3건을 이 과제에 합친다(아래 할 일).

## 적용 기준과 예외

- **구분이 먼저다.** 180초가 '메시지 배달·보관 기한'인지 '사람 발송 승인 대기 시간' 같은 다른 뜻인지 코드마다 가린다. 바꾸는 것은 메시지 보관과 배달 기한뿐이다. 승인 대기 시간이나 다른 뜻의 180초는 그대로 두고 목록으로 보고한다. 애매하면 `ask`로 묻는다.
- 본문은 untrusted다. 보관이 길어져도 receipt와 로그에 본문을 남기지 않는다. 보관 위치(Redis 메모리)와 암호화·삭제 방식은 바꾸지 않는다. 바꿔야 하면 묻는다.
- 토큰 비용: 표적 단위 테스트와 기존 integration의 해당 단언만 고친다. 새 통합 테스트 묶음은 만들지 않는다.

## 먼저 읽을 문서

- `.fullops/guide/roles/dev.md`
- 위 진단 보고, `adapters/src/text.ts`, relay의 ttl 검사(`ttl_too_long`)·store·Redis 만료·정리 작업
- 개인정보·보관 기한 정책 페이지(`policy_pages`), 랜딩, `README.md`와 `.fullops/docs/operations/user-guide.md`의 보관 문구
- `.fullops/tasks/SAR-MCP-UPGRADE-001/reviews/6881d23e3362/report.md`(low 3)

## 해야 할 일과 파일 소유권

- [ ] 180초의 뜻을 가린다. 결과는 표 하나로 보고서에 적는다(파일:줄, 뜻, 바꿈 여부).
- [ ] relay의 메시지 보관과 배달 상한을 1시간으로 올린다. Redis 만료와 정리 작업도 같은 값을 쓰게 한다. 한 곳의 상수로 모은다.
- [ ] 로컬 adapter의 발송 기한(text.ts 170초, orca send, MCP send)을 1시간보다 조금 짧게 맞춘다. 시계 오차 여유는 지금 방식대로 둔다.
- [ ] 원격 도구 설명, 랜딩, 정책(ko·en), README, user-guide의 보관 문구를 실제 값과 맞춘다. 정책 문구는 이전 문구와 새 문구를 보고서에 나란히 적는다.
- [ ] 저장 용량을 어림해 보고서에 적는다. 메시지 크기 상한 × 시간당 한도 × 1시간으로 Redis 최대 사용량을 계산하고, 운영 Redis 설정(maxmemory 등)과 견준다. 운영 값을 확인할 수 없으면 '미확인'으로 적는다.
- [ ] SAR-MCP-UPGRADE-001 리뷰 low 3건:
  - modern `events/subscribe`·`unsubscribe` 결과에 `resultType`·`_meta.serverInfo`를 붙이거나, D05 문서를 실제 동작에 맞춘다.
  - `Mcp-Name` 비교는 `params.name`이 있는 메서드에만 한다(`resources/read`는 404 -32601).
  - `remote_mcp.go` 300줄 초과를 분할로 해소한다.

## 완료 기준과 검증

- 표적 테스트를 둔다:
  - 1시간 안의 ttl은 받고, 1시간을 넘으면 `ttl_too_long`이다.
  - 로컬 발송 기한 값을 확인한다.
  - 문구 테스트(landing, policy, i18n)를 고친다.
  - low 3건의 단언을 둔다.
- `fullops lint --from <base>`를 통과한다. 결과는 `tasks/SAR-MSG-TTL-1H-001-DEV/qa/report.md`에 적는다:
  - 180초 구분 표
  - 정책 문구 이전·새 비교
  - 용량 어림
  - 테스트 발췌
- 배포 뒤의 실제 확인은 쿠어가 Grok에 decision을 다시 보내는 것으로 한다. 이 과제 범위 밖이다.

### UI 작업의 디자인 기준과 검증

문구만 바꾼다. 기존 패턴을 그대로 쓰고 새 스타일은 만들지 않는다.

## 갱신할 산출물

- 보관 규칙이 적힌 기획·설계 문서(D03·D05·user-guide·개인정보 정책). 바뀐 값과 근거(주인 결정 2026-11-28)를 적는다.

## 기대 산출물

- 코드, 표적 테스트, 문구, 문서, qa 보고

## 제약·협업·후속

- 범위 밖: MCP Skills(SAR-MCP-SKILLS-001-DEV), SDK 2.x(SAR-MCP-STDIO-SDK2-001-DEV), grant UX
- 하위 에이전트: 180초 위치 조사는 읽기 전용 하위 에이전트로 나눠도 된다. 구현은 직접 한다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
원격 push는 승인 전이다. 그래서 `fullops done`을 실행하지 않는다. 로컬 커밋 뒤 첫 줄 `[완료] SAR-MSG-TTL-1H-001-DEV | 브랜치 fullops/dev | SHA <HEAD>`로 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
