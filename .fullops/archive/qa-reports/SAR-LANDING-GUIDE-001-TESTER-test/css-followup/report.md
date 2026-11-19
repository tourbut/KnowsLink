---
title: 회원 CTA 대비와 공개 예시 간격 수정의 독립 수락 QA
status: review
updated: 2026-10-15
owner: tester
tasks: [SAR-LANDING-GUIDE-001-TESTER]
summary: CSS 두 결함의 독립 수정후 판정과 원증거 및 gate 재사용을 구분한다
---

# SAR-LANDING-GUIDE-001-TESTER CSS 후속

수정 범위의 독립 동작·시각 판정은 PASS다. 원 V-LG-01의 예시 간격과 V-LG-02/F-LG-01의 회원 CTA 대비는 수정 후보에서 해소됐다. 과거 FAIL을 변경하지 않았다. 전체 제품 수락·고정 SHA 코드 리뷰·designer 수락·main 통합·운영 배포는 이 보고의 완료 선언에 포함하지 않는다.

## 입력과 실행 경계

- 제품·fixture source: `f2b8cb622dc397034423582978e77560e0355771`.
- DEV 완료 및 이번 lint base: `fc4ad58b65471e2de298cb9f1aba5d2253ea8135`.
- 새 QA 실행 HEAD: `db7fd8ad381585df28278e732cce45188abc543d`. 최종 증거 커밋은 별도 인계한다.
- task `task_5d4833216e3f`, dispatch `ctx_cd7d7e4e46a6`, attempt `53ebe1ae346444fc9ccf3d04466cdc88`.
- codex/gpt-6-astra/high, test lite, subagent standard. 단일 fixture의 짧은 순차 검사이므로 하위 worker는 사용하지 않았다.
- fullops-common-0.3.3, project.md, document-writing.md, frontend-design과 D04 공개 화면·여섯 색 토큰·24px·44px·focus 규칙을 적용했다.

전용 Chrome155.0.8059.40/CDP1.3, Node24.13.0을 사용했다. 새 profile은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-css-tester-20261015/chrome-profile`, CDP는 `127.0.0.1:18901`이다. DEV CDP18897과 기존 profile·서버는 제어하지 않았다. 모든 자료와 자기 Chrome을 보존하며 PID는 같은 private 폴더의 `chrome.pid`에 남겼다. 정리·종료 판단은 coor에 인계한다.

허용 origin은 GET 전용 합성 fixture `http://127.0.0.1:18896`이다. 파일은 `start-google.html`, `human-home-empty.html`, `human-home.html`, `human-account.html` 네 개다. DOM 읽기·스크롤·details 펼치기·viewport·pseudo-state와 자기 Chrome 설정의 확대만 실행했다. 폼 제출·외부 링크 이동·로그인·DB·계정·메일·메시지·설치·배포·삭제·이동은 실행하지 않았다. 탐색·shadow·CUAWright·새 설치는 OFF다.

## 실제 실행과 원 종료코드

명령 cwd는 tester 레포 루트다. 아래 파일 경로는 이 보고 폴더 기준이다. 종료코드는 파이프 없이 원 프로세스에서 받았다.

| 명령 | 원 exit | 결과 |
|---|---:|---|
| `node <css-followup>/browser.mjs` 최초 | 1 | 새 Chrome 기동 직후 CDP 연결 `ECONNREFUSED`. [원 로그](browser-startup-failure.log). 제품 검사 전 환경 준비 실패다. |
| 같은 명령, `/json/version` 응답 확인 후 | 0 | 독립 44 checks PASS. [로그](browser.log), [측정·요청·이미지 hash](results.json). |
| `node <css-followup>/zoom.mjs inspect` | 0 | 자기 Chrome Page zoom100%와200% 선택지 확인. [설정](zoom-settings-before.json). |
| `node <css-followup>/zoom.mjs apply` | 0 | native ArrowDown5회로 실제 browser200%. 새 예시 경계만 검사했다. [결과](zoom-results.json). |
| `python3 <css-followup>/audit.py` | 0 | 18 provenance checks PASS,212경로 bounded 원천 기록. [동일성](identity.json), [원천 기록](source-audit.json). |

