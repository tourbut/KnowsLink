---
title: SAR-PUBLIC-MESSAGES-001-DEV-FIX-3 — 일반 사용자 정리 동작과 snapshot 일관성의 최종 회귀 수정
status: draft
updated: 2026-10-07
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX-3]
summary: 일반 사용자 정리 동작과 snapshot 일관성의 최종 회귀 수정
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX-3 — 일반 사용자 정리 동작과 snapshot 일관성의 최종 회귀 수정

## 대상·기준·복귀

DEV는 fixed758e9f638501f11000ed17c558a9eb54b4350eb3의 admission/store를 분석하고 최소수정·회귀검사·기술산출물을 같은 과제에서 끝낸다. 기술방법은 DEV가 결정한다. 제품한도/완료조건 변경은 coor 경유 designer에게 질문한다. 기준 fullops-common-0.3.3의 README/coding-style/testing/security, FULLOPS/project/document-writing/orca-agents/contexts/dev를 읽는다. 기술정본 D03/D05/D06/D07/D09/D10, 제품 SAR-PUBLIC-SERVICE PS11/운영 기본값, DEV-FIX-2 phase/완료전문, Sol 진행분 provenance/sol-narrow-tests.log/sol-review-probe_test.go.txt를 읽는다. 새 의존성 없으면 기존 표준 Go 근거를 재사용한다.
복귀 coor term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. 새 세션/route Opus5.5high다. 이전 세션05c78055 자체PASS는 독립QA가 아니다.

## 목표·완료 조건

- 실제 관측을 분석한다. 동일 DB commit 시각에서 늦은 이전snapshot이 최신snapshot을 덮는다. 다른instance의 새자격/재시작nil snapshot에서 신규16포화중 유효자기정리가 갱신전429·실제1s sweep후200이다. 같은회원다른agent ACK가 handler invalid_lease409일 때 회원자기 key-revoke가 로컬입장에서429이며 정상재시도18건 모두 미입장이다. 원관측 프로브의PASS는 관측성립이지 제품적합성PASS가 아니다. 원증거를 수정하지 않는다.
- 정리분류와 실제handler의 권한/lease/상태, snapshot저장/refresh의 전체caller를 대조한다. 유효사용자의 자기철회/로그아웃/deny/unpair/ACK가 부당하게 거절되지 않도록 공통 원인을 고친다. 신규16/정리4·rate·body/기한·권한/CSRF/현재키/관계세대·crash회수 기준은 낮추지 않는다. 무제한fallback/queue로 문제를 옮기지 않는다.
- 자체앱의 정상사용자 기능회귀를 추가하고 원758 RED/새코드PASS를 구분한다. 저장된 기존프로브는 근거로만 읽는다. 차단된 응답 내용을 복원/우회하거나 외부서비스 탐색·악용절차·공격도구를 만들지 않는다. 자기격리 Postgres/로컬 앱만 사용한다. 상태순서·다중instance·새자격·정리caller·종료실패 회수와 영향 Node/MCP 회귀를 완료한다. 무관한기능/의존성/추상화를 추가하지 않는다.
- 최종 code/기록HEAD의 make lint/test·영향 verify-mvp·FullOps lint --from758·strict·diff 결과에 SHA/exit/ERROR/WARNING/실행불가·규모/DEP/SLOP/DESIGN 처리와 원RED/중간실패를 남긴다. phase SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md 및 기술영향 D03/D05/D06/D07/D09/D10을 갱신한다. D12/D13은 미검증·이행조건 연결만하며 운영수락은 만들지 않는다.
- 현재 inbox에 완료 전문을 쓰고 work.py finish로 archive/빈inbox를 확인한다. 최종 cleanHEAD·일반push·실제 worker_done으로 fullSHA/검증/제약/후속을 보고한 뒤 idle한다.

## 소유·제약·후속

제품/자기tests·기술정본/phase/context/PLANS·자기inbox/logs만 수정한다. WorkersFree·유료전환/구독/초과과금금지·기존 공유서버/Tunnel 유지다. 운영배포/실메일/외부계정/운영자료삭제는 없다. 자기fixture만 회수한다. 별도 허용범위 정적리뷰·사용자 승인 Sol 독립QA·designer UI 영향 확인 뒤 coor가 main merge/origin/main 일반push/조상 확인/유휴clean 역할sync를 수행한다. 실메일·운영공개·실24h·노우↔다닷은 별도 미검증이다. coor가 공개준비를 병행한다. 과거 실패·플랫폼 표시차단·334원본hash를 보존한다.
