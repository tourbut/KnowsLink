---
title: FULLOPS-UPDATE-1.24.0
status: draft
updated: 2026-11-19
owner: coor
tasks: [FULLOPS-UPDATE-1.24.0]
summary: 2.x 전환 준비로 Python판 1.14.1–1.24.0을 적용하고 이슈 자동 모드를 끈다
---

# FULLOPS-UPDATE-1.24.0

## 버전

- 업데이트 전 설치: Python판 `fullops-squad` 1.24.0 패키지가 캐시(`~/.claude/plugins/cache/fullops-squad/fullops-squad/1.24.0`)에 있다. hook runtime `current.json`도 1.24.0이다. Claude Code 호스트의 활성 플러그인은 `fullops@fullops` 2.2.4다. `fullops-squad` 마켓플레이스는 이미 제거되어 있다.
- 레포 적용 버전: 1.14.0 → 1.24.0.
- 대상 버전: 1.24.0. 사용자 2026-11-19 지시("fullops 2.x 전환 준비")에 따른다.

## 설치

`deps.py --host claude-code` 실행 전 `--check`는 "hook runtime missing or outdated"와 "host dependency installation incomplete"를 알렸다. `deps.py --host claude-code`는 "의존성 설치 완료"로 끝났다. 이어서 `--check --host claude-code`는 `호스트 의존성 검증 완료`를 출력했다. 갱신 대상은 `@latest` npm 도구, 공통 스킬, ponytail·mattpocock-skills·pstack이다. 2.x 플러그인은 바꾸지 않았다. `frontend.json`의 tools가 비어 있어 선택형 프런트 도구는 해당 없음이다.

## 필수 적용

`update.py --repo . --from 1.14.0`가 미충족 10건을 냈다. 모두 적용했다.

| 항목 | 변경 |
|---|---|
| 1.18.0 `## 작업 검증`·`## 자주 틀리는 것` | `project.md`. 명령은 `go test ./cmd/... ./internal/...`이다. 2026-11-19 `d79a2486`에서 실행해 모든 패키지가 `ok` 또는 `[no test files]`였다. `자주 틀리는 것`은 비워 두었다. |
| 1.18.0 verifier 절 | `rules/delegation.md`. 템플릿 본문으로 맞추고 front matter는 보존했다. |
| 1.20.0 `## 주요 흐름` | `project.md`. 표는 비워 두었다. |
| 1.20.0 `testing.md` 세트 정본 | `rules/common/testing.md`. 템플릿 본문으로 맞췄다. |
| 1.20.0 `fullops-common-0.3.7` | `rules/common/README.md`, `project.md`. |
| 1.24.0 상태 어휘·`## 위임 범위` | `PLANS.md` 첫머리에 과제 표를 새로 두었다. 기존 PLANS는 서술형이라 표가 없었다. 아래 절의 살아 있는 대기 항목 6개를 `대기(사람: …)`·`보류(…)`로 옮겼다. 상시 위임은 없다. |
| 1.24.0 `.gitignore` 두 줄 | `docs/evaluations/jev/history/`, `docs/evaluations/jev/*.prior-*.json`. |

## 릴리스 판정

