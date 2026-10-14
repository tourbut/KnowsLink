---
title: SAR-REFACTOR-002-OPS — 검수된 Astra 리팩터링 운영 배포·백업·보존·health 검증
status: draft
updated: 2026-10-14
owner: ops
tasks: [SAR-REFACTOR-002-OPS]
summary: 검수된 Astra 리팩터링 운영 배포·백업·보존·health 검증
attempt: 3344417eb706415eaf6e4411c7e5fc35
base: 8a93e38e190f1b9e7b6bb183e578e410f5dd97c6
subagent_level: standard
test_level: lite
---

# SAR-REFACTOR-002-OPS — 검수된 Astra 리팩터링 운영 배포

- From/To: coor/ops. 상태 ready. Purpose: implementation.
- Task key SAR-REFACTOR-002-OPS. Test level lite. Subagent level standard(작은 직렬 운영은 직접 수행).
- Agent/model/effort codex/gpt-6-astra/high. 같은 사용자 지정 Astra 리팩터링의 배포 후속이다.
- 배포할 고정 제품 SHA: 8a93e38e190f1b9e7b6bb183e578e410f5dd97c6 (수락·main/origin 통합 완료).
- 담당 repo 0b08ec4c-9e3a-4613-8197-5a835e545335; C:/Users/shin/orca/workspaces/KnowsLink/ops; fullops/ops.
- 복귀 coor C:/Users/shin/orca/workspaces/KnowsLink/coor / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061 / run run_5534969da2d0. 실제 task/dispatch는 시작 preamble을 사용한다.
- 병합 책임 coor. 운영 결과 문서도 독립 fixed 리뷰 후 main 일반 통합/push.

## 승인과 현재 근거

사용자 최신 요청은 '배포해'다. 기존 link.knowslog.com의 relay에 고정 검수 제품을 실제 배포하고 private 백업·격리 복원·보존/health/smoke 검증·실패 복구까지 승인된 작업으로 끝까지 수행한다. 배포 직전 재승인은 필요 없다. 새 실계정 메시지 발송·agent 재등록/키 교체·친구 변경·Grok/Dots 호스트 원격 설치·다른 서비스 변경은 포함하지 않는다.

실제 운영 마지막 관측은 2026-10-13 제품2747fe98653d5c6946f41c219fc256395333d7fe다. 현재 SHA·프로세스·데이터를 새로 확인하며 과거 결과를 현재 상태로 주장하지 않는다. 동일 서버의 Orca/s8/mcp/host cloudflared는 공유 자원으로 보존한다. 현재 접속값은 기존 비공개 환경/승인된 helper에서만 메모리로 읽고 값은 출력하지 않는다.

수락 DEV3e5af5e와 QAe9b5cba가 배포 SHA8a에 포함됐다. 최종 소스/문서와 독립 QA 리뷰는 아래 정본이다. 최종 배포 SHA와 수락 후보의 실행 입력 동일성을 직접 확인해 기존 exact Linux lint/test/build·Postgres race·독립 QA를 재사용한다. 문서만 다른 SHA의 신규 실행 PASS로 바꾸지 않으며 제품 전체 gate를 이유 없이 반복하지 않는다.

## 적용 기준과 예외

AGENTS.md·FULLOPS.md·fullops-common-0.3.3의 README/코딩/테스트/보안·project.md·document-writing.md·D07/D12/D13·ops context를 직접 읽는다. 기존 사람 친구·agent 권한·writer guard·서명·Inbox/ACK·Events·24h 보관 경계를 보존한다. UI 외형/스키마/의존성 변경은 없다.

worker 삭제·이동 정리 금지를 유지한다. beta.sh deploy/backup/restore-verify는 내부 삭제가 있어 실행하지 않는다. 이전 운영의 무삭제 bundle→정확 SHA checkout→build→relay stop→새 private backup/list→network-none 격리 restore→기존 up 절차를 직접 확인해 재사용한다. 이번 승인된 배포에 필요한 대상 relay 컨테이너 교체와 고정 tracked checkout은 허용된 배포 범위다. 원본/backup·volume·clone·정본·자격 파일의 삭제/이동/덮어쓰기 정리는 허용하지 않는다. hook 차단은 원문/이유를 보존해 coor에 보고하고 다른 표현/스크립트로 우회하지 않는다.

