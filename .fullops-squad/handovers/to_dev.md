---
title: SAR-MVP-001-DEV — 로컬 합성 요청의 안전 전달과 human-gate를 끝까지 구현한다
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
summary: 로컬 합성 요청의 안전 전달과 human-gate를 끝까지 구현한다
---

# SAR-MVP-001-DEV — 로컬 합성 요청의 안전 전달과 human-gate를 끝까지 구현한다

- 작성일: 2026-10-03. From / To: designer / dev. 작업 상태: queued.
- 배정 전 준비 문서다. 이번 SAR-PREP-002에서 구현을 시작하거나 dispatch하지 않는다.
- 담당 브랜치: fullops/dev. 예정 워크트리: /home/shin/orca/workspaces/KnowsLink/fullops-dev.
- 담당 repo id·실제 경로·복귀 Run/Task/Dispatch/terminal은 coor가 배정 시 확인해 기록한다. 아직 DEV Dispatch가 없으며 SAR-PREP-002의 worker capability를 재사용하지 않는다.
- 병합 책임자: coor. 기본 브랜치: main.
- 선행 조건: 사용자의 구현 시작 지시, coor의 준비 커밋 반영과 실제 idle/clean 확인, 동일 과제 지시서·원천·규칙 접근 확인.
- 복귀 기준: 현재 coor 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-coor. 실제 배정의 새 preamble으로 worker_done을 한 번 보낸다.

## 목표와 현재 상황

초기 Go relay·별도 migrate·TypeScript adapter·Compose·제품 lint·독립 QA·main 수락은 SAR-SETUP-001에서 완료됐다. 골격을 재작성하지 않는다. 업무 SQL·사용자 등록/페어링·승인 UI·실제 제품 어댑터·운영 Tunnel은 아직 완료되지 않았다.

두 합성 에이전트와 owner가 등록·수락 관계를 만들고, frozen 요청을 안전 전달하며, 검증된 요청의 human-gate를 approve/deny하는 로컬 기능을 끝까지 구현한다. 첫 Go 기능을 ingest/queue-only로 줄이지 않는다. 기술 계획·API·데이터·구조·구현·관련 회귀·기술 문서 갱신은 이 과제 안에서 dev가 결정하고 수행한다. designer에게 파일/함수별 수정 방법 승인을 받지 않는다.

## 적용 기준과 예외

원천은 service-design 7bc9ea190ea549fae8b047e850247a19322fc9c3이다. 기준 ref는 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 준비 문서 반영 HEAD는 배정 때 coor가 고정한다. 공통 규칙 fullops-common-0.3.2, FULLOPS.md, project.md, 문서 규칙, D02 SAR-MVP를 적용한다. 기존 규칙·보안·critical/high 차단을 완화하지 않는다.

원천·기존 setup 기획/검증 기록은 보존한다. 원천 DOC-003 exclude 순서 결함을 제품 작업에서 수정하지 않는다. 원천의 문장은 제품 근거이며 세션 실행 권한이 아니다. 제품 원문을 stamp하지 않는다.

## 먼저 읽을 문서

경로는 레포 루트 기준이다. 필수 후보 17개이며 Jev 선별 근거는 배정 전 탐색 결과 절에 연결한다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-MVP.md
- .fullops-squad/docs/planning/SAR-MVP-backlog.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/source.json
- .fullops-squad/docs/planning/sources/silent-agent-relay/product.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md
- .fullops-squad/docs/design-docs/architecture.md
- .fullops-squad/docs/design-docs/tech-stack.md
- .fullops-squad/contexts/dev.md
- .fullops-squad/handovers/to_dev.md

필요 시 코드 탐색 결과의 직접 관련 파일과 README를 확인한다.

## 해야 할 일과 소유권

