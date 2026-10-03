---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW — 등록 안내·installer 보완 423db6a의 delta 독립 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW]
summary: 등록 안내·installer 보완 423db6a의 delta 독립 리뷰
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW — 423db6a 수정 delta 독립 리뷰

- coor→ops, 2026-10-04, ready.
- 기록 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops. 복귀 coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7; preamble task/dispatch를 따른다.
- 기준8e46c5a846e6d190e484e48be40b3dc368001a2b → 고정423db6a2a388ea63610462f9d3a5f4c619dd781b. 원본리뷰96d9673·read-only snapshot /tmp/knowslink-bot-catalog-review-423db6a.
- 구현자285f1146-41cc-44cf-9ff4-bbe496639170와 다른 실제reviewer session을 기록한다. 최신 코드·문서 수정은 기존 Bot등록 과제 후속이다.

## 먼저 읽기·기준

FULLOPS.md, 공통규칙0.3.2 README/coding-style/testing/security, project.md, 문서작성규칙, docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md, docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md, scripts/install_bot_mcp.sh와 관련 README/D10/D12/D13. fullops-review/open-code-review-delegate 적용. 최초리뷰의 불변 코드 근거는 재사용한다.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW-*.json. keep 모두 읽고 원본F01-F08 판정은 현재고정SHA와 대조한다.

## 해야 할 일·소유권·완료

- [ ] 준비 OCR key SAR-MVP-002-BOT-CATALOG-DEV-FIX, exact base8e46c5a/head423db6a. snapshot detached/clean/HEAD 확인 후 읽기전용 유지, 결과는 OPS 체크아웃.
- [ ] 이전 F01-F08의 실제 최신수정 판정을 기록한다. 특히 원인단정→가설, 공식TeamBots와개인계정미확정, 이슈 AddMcpServer 부재·조건부 UI·Settings Plugins 삭제·workspace 지속성한계, installer fresh추출/절대prefix/비소유파일보호/검증후교체·실패보존. docclaim과실제교체실패경계가 일치하는지 확인한다.
- [ ] diff 모든files/result 판정과 skipped근거/independence/report/lint/check를 완성한다. 변화없는 코어/Go/UI/원본CLI는 전체 재검증하지 않는다. 필요한 작은회귀는 별도scratch에서 실행하거나 TESTER 기록을 연결한다. fixed423db6a의 cleanDEV에서 lint --from8e46c5a를 새review lint.json에 기록. check exactrefs --task-key SAR-MVP-002-BOT-CATALOG-DEV-FIX exit0.
- [ ] 미해결 critical/high는 차단, 공개댓글 전 수정이 필수인 남은문구를 구분한다. source코드/이전review/타인박스/PLANS/board는 수정금지. 새review directory와자기완료기록만 소유한다. 실제앱계정등록/개인UI/도구카탈로그/모델status는 미검증으로 분리한다.

worker는 완료 기준까지 진행한다. work.py finish·커밋 후 [완료] SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW | final SHA <Git전체해시> | 원본F01-F08판정·검증·미검증 형식 worker_done. 실제relay/유료/FullOps업데이트·이슈외부댓글 금지. 병합담당coor.
