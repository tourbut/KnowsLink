---
title: SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 베타 연결의 독립 QA
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-001-TESTER-PUBLIC]
summary: 최신 gate와 보호된 공개 베타 연결의 독립 QA
---

# SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 베타 연결의 독립 QA

- From / To: coor / tester. 상태 ready.
- 기록 checkout / 브랜치: /home/shin/orca/workspaces/KnowsLink/fullops-tester / fullops/tester.
- 복귀: run_8ca8bc058ab7, coor term_9afa8217-862c-404d-9a43-2122427113fc. 새 task/dispatch는 preamble을 따른다.
- 승인 범위: 읽기 전용 최신 고정 gate 검증, 외부 미인증 HTTP 검사, 공유 서비스 회귀, 자기 QA 기록 커밋. Access 쓰기·배포/DNS/connector 변경·제품 소스·실데이터·비밀 출력은 금지다.

## 기준과 읽을 문서

고정 설정 head28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2, basef824015다. readonly snapshot /tmp/knowslink-beta-review-28bd1bb 또는 실제 deployed checkout /home/shin/deploy/knowslink의 SHA를 확인한다. 상태 /home/shin/deploy/knowslink-state는 사용자0600 파일과 실제 cloudflared 설정을 읽는 데만 사용한다. 사용자 이메일과 토큰/credentials 값은 보고서나 Git에 넣지 않는다.

FULLOPS, fullops-common-0.3.2 README/세규칙, project.md 정본, 문서 작성 규칙과 fullops-test를 적용한다. Jev 결과 docs/evaluations/jev/SAR-BETA-001-TESTER-PUBLIC-*의 keep 전부를 읽는다: deploy/knowslink/{access_apply.py,verify.py,beta.sh,compose.ops.yaml,tunnel/config.yml.tmpl}; 이전 TESTER.md, REVIEW-FINAL-review/report.md; OPS phase; ops-guide; project; 공통 README/세규칙; FULLOPS; 현재 인박스. 이전 보고의 held는 역사적 조건이며 현재 통과로 바꾸지 않는다. 충돌/주의 후보는 없다.

## 해야 할 일과 선행 조건

- [ ] 최신 gate의 aud 불일치 일반/최적화 옵션 거부, proof ID 불일치·missing/stale/future proof 차단을 안전한 임시 상태에서 narrow 독립 확인한다. 제품 전체 QA/백업복원을 반복하지 않는다. PASS/실패 및 실제 SHA/명령 종료코드만 기록한다.
- [ ] gate 검증 완료 시 외부 연결을 열 수 있는지 coor 결정이 필요하므로 preamble ask로 gate 통과/실패와 SHA를 한 번 제출한다. ask는 public_ready 질문이다. 회신까지 기존 앱/DNS/connector를 변경하지 않는다.
- [ ] coor가 보호된 공개 연결 적용 완료를 회신한 뒤 https://link.knowslog.com의 미인증 root/owner/API/health/unknown 경로가 제품 200을 반환하지 않고 Cloudflare Access로 redirect/deny하는지 실제 확인한다. verify.py public의 구현과 기대 조건을 확인하여 실행한다.
- [ ] 실제 Access 앱/정책/IdP/단일 allow 사용자와 origin JWT aud/team 일치, no-bypass, DNS/별도connector 및 기존 orca Tunnel 보존을 로컬 비밀 자료의 값 출력 없이 확인한다. 새OAuth관리GET은 OPS가 보존한 access.live.json과 실제config를 대조한다. GET 재확인이 필요하면 coor에 요청한다.
- [ ] verify.py regression으로 orca/s8 200, mcp401 및 공유 cloudflared PID를 확인한다. 사용자 이메일 OTP 로그인은 인간 검사로 held하며 대신하지 않는다.
- [ ] 새 docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md와 관련 test/·scenarios 및 tester context/inbox/logs만 작성한다. lint/strict/whitespace 종료코드와 제품 diff 무변경을 확인한 뒤 finish/commit/worker_done한다.

## 수락 및 재사용

기존 local QA1762b43의 실행f824015와 제품78b1d92·QA659f4b0·리뷰311381f를 원래 조건으로 재사용한다. 제품과 Compose/verify/tunnel template 불변인 근거를 확인한다. 최신 gate와 실제 외부 보호만 새 실행으로 표시한다. 검사 소유자는 tester, 실제 리소스 적용은 OPS, 인간 로그인은 사용자다. 미해결 critical/high와 필수 실패는 수락을 차단한다. 캡처/영상은 자동 HTTP/설정 판정에 불필요하므로 만들지 않는다.

## 기대 산출물·완료 보고

D11/D12/D13은 읽기 전용 참조이며 새 QA 기록만 작성한다. 전체 제품·실제 벤더·실데이터/장기운영 held는 유지한다. 완료 body 첫 줄은 [완료] SAR-BETA-001-TESTER-PUBLIC | SHA <전체 보고 커밋>이다. 인간 로그인 미실행을 명시한다. 아래에 완료 전문을 작성한다.
