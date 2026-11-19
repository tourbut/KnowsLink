---
title: SAR-MVP-001-REVIEW-FINAL — RF-01 수정 독립 최종 코드 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-REVIEW-FINAL]
summary: RF-01 해소와 고정 SHA 독립 코드 리뷰 및 기존 증거 재사용 범위를 기록한다
---

# SAR-MVP-001-REVIEW-FINAL — RF-01 수정 독립 최종 코드 리뷰

## 판정

고정 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 코드 리뷰 수락이 가능하다.
신규 human 전달 우회와 RF-01 legacy claim 우회는 해소됐다. 미해결 critical/high는 0건이다.
이 판정은 호스트 AI의 내용 검토다. OCR은 파일 선택·규칙 해석과 기록 검사만 수행했다.
독립 tester의 새 SHA narrow QA 확인과 제품 최종 수락은 coordinator 후속이다.
공개 정책 제안·실벤더·실데이터·실제 배포의 기존 held를 유지한다.

## 고정 범위와 독립성

- base / merge-base: `0e4b5de7fa5eccb6a4391dedce3742e1845d1d00`. OPS 독립 리뷰 `c0e37c0`는 이 main 기준에 이미 수락됐다.
- head: `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`. RF-01 제품 수정 커밋은 `3c8dc73333f063f5836151bc6563c969f1776513`이다.
- snapshot: `/tmp/knowslink-mvp-review-78b1d92`. 작업 시작 전에 존재한 별도 snapshot의 head·detached·clean을 직접 확인했다.
- 추적 파일의 쓰기 비트를 제거했다. 검증 뒤에도 추적 파일 변경과 쓰기 비트는 0개다.
- 기록 체크아웃: `/home/shin/orca/workspaces/KnowsLink/fullops-dev`.
- 구현자: 실제 Claude session `1043ed6a-289c-456a-a1bb-aeaf0b0d4a5d`.
- 검토자: 실제 Codex session `01a10135-a619-7441-85a7-2ac49e68dad9`, `gpt-6.1-sol`.
- 이전 구현 세션 `4b9e5e7e-4fa5-442a-8a7a-d9c1d3c51c33`와 `01a10055-e46e-7d80-9ef1-12c1aabc69ce`도 보존했다.
- provider JSONL의 실제 session ID와 cwd를 확인했다. 상세 근거는 [independence.json](independence.json)이다.
- OCR `open-code-review v1.12.11 (a758d9c)`. rule SHA256 `ab2116fba36e3debc281c1c96acc502294b96a7db28fdd8cc2df3a3bceb92aa0`.

제품·기술/기획 정본·PLANS·board를 수정하지 않았다.
설치·빌드 결과는 snapshot의 Git 제외 경로에만 생성됐다. probe는 `/tmp` Go overlay를 사용했다.
새 후보의 GREEN은 제품 코드를 변경하지 않은 실행이다.
RED만 이전 `5002db6`의 store.go를 외부 overlay에 매핑했다. 이는 현재 후보의 실행 PASS가 아니다.

## 적용 기준

`fullops-review`, `open-code-review-delegate`, `fullops-work`, `orchestration`과 정규 DEV 인박스를 적용했다.
공통 기준은 `fullops-common-0.3.2`다. FULLOPS·공통 README·coding-style/testing/security·project·문서 작성 규칙·review/rule을 읽었다.
snapshot과 기록 체크아웃의 규칙은 같으며 SHA256을 [evidence.json](evidence.json)에 기록했다.

제품 원천은 D02 `docs/planning/product-specs/SAR-MVP.md`, service-design `7bc9ea1`의 frozen protocol C1–C5다.
D03 architecture/tech-stack, D05 interface-design, D06 data-model, D07 database-design, D08 db-schema, D09 crud-design, D10 module-design의 이전 검토 근거를 재사용했다.
RF-01에서 변경된 D03/D05/D06/D09/D10과 DEV 실행 기록·완료 로그는 새 SHA의 diff와 구현을 대조했다.
API·라이브러리 변경은 없어 이전 탐색·API 근거를 재사용했다. 보안·high 차단 기준을 완화하는 예외는 없다.

