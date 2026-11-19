---
title: SAR-MVP-001-REVIEW-FIX 독립 코드 리뷰
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-REVIEW-FIX]
summary: 신규 human 전달 우회 수정과 기존 human claim의 잔여 C1 high 및 수락 차단을 기록한다
---

# SAR-MVP-001-REVIEW-FIX — 고정 수정 후보 독립 리뷰

## 판정

고정 후보 `4262d02fdd7b0b57a804d9e550597852950ffeae`의 수락을 차단한다.
신규 human 전달의 send·lease·persist·ACK·claim 우회는 수정됐다.
그러나 수정 전 발급된 human claim을 새 코드가 authorize·result에서 계속 허용한다.
미해결 high 1건이다. critical은 발견하지 않았다.
제품 파일을 수정하지 않았다. 구현 담당자 수정과 새 고정 SHA의 독립 재검증이 필요하다.
이 판정은 호스트 AI 검토 결과다. OCR의 자동 내용 판정이 아니다.

## 고정 범위와 독립성

- base / merge-base: `ffca87c9ec501e9f313c50f60124b686d27af728`.
- head: `4262d02fdd7b0b57a804d9e550597852950ffeae`. 제품 코드는 DEV 완료 `d7e2149`와 같은 후보다.
- snapshot: `/tmp/knowslink-mvp-review-4262d02`. detached·clean을 확인했다.
- 기록 체크아웃: `/home/shin/orca/workspaces/KnowsLink/fullops-dev`.
- 검토자: Codex `gpt-6.1-sol` high, 실제 session `01a1011b-a9a1-7673-94e8-86d370d9684b`.
- 수정 구현자: Claude session `4b9e5e7e-4fa5-442a-8a7a-d9c1d3c51c33`. provider JSONL의 sessionId를 직접 확인했다.
- 최초 구현자: Codex session `01a10055-e46e-7d80-9ef1-12c1aabc69ce`. 현재 검토자는 두 구현 세션과 다르다.
- 상세 세션 근거와 읽기 전용 상태: [independence.json](independence.json).
- OCR: `open-code-review v1.12.11 (a758d9c) linux/amd64`.
- OCR rule SHA256: `ab2116fba36e3debc281c1c96acc502294b96a7db28fdd8cc2df3a3bceb92aa0`.

snapshot 추적 파일은 변경하지 않았다. 설치·빌드 결과는 Git 제외 경로에만 생겼다.
추가 재현은 `/tmp`의 Go test-only overlay를 사용했다.
새 head의 제품 메서드는 그대로이며 테스트 함수만 외부 임시 파일에서 추가했다.
기존 코드를 고치는 overlay는 사용하지 않았다. 이전 구현의 fixture 생성 단계만 원래 `a6a10c7` 제품 소스를 썼다.

## 적용 기준과 요구사항

`fullops-review`·`open-code-review-delegate`·`fullops-work`와 현재 역할 지시서를 적용했다.
공통 규칙은 `fullops-common-0.3.2`다.
FULLOPS·공통 README 및 coding-style/testing/security·project·문서 작성 규칙·review/rule을 확인했다.
snapshot과 기록 체크아웃의 같은 문서가 바이트 단위로 일치한다. SHA256은 [evidence.json](evidence.json)에 기록했다.

제품 기준은 D02 `docs/planning/product-specs/SAR-MVP.md`와 고정 원천 service-design `7bc9ea1`의 `protocol.md` C1–C5다.
D03 architecture/tech-stack, D05 interface-design, D06 data-model, D07 database-design, D08 db-schema, D09 crud-design, D10 module-design을 대조했다.
DEV 최초 완료 전문과 human 수정 후속, DEV 실행 기록, TESTER QA·시나리오·probe, designer D04·직접 UI 기록, 공개 기획 기록을 검토했다.
기존 탐색·API 근거를 재사용했다. 신규 외부 SDK 선택이나 제품 정책 결정은 하지 않았다.
공통 기준·보안·critical/high 차단을 낮추는 예외는 없다.

## 파일 커버리지와 생략 영향

전체 변경은 113개다. OCR 대상은 92개이며 제외는 21개다.
최종 `result.json`은 reviewed 94개, skipped 19개, pending 0개다.
전체 변경 대비 reviewed 비율은 83.19%다. 모든 `(path,status)`의 처리 기록 비율은 100%다.

제품 소스·테스트·SQL·생성 Go 코드·Makefile·의존성 변경은 모두 직접 검토했다.
OCR의 unsupported_ext였던 go.mod/go.sum도 직접 검토했다.
제외된 원래 targeted.log는 직접 읽었다. 보완 pending PNG와 revoked PNG는 `view_image`로 직접 관찰했다.

