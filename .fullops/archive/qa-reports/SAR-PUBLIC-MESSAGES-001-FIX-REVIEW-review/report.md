---
title: SAR-PUBLIC-MESSAGES-001-FIX-REVIEW 리뷰
status: review
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-REVIEW]
summary: 수정 후보 dfc70ca의 독립 delta 보안 리뷰로 H-1·M-1·F-UI는 해소됐고 새 high H-2로 수락 불가
---

# SAR-PUBLIC-MESSAGES-001-FIX-REVIEW 리뷰

- 검토자 / CLI / 모델: ops / Claude Code CLI / claude-opus-5-5. 실제 세션 `3cee40a4-bf54-4f2c-bad3-1559b8f6b006`.
- 구현자 세션: 수정 DEV `9ebf9a10-609b-475a-b1cd-f9edad84c9a1`, 원 DEV `01a11106-b736-7db1-b518-a65b64dbc5fb`. 이 검토 세션과 둘 다 다르다. 원 REVIEW-2 세션 `9bc44cbf-9e5c-4a16-a098-ae2fe4cab6ce`과도 다르다.
- snapshot: `/tmp/knowslink-messages-fix-review-dfc70ca`. 리뷰 전후 detached·porcelain 0줄·HEAD `dfc70caa748a90614b02d48c78b4651345938339`이다. 읽기만 했다. 설치·lint·재현은 별도 scratch clone에서 했다.
- base SHA / head SHA / merge-base: `09c523da8a3407288d9f5d711e1834af12bc7808` / `dfc70caa748a90614b02d48c78b4651345938339` / `09c523da8a3407288d9f5d711e1834af12bc7808`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate. `.fullops-squad/review/rule.json` (sha256 `ab2116fb…`). OCR LLM은 사용하지 않았다. 판정은 검토자 AI의 판단이며 OCR 자동 판정이 아니다.
- 공통 규칙: `fullops-common-0.3.3` README·coding-style·testing·security, FULLOPS.md, project.md, document-writing.md, contexts/ops.md. 지시서 준비 커밋 `b5a44fc`에서 읽었다. 예외 없음.
- 요구사항·완료 기준 원천: SAR-PUBLIC-SERVICE.md PS-04/06/07·PS-08–11과 운영 기본값 표(동시 처리 "신규 작업이 정리 budget을 소비하지 못한다", 정리 budget도 인증·현재 권한·유효 lease를 지킨다), SAR-MVP.md C1–C5, UX06/07, DEV-FIX 실행 기록 `SAR-PUBLIC-MESSAGES-001-DEV-FIX.md`·QA 증거, 원 REVIEW-2 report/result/slot-repro.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 91 / 53 / 38 / 53 / 38.
- lint(`lint.json`): ERROR 0 / WARNING 5 / 실행 불가 0. scratch clone(dfc70ca, `make install` exit0, 실행 뒤 porcelain 0)에서 `lint.py --from 09c523d` exit0. product-lint exit0, product-test exit0.
- SIZE-002: 지시서에 예상 규모 수치는 없다. 실제 추가 568줄이다. 제품 코드는 약 190줄이고 나머지는 회귀 시험·기록이다. H1·M1·F-UI는 같은 입장/gate 경계라 나누면 미완료 후보가 된다는 DEV 판단에 동의한다.
- DEP-001: 해당 없음. go.mod·go.sum·package 파일 변경 없음을 diff로 확인했다.
- UI 디자인: template·CSS 변경 없음. 바뀐 것은 gate 결정 뒤 303 Location뿐이다. 디자인 lint·테마 전환은 미구성이다(project.md). DESIGN 경고 0. 직접 시각 검수는 designer UX07 재검수 몫이다.

## 검토 범위

result.json의 91개 (path, status)를 모두 기록했다. 제품 Go 4개(capacity.go·cleanup_admission.go·http.go·store.go)는 diff와 관련 함수 전체를 읽었다. 시험 Go 4개도 읽었다.
`boundedHTTP`의 모든 호출자를 확인했다. `Handler()`와 시험뿐이다. 정리 handler와 `cleanupTarget`의 소유 조건을 하나씩 대조했다: `operate`(/v1 key-revoke·unpair·owner-revoke는 owner 자격), `/v1/text/{ack}`·`/v1/test/ack`(`leased`), `memberAction`(/home 회원 owner), `gateDecision`·`gateOwner`(gate owner·CSRF), logout.
`principal`·`session`·`hashToken`·`transaction`·`hit`·rate 표·`cmd/relay/main.go` 서버 timeout도 읽었다. 문서 16개와 JSON 29개는 내용·SHA·비밀 부재를 확인했다. JSON 29개는 모두 파싱된다.
제외 38개(.log/.exit/.txt/.head)는 rule.json exclude 증거다. DEV-FIX `.exit` 12개 값과 head.txt, mvp-final·review-repro·control 로그의 판정 줄을 직접 열어 확인했다.

