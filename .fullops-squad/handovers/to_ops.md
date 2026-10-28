---
title: "SAR-AUTOSETUP-001-OPS — 자동 설정 0.1.1 게시, relay main 배포, 서버 override를 repo compose로 대체"
status: draft
updated: 2026-10-28
owner: ops
tasks: [SAR-AUTOSETUP-001-OPS]
summary: "자동 설정 0.1.1 게시, relay main 배포, 서버 override를 repo compose로 대체"
attempt: 1bb285aba35e4ee9832219df4824868b
base: 1adcc7b9974e14f9f746bf3483c24777223acf76
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-OPS — 자동 설정 0.1.1 게시, relay main 배포, 서버 override를 repo compose로 대체

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-28
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D044 '운영 교체 진행'(2026-10-28). (1) 이 지시서 기준 main 커밋에서 `package_plugin.py --release`로 0.1.1을 만들어 서버 `knowslink-state/install/0.1.1/`에 게시(0.1.0 폴더는 지우지 않음). (2) 공개 0.1.1 세 파일 확인 뒤 relay를 같은 main 커밋으로 전환(과제 고유 백업·rollback 태그 `knowslink-rollback:bf67dfcb-autosetup011-20261028` 형식·writer stop·격리 복원·relay up). (3) 서버 checkout이 새 `compose.ops.yaml`(relay install env·:ro 마운트 포함)을 갖게 된 뒤 서버 로컬 `compose.install.yaml`을 D12 절차대로 빼고(개명 보관, 삭제 안 함) relay를 다시 띄운다. / 추가 승인 필요: 실계정 로그인·발송, 기존 volume·backup·0.1.0 파일 삭제, ingress·env·다른 서비스 변경, DB migration 발생, force-push.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_189d77c7-a247-4d47-bf52-920daaedab82 / run_1701a184c359 / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 relay `bf67dfcb`, 공개 `/install/0.1.0/`, 서버 로컬 `compose.install.yaml` override다(D12 맨 위 SAR-CONNECT-SIMPLE-001-OPS 절). main에는 그 뒤 FOLLOWUP-DEV(`compose.ops.yaml` install env·마운트, zod 4.6.7 patch), O04-DEV(설치기 O-04 수정, 0.1.1 상향, 화면 명령 0.1.1), tester QA3 PASS(O-04 해소, high 0, O-05: 0.1.0은 플러그인 3개 환경에서 첫 설치도 실패)가 들어 있다. 절차 후보는 D12 `### 자동 설정 0.1.1 교체 후보` 절과 FOLLOWUP-DEV가 적은 override 제거 절차다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing·security)이다. 기준 커밋 main(이 지시서 커밋). 제품 정본: `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md`. 결정 D043·D044는 `decision.py list`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/operations/ops-guide.md`(맨 위 운영 적용 절, `### 자동 설정 0.1.1 교체 후보`, override 제거 절차)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-OPS-test/report.md`(직전 배포 순서·명령)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-QA-test/report.md`(QA3 절)
- `.fullops-squad/docs/operations/transition.md`(rollback)

## 해야 할 일과 파일 소유권

- [ ] 제품 diff `bf67dfcb..<기준>`에서 schema·`db/`·의존성 변경을 확인한다(예상: zod patch만, migration 없음). migration이 있으면 멈추고 escalation.
- [ ] 0.1.1 release를 만들어 서버 `install/0.1.1/`에 두고, 공개 세 파일 200·설치기 안 SHA-256 일치·0.1.0 그대로를 확인한다.
- [ ] relay를 기준 커밋으로 전환한다(직전과 같은 순서). 회원 화면 `use=claude` 명령이 `install/0.1.1/`을 가리키는지 확인한다.
- [ ] 서버 checkout의 `compose.ops.yaml`로 install env·마운트가 들어오는지 `docker compose ... config`로 확인한 뒤 `compose.install.yaml`을 개명 보관하고 relay를 다시 띄운다. 이후 `/install/0.1.0/`·`/install/0.1.1/` 모두 200인지 본다.
- [ ] live 확인(실계정 없이): healthz, restart 0, 첫 기동·5분 뒤 error/panic 0, 기존 회원·agent·친구 보존 비교.
- [ ] D12를 현행으로 갱신한다(후보 절 → 적용 결과, override 절 → repo compose 기준 호출 규칙). D13 현재 운영 상태 갱신.

소유 파일: D12, D13, `qa-reports/SAR-AUTOSETUP-001-OPS-test/`, 이 인박스.

## 완료 기준과 검증

