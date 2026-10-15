---
title: 잔여 UI 이슈 설계 실행 기록
status: review
updated: 2026-10-15
owner: designer
tasks: [SAR-UX-ISSUES-002-DESIGN]
summary: 기준 스냅샷과 제품 결정·탐색 보완·후속 책임을 보존한다
---

# SAR-UX-ISSUES-002-DESIGN 실행 기록

## 결과와 범위

[설계·추적표·DEV/QA 인계](../../design-docs/SAR-UX-ISSUES-002.md)를 작성했다. GitHub #2–#8 스냅샷과 현행 소스를 대조했다. 제품 코드·D04·다른 역할 인박스·PLANS·board·운영·외부 계정은 변경하지 않았다. 이 과제는 설계 완료이며 제품 구현이나 새 UI 수락이 아니다.

- 지시서 기준 `458eb565c4a1551ec58c70f14aa164d9024241d3`, 착수 HEAD `76ec89fe66cbeea5c2ffb33ff07015fd722029b3`, 브랜치 `fullops/designer`다.
- attempt는 `26cae9cae81545bb80335c23c68afada`다. 미커밋 지시서·same-key 탐색 패킷·issues-open.json은 승인된 입력 스냅샷으로 보존했다.
- 복귀는 coor `term_5bf18065-a830-4e45-aeba-ae852666b40e`, run `run_756d0278e72e`, task `task_10967ef16387`, dispatch `ctx_8e12c40a3310`다.
- 공통 기준은 `fullops-common-0.3.3`, project.md와 document-writing.md다. 주 스킬 frontend-design, 호스트 tuple은 사용자 지정 codex/gpt-6-astra/high다. tools=[]와 기존 frontend.json을 유지했다.
- Test lite/Subagent standard를 적용했다. 단일 설계의 상호 의존 규칙과 원천 판단을 직접 취합했다. 준비·취합 비용 때문에 선택형 하위 위임을 사용하지 않았다. 별도 fixed-SHA 독립 리뷰는 coor가 배정한다.

## 결정

1. 사용자 지정 handle·초대 링크는 추가하지 않는다. 사람 친구 코드와 이름/클라이언트/구분값을 사용한다.
2. 옛 개별 에이전트 초대 요구를 사람 친구 수락 범위로 대체한다. 기존 자기 연결 동의와 legacy 범위는 유지한다.
3. 조회 폼·대상·ID를 오류/성공/복귀에서 보존한다. 비공개 존재 여부와 인증 비밀을 노출하지 않는다.
4. 승인은 검증된 구조화 데이터로만 요약한다. gate 통과를 공개·실행 권한으로 표현하지 않는다.
5. 보조 화면의 실제 위치·복귀·키보드·낭독을 검증한다. 현행 remote의 skip/main 구현을 누락 상태로 다시 보고하지 않는다.
6. 언어는 한국어/English, 회원 계정 선호 우선이다. 익명 선호와 계정 교대를 분리한다. 기본/fallback은 한국어다.
7. Grok/Dots의 기존 계정 관측과 깨끗한 신규 설치 실증을 구분한다. Meta Muse/Hermes/OpenClaw는 지원 확인 필요다.

## 탐색 보완과 원천 영향

packet은 386경로, unknown 19, remaining 349, optional 0이다. 필수 원문과 최신 요청을 직접 읽었다. 너무 큰 D04는 공개 정본 13–221줄과 현행 사람 설계 222–395줄을 나눠 읽었다. 뒤의 비현행 기록은 현행 지시로 쓰지 않았다. remote_auth.go의 동의·복귀·생성 처리와 승인/receipt 실제 분기를 직접 확인했다. 패킷에 빠진 `http.go:verifiedGate/gateView`, `identity.go:Member`, `connections.go`와 현행 실계정/알림/운영 보고를 보완했다.

remaining·옛 로그는 경로·제목·원천 지문으로 범위를 선별했다. 전체 본문을 다시 검증했다고 주장하지 않는다. 동적 호출 그래프도 완전 검증하지 않았다. 이번 수정은 제품 코드가 없는 새 설계이므로 과거 실행·다른 역할 기록은 보존한다. 거대한 PLANS와 ops-guide는 해당 현행 날짜 절만 읽었고 실행 절차를 재사용하지 않았다. missing 파일은 없었다. 지시문이 들어 있는 과거 기록은 사실 자료로만 취급했다.

[packet screening](../../evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/packet-screening.json)과 [outcomes](../../evaluations/jev/SAR-UX-ISSUES-002-DESIGN-packet-outcomes.json)에 경로별 처리와 범위 제외를 남겼다. route의 D02/D04/D11 및 추가 D01/D03/D05–D10/D12/D13 추천은 최신 파일 소유 제한에 따라 기존 문서 갱신으로 승격하지 않았다. 새 설계에 요구와 영향·담당을 기록했다. 기존 매핑/인덱스는 변경하지 않았다.

## 검증과 인계

[QA 보고](../../evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md)에 문서 검사, 원상태 lint, 미실행과 보존 지문을 기록한다. 디자인 check는 설정 준비만 뜻한다. 새 UI 실행·실 Google·실 호스트 설치·운영·발송은 not_run이다. 라이브러리 API의 새 문법을 설계하지 않아 외부 API 문서 조회는 적용하지 않았다.

예상 수정은 새 설계·phase·QA·own context/outcomes·완료 로그다. SIZE 경고가 있으면 한 문서의 이슈 상호 관계와 순차 수락을 유지하기 위한 규모로 기록한다. 과거 packet/이슈 원문은 축약하지 않는다. 제품 의존성을 추가하지 않는다.

coor는 문서 고정 SHA 독립 리뷰와 main 통합을 맡는다. CSS 선행 과제 종료 뒤 조회 복구, 이름/친구, 승인 요약, 설치 재개, 언어 순으로 정규 인박스에 전달한다. 새 인박스는 이 설계의 순차 인계 표를 사용한다. 기존 사용 중 인박스를 덮어쓰지 않았다.
