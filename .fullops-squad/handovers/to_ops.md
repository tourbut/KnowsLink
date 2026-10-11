---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS — 독립 검수·main 수락 후 OAuth 브라우저 복귀 수정을 기존 서버에 보존 배포
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS]
summary: 독립 검수·main 수락 후 OAuth 브라우저 복귀 수정을 기존 서버에 보존 배포
attempt: f5a42297abfc45748fe3746cb2fc975e
base: eac603b6a41d2ce002e698600ac66b8fc67ad762
subagent_level: standard
test_level: lite
---
# SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS — 검수된 OAuth 복귀 수정의 기존 서버 보존 배포

- Purpose implementation(운영). 상태 ready: 독립 tester·designer 시각 수락과 main push 완료.
- ops C:/Users/shin/orca/workspaces/KnowsLink/ops, fullops/ops. fresh Claude Sonnet5.5 high(Jev). test lite/subagent standard.
- coor 복귀 term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. task/dispatch는 preamble을 따른다. 기록 main 병합·push는 coor 담당이다.

## 승인 범위와 보존

기존 .env.server Linux 서버에서 승인된 Google/Dots 연결의 OAuth 브라우저 복귀 오류만 보완 적용한다. 사용자가 기존 서버 배포와 테스트를 승인했다. 추가 승인을 반복 요청하지 않는다. 기존 비용 안에서 진행한다.

로컬 Docker·새 유료 API/서비스/billing·삭제/이동/cleanup/복구·force reset/clean은 금지한다. F: 사고 자료·private backup·DB·키·agent·pair·Google 설정·Tunnel credential·공유 서비스를 보존한다. 실계정 승인·외부 메시지·새 설치 생성은 부모가 배포 후 수행한다. 이번 worker는 수행하지 않는다.

## 원천과 적용 기준

배포 제품 SHA는 eac603b6a41d2ce002e698600ac66b8fc67ad762다. 실제 Chromium 독립 QA/고정 SHA 리뷰는 tester 기록96c7b823으로 수락됐다. designer 시각 기록610b2b0 및 별도 coor fixed-SHA 리뷰가 통과했다. main/origin aee03bcf0a039cd3918e4042d4041805ecd8c5d8에 세 완료 SHA 조상 관계와 push read-back을 확인했다. `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` 및 `-test/report.md`, `SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md`와 COOR 폴더 gate 영수증을 읽는다. F1 medium(원인별 실패 안내 부족)은 reviewer 비차단 잔여로 보존한다. 기존 미완료 install1은 삭제하지 않는다.

fullops-common-0.3.3과 FULLOPS, common README+세 규칙, project, orca-agents, delegation, document-writing을 읽는다. 원천은 현재 ops-guide/transition, SAR-SSE-INSTALL-001-OPS-LIVE, 이전 OPS 실행 기록과 logs/2026-10-11_to_ops.md, DEV 최종 OAuth 복귀 실행 기록 및 tester 수락 보고다. 고정 제품 SHA와 정확한 gate 경로는 선행 수락 뒤 coor가 연결한다. 기존 Context7/공식 CLI 근거는 동일 명령에 재사용한다.

운영 배포 /home/shin/deploy/knowslink, state /home/shin/deploy/knowslink-state, 도메인 https://link.knowslog.com. 마지막 배포 SHA 121cde923798c745af1f952347fcec1cae40b853. Tunnel config SHA256 f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512. 실제 착수 시 다시 읽어 현재값과 대조한다. 변경된 상태를 덮어쓰지 않는다.

SSH helper D:/workspace/KnowsLink/.git/fullops-gate/coor_remote.py와 기존 known_hosts RejectPolicy를 재사용한다. .env.server 값·code/token·메시지·DB 내용을 출력하거나 Git에 넣지 않는다. PowerShell literal stdin과 구조 인자로 전달한다. 원 helper·플러그인 캐시를 수정하지 않는다.

## 해야 할 일·소유권·완료 기준

