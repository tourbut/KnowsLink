---
id: handover-dev
title: "SAR-ORCA-HOST-002-DEV — SAR-ORCA-HOST-001-DEV 리뷰 후속: watch 잠금의 pid 재사용 위험, README 줄바꿈, 공백 경로 따옴표"
status: approved
updated: 2026-11-28
owner: dev
summary: "SAR-ORCA-HOST-001-DEV 리뷰 후속: watch 잠금의 pid 재사용 위험, README 줄바꿈, 공백 경로 따옴표"
key: SAR-ORCA-HOST-002-DEV
role: dev
goal: "SAR-ORCA-HOST-001-DEV 리뷰 후속: watch 잠금의 pid 재사용 위험, README 줄바꿈, 공백 경로 따옴표"
base: 5e52f8d6876e05ec3afc8a84086ca5737caba63d
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-HOST-002-DEV — SAR-ORCA-HOST-001-DEV 리뷰 후속: watch 잠금의 pid 재사용 위험, README 줄바꿈, 공백 경로 따옴표

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: SAR-ORCA-HOST-001-DEV 리뷰의 후속 수정(D068 리뷰 accept, medium 1·low 2)과 이 PC 설치본 갱신(install 재실행).
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 새 agent 생성, grant 변경, agent 폴더 삭제. 하지 않는다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)

## 현재 상황과 확인 근거

- SAR-ORCA-HOST-001-DEV(`6943ee6d`)는 리뷰 accept 뒤 main에 로컬 병합됐다. push는 아직 하지 않았다.
- 리뷰 보고: `.fullops/tasks/SAR-ORCA-HOST-001-DEV/reviews/6943ee6de30d/report.md`
- 이 PC에는 `%LOCALAPPDATA%\KnowsLink\orca\`에 설치돼 있다. 쿠어 watch는 `HKCU Run KnowsLinkOrcaWatch`로 돈다. agent는 `agent_5513ec51367fed3748903d`다.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다. 키는 argv와 로그에 넣지 않는다.
- 새 LLM 리뷰와 통합 테스트는 만들지 않고, 표적 단위 테스트만 둔다.

## 먼저 읽을 문서

- 위 리뷰 보고
- `adapters/src/orca*.ts`, `adapters/README.md`의 Orca 절

## 해야 할 일과 파일 소유권

- [ ] **medium.** `orca-watch.lock`이 pid와 uptime만 보고 살아 있다고 판정한다. 같은 부팅에서 pid가 재사용되면 start가 생략되고, stop과 install이 무관한 프로세스를 kill할 수 있다. 고칠 것:
  - 프로세스 이미지와 명령줄이 설치된 node와 watch인지 확인한다.
  - 아니면 낡은 잠금으로 보고 지운 뒤 진행한다. kill하지 않는다.
  - watch가 정상 종료할 때 잠금을 지운다.
  - 표적 테스트를 둔다.
- [ ] **low 1.** `adapters/README.md` 457-469의 `<설치>\node.exe`에서 `\n`이 실제 줄바꿈으로 깨졌다. 바로잡는다.
- [ ] **low 2.** 공백이 있는 경로에서 toolsAdd 한 줄의 따옴표가 중첩된다. 고치고 테스트를 둔다.
- [ ] **lint 경고.** `orca.ts`가 300줄을 넘었다. 수정 범위와 겹치면 나눈다. 아니면 남은 일에 적는다.
- [ ] **설치본 갱신.** 이 PC에서 install을 다시 실행해 갱신한다. 다음을 확인한다:
  - watch가 하나만 도는지
  - agent가 그대로인지
  - 무관한 프로세스를 kill하지 않았는지

## 완료 기준과 검증

- 표적 테스트(`orca.test.js` 등)와 `fullops lint --from <base>`가 통과한다.
- `tasks/SAR-ORCA-HOST-002-DEV/qa/report.md`에 테스트 발췌와 설치본 갱신 결과를 적는다.

### UI 작업의 디자인 기준과 검증

해당 없음. adapter와 문서만 고친다.

## 갱신할 산출물

- `adapters/README.md`의 Orca 절

## 기대 산출물

- 코드, 표적 테스트, README, qa 보고

## 제약·협업·후속

- 범위 밖: MCP 업그레이드(SAR-MCP-UPGRADE-001), macOS·Linux 자동 시작
- 하위 에이전트는 쓰지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
원격 push는 승인 전이다. 그래서 `fullops done`이 push를 하면 실행하지 않는다. 로컬 커밋 뒤 첫 줄 `[완료] SAR-ORCA-HOST-002-DEV | 브랜치 fullops/dev | SHA <HEAD>`로 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (dev, 2026-11-28)

- 브랜치 / SHA: `fullops/dev` / 코드 `f673644f`, qa front matter `b55e6126`
- 변경 이유: 리뷰 medium 1·low 2 후속.
  - medium: `liveWatch`가 pid·부팅 시각·생존에 더해 holder의 이미지가 node이고 명령줄이 `orca watch`인지 확인한다(Windows `Get-CimInstance Win32_Process`, 그 밖 `ps`). 아니면 낡은 잠금으로 지우고 진행하며 kill하지 않는다. watch는 정상 종료·SIGINT·SIGTERM에서 자기 잠금만 지운다(`releaseWatch`).
  - low 1: README `<설치>\node.exe` 깨짐 3곳 복구, 잠금 설명·한계 갱신.
  - low 2: `toolsAdd()`가 따옴표를 가진 명령을 작은따옴표로 감싼다.
  - lint 경고: 잠금·설치·autostart를 `adapters/src/orca-host.ts`로 분리. `orca.ts` 521→343줄.
- 검토 필요: probe 실패(PowerShell 없음 등)는 "watch 아님"으로 본다. kill 안전 쪽이지만 그때 watch가 둘 뜰 수 있다.
- 검증: `node dist/orca.test.js` PASS, `npm run lint` 통과, `fullops lint --from 5e52f8d6` 통과(ERROR 0). 설치본 갱신: agent `agent_5513ec51367fed3748903d` 그대로, watch 1개(pid 105428), `orca start`가 running 보고, kill은 검증된 기존 watch 96412만.
- 못 한 것: 실제 로그오프·재로그온 pid 재사용 재현, macOS·Linux `ps` 경로 실행.
- 스스로 정한 선택: 생존 판정에서 폴더 인자는 맞추지 않는다(상대·대소문자 경로 문제). 되돌리기: `isWatch`에 폴더 비교 추가.
- 산출물: `.fullops/tasks/SAR-ORCA-HOST-002-DEV/qa/report.md`
- 남은 일:
  - `orca.ts` 343줄, `orca.test.ts` 329줄로 300줄 경고 남음(watch·send를 더 나누거나 테스트 분할).
  - 개발 node 버전이 설치본과 다르면 봇 MCP가 `node.exe`를 잡고 있어 install이 복사에서 실패한다(기존 동작). 설치본 node로 install을 돌려 피했다. 오류 문구가 `failed.`뿐이라 원인이 안 보인다.
