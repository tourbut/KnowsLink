---
title: 이름 표시·친구 코드·친구 요청 구현 기록
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-NAMES-DEV]
summary: SAR-UX-ISSUES-002 2단위(U2/U4+해당 U7)의 기술 계획·결정·검증 연결
---

# 이름 표시·친구 코드·친구 요청 구현 기록

기준 main `8d0db4b`, 지시서 커밋 `6275431`, attempt `9a5aabbe90544e148909a807a7a14d18`에서 시작했다. 새 schema·공개 API·권한은 추가하지 않았다.

## 기술 계획

1. 이름 라벨 하나(`human-agent-label`: 이름 / 종류 / 구분값)를 조회·승인·legacy·연결 화면에 쓴다. 충돌 확장은 공용 `distinguish`로 모은다.
2. 상대 agent 이름은 `visibleAgent`로만 얻는다. 자기 agent와 현재 활성·비차단 친구만 이름을 보이고 그 밖은 ‘확인할 수 없는 에이전트’와 상세의 ID다.
3. 친구 코드는 친구·내 계정에 라벨 있는 읽기 전용 `textarea`로 보인다. 긴 코드가 줄바꿈되어 잘리지 않는다.
4. 친구 요청은 handler에서 앞뒤 공백만 제거한다. 빈 코드·동의 누락은 같은 화면 422와 입력 유지로 처리하고 매처의 정확 비교는 그대로 둔다.
5. 조회 선택기는 닫힌 select의 잘림을 없애려고 radio fieldset으로 바꾼다.
6. 친구 거부 화면은 친구 목록으로 복귀한다. 익명 429 조회 화면은 `Unavailable`을 쓴다.

## 결정

- D04 A-07은 읽기 전용 값을 `code`로 쓰라고 하지만 U2(더 새 결정)는 읽기 전용 입력을 요구한다. O-06의 입력 잘림을 피하려고 `input` 대신 `textarea`를 골랐다.
- 새 대화의 select는 이미 이름이 먼저라 바꾸지 않았다. 지시 범위는 조회 선택기다.
- legacy 화면 제목을 ‘기존 연결’로 바꿔 ‘내 에이전트’와 title이 겹치지 않게 했다.
- 하위 에이전트: 구현은 같은 템플릿 파일을 공유해 파일 소유를 나눌 수 없어 직접 했다. 읽기 전용 리뷰 하나를 배정했고 low 4건을 `516bf62`에서 고쳤다.

검증과 fixture는 [DEV 보고](../../evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md)를 따른다.
