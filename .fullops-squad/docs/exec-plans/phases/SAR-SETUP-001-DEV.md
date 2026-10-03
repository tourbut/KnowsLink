---
title: 초기 구성과 lint 구현 실행 기록
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV]
summary: 기술 계획·수락 기준별 검증·실행 한계와 후속 수락 절차를 기록한다
---

# SAR-SETUP-001-DEV — 초기 구성과 제품 lint

## 기술 계획

기준 ref는 `729446d8da57`이다. 규칙은 `fullops-common-0.3.1`이다. D02와 고정 원천의 Go·TypeScript·Postgres·Compose 결정을 유지한다.

1. 루트 Go module과 `cmd/relay`, `cmd/migrate`, `internal/config`를 구성한다. relay는 DB ping과 `/healthz`만 제공한다. migration은 별도 프로세스의 goose SQL-only 경로로 실행한다.
2. 벤더 연결 없는 TypeScript adapter를 구성한다. npm lock, ESLint, Prettier, TypeScript 검사를 고정한다.
3. Compose의 네 서비스를 구성한다. Postgres 포트는 게시하지 않는다. relay 게시 주소는 loopback이다. cloudflared는 선택 profile이다.
4. `make lint`로 Go 포맷·vet, TypeScript 포맷·lint·타입과 설정 검사를 실행한다. 임시 복제본의 위반 주입으로 실패 전파를 검증한다.
5. 직접 테스트·빌드·Compose 검사 결과와 한계를 기록한다. D03과 개발 정본을 갱신한다. 완료 기록과 코드를 커밋한 뒤 기준 ref의 FullOps lint를 실행한다.

업무 스키마와 쿼리는 만들지 않는다. sqlc 설정만 마련하고 실제 생성 검사는 미적용으로 기록한다. 빈 migration은 명시적 no-op이며 실제 migration 적용 성공으로 보고하지 않는다.

## 적용 스킬과 문서 근거

caveman full은 세션 응답에만 적용한다. ponytail full, orchestration, fullops-work, fullops-test, fullops-deliverables를 적용한다. 영속 문서는 한국어 STE와 front matter 규칙을 따른다.

지시서의 keep 문서, 공통 규칙 세 문서, D02, 원천 product·architecture·protocol·decisions·mvp-checklist를 직접 읽었다. 제품 코드는 착수 HEAD `dbe0b40`에 없다. 기존 Jev 탐색 실패 기록을 재사용했다. 원천은 수정하지 않는다.

## 검증 결과

브랜치는 `fullops/dev`다. 검증 환경은 Linux amd64, Go 1.27.1, Node 22.22.2, npm 10.9.7, Python 3.12.3, GNU Make 4.3, Docker 29.4.3, Compose 5.1.3, Postgres 17.11이다. 직접 구현과 테스트를 수행한 Dispatch는 `ctx_67f98ed4cd42`, Task는 `task_a7fd5d1b8806`이다.

