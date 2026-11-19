---
title: SAR-MVP-002-BOT-CATALOG-DEV-TESTER — 실제 Bot Command 등록 후보의 독립 설치 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-TESTER]
summary: 후보 8e46c5a의 Command 설치와 실패 전파를 독립 검증한다
---

# SAR-MVP-002-BOT-CATALOG-DEV-TESTER — 실제 Bot Command 등록 후보의 독립 설치 QA

## 판정

판정 후보는 `8e46c5a846e6d190e484e48be40b3dc368001a2b`다. 기준은 `8c95bde1240bf52819ce987b2d387d437610c6e8`다.
실행 위치는 기록 체크아웃 밖의 detached clone이다. clone HEAD는 시작과 끝이 같다. clone의 추적 파일 diff는 없다.
`sh scripts/install_bot_mcp.sh`의 첫 실행과 재실행 종료코드는 0이다. 임시 `KNOWSLINK_PREFIX`만 사용했다.
ZIP SHA256은 `b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad`다. 설치 출력이 같은 값을 찍었다.
준비된 Node는 `v22.22.2`다. 그 Node의 npm은 `10.9.7`이다. 등록 Command와 Arguments는 그 prefix의 절대 경로다.
레포 밖 cwd에서 `env -i`로 그 Node와 bundle을 직접 실행했다. 도구는 2개다. `knowslink_status`는 held다. stderr 길이는 0이다.
미지원 `uname`, 다운로드 실패, checksum 실패, 빌드 실패의 종료코드는 1이다. 그 출력에 Ready가 없다.
`verify_grok_plugin.py` 실패 shim의 종료코드는 1이다. grok stdout와 stderr가 보인다. `GROK_CONFIG` 계열 4개 이름은 자식 환경에 없다. 실패 출력에 PASS는 없다.
새 결함은 없다. 새 critical/high는 없다.
실제 앱 계정 등록, 승인 카드의 Command 제공, 앱 카탈로그 노출, Bot 세션 호출, 실제 aarch64 실행은 미검증이다.
제품 코드와 준비 리뷰 파일은 수정하지 않았다. 증거는 [SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/)다.

## 기준

날짜는 2026-10-04다. 기록 브랜치는 `fullops/tester`다. 착수 HEAD는 `10949a41356907df8f20c5e5a4b3ce024f427ce7`다.
공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다. `fullops-test`와 문서 작성 규칙을 적용했다.
Jev find, documents-find, context, route를 읽었다. route는 implementation/tester다. keep 문서를 읽었다.
민감 원문은 외부로 보내지 않았다. 통과 조건은 명령 종료코드와 출력 값으로 판정한다.
로컬 stdio의 held를 앱 카탈로그 성공으로 기록하지 않는다.

## 환경

| 항목 | 값 |
|---|---|
| 호스트 | Linux x86_64 |
| 기록 체크아웃 Node | v22.22.2, npm 10.9.7 |
| 준비 Node | prefix의 `node/bin/node` v22.22.2, npm 10.9.7 |
| Python | 3.12.3 |
| clone | `/tmp` 아래 detached `8e46c5a846e6d190e484e48be40b3dc368001a2b` |
| prefix | `/tmp` 아래 임시 디렉터리. `/workspace/.knowslink`는 실행 전후에 없다 |

npm cache는 임시 디렉터리다. 사용자 `~/.grok`의 `config.toml`, `trusted_folders.toml`, `installed-plugins` 내용 hash는 실행 전후가 같다. 그 목록에 knowslink 이름은 없다.

## 설치와 재실행

첫 실행 cwd는 clone 루트다. PATH 앞의 `xz`는 호출되면 exit 1을 기록한다. sentinel은 비어 있고 설치 종료코드는 0이다.
공식 `SHASUMS256.txt`의 x64 값과 arm64 값은 스크립트 pin과 같다. 스크립트 checksum 줄은 `node-v22.22.2-linux-x64.tar.gz: OK`다.
같은 tarball을 다시 받아 `sha256sum -c`를 실행했다. 종료코드는 0이다.
패키지 검증 줄은 `extracted standalone MCP test exit=0`이다. MCP PASS 줄은 2개다. 하나는 패키지 검증이고 하나는 `env -i` 경계 검사다.
등록 값은 아래와 같다.

```text
Name:        knowslink
Type:        Command
Command:     <prefix>/node/bin/node
Arguments:   <prefix>/knowslink/dist/plugin.js
Environment: none (held is the default)
```

두 경로는 절대 경로이고 실행 파일이 그 자리에 있다.
prefix의 `unrelated/keep.txt`와 `user-settings/keep.txt` hash는 첫 실행 뒤와 재실행 뒤가 같다.
재실행은 PATH 앞의 curl이 호출되면 exit 99다. curl sentinel은 비어 있다. checksum OK 줄은 없다. 종료코드는 0이다. Command와 Arguments는 같다. bundle mtime은 바뀌었다.

## 독립 Command 프로브

프로브 cwd는 clone 밖이다. 자식은 `env -i <준비 node> <준비 plugin.js>`다. `/proc` 환경 키 개수는 0이다.
보낸 호출은 initialize, `tools/list`, `knowslink_status`다. `knowslink_pull_once`는 호출하지 않았다.
도구 이름은 `knowslink_pull_once`와 `knowslink_status`다. status 본문은 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다. `isError`는 false다.
stdout의 내용 있는 줄은 JSON이다. stderr 길이는 0이다. 프로세스 종료코드는 0이다. 프로토콜은 `2025-11-25`다.
이 결과는 로컬 stdio다. 앱 등록 성공이 아니다.