audit 최초에는197경로를 기록했다. packet unknown에만 있던15개 과거 문서의 확인 누락을 발견해 범위를 보완했다. 제품 검사는 반복하지 않았다. 브라우저 검사의 의도적 navigation은12회이며, 나머지는 CSS 측정·viewport·스크롤·pseudo-state다. Jev 조작 추론이나 비용은 없고 결정 가능한 조건을 기존 CDP 검사 코드로 확인했다. 기존 DEV 스크립트에서 독립 profile·실행 메타데이터·빈 배열 방지·대표 영향·경계 캡처를 보강했다.

## 수정 후 판정

| 항목 | 수정 후 독립 결과 |
|---|---|
| V-LG-02 / F-LG-01 회원 CTA | 1440·390·320, 기본·hover·focus-visible 모두 글자 `rgb(255,255,255)`, 배경 `rgb(0,90,98)`, 대비7.958059850759797:1. 4.5 이상 PASS. |
| CTA 크기·focus | 높이51.59375px,44px 이상. focus outline3px solid accent, offset2px. pseudo-state 검사이며 실제 Tab 사용성 재시험이 아니다. |
| V-LG-01 예시 하단 | 세 폭 모두 `.examples` 다음 h3까지24px. D04 기존 간격 토큰과 일치한다. |
| 공개 영향 | hero CTA 대비·44px PASS. 열린 details에서 가로overflow 없음. |
| 회원 영향 | 빈 홈·populated 홈·계정 가로overflow 없음. populated 일반 대화 링크는 존재하며 ink색·44px 이상 유지. 계정의 서비스 안내 href `/` 유지. |
| 새 간격 실제200% | setting200%, DPR2, visualViewport.scale1, CSS zoom1. innerWidth712, scrollWidth704, 예시 간격24px PASS. |
| 요청 경계 | browser.mjs 기록 요청은 fixture origin의 GET뿐이다. auth 성공이나 실제 회원 데이터의 근거가 아니다. |

직접 연 이미지: `member-empty-1440.png`, `member-focus-320.png`, `examples-1440.png`, `examples-390.png`, `examples-320.png`, `member-populated-1440.png`, `member-account-320.png`, `browser-zoom-200.png`.

흰 CTA 문구는 청록 배경에서 읽힌다. 320px focus 외곽선은 버튼을 둘러싸며 잘리지 않는다. 예시 세로선 아래와 '대화에서 알아둘 점' 제목 사이에 여백이 생겨 별도 구획으로 읽힌다. 모바일·200%에서도 한글 본문은 정상 줄바꿈하며 겹침·가로 잘림이 없다. 긴 본문은 세로 스크롤하며 캡처 위·아래의 viewport 경계는 콘텐츠 손실이 아니다. populated 목록의 선택 배경과 이름 계층, 계정 폼의 라벨·입력 폭은 유지된다. 계정 fixture의 합성 오류 안내는 기존 정적 샘플 상태다.

원 tester [report.md](../report.md)의 F-LG-01 medium1.846967:1과5149ac1 기록은 그대로다. 원 designer `SAR-LANDING-GUIDE-001-VISUAL/report.md`의 V-LG-01 low0px 및 V-LG-02 medium1.846967:1도 그대로다. 해당 designer 원문은 designer 체크아웃에서 읽었으며 다른 역할 파일을 수정하지 않았다. 새 해소 판정은 이 후보·네 fixture·측정 범위에만 적용한다. 새 critical/high 발견은 없다.

## source HTML과 재사용 증거

`identity.json`은 source f2b8cb6의 세 파일 Git blob SHA256과 manifest, QA HEAD의 동일성을 확인한다. 네 served HTML은 private fixture 파일 및 manifest SHA256과 byte 단위로 같다. 렌더된 memberStyle 전체와 공개 추가 style은 source Go 문자열과 일치한다. 이는 DEV가 기존 Go 테스트로 생성한 합성 정적 HTML이다. 이번 QA에서 Go template을 다시 생성하지 않았고 인증·HTTP 헤더·서버 session·DB 동작을 검증한 것으로 해석하지 않는다.

