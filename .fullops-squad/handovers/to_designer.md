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

- [x] DEV인계의일반회원localhost화면을합성메일로구동하고등록/재확인부터두agent연결·키관리·명시적관계까지브라우저에서직접확인한다. CLI/자동QA결과를시각관측으로바꾸지않는다.
- [x] 새전용브라우저page를만들어사용자기존CF/계정탭은변경하지않는다. 이검수동안coor/QA는브라우저조작을하지않는다. Orca캡처장애시같은실제localhost UI의별도Chromium/Playwright를허용하되engine/PNG/실제직접관측과장애근거를분리한다.
- [x] UX04 연결대기/준비지문/승인/완료/취소/만료/미지원·권한/실행위치·키활성/철회와재확인오류, UX05 초대/수락/거절/양측철회/수락전상태와이메일비노출을확인한다.
- [x] 기존memberStyle정보구조·label/키보드/읽을수있는오류/긴ID복사·색외상태표시·세션종료와agent철회구분을관측한다. 정지캡처로판정가능한항목만PNG를만든다. 인증메일/코드/grant/token/privatekey는촬영전마스킹하고마스킹방법을기록한다.
- [x] SHA·fixture/port·관측/캡처경로·PASS/FAIL/미검증을docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md에남기고D04원천/과제보고를규약대로stamp한다. 기존UX01–03증거는동일의존성확인시원래SHA로만재사용한다.

별도theme/Tailwind/shadcn/디자인전용lint는없으므로해당없음과실제검수영향을보고한다. 채팅버블/composer/긴타임라인/무조건성공배너를쓰지않는다. 실제메일로그인/공개/노우↔다닷은미검증이다. 기술결함은위치와관측을보고하고DEV에인계한다. 과제별mockup검수보고/캡처/자기기록만수정한다.

## 완료 보고

검수 역할 작업은 완료했다. 대상 제품 fixedSHA는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`이며 기준 ref는 `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`다. 기록 시작 HEAD는 `d111fd4fde1538492c3c335e78ed510408718b2e`다. 최종 기록 HEAD는 worker_done의 SHA와 검증 JSON으로 고정한다. 제품 코드·제품 규칙은 바꾸지 않았다.

일반 코드 가입·재로그인·재확인에서 두 자기 agent 연결과 별도 교차 회원 agent를 실제 localhost 브라우저로 검수했다. 대기·prepared 지문·승인 후 미연결·실제 complete·취소·fixture 만료·미지원/권한 거부·회전·선택 키 철회·agent 철회와 명시 관계 pending/active/denied/revoked를 직접 관측했다. 교차 회원 수신 수락과 active 관계의 발신/수신 양측 철회도 확인했다. 상대 이메일은 관계 화면에 없었다. 로그아웃 뒤 재로그인한 자기 agent 키는 활성으로 남았다.

지정 키보드 이동·label·readonly ID 전체 선택과 상태 문자 표시는 PASS다. UI 전체 시각 수락은 FAIL/보류다. medium F-UI-01은 390px 홈에서 scrollWidth 596px인 지문 overflow다. medium F-UI-02는 오류 문구가 홈/재확인을 요구하지만 링크나 다음 컨트롤이 없는 화면이다. low는 종료 연결의 취소 버튼 잔존과 관계 기한의 UTC raw 표시 및 작은 native select다. 실제 안전 거부를 오류 UX FAIL과 구분한다. 이번 관측에서 새 critical/high는 없으며 다른 리뷰·QA의 finding을 닫지 않는다. 기술 수정 방법은 DEV 책임이다.

새 fixture `/tmp/knowslink-agents-ui-d1eef9b`, Compose `knowslink-agents-ui-d1eef9b-8398`, DB 57187/relay 59219/SMTP 52591을 사용했다. Orca 자기 신규 page만 조작했다. Orca 1×1/blank 캡처와 독립 CDP timeout을 보존하고, 허용된 Playwright 1.60.0/Chromium 147.0.7727.15의 실제 localhost UI로 보완했다. UI PNG 31개를 직접 열었으며 장애 PNG 2개는 PASS 근거에서 제외했다. PNG는 촬영 전에 이메일·코드·연결 수단을 불투명 mask로 가렸다. private key·token·credential·실제 코드를 기록하지 않았다.

만료와 오래된 재확인은 own relay를 중지한 전용 DB 시간 필드 설정으로 관측했다. 실제 TTL 대기와 장시간 타이머는 미검증이다. 자기 프로세스·Compose DB/volume/network·합성 메일/키/자격 파일·Orca page를 회수했다. 기존 탭·사용자 자격·운영 자원·공유 서비스를 보존했다.

산출물은 `docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md`와 PNG 폴더, `docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI.md`, `docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/`의 manifest·fixture/정리/검증 로그다. D04 원천에는 이번 결과 링크를 추가하고 stamp했다. 자기 context와 PLANS의 자기 결과를 append했다. 현재 인박스 전문은 work.py finish로 날짜 로그에 보존하고 빈 인박스를 확인한다.

제품 무변경 diff·문서 strict·Git 공백·PNG/JSON 형식·마스킹을 확인한다. 새 기록 고정 HEAD에서 FullOps lint를 기준 ref로 실행하고 원 종료코드와 ERROR/WARNING/실행불가를 검증 JSON에 남긴다. 등록 product-lint/product-test는 기준 ref 이후 기존 제품 변경으로 실행하며 이번 역할의 제품 코드 변경은 아니다. theme/Tailwind/shadcn·디자인 전용 lint는 해당 없음이다. Go CSS를 자동 DESIGN 검사로 수락하지 않는다.

실제 일반 이메일·운영 공개·노우↔다닷·외부 앱/OAuth·실제 왕복·자동 wake·24h 초대 만료·전체 API 철회/용량/동시성·모든 세션 종료·스크린리더·OS clipboard는 미검증이다. 기존 UX01–03 증거를 이번 실행 PASS로 재사용하지 않았다. 변경된 홈은 새로 검수했다. DEV가 medium 두 건을 수정한 고정 SHA를 제공하면 designer가 직접 재검수한다. coor는 독립 QA·리뷰·운영 공개 조건을 유지한다.
