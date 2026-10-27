---
title: SAR-CONNECT-SIMPLE-001-LIVE 운영 실연결 검증 결과 (CS-AC-18)
status: draft
updated: 2026-10-27
owner: tester
tasks: [SAR-CONNECT-SIMPLE-001-LIVE]
summary: 운영 relay에서 회원 A(Grok Bot)와 회원 B(Claude Code)의 새 agent 연결, 친구 해제 후 재요청·수락, 서로의 agent 이름 표시, 시험 메시지 1건 수신을 확인한 CS-AC-18 결과와 옛 시나리오 기대 4건 갱신 기록
---

# SAR-CONNECT-SIMPLE-001-LIVE 결과

## 판정

**CS-AC-18 PASS.** 운영 relay `https://link.knowslog.com`(bf67dfcb)에서 두 회원의 새 agent가 연결되고, 친구 요청·수락 뒤 서로의 agent가 이름과 함께 보였으며, 승인된 시험 메시지 1건이 도착했다. 미해결 high 이상 결함은 0이다. 사람이 꼭 해야 한 행동은 회원 A 동의 링크 승인 한 번이다(아래 '행동 기록').

기준: main `3d67c5d8`(제품 `bf67dfcb`), 브랜치 `fullops/tester`. 테스트 레벨 lite, 하위 에이전트 없음. 계정 주소는 문서에 적지 않는다(회원 A·B).

## CS-AC-18 항목별 결과

| 항목 | 판정 | 근거 |
|---|---|---|
| 회원 B 연결 완료(키 등록됨) | PASS | Claude Code 클라이언트(공개 배포 0.1.0 ZIP, SHA-256 `767c0e87…`이 OPS 보고와 일치)의 `setup-connect --host claude` 출력이 `waiting`→`connected`(agent 이름 `Claude Code`). '내 에이전트' 카드는 '키 등록됨'. |
| 회원 B 지문 대조 | PASS | 클라이언트 출력 `SHA256:55bFmJcb…`과 동의 화면 지문이 일치(전체 문자열 비교). 앞 8자만 기록. |
| 회원 B 기본 이름(A-04) | PASS | 동의 화면 이름 칸의 기본값이 `Claude Code`. 바꾸지 않고 승인. |
| 회원 B W-01 링크 재방문 | PASS | '이 연결은 완료됐습니다 / Claude Code 에이전트에 키가 등록됐습니다. 이 링크는 다시 쓸 수 없습니다. 메시지를 받을 세션은 앱에서 연결하세요.' |
| 회원 B CS-AC-05 실패 시도가 목록에 안 늘음 | PASS | 두 번째 시도를 시작해 Google 로그인·지문 단계 후 '취소'. 클라이언트는 `failed`/`cancelled`와 문구 '연결을 취소했습니다. 다시 연결하면 새 링크를 만듭니다.'(사유 코드·문구·다음 행동 있음). 취소 전후 agent 카드 수가 4로 같다. |
| 회원 A 연결 완료(키 등록됨) | PASS | Grok Bot `링크` 봇의 새 Command MCP `knowslink-live`(새 설치, main `3d67c5d8`, defaults.json, 환경 변수 없음). `knowslink_connect` 뒤 `knowslink_connect_status`가 `connected`, 이름 `Grok Bot`, 지문 앞 8자 `-5iUJjtA`. '내 에이전트' 카드 '키 등록됨'. |
| 회원 A 지문 대조 | PASS(부분) | 봇이 보여 준 지문(앞 8자 `-5iUJjtA`)과 연결 완료 뒤 `knowslink_connect_status`의 지문이 같다. 동의 화면에서의 대조는 사용자가 승인했으므로 에이전트가 직접 하지 못했다(행동 기록). |
| 회원 A 기본 이름 | PASS | 새 agent 이름 `Grok Bot`(host=grok 기본 이름). 이름 입력은 사람이 한 행동 쪽에서 생략됐다. |
| 회원 A W-01 링크 재방문 | PASS | 같은 완료 화면 '이 연결은 완료됐습니다 / Grok Bot 에이전트에 키가 등록됐습니다…'. |
| 친구 요청·수락 | PASS | 아래 '친구 흐름'. |
| 서로의 agent가 이름과 함께 보임 | PASS | 해제 전·재수락 뒤 모두 A는 B의 `Claude Code`를, B는 A의 `Grok Bot`을 친구 화면에서 본다. |
| 도구 결과의 표시 이름(A-06) | PASS | Claude 쪽 `contact`: `{"name":"Grok Bot","friend":"Knowslog"}`. 수신 결과 `sender:{friend:"Knowslog",name:"Grok Bot"}`, `me:{name:"Claude Code"}`, `peer:{name:"Grok Bot"}`. Grok 쪽 `knowslink_status`: `me` 표시 이름 `Grok Bot`과 `sharedAgent` 안내. 한계는 관찰 2. |
| 시험 메시지 1건 | PASS | 아래 '시험 메시지'. |

