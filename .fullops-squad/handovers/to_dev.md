---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX — 조건부 리뷰의 앱 등록 안내와 설치 파일 보존 경계를 보완한다
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX]
summary: 조건부 리뷰의 앱 등록 안내와 설치 파일 보존 경계를 보완한다
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX — 등록 안내·installer 경계 보완

- 작성일 2026-10-04; coor→dev; 상태 ready.
- 작업: /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev. 복귀 coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch preamble.
- 승인/범위: 사용자의 실제 Bot 재시험 실패 조치의 같은 후속. 리뷰 지적 수정·관련 검증·push/댓글 준비. 실제계정 등록/릴레이·실데이터·유료inference·FullOps업데이트 제외.
- 기준제품8e46c5a846e6d190e484e48be40b3dc368001a2b, 리뷰96d9673e1a6dd80db1e5d9c5166efb64aafa0346. 착수HEAD를 lint 기준으로 기록.
- 병합 책임 coor, main/origin.

## 먼저 읽기·기준

FULLOPS.md, project.md, rules/common/README와 연결3규칙(0.3.2), 문서작성규칙, docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md 및 result.json, 해당 DEV 실행기록/README/installer, 공식 team-bots와 settings-and-notifications. ponytail full, 필요한 fullops-deliverables 및 writing-for-agents를 적용한다. 원본 QA8e46c5a는 별도TESTER에서 진행중이고 결과는 새 수정의 불변 부분에만 재사용한다.

Jev find/documents-find/context는 docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-FIX-*.json. keep 모두 읽는다. 기존 DEV 실행기록/README의 원인단정과 신규리뷰 판단이 충돌하면 현재 리뷰의 관측/미확정 구분으로 정정하고 기존전문은 보존한다.

## 해야 할 일·소유권

- [ ] F01/F02/F04 정정. 이슈1에서 AddMcpServer 호출/검색 도구 부재가 이미 확인됐다. 자연어 Add→승인카드, 개인계정 Command 제공은 공식근거가 없다. 문서/ops-guide/전환/모듈/댓글초안에 이 사실을 넣고 확인되지 않은 카드 노출을 기대값으로 단정하지 않는다. 카탈로그 경계 원인은 등록미확인 가설로 낮추고 CLI만 성공한 관측사실과 구분한다. Settings→Plugins는 공식문서와 모순이므로 삭제한다.
- [ ] F03/F05 반영. Command 계약의 공식근거는 Team Bots의 Info→Setup→Plugins→Add/채팅요청과 Command 표에 한정한다. 개인계정에 같은 UI가 있다고 가정하지 않는다. 사용자 질문(실제앱 추가메뉴)은 미응답이며 필요한 UI검증을 정확히 조건부 인계한다. /workspace 지속성은 보장하지 않고 업데이트/Reset 뒤 명령/파일 검증과 재설치 복구를 안내한다.
- [ ] F06-F08 최소보완. stale package 경로에서 기존 파일이 새bundle에 섞이지 않도록 새 추출영역을 쓴다. 추출/검증 실패로 기존 준비bundle을 먼저 삭제하지 않는다. 상대 KNOWSLINK_PREFIX는 절대경로로 정규화하거나 변경전 명시적으로 거절한다. prefix/node가 기존 사용자 파일 등 비정상 형태면 임의삭제하지 않고 실패한다. installer가 소유·교체하는 경로를 문서에 명시하고 prefix의 무관파일은 보존한다. 과한 프레임워크/대규모리팩터링은 금지한다.
- [ ] 관련 회귀: stale marker가 새bundle로 유입되지 않음, 잘못된 archive/추출/검사 실패에 기존preparedfile 보존, 상대prefix 경계, 기존node파일 보존·실패, 신규/재실행·env-i tools2/statusheld. 신규실험은 임시prefix/clone, 사용자실제/workspace와 ~/.grok는 수정하지 않는다. 코어/Node20/Go/UI/전체CLI QA는 반복하지 않는다.
- [ ] 수정문서/D10/D12/D13/실행기록 및 댓글초안을 실제근거에 맞게 보완한다. 이전 리뷰와 원본 QA/완료아카이브는 소급 덮어쓰지 않는다. 기존잘못된 주장에 대한 정정은 현재 실행문서·후속로그에 남긴다. ZIP README 변경으로 새SHA256을 기록한다.

제품범위 scripts/install_bot_mcp.sh와 해당설치README/기술정본/실행기록. 리뷰파일·타역할인박스·PLANS·board는 수정하지 않는다.

## 완료·보고

새도구카탈로그/계정등록 성공은 원격재시험 전 미검증이다. 준비 산출물과 관련경계회귀·product/FullOps lint exit0을 완료하고 commit한다. work.py finish로 자기인박스/전문을 보존한 뒤 [완료] SAR-MVP-002-BOT-CATALOG-DEV-FIX | final SHA <전체Git해시> | 변경·검증·미검증·새ZIP SHA256 형식 worker_done. coor가 수정 고정SHA delta 독립리뷰와 필요한 좁은후속 QA만 배정한다.
