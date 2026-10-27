---
title: "SAR-CONNECT-SIMPLE-001-OPS — 연결 단순화(SAR-CONNECT-SIMPLE-001) 운영 배포: relay 먼저, 클라이언트·설치기 확인"
status: draft
updated: 2026-10-27
owner: ops
tasks: [SAR-CONNECT-SIMPLE-001-OPS]
summary: "연결 단순화(SAR-CONNECT-SIMPLE-001) 운영 배포: relay 먼저, 클라이언트·설치기 확인"
attempt: f3d4ffbc37374ae0ad4f81329276933a
base: bf67dfcbb4f7164149901d847edc0aa0e5783dc7
subagent_level: standard
test_level: lite
---

# SAR-CONNECT-SIMPLE-001-OPS — 연결 단순화·자동 설정 0.1.0 운영 배포(relay 먼저, 설치 파일·/install ingress)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-27
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D042 '배포 진행', D043 '자동 설정 함께 공개'(2026-10-27). link.knowslog.com의 KnowsLink checkout·relay를 main `bf67dfcb` 제품 코드로 전환한다. 배포 전 과제 고유 Postgres backup·rollback 태그(현재 운영 `4fdd149b`)를 만든다. `scripts/package_plugin.py --release`(운영 relay URL 기본값) 산출물 0.1.0을 서버 install 폴더에 두고 운영 env에 `KNOWSLINK_INSTALL_DIR`만 추가한다. live 공개 ingress에 repo `deploy/knowslink/tunnel/public-ingress.yml`의 `/install/<semver>/knowslink-*.(ps1|zip)` 한 경로만 추가한다(반영 전후 `scripts/check_public_ingress.py`). / 추가 승인 필요: 실계정 로그인·발송, 기존 volume·backup 삭제, 그 밖의 env·ingress·다른 서비스 변경, DB migration 발생, force-push, Cloud Console 변경.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 `4fdd149b`(SAR-PUBLISH-001-P03, D12 'SAR-PUBLISH-001-P03 운영 적용 결과' 절)다. main `bf67dfcb`에는 그 뒤 SAR-AUTOSETUP-001(자동 설정 화면·`/install` 경로, QA 일부 완료·D040 열림)과 SAR-CONNECT-SIMPLE-001(연결 한도 제거·연결 시도 정리·화면 수정, 독립 리뷰 critical/high 0·tester QA·designer 재검수 PASS)이 들어 있다. 제품 diff는 59 파일이며 `deploy/knowslink/tunnel/public-ingress.yml`, `scripts/install_bot_mcp.sh`, `scripts/package_plugin.py`, `scripts/autosetup/knowslink-setup.ps1`, `adapters/package.json`을 포함한다. DEV 주의: relay를 클라이언트·설치기보다 먼저 배포한다(이전 relay는 `/v1/connect/start`의 `host`를 422로 거절). 설치기는 relay가 제공하므로(`internal/relay/install_files.go`, `scripts/install_bot_mcp.sh`) relay 전환이 곧 설치기 전환이다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `bf67dfcb`. 제품 정본: `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md`(CS-AC, L-03·L-04), `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`. 결정 D040·D042·D043은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/operations/ops-guide.md`(맨 위 'SAR-CONNECT-SIMPLE-001-DEV agent 용량 경계', 'SAR-PUBLISH-001-P03 운영 적용 결과' 절 — 지난 배포 순서)
- `.fullops-squad/docs/operations/transition.md`(rollback 절차)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md`('배포 주의' 절)
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTOSETUP-001-DEV.md`(배포 위치 행), `.fullops-squad/docs/design-docs/interface-design.md` `/install` 줄
- `.fullops-squad/handovers/logs/2026-10-22_to_ops.md`(SAR-PUBLISH-001-OPS 지시서·완료 보고 — 같은 형식의 직전 배포)

## 해야 할 일과 파일 소유권

- [ ] 배포 전: 과제 고유 private backup, rollback 태그(`knowslink-rollback:4fdd149b-connect001-20261027` 형식), 현재 상태(healthz, 활성 owner의 철회되지 않은 agent 수) 기록. agent 수가 200 이하인지(거절 임계 2000이 10배 이상, L-03) 확인한다.
- [ ] 제품 diff `4fdd149b..bf67dfcb`에서 schema·`db/`·`compose.yaml`·의존성 변경 여부를 확인한다. migration이 생기면 멈추고 escalation한다.
- [ ] relay를 `bf67dfcb`로 전환한다(지난 배포와 같은 무삭제 bundle·build·writer stop·backup·격리 복원·relay up 순서).
- [ ] `package_plugin.py --release <서버 install 폴더>`로 0.1.0을 만들고 운영 env에 `KNOWSLINK_INSTALL_DIR`를 추가해 relay를 다시 띄운다(이전 env는 0600으로 보관).
- [ ] live ingress에 `/install/...` 한 경로를 추가하고 `https://link.knowslog.com/install/0.1.0/knowslink-setup.ps1`이 200(`text/plain; charset=utf-8`), zip 2개가 200, 다른 이름·버전은 404, 그 밖의 경로 노출이 바뀌지 않았는지 확인한다.
- [ ] live 확인(실계정 로그인·발송 없이): healthz, `/`, `/home/agents` 비로그인 리다이렉트, `/v1/connect/start`가 `host` 필드를 받는지(422 아님), `install_bot_mcp.sh` 응답, 기존 친구·메시지 회귀 없음(서버 측 상태로 확인 가능한 범위), 첫 기동 로그에 오류·`agent capacity warning` 없음.
- [ ] D12 'agent 용량 경계' 절의 '운영 미적용(DEV 후보)'을 현행으로 바꾸고 운영 적용 결과 절(시각·중단 시간·백업·태그·install 폴더)을 더한다. D13 현재 운영 상태를 갱신한다.

