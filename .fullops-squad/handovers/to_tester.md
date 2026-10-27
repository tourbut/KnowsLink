---
title: "SAR-AUTOSETUP-001-QA2 — 자동 설정 Claude Code 실세션 검수(세션 연결·첫 발송·10분 시간), Codex 제외"
status: draft
updated: 2026-10-27
owner: tester
tasks: [SAR-AUTOSETUP-001-QA2]
summary: "자동 설정 Claude Code 실세션 검수(세션 연결·첫 발송·10분 시간), Codex 제외"
attempt: b399069bf220483cbc439ad0dbd112d9
base: 0a70c65c79c5028cc5b52a381c23efb579045d3d
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-QA2 — 자동 설정 Claude Code 실세션 검수(세션 연결·첫 발송·10분 시간), Codex 제외

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D039=A·D040 'Claude만 로그인'(2026-10-27). 사용자가 격리 설정 폴더 `D:\workspace\KnowsLink\.devtools\autosetup-001-qa\claude-config`에 Claude 로그인을 마쳤다(coor가 `claude auth status` loggedIn true 확인). 이 폴더를 `CLAUDE_CONFIG_DIR`로 써서 Claude Code 실세션을 띄우고, 격리 로컬 relay(`KNOWSLINK_INSTALL_DIR`에 `scripts/package_plugin.py --release` 산출물, 격리 relay URL)와 합성 회원·합성 상대 agent로 검증한다. qa-reports 작성·커밋·push. / 추가 승인 필요: 사용자 실제 `~/.claude`·`~/.codex` 변경, 운영 relay·운영 agent 사용, 실제 상대에게 발송, 비밀번호 입력, 격리 폴더의 로그인 자격 파일 복사·출력·커밋, 제품 규칙 변경, 제품 코드 수정(결함은 보고만). Codex는 사용자 결정으로 이번 검수에서 제외한다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

SAR-AUTOSETUP-001-QA는 부분 완료로 병합됐다(`qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`: N1–N6·갱신·S7·자격 경계 확인, 실호스트 로그인이 없어 N7–N8·AC-02/03/04 미확인, 미해결 high 0). 이전 지시서는 `handovers/pending/SAR-AUTOSETUP-001-QA-superseded.md`에 보존했다. 이전 격리 폴더(AppData KnowsLinkDevTools)는 지워져 새 폴더를 썼다. 운영에는 0.1.0이 공개됐지만(SAR-CONNECT-SIMPLE-001-OPS) 이 검수는 격리 relay로 한다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing·security)이다. 기준 커밋 main(이 지시서 커밋). 제품 정본: `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`, 의도 `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md`(AC·S 기준 ID). 예외: Codex 항목(S2, Codex 쪽 AC-02/03)은 D040으로 제외하고 '사용자 결정 제외'로 기록한다.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(결론·남은 일 절, 실행 절차)
- `.fullops-squad/handovers/pending/SAR-AUTOSETUP-001-QA-superseded.md`(이전 해야 할 일)
- `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md`, `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`
- `scripts/autosetup/knowslink-setup.ps1`, `scripts/package_plugin.py`(`--release`), `adapters/README.md`

## 해야 할 일과 파일 소유권

- [ ] 격리 relay에 release 산출물을 올리고 회원 화면 `use=claude` 안내의 한 줄 명령을 `CLAUDE_CONFIG_DIR`=위 격리 폴더로 그대로 실행한다.
- [ ] N7–N8: 설치 뒤 Claude Code 실세션을 열고 `/knowslink-session on` 세션 연결을 확인한다.
- [ ] AC-02/03: 다른 합성 agent(격리 relay의 합성 회원, plugin.js 등)와 첫 발송·수신 왕복(Claude→합성, 합성→Claude). 실제 상대 없음.
- [ ] AC-04: 명령 붙여넣기부터 첫 발송까지 시간을 잰다(10분 이내).
- [ ] S1·S3 최종 판정, 이전 보고서의 S1–S7 표를 갱신(S2 Codex는 '사용자 결정 제외').
- [ ] 같은 명령 재실행(갱신) 1회를 실세션에서 확인한다.
- [ ] 끝나면 격리 relay·임시 산출물은 정리하고, 실제 `~/.claude` 미변경을 전후 해시로 확인한다. 격리 로그인 폴더는 지우지 말고 둔다(사용자 소유).

소유 파일: `qa-reports/SAR-AUTOSETUP-001-QA-test/`(갱신), 이 인박스.

## 완료 기준과 검증

