---
id: handover-ops
title: "SAR-ORCA-HOST-001-OPS — 8272fb91을 운영에 배포하고 웹 'Orca 봇 (내 컴퓨터)' 항목과 이 PC 설치본 버전을 확인한다"
status: approved
updated: 2026-11-28
owner: ops
summary: "8272fb91을 운영에 배포하고 웹 'Orca 봇 (내 컴퓨터)' 항목과 이 PC 설치본 버전을 확인한다"
key: SAR-ORCA-HOST-001-OPS
role: ops
goal: "8272fb91을 운영에 배포하고 웹 'Orca 봇 (내 컴퓨터)' 항목과 이 PC 설치본 버전을 확인한다"
base: 8272fb911e1a0cf660e80390b7e6ef4ba32e2e3c
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-HOST-001-OPS — 8272fb91을 운영에 배포하고 웹 'Orca 봇 (내 컴퓨터)' 항목과 이 PC 설치본 버전을 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인 승인(쿠어 msg_2029fb016f50, 2026-11-28 "추천대로 수정 뒤 push하고 배포해"). 다음 순서로 한다:
  1. 백업
  2. `8272fb91` 운영 배포
  3. verify
  4. 실패하면 직전 배포 `0899e5bd`로 되돌린다
  ops 브랜치 push(`fullops done`)도 승인됐다.
- 추가 승인이 필요한 행위: DB 마이그레이션, 데이터 정리, 부하 테스트, `.env`·Tunnel·DNS·방화벽 변경, 홈서버 stack 삭제. 생기면 멈추고 `ask`로 묻는다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)

## 현재 상황과 확인 근거

- 운영 서버는 Vultr `knowslink-relay-01`이고 `0899e5bd`가 돌고 있다(SAR-ORCA-BOT-001-OPS).
- 배포 대상은 `8272fb91`(origin/main)이다. 다음 두 과제가 들어 있다:
  - SAR-ORCA-HOST-001-DEV: 리뷰 accept
  - SAR-ORCA-HOST-002-DEV: 리뷰 면제 D070
- lead가 확인한 `0899e5bd..8272fb91`의 변경:
  - `db/`, `deploy/`, `compose.yaml`, `Dockerfile`, `tunnel/`: 변경 없음
  - 서버 코드: `internal/relay/human_pages.go`·`i18n_en_human.go` 웹 문구만
  마이그레이션 없는 배포 경로를 쓰되, 배포 전에 직접 다시 확인한다.
- 이 PC 설치본: `%LOCALAPPDATA%\KnowsLink\orca\`. dev가 002 빌드로 이미 갱신했다(watch pid 105428, agent `agent_5513ec51367fed3748903d`).

## 적용 기준과 예외

- D12 운영 절차와 공통 보안 규칙을 따른다. 키·credential·`.env` 값을 argv·로그·보고에 넣지 않는다.
- SSH 스크립트는 파일로 올려 실행한다(D12 스크립트 주의). 허용 규칙은 이 워크트리 `.claude/settings.local.json`에 있다. 서버 명령은 Bash 도구로 실행한다.

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`
- [D12 운영자설명서](../docs/operations/ops-guide.md)의 SAR-ORCA-BOT-001-OPS 절
- `.fullops/tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`, `deploy-b4c76b64.sh`. 같은 절차를 대상 SHA만 바꿔 쓴다.

## 해야 할 일과 파일 소유권

- [ ] 배포 전 `beta.sh backup`과 rollback 태그를 만든다. 해시를 기록한다.
- [ ] `8272fb91`을 배포한다. `verify.py local`·`public` 통과와 relay error 0줄을 확인한다.
- [ ] 공개 웹 '어디에서 사용할까요'에 'Orca 봇 (내 컴퓨터)' 항목이 보이는지 확인한다(ko·en). 로그인 없이 볼 수 있는 페이지면 curl로, 아니면 HTML 발췌로 확인한다.
- [ ] 운영 읽기 전용으로 확인한다: 쿠어와 Grok Bot agent가 revoked=false이고, grant 3개가 그대로다.
- [ ] 이 PC 설치본의 버전과 상태를 기록한다:
  - bundle SHA256
  - `plugin.js` 버전
  - watch 프로세스 1개
  - HKCU Run 값
  설치본 갱신이 필요하면 `orca install`로 한다.
- [ ] D12에 이번 배포 절을 짧게 더한다.

## 완료 기준과 검증

- `tasks/SAR-ORCA-HOST-001-OPS/qa/report.md`에 다음을 적는다:
  - 배포 SHA
  - 백업 해시
  - verify 발췌
  - 웹 항목 확인 발췌
  - agent·grant 상태
  - 설치본 버전
- 실패하면 `0899e5bd`로 되돌리고 원인을 보고한다.

### UI 작업의 디자인 기준과 검증

해당 없음. 운영 작업이다.

## 갱신할 산출물

- D12 운영자설명서(이번 배포 절)

## 기대 산출물

- `tasks/SAR-ORCA-HOST-001-OPS/qa/report.md`, D12 갱신

## 제약·협업·후속

- 범위 밖: 코드 수정, 재부팅 자동 시작 확인(사람 단계), Grok에 decision 실발송 시험

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
push가 승인됐다. 그래서 `fullops done`을 실행하고, 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치: `fullops/ops`
- SHA: `58bfdf2eff7c`(lint 통과 head). 운영 배포 SHA는 `8272fb911e1a0cf660e80390b7e6ef4ba32e2e3c`.
- 변경 이유: `8272fb91` 운영 배포(복귀 없음). 코드 변경 없이 D12, 검증 보고, 배포·상태 확인 스크립트, ops 역할 지침만 바꿨다.
- 검토 필요: D12가 상한 500줄을 넘는다(SIZE-001 경고, 702줄). 인수 내용과 달리 `connect_guide.go` 선택지 1줄도 서버 변경에 있었다.
- 검증한 것: 백업 SHA256 `2619ad7b…e90bec`, `verify.py local`·`public`, relay error 0줄, agent 2개 `revoked=false`·grant 3개, 웹 항목(container binary 문자열과 template 렌더 ko·en), 설치본 bundle이 `8272fb91` 빌드와 동일, watch 1개, HKCU Run.
- 못 한 것: 로그인한 운영 웹 화면, rollback 실행, 재부팅 자동 시작, decision 실발송.
- 선택과 되돌리는 방법: bundle이 같아 `orca install`을 실행하지 않았다. 운영 복귀는 `beta.sh deploy 0899e5bd29cc11ba4ee8e70419185c97ffa536d1`.
- 산출물·로그: [검증 보고](../tasks/SAR-ORCA-HOST-001-OPS/qa/report.md), [D12](../docs/operations/ops-guide.md).
- 남은 일: 재부팅 자동 시작 확인(사람), grant 만료 2026-12-28T05:08Z.