소유 파일: D12 `.fullops-squad/docs/operations/ops-guide.md`, D13 `.fullops-squad/docs/operations/transition.md`, `qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/`, 이 인박스.

## 완료 기준과 검증

- 운영이 `bf67dfcb` 제품 코드이고 healthy, restart 0, migrate 0, `/install/0.1.0/` 3파일 200, 기존 회원·agent·친구 보존. 실패하면 rollback 절차(D13)를 따르고 보고한다.
- 증명 방법: `qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md`(front matter 필수, 명령·시각·종료코드 한 줄씩, 비밀 값·이메일 없음).
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/ops --from bf67dfcb` ERROR 0. `work.py finish`, `worker_done`.
- worker 완료는 배포·서버 측 확인까지다. 실계정 연결 시험(CS-AC-18)과 자동 설정 실세션 QA는 tester가 맡는다.

### UI 작업의 디자인 기준과 검증

해당 없음(배포). live 화면은 상태코드·리다이렉트만 확인한다.

## 갱신할 산출물

D12 ops-guide, D13 transition.

## 기대 산출물

운영 반영, D12·D13 갱신, `qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md`.

## 제약·협업·후속

하위 에이전트 배정 안 함(직렬 운영 작업). 배포 뒤 coor가 tester에 CS-AC-18 실연결(Grok=회원 A, Claude Code=회원 B)과 옛 CONNECT 시나리오 기대값 4건 갱신을 배정한다. ops는 실계정 로그인·발송을 하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-CONNECT-SIMPLE-001-OPS-packet.json` / SHA `bf67dfcbb4f7164149901d847edc0aa0e5783dc7` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/ops-guide.md` (document_update, document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/operations/transition.md` (document_update, document_read) · 줄 109, 113, 127, 137, 220 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 16, 16, 137, 201, 201, 281 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 42, 42, 45 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 102 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTOSETUP-001-DEV.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 47, 108, 120, 125, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 52, 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 119, 127, 184 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 209, 244, 314 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 36 · inferred
- `adapters/README.md` (document_read) · 줄 60, 60, 134, 300 · inferred
- `docs/README.md` (document_read) · 줄 38, 38 · inferred
- `docs/operations.md` (document_read) · 줄 19, 27, 27, 33, 33, 38, 53, 54 · inferred
- `docs/system-design.md` (document_read) · 줄 40, 73, 74, 302, 304, 318 · inferred
나머지 참조 후보 38건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 5건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / SHA: `fullops/ops` / 완료 `2b4e2c7b`(문서·증거 커밋, 제품 코드 변경 없음). 이후 finish 기록 커밋이 HEAD다.
- 변경 이유: D042·D043에 따라 운영 relay를 main `bf67dfcb`로 전환하고 자동 설정 0.1.0을 `/install/0.1.0/`에 공개했다. 중단 약18초(13:27:25Z~13:27:43Z), healthy·restart0·migrate0, migration 없음.
- 지시와 달라진 판단: `compose.yaml`이 relay에 `KNOWSLINK_INSTALL_DIR`를 넘기지 않고 폴더를 마운트하지 않아 env 추가만으로는 404였다. coor 승인(옵션 A)으로 서버 로컬 `compose.install.yaml`(env 1줄 + 읽기 전용 마운트 1개)을 썼다. 운영 `.env`는 바꾸지 않았다.
- 검토 필요: 위 override가 repo 밖이다(D12에 경로·내용·호출 규칙·되돌리는 법 기록). ingress는 `/install/...` 한 대안만 추가했다.
- 검증: 백업 `4fdd149b-pre-connect001-20261027T132730Z.dump`(0600·11464bytes)의 격리 복원 digest 일치. 공개 probe 21/21(설치 3파일 200+SHA 일치, 이웃 404 9개, `host` 필드 수용, `/home/agents` 303, 차단 경로 유지), 보존 비교 39/39, 첫 기동·5분 뒤 error/panic 0, `agent capacity warning` 0, live agent 30. `lint.py --from bf67dfcb`(HEAD `2b4e2c7b`): ERROR 0·WARNING 2(SIZE-001)·실행 불가 1(`make lint`, Windows `make` 부재). `make test`는 skipped.
- 못 한 것: 실계정 로그인·발송, CS-AC-18 실연결, 설치기 실세션 QA, 실제 rollback, 브라우저 렌더링.
- 산출물·로그: [report.md](../docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md), D12 `ops-guide.md`, D13 `transition.md`. 서버 작업물 `knowslink-state/connect001-ops-20261027/`, 로컬 원본 `.devtools/connect-simple-001-ops/`.
- 남은 일: tester: 배포된 운영에서 CS-AC-18 실연결(Grok=회원 A, Claude Code=회원 B)과 자동 설정 실세션 QA, 옛 CONNECT 시나리오 기대값 4건 갱신. dev: `compose.install.yaml`의 두 줄을 `compose.ops.yaml`로 옮겨 repo가 추적하게 한다(override 파일 제거).
- 의존성: 변경 없음(lockfile·go.mod 무변경, 제품 코드 무변경).

원본 안내:
항목당 1~2줄로 쓰고 상세는 `qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md` 링크로 대신한다. worker_done 본문 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-OPS | 브랜치 fullops/ops | SHA <브랜치 HEAD>` 한 줄로 쓰고 SHA는 하나만 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