| 릴리스 | 판정·근거 |
|---|---|
|1.14.1 Orca 평면 응답|해당 없음. 설정 변경 없음.|
|1.15.0 lint 단계 실행|적용. `lint/README.md`에 두 단계 실행 줄을 넣었다. 시험 명령은 `make test` 하나라 `parallel_tests`는 해당 없음이다. 역할 인박스가 비어 있어 `--level` 정리는 해당 없음이다.|
|1.15.1 모션 템플릿|해당 없음.|
|1.16.0 delegation·lint README|적용. pnpm 없음, Unity 아님, `.fullops-squad/` CRLF 커밋 0개다. 섀도 판정은 선택이라 해당 없음이다.|
|1.16.1 pstack|적용. deps가 설치·갱신했다.|
|1.17.0 `.gitkeep` 폴더|보류(사용자). `.gitkeep`만 있는 폴더 7개(adr, diagrams, exec-plans/active·completed·discarded·sprint-contracts, planning/meeting-logs)가 있다. 삭제라서 지우지 않았다. SHA 없는 hold·리뷰어 pending은 `integration.py status`의 pending이 비어 해당 없음이다. 이슈 intent는 아래 이슈 모드 절을 따른다. 레포 경로에 cmd 메타문자가 없다.|
|1.18.0 템플릿 병합|적용. 필수 항목 외에 FULLOPS.md 문서 안내 행, `handovers/_TEMPLATE.md`의 `증명 방법`·`verifier`, `lint/README.md`의 `output_tail`·`rules` 문구를 맞췄다. `_TEMPLATE.md`의 레포 고유 절(0.9.12 통합, 고정 기준과 탐색 패킷)과 `## 규칙 변경 기록`, `lint.json`은 보존했다.|
|1.19.0 `_REPORT.md`|해당 없음. 스크립트가 레포 사본을 읽지 않는다.|
|1.20.0 주요 흐름|적용. `lint/README.md` `--level full` 문장, `scenarios/README.md`, `intents/_TEMPLATE.md` `확인할 흐름`을 맞췄다. `lint.json`과 `test_level`은 바꾸지 않았다.|
|1.21.0 intent `배경`|적용.|
|1.22.0 `local/`|해당 없음. `.gitignore`에 이미 있다.|
|1.22.1|해당 없음.|
|1.23.0 옛 옵션|해당 없음. 현행 문서·스크립트에 `--handover`, `--host both`, `issue_mode.py --token` 사용처가 없다.|
|1.24.0 PLANS·gitignore·settle·템플릿|적용. `integration.py settle`은 `{"released": 0}`이다. 이미 추적 중인 `*.prior-*.json` 사본은 Git에서 빼지 않았다(삭제라서 사용자 확인 필요).|
|setup dry-run|생성 예정 0개.|

## 이슈 자동 모드

사용자 지시로 껐다. 소유 세션 `01a12dd8-…`로 `issue_mode.py disable`을 실행했다. 결과는 enabled false, expires 0이다. lease는 2026-10-12에 이미 만료됐고 poll 프로세스도 없었다.

queued 잡은 사용자가 말한 #8 하나가 아니라 #3–#8 여섯 개였다. GitHub 이슈 #1–#9는 모두 2026-10-20(#9는 10-22)에 `completed`로 닫혔다. 그래서 처리하지 않고 버렸다. 6개를 cancelled로 기록했다. 다음 poll의 `fresh()`가 닫힌 이슈에 내리는 판정과 같다. #2의 held 잡은 기존 사유를 보존했다. `status.json`의 표시 job은 이제 null이다.

## 검증

- `update.py --repo . --from 1.14.0 --verify`: `필수 적용 10/10 충족`, exit 0.
- `deps.py --check --host claude-code`: `호스트 의존성 검증 완료`.
- `go test ./cmd/... ./internal/...`: 통과.
- PLANS.md 표 검사: 모든 표의 열 수·구분선이 맞다. 사용자가 짚은 856·857줄도 오류가 없어 고치지 않았다.
- 제품 코드 변경 없음.

## 보류와 동기화

- 역할 인박스 `to_*.md` 5개는 비어 있다. 통합 pending은 0이다.
- dev 워크트리에 미추적 파일이 있다(`SAR-MOTION-INTRO-001-DEV-test/` 아래 5개 폴더, `videos/service-intro/`). 사용자 자료라 보존하고 동기화만 한다.
- 2.x 전환 뒤 이 레포의 Python판 hook·스크립트 경로는 쓰지 않는다. 다음 개발은 2.x를 로드한 새 coordinator 세션에서 한다.
