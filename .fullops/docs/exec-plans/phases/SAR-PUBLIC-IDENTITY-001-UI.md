---
title: 일반 이메일 신원 직접 UI 검수 실행 기록
status: review
updated: 2026-10-05
owner: designer
tasks: [SAR-PUBLIC-IDENTITY-001-UI]
summary: 독립 실행 환경과 Orca 캡처 복구 및 별도 Chromium 시각 검수와 정리를 보존한다
---

# SAR-PUBLIC-IDENTITY-001-UI — 실행 기록

## 범위와 기준

제품은 detached `59b66ada8b36802484cc6d7e22523257b50572cc`다. 기록 준비는 `1207bdf4542fac98d228f86de79aad9f4e126ce8`이며 기준은 `94533b207b456c0560800fe30a7c90b2b5887c6e`다. Task `task_1a9cb470bd74`, Dispatch `ctx_a3efa466a50f`, Run `run_8ca8bc058ab7`를 사용했다.

designer 인박스 전문과 coor 탐색 JSON을 먼저 읽었다. FULLOPS.md의 필수 규칙과 프로젝트·제품·UX 정본을 읽었다. member.go·identity.go·http.go·mail_sink.py, README·사용자 가이드·DEV 실행 기록으로 실제 경로를 확인했다. 초기 읽기 명령의 잘못된 `coding.md` 경로는 실제 정본 `coding-style.md`로 바로잡았다.

수정 소유권은 QA 보고, 이 실행 기록, UI 증거, designer 인박스와 해당 전문 로그다. 제품 코드·수치·기술 승인·PLANS/board·다른 인박스·운영·CF·다른 QA 실행 환경은 변경하지 않았다. 컨텍스트 갱신은 이번 소유 범위에 없어 하지 않았다.

## 독립 실행과 환경 오류

1. `/tmp/knowslink-public-identity-ui-59b66ad`의 HEAD와 깨끗한 작업 트리를 확인했다.
2. 별도 PostgreSQL 17 컨테이너 `knowslink-identity-ui-59b66ad-pg`와 loopback port를 만들었다. 기존 migration의 Up SQL을 적용했다.
3. 초기 준비에서 `pg_isready`가 임시 bootstrap 서버를 먼저 감지했다. DB 생성 전 SQL은 exit 2로 실패했다. 자기 컨테이너만 제거하고 실제 대상 DB `SELECT 1` 준비 확인으로 재시작했다. 제품 실패로 기록하지 않았다.
4. 후보의 `go build -o <자기 scratch>/relay ./cmd/relay`는 exit 0이었다. 소스·제품 build 폴더는 변경하지 않았다.
5. 후보 mail_sink.py와 자기 relay를 loopback에서 실행했다. 환경 변수의 제품 인증값은 출력하거나 저장하지 않았다. 합성 signup·trial·trusted IP header는 비웠다.
6. fixture 상태별 검사에서 자기 DB challenge Exp·session Verified를 조정했다. 화면 재검수 시 자기 send·HTTP bucket을 초기화했다. 시간 경계와 모든 한도 QA로 보고하지 않았다.
7. 발송 실패 검사 뒤 sink를 같은 port로 빨리 재기동하자 `OSError: [Errno 98] Address already in use`가 났다. 수정 캡처 시도는 대기 화면 timeout으로 exit 1이었다. 새 자기 SMTP port를 선택해 같은 후보 바이너리·같은 자기 DB로 relay를 재시작했다. 최종 시도는 exit 0이었다.

scratch는 `/tmp/knowslink-public-identity-ui-59b66ad-run`이다. 실행 helper와 바이너리는 레포 밖 자기 자원이다. 제품 소스 변경이나 제품 의존성 설치는 없었다.

## Orca와 직접 시각 보완

Orca `orca-cli`·`orchestration`의 버전 맞춤 가이드와 `references/browser.md`를 읽었다. 새 page `56050624-50c1-4700-b313-a3ed78e71059`만 명시 제어했다. 기존 Cloudflare 페이지에 goto·입력·클릭을 하지 않았다.

