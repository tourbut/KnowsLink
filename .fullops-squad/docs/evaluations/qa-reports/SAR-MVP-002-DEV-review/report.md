# SAR-MVP-002-DEV 리뷰

- 검토자 / CLI / 모델: ops 독립 검토자 / Claude Code / claude-sonnet-5-5. 세션 `0f2048f7-ec36-47bb-8707-f202b5c14e02`.
- 구현자 세션: Codex `01a101f4-2c66-7843-a503-b808214ee39f`. rollout의 첫 session_meta cwd는 `/home/shin/orca/workspaces/KnowsLink/fullops-dev`다. 두 세션은 다르다.
- snapshot: `/tmp/knowslink-plugin-review-552586b`. 고정 head `552586b6e886f95bffa9a000a031ea03070afedb`, detached, 작업 트리 clean. 읽기 전용으로 유지했고 리뷰 중 변경하지 않았다.
- base SHA / head SHA / merge-base: `dbdd70086971285b790683f362702e5a9ff55acd` / `552586b6e886f95bffa9a000a031ea03070afedb` / `dbdd70086971285b790683f362702e5a9ff55acd`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c). `.fullops-squad/review/rule.json` sha256 `ab2116fb…2aa0`. FullOps 0.9.13 `review.py`. 공통 규칙 fullops-common-0.3.2 README와 coding-style·testing·security, FULLOPS.md, project.md, 문서 작성 규칙을 적용했다. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, `docs/exec-plans/phases/SAR-MVP-002-DEV.md`, `adapters/README.md`와 공식 Cursor·Grok Bot 문서.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 36 / 28 / 8 / 28 / 8. 제외 8개는 DEV 검증 로그(`*.txt`)이며 skipped 사유를 파일별로 기록했다.
- lint(`lint.json`): ERROR 0 / WARNING 3 / 실행 불가 0. coor가 552586b 깨끗한 HEAD에서 실행했다. `product-lint: make lint` passed, exit 0.

## 검토 범위

- 28개 reviewed. 제품 소스 `core.ts`, `mcp.ts`, `index.ts`, `synthetic.ts`, `mcp.test.ts`, `package_plugin.py`, Makefile, manifest, `mcp.json`, SKILL.md, README, 설계 문서 5건과 실행·인계 기록을 읽었다. 그 소스와 `internal/relay/http.go`의 claim·gate-consume·authorize 서버 검사를 대조했다.
- 8개 skipped는 OCR이 확장자로 제외한 DEV 원시 로그다. 결과 줄·exit·해시와 `verification.json` 연결로 존재를 확인했다. 생략 영향은 없다.
- `package-lock.json`은 생성물이라 줄 단위 검토를 생략했다. `npm ls --all`(invalid·extraneous·missing 0), SDK 1.32.0·esbuild 0.28.2 고정, 96개 prod 패키지의 LICENSE 확보로 대신했다.
- 기존 검증 증거를 재사용했다. `c4ebbec` 이후 제품 소스(core·mcp·test·synthetic)는 불변이다. `5bd1bd2..552586b`에는 제품 변경이 없고 문서·기록만 있다. `make lint/test/verify-mvp`(c4ebbec)와 `make plugin`(5bd1bd2)의 exit 0이 그대로 유효하다. Go/UI는 불변이며 Chrome QA 9584aaf는 원래 SHA의 근거로 재사용했다. 전체 검사는 반복하지 않았다.

### 영역별 결과

| 영역 | 결과 |
|---|---|
| plugin manifest | `name` kebab-case, `skills`·`mcpServers` 경로, `mcpServers` 래퍼, `${CURSOR_PLUGIN_ROOT}`가 공식 plugin reference와 일치한다. marketplace `source: "knowslink"` 상대 경로도 일치한다. |
| 공식 설치 근거 | Dashboard → Plugins & MCPs → Team Marketplaces → Import from Repo, Default Off, 로컬 `~/.cursor/plugins/local`을 공식 문서에서 다시 확인했다. Grok Bot connect 페이지는 stdio MCP·ZIP 업로드·Node runtime을 언급하지 않는다. README는 이를 주장하지 않고 미확인으로 적었다. 일치한다. |
| ZIP 재현성 | 고정 timestamp·권한·정렬된 notice다. snapshot을 `git archive`로 임시 복사해 build + `package_plugin.py --verify`를 실행했다. ZIP은 8개 파일이고 압축 해제 bundle의 MCP 검사가 exit 0이다. 내 환경이 DEV의 `node_modules`를 symlink로 사용해 bundle 주석 경로가 달랐다. 경로 정규화 후 `dist/plugin.js`는 DEV ZIP과 바이트 동일했다. 나머지 7개 파일도 동일하다. 기록된 ZIP SHA `0e671d1a…9117`과 내 ZIP SHA는 이 경로 차이로 다르다. SHA 자체를 재현한 것은 아니며, 내용 동일성으로 확인했다. 배포된 bundle에 호스트 경로는 0건이다. |
| MCP stdio·stdout | `console.log`·`process.stdout` 사용이 없다. gate 로그 콜백은 mcp 경로에서 빈 함수다. 도구 출력은 5필드뿐이다. 테스트가 stderr 빈 값을 단언한다. stdout 오염은 SDK 클라이언트 파싱 실패로 이어진다. |
| secret·입력·URL | 도구 입력은 비어 있다. 설정은 서버 환경에서만 읽고 `held` 기본·`production` 값 held다. HTTP loopback root만 허용하고 userinfo·path·query·hash를 거부한다. `redirect:"error"`다. 테스트 5종이 요청 0건까지 단언한다. F-03 참조. |
| persist/ACK/claim/gate/result | 서명 검증 → persist → ACK → claim 순서이고, 각 실패 시 후속 호출이 없음을 테스트한다. `relay.result`는 ACK 후 종료한다. 서버 쪽 claim은 중복 시 `already_claimed`다. gate 대기는 SQL e2e가 실제 bundle로 검증한다. F-01·F-02 참조. |
| 회귀 영향 | `index.ts` CLI 계약이 유지된다. core 이동은 기능 변경 없는 추출이다. 신규 dependency는 SDK 1.32.0 하나이며 pin이다. `make test`에 npm test가 추가됐다. |
| 문서-동작 | README·SKILL·설계 문서의 상태 값·제한·재개 조건이 코드와 일치한다. 불일치를 찾지 못했다. |

