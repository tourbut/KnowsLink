---
title: 유효 자격 정리 flood H-2 수정 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX-2]
summary: OPS FIX-REVIEW H-2(유효 자격의 타 owner·lease 없는·반복 정리가 DB 입장 동안 정리 채널 점유)와 L-2의 원인·수정·회귀 검증을 기록한다
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX-2 — 유효 자격 정리 flood H-2 수정 기록

## 기준과 기술 계획

- 기준 ref: fixed `dfc70caa748a90614b02d48c78b4651345938339`. 지시서 준비 커밋 `283695e`에서 FULLOPS·project·fullops-common-0.3.3·document-writing·contexts/dev, SAR-PUBLIC-SERVICE PS-11과 운영 기본값 표(동시 처리·정리 budget), OPS FIX-REVIEW report·review-foreign-cleanup_test.go.txt, 원 DEV-FIX 실행 기록을 읽었다.
- 원인 재확인: `boundedHTTP`는 `cleanupCredential`과 유효 자격 색인만으로 로컬 정리 채널을 골랐다. 자기 기록(`cleanupTarget`)과 rate는 DB 입장 transaction 안에서 판정했다. 정리 슬롯은 pool 획득·전역 row lock 대기 동안 잡혀 있었다. rate 거부도 transaction 안이므로 슬롯을 보호하지 못했다. 같은 owner의 유효한 자기 정리 반복도 4개를 모두 잡을 수 있었다.
- 계획: (1) DB 전 로컬 선택에 transaction과 같은 자기 기록 대조를 커밋 snapshot으로 적용한다. (2) 정리 공정성 단위를 owner로 하고 로컬 채널과 공유 Clean 기록에서 owner당 동시 1개로 제한한다. (3) snapshot은 commit 순서로만 교체한다. (4) 종료 실패 기록이 그 owner를 30s 동안 막지 않게 한다. 제품 수치(신규16·정리4·rate·10s·30s·8/32KiB)·wire·응답 코드는 바꾸지 않는다.
- 공통 경계: 모든 HTTP 요청은 `boundedHTTP` 한 곳을 지난다. `requestBuckets`·`enterHTTP`·`cleanupTarget`의 호출자는 이 경계와 시험뿐이다. 수정은 이 경계·분류·commit 함수(`transaction`)에만 둔다.
- 라이브러리 근거: 새 의존성은 없다. 표준 `sync.Map.LoadOrStore`·`sync/atomic.Pointer.CompareAndSwap`을 go1.27.1 설치본의 `go doc`으로 확인했다. Context7 조회 대상이 아니다.

## 원인과 수정

### H-2 유효 자격의 타 owner·lease 없는·반복 정리가 정리 채널을 점유

- 로컬 분류(`cleanup_admission.go` `cleanupOwner`): `transaction`은 commit 뒤 상태를 `committed{st, at, creds}` snapshot으로 둔다. creds는 자격 해시→`{kind, principal, 세션 만료}` 색인이다. DB 전에 색인으로 principal을 O(1)로 찾는다. 경로별 자격 종류를 `requestBuckets`와 같게 맞춘다(세션 경로는 세션만, /v1은 agent bearer 또는 owner, /owner는 owner만). 그 뒤 같은 `cleanupTarget`을 snapshot에 읽기 전용으로 적용한다. 통과한 요청만 정리 채널을 쓴다. 타 owner·lease 없는·위조·만료 정리는 DB 대기 동안 신규 채널만 쓴다.
- owner 공정성(`capacity.go`): 로컬은 `Service.cleaning`(sync.Map)으로 owner당 정리 슬롯 1개다. 이미 잡혀 있으면 DB 호출·rate 없이 기존 `429 capacity`(Retry-After 1)다. 공유는 `admission.Owner`를 저장하고 `enterHTTP`가 같은 owner의 유효 Clean 기록이 있으면 거부한다. 다중 프로세스에서도 owner 하나가 공유 정리 4개를 모두 잡지 못한다.
- 공정성 단위: agent 자격은 agent의 owner, 세션은 회원의 owner, owner 자격은 그 owner다. 한 회원이 agent 5개·세션 여러 개를 써도 정리 슬롯은 1개다. rate key(principal별)는 바꾸지 않았다.
- 상한·비용: 정리 슬롯을 모두 막으려면 정리 대상 자기 기록을 가진 서로 다른 owner 4명이 동시에 필요하다. 이 다중 계정 경우는 기존 신원 한도(이메일 확인·IP rate)가 비용을 만든다. 제품 quota 변경이 아니므로 designer 질문은 하지 않았다. 같은 owner의 두 번째 동시 정리는 1초 뒤 재시도 안내를 받는다. adapter ACK는 delivery마다 순차라 영향이 없다. snapshot 저장은 이전 색인 생성과 같은 비용이다(L-1 보존). DB 전 분류는 해시 조회와 map 조회만 한다. 새 무제한 대기·DB 기록은 없다. `cleaning` 항목 수는 진행 중 정리 요청 수 이하다.

