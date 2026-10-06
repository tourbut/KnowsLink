---
title: SAR-PUBLIC-MESSAGES-001-UI-FIX — 수정 gate 직접 재검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI-FIX]
summary: 수정 gate Deny 결과·홈 복귀의 직접 판정과 HTTP 포화 한계를 보존한다
---

# SAR-PUBLIC-MESSAGES-001-UI-FIX — 수정 gate 직접 재검수

## 판정과 고정 대상

**F-UI-MSG-01의 정상 pending gate Deny 직후 결과·홈 복귀는 로컬 직접 시각 PASS다.**
제품 대상과 검사 기준 ref는 `dfc70caa748a90614b02d48c78b4651345938339`다.
실제 일반회원의 desktop1280/mobile390 클릭·Enter 네 조건에서 canonical 결과 GET200·denied 저장을 확인했다.
결정 버튼은 사라지고 비활성 문구가 보인다. 결과 링크를 실제 클릭하면 `/home`의 denied 목록으로 돌아간다.
owner 경로도 자기 자격의 좁은 desktop 조건에서 denied 결과 GET200·비활성·`/owner` 복귀를 확인했다.

이 PASS는 신규 HTTP rate 포화에서 결과 GET까지 항상200이라는 뜻이 아니다.
아래 L-UI-FIX-1에서는 Deny가 저장됐으나 후속 조회가429다. 이 조건은 즉시 결과 PASS에서 제외한다.
원본 [09c 직접 검수](SAR-PUBLIC-MESSAGES-001-UI.md)의 UX07 FAIL·F-UI-MSG-01 medium은 그 원래 SHA의 결과로 유지한다.
원본 mobile Deny 오류의 effectiveviewport980과 [R-UI-1](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md)도 보존한다.
원 OPS H-1 high·M-1 medium·review check exit1을 이 UI 과제로 해소하지 않는다.
제품 전체 수락·main 통합은 coor가 독립 TESTER/OPS 결과와 함께 판단한다.

## 기준과 독립 실행

`fullops-common-0.3.3` 네 규칙과 FULLOPS·project·document-writing·designer context를 적용했다.
제품 근거는 [PS04/06/07/08–11](../../planning/product-specs/SAR-PUBLIC-SERVICE.md), [MVP C1–C5](../../planning/product-specs/SAR-MVP.md), [UX06/07](SAR-PUBLIC-SERVICE-UX.md)다.
DEV-FIX phase·QA·완료 전문과 원UI report·manifest·supplement·UI-RECORD-REVIEW를 직접 읽었다.
Jev의 cleanup_admission 미전송은 제외 근거가 아니다. capacity·cleanup_admission·http·store·member·member_receipt의 관련 경계도 직접 읽었다.
제품 코드·제품/UX 규칙·기술 정본·타인 인박스는 바꾸지 않았다.

기록 시작 HEAD는 `359653a`다. 검수는 별도 clean detached clone `/tmp/sar-messages-ui-fix-vnsxs7mb/snapshot`에서 실행했다.
[고정 snapshot](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/snapshot-identity.json)은 exact SHA·detached·제품 diff0을 기록한다.
변경 없는 relay/migrate를 빌드하고 자기 Postgres17 컨테이너와 localhost SMTP sink에만 연결했다.
실제 주소는 `http://localhost:39033`이었다. 두 비영속 Chrome context에서 서로 다른 합성 일반회원이 실제 이메일 코드로 로그인했다.
`/auth/start`·`/auth/verify`로 받은 Secure·HttpOnly cookie를 사용했다. 일반회원 경로에 owner Basic 인증을 쓰지 않았다.
agent·키·활성 관계는 자기 SQL fixture다. 업무 부모·pull·persist·ACK·claim·H는 실제 서명 HTTP로 만들었다.
agent ID는 각각128자, 요청/gate ID는36자다. typed body는 window/granularity_min의 비민감 합성 값이다.
시간·철회·원문 부재·한도 조건은 자기 relay 중지 후 자기 DB만 변경하고 실제 handler/sweep으로 관찰했다.
owner 회귀만 자기 owner credential을 별도로 만들고 `/owner` context에서 Basic 인증했다.

## 직접 관찰