## 실패 전파

| 조건 | 종료코드 | Ready | 추가 관찰 |
|---|---|---|---|
| `uname` Darwin/arm64 | 1 | 없음 | stderr `unsupported platform: Darwin arm64`. prefix 경로가 생기지 않음 |
| curl exit 1 | 1 | 없음 | `node` 디렉터리가 생기지 않음 |
| tarball 내용이 pin과 다름 | 1 | 없음 | `sha256sum` FAILED가 보임 |
| 복사본 `mcp.ts` 문법 오류 | 1 | 없음 | `tsc` TS1109. 제품 트리는 바꾸지 않음 |

실제 aarch64 호스트는 없다. arm64 Node 바이너리는 실행하지 않았다. arm64 pin이 공식 checksum과 같은 것만 확인했다.

## verify 스크립트

후보의 `scripts/verify_grok_plugin.py`를 grok shim으로 실행했다. shim은 stdout 표식, stderr 표식, exit 7이다.
부모에는 `GROK_HOME`, `GROK_CONFIG`, `GROK_CONFIG_PATH`, `GROK_MANAGED_CONFIG_URL`을 넣었다. 자식 환경 이름에는 그 4개가 없다.
stderr는 `grok plugin validate ... exit=7`과 두 표식을 포함한다. PASS 문구는 없다.
소스의 PASS 문장은 `installed bundle held under this node`다. 이전 문장 `installed copy held`는 없다.
이 실행은 실패 경로다. 이 ZIP의 `grok plugin install`과 `grok mcp doctor`는 다시 실행하지 않았다.

## 재사용한 증거

| 증거 | 대상 SHA | 이번 후보와의 비교 | 이번 실행 |
|---|---|---|---|
| Node 20 `EBADENGINE` | `5506d646c47c2d64b35d1254ddfdec5e2003084d` | `adapters/package.json`, `adapters/.npmrc`, `adapters/src` diff 종료코드 0 | 재사용. 새 Node 20 실행이 아님 |
| Go race `make test` | `c4ebbecd3ef91be10ecbb517fe451e02769358bc` | `cmd`, `internal` diff 종료코드 0 | 재사용 |
| `make verify-mvp`와 MCP 경계 기록 | `552586b6e886f95bffa9a000a031ea03070afedb` | `cmd`, `internal`, `adapters/src` diff 종료코드 0 | 재사용. 이번 설치의 MCP PASS는 그 소스의 새 실행이다 |
| Chrome owner UI | `9584aafcbb5fee88dcc6d618caf660884f6a527d` | HTML, CSS, TSX, JSX, `cmd`, `internal` 이름 diff가 비어 있음 | 재사용. 새 브라우저 QA가 아님 |
| 실제 Grok CLI 설치 | `5506d646c47c2d64b35d1254ddfdec5e2003084d` | 플러그인 런타임 소스는 같다. ZIP과 `verify_grok_plugin.py`와 `adapters/README.md`는 다르다 | doctor를 이 ZIP의 결과로 기록하지 않음. 전체 CLI 설치 QA는 반복하지 않음 |

`5506d64` 대비 이름 변경은 `Makefile`, `adapters/README.md`, `scripts/install_bot_mcp.sh`, `scripts/verify_grok_plugin.py`다.

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `sh -n scripts/install_bot_mcp.sh` | 0 | [sh_n.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/sh_n.out) |
| 공식 SHASUMS256 다운로드 | 0 | [shasums.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/shasums.out) |
| 첫 설치 | 0 | [install_first.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_first.out) |
| ZIP `sha256sum` | 0 | [zip_hash.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/zip_hash.out) |
| x64 tarball `sha256sum -c` | 0 | [tarball_check.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/tarball_check.out) |
| 독립 `env -i` 프로브 | 0 | [probe.json](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/probe.json) |
| 재실행 | 0 | [install_rerun.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_rerun.out) |
| 미지원 uname | 1 | [install_unsupported.err](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_unsupported.err) |
| 다운로드 실패 | 1 | [install_download_fail.err](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_download_fail.err) |
| checksum 실패 | 1 | [install_checksum_fail.err](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_checksum_fail.err) |
| 빌드 실패 | 1 | [install_build_fail.out](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/install_build_fail.out) |
| verify 실패 shim | 1 | [verify_shim.err](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/verify_shim.err) |

자기 종료코드 목록은 [exits.txt](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/exits.txt)다. 첫 `1 install_first`의 설명은 [harness-note.txt](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/harness-note.txt)다.
재현 스크립트는 [run.sh](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/run.sh)와 [probe.mjs](SAR-MVP-002-BOT-CATALOG-DEV-TESTER-test/probe.mjs)다.

## 한계

실제 Grok Bot 계정의 custom MCP 등록, 승인 카드의 Command 종류, Installed 목록, 새 대화의 `knowslink_status`는 원격 계정이 없어 미검증이다.
실제 aarch64 기계에서 Node arm64 바이너리를 실행하지 않았다.
유료 inference와 실제 relay는 호출하지 않았다. DEC-02와 calendar effect는 held다.
이 ZIP에 대한 `grok plugin install`과 `mcp doctor`는 실행하지 않았다. 이전 CLI 설치 증거는 ZIP `fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75`다.
