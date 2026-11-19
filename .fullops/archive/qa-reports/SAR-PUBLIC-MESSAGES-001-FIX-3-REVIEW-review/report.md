---
title: SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW 리뷰
status: review
updated: 2026-10-07
owner: ops
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW]
summary: 고정 d089의 독립 정적 리뷰로 기존 입장 결함 수정 근거와 UTF-8 stdin medium 수정 요청을 기록한다
---

# SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW 리뷰

- 검토자 / CLI / 모델: OPS / Codex / 사용자 지정 gpt-6.1-sol high. 실제 `CODEX_THREAD_ID`는 `01a116a3-2a1d-7961-92fa-33417bc9c46c`다.
- 구현자 실제 세션: `aa85544d-18c3-43d3-95d0-b729aa9e9e8c`. 서로 다른 실제 세션이다.
- base / merge-base: `68b0d6a0c854fdaec6828a232dd3945814be1404`. head: `d08903a55c3638128827010400e66e9d45b61d7c`.
- 읽기 전용 managed snapshot: `/home/shin/orca/workspaces/KnowsLink/.fullops-review-4e43dd6838fe4b8e938312708eaf0bcd`. exact HEAD·detached·clean을 직접 확인했다. 제품 읽기는 이곳에서 했다. 실행·설치·파일 쓰기는 하지 않았다.
- 기록 checkout: 등록 OPS `/home/shin/orca/workspaces/KnowsLink/fullops-ops`, 준비 HEAD `cebc32c2e3ae3b15ff5fd7238de1c5ab96eaf7b4`. Run `run_8ca8bc058ab7`, task `task_886bd365fe36`, dispatch `ctx_4dade3f47e35`다.
- OCR: open-code-review v1.12.11 (a758d9c), delegate. 규칙 hash `ab2116fba36e3debc281c1c96acc502294b96a7db28fdd8cc2df3a3bceb92aa0`. 여섯 rule group을 읽었다. OCR LLM은 실행하지 않았다. 아래 판정은 검토자 AI의 판단이다.
- 적용 기준: fullops-common-0.3.3 README/coding-style/testing/security, FULLOPS/project/document-writing/orca-agents/contexts/ops. OPS 준비 HEAD와 제품 snapshot의 공통 규칙·project·review 규칙 diff는 없다. 예외는 없다.
- 요구사항: `docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` PS04/06/07/08–11, `docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` UX06/07. DEV 원본/FIX/FIX2/FIX3 실행 기록, 원 TESTER09c, OPS09c/dfc, 미완료 FIX2, designer09c/dfc 보고서를 원 SHA로 대조했다.
- 파일별 coverage와 총계는 [coverage.json](coverage.json)과 [result.json](result.json)에 고정했다. 전체 변경 702개의 `(path,status)`를 reviewed 또는 사유 있는 skipped로 기록했다. 제품 변경 32개는 모두 검토했다.
- 후보 [lint.json](lint.json): exact base/head, ERROR0 / WARNING14 / 실행 불가0. `product-lint`와 `product-test`는 각 passed / exit0이다. COOR 후보 원본과 byte/hash가 같다.
- SIZE-002: main 이후 누적 추가6325줄/예산400이다. 원 DEV 예상 제품/시험12–18파일, 실제20파일의 중앙 HTTP 경계·bundle 수정 근거를 읽었다. FIX2 예상 제품+120/−35·시험+280, 실제 제품+122/−35를 읽었다. FIX3 예상 제품+90/−30·시험+190, 실제 제품+95/−34·시험+188/−13을 읽었다. 누적 문서·원실패·검사·시각 증거를 포함한 범위다. 공통 입장·snapshot 수정과 회귀를 나누어 미완성 후보를 만들지 않은 판단은 타당하다. 원 증거와 coordinator 준비 기록은 줄이지 않았다.
- DEP-001: 해당 없음. `go.mod`, `go.sum`, adapter package/lock의 제품 diff가 없다. SQL 생성물은 기존 sqlc의 `ReadRelay` 17줄이며 새 의존성이 아니다. relayBase는 기존 core로 이동했다.
- UI 디자인: D04 일반서비스 UX·MESSAGES UI/FIX와 기존 memberStyle/Go template을 적용한다. 누적 변경에는 receipt/gate UI가 있다. FIX3은 template·CSS를 바꾸지 않았다. 별도 디자인 lint·공용 다중 테마가 없으므로 미실행/해당 없음이다(project.md). DESIGN 경고0은 Go 문자열 CSS의 시각 수락을 뜻하지 않는다. 최신 d089의 직접 영향 검수는 designer 몫이다.