| 조건 | 관찰과 판정 | PNG 접두사 |
|---|---|---|
| Deny 클릭·Enter | PASS. 두 폭에서 즉시 canonical GET200·denied 저장·비활성 문구·실제 홈 클릭을 확인했다. GET405는 없다. | `gate-before-deny-*`, `gate-deny-*-immediate`, `home-after-deny-*` |
| owner Deny | PASS, desktop만. canonical `/owner/gates/<id>` GET200·denied·비활성·실제 `/owner` 복귀다. | `owner-gate-*` |
| pending·approve | PASS. 검증 typed body·정책·KST 기한·hint 경고와 결정 버튼이 보인다. Approve 뒤 approved·비활성을 확인했다. | `gate-pending`, `gate-approved` |
| 긴 ID·390px·키보드 | PASS. 실제 viewport/scrollWidth390이며 ID·JSON·버튼·오류가 폭 안에서 읽힌다. Tab은 Approve·Deny·홈 링크다. Enter는 실제 CDP keyDown/keyUp이다. | 모든 mobile 및 `observations.json` |
| HTML·명령문 hint | PASS, 화면 관측 한정. script/명령 문자열은 글자다. script요소0·injected=false이며 본문이 판단 근거다. 외부 업무를 실행하지 않았다. | pending·결정 결과 |
| expired·revoked·unavailable | PASS. 상태와 원문 부재/검증 본문을 구분하고 결정 버튼이 없다. 홈 링크가 읽힌다. | `gate-expired`, `gate-revoked`, `gate-unavailable` |
| 잘못된 부모 서명 | PASS, 승인 불가 기준만. 상태는 pending으로 남지만 원문 부재와 비활성이 보인다. | `gate-invalid-signature` |
| 타 회원·잘못된 CSRF Deny | PASS, 자기 fixture 화면 범위.403 안전 오류·내용 비노출·홈 링크다. CSRF 실패의 저장 상태는 pending이다. | `foreign-gate-refusal`, `gate-csrf-error` |
| 신규16 포화의 Approve | PASS, 거부 화면.429·1초 뒤 수동 재시도·성공 아님·홈 링크다. gate는 pending이다. | `gate-new-capacity-approve` |
| 정리4 포화의 Deny | PASS, 거부 화면.429·1초 뒤 수동 재시도·성공 아님·홈 링크다. gate는 pending이다. | `gate-cleanup-capacity-deny` |
| cleanup rate100 포화의 Deny | PASS, 거부 화면.429·KST 재시도 시점·홈 링크다. gate는 pending이다. | `gate-cleanup-rate-deny` |
| 신규 rate200 포화 뒤 Deny | LIMITED. denied 저장은 성공하지만 canonical GET429로 즉시 상태 확인은 지연된다. 포화 fixture 해제 뒤 실제 홈 클릭·denied GET200을 확인했다. | `gate-new-rate-deny`, `gate-denied-after-rate-recovery` |
| receipt HTTP rate 영향 | PASS, 오류 화면만. 실제 `/home/receipts` GET429·KST 재시도 시점·홈 링크를 보인다. 자기 포화 해제 뒤 링크로 홈에 돌아갔다. | `receipt-http-rate` |

### L-UI-FIX-1 — 신규 rate 포화 뒤 결과 확인 지연

자기 DB의 `http:new` rolling60s 기록을200개로 만들었다. cleanup rate는 남겨 두었다.
유효 Deny 뒤 DB는 denied다. 브라우저는 canonical `/home/gates/<id>`로 갔지만 실제 GET은429다.
화면은 요청 제한·KST 재시도 시점·홈 링크이며 denied 상태는 아직 보이지 않는다.
따라서 저장 실패·405 재발·권한 우회로 기록하지 않는다. 즉시 결과200의 PASS 조건에도 넣지 않는다.
자기 포화를 해제한 뒤 홈 링크를 실제 클릭했다. canonical을 다시 조회하면 denied·비활성·홈 링크가 보인다.
이 관측은 승인된 HTTP rate 정책 아래의 결과 조회 한계다. 모든 포화 조건을 시각 수락했다고 주장하지 않는다.
coor는 TESTER의 cleanup/HTTP 경계 결과와 함께 이 제한을 통합 근거에 연결한다.
새 제품 수치·rate 예외·기술 수정 방법은 이번 직접 검수에서 결정하지 않았다.

## 변경 없는 UX06 증거의 재사용 경계

