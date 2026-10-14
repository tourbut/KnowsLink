---
title: Astra 리팩터링 운영 배포와 복구 후보 검증
status: draft
updated: 2026-10-14
owner: ops
tasks: [SAR-REFACTOR-002-OPS]
summary: 수락 제품8a93e38의 실제 배포·백업·격리복원·보존·공개 검사와 이전 이미지 조회 실패 및 복구 후보를 기록한다.
---

# SAR-REFACTOR-002-OPS 결과

수락 제품 `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6`을 기존 운영에 배포했다. 새 private 백업·격리 복원·writer guard·권한/Inbox/Events·공유 자원 보존과 안전한 공개 검사를 통과했다. 이전 이미지 주소는 전환 뒤 조회할 수 없어 같은 이전 소스로 별도 복구 후보를 만들었다. 원 이미지 동일성과 실제 rollback 성공은 주장하지 않는다.

## 범위·기준·시간

- 사용자 실제 배포 승인과 현재 역할 인박스를 적용했다. 제품 변경·계정 송수신·재등록·키 교체·다른 서비스 변경·삭제/이동 정리는 수행하지 않았다.
- 준비 HEAD는 `63ba2167352a4b6b1ec2df8782841ba768963176`, 기준 제품은 `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6`, 브랜치는 `fullops/ops`다. 결과 문서 커밋은 worker_done의 `SHA:`이며 배포 제품과 다르다.
- Task `task_6e35fa7983b5`, Dispatch `ctx_82a534ee9cf7`, worker `term_d9fae13d-e296-47fa-b737-f71ae73a48d0`, 복귀 Run `run_5534969da2d0`이다.
- AGENTS/FULLOPS, fullops-common-0.3.3, project.md, D07/D12/D13, document-writing과 현재 OPS context를 읽었다. fullops-orca/work/deliverables를 적용했다. 단일 운영 writer 전환은 직접 직렬 수행했다. coor의 독립 결과 리뷰는 후속이다.
- 2026-10-14 11:27:54Z의 fresh baseline은 `2747fe98653d5c6946f41c219fc256395333d7fe`, detached·clean·relay/Postgres healthy였다. 과거10월13일 관측을 현재 상태로 사용하지 않았다.

## 실제 배포·백업·복원

| 항목 | 실제 결과 |
|---|---|
| bundle | main이 고정8a93e38임을 확인했다. SHA256 `b1d19faf32f50332660d24d03618cd21fb328448da37d62d913b624b7dcdfc4e`; 서버 hash·bundle verify·FETCH_HEAD·이전 SHA 조상 확인 |
| 실행 입력 | 이전2747→8a의 DB/migration·deploy·Compose·Dockerfile·go.mod/go.sum·Makefile·sqlc diff0. 제품 수정은 기존 Inbox·RemoteEvents·sweep 세 경로 |
| build | 고정 clean checkout에서 relay/migrate build EXIT0. 기존 relay는 build 중 healthy |
| 중단 | STOP_BEGIN `2026-10-14T11:30:57Z`, UP_DONE `2026-10-14T11:31:13Z`; KST20:30:57~20:31:13, 약16초 |
| writer 고정 | relay stop 뒤 epoch1001476이2초 간격으로 동일. 다른 DB client0. frozen 관측의 client1은 검사 psql 자신 |
| 새 백업 | `/home/shin/deploy/knowslink-state/backups/2747fe9-pre-refactor002-20261014T113101Z.dump`; 상위0700·파일0600·20039bytes |
| 백업 SHA256 | `3e92b655bd884378cd04d9624804d60d5d607c5c07fffe865525301cea9dc6cf` |
| list·복원 | pg_restore list EXIT0·TOC25줄. network none·host port 없음·postgres 사용자로 복원 EXIT0 |
| 전체 상태 | 복원 data::text digest와 frozen 운영 data digest 일치. goose2·guard1·singleton1행·epoch1001476·HumanModel1 |
| 거부 경계 | 격리 DB에서 HumanModel 제거 UPDATE가 `human authorization writer downgrade refused`로 거부됐다. live DB에는 거부용 쓰기를 하지 않았다 |
| 재개 | 미리 만든 이미지로 Compose `up -d --no-build --wait relay`; migrate EXIT0·기존 goose2 유지. local 검사 EXIT0 |
| 실제 relay 이미지 | `sha256:482d3197fe7ca26eef7918a2061feede01bd5199dfd8be37af5f46c4e381343d`; healthy·restart0 |
| 로그 관측 | 11:34:22Z, 기동 약3분 뒤 relay error/panic0·Google configured=true·기동 후 Tunnel error0 |

