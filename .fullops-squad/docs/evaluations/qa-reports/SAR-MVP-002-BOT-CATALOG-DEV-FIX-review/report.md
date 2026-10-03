---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-002-BOT-CATALOG-DEV-FIX-REVIEW]
summary: 기존 조건부 리뷰의 등록 안내와 설치 경계 보완 delta를 독립 검토한다
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX 리뷰

- 검토자 / CLI / 모델: ops 검토 세션 `71164814-f769-4f08-b4e3-c0334aa7a429` / Claude Code / claude-sonnet-5-5. 구현자 세션 `285f1146-41cc-44cf-9ff4-bbe496639170`과 다르다. 원본 리뷰 세션 `103a7c86-…`과도 다른 새 세션이다.
- base SHA / head SHA / merge-base: `8e46c5a846e6d190e484e48be40b3dc368001a2b` / `423db6a2a388ea63610462f9d3a5f4c619dd781b` / `8e46c5a846e6d190e484e48be40b3dc368001a2b`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11. `review/rule.json`(`rule_sha256` ab2116fb…), `fullops-common-0.3.2`(`rules/common/README.md`와 연결 세 규칙), `project.md`, open-code-review-delegate 모드(OCR LLM 미사용). 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, 원본 리뷰 `SAR-MVP-002-BOT-CATALOG-DEV-review/report.md`·`result.json`(96d9673), FIX 실행 기록, DEV 완료 로그, 이슈1 원본, 공식 Grok Bot 문서. Jev find·documents-find·context(`…-FIX-REVIEW-*`) 모두 읽음. 참고용이다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 37 / 30 / 7 / 37 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: 0 / 1 / 0. 깨끗한 423db6a 클론(`npm ci` 후)에서 `lint.py --from 8e46c5a`, product-lint `make lint` 통과. WARNING 1건은 `PLANS.md` 582줄(상한 500, coor 소유 기존 파일의 증가)이다.

snapshot `/tmp/knowslink-bot-catalog-review-423db6a`: HEAD `423db6a2…`, detached, `git status --porcelain` 0줄. 시작과 끝에 확인했다. snapshot에는 쓰기·설치를 하지 않았다. 실행 시험은 스크래치의 별도 클론에서 했다.

## 검토 범위

result.json의 37개 `(path, status)`를 모두 reviewed로 기록했다. 구현 변경은 `scripts/install_bot_mcp.sh`, `adapters/README.md`, D10·D12·D13, FIX 실행 기록, 이전 DEV 기록의 정정 1줄, DEV 완료 로그다. 나머지는 coor 소유 준비·병합분(PLANS·board·Jev JSON·원본 리뷰 파일·to_tester·이전 OPS 로그)이며 읽고 일치만 확인했다. 제외된 증거 txt 7개는 내용을 읽고 독립 재현과 대조했다.

delta에서 바뀌지 않은 `adapters/src`·Go 코어·UI·`Makefile`·`verify_grok_plugin.py`와 원본 CLI 검증은 `git diff --stat`로 변경 없음을 확인하고 재검증하지 않았다. 해당 근거는 8e46c5a의 원본 리뷰·DEV 기록을 재사용한다.

공식 근거는 2026-10-04에 `https://docs.x.ai/grok-bot/{team-bots,computer-and-apps,settings-and-notifications}.md`를 다시 받아 대조했다. README·FIX 기록의 인용 문장(`Installed connectors are account-wide`, `Plugins are not a settings section`, `designed to survive`, `replaceable`, `very recent changes may be lost`, Setup → Plugins → Choose Add or ask the Bot in chat, Plugins 표의 Command 행)은 원문과 일치한다.

## 원본 F-01~F-08 판정 (고정 SHA 423db6a 기준)

