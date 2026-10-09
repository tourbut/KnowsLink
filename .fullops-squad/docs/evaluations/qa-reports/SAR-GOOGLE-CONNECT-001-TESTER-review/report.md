---
title: SAR-GOOGLE-CONNECT-001-TESTER 독립 코드 리뷰
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-GOOGLE-CONNECT-001-TESTER]
summary: "고정 SHA 152217f의 Google 클라이언트 연결 변경 독립 리뷰, 서버 lint 증거 검증, 수락 결론"
---

# SAR-GOOGLE-CONNECT-001-TESTER 리뷰

- 검토자 / CLI / 모델: tester / Claude Code / claude-sonnet-5-5 high. 실제 세션 `b46475f6-dee2-49f6-936f-fe1513a39142`. 구현자 세션 `01a12123-fd98-7b80-880e-6c2802149a6d`와 다르다.
- base SHA / head SHA / merge-base: `f9f7675be6c9502bd6ab2f810bd42ff62a26a174` / `152217f63cced85f620695961955d680ce27864e` / `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`.
- 코드 snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-8493ad160dbd469aa860e952816924c9`. `review.py snapshot`이 만든 clean detached worktree이며 읽기 전용으로 사용했다. 결과·report·lint·QA는 이 tester 기록 체크아웃에만 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13. 규칙은 `.fullops-squad/review/rule.json`(result.json의 rule_sha256). 공통 규칙은 `fullops-common-0.3.3`(`rules/common/README.md`, coding-style, testing, security)과 `project.md`, `docs/agents/document-writing.md`를 직접 읽었다. 제품 정본은 고정 SHA의 DEV 실행 기록·QA 보고·D12 Google 연결 절이다. 지시서·route·packet은 이 체크아웃의 미커밋 원천 snapshot으로 읽었다. 예외는 없다. 규칙 4개·project.md·review/rule.json·document-writing.md는 snapshot과 이 체크아웃이 줄바꿈(CRLF)만 다르고 내용이 같음을 비교로 확인했다.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_tester.md`, 고정 SHA의 `docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md`, `docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-DEV/report.md`, D12 `ops-guide.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 변경 57개 파일. OCR 대상 50, 증거 제외 7(DEV qa-reports 로그 6과 server-exits.txt). reviewed 36, skipped 21. skipped는 제외 증거 7, 생성된 Jev 증거 13, 보존 handover 로그 1이다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 8 / 실행 불가 0. 아래 "lint 증거 검증"에 근거를 적었다.
- SIZE-002: 지시서의 예상 규모는 없다. 실제 추가 1126줄(lint 대상 33개 파일)이 상한 400줄을 넘는다. 근거는 아래 "SIZE·DEP·SEC 경고 판정"에 적었다.
- DEP-001: `adapters/package.json`의 `scripts.test`만 바뀌었다. 의존성·lock·`go.mod`·`go.sum`·`db/` 변경은 없다(`git diff --stat`으로 확인). 새 런타임·개발 의존성은 없고 표준 라이브러리와 기존 패키지만 쓴다.
- UI 디자인: 해당 없음(변경 아님)으로 판정했다. 서버 렌더 HTML 한 템플릿(`device`)과 기존 템플릿의 폼 일부가 바뀌었으며 `project.md`가 Tailwind·공용 컴포넌트·디자인 lint 미구성을 명시한다. 새 템플릿은 기존 `memberStyle`을 재사용한다. DEV는 desktop과 mobile 390×844 Chrome 직접 검수를 기록했으나 캡처가 없다. 이 리뷰는 시각 판정을 하지 않았고 DESIGN lint 결과도 없다. 시각 수락은 designer 담당이다.

## 검토 범위

result.json의 모든 `(path,status)`에 reviewed 또는 skipped와 근거를 기록했다.

- 읽은 코드: `device.go`·`google.go`·`member_agents.go`·`member.go`·`capacity.go`·`connections.go` 전체 diff, 호출 흐름의 `identity.go` signIn·session, `store.go` hashToken·sweep, `http.go` Origin 보호 경계, `core.ts` `relayBase`·`request`(redirect "error"), `login.ts`, `private-files.ts`, `connect.ts`, `mcp.ts`, `text.ts` 전체 diff와 모든 테스트 파일, ingress 조각과 self-check, Makefile·package.json·`.gitattributes`.
- 읽은 문서: 설계 문서 5개(architecture·crud-design·data-model·interface-design·module-design)의 diff, D11·D12·D13 diff 전문, PLANS·contexts diff, DEV 실행 기록, DEV QA 보고, README 2개.
- skipped 21개의 사유: 제외 증거 7은 `result.json`의 reason대로 같은 폴더 요약으로 존재를 확인한다. 실제 확인 방법은 DEV 보고 표의 종료코드 4개(`server-exits.txt`)를 내가 서버 재실행한 종료코드와 대조한 것이다. Jev 증거 13은 packet-outcomes 요약(181 항목: completed 28, no_change 153)과 head 일치만 확인했다. 영향은 추적성 세부 미검증이다. 보존 handover 로그 1은 앞부분만 읽었다.
- 삭제·이름 변경 영향: `handovers/to_dev.md`는 215줄 비움 삭제다. 같은 SHA의 `handovers/logs/2026-10-10_to_dev.md`가 전문을 보존한다. 다른 삭제는 없다. `connect.ts`·`text.ts`의 인라인 파일 검사를 `privatePath`로 통합한 변경은 같은 기준(소유자, 그룹·타 사용자 권한, 링크 거부, 크기 8192)을 유지한다.
- 고위험 경계의 실제 판정은 다음 절에 적었다.

## 발견 사항

미해결 critical/high는 없다. result.json `findings`에 4건을 기록했다.

| 번호 | 심각도 | 파일·줄 | 내용 | 상태 |
|---|---|---|---|---|
| 1 | medium | `internal/relay/device.go:21-45` | 링크 기반 사회공학 잔여 위험. 공격자가 자기 공개키로 시작한 요청의 링크를 받은 피해자가 Google 로그인 뒤 지문 비교 없이 승인하면 피해자 계정에 공격자 키의 새 agent가 생긴다. 페이지가 "자기 클라이언트에서 시작한 요청만 승인"을 경고하고 지문을 보인다. 새 agent는 관계가 없어 상대 수락 없이는 메시지가 오가지 않고, 기존 agent·키는 바뀌지 않는다. 설계 한계이며 수락 가능이다. | 열림(수락) |
| 2 | medium | `internal/relay/device.go:36-45`, `connections.go` sweep | 익명 start가 전역 Device 요청 2000개 상한을 쓰고 만료 뒤 24h 보존한다. 익명 한도(IP 30/min, 전역 200/min)로 약 10분이면 상한이 차고 최대 24h 동안 신규 클라이언트 연결이 `capacity`로 거부된다. 기존 키·이메일 로그인·기존 agent에는 영향 없다. 미결합·미소비 요청의 조기 삭제 또는 IP별 상한을 권고한다. | 열림(권고) |
| 3 | low | `internal/relay/device_test.go:41-45` | `finishGoogleLogin`의 "기존 다른 회원 거부" 분기(`m.Owner != c.Owner`)를 DEV 검사가 검증하지 않는다. DEV의 foreign은 새 subject라 앞선 `m == nil`에서 거부된다. 서버 mutation에서 이 조건을 제거해도 DEV 검사는 통과했고 tester HTTP 검사만 실패했다. 제품 동작은 올바르다. 기존 회원이 선점된 요청을 가로채면 `c.Owner`가 바뀌므로 회귀 가치가 높다. QA의 `device_tester_qa_integration_test.go` S1을 DEV 회귀에 추가하길 권고한다. | 열림(권고) |
| 4 | low | `internal/relay/device.go:47-58` | 동의 때 agent를 먼저 만든다. 클라이언트가 complete 응답을 잃으면 키 없는 agent가 활성 5개·기록 10개 한도를 소비한다. 사용자 안내(D11)가 홈에서 폐기하도록 적었고 `login.ts` 오류 문구도 같다. | 열림(문서화됨) |

기대 동작을 확인한 항목(결함 아님): 결합 전 `c.Owner == ""`이면 첫 callback이 결합한다. 두 번째 계정은 `finishGoogleLogin`이 `invalid_auth`로 거부하고 세션도 만들지 않으며 신규 회원도 만들지 않는다.

## 고위험 경계의 실제 판정

| 경계 | 판정 | 근거 |
|---|---|---|
| Google 인증(state·nonce·PKCE·ID token) | 통과 | 기존 `consumeGoogleAttempt`·`verify`는 변경 없다. 연결 ID는 시도 레코드(`GoogleAttempts`)에 서버 쪽으로 묶인다. 서버 `TestGoogleHTTP`·`TestGoogleDeviceHTTP` PASS. |
| Origin/CSRF | 통과 | 모든 신규 POST가 기존 `crossOrigin.Handler(mux)` 뒤에 있다. tester S7: cross-site `/auth/google`(connection 포함)·`/home/device-confirm`은 403. 비브라우저 API(`/v1/connect/*`)는 Origin 헤더가 없고 서명·token으로 보호한다. |
| 만료·replay | 통과 | 요청 TTL 10분. 만료 시 페이지·Google 재시작·poll 모두 401(S5). start/approve/complete/callback replay 거부는 DEV 검사와 서버 재실행으로 확인. 시도 레코드는 1회 소비. |
| 회원·클라이언트 분리 | 통과 | 선점·경쟁 계정 거부(S1), 타 회원의 보기·재시작·취소·승인·`/home/confirm` 거부(S2). 같은 회원의 두 요청은 별도 agent·키·credential이고 pairing은 없다(DEV `TestGoogleDeviceHTTP`가 `Agents==2, Pairs==0` 확인). 한계는 finding 3. |
| 개인키 저장·권한 | 통과 | 키는 클라이언트에서만 생성하고 서버는 공개키만 받는다. complete는 개인키 서명(`KNOWSLINK-CONNECT` 바이트)이 필요해 token만으로 credential을 받을 수 없다(S3). 브라우저 URL은 token 해시다. 로컬 폴더는 mkdir 실패로 덮어쓰기를 막고 POSIX 소유자·모드와 Windows ACL(.NET, 상속 차단, 소유자·SYSTEM·Administrators만 허용)을 검사한다. Windows 로컬 `connect.test`·`login.test` PASS. |
| 리다이렉트 경계 | 통과 | `next`는 `"/connect/"+a.Connection`이며 값은 서버가 찾은 연결 키다. 잘못된 connection 값은 303을 만들지 않는다(S6). `relayBase`와 `request`의 `redirect:"error"`로 클라이언트가 외부로 따라가지 않는다. |
| 사용자 화면·동의 | 통과(시각 미판정) | 동의는 Google 계정 + 최근 인증(5분) + 지문 표시 + 전용 POST로만 이뤄진다. 기존 `/home/confirm`은 Device를 승인하지 못한다. Google `prompt=select_account consent`. 시각 품질은 이 리뷰의 범위 밖이다. |
| 운영 ingress | 통과(적용 미검증) | `public-ingress.yml` 정규식은 공개 경로만 허용하고 `/owner`·`/v1/owners`·`/v1/authorize`·`/healthz`·`/v1/test`를 막는다. self-check exit0. 조각은 edge Access를 바꾸지 않으며 실제 Cloudflare 적용·경로 우선순위 검증은 OPS 단계다. 이 리뷰는 운영 Access·Tunnel을 건드리지 않았다. |
| 새 의존성·SQL | 통과 | 변경 파일 목록에 `go.mod`·`go.sum`·`package-lock.json`·`db/`가 없다. 상태는 기존 JSON 행에 `Device bool`·`Connection string`만 더한다. |
| held·기존 로그인 호환 | 통과 | 새 MCP 도구는 `public-node` 설정이 없으면 `held`를 반환한다. 서버 `TestEmailIdentity`·`TestGoogleHTTP` PASS. |

## lint 증거 검증

원천: `D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/google-connect-final-lint.json`을 `lint.json`으로 보존했다(변경 없는 복사).

- head `152217f6…`는 이 리뷰의 head와 같다. base·merge_base는 `f9f7675…`와 같다. `test_level`은 lite다.
- `config_sha256` `9c49bb2d…`는 snapshot의 `.fullops-squad/lint/lint.json`의 sha256과 같다. tester 체크아웃의 현재 `lint/lint.json`(`49572df6…`)과는 다르다. 이 증거는 검토 SHA의 설정으로 계산됐으므로 일치 기준은 snapshot이다.
- 등록 명령 `product-lint`(`make lint`)와 `product-test`(`make test`)는 passed, exit 0이다. 실행 불가 0, ERROR 0, WARNING 8이다. 서버 clean clone에서 실행했다는 DEV 기록과 `google-connect-server-handoff.json`의 `final_gate_command`가 일치한다.
- 변경 없는 등록 전체 검사는 반복하지 않았다. 대신 위험 경계의 좁은 검사를 직접 실행했다(QA 보고).

### SIZE·DEP·SEC 경고 판정

| 경고 | 판정 |
|---|---|
| SIZE-001 `PLANS.md` 1296줄 | 기존 큰 문서의 +7줄. 수락. |
| SIZE-001 `mcp.ts` 319줄 | 새 도구 2개(+72줄). 한 도구씩 분리하면 오히려 설정이 늘어 수락. 상한 초과는 경고로 유지한다. |
| SIZE-001 `member.go`(357)·`member_agents.go`(331) | 기존 큰 모듈에 템플릿·라우트 추가. 수락. |
| SIZE-002 추가 1126줄 > 400 | 지시서의 예상 규모가 없다. 추가의 대부분은 테스트(`device_test`·`device_integration_test`·`login.test` 417줄)·`login.ts`·`private-files.ts`·`device.go`와 문서다. 한 기능의 신뢰 경계를 분리하면 검사가 약해져 분할 불필요로 판정한다. |
| DEP-001 `package.json` | `scripts.test`에 `login.test.js`를 추가한 것이다. 의존성 변경 없음. |
| SEC-001 `device_integration_test.go:47` | `ClientSecret: "local-fixture"`는 로컬 합성 OAuth 설정이며 실 자격이 아니다. 오탐. |
| SLOP-004 `check_public_ingress.py:19` | 사용자가 읽는 성공 출력이다. 디버그 잔재가 아니다. |

## 검증 및 남은 제약

별도 문서 `qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md`에 명령·종료코드·증거를 적었다. 요약은 다음과 같다.

- 직접 독립 실행(운영 서버의 별도 임시 공간·격리 DB, Linux): TypeScript 빌드·`login.test`·`connect.test`, `go test -race ./internal/relay -run TestDevice`, migration, DEV의 Google/이메일 integration 3종, 내가 쓴 HTTP 경계 검사 `TestTesterGoogleDeviceBoundaries`. 모두 exit 0.
- 직접 독립 실행(Windows 로컬, Node 22.22.2, Docker 없음): `connect.test`·`login.test`·`check_public_ingress.py`. 모두 exit 0.
- mutation 두 건으로 tester 검사가 경계 회귀를 잡는지 확인했다. poll 서명 검사 제거와 기존 다른 회원 거부 분기 제거는 모두 tester 검사만 실패했다.
- DEV 증거 재사용: DEV의 `make lint`·`make test` 종료코드와 lint.json은 재실행하지 않고 위 lint 증거와 대조했다.
- 실행하지 않은 것: 실제 Google 계정 로그인, 실제 Cloudflare Access/Tunnel 적용, 외부 Bot 설치·왕복 전달, 시각 검수. 합성 검사를 그 성공으로 주장하지 않는다. 같은 Google 계정의 두 실제 클라이언트·관계 수락·왕복은 OPS 수락 조건으로 남는다.

## 대화 미참조 인계 점검

이전 대화 없이 `.fullops-squad/FULLOPS.md`와 지시서에서 시작해 아래를 찾았다.

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `handovers/to_tester.md`, 고정 SHA의 `exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` "기술 계획과 원천" | 확인 | DEV 지시서 로그의 "기준 SHA" 한 줄(`c3f918d…`)이 front matter의 base(`f9f7675…`)와 다르다. 본문 표기 불일치이며 동작에는 영향이 없다. |
| 구조와 구현/미완료 상태 | `design-docs/module-design.md` "Google 연결 모듈", `crud-design.md`·`data-model.md`·`interface-design.md`의 Google 연결 절 | 확인 | 코드와 일치한다(Device 상태 흐름·2000개 상한·proof 바이트). |
| 실행·검증 방법과 증거 | DEV `qa-reports/SAR-GOOGLE-CONNECT-001-DEV/report.md`, 서버 로그 4개, `fullops-gate/google-connect-final-lint.json`, 이 리뷰의 `lint.json` | 확인 | DEV 로그의 종료코드와 서버 재실행 결과가 일치한다. |
| 운영·복구(해당 시) | `operations/ops-guide.md` "Google 연결 공개 경로 적용", `transition.md` | 확인 | 5단계 적용·복구 순서가 있다. 실제 Access/Tunnel 적용은 미실행이다(OPS). |
| 다음 작업·담당·재개 조건 | `ops-guide.md`·`transition.md`의 "부모 담당" 문장, `PLANS.md` | 확인 | main 통합·OPS 적용·실제 Google/Bot 수락이 후속이다. finding 2·3은 DEV 후속 후보다. |
| 로컬 링크·절 접근/지원 한계 | `README.md`·`adapters/README.md`의 user-guide 링크 | 확인 | 링크 대상이 존재한다. anchor 정밀 검사는 하지 않았다. |
| snapshot 정리 후 정본 접근 | 이 리뷰 폴더(`result.json`·`report.md`·`lint.json`)와 QA 폴더 | 확인 | 결과는 snapshot 밖 기록 체크아웃에 있다. snapshot 정리는 부모가 release 뒤 `review.py cleanup`으로 처리한다. |

## 검토 결론

수락 가능이다. 근거: 미해결 critical/high 없음, 고위험 경계 9개 모두 통과, lint 증거의 head·base·config가 검토 SHA와 일치하고 ERROR 0이며 등록 test·lint exit 0, 서버 독립 QA와 Windows 로컬 검사 통과. 다음 조건을 남긴다.

1. finding 2(2000개 상한 잠금)와 finding 3(회귀 검사 간극)은 병합을 막지 않는 권고다. DEV 후속 과제로 coor가 배정한다.
2. finding 1은 설계 한계다. D11에 "다른 사람이 보낸 링크를 승인하지 않는다"를 이미 적었다. 공개 전 운영 안내에 한 번 더 강조하길 권한다.
3. 실제 Google·Cloudflare·외부 Bot 수락과 시각 검수는 이 결론에 포함되지 않는다.

이 결론은 AI 검토자의 판단이며 `review.py check` 통과가 내용을 보증하지 않는다.
