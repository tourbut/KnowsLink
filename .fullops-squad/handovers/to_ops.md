---
title: SAR-MVP-003-BIDIRECTIONAL-OPS — 권한 확보 후 시험 Access 인증과 공개 HTTPS 검증 재개
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-OPS]
summary: 권한 확보 후 시험 Access 인증과 공개 HTTPS 검증 재개
---

# SAR-MVP-003-BIDIRECTIONAL-OPS — 권한 확보 후 시험 인증 적용

- 상태 ready. 기존 과제 후속이다. 이전 결과 fa221886351b12404702f34348c5be45a520541f, main/origin39993d0631f8ae1569d2fdea0c1246510ca38223. 이전 세션은 released라 fresh Sonnet5.5 high(기존 route)로 재개한다.
- 담당 /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops. repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. 병합coor/main/origin.
- 복귀 run_8ca8bc058ab7, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, task/dispatch/handle은실제preamble.
- 승인: 기존link 실제text왕복시험·최소24h machineauth·배포·GitHub댓글. 추가 사용자 지시로 로그인된 Orca browser에서 관리API token발급완료. 신규비용·owner보호해제·shared변경·업무효과·dots·FullOps업데이트 제외.

## 먼저 읽을 문서와 기준
FULLOPS.md, rules/common/README.md 및 연결3규칙(0.3.2), project.md, docs/agents/document-writing.md, fullops-work/ponytail. docs/operations/ops-guide.md 11~13 및 docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md 재개절차. deploy/knowslink/access_trial_plan.py. 같은key의 docs/evaluations/jev/*find/context keep와 충돌확인. 기존route/model을재사용(제품판단/제품코드변경없음).

## 새 관측과 비밀 취급
coor가 사용자 승인 Orca browser에서 knowslink-trial-ops-24h 사용자API token발급. 계정4d545b8037b67ac95d83fe4276ce4aca 하나, Access ServiceTokens Edit와AppsPolicies Edit 두권한만. verify active·expires_on2026-10-04T23:59:59Z, servicetokens/apps/policies GET모두성공. 브라우저secret화면을목록으로이동했고임시secret snapshot삭제.
관리 token은 /home/shin/deploy/knowslink-state/cf-service-token-api.env (0600, CF_API_TOKEN 단일line)에있다. Python으로메모리에읽고requests/urllibheader로사용하라. cat/source/print/argv/envexport/오류본문출력금지. 토큰값은로그·Git·issue에남기지마라. Grok에관리token전달금지. API응답clientsecret은직접0600파일로저장하고출력금지.
Grok사전회신5976609951: 오너private파일배치수단가능. CLI시험우선, 기존MCP삭제/재등록불필요. 실제파일전달아직안했고CodexactualGrok용송신대기.

## 해야 할 일과 소유권
- [ ] 이전OPS baseline/root app/policy/공유서비스·state권한재확인. 이미deploy0911c2c+allowlist이며변경없는배포·loopback·제품전체QA를반복하지않는다.
- [ ] 이전OPS재개절차대로 distinct24h service tokens2개→non_identity reusablepolicy→link.knowslog.com/v1/test/* selfhosted app/trialAUD. 검수된access_trial_plan.py본문사용. 이름충돌확인. secret즉시 agent별privateconfig정확삽입. root app/policy/IdP/team 변경금지.
- [ ] 기존원점config백업유지. trialrule/trialAUD·ownerrule/ownerAUD requiredtrue 후보 pinnedimage validate/rule후knowslink connector만재기동. beta render-config/expose 금지. 실패 rollback.
- [ ] 공개negative(no/invalidCF, validCF+wrongrelay, token만owner/signup/pair/business불가)와positiveregistry/keys/send/pull. 합성text두localclients HTTPS왕복ID/nonce확인(actualGrok아님). businesspull금지.
- [ ] grok-export의key/environment0600·dir0700·ownedregular확인. CFheader삽입완료, /workspace/.knowslink-trial/trial_grok/key.pem절대경로. 준비파일경로와만료만보고. 외부전달은coor/user안전채널담당. 시험종료전바로tokenrevoke하지말고종료revoke/원점복구/allowlist비움절차명시.
- [ ] D12/D13·OPS실행기록에현재차단해소/공개실제검증/만료와복귀근거추가. 원본당시차단기록보존. Grok최종댓글초안정확명령포함. 제품소스변경없음.
- [ ] selflint --from착수SHA0, work.pyfinish로완료보고전문/archive/빈inbox, commit·worker_done key와SHA40자리. issue게시/main병합은coor.

## 완료 기준과 후속
운영적용+실제공개positive/negative 모두통과하면transportready. actualGrok미검증·수동pull·TTL180초유지. 새로운설명되지않는실패/권한부족/경계변경이면ask. 사용자자료삭제금지. 고정검수된제품0911c2c불변(최신문서HEAD와동일제품), 기존리뷰/QA재사용. 산출물D12/D13 및조건부D10. 공개실패를성공으로표시하지않는다.

## 완료 보고
worker가결과/검증/못한것/만료/후속을채운다.
