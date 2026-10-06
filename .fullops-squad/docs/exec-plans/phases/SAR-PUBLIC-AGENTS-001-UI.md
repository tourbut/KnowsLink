---
title: SAR-PUBLIC-AGENTS-001-UI — 고정 후보 직접 UI 검수 기록
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI]
summary: 일반 회원 연결 화면의 격리 실행·직접 관측·장애·완료 처리와 후속을 기록한다
---

# SAR-PUBLIC-AGENTS-001-UI — 고정 후보 직접 UI 검수 기록

## 범위·결정

현재 인박스는 `handovers/to_designer.md`다. Task는 `task_8398bc20dd2b`, Dispatch는 `ctx_f6a28123e265`다. 제품 고정 SHA는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`이며 lint 기준은 `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`다. 적용 기준은 fullops-common-0.3.3·FULLOPS·project·PS04–07/PS11/개인정보·UX04–05다.

제품과 제품 규칙은 바꾸지 않았다. 정본 [직접 UI 보고](../../design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md)에 관측·finding·판정·재현·마스킹·미검증을 기록했다. D04 기존 원천에는 이번 결과 링크만 추가했다. medium F-UI-01/02로 필수 UX 시각 수락은 FAIL/보류다. 작업 완료와 제품 수락을 분리한다.

## 독립 실행·장애

새 `/tmp/knowslink-agents-ui-d1eef9b` clone과 자기 Compose project·DB·SMTP·relay를 사용했다. 처음 Compose interpolation은 DATABASE_URL 누락으로 실패했다. 설정을 보완했다. 첫 host migration은 DB가 internal network에만 있어 ping 실패였다. 기존 verify_mvp의 loopback port와 ingress network fixture 패턴을 적용한 뒤 migration은 exit 0이다. 이 실행 준비 오류를 제품 결함으로 표시하지 않았다.

Orca DOM 조작은 됐지만 viewport PNG는 1×1, full PNG는 배경뿐이었다. tab switch 뒤에도 blank였다. 이 두 PNG는 장애 증거로 보존했다. 전용 수동 CDP Chromium은 Page.goto timeout이었다. 전용 프로세스를 종료한 뒤 Playwright launch의 정상 engine으로 검수했다. 설치된 Playwright 1.60.0 API는 기존 UX 검수에서 확인한 launch/newContext/screenshot 근거를 재사용했다. 새 라이브러리 API·제품 SDK를 선택하지 않았다.

검수 중 중복 reauth 버튼과 section selector를 단일 요소로 좁혔다. context 전환의 시작 시점에는 새 page가 about:blank여서 fill timeout이 발생했다. 새 전용 page의 실제 시작 주소로 이동한 뒤 정상 진행했다. 이러한 조작 오류를 기능 FAIL로 기록하지 않았다. 재확인 60초 rate 오류는 실제 서비스 화면으로 별도 기록했고, 안내 시각 뒤 정상 재확인했다.

## 검증·정리·완료

UI PNG 31개와 장애 PNG 2개의 경로·해시·크기를 [manifest](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/evidence.json)에 보존했다. fixture 시간 경계 변경은 own relay를 멈추고 적용했다. 실제 TTL 대기는 미검증이다. 제품 고정 SHA 대비 tracked product diff는 exit 0이다.

strict deliverables·Git 공백·PNG/JSON 형식·마스킹·민감값 비노출을 확인한다. 새 기록 고정 HEAD의 FullOps lint를 기준 ref에서 실행한다. product-lint/product-test는 기준 ref 이후의 기존 제품 변경 때문에 실행되며 이번 역할의 제품 변경을 뜻하지 않는다. 고정 HEAD·원 명령 exit code·ERROR/WARNING/실행불가를 검증 JSON에 보존한다. 제품 코드를 바꾸지 않은 designer가 DEV 전체 verify-mvp를 반복하지 않는다.

자기 신규 자원만 회수하고 사용자 기존 탭·운영 자원·인증값은 보존한다. 인박스 완료 보고 전문을 `work.py finish`로 날짜 로그에 보존한 뒤 빈 인박스를 확인한다. 최종 기록 SHA와 결과 링크를 preamble의 worker_done으로 한 번 보내고 종료한다. 후속은 DEV의 medium 수정과 designer의 직접 재검수이며 PLANS에 자기 결과를 append한다.

## 초기 기록 HEAD의 검사

기록 `f9c5eded2db2382257a26a97bef8a7923d9785e0`에서 FullOps lint는 exit 1, ERROR 1/WARNING 8/실행 불가 0이었다. 등록 product-lint와 product-test는 모두 exit 0이다. 유일한 ERROR는 준비 HEAD에서 상속된 OPS 리뷰 빈 양식 `SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md`의 front matter 누락(DOC-003)이다. 초기 실패를 [JSON](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/lint-initial-f9c5ede.json)과 [원 출력](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/lint-initial-f9c5ede.txt)에 보존했다. 다른 역할 소유 파일이므로 preamble ask로 본문·판정 보존과 메타데이터만 보정하는 범위를 확인했다.

WARNING 8은 기존 제품 SIZE-001 5건, SIZE-002 1건, DEP-001 1건, SEC-001 1건이다. DEV가 기록한 기존 이유를 유지한다. 추가 줄 수는 기존 제품 변경과 이번 관측·문서·증거가 기준 ref 이후에 함께 포함된 결과다. 이번 역할은 라이브러리·dependency/version/lock을 바꾸지 않았다. 실제 CLI 검수·기존 source 읽기는 범위를 늘린 제품 변경이 아니다. 큰 증거 묶음은 같은 고정 SHA의 지정 상태를 추적하기 위해 한 결과로 보존했다.

문서 strict는 검사 13/미작성 0/문제 0/경고 0, exit 0이다. PNG 33개의 manifest·형식과 task JSON 검사는 exit 0이다. UI PNG 31개는 직접 열었고 인증값은 마스킹했다. 작업 시작 대비 제품 diff와 Git 공백 검사도 exit 0이다. 자기 자원 정리의 첫 확인은 Playwright SIGTERM 종료 지연으로 실패했다. 해당 자기 프로세스 그룹에만 SIGKILL을 보낸 뒤 정리 JSON에서 종료·port 폐쇄·Compose 자원 부재를 확인했다. 사용자 프로세스에는 신호를 보내지 않았다.

## 상속된 리뷰 양식 보정과 검사 통과

초기 기록 f9c5ede의 FullOps는 상속된 OPS 리뷰 빈 양식의 DOC-003 한 건으로 실패했다. product-lint/product-test는 exit 0이었다. coor는 preamble ask 답변으로 해당 파일의 metadata-only stamp를 허용했다. OPS 완성 보고서 70f26bc는 coor가 통합 때 유지한다. 빈 양식의 본문과 판정은 바꾸지 않았다. 원 본문 SHA256과 일치 검사는 exit 0이며 `review-metadata-only.json`에 보존했다. stamp의 front matter 구분 개행을 본문에서 분리해 비교했다. 이 예외는 UX FAIL이나 OPS 리뷰 완료를 PASS로 바꾸지 않는다.

기록 고정 `f2f56b6a57309eb675eb9a23dc68782d815408a6`의 FullOps lint는 exit 0, ERROR 0/WARNING 8/실행 불가 0이다. 등록 product-lint/product-test도 모두 exit 0이다. 기준 ref는 지시서의 d2f7ba5다. 최종 완료 로그 커밋 HEAD의 재검사 결과는 worker_done과 별도 final lint JSON으로 고정한다. 최종 기록 SHA와 제품 fixedSHA를 혼동하지 않는다.
