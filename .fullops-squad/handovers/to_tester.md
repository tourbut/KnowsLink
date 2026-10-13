---
title: SAR-REFACTOR-001-TESTER — Sol 리팩토링 최종 고정 SHA의 독립 회귀 검증
status: draft
updated: 2026-10-13
owner: tester
tasks: [SAR-REFACTOR-001-TESTER]
summary: Sol 리팩토링 최종 고정 SHA의 독립 회귀 검증
attempt: 157f6cde8bc24988ad213fab2461510c
base: 0eb7735f0963272b62266992791f9f09152f9266
subagent_level: standard
test_level: lite
---

# SAR-REFACTOR-001-TESTER — 고정 Sol 리팩토링의 독립 회귀 검증

- 상태: ready. Purpose: review/test. Test level lite. Subagent level standard이나 이번 좁은 QA는 직접 수행한다.
- 기준 diff: 0eb7735f0963272b62266992791f9f09152f9266 → cd4dae810ae262759bcff29895165f1a64f8c51a.
- 실제 제품 후보: cd4dae810ae262759bcff29895165f1a64f8c51a. 준비 HEAD는 문서만 다를 수 있으므로 제품 동일성을 직접 확인한다.
- 담당 repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester.
- 복귀 coordinator term_549a51c8-0972-47ae-b2eb-9b504c3b6330 / run_5534969da2d0. 주입된 실제 task/dispatch 주소로 보고한다.
- CLI/model/effort: codex/gpt-6.1-sol/high. 사용자 Sol 지정 유지. 구현자01a13509-e1a3-7b62-bf9d-5a48405d706a와 다른 실제 본인 세션을 기록한다.

## 적용 기준과 먼저 읽을 문서

공통 fullops-common-0.3.3 README/coding-style/testing/security, project.md, AGENTS.md, docs/planning/product-specs/SAR-HUMAN-PEER.md, docs/design-docs/module-design.md, docs/exec-plans/phases/SAR-REFACTOR-001-DEV.md, DEV-test/report.md와 gate/remote.py를 읽는다. packet partial/unknown을 일반 rg로 보완한다. 구조·보안 기준 변경은 없다. 리팩토링·코드 리뷰에 Context7는 불필요하며 실제 외부 API 근거가 필요한 때만 조회한다.

## 해야 할 일과 소유권

- [ ] 고정 diff와 실제 호출자/테스트를 읽고 공유 connectionResponse의 retry/error 우선순위, humanAgents의 owner·철회·정렬·차단 경계, validText의 UTF-8/공백/surrogate·confirm gate, MCP 포장/held wire 보존을 독립 확인한다. API·CSRF·서명·Inbox/ACK·TTL·MCP Events와 template/UI의 변화 유무를 확인한다.
- [ ] 직접 회귀 검증을 한다. DEV exactcd4 Linux make lint/test/build EXIT0·FullOps errors0/warnings3를 실제 원본 로그와 SHA로 확인하고 재사용 여부/사유를 적는다. 같은 전체 등록 검사를 불필요하게 반복하지 않는다. 독립 focused Go/TS 경계 검사를 실행하고 종료코드를 보존한다.
- [ ] Windows Go race는 cgo 부재, 보조 Node24는 지원환경 아님을 구분한다. 기존 검증된 SSH helper와 고유Linux clone을 읽기/테스트에 재사용할 수 있다. AppData/Local/KnowsLinkDevTools/refactor001/cd4dae810ae262759bcff29895165f1a64f8c51a.base의 clone에서 HEAD·clean·Node22/npm10 확인 뒤 정확한 후보 focused 검사를 실행한다. DEV가 같은 clone을 바꾸고 있다면 별도 고유 clone으로 검증한다. 비밀값·환경 전체·인증URL은 출력하지 않는다.
- [ ] qa-reports/SAR-REFACTOR-001-TESTER-test/report.md에 검사 명령·종료코드·고정 SHA·PASS/FAIL/not_run·발견severity·호환 보존·검증 한계를 기록한다. packet outcomes를 완료하고 work.py finish·worker_done까지 처리한다.

## 완료 기준과 제약

제품 코드는 읽기 전용이고 변경하지 않는다. test helper/보고서/시나리오와 tester context만 소유한다. C/H 발견 시 coor에 즉시 보고하고 판정에서 숨기지 않는다. UI template/CSS/문구 변화0이면 기존HTTP/template회귀로 충분하며 새디자인검수는 해당없음이다. 독립 코드 리뷰는 root가 fixed snapshot에서 별도로 수행 중이다.
운영·실계정 발송·배포·키변경·로컬Docker·공유자원 변경·Windows 삭제/이동은 금지한다. 기존 Linux test 자체 mktemp 정리 trap은 필수 검사 승인 범위다. 권한 밖 삭제/새정리 명령은 하지 않는다. hook 차단을 우회하지 않는다. 새 의존성·보안 완화·테스트 약화는 금지다. 모델 capacity 대기로 종료됐으면 실패/미완료로 보고하고 기존 결과를 보존한다.

## 산출물·완료 보고

갱신할 제품 정본 없음. 독립 QA 보고와 기존 gate 검증 근거만 남긴다. 실제 세션ID·후보SHA·기록SHA·검증종료코드·severity·not_run·packet outcomes를 기재한다. 코드를 수정하지 않고 완성된 QA 결과를 커밋하고 work.py finish 후 현재 주입된 parent 주소로 worker_done을 보낸다.
