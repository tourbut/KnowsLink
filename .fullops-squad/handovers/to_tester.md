---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER — 좁은 timeout검증
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER]
summary: 기본10초본문중단 최소수정의 독립 검증과 재사용 경계를 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER — 좁은 timeout검증
- ready. 기준cd60e7f87eb5ce137eca887980f232b3f67a18d0 → 대상711f2532be423d1ca7707463a20fdc168f50bece. 원본리뷰0c367301·QAfa168938 불변근거재사용. timeout DEV만 새 대상이다.
- coor/main/origin 병합, 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9 run_8ca8bc058ab7 task/dispatch preamble. 역할 /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester.
- FULLOPS/common0.3.2 README/연결3규칙/project/문서작성·해당검증스킬. 원본QA bodytimeout/held.mjs와 최신TIMEOUT기록 먼저읽고 필요한33줄core변경·관련test/README만깊게확인. 전체Go/SQL/Chrome/CLI 반복금지. 불변성diff와 원래SHA로재사용.
- 소유: 자기review/QA폴더·자기완료로그/inbox만. source·타인박스/원본리뷰·원본QA·PLANS/board·실사용자state/배포/CF/Grok 수정금지.
- selflint --from착수HEAD0, work.pyfinish/archive/commit, worker_done 과제key와 final SHA 전체40자리명시. actualGrok/공개Access/운영배포未검증유지.
- [ ] fullops-test. detached711f253 clone의 준비Node22에서 원본held.mjs의 stalledpartialrealHTTP/default6인자10000ms 실패항목을 좁게 재시험(2~3회, 각10500ms내TimeoutError). 강제GC의실제socket기본10초를1회, 관련짧은100ms stream/header·cleanup/socketclosed/timer0/noUnhandled·MCPbusy해제/held노네트워크를검증. DEVgreen로그만복사실행해독립대체금지.
- [ ] 원본전체held옵저버재사용가능하나 SQL/fullGo/전체CLI/전체QA는복제금지. 관련diff확인하고 이전fa168938 양방향ID/내용성공증거재사용. 새assert기본timeout·deadline실패시단순medium이라수락으로바꾸지않음.
- [ ] docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md 및 -test/에 짧은명령/exit/측정/판정/재사용경계. 외부실제Grok미검증.

## Jev keep·충돌
자기key find/documents-find/context의keep모두확인. 지시 전제와 충돌 — 먼저 확인: 원본fa168938 QA는10초실패, 최신711f253 DEV는고쳤다고주장. 이를독립재검증한다. 구현자latest실제세션91bb4020-55f1-4789-85f4-b53d2b3ce838이며 reviewer와달라야한다. 변경없는Go/SQL/두MCP왕복불변만원본재사용.
