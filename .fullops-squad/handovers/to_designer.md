---
title: SAR-CONNECT-SIMPLE-001-ART-RECHECK — V-CS-01–06 수정 확인과 CS-AC-17 최종 판정
status: draft
updated: 2026-10-27
owner: designer
tasks: [SAR-CONNECT-SIMPLE-001-ART-RECHECK]
summary: V-CS-01–06 수정 확인과 CS-AC-17 최종 판정
attempt: 47c7de52bbd746769bf7d08fc85ddb82
base: 06af3fe2318d1f6cf1e423e1f5257ae656c51895
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-ART-RECHECK — V-CS-01–06 수정 확인과 CS-AC-17 최종 판정

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: fullops/designer에서 재검수 보고·판정 캡처 커밋, main 06af3fe2 격리 로컬 relay·합성 fixture 실행, 커밋·push.
- 추가 승인이 필요한 행위: 운영 relay·agent 사용, 실제 발송, 제품 코드 수정, D02 변경.
- 상위 의도: 없음. D02 SAR-CONNECT-SIMPLE-001 CS-AC-17.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

SAR-CONNECT-SIMPLE-001-ART(9ed04461)가 V-CS-01 medium으로 CS-AC-17 FAIL을 냈다. SAR-CONNECT-SIMPLE-001-UI-FIX-DEV(061c4142)가 V-CS-01–06을 고쳐 main 06af3fe2에 통합됐다. 보고: `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-test/report.md`.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`, `testing.md`. 디자인 정본 `DESIGN.md`, D04 연결 단순화 절.
- 기준 main `06af3fe2318d`. ID: `V-CS-` → `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`, `CS-AC-` → D02 SAR-CONNECT-SIMPLE-001.
- 예외: 없음.

## 먼저 읽을 문서

- 이전 ART 보고 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`(fixture 띄우는 방법과 결함 표)
- UI-FIX-DEV 보고

## 해야 할 일과 파일 소유권

- [ ] V-CS-01: EN Google 실패 화면 1폭(1440) 캡처로 h1·`<title>`이 'Google sign-in'인지 확인.
- [ ] V-CS-02–06: KO·EN 390 한 폭에서 동의 직후 문장, 버튼 aria 이름(접근성 트리), agent_rate 복귀 링크, 클라이언트 줄 표시 이름, 시도 최신순을 확인.
- [ ] CS-AC-17 최종 판정(PASS/FAIL).

파일 소유권: `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-RECHECK-test/`만 쓴다.

## 완료 기준과 검증

- 보고 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-RECHECK-test/report.md`(front matter `deliverables.py --stamp`), 판정에 연결된 캡처만 커밋. FullOps lint `--from 06af3fe2318d` ERROR 0.
- 이전 검수 범위 전체를 반복하지 않는다(공통 테스트 기준 `담당과 반복 범위`).

### UI 작업의 디자인 기준과 검증

검수만 한다. 디자인 lint 해당 없음, 테마 전환 제품에 없음.

## 갱신할 산출물

없음.

## 기대 산출물

재검수 보고와 판정 캡처.

## 제약·협업·후속

운영 relay·agent는 쓰지 않는다. 하위 위임 불필요.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-ART-RECHECK-packet.json` / SHA `06af3fe2318d1f6cf1e423e1f5257ae656c51895` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
나머지 참조 후보 0건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

worker_done 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-ART-RECHECK | 브랜치 fullops/designer | SHA <HEAD>`.
