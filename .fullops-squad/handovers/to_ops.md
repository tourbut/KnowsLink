---
title: SAR-RELEASE-001-OPS — main 통합 제품을 기존 운영 relay에 배포하고 백업·보존·health·public smoke를 확인한다
status: draft
updated: 2026-10-17
owner: ops
tasks: [SAR-RELEASE-001-OPS]
summary: main 통합 제품을 기존 운영 relay에 배포하고 백업·보존·health·public smoke를 확인한다
attempt: a4b308ad7cbf4e55bc583722a80805d7
base: 1dd6800d5d4238e2060f99e98eaab437ea06682c
subagent_level: standard
test_level: lite
---

# SAR-RELEASE-001-OPS — main 통합 제품을 기존 운영 relay에 배포하고 백업·보존·health·public smoke를 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-17
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: 사용자 지시 "운영배포진행해"(2026-10-17). 기존 link.knowslog.com 서비스의 KnowsLink checkout·relay/migrate를 배포 SHA로 전환한다. 과제 고유 private backup·격리 restore fixture·작업 디렉터리를 만든다.
- 추가 승인이 필요한 행위: DNS·Cloudflare Access·계정 정책 변경, 삭제·prune·volume 삭제, 실계정 로그인·발송, pull/ACK/token 소비, 운영 DB를 과거 dump로 덮는 복구. 필요하면 `ask`로 멈춘다.
- 담당 repo id / 워크트리 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main·origin/main
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_e3a3e5f6-e163-420f-bba6-51b428f7671f / run·task·dispatch는 실제 preamble 값

## 현재 상황과 확인 근거

