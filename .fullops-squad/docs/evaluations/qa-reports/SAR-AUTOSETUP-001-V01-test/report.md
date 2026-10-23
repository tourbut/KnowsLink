---
title: SAR-AUTOSETUP-001-V01 검증 보고
status: draft
updated: 2026-10-23
owner: dev
tasks: [SAR-AUTOSETUP-001-V01]
summary: 320px 명령 입력란 rows 4 수정 검증
---

# SAR-AUTOSETUP-001-V01 검증

변경: `internal/relay/human_pages.go`의 `#setup-command` `rows="3"`을 `rows="4"`로 바꿨다. 새 스타일 요소는 없다.

- `go test ./internal/relay/...` (Windows, dev 워크트리): ok
- 회귀 테스트: `TestAutosetupGuideOnlyForClaudeAndCodex`가 `rows="4"`를 확인한다.

## 입력란 높이 측정

방법: `human-agents` 템플릿과 `memberStyle`을 iframe(폭 320·390)에 렌더하고 Chrome headless로 `scrollHeight/clientHeight`를 읽었다. 임시 측정 코드는 삭제했다. 운영 relay는 쓰지 않았다.

| 폭 | 언어 | claude | codex |
|---|---|---|---|
| 320 | KO | 106/106 | 106/106 |
| 320 | EN | 106/106 | 106/106 |
| 390 | KO | 106/106 | 106/106 |
| 390 | EN | 106/106 | 106/106 |

8조합 모두 scrollHeight ≤ clientHeight. 가로 넘침 없음(문서 scrollWidth 305·375).
한계: 합성 템플릿 렌더이며 실제 relay 로그인 화면은 아니다. designer가 320px 4조합을 재확인한다.