## 전체 커버리지와 이전 리뷰 재사용

고정 base/head 전체 변경은 126개다. reviewed 101개, skipped 25개, pending 0개다.
reviewed 비율은 80.16%다. 모든 `(path,status)`의 처리 기록은 100%다.
개별 사유는 [result.json](result.json), Git blob 비교와 재사용 방식은 [evidence.json](evidence.json)에 있다.

이전 `SAR-MVP-001-REVIEW-FIX`의 base `ffca87c9…` / head `4262d02…` 전체 113개, reviewed94/skipped19를 명시적으로 재사용했다.
그중 새 범위에 남고 head blob이 같은 94개는 원래 판정과 사유를 연결했다.
새 SHA에서 새로 검토한 25개와 변경 DEV 누적 아카이브 1개 및 원시 로그 6개의 처리를 별도로 기록했다.
현재 base에 이미 들어간 OPS 자료는 현재 변경 목록에서 빠진다. 이를 누락으로 처리하거나 재검토했다고 주장하지 않는다.

skipped는 이전 범위의 18개 생략과 DEV 누적 로그 1개 및 원래 리뷰 원시 로그 6개다.
누적 아카이브의 관련 DEV·실패 리뷰·RF-01 후속 전문은 읽었지만 초기 setup 전체를 다시 감사하지 않았다.
원시 로그는 고정 내용·SHA256·exit 표식·보고서/manifest의 정합성을 확인했다.
원래 제외 증거 21개의 SHA256은 기존 evidence와 일치한다.
이 생략은 이전 전체 QA를 새 후보의 독립 QA PASS로 옮기는 근거가 아니다.

원래 실패 `SAR-MVP-001-REVIEW-FIX-review`와 interrupted `SAR-MVP-001-INTEGRATION-review`의 26개 추적 파일은 고정 head와 바이트 단위로 같다.
원래 RF-01 미해결 high와 check exit1을 그대로 보존했다. 새 결과 key만 이번 SHA의 해소를 기록한다.

## RF-01 해소와 저장 호환

`store.go:271–278`의 `parentRouting`은 저장된 `Message.Deliver == "agent"`를 요구한다.
`authorize`와 `gate-consume`는 `parentFor`를 통해 이 경계를 실행한다.
`relay.approval.request` H와 `relay.result` R의 ingest도 같은 경계를 실행한다.
검사는 부모 token 소비·새 gate·새 R·completion 변경 전에 수행한다.
실제 Service transaction은 오류 시 원래 상태를 다시 읽고 안전 cleanup만 저장한다.

실제 이전 `a6a10c7` State 메서드로 발급한 fixture는 human 요청 2개와 agent 요청 1개의 claimed·delivered 상태를 담는다.
이전 JSON에는 Deliver가 없어 새 구조에서 빈 값으로 읽힌다. private key는 fixture 전용 합성 값이다.
fixture의 부모 봉투 서명을 원래 등록 공개키로 검증했다. 키 교체나 kid 재할당을 사용하지 않았다.
승인 gate fixture는 이전 owner POST 처리기와 같은 필드 변경이며 실제 owner HTTP 실행 증거로 주장하지 않는다.

같은 unit 회귀를 이전 store.go overlay에서 실행하면 legacy human authorize가 허용돼 exit1이다.
새 head 그대로 실행하면 authorize·R·H·gate-consume 거부와 상태 무변경, 같은 loaded State의 새 agent 정상 흐름이 exit0이다.
추가 test-only overlay는 human과 gated 부모의 저장 Deliver를 명시적 human으로 설정하고 같은 경계 거부를 확인했다.
TTL probe는 원래 signed route·TTL과 301초 이후 sweep의 token·원문/inbox 제거 및 gate expired를 확인했다.