- 마지막 운영 배포는 SAR-REFACTOR-002-OPS의 8a93e38이다(참고값, fresh baseline으로 다시 확인).
- 배포 SHA는 main 1dd6800이다. 8a93e38 이후 제품 변경은 39파일 +2051/−121이다. 범위: 공개 첫페이지·설치 안내(SAR-LANDING-GUIDE-001), 조회 receipt·선택기(SAR-UX-ISSUES-002-RECEIPT), Claude Code 세션 메시지와 새 API `GET /v1/contacts/{agent}`·`POST /v1/text/seen`·receipt `session_received`(SAR-SESSION-PEER-001/002).
- `db/`·`deploy/` 변경은 없다. 새 migration이 없어야 한다. 실제 diff로 확인하고 다르면 보고한다.
- 각 기능은 독립 리뷰·tester QA를 거쳐 main에 통합됐다. SAR-SESSION-PEER-002 코드 리뷰는 사용자가 생략을 선택했다(waive 기록). 새 API 권한 경계는 tester가 직접 HTTP 20건으로 확인했다.
- 대기하던 SAR-LANDING-GUIDE-001-OPS는 미착수로 이 과제에 흡수했다(handovers/logs/2026-10-17_to_ops.md). 그 지시서의 운영 계약을 이 과제가 이어받는다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md` 및 연결된 코딩·테스트·보안 규칙이다. 기준 커밋은 main 1dd6800이다. AGENTS.md의 Windows 파일 작업 안전, project.md, D07/D12/D13을 유지한다. 예외 없음.

## 먼저 읽을 문서

- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md — 직전 배포 절차·원 이미지 조회 실패·교훈.
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/commands.txt, public-smoke.json.
- .fullops-squad/handovers/logs/2026-10-17_to_ops.md — 흡수한 첫페이지 OPS 지시서(공개 페이지 검사 항목).
- .fullops-squad/docs/operations/ops-guide.md, transition.md.
- .fullops-squad/contexts/ops.md.
- README.md, docs/operations.md.
- deploy/knowslink/compose.ops.yaml, deploy/knowslink/verify.py, deploy/knowslink/beta.sh(삭제 포함 절차는 데이터로만 읽고 실행하지 않는다).
- .fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md — 새 API 계약.

직전 private helper는 C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002-ops/에 있다. known_hosts/RejectPolicy·자격 메모리 사용·비밀 비출력 계약을 유지한다.

### 지시 전제와 충돌 — 먼저 확인

jev_context(SAR-RELEASE-001-OPS-context.json)가 충돌로 표시한 문서는 두 개다. 흡수한 첫페이지 OPS 지시서(logs/2026-10-17_to_ops.md)는 "QA 수락·main 통합 전 배포하지 않는다"를 전제로 한다. 그 선행 조건은 main 7551ed9 통합으로 충족됐고, 이번 사용자 지시로 배포한다. 그 지시서의 운영·보존 계약만 이어받는다.
verify.py의 Access302·projects baseline은 현재 404/401 계약과 맞지 않아 직전에도 미실행이었다. 과거 실패를 보존하고 현재 public-smoke 계약과 새 페이지·새 API 검사를 사용한다. 검사를 약화하지 않는다.

## 해야 할 일과 파일 소유권

- [x] fresh 운영 SHA·image ID·healthy·현재 요청/DB/공유 서비스 baseline을 읽기 전용으로 확인한다.
- [x] 배포 SHA bundle/hash·서버 clean checkout·source ancestor·runtime/DB/배포 입력 diff를 확인한다. migration이 있으면 멈추고 보고한다.
- [x] 새 build 전에 현재 실행 이미지에 과제 고유 rollback 태그를 붙이고 image inspect로 가용성을 확인한다. 기존 tag를 덮지 않는다.
- [x] relay가 정상일 때 새 이미지를 먼저 build한다. 전환 직전 단일 writer stop·새 private dump·네트워크 없는 고유 격리 복원과 digest 검증을 한다. 기존 backup/volume은 보존한다.
- [x] 미리 만든 이미지로 relay/migrate만 재개하고 health·restart·로그·schema·상태 digest·공유 컨테이너 보존을 확인한다. 실패하면 rollback 태그 이미지로 서비스를 복구하고 원 실패를 남긴다. 운영 DB를 과거 dump로 덮지 않는다.
- [x] 기존 안전 public smoke와 새 공개 첫페이지(소개·설치 details·로그인 provider·무인증 회원 경계·보안 헤더)를 실제 HTTPS에서 확인한다.
- [x] 새 API의 무인증 경계를 확인한다: `GET /v1/contacts/{agent}`·`POST /v1/text/seen`가 인증 없이 401/403이다. 실계정 토큰은 쓰지 않는다.
- [x] D12/D13·README 현재 운영 기준·docs/operations·own QA/context를 갱신한다. 제품 source·설계·다른 역할 문서는 수정하지 않는다.
- [x] 실행 명령 자신의 종료코드, UTC/KST 전환 시각·중단 시간·새 백업 mode/hash·이미지·보존 assertions·원 실패·미실행·복구 한계를 남긴다.
- [ ] work.py finish·고정 결과 SHA·worker_done으로 회신한다.

소유권: `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/`, `.fullops-squad/contexts/ops.md`, 이 인박스, D12/D13 원천, README 운영 절, docs/operations.md, .fullops-squad/docs/operations/. runtime 제품 변경이 필요하면 멈추고 보고한다(DEV 인계). main 병합·push·현황판은 coor 담당이다.

## 완료 기준과 검증

- 운영이 1dd6800 이미지로 healthy이고 restart 0, 기존 데이터·Events·공유 서비스 보존 assertion이 통과한다.
- 새 private 백업과 격리 복원 검증이 남는다. rollback 태그 이미지가 inspect된다.
- public smoke 전부와 새 첫페이지·새 API 무인증 경계 검사가 통과한다.
- 제품 코드 gate는 각 기능의 수락 증거(main 통합 기록)를 재사용한다. 이번 운영 검사와 구분해 적는다. lint.py는 결과 문서 HEAD에서 ERROR 0이어야 한다.
- 실계정 로그인·메시지 왕복은 이번 범위 밖이다. 미실행으로 적는다.

### UI 작업의 디자인 기준과 검증

해당 없음. 운영 배포이며 UI 코드를 바꾸지 않는다. 공개 페이지는 배포 후 HTTPS에서 내용·헤더만 확인한다.

## 갱신할 산출물

D12, D13, D10(route 추천). D10은 운영 변경이 실제 영향을 줄 때만 갱신한다.

## 기대 산출물

`.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/`(report.md·commands.txt·baseline/after/preservation/public-smoke JSON·로그), 갱신한 D12/D13·운영 문서. 비공개 자원은 AppData/Local/KnowsLinkDevTools/release001-ops 및 서버 knowslink-state 아래 과제 고유 디렉터리에 둔다. credential·token·cookie·실메시지·DB row는 Git·터미널·보고서에 출력하지 않는다.

## 제약·협업·후속

운영 전환은 직렬 처리한다. 하위 위임은 하지 않는다(직렬 운영 작업). 기존 사용자 소유 터미널을 종료하거나 작업을 주입하지 않는다. Docker/Compose 버전이 직전과 다르면 Context7 또는 공식 문서로 옵션을 확인한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-RELEASE-001-OPS-packet.json` / SHA `4d3e54c1e4210809829d945a65dbfb73d4620fe2` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 73 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md` (document_read) · 줄 54, 61, 70, 98 · inferred · 필수
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 2, 2, 2, 6, 7, 7, 7, 14, 14, 14, 14, 25, 31, 31, 38, 43, 55, 55, 56, 64, 65, 71, 71, 77, 78, 91, 91, 98, 104, 106, 109, 111 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 15, 81 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `docs/operations.md` (document_read) · 줄 6, 19, 20, 33, 39, 48, 128 · inferred · 필수
- `.fullops-squad/contexts/archive/tester-2026-10-14-refactor002.md` (document_read) · 줄 28, 59, 71 · inferred
- `.fullops-squad/contexts/coor.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 63, 75, 89, 103, 105 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 117 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 31, 101, 229 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 34, 197, 216, 231 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 66, 85, 96, 97, 227, 230, 273, 279, 284, 292 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 25, 29, 118 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 86, 106, 166, 235, 244, 258, 270 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 68, 74, 91, 97 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 32, 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.3.0-FINAL-review/report.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.0-review/report.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.2-review/report.md` (document_read) · 줄 14, 29, 30, 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.4-FINAL-REVIEW-review/report.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 49, 95 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 73, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 32, 74, 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 20, 60, 75, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 50, 59, 65, 65, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/report.md` (document_read) · 줄 94 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 29, 33, 56, 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-QA-DOCS-review/report.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 30, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-test/report.md` (document_read) · 줄 40, 49 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 16, 83 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 24, 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV/report.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 12, 14, 16, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 16, 74, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md` (document_read) · 줄 115, 196 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/v01-v02-report.md` (document_read) · 줄 44 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-LIVE/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CSS-REVIEW-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume/report.md` (document_read) · 줄 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 14, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/report.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/web-20261014-232426-595664/report.md` (document_read) · 줄 7, 12, 13, 14, 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/dev-final-handoff.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/report.md` (document_read) · 줄 12, 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md` (document_read) · 줄 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 22, 151 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 49, 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FINAL-RECORDS-review/report.md` (document_read) · 줄 40, 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 53, 99, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 116 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/scenario.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 38, 89, 90, 93, 148, 164 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 24, 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 15, 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 29, 114 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-CODE-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-FINAL-review/report.md` (document_read) · 줄 22, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-FINAL-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-RESULT-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md` (document_read) · 줄 28, 67, 76, 80, 86, 122 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md` (document_read) · 줄 34, 70, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CLEANUP/report.md` (document_read) · 줄 14, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/report.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md` (document_read) · 줄 14, 57, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/verified-qa-report.md` (document_read) · 줄 12, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-INTEGRATION/report.md` (document_read) · 줄 7, 34, 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS-REVIEW-review/report.md` (document_read) · 줄 39, 59 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-QA-REVIEW-review/report.md` (document_read) · 줄 54, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md` (document_read) · 줄 12, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/report.md` (document_read) · 줄 2, 6, 10, 18, 20, 67, 77, 83, 83, 98, 102, 104 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REPO-DOCS-001-DEV/report.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REPO-DOCS-001-REVIEW-review/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-LIVE/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 20, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 45, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 34, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 28, 66, 72, 107, 109, 121 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-2-review/report.md` (document_read) · 줄 25, 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-3-review/report.md` (document_read) · 줄 20, 32, 59, 69, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DESIGN/report.md` (document_read) · 줄 35, 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-DEV-test/report.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/report.md` (document_read) · 줄 37, 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-FINAL-review/report.md` (document_read) · 줄 23, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 7, 13, 34, 45, 63, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-REVIEW-review/report.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 26, 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 7, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 84 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md` (document_read) · 줄 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-REVIEW-review/report.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 12, 35 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 48, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-FIX-REVIEW-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 49 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 38, 50 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.12.md` (document_read) · 줄 14, 22, 32, 36, 38, 42, 48 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.13.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 39, 44, 49, 51, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 15, 41 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 32, 57 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.3.0.md` (document_read) · 줄 34, 40, 42, 46, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.0.md` (document_read) · 줄 28, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.2.md` (document_read) · 줄 14, 16, 27, 33, 39, 49 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 18, 38, 54, 72, 76, 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.6.0.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md` (document_read) · 줄 14, 116 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 14, 42, 46, 97, 115 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-DEV.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 23, 30, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 41 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 22, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 16, 21, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-TESTER.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 52, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 34, 40, 48 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-VISUAL.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 18, 112 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md` (document_read) · 줄 42, 51 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 59, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 170 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 47, 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 14, 35, 46, 80, 223 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 14, 38, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md` (document_read) · 줄 50, 50, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 22, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 88 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 40, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 2, 10, 14, 18, 21, 22, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 30, 169, 171 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-DEV.md` (document_read) · 줄 43, 44, 64 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-OPS.md` (document_read) · 줄 22, 23, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-OPS.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REPO-DOCS-001-DEV.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15, 69, 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 30, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DEV.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-TESTER.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 28, 31, 40, 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 23, 32, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 23, 49, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SYSTEM-DIAGRAM-001-DEV.md` (document_read) · 줄 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 28, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 66 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 19, 19, 31, 31, 43, 61, 135, 143, 149, 152 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 203 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14, 22 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 14, 16 · inferred
- `.fullops-squad/docs/planning/SAR-SETUP-001-request.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 38, 125 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 134 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 41, 49, 140, 147 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 28, 104, 110 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 56, 68, 80, 84, 94, 115, 189, 213, 215, 219 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 27, 90, 149, 196, 197, 259, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 24, 72, 161, 197, 229, 237, 237, 253, 254, 269, 270, 296, 299 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 24, 75, 125, 167, 207 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 93, 126, 156 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 43, 93, 121, 140, 154 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 22, 39, 64, 80, 87, 105, 142, 159 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 76, 85, 86 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 32, 235, 257, 276, 284 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 28, 77, 173, 203, 232, 254, 283, 336, 384 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 97, 142, 162, 184, 207, 240, 252, 254 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 222, 226, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 37, 88, 92, 111, 261, 265, 267, 309 · inferred
- `.fullops-squad/handovers/logs/2026-10-12_to_ops.md` (document_read) · 줄 32, 36, 46, 49, 49, 79 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 53, 230, 232, 270, 294, 300 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 35, 265, 271, 271 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 64, 356, 360, 362, 403 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 32, 37, 93, 191, 219, 257, 343 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 29, 49, 49, 147, 355, 377, 380, 516, 520 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 32, 72, 134, 156, 166, 205 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 304 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 71, 89, 103 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 30, 34, 51, 66, 72, 73, 84, 88 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 7, 12, 13, 22, 26, 28, 32, 36, 40 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 32, 66, 129, 148 · inferred
- `README.md` (document_read) · 줄 36, 123, 145, 154, 202 · inferred
- `adapters/README.md` (document_read) · 줄 46, 75, 90, 290, 298, 363, 386 · inferred
- `adapters/hooks/register.test.ts` (impact_check) · 줄 57, 91, 245 · inferred
- `adapters/hooks/register.ts` (impact_check) · 줄 60, 265 · inferred
- `adapters/manifest.json` (impact_check) · 줄 24 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 12, 40, 44 · inferred
- `adapters/src/inbox.ts` (impact_check) · 줄 268 · inferred
- `adapters/src/index.ts` (impact_check) · 줄 6, 25 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 25, 41, 43 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 27, 32, 71, 83, 94, 135, 161 · inferred
- `adapters/src/public-check.ts` (impact_check) · 줄 15, 18, 23, 29, 47, 59, 68, 84 · inferred
- `adapters/src/trial-check.ts` (impact_check) · 줄 12, 113 · inferred
- `adapters/src/trial-cli.ts` (impact_check) · 줄 3, 33 · inferred
- `adapters/src/trial-setup.ts` (impact_check) · 줄 50, 78, 124 · inferred
- `cmd/migrate/main.go` (impact_check) · 줄 2, 18 · inferred
- `cmd/migrate/main_test.go` (impact_check) · 줄 2 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 2, 19 · inferred
- `cmd/relay/main_test.go` (impact_check) · 줄 2 · inferred
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 26, 49, 91 · inferred
- `deploy/knowslink/test_public_config.sh` (impact_check) · 줄 2, 18, 19, 25, 32, 35, 38 · inferred
- `deploy/knowslink/verify.py` (impact_check) · 줄 1, 63, 66, 79 · inferred
- `docs/development.md` (document_read) · 줄 16, 75, 152 · inferred
- `docs/system-design.md` (document_read) · 줄 24, 73, 75, 77, 79, 126, 225, 294, 304 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 29, 32, 303, 304 · inferred
- `internal/relay/device.go` (impact_check) · 줄 1, 21, 26, 120, 121 · inferred
- `internal/relay/human.go` (impact_check) · 줄 129 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 8, 16, 17, 18, 23, 24, 27, 28, 31, 32, 34, 35, 38, 39, 41, 42 · inferred
- `internal/relay/human_pages_test.go` (impact_check) · 줄 52 · inferred
- `internal/relay/landing_integration_test.go` (impact_check) · 줄 54, 62 · inferred
- `internal/relay/landing_pages.go` (impact_check) · 줄 7, 11, 15, 40, 46, 52, 54, 75 · inferred
- `internal/relay/landing_test.go` (impact_check) · 줄 1, 26, 34, 36, 41 · inferred
- `internal/relay/member_agents.go` (impact_check) · 줄 61, 62, 210, 211, 330, 331 · inferred
- `internal/relay/member_pages.go` (impact_check) · 줄 49, 54, 57, 58, 59, 65, 66, 87, 88, 96, 97, 110, 112, 116, 148 · inferred
- `internal/relay/protocol.go` (impact_check) · 줄 264, 266, 289, 290 · inferred
- `internal/relay/protocol_test.go` (impact_check) · 줄 48, 51, 61 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 54, 85, 110, 123, 210, 224, 243 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 5, 7, 8, 15, 20 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 77 · inferred
- `scripts/mail_sink.py` (impact_check) · 줄 38, 53 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 123, 169, 175, 197 · inferred
- `scripts/run_trial.py` (impact_check) · 줄 9, 33 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 213, 234 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 189 · inferred
- `scripts/verify_grok_plugin.py` (impact_check) · 줄 13, 42 · inferred
- `scripts/verify_setup.py` (impact_check) · 줄 23, 36, 53 · inferred
미확인 22건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
제품 `1dd6800d5d4238e2060f99e98eaab437ea06682c`를 기존 link.knowslog.com 운영 relay에 배포했다. 준비 HEAD는 `2ef95cb`, 브랜치는 `fullops/ops`다. 상세와 원 로그는 [운영 보고](../docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/report.md)와 같은 폴더의 JSON/log다.

