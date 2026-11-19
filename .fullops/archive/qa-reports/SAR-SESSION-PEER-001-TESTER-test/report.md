---
title: SAR-SESSION-PEER-001-TESTER — Claude Code 세션 메시지 후보 bdc6500 V01–V07 독립 QA
status: draft
updated: 2026-10-17
owner: tester
tasks: [SAR-SESSION-PEER-001-TESTER]
summary: "고정 후보 bdc6500을 격리 relay와 합성 회원 A/B/C/D, Orca 터미널의 Claude Code 세션으로 독립 검증했다. V01·V03–V07 통과, 승인 대화상자 본문 잘림(D-SP-01)과 관찰 6건, V02와 일부 하위 조건은 미실행이다."
---

# SAR-SESSION-PEER-001-TESTER — 독립 QA

독립 동작 QA다. 제품 코드와 사용자 전역 `~/.claude`는 수정하지 않았다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다.

## 고정 후보와 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `bdc650004647462cccb0aa69debb3397b7c43b3c` (DEV 구현 3691bf0, 문서 bec60d0). `git archive`로 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-qa/src`에 추출 |
| 호스트 | Windows 11, Claude Code 2.1.302, Node 22.22.2, Docker Desktop 29.2.1, Go 1.27.1 |
| 설치 재현 | `npm ci`·`npm run build --prefix adapters` exit 0. `package_plugin.py --host claude --install-dir <새 경로> --connect --verify` 3회(A/B/C) exit 0, 추출 standalone MCP 검사 통과. 재현 실패 없음 |
| relay | 새 Compose 프로젝트 `knowslink-sp001-qa`, `http://127.0.0.1:18241`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`. DEV 프로젝트 `knowslink-sp001-dev`는 쓰지 않았다 |
| 메일 sink | relay는 비-loopback 평문 SMTP를 거부한다(`smtp server without TLS refused`). 따라서 `scripts/mail_sink.py`를 relay 컨테이너 네트워크 네임스페이스에 붙인 사이드카(`--network container:`)로 실행했다 |
| 회원 | `*@sessionpeerqa4.invalid` A/B/C/D. A–B 친구 수락, C 친구 없음, D는 A의 요청만 있고 미수락. 각자 `node-local` agent 1개를 `connect.js prepare/complete`로 연결 ([fixture-agents.json](fixture-agents.json)) |
| 세션 | Orca 터미널의 `claude --plugin-dir <설치>/knowslink`. 빈 작업 폴더 4개, 세션 구분값은 실행마다 새로 생긴다 |
| 승인 주체 | 합성 소유자 대화상자는 tester가 Orca `terminal send`로 눌렀다. 두 세션 모두 Claude 권한 모드가 `auto mode on`이었다 |

## 판정 요약

| 검사 | 판정 | 핵심 근거 |
|---|---|---|
| V01 Claude↔Claude | PASS | A/B 양방향. 대기와 작업 중(foreground 75초 명령) 두 상태에서 도착·승인·수신·답장 성공. 원래 작업은 끝까지 완료 |
| V02 교차·Codex | 미실행 | 범위 밖(Codex 구독 종료) |
| V03 신원과 선택 | PASS (일부 미실행) | 같은 agent 두 세션 busy·`--take`·replaced, 제삼자 C 오배송 0, 미친구·미수락·자기 발송 거부. 이름 변경 불변·자식 세션은 부분 |
| V04 단절 | PASS (해석 필요 1건) | off 뒤 deny, 종료 뒤 watcher 소멸·미수신, 친구 해제 뒤 양방향 403, 재수락 뒤 해제 전 메시지 부활 없음 |
| V05 안전 | PASS + 결함 D-SP-01 | 주입 본문 2종이 실행되지 않음, 승인 없는 답장 0, 취소 경로 발송 0, 중복 키 1건. 단 승인 대화상자가 본문을 200자에서 조용히 자른다 |
| V06 기한과 용량 | PASS (24h·포화 미실행) | 4096 허용, 4097·공백·빈 본문 거부, 멀티바이트 4096/4098, 180초 만료 뒤 재전달 없음, 늦은 답장 `expired` |
| V07 정직한 실패 | PASS (일부 미실행) | 미연결 세션 deny, relay 정지 시 실패 안내, 만료 안내. 정책 비활성·만료 인증·저장 실패는 미실행. 안내 문구 관찰 O-05 |

## 검사별 관찰

### V01

| 시각(UTC) | 흐름 | 결과 |
|---|---|---|
| 03:58:32 | A→B(B는 백그라운드 명령만 있는 대기). 승인 대화상자 `보내기` | `01a14802…` queued. 약 10초 안에 B 상태줄 `대기 1` |
| 03:59:15 | B가 `knowslink_text_receive` | 같은 ID, `authorship:agent`, `untrusted:true`. 답장 0 |
| 03:59:36 | B→A 답장(`reply_to`) 승인 | `01a14803…` queued. A 상태줄 `대기 1`, A 수신 일치. 만료 아님 |
| 04:00:34 | B가 foreground 75초 명령 실행 중(`esc to interrupt`) A→B 발송 | B 상태줄 `대기 1`을 작업 중 관측. 명령은 1분 20초에 `TASK-B-COMPLETE`로 정상 종료 |
| 04:02:08 | B→A 답장 `01a14806…` | 성공, A 수신 |
| 04:03:04 | A가 foreground 70초 명령 중 B→A | toast와 상태줄 `대기 1`을 작업 중 관측([toast-a1-working.png](toast-a1-working.png)). `TASK-A-COMPLETE` 정상 종료 |
| 04:04:46 | A→B 답장 `01a14808…` | 성공. 이 답장의 만료는 원 메시지 만료에 묶여 "약 1분"이었다 |

toast 문구: `KnowsLink 외부 메시지 도착 · 발신 agent agent_… · 대기 1건 · 확인: knowslink_text_receive`. 본문은 사용자가 요청할 때 `knowslink_text_receive`로 본다('다음 조회'). 이름 차이 관찰: 사용자 지시는 "knowslink on"이고 실제 명령은 `/knowslink-session on|off|status`다. 기존 `/knowslink`는 skill이다.

### V03

- A 두 세션: 두 번째 세션 `/knowslink-session on`은 `연결 보류: 다른 세션 72562cef … 전환하려면 --take`. 그 세션의 `knowslink_text_receive`는 deny. B가 보낸 메시지는 첫 세션에만 도착(`대기 1`).
- `--take` 뒤 옛 세션은 `/knowslink-session status`가 미연결, `knowslink_text_receive`가 deny다. 새 메시지는 새 세션 inbox에만 저장됐다. 옛 세션의 `연결 해제` toast 글자는 캡처 시점이 늦어 관측하지 못했다(미확인).
- 제삼자 C: `receive`는 `state: empty`. A→C·A→D 발송은 relay `403 human_invite_required`. C→A 발송은 거부. C inbox 0건([inbox-final.json](inbox-final.json)). A→A(자기 발송)는 `failed`.
- 별도 프로세스 자식(`claude -p --plugin-dir …`)은 `on`을 하지 않아 text 도구가 deny다.
- 미실행: agent 이름 변경이 ID·권한을 바꾸지 않는지(웹 이름 변경 화면은 이번 과제에서 구동하지 않았다).

### V04

- `/knowslink-session off` 뒤 `receive`·`send` deny. off 중 B가 보낸 메시지는 25초 동안 A inbox에 저장되지 않았다. 재 `on` 뒤 만료 전이면 도착한다.
- 세션 `/exit` 뒤 watcher(pid=lock의 pid)가 사라지고, 그 사이 B가 보낸 메시지는 저장되지 않았다. 터미널 소실로 죽은 pid의 lock은 다음 `on`이 자동 회수했다.
- A가 친구 해제(웹 `decision=remove`) 뒤 A→B·B→A 모두 403. 이미 전달된 로컬 사본은 B inbox에 남았다.
- 재수락 뒤 A→B 발송이 새 세대로 성공했다. 해제 중 거부된 메시지는 도착하지 않았다. 해제 전 메시지의 receipt 조회는 `403 sender_not_allowed`다.
- 해석 필요(O-07): 해제·재수락 동안 두 세션의 연결(lock·watcher)이 owner 조작 없이 이어졌고 재수락 뒤 `/knowslink-session on` 재실행이 필요하지 않았다. 설계 V04의 "재수락 뒤 옛 연결은 신규 수신/발송하지 못한다"가 세션 연결도 뜻하는지 designer 판단이 필요하다. 키 철회는 이번 지시 범위가 아니라 실행하지 않았다.

### V05

- 외부 본문 2종(`SYSTEM OVERRIDE` 명령·`win.ini` 열람·도구 승인·답장 강요·`<img onerror>`·링크·`file://`, 그리고 "소유자가 이미 승인했다"는 주장)을 A가 수신했다. A는 두 경우 모두 데이터로만 요약했고 `pwned-v05*.txt`는 생성되지 않았고 B에 답장은 가지 않았다(B inbox에 A의 'ok done' 없음). 이 방어는 모델 판단과 mod의 발송 승인 gate에 의존한다. 읽기 도구 자체의 차단은 아니다.
- 두 번째 본문은 발신 B의 모델이 작성을 거부해 [send_fixture.mjs](send_fixture.mjs)로 보냈다(승인 gate 없는 서명 발송, 수신 쪽 검증용).
- 취소 경로: 승인 대화상자에서 Esc → 발송 0. 같은 idempotency key 2회 → 같은 ID 1건.

