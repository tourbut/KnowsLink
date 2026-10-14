---
title: SAR-REFACTOR-002-TESTER — Astra 리팩터링 고정 제품 후보 독립 회귀·성능 검증
status: draft
updated: 2026-10-14
owner: tester
tasks: [SAR-REFACTOR-002-TESTER]
summary: Astra 리팩터링 고정 제품 후보 독립 회귀·성능 검증
attempt: 23e80de7e6e84386b7b3d54b114b3f26
base: bcb1dd6629b508a30ec39a5d3141c8ce0a558597
subagent_level: standard
test_level: lite
---

# SAR-REFACTOR-002-TESTER — 고정 후보 독립 회귀·성능 검증

- From/To: coor/tester. 상태: ready. Purpose: review.
- Task key: SAR-REFACTOR-002-TESTER. Test level: lite. Subagent level: standard.
- Agent/model/effort: codex/gpt-6-astra/high. 사용자 최신 Astra high 지정.
- 제품 후보: bcb1dd6629b508a30ec39a5d3141c8ce0a558597. 비교 제품 기준: 8f9370bb62c9133c73e10ea2307b4cd696e10f7e. 준비/QA 커밋이 추가돼도 제품 입력 동일성을 직접 확인한다.
- 담당 repo 0b08ec4c-9e3a-4613-8197-5a835e545335. C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester.
- 복귀 coor C:/Users/shin/orca/workspaces/KnowsLink/coor, terminal term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061, run run_5534969da2d0. 실제 task/dispatch는 시작 메시지.
- 독립 코드 리뷰는 별도 Astra 세션 SAR-REFACTOR-002-CODE-REVIEW가 수행한다. tester는 독립 동작 QA를 수행하고 제품 코드를 수정하지 않는다.

## 적용 기준과 예외

AGENTS.md, FULLOPS.md, fullops-common-0.3.3 README·coding-style·testing·security, project.md, document-writing.md, fullops-test를 읽는다. 기존 HP-01–16·HTTP/MCP·서명·권한·TTL·수신/ACK·writer guard 기준을 유지한다. 직접 재현과 실패 경계를 새로 검사하며 원 개발자의 PASS를 독립 QA로 바꿔 적지 않는다. 별도 reviewer가 code review를 담당하므로 같은 snapshot review 기록을 만들 필요는 없다.

## 먼저 읽을 문서와 파일

- docs/system-design.md
- .fullops-squad/project.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/adapter-report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/remote-report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/gate/remote.py
- internal/relay/store_retention_test.go
- internal/relay/store_retention_integration_test.go
- internal/relay/remote_events_test.go
- adapters/src/inbox-files.test.ts

## 해야 할 일과 파일 소유권

- [ ] 정확한 제품 입력과 실제 환경을 기록한다. 테스트·QA만 소유하며 DEV 제품 파일·인박스·coor PLANS/board를 수정하지 않는다.
- [ ] 보관 sweep의 만료 직전/정각·철회·delivered·claim·24h 제거 경계를 독립 확인한다. 보관 사본의 조회 권한과 활성 transport 인가가 약화되지 않는지도 검사한다.
- [ ] Inbox housekeeping EACCES 실패 시 unread 원본/본문이 남고 재시도에 수신되는지 확인한다. concurrent stat ENOENT만 허용하고 다른 오류는 전달되는지, 중복/읽음 marker/빈 inbox 경계도 확인한다.
- [ ] callback 마지막 lease 만료 뒤 새 도착이 즉시 다시 wake되고 새 event ID/첫 시도·3회 한도·활성 lease·성공 뒤 idle·철회/만료가 유지되는지 검사한다. 단일 claim의 per-agent pending 재사용이 요청 간 캐시가 아닌지 확인한다.
- [ ] 관련 Go/TS 새 회귀를 직접 실행하고 최소한 한 가지 이상의 독립적인 경계 입력을 추가/assert로 검증한다. 기존 tests를 그대로 재실행한 결과와 새 독립 checks를 구분한다. race/실제 DB 검증은 구현자 exact Linux gate의 범위와 신규 실행을 구분해 필요한 짧은 확인을 수행한다.
- [ ] before/after benchmark fixture·환경·범위를 검토하고 짧게 재측정한다. 개선 수치를 end-to-end throughput으로 과장하지 않는다. 큰 sweep 이득과 전체 JSONB/row lock 한계를 구분한다.
- [ ] 최종 보고에 실제 SHA·명령·종료코드·PASS/실패/미실행·독립 근거·발견 사항·남은 위험을 기록한다. 필수 unexplained 실패는 완료 수락으로 표시하지 말고 coor에 재현과 함께 즉시 전달한다.
- [ ] packet-outcomes와 uncertainty_review, 상세 phase, 전문 완료 보고를 작성하고 work.py finish 후 커밋한다. worker_done 본문에 SHA: <40자리>를 명시한다.

