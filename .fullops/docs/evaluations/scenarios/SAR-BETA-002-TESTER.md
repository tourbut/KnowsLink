---
title: SAR-BETA-002-TESTER — 루프백 소유자 브라우저 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-002-TESTER]
summary: 루프백 소유자 화면의 브라우저 절차와 통과 조건을 정의한다
---

# SAR-BETA-002-TESTER — 루프백 소유자 브라우저 시나리오

실행 위치는 tester 워크트리다. 대상은 `/home/shin/deploy/knowslink`의 detached HEAD `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`다.
제품 경로 `cmd`, `internal`, `adapters`, `db`, `Dockerfile`, `compose.yaml`, `Makefile`, `scripts`는 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`와 같아야 한다.
판정 명령은 `node .fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER-test/browser.mjs`다. 종료코드 0이 통과다.
브라우저 도구는 저장소 밖 `/tmp/sar-beta-002-browser`의 `playwright-core`다. 환경 변수 `SAR_BETA_PW`가 있으면 그 경로를 쓴다. 실행 파일은 `/usr/bin/google-chrome`이다.
저장소의 package 의존성은 바꾸지 않는다. `beta.sh owner-login`은 실행하지 않는다. trace, HAR, 스크린샷은 남기지 않는다.
credential은 fixture 파일에서 메모리로만 읽는다. 로그와 결과 JSON에 쓰지 않는다.

## 절차

1. 배포 HEAD와 porcelain을 확인한다. 제품 경로 diff가 비어 있는지 확인한다. `http://127.0.0.1:8080/healthz`가 200인지 확인한다. healthz는 브라우저 판정의 전제다.
2. `build/qa-fixture.json`의 바이트와 파일 모드를 보존한다. 값은 출력하지 않는다.
3. 보존한 fixture owner에게 상태 `expired`인 gate가 있으면 그 owner의 새 브라우저 컨텍스트에서 해당 gate를 GET만 한다. 통과는 상태 200, 본문 `상태: expired`, `승인·거절 버튼 비활성: expired`, 버튼 0개, POST 0개다.
4. `adapters/dist/synthetic.js http://127.0.0.1:8080 --seed`로 합성 레코드만 추가한다. 정확한 owner의 새 컨텍스트에서 `/owner`와 fresh gate를 연다.
5. fresh gate에서 typed body의 `granularity_min`과 `2026-10-03T10:00:00Z`, intent `schedule.query`, 정책 `disclosure policy absent`, 만료 시각, 안내 문장, 활성 버튼 2개를 확인한다.
6. `Approve 승인`을 누른다. 통과는 브라우저가 관찰한 POST 303, 상태 `approved`, 버튼 0개, 비활성 문구, 정책과 안내 문장과 typed body 유지, 새로고침 후 같은 상태다.
7. 부모 메시지의 completion이 빈 값인지 확인한다. 클릭이 일정 실행 결과를 만들지 않아야 한다.
8. 다른 seed의 fresh gate를 이전 seed owner 컨텍스트와 같은 seed의 다른 owner 컨텍스트로 연다. 통과는 403 `sender_not_allowed`이고 `granularity_min`과 버튼이 없다.
9. 맞는 owner의 다른 새 컨텍스트로 같은 gate를 연다. 통과는 200과 `상태: pending`이다.
10. `Deny 거절`에 6번과 7번의 조건을 적용한다. 상태는 `denied`다.
11. 3번의 기존 만료 gate가 없을 때만 새 gate를 만들고 출력된 만료 시각까지 실제로 기다린다. 그 뒤 같은 owner로 3번의 읽기 전용 조건을 확인한다. 시계나 라이브 DB 시각은 바꾸지 않는다.
12. finally에서 fixture 바이트와 모드를 복원한다. `relay_state`는 1행으로 남는다. 배포 HEAD와 porcelain은 시작과 같다.

## 판정

| 항목 | 통과 | 실패 |
|---|---|---|
| 승인 | dashboard, fresh 표시, Approve 후 `approved`와 새로고침, POST 303, completion 빈 값 | 버튼이 남거나 completion이 생기거나 상태 200이 아님 |
| 거절 | Deny 후 `denied`와 새로고침, POST 303, completion 빈 값 | 상태가 pending으로 남음 |
| 다른 owner | 이전 owner와 같은 seed의 다른 owner가 403 `sender_not_allowed` | 200 또는 원문 표시 |
| 복구 | 맞는 owner의 새 컨텍스트가 200 | 403이 유지됨 |
| 만료 | 기존 expired gate 또는 실제 대기 후 버튼 0과 POST 0 | 만료 화면에 Approve 버튼이 있음 |
| 복원 | fixture 해시와 모드가 같고 배포 트리가 clean | fixture 교체 또는 배포 소스 변경이 남음 |

## 시나리오 밖

이메일 OTP, 공개 Access 로그인 뒤 UI, 모바일 실기기는 이 시나리오 밖이다. 담당은 사용자 또는 후속 관측이다.
공개 negative `c993d599efab0bfdc5741bc9cb02053ca9afd856`, 로컬 런타임 `1762b430bed1c0584fecd163ae81567a4a5d04a9`, 제품 QA `659f4b06dfda1bd59997f8f2b06cba23b48026a6`, 리뷰 `7ba9df046611c67109b53ad2b43812543b937f95`와 `311381feb0f18203f60329605d762257bfff8421`은 원래 실행 조건으로 재사용한다.
Jev 웹 조작은 이 판정의 도구가 아니다. DOM과 버튼은 `browser.mjs`의 checks가 판정한다.
이 절차의 컨텍스트는 새 Playwright 컨텍스트다. 사용자 브라우저 프로필의 캐시 관측이 아니다.
