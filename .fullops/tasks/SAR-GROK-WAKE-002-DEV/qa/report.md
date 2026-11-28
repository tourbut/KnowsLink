---
id: qa-SAR-GROK-WAKE-002-DEV
title: SAR-GROK-WAKE-002-DEV QA 보고
status: approved
updated: 2026-11-28
owner: dev
summary: relay 직접 Grok Bot routine wake의 표적 단위 테스트와 격리 relay 통합 1회 결과
---

# SAR-GROK-WAKE-002-DEV QA 보고

테스트 레벨 lite. 표적 단위 테스트와 격리 Postgres·Redis 통합 1회로 확인했다. 운영 배포, 운영 `.env` 변경, DB 마이그레이션, 실제 Grok 계정 시험은 하지 않았다.

## 구현 요약

- `internal/relay/routine_wake.go`: 등록 저장(`Remote.Routines[agentID]`, key는 `KNOWSLINK_SEAL_KEY` 봉인, AAD=agent ID), 도착 감지(`eventPending`), 상한, 결과 판정, 재시도, POST.
- 배달은 기존 `RemoteEvents` LISTEN 루프에서 MCP Events 배달 다음에 돈다. seal 키가 없으면 claim하지 않는다.
- 웹 `/home/agents#routines`: 원격 agent별 등록(`POST /home/routine-set`)·해제(`POST /home/routine-remove`). ko·en.
- 문서: `README.md`, `adapters/README.md`, D05 `interface-design.md`, D07 `database-design.md`, D12 `ops-guide.md`.

## 결정(추천안으로 진행, 되돌릴 수 있음)

| 결정 | 선택 | 근거 | 되돌리는 법 |
|---|---|---|---|
| 저장 위치 | `remoteState.Routines` JSON 필드(additive optional) | SQL migration 불필요. 원격 설치가 사라지면 `sweepRemote`가 함께 지운다 | 필드 삭제. 이전 코드는 이 필드를 버린다(재등록 필요) |
| 연결 실패와 timeout 구분 | `httptrace.WroteRequest` 성공 여부로 판단. 요청이 나가기 전 실패만 '연결 실패'(30초 뒤 1회 재시도). 나간 뒤 응답 없음은 재시도 안 함 | analysis의 '시작 여부 불명이면 재시도 안 함'을 timeout 외 끊김에도 적용해 중복 run을 막는다 | `routinePost` |
| 2분 간격 | 새 알림(첫 시도) 사이 2분. 재시도(30초)는 같은 알림의 일부다. 시간당 20회는 재시도를 포함한 모든 POST를 센다 | 주인 상한을 넘지 않으면서 analysis의 30초 재시도를 지킨다 | 상수 3개 |
| 다시 울리는 조건 | 대기 집합에 아직 알리지 않은 메시지 ID가 있을 때만. 줄기만 한 집합(일부 수신)은 울리지 않는다 | `claimRemote`의 `Seen`/`Pending`을 따르되, 수신 중에 다시 깨우는 일을 막는다 | `routineFresh` |
| 429 | 현재 시간 창 끝까지 멈춘 뒤 아직 대기 중이면 한 번 울린다 | '다음 시간 창까지 멈춤' | `finishRoutine` |
| 기타 4xx(400·413 등) | 이 집합은 포기, 등록은 유지 | run이 시작되지 않았고 반복해도 같은 결과 | `finishRoutine` |
| 재등록 | URL·key·상태를 새로 쓰되 시간당 계수는 유지 | 재등록으로 상한을 비우지 못하게 한다 | `setRoutine` |
| 실패 표시 | 재인증 없음(403)·형식 오류(422)·동의 없음은 '내 에이전트' 화면을 그대로 다시 그리고 맨 위에 `role="alert"` 오류를 보인다. 재인증 전이면 절 안에 '로그인 다시 확인' 버튼을 미리 보인다 | grant-create가 일반 거절 화면으로 넘어가는 문제(SAR-ORCA-GRANT-UX-001-DEV)를 이 화면에서는 피한다 | 핸들러 한 case |
| 키 없을 때 화면 | 폼을 숨기고 '암호화 키가 설정되지 않아 등록할 수 없습니다'를 보인다. 서버도 `routine_unavailable`로 거절 | 평문 저장 금지 | 템플릿 조건 |

