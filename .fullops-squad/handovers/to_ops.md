---
title: "SAR-PUBLISH-001-OPS — 캘린더 제거·정책 페이지 운영 배포, 저장 토큰 정리 확인 후 캘린더 키 제거, live ingress /privacy·/terms"
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-OPS]
summary: "캘린더 제거·정책 페이지 운영 배포, 저장 토큰 정리 확인 후 캘린더 키 제거, live ingress /privacy·/terms"
attempt: c16edb7c37184a448fff14cb8fb59255
base: 7c6ef3734a38a8f273714aa2915dc2f14f83a7e4
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-OPS — 캘린더 제거·정책 페이지 운영 배포, 저장 토큰 정리 확인 후 캘린더 키 제거, live ingress /privacy·/terms

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D020 '배포 진행'(2026-10-22). 기존 link.knowslog.com의 KnowsLink checkout·relay를 main `7c6ef373` 제품 코드로 전환한다. 배포 전 과제 고유 Postgres backup·rollback 태그(4787e85)를 만든다. 첫 기동 정리 기록 확인 뒤 운영 env 정본의 `KNOWSLINK_CALENDAR_KEY` 줄만 지운다(이전 env는 0600으로 보관, 값은 어디에도 남기지 않는다). live 공개 ingress에 `/privacy`·`/terms` 두 경로만 추가한다(repo 후보 `deploy/knowslink/tunnel/public-ingress.yml`, 반영 전후 `scripts/check_public_ingress.py`). / 추가 승인 필요: Cloud Console 변경(coor가 사용자와 한다), 실계정 로그인·발송, 기존 volume·backup 삭제, 그 밖의 ingress·다른 서비스 변경, DB migration 발생, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 main 4787e85(서버 캘린더 연결 포함, KNOWSLINK_CALENDAR_KEY 있음)다. main 7c6ef373에 SAR-PUBLISH-001(캘린더 제거·일회 정리·/privacy·/terms, 링크 스타일 수정 포함)이 독립 리뷰·tester QA 62/62·designer 시각 검수를 통과해 있다. 2026-10-22 사용자 결정 D020 '배포 진행'. 절차 정본: D13 `.fullops-squad/docs/operations/transition.md` 맨 위 SAR-PUBLISH-001-DEV 인계 1·3·4번, D12 `.fullops-squad/docs/operations/ops-guide.md` 'SAR-PUBLISH-001 운영 후보' 절.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `7c6ef373`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-REVIEW-review/report.md`(low 4건)
- `.fullops-squad/docs/operations/transition.md`(맨 위 SAR-PUBLISH-001-DEV 인계)
- `.fullops-squad/docs/operations/ops-guide.md`(SAR-PUBLISH-001 운영 후보 절)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [ ] 배포 전: backup, rollback 태그, 현재 상태(healthz, 연결 수) 기록.
- [ ] main `7c6ef373` 제품 코드로 relay를 전환한다(키는 그대로 둔 채). 첫 기동 relay 로그의 `calendar cleanup` 한 줄을 확인하고 개수만 기록한다(토큰·이메일 값 금지). 리뷰 low: commit과 revoke 사이에서 멈추면 기록 줄이 없을 수 있다. 줄이 없으면 상태 JSON에 캘린더 연결이 0인지로 판정하고 그 판단을 적는다.
- [ ] 정리 확인 뒤 운영 env에서 `KNOWSLINK_CALENDAR_KEY`를 지우고 relay를 다시 띄운다. 재기동 뒤 정리가 다시 아무것도 하지 않는지(idempotent) 확인한다.
- [ ] live ingress에 `/privacy`·`/terms`를 반영하고 `https://link.knowslog.com/privacy`·`/terms`가 비로그인 200인지, 다른 경로 노출이 바뀌지 않았는지 확인한다.
- [ ] live 확인: `/`·`#start`의 정책 링크(로그인 활성 시 보이는 줄 포함, tester 관찰), `/home/calendar/connect` 404, healthz, 기존 친구·메시지 동작 회귀 없음(실계정 발송 없이 확인 가능한 범위).
- [ ] P-03(백업 보존 기간)과 P-04(이메일 업체: 운영 SMTP 꺼짐 확인 → 행 생략 유지)를 실제 운영 값으로 확인해 보고에 적는다. 문안 변경이 필요하면 고치지 말고 보고만 한다.
- [ ] D12·D13의 'SAR-PUBLISH-001 운영 후보(운영 미적용)'를 현행 운영 상태로 갱신한다.

소유 파일: D12, D13, `qa-reports/SAR-PUBLISH-001-OPS-test/`, 이 인박스.

## 완료 기준과 검증

- 운영이 `7c6ef373` 제품 코드이고 healthy, 캘린더 연결·토큰 0, 운영 env에 캘린더 키 없음, live `/privacy`·`/terms` 200. 실패하면 rollback 절차(D13 4번)를 따르고 보고한다.
- 증거는 `qa-reports/SAR-PUBLISH-001-OPS-test/report.md`(명령·시각·종료코드 한 줄씩, 비밀 값 없음).
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/ops --from 7c6ef373` ERROR 0. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(배포). live 화면은 링크·상태코드만 확인한다.

## 갱신할 산출물

D12 ops-guide, D13 transition.

## 기대 산출물

운영 반영, D12·D13 갱신, OPS 보고서.

## 제약·협업·후속

배포 뒤 coor가 사용자와 Cloud Console(calendar.freebusy scope 제거, Calendar API 비활성화, 브랜딩 홈페이지·개인정보처리방침·약관 URL, 앱 게시)을 한다. ops는 Cloud Console을 바꾸지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-OPS-packet.json` / SHA `7c6ef3734a38a8f273714aa2915dc2f14f83a7e4` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 15, 15, 19, 20, 20, 21, 21, 21, 21, 32, 35, 66, 66, 70, 70, 80, 88, 88, 108, 136, 145 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 106, 107, 107, 111, 117 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 197, 244, 276, 277, 277, 277, 285 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 127, 127 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 108 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 335, 335, 345 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 57, 57, 59 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 35, 120 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 9, 18, 18, 23, 24, 47, 47, 129, 217, 241, 321, 335, 359, 366, 372, 394, 394 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 18, 40, 47 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 36 · inferred
- `docs/operations.md` (document_read) · 줄 19, 19, 25, 27, 27, 27, 38, 53, 54, 54 · inferred
- `internal/relay/policy/terms.en.html` (impact_check) · 줄 1, 3, 16, 37, 42, 45 · inferred
- `internal/relay/policy_pages_test.go` (impact_check) · 줄 31, 31, 32, 32, 56, 56, 57, 57, 72, 72, 76, 76, 87, 87, 93, 93 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 1, 5, 9, 9, 18, 18, 24 · inferred
나머지 참조 후보 62건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
