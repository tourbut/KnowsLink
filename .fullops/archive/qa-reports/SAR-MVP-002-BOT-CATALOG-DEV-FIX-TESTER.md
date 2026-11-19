---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER — 423db6a installer 보존 경계 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER]
summary: "423db6a의 파일 보존, 실패 경계, 새 ZIP, 레포 밖 held를 독립 검증한다"
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER — 423db6a installer 보존 경계 QA

## 판정

판정 후보는 `423db6a2a388ea63610462f9d3a5f4c619dd781b`다. 변경 기준은 `8e46c5a846e6d190e484e48be40b3dc368001a2b`다.
실행 위치는 기록 체크아웃 밖의 detached clone이다. clone HEAD는 시작과 끝이 같다. clone의 추적 파일 diff는 없다.
첫 설치와 재실행의 종료코드는 0이다. ZIP SHA256은 `d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609`다.
등록 Command와 Arguments는 그 prefix의 준비된 파일이다. 레포 밖 `env -i` 프로브는 도구 2개, held, stderr 길이 0이다.
상대 prefix, `node` 사용자 파일, `node` 링크, bundle이 없는 `knowslink` 폴더는 변경 전에 exit 1이다. 내용은 남는다.
bundle 추출 실패와 `env -i` 검증 실패는 exit 1이다. 기존 bundle이 남는다. Ready가 없다. stage가 없다.
swap이 이전 `knowslink`를 stage로 옮긴 뒤의 실패 두 건은 exit 1이고 Ready와 stage가 없다. 이전 `knowslink` 디렉터리도 없다.
이 관찰은 low 결함 1건이다. 새 critical/high는 없다. 제품 코드는 수정하지 않았다.
실제 앱 등록, 개인 UI, 카탈로그, 모델 status 호출, relay, 유료 inference는 미검증이다.
증거는 [SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/)다.

## 기준

날짜는 2026-10-04다. 기록 브랜치는 `fullops/tester`다. 착수 HEAD는 `ebfd5295538b802db1721de71011933ef529a4f5`다.
공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다. `fullops-test`와 문서 작성 규칙을 적용했다.
Jev find, documents-find, context를 읽었다. keep은 FULLOPS, project, 공통 README, FIX 실행 기록이다.
context는 원본 QA `SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md`를 과제와 충돌로 표시한다. 원본 QA의 성공을 이번 파일 보존 성공으로 쓰지 않는다.
구현자의 경계 회귀 19개 스크립트는 실행하지 않았다. 판정은 임시 prefix의 실제 파일과 짧은 shim이다.
로컬 stdio의 held를 앱 카탈로그 성공으로 기록하지 않는다.

## 환경

| 항목 | 값 |
|---|---|
| 호스트 | Linux x86_64 |
| 재사용 Node | v22.22.2, npm 10.9.7. 바이너리 SHA256 `81925c0995b5c1427b5d538e6a90ca2fdc4daffb786b09af749beaf7369d4e90` |
| clone | `/tmp` 아래 detached `423db6a2a388ea63610462f9d3a5f4c619dd781b` |
| prefix | `/tmp` 아래 임시 절대 경로. `/workspace/.knowslink`는 실행 전후에 없다 |

Node 트리는 DEV FIX scratch의 기존 prefix에서 복사했다. 그 트리는 이전 설치에서 tarball checksum OK 뒤에 남은 것이다.
이번 설치 출력에 `tar.gz: OK`가 없다. 공식 SHASUMS 다운로드는 반복하지 않았다.
HOME과 npm cache는 임시 디렉터리다. 사용자 `~/.grok`의 `config.toml`, `trusted_folders.toml`, `installed-plugins` hash는 실행 전후가 같다.
기록 체크아웃의 `git status --porcelain`은 실행 중 비어 있었다.

## 설치와 레포 밖 프로브

