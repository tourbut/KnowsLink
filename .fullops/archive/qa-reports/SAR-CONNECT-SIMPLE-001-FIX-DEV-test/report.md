---
title: SAR-CONNECT-SIMPLE-001-FIX-DEV 자체 검사 보고
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-FIX-DEV]
summary: "REVIEW2 후속 4건(complete agent_rate retry_at, 구버전 시도 agent ID 비재발급, D12 배포 순서, attempt-remove 음성 테스트)을 고치고 Linux 격리 make lint·make test·integration EXIT=0으로 검사했다."
---

# SAR-CONNECT-SIMPLE-001-FIX-DEV 자체 검사 보고

독립 리뷰 [REVIEW2](../SAR-CONNECT-SIMPLE-001-DEV-REVIEW2-review/report.md)의 medium 1·low 3 수정에 대한 DEV 검사다. 테스트 레벨은 lite다. 운영 relay·운영 agent·실계정은 쓰지 않았다.

- 브랜치 `fullops/dev`, 코드 SHA `fd0eb2d5b02e`(기준 `116e3e8a`).
- 원시 로그: `D:/workspace/KnowsLink/.devtools/connect-simple-001-fixdev/`(gitignore 공용 폴더). 파일 이름은 `<SHA>-lint.log`·`-test.log`·`-integration.log`·`fullops-lint.log`다.

## 발견 사항별 수정과 검사

| ID | 수정 | 검사 | 결과 |
|---|---|---|---|
| M1 | `member_agents.go` `connectAPI`의 두 오류 응답이 `errorBody(e)`를 쓴다. `agent_rate` 거절에 `retry_at`이 붙는다 | integration `TestGoogleDeviceHTTP`: 24시간 안 agent 19개 회원이 두 시도를 동의하고 키를 저장한다. 첫 complete 200, 두 번째 429 `"error":"agent_rate"`·`"retry_at"` | 통과. 수정 전 코드로 RED 실행은 하지 않았다(Linux 게이트 1회만 실행). 수정 전 응답은 `retry_at`이 없어 이 단언이 실패한다 |
| L1 | `connections.go` `connection()`의 예약 판정을 `Connection.Name`이 있는 시도(이 흐름의 동의)로 한정한다. 구버전 시도의 agent가 철회·정리된 뒤 complete는 `invalid_auth`다(C1) | `TestOldAttemptNeverRemakesAgent`. 수정 전 코드에서 실패하고 수정 뒤 통과함을 확인했다 | 통과 |
| L2 | D12 `ops-guide.md` SAR-CONNECT-SIMPLE-001 절에 'relay를 클라이언트·설치기보다 먼저 배포' 한 줄 | 문서 | 반영 |
| L3 | 코드 변경 없음 | integration `TestGoogleDeviceHTTP`: 다른 회원의 시도와 완료(consumed)된 시도의 `/home/attempt-remove`가 403이고 두 시도가 남는다 | 통과 |

D10 `module-design.md`의 A-01–A-03 행에 예약 한정·`retry_at`·정리 거절과 새 테스트를 더했다.

## 필수 검사

| 검사 | 환경 | 결과 |
|---|---|---|
| `make lint` | Linux 격리 clone, `fd0eb2d5` | EXIT=0 |
| `make test` | 같은 clone | EXIT=0 |
| `go test -race -tags=integration ./internal/relay` | 같은 clone, 격리 Postgres 17 컨테이너(검사 뒤 정지) | EXIT=0, `ok … 147.414s` |
| `go test ./internal/...` | Windows, Go 1.27.1 | 통과 |
| FullOps `lint.py --from 116e3e8a` | `fd0eb2d5`의 깨끗한 detached worktree | ERROR 0, WARNING 2(SIZE-001 `PLANS.md`·`ops-guide.md` 기존 500줄 초과), 실행 불가 2(`make`가 Windows에 없음. 위 Linux 결과로 대신한다) |

dev 작업 트리에는 다른 과제(SAR-MOTION-INTRO)의 미추적 파일이 있다. 그래서 FullOps lint를 깨끗한 worktree에서 돌렸다. 그 파일은 옮기거나 지우지 않았다.

## 미실행과 후속

- tester 독립 QA(CS-AC-01–16), designer 시각 검수(CS-AC-17), 배포 뒤 실연결(CS-AC-18)은 coor가 이 SHA 뒤에 배정한다.
- 화면·문구 변경은 없다.
