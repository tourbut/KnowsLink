---
title: SAR-MVP-002-BOT-CATALOG-DEV-TESTER — 실제 Bot Command 등록 후보8e46c5a의 독립 설치 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-TESTER]
summary: 실제 Bot Command 등록 후보8e46c5a의 독립 설치 QA
---

# SAR-MVP-002-BOT-CATALOG-DEV-TESTER — 앱 Command MCP 등록 준비 독립 설치 QA

- 날짜: 2026-10-04; From: coor; 상태 ready.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble.
- 사용자 범위: 실제 Bot 카탈로그 실패 조치·검증·push·이슈댓글. 실제relay/실데이터/유료inference/FullOps업데이트 제외.
- 고정후보 8e46c5a846e6d190e484e48be40b3dc368001a2b; 기준8c95bde. 병합책임 coor/main/origin.

## 먼저 읽기와 기준

FULLOPS.md, project.md, rules/common/README.md와 coding-style/testing/security, 문서작성규칙, docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md, handovers/logs/2026-10-04_to_dev.md 해당 전문, adapters/README.md, scripts/install_bot_mcp.sh, scripts/verify_grok_plugin.py 및 기존 관련 QA. 공통규칙0.3.2와 현재 준비SHA 규칙을 적용한다. 실제 제품은 고정후보를 확인한다.

- 역할/기록 워크트리: tester, /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-TESTER-*.json. keep 모두 읽는다. 이전 QA는 CLI 설치 성공과 앱 미검증만 증명하므로 이번 앱 등록 성공 근거로 사용하지 않는다.

## 해야 할 일·소유권

- [ ] fullops-test를 적용한다. 별도 detached clone8e46c5a와 자기 임시 KNOWSLINK_PREFIX를 사용한다. 사용자 ~/.grok·실제/workspace/.knowslink 및 타역할 체크아웃 변경 금지.
- [ ] 신규 scripts/install_bot_mcp.sh를 첫실행·재실행한다. 다운로드checksum·Node22·npm10·패키지검증·env-i PASS·출력 등록값을 실제명령exit와 함께 기록한다. ZIP 예상 b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad. 준비node/bundle을 레포밖cwd·env-i로 직접 실행해 tools2와 statusheld/stderr0을 독립assert한다.
- [ ] 기존 prefix의 무관 marker 및 사용자설정이 보존되는지 확인한다. 설치/재실행 후 command/argument 절대경로가 실제준비파일을 가리키는지 확인한다. xz없는 환경·미지원uname 실패1·다운로드/체크섬/빌드 실패전파에서 거짓Ready가 없는지 관련 낮은비용 경계만 검증한다. 실제arm64호스트 없으면 미검증으로 분리한다.
- [ ] Grok 검증 스크립트 변경 F01-F03의 오류출력·GROK_CONFIG계열 제거·PASS문구를 작은 실패shim/관련스모크로 확인한다. Node20실패·기존Go/UI/relay/fullCLI설치 증거는 변경없어 재사용하고 전체QA를 복제하지 않는다.
- [ ] docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md 및 -test/에 명령/exit/입력/판정/재현 증거를 작성한다. 실제앱등록/Command개인계정지원/카탈로그 노출/Bot세션status는 원격계정 없으므로 미검증이다. 로컬직접status를 앱성공으로 기록하지 않는다.

## 완료 기준·보고

새installer와 실제app-command 경계·실패전파를 독립 검증한다. 제품코드/리뷰파일/타인박스/PLANS/board 수정 금지. 자기QA와테스트/완료로그만 소유, 결함은 재현으로 보고한다. D01-D13 갱신 없음. 자기 lint --from착수HEAD exit0과 work.py finish/아카이브/커밋 뒤 [완료] SAR-MVP-002-BOT-CATALOG-DEV-TESTER | final SHA <전체Git해시> | 검증·결함·미검증 worker_done. 유료inference/relay 호출하지 않는다.
