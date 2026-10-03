---
title: SAR-MVP-001-REVIEW-FIX — 수정 후보의 별도 세션 고정 SHA 독립 코드 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-REVIEW-FIX]
summary: 수정 후보의 별도 세션 고정 SHA 독립 코드 리뷰
---

# SAR-MVP-001-REVIEW-FIX — 기존 MVP 수정 후보의 독립 검증

복귀 Run run_8ca8bc058ab7, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor다. 새 preamble의 실제 Task/Dispatch/capability만 사용한다. 신규 기능·정책 수치 확정·실제 배포는 시작하지 않는다.

기존 합성 MVP의 수정 후보 전체를 별도 세션에서 독립 코드 리뷰한다. 구현은 하지 않는다. 고정 base ffca87c9ec501e9f313c50f60124b686d27af728, head 4262d02fdd7b0b57a804d9e550597852950ffeae다. 제품 코드는 d7e2149와 같고 기존 QA·기획·직접 UI 결과를 포함한다. 읽기 전용 깨끗한 detached snapshot은 /tmp/knowslink-mvp-review-4262d02이다. 기록 체크아웃은 상설 fullops-dev다. 설치된 review.py prepare로 새 key SAR-MVP-001-REVIEW-FIX를 사용한다. 원래 중단 리뷰는 성공으로 바꾸거나 덮어쓰지 않는다.

실제 구현 세션은 기존 Codex 01a10055-e46e-7d80-9ef1-12c1aabc69ce와 C1 수정 Claude 4b9e5e7e-4fa5-442a-8a7a-d9c1d3c51c33이다. 후자는 실제 /home/shin/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-dev 메타데이터로 확인했다. independence에는 실제 수정 구현 세션과 자신의 서로 다른 실제 session ID·snapshot 경로/head/read_only를 기록하고 원래 구현자도 보고서에 남긴다.

fullops-review·open-code-review-delegate 스킬과 snapshot의 FULLOPS·공통 README와 세 규칙·project·문서 작성 규칙·review/rule.json·D02/protocol C1–C5·D03/D05–D10·DEV 완료 원문/수정 기록·TESTER QA·공개 기획/직접 UI 기록을 읽는다. 공통 fullops-common-0.3.2다. 기존 탐색 근거를 재사용한다. 전체 main 대비 제품 diff와 문서·검증 조건을 reviewed/skipped로 검토한다. 새로운 후보의 변화 없는 동작 증거는 원래 실행 SHA·조건으로 연결하고 신규 실행으로 주장하지 않는다.

특히 reviewer msg_4fbcac80f76c의 C1 high가 실제 해소됐는지 확인한다. agent/human delivery route 저장과 lease/persist/ACK/claim 인증 경계, H 이외 direct deliver:human의 403 제한이 frozen 정본 요구를 낮추거나 새로운 제품 결정을 만들지 않는지 판정한다. strict wire·현재 권한·PoP·owner/agent 경계·approval/result 결속·idempotency·TTL·철회·lease/claim·CSRF·DB 경합/epoch CAS도 확인한다. QA-06 stale epoch의 기존 targeted 검사·DEV 자동 회귀 재사용 가능성을 판정한다. 공개 기획의 제안값은 확정하지 않으며 별도 기획 문서 결과를 검토·통합할 수 있는지 보고한다. 전체 제품 공개·실벤더·신원·정책·singleton 처리량·WAL/backup held는 합성 수락과 분리한다. 독립 QA 후속은 coor가 전달한다.

소유는 새 docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/·현재 review inbox/logs다. 제품·기획/기술 정본·PLANS·board·원래 리뷰 결과는 수정하지 않는다. snapshot 추적 파일은 읽기 전용이며 ignored 의존성 설치만 허용한다. snapshot lint를 위 base 기준으로 실행하고 실제 head의 ERROR 0·product-lint passed를 확인한다. 실제 review.py check를 고정 base/head로 실행한다. 미해결 critical/high·필수 증거 결함이면 수락을 차단한다. 검토 결론·발견 사항·재현·검증 범위·skipped 영향·independence를 완성한다. 원래 report.md는 중단 템플릿이며 result pending을 과거 완료로 주장하지 않는다. 새 리뷰가 최종 수락을 판정한다.

완료 전문에 기록 경로와 검토 head를 연결하고 work.py finish로 현재 review inbox를 비운다. 소유 파일만 커밋한다. 새 worker_done 본문에 `[완료] SAR-MVP-001-REVIEW-FIX | SHA <보고 커밋> | 리뷰 head 4262d02 | ...`를 포함하고 outcome을 명시한다. 구현자 새 수정이 필요하면 coor에 보고한다. main 병합·push는 coor가 필수 QA·시각 재사용·리뷰를 확인해 수행한다.
