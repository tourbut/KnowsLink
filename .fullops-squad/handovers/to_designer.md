---
title: SAR-MVP-001-UI — 현재 후보 직접 UI 검수
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-001-UI]
summary: 기존 MVP 승인 화면의 V-01–04 시각 수락 조건을 확인한다
---

# SAR-MVP-001-UI — 현재 후보 직접 UI 검수

사용자는 현재 작업의 완료·검증·병합·원격 공유 후 작업을 멈추라고 지시했다. 이번 검수는 기존 MVP 완료에 필요한 필수 검수이며 신규 기능·배포·정책 확정은 하지 않는다.

제품 후보 a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92와 tester QA c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8의 V-01–04 캡처를 직접 본다. .fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/의 pending·approved·denied·expired·revoked·unavailable·unauthorized PNG를 view_image 등 실제 시각 도구로 확인한다. HTML이나 tester 판정만으로 시각 검수 완료를 선언하지 않는다. 읽을 PNG는 기존 tester 체크아웃의 같은 경로에도 있다.

먼저 FULLOPS.md, rules/common/README.md와 coding-style.md/testing.md/security.md, project.md, docs/agents/document-writing.md, docs/planning/product-specs/SAR-MVP.md, docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md, contexts/designer.md를 읽는다. 규칙 fullops-common-0.3.2, lint 기준 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 기존 DEV/TESTER 탐색 근거와 관련 문서를 재사용한다. 후보에서 디자인 기준을 낮추지 않는다.

V-01 verified typed body와 정책·발신/대상·pending 승인/거절, V-02 결정 결과와 중복 승인 불가, V-03 expired/revoked와 승인 불가, V-04 원문 부재/권한 불명 차단을 판정한다. hint 강조로 typed body가 숨거나 채팅 composer·버블·장기 timeline·색만 상태·GET 승인 링크가 있으면 실패다. unauthorized 캡처가 단순 invalid_auth라는 점은 기존 완료 조건과 대조한다.

소유 파일은 docs/design-docs/mockups/의 D04 원천, docs/evaluations/qa-reports/SAR-MVP-001-UI.md, docs/exec-plans/phases/SAR-MVP-001-UI.md, 자기 contexts와 인박스·아카이브다. 코드·원천·PLANS·board·기존 QA 증거는 수정하지 않는다. 기존 PNG를 재사용하고 새 캡처·영상은 결함 판정에 필요한 경우만 만든다. 검사 결과 PASS/FAIL/held와 이미지 경로·고정 SHA를 기록한다. strict·git diff --check·필요 lint를 수행하고 소유 파일만 커밋한다. work.py finish와 새 preamble worker_done을 한 번 보내고 종료한다. 코드 결함은 coor에 보고한다.

기존 designer 워크트리를 사용한다. 복귀 Run run_8ca8bc058ab7, coordinator 경로 /home/shin/orca/workspaces/KnowsLink/fullops-coor다. 지시서의 진행 중 작업은 이번 검수까지이며 후속 기능이나 배포를 시작하지 않는다.

## 실패한 검수 기록 마무리 재개 — 2026-10-03

이전 Dispatch ctx_b7d073da41ef는 failed이며 execution-host worker-show observation은 exited다. 기존 UI 보고서·D04·보완 PNG·실행 기록을 보존해 같은 과제를 마무리한다. 설치된 FullOps 0.9.12를 사용한다. 기준 ref는 준비 HEAD 0a5b044다. 기존 직접 열람은 원래 관찰로 유지하고 이번 세션에서 필요한 이미지만 직접 확인한다. 문서 front matter·상대 링크·strict·소유파일 커밋·work.py finish·새 worker_done을 완료한다. 기존 제품 후보 a6a10c7과 QA c59537b의 시각 검수이며 신규 코드 실행이 아니다. reviewer msg_4fbcac80f76c의 deliver:human 인증 경계 high는 DEV 수정·독립 재검증 대상이다. 이를 UI PASS로 해소하지 않는다. DEV 수정 후보가 UI에 영향을 주면 coor가 같은 과제 후속을 전달한다. 제품/기획/PLANS/board는 수정하지 않는다.
