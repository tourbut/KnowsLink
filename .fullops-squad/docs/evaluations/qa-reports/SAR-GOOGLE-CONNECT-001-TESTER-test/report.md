---
title: SAR-GOOGLE-CONNECT-001-TESTER 좁은 독립 QA
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-GOOGLE-CONNECT-001-TESTER]
summary: 운영 서버 별도 공간과 Windows 로컬에서 Google 클라이언트 연결의 핵심 보안 경계를 짧게 독립 검증한 결과
---

# SAR-GOOGLE-CONNECT-001-TESTER 좁은 독립 QA

테스트 레벨은 lite다. 대상 SHA는 `152217f63cced85f620695961955d680ce27864e`(기준 `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`)다. 소스는 이 SHA의 `git archive`(sha256 `c515ccc0999b93de15af953a9e73c352498010fb193e6d66b279f72803179dd7`)다. 검사는 제품 코드를 바꾸지 않았다. 추가한 파일은 tester 테스트 코드 하나뿐이며 서버 임시 복사본에만 넣었다.

## 접근과 보존

- 서버 접속은 이 체크아웃의 `.env.server`(main 정본 링크)를 비공개로 읽는 SSH(Paramiko 5.0.0)로 했다. 이 Windows의 Python에는 Paramiko가 없어 스크래치 디렉터리의 `--target`에 설치했다(전역 환경 변경 없음). IP·USER·PW는 출력·로그·Git·셸 인자에 넣지 않았고 출력은 값 치환으로 걸렀다. 처음 보는 서버 키는 이 세션의 비공개 known_hosts에 고정했다.
- 별도 공간은 서버의 `/tmp/kl-gc-tester560836`(권한 700)이다. 격리 DB는 `postgres:17` 컨테이너 3개(`kl-gc-tester560836`, `-2`, `-3`)이며 loopback 임의 포트·256MiB·1CPU·합성 계정이었다. 각 컨테이너는 EXIT trap과 라벨 `knowslink.tester.qa=tester560836`로 지웠다. 마지막 정리에서 라벨 컨테이너 없음·임시 디렉터리 삭제를 확인했다. 실행 중이던 `knowslink-relay-1`·`knowslink-cloudflared-1`·`knowslink-postgres-1`과 `myportfolio-*` 컨테이너 목록은 검사 전후 같다. 운영 DB·키·Tunnel·다른 서비스는 읽거나 바꾸지 않았다.
- 로컬 Docker는 기동하지 않았다. 로컬 실행은 Docker 없는 Node/Python/Go 컴파일뿐이다.
- Context7 대신 공식 문서를 쓸 라이브러리 API 확인은 필요하지 않았다. 새 외부 API 호출 코드를 쓰지 않았다(테스트는 기존 `oauth2`·`pgx` 사용 패턴을 DEV 테스트에서 그대로 가져왔다).

## 실행한 검사와 종료코드

| 구분 | 명령(요약) | 환경 | 종료코드 | 증거 |
|---|---|---|---|---|
| 직접 독립 | `npm ci --prefix adapters`, `npm run build --prefix adapters` | 서버 Node 22.22.2/npm 10.9.7 | 0, 0 | `server-exits.txt` |
| 직접 독립 | `node adapters/dist/login.test.js` | 서버 | 0 | `server-ts-login.log` |
| 직접 독립 | `node adapters/dist/connect.test.js` | 서버 | 0 | `server-ts-connect.log` |
| 직접 독립 | `go test -race ./internal/relay -run TestDevice -count=1` | 서버 Go 1.27.1 | 0 | `server-go-unit.log` |
| 직접 독립 | `go run ./cmd/migrate` | 서버, 격리 DB | 0 | `server-exits.txt` |
| 직접 독립(DEV 테스트 코드 재실행) | `go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP\|TestGoogleHTTP\|TestEmailIdentity' -count=1 -v` | 서버, 격리 DB | 0 | `server-dev-integration.log` (세 최상위 테스트 PASS) |
| 직접 독립(tester 신규) | `go test -race -tags integration ./internal/relay -run TestTesterGoogleDeviceBoundaries -count=1 -v` | 서버, 새 격리 DB | 0 | `server-tester-integration.log`, 테스트 코드 `device_tester_qa_integration_test.go`(sha256 `cbd443116cc0b9c72297f8b4c73ee98a82a019fa191944d426d42112e52acbf0`) |
| 직접 독립 | `node adapters/dist/connect.test.js`, `login.test.js` | Windows 로컬 Node 22.22.2(`KnowsLinkDevTools`), 스크래치 복사본 | 0, 0 | `local-win-ts-connect.log`, `local-win-ts-login.log` |
| 직접 독립 | `python scripts/check_public_ingress.py` | Windows 로컬 | 0 | 출력: `Public ingress: Google/member paths only; ...` |
| 직접 독립(컴파일) | `gofmt -l`, `go vet -tags integration ./internal/relay` | Windows 로컬 | 0 | 출력 없음 |
| Mutation 2건 | 아래 참고 | 서버 복사본 | 둘 다 1(테스트 실패, 기대됨) | `server-mutation.txt` |

