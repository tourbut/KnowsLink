---
title: SAR-SETUP-001-TESTER — 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-SETUP-001-TESTER]
summary: dev 완료 SHA 0cc10b0의 수락 기준별 독립 QA 결과와 한계
---

# SAR-SETUP-001-TESTER — 독립 QA 보고서

## 판정

- 대상 SHA: `0cc10b083771be9b3423833b222c57d426315333` (dev 완료 SHA). 깨끗한 독립 clone의 detached checkout에서 실행했다.
- 제품 수락 판정: 수락 가능. D02의 SETUP-01–04, LINT-01–03, DOC-01, SCOPE-01이 모두 통과했다. 미해결 critical/high 결함은 없다.
- 이 판정은 병합 승인이 아니다. 업무 SQL·sqlc·UI·실제 Tunnel은 미실행이다. MVP 보안·프로토콜 구현 완료를 뜻하지 않는다.
- 시나리오: [SAR-SETUP-001-TESTER.md](../scenarios/SAR-SETUP-001-TESTER.md). 로그: [SAR-SETUP-001-TESTER-test/](SAR-SETUP-001-TESTER-test/).
- DEV 로그는 재인용하지 않았다. 아래 모든 결과는 tester가 이번에 직접 실행했다.

## 환경

- 실행 날짜 2026-10-03, Linux amd64. Go 1.27.1, Node 22.22.2, npm 10.9.7, Python 3.12.3, GNU Make 4.3, Docker 29.4.3, Compose 5.1.3.
- FullOps 플러그인 0.9.10. FullOps lint 기준 ref `729446d8da57`.
- clone은 레포 외부 임시 경로에 두었다. 명령은 `run.py`(셸 미사용)로 실행해 실제 종료코드를 기록했다. `| tail`을 쓰지 않았다.
- 실제 키·token·payload를 쓰지 않았다. 설정은 `.env.example`의 더미 값이다.

## 수락 기준별 결과

| ID | 결과 | 시나리오 | 증거와 종료코드 |
|---|---|---|---|
| SETUP-01 | 통과 | N-01, N-02 | `make install` 0, `make build` 0 ([main.log](SAR-SETUP-001-TESTER-test/main.log)). 설치·빌드·lint 재실행 뒤 `git status --porcelain` 빈 값. 무시 대상은 dist·node_modules·build·`__pycache__` |
| SETUP-02 | 통과 | N-04, C-04 | relay·migrate는 `DATABASE_URL` 없이 종료코드 1. adapter 실행 0. relay 기동 중 migrate는 `no SQL migrations; no-op`. relay 코드에 goose·migration 호출 없음 ([boundary.log](SAR-SETUP-001-TESTER-test/boundary.log), [live.log](SAR-SETUP-001-TESTER-test/live.log)) |
| SETUP-03 | 통과 | C-03, C-04 | 기본 profile 서비스 postgres·migrate·relay. `--profile tunnel`에서 cloudflared 추가로 네 서비스. Postgres `ports` 없음·internal 네트워크, migrate `restart: no`, relay는 migrate 완료 후 시작, 게시 `127.0.0.1:8080`. 실제 기동 컨테이너 포트도 Postgres `5432/tcp`(미게시) |
| SETUP-04 | 통과 | C-01, C-02, B-02 | `POSTGRES_PASSWORD` 누락 시 Compose 종료코드 1. 잘못된 URL relay 1. `.env`는 미추적·무시. 비밀값 패턴 일치는 더미 `example-local-only`뿐. cloudflared token 없음 255 ([runtime.log](SAR-SETUP-001-TESTER-test/runtime.log)) |
| LINT-01 | 통과 | N-03 | `make lint` 0, `make test` 0 (race, 3개 패키지) |
| LINT-02 | 통과 | F-01–F-05 | Go 포맷 2(`Go formatting violations`), TS 미사용 변수 2(`@typescript-eslint/no-unused-vars`), 타입 오류 2(`TS2322`), TS 포맷 2. 모든 주입은 원복 뒤 Git 상태 빈 값, `make lint` 0 ([injection.log](SAR-SETUP-001-TESTER-test/injection.log)) |
| LINT-03 | 통과 | F-06, F-07, N-02 | `make verify` 0 ([verify.log](SAR-SETUP-001-TESTER-test/verify.log)). 등록 명령 전파: `0cc10b0` 기준 `lint.py`가 TS 타입 오류에서 종료코드 1, `[failed] product-lint`, ERROR 1 ([fullops-injected.json](SAR-SETUP-001-TESTER-test/fullops-injected.json)). 검사 모드는 자동 수정 안 함. 정상 재실행에서 추적 파일 변경 없음 |
| DOC-01 | 통과 | D-01 | README·project.md·D03의 Go 1.27.1·Node 22.22.2·npm 10.9.7·Postgres 17·명령 목록이 `go.mod`·`.nvmrc`·`Makefile`·`compose.yaml`과 일치. 미구현·미적용 항목을 명시 |
| SCOPE-01 | 통과 | B-01–B-03 | 실제 relay: `GET /healthz` 200, `POST /healthz` 405, `POST /v1/messages`·`/webhook`·`/approve` 404. adapter `webhook:false`, `evidenceFetch:false`. 원천 diff 없음. `lint.json` 변경은 `commands` 추가뿐 ([scope.log](SAR-SETUP-001-TESTER-test/scope.log)) |
| QA-01 | 통과 | 전체 | 이 보고서 |

