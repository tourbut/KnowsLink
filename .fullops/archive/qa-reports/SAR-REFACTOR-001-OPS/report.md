---
title: 수락 리팩토링 운영 배포·보존 검증 및 클라이언트 안내
status: draft
updated: 2026-10-13
owner: ops
tasks: [SAR-REFACTOR-001-OPS]
summary: 2747fe9 운영 적용과 백업·writer guard·회원·친구·Events 보존 및 Grok/Dots 후속을 기록한다
---

# SAR-REFACTOR-001-OPS 운영 적용 결과

수락 제품 `2747fe98653d5c6946f41c219fc256395333d7fe`를 기존 운영 서버에 적용했다. 회원·agent·친구·키·Remote Inbox·Events·Google 설정·writer guard·Tunnel·공유 서비스의 보존 대조가 통과했다. 제품 코드와 DB schema는 변경하지 않았다. 클라이언트 업데이트는 안내만 작성했다.

## 대상과 권한

- 사용자 운영 배포 명시 승인과 역할 인박스의 명시 배포 지시를 적용했다. 기준은 `2747fe98653d5c6946f41c219fc256395333d7fe`다. 준비 HEAD는 `a879c54`이며 역할 브랜치는 `fullops/ops`다.
- 적용 기준은 `fullops-common-0.3.3`, AGENTS.md, FULLOPS.md, project.md, D07, D12, D13이다. 테스트 레벨은 lite다. 지시대로 직접 수행했고 하위 위임은 하지 않았다.
- Dispatch는 `ctx_83153b8b8e1b`, Task는 `task_7aa1eabd0631`, Run은 `run_5534969da2d0`다. 송신 핸들은 `term_a25928cb-f090-435d-96ce-eb5433c1b77f`다.
- [CODE 수락](../SAR-REFACTOR-001-CODE-review/report.md), [FINAL 수락](../SAR-REFACTOR-001-FINAL-review/report.md), [독립 QA](../SAR-REFACTOR-001-TESTER-test/report.md)를 읽었다. 고정 제품 `cd4dae8`의 기존 Linux 등록 게이트를 재사용한다. 이번 문서 HEAD의 새 제품 테스트 실행으로 표시하지 않는다.

## 배포와 백업

| 항목 | 결과 |
|---|---|
| 이전 제품 | `26094e90648aa71e8d07b381e9891acb37438730`, detached·clean·relay healthy |
| 실제 적용 제품 | `2747fe98653d5c6946f41c219fc256395333d7fe`, detached·clean |
| 제품 diff | Go/TS 코드·회귀 11개 파일. `db`, `deploy`, Compose, Dockerfile, go.mod/go.sum, Makefile, sqlc.yaml diff 0 |
| bundle | `main ^26094e9`, SHA256 `4753f5cc05bca30e571782d1a8f10285bcdea003216904182e08fbd990e5bb0c`. 서버 hash·bundle verify·FETCH_HEAD·이전 SHA 조상 관계 확인 |
| 빌드 | 고정 checkout에서 relay/migrate build EXIT=0. 빌드 동안 기존 relay는 가동했다 |
| 중단 구간 | STOP_BEGIN `2026-10-13T13:21:02Z`, UP_DONE `13:21:18Z`. 약 16초 뒤 healthy 확인 |
| 쓰기 고정 | relay만 stop. epoch `912460`이 2초 간격으로 같았다 |
| 새 백업 | `/home/shin/deploy/knowslink-state/backups/26094e9-pre-refactor-20261013T132105Z.dump`, 0600, 15693 bytes |
| 백업 SHA256 | `792523eebdb6d007fd6f01cf45ebc75b282942dd96eabaf60424602b74cc95bd` |
| 백업 검사 | `pg_restore --list` EXIT=0, TOC 25줄. network none 복원 EXIT=0 |
| 격리 복원 | goose 2·guard 1·relay_state 1행·epoch 912460·HumanModel 1. HumanModel 제거 쓰기는 `human authorization writer downgrade refused`로 거부 |
| 기동 | `beta.sh up` EXIT=0. migrate 종료 0이며 적용할 새 migration은 없었다. `verify.py local` EXIT=0 |
| relay 이미지 | `sha256:220edb718b2fe1925eed243084e41fedf641c86a0342760f682672b91379ed3d`, healthy·restart 0 |
| 실행 로그 | [build.log](build.log), [deploy-verified.log](deploy-verified.log), [post-check.log](post-check.log) |