## 검토 범위

기준 main 이후 미수락 변경을 diff와 최종 소스로 읽었다. `capacity.go`, `cleanup_admission.go`, `store.go`, `public_text.go`, `member_receipt.go`는 전체를 읽었다. HTTP/gate/member handler와 adapter transport·CLI·MCP의 변경 및 호출자를 읽었다. 변경 Go 시험·adapter 시험과 실제 SQL/생성물을 대조했다.

| 경계 | 직접 확인한 내용 | 판정 |
|---|---|---|
| 익명 slow body H-1 | 본문8/32KiB·socket/context10s 수신이 모든 슬롯보다 앞이다. 익명 body는 DB·입장 기록을 쓰지 않는다 | 정적 수정 근거 확인 |
| 자격·소유 H-2 | `cleanupOwner`가 commit snapshot의 kind/만료/자기 기록을 확인한 뒤 슬롯을 고른다. 타 owner·lease 없는 ACK는 신규다 | 정적 수정 근거 확인 |
| 공정성 | 로컬 LoadOrStore와 공유 enterHTTP 모두 동일 단위당1개다. owner 제어와 소속 agent ACK를 두 단위로 나눈다. global clean4는 유지한다 | owner 철회가 agent ACK 하나에 묶이지 않음 |
| snapshot 순서 | commit 뒤 immutable 상태를 게시한다. CAS가 `(DB 시각,epoch)` 이전/동일 값을 거부한다. refresh SELECT가 row epoch/clock/data를 함께 읽는다 | 같은 시각 역전 수정 확인 |
| stale/restart | 신규 채널이 full이고 자격을 증명 못할 때만 refresh한다. 읽기 동시1·취소10s·요청 도착 뒤 시작한 읽기 공유다. 잠금/입장 기록을 만들지 않는다 | 포화 직전 유효 자기 정리의 재판정 확인 |
| DB 원자성·종료 | row lock·epoch CAS·commit·rollback을 유지한다. 실패 작업은 원 상태/sweep/rate로 복원한다. 실패 finish는 orphan으로 다음 commit에 회수한다 | 권한·부분 작업 복원 보존 |
| 현재 자격·lease | public 양측 member/owner/key·active pair·세대 확인, text와 frozen/test 구분, persist 전 ACK 거부, TTL/30s lease·3회 제한 | 보존 |
| gate·CSRF | same-origin·현재 owner·CSRF·서명/body/digest/세대/기한을 다시 확인한다. hint는 html/template로 escape한다. deny 성공 뒤 canonical303, 정책 executable/disclosure=false | 보존 |
| rate/용량 | 안정된 principal rate를 입장 transaction에서1회 소비한다. 로컬 capacity 거부는 DB 입장 전이다. queue100/gate100/receipt20000/claim4·HTTP16/4·crash30s 유지 | 문서·기존 시험 계약과 일치 |
| 일반 text | 별도 서명 domain·closed wire·4096bytes/180s·현재 인가 후 멱등·한 번 관련 답장·본문 삭제·metadata24h | CLI stdin의 M-UTF8-01 제외 |

정리 호출자는 `boundedHTTP`와 시험뿐이다. `remember` 제품 호출자는 transaction/refresh다. requestBuckets·enterHTTP·cleanupTarget·requestHit와 모든 handler를 검색해 대조했다. ACK/text/test ACK, 키/agent/owner 철회, unpair, cancel, invite/gate deny, logout을 확인했다. snapshot은 인가의 최종 권위가 아니다. transaction과 handler가 현재 DB에서 다시 검증한다.

기존 UI 원본49개와 FIX42개의 manifest hash를 확인했다. PNG를 다시 직접 시각 검수하지 않았다. 이미 제외된 로그·이미지와 추가 raw evidence는 파일별 skipped 이유를 남겼다. 원본334 hash도 모두 일치했다. 과거 pending·quota·플랫폼 차단 결과를 새 PASS로 바꾸지 않았다. [evidence-check.json](evidence-check.json)에 검사·해시·provenance를 남겼다.

