---
title: SAR-ORCA-HOST-002-DEV QA 보고
status: draft
updated: 2026-11-28
owner: dev
summary: SAR-ORCA-HOST-002-DEV 표적 테스트와 이 PC 설치본 갱신 결과
---

# SAR-ORCA-HOST-002-DEV QA 보고

- 테스트 레벨: lite (표적 단위 테스트 + 이 PC 설치본 갱신)
- 날짜: 2026-11-28

## 표적 테스트

`npm run build && node dist/orca.test.js` (Node 24.13.0)

```text
PASS: orca schema, signing bytes without kind unchanged, doorbell metadata only, ring once, retry after failure, --item answer, untrusted receive, fixed install defaults, autostart argv, one watch per folder
```

추가 단언(`adapters/src/orca.test.ts` 5절):

- `isWatch`: node 이미지 + `orca watch` 명령줄만 참. `orca start`, notepad, `null`은 거짓.
- 이번 부팅·살아 있는 다른 pid의 잠금이라도 probe가 watch가 아니면(`cmd.exe`, 읽기 실패) `liveWatch`는 null, `claimWatch`는 낡은 잠금을 바꿔 잡는다(kill 없음).
- probe가 watch면 `liveWatch`·`claimWatch`가 그 holder를 돌려준다.
- `releaseWatch`는 자기 pid 잠금만 지운다.
- `toolsAdd`: 공백 없는 경로는 큰따옴표 한 겹, 공백 있는 경로는 작은따옴표로 감싸 중첩이 없다.

실제 Windows probe: 자기 pid는 `{exe:"C:\Program Files\nodejs\node.exe", cmd:...}`이고 `isWatch` false, 없는 pid는 null.

`npm run lint`(eslint) 통과.

## 설치본 갱신 (이 PC)

1. 개발 node(v24)로 `orca install` 실행 → 기존 watch(pid 96412, 명령줄이 설치본 `orca watch`로 확인됨)를 멈춘 뒤 `node.exe` 복사에서 실패(`KnowsLink Orca install failed`). 봇 MCP(pid 94184)가 설치본 `node.exe`(v22.22.2)를 잡고 있고 크기가 달라 복사를 시도했다. 기존 동작이며 이번 변경과 무관하다.
2. 설치본 node로 다시 실행: `<설치>\node.exe adapters\dist\plugin.js orca install <폴더> --bot <봇>` → exit 0.
   - `agent`: `agent_5513ec51367fed3748903d` (그대로)
   - `watch`: `{"state":"started","pid":105428}`, 잠금 `{"pid":105428,...}`
   - Run 값 `KnowsLinkOrcaWatch` = `<설치>\node.exe <설치>\dist\plugin.js orca start`
3. `orca start` 재실행 → `{"state":"running","pid":105428}`. `orca watch` 명령줄의 node 프로세스 수 1.
4. kill 대상: 두 실행에서 `process.kill`은 1번의 검증된 watch 96412만 받았다. 봇 MCP 94184는 계속 살아 있다. 전후 pid 비교에서 사라진 나머지는 스냅샷용 단명 프로세스다(개별 식별은 못 함).

## 못 한 것

- 같은 부팅 pid 재사용의 실제 재현(로그오프·재로그온)은 하지 않았다. 단위 테스트의 probe 주입으로 대신했다.
- macOS·Linux `ps` 경로는 실행하지 않았다.
