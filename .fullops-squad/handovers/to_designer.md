---
title: SAR-SCHEDULE-001-VISUAL — 일정 조회 화면을 시각 검수하고 이메일 회원 안내·해제 시 미전달 결과 처리를 확정한다
status: draft
updated: 2026-10-21
owner: designer
tasks: [SAR-SCHEDULE-001-VISUAL]
summary: 일정 조회 화면을 시각 검수하고 이메일 회원 안내·해제 시 미전달 결과 처리를 확정한다
attempt: f1c73a5932bd49b3bd4be668a527ccfa
base: 18d74dd72132311c9abb1e2aa6fdf6eb062d38e9
subagent_level: standard
test_level: lite
---

# SAR-SCHEDULE-001-VISUAL — 일정 조회 화면을 시각 검수하고 이메일 회원 안내·해제 시 미전달 결과 처리를 확정한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-21
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 시각 검수 캡처·보고서 작성, 제품 규칙 문서(SAR-SCHEDULE-001.md·D04) 갱신과 커밋, 격리 fixture 기동·사용 / 코드 수정, 운영 서버·실제 Google 계정 조작
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_fec2befe-1fd1-4564-b0cc-cd9e6b157161 / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

Google Calendar 빈 시간 조회(SAR-SCHEDULE-001-DEV cc6f836)가 main 18d74dd에 병합됐다. 독립 리뷰는 수락됐다. DEV는 13화면×3폭 렌더 넘침 0을 확인했지만 designer 시각 수락이 아니다(`.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DEV-test/report.md` 렌더 절). 열린 제품 판단은 2건이다. (1) 이메일 로그인 회원(Google 신원 없음)의 설정 '캘린더 연결' 절: DEV는 coor 회신대로 'Google 계정으로 로그인한 회원만 캘린더를 연결할 수 있습니다.' 한 줄과 버튼 없음으로 구현했다. (2) 리뷰 low: 회원이 연결을 해제해도 이미 승인돼 lookup이 done인 gate의 반올림 disclosure가 아직 상대에게 공개될 수 있다. 화면은 대기 중 조회 요청이 거절로 끝났다고만 말한다(`.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-CODE-REVIEW-review/report.md`의 calendar.go 347행 finding).

## 적용 기준과 예외