- 운영이 기준 커밋, healthy·restart 0·migrate 0, `/install/0.1.1/` 세 파일 200·SHA 일치, 화면 명령 0.1.1, override 없이 install 동작, 기존 데이터 보존. 실패하면 rollback 절차를 따르고 보고한다.
- 증명 방법: `qa-reports/SAR-AUTOSETUP-001-OPS-test/report.md`(front matter 필수, 명령·시각·종료코드, 비밀 없음).
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.10.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/ops --from <이 지시서 기준 SHA>` ERROR 0. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(배포). 화면은 명령 문자열의 버전만 확인한다.

## 갱신할 산출물

D12 ops-guide, D13 transition.

## 기대 산출물

운영 반영, D12·D13 갱신, OPS 보고서.

## 제약·협업·후속

하위 에이전트 배정 안 함(직렬 운영). 실계정 시험은 하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-OPS-packet.json` / SHA `1adcc7b9974e14f9f746bf3483c24777223acf76` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 75 · inferred · 필수
- `.fullops-squad/docs/operations/ops-guide.md` (document_update, document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 16, 137, 201 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 102 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 125, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_update, document_read) · 줄 15, 117, 121, 135, 145, 228 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 209, 244, 314 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md` (document_read) · 줄 125, 155, 180 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 18 · inferred
- `adapters/README.md` (document_read) · 줄 48, 55, 59, 105, 106, 134, 247, 300 · inferred
- `deploy/knowslink/compose.ops.yaml` (direct_edit) · 줄 전체/미확인 · inferred
- `docs/README.md` (document_read) · 줄 38, 38 · inferred
- `docs/operations.md` (document_read) · 줄 27, 33, 33 · inferred
- `docs/system-design.md` (document_read) · 줄 74, 303, 318 · inferred
- `scripts/package_plugin.py` (direct_edit) · 줄 25, 32, 100, 141, 164, 200 · inferred
나머지 참조 후보 50건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄, 상세는 보고서 링크.

- 브랜치 / SHA: `fullops/ops` / 기록 커밋 `e5bf4308`(기준 `1adcc7b9`; 운영 배포 제품은 main `c6dfc39ac4f6dabddc5864ad611437b84ded383e`).
- 변경 이유: 사용자 결정 D044. 자동 설정 0.1.1(O-04 수정)을 서버에 게시하고 relay를 `c6dfc39a`로 전환했으며, 서버 로컬 `compose.install.yaml`을 repo `compose.ops.yaml`로 대체했다. 서버 변경: checkout `bf67dfcb`→`c6dfc39a`, relay 이미지, `install/0.1.1/` 3파일, 새 백업 1개, override 개명 보관.
- 검토 필요: (1) 호스트 사용자 서비스 `cloudflared-orca`가 13:21:38Z에 이 과제와 무관하게 재시작됐다(다른 곳의 config 수정, 보고서 보존 절). (2) `compose.install.yaml.retired`가 서버에 남는다. 삭제하지 않았다. (3) 이전 지시서 문구의 live agent 30개는 32개로 늘었다(정상 증가).
- 검증한 것: 제품 diff에 migration·Go 의존성·Dockerfile 변경 0. 공개 `/install/0.1.1/` 세 파일 200·SHA-256 일치·설치기 안 ZIP 해시 일치, `/install/0.1.0/` 그대로 200(26/26 ×3회). relay healthy·restart 0·migrate 0, 중단 약18초(13:19:00Z~13:19:18Z), 첫 기동·5분 뒤 error/panic 0. 화면 명령 0.1.1은 실행 중 relay 바이너리 문자열로 확인(0.1.1 1건, 0.1.0 0건). override 없는 `config`==있는 `config`, 개명 후 relay 재생성 없음. 보존 8단계 비교 41/41.
- 못 한 것: 실계정 로그인·발송, 회원 화면 렌더링, 설치기 실세션 QA(AC-07·S5), 실제 rollback. 호스트 `make lint`는 Windows에서 실행 불가(unavailable, 문서·운영 기록만 변경).
- 산출물·로그: [보고서](../docs/evaluations/qa-reports/SAR-AUTOSETUP-001-OPS-test/report.md), 같은 폴더 `commands.txt`·`*.log`·`probe-*.json`·`preservation.json`. D12 `ops-guide.md`, D13 `transition.md` 갱신. 원본은 `.devtools/autosetup-001-ops/`.
- 남은 일: tester가 플러그인 2개 이상 환경에서 0.1.1 갱신 재실행(AC-07·S5)을 실세션으로 확인한다. 안정 확인 뒤 coor가 서버 `compose.install.yaml.retired` 정리 여부를 정한다(지금은 보관).
 worker_done 본문 첫 줄은 `[완료] SAR-AUTOSETUP-001-OPS | 브랜치 fullops/ops | SHA <브랜치 HEAD>` 한 줄로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
