---
title: SAR-SESSION-PEER-001-RETEST — 최종 후보 df52535 수정 부분 재확인
status: draft
updated: 2026-10-17
owner: tester
tasks: [SAR-SESSION-PEER-001-RETEST]
summary: "df52535에서 승인 창 전체 본문·reply_to 표시, 세션 교체 뒤 deny, 같은 세션 subagent deny, V01 왕복은 통과했다. 새 관찰 1건: 호스트 승인 창이 약 1.9천 자를 넘는 본문을 '…'로 자른다."
---

# SAR-SESSION-PEER-001-RETEST — 재확인

독립 동작 QA(lite)다. 제품 코드와 사용자 전역 `~/.claude`는 수정하지 않았다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다.

## 고정 후보와 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `df525350b8d38a86c92a4f8aeb5c37449741ff9a` (`fullops/dev`). `git archive`로 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-retest/src`에 추출 |
| 호스트 | Windows 11, Claude Code 2.1.302, Node 22.22.2, Docker Desktop |
| 설치 | `npm ci`·`npm run build --prefix adapters` exit 0. `package_plugin.py --host claude --install-dir … --connect --verify` A·B exit 0 ([pkg-a.log](pkg-a.log)) |
| relay | 새 Compose 프로젝트 `knowslink-sp001-rt`, `http://127.0.0.1:18251`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`. 메일 sink는 relay 네트워크 네임스페이스 사이드카 |
| 회원 | `*@sessionpeerrt.invalid` A·B, 친구 수락, `node-local` agent 각 1개 ([fixture.py](fixture.py), [fixture-agents.json](fixture-agents.json)) |
| 세션 | Orca 터미널의 `claude --plugin-dir <설치>/knowslink`. 빈 작업 폴더 `work-a1`·`work-a2`·`work-b`. 모두 `auto mode on`. 합성 소유자 승인은 tester가 Orca `terminal send`로 눌렀다 |

## 판정 요약

| 검사 | 판정 | 핵심 근거 |
|---|---|---|
| D-SP-01 승인 창 전체 본문·길이 | PASS (255자·1387자), 관찰 O-RT-01 (3256자) | 아래 표 |
| 취소 시 deny | PASS | 255자 본문 `취소` → "소유자가 이 발송을 승인하지 않았습니다", 발송 0 |
| 답장 승인 창 reply_to | PASS | `답장(reply_to 01a14845-…)` 표시 |
| 세션 교체 뒤 원래 세션 deny | PASS | `--take` 뒤 옛 세션 receive·send 모두 deny, 새 세션만 수신 |
| O-04 같은 세션 subagent deny | PASS | subagent의 receive·send 모두 deny, 메인 대화는 수신 |
| V01 핵심 왕복 | PASS | 발송 승인→도착(`대기 1`)→수신→답장 승인→수신 |

## D-SP-01

| 본문 | 승인 창 | 판정 |
|---|---|---|
| 255자 (`0123456789`×25+`TAIL5`) | `본문 전체 255자 「…TAIL5」`: 끝 표식까지 전부 표시 ([ev-approval-255.txt](ev-approval-255.txt)) | PASS |
| 1387자 | `본문 전체 1387자`, 창 끝에 `」 을 그대로 보낼까요?` 까지 표시 ([ev-approval-1387-screen.txt](ev-approval-1387-screen.txt)) | PASS |
| 3256자(4096자를 요청했으나 모델이 옮겨 적으며 줄었다. 창이 보고한 길이가 실제 인자 길이다) | `본문 전체 3256자`이지만 약 1893자 뒤에서 `…`로 끝나고 질문 끝(`을 그대로 보낼까요?`)이 보이지 않는다 ([ev-approval-4096-screen.txt](ev-approval-4096-screen.txt)). 취소해 발송 0 | 관찰 O-RT-01 |

- 코드: `adapters/hooks/register.ts`는 본문 전체를 `ui.ask` 질문에 넣고(200자 `slice` 제거), `TEXT_MAX` 초과는 승인 창 전에 deny한다. 3256자 사례의 `…`는 훅이 아니라 호스트(Claude Code) 승인 창의 표시 한도로 보인다. 터미널 크기에 따른 한도인지는 확인하지 못했다.
- 4096자 정확 근처 표시와 4097자 초과 deny는 직접 구동하지 못했다. 4097 이상은 코드 경로(`args.text.length > TEXT_MAX`)와 DEV 단위 시험에 맡긴다. 이전 QA에서 클라이언트가 4097을 거부한 결과를 재사용한다.
- 잘림이 이제 '…'와 전체 길이로 드러나므로 이전의 조용한 잘림(D-SP-01)은 해소됐다. 다만 약 1.9천 자를 넘는 본문의 뒷부분은 소유자가 승인 창에서 볼 수 없다.

## 세션 교체

1. A 세션 1(`f710e889`)이 연결됨. 같은 agent로 세션 2(`ce51537b`)에서 `/knowslink-session on --take` → 연결됨.
2. 세션 1의 상태줄 표시가 즉시 사라졌다.
3. B가 A에게 보낸 `RT-take msg from B`: 세션 1에서 `knowslink_text_receive` → `KnowsLink: 이 세션은 연결되지 않았습니다. /knowslink-session on 으로 먼저 연결하세요.` deny. `knowslink_text_send`도 같은 deny이며 승인 창 없음.
4. 세션 2 상태줄 `대기 1`, 세션 2의 `receive`가 그 메시지를 받음(`authorship: agent`, `untrusted: true`).
5. 세션 1이 보내려던 `RT-old-session-send`는 B 쪽 inbox에 없다(B inbox는 V01 ping 한 건뿐, `claude-b-private/inbox`).

## O-04 subagent

- 세션 2의 Agent 도구 subagent가 `knowslink_text_receive`를 호출 → `<tool_use_error>KnowsLink: 하위 에이전트에서는 KnowsLink 메시지를 쓸 수 없습니다. 메인 대화에서 사용하세요.</tool_use_error>`.
- 같은 subagent가 `knowslink_text_send`를 호출 → 같은 deny. 승인 창 없음. B 쪽 inbox에 도착 0 ([ev-subagent-send-screen.txt](ev-subagent-send-screen.txt)).
- 전경 실행을 요청했으나 Claude Code가 백그라운드로 실행했다. 두 번 모두 백그라운드였고 전경만의 경로는 따로 확인하지 못했다.
- 메인 대화는 같은 세션에서 `knowslink_text_receive`로 메시지를 받았다(만료 후라 `expired: true`로 표시, 수신 자체는 성공).

## V01 왕복

| 흐름 | 결과 |
|---|---|
| A→B `RT-V01 ping from A` (승인 창: `새 메시지 · 본문 전체 18자`) | `보내기` 후 `queued`, id `01a14845-8eae-…`, TTL 180초 |
| B 상태줄 | 약 10초 안에 `대기 1` |
| B `knowslink_text_receive` | 같은 id, 본문 일치, `untrusted:true`, 답장 자동 발송 없음(`Reply not sent`) |
| B→A 답장 (승인 창: `답장(reply_to 01a14845-8eae-…) · 본문 전체 18자`) | `보내기` 후 도착, A 상태줄 `대기 1` |
| A `knowslink_text_receive` | `RT-V01 pong from B`, `reply_to` 일치 |

toast 글자는 이번에 캡처하지 않았다(상태줄 `대기 N`으로 도착만 확인). 이전 QA의 toast 확인을 재사용한다.

## 관찰

| ID | 심각도 | 내용 | 기대 |
|---|---|---|---|
| O-RT-01 | low–medium | 호스트 승인 창이 약 1.9천 자(이 터미널에서 약 11줄)를 넘으면 본문을 `…`로 자르고 질문 끝을 보여주지 않는다. 훅은 최대 4096자까지 승인 창에 넘긴다 | 승인 창으로 전체를 볼 수 없는 길이는 deny하거나 한도를 보이는 길이(예: 1.5천 자)로 낮추기, 또는 승인 창에 '일부만 보임' 안내 |

## 재사용한 이전 결과

V01 대기·작업 중 상태, V03 제삼자 오배송 0·미친구 거부, V04 off·종료·친구 해제, V05 주입 본문, V06 한도·만료, V07 실패 안내, O-02/03/05/06은 이번에 재실행하지 않았다. 근거: 지시서가 이 수정 범위(승인 창 본문 표시, 세션 교체 해제, subagent deny)를 한정했고, 이번 V01 왕복에서 발송·수신·답장 경로가 회귀 없이 동작했다. `git diff bdc6500..df52535` 전체를 독립 검토하지는 않았다.

## 실행한 자동 검사

| 명령 | 결과 |
|---|---|
| `npm ci --prefix adapters` / `npm run build --prefix adapters` | exit 0 ([npm-ci.log](npm-ci.log), [npm-build.log](npm-build.log)) |
| `npm test --prefix adapters` | exit 0, session-lock 포함 ([npm-test.log](npm-test.log)) |

Go 테스트와 `make lint/test` 전체는 실행하지 않았다. 이번 범위는 Claude 어댑터 훅이다.

## 정리 상태

세션 3개 `off` 후 `/exit`, 터미널 닫음. Compose 컨테이너 `knowslink-sp001-rt-{relay-1,postgres-1}`와 sink `knowslink-sp001-rt-sink`를 정지했고 삭제하지 않았다. 도구·로그는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-retest/`에 남아 있다. 사용자 전역 설정은 건드리지 않았다.
