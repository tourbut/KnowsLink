---
title: KnowsLink 기능 단위 백로그
status: draft
updated: 2026-10-05
owner: designer
tasks: [SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-PUBLIC-SERVICE-001]
summary: 일반 서비스 기능 순서와 현재 인박스 및 후속 대기와 기존 held를 연결한다
---

# KnowsLink 기능 단위 백로그

## 현재 일반 서비스 작업 순서 — SAR-PUBLIC-SERVICE-001

현재 제품 정본은 [일반 서비스 D02](product-specs/SAR-PUBLIC-SERVICE.md)이며 기존 MVP와 frozen 안전 경계는 함께 적용한다. 사용자 지시는 일반 서비스 완성 뒤 OpenAI dot “다닷” 연결이다. 아래 이전 MVP 기능·DEC 표는 당시 범위와 held 기록이다. 이번 순서와 DEC-03 기본값이 충돌하는 항목에 우선한다. 과거 미승인 수치를 승인됐다고 소급하지 않는다.

| 키·담당 | 제품 완료 조건 | 선행·현재 상태 |
|---|---|---|
| SAR-PUBLIC-IDENTITY-001-DEV, dev | PS-01–04, 신원/세션 남용 제한과 자기 owner UI. 일반 이메일 실제 확인·재로그인·교차 계정 거부 | [DEV 인박스](../../handovers/to_dev.md) ready. coor가 최신 main 포함과 dispatch 기준 SHA를 고정한다. |
| SAR-PUBLIC-IDENTITY-001-TESTER, tester | 같은 후보의 일반 이메일·세션·회원 분리 독립 QA | [TESTER 인박스](../../handovers/to_tester.md) waiting. DEV 고정 후보와 허가된 일반 이메일 환경 이후 실행한다. |
| SAR-PUBLIC-SERVICE-OPS-READINESS, ops | 신원 제공자·owner-origin 경계·공유 서비스·실제 자원·복구 근거 | 완료 SHA `0313deae0dec9af813b70ee9685e1a6d0a2b84d7`를 확인했다. 신원 바인딩·자원/처리량·백업 복구 미충족을 후속 공개 조건으로 반영한다. OPS 후속 인박스 작성은 coor가 한다. |
| SAR-PUBLIC-AGENTS-001-DEV, dev | PS-04–07·연결/키/관계 한도와 계정 비활성화. 일반 사용자가 자기 클라이언트에서 연결·회전·철회 | identity 기능 수락 뒤 coor가 PLANS 대기에서 현재 DEV 인박스로 옮긴다. |
| SAR-PUBLIC-MESSAGES-001-DEV, dev | PS-08–11, 실제 일반 신원 두 클라이언트 왕복·실패·기존 gate 회귀 | agents 기능 수락·지원 실제 클라이언트 확인 뒤 준비한다. |
| SAR-PUBLIC-SERVICE-OPS, ops | PS-12/13 배포 후보·자원 보호·일반 신원 공개 경계·격리 복구·기존 owner/공유 회귀 | readiness 완료 아카이브 뒤 빈 OPS 인박스에 작성한다. 수락된 구현 SHA·독립 검증 전 공개 변경은 하지 않는다. |
| SAR-PUBLIC-SERVICE-ACCEPT-001, coor/tester/designer/reviewer | PS-01–13 전체 수락. fixed-SHA 독립 QA·read-only 별도 세션 리뷰·직접 UI/일반 이메일 사람 확인·운영 근거 | 각 기능 후보를 수락한 뒤 하나의 운영 후보에서 확인한다. 미해결 critical/high와 필수 실패는 차단한다. |
| SAR-DOTS-DADAT-001, dev/ops/tester | PS-14. 동일 이메일의 Grok Bot “노우”↔다닷 일반 온보딩·별도 agent 자격·지원 앱 도구/연결 상태·양쪽 관련 ID 왕복 | 일반 서비스 수락 이후, 계정 화면·사용자 연결 동의 확보. 자동 wake·이벤트는 별도 제품 범위 결정이다. |

