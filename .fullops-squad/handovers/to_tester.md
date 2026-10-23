---
title: "SAR-AUTOSETUP-001-QA — Claude Code·Codex 자동 설정 큰 틀 검수(S1–S7, 호스트별 1회)"
status: draft
updated: 2026-10-23
owner: tester
tasks: [SAR-AUTOSETUP-001-QA]
summary: "Claude Code·Codex 자동 설정 큰 틀 검수(S1–S7, 호스트별 1회)"
attempt: 19d219b00469434185a3811d904f7b04
base: 4fefdf68fa55d2405bb7aa7c94a73439bf7ab235
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-QA — Claude Code·Codex 자동 설정 큰 틀 검수(S1–S7, 호스트별 1회)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-23
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/tester의 qa-reports·시나리오 문서 작성, 격리 로컬 relay(KNOWSLINK_INSTALL_DIR에 `scripts/package_plugin.py --release` 산출물)·합성 회원으로 검증, 이 PC의 실제 Claude Code·Codex 로그인을 **검증용 별도 CLAUDE_CONFIG_DIR/CODEX_HOME**에서 사용(사용자 결정 D039=A), 커밋·push. 추가 승인 필요: 사용자 실제 ~/.claude·~/.codex 설정 변경, 운영 relay·배포물 공개 게시, 실제 상대에게 발송, 비밀번호 입력, 제품 규칙(AS·D031–D039) 변경, 제품 코드 수정(결함은 보고만)
- 상위 의도: `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md` — S1–S7 큰 틀 검수(호스트별 1회)
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

dev SAR-AUTOSETUP-001-DEV 완료 c74c90e1 → main 8ec61ab3 병합(사용자 D038 리뷰 생략, 병합자 diff 확인). dev 자체 검사: AC-05–AC-10, F-01–F-09 KO·EN, AS-10 PASS(격리 fixture). dev가 실세션으로 확인하지 못한 부분: AC-01 N7–N8 사람 행동 표, AC-02/03 첫 발송, AC-04 10분 측정. 근거: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-DEV-test/report.md`.

## 적용 기준과 예외

- 공통 규칙: `.fullops-squad/rules/common/README.md`, `testing.md`(lite: 새 사용자 흐름 큰 틀 검수 1회). 기준 main 4fefdf68.
- `AS-`·`AC-`·`F-`·`N`/`U` 흐름 → `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md`(D02). `S1–S7` → intent. 결정 `D031–D039` → `fullops-decisions`(PLANS.md 기록). 산출물 `D04` → `docs/design-docs/mockups/service-member-ui.md`.
- 예외: 실 호스트 로그인 사용은 D039=A로 승인. 실제 Google 로그인은 합성 회원 경로를 쓰고 사용자 실계정 비밀번호를 입력하지 않는다(Google 로그인이 꼭 필요하면 ask).

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md`
- `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-DEV-test/report.md` (+ 같은 폴더 `scenarios.py`, `drive.py`)
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTOSETUP-001-DEV.md`
- `scripts/autosetup/knowslink-setup.ps1`, `scripts/package_plugin.py`(`--release`)
- `adapters/README.md`, 루트 `README.md` 자동 설정 절, D11 해당 절

## 해야 할 일과 파일 소유권

- [ ] 격리 relay에 release 배포물을 올리고 회원 화면 `use=claude`·`use=codex` 안내의 한 줄 명령을 그대로 실행(호스트별 1회, 각각 새 격리 설정 폴더)
- [ ] AC-01: N1–N8 사람 행동 표(단계·행동·입력값 종류) 작성, 기준선(intent 문제 절 추정)과 비교
- [ ] AC-02/03: 설치 직후 세션 연결 → 다른 합성 agent에게 첫 발송·수신까지(Claude→Codex, Codex→Claude)
- [ ] AC-04: 명령 붙여넣기부터 첫 발송까지 시간 측정(10분 이내)
- [ ] S1–S7 각각 PASS/FAIL/미확인 판정, F-0x 중 실세션에서 새로 나온 실패가 있으면 재현 절차
- [ ] 같은 명령 재실행(갱신) 1회 확인 — dev AC-07 결과를 실세션에서 재확인
- [ ] 결함은 qa-reports에 재현·영향·심각도로 기록(제품 코드 수정 안 함)

## 완료 기준과 검증

- 정본: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(front matter 필수). S1–S7 판정표, 호스트별 행동 표, 소요 시간, 명령·종료코드·기준 SHA.
- 원시 로그는 `raw/`(gitignore) 또는 `.txt`. 키·토큰·실제 이메일·지문 전체 값은 적지 않는다(지문은 앞 8자).
- 증명 방법: S1–S7 판정표와 실측 시간·행동 수. 통과 조건 = AC-01–AC-04 PASS, 미해결 high 0.
- lint: `python C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.3/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from 4fefdf68`, ERROR 0.
- 검증 뒤 격리 설정 폴더·임시 relay는 tester가 정리(격리 경로만, 실제 ~/.claude·~/.codex 미변경을 전후 해시로 확인).
- UI 작업: 해당 없음(시각 검수는 designer SAR-AUTOSETUP-001-ART).

## 갱신할 산출물

없음(QA 보고서만). D04·D02 변경 필요가 보이면 완료 보고 `남은 일`에 적는다.

## 기대 산출물

`.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`

## 제약·협업·후속

designer SAR-AUTOSETUP-001-ART와 병렬. 파일 소유권 겹침 없음(tester는 qa-reports/SAR-AUTOSETUP-001-QA-test만). 하위 에이전트 배정 안 함.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-QA-packet.json` / SHA `4fefdf68fa55d2405bb7aa7c94a73439bf7ab235` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 75 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-DEV-test/report.md` (document_read) · 줄 16 · inferred · 필수
- `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md` (document_read) · 줄 31, 56 · inferred · 필수
- `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 16 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read, document_update) · 줄 101 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
- `adapters/hooks/register.test.ts` (impact_check) · 줄 77, 78, 133, 146, 172, 176, 181, 185, 195, 199, 220, 225, 242, 246, 256, 260, 284, 288, 305 · inferred
- `internal/relay/i18n_en_human.go` (impact_check) · 줄 33, 53, 54, 57, 61, 79, 80, 98, 99, 131, 163, 168, 173 · inferred
- `internal/relay/i18n_test.go` (impact_check) · 줄 357, 358, 366, 367, 378, 416, 445, 446, 482 · inferred
- `scripts/autosetup/knowslink-setup.ps1` (direct_edit) · 줄 187, 188, 189, 195, 196, 200, 209, 211, 223, 228, 240, 252, 352 · inferred
- `scripts/package_plugin.py` (direct_edit) · 줄 25, 32, 99, 139, 162, 198 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 31, 45, 46, 47, 47, 52, 56, 58, 76, 77, 80, 81, 104, 113, 115 · inferred
나머지 참조 후보 50건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 5건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-test/report.md`(역할별로 나눈 키는 `<키>-DEV` 그대로) 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