실제 Postgres/HTTP `legacy_unrouted_claim_parent_boundaries`는 경로 미기록 부모의 네 경계 403과 무변경을 확인했다.
동일 claim의 경로만 agent로 복구하면 authorize·R·consume 200이다. 다른 인증 장애 때문에 거부된 것이 아님을 확인했다.
이 HTTP 테스트는 실제 이전 DB import가 아니라 유효 claim의 저장 경로를 지운 회귀다.
실제 이전 serialized fixture는 unit에서 검증했다. 두 근거의 범위를 구분한다.

경로 미기록 legacy agent도 fail-closed로 거부한다. 이는 C1 권한 불명 차단과 C3 재시작 후 claim 재발급 금지에 맞는다.
원문 서명에서 경로를 추정해 실행권을 복구하지 않는다. SQL/wire/저장 구조를 다시 바꾸지 않는다.
전환 시 이미 수락된 일부 요청은 잔여 TTL 동안 처리 결과를 내지 못할 수 있다. 원래 요청 수명은 최대 300초다.
300초가 자동 재전송 성공이나 완료 시간 보장이라는 뜻은 아니다.
같은 ID/key의 replay는 receipt-only이며 새 claim을 만들지 않는다. 새 시도에는 별도 신규 서명 요청·ID/key가 필요하다.
실제 외부 효과의 복구·exactly-once는 보장하지 않는다. payload 정리는 기존 transaction sweep 때 수행하며 wall-clock 즉시 삭제를 주장하지 않는다.

## 신규 C1와 정상 경로 회귀

새 후보는 저장 route를 lease·persist·ACK·claim에서도 검사한다. H 이외 direct human send는 403이다.
`agent_cannot_process_human_delivery`의 실제 DB/HTTP 회귀가 통과했다.
원래 direct human inbox 미지원에 대한 안전 차단 판정을 재사용한다. frozen human wire 값과 H owner gate를 유지한다.
신규 경계를 우회해 agent 전달로 바꾸거나 owner 승인을 대체하지 않는다.

정상 agent claim·authorize·R, owner gate CSRF/결정/consume·실행/공개 false, 현재 key/pair/owner·철회·세대·TTL·lease/ACK·DB 경합 회귀가 통과했다.
H/R 부모 binding, optional result/error·positive done 차단, strict wire와 registry의 변화 없는 근거는 이전 전체 검토에 연결했다.
RF-01과 원래 신규 C1 high는 새 result에서 resolved true다. 추가 미해결 critical/high는 발견하지 않았다.

## 새 실행과 종료코드

[execution.json](execution.json)과 [probe-execution.json](probe-execution.json)은 실제 subprocess 종료코드를 보존한다.
로그 출력은 명령에 직접 연결했다. `| tail`로 종료코드를 가리지 않았다.

| 검사 | 결과 | 근거 |
|---|---|---|
| make install | exit0 | install.log |
| 기존 store.go overlay + 현재 legacy 회귀 | 의도한 RF-01 RED exit1 | legacy-red.log |
| 현재 head legacy unit/race, -count=1 | GREEN exit0 | legacy-green.log |
| explicit human·signed fixture·TTL test-only overlay | exit0 | legacy-human-ttl.log·legacy-human-ttl-probe.txt |
| make test | exit0, unit/race | unit.log |
| make verify-mvp | exit0, 실제 별도 Compose/Postgres/HTTP·TS·Go UI seed | mvp.log |
| lint.py --from 0e4b5de… | exit0, product-lint passed, ERROR0/WARNING3/실행 불가0 | lint.json·lint.log |
| deliverables.py --repo snapshot --strict | exit0, 검사13/미작성2/문제0/경고0 | strict.log |

