---
title: SAR-ORCA-BOT-001 — Orca 봇과 Grok Bot이 KnowsLink로 고정 형식 메시지를 주고받는다
status: approved
updated: 2026-11-27
owner: coor
tasks: [SAR-ORCA-BOT-001]
summary: Orca 봇 adapter와 같은 주인의 자기 agent 사이 사전 승인(grant) 묶음의 목표·범위·한도·완료 조건·열린 질문
---

# SAR-ORCA-BOT-001 — Orca 봇과 Grok Bot이 KnowsLink로 고정 형식 메시지를 주고받는다

relay·웹·adapter·Grok skill에 걸치는 묶음이다. 사용자가 승인하면 `status`를 `approved`로 바꾸고 `## 승인`에 근거를 적는다.

## 배경

- 설계 정본: orca-helm `docs/orca-bots.md` 3절(PR #10, 2026-11-26판). 사전 승인 한도는 추천안으로 잠정 채택됐고, 구현 전에 주인이 최종 확인한다(같은 문서 3.5·9절).
- 지금 KnowsLink는 발송마다 사람 승인이 원칙이라 봇끼리 자동으로 주고받을 수 없다.
- 쿠어(Orca 봇)가 주인 지시를 2026-11-27 orchestration 메시지(msg_ffc613fd5eaa)로 전달했다.

## 바라는 결과

Orca 봇이 주인 Google 계정의 자기 agent `쿠어 (Orca)`로 KnowsLink에 참여한다. Grok Bot과 Orca 봇이 `orca.decision.v1`·`orca.task.v1`·`orca.answer.v1` 세 형식을 사람 발송 승인 없이 주고받는다. 그 밖의 형식과 한도를 넘는 발송은 지금처럼 발송별 승인을 받는다.

## 범위

1. adapter(`adapters/src/orca.ts`와 CLI 진입점 `plugin.js orca <cmd>`)
   - `setup-connect <URL> <폴더> --host orca`: 기본 폴더 `%LOCALAPPDATA%\KnowsLink\orca-agent`. 비공개 폴더·ACL 검사, 브라우저 Google 동의·지문 비교는 기존과 같다.
   - `watch <폴더> --host orca --bot <이름>`: 기존 polling/SSE 수신, 서명 검증, Inbox durable 저장, persist·ACK 뒤 도어벨 메일.
   - `receive <폴더>`: Inbox에서 하나를 읽어 `untrusted: true`로 출력한다(형식 이름, item, 본문).
   - `send <폴더> --to <agent> --kind <세 형식> --file <json>`: 스키마 검사 뒤 사전 승인 경로로 보낸다.
   - Orca 봇 채팅 안 MCP는 쓰지 않는다(orca-helm 6.4절).
2. 인증: 기존 Ed25519 agent 키·credential, 자기 agent(`self:` 출처). 키는 agent 폴더에만 두고 argv·로그·메일에 넣지 않는다.
3. 도어벨 메일: ACK 뒤 `orca orchestration send --to bot:<이름> --type status --subject "KnowsLink: 새 메시지 N건" --body "<고정 문장>"`. 본문에는 검증한 ID, 대기 수, receive 명령만 넣는다. item이 있는 `orca.answer.v1`은 `--item <n>`을 붙인다. 메시지 본문과 발신자 이름은 넣지 않는다. 같은 ID는 한 번만 깨운다(기존 wake marker). 메일이 실패하면 다음 polling에서 다시 보낸다. `orca` 실행 파일 찾기는 orca-helm `lib/orca.mjs`의 `resolveOrcaBin`·`runOrca`를 참고한다.
4. 사전 승인(relay·웹·adapter)
   - grant는 (보내는 agent, 받는 agent, 형식) 한 쌍이다. 두 agent는 같은 owner의 활성 자기 agent다. 친구·legacy 관계에는 적용하지 않는다.
   - 웹 "내 에이전트"에서 주인만 만들고 끈다. agent API와 모델의 `confirmed:true`로는 만들 수 없다.
   - 한도는 Q1의 답을 따른다. 넘거나 grant가 없으면 `approval_required`로 거부한다.
   - grant 철회, 키 철회, agent 철회로 끝나고 되살아나지 않는다(generation).
   - 본문은 여전히 untrusted다. receipt·로그에는 형식·ID·크기만 남긴다.
5. Grok 쪽 skill: `orca.decision.v1`을 주인에게 보여 준다(선택지는 번호). 답은 `orca.answer.v1`로 돌려준다. 주인 지시는 `orca.task.v1`로 Orca 봇에 넘긴다.

## 제약

- 토큰 비용: 새 LLM 리뷰, 통합 테스트, 상시 eval을 더하지 않는다. 테스트는 adapter·relay 단위의 표적 테스트만 돌린다.
- 운영 배포, 운영 계정 연결, 원격 push는 쿠어를 거쳐 주인 승인을 받은 뒤에만 한다.
- 되돌릴 수 있는 기술 결정은 dev가 추천안으로 정한다.

## 성공 기준

| ID | 기준 |
|---|---|
| S1 | grant 없음, 형식 밖, 본문 한도 초과, 시간당 한도 초과, 기한 지남, 다른 owner, 철회 뒤 발송이 모두 `approval_required`로 거부된다 |
| S2 | agent API와 `confirmed:true`로 grant를 만들 수 없다 |
| S3 | 도어벨 메일에 본문·발신자 이름이 없다. 같은 ID는 한 번만 깨운다. 메일 실패 뒤 다음 polling에서 다시 깨운다 |
| S4 | item이 있는 answer는 `--item`으로 메일을 보낸다 |
| S5 | `receive` 출력에 `untrusted: true`가 있다. receipt·로그에 본문이 없다 |
| S6 | 격리 relay에서 Grok 역할 agent와 Orca 역할 agent가 왕복 1회를 한다: task → 도어벨 → receive, decision → answer → `--item` 메일 |

운영 계정 왕복 1회는 주인 승인 뒤의 별도 단계다.

## 과제 구성

- `SAR-ORCA-BOT-001-DEV`(dev): 범위 1–5, 표적 테스트, S6 격리 왕복. 웹 화면은 기존 "내 에이전트" 패턴을 따른다.
- tester 독립 QA와 designer 시각 검수는 배정하지 않는다(토큰 비용 제약). 병합 전 리뷰는 `fullops merge`의 기본 판정을 따른다.

## 열린 질문

| ID | 질문 | 추천 |
|---|---|---|
| Q1 | 사전 승인 한도 | 본문 2 KB, 시간당 30건, 기한 30일. 기한 갱신은 주인이 다시 동의한다 |
| Q2 | 묶음 구성 | dev 과제 하나로 배정하고 tester·designer는 배정하지 않는다 |

## 승인

2026-11-27 사용자가 질문 카드로 승인했다. Q1은 추천안(본문 2 KB, 시간당 30건, 기한 30일, 갱신은 주인 재동의)으로 확정했다(D057). Q2는 dev 과제 하나로 배정한다(D058).
