---
title: SAR-PUBLIC-AGENTS-001-UI — 일반 회원 agent 연결·키·관계 UX04–05 직접 검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI]
summary: 일반 회원 agent 연결·키·관계 UX04–05 직접 검수
---

# SAR-PUBLIC-AGENTS-001-UI — 일반 회원 agent 연결·키·관계 UX04–05 직접 검수

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

담당 checkout /home/shin/orca/workspaces/KnowsLink/fullops-designer, branch fullops/designer다. 새 Codexgpt-6.1-solhigh세션으로UX04–05와변경된자기홈을직접검수한다. 제품수치/코드/권한변경없이판정한다. 실행fixture는정확한d1의별도 /tmp/knowslink-agents-ui-d1eef9b다. 자기전용SMTPsink·DB·relayport·Compose project를쓴다.

- [ ] DEV인계의일반회원localhost화면을합성메일로구동하고등록/재확인부터두agent연결·키관리·명시적관계까지브라우저에서직접확인한다. CLI/자동QA결과를시각관측으로바꾸지않는다.
- [ ] 새전용브라우저page를만들어사용자기존CF/계정탭은변경하지않는다. 이검수동안coor/QA는브라우저조작을하지않는다. Orca캡처장애시같은실제localhost UI의별도Chromium/Playwright를허용하되engine/PNG/실제직접관측과장애근거를분리한다.
- [ ] UX04 연결대기/준비지문/승인/완료/취소/만료/미지원·권한/실행위치·키활성/철회와재확인오류, UX05 초대/수락/거절/양측철회/수락전상태와이메일비노출을확인한다.
- [ ] 기존memberStyle정보구조·label/키보드/읽을수있는오류/긴ID복사·색외상태표시·세션종료와agent철회구분을관측한다. 정지캡처로판정가능한항목만PNG를만든다. 인증메일/코드/grant/token/privatekey는촬영전마스킹하고마스킹방법을기록한다.
- [ ] SHA·fixture/port·관측/캡처경로·PASS/FAIL/미검증을docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md에남기고D04원천/과제보고를규약대로stamp한다. 기존UX01–03증거는동일의존성확인시원래SHA로만재사용한다.

별도theme/Tailwind/shadcn/디자인전용lint는없으므로해당없음과실제검수영향을보고한다. 채팅버블/composer/긴타임라인/무조건성공배너를쓰지않는다. 실제메일로그인/공개/노우↔다닷은미검증이다. 기술결함은위치와관측을보고하고DEV에인계한다. 과제별mockup검수보고/캡처/자기기록만수정한다.

## 완료 보고

검증·관측·finding·fixedSHA·기록HEAD·미검증·산출물·후속을전문으로적는다.
