---
title: "SAR-AUTOSETUP-001 — C-03 클라이언트 자동 설정: intent·흐름·카피 설계와 DEV 지시서 작성"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-AUTOSETUP-001]
summary: "C-03 클라이언트 자동 설정: intent·흐름·카피 설계와 DEV 지시서 작성"
attempt: 2b2de46029f84b41a07d22aa87fa0755
base: cfdd642b3dfe5d8852df8a869d77a3a2499a5a02
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001 — C-03 클라이언트 자동 설정: intent·흐름·카피 설계와 DEV 지시서 작성

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/designer의 기획·설계 문서(intent, D02, D04 해당 절)와 dev 지시서 초안 작성·커밋·push. 추가 승인 필요: 제품 코드 수정, 운영 배포, 실계정 사용, intent의 approved 처리(사용자 승인 사항), 기존 확정 규칙(LG-04, HP-10, SP 규칙, D023–D030) 변경
- 상위 의도: `docs/planning/intents/SAR-AUTOSETUP-001.md`(이번 과제에서 초안 작성, 전체 S 기준)
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

사용자 결정 D027(2026-10-22)로 C-01 Codex 세션 수신 뒤 C-03 클라이언트 자동 설정 확대가 다음이다(`.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` 2절 C-03: 대상은 Claude Code와 C-01 뒤의 Codex, 키 복사·자격 공유·소유자 동의 없는 연결 금지 LG-04). C-01은 main에 병합됐다. 현재 설치는 호스트별로 `adapters/README.md`의 `package_plugin.py --host claude|codex --install-dir … --connect --verify`, Google 로그인·지문 확인, 세션 연결 명령을 사용자가 단계별로 따라 한다(`.fullops-squad/docs/operations/user-guide.md` Claude Code·Codex 절, 회원 화면 '어디에서 사용할까요?' 안내). 사용자가 "다음 작업 진행해"로 착수를 지시했다. FullOps 1.8.3 규칙상 두 역할 묶음이라 intent를 먼저 쓰고 사용자 승인 뒤 dev를 배정한다.

## 적용 기준과 예외

공통 기준 `.fullops-squad/rules/common/README.md`(fullops-common-0.3.5, testing.md 큰 틀 검수), 문서 작성 `.fullops-squad/docs/agents/document-writing.md`. 기준 main `cfdd642b3dfe5d8852df8a869d77a3a2499a5a02`. 제품 정본: `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md`(C-03, NG-05 검증 안 된 클라이언트 표기 금지), `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`(LG-03·04, HP), `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`(SP-12), SoT `docs/planning/sources/silent-agent-relay/`. intent 템플릿 `docs/planning/intents/_TEMPLATE.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md`(C-03·NG-05)
- `adapters/README.md`(현재 설치 절차)
- `.fullops-squad/docs/operations/user-guide.md`(Claude Code·Codex·Grok·Dots 절)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(D04, '어디에서 사용할까요?' 안내)
- `scripts/package_plugin.py`(설치물 생성기, 읽기 전용)
- `.fullops-squad/docs/planning/intents/_TEMPLATE.md`
- `.fullops-squad/contexts/designer.md`

## 해야 할 일과 파일 소유권

- [ ] `docs/planning/intents/SAR-AUTOSETUP-001.md` 초안(status draft): 배경·문제·바라는 결과·성공 기준 S1…(예: 새 사용자가 Claude Code/Codex에서 설치부터 첫 메시지 발송까지 걸리는 단계 수·수동 입력 수)·대상 흐름·제약(LG-04, 소유자 동의, 자격 복사 금지, NG-05)·범위 밖·열린 질문
- [ ] D02 `product-specs/SAR-AUTOSETUP-001.md`: 자동 설정의 제품 규칙(한 명령/한 화면에서 무엇이 자동이고 무엇이 반드시 사람 확인인지), 실패·재시도·기존 설치 업데이트 흐름, 수락 기준
- [ ] D04 해당 절: 회원 화면 '어디에서 사용할까요?'의 Claude Code·Codex 안내 변경안과 카피(KO/EN)
- [ ] dev 지시서 초안 `docs/planning/SAR-AUTOSETUP-001-DEV-brief.md`: 목표·확정 규칙·범위·완료 조건·증명 방법(기술 계획은 dev 몫)
- [ ] contexts/designer.md 3줄

파일 소유권: 위 문서와 designer 인박스만. 제품 코드는 고치지 않는다.

## 완료 기준과 검증

