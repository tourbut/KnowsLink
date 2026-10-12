---
title: 사람 peer DEV 실행 환경 준비 기록
status: review
updated: 2026-10-12
owner: dev
tasks: [SAR-HUMAN-PEER-001-DEV]
summary: 합성 인증·전달과 21개 화면 증거, f567abf 실패 보존 및 7f83fd6 전체 SQL 회귀 종료코드 0
---

# 실행 환경 준비

이 기록은 DEV 자체 검사다. 독립 QA·코드 리뷰·디자인 수락을 주장하지 않는다.
이 child의 task는 `task_8988765c2d4f`, dispatch는 `ctx_9ded146db1ea`다.
준비 SHA는 `59d419d0b7c4fca46cabafc04b5a47c6233c944a`다.
구현 후보 `4eaefa6e0a5cf623f0d8d8080d918086f86030ec`에서 focused SQL integration과 실제 인증·API·전달 자체 검사가 통과했다. 세 CSS 폭의 정상 viewport 캡처 21개를 보존했다. `f567abff18b8f6ce77fc79e6601ba7163ae7f9ff` 전체 tagged SQL 회귀의 실패를 보존한 뒤, 수정 후보 `7f83fd6f3ae9945e292a86212c0ab566ea282b7a`의 전체 회귀가 118.664초 후 종료코드 0으로 통과했다. 부모의 최종 제품 게이트와 독립 QA·리뷰·디자인 판정은 별도다.

## 실제 환경

- Windows: Go 1.27.1, 기본 Node 24.13.0/npm 11.8.0. 기존 별도 Node 22.22.2를 확인했다.
- Windows PATH에는 make·psql이 없다. 기존 w64devkit에는 GNU Make 4.4.1이 있다. 해당 조합의 제품 게이트는 실행하지 않았다.
- SSH Linux: 모듈 밖 기본 Go 1.22.2, 모듈 안 자동 선택 Go 1.27.1, Node 22.22.2/npm 10.9.7, GNU Make 4.3, Python 3.12.3.
- Docker 29.4.3, Compose 5.1.3. `compose.yaml`의 고정 Postgres 17 이미지로 실제 서버 17.11을 확인했다.
- Linux 호스트에도 psql은 없다. 격리 Postgres 컨테이너의 psql을 사용했다.
- DEV `.env.server`의 접속값은 메모리에서만 읽었다. 기존 `C:/Users/shin/.ssh/known_hosts`의 대상 일치와 SSH 접속을 확인했다. 호스트명·사용자명·키·비밀번호는 기록하지 않았다.
- Paramiko는 기존 설치를 재사용했다. 명시적 known_hosts와 `RejectPolicy`를 사용했다. 신뢰 파일을 바꾸지 않았다.

## 실제 동기화 범위

기존 `frontpolish001-coor/remote.py`와 `make.cmd`는 읽기만 했다. 두 파일의 ROOT는 coor다.
실행한 새 helper의 ROOT는 `C:/Users/shin/orca/workspaces/KnowsLink/dev`다.
`git bundle create <외부 파일> HEAD`와 원격 새 clone은 지정 HEAD의 추적 Git 내용만 재현한다.
미커밋·미추적·ignored 파일은 동기화하지 않는다. `.env`, `.env.server`, `.fullops-squad/.env`는 ignored임을 확인했다. 추적된 `.env.example`, `.fullops-squad/.env.example`는 clone에 포함된다. 로컬 작업 트리 변경을 원격 후보로 주장하지 않는다.
helper는 요청 SHA가 현재 부모 HEAD인지 확인한다. 원격 HEAD 일치와 깨끗한 추적 상태도 확인한다.

외부 helper 위치:
`C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/`

- `remote.py`: probe/prepare/runtime/auth/collect/make 및 자기 검사.
- `runtime_remote.py`: 새로운 전용 내부망·Postgres·migration·loopback relay·SMTP sink를 준비한다.
- `auth_remote.py`: 실제 OTP·세션과 인증 실패 경계를 검사한다.
- `make.cmd`: 부모가 FullOps gate에서 사용하는 선택형 bridge다. child는 최종 lint/test 게이트를 실행하지 않았다. 부모 게이트 결과는 부모 보고서에서 확인한다.

## 실행과 증거

실행 cwd는 DEV 루트다. 아래 `python3`는 실제 Python 3.13 경로로 해석됐다.

