---
title: "SAR-PUBLISH-001-D04 — D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-PUBLISH-001-D04]
summary: "D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)"
attempt: b0cfa862fa4241b59e1971229fd3e6ef
base: 179bdfe8d1b0b9a70f8e45c979bd19206b656c86
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-D04 — D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: D04 문서 수정과 커밋 / 제품 코드 수정, 운영 서버·운영 env·live Tunnel·Cloud Console·실계정 사용, 정책 문안 변경, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

main 179bdfe8: SAR-PUBLISH-001-DEV(캘린더 제거·/privacy·/terms) 병합, 독립 리뷰 수락, tester QA 62/62 PASS(`.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-TESTER-test/report.md`), designer 시각 검수 PASS(medium 1, low 3, 문서 D-01: `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`). 운영 배포 전에 시각 지적을 고친다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `179bdfe8d1b0b9a70f8e45c979bd19206b656c86`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`(V-01~V-04, D-01)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [ ] D-01: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`의 정책 페이지 하단 배치도 줄을 `개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로`로 고친다. 첫페이지 하단에는 '홈'이 없다고 한 줄 적는다.
- [ ] V-01~V-03 수정 방향(본문 링크 inline, 하단 aria-current 굵게, 설정 링크 gap)이 D04 화면 기준과 다르면 D04에 한 줄로 반영한다. 같으면 바꾸지 않는다.

소유 파일: D04, 이 인박스.

## 완료 기준과 검증

- D04가 실제 화면(정책 하단 '홈' 위치)과 맞는다. front matter는 `deliverables.py --stamp --task SAR-PUBLISH-001-D04`로 갱신한다.
- 기록 커밋 뒤 `lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 179bdfe8d1b0b9a70f8e45c979bd19206b656c86` ERROR 0(스크립트 경로는 dev 지시서와 같다). `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(문서 변경).

## 갱신할 산출물

D04(service-member-ui).

## 기대 산출물

갱신된 D04.

## 제약·협업·후속

dev의 V-01~V-03 수정(SAR-PUBLISH-001-FIX)과 병렬이다. dev 수정 뒤 시각 재확인은 coor가 별도로 정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-D04-packet.json` / SHA `179bdfe8d1b0b9a70f8e45c979bd19206b656c86` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 77 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 184 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 136 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
나머지 참조 후보 4건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
