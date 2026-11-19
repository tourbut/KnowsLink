---
title: SAR-MVP-002-INSTALL-FIX-DEV-TESTER — 설치 수정 후보의 실제 Grok CLI 독립 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-INSTALL-FIX-DEV-TESTER]
summary: 후보 5506d64의 실제 Grok CLI 설치와 Node 경계 및 ZIP hash를 독립 검증한다
---

# SAR-MVP-002-INSTALL-FIX-DEV-TESTER — 설치 수정 후보의 실제 Grok CLI 독립 QA

## 판정

판정 후보는 `5506d646c47c2d64b35d1254ddfdec5e2003084d`다. 기준은 `376981441d1b0650897847273e1773f8473949e5`다.
실행 위치는 기록 체크아웃 밖의 detached clone이다. clone HEAD는 시작과 끝이 같다. clone의 추적 파일 diff는 없다.
ZIP SHA256은 `fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75`다. 두 번의 `make plugin`과 `make verify-grok-plugin`이 같은 값을 출력했다. DEV 기록과 같다.
`make install`, `make plugin`, `make verify-grok-plugin`의 자기 종료코드는 0이다.
문서 명령을 임시 HOME에서 따로 실행했다. validate, `install --trust`, list, doctor가 통과했다. MCP 서버는 1개다. 도구는 2개다.
설치본 `knowslink_status`는 `held`다. stdout의 내용 있는 줄은 JSON-RPC다. stderr는 비어 있다.
Node 20.19.2의 `npm ci` 종료코드는 1이고 코드는 `EBADENGINE`이다. Node 22.22.2 빌드의 종료코드는 0이다.
새 결함은 없다. 새 critical/high는 없다.
실제 Bot 계정, 앱 동적 카탈로그, hosted Node, 유료 inference, 실제 relay는 미검증이다.
제품 코드와 준비 리뷰 파일은 수정하지 않았다. 증거는 [SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/)다.

## 충돌 후보

Jev context는 [SAR-MVP-002-INSTALL-FIX-DEV.md](../../exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md)를 conflict로 표시했다.
지시서는 이전 QA가 Grok 설치를 실행하지 않았고, DEV 기록의 주장을 독립 검증하라고 적는다.
DEV 기록은 이전 QA의 미실행과 DEV의 임시 HOME 설치를 서로 다른 실행으로 적는다. 제품 규칙 충돌은 없다. ask하지 않았다.
이번 실행은 DEV의 hash, Node 20 exit 1, Node 22 exit 0, doctor healthy, 도구 2개, 설치본 held를 다시 확인했다.

## 기준

날짜는 2026-10-04다. 기록 브랜치는 `fullops/tester`다. 시작 HEAD는 `b51d1d13ce11ebcd2f6e7abfdd4a8fe9d852ffff`다.
공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다. `fullops-test`와 문서 작성 규칙을 적용했다.
Jev find, documents-find, context, route, model-tester를 읽었다. route는 implementation/tester다. keep 문서를 읽었다.
민감 원문은 외부로 보내지 않았다. 통과 조건은 명령 종료코드와 도구 결과로 판정한다.

## 환경

| 항목 | 값 |
|---|---|
| grok | 1.0.46 (`2765805b9442`) |
| Node 22 | v22.22.2, npm 10.9.7. 공식 SHASUMS256 확인 뒤 시험 HOME에 풀었다 |
| Node 20 | v20.19.2, npm 10.8.2. 공식 SHASUMS256 확인 뒤 다른 시험 HOME에 풀었다 |
| Python | 3.12.3 |
| Make | GNU Make 4.3 |
| Go | go1.27.1 |
| clone | `/tmp` 아래 detached `5506d646c47c2d64b35d1254ddfdec5e2003084d` |

