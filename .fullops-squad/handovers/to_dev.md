---
title: SAR-UX-ISSUES-002-RECEIPT-DEV — 요청 상태 조회에서 입력·선택·복귀와 오류 재시도를 보존하고 존재 비노출·계정 격리·접근성을 유지한다
status: draft
updated: 2026-10-15
owner: dev
tasks: [SAR-UX-ISSUES-002-RECEIPT-DEV]
summary: 요청 상태 조회에서 입력·선택·복귀와 오류 재시도를 보존하고 존재 비노출·계정 격리·접근성을 유지한다
attempt: c87837bfcc2040289e6986cc4ed64881
base: e939cfbb1939adb4126ad379317c26ec39b692b3
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-RECEIPT-DEV — 요청 상태 조회에서 입력·선택·복귀와 오류 재시도를 보존하고 존재 비노출·계정 격리·접근성을 유지한다

## 목적·고정 범위

사용자 "계속 이어서 이슈 해결하면서 진행해"의 첫 후속 기능이다. 수락된 .fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md U6와 해당 U7·조회 선택기의 U2 이름 표시만 구현한다. 첫페이지 CSS수정 fc4 및 설계11ea와 그 QA/리뷰의 main/origin 통합을 coor가 확인한 뒤 dispatch한다. 기존 CSS/첫페이지/사람친구/권한/서명·프로토콜 계약은 유지한다. #2/#4 전체추가·승인요약·설치온보딩·언어전환은 다음기능이므로섞지않는다.

역할dev/codex gpt-6-astra high는기존사용자고정. Testlite/Subagentstandard. 주frontend-design·기존D04 여섯토큰/44px/대비·단순한국어UI 유지. 완료기준·요구는설계정본 U6-A/U7에서확인하고 DEV가짧은기술계획/모듈·저장·복귀방식/필요테스트를같은과제에서결정한다. coor의후보경로는구현지시가아니다.

## 사용자 완료 조건

- 조회의 제목/진입을 '요청 상태 조회'로 명확히 하고 온라인 확인과구별한다. 기존내에이전트/관리진입에서접근하며 주요메뉴세개를늘리지않는다.
- 처음·성공·없는/타인/접근불가/보존종료·형식오류·429/일시장애에선택한내에이전트와안전한입력 ID를보존하고같은화면에서수정/재조회한다. HTMLescape·안전길이/형식거절·소유에이전트목록만표시·비공개존재비노출/동일응답경계를유지한다.
- 성공시전달/처리/관련답장을분리하고 queued·ACK·답장·온라인을혼동하지않는다. 선택기는허용된이름/종류/안정구분값으로동명을구별하고권한없는이름을새로조회하지않는다.
- 새로고침·뒤로가기·재조회·복귀시해당폼위치와입력을복구한다. 세션만료후동일회원의승인된조회문맥만복원하고 다른회원/익명에이전계정정보를복원하지않는다. 장기조회이력/자동발송/자동polling/새권한확대없음. 보안상안전한기술방법을찾기어려우면구체적문제와기존정본충돌만coor에알린다.
- 오류필드/결과로의키보드focus·라벨/skip/main/현재위치·안전한돌아가기·긴한글/ID/320·390·1440/200%의필요범위를검증한다. 내부기술JSON은접고사용자다음행동은보이게한다.

## 소유권·검증·인계

기본후보 internal/relay/member_receipt.go/member_pages.go 및실제필요한회원진입/안전한복귀·인가뷰/기존테스트/helper다. 범위안의필요콜러·테스트·fixture·기술문서수정은DEV가판단해모두완료한다. 기존러너/fixture/매개변수테스트/표준라이브러리우선, 새의존성/불필요추상화/설정/JS/CSP완화금지. 핵심오류/계정격리/비노출의의미있는실행검사를남기고문구/구현문자열만복제한테스트는쓰지않는다.

related검사는개발중짧게, 안정후등록make lint/test전체gate한번과필요DB/auth실행을기존격리Linux/helper로수행한다. 최신제품소스가바뀌므로fc4제품PASS를새실행으로재명명하지않는다. 현재SHA FullOps lint의원exit/status/WARNING/불가사유를보존하고ERROR를고친다. 최종source/환경·새QA에필요한기능fixture/origin/허용동작·계정A/B/C합성조건·원실패·복구를한번준비해tester/designer에공유한다.

own same-key route/find/context/packet/history/outcomes·QA/phase/context/inbox/log는추적증거소유권에포함하며제공미커밋스냅샷도원본/이력보존해완료커밋에포함한다. 변경에영향받는D05/D10/D11과필요D09/사용안내의정본만갱신한다. 다른role/다른리뷰결과/PLANS/board/GitHub변경금지. 개발산출물은레포에커밋하며임시검증자료는별도private task디렉터리에보존한다. 삭제/이동·forcepush·실계정/실메일/상대발송·운영/배포금지.

먼저 읽기: FULLOPS/common 세 규칙/project/docwriting/contextsdev, 새 설계 U6/U7/U2와 D04, member_receipt.go/member_pages.go/human_http.go/human_pages.go, 현재 회원 auth/조회 회귀 테스트, interface/module/테스트 정본. find/context의 큰 파일·partial/unknown은 필요한 원문을 직접 확인한다. 외부 SDK 근거는 Context7을 사용하고 실패하면 공식 자료로 확인한다.

완료: 기술계획→구현/실패경계검사→현재SHA gate/보고→packetoutcomes→work finish아카이브/빈인박스→최종commit/실제worker_done. 후속코드리뷰/독립QA/직접시각/main통합은coor책임. 복귀run_756d0278e72e/coor term_5bf18065-a830-4e45-aeba-ae852666b40e, repo0b08ec4c-9e3a-4613-8197-5a835e545335/dev.
