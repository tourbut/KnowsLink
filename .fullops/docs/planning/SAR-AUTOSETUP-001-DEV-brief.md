---
title: SAR-AUTOSETUP-001-DEV — Claude Code·Codex 자동 설정 구현 지시서 초안
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-AUTOSETUP-001]
summary: "Claude Code·Codex 자동 설정 DEV 지시서 초안: 목표·확정 규칙·범위·완료 조건·증명 방법"
---

# SAR-AUTOSETUP-001-DEV — Claude Code·Codex 자동 설정 구현 지시서 초안

이 문서는 coor가 dev 인박스(`handovers/to_dev.md`)를 쓸 때 옮길 초안이다. intent 승인 전에는 배정하지 않는다. intent 열린 질문 Q1–Q6의 답이 추천안과 다르면 coor가 해당 줄을 고친 뒤 배정한다.

## 목표

Claude Code·Codex 사용자가 회원 화면의 명령 하나와 사람 확인 단계만으로 설치·Google 연결·세션 연결·첫 발송을 마친다. 같은 명령으로 기존 설치를 갱신한다. 상위 의도는 [intent SAR-AUTOSETUP-001](intents/SAR-AUTOSETUP-001.md)이고 DEV가 맡는 성공 기준은 S1–S7 전체다.

## 확정 규칙(제품)

정본은 [D02 SAR-AUTOSETUP-001](product-specs/SAR-AUTOSETUP-001.md)이다. 요약:

- AS-01–AS-13: 한 명령, 설치 계획 먼저, 자동으로 하는 일과 사람이 반드시 하는 일, 설치기의 Google 연결 시작, 자격 경계, 호스트별 agent, 세션 연결 분리, 버전·OS 확인, 기존 갱신, 실패와 다시 실행, 출력 언어·기록, 제거 안내.
- 흐름 N1–N8, U1–U4, 실패 F-01–F-09.
- 화면 카피: [D04 자동 설정 절](../design-docs/mockups/service-member-ui.md#claude-codecodex-자동-설정-안내--sar-autosetup-001).
- 바꾸지 않는 규칙: LG-04, HP-10, SP-01–12, D023–D030. 기존 세션 연결·발송 승인 대화상자·Codex 4자리 코드 동작은 그대로다.

## 범위

- 포함: Windows 11의 Claude Code·Codex CLI 설치기, 고정 버전 배포물과 checksum, 호스트 등록, 설치기의 Google 연결 시작·키 저장 확인, 기존 갱신, 실패 F-01–F-09, KO·EN 출력, 회원 화면 `use=claude` 안내 변경과 `use=codex` 선택지 추가, 6번 선택지 보조 문구 변경, adapters·루트 README·D11의 해당 절 갱신.
- 제외: macOS·Linux 지원 문구, Claude Desktop·Grok CLI·Cursor·Hermes, Grok Bot installer·Dots 흐름 변경, 제거 명령, 공개 marketplace, 운영 배포(ops 별도 승인), 실계정 사용(사용자 별도 승인).
- 기술 계획은 DEV가 정한다: 설치기 형태, 배포 위치, 호스트 CLI 등록 방법, 기존 설치 찾기, 원자적 교체 방법, 설치 기록 형식. 호스트 CLI·플러그인 API는 Context7 또는 공식 문서로 확인하고 지시서에 근거를 남긴다. 기존 `scripts/package_plugin.py`의 `install()` 보호 조건(절대 경로, 기존 폴더·private 폴더 거부, Claude `KNOWSLINK_AUTO_RECEIVE=off`, Codex `KNOWSLINK_SESSION_HOST=codex`)을 낮추지 않는다.

## 착수 때 할 일

1. 현재 수동 절차의 기준선을 깨끗한 Windows 11에서 측정한다. 사람 행동 수·직접 입력한 경로 수·명령 수·소요 시간을 기록한다. intent 문제 절의 표는 문서 기준 추정이다.
2. 선택한 기술 계획이 AS 규칙과 충돌하면 구현 전에 coor를 통해 designer에게 묻는다.

## 완료 조건

- D02 AC-01–AC-10을 DEV 자체 검사로 확인한다. AC-04 시간 측정과 AC-01 행동 표는 tester 큰 틀 검수에서 다시 확인한다.
- 실패 주입: F-01–F-09를 격리 fixture에서 KO·EN으로 확인한다. 실제 운영 relay·실계정을 쓰지 않는다.
- 자격 경계(AC-05·06): 설치 전후 해시 비교로 다른 agent 폴더·호스트 인증 파일·다른 MCP 항목의 변경 0을 보인다.
- 갱신(AC-07): 같은 명령 두 번 실행 뒤 agent ID·키 지문·친구 수·웹 agent 수가 같다.
- lint·product-test 통과, 별도 세션 고정 SHA 코드 리뷰.

## 증명 방법

- `qa-reports/SAR-AUTOSETUP-001-DEV/report.md`에 AC별 결과·명령·종료코드·고정 SHA를 적는다.
- 사람 행동 표(AC-01)는 단계·행동·입력값 종류로 쓴다. 키·토큰·실제 이메일은 적지 않는다.
- 화면 캡처는 만들지 않는다. 화면 시각 검수는 designer 후속이다.

## 후속

- tester: 새 사용자 흐름이므로 통합 후보에서 큰 틀 검수 1회(S1–S7, 호스트별 1회).
- designer: D04 절의 시각 검수(1440·390·320, KO·EN), AC-10 통과 뒤 U8 표 갱신.
- ops: 배포물 게시와 운영 반영은 사용자 승인 뒤 별도 과제.
