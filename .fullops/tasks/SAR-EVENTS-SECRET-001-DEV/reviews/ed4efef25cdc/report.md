# SAR-EVENTS-SECRET-001-DEV 리뷰 보고

- base: 4520d88ec290ff2203b7c587ab94051bf2f243bb
- head: ed4efef25cdc6090a23e7badc3d6a0df34f3cade
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-f06015b6
- reviewer session: 5cd78d16-01e6-4b1f-b73d-69dbc9d74ff9 (구현 세션과 별개)

## 요약

lead 리뷰 초점(D083) ①~⑤를 코드로 따라가며 확인했다. 제품 파일 9개와 문서 4개를 열람했고, 하네스 기록 2개는 preset대로 건너뛰었다.

- ① `seal.go`: `Seal`마다 `crypto/rand`로 GCM nonce 12바이트를 새로 만든다. 재사용 경로가 없다. 재봉인은 구독 refresh마다(구독 128개, 시간당 1회 수준)라서 무작위 nonce 한도(2^32)와 거리가 멀다. AAD는 구독 ID(`hash(install‖URL)`)다. 키는 StdEncoding base64 32바이트만 받는다. 오류 문구·기동 로그(`configured=%t`)·구독 응답·`fault` 코드에 키나 평문이 없다.
- ② 키 없음: `NewSealer("")`는 nil을 돌려준다. `events/subscribe`는 secret 검사·상태 쓰기 전에 `events_webhook_unavailable`(-32602)로 거절한다. 그래서 평문을 쓰는 경로가 없다. unsubscribe는 키 검사 앞에 있어 키가 없어도 동작한다. 잘못된 키는 `run`이 오류를 돌려줘 기동이 실패한다. D05·D12 서술과 같다.
- ③ 평문 이행: 상태의 모든 쓰기 경로(최종 트랜잭션의 `sub.Secret`, `OldSecret`)가 `Sealer.Seal`을 거친다. 봉인 값은 `sealed:v1:`로 시작하므로 `whsec_` fallback과 겹치지 않는다. 그래서 공격자가 요청으로 평문 저장 값을 만들 수 없다. fallback은 배포 전 레코드에만 닿는다. lease 단계 새 구독은 `Secret==""`이고, 검증 실패 시 삭제된다. `OldSecret`은 `sweepRemote`가 1분 뒤 비운다(remote_auth.go:120). 그래서 "약 1시간 뒤 평문 없음"이라는 D07 서술이 성립한다. 키 분실 시 배달은 410처럼 구독을 지운다.
- ④ Beta: `extensions["io.modelcontextprotocol/events"]` 광고는 기존 `capabilities.events` 옆에 추가만 했다. 선언 판정(`remoteDeclaresEvents`)은 modern 요청에서만 Beta로 본다. `events/poll`은 기존 `remotePrincipal` 인증 트랜잭션 안에서 자기 설치의 agent 대기 집합만 읽는다. 외부 호출이 없어 SSRF 면이 없다. `Mcp-Name` 비교 대상에 추가했다. subscribe의 URL 검사·callback 검증(SSRF 경계)은 바뀌지 않았다. webhook·poll 본문은 `{eventId,name,timestamp,data:{arrival:true}}`(+webhook `cursor:null`)뿐이다. cursor·eventId는 대기 ID 집합의 해시라서 본문·발신자·메시지 ID를 드러내지 않는다.
- ⑤ Dots: 미선언 요청의 `events/list`는 `["webhook"]`, `webhook-id`=`eventId`, 검증 실패 -32015, 본문 형태가 그대로다. 평문 구독의 OpenAI 형식 refresh는 성공하고 재봉인된다(integration `remoteLegacySecretReseal`).

## 검증

- 다시 돌린 표적 단위 테스트: `go test ./internal/relay/ -run 'Seal|RemoteSecret|RemoteSubscribe|RemoteEvents|RemotePoll' -count=1 -v`의 5개가 모두 PASS였다(`TestSealerRoundTripAndRefusals`, `TestRemoteSecretMigrationReadsPlainTextUntilResealed`, `TestRemoteSubscribeRefusedWithoutSealKey`, `TestRemoteEventsDualShapeAndMetadataOnlyBody`, `TestRemotePollIsStatelessAndMetadataOnly`).
- integration(`-tags=integration TestRemote*`, 격리 Postgres·Redis)은 다시 돌리지 않았다. QA 보고의 통과 기록과 테스트 diff 내용으로 확인했다.
- D12 이행 확인 SQL: `relay_state.state`가 `json` 열(00003 migration)이고 필드명이 `Remote`/`Subscriptions`/`Secret`(태그 없음)이라서 `json_each` 쿼리가 맞다.

## 발견 사항

없음(high 0, medium 0, low 0).

참고(결함 아님): 키가 없을 때도 미선언(Dots) 요청의 `events/list`는 Dots 형태를 지키려고 `["webhook"]`을 광고한다. 이때 subscribe는 거절된다. D12가 "키를 먼저 넣는다"를 굵게 적어 운영 순서로 막고 있다. Beta `events/poll` 요청 필드가 SEP-3415 원문과 일치하는지는 원문을 다시 대조하지 않았다.

## 결론

accept. 운영 배포 전 `KNOWSLINK_SEAL_KEY` 주입(ops 단계)이 남아 있다.
