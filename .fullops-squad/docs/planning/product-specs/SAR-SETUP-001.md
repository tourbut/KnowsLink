---
id: D02
title: KnowsLink 초기 구성과 lint 요구사항
status: review
updated: 2026-10-03
owner: designer
tasks: [SAR-SETUP-001]
downstream: [D03]
summary: "고정 원천의 초기 구성과 lint 범위, 수락 기준, 역할별 산출물을 정의한다"
---

# SAR-SETUP-001 — 프로젝트 초기 구성과 lint 요구사항

## 목적과 근거

KnowsLink는 승인된 에이전트 사이의 선택적 인간개입 전달 제품이다. 연락처 단위는 에이전트이며 사람은 owner다.
이번 과제의 사용자는 이후 기능을 구현하고 검증할 개발자다. 목표는 고정 기획에 맞는 개발 시작점과 실제 실패를 검출하는 lint를 제공하는 것이다.
사용자 요청은 [요청 원문](../SAR-SETUP-001-request.md)이다. 전체 MVP 구현이나 운영 배포를 요청한 것으로 확대하지 않는다.

- 원천 위치: `.fullops-squad/docs/planning/sources/silent-agent-relay/`. 요청의 `sources/silent-agent-relay/`는 이 디렉터리를 가리킨다.
- 원천 커밋: `404ff834c0607055d63d2053bf7771d2f46ad3ae`. [source.json](../sources/silent-agent-relay/source.json)에 고정했다.
- 작업 기준 ref: `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`. 요청·원천은 준비 커밋 `0cc35f0`에 있다.
- 적용 기준: `fullops-common-0.3.1`, `.fullops-squad/project.md`, 문서 작성 규칙. 보안·검증 예외는 없다.
- 문서 상태는 검토 대기다. 원천의 확정 결정은 유지하며 이번 문서에 별도 사용자 승인을 받았다고 주장하지 않는다.

## 원천에서 유지할 결정

| 항목 | 확정 내용 | 근거 |
|---|---|---|
| 이름·제품 | 공개명 KnowsLink, 내부명 `silent-agent-relay`, 호스트 `relay.knowslog.com`; 채팅·장기 메시지 아카이브가 아님 | [product.md](../sources/silent-agent-relay/product.md) |
| 런타임 | Go relay, TypeScript adapters/plugins, Postgres | [architecture.md](../sources/silent-agent-relay/architecture.md) |
| 패키징 | 집 미니서버의 Docker Compose: `postgres`, one-shot `migrate`, `relay`, `cloudflared`; Tunnel은 ingress 전용, Postgres 공개 노출 금지 | [decisions.md](../sources/silent-agent-relay/decisions.md) |
| DB 도구 | `jackc/pgx/v5` + `pgxpool`, `pressly/goose/v3` SQL-only, 별도 `cmd/migrate`, `sqlc`의 `sql_package: pgx/v5` | [architecture.md](../sources/silent-agent-relay/architecture.md) |
| UI | 향후 human-gate는 Go `net/http` + `html/template`; 별도 TypeScript frontend 없음 | [product.md](../sources/silent-agent-relay/product.md) |
| 어댑터 | Grok Bot, Claude Code, Codex, Dots 순서; TypeScript pull-default; 공식 inbound API 가정과 벤더 코어 수정 금지 | [mvp-checklist.md](../sources/silent-agent-relay/mvp-checklist.md) |
| 보안 | frozen `relay.v1` 및 C1–C5 유지; evidence fetch/preview와 webhook은 MVP OFF | [protocol.md](../sources/silent-agent-relay/protocol.md) |

`project.md`의 기술 스택 미정은 원천 반영 전 상태다. dev는 잠긴 선택을 다시 결정하지 않는다.
Go·Node·도구 버전, 패키지 관리, 디렉터리 배치, 이미지, 개발용 포트, 설정 전달 방식, lint 도구는 dev가 결정하고 근거를 남긴다.
Workers/DO, SQLite, 별도 프런트엔드, Atlas auto-apply, ORM AutoMigrate, Go-code migration으로 대체하지 않는다.

## 이번 범위