Orca에서 실제 이메일 fixture 입력·SMTP 코드 확인·오답·재발송·만료·홈·전체 종료·재확인·재로그인·현재 종료·메일 실패를 조작했다. 화면 밖 버튼의 native click이 이동하지 않은 경우 해당 버튼에 focus와 Enter를 사용했다. 이 도구 동작을 제품 버튼 결함으로 판정하지 않았다.

blank PNG를 확인해 escalation `msg_3a0f2a6297e4`, `msg_bd48f3e31e99`와 coordinator ask를 보냈다. 한 장의 복구를 전체 복구로 오인하지 않도록 뒤의 재발을 함께 보고했다. coor는 브라우저 병행 조작을 중지하고 전용 페이지 독점 뒤 별도 실제 Chromium 대안을 허용했다. 허용 답변은 이 Dispatch의 Run 질문/답변에 있다.

Orca 실제 정상 PNG 2개와 blank PNG 1개를 보존했다. 최종 별도 Chromium `147.0.7727.15`·Playwright core `1.63.0`의 비영속 context에서 같은 서비스로 실제 요청을 보냈다. 기존 설치를 읽기 전용으로 사용했고 다른 QA 프로필이나 실행 상태에는 접속하지 않았다. 지정 폭은 390×844와 1280×900이다. 최종 PNG 15개를 직접 열어 검수했다. 캡처 전에 이메일 텍스트만 치환하고 입력·홈 신원값을 가렸다. 로그인 코드·메일·cookie·token 전문은 기록하지 않았다.

Context7 library resolve는 월간 quota 초과로 실패했다. API 확인은 공식 [BrowserType.launch](https://playwright.dev/docs/api/class-browsertype#browser-type-launch), [Browser.newContext](https://playwright.dev/docs/api/class-browser#browser-new-context), [screenshots](https://playwright.dev/docs/screenshots)를 2026-10-05 확인했다. launch의 executablePath/headless, 새 context viewport, screenshot의 fullPage·mask만 사용했다. 제품 SDK 선택이나 기술 승인을 수행하지 않았다.

## 결과·근거

[직접 UI 보고](../../evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md)가 판정 정본이다. 로컬 fixture UX01–03은 PASS다. 기존 gate의 UTC 기본 만료 문자열에 low F-UI-01을 남겼다. 이 UI에서 새 critical/high는 없다. 별도 보안 리뷰 finding은 이번 검사로 닫지 않는다.

자동 결과는 증거 폴더의 `playwright-observations.json`에 상태·경로·브라우저 버전·폭·같은 회원 여부만 남겼다. `orca-commands.jsonl`은 비밀 인자 없이 명령 종류·page·exit를 남겼다. `fixture-events.jsonl`은 조정 항목과 시각만 남겼다. `visual.log`·`visual.exit`는 마지막 실제 브라우저 실행 exit 0의 근거다. PNG의 실제 판정은 자동 결과와 구분한다.

실제 이메일·운영 SMTP·공개 배포·QA-P06·QA-P07·최종 동일 이메일 노우↔다닷은 미실행이다. 미실행을 로컬 코드 단계 실패나 전체서비스 PASS로 바꾸지 않는다. 코드 단계의 통합 판단과 운영 최종 수락을 구분한다. 새 보안 수정 후보의 UI delta는 별도 후속이다.

## 정리와 완료

자기 relay·sink를 종료했다. 자기 DB 컨테이너와 anonymous volume을 제거했다. 자기 private 메일을 지웠고 세 loopback port 종료를 확인했다. 전용 Orca 페이지를 닫았다. 후보 checkout은 계속 깨끗한 detached59다. 증거 `cleanup.json`에 종료 사실과 명령 exit 0을 남겼다.

기록 파일과 증거만 검증·커밋한다. 제품 코드 변경이 없어 product-lint 적용 대상은 없다. 자체 파일의 front matter·로컬 링크·JSON·PNG 형식·민감값 비노출·Git 공백을 확인한다. 역할 완료 보고 전문을 인박스에 쓰고 `work.py finish --role designer --key SAR-PUBLIC-IDENTITY-001-UI`로 보존한다. 결과 SHA는 이 Run의 worker_done과 해당 Git 커밋을 정본으로 사용한다.