## 발견 사항

### M-UTF8-01 — medium, 미해결: CLI stdin 청크 경계에서 UTF-8 본문이 변형된다

- 위치: `adapters/src/text.ts:215–218`, 특히 `input += String(chunk)`.
- 조건: `text.js send`의 stdin Buffer가 멀티바이트 문자 중간에서 나뉜다. stream의 청크 경계는 문자 경계를 보장하지 않는다.
- 영향: 입력 `한`의3bytes가 두 청크로 나뉘면 `���`의9bytes로 바뀐다. sendText는 이 변형된 본문을 검증·서명한다. 상대는 원문과 다른 내용을 정상 메시지로 받거나4096bytes 이내 원문도 길이 초과로 거부된다. MCP 문자열 입력은 이 CLI 청크 변환을 거치지 않는다.
- 검증: Node stdlib의 같은 `String(Buffer)` 연산으로 분할 결과를 assert했다. 네트워크·메일·DB·부하·차단된 리뷰 시험을 실행하지 않았다. [encoding-check.json](encoding-check.json)에 명령과 exit0을 남겼다. 최종 CLI 전체 실행으로 오인하지 않는다.
- 수정 요청: DEV가 기존 stdin에 streaming UTF-8 decoder를 적용하거나 bytes를 상한까지 모은 뒤 한 번 decode한다. byte 상한·명시 승인·비공개 입력을 유지한다. 분할 한글/emoji와4096bytes 경계의 작은 회귀를 추가한다.
- 범위: 현재 요구 안의 좁은 기술 결함이다. 제품 규칙 변경이 아니다. OPS는 제품 코드를 수정하지 않는다. coor가 DEV 후속을 배정하고 새 fixed SHA에 검사·delta 리뷰를 연결한다.

새 critical/high는 찾지 않았다. 원 H-1·H-2의 원 SHA 실패는 보존한다. 위 표와 저장된 회귀는 d089에서 해당 원인 수정의 정적 근거다. 새 독립 실행 QA의 통과로 대신 쓰지 않는다.

기존 L-1(전체 JSONB·전역 row lock·자격 색인 비용), L-A(합성 경로의 ID 재등록), L-B(legacy invite-decision의 세대 미결속)는 해소하지 않았다. refresh는 포화 시 전체 JSONB 읽기를 더한다. 다수 owner 공모·pool 대기·운영 크기/CPU·실부하/edge 보호는 공개 수락 전 OPS 담당이다. 합성 가입 unset·운영 DB 합성 owner0을 유지해야 한다. 이 low/운영 후속을 새 PASS로 표시하지 않는다.

## 검증 및 남은 제약

후보 검사 재사용: d089의 `product-lint`/`product-test` 각exit0을 exact JSON에서 확인했다. 마지막 코드5d1924c와 d089의 제품 경로 diff는0이다. DEV-FIX3 최종 lint/test/MVP의 .exit는0이다. MVP는 PASS47·FAIL0이며 stale/restart·owner/agent 두 단위·flood 0/18·finish 회수를 포함한다. 이는 DEV 실행의 재사용이다. reviewer/tester의 새 실행으로 표시하지 않는다.

원758 RED의 stale/restart429·동일시각 snapshot 역전·정리 슬롯1/2, epoch/persist/local/shared 변형 검출, 중간9620723 MVP exit2도 보존했다. 신규 여유 때 불필요한 refresh가 pool을 기다리던 중간 실패와 최종 포화 때만 refresh하는 수정이 코드·기록과 일치한다. 검증 결함 때문에 공격/부하 도구를 새로 만들거나 차단된 원 시험을 재실행할 필요는 없었다.

경고14개: SIZE-001 10개는 기존 큰 handler/시험/누적 PLANS이며 수정 경계와 원 기록 보존 이유를 확인했다. SEC-001 두 개는 시험 fixture `LeaseToken="lease"`와 `ClaimToken="active"`다. 실 credential이 아니다. SLOP-004는 plugin 검증의 PASS 출력이다. SIZE-002는 위 누적 범위 설명으로 처리했다. dependency·DESIGN 경고는 없다.

