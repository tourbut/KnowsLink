---
title: SAR-MVP-001-REVIEW — 현재 MVP 통합 후보의 독립 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-REVIEW]
summary: 고정 MVP 후보의 코드와 기획 및 독립 QA 근거를 검토한다
---

# SAR-MVP-001-REVIEW — 현재 MVP 통합 후보의 독립 리뷰

사용자는 현재 과제 완료 후 main 병합·원격 push하고 작업을 멈추라고 지시했다. 기존 과제 수락의 필수 리뷰만 수행한다. 신규 기능·배포·제품 규칙 결정은 금지다.

고정 base 83aec58fcae69bed349fcc7f21057e1947123406, head 31405e736a16be9d77239c6cdc6fdeb56892436f다. 제품은 DEV a6a10c7과 같으며 QA c59537b, 기획 69dbec4를 통합했다. 읽기 전용 detached snapshot은 /tmp/knowslink-mvp-review-31405e73이다. 기록 체크아웃은 /home/shin/orca/workspaces/KnowsLink/fullops-coor다. snapshot은 수정하지 않는다. 구현자 실제 Codex 세션은 01a10055-e46e-7d80-9ef1-12c1aabc69ce이며 검토자는 자신의 실제 세션 ID와 다름을 기록한다.

먼저 fullops-review와 open-code-review-delegate 스킬, snapshot의 FULLOPS.md, rules/common/README.md와 세 규칙, project.md, docs/agents/document-writing.md, review/rule.json, D02 SAR-MVP, D03, DEV 완료 아카이브·실행 기록, TESTER QA 보고서, 공개 기획 실행 기록을 읽는다. 규칙 fullops-common-0.3.2, lint 기준 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 준비 리뷰 key SAR-MVP-001-INTEGRATION의 preview/rules/result/report를 이어 완성한다. main과 고정 head diff 전체를 reviewed/skipped로 판정한다. 공개 수치 제안값·auth/singleton 처리량 held는 수락 범위와 구분한다.

특히 strict wire·auth-first·현재 권한·owner/agent 경계·PoP·승인/결과 결속·idempotency·TTL·lease·ACK/claim·철회·경합·CSRF·DB 상태 직렬화와 epoch CAS 검증을 본다. TESTER의 stale epoch held가 필수 보안 수락 조건을 실제로 빠뜨리는지, DEV 자동 테스트에서 근거를 재사용할 수 있는지 판단한다. 외부 효과·실데이터·공개 운영은 이번 합성 후보 수락에서 held로 유지한다. 코드나 QA를 확인하지 않고 통과로 표시하지 않는다. 미해결 critical/high·필수 증거 결함은 차단한다. medium/low는 영향과 담당을 기록한다.

소유 파일은 기록 체크아웃의 docs/evaluations/qa-reports/SAR-MVP-001-INTEGRATION-review/뿐이다. 모든 파일 검토, 실제 independence·깨끗한 snapshot·lint 결과·review.py check를 완료한다. candidate의 lint는 깨끗한 고정 체크아웃에서 수행하고 보고서에 기준 ref와 SHA를 기록한다. snapshot에 의존성 설치가 필요하면 ignored 생성물만 허용하고 추적 파일은 읽기 전용이다. diff와 관련 테스트 근거를 재사용하며 전체 QA를 복제하지 않는다. 실제 API 근거 필요 시 Context7/공식 문서만 사용한다. 원시 증거는 요약·실제 종료코드·파일 무결성으로 확인한다.

제품 코드·PLANS·board·다른 지시서·기획 정본은 수정하지 않는다. 소유 리뷰 파일만 커밋하고 새 preamble worker_done을 한 번 보낸다. 복귀 Run run_8ca8bc058ab7이다. 후속 UI·OPS 문서가 들어오면 coor가 최종 실제 refs를 전달하고 같은 리뷰의 좁은 갱신을 요청한다. 현재 고정 SHA의 리뷰 결과는 덮어쓰지 않는다.