## 검증과 완료 조건

합성·격리 fixture만 사용한다. 실계정·외부 메시지·운영 DB·배포는 하지 않는다. Windows에서 race/도구가 미지원이면 기존 안전 SSH 패턴으로 고유 Linux clone·합성 loopback DB를 사용하되 known_hosts/RejectPolicy·HEAD/clean·출력 가림을 유지한다. DEV gate helper와 기록은 재사용할 수 있지만 실행 SHA/owner를 구분한다. 레포 밖 임시 자료는 보존하고 경로를 적는다. 삭제·이동 정리를 실행하거나 안전 guard를 완화/우회하지 않는다. 파일 편집 중 guard가 문자열 오탐으로 차단해도 원문을 보존해 coor에 보고한다.

DEV 필수 lint/test/build는 DEV가 최종 안정 후보에서 실행한다. tester는 지정 경계 QA를 독립 실행하며 같은 전체 gate를 이유 없이 반복하지 않는다. UI 외형 변경이 없어 새 시각 캡처는 해당 없음이다. 설명되지 않는 실패·critical/high가 없고 독립 경계 검증이 완료돼야 QA를 수락한다. 최종 main 통합과 push는 coor가 코드 리뷰·QA·문서 delta를 함께 확인한 뒤 수행한다.

## 기대 산출물

- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md 및 실행 근거
- .fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-TESTER.md
- 필요한 최소 독립 테스트 또는 assert 기반 검사, packet-outcomes, 전문 archive·빈 인박스

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-REFACTOR-002-TESTER-packet.json` / SHA `f5f883402a283e6884cd81d12f0226b4da2c6491` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/adapter-report.md` (document_read) · 줄 2, 6 · inferred · 필수 · {'relevant': 0.77, 'evidence': 0.88, 'contradicts': 0.23, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md` (document_read) · 줄 6 · inferred · 필수 · {'relevant': 0.68, 'evidence': 0.81, 'contradicts': 0.19, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/remote-report.md` (document_read) · 줄 6 · inferred · 필수 · {'relevant': 0.83, 'evidence': 0.87, 'contradicts': 0.3, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 17, 17, 22, 32, 33, 34, 35, 60, 60, 61, 61 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.77, 'evidence': 0.87, 'contradicts': 0.19, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `adapters/src/inbox.ts` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `docs/system-design.md` (document_read) · 줄 21 · inferred · 필수 · {'relevant': 0.51, 'evidence': 0.68, 'contradicts': 0.13, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_events.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `internal/relay/store.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/PLANS.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 91, 93 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 6, 18, 19, 20, 21, 22, 23, 25, 27, 29, 30, 32, 33, 35, 36, 38, 39, 41, 42, 44, 45, 47, 48, 50, 51, 53, 54, 56, 57, 59, 61, 62, 64, 66, 69, 71, 74, 76, 77, 80, 80, 81, 81 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 104 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 52, 53, 54, 55, 56, 57, 58, 101 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 14, 50, 87 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 14, 64, 82 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 25, 87, 118 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 98 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 61, 67 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 14, 20, 91, 174, 177, 193, 195, 199, 201, 203 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 39, 41 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.2-review/report.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-COPY-REVIEW-review/report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 6, 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-TESTER.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 2, 6, 10, 28, 53, 65, 73, 74, 75, 76, 77, 78 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 2, 6, 10, 19, 24, 30, 35, 64, 71, 72, 73, 74, 75, 76, 77, 78, 79, 89, 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 2, 6, 10, 22, 100, 101 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/report.md` (document_read) · 줄 6, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-QA-DOCS-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 6, 10, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-test/report.md` (document_read) · 줄 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 73, 83 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-DOC-review/report.md` (document_read) · 줄 2, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-UI-DELTA-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 136 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 2, 6, 10, 21, 27, 29, 53, 59, 61, 62, 66, 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 2, 6, 10, 20, 29, 34, 36, 68, 70, 71, 72, 73, 74, 75, 108, 109, 110, 111, 112, 113, 114, 117 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 2, 6, 10, 18, 30, 51, 52, 53, 55, 56, 57, 91, 92, 93, 94, 95, 96, 97 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-UI.md` (document_read) · 줄 21, 31, 32, 33, 34, 35, 36, 37, 41, 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 23, 30, 93, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 129, 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-review/report.md` (document_read) · 줄 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 24, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 123, 124 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 19, 103, 104, 105, 106, 107, 108, 109, 115 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/report.md` (document_read) · 줄 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 23, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md` (document_read) · 줄 2, 6, 10, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md` (document_read) · 줄 2, 6, 10, 72, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PREP-002-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FINAL-RECORDS-review/report.md` (document_read) · 줄 14, 26, 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 39, 105, 112 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER-test/probe-failures.md` (document_read) · 줄 2, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 2, 6, 10, 27, 49, 150, 151 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/scenario.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 22, 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-QA-RECORD-REVIEW-review/report.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 7, 39, 85, 89, 90, 91, 92, 137, 143, 162 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 119 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 2, 6, 10, 33, 72 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW-review/report.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 6, 10, 29 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 117 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER-test/probe-failures.md` (document_read) · 줄 2, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 2, 6, 10, 31, 38, 42, 69, 102, 109 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-CODE-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md` (document_read) · 줄 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-FINAL-review/report.md` (document_read) · 줄 6, 6, 12, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-FINAL-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-RESULT-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md` (document_read) · 줄 6, 10, 19, 19, 124 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-RESULT-review/report.md` (document_read) · 줄 1 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md` (document_read) · 줄 6, 6, 10, 10, 24, 66, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CLEANUP/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 58, 68 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/report.md` (document_read) · 줄 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 15, 38, 47, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 2, 6, 10, 17, 31, 32, 34, 36, 37, 39, 83, 92, 98, 101, 102, 103, 112 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 20, 41, 43, 53 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 6, 28, 59, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/README.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 2, 6, 10, 21 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 2, 6, 10, 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-002-TESTER.md` (document_read) · 줄 2, 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-HUMAN-PEER-001-TESTER.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 2, 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 15 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 2, 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 2, 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 6, 10, 12, 13, 18 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 2, 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 2, 6, 10, 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 6, 10, 12, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-TESTER.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-TESTER.md` (document_read) · 줄 2, 6, 10, 18, 28, 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 151 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV-REVIEW.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 17, 136 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT.md` (document_read) · 줄 14, 74 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 75, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 23, 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 91 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 89 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 18, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 17, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 7, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 2, 6, 10, 14, 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 36, 42, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 49 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 40, 60, 78, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-DEV.md` (document_read) · 줄 6, 10, 46, 52, 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-OPS.md` (document_read) · 줄 6, 10, 18, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-DEV.md` (document_read) · 줄 6, 10, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 2, 6, 10, 20, 26, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-TESTER.md` (document_read) · 줄 2, 6, 10, 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 36, 37, 39, 47 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 6, 10, 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 23, 93 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 23, 63, 85, 86, 87, 88 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 97 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 131 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 96 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 183, 189, 266, 268 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 205, 278, 292 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 53, 72, 100, 119, 129, 259, 297 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 53, 132, 145 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 34, 43, 51, 57, 59, 62, 66, 70, 83, 91, 95, 99, 107, 109, 112, 116, 120, 131, 139, 140, 144, 150, 154, 157, 161, 165, 176, 183, 192, 194, 197, 201, 205, 213, 217, 224 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 35, 81 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 147 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 41, 71, 75, 81 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 156, 253, 284, 318, 426, 451 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 12, 17, 20, 36, 74, 76, 109, 130, 203, 212, 250, 283 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 37, 49, 56, 57, 70, 74, 77, 81, 85, 124, 132, 140, 144, 147, 151, 155, 180, 191, 199, 209, 212, 216, 220, 236 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 48, 72, 104, 105, 108, 109, 110, 111, 112, 114, 117, 119, 120, 122, 126, 131, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 173, 175, 185, 267 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 61, 170, 171, 174, 175, 176, 177, 178, 180, 183, 185, 188, 190, 196, 197, 198, 199, 200, 201, 231 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 23, 34, 40, 67, 72, 73, 101, 102, 103, 107, 108, 109, 111, 113, 114, 116, 117, 118, 122, 124, 125, 127, 131, 136, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 157, 171, 173, 183, 263 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 9, 10, 13, 17, 25, 45, 57, 66, 95, 100, 101, 103, 106, 107, 108, 114, 123, 125, 126, 127, 133, 135, 140, 142, 159, 185, 186, 266, 311 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 9, 10, 13, 17, 25, 101, 105, 133, 134, 138, 140, 142, 144, 146, 147, 151, 153, 154, 156, 157, 158, 159, 160, 160, 161, 164, 165, 169, 173, 174, 185, 197, 198, 203, 245, 259, 267, 267, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 36, 43, 43, 62, 62, 66, 66, 67, 67, 70, 70, 74, 74, 82, 82, 105, 105, 114, 114, 117, 118, 118, 119, 119, 151, 152, 153, 154, 158, 159, 160, 165, 166, 167, 170, 171, 172, 177, 178, 179, 181, 183, 184, 186, 188, 189, 193, 195, 196, 198, 202, 205, 207, 208, 211, 212, 212, 213, 213, 214, 215, 218, 219, 220, 224, 225, 226, 230, 233, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 260, 266, 274, 289, 291, 298, 299, 301, 302, 308 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 49 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 6, 14, 17, 74, 75, 76, 82, 119, 120, 123, 124, 127, 133, 134, 150, 152, 155, 156, 157, 158, 159, 160, 161, 162, 162, 165, 166, 173, 174, 175, 181, 186, 188, 203, 204, 207, 211 · inferred
- `README.md` (document_read) · 줄 36 · inferred
- `adapters/README.md` (document_read) · 줄 232 · inferred
- `adapters/src/inbox-files.test.ts` (direct_edit, impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.85, 'evidence': 0.83, 'contradicts': 0.19, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `deploy/knowslink/verify.py` (impact_check) · 줄 전체/미확인 · inferred
- `docs/README.md` (document_read) · 줄 37 · inferred
- `docs/development.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/operations.md` (document_read) · 줄 10 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 177 · inferred
- `internal/relay/remote_events_test.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown
- `internal/relay/store_retention_integration_test.go` (direct_edit) · 줄 17 · inferred
- `internal/relay/store_retention_test.go` (direct_edit, impact_check) · 줄 11, 37, 76 · inferred · {'relevant': 0.6, 'evidence': 0.57, 'contradicts': 0.25, 'injection': 0.02, 'decision': 'keep', 'reason': None}
미확인 19건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

실제 독립 검증·제품 입력 SHA·결함/미실행·성능 측정 범위·QA 근거를 적는다.
