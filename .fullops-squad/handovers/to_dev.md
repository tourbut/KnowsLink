---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT — 기본 body timeout 수정
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT]
summary: 독립 QA의 기본 10초 응답본문 중단 실패를 최소 수정한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT — 기본 body timeout 수정
- ready coor→dev. /home/shin/orca/workspaces/KnowsLink/fullops-dev fullops/dev, 복귀coor /home/shin/orca/workspaces/KnowsLink/fullops-coor term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9 run_8ca8bc058ab7. Task/dispatch preamble. 병합coor/main/origin.
- 제품고정cd60e7f87eb5ce137eca887980f232b3f67a18d0, 리뷰0c36730106025be95709236b63f080a510a4aebe, 최신QA fa16893870fcaf33e968065e88f2e250a845d09c. 실제배포미변경. 이번기준착수HEAD를fullSHA기록.
- FULLOPS.md/common0.3.2 README·연결3규칙/project.md/문서작성/ponytail full/diagnosing-bugs 먼저읽음. QA보고 bodytimeout과 held.mjs·body-timeout-followup.txt 재현근거를먼저확인. Jev keep/충돌확인.
- [ ] 기본6인자 TestTransport timeout10000이 첫청크후 body멈춤에서 15001ms hung, signalabort10001ms에도 readerpending이라는 실패를 재현. 100ms합성시험통과를기본성공으로대체금지.
- [ ] timeout경로만최소수정. body읽기까지총deadline/종료/cancel/후속busy해제를 보장. 비밀출력·unhandledrejection·타이머잔존·cleanup문제방지. 새대형추상화금지.
- [ ] 실제기본10000ms stalledbody 회귀와관련header/body/64KiB/redirect/MCPbusy·no-network held 검증. 반복필요하면 boundedfewruns, 성공만선별금지. 원본QA/리뷰/로그변경금지. 관련source·test·최신정정실행기록만소유. 다른low5·timeout무관제품확대/타인인박스/PLANS/board금지.
- [ ] DEV원본댓글/문서10초주장을실제코드와맞춤. 원본문서기록은보존하고정정링크. 원본두MCP/SQL왕복성공은무관source불변확인후재사용. 실제Cloudflare/Grok/privatecredential배포·FullOps업데이트금지.
- [ ] lint --from착수HEAD0, Jevscore, work.pyfinish/archive/commit, worker_done final SHA40자리와재현/관련결과. 새fixedSHA delta독립리뷰·좁은TESTER 재시험은coor배정.

## Jev 근거
자기 key find/documents-find/context의 keep 전부확인. 지시 전제와 충돌 — 먼저 확인: 원본DEV기록은10초body중단성공, 최신QA는기본10000ms실패다. 성공선언 대신 actual재현조건/fetch·stream·signal을진단하고소스또는시험결함의근거를확인한다. 원본기록전문보존.
