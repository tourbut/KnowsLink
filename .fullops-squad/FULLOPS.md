---
title: FullOps Squad 하네스 지도
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-098, FULLOPS-UPDATE-099, FULLOPS-UPDATE-0.9.10]
summary: 작업 규약과 역할별 참조 및 업데이트 절차를 안내한다
---

# FullOps Squad — 하네스 지도

현재 레포 루트의 `.fullops-squad/fullops.json`이 있을 때만 이 규약을 적용한다.
규약은 한국어로 관리하고 식별자·코드베이스의 기존 언어 규칙은 유지한다.
문서를 작성하거나 수정하기 전에 [문서 작성 규칙](docs/agents/document-writing.md)을 읽는다.
작성자는 front matter와 한국어 STE 작성 원칙을 적용한다. 검토자는 같은 규칙으로 문서를 확인한다.

| 작업 | 먼저 읽을 문서 (`.fullops-squad/` 기준) |
|---|---|
| 공통 코딩·테스트·보안 기준과 우선순위 | [rules/common/README.md](rules/common/README.md) 및 연결된 세 규칙 |
| 기술 기준·검증 명령·문서 정본 | `project.md` |
| 역할·워크트리·에이전트 배정 | `orca-agents.md` |
| 현재 할 일·우선순위 | `PLANS.md` |
| 프로젝트 단계·진행 현황판 | `board/board.json` (coordinator가 관리) · `board/index.html` |
| 역할별 작업 명세와 완료 처리 | `handovers/_TEMPLATE.md` · `fullops-work` 스킬 |
| Orca 착수 확인·회신·검토 | `orca-agents.md` · `fullops-orca` 스킬 |
| 코드 변경의 lint 게이트 | `lint/README.md` · `lint/lint.json` |
| 병합 전 delegate 코드·문서 검토 | `review/rule.json` · `fullops-review` 스킬 |
| 플러그인 설치·레포 적용 갱신 | `update-fullops` 스킬 · `docs/exec-plans/phases/FULLOPS-UPDATE-<버전>.md` |
| 과거 결정·교훈 | `contexts/<role>.md` |
| 산출물과 원천 매핑 | `docs/deliverables/README.md` |
| 외부 스킬 설정 | `docs/agents/{issue-tracker,triage-labels,domain}.md` |

작업의 흐름은 요구사항 → 설계 → 역할별 핸드오버 → worker 구현·검증 → 직접 회신 → 검토·병합이다.
문서에 적힌 코드·API·브랜치·도구의 존재는 작업 전에 직접 확인한다.
코드 탐색은 지시서의 `먼저 읽을 문서`에서 시작하고, 필요한 범위만 검색해 읽는다. 레포 전체를 훑어 읽지 않는다.
라이브러리·SDK·프레임워크의 API나 설정을 확인할 때는 프로젝트 버전을 기준으로 Context7의 `resolve-library-id` → `query-docs`를 사용한다. 확인한 library ID·버전·출처·필요한 API 근거를 지시서에 남기고 같은 근거는 재사용한다. 버전 변경·근거 부족 시 다시 조회하며, 미지원·연결 실패 시 공식 문서로 확인한다. 비공개 코드·비밀값을 조회문에 포함하지 않는다.
실질적 개발 작업은 지시서에 목적·범위·완료 기준·산출물·복귀 주소를 남긴다.
이미 허가된 작업은 다시 승인받지 않는다. 미승인 외부 발송·배포 등은 실행 직전에 별도로 확인한다.
작업자와 검토자는 `rules/common/README.md` 및 연결된 세 규칙을 읽고 `project.md`의 정본과 함께 적용한다. 지시서에 규칙 식별자·프로젝트 문서 경로·기준 커밋 또는 스냅샷·예외를 기록하고 worker와 검토자가 같은 버전을 읽는지 확인한다.
프로젝트 보안·아키텍처·테스트 기준이 외부 스킬보다 우선하며, 충돌은 지시서에 기록한다. 공통 기본값으로 기존 기준을 낮추지 않고 보안·권한·미해결 critical/high 차단을 임의 완화하지 않는다.
ponytail full을 적용하되 검증·산출물·핸드오버 기록은 생략하지 않는다.
하네스 플러그인 캐시는 읽기 전용 자원으로 취급한다. 개발 산출물은 이 레포에 Git으로 관리한다.