1. 서버 clean·HEAD·서비스·env/credential/Tunnel hash·DB 집계 보존 기준을 읽기 전용으로 확보한다. 운영 설치1은 실제 사용자 승인 결과이므로 삭제하지 않는다. private 경로에는 원문, Git에는 비식별 요약만 보존한다.
2. 기존 state/backups에 고유 private DB backup을 생성하고 권한과 pg_restore --list를 확인한다. 이전 backup은 유지한다. 복구 실행은 하지 않는다.
3. 검수된 고정 SHA가 실제 origin/main 조상인지 확인한다. 서버 origin의 사라진 로컬 경로를 맹목적으로 fetch하지 않는다. 이전에 검증된 bundle/SFTP 방식으로 정확한 SHA를 전달하고 clean 배포에 안전하게 반영한다. migration이 있거나 서버 상태가 예상과 다르면 ask한다. 원 상태를 강제 덮어쓰지 않는다.
4. 기존 Compose project knowslink, compose.yaml+deploy/knowslink/compose.ops.yaml, state/.env와 UID/GID/STATE_DIR 계약을 유지한다. secret 출력 없는 config --quiet로 검증한다. relay만 최소 build/recreate한다. DB·Tunnel·다른 서비스·volume은 보존한다. 수정에 새 ingress 경로가 필요하지 않으면 config를 바꾸지 않는다. 제품 코드 수정이 필요하면 DEV에 반환한다.
5. local/public health, OAuth metadata, 무인증 MCP/SSE·polling 거부와 private 경로 404, 기존 공유 서비스의 이전 상태를 비교한다. 제품 전체 QA는 정확한 고정 SHA 증거를 재사용하고 반복하지 않는다. 실제 Google/Dots OAuth와 인증된 SSE/Events는 배포 후 부모 검증으로 남긴다. smoke 실패나 service degradation은 현재 상태·원인·담당·재개 조건을 보고한다. DB restore나 무조건 rollback은 하지 않는다.
   추가 비용/인증 없는 합성 43자 consent 경로 한 건으로 새 HTML 실패 안내·401·no-store/no-referrer·기존 CSP를 확인한다. 사용자 원 링크나 실제 attempt/code/token을 재사용하지 않는다. 새 설치를 만들지 않는다.
6. 자신의 exec plan·QA 요약·contexts/ops와 D12 ops-guide/D13 transition만 실제 변경 영향에 맞게 갱신한다. 제품 source와 다른 역할 증거는 수정하지 않는다. 새 script/config가 없으면 원 exact-SHA product gate와 Windows 실행 불가를 정직하게 구분한다. 새 실행 파일/설정은 해당 registered gate와 고정 SHA 리뷰를 받는다. 예상 문서 변경은 5파일 이내다.
7. 모든 packet 비선택 path/category와 partial/unknown을 수동 근거로 outcomes 처리한다. work.py finish로 원 지시서와 완료 보고를 보존한다. worker_done은 최종 기록40SHA 한 개와 운영40SHA·검증 결과·미실행을 분리해 보고한다.

## 지시 전제와 충돌 — 먼저 확인

탐색은 ops-guide와 이전 OPS 로그를 conflict로 분류했다. 이전 로그의 ba3754d 시작 상태와 ingress 변경 절차는 완료된 과거 작업이며 이번 작업에서 반복하지 않는다. 현재 기준은 OPS-LIVE의 121cde9 및 실제 read-only 착수 확인이다. OPS-LIVE 본문은 탐색 서비스 미전송이므로 로컬 원문을 직접 읽는다.