기존 helper `AppData/Local/KnowsLinkDevTools/humanpeer001-ops/remote.py`의 known_hosts·RejectPolicy를 유지했다. 자격은 메모리에서만 읽었다. 서버에 없는 origin과 `beta.sh deploy`·`backup` 내부 rm 때문에 무삭제 bundle→checkout→build→relay stop→고유 backup/list→격리 복원→up 순서를 재사용했다. `deploy-history.log`에 수동 적용을 append했다. prepare/seed/owner-login/restore-verify/expose/Access·DNS·Tunnel 변경 명령은 실행하지 않았다.

복원 컨테이너 `knowslink-refactor001-rehearse-20261013`은 `--network none`, postgres 사용자, 호스트 포트 없음으로 실행했다. 검증 후 Postgres를 정상 종료했고 컨테이너도 exited다. 삭제하지 않았으므로 복원된 비공개 데이터가 Docker 관리 권한 안에 남는다. 백업과 기존 임시자료도 보존했다.

## 보존 대조

전후 관측은 `13:21:01Z`와 `13:21:19Z`다. [before.json](before.json), [after.json](after.json), [preservation.json](preservation.json)에 개인값 없는 개수와 SHA256만 남겼다. DB 값은 read-only transaction으로 메모리에 읽고 로그로 출력하지 않았다.

| 대상 | 전 → 후 / 판정 |
|---|---|
| schema | SHA256 `0fe69ca9d3c55799a0f97987e862ce8b0589ee3871b28b50a5aa1627d6bd7002` 동일. pg_dump의 임의 restrict 표식만 정규화 |
| writer 경계 | goose 2·guard 1·HumanModel 1·singleton 1행 동일 |
| 회원·agent·친구·관계 | Members 2·Agents 28·Friends 1·Pairs 13·Owners 24·Identities 2. 각 전체 map digest 동일 |
| 키·철회·이름·친구 인가 | Agents·Members·Friends·Pairs 전체 digest 동일. 키는 Agents 안에 있으므로 공개키/credential·철회 상태도 같은 바이트 구조로 보존 |
| 메시지·ACK 보존 | Messages 15·Idempotency 15·각 digest 동일. 새 전송·receive·ACK를 하지 않았다 |
| 세션 | Sessions 0·digest 동일. 새 Google 로그인은 하지 않았다 |
| Remote | Installs 2·Tokens 2·Clients 1·Subscriptions 1·Codes/Attempts 0. 각 map digest 동일 |
| Remote Inbox | 설치별 Inbox를 포함한 Installs digest 동일. 합계는 before/after의 remote_inbox_count 참조 |
| 설정·자격 | `.env`, tunnel.uuid, access.aud, tunnel/config.yml SHA256와 mode 동일. Google configured=true |
| 백업 | 기존 14개 파일의 hash·mode·크기 동일. 새 고유 백업 1개만 추가해 15개 |
| 공유 자원 | 기존 relay/migrate 외 모든 기존 컨테이너의 ID·이미지·상태·시작시각·restart 동일. KnowsLink Postgres/cloudflared도 동일 |
| 호스트 connector·listener | cloudflared PID 506937 동일. 8080/5432 listener는 기존 loopback만 유지 |

Events 구독은 재시작 전후 같은 identity digest·만료 `2026-10-13T14:09:16.701715Z`·Verified `13:09:16.701715Z`였다. attempts=1, pending=false, lease=false, Next=`2026-10-13T00:45:51.130657Z`도 같다. 이전 NOTIFY 기록의 만료보다 연장된 현재 구독은 관측했지만, 이번 배포 중 갱신 요청이나 새 delivery retry는 관측하지 않았다. 구독 지속 저장은 확인했고 다음 자동 갱신·호스트 턴·장시간 수신 성공은 not_run이다. URL·secret·install ID는 출력하지 않았다.

## 공개·보호·공유 smoke