`grok plugin install --help`는 `--trust`를 즉시 신뢰 옵션으로 보여 준다.
문서의 Node 22 파이프는 curl, grep, `sha256sum -c`의 `PIPESTATUS`가 0 0 0이다. 출력은 `node-v22.22.2-linux-x64.tar.xz: OK`다.
받은 디렉터리의 `sha256sum -c` 종료코드는 0이다. Node 20 tarball의 `sha256sum -c` 종료코드도 0이다.
빌드와 doctor의 `node`는 그 시험 HOME의 `v22.22.2`다. `go mod download`와 `npm ci`의 cache는 시험 영역이다.

## 패키지

ZIP 항목은 11개다. timestamp는 1980-01-01이다.
항목은 두 marketplace 색인, Cursor manifest, `mcp.json`, Grok manifest, `.mcp.json`, README, skill, `dist/plugin.js`, package.json, `THIRD_PARTY_NOTICES.txt`다.
`node_modules`, fixture, `.env`, PEM, secret, credential 경로는 없다.

## Grok 설치

임시 HOME과 빈 cwd에서 문서 명령을 실행했다. `make verify-grok-plugin`의 내부 실행과 별도다.

| 명령 | 종료코드 | 관찰 |
|---|---|---|
| `grok plugin validate` | 0 | `Plugin manifest is valid.`, name `knowslink`, version `0.1.0`, MCP servers |
| `grok plugin install <경로> --trust` | 0 | `Installed 1 plugin(s) from ...: knowslink` |
| `grok plugin list --json` | 0 | name `knowslink`, version `0.1.0` |
| `grok mcp doctor knowslink --json` | 0 | `healthy_count` 1, command found가 시험 Node 22.22.2, handshake `2025-11-25`, `2 tools discovered` |

설치 경로는 임시 HOME의 `.grok/installed-plugins/knowslink-07ec00a3`다.
준비 폴더를 삭제한 뒤 doctor 종료코드는 0이고 `healthy_count`는 1이다.
`--trust`를 뺀 install의 종료코드는 1이다. stderr는 `--trust`로 다시 실행하라고 적는다. list는 `[]`다. doctor는 `MCP server 'knowslink' not found.`로 종료코드 1이다. MCP는 켜지지 않았다.
`grok plugin marketplace add <폴더>` 종료코드는 0이다. 이어서 `grok plugin install knowslink --trust` 종료코드는 0이다. list의 marketplace는 `extract-market`다. doctor `healthy_count`는 1이다.
`grok plugin uninstall knowslink` 종료코드는 0이다. 이후 list는 `[]`다.

설치본 `dist/plugin.js`에 줄 단위 JSON-RPC로 initialize, `tools/list`, `knowslink_status`를 보냈다. 환경은 `KNOWSLINK_MODE=held`다.
서버 프로토콜은 `2025-11-25`다. 도구 이름은 `knowslink_pull_once`와 `knowslink_status`다.
status 본문은 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다. `isError`는 false다.
stdout의 내용 있는 줄은 모두 JSON이다. stderr 길이는 0이다. 이 프로브는 `knowslink_pull_once`를 호출하지 않았다.
`make plugin --verify`와 `make verify-grok-plugin`은 같은 빌드의 `mcp.test.js`를 실행했다. PASS 줄은 held와 network 없음이다. 두 명령의 종료코드는 0이다.

## 설정 보존

설치 전 임시 HOME에 marker 줄, `qa-marker.txt`, `installed-plugins/unrelated/keep.txt`를 두었다.
설치 뒤 marker 줄과 두 파일은 그대로다. CLI가 config.toml에 `[plugins] enabled = ["knowslink"]`를 추가했다.
marketplace 경로에서도 marker 줄은 남았다. CLI가 `[[marketplace.sources]]`를 추가했다.
uninstall 뒤 플러그인 list는 비어 있다. config.toml에는 marker, marketplace source, enabled 이름이 남는다.
사용자 `~/.grok/installed-plugins`에 knowslink 이름은 없다.
QA 전후 비교에서 크기 또는 mtime이 바뀐 사용자 파일은 `models_cache.json`뿐이다. 크기는 9617이다.
임시 HOME에는 `models_cache.json`이 없다. QA용 grok이 끝난 뒤에도 그 사용자 파일의 mtime은 갱신됐다. 갱신 주체는 이 세션의 Grok 프로세스다. 설치 시험은 사용자 플러그인과 사용자 config를 바꾸지 않았다.

