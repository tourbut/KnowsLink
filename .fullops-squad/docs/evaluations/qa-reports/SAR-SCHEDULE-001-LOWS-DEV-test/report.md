---
title: SAR-SCHEDULE-001-LOWS-DEV DEV 보고
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-SCHEDULE-001-LOWS-DEV]
summary: "일정 조회 리뷰 low 5건(done 재전송, 봉인 AAD, 미저장 token revoke, revoke body close, owner-revoke 연결 해제)의 수정·실패 테스트·원격 검사 결과"
---

# SAR-SCHEDULE-001-LOWS-DEV DEV 보고

## 대상과 결론

대상은 fullops/dev 제품 커밋 `08819e3`이다. 기준은 지시서 base `18d74dd`이다. 원격 검사는 `631996d`에서 실행했다. `08819e3..631996d`은 `.fullops-squad/` 아래 파일만 바꾼다. 리뷰 low 6건 중 기술 5건을 고쳤다. 남은 L2(해제 시 미전달 disclosure)는 designer 과제 SAR-SCHEDULE-001-VISUAL 몫이다. DEV 완료이며 제품 수락이 아니다.

지시서 번호와 리뷰 번호는 다르다. 아래 표는 지시서 번호(L1~L5)를 쓰고 리뷰 번호를 괄호에 적는다.

## 수정과 실패 테스트

| 지시서(리뷰) | 수정 | 테스트 | 수정 전 |
|---|---|---|---|
| L1(L1) done 재전송 | `store.go` `ingest`: done disclosure 대조를 idempotency 분기 뒤로 옮겼다. 같은 key·같은 내용은 기존 receipt, 다른 내용은 `idempotency_conflict`다 | `TestCalendarApprovalRunsOneLookup` | 403 `disclosure_denied` |
| L2(L3) 봉인 AAD | `calendar.go` `sealToken/openToken(key, data, member)`: member ID를 GCM additional data로 묶었다 | `TestSealedTokenRoundTripAndTamper` | 다른 member로 열림 |
| L3(L4) 미저장 token revoke | `google.go` callback: scope 부족이나 `finishGoogleLogin` 거부로 저장하지 않는 refresh token을 `revokeToken`한다. 회원 자신의 연결된 계정 token은 같은 grant라 revoke하지 않는다 | `TestCalendarConnection`(다른 계정 동의) | revoke 0건 |
| L4(L5) body close | `calendar.go` `revokeToken`: 응답이 있으면 상태와 관계없이 body를 닫는다 | `TestRevokeClosesBodyOnEveryStatus`(200·400·503) | 400에서 열린 채 |
| L5(L6) owner-revoke | `http.go` `owner-revoke`: 비활성화 전에 회원을 찾고 `disconnectCalendar`, handler가 lock 밖에서 `revoke`한다. 응답은 `{"state":"revoked"}` 그대로다 | `TestCalendarRevocationAndFailureClasses/owner_revoke_disconnects_and_revokes` | 봉인 token 남음 |

RED 근거:

- 단위(L2·L4): 로컬에서 해당 수정만 되돌려 실행했다. 두 테스트 FAIL, EXIT=1. 복원 뒤 PASS, EXIT=0. 로그는 `raw/red-unit.log`·`raw/green-unit.log`.
- 통합(L1·L3·L5): 제품 수정만 되돌린 임시 커밋 `11b6b03`(봉인 서명만 유지, 로컬 전용, 삭제함)을 원격 격리 clone에서 실행했다. 세 테스트가 위 '수정 전' 메시지로 FAIL, EXIT=1. 같은 테스트 함수의 다른 subtest는 PASS. 로그는 `raw/11b6b03d…-integration-red.log`.

## 확인과 판단