### L-2 색인 stale·restart·Store 순서

- Store 순서 역전: `remember`가 DB commit 시각(`row.Now`, 전역 row lock 순서로 비감소)을 비교해 더 이른 snapshot이 늦게 저장되어도 새 snapshot을 덮지 않는다(CAS 반복).
- 다른 프로세스의 새 자격·재시작 직후 nil: snapshot이 모르면 신규 채널로 간다. 위조 token과 같은 경로라 M-1 보호를 유지한다. 공유 Clean·rate는 transaction이 정확히 판정한다. 다음 로컬 commit이 snapshot을 갱신한다. `cmd/relay`의 1s retention sweep이 commit이므로 DB가 정상이면 1s 안에 갱신된다.
- 막 철회된 자격: snapshot이 아직 유효로 보면 owner당 1개 슬롯을 한 transaction 동안 쓴다. transaction이 신규로 판정하면 기존처럼 신규 슬롯으로 돌려준다.
- 종료 실패 기록: 경합에서 종료 transaction(2s)이 실패하면 공유 기록이 30s 남아 그 owner의 정리를 막았다(첫 flood 실행에서 13/18 capacity로 발견). 실패 token을 `Service.orphans`에 두고 같은 프로세스의 다음 commit이 지운다. crash 회수 30s는 그대로다.

### coor 추가 인계: TESTER 원본 QA의 HTTP 입장 기록 잔류(medium)

- 출처: coor handoff(사용자의 Grok tester 중단 뒤 원본 증거 추가). 원본 QA 기록 `60b9892`, 후보 `09c523d`, `TestQAMessagesIndependent` 17차 exit1. 원본 로그·`probe_test.go.src`는 fullops-tester 체크아웃에서 읽기만 했다. 원본 실패는 바꾸지 않는다.
- 현상: receipt 20000 거절(409) 뒤 신규 `st.HTTP` 기록 1개가 남았다. 남은 기록은 30s 만료까지 신규 16 중 1개를 차지한다.
- 원인: 20000개 메시지 상태는 transaction마다 JSON 왕복 비용이 크다. 종료 transaction의 2s 기한 안에 row lock·처리를 끝내지 못하면 오류를 버리고 기록을 30s 만료에 맡겼다. H-2 첫 flood 실행의 13/18 실패와 같은 경로다.
- 수정: 위 L-2의 `orphans`·`reclaim`이다. 종료 실패 token은 같은 프로세스의 다음 commit(입장 transaction 10s 기한, 1s sweep)이 지운다. 신규·정리 기록 모두 해당한다. 상태 크기 자체의 비용은 L-1(전역 lock·전체 JSON)로 보존한다.
- 회귀: `TestFailedFinishReclaimedAtNextCommit`이 handler 안에서 실제 row lock을 잡아 종료 transaction을 2s 기한으로 실패시킨다. 정리(`/v1/key-revoke`)와 신규(`/v1/pull`) 모두 다음 commit에서 공유 기록 0, 같은 owner 재입장 가능, orphan 0을 확인한다. 새 fixed 독립 QA는 TESTER 후속이며 이 DEV 검사로 대체하지 않는다.

### /v1/connect·정리 caller 대조

