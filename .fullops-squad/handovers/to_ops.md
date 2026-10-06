---
title: SAR-PUBLIC-MESSAGES-001-FIX-REVIEW — 수정 HTTP 정리 입장·gate 후보의 독립 delta 보안 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-REVIEW]
summary: 수정 HTTP 정리 입장·gate 후보의 독립 delta 보안 리뷰
---

# SAR-PUBLIC-MESSAGES-001-FIX-REVIEW — 수정 HTTP 정리 입장·gate 독립 delta 보안 리뷰

- 상태: ready. base09c523da8a3407288d9f5d711e1834af12bc7808 → fixed dfc70caa748a90614b02d48c78b4651345938339. 원본 DEV37/OPS6e/UIf154의 실패/조건은 불변이다.
- 담당 ops / /home/shin/orca/workspaces/KnowsLink/fullops-ops / fullops/ops. 결과는 기록checkout에만 작성한다.
- 복귀 coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7 / repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. Task/Dispatch는 preamble 기준.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3 README/coding-style/testing/security·FULLOPS·project·document-writing·contexts/ops를 전달SHA에서 읽는다. 제품 PS11/04/06/07/08–10 SAR-PUBLIC-SERVICE·MVP C1–C5·UX06/07 정본, DEV-FIX phase/완료로그/QA증거와 원REVIEW-2 report/result/repro를 읽는다. review/rule.json 및 FIX-REVIEW의preview/rules/result를 적용한다. external open-code-review-delegate 스킬을 읽는다. 코드기준은 capacity/cleanup_admission/store/http와그모든호출자·새 unit/integration시험·현재인증/CSRF/owner판정이다.

## 할 일과 완료 조건

- [ ] 별도세션 독립 delta 리뷰를 끝낸다. 새 구현자 actualClaude9ebf9a10-609b-475a-b1cd-f9edad84c9a1과 원DEV01a11106-b736-7db1-b518-a65b64dbc5fb를 자기리뷰로 대신하지 않는다. independence에는새구현자/actual자기session을 기록한다. read-only detached snapshot `/tmp/knowslink-messages-fix-review-dfc70ca`의clean/fullHEAD를확인하고 읽기만한다. 설치/실행은 별도scratch다.
- [ ] H1 본문이슬롯전수신으로옮겨졌을때 자원/연결/대기 상한·10s/8/32KiB와익명slowbody뒤유효cleanup 보호를재현한다. M1stateless자격/CSRF/closedbody+transaction ownership+로컬 live색인의stale/multi-instance/revocation/swap경합이유효정리를막거나익명/타owner/잘못된본문이cleanup budget을점유하지않는지검토한다. 모든정리caller ACK/key/agent/unpair/deny/cancel/gate/logout·owner/reply/text와공유16+4/rate/crash30s/restart를필요범위로확인한다. 새보호가다른층의무제한점유로옮긴것인지판정한다. edge미검증을통과근거로삼지않는다.
- [ ] F-UI-MSG-01 canonical redirect·현재권한/CSRF/정책회귀와DEV대조군/변형RED·최종fixed통과를확인한다. 기존review checkexit1/H1/M1을고치지않고새fixed해소근거를쓴다. L1전역lock/부하·L-A/L-B 공개전조건과실운영미검증은보존한다.
- [ ] 모든(path,status)의reviewed/skipped구체근거·원천/공통기준/규모·DEP/DESIGN/SLOP/SIZE판단·미검증·원실패를 report/result에기록한다. 새fixedlint --from09c의head/등록kind명령exit를확인하고재현필요영향검증만실행한다. COPY로highresolved처리하는probe는수락이아니며사용하지않는다. check실패면그대로실패/남은수정담당을보고한다.
- [ ] work.py finish전문/빈인박스·최종recordsHEADlint/test/strict/diff·일반push뒤 fullSHA/리뷰결론/근거를worker_done으로보고한다. 최종JSON레포밖이면coor가영속화할정확한경로를보고한다.

## 소유권·제약

OPS는새FIX-REVIEW-review/자기phase/context/PLANS·정규인박스/log만작성한다. 제품코드·기술/제품정본·기존결과·타인인박스는변경하지않는다. 시각판정은designer·동작QA는tester별도다. 새세션Opus5.5high는fullops-review의인증/데이터/동시성고성능규정이며routeSonnet추천override근거를PLANS에보존했다. 기록완료와main제품수락을구분한다. Workers Free만허용·과금/유료구독/운영배포/실메일/운영자료/서버Tunnel변경금지. 승인된읽기전용검토/자기격리검사/자기자원회수/문서커밋·일반push는재승인없이완료한다. 제품규칙질문만coor경유designer, 기술소견/수정요청은coor경유DEV다.
