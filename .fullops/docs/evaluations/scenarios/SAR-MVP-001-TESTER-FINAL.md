---
title: SAR-MVP-001-TESTER-FINAL — RF-01 좁은 독립 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: "legacy claim 차단과 정상 agent, owner, 현재 권한의 좁은 재검증 절차를 정의한다"
---

# SAR-MVP-001-TESTER-FINAL — RF-01 좁은 독립 QA 시나리오

제품 대상은 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`다. 이 트리의 제품 경로는 그 커밋과 같아야 한다.
실행 위치는 tester 워크트리 `/home/shin/orca/workspaces/KnowsLink/fullops-tester`다.
명령은 기존 `run.py`로 실행하고 종료코드를 그대로 남긴다.
합성 credential, claim token, CSRF 값은 로그에 남기지 않는다. Tunnel과 운영 DB는 사용하지 않는다.
`4262d02` 보고서와 test 증거, 원래 held, UI PNG는 원래 SHA와 조건으로 재사용한다. 파일을 덮어쓰지 않는다.

## 절차

1. `78b1d92`와 HEAD의 제품 경로 diff가 비어 있는지 확인한다.
2. `internal/relay/testdata/legacy_claims.json`을 읽는다. fixture를 다시 생성하지 않는다.
3. 고유 Compose project에서 relay를 띄운다. Postgres 비밀번호와 Tunnel token은 로그에 남기지 않는다.
4. fixture의 `Receipt.exp`와 `Gate.Exp`만 현재 시각보다 뒤로 옮긴 뒤 `relay_state`에 로드한다. `Deliver`와 claim token은 유지한다.
5. 같은 SQL에서 `Deliver` 부재, delivered claim, 승인 gate 미소비를 확인한다.
6. human claim과 경로 미기록 agent claim으로 authorize, H, R을 호출한다. 통과는 `403 sender_not_allowed`다.
7. 승인 legacy gate를 parent claim으로 consume한다. 통과는 `403 sender_not_allowed`이고 gate는 미소비다.
8. 같은 상태에서 새 agent 메시지를 send, pull, persist, ACK, claim, authorize한다. 통과는 `200`이다.
9. 그 부모로 owner gate를 열고 잘못된 CSRF를 거절한 뒤 approve, consume, denied result를 확인한다.
10. agent credential로 owner-revoke를 호출한다. 통과는 `401 invalid_auth`다.
11. 해당 project만 `down --volumes` 한다.

## 판정

| 항목 | 통과 | 실패 |
|---|---|---|
| legacy human, unrouted | authorize, H, R이 `403 sender_not_allowed` | `200` 또는 claim, completion, gate가 변함 |
| legacy consume | `403 sender_not_allowed`, 미소비 | consume `200` |
| 새 agent, owner, 현재 권한 | 위 8–10의 상태 코드 | 그 경계의 다른 상태 코드 |
| 환경 실패 | 제품 판정으로 쓰지 않는다 | 종료코드와 원인을 로그에 남긴다 |

전체 MVP, `make verify-mvp`, 새 화면 캡처, 원래 held의 해소는 이 시나리오 밖이다.
