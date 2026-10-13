---
title: SAR-REPO-DOCS-001-DEV — README 상세 설치·업데이트와 루트 docs 현행 설명 자료 작성
status: draft
updated: 2026-10-13
owner: dev
tasks: [SAR-REPO-DOCS-001-DEV]
summary: README 상세 설치·업데이트와 루트 docs 현행 설명 자료 작성
attempt: 7644e183f44c4b56b5e25e66f10b05d2
base: 473da85ad6edc56d3833fd42cf540e8d7887920a
subagent_level: standard
test_level: lite
---

# SAR-REPO-DOCS-001-DEV — 현재 README 설치 안내와 루트 docs 설명 자료

- 작성일: 2026-10-13. From/To coor/dev. 상태: ready. Purpose documentation. Test level lite. Subagent standard이나 이번 문서 작업은 직접 수행한다.
- 승인: README.md 수정·루트 docs/ 신규 설명 문서·필요한 기존 adapter 안내의 사실 정정·문서 검증/완료 기록. 제품 코드·의존성·운영 설정·실계정 발송·클라이언트 직접 갱신·배포·삭제/이동 금지.
- 기준: 473da85ad6edc56d3833fd42cf540e8d7887920a. 현재 운영 제품2747fe98653d5c6946f41c219fc256395333d7fe. 현재 실제 수락과 과거 후보/미수락 표현을 구분한다.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev. CLI/model/effort codex/gpt-6.1-sol/high (사용자 Sol 지정 유지).
- 복귀: run_5534969da2d0 / coor term_549a51c8-0972-47ae-b2eb-9b504c3b6330. 실제 주입 task/dispatch 주소로 보고한다. 병합 책임자coor, origin/main.

## 범위와 완료 기준

README에 현재 KnowsLink 목적(사람 친구가 되면 동의된 하위 에이전트들이 대화), 현재 기능/제한, 자세한 설치 절차, 문서 목차를 작성한다. 사용자용 웹 로그인·이름·agent 연결·친구 수락·대화와 개발환경 설치를 구분한다.

신규 설치는 Grok 로컬 Command MCP의 준비 도구/Node22·별도 깨끗한 clone·install_bot_mcp.sh·실제Command/Arguments·public-node/RELAY_URL/새고유연결폴더 설정·명시 Google 연결동의·connected 확인까지, Dots 기존Remote HTTPS /mcp 등록·Google동의·기존사용자와 신규사용자 구분·도구/Events 상태 확인까지 설명한다. 비밀값을 채팅/argv/Git에 넣지 않도록 실제 저장 계약에 맞춘다. 계정정책/UI 미확인 부분을 확정하지 않는다.

기존 설치 업데이트는 Grok 같은 기존prefix/연결폴더/key/agent ID·webhook routine을 보존하는 새bundle 적용과 MCP·정확한해당watcher 재시작, Dots 서버요청부터반영·metadata변경때만기존연결Refresh를 설명한다. installer 출력의 새connect/newfolder 안내를 기존설치에 적용하지 않는다. 친구/agent 재등록은 이번 리팩토링 업데이트에 불필요하다. 수신표시/ACK와 자동답장을 구분한다. 부하·장시간알림/다음Events갱신을 검증됐다고 쓰지 않는다.

루트 docs/README.md(독자별목차), docs/architecture.md(구성·레포지도·회원/친구/agent·송수신/서명/Inbox/ACK·SSE와MCPEvents·보관기한), docs/development.md(필수버전·clone/install·.env.example기반로컬실행·privatePostgres/Google/SMTP모드·build/test/격리QA·실패대처), docs/operations.md(배포/백업·schema/writer보호·동일writer복구·Grok/Dots유지보수·알림문제판별)를 기본으로 작성한다. README의 설치 절차를 docs에 통째복제하지 않고 링크한다. 필요할 때만 파일을 더한다. 설명에 도움이 되는 작은 Mermaid 흐름도는 허용한다. 모든명령/경로/환경변수는 실제코드와대조한다.

