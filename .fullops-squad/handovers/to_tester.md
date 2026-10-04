---
title: SAR-MVP-003-BIDIRECTIONAL-TESTER — 고정 양방향 시험 검증
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-003-BIDIRECTIONAL-TESTER]
summary: 실제 시험 메시지 transport와 운영 준비의 독립 검증을 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-TESTER — 고정 양방향 시험 검증
- ready, tester 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester.
- coor 담당 병합 main/origin. 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7, task/dispatch preamble.
- 고정 대상 cd60e7f87eb5ce137eca887980f232b3f67a18d0, 변경 기준 f2849486ed48295e239714700e651d30c32f1c2c. 구현자 Codex 세션01a104d1-8cc5-7430-bcf1-1a2732c5183f. 실제Grok왕복/운영배포/Access인증은 아직미완료.
- FULLOPS.md, fullops-common-0.3.2 README와 연결3규칙, project.md, 문서작성규칙, DEV실행기록/D03D05D10/D12 13장 먼저 읽는다. Jev 관련 결과 keep/충돌 확인. 기존 owner-only보호는 유지하며 사용자 시험승인은 이번text왕복에 한정.
- 자기 기록만 수정. 제품/타인인박스/PLANS/board/기존원본리뷰 변경금지. 비밀값·실제 private trialstate는 열람출력/외부전달하지 않는다. 시험은 자기 격리자원. FullOps 업데이트 제외.
- self lint --from 착수HEAD exit0, work.pyfinish/archive/commit, worker_done final SHA 전체40자리.
## 할 일·완료
- [ ] fullops-test 적용. 별도cd60e7f detachedclone에서 시험시나리오 실행. 최소 추가하니스로 실제Postgres/두 독립MCP send→receive→reply→receive ID/from/to/text/held업무효과없음을검증. DEV스크립트만복사해독립증거대체금지. 기존verify-mvp 재사용가능하되 추가관찰assert와대상기록.
- [ ] 기본held no-network, trial crossagent/auth/key/signature/TTL/dup·idempotency, remote고정origin·redirect/loopback/10초 bodytimeout/64KiB, 실패cleanup/비밀출력없음, run_trial 0600 launcher와 Codex stdin/receive 명령, standalone배포bundle도구4개. 기술문서·Grok댓글 실제실행절차 검증.
- [ ] 기존Go/UI불변증거는 diff로의존동일성확인후재사용. UI변경없으므로 불필요Chrome복제금지. 실제Cloudflare 서비스/사용자credential/Grok계정 변경금지, 실제왕복미검증.
- [ ] docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md와 -test/에 명령exit/assert/남은외부검증정리. 새결함은고치지않고보고. 실패항목의의미를사후변경하지않음.

## Jev 목록·충돌
자기 key의 find/documents-find/context JSON keep 모두 읽는다. 지시 전제와 충돌 — 먼저 확인: DEV실행기록은 준비성공과 actual계정차단을 함께 적는다. 실제Grok왕복성공으로 확장하지 않는다. 신규QA/리뷰만 자기검증범위의근거다.
