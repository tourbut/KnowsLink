---
title: Grok Bot webhook routine wake 후보 실행 기록
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-GROK-WEBHOOK-001-DEV]
summary: 사용자 메시지로 보이지 않는 webhook routine 깨우기 후보의 기술 계획·구현·검증·미검증 범위를 기록한다
---

# SAR-GROK-WEBHOOK-001-DEV — webhook routine wake 후보

## 기준

- 지시서: `handovers/to_dev.md`, attempt `590f4c27d88d410c8c15fe232581b35a`, 기준 SHA `7a8c2b9194edb77b52db41a14b2325af743f0d6e`.
- 규칙: `fullops-common-0.3.3`, `project.md`, `contexts/dev.md`, `docs/agents/document-writing.md`. 테스트 레벨 lite. 하위 에이전트는 사용하지 않았다.
- 공식 근거: [cursor/plugins pstack make-bot-ui SKILL.md](https://github.com/cursor/plugins/blob/main/pstack/skills/make-bot-ui/SKILL.md)(2026-10-11 WebFetch). URL `https://api2.cursor.sh/automations/webhook/<id>`(query 없음, ID 추측 금지), POST, `Content-Type: application/json`, `Authorization: Bearer <key>`, `X-Automation-Key: <key>`, 8초 단일 시도·재시도 없음, 200은 routine wake, routine 턴은 `<webhook_event>`의 `body`를 untrusted로 읽음, 무해한 probe 한 번 뒤 live 주장. Context7은 직전 과제에서 월간 quota 초과 기록이 있어 공식 문서로 확인했다.
- Bot 보고 근거: `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-DOTS-GROK-LIVE/webhook-wake-research.json`(UpdateRoutine `trigger.type=webhook`, `webhook_event`).

## 기술 계획

1. 새 class를 만들지 않는다. 기존 `GrokWake`에 webhook 모드를 추가한다. 호출자 `inbox.ts`·`mcp.ts`·`text.ts`는 `GrokWake.fromEnv`만 사용하므로 바꾸지 않는다.
2. `KNOWSLINK_GROK_WEBHOOK_FILE`(절대경로)이 있으면 그 process의 wake는 webhook만 사용한다. `sendPrompt`와 동시에 울리지 않고 실패 뒤 fallback도 하지 않는다. 변수가 없으면 기존 loopback/off 동작이 그대로다. 기존 `sendPrompt` 코드·설정은 제거하지 않는다.
3. URL·key는 비공개 파일(`privatePath`: symlink 거부, 0600·소유자 또는 Windows ACL)에서 요청마다 읽는다. URL은 고정 host 정규식만 허용한다. key는 출력 가능한 ASCII 8–4096자다.
4. 요청은 기존 `fetch`·`AbortSignal.timeout`·`redirect: "manual"`을 공통 `request`로 묶어 재사용한다. 결과 분류(retry/rejected/uncertain)와 ID별 `.wake` marker는 기존 계약을 따른다.
5. body는 `{event:"knowslink_text_arrived", ids, pending, next}`이다. ids는 UUID 검증 후 최대 10개다. 받은 text·발신자는 넣지 않는다.
6. `wake-check`는 webhook 모드에서 파일만 검증한다. POST는 routine 턴을 시작하므로 probe는 owner·coor의 수동 단계로 남긴다.

설치 스크립트 `scripts/install_bot_mcp.sh`는 bundle만 배치하고 env를 정하지 않는다. watcher env는 owner가 시작 명령에 붙인다. 그래서 변경하지 않았다.

## 변경

- `adapters/src/grok-wake.ts`: webhook 모드, 공통 `request`, `ENOTFOUND`도 확실한 미전달(`retry`)로 분류, `check()` 파일 검사.
- `adapters/src/grok-wake.test.ts`, `adapters/package.json` test script: 단위 검사 추가.
- `adapters/README.md`: 지원 상태 표 정정, webhook 절 추가.
- D03 `architecture.md`: "webhook 공식 기능 없음" 일반화를 확인된 후보 범위로 정정하고 webhook 절을 추가했다. 지시서의 일반화 정정 요구에 따른 것이다.
- D10 `module-design.md`: 책임·검증 표 행 추가. D12 `ops-guide.md`: 전환·복귀 절차 추가.

보안 판단:
- SSRF: 사용자 입력 URL이 아니라 owner의 비공개 파일이다. 그래도 host·scheme·path 형식을 고정한다. userinfo·port·query·fragment·`..`는 정규식에서 거부된다. redirect는 따르지 않으므로 secret이 다른 host로 가지 않는다.
- 비밀값: key·URL은 결과·상태·marker·로그에 넣지 않는다. 응답 본문은 읽지 않고 cancel한다.
- 프롬프트 삽입: routine에는 UUID·수·고정 안내만 간다. routine prompt는 owner가 쓴다.
- 제품 원천 C5 "optional webhook OFF": relay가 vendor로 webhook을 보내는 기능은 계속 없다. 이 변경은 owner Bot 컴퓨터의 adapter가 owner 자신의 routine에 보내는 opt-in·기본 off wake다. 제품 규칙 변경으로 보지 않았다. coor·designer가 다르게 판단하면 후속에서 조정한다.

## 검증

| 검사 | 위치·대상 | 결과 |
|---|---|---|
| `npm run check`(prettier·eslint·tsc) | Windows 작업 트리 | 종료코드 0 |
| `npm test`(build + 7개 test) | Windows Node 24, 작업 트리 | 종료코드 0. 기존 inbox loopback wake 회귀 포함 PASS |
| `node dist/grok-wake.test.js` 변이 | URL host 완화(`cursor\.sh.*`), `redirect: "follow"`, ENOTFOUND 분류 제거 | 세 변이 모두 AssertionError. 원본 복구 뒤 PASS |
| 등록 product-lint/product-test | Linux 고유 clean clone, 고정 SHA | 아래 Linux gate 절 |

단위 검사 범위: 미설정 null, 상대경로 invalid, webhook+agent 동시 설정 시 webhook만, 2xx accepted_unverified, header·body key 집합·비UUID 제거, 302·401 rejected, 503 uncertain, ENOTFOUND retry, 50ms timeout uncertain 단일 시도, http·다른 host·userinfo·query·`..`·짧은 key·공백 key·key 없음 무전송, Linux 0644 파일 무전송, `check()` 무전송, 결과 JSON에 key·URL ID 없음.

### Linux gate

기존 Linux 서버에서 고유 `/tmp/knowslink-webhook-f39b677-ctxbcbfa2f4992f-2/repo` clean clone을 사용했다. 이 과제 commit의 git bundle에서 `f39b6777abe5168ed0eba0f5fab138f5a64f5ad2`를 detach checkout했다. 예전 product-byte shim은 쓰지 않았다. Node v22.22.2, npm 10.9.7, go1.22.2. `npm ci` 뒤 FullOps 1.3.0 `lint.py --from 7a8c2b9`를 실행했다. DB·container·운영 서비스는 사용하지 않았다.

| 항목 | 결과 |
|---|---|
| product-lint `make lint` | passed |
| product-test `make test`(Go race + `npm test --prefix adapters` 7개 + shell guard) | passed |
| lint.py | 종료코드 0, ERROR 0, WARNING 3, 실행 불가 0. 실행 전후 clone clean |

WARNING 처리:
- SIZE-001 `PLANS.md` 1538줄: coor 소유 파일이며 이 과제는 수정하지 않았다. 기준 이후 coor commit의 diff다.
- DEP-001 `adapters/package.json`: test script에 `grok-wake.test.js`만 추가했다. 의존성 변경은 없다.
- SLOP-001 `grok-wake.ts:163`: 주석의 영어 단어 "any"를 TypeScript `any`로 오인한 경고다. 주석 문구를 바꿨다.

첫 실행(`…-ctxbcbfa2f4992f`)은 lint 도구 tar에 `assets/`를 빠뜨려 lint.py가 시작 전에 실패했다(제품 검사 미실행). 같은 SHA로 새 고유 경로에서 다시 실행했다. 두 임시 경로는 삭제하지 않고 보존했다. archive 뒤 최종 HEAD의 gate 결과는 worker_done 보고에 적는다.

## 미검증과 후속

- 실제 routine 생성·key 발급·probe·routine 턴 표시·노우 receive는 실행하지 않았다. coor·owner 후속이다. 그 전에는 현재 `sendPrompt` watcher 설정을 유지한다.
- api2.cursor.sh의 2xx 외 응답 코드 의미는 공식 문서에 없다. 4xx를 rejected, 5xx를 uncertain으로 둔 것은 기존 계약을 따른 판단이다.
- 공식 문서는 POST 실패 시 로컬 로그에 남기고 routine에서 drain하라고 권한다. 이 구현은 이미 local inbox가 그 역할을 한다. routine은 `knowslink_text_receive`로 inbox를 비운다.
- 독립 tester QA·코드 리뷰는 고정 SHA snapshot에서 후속으로 수행한다. 제품 배포·Bot 설정 변경은 하지 않았다.
