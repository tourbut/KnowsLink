---
title: 회원 웹 상용 UI 구현 실행 기록
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SERVICE-UI-001-DEV]
summary: 로그인·내 에이전트·연결 관리 SSR 화면의 상용 디자인 구현 계획·변경·검증·시각 검수 환경과 미검증 범위를 기록한다
---

# SAR-SERVICE-UI-001-DEV — 회원 웹 상용 UI 구현

## 기준

- 지시서: `handovers/to_dev.md`, attempt `69804b27a2274875a19179a73e833ae9`, 기준 SHA `b6a4f4da545cfdd3f57862f8724f3eacbe00c8e2`(작업 시작 HEAD `bd24be0`).
- 디자인 정본: D04 [service-member-ui.md](../../design-docs/mockups/service-member-ui.md)(main 수락본). R-01–R-11, A-01–A-10, I-01–I-08을 적용했다. D04는 수정하지 않았다.
- 규칙: `fullops-common-0.3.3`, `project.md`, `contexts/dev.md`, `docs/agents/document-writing.md`. 테스트 레벨 lite. Ponytail full.
- 스킬: `C:/Users/shin/.agents/skills/frontend-design/SKILL.md`를 읽었다. D04가 색·글꼴·간격을 고정하므로 정본을 따랐다. 강조는 지문 블록 한 곳에만 두었다. 그라디언트·그림자·알약 버튼·장식 통계·대문자 라벨은 쓰지 않았다.
- 하위 에이전트: 읽기 전용 조사 1개(packet 문서·옛 문구 사용처 조사). 구현은 직접 했다.
- 외부 라이브러리 문서 조회는 없다. 새 의존성이 없고 표준 HTML·CSS만 사용했다.

## 기술 계획

1. 구조는 그대로 둔다. 기존 `memberStyle` 상수를 D04 8절 토큰의 단일 정의로 바꾼다. gate(`http.go`)와 OAuth 동의 화면(`remote_auth.go`)은 같은 상수를 이미 쓰거나 쓰도록 바꾼다. 새 파일·패키지·JS·폰트·이미지는 없다.
2. 상태 배지는 `badge` 템플릿 함수 하나로 만든다. 글자·기호·색을 함께 쓴다.
3. 홈 데이터는 기존 transaction 안에서 계산한다: 관계 정렬(R-11), 활성 에이전트 수, 받은 초대·대기 승인 요청 건수(R-10). 새 상태·API는 없다.
4. form 계약을 먼저 고정한다. 재설계 전 템플릿으로 합성 fixture를 렌더하여 `testdata/member_forms.golden`을 커밋(`37b85a4`)한 뒤 화면을 바꾼다.
5. 검증은 lite다. 기존 단위·integration 문구 검사를 새 문구로 고치고, 고정 golden·외부 자원 0건·Orca 내장 브라우저 폭 측정을 추가한다.

## 변경

| 파일 | 변경 |
|---|---|
| `internal/relay/member.go` | `memberStyle` 토큰(8.1 색, 8.2 글꼴·`keep-all`+`break-word`, 식별자·지문·코드만 `anywhere`, 8.3 간격, 8.4 버튼 44px, 8.5 배지·안내, 8.6 포커스). `badge` 함수. 모든 회원 화면 템플릿 재구성. 홈 관계 정렬·건수 계산. 문구 표(6절) 적용 |
| `internal/relay/member_agents.go` | `memberAgent.ActiveKeys`, `memberPair.Mine/Peer`. 거부·연결 화면 제목의 "agent"를 "에이전트"로 |
| `internal/relay/http.go` | gate 템플릿이 `memberStyle`을 쓴다. 문구·hint/body 구분·비활성 버튼 숨김·CSRF hidden은 같다 |
| `internal/relay/member_ui_test.go`, `testdata/member_forms.golden` | 합성 fixture 29개, form 계약 golden, 외부 자원 금지, 시각 검수용 HTML 출력 |
| `connections_test.go`, `*_integration_test.go` 4개 | 바뀐 문구·마크업에 맞춘 기대값과 helper 정규식. 검사 의도는 같다 |
| README, D10, D06(interface-design), user-guide | 바뀐 버튼 이름과 "계정 비활성화 준비 중" 제거 반영 |