과거 ops-guide/transition의 held·Access·미적용 설명은 현재 OPS-LIVE와 구분한다. 최근 실제 배포는 위 121cde9이며 사용자 승인 이후 설치1/token0이 관찰됐다. 과거 수락을 새 OAuth 성공으로 간주하지 않는다. verify.py regression baseline 결함과 medium SSE 헤더 결함은 이미 보존된 잔여이며 이번 성공으로 바꾸지 않는다. 새 결함은 coor에 보고한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-packet.json` / SHA `bbba78bef50a35a3a3f572a22613ed63f1e36069` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 15, 15, 15, 20, 20, 20, 41, 41, 41, 43, 43, 43 · inferred · 필수 · {'relevant': 0.89, 'evidence': 0.87, 'contradicts': 0.23, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 2, 6, 6, 6, 10, 10, 10, 26, 26, 26, 27, 27, 27, 77, 81, 81, 81 · inferred · 필수 · {'relevant': 0.79, 'evidence': 0.77, 'contradicts': 0.17, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 7, 13, 13, 13, 13, 21, 27, 27, 27, 29, 29, 41, 50 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 6 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `deploy/knowslink/beta.sh` (impact_check) · 줄 전체/미확인 · unknown · 필수 · {'relevant': 0.18, 'evidence': 0.39, 'contradicts': 0.6, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `deploy/knowslink/compose.ops.yaml` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.48, 'evidence': 0.71, 'contradicts': 0.2, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `deploy/knowslink/verify.py` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.54, 'evidence': 0.79, 'contradicts': 0.41, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_auth.go` (impact_check) · 줄 전체/미확인 · unknown · 필수 · {'relevant': 0.2, 'evidence': 0.27, 'contradicts': 0.27, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 6, 66, 66, 66, 68, 68, 68 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 6, 6, 28, 28, 28, 28, 29, 30 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 34, 67 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 6, 6, 66, 66, 66 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 7, 14, 20, 23, 25, 115, 155, 163, 179, 179 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 14, 38 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 7, 14, 19 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 7, 14, 39 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 7, 14, 23, 28, 190 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 87 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 82 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 7, 14, 24, 55, 188, 199, 205 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 7, 14, 16, 19, 21, 21, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 43, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 23, 27, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 2, 6, 10, 17, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 72 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 2, 6, 6, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 15, 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 10, 10, 10, 12, 12, 12, 20, 20, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 99, 125, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 24, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 27, 33, 35, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 18, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 2, 6, 10, 17, 55, 56, 57, 59, 60, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 68, 72, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 36, 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 24, 53, 204, 221 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 6, 10, 38, 48, 55, 76, 82, 82, 86, 92, 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 2, 6, 10, 29, 55 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 2, 6, 6, 6, 10, 10, 10, 60, 60, 60, 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 28, 28, 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 14, 21, 21, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 2, 6, 10, 21, 22, 24, 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 2, 2, 6, 10, 16, 22, 27, 44, 48, 57, 72, 76 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update, document_read) · 줄 34 · inferred · {'relevant': 0.79, 'evidence': 0.83, 'contradicts': 0.72, 'injection': 0.13, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/operations/transition.md` (document_update, document_read) · 줄 15, 25, 108 · inferred · {'relevant': 0.19, 'evidence': 0.3, 'contradicts': 0.74, 'injection': 0.08, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 13, 22, 25, 57, 127 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 136 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 191 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 9, 12, 16, 20, 33, 39, 71, 104, 126 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 38, 41, 52, 54 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 43, 51, 57 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 60, 89 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 210 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 113, 157 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 115, 179, 210, 307 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 116, 158 · inferred
- `.fullops-squad/handovers/logs/2026-10-10_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-10_to_tester.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 9, 10, 13, 13, 17, 18, 25, 25, 34, 40, 49, 69, 77, 81, 143, 165, 165, 165, 166, 166, 166, 169, 169, 169, 169, 173, 173, 173, 174, 180, 180, 180, 180, 189, 195, 195, 195, 199, 210, 210, 210, 212, 212, 212, 238, 239, 241, 241, 241, 242, 242, 242, 243, 244, 245, 246, 247, 247, 247, 254, 259, 260, 261, 261, 261, 262, 262, 262, 263, 264, 302, 302, 302 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 9, 10, 14, 18, 26, 35, 54, 61, 80, 83, 85, 87, 91, 93 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 36, 36, 46, 48, 71, 83, 97, 103, 175, 175, 175, 176, 176, 176, 179, 179, 179, 179, 183, 183, 183, 184, 191, 191, 191, 191, 201, 201, 201, 207, 228, 228, 228, 230, 230, 230, 231, 231, 231, 286, 287, 313, 314, 315, 316, 342, 366, 367, 368, 369, 434, 434, 434, 440, 440, 440, 445, 445, 445 · inferred
- `adapters/README.md` (document_read) · 줄 6, 14, 24, 50, 70, 234 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14 · inferred
- `deploy/knowslink/tunnel/public-ingress.yml` (impact_check) · 줄 전체/미확인 · inferred
- `internal/config/config.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/device.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/testdata/consent_browser.mjs` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 146 · inferred
미확인 12건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

완료 전문과 최종 기록 SHA는 worker_done 전에 보존한다.