| ID | 원본 | 판정 | 근거 |
|---|---|---|---|
| F-01 | medium, AddMcpServer 부재·카드 전제 | **해소** | README 앱 등록 절 관측 사실, D10, D12 3단계, D13, 댓글 초안 모두 이슈1의 `AddMcpServer`·`InstallPlugin` 호출 불가·검색 0건을 적는다. 카드는 제3자 글이며 공식 계약이 아니라고 쓰고, 2단계를 Setup → 채팅 → 중단·회신의 조건부 분기로 바꿨다. 카드 노출을 기대값으로 단정하지 않는다. |
| F-02 | medium, 원인 단정 | **해소** | README·D10·D13·FIX 표·댓글 초안이 '미확정 가설'로 낮췄고 관측 사실(CLI healthy·앱 14개에 없음)과 공식 문서가 말하는 범위(account-wide)를 분리했다. 이전 DEV 기록은 본문을 보존하고 정정 안내와 '초안 미게시'를 위에 붙였다. 남은 단정 문구는 그 보존본과 아카이브 로그뿐이다. |
| F-03 | low, 근거 범위 | **해소** | README 공식 근거 범위를 Team Bots로 한정하고 개인 계정 UI·Command 종류 미확인을 명시했다. |
| F-04 | low, Settings → Plugins | **해소** | 활성 문서에서 삭제했다. 공식 `Plugins are not a settings section` 근거 링크를 달았다. 남은 언급은 정정 표와 보존 기록뿐이다. |
| F-05 | low, /workspace 단정 | **해소** | README·D12가 '유지는 보장하지 않는다'로 바꾸고 update·recover·Reset 뒤 확인 명령과 재설치 복구를 안내한다. 공식 인용과 일치한다. |
| F-06 | low, stale package 혼입·기존 bundle 선삭제 | **해소** | 독립 재현: `package/knowslink/STALE`과 기존 bundle의 STALE을 둔 prefix 재실행 exit0, 새 bundle에 STALE 0개, `package/`·`notes.txt` 보존, stage 없음. 추출 실패·`env -i` 실패 시 기존 bundle이 보존됨은 DEV 회귀 로그와 코드(교체가 검사 뒤)로 확인했다. |
| F-07 | low, 상대 prefix | **해소** | 독립 재현: `KNOWSLINK_PREFIX=rel/pfx` → exit1 `must be an absolute path`, 폴더 미생성. |
| F-08 | low, 비소유 파일 삭제·archive 잔존 | **해소(경계 있음)** | 독립 재현: `node`가 사용자 파일이면 exit1로 보존, 심볼릭 링크도 거절. archive는 stage 안에 있어 종료 때 지워진다. README가 관리 경로·거절 조건을 적는다. 경계는 아래 N-01~N-03이다. |

## 발견 사항

critical/high/medium 없음. 신규 low 3건이며 모두 미해결이다(제품 수정 금지 범위, 후속 요청).

| ID | 심각도 | 파일·줄 | 재현·근거 | 영향 |
|---|---|---|---|---|
| N-01 | low | `scripts/install_bot_mcp.sh:54`, `adapters/README.md:58`, D10 58행, FIX 기록, 댓글 초안 | Node는 체크섬·추출 직후(54행) 교체되고 bundle 빌드·package 검증·`env -i` 검사는 그 뒤다. 구 Node(v0.0.0)와 bundle이 있는 prefix에서 Node 교체 뒤 npm이 실패하게 하면(shim) exit1, Node만 새 v22.22.2, bundle은 기존이다. 문서는 '검사를 통과한 Node와 bundle만 교체'·'실패하면 기존 준비물을 그대로 둔다'고 쓴다. DEV 회귀 19개에는 이 경우가 없다(구 Node 시험은 archive 실패만, 추출·검사 실패 시험은 Node를 재사용했다). | 문구와 실제 교체 경계가 어긋난다. 새 Node는 고정 SHA256으로 검증돼 실해는 작다. |
| N-02 | low | `scripts/install_bot_mcp.sh:41-44` | `swap`은 `mv` 두 번이다. 첫 `mv` 직후 TERM을 주입하면(mv shim) trap이 `exit 1`을 거쳐 EXIT trap으로 stage를 지운다. 결과 exit1, `$PREFIX/knowslink` 없음. old와 new bundle이 함께 사라진다. 실제 창은 마이크로초이며 1단계 재실행으로 복구된다. | 신호가 이 구간에 걸리면 '실패해도 기존 보존'이 거짓이다. 확률이 매우 낮다. |
| N-03 | low | `scripts/install_bot_mcp.sh:26-31,36-38`, `adapters/README.md:58` | 비소유 판정은 마커 파일 하나다. 마커가 있는 `knowslink/`는 사용자 파일이 섞여도 통째 교체된다(`knowslink/mine.txt` 소실, 형제 `other.txt` 보존). 문서는 관리 경로를 적었으나 '폴더 안 파일도 교체'를 쓰지 않았다. SIGKILL 뒤 `.stage.XXXXXX`가 남고 다음 실행(exit0)도 지우지 않는다. | 문서 한 줄 부족과 잔존 디렉터리. 기존 bundle·등록 값에 영향 없음. |