[public-smoke.json](public-smoke.json)의 실제 HTTPS 검사 21개가 모두 통과했다. 응답 본문·쿠키·callback URL은 저장하지 않았다.

- `/` 200·no-store·CSP·no-referrer·nosniff. 무세션 `/home`, `/home/friends`, `/home/agents`는 303과 `/?n=expired`다.
- `/healthz`, `/owner`, `/v1/owners`, 정확한 `/v1/keys`, `/v1/authorize`, `/v1/test/pull`, `/_probe`는 404다.
- 무인증 `/v1/text/events`와 query token은 401 JSON이다. `/mcp` initialize POST도 401이다.
- OAuth well-known 두 경로는 200 JSON, 합성 43자 consent는 401 HTML, GET connect/start는 405다.
- 공유 서비스 orca/s8/mcp는 로컬 PC HTTPS와 서버 curl에서 200/200/401이다. 서버 Python urllib은 세 호스트에서 전후 모두 403이었다. 이 클라이언트 차이를 서비스 장애나 PASS로 바꾸지 않았다.
- `13:23:12Z` 검사에서 새 relay 전체 로그의 error/panic 0건이다. Google configured=true이며 up 완료 후 KnowsLink cloudflared error도 0건이다. 이는 약 2분 관측이며 장시간 보장이 아니다.

`verify.py regression`은 저장 baseline projects에 knowslink가 남는 기존 결함을 읽고 반복하지 않았다. 이번 직접 전후 컨테이너·connector·curl 비교로 보완했다. `verify.py public`의 과거 Access 302 전제도 현재 member allowlist+404와 달라 not_run이다. 보호 경로는 위 실제 smoke로 검사했다. Access dashboard 정책 GET은 하지 않았으므로 정책 객체의 신규 감사를 주장하지 않는다.

## Grok 업데이트 안내 — owner 후속, 이번 실행 없음

Grok은 대화 컴퓨터에서 실행하는 로컬 Command MCP다. [기존 실제 설치 보고](../SAR-AUTO-RECEIVE-001-LIVE/report.md)는 새 MCP process·9개 도구·기존 public-node/연결 폴더·watcher 유지를 기록했다. [최신 실계정 기록](../SAR-HUMAN-PEER-002-LIVE/report.md)과 [NOTIFY 기록](../SAR-HUMAN-PEER-003-NOTIFY/report.md)은 기존 agent와 webhook routine 수신을 확인한다. 설치 보고와 직접 관측의 구분을 유지한다. 이번에는 Bot 파일·현재 PID·등록 UI를 새로 조사하지 않았다.

서버 배포는 Grok의 이미 실행 중인 Node bundle을 갱신하지 않는다. 이번 TS 공통 처리도 받으려면 owner가 다음 유지보수를 수행한다.

1. 기존 Command/Arguments, `KNOWSLINK_AGENT_FOLDER`, receive mode, webhook 파일 경로와 정확한 watcher PID를 비공개 위치에서 확인한다. 키·agent ID·Inbox와 routine을 그대로 둔다.
2. 수락 SHA `2747fe98653d5c6946f41c219fc256395333d7fe`의 별도 깨끗한 소스에서 기존 `sh scripts/install_bot_mcp.sh`로 같은 설치 prefix의 bundle을 준비·검증한다. 기본 Command는 `/workspace/.knowslink/node/bin/node`, Arguments는 `/workspace/.knowslink/knowslink/dist/plugin.js`다. 실제 등록 값이 다르면 그 기존 prefix를 사용한다. installer는 자체 stage/설치물의 rm/mv를 포함하므로 이번 worker의 삭제금지 범위에서 실행하지 않았다. owner는 별도 유지보수 범위를 확인한 뒤 실행한다.
3. installer가 출력하는 신규 연결 안내를 기존 연결에 적용하지 않는다. 기존 `public-node`, `RELAY_URL`, **같은 기존 연결 폴더**를 유지한다. `knowslink_connect`를 다시 실행하거나 새 agent/key/친구 관계를 만들지 않는다.
4. 같은 Command 등록의 MCP process를 재시작한다. 기존 연결 폴더 watcher의 정확한 PID만 종료하고 새 bundle로 watcher 하나를 재시작한다. `pkill -f`나 다른 연결 폴더의 process 종료는 사용하지 않는다. 기존 `KNOWSLINK_GROK_WEBHOOK_FILE`, receive mode와 routine prompt를 그대로 적용한다. loopback wake로 전환하지 않는다.
5. `knowslink_status`의 connected/autoReceive metadata와 webhook 모드 `plugin.js wake-check`를 확인한다. wake-check는 파일 검증이며 발송하지 않는다. 표시되는 기존 agent ID가 같은지 확인한다. 실제 수신 알림 확인용 새 메시지는 별도 승인이 있는 후속에서만 보낸다.

