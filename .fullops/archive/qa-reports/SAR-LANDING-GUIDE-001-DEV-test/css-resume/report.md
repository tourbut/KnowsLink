---
title: 재개된 공개 예시 간격과 회원 CTA 대비 수정 검증
status: review
updated: 2026-10-15
owner: dev
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 두 CSS 수정의 실제 Chrome 측정과 exact SHA Linux 검사 및 독립 검수 인계를 기록한다
---

# CSS 수정 DEV 결과

제품·기존 fixture 출력 변경 SHA는 `f2b8cb622dc397034423582978e77560e0355771`이다. 비교 base는 `fe2387a7676f7997aea76a2606f4cc2dbcd4b56b`, 재개 HEAD는 `dc5ec35b2c0fa297ba91387b9794b0be8398872c`, attempt는 `9a48a4180041493b93ddce9b5d1bf6f2`다. 완료 기록 커밋의 최종 SHA와 현재 SHA lint 영수증은 worker_done에 전달한다.

## 원인과 최소 변경

회원 공용 `.conversation-list a`가 `.primary`보다 우선하여 흰 글자를 ink로 덮었다. 일반 링크 색상만 `.conversation-list a:not(.button)`으로 분리했다. 기존 primary의 canvas/accent 토큰, 크기, border, hover와 focus 규칙은 유지했다. 실제 사용하는 빈 대화 CTA는 `human-conversation-list`의 `/home/agents` 링크다. populated 대화 목록의 일반 링크는 기존 ink를 유지한다.

공개 `.examples`는 다음 h3와 간격이 없었다. D04의 기존 24px 간격을 margin-bottom으로 적용했다. 모바일 gap과 공개 섹션 간격은 그대로다. 카피·템플릿 구조·인증 handler·DB·adapter는 변경하지 않았다. 기존 `TestHumanEmptyChoicesOfferRecovery`에 환경 변수로 선택하는 HTML 출력을 추가하여 결함이 발생하는 빈 회원 화면을 기존 fixture로 재현한다. 새로운 제품 기능이나 문자열 일치 테스트는 없다.

## 실제 검사

| 대상 | 결과 | 근거 |
|---|---|---|
| 관련 Go 검사 | exit 0 | `go test ./internal/relay -run 'TestLanding\|TestMemberUIFixtures\|TestHuman' -count=1`, [로그](related-tests.log) |
| exact source SHA make lint | exit 0 | 새 Linux clone의 [lint 로그](f2b8cb622dc397034423582978e77560e0355771-lint.log) |
| exact source SHA make test | exit 0 | 새 Linux clone의 [test 로그](f2b8cb622dc397034423582978e77560e0355771-test.log) |
| 실제 Chrome 계산·배치 | 38 checks PASS | [실행 코드](browser.mjs), [측정값·캡처 hash](results.json) |
| 회원 CTA | 모든 폭·기본/hover/focus-visible에서 7.958059850759797:1 | 흰 글자 `rgb(255,255,255)`, 청록 배경 `rgb(0,90,98)`, 16px, 높이 51.59375px |
| 공개 CTA·예시 경계 | 7.958059850759797:1, 경계 24px | CSS viewport 1440×900, 390×844, 320×844 |
| 대표 영향 | 공개 열린 details, 빈 회원·회원 계정 overflow 없음; 일반 대화 링크 ink 유지 | same results.json |
| focus-visible | solid 3px / offset 2px / accent | CDP pseudo-state 검사. 실제 키보드 완료 주장이 아님 |

DEV는 `member-empty-1440.png`, `examples-390.png`, `examples-320.png`를 직접 열어 흰 CTA 글자와 예시·후속 제목의 분리, 줄바꿈·겹침을 확인했다. 독립 designer 판정은 후속이다. 320px 예시 본문은 세로 스크롤되며 가로 넘침은 없다.

원 tester의 1.846967:1 실패와 designer V-LG-01 0px 실패는 원 보고에 보존한다. source SHA의 첫 FullOps lint 시도는 완료 문서가 미커밋이어서 exit 2, `작업 트리가 깨끗하지 않습니다. 커밋 후 실행하세요`였다. 제품 검사 실패나 PASS가 아니며 최종 clean commit에서 다시 실행한다. 이후 직접 실행한 make lint/test는 각각 원 종료0이다. 파이프 없이 helper가 SSH channel 종료코드를 반환하고 로그에 EXIT를 기록한다.

## 환경과 재현

기존 승인 helper `../gate/remote.py`를 읽고 HERE만 새 외부 경로로 바꿨다. [보존한 helper](remote.py)는 known_hosts와 RejectPolicy, exact HEAD·clean clone 검사를 유지한다. 저장소 검사 설정은 변경하지 않는다. make.cmd는 등록 `make lint/test`를 해당 helper로 전달하고 원 종료코드를 반환한다.

