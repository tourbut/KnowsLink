---
title: "SAR-PUBLISH-002-DOCS — 공개 뒤 D04 정책 하단 '홈'·D11 정책 절 갱신"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-PUBLISH-002-DOCS]
summary: "공개 뒤 D04 정책 하단 '홈'·D11 정책 절 갱신"
attempt: 4a0ef7137c814ca49c9d82ccd36cda9c
base: 3e04e6984f25bcb50e605ae69f860a33f6f8e4a4
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-002-DOCS — 공개 뒤 D04 정책 하단 '홈'·D11 정책 절 갱신

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: D04·D11 문서만 고친다. 제품 코드·정책 HTML·다른 정본은 바꾸지 않는다. / 추가 승인 필요: 제품 코드 변경, 정책 문안 변경, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

2026-10-22 SAR-PUBLISH-001이 끝났다. 운영은 main 4fdd149b 제품 코드다. 서버 캘린더 연결을 제거했고(첫 기동 때 연결 0), `/privacy`·`/terms`(KO/EN, 백업 14일 문장)가 비로그인 200으로 공개됐다. Google OAuth 앱 knowslink-auth는 calendar.freebusy scope 삭제, Calendar API 사용 중지, 브랜딩 URL 입력 뒤 "프로덕션 단계"로 게시됐다(D022). 남은 후속: designer 시각 검수 D-01(D04 배치도 누락), D11 사용 문서의 'DEV 후보' 문구, tester가 보지 못한 키보드·200%·서버 경유 언어 화면(`.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`, `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-TESTER-test/report.md` 59행).

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `3e04e698`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`(D-01·V-01–V-04)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-TESTER-test/report.md`(검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [ ] D04 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` 정책 페이지 하단 배치도 줄을 `개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로`로 고치고, 첫페이지 하단에는 '홈'이 없다고 한 줄 적는다(D-01, `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md` 39–53행).
- [ ] D11 `.fullops-squad/docs/operations/user-guide.md` '개인정보처리방침과 서비스 약관 — SAR-PUBLISH-001 DEV 후보' 절을 운영 공개 상태로 고친다: 제목에서 'DEV 후보'를 빼고, '이 절은 DEV 후보의 사용법이다…' 문장을 '2026-10-22부터 공개 주소에서 열린다'로 바꾼다. '법률 검토 전 초안' 문장은 P-05가 열려 있으므로 유지한다.
- [ ] 같은 파일의 캘린더 연결 삭제 안내 문단을 실제 결과에 맞춘다: 운영 첫 기동 때 남은 연결이 0이었다(`.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`). 안내 대상이 없다는 사실만 적고 절차 문장은 줄인다.

소유 파일: D04, D11, 이 인박스.

## 완료 기준과 검증

- 위 세 항목 반영. 다른 파일 변경 없음.
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.2/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 3e04e698` ERROR 0. 산출물 front matter는 `deliverables.py --stamp --task SAR-PUBLISH-002-DOCS`. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(문서).

## 갱신할 산출물

D04, D11.

## 기대 산출물

D04·D11 갱신 커밋.

## 제약·협업·후속

tester SAR-PUBLISH-002-QA와 병렬. 파일이 겹치지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-002-DOCS-packet.json` / SHA `3e04e6984f25bcb50e605ae69f860a33f6f8e4a4` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 16, 26 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 119, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 30, 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14, 14, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 133, 137 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_update, document_read) · 줄 113 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 23, 23, 44 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
나머지 참조 후보 16건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
