---
title: SAR-HOST-VERIFY-001-LIVE — 운영 Claude Code 작업 세션·Grok Bot 신규 설치와 실계정 왕복·알림 시나리오
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-HOST-VERIFY-001-LIVE]
summary: "운영 link.knowslog.com(제품 7b4f9bc)에서 두 호스트의 깨끗한 신규 설치·연결·친구·왕복·알림을 단계별로 검증하는 순서, 사용자 승인 지점, 증거 형식, 중단 조건"
---

# SAR-HOST-VERIFY-001-LIVE — 실계정 호스트 검증 시나리오

대상은 운영 `https://link.knowslog.com`의 제품 `7b4f9bcf11bb7ab13207d8e635060a6e90e7f8c2`다. 기록 기준은 main `e9bbdea`다. 결과는 [QA 보고서](../qa-reports/SAR-HOST-VERIFY-001-LIVE/report.md)와 같은 폴더의 `evidence.json`에 남긴다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 실계정 로그인·친구 요청·발송이 비가역이고 단계마다 사용자 승인이 필요하다.

## 공통 조건

- 허용 origin은 `link.knowslog.com`과 Google 로그인이다. 운영 서버·DB·ingress를 바꾸지 않는다. 서버 API를 직접 호출해 에이전트를 흉내 내지 않는다.
- 로그인은 사용자가 직접 한다. tester는 자격을 입력하지 않는다.
- 다음 행동은 실행 직전에 `orca orchestration ask`로 멈추고 사용자 승인을 받는다: 실계정 로그인 시작, 연결 동의, 친구 요청·수락, 메시지 발송, Grok 앱 설정 변경(MCP 등록·환경·루틴), 기존 설치·키의 철회·삭제.
- 기존 설치(Grok `/workspace/.knowslink`, 기존 agent·키·watcher·루틴, 사용자 `~/.claude` 설정, 다른 프로젝트)를 바꾸지 않는다. 새 폴더만 쓴다.
- 비밀 제외: 개인 키, credential, 쿠키, 친구 코드, 연결 token, webhook URL·key를 증거·채팅·argv에 넣지 않는다. 메시지 본문은 비민감 고정 문구만 쓴다.
- 명령은 자신의 종료코드를 그대로 남긴다(파이프로 가리지 않는다).

## 계정 배정 (첫 ask로 확인)

| 쪽 | 제안 계정 | 근거 |
|---|---|---|
| Claude Code 작업 세션 | 사용자 확인 필요 | 기존 Dots 계정 scshin88@gmail.com(신성철)을 쓰면 Grok 계정과 이미 친구다. 그러면 '새 친구 요청·수락'이 아니라 '기존 친구 관계에 새 agent 자동 포함'을 검증한다 |
| Grok Bot | knowslog01@gmail.com(Knowslog, 기존 Grok 계정) | 같은 Bot에 새 Command MCP·새 연결 폴더로 새 agent를 만든다 |

새 친구 요청·수락을 실제로 검증하려면 두 계정이 아직 친구가 아니어야 한다. 기존 친구 해제는 기존 관계 변경이라 이 시나리오의 기본 경로가 아니다.

## 단계 — Claude Code 작업 세션 (C)

