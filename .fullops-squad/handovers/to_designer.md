---
title: SAR-CONNECT-SIMPLE-001-ART — 연결 단순화 화면 직접 시각 검수(CS-AC-17)
status: draft
updated: 2026-10-27
owner: designer
tasks: [SAR-CONNECT-SIMPLE-001-ART]
summary: 연결 단순화 화면 직접 시각 검수(CS-AC-17)
attempt: c69d175bf4c24aa5b453f4782b3bf45a
base: 3f1da403086e4449a7b3df214c441cd1b5972db6
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-ART — 연결 단순화 화면 직접 시각 검수(CS-AC-17)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: designer 초안(SAR-CONNECT-SIMPLE-001 tester 초안의 ART 절) → coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: fullops/designer에서 시각 검수 보고서·판정 캡처 커밋, main 3f1da403의 격리 로컬 relay와 합성 회원 fixture 실행, 커밋·push.
- 추가 승인이 필요한 행위: 운영 relay·운영 agent 사용, 실제 발송, 제품 코드 수정(결함은 보고와 수정 지시안만), D02 제품 규칙 변경.
- 상위 의도: 없음. D02 SAR-CONNECT-SIMPLE-001 CS-AC-17.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

SAR-CONNECT-SIMPLE-001-DEV(0fff4ca0)와 리뷰 후속 FIX-DEV(987152d4)가 main 3f1da403에 통합됐다. 독립 코드 리뷰는 critical/high 0이다. UI 변경(M·W)이 있어 designer 직접 시각 검수를 배정한다. tester는 같은 SHA에서 SAR-CONNECT-SIMPLE-001-QA를 병행한다.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`, `testing.md`, `security.md`. 디자인 정본: `DESIGN.md`와 D04 `docs/design-docs/mockups/service-member-ui.md`.
- 기준: main `3f1da403086e4449a7b3df214c441cd1b5972db6`. ID 정의: `M-`·`W-`·`CS-AC-` → [D02 SAR-CONNECT-SIMPLE-001](../docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md). `D01~D13` → `docs/deliverables/README.md`.
- 예외: 없음.

## 먼저 읽을 문서

1. [D04 연결 단순화 절](../docs/design-docs/mockups/service-member-ui.md#연결-단순화와-내-에이전트-정리--sar-connect-simple-001)
2. [D02 SAR-CONNECT-SIMPLE-001](../docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md)의 M·W·E 표
3. DEV 보고 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md`(fixture 띄우는 방법)

## 해야 할 일과 파일 소유권

- [ ] '내 에이전트'의 진행 중인 연결 절('이 시도 정리' 포함)과 이름을 정할 에이전트 절.
- [ ] 연결 링크 화면: 대기·동의 직후·완료·만료·취소, 완료 링크 재방문(W-01).
- [ ] 연결 중 Google 실패 화면(W-02 등 fixture로 가능한 원인).
- [ ] 위 화면을 KO·EN × 1440·390·320px로 보고 D04 문구 표와 '자기 비평과 수락' 기준으로 판정한다. 결함은 위치·재현·수정 지시안과 심각도를 적는다.

파일 소유권: designer는 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/`만 쓴다. 제품 코드·D02·D04는 고치지 않는다(필요하면 남은 일에 적는다).

## 완료 기준과 검증

- 보고 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`에 화면×언어×폭마다 PASS/FAIL/not_run과 근거. front matter는 `deliverables.py --stamp`로 쓴다.
- 캡처는 판정 항목에 연결된 것만 커밋한다. 로그·임시 파일은 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART/`.
- FullOps lint `--from 3f1da403086e4449a7b3df214c441cd1b5972db6` ERROR 0.
- fixture는 tester와 다른 포트·DB로 띄운다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 토큰 정본: `DESIGN.md`, D04 위 절.
- 공용 컴포넌트·예외: 검수만 한다. 예외가 보이면 근거와 함께 보고.
- 디자인 lint: 해당 없음(코드 변경 없음).
- 테마 전환: 제품에 테마 전환이 있으면 라이트·다크 각 1폭 확인, 없으면 not_run 사유.

## 갱신할 산출물

없음(QA 보고서만). D04 변경이 필요해 보이면 남은 일에 적는다.

## 기대 산출물

`docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`와 판정 캡처.

## 제약·협업·후속

운영 relay·운영 agent는 쓰지 않는다. 하위 위임은 필요 없다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-ART-packet.json` / SHA `3f1da403086e4449a7b3df214c441cd1b5972db6` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update, document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 25, 60 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 16, 26 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 15, 174 · inferred
- `.fullops-squad/docs/planning/SAR-AUTOSETUP-001-DEV-brief.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 108 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 23, 44 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
나머지 참조 후보 5건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 판정 요약 / 결함 / 남은 일. worker_done 본문 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-ART | 브랜치 fullops/designer | SHA <HEAD>` 한 줄로 쓴다.
