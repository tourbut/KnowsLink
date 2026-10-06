---
title: SAR-PUBLIC-MESSAGES-001-TESTER 프로브 중간 실패
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-TESTER]
summary: 17차까지 프로브 기대값 수정과 남은 제품 결함을 구분한다
---

# 프로브 중간 실패

후보는 `09c523da8a3407288d9f5d711e1834af12bc7808`다. 제품 코드는 수정하지 않았다. 아래 1–7은 프로브 기대값이다. 8은 제품 결함이며 17차에도 남았다.

| 차수 | 종료 | 구분 | 내용 |
|---|---|---|---|
| 8 | 1 | 프로브 | 프로브 시계로 TTL을 계산해 180초가 422였다. DB `clock_timestamp()`로 고친 뒤 200과 422가 맞았다. |
| 9 | 1 | 프로브 | receipt filler의 `Accepted`가 0이라 다음 transaction이 지웠다. 20000 다음은 200이었다. `Accepted`를 채운 뒤 409다. |
| 10 | 1 | 프로브 | 이미 남은 HTTP 행 때문에 신규 입장이 바로 429였다. 선택 대기가 그 429를 멈춤으로 보았다. |
| 11 | 1 | 둘 다 | 남은 신규 1개를 `t.Error`로 기록한 뒤, 해제 후 그 행을 새 누수로 `t.Fatal`했다. 이후에는 검사 전 개수와 같으면 통과한다. |
| 12 | 1 | 프로브 | persist 없이 ACK해서 `invalid_lease` 409였다. persist를 슬롯을 심기 전에 실행하면 ACK는 200이다. |
| 13–14 | 1 | 프로브 | pipe 본문이 슬롯을 점유하기 전에 조회 IP가 `rate_limited` 429를 받았다. `capacity`만 점유로 인정하고, 막힌 reader를 쓴다. |
| 15 | 1 | 프로브 | unpair 뒤 sweep가 denied gate를 revoked로 바꿨다. deny 직후와 unpair 직후를 나눠 확인한다. |
| 16 | 1 | 둘 다 | H-1, M-1, 슬롯 잔류는 유지됐다. 마지막 GET receipt를 303으로 기대해 429에서 멈췄다. GET은 신규 작업이라 429가 맞다. |
| 17 | 1 | 제품 | `t.Fatal`은 없다. 실패는 슬롯 잔류 1개, H-1 신규 16, M-1 form deny, H-1 정리 4다. |

17차 로그는 `probe-go-attempt17.log`다. 결함 문장은 `admission-defects.json`이다.