Context7 Docker Compose resolve는 2026-10-14 quota 초과로 실패했다. 공식 https://docs.docker.com/reference/cli/docker/compose/up/ 와 https://docs.docker.com/reference/cli/docker/compose/stop/ 를 확인했다. up --wait의 health 대기·기존 mounted volume 보존, stop의 지정 service 비삭제 정지를 근거로 사용한다. 실제 서버 CLI 버전/help와 현재 compose 소스를 확인하고 적용한다. 비밀값은 외부 조회문에 넣지 않는다.

## 먼저 읽을 문서

- docs/operations.md
- docs/system-design.md
- .fullops-squad/project.md
- .fullops-squad/docs/design-docs/database-design.md
- .fullops-squad/docs/operations/ops-guide.md
- .fullops-squad/docs/operations/transition.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/commands.txt
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-INTEGRATION/report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-QA-REVIEW-review/report.md
- deploy/knowslink/beta.sh
- deploy/knowslink/verify.py
- deploy/knowslink/compose.ops.yaml
- deploy/knowslink/tunnel/public-ingress.yml
- README.md의 기존 Grok/Dots 업데이트 절

필요한 기존 보존/운영 명령만 연결해 읽는다. 과거 문서의 실제 자격이나 hostnames 전체를 출력하지 않는다. 새 인박스와 현재 packet을 전문 확인한다.

## 해야 할 일과 소유권

- [ ] 현재 서버 SHA·clean·relay/DB/tunnel health·공유 서비스·writer/schema·비공개 파일 권한을 새로 관측한다. 이전 소스와 배포 SHA의 runtime/config/migration diff를 확인하고 현재 정상복구 후보/이미지를 보존한다.
- [ ] 변경 전 기존 회원·agent 키/철회·친구/generation/suppression·Remote OAuth/설치/Inbox/Events·기존 설정·공유 서비스의 값 없는 count/digest를 기록한다. 만료 가능한 token/세션/메시지/구독은 현재 시간을 기준으로 분리하고 과거 자격을 되살리지 않는다.
- [ ] 같은 고정 SHA의 Git bundle을 고유 위치에 전송하고 hash/bundle verify/HEAD/clean을 확인한다. 현재 서비스가 살아 있는 동안 새 relay/migrate 이미지를 build한다. 원격 origin이 없으면 기존 안전 bundle 방법을 쓴다.
- [ ] relay와 모든 관련 writer만 정지하고 안정된 epoch를 확인한다. 기존 운영 DB에 pg_dump -Fc를 새 고유0700/0600 private 경로로 생성한다. hash·크기·pg_restore list를 확인한다. private dump/데이터/credential/본문/URL/secret/DB row 전문은 Git·터미널·보고서에 출력하지 않는다.
- [ ] network none·host port 없는 고유 격리 Postgres에 복원한다. schema/goose/HumanModel/writer guard·권한 상태를 확인하고 downgrade writer 거부를 검사한다. 실제 운영 DB에 거부용 쓰기를 시험하지 않는다. 검증 후 해당 fixture만 정상 stop하며 삭제하지 않는다.
- [ ] 새 migration 없음/호환을 확인한 뒤 현재 절차로 지정 relay를 기동하고 실제 이미지ID·고정 SHA·health·restart·짧은 error 관측을 기록한다. 권한/Google/Inbox/친구/기존 원격 설치와 자격 보존 및 공유 서비스 불변을 대조한다. 실패 시 쓰기를 멈추고 새 최신 백업을 보존한 뒤 같은 writer 규약의 직전 이미지/코드로 복구한다. 과거 DB로 운영 권한을 덮어쓰지 않는다.
- [ ] local readiness/DB health와 공개 HTTPS 경로 smoke·인증 거부·owner/test/probe 비공개·MCP/OAuth metadata를 확인한다. 실제 송신·receive·ACK·재등록·로그인 토큰 소비를 일으키는 검사는 실행하지 않는다. 공개 HTTP 입장/정리로 바뀌는 volatile 상태와 영구 권한 보존을 구분해 사유를 기록한다.
- [ ] D12/D13·ops phase·배포 QA·docs/operations.md를 실제 적용 SHA/시간/backup/복구 기준에 맞게 갱신한다. README의 이전 '이번 리팩토링' 설명이 새 동작과 충돌하면 현재 고정 SHA와 범위를 좁혀 고친다. Grok 로컬 Inbox 수정은 기존 prefix/bundle·watcher 갱신이 필요하고 Dots 서버 변경은 기존연결에서 적용됨을 구분한다. 사용자 계정에 대신 메시지를 보내거나 새 연결을 만들지 않는다.
- [ ] 모든 packet path/category와 uncertainty를 처리해 outcomes를 남기고 완료 보고 전문을 work.py finish로 보존한다. 코드/운영 절차 helper 변경 여부·기존 gate 입력 동일성/재사용·신규 실제 배포 검증을 구분한다. 고정 결과 커밋 후 worker_done 본문에 SHA: <40자리>·배포 제품 SHA·검증·남은 한계를 명시한다.

