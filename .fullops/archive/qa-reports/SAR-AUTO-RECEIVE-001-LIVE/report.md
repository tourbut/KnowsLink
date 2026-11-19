---
title: 실제 노우 자동 수신 1차 통합 검증
status: draft
updated: 2026-10-11
owner: coor
tasks: [SAR-AUTO-RECEIVE-001]
summary: 화면 수신 유도 없이 실제 노우 질문·관련 답장·자동 inbox 저장을 검증한 1차 수락 기록
---

# 실제 노우 자동 수신 1차 통합 검증

## 결과와 범위

2026-10-11 00:18 KST에 실제 운영 도메인 https://link.knowslog.com에서 한 건의 관련 왕복을 확인했다. Bot 설치 SHA는 25adea2a209475b11c1203de4cf7f2d0bbad3544다. 제품 변경 SHA는 f8dd6659d7a22466e7c84a67431448fd7727e7c3다. 테스트 레벨 lite, 하위 위임 standard를 유지했다.

사용자 요청은 현재 세션에서 노우의 작업 상태를 질문하고 그록봇 화면을 통한 수신 유도 없이 답장을 받는 것이다. 이번 검사는 설치와 단일 비민감 답장 권한을 먼저 설정했다. 질문 송신 이후에는 그록봇 화면을 읽거나 조작하지 않았다. 실제 stdio MCP client를 연결해 자동 알림과 receipt를 관찰했다.

## 설치 보고와 직접 관찰을 구분한다

Bot의 일회 설치 보고: install_bot_mcp.sh exit0, 새 MCP process connected/도구9개/autoReceive running, 기존 public-node 환경·연결 폴더 유지, wake-check ok/configuredListed true, 실제 노우 UUID a87f86c1-4a66-4b38-9249-c09fbf942f9e, watcher PID170289. 키 없는 과거 미완성 설치 폴더는 knowslink.prev-20261009로 보존했다고 보고했다. 이 항목은 노우의 설치 보고이며 이 세션의 독립 shell 측정값으로 표시하지 않는다.

직접 관찰: 질문01a12664-c5f6-756e-b18b-c3e6c3fbb973을 agent_943334beca406f0c3417d2에서 agent_077c666294c4eb28b783f8로 한 번 송신했다. 15:18:21 UTC 접수 뒤 15:18:32에 transport delivered/completion received를 확인했다. 이 상태만으로 노우 턴을 판정하지 않았다.

15:18:40 UTC에 관련 답장01a12665-0c5c-7728-9daf-65b0d9f6f994의 MCP automatic_notice가 도착했다. 15:18:44에는 수동 receive 호출 전에 autoReceive pending1/lastReceivedId 일치를 확인했다. 이어 같은 inbox에서 답장을 표시했다. from, reply_to, untrusted:true를 assert로 확인했다. 원 질문 receipt의 completion reply_received와 reply_id도 일치했다. 전체 harness exit0이다. 답장은 현재 작업 상태의 요약이었으며 실행 지시로 사용하지 않았다. Git 증거에는 본문을 제외했다.

## 근거와 제한

같은 폴더의 evidence.json에 시간·ID·상태·자동 알림·수동 표시 전 inbox 상태를 보존했다. 원문이 포함된 임시 실행 로그는 Git 공용 디렉터리 fullops-gate 안의 로컬 파일이며 원격에 올리지 않는다. 실제 노우의 관련 답장과 자동 저장까지 확인했으므로 1차 자동 수신 왕복은 수락한다.

Gateway의 개별 .wake marker와 노우 내부 턴 trace는 직접 수집하지 않았다. 원문을 전달하지 않는 loopback gateway 경로는 미문서화 기능이다. vendor 업데이트 뒤에는 다시 확인해야 한다. 장시간·재부팅·Bot Update/Reset·동시 대화 검사는 이번 lite 범위 밖이다. 기존 독립 QA의 중복·재시작 자동 검사를 이번 실제 장시간 측정으로 바꾸어 표현하지 않는다.

현재 세션의 검증용 stdio MCP 연결은 확인 후 정상 종료했다. 이 클라이언트의 자동 수신은 MCP process가 연결된 동안 동작한다. Bot의 상시 watcher는 종료하지 않았다. Codex 대화 턴의 자동 wake는 이번 검사 범위가 아니다. 서버·Tunnel·DB·기존 키·관계를 변경하거나 새 유료 API·routine·결제 동의를 사용하지 않았다.