화면별 요점:
- 로그인: 한 줄 설명과 3줄 목록, "Google로 계속" 주 버튼, 이메일 로그인이 켜져 있으면 "또는" 구분선 뒤 입력. 둘 다 꺼지면 "로그인을 준비하고 있습니다". 오류는 입력과 `aria-describedby`로 연결한다.
- 홈: 헤더(제품명·마스킹 이메일·로그아웃), 확인 필요 띠(받은 초대·대기 승인 요청이 있을 때만), 섹션 이동(960px 이상 왼쪽 sticky 열), 에이전트, 관계, 승인 요청, 연결 확인, 계정과 보안.
- 에이전트: 접힌 `details` 행(식별자·배지·활성 키 수). 펼치면 식별자(`user-select:all`), 키 목록, 연결 추가, 위험 구역 `details`(결과 문장, 키별 "이 키 철회", "이 에이전트 철회 (모든 키·관계 종료)"). 위험 버튼은 위험 색 외곽선이고 `aria-label`에 대상 ID를 넣는다.
- 관계: 내 에이전트 → 상대, 받은/보낸 초대 라벨, 상태 배지, 대기 중일 때만 기한. "관계 세대" 숫자는 화면에 없고 hidden 값은 그대로다. 새 초대 폼은 접힌 `details`이며 관계가 없을 때만 펼친다.
- 연결 관리: 1·2·3단계 진행 표시(`aria-current="step"`), 만료·취소는 진행 표시 대신 안내. 연결 수단은 `code` 블록이며 바로 위에 "대화·로그·캡처에 붙이지 마세요". 준비·마무리 명령 두 블록. 지문은 큰 모노 블록.
- 기기 승인: 위험 색 왼쪽 선 경고 띠, 큰 지문, 권한 목록, 주 버튼과 취소.
- receipt·gate: 2열 정의 목록과 토큰만 공유. 문장 의미는 같다.

### 의도적 판단과 미적용

- 홈 헤더 로그아웃(4.2)과 계정 섹션 로그아웃(5.8)을 둘 다 두었다. golden의 유일한 차이(홈 `post /auth/logout` 1개 추가)다. 같은 action·method이며 hidden 값이 없다.
- 받은 초대에는 기존처럼 수락·거절과 별도의 "관계 철회" form이 함께 있다. form 계약 보존을 우선했다. 위치 조정은 designer 검수 의견에 따른다.
- 확인 필요 띠에 "재인증 필요"는 넣지 않았다. 5분이 지나면 항상 참이 되어 띠가 계속 보이기 때문이다. 대신 `Recent`가 아니면 에이전트 영역 위에 재확인 안내와 버튼을 먼저 보인다(5.8).
- 서버 안내문(`safeProblem`, `agentLimits`, `keyRecordsFull`, `transportNames`)의 "agent" 용어는 바꾸지 않았다. D04 6절은 통일을 "허용"으로 정했다. 정책 검사가 문장을 고정하므로 별도 과제가 적합하다.
- 만료 시각은 기존 `clock()` 형식(`2026-10-12 07:55:00 KST`)을 유지했다. 짧은 날짜 표기(5.4)는 적용하지 않았다.
- 연결 방식은 select를 유지했다. 390px에서 긴 선택지 "키 교체 (…)"는 닫힌 select 안에서 잘릴 수 있다. 기본값 "새 키 추가 (최대 3개)"는 잘리지 않는다.
- OAuth 동의 화면(`remote_auth.go`)은 마크업을 바꾸지 않았다. 새 토큰이 적용되며 1440·390 렌더에서 overflow가 없었다.
- 토큰 값은 D04 8.1과 같다. 조정 값은 없다.

## 검증

| 검사 | 위치·대상 | 결과 |
|---|---|---|
| form 계약 golden | Windows, `37b85a4`(재설계 전 템플릿)에서 생성 → 이후 커밋에서 비교 | 차이는 홈 헤더 로그아웃 1건(의도). 원격 동의 fixture 2행은 `678736d`에서 추가 |
| `go test ./internal/relay`, `go vet -tags 'integration browser'` | Windows Go 1.27.1 | 종료코드 0 |
| 외부 자원 0건 | `TestMemberUIFixtures`, 29개 화면 | `src=`·`<link`·`url(`·`<script`·`@import`·`<no value>` 없음 |
| 폭·터치 크기 | Orca 내장 브라우저, 29개 fixture × 1440/390/320, 모든 `details` 펼침 | 가로 overflow 요소 0, 44px 미만 버튼·입력·select·summary·이동 링크 0. `scrollWidth`는 모두 뷰포트 폭 이하 |
| 키보드 | 같은 브라우저, 홈 390 | 첫 Tab이 "본문으로 건너뛰기"(보임, 3px 외곽선) → Enter로 `#main`. Shift+Tab으로 에이전트 `summary` 도달, Enter 닫힘·Space 열림. 위험 구역 `summary` 외곽선 `--danger` |
| 등록 product-lint·product-test, integration | Linux 고유 clone | 아래 Linux gate 절 |

