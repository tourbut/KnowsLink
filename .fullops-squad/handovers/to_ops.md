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

- [x] 별도세션 독립 delta 리뷰를 끝낸다. 새 구현자 actualClaude9ebf9a10-609b-475a-b1cd-f9edad84c9a1과 원DEV01a11106-b736-7db1-b518-a65b64dbc5fb를 자기리뷰로 대신하지 않는다. independence에는새구현자/actual자기session을 기록한다. read-only detached snapshot `/tmp/knowslink-messages-fix-review-dfc70ca`의clean/fullHEAD를확인하고 읽기만한다. 설치/실행은 별도scratch다.
- [x] H1 본문이슬롯전수신으로옮겨졌을때 자원/연결/대기 상한·10s/8/32KiB와익명slowbody뒤유효cleanup 보호를재현한다. M1stateless자격/CSRF/closedbody+transaction ownership+로컬 live색인의stale/multi-instance/revocation/swap경합이유효정리를막거나익명/타owner/잘못된본문이cleanup budget을점유하지않는지검토한다. 모든정리caller ACK/key/agent/unpair/deny/cancel/gate/logout·owner/reply/text와공유16+4/rate/crash30s/restart를필요범위로확인한다. 새보호가다른층의무제한점유로옮긴것인지판정한다. edge미검증을통과근거로삼지않는다.
- [x] F-UI-MSG-01 canonical redirect·현재권한/CSRF/정책회귀와DEV대조군/변형RED·최종fixed통과를확인한다. 기존review checkexit1/H1/M1을고치지않고새fixed해소근거를쓴다. L1전역lock/부하·L-A/L-B 공개전조건과실운영미검증은보존한다.
- [x] 모든(path,status)의reviewed/skipped구체근거·원천/공통기준/규모·DEP/DESIGN/SLOP/SIZE판단·미검증·원실패를 report/result에기록한다. 새fixedlint --from09c의head/등록kind명령exit를확인하고재현필요영향검증만실행한다. COPY로highresolved처리하는probe는수락이아니며사용하지않는다. check실패면그대로실패/남은수정담당을보고한다.
- [x] work.py finish전문/빈인박스·최종recordsHEADlint/test/strict/diff·일반push뒤 fullSHA/리뷰결론/근거를worker_done으로보고한다. 최종JSON레포밖이면coor가영속화할정확한경로를보고한다.

## 소유권·제약

OPS는새FIX-REVIEW-review/자기phase/context/PLANS·정규인박스/log만작성한다. 제품코드·기술/제품정본·기존결과·타인인박스는변경하지않는다. 시각판정은designer·동작QA는tester별도다. 새세션Opus5.5high는fullops-review의인증/데이터/동시성고성능규정이며routeSonnet추천override근거를PLANS에보존했다. 기록완료와main제품수락을구분한다. Workers Free만허용·과금/유료구독/운영배포/실메일/운영자료/서버Tunnel변경금지. 승인된읽기전용검토/자기격리검사/자기자원회수/문서커밋·일반push는재승인없이완료한다. 제품규칙질문만coor경유designer, 기술소견/수정요청은coor경유DEV다.

## 완료 보고

- 결과: 고정 `09c523da8a3407288d9f5d711e1834af12bc7808..dfc70caa748a90614b02d48c78b4651345938339` 독립 delta 보안 리뷰를 완료했다. 결론은 수락 불가다. 원 H-1·M-1·F-UI-MSG-01은 해소, 새 미해결 high H-2, low L-1 확장·L-2, 기존 L-A/L-B 보존이다. 제품 코드는 수정하지 않았다.
- 독립성: 리뷰 세션 `3cee40a4-bf54-4f2c-bad3-1559b8f6b006`(Claude Code, claude-opus-5-5). 수정 DEV `9ebf9a10-609b-475a-b1cd-f9edad84c9a1`·원 DEV `01a11106-b736-7db1-b518-a65b64dbc5fb`과 다르다. snapshot은 리뷰 전후 detached dfc70ca·porcelain 0줄이며 읽기만 했다. 설치·lint·재현은 scratch clone에서 했다.
- H-1 해소: 상한 본문을 슬롯 전 10s 안에서 받는다. 연결당 32KiB·10s로 header 수신과 같은 계층이며 새 무제한 층이 아니다. TestHTTPSlowBodyAndCleanupAdmission PASS, DEV 대조군 09c exit1.
- M-1 해소(익명·위조 범위): 익명·위조·잘못된 본문·cross-site·잘못된 gate CSRF는 정리 채널·Clean·cleanup rate를 쓰지 않는다.
- F-UI-MSG-01 해소: 결정 성공 뒤에만 고정 back의 정식 gate 화면으로 303한다. CSRF·현재 권한·정책 무변경. 시각 판정은 designer UX07 몫이다.
- H-2 high: `internal/relay/capacity.go:159-215`. 로컬 정리 채널을 유효 자격만으로 고르고 소유·rate는 DB 입장 transaction 안에서 판정한다. 그동안 슬롯을 잡는다. 실제 Postgres 재현: row lock 중 타 owner 철회 4건이 4/4 점유, 유효 자기 철회 429 capacity. lock 없는 flood에서 신규 flood 대조 0/18, 타 owner 정리 flood 18/18, agent 자격의 lease 없는 ACK flood 18/18. 인증 회원 하나가 프로세스의 모든 ACK·철회·deny·unpair·로그아웃을 막는다. PS-11 위반이다.
- L-1 확장: 커밋마다 전체 자격 색인 재생성. L-2: 색인 stale·restart nil·Store 순서 역전. 관찰: `/v1/connect/*` 정리 경로는 HTTP 정리 분류 밖(09c와 같음).
- 검사(dfc70ca scratch): make install exit0, lint.py --from 09c523d exit0(ERROR0/WARNING5/실행불가0, product-lint/test exit0), make verify-mvp 2회 exit0(PASS 42·FAIL/SKIP 0, 자기 Compose 자원 회수). review.py check --task-key SAR-PUBLIC-MESSAGES-001-DEV-FIX exit1(H-2 정상 차단). 형식 probe 미사용.
- 경고 판단: SIZE-001 4·SIZE-002(568줄) 분리 사유 수락. DEP0·DESIGN0·SEC/SLOP 0.
- 산출: `docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/`의 report.md·result.json·lint.json/log/exit/head·install.exit·verify-mvp(.log/.exit/.head, run1-*)·review-foreign-cleanup_test.go.txt·review-check.log/.exit, PLANS·contexts/ops append. 원 REVIEW-2 기록·DEV 실패 증거는 변경하지 않았다.
- 미검증: 운영 Tunnel/edge 본문 완충·연결 수 상한, 실부하·상태 크기·CPU, 실제 브라우저 시각, 실메일·실24h·외부 계정·운영 배포. verify-grok-plugin은 adapter 변경이 없어 실행하지 않았다.
- 후속: DEV가 H-2(권장 L-2·/v1/connect 확인)를 고치고 자기 대상 반복·타 owner flood 회귀 검사를 넣은 새 fixed SHA를 낸다. OPS가 새 리뷰 키·별도 세션으로 재검토한다. 이 완료는 main 제품 수락이 아니다.