| 영역 | 확인 내용 | 결과 |
|---|---|---|
| H1 본문 선수신 | 8/32KiB 상한 본문을 슬롯 전에 10s 안에서 받는다. 익명 slow 연결은 rate·슬롯·DB를 쓰지 않는다 | 해소 |
| 다른 층 이동 여부(익명) | 연결당 최대 32KiB·10s다. header 수신(1MB·5s)과 같은 연결 계층이며 새 대기열·공유 상태가 없다. 연결 수 상한은 기존 edge/OS 몫이다 | 새 무제한 점유 없음 |
| M1 익명·위조·형식 오류 | 익명·위조 token·잘못된 본문·cross-site·잘못된 gate CSRF는 로컬 정리 채널·공유 Clean·cleanup rate를 쓰지 않는다 | 해소 |
| 유효 자격의 타 owner·lease 없는 정리 | DB 입장 transaction 동안 로컬 정리 채널 4개 중 하나를 잡는다. 재분류는 그 뒤다 | H-2 |
| 로컬 색인 stale·restart·swap | 다른 프로세스의 새 자격, restart 직후 nil 색인, 커밋 뒤 Store 순서 역전 | L-2 |
| 정리 handler 소유 대조 | ACK·text/test ACK·key/agent 철회·owner-revoke·unpair·invite deny·cancel·gate deny·logout | 일치 |
| F-UI-MSG-01 | 결정 성공 뒤에만 `back+"/gates/"+id`로 303한다. back은 고정 literal이다. CSRF·현재 권한·verifiedGate·정책은 무변경이다 | 해소 |
| 공유 16+4·rate·crash 30s | `enterHTTP`·Exp·rate 표 무변경. 두 인스턴스 경합 시험 PASS | 결함 없음 |

## 발견 사항

### H-1 — high, 해소: 익명 slow body의 슬롯 점유 (원 REVIEW-2)

- 근거: `internal/relay/capacity.go:146-173`. 상한 본문 수신·ParseForm·정리 자격 판정이 모두 로컬 채널 `select`보다 앞이다.
- 검증: `TestHTTPSlowBodyAndCleanupAdmission`이 이 리뷰의 dfc70ca 실행에서 PASS다. 익명 /v1/ack 8·인증 key-revoke 8·익명 /auth/start 24개 slow 연결 뒤에도 유효 정리와 신규가 200이다. DEV 대조군은 09c에서 exit1(`cleanup behind slow bodies got 429`)이다.
- 원 REVIEW-2의 report·result·review-check exit1은 고치지 않았다. 이 기록이 새 fixed의 해소 근거다.

### M-1 — medium, 해소(익명·위조 범위): 경로만 보는 정리 입장 (원 REVIEW-2)

- 근거: `internal/relay/cleanup_admission.go:15-30`(무상태 자격·출처·CSRF·본문), `:33-67`(자기 기록), `capacity.go:29-61`.
- 검증: 위 통합 시험의 7가지 거부 요청은 Clean 기록 0이다. 단위 시험 `TestCleanupAdmissionNeedsVerifiedOwnRecord`는 12개 행렬이다. DEV M1 변형 대조군은 exit1이다.
- 남은 문제: 유효 자격을 가진 행위자의 같은 위험은 H-2로 따로 기록한다.

### H-2 — high, 미해결: 유효 자격의 타 owner·lease 없는 정리 요청이 DB 입장 동안 로컬 정리 채널을 점유한다

- 위치: `internal/relay/capacity.go:159-215`, `internal/relay/cleanup_admission.go:15-30`.
- 원인: 로컬 채널 선택은 `cleanupCredential`과 `liveCredential`만 본다. 자기 기록 여부(`cleanupTarget`)와 rate는 DB 입장 transaction 안에서 판정한다. 정리 슬롯은 그 transaction(pool 획득·전역 row lock 대기 포함) 동안 잡혀 있다. 재분류와 슬롯 반환은 그 뒤다. rate 거부 요청도 transaction 안에서 거부되므로 rate가 점유를 막지 못한다.
- 재현(dfc70ca scratch + 리뷰 시험 1개, 실제 Postgres, verify-mvp 안에서 실행):
  - A(결정적): 전역 row lock을 잡은 동안 agent_b owner가 agent_a의 키 철회 4건을 보낸다. 로컬 정리 채널 4/4가 찬다. 유효한 agent_a owner의 자기 키 철회는 `429 {"error":"capacity"}`다. 타 owner 요청은 cleanup rate 0·신규 rate 4로 정확히 집계됐다. 공유 Clean 기록은 쓰지 않는다.
  - B(lock 없음, 3s flood 48 연결): 같은 principal의 신규 작업 flood 중 유효 정리 probe의 capacity 429는 0/18이다. 타 owner 키 철회 flood 중에는 18/18이다.
  - C(공개 회원 경로, 2회차): owner 자격 없이 agent 자격만으로 보유하지 않은 lease(`{"id":"none","token":"none"}`)를 `/v1/ack`로 반복한다. 유효 정리 probe의 capacity 429는 18/18이다. 일반 회원이 가진 agent 자격만으로 충분하다.
  - A·B는 1회차와 2회차에서 같은 결과(4/4, 0/18 대 18/18)다.
