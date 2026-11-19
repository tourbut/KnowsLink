---
title: SAR-SERVICE-UI-001-TESTER 독립 기능 QA 증거
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SERVICE-UI-001-TESTER]
summary: 회원 웹 고정 후보 c7828b2의 lite 기능·보안·키보드·반응형 QA 실행 명령과 결과와 미검증 범위를 기록한다
---

# SAR-SERVICE-UI-001-TESTER 독립 기능 QA

## 대상과 환경

- 대상 SHA `c7828b227b5715f7bdf02bff46bf9cf102a758d9`. 비교 기준 `8539ef5b6d3e7ed99259762c62d164b8ba40fab0`(merge-base `bd24be06b20ff21cb9e8de4680345b29fa9b139c`).
- 테스트 레벨 lite. 하위 에이전트는 쓰지 않았다(검증이 하나의 순차 흐름이다).
- 실행 세션 `1fe0133c-f609-42f7-a774-bf5b8b5765c2`. 구현자 세션 `5953ccd3-7f4f-4c5e-b63d-c02d36ee35ac`와 다르다.
- tester 워크트리 HEAD `030c5bd`는 `c7828b2` 위에 문서만 더했다. `git diff --stat c7828b2 HEAD -- . ':!.fullops-squad'`가 비어 있어 제품 바이트가 같다. 합성 HTML은 이 워크트리에서 만들었고 리뷰 snapshot에는 쓰지 않았다.
- 합성 fixture 29개. 실제 계정·비밀값·운영 데이터는 쓰지 않았다. 이 QA는 실제 Google 로그인 왕복의 성공 증거가 아니다.

## 실행 결과

| 검사 | 위치·명령 | 결과 |
|---|---|---|
| Linux `make lint`, `make test` | 고유 clone `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2/repo`(git bundle → detach `c7828b2`, 실행 전후 clean). 절차 `linux-prep.sh`·`linux-gate.sh` | 종료코드 0, 0. [linux-exits.txt](linux-exits.txt), [lint 로그](linux-product-lint.log), [test 로그](linux-product-test.log) |
| FullOps `lint.py` | 같은 clone, `--from 8539ef5…`. 절차 `linux-lint.sh` | 종료코드 0, ERROR 0, WARNING 4(SIZE-001 3건, SIZE-002 1건), 실행 불가 0. 결과 `../SAR-SERVICE-UI-001-TESTER-review/lint.json` |
| 표적 Go 테스트 | 같은 clone, `go test ./internal/relay -run 'TestMemberUIFixtures\|TestMemberPagesShowNextSteps' -count=1 -v` | 둘 다 PASS |
| integration·browser 태그 컴파일 | 같은 clone, `go vet -tags 'integration browser' ./internal/relay`, `go test -tags 'integration browser' -c` | 종료코드 0, 0. 절차 `linux-vet.sh` |
| form 계약 golden 대조 | `git show 37b85a4:…member_forms.golden`(재설계 전 템플릿)과 `c7828b2`의 golden `diff` | 차이는 홈 `post /auth/logout` 2행(`home`·`home-empty`)과 원격 동의 fixture 2행뿐. 새 action·hidden 필드 변경은 없다 |
| 폭 측정 | Orca 내장 브라우저, fixture 29 × 1440·390·320 = 87건, 모든 `details` 펼침. `measure.js`, 결과 [browser-measure.jsonl](browser-measure.jsonl) | 가로 overflow 0건(`scrollWidth<=innerWidth`), 요소 밖 넘침 0건, 44px 미만 버튼·입력·select·summary·이동 링크 0건, 라벨 없는 입력 0건, `lang=ko`, `h1` 1개 |
| 외부 자원·인라인 | `home` 390에서 `[src],link,script,img,iframe,object,embed` 수와 `[style]`·`on*=` 속성 수 | 모두 0 |
| 한국어 줄바꿈 | `home` 390 계산 스타일 | `body` `word-break:keep-all`·`overflow-wrap:break-word`, `anywhere`는 `code`·`pre`에만 있다 |
| 헤딩 구조 | `home` 390 | h1 → h2(에이전트·관계·승인 요청·연결 확인·계정과 보안) → h3(에이전트 행·새 관계 초대·로그아웃) → h4. 건너뜀 없음 |
| 버튼 이름 | `home` 390 | 같은 이름의 버튼 중복 0건. 수락·거절·철회·발급은 대상 식별자를 `aria-label`에 포함 |
| 키보드 | `home` 390, `keyboard-walk.jsonl` | 첫 Tab이 "본문으로 건너뛰기"(`y=8`, 3px 외곽선). Enter로 `#main`. Tab 순서가 화면 위→아래와 같다. 28번 이동 모두 3px 외곽선. 닫힌 `details` 안의 버튼은 Tab으로 닿지 않는다. 에이전트 `summary`는 Enter 열림·Space 닫힘·Enter 열림. 위험 구역 `summary`와 그 안의 버튼 외곽선이 `rgb(160,22,22)`(`--danger`). 위험 구역을 펼친 뒤에만 "이 키 철회"·"이 에이전트 철회" 버튼에 닿는다 |
| 텍스트 간격 확대 | `home`·`connection-prepared`·`device-prepared`·`receipt`·`connection-token`·`start-both-problem` × 390·320, 줄 높이 1.5·글자 간격 .12em·단어 간격 .16em·문단 간격 2em 주입 | 가로 넘침 0건, 잘린 요소 0건 |
| 확대 간접 확인 | 320px 폭 | 1440px 기준 450% 확대와 같은 폭이다. 가로 스크롤 0건. 실제 브라우저 200% 확대는 직접 조작하지 않았다 |
| 홈 길이 | `home` 390(에이전트 4·관계 5·승인 요청 2) | 문서 높이 3984px. D04 관찰(O-01)의 7986px(에이전트 4)보다 짧다 |