| D02 ID | 결과와 명령 | 종료코드·증거 |
|---|---|---|
| SETUP-01 | `make install`, `make build`; Go module과 npm lock 사용 | 0, [install.log](../logs/SAR-SETUP-001-DEV/install.log), [checks.log](../logs/SAR-SETUP-001-DEV/checks.log); 깨끗한 체크아웃 재현 기록은 후속 절에 추가 |
| SETUP-02 | `make test`, `make verify-runtime`; 별도 migrate와 no-op, adapter 종료 확인 | 0, [checks.log](../logs/SAR-SETUP-001-DEV/checks.log), [runtime.log](../logs/SAR-SETUP-001-DEV/runtime.log) |
| SETUP-03 | Compose config 구조 검사와 실제 로컬 컨테이너 기동 | 0, config는 네 서비스, runtime은 Tunnel 비활성의 postgres·migrate·relay |
| SETUP-04 | 필수 Compose 값 누락, relay·migrate URL 누락·DB 실패, 빈 Tunnel token | 오류를 검출했다. Go 프로세스 1, cloudflared 255; [runtime.log](../logs/SAR-SETUP-001-DEV/runtime.log) |
| LINT-01 | `make lint`의 Go 포맷·vet·module 무결성, TypeScript 포맷·ESLint·타입, 설정 검사 | 0, [checks.log](../logs/SAR-SETUP-001-DEV/checks.log) |
| LINT-02–03 | `make verify`; 등록된 `product-lint` 명령을 그대로 실행, 위반 네 종류·원복 검사 | 주입별 통합 명령 2, 원복 0, 검증 스크립트 0; [lint-injection.log](../logs/SAR-SETUP-001-DEV/lint-injection.log) |
| DOC-01 | D03·README·project.md 갱신, deliverables strict | 검사 13, 미작성 11, 문제 0, 경고 0; FullOps 최종 결과는 별도 완료 회신 |
| SCOPE-01 | 업무 endpoint HTTP 404 unit test, 미구현 adapter, 원천 변경 없음 | unit/race test 0; 실데이터 경로와 외부 연결 없음 |
| QA-01 | 완료 SHA 이후 tester가 독립 검증 | 후속 대기; 직접 검사로 독립 QA를 대체하지 않음 |

제품 lint는 자동 수정하지 않는다. 임시 복제본에서 Go 포맷 위반, TypeScript 미사용 변수, 타입 불일치 TS2322, 포맷 위반을 각각 검출했다. 각 주입에서 실제 진단과 비영 종료코드를 확인했다. 복제본 원복 후 통합 명령이 0으로 끝났다.

runtime 검증은 Compose 의존 순서와 실제 DB ping을 확인했다. 빈 migration 뒤 public 테이블 수는 0이었다. DB 중지 뒤 healthcheck는 HTTP 503이었다. 검증 컨테이너와 네트워크는 종료했다. 임시 tmpfs override는 검증 저장소만 바꿨다. 제품 Compose의 영속 volume은 수정하거나 삭제하지 않았다.

로그는 subprocess의 실제 returncode를 저장했다. 명령을 `tail`과 연결하지 않았다. 실행 뒤 로그 일부를 읽는 동작은 판정 명령이 아니다. [environment.log](../logs/SAR-SETUP-001-DEV/environment.log)에 실제 도구·의존성 버전이 있다.

## 실패 원인과 수정

첫 runtime 검증은 cloudflared 누락 오류의 기대 문구를 `credentials`로 잘못 지정해 실패했다. 실제 2026.9.1은 `requires the ID or name of the tunnel`로 실패하며 종료코드는 255다. 검사 기대 문구를 실제 필수 설정 오류에 맞춘 뒤 전체 runtime 검증을 다시 통과했다. [runtime-initial.log](../logs/SAR-SETUP-001-DEV/runtime-initial.log)는 실패 증거를 보존한다.

처음 조회한 ESLint 9.39.3은 npm 지원 종료 경고를 반환했다. 최종 구성은 지원되는 10.12.0으로 갱신했다. 최종 설치 감사는 취약점 0이다. Go의 `./...` 탐색은 node_modules의 외부 Go 파일도 포함했다. vet와 test 대상을 실제 제품 경로로 제한했다. 기존 FullOps 규칙과 exclude는 완화하지 않았다.

## 미적용과 후속

업무 SQL이 없어 실제 migration Up 적용과 sqlc 생성은 미적용이다. 빈 SQL의 no-op을 migration 적용 성공으로 보고하지 않는다. 운영 Tunnel·배포·네 벤더 통합과 업무 MVP는 범위 밖이다. UI가 없어 직접 시각 검수는 미적용이며 UI 구현 뒤 수행한다.

독립 QA와 고정 SHA의 독립 리뷰는 coordinator의 후속 과제다. 미해결 critical/high 차단은 유지한다. D03은 `review`이며 제품 수락이나 병합 승인으로 표시하지 않는다.
