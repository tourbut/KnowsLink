---
title: HTTP 입장 H1/M1·gate Deny 405 수정 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX]
summary: OPS REVIEW-2의 slow-body 슬롯 점유(H-1)·경로만 보는 정리 입장(M-1)과 designer F-UI-MSG-01 gate Deny 405의 원인·수정·회귀 검증을 기록한다
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX — HTTP 입장 H1/M1·gate Deny 405 수정 기록

## 기준과 기술 계획

- 기준 ref: 원 fixed `09c523da8a3407288d9f5d711e1834af12bc7808`. 지시서 준비 커밋 `9fd8665`에서 FULLOPS·project·fullops-common-0.3.3·contexts/dev·SAR-PUBLIC-SERVICE PS11·SAR-MVP HTTP 표와 REVIEW-2 report·slot-repro를 읽었다.
- 추가 지시: designer handoff msg_6194f4825058 F-UI-MSG-01(coor 인박스의 "같은 과제 추가 수정" 절). 같은 과제에서 처리했다.
- 계획: (1) 본문을 모든 슬롯보다 먼저 받는다. (2) 정리 분류를 경로가 아니라 검증된 principal·자기 기록에 묶는다. (3) 로컬 정리 채널은 위조할 수 없는 자격만 고른다. (4) gate 결정 뒤 정식 화면으로 보낸다. 제품 수치(신규16·정리4·rate·10s·30s·8/32KiB)와 wire·응답 코드는 바꾸지 않는다.
- 공통 경계: 모든 HTTP 요청은 `boundedHTTP` 한 곳을 지난다. 다른 호출자는 없다(`Handler()`와 시험만). 수정은 이 경계와 그 분류 함수에만 둔다.

## 원인과 수정

### H-1 slow body가 슬롯을 점유

- 원인: `boundedHTTP`가 로컬 채널 슬롯을 잡은 뒤 `io.ReadAll`로 본문을 읽었다. 인증·rate는 본문 뒤 DB에서 했다. 익명 송신자가 본문을 늦게 보내면 슬롯이 10s 동안 묶였다.
- 수정: 상한 본문(form·connect 8KiB, 나머지 32KiB)을 같은 10s read deadline 안에서 먼저 받는다. 그 뒤 로컬 채널→DB 입장 transaction을 한다. 느린 송신자는 자기 연결만 점유한다. 이 점유는 header 수신(ReadHeaderTimeout 5s, header 최대 1MB)과 같은 연결 계층이다. 새 대기열·DB 기록·공유 상태는 없다. 연결 수 자체의 상한은 edge/OS 계층이며 이번에 확인하지 않았다.
- 본문을 다 받은 malformed 요청은 이전처럼 입장에서 rate를 한 번 쓴다. 로컬 채널 포화와 deadline으로 입장 전에 끝난 요청은 rate를 쓰지 않는다(기존 "입장 전 로컬 거부는 DB 호출 없음"과 같다).

### M-1 정리 입장을 경로만으로 판정

- 원인: `cleanupRequest`는 경로·form 값만 보았다. 익명·위조 자격 요청도 로컬 정리 채널과 공유 Clean 기록을 썼다.
- 수정(`cleanup_admission.go`):
  - `cleanupCredential`: 정리 경로·같은 출처(`CrossOriginProtection.Check`)·gate HMAC CSRF·/v1 closed 본문을 상태 없이 확인하고 해당 자격을 돌려준다.
  - `cleanupTarget`: transaction 안에서 각 정리 handler의 소유 조건을 그대로 대조한다. ACK는 `leased`(유효 lease·text/일반 경로 일치), 키·agent·unpair·초대 거절·cancel·gate는 자기 owner, owner-revoke는 owner principal, logout은 세션이다. handler가 여전히 최종 권한이다. 이 함수는 budget만 고른다.
  - `requestBuckets`는 principal이 있고 위 둘을 통과할 때만 cleanupRate와 Clean 기록을 준다. 나머지는 anonymousRate 또는 memberRate·신규 기록이다.
  - 로컬 채널 선택: DB 전에는 principal을 확인할 수 없다. 커밋 직후 상태에서 유효 자격 해시 색인을 만들어 `Service.live`에 둔다. 색인에 있는 자격만 정리 채널을 쓴다. 위조 token은 정리 채널에 들어가지 못한다. transaction 판정이 신규면 신규 슬롯으로 바꾸고 정리 슬롯을 돌려준다. 신규가 가득하면 429 capacity다.
