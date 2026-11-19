---
title: SAR-MVP-001-TESTER — 첫 안전 전달 기능의 독립 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER]
summary: 고정 합성 후보의 QA-01–11 절차와 통과 조건 및 held 경계를 정의한다
---

# SAR-MVP-001-TESTER — 첫 안전 전달 기능의 독립 QA 시나리오

제품 대상은 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`다. 준비 커밋은 지시서와 route만 바꾼다. 제품 경로는 그 후보와 같아야 한다.
실행 위치는 tester 워크트리 `/home/shin/orca/workspaces/KnowsLink/fullops-tester`다. 명령은 `run.py`로 실행하고 종료코드를 그대로 남긴다.
합성 credential과 private key는 로그·문서·캡처에 남기지 않는다. 외부 발송·운영 DB·Tunnel은 사용하지 않는다.
초기 골격 시나리오의 404·빈 migration 기대값은 이 기능의 통과 조건이 아니다.

## 공통 절차

1. `git diff a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92 --` 제품 경로의 출력이 비어 있는지 확인한다.
2. `npm ci --prefix adapters`로 lockfile 의존성을 맞춘 뒤 `make lint`, `make test`, `make build`를 실행한다.
3. `make verify-mvp`로 고유 Compose project의 migration·Postgres 경합·Go UI·TypeScript stub을 실행한다.
4. `isolated.py`가 다른 고유 project를 띄우고 `probe.mjs`로 QA-01–11의 HTTP·SQL 시간 제어·화면 HTML을 실행한다. 종료 시 그 project만 `down --volumes` 한다.
5. 저장된 HTML을 headless Chrome으로 PNG에 찍는다. CSRF 값은 저장 전에 지운다.
6. DEC-02·DEC-03·Free N·실벤더·WAL 완전 삭제는 held로 남긴다. 합성 통과로 바꾸지 않는다.

## 판정

| ID | 통과 조건 | 실행 |
|---|---|---|
| QA-01 | owner 가입·PoP·rotate가 되고, agent credential의 owner 동작·URL kid·이전 kid 재할당·이전 key 서명이 거부된다. | probe |
| QA-02 | 수락 전 전달 거부, deny, pending의 active 비포함, 동시 accept 1건, active 재초대의 같은 세대. Free N은 추정하지 않는다. | probe |
| QA-03 | 중복 키·deliver:both·unknown·A2A 필드·잘못된 Unicode·대문자 from 거부. 인증 없는 재전송은 receipt digest를 돌려주지 않는다. | `make test`의 frozen parse와 probe |
| QA-04 | 같은 key+digest는 receipt만, 다른 digest는 409, 301초와 id 충돌은 reservation을 남기지 않는다. | `make verify-mvp`와 probe |
| QA-05 | lease 30초와 exp 최소, 성공 grant만 attempts 증가, 3회 window, 늦은 ACK·claim 1회·재시작 재claim 거부, max attempts. | probe와 integration |
| QA-06 | key revoke 확정 후 기존 대기는 임대되지 않는다. 교체 key도 그 대기를 복구하지 않는다. 교체 뒤 새 전송은 임대된다. unpair 중 ACK는 거부되고 재수락은 새 세대다. 이전 idempotency replay는 거부된다. DB 시계 이상은 503이고 없는 credential은 401이다. 고의 stale epoch는 held다. | probe |
| QA-07 | 무정책 query의 disclosure/executable false, commit stub 비실행, done 거부, H의 claim·digest·중복·재귀·exp 결속. | probe와 synthetic |
| QA-08 | 검증 body·정책 표시, hint 비사용, GET 비승인, agent 401, Basic 인증, CSRF, 원자적 approve/deny, consume 1회, disclosure false. | probe |
| QA-09 | 방향·optional result·stack 거부, denied 1회, 같은 parent의 두 번째 result 거부, 결과 pull이 추가 명령을 만들지 않음. transport 만료는 서명 결과가 아니다. | probe |
| QA-10 | evidence 수신이 fetch가 아님, high가 앞 메시지를 넘지 않음, webhook 404, registry에 A2A 승인 대체 없음, 미설정 adapter의 webhook/evidence false. 실연결은 held. | probe |
| QA-11 | exp가 원문만 지우고 receipt는 남긴다. 저장 시각 25시간은 지우고 23시간은 남긴다. 벽시계 24시간을 기다리지 않는다. expired/revoked/unavailable/401 문구를 캡처한다. | probe |
| V-01–04 | pending, approved, denied, expired, revoked, 원문 부재, 권한 불명 화면을 같은 후보에서 저장한다. designer 판정은 별도다. | probe HTML·PNG |

## 시나리오 밖

- 실데이터 positive silent done, 공개 한도, 실제 신원 인증, singleton 처리량, 벤더 연결, WAL/backup 완전 삭제는 held다.
- 제품 코드 수정·운영 배포·실데이터는 하지 않는다.
- 직접 시각 수락은 designer가 같은 후보의 캡처로 판정한다.