프로토콜·credential·agent ID·친구 모델은 변하지 않았으므로 재인증·친구 재등록은 요구하지 않는다. 로컬 bundle을 갱신하지 않아도 이전 기능 보존 계약으로 서버를 사용할 수 있다. 새 로컬 코드 적용 여부와 기존 호환 사용을 구분한다. 기존 webhook routine과 Inbox/ACK marker는 삭제하지 않는다.

## Dots 업데이트 안내 — 기존 Remote MCP 유지

Dots는 [기존 실계정 연결](../SAR-HUMAN-PEER-002-LIVE/report.md)과 [Events 자동화](../SAR-HUMAN-PEER-003-NOTIFY/report.md)가 있는 remote app이다. 로컬 Grok ZIP 설치 대상이 아니다. 서버 endpoint `https://link.knowslog.com/mcp`, remote 도구·Events schema·OAuth 경로 변경은 0이다. 기존 app·Google 연결·설치 agent/key·친구·Events 자동화를 유지한다.

이번 서버 구현은 이후 remote 요청부터 적용된다. app 재등록·재인증·친구 재등록은 이번 변경에 필요하지 않다는 판단은 실제 diff와 보존 관측에서 도출했다. Dots 내부 cache·다음 자동 갱신·화면 수신을 새로 검증했다는 뜻은 아니다. owner가 auth 만료 등 별도 오류를 받으면 해당 오류를 확인한다.

공식 OpenAI 문서는 tool names/descriptions/schemas/annotations/auth/UI resources가 바뀔 때 서버 배포→Plugins의 기존 연결에서 Refresh→새 대화 검사를 요구한다. Events 목록을 바꿀 때도 rescan한다. 이번에는 해당 metadata 변경이 없으므로 refresh는 필수 단계로 추가하지 않는다. 향후 이 metadata가 바뀌면 기존 연결을 refresh하며, app을 삭제하거나 새 Google 설치를 만들지 않는다. 새 대화 확인이 필요해도 기존 자동화·구독은 유지한다.

## 공식 근거와 적용 범위

Context7 resolve-library-id는 `Monthly quota exceeded`를 반환했다. 공식 원문을 2026-10-13에 열어 대체 확인했다. OpenAI Docs 스킬을 적용했다.