1. Go relay와 TypeScript adapter의 최소 개발 골격, 의존성·버전·잠금 방식, 재현 가능한 설치·검사·빌드 명령을 구성한다.
2. relay와 별도 `cmd/migrate`의 최소 실행 진입점 및 잠긴 DB 도구의 설정 경계를 마련한다. 도메인 테이블이나 가짜 업무 쿼리를 도구 실행 목적으로 만들지 않는다.
3. 네 Compose 서비스의 구성과 개발용 설정 예시를 제공한다. 실제 Tunnel 연결이나 운영 자원 변경은 수행하지 않는다.
4. Go·TypeScript의 포맷·정적 검사·타입 검사와 추가한 설정에 필요한 검사를 구성한다. 하나의 문서화된 명령으로 필수 검사를 실행할 수 있어야 한다.
5. 기존 FullOps lint의 `commands`에 실제 제품 검사 명령을 연결한다. 기존 검사와 보안 기준을 낮추지 않는다.
6. README와 기술 정본에 설치, 검사, 로컬 실행, 중지, 필요한 환경 변수 이름, 미구현 기능과 검증 한계를 기록한다.

골격에는 기동 확인 또는 상태 확인에 필요한 최소 동작만 둔다. 업무 요청의 수락·승인·외부 실행을 성공한 것처럼 응답하지 않는다.
어댑터는 벤더 API에 연결하지 않는 최소 TypeScript 골격으로 충분하다. 네 제품의 어댑터 구현은 요구하지 않는다.
향후 첫 기능 MVP는 원천대로 human-gate까지 포함한다. 이번 설정 완료를 ingest-only MVP 또는 기능 MVP 완료로 표시하지 않는다.

## 수락 기준

| ID | 통과 조건 | 확인 증거 |
|---|---|---|
| SETUP-01 | 깨끗한 체크아웃에서 명시된 버전과 문서 명령으로 Go·TypeScript 의존성 설치 및 빌드가 재현된다 | dev SHA, 도구 버전, 명령, 종료코드, 생성물로 인한 Git 변경 여부 |
| SETUP-02 | Go relay와 별도 `cmd/migrate`, TypeScript adapter 골격이 존재하며 잠긴 기술 선택을 따른다 | 실제 경로와 기술 정본의 일치, 최소 실행 확인; API 시작 시 migration 자동 실행 없음 |
| SETUP-03 | Compose가 네 서비스를 정의하고 구성 검증에 통과한다. `migrate`는 one-shot이며 Postgres를 공개 게시하지 않는다 | Compose 구성 검사와 서비스·포트·의존 순서 확인; 실제 연결은 필요 없음 |
| SETUP-04 | 실제 키 없이 예시 설정으로 검사 가능하다. 필수값 누락 시 안전하게 실패하거나 해당 외부 서비스가 비활성이다 | 누락 설정 시나리오, 비밀값·실제 Tunnel token 미포함 확인 |
| LINT-01 | Go 포맷·정적 검사와 TypeScript 포맷·lint·타입 검사가 문서화된 통합 명령에 포함된다 | 정상 골격의 종료코드 0, 검사 대상과 제외 경로 |
| LINT-02 | 임시 Go 포맷 위반, TypeScript lint 위반 및 타입 오류가 각 검사에서 0이 아닌 종료코드로 검출된다 | 복제본 또는 임시 fixture에서 오류 주입·검출·원복 후 재통과; 관련 없는 환경 오류는 인정하지 않음 |
| LINT-03 | 통합 명령과 FullOps 등록 명령이 검사 실패를 숨기지 않는다. 변경하지 않은 파일·원천·의존성·생성물을 합리적으로 구분한다 | 명령 연결 확인, 정상·실패 실행, 자동 수정으로 실패를 감추지 않는 검사 모드 |
| DOC-01 | README·`project.md`·D03에 실제 경로, 기술 선택, 명령, 구현 상태가 일치한다 | 요구사항 ID와 구현·검증 매핑, 고정 원천 링크 |
| SCOPE-01 | 업무 relay, 승인 UI, 실데이터, 외부 발송·배포가 활성화되지 않는다. 원천은 변경하지 않는다 | 기준 ref 이후 diff, 실행 진입점 및 설정 검토 |
| QA-01 | tester가 dev 완료 SHA를 받은 뒤 독립 환경에서 위 기준을 재현하고 결과를 기록한다 | SHA가 적힌 시나리오·QA 보고서, 실패·미실행·통과 구분 |