## Node 20

시험 영역에 복사한 adapter에서 Node v20.19.2와 npm 10.8.2로 `npm ci`를 실행했다. 종료코드는 1이다.
stderr 코드는 `EBADENGINE`이다. 요구는 node `>=22.22.2 <23`, npm `>=10.9.7 <11`이다. 실제는 npm `10.8.2`, node `v20.19.2`다.
npm 로그는 시험 cache에 있다.

## 재사용한 증거

| 증거 | SHA | 이번 실행 |
|---|---|---|
| Go race `make test` | `c4ebbecd3ef91be10ecbb517fe451e02769358bc` | 재사용. `552586b`와 후보의 `cmd`, `internal`, `adapters/src` diff 종료코드 0 |
| Chrome owner UI | `9584aafcbb5fee88dcc6d618caf660884f6a527d` | 재사용. HTML, CSS, TSX, JSX, Go 이름 diff가 비어 있다. 새 브라우저 QA로 기록하지 않음 |
| 리뷰 `0055a5b` | `0055a5b992998a919e155db226a03eeb12e08a3f` | 코어 판정만 재사용. 대상은 `552586b`다 |
| 이전 QA 기록 | `0fb32cd45ee77bbe9bdad8629f4bf2bdff6b2264` | Grok 설치는 그 기록이 실행하지 않았다. 설치 판정은 이번 실행이다 |

`adapters/package.json`의 차이는 format 검사 경로에 `.mcp.json`과 `.grok-plugin/plugin.json`을 더한 한 줄이다.
MCP 경계 PASS는 이번 빌드의 `make plugin --verify`와 `make verify-grok-plugin`이 다시 출력했다.

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `make install` | 0 | [make_install.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/make_install.log) |
| `make plugin` | 0 | [make_plugin.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/make_plugin.log) |
| `make plugin` 재실행 | 0 | [make_plugin_again.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/make_plugin_again.log) |
| `make verify-grok-plugin` | 0 | [make_verify_grok.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/make_verify_grok.log) |
| `sha256sum` ZIP | 0 | [zip_hash.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/zip_hash.log) |
| Node 22 공식 파이프 | 0 | [node22_pipe.out](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/node22_pipe.out) |
| Node 20 `npm ci` | 1 | [node20_npm_ci.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/node20_npm_ci.log) |
| validate, install, list, doctor | 0 | [doc_validate.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/doc_validate.log), [doc_install.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/doc_install.log), [doc_list.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/doc_list.log), [doc_doctor.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/doc_doctor.log) |
| 준비 폴더 삭제 후 doctor | 0 | [doc_doctor_after_rm.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/doc_doctor_after_rm.log) |
| trust 생략 install | 1 | [notrust_install.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/notrust_install.log) |
| marketplace add, install, doctor, uninstall | 0 | [market_add.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/market_add.log), [market_install.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/market_install.log), [market_doctor.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/market_doctor.log), [market_uninstall.log](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/market_uninstall.log) |
| 설치본 status | 0 | [status-protocol.json](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/status-protocol.json) |

자기 종료코드 목록은 [exits.txt](SAR-MVP-002-INSTALL-FIX-DEV-TESTER-test/exits.txt)다.

## 한계

실제 Bot 계정 설치, Bot 앱의 동적 도구 카탈로그, 모델이 중간에 있는 새 Bot 세션의 `knowslink_status`, hosted Node는 원격 재시험 전까지 미검증이다.
유료 inference와 실제 relay는 호출하지 않았다. DEC-02와 calendar effect는 held다.
`make verify-grok-plugin`의 held 검사는 로컬 stdio다. 운영 relay에 접속하지 않는다.
uninstall 뒤 config.toml에 enabled 이름과 marketplace source가 남는 동작은 grok 1.0.46에서 관찰했다. 플러그인 list는 비어 있다. KnowsLink 결함으로 판정하지 않았다.
