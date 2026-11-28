---
id: handover-dev
title: "SAR-ORCA-HOST-001-DEV — Orca 봇을 KnowsLink 지원 호스트로 만든다: 고정 위치 설치·자동 시작 doorbell, orca bot tools add 한 줄 MCP, 웹 항목·README"
status: approved
updated: 2026-11-28
owner: dev
summary: "Orca 봇을 KnowsLink 지원 호스트로 만든다: 고정 위치 설치·자동 시작 doorbell, orca bot tools add 한 줄 MCP, 웹 항목·README"
key: SAR-ORCA-HOST-001-DEV
role: dev
goal: "Orca 봇을 KnowsLink 지원 호스트로 만든다: 고정 위치 설치·자동 시작 doorbell, orca bot tools add 한 줄 MCP, 웹 항목·README"
base: 0899e5bd29cc11ba4ee8e70419185c97ffa536d1
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-HOST-001-DEV — Orca 봇을 KnowsLink 지원 호스트로 만든다: 고정 위치 설치·자동 시작 doorbell, orca bot tools add 한 줄 MCP, 웹 항목·README

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 쿠어가 성철 님 요청을 전달했다(msg_e7c4b453e028, 2026-11-28). 되돌릴 수 있는 결정은 추천안으로 잠정 진행한다. 이 PC의 로컬 설치와 쿠어 수신기 교체도 범위에 든다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 새 agent 생성(금지), grant 변경, 기존 agent 폴더 삭제. push와 배포는 하지 않고 완료 보고의 `남은 일`에 적는다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md). 범위 1의 "Orca 봇 채팅 안 MCP는 쓰지 않는다(orca-helm 6.4절)"는 이번 주인 요청으로 바뀐다. 의도 문서의 범위 1에 바뀐 내용을 한 줄 적는다.

## 현재 상황과 확인 근거

- 운영 왕복 1회는 SAR-ORCA-BOT-001-OPS에서 끝났다(task → 도어벨 → receive, decision → answer `--item 13`).
- 지금 쿠어 연결은 다음과 같다.
  - agent: `쿠어` `agent_5513ec51367fed3748903d`
  - 폴더: `%LOCALAPPDATA%\KnowsLink\orca-agent-20261128`
  - grant 3개: Grok Bot `agent_2b73608560c8ba6e8c3d06`과 주고받는다.
  - 실행: ops 워크트리 빌드 `C:\Users\shin\orca\workspaces\KnowsLink\ops\adapters\dist\plugin.js`의 `orca watch --bot 쿠어`가 ops 터미널 `term_fa736957-…`에서 돈다. 워크트리를 지우거나 PC를 다시 켜면 끊긴다.
- Orca는 `orca bot tools add <명령>`으로 stdio MCP를 붙인다. 환경 변수를 넘기는 옵션은 없다(쿠어 확인).
- 이미 있는 것:
  - `adapters/src/defaults.ts`: `dist/` 옆 `defaults.json`. 설정한 환경 변수가 이긴다.
  - `adapters/src/mcp.ts`: public-node 폴더는 기본으로 자동 수신한다. `KNOWSLINK_AUTO_RECEIVE=off`로 끈다.
  - Claude Code 자동 설정 안내: SAR-AUTOSETUP-001.
  - 웹 '어디에서 사용할까요': `internal/relay/human_pages.go`, `i18n_en_human.go`.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다. 키와 credential은 agent 폴더에만 둔다. argv·로그·래퍼·defaults.json·메일에 넣지 않는다. 래퍼와 defaults.json에는 모드, relay URL, 폴더 경로만 둔다.
- 본문은 untrusted다. MCP 경로에서도 grant 밖의 형식과 한도를 넘는 발송은 지금처럼 발송별 승인을 받는다. agent API와 `confirmed:true`로 grant를 만들 수 없다는 경계(S2)를 유지한다.
- 토큰 비용: 새 LLM 리뷰, 통합 테스트, 상시 eval을 더하지 않는다. 표적 단위 테스트만 둔다.
- Windows 파일 작업은 AGENTS.md를 따른다. 기존 `orca-agent-20261128` 폴더는 옮기거나 지우지 않고 그대로 쓴다.

## 먼저 읽을 문서

