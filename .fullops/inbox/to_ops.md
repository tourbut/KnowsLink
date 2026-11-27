---
id: handover-ops
title: SAR-WIN-LINT-001-OPS — Windows lead 체크아웃에서 fullops lint의 make lint·make test가 실행되고 통과하게 한다(검사 기준 유지)
status: approved
updated: 2026-11-27
owner: ops
summary: Windows lead 체크아웃에서 fullops lint의 make lint·make test가 실행되고 통과하게 한다(검사 기준 유지)
key: SAR-WIN-LINT-001-OPS
role: ops
goal: Windows lead 체크아웃에서 fullops lint의 make lint·make test가 실행되고 통과하게 한다(검사 기준 유지)
base: a09df330361b8537cc0cf2f4778840fb15d679a8
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-WIN-LINT-001-OPS — Windows lead 체크아웃에서 fullops lint의 make lint·make test가 실행되고 통과하게 한다(검사 기준 유지)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-27
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: Windows에서 `make lint`·`make test`가 돌도록 하는 도구·환경·레포 설정 정비(사용자 결정 D062).
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 검사 항목을 빼거나 약하게 하는 변경(`config/lint.json` 명령 삭제, Makefile 검사 생략 등). 필요하면 착수 전에 `ask`로 묻는다.

## 현재 상황과 확인 근거

- SAR-ORCA-BOT-001-DEV(`73f3ac3e`, 리뷰 accept)를 `fullops merge`로 병합하면, 병합 뒤 lint가 `product-lint: make lint`·`product-test: make test` 실행 불가(Windows PATH에 make 없음)로 미통과가 되고 M-07로 되돌아간다. ERROR는 0이다. 이 과제는 그 병합을 막고 있다(보류 담당 coor).
- `D:/workspace/KnowsLink/.devtools/w64devkit/bin`에 GNU Make 4.4.1이 있다. 이것을 PATH에 넣고 main 체크아웃에서 `make lint`를 돌리면 `lint-go`의 gofmt가 실패한다. main 체크아웃 작업 파일 1,359개가 `i/lf w/crlf`다. `.gitattributes`는 `* text=auto eol=lf`이지만, 그 규칙 전에 체크아웃된 파일이 그대로 남았다. `git checkout-index --force --all`은 index stat이 같아 다시 쓰지 않았다.
- 이전 기록(역할 메모): w64devkit make로 `make lint`를 돌리면 `cloudflared`에서 멈춘다. adapters npm은 Node 22(`D:/workspace/KnowsLink/.devtools/node-v22.22.2-win-x64`)가 필요하다.
- 같은 코드는 Linux clean clone(Docker `node:22`, Go 1.27.1)에서 `make lint`·`make test`를 통과했다. 근거는 dev 워크트리의 `.fullops/tasks/SAR-ORCA-BOT-001-DEV/qa/report.md`(`73f3ac3e`)다.

## 적용 기준과 예외

- 공통 기준: `.fullops/guide/rules/common/README.md`.
- 프로젝트 필수 검사는 유지한다. 완화하지 않는다. Windows에서 돌 수 없는 하위 검사가 있으면 같은 검사를 Windows에서 돌리는 방법을 찾는다. 찾지 못하면 빼지 말고 `ask`로 보고한다.

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`, `.fullops/guide/project.md`의 `## 작업 검증`
- `Makefile`, `.fullops/config/lint.json`, `.gitattributes`
- `AGENTS.md`의 `Windows 파일 작업 안전`·`작업 파일 위치` 절(삭제·재체크아웃 전에 따른다)

## 해야 할 일과 파일 소유권

- [ ] main 체크아웃(`D:/workspace/KnowsLink`)과 역할 워크트리의 CRLF 작업 파일을 LF로 다시 체크아웃한다. 깨끗한 트리에서만 하고, 방법과 전후 `git status`를 기록한다.
- [ ] `fullops lint`가 부르는 `make`가 Windows에서 찾아지게 한다. 사용자 PATH를 바꾸지 않는 방법(레포 스크립트 래퍼나 lint.json의 실행 경로 등)을 우선한다. 고른 방법과 이유를 적는다.
- [ ] `make lint`의 Windows 실패 지점(cloudflared, docker compose 등)을 하나씩 해결한다.
- [ ] `make test`를 Windows에서 실행한다(race는 cgo·gcc가 필요하다. w64devkit gcc로 되는지 확인).
- [ ] 결과를 `.fullops/guide/project.md`의 `## 작업 검증`과 역할 메모에 짧게 남긴다.

## 완료 기준과 검증

- 증명 방법: lead 체크아웃(main, `D:/workspace/KnowsLink`)에서 `fullops lint --from <기준>`을 실행하면 `product-lint`·`product-test`가 실행되고 통과한다(실행 불가 0). 출력 발췌를 `tasks/SAR-WIN-LINT-001-OPS/qa/report.md`에 붙인다.
- 검사 항목 수가 줄지 않았음을 Makefile·lint.json diff로 보인다.

## 갱신할 산출물

- 없음. 운영 문서(`project.md`의 `작업 검증`, 역할 메모)만 갱신한다.

## 기대 산출물

- 래퍼·설정 변경(있으면), `tasks/SAR-WIN-LINT-001-OPS/qa/report.md`.

## 제약·협업·후속

- 범위 밖: SAR-ORCA-BOT-001-DEV 코드 수정, 원격 push.
- 받아 둔 도구와 로그는 `D:/workspace/KnowsLink/.devtools/`(gitignore)에 둔다. 레포 밖에 만들지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `tasks/SAR-WIN-LINT-001-OPS/qa/report.md` 링크로 대신한다. worker_done 본문은 `fullops done`이 출력한 한 줄을 그대로 쓴다. 원격 push가 승인 전이면 `fullops done` 대신 로컬 커밋 뒤 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-11-27)

- 브랜치 / SHA: `fullops/ops`, 코드 `0f1be95a`·`1eedce84`. 원격 push 안 함(`fullops done` 대신 worker_done만).
- 변경 이유: `scripts/make.py` 래퍼, 0600 검사의 Linux 컨테이너 실행, CRLF 작업 파일 LF 정리. 검사 항목 변경 없음.
- 검토 필요: 이 과제 병합 때만 coor 세션 PATH 접두가 필요하다.
- 검증: ops 워크트리 `fullops lint` 통과(실행 불가 0). main 체크아웃의 전체 `fullops lint`는 병합 전이라 미실행.
- 선택·되돌리기, 로그, 남은 일: [보고서](../tasks/SAR-WIN-LINT-001-OPS/qa/report.md).

