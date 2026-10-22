---
title: SAR-PUBLISH-002-QA — 운영 KO/EN 정책·첫페이지 키보드·200%·언어 전환 읽기 전용 QA
status: draft
updated: 2026-10-22
owner: tester
tasks: [SAR-PUBLISH-002-QA]
summary: 운영 KO/EN 정책·첫페이지 키보드·200%·언어 전환 읽기 전용 QA
attempt: ff50f2d4c18945c38b98415162fba94d
base: 3e04e6984f25bcb50e605ae69f860a33f6f8e4a4
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-002-QA — 운영 KO/EN 정책·첫페이지 키보드·200%·언어 전환 읽기 전용 QA

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 운영 `https://link.knowslog.com/`·`/privacy`·`/terms`는 읽기·스크롤·키보드 포커스 이동·언어 전환 링크 이동만 한다(쿠키 외 상태 변경 없음). 로그인 회원 저장 언어는 격리 fixture(`TestPublishBrowserFixture` 또는 loopback)로 확인한다. / 추가 승인 필요: 실계정 로그인, 폼 제출, 운영 데이터 변경, 제품 코드 변경, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
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

- [ ] KO/EN 각각 `/`, `/privacy`, `/terms`: 언어 전환 링크(`/lang/en?back=…`, `/lang/ko?back=…`)가 같은 페이지로 돌아오고 문서 언어가 바뀌는지.
- [ ] 키보드 Tab 순회: 하단 링크(개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로)와 `#start` 정책 링크에 포커스가 보이고 순서가 시각 순서와 같은지. 현재 페이지 링크의 `aria-current` 표시.
- [ ] 200% 확대(또는 640px 상당 폭): 가로 스크롤·잘림·겹침 없음. 문장 속 링크 줄 간격이 정상인지(V-01 수정 e6ff5a42 회귀 확인).
- [ ] 로그인 회원 저장 언어: 격리 fixture에서 EN 저장 회원이 `/privacy`를 영어로 보는지.
- [ ] Google 로그인 시작(`/` 로그인 버튼의 이동 대상)이 calendar 범위를 요청하지 않는지 URL의 `scope` 값만 읽어 확인한다(동의 화면에서 계정 선택·진행 금지).

소유 파일: `qa-reports/SAR-PUBLISH-002-QA-tester-test/`, 시나리오 추가 시 `docs/evaluations/scenarios/` 새 파일, 이 인박스.

## 완료 기준과 검증

- 보고 `qa-reports/SAR-PUBLISH-002-QA-tester-test/report.md`: 항목별 PASS/FAIL·근거 한 줄, 실패는 재현·화면 위치. 캡처는 `raw/`(gitignore) 또는 필요한 최소만.
- 결함은 고치지 말고 보고한다. 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.2/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from 3e04e698` ERROR 0. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

하단 링크·포커스·확대 판정은 WCAG 2.2 AA 기준(포커스 보임, 1.4.10 reflow)으로 한다.

## 갱신할 산출물

없음.

## 기대 산출물

QA 보고서.

## 제약·협업·후속

designer SAR-PUBLISH-002-DOCS와 병렬. 웹 실행기는 Orca·Playwright·ego 중 선택(scenarios README).

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-002-QA-packet.json` / SHA `3e04e6984f25bcb50e605ae69f860a33f6f8e4a4` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-TESTER-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/scenarios/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 77 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 133 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `internal/relay/policy/privacy.en.html` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/policy/privacy.ko.html` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 7건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
