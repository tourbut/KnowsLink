---
title: SAR-PUBLIC-MESSAGES-001-DEV — 일반 회원의 연결 확인 text·관련 회신·receipt·실패 복구·gate와 자원 보호 구현
status: draft
updated: 2026-10-06
owner: coor
tasks: [SAR-PUBLIC-MESSAGES-001-DEV]
summary: 일반 회원 메시지·receipt·gate 구현과 자동 검사·독립 검수 인계
---

# SAR-PUBLIC-MESSAGES-001-DEV — 일반 회원 메시지·receipt·gate 구현

- 작성일: 2026-10-06
- From / To: coor / dev
- 상태: ready
- 담당: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb / /home/shin/orca/workspaces/KnowsLink/fullops-dev / fullops/dev
- 병합: coor / main·origin/main
- 복귀: 같은 repo / /home/shin/orca/workspaces/KnowsLink/fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. Task·Dispatch는 worker-start preamble 기준.

## 현재 상황과 확인 근거

사용자가 다음 메시지 과제 진행을 요청했다. IDENTITY와 AGENTS 로컬 코드·독립 QA·보안 리뷰·직접 UI는 main 수락됐다. 기준 SHA `7efbaa349a5857eb1eac859a955ec3a09c91f800`은 main/origin/main 및 5역할의 clean HEAD다. 실제 이메일·운영 공개·실24h·최종 노우↔다닷은 미검증이다. 기존 MESSAGES 대기는 이번 사용자 지시로 재개한다. 제품 수치·규칙은 이미 designer 정본에 있다.

## 적용 기준과 예외

fullops-common-0.3.3의 README·coding-style·testing·security, FULLOPS.md, project.md, document-writing.md를 같은 기준 SHA에서 적용한다. 제품 정본은 SAR-PUBLIC-SERVICE PS-08–11과 UX-06/07이며 PS-04/06/07 권한·세대·철회 및 frozen SAR-MVP C1–C5도 유지한다. 예외로 수락 수준을 낮추지 않는다. 기술 분석·짧은 계획·구현·검사·기술 문서 갱신은 같은 DEV 과제에서 결정한다. 제품 기준 변경만 coor 경유 designer에게 질문한다.

## 먼저 읽을 문서

- .fullops-squad/FULLOPS.md, rules/common/README.md와 연결된 세 규칙, project.md, contexts/dev.md, docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md, SAR-MVP.md
- .fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md
- .fullops-squad/docs/design-docs/architecture.md, tech-stack.md
- .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md, SAR-PUBLIC-AGENTS-001-DEV-FIX.md
- adapters/README.md

먼저 확인할 기술 전제: connections.go는 기존 연결 전용 경계이며 메시지 전달 권한을 자동 부여하지 않는다. AGENTS 실행 기록의 text/receipt/gate 후속 미구현은 이번 구현의 선행 상황이다. 두 파일의 전제를 직접 확인한 뒤 기존 PS04/06/07와 frozen 계약을 유지한다. Jev의 충돌 가능성 표시는 이 범위 차이로 해석하며 새 제품 결정을 만들지 않는다.

먼저 확인할 코드: internal/relay/connections.go·protocol.go·test_messages.go·member.go·store.go, adapters/src/core.ts·connect.ts.

Jev code/doc find·context 결과는 docs/evaluations/jev/SAR-PUBLIC-MESSAGES-001-DEV-*.json에 남긴다. 기존 문서의 후속/미검증 설명은 이번 기능의 구현 전제이며 과거 PASS로 해석하지 않는다.

## 해야 할 일과 파일 소유권

- [x] 기존 일반 연결·text/relay/MCP/Node CLI·receipt/gate 흐름과 실제 지원 인터페이스 근거를 확인하고 짧은 기술 계획을 실행 기록에 쓴다. 실제 클라이언트 호환성을 문서나 fixture만으로 주장하지 않는다.
- [x] PS-08/09의 비민감 연결 확인 text(4096 UTF-8 bytes·TTL180s)·관련 답장·요청/답장 ID·현재 수신/처리/실패 상태와 수동 pull 안내를 일반 회원 경로에 구현한다. frozen relay.v1 업무 wire와 분리하며 자기 agent·현재 활성 관계·키·세대를 매 단계 확인한다.
- [x] PS-10 회원 human gate의 검증된 typed body·정책·기한·현재 권한 approve/deny와 PS-11의 메시지/gate/receipt/HTTP·claim 자원 상한 및 안전 정리 budget을 구현한다. 기존 신원/AGENTS 한도는 재사용한다.
- [x] 변경 동작·오류·경계·동시성·재시작·high priority·타회원 거부·철회/재수락·idempotency 충돌·XSS/신뢰하지 않는 text·안전 정리 포화 회귀를 자동 검사한다. 자기 격리 fixture/자원만 사용하고 회수한다.
- [ ] 영향 기술 산출물·실행 기록·PLANS·context를 갱신하고 최종 커밋·일반 역할push·고정HEAD lint/test·strict 검사 후 완료 전문과 worker_done을 보낸다.