첫 실행은 복사한 Node가 `v22.22.2`라 다운로드를 건너뛴다. 그 다음 `npm ci`, build, ZIP 생성, stage 안 `env -i`, `knowslink` 교체, Ready가 이어진다.
ZIP 명령의 종료코드는 0이다. hash는 기대값과 같다. ZIP 목록에 `STALE` 항목이 없다.
Command는 `<prefix>/node/bin/node`다. Arguments는 `<prefix>/knowslink/dist/plugin.js`다. 두 파일은 그 경로에 있다. Node 버전은 `v22.22.2`다.
프로브 cwd는 clone 밖이다. 자식은 `env -i <준비 node> <준비 plugin.js>`다. 환경 키 개수는 0이다.
도구 이름은 `knowslink_pull_once`와 `knowslink_status`다. status 본문은 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다.
`isError`는 false다. stderr 길이는 0이다. 종료코드는 0이다. `knowslink_pull_once`는 호출하지 않았다.
설치 성공 출력의 stderr에는 esbuild 진행 줄이 있다. 그 줄은 레포 밖 프로브의 stderr가 아니다.

## 재실행과 거절

재실행 prefix에는 `package/knowslink/STALE`, 기존 bundle의 `STALE`, `notes.txt`를 넣었다.
종료코드는 0이다. Ready가 있다. 다운로드 checksum 줄은 없다. 새 bundle 아래에 `STALE`이 없다.
`package/knowslink/STALE`과 `notes.txt`는 남는다. stage 디렉터리는 없다. 같은 경로의 레포 밖 프로브 종료코드는 0이다.
상대 prefix `rel/pfx`의 종료코드는 1이다. stderr는 `KNOWSLINK_PREFIX must be an absolute path: rel/pfx`다. `rel` 디렉터리는 생기지 않는다. Ready가 없다.
`node`가 파일 `mine`이면 종료코드는 1이다. 파일 내용이 남고 prefix에는 그 파일만 있다.
`node`가 `node-target`을 가리키는 symlink이면 종료코드는 1이다. 링크와 대상 내용이 남는다. `knowslink`와 stage는 없다.
bundle이 없는 `knowslink` 디렉터리의 `x`는 `mine`으로 남는다. `plugin.js`는 없다. 같이 둔 Node 바이너리 hash는 같다.
세 거절의 stderr는 `refusing to replace`로 시작한다. Ready가 없다.

## swap 이전 실패

추출 shim은 `python3 -m zipfile`만 exit 1로 만든다. `package_plugin.py`는 실제 Python으로 통과한다.
종료코드는 1이다. stderr에 `shim: extract failed`가 있다. 기존 `OLD-BUNDLE`과 `plugin.js` hash가 남는다. `notes.txt`가 남는다. Ready와 stage가 없다.
`env` shim은 첫 인자가 `-i`일 때만 exit 1이다. stderr에 `shim: env -i check failed`가 있다.
종료코드는 1이다. 기존 bundle hash가 남는다. Ready와 stage가 없다.

## swap 구간 관찰

FIX 기록은 실패하면 기존 준비물을 보존하고 stage를 지운다고 적는다. 재시험 초안은 실패하면 기존 준비물을 그대로 둔다고 적는다.
스크립트는 기존 디렉터리를 `stage/old-*`로 옮긴 다음 새 디렉터리를 prefix로 옮긴다. EXIT trap은 stage를 지운다. TERM은 `exit 1`이다.

두 번째 `mv`가 `shim: swap mv failed`로 exit 1인 실행:

| 항목 | 관찰 |
|---|---|
| 종료코드 | 1 |
| Ready | 없음 |
| stage | 없음 |
| `knowslink` | 없음. `plugin.js`는 MISSING. `OLD-BUNDLE`도 없음 |
| `node` 바이너리 hash | 실행 전과 같음 |
| `notes.txt`, `package/knowslink/STALE` | 남음 |

