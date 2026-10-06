---
title: SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 agent 연결·키·관계 후보의 독립 동작 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-TESTER]
summary: 일반 회원 agent 연결·키·관계 후보의 독립 동작 QA
---

# SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 agent 연결·키·관계 후보의 독립 동작 QA

- 작성일: 2026-10-06
- From: coor
- 상태: ready
- repo: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_6895aaf1-7b43-4fe0-a416-76f1255a5946, run_8ca8bc058ab7. Task/Dispatch는 착수 영수증/preamble이다.
- 대상 고정 SHA: d1eef9bb90b9726149980320c42fb1fdbcaf584a. 기준: d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b.
- 적용 기준: fullops-common-0.3.3·FULLOPS·project·document-writing·PS04–07/해당PS11/개인정보·UX04–05. 제품 구현이 늘린 규칙을 정본으로 간주하지 않는다. 기존 실패와 후속 실제 이메일/공개/플랫폼 조건을 보존한다.

## 현재 상태와 먼저 읽을 문서

DEV는 최종d1eef9b에서 product-lint/product-test·verify-mvp·strict 모두exit0, ERROR0/WARNING8/실행불가0을 보고했다. 초기2ac91a7 근거와 최종 malformed rate 보완을 구분한다. 완료로그·빈인박스·후행공백 정규화 사유를 확인했다. 제품 후보는 coor에만 반영했으며 main 수락은 리뷰/QA/UI 대기다.

1. .fullops-squad/FULLOPS.md, project.md, rules/common/README.md 및 코딩/테스트/보안 세 규칙, 자기 contexts/<role>.md.
2. .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md와 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md.
3. .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md, docs/design-docs/interface-design.md.
4. README.md, adapters/README.md, internal/relay/connections.go, member_agents.go, connections_integration_test.go, adapters/src/connect.ts.

Jev code/documents-find와 context는 자신의 과제 키 docs/evaluations/jev/<키>-*.json을 사용한다. 충돌/지시문 경고는 없다. connect.ts의 민감/크기 제외 원문은 안전하게 직접 읽고 비밀값을 출력하지 않는다. UI의 adapters/README omit은 설치 위치 판단에 필요하므로 참조한다. 위 정본과 diff에서 필요한 호출만 좁혀 확인한다.

## 공통 제약과 결과 처리

Workers Free만 허용한다. 유료 전환/구독/초과 과금·실메일·운영 공개/배포·외부 계정 연결·기존 서버/Tunnel/공유서비스 변경은 실행하지 않는다. localhost 격리 fixture·자기 신규 시험자원 생성/회수·기록·커밋·일반 역할push는 허가됐다. 사용자 자료/자격을 변경하지 않는다. 실제 이메일/플랫폼을 fixture PASS로 표시하지 않는다. 새 frontend·의존성·제품수치·scope변경 없이 검수한다.

검사 명령 자신의exit code·fixedSHA·환경·증거를 남긴다. 기존 동일의존성 증거는 원래SHA로 재사용하고 변경영향/증거결함만 재검증한다. 멈춤은 설명되지 않는 실패·범위밖행위·제품 기준 불명확성에 한정한다. 기술 검수는 담당자가 수행하고 제품 판단만 coor를 통해designer에게 질문한다.

자신의 인박스/완료로그/contexts·과제별 보고/증거만 수정한다. 제품 코드는 수정하지 않는다. 문제는 재현/위치/영향으로 보고한다. PLANS는 자기결과만append하며 board는coor 소유다. work.py finish로 지시서/완료보고 전문을 아카이브하고 빈인박스 확인 후커밋한다. 새 기록 고정HEAD의 FullOps lint와 git diff --check, 필요 문서strict 결과를 남긴다. preamble의worker_done을 한 번 보내고idle로 둔다. critical/high 및 필수 실패는 수락을차단한다.
## 담당과 완료 기준

담당 checkout /home/shin/orca/workspaces/KnowsLink/fullops-tester, branch fullops/tester다. 사용자 지정 새 Grok4.7high 세션으로 검사한다. 실행은 별도격리 /tmp/knowslink-agents-qa-d1eef9b의 정확한d1을사용한다. 자기Compose project/DB/포트를사용하고운영자원에접속하지않는다. 제품작업트리 불변을before/after HEAD/status로확인한다.

- [ ] fullops-test를적용해과제별 시나리오·결과/원시로그/명령별exit·SHA를보존한다. 먼저DEV검사내용을확인하고단순복사 대신독립실행/미충족경계를정한다.
- [ ] 실제격리Postgres/HTTP로교차owner·agent 자격owner우회·최근5분·grant owner/agent/client/PoP 바꿔치기·취소/만료/1회완료·동시완료를확인한다.
- [ ] 키3개/회전/선택철회·옛credential의enqueue/pull/ACK/exec/result거부와연결실패시기존키보존, 동일owner2agent의명시적첫수락·타회원초대수신·거절/양측철회·동시수락/옛generation을확인한다.
- [ ] agent200/owner5·activepair400/owner20·pending200/송신owner10/24h 경계/포화/재시작·안전정리·invalid/malformed요청도rate소비·credential회전budget불변을확인한다. 독립검사가부족하면최소회귀를작성하되제품소스는변경하지않는다.
- [ ] Node prepare→owner confirm→complete의로컬 실제CLI·0700/0600·키출력금지·재사용/redirect/파일보존·안내정확성을확인한다. 필요fixture는합성값만쓰고인증값로그를남기지않는다.
- [ ] make lint/test/verify-mvp의필요회귀를고정d1에서한번실행한다. corrected relay-stop/Go/up순서를유지하고관련신원/gate회귀를확인한다. 외부실행/실이메일/플랫폼은SKIP으로구분한다.

과제보고 docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md와 -test/증거, scenarios/과제시나리오를남긴다. 영상/전체UI캡처는designer소유이며중복하지않는다. 담당D10/D12의제품기술/운영정본을QA가수정하지않고추적근거만결과에연결한다. 테스트실패는원문과재실행이유를보존한다.

## 완료 보고

검증·관측·finding·fixedSHA·기록HEAD·미검증·산출물·후속을전문으로적는다.