DEV 소유: cmd/·internal/·adapters/·db/·scripts/·제품 설정과 기술 docs/design-docs(기획 mockups 제외)·QA 자동 검사·docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md 및 자기 context/PLANS/인박스·아카이브. 제품 기획/UX 정본과 타 역할 인박스는 수정하지 않는다. 예상 범위는 relay·adapter·회원 UI·관련 검사/기술 문서이며 실제 규모와 SIZE 경고 처리 근거를 완료 보고에 쓴다. 추측한 모듈/기능은 추가하지 않는다.

## 완료 기준과 검증

DEV는 승인된 로컬 코드 범위를 끝까지 구현한다. make lint/test 및 영향 SQL·adapter·UI 통합 make verify-mvp, 필요 변경 경로 검사·git diff --check·deliverables strict를 수행한다. 커밋 뒤 최종 고정 HEAD에서 FullOps lint --from 기준SHA를 실행하고 HEAD·각 명령 exit·ERROR/WARNING/실행불가를 기록한다. 완료 보고 마지막 metadata/archive 커밋도 고정HEAD 증거를 갖춘다. 기존 실패를 지우거나 재실행 PASS로 바꾸지 않는다.

coor는 고정 후보에서 구현자와 다른 별도 OPS 세션의 독립 보안 리뷰, tester의 PS08–11 독립 QA, designer의 UX06/07 직접 UI 검수를 배정한다. 필수 실패와 미해결 critical/high는 main 수락을 차단한다. 전체 운영 PS13·최종 PS14·실제 메일/외부 플랫폼 왕복은 로컬 코드 완료와 분리한다. 실제 지원 클라이언트 코드/로컬 실제 프로세스 왕복 근거는 확보하고 운영/사용자 관측이 필요하면 정확한 담당·재개조건을 남긴다. 이미 검증된 변경 없는 증거는 원래 SHA/조건으로 연결하고 영향 범위만 재검증한다.

### UI 작업의 디자인 기준과 검증

기존 Go template/CSS·회원 공용 화면 패턴과 UX06/07를 재사용한다. queued를 실제 수신으로 표시하지 않는다. 요청ID·관련답장ID·TTL·수동 pull/실패 복구 동작을 보여준다. gate의 hint는 검증본문을 대체하지 않는다. 채팅 버블·자유 composer·장기 타임라인·무조건 성공 배너를 추가하지 않는다. project.md의 별도 design lint/테마 미구성을 확인하고 make lint 및 해당 없음 근거를 기록한다. 직접 시각 판정은 designer 후속이며 DEV 자동검사를 시각 PASS로 보고하지 않는다.

## 갱신할 산출물

route 추천 D05·D09·D10과 실제 영향 D03·D06·D07·D08을 확인해 해당 원천을 갱신한다. 변경 없는 ID는 근거를 쓴다. D01/D02/D04 제품 정본은 designer 소유다. front matter는 deliverables.py --stamp로 작성한다.

## 제약·협업·후속

Workers Free만 허용한다. 유료 플랜·구독·초과 과금은 금지한다. 기존 서버/Tunnel 유지. 운영 배포·실메일 발송·실제 외부 계정 연결·운영 자료 삭제는 이번 구현 범위가 아니다. 로컬 비파괴 구현·격리 검사·커밋·일반 push는 승인됐다. low L-A/L-B와 공개전 합성가입unset·운영DB 합성owner0 조건은 보존한다. 공개 수락이나 실24h/실메일/노우↔다닷 성공을 선언하지 않는다. 제품 모호성·범위 밖 행위는 질문하며 기술 버그는 같은 과제에서 해결한다.

## 완료 보고

브랜치/full SHA·변경 이유·기술 판단·실제 클라이언트 근거·검증한 것/미검증·산출물·명령 exit·원본 실패·SLOP/DEP/DESIGN/SIZE 처리·후속 담당/재개조건을 전문으로 남긴다. work.py finish로 로그 보존·인박스 비움·최종 커밋 뒤 동일 Run으로 worker_done을 보낸다.

## DEV 기술 계획

기준과 분석·짧은 계획은 `docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md`에 기록했다. 별도 knowslink.text.v1·기존 shared Message·회원 receipt·공통 용량·Node CLI/MCP를 같은 과제로 구현한다. 제품 정본·frozen relay.v1은 유지한다.
