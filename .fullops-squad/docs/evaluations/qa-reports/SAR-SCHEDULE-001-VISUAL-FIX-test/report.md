---
title: SAR-SCHEDULE-001-VISUAL-FIX DEV 보고
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-SCHEDULE-001-VISUAL-FIX]
summary: "일정 조회 시각 검수 수정 4건(V1 조회 상세 한 줄, V2 summary 기호, V3 reconnect 해제 details, D2 해제 문구와 미전달 결과 삭제)의 구현·실패 테스트·검사 결과"
---

# SAR-SCHEDULE-001-VISUAL-FIX DEV 보고

## 대상과 결론

대상은 fullops/dev 제품 커밋 `44ecde1`이다. 기준은 지시서 base `813e2d3`이다. 수정 4건(V1·V2·V3·D2)을 구현했다. DEV 완료이며 제품 수락이 아니다. 시각 수락은 designer, 독립 QA는 tester 몫이다.

## 수정과 테스트

| 항목 | 수정 | 테스트 |
|---|---|---|
| V1 | `human_pages.go` `human-gate`: 캘린더가 연결되지 않은 `schedule.query`(이메일 회원 포함, reconnect 포함)는 범위 두 줄 대신 '이 요청은 승인할 수 없습니다. 거절하거나 기한이 지나면 조회 없이 끝납니다.' 한 줄. `schedule.commit`은 `.Calendar`가 비어 기존 U5 두 줄 그대로 | `TestCalendarUnconnectedAndCommitStub`(한 줄 있음, 두 줄 없음, commit U5 유지) |
| V2 | `member_pages.go` 공용 CSS: 자식 요소가 없는 `summary`는 `flex-wrap:nowrap;align-items:baseline;padding-block:9px`, `summary::before{flex:none}`. 기호가 첫 줄 기준선에 붙고 넘친 문구만 다음 줄로 간다. `h3`·`span`을 가진 summary(에이전트 행, 랜딩)는 바뀌지 않는다 | 렌더 probe의 `split` 측정(아래) |
| V3 | `human-calendar-disconnect` 템플릿을 새로 두고 connected·reconnect 둘 다 쓴다. '캘린더 다시 연결' 주 버튼은 그대로 | `TestHumanAccountCalendarStates`(모든 해제 form이 위험 details 안, 안내 문장 뒤) |
| D2 문구 | 해제 안내와 `?n=calendar-off` 안내를 report D2 KO/EN 문장으로 교체 | 위 테스트, `TestEnglishScreensHaveNoKorean` |
| D2 서버 | `calendar.go` `disconnectCalendar`: 조건을 `pending || (approved && !Consumed)`로 바꿨다. 조회 완료·미consume gate도 `denied`, `Lookup.State=failed`, disclosure nil. 요청자는 대기 요청과 같은 `denied` 결과를 받는다. consume된 gate는 `approved`로 남는다 | `TestCalendarRevocationAndFailureClasses/disconnect_drops_undelivered_result`, `disconnect_keeps_consumed_result` |

RED 근거: 테스트만 담은 커밋 `eeec681`(제품 코드는 base)을 원격 격리 clone에서 실행했다. `TestCalendarUnconnectedAndCommitStub`(V1 한 줄 없음)와 `disconnect_drops_undelivered_result`(`state:approved lookup:done`)가 FAIL, EXIT=1. 같은 함수의 다른 subtest는 PASS. 로그 `raw/eeec6819…-integration-red.log`. V3 단위 단언은 base 템플릿의 reconnect 해제 form이 details 밖이라 실패하는 조건이다(실행 기록 없음).

## 검사

| 검사 | 위치·대상 | 결과 |
|---|---|---|
| `go vet ./internal/relay`(태그 없음·`integration`·`integration,calendarfixture`) | 로컬 Windows | EXIT=0 |
| `go test ./internal/relay` | 로컬 Windows `44ecde1` | EXIT=0 |
| `make lint` / `make test` / integration | 원격 Linux 격리 clone `44ecde1` | [lint](44ecde1223d5dc1fd0e936e17f1879600b7d3750-lint.log) EXIT=0 / [test](44ecde1223d5dc1fd0e936e17f1879600b7d3750-test.log) EXIT=0 / [integration](44ecde1223d5dc1fd0e936e17f1879600b7d3750-integration.log) EXIT=0 (`-skip TestRemoteConsentBrowserReturn`, `ok internal/relay 174.4s`, FAIL 0) |
| 렌더 probe | 로컬 headless Chrome, 10화면 × 1440·390·320 | 30개 가로 넘침 0, text-only summary 분리 0 |
| FullOps lint `--from 813e2d3` | 로컬, 깨끗한 임시 detached worktree | 아래 상세 |

렌더: [probe.mjs](probe.mjs)는 DEV probe를 복사해 text-only summary의 문구 시작 x가 summary 왼쪽에서 20px 안이면(기호와 다른 줄) `split`으로 기록한다. 대상은 gate-query-{connected, unconnected, approved}, human-account-calendar-{connected, reconnect}와 각 영어판이다. 결과 `raw/probe.json`, 캡처 `raw/*.png` 30장. 캡처의 details는 닫힌 상태다. 열린 해제 안내 문장은 단위 테스트로만 확인했다.

이후 커밋은 `.fullops-squad/` 아래 보고·지시서·패킷 결과만 바꾼다. 원격 실행은 [remote.py](remote.py)(LOWS-DEV 복사본, 경로·접두어만 바꿈)다. gate 로그 SHA256은 [gate-log-sha256.txt](gate-log-sha256.txt)다.

## 확인과 판단

- 로컬 PATH의 `gofmt`가 toolchain `$(go env GOROOT)/bin/gofmt`와 정렬 결과가 달라 원격 `make lint`가 처음 `2e5020f`에서 EXIT=2였다. GOROOT gofmt로 다시 맞춰 `44ecde1`에서 다시 실행했다.
- `gateScopes["schedule.query"]` 문구는 이제 화면에 나오지 않는다. 테스트가 금지 문자열로 쓰므로 남겼다.
- 의존성 변경 없음. `go.mod`·`package.json`을 바꾸지 않았다.

## 미검증·남은 일

- 열린 details 상태의 시각 캡처와 시각 수락은 designer 몫이다.
- `:has()`를 지원하지 않는 브라우저는 V2 이전 배치(기호 단독 줄)로 보인다. 기능 손실은 없다.
- tester 독립 QA(SC-09, LOWS-DEV 포함)와 designer 영향 화면 재확인은 coor가 배정한다.
