---
title: SAR-SCHEDULE-001-TESTGAP-DEV DEV 보고
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-SCHEDULE-001-TESTGAP-DEV]
summary: "독립 QA not_run 세 경로(SC-04 친구 아님·철회 세대, SC-06 친구 해제·차단, L3 자기 계정 재동의 실패)의 integration 테스트와 검사 결과"
---

# SAR-SCHEDULE-001-TESTGAP-DEV DEV 보고

## 대상과 결론

대상은 fullops/dev 완료 SHA `4a708c1`이다. 기준은 지시서 base `884c6bc`다. 세 경로에 integration 테스트를 추가했고 모두 첫 실행부터 기대대로 통과했다. 제품 코드는 바꾸지 않았다. 따라서 D10·`docs/development.md` 갱신과 tester 재QA는 필요 없다.

## 추가한 테스트

파일: `internal/relay/calendar_testgap_integration_test.go`(태그 `integration`, 기존 `calendarSetup`·`gateFor`·합성 provider 재사용).

| 항목 | 테스트 | 확인한 것 |
|---|---|---|
| T1 SC-04 친구 아님 | `TestCalendarQueryOutsideRelationship/stranger_refused_at_accept` | 관계 없는 `agent_c`의 `schedule.query`가 수락 단계에서 403. 메시지·gate 없음, free/busy 호출 0 |
| T1 SC-04 철회 세대 | `.../revoked_generation_not_delivered` | gen1에서 수락한 요청이 unpair 뒤 다음 세대가 active여도 `failed:revoked`, envelope 없음, claim 403, 요청자 receipt 403, 호출 0 |
| T2 SC-06 | `TestCalendarFriendEndedStopsQuery/{remove,block}_approved_{false,true}` | `/home/friends/decision` 해제·차단 직후 gate `revoked`. 대기 요청은 승인 409·조회 0. 승인·조회 완료 뒤 미consume 요청은 consume 409·disclosure 삭제·조회 1회 그대로. 요청자 쪽은 `relay.result` 전송 403, 요청자 agent receipt 403 |
| T3 L3 | `TestCalendarOwnAccountRefusalKeepsGrant` | 이미 연결된 자기 계정(`sub-bob`)의 재동의가 `signIn` 거부(시작 세션 종료)로 403일 때 provider revoke 0, 기존 봉인 토큰·`connected` 유지. 다른 계정 revoke는 기존 `TestCalendarConnection`이 짝이다 |

판단:

- 철회 세대 재연결은 human 회원이 있으면 `/v1/invites`가 `human_invite_required`로 막힌다. 그래서 다음 세대 active pair를 상태로 직접 만들었다. 첫 실행 기록은 `raw/4df22f6…-integration-red.log`이며 테스트 준비 오류이고 제품 결함이 아니다.
- T3의 "권한 일부 거절"은 `google.go:225` 분기다. 이 분기는 `include_granted_scopes=true`에서 응답 scope에 freebusy가 없을 때만 탄다. 그때 기존 grant에도 freebusy가 없으므로 revoke해도 공유 grant 손실이 아니다. 합성 provider는 freebusy 없는 동의에 refresh token을 주지 않는다. 그래서 테스트는 `shared` 검사가 지키는 결과 거부 분기(`google.go:240`, `256`)를 실행했다.

## 검사

| 검사 | 위치·대상 | 결과 |
|---|---|---|
| `go vet ./internal/relay`(태그 없음·`integration`·`integration,calendarfixture`) | 로컬 Windows `4a708c1` | EXIT=0 |
| `go test ./internal/relay` | 로컬 Windows `4a708c1` | EXIT=0 |
| `make lint` / `make test` / integration | 원격 Linux 격리 clone `4a708c1` | [lint](4a708c159919a843db6eb32afc5c39882a62a62d-lint.log) EXIT=0 / [test](4a708c159919a843db6eb32afc5c39882a62a62d-test.log) EXIT=0 / [integration](4a708c159919a843db6eb32afc5c39882a62a62d-integration.log) EXIT=0 (`-skip TestRemoteConsentBrowserReturn`, `ok internal/relay 156.4s`, FAIL 0) |
| FullOps lint `--from 884c6bc` | 로컬, 깨끗한 임시 detached worktree | 아래 상세 |

원격 실행은 [remote.py](remote.py)(VISUAL-FIX 복사본, 로컬 로그 경로·clone 접두어만 바꿈)다. gate 로그 SHA256은 [gate-log-sha256.txt](gate-log-sha256.txt)다. `4a708c1` 이후 커밋은 `.fullops-squad/` 아래 보고·지시서·패킷 결과만 바꾼다.

FullOps lint 상세: 완료 커밋 뒤 기록한다.

## 미검증·남은 일

- T3 테스트가 `shared` 검사를 지우면 실패하는지(mutation)는 실행하지 않았다. 단언은 provider revoke 목록이 비어 있는지다.
- 의존성 변경 없음. `go.mod`·`package.json`을 바꾸지 않았다.
- tester 재QA 불필요(테스트만 추가). SC-09 실계정은 계속 D011 대기다.
