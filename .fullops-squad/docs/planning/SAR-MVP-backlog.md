---
title: KnowsLink 기능 단위 백로그
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
summary: 검증 가능한 MVP 기능의 우선순위와 queued 인계 및 결정 보류를 연결한다
---

# KnowsLink 기능 단위 백로그

## 실행 기준

제품 정본은 [D01](business-plan.md)과 [새 MVP D02](product-specs/SAR-MVP.md)다. 원천은 [service-design 7bc9ea1](sources/silent-agent-relay/source.json)이다. [기존 setup D02](product-specs/SAR-SETUP-001.md)는 초기 골격 완료 이력으로 보존한다. setup을 다시 만들지 않는다.

이번 과제는 개발 준비까지 완료한다. 아래 기능은 아직 구현되지 않았다. 첫 DEV/TESTER 지시서는 queued이며 이번에 dispatch하지 않는다. 기술 계획만의 별도 승인 과제를 만들지 않는다. DEV가 각 기능의 기술 계획·구현·관련 회귀·기술 문서 갱신을 함께 맡는다.

P0는 첫 로컬 안전 흐름이다. P1은 실제 제품 연결과 정책을 갖춘 유용한 업무 흐름이다. P2는 운영 공개 준비다. 우선순위는 일정·가격·quota 확정이 아니다.

## 기능 목록

| 기능 키·우선순위 | 목표·제품 규칙 | 포함 / 제외 | 사용자 완료 조건 | 담당·선행 조건·상태 |
|---|---|---|---|---|
| SAR-MVP-001, P0 | 승인 관계에서 합성 요청을 안전 전달하고 owner가 gate를 결정한다. D02 MVP-01–16, C1–C5와 frozen wire 유지 | 등록·키/철회·초대/수락·contacts, ingest·queue·receipt, pull stub 1개·ACK/claim, result·approve/deny UI 포함. 실데이터·실제 calendar/tool effect·벤더 연결·운영 배포 제외 | 로컬 두 agent/owner가 관계를 만든다. 정책 없는 query는 정보를 공개하지 않는다. 검증된 M에 묶인 gate의 approve/deny와 만료·철회를 확인한다. approve 뒤에도 stub 실행·공개 정책 우회가 없다. A가 receipt/transport/처리 결과를 구분한다 | dev 구현, tester 독립 QA, designer 직접 UI 검수, coor 독립 리뷰·수락 조정. 준비 문서 반영과 coor의 명시적 구현 시작 지시가 선행. DEV queued, TESTER는 DEV 고정 SHA 대기 |
| SAR-MVP-002, P1 | Grok Bot 외부 인터페이스로 동일 안전 흐름을 실제 어댑터에 연결한다. MVP-06–14 | TypeScript pull-default와 지원 인터페이스 확인 포함. 벤더 코어 수정·공식 inbound API 가정·실제 위험 효과 제외 | 해당 제품에서 요청 수신·durable persist·ACK/claim·human-gate와 결과 경로가 증명된다. 사용 불가 인터페이스를 성공으로 표시하지 않는다 | dev, tester. 001 수락·인터페이스와 owner 권한 확인·사용자의 실제 연결 승인 뒤 재개. queued, 인터페이스 미확인 |
| SAR-MVP-003, P1 | 정책이 허용하는 일정 조회를 조용히 완료한다. MVP-07/08/11/15 | pair별 disclosure·반환 schema·도구 gateway 집행·허용 silent done 포함. schedule.commit effect·임의 한도·무정책 정보 공개 제외 | 정책 허용 범위의 query가 B-human 미통지로 done을 반환한다. 범위 밖·철회·누적 제한 위반은 거부한다. 결과에는 허용 필드만 있다 | designer 정책 결정, dev 구현, tester 검증. 001 수락과 DEC-02 확정이 선행. held |
| SAR-MVP-004, P1 | Claude Code 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | TypeScript pull 어댑터 포함. 원천 wire 변경·새 업무 정책 제외 | 해당 제품 연결과 기존 safety 회귀를 고정 SHA에서 확인한다. 공유 조정 없는 다중 adapter 활성화는 없다 | dev, tester. 002의 Grok Bot 우선순위 처리 이후 인터페이스 확인. queued |
| SAR-MVP-005, P1 | Codex 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | 004와 같은 안전 경계 | 실제 연결·gate·result와 다중 adapter 조정 조건을 증명한다 | dev, tester. 004 우선순위 처리와 인터페이스 확인. queued |
| SAR-MVP-006, P1 | Dots 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | 004와 같은 안전 경계 | 실제 연결·gate·result와 기존 회귀를 증명한다 | dev, tester. 005 우선순위 처리와 인터페이스 확인. queued |
| SAR-MVP-007, P2 | 수락된 기능을 집 미니서버의 공개 ingress에 안전하게 준비한다. MVP-12/15, 원천 배포 잠금 | Compose·Postgres 비공개·별도 migration·Tunnel ingress·운영 문서 포함. 제한 없는 공개·승인 없는 배포·유료화 제외 | ops가 설정·수락 SHA·안전 제한·운영 절차를 검증한다. 실제 공개는 명시적 승인 후 별도 실행한다 | ops, dev는 필요한 제품 제한 구현, tester는 운영 후보 검증. DEC-03/05, 독립 QA·직접 UI 검수·독립 리뷰·수락 선행. held |

