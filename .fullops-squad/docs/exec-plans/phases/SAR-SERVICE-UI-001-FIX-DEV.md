---
title: 회원 웹 글자 링크 44px 터치 영역 수정 기록
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SERVICE-UI-001-FIX-DEV]
summary: D04 9절 44px 기준을 채우지 못한 글자 링크(V-01)를 공통 CSS 한 줄로 고친 계획·변경·측정·검증과 미검증 범위를 기록한다
---

# SAR-SERVICE-UI-001-FIX-DEV — 글자 링크 44px 수정

## 기준

- 지시서 `handovers/to_dev.md`, attempt `63b11e3217af4852a1abe19766e7f2bd`, 기준 SHA `c2fd5f4c37f17c3833ea2e5d9134feeb01561257`(작업 시작 HEAD `f67ba65`).
- 결함: [SAR-SERVICE-UI-001-VISUAL 보고서](../../evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md) V-01(medium). 글자 링크 11종이 높이 21px이며 D04 9절의 44px 이상 기준에 미달했다. 원 보고서의 발견·한계는 지우지 않았다. D04는 완화하지 않았다.
- 규칙: `fullops-common-0.3.3`, `project.md`. 테스트 레벨 lite, 하위 위임 없음(작은 직렬 CSS 수정). 새 라이브러리·JS·색·기능 없음. Context7 조회는 표준 CSS만 써서 하지 않았다.

## 기술 계획

1. 링크 호출자 확인: `member_pages.go`의 모든 템플릿, `http.go` gate 템플릿, `remote_auth.go`의 OAuth 동의·반환 화면이 모두 `memberStyle` 하나를 쓴다. 따라서 공통 `a` 규칙 한 곳이 모든 호출자를 덮는다. 호출자별 class는 만들지 않는다.
2. `a{color:var(--accent)}`를 `a{display:inline-flex;align-items:center;min-height:44px;color:var(--accent)}`로 바꾼다. 줄글 속 링크는 없다. 모든 링크가 단독 줄 또는 목록·행의 항목이다.
3. 같은 값을 따로 가진 `.toc a` 규칙을 지운다. `.button`(이미 `inline-flex`·44px)과 건너뛰기 링크(`position:absolute`)는 같은 규칙과 겹쳐도 크기가 같거나 커지지 않는다. 버튼형 `button.link`(로그아웃)는 `a`가 아니므로 영향이 없다.
4. 회귀: `TestMemberUIFixtures`가 29개 화면 모두에 이 규칙이 렌더되는지 확인한다. 높이 자체는 실제 DOM으로 측정한다.

## 변경

| 파일 | 변경 |
|---|---|
| `internal/relay/member_pages.go` | 공통 `a` 규칙에 44px 터치 영역 추가, 중복된 `.toc a` 규칙 제거 |
| `internal/relay/member_ui_test.go` | fixture 렌더마다 44px 링크 규칙 존재 확인 |
| 이 문서, [measure_links.py](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-DEV/measure_links.py), [measurements.json](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-DEV/measurements.json) | 기록과 재현 가능한 측정 |

form·CSRF·owner·recent 5분·CSP 계약은 바꾸지 않았다. 템플릿 마크업을 바꾸지 않아 `member_forms.golden`은 그대로 통과한다.

## 검증

대상 제품 SHA는 `609fcfd9b5049f093cec4162a3a83e4801d2e2df`다. 이 뒤의 커밋은 문서·증거만 바꾼다.

| 검사 | 위치 | 실제 종료코드·결과 |
|---|---|---|
| `go test ./internal/relay -run TestMemberUIFixtures -count=1` | Windows Go 1.27.1 | 0 (form golden·외부 자원 금지·새 44px 규칙) |
| `gofmt -l internal/relay`, `git diff --check` | Windows | 출력 없음, 0 |
| 모든 `a[href]` 및 버튼·입력·select·summary 높이·overflow 측정 | Orca 내장 브라우저, fixture 29개 × 1440/390/320, 모든 `details` 펼침 | 0. 87개 키, 링크 198건 측정. 44px 미만 0건(링크·기타 모두), overflow 요소 0건, `scrollWidth - innerWidth` 최대 0 |
| `make install`, `make lint`, `make test` | 기존 Linux 서버, 고유 clone `/tmp/knowslink-uifix-609fcfd`, 대상 SHA detach checkout | 0, 0, 0. clone clean |

- 측정에 포함된 링크 종류는 "홈으로 돌아가기"·"홈으로 이동"·"상태 새로고침"·"작업 화면으로 돌아가기"·"다른 주소로 시작"·"로그인 화면으로 이동"·확인 필요 링크 2종·승인 요청 ID 링크·원격 동의 "내 agent 관리"·섹션 이동·건너뛰기 링크다. V-01의 11종을 모두 포함한다. 이전 측정이 버튼만 센 누락을 반복하지 않도록 `a[href]` 전체를 센다.
- 이전 `make verify-mvp`·Postgres integration은 다시 실행하지 않았다. 이 변경은 CSS 한 규칙이며 Go 로직·마크업·form·HTTP를 바꾸지 않는다. 이전 결과는 SAR-SERVICE-UI-001-DEV 문서의 `0b485ad`에 묶여 있고 이번 SHA의 결과가 아니다.
- Windows에는 `make`가 없다. Linux 결과를 Windows 결과로 대신하지 않았다.
- 서버의 이전 clone·container·임시 증거는 삭제하지 않았다. 새 clone은 위 경로에 남아 있다. 서버 접속값은 출력과 Git에 넣지 않았다.

## 미검증과 후속

- tester의 새 SHA 작은 diff 리뷰와 link·overflow 회귀, designer의 V-01 링크 직접 재측정은 후속이다. 이 기록은 시각 수락이 아니다.
- 390px에서 "취소" 버튼과 링크 사이 간격은 `p` 여백 12px과 링크 44px로 달라졌다. 실제 모습은 designer 직접 확인이 필요하다.
- 기존 low V-02–V-12는 이 과제 범위 밖이며 그대로 남아 있다.
- 실제 Google 로그인·운영 렌더는 OPS 배포 뒤 coor 확인이다. main 통합·배포는 coor 후속이다.