[원천 동일성](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/source-identity.json)은09c와dfc의 member.go·member_receipt·public_text·identity·member_agents·text.ts·package/lock·Go 의존성·Compose의 동일성을 확인한다.
`internal/relay/text.go`는 양쪽에 없다. 탐색 오기를 absent로 보존하고 동일성 근거에서 제외했다.
receipt 표시와 text wire의 변경 없는 원래 UX06 시각 관측은09c의 SHA·실행 조건·원PNG/manifest로만 연결한다.
queued/leased/delivered·답장·TTL/수동 확인·만료/철회/소진의 원PASS를 dfc의 새 실행 PASS로 바꾸지 않는다.
capacity/store/http의 공통 경계는 변경됐다. 전체 의존성 diff0이나 전체 UX06 동작 동일성을 주장하지 않는다.
이번에는 영향을 받은 gate 오류와 receipt HTTP rate 화면을 실제 새 후보에서 확인했다.
공통 경계의 전체 권한·경합·rate·다중 인스턴스 QA는 TESTER 담당이다.

## 증거·실패·마스킹·회수

[manifest](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/manifest.json)는42개 PNG의 exact SHA·조건·viewport·pixel 크기·hash·직접 판정을 연결한다.
42개를 모두 `view_image`로 직접 열었다. 자동 폭 수치를 직접 시각 PASS로 삼지 않았다.
20개 mobile 화면 모두 실제 viewport/scrollWidth390이다. owner 모바일은 이번에 실행하지 않았다.
영상은 필요하지 않았다. 정지 상태와 실제 버튼·키보드의 최종 응답으로 판정 가능했다.

첫 브라우저 run은18개 관측과40개 PNG를 저장한 뒤 Chrome profile 회수의 `OSError: [Errno 39] Directory not empty: 'Default'`로 exit1이었다.
제품 assertion 실패가 아니다. 이 원래 log/exit를 유지하고 전체 run exit0으로 바꾸지 않았다.
별도 회수에서 자기 Chrome 프로세스 부재·profile 제거·relay/SMTP listener 종료를 확인했다.
첫 회수 확인은 자기 검사 명령의 문자열을 Chrome argv로 잘못 셌다. 판별을 executable 첫 argv로 좁혀 확인했다.
receipt HTTP 후속은 exit0이며 실제 로그인·429·홈 클릭과2개 PNG를 보존한다.

원 UI/DEV-FIX/OPS 원천216개 파일은 [SHA256 대조](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/original-preserved.json)에서 byte 동일이다.
기존 designer archive의 prefix hash도 보존한다. 이번 지시서·완료 전문만 `work.py finish`로 append한다.
이메일·code·CSRF/연결 입력은 캡처 전에 가렸다. 개인키·credential·메일 전문은 Git에 넣지 않았다.
[cleanup](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/cleanup.json)은 자기 container/익명 volume·listener·browser context/profile·메일·비밀 설정·clone·binary 회수를 기록한다.
공유8개 컨테이너의 ID/이름은 전후 동일이다. 운영 DB·서버·Tunnel·실메일·외부 계정은 사용하지 않았다.

## 검사와 인계

기존 Go template/memberStyle·UX 정본을 유지했다. theme/design lint·Tailwind/shadcn은 미구성이다.
제품 구현과 의존성 변경은0이다. DESIGN 자동검사는 직접 시각 수락을 대체하지 않는다.
이번 산출물은 D04 연결·과제 report/QA/phase·context·PLANS·자기 archive다.
최종 archive HEAD에서 FullOps `--from dfc70caa748a90614b02d48c78b4651345938339`의 product-lint/test·strict·ref-to-HEAD diff를 실행한다.
결과 JSON/log/exit·경고·HEAD·원격 동일성은 레포 밖에 남겨 worker_done에 고정한다. 자기 SHA를 자기 커밋에 순환 기록하지 않는다.

실메일·운영 PS08/13/14·운영 공개/배포·실24h·노우↔다닷·부하/복원·전체 API/스크린리더는 미검증이다.
schedule.commit 개별 화면·OS clipboard·실제 사람의 일반 이메일 로그인도 이번 PASS 범위가 아니다.
coor는 fixed dfc의 TESTER/OPS 독립 결과와 마지막 기록 검사·push를 확인한 뒤 main 통합을 판단한다.
원09c 실패·기존 보류를 유지하고 다음 제품 과제는 PLANS에서 기다린다.