`server-exits.txt`에는 `tester-integration=0`이 세 줄 있다. 첫 줄은 무효다. 처음 업로드한 tar에 tester 테스트 파일이 빠져 `no tests to run`으로 exit 0이 나왔다. 파일을 따로 올려 재실행한 둘째 줄이 첫 유효 결과다. 셋째 줄은 테스트를 보강(기존 회원 경쟁 시나리오 추가)한 최종판이다. 최종 로그는 마지막 실행(`--- PASS: TestTesterGoogleDeviceBoundaries`)이다. 첫 줄과 둘째 줄을 PASS 근거로 쓰지 않는다.

DEV 증거 재사용과 직접 실행의 구분: DEV의 `make lint`·`make test`와 서버 lint/test 로그는 재실행하지 않았다. 최종 lint JSON은 리뷰의 `lint.json`으로 대조했다. 위 표의 모든 행은 이번에 직접 실행했다. 테스트 소스 중 `TestGoogleDeviceHTTP` 등은 DEV가 작성한 코드이며 독립 실행은 같은 SHA의 코드를 새 환경에서 돌린 것이다.

## tester 신규 검사가 덮는 경계

`TestTesterGoogleDeviceBoundaries`는 실제 `Service.Handler()`에 서명한 로컬 OAuth 서버(PKCE 검증)와 Postgres를 붙이고 브라우저 쿠키 항아리 셋과 익명 브라우저로 시나리오 8개를 실행한다.

1. S1 선점·경쟁: 기존 로그인 회원 B와 신규 계정 D가 미결합 요청에 Google을 시작하고 A가 먼저 callback한다. B와 D의 callback은 403이고 B의 세션은 그대로이며 D는 세션이 없고 회원 수가 늘지 않으며 결합은 A에 머문다.
2. S2 타 회원: 익명·로그인한 B가 보기, 재시작, 취소, 승인, 기존 `/home/confirm`, `/home/connections/{id}`를 시도해도 수락되지 않고 요청 상태가 변하지 않는다.
3. S3 token API 오용: 승인 전 `/v1/connect/info`·`prepare`·`complete`는 401. 위조 서명·다른 client의 poll은 401. 올바른 poll은 `prepared`만 보이고 agent·credential이 새지 않는다.
4. S4 취소: 소유자 취소 뒤 승인과 poll이 거부된다.
5. S5 만료: 페이지·재시작·poll이 거부된다.
6. S6 잘못된 connection 값(URL·상대경로·임의 문자열)은 303을 만들지 않는다.
7. S7 cross-site 브라우저 POST(`/auth/google`·`/home/device-confirm`)는 403.
8. S8 거부된 경로가 agent·pair를 만들지 않는다.

성공 경로(시작, Google 결합, 지문 표시, 명시적 동의, 서명 poll, complete, 단일 사용, 같은 회원의 두 별도 agent·credential, 최근 인증 만료 422, cross-site 동의 403, 기존 `/home/confirm` 우회 403)는 DEV의 `TestGoogleDeviceHTTP`를 서버에서 직접 재실행해 확인했다.