실메일·공개·외부 계정/노우↔다닷·실24h·새 운영 측정·직접 시각 검수는 실행하지 않았다. 최신 fixed의 tester 독립 QA와 designer 영향 판정은 coor가 별도로 확인한다. M-UTF8-01은 수정 요청으로 남는다. check는 기록·독립성·exact SHA·high 차단 검사이며 기능/공개 수락이 아니다.

## 대화 미참조 인계 점검

산출물 인덱스 `docs/deliverables/README.md`에서 D01–13을 찾았다. 실제 source와 변경 영향을 대조했다. [source-documents-check.json](source-documents-check.json)은 로컬 링크/파일 존재와 읽기 범위를 기록한다. 외부 링크를 새로 조회하거나 신규 제공자 사실을 검증하지 않았다.

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | D01 business-plan; D02 SAR-PUBLIC-SERVICE PS08–11; D04 일반서비스 UX·MESSAGES UI/FIX | 확인 | 일반 신원·명시 text·한도·gate와 무료/후속 경계가 있다. 기획 규칙을 변경하지 않았다 |
| 구조와 구현/미완료 상태 | D03 architecture/tech-stack; D05 interface; D06 data-model; D07 database; D09 CRUD; D10 module | 확인 | 공통 admission·ReadRelay·epoch·owner/agent 단위와 DEV-FIX3가 연결된다. 원09c QA 절은 원 시점 기록이다 |
| 실제 schema | D08 generated/db-schema; db/migrations/00001_relay.sql | 확인 | table/migration 무변경. ReadRelay는 질의 추가다. D08 재생성은 해당 없음 |
| 실행·검증 방법과 증거 | project·Makefile·DEV/FIX/FIX2/FIX3 phase·후보 lint·TESTER09c·UI dfc | 확인 | 정확한 후보와 역사적 실행을 분리한다. 새 독립 QA/UI가 남았다 |
| 사용자 사용 방법 | D11 user-guide; 루트 README; adapters README·KnowsLink skill | 확인/보완 필요 | D11은 기존 베타/신원 안내다. 새 text 절차는 README/adapter에 있다. 공개 전 OPS가 D11에서 해당 정본을 연결해야 한다 |
| 운영·복구 | D12 ops-guide·D13 transition; SAR-PUBLIC-SERVICE-OPEN-PREP | 확인 | 기존0911 배포의 백업/격리복원·보호 유지 근거다. 새 후보 공개/실메일/일일백업·RPO/RTO 수락이 아니다 |
| 다음 작업·담당·재개 조건 | PLANS·DEV-FIX3 미검증/인계·OPEN-PREP | 확인 | coor가 M-UTF8-01 DEV 수정·새 SHA 검사와 독립 QA/UI 뒤 main 판정한다. SMTP/권한/공개는 별도 OPS 후속이다 |
| 로컬 링크·절 접근/지원 한계 | source-documents-check.json; packet-outcomes.json | 확인 | 링크 파일을 검사했다. 자동 anchor/render 검증과 외부 URL 재조회는 하지 않았다. D02 PS08–11/UX06–07은 절/표로 직접 찾았다 |
| snapshot 정리 후 정본 접근 | result/report/lint/coverage/evidence-check; Git fixed SHA | 확인 | 정본 기록은 OPS checkout에 있다. snapshot cleanup은 reviewer release 뒤 coor가 managed 절차로 수행한다. 이 작업은 삭제하지 않았다 |

## 검토 결론

**수정 요청이다.** 새 critical/high0이지만 CLI 원문 보존의 M-UTF8-01 medium이 미해결이다. d089의 무조건 main/공개 수락을 승인하지 않는다. DEV 수정 또는 coor가 책임자와 합의한 명시적 수락 판단이 필요하다. 기존 요구를 낮추는 예외는 이 리뷰가 만들지 않는다.

본 과제의 독립 정적 검토와 기록 작성은 완료했다. 결과·packet outcomes·exact base/head `review.py check`·archive·마지막 역할 커밋 검사·일반 push를 완료 근거로 연결한다. 제품 QA·직접 시각 수락·새 코드 수정·main 통합·공개는 해당 담당 후속이다. 준비 handover/packet/PLANS와 report 템플릿 원본은 [prepared-inputs.json](prepared-inputs.json), 별도 template 원문, 작업 전문 archive로 보존했다.
