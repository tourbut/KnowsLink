---
id: handover-ops
title: SAR-ORCA-BOT-001-OPS — b4c76b64를 운영에 배포하고 쿠어 (Orca) agent를 운영에 연결해 Grok과 왕복 1회를 확인한다
status: approved
updated: 2026-11-28
owner: ops
summary: b4c76b64를 운영에 배포하고 쿠어 (Orca) agent를 운영에 연결해 Grok과 왕복 1회를 확인한다
key: SAR-ORCA-BOT-001-OPS
role: ops
goal: b4c76b64를 운영에 배포하고 쿠어 (Orca) agent를 운영에 연결해 Grok과 왕복 1회를 확인한다
base: e82012a9df887634e401600beb1d1199b44b0799
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-BOT-001-OPS — b4c76b64를 운영에 배포하고 쿠어 (Orca) agent를 운영에 연결해 Grok과 왕복 1회를 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 사용자 결정 D063(2026-11-28, "셋 다 진행"). 이 순서로 한다: ① 운영 배포, ② 받는 쪽 adapter 갱신, ③ 쿠어 agent 운영 연결, grant 생성, Grok과 왕복 1회. origin push(main·역할 브랜치)도 승인됐다.
- 추가 승인이 필요한 행위: DB 마이그레이션, 데이터 정리·backfill, 운영 부하 테스트, 홈서버 stack 삭제, `.env`·Tunnel·DNS·방화벽 변경. 생기면 멈추고 `ask`로 묻는다.
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md) — 운영 계정 왕복(성공 기준 S6의 운영판).

## 현재 상황과 확인 근거

- 운영: Vultr `knowslink-relay-01`, 제품 `40aafe6f`(goose 4, Redis). 접속·경로·compose 호출은 [D12](../docs/operations/ops-guide.md)의 SAR-VULTR-MIGRATE-001·SAR-RELEASE-004-OPS 절을 따른다.
- 배포 대상 `b4c76b64`(origin/main 조상)에는 SAR-ORCA-BOT-001-DEV(grant·Orca adapter)와 SAR-WIN-LINT-001-OPS가 들어 있다. `40aafe6f..b4c76b64` 사이에 `db/` 변경이 없다. 마이그레이션 없는 배포 경로를 쓰되 직접 다시 확인한다.
- 리뷰 low 2(`.fullops/tasks/SAR-ORCA-BOT-001-DEV/reviews/73f3ac3eb3f9/report.md`): 이전 버전 adapter는 `kind`가 있는 envelope를 버린다. 받는 쪽 adapter(Grok 역할 agent와 Orca 역할 agent 양쪽)를 새 버전으로 갱신하기 전에 grant를 만들지 않는다.
- 리뷰 low 1: Grok skill의 자리 표시(`<plugin.js>`·`<this agent folder>`·`<Orca agent>`)를 채우는 방법이 아직 없다. 운영 왕복에서 실제 값을 쓰고, 쓴 값을 보고에 적는다(SAR-ORCA-BOT-002-DEV 후속의 근거가 된다).
- 격리 왕복 절차와 스크립트: `.fullops/tasks/SAR-ORCA-BOT-001-DEV/qa/report.md`, `s6.py`.

## 적용 기준과 예외

- D12 운영 절차, 공통 보안 규칙. 키·credential·`.env` 값을 argv·로그·보고·메일에 넣지 않는다.
- SSH 스크립트는 heredoc 대신 파일로 올려 실행하거나 `docker compose exec`에 `</dev/null`을 붙인다(D12 스크립트 주의).

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`
- [D12 운영자설명서](../docs/operations/ops-guide.md) 첫 세 절
- `adapters/README.md`의 Orca 절, `adapters/skills/knowslink/SKILL.md`
- orca-helm `docs/orca-bots.md` 3.3절(도어벨): `C:\Users\shin\orca\workspaces\orca-helm\shrimpfish\docs\orca-bots.md`

## 해야 할 일과 파일 소유권

- [ ] 배포 전 `beta.sh backup`과 rollback 태그. 해시·크기를 기록한다.
- [ ] `b4c76b64` 배포. `verify.py local`·`public` 통과, relay error 0줄을 확인한다.
- [ ] 받는 쪽 adapter 갱신: 이 PC의 Grok 역할 agent 설치와 새 Orca agent 폴더가 `b4c76b64`의 adapter를 쓰는지 확인한다.
- [ ] `plugin.js orca setup-connect https://<운영 URL> %LOCALAPPDATA%\KnowsLink\orca-agent --host orca`로 `쿠어 (Orca)`를 연결한다. 브라우저 Google 동의·지문 비교·이름은 사람 단계다(아래 '사람 단계').
- [ ] grant 생성: 웹 "내 에이전트"에서 주인이 만든다. Grok→쿠어는 `orca.task.v1`·`orca.answer.v1`, 쿠어→Grok은 `orca.decision.v1`이다(사람 단계).
- [ ] `watch --host orca --bot 쿠어` 실행과 운영 왕복 1회: task → 쿠어 도어벨 메일 → receive, decision → Grok → answer → `--item` 메일.
- [ ] D12에 이번 적용 절을 짧게 더한다.

## 사람 단계

사용자가 직접 해야 하는 단계가 오면 `ask`로 보낸다. 무엇을 어디서 누르는지 한두 줄로 적는다. 예: "브라우저에 열린 KnowsLink 동의 화면에서 터미널 지문 `ab12…`와 같은지 보고 '동의'". lead가 쿠어를 거쳐 주인에게 전한다. 답이 올 때까지 기다린다. 사람 단계를 자동화하지 않는다.

## 완료 기준과 검증

- 증명 방법: `tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`에 다음을 기록한다.
  - 배포 SHA, 백업 해시, `verify.py` 결과 발췌
  - 연결된 agent 이름과 ID
  - grant 3개
  - 왕복의 도어벨 argv(본문 없음 확인)
  - `receive` 출력의 `kind`·`untrusted: true`
  - answer 메일의 `--item`
- 실패하면 직전 배포(`40aafe6f`)로 되돌리고 원인을 보고한다. 이번 배포에는 마이그레이션이 없으므로 코드 복귀가 가능하다. 직접 다시 확인한다.

### UI 작업의 디자인 기준과 검증

해당 없음. 운영 작업이다.

## 갱신할 산출물

- D12 운영자설명서(이번 적용 절).

## 기대 산출물

- `tasks/SAR-ORCA-BOT-001-OPS/qa/report.md`, D12 갱신.

## 제약·협업·후속

- 범위 밖: 코드 수정(리뷰 low 1–3은 SAR-ORCA-BOT-002-DEV), Grok 직접 webhook 경로 정리.
- 쿠어 Orca 봇에게 보내는 실제 메일은 도어벨 형식뿐이다. 그 밖의 시험 메일을 보내지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `tasks/SAR-ORCA-BOT-001-OPS/qa/report.md` 링크로 대신한다. push가 승인됐으므로 `fullops done`을 실행하고, 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