## D04 필수 요구(R-01–09)와 수락 기준 대조

| 항목 | 판정 | 근거 |
|---|---|---|
| R-01 토큰 공통 적용·한 곳 정의 | 충족 | `member_pages.go` `memberStyle` 한 상수. 회원 홈·로그인·확인 대기·연결 관리·기기 승인·receipt·오류·gate·OAuth 동의가 같은 상수를 쓴다(`grep memberStyle`) |
| R-02 줄바꿈 | 충족 | 위 계산 스타일 |
| R-03 홈 4.2 순서 | 부분 충족 | 헤더·확인 필요 띠·섹션 이동·에이전트·관계·승인·연결 확인·계정이 모두 있다. 960px 미만에서 섹션 이동(y=89)이 확인 필요 띠(y=199)보다 위다. D04 4.2는 띠를 앞에 둔다. 영향은 낮고 시각 판정은 designer 몫이다. 960px 이상은 왼쪽 열이라 차이가 없다 |
| R-04 접힌 행·위험 구역 | 충족 | 에이전트 행 `details`, 위험 구역 `details`, 결과 문장, 키별·에이전트 철회 |
| R-05 3단계·지문 강조 | 충족 | `aria-current="step"`, `.fp.big` 지문. `connection-*`·`device-*` fixture 29개 중 해당 화면 측정 통과 |
| R-06 문구 표 | 충족 | 문구 대조(6절 표의 모든 현재→제안 행). "계정 비활성화 준비 중"·"다닷"·"관계 세대" 화면 문구 제거 |
| R-07 오류 화면 다음 동작 | 충족 | `refusal`·`connection` 만료·취소·`gate-closed` 모두 링크가 있다. 만료·취소에 "연결 취소" 버튼 없음(F-UI-03) |
| R-08 A-01–A-10, I-01–I-08 | 대부분 충족 | 아래 접근성·불변식 표 |
| R-09 로그인 설명 3줄 | 충족 | `start` 템플릿 목록 3개와 한 줄 설명 |
| R-10(권장) 확인 필요 띠 | 충족 | 받은 초대·대기 승인 요청 건수 링크 |
| R-11(권장) 관계 정렬 | 충족 | `homePage`의 `rank` 정렬. `home` fixture 첫 행이 받은 초대 |

범위 밖(후속): R-12 다크 모드·실제 접속 상태·클라이언트 복사 버튼. 요구하지 않았다.

### 접근성·불변식

