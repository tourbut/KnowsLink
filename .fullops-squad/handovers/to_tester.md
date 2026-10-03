---
title: SAR-MVP-001-TESTER-FIX — 수정된 기존 MVP 후보의 독립 QA
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: 수정된 기존 MVP 후보의 독립 QA
---

# SAR-MVP-001-TESTER-FIX — 기존 MVP 수정 후보의 독립 검증

복귀 Run run_8ca8bc058ab7, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor다. 새 preamble의 실제 Task/Dispatch/capability만 사용한다. 신규 기능·정책 수치 확정·실제 배포는 시작하지 않는다.

제품 후보 4262d02fdd7b0b57a804d9e550597852950ffeae의 독립 수정 QA다. 제품은 DEV d7e2149와 같으며 UI e238777·기존 QA c59537b를 포함한다. 기준 ref 4262d02, 공통 fullops-common-0.3.2를 적용한다. 현재 상설 tester에서 제품 코드는 수정하지 않는다.

먼저 FULLOPS·공통 README와 세 규칙·project·문서 작성 규칙·contexts/tester·D02 SAR-MVP·protocol C1–C5·DEV 수정 실행 기록·기존 QA 보고서·UI 보고서를 읽는다. 기존 tester find/documents-find/context를 같은 범위 탐색 근거로 재사용한다. 수정 구현과 reviewer 원래 targeted.log를 대조한다. 기존 SHA·조건·held·실패 기록은 보존한다.

실제 Postgres/HTTP에서 deliver:human의 agent pull·persist·ACK·claim 차단과 owner gate·agent delivery·approval/result·current-auth·철회·TTL·lease·CSRF 관련 회귀를 검증한다. H 이외 direct deliver:human의 신규 403 처리와 frozen 프로토콜의 정합성은 독립 reviewer도 확인한다. 실제 DB 경계 테스트가 실행됨을 증명하며 환경 실패를 결함 재현이나 성공으로 바꾸지 않는다. 변경 의존성이 넓으면 해당 QA 범위를 넓힌다. 기존 QA의 변화 없는 동작·8개 시각 PNG는 실제 제품 파일/설정 동일성 확인 후 원래 SHA·조건으로 재사용한다. 전체 장시간 QA·불필요한 캡처는 반복하지 않는다. 원래 held와 새 실행 PASS를 구분한다.

소유는 docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md와 같은 키 test 증거 폴더·필요 scenarios·자기 contexts·현재 inbox/logs다. PLANS·board·제품·DEV/기획/원래 QA 증거는 수정하지 않는다. 명령 자신의 종료코드를 남긴다. make lint와 영향 자동/실제 DB 회귀·strict·공백 및 커밋 후 FullOps lint를 수행한다. 보고서에 대상 SHA·수락 항목·pass/fail/held·범위·증거 재사용 정합성과 후속을 적는다. 실패/결함은 coor에 알리고 제품 코드를 고치지 않는다. work.py finish·소유파일 커밋·새 worker_done을 수행한다. 본문에 `[완료] SAR-MVP-001-TESTER-FIX | SHA <보고 커밋> | 제품 4262d02 | ...`를 포함한다. 독립 QA 완료는 제품 최종 수락이나 공개 배포 완료가 아니다.

## 완료 보고

상태: 독립 QA를 마쳤다. 제품 `4262d02fdd7b0b57a804d9e550597852950ffeae`의 신규 human transport 차단과 관련 회귀는 통과했다. legacy claim의 authorize·result 수락은 미해결 high이며 제품 수락을 차단한다. 병합 승인과 공개 배포 완료는 아니다. 보고 커밋 SHA는 worker_done 본문에 적는다.

변경 이유: 지시서가 고정 제품 4262d02의 독립 QA를 요구했다. 실제 Postgres와 relay HTTP에서 H 이외 `deliver:human`의 pull·persist·ACK·claim 차단과 owner gate·agent delivery·approval/result·current-auth·철회·TTL·lease·CSRF 회귀를 실행했다. 제품 코드, PLANS, board, 원래 QA 증거는 수정하지 않았다.

지시와 달라진 판단: coor 메시지 msg_ffe2067ae786의 legacy claim high를 이번 판정에 포함했다. 전체 QA는 반복하지 않았다. probe의 첫 종료코드 1은 tester SQL 중단이라 제품 실패로 세지 않았다. 세 번째 probe의 종료코드 1이 제품 판정이다.

수락 기준별 결과: 신규 direct·stored human 차단과 이번 회귀 17항은 통과다. FIX-legacy-claim-authorize-result는 실패다. authorize 200, `relay.result` 200, completion `denied`, gate 0이다. `make lint`, `make test`, `make verify-mvp`의 종료코드는 0이다. `make test`는 integration tag를 포함하지 않는다. 실제 DB 증거는 `make verify-mvp`와 HTTP probe다.

재사용: 승인 화면 템플릿과 설정·adapter·cmd·db는 원래와 같다. PNG 8개 blob은 `c59537b` 7개와 `e238777` 1개와 같다. 새 캡처는 없다. QA-03 HTTP 7종은 `protocol.go`가 같아 `c59537b` 조건으로 남긴다. QA-11 24시간과 QA-10 실연결은 이번 통과로 옮기지 않았다.

held: DEC-02, DEC-03, Free N, 실adapter, A2A 현행 검토, WAL 또는 backup 삭제, 고의 stale epoch를 유지한다. designer 시각 판정은 `e238777` 기록이다. 이번 실행의 시각 통과로 바꾸지 않는다.

lint: 기준 ref `4262d02`. 보고 커밋 `bcc8d34009f00e749975546f66e162b72e8b2ba2`의 FullOps lint 종료코드는 0이다. ERROR 0, WARNING 0, 실행 불가 0. product-lint는 passed다. JSON은 증거 폴더의 `fullops-lint.json`이다.

산출물: [QA 보고서](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md), [시나리오](../docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md), [증거](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX-test/), [컨텍스트](../contexts/tester.md).

후속: dev가 legacy claim high를 수정한다. 새 제품 SHA에서 그 항목을 다시 검증한다. coor는 high가 열린 동안 제품 수락을 차단한다. 원래 held는 유지한다.
