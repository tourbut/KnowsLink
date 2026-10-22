---
title: "SAR-CODEX-RECV-001-DEV — C-01 Codex 세션 수신: SP-01–12·V02 기준 Claude↔Codex·Codex↔Codex 근거별 지원 구현·검증"
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-CODEX-RECV-001-DEV]
summary: "C-01 Codex 세션 수신: SP-01–12·V02 기준 Claude↔Codex·Codex↔Codex 근거별 지원 구현·검증"
attempt: 0dda905847aa4579adee0fdf140f62f3
base: 9c397d98bf2255a0624728b85055f97be6a22eb8
subagent_level: standard
test_level: lite
---

# SAR-CODEX-RECV-001-DEV — C-01 Codex 세션 수신: SP-01–12·V02 기준 Claude↔Codex·Codex↔Codex 근거별 지원 구현·검증

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/dev 브랜치의 adapters·internal·scripts·기술 문서 수정, 격리 로컬 relay와 이 PC의 Codex CLI(0.163.0, 이미 로그인됨)·Claude Code 세션으로 테스트 계정 간 검증, 커밋·push. 추가 승인 필요: 운영 relay 배포·운영 호스트 변경, 실제 다른 사람에게 발송, 새 실계정 생성, Codex 전역 설정(~/.codex) 변경, daemon 상시 설치, 사용자 제품 규칙(SP-01–12·HP 발송 승인) 변경
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

사용자가 2026-10-22 이슈 #9 결정 D027(추천 A)로 다음 개발 1순위를 C-01 Codex 세션 수신으로 정했다(`.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` 2절 C-01). Claude Code 세션 연결은 SAR-SESSION-PEER-001/002로 구현·검증됐다(`.fullops-squad/docs/design-docs/module-design.md` "SAR-SESSION-PEER-001-DEV Claude Code 세션 연결" 절, adapters/hooks/register.ts·adapters/src/text.ts·session-lock.ts). Codex 쪽은 송신 stdin·수신 CLI(trial-cli.ts)만 있고, 실행 중 Codex 세션 수신은 미구현이다. 설계 `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`의 Codex 행(App Server·`codex queue --thread --message`)과 V02가 기준이다. 이 PC의 Codex CLI는 0.163.0, ChatGPT 로그인 상태다(설계 당시 0.162.1).

## 적용 기준과 예외

공통 기준 `.fullops-squad/rules/common/README.md`와 연결 규칙. 기준 커밋 main `9c397d98`. 제품 정본: `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` SP-01–12·V01–07, `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` HP-01–16(HP-10 발송마다 승인, HP-15 도구 권한), SoT `docs/planning/sources/silent-agent-relay/decisions.md` 2026-10-22 Reply-loop invariant(자동 답장 0). 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`(SP-01–12, Codex 근거 표, V02)
- `.fullops-squad/docs/design-docs/module-design.md`(세션 연결 절)
- `adapters/hooks/register.ts`, `adapters/src/text.ts`, `adapters/src/session-lock.ts`, `adapters/src/trial-cli.ts`, `adapters/README.md`
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`
- `.fullops-squad/contexts/dev.md`

## 해야 할 일과 파일 소유권

- [ ] 기술 계획: Codex 수신 경로 선택(`codex queue`, App Server remote, MCP 다음 조회 중) 근거와 한계를 `docs/exec-plans/phases/SAR-CODEX-RECV-001-DEV.md`에 짧게 적는다. 친구 본문을 소유자 지시로 넣지 않는 untrusted 데이터 경계를 반드시 포함한다
- [ ] Codex 세션 연결·수신 구현: SP-01(명시 선택·이름), SP-04(한 세션 한 수신), SP-06(자동 답장 0), SP-12(수준별 지원 표기), 발송마다 승인 유지. Claude Code 쪽 기존 동작은 회귀 없이 유지
- [ ] 테스트: 단위 테스트 추가, 기존 `adapters` 테스트·Go 테스트 통과
- [ ] V02 실측: 격리 로컬 relay에서 Claude→Codex, Codex→Claude, Codex↔Codex를 각각 기록. remote 준비 세션과 독립 TUI를 별도 행으로, 기존 세션 수신 / 다음 조회 / 새 세션 수준을 구분해 report.md에 적는다
- [ ] 기술 문서 갱신: D10 module-design.md, D09/D05 해당 절, adapters/README.md 지원 표

파일 소유권: adapters/, internal/relay(필요 시), scripts/, 위 문서, dev 인박스. 예상 규모: 코드 300~800줄.

## 완료 기준과 검증

- Codex 세션 수신이 근거가 있는 수준으로만 지원 표기된다(SP-12). 근거 없는 방향·수준은 "미지원/미검증"으로 남긴다. 그 자체는 실패가 아니다.
- 친구 본문이 Codex 모델에 소유자 지시로 들어가지 않는다(외부 데이터 표시). 발송은 매번 소유자 승인, 자동 답장 0.
- 검증(dev 담당, 완료 SHA에서 1회): 제품 테스트(make lint/test는 Linux clone 또는 WSL, 이전 과제 방식), FullOps `python <FullOps 1.8.2 scripts>/lint.py --repo . --from 9c397d98` ERROR 0, 결과는 `docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`. 종료코드를 가리지 않는다.
- worker 완료는 구현·검증·문서까지다. 독립 tester QA와(대규모면) 독립 코드 리뷰는 coor가 후속 배정한다. 운영 배포는 사용자 결정이다.
- 멈출 조건: 추가 승인 행위, 설명되지 않는 검증 실패, SP/HP 제품 규칙을 바꿔야만 하는 경우(ask).

### UI 작업의 디자인 기준과 검증

웹 UI 변경 없음 예상. Codex 쪽 표시(터미널 출력·상태)는 SP 규칙의 문구를 따른다. 웹 화면을 바꾸게 되면 D04 기준과 이유를 완료 보고에 적는다.

## 갱신할 산출물

D10(module-design.md), D09·D05 해당 절(route 추천, 해당 없으면 사유와 함께 없음), adapters/README.md. D11 사용자 설명서의 Codex 안내는 지원 판정 뒤 후속.

## 기대 산출물

- adapters 코드·테스트
- `docs/exec-plans/phases/SAR-CODEX-RECV-001-DEV.md`
- `docs/evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md`
- D10·관련 문서

## 제약·협업·후속

하위 위임은 standard 범위에서 필요할 때만. Codex 세션은 이 PC의 사용자 ChatGPT 구독을 쓰므로 실측 세션 수를 필요한 만큼으로 줄인다. `남은 일`에 tester 독립 QA(V02 재현)와 D11 갱신을 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CODEX-RECV-001-DEV-packet.json` / SHA `9c397d98bf2255a0624728b85055f97be6a22eb8` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 113, 127, 148 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 225, 360 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/hooks/register.ts` (direct_edit) · 줄 44, 70, 88, 112, 231 · inferred
- `adapters/src/session-lock.test.ts` (impact_check) · 줄 7, 8, 10, 11, 15, 16, 17, 19, 21, 24, 28, 29, 30, 32, 33, 34, 35, 45, 46, 49, 50, 52 · inferred
- `adapters/src/session-lock.ts` (direct_edit) · 줄 10, 18, 29, 57, 64 · inferred
- `adapters/src/text.ts` (direct_edit) · 줄 전체/미확인 · unknown
- `adapters/src/trial-cli.ts` (direct_edit) · 줄 3 · inferred
- `docs/system-design.md` (document_read) · 줄 71, 125 · inferred
나머지 참조 후보 22건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 5건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; adapters/src/text.ts: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