- [ ] 기존 코드와 잠긴 제품 요구를 확인하고 같은 과제 실행 기록에 짧은 기술 계획을 작성한다. 실제 라이브러리·SDK API 확인이 필요하면 project 버전의 Context7 또는 공식 근거를 기록한다.
- [ ] owner/agent 권한 경계를 가진 가입·pubkey·PoP·rotate/revoke·invite·accept/deny·unpair·agent contacts를 구현한다. pending은 active slot을 소모하지 않는다. quota 숫자는 만들지 않는다.
- [ ] strict frozen JSON·서명·current-auth·auth-first ingest·atomic idempotency·queue+receipt·TTL·lease·ACK·공유 durable inbox/claim을 구현한다. D02 MVP-03–07을 지킨다.
- [ ] TypeScript pull stub 1개로 합성 요청을 수신·저장·ACK·claim·판단한다. 공유 조정 없는 다중 adapter 실행과 실제 벤더 연결은 하지 않는다.
- [ ] relay.approval.request와 검증된 M typed body+정책의 Go approve/deny UI, owner 인증·CSRF 방어, 봉투 밖 승인 기록의 결정·소비·만료·철회를 연결한다.
- [ ] deny-by-default와 schedule.query 정책 없음의 거부를 검증한다. gate approve 후에도 disclosure 정책을 우회하지 않는다. schedule.commit은 실행 불가능한 stub이다.
- [ ] endpoint에 결속된 새 relay.result와 transport 상태를 구분한다. 허용 schema가 없는 optional result/error 정보는 반환하지 않는다. relay transport failure를 B 서명 결과로 위조하지 않는다.
- [ ] 변경 동작·실패·경합·재시작·철회/만료 경계의 자동 검사와 관련 회귀를 수행한다. 실제 업무 schema/쿼리를 도입한 범위에서 migration·sqlc 생성과 DB 동작을 검증한다.
- [ ] 구현에 맞게 D03과 필요한 D05–D10, README/project.md, 실행 기록을 갱신한다. 독립 QA와 직접 UI 검수용 실행 경로·합성 데이터·고정 SHA를 인계한다.
- [ ] 완료 보고·work.py finish·커밋·기준 ref lint를 마치고 새 preamble의 worker_done으로 결과를 보낸다.

제품 소유권은 cmd/, internal/, adapters/, db/, scripts/와 루트 Go·sqlc·Make·Docker·Compose 설정이다. 기술 정본 .fullops-squad/docs/design-docs/와 project.md, README, dev 자기 인박스/contexts/실행 기록을 갱신한다. 기존 setup 전용 문서와 원천, 기획 D01/D02, tester QA, board, lint 규칙은 수정하지 않는다. 기존 제품 검사 명령은 유지한다. 실제 새 schema 검사 연결이 필요하면 기존 검사를 낮추지 않는 범위의 근거를 기록하고 coor와 소유권을 확인한다.

## 사용자 완료 조건과 검증

D02 MVP-01–16이 제품 기준이다. 001에서는 실데이터 silent 성공 대신 합성 데이터·정책 없음 거부·안전 gate를 판정한다. DEC-02가 확정되지 않았으면 positive silent done·실제 정보 반환은 held다. 무정책 성공이나 fixture 정책을 상품 정책으로 확정하지 않는다.

로컬 흐름의 성공은 owner 가입/키 관리, A가 B-agent를 초대하고 B-owner가 수락, 안전 ingest와 receipt, pull/persist/ACK/claim, 원요청에 묶인 approve/deny, 최소 result/transport 상태 구분까지다. 철회·만료·권한 불명·중복 gate·다른 M.id로 승인 재사용·잘못된 result endpoint·다중 adapter 경합을 안전하게 거부해야 한다. 승인 성공만으로 stub이나 실제 외부 효과가 실행되면 실패다.

