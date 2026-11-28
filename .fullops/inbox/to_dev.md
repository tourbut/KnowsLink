---
id: handover-dev
title: "SAR-GROK-INSTALL-001-DEV — Grok Bot을 git clone 없이 한 줄 명령으로 설치·갱신한다(체크섬 검증, 웹 '내 에이전트 → Grok Bot'에 명령 표시, README 갱신)"
status: approved
updated: 2026-11-28
owner: dev
summary: "Grok Bot을 git clone 없이 한 줄 명령으로 설치·갱신한다(체크섬 검증, 웹 '내 에이전트 → Grok Bot'에 명령 표시, README 갱신)"
key: SAR-GROK-INSTALL-001-DEV
role: dev
goal: "Grok Bot을 git clone 없이 한 줄 명령으로 설치·갱신한다(체크섬 검증, 웹 '내 에이전트 → Grok Bot'에 명령 표시, README 갱신)"
base: b3b9cc38a206a94600a8740c0642a91a292caf3c
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-GROK-INSTALL-001-DEV — Grok Bot을 git clone 없이 한 줄 명령으로 설치·갱신한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인 제안(쿠어 전달 msg_c9d6eebb7698, 2026-11-28). Grok Bot 설치를 레포 clone과 `scripts/install_bot_mcp.sh` 없이 한 줄 명령으로 바꾼다. 되돌릴 수 있는 기술 결정은 추천안으로 정한다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포(운영 `knowslink-state/install/`에 파일 올리기 포함), GitHub Release·CI secret 생성. 하지 않고 완료 보고의 `남은 일`에 적는다.
- 상위 의도: 없음(단일 과제). 배경은 [SAR-ORCA-BOT-001](../docs/planning/intents/SAR-ORCA-BOT-001.md)의 운영 왕복이 Grok Bot 설치에서 막힌 것이다.

## 현재 상황과 확인 근거

- 지금 Grok Bot 설치: Bot 컴퓨터에서 레포를 clone하고 `sh scripts/install_bot_mcp.sh`를 실행한다. Node `v22.22.2`를 고정 SHA256으로 받고, bundle을 빌드·검증해 `/workspace/.knowslink/{node,knowslink}`로 교체한다(`adapters/README.md`의 `Grok Bot 앱 등록` 절).
- Claude Code·Codex는 이미 한 줄 설치가 있다: `irm https://link.knowslog.com/install/0.1.1/knowslink-setup.ps1`(루트 `README.md` 210행 부근, SAR-AUTOSETUP-001). relay가 운영 state의 `install/<버전>/`에서 내려준다(D12 SAR-VULTR-MIGRATE-001 절: state에 `install/`(0.1.0, 0.1.1)).
- 주인·쿠어 추천안:
  - 버전별 Grok bundle과 SHA256SUMS를 `link.knowslog.com`(또는 GitHub Release)에서 받게 한다.
  - 웹 '내 에이전트 → Grok Bot'에 `curl -fsSL https://link.knowslog.com/install/grok.sh | sh` 같은 한 줄을 보인다.
  - 같은 명령으로 갱신한다.
  - README의 Grok 신규 설치·업데이트 절을 그 한 줄 기준으로 고친다.

## 적용 기준과 예외

- 공통 기준: `.fullops/guide/rules/common/README.md`.
- 보안 기준은 기존 Claude/Codex installer와 같게 둔다. 버전을 고정한 URL을 쓰고, 받은 bundle은 체크섬을 검증하며, 실패하면 기존 설치를 그대로 둔다. `~/.grok`과 앱 자격 증명은 바꾸지 않는다. Node SHA256 고정을 유지한다.
- 새 공용 서비스와 새 의존성은 피한다. 기존 `install/<버전>/` 배포 경로를 먼저 재사용한다. CI·GitHub Release가 꼭 필요하면 이유를 적고, 생성 자체는 승인 대상으로 남긴다.

## 먼저 읽을 문서

