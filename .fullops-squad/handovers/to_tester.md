---
title: "SAR-CONNECT-SIMPLE-001-LIVE — 운영 실연결 검증 CS-AC-18(회원 A Grok·회원 B Claude Code, 친구 요청·수락)과 옛 CONNECT 시나리오 기대 4건 갱신"
status: draft
updated: 2026-10-27
owner: tester
tasks: [SAR-CONNECT-SIMPLE-001-LIVE]
summary: "운영 실연결 검증 CS-AC-18(회원 A Grok·회원 B Claude Code, 친구 요청·수락)과 옛 CONNECT 시나리오 기대 4건 갱신"
attempt: b6e2ee12e38d415aab56a140569e489a
base: e5819ac7966ffd8c3a9779835d65b5c2cb02b760
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-LIVE — 운영 실연결 검증 CS-AC-18과 옛 CONNECT 시나리오 기대 4건 갱신

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정(2026-10-26, coor 전달)과 D042 배포 진행. 운영 relay `https://link.knowslog.com`(bf67dfcb)에서 회원 A(Grok Bot 쪽, 운영 테스트 계정)와 회원 B(Claude Code 쪽, 사용자 본인 계정)의 실제 연결, 두 회원 사이 친구 요청·수락. ego 브라우저에서 에이전트가 Google 계정 선택·동의·지문 비교·이름·연결 동의·친구 코드 요청·수락을 모두 직접 진행한다(사용자에게 넘기지 않는다). 두 회원의 실제 주소는 coor가 별도 메시지로 전달한다(문서·커밋에 적지 않는다, 회원 A/회원 B로만 쓴다). 새로 만든 agent와 연결 기록은 그대로 둔다. / 추가 승인 필요: 실제 메시지 발송(두 회원 사이 시험 1건도 coor 경유 사용자 승인 뒤), 비밀번호·2단계 인증 코드·새 기기 확인 입력, 사용자 실제 `~/.claude`·`~/.codex` 설정 변경, 두 회원 외 계정, 기존 운영 agent(`agent_c18976751d5a96d8adfbb0`, 이름 없는 agent, 기존 Grok `knowslink` 연결) 철회·정리·변경, 제품 코드 수정(결함은 보고만), 제품 규칙 변경.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

