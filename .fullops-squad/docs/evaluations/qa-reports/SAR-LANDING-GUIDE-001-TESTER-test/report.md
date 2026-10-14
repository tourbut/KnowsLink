---
title: 공개 안내 독립 기능 QA와 회원 버튼 대비 실패
status: review
updated: 2026-10-14
owner: tester
tasks: [SAR-LANDING-GUIDE-001-TESTER]
summary: 합성 인증·반응형·실제 키보드·브라우저 200% 결과와 기존 회원 대비 결함을 구분한다
---

# SAR-LANDING-GUIDE-001-TESTER

독립 검증을 완료했다. 공개 안내·합성 인증의 핵심 검사는 통과했다. 회원 홈의 기존 버튼 대비는 실패했다. 전체 제품 수락은 미완료이며 DEV의 좁은 CSS 후속과 designer 수락이 필요하다. 제품 파일은 수정하지 않았다.

## 고정 입력과 권한

- 구현 및 lint 기준: `fe2387a7676f7997aea76a2606f4cc2dbcd4b56b`.
- 실제 서버 fixture: `4795d619513fde4fbcd4a55b94ec91b056ebbaf0`.
- 착수 기록 HEAD: `297e6d3679f6d866171a35775b615685221ffd47`.
- task `task_09a03f21ce42`, dispatch `ctx_02e6a3423e6e`, attempt `a43025df15404705846c60b97cdb8dc2`.
- codex/gpt-6-astra/high, test lite, subagent standard. 이 단일 fixture 검증은 하위 worker 없이 수행했다.
- 탐색·웹 검사는 FullOps 설치 도구1.4.3을 사용했다. 완료 직전 호스트 캐시가1.4.4로 바뀌어1.4.3 work.py가 errno2로 실패했다. 새 fullops-work/CLI help를 읽고 finish/lint만1.4.4로 실행한다. 저장소 규약1.4.2와 설정·검사 수준은 변경하지 않았다.
- AGENTS/FULLOPS, fullops-common-0.3.3, project.md, document-writing.md, D02 LG-01–08 및 D04 공개 첫페이지 절을 적용했다.

[identity.json](identity.json)은 fixture→구현→착수 HEAD의 cmd/internal/adapters/db/deploy/scripts와 의존성·빌드·설정 입력 diff0을 확인한다. `memberStyle`은 landing 이전 eae9a8b와 동일하다. `member_pages.go` 파일 전체는 start template 분리·parse 연결 때문에 바뀌었다. 최초 escalation의 파일 전체 diff0 표현은 msg_85d7c376f307로 정정했다.

허용 origin은 `http://127.0.0.1:18781`~`18784`다. Google/both/email/none 순서다. 합성 DB와 메일·Google 취소 mock만 사용했다. GET은 `/`, `/auth/verify`, `/auth/google/callback`, `/home`, `/home/name`, `/home/account`와 합성 helper `/fixture/code`, `/fixture/google-cancel`이다. POST는 `/auth/start`, `/auth/verify`, `/auth/google`, `/home/profile`, `/fixture/expire`만 실행했다. 실제 로그인·메일·메시지·설치·운영 배포는 실행하지 않았다. 요청·cookie·code 원문을 증거에 저장하지 않았다.

## 실행 결과

명령 cwd는 tester 레포 루트다. 아래 코드는 명령 자신의 종료코드다. 파이프로 종료코드를 숨기지 않았다.

| 실행 | 결과와 범위 |
|---|---|
| `jev_test_web.py --repo . --key SAR-LANDING-GUIDE-001-TESTER --scenario .fullops-squad/docs/evaluations/scenarios/web/landing-public-tester.json --validate` | exit0, 시나리오 유효 |
| 같은 실행기 `--page d3d68c8f-44ca-4ca6-9961-e73d1985260c --assert-only` | exit0, 독립 checks7개 PASS; 모델·조작 없음, 비용0 |
| `python3 <이 폴더>/functional.py` 최초 | exit1, none 검사에서 CSS `:disabled` 문자열까지 잘못 잡은 테스트 기대값 오류 |
| `python3 <이 폴더>/functional.py` 수정 후 | exit0, HTTP34개/요청30개 PASS |
| `node <이 폴더>/browser.mjs` | exit0, 실제 Chrome DOM22개 PASS, navigation/input12회 |
| CUA native 입력 | Enter·Space·Tab·Shift+Tab·앵커·폼 Enter의 실제 효과7개 확인, [기록](native-input.json) |
| `node <이 폴더>/contrast.mjs` | exit1, 공개 CTA PASS, 기존 회원 CTA FAIL |
| `node <이 폴더>/zoom.mjs inspect` 및 `apply` | 각각 exit0, 실제 Chrome Page zoom200% PASS |
| `node <이 폴더>/zoom.mjs expanded` | exit0, 열린 설치·오류 로그인 폼·하단 접근과 overflow PASS |
| `python3 <이 폴더>/source_audit.py` | exit0,224경로 bounded 원천 기록·runtime 동일성·기존 CSS 동일성 확인 |