게시 전 필수 수정: 없음. 권고 1건: 댓글 초안의 '실패하면 기존 준비물을 그대로 둡니다.'를 'bundle 단계가 실패하면 기존 bundle은 그대로 두고, Node는 체크섬 통과 뒤 교체되었을 수 있습니다.'로 고치면 N-01과 맞는다. N-02·N-03은 후속 개선이며 게시를 막지 않는다.

## 검증 및 남은 제약

독립 확인(전부 423db6a 클론, 스크래치 prefix, 사용자 `/workspace`·`~/.grok` 미사용):

- `sh scripts/install_bot_mcp.sh` 첫 실행: exit0, tar.gz checksum OK, PASS 2줄, 등록 값 절대 경로 출력, ZIP SHA256 `d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609`로 DEV 값과 동일. 실행 뒤 `.stage.*` 없음.
- 등록 command를 `/`에서 `env -i`로 JSON-RPC 호출: tools `knowslink_status`·`knowslink_pull_once`, status `held`, stderr 0.
- 경계 재현: 재실행 stale·무관 파일 보존, 상대 prefix, `node` 사용자 파일·심볼릭 링크 거절(F-06~F-08). shim으로 N-01·N-02·N-03.
- `lint.py --from 8e46c5a`: ERROR 0·WARNING 1·실행 불가 0. `make lint` passed.
- 공식 문서 3종 인용 대조 일치.

실행하지 않은 검증(미검증): 실제 Grok Bot 계정의 Command 등록, 개인 계정 Setup 화면·추가 메뉴, 승인 카드 노출, `AddMcpServer` 제공 여부, 앱 카탈로그 노출, Bot 세션의 `knowslink_status`, 실제 aarch64 호스트, update·Reset 뒤 `/workspace` 유지. 원격 계정 접근이 없어 owner의 실제 재시험이 판정한다. 실제 relay·`knowslink_pull_once`·유료 inference·FullOps 업데이트·이슈 외부 댓글은 실행하지 않았다. 전체 Go·코어·UI·CLI 테스트는 복제하지 않았다. 원본 TESTER QA는 8e46c5a 대상이라 423db6a의 installer 변경을 덮지 않으며, 이번 scratch 재현과 DEV 회귀 로그가 그 공백을 메운다. D01-D13은 갱신하지 않았다. 제품 코드·이전 리뷰는 수정하지 않았다.

## 검토 결론

수락한다. critical/high/medium이 없고 원본 F-01~F-08이 423db6a에서 모두 해소되었다. 코드의 절대 prefix·비소유 항목 거절·fresh stage·검증 뒤 교체·실패 시 bundle 보존은 독립 재현됐다. 문서는 원인을 미확정 가설로, 공식 근거를 Team Bots로 한정하고 CLI·로컬 성공을 앱 성공으로 표시하지 않는다. 남은 N-01~N-03은 low이며 N-01 문구 1건만 게시 전 고치기를 권한다. check 통과는 기록 검사이며 앱 노출을 보증하지 않는다.
