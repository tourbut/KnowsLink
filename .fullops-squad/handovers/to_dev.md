---
title: SAR-MVP-002-BOT-CATALOG-DEV — CLI 설치 이후 실제 Grok Bot 도구 카탈로그 등록 실패를 진단·수정한다
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-002-BOT-CATALOG-DEV]
summary: CLI 설치 이후 실제 Grok Bot 도구 카탈로그 등록 실패를 진단·수정한다
---

# SAR-MVP-002-BOT-CATALOG-DEV — 실제 Grok Bot 카탈로그 등록 실패 조치

- 작성일: 2026-10-04; From / To: coor / dev; 상태: ready.
- 승인: 사용자 재시험 실패 진단·수정·검증·GitHub push 및 이슈 댓글 인계. 기존 실제 외부 relay·실데이터·업무 효과·유료 inference held 유지. FullOps 업데이트 제외.
- 워크트리 / 브랜치: /home/shin/orca/workspaces/KnowsLink/fullops-dev / fullops/dev.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble을 따른다. 병합 책임 coor, main/origin.
- 기준: 4d6ccfd003ee0823928f044b8528723468fba6ba, 착수한 HEAD를 lint 기준으로 쓴다.

## 현재 상황과 증거

이슈1 추가 댓글5971027648/5971034506을 원본 docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-issue-1.json에 보존했다. 대상은 xAI 공식 Grok Bot https://docs.x.ai/grok-bot.
실제 Bot 컴퓨터에서 CLI1.0.40·Node22.22.2 설치·ZIP hash fb745c66…·doctor 서버1 도구2가 성공했다. 다만 Bot 채팅 GetDynamicTools knowslink 검색0, 설치 MCP14개에 knowslink 없음, knowslink_status 호출 불가다. 앞선 CLI 형식 수정은 성공했지만 앱 도구 등록은 해결되지 않았다. xz 없어서 첫 Node설치 실패 후 xz-utils 설치로 성공했다. 같은 CLI 검증만 반복해 앱 실패를 해결했다고 보고하지 않는다.

## 먼저 읽기와 기준

FULLOPS.md, project.md, rules/common/README.md와 coding-style/testing/security, 문서작성 규칙, orca-agents.md, 해당 이전 DEV/REVIEW/TESTER 실행 기록·QA 보고서, adapters/README.md, adapters/skills/knowslink/SKILL.md, adapters/src/mcp.ts 및 현 설치 설정/패키징 스크립트. 적용 공통규칙0.3.2. diagnosing-bugs·ponytail full을 적용한다. Grok/앱 공식 최신 문서로 지원 API/등록 계약을 확인한다. MCP SDK는 고정1.32.0 Context7 근거를 재사용/필요 시 갱신한다.

coor 공식 조사(2026-10-04): https://docs.x.ai/grok-bot/computer-and-apps 의 Connect an app은 sidebar Marketplace→Browse→Add→auth→chat @connector, installed connectors account-wide를 명시한다. https://docs.x.ai/grok-bot/skills-routines-and-automations 는 Marketplace→Your plugins→Manage plugins and skills→Private skills와 /skill을 명시한다. 위 문서는 임의 stdio CLI 설치가 앱에 등록된다고 명시하지 않는다. 이 사실은 조사 시작점이고 DEV가 앱 실제 지원 계약을 확인한다.

Jev find/documents-find/context는 docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-*.json. keep 전부 읽는다. 지시 전제와 충돌 — 먼저 확인: 이전 INSTALL-FIX 실행 기록. CLI 성공을 실제 Bot 앱 도구노출로 확장한 주장은 이번 재시험으로 반증됐으므로 현재 이슈 원본을 적용한다. 제품 규칙 자체 충돌이면 ask한다.

## 해야 할 일과 소유권

- [ ] diagnosing-bugs의 실제 증상 검증 루프부터 확보한다. 원격 계정 없으면 공식/실제 설치 경계와 필요한 최소 redacted 증거·UI 단계를 정확히 기록한다. 거짓 mock 카탈로그나 CLIdoctor PASS로 앱 증상 해결을 대신하지 않는다.
- [ ] 앱이 사용하는 플러그인/커넥터 등록 계약을 공식 문서/현장 도움말/근거로 확인한다. Bot computer의 CLI ~/.grok와 앱 account catalogue를 구분한다. marketplace Git import/지원 custom MCP/Private skill/HTTP 요구 등 가설은 증거로 판정하고 없는 API·설정경로를 만들어 쓰지 않는다.
- [ ] 사용자 목표(실제 Bot이 KnowsLink 플러그인을 사용)를 충족하는 최소 설치 패키지와 호출 경로를 준비·가능한 범위 구현·검증한다. CLI-only 플러그인으로 조용히 목표를 변경하지 않는다. account MCP 등록이 필수면 정확한 URL/설정/필요 자격 증명과 준비 가능한 전체 산출물을 먼저 만든다. 승인되지 않은 새 공용 서비스 공개/유료/실제 relay 활성화는 하지 않는다. held-only 배포 준비는 허용한다.
- [ ] 앱 직접 권한·원격 등록 API가 없다면 독립 검증 가능한 호출 수단·계정에 등록할 구체적 설치물과 실제 UI의 최소 절차를 완성한다. 실제 Bot 상태 호출을 가능하게 하는 단계가 핵심이다. 대안 Skill/terminal 경로는 앱 카탈로그 등록과 다름을 명시하고 제품 목표 변경이 필요한 경우 ask한다.
- [ ] Node 준비의 xz 의존·Linux아키텍처 및 직전 리뷰low5의 관련 사항을 최소 수정으로 함께 해결한다. 같은 증상을 여러 무관 테스트로 과장하지 않는다.
- [ ] 고정 SHA·검증결과·실제 앱 등록 재시험 명령/단계와 실패 회신 항목을 이슈 댓글 초안으로 제공한다. coor가 리뷰/QA/푸시 뒤 게시한다. 원격 성공 미확인이면 해결 완료라고 쓰지 않는다.

소유권: adapters/와 설치 관련 최소 scripts/Makefile/README, 프로젝트 기술 정본(D10/D12/D13 관련). PLANS/board·다른 역할 인박스는 coor 소유다.

## 완료 기준·검증·산출물

Jev 추천 D13/D10/D12를 실제 변경에 맞게 갱신한다. 기술 분석·설계·구현·단위/통합 검증을 DEV가 같은 과제에서 완료한다. 앱 등록 실패를 드러내는 재현과 적용 전후 신호를 구분한다. 직접 검증 가능한 구현·패키지·설치 지침·댓글 초안을 커밋한다. product lint 및 FullOps lint를 착수 HEAD 대비 exit0으로 완료한다. stderr/원문/credential 보호·KNOWSLINK_MODE=held 보존. coor가 고정SHA의 독립 리뷰·TESTER QA를 배정한다. 기존 코어/Go/UI 불변 증거는 재사용한다. 실제 계정/도구호출 미검증은 별도 후속으로 기록한다.

## 완료 보고

work.py finish로 지시서와 전문을 logs에 보존·인박스 비우기·커밋한다. worker_done body는 반드시 `[완료] SAR-MVP-002-BOT-CATALOG-DEV | 브랜치 fullops/dev | final SHA <40자리 Git SHA> | 변경 ... | 검증 ... | 미검증 ... | 후속 ...`로 쓴다. ZIP은 SHA256으로 별도 표기하여 Git commit SHA 파서와 혼동하지 않는다.
