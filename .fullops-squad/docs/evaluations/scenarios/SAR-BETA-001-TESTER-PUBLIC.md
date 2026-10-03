---
title: SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 연결의 좁은 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-001-TESTER-PUBLIC]
summary: 최신 Access gate와 공개 negative HTTP의 좁은 QA 절차와 통과 조건을 정의한다
---

# SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 연결의 좁은 QA 시나리오

실행 위치는 tester 워크트리다. gate 대상은 `/tmp/knowslink-beta-review-28bd1bb`의 `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`다.
공개 HTTP 대상은 `https://link.knowslog.com`이다. 상태 디렉터리는 `/home/shin/deploy/knowslink-state`다. 내용은 모드와 일치 여부만 확인하고 값 출력은 하지 않는다.
제품 코드, 배포 소스, Cloudflare 쓰기, DNS 변경, connector 변경은 이 시나리오 밖이다.
기존 로컬 QA `1762b430bed1c0584fecd163ae81567a4a5d04a9`와 제품 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`는 원래 조건으로 재사용한다. `verify.py`, `beta.sh`, `compose.ops.yaml`, tunnel 템플릿 blob이 `f824015314c66bcab42940cfe3db2edabb22e1dd`와 같을 때 백업과 전체 제품 QA를 반복하지 않는다.

## 절차

1. snapshot과 `/home/shin/deploy/knowslink`의 HEAD가 `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`인지 확인한다. porcelain이 0줄인지 확인한다.
2. 제품 경로와 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 diff가 비어 있는지 확인한다.
3. `qa-reports/SAR-BETA-001-TESTER-PUBLIC-test/gates.py`를 `python3 -B`로 실행한다. 종료코드 0이 gate 통과다.
4. gate 결과를 preamble `ask`로 coor에 한 번 제출한다. 회신 전에 앱, DNS, connector를 바꾸지 않는다.
5. coor가 보호된 연결 적용 완료를 회신한 뒤 `public_http.py`를 실행한다. 기본 해석기 NXDOMAIN은 `dig @1.1.1.1`과 `dig @8.8.8.8`의 결과와 따로 기록한다.
6. `PYTHONOPTIMIZE`가 없는 `python3 -B`로 배포 트리의 `verify.py public`과 `verify.py regression`을 실행한다. `shared-baseline.json`은 덮어쓰지 않는다.
7. `proof_check.py`로 보존된 `access.live.json`과 실제 tunnel 설정을 대조한다. 나이가 10분을 넘거나 불일치가 있으면 coor에 읽기 전용 GET 재확인을 요청한다.
8. 인간 이메일 로그인은 실행하지 않는다. held로 남긴다.

## 판정

| 항목 | 통과 | 실패 |
|---|---|---|
| aud 불일치 | 일반과 `-O`, `-OO`, `PYTHONOPTIMIZE=1`의 종료코드가 1이다 | 최적화 실행의 종료코드가 0이다 |
| proof | ID 불일치, missing, 660초, 미래 mtime의 종료코드가 1이다. fresh의 종료코드는 0이다 | 차단 사례가 0을 반환한다 |
| 공개 HTTP | 미인증 `/`, `/owner`, `/v1/registry`, `/healthz`, `/_probe`가 302, 401, 403이다. 302의 호스트는 `cloudflareaccess.com`으로 끝난다. 본문에 `relay`가 없다 | 상태 200 또는 제품 본문 |
| 우회 헤더 | 가짜 JWT, service token, Bearer의 `/healthz`가 같은 거부다 | 헤더만으로 상태 200 |
| HTTP 80 | `/`가 301이고 Location이 `https://link.knowslog.com`이다 | 제품 응답 |
| 공유 서비스 | regression 종료코드 0. orca 200, s8 200, mcp 401. 호스트 cloudflared PID 유지 | 코드 또는 PID 변화 |
| 증명 | 0600 증명과 설정의 ID, aud, team, 단일 allow, IdP 1개, bypass 없음이 일치한다 | 불일치 또는 10분 초과를 통과로 기록 |
| 해석기 | 기본 해석기 실패와 공개 해석기 성공을 서로 다른 결과로 적는다 | NXDOMAIN을 공개 DNS 실패로 적는다 |

## 시나리오 밖

인간 이메일 OTP 로그인, Access 앱 생성, DNS 쓰기, connector 재시작, 백업, 복원, 전체 MVP, 화면 캡처는 이 시나리오 밖이다.
Jev 웹 조작은 이 경계의 판정 도구가 아니다. 판정은 위 스크립트의 종료코드와 상태코드다.