- 한계: 색인은 이 프로세스의 마지막 커밋 상태다. 다른 프로세스에서 막 만든 자격은 이 프로세스의 다음 transaction 전까지 신규 채널을 쓴다. 이때도 공유 Clean 기록·cleanup rate는 transaction이 정확히 판정한다. 막 철회된 자격은 한 transaction 동안 정리 채널을 쓸 수 있으나 rate·기록은 신규로 집계된다.

### F-UI-MSG-01 gate Deny 뒤 405

- 원인: `gateDecision`이 `r.URL.Path`로 303했다. Deny form은 `/…/gates/{id}/deny`로 POST하므로 브라우저가 POST 전용 경로를 GET해 405가 됐다. 거절 저장은 성공했다.
- 수정: 두 결정 모두 `back+"/gates/"+id`(정식 gate 화면)로 303한다. 이 화면은 결정 상태·비활성 버튼·작업 화면 링크를 보여 준다. 화면 template·CSS는 바꾸지 않았다.

## 자동 검증 증거

증거: [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-DEV-FIX/). `.exit`는 명령 반환값이다. 로그는 후행 공백만 제거했다.

| 검사 | 대상 | 결과 |
|---|---|---|
| REVIEW-2 원 재현(slot-repro, DB 없음) | 09c | exit0, 정리·신규 probe 429 capacity 재현 (review-repro-09c) |
| 같은 재현 | 수정 코드 SHA | exit1. probe가 로컬 슬롯을 지나 DB 입장까지 간다(DB 없는 Service라 EOF). 429 capacity 아님 (review-repro-fix) |
| 새 통합 검사 대조군 | 09c + 새 시험 | exit1 `cleanup behind slow bodies got 429` (control-09c) |
| M-1 변형(H1 수정 유지, 경로만 분류) | 수정 코드 변형 | exit1 `admitted 4 want 7` (control-m1-mutant) |
| F-UI 변형(옛 redirect) | 수정 코드 변형 | exit1 member/owner 모두 `deny redirected to …/deny` (control-ui-mutant) |
| make lint / make test | 코드 SHA | exit0 / exit0 |
| make verify-mvp | 코드 SHA | exit0. 격리 Compose·Postgres·실제 Node/MCP. 자기 project volume·network 회수 |
| FullOps lint.py --from 09c | 코드 SHA | exit0. ERROR0·WARNING5·실행불가0. product-lint/test passed |

