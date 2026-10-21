---
title: SAR-PUBLISH-001-TESTER — 캘린더 제거·일회 정리·정책 페이지 독립 QA 시나리오
status: draft
updated: 2026-10-22
owner: tester
tasks: [SAR-PUBLISH-001-TESTER]
summary: "main 17a360b2의 서버 캘린더 제거(D018)·시작 시 일회 정리·/privacy·/terms(D017)를 격리 fixture에서 다시 돌리는 절차와 통과 조건"
---

# SAR-PUBLISH-001-TESTER — 독립 QA 시나리오

결과는 [QA 보고서](../qa-reports/SAR-PUBLISH-001-TESTER-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 정리 전 상태를 만들고 실행 횟수를 세야 하므로 fixture API 호출 순서를 고정해야 한다. 테스트 레벨 lite. 제품 코드는 고치지 않는다. 운영 서버·실계정·실메일은 쓰지 않는다.

## 실행

1. `python remote.py prepare` → `python remote.py fixture`(기준 SHA 고정, 자체 Postgres 컨테이너, loopback). 새 fixture마다 `walk.py`는 한 번만 돌린다. 이유: 승인 단계가 pending gate를 소비하고 `seed`는 그것을 되돌리지 않는다.
2. `python explore.py walk.py > raw-walk.out`, `python show.py`. 62개 항목이 모두 PASS여야 한다.
3. 정책 문안: `parity.py <spec> raw`(DEV 도구 재사용)로 렌더 4개와 정본의 일치를 확인한다.
4. 3폭: `measure.py <raw 절대경로> <출력 절대경로>`. 1440·390·320에서 가로 넘침이 없어야 한다(영어 위탁 표는 표 영역 안 스크롤만 허용).
5. 회귀: `remote.py make lint|test`, `remote.py integration`. 모두 EXIT=0, `--- FAIL` 0.
6. 끝: `python explore.py stop.py`로 fixture를 멈추고 컨테이너를 내린다.

## covers와 독립 checks

| covers | 통과 조건 |
|---|---|
| RC-01 화면·경로 | 설정(KO/EN)·승인 상세에 캘린더·바쁜 시간 단어 0. `/home/calendar/connect`·`/disconnect`가 POST·GET 모두 404. 로그인 scope는 `internal/relay/google.go`의 `openid`·`email`뿐이고 offline 옵션 없음(fixture는 `/auth/google`을 막으므로 코드 읽기와 integration 단언 `google_integration_test.go:92`로 확인) |
| RC-02·RC-03 정리 | 키 있음·키 있음+revoke 실패·키 없음 3사례에서 연결 0, lookup·disclosure·키 값 0, 미전달 결과 gate `denied`, 정리 기록 한 줄이 정해진 형식이고 `@`·긴 토큰이 없다. revoke 요청은 키 있는 두 사례에서만 2건. 같은 정리를 다시 실행하면 기록 빈 줄, revoke 요청 증가 0, 상태 불변 |
| RC-04 승인 | 정리 뒤 pending `schedule.query` gate에 승인 버튼이 있고 승인이 받아진다. 승인 뒤 화면에 승인 버튼·캘린더 단어 없음. 승인 gate를 consume하면 `disclosure:false`. 정리된 gate consume은 `invalid_gate`. `schedule.commit` 동작은 integration `TestPostgresSafety` |
| RC-05 회귀 | 친구·메시지·승인 통합 검사 `--- FAIL` 0 |
| PC-01 | `/privacy`·`/terms` 비로그인·로그인 모두 200, KO/EN |
| PC-02 | 렌더 본문이 정본 4개 문안과 단어 단위로 일치, 영어판 한글 0, `{ }` 0 |
| PC-03 | 첫페이지 하단·설정·두 페이지 상호 링크 있음. `#start` 안내 한 줄은 Google 또는 이메일 로그인이 켜질 때만 렌더된다(`landing_pages.go:70`). fixture에는 로그인 설정이 없어 한 줄이 없는 것이 정상이며 `TestPolicyLinks`(make test)가 맡는다 |
| PC-04 | 고유 title, h1 1개, 3폭 가로 넘침 0 |
| 관찰만 | commit~revoke 사이 crash 시 정리 기록 줄 없음(리뷰 low): fixture에는 그 사이에서 프로세스를 끊는 지점이 없어 재현하지 않는다 |