## 제품 기획과 기술 계획의 책임

제품 기획 역할은 designer다. 기술 계획 역할은 dev다. coordinator는 coor다.
designer는 제품 목표·규칙·수치·방향·우선순위·사용자 완료 조건을 결정한다.
dev는 기존 제품 요구 안의 기술 계획·구조/API·버그 분석·구현·테스트·기술 문서 갱신을 같은 과제에서 수행한다.
coor는 요청 분류·인계·질문 전달·진행 관리·병합을 조정한다. 제품 판단이나 기술 판단을 대신하지 않는다.
신규 개발 요청은 `jev_route.py`의 product/implementation/unresolved 책임 분류를 적용한다. setup 갱신·동기화·현황판 정리는 직접 처리한다.
독립 코드 리뷰는 구현자와 다른 검토자의 별도 세션과 고정 SHA의 깨끗한 detached snapshot을 사용한다.
snapshot은 읽기 전용으로 유지한다. 리뷰 결과는 별도 기록 체크아웃에 작성한다.
tester의 독립 동작 QA와 필요한 직접 시각 검수를 유지한다. 미해결 critical/high는 수락·병합을 차단한다.
전환 전 지시서·원천·리뷰·검증 기록은 보존한다. 상세 책임과 인계는 `orca-agents.md`를 따른다.

## 완료 결과의 통합과 역할 인박스 — FullOps 0.9.12

실행 지시서는 역할별 `handovers/to_<role>.md` 하나로 고정한다. 다른 과제로 사용 중인 인박스를 덮어쓰지 않는다. 다음 과제는 PLANS.md에 대기시킨다. 과제명 파일·pending·logs는 참조 자료이며 배정 지시서로 사용하지 않는다. 역할 작업 완료 시 `work.py finish`로 지시서와 완료 보고 전문을 logs에 보존하고 인박스를 비운다.

coor는 worker_done을 받으면 고정 SHA의 필수 검토·검증을 확인하고 실제 기본 브랜치 main에 병합한다. 원격 origin/main에 일반 push한 뒤 완료 SHA의 로컬·원격 조상 관계를 확인한다. 역할 브랜치 push만으로 통합을 완료하지 않는다. 전체 제품 수락과 개별 문서 결과의 통합은 구분한다. 미해결 critical/high와 필수 검증 실패는 계속 차단한다.

병합 뒤 coor를 포함한 모든 등록 역할의 실제 worker 상태와 작업 트리를 확인한다. 쉬고 있으며 깨끗한 역할만 최신 main으로 동기화한다. 진행 중·미커밋 변경·상태 확인 불가인 역할은 PLANS.md에 최신 기본 SHA와 예약을 남긴다. 다음 dispatch 전에 최신 로컬·원격 main 포함을 확인한다. 진행 중 체크아웃과 사용자 자료를 덮어쓰지 않는다.

완료 메시지는 Git 공용 디렉터리의 fullops-integration에 보존된다. `integration.py --repo <레포> status`로 미통합 결과를 확인한다. 검수 대기·실패·충돌·원격 오류·사용자 제한은 PLANS.md에 메시지 ID·SHA·사유·담당·재개 조건을 기록한다. 같은 내용을 `integration.py hold`로 남긴다. 조건 충족 시 resume하고 병합·push를 이어간다. 성공 결과를 통합한 뒤 release·ack하고 다음 독립 과제를 배정한다. 사용자의 현재 과제 완료 뒤 중지 지시는 유지한다.

기존 활성 과제명 지시서는 작업 중 이동하지 않는다. 해당 과제 완료 뒤 logs에 보존하고 다음 과제부터 정규 역할 인박스를 사용한다.
