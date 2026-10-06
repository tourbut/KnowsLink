---
title: SAR-PUBLIC-AGENTS-001-REVIEW — 일반 회원 agent 연결·키·관계 후보의 독립 보안 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-AGENTS-001-REVIEW]
summary: 일반 회원 agent 연결·키·관계 후보의 독립 보안 리뷰
---

# SAR-PUBLIC-AGENTS-001-REVIEW — 일반 회원 agent 연결·키·관계 후보의 독립 보안 리뷰

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

담당 checkout /home/shin/orca/workspaces/KnowsLink/fullops-ops, branch fullops/ops다. 별도 새 Claude Opus5.5high 검토 세션을 쓴다. DEV actual session은 01a10f52-ac0f-75a0-b253-9a926a8e5650이다. 자신의 실제 세션ID를 확인해 result.independence에 기록한다. snapshot /tmp/knowslink-agents-review-d1eef9b는 clean detached d1eef9b이며 read-only로 유지한다. 설치/빌드/검사는 별도scratch의 동일d1에서 실행한다. 결과는 OPS 체크아웃에만 적는다.

- [ ] fullops-review와 open-code-review-delegate를 적용하고 준비된 SAR-PUBLIC-AGENTS-001-REVIEW-review preview/rules/result/report 전 파일을 검토한다. 기준d2..최종d1을 사용하고 일찍 발견한 문제에서 멈추지 않는다.
- [ ] owner/session/recent auth·grant 결속/PoP/재사용/상태·key credential/옛key 전API차단·agent/pair용량/동시성/재시작·rate/CSRF/입력실패·보존과화면개인정보·Node파일/URL 경계를 검토한다. PS와 실제 구현 차이를 확인한다.
- [ ] 최종 제품d1의 FullOps lint JSON을 준비result와 연결한다. 깨끗한 d1 DEV 또는 별도scratch에서명령을 실행해 head/exit일치를 확인하고 readonlysnapshot은변경하지 않는다. SIZE실제범위/DEP testscript-only/SEC 합성token/DESIGN없음근거와 테스트 충분성을판단한다.
- [ ] 모든(path,status)을 reviewed 또는 concrete skipped로 기록하고 원시로그는요약exits와무결성을대조한다. report frontmatter를stamp하고 findings/수락결론·실제미검증·후속을남긴다. 실제독립성·snapshot head/read_only를기록한다.
- [ ] review.py check --key SAR-PUBLIC-AGENTS-001-REVIEW --from d2f7ba5 --to d1eef9b --task-key SAR-PUBLIC-AGENTS-001-DEV를통과시킨다. check는기록검사이며 제품판정과구별한다.

갱신 산출물은 리뷰디렉터리/result/report/lint·자기완료기록이다. 설계문서/제품코드수정은없다. 신규critical/high나구현기준불일치는수락보류로보고한다.

## 완료 보고

검증·관측·finding·fixedSHA·기록HEAD·미검증·산출물·후속을전문으로적는다.