CSP 헤더 코드(`http.go` 공통 middleware, `capacity.go`)는 바꾸지 않았다. owner 분리 검사(`TestPublicAgentHTTP`, `TestEmailIdentity`의 타 회원 비노출)는 기대 문구만 바꾸고 그대로 둔다.

### Linux gate

기존 Linux 서버에서 과제 전용 helper(paramiko, `.env.server` 값은 출력에서 가림)로 실행했다. 커밋의 git bundle을 고유 `/tmp` 경로에 clone하고 해당 SHA를 detach checkout했다. Node v22.22.2, Go toolchain go1.27.1. 실행 전후 clone은 clean이었다. 운영 relay·DB·Tunnel은 사용하지 않았다.

| 대상 SHA | 경로 | 검사 | 결과 |
|---|---|---|---|
| `678736d` | `/tmp/knowslink-serviceui-678736d-ctxd554b501` | `make lint`, `make test` | 0, 0 |
| `678736d` | 같은 경로, 격리 Postgres17 container `knowslink-dev-serviceui-ctxd554b501` | `go test -race -tags 'integration browser' ./internal/relay` (`TEST_BROWSER=/usr/bin/google-chrome`) | 1. `TestPublicAgentHTTP` helper가 옛 마크업(`<code>agent_…</code>`, `readonly value=`)을 정규식으로 찾다가 panic했다. 제품 결함이 아니라 검사 helper의 마크업 의존이다 |
| `0b485ad` | `/tmp/knowslink-serviceui-0b485ad-ctxd554b501` | helper 정규식을 새 마크업(`<summary><h3><code class="select">`, `aria-labelledby="token-l">`)으로 수정 뒤 `make lint`, `make test` | 0, 0 |
| `0b485ad` | 같은 경로, container `knowslink-dev-serviceui-0b485ad-ctxd554b501` | migration, 위 integration+browser 전체 | 0, 0 (`ok … 114.875s`). owner 분리·CSRF·gate·OAuth 동의 브라우저 반환 회귀 포함 |

container 두 개는 지시서의 "컨테이너 정리 금지"에 따라 `docker stop`만 하고 삭제하지 않았다. 이름은 위 표와 같다. 정리는 coor·OPS 판단이다. 최종 archive HEAD의 등록 gate(lint.py) 결과는 완료 보고에 적는다.

## 시각 검수 환경 (designer 후속)

비밀값·실계정 없이 합성 데이터만 쓴다.

1. 고정 SHA를 checkout한다. Go 1.22 이상이 필요하다.
2. 빈 폴더를 만들고 `KNOWSLINK_UI_FIXTURE_DIR=<빈 폴더 절대경로> go test ./internal/relay -run TestMemberUIFixtures -count=1`을 실행한다. 29개 HTML이 생긴다.
3. Orca 내장 브라우저에서 `file:///<폴더>/<이름>.html`을 연다. `orca exec --command "set viewport 1440 900"`(또는 390 844)으로 폭을 바꾼다.
4. 주요 화면: `start-google`, `start-both-problem`, `verify`, `home`(에이전트 4·관계 5·승인 요청 2, 재확인 완료), `home-empty`(재확인 필요), `connection-token`·`-prepared`·`-approved`·`-consumed`·`-expired`, `device-prepared`·`-approved`·`-login`, `receipt`, `refusal-*`, `gate-active`·`gate-closed`, `remote-consent`.
5. 측정 예: `orca eval --expression "document.documentElement.scrollWidth<=innerWidth"`.

정적 파일에는 응답 CSP 헤더가 없다. CSP는 서버 middleware가 붙이며 이번 변경과 무관하다. 실제 서버 렌더가 필요하면 기존 integration 검사처럼 격리 DB로 relay를 띄운다. 실제 Google 로그인 왕복은 사용자 소유 계정이 필요하다.

## 미검증과 후속

- designer 직접 시각 검수(1440·390), tester 독립 기능 QA·코드 리뷰, main 통합, OPS 배포는 후속이다. 이 기록은 UI PASS가 아니다.
- 200% 확대·텍스트 간격 확대(A-08)와 스크린 리더 실제 낭독은 측정하지 않았다. 320px overflow 0으로 간접 확인만 했다.
- 실제 Google OAuth 왕복·실계정 데이터·운영 화면은 보지 않았다.
- 옛 검수 스크립트 `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/render_check.go`는 정규식으로 `memberStyle`을 찾는다. 상수 이름은 같지만 여러 줄이 되었으므로 그 스크립트는 재사용 전에 확인한다. 과거 증거라 수정하지 않았다.
- D04 관련 designer 문서 `SAR-PUBLIC-SERVICE-UX.md` 22행은 "계정 비활성화"를 홈 항목으로 적는다. designer 소유라 수정하지 않았다.