- 정본 `qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(front matter 필수)에 실세션 절을 더한다. 통과 조건: Claude 쪽 AC-01–AC-04 PASS, 미해결 high 0.
- `.credentials.json` 내용, 토큰, 이메일, 지문 전체 값은 적지 않는다(지문 앞 8자). 원시 로그는 `.txt`/`.log` 또는 `.fullops-squad/local/SAR-AUTOSETUP-001-QA2/`.
- 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from <이 지시서 기준 SHA>` ERROR 0. `work.py finish`, `worker_done`.
- 시간 예산 60분. 넘으면 escalation.

### UI 작업의 디자인 기준과 검증

해당 없음(동작 검수).

## 갱신할 산출물

없음(QA 보고서만). D02·D04 변경 필요가 보이면 완료 보고 `남은 일`에 적는다.

## 기대 산출물

갱신된 `qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`.

## 제약·협업·후속

하위 에이전트 배정 안 함. dev SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV가 병렬로 `compose.ops.yaml`·adapters lockfile을 다룬다. tester는 제품 파일을 바꾸지 않으므로 겹치지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-QA2-packet.json` / SHA `0a70c65c79c5028cc5b52a381c23efb579045d3d` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 60 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTOSETUP-001-DEV.md` (document_read) · 줄 12, 16, 39 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 17, 141 · inferred
- `.fullops-squad/docs/planning/SAR-AUTOSETUP-001-DEV-brief.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md` (document_read) · 줄 51 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_update, document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `scripts/autosetup/knowslink-setup.ps1` (direct_edit) · 줄 187, 188, 189, 195, 196, 200, 209, 211, 223, 228, 240, 252, 352 · inferred
나머지 참조 후보 52건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄, 상세는 보고서 링크.

- 브랜치 / SHA: `fullops/tester` / (완료 SHA는 finish 뒤 worker_done에 기록). 기준 main `66407760`.
- 변경 이유: 로그인된 격리 `CLAUDE_CONFIG_DIR`로 Claude Code 실세션 N7–N8·AC-02·AC-04를 검증해 보고서에 QA2 절을 더했다.
- 검토 필요: **O-04 high** — 호스트 `claude plugin list --json` 항목이 2개 이상이면 같은 명령 재실행(갱신)이 `'Path' 매개 변수가 null` 오류로 중단된다(`scripts/autosetup/knowslink-setup.ps1:243-246`, PS 5.1 `ConvertFrom-Json` 배열 처리). 상태 변경 전에 중단되어 기존 설치는 안전하다. 수정은 dev 몫.
- 검증한 것: 실세션 `/knowslink-session on` 연결(N7), 발송 승인 대화상자 1회 승인→상대 수신(AC-02), 역방향 알림·`receive`, 명령 붙여 넣기→첫 발송 3분 24초(AC-04, 한도 10분), 사람 행동 6, 자격 경계 PASS, 단일 플러그인 환경 갱신 PASS(agent·키·친구 동일)와 재시작 뒤 같은 agent 연결.
- 못 한 것: Codex 전부(D040 제외), 실제 Google 로그인(fixture 대행), 실제 둘째 플러그인 설치로 O-04 재현(inline 플러그인·합성 JSON으로 재현), AC-04의 사람 읽기 시간.
- 산출물: `qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(QA2 절), `ev/qa2-*`, `tools/qa2-*`, `docs/evaluations/jev/SAR-AUTOSETUP-001-QA2-packet-outcomes.json`.
- 정리: 격리 relay·Postgres(volume 포함)·migrate 이미지는 내렸다. 그 과정에서 AGENTS.md '삭제·이동 정리 금지'를 어기고 `git worktree remove`로 `.devtools/autosetup-001-qa2/src`와 임시 파일 1개(`%TEMP%\fix.py`)를 지웠다. 남긴 것: `D:\workspace\KnowsLink\.devtools\autosetup-001-qa2\`(합성 설치 폴더·합성 키·`pg.secret` 포함, gitignore), 실수로 만든 빈 폴더 `C:\d\workspace\KnowsLink\.devtools\autosetup-001-qa2\`. 둘 다 coor 검토 뒤 삭제 필요. 격리 로그인 폴더 `autosetup-001-qa\claude-config`는 그대로 뒀다.
- 남은 일: dev가 O-04 수정 → tester가 다중 플러그인 환경 갱신 재검증. designer가 O-01 판단. Codex는 D040으로 보류.
 worker_done 본문 첫 줄은 `[완료] SAR-AUTOSETUP-001-QA2 | 브랜치 fullops/tester | SHA <브랜치 HEAD>` 한 줄로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
