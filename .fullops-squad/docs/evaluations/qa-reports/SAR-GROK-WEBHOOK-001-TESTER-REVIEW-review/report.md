---
title: SAR-GROK-WEBHOOK-001-TESTER-REVIEW 독립 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-GROK-WEBHOOK-001-TESTER]
summary: webhook routine wake 고정 후보 b86f397의 독립 lite QA와 코드 리뷰 결과를 기록한다
---

# SAR-GROK-WEBHOOK-001-TESTER-REVIEW 리뷰

- 검토자 / CLI / 모델: tester, Claude Code, Claude Sonnet 5.5. 실제 세션 `254e09d8-cebb-4a7d-a717-ff1820a00fe6`. 구현자 세션 `c7162b6d-560e-40a5-9717-515db79c0d0a`와 다르다.
- base SHA / head SHA / merge-base: `f0fdf9daf24a7def3b5ec317692ed8bcb720ec12` / `b86f39711c2f5240d37332f44a9930d4dc2653fa` / `0de25a5a20044e4a30e72bca6b65b8202ec314f4`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-667bd7fbdeea48ac867606b5efcec14b` (`review.py snapshot`, detached, head 일치, 수정하지 않음). DEV의 최신 UI 커밋은 포함하지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13 delegate(LLM 호출 없음). `rules.json` 그룹 3개(md, json, ts)를 한 번씩 읽고 적용했다. 공통 규칙 `fullops-common-0.3.3`(`rules/common/README.md`, `testing.md`, `security.md`), `project.md`, `FULLOPS.md`를 이 checkout(`bd24be0`)에서 읽었다. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_tester.md`, `docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md`, 공식 [make-bot-ui SKILL.md](https://github.com/cursor/plugins/blob/main/pstack/skills/make-bot-ui/SKILL.md)(2026-10-11 WebFetch), D03/D10/D12, 제품 원천 C5(`docs/planning/sources/silent-agent-relay/architecture.md` 37줄·100줄).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 12 / 12 / 0 / 12 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: head `b86f39711c2f`, 파일 8, ERROR 0, WARNING 1(DEP-001), 실행 불가 0. 실제 명령 종료코드 0(아래).
- SIZE-002: 발생하지 않았다. 해당 없음.
- DEP-001: `adapters/package.json`은 test script에 `node dist/grok-wake.test.js`를 붙인 1줄 변경이다. `dependencies`·`devDependencies` 변경 없음(diff 확인). 새 의존성은 없다. fetch·AbortSignal은 Node 표준이다.
- UI 디자인: 해당 없음. UI 파일 변경 없음(이 후보는 adapter와 문서뿐).

## 검토 범위

12개 파일 전부 `reviewed`다. 이유는 `result.json`의 파일별 `reason`에 있다. 제품 코드는 `adapters/src/grok-wake.ts`(diff 전체와 호출자 `inbox.ts` ring·marker, `mcp.ts` wake-check, `index.ts`)와 `grok-wake.test.ts`를 줄 단위로 읽었다. 문서는 README webhook 절, D03·D10·D12 갱신, phase 문서를 구현·공식 규격과 대조했다. `to_dev.md` 삭제는 인박스 비움이다. 로그 전문의 의미 검토는 하지 않았다. 비밀값 패턴만 diff에서 검색했고 발견되지 않았다.

## 발견 사항

미해결 critical/high 없음.

**F1 — low — `adapters/src/grok-wake.ts:108-111`**
- 조건: Bot 컴퓨터의 DNS가 일시 실패(`EAI_AGAIN`)하거나 네트워크 경로가 없다(`ENETUNREACH`/`EHOSTUNREACH`).
- 재현: 프로브에서 `[2001:db8::1]`로 fetch하면 4ms 만에 실패하고 `uncertain/webhook_no_response`가 된다. `ECONNREFUSED`·`ENOTFOUND`만 `retry`다.
- 영향: 확실히 전달되지 않았는데 marker가 남아 `.wake` 만료(24h)까지 같은 ID를 다시 울리지 않는다. 메시지는 inbox에 남아 수동 receive나 다음 새 ID의 깨우기에서 읽힌다. 중복 전송·비밀 노출은 없다. 안전한 방향의 실패다. README의 "DNS 실패는 retry"는 `ENOTFOUND`에만 맞다.
- 수정 상태: 미수정(tester는 제품 코드를 고치지 않는다). 후속 DEV가 `EAI_AGAIN`·`ENETUNREACH`·`EHOSTUNREACH`를 retry 코드에 더하거나 문구를 좁히면 된다. 수락 차단 아님.

### 항목별 독립 검사 결과

| 항목 | 결과 | 근거 |
|---|---|---|
| 기본 off·legacy 보존 | PASS | 변수 없음·빈 값은 `null`. legacy 단독은 `sendPrompt` 동작·고정 문장 그대로 |
| 고정 HTTPS host/path | PASS | 대문자 host, 끝 `/`, 129자 ID, `.`, port, fragment, 끝 개행, `%2F`, 다른 path, subdomain, http, userinfo, `?x`, `..` 모두 무전송 |
| 비밀 파일/헤더 | PASS | `privatePath`(symlink 거부, 0640·0604 거부, 8192 초과 거부, 누락 거부). Bearer와 `X-Automation-Key` 동일 key. key 4097자·유니코드·CRLF 주입 거부. 호출마다 재읽기라 key 교체가 재시작 없이 적용 |
| redirect 차단 | PASS | 실제 fetch로 302 응답 시 두 번째 서버에 요청 0건, `rejected` |
| 8초 단일 시도 | PASS | `fromEnv`가 8000ms 고정. 무응답(300ms 모의)·연결 reset 모두 1회 `uncertain` |
| 중복 marker | PASS | 같은 batch 재tick은 POST 0건. 새 ID만 1회. uncertain은 재전송 없음. retry는 marker를 풀고 복구 뒤 다시 울림 |
| 미전달/불확실 분류 | PASS(F1 제외) | 2xx accepted_unverified, 3xx·4xx rejected, 5xx·무응답·reset uncertain, ECONNREFUSED·NXDOMAIN retry |
| 로그 비밀 미노출 | PASS | 결과·status·marker·wake-check 출력에 key·routine ID·메시지 본문 없음. 서버가 key를 echo해도 응답 body를 읽지 않는다 |
| received text 미전송 | PASS | body 키는 `event/ids/pending/next`뿐. UUID 아닌 ID 제거, 최대 10개. 본문이 든 메시지로 tick해도 POST에 없음 |
| webhook 배타 | PASS | webhook+agent 동시 설정 시 gateway 요청 0건. webhook 5xx 뒤에도 `sendPrompt` fallback 없음. `wake-check`도 webhook만 검사 |
| `wake-check` | PASS | 정상 `{"state":"ok","mode":"webhook"}` 종료 0. 불량 파일·상대경로는 종료 1, 무전송 |
| body event명과 routine 정의 | 인계 | event명 `knowslink_text_arrived`. owner가 routine prompt에 이 event명, `knowslink_text_receive` 호출, body untrusted 규칙을 넣어야 한다. 공식 문서상 `<webhook_event>`의 `body`는 JSON 문자열이므로 prompt가 파싱을 지시해야 한다. 실제 routine 정의와의 일치는 미검증이다 |

### 공식 규격 대조

공식 SKILL.md: URL `https://api2.cursor.sh/automations/webhook/<id>`(query 없음), POST, `Content-Type: application/json`, `Authorization: Bearer <key>`, `X-Automation-Key: <key>`, 8초, 재시도 없음, 200이 wake, body는 untrusted. 구현과 일치한다. 200 외 응답 의미는 공식에 없다. 4xx rejected·5xx uncertain은 기존 계약을 따른 DEV 판단이며 DEV 문서가 이를 밝혔다.

### C5 "optional webhook OFF" 해석 검토

C5의 webhook은 `architecture.md:37`과 `decisions.md:33`에서 relay 쪽 선택 기능이다. 재검토 조건은 endpoint 소유 확인·SSRF·wake 인증·재전송 방지다. 이 후보는 owner Bot 컴퓨터의 adapter가 owner 자신의 routine에 보내는 opt-in 기본 off wake이며 relay→vendor 전송이 아니다. 재검토 조건도 충족한다: endpoint는 owner 비공개 파일, SSRF는 host/path 고정과 redirect 차단, wake 인증은 sender key, 재전송 방지는 ID marker. wake는 ACK·실행·승인 권한이 아니다. 사용자의 "대체 방식 확인" 요청 범위 안이라고 판단한다. 다만 제품 해석이므로 coor가 designer 확인 여부를 정한다. DEV 문서의 같은 단서에 동의한다.

## 검증 및 남은 제약

| 검사 | 위치·대상 | 실행 주체 | 결과 |
|---|---|---|---|
| 고유 Linux clone | `/tmp/knowslink-webhook-tester-b86f397-s254e09d8/repo`, bundle에서 `b86f397` detach, Node v22.22.2, go1.27.1 | tester | `git rev-parse HEAD` 일치, status 0줄 |
| `lint.py --from f0fdf9d` | 같은 clone | tester | 종료코드 0, ERROR 0, WARNING 1(DEP-001), 실행 불가 0. 결과는 `lint.json` |
| product-lint `make lint` | 같은 clone | tester(lint.py 경유) | passed, 종료코드 0 |
| product-test `make test` | 같은 clone | tester(lint.py 경유) | passed, 종료코드 0. Go race, adapters 7개 테스트(grok-wake 포함), shell guard |
| 독립 경계 프로브 | `adapters/dist`, 실제 fetch와 loopback 서버 | tester | 58 PASS / 0 FAIL, 종료코드 0. 정보 2건(F1) |

- DEV의 Linux gate 증거(SHA `f39b677`)는 재사용하지 않았다. lint/test는 `b86f397` 바이트에서 직접 실행했다.
- 증거: `qa-reports/SAR-GROK-WEBHOOK-001-TESTER-test/webhook-boundary-probe.mjs`, `probe-and-lint.log`. 프로브의 key는 합성 값이다.
- 미검증(사유): 실제 Bot routine·key 발급·probe·routine 턴 표시·`knowslink_text_receive` 왕복은 운영·coor 후속이며 Grok UI를 열지 않았다. Windows ACL 경로는 Bot이 Linux이고 기존 `privatePath` 코드라 제외했다. 실제 api2.cursor.sh 응답 코드는 확인할 수 없다. 부하·장시간·E2E는 lite 범위 밖이다. main/배포는 미수락 그대로다.
- 임시 Linux 경로는 삭제 금지 지침에 따라 보존한다. snapshot은 `review.py cleanup` 전까지 유지한다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `phases/SAR-GROK-WEBHOOK-001-DEV.md` 기준·기술 계획, D03 webhook 절 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | D10 책임·검증 표, `adapters/README.md` webhook 절 | 확인 | DNS 문구 정밀도(F1) |
| 실행·검증 방법과 증거 | phase 검증 표, 이 보고서, `SAR-GROK-WEBHOOK-001-TESTER-test/` | 확인 | 없음 |
| 운영·복구 | D12 "Bot webhook routine wake 전환" | 확인 | probe는 coor 수동 단계로 명시됨 |
| 다음 작업·담당·재개 조건 | phase "미검증과 후속" | 확인 | 실제 routine 시험은 coor·owner 후속 |
| 로컬 링크·절 접근/지원 한계 | ops-guide→README anchor | 확인 | 대상 헤딩 존재 확인. 렌더러별 anchor 지원은 미점검 |
| snapshot 정리 후 정본 접근 | 리뷰 폴더와 `-test/` | 확인 | 결과·보고서·lint는 snapshot 밖 tester checkout에 있다 |

## 검토 결론

수락 가능(조건부). `b86f397`의 webhook wake 후보는 요청 범위의 경계(기본 off, 고정 URL, 비밀 파일·헤더, redirect 차단, 8초 단일 시도, marker 중복 방지, 분류, 비밀·본문 미노출, legacy 보존)를 독립 실행으로 충족한다. 미해결 critical/high는 없고 low 1건(F1)이 있다. 실제 Bot routine·key·probe는 미검증이며 main 미수락·미배포다. check 통과는 기록 검사이며 이 판단은 AI 검토 결과다.