- [xAI Team Bots](https://docs.x.ai/grok-bot/team-bots#set-up-what-the-bot-needs): Command는 각 대화의 컴퓨터에서 실행한다. 개인 Bot의 실제 설치는 위 LIVE 관측이 근거다.
- [OpenAI Refresh metadata](https://developers.openai.com/plugins/deploy/connect-chatgpt#refresh-metadata): metadata 변경에 따른 Refresh와 새 대화 검사다. 변경 없는 이번 서버 코드 반영 판단은 diff 기반 추론이다.
- [OpenAI MCP Events](https://developers.openai.com/plugins/build/mcp-events#manage-subscriptions): 구독 상태의 재시작 보존·접근 재검사와 갱신 계약이다. 이번 metadata 관측을 호스트의 향후 성공 보장으로 확장하지 않는다.

## 실패·미실행·복구·후속

기록 커밋 `c914962066f0cd5aeeb551d4316b09b3dd5b9235`에서 설치 FullOps1.4.1 lint를 기준2747fe9로 실행했다. 결과는 ERROR0·WARNING1·unavailable2·EXIT1이다. Windows에 make가 없어 product-lint/product-test는 실행하지 못했다. 경고는 기존 PLANS1749줄의 SIZE-001이다. 원 결과는 [lint-windows.json](lint-windows.json)에 보존했다. 새 제품 코드·테스트·설정·의존성 변경이 없으므로 새 제품 등록검사의 적용 변화는 없다. 기존 exact cd4dae8의 실제 Linux product-lint/product-test EXIT0은 [gate 재사용 근거](gate-reuse.json)로 연결하며 새 Windows PASS로 표시하지 않는다. 산출물 strict 검사는13개·문제0·경고0이다. 로그 사본의 끝 공백을 정규화한 뒤 git diff --check를 통과했다.

Git 로그 사본은 UTF-8/LF와 행 끝 공백만 정규화했다. [evidence-copy-audit.json](evidence-copy-audit.json)은 원본·사본 hash와 범위를 남긴다. AppData 원본은 변경하지 않았다. 원본 바이트 동일성을 주장하지 않는다.

착수 packet336경로/338쌍의 수동 보완은 [원 audit](packet-audit-dispatched.json)과 [원 outcomes](packet-outcomes-dispatched.json)에 보존했다. 실제 기술 계획과 복귀 주소를 채운 뒤 공식 `jev_packet.py --force --write-inbox`로 identity를 갱신했다. 갱신 packet은138경로/140쌍이며 stale producer는 현재 분류로 사용하지 않았다. unknown19건·remaining105경로를 포함한150경로의 존재/hash/관련 검색 구간을 추가 확인했다. 현재 outcomes140쌍의 unknown은0이다. 과거 전문 재독은 주장하지 않는다. `work.py finish`가 완료 계약을 통과해 [날짜별 인박스 전문](../../../../handovers/logs/2026-10-13_to_ops.md)에 지시서와 결과를 보존했고 현재 ops 인박스는 비었다.

첫 after-checkout 관측은 old HEAD assertion이 남아 EXIT=1이었다([before-build-observe.log](before-build-observe.log)). 다음 사전 관측은 업로드 파일의 CRLF 때문에 heredoc 종료를 읽지 못해 EXIT=1이었다([deploy.log](deploy.log)). 둘 다 relay stop 전에 실패했다. LF 고유 사본과 정확한 NEW HEAD 조건으로 수정한 뒤 같은 배포를 완료했다. 실패 로그·원 파일·부분 before.json은 보존했다. 격리 initdb의 locale 경고는 보존했으며 복원·guard·종료는 EXIT=0이었다.

실계정 송신/receive/ACK·새 Google 로그인·인증 SSE 재연결·Grok bundle/process 직접 갱신·Dots refresh/자동화 수정·실제 화면/모바일/키보드/200%·부하·장시간 wake·향후 구독 갱신은 not_run이다. UI template/CSS 변경이 없으므로 새 직접 시각 검수는 추가하지 않았다. 기존 QA/시각/held와 전체 HP-16 수락 경계는 유지한다.

문제 발생 시 relay를 멈추고 현재 DB를 새 고유 백업한다. 이번과 직전 `26094e9`는 같은 schema·HumanModel writer이므로 검토한 직전 코드로만 복귀할 수 있다. HumanModel 이전 writer와 자동 Down은 금지한다. 현재 친구·철회·generation·키·Remote/Inbox를 잃는 과거 DB 덮기는 하지 않는다. 이번에 실제 rollback을 실행하지 않았다.

후속 담당 coor는 이 문서 결과의 고정 SHA 독립 snapshot 리뷰·main 일반 통합을 수행한다. owner는 승인된 후속에서 Grok 로컬 bundle/process 갱신과 필요 시 실제 수신 확인을 맡는다. 이번 OPS는 외부 Bot에 메시지를 보내지 않았다.

원 작업 파일은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor001-ops/`에 보존했다. 서버 bundle·관측·backup.list·부분 실패 자료는 `/home/shin/deploy/knowslink-state/refactor001-ops-20261013/`에 있다. 복원된 비공개 DB는 위 stopped 격리 컨테이너 안에 보존했다. 레포 증거는 같은 폴더의 JSON/log와 [commands.txt](commands.txt)다. 탐색 보완과 처리 결과는 [packet audit](packet-audit.json) 및 [packet outcomes](../../jev/SAR-REFACTOR-001-OPS-packet-outcomes.json)로 연결한다.