## 행동 기록 (순서대로)

에이전트가 대신한 행동(A=회원 A, B=회원 B):

1. B: 공개 배포 ZIP 내려받기·SHA-256 확인·압축 해제, 새 폴더로 `setup-connect` 실행, 동의 링크 열기, 지문 대조, Google 계정 선택(지정한 회원 B 항목이 선택됨을 확인), 기본 이름 확인, 승인, 완료 확인, W-01 재방문.
2. B: 두 번째 시도 시작·로그인·취소, 실패 문구·카드 수 확인.
3. A: Grok Bot 앱의 `링크` 봇에 지시로 새 설치(`/workspace/knowslink-source-live`에 main `3d67c5d8` checkout, `install_bot_mcp.sh` `Ready`, node v22.22.2), 새 MCP `knowslink-live` 등록(승인 카드 없이 바로 등록·도구 9개), `knowslink_status`, `knowslink_connect(confirmed:true)`로 동의 링크·지문 받기.
4. A: 로그아웃 후 Google 계정 선택(지정한 회원 A 항목이 선택됨을 확인)·로그인, W-01 완료 화면 확인, '내 에이전트' 카드 확인.
5. 친구 해제(B 화면의 '친구 해제', 사용자 승인 뒤), A 친구 코드로 B가 요청, A가 수락. 코드 값은 기록하지 않는다.
6. 시험 메시지 1건(사용자 승인 뒤), 수신 확인.

사람이 꼭 한 행동:

- A 동의 링크 승인 1회. 동의 링크를 에이전트가 열어 보기 전에 사용자가 직접 승인했다(coor 전달 답, 2026-10-27). 그래서 A의 동의 화면 지문 대조와 이름 입력은 에이전트가 하지 못했다. 결함이 아니다. 사용자가 비밀번호·2단계 인증·새 기기 확인을 입력한 일은 없다(Google이 요구하지 않았다).
- 그 밖에 사용자가 꼭 해야 한 행동은 없었다. 승인 요청 3건(Grok 새 설치, 친구 해제, 시험 메시지)은 coor 경유로 사용자가 답했다.

## 친구 흐름

- 시작 상태: A와 B는 과거 시험의 친구였다(양쪽 친구 화면에 서로 한 명). 해제 전에 A는 B의 `내 에이전트 1`·`OpenAI Dots`·`그록봇-링크`·`Claude Code`를, B는 A의 `Grok 노우`·`Grok Bot`·`링크` 등을 이름과 함께 봤다(새로 만든 agent 두 개가 연결 직후 자동으로 보임).
- 사용자 승인(2026-10-27) 뒤 두 회원 사이 관계 하나만 해제했다('친구 해제'). 해제 뒤 양쪽 화면은 '아직 친구가 없습니다', 기록 '해제 / 다시 연결하려면 새 친구 요청과 수락이 필요합니다.'.
- B가 A의 친구 코드로 요청(동의 체크 후 '친구 요청 보내기', 보낸 행 '요청 대기 / 입력한 코드 …7YVA'), A 화면 '신성철의 친구 요청 · 서로의 현재·앞으로 추가할 에이전트가 연결됩니다', A가 '수락'.
- 수락 뒤 A가 B의 친구 화면에서 `Claude Code`를 포함한 B의 agent 목록을 이름과 함께 봤다. 다른 친구·agent는 건드리지 않았다.

## 시험 메시지

- 사용자 승인(coor 경유): 발신 A의 `Grok Bot` agent → 수신 B의 `Claude Code` agent, 본문 'KnowsLink 연결 시험입니다. 답장은 필요 없습니다.' 1건. 재전송·답장 없음.
- 발신: Grok `knowslink-live`의 `knowslink_text_send`, 승인 카드는 없었다. 접수 시 상태 `queued`, 요청 ID 앞 8자 `01a17db2`, 접수 23:10:39 KST, 기한 23:13:38 KST.
- 수신: B의 격리 agent 폴더에서 MCP `knowslink_text_receive` 호출 결과 `state:received`, 본문 일치, `untrusted:true`, `authorship:agent`, `sessionReceipt:recorded`, 만료 전. 한 번만 받았고 답장하지 않았다.

## 결함과 관찰

high 이상 결함은 없다.

