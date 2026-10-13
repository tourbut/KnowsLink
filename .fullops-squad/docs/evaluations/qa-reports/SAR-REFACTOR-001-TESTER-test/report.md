---
title: Sol 리팩토링 독립 경계 QA
status: review
updated: 2026-10-13
owner: tester
tasks: [SAR-REFACTOR-001-TESTER]
summary: 고정 cd4dae8의 원본 게이트를 대조하고 Linux에서 연결 응답·사람 선택·text·MCP 경계를 독립 검증했다.
---

# SAR-REFACTOR-001-TESTER 독립 QA

판정은 이번 기능 보존 리팩토링 범위에서 PASS다. 새 제품 결함과 미해결 critical/high 발견은 0건이다. 제품 코드는 변경하지 않았다. 전체 운영 수락과 coordinator의 독립 코드 리뷰는 별도다.

## 대상과 독립성

- 기준 diff: `0eb7735f0963272b62266992791f9f09152f9266` → `cd4dae810ae262759bcff29895165f1a64f8c51a`.
- 준비 HEAD: `eea94316418cfd8ccc501e6ca66ed82a2abc3565`. 제품 경로는 후보와 동일하다. DEV 최종 문서 `e65146b`도 제품이 동일하다.
- tester 실제 세션: `01a13534-4add-7a22-bb3f-5e40c8addb09`. Codex `gpt-6.1-sol` / high. Task `task_bf0eee88d603`, Dispatch `ctx_6e75a377f476`.
- 원 구현 부모 `01a13509-e1a3-7b62-bf9d-5a48405d706a`, Go `01a1350c-1383-7323-9d75-748304b33792`, TS `01a1350c-7f6b-7f03-822b-e3a3ec516f30`, 완료 재개 `01a1351f-c007-7f22-ac7c-1837dca51f0c`와 다른 세션이다.
- 적용 기준: `fullops-common-0.3.3`, AGENTS.md, FULLOPS.md, project.md, 공통 세 규칙, D02 SAR-HUMAN-PEER, D10 module-design, DEV 실행 기록과 QA 원본. Test level lite. 좁은 QA는 직접 수행했다.

## 원본 전체 게이트 재사용