### Mutation 확인

서버 복사본의 소스를 바꿔 tester 검사가 실제로 잡는지 확인했다. 복사본은 정리했다.

- poll의 `deviceProof` 검사 제거: `TestTesterGoogleDeviceBoundaries`가 S3에서 `got 200 want 401`로 실패했다.
- `finishGoogleLogin`의 `m.Owner != c.Owner` 거부 조건 제거: 첫 시도(조건을 `false`로)는 컴파일 오류였고, 변수를 쓰게 고친 두 번째에서 보강 전 테스트는 통과했다. 이것이 리뷰 finding 3이다. 보강한 테스트는 `existing competing member got 200 session-kept=false`로 실패했다. DEV 검사는 같은 변경에서 통과한다.

## 결과

- 실패 0. 직접 독립 검사는 모두 통과했다.
- 제품 결함 재현은 없다. 리뷰 finding 2(2000개 상한)는 코드와 한도 숫자를 읽어 도출했고 부하로 재현하지 않았다.
- 수락 판정: 이 lite QA의 통과 조건(핵심 연결 성공, 명시적 동의, 만료·replay, 다른 회원·클라이언트 분리, 필요한 기존 로그인 회귀)은 충족했다.

## 생략과 한계

- 실제 Google 계정·실제 OAuth 서버·실제 Cloudflare Access/Tunnel·외부 Bot 성공은 실행하지 않았다. 합성 OAuth(서명한 로컬 RS256 제공자)와 합성 DB의 결과를 그 성공이라고 부르지 않는다.
- 시각 검수와 캡처는 하지 않았다. 디자인 판정은 designer 담당이다.
- 부하·장시간·동시성 강도는 검증하지 않았다.
- 실제 Windows 사용자 환경의 `agent.json` ACL 외 경로(다른 사용자 계정 시나리오)는 DEV의 `icacls` 기반 테스트를 재실행한 것으로 갈음했다.
- 서버 clean clone의 등록 전체 lint는 반복하지 않았다(lint 증거 검증으로 대체).

## 재현

`device_tester_qa_integration_test.go`를 SHA `152217f`의 `internal/relay/`에 복사한다. 격리 Postgres 17을 띄우고 `TEST_SYNTHETIC_DATABASE=1`·`TEST_DATABASE_URL`을 설정해 `go run ./cmd/migrate` 뒤 `go test -race -tags integration ./internal/relay -run TestTesterGoogleDeviceBoundaries -count=1 -v`를 실행한다. 운영 `DATABASE_URL`을 쓰지 않는다.

## tester 체크아웃 lint.py 시도(done-gate)

HEAD `165ec4c`(이 체크아웃은 제품 변경 없이 기록만 추가)에서 `lint.py --repo . --from f9f7675…`를 실행했다. 첫 실행은 `make`·Node 미탐색으로 실행 불가였다. 둘째 실행은 로컬 도구(Node 22.22.2, GNU Make 4.4.1)를 PATH에 넣었다. 첫 번째 ERROR(DOC-003)는 report.md front matter 형식이었고 `deliverables.py --stamp`로 고쳤다. 남은 ERROR 2건은 `product-lint`(gofmt가 `internal/relay` 전체 파일을 나열)와 `product-test`(`TestFrozenParsingAndSigning`의 `registry integrity failure`)다. 원인은 이 체크아웃의 `core.autocrlf=true`로 Go 파일이 CRLF 작업 사본이라는 점이다(`git ls-files --eol`: `i/lf w/crlf`). 이 체크아웃의 HEAD에는 DEV가 추가한 `.gitattributes`(eol=lf)가 아직 없다. 이 두 실패는 변경하지 않은 제품 파일의 환경 문제이며 내 변경과 무관하다. 검토 SHA의 등록 lint/test 증거는 위 리뷰의 서버 clean clone `lint.json`(exit 0)이다. 이 체크아웃에서 통과 기록을 만들지 못했으므로 done-gate 통과로 주장하지 않는다. 필요하면 main 통합 뒤 이 체크아웃을 최신 main으로 동기화하여 재실행한다.