README의 과거 후보 전용문구·이미지원하는기능의후속표현·사람friend가아닌agent별관계 안내는 현재와충돌하지않게 정리하고 유효한legacy/합성테스트/운영한계는 보존한다. adapters/README.md의 구식상태가 새문서와 충돌하면 필요한 절만정정하거나 새정본으로 명확히 연결한다. 전체과거기록 재작성/삭제 금지. 사용자 문서에 FullOps task/session/dispatch/개인메일·실제app/구독ID·운영private경로를 노출하지 않는다. 실제운영사실은 날짜와 제품SHA·공개주소 및 검증범위만 쓴다. 내부 심층정본은 상대링크로 연결한다.

## 먼저 읽을 근거

AGENTS.md, FULLOPS.md, rules/common/README.md 및 coding-style/testing/security, project.md, docs/agents/document-writing.md, docs/deliverables/README.md, contexts/dev.md.
README.md, adapters/README.md, Makefile, .env.example, compose.yaml, scripts/install_bot_mcp.sh, adapters/src/index.ts, adapters/src/login.ts, adapters/src/mcp.ts, adapters/skills/knowslink/SKILL.md, .fullops-squad/docs/design-docs/module-design.md, .fullops-squad/docs/operations/user-guide.md, .fullops-squad/docs/operations/ops-guide.md, .fullops-squad/docs/operations/transition.md, .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md, .fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-003-NOTIFY/report.md.

공통 fullops-common-0.3.3와 현재473da85 정본을 적용한다. 기술설명은 기존구현에 국한하고 새제품판단은 하지 않는다. 이전 대화 미참조로 README→docs목차→설치/구조/개발/운영을 따라갈 수 있어야 한다. 모든 D01–13을 새로 만들거나 조립문서 원문을 복제하지 않는다.

Context7 월quota초과를 이전coor/OPS가 확인했다. OpenAI https://developers.openai.com/plugins/deploy/connect-chatgpt 의 Add MCP/Refresh metadata 및 https://developers.openai.com/plugins/build/mcp-events, xAI https://docs.x.ai/grok-bot/team-bots 공식원문을2026-10-13 조회한근거가 OPS보고에 있다. 같은계약은 재사용하고 부족한현재외부API/CLI근거만 Context7 resolve→query(실패면공식문서)로 보완한다. OpenAI안내에는 openai-docs스킬을 적용한다. 인증값/비공개코드 검색전송 금지.

## 검증·완료

수정문서 로컬파일/fragment 링크·코드펜스/예제경로·실제CLI옵션/환경변수·마이그레이션수동순서·새설치와업데이트구분을 검사한다. 네트워크/계정/삭제를 일으키는 설치/배포명령은 검증 목적으로 실행하지 않는다. 기존통과제품게이트는 문서만변경에 재사용하고 전체제품테스트 반복불필요. git diff --check, 문서 정적 lint와 필요한deliverables strict, 명령의 실제소스 대조 및 재현가능한 좁은링크검사면 충분하다. 미실행/지원환경한계는 솔직히기록한다. 단순문서에 새제품테스트/의존성/상설검사도구 추가 금지.

탐색partial/unknown·남은후보는 관련구간만 수동보완하고 원추천과 실제범위를 outcomes에남긴다. 과거전문전수재독을주장하지않는다. 이 지시서는 초기부터완결되어 있으므로 의미변경없는상태/완료체크만 갱신하고, 실질범위변경으로 instruction hash가바뀌면 공식packet 갱신 및 원본보존 후 finish한다.

phase와 qa-reports/SAR-REPO-DOCS-001-DEV/report.md에 변경/검증/한계·실제작성자세션ID를 남긴다. 문서외코드변경0을확인한다. work.py finish로 현재인박스와완료보고를 logs보존 후 비우고 commit한다. 완료본문은 [완료] SAR-REPO-DOCS-001-DEV | 브랜치 fullops/dev | SHA <최종40자리> 형식으로 보고한다. coor가 다른실제세션 fixedSHA snapshot문서리뷰·main일반통합/push를 맡는다.