- `adapters/README.md`의 `Grok Bot 앱 등록` 절, `scripts/install_bot_mcp.sh`, `scripts/package_plugin.py`
- 루트 `README.md`의 Grok Bot 절과 Claude Code·Codex 한 줄 설치 절
- 기존 installer가 운영 state의 `install/`에서 어떻게 서빙되는지(relay 코드, D12)
- [SAR-AUTOSETUP-001 intent](../docs/planning/intents/SAR-AUTOSETUP-001.md) Q1 결정(한 줄·고정 버전·checksum)
- `.fullops/guide/roles/dev.md`

## 해야 할 일과 파일 소유권

- [ ] Grok installer `install/<버전>/knowslink-grok.sh`(이름은 dev가 정함)와 Grok bundle 산출물, `SHA256SUMS`를 만드는 빌드 명령. 레포 clone·make·Go 없이 `curl`·`tar`·`gzip`·`sha256sum`(또는 `shasum`)·`sh`만 쓴다. 신규 설치와 갱신이 같은 명령이다. 실패하면 기존 설치를 보존한다.
- [ ] `install_bot_mcp.sh`의 Node 고정·`env -i` 경계 검사·prefix 규칙을 installer로 옮기거나 공유한다. 중복을 만들지 않는다.
- [ ] 웹 '내 에이전트 → 어디에서 사용할까요? → Grok Bot'에 한 줄 명령과 갱신 안내를 보인다(ko·en). 기존 Claude/Codex 표시 패턴을 쓴다.
- [ ] 루트 README와 `adapters/README.md`의 Grok 신규 설치·업데이트 절을 한 줄 기준으로 고친다. clone 방식은 개발자용으로만 남긴다.
- [ ] 운영 반영 절차(어떤 파일을 운영 state `install/<버전>/`에 어떻게 올리는지)를 D12 또는 README 운영 절에 적는다. 실행은 하지 않는다.

## 완료 기준과 검증

- 증명 방법: 깨끗한 Linux 컨테이너(Grok Bot 컴퓨터 대용, x86_64·가능하면 aarch64)에서 로컬로 서빙한 installer 한 줄을 실행해 다음을 보인다. 증거는 `tasks/SAR-GROK-INSTALL-001-DEV/qa/report.md`에 남긴다.
  1. 신규 설치가 된다. `env -i <prefix>/node/bin/node --version` 결과가 `v22.22.2`이고 `<prefix>/knowslink/dist/plugin.js`가 있다.
  2. 같은 명령을 다시 실행하면 갱신으로 처리된다.
  3. 체크섬 불일치·다운로드 실패에서는 기존 설치가 그대로 남는다.
- 표적 테스트만 돌린다. 새 통합 테스트 묶음이나 상시 eval은 더하지 않는다. `fullops lint --from <기준>`을 이 워크트리에서 통과시킨다(Windows는 `python3 scripts/make.py` 래퍼, Docker Desktop 필요).

### UI 작업의 디자인 기준과 검증

- 기존 '어디에서 사용할까요?' Claude Code·Codex 선택지의 명령 블록 컴포넌트와 문구 패턴을 그대로 쓴다. 새 스타일을 만들지 않는다. designer 검수는 배정하지 않는다.

## 갱신할 산출물

- route 후보 D10(구현), D11(사용자 안내), D03(설계). 실제로 바뀐 것만 갱신한다.

## 기대 산출물

- installer·빌드 스크립트·웹 문구·README·운영 반영 절차. 증거는 `tasks/SAR-GROK-INSTALL-001-DEV/qa/report.md`.

## 제약·협업·후속

- 범위 밖: Grok Bot 실제 계정 설치·연결(SAR-ORCA-BOT-001-OPS가 이 과제 뒤에 이어서 한다), 운영 배포, Remote HTTPS MCP 전환, Marketplace 등록.
- ops 인박스는 SAR-ORCA-BOT-001-OPS가 쓰고 있다(대기 중). 운영 파일은 건드리지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `tasks/SAR-GROK-INSTALL-001-DEV/qa/report.md` 링크로 대신한다. 원격 push는 승인 전이므로 `fullops done` 대신 로컬 커밋 뒤 worker_done만 보낸다. 첫 줄은 `[완료] SAR-GROK-INSTALL-001-DEV | 브랜치 fullops/dev | SHA <HEAD>` 한 줄이다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
