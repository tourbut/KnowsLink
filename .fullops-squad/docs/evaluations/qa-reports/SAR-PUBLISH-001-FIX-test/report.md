---
title: SAR-PUBLISH-001-FIX DEV 검증 보고
status: review
updated: 2026-10-22
owner: dev
tasks: [SAR-PUBLISH-001-FIX]
summary: 시각 검수 V-01~V-03 수정(문단 링크 인라인화, 정책 하단 현재 링크 표시, 설정 링크 간격)과 3폭 재측정, 원격 Linux lint·test EXIT=0
---

# SAR-PUBLISH-001-FIX DEV 검증 보고

기준 main `179bdfe8`, 수정 커밋 `2943acc6b9e6abbc2af36088a7f6095e203a5bad`(브랜치 `fullops/dev`). designer 재확인은 coor가 배정한다.

## 변경

| ID | 수정 | 파일 |
|---|---|---|
| V-01 | `.policy p a,.policy li a,.auth p a{display:inline;min-height:0}`. 전역 `a` 44px 규칙은 그대로 두고 문단 속 링크만 예외로 좁혔다. 첫페이지 `#start` 정책 안내, 정책 본문 링크(§6 Google 링크·상호 링크)가 해당한다 | `internal/relay/member_pages.go` |
| V-02 | `.landing footer a[aria-current]{font-weight:600;text-decoration:none}`. 새 색·글꼴 없음 | `internal/relay/landing_pages.go` |
| V-03 | 설정 링크 줄 `<p class="links">` + `.links{display:flex;flex-wrap:wrap;gap:0 16px}` | `internal/relay/human_pages.go`, `member_pages.go` |
| 테스트 | 설정 링크 줄 단언을 새 마크업에 맞춤 | `internal/relay/policy_pages_test.go` |

V-04는 지시대로 수정하지 않았다. 정책 문안·D04는 건드리지 않았다.

## 검증

| 검사 | 결과 |
|---|---|
| 원격 Linux `make lint` @`2943acc6` ([remote-lint.log](remote-lint.log)) | EXIT=0 |
| 원격 Linux `make test` @`2943acc6` ([remote-test.log](remote-test.log)) | EXIT=0 |
| Windows `go test ./internal/relay -count=1` (`KNOWSLINK_UI_FIXTURE_DIR` 지정) | ok |
| `node measure.mjs`: fixture 8화면(KO/EN `#start`·`/privacy`·`/terms`·설정) × 1440·390·320, Chrome headless=new CDP | EXIT=0, [measure.json](measure.json) |

측정 결과(measure.json):
- 가로 넘침: 24건 모두 `hOverflow=false`.
- V-01: 문단 링크 `display:inline`. 줄 높이 `#start` 22px, 정책 본문 26px(수정 전 약 44px 간격).
- 44px 유지: 문단 링크를 뺀 보이는 모든 `a`의 높이가 44px 이상이다(`nonInlineBelow44` 24건 모두 빈 목록). 첫페이지 다른 링크·정책 하단·설정 링크 줄 포함.
- V-02: 하단 현재 링크 `font-weight:600`, `text-decoration:none`(KO/EN 정책 6건 × 3폭).
- V-03: 설정 링크 사이 가로 간격 16px(1440·390 KO). 320 KO와 390·320 EN에서는 마지막 링크가 다음 줄로 내려간다(줄바꿈, 넘침 없음).
- 3폭 캡처 24장 중 `#start` KO 390, 설정 EN 320 두 장을 직접 봤다. 나머지는 측정값으로 판정했다. 캡처는 로컬 임시 폴더에만 두고 커밋하지 않는다.

## 검증하지 못한 것

- 서버 경유 화면(언어 쿠키 전환·로그인 리다이렉트)은 보지 않았다. 단위 테스트가 쓴 fixture HTML을 렌더했다.
- 실제 키보드 순회·200% 확대는 보지 않았다.
- 수정 전 상태와의 자동 비교는 하지 않았다. 수정 전 값은 designer 보고(V-01 약 44px)를 근거로 했다.
- 독립 시각 재확인은 designer 몫이다.

## 의존성·산출물

의존성 변경 없음(CSS·템플릿만 수정). 갱신할 산출물 없음(D10 문서 변경 불필요).
