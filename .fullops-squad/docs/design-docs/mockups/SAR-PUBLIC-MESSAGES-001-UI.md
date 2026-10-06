---
title: SAR-PUBLIC-MESSAGES-001-UI — 일반회원 receipt·gate 직접 검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI]
summary: 고정 일반회원 receipt·gate 직접 시각 판정과 Deny 후405 결함·원본 증거를 보존한다
---

# SAR-PUBLIC-MESSAGES-001-UI — 일반회원 receipt·gate 직접 검수

## 판정과 고정 대상

**UX06은 아래 로컬 시각 조건에서 PASS다. UX07은 F-UI-MSG-01 medium으로 FAIL이며 시각 수락을 보류한다.**
제품 고정 SHA와 검사 기준 ref는 `09c523da8a3407288d9f5d711e1834af12bc7808`이다.
기록 시작 HEAD는 `a7443d5`다. 제품 코드·기술 정본·제품/UX 규칙은 변경하지 않았다.
worker의 검수 기록 완료와 제품 수락은 별개다. OPS가 같은 후보에서 발견한 H-1 high·M-1 medium도 유지한다.
이 보고서는 DEV의 자동검사·390px 측정을 직접 시각 PASS로 바꾸지 않는다.

적용 기준은 `fullops-common-0.3.3`의 네 문서, FULLOPS·project·document-writing·designer context다.
제품 근거는 [PS04/06/07·PS08–11](../../planning/product-specs/SAR-PUBLIC-SERVICE.md), [MVP C1–C5](../../planning/product-specs/SAR-MVP.md), [UX06/07](SAR-PUBLIC-SERVICE-UX.md)다.
DEV 실행 기록·완료 로그·QA 원본·mobile-width.py·README·adapters README와 실제 `member_receipt.go`·member/http/public_text/capacity/store·text.ts를 읽었다.
Jev의 `member_receipts.go` 오기와 c6 준비 실패는 원본 JSON에 보존했다. 필수 adapters README와 capacity/store를 제외하지 않았다.
[정본 동일성](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/source-identity.json)은 읽은 체크아웃과 고정 후보의 내용 일치를 확인한다.

## 실제 화면과 자기 fixture

별도 detached clone의 변경 없는 relay를 빌드했다. 자기 Postgres 17 컨테이너 `sar-messages-ui-28d645895e`와 기존 로컬 SMTP sink를 사용했다.
Chrome의 서로 다른 비영속 context에서 두 합성 일반회원이 실제 `/auth/start`·코드 입력·`/auth/verify`로 로그인했다.
브라우저가 받은 Secure·HttpOnly 회원 cookie로 실제 `/home/receipts`와 `/home/gates`에 접근했다. owner Basic 인증을 사용하지 않았다.
회원·세션은 실제 로그인 결과다. 검수용 agent의 키와 활성 관계는 자기 SQL fixture로 준비했다. 일반 온보딩 재검수로 보고하지 않는다.
서명 text 송신·pull·persist·ACK·관련 답장 및 서명된 업무 부모·claim·H 생성은 변경 없는 실제 HTTP handler로 처리했다.
시간·철회·원문 부재·포화 조건은 자기 relay를 중지하고 자기 DB만 수정한 뒤 실제 sweep·검증으로 관찰했다.

실제 주소는 `http://localhost:46743`이었다. 데스크톱은 1280×900, 모바일은 390×844다.
agent 식별자는 각각 128자이며 요청/답장 ID는 36자다. 긴 ID의 줄바꿈·상태·버튼·오류를 PNG에서 직접 확인했다.
[manifest](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/manifest.json)는 SHA·조건·viewport·이미지 hash·직접 판정을 연결한다.
PNG는 49개다. 유효 화면 47개와 초기 만료 fixture 오류 보존 2개다. 모든 유효 화면을 직접 열었다.
영상은 필요하지 않았다. 이번 판단은 정지 화면의 상태와 실제 버튼 결과로 가능했다.

## UX06 직접 판정

