---
title: 운영·업데이트·복구
status: draft
updated: 2026-11-12
owner: ops
tasks: [SAR-REFACTOR-002-OPS, SAR-RELEASE-001-OPS, SAR-RELEASE-002-OPS, SAR-RELEASE-003-OPS, SAR-PUBLISH-001-DEV, SAR-PUBLISH-001-OPS, SAR-CONNECT-SIMPLE-001-OPS, SAR-AUTOSETUP-001-OPS, SAR-MSG-EPHEMERAL-001, SAR-VULTR-MIGRATE-001]
summary: 현재 배포·백업·writer 보호와 클라이언트 갱신 및 복구 경계를 설명한다.
---

# 운영·업데이트·복구

2026-10-28 운영 제품은 `c6dfc39a`다(정본: [D12 현재 운영 제품](../.fullops-squad/docs/operations/ops-guide.md#sar-autosetup-001-ops-운영-적용-결과-2026-10-28)).
공개 주소는 `https://link.knowslog.com`이며 Dots MCP endpoint는 `/mcp`다.
2026-11-12부터 운영 서버는 Vultr 서울 인스턴스다. 홈서버 스택은 정지 상태로 보존한다([이전 기록](../.fullops-squad/docs/operations/ops-guide.md#sar-vultr-migrate-001-운영-서버-vultr-이전-2026-11-12)).
이 문서는 현재 절차의 설명이다. 명령 실행·계정 발송·배포 승인을 대신하지 않는다.
실제 경로·비밀 설정·공유 자원 목록은 접근이 제한된 운영 자료에서 확인한다.

## 현재 확인 범위

[운영 적용 기록](../.fullops-squad/docs/evaluations/qa-reports/SAR-AUTOSETUP-001-OPS-test/report.md)은 배포·백업·복원 리허설·보존·HTTPS smoke와 자동 설정 0.1.1 게시를 기록한다. 공개 ingress 변경은 없었다.
relay health, schema·HumanModel·writer guard, 기존 회원·agent·친구·Remote 설치·구독 보존을 확인했다.
실제 친구 대화와 질문 없는 화면 알림은 별도 두 계정 시험의 근거다.
그 시험 제품은 과거 `26094e9`다. `bf67dfcb` 운영에서 새 Grok/Claude agent 연결과 시험 메시지 1건 수신을 확인했다([CS-AC-18](../.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-LIVE-test/report.md)).
이번 배포는 직전 `bf67dfcb`의 상태 보존을 확인했고 실계정 로그인·발송은 실행하지 않았다.
D018 캘린더 제거·`KNOWSLINK_CALENDAR_KEY` 삭제·`/privacy`·`/terms` ingress 추가는 SAR-PUBLISH-001-OPS(`7c6ef373`, 2026-10-22)에서 완료했다. 저장된 캘린더 연결은 0이었다.
relay의 이 키 설정과 기동 때 캘린더 정리 코드도 이후 제거했다.
`calendar.freebusy` scope 삭제·Calendar API 비활성화·브랜딩 URL 등록과 앱 게시는 소유자가 Cloud Console에서 직접 한다([인수인계](../.fullops-squad/docs/operations/transition.md)).
새 계정 설치·부하·장시간 알림·다음 Events 갱신·앱 재시작 성공을 보장하지 않는다.
자동 답장·외부 업무 실행은 활성화하지 않았다.

## 배포 준비와 백업

운영 정본은 [상세 운영자 안내](../.fullops-squad/docs/operations/ops-guide.md), [인수인계](../.fullops-squad/docs/operations/transition.md), [DB 설계](../.fullops-squad/docs/design-docs/database-design.md)다.
`deploy/knowslink/`의 도구는 실제 운영 위치와 권한을 확인한 담당자만 사용한다.
과거 문서의 owner-only Access 전제를 현재 회원 경로에 그대로 적용하지 않는다.

1. 독립 리뷰·필수 QA를 통과한 고정 SHA와 제품 diff를 확인한다. 미해결 critical/high가 있으면 배포하지 않는다.
2. 현재 SHA·서비스 health·schema/writer·자격 파일 권한·공유 컨테이너와 ingress 상태를 기록한다. 비밀값은 출력하지 않는다.
3. 같은 고정 SHA의 깨끗한 소스와 새 이미지를 준비한다. 현재 서비스와 DB를 보존한다.
4. 쓰기를 멈춘다. relay와 다른 writer가 정지했는지 확인한다. 안정된 epoch도 대조한다.
5. 현재 DB를 **새 고유 파일**로 `pg_dump -Fc` 백업한다. 폴더 0700·파일 0600을 적용한다.
6. hash·크기·`pg_restore --list`를 기록한다. 네트워크 없는 격리 DB에 복원해 schema·guard·권한 상태를 확인한다.
7. 필요한 migration Up을 적용한다. 새 reader/writer만 기동한다. DB와 친구·철회·generation·Remote/Inbox를 보존한다.
8. health·인증 거부·공개 allowlist·공유 서비스 보존을 검사한다. 승인된 범위의 수락 뒤 트래픽을 재개한다.

최근 배포(`bf67dfcb`→`c6dfc39a`)는 새 migration이 없었다. 백업과 network-none 복원에서도 goose 2·guard 1·HumanModel 1을 확인했다.
일반 `beta.sh deploy <sha>`는 migration 변경을 거부하며 모든 이행을 자동으로 처리하지 않는다.
`beta.sh deploy`, `backup`, `restore-verify`는 내부 삭제·교체·임시 컨테이너 정리를 포함한다.
삭제 금지 상황에서 실행하지 않는다. 최신 운영은 별도 bundle과 수동 보존 절차를 사용했다.
그 수동 실행의 실제 명령·검증·실패 원본은 운영 적용 기록에 있다.

build 전에 현재 실행 이미지에 과제 고유 rollback 태그를 붙이고 `docker image inspect`로 확인한다. 기존 태그를 덮지 않는다.
새 공개 경로를 추가하면 `deploy/knowslink/tunnel/public-ingress.yml`의 `path` 정규식과 `scripts/check_public_ingress.py`를 함께 바꾼다. 둘을 바꾸지 않으면 새 경로가 공개 주소에서 404다.
live Tunnel config는 원본을 백업하고 `cloudflared tunnel --config <후보> ingress validate`로 확인한 뒤 적용한다. cloudflared 재시작은 약1~2초 끊김을 만든다.
공개 주소를 짧은 간격으로 반복 호출하지 않는다. 약2~2.5회/초 관찰 루프가 두 번의 Tunnel 재시작에서 모두 1~2분 동안 429를 만들었다. 재시작 관찰은 공개 주소가 아니라 서버 측 cloudflared 등록 로그로 하고, 공개 확인은 1초 이상 간격으로 한다.

백업에는 키·credential·callback secret이 포함될 수 있다. SAR-MSG-EPHEMERAL-001 이행 뒤 새 백업에는 메시지·Remote 설치 Inbox가 없다(Redis에만 있음). 백업 시점의 본문 없는 gate metadata(ID·digest·대상·상태)와 Dots 구독의 메시지 ID는 들어간다. 이행 전 백업에는 남아 있다.
출력·공유·Git 저장을 금지한다. 접근권한과 보존기한을 별도로 관리한다.
복원 리허설 통과는 실제 rollback 실행이나 WAL의 완전 삭제 증거가 아니다.

## Schema와 writer 보호

`db/migrations/00002_human_writer.sql`은 HumanModel을 낮추는 JSON writer를 거부한다.
Down은 의도적으로 실패하며 guard를 제거하지 않는다.
`db/migrations/00003_relay_state_json.sql`은 `data`(jsonb)를 `state`(json, lz4)로 바꾸고 guard를 `human_model` 숫자 컬럼으로 옮긴다. Down도 실패한다.
이 이행 뒤에는 구 코드가 `data` 컬럼을 찾지 못해 SQL 오류로 멈춘다. 복귀 경로는 이행 전 백업 복원뿐이며 백업 이후 변경을 잃는다. 기본은 전진 수정이다.
절차는 [운영자설명서](../.fullops-squad/docs/operations/ops-guide.md#sar-db-wait-001-저장-형식-이행-운영-미적용)를 따른다. 운영 SQL은 `state`와 `human_model`을 쓴다.
구 writer와 새 writer의 rolling 쓰기를 병행하지 않는다. 코드만 오래된 버전으로 되돌리지 않는다.

신규 schema 이행은 **쓰기 중단·최신 private 백업·격리 복원 리허설·Up·새 writer 검사·재개** 순서다.
검사 대상은 legacy ID·키·철회·정확한 pair와 새 친구 권한의 거부 경계다.
transport 삭제와 1시간 웹 사본 만료(Redis)는 이행 전 백업/WAL의 영구 삭제를 뜻하지 않는다.

복구가 필요하면 다음 경계를 지킨다.

1. writer와 전달을 멈춘다. 실패 시점의 현재 DB도 새 고유 백업으로 보존한다.
2. 같은 schema·HumanModel·새 필드를 보존하는 검토된 reader/writer를 선택한다.
3. 이번 리팩터링의 직전 `2747fe98653d5c6946f41c219fc256395333d7fe`는 같은 writer다. 임의의 구버전에는 이 판단을 적용하지 않는다.
4. 친구·차단·suppression·generation·철회·키·Remote/Inbox가 보존되는지 확인한다. 권한이 불명확하면 read-only/전달 차단을 유지한다.
5. 필요한 재인증·재연결·명시적 새 수락과 독립 확인 뒤 전달을 연다.

과거 백업 전체 복원은 guard로 탐지할 수 없고 폐기한 권한을 되살릴 수 있다.
현재 권한을 과거 DB로 덮지 않는다. 이번 운영 기록은 실제 rollback 성공을 주장하지 않는다.
이번 전환 후 이전 이미지 주소는 조회할 수 없었다. 같은 직전 SHA로 별도 복구 후보 이미지를 다시 만들었다.
원 이미지 동일성을 주장하지 않는다. 복구 후보의 태그·이미지 ID·원 실패는 운영 적용 기록을 따른다.
다음 배포는 build 전에 현재 이미지를 고유 rollback 태그에 고정하고 실제 조회 가능 여부를 확인한다.

## 메시지 임시 저장소 Redis (SAR-MSG-EPHEMERAL-001)

메시지·대화 사본·Dots Inbox는 Postgres가 아니라 디스크 저장 없는 Redis에만 있다(D049~D052). 운영 적용 전이다. 이행 순서는 [운영자설명서](../.fullops-squad/docs/operations/ops-guide.md#sar-msg-ephemeral-001-메시지-휘발-저장-이행-운영-미적용)를 따른다.

- Compose 서비스 `redis`: RDB·AOF 끔, `maxmemory 128mb`·`noeviction`, 포트 미게시, internal network. relay는 `REDIS_URL=redis://redis:6379/0`을 쓴다.
- Redis 재시작은 미전달 메시지·대화 사본·receipt·멱등 기록을 지운다. 이것은 받아들인 동작이다(D051). 재시작 중 메시지 기록이 있거나 새로 쓰는 요청은 `unavailable`, `/healthz`는 503이다. 복구 뒤 새 송신은 정상이다.
- Redis를 백업하지 않는다. `beta.sh backup`은 상태 행에 메시지 key가 남았는지 `message_state_rows`로 보인다(이행 뒤 0).
- 로그·오류 응답에는 본문을 쓰지 않는다. relay 로그는 기동 설정과 고정 오류 문장만 쓴다.

## DB 대기 상한

relay는 `lock_timeout` 1.5초, 유휴 transaction 종료 5초, 슬롯 대기 2초를 쓴다. 값과 이유는 [DB 동시성](system-design.md#db-동시성저장보관)을 따른다.
503 `unavailable`이 이어지면 `pg_stat_activity`에서 lock 대기와 `idle in transaction` 세션을 먼저 확인한다. 멈춘 유휴 세션은 5초 뒤 스스로 끝난다.
잠금을 쥔 세션이 결과를 보내는 중에 상대 프로세스만 멈추면 keepalive로 감지하지 못한다. 이 경우 해당 relay를 재시작한다.

## Grok·Dots 유지보수

설치 본문은 [Grok 기존 업데이트](../README.md#grok-기존-설치-업데이트), [Dots 기존 업데이트](../README.md#dots-기존-연결-업데이트)를 따른다.
Grok은 로컬 Node bundle이다. 서버만 배포하면 실행 중인 MCP·watcher는 바뀌지 않는다.
Inbox 정리 오류의 본문 보존 수정은 로컬 bundle 갱신 뒤 적용된다. 서버 배포만으로 해당 수정이 적용됐다고 판단하지 않는다.
같은 prefix·Command·연결 폴더·key·agent ID·Inbox·webhook routine을 유지한다.
MCP와 정확한 해당 폴더 watcher를 재시작한다. 넓은 process 종료나 신규 connect를 적용하지 않는다.
installer의 신규 폴더 안내는 업데이트 절차가 아니다.

Dots는 기존 Remote HTTPS 연결을 유지한다. 변경된 서버 처리는 다음 remote 요청부터 적용된다.
이번 리팩토링에는 tool·auth·UI·Events metadata 변경이 없으므로 Refresh·재등록·친구 재수락은 필수가 아니다.
이 판단은 diff·서버 보존 관측에 근거한다. Dots cache·다음 자동 갱신을 새로 검증했다는 뜻은 아니다.
향후 metadata가 바뀌면 [공식 Refresh 절차](https://developers.openai.com/plugins/deploy/connect-chatgpt#refresh-metadata)로 기존 연결을 새로 고친다.
Events 목록 변경도 rescan한다. 기존 앱·설치·자동화를 삭제해서 갱신하지 않는다.

## 수신 알림 문제 판별

본문·키·token·callback URL/secret을 로그나 대화에 노출하지 않는다.
진단 자체가 receive·ACK·새 송신을 일으키는지 먼저 확인한다. 상태 조회와 실제 시험을 구분한다.

| 단계·증상 | 확인·다음 조치 |
|---|---|
| `held`·로컬 도구 없음 | Command/Arguments·Node·bundle·mode 확인. 신규 설치만 Google connect 진행 |
| `configured_unverified` | 설정 상태일 뿐이다. 기존 자격·웹 agent·연결 완료 기록을 대조 |
| `queued`·수신 없음 | 상대 process, 현재 친구/자기 권한·키·기한 확인. queued를 delivered로 보고하지 않음 |
| `failed:expired` | 최대 180초 전달 기한 경과. 현재 연결을 복구하고 새 승인 뒤 새 메시지 발송 |
| `autoReceive: off/stopped/backoff` | 자동 수신 환경·private 폴더·process·lastError 확인. 비밀값 출력 금지 |
| 로컬 Inbox pending·화면 없음 | 수신과 저장은 진행됨. host notice와 채팅 턴을 별도로 확인 |
| `hostNotice: sent_unverified` | MCP notice 전송만 의미. 화면 표시·새 턴 성공이 아님 |
| `hostWake: accepted_unverified` | gateway/webhook HTTP 수락만 의미. routine 실행·알림 표시를 따로 확인 |
| Grok 질문해야 보임 | watcher 하나의 생존·정확한 연결 폴더·webhook 파일 권한·routine의 채팅 표시 지시 확인 |
| Dots Events 목록만 있음 | 실제 구독·callback Verified·만료/갱신·해당 대화 자동화 지시 확인 |
| Dots delivered인데 화면 없음 | 서버 Inbox/ACK와 host 턴·표시 지연을 분리. receive 성공만으로 UI 성공 판정 금지 |
| 중복 우려·불확실한 wake | `uncertain`·`rejected`는 자동 재발송하지 않음. 기존 ID·marker 보존 |

Grok webhook routine은 고정 metadata doorbell 뒤 `knowslink_text_receive`로 본문을 읽는다.
받은 본문을 사용자 권한 prompt나 외부 도구 명령으로 넣지 않는다.
`KNOWSLINK_GROK_WEBHOOK_FILE`은 URL/key 자체가 아니라 소유자 전용 비공개 파일 경로다.
webhook 모드의 `plugin.js wake-check`는 파일 검사이며 발송하지 않는다.
routine이 알림을 숨기도록 설정됐으면 수신이 성공해도 채팅에 보이지 않을 수 있다.

Dots는 `knowslink.text.arrived` 구독 후 receive 결과를 실제 대화에 표시하고 성공 뒤 ACK하도록 명시한다.
현재 구독 저장은 서버 재시작 전후 확인했다. 다음 host 갱신과 장시간 알림은 미검증이다.
수신·ACK가 자동 답장 승인을 대신하지 않는다. 별도 승인 없는 새 시험 메시지도 보내지 않는다.

## 보호 경로와 인계

최근 공개 smoke에서 `/`·OAuth metadata는 성공하고 무인증 `/mcp`·SSE는 401이었다.
관리·시험 정확 경로는 404였다. local health는 공개 health 노출과 구분한다.
다른 공유 서비스·Tunnel·Google 설정을 보존한 근거는 운영 적용 기록에 있다.
과거 검증 도구의 Access 302 전제나 오래된 baseline과 현재 보호 계약이 다르면 실제 경로 기준으로 보완한다.

후속 담당에게 현재 SHA·검사 명령/종료코드·실패·미실행·새 백업·복구 경계·남은 자료 위치를 넘긴다.
장시간·부하·다음 Events 갱신과 실제 앱 재시작 검증은 승인된 별도 시험으로 남긴다.