```powershell
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py self-check
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py probe
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py prepare --sha 59d419d0b7c4fca46cabafc04b5a47c6233c944a
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py runtime --sha 59d419d0b7c4fca46cabafc04b5a47c6233c944a
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py auth --sha 59d419d0b7c4fca46cabafc04b5a47c6233c944a
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py collect --sha 59d419d0b7c4fca46cabafc04b5a47c6233c944a
```

부모는 준비 SHA에서 make install/build와 격리 fixture 준비만 허가했다.
`make install`, `make build`, SQL migration, `/healthz`가 종료코드 0을 반환했다.
두 합성 회원의 SMTP sink OTP 요청 303, 틀린 코드 401, 올바른 검증 303, 인증 home 200을 확인했다.
미인증 home 303과 잘못된 이메일 422도 확인했다.
비밀값은 출력하지 않았다. `.invalid` 수신자와 로컬 SMTP sink만 사용했다. Google 설정은 비웠다.
요청 허용 범위는 GET `/`, `/healthz`, `/home` 및 POST `/auth/start`, `/auth/verify`다.

원격 준비 clone:
`/tmp/knowslink-humanpeer001-59d419d0b7c4-02300ad1`

성공 fixture:
`/tmp/knowslink-humanpeer001-59d419d0b7c4-02300ad1/qa-34e6309f`

- network: `qa-34e6309f-02300ad1`, internal, DB 공개 포트 없음.
- Postgres: `qa-34e6309f-02300ad1-postgres`.
- relay: 원격 `127.0.0.1:38135`, PID 4024298.
- SMTP sink: 원격 `127.0.0.1:51795`, PID 4024297.
- 환경과 세션·합성 이메일 원문은 원격 0700 디렉터리 안의 0600 파일에만 보존했다. Git 증거에는 복사하지 않았다.

공개 가능한 증거는 이 디렉터리 아래 준비 SHA 폴더에 수집했다.
`qa-34e6309f/effective-go.log`, `build.log`, `migrate.log`, `postgres-version.log`, `auth-selfcheck-*.json`과 실행 로그에 근거가 있다.
명령의 원래 종료코드를 보존했다. 출력 축약은 종료코드 확인 후 표시만 했다.

최초 probe는 호스트 psql 부재에서 127로 끝났다. 후속 probe가 나머지 버전을 확인했다.
최초 runtime 시도는 내부망의 published port 조회 실패로 1을 반환했다.
해당 `qa-da1c0627` 디렉터리·망·Postgres를 보존했다. 후속 fixture는 공개 포트를 만들지 않고 내부 IP로 접속했다.
실패 기록을 성공으로 바꾸지 않았다.

## 기준 환경 준비 시점의 실행 계약

부모가 구현을 커밋하고 고정 SHA와 준비 신호를 전달하면 새 clone을 만든다.
부모가 `make lint`, `make test`, FullOps lint.py의 최종 실행과 커밋을 소유한다.
선택형 bridge는 외부 helper 폴더를 PATH 앞에 붙여 사용한다. bridge는 원래 make 종료코드를 반환한다.
이 worker는 부모와 조율한 고정 후보에서만 실제 합성 auth/API/DB와 320/390/1440 캡처를 수행한다.
기존 Orca 브라우저에서 task loopback만 연다. 고정 page ID를 사용한다.
정적 HTML fixture와 실제 API/DB 화면의 검증 결과를 구분한다.
탐색·shadow·CUAWright는 OFF다. 제품 수락은 coordinator의 독립 tester/reviewer/designer 후속이다.

저장소 코드·기존 제품 파일은 수정하지 않았다. 이 child는 Git 커밋을 하지 않는다.
FullOps product-lint/product-test gate는 부모의 구현 커밋에서 실행할 항목이다. 이 연구의 환경 준비 증거에는 적용하지 않았다.
삭제·이동 정리를 하지 않았다. 실패 시도와 모든 실행 자원은 보존한다.

## 문서 근거