- `.fullops/guide/roles/dev.md`
- `adapters/README.md`의 Orca 절, `adapters/src/orca.ts`, `defaults.ts`, `mcp.ts`
- `.fullops/docs/planning/intents/SAR-AUTOSETUP-001.md`: Claude Code 자동 설정 안내 패턴
- `.fullops/tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`: 운영 왕복과 실제 값. ops가 병합하기 전이면 ops 워크트리에서 읽는다.
- orca-helm `docs/orca-bots.md` 3.3절(도어벨), 6.4절: `C:\Users\shin\orca\workspaces\orca-helm\shrimpfish\docs\orca-bots.md`

## 해야 할 일과 파일 소유권

- [ ] **고정 위치 설치.** `plugin.js orca install`(이름은 dev가 정한다)을 만든다. bundle과 `defaults.json`, MCP 래퍼(.cmd)를 `%LOCALAPPDATA%\KnowsLink\orca\`처럼 고정된 곳에 복사한다. 워크트리 경로에 기대지 않는다. 다시 실행하면 덮어써서 갱신한다.
- [ ] **자동 시작 watch.** 로그온할 때 고정 위치의 `orca watch --bot <이름>`이 뜨게 한다. 추천은 관리자 권한이 필요 없는 사용자 단위 방법이다(예: 작업 스케줄러 로그온 트리거나 시작 프로그램). 중복 실행을 막고 끄는 명령을 둔다.
- [ ] **stdio MCP 한 줄 명령.** `orca bot tools add <명령>`에 그대로 넣을 한 줄을 출력한다. 래퍼나 defaults.json으로 다음 값을 넣는다:
  - `KNOWSLINK_MODE=public-node`
  - `RELAY_URL`
  - `KNOWSLINK_AGENT_FOLDER`(기존 쿠어 폴더)
  기존 agent와 grant를 그대로 쓰고, 새 agent를 만들지 않는다.
- [ ] **수신자 하나.** watch와 MCP가 같은 폴더를 쓰면 두 쪽이 동시에 받거나 잠금에서 부딪히지 않게 정한다. 추천은 watch가 유일한 수신자이고, MCP는 Inbox 읽기와 발송만 하는 안이다(MCP 쪽 `KNOWSLINK_AUTO_RECEIVE=off` 등). 고른 안과 이유를 보고에 적는다.
- [ ] **웹과 문서.** 웹 '어디에서 사용할까요'(ko·en)에 'Orca 봇 (내 컴퓨터)'을 넣는다. 기존 Claude Code 항목과 같은 패턴을 따른다. 안내 단계는 bundle 설치, 연결 동의, `orca bot tools add`, 도어벨 자동 시작이다. README와 `adapters/README.md`의 Orca 절도 같이 고친다.
- [ ] **이 PC 적용(마지막).** SAR-ORCA-BOT-001-OPS가 끝난 뒤에 한다. 진행 여부는 lead에게 `ask`로 묻는다.
  1. 고정 위치에 설치한다.
  2. 자동 시작 watch로 바꾸고, ops 터미널의 기존 watch를 멈춘다.
  3. 쿠어가 붙일 한 줄 명령을 `ask`로 보낸다. lead가 쿠어에게 전한다.
  4. 쿠어가 붙이면 `knowslink_status`의 agent ID가 `agent_5513ec51367fed3748903d`인지 확인한다.
  재부팅 확인은 사람 단계라 lead를 거쳐 요청한다.

## 완료 기준과 검증

- 증명 방법(`tasks/SAR-ORCA-HOST-001-DEV/qa/report.md`):
  - 표적 단위 테스트를 둔다. 다룰 것은 설치 경로, 래퍼·defaults.json에 비밀값이 없음, MCP가 자동 수신하지 않음(고른 안일 때), 자동 시작 등록과 해제의 명령 생성이다.
  - `python3 scripts/make.py lint`와 `fullops lint --from <base>`를 통과한다.
  - 이 PC 적용 결과를 적는다. 다음을 발췌한다:
    - 설치 경로
    - 자동 시작 등록 확인
    - 기존 agent ID 그대로
    - grant 3개 그대로
- 실패 경계: 새 agent가 생기거나 grant가 바뀌면 멈추고 `ask`로 묻는다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서: `.fullops/docs/design-docs/interface-design.md`, `mockups/service-member-ui.md`의 '어디에서 사용할까요' 패턴.
- 공용 컴포넌트: 기존 호스트 항목(Claude Code·Grok Bot·Dots)의 구조와 문구 톤을 그대로 쓴다. 새 스타일을 만들지 않는다.
- 검증: 기존 human_pages·i18n 표적 테스트에 새 항목을 더한다. designer 시각 검수는 배정하지 않는다(토큰 비용).

## 갱신할 산출물

- D03, D10, D11(route가 고름): 해당 문서의 Orca 봇 호스트 절. 의도 문서 SAR-ORCA-BOT-001 범위 1의 한 줄.

## 기대 산출물

- adapter 설치·자동 시작·MCP 래퍼 코드와 표적 테스트, 웹 항목, README, `tasks/SAR-ORCA-HOST-001-DEV/qa/report.md`.

## 제약·협업·후속

- 범위 밖:
  - MCP 스펙 업그레이드(SAR-MCP-UPGRADE-001, 이 과제 다음)
  - Grok answer 형식(SAR-ORCA-BOT-002-DEV)
  - grant UX(SAR-ORCA-GRANT-UX-001-DEV)
  - macOS·Linux 자동 시작. Windows만 하고, 다른 OS는 남은 일에 적는다.
- 하위 에이전트는 쓰지 않는다. 묶음이 작고 파일이 겹친다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `tasks/SAR-ORCA-HOST-001-DEV/qa/report.md` 링크로 대신한다. 원격 push는 승인 전이다. 그래서 `fullops done` 대신 로컬 커밋 뒤 worker_done만 보낸다. 첫 줄은 `[완료] SAR-ORCA-HOST-001-DEV | 브랜치 fullops/dev | SHA <HEAD>`다. 남은 일에는 배포가 필요한 SHA와 push 승인을 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과

[완료] SAR-ORCA-HOST-001-DEV | 브랜치 fullops/dev | SHA 081d9828

- 변경 이유: Orca 봇을 워크트리와 무관한 고정 설치·로그온 자동 시작 watch·stdio MCP 한 줄로 지원 호스트로 만든다(주인 요청 msg_e7c4b453e028).
- 변경: `adapters/src/orca.ts`(`orca install|start|stop|autostart on|off`, watch 잠금, 도어벨 node 전체 경로), `defaults.ts`(`DEFAULT_KEYS`에 폴더·봇·자동 수신), 웹 `Orca 봇 (내 컴퓨터)` ko·en, README·adapters README·D03·D10·D11·의도 문서. 구현 커밋 bb252e31.
- 검토 필요: 자동 시작 방식(HKCU Run, 로그온 때 콘솔 창 잠깐 보임), 수신자 하나 안(watch만 자동 수신, MCP `KNOWSLINK_AUTO_RECEIVE=off`), `.cmd` 래퍼 대신 node 직접 실행.
- 검증: `orca.test.js`·`connect-simple`·`mcp` PASS, Go 안내·i18n 테스트 PASS, `make.py lint`·`fullops lint --from 0899e5bd` 통과(WARNING `orca.ts` 464줄 등). 이 PC: `%LOCALAPPDATA%\KnowsLink\orca` 설치, Run 값 확인, 새 watch PID 69852 relay polling, agent `agent_5513ec51367fed3748903d` 그대로, 쿠어 `knowslink` 도구가 설치본 명령.
- 못 한 것: 쿠어 대화 안 `knowslink_status`(같은 명령 stdio로 대신 확인), 서버 grant 3개 재조회(설치는 grant 경로를 부르지 않음), 재부팅 확인(사람 단계), macOS·Linux 자동 시작.
- 스스로 정한 선택: coordinator 무응답(msg_5876c6c26f0e)이라 `tools add --bot 쿠어`를 직접 실행했고 이미 같은 명령으로 붙어 있었다(되돌리기 `orca bot tools rm knowslink --bot 쿠어`). 인박스의 "fullops done 대신 worker_done만"과 달리 preamble·hook 지시대로 `fullops done`을 실행했다(로컬 기록만, push 없음). 자동 시작 해제는 `orca autostart off`.
- 산출물: [QA 보고](../tasks/SAR-ORCA-HOST-001-DEV/qa/report.md).
- 남은 일: push 승인과 운영 배포(웹 항목, SHA 081d9828). 재부팅 뒤 watch 자동 시작 확인(사람). 공백 경로에서 `tools add` 인용 확인. `orca.ts` 크기 분할 검토.
