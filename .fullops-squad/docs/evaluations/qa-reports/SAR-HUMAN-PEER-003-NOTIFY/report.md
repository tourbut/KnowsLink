---
title: Grok·Dots 질문 없는 수신 알림 검증
status: approved
updated: 2026-10-13
owner: coor
tasks: [SAR-HUMAN-PEER-003-NOTIFY]
summary: Grok 루틴 알림과 Dots 이벤트 구독을 연결하고 실제 질문 없는 화면 알림을 확인했다
---

# 결과

사용자는 Grok과 Dots 모두 질문해야 메시지를 확인한다고 보고했다. 기존 운영 제품26094e9와 두 계정·키·친구 관계를 보존했다. 제품 코드는 변경하지 않았다. 수신 알림만 설정했으며 자동 답장은 활성화하지 않았다.

Grok Bot은 기존 watcher가 메시지를 질문 전에 받았지만 루틴이 결과를 채팅에 알리지 않도록 설정돼 있었다고 보고했다. 새 수신만 노우 채팅에 알리도록 루틴을 수정했다. probe01a13094-80f4-7960-885c-1ad73617fbdf 송신 뒤 Grok에 질문을 보내지 않았고, 23:47:06 KST 실제 native 앱에서 일치 본문과 수신 알림을 확인했다.

Dots 실제 다닷에 knowslink.text.arrived 이벤트 구독을 연결했다. 자동화6accf332deac81918d3749a3f236f405가 enabled=true를 보고했다. 서버 읽기 전용 SQL로 실제 subscription 존재와 callback Verified를 확인했다. 첫 시험은 이벤트 실행이 늦어 failed:expired였다. 첫 실패를 통과로 바꾸지 않았다. 두 번째 시험은 질문 없이 delivered/persisted=true/completion=received·Inbox0을 확인했으나 native 화면에 알림이 바로 나타나지 않았다. 이후 다닷은23:58:45 알림 저장·23:58:57 ACK·turn_trigger=mcp_event를 보고했고 알림도 나중에 화면에 나타났다. 지연 원인은 확정하지 않았다.

Dots 자동화는 실제 다닷 대화로 알림을 전송하고 성공과 대화 기록을 확인한 뒤 ACK하도록 보완했다. 기존 구독을 유지하고 prompt만 수정했다는00:03:07 저장 보고를 받았다. 세 번째 시험01a130a4-f995-7ae0-a59c-26630af29956 뒤 Dots에 질문하거나 root receive를 호출하지 않았다. 실제 native 화면에 KL-DOTS-VISIBLE-20261013 본문과 일치 ID가 나타났다. screenshot을 직접 열어 사용자에게 보이는 알림을 확인했다. 이 범위에서 양쪽 질문 없는 수신 알림을 수락한다. 마지막 메시지의 ACK 결과를 별도로 조회하지 않았다.

## 근거와 한계

[evidence.json](evidence.json), grok-observation.txt, server-readonly.log, server-readonly-second.log, inbox-readonly.log, message-readonly.log에 근거를 보존했다. SQL은 BEGIN READ ONLY로 실행했고 callback URL·secret·키·쿠키는 조회하지 않았다. screenshot 원본은 evidence.json의 외부 경로에 보존했다. 명령은 기존 OPS SSH helper를 재사용했다. git diff --check는 종료0이다.

첫 만료와 두 번째 화면 표시 지연은 남은 위험이다. 장시간 수신, background 화면 갱신, 앱/컴퓨터 재시작과 구독 갱신은 재검증하지 않았다. 현재 구독 만료는2026-10-12T15:48:20Z이며 호스트 갱신 여부는 미검증이다. 자동 답장·외부 업무 실행은 별도 사용자 지시가 필요하다. 새 유료 서비스·도구 설치·SQL 쓰기·키 재발급·삭제·이동은 없다.
