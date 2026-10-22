---
title: "SAR-CODEX-RECV-001-DEV — tester D-CR-01(medium) 수정: 새 MCP 프로세스의 첫 호출 경쟁에서 CodexSession이 둘 생겨 off 뒤 lock·자동 수신이 남는 문제. codexSession() 생성을 한 번으로 보장하고 RC-4·RC-5·tools/orphan.mjs로 회귀 확인"
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-CODEX-RECV-001-DEV]
summary: "tester D-CR-01(medium) 수정: 새 MCP 프로세스의 첫 호출 경쟁에서 CodexSession이 둘 생겨 off 뒤 lock·자동 수신이 남는 문제. codexSession() 생성을 한 번으로 보장하고 RC-4·RC-5·tools/orphan.mjs로 회귀 확인"
attempt: b912698f5bdf46839aa0a9dfe017bca0
base: 4c4e1f8912764bee7863847f12713c34a80702b1
subagent_level: standard
test_level: lite
---

# SAR-CODEX-RECV-001-DEV — tester D-CR-01(medium) 수정: 새 MCP 프로세스의 첫 호출 경쟁에서 CodexSession이 둘 생겨 off 뒤 lock·자동 수신이 남는 문제. codexSession() 생성을 한 번으로 보장하고 RC-4·RC-5·tools/orphan.mjs로 회귀 확인

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/dev의 adapters 수정·테스트, 격리 로컬 relay와 이 PC의 Codex CLI로 회귀 실측, 커밋·push. 추가 승인 필요: 운영 relay·실계정·실제 상대 발송, ~/.codex·~/.claude 전역 설정 변경, SP·HP 제품 규칙 변경
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

첫 시도(611356a9, 코드 51acd223)의 tester 독립 QA SAR-CODEX-RECV-001-TESTER 결과는 critical/high 0, medium D-CR-01이다(`.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md` "D-CR-01" 절). 새 Codex MCP 프로세스의 첫 호출 두 개(status+on, on+on)가 겹치면 `adapters/src/mcp.ts`의 `codexSession()`이 await 사이에 `CodexSession`을 둘 만든다. 결속된 고아 객체는 모듈 변수와 달라 `off`가 lock·자동 수신을 풀지 못한다(SP-09). 발송은 안전하다. 사용자가 독립 코드 리뷰를 생략했다(waiver). 이 수정 뒤 coor가 병합한다. 이 브랜치는 main을 병합한 `4c4e1f89`이다.

## 적용 기준과 예외

공통 기준 `.fullops-squad/rules/common/README.md`. 기준 `4c4e1f8912764bee7863847f12713c34a80702b1`. 제품 정본 `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` SP-02·SP-09. 이전 시도 기록은 `handovers/logs/2026-10-22_to_dev.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md`(D-CR-01, RC-4·RC-5, O-5)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/tools/approval-matrix.mjs`, `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/tools/orphan.mjs`(회귀 도구)
- `adapters/src/mcp.ts`(`codexSession`), `adapters/src/codex-session.ts`, `adapters/src/codex-session.test.ts`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`

## 해야 할 일과 파일 소유권

- [ ] `codexSession()`이 프로세스당 `CodexSession`을 정확히 하나만 만들도록 고친다(생성 promise 공유 등 가장 작은 변경). 동시 `on`+`on`은 대화상자·결속이 하나만 진행되거나 둘째가 명확한 결과를 받게 한다
- [ ] 단위 테스트: 동시 첫 호출에서 인스턴스 1개, `off` 뒤 lock 해제·수신 중지
- [ ] 회귀 실측: tester의 RC-4·RC-5와 `orphan.mjs`를 수정 빌드로 실행해 PASS 기록
- [ ] DEV report.md에 D-CR-01 수정 절 추가, D10 module-design.md 해당 행 갱신(필요 시)

파일 소유권: adapters/src/mcp.ts, codex-session.ts(필요 시), 테스트, DEV report, D10. 예상 규모 50줄 이하.

## 완료 기준과 검증

- D-CR-01 재현 절차에서 고아 결속·남은 lock·off 뒤 수신이 0이다.
- `npm run check --prefix adapters`, `npm test --prefix adapters` exit 0. 제품 make lint/test는 Linux clone(이전 방식) exit 0.
- FullOps lint `--from 4c4e1f89` ERROR 0. 결과는 DEV report와 완료 보고에.
- 멈출 조건: 추가 승인 행위, 설명되지 않는 실패.

### UI 작업의 디자인 기준과 검증

해당 없음. 웹 UI 변경 없음.

## 갱신할 산출물

D10(module-design.md, 동작이 바뀌면). 그 밖 없음.

## 기대 산출물

- adapters 수정·테스트
- DEV report.md D-CR-01 절

## 제약·협업·후속

하위 위임 없음. `남은 일`에 coor 병합과 designer D11 Codex 안내(YOLO 불가·코드 입력)를 적는다. tester 재QA는 필요 없다(회귀 도구 결과로 갈음, coor가 병합 때 판단).

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-22_to_dev.md

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CODEX-RECV-001-DEV-packet.json` / SHA `4c4e1f8912764bee7863847f12713c34a80702b1` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-RECV-001-TESTER-test/report.md` (document_read) · 줄 23, 24, 104, 104, 105, 105, 105, 106, 109 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 85 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 86, 315 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.2.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 29, 73 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 57, 123 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 42, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 46, 76 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/hooks/register.test.ts` (impact_check) · 줄 8, 43, 245, 287 · inferred
- `adapters/hooks/register.ts` (impact_check) · 줄 200, 261 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14, 36, 50 · inferred
- `adapters/src/codex-session.ts` (direct_edit) · 줄 29 · inferred
- `adapters/src/mcp.ts` (direct_edit) · 줄 72, 85, 483 · inferred
- `docs/system-design.md` (document_read) · 줄 71 · inferred
- `scripts/verify_grok_plugin.py` (impact_check) · 줄 35, 38 · inferred
나머지 참조 후보 7건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