1단계(격리 fixture) CS-AC-01–17은 통과했다(`qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md`, ART 재검수 PASS). ops가 2026-10-27 운영 relay를 bf67dfcb로 전환했고 자동 설정 0.1.0을 `/install/0.1.0/`에 공개했다(`qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md`, D12 ops-guide 해당 절). 운영 live agent 30개. 이 과제는 2단계(CS-AC-18)와 1단계 QA가 남긴 옛 시나리오 기대 4건 갱신이다. 2단계 조건 정본은 `handovers/pending/SAR-CONNECT-SIMPLE-001-tester.md`의 '실연결 재검증 — 2단계' 절이다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`, security.md)이다. 기준 커밋 main `e5819ac7`. 제품 정본·완료 조건(CS-AC, L-·A-·W-·T- 규칙 ID): `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/handovers/pending/SAR-CONNECT-SIMPLE-001-tester.md`('실연결 재검증 — 2단계' 절)
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md`(CS-AC-18, E 표, W-01, A-06)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md`('이전 시나리오 회귀' 절의 기대 변경 4건)
- `.fullops-squad/docs/evaluations/scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md`, `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md`
- `adapters/README.md`(Claude Code 연결), `.fullops-squad/docs/operations/user-guide.md`(Grok Bot 연결)

## 해야 할 일과 파일 소유권

- [ ] coor가 보낸 두 회원 주소를 받기 전에는 로그인하지 않는다(메시지가 없으면 `ask`).
- [ ] 회원 A로 Grok Bot 쪽 agent를 새로 연결한다. 기존 Grok `knowslink` 연결은 바꾸지 않는다. 사용자 실제 Grok 설정을 바꿔야 하면 멈추고 `ask`한다.
- [ ] 회원 B로 Claude Code 쪽 agent를 새로 연결한다. 실제 `~/.claude`는 바꾸지 않는다(테스트 폴더의 agent 폴더·격리 설정만 쓴다).
- [ ] 두 연결 각각: 키 등록, 기본 이름, 링크 재방문 완료 화면(W-01), 실패 시도가 '내 에이전트'에 늘지 않음(CS-AC-05), 실패가 나면 E 표 사유·문구·다음 행동. 지문은 클라이언트 출력과 브라우저 화면을 직접 대조(앞 8자만 기록).
- [ ] 친구 코드 요청·수락 뒤 두 회원 화면에 서로의 agent가 이름과 함께 보이는지, 도구 결과의 표시 이름(A-06).
- [ ] 에이전트가 대신한 행동과 사람이 꼭 해야 한 행동(있다면)을 순서대로 기록한다.
- [ ] 시험 메시지 1건은 `ask`로 coor에 발송 승인을 요청하고, 승인 뒤에만 보낸다. 승인이 없으면 not_run.
- [ ] 옛 시나리오 기대 4건(QA 보고 '이전 시나리오 회귀' 절)을 현재 규칙에 맞게 시나리오 문서에서 갱신한다(제품 코드·테스트 코드는 바꾸지 않는다).

소유 파일: `qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/`, `docs/evaluations/scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md`, `docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md`, 이 인박스.

## 완료 기준과 검증

- CS-AC-18 PASS/FAIL/not_run과 근거가 `qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md`(front matter 필수)에 있다. 결함은 심각도·재현 절차. 통과 조건: 두 회원 연결 완료·친구 성립·서로의 agent 이름 표시, 미해결 high 0.
- 이메일 주소·키·토큰·쿠키·지문 전체 값은 어디에도 적지 않는다. 캡처에 이메일이 보이면 가리거나 커밋하지 않는다. 원시 로그·임시 파일은 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-LIVE/`.
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from e5819ac7` ERROR 0. `work.py finish`, `worker_done`.
- Google이 비밀번호·2단계 인증·새 기기 확인을 요구하면 입력하지 않고 `ask`한다. 시간 예산 90분, 넘으면 escalation.

### UI 작업의 디자인 기준과 검증

해당 없음(동작 검증). 화면은 D04 문구로 판정만 하고 시각 검수는 하지 않는다.

## 갱신할 산출물

없음(QA 보고서와 시나리오만). D02·D04·D12 변경 필요가 보이면 완료 보고 `남은 일`에 적는다.

## 기대 산출물

`qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md`, 갱신한 시나리오 2개.

## 제약·협업·후속

하위 에이전트 배정 안 함(실계정 브라우저 세션은 직렬). SAR-AUTOSETUP-001-QA(D040, `handovers/pending/SAR-AUTOSETUP-001-QA-tester.md`로 대기)와 별개이며 그 결과를 이 과제 PASS로 쓰지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-LIVE-packet.json` / SHA `e5819ac7966ffd8c3a9779835d65b5c2cb02b760` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_update, document_read) · 줄 185 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 227 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 141 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 15 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-COMPETE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 98 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 100 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/pending/SAR-CONNECT-SIMPLE-001-tester.md` (document_update, document_read) · 줄 21, 39 · inferred
나머지 참조 후보 16건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / SHA: `fullops/tester` / `ce40f8d2`(기록 커밋). 기준 main `3d67c5d8`.
- 변경 이유: 운영 relay에서 CS-AC-18 실연결 검증과 옛 CONNECT 시나리오 기대 4건 갱신.
- 판정: CS-AC-18 PASS, high 이상 결함 0. 회원 B(Claude Code) 연결·지문 일치·기본 이름·W-01·취소 시도 목록 불변, 회원 A(Grok Bot 새 설치 `knowslink-live`) 연결·W-01, 친구 해제 후 재요청·수락, 서로의 agent 이름 표시, 시험 메시지 1건 수신(sender 이름 `Grok Bot`).
- 검토 필요: 회원 A 동의 화면 승인은 사용자가 직접 했다(사람이 꼭 한 행동). Grok Bot 컴퓨터에 새 prefix `/workspace/.knowslink-live`·소스 폴더·MCP `knowslink-live`를 더했다(제거 방법은 보고서).
- 검증한 것과 못 한 것: 위 항목 확인. 못 한 것: 회원 A 동의 화면의 에이전트 직접 지문 대조·이름 입력, 자동 설정 설치기·Claude Code 세션 안 `/knowslink-session`, 양방향 대화. 캡처 없음.
- 산출물: [report](../docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md), [시나리오 UX-002](../docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md), [시나리오 CONNECT-SIMPLE](../docs/evaluations/scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md). 원시 로그는 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-LIVE/`.
- 남은 일: 관찰 1(새 설치 `npm audit` 4건 보고)을 dev가 의존성 범위로 확인. D02·D04·D12 변경 필요 없음. 새 agent·연결 기록·Grok 새 설치는 사용자 지시 전까지 유지.
