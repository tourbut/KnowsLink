---
title: "SAR-CODEX-RECV-001-DOCS — D11 사용자 설명서에 Codex 세션 연결 안내(확인 코드 승인, YOLO 불가, 지원 수준) 추가"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-CODEX-RECV-001-DOCS]
summary: "D11 사용자 설명서에 Codex 세션 연결 안내(확인 코드 승인, YOLO 불가, 지원 수준) 추가"
attempt: 173c014435144577a8eaa7a5cf99ab00
base: b102960514f5c3ef782469eee9ea256795b5472b
subagent_level: standard
test_level: lite
---

# SAR-CODEX-RECV-001-DOCS — D11 사용자 설명서에 Codex 세션 연결 안내(확인 코드 승인, YOLO 불가, 지원 수준) 추가

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/designer의 D11 문서 수정·커밋·push. 추가 승인 필요: 제품 코드·웹 화면 수정, 운영 배포, 실계정 사용
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

Codex 세션 수신(C-01, 결정 D027)이 main `b1029605`에 병합됐다(dev bd1f3070, tester QA critical/high 0, D-CR-01 수정). 사용자 안내가 아직 없다. 사실 근거: `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`(V02 판정·관찰 O-1·O-2), `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md`, `adapters/README.md`.

## 적용 기준과 예외

문서 작성 `.fullops-squad/docs/agents/document-writing.md`. 기준 main `b102960514f5c3ef782469eee9ea256795b5472b`. D11 정본 `.fullops-squad/docs/operations/user-guide.md`. 제품 규칙 `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` SP-12(근거 없는 지원 표기 금지). 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/operations/user-guide.md`(기존 Claude Code 세션 연결 안내 절 형식)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`
- `adapters/README.md`(설치·지원 표)
- `.fullops-squad/contexts/designer.md`

## 해야 할 일과 파일 소유권

- [ ] `.fullops-squad/docs/operations/user-guide.md`에 Codex 연결 절 추가: 설치, 연결(`knowslink_session on`·take·off·status), 승인 창마다 표시되는 4자리 확인 코드 입력, 거절은 Esc 취소, Codex 도구 실행 승인 창과 KnowsLink 승인 창이 따로 뜬다는 점
- [ ] 제한 명시: YOLO(전체 자동) 세션에서는 연결·발송 불가, 도착 알림 없음(다음에 받기 요청 때 본문 표시), 한 agent는 한 세션만 받음
- [ ] front matter stamp(`deliverables.py --stamp --task SAR-CODEX-RECV-001-DOCS`), contexts/designer.md 3줄

파일 소유권: D11과 designer 인박스·contexts만.

## 완료 기준과 검증

- 안내가 DEV/tester 보고의 실측 범위를 넘지 않는다(SP-12). 미검증(하위 에이전트, 4096바이트 대화상자 전체 표시)은 안내에 지원으로 적지 않는다.
- lint: designer 워크트리 `python <FullOps 1.8.2 scripts>/lint.py --repo . --from b1029605` ERROR 0. 문서 과제라 QA report.md는 만들지 않는다.

### UI 작업의 디자인 기준과 검증

해당 없음. 문서만 바꾼다.

## 갱신할 산출물

D11(`docs/operations/user-guide.md`). 그 밖 없음.

## 기대 산출물

- `.fullops-squad/docs/operations/user-guide.md`

## 제약·협업·후속

하위 위임 없음. 남은 일이 없으면 "없음"으로 적는다(tester 후속 없음).

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CODEX-RECV-001-DOCS-packet.json` / SHA `b102960514f5c3ef782469eee9ea256795b5472b` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 119, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-CODEX-RECV-001-DEV.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 137 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 74, 81 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 48 · inferred
- `docs/README.md` (document_read) · 줄 38 · inferred
나머지 참조 후보 0건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
