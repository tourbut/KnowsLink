---
title: SAR-MVP-003-BIDIRECTIONAL-OPS — 기존 link 시험 운영적용
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-OPS]
summary: 검수된 시험 transport를 기존 배포와24시간 machine인증에 적용하고 실제 접근을 확인한다
---

# SAR-MVP-003-BIDIRECTIONAL-OPS — 기존 link 시험 운영적용
- queued, coor가 최신711f253 delta리뷰/좁은QA 수락·main/origin확인 후 dispatch하면ready. /home/shin/orca/workspaces/KnowsLink/fullops-ops fullops/ops. 병합coor/main/origin.
- 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9 run_8ca8bc058ab7 task/dispatch preamble.
- 사용자승인: Codex↔Grok 실제시험text, 주소 link.knowslog.com. 필요한기존배포갱신·시험인증·privateconfig준비·GitHubpush/issue댓글 승인. D12의과거별도servicetoken승인요청은 이번실제시험승인이우선한다. 신규비용/다른서버/owner보호해제/업무effect/dots/FullOps업데이트 제외.
- 수락제품711f2532be423d1ca7707463a20fdc168f50bece 및 전체cd60review0c367301/timeoutreviewb15740de, QAfa168938 불변왕복과 최신좁은QA(coor dispatch에SHA전달).
- FULLOPS/common0.3.2 README·연결3규칙/project/documentwriting/fullops-work/ponytail, D12 11~13장·DEV실행기록2개·access_trial_plan.py 먼저읽음. Cloudflare One 관련스킬/공식자료가필요하면적용. Jevkeep/충돌확인.

## 승인 작업·완료
- [ ] 실제state와현재배포28bd1bb/기존rootAccess/정책/Tunnel/DNS·공유컨테이너를readonly확인·비밀없는baseline보존. /home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL의DEV생성trial_codex/trial_grok key/config/pair 재사용. 폴더0700/key·config0600 소유regular심링크없음 직접확인(review low4 완화).
- [ ] 기존beta.sh deploy로coor확인된origin/main reviewedSHA를반영하고 시험allowlist 두agent만명시. 기존DB백업/rollback/restart/health/logsecret없는근거. 새서버/Tunnel/DNS생성·초기expose재실행금지.
- [ ] D12/access_trial_plan의rendered검수본문으로 두distinct 24h CFservice tokens·trial전용non_identity ServiceAuth policy·link.knowslog.com/v1/test/* 앱과trialAUD 생성/GET검증. rootowner앱·이메일policy·IdP/team조직 기존자원변경금지. 이름충돌이면기존우리trial자원인지직접확인하며 임의사용/덮어쓰기금지. secret응답은로컬0600에즉시저장, toolstdout/chat/issue/trace/Git에출력금지. 로컬관리OAuth/API읽기파일을읽어야하면허가된file-storebridge처럼메모리사용·필요권한확인·출력금지. 무권한이면preparedbody와권한명만보고.
- [ ] 새원점tunnel후보(검수된prefixregex·trial/owner각AUD·required true)를pinnedimagevalidate/rule/negative로확인. 원점config백업후knowslink connector만재기동. beta render-config 단일AUD로trialrule덮어쓰기금지.
- [ ] 실제공개negative: no/invalidCFauth거부·validCF/wrongrelayauth거부·token만으로owner/signup/pair/business불가. 두유효CF/agent 자격의registry/keys/전용send/pull권한검증. 원문은합성testtext만. 필요하면두localclient→실제HTTPS 왕복nonce/ID로transport검증하되 actualGrok아님명시. originalbusinesspull 금지(reviewlow1완화). actualGrok용시험송신은private전달/parent준비전TTL유실을피해대기.
- [ ] agent별privateconfig CFheader를correctagent에삽입하고 Grok측path를/workspace/.knowslink-trial/trial_grok/key.pem에맞춘환경파일·key 최소export를state0700내준비. owner/관리token/Codexkey 전달금지. secret전달수단이issue5976122078에아직미확인이라자동외부전달금지. 사용자에게준비파일경로와필요private수단만보고. credential값없이실제CF만료시각/TTL/복구·종료절차를정리.
- [ ] D12/D13/D10 실제운영결과·실행기록 SAR-MVP-003-BIDIRECTIONAL-OPS.md. Grok최종issue댓글초안은고정reviewedSHA·installhash·Pythonlauncher정확값·private파일수신후chmod·수동receive/send회신ID·실패로그·준비완료reply조건·TTL180초에맞춘coor송신협업을포함. 아직actualGrok미검증/자동wake없음/업무effect없음.
- [ ] selflint --from착수HEAD0, 변경제품소스없음확인, work.pyfinish/archive/commit·worker_done 과제key final SHA40자리/실제접근근거/남은private전달차단. 외부issue게시와main통합은coor.
- 실패면원인별최소rollback/시험자원정리와기존owner서비스정상확인, user자료삭제금지. 단기token만료에맡기지않고정상시험종료revoke/시험앱·policy제거/ownerconfig복원/allowlist비움절차을남긴다. 실제Grok준비미확인이면transportready까지만완료하며왕복완료아님.

## Jev 목록·충돌
자기key의find/documents-find/context keep모두확인. 지시 전제와 충돌 — 먼저 확인: D12/DEV원본기록은서비스인증·후보배포미완료를적는다. 이번사용자승인으로 검수후그구체적준비작업을재개하나 실제계정왕복성공으로확장하지않는다. 원본베타owner-only보호와공유서비스는보존한다.
