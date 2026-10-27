---
title: "SAR-AUTOSETUP-001-QA3 — O-04 수정 재검증(다중 플러그인 호스트 갱신, 0.1.0→0.1.1)"
status: draft
updated: 2026-10-28
owner: tester
tasks: [SAR-AUTOSETUP-001-QA3]
summary: "O-04 수정 재검증(다중 플러그인 호스트 갱신, 0.1.0→0.1.1)"
attempt: f70aac227f5a44ceba37d741507aafb6
base: 1e18baa0cecec0d0577c7dab03a7f40d437d33ae
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-QA3 — O-04 수정 재검증(다중 플러그인 호스트 갱신, 0.1.0→0.1.1)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-28
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: QA2와 같다(D039=A, D040 'Claude만 로그인'). 사용자가 로그인한 격리 설정 폴더 `D:\workspace\KnowsLink\.devtools\autosetup-001-qa\claude-config`를 `CLAUDE_CONFIG_DIR`로 쓰고, 격리 로컬 relay·합성 회원으로 검증한다. 격리 폴더 안에서만 둘째 플러그인을 추가·제거해도 된다. qa-reports 작성·커밋·push. / 추가 승인 필요: 사용자 실제 `~/.claude`·`~/.codex` 변경, 운영 relay·운영 agent, 실제 상대 발송, 비밀번호 입력, 자격 파일 읽기·복사·출력·커밋, 제품 코드 수정(결함은 보고만). 저장소 `AGENTS.md`의 'Windows 파일 작업 안전'을 따른다. 삭제·`git worktree remove`·이동 정리는 하지 않고, 남길 임시 경로를 완료 보고에 적는다(QA2에서 이 규칙을 어겼다).
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

QA2가 찾은 O-04(high)를 dev가 고쳤다(main `1e18baa0`, dev d8c996b4). `Get-HostJson`이 PS5.1 배열을 펼치고, 예상 밖 오류는 `[F-10]` S6 형식으로 끝난다. 공개 0.1.0은 고정이라 플러그인·설치기·화면 명령을 0.1.1로 올렸다. 근거: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-O04-DEV-test/report.md`, D12 운영 교체 후보 절. 실제 다중 플러그인 호스트 실행은 dev가 하지 않았다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing·security)이다. 기준 커밋 main(이 지시서 커밋). 제품 정본: `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`(AC-07, S5, S6). Codex는 D040으로 제외.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(`QA2` 절, `### O-04`, 재현 도구 `tools/qa2-repro.ps1`)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-O04-DEV-test/report.md`
- `scripts/autosetup/knowslink-setup.ps1`, `scripts/package_plugin.py`(`--release`)

## 해야 할 일과 파일 소유권

- [ ] 격리 relay에 main의 0.1.0·0.1.1 release 산출물을 함께 올린다.
- [ ] 다중 플러그인 환경(QA2와 같이 Orca `CLAUDE_CODE_PLUGIN_DIRS` inline 플러그인이 둘째 항목인 상태, 가능하면 격리 폴더에 실제 둘째 플러그인도 하나 설치)에서 0.1.1 명령을 재실행해 갱신이 끝나는지 확인한다(AC-07, S5). agent ID·키 파일 해시·웹 agent 수·친구 수가 그대로인지 전후로 비교한다.
- [ ] 0.1.0으로 설치된 상태에서 0.1.1 명령을 실행하는 갱신(실제 사용자 경로)을 1회 확인한다.
- [ ] 갱신 뒤 Claude Code를 다시 열어 `/knowslink-session on`이 같은 agent로 연결되는지 확인한다.
- [ ] 예상 밖 오류가 `[F-10]` 형식(실패 단계·바뀐 것·다음 행동)으로 끝나는지 합성 실패 1건으로 확인한다.
- [ ] 보고서에 QA3 절을 더하고 O-04 판정을 갱신한다.

소유 파일: `qa-reports/SAR-AUTOSETUP-001-QA-test/`(갱신), 이 인박스.

## 완료 기준과 검증

