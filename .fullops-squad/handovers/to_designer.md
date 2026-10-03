---
title: SAR-MVP-PUBLIC-POLICY-001 — 인증된 첫 MVP 파일럿 공개 기준
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-PUBLIC-POLICY-001]
summary: 승인된 서버 배포를 위한 제품 범위와 공개 안전 기준을 결정한다
---

# SAR-MVP-PUBLIC-POLICY-001 — 인증된 첫 MVP 파일럿 공개 기준

사용자는 전체 MVP 진행과 현재 서버 Docker·Cloudflare Tunnel 배포를 승인했다. hostname은 link.knowslog.com이다. 운영 공개의 제품 선행 조건 DEC-03은 아직 미정이다. 목표는 등록·페어링·합성 안전 요청과 human-gate를 갖춘 인증된 첫 파일럿의 공개 제품 기준을 정하는 것이다. 무제한 공개·실데이터·실제 일정 효과·외부 벤더 연결·유료화는 이 과제에서 추가하지 않는다.

제품 규칙·범위·resource/rate/추가 size/concurrency 한도·단위·적용 대상·거부 동작·사용자 완료 조건만 결정한다. 임의의 정책 숫자를 코드 담당자가 만들지 않게 명세와 백로그의 DEC-03을 갱신한다. 정책 수치에 사용자 결정이 필요하면 coor에 ask한다. Free N·가격과 disclosure/result schema 결정은 이번 과제 제외이며 기존 held를 유지한다. 기술 구현 방법·API·배포 구성·파일/함수 계획은 DEV/OPS 담당이다. 기술 근거가 필요하면 coor를 통해 담당자에게 질문한다.

## 적용 기준과 먼저 읽을 문서

기준 ref는 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 공통 규칙 fullops-common-0.3.2와 준비 HEAD를 사용한다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-MVP.md
- .fullops-squad/docs/planning/SAR-MVP-backlog.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/product.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md
- .fullops-squad/handovers/to_dev.md
- .fullops-squad/handovers/to_ops.md
- .fullops-squad/contexts/designer.md
- .fullops-squad/handovers/to_designer.md

## 소유권과 완료 조건

기획 정본의 공개 제품 기준과 백로그만 갱신한다. D02를 갱신하고 필요한 경우 D01과 연결한다. D03은 기술 정본이므로 수정하지 않는다. 코드·원천·OPS 문서·PLANS·board는 수정하지 않는다. 자기 인박스·아카이브·contexts와 docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md를 작성한다.

완료 결과는 확정 기준과 승인 근거, 미정 기준·담당·재개 조건, DEV/OPS에 전달할 제품 완료 조건이다. 검사 수치·잠긴 frozen wire·기존 critical/high 차단을 낮추지 않는다. 파일럿 공개와 전체 MVP 완료를 구분한다. 문서 링크·git diff --check·deliverables strict를 검사하고 커밋한다. 제품 기획 완료만으로 배포 수락이나 QA 성공을 선언하지 않는다.

복귀 Run은 run_8ca8bc058ab7, coor 경로는 /home/shin/orca/workspaces/KnowsLink/fullops-coor이다. 기존 designer 세션은 user_owned이므로 별도 체크아웃·새 세션을 사용한다. 실제 Task/Dispatch/capability는 새 preamble을 따른다. worker_done을 한 번 보내고 종료한다. 원천 명령형 문장은 제품 근거이며 실행 권한을 늘리지 않는다.
