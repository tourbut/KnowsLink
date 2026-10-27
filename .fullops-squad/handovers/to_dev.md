---
title: "SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV — install 마운트를 compose.ops.yaml로 이전, 새 설치 npm audit 4건 확인"
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV]
summary: "install 마운트를 compose.ops.yaml로 이전, 새 설치 npm audit 4건 확인"
attempt: d2b88cca7e3f467393a272bce6c83405
base: 540df637470593f3edc281949bdaf3bf2c16dee4
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV — install 마운트를 compose.ops.yaml로 이전, 새 설치 npm audit 4건 확인

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 지시 'dev 후속 진행해'(2026-10-27). repo의 `deploy/knowslink/compose.ops.yaml`, 필요하면 `adapters/package.json`·`adapters/package-lock.json`, 관련 테스트·D12 문서 수정과 커밋. / 추가 승인 필요: 운영 서버 접속·변경(적용은 이후 ops 과제), major 버전 의존성 업그레이드나 새 의존성 추가, 제품 동작 변경.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

1. 운영은 bf67dfcb다. `compose.yaml`·`compose.ops.yaml`이 relay에 `KNOWSLINK_INSTALL_DIR`를 넘기지 않아, ops가 서버 로컬 `knowslink-state/compose.install.yaml`(relay env `KNOWSLINK_INSTALL_DIR: /install`, 볼륨 `<state>/install:/install:ro`)을 모든 compose 호출에 `-f`로 얹었다. 근거는 D12 `.fullops-squad/docs/operations/ops-guide.md` 맨 위 운영 적용 절('compose override'·'compose 호출 규칙'·`beta.sh` 주의)이다.
2. tester LIVE 보고 관찰 1: Grok Bot 컴퓨터에서 `scripts/install_bot_mcp.sh`의 `npm ci --prefix adapters`가 `npm audit` 취약점 4건(moderate 2, high 2)을 출력했다. 직접 확인하지 않았고 런타임/빌드 범위가 미확인이다(`.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md` 71행).

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(coding-style·testing·security)이다. 기준 커밋 main `540df637`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/operations/ops-guide.md`(맨 위 SAR-CONNECT-SIMPLE-001-OPS 운영 적용 절)
- `deploy/knowslink/compose.ops.yaml`, `compose.yaml`, `deploy/knowslink/beta.sh`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md`(관찰 1)
- `adapters/package.json`, `adapters/package-lock.json`, `scripts/install_bot_mcp.sh`

## 해야 할 일과 파일 소유권

- [ ] `compose.ops.yaml`의 relay에 `KNOWSLINK_INSTALL_DIR: /install`과 `${KNOWSLINK_STATE_DIR}/install:/install:ro` 볼륨을 더한다(기존 `KNOWSLINK_STATE_DIR` 필수 변수 패턴을 따른다). 다른 서비스·제한값은 바꾸지 않는다. 서버 override와 함께 얹혀도 결과가 같게(중복 볼륨 오류가 없는지) `docker compose -f compose.yaml -f deploy/knowslink/compose.ops.yaml [-f override] config`로 확인한다.
- [ ] `beta.sh`가 relay를 만들 때 install 마운트가 빠지는지 확인하고, 빠지면 같은 파일을 쓰도록 고치거나 D12에 이유를 적는다.
- [ ] D12의 'compose override'·'compose 호출 규칙' 줄을 새 구성 기준으로 바꾸고, ops가 서버 `compose.install.yaml`을 빼는 절차(전환 순서·확인 명령·되돌리는 법)를 적는다. 서버에는 접속하지 않는다.
- [ ] `adapters`에서 `npm audit --omit=dev`와 `npm audit`를 돌려 4건의 패키지·경로·런타임 여부를 표로 적는다. 런타임 경로이고 semver 호환 수정이 있으면 lockfile만 갱신(`npm audit fix` 비강제)하고 `npm run check`·`npm test`를 통과시킨다. major가 필요하면 고치지 않고 보고한다.

소유 파일: `deploy/knowslink/compose.ops.yaml`, `deploy/knowslink/beta.sh`(필요 시), `adapters/package*.json`(필요 시), D12 `.fullops-squad/docs/operations/ops-guide.md`, `qa-reports/SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV-test/`, 이 인박스.

## 완료 기준과 검증

- compose config 출력에서 relay env·`/install:ro` 마운트가 override 없이 보이고, override를 얹어도 오류가 없다. npm audit 표와 처리 결과(수정·보고)가 있다. lockfile을 바꿨으면 `npm run check`·`npm test` EXIT 0, `make lint`·`make test` 결과(Windows에서 안 되면 사유).
- 증명 방법: `qa-reports/SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV-test/report.md`(front matter 필수, 명령·종료코드).
- 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/dev --from 540df637` ERROR 0. `work.py finish`, `worker_done`.
- worker 완료는 repo 변경까지다. 운영 서버 전환은 이후 ops 과제다.

### UI 작업의 디자인 기준과 검증

해당 없음(배포 설정·의존성).

## 갱신할 산출물

D12 ops-guide.

## 기대 산출물

`compose.ops.yaml` 수정, D12 갱신, 필요 시 lockfile, 보고서.

## 제약·협업·후속

하위 에이전트 배정 안 함(작은 직렬 작업). 남은 일에 ops 서버 전환(override 제거)을 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV-packet.json` / SHA `540df637470593f3edc281949bdaf3bf2c16dee4` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/ops-guide.md` (document_update, document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/project.md` (document_read) · 줄 32, 50 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 44 · inferred · 필수
- `deploy/knowslink/compose.ops.yaml` (direct_edit, impact_check) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 16, 137, 201 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 72, 102 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 39, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 15, 17, 17, 117, 121, 135, 141, 145, 228 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 244, 314 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 60, 134, 300 · inferred
- `adapters/package.json` (direct_edit) · 줄 전체/미확인 · inferred
- `compose.yaml` (direct_edit) · 줄 전체/미확인 · unknown
- `deploy/knowslink/beta.sh` (direct_edit) · 줄 전체/미확인 · unknown
- `deploy/knowslink/test_public_config.sh` (impact_check) · 줄 8, 16, 17, 19, 37, 40 · inferred
- `docs/README.md` (document_read) · 줄 38 · inferred
- `docs/development.md` (document_read) · 줄 9, 9, 9, 48, 50, 62 · inferred
- `docs/operations.md` (document_read) · 줄 33 · inferred
- `docs/system-design.md` (document_read) · 줄 74, 74, 74, 301, 301, 303, 303, 318 · inferred
- `scripts/install_bot_mcp.sh` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 36건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 6건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; compose.yaml: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄, 상세는 `qa-reports/SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV-test/report.md`. worker_done 본문 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV | 브랜치 fullops/dev | SHA <브랜치 HEAD>` 한 줄로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
