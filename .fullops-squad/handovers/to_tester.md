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
- [x] fullops-test. detached711f253 clone의 준비Node22에서 원본held.mjs의 stalledpartialrealHTTP/default6인자10000ms 실패항목을 좁게 재시험(2~3회, 각10500ms내TimeoutError). 강제GC의실제socket기본10초를1회, 관련짧은100ms stream/header·cleanup/socketclosed/timer0/noUnhandled·MCPbusy해제/held노네트워크를검증. DEVgreen로그만복사실행해독립대체금지.
- [x] 원본전체held옵저버재사용가능하나 SQL/fullGo/전체CLI/전체QA는복제금지. 관련diff확인하고 이전fa168938 양방향ID/내용성공증거재사용. 새assert기본timeout·deadline실패시단순medium이라수락으로바꾸지않음.
- [x] docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md 및 -test/에 짧은명령/exit/측정/판정/재사용경계. 외부실제Grok미검증.

## 완료 보고

[완료] SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER. 최종 SHA 전체 40자리는 worker_done 본문에 적는다.
착수 HEAD는 `a9b033d9460c78b06b60699ec9373e89127f829d`다. 판정 후보는 `711f2532be423d1ca7707463a20fdc168f50bece`다. 실행 위치는 `/tmp/sar-mvp-003-timeout-qa`의 detached clone이다.
기본 6인자 `TestTransport`의 실제 HTTP 부분 본문 3회는 `TimeoutError`다. 시간은 10007ms, 10008ms, 10007ms이고 모두 10500ms 안이다. 저장 `timeoutMs`는 10000이다.
강제 GC 100ms 1회는 `TimeoutError` 10007ms다. cleanup은 `TimeoutError` 10017ms, socket close 10023ms, Timeout 0, unhandled 0, 다음 본문 `{"ok":true}`다.
100ms stream은 `TimeoutError` 104ms다. 302 header는 `TypeError` 13ms다. MCP 첫 수신은 `failed` 10026ms이고 겹친 호출은 `busy`이며 이후 state는 `failed`다. held canary hit는 0이다.
`cmd`, `internal`, `db`, `scripts`, Makefile, package, `mcp.ts`, `test-transport.ts`의 diff 종료코드는 0이다. 원본 QA `fa16893870fcaf33e968065e88f2e250a845d09c`의 왕복 ID 네 건을 재사용했다. 이번 SHA에서 그 ID를 다시 만들지 않았다.
SQL, Go, 전체 CLI, 전체 QA, 실제 Grok, Cloudflare, 공개 Access, 운영 배포는 실행하지 않았다. 새 결함과 새 critical/high는 없다. 제품 코드, 원본 QA, 원본 리뷰는 수정하지 않았다.
보고서: `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md`.
lint는 보존 커밋 이후 `--from a9b033d9460c78b06b60699ec9373e89127f829d`로 실행한다. ERROR와 WARNING 개수는 worker_done 본문에 적는다.

## Jev keep·충돌
자기key find/documents-find/context의keep모두확인. 지시 전제와 충돌 — 먼저 확인: 원본fa168938 QA는10초실패, 최신711f253 DEV는고쳤다고주장. 이를독립재검증한다. 구현자latest실제세션91bb4020-55f1-4789-85f4-b53d2b3ce838이며 reviewer와달라야한다. 변경없는Go/SQL/두MCP왕복불변만원본재사용.
