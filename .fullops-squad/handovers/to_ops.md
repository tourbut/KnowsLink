---
title: SAR-PUBLISH-001-RETENTION — 운영 백업 14일 자동 정리 구성과 14일 초과분 삭제(D021)
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-RETENTION]
summary: 운영 백업 14일 자동 정리 구성과 14일 초과분 삭제(D021)
attempt: 78ed9cb2768648068aa99bec686c157a
base: d5695ff1b7c6c738f527a05978e69248d1404e16
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-RETENTION — 운영 백업 14일 자동 정리 구성과 14일 초과분 삭제(D021)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D021(2026-10-22). 운영 서버에 KnowsLink 백업 14일 보존 정리(systemd timer 또는 cron, 기존 운영 방식에 맞춤)를 만들고, `knowslink-state/backups/`의 14일 초과 백업을 지운다. 지우기 전에 대상 목록(이름·날짜·크기만)을 기록한다. 가장 최근 백업과 rollback용 SAR-PUBLISH-001-OPS 백업은 남긴다. / 추가 승인 필요: KnowsLink 백업 디렉터리 밖 파일 삭제, DB·volume 변경, 다른 서비스 변경, Cloud Console, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 main 7c6ef373(SAR-PUBLISH-001 배포 완료)다. ops 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md` P-03 절: 정책 문안은 백업을 정기 교체·삭제한다고 하지만 실제로는 자동 정리가 없고 `knowslink-state/backups/`에 백업 20개(14일 초과 6개)가 쌓여 있다. 2026-10-22 사용자 결정 D021: 14일 초과 백업을 자동 삭제하고 기존 초과분도 지우며 문안을 '최대 14일 보관'으로 바꾼다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `d5695ff1`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`(P-03·P-04 절)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [x] 정리 스크립트(또는 unit)를 repo에 두고(`deploy/knowslink/` 아래, 헤더 주석 1~3줄) 운영에 설치한다. 규칙: 수정 시각 14일 초과 백업 파일만 삭제, 디렉터리 밖은 건드리지 않음, 삭제 개수만 로그.
- [x] 첫 실행 전 대상 목록을 기록하고 실행한다. 실행 뒤 남은 백업 수와 가장 오래된 날짜를 기록한다.
- [x] 다음 실행 시각(timer/cron)이 잡혔는지 확인한다.
- [x] D12 `.fullops-squad/docs/operations/ops-guide.md`에 백업 보존 14일과 정리 방법·확인 명령을 적는다.

소유 파일: `deploy/knowslink/` 아래 새 정리 파일, D12, `qa-reports/SAR-PUBLISH-001-RETENTION-test/`, 이 인박스.

## 완료 기준과 검증

- 운영에 14일 초과 백업이 없고 자동 정리가 예약돼 있다. 증거는 `qa-reports/SAR-PUBLISH-001-RETENTION-test/report.md`(비밀·개인정보 없음).
- 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/ops --from d5695ff1` ERROR 0. `work.py finish`, `worker_done`.
- 정책 문안 변경은 dev(SAR-PUBLISH-001-P03)가 병렬로 한다. 문안 배포는 이 과제가 끝난 뒤 coor가 따로 배정한다.

### UI 작업의 디자인 기준과 검증

해당 없음(운영 작업).

## 갱신할 산출물

D12 ops-guide.

## 기대 산출물

운영 정리 예약, 정리 파일, D12, 보고서.

## 제약·협업·후속

dev 문안 수정과 병렬. 파일이 겹치지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-RETENTION-packet.json` / SHA `d5695ff1b7c6c738f527a05978e69248d1404e16` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 77 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 184 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 133 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
나머지 참조 후보 5건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치 / SHA: `fullops/ops` / 완료 SHA `2641bf1f7533372650e13bd4df3d71010b8a90be`(이 보고 기록 커밋은 그 뒤).
- 변경 이유: D021. 백업 자동 정리가 없어 문안과 달랐다. 서버 사용자 cron(매일 03:17 KST)과 `deploy/knowslink/prune-backups.sh`를 만들었다. 14일 초과 백업 9개를 지웠고 11개가 남았다(가장 오래된 2026-10-10, 14일 초과 0).
- 지시와 달라진 판단: 14일을 정확히 `-mmin +20160`으로 잡아 지시서의 ‘6개’ 대신 9개를 지웠다(2026-10-07 백업 3개가 14.0~14.7일). `-mtime +14`는 15일 이상만 지워 ‘최대 14일 보관’ 문안과 맞지 않는다. 가장 최근 1개는 14일이 넘어도 남기는 규칙을 더했다.
- 검토 필요: 위 9개 기준. rollback용 `4787e85-pre-publish001` 백업은 이번에 남았지만 2026-11-05 이후에는 최신 1개가 아니면 자동 삭제된다. 더 둘 복구점은 따로 복사한다.
- 검증한 것: 임시 fixture 로컬 실행(경계·밖 파일·빈 폴더·전부 초과), 서버 `bash -n`·`--list` 목록 기록 뒤 실행, crontab 줄·cron active·다음 실행 2026-10-23 03:17 KST. 못 한 것: 예약된 자동 실행은 아직 한 번도 안 돌았다. 의존성 변경 없음(표준 bash·find·ls만 사용, 파일별 의존성 추가 0).
- FullOps 검사: HEAD `2641bf1f`, 기준 `d5695ff1`, ERROR 0·WARNING 1(`ops-guide.md` SIZE-001, 기존 초과 문서)·실행 불가 2(`make lint`·`make test`, Windows `make` 부재). 제품 코드·테스트는 바꾸지 않았다.
- 산출물·로그: [보고서](../docs/evaluations/qa-reports/SAR-PUBLISH-001-RETENTION-test/report.md), D12 `ops-guide.md`(백업 보존 14일 절), `deploy/knowslink/prune-backups.sh`. 서버: `~/.local/bin/knowslink-prune-backups.sh`, `knowslink-state/backup-prune.log`, `knowslink-state/crontab.pre-retention001`.
- 남은 일: coor가 dev SAR-PUBLISH-001-P03 문안('최대 14일 보관') 병합 후 배포를 배정한다. 2026-10-23 이후 `backup-prune.log`에 첫 자동 실행 줄이 생겼는지 확인한다. D13 `transition.md`의 기존 백업 수 설명은 소유 밖이라 고치지 않았다.
