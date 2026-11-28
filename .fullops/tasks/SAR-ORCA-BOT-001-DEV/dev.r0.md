---
id: handover-dev
title: SAR-ORCA-BOT-001-DEV — Orca 봇 adapter와 같은 주인 자기 agent 사이 사전 승인(grant)을 구현해 Grok Bot ↔ Orca 봇이 고정 형식 메시지를 주고받게 한다
status: approved
updated: 2026-11-27
owner: dev
summary: Orca 봇 adapter와 같은 주인 자기 agent 사이 사전 승인(grant)을 구현해 Grok Bot ↔ Orca 봇이 고정 형식 메시지를 주고받게 한다
key: SAR-ORCA-BOT-001-DEV
role: dev
goal: Orca 봇 adapter와 같은 주인 자기 agent 사이 사전 승인(grant)을 구현해 Grok Bot ↔ Orca 봇이 고정 형식 메시지를 주고받게 한다
base: d8096fcada5cfce161f09fa3d578f05dd8ca045c
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-ORCA-BOT-001-DEV — Orca 봇 adapter와 같은 주인 자기 agent 사이 사전 승인(grant)을 구현해 Grok Bot ↔ Orca 봇이 고정 형식 메시지를 주고받게 한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-27
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 상위 의도 `## 범위` 1–5, 표적 테스트, 격리 relay 왕복 1회.
- 추가 승인이 필요한 행위: 운영 배포, 운영 계정 연결, 원격 push. 하지 않고 완료 보고의 `남은 일`에 적는다(lead가 쿠어를 거쳐 주인 승인을 받는다).
- 상위 의도: [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md) — 성공 기준 S1–S6 전부.

## 현재 상황과 확인 근거

