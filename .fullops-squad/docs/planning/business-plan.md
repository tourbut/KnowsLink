---
id: D01
title: KnowsLink 서비스 개요
status: review
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
downstream: [D02]
summary: KnowsLink의 가치와 freemium 방향 및 미정 결정의 담당과 재개 조건을 정리한다
---

# KnowsLink 서비스 개요

## 가치와 대상 사용자

KnowsLink는 Knows 가족에서 승인된 에이전트끼리만 짧게 잇는 연결 제품이다. 개인 퍼블릭넷에서 두 에이전트의 owner가 관계를 허용한 뒤 구조화된 업무 요청을 전달한다. 연락처 단위는 에이전트다. 사람은 에이전트의 owner로 연결된다.

사용자는 상대 에이전트를 선택하고 처리 결과를 받는다. 수신 에이전트는 허용된 요청만 조용히 처리한다. 인간의 권한이나 판단이 필요한 요청은 owner에게 올린다. 결제·삭제·권한 부여·외부 발송·약속은 인간 게이트를 거친다. 채팅 화면이나 장기 메시지 아카이브는 제품 목표가 아니다.

제품 판단의 근거는 service-design `7bc9ea190ea549fae8b047e850247a19322fc9c3`이다. [원천 고정 기록](sources/silent-agent-relay/source.json), [제품 결정](sources/silent-agent-relay/product.md), [결정 로그](sources/silent-agent-relay/decisions.md)를 따른다. 원천 문서는 수정하지 않는다.

## 비즈니스 방향

[비즈니스 모델 원천](sources/silent-agent-relay/business-model.md)은 Tailscale형 freemium 방향을 잠갔다. Personal free는 페어링 슬롯 제한을 둔다. Personal Pro는 더 많은 슬롯을 제공한다. 설명 방향은 “친구 N명까지”다. Free N과 Pro 가격은 미정이다. 슬롯 단위의 세부 정의도 보류한다. 원천의 약 $10 메모는 위험 가설이며 확정 가격이 아니다.

MVP 피칭에서 메시지 볼륨 과금·Slack-seat형 과금·마켓플레이스 수수료는 생략한다. 결제나 구독 구현을 이번 준비 범위에 넣지 않는다. 양측 채택, 반복 사용, 지원·보안 비용은 검증할 가설이다. 가입자 수·매출·일정·성능 목표를 임의로 설정하지 않는다.

## MVP 범위와 개발 준비 상태

MVP는 등록·페어링·철회, frozen `relay.v1`, 짧은 TTL 큐, lease·ACK·공유 실행 claim, receipt, 인간 approve/deny까지 포함한다. 첫 Go 기능 구현을 ingest/queue-only로 줄이지 않는다. TypeScript 어댑터는 pull-default다. 대상 우선순위는 Grok Bot, Claude Code, Codex, Dots다. MVP의 최소 어댑터 범위는 pull stub 1개다. 실제 제품 연결은 외부 인터페이스 확인이 필요하다.

집 미니서버의 Go relay·Postgres·Docker Compose와 Cloudflare Tunnel ingress 방향은 원천의 확정 조건이다. 기술 계획과 구현 설정은 dev와 ops가 맡는다. 운영 배포와 외부 발송은 이번 요청의 승인 범위가 아니다.

SAR-SETUP-001 초기 골격·제품 lint·독립 QA·main 수락은 완료됐다. 업무 SQL, 사용자 등록·페어링 동작, human-gate UI, 실제 어댑터, 운영 Tunnel은 아직 완료되지 않았다. 이번 문서는 제품 기획의 개발 준비 결과다. 제품 MVP 구현 완료를 뜻하지 않는다.

## 사용자 성공 조건과 후속

owner는 에이전트와 키를 관리하고 상대 초대를 명시적으로 수락하거나 거절할 수 있어야 한다. 허용 관계 안에서 요청의 전달 상태와 처리 결과를 구분할 수 있어야 한다. owner는 검증된 요청 본문과 적용 정책을 보고 approve/deny를 결정할 수 있어야 한다. 철회·만료·중복 요청은 이전 실행 권한을 되살리지 않아야 한다.

상세 수락 기준은 [MVP 요구사항](product-specs/SAR-MVP.md)이다. 실행 순서와 첫 기능은 [기능 단위 백로그](SAR-MVP-backlog.md)다. 기술 설계 D03과 필요한 D05–D10은 후속 dev가 실제 구현과 함께 갱신한다.

## 미정 결정

| 항목 | 담당 | 영향 | 재개 조건 |
|---|---|---|---|
| Free N·Pro 가격·슬롯 단위 | designer가 사용자 결정을 coor 경유 수집 | 유료화·슬롯 상품 정책 보류 | 사용자 결정과 원천 변경 근거 확보 |
| 일정 disclosure·출력 allowlist·범위 수치 | designer, dev는 구현 가능성 근거 제공 | 실데이터 silent 조회·정보 반환 보류 | 반환 필드·window·granularity·누적 한도와 정책 승인 기록 확보 |
| 자원·rate·추가 size·concurrency 수치 | designer가 정책 결정, dev가 측정 근거 제공 | 무제한 공개 배포 금지 | 수치·단위·적용 범위·검증 근거 확정 |
| 실제 어댑터 인터페이스 | dev | 제품별 실제 연결 보류 | 지원 인터페이스·권한·통합 가능성 확인 |
| 운영 설정·배포 시점·외부 발송 | ops와 coor, 사용자가 실행 승인 | 운영 공개 보류 | 안전 제한·수락 SHA·운영 설정·명시적 실행 승인 확보 |

미정 결정은 개발 준비와 합성 데이터의 로컬 검증을 막지 않는다. 보류된 경로를 허용하는 기본값으로 바꾸지 않는다.

## 개정 이력

- 2026-10-03: SAR-PREP-002에서 최신 원천의 제품 가치·사업 방향·미정 결정과 개발 준비 범위를 정리했다.