| ID | 판정 | 근거·비고 |
|---|---|---|
| A-01 | 충족 | 위 헤딩 구조, `lang=ko`, 페이지별 `title` |
| A-02 | 부분 | `header`·`nav`·`main`·`aria-labelledby` 있음. 건너뛰기 링크 동작 확인. OAuth 동의 두 화면(`remote-consent`·`remote-login`)에는 `main`이 없다(DEV가 마크업을 바꾸지 않았다고 기록, R-01 목록 밖). 리뷰 F-1 |
| A-03 | 충족 | 위 키보드 행 |
| A-04 | 계산 승계 | 토큰 값이 D04 8.1과 같음을 코드로 확인했다. 대비 픽셀 측정은 하지 않았다 |
| A-05 | 충족 | 배지는 기호+글자(`✓ 연결 완료`, `… 미연결`, `— 철회`) |
| A-06 | 충족(코드) | `aria-describedby="email-help problem"`, `code-help problem`, `role="alert"`가 입력 위에 있다 |
| A-07 | 부분 | 값 옆에 라벨 텍스트가 있고 `user-select:all`이다. 연결 수단 `code`의 `aria-labelledby`는 `code` 역할이 이름을 받지 않을 수 있다. 실제 스크린 리더 미검증. 리뷰 F-2 |
| A-08 | 부분 | 320px·텍스트 간격 확대 통과. 200% 확대 직접 조작·스크린 리더 낭독은 미검증(DEV도 같음) |
| A-09 | 충족 | 위험 `summary` "이 에이전트 철회·키 철회", 버튼 `aria-label`에 대상 포함 |
| A-10 | 충족 | `display:none` 숨김 텍스트 없음(`.skip`은 화면 밖 위치) |
| I-01 CSRF·form 계약 | 충족 | golden 대조, `gate`의 `csrf` hidden 유지, `crossOrigin.Handler`(`http.go:189`) 변경 없음 |
| I-02 CSP | 충족 | 두 CSP 설정(`http.go:193`, `capacity.go:169`)이 diff에 없다. 인라인 `<style>`만 쓰고 `style=` 속성 0건 |
| I-03 owner 분리 | 충족(코드) | `homePage`의 `agents`·`pairs`·`gates` 수집은 변경이 없다(owner 비교). 새 건수 계산은 이미 걸러진 목록만 센다. integration 재실행은 아래 한계 |
| I-04 재인증 | 충족(코드) | 서버 검사 코드 변경 없음. `Recent`가 거짓이면 에이전트 위와 계정에 재확인 안내와 버튼이 보인다 |
| I-05 파괴 확인 | 충족 | 별도 모달·확인 화면 없음. 위험 구역을 펼친 뒤에만 단일 `POST` 버튼에 닿는다(키보드로 확인) |
| I-06–I-08 | 충족 | "온라인" 등 없는 상태 표시 없음, 마스킹 이메일, 24시간 보존 안내 유지 |

## DEV 공개 편차의 영향

- 서버 안내문(`safeProblem`, `agentLimits` 등)에 영문 "agent"가 남는다. D04 6절이 허용한 통일이며 정책 검사가 문장을 고정한다. 영향 낮음, 별도 과제.
- 만료 시각이 `2026-10-12 07:55:00 KST` 형식이다. D04 5.4는 짧은 날짜를 "우선"으로 적었다. F-UI-04(KST 유지)는 충족한다. designer 판단.
- 확인 필요 띠에 "재인증 필요"가 없다. 대신 에이전트 영역 위 안내가 항상 보인다. DEV의 이유(5분 뒤 항상 참)가 타당하다.
- 홈 헤더와 계정 섹션에 로그아웃이 둘이다. 같은 action·method·hidden 없음. 중복이며 위험은 없다.
- 받은 초대에 "관계 철회"가 수락·거절 옆에 함께 있다. form 계약 보존이다. 위치는 designer 검수 의견.
- 연결 방식 select의 긴 선택지는 닫힌 상태에서 390px에 잘릴 수 있다(브라우저 select 동작). 펼친 목록의 전체 문구는 확인하지 않았다.

## 미검증과 한계

- 실제 Google 로그인 왕복은 실행하지 않았다. 사용자 세션 로그아웃·실 agent 생성·폐기는 하지 않았다.
- fixture는 정적 HTML이다. 응답 CSP 헤더·쿠키·Strict 세션은 이 측정에 없다. 서버 코드 diff로만 확인했다.
- integration·browser 태그 테스트(owner 분리·CSRF·gate·OAuth 반환 회귀)는 이번에 DB 없이 실행하지 않았다. 컴파일(`go vet`)만 통과했다. DEV가 `0b485ad`에서 같은 태그로 통과(`ok … 114.875s`)한 증거가 있고, `0b485ad`→`c7828b2`의 제품 diff는 `member.go`의 템플릿 블록을 `member_pages.go`로 옮긴 것과 `README.md` 문구뿐이다. 줄 단위로 대조했다(`member.go`에서 사라진 줄은 `"html/template"` 임포트 한 줄 외에 모두 `member_pages.go`에 있다). 이 증거는 DEV 실행이며 이 QA의 자기 실행으로 기록하지 않는다.
- 실제 스크린 리더, 200% 브라우저 확대, 대비 픽셀 측정, 다크 모드는 범위 밖이다.
- 직접 시각 검수(정보 위계·여백·색 인상)는 designer 후속이다. 이 QA가 시각 수락을 대신하지 않는다.
- 캡처·영상은 만들지 않았다.
- Context7: 새 외부 라이브러리·SDK가 없다(표준 HTML·CSS·Go `html/template`). 외부 문서 조회가 필요하지 않았다.

## 보존 자원(삭제하지 않음)

- Linux `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2`(clone·로그·lint.json). `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c`는 첫 시도의 빈 폴더다(번들 업로드 경로 오류로 clone 실패, 이후 `-r2`를 썼다). `/tmp/knowslink-serviceui-tester-s1fe0133c.bundle`, `/tmp/knowslink-serviceui-tester-tools-s1fe0133c.tar`.
- 로컬 snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-c1644dfcf27644b3bc46f3c1788d9816`(읽기 전용 detached).
- Docker 컨테이너는 만들지 않았다.