- `/v1/connect/`에는 `info`·`prepare`·`complete`만 있다(`member_agents.go` `agentRoutes`). `/v1/connect/{cancel,key-revoke,agent-revoke,unpair}` 경로는 없다. cancel·key-revoke·agent-revoke·unpair는 `/home/*` `memberAction`이며 이미 정리 분류·cleanupRate·reauth·현재 owner 검사를 쓴다. 변경하지 않았다.
- cleanupRate를 쓰는 모든 handler를 대조했다: `memberAction`(cancel·key-revoke·agent-revoke·unpair·invite deny), `gateOwner`(gate deny), `logout`, `/v1` `operate`(ack·key-revoke·unpair·owner-revoke), `/v1/text/ack`, `/v1/test/ack`. 누락은 `/v1/invite-decision`의 JSON `decision:deny` 하나였다. `/home/invite-decision` deny와 같게 정리 분류와 `cleanupTarget`(수신 owner)에 추가했다. handler 권한·Generation 검사(L-B)는 바꾸지 않았다.
- 정책·CSRF·closed schema·현재 키·관계 세대·lease·클라이언트 자격 검사는 handler에 그대로 있다. 분류는 budget만 고른다.

## 자동 검증 증거

증거: [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2/). `.exit`는 명령 반환값이다. 대조군은 자기 격리 Compose Postgres(같은 migration)에서 scratch worktree로 실행하고 자원을 회수했다.

| 검사 | 대상 | 결과 |
|---|---|---|
| OPS H-2 재현 시험 원문 | fixed dfc70ca | exit0 = 재현. A 4/4·own 429 capacity, B new 0/18·foreign 18/18, C no-lease 18/18 (control-dfc-review-repro) |
| 같은 재현 시험 | 수정 코드 `0ae40d5` | exit1 `not reproduced: held 0`. 이 시험이 row lock을 10s 넘게 잡아 probe가 503이 된다. 정리 채널 점유 0/4 (fix-review-repro) |
| 새 flood 회귀의 변형 | `0ae40d5` + 로컬 자기 기록 대조·owner 상한 제거 | exit1 `attack admission: new 0 clean 4` (control-h2-mutant) |

새·변경 검사:
- `TestValidCredentialCleanupFlood`(통합): A. 전역 row lock을 잡은 채 타 owner 키 철회 4·agent lease 없는 ACK 4는 신규 채널 8개로 대기한다. 같은 owner의 자기 키 철회 4는 정리 슬롯 1개만 잡고 3개는 즉시 429 capacity다. 다른 owner 3명의 자기 철회가 나머지 정리 슬롯 3개를 잡는다. lock을 풀면 모두 200이다. rate는 타 owner 신규 4·no-lease 신규 4·반복 owner cleanup 1·다른 owner 각 1이다. B. lock 없이 48 연결 flood(타 owner·no-lease·자기 반복·신규 대조) 중 다른 owner의 유효 정리 18회가 모두 200이다. 종료 뒤 공유 기록·로컬 채널·owner 표시가 0이다.
- `TestCleanupSnapshotAcrossProcessesAndRestart`(통합): 다른 Service에서 만든 owner 자격은 이 프로세스의 다음 commit 전까지 정리 채널을 쓰지 않는다. 재시작한 빈 Service도 같다. 한 번 commit하면 owner를 정확히 찾는다.
- `TestFailedFinishReclaimedAtNextCommit`(통합): 실제 row lock으로 종료 transaction을 실패시킨 정리·신규 기록이 다음 commit에서 지워지고 같은 owner가 다시 입장한다.
- `TestRememberKeepsNewestCommit`(unit): 늦게 저장되는 이전 commit이 새 snapshot을 덮지 않는다.
- `TestCleanupAdmissionNeedsVerifiedOwnRecord`(unit): 12개 행렬마다 로컬 `cleanupOwner`가 transaction 판정과 같다. 세션 쿠키의 /v1 사용·agent 자격의 /owner 사용·idle 만료 세션은 정리 채널을 쓰지 않는다.
- `TestSharedCapacitiesAndCleanupClassification`(unit): owner당 공유 Clean 1, `/v1/invite-decision` deny/accept 분류, 재시작 직렬화 뒤 Clean 상한.
- 기존 `TestHTTPSlowBodyAndCleanupAdmission`·`TestHTTPConcurrencyAcrossInstances`: 자기 정리 4건을 서로 다른 owner 4명으로 바꿨다. 같은 owner 4건 동시 정리는 이제 H-2 반복 flood와 같은 형태라 1건만 입장한다. 5번째 owner는 로컬·공유 정리 상한에서 429이며 같은 owner 반복도 429이고 rate를 쓰지 않는다. 16/4 두 인스턴스 경합·종료 회수 의미는 같다.

