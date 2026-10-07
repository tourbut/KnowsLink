---
title: SAR-PUBLIC-MESSAGES-001-DEV-FIX-2 — 유효자격 정리 flood H2 원인 수정과 모든 정리 경로 보호 회귀
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX-2]
summary: 유효자격 정리 flood H2 원인 수정과 모든 정리 경로 보호 회귀
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX-2 — 유효자격 정리 flood H2 원인 수정·회귀

- 상태: ready. 같은MESSAGES의후속기술과제이며DEV가분석/계획/구현/테스트를함께완료한다.
- 담당 dev / /home/shin/orca/workspaces/KnowsLink/fullops-dev / fullops/dev.
- 복귀 coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7 / repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. Task/Dispatch는preamble기준.

## 현재 상황·적용 기준과 먼저 읽을 문서

기준 fixeddfc70caa748a90614b02d48c78b4651345938339다. 원09c H1/M1 및UX07 FAIL과dfc 새H2high를보존한다. OPS c219의FIX-REVIEW는실제Postgres 타owner/no-lease flood에서유효정리18/18이429, 신규flood대조0/18로수락불가/checkexit1이다. 원H1/M1익명범위/redirect는dfc에서해소판정을받았지만새H2가남았다. 현재UI는dfc UX07를마무리하고QA는09c 원본을마무리한다. 다음최종fixed의독립delta/좁은QA·필요UI수락전main보류다.

fullops-common-0.3.3 네규칙·FULLOPS/project/document-writing/contexts/dev를전달SHA에서읽는다. 제품SAR-PUBLIC-SERVICE PS11및정리budget조건·PS04/06/07/08–10·SAR-MVP C1–C5·UX06/07를읽는다. FIX-REVIEW-review report/result/review-foreign-cleanup_test.go.txt/verify-mvp.log/exit·원REVIEW-2 및DEV-FIX phase/완료로그/QA·Jev code/doc/context를읽는다. 공통경계 capacity/cleanup_admission/store/http/member/connect/currentprincipal/CSRF/lease/rate와 모든정리caller가범위다.

## 할 일과 완료 기준

- [ ] H2를실제자기격리DB에서재현하고짧은기술계획을자기phase에쓴다. 로컬정리선택은현재live자격만보고 자기대상/유효lease/rate는DB입장후에본다. DBpool/rowlock대기도4정리슬롯을점유하여한회원의타대상/no-lease flood가모든정리를차단한다. rate제한만으로슬롯보호되지않음을확인하고공통경계의모든caller를추적한다.
- [ ] PS11의현재인가/유효lease/자기기록과신규16+정리4분리를지키며타owner/no-lease/잘못된자격/본문/CSRF/flood가안전정리budget을독점하지않게근본수정한다. 유효자기대상반복/멱등철회flood도다른principal의정리를굶기지않게검사한다. 구조/알고리즘/보호값은DEV판단이며상한·공정성·비용·UX영향과근거를기록한다. 제품quota/제품규칙변경이필요하면coor경유designer판단이다. 단순증상guard나미확인edge기대·다른층의새무제한대기/상태로옮기지않는다.
- [ ] 같은수정에서 L2 프로세스색인stale/새자격crossinstance/restartnil/커밋후Store순서역전을분석하고보호회귀를포함한다. /v1/connect cancel/key-revoke/agent-revoke/unpair도같은정리경로로분류누락/현재권한을확인하고필요수정한다. 정리handler의정책/CSRF/closedschema/현재키·관계세대·lease/클라이언트credentials의권한을낮추지않는다.
- [ ] 결정적rowlock재현과실제flood(타owner·agentno-lease·자기반복·신규대조)의유효정리성공/신규16정리4/rate/cleanup/fairness을자동회귀로남긴다. H1익명/인증slowbody·본문8/32KiB/deadline10s·공유다중instance/재시작/crash30s/철회·lease·UX07 canonical redirect/CSRF·기존text/reply/receipt/claim 영향회귀를필요범위로검사한다. 원fixed반례/수정후통과의직접exit/HEAD를남긴다. 원래실패를고치거나high를resolved로뒤집지않는다.
- [ ] 관련기술정본 D10 및실제영향 D03/05/06/07/09·자기phase docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md·QA·context/PLANS를현재코드와일치시킨다. 변경없는산출물은이유를쓴다. 제품/UX정본·기존리뷰/QA/UI·타인인박스는변경하지않는다. L1전역lock/색인전체스캔비용·기존L-A/L-B·공개전합성가입unset/합성owner0·실운영후속은보존한다.
- [ ] 최종코드와archive/metadata마지막HEAD까지make lint/test·영향verify-mvp·FullOps --fromdfc lint/test ERROR/WARNING/실행불가·strict/diff·일반push를완료한다. 자신의명령exit를보존하고후행명령/파이프로실패를숨기지않는다. work.py finish전문/빈인박스뒤worker_done fullSHA/최종JSON경로/원인/회귀/실패·미검증/후속을한번보고한다.

## 제약과 완료 보고

Workers Free만허용·유료플랜/구독/초과과금금지·기존서버/Tunnel·wire·기존제품한도를유지한다. 운영배포/실메일/외부계정/실24h/운영자료삭제/노우↔다닷은이번과제가아니다. 승인된로컬비파괴개발/자기격리검사·자기자원회수·커밋/일반push는재승인없이완료한다. 새독립OPS 세션/read-only fixed리뷰·TESTER좁은QA·필요직접UI와최종main판정은coor후속이다. 제품/범위질문만coor경유designer로보내며기술분석/방법을상시제품승인으로돌리지않는다. 예상변경규모·DEP/SLOP/DESIGN/SIZE경고처리·원RED/환경실패를전문에기록한다.

## 탐색 보완·지시 전제와 충돌 — 먼저 확인

coor가Jev에쓴 internal/relay/connect.go는없는경로다. 원거부JSON을보존하며실제 internal/relay/connections.go·member_agents.go·connections_integration_test.go가필수keep이다. 기존DEV-FIX 자동PASS는H2발견전조건이며현재독립FIX-REVIEW high/checkexit1을우선한다. sensitive/oversized로Jev미전송된cleanup_admission.go는직접읽는다. 기술전제는이번과제에서분석하고제품규칙변경만coor경유질문한다.

## 원본 TESTER 마감 후 추가 근거

사용자가 Grok 토큰 만료로 원본 TESTER 추가 실행을 중단했다. coordinator 마감 기록과 msg_78c8fcd13251을 참고한다. 기록60b9892/제품09c의 독립 probe17차 exit1은 보존한다. receipt20000 거절409 뒤 신규 st.HTTP 행1개가 해제 transaction2초 timeout으로30초 만료까지 잔류한 medium 결함을 같은 FIX-2의 공통 admission/Store 원인 분석·수정·검사에 포함한다. 원본 QA 로그/프로브는 fullops-tester에서 읽기만 한다. 새 H2 후보의 독립 QA는 보류이며 DEV 자체 검사로 대체하지 않는다.
