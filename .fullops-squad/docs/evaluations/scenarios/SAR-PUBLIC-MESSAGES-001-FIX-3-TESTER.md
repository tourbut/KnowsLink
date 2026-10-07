---
title: 최종 메시지 후보의 독립 기능 QA 준비와 중단 범위
status: draft
updated: 2026-10-07
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER]
summary: 준비한 교차 owner·정리·snapshot 시험과 컴파일 중단 상태를 연결한다
---

# SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER — 미실행 시나리오

대상은 `d08903a55c3638128827010400e66e9d45b61d7c`다. [실제 결과](../qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md)는 QA 미완료다.
기존 [원09c 시나리오](SAR-PUBLIC-MESSAGES-001-TESTER.md)의 요구값과 실패 증거를 유지한다.
새 fixture는 이메일 inbox로 두 owner를 만들고 별도 agent를 정상 연결한다. Node/MCP의 요청·회신·receipt ID를 확인한다.
신규 로컬 16과 DB 16을 포화시킨 뒤 restart snapshot에서 자기 cancel·invite deny·unpair·key/agent revoke·logout/all의 303과 실제 상태를 확인한다.
같은 DB 시각에서 epoch 11 뒤 10은 11을 유지해야 한다. 같은 시각 epoch 12는 갱신돼야 한다. 더 늦은 시각의 복원 epoch 1도 갱신돼야 한다.
원 H1/M1·HTTP 행 잔류와 전체 PS08–11 경계는 원 프로브 및 기존 실제 DB integration으로 검증하도록 준비했다.
시험 원문은 `../qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER-test/probe_test.go.src`와 `normal_test.go.src`다. runner·Node helper도 같은 디렉터리에 둔다.
이번 runner는 unused import로 컴파일 exit 1이다. 기능 assertion은 미실행이다. 사용자 중단 지시에 따라 수정·재실행하지 않았다.
실메일·운영·vendor·실24h·시각 수락은 별도 미검증이다. 이 준비 문서는 PASS 증거가 아니다.