Context7 resolve-library-id가 quota 오류를 반환했다. 공식 문서로 대체했다.
[Paramiko SSH client](https://docs.paramiko.org/en/stable/api/client.html)로 known_hosts·RejectPolicy·명령 실행 API를 확인했다.
[Paramiko channel](https://docs.paramiko.org/en/stable/api/channel.html)로 명령 종료코드 API를 확인했다.
[Docker network create](https://docs.docker.com/reference/cli/docker/network/create/)로 내부망의 격리와 호스트 접속 범위를 확인했다.
[Docker container run](https://docs.docker.com/reference/cli/docker/container/run/)으로 전용 network·env-file 실행 인자를 확인했다.
Orca의 설치 바이너리에서 orchestration·orca-cli 및 browser reference를 읽었다.
`--type update`는 설치 CLI에서 지원하지 않아 실패했다. 진행 보고는 지원되는 `--type status`로 전달했다.

## 브라우저 경로 준비

loopback SSH forward의 `/healthz`가 Windows에서도 200을 반환했다.
기준 fixture를 향한 local origin은 `http://localhost:8292`다. forward 프로세스 PID는 48660이다.
Orca 1.4.224의 별도 isolated profile `4359b1d6-c7ae-48a2-b342-1e672c66896a`를 만들었다.
고정 page ID는 `d6f73e04-f39d-4a0e-9cc2-9588a1d153a4`다.
실제 합성 OTP로 만든 세션을 메모리에서 읽고 해당 profile에 Secure·HttpOnly·SameSiteStrict cookie로 주입했다.
브라우저 `/home` 도착과 main 1개를 확인했다. 세션값은 출력하거나 증거에 복사하지 않았다.
기준 UI의 캡처는 만들지 않았다. 후속 구현 화면의 캡처와 혼동하지 않는다.
`forward.py`, `browser_session.py`, `capture.py`는 외부 helper 폴더에 보존했다.
캡처 runner의 SHA·origin·고정 page·GET 경로 검증과 Python compile만 실행했다.
320/390/1440 실제 캡처는 부모의 고정 구현 후보와 준비 신호를 기다린다.
[Paramiko transport](https://docs.paramiko.org/en/stable/api/transport.html)로 기존 SSH의 loopback 채널 API를 확인했다.

## Git bundle 범위 확인

`git bundle list-heads`는 준비 SHA 하나를 `HEAD`로 표시했다.
이 bundle은 해당 HEAD와 그 이전의 도달 가능한 Git 이력도 포함한다. 단순한 현재 파일 tar가 아니다.
준비 HEAD의 도달 가능한 커밋은 944개다. 다른 branch의 최신 미통합 내용과 작업 트리는 포함하지 않는다.
`git rev-list <준비 SHA> -- .env .env.server .fullops-squad/.env`는 비어 있었다.
이 세 비공개 환경 파일은 준비 HEAD의 이력에도 없음을 확인했다. 환경 파일 내용은 읽어 출력하지 않았다.


## 구현 후보 실제 실행

부모는 `4eaefa6e0a5cf623f0d8d8080d918086f86030ec`를 전달하고 prepare/runtime/integration과 실제 auth/UI fixture 실행을 허가했다.
이 후보는 독립 리뷰 전 상태다. 부모는 기술 문서와 이후 회귀 fixture 호환성을 별도로 갱신한다.

```powershell
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py prepare --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py runtime --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py integration --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py auth --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py fixture --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py collect --sha 4eaefa6e0a5cf623f0d8d8080d918086f86030ec
```

실행 당시 integration action은 다음 focused 명령을 실행했다. 이후 부모 요청으로 기본 동작을 전체 suite로 바꿨다.

```text
go test -race -tags=integration ./internal/relay -run TestHumanPostgres -count=1
```

- make install/build, SQL migration, healthz: 종료코드 0.
- focused `TestHumanPostgres`: 3.101초, 종료코드 0. 실제 tagged 테스트 존재를 먼저 확인했다.
- 별도 DB `humanpeer_integration_9e418000ba45`를 만들었다. relay DB와 공유하지 않았다.
- 두 합성 회원의 실제 로컬 SMTP OTP·잘못된 코드 401·이름 onboarding·home 200이 통과했다.
- 공개 폼으로 세 에이전트를 만들었다. 같은 이름 두 개와 긴 한국어 이름 하나를 사용했다.
- 기존 adapter prepare/confirm/complete로 각 공개키와 credential 연결을 완료했다. 개인키는 Node 클라이언트의 private 파일에만 보존했다.
- 실제 폼으로 친구 코드를 발급하고 요청·명시적 수락을 완료했다.
- 다른 Origin의 composer 요청은 403이었다. 같은 Origin의 웹 owner 작성은 303이었다.
- 실제 Node adapter가 owner 본문을 검증·저장하고 persist/ACK를 완료했다.
- 반대 agent가 자기 로컬 키로 서명한 답장을 보냈다. 다른 adapter가 검증·저장·persist/ACK를 완료했다.
- 메시지 본문·키·credential·OTP는 로그에 쓰지 않았다. 직접 State JSON/SQL 주입은 하지 않았다.

DB/container/network 생성과 빈 schema migration은 fixture provisioning이다.
제품 상태 변경은 실제 SMTP 인증·HTTP 폼·공개 연결 API·기존 adapter로 실행했다.
Go integration 테스트의 내부 fixture·시간 제어는 테스트 소스에 속한다. 실제 벽시계 180초 대기 검사로 표시하지 않는다.

원격 clone: `/tmp/knowslink-humanpeer001-4eaefa6e0a5c-c63728a4`

실제 fixture: `/tmp/knowslink-humanpeer001-4eaefa6e0a5c-c63728a4/qa-c8a0339a`

- Postgres: `qa-c8a0339a-c63728a4-postgres`, 내부망 `qa-c8a0339a-c63728a4`, 공개 DB 포트 없음.
- remote relay: loopback 58495, PID 4071560.
- remote SMTP sink: loopback 42473, PID 4071559.
- Windows SSH forward: PID 36448, `http://localhost:14228`.
- 미리보기: `http://localhost:14228/home`.
- 기존 isolated Orca profile/page에서 Alice 합성 세션을 사용한다. 세션값은 전달하지 않았다.

증거는 이 디렉터리의 구현 SHA 폴더에 있다.
`qa-c8a0339a/humanpeer_integration_9e418000ba45/integration-result.json`은 정확한 명령·종료코드·로그를 연결한다.
`qa-c8a0339a/auth-selfcheck-*.json`, `human-fixture-result.json`은 인증·실제 요청·adapter 동작 범위를 연결한다.

## 실제 화면 캡처와 한계

정상 증거: `4eaefa6e0a5cf623f0d8d8080d918086f86030ec/captures-1791801092715284800/manifest.json`.
대화 목록·친구·내 에이전트·새 대화·실제 대화의 5개 화면을 3개 CSS 폭에서 캡처했다.
15개 상태의 전체 문서를 21개 viewport PNG로 나눴다. 이미지 합성·변형은 하지 않았다.
모바일 viewport는 320×640과 390×844, DPR 1이다.
데스크톱 viewport는 CSS 1440×900, DPR 0.75이며 PNG는 1080×675다.
Orca 1.4.224와 Chrome 150.0.7871.250을 사용했다. 시나리오의 탐색·shadow·CUAWright는 OFF다.
각 PNG의 실제 크기·SHA-256·scroll 위치와 실제 viewport·DPR을 manifest에 남겼다.
DOM 측정은 가로 넘침 0, 44px 미만의 보이는 제어 0이었다. 측정은 독립 디자인 수락이 아니다.
대화의 메시지·composer와 내 에이전트 일부 화면을 직접 읽어 캡처의 반복 타일이 사라졌음을 확인했다.

최초 `captures-1791800717129804500`의 full-page PNG는 native surface를 넘는 영역에서 타일이 반복됐다.
해당 증거는 부적합으로 보존했다. 같은 현상이 DPR 1의 1440px viewport 오른쪽에서도 나타났다.
DOM의 top header/main/brand는 각각 하나였고 가로 넘침은 없었다. 제품 중복 렌더링 결함으로 판정하지 않았다.
지원되는 `set viewport 1440 900 0.75`와 개별 viewport screenshot으로 반복 타일을 피했다.
데스크톱의 CSS 반응형 검사는 1440 폭이다. native 1440px·DPR 1 이미지 검사로 표시하지 않는다.
독립 designer가 필요하면 실제 1440px surface의 DPR 1 캡처를 추가한다.

키보드 시도는 통과하지 않았다. `keyboard-blocked.json`을 정상 증거 폴더에 보존했다.
Orca 문서의 `document.hasFocus()`가 false였고 Tab 뒤 activeElement는 BODY였다.
명시적 DOM focus는 skip link를 선택했지만 Enter는 활성화되지 않았다.
이는 native 포커스 부재 상태의 실행기 제한 근거다. 제품 focus-visible/키보드 결함 또는 PASS로 표시하지 않는다.
coor/designer가 실제 보이는 browser surface에 포커스를 준 뒤 키보드 검사를 수행한다.

## 최종 전체 회귀

부모 요청으로 integration action의 기본 명령을 아래와 같이 확장했다.
각 실행은 새로운 별도 DB를 만들고 migration을 적용하며 원래 종료코드와 비밀값 제거 로그를 보존한다.

```text
go test -race -tags=integration ./internal/relay -count=1
```

선택형 `--pattern TestHumanPostgres`는 focused 명령을 재현한다.
전체 suite는 UI child의 기존 identity/connections/public_messages/google/device/remote fixture 호환성이 정리된 최종 후보에서 실행한다.
`f567abff18b8f6ce77fc79e6601ba7163ae7f9ff`를 정확히 prepare한 뒤 전체 suite를 실행했다. 113.318초 후 종료코드 1이었다. focused PASS를 전체 회귀 PASS로 바꾸지 않는다.
부모는 최종 SHA에서 FullOps product-lint/product-test와 커밋·통합을 수행한다.
후보 SHA가 바뀌면 관련 제품/테스트/설정 의존성 delta를 확인하고 기존 증거의 재사용 또는 재실행 근거를 기록한다.

### f567abf 실행 분리와 캡처 재사용

부모 요청에 따라 정확한 `f567abff18b8f6ce77fc79e6601ba7163ae7f9ff`를 준비했다.
원격 base는 `/tmp/knowslink-humanpeer001-f567abff18b8-0b335c67`이다.
부모 make bridge는 base의 `repo`를 사용한다. child runtime은 `qa-8154ee89/repo`에 같은 SHA를 별도 clone하여 install/build했다.
이 분리로 부모 make와 child build/test가 같은 생성물 경로를 쓰지 않는다.
install/build, migration, healthz는 종료코드 0이다.
전용 Postgres는 `qa-8154ee89-0b335c67-postgres`, 내부망 `qa-8154ee89-0b335c67`, private IP `172.28.0.2:5432`이며 공개 DB 포트는 없다.
relay는 loopback 43813/PID 4121268, SMTP sink는 loopback 45051/PID 4121267이다.
환경 파일·DB·컨테이너·프로세스·로그는 보존한다.

```powershell
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py prepare --sha f567abff18b8f6ce77fc79e6601ba7163ae7f9ff
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py runtime --sha f567abff18b8f6ce77fc79e6601ba7163ae7f9ff
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py integration --sha f567abff18b8f6ce77fc79e6601ba7163ae7f9ff
```

4eaefa6→f567abf의 제품 delta를 확인했다. human/member/remote 이름 input의 maxlength는 40에서 80 UTF-16 unit으로 바뀌어 서버의 40 Unicode 문자 제한과 맞춘다.
member agent의 두 제한 안내 문구가 사람 친구 모델을 설명하도록 바뀌었다. 기존 캡처는 이 제한 상태를 표시하지 않는다.
remote MCP receive는 기존 signed wire를 그대로 두고 응답에 메시지별 authorship과 untrusted를 추가한다. 캡처된 local adapter 대화 흐름과 다른 경로다.
다른 변경은 회귀/기존 integration fixture 호환성, consent browser fixture 이름, 문서/schema generator, 배포 skill/package 안내다.
CSS·레이아웃·캡처한 화면의 메시지 렌더 구조는 바뀌지 않았다. 부모 지시에 따라 기존 21개 PNG를 재사용하며 원래 4eaefa6 SHA를 유지한다.
새 maxlength 입력 경계·수정 안내 상태·remote MCP 응답은 기존 PNG가 검증하지 않는다. 이 범위는 새 후보 테스트와 독립 검수 대상이다.

### f567abf 전체 SQL 실패 기록

별도 DB `humanpeer_integration_ac01446c7d41`의 생성과 migration은 종료코드 0이었다.
전체 `go test -race -tags=integration ./internal/relay -count=1`은 종료코드 1이었다.
증거는 `f567abff18b8f6ce77fc79e6601ba7163ae7f9ff/qa-8154ee89/humanpeer_integration_ac01446c7d41/`의 `integration-result.json`, `go-test-relay-integration.log`다.
로그는 잘라내지 않은 비밀값 제거 원본이며 실패 출력·native exit를 보존한다.

- `TestPublicAgentHTTP/same_owner_legacy_refusal_generation_expiry_and_cleanup_budget`: `connections_integration_test.go:180`, expected 403 / actual 303.
- `TestPostgresSafety/owner_gate_csrf_binding_and_no_effect`: `integration_test.go:226`, `verified UI missing`.
- `TestGateFailureStates/denied`, `/expired`, `/revoked`, `/unavailable`: `integration_test.go:403`, `unsafe ... UI`, actual status 200.

읽기 triage에서 gate 검사는 raw `pending`과 연속 `상태: <state>` 문자열을 기대하지만 현재 페이지는 badge markup을 사용한다.
실패 출력에는 승인 control이 없는 상태 설명이 보인다. stale assertion 가능성이며 제품 안전성 PASS로 해석하지 않는다.
same-owner legacy 거부는 현재 redirect 결과와 실제 pair 상태를 함께 검사해야 한다.
부모에게 실패와 위치를 전달했다. child는 제품·fixture 소스를 수정하지 않았다. 수정 후보에서 필요한 회귀를 다시 실행한다.

### 7f83fd6 전체 회귀 재실행

부모는 legacy 반복 accept의 idempotent 303 결과와 실제 상태 불변을 함께 검사하도록 fixture를 수정했다.
gate fixture는 localized badge·검증 본문·승인 control과 거부/만료/철회/부재 상태의 승인·거절 control 부재를 검사한다.
실제 A2 self-send/pull/persist/ACK와 production transaction을 사용하는 `TestHumanLegacyPostgres`도 추가했다.
이 변경을 포함한 `7f83fd6f3ae9945e292a86212c0ab566ea282b7a`에서 전체 suite를 다시 실행했다. 118.664초 후 종료코드 0이었다.

새 원격 base는 `/tmp/knowslink-humanpeer001-7f83fd6f3ae9-a0d3aeba`, private source는 `qa-2c401650/repo`다.
install/build/migration/healthz는 종료코드 0이다.
Postgres는 `qa-2c401650-a0d3aeba-postgres`, 내부망은 `qa-2c401650-a0d3aeba`, private DB 주소는 `172.29.0.2:5432`다.
공개 DB 포트는 없다. relay는 loopback 45897/PID 4143191, SMTP sink는 loopback 53909/PID 4143190이다.

```powershell
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py prepare --sha 7f83fd6f3ae9945e292a86212c0ab566ea282b7a
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py runtime --sha 7f83fd6f3ae9945e292a86212c0ab566ea282b7a
python3 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/remote.py integration --sha 7f83fd6f3ae9945e292a86212c0ab566ea282b7a
```

f567abf→7f83fd6의 제품 delta는 이름에서 모든 Cf 대신 Control/Bidi_Control만 금지하여 ZWJ를 허용하는 경계 수정,
`unavailable` badge의 한국어 표시, 그리고 http.go의 gofmt다. 나머지는 테스트·문서다.
기존 합성 한국어 이름·활성 대화의 캡처 내용과 CSS는 바뀌지 않는다.
기존 21개 PNG의 SHA를 유지한다. 새 이름 경계와 unavailable gate 상태는 기존 PNG의 검증 범위에 포함하지 않는다.

별도 DB는 `humanpeer_integration_9034dbf57cc8`이며 생성·migration·전체 go test가 모두 종료코드 0이다.
증거는 `7f83fd6f3ae9945e292a86212c0ab566ea282b7a/qa-2c401650/humanpeer_integration_9034dbf57cc8/`의 `integration-result.json`과 `go-test-relay-integration.log`다.
명령에는 `-run` 필터가 없다. HumanPostgres/Legacy 소스의 환경 사전조건과 runner의 `TEST_SYNTHETIC_DATABASE=1` 및 전용 `TEST_DATABASE_URL` 설정을 확인했다.
진행 중 추가 제품 변경 없이 전체 relay tagged 테스트를 실행했다. 테스트별 verbose 실행 기록은 만들지 않았다.

마지막 evidence integrity 자체 검사는 원래 PNG 21개의 해시·크기, 15개 상태의 landmarks/overflow/target size, 텍스트 증거의 서버 비밀번호 부재를 확인했다.
외부 helper의 assert self-check와 Python 구문 검사도 종료코드 0이다.
이 child는 제품 소스·git commit·최종 FullOps gate를 소유하지 않는다. 모든 임시 파일과 실행 리소스는 보존했다.
남은 한계는 native focus 없는 키보드 미완료, CSS1440/DPR0.75 캡처 조건, 원래 캡처 SHA 이후의 명시된 경계 변경이다.
