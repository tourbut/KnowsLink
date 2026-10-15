---
title: 요청 상태 조회 DEV 결과
status: draft
updated: 2026-10-15
owner: dev
tasks: [SAR-UX-ISSUES-002-RECEIPT-DEV]
summary: 요청 조회 복구의 DB·인증 검사와 합성 fixture 및 독립 브라우저 검수 인계
---

# 요청 상태 조회 DEV 결과

제품 source는 `1e49033c6a036112019474fd010015a31996b7eb`다. 기준 ref는 `e939cfbb1939adb4126ad379317c26ec39b692b3`이며 착수 HEAD는 `d4a075c9ce035379013fa1a397f524e5cc54762a`다. 이후 완료 문서 커밋은 제품 소스를 바꾸지 않는다. 최종 HEAD와 FullOps gate의 실제 결과는 worker_done 및 외부 `receipt-ux002-dev/final-lint.json`에 연결한다.

## 변경과 경계

최초·성공·형식 오류·없는/타인/보존 종료·429에 같은 GET 편집 폼을 제공한다. 안전한 ID 입력과 소유 agent 선택을 유지한다. 128 UTF-8 bytes를 넘는 입력이나 제어/방향 문자는 해당 ID만 비우고 선택을 유지한다. 소유하지 않은 agent를 담은 이전 URL은 그 agent와 ID를 반사하지 않는다. 철회된 선택은 다른 agent로 자동 변경하지 않는다.

현재 소유 agent의 이름·종류·안정 구분값만 선택기에 표시한다. 구분값이 충돌하면 전체 ID로 확장한다. 결과는 전달·처리·관련 답장을 구분하고 원래 상태·ID는 기술 상세에 둔다. 새 조회는 발송·ACK·원문 공개·권한을 만들지 않는다. 기존 24시간 보존·서명·친구·멱등성 계약은 유지한다.

`__Host-kl_receipt`는 마지막 조회의 회원·agent·안전한 원입력·만료만 담는 입력 쿠키다. Secure/HttpOnly/Strict이고 최대 70분이며 조회 때 갱신한다. 서명된 인증 수단이 아니다. 서버가 로그인 회원·소유권·길이·만료를 다시 검사한다. 명시 로그아웃은 지우고 다른 회원·익명에게 이전 계정 내용을 복원하지 않는다. 결과·이메일·이름·키·세션 token은 저장하지 않는다. 이메일과 Google의 기존 로그인 검증을 유지하며 후자는 attempt의 복귀 표시와 동일 출처 링크를 사용한다. device/remote 복귀가 우선한다.

상위 HTTP admission의 429도 해당 transaction에서 인증한 폼을 사용한다. 제한을 우회하거나 추가 admission slot을 열지 않는다. DB가 로그인 자체를 확인할 수 없는 503 또는 로컬 동시 처리 한도는 이전 입력을 표시하지 않고 쿠키와 URL을 유지한다. 서비스 복구 후 다시 인증하고 복원한다. DB commit 실패처럼 이미 인증한 view가 있으면 입력을 유지한다. 권한을 확인할 수 없는 캐시로 사적 입력을 노출하지 않는 fail-closed 경계다.

D05/D10/D11 정본을 갱신했다. D09 메시지 CRUD·DB 테이블·영속 조회 이력은 변하지 않는다. GoogleAttempt의 짧은 복귀 bool만 D05/D10에 기록했다. 기존 디자인의 여섯 토큰·44px·글꼴·간격·CSP를 그대로 사용한다. 새 dependency·CSS·JS·언어 기능·승인 요약·설치 온보딩은 없다.

## 실행 근거와 원 실패