| # | 심각도 | 내용 | 재현 |
|---|---|---|---|
| 1 | 관찰 | 새 설치 스크립트가 출력한 `npm audit` 경고: 취약점 4건(moderate 2, high 2). 봇의 보고이며 직접 확인하지 않았다. 의존성 범위가 런타임인지 빌드 도구인지 모른다. | `install_bot_mcp.sh` 실행 출력. dev 확인 필요. |
| 2 | 관찰 | Grok `knowslink_status` 결과에 연락처·친구 목록이 없다. 상대 agent 이름 표시(A-06)는 Claude 쪽 `contact`·수신 결과로 확인했다. Grok 쪽에서 상대 이름을 보는 도구 경로는 검증하지 못했다. | 봇이 `knowslink_status` 결과를 보고. |
| 3 | 관찰 | 기존 `knowslink-link`(옛 번들 0.1.0, 폴더에 10/26 agent.json)는 `knowslink_connect`가 사유 코드 없이 `failed`로 끝났다. 새 번들 이전 설치라서 당연하며 D02 E 표 위반이 아니다. 새 설치(`knowslink-live`)는 사유 없이 실패한 적이 없다. | 봇이 옛 연결로 `knowslink_connect`를 실행. |
| 4 | 관찰 | 다른 회원 세션으로 연결 링크를 열면 '연결 수단이 만료·취소·사용되었거나 유효하지 않습니다. 아직 새 키가 연결되지 않았습니다.'와 로그인 화면 링크만 보인다(이름·지문 없음, 정보 노출 없음). 로그아웃 상태도 같은 문구다. 지정 계정으로 로그인한 뒤 링크를 다시 열어야 한다. D04가 현행 문구 유지로 정한 동작이다. | B 세션에서 A의 링크 열기, 로그아웃 후 열기. |
| 5 | 관찰 | A의 '내 에이전트'에는 이전 시험으로 남은 이름 없는 agent 3개(키 등록 2·미사용 1)가 있다. 이 과제가 만든 것이 아니며 건드리지 않았다. | `/home/agents`. |

## 변경한 것과 되돌리는 방법 (비밀 없음)

- Grok Bot `링크` 봇 컴퓨터에 새 prefix `/workspace/.knowslink-live`(node, bundle, `defaults.json`, 연결 폴더 `connections/grok-agent`), 소스 checkout `/workspace/knowslink-source-live`(main `3d67c5d8`), 새 Command MCP `knowslink-live`(Command `/workspace/.knowslink-live/node/bin/node`, Arguments `/workspace/.knowslink-live/knowslink/dist/plugin.js`, 환경 변수·헤더 없음)를 더했다. 기존 `/workspace/.knowslink`, `knowslink`, `knowslink-link` 연결, 10/26 agent 폴더는 바꾸거나 지우지 않았다.
- 제거: Grok Bot 앱 설정에서 `knowslink-live`를 삭제하고 위 두 폴더를 소유자가 정리한다. 새 agent(`Grok Bot`, `Claude Code`)는 사용자 지시 없이 철회하지 않는다. 필요하면 '내 에이전트'에서 철회한다.
- 운영 relay에 새 agent 2개, 취소된 연결 시도 1건, 친구 관계 해제·재수락 1회, 시험 메시지 1건이 남는다. 기존 운영 agent `agent_c18976…`과 이름 없는 agent, 기존 Grok `knowslink` 연결은 건드리지 않았다.
- Claude 쪽은 테스트 폴더의 agent 폴더만 썼다. 실제 `~/.claude`·`~/.codex` 설정은 바꾸지 않았다. 원시 로그·agent 폴더(개인키 포함)는 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-LIVE/`에 있고 Git 제외 폴더다. 캡처는 만들지 않았다.

## 옛 시나리오 기대 4건 갱신

[SAR-UX-ISSUES-002-CONNECT-TESTER](../../scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md)에 '현행 규칙 기대 갱신' 절을 더하고 [SAR-CONNECT-SIMPLE-001-TESTER](../../scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md)의 회귀 행·CS-AC-18 인계 줄을 맞췄다. 갱신한 기대: 선택지 6개, Claude Code 자동 설정 안내, agent 0개 화면의 `use` 6개, 여섯째 agent 생성 허용. 제품 코드와 옛 `checks.py`·`mutate.py`는 바꾸지 않았다.

## 한계

- 자동 설정 설치기(`knowslink-setup.ps1`)와 Claude Code 세션 안의 `/knowslink-session on`은 이 과제 범위가 아니다(D040 별도). 클라이언트는 `setup-connect` 진입점과 MCP stdio 호출로 확인했고 실제 Claude Code 대화 호출은 하지 않았다.
- 회원 A 동의 화면의 에이전트 직접 지문 대조는 하지 못했다(사람이 승인).
- 실제 Google 비밀번호·2단계 인증 화면은 나타나지 않았다.
- 시험 메시지는 A→B 한 방향 1건이다. 답장·대화는 시험하지 않았다.

## 남은 일

- D02·D04·D12 변경이 필요한 점은 없다. 관찰 1(npm audit)은 dev가 의존성 범위를 확인하면 된다.