| 조건 | 관찰과 판정 | PNG 이름 접두사 |
|---|---|---|
| 선택 agent·명시 송신 안내 | PASS. 홈은 자기 agent 선택·요청 ID 조회와 Node client 명시 send/reply·비공개 stdin·manual receive를 안내한다. composer·채팅버블·장기 목록은 없다. | `home-tool-guidance` |
| queued·leased·실제 수신 | PASS. queued는 수신 미확인, leased는 ACK 없음, delivered는 클라이언트 저장·ACK 확인이다. 처리 결과 없는 상태와 received를 구분한다. | `receipt-queued`, `receipt-leased`, `receipt-delivered` |
| 관련 답장 | PASS. 원요청 화면은 별도 답장 ID와 답장의 queued/delivered를 보인다. 답장 ACK 후 reply_received다. 답장 receipt는 원요청 ID를 보인다. | `receipt-reply-queued`, `receipt-reply-delivered`, `reply-parent-id` |
| TTL·수동 확인 | PASS. KST 기한·TTL·text 최대180초·idle pull10초 이상·자동 wake 없음이 읽힌다. 실제 수동 확인 버튼이 같은 agent/요청을 조회한다. | `receipt-manual-refresh` 및 각 receipt |
| 오프라인·만료·철회·수신 소진 | PASS. 오프라인 가능성과 pull 미실행을 단정하지 않고 만료 원인으로 안내한다. 새 명시 요청·현재 키/관계 확인·옛 요청 복구 불가·클라이언트 확인을 보인다. | `receipt-expired`, `receipt-revoked`, `receipt-max-attempts` |
| 제한·잘못된 키·conflict | PASS, 표시 범위를 한정한다. 실제 429 화면은 KST 재시도 시점과 홈 복귀를 제공한다. 잘못된 키·conflict의 다음 동작은 receipt의 상시 안내다. 실제 API 401 invalid_auth·409 idempotency_conflict도 기록했다. 별도 웹 송신 오류 화면은 존재하지 않으므로 PASS를 주장하지 않는다. | `receipt-rate` 및 receipt 안내 |
| 타 회원·긴 ID·키보드 | PASS. 다른 회원은 요청 내용 없이 안전한 거부·홈 복귀를 본다. 128자 ID와36자 요청이 화면 안에서 줄바꿈한다. Tab 순서는 수동 확인·홈 복귀다. | `foreign-receipt-refusal` 및 각 receipt |

## UX07 직접 판정과 결함

| 조건 | 관찰과 판정 | PNG 이름 접두사 |
|---|---|---|
| 검증 본문·판단 근거 | PASS. 발신/대상·원요청·schedule.query·정책·KST 기한·검증된 window/granularity_min JSON을 보인다. hint에는 판단 근거가 아님을 명시한다. | `gate-pending` |
| HTML·명령문 | PASS, 로컬 관측 한정. 다른 내용의 hint와 script/명령 문자열은 글자로 보인다. script 요소0·injected=false다. 본문을 숨기지 않는다. 셸·업무 효과는 실행하지 않았다. | `gate-pending`, `gate-approved`, `gate-denied` |
| approve·deny 결과 | Approve는 PASS. approved와 비활성 결정 영역을 보인다. Deny의 실제 직후 화면은 FAIL이다. canonical gate를 직접 열면 denied와 비활성 버튼이 보인다. | `gate-approved`, `gate-deny-redirect`, `gate-denied` |
| 원문 부재·만료·철회·잘못된 서명 | PASS. unavailable/expired/revoked 또는 원문 부재 문구와 승인 불가가 보인다. approve 버튼이 없다. 서명 실패 fixture의 상태명은 pending으로 남지만 본문과 결정은 비활성이다. | `gate-unavailable`, `gate-expired`, `gate-revoked`, `gate-invalid-signature` |
| 권한·오류·복구 | PASS, Deny 직후 오류는 제외한다. 타 회원은 자료 없이 안전한 거부를 본다. CSRF 실패는 승인되지 않음·홈 복귀다. 동시 한도는1초 뒤 수동 재시도·성공 아님·홈 복귀다. | `foreign-gate-refusal`, `gate-csrf-error`, `gate-capacity-error` |
| 모바일·키보드 | PASS, Deny 결과 결함은 유지한다. 390px에서 typed body·ID·오류·버튼이 읽힌다. Tab은 Approve·Deny·홈 복귀이며 Enter로 Deny해도 같은405가 재현된다. | 각 gate 및 `supplement-result.json` |

### F-UI-MSG-01 — medium: Deny 뒤 결과 대신405