- 전환: STOP `2026-10-17T09:37:48Z`(18:37:48 KST) → UP `09:38:09Z`, 약21초. relay `sha256:4893d0f5…0098999` healthy·restart0, migrate EXIT0, 새 migration 없음. 제품 diff는 지시서와 같은 39파일 +2051/−121이고 db/deploy 입력 diff는 0이다.
- 입력: bundle SHA256 `f8e1b2de…aba6c60`, 서버 hash·bundle verify·FETCH_HEAD·8a93e38 조상 확인. 서버 checkout은 `1dd6800` clean.
- rollback: build 전에 이전 relay 이미지 `sha256:482d3197…381343d`에 `knowslink-rollback:8a93e38-release001-20261017` 태그를 붙여 inspect했다. network-none read-only로 relay/migrate 파일과 SHA256도 확인했다. 실제 rollback은 실행하지 않았다.
- 백업: `backups/8a93e38-pre-release001-20261017T093752Z.dump`, 0600·11104bytes, SHA256 `8498f8eb…ca89ea3452`, TOC25줄. frozen epoch1297596·다른 DB client0에서 생성했다. network-none 격리 복원의 전체 data digest가 frozen 값과 같고 goose2·guard1·HumanModel1·downgrade 거부를 확인했다. fixture `knowslink-release001-rehearse-20261017`(exited0)을 보존했다. 운영 DB는 복원하지 않았다.
- 보존: [preservation.json](../docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/preservation.json) assertion 38/38 통과. Rates만 volatile 차이. 기존 백업16개·공유 컨테이너·Postgres·Tunnel credential·파일 mode 보존.
- 공개 smoke: 기존 21/21, 새 첫페이지(소개·설치 details 2·Google 로그인 폼·보안 헤더·금지 문자열 없음) 통과. 공개 응답의 Cloudflare Web Analytics script 1개는 edge 주입이고 origin에는 없다(관찰만).
- **지시서와 다른 발견**: 새 API `GET /v1/contacts/{agent}`·`POST /v1/text/seen`이 공개 주소에서 404였다(loopback은 401). 공개 Tunnel `path` 정규식에 두 경로가 없었다. coor `ask`에 사용자 결정 "repo+live 함께 반영"을 받아 두 경로만 추가했다: repo `deploy/knowslink/tunnel/public-ingress.yml`·`scripts/check_public_ingress.py`, live `config.yml`(원본 백업 SHA256 `f5abe860…458512`, 0600). `cloudflared tunnel ingress validate`/`rule` 통과 후 2026-10-17T09:44:36Z에 cloudflared를 재시작했다. 끊김 약2초(502/530 → 200), 연결 4개 등록·새 프로세스 error0이다. 이후 두 경로 무인증 401(%3A 인코딩 포함), 이웃 경로 6개 404, live와 repo 규칙 동일(SHA256 `f97e0fd9…cb7`).
- 검사: 로컬 `check_public_ingress.py` 통과. Windows `test_public_config.sh`는 NTFS 파일 mode(644≠600)로 실패했다. 같은 파일을 `git archive 1dd6800` 사본에 덮고 서버 Linux에서 두 검사 모두 EXIT0이다.
- 실패·사고 기록: (1) ingress 적용 스크립트 첫 시도가 `grep -m1` SIGPIPE로 종료141, live 변경 전이었다. (2) 내 관찰 루프(약2.5회/초)가 09:44:51Z부터 relay 경유 공개 경로를 약2분 429로 만들었다. 같은 시각 로컬 PC도 429였고 09:47:05Z에 회복했다. 다른 사용자 영향은 알 수 없고 단위(IP별/공유)는 확인하지 않았다. (3) post-check 재실행의 assertion은 의도한 cloudflared 재시작의 종료 ERR 5줄을 세어 실패했다. 새 프로세스 기준 0건이다.
- 미실행: 실계정 로그인·메시지 왕복, 인증된 contacts/seen 성공 경로, pull/ACK/token 소비, 실제 rollback, Access/DNS 변경, `verify.py regression/public`(과거 전제 불일치), `beta.sh deploy/backup/restore-verify`(삭제 포함). 제품 코드 gate는 각 기능의 수락 증거를 재사용한다(이번에 재실행하지 않음).
- 문서: D12 `ops-guide.md`·D13 `transition.md`·README 운영 SHA·`docs/operations.md`·`contexts/ops.md`를 갱신했다. D10은 갱신하지 않았다. `docs/README.md`가 8a93e38을 가리키는 오래된 포인터는 ops 소유가 아니라 수정하지 않았다(coor/DEV 후속).
- 후속(coor/DEV): ingress 수정 커밋 `21a6e87`은 main 병합 뒤 서버 checkout에 도착한다(live는 이미 반영). 새 공개 경로 변경은 ingress 규칙·검사를 함께 바꾸는 gate가 필요하다. Grok 로컬 bundle과 Claude Code hook 패키지는 서버 배포로 갱신되지 않는다.
- 비공개 자원: 로컬 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/release001-ops/`, 서버 `/home/shin/deploy/knowslink-state/release001-ops-20261017/`. 삭제·이동·prune 없음.