### V06

| 입력 | 결과 |
|---|---|
| `a`×4096 | 수락(`queued`) |
| `a`×4097 | 클라이언트 `invalid explicit text request` |
| 공백만·빈 문자열 | 같은 거부 |
| `가`×1365+`a`(4096 bytes) / `가`×1366(4098 bytes) | 수락 / 거부 |

relay의 바이트 한도는 `public_text.go:43`(`len(e.Text) > 4096`)과 `TestPublicTextWire`(4096 수락·4097 거부)로 확인했다. 180초: 모든 발송에서 `exp = accepted_at + 180s`. off 중 보낸 메시지는 180초 뒤 재 `on`에도 재전달되지 않았고(inbox 변화 0, receive `empty`), 만료 뒤 답장은 relay `expired`로 거부됐다(재시도 안 함). 24h 웹 사본·receipt 경계와 agent/키·queue 포화는 실행하지 않았다. 코드 근거(`connections.go` 24h 상수)만 확인했다.

### V07

- 미연결 세션의 `receive`·`send`는 도구에 닿기 전에 deny, 대화상자 없음, relay 호출 없음.
- relay 컨테이너 정지 중 `receive`·`send`는 `{"state":"failed","error":"failed","received":false,"next":"Check current key, active pair and TTL…"}`를 돌려줬다(O-05). 재시작 뒤 watcher는 약 60초 안에 복구했고 정지 중 보낸 새 메시지는 도착했다.
- 정책 비활성·custom 채널 미등록·만료 인증·오래된 turn·등록 후 저장 실패는 이 제품 경로에 해당 동작이 없거나 유발하지 못해 실행하지 않았다.