생략은 누적 역할 아카이브 3개, 나머지 원시 로그 10개와 PNG 6개다.
아카이브의 MVP 완료 전문·수정 후속·SHA·held는 발췌 대조했지만, 초기 setup 등 과거 누적 전문을 다시 감사하지 않았다.
원시 로그는 존재·SHA256·exit 표식·최종 요약을 검증하고 원래 QA 보고서와 대조했다.
PNG 7개는 원래 QA `c59537b`의 Git blob과 일치하며 크기도 확인했다.
다른 PNG의 새 직접 열람·새 캡처·전체 QA 재실행은 주장하지 않는다.
개별 생략 사유는 [result.json](result.json), 증거 무결성은 [evidence.json](evidence.json)에 있다.
이 생략은 기존 실행의 전 항목을 새 후보 PASS로 확장할 근거가 아니다.

이전 `SAR-MVP-001-INTEGRATION-review/report.md`는 중단 템플릿이다.
해당 result의 pending은 완료로 바꾸지 않았다. 원래 기록을 보존했다.
새 리뷰만 이번 head의 수락 또는 차단 결론을 제시한다.

## 발견 사항

### C1 기존 신규 human 전달 우회 — high, 해당 경로 해소

원래 reviewer `msg_4fbcac80f76c`는 `deliver:human` 요청이 agent pull/persist/ACK/claim으로 delivered가 되고 owner gate는 0개임을 재현했다.
원래 head는 `31405e736a16be9d77239c6cdc6fdeb56892436f`다.
원시 증거는 이전 리뷰의 `targeted.log`이며 원래 실행 SHA로 유지한다.

수정 후보는 `Message.Deliver`를 저장한다.
leaseMessage·leased·claim은 `Deliver==agent`를 요구한다.
H 외 직접 human send는 서명·발신 인증 뒤 `403 sender_not_allowed`로 거부한다.
새 후보의 실제 Postgres/HTTP 회귀 `agent_cannot_process_human_delivery`가 통과했다.
저장된 human 요청의 위조 lease에 persist/ACK 409와 claim 403도 확인했다.
다만 이미 발급된 claim의 아래 잔여 문제 때문에 C1 전체 해소는 선언하지 않는다.

### RF-01 기존 human claim의 부모 처리 우회 — high, 미해결

위치: `internal/relay/store.go:271–276`, 관련 호출은 `http.go`의 authorize/gate-consume와 ingest의 H/R이다.
`parentRouting`은 delivered·Claimed·현재 권한을 검사하지만 `Deliver==agent`를 검사하지 않는다.
이전 구현은 human 요청을 agent에게 claim할 수 있었다.
그 상태를 새 구조로 읽으면 `Deliver`는 미기록 값이지만 기존 ClaimToken과 Claimed는 유지된다.

재현은 원래 `a6a10c7`의 실제 State ingest/leaseMessage/operate로 signed human 요청을 전달·persist·ACK·claim한 뒤 상태를 직렬화한다.
새 head의 그대로인 State 코드로 읽고 동일 B credential과 기존 token을 사용한다.
새 claim 호출은 거부된다. 기존 token의 authorize는 성공한다.
새 B 서명 `relay.result`도 수락되고 부모 completion은 denied이며 owner gate는 0개다.
원래 B private key를 임시 파일에서 재사용했다. kid 재할당이나 공개키 교체로 재현하지 않았다.

[legacy-claim.log](legacy-claim.log)의 결정적 출력:

```text
HEAD: new claim correctly denied
REPRODUCED: old deliver=human claim authorizes and accepts agent result; completion=denied; owner gates=0
```

fixture 생성과 새 head 검사는 각각 종료코드 0이다.
이는 결함 존재를 확인한 관찰 테스트 성공이다. 보안 조건 PASS가 아니다.
재현 함수는 [legacy-claim-probe.txt](legacy-claim-probe.txt)에 보존했다.
새 HTTP 서버/DB를 통한 이 잔여 결함 재현은 하지 않았다.
현재 HTTP authorize·send 호출이 그대로 이 메서드를 실행하며 성공 시 200을 반환하는 연결을 직접 확인했다.

영향은 수정 전 DB·claim을 유지하는 전환에서 agent가 human 전달의 처리 결과를 계속 확정하는 C1 위반이다.
실제 외부 도구 실행·일정 공개는 false다. 그런 효과나 새 human lease가 발생했다고 주장하지 않는다.
필요한 수정은 기존 human/경로 미기록 claim의 차단과 모든 부모 처리 경계의 agent 경로 확인이다.
구현 담당자는 legacy claimed human·경로 미기록 부모의 authorize/H/R/consume 거부와 정상 agent 회귀를 추가해야 한다.
새 snapshot과 새 review key에서 독립 재검증한다.
coor에 escalation `msg_5d70016552f0`으로 전달했다.