SQL·sqlc 검사는 실제 도입한 파일 범위에 맞춘다. 업무 스키마가 없어 실행할 수 없는 생성 검사는 미적용 사유를 기록한다.
Compose 구성 검증과 실제 컨테이너 기동은 구분한다. DB 접속과 migration 실행을 검증했다면 별도 증거를 남긴다.
도구 미설치·네트워크 제한으로 필수 검사를 못 하면 미검증으로 기록하고 coordinator에게 전달한다. 전체 수락 완료로 보고하지 않는다.
제품 lint의 실패 검출은 필수다. FullOps lint는 merge-base의 설정을 사용하므로 새 `commands`를 등록한 것만으로 해당 명령의 실행을 증명하지 못한다.
따라서 dev와 tester는 새 제품 검사 명령을 직접 실행하고 별도로 FullOps lint 결과를 남긴다.

## 이번에 구현하지 않을 기능과 미결정 사항

- 에이전트 등록·키 관리·페어링, ingest·서명·멱등·receipt·큐·lease·ACK·공유 실행 claim, 승인 UI·승인 기록, 실제 도구 실행은 후속 기능 과제다.
- C1–C5가 코드로 구현·검증됐다고 주장하지 않는다. 실데이터 silent 경로를 열지 않는다.
- 일정 정보는 disclosure policy가 없으면 deny다. `schedule.commit`은 승인으로 활성화할 수 없는 stub이다.
- 원천의 고정값(`MAX_TTL` 300s, lease 30s, attempts 3, receipt 24h)을 변경하지 않는다. 이번 과제에서 해당 동작의 구현을 요구하지 않는다.
- Free N, Pro 가격, slot-unit, disclosure 수치, 추가 rate·size·concurrency 제한은 미정으로 유지한다. 무제한 운영 배포로 대신하지 않는다.
- 운영 이미지·포트 등 구현 설정은 dev가 정할 수 있다. 비즈니스·보안 정책 수치를 임의로 확정할 권한은 없다.
- CI·배포 자동화, 운영 Tunnel·DNS·미니서버 변경, 병원·폐쇄망 지원, 외부 메시지 발송은 필수 산출물이 아니다.

## 역할과 산출물 추적

| 역할 | 과제·소유권 | 산출물 |
|---|---|---|
| designer | SAR-SETUP-001; 범위·수락 기준, dev/tester 지시서 | D02 이 문서, 역할 인박스, 실행 기록 |
| dev | SAR-SETUP-001-DEV; 기술 설계·초기 구성·lint·직접 검증 | D03 `docs/design-docs/architecture.md`, `docs/design-docs/tech-stack.md`; README·project.md; 실제 코드·설정 |
| tester | SAR-SETUP-001-TESTER; dev 완료 SHA 이후 독립 검증 | `docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md`, `docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` |
| coor | 설계 커밋 전달, dev 완료 후 tester dispatch, 검토·병합 조정 | 진행 기록과 dispatch 정보 |

Jev route는 키 부재로 산출물을 선택하지 못했다. 원천과 과제 범위를 읽고 D02·D03만 선택했다.
D01은 사업 모델 변경이 없어 갱신하지 않는다. D04–D10은 화면·API·데이터·업무 모듈 구현이 없어 생성하지 않는다.
D11–D13은 사용자 기능·운영 배포·운영 이관이 없어 미작성 상태를 유지한다. 개발 환경 사용법은 README와 D03에 기록한다.
QA 원천과 상세 실행 기록은 D01–D13 외 문서로 관리한다. 원천 스냅샷은 산출물 완료 상태로 재분류하지 않는다.

## 개정 이력

- 2026-10-03: 고정 원천을 기준으로 초기 구성·lint와 향후 기능 MVP의 경계를 정했다. SETUP-01–04, LINT-01–03, DOC-01, SCOPE-01, QA-01을 등록했다.
