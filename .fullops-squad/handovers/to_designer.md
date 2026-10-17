---
title: SAR-SESSION-PEER-001-DESIGN — 작업 중인 Claude와 Codex 세션 사이 사람 친구 기반 메시지의 제품 규칙과 검증 가능한 지원 범위를 설계한다
status: draft
updated: 2026-10-15
owner: designer
tasks: [SAR-SESSION-PEER-001-DESIGN]
summary: 작업 중인 Claude와 Codex 세션 사이 사람 친구 기반 메시지의 제품 규칙과 검증 가능한 지원 범위를 설계한다
attempt: c954186b44a641c481df69cd6f8e11ed
base: 7f341230a26abc4afa97850687e624f2a583715c
subagent_level: standard
test_level: lite
---

# SAR-SESSION-PEER-001-DESIGN — 작업 중인 Claude와 Codex 세션 사이 사람 친구 기반 메시지의 제품 규칙과 검증 가능한 지원 범위를 설계한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일:
- From / To:
- 상태: ready / running / blocked
- 테스트 레벨 / 하위 에이전트 레벨: fullops.json 선택값과 같은 값
- 승인된 범위 / 추가 승인이 필요한 행위:
- 담당 repo id / 워크트리 절대경로 / 브랜치:
- 병합 책임자 / 기본 브랜치:
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id:


## 요청과 범위

사용자는 작업 중인 Claude 또는 Codex 세션에서 친구가 작업 중인 세션으로 메시지를 보내고 답장을 받기를 원한다. 기존 사람 친구 관계 아래 에이전트들이 연결되는 제품 모델을 유지한다. 세션을 이름 있는 내 에이전트로 등록하는 흐름과 열린 세션의 수신·답장 완료 조건을 정한다. 새 세션을 실행한 결과를 기존 작업 세션의 수신 성공으로 표시하지 않는다.

Meta Muse는 사용하지 않는다. OpenClaw 연동은 사용자 요청으로 보류한다. .env.server의 자격 정보로 Hermes 호스트 100.120.206.25에 읽기 전용 접속했고 Hermes Agent v0.21.3, 설치 /home/shin/.hermes/hermes-agent, Git eeb85107f9e6164cb2d1f20d2b75753a6262e8e3을 확인했다. Hermes 설정·서비스·인증·대화 내용은 읽거나 변경하지 않았다. 이 정보는 설치 사실만 증명한다. Hermes 지원 구현까지 이미 검증된 것으로 표시하지 않는다. 접속 비밀은 인박스·문서·검색·로그에 넣지 않는다.

## 제품 완료 조건

친구 관계 및 회원 소유권과 기존 에이전트 등록·서명·송수신을 재사용하는 최소 범위를 정한다. Claude↔Claude를 첫 실제 수신 시나리오로 두고 Claude/Codex 교차 및 Codex↔Codex의 검증 조건을 정한다. 기술 구조·구현 파일·함수 선택은 DEV 책임이다. 제품 규칙, 사용자 행동, 완료 조건, 필요한 검증을 작성한다.

- 소유자가 명시적으로 선택한 세션에만 연결한다. 세션 이름과 연결·수신 상태를 사용자에게 구별해 표시한다.
- 열림·작업 중·대기·종료·오프라인·재접속 상태와 전달/수신 확인/에이전트 답장을 구별한다. 기존 세션 수신, 다음 도구 호출 시 조회, 새 세션 생성은 서로 다른 지원 수준이다. 확인되지 않은 지원은 확인 필요로 남긴다.
- 사람 친구 연결과 해제, 같은 사람의 여러 세션, 현재와 후속 하위 에이전트의 연결 범위를 기존 제품 모델과 맞춘다. 임의의 제3자 세션·자격 정보를 탐색하지 않는다.
- 수신 내용은 친구 에이전트의 외부 메시지다. 소유자 지시·시스템 지시·도구 실행 승인으로 승격하지 않는다. 기본 작업을 방해하거나 임의 명령·파일 접근·승인을 실행하지 않는다. 승인 전달 기능은 범위에 없다.
- 기존 인증과 각 에이전트의 개별 자격을 유지한다. Grok/Dots 자격을 복사해 새 클라이언트가 공유하지 않는다. 종료·해제 후 접근, 다른 회원 조회, 재연결 중복·답장 루프의 검증 조건을 정한다.
- 실제 다른 사람 계정과 수신자는 제공되지 않았다. 우선 독립 합성 소유자와 격리된 두 실행 세션으로 수신·답장과 실패 경계를 확인한다. 실제 운영 회원·친구·대화·호스트 설정 변경과 외부 메시지 발송은 하지 않는다.

