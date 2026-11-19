---
title: SAR-CONNECT-SIMPLE-001-DEV 자체 검사 보고
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-DEV]
summary: "에이전트 한도 제거·키 저장 전 시도 분리·기본 이름·재시도 폴더·실패 사유·연결 화면 구현을 격리 Postgres·합성 fixture·Windows C:/D: 실측으로 검사했다. CS-AC-01–16의 DEV 검사가 통과했고 tester QA·designer 시각 검수(CS-AC-17)는 후속이다."
---

# SAR-CONNECT-SIMPLE-001-DEV 자체 검사 보고

DEV 자체 검사다. 제품 규칙은 [D02](../../../planning/product-specs/SAR-CONNECT-SIMPLE-001.md), 모듈 책임은 [D10 절](../../../design-docs/module-design.md#sar-connect-simple-001-연결-단순화-책임)이다. 운영 relay·운영 agent·실계정·실제 상대 발송·사용자 실제 `~/.claude`·`~/.codex`는 쓰지 않았다.

- 브랜치 `fullops/dev`. 코드 후보 SHA `310fdd73428deec2fadb4c7a8a1172901af52ebe`(기준 `9ad346fe`).
- 원시 로그: `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-DEV/`(gitignore). Linux 게이트 작업 폴더: `D:/workspace/KnowsLink/.devtools/connect-simple-001-dev/`.

## 필수 검사

| 검사 | 환경 | 결과 |
|---|---|---|
| `make lint` | Linux 격리 clone, `310fdd73` | EXIT=0 (`310fdd73…-lint.log`) |
| `make test` | 같은 clone | EXIT=0 (`310fdd73…-test.log`) |
| Go 통합 전체 `go test -race -tags=integration ./internal/relay` | 같은 clone, 격리 Postgres 17 컨테이너(검사 뒤 정지) | EXIT=0, `ok … 151.057s` (`310fdd73…-integration.log`) |
| `npm run check`, `npm test`(13개) | Windows, Node 22.22.2 | EXIT=0 (`parent-check.log`, `parent-test.log`) |
| `scripts/package_plugin.py --verify` | Windows | EXIT=0 (`parent-verify.log`) |
| FullOps `lint.py --from 9ad346fe` | `310fdd73`의 깨끗한 detached worktree | ERROR 0, WARNING 19 (`fullops-lint.log`) |

FullOps lint를 깨끗한 worktree에서 실행한 이유: dev 작업 트리에 다른 과제(SAR-MOTION-INTRO)의 미추적 캡처·`videos/`가 있어 lint가 실행을 거부했다. 그 파일은 이 과제 소유가 아니라 옮기거나 지우지 않았다. lint의 `make lint`·`make test` 항목은 Windows에 `make`가 없어 '실행 불가'로 표시되며 위 Linux 결과로 대신한다.

## 수락 기준별 DEV 검사

| ID | 검사 | 결과 |
|---|---|---|
| CS-AC-01 | `TestAgentAndPairCapacity`(같은 회원 활성 7개), integration `TestPublicAgentHTTP`(여섯 번째 이후 303), `TestConnectScreens`·`TestSaturationGuidance`(한도 문구 '회원당'·'먼저 철회' 없음) | 통과 |
| CS-AC-02 | `TestAgentAndPairCapacity`(철회 10개 보존 회원), integration `revoked_records_never_refuse_restart_and_retention`(철회 6개 + 동시 8건 모두 303) | 통과 |
| CS-AC-03 | `TestAgentAndPairCapacity`(전체 201개 생성 허용, 활성 2000에서 `service_capacity` 문구만) | 통과 |
| CS-AC-04 | `TestAgentRate`(21번째 `agent_rate`, 재시도 시각, 철회 포함, 기존 credential 유지), `TestAttemptsAreNotAgents`(키 저장 전 시도 5회 비집계), integration `concurrent_agent_rate_and_restart`(동시 8건 중 5건, 재시작 뒤 429) | 통과 |
| CS-AC-05 | `TestAttemptsAreNotAgents`(대기·동의 미저장·취소·만료·정리 5회 뒤 agent 수 같음), integration `TestGoogleDeviceHTTP`(동의 뒤 agent 0) | 통과 |
| CS-AC-06 | `TestAttemptsAreNotAgents`(기한 + 24시간 정리), `TestDeviceFailureBoundaries`(기한 뒤 완료), integration(정리 버튼 뒤 poll 401) | 통과 |
| CS-AC-07 | `TestDefaultNames`(호스트별·번호·모름·영어), integration(동의 화면 `value="Claude Code"`), Dots 기본값 | 통과 |
| CS-AC-08 | `TestConnectScreens`(이름 절이 접히지 않은 폼) | 통과. 화면 확인은 tester |
| CS-AC-09 | adapters `connect-simple.test.ts` 출력 스냅숏(connect_status·status·send·receive의 `me`/`peer`, 키·토큰 없음), integration 자기 contact 라벨 | 통과 |
| CS-AC-10 | Windows 실측(아래), `connect-simple.test.ts`(남은 시도 이동·빈 폴더·`private_folder`) | 통과 |
| CS-AC-11 | `connect-simple.test.ts`(사유 8개 KO/EN), `grep-reasonless-failed.txt`(세 경로 사유 없는 `failed` 0건) | 통과. `signed_out`은 웹 전용 |
| CS-AC-12 | integration(시작 회원 재방문 '이 연결은 완료됐습니다'·이름, 로그아웃 401 일반 문구·이름 없음) | 통과 |
| CS-AC-13 | `TestConnectFailureMessages`(세 원인 서로 다른 문구, 이메일 없음, 계정 선택 폼) | 통과. Google 실제 화면은 tester |
| CS-AC-14 | `connect-simple.test.ts`(release `defaults.json`만으로 동의 전 요청 0건, 동의 뒤 start, env 우선) | 통과. 실제 Grok Bot 미실행 |
| CS-AC-15 | `package_plugin.py --verify` EXIT=0, C-10 원인 아래 | 통과 |
| CS-AC-16 | `make test`·통합 전체(지문·동의·비공개 폴더·발송 승인·친구 권한 기존 테스트), ACL 결과 비교 | 통과 |

## C-10 원인과 수정

- 원인: main 체크아웃 `D:/workspace/KnowsLink`의 `internal/relay/registry.json` 작업 파일이 CRLF다(`git ls-files --eol`: `i/lf w/crlf`). `.gitattributes eol=lf` 전에 받은 파일이 다시 쓰이지 않았다. `mcp.test.ts`가 그 바이트를 fake relay로 내보내면 클라이언트는 고정 LF hash와 달라 `/v1/registry` 뒤 멈춘다. 제품 동작은 맞고 테스트가 체크아웃 줄바꿈에 의존한 결함이다.
- 수정: fixture를 커밋된 LF 형태로 읽고 고정 hash 일치를 먼저 assert한다. 경로 순서 검사는 그대로다.
- 증거: `c10-repro-crlf.log`(수정 전 CRLF로 같은 실패), `c10-after-fix.log`(같은 CRLF로 통과, 1바이트 변경 registry는 실패).
- 남은 일: main 체크아웃의 CRLF 파일을 다시 받을지는 사용자·coor 결정이다(파일 재작성).
- 관찰: 첫 `npm test`에서 `mcp.test.js`가 한 번 실패했고 원인을 남기지 못했다. 바로 21회 연속 통과했다(`mcp-loop-*.log`, `mcp-cold-*.log`).

## L-05 영향

legacy pair는 명시적 초대(`/v1/invites`·`/home/invites`)로만 생긴다. 친구 전달 권한은 회원 쌍의 `Friends` 기록으로 계산한다(`authority`). agent 수가 늘어도 legacy pair 전체 한도(active 400, pending 200)를 쓰지 않고 친구 전달도 바뀌지 않는다. designer 판단 요청은 필요하지 않았다.

## Windows 실측(R-01·R-02)

로컬 fake relay, C: `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-DEV/measure/`, D: `D:/workspace/KnowsLink/.devtools/connect-simple-001-dev/measure/`(`windows-measure.txt`, EXIT=0).

| 경우 | C: | D: |
|---|---|---|
| 취소된 시도 파일이 남은 폴더 | 연결 성공, `.unused-*`에 파일 보존 | 같음 |
| 대기 중 시도가 남은 폴더 | 같은 시도 재개, 새 start 0건 | 같음 |
| 미리 만든 빈 폴더 | 연결 성공, 이동 없음 | 같음 |
| 새 폴더 | 연결 성공 | 같음 |

8개 폴더의 권한 결과가 같다: 소유자 현재 사용자, 상속 차단, Allow는 현재 사용자 FullControl 하나. 수정 전 D: 실패 원인은 `acl-probe-setowner.txt`: D: 새 폴더는 WRITE_OWNER가 없어 같은 소유자를 다시 쓰는 것도 거부된다. 소유자가 다를 때만 `SetOwner`를 부르도록 고쳤다. 검사 조건은 완화하지 않았다.

## 미실행과 이유

- 실제 Grok Bot 앱·Linux Bot에서 `install_bot_mcp.sh`, 실제 relay로 설치기 끝까지 실행: tester 범위(CS-AC-14·18).
- T-04 실측: 이 컴퓨터에 Grok Bot 앱이 없다. 문서상 Bot별 신호가 없어 같은 agent를 쓰고 `knowslink_status`가 `sharedAgent`로 알린다.
- 회원 화면 1440·390·320px·키보드 확인(CS-AC-17): designer 시각 검수.
- `knowslink_text_receipt`·`knowslink_session` 결과의 이름: A-06 목록(상태·연결·보낸·받은) 밖이라 넣지 않았다.

## lint 경고와 의존성

- DEP-001 `adapters/package.json`: 의존성 변경 없음. `test` 스크립트에 새 테스트 파일 실행만 더했다.
- SEC-001 `connect-simple.test.ts:73`: fake relay 응답의 고정 문자열 `lease_token: "lease"`이다. 비밀값이 아니다.
- SIZE-001·SIZE-002: 이미 300줄을 넘던 파일에 기능이 더해졌다. 변경 추가 1734줄(지시서 예상 600–900줄)은 새 테스트(Go 207줄, TS 419줄)와 문서가 절반 가까이다. 커밋은 relay·adapters·문서로 나눴다(`e7f4826f`, `41cec436`, `0c21f209`, `ddda8706`, `310fdd73`).
- SLOP-002 조건부 객체 스프레드 3곳: 선택 필드(`retry_at`·`agent`)를 JSON에 넣을 때만 쓴다. 바꾸지 않았다.

## 배포 주의

relay를 클라이언트·설치기보다 먼저 배포한다. 이전 relay는 `/v1/connect/start`의 `host`를 422로 거절하고 클라이언트는 `unreachable`로 보인다.