## 완료 기준 대응

| 기준 | 테스트 |
|---|---|
| URL 형식 거부 | `TestRoutineRegistrationRefusesBadInputAndSealsKey`: http, 다른 호스트·접미 호스트, 빈 ID, `..`, query, userinfo, 다른 포트, 129자 ID 거부. 통합: `https://127.0.0.1/...` 등록이 422 + 오류 문구 |
| SSRF 거부 | `TestRoutinePostHeadersBodyAndOutcomes`: 운영 클라이언트(`remoteHTTPClient`)로 127.0.0.1·10.0.0.1·::1·169.254.169.254 POST가 아무것도 보내지 않고 0 |
| 재인증 없는 등록 거부 | 통합 `TestRoutineWakePostgres`: 6분 전 확인 세션의 등록·해제가 403, `role="alert"`, 재인증 버튼, 저장 없음, 응답에 key 없음 |
| 상한(2분, 시간당 20) | `TestRoutineSpacingHourlyCapAndNoRepeat`: 같은 집합 재울림 없음, 30초 뒤 새 도착은 2분까지 대기 후 묶어 1회(pending 2), 줄어든 집합 무시, 21번째 POST는 다음 시간 창까지 대기 |
| 판정별 동작 | `TestRoutineResultHandling`: 200·503×2·연결 실패→200·응답 없음·400·429·401·403·404·410·키 분실 각각의 POST 수, 다시 등록 필요, 429 재개, 재등록 재개 |
| 봉인 왕복, key 비노출 | 봉인 값이 같은 agent로만 열림, 상태 JSON·화면 view에 key·전체 URL 없음(끝 4자만). 통합 화면에도 key·전체 URL 없음 |
| 본문에 본문·발신자 없음 | 단위·통합 모두 본문 키가 정확히 `event,next,pending`, 메시지 본문·발신자 agent ID 없음. 헤더 `Authorization: Bearer`·`X-Automation-Key` 확인 |
| 격리 relay 통합 1회 | `TestRoutineWakePostgres`: 웹 로그인 → 원격 agent 등록 → 도착 1건 → 실제 `RemoteEvents` 루프 → 가짜 routine 서버가 POST 정확히 1회(1.5초 더 기다려도 1회), 경로·Host·헤더·본문 확인, 화면 '마지막 성공', 해제 |

## 실행 기록

```
gofmt -l cmd internal                                  # 출력 없음
go vet ./internal/... ./cmd/... (+ 통합·fixture 태그)   # 통과
go test -count=1 ./internal/... ./cmd/...              # ok (relay 3.4s)
# 격리 컨테이너 postgres:17-alpine·redis:8.10.2-alpine(kl-grokwake-*), go run ./cmd/migrate 뒤
TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL=<격리> TEST_REDIS_URL=<격리> \
  go test -count=1 -tags=integration -run 'TestRoutineWakePostgres|TestRemote|TestHuman' ./internal/relay/   # ok 6.8s
```

컨테이너와 받은 redis 이미지는 실행 뒤 삭제했다. 통합 테스트는 가짜 routine 서버로 가는 다이얼만 바꾼다(`routineClient`, 테스트 전용). 저장 URL은 실제 고정 호스트 그대로다.

## 하지 않은 것

- 운영 배포. 새 env는 없다(`KNOWSLINK_SEAL_KEY`는 65a752fe부터 운영에 있음).
- 배포 뒤 주인 단계: 웹에서 routine 등록, 실제 도착 알림 1회 확인(원격 routine run이 `knowslink_receive`를 부르는지).
- designer 시각 검수(지시서대로 배정하지 않음). 화면은 grant 절과 같은 구조·클래스만 썼다.
