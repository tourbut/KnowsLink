---
title: SAR-MVP-001-TESTER-FINAL — RF-01 좁은 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: "제품 78b1d92의 legacy claim 차단과 정상 agent, owner, 현재 권한 회귀를 기록한다"
---

# SAR-MVP-001-TESTER-FINAL — RF-01 좁은 독립 QA 보고서

## 판정

제품 후보는 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`다. 검사 트리의 제품 경로는 이 커밋과 같다.
`a6a10c7`가 직렬화한 human claim과 경로 미기록 agent claim은 authorize, `relay.approval.request`, `relay.result`에서 `403 sender_not_allowed`다.
승인된 legacy gate의 consume도 `403 sender_not_allowed`다. gate는 소비되지 않았고 completion은 비어 있으며 claim token은 남아 있다.
같은 Postgres 상태에서 새 agent 전달의 send, pull, persist, ACK, claim, authorize는 통과했다.
owner gate, approval, denied result, current-auth 회귀도 통과했다.
`4262d02`에서 열린 legacy claim high는 이 경계에서 재현되지 않았다.
이 판정은 전체 MVP 수락이나 공개 배포 완료가 아니다. 원래 held는 유지한다. 제품 코드는 수정하지 않았다.
시나리오: [SAR-MVP-001-TESTER-FINAL.md](../scenarios/SAR-MVP-001-TESTER-FINAL.md). 로그: [SAR-MVP-001-TESTER-FINAL-test/](SAR-MVP-001-TESTER-FINAL-test/).

## 환경

날짜는 2026-10-03이다. 브랜치는 `fullops/tester`다. 검사 시작 HEAD는 `11e3ff3689dc32fc879b892e25845609c6374f0b`다.
공통 기준은 `fullops-common-0.3.2`다. 적용 문서는 `rules/common/README.md`와 코딩·테스트·보안 규칙, `project.md`, 문서 작성 규칙, protocol C1이다.
Go 1.27.1, Node v22.22.2, Python 3.12.3, Docker 29.4.3, Compose v5.1.3이다. 기록은 [versions.log](SAR-MVP-001-TESTER-FINAL-test/versions.log)다.
합성 credential, CSRF, claim token, DB 비밀번호는 로그에 남기지 않았다. Tunnel profile은 켜지 않았다. Compose project는 종료 시 그 project만 `down --volumes` 했다.
Jev는 다시 실행하지 않았다. 근거는 `SAR-MVP-001-TESTER-find.json`, `SAR-MVP-001-TESTER-documents-find.json`, `SAR-MVP-001-TESTER-context.json`, `SAR-MVP-001-TESTER-FIX-route.json`이다.

## 실행한 경계

fixture는 `internal/relay/testdata/legacy_claims.json`이다. DEV가 `a6a10c7` detached worktree에서 생성한 파일을 그대로 읽었다. 이 QA는 fixture를 다시 만들지 않았다.
파일의 메시지 네 개에는 `Deliver` 키가 없다. human 요청, 그 요청의 승인 gate, 경로 미기록 agent 요청이 들어 있다. claim은 이전 State 메서드가 발급한 값이다.
원래 만료는 `2026-10-03T10:04:00Z`다. 이 QA의 시계는 그 TTL을 지난 뒤였다.
로드 전에 `Receipt.exp`와 `Gate.Exp`만 4분 뒤로 옮겼다. retention sweep가 만료된 claim token을 지우면 경로 거부와 만료가 섞이기 때문이다.
`Deliver` 키, claim token, 승인 상태는 바꾸지 않았다.
로드와 형태 확인은 한 SQL이다. 결과는 `1|4|true|true|approved|false`다. 행 1개, 메시지 4개, `Deliver` 키 없음, 세 claim은 delivered이고 token이 있다. gate는 approved이고 소비되지 않았다.
relay의 1초 retention은 영값 `Deliver`를 빈 문자열로 다시 저장한다. 빈 문자열은 `agent`가 아니다. HTTP 거부는 그 뒤에도 `sender_not_allowed`다.

| 검사 | 결과 | 관찰 |
|---|---|---|
| FIX-load-fixture | 통과 | 원자적 로드 결과 `1\|4\|true\|true\|approved\|false` |
| FIX-legacy-human | 통과 | authorize, H, R이 `403 sender_not_allowed` |
| FIX-legacy-unrouted | 통과 | 경로 미기록 agent claim도 같은 세 거절 |
| FIX-legacy-consume | 통과 | 승인 gate consume `403 sender_not_allowed` |
| FIX-legacy-unchanged | 통과 | completion 공백, `Deliver` 공백, claim 유지, gate 미소비, state `approved` |
| FIX-legacy-still-denied | 통과 | 새 정상 흐름 뒤 human authorize가 다시 `403 sender_not_allowed` |
| REG-current-auth | 통과 | agent credential의 owner-revoke는 `401 invalid_auth` |
| REG-agent-delivery | 통과 | 새 agent send, persist, ACK, claim, authorize `200`. executable false, disclosure false |
| REG-owner-approval-result | 통과 | H `200`. owner 페이지 `200`. agent 페이지 `401`. 잘못된 CSRF `403`. approve `303`. consume `200`. result `200`. completion `denied` |

메모리 단의 `TestLegacyClaimsCannotReachParentBoundaries` 종료코드는 0이다. 이 명령은 Postgres HTTP 증거를 대신하지 않는다. 로그는 [unit.log](SAR-MVP-001-TESTER-FINAL-test/unit.log)다.

## 명령과 종료코드

| 명령 | 종료코드 | 로그 |
|---|---|---|
| 제품 diff `78b1d92`와 HEAD | 0 | [versions.log](SAR-MVP-001-TESTER-FINAL-test/versions.log). 제품 경로 출력은 비어 있다 |
| `4262d02`와 `78b1d92`의 제품 diff | 0 | 같은 로그. 변경은 `store.go`, `integration_test.go`, `legacy_test.go`, `legacy_claims.json`이다 |
| `go test` legacy fixture | 0 | [unit.log](SAR-MVP-001-TESTER-FINAL-test/unit.log) |
| `isolated.py` 1회 | 1 | [probe.log](SAR-MVP-001-TESTER-FINAL-test/probe.log). `?` 형태 조회가 기대를 벗어났다. 제품 판정 아님 |
| `isolated.py` 2회 | 1 | 같은 로그. retention 재저장과 형태 조회가 겹쳤다. 제품 판정 아님 |
| `isolated.py` 3회 | 1 | 같은 로그. 형태 값 `1\|4\|true\|true\|approved\|false`를 `t`로 비교했다. 제품 판정 아님 |
| `isolated.py` 4회 | 0 | 같은 로그의 마지막 `[exit 0]`. project `knowslink-final-ed42769aa3`. `down --volumes` 종료코드 0. 판정 실행 |

`make verify-mvp`와 `make verify-runtime`은 실행하지 않았다. 설정, migration, adapter, cmd, db, owner HTML은 `4262d02`와 같다. 이번 실제 DB 증거는 고유 Compose의 HTTP probe다.
앞의 세 probe에서 제품 호출은 이미 기대 거절과 정상 회귀를 반환했다. 실패 줄은 tester의 형태 비교다. 네 번째 실행만 판정이다.

## 재사용

`4262d02` QA 보고서와 `SAR-MVP-001-TESTER-FIX-test/`는 수정하지 않았다. 그 실행의 pass, fail, held는 `4262d02` 조건으로 남긴다.
`a6a10c7`의 QA-01–11과 8개 held도 그 SHA의 기록으로 남긴다.
owner HTML, Compose, Makefile, adapter, cmd, db는 `4262d02`와 `78b1d92` 사이에서 바뀌지 않았다. PNG 8개의 동일성 증거는 `c59537b` 7개와 `e238777` 1개의 원래 SHA와 조건으로 재사용한다. 새 캡처는 없다.
designer 시각 판정은 `e238777`이다. 이번 실행의 시각 통과로 바꾸지 않는다.
QA-11 24시간, WAL, QA-10 실연결과 실adapter는 이번 통과로 옮기지 않았다.

## held

DEC-02 실데이터 silent done, DEC-03 공개 한도와 신원, Free N, 실adapter, A2A 현행 검토, WAL 또는 backup 삭제, 고의 stale epoch를 유지한다.
공개 배포와 전체 제품 수락은 이번 좁은 QA의 통과 조건이 아니다.

## 후속

coor가 고정 SHA 독립 리뷰와 main 통합을 이어간다. 원래 held가 남아 있는 동안 전체 제품 수락은 보류다.
이 경계에서 새 high는 없다. 제품 코드의 추가 수정은 이 QA가 요구하지 않는다.