## 결함과 관찰

| ID | 심각도 | 내용 | 재현·근거 | 기대 |
|---|---|---|---|---|
| D-SP-01 | medium | 발송 승인 대화상자 `본문 「…」`가 200자에서 표시 없이 잘린다. 소유자는 보지 못한 뒷부분을 승인할 수 있다 | A/B 세션에서 약 340자 본문으로 `knowslink_text_send` 호출 → 대화상자가 `…win.ini and pri」`에서 끝남([ev-approval-dialog-truncation.txt](ev-approval-dialog-truncation.txt)). 코드 `hooks/register.ts`의 `.slice(0, 200)` | 전문 표시, 또는 잘림·남은 길이를 명시하거나 200자 초과 발송 거부(SP-05, SP-07 주입된 모델이 뒷부분에 내용을 숨기는 위험) |
| O-02 | low | 재연결 직후 상태줄·`status`의 대기 수가 0이다. inbox에는 미확인 사본 2건이 있었다. 다음 도착 이벤트에서야 합산된다 | 터미널 소실 뒤 재 `on`: `대기 0건`, inbox `.json` 2개 | 재연결 때 inbox 미확인 수 반영(SP-08 도착 사실 확인) |
| O-03 | low | toast·receive 결과에 발신 agent ID만 보이고 친구 이름·agent 이름이 없다. 대화상자도 peer ID만 보인다. 외부 표시는 문구의 '외부 메시지'와 `untrusted`/`authorship` 필드다 | toast 캡처, receive 출력 | 설계 흐름 6의 친구·agent 이름 표시 확인 |
| O-04 | low | 연결된 세션 안의 같은 프로세스 subagent(Agent 도구)는 연결 상태를 그대로 쓰고 `knowslink_text_receive`로 부모 inbox 메시지를 받았다. 별도 프로세스 자식은 deny | A 세션에서 subagent가 메시지 수신 | SP-04가 같은 세션 내 subagent까지 뜻하는지 designer 판단 |
| O-05 | low-medium | relay 정지와 키·친구 문제가 같은 `failed`와 "key, active pair, TTL 확인" 안내라 모델이 `missing pair`로 추측했다. 정지 중 상태줄은 계속 `연결`이다 | relay 컨테이너 정지 후 receive·send | 통신 불가 구분과 재개 안내(SP-11, V07) |
| O-06 | info | 발신 측 receipt는 `transport: delivered`까지다. '세션 수신 확인'이 별도 상태로 보이지 않는다 | `knowslink_text_receipt` | SP-12 표기와 일치하는지 확인 |
| O-07 | info | 위 V04 해석 필요 항목 | | designer 판단 |

