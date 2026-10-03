---
title: SAR-MVP-002-DEV 실행 기록
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-002-DEV]
summary: 공식 제품 후보와 안전 pull 준비 및 검증과 재개 조건을 기록한다
---

# SAR-MVP-002-DEV — 공식 인터페이스 조사와 안전 pull 준비

## 기준과 승인 범위

기준 ref는 `dbdd70086971285b790683f362702e5a9ff55acd`다. 착수 HEAD는 `0f58336f5ec6c678dfdf416e1cc8fee3e1ef7472`이며 브랜치는 `fullops/dev`다. `fullops-common-0.3.2`, FULLOPS.md, project.md, D02·D03·D05와 frozen 원천을 직접 읽었다. Context7 조회에는 비공개 코드·비밀값을 보내지 않았다.

공식 읽기 전용 조사·로컬 합성 검증·DEV 기술 문서·커밋·역할 브랜치 일반 push는 승인됐다. 실제 Grok 계정 연결·업무 발송·실데이터·유료 API·운영 배포·공유 서비스 변경은 held다. 독립 리뷰와 TESTER QA는 coor가 고정 후보에서 배정한다.

## 짧은 기술 계획

1. 원천의 이름과 공식 제품 정체성을 구분한다. 현재 원천은 URL·배포판을 특정하지 않는다.
2. 공식 shell·MCP·headless 경로와 인증·owner 권한을 확인한다. 지원되지 않은 inbound API를 만들지 않는다.
3. 기존 TypeScript Adapter의 서명·shared persist·ACK·claim·gate·최소 denied result 경계를 재사용한다. 대상 확정 뒤 공식 Cursor plugin manifest·stdio MCP·skill을 준비한다.
4. 로컬 합성 회귀와 해당 변경 검사를 실행한다. 문서와 held의 담당·재개 조건을 보존한다.

## 공식 근거와 판정 — 확인일 2026-10-03

Context7 `resolve-library-id`의 Grok Bot 후보는 `/websites/x_ai_grok-bot`다. 같은 ID의 `query-docs`로 persistent cloud computer·shell·MCP policy·member approval 근거를 조회했다. 버전 지정 없는 live 문서이며 특정 앱 빌드를 검증한 근거가 아니다.

