---
title: SAR-SETUP-001-TESTER — 초기 구성과 제품 lint 독립 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-SETUP-001-TESTER]
summary: SAR-SETUP-001 초기 구성과 제품 lint의 독립 QA 시나리오와 통과 조건
---

# SAR-SETUP-001-TESTER — 초기 구성과 제품 lint 독립 QA 시나리오

대상은 dev 완료 SHA `0cc10b083771be9b3423833b222c57d426315333`이다. 제품에 화면이 없어 Jev 웹·Unity 조작 테스트를 쓰지 않는다. 모든 시나리오는 셸 명령과 종료코드로 판정한다.
요구사항은 D02 `SAR-SETUP-001`이다. 결과는 [QA 보고서](../qa-reports/SAR-SETUP-001-TESTER.md)와 `../qa-reports/SAR-SETUP-001-TESTER-test/`에 있다.

## 공통 조건

1. 실행 위치는 레포 밖의 임시 Git clone이다. clone을 `0cc10b0` 에 detached checkout하고 `git status --porcelain`이 비어 있음을 확인한다.
2. 실제 키·Tunnel token·운영 DB를 쓰지 않는다. 설정은 `.env.example`의 더미 값만 쓴다.
3. 명령은 `run.py`로 실행한다. `run.py`는 셸을 거치지 않고 명령 자신의 종료코드를 로그에 기록한다. 파이프로 종료코드를 가리지 않는다.
4. 위반 주입은 clone의 복사본에서만 한다. 제품 파일을 tester 체크아웃에서 수정하지 않는다.

## 시나리오