- 설계 정본: orca-helm `docs/orca-bots.md` 3절(PR #10). 로컬 사본은 `C:\Users\shin\orca\workspaces\orca-helm\shrimpfish\docs\orca-bots.md`(브랜치 `tourbut/orca-bots-polish`).
- 사전 승인 한도는 사용자가 2026-11-27 확정했다: 본문 2 KB, 시간당 30건, 기한 30일, 갱신은 주인 재동의(결정 D057). 묶음은 이 dev 과제 하나다(D058).
- 지금 adapter에는 claude·codex·grok 호스트가 있다(`adapters/src/connect.ts`, `setup-connect.ts`, `inbox.ts`, `grok-wake.ts`). Orca 호스트는 없다.
- relay는 Go(`cmd/`, `internal/`), DB 마이그레이션은 `db/`(goose)다.

## 적용 기준과 예외

- 공통 기준: `.fullops/guide/rules/common/README.md`와 연결된 코딩·테스트·보안 규칙.
- 기존 발송별 승인·서명·키 경계를 완화하지 않는다. 예외는 상위 의도 `## 범위` 4의 grant 하나뿐이며 근거는 D057이다.
- 결정 ID `D057`·`D058`은 사용자 결정이다(산출물 ID D01–D13과 다르다).

## 먼저 읽을 문서

- [상위 의도 SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)
- orca-helm `docs/orca-bots.md` 3절(위 로컬 경로)과 6.4절. `lib/orca.mjs`의 `resolveOrcaBin`·`runOrca`
- `.fullops/guide/roles/dev.md`: 이 역할의 지난 결정·교훈
- `adapters/README.md`, `.fullops/docs/design-docs/architecture.md`(필요 시 확인)

## 해야 할 일과 파일 소유권

- [x] relay: 자기 agent 사전 승인 grant 저장·검사·철회(generation), 한도(스키마·2 KB·시간당 30건·30일), 그 밖은 `approval_required`. `internal/`·`db/` 마이그레이션.
- [x] 웹 "내 에이전트": 주인만 grant 만들기·끄기. agent API·`confirmed:true`로는 만들 수 없다.
- [x] adapter: `adapters/src/orca.ts`와 `plugin.js orca <cmd>` 진입점(`setup-connect`·`watch`·`receive`·`send`), 도어벨 메일(`orca orchestration send --type status`, `--item`).
- [x] Grok 쪽 skill 지시: decision → 번호 선택지 → answer, 주인 지시 → task.
- [x] 표적 테스트(adapter·relay 단위)와 격리 relay 왕복 1회.
- [x] 관련 기술 문서(architecture·interface-design·adapters README) 갱신.

## 완료 기준과 검증

- 증명 방법:
  - relay: 레포 루트에서 `go test ./cmd/... ./internal/...`(Windows는 make·race 없이). grant 없음·형식 밖·2 KB 초과·시간당 한도·기한 지남·다른 owner·철회 뒤 → `approval_required`, agent API로 grant 생성 거부(S1·S2).
  - adapter: Node 22로 해당 테스트 파일만 실행. 도어벨 본문에 메시지 본문·발신자 이름 없음, 같은 ID 한 번만, 메일 실패 뒤 다음 polling에 재시도, `--item`(S3·S4), `receive`의 `untrusted: true`(S5).
  - 격리 relay 왕복 1회: task → 도어벨 → receive, decision → answer → `--item` 메일(S6). 증거는 `tasks/SAR-ORCA-BOT-001-DEV/qa/report.md`.
- 새 LLM 리뷰, 통합 테스트 묶음, 상시 eval을 더하지 않는다(토큰 비용 제약).
- `fullops lint --from <기준 SHA>`는 이 워크트리 루트에서 실행한다. Windows에서 `product-lint: make lint` 실행 불가는 기존과 같다. 제품 검사는 격리 Linux clone에서 돌릴 수 있으면 돌리고, 못 하면 이유를 적는다.

### UI 작업의 디자인 기준과 검증

- 웹 "내 에이전트"의 grant 만들기·끄기만 해당한다. 기존 화면의 컴포넌트·토큰을 그대로 쓰고 새 스타일을 만들지 않는다. designer 시각 검수는 배정하지 않는다(D058).

## 갱신할 산출물

- route 후보 D02·D03·D10. 실제로 바뀐 것만 갱신한다(grant 정책은 D03 설계 문서, 화면·API는 해당 설계 문서).

## 기대 산출물

- 코드: `internal/`·`db/`·웹·`adapters/src/orca.ts`와 테스트, Grok skill 지시 파일.
- 증거: `tasks/SAR-ORCA-BOT-001-DEV/qa/report.md`.

## 제약·협업·후속

- 병렬 방침 on: relay·웹과 adapter·skill은 파일이 겹치지 않으면 하위 에이전트로 나눠도 된다. 쓰기 하위 에이전트가 끝날 때까지 턴을 끝내지 않는다.
- 범위 밖: 운영 배포·운영 계정 연결·원격 push(승인 필요), Grok 직접 webhook 경로 정리(orca-helm 9절 3단계), Orca 봇 채팅 안 MCP.
- 키·credential은 agent 폴더에만 둔다. argv·로그·메일·증거 파일에 넣지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `tasks/SAR-ORCA-BOT-001-DEV/qa/report.md` 링크로 대신한다. worker_done 본문은 `fullops done`이 출력한 한 줄을 그대로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치 / SHA: `fullops/dev` / 3b243345(구현)·441d1a41(vet 태그) 뒤 완료 기록 커밋. push 안 함(coordinator 지시 B).
- 변경 이유: 같은 owner 자기 agent 사이 `orca.task|decision|answer.v1`만 relay가 grant로 통과시키고, Orca 봇은 본문 없는 도어벨로 깨운다(D057).
- 검토 필요: `internal/relay/grants.go`·`public_text.go`의 grant 검사와 `kind` 서명 범위, 웹 `/home/grant-create`(5분 재확인·동의)·`/home/grant-revoke`(정리 예산).
- 검증한 것: `go test ./cmd/... ./internal/...` 통과, adapter 표적 7개·`npm run check` 통과, 격리 relay 왕복 S1–S6 통과([QA](../tasks/SAR-ORCA-BOT-001-DEV/qa/report.md)). `fullops lint` ERROR 0, make 2개는 Windows 실행 불가. 같은 SHA Linux clean clone에서 `make test` 통과, `make lint` 통과(`check_compose.py`만 host docker).
- 못 한 것: 실제 Orca 봇(쿠어) 메일, 운영 relay·계정 왕복, `fullops done`(push 포함이라 실행 안 함).
- 스스로 정한 선택: 새 wire 대신 `knowslink.text.v1` 선택 필드 `kind`(없으면 바이트 동일). 상태 JSON `Grants` 맵이라 DB 마이그레이션 없음. grant 메시지는 웹 대화 사본을 만들지 않음. 기본 agent 이름 `Orca Bot`. 되돌리기는 3b243345 revert.
- 산출물: D03 architecture·D05 interface-design·D10 module-design 새 절, adapters README·Grok skill.
- 남은 일: 운영 배포·쿠어 agent 연결·grant 생성과 운영 왕복 1회(주인 승인 필요), 원격 push.