- source1e49033 관련 Go 검사: `go test ./internal/relay -run 'TestReceipt|TestMemberReceiptView|TestMemberUIFixtures|TestHuman|TestGoogle' -count=1`, 종료0. [최종 로그](final-related-tests.log).
- source1e49033 실제 Linux/Postgres race 검사: `go test -race -tags=integration ./internal/relay -run 'TestReceiptRecoveryHTTP|TestEmailIdentity|TestGoogleHTTP|TestGoogleDeviceHTTP|TestPublicMessageHTTPBoundaries|TestPublicHTTPAdmissionAndGateSafety|TestHumanPostgres' -count=1 -v`, 19.711초·종료0. [로그](1e49033c6a036112019474fd010015a31996b7eb-integration.log).
- coverage: 성공/빈/잘못된/없는/타인/보존 종료·만료 receipt, 429 입력 보존, 세션 만료·동일회원 재로그인·다른회원 교대·익명 비노출, DB 장애 뒤 복구, logout 초기화, 동명 구분값 충돌, HTML escape와 단일 autofocus 대상.
- 첫 form golden 검사는 기존 hidden 재조회와 legacy 인라인 폼의 의도된 변경 때문에 실패했다. 소유 목록 fixture를 보강하고 승인된 두 계약만 golden에 반영했다. 다른 form 계약은 그대로다.
- sourcecc84b24 SQL 검사는 테스트가 B의 세션을 A 브라우저에 공유한 뒤 로그아웃하고 B로 조회해 303을 받은 것으로 실패했다. 제품의 세션 종료 동작은 정상이었다. 보존기간 검사를 logout 이전으로 옮긴0e4eba4 및 최종1e49033에서 통과했다. [원 실패](cc84b24e11b92414a9d2c65e342f378abeaba2ea-integration.log)를 보존한다.
- product-lint/product-test는 최종 완료 기록 커밋에 대해 정본 lint.py가 등록된 make lint/test를 실행한다. [helper](remote.py)는 exact HEAD의 clean Linux clone·known_hosts RejectPolicy·원 SSH 종료코드를 유지한다. final-lint.json과 HEAD별 lint/test 원로그는 private 경로에 보존한다. 이전 제품 PASS를 이번 실행으로 재명명하지 않는다.

## QA fixture 인계

`http://127.0.0.1:18928`은 task-owned Windows SSH forwarding origin이다. 실제 fixture는 source1e49033의 `/tmp/knowslink-receipt002-1e49033c6a03-a0fafa26`에 있고 서버 loopback `http://127.0.0.1:44077`을 사용한다. 합성 DB container는 `kl-receipt002-fixture-1e49033c6a03-76b5ef22`다. [manifest](browser-fixture.json), [시작 로그](fixture.log), [익명 접근 smoke](fixture-smoke.json)를 보존했다. 서버와 forwarding은 시작 후24시간 동안만 유지하며 종료·재준비는 coor 후속이다.

생성은 `remote.py prepare` 뒤 `remote.py fixture`다. 내부 명령은 `go test -c -tags=integration,receiptfixture`와 `RECEIPT_FIXTURE_MANIFEST=<private path> <binary> -test.run=^TestReceiptBrowserFixture$ -test.timeout=25h -test.v`다. `TEST_SYNTHETIC_DATABASE=1`과 고유 Postgres DB가 필수다. 기존 production DB에는 실행하지 않는다. [forward.py](forward.py)는 manifest의 SHA·task 경로·loopback target을 검사한다. SSH 값은 기존 비공개 파일에서 읽으며 보고·argv·로그에 출력하지 않는다.

허용 요청은 GET `/`, `/home/receipts`, `/home/agents`, `/home/legacy`와 아래 합성 POST뿐이다. fixture의 다른 경로·메시지 발송·외부 인증·실메일은 서버가 거부한다. 같은 origin의 요청으로 실행하고 `Origin: http://127.0.0.1:18928`을 사용한다. 포트가 달라도 쿠키는 공유되므로 독립 브라우저 profile/context를 써야 한다.

| POST | 합성 효과 |
|---|---|
| `/fixture/as/a`, `/fixture/as/b`, `/fixture/as/c` | A/B/C의 새 유효 세션을 발급한다. 실제 이메일/OIDC 검증을 대체하는 browser fixture 전용 주입이며 기존 조회 쿠키는 유지한다. |
| `/fixture/expire` | 현재 합성 세션의 idle 만료를 만든다. 다음 GET은 로그인 화면으로 간다. |
| `/fixture/limit` | 현재 합성 회원의40회 rate를 채운다. 다음 조회는429다. |
| `/fixture/outage` | product handler의 DB context를 취소하여503을 만든다. |
| `/fixture/recover` | 위 outage를 해제하고 합성 rate를 비운다. |
| `/fixture/revoke` | A의 agent_a를 철회한다. 마지막에 실행한다. 새 fixture 없이는 복원하지 않는다. |

