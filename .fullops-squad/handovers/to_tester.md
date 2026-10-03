---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER — 423db6a installer 파일보존과 새 패키지의 좁은 독립 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER]
summary: 423db6a installer 파일보존과 새 패키지의 좁은 독립 QA
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER — 423db6a의 좁은 경계 QA

- coor→tester, 2026-10-04, ready. 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester.
- 복귀 coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. preamble task/dispatch를 따른다.
- 고정423db6a2a388ea63610462f9d3a5f4c619dd781b, 변경기준8e46c5a846e6d190e484e48be40b3dc368001a2b. 원본 QA87cfb7c0f2361e450df1ec82190185b8a0d2feb2의 불변 증거를 재사용한다.
- 사용자 승인된 같은 Bot등록 조치 후속이다. 실제계정/유료inference/relay/업무효과/FullOps업데이트 제외, 자기임시clone/prefix만 사용. 병합담당coor/main/origin.

## 먼저 읽기·기준

FULLOPS.md, 공통0.3.2 README와 연결3규칙, project.md, 문서작성규칙, 해당 FIX 실행기록·installer와 원본QA의 판정/재사용경계. fullops-test 적용. 관련 좁은 코드만 조사하고 원본전체QA·공식문서 조사·동일 Node다운로드/전체Go/UI를 복제하지 않는다. 새대형harness를 만들지 않고 간단한 실제파일/실패shim으로 판정한다.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-*.json. keep 모두 읽는다. 지시 전제와 충돌 — 먼저 확인: 원본QA는8e46c5a 대상이고 installer는이번에바뀌었다. 최신423db6a가수락대상이며 원본QA는명시한불변근거에만재사용한다. 원본보고의성공을새파일보존경계성공으로 확장하지 않는다.

## 해야 할 일·소유권·완료

- [x] 별도 detached clone423db6a·임시절대prefix에서 실제 installer를 실행한다. 이미검증된 Node22 runtime을 시험prefix에 재사용해도 되며 실제단계를 기록한다. 최종ZIP SHA256 d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609, 출력command/argument가 실제준비파일을 가리키고 레포밖env-i statusheld/tools2/stderr0인지 독립assert.
- [x] old package/knowslink/STALE 및 기존bundle STALE이 새bundle로 섞이지 않음, prefix 무관file/old package 보존, 자기stage 정리, 재실행Ready 정상.
- [x] 상대prefix는 변경전거절, prefix/node 사용자file·링크나 knowslink 비소유폴더 거절/내용보존. bundle 추출 실패와 env-i 검증 실패에서 기존bundle 보존·Ready없음·stage정리. 구현자 19회귀 스크립트만 복사실행한 것으로 독립QA를 대체하지 않는다.
- [x] 변경범위의 실제실패 신호를 기록한다. swap 자체실패에서 무엇이 보존/미보장인지 관찰된범위와 문서주장을 구분한다. 테스트항목을 결과를 보고 줄이거나 의미를 바꾸지 않는다.
- [x] docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md와 -test/에 간결한 명령·exit·assert·판정·재사용·미검증 기록. 실제app등록/개인UI/카탈로그/모델status호출은 미검증. 원본code/리뷰파일/타인박스/PLANS/board 수정금지.
- [ ] 자기기록 lint --from착수HEAD exit0, work.py finish로 아카이브/인박스비움, commit 후 [완료] SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER | final SHA <전체Git해시> | 좁은경계판정·재사용·미검증 worker_done.

새결함은 직접고치지 않고 재현/수락영향으로 coor에 보고한다. 구현변경 후 남은실제계정재시험 단계와 로컬검증을 분리한다.
