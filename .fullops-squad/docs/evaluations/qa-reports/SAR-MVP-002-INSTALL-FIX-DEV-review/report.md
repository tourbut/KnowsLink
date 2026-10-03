---
title: SAR-MVP-002-INSTALL-FIX-DEV 리뷰
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-002-INSTALL-FIX-DEV-REVIEW]
summary: Grok 설치 경로 수정 후보의 독립 리뷰 결과를 기록한다
---

# SAR-MVP-002-INSTALL-FIX-DEV 리뷰

- 검토자 / CLI / 모델: ops 독립 검토자 / Claude Code / claude-sonnet-5-5. 세션 `f0c560f5-2941-4a4d-96b1-855642ac193c`.
- 구현자 세션: `08407653-b74b-4e95-abc2-4074c29956cb`(지시서 기록). 두 세션은 다르다.
- snapshot: `/tmp/knowslink-install-review-5506d64`. HEAD `5506d646c47c2d64b35d1254ddfdec5e2003084d`, detached, 작업 트리 clean. 리뷰 중 수정·npm 설치를 하지 않았다.
- base SHA / head SHA / merge-base: `376981441d1b0650897847273e1773f8473949e5` / `5506d646c47c2d64b35d1254ddfdec5e2003084d` / `376981441d1b0650897847273e1773f8473949e5`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c). `.fullops-squad/review/rule.json` sha256 `ab2116fb…2aa0`. FullOps 0.9.13 `review.py`. fullops-common-0.3.2 README와 coding-style·testing·security, FULLOPS.md, project.md, 문서 작성 규칙을 적용했다. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, `docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md`, DEV 지시서 전문(`handovers/logs/2026-10-04_to_dev.md`), 이슈1 원본 JSON, CLI 동봉 안내서 `~/.grok/docs/user-guide/09-plugins.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 18 / 11 / 7 / 18 / 0. 제외 7개는 DEV 검증 로그 6개와 `adapters/.npmrc`이며 모두 직접 읽었다.
- lint(`lint.json`): ERROR 0 / WARNING 0 / 실행 불가 0. DEV 워크트리의 깨끗한 `5506d64`에서 `--from 3769814`로 실행했다. `product-lint: make lint` passed, exit 0.

## 검토 범위

- 제품 변경 `Makefile`, `adapters/.grok-plugin/plugin.json`, `adapters/.mcp.json`, `adapters/.npmrc`, `adapters/package.json`, `adapters/README.md`, `scripts/package_plugin.py`, `scripts/verify_grok_plugin.py`를 전부 읽었다. 관련 계약으로 `adapters/mcp.json`, `.cursor-plugin/plugin.json`, `src/mcp.ts`(held 기본값 7줄), `src/mcp.test.ts`(시험 환경·process.execPath)를 대조했다.
- 기록 변경 실행 기록·완료 로그·project.md·인박스 비움을 실제 동작과 대조했다. 인박스 비움은 `work.py finish` 결과이고 전문이 logs에 있다.
- 증거 재사용: 코어(Go)·UI·`src/*.ts` 제품 소스는 `3769814..5506d64`에서 변경이 없다. 이전 QA `0fb32cd`·리뷰 `0055a5b`의 해당 증거는 코어에만 재사용했다. 설치·패키지 변경은 아래처럼 새로 실행했다.

### 영역별 결과

| 영역 | 결과 |
|---|---|
| Grok 지원 형식 | 안내서 09-plugins.md는 `.grok-plugin/` 또는 `.claude-plugin/` manifest·`.mcp.json`·`skills/` 기본 탐색을 말한다. 실제 `grok plugin validate`가 `name knowslink`, `version 0.1.0`, `1 skill dir(s)`, `MCP servers`를 인식했다. 일치한다. |
| 이전 결함 재현 | 이전 구성(`.grok-plugin`·`.mcp.json` 제거)을 임시 HOME에 설치하면 이름 `knowslink-d99e6650`, version null, `mcpServers: 0`, `grok mcp doctor knowslink`는 `not found`다. DEV 주장이 독립 재현됐다. `validate`는 이 경우에도 exit 0이므로 음성 검사로 쓸 수 없다. verify 스크립트는 list의 name·version과 doctor를 단언한다. |
| MCP discovery·GROK_PLUGIN_ROOT | 수정 구성을 설치하면 doctor가 `plugin: knowslink` 출처, `node <설치 폴더>/dist/plugin.js`, command found(`/usr/bin/node`), server started, handshake OK(protocol 2025-11-25), `2 tools discovered`, healthy 1을 보였다. `${GROK_PLUGIN_ROOT}`가 설치 경로로 치환된다. 안내서는 이 변수를 hook 환경에서만 설명한다. MCP 치환은 실제 실행으로 확인했다. |
| trust | 안내서와 같이 `--trust` 없이는 설치가 확인 요구에서 멈춘다. README와 verify는 `--trust`를 쓴다. 임시 HOME 설치는 플러그인 서버를 활성화했다. |
| marketplace | `.grok-plugin/marketplace.json`(source `knowslink`)을 풀어 `grok plugin marketplace add <폴더>` → `grok plugin install knowslink --trust`가 exit 0이다. 설치 뒤 준비 폴더를 지워도 doctor가 healthy다. `grok plugin uninstall knowslink`도 동작했다. README 56·58줄이 맞다. |
| README 설치 명령 | 56줄까지 명령을 그대로 실행해 기대 출력이 모두 맞았다. 상대 경로 `old/knowslink`는 GitHub shorthand로 해석돼 git clone을 시도한다. README는 절대 경로 `/tmp/knowslink-plugin/...`을 쓰므로 영향이 없다. |
| 패키지·checksum | 스냅샷의 `git archive`를 스크래치 사본으로 풀어 `npm ci --prefix adapters` 뒤 `make verify-grok-plugin`을 실행했다. exit 0, ZIP SHA256 `fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75`로 DEV 값과 같다. 항목은 11개이며 credential·fixture·node_modules·.env가 없다. |
| 설치 폴더·검증 스크립트 격리 | `verify_grok_plugin.py`는 임시 HOME에서만 설치한다. 실행 전후 `~/.grok/installed-plugins` 파일 목록 해시가 같다(`526474ad…`). 내 수동 시험도 임시 HOME을 사용했다. 환경 변수 격리 범위는 F-03. |
| 실패 전파 | 음성 시험 2종: `.mcp.json`을 `dist/missing.js`로 틀리면 verify exit 1, PATH에서 grok을 빼면 `FileNotFoundError` exit 1. 둘 다 fail closed다. 진단 출력 손실은 F-01. |
| 런타임 strict | `.npmrc` `engine-strict=true`와 engines `>=22.22.2 <23`·npm `>=10.9.7 <11`은 이 호스트의 Node 22.22.2/npm 10.9.7에서 `npm ci` exit 0이다. Node 20의 EBADENGINE은 DEV 로그(exit 1)를 읽었고 Node 20이 없어 재실행하지 못했다. |
| Bot 앱 지원 주장 | README 60줄과 실행 기록 `미검증 항목`이 CLI 설치와 앱 동적 카탈로그를 구분한다. CLI 성공을 앱 설치 성공으로 주장하지 않는다. 실행 기록 34줄의 단정은 F-04. |
| 호환성 | Cursor 파일(`.cursor-plugin/plugin.json`, `mcp.json`)과 marketplace는 ZIP에 유지된다. `adapters/README.md`에서 Cursor Dashboard 팀 marketplace 안내가 사라졌다. 해당 안내는 이전 이슈에서 검증되지 않은 것이며 Cursor 로컬 검사 문장은 남았다. |

## 발견 사항

| ID | 심각도 | 파일·줄 | 설명 | 상태 |
|---|---|---|---|---|
| F-01 | low | `scripts/verify_grok_plugin.py:22-24` | `capture_output=True`와 `check=True` 조합으로 실패 시 grok의 stdout·stderr가 버려진다. 재현: `.mcp.json` 경로를 틀리게 패키지하면 exit 1이지만 doctor 사유가 출력되지 않는다. fail closed라 차단 아님. | 미해결. `CalledProcessError`의 stdout·stderr 출력 권장. |
| F-02 | low | `scripts/verify_grok_plugin.py:33-35` | PASS 문구 `installed copy held`가 과장이다. 설치본 bundle을 시험 자체의 환경과 `process.execPath`로 실행한다. Grok이 설치본 `.mcp.json` env와 PATH node로 띄운 상태는 검사하지 않는다. held는 `mcp.ts:7` 기본값이라 안전 영향은 없다. | 미해결. 문구 하향 또는 doctor의 command found 경로 단언. |
| F-03 | low | `scripts/verify_grok_plugin.py:19-20` | 격리가 HOME 재지정과 `GROK_HOME` 제거뿐이다. grok 바이너리에 `GROK_CONFIG`·`GROK_CONFIG_PATH`·`GROK_MANAGED_CONFIG_URL`이 있다. 이 호스트에는 없어 누출을 재현하지 못했다. 쓰기 영향은 미검증이다. | 미해결. 같은 변수를 `environment.pop`에 추가 권장. |
| F-04 | low | `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md:34` | `Bot 앱 도구 목록에 InstallPlugin 등이 없던 것은 앱의 동적 카탈로그 문제`를 사실로 적는다. 136줄은 같은 노출을 재시험 전까지 미판정으로 둔다. 근거는 이슈 본문의 도구 부재뿐이다. | 미해결. 추정 표현으로 낮춰 두 곳을 일치시킨다. |
| F-05 | low | `adapters/README.md:10-13` (같은 명령: 실행 기록 85-88) | Node 준비 절차가 linux-x64와 xz를 가정한다. Bot 컴퓨터의 아키텍처·xz는 확인되지 않았다. 재시험 실패 시 원인 분류를 어렵게 한다. | 미해결. 댓글에 `uname -m` 출력 요청을 넣거나 파일명을 `uname -m`으로 고른다. |

critical 0, high 0, medium 0, low 5. 모든 줄 번호는 고정 snapshot과 `git show 5506d64`의 실제 파일에서 확인했다.

## 검증 및 남은 제약

- 직접 실행(스크래치 사본, snapshot 비수정): `npm ci --prefix adapters` exit 0, `make verify-grok-plugin` exit 0, `npm run check --prefix adapters` exit 0, `make lint-config` exit 0, README 설치 명령 순서 실행(임시 HOME), marketplace 경로 설치·삭제 후 동작·uninstall, 이전 구성 재현, 음성 시험 2종. DEV 깨끗한 `5506d64`에서 `lint.py --from 3769814` exit 0.
- 실행하지 않은 검증: Go 테스트와 `make verify-mvp`(코어 불변, 이전 증거 재사용), Node 20 `npm ci`(Node 20 없음, DEV 로그만 확인), 실제 Bot 계정 설치·Bot 앱 Marketplace·`knowslink_status` 세션 호출·hosted Node(모두 held 또는 유료 inference), `GROK_CONFIG*` 누출 시험.
- 재시험 댓글 초안의 명령은 README와 같은 형태이며 Step 2–3의 기대 출력이 실제와 일치했다. Step 0·1은 Bot 컴퓨터에서만 확인할 수 있다. F-05와 함께 실패 시 `uname -m`·`node --version` 회신을 요청하면 분류가 쉬워진다.
- 이 리뷰는 TESTER의 독립 동작 QA와 원격 재시험을 대체하지 않는다.

## 검토 결론

수락 가능. 미해결 critical/high/medium이 없고 low 5건은 후속 개선이다. 설치 수정의 핵심 주장(이전 ZIP은 MCP 0개, 수정 ZIP은 `grok plugin install`로 1개 서버와 2개 도구를 노출, ZIP 해시 재현, 격리·실패 전파)을 독립 환경에서 재현했다. 실제 Bot 앱 노출은 증명되지 않았고 문서가 그렇게 주장하지도 않는다. check 통과는 AI 검토의 내용·테스트 성공을 자동 보증하지 않는다.