새 과제는 각 역할 인박스가 비면 작성한다. PLANS·board의 관리자는 coor다. designer는 이 표와 [실행 기록](../exec-plans/phases/SAR-PUBLIC-SERVICE-001.md)에 대기 제안을 남기며 coor가 PLANS에 반영한다. 과제명 파일·logs를 현재 지시서로 dispatch하지 않는다.

새 DEC-03은 일반 서비스 D02의 2026-10-05 운영 기본값이다. DEV/OPS 기술 근거가 공개 조건이다. DEC-01 상품·DEC-02 실일정 disclosure는 held다. DEC-04는 실제 지원 클라이언트·다닷 인터페이스 확인이다. 기존 수동 시험을 자동 wake 성공으로 바꾸지 않는다. 사용자 이메일 입력·다닷 계정 연결·새 비용/실데이터 별도 승인은 해당 사람 확인/외부 경로만 대기시킨다.

## 이전 MVP 실행·결정 기록

## 실행 기준

제품 정본은 [D01](business-plan.md)과 [새 MVP D02](product-specs/SAR-MVP.md)다. 원천은 [service-design 7bc9ea1](sources/silent-agent-relay/source.json)이다. [기존 setup D02](product-specs/SAR-SETUP-001.md)는 초기 골격 완료 이력으로 보존한다. setup을 다시 만들지 않는다.

SAR-PREP-002는 개발 준비까지 완료한 이력이다. 2026-10-03 사용자 구현 시작 승인으로 첫 DEV 지시서는 ready다. 실제 진행과 배정은 coor가 관리한다. TESTER는 고정 구현 후보를 기다린다. 기획 문서가 구현 완료를 뜻하지 않는다. 기술 계획만의 별도 승인 과제를 만들지 않는다. DEV가 각 기능의 기술 계획·구현·관련 회귀·기술 문서 갱신을 함께 맡는다.

P0는 첫 로컬 안전 흐름이다. P1은 실제 제품 연결과 정책을 갖춘 유용한 업무 흐름이다. P2는 운영 공개 준비다. 우선순위는 일정·가격·quota 확정이 아니다.

## 기능 목록