기존 `humanpeer001-ops/remote.py`의 known_hosts·RejectPolicy·메모리 자격을 재사용했다. 비밀값·원 DB 행·본문·callback URL/secret을 로그에 출력하지 않았다. `beta.sh deploy/backup/restore-verify`는 실행하지 않았다. `up --no-build`는 이미 만든 이미지만 사용해 중단 중 중복 build를 피한 선택이다. Compose v5.1.3의 실제 help에서 지원을 확인했다. 지시서의2026-10-14 Context7 quota 실패와 [공식 up](https://docs.docker.com/reference/cli/docker/compose/up/)·[stop](https://docs.docker.com/reference/cli/docker/compose/stop/) 근거를 재사용했다.

복원 fixture `knowslink-refactor002-rehearse-20261014`는 Postgres 정상 stop 뒤 container exited0이다. private 복원 데이터는 컨테이너에 보존됐다. 백업/list는 서버 private 경로에만 남겼다. 운영 DB 복원·Down·volume 삭제·prune·다른 서비스 변경은 없다.

## 보존·현재 시점

[before](before.json), [frozen](frozen.json), [after](after.json), [preservation](preservation.json)을 대조했다. read-only transaction으로 값은 메모리에만 읽고 count/digest만 기록했다.

- 회원2·agent28·owner24·친구1·pair13·신원2의 전체 map digest가 같다. 이름·키·credential·철회·friend/generation·suppression을 포함하는 상태를 보존했다.
- 메시지11·멱등11·세션0·연결0·GoogleAttempts0·Challenges0·Gates0의 digest가 같다. 별도 Conversations/Suppressions map이 없는 현재 상태에서는 해당 값이 메시지/회원/친구 전체 map에 포함된다. 과거15개 메시지 관측을 복원하지 않았다.
- Remote Clients1·Installs2·Tokens2·Subscriptions1·Codes/Attempts0과 Inbox합계0이 같고 각 전체 map digest도 같다. token 재발급·소비·설치 추가를 하지 않았다.
- Events Exp `2026-10-14T11:55:42.375744Z`, Verified `2026-10-14T10:55:42.375744Z`, attempts1·pending=false·lease=false·Next와 identity digest가 같다.10월13일보다 연장된 현재 구독을 관측했지만 이번 작업이 갱신한 것은 아니다.
- schema SHA256·goose/guard/HumanModel·Google env·Tunnel config/credential·파일 mode가 같다. 기존 백업15개 hash/크기/mode를 보존했고 새 백업1개만 추가했다.
- relay/migrate 외 기존 컨테이너31개의 ID·image·running·health·시작시각·restart가 같다. KnowsLink Postgres/cloudflared도 포함한다. host cloudflared PID506937과 loopback8080/5432 listener도 같다.
- Rates만 digest가 다르고 개수4는 같다. 기존 `identity.go`의 요청 시각 기록과 window 정리 및 local 보호 검사가 갱신하는 volatile map이다. 권한 변경으로 해석하지 않는다. HTTP map은 검사 전후 빈 상태다. 모든 비교를 byte 불변이라고 과장하지 않는다.

보존 assertion33개가 통과했다. 만료 가능한 데이터도 이번 짧은 전환 중 일치했다. 향후 정상 만료·cleanup·사용자 요청에 따른 변경까지 금지하는 보장은 아니다.

## 공개·공유 검사

[public-smoke.json](public-smoke.json)의21개 HTTPS 검사가 통과했다. 실제 송신·pull/receive·ACK·token 소비·로그인을 하지 않았다. 응답 본문·cookie·redirect URL은 보존하지 않았다.

- `/`200과 no-store/CSP/no-referrer/nosniff; 무세션 `/home`, `/home/friends`, `/home/agents`303과 expired redirect.
- `/healthz`, `/owner`, `/v1/owners`, 정확 `/v1/keys`, `/v1/authorize`, `/v1/test/pull`, `/_probe`404.
- 무인증 SSE와 query token401, 무인증 MCP initialize POST401.
- OAuth well-known 두 경로200, 합성43자 consent401, GET connect/start405.
- 공유 orca/s8/mcp는 로컬 HTTPS와 서버 curl에서200/200/401. 서버 urllib의403/403/403은 baseline부터 동일한 client 차이로 보존했다.

과거 `verify.py regression`의 projects baseline 문제와 `verify.py public`의 Access302 전제는 현재 계약과 달라 실행하지 않았다. 실제 전후 자원 비교·현재404/401 보호 검사로 보완했다. Access 정책 객체의 신규 감사나 새 UI 시각 수락을 주장하지 않는다. UI/template/CSS diff가 없어 시각 검사를 추가하지 않았다.

## 이전 이미지 조회 실패와 복구 후보

fresh baseline과 build 직전의 실제 container Image는 `sha256:220edb718b2fe1925eed243084e41fedf641c86a0342760f682672b91379ed3d`였고 당시 image inspect는 성공했다. 전환 후 같은 주소의 image inspect는 종료1이었다. [원 실패](post-check.log)와 [새 image 목록](image-check.log)을 보존했고 coor에 msg_190ddfe4ee12로 즉시 알렸다. 과거 보고만으로 보존됐다고 주장하지 않는다. 이 작업에서 image 삭제/prune 명령은 없었다. 태그 없는 이전 주소가 전환 이후 조회되지 않는 원인의 별도 Docker 내부 감사를 수행하지 않았다.

검토된 직전 소스 `2747fe98653d5c6946f41c219fc256395333d7fe`를 고유 `rollback-source/`에 git archive로 추출했다. pinned Dockerfile과 의존성으로 별도 `knowslink-rollback:2747fe9-refactor002-20261014`를 빌드했다. [빌드 EXIT0](rollback-image.log)의 image ID는 `sha256:12f1313a1460997710d727eb543eebdf9b4dee5b8f8ffdc01dc25d50442bc72a`다. build cache를 사용했어도 원 manifest 동일성으로 해석하지 않는다.

현재 live checkout·service tag·env·DB는 유지했고 healthy/restart0을 다시 확인했다. coor msg_effcf5b52563도 같은 소스의 별도 rollback 준비를 승인된 범위로 확인했다. [후속 검사](post-check-recovered.log)는 원 이미지 보존 false와 재빌드 후보를 분리한다. 원 실패를 PASS로 바꾸지 않았다. 다음 배포는 build 전에 고유 rollback 태그와 실제 가용성을 확보한다. 실제 rollback은 미실행이며 과거 backup을 현재 DB에 덮지 않는다.

[복구 후보 artifact 검사](rollback-verify.log)도 EXIT0이다. network none·read-only 컨테이너에서 relay/migrate 실행 파일 존재·SHA256을 확인했다. `knowslink-refactor002-rollback-inspect-20261014`는 exited0으로 보존했다. DB에 연결하거나 앱을 기동한 검사는 아니다.

## gate 재사용·문서 검사·탐색

수락 DEV `3e5af5e5ef9ca993baa6745aaf5b2b1f157ace46`와 배포8a의 cmd/internal/adapters/db/deploy/scripts 및 build/config/dependency 입력 diff0을 직접 확인했다. [gate-reuse.json](gate-reuse.json)은 실제 final DEV Linux make lint/test EXIT0과 입력 동일성을 연결한다. 별도 build·Postgres race·독립 QA는 [최종 제품 리뷰](../SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md)와 [QA 리뷰](../SAR-REFACTOR-002-QA-REVIEW-review/report.md)의 원 SHA/조건으로 재사용한다. 새 문서 HEAD의 실행 PASS로 바꾸지 않는다.

이번 제품 코드·테스트·설정·의존성 변경은 없다. task-local 운영 명령은 LF/syntax·직접 실행·보존 assertion·공개 검사로 검증했다. 변경 문서는 D12/D13·README·operations·OPS phase/context/QA다. 산출물 strict는13개·문제0·경고0이고 JSON/packet72쌍·frozen-to-after assertion·git diff --check를 통과했다.

Windows 현재 PATH에 make가 없다. 이 문서 결과의 최종 커밋에서 설치 FullOps1.4.2 lint를 기준8a로 실행하며 원 결과는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002-ops/final-lint.json`에 보존한다. 실제 반환 결과는 worker_done에서 전달한다. 등록 make lint/test가 unavailable이면 새 PASS로 보고하지 않는다. 제품 입력 변경이 없어 지시서에 따라 위 exact DEV Linux gate를 재사용하고 반복 전체 제품 검사를 추가하지 않는다. 검사기/설정/환경을 완화하거나 통과를 위조하지 않는다.

착수 packet의 과거 head8a와 실제 준비HEAD63ba 불일치를 확인했다. 공식 `jev_packet.py --force --write-inbox`로 현재HEAD packet을 만들었고 필수19개 경로를 유지했다. 현재69경로·72개 path/category별 outcomes와97경로의 존재/hash/관련 제목·검색 구간을 기록했다. partial/unknown18개·과거 producer4개는 현재 의미 검색 성공으로 승격하지 않았다. 큰 D12/PLANS·과거 archive는 필요한 구간을 수동 확인했다. PLANS는 coor 소유라 수정하지 않았다. 과거 CODE/FINAL 실패와 수락 FINAL-MAIN/QA, 운영 미배포 설명과 이번 실제 배포 시점을 구분했다. 과거 문서 전문 재감사는 주장하지 않는다.

## 후속·보존 위치

Grok Inbox 수정은 같은 prefix·Command·기존 연결 폴더·agent/key·Inbox·webhook routine을 보존한 별도 로컬 bundle 갱신과 MCP/정확한 watcher 재시작이 필요하다. 이번에 실행하지 않았다. Dots의 Remote 처리 수정은 기존 endpoint/app의 다음 요청부터 적용된다. wire/tool/OAuth/Events metadata diff0이므로 새 연결·친구 재등록을 요구하지 않는다. host cache·화면·다음 자동 갱신은 미검증이다.

실계정 송수신·Google 로그인·인증 SSE·외부 callback/호스트 turn·장시간·최대부하·실제 rollback은 미실행이다. JSONB singleton/row lock 한계와 로컬 claim 뒤 crash/read/write window는 남는다. coor는 고정 결과 문서의 독립 리뷰와 main/origin 통합을 수행한다.

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002-ops/`.
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/refactor002-ops-20261014/`; bundle·관측·backup.list·rollback-source와 archive 보존.
- 격리 복원: 종료한 `knowslink-refactor002-rehearse-20261014`; 비공개 복원 데이터 보존.
- Git 증거: 이 폴더의 redacted JSON/log·[commands.txt](commands.txt). [사본 감사](evidence-copy-audit.json)는 UTF-8/LF·끝 공백 정규화와 원 hash를 구분한다. 원본 바이트 동일성을 주장하지 않는다.