브라우저는 설치된 Chrome154.0.8037.98/CDP1.3, Node24.13.0이다. 전용 headless profile은 `.fullops-squad/tmp/landing-tester-chrome`, page target은 `89708B0AB5EF43F8F03C7E4C08C19D45`다. Orca는 별도 isolated profile `4fc519c9-c493-4db2-adfc-445a88c671bf`를 사용했다. CUA는 별도 task tab의 실제 키보드만 검증했다. CUA session 이름만으로 storage 격리를 주장하지 않는다.

각 시나리오는 의도적 동작60개/20분 안에서 끝났다. CDP의12는 navigation/input만 센 값이며 viewport·details setup·scroll은 별도다. 이 동작을 포함해도60개 미만이다. native CUA는16동작이다. 새 탐색·shadow·CUAWright는 OFF다. 브라우저 네트워크 기록은 loopback만 포함한다. 외부 문서 링크는 href와 로컬 README 제목을 대조했고 실제 외부 페이지로 이동하지 않았다.

재현 순서는 fresh 합성 fixture에서 functional.py, browser.mjs, contrast.mjs, zoom.mjs inspect/apply/expanded다. functional.py는 처음에 이름이 없는 landing 합성 계정을 요구하며 이름 저장 뒤 상태를 보존한다. 같은 DB에서 무조건 재실행하는 reset 스크립트는 아니다. 현재 fixture는 시각 QA를 위해 그대로 유지한다.

## 요구사항별 판정

| 조건 | 판정·증거 |
|---|---|
| LG-01 공개 읽기·안전 경계 | 다섯 구획, 단일h1/ko, public200/private303, no-store/CSP, 비공개 이름·이메일 비노출 PASS |
| LG-02 로그인한 공개 `/` | 합성 POST→확인→이름 미설정303→이름 저장→회원 홈·계정→서비스 안내 PASS |
| LG-03–06 카피·설치·가상 예시 | D02/D04와 README/architecture/system-design/operations의 실제 원문 대조. 현재·향후 agent, 매 발송 승인, process/routine, Events/구독/표시, 신규/업데이트,180초/24h 구분 PASS |
| LG-07 설정·인증 오류 | 네 provider 구성, 형식422·wrong code401·미설정 메일503·Google 취소401·challenge410·세션 만료·CSRF403 PASS |
| LG-08 반응형 | 실제 CSS viewport1440×900/390×844/320×844, CTA 첫 화면·열린 details overflow·44px·URL 문자열 PASS |
| LG-08 native 키보드 | CUA의 실제 접기·역방향 focus·앵커·form 제출·코드 오류 복구 PASS. DOM dispatchEvent로 대체하지 않음 |
| LG-08 실제200% | 전용 profile의 Chrome settings select#zoomLevel을 native ArrowDown5회로100→200 변경. device override 해제·DPR2·visualViewport.scale1·CSS zoom1. innerWidth712, scrollWidth704. 열린 설치와 오류 폼도 PASS |
| 모션·공개 대비 | reduced-motion에서 animation/transition 없음 PASS. 실제 hero 흰 글자/#005A62 대비7.958060 PASS. 나머지 기존 token 대비는 DEV의 동일 CSS 근거 재사용 |
| 회원 주 이동 | /home 렌더와 세 주 메뉴 href·계정 왕복 PASS. /home/friends와 /home/agents 실제 진입은 fixture allowlist 밖이므로 not_run |
| 공유 auth/연결/Remote 성공 복귀 | 동일 runtime 입력의 DEV DB race integration 및 별도 code review 근거 재사용. 이 fixture에는 해당 성공 mock/endpoint가 없어 새 브라우저 실행은 not_run |

DEV의 make lint/test와 DB race integration은 원 SHA4795의 실행 기록이다. 입력 동일성을 확인해 인증 성공/Remote/연결 경계 근거를 재사용했다. 이를 tester의 새 실행이라고 부르지 않는다. QA 코드 추가에 대한 현재 기록 HEAD의 필수 lint/test는 아래 완료 게이트에서 별도로 실행한다.

## 실패 F-LG-01 — 기존 회원 CTA 대비, medium

1. 합성 이메일 인증 후 이름을 저장한다.
2. 활성 agent가 없는 회원으로 `/home`을 연다.
3. 최근 대화 영역의 ‘에이전트 추가’ 버튼을 확인한다.

실제 `.conversation-list .primary`는 글자 `rgb(24,43,49)`, 배경 `rgb(0,90,98)`,16px/600이다. 대비1.846967은4.5 미만이다. [contrast-results.json](contrast-results.json)과 [회원 홈 PNG](member-home-1440.png)에 재현 근거가 있다. `.conversation-list a`가 `.primary`의 글자색을 덮는다. 해당 공용 CSS 문자열은 landing 이전과 같다. 공개 CSS의 신규 회귀는 아니다.

