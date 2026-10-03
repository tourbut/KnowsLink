---
title: SAR-MVP-002-BOT-CATALOG-DEV 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-BOT-CATALOG-DEV, SAR-MVP-002-BOT-CATALOG-DEV-REVIEW]
summary: 실제 Bot 앱 Command MCP 등록 준비의 독립 리뷰 결과를 기록한다
---

# SAR-MVP-002-BOT-CATALOG-DEV 리뷰

- 검토자 / CLI / 모델: ops 검토 세션 `103a7c86-8e5c-415a-976d-a7f06607a499` / Claude Code / claude-sonnet-5-5. 구현자 세션 `fcfec882-634c-458a-b16b-b2b285aea791`과 다르다.
- base SHA / head SHA / merge-base: `8c95bde1240bf52819ce987b2d387d437610c6e8` / `8e46c5a846e6d190e484e48be40b3dc368001a2b` / `8c95bde1240bf52819ce987b2d387d437610c6e8`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11. `review/rule.json`(`rule_sha256` ab2116fb…), `fullops-common-0.3.2`(`.fullops-squad/rules/common/README.md`와 연결 세 규칙, 기준 커밋 8c95bde), `project.md`, open-code-review-delegate 모드(OCR LLM 미사용). 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, 실행 기록 `SAR-MVP-002-BOT-CATALOG-DEV.md`, DEV 로그 `handovers/logs/2026-10-04_to_dev.md`, 이슈1 원본 `docs/evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-issue-1.json`, 공식 Grok Bot 문서(아래). Jev find·documents-find·context(REVIEW 키) 전부 읽음. 모두 참고용이다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 18 / 12 / 6 / 18 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 0 / 0. 깨끗한 8e46c5a 클론(`npm ci` 후)에서 `lint.py --from 8c95bde`, product-lint `make lint` 통과. 첫 시도는 node_modules가 없는 클론이라 `prettier: not found`로 ERROR 1이었다. 의존성 설치 후 재실행했고 이 값이 `lint.json`이다.

snapshot `/tmp/knowslink-bot-catalog-review-8e46c5a`: HEAD `8e46c5a8…`, detached, `git status --porcelain` 0줄. 시작과 끝에 확인했고 쓰기·npm 설치는 하지 않았다. 실행 시험은 별도 scratch 클론(스크래치 폴더)에서 했다.

## 검토 범위

result.json의 18개 `(path, status)`를 모두 reviewed로 기록했다. 문서 8개와 코드 4개는 diff 전체를 읽었다. 제외된 증거 txt 6개는 내용을 읽고 독립 재현 결과와 대조했다. 대조 방법은 각 파일 `reason`에 있다. 삭제된 `handovers/to_dev.md`는 logs 보존본과 일치함을 확인했다.

공식 근거는 2026-10-04에 `https://docs.x.ai/grok-bot/{team-bots,computer-and-apps,settings-and-notifications,teams-and-enterprises,approvals-security-and-privacy,faq,troubleshooting,bots,computers,security-faq}.md`를 다시 받아 대조했다. Node 체크섬은 `https://nodejs.org/dist/v22.22.2/SHASUMS256.txt`와 대조했다.

## 판단 요약

| 질문 | 판단 |
|---|---|
| Team Bots의 Command 지원 | 공식 근거 있음. `team-bots.md` Plugins 표·문단이 Custom MCP server **Command**를 명시하고 Bot 컴퓨터에서 실행됨·환경 변수 필요 시 owner 채팅 한정을 설명한다. 발췌와 문서 설명은 원문과 일치한다. |
| 개인 계정 UI의 Command 지원·등록 경로 | 공식 문서에 없음. 미확인. README는 근거 범위를 Team Bots로 한정하지 않았다(F-03). |
| 자연어 Add 요청 → 승인 카드 | 공식 근거 없음. 제3자 글만 있다. 이슈1 첫 보고에서 같은 Bot에 `AddMcpServer`가 없었다. 신규 문서는 이 기존 이슈를 적지 않았다(F-01). “반드시 카드가 나온다”는 문장은 없으나 ops-guide와 댓글 초안에는 단서가 없다. |
| 카탈로그 미노출 원인 | 가설. 공식 문서는 계정 전체 설치만 보증한다. 실행 기록 표는 가설로 쓰나 README·D10·D13·댓글 초안은 사실로 쓴다(F-02). |
| 앱 증상 해결 주장 | 과장 없음. CLI doctor·로컬 status·env -i 성공을 앱 해결로 표시하지 않고 실제 계정 재시험·이슈 열어 둠을 유지한다. |
| 기존 이슈 구분 | AddMcpServer 도구 부재(첫 보고)와 InstallPlugin 등 부재의 원인은 미확정이다. CLI 설치와 앱 카탈로그를 분리한 점은 정확하다. |

## 발견 사항

critical/high 없음. 모두 미해결(제품 수정 금지 범위)이며 수정 상태는 DEV 후속 요청이다.