일반회원이 정상 pending gate의 **Deny 거절**을 누르면 `/home/gates/<id>/deny`로 이동한다.
브라우저 GET 응답은 **405**, 화면에는 `Method Not Allowed`만 있다. 거절 결과·요청 정보·홈 복귀 링크가 없다.
DB의 gate는 이미 `denied`다. canonical `/home/gates/<id>`를 직접 열면 denied 결과와 홈 복귀가 정상이다.
따라서 거절 자체 실패나 권한 우회는 아니다. 사용자가 거절 완료를 확인하고 다음 동작을 선택하는 UX07 완료 조건이 실패한다.
마우스 클릭의 실제 오류를 desktop/mobile PNG로 보존했다. 키보드 Enter의405·저장된 denied는 [후속 원관측](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/supplement-result.json)에 있다.
DEV가 기존 요구 안에서 수정 방법을 결정한다. 새 고정 후보에서 designer가 Deny 직후 결과·홈 복귀를 직접 재검수한다.
제품 규칙·API 설계·코드는 이 과제에서 수정하지 않았다.

## 초기 실패·보존·회수

임시 Go fixture 작성은 설계 역할 코드 금지 hook에 의해 차단됐다. 기존 relay와 허가된 QA 실행 기록·SQL fixture 방식으로 진행했다.
첫 브라우저 실행부터 다섯 번째까지는 각각 exit1이다. 종료된 fixture process, MIME 코드 추출, email 변수 충돌, H의 reason 필드 누락, SMTP 재시작 bind 실패가 원인이다.
여섯 번째 실행은 exit0이다. 초기 실패 로그·exit와 SMTP 장애 원문을 모두 보존했다. 제품 결함으로 바꾸지 않았다.
초기 만료 fixture는 exp만 과거로 옮겨 음수 TTL을 만들었다. 이2개 PNG를 `initial-fixture`로 보존했다.
accepted_at도181초 전으로 옮긴 새 실제 화면은 TTL180이다. 수정된 만료 fixture만 UX06 판정에 사용했다.
첫 supplement는 navigation status assertion으로 exit1이다. 실제 Network 응답·Enter 입력을 기록한 후속은 exit0이다.

로그는 [자기 실행 증거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/)에 있다. 명령 자신의 종료코드를 각 `.exit`에 보존했다.
DEV 원본45파일은 [SHA256 대조](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/original-dev-preserved.json)에서 동일하다.
[cleanup](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI/cleanup.json)은 자기 DB 컨테이너/익명 volume·relay·SMTP·Chrome context/profile·메일·비밀 설정·clone을 회수한 결과다.
공유8개 컨테이너의 ID·이름은 자기 정리 전후 동일하다. 기존 .env·운영 DB·서버·Tunnel·외부 계정은 사용하지 않았다.
이메일·code·연결수단·CSRF 입력은 캡처 전에 가렸다. 클라이언트 private key·credential·메일 code는 로그와 Git에 남기지 않았다.

## 검사·미검증·인계

기존 Go template/memberStyle을 기준으로 판정했다. 별도 theme/design lint·Tailwind/shadcn은 해당 없음이다.
FullOps 자동 DESIGN 결과를 직접 시각 수락으로 해석하지 않는다. 모든 유효390px 화면은 scrollWidth390이며 직접 이미지 판정도 별도로 수행했다.
고정 기록 HEAD의 FullOps `--from 09c523d`·등록 product-lint/product-test·strict·diff 검사와 경고는 [실행 기록](../../exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI.md)에 연결한다.
마지막 archive/metadata SHA도 필수검사를 실행한다. 자기 SHA를 자기 커밋에 순환 기록하지 않고 worker_done의 fullSHA와 원본 검증 경로로 고정한다.

실메일·운영 PS08/13/14·운영 공개/배포·실24h·외부 Grok Bot/다닷·노우↔다닷·운영 부하/복원은 미검증이다.
실제 사용자 이메일의 사람 로그인·스크린리더 전체·OS clipboard·schedule.commit 개별 화면·전체 API 권한/자원 QA도 이 직접 검수의 PASS 범위가 아니다.
TESTER 독립 QA·OPS 보안/운영 수락은 별도다. coor가 원본 high/medium과 새 F-UI-MSG-01의 DEV 수정·새 fixed 재검수·통합을 조정한다.
수정된 화면의 PASS로 원본09c FAIL·OPS H-1/M-1·초기 실패를 덮어쓰지 않는다. 다음 과제는 PLANS에 대기시킨다.
