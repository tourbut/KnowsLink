---
title: SAR-MVP-003-BIDIRECTIONAL-REVIEW — 고정 양방향 시험 검증
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-REVIEW]
summary: 실제 시험 메시지 transport와 운영 준비의 독립 검증을 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-REVIEW — 고정 양방향 시험 검증
- ready, ops 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops.
- coor 담당 병합 main/origin. 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7, task/dispatch preamble.
- 고정 대상 cd60e7f87eb5ce137eca887980f232b3f67a18d0, 변경 기준 f2849486ed48295e239714700e651d30c32f1c2c. 구현자 Codex 세션01a104d1-8cc5-7430-bcf1-1a2732c5183f. 실제Grok왕복/운영배포/Access인증은 아직미완료.
- FULLOPS.md, fullops-common-0.3.2 README와 연결3규칙, project.md, 문서작성규칙, DEV실행기록/D03D05D10/D12 13장 먼저 읽는다. Jev 관련 결과 keep/충돌 확인. 기존 owner-only보호는 유지하며 사용자 시험승인은 이번text왕복에 한정.
- 자기 기록만 수정. 제품/타인인박스/PLANS/board/기존원본리뷰 변경금지. 비밀값·실제 private trialstate는 열람출력/외부전달하지 않는다. 시험은 자기 격리자원. FullOps 업데이트 제외.
- self lint --from 착수HEAD exit0, work.pyfinish/archive/commit, worker_done final SHA 전체40자리.
## 할 일·완료
- [x] fullops-review/open-code-review-delegate 적용. Jev sonnet-only 추천보다 인증/데이터변경의 high-performance 규정을 우선하여 fresh claude-opus-5-5 high로 리뷰. 실제reviewer_session과 구현자 서로다름증명.
- [x] read-only clean detached /tmp/knowslink-bidirectional-review-cd60e7f는 대상SHA 고정. 실행검증은 별도scratchclone, snapshot쓰기금지.
- [x] 준비된 review key SAR-MVP-003-BIDIRECTIONAL의 preview/rules/result 전체path/status 검토. 원본문서와 실제구현/테스트/댓글/Access생성계획 정합을 대조. trial intent와 기본held/business 경계, two-agent allowlist/TTL/서명/원문수명/비밀소유권/SSRF/응답상한·body timeout/재전송/ACKclaim중단을 검토. 기존원점Access가trial앱JWT를받는후보도검토. trial signup/owner/admin/pairing 우회가 없어야함.
- [x] 실제target lint --from f284948을 별도cleanclone에서수행해 리뷰lint.json저장. 정확한SHA check --from f284948 --to cd60e7f --task-key SAR-MVP-003-BIDIRECTIONAL exit0. ERROR/critical/high는 수락차단. medium/low 영향과 게시/배포가능범위를명시. 실제Grok미검증을확장금지.
- [x] report.md/result.json/근거를 자기review폴더에작성. 코드결함은 재현·등급·차단으로보고. 전체코드를재작성하지않음.

## Jev 목록·충돌
자기 key의 find/documents-find/context JSON keep 모두 읽는다. 지시 전제와 충돌 — 먼저 확인: DEV실행기록은 준비성공과 actual계정차단을 함께 적는다. 실제Grok왕복성공으로 확장하지 않는다. 신규QA/리뷰만 자기검증범위의근거다.

## 완료 보고

- 결과: 고정 head `cd60e7f87eb5ce137eca887980f232b3f67a18d0`(base `f2849486ed48295e239714700e651d30c32f1c2c`)을 수락 가능으로 판정했다. critical/high/medium 0, low 5(비차단·미해결)다.
- 독립성: 검토자 Claude Code `claude-opus-5-5` high 세션 `c4c411b8-d66d-4356-981b-9e2279440b02`, 구현자 Codex 세션 `01a104d1-8cc5-7430-bcf1-1a2732c5183f`. snapshot `/tmp/knowslink-bidirectional-review-cd60e7f`는 detached·clean·읽기 전용으로 보존했다. 실행은 별도 scratch clone에서 했다.
- 검증: target lint `--from f284948` exit 0(ERROR 0/WARNING 2 SIZE-001/실행 불가 0), `make test` exit 0, `make verify-mvp` exit 0, access plan selftest exit 0, 리뷰어 probe 13개 우회 경로 403·allowlist 제거 fail-closed. `review.py check --task-key SAR-MVP-003-BIDIRECTIONAL` exit 0(reviewed 47/47).
- low: 업무 `/v1/pull`의 trial lease 교차, trial claim policy 라벨, cloudflared path regex 비고정, run_trial 키 파일 권한·TOCTOU, 비정상 응답 본문 미해제. 상세는 [리뷰 보고서](../docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md).
- 남은 일: 실제 Grok 왕복·실제 Access token/path 앱/trial AUD 인증 호출·운영 후보 배포는 미검증이다. tester QA와 coor 병합·D12 13장 적용이 후속이다. 제품 코드는 수정하지 않았다.

