---
title: SAR-CONNECT-SIMPLE-001-UI-FIX-DEV — 연결 단순화 화면 검수·QA 후속 수정(V-CS-01–06)
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-UI-FIX-DEV]
summary: 연결 단순화 화면 검수·QA 후속 수정(V-CS-01–06)
attempt: f24c4be2e09f4fa6a0417a2a146d4de2
base: e426f03be373426555dcd1c262407354912483cf
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-UI-FIX-DEV — 연결 단순화 화면 검수·QA 후속 수정(V-CS-01–06)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor(ART·QA 후속) → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: fullops/dev에서 relay 화면·i18n·정렬 코드와 테스트, D04 표의 V-CS-02 문장(designer가 정한 문장 그대로), D10 수정, 격리 로컬 relay·합성 fixture 검증, 커밋·push.
- 추가 승인이 필요한 행위: 운영 배포, 운영 relay 데이터 변경, 실제 발송, D02 제품 규칙 변경, 위 문장 외 문구 결정.
- 상위 의도: 없음. D02 SAR-CONNECT-SIMPLE-001 CS-AC-17 수락 조건(V-CS-01).
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

main e426f03b에 SAR-CONNECT-SIMPLE-001 DEV·리뷰 후속·ART·QA가 통합됐다. designer ART는 CS-AC-17을 V-CS-01(medium) 때문에 FAIL로 판정했고 low 5건을 남겼다. tester QA는 critical/high/medium 0, low 1(V-CS-02와 같은 조사 문제)이다. 결함 표와 수정 지시안은 ART 보고서에 있다.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`, `coding-style.md`, `testing.md`, `security.md`. 디자인 정본 `DESIGN.md`, D04 `docs/design-docs/mockups/service-member-ui.md`.
- 기준 main `e426f03be373`. ID: `V-CS-` → ART 보고서, `CS-AC-`·`W-` → [D02 SAR-CONNECT-SIMPLE-001](../docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md).
- 예외: 없음. 다른 회원 정보 비노출은 유지한다.

## 먼저 읽을 문서

- ART 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`의 결함 표(V-CS-01–06 위치·재현·수정 지시안)
- QA 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md` 결함 1
- D04 연결 단순화 절

## 해야 할 일과 파일 소유권

- [ ] V-CS-01(medium): EN 사전에 `"Google 로그인": "Google sign-in"` 추가. EN refusal 렌더 h1·`<title>`에 한글이 없음을 단언하는 테스트를 `TestConnectFailureMessages`에 추가.
- [ ] V-CS-02(low): `member_pages.go` device outcome approved KO 본문을 ART 보고서의 designer 결정 문장으로 교체(EN 그대로). D04 표의 같은 칸도 같은 문장으로.
- [ ] V-CS-03(low): 새 두 절 버튼에 기존 형식 `aria-label`('이 시도 정리: <클라이언트> · <HH:MM> 시작', '이름 저장: 이름 없는 에이전트 · <구분값>', EN 같은 형식). 보이는 문구는 그대로.
- [ ] V-CS-04(low): `agent_rate` 거절 재방문의 복귀 링크를 '내 에이전트로 돌아가기'(`/home/agents`).
- [ ] V-CS-05(low): device 화면 클라이언트 줄을 `hostNames` 표시 이름으로(모르면 현행 문구). '내 에이전트' 시도 제목과 같은 값.
- [ ] V-CS-06(low): `ownAttempts`를 시작 시각 최신순, 같으면 ID순으로 정렬.

예상 변경: 코드 40–100줄, 테스트 포함.

## 완료 기준과 검증

- 6건 반영, 새·기존 테스트 통과. Linux 격리 `make lint`·`make test`. FullOps lint는 깨끗한 detached worktree에서 `--from e426f03be373` ERROR 0.
- 증거 `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-test/report.md`(front matter는 `deliverables.py --stamp`).
- 수정 확인 범위: designer EN 1폭 재검수(V-CS-01)는 coor가 이 SHA 뒤 배정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 토큰: `DESIGN.md`, D04 위 절. 새 스타일·토큰은 만들지 않는다.
- 공용 컴포넌트: 기존 refusal·device 템플릿과 카드 aria 형식을 재사용한다.
- 디자인 lint: lint.json에 디자인 검사가 없으면 해당 없음.
- 테마 전환: 제품에 없음.

## 갱신할 산출물

D04(V-CS-02 문장 한 칸), D10(필요 시). 그 외 없음.

## 기대 산출물

relay 코드·테스트, D04 한 칸, QA 보고.

## 제약·협업·후속

dev 작업 트리의 SAR-MOTION-INTRO 미추적 파일은 건드리지 않는다. 하위 위임 불필요. 이전 CONNECT 시나리오 기대 갱신은 tester 몫이라 이 과제 범위 밖이다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-packet.json` / SHA `e426f03be373426555dcd1c262407354912483cf` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md` (document_read) · 줄 24, 58, 75 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 20, 24 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 23, 32, 44 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_read) · 줄 87, 124, 143, 170 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/system-design.md` (document_read) · 줄 70, 78, 113 · inferred
- `internal/relay/google.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/human_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/i18n_en_member.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member_agents.go` (direct_edit) · 줄 64, 68, 71, 105, 128, 196, 400, 404 · inferred
- `internal/relay/member_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 25건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 6건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; internal/relay/google.go: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일. worker_done 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-UI-FIX-DEV | 브랜치 fullops/dev | SHA <HEAD>`.

- 브랜치 / SHA: `fullops/dev` / 코드 `659380f1`·`9f279b4f`, 보고 `165ec732`(finish 커밋이 최종 HEAD).
- 변경 이유: ART V-CS-01(medium)·V-CS-02–06(low), QA 결함 1. EN 사전 `Google 로그인`, device approved KO 문장(D04 같은 칸), 두 절 버튼 `aria-label`, `agent_rate` 거절 재방문 복귀 `/home/agents`, device 클라이언트 줄 `clientName`, `ownAttempts` 최신순·ID순.
- 검토 필요: 독립 코드 리뷰(relay 화면·정렬 Go 변경), designer EN 1폭 재검수(V-CS-01).
- 검증한 것: Linux 격리 clone `9f279b4f` `make lint`·`make test` EXIT=0, `659380f1` integration EXIT=0(`ok 161.809s`, 그 뒤 차이는 사전 공백·front matter), Windows `go test ./internal/relay`, FullOps lint 깨끗한 worktree `--from e426f03be373` ERROR 0.
- 못 한 것: dev 작업 트리 FullOps lint(SAR-MOTION-INTRO 미추적 파일로 실행 거부), 캡처(lite), 임시 자료 정리(AGENTS.md 삭제 금지).
- 산출물·로그: [보고](../docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-test/report.md), 원시 로그 `D:/workspace/KnowsLink/.devtools/connect-simple-001-uifix-dev/`.
- 남은 일: 리뷰·재검수 배정(coor). 남은 임시 자료 경로는 보고서 끝에 있다.
