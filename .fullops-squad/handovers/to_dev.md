---
title: SAR-AUTOSETUP-001-V01 — 320px 명령 입력란이 마지막 줄을 가리는 V-01 수정
status: draft
updated: 2026-10-23
owner: dev
tasks: [SAR-AUTOSETUP-001-V01]
summary: 320px 명령 입력란이 마지막 줄을 가리는 V-01 수정
attempt: 0b4e9fd12ba9436dbe63c1736fbb0d0f
base: 708635278707993b3fc02af2c1052e86d27f41ab
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-V01 — 320px 명령 입력란이 마지막 줄을 가리는 V-01 수정

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-23
- From / To: coor → dev
- 상태: done
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/dev의 `internal/relay/human_pages.go` 명령 입력란 수정과 관련 테스트, 커밋·push. 추가 승인 필요: 새 색·토큰·스타일 요소 추가, 다른 화면 변경, 운영 배포
- 상위 의도: `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md` — S1(회원 화면 명령 안내)
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

designer SAR-AUTOSETUP-001-ART(main 70863527) V-01 medium: 320px에서 명령이 4줄이 되는데 `#setup-command` 입력란은 `rows="3"`이라 마지막 줄(`-App … -Lang …`)이 숨는다(scrollHeight 106 > clientHeight 83, KO·EN, claude·codex). 390px도 여유 0. 근거: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md` V-01.

## 적용 기준과 예외

- 공통 규칙 `.fullops-squad/rules/common/README.md`, 기준 main 70863527.
- 디자인 정본 D04 `docs/design-docs/mockups/service-member-ui.md`. 예외 없음.

## 먼저 읽을 문서

- `internal/relay/human_pages.go` 37행 `human-setup-command`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md` V-01
- `internal/relay/autosetup_guide_test.go`

## 해야 할 일과 파일 소유권

- [ ] `#setup-command` 입력란을 320px에서도 명령 전체가 보이게 고친다(designer 권고: `rows="4"` 또는 좁은 폭에서만 min-height). 기존 토큰·클래스 안에서, 새 스타일 요소 없음
- [ ] 회귀 테스트 1개(`autosetup_guide_test.go`에 행 수 또는 속성 확인) 추가
- [ ] 320·390px에서 scrollHeight ≤ clientHeight 확인(KO·EN, claude·codex) — 격리 로컬 relay

예상 변경: 10줄 이내.

## 완료 기준과 검증

- `go test ./internal/relay/...` 통과(Linux clone `make test` 또는 Windows go test), lint `.../1.8.3/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/dev --from 70863527` ERROR 0.
- 증명 방법: 320·390px × KO·EN × claude·codex 8조합에서 scrollHeight ≤ clientHeight 수치를 report.md에 기록.
- UI 작업: 디자인 정본 D04, 공용 `.share` 클래스 유지, 테마 변경 없음.

## 갱신할 산출물

없음

## 기대 산출물

`internal/relay/human_pages.go`, 테스트, `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-V01-test/report.md`

## 제약·협업·후속

tester SAR-AUTOSETUP-001-QA가 같은 화면으로 검수 중이다(tester 워크트리 별도, 파일 겹침 없음). 후속: designer가 320px 4조합만 재확인. 하위 에이전트 배정 안 함.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-V01-packet.json` / SHA `708635278707993b3fc02af2c1052e86d27f41ab` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `internal/relay/autosetup_guide_test.go` (direct_edit) · 줄 15 · inferred
- `internal/relay/human_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 0건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / SHA: fullops/dev / (finish 뒤 아래 커밋 SHA)
- 변경 이유: 320px에서 4줄 명령의 마지막 줄이 숨던 V-01. `#setup-command` rows 3→4.
- 검토 필요: 없음. 새 스타일 요소 없음.
- 검증: `go test ./internal/relay/...` ok. 8조합 scrollHeight 106 = clientHeight 106. [report](../docs/evaluations/qa-reports/SAR-AUTOSETUP-001-V01-test/report.md)
- 의존성 변경 없음(rows 속성만 변경).
- 남은 일: designer가 320px 4조합 재확인.

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-test/report.md`(역할별로 나눈 키는 `<키>-DEV` 그대로) 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