A는 agent_a와 agent_a2, B는 agent_b, C는 agent_c를 소유한다. 모두 같은39자 한글 이름과 같은 종류로 표시하여 구분값을 확인할 수 있다. A/B 참여 ID는 `0199a3f2-4c10-7a11-8b22-334455667788`, B/C만 참여한 ID는 `0199a3f2-4c10-7a11-8b22-3344556677aa`, 없는 ID는 `0199a3f2-4c10-7a11-8b22-334455667799`다. 수신 완료 합성 metadata이며 실제 메시지를 발송하지 않는다. source의 실제 이메일·Google/device 검증은 별도 SQL 회귀 로그가 증거다.

A로 시작하여 `/home/receipts?agent=agent_a&id=0199a3f2-4c10-7a11-8b22-334455667788`을 연다. ID 수정·다시 조회·뒤로/새로고침·입력 복귀를 확인한다. 존재하지 않는/타인 ID의 동일 문구와 소유 선택 보존을 확인한다. limit/recover, expire/as/a, outage/recover, as/b로 계정 교대를 확인한다. 마지막 revoke에서는 다른 agent로 조용히 선택이 바뀌지 않아야 한다. POST 자체는204이므로 다음 GET을 명시적으로 한다.

정적 합성 HTML은 private `fixtures/`에 있고 [파일 hash](fixture-manifest.json)에 source1e49033을 기록했다. 기존 Go test의 `KNOWSLINK_UI_FIXTURE_DIR` 출력이며 실제 인증/DB 증거와 구분한다. first/invalid/unavailable/limited/success/revoked 상태를 제공한다.

## 미실행과 coordinator 인계

추가 headless Chrome·정적 서버·문서 다운로드·CDP 검사 생성의 묶음 exec 요청이 자동 승인 검토에서 `blocked by policy`로 거부됐다. 상세 사유는 제공되지 않았고 요청은 실행되지 않았다. 원 오류는 도구 기록에 보존한다. escalation `msg_86ac8f6f320b`, coordinator 답변 `msg_8a42617af54d`는 재시도/다른 셸 우회 금지와 기존 DEV 완료 조건 유지를 명시했다. 추가 headless 기동은 DEV 필수가 아니며 독립 tester가 기본 Orca 브라우저로 검증한다. 후속 승인된 fixture 준비는 브라우저를 기동하지 않는다.

따라서320/390/1440 실제 렌더, 실제200% 확대, 키보드·스크린리더·보조 브라우저는 이번 DEV에서 not_run이다. template와 autofocus 구조 검사는 실제 포커스·시각 수락을 대신하지 않는다. 제품 수락 기준은 낮추지 않는다. designer 직접 시각 검수, tester 독립 QA, 작성자와 다른 세션의 고정 SHA 코드 리뷰와 main 통합은 coor 책임이다.

Context7 Go resolve는 monthly quota exceeded, web 도구는 token_expired여서 외부 문서 조회를 완료하지 못했다. 기존 표준 라이브러리·setCookie/html/template 패턴을 재사용했으며 API·SDK dependency 변경은 없다. frontend.py configuration_ready는 인증/브라우저/MCP PASS가 아니다.

제공된 same-key route/find/context/packet/history 원본을 커밋했다. [source-audit](source-audit.json)는 패킷 비선택 후보와 unknown/잔여 경로의 hash·관련절 검색 근거다. 큰 D04는 현행 토큰·접근성·U6 연결을 직접 읽었다. 과거 인박스·리뷰·OPS 기록을 현재 실행 지시나 이번 PASS로 사용하지 않았다. 패킷에 빠진 실제 영향 경로 capacity.go/member.go/google.go/identity.go와 새 테스트를 직접 확인했다. 기술 계획은 완료 보고 안에 보존하여 최초 지시 identity를 유지한다. packet outcomes와 work.py finish로 원 인박스와 결과 전문을 보존한다.

자료 보존 루트는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/receipt-ux002-dev/`다. 기존 파일 삭제·이동, 운영 변경·배포·실계정·실메일·상대 발송·GitHub·PLANS·board 변경은 하지 않았다.

검증 로그 보존: 첫 SQL 실패 출력의 공백만 있는 줄이 git diff --check에 걸렸다. 표시용 .log의 줄 끝 공백만 정리했다. 원 출력 전체와 SHA256은 initial-integration-raw.json 및 private 원로그에 그대로 보존했다. 테스트의 실패·종료코드는 바꾸지 않았다. 산출물 strict 검사는13종·문제0·경고0이었다.
