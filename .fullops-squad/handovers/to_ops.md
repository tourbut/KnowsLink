---
title: SAR-HUMAN-PEER-001-OPS — 수락한 사람 peer·UI를 기존 운영에 안전하게 적용하고 실제 계정 연결을 인계한다
status: draft
updated: 2026-10-12
owner: ops
tasks: [SAR-HUMAN-PEER-001-OPS]
summary: 수락한 사람 peer·UI를 기존 운영에 안전하게 적용하고 실제 계정 연결을 인계한다
attempt: f7ff48dfbefc43d2ab71fae18ad90c57
base: 26094e90648aa71e8d07b381e9891acb37438730
subagent_level: standard
test_level: lite
---

# SAR-HUMAN-PEER-001-OPS — 수락 버전 운영 적용과 실제 계정 보존

- Purpose: 실제 사용자 계정에서 새 사람 peer·이름·대화 UI를 사용할 수 있게 수락된 main을 기존 운영에 적용한다.
- Task key: SAR-HUMAN-PEER-001-OPS
- From / To: coor / ops. 상태: ready. test lite / subagent standard. 직렬 배포는 추가 위임 없이 수행한다.
- 담당: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/ops, fullops/ops.
- 복귀: coor C:/Users/shin/orca/workspaces/KnowsLink/coor, term_47da6145-4e1e-4cc1-93ea-20bc0dd590b9, run_5534969da2d0. 실제 task/dispatch는 preamble을 사용한다.
- 병합 책임: coor / main. 운영 제품 고정 SHA: 26094e90648aa71e8d07b381e9891acb37438730. 준비 기록 SHA와 구분한다.

## 현재 상황과 승인 범위

사용자가 사람 단위 peer·UI 전면개편과 실제 자기 이메일 등록을 요청했고 Google 본인 로그인을 완료했다. coor 운영 GET에서 기존 회원의 활성 agent 2개, 연결 키 있는 agent 1개와 기존 관계 1개를 확인했다. 기존 운영은 이전 UI다. 이메일 원문·cookie·키·OTP를 로그에 남기지 않는다. 새 제품은 main/origin 반영 완료. exact26094e9 make lint/test native0, FullOps ERROR0/WARNING19/unavailable0이다. 원254ff6a 독립 QA/코드리뷰 critical/high0, UI12dc0c7 독립 delta·scopedV01–V03 시각 수락 완료. low3 공백과 실제keyboard/200%, scenario not_run은 보존한다.

새 계정·agent를 합성으로 만들거나 실제 Google 계정을 이메일로 병합하지 않는다. 기존 사용자/키/agent/관계/철회/메시지와 다른 공유 서비스 보존을 최우선으로 한다. 운영 적용과 검증에 필요한 비파괴 SSH/SFTP·relay 교체·새 migration은 허용한다. 배포 전 새로운 보호된 DB backup을 만들고 검증하되 raw 내용·개인정보를 출력하지 않는다. 기존 backup·컨테이너·자료·코드·파일 삭제/이동은 금지한다. 로컬 Docker 실행 금지. 운영 seed/owner-login/restore·DB 과거 덮어쓰기·구버전 writer rollback 금지. beta.sh entrypoint에 삭제가 포함되면 사용하지 말고 안전한 기존 수동 절차를 적용한다. ingress/OAuth/외부 과금/공유 서비스·키 변경은 범위 밖이다. 문제가 있으면 자동 다운그레이드 대신 현재 권한을 보호하고 ask한다.

## 적용 기준과 정본

fullops-common-0.3.3 README와 coding-style/testing/security, project.md, document-writing.md 및 AGENTS Windows 금지를 같은 base26094e9에서 읽는다. shell hook 차단을 우회하지 않는다. Env:ORCA* 같은 환경 전체 열거 금지. 자격은 기존 private .env.server에서 메모리로 읽고 값 출력 없이 알려진 SSH host key를 검증한다.

## 해야 할 일과 파일 소유권

