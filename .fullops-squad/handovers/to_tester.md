---
title: SAR-MVP-002-DEV-TESTER — Grok Bot 플러그인 고정 후보의 로컬 패키지와 안전 경계를 독립 검증한다
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-002-DEV-TESTER]
summary: Grok Bot 플러그인 고정 후보의 로컬 패키지와 안전 경계를 독립 검증한다
---

# SAR-MVP-002-DEV-TESTER — Grok Bot 플러그인 로컬 독립 QA

- 상태: ready
- From / To: coor / tester
- 담당: /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester
- 기준 ref: dbdd70086971285b790683f362702e5a9ff55acd
- 고정 제품 후보: 552586b6e886f95bffa9a000a031ea03070afedb
- 복귀: Run run_8ca8bc058ab7, coor term_8b2f910b-c0dc-41de-bceb-03865daa87eb. 실제 Task/Dispatch는 preamble을 따른다.

## 적용 기준과 예외

fullops-common-0.3.2 README와 coding-style/testing/security, FULLOPS.md, project.md와 문서 작성 규칙을 읽는다. 고정 후보의 D03/D05/D10과 사용자 설치 README를 적용한다. 실제 Grok Bot 계정·marketplace 등록·hosted runtime·외부 연결·실데이터·유료 API·운영 변경은 held다. 합성 로컬 검증과 소유 기록·커밋·역할 push는 승인됐다. UI/Go 제품은 불변이며 Chrome QA 9584aaf는 원래 SHA의 근거로 재사용한다.

Jev 근거: docs/evaluations/jev/SAR-MVP-002-DEV-TESTER-{find,documents-find,context,route}.json. 후보는 모두 keep이다. sensitive/oversized 원문은 외부 미전송이며 로컬에서 직접 읽는다.

## 먼저 읽을 문서

- .fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV.md
- .fullops-squad/docs/exec-plans/logs/SAR-MVP-002-DEV/verification.json 및 official-sources.json
- adapters/README.md, adapters/.cursor-plugin/plugin.json, adapters/mcp.json, adapters/skills/knowslink/SKILL.md
- adapters/src/core.ts, adapters/src/mcp.ts, adapters/src/mcp.test.ts, adapters/src/synthetic.ts
- scripts/package_plugin.py, Makefile, adapters/package.json
- .fullops-squad/docs/design-docs/interface-design.md, module-design.md
- .fullops-squad/contexts/tester.md, project.md, 규칙과 현재 인박스

## 갱신할 산출물

없음. 소유 QA/리뷰 기록과 상세 실행 기록만 작성한다.

## 해야 할 일과 파일 소유권

- [ ] fullops-test를 읽고 안정된 고정 후보 552586b의 별도 깨끗한 clone에서 실행한다. 현재 기록 체크아웃은 시험 대상 SHA로 바꾸지 않는다.
- [ ] make install·make plugin으로 ZIP을 만들고 압축 해제 standalone MCP initialize/discovery/status held/pull held와 package에 secret·fixture·node_modules가 없음을 확인한다. DEV ZIP sha256 0e671d1a89c141d896034fff31619b9cd2148b73b567adbc3a97126031989117의 재현 또는 차이 원인을 기록한다.
- [ ] 새 MCP 동작의 자동 경계 검사와 make verify-mvp의 격리 실제 SQL/MCP 경로를 실행한다. loopback 제한·redirect 차단·권한·persist/ACK/claim 실패 전파·gate·최소 denied result·동시 호출 경계에서 실패/정상 동작을 판정한다. make lint/test/build 등 이미 DEV가 검증한 범위는 변경·새 결함 영향이 없으면 원래 SHA 근거로 재사용한다.
- [ ] 기존 배포·실데이터·계정·공유 fixture·사용자 도구 설정은 변경하지 않는다. UI가 불변이면 Chrome QA 9584aaf를 원래 SHA 증거로 재사용하고 새 브라우저 QA로 표시하지 않는다.
- [ ] 시나리오·QA·명령 자신의 종료코드·HEAD·제약을 docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md와 같은 -test/ 증거, 필요 scenarios/에 기록한다. 파이프로 종료코드를 가리지 않는다.
- [ ] work.py finish로 인박스와 완료 전문을 보존하고 커밋·기준 lint·strict 통과·역할 push 뒤 worker_done을 보낸다.

소유: tester 시나리오/QA/증거, contexts/tester.md, docs/exec-plans/phases/SAR-MVP-002-DEV-TESTER.md, tester 인박스/완료 로그. 제품 코드·PLANS/board·운영 배포 변경 금지. 결함은 재현 근거로 DEV에 전달한다.

## 완료 기준과 검증

정상/실패 package·stdio MCP·SQL gate 경계가 독립 실행으로 확인되고 명령·SHA·자기 exit·실패/미실행·재사용 근거가 연결돼야 한다. 실제 Grok Bot 설치/도구 검색/hosted runtime은 계정 없는 로컬 실행과 구분해 held로 남긴다. 결함 0 또는 해결되지 않은 결함의 severity와 재개 조건을 보고한다.

## 완료 보고

첫 줄 [완료] SAR-MVP-002-DEV-TESTER | SHA <전체 기록 SHA>. 대상 후보·QA 결론·검사 exit·ZIP hash·불변 UI 재사용·실제 연결 held와 독립 리뷰 대기를 적는다.