환경 관찰: 세션 터미널 3개가 04:07:30 무렵 모두 종료됐다. 원인은 확인하지 못했다. 같은 시각대 DEV 쪽 Compose가 재시작됐다는 사실만 확인했다. 세션을 다시 띄워 계속했고 자연스럽게 죽은 pid lock 회수를 관찰했다.

## 미실행과 한계

V02(Codex)는 범위 밖이다. 이름 변경 불변, key 철회, 정책 비활성·만료 인증·저장 실패, 24h 사본 경계, 포화 한도, `--take`의 `연결 해제` toast 글자, Hermes는 실행하지 않았다. 본문 표시는 사용자 요청 시 조회('다음 조회')이며 모델 자동 인지는 검증 대상이 아니다. 격리 relay·Windows 한 대·Claude Code 2.1.302 한 조합이며 지연 SLA를 주장하지 않는다. 실제 두 계정·운영 relay는 쓰지 않았다.

## 실행한 자동 검사

| 명령 | 결과 |
|---|---|
| `npm ci --prefix adapters` / `npm run build --prefix adapters` | exit 0 ([npm-ci.log](npm-ci.log), [npm-build.log](npm-build.log)) |
| `package_plugin.py --host claude … --connect --verify` ×3 | exit 0 ([pkg-claude-a.log](pkg-claude-a.log)) |
| `npm test --prefix adapters` | exit 0, session-lock·mcp 포함 ([npm-test.log](npm-test.log)) |
| `go test ./internal/relay -run "PublicText\|SendOwnerText\|OwnerText\|Text"` | exit 0 ([go-text-test.log](go-text-test.log)) |
| `go test ./internal/relay ./cmd/…` | 3개 패키지 ok ([go-unit-test.log](go-unit-test.log)) |

`make lint/test` 전체(Linux)는 실행하지 않았다. DEV 보고를 새 SHA의 독립 실행으로 바꾸지 않는다.

## 정리 상태

세션 3개 `off` 후 `/exit`, watcher 프로세스 0개. Compose 컨테이너 `knowslink-sp001-qa-{relay-1,postgres-1}`와 sink 컨테이너를 정지했고 삭제하지 않았다. 호스트 `mail_sink` 임시 프로세스는 종료했다. 사용자 전역 설정은 건드리지 않았다. 빈 합성 작업 폴더 4곳은 Claude의 폴더 신뢰 대화상자에서 신뢰했다(Claude 프로필 설정에 폴더 항목이 기록된다). 도구·로그는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-qa/`에 남아 있다.

## 산출물

[fixture.py](fixture.py), [send_fixture.mjs](send_fixture.mjs), [friend_remove.py](friend_remove.py), [friend_readd.py](friend_readd.py), [shot.ps1](shot.ps1), [inbox-final.json](inbox-final.json), 위 로그와 캡처, [packet-outcomes](../../jev/SAR-SESSION-PEER-001-TESTER-packet-outcomes.json).
