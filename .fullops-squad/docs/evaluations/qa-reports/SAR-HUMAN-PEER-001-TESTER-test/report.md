---
title: SAR-HUMAN-PEER-001 독립 QA 결과
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-HUMAN-PEER-001-TESTER]
summary: 254ff6a 합성 A/B/C HTTP QA는 실행한 단언이 통과했고 키보드·벽시계·B의 둘째 agent는 미실행이다
---

# SAR-HUMAN-PEER-001-TESTER QA

판정 대상 SHA는 `254ff6a98b12b7e51f0835491300f7695425d264`다. 환경 정본은 `C:\Users\shin\AppData\Local\KnowsLinkDevTools\humanpeer001-runtime-ctx9ded\final-report.md`다. `runtime/report.md`와 그 폴더의 이전 SHA 캡처는 실패·회귀 이력이다. 이번 통과로 다시 이름 붙이지 않는다.

원격 트리는 검수 시작 시 clean이었다. 원격 HEAD는 254ff6a다. origin은 `http://127.0.0.1:41381`이다. healthz는 200이다. 컨테이너 이름은 `qa-ccaa8cd7-c30163b6`다. Postgres 이름은 `qa-ccaa8cd7-c30163b6-postgres`다. final-report 당시 SMTP는 `52205`다. 이후 프로세스 재시작이면 라이브 포트는 `runtime.json`의 `smtp_port`를 따른다. 통과 실행은 그 환경의 로컬 sink에서 OTP를 읽었다.

tester loopback은 `http://localhost:13238`다. forward 작업 ID는 `01a12fe7-f6c9-7930-81e9-340bbf7512a7`다. 이 프로세스를 종료하지 않았다. designer 전용 profile은 `3168fe9e-11d7-4ea3-b95f-86ababc9b56a`다. page는 `f7f75011-b98e-46d9-85f1-bfd1582f3665`다. URL은 `http://localhost:13238/home`다. 제목은 `대화 · KnowsLink`였다. 이후 키보드 검수는 이 page를 조작하지 않았다.

## 명령과 종료 코드

| 명령 | 종료 | 근거 |
|---|---|---|
| 원격 `make lint` | 0 | `1791806105370100800-254ff6a98b12b7e51f0835491300f7695425d264-make_remote` |
| 원격 `make test` | 0 | 같은 로그. owner-text와 shell_guard를 포함한다 |
| 격리 `TestHumanPostgres`, `TestHumanLegacyPostgres` | 0 | DB `humanpeer_integration_f46538af714c`, 3.622초. 패키지 전체는 아니다 |
| Windows `lint.py` | 1 | ERROR 0, WARNING 17, make unavailable 2. 제품 실패로 보지 않는다 |
| 독립 HTTP `independent_qa_remote.py` | 0 | `1791806343903484100-254ff6a98b12b7e51f0835491300f7695425d264-independent_qa_remote` |
| `review.py check --task-key SAR-HUMAN-PEER-001-DEV` | 0 | 문구 3건 반영 뒤 재실행. reviewed 88, skipped 155, WARNING 17 |

로컬 로그 디렉터리는 `C:\Users\shin\AppData\Local\KnowsLinkDevTools\humanpeer001-runtime-ctx9ded`다.

## 독립 HTTP

결과 필드는 `PASS`다. `secrets_in_output`은 false다. 실패한 check는 없다. 증거 디렉터리는 `/tmp/knowslink-humanpeer001-254ff6a98b12-c30163b6/qa-ccaa8cd7/tester-independent-b0128515`다.

회원은 새 합성 A/B/C다. 이메일 도메인은 `humanpeer.invalid`다. helper auth가 만든 alice/bob `private-sessions.json`은 이 실행이 덮어쓰지 않았다. designer fixture와 독립 회원은 같은 DB에 함께 있다.

실행한 단언은 시나리오 문서의 PASS 항목과 같다. 핵심은 healthz 200, 익명 홈 303, 익명 invite 401, OTP 후 303, 이름 한계, agent 생성, 수락 전 발신 403, 없는 코드의 이름 비노출, 제삼자 수락 403, 다른 Origin 403, 수락 303, 양방향 본문, authorship 구분, 멱등 303, 본문 충돌 409, 공백·과대 422, 추가 agent 발신 303, 반대 요청 비자동 수락, 취소와 거절, 차단 후 읽기 전용, 옛 attempt 403, 해제 후 친구 미복구, 새 generation 발신 303, 브라우저 세션의 contacts 거부, legacy 페이지에 새 본문 없음이다.

앞선 실패 디렉터리는 제품 결함이 아니다. `tester-independent-6d7b1dcb`는 OTP 대기 시간이 짧았다. `tester-independent-53f38cf8`는 주소 대소문자가 sink와 달랐다. `1791806196281424800` 실행은 친구 상세에서 상대 agent가 아니라 자신의 추가 agent 이름을 찾았다. 세 경로는 보존한다.

B의 두 번째 agent는 만들지 않았다. 수락 이후 추가 agent 발신은 A의 둘째 agent로 통과했다.

## 키보드와 200%

키보드 profile은 `bab6d5ab-3906-4d38-8dd6-42180a6e53ed`다. page는 `c41d2fbf-cdac-4a1d-9dd6-00e2290ec740`다. 친구 화면 제목은 `친구 · KnowsLink`다. `input.focus()`는 `#friend-code`를 잡는다. `orca keypress Tab` 뒤 `document.activeElement`는 그 input 또는 body에 남는다. 엔진이 DOM 포커스를 옮기지 않는다. Tab, Enter, Space의 제품 판정은 not_run이다. 사유는 `blocked_engine_focus`다. 제품 결함으로 보고하지 않는다.

CSS `zoom=2` 측정값은 `{"zoom":"2","clientWidth":1188,"scrollWidth":1188,"bodyScroll":594}`다. 그 너비에서 가로 overflow는 없다. 브라우저 확대와 320/390/1440 viewport 변경 API는 없다. 둘 다 not_run이다. 390/1440 캡처는 designer가 연 page에서 본다.

## 미실행

공유 DB에서 24시간 벽시계, 180초 큐 만료, 차단 직전 queued/leased/gate, goose Down, 재로그인, 같은 agent 자기 발신 거부, 차단 코드 조회, B의 둘째 agent, 브라우저 뒤로 가기, CSP와 rate flood는 not_run이다. 자동 검사와 코드 리뷰가 일부 경계를 덮는다. 그 사실을 라이브 PASS로 적지 않는다.

## 코드 리뷰와 결함

리뷰 세션은 `01a12fd5-0d90-7803-85ed-d7546b4f7189`다. 작성 세션은 `01a12f82-ca28-7de1-8f74-acbfa1a6c823`다. snapshot은 `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-cba0ca432fa64d08903c989aa30d9464`다. 미해결 critical/high는 없다. low 문구 3건은 리뷰 보고서에 있다. 제품 소스는 수정하지 않았다.

Windows `core.longpaths=true`는 공유 Git 설정이다. 제품 커밋이 아니다. snapshot, forward, 실패 증거, 키보드 profile을 삭제하거나 이동하지 않았다.
