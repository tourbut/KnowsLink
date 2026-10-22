---
title: SAR-CODEX-RECV-001-TESTER — Codex 세션 수신 후보 611356a9의 V02·SP 승인 경계 독립 QA
status: draft
updated: 2026-10-22
owner: tester
tasks: [SAR-CODEX-RECV-001-TESTER]
summary: Codex 세션 수신 후보 611356a9의 V02·SP 승인 경계 독립 QA
attempt: 16b1c9f8d37a4746a19bc5b5c7d53691
base: d8c112cd75970f3a937388830658e6156e7b569b
subagent_level: standard
test_level: lite
---

# SAR-CODEX-RECV-001-TESTER — Codex 세션 수신 후보 611356a9의 V02·SP 승인 경계 독립 QA

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 후보 `611356a9995371d7d1aef45feb74ef7031949684`을 tester 워크트리 밖 임시 디렉터리(`git worktree add --detach` 또는 `git archive`)로 빌드해 격리 로컬 relay(합성 회원 `*.invalid`)와 이 PC의 Codex CLI·Claude Code로 실측, QA 보고서 커밋·push. 추가 승인 필요: 운영 relay·실계정·실제 상대 발송, `~/.codex`·`~/.claude` 전역 설정 변경, 제품 코드 수정(결함은 보고만)
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

dev가 Codex 세션 수신(C-01, 사용자 결정 D027)을 fullops/dev `611356a9995371d7d1aef45feb74ef7031949684`로 완료했다. 지원 수준은 "세션 결속 로컬 사본 + 다음 조회 수신"이고, 승인은 MCP elicitation의 일회 4자리 확인 코드 입력만 수락한다(YOLO 자동 수락 결함 F1 수정). 사용자가 독립 코드 리뷰를 생략했으므로(waiver `qa-reports/SAR-CODEX-RECV-001-DEV-REVIEW-review/waiver.json`) 이 QA가 병합 전 유일한 독립 검증이다. DEV 보고: `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`. 후보는 아직 main에 없다.

## 적용 기준과 예외

공통 기준 `.fullops-squad/rules/common/README.md`, testing.md `담당과 반복 범위`. 대상 고정 SHA `611356a9995371d7d1aef45feb74ef7031949684`(fullops/dev). 제품 정본: `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` SP-01–12·V02, `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` HP-10. SoT decisions.md 2026-10-22 Reply-loop invariant. 예외: 대상 코드는 tester 브랜치에 없으므로 고정 SHA를 별도 경로로 읽는다.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`(DEV 실측 표·결함·관찰, 같은 SHA의 `git show 611356a9:<경로>`)
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`(SP·V02)
- `adapters/src/codex-session.ts`, `adapters/src/mcp.ts`, `adapters/README.md`(후보 SHA 기준)
- DEV 도구 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/codex-recv-001-dev/`(fixture.py·drive.mjs 재사용 가능, 결과는 독립 판정)
- `.fullops-squad/contexts/tester.md`

## 해야 할 일과 파일 소유권

- [ ] 후보 SHA 빌드·설치(package_plugin.py --host claude|codex --connect --verify), 새 격리 Compose 프로젝트와 합성 회원 A(Claude)·B·C(Codex)
- [ ] V02 재현: Claude→Codex, Codex→Claude, Codex↔Codex 각 1회 이상. App Server 세션과 독립 TUI(on-request)를 별도 행으로. 수준(도착 표시/로컬 사본/다음 조회)을 독립 판정
- [ ] 승인 경계(이번 QA 핵심): on/take/발송 각각에서 거절·Esc 취소·입력 없는 accept·틀린 코드 → 결속·발송 0; YOLO TUI에서 on·발송 거부; 승인 창 열린 동안 off/take 시 발송 0; 모델이 confirmed:true를 넘겨도 대화상자 없이는 발송 0; 친구 본문 안 지시문이 실행·자동 답장되지 않음(SP-06/07)
- [ ] 회귀: Claude Code `/knowslink-session` 기존 동작 1회 왕복, `npm test --prefix adapters` 후보 SHA 통과
- [ ] 보고서 `docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md`(front matter 필수, 원시 로그는 .txt/.log 또는 raw/)

파일 소유권: 위 보고서 폴더와 tester 인박스·contexts/tester.md만.

## 완료 기준과 검증

- 위 항목마다 PASS/FAIL/NOT RUN과 근거 파일. FAIL은 재현 절차·등급(critical/high/medium/low).
- critical/high가 있으면 병합 차단으로 보고한다(제품 코드는 고치지 않는다).
- tester lint: tester 워크트리 cwd `python <FullOps 1.8.2 scripts>/lint.py --repo . --from d8c112cd` ERROR 0.
- 정리: 실측 뒤 Codex/Claude 세션 off·종료, Compose 정지. 임시 빌드 경로는 지우지 말고 보고서에 경로를 남긴다.
- 멈출 조건: 추가 승인 행위, 설명되지 않는 환경 실패.

### UI 작업의 디자인 기준과 검증

해당 없음. 웹 UI 변경 없는 후보다. 터미널 대화상자 문구는 SP 규칙과 대조한다.

## 갱신할 산출물

없음. route 추천 D10·D12·D03은 읽기 전용(설계·구현 문서는 dev 소유). QA 보고서만 쓴다.

## 기대 산출물

- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md`와 증거 파일

## 제약·협업·후속

하위 위임 없음. Codex 세션은 사용자 ChatGPT 구독을 쓰므로 필요한 만큼만 연다. `남은 일`에는 결함이 있으면 dev 수정, 없으면 coor 병합과 D11 Codex 안내(designer)를 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CODEX-RECV-001-TESTER-packet.json` / SHA `d8c112cd75970f3a937388830658e6156e7b569b` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
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

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