| 기능 키·우선순위 | 목표·제품 규칙 | 포함 / 제외 | 사용자 완료 조건 | 담당·선행 조건·상태 |
|---|---|---|---|---|
| SAR-MVP-001, P0 | 승인 관계에서 합성 요청을 안전 전달하고 owner가 gate를 결정한다. D02 MVP-01–16, C1–C5와 frozen wire 유지 | 등록·키/철회·초대/수락·contacts, ingest·queue·receipt, pull stub 1개·ACK/claim, result·approve/deny UI 포함. 실데이터·실제 calendar/tool effect·벤더 연결·운영 배포 제외 | 로컬 두 agent/owner가 관계를 만든다. 정책 없는 query는 정보를 공개하지 않는다. 검증된 M에 묶인 gate의 approve/deny와 만료·철회를 확인한다. approve 뒤에도 stub 실행·공개 정책 우회가 없다. A가 receipt/transport/처리 결과를 구분한다 | dev 구현, tester 독립 QA, designer 직접 UI 검수, coor 독립 리뷰·수락 조정. 준비 문서 반영과 coor의 명시적 구현 시작 지시가 선행. DEV ready, TESTER는 DEV 고정 SHA 대기 |
| SAR-MVP-002, P1 | Grok Bot 외부 인터페이스로 동일 안전 흐름을 실제 어댑터에 연결한다. MVP-06–14 | TypeScript pull-default와 지원 인터페이스 확인 포함. 벤더 코어 수정·공식 inbound API 가정·실제 위험 효과 제외 | 해당 제품에서 요청 수신·durable persist·ACK/claim·human-gate와 결과 경로가 증명된다. 사용 불가 인터페이스를 성공으로 표시하지 않는다 | dev, tester. 001 수락·인터페이스와 owner 권한 확인·사용자의 실제 연결 승인 뒤 재개. queued, 인터페이스 미확인 |
| SAR-MVP-003, P1 | 정책이 허용하는 일정 조회를 조용히 완료한다. MVP-07/08/11/15 | pair별 disclosure·반환 schema·도구 gateway 집행·허용 silent done 포함. schedule.commit effect·임의 한도·무정책 정보 공개 제외 | 정책 허용 범위의 query가 B-human 미통지로 done을 반환한다. 범위 밖·철회·누적 제한 위반은 거부한다. 결과에는 허용 필드만 있다 | designer 정책 결정, dev 구현, tester 검증. 001 수락과 DEC-02 확정이 선행. held |
| SAR-MVP-004, P1 | Claude Code 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | TypeScript pull 어댑터 포함. 원천 wire 변경·새 업무 정책 제외 | 해당 제품 연결과 기존 safety 회귀를 고정 SHA에서 확인한다. 공유 조정 없는 다중 adapter 활성화는 없다 | dev, tester. 002의 Grok Bot 우선순위 처리 이후 인터페이스 확인. queued |
| SAR-MVP-005, P1 | Codex 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | 004와 같은 안전 경계 | 실제 연결·gate·result와 다중 adapter 조정 조건을 증명한다 | dev, tester. 004 우선순위 처리와 인터페이스 확인. queued |
| SAR-MVP-006, P1 | Dots 대상에 동일 사용자 흐름을 연결한다. MVP-06–14 | 004와 같은 안전 경계 | 실제 연결·gate·result와 기존 회귀를 증명한다 | dev, tester. 005 우선순위 처리와 인터페이스 확인. queued |
| SAR-MVP-007, P2 | 수락된 기능을 집 미니서버의 공개 ingress에 안전하게 준비한다. MVP-12/15, 원천 배포 잠금 | Compose·Postgres 비공개·별도 migration·Tunnel ingress·운영 문서 포함. 제한 없는 공개·승인 없는 배포·유료화 제외 | ops가 설정·수락 SHA·안전 제한·운영 절차를 검증한다. `link.knowslog.com`의 누구나 가입 가능한 첫 인증·합성 파일럿은 기존 사용자 배포 승인에 따라 수락과 확정 기준 충족 뒤 실행한다 | ops, dev는 필요한 제품 제한 구현, tester는 운영 후보 검증. DEC-03/05, 독립 QA·직접 UI 검수·독립 리뷰·수락 선행. held |

002·004–006은 어댑터 확장 순서를 유지한다. MVP 최소 stub 1개 수락에 네 벤더의 실제 연결을 모두 요구하지 않는다. 인터페이스가 불가능하면 담당 dev가 근거와 대안을 coor에 보고한다. 대상 우선순위 변경은 designer의 제품 결정이다.

001은 ingest/queue만 끝내지 않는다. 세부 구현 순서는 dev의 같은 과제 기술 계획이다. 초기 골격과 합성 fixture는 업무 요청 실행 완료 증거가 아니다. 003의 실데이터 silent 성공은 001의 종료 조건으로 당기지 않는다.

## 첫 기능의 인계

- 준비 지시서: [SAR-MVP-001-DEV](../../handovers/to_dev.md), [SAR-MVP-001-TESTER](../../handovers/to_tester.md).
- 현재 실행 기준: 사용자의 구현 시작 승인은 확보됐다. coor가 현재 지시서·규칙과 승인 범위로 001-DEV를 배정하고 진행을 관리한다. 이번 제품 정책 기록은 기존 DEV 범위를 확장하지 않는다.
- dev는 실제 동작·API·데이터·모듈에 따라 D03, 필요한 D05–D10을 작성한다. D08은 실제 업무 schema와 생성 명령이 있을 때만 만든다.
- DEV 성공 회신 뒤 coor가 고정 구현 SHA와 독립 QA의 통합 후보 SHA를 기록한다. tester는 그 후보를 독립 검증한다. designer는 같은 UI 후보를 직접 검수한다. 독립 코드 리뷰는 구현자와 다른 세션의 깨끗한 read-only detached snapshot에서 한다.
- DEV 완료와 001 기능 수락, 전체 MVP 수락, 운영 공개는 다른 상태다. 미해결 critical/high·held·정책 보류를 완료로 바꾸지 않는다.