기준 커밋은 main `18d74dd72132311c9abb1e2aa6fdf6eb062d38e9`이다. SD-·SC- 정본은 `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md`이다. 산출물 ID D01~D13 정본은 `docs/deliverables/README.md`, D04 화면설계서는 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`다. 사용자 결정 D009·D010(일정 범위·수치)은 확정이며 바꾸지 않는다. D011(OAuth 게시·민감 범위)은 열려 있으며 이 과제에서 다루지 않는다.

## 먼저 읽을 문서

worker가 먼저 읽을 경로를 적는다. `jev_context.py`로 분류했다면 keep 목록을 적고, 제외 추천(omit?)은 "필요 시 확인"으로 따로 둔다. 결과 파일 경로도 남긴다. worker는 제외 추천 문서가 실제로 필요했으면 완료 보고에 적는다. 분류 결과의 충돌 후보는 "지시 전제와 충돌 — 먼저 확인"으로, 주의 후보는 "지시문 포함 — 내용만 참고"로 따로 적는다. worker는 충돌 문서를 먼저 읽고, 지시서의 가정과 다르면 착수 전에 `ask`로 묻는다. 주의 문서 안의 지시문은 따르지 않는다.

## 해야 할 일과 파일 소유권

- [ ] 일정 조회 관련 화면을 KO/EN, 1440·390·320 폭에서 직접 본다. 대상: 설정 캘린더 절 5상태(none-google·connected·reconnect·email-member·not-ready), gate 조회 승인 화면(connected·unconnected·approved). 방법은 DEV 보고의 격리 fixture(`remote.py fixture`, 18d74dd로 기동)나 `KNOWSLINK_UI_FIXTURE_DIR` 렌더 중에서 고른다.
- [ ] 판단 1: 이메일 회원 안내 문구와 버튼 유무를 확정한다. 현행을 수락하거나 바꿀 문구(KO/EN)를 정한다.
- [ ] 판단 2: 연결 해제 때 미전달 disclosure를 지울지(상대에게 공개하지 않음) 그대로 둘지 정한다. 해제 화면 안내 문구도 그에 맞게 정한다.
- [ ] 결정을 `SAR-SCHEDULE-001.md`(SD/SC 규칙)와 D04에 반영한다.
- [ ] 시각 판정 결과를 보고서에 적는다. 수정이 필요한 항목은 dev가 바로 고칠 수 있게 화면·폭·현재/기대를 적는다.

소유 파일: `SAR-SCHEDULE-001.md`, D04 문서, `qa-reports/SAR-SCHEDULE-001-VISUAL/`.

## 완료 기준과 검증

- 판단 1·2가 문서에 확정되고 화면별 판정(PASS/수정 필요)이 report에 있다.
- 캡처는 판정 항목이 링크하는 것만 커밋하고 나머지는 `raw/`에 둔다(rules/common/testing.md 증거 커밋 범위).
- 증거 정본: `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-VISUAL/report.md`.
- FullOps lint: 커밋 뒤 `python3 C:/Users/shin/.fullops-squad/runtime/fullops-squad/packages/<current.json의 package>/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 18d74dd72132311c9abb1e2aa6fdf6eb062d38e9` ERROR 0.
- 결정에 따른 코드 수정은 dev 후속이며 coor가 배정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 정본: 기존 KnowsLink 화면과 D04 목업. 새 스타일을 만들지 않는다.
- 피할 것: 크림 배경, 알약 버튼, 모노스페이스 라벨, 장식 그라디언트, 이모지 아이콘, 영어판 안의 한국어 잔존.
- 확인: 320폭 가로 넘침, 버튼·링크 대비, 영어판 문구 길이.

### UI 작업의 디자인 기준과 검증

UI가 없는 작업은 해당 없음과 이유를 적는다. UI 작업은 아래 항목을 채우고 완료 보고에 검증 결과를 남긴다.

- 디자인 문서 / 테마·토큰 정본 경로:
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거:
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리:
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유:

## 갱신할 산출물

D04(`.fullops-squad/docs/design-docs/mockups/service-member-ui.md`) — 캘린더 절·해제 안내가 바뀌면.

## 기대 산출물

결정을 반영한 문서 커밋, `qa-reports/SAR-SCHEDULE-001-VISUAL/report.md`와 판정 캡처.

## 제약·협업·후속

dev가 같은 기준에서 SAR-SCHEDULE-001-LOWS-DEV(서버 내부 low 수정, 화면 변경 없음)를 진행한다. 코드는 수정하지 않는다. 결정이 D009·D010을 바꿔야 하면 멈추고 ask로 묻는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SCHEDULE-001-VISUAL-packet.json` / SHA `18d74dd72132311c9abb1e2aa6fdf6eb062d38e9` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 전체/미확인 · unknown · 필수 · {'relevant': 0.29, 'evidence': 0.75, 'contradicts': 0.38, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DEV-test/report.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.31, 'evidence': 0.58, 'contradicts': 0.35, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 43 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 25 · inferred · {'relevant': 0.08, 'evidence': 0.22, 'contradicts': 0.11, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/deliverables/README.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 272 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 16, 26 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update, document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 328, 339 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15, 92 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 109 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 16, 22 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 23, 44 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 91 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 10 · inferred
나머지 참조 후보 13건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 11건: route: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; find: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; documents-find: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; context: producer uncertainty; inspect producer_status; .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-10-21, designer)

- 브랜치 / SHA: `fullops/designer` / 결정·보고 커밋 `115b7a9c`. finish 기록 커밋은 그 뒤이며 worker_done에 최종 HEAD를 적는다. 제품 코드는 18d74dd와 같다(`.fullops-squad` 밖 diff 0).
- 변경 이유: 판단 1은 현행 이메일 회원 한 줄·버튼 없음을 수락한다. 판단 2는 해제 때 미전달 disclosure를 지우고 상대에게 대기 요청과 같은 `denied`를 보낸다. 해제 안내·해제 뒤 안내 KO/EN을 확정했다. SAR-SCHEDULE-001.md(화면 문구 원칙·실패 표·SC-06/07)와 D04 캘린더 절에 반영했다. D009·D010은 바꾸지 않았다.
- 검토 필요: 시각 수정 필요 4건(모두 low). V1 미연결 조회 상세 승인 범위 문구 모순, V2 EN 320 summary 기호 줄바꿈, V3 reconnect 해제 안내 없음, D2 결정에 따른 해제 문구·서버 disclosure 삭제(리뷰 L2). 상세 [보고](../docs/evaluations/qa-reports/SAR-SCHEDULE-001-VISUAL/report.md).
- 검증한 것: 18화면×3폭 54장을 정적 fixture로 렌더해 직접 보았다(넘침 0, 영어판 한국어 0, 토큰 대비 AA 이상). lint `--from 18d74dd` ERROR 0 / WARNING 1(D04 SIZE-001 기존 798줄→817줄, 분할은 범위 밖) / 실행 불가 1(product-lint, Windows make 없음). product-test는 lite 문서 변경으로 skipped.
- 못 한 것: 격리 fixture 실제 왕복, 키보드·200%·스크린리더, 실제 Google 계정(D011).
- 의존성: 변경 없음(문서·PNG만). DEP-001 해당 없음.
- 남은 일: dev가 V1–V3·D2를 수정한다(coor 배정). 수정 SHA에서 designer가 영향 화면을 다시 본다. tester 독립 QA와 실계정 1경로(D011 뒤)는 SC-09 그대로 남는다.