## Frozen direct human 제한의 정합성

wire는 여전히 agent/human 두 값을 명시적으로 파싱한다. `deliver:both`는 거부한다.
seed의 schedule.commit human은 전형적 경로이며 실행 가능한 도구를 뜻하지 않는다.
현재 로컬 합성 요구는 agent M 수신 뒤 새 human H gate를 검증한다.
직접 human inbox가 없는 구현이 H 외 human 요청을 403으로 거부하는 것은 지원하지 않는 경로의 안전한 차단이다.
agent로 변경·자동 전달·owner 승인 대체·stub 실행 허용이 없으므로 frozen 보안 기준을 낮추지 않는다.

이 판정은 direct human inbox를 완성했다거나 human wire 값을 폐기했다는 뜻이 아니다.
직접 human 업무 경로의 추가 지원은 후속 요구·구현·독립 수락 대상이다.
이번 리뷰에서 새 제품 결정이나 코드 변경을 만들지 않았다.
새 제한만으로 기존 claim의 유효성을 자동 취소하지 못하는 RF-01은 별도로 차단한다.

## C1–C5 검토와 QA 조건

| 기준 | 확인한 구현·증거 | 현재 한계 |
|---|---|---|
| C1 | owner/agent credential 분리, owner-bound PoP, `(from,kid)` lookup, rotate/revoke, current key/pair/owner, 세대·epoch CAS, 각 신규 delivery 경계 | RF-01 high 미해결. 실제 사용자 신원과 공개 인증 held |
| C2 | 권위 있는 M·claim·digest·endpoints·generation·exp 결속, H 재귀·중복 거부, verified typed body+정책, owner POST·HMAC CSRF, 원자적 decision/consume | approve가 실행·공개를 허용하지 않음. legacy 부모 경계는 RF-01 |
| C3 | auth/routing 이후 digest·atomic idempotency·TTL/id·queue+receipt, receipt-only replay, row lock·epoch CAS, persist/ACK/claim, 30s·3회·마지막 lease·경합·재시작·비정상 시계 | 실제 외부 효과 exactly-once를 주장하지 않음 |
| C4 | query 기본 비공개, result endpoint 반전·현재 부모 권한, optional result/error·done 거부, exp/철회/result 뒤 payload 정리, metadata 24h | 정책·실데이터 positive done·WAL/backup 완전 삭제 held |
| C5 | duplicate/Unicode/unknown/null/비정규 값 차단, 동일 parsed Raw의 검증·처리, hash-locked registry, 구현 subset, evidence fetch/preview·webhook OFF | 자원/rate/추가 size/concurrency 수치와 singleton 공개 처리량 held |

기존 QA는 `a6a10c7`에서 31 pass·8 held·0 fail이다.
새 후보와 달라진 제품 파일은 store/http/integration_test/protocol_test 4개다.
따라서 기존 QA 전체를 새 후보의 독립 QA PASS로 옮기지 않는다.
이번 snapshot 실행은 실제 변경 경계의 회귀 근거이며 tester 후속을 대체하지 않는다.
coor는 Grok QA `task_33e3336872e4 / ctx_fcf73eae42eb`를 같은 후보에 배정했다고 전달했다.
새 QA 결과는 현재 고정 snapshot에 포함되지 않았으며 아직 수신하지 않았다.
RF-01 때문에 정상 신규 경로 QA 통과만으로 제품 수락을 열 수 없다.

QA-06 stale epoch의 원래 targeted 검사는 zero-row SaveRelay와 HTTP 503·owner signup rollback을 관찰했다.
SQL, 생성 SaveRelay, transaction/current/sweep이 원래 `31405e7`과 일치한다.
그 동일성을 evidence.json에 기록했다. 변경 없는 targeted 근거는 원래 SHA로 재사용할 수 있다.
이번 DEV 자동 DB 회귀는 경합·철회·세대·시계 차단을 다시 검사했다.
이는 고의 stale epoch를 새 tester가 실행했다는 뜻이 아니다. 원래 tester held는 소급 수정하지 않는다.

## 직접 시각 증거와 문서 통합

기존 designer는 `a6a10c7` 제품과 `c59537b` QA 증거의 PNG 7개 및 보완 PNG 1개를 직접 검수했다.
보완 pending은 기존 저장 HTML의 1280×1100 정지 렌더이며 새 서버 실행이 아니다.
이번 검토자는 full pending과 revoked를 직접 열었다. typed body·정책·두 POST 버튼·철회 후 원문 부재와 버튼 차단을 확인했다.

