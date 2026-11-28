---
title: SAR-ORCA-HOST-001-DEV 검증 보고
status: review
updated: 2026-11-28
owner: dev
tasks: [SAR-ORCA-HOST-001-DEV]
summary: Orca 봇 고정 설치·자동 시작 watch·stdio MCP 한 줄·웹 항목을 구현하고 표적 테스트와 lint를 통과했다
---

# SAR-ORCA-HOST-001-DEV 검증 보고

## 구현

| 항목 | 내용 |
|---|---|
| 고정 위치 설치 | `plugin.js orca install <연결 폴더> --bot <이름> [--dir <경로>]`. 기본 `%LOCALAPPDATA%\KnowsLink\orca\`에 `dist\plugin.js`, 실행 중인 `node.exe`, `defaults.json`. 다시 실행하면 덮어쓴다. 같은 크기의 `node.exe`는 복사하지 않는다(봇 MCP가 잠근다). |
| 자동 시작 | HKCU `...\CurrentVersion\Run`의 `KnowsLinkOrcaWatch` = `<설치>\node.exe <설치>\dist\plugin.js orca start`. 관리자 권한 없음. `orca start`는 watch를 detached·창 없이 띄우고 출력은 `<설치>\watch.log`. |
| 중복 방지·끄기 | watch가 `<폴더>\orca-watch.lock`(pid·시작 시각)을 잡는다. 살아 있는 holder가 있으면 `orca watch already running`. 이번 부팅 전 잠금·죽은 pid는 버린다. `orca stop`, `orca autostart off`(Run 값 삭제+stop), `orca autostart on`. |
| stdio MCP 한 줄 | 설치 출력 `toolsAdd`: `orca bot tools add "<설치>\node.exe <설치>\dist\plugin.js" --name knowslink`. 값은 `defaults.json`이 넣는다(`KNOWSLINK_MODE=public-node`, `RELAY_URL`=폴더의 `agent.json` 값, `KNOWSLINK_HOST=orca`, `KNOWSLINK_AGENT_FOLDER`, `KNOWSLINK_ORCA_BOT`, `KNOWSLINK_AUTO_RECEIVE=off`). |
| 수신자 하나 | 추천안. watch만 자동 수신한다. MCP는 `KNOWSLINK_AUTO_RECEIVE=off`라 loop가 없고, `knowslink_text_receive`는 같은 Inbox에서 읽고 비었을 때만 요청받은 1건을 당긴다(`orca receive`와 같은 경로). 이유: 도어벨은 watch만 울리므로 두 loop가 돌면 MCP가 먼저 가져간 메시지는 도어벨 없이 Inbox에 놓인다. |
| 웹 | '어디에서 사용할까요'에 `Orca 봇 (내 컴퓨터)`(ko·en). Claude Code 항목과 같은 `dl.kv` 구조, 새 스타일 없음. |

## 검증

| 검사 | 결과 |
|---|---|
| `node dist/orca.test.js` | PASS: 기존 S3–S5 + defaults 키가 `DEFAULT_KEYS` 안·비밀 키워드 없음·`AUTO_RECEIVE=off`, `installedCommand` 따옴표, `autostartArgs` add/delete argv, `claimWatch` 두 번째 거부·부팅 전 잠금 재사용 |
| `connect-simple.test.js`, `mcp.test.js` | PASS(defaults 변경 회귀) |
| reg.exe 인용 | 공백 경로 값을 임시 이름 `KnowsLinkQuoteProbe`로 add→query→delete. 저장 값 `"C:\Program Files\x\node.exe" C:\a\plugin.js orca start` 그대로 |
| Go | `TestAutosetupGuideOnlyForClaudeAndCodex`(orca ko·en 문구·명령 박스 없음), 선택지 순서, i18n 3종 PASS |
| `python3 scripts/make.py lint` | 통과 |
| `fullops lint --from 0899e5bd` | 통과(ERROR 0). WARNING: `orca.ts` 464줄, `i18n_en_human.go` 312줄, PLANS.md 크기, product-test 시간 |

## 이 PC 적용

lead 승인 `msg_adfdbb26e504`(2026-11-28). 순서: 기존 watch 중지 → install → 확인 → 쿠어용 한 줄 ask(`msg_ce655a143a82`).

| 항목 | 결과 |
|---|---|
| 기존 watch 중지 | ops 빌드 `node.exe .../ops/adapters/dist/plugin.js orca watch .../orca-agent-20261128 --bot 쿠어` PID 74944 하나만 있었고 `Stop-Process`로 끝냄 |
| install | Node 22.22.2로 `bb252e31` 빌드 실행, 종료 0. 출력 `installed` `C:\Users\shin\AppData\Local\KnowsLink\orca`, `agent` `agent_5513ec51367fed3748903d`, `watch` `{"state":"started","pid":69852}` |
| 설치 파일 | `dist\plugin.js` 1692629 B, `node.exe` 87074816 B, `defaults.json` 274 B(모드·relay URL·폴더·봇 이름·`AUTO_RECEIVE=off`만), `watch.log` |
| 자동 시작 등록 | `reg query`: `KnowsLinkOrcaWatch REG_SZ C:\Users\shin\AppData\Local\KnowsLink\orca\node.exe C:\Users\shin\AppData\Local\KnowsLink\orca\dist\plugin.js orca start` |
| 새 watch | PID 69852 argv `...\KnowsLink\orca\node.exe ...\orca\dist\plugin.js orca watch ...\orca-agent-20261128 --bot 쿠어`, 잠금 `{"pid":69852,...}`. 15초 동안 30회 표본 중 11회 `104.21.54.184:443`(`link.knowslog.com`) 연결. `watch.log`는 수신 때만 쓰므로 비어 있다 |
| 중복 방지 | 설치된 `orca start` 재실행 → `{"state":"running","pid":69852}`, 새 프로세스 없음 |
| MCP | 설치된 node+bundle을 `KNOWSLINK_*` 환경 없이 stdio로 띄워 `knowslink_status`: `me.agent` `agent_5513ec51367fed3748903d`(이름 쿠어), `autoReceive.state` `off` |
| agent·grant | 새 agent 없음(install·MCP 둘 다 같은 ID). install은 grant 경로를 부르지 않는다. 서버의 grant 3개는 이 worker가 다시 읽지 않았고 쿠어 확인을 요청했다 |

미확인: 쿠어 봇의 `tools add` 결과, 재부팅 뒤 자동 시작(사람 단계).

## 스스로 정한 선택과 되돌리는 방법

- `.cmd` 래퍼 대신 `node.exe`+`plugin.js` 직접 실행과 `defaults.json`. cmd.exe는 인용 인자를 망가뜨리고 종료 때 배치 확인을 묻는다. 되돌리려면 `orcaDefaults`의 값을 래퍼의 `set`으로 옮긴다.
- 자동 시작은 작업 스케줄러 대신 HKCU Run. `schtasks /SC ONLOGON`은 비관리자에게 막힐 수 있다. 로그온 때 콘솔 창이 잠깐 보인다.
- 설치 순서는 연결(setup-connect) → install. install이 연결 폴더의 relay URL과 agent ID를 읽는다.
- `node.exe`를 설치 폴더에 복사한다(.devtools·워크트리 정리에 영향받지 않게).

## 못 한 것과 한계

- macOS·Linux 자동 시작 없음.
- Orca의 `tools add`가 공백 경로의 따옴표를 어떻게 나누는지 확인하지 않았다(이 PC 경로에는 공백이 없다).
- 같은 부팅 안 pid 재사용 시 잠금을 살아 있다고 볼 수 있다. `watch.log`는 순환하지 않는다.
- 로컬 MCP `knowslink_text_send`는 고정 형식(`kind`)을 받지 않는다. 고정 형식은 지금처럼 `orca send` CLI다.
