---
title: SAR-MVP-003-TRIAL-CLEANUP — 실제 왕복 시험 종료와 임시 인증 정리
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-MVP-003-TRIAL-CLEANUP]
summary: 실제 왕복 시험 종료와 임시 인증 정리
---

# SAR-MVP-003-TRIAL-CLEANUP — 실제 시험 종료
- ready. 실제 왕복이 Grok결과5993585072와 여기직접송수신ID로확정됐다. 성공댓글5993750758. 사용자가원래승인한시험종료절차(D12/이전OPS지시서)를실행한다. 새로운제품/유료자원/업무effect/dots/FullOps업데이트제외.
- ops /home/shin/orca/workspaces/KnowsLink/fullops-ops fullops/ops, repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. base main/origin439d74a19e1628db268aa7a2788ff6b29c9b2d04. 복귀run_8ca8bc058ab7/preambleworker_done, 병합coor/main/origin.
- freshSonnet5.5high route. 기존RENEW세션user_takeover retained이므로사용자터미널재사용/닫기금지.

## 먼저 읽을 문서·기준
FULLOPS.md, commonREADME연결3규칙0.3.2/project/document-writing/fullops-work/ponytail. docs/operations/ops-guide.md13절, exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md와-OPS-RENEW.md 종료절차, SAR-MVP-003-ACTUAL-TRIAL.md최종대조. 정본직접링크와자기route, 추가탐색은jev_find/context규약.

## 승인된 범위와 실제 private 경로
관리token /home/shin/deploy/knowslink-state/cf-service-token-api.env0600 CF_API_TOKEN은2026-10-05T23:59:59Z만료예정. Python메모리에서만읽고값출력/argv/export/Git/issue금지. trial state=/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL. 이자원파일trial-access.json과liveGET을대조한다.
우리trialtoken UUID codex283b0bdf-3c2b-4c99-b9ce-96516d226cd9/grok647039da-10bb-4459-a595-6e0caa545adb, trial앱84b33961-6f78-42d3-b27b-06e064a05c2b, reusablepolicy884305ff-26ef-4fc5-828a-8145e7a5706c. 이름/type/path/UUID검증전삭제금지. rootowner앱·정책·IdP/team/Tunnel/DNS·공유서비스·DB사용자자료보존.

## 해야 할 일과 파일 소유권
- [ ] 실제trial자원/rootbaseline/공유서비스·현재ownerbackup/liveconfig/allowlist확인.
- [ ] 위2token만revoke/delete, trial앱·trialpolicy만삭제. GET으로우리trial자원제거와rootowner앱/정책불변확인. 관리사용자API token은무효화endpoint권한별도라값삭제/발급재시도하지않고만료기록만유지.
- [ ] /home/shin/deploy/knowslink-state/tunnel-bak-pre-trial/config.yml owner-only백업을live에복원0600, pinnedcloudflared validate/rule확인후knowslinkconnector만재기동. relaystate .env KNOWSLINK_TEST_AGENTS줄만비움/제거후기존composerelay갱신. 제품source/DB/newTunnel/DNS/초기expose변경금지.
- [ ] localtrialrouteclosed·공개owner302·이전시험serviceCFheaders로trial/owner거부확인. 공유서비스orca200/s8200/mcp401·relaypostgreshealthy·owner동작regression확인. 성공실제왕복source검증은반복하지않는다.
- [ ] 로컬시험privatekey/archive/ownerJSON은사용자파일이므로이번작업에서삭제하지말고0600보존·credential현재더는연결불가명시. Grok측폴더/받은tar제거는Bot에게댓글로요청됐고직접원격삭제하지않는다.
- [ ] D12/D13/새종료실행기록 docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md: 삭제자원증거/원점복원·allowlist/public/local/공유정상/남은MCP자동wake dots/파일삭제요청만기록. 실제왕복성공결과원본보존.
- [ ] selflint --from착수HEAD0, 제품변경없음·work.pyfinish전문archive/빈inbox/commit·worker_done key/fullSHA40자리. issue최종운영댓글/main병합은coor.

## 실패·완료
범위이탈/설명불가실패/삭제ID불일치/backup불일치면ask와현재상태보존. trial인증정리+기존owner보호복귀+공유회귀통과가완료다. 사용자자료/기존owner자원삭제금지. 실패를성공으로보고하지않는다.

## 완료 보고
worker가검증/정리후상태/한계를채운다.