page template, gateView, gateDecision, Compose·Dockerfile·TypeScript adapter는 이전 고정 후보와 같다.
기존 PNG 7개도 원래 QA blob과 일치한다. 시각 판정의 unchanged 화면 범위는 재사용할 수 있다.
claim 처리 경계 변화와 RF-01은 이미지로 증명하지 않는다.
401 invalid_auth는 인증 실패의 자료 미노출·승인 불가만 증명한다. 인증된 owner의 권한 저장소 장애 전체를 증명하지 않는다.

D04·기획 공개 기준의 문서 결과는 제품 수락과 분리해 통합 검토할 수 있다.
공개 가입·hostname의 승인과 DEC-03 수치 제안을 구분했다.
권장안의 32KiB·quota·rate·concurrency 등을 확정 정책이나 구현 완료로 바꾸지 않았다.
안전 정리 budget·부모/H/R 종료·NAT/source IP·서버 자원 측정은 담당 DEV/OPS 근거와 사용자 결정 후속이다.
Free N·가격·실벤더·실데이터·positive silent done·운영 공개 held를 유지한다.
D04 원천은 review로 존재한다. 공유 인덱스와 D02의 이전 미작성 문장은 소유 경계로 보존된 이력이다.
strict의 미작성 4를 새 D04 부재로 오인하지 않는다. coor가 통합 때 인덱스 정합성을 반영한다.
제품·기획·PLANS·board·기존 리뷰는 수정하지 않았다.

## 실제 실행과 종료코드

새 실행 위치는 위 readonly snapshot이다. 로그는 별도 기록 체크아웃에 저장했다.
subprocess.run의 실제 종료코드를 기록했다.
Docker 출력의 줄 끝 공백은 저장 뒤 정리했다. 상태·명령·출력 순서·종료코드는 보존했다. 명령 뒤 `| tail`로 종료코드를 가리지 않았다.

| 검사 | 결과 | 근거 |
|---|---|---|
| make install | exit 0 | install.log |
| lint.py --from ffca87c… --out lint.json | exit 0; product-lint passed; ERROR 0·WARNING 3·실행 불가 0 | lint.json·lint.log |
| make test | exit 0; unit/race | unit.log |
| make verify-mvp | exit 0; 실제 별도 Compose migration·Postgres/HTTP/race·TS adapter·Go UI seed | mvp.log |
| legacy claim test-only overlay, 원래 a6a10c7 fixture | exit 0; 실제 이전 메서드로 human claim 생성 | legacy-claim.log |
| legacy claim test-only overlay, 현재 4262d02 | exit 0; RF-01 재현 | legacy-claim.log |
| deliverables.py --repo snapshot --strict | exit 0; 검사 13·미작성 4·문제 0·경고 0 | 실행 확인; D04 인덱스 한계 위 참조 |
| review.py check --from ffca87c… --to 4262d02… | exit 1; 실제 오류는 critical/high 미해결 사항 | 기록 구조 및 high 차단 확인 |

SIZE-001 WARNING 3은 http.go 487줄, integration_test.go 449줄, store.go 420줄이다.
이는 lint 계산 줄 수이며 실제 전체 파일 줄 수와 구분한다. 이번 리뷰는 파일 분할을 구현하지 않는다.
make build는 verify-mvp의 선행 target으로 성공했다.
make verify, verify-runtime, generate/schema 재생성, 독립 tester 전체 probe와 새 캡처는 실행하지 않았다.
관련 setup·SQL·생성 경로는 직접 diff와 기존 근거를 대조했다. 새 실행 PASS로 기록하지 않는다.
검사 자원은 고유 knowslink-mvp project만 회수했다. 기존 서버 서비스·Tunnel을 변경하지 않았다.

## 완료와 후속

리뷰 기록은 완성하되, check의 미해결 high 실패를 통과로 바꾸지 않는다.
check는 refs·coverage·independence 검사 뒤 high에서 중단했다.
lint 후속 검사는 check에서 실행되지 않았다. 실제 snapshot lint 및 [record-validation.json](record-validation.json)의 직접 정합성 검사를 별도 근거로 남긴다.
현재 reviewer inbox에 완료 전문을 쓰고 fullops-work의 work.py finish로 로그를 보존한다.
소유 리뷰 파일·inbox/logs만 커밋하고 새 Dispatch에 outcome failed로 보고한다.
실패 의미는 후보 수락과 필수 review gate 미충족이다. 리뷰 수행·증거·차단 인계는 완료했다.

coor는 RF-01을 구현 담당자에게 인계한다. 수정 후 고정 head의 독립 QA와 새 독립 리뷰를 조정한다.
시각 증거의 재사용 범위와 기획 문서 통합을 전체 제품 수락·공개와 구분한다.
미해결 high·held·사용자의 현재 과제 완료 뒤 중지 지시는 유지한다.
