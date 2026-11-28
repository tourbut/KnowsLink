# SAR-ORCA-HOST-001-DEV 리뷰 보고

- base: b8b2bb7f0b9963b0fc1491e30d0b07708c283cb2
- head: 6943ee6de30d12762af2d6e8f192a8ff7a32b574
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-68284280

## 요약

D068 보안 경계 중심 1회 리뷰. 변경 16개 파일을 읽었고 하네스 기록 2개는 건너뛰었다. 표적 테스트 `dist/orca.test.js`를 스냅샷 밖 복사본(Node 22.22.2)에서 다시 돌려 통과했다(schema, doorbell 메타데이터만, untrusted receive, 고정 설치 defaults, autostart argv, 폴더당 watch 하나). 새 통합 테스트는 만들지 않았다.

| 초점 | 판정 |
|---|---|
| ① knowslink_send 우회 | 통과. MCP 입력에 `confirmed`가 없고 `body` 키는 `orcaBody` 허용 목록만, 2 KB 상한. `sendKind`는 confirmed 없이 서명하며 relay `grants.go:134-139`가 owner·쌍·kind·기한·시간당 한도·본문을 확인해 밖이면 거절한다. `publicFailure`가 `relay HTTP 403 approval_required`를 그대로 돌려준다. |
| ② 키·credential 노출 | 통과. `.cmd` 래퍼는 없다(Run 값은 `node.exe plugin.js orca start`). `defaults.json`은 모드·relay URL·폴더·봇·AUTO_RECEIVE만, Run 값과 `watch.log`는 경로·ID·kind·크기·대기 수만 담는다. |
| ③ 사용자 단위·경로 | 통과. `reg.exe add HKCU\...\Run`을 `shell:false`로 실행, 관리자 권한 없음. 경로는 argv로 전달되어 셸 주입이 없다. `--dir`는 사용자 지정 위치에 사용자 권한으로 쓸 뿐이다. |
| ④ receive untrusted | 통과. `orcaView`가 `untrusted: true`, 도구 설명에 지시 금지 명시, 수신 경로 로그에 본문 없음. |
| ⑤ 수신자 하나 | 대체로 통과. watch가 `orca-watch.lock`을 잡고 MCP는 `KNOWSLINK_AUTO_RECEIVE=off`로 자동 loop를 돌리지 않는다. 잠금 생존 판정 약점은 발견 1. |

## 발견 사항

1. medium `adapters/src/orca.ts:227-236`·`:395` — 잠금이 pid·uptime만으로 생존을 판정하고 종료 때 지워지지 않는다. 같은 부팅 안 로그오프·재로그온(또는 Fast Startup) 뒤 pid가 재사용되면 start가 watch를 띄우지 않고, stop·autostart off·install이 무관한 프로세스를 kill한다. named pipe 잠금 또는 kill 전 명령줄 확인을 권한다.
2. low `adapters/README.md:457-469` — `<설치>\node.exe`의 `\n`이 실제 줄바꿈이 되어 Run 값·toolsAdd 예시·`orca autostart off` 명령 세 곳이 깨졌다.
3. low `adapters/src/orca.ts:350,449` — 설치 경로에 공백이 있으면 `toolsAdd` 한 줄이 중첩 따옴표가 된다.

## 결론

accept. high 0. 보안 경계 ①–⑤에서 우회·노출은 찾지 못했다. 발견 1은 후속 수정으로 다룬다.