- [ ] 기존 운영 실제 SHA/clean/health/공유 서비스·회원/agent/key/관계 상태를 읽고 짧은 기술 배포·복구 계획을 phase에 기록한다.
- [ ] DB migration00002_human_writer와 기존 backup 제한을 읽고 writer fail-closed의 실제영향과 전진 복구 방식을 확인한다. 새 backup과 검증을 먼저 수행한다. 백업 없는 migration 금지.
- [ ] exact26094e9 bundle을 private 전송·검증하고 수락 버전으로 relay만 갱신한다. 기존 서버 origin 변경 없이 진행한다. 관계 자동 승격·이메일 병합·새 실계정 생성 금지.
- [ ] migration/health/local verify/public smoke/익명 인가와 새 UI 공개 확인, 전후 기존 회원/agent/key/관계 보존을 비밀 없는 집계로 대조한다. 기존 regression 실패를 새 PASS로 덮지 않는다.
- [ ] D10/D12/D13 및 필요하면 D11의 현행 운영상태를 정확히 갱신한다. 소유 경로는 .fullops-squad/docs/operations/, contexts/ops.md, docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md, docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-OPS/와 자기 packet outcomes·archive다. 제품 코드·정책·다른 worker 기록·coor PLANS/board 수정 금지. helper는 workspace 또는 task 전용 AppData/Local/KnowsLinkDevTools 경로에 보존한다.
- [ ] 실제 Google 계정의 새 /home 흐름은 coor가 기존 로그인 탭에서 확인할 수 있게 정확한 URL/운영SHA/예상 첫이름화면을 인계한다. owner의 표시 이름/agent 명칭 선택·친구 등록은 coor 후속이다. worker가 실제 친구에게 요청/메시지 발송하지 않는다.

## 완료 기준과 검증

운영 exactSHA와 migration 버전, relay healthy, 새 공개 UI, 기존 계정/2agents/keys/legacy관계 보존, 다른 공유 서비스 영향 없음, 각 명령 native exit code와 보호된 backup 경로/검증 근거를 남긴다. 재실행은 운영 영향 검사만 수행한다. 이미 수락된 제품 전체 QA를 반복하지 않는다. 운영 실패/미확인·실계정 최종 UI 미검사를 명시한다. UI source 변경 없으므로 디자인 재설계/토큰 작업은 해당 없음. D04를 배포 smoke의 정본으로만 참고한다.

## 기대 산출물과 완료 계약

예상 문서 추가 100–200줄, 제품 변경0. D10/D12/D13 갱신, 필요한 D11 현행상태. 새 의존성 없음. 먼저 읽을 문서는 탐색패킷에서 제공한다. partial/unknown은 원문 검색으로 보완하고 path/category별 completed/no_change outcomes와 uncertainty_review를 남긴다. 완료 전문을 인박스에 쓰고 work.py finish로 archive/빈인박스를 확인한다. clean commit40SHA·FullOps 오류/경고/unavailable사유와 actual session ID를 보고한다. 원 규칙의 make는 문서만 변경하므로 Windows unavailable이면 사유를 남기고 제품 exact26094e9 native0를 재사용한다. 중간 승인 없이 범위를 끝까지 진행한다. 실제 배포 실패·백업 불가·소유권 밖 수정·정책변경만 ask한다. 완료 직전 mailbox 확인 후 현재 Task/Dispatch의 worker_done succeeded/failed를 한 번 보내고 idle한다.

## 지시 전제와 충돌 — 먼저 확인
과거 OPS/OPS-LIVE의 새 backup 없음은 migration 없는 당시 과제의 제한이다. 이번에는 새 writer migration이 있어 보호된 backup 검증이 필수다. beta.sh backup/restore-verify/deploy는 실패 시 rm·기존 검증 컨테이너 삭제를 포함하므로 호출 금지다. 삭제 없는 pg_dump와 격리 고유 검증 자원을 사용하는 수동 절차를 정한다. verify.py regression의 기존 projects baseline 실패는 새 실패와 구분하여 직접 공유 서비스 전후 상태를 확인하고 원 exit1을 보존한다. beta.sh up는 내용 확인 후 허용한다. 원문이 다른 위험을 드러내면 ask한다. keep/partial 원문은 모두 필수이며 context가 넘긴 인박스 뒤쪽도 직접 읽는다.

