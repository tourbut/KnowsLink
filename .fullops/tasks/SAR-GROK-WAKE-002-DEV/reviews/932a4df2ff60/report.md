# SAR-GROK-WAKE-002-DEV 리뷰 보고

- base: 1539c6c9594becc18850835ef9d494aadc49bf6b
- head: 932a4df2ff602ff8c372bbef0f088fe3f4067dc4
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-424f5ea9

## 요약

결론 accept, high 0, low 2. 리뷰 초점 D087의 보안 경계 6항목을 코드로 확인했다.

1. SSRF: 등록 때 `routineURLPattern`(고정 호스트 `api2.cursor.sh`, 포트 없음)과 key `^[\x21-\x7e]{1,512}$`(CR/LF 헤더 주입 불가)를 검사한다. 발신 경로는 `deliverRoutine` 하나이고 운영은 `routineClient = remoteHTTPClient`(443만, 공인 IP만, 다이얼마다 DNS 검사, redirect 거부, Proxy nil, 5초)다.
2. key: `Sealer.Seal(key, agent)`로 AAD=agent ID 봉인, `Open`도 같은 AAD. 화면은 URL 끝 4자·상태·마지막 성공만, 실패 화면 폼은 값을 다시 채우지 않는다. 로그 출력 없음, 리다이렉트 `/home/agents?n=routine#routines`에 key 없음. Sealer nil이면 등록 거부(`routine_unavailable`)·폼 숨김·`remoteListen`이 claim 건너뜀.
3. 등록·해제: `memberHit` 세션, `reauthWindow` 5분(아니면 403·`role=alert`·재확인 버튼), consent, 전역 `crossOrigin`(Sec-Fetch-Site/Origin) CSRF, `a.Owner == owner`·미철회·현재 원격 설치 확인. 다른 owner agent 등록·해제 403.
4. 상한·판정: 2분 간격, 시간당 20회(재시도 포함), 2xx 성공 / 0·5xx 30초 뒤 1회 / 보냈지만 응답 없음 무재시도 / 429 창 끝까지 / 401·403·404·410·봉인 실패 Broken / 기타 4xx 버림. 재시도 폭주 경로 없음. 해제 뒤 재등록 카운터 초기화(low 1)와 finish 실패 뒤 1회 중복(low 2)만 남는다.
5. 본문: `routineBody`는 고정 event·대기 수·고정 지시문뿐. 본문·발신자·ID 없음.
6. remote_events.go: `claimRemote` 배달이 먼저이고 Delivery가 없을 때만 routine을 claim, `Next`는 더 이른 쪽만 택한다. remote_auth.go: `Routines` omitempty 필드 추가와 sweep의 routine 정리만이며 기존 Events·OAuth 경로는 바뀌지 않았다(관련 Remote* 테스트 통과).

문서(README·adapters/README·D05·D07·ops-guide)는 코드와 일치한다. 단 low 1의 '다시 등록해도 시간당 횟수는 이어진다'는 해제 뒤 재등록에는 맞지 않는다.

## 검증

- `go vet ./internal/relay/` 통과.
- `go test ./internal/relay/ -run 'Routine|Remote|Human|Agents' -count=1` 통과(Routine 4종, Remote·Human 회귀 포함).
- 통합 `TestRoutineWakePostgres`(build tag integration, 격리 Postgres·Redis)는 lite 레벨이라 다시 돌리지 않았다. 구현자 qa/report.md의 실행 기록을 근거로 둔다.

## 발견 사항

- low `internal/relay/routine_wake.go:112-118` 해제 뒤 재등록이 2분 간격·시간당 20회 카운터를 0부터 시작한다(주인 재확인 필요라 외부 우회 아님). 문서와 다름.
- low `internal/relay/routine_wake.go:191-209` 첫 POST 뒤 finish 트랜잭션 실패·프로세스 중단 시 lease 만료 후 같은 집합을 1회 더 POST해 run이 두 번 시작될 수 있다.