ops 소유: task-local 운영 helper/증거/phase, D12/D13, 관련 공개 운영/업데이트 문서. 제품 코드·테스트·공통 설정/guard·coor PLANS/board·다른 인박스는 수정하지 않는다. 운영 스크립트의 일반 기능 리팩터링은 범위 밖이다. 예상 변경은 관련 문서5~10개와 task-local 증거이며 새 dependency/운영 스키마는 없다.

## 완료 기준과 검증

worker 완료는 실제 배포8a93e38 + 새 private backup/list/격리복원 + 권한/기존 연동/공유 자원 보존 + 실제health/안전publicsmoke + 실패 원본/원격 종료코드 + 원인/복구/클라이언트 후속 문서 + 전문archive/빈inbox/worker_done이다. routine 배포 범위는 중간 확인 없이 수행한다. 현재 운영 drift·유출·무관 서비스 영향·비호환 DB·설명되지 않는 실패·hook 차단만 즉시 coor에 알리고 필요한 보호 상태를 유지한다.

coor는 운영 결과를 별도 실제 세션 fixed snapshot 리뷰하고 필요한 보존/health 증거를 대조한 뒤 문서 결과도 main/origin에 통합한다. 실제 계정 대화·장시간 수신·최대 부하·실제 rollback 성공은 새로 시험하지 않았으면 미확인으로 적는다. '배포 성공'을 모든 Grok 로컬 설치 업데이트나 자동 알림 화면 성공으로 확대하지 않는다.

## 기대 산출물

- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md와 redacted before/after/preservation/public-smoke·명령/종료코드·backup/restore 영수증
- .fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-OPS.md
- 영향받는 D12/D13·docs/operations.md 및 필요한 README 업데이트 안내
- packet-outcomes·전문 archive·빈 인박스·고정 결과 SHA

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-REFACTOR-002-OPS-packet.json` / SHA `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 7, 14, 18 · inferred · 필수 · {'relevant': 0.17, 'evidence': 0.46, 'contradicts': 0.34, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/commands.txt` (impact_check) · 줄 1 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md` (document_read) · 줄 6, 10, 19, 124 · inferred · 필수 · {'relevant': 0.19, 'evidence': 0.49, 'contradicts': 0.74, 'injection': 0.08, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md` (document_read) · 줄 6, 10, 18, 26, 46, 74, 97 · inferred · 필수 · {'relevant': 0.45, 'evidence': 0.65, 'contradicts': 0.52, 'injection': 0.09, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-INTEGRATION/report.md` (document_read) · 줄 2, 6, 12, 14, 24 · inferred · 필수 · {'relevant': 0.18, 'evidence': 0.33, 'contradicts': 0.48, 'injection': 0.09, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-QA-REVIEW-review/report.md` (document_read) · 줄 6, 10, 16, 24, 62 · inferred · 필수 · {'relevant': 0.53, 'evidence': 0.72, 'contradicts': 0.39, 'injection': 0.11, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 2, 6, 7, 14, 14, 17, 18, 48, 49, 50, 51, 52, 83, 84 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.59, 'evidence': 0.81, 'contradicts': 0.19, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `README.md` (document_read) · 줄 25, 36, 165 · inferred · 필수
- `deploy/knowslink/beta.sh` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `deploy/knowslink/compose.ops.yaml` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `deploy/knowslink/tunnel/public-ingress.yml` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `deploy/knowslink/verify.py` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `docs/system-design.md` (document_read) · 줄 6, 16, 17, 24, 279, 310, 318 · inferred · 필수 · {'relevant': 0.3, 'evidence': 0.5, 'contradicts': 0.24, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/tester-2026-10-14-refactor002.md` (document_read) · 줄 6, 80, 81, 83, 84 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28, 29 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 91, 93 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 28, 29 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 7, 14, 18, 133, 197 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 14 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 14, 22, 24, 30, 237 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 102 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.2-review/report.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-CODE-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md` (document_read) · 줄 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-FINAL-review/report.md` (document_read) · 줄 6, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-FINAL-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-RESULT-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-RESULT-review/report.md` (document_read) · 줄 1 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md` (document_read) · 줄 6, 10, 24, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CLEANUP/report.md` (document_read) · 줄 6, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CODE-REVIEW-review/report.md` (document_read) · 줄 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/adapter-report.md` (document_read) · 줄 2, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/remote-report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/report.md` (document_read) · 줄 2, 6, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/verified-qa-report.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-REVIEW-review/report.md` (document_read) · 줄 1 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 125, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-DEV.md` (document_read) · 줄 6, 10, 46, 52, 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-OPS.md` (document_read) · 줄 6, 10, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-DEV.md` (document_read) · 줄 6, 10, 19, 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-TESTER.md` (document_read) · 줄 2, 6, 10, 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 79 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update, document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_update, document_read) · 줄 23, 27, 41, 51, 134 · inferred · {'relevant': 0.24, 'evidence': 0.35, 'contradicts': 0.6, 'injection': 0.08, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 46, 81, 151 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 84 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 9, 10, 13, 17, 25, 45, 57, 66, 266, 311 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 9, 10, 13, 17, 25, 101, 105, 156, 157, 158, 159, 160, 197, 198, 245, 259, 267, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 36, 43, 62, 66, 67, 70, 74, 82, 105, 114, 117, 118, 119, 211, 212, 213, 298 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 49, 80, 81, 84, 84, 88, 89, 96, 96, 99, 100, 134, 147, 156, 157, 158, 164, 237, 238, 239, 240, 241, 242, 243, 244, 285, 286, 343 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_tester.md` (document_read) · 줄 9, 10, 13, 13, 17, 18, 25, 28, 29, 33, 33, 43, 44, 45, 46, 71, 72, 78, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 228, 229, 230, 298, 300 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 98 · inferred
- `adapters/README.md` (document_read) · 줄 82, 246 · inferred
- `docs/README.md` (document_read) · 줄 19, 20, 22, 36, 36, 37 · inferred
- `docs/architecture.md` (document_read) · 줄 36, 134 · inferred
- `docs/development.md` (document_read) · 줄 111, 158 · inferred
- `docs/operations.md` (document_update, document_read) · 줄 19, 19 · inferred · {'relevant': 0.82, 'evidence': 0.84, 'contradicts': 0.4, 'injection': 0.08, 'decision': 'keep', 'reason': None}
미확인 17건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

실제 적용제품 SHA·UTC/KST·중단시간·health·backup/복원·보존·공개검사·공유자원·재사용gate·실패/미실행·Grok/Dots후속·임시자료보존/정지·결과 SHA를 채운다.