새 검사:
- `TestHTTPSlowBodyAndCleanupAdmission`(통합): 익명 /v1/ack 8·인증 key-revoke slow 8·익명 /auth/start 24 연결을 유지한 채 유효 owner key-revoke와 신규 /v1/registry가 200이다. slow 연결은 rate를 쓰지 않는다. 본문 완성 뒤 32개 익명 요청은 30 허용·2개 429이며 ip bucket은 31개다. 이어 익명·위조 token·타 owner 키·잘못된 본문·cross-site·lease 없는 ACK·잘못된 gate CSRF 7건을 붙잡아도 Clean 기록 0이다. 자기 정리 4건은 모두 Clean으로 입장하고 5번째는 로컬 정리 채널에서 429다. 타 owner 요청은 cleanup rate를 쓰지 않고 신규 rate 2회다. 종료 뒤 입장 기록 0이다.
- `TestCleanupAdmissionNeedsVerifiedOwnRecord`(unit): 유효 자격 색인(철회 agent·위조 제외)과 12개 분류 행렬.
- `TestEmailIdentity/member_gate_deny_form_returns_to_result`(통합): 실제 렌더된 form의 formaction·csrf로 Deny를 보내 303 `/home/gates/{id}`, 그 화면 200 `상태: denied`·버튼 비활성·`href="/home"`, 홈 목록 유지를 확인한다. `TestPublicHTTPAdmissionAndGateSafety`는 owner deny Location을 확인한다.
- 기존 `TestHTTPConcurrencyAcrossInstances`의 정리 요청은 이제 자기 키 대상 본문을 보낸다. 빈 `{}` 철회는 신규 작업이기 때문이다. 두 인스턴스 16/4 경합·종료 회수 의미는 같다.

원본 실패 보존: mvp-first(exit2)는 새 시험이 slow 연결 종료 직후 rate를 확인한 경쟁이었다. 본문을 완성하고 응답을 읽도록 고쳤다. mvp-second(exit2)는 5번째 정리 요청의 rate 기대값 오류였다. 로컬 채널 거부는 rate를 쓰지 않는 기존 규칙이 맞다. 두 실패 모두 시험 작성 오류이며 제품 결함이 아니다. 아카이브 커밋 `4e64efd`의 첫 FullOps 게이트는 exit1이었다(DOC-003, 이 기록의 summary 따옴표 형식). `deliverables.py --stamp`로 front matter를 다시 써서 고쳤다(fullops-archive-4e64efd). 대조군·변형은 파일 분리 전 동일 로직 트리에서 실행했다. 원 재현의 수정 측(review-repro-fix)은 코드 SHA에서 다시 실행했다.

## 경고 처리와 산출물

- SIZE-001 4건: PLANS.md(coor 기록 누적), identity_integration_test 359, public_messages_integration_test 407, store.go 494. 새 정리 분류는 `cleanup_admission.go`, 새 시험은 `admission_integration_test.go`·`cleanup_admission_test.go`로 분리해 기존 파일 증가를 줄였다. 기존 파일의 증가는 F-UI 회귀 28줄·owner deny 확인·색인 저장 4줄이다. 무관한 분할은 하지 않았다.
- SIZE-002: 추가 469줄 중 제품 코드 약 190줄, 나머지는 회귀 시험이다. H1·M1·F-UI는 같은 입장/gate 경계라 나누면 미완료 후보가 된다.
- DOC-001은 새 시험 파일 header로 해소했다. DEP0·새 의존성 없음. DESIGN: template·CSS 변경 없음. 별도 theme/design lint는 미구성(project.md)이다.
- 기술 산출물: D10 module-design, D03 architecture·tech-stack, D05 interface-design, D06 data-model, D09 crud-design을 갱신했다. D07 database-design은 변경 없음이다. table·SQL·migration·JSONB 필드가 그대로이며 색인은 저장하지 않는 프로세스 메모리 값이다.

## 미검증·인계

- 미검증: 운영 Tunnel/edge 본문 완충과 연결 수 상한, 실부하·상태 크기·CPU(L-1), 실제 브라우저의 키보드 Enter 직접 확인, 운영 배포·실메일·외부 계정·노우↔다닷·실24h. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다.
- 보존: 원 OPS REVIEW-2 H-1/M-1·review check exit1과 원 QA/UI 실패는 09c 기록 그대로다. L-1·L-A·L-B·합성 가입 unset·운영 DB 합성 owner 0 조건을 유지한다.
- 후속(coor 배정): 새 fixed SHA의 OPS 별도 세션 delta 보안 리뷰, TESTER 좁은 QA, designer UX07 직접 재검수, 최종 검사 수락 뒤 main 판정. 미해결 critical/high는 계속 main을 차단한다.
