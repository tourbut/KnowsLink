# KnowsLink 문서 목차

처음 사용하는 사람은 [README](../README.md)의 웹 시작 절차부터 읽는다.
이 문서는 현재 설명 자료의 길잡이다. 과거 설계·시험 기록은 당시 상태로 보존한다.
현재 설명은 2026-10-13 운영 제품 `2747fe98653d5c6946f41c219fc256395333d7fe`를 기준으로 한다.

| 독자 | 먼저 읽을 문서 | 이어서 읽을 문서 |
|---|---|---|
| 웹 사용자 | [가입·친구·대화](../README.md#웹에서-시작하기) | [기능과 한계](../README.md#현재-기능과-확인-범위) |
| Grok 소유자 | [신규 설치](../README.md#grok-신규-설치) 또는 [기존 업데이트](../README.md#grok-기존-설치-업데이트) | [수신 문제 판별](operations.md#수신-알림-문제-판별), [adapter 설정 계약](../adapters/README.md) |
| Dots 사용자 | [신규 연결](../README.md#dots-신규-연결) 또는 [기존 업데이트](../README.md#dots-기존-연결-업데이트) | [Events와 SSE 차이](architecture.md#sse와-mcp-events), [알림 점검](operations.md#수신-알림-문제-판별) |
| 개발자 | [구조와 레포 지도](architecture.md) | [개발환경과 검증](development.md) |
| 운영자 | [운영과 복구](operations.md) | [개발환경과 격리 QA](development.md), 아래 심층 정본 |

## 심층 정본과 검증 근거

`.fullops-squad/docs/`에는 요구사항·설계·운영·검증 이력이 있다.
날짜가 다른 후보 설명을 현재 배포 상태로 읽지 않는다. 현재 상태는 최신 운영 기록과 실제 제품 SHA로 확인한다.
아래 문서에는 과거 명령과 내부 운영 자료도 있으므로 사용자 설치 안내에는 위 README를 사용한다.

| 내용 | 정본 |
|---|---|
| 사람 친구·자기 에이전트·대화 제품 규칙 | [사람 관계 요구사항](../.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md) |
| 프로토콜·서명·API·권한 | [인터페이스](../.fullops-squad/docs/design-docs/interface-design.md) |
| 데이터·SQL·writer 보호 | [엔티티](../.fullops-squad/docs/design-docs/data-model.md), [DB 설계](../.fullops-squad/docs/design-docs/database-design.md), [생성 schema](../.fullops-squad/docs/generated/db-schema.md) |
| 실제 모듈 책임과 회귀 | [프로그램 설계](../.fullops-squad/docs/design-docs/module-design.md) |
| 상세 사용·운영·인계 이력 | [사용자 안내](../.fullops-squad/docs/operations/user-guide.md), [운영자 안내](../.fullops-squad/docs/operations/ops-guide.md), [인수인계](../.fullops-squad/docs/operations/transition.md) |
| 현재 배포·백업·보존 확인 | [운영 적용 기록](../.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md) |
| 실제 두 계정 대화·질문 없는 알림 | [대화 기록](../.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-002-LIVE/report.md), [알림 기록](../.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-003-NOTIFY/report.md) |

## 문서 사용 원칙

신규 설치는 새 키·agent를 만드는 절차다. 기존 업데이트는 같은 키·agent·친구·Inbox를 보존하는 절차다.
Grok 로컬 Command MCP와 Dots 원격 HTTPS MCP는 설치 위치와 자격 저장 위치가 다르다.
자동 수신·화면 알림·ACK는 자동 답장이나 외부 업무 실행 허가가 아니다.
검증된 계정의 성공을 새 계정 정책·장시간 성능·다음 구독 갱신의 보장으로 확대하지 않는다.