coor는 msg_7f62757d938a에서 medium 실패 보존과 DEV의 좁은 CSS 후속을 지시했다. designer의 V-LG-01 low 간격 발견과 함께 후속 처리한다. tester는 제품을 고치지 않는다. 수정 SHA/fixture가 준비되면 CTA 대비·해당 예시 간격·대표 화면과 필요한 필수 영향만 재검증한다.

## 원 실패·미실행 보존

`http-first-failure.json`은 none 상태의 `disabled` 문자열 검사가 CSS selector를 오인한 원 실패다. 실제 button/input 속성으로 범위를 바로잡았으며 ‘가짜 disabled 로그인 없음’ 조건은 유지했다. 최초 route bind는 과제/역할/SHA 불일치로 실패했다. 공식 route 재분류→bind→code/doc find→context→packet으로 갱신했고 JSON을 수동 보정하지 않았다.

Orca native focus와 CUA Ctrl+plus는 효과가 없었다. CUA DOM read-only API의 document.hasFocus도 미지원이었다. 해당 실패는 [native-input.json](native-input.json)에 보존한다. 실제 keyboard와200%의 후속 성공은 별도 도구·전용 profile 증거이며 원 한계를 삭제하지 않는다. 기존 DEV의 nil DB panic, native click 실패, CSP fetch 실패, 포트간 cookie 공유·타일 캡처·fractional rounding 실패도 원 DEV 보고서에 남는다.

회원 친구/agent 실페이지·실제 Google 계정·실메일·Remote 설치·실상대 발송·장시간·재시작·실호스트 수락·스크린리더는 not_run이다. fixture allowlist와 과제 권한 밖이다. 운영 성공을 추정하지 않는다. 새 critical/high 발견은 없으며 medium1은 미해결이다.

## 산출물·탐색·완료 게이트

[browser-results.json](browser-results.json)은 각 캡처의 viewport/DPR/zoom100·SHA256와 요청 path를 기록한다. native200%는 [zoom-results.json](zoom-results.json), [열린 상태](zoom-expanded-results.json), [PNG](browser-zoom-200.png)로 분리했다. `member-home-1440.png`, `member-account-1440.png`는 coor에 msg_da72256ad625로 전달했다. designer의 직접 판정은 별도다.

[source-audit.json](source-audit.json)은 packet224개 경로의 HEAD hash·제목·관련 발췌다. 역사 문서 전체를 재감사했다는 뜻은 아니다. 필수 D04의 공개 절1–220과 회원 토큰·auth/fixture 원문은 직접 읽었다. 과거 최소 로그인·queued·303 기록은 현행 LG-02/ready200과 구분했다. frontend Sol 기본값은 최신 사용자 Astra 지정으로 해소했다. README의 실제 운영 이력은 새 페이지 수락과 별개다. oversized ops-guide의 첫 현행 절도 읽어 운영 배포·계정 검증이 이번 범위가 아님을 확인했다. packet partial/unknown19와 문서 잔여 후보를 보존하며 outcomes에 처리 이유를 연결한다.

Context7의 Chrome DevTools 조회는 월 한도 오류였다. 공식 [CDP protocol schema](https://raw.githubusercontent.com/ChromeDevTools/devtools-protocol/master/json/browser_protocol.json)로 Input/Page 명령을 확인했다. 새로운 library·패키지는 설치하지 않았다.

커밋 뒤 `python3 <설치 플러그인>/scripts/lint.py --repo <tester> --from fe2387a7676f7997aea76a2606f4cc2dbcd4b56b`를 실행한다. 최종 정확한 SHA·종료코드·ERROR/WARNING은 worker_done과 `.fullops-squad/tmp/landing-tester-gate/final-lint.json`에 연결한다. Windows에 make가 없으므로 [gate.py](gate.py)가 같은 QA HEAD의 새 Linux clone에서 등록 make lint/test를 실행한다. 기존 fixture clone의 동일 package-lock과 설치된 node_modules만 재사용하며 make install은 실행하지 않는다. known_hosts/RejectPolicy와 private 값 가림을 유지한다.

공유 fixture·DB·forward는 designer 검수 완료까지 유지한다. 원격 fixture와 종료 책임은 DEV의 browser-fixture.json/report.md 및 coor에 그대로 있다. tester의 Chrome profile·gate bundle·로그도 `.fullops-squad/tmp/`에 보존한다. task 전용 두 임시 하위 폴더만 Git 로컬 exclude에 등록해 credential/cache를 커밋하지 않는다. 삭제·이동 정리는 없다.

전체 수락·배포·PLANS/board 갱신은 coor 후속이다. work.py finish1.4.4는 exit0으로 전문을 handovers/logs/2026-10-14_to_tester.md에 보존하고 인박스를 비웠다. 최종 커밋·gate 뒤 한 번 worker_done으로 인계한다.