- 기존 봉인 이행 불필요: 캘린더는 운영에 연결되지 않았다. `operations/transition.md`의 SAR-SCHEDULE-001-DEV 절은 '운영 연결 전 초안'이고, 운영 연결은 D011 뒤 과제다. 운영 DB·env는 조회하지 않았다(지시서 승인 범위 밖). AAD 도입 전 봉인 값은 열리지 않으므로 운영 연결 전 이 변경이 들어가야 한다.
- owner-revoke 도달성: 사람 회원의 owner 기록은 credential이 비어 있다(`identity.go` `signIn`). 따라서 현재 `/v1/owner-revoke`로 사람 회원 owner를 비활성화할 수 없다. 이번 수정은 방어 경로다. 테스트는 owner credential을 직접 넣어 경로를 연다.
- L3 제외 조건: Google은 같은 계정의 grant를 token 하나로 철회한다(E7). 회원 자신의 계정이 이미 연결된 상태에서 만료 등으로 거부된 token을 revoke하면 기존 연결도 끊긴다. 그래서 그 경우만 revoke하지 않는다. 이 제외 분기는 자동 검사가 없다.
- 의존성: `go.mod`·`package.json` 변경 없음. stdlib `crypto/cipher` AAD 인자를 쓴다.

## 검사

| 검사 | 위치·대상 | 결과 |
|---|---|---|
| `go vet ./internal/relay`(태그 없음·`integration`·`integration,calendarfixture`) | 로컬 Windows `08819e3` 이후 작업 트리 | EXIT=0 |
| `go test ./internal/relay` | 로컬 Windows | EXIT=0 |
| `git diff --check 18d74dd HEAD` | 로컬 | EXIT=0 |
| `make lint` | 원격 Linux 격리 clone `631996d` | EXIT=0 |
| `make test` | 원격 | EXIT=0 |
| integration `go test -race -tags=integration ./internal/relay -skip TestRemoteConsentBrowserReturn` | 원격, 자체 Postgres 컨테이너 | EXIT=0 (`ok internal/relay 166.6s`, FAIL 0) |
| fixture-compile(`integration,calendarfixture` vet·test -c) | 원격 | EXIT=0 |
| FullOps lint `--from 18d74dd` | 로컬, `012549f`의 깨끗한 임시 worktree | ERROR 0 / WARNING 3 / 실행 불가 2, EXIT=1 |

원격 실행은 [remote.py](remote.py)다. DEV 과제의 remote.py를 복사해 `KL_REF`와 integration 실행 필터를 더했다. 완료 SHA의 gate 로그만 이 폴더에 둔다: [lint](631996d3e5363dc4fca57365e641cbf28e23e043-lint.log), [test](631996d3e5363dc4fca57365e641cbf28e23e043-test.log), [integration](631996d3e5363dc4fca57365e641cbf28e23e043-integration.log), [fixture-compile](631996d3e5363dc4fca57365e641cbf28e23e043-fixture-compile.log). SHA256은 [gate-log-sha256.txt](gate-log-sha256.txt), raw 로그 SHA256은 `raw/sha256.txt`다.

FullOps lint 상세:

- 역할 워크트리에는 지시서가 보존하라고 한 기존 미추적 파일(SAR-MOTION-INTRO 증거·`videos/service-intro/`)이 있다. lint.py는 깨끗하지 않은 트리를 거부한다. 그래서 같은 HEAD의 detached 임시 worktree에서 실행하고 지웠다. 출력은 `raw/fullops-lint.txt`.
- 실행 불가 2: product-lint·product-test. Windows에 `make`가 없다. 같은 명령은 원격 Linux에서 EXIT=0이다. EXIT=1은 이 실행 불가 때문이다.
- SIZE-001 3건: `calendar.go` 387줄(+6), `calendar_integration_test.go` 458줄(+29), `http.go` 729줄(+10). 기존 큰 파일의 좁은 델타다. 리뷰가 수락한 원래 크기에서 분할하지 않는다.
- 이후 커밋은 이 보고만 바꾼다.

## 미검증·남은 일

- 실제 Google 계정의 revoke·AAD 동작은 확인하지 않았다(D011 열림). 합성 제공자 통과는 실계정 수락이 아니다.
- L3의 '자기 연결 계정 token 유지' 분기는 테스트가 없다.
- tester 독립 QA와 delta 리뷰 필요 여부 판정은 coor 몫이다.