## 발견 사항

| ID | 심각도 | 파일·줄 | 설명 | 상태 |
|---|---|---|---|---|
| F-01 | medium | `adapters/src/mcp.ts:268-283` | `knowslink_pull_once`가 gate 결정이나 exp까지 한 호출에서 대기한다. 호스트 tool timeout이 더 짧으면 클라이언트는 timeout을 받고 서버 작업과 `busy`가 남는다. 취소 신호는 처리하지 않는다. claim은 재발급되지 않아 fail-safe다. README가 busy·재시도 금지를 설명한다. 실제 Bot 앱의 timeout은 미확인(held)이다. | 미해결. 실제 연결 재개 전 앱 timeout 확인과 필요 시 비차단 상태 조회 설계. |
| F-02 | low | `adapters/src/core.ts:170-187` | gate가 거절·만료·대기 만료여도 `/v1/authorize`로 진행한다. relay authorize는 claim만 검사한다. 현재 authorize는 `executable:false`/`disclosure:false` 고정이고 결과가 denied뿐이라 효과가 없다. 기존 동작을 이동한 것이다. | 미해결. disclosure·calendar effect 추가 과제는 gate-consume 성공을 authorize 전제로 강제해야 한다. |
| F-03 | low | `adapters/src/core.ts:204-213` | host 허용 목록에 DNS 이름 `localhost`가 있다. 비정상 hosts 구성이면 Bearer가 loopback 밖으로 갈 수 있다. 합성 모드와 시험 키 한정이며 기본 held다. | 미해결. 숫자 loopback만 허용하면 더 단단하다. |

critical 0, high 0, medium 1, low 2. 재현: F-01은 합성 relay에서 gate를 pending으로 두고 호스트 timeout보다 길게 대기시키는 구성이 필요하다. 이번 리뷰는 코드 경로와 서버 검사의 정적 대조로 확인했고 앱 timeout은 실행하지 않았다.

lint WARNING 3건의 의미는 다음과 같다.

- `PLANS.md` SIZE-001(511줄)은 coor 소유 파일이 준비 커밋에서 커진 것이며 DEV 변경이 아니다.
- `mcp.test.ts:66` SEC-001은 더미 문자열 `"synthetic-lease"`다. 실제 credential이 아니다.
- `SAR-MVP-002-DEV.md:85` SEC-001은 같은 더미 값을 설명하는 문장이 일으킨 오탐이다.
- 실제 비밀값은 없다. 소스·manifest·ZIP·로그에서 credential 형태 값을 찾지 못했다.

## 검증 및 남은 제약

- 직접 실행: snapshot 상태 확인(HEAD·detached·clean), 독립 build·package·압축 해제 MCP 검사(exit 0), `npm ls --all`, 공식 문서 3건 재조회, git diff 범위 확인. 위 실행은 snapshot 밖 임시 사본에서 했고 snapshot은 수정하지 않았다.
- 재사용: DEV의 `make lint/test/verify-mvp/plugin` exit 0과 SQL e2e PASS(`mcp-sql-mvp.txt`), coor lint.json(ERROR 0), Chrome QA 9584aaf.
- 실행하지 않은 것: 실제 Grok Bot 계정·앱 빌드·marketplace 등록·hosted Node·tool timeout·외부 연결·실데이터·유료 API. 모두 held다. 이 리뷰는 TESTER 독립 동작 QA와 직접 시각 검수를 대체하지 않는다.
- skipped 8건(로그)은 수락 판단에 영향이 없다.
- 남은 일: F-01은 실제 연결 재개 전, F-02는 disclosure·effect 과제 전, F-03은 다음 어댑터 수정 때 처리한다. PLANS.md는 coor가 갱신한다.

## 검토 결론

수락 가능하다. 고정 후보 `552586b6e886f95bffa9a000a031ea03070afedb`는 코드·문서·공식 설치 근거 리뷰에서 미해결 critical/high가 없다. medium 1·low 2는 실제 연결 재개 전 후속으로 남긴다. 실제 연결과 Bot 앱 설치 지원은 증명되지 않았고 held가 유지된다. check 통과는 기록 검사이며 AI 판단 결과다.