증명 방법: intent의 성공 기준이 관찰 가능한 수치·흐름으로 적혀 있고, D02 수락 기준이 각 S 기준에 연결된다.
- 지원 표기는 실측 근거가 있는 호스트(Claude Code, Codex)만이다. 다른 클라이언트 확대는 범위 밖(NG-05).
- 사용자 결정이 필요한 것은 intent `열린 질문`에 선택지·추천과 함께 적는다. 임의로 확정하지 않는다.
- lint: designer 워크트리 `python <FullOps 1.8.3 scripts>/lint.py --repo . --from cfdd642b` ERROR 0. 문서 과제라 QA report는 만들지 않는다.

### UI 작업의 디자인 기준과 검증

설계 문서만 쓴다. D04 변경안은 기존 디자인 정본(`docs/design-docs/mockups/service-member-ui.md`)의 토큰·컴포넌트를 재사용한다. 화면 구현은 후속 dev.

## 갱신할 산출물

D02(신규 SAR-AUTOSETUP-001), D04(회원 안내 절). route 추천 D10·D03·D05는 dev 후속 몫이라 이번에 쓰지 않는다.

## 기대 산출물

- `docs/planning/intents/SAR-AUTOSETUP-001.md`
- `docs/planning/product-specs/SAR-AUTOSETUP-001.md`
- D04 해당 절
- `docs/planning/SAR-AUTOSETUP-001-DEV-brief.md`

## 제약·협업·후속

하위 위임 없음. 남은 일: 사용자 intent 승인 → coor가 dev 배정. tester 후속은 lite 규칙상 새 사용자 흐름이므로 통합 후보에서 큰 틀 검수 1회로 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-packet.json` / SHA `cfdd642b3dfe5d8852df8a869d77a3a2499a5a02` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read, document_update) · 줄 109 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 56 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 59, 115 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 169, 200 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 119, 125 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 137, 177 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_update, document_read) · 줄 24, 37, 178, 186, 284 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read, document_update) · 줄 23 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read, document_update) · 줄 67, 71, 75, 117 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read, document_update) · 줄 106 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read, document_update) · 줄 364 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
- `README.md` (document_read) · 줄 26, 169, 197 · inferred
- `adapters/README.md` (document_update) · 줄 전체/미확인 · inferred
- `docs/README.md` (document_read) · 줄 20, 38 · inferred
- `scripts/package_plugin.py` (direct_edit) · 줄 25, 32, 99, 136, 172 · inferred
나머지 참조 후보 25건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-test/report.md`(역할별로 나눈 키는 `<키>-DEV` 그대로) 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-10-22, designer)

- 브랜치 / SHA: `fullops/designer` / 산출물 커밋 `b7f472bf`(기준 main `cfdd642b`). packet-outcomes·완료 보고 커밋이 뒤따른다.
- 변경 이유: D027에 따라 C-03 자동 설정을 Claude Code·Codex로 한정해 intent 초안·D02 규칙·D04 안내 변경안·DEV 지시서 초안을 썼다. 제품 코드는 바꾸지 않았다.
- 결정: 명령 하나 + 설치 계획 동의 뒤 자동(AS-01–03). Google·지문·이름·세션 연결·발송 승인은 사람(AS-04). 재실행은 갱신이며 agent·키 보존(AS-10). 실패 F-01–F-09는 재실행으로 이어감(AS-11). 지원 문구는 Claude Code 2.1.302·Codex CLI 0.163.0·Windows 11만(NG-05).
- 검토 필요: intent Q1–Q6(배포 형태, 설치기의 Google 연결 시작, 호스트 등록 범위, OS, 미검증 버전, Codex 선택지 시점). 모두 추천안 기준으로 D02·D04를 썼고 사용자 승인 전이다.
- 검증한 것: lint `--from cfdd642b` ERROR 0, WARNING 2(SIZE-001 PLANS.md·service-member-ui.md 기존 초과), `product-lint: make lint` 실행 불가(Windows에 make 없음), product-test skipped. 못 한 것: 현재 수동 절차 단계 수·시간 실측(문서 기준 추정, DEV 착수 때 실측).
- 산출물: [intent](../docs/planning/intents/SAR-AUTOSETUP-001.md), [D02](../docs/planning/product-specs/SAR-AUTOSETUP-001.md), [D04 절](../docs/design-docs/mockups/service-member-ui.md#claude-codecodex-자동-설정-안내--sar-autosetup-001), [DEV 지시서 초안](../docs/planning/SAR-AUTOSETUP-001-DEV-brief.md), [designer 컨텍스트](../contexts/designer.md). 의존성 변경 없음(문서만).
- 남은 일: 사용자 intent 승인과 Q1–Q6 답 → coor가 DEV 지시서 초안으로 dev 배정. 구현 뒤 tester 큰 틀 검수 1회(새 사용자 흐름, S1–S7), designer D04 시각 검수와 U8 표 갱신.
