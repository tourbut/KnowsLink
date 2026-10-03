---
title: SAR-SETUP-001-DEV-099 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV-REVIEW]
summary: 고정 구현 SHA의 독립 코드 리뷰와 조건부 수락 근거를 기록한다
---

# SAR-SETUP-001-DEV-099 리뷰

- 검토자 / CLI / 모델: Orca dispatch `ctx_53a98f9ed8f8`(task `task_158decad4fc0`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 ID는 `ba47d7dc-4c3e-472c-bdaa-3afa95e7a285`이다.
- 구현자: dispatch `ctx_67f98ed4cd42`(task `task_a7fd5d1b8806`) / Codex `gpt-6.1-sol` high. 구현 세션 ID는 `01a0ffa6-bcbc-7383-a8e5-f14521f0dfa3`이다.
- base SHA / head SHA / merge-base: `dbe0b40076af4d440bb263ca4d02d671780d2514` / `0cc10b083771be9b3423833b222c57d426315333` / `dbe0b40076af4d440bb263ca4d02d671780d2514`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c) / `review/rule.json`(rule_sha256 `ab2116fb…2aa0`) 7개 그룹, 공통 규칙 `fullops-common-0.3.2`(구현 당시 0.3.1).
- 요구사항·완료 기준 원천: snapshot의 D02 `docs/planning/product-specs/SAR-SETUP-001.md`, `handovers/logs/2026-10-03_to_dev.md`, D03 두 문서, `docs/exec-plans/phases/SAR-SETUP-001-DEV.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 50 / 38 / 12 / 50 / 0.
- lint(`lint.json`, base `dbe0b40`) ERROR / WARNING / 실행 불가와 사유: 0 / 2 / 0. 종료코드 0. 지시서 기준 `729446d8da57`의 결과(`lint-729446d8da57.json`)도 0 / 2 / 0, 종료코드 0이다.

## 독립성 근거

구현 세션은 Orca `worker-show --dispatch ctx_67f98ed4cd42`에서 provider codex와 worktree `fullops-dev`를 확인했다. 그 다음 `~/.codex/sessions/2026/10/03/rollout-2026-10-03T11-45-12-01a0ffa6-….jsonl`의 `session_meta`를 직접 읽었다. 이 파일의 시각은 dispatch 시각 02:45:11Z와 같다. cwd는 `fullops-dev`이고 지시문의 task ID는 `task_a7fd5d1b8806`이다.

리뷰 세션은 Claude Code의 현재 세션이다. provider와 세션이 모두 구현자와 다르다.

snapshot `/tmp/SAR-SETUP-001-review-0cc10b0`은 HEAD `0cc10b0…`, detached 상태다. 작업 시작과 lint·go test 실행 뒤 `git status --porcelain`이 빈 값임을 확인했다. snapshot을 변경하지 않았다. 결과는 기록 체크아웃 `fullops-coor`에만 썼다.

## 검토 범위

result.json의 50개 `(path, status)`를 모두 reviewed로 기록했다. 파일마다 확인 근거가 있다.

- Go 6개: 전체 파일을 읽었다. config, pgx, goose 호출 경로와 오류·context·timeout·종료 처리를 검토했다.
- Compose, Dockerfile, Makefile, `.dockerignore`, `.gitignore`: 서비스 경계를 검토했다. Postgres 포트 미게시와 internal 네트워크를 확인했다. one-shot migrate와 의존 순서, loopback 게시를 확인했다. cloudflared는 profile에서만 실행되고 token 기본값이 비어 있다. non-root 런타임과 build context 제한도 확인했다.
- TypeScript adapter와 설정 5개: 벤더 연결이 없다. strict와 noEmitOnError를 확인했다. check는 `&&` 체인으로 실패를 전파한다.
- 스크립트 3개: 임시 복제본, 고유 Compose project, tmpfs override, `finally down`, 네트워크 차단 cloudflared를 확인했다. 실제 자원은 변경되지 않는다.
- 문서 12개: 경로, 버전, 명령, 상태를 실제 코드·설정과 대조했다. `handovers/to_dev.md` 153줄 삭제는 `handovers/logs/2026-10-03_to_dev.md`에 원문이 보존된 인박스 정리다.
- OCR 제외 12개: `.env.example`, `.npmrc`, `go.mod`는 직접 읽었다. 비밀값이 없다. `go.sum`은 `go mod verify` 로그와 리뷰 세션의 `go vet`·`go test` 성공으로 무결성을 확인했다. 원시 로그 8개는 command·exit 줄, 기대 진단, 비밀값 패턴을 검색했다. 주장된 종료코드가 로그와 일치한다. install의 audit 0건, clean clone의 `git status --porcelain: ''`, runtime의 public 테이블 0, 503, cloudflared 255를 직접 확인했다.

## 발견 사항

critical, high, medium 발견은 없다.

1. low, 미해결. `cmd/relay/main.go:36-46`, `cmd/migrate/main.go:43-63`.
   - 재현: 잘못된 비밀번호나 연결 거부가 일어난다.
   - 증상: 둘 다 `database ping failed`만 출력한다.
   - 영향: 원본 오류를 버리므로 비밀값 노출은 막는다. 하지만 운영자가 인증·DNS·연결 거부를 구분할 수 없다. 개발 골격 범위의 수락에는 영향이 없다.
   - 권고: 후속 과제에서 DSN을 빼고 SQLSTATE나 오류 분류만 기록하는 방안을 검토한다.
2. low, 미해결. `scripts/check_compose.py:22-35`.
   - 재현: `PYTHONOPTIMIZE=1 make lint-config`를 실행한다.
   - 증상: assert가 제거되어 Compose 경계 검사 없이 통과한다.
   - 영향: 현재 Makefile과 기본 환경에서는 영향이 없다. 실패 은닉 경로만 남는다.
   - 권고: 명시적인 `if … raise SystemExit`로 바꾼다.

### 주요 보안 경계 확인

- DB·설정 오류 메시지에 URL·비밀번호가 없다. `config_test.go`와 `main_test.go`가 비노출을 검사한다.
- Postgres는 `database`(internal) 네트워크에만 있고 호스트 포트가 없다.
- relay는 `127.0.0.1`에만 게시한다. 업무 경로 `/ingest`, `/approve`, `/exec`는 404다.
- cloudflared는 `tunnel` profile이다. token이 없으면 종료코드 255로 실패한다. 검증은 `--network none`에서 실행했다.
- migrate는 전역 Go registry를 비활성화하고 `.go` 파일을 거부한다. 빈 디렉터리는 no-op이며 goose 테이블도 만들지 않는다.

### lint 실패 전파 확인

`product-lint`는 `make lint`다. Makefile 각 줄의 실패가 make를 중단한다. `lint-go`는 `set -eu`와 `gofmt -l` 결과로 실패한다. adapters `check`는 `&&` 체인이다. `lint-injection.log`에서 결과를 확인했다. 정상 실행은 0이다. Go 포맷, ESLint, TS2322, Prettier 주입은 각각 2와 기대 진단을 남겼다. 원복 뒤 0이다. 자동 수정 모드는 없다.

## 검증 및 남은 제약

- FullOps lint (지시서 기준): 설치된 0.9.10 `lint.py --repo <snapshot> --from 729446d8da57` 결과다. 파일 52개를 검사했다. 종료코드 0, ERROR 0, WARNING 2(LINT-001, LINT-000), 실행 불가 0이다. 결과는 [lint-729446d8da57.json](lint-729446d8da57.json)에 있다.
- FullOps lint (check 기준): `review.py check`는 lint의 base가 리뷰 base와 같아야 통과한다. 그래서 `--from dbe0b40…`로 다시 실행했다. 파일 46개를 검사했다. 결과는 같다: 종료코드 0, ERROR 0, WARNING 2, 실행 불가 0. 결과는 [lint.json](lint.json)에 있다. `729446d`는 `dbe0b40`의 조상이다. 두 기준 모두 merge-base의 `commands`가 비어 있다.
- `review.py check --key SAR-SETUP-001-DEV-099 --from dbe0b40… --to 0cc10b0…`: 종료코드 0이다. 출력은 `기록 검사 통과: reviewed=50, skipped=0, total=50, lint WARNING=2`이다.
- WARNING은 기준 ref의 `commands`가 비어 있어서 생긴다. 새 `product-lint`는 이 검사에서 실행되지 않는다. 대신 dev의 직접 실행 로그가 증거다. 병합 뒤부터 FullOps lint가 `make lint`를 실행한다.
- 초기 커밋 `bd9c8dc`의 SEC-001 ERROR 2건(로그 표시 형식)은 head에서 해소되었다.
- 리뷰 세션 보조 확인: snapshot에서 `gofmt -l cmd internal`(빈 출력, 0), `go vet ./cmd/... ./internal/...`(0), `go test -count=1 ./cmd/... ./internal/...`(0)을 실행했다. snapshot은 변경되지 않았다.
- dev 증거는 원래 조건의 기록이다. 코드 SHA는 `929832a`, 환경은 Go 1.27.1, Node 22.22.2, Docker 29.4.3이다. `929832a..0cc10b0` 사이에는 문서와 로그만 바뀌었고 제품 코드는 같다. 이 기록을 QA의 독립 실행으로 표시하지 않는다.
- 실행하지 않은 검증: `make lint`, `make verify`, `make verify-runtime`. 이 명령들은 snapshot에 `npm ci`와 Compose 빌드가 필요해 snapshot을 바꾼다. 새 실패나 증거 결함도 없어 재실행하지 않았다. 독립 재현은 SAR-SETUP-001-TESTER(QA-01)가 맡는다.
- 미적용: sqlc generate와 실제 SQL Up. 업무 SQL이 없으며 D02가 허용한다.
- 공통 규칙 0.3.2 영향: 0.3.2는 testing.md에 `담당과 반복 범위` 절만 추가했다. 이 절은 직접 검사와 독립 QA의 구분, 증거 재사용 시 원래 SHA 유지를 요구한다. 구현 기록은 이미 이를 지키며 직접 검증을 독립 QA로 표시하지 않는다. 재검토 범위는 늘지 않는다. `project.md`와 D03의 `0.3.1` 기재는 구현 당시 기준으로 유지해도 된다.

## 검토 결론

조건부로 수락한다. critical·high 발견이 없다. low 2건은 비차단이다. 50개 파일을 모두 검토했다. FullOps lint ERROR는 0이다.

남은 수락 조건은 두 가지다.

1. SAR-SETUP-001-TESTER의 독립 QA(QA-01)가 head `0cc10b0`을 통과해야 한다.
2. 병합 책임자가 최신 SHA와 허가된 병합 범위를 확인해야 한다.

이 결론은 AI 검토자의 판단이다. OCR의 자동 판정이 아니다. `review.py check` 통과는 기록 형식만 검사한다. 내용의 정확성이나 테스트 성공을 보증하지 않는다.