외부 자료는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-css-dev/`에 보존한다. fixtures는 source SHA에서 기존 Go 테스트가 렌더한 합성 HTML 40개다. [fixture manifest](fixture-manifest.json)는 HTML과 소스 파일 SHA256을 기록한다. loopback `http://127.0.0.1:18896/`에서 `start-google.html`, `human-home-empty.html`, `human-home.html`, `human-account.html`을 연다. 서버는 이 HTML 디렉터리만 제공한다. POST는 지원하지 않으며 제품 서버·계정·DB가 아니다.

새 Chrome headless profile은 위 외부 경로의 `chrome-profile`, CDP는 loopback 18897이다. `browser.mjs`는 관측·스크롤·details 열기와 CSS pseudo-state만 사용한다. 요청 기록은 fixture origin의 GET뿐이다. 실제 인증·메일·메시지·배포·삭제·이동·새 dependency는 없다. 기존 4795 fixture나 프로세스를 중지·교체하지 않았다. 새 static server와 Chrome PID는 외부 `http.pid`와 `chrome.pid`에 보존하며 후속 검수용으로 유지한다. 종료·보존 판단은 coor에 인계한다.

설치된 Node 22.22.2와 Chrome을 재사용했다. Context7 resolve-library-id는 월간 quota 초과여서 사용할 수 없었다. CDP의 [Runtime.evaluate](https://chromedevtools.github.io/devtools-protocol/tot/Runtime/#method-evaluate)와 [CSS.forcePseudoState](https://chromedevtools.github.io/devtools-protocol/tot/CSS/#method-forcePseudoState) 공식 문서와 기존 tester 실행 코드를 근거로 측정했다. 프레임워크·SDK·라이브러리 변경은 없다.

## 재사용과 남은 검수

`git diff fe2387a HEAD --name-only -- internal adapters db cmd go.mod go.sum Makefile compose.yaml scripts`의 제품 결과는 member_pages.go, landing_pages.go, human_pages_test.go 세 파일뿐이다. CSS 밖 인증·DB·adapter 입력과 의존성이 같으므로 원 source4795의 DB/인증 QA와 tester5149ac1의 실제 native keyboard·Chrome native200% PASS를 영향 밖 자료로 재사용한다. 이번 SHA에서 그 검사를 재실행한 것으로 쓰지 않는다. 구형 DEV/designer 문서의 native200% 대기는 최신 tester PASS로 해소됐으며 원 한계 기록은 보존한다.

고정 source 및 최종 기록 SHA의 독립 코드 리뷰, tester의 변경 영역 확인, designer의 예시 경계 판정, main 통합은 coor 후속이다. 이번 DEV 측정으로 독립 수락이나 운영 완료를 선언하지 않는다. D13 transition은 운영 인수 문서이므로 이번 CSS-only 범위에서 변경하지 않았다. PLANS/board와 다른 역할 파일도 읽기만 했다.

## 최종 기록 gate의 기존 문서 오류와 좁은 조정

기록 SHA `2cc87e3ac7011d0e54c4fac7b05a5b0d819df12c`의 실제 Linux product-lint/product-test는 각각 종료0이었다. FullOps 최종 판정은 기존 `SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md`의 front matter 누락 DOC-003 때문에 종료1이었다. [원 실패 JSON](2cc87e3-final-lint-doc003.json)을 보존한다. 기존 PLANS 1873줄의 SIZE-001 경고1은 coor 소유의 누적 이력이며 축약하지 않았다.

coor 답변 `msg_ed3a7dcf3e7d`는 이 한 문서의 front matter만 공식 stamp로 추가하도록 소유권을 확장했다. 원 본문·판정·세션·SHA는 바꾸지 않았다. 제품 경로는 2cc87e3과 동일하다. `lint.py --help`와 실제 실행부에는 다른 HEAD 검사 영수증의 재사용 옵션이 없고 등록 두 명령을 항상 실행한다. 최종 current SHA done-gate를 보존하기 위해 정본·helper·상태를 우회하지 않고 같은 Linux 방식으로 두 명령을 재실행한다. 2cc87e3 실행을 새 SHA 실행으로 재명명하지 않는다. work.py finish는 이미 성공했으며 빈 인박스를 다시 채우지 않는다. 추가 조정과 최종 SHA 결과는 이 보고와 외부 완료 영수증에 연결한다.

[source-audit.json](source-audit.json)은 packet39개 항목과 unknown/잔여 경로의 hash·관련 발췌를 보존한다. 역사 문서 전수 재감사를 뜻하지 않는다. 오래된 cache/중지/미착수/200% 대기는 최신 재개 절과 tester 근거로 해소한다. 큰 D04는 공개 첫페이지·색/간격/접근성 절을 직접 확인했다. 나머지 과거 인박스는 현재 실행 지시가 아니며 전체 기능·운영 수락을 바꾸지 않는다.