| ID | 덮는 기준 | 단계 | 통과 조건 |
|---|---|---|---|
| N-01 정상 설치·빌드 | SETUP-01 | `make install`, `make build`를 순서대로 실행한다. 도구 버전은 `go version`, `node -v`, `npm -v`로 기록한다. | 종료코드 0. 버전은 Go 1.27.1, Node 22.22.2, npm 10.9.7이다. |
| N-02 재실행 안정성 | SETUP-01, LINT-03 | N-01 뒤 `git status --porcelain`을 실행한다. `make lint`와 `git status --porcelain --ignored`를 다시 실행한다. | 추적 파일 변경 없음. 무시 대상은 `adapters/dist/`, `adapters/node_modules/`, `build/`, `scripts/__pycache__/`뿐이다. |
| N-03 제품 lint | LINT-01 | `make lint`, `make test`를 실행한다. | 종료코드 0. Go 포맷·vet·module 검증, TypeScript 포맷·ESLint·타입, YAML·Compose·스크립트 구문을 검사한다. |
| N-04 최소 실행 | SETUP-02 | `./build/relay`와 `./build/migrate`를 `DATABASE_URL` 없이 실행한다. `npm run start --prefix adapters`를 실행한다. | Go 프로세스는 종료코드 1과 `DATABASE_URL is required`를 출력한다. adapter는 `unimplemented` JSON을 출력하고 종료코드 0이다. |
| F-01 Go 포맷 위반 | LINT-02 | 복사본의 `internal/config/config.go`에 `func   badFormat( ){ }`를 추가하고 `make lint`를 실행한다. | 종료코드 2와 `Go formatting violations`. |
| F-02 TypeScript lint 위반 | LINT-02 | `adapters/src/index.ts`에 미사용 변수를 추가하고 `make lint`를 실행한다. | 종료코드 2와 `no-unused-vars` 진단. |
| F-03 TypeScript 타입 오류 | LINT-02 | `const typeMismatch: number = "text";`를 추가하고 `make lint`를 실행한다. | 종료코드 2와 `TS2322`. |
| F-04 TypeScript 포맷 위반 | LINT-02 | 따옴표·공백이 틀린 `console.info(   'x'  )`를 추가하고 `make lint`를 실행한다. | 종료코드 2와 Prettier 진단. |
| F-05 원복·재통과 | LINT-02, 복구 | 각 주입 뒤 `git checkout -- .`를 실행하고 `git status --porcelain`, `make lint`를 실행한다. | Git 상태 비어 있음. `make lint` 종료코드 0. |
| F-06 dev 검증 스크립트 | LINT-02, LINT-03 | clone에서 `make verify`를 실행한다. | 종료코드 0. 위반별 `make lint` 종료코드 2, 원복 후 0. |
| F-07 FullOps 등록 명령 전파 | LINT-03 | 복사본에 TypeScript 타입 오류를 커밋한다. `lint.py --repo <복사본> --from 0cc10b0`을 실행한다. 기준 `0cc10b0`의 `lint.json`에 `product-lint`가 있다. | 종료코드 1과 `[failed] product-lint`, ERROR 1. 실행 뒤 복사본을 `0cc10b0`로 되돌린다. |
| C-01 필수값 누락 | SETUP-04 | `POSTGRES_PASSWORD`·`DATABASE_URL` 없이 `docker compose config`를 실행한다. `DATABASE_URL=not-a-url`로 relay를 실행한다. | Compose는 `POSTGRES_PASSWORD is required`로 종료코드 1. relay는 URL 오류와 종료코드 1. |
| C-02 비밀값 비추적 | SETUP-04 | `git ls-files --error-unmatch .env`, `git check-ignore -v .env .env.local`, 비밀값 패턴 `git grep`을 실행한다. | `.env`는 미추적이고 무시된다. 패턴 일치는 더미 `example-local-only`뿐이다. |
| C-03 Compose 구성 | SETUP-03 | `docker compose --env-file .env.example config --format json`을 구조로 확인한다. 기본 profile과 `--profile tunnel`의 서비스 목록을 비교한다. | 기본은 postgres·migrate·relay. `tunnel` profile에 cloudflared. Postgres `ports` 없음, `database` 네트워크 internal. migrate `restart: no`, relay는 `service_completed_successfully` 의존. relay 게시는 `127.0.0.1`. |
| C-04 실제 기동·migration 경계 | SETUP-02, SETUP-03 | 고유 프로젝트명·포트 18080으로 `up --build --wait relay`를 실행한다. migrate 로그, public 테이블 수, 컨테이너 포트를 확인한다. 실행 뒤 `down -v`로 이 프로젝트만 정리한다. | relay 기동 성공. migrate는 `no SQL migrations; no-op`. public 테이블 0. Postgres 호스트 게시 없음. |
| B-01 동작 경계 | SCOPE-01 | 실행 중 relay에 `GET /healthz`, `POST /healthz`, `POST /v1/messages`, `POST /webhook`, `GET /approve`를 보낸다. | 200, 405, 404, 404, 404. 업무·승인·webhook 응답 없음. |
| B-02 외부 기능 비활성 | SCOPE-01, SETUP-04 | `go test -race -count=1 -v`를 실행한다. `adapters/src/index.ts`와 코드의 webhook·evidence 항목을 검색한다. `make verify-runtime`을 실행한다. | 테스트 통과. adapter의 `webhook`·`evidenceFetch`는 false. cloudflared는 token 없이 종료코드 255. |
| B-03 원천·범위 diff | SCOPE-01 | `git diff --quiet 729446d 0cc10b0 -- .fullops-squad/docs/planning/sources`를 실행한다. 변경 파일 목록과 `lint.json` diff를 확인한다. | 원천 diff 없음. `lint.json` 변경은 `commands` 추가뿐. exclude·size·rules 변경 없음. |
| D-01 문서 일치 | DOC-01 | README·project.md·D03의 버전·경로·명령을 실제 `go.mod`·`Makefile`·`compose.yaml`과 대조한다. | 일치. 미구현 항목은 미구현으로 표시. |

## 시나리오 밖 항목

- 업무 SQL·goose Up 적용·sqlc 생성: 파일이 없어 미적용이다. 후속 기능 과제의 검사 범위다.
- 직접 시각 검수: UI가 없어 해당 없다. designer의 후속 조건은 UI 구현 뒤 시작한다.
- 실제 Tunnel 연결·운영 배포·네 벤더 어댑터: 범위 밖이다. 검증하지 않았다.
- 기준 `729446d8da57`의 FullOps lint: tester 최종 커밋에서 실행하며 QA 보고서에 별도로 기록한다. 이 결과는 제품 동작의 증거가 아니다.
