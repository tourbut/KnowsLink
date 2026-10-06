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
