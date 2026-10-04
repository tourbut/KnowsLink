---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-REVIEW — 좁은 timeout검증
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-REVIEW]
summary: 기본10초본문중단 최소수정의 독립 검증과 재사용 경계를 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-REVIEW — 좁은 timeout검증
- ready. 기준cd60e7f87eb5ce137eca887980f232b3f67a18d0 → 대상711f2532be423d1ca7707463a20fdc168f50bece. 원본리뷰0c367301·QAfa168938 불변근거재사용. timeout DEV만 새 대상이다.
- coor/main/origin 병합, 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9 run_8ca8bc058ab7 task/dispatch preamble. 역할 /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops.
- FULLOPS/common0.3.2 README/연결3규칙/project/문서작성·해당검증스킬. 원본QA bodytimeout/held.mjs와 최신TIMEOUT기록 먼저읽고 필요한33줄core변경·관련test/README만깊게확인. 전체Go/SQL/Chrome/CLI 반복금지. 불변성diff와 원래SHA로재사용.
- 소유: 자기review/QA폴더·자기완료로그/inbox만. source·타인박스/원본리뷰·원본QA·PLANS/board·실사용자state/배포/CF/Grok 수정금지.
- selflint --from착수HEAD0, work.pyfinish/archive/commit, worker_done 과제key와 final SHA 전체40자리명시. actualGrok/공개Access/운영배포未검증유지.
- [ ] fullops-review/OCRdelegate. read-only clean detached /tmp/knowslink-timeout-review-711f253, 실행은별도scratch. 원본별도reviewer세션c4c411b8 계속사용가능, 최신DEV실제sessionId는완료로그/원본sessionmetadata에서확인·서로다름기록.
- [ ] reviewkey SAR-MVP-003-BIDIRECTIONAL-TIMEOUT의 모든path/status를처리하되 원본47개코드불변은0c367301증거재사용. 새타이머/abortPromise race/reader cancel/finally/busy·cleanup/비밀/noUnhandled을검토. 원본QA의medium해소판정명시.
- [ ] target711f253 cleanclone lint --fromcd60e7f→lint.json. exactrefs check --task-key SAR-MVP-003-BIDIRECTIONAL-TIMEOUT0. criticalhigh/error차단. narrowreal10sec실행필요시1회만, 반복QA는tester. 결과report/result작성.

## Jev keep·충돌
자기key find/documents-find/context의keep모두확인. 지시 전제와 충돌 — 먼저 확인: 원본fa168938 QA는10초실패, 최신711f253 DEV는고쳤다고주장. 이를독립재검증한다. 구현자latest실제세션91bb4020-55f1-4789-85f4-b53d2b3ce838이며 reviewer와달라야한다. 변경없는Go/SQL/두MCP왕복불변만원본재사용.
