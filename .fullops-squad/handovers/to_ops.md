---
title: SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW — 만료된 시험 인증 갱신과 Grok private 전달 묶음 준비
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW]
summary: 만료된 시험 인증 갱신과 Grok private 전달 묶음 준비
---

# SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW — 시험 인증 갱신
- ready. 승인된실제Codex↔Grok text시험을2026-10-05재개한다. 이전진행못함 사용자지시. FullOps업데이트/업무effect/dots/새비용제외. 기존Orca브라우저발급승인유효하며coor가같은최소관리token을갱신했다.
- ops /home/shin/orca/workspaces/KnowsLink/fullops-ops fullops/ops, repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. 병합coor/main/origin. run_8ca8bc058ab7 preamble worker_done, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor.
- 기준main/origin31cf1a5b6ab0e9ac607f86ad06d7ad797ceeab11. 제품0911c2c불변, 이전OPS2b70909a9a82b5f03daf100b289f1c15fb3f16f1의공개HTTPS검증통과·실제Grok未실행. agent pair/key재사용.

## 먼저 읽을 문서·기준
FULLOPS.md, commonREADME연결3규칙0.3.2/project/document-writing/fullops-work/ponytail. D12 docs/operations/ops-guide.md 11~13, docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md 재개결과/만료/종료, deploy/knowslink/access_trial_plan.py. 원천범위탐색이추가필요하면자기key jev_find/jev_context와keep/충돌규약수행. 기존문서주소확정으로직접참조한다. 모델route의opsSonnet5.5high, fresh(이전released·하루지남).

## 현재 관측
2026-10-05T08:04Z 현재 기존두service tokens는06:42Z에만료, 관리token도만료. coor가Orca로그인browser에서새 knowslink-trial-ops-20261005 사용자API token생성: 계정4d545b8037b67ac95d83fe4276ce4aca하나+ServiceTokensEdit/AppsPoliciesEdit두권한만, active확인,만료2026-10-05T23:59:59Z. 파일 /home/shin/deploy/knowslink-state/cf-service-token-api.env 0600 CF_API_TOKEN단일line으로교체완료. secret화면종료,값출력없음.
Grok installer0911c2c/hash통과, parentMCPtools2불변이지만CLI시험가능. /workspace/.knowslink-trial/trial_grok는0700 box:box 준비. 사용자파일전달未완료. 댓글5977433403의지난만료절차는이번새결과로갱신해야한다.

## 소유권·할 일
- [ ] 기존배포/공유서비스/rootowner앱·정책/원점설정·private권한baseline확인. 제품소스/배포재빌드/DB/newTunnel/DNS/owner정책변경금지.
- [ ] 관리API token을Python메모리만읽어동일scope로두시험service token을24h갱신. 공식API가existingrefresh를지원하면그검증된방법을우선하고 UUID유지/secret변경/expiry응답을확인. 불가하면우리만료token2개만삭제재생성·기존우리trialpolicy의includeUUID두개만갱신. 기존앱/path/AUD/Tunnel변경불필요하면보존. CFAPI primary문서근거확인, 값은stdout/argv/envexport/Git/issue/log에출력금지.
- [ ] agent별CFheader를trial_codex/trial_grok/grok-export의정확한privateconfig에갱신·ownedregular0600/key0600/dir0700확인. AGENT_KEY_FILE Grokexport=/workspace/.knowslink-trial/trial_grok/key.pem 유지. 메타trial-access.json만료/ID갱신.
- [ ] 실제공개negative(no/invalidCF/wrongrelayauth, validCF로owner/business거부) 및positive registry, 두localclients합성textHTTPS왕복. actualGrok아님. 이전제품전체QA/배포/loopback재실행금지. TTL180/업무pull금지. 예상403/401/owner302판정,실패최소복구/ask.
- [ ] 사용자전달용private archive /home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/knowslink-grok-trial-20261005.tar.gz 0600생성. 내용은grok-export의key.pem/environment.json두개만, archive상대basename만/소유권명령후처리. Codexkey/owners/adminAPI token포함금지. archive목록만검증, 외부upload하지말고로컬경로/만료만보고.
- [ ] D12/D13와새실행기록 docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md: 실제갱신/expiry/negativepositive/파일권한·archive안전내용·원본증거재사용·종료revoke정리. 원본지난기록보존. exactCLI tar추출명령/Grok준비회신·새roundkey案을댓글초안으로기록. 실제Grok아직미검증.
- [ ] selflint --from착수HEAD0·제품변경없음·work.pyfinish 전문archive/빈inbox·commit·worker_done fullSHA40자리. main/issue/외부private전달은coor/user.

## 완료 기준/제약
갱신된private자격으로publictransportready까지. 실제Grok용송신은파일전달/parent준비회신전대기. 관리token과serviceexpiry차이를명시하고시험종료시우리trialtokens/policy/app삭제·원점owner복원·allowlist비움절차유지. user자료/기존owner/공유자원보존. 권한부족/검증실패/범위이탈이면ask하고근거보존.

## 완료 보고
검증/새만료/파일경로/남은실제Grok전달을채운다.