fc4와 새 QA HEAD의 cmd/internal/adapters/db/deploy/scripts·의존성·빌드·lint 설정 입력은 diff0이다. `identity.json`에 경로·원 diff exit0·Linux 원 gate receipt와 로그 hash를 보존했다. 기존 DEV helper는 exact HEAD·clean clone·known_hosts/RejectPolicy 및 원 exit 전파를 유지한다. 원 Linux source clone·install·gate 조건은 DEV css-resume 보고서와 private `landing-guide001-css-dev`의 fc4 로그에 연결된다. make install은 fc4 로그 exit0이며 기존 Node22.22.2/Go1.27.1 환경이다. 이번에 Linux 환경에 접속하거나 새 전체 gate를 실행하지 않았다.

재사용하는 product-lint/product-test는 **fc4ad58b65471e2de298cb9f1aba5d2253ea8135에서 실행된 Linux exit0**이다. 원 FullOps는 ERROR0/WARNING1/unavailable0이며 경고는 기존 PLANS SIZE-001이다. 이를 새 QA HEAD의 실행이나 새 독립 검사로 표시하지 않는다.

fe2387a에서 f2b8cb6까지 제품 변경은 두 CSS와 기존 fixture 출력용 테스트5줄뿐이다. 인증·DB·protocol·adapter·의존성 입력은 불변이다. 이전 HTTP34·키보드22 및 native 실제 입력·기존200% PASS는 원4795 fixture/5149ac1 조건으로 재사용한다. Chrome154의 과거 키보드 증거를 이번 Chrome155 실행으로 표시하지 않는다. 새 margin의 확대 영향만 위 새200% 검사로 확인했다.

## 미실행·완료 게이트

실계정·실Google·메일·송신·DB·운영·배포·스크린리더·회원 친구/agent 실페이지·전체 키보드 회귀는 not_run이다. 허용 범위 밖이거나 CSS 변경 영향 밖이다. 테마 전환과 별도 디자인 lint는 프로젝트에 없고, 일반 DESIGN lint는 Go 문자열 CSS 판정을 대신하지 않는다.

Context7 `resolve-library-id`는 월 quota 초과였다. 기존 검사 코드를 재사용하며 공식 [CDP CSS](https://chromedevtools.github.io/devtools-protocol/tot/CSS/#method-forcePseudoState)와 [Page](https://chromedevtools.github.io/devtools-protocol/tot/Page/#method-captureScreenshot)를 확인했다. 새 dependency는 없다.

packet의 stale producer와 partial/unknown19는 정본을 덮지 않고 보존했다. 필수 규약·D04·변경 source·원 실패 보고는 직접 읽었다. 나머지212경로는 제목·관련 발췌·hash로 영향 여부를 좁혀 확인했으며 전체 역사 재감사가 아니다. PLANS/ops-guide의 과거 정지·운영 기록은 최신2026-10-15 재개 지시와 별개다. 결과는 같은 키 outcomes에 연결한다.

coor는 착수 당시 미커밋인 same-key route/packet/history 입력도 이번 QA의 추적 증거로 함께 커밋하도록 답했다. 제품·다른 역할·PLANS·board는 변경하지 않는다. work.py finish로 전문을 보존한 뒤 커밋하고 현재 clean HEAD에서 지정 lint를 실행한다. Windows에는 make가 없으므로 원 unavailable·exit를 private `landing-guide001-css-tester-20261015/final-lint.json`과 `final-lint.log`에 남긴다. 원 fc4 Linux PASS를 이 상태와 구분한다. 최종 SHA와 실제 gate 결과는 같은 private 폴더의 `final-handoff.md` 및 worker_done에 기록한다.

work.py finish의 첫 시도는 exit1, `packet task/role/attempt/instruction identity differs`였다. 완료 체크 시 기존 체크리스트 문구를 바꾼 것이 원인이다. 문구를 원본으로 복원하고 체크 상태만 유지했다. 원 packet·해시·기준은 바꾸지 않았다.
