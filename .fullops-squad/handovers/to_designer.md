---
title: SAR-AUTOSETUP-001-ART — 회원 화면 자동 설정 안내(D04) 시각 검수와 U8 표 갱신
status: draft
updated: 2026-10-23
owner: designer
tasks: [SAR-AUTOSETUP-001-ART]
summary: 회원 화면 자동 설정 안내(D04) 시각 검수와 U8 표 갱신
attempt: a1134cf0f55b41a39c9e2e4e455af1ad
base: 4fefdf68fa55d2405bb7aa7c94a73439bf7ab235
subagent_level: standard
test_level: lite
---

# SAR-AUTOSETUP-001-ART — 회원 화면 자동 설정 안내(D04) 시각 검수와 U8 표 갱신

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-23
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fullops/designer의 D04(`docs/design-docs/mockups/service-member-ui.md`) U8 표 갱신·qa-reports 작성, 격리 로컬 relay·합성 회원 화면 캡처, 커밋·push. 추가 승인 필요: 제품 코드 수정(결함은 보고만), 제품 규칙 변경, 운영 relay 접근·실계정 로그인
- 상위 의도: `.fullops-squad/docs/planning/intents/SAR-AUTOSETUP-001.md` — S1·S2 중 회원 화면 안내의 시각·카피 부분
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

dev SAR-AUTOSETUP-001-DEV(c74c90e1, main 8ec61ab3)가 회원 화면 `use=claude` 안내 변경, `use=codex` 선택지 추가, 6번 선택지 보조 문구 변경을 구현했다(`internal/relay/connect_guide.go`, `human_pages.go`, `i18n.go`, `i18n_en_human.go`). dev는 화면 캡처를 만들지 않았다(brief: 시각 검수는 designer 후속).

## 적용 기준과 예외

- 공통 규칙: `.fullops-squad/rules/common/README.md`, `testing.md`. 기준 main 4fefdf68.
- 디자인 정본: D04 `docs/design-docs/mockups/service-member-ui.md#claude-codecodex-자동-설정-안내--sar-autosetup-001`, 기존 회원 화면 토큰·공용 스타일.
- `AS-`·`AC-10` → `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md`.

## 먼저 읽을 문서

- D04 자동 설정 절과 U8 표
- `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md`(AC-10)
- `internal/relay/connect_guide.go`, `internal/relay/i18n.go`, `internal/relay/i18n_en_human.go`, `internal/relay/human_pages.go`
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-DEV-test/report.md`

## 해야 할 일과 파일 소유권

- [ ] 격리 로컬 relay + 합성 회원으로 `use=claude`·`use=codex`·선택지 목록 화면을 1440·390·320 폭, KO·EN으로 캡처
- [ ] D04 카피·순서·강조와 실제 화면 비교, 한 줄 명령 복사·줄바꿈·가로 넘침·대비 확인
- [ ] AC-10 판정 후 D04 U8 표 갱신(`deliverables.py --stamp`로 front matter)
- [ ] 차이는 수정 필요(dev)/수용으로 분류해 report에 기록

## 완료 기준과 검증

- 정본: `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md`. 화면×폭×언어 판정표, 캡처 경로.
- 증명 방법: 3화면(use=claude·use=codex·선택지 목록) × 3폭 × 2언어 = 18조합 판정과 U8 표 갱신 커밋.
- 통과 조건: 가로 스크롤 0, 명령 줄 복사 가능, D04와 카피 일치 또는 차이 기록.
- lint: `.../1.8.3/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 4fefdf68`, ERROR 0.
- UI 작업: 디자인 정본 = D04, 코드 수정 없음(검수만), 테마 변경 없음.

## 갱신할 산출물

D04(U8 표·자동 설정 절 검수 결과)

## 기대 산출물

`.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md`, D04 갱신

## 제약·협업·후속

tester SAR-AUTOSETUP-001-QA와 병렬, 파일 겹침 없음. 피할 스타일: 기존 회원 화면에 없는 새 색·카드 그림자·알약 버튼 제안 금지, 수정 제안은 기존 토큰 안에서. 하위 에이전트 배정 안 함.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTOSETUP-001-ART-packet.json` / SHA `4fefdf68fa55d2405bb7aa7c94a73439bf7ab235` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 179 · inferred
- `internal/relay/connect_guide.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 97, 98, 453, 526 · inferred
- `internal/relay/human_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/i18n.go` (direct_edit) · 줄 24, 28, 42, 53, 75, 83, 109, 145, 178, 191, 199, 228, 238, 287, 306 · inferred
- `internal/relay/i18n_en_human.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/i18n_test.go` (impact_check) · 줄 51, 137, 282, 288, 298, 301, 301, 343, 348, 349, 376, 384 · inferred
- `internal/relay/member_pages.go` (impact_check) · 줄 19, 29, 51, 57, 58, 68, 91, 107, 110, 118, 118, 121, 121, 135 · inferred
- `internal/relay/policy/privacy.en.html` (impact_check) · 줄 24, 26, 27, 28, 29, 30, 31, 41, 43, 44 · inferred
- `internal/relay/policy/privacy.ko.html` (impact_check) · 줄 24, 26, 27, 28, 29, 30, 31, 41, 43, 44 · inferred
- `videos/knowslink-motion/index.html` (impact_check) · 줄 171, 181 · inferred
나머지 참조 후보 0건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 5건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/product-specs/SAR-AUTOSETUP-001.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / 완료 SHA: `fullops/designer` / `4c1d598000f5b20185007415d7edf5a517f3fe0c`(기준 main `11d91648`). 이후 finish 로그 커밋이 따른다.
- 변경 이유: 회원 화면 자동 설정 안내의 시각·카피 검수와 D04 U8 표 갱신. 제품 코드는 고치지 않았다.
- 결과: 3화면×3폭×2언어 18조합 중 14 PASS, 4 FAIL(V-01, 320px claude·codex × KO·EN). D04 카피 48행 일치(차이 0), 가로 스크롤 0, 명령은 claude·codex에만 있음(AC-10 PASS), 도움말 대비 6.19:1.
- 검토 필요: V-01(medium, **수정 필요 dev**) 320px에서 명령이 4줄인데 입력란이 `rows=3`이라 마지막 줄(`-App … -Lang …`)이 가려진다. 입력란을 4행으로 늘리거나 좁은 폭에서 `min-height`를 키운다(새 토큰 없음). V-02–V-04는 low로 수용했다.
- 검증한 것과 못 한 것: 격리 relay·Postgres 17·loopback SMTP sink·합성 회원으로 Chromium 1248 CDP 실측(`4c1d5980`의 shoot.mjs·compare.mjs). 실제 Ctrl+C 클립보드·200%·키보드 순서는 못 했다. lint: HEAD `4c1d5980` ERROR 0, WARNING 1(SIZE-001, D04 924줄·기존 912줄에 12줄 추가), `product-lint`·`product-test`는 Windows에 `make`가 없어 실행 불가·생략이다. relay·DB 컨테이너는 검수 뒤 정리했다.
- 산출물·로그: [report.md](../docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md), [measure.json](../docs/evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/measure.json), D04 `service-member-ui.md`('시각 검수 결과' 소절). 의존성 변경 없음(제품 코드·패키지 파일 미변경, 검수 도구는 Node·Python 표준 기능만 사용).
- 남은 일: dev — V-01 수정(`internal/relay/human_pages.go` 입력란 행 수). 수정 뒤 designer가 320px 4조합만 재확인.

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-test/report.md`(역할별로 나눈 키는 `<키>-DEV` 그대로) 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