002·004–006은 어댑터 확장 순서를 유지한다. MVP 최소 stub 1개 수락에 네 벤더의 실제 연결을 모두 요구하지 않는다. 인터페이스가 불가능하면 담당 dev가 근거와 대안을 coor에 보고한다. 대상 우선순위 변경은 designer의 제품 결정이다.

001은 ingest/queue만 끝내지 않는다. 세부 구현 순서는 dev의 같은 과제 기술 계획이다. 초기 골격과 합성 fixture는 업무 요청 실행 완료 증거가 아니다. 003의 실데이터 silent 성공은 001의 종료 조건으로 당기지 않는다.

## 첫 기능의 인계

- 준비 지시서: [SAR-MVP-001-DEV](../../handovers/to_dev.md), [SAR-MVP-001-TESTER](../../handovers/to_tester.md).
- 지금 가능한 다음 실행 단위: 사용자가 구현 시작을 지시하면 coor가 준비 커밋·현재 지시서·규칙을 dev 워크트리에 반영하고 001-DEV를 배정한다. 이 과제에서는 배정하지 않는다.
- dev는 실제 동작·API·데이터·모듈에 따라 D03, 필요한 D05–D10을 작성한다. D08은 실제 업무 schema와 생성 명령이 있을 때만 만든다.
- DEV 성공 회신 뒤 coor가 고정 구현 SHA와 독립 QA의 통합 후보 SHA를 기록한다. tester는 그 후보를 독립 검증한다. designer는 같은 UI 후보를 직접 검수한다. 독립 코드 리뷰는 구현자와 다른 세션의 깨끗한 read-only detached snapshot에서 한다.
- DEV 완료와 001 기능 수락, 전체 MVP 수락, 운영 공개는 다른 상태다. 미해결 critical/high·held·정책 보류를 완료로 바꾸지 않는다.

## 미정 결정과 재개 조건

| ID | 미정 항목 | 결정 담당 | 막는 경로·현재 안전 동작 | 재개 조건 |
|---|---|---|---|---|
| DEC-01 | Free N·Pro 가격·slot-unit refinement | designer가 coor 경유 사용자 결정 수집 | 상품·슬롯 quota·유료화. pending이 active slot을 소모하지 않는 불변식만 유지 | N·가격·단위의 사용자 결정과 원천 변경 근거 확보 |
| DEC-02 | schedule.query disclosure, 부모 intent별 result/error schema, 반환 필드·window·granularity·누적 한도 | designer, dev가 기술 근거 제공 | 실데이터 조회·정보 반환·positive silent done. 현재 policy 없음은 deny, optional result/error 공개는 보류 | 필드·수치·권한·누적 정책과 승인 근거 기록 뒤 003 재개 |
| DEC-03 | resource/rate/추가 size/concurrency 상한 | designer가 정책 결정, dev가 측정·운영 근거 제공 | 무제한 운영 공개. 알려진 frozen field 한도는 유지하고 추가 수치를 만들지 않음 | 수치·단위·적용 범위·거부 동작·검증 근거 확정 뒤 제한 구현과 QA |
| DEC-04 | 제품별 실제 외부 인터페이스·통합 세부 | dev | 실제 어댑터 연결. 로컬 stub 성공만 주장 | 지원 인터페이스·인증·현재 버전의 공식 근거·권한 확인. 제품 방향 변경 필요 시 designer에 질문 |
| DEC-05 | 미니서버 port/path/process manager, image/digest/env, Tunnel 운영 설정·공개 시점 | dev/ops가 기술 설정 결정, coor가 실행 승인 수집 | 운영 Tunnel·배포·외부 발송. 기존 로컬 설정은 골격 증거로만 재사용 | 수락된 고정 SHA·제한 구현·운영 설정·명시적 실행 승인 확보 |

DEC-01–05의 값을 이번 기획에서 추정하지 않았다. 각 보류는 해당 경로만 막는다. 001의 합성 데이터·로컬 폐쇄 검증과 문서 준비는 진행할 수 있다. 공개 정책이나 quota를 fixture 숫자로 확정하지 않는다.

## 원천·요구사항·후속 증거 추적

| 근거 | D02 ID | 기능 | DEV/QA/시각 검수 연결 |
|---|---|---|---|
| product 등록·Pairing, protocol C1/C5 | MVP-01/02/07 | 001 | DEV 직접 권한·동시성 회귀, TESTER QA-01/02/06 |
| protocol 봉투·서명·Ingest·멱등·운영 상수 | MVP-03–06/14 | 001 | TESTER QA-03–05/10, DEV frozen fixture와 경계 검사 |
| product silent/escalate, protocol C2/C4 | MVP-08–11/15/16 | 001, 003 | TESTER QA-07–09/11, designer V-01–04 |
| product 어댑터·protocol C3/C5 | MVP-12/13 | 001, 002, 004–006 | TESTER QA-05/10, 각 실제 인터페이스 근거 |
| architecture 배포·저장, business-model TBD | MVP-12/15, DEC-01/03/05 | 007 및 미래 상품 과제 | ops 고정 후보 검증·배포 승인. 이번 준비에서 실행하지 않음 |

QA 보고서와 시나리오 파일은 tester 실행 시 작성한다. 현재 링크가 없는 미래 보고서를 완료 산출물로 표시하지 않는다.

## 개정 이력

- 2026-10-03: SAR-PREP-002에서 검증 가능한 기능과 결정 보류를 분리했다. 첫 기능 DEV/TESTER를 queued로 준비했다.