## 미정 결정과 재개 조건

| ID | 미정 항목 | 결정 담당 | 막는 경로·현재 안전 동작 | 재개 조건 |
|---|---|---|---|---|
| DEC-01 | Free N·Pro 가격·slot-unit refinement | designer가 coor 경유 사용자 결정 수집 | 상품·슬롯 quota·유료화. pending이 active slot을 소모하지 않는 불변식만 유지 | N·가격·단위의 사용자 결정과 원천 변경 근거 확보 |
| DEC-02 | schedule.query disclosure, 부모 intent별 result/error schema, 반환 필드·window·granularity·누적 한도 | designer, dev가 기술 근거 제공 | 실데이터 조회·정보 반환·positive silent done. 현재 policy 없음은 deny, optional result/error 공개는 보류 | 필드·수치·권한·누적 정책과 승인 근거 기록 뒤 003 재개 |
| DEC-03 | 누구나 가입 가능한 인증 공개 서비스의 resource/rate/추가 size/concurrency 상한 | designer가 권장안 작성·coor 경유 사용자 결정 수집, DEV/OPS가 측정·운영 근거 제공 | 공개 범위·거부 원칙 결정, 숫자는 제안 단계. [D02 권장안](product-specs/SAR-MVP.md#초기-안전-한도-권장안--사용자-승인-전-제안값). 실제 공개 held, 로컬 합성 구현·OPS 읽기 전용 준비는 계속 | 사용자 승인/조정 값과 안전 종료·서버 자원 보호 근거를 반영해 수치·단위·대상·거부 확정, 제한 구현·독립 QA·UI 검수·독립 리뷰·수락 뒤 공개 |
| DEC-04 | 제품별 실제 외부 인터페이스·통합 세부 | dev | 실제 어댑터 연결. 로컬 stub 성공만 주장 | 지원 인터페이스·인증·현재 버전의 공식 근거·권한 확인. 제품 방향 변경 필요 시 designer에 질문 |
| DEC-05 | 미니서버 port/path/process manager, image/digest/env, Tunnel 운영 설정·공개 시점 | dev/ops가 기술 설정 결정, coor가 실행 승인 수집 | 운영 Tunnel·배포·외부 발송. 기존 로컬 설정은 골격 증거로만 재사용 | 수락된 고정 SHA·확정 제한 구현·운영 설정 검증. 현재 서버 Docker·Tunnel 첫 파일럿 승인은 확보됐으며 범위 확대에는 새 승인 필요 |

DEC-01/02/04의 기존 held·queued와 재개 조건은 유지한다. DEC-03의 초기 수치는 제안이며 승인 전 기본값이 아니다. DEC-05의 서버 Docker·Tunnel 배포 승인은 확보됐고 hostname은 link.knowslog.com이다. 수락 SHA·제한 집행·운영 검증은 남아 있다. 각 보류는 해당 경로만 막는다. 001의 합성 데이터·로컬 폐쇄 검증과 문서 준비는 진행할 수 있다. 공개 정책이나 quota를 fixture 숫자로 확정하지 않는다. 공개 가입은 사용자 결정이며 가입을 초대 전용으로 바꾸지 않는다. 실제 공개와 전체 MVP 완료는 구분한다.

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

- 2026-10-03: SAR-MVP-PUBLIC-POLICY-001에서 공개 가입·link.knowslog.com·인증 합성 파일럿 범위와 DEC-03 제안값·확정 조건을 정리했다. 사용자 승인과 기술 근거 전에는 공개 held를 유지한다.

- 2026-10-05: 일반 서비스의 검증 가능한 기능 순서·ready 인박스·진행 OPS 보존과 PLANS 후속 제안을 추가했다.