첫 rename 직후 SIGTERM인 실행도 종료코드 1이다. stderr는 `shim: sent SIGTERM`이다. Ready와 stage가 없다. `knowslink`와 `OLD-BUNDLE`이 없다. `node` hash와 `notes.txt`는 남는다.
보존 주장에 대한 검사 3개가 실패했다. 목록은 [asserts.txt](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/asserts.txt)다. `qa_fail`은 3이다.
이 구간에서 이전 `knowslink`는 보존되지 않는다. 추출 실패와 `env -i` 실패의 보존 결과와 다르다.
등급은 low다. 영향은 그 짧은 rename 구간의 이전 bundle 소실이다. prefix의 다른 파일과 Node는 남는다. 다시 설치하면 준비할 수 있다.
새 critical/high는 아니다. 제품 코드는 고치지 않았다. 병합 판단은 coor가 한다.

## 재사용한 증거

원본 기록은 `87cfb7c0f2361e450df1ec82190185b8a0d2feb2`다. 그 판정 후보는 `8e46c5a`다.

| 증거 | 이번 비교 | 사용 |
|---|---|---|
| Node 20 `EBADENGINE` | `adapters/package.json`, `adapters/.npmrc`, `adapters/src` diff 종료코드 0 | 재사용. 새 Node 20 실행이 아님 |
| Go race `make test` | `cmd`, `internal`이 위 diff에 포함되고 종료코드 0 | 재사용 |
| `make verify-mvp`와 MCP 경계 기록 | 같은 core diff 종료코드 0 | 재사용. 이번 설치의 MCP PASS는 그 소스의 새 실행이다 |
| Chrome owner UI | HTML, CSS, TSX, JSX 이름 diff가 비어 있음 | 재사용. 새 브라우저 QA가 아님 |
| 원본 설치 성공, 실패 전파, ZIP `b7882df7` | 제품 diff는 `adapters/README.md`와 `scripts/install_bot_mcp.sh` | 이번 보존 경계의 결과로 쓰지 않음 |
| 실제 Grok CLI 설치 | 이전 ZIP `fb745c66` | 이 ZIP의 `grok plugin install`과 `mcp doctor`로 쓰지 않음 |

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| 첫 설치 | 0 | [install-first.out](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/install-first.out) |
| ZIP `sha256sum` | 0 | [zip.sha256](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/zip.sha256) |
| 레포 밖 `env -i` 프로브 | 0 | [probe-first.json](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/probe-first.json) |
| STALE 재실행 | 0 | [install-rerun.out](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/install-rerun.out) |
| 재실행 프로브 | 0 | [probe-rerun.json](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/probe-rerun.json) |
| 상대 prefix | 1 | [relative.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/relative.err) |
| `node` 사용자 파일 | 1 | [nodefile.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/nodefile.err) |
| `node` symlink | 1 | [nodelink.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/nodelink.err) |
| 비소유 `knowslink` | 1 | [kdir.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/kdir.err) |
| bundle 추출 실패 | 1 | [badextract.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/badextract.err) |
| `env -i` 검증 실패 | 1 | [badcheck.err](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/badcheck.err) |
| swap 두 번째 `mv` 실패 | 1 | [swapfail.observe](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/swapfail.observe) |
| swap 첫 rename 뒤 SIGTERM | 1 | [swapterm.observe](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/swapterm.observe) |

자기 종료코드 목록은 [exits.txt](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/exits.txt)다.
재현 스크립트는 [run.sh](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/run.sh)와 [probe.mjs](SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER-test/probe.mjs)다.

## 한계

실제 Grok Bot 계정의 custom MCP 등록, 개인 Setup 화면, 승인 카드, Installed 목록, 새 대화의 `knowslink_status`는 미검증이다.
실제 aarch64 실행, 컴퓨터 update·Reset 뒤 `/workspace` 유지, `grok plugin install`, `mcp doctor`는 실행하지 않았다.
유료 inference와 실제 relay는 호출하지 않았다. `knowslink_pull_once`는 호출하지 않았다.
이번 Node 바이너리가 공식 tarball과 같은지는 이번 실행에서 다시 확인하지 않았다. 버전과 기존 checksum OK prefix의 재사용만 기록한다.