- 통과 조건: 다중 플러그인 갱신·0.1.0→0.1.1 갱신 PASS, O-04 해소, 미해결 high 0.
- 자격 파일 내용·토큰·이메일·지문 전체 값은 적지 않는다. 원시 로그는 `.txt`/`.log` 또는 `.fullops-squad/local/SAR-AUTOSETUP-001-QA3/`.
- 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from <이 지시서 기준 SHA>` ERROR 0. `work.py finish`, `worker_done`.
- 시간 예산 45분. 넘으면 escalation.

### UI 작업의 디자인 기준과 검증

해당 없음(동작 검수).

## 갱신할 산출물

없음(QA 보고서만).

## 기대 산출물

갱신된 `qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`.

## 제약·협업·후속

하위 에이전트 배정 안 함. 통과하면 coor가 사용자 승인을 받아 ops 운영 교체(0.1.1 게시·relay 배포·서버 override 제거)를 배정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-QA3-packet.json` / SHA `1e18baa0cecec0d0577c7dab03a7f40d437d33ae` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-O04-DEV-test/report.md` (document_update, document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 30, 32, 36, 75 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 141 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 122, 132, 185 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 99, 221, 225, 251, 283, 317 · inferred
- `internal/relay/i18n.go` (impact_check) · 줄 35, 47, 76, 100 · inferred
- `internal/relay/i18n_en_human.go` (impact_check) · 줄 33, 53, 54, 57, 61, 79, 80, 98, 99, 163, 168, 173 · inferred
- `internal/relay/i18n_test.go` (impact_check) · 줄 357, 358, 366, 367, 378, 416, 445, 446, 482 · inferred
- `internal/relay/landing_integration_test.go` (impact_check) · 줄 30, 43, 46, 97 · inferred
- `internal/relay/lang_integration_test.go` (impact_check) · 줄 66, 77, 82, 88, 109, 163, 194 · inferred
- `scripts/autosetup/knowslink-setup.ps1` (direct_edit) · 줄 193, 194, 195, 201, 202, 206, 215, 217, 229, 234, 247, 259, 359 · inferred
- `scripts/autosetup/test-get-registered.ps1` (impact_check) · 줄 1, 2, 2, 3, 7, 7, 15, 24, 28, 40, 45 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 9, 37, 143, 147 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 31, 45, 46, 47, 47, 52, 56, 58, 76, 77, 80, 81, 104, 113, 115 · inferred
나머지 참조 후보 35건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / SHA: fullops/tester / 완료 커밋(worker_done에 기록)
- 변경 이유: O-04 수정(0.1.1) 재검증. 제품 코드 변경 0, 보고서·증거·도구만.
- 결과: 호스트 항목 3개(inline+실제 로컬 플러그인+knowslink)에서 0.1.0→0.1.1 갱신·같은 명령 재실행 PASS, agent/키 해시/웹 agent/친구 동일, 재시작 뒤 `/knowslink-session on` 같은 agent, `[F-10]` 안내 PASS. O-04 해소, 미해결 high 0.
- 검토 필요: O-05(0.1.0은 항목 3개에서 첫 설치도 실패, 0.1.1에서 해소), O-06(갱신 뒤 호스트 플러그인 캐시·훅·스킬은 0.1.0 유지, MCP만 0.1.1) — dev 판단.
- 못 한 것: 0.1.1 신규 설치(3개 항목) 직접 실행, 동의 뒤 단계 F-10, 대화형 세션 화면.
- 산출물: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(QA3 절), `ev/qa3-*`, `tools/qa3-*`.
- 남은 임시 경로(삭제 안 함): `D:\workspace\KnowsLink\.devtoolsutosetup-001-qa3\`, Docker `kl-autosetup-qa3-pg`, 포트 18443 fixture 프로세스, 격리 로그인 폴더의 knowslink 0.1.1·`qa3-other` 플러그인.
- 남은 일: coor 사용자 승인 후 ops 운영 교체(0.1.1 게시·relay 배포·서버 override 제거).
항목당 1~2줄, 상세는 보고서 링크. worker_done 본문 첫 줄은 `[완료] SAR-AUTOSETUP-001-QA3 | 브랜치 fullops/tester | SHA <브랜치 HEAD>` 한 줄로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