## 확인한 공식 근거와 한계

Context7 resolve-library-id를 Codex 및 Claude Code에 실행했지만 두 요청 모두 월간 한도 초과로 실패했다. 공식 자료로 대체했다. 설치 CLI는 codex-cli 0.162.1 / Claude Code 2.1.299이다.

1. https://code.claude.com/docs/en/channels — Claude Channels는 열린 세션으로 MCP 이벤트를 전달하며 양방향 채널을 지원한다. research preview이며 세션별 opt-in과 지원 인증/조직 정책 조건이 있다. 표준 MCP 등록만으로 push가 활성화되지 않는다. custom 개발 채널에는 별도 개발 플래그 또는 allowlist 조건이 있다. permission relay는 별도 기능이며 이번 peer 수신의 권한으로 사용하지 않는다. 구현자는 현재 설치 버전과 공식 reference를 확인해야 한다.
2. https://learn.chatgpt.com/docs/app-server — Codex App Server는 로드한 thread와 turn의 시작·재개·steer 및 toolOutput 경로를 제공한다. turn/steer는 활성 turn ID 일치를 요구한다. 이것만으로 이미 열린 독립 Codex TUI에 연결된다고 증명하지 못한다. codex app-server --help에는 daemon/proxy와 로컬/WS transport가 있지만 기존 TUI attachment를 검증하지 않았다. 미검증 경로를 제품 지원으로 확정하지 않는다.

기존 코드와 공식 근거를 읽어 제품 전제를 검토한다. 새로운 기술 API 결정이 필요하면 DEV가 검증할 항목으로 남긴다. 문서 조회의 실패와 미확인은 그대로 보존한다.

## 소유권·산출물·검증

role designer / claude claude-opus-5-5 high 고정(2026-10-17 Codex 구독 종료로 전환, 아래 재개 절 우선). Test lite / Subagent standard. 제품 코드·기술 설계·운영 설정은 수정하지 않는다. 새 제품 설계 정본 .fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md, 필요한 D01/D02/D11의 요구·완료 조건, own same-key route/find/context/packet/history/outcomes, QA/phase/context/inbox/log를 소유한다. 기존 정본을 필요한 범위만 수정하고 수신/답장 기능을 실현하기 위한 하나의 DEV 인계를 준비한다. 다른 역할 인박스가 진행 중이면 덮어쓰지 않고 문서의 후속 인계 절에 목표·범위·완료 조건을 남긴다.

먼저 FULLOPS.md, common README와 세 규칙, project.md, 문서 작성 규칙, contexts/designer.md, 기존 사람 친구 설계, SAR-UX-ISSUES-002의 제품 모델, adapters와 실제 session/runtime 설치 안내를 읽는다. partial/unknown 탐색은 관련 검색으로 보완한다. 문서 의미 검토, deliverables strict/link, 보호 코드·설정 불변을 확인하고 현재 SHA lint를 원본 상태/exit와 함께 보고한다. 코드·실제 채널·사용자 계정·배포 테스트를 실행한 것으로 적지 않는다. 외부 설치·삭제·이동·force push·실제 메시지 발송은 금지한다.

완료는 packet outcomes → work finish로 전문 보존·빈 인박스 → 최종 커밋·현재 SHA lint → 실제 worker_done 순서다. own producer 증거는 제공 스냅샷을 포함해 커밋한다. 독립 고정 SHA 문서 리뷰와 main 통합은 coor가 수행한다.