| 후보·공식 출처 | 확인 내용 | 현재 판정 |
|---|---|---|
| [Grok Bot overview](https://docs.x.ai/grok-bot/overview), [computer](https://docs.x.ai/grok-bot/computer-and-apps) | cloud computer에 shell·파일·브라우저가 있다. `/workspace`와 credential은 동일 사용자의 여러 Bot이 공유한다 | shell에서 기존 pull 프로그램을 실행하는 준비는 기술적으로 가능하다. 사용자 대상 확정·실제 컴퓨터 runtime·network 권한은 미확인이다 |
| [Bot routines](https://docs.x.ai/grok-bot/skills-routines-and-automations), [chat](https://docs.x.ai/grok-bot/chat-and-collaboration) | 대화·저장 skill·schedule routine과 일부 Slack/GitHub event trigger를 제공한다 | 이 문서에서 arbitrary relay inbound endpoint나 typed Bot task API는 확인되지 않았다. 부재 전체를 증명하지는 않는다 |
| [Bot security](https://docs.x.ai/grok-bot/security), [get started](https://docs.x.ai/grok-bot/get-started) | Cursor 계정과 앱이 필요하다. Bot은 member의 권한으로 작동한다. Auto Review·network policy와 local execution 승인은 별도다 | vendor 승인은 KnowsLink owner accept/approve를 대체하지 않는다. shell 지원만으로 안전 통합 수락을 주장하지 않는다 |
| [Grok Build](https://docs.x.ai/build/overview), [headless/ACP](https://docs.x.ai/build/cli/headless-scripting), [MCP](https://docs.x.ai/build/features/mcp-servers) | coding CLI는 headless와 stdio ACP·MCP를 제공한다. 로그인 또는 inference key는 별도다 | Grok Bot과 다른 후보다. 사용자 대상 변경 없이 CLI를 002 대상으로 선택하지 않는다 |
| [xAI inference reference](https://docs.x.ai/developers/rest-api-reference/inference) | 별도 inference API 제품이다 | Bot 메시지 수신 API로 대체하지 않는다. 유료 호출은 실행하지 않았다 |

공식 Bot overview·get-started·chat의 갱신 표시는 2026-09-21이다. computer·routines는 2026-09-14이며 security는 2026-09-28이다. 문서 조회일과 앱 설치 버전을 구분한다. 현재 Grok Bot 앱 빌드·hosted Node/Docker 설치 상태는 미확인이다.

## 로컬 CLI 관찰

`command -v grok`는 `/home/shin/.local/bin/grok`다. symlink의 실제 파일은 `/home/shin/.grok/downloads/grok-1.0.46-linux-x86_64`다. `grok --version`은 `grok 1.0.46 (2765805b9442) [stable]`을 출력했다. `grok --help`의 제목은 `Grok Build TUI`이며 `-p`, `--output-format`, `--permission-mode`, `agent`를 제공한다. 두 명령의 종료코드는 0이다. 모델 prompt·로그인·MCP 등록·전역 설정 변경은 실행하지 않았다. 공식 live 문서의 옵션을 설치 1.0.46 지원이라고 자동 간주하지 않는다.

## 제품 정체성 질문

Orca ask에서 공식 docs.x.ai/grok-bot, 설치 CLI, 다른 대상 URL을 구분해 질문했다. coor는 사용자에게 확인을 요청했고, 답 전에는 후보별 조사·기존 합성 경계 검증·문서만 진행하라고 회신했다. 이후 coor 메시지 `msg_1d99c6fd5e89`에서 사용자가 공식 Grok Bot과 플러그인 제작을 확정했다. `msg_d8b63570d547`은 같은 DEV 과제의 구현 책임을 확인했다. 정체성 held는 해제됐고 실제 계정/연결 held는 유지한다.

## 검증 기록

준비 HEAD `0f58336f5ec6c678dfdf416e1cc8fee3e1ef7472`에서 기존 `make verify-mvp`는 exit 0이었다. 코드 수정 후 후보 SHA에서 새 MCP 경로와 SQL 회귀를 별도로 실행한다. UI 코드는 변경하지 않았다. 기존 Chrome QA `9584aaf`는 재사용 근거이며 이번 SHA에서 새 직접 시각 검수를 실행한 증거가 아니다.

## 후속과 재개 조건

- coor/designer: 공식 Grok Bot 대상과 플러그인 제작은 사용자에게 확인했다. 제품 우선순위와 정책 변경은 기존 책임을 유지한다.
- DEV: 공식 plugin/MCP package와 설치 안내를 준비했다. 실제 앱 빌드·MCP discovery·hosted Node와 승인된 network는 실제 연결 재개 때 대조한다.
- owner/OPS: 실제 연결 승인, 계정 권한, 안전 secret 전달, 격리 relay 도달 경로가 필요하다. hosted loopback은 이 개발 호스트의 loopback과 다르다. 현재 보호된 베타·Tunnel·Tailscale을 재설정하지 않는다.
- coor: 고정 SHA 독립 리뷰·TESTER QA와 PLANS/board 갱신을 담당한다. actual connection held와 DEC-02의 무정책 deny·calendar stub·외부 exactly-once 비주장을 유지한다.

## 확정된 플러그인 기술 계획과 구현

Cursor 공식 `.cursor-plugin/plugin.json`, `mcp.json`, skill 구조를 사용한다. SDK 1.32.0이 stdio MCP initialize·discovery·tool 호출을 처리한다. esbuild 0.28.2로 고정 SDK·core를 하나의 ESM bundle에 포함한다. package 생성은 허용 목록과 고정 ZIP timestamp를 사용한다. 설치 대상에는 Node 22.22.2만 필요하며 신규 runtime 의존성 다운로드는 없다.

기본 mode는 held다. `synthetic-loopback`만 기존 Adapter를 사용하며 원문·claim·credential을 모델에 반환하지 않는다. 입력 없는 pull-once 도구는 서명·persist·ACK·claim·gate·최소 denied R을 순서대로 수행한다. gate 승인 권한과 무정책 deny를 유지한다. 외부 URL·redirect·production mode는 차단한다.

초기 bundle 테스트 exit 1은 index CLI의 direct-run guard가 bundle 경로와 일치해 stdout에 합성 CLI 출력을 기록한 결함이었다. 공통 core와 index CLI를 분리했다. 같은 bundle handshake 테스트가 exit 0으로 바뀌었다. 이 실패를 환경 문제나 PASS로 바꾸지 않는다.

Cursor plugin reference의 package 형식과 설치 경로는 [공식 문서](https://cursor.com/docs/reference/plugins)로 확인했다. Bot의 connector policy와 account-wide plugin 지원은 [Bot 공식 문서](https://docs.x.ai/grok-bot/teams-and-enterprises)를 따른다. Cursor team marketplace의 repository import와 로컬 개발 경로는 [공식 설치 문서](https://cursor.com/docs/plugins)다. 실제 Grok Bot 앱의 ZIP 직접 업로드나 Cursor IDE 로컬 import의 Bot 자동 전달은 확인되지 않았다. package 구조 검증과 실제 계정 설치/도구 검색은 구분한다.

[사용자 설치 안내](../../../../adapters/README.md)에 artifact 구조, marketplace를 통한 owner/admin 설치, hosted runtime, synthetic-only 설정, 상태 의미와 actual connection held를 기록했다. 실제 marketplace 등록·계정 로그인·MCP 활성화·유료 호출·운영 연결은 실행하지 않았다.

## 고정 SHA·검사별 종료코드·산출물

제품 구현 후보는 `c4ebbecd3ef91be10ecbb517fe451e02769358bc`다. 패키징 검증과 제3자 라이선스 고지 후보는 `5bd1bd2f4ea32af625831a284dadd1551821530e`다. 후자에서 core/MCP/test/synthetic TS는 전자와 동일하다. 완료 기록 보존 커밋은 문서·로그·빈 인박스만 포함하며 최종 전체 SHA는 worker_done에 명시한다.

| 명령 | 실행 대상 | 자기 종료코드 | 증거 |
|---|---|---|---|
| `make verify-mvp` (기존 경계) | 준비 `0f58336f5ec6c678dfdf416e1cc8fee3e1ef7472` | 0 | [baseline](../logs/SAR-MVP-002-DEV/baseline-verify-mvp.txt) |
| `make lint` | `c4ebbecd3ef91be10ecbb517fe451e02769358bc` | 0 | [제품 검사](../logs/SAR-MVP-002-DEV/product-lint.txt) |
| `make test` | 같은 구현 후보 | 0 | [Go race·SDK MCP](../logs/SAR-MVP-002-DEV/product-test.txt) |
| `make verify-mvp` (MCP bundle 연동) | 같은 구현 후보 | 0 | [실제 SQL·gate·result](../logs/SAR-MVP-002-DEV/mcp-sql-mvp.txt) |
| `make plugin` | `5bd1bd2f4ea32af625831a284dadd1551821530e` | 0 | [bundle·압축 해제 실행](../logs/SAR-MVP-002-DEV/package.txt) |
| `lint.py --repo . --from dbdd70086971285b790683f362702e5a9ff55acd` (FullOps 0.9.13) | 같은 패키징 후보 | 0 | [고정 lint JSON](../logs/SAR-MVP-002-DEV/lint-candidate.json) |

모든 검사 명령은 결과를 파일에 저장하고 `$?`를 즉시 보존한 뒤 그 코드로 종료했다. 완료된 로그의 tail 조회는 검사의 통과 근거로 사용하지 않았다. `make plugin`은 `make build`와 압축 해제한 standalone bundle의 실제 SDK initialize/discovery/tool 호출을 포함한다. 기존 Go race·SQL 회귀는 frozen field·current-auth·owner 분리·human routing·claim·gate 경계를 검증했다. MCP 경계 검사는 서명·persist·ACK·claim 실패 뒤 판단/결과가 실행되지 않음을 확인했다.

FullOps candidate 결과는 ERROR 0, WARNING 2, 실행 불가 0이다. 등록 `product-lint: make lint`의 자기 exit도 0이다. SIZE-001은 준비 커밋이 이미 변경한 coor 소유 PLANS.md의 511줄 경고다. DEV는 PLANS를 변경하지 않았다. SEC-001은 테스트의 `lease_token: "synthetic-lease"` 더미이며 실제 credential이 아니다. 경고를 숨기거나 lint 설정을 완화하지 않았다.

archive는 [build/knowslink-grok-bot-plugin.zip](../../../../build/knowslink-grok-bot-plugin.zip)이다. SHA256은 `0e671d1a89c141d896034fff31619b9cd2148b73b567adbc3a97126031989117`이다. 소스와 생성 명령은 Git에 보존하고 archive는 ignored build artifact로 유지한다. 같은 소스·고정 dependency에서 `make plugin`으로 재생성한다. package에는 runtime dependency의 LICENSE 전문을 THIRD_PARTY_NOTICES.txt로 포함한다. package SHA와 명령/HEAD/exit 연결은 [verification.json](../logs/SAR-MVP-002-DEV/verification.json)에 있다.

공식 URL·Context7 library ID·버전 한계·제품 확정 메시지는 [official-sources.json](../logs/SAR-MVP-002-DEV/official-sources.json)에 보존했다. `npm install` 두 번의 종료코드는 0이며 각각 audit vulnerabilities 0이었다. `quick_validate.py adapters/skills/knowslink`도 exit 0이었다.

D03(architecture/tech-stack), D05(interface), D10(module), D13(transition)의 실제 영향만 갱신했다. contexts/dev.md와 사용자 package README를 갱신했다. 기술 문서 front matter는 deliverables.py stamp로 썼다. 완료 기록을 갱신한 작업 트리에서 `deliverables.py --repo . --strict`는 exit 0이었다. [strict 로그](../logs/SAR-MVP-002-DEV/deliverables-strict.txt)에 보존했다.

기존 Chrome QA `9584aaf`는 UI·Go 코드 변경 없음의 재사용 근거다. 새 Bot 앱 UI 검수·실제 plugin install·marketplace indexing·hosted stdio startup은 미실행이다. coor가 고정 후보의 독립 리뷰와 TESTER QA를 배정한다. 미해결 critical/high는 계속 수락 차단이다. 이번 worker 성공은 공식 plugin package 준비와 합성 로컬 검증 완료이며 actual connection 성공이나 전체 MVP 수락이 아니다.

완료 로그를 staging한 첫 `git diff --cached --check`는 Compose progress의 줄 끝 공백 때문에 exit 2였다. 커밋이 만들어지지 않은 작업 트리에서 뒤따른 lint는 깨끗한 HEAD 조건으로 exit 2였다. 두 결과는 제품 검사 실패나 PASS가 아니다. 제품 로그의 줄 끝 공백만 정리했고 원본/정리본 SHA256을 log-normalization.json에 보존했다. 원본 `/tmp/SAR-MVP-002-DEV-*.log`는 그대로 남겼다. 내용·명령·종료코드·테스트 수는 바꾸지 않았다. 최종 커밋의 깨끗한 HEAD에서 lint를 다시 실행한다.