## 실행한 명령과 종료코드 요약

| 명령 | 종료코드 |
|---|---|
| `make install` / `make lint` / `make test` / `make build` | 0 / 0 / 0 / 0 |
| `make lint` 재실행 | 0 |
| 주입 4종의 `make lint` | 2 / 2 / 2 / 2 |
| 원복 후 `make lint` ×4 | 0 ×4 |
| `make verify` | 0 |
| `make verify-runtime` | 0 (`PASS: local startup, empty SQL no-op, readiness, missing settings, DB failure, disabled Tunnel, inert adapter`) |
| `lint.py --from 0cc10b0` (타입 오류 커밋) | 1 |
| `docker compose config` (필수값 없음) | 1 |
| `docker compose up --build --wait relay` / `down -v` (고유 프로젝트) | 0 / 0 |

## FullOps lint와 직접 제품 검사의 구분

- 직접 제품 검사: 위 표의 `make lint`·`make test`·`make verify`·`make verify-runtime`이다. 대상은 dev SHA `0cc10b0`이다.
- 기준 `729446d8da57`의 `lint.json`에는 `commands`가 비어 있다. 이 기준의 FullOps lint는 `product-lint`를 실행하지 않는다. 따라서 FullOps 통과는 제품 lint 실행의 증거가 아니다.
- 등록 명령의 실행과 실패 전파는 기준을 `0cc10b0`로 바꾼 임시 복사본에서 따로 증명했다(LINT-03). 이 복사본은 제품 저장소에 반영하지 않았다.
- 원천 스냅샷의 DOC-003 차단: 기준 `729446d8da57`에 원천 전용 exclude가 있다. 원천 파일은 변경하지 않았다.
- tester 최종 SHA의 FullOps lint 결과는 아래 "tester 최종 검사" 절에 기록한다.

## 결함과 관찰

결함 없음. 아래는 수락에 영향이 없는 관찰이다.

1. `docker compose config --services`의 기본 결과는 세 서비스다. cloudflared는 `tunnel` profile에만 있다. D02의 "네 서비스 정의"와 일치하며 기본 실행에서 Tunnel이 꺼져 있다는 설계와도 일치한다.
2. `make lint`의 TypeScript 검사 대상은 `adapters/src`와 설정 파일이다. Python 스크립트는 `py_compile` 구문 검사만 한다. 스크립트 로직은 `make verify`·`make verify-runtime` 실행으로 확인된다.
3. 실행 전에 로컬 Docker에 postgres·cloudflared 이미지가 이미 있었다. 빈 이미지 캐시에서의 첫 pull은 검증하지 않았다.

## 미실행과 한계

- 업무 SQL이 없어 goose Up 적용과 `sqlc generate`는 실행하지 않았다. 빈 migration no-op은 적용 성공이 아니다.
- UI가 없어 직접 시각 검수와 캡처·영상은 해당 없다.
- 실제 Tunnel·운영 배포·네 벤더 어댑터·C1–C5 보안 동작은 범위 밖이며 검증하지 않았다.
- ARM 등 다른 플랫폼과 다른 Docker 버전은 실행하지 않았다.
- 문서 메타데이터 검사는 의미적 수락을 대신하지 않는다.