- 영향: 이메일 인증 회원 하나가 자기 자격과 타인 대상 정리 경로를 반복해 relay 프로세스의 모든 유효 ACK·철회·deny·unpair·로그아웃을 429 capacity로 만든다. 키 유출 뒤 피해자의 즉시 철회도 막힌다. PS-11 "신규 작업이 정리 budget을 소비하지 못한다"와 "정리 budget도 인증·현재 권한·유효 lease를 지킨다"를 위반한다. 수정 기록의 interface-design 서술("나머지 요청은 같은 경로라도 신규 입장")도 로컬 채널에서는 성립하지 않는다. 비인증에서 인증 회원으로 비용이 올랐지만 영향 범위는 원 H-1과 같다. 전역 슬롯이 없던 base 7efbaa3 대비 이 후보 계열의 회귀다.
- 판정: 새 보호가 익명 점유를 "유효 자격 한 개의 무제한 반복"으로 옮겼다. rate·자격 철회가 자동으로 막지 못한다.
- 수정 방향(DEV 판단): 정리 슬롯을 잡기 전에 자기 기록과 principal 공정성을 확정한다. 예: 커밋 상태에서 정리 대상 소유 색인까지 만들어 로컬 선택에 쓰고, principal(회원)당 로컬 정리 동시 1개 같은 상한을 둔다. 또는 정리 판정 transaction을 슬롯 밖에서 짧게 끝낸다. 자기 대상 반복(멱등 철회) flood와 타 owner flood의 회귀 검사를 추가한다.

### L-1 — low, 미해결(확장): 요청당 전역 lock 비용

- 원 REVIEW-2 L-1을 보존한다. 이 후보는 `transaction` 커밋마다 `liveCredentials()`로 owner·agent×key·세션 전체 map을 새로 만든다(`store.go:171-172`). 입장·해제 transaction 모두 해당한다. 상태 크기에 비례한 CPU·할당이 요청마다 늘어난다. OPS 공개 전 측정에 포함한다.

### L-2 — low, 미해결: 로컬 자격 색인의 stale·restart·Store 순서

- 위치: `internal/relay/store.go:171-172`, `cleanup_admission.go:84-121`.
- 내용: 색인은 이 프로세스의 마지막 커밋 상태다. 다른 프로세스에서 막 만든 세션, restart 직후 nil 색인, 늦게 끝난 이전 transaction의 Store가 새 색인을 덮는 경우가 있다. 이때 유효 정리는 신규 채널로 간다. 신규 16개가 DB 대기 중이면 429가 된다. 막 철회된 자격은 한 transaction 동안 정리 채널을 쓴다(H-2와 같은 점유). 공유 Clean 기록·rate는 transaction이 정확히 판정한다.
- 처리: H-2 수정 때 함께 판단한다. 단독으로는 수락 차단이 아니다.

### 보존하는 기존 low·조건

- L-A(삭제된 회원 agent ID의 합성 owner 재등록)·L-B(`/v1/invite-decision` Generation 미결속)는 이 delta에 코드 변경이 없다. 미해결로 보존한다.
- 공개 전 `KNOWSLINK_SYNTHETIC_SIGNUP` 미설정과 운영 DB 합성 owner 0 확인은 계속 필수다.
- 관찰(이 delta 무변경): `/v1/connect/{cancel,key-revoke,agent-revoke,unpair}`는 handler 안에서 cleanup rate를 쓰지만 HTTP 입장 분류(`cleanupRequest`)에는 없다. 09c와 같다. H-2 수정 범위에서 PS-11 정리 경로 목록과 함께 확인을 권한다.

### lint 경고 판단

- SIZE-001 4건(PLANS·identity_integration_test·public_messages_integration_test·store.go)과 SIZE-002: DEV가 새 분류·시험을 새 파일로 분리했다. 기존 파일 증가는 회귀 시험과 색인 저장 4줄이다. 수락한다.
- SEC·SLOP·DOC·DESIGN·DEP 경고는 0이다.