[gate-audit.json](gate-audit.json)은 DEV의 원본 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor001/cd4dae810ae262759bcff29895165f1a64f8c51a-<target>.log`와 Git의 [DEV gate](../SAR-REFACTOR-001-DEV-test/gate/)를 대조한다. install/lint/test/build의 내용은 LF 정규화 뒤 동일하다. 원본과 사본의 SHA256을 각각 남겼다. Windows checkout의 개행 변환 때문에 바이트 동일성을 주장하지 않는다.

| 검사 | 원본 결과 | 재사용 근거 |
|---|---|---|
| `make install` | EXIT=0 | 고정 후보의 lockfile 설치 |
| `make lint` | EXIT=0 | format/vet/modules, adapter format/ESLint/types, 설정·공개 ingress 검사 |
| `make test` | EXIT=0 | Go race, TS 등록 suite, shell guard 4 tests. Go 결과 일부는 정상 캐시다. |
| `make build` | EXIT=0 | Go relay/migrate, TypeScript, plugin bundle |
| FullOps 1.4.1 등록 게이트 | product-lint/product-test exit 0 | JSON head가 정확한 후보다. ERROR 0, WARNING 3, unavailable 0이다. |

DEV FullOps 기준은 `bc0fd031efd252d9755681d8b429eb31c045b7e4`다. tester 기준 `0eb7735f0963272b62266992791f9f09152f9266`와 제품 경로가 같음을 확인했다. 경고는 PLANS/core/human_test 크기 3건이다. 원래 사라진 1.4.0 경로의 실패는 DEV 기록에 그대로 둔다. lint.log의 public config 거절 문구는 예상 부정 검사이며 전체 EXIT=0과 함께 판정했다. 재사용 결과를 tester의 새 전체 게이트 실행으로 표시하지 않는다.

## 이번 독립 실행

실행 위치는 기존 고유 clone `/tmp/knowslink-refactor001-cd4dae810ae2-45997998/repo`다. 검증 전후 HEAD는 정확한 후보이며 `git status --porcelain`은 비어 있다. Node `v22.22.2`, npm `10.9.7`, Go `go1.27.1 linux/amd64`, CGO_ENABLED=`1`을 확인했다. SSH helper의 known_hosts와 RejectPolicy, 비밀값 가림을 유지했다. 접속값이나 인증 URL은 기록하지 않았다.

전체 명령과 실제 종료코드는 [commands.json](commands.json)에 있다. 모든 실행은 `set -e`와 exact HEAD·clean 검사를 먼저 수행했다. SSH channel의 종료코드를 직접 기록했다. 파이프로 명령 종료코드를 가리지 않았다.

| 검사 | 결과·증거 |
|---|---|
| 지원 환경·exact SHA·clean 검사 | 0 · [environment-verified.log](environment-verified.log) |
| `go test -race ./internal/relay -count=1 -v -run '^(TestConnectionResponseKeepsRetryAndErrorPriority\|TestDeviceFailureBoundaries\|TestHumanAgentListsKeepOwnerAndFriendBoundaries\|TestHumanFriendTransitionsAndDiscovery\|TestHumanGateRevocationAndReaccept\|TestFrozenParsingAndSigning\|TestDeviceConsentTemplate\|TestMemberUIFixtures)$'` | 0 · 실제 최상위 8 tests 통과 · [go-boundaries.log](go-boundaries.log) |
| `npm run build --prefix adapters` | 0 · 현재 소스의 TS·MCP bundle을 다시 만들어 focused 검사의 stale artifact 가능성을 없앴다 · [ts-build.log](ts-build.log) |
| `node adapters/dist/mcp.test.js && node adapters/dist/owner-text.test.js` | 0 · [ts-boundaries.log](ts-boundaries.log) |
| 독립 `node --input-type=module` text matrix | 0 · 15개 경계 · [text-matrix.log](text-matrix.log) |
| 실행 뒤 exact SHA·clean 확인 | 0 · [post-clean.log](post-clean.log) |

첫 환경 프로브는 JS 확인식의 셸 따옴표 오류로 EXIT=1이었다. [environment.log](environment.log)와 원 명령을 보존했다. 이는 제품 실패가 아니다. 확인식을 고친 뒤 환경과 나머지 검사가 통과했다. 증거 대조의 첫 바이트 동일성 assert도 checkout 개행 차이로 실패했으며 LF 정규화 대조로 구분했다.

## 보존한 동작과 경계

| 대상 | 독립 확인 |
|---|---|
| `connectionResponse` | deviceAPI/connectAPI 두 호출자를 읽었다. time.Time retry가 transaction 오류보다 먼저다. 명령 error map이 transaction 오류보다 먼저다. 일반 transaction 오류와 성공 응답도 유지한다. 429 JSON과 Retry-After의 반올림이 일치한다. |
| `humanAgents` | owner 필터·철회 제외·ID 순 정렬을 유지한다. 키 없는 agent는 목록에서 연결 필요로 보인다. friend/new의 현재 친구·block·상대 회원 자격 검사는 helper 호출 전에 유지된다. 차단/해제·키/회원 철회 gate 회귀도 통과했다. |
| `validText` | 일반·trial 송신과 두 수신 경로가 공유 함수를 호출한다. 4096 bytes ASCII/Korean/emoji는 허용한다. 초과·빈 값·ASCII/Unicode 공백·고립/삽입 surrogate는 거부한다. 정상 surrogate pair는 허용한다. 앞뒤 공백도 원래 byte 한도에 포함한다. 명시 confirm=false는 relay 요청 없이 거부한다. |
| MCP 포장 | SDK stdio 실제 bundle의 도구 9개 discovery, text content JSON의 state/transport/actualConnection/webhook/evidenceFetch, isError와 기본 held를 검사했다. invalid config와 redirect를 거부하고 held에서 network/key read를 하지 않는다. |
| 서명·Inbox·ACK | frozen signing 회귀가 통과했다. owner wire는 agent 서명을 만들지 않는다. recipient/schema/surrogate 실패는 persist/ACK를 하지 않는다. 정상 owner text는 local store 이후 persist/ACK를 한다. 공유 helper 밖의 서명·저장 순서는 diff가 없다. |
| API·CSRF·TTL·Events·UI | 변경 제품 파일 11개의 diff를 확인했다. handler 외부 API·CSRF·TTL·MCP Events 계약, template/CSS/문구, schema, 의존성·배포 설정의 변경은 없다. 기존 exact SHA 전체 테스트를 재사용하고 frozen/template 회귀를 새로 실행했다. |

## 패킷 보완과 한계

packet 183개 경로와 189개 path/category 쌍의 정본 존재·검색 근거를 [packet-audit.json](packet-audit.json)에 남겼다. 큰 문서·코드의 unknown과 documents-find 잔여 후보는 일반 `rg`로 보완했다. 추가 보완 180개 경로는 모두 존재한다. 과거 기록은 시점과 관련 검색 구간을 확인했으며 전문 재독·현재 SHA의 새 실행을 주장하지 않는다. packet direct_edit는 DEV 변경 seed이므로 tester는 제품을 편집하지 않았다. 모든 쌍의 처리 결과는 [packet outcomes](../../jev/SAR-REFACTOR-001-TESTER-packet-outcomes.json)에 있다.

Windows Go race는 cgo 부재 때문에 실행하지 않았다. DEV Node24/npm11 결과는 보조 근거이며 최종 지원환경 증거로 쓰지 않았다. Docker/실제 SQL migration·운영 배포·실계정 메시지·Google/외부 Bot 왕복·부하·백업 복원·장시간 wake는 not_run이다. 관련 계약이 바뀌지 않은 이번 lite 리팩토링 범위를 넘는다. UI 변화가 없어 새 직접 시각 판정은 해당 없음이다. 전체 제품 HP-16 수락과 과거 held 항목을 이번 QA로 해제하지 않는다.

새 코드·테스트 파일은 만들지 않았다. 실행용 코드는 stdin으로만 전달했고 재현 명령을 JSON에 보존했다. 따라서 이번 문서·로그 변경에 새 제품 등록 게이트는 적용하지 않는다. 지시대로 exact 후보의 기존 전체 게이트를 재사용했다. `git diff --check`와 완료 패킷 검사를 수행한다. 최종 기록 커밋은 worker_done의 SHA로 식별한다. bundle·clone·기존 증거와 이번 로그를 모두 보존한다. 후속 fixed-SHA 리뷰·main 통합은 coordinator가 담당한다.


## 2026-10-13 후속: text matrix escape 재현성

판정은 기록 오류 없음과 실제 검증 값 차이 없음이다. 따라서 commands.json을 정정하지 않았다. 원 commands.json·text-matrix.log·환경 프로브 실패 environment.log는 그대로 보존했다. 고정 d8e0904의 Git 원본 SHA256도 [escape-audit.json](escape-audit.json)에 남겼다.

이중 역슬래시는 JSON 문자열의 정상 escape다. JSON을 한 번 디코딩하면 JS 소스의 각 escape는 역슬래시 한 개다. 예를 들어 JSON 파일의 `\\ud800`은 디코딩 뒤 JS 소스 `\ud800`이 된다. JS가 이 소스를 읽으면 문자열 값은 UTF-16 code unit `55296` 하나다. literal backslash 문자는 값에 남지 않는다. JSON 파일의 command 값을 따옴표째 셸에 복사하는 방식은 재현 절차가 아니다.

원 tester 세션 `01a13534-4add-7a22-bb3f-5e40c8addb09`의 실제 도구 호출 `call_b99a82b25ff94921b4f0bdf1d53c650f`를 읽었다. 도구 cmd 문자열을 디코딩한 뒤 Python AST에서 node/checks/full 대입만 복원했다. 전체 원 스크립트나 SSH 접속 코드는 다시 실행하지 않았다. 복원한 SSH exec command와 commands.json에서 한 번 디코딩한 command는 UTF-8 bytes까지 같다. 두 명령의 SHA256은 `faf7b9403ee59ad31d01ab89ccdfc23e8564e7213e6efb54eba8de792fce69ba`다.

원 SSH helper는 이 명령을 client.exec_command에 그대로 전달한다. 명령의 `<<'NODE'`는 Node stdin의 shell backslash 확장을 막는다. 원 SSH 패킷을 캡처했다고 주장하지 않는다. 도구 원문·원 helper·원 JSON·원 실행 로그를 대조한 결과다.

원 cases 선언을 변경 없이 추출했다. 로컬 Node에서 15개 값의 string 타입·boolean 기대 타입·UTF-16 길이·UTF-8 bytes를 독립 assert했다. 공백과 surrogate의 code units도 assert했다. 결과는 EXIT=0이다. 원 명령·stdin·stdout·stderr·Node 버전은 escape-audit.json에 있다.

| 원 case | 확인한 실제 값 | UTF-8 bytes |
|---|---|---:|
| space | code units [32, 10, 9], trim 결과 빈 문자열 | 3 |
| Unicode space | code units [160, 12288], trim 결과 빈 문자열 | 5 |
| high surrogate | code unit [55296] | 3 |
| low surrogate | code unit [57343] | 3 |
| embedded surrogate | code units [111, 107, 55296] | 5 |
| valid pair | code units [55357, 56832], code point 128512 | 4 |

surrogate 단독 값의 3 bytes는 Node UTF-8 인코딩 시 replacement 문자가 차지하는 길이다. 문자열 값에 surrogate가 실제로 존재하는지는 code unit assert로 별도 확인했다. ASCII/Korean/emoji 최대값은 각각 4096 bytes이고 초과값은 4097/4098/4097 bytes다. 앞뒤 공백을 포함한 최대값과 초과값은 4096/4097 bytes다.

재현할 때 Python json.loads로 commands.json을 읽고 name이 text-matrix인 항목의 command 값을 기존 SSH helper run에 전달한다. JSON escape를 추가로 해제하거나 raw JSON 문자열을 복사하지 않는다. value/type/bytes assert만 재현하려면 escape-audit.json의 runtime_audit.stdin을 디코딩하여 `node --input-type=module`의 stdin으로 전달한다. 두 재현 경로를 구분한다.

이번 로컬 검사는 JS 값 의미만 확인했다. 지원 환경의 validText 동작 증거는 기존 exact cd4dae8 Linux Node22 실행 EXIT=0을 재사용한다. 값 차이가 없어 Linux matrix와 전체 제품 게이트를 반복하지 않았다. 새 제품·테스트 코드 파일은 없으므로 코드 변경용 FullOps 등록 게이트는 해당 없음이다. 제품·의존성·운영·실계정·삭제·이동 변경은 없다. 새 critical/high 발견은 0건이다.

후속 세션은 `01a13541-de34-7933-af11-5e299eba76af`다. Task는 `task_36fa008bd9f6`, Dispatch는 `ctx_05040925a2d9`다. 최종 delta의 독립 리뷰와 main 반영은 coordinator가 맡는다.

후속 완료 계약을 위해 새 attempt의 탐색 패킷을 생성했다. 279개 경로·280개 쌍을 경로·해시·키워드 줄로 처리했다. 과거 attempt producer 4개는 현재 근거로 쓰지 않았다. 큰 과거 문서는 관련 줄만 수동 검색했고 전문 재독을 주장하지 않는다. 원 packet과 outcomes는 jev/history에 보존했다. 상세는 현재 packet-outcomes.json의 unknown_review에 있다.
