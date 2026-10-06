---
title: SAR-PUBLIC-MESSAGES-001-DEV-FIX — HTTP 입장 H1/M1 원인 수정과 안전 정리 보호 회귀 검사
status: draft
updated: 2026-10-06
owner: coor
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX]
summary: HTTP 입장의 익명 slow-body·정리슬롯 H1/M1 원인 수정과 회귀검사 인계
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX — HTTP 입장 H1/M1 수정·안전 정리 보호

- 상태: ready. 같은 MESSAGES 기능의 보안 리뷰 후속 DEV 원인 분석·기술 계획·구현·테스트다.
- 담당: dev / /home/shin/orca/workspaces/KnowsLink/fullops-dev / fullops/dev
- 복귀: repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb / fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. Task/Dispatch는 preamble 기준.

## 현재 상황·적용 기준과 먼저 읽을 문서

기준 ref는 fixed `09c523da8a3407288d9f5d711e1834af12bc7808`이다. 원DEV final37f9a1e/제품ef5c571과 제품diff0이다. OPS 원리뷰6e1dac284e2fa8b42e16a5e8e59b11e0354a282b는 수락 불가 H1high/M1medium이며 원본checkexit1·형식probe비수락·원실패를 보존한다. 현재QA/UI는09c 원본검사를계속하며수정후영향범위검수는별도다. main/origin68b0d6a 제품통합보류다.

fullops-common-0.3.3의 README·coding-style/testing/security, FULLOPS·project·document-writing·contexts/dev를 같은 전달SHA에서읽는다. SAR-PUBLIC-SERVICE.md PS11/PS04·06·07·08–10와UX06/07, SAR-MVP.md C1–C5, 원DEV실행기록/완료로그와 REVIEW-2-review/report.md·result.json·slot-repro_test.go.txt·slot-repro.log/.exit를읽는다. 코드keep는capacity.go/http.go/store.go/identity.go/member.go의현재인증/세션/rate·관련public_text/회원gate 경계다. Jev find/doc/context 결과도확인한다. 제품규칙변경만coor경유designer판단이다.

## 해야 할 일·소유권

- [ ] OPS의 H1/M1을 직접재현하고 같은과제실행기록에짧은기술계획을적는다. H1: boundedHTTP가 local slot을 body읽기전에잡아익명slowbody4개로정리ACK/revoke/deny/unpair/logout전체429·16개로신규전체429, rate소모0이다. M1: cleanupRequest가path만보고익명도정리수용량을사용한다. 경계함수의모든호출자를확인하고원인을공통경계에서수정한다.
- [ ] PS11상한과현재인가를유지하며본문수신대기가인증된유효정리budget을잠그지못하고익명/잘못된자격/타회원/CSRF/잘못된본문 요청이정리입장을소비하지않게한다. 현재principal·작업권한·schema·rate·body상한/deadline·HTTP16/정리4·공유동시슬롯/크래시회수 의미를일관되게지킨다. 무제한대기/새자원상태를다른층으로옮겨수락하지않는다. 구체구조/API/토큰순서는DEV가정한다. 외부edge완충미확인을해결로대신하지않는다.
- [ ] H1slowbody·M1익명/거부자격·유효ACK/deny/철회/unpair/logout·정리/신규포화·동시성/재시작/기한/본문경계·rate실패집계·기존text/receipt/gate인가회귀를필요범위로자동검사한다. 자기격리fixture만사용·회수한다. 설명되지않는제품실패는분석/수정하고원본실패를보존한다.
- [ ] 관련기술정본/모듈/공용입장설명·실행기록/QA근거/PLANS/context를현재구현에맞게갱신한다. 제품기획·UX규칙·기존OPS결과·다른역할인박스는수정하지않는다. L1전역lock비용·기존L-A/L-B·공개전합성가입unset/운영DB합성owner0는보존한다.
- [ ] 변경관련검사와make lint/test·영향통합verify-mvp·strict/diff를완료하고최종커밋/일반push·동일최종HEAD의FullOps --from09c lint/test(exit/HEAD/ERROR/WARNING/실행불가)를남긴다. 완료metadata/archive마지막SHA까지검사한다. work.py finish로전문아카이브/빈인박스뒤worker_done에fullSHA/근거/미검증을보고한다.

## 완료 기준·검수 인계·산출물

DEV가원인분석/계획/구현/검사를같은과제에서완료한다. 제품완료는별도새fixedSHA OPS독립delta리뷰·TESTER좁은QA·필요영향UI와최종검사를수락한뒤coor가판정한다. 현재H1/M1의원리뷰는원fixed09c에남기고PASS로고치지않는다. 필수실패/critical/high는main수락차단이다. 새fixed에서해소판정을받는다. 기술산출물D10추천과실제영향D03/05/06/07/09를갱신하고변경없는ID는근거를쓴다. 실행기록은docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md다. 새코드규모/DEP/SLOP/DESIGN/SIZE경고처리를전문에기록한다.

## 제약

Workers Free만허용하며유료전환/구독/초과과금은금지다. 서버/Tunnel·frozenwire·제품한도는유지한다. 운영배포·실메일·외부계정·운영자료삭제·실24h·노우↔다닷시험은이번과제가아니다. 비파괴로컬수정/격리검사/자기fixture회수/커밋/일반push는승인됐다. UI변경이면기존memberStyle/UX정본을유지하고별도theme/designlint미구성을기록하며직접시각은designer후속이다. 계획/파일단위별승인을묻지않는다.

## 완료 보고

브랜치/fullSHA·H1/M1원인/수정/재현회귀·원본실패·최종HEAD각명령exit/경고·기술문서/증거·미검증·후속담당/재개조건을전문으로쓸것.

Jev context 충돌 가능성: 원DEV 실행 기록의 자체 PASS는 H1/M1 발견 전 조건이다. 현재 독립OPS H1/M1 미해결과 원본 check exit1이 이번 수정 전제이며 과거 검사를 덮어쓰지 않는다. keep 후보는 지시서의 먼저읽기 범위에 포함한다.