SIZE-001 WARNING은 http.go 487줄·integration_test.go 478줄·store.go 420줄이다. lint 계산 줄 수다.
기존 큰 파일에 대한 비차단 경고이며 이번 기록 과제에서 제품 분할을 하지 않는다.
make build는 verify-mvp 선행 target으로 실행됐다. 고유 Compose project `knowslink-mvp-69416484a9`만 회수했다.
make verify/verify-runtime, generate/schema 재생성, 전체 tester probe·새 캡처는 반복하지 않았다.
관련 setup·SQL·생성물·환경 파일은 이전 head와 같고 기존 실행 근거를 유지한다. 새 실행 PASS로 쓰지 않는다.

## QA·시각·기획·공개 후속

coordinator `msg_9aa228d35a66`는 이전 4262d02 QA `36bd4ae` 완료를 전달했다.
그 QA는 신규 경계·정상 agent/owner/current-auth 통과와 legacy high를 확인했다는 전달이다.
새 78b1d92 narrow QA는 같은 TESTER-FIX 후속으로 착수했고 `SAR-MVP-001-TESTER-FINAL.md`를 예정했다.
최종 tester 결과는 이 고정 head에 없으며 아직 수신하지 않았다. 본 리뷰의 자동 회귀가 독립 tester 전체 수락을 대체하지 않는다.

기존 QA `a6a10c7`의 31 pass·8 held·0 fail은 원래 SHA·조건으로 보존한다.
QA-06 stale epoch의 원래 targeted 근거는 SQL·생성 SaveRelay·transaction/current/sweep의 hash 동일성으로 재사용한다.
이는 새 tester의 stale epoch 실행이나 원래 held 해제가 아니다.

기존 designer UI `e238777`의 직접 PNG 검수는 unchanged 화면 범위에서 재사용한다.
page/gateView/gateDecision, Dockerfile/Compose, TS adapter와 기존 PNG 무결성을 직접 확인했다.
이번에 이미지·화면을 다시 열거나 캡처하지 않았다. 원래 직접 시각 PASS를 새 검수로 바꾸지 않는다.
보안 경계 변경은 회귀와 코드로 판정했다.

D04와 공개 기획 문서 결과는 별도 문서 수락·통합 대상이다. D04 인덱스는 이번 head에서 review로 연결됐다.
DEC-03 한도·rate·concurrency와 budget 제안은 사용자 승인 전이며 구현 완료·확정 정책으로 쓰지 않는다.
DEC-02·positive silent done·실벤더/실데이터·공개 처리량·WAL/backup 완전 삭제·운영 공개 held를 유지한다.
실제 배포 재개와 main 병합은 coordinator가 필수 QA/수락 조건을 확인해 처리한다.

## 기록 검사와 완료

새 key의 review.py check는 고정 base/head·누락 없는 coverage·실제 독립 세션·snapshot·lint를 검사한다.
check는 exit0으로 통과했다. check 통과는 내용 검토나 테스트 성공을 자동 보증하지 않는다. 실제 결과는 [check-status.json](check-status.json)에 기록했다.
현재 정규 DEV 인박스에 완료 전문을 작성하고 work.py finish를 실행했다.
기존 SAR-MVP-001-REVIEW-FIX key 중복으로 exit1이 발생했다. 원래 실패 리뷰 아카이브는 보존했다.
현재 지시서의 명시적 예외에 따라 실패 근거와 이번 지시서·완료 전문을 2026-10-03_to_dev.md의 별도 후속에 append했다.
append 전문과 인박스 무변경을 검증한 뒤 인박스를 비웠다. [archive-status.json](archive-status.json)에 결과를 기록했다.
소유 리뷰·inbox/logs만 커밋한다. 제품·PLANS/board·원래 실패 리뷰는 보존한다.
커밋 뒤 기준 base의 lint와 strict·공백을 확인하고 새 Dispatch로 worker_done을 한 번 보낸다.
