---
title: SAR-MVP-002-BOT-CATALOG-DEV-REVIEW — 실제 Bot Command 등록 후보8e46c5a의 독립 코드·문서 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-REVIEW]
summary: 실제 Bot Command 등록 후보8e46c5a의 독립 코드·문서 리뷰
---

# SAR-MVP-002-BOT-CATALOG-DEV-REVIEW — 앱 Command MCP 등록 준비 독립 리뷰

- 날짜: 2026-10-04; From: coor; 상태 ready.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble.
- 사용자 범위: 실제 Bot 카탈로그 실패 조치·검증·push·이슈댓글. 실제relay/실데이터/유료inference/FullOps업데이트 제외.
- 고정후보 8e46c5a846e6d190e484e48be40b3dc368001a2b; 기준8c95bde. 병합책임 coor/main/origin.

## 먼저 읽기와 기준

FULLOPS.md, project.md, rules/common/README.md와 coding-style/testing/security, 문서작성규칙, docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md, handovers/logs/2026-10-04_to_dev.md 해당 전문, adapters/README.md, scripts/install_bot_mcp.sh, scripts/verify_grok_plugin.py 및 기존 관련 QA. 공통규칙0.3.2와 현재 준비SHA 규칙을 적용한다. 실제 제품은 고정후보를 확인한다.

- 역할/기록 워크트리: ops, /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-REVIEW-*.json. keep 모두 읽는다. 이전 QA는 CLI 설치 성공과 앱 미검증만 증명하므로 이번 앱 등록 성공 근거로 사용하지 않는다.

## 해야 할 일·소유권

- [ ] fullops-review/open-code-review-delegate를 적용한다. 준비 리뷰 키 SAR-MVP-002-BOT-CATALOG-DEV, 기준8c95bde·대상8e46c5a. 읽기전용 detached snapshot /tmp/knowslink-bot-catalog-review-8e46c5a에서 HEAD/clean을 확인한다. 구현자 fcfec882-634c-458a-b16b-b2b285aea791와 다른 실제 reviewer session을 기록한다. snapshot 변경·npm설치 금지; 테스트가 필요하면 별도scratch에서 한다.
- [ ] 모든 변경과 관련 계약을 검토한다. 특히 installer의 다운로드 고정checksum·아키텍처·xz불필요·재실행/기존파일/실패전파·prefix경계·env-i 실제명령·PATH·held 기본값·앱등록 지침을 확인한다.
- [ ] 공식 Team Bots 문서의 Command server 지원과 실제 개인계정의 UI 미검증을 구분한다. 이슈에서 AddMcpServer tool 부재가 있었으므로 같은 Bot에 자연어 Add만 요청하면 승인카드가 반드시 나타난다고 주장하는지 확인한다. 실제 UI 경로는 공식 근거/확인된 단계만 쓰고 앱증상 미해결을 CLI나 로컬status 성공으로 과장하지 않는다.
- [ ] result/report/lint/check를 완성한다. 모든 제외파일도 확인/명시적skipped사유, 실제파일줄 findings. DEV clean8e46c5a에서 lint --from8c95bde를 리뷰 lint.json으로 실행한다. review check exactrefs --task-key SAR-MVP-002-BOT-CATALOG-DEV exit0. critical/high 미해결은 차단이다.

## 완료 기준·보고

공식근거의 적용범위·신뢰경계·기존기준·진짜 앱미검증을 판단한다. 고정제품 코드는 수정하지 않는다. 준비리뷰 디렉터리/자기완료로그만 편집한다. D01-D13 갱신 없음. work.py finish 후 커밋하고 [완료] SAR-MVP-002-BOT-CATALOG-DEV-REVIEW | final SHA <전체Git해시> | 판정·검증·남은계정단계 형식으로 worker_done. ZIP 해시는 SHA256으로 구분한다.