원본 실패 보존: 첫 개발 실행(자기 격리 DB)에서 flood B 타 owner 단계가 13/18 capacity로 실패했다. 원인은 경합 중 종료 transaction 실패로 남은 owner의 공유 Clean 기록이었다. reclaim으로 고친 뒤 3회 반복과 전체 통합 실행이 통과했다. flood 중 유효 정리 지연(최대 약 3s)은 전역 row lock 대기열(L-1)이며 신규 대조 flood에서도 같다. 10s 기한 안에서 모두 200이다.

## 최종 검사

| 검사 | 대상 HEAD | 결과 |
|---|---|---|
| make lint | 코드 `0ae40d5` | exit0 (lint-final) |
| make test | 코드 `0ae40d5` | exit0, Go race·adapter 시험 (test-final) |
| make verify-mvp | 코드 `0ae40d5`(코드 경로 무변경 확인). 시험 파일 분리 뒤 재실행은 mvp-split | exit0. Go 통합 `--- PASS` 45·FAIL/SKIP 0, 실제 Node/MCP·Go owner UI. flood B 4종 모두 유효 정리 거부 0/18(최대 3.0s). 자기 Compose container·volume·network 회수 (mvp-final) |
| FullOps lint.py --from dfc70ca, deliverables strict, diff check | 기록 마지막 HEAD | 완료 보고에 기록한다. 기록 안에 자기 SHA를 순환 기록하지 않는다 |

## 경고 처리와 산출물

- 예상 규모: 제품 Go 약 +120/−35줄, 회귀 시험 약 +280줄이다. 실제: 제품 +122/−35, 전체 internal +404/−69다.
- SIZE-001: store.go 508(이전 494)·integration_test 481(이전 478)·public_messages_integration_test 410(이전 407)·PLANS(coor 누적)은 기존 대상이며 증가는 owner 공정성·회수 코드와 fixture helper·owner 4명 변경이다. 새 H-2/L-2 시험은 `cleanup_flood_integration_test.go`(195줄)로 분리해 `admission_integration_test.go`를 300줄 아래(183)로 유지했다. 무관한 분할은 하지 않았다. SIZE-002(추가 약 1.9천 줄)는 대부분 증거 로그이며 제품·시험 코드는 같은 입장 경계라 나누지 않았다. ERROR·DEP·SLOP·DESIGN 0, template·CSS 변경 없음.
- 기술 산출물: D10 module-design, D03 architecture, D05 interface-design, D06 data-model, D07 database-design(선택 JSONB 필드 Owner, table·SQL·migration 무변경), D09 crud-design을 갱신했다. D03 tech-stack은 새 기술·의존성이 없어 변경 없음이다. D08 테이블 정의는 변경 없음이다.

## 미검증·인계

- 미검증: 운영 Tunnel/edge와 실부하·상태 크기·CPU(L-1), 다수 실제 계정 공모(서로 다른 owner 4명 이상) 실측, 운영 배포·실메일·외부 계정·노우↔다닷·실24h. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다.
- 보존: OPS FIX-REVIEW H-2 high·review check exit1과 원 REVIEW-2·QA·UI 기록은 수정하지 않았다. L-1(전역 lock·색인 비용)·L-A·L-B·공개 전 합성 가입 unset·운영 DB 합성 owner 0 조건을 유지한다.
- 후속(coor 배정): 새 fixed SHA의 OPS 별도 세션 delta 보안 리뷰, TESTER 좁은 QA, 필요 시 designer UX07 수락, 그 뒤 main 판정. 미해결 critical/high는 계속 main을 차단한다.