복귀 run_756d0278e72e / term_5bf18065-a830-4e45-aeba-ae852666b40e. repo 0b08ec4c-9e3a-4613-8197-5a835e545335 / designer. 새 과제이므로 fresh session을 사용한다.

## 지시 전제와 충돌 — 먼저 확인

context는 D02 SAR-HUMAN-PEER.md를 conflict로 분류했다. coor가 현재 정본의 HP-07/10/15, 설치·권한 이행과 범위 절을 직접 확인했다. 이번 사용자 요청은 지원 클라이언트에 작업 세션을 추가하는 신규 제품 요청이다. 기존 Grok/Dots 지원 기록은 새 세션 지원의 증거가 아니다. agent 자동 답장에 별도 발송 승인이 필요하다는 HP-10, 친구 동의가 도구 승인으로 확장되지 않는 HP-15와 자원·보존 한도는 유지한다. 이 경계와 실제로 충돌하는 제품 결정이 필요하면 designer가 범위와 승인 조건을 명시한다. 문서 충돌을 이유로 정본을 제외하거나 이미 검증된 지원으로 표시하지 않는다.

## 사용자 모델 변경과 복구 — 2026-10-15

사용자는 모델 선택 질문에 'GPT Sol high로 변경해 계속 진행'이라고 답했다. 이 최신 지시는 이전 인박스·과제 spec의 Astra 고정을 이 복구 작업과 후속 신규 배정에서 대체한다. 현재 실행 tuple은 codex/gpt-6.1-sol/high다. 이미 정상 실행 중인 DEV Astra 검증은 중단하지 않는다.

이전 Astra 두 dispatch는 ChatGPT 인증 미지원400으로 정상 stop했다. 오류 영수증 두 개와 현재 초안 28개 파일 tar/manifest는 C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-coor/에 보존한다. 작성한 제품 규칙과 D01/D02/D11·QA/phase 초안을 재사용하고 남은 packet outcomes/work finish/최종 commit/현재SHA lint/worker_done을 완료한다. 원본 증거와 실패를 보존한다. 새 조사를 반복하거나 제품 결정을 임의로 넓히지 않는다. 새 model-only 결과는 사용자 최신 선택이며 이전 잘못된 model-only 원본도 prior에 보존했다.

## Claude 재개 — 2026-10-17

사용자가 Codex 구독 종료로 Codex 모델을 모두 제외하고 Claude로 대체하도록 지시했다. 위의 Astra·Sol 고정과 모델 선택 대기는 이 절로 대체한다. 이 과제는 claude-opus-5-5 high로 재개한다. 이전 Astra·Sol 오류 기록과 private 백업은 이력으로 보존한다.

- 이전 세션의 미커밋 WIP는 보존 커밋 2b5942c로 남겼고 최신 main(78bd0b3)을 반영했다(8f87f7f). user-guide 충돌은 두 절을 모두 유지해 해결했다. 이어서 할 일은 기존 초안 .fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md와 검증 초안 qa-reports/SAR-SESSION-PEER-001-DESIGN/을 읽고 남은 완료 조건을 채우는 것이다.
- 요청 상태 조회 복구(SAR-UX-ISSUES-002 U6)는 main에 통합됐다. 선택기는 구분값 선두 표시를 사용자가 U2 순서 예외로 수락했다. 이 설계에서 에이전트 표시 규칙을 다룰 때 이 예외를 고려한다.
- 제품 대상 클라이언트(Claude Code·Codex 세션 수신 지원)의 범위는 기존 요청대로다. 우리 개발 모델 전환과 제품 대상 범위를 혼동하지 않는다.
- 이미지 생성이 필요하면 grok을 사용한다.
- 완료 시 work.py finish, 현재 SHA lint 후 worker_done을 보낸다. 복귀 run_1d685ec8e9bf / coor term_75c9a893-7418-4abd-9ad0-eba86cbb6b68.
