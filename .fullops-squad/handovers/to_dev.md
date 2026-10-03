---
title: SAR-MVP-002-DEV — Grok Bot 지원 인터페이스를 확인하고 안전한 실제 어댑터 연결을 준비한다
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-002-DEV]
summary: Grok Bot 지원 인터페이스를 확인하고 안전한 실제 어댑터 연결을 준비한다
---


# SAR-MVP-002-DEV — Grok Bot 인터페이스 확인과 안전 어댑터 준비

- 상태: ready
- From / To: coor / dev
- 사용자 근거: 베타 완료 뒤 다음 작업 진행 요청. 백로그 002 우선순위를 적용한다.
- 담당: /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev
- 기준: main e732fedb8a7f80b9813219bf2dbc65fc029ff272. 배정 준비 커밋도 반영한다.
- 복귀: Run run_8ca8bc058ab7, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor. 실제 task/dispatch/터미널은 Orca preamble을 따른다.

## 현재 상황과 확인 근거

001 제품 78b1d92 수락, 폐쇄 합성 베타 배포 28bd1bb와 독립 Chrome QA 9584aaf 완료. 다음 백로그는 002다. 현재 어댑터는 synthetic stub이며 실제 Grok Bot 연결 증거가 없다. Jev find absent는 새 통합 부재의 후보이며 코드 자체가 없다는 뜻이 아니다.

## 적용 기준과 예외

fullops-common-0.3.2 README 및 coding-style/testing/security, FULLOPS.md, project.md, 기술 정본 architecture.md와 tech-stack.md를 기준 커밋에서 읽는다. Ponytail full을 적용한다. 기존 frozen wire·권한·공유 claim·무정책 deny를 유지한다. 기술 계획과 구현 판단은 DEV가 같은 과제에서 맡는다.

## 먼저 읽을 문서

Jev 근거: docs/evaluations/jev/SAR-MVP-002-DEV-{find,documents-find,context,route}.json. context keep 전체를 읽는다. 경로는 레포 기준이다.

- adapters/src/index.ts, adapters/src/synthetic.ts, adapters/package.json
- scripts/verify_mvp.py
- .fullops-squad/docs/design-docs/interface-design.md
- .fullops-squad/docs/design-docs/architecture.md, .fullops-squad/docs/design-docs/tech-stack.md
- .fullops-squad/docs/planning/product-specs/SAR-MVP.md
- .fullops-squad/docs/planning/SAR-MVP-backlog.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/product.md, protocol.md
- .fullops-squad/contexts/dev.md, project.md, FULLOPS.md
- .fullops-squad/rules/common/README.md 및 연결된 세 규칙
- 현재 역할 인박스

## 해야 할 일과 파일 소유권

- [ ] 원천의 Grok Bot 대상을 확인하고 현재 설치된 Grok CLI와 같은 제품인지 근거로 구분한다. 지원 외부 인터페이스·버전·인증·owner 권한을 조사한다. SDK/API는 Context7 resolve→query를 우선한다. 미지원이면 공식 문서/공식 소스로 확인하고 출처를 남긴다. 추정 inbound API를 만들지 않는다.
- [ ] 같은 과제 실행 기록에 짧은 기술 계획과 가능/불가 판정을 적는다. 가능할 때 필요한 TypeScript pull-default 어댑터와 합성 로컬 검증을 구현한다. 불확실한 제품 정체성·제품 변경은 coor에 질문한다. 설치된 Grok CLI를 이름만으로 대상이라고 확정하지 않는다.
- [ ] 실제 지원 경로에서 수신·durable persist·ACK/claim·gate·최소 result의 기존 안전 경계를 검증한다. 외부 실행 권한이 없으면 준비와 합성 검증까지만 진행하고 실제 연결 held를 남긴다. 불가한 경우 근거·최소 대안·재개 조건을 구체적으로 보고한다.
- [ ] 영향을 받은 D05/D10 및 필요 D03/D13 기술 문서와 contexts/dev.md, docs/exec-plans/phases/SAR-MVP-002-DEV.md를 갱신한다. 기획 정본과 PLANS/board는 coor에게 변경 요청한다.
- [ ] 변경 영향에 맞는 자동 검사·회귀·lint·strict를 실행하고 자신의 종료코드와 실제 HEAD를 남긴다. work.py finish로 전문 보존·인박스 비우기·커밋·역할 브랜치 일반 push 뒤 worker_done을 보낸다.

파일 소유권: 기존 DEV 제품 경로 adapters/, internal/, cmd/, scripts/, 관련 설정과 기술 문서. 운영 배포·Cloudflare·공유 서비스·전역 도구/사용자 설정은 수정하지 않는다.

## 완료 기준과 검증

지원 가능 여부, 실제 대상 정체성, 현재 버전 공식 근거, 인증·권한·로컬 검증과 실제 연결의 차이가 명확해야 한다. 지원되는 경로가 있을 때 가능한 준비 코드를 끝까지 구현한다. 불가하면 없는 인터페이스를 가짜 성공으로 대신하지 않는다.
DEV는 변경 동작 직접 검증을 수행한다. 코드 변경은 고정 후보 SHA에서 구현자와 다른 세션의 독립 리뷰와 TESTER 독립 QA 뒤 coor가 수락한다. UI 변경 없으면 기존 시각 증거를 재사용하고 근거 SHA를 구분한다. 문서만인 경우 근거·규약 검토를 적용한다. 기존 실패/held를 보존한다.

## 갱신할 산출물

route 추천 D05/D10/D13. 실제 영향만 갱신하며 D03 필요 시 추가한다. front matter는 deliverables.py stamp로 유지한다.

## 제약·협업·후속

실제 외부 업무 발송·실데이터·유료 API 호출·운영 연결 활성화는 이번 다음작업 승인에 포함하지 않는다. 읽기 전용 공식 조사, 로컬 코드·합성 검증·커밋·push는 승인됐다. 실제 연결용 비밀/계정 필요 시 준비를 완성한 뒤 최소 필요사항을 보고한다. 기존 베타와 Tailscale을 유지한다. DEC-02 미정 정책, positive silent done, calendar effect, 외부 exactly-once를 확정하지 않는다. 정책 변경이나 대상 우선순위 변경은 제품 결정이다.

## 완료 보고

브랜치·고정 SHA·변경 이유·공식 근거·검증 HEAD/명령/종료코드·실행 불가·후속과 재개 조건을 적는다. worker_done body 첫 줄은 `[완료] SAR-MVP-002-DEV | SHA <전체 완료 SHA>`로 한다.
