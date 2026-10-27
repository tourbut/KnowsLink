---
title: SAR-CONNECT-SIMPLE-001-DEV-REVIEW2 리뷰
status: draft
updated: 2026-10-27
owner: reviewer
summary: 0fff4ca07686 병합 전 독립 리뷰 보고서
---

# SAR-CONNECT-SIMPLE-001-DEV-REVIEW2 리뷰

- 검토자 / CLI / 모델: Orca dispatched reviewer / Claude Code / `claude-opus-5-5`. 실제 Claude 세션 ID `ef5b1ddf-2320-42ef-8001-224dff4c58d7`. 기록상 reviewer 세션 식별자는 prepare 값 `reviewer-ecfad14e-fbbb-4268-bb09-08a180378aa5`를 유지한다. 구현자 세션 `445c5ac6-1bb5-4e2e-a374-2700e51c3b87`과 다르다.
- base SHA / head SHA / merge-base: `3fa283e0c69208b9c4e81aa4c73873cefca20a6e` / `0fff4ca0768622c7e8f7128dca3d671372d5f144` / `ec33a3fc45692763d9ae3e261307a1ddd276e47e`. snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-c2abd5fcb84c40629f94ed889951b536`의 HEAD가 `0fff4ca0`이고 작업 트리가 깨끗함을 확인했다. snapshot은 바꾸지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.3 (375f64c). `rules.json` rule_sha256 `f1061481…94de8`. 공통 규칙 `.fullops-squad/rules/common/`(README, coding-style, security, testing)를 snapshot 기준으로 적용했다. OCR의 LLM 판정은 쓰지 않았다. 아래 판단은 AI 검토자의 판단이다.
- 요구사항·완료 기준 원천: [D02 SAR-CONNECT-SIMPLE-001](../../../planning/product-specs/SAR-CONNECT-SIMPLE-001.md) CS-AC-01–18, 지시서 `.fullops-squad/handovers/logs/2026-10-26_to_dev.md`, DEV 자체 검사 [report.md](../SAR-CONNECT-SIMPLE-001-DEV-test/report.md).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 57 / 57 / 0 / 52 / 5. skipped 5건은 prepare가 채운 하네스 기록(`docs/evaluations/jev/` 4건, `handovers/logs/` 1건)이다. 지시서 범위(relay 한도·시도·이름·화면, adapters 재시도·사유, 설치기, 문서) 밖 변경은 없다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 20 / 2. 실행 불가는 `make lint`·`make test`(Windows에 make 없음)다. 대체 근거는 DEV의 Linux 격리 clone `310fdd73` 실행(`make lint`·`make test`·`go test -race -tags=integration` EXIT=0)이다. `310fdd73` 뒤 커밋(`d1a401b5`, `0fff4ca0`)은 `.fullops-squad` 문서만 바꿨다(diff로 확인). 리뷰어는 lint·테스트를 다시 돌리지 않았다.
- SIZE-002: 지시서 예상 600–900줄 / 실제 코드 추가 1738줄(lint 집계). 새 테스트(Go 207, TS 440)와 문서가 절반 가까이다. 커밋이 relay·adapters·문서로 나뉘어 있어 분할 요청은 하지 않는다. SIZE-001 경고는 기존 300줄 초과 파일에 기능이 더해진 것이다. 이번 병합을 막지 않는다.
- DEP-001: `adapters/package.json`은 `test` 스크립트에 `connect-simple.test.js`만 더했다. 의존성 변경 없음. SEC-001(`connect-simple.test.ts:73`)은 fake relay의 `lease_token: "lease"` 고정 문자열이라 비밀값이 아니다. SLOP-002 3곳은 선택 필드 직렬화로 수락한다.
- UI 디자인: 회원 화면(`human_pages.go`, `member_pages.go` device·refusal) 변경이 있다. 디자인 정본은 D04 연결 단순화 절이다. 1440·390·320px·키보드 검증(CS-AC-17)은 designer 시각 검수로 PLANS.md에 대기 중이다. 이 리뷰는 템플릿 구조·문구·EN 번역 항목만 확인했다.

## 검토 범위

result.json의 모든 (path, status)에 reviewed/skipped와 사유를 기록했다. 중점 위험 경로별 확인 결과는 다음과 같다.

1. relay 한도 제거와 대체 장치(`connections.go`, `member_agents.go`, `store.go`, `http.go`, `remote_auth.go`): 회원당 5·기록 10·전체 200 상한이 모두 사라졌다. 남은 거절은 `agentLimit`의 `service_capacity`(활성 agent 2000)와 `agent_rate`(회원당 24h `Created` 20개, 철회 포함)뿐이다. agent 생성 경로 4곳(웹 추가, device complete, Dots 원격, owner API `operateAs`)이 모두 `agentLimit` 뒤 `addAgent`를 지난다. 그 밖의 `st.Agents[...] =` 쓰기는 없다(grep). 상태는 한 행 잠금 트랜잭션이라 동시 요청은 직렬화된다. 동의 시점 검사는 자리를 예약하지 않지만 complete가 다시 검사하므로 여러 시도를 동의해 20개를 넘길 수 없다. `Retry` 계산(`made[len-20]+24h`)은 정확하다. 철회 agent는 철회 뒤 24h 이상 남으므로 창 안의 생성 기록이 sweep으로 사라지지 않는다. L-05: 친구 전달은 `Friends` 기반이고 legacy pair는 명시 초대로만 생긴다는 DEV 분석을 코드에서 확인했다. `/v1/contacts/{own}` 허용은 자기 agent만이며 상대 label 403 테스트는 유지된다.
2. 키 저장 전 시도와 정리: `approveDevice`는 agent를 만들지 않고 ID·이름만 예약한다. `refused` 상태는 `fail()`이 오류 없는 값을 돌려 트랜잭션에 커밋된다. `/home/attempt-remove`는 로그인 세션, cross-origin 차단, `c.Device && c.Owner == owner && state != consumed`를 확인한다. cleanup admission도 같은 owner만 허용한다. 다른 회원의 시도는 지울 수 없다. 정리 뒤 complete는 연결 기록이 없어 401이며 고아 agent가 생기지 않는다. sweep은 approved device 시도를 기한+24h까지 두고 그때 지운다. 늦은 complete는 `open()`으로 허용된다.
3. Google 실패·링크 재방문: `connectFailure`는 다른 계정 시작, 가입 닫힘, 기한 지남을 서로 다른 문구로 보이고 상대 계정의 이메일·이름을 넣지 않는다. '다른 계정 선택'은 `/auth/google`로 가며 `prompt=select_account consent`를 쓴다. `devicePage`는 시작한 회원(세션 유효, `Members[id].Owner == c.Owner`)에게만 결과·이름·거절 사유를 보이고 그 밖에는 기존 401 일반 문구다. 새 KO 문구의 EN 항목이 `i18n_en_member.go`에 있다.
4. adapters 재시도 폴더와 ACL: `privateDirectory`는 EEXIST이면서 링크가 아니고 비어 있는 폴더만 받아들인다. 내용이 있으면 원래 오류를 던져 키를 덮어쓰지 않는다. `aside()`는 비어 있지 않은 폴더를 `rename`만 하고 지우지 않는다. `agent.json`이 있으면 fresh여도 먼저 connected로 끝난다. Windows ACL은 소유자가 이미 현재 사용자일 때만 `SetOwner`를 건너뛴다. 이어서 `privatePath`가 소유자·허용 규칙을 그대로 검증하므로 권한이 약해지지 않는다. `complete()`의 `agent.json` 쓰기는 `wx`라 기존 키 파일을 덮지 않는다.
5. 구버전 relay 422: 클라이언트는 알려진 host만 보낸다. 배포 순서는 D06 `interface-design.md:36`과 DEV 보고 '배포 주의'에 있다. D12에는 없다(발견 사항 3).
6. C-10: `mcp.test.ts`는 CRLF를 LF로 바꾼 뒤 고정 sha256 `2b25c6b5…`와 같은지 먼저 단언한다. 내용 고정은 이전보다 강하다. relay의 `RegistryValid()`는 embed 원본 바이트를 그대로 검사하고 `.gitattributes`의 `eol=lf`도 유지된다. 약화가 아니다.

테스트 삭제·완화: 제거된 단언은 회원당 5·기록 10·전체 200 거절, 자기 contact 403, 기한 시각 complete 거절이다. 각각 L-01·L-02·L-03, A-06, A-02의 새 규칙에 따른 것이다. 대체 테스트(`TestAgentRate`, `TestAttemptsAreNotAgents`, 기한+24h complete 거절, `concurrent_agent_rate_and_restart`, 상대 label 403 유지)가 있다. 검증 무력화는 찾지 못했다.

## 발견 사항

| 심각도 | 파일·줄 | 재현 조건 | 영향 | 수정 상태 |
|---|---|---|---|---|
| medium | `internal/relay/member_agents.go:486-489` | agent 19개 회원이 두 시도를 동의하고 둘 다 키를 저장한다. 두 번째 `/v1/connect/complete`가 `agent_rate`로 거절된다. | 응답이 `errorBody`를 거치지 않아 `retry_at`이 없다. 클라이언트는 시각 대신 '잠시'를 보이고 `next`는 없는 retry_at을 기다리라고 한다. D02 E `agent_rate`의 `<time>`과 CS-AC-04의 '다시 가능한 시각' 위반이다. | 미해결. 이 경로(472·488행)도 `errorBody(e)`를 쓰고 테스트 1개를 더한다. |
| low | `internal/relay/connections.go:175` | 배포 직전 10분 안에 구버전이 동의(agent 생성)한 시도의 agent를 회원이 철회한다. sweep이 24h 뒤 지운 뒤 기한+24h 안에 클라이언트가 complete한다. | `reserved`로 판정돼 같은 agent ID가 다시 생긴다. 'agent ID 재발급 없음'(C1) 불변식 위반이다. 키 없던 agent라 (from,kid) 재사용은 없다. | 미해결. 새 흐름 표식(예: `c.Name`)으로 한정하거나 agent 철회 때 그 device 시도를 취소한다. |
| low | `.fullops-squad/docs/operations/ops-guide.md:14-24` | 운영 담당이 D12만 보고 클라이언트·설치기를 relay보다 먼저 배포한다. | 구 relay가 `host`를 422로 거절하고 사용자는 `unreachable`을 본다. | 미해결. D12 절에 배포 순서 한 줄을 더한다. |
| low | `internal/relay/member_agents.go:312-319` | 해당 없음(테스트 공백). | `attempt-remove`의 다른 회원·consumed 거절에 음성 테스트가 없다. 현재 코드는 맞다. | 미해결. 음성 테스트 1개를 권한다. |

critical·high는 없다.

## 검증 및 남은 제약

- 리뷰어는 snapshot과 기존 증거만 읽었다. lint·테스트·임시 테스트를 실행하지 않았다(지시서 제약).
- 재사용 증거: `lint.json`(coor 투입, base·head 일치, ERROR 0). 등록 검사 `make lint`·`make test`는 Windows 실행 불가이며, DEV의 Linux 격리 clone `310fdd73` 결과(EXIT=0)와 Windows `npm run check`·`npm test`·`package_plugin.py --verify` EXIT=0을 DEV 보고에서 인용한다. 원시 로그는 gitignore된 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-DEV/`에 있어 이 리뷰에서 열지 않았다. `310fdd73`→`0fff4ca0` 차이가 문서뿐임을 git diff로 확인했다.
- 미실행 범위: 실제 Grok Bot 설치·T-04 실측·CS-AC-17 시각 검수·CS-AC-18 실연결은 tester·designer 후속이다(PLANS.md 대기). DEV 보고의 `mcp.test.js` 1회 원인 미상 실패는 21회 연속 통과로 기록돼 있다. 재발하면 원인 조사가 필요하다.
- jev packet·route JSON 4건과 지시서 로그는 skipped다. 존재는 result.json 목록과 DEV 보고의 링크로 확인했다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | D02 `product-specs/SAR-CONNECT-SIMPLE-001.md` 상태와 근거·제품 규칙 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | D10 `module-design.md` SAR-CONNECT-SIMPLE-001 절, D09 `data-model.md`, `crud-design.md` 새 절 | 확인 | 미완료(tester·designer)는 PLANS.md에 있다 |
| 실행·검증 방법과 증거 | `qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md` | 확인 | 원시 로그는 gitignore 로컬 경로라 다른 체크아웃에서 열 수 없다 |
| 운영·복구(해당 시) | D12 `ops-guide.md` 'agent 용량 경계' 절 | 확인 | 배포 순서 누락(발견 사항 3) |
| 다음 작업·담당·재개 조건 | `PLANS.md` SAR-CONNECT-SIMPLE-001-DEV 4줄 | 확인 | registry.json CRLF 재체크아웃은 사용자 결정 대기 |
| 로컬 링크·절 접근/지원 한계 | D02→D04·D10 anchor, D09→D12 anchor | 확인 | 한글 anchor는 GitHub 규칙 기준이며 렌더러별 차이는 미확인 |
| snapshot 정리 후 정본 접근 | coor 리뷰 폴더의 result.json·report.md·lint.json | 확인 | 없음 |

## 검토 결론

조건부 수락한다. critical·high가 없다. 핵심 위험 경로(한도 대체 장치의 경쟁·우회, 시도 정리 권한, 실패 화면의 정보 노출, 폴더·ACL 보호, C-10 테스트)에서 차단 결함을 찾지 못했다. medium 1건(complete 시점 `agent_rate`의 `retry_at` 누락)과 low 3건은 후속 수정 과제로 PLANS.md에 남기기를 권한다. `review.py check` 통과는 기록 검사이며 이 판단의 내용이나 테스트 성공을 자동 보증하지 않는다.