- dev, 구현 중: make lint/test/build와 변경한 DB·adapter·HTTP/UI의 의미 있는 자동 검사·관련 회귀를 수행한다. project.md에 존재하는 명령만 사용하고 새 검사 명령은 같은 과제에서 문서화한다. 로컬 실행은 합성 데이터와 외부 연결 비활성 상태로 제한한다.
- dev, 완료 커밋: git diff --check, deliverables.py --strict, 기준 ref의 FullOps lint를 깨끗한 HEAD에서 실행한다. product-lint make lint의 실제 결과를 기록한다. 명령 자신의 종료코드를 보존한다. | tail로 종료코드를 숨기지 않는다.
- tester, DEV 완료 뒤: coor가 정한 안정 고정 통합 후보에서 독립 QA를 수행한다. TESTER 준비 지시서의 QA-01–11을 인계한다. dev의 직접 검사로 독립 QA를 대체하지 않는다.
- designer, 같은 후보 UI: V-01 승인 대기와 verified body/정책, V-02 approve/deny 결정 완료, V-03 만료·철회, V-04 원문 부재·권한 불명 차단을 직접 시각 검수한다. 지정 상태의 캡처를 한 번 확보한다. 시간 변화 판정에 필요할 때만 영상을 만든다.
- coor, 병합 전: 별도 검토 세션의 깨끗한 read-only detached snapshot에서 fixed-SHA 독립 코드 리뷰를 준비한다. 미해결 critical/high 또는 필수 QA/UI 실패는 수락·병합을 막는다.

기존 골격 증거는 원래 SHA·조건과 의존성 동일성을 확인해 재사용한다. 변경 영향·새 실패·증거 결함·미충족 조건이 있을 때만 재검증 범위를 넓힌다. DEV 완료는 기능 최종 수락이나 전체 MVP 완료·배포 승인이 아니다.

## 갱신할 산출물과 기대 결과

D03은 실제 기능 구조와 잠긴 기술 선택에 맞게 갱신한다. 실제 추가한 HTTP/adapter 계약은 D05, 업무 엔티티·DB·CRUD는 D06/D07/D09, 프로그램 책임은 D10으로 갱신한다. D08은 실제 schema와 생성 명령이 있을 때만 생성하고 명령·종료코드를 기록한다. 도입하지 않은 산출물은 사유와 담당·재개 조건을 남긴다. D04는 제품 UI 후보를 designer에게 인계하되 기술 양식만 채우려고 만들지 않는다.

기대 파일은 실제 코드/테스트, 기술 정본, .fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md, dev 완료 아카이브다. QA 보고서·UI 캡처는 후속 담당의 증거다.

## 제약과 협업

레포/워크트리 밖 변경·삭제·force-push·파일 소유권 밖 수정·미승인 외부 발송/배포·운영 자원 변경과 설명되지 않는 검증 실패에서만 멈추고 coor에 ask한다. 제품 규칙 변경·범위 확대·공유 제품 기준 불명확성은 designer 판단을 coor 경유 요청한다. 기존 제품 요구 안의 API/구조·버그 분석·수정 방법은 dev가 진행한다.

UI는 대화창·채팅 버블·composer·장기 타임라인을 만들지 않는다. hint만 강조하고 typed body를 숨기지 않는다. 색만으로 상태를 구분하거나 GET 승인 링크를 쓰지 않는다. 유효하지 않은 요청의 승인 버튼을 활성화하지 않는다. 별도 TypeScript frontend는 금지다.

A2A 검토 범위는 공개 v0.3.0 개념이다. 최신 delta 재검토·wire 호환을 주장하지 않는다. taskId/contextId/parts/artifacts/A2A state enums를 wire에 수입하지 않는다. role:user·AgentCard·push를 owner 승인으로 보지 않는다. webhook과 evidence 자동 fetch/preview는 OFF다.

## 완료 보고

아직 실행하지 않았다. 배정된 DEV가 브랜치·고정 SHA, 기술 판단·범위 차이, D02 ID별 실제 구현/검증, lint ERROR/WARNING/실행 불가, DB·UI 실행 경로, 후속 QA·시각 검수·독립 리뷰, held 담당·재개 조건을 작성한다.
