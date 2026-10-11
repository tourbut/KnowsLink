---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER — OAuth 승인 후 실제 Chromium callback과 token교환 독립 수락 검수
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER]
summary: OAuth 승인 후 실제 Chromium callback과 token교환 독립 수락 검수
attempt: 0a2d62611ca542a5b77e279c3c4f1734
base: c399a5e52c714d4ee61a7cbca03d32ffbfed9bcf
subagent_level: standard
test_level: lite
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER — OAuth 승인 callback 독립 검수

- Purpose review/QA. DEV 최종 SHA eac603b6a41d2ce002e698600ac66b8fc67ad762. Test lite, subagent standard.
- tester C:/Users/shin/orca/workspaces/KnowsLink/tester/fullops/tester. fresh Claude Sonnet5.5 high(Jev).
- coor 복귀 term_e61d3e14-29e9-4954-943a-4a75707c82de/run_dc9b4e6dd18a. task/dispatch preamble 참조.
- 비교 기준 c399a5e52c714d4ee61a7cbca03d32ffbfed9bcf. 구현자 실제 세션 e20cc179-625e-43fa-84fd-b431d432b9c6과 다른 실제 세션으로 검수한다. Review SHA는 위 DEV 최종40SHA이며 coor 준비 merge의 HEAD를 제품 리뷰 SHA로 오인하지 않는다.
- 사용자 승인: 기존 Linux의 격리 테스트, 기존 비용 안의 수락 검증. 운영 배포·실계정 Google/Dots조작·새비용·삭제/이동/cleanup은 금지한다.

## 요구와 원천

DEV 실행 기록 `docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md`와 driver/test source를 먼저 읽는다. 정확한 최종 SHA의 등록 gate는 `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/dev-eac603b6-final-gate.json`과 pass 영수증이다. product-lint/test 각각0, ERROR0/WARNING3/unavailable0이며 첫 ANTI-004 실패는 build tag 분리로 수정됐다. Linux Chrome148 합성 provider/callback에서 이전 handler 실패·수정 성공을 DEV가 재현했다. independent 실행으로 확인한다.

기존 격리 helper `C:/Users/shin/AppData/Local/KnowsLinkDevTools/oauth001/remote.py`는 읽기 참고만 한다. DEV의 `.sh`에는 자기 임시 PG에 대한 `docker rm -f`/trap 정리가 있었으므로 그대로 실행하지 않는다. 본인 고유 helper·고유 경로/PG를 준비하고, 삭제/이동/cleanup/trap 제거 없이 보존한다. 기존 운영 서비스 목록 전후 동일 여부를 확인한다. 운영 DB·실계정은 사용하지 않는다.

실제 Dots remote consent 승인 후 authorization_failed. 운영 집계 install1 생성/token0, code1분 만료와 attempt소멸. 유력 원인은 승인POST→외부callback303와 CSP form-action self+Google 충돌이다. 최초실패가 만료인지 callback차단인지는 실제 Chromium 회귀로 검증한다. DEV handover archive와 DEV 최종exec plan·QA를 읽고 실제 변경이 요구를 충족하는지 검수한다.

먼저 FULLOPS, common README+3규칙, project, orca-agents, delegation, document-writing, review rule/lint계약을 읽는다. 기술원천은 remote_auth.go/google.go/http.go/capacity.go와 관련 HTTP/단위테스트, SSE-INSTALL-001-REMOTE 계획·OPS-LIVE·이번 DEV 실행 기록이다. 과거 수락 증거와 새브라우저증거를 분리한다.

## 완료 기준·소유권

1. 구현자와 다른 실제 세션, 깨끗한 별도 detached fixedSHA 읽기 전용 snapshot에서 fullops-review/OCR delegate 검토. snapshot은변경하지 않고 보고·QA만 tester정본에 쓴다. source 수정금지. 발견은 coor에 전달.
2. Chromium의 합성 Google/provider·registered host callback/isolatedDB로 승인POST부터 callback도착과 token교환 성공을 직접 확인한다. 기존 실제계정/운영DB에는 접속하지 않는다. 원 code/token/credential을 기록하지 않고 codepresent/교환status 등으로 요약한다. DEV의 재현경로를 재사용할 수 있지만 독립 실행과 대상 SHA를 명시한다.
3. 허용callback/PKCE/resource/issuer/CSRF/Google account binding·expired/reused consent의 실패 경계와 안전한 사용자 복귀안내를 검수한다. 전역wildcard CSP·CSRF/SameSite완화·로그인생략·새토큰노출이 없는지 확인한다. 공통middleware변경이면 영향Google/SSE회귀를 확인한다.
4. DEV exact finalSHA 필수registered lint/test 근거를 검증한다. 바이트동일/영향없는제품 검사를 불필요하게 반복하지 않고 독립핵심 QA에 집중한다. 필요한 새 lint와 review.py check의 head를 일치시킨다.
5. 기존medium SSE헤더·verify.py baseline결함을 새결함과 분리한다. 미해결C/H 차단. 운영계정·실제DotsEvents미검증을 성공으로 바꾸지 않는다.
   반환 화면의 자동 복귀와 링크가 desktop/mobile에서 사용 가능하고, 오류 화면에 재시작 안내가 있는지 캡처한다. 직접 시각 수락은 designer 담당이며 캡처와 관찰 근거를 넘긴다.
6. owns 자신의exec plan/QA/review/scenarios/contexts/tester만. 원문·failed·skipped보존, partial/unknown과packet모든비선택path/category outcomes처리. work.py finish 후 명확한최종40SHA worker_done. 부모가 main수락/push·OPS배포·실계정연결을 맡는다.

## 지시 전제와 충돌 — 먼저 확인

탐색 결과는 `docs/operations/ops-guide.md`를 conflict로 분류했다. 과거 운영 가이드의 배포·검증 절차보다 이번 지시서의 격리 검수 범위를 우선한다. 현재 배포 상태는 OPS-LIVE 기록을 참조한다. 과거 가이드의 수락 기록을 이번 브라우저 검증의 성공 근거로 사용하지 않는다. 운영 변경이 필요하면 직접 적용하지 않고 coor에 보고한다.

Google 구현과 일부 테스트·OPS-LIVE 본문은 탐색 서비스에 전송되지 않았다. 로컬 원문을 직접 읽고 확인하며, 미전송을 검증 완료로 해석하지 않는다. DEV 최종 SHA를 받은 뒤 변경된 원천과 패킷을 다시 연결한다.

## 완료 보고

완료 후 전문을 기록하고 work.py finish로 보존한다.
