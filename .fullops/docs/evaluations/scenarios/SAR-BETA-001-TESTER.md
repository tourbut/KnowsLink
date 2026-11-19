---
title: SAR-BETA-001-TESTER — 베타 로컬 런타임 좁은 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-001-TESTER]
summary: 베타 로컬 런타임의 좁은 QA 절차와 통과 조건을 정의한다
---

# SAR-BETA-001-TESTER — 베타 로컬 런타임 좁은 QA 시나리오

실행 위치는 tester 워크트리다. 대상은 `/home/shin/deploy/knowslink`의 현재 detached 체크아웃과 `/home/shin/deploy/knowslink-state`다.
명령은 `qa-reports/SAR-BETA-001-TESTER-test/run.py`다. 종료코드 0이 판정 통과다.
credential, CSRF, API 토큰, 사용자 이메일은 로그에 쓰지 않는다. `owner-login`과 `seed`는 실행하지 않는다.
제품 QA, 시각 판정, 최종 리뷰는 원래 SHA와 조건으로 재사용한다. 파일을 덮어쓰지 않는다.

## 절차

1. 배포 HEAD를 기록한다. `437f1432a158670a485413c1aba159debc3759e5`가 조상인지 확인한다. 판정 실행 동안 HEAD가 바뀌면 그 실행은 판정이 아니다.
2. 제품 경로와 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 diff가 비어 있는지 확인한다.
3. relay 포트가 `127.0.0.1:8080`인지 확인한다. Postgres 호스트 매핑이 없는지 확인한다.
4. postgres·relay의 재시작과 메모리 제한, relay의 `cap_drop`·`read_only`·`no-new-privileges`를 확인한다.
5. `.env`와 tunnel 자격 파일이 0600이고 상태 디렉터리가 0700인지 확인한다. 파일 내용은 읽지 않는다.
6. `verify.py local`과 `verify.py regression`을 실행한다. `shared-baseline.json`은 덮어쓰지 않는다.
7. 합성 fixture로 음성만 호출한다. 인증 없는 owner는 401이다. agent credential의 owner 경로와 owner-revoke는 401이다. 잘못된 CSRF는 403이고 gate 상태는 그대로다.
8. `beta.sh backup`으로 새 0600 dump를 만든다. 기존 dump는 남긴다.
9. `beta.sh restore-verify`로 네트워크 없는 임시 컨테이너에 복원한다. 라이브 DB와 `knowslink_postgres-data`는 유지한다. 임시 컨테이너만 제거한다.
10. `access.aud`가 없을 때만 `beta.sh expose`를 실행한다. 통과는 종료코드 1과 `Access app not recorded`다. DNS와 connector는 생기지 않아야 한다.
11. `verify.py regression`을 다시 실행한다. 배포 체크아웃은 clean이어야 한다.

## 판정

| 항목 | 통과 | 실패 |
|---|---|---|
| loopback·Postgres | relay는 `127.0.0.1:8080`만 연다. Postgres는 호스트에 게시되지 않는다 | `0.0.0.0` 또는 Postgres 호스트 매핑 |
| 비밀·제한 | `.env`와 자격 파일 0600. 재시작과 메모리 제한이 구성과 같다 | 모드가 느슨하거나 제한이 없다 |
| 합성 음성 | 401 `invalid_auth`, 403 `invalid_csrf`. owner는 활성으로 남는다 | 200 또는 303으로 결정이 저장된다 |
| 백업·복원 | dump 0600. 임시 복원 종료코드 0. 라이브 health 200. 볼륨 유지 | 라이브 DB 교체 또는 볼륨 삭제 |
| expose | Access 기록이 없으면 DNS 전에 종료코드 1 | DNS 생성 또는 connector 기동 |
| 공유 서비스 | regression 종료코드 0 | 코드, PID, myportfolio 상태 변화 |

## 추가 게이트

12. `gates.py`를 실행한다. Access 증명이 없거나 10분이 지났거나 aud·이메일·IdP가 맞지 않으면 0이 아닌 종료로 끝난다.
13. 임시 상태의 `expose`는 DNS 명령 전에 종료코드 1이다.
14. 임시 clone에서 없는 SHA와 migration이 다른 SHA의 `deploy`는 backup 전에 종료코드 1이다. 라이브 체크아웃은 바꾸지 않는다.

## 시나리오 밖

`verify.py public`, 본인 이메일 로그인, 원래 앱 Basic의 인간 확인은 보호 연결 뒤의 별도 QA와 인간 검사다.
전체 MVP, `make verify-mvp`, 새 화면 캡처, DEC-02, DEC-03, 실adapter, 실제 벤더, WAL 삭제는 이 시나리오 밖이다.
Jev 웹 조작은 이 경계의 판정 도구가 아니다. 판정은 `run.py`의 checks다.