| ID | 단계 | 승인 | 통과 조건 |
|---|---|---|---|
| C1 | 준비물 기록 | — | Claude Code·OS·Node 22.22.2·Python·Git 버전과 경로 기록 |
| C2 | 새 폴더 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/host-verify-001/`에 공개 저장소 clone, `7b4f9bc` detach, `npm ci`·build, `package_plugin.py --host claude --install-dir <새 절대경로> --connect --verify` | 사전 승인됨 | 각 명령 exit 0, `git status --short` 빈 출력, 설치 폴더·`-private` 폴더 신규 |
| C3 | `claude plugin validate`·`claude plugin test <설치>/knowslink` | — | exit 0 |
| C4 | 새 작업 폴더에서 `claude --plugin-dir <설치>/knowslink` 세션 시작(사용자 설정 불변), 도구 검색 | — | knowslink MCP 도구 목록과 `/knowslink-session` 명령 확인 |
| C5 | `knowslink_connect {"confirmed":true}` → 사용자가 URL을 열어 Google 로그인·지문 대조·이름·동의 | ask | `knowslink_connect_status` connected |
| C6 | 웹 '내 에이전트'에서 새 agent 이름·키 등록 확인 | — | 새 agent가 이름과 함께 등록됨 |
| C7 | `/knowslink-session on`·`status` | — | `KnowsLink 연결됨: 세션 <구분값> · 내 agent <agent>` |
| C8 | 친구(배정 결과에 따름): 새 요청·수락 또는 기존 관계에 새 agent 표시 | ask | 상대 화면에 새 agent가 보임 |
| C9 | Claude → Grok 새 agent 발송. 승인 대화상자의 `보내기`는 사용자가 결정 | ask | 승인 창에 상대 이름·본문 전체, 발송 성공·request ID |
| C10 | Grok 답장 도착 → toast·상태줄 `대기 N` → `knowslink_text_receive`로 본문·발신자·`sessionReceipt` | — | toast 문구와 본문·reply_to 일치 |
| C11 | 재개: `/knowslink-session off` → `on`, `status` | — | 재연결 문구와 대기 수 표시 |

## 단계 — Grok Bot (G)

| ID | 단계 | 승인 | 통과 조건 |
|---|---|---|---|
| G1 | Grok 앱 제품/버전·OS·Bot 컴퓨터 OS/arch, 계정 권한(Command 등록 메뉴) 기록 | — | 값 또는 '미기록' |
| G2 | Bot에 README '신규 설치' 실행 요청: 새 소스 폴더 `/workspace/knowslink-source-hv01`, **별도 prefix** `KNOWSLINK_PREFIX=/workspace/.knowslink-hv01`(기본 prefix는 기존 설치를 교체하므로 금지) | ask(앱 조작) | `Ready`, Node `v22.22.2`, `dist/plugin.js` 존재, 소스 SHA 기록 |
| G3 | 앱에 새 Command MCP `knowslink-hv01` 등록(기존 `knowslink` 항목 불변) | ask(설정 변경) | `knowslink_status`가 held로 실행됨(도구 검색) |
| G4 | 환경 `public-node`·RELAY_URL·새 폴더 `/workspace/.knowslink-hv01/connections/grok-hv-01` 설정·재시작 | ask(설정 변경) | 상태 변화 기록 |
| G5 | `knowslink_connect` → 사용자가 Google 로그인·지문 대조·동의 | ask | connect_status connected, 웹에 새 agent 이름 |
| G6 | 친구: Claude 쪽과 같은 배정 결과 | ask | 서로 새 agent 표시 |
| G7 | 수신·답장: C9 메시지를 새 연결로 받고 한 번 답장 | ask | 답장 ID·reply_to |
| G8 | 질문 없는 화면 알림: 새 연결의 webhook routine·watcher 설정(새 설치는 자동 생성 없음) | ask(설정 변경) | 질문 없이 Grok 채팅에 도착 알림 표시 |
| G9 | 재개: MCP 재시작 뒤 `knowslink_status` | — | 같은 agent 유지 |

## 증거 형식

`evidence.json`은 호스트별로 U8 후속 검증표 항목을 각각 기록한다: `product_version`, `os`, `account_permission`, `source_sha`, `package_check`, `clean_app_install`, `tool_discovery`, `auth`, `name_registration`, `friend_acceptance`, `round_trip`, `resume`, `notification`, `verified_at`, `evidence`. 각 값은 `PASS`·`FAIL`·`not_run`·`미기록`과 근거(명령·종료코드·화면 관측·ID)다. 소스/패키지 PASS로 앱 설치 PASS를 채우지 않는다. 화면 관측은 캡처 경로 또는 화면 텍스트로 남긴다. '신규 설치 검증 완료'는 한 호스트의 모든 단계가 같은 명시 환경에서 PASS일 때만 쓴다.

## 중단 조건

- 설치·빌드·검사 명령이 0이 아닌 코드로 끝나면 다음 단계로 가지 않는다. 원 오류를 보존하고 우회하지 않는다.
- 설치가 기존 경로(`/workspace/.knowslink`, 기존 연결 폴더, `~/.claude`)를 건드릴 징후가 있으면 멈추고 ask한다.
- 운영 5xx·로그인 루프·결제/요금 화면·Cloudflare 로그인 화면이 보이면 멈추고 escalation한다.
- 승인 답이 없거나 거절이면 해당 단계와 그 뒤 단계는 not_run이다.
- 상대가 의도하지 않은 대상(기존 agent 등)으로 발송될 징후가 보이면 발송하지 않는다.
