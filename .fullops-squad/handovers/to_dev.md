---
title: SAR-GROK-WEBHOOK-001-DEV — 사용자 메시지 없는 webhook routine 깨우기 후보 구현
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-GROK-WEBHOOK-001-DEV]
summary: 사용자 메시지 없는 webhook routine 깨우기 후보 구현
attempt: 590f4c27d88d410c8c15fe232581b35a
base: 7a8c2b9194edb77b52db41a14b2325af743f0d6e
subagent_level: standard
test_level: lite
---

# SAR-GROK-WEBHOOK-001-DEV — 사용자 메시지 없는 webhook routine 깨우기 후보 구현

- 작성일: 2026-10-11
- From / To: coor / dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 담당 repo id: 0b08ec4c-9e3a-4613-8197-5a835e545335
- 담당 워크트리: C:/Users/shin/orca/workspaces/KnowsLink/dev; 브랜치: fullops/dev
- 병합 책임자: coor; 기본 브랜치: main
- 복귀: C:/Users/shin/orca/workspaces/KnowsLink/coor, term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 run/task/dispatch는 Orca preamble을 사용한다.

## 적용 기준과 예외
fullops-common-0.3.3: rules/common/README.md와 coding-style.md, testing.md, security.md를 읽는다. project.md와 FULLOPS.md, contexts/dev.md, docs/agents/document-writing.md를 읽는다. 기준 SHA는 front matter의 base다. Ponytail full 적용. 실제 사용 중인 기존 에이전트 세션은 종료하지 않는다.
추가 과금 금지. 로컬 Docker 금지. 검증은 기존 Linux 서버에서 고유 후보 체크아웃으로 실행한다. 기존 DB/키/agent/pair/Google/Tunnel 자격/공유 서비스를 보존한다. 파일·컨테이너·임시 증거 삭제/이동/정리 금지. F:/KnowsLink-Incident-20261011-110130 보존. 비밀값은 로그/Git/채팅에 쓰지 않는다.
웹은 Orca 내장 브라우저만 사용한다. Computer Use는 Grok Bot 설정 제어에만 사용한다. Grok UI를 열어 수신을 유발하거나 수신 결과를 읽지 않는다. 제품 배포·Bot 설정 변경은 후속 검수 이후 coor/ops 담당이다.

