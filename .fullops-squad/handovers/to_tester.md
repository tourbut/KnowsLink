---
title: SAR-MVP-001-TESTER-FIX — RF-01 수정 좁은 독립 QA
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: legacy claim 경계의 새 후보를 독립 재검증한다
---

# SAR-MVP-001-TESTER-FIX — 78b1d92 후속

고정 제품 78b1d92c8aa626245d3349ffaf7367d28f1dd3ef의 RF-01 수정만 좁게 독립 재검증한다. 이전 4262d02 QA의 report/test와 original31pass8held·UI PNG 동일성 증거는 원래 SHA/조건으로 재사용하고 덮어쓰지 않는다. 새 기록은 SAR-MVP-001-TESTER-FINAL.md 및 SAR-MVP-001-TESTER-FINAL-test/에 작성한다. inbox key는 같은 SAR-MVP-001-TESTER-FIX 후속이다.

먼저 FULLOPS·공통 README와 세 규칙·project·문서 작성 규칙·기존 TESTER-FIX 보고·실패 REVIEW-FIX legacy 증거·DEV 후속 phase/logs 및 protocol C1을 읽는다. 이전 같은 key Jev find/context를 재사용한다. 실제 a6a10c7의 serialized human/unrouted fixture와 parentRouting의 차단·authorize/result·H/R 부모/owner consume 경계를 새 SHA의 실제 Postgres/HTTP에서 검증한다. 정상 새 agent·owner gate·approval/result·현재 권한 영향 회귀를 필요한 범위만 검증한다. 제품 파일과 PLANS/board는 수정하지 않는다. 새 전체 MVP 반복·새 기능·실제 배포는 시작하지 않는다.

관련 QA 보고·시나리오·테스트 증거·context/inbox/logs만 소유한다. fullops lint ERROR0/product-lint passed·strict·공백·실제 종료코드를 보존한다. 원래 held를 pass로 바꾸지 않고 신원/벤더·정책·처리량·WAL·배포 한계를 유지한다. 같은 key finish 중복 차단 시 실패 근거와 별도 후속 전문을 logs에 보존하고 inbox를 비운다. worker_done에 [완료] SAR-MVP-001-TESTER-FIX | SHA <보고 커밋> | 제품78b1d92 | 판정·남은 held를 포함한다. Run run_8ca8bc058ab7/coor term_9afa8217-862c-404d-9a43-2122427113fc로 새 preamble을 사용한다. 독립 리뷰·main 통합은 coordinator가 후속한다.

## 완료 보고

상태: RF-01 좁은 독립 QA를 마쳤다. 제품 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`에서 legacy human claim과 경로 미기록 claim의 authorize, H, R, consume은 `403 sender_not_allowed`다. 새 agent, owner gate, approval/result, current-auth는 통과했다. 원래 held는 유지한다. 전체 MVP 수락과 공개 배포 완료는 아니다. 보고 커밋 SHA는 worker_done 본문에 적는다.

변경 이유: 지시서가 같은 과제의 후속으로 `78b1d92`만 좁게 재검증하라고 했다. 실제 Postgres와 relay HTTP에 `a6a10c7` fixture를 로드했다. 제품 코드, PLANS, board, `4262d02` report/test는 수정하지 않았다.

지시와 달라진 판단: fixture TTL이 QA 시각보다 이전이었다. `Receipt.exp`와 `Gate.Exp`만 4분 뒤로 옮겨 sweep가 claim을 지우지 않게 했다. `Deliver`와 claim token은 fixture 값이다. probe 1–3회 종료코드 1은 tester 형태 비교다. 4회 종료코드 0이 판정이다.

수락 기준별 결과: legacy 경계와 이번 회귀 9항은 통과다. `go test` legacy fixture 종료코드는 0이다. 그 명령은 HTTP 증거를 대신하지 않는다. `make verify-mvp`는 실행하지 않았다. 설정과 migration은 `4262d02`와 같다.

재사용: `4262d02` 보고서와 test 증거는 유지했다. PNG 8개는 `c59537b`와 `e238777` 조건으로 남긴다. 새 캡처는 없다.

held: DEC-02, DEC-03, Free N, 실adapter, A2A 현행 검토, WAL 또는 backup 삭제, 고의 stale epoch를 유지한다. designer 시각 판정은 `e238777`이다.

lint: 기준 ref `4262d02`. 보고 커밋 `f7ff0b26be890bd73852dc0e8498e66f1475ed0a`의 FullOps lint 종료코드는 0이다. ERROR 0, WARNING 1, 실행 불가 0. product-lint는 passed다. WARNING은 `SIZE-001` `internal/relay/integration_test.go` 478줄이다. 제품 파일이므로 수정하지 않았다. JSON은 증거 폴더의 `fullops-lint.json`이다.

산출물: [QA 보고서](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md), [시나리오](../docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md), [증거](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL-test/), [컨텍스트](../contexts/tester.md).

후속: coor가 독립 리뷰와 main 통합을 이어간다. 원래 held가 열린 동안 전체 제품 수락은 보류다. 같은 key `work.py finish`가 중복으로 실패하면 그 종료코드와 이 전문을 logs의 별도 후속 절에 남기고 inbox를 비운다.