## 검증 및 남은 제약

| 검사 | 대상 HEAD | 결과 | 증거 |
|---|---|---|---|
| `make install` | dfc70ca scratch clone | exit0 | install.exit |
| FullOps `lint.py --from 09c523d` (product-lint, product-test 포함) | dfc70ca, porcelain 0 | exit0, ERROR0/WARNING5/실행불가0, `make lint`·`make test` 각 exit0 | lint.json, lint.log, lint.exit, lint.head |
| `make verify-mvp` 1회차 + H-2 시험 A/B | dfc70ca + 리뷰 시험 1개(untracked) | exit0, `--- PASS` 42(제품 41 + 리뷰 1), FAIL/SKIP 0. 제품 시험 전부 PASS. 자기 Compose volume/network 회수 | run1-verify-mvp.log/.exit/.head |
| `make verify-mvp` 2회차 + H-2 시험 A/B/C | 같은 조건 | exit0, `--- PASS` 42, FAIL/SKIP 0. A 4/4·own 429 capacity, B 0/18 대 18/18, C 18/18. 자기 자원 회수(남은 container·volume 0) | verify-mvp.log/.exit/.head |
| H-2 재현 시험 원문 | — | — | review-foreign-cleanup_test.go.txt |
| `review.py check --task-key SAR-PUBLIC-MESSAGES-001-DEV-FIX` | records(리뷰 기록 작성 뒤) | exit1 `리뷰 처리 실패: critical/high 미해결 사항이 있습니다`. H-2 때문이며 정상 차단이다. 형식 probe는 쓰지 않았다 | review-check.log/.exit |

로그는 각 명령의 stdout/stderr다. `.exit`는 명령이 반환한 값이다. `git diff --check`를 위해 후행 공백만 제거했다.

- 제품 코드는 DEV 코드 SHA `1fdeaa1`과 dfc70ca 사이에 차이가 없다. DEV mvp-final(exit0)·lint-final(exit0)은 그 조건으로 유효하다. 이 리뷰는 lint·verify-mvp를 dfc70ca에서 다시 실행했다.
- DEV 원본 실패(mvp-first/second exit2, fullops-archive-4e64efd exit1)와 대조군 RED(control-09c·control-m1-mutant·control-ui-mutant exit1)는 기록 그대로 보존됐다. review-repro-fix exit1은 DB 없는 Service의 EOF라 단독 해소 근거로 쓰지 않았다. 해소 판정은 실제 Postgres 시험으로 했다.
- 원 REVIEW-2 H-1/M-1·review-check exit1·probe 기록은 변경하지 않았다. 이번 리뷰는 result 사본의 high를 바꾸는 probe를 쓰지 않았다.
- 실행하지 않음: 운영 Tunnel/edge 본문 완충과 연결 수 상한, 실부하·상태 크기·CPU(L-1), 실제 브라우저 시각 검수(designer UX07), 실메일·실24h·외부 계정·운영 배포. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다. edge 미검증을 통과 근거로 쓰지 않았다.
- skipped 38개는 원시 로그·종료코드 증거다. 판단에는 `.exit` 값과 판정 줄을 사용했다. 생략 영향은 없다.
- records 최종 HEAD의 lint/test·deliverables strict·diff 검사는 완료 보고에 기록한다. 기록 안에 자기 SHA를 순환 기록하지 않는다.

## 검토 결론

수락 불가다. 원 H-1(익명 slow body)과 M-1(익명·위조 정리 입장)은 dfc70ca에서 해소됐다. F-UI-MSG-01 gate Deny 405도 해소됐다.
그러나 새 미해결 high H-2가 main 수락을 차단한다. 인증 회원 하나가 타인 대상·lease 없는 정리 요청을 반복해 DB 입장 동안 로컬 정리 채널 4개를 모두 점유한다. 실제 Postgres에서 유효 정리 18/18이 429 capacity였다. PS-11 정리 budget 보장이 무효가 된다.
L-1(확장)·L-2와 기존 L-A/L-B는 OPS 공개 전 조건으로 보존한다.
재개 조건: DEV가 H-2(권장 L-2·/v1/connect 정리 경로 확인 포함)를 고치고 자기 대상 반복·타 owner flood 회귀 검사를 추가한 새 fixed SHA를 낸다. 그 뒤 OPS가 별도 세션의 새 리뷰 키로 재검토한다.
이 결론은 기술 리뷰 판단이다. TESTER 독립 QA·designer 직접 시각 검수·운영 공개 수락을 대신하지 않는다.