| ID | 심각도 | 파일·줄 | 재현·근거 | 영향 |
|---|---|---|---|---|
| F-01 | medium | `adapters/README.md:57`, `.fullops-squad/docs/operations/ops-guide.md:193`, 실행 기록 125·138행 | 이슈1 첫 보고: `AddMcpServer`·`InstallPlugin` 등 호출 불가, 검색 0건. 공식 문서 10종에 채팅 요청·**Add MCP Server** 카드 없음. 근거는 제3자 글뿐(실행 기록 48행). README 61행은 카드의 Command 종류 부재만 다룬다. | owner가 자연어 요청으로 카드가 반드시 생긴다고 읽을 수 있다. AddMcpServer가 없으면 2단계가 막히고 원인 진단이 어렵다. |
| F-02 | medium | `adapters/README.md:38`, `design-docs/module-design.md:56`, `operations/transition.md:54`, 실행 기록 50·86행, DEV 로그 125행 | 공식 `computer-and-apps.md`: Installed connectors are account-wide. 앱 카탈로그가 등록된 connector만 읽는다·CLI plugin을 읽지 않는다는 문장 없음. 가설 1은 앱에서 반증·재현하지 않음. | 틀린 원인 단정이 이슈 댓글(게시 예정)에 들어간다. 재시험 실패 시 기록이 오도한다. |
| F-03 | low | `adapters/README.md:38` | Command 근거는 `team-bots.md#plugins`. 개인 계정 문서 없음. | 개인 owner가 Command 종류가 있다고 기대한다. 후속 분기는 실행 기록에만 있다. |
| F-04 | low | `adapters/README.md:59`, 실행 기록 129행 | `settings-and-notifications.md`: Plugins are not a settings section. | 존재하지 않는 UI 경로를 안내한다. |
| F-05 | low | `adapters/README.md:48` | `computer-and-apps.md`: Files … are designed to survive; manually installed packages are replaceable; Reset은 최근 변경 손실. 실행 기록은 미검증으로 분류. | 영속성 단정. 재실행으로 복구 가능함을 적지 않았다. |
| F-06 | low | `scripts/install_bot_mcp.sh:35-38` | `$PREFIX/package/knowslink/STALE`을 만든 뒤 재실행하면 exit0이고 `$PREFIX/knowslink/STALE`이 남는다. 기존 bundle을 먼저 지워 추출 실패 시 서버가 깨진다. | 오염된 bundle 배포 가능. 낮은 확률. |
| F-07 | low | `scripts/install_bot_mcp.sh:6,44-47` | `KNOWSLINK_PREFIX=../rel/pfx` → exit0, `Command: ../rel/pfx/node/bin/node`. | 앱에 상대 경로가 등록된다. |
| F-08 | low | `scripts/install_bot_mcp.sh:18-27`, README 48행 | `$PREFIX/node`가 파일이면 exit0로 디렉터리로 교체된다. 재사용은 버전 문자열만 본다. 체크섬 실패 시 tar.gz가 남는다. | PREFIX 안 사용자 파일 소실. README에 삭제 범위 없음. |

## 검증 및 남은 제약

독립 확인(전부 8e46c5a 클론, 스크래치):

- `sh scripts/install_bot_mcp.sh`: exit0, tar.gz checksum OK, PASS 2줄, 등록 값 출력, ZIP SHA256 `b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad`로 DEV 값과 동일.
- 재실행 exit0, 다운로드·체크섬 줄 없음(재사용). 잘못된 체크섬은 exit1이며 `$PREFIX/node` 미생성. 클론 밖 cwd 실행은 `npm ci` 오류로 exit1. 미지원 플랫폼은 DEV 로그와 스크립트 case로 확인했다.
- `set -eu`로 `env -i` 경계 검사 실패가 종료 코드로 전파됨을 확인(`sh -c 'set -e; env -i false'` exit1). 테스트는 `process.execPath`와 명시 env로 서버를 띄우므로 PATH·상속 env 없이 실행된다.
- 고정 체크섬: `978978a6…`(x64)·`b2f3a96f…`(arm64)가 공식 SHASUMS256과 일치. `.tar.gz`·`--strip-components=1`은 xz 불필요. 필요한 도구는 curl·tar·gzip·sha256sum·python3·npm 레지스트리 접근이다.
- `make lint-config` exit0(`sh -n` 포함). 실제 grok 1.0.46으로 `verify_grok_plugin.py` PASS, 실패 shim은 grok 출력 표시 후 exit1. `GROK_CONFIG`·`GROK_CONFIG_PATH`·`GROK_MANAGED_CONFIG_URL` 문자열이 grok 바이너리에 있음.
- lint: ERROR 0·WARNING 0·실행 불가 0.

실행하지 못한 검증: 실제 Bot 계정 등록·승인 카드·앱 카탈로그·Bot 세션 `knowslink_status`(계정 접근 없음, 사용자 범위 밖), 실제 aarch64 호스트, `/workspace` 유지, hosted Node. 위 항목은 owner의 재시험이 판정한다. 실제 relay·`knowslink_pull_once`·유료 inference·FullOps 업데이트는 실행하지 않았다. D01-D13은 갱신하지 않았다. 제품 코드·외부 댓글은 수정·게시하지 않았다.

## 검토 결론

조건부 수락한다. critical/high는 없다. 코드의 설치·체크섬·실패 전파·`env -i` 검증은 독립 재현됐고 CLI·로컬 성공을 앱 성공으로 과장하지 않는다. 공식 근거는 Team Bots의 Command 지원에 한정되며 개인 UI·채팅 Add 요청·AddMcpServer는 미확정이다. 게시 전 F-01·F-02 문구를 완화하고 F-04를 지우기를 권고한다. F-06~F-08은 후속 개선이다. check 통과는 기록 검사이며 앱 노출을 보증하지 않는다.