## 완료 기준과 검증
현재 범위 완료 후 docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md에 기록하고, packet-outcomes의 모든 필수 path/category와 partial 불확실성을 처리한다. work.py finish로 전문 보존·빈 인박스를 확인한다. 커밋 후 fullops lint를 base 기준으로 실행하고 SHA/ERROR/WARNING/불가 사유를 기록한다. worker_done은 actual preamble로 보낸다. 구현과 제품 최종 수락을 구분한다.
멈출 조건: 파괴 작업, 범위 밖 파일 변경, 설명되지 않는 필수 검증 실패, 제품 규칙 변경만 ask/escalation한다. 비파괴 작업은 중간 승인을 기다리지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-GROK-WEBHOOK-001-DEV-packet.json` / SHA `7a8c2b9194edb77b52db41a14b2325af743f0d6e` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.64, 'evidence': 0.72, 'contradicts': 0.24, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.26, 'evidence': 0.47, 'contradicts': 0.09, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-DOTS-GROK-LIVE/webhook-wake-research.json` (impact_check) · 줄 10, 10, 14, 14, 18, 25, 26, 27, 28, 32, 32, 37, 38, 39 · inferred · 필수 · {'relevant': 0.88, 'evidence': 0.87, 'contradicts': 0.1, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 7, 7, 14, 14, 14, 14, 31, 31, 38, 39, 39, 43, 44, 46, 50 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 51 · inferred · 필수 · {'relevant': 0.46, 'evidence': 0.86, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.15, 'evidence': 0.45, 'contradicts': 0.18, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.24, 'evidence': 0.67, 'contradicts': 0.1, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.66, 'evidence': 0.84, 'contradicts': 0.09, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.47, 'evidence': 0.76, 'contradicts': 0.1, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `adapters/src/inbox.ts` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.47, 'evidence': 0.59, 'contradicts': 0.25, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 20, 25, 81, 97, 202 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 106, 112, 137, 250 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 69 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-LIVE/report.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 143 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 120 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md` (document_read) · 줄 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 81 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md` (document_read) · 줄 69, 102 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 41, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md` (document_read) · 줄 23, 62, 62, 73, 73, 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md` (document_read) · 줄 113 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 133 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 129 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 94 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 26, 33, 45, 51 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 28, 42, 233, 392 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 19, 27, 125 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 138 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 30, 37, 38, 42, 58 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 11, 33, 118, 121, 136 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 36, 113 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 364, 370, 384 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 97 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `adapters/.grok-plugin/plugin.json` (impact_check) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_update, document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.57, 'evidence': 0.58, 'contradicts': 0.28, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 26 · inferred
- `adapters/src/grok-wake.ts` (direct_edit, impact_check) · 줄 전체/미확인 · unknown
- `adapters/src/index.ts` (impact_check) · 줄 13 · inferred
- `adapters/src/mcp.ts` (direct_edit, impact_check) · 줄 23, 366 · inferred · {'relevant': 0.4, 'evidence': 0.51, 'contradicts': 0.45, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `internal/relay/text_events.go` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/install_bot_mcp.sh` (direct_edit, impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.27, 'evidence': 0.37, 'contradicts': 0.46, 'injection': 0.04, 'decision': 'keep', 'reason': None}
미확인 10건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
완료 시 변경 이유/실제 SHA/검증 명령과 종료코드/미검증/후속/산출물·로그를 전문으로 기록한다.

### 결과 (dev, 2026-10-11)

- 상태: 구현 후보 완료. 제품 최종 수락·실제 routine 왕복은 미완료다.
- 변경 이유: `sendPrompt` doorbell은 노우 대화에 사용자 메시지로 보인다. 공식 make-bot-ui 스킬의 webhook trigger routine은 routine 턴으로 깨운다. 기존 경로를 유지한 채 opt-in 후보를 추가했다.
- 구현: `adapters/src/grok-wake.ts`에 `KNOWSLINK_GROK_WEBHOOK_FILE` 모드를 추가했다. 비공개 `{url,key}` 파일을 요청마다 읽는다. URL은 `https://api2.cursor.sh/automations/webhook/<id>` 형식만 허용한다. Bearer·`X-Automation-Key` POST JSON 8초 단일 시도, redirect 미추종, 메타데이터(UUID·수·고정 안내)만 보낸다. 설정 시 `sendPrompt`와 배타이고 fallback이 없다. 미설정이면 기존 동작이다. `wake-check`는 webhook 모드에서 무전송 파일 검사다. 호출자·설치 스크립트는 변경이 필요 없어 그대로 두었다.
- 검사: `adapters/src/grok-wake.test.ts`(정상·3xx/4xx/5xx·DNS 실패·timeout·미설정·동시 설정·불신 URL/key·0644 파일·check 무전송·결과 비밀값 부재). 변이 3종(URL host 완화, redirect follow, ENOTFOUND 분류 제거) 검출. 중복·재시작 무재전송은 기존 marker 계약과 inbox.test 회귀로 확인했다.
- 문서: README 지원 표 정정·webhook 절, D03 architecture 일반화 정정·webhook 절, D10 module-design 행, D12 ops-guide 전환 절차, phase 기록, contexts/dev, packet outcomes(71항목, uncertainty_review 포함).
- 검증: Windows `npm run check` 0, `npm test` 0. Linux 고유 clean clone `f39b6777abe5168ed0eba0f5fab138f5a64f5ad2`에서 등록 product-lint passed, product-test passed, lint.py 종료코드 0 / ERROR 0 / WARNING 3 / 실행 불가 0. WARNING 처리는 phase 기록에 있다. 최종 archive HEAD의 gate는 worker_done에 적는다.
- 판단: 제품 원천 C5 "optional webhook OFF"는 relay→vendor webhook에 대한 규칙으로 보았다. 이 adapter-local opt-in wake는 기본 off라 규칙 변경이 아니라고 판단했다. 다르게 보면 coor·designer가 조정한다.
- 미검증·후속: 실제 routine 생성·key 발급·무해 probe·routine 턴 표시·노우 receive(coor·owner). 독립 tester QA·코드 리뷰(고정 SHA snapshot). 대체 확인 전 현재 `sendPrompt` 설정 유지. 배포·Bot 설정 변경은 하지 않았다.
- 산출물: [phase 기록](../docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md), `docs/evaluations/jev/SAR-GROK-WEBHOOK-001-DEV-packet-outcomes.json`. Linux 임시 경로 `/tmp/knowslink-webhook-f39b677-ctxbcbfa2f4992f{,-2}`는 보존했다.

## 현재 상황과 확인 근거
사용자는 사용자 메시지로 보이는 Grok sendPrompt 깨우기를 대체하되 대체 방식 확인까지 기존 경로를 유지하도록 지시했다. 실제 Bot 보고로 UpdateRoutine trigger.type=webhook 및 webhook_event 루틴 이벤트를 확인했다. 아직 루틴/키 생성·실전 시험은 하지 않았다.
근거: docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-DOTS-GROK-LIVE/webhook-wake-research.json. 공식 규격 https://github.com/cursor/plugins/blob/main/pstack/skills/make-bot-ui/SKILL.md. Context7은 이전 조회 quota exceeded로 공식 문서 fallback 기록이 있다.

## 해야 할 일과 파일 소유권
- [ ] adapters/src/grok-wake.ts, mcp.ts 및 모든 호출자·tests·env docs·Bot 설치 스크립트를 좁게 읽고 기술 계획을 같은 과제에 작성한다.
- [ ] 기존 SSE/polling/inbox/ACK와 logging 알림을 유지하면서 메타데이터만 전송하는 webhook 깨우기를 최소 구현한다. 기존 sendPrompt는 설정된 새 경로 검증 전 제거/자동 비활성화하지 않는다. 동시 이중 깨우기·실패시 무조건 legacy fallback으로 중복 호출하지 않는다.
- [ ] 공식 규격의 HTTPS api2.cursor.sh/automations/webhook/<실제 ID> URL을 검증하고 Authorization Bearer 및 X-Automation-Key로 POST JSON을 8초 단일 시도한다. URL/키는 비공개 파일/env로 받으며 미설정이면 현 동작 그대로다. redirect/SSRF/키노출/received text 프롬프트 삽입을 차단한다. URL ID를 임의 생성하거나 비밀값을 산출물에 쓰지 않는다.
- [ ] 정상·오류·timeout·미설정·중복·불신 입력·로그 비밀값 부재의 관련 자동 검사를 추가한다. 실제 Bot 루틴 생성 및 배포는 검수 이후 후속이다.
- [ ] 설치/운영 설명과 D10/D12 영향 원천을 필요한 범위만 갱신한다. 기존 'webhook 없음' 일반화는 검증된 후보 범위로 정정한다.
파일 소유권: adapters/src/grok-wake.ts 및 관련 adapter tests/mcp 호출부/설치 문서·설정 예시, 필요 Bot 설치 script, 본인 phase/contexts/로그/packet outcomes, D10/D12 영향 원천. 회원 UI member.go 및 디자인 문서는 designer/UI 후속 소유이므로 수정하지 않는다. PLANS/board는 coor 소유다.

## 검증 담당과 후속
DEV는 실제 고정 후보 SHA를 Linux 고유 체크아웃에서 product lint/test 및 변경 경계 시험한다. 예전 eac603b6 product-byte shim을 새 코드 검증처럼 재사용하지 않는다. 원본 서버 배포 코드/상태 env 변경 금지. 테스트 자원 삭제 금지. 독립 tester는 이후 fixed SHA read-only snapshot에서 QA와 code review, 실전 Bot routine 시험은 coor가 담당한다.
UI 작업 해당 없음. 추가 의존성 필요 시 표준 라이브러리 대안부터 확인한다.

## 갱신할 산출물
D10, D12 영향 원천만 확인 후 필요한 경우 갱신한다. 나머지는 read-only 참조로 유지한다.
