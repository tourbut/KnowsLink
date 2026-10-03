---
title: SAR-MVP-002-INSTALL-FIX-DEV-REVIEW — 설치 수정 후보 5506d64의 독립 코드·문서 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-INSTALL-FIX-DEV-REVIEW]
summary: 설치 수정 후보 5506d64의 독립 코드·문서 리뷰
---

# SAR-MVP-002-INSTALL-FIX-DEV-REVIEW — 설치 수정 고정 SHA 독립 코드·문서 리뷰

- 작성일: 2026-10-04; From: coor; 상태: ready
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble.
- 병합 담당 coor, main/origin. 사용자 승인: 이슈1 진단·수정·검증·푸시·완료 후 댓글 게시. worker는 외부 댓글을 쓰지 않는다.
- 실제 Bot 계정/유료 inference/relay/DEC-02/calendar·실데이터는 held, FullOps 업데이트 제외.

## 적용 기준과 먼저 읽기

FULLOPS.md, project.md, rules/common/README.md 및 coding-style.md/testing.md/security.md, docs/agents/document-writing.md, docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md, 해당 DEV handovers/logs/2026-10-04_to_dev.md, adapters/README.md, scripts/package_plugin.py, scripts/verify_grok_plugin.py. 공통 0.3.2. 현재 준비 SHA에서 규칙을 읽고 제품 대상은 아래 고정 SHA를 쓴다. 이전 QA0fb32cd·리뷰0055a5b는 코어/Go/UI 변경 없는 부분의 증거로 재사용 가능하다. 설치 및 패키지 변경은 새로 확인한다.

- 역할/워크트리: ops, /home/shin/orca/workspaces/KnowsLink/fullops-ops, fullops/ops
- 고정 후보: 5506d646c47c2d64b35d1254ddfdec5e2003084d, 기준 376981441d1b0650897847273e1773f8473949e5.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-INSTALL-FIX-DEV-REVIEW-*.json. keep 전부 읽는다. 지시 전제와 충돌 — 먼저 확인: 현재 DEV 실행 기록과 이전 QA(설치 미실행); current DEV 기록의 주장 자체를 독립 검증하며 이전 QA는 불변 코어에만 재사용한다. 근거가 틀리면 결함으로 보고하고 제품 규칙 충돌이면 ask한다.

## 해야 할 일과 소유권

- [ ] fullops-review와 open-code-review-delegate를 적용한다. OCR 준비 키는 SAR-MVP-002-INSTALL-FIX-DEV, 경로 docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review이다. 구현 세션08407653-b74b-4e95-abc2-4074c29956cb와 다른 실제 reviewer 세션을 기록한다.
- [ ] 읽기 전용 /tmp/knowslink-install-review-5506d64의 detached·clean·HEAD를 확인하고 모든 변경/관련 계약을 검토한다. snapshot을 수정하거나 npm 설치하지 않는다. 결과는 OPS 기록 체크아웃의 준비 리뷰 디렉터리에 쓴다.
- [ ] Grok plugin 지원 형식·MCP discovery·GROK_PLUGIN_ROOT·trust·marketplace·doctor 근거와 README 설치 명령을 확인한다. 런타임 strict와 checksum·설치 폴더·검증 스크립트의 격리/실패 전파도 확인한다. 봇 앱 지원을 CLI 성공으로 주장하지 않는지 확인한다. findings 위치는 파일별 실제 줄로 확인한다.
- [ ] result.json 모든 파일 판정·independence·conclusion과 보고서를 채운다. DEV 깨끗한 고정5506d64에서 lint --from3769814를 실행해 리뷰 lint.json에 저장한다. 리뷰 check --from3769814 --to5506d64 --task-key SAR-MVP-002-INSTALL-FIX-DEV exit0을 확인한다.

## 완료 기준과 결과

미해결 critical/high는 수락 불가다. 검증 증거와 생략 영향을 판단하고 실패/미검증을 숨기지 않는다. source 코드·타 역할 인박스·PLANS·board는 수정하지 않는다. 기술 결함이면 재현과 수락 조건을 coor에게 보고한다. 보고서/result/lint/check/완료 로그만 소유한다. D01–D13 갱신 없음. work.py finish로 정규 인박스와 전문을 보존·비우고 커밋한 SHA를 worker_done으로 보낸다.
