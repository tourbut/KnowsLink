---
title: Astra 리팩터링 독립 동작 QA
status: review
updated: 2026-10-14
owner: tester
tasks: [SAR-REFACTOR-002-TESTER]
summary: 고정 bcb 제품의 Inbox 오류 복구·TTL·보관·ACK·callback 경계를 독립 검사하고 짧은 CPU 벤치를 재측정했다.
---

# SAR-REFACTOR-002-TESTER

지정 동작 QA는 PASS다. 새 제품 결함과 미해결 critical/high 발견은 0건이다. 제품 코드는 수정하지 않았다. 별도 Astra 코드 리뷰, 최종 DEV 문서·등록 gate, main 통합은 coordinator가 취합한다.

## 입력과 독립성

- 고정 제품: `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`.
- 비교 제품: `8f9370bb62c9133c73e10ea2307b4cd696e10f7e`.
- 착수 HEAD: `8c08a1b0e69307c05f6d55a9674ac91925f036af`. 제품 경로 diff가 비어 있음을 [입력 증거](input-identity.json)에 기록했다.
- 실제 tester 세션: `01a139b0-8386-7221-bb61-44399fb01ea1`. Task `task_677998bdba45`, Dispatch `ctx_2b4b7f958935`, terminal `term_52277020-d274-4fbd-a842-384f1b110aa6`.
- 사용자 고정 codex/gpt-6-astra/high, Purpose review, Test level lite. 작성 DEV와 별도 세션이며 코드 snapshot 리뷰를 중복하지 않았다. 좁은 QA를 직접 수행해 선택형 하위 위임은 사용하지 않았다.
- 적용: AGENTS.md, FULLOPS.md, fullops-common-0.3.3 세 규칙, project.md, document-writing.md, fullops-test와 fullops-work, 현행 system-design, DEV 감사·adapter·remote·retention 회귀 원문.

독립 Linux clone은 `/tmp/knowslink-qa-refactor002-269a21604c94/repo`다. Go1.27.1 linux/amd64, CGO=1, Node22.22.2, npm10.9.7, AMD Ryzen 3 3200G다. SSH는 기존 known_hosts와 RejectPolicy를 사용했다. 접속 인증값은 메모리에만 읽고 출력에서 가렸다.

제품 clone은 실행 전후 exact SHA와 clean을 확인했다. 새 Go 파일은 clone 밖에 두고 `-overlay`로만 연결했다. 독립 mjs도 clone 밖에서 빌드된 제품 모듈을 불러왔다. 새 assertion 두 파일의 SHA256은 입력 증거에 있다. Node 의존성은 lockfile로 설치하고 TS를 다시 빌드해 stale dist를 배제했다.

## 실행 결과

전체 재현 명령과 각 SSH 원본 종료코드는 [commands.jsonl](commands.jsonl)에 있다. 실행기는 [run.py](run.py)다. 셸 `set -e`와 SSH channel 종료코드를 사용했고 로그 축약 파이프는 사용하지 않았다.

| 검사 | 새 실행·종료코드 | 증거 |
|---|---|---|
| 환경·후보 clean·TS build | 0 | [prepare.log](prepare.log) |
| 관련 Go 회귀 및 새 경계 | `go test -race -overlay ... ./internal/relay -count=1 -v`, 0. 최상위 16개 통과 | [qa-race.log](qa-race.log) |
| 기존 Inbox·owner wire 회귀 | Node22, 0 | [qa-inbox.log](qa-inbox.log) |
| 새 Inbox 오류 matrix | Node22, 0 | [inbox-boundaries.mjs](inbox-boundaries.mjs), 같은 로그 |
| 격리 DB migration·HTTP/claim 회귀 | `go test -race -tags=integration ... -count=1 -v`, 0. 4개 통과, 5.512초 | [database.log](database.log) |
| DB 정상 stop | 0. 삭제하지 않음 | [database-stop.log](database-stop.log) |
| before/after CPU 벤치 | 각각 0 | [before](benchmark-before.log), [after](benchmark-after.log) |
| 최종 제품 SHA·clean | 0 | [post-clean.log](post-clean.log) |

DB는 새 Postgres17 컨테이너 `kl-qa-refactor002-6405308e3d1f`다. digest는 run.py의 고정 IMAGE다. CPU 1개·512MiB 제한, 동적 loopback 포트, 합성 자격만 사용했다. 실제 실행은 TestHumanPostgres, TestSharedExecutionClaimHTTP, TestHTTPConcurrencyAcrossInstances, TestRemoteOAuthHTTP다. 사용자 관계·이력·writer 보호, 실행 claim, 다중 Service 경합, Remote OAuth/MCP 계약을 확인했다. 운영 DB와 실계정은 사용하지 않았다.

## 새 독립 assertion

[independent_test.go](independent_test.go)는 기존 fixture를 사용하되 다음 경계와 기대값을 새로 작성했다.

1. queued/delivered × active/revoked 4조합에서 TTL 1ns 전과 정각을 확인했다. 정상 queued는 조기 만료되지 않는다. 철회된 transport와 arrival은 차단된다. 정각에는 transport·lease·claim token이 비워지고 delivered receipt와 이미 받은 `received` completion은 보존된다. 미완료 claim의 completion은 `failed:expired`가 된다.
2. 24시간 1ns 전에는 참여자 두 명만 대화 사본을 본다. 제3자는 보지 못한다. 철회 뒤에도 참여자 이력은 남으며 활성 전송 인가를 복구하지 않는다. 정각에는 message·idempotency가 제거된다.
3. 실제 persist/ACK 경로에서 persist 이전 ACK를 거부하고 persist 이후 ACK를 수락한다. sweep claim 필드 검사는 합성 직접 상태이며 실제 execution claim 경계는 별도 DB HTTP 회귀가 담당한다.
4. callback을 정확한 15초 lease 경계마다 3회 claim한다. 마지막 lease 1ns 전 새 도착은 lease를 우회하지 않는다. 정각에는 기존 event의 네 번째 시도 없이 즉시 재예약한다. 새 event ID와 첫 시도 예산을 확인했다. 과거 event의 늦은 성공 응답은 새 event lease를 지우지 않는다. 새 event 성공 뒤에는 subscription 만료까지 idle이다.
5. 같은 마지막 lease 시점의 관계 철회, 메시지 만료, 구독 만료, 설치 만료, 새 도착 없음은 새 callback을 만들지 않는다. 호출 간 pending 캐시 재사용 방지는 기존 TestRemoteClaimsRecheckCurrentAuthority를 이번 프로세스에서 독립 재실행했다.

새 Inbox matrix는 unread 본문 2개와 fresh tmp marker를 실제 파일로 만든다. stat에 EACCES/EIO/EPERM/EACCES를 차례로 주입해 매번 오류 전파·두 원본 JSON·본문 전체 동일성을 확인한다. ENOENT만 통과하고 오래된 순서로 첫 본문을 반환한다. 오류 제거 뒤 두 번째 본문도 반환한다. read marker는 비어 있고 같은 ID 재저장은 false다. 빈 Inbox에서는 오류 stat을 호출하지 않고 null을 반환한다. mock은 marker stat에만 적용하며 finally에서 복원했다. 오래된 marker의 파일 제거 분기는 실행하지 않았다.

## 첫 실패와 판정 수정

첫 새 Go 검사에서 delivered 2조합이 실패했다. [원 로그](qa-race-attempt1.log)와 [원 테스트](independent-attempt1.go.txt)를 보존했다. 원 종료코드는 1이다.

원인은 테스트 기대값이다. 실제 public text ACK는 `http.go`의 persist/ack 처리에서 completion을 `received`로 만든다. `store.go`는 빈 completion만 만료 상태로 바꾼다. 최초 테스트는 ACK 완료도 `failed:expired`가 되어야 한다고 잘못 기대했다. 제품을 바꾸지 않고 delivered의 기대값을 `received`로 바로잡았다. 이는 기존 완료 보존 계약을 더 엄격히 검사한다. 수정 후 같은 고정 제품과 동일 race 명령이 0이었다. 제품 RED/GREEN 수정으로 주장하지 않는다.

## 짧은 성능 재측정

동일 Linux 호스트·Go·fixture·명령으로 before를 먼저, after를 다음 실행했다. baseline clone은 비교 제품 SHA이며 후보의 두 benchmark 파일만 복사했다. 실행은 `-run '^$' -bench '^(BenchmarkRetainedHistory|BenchmarkRemoteClaimSubscriptions)$' -benchmem -benchtime=3x -count=3`이다. [중앙값](benchmark-medians.json)은 각 세 표본에서 계산했다.

| 10000개 fixture | 이전 ns/op | 이후 ns/op | 관측 변화 |
|---|---:|---:|---:|
| 만료 보관 이력 sweep | 56907942 | 1130675 | 98.0% 감소 |
| JSON 포함 transaction CPU | 197382077 | 148485943 | 24.8% 감소 |
| Remote claim, 구독 4개 | 1633501 | 512394 | 68.6% 감소 |
| Remote claim, 구독 1개 | 465271 | 504850 | 8.5% 증가 |

fixture는 100명·200 agent, 만료 delivered 이력 100/1000/10000개다. 최대 이력 JSON은 6925918bytes다. claim은 동일 agent의 1/4 구독과 delivered 기록 100/10000개다. 단일 구독은 재사용 이득이 없으며 모두 빨라졌다고 주장하지 않는다. 짧은 공유 호스트 측정이라 신뢰구간·장시간 안정성·운영 처리량을 주장하지 않는다.

sweep의 큰 이득은 전체 DB transaction 이득과 다르다. 새 측정은 CPU만 포함한다. 전체 JSONB 직렬화·SQL 변환·전역 row lock은 남는다. DEV의 별도 실제 DB 결과는 10000개에서 588663216→542726660ns/op, 약 7.8% 감소다. 해당 수치는 재사용이며 tester의 새 DB 벤치가 아니다. HTTP 입장/업무/퇴장 합계나 end-to-end throughput으로 환산하지 않는다.

## DEV 근거 재사용과 필수 gate

[dev-evidence-index.json](dev-evidence-index.json)은 DEV 원본 경로와 복사본 SHA256을 기록한다. [DEV integration](dev-integration.log)은 exact bcb에서 전체 tagged race 111.338초·native0이다. 최초 build 누락 MODULE_NOT_FOUND/native1은 DEV 원본에 보존되어 있다. 새 4개 DB 회귀와 이 전체 회귀의 실행자를 구분한다.

[DEV FullOps](dev-fullops.json)의 head는 exact bcb다. product-lint/product-test는 각각 exit0이다. 전체 lint는 ANTI-004의 `inbox.pending()` 오인으로 ERROR1/WARNING2/exit1이므로 FullOps 전체 PASS로 재사용하지 않는다. DEV가 최종 테스트·문서 기록을 수정하고 별도 등록 gate를 진행한다. 이 tester는 고정 bcb 원본 테스트를 유지한다.

새 QA 코드 때문에 세션 done-gate는 최종 기록 커밋의 `lint.py --from bcb1dd6629b508a30ec39a5d3141c8ce0a558597`를 요구한다. 최종 커밋 후 외부 task-local make.cmd가 같은 기록 HEAD의 깨끗한 Linux clone에서 실제 make lint/test를 실행한다. 전체 DB integration과 build를 다시 돌리지 않는다. gate의 정확한 결과 SHA·native exit·ERROR/WARNING은 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/qa-refactor002/final-lint.json`과 worker_done에 남긴다. 이 문서 작성 시점에는 그 후속 실행의 PASS를 선기록하지 않는다.

## 패킷·보존·한계

현재 착수 HEAD에서 find/documents-find/context/packet을 정식 force 갱신했고 이전 결과는 history에 보존했다. route bind는 이전 producer SHA 불일치로 거부되었다. 해당 unknown을 숨기지 않고 현재 지시서의 task·role·attempt·고정 제품으로 수동 확인했다. [packet-audit.json](packet-audit.json)은 220개 경로의 존재·해시·제목·관련 구간을 기록한다. outcomes는 비선택 217 path/category를 처리했다. 큰 과거 문서의 전문 재독과 완전한 호출 그래프를 주장하지 않는다. 실제 관련 store/remote/Inbox 원문·새 경계 실행이 미전송·partial을 보완한다.

레포 밖 보존 자료:

- 로컬 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/qa-refactor002/`: bundle, remote.txt, container.txt, 최종 gate 경로·로그·결과.
- 원격 `/tmp/knowslink-qa-refactor002-269a21604c94/`: product bundle, clean repo, baseline, 두 독립 검사 소스, overlay.
- 합성 Inbox `/tmp/kl-inbox-files-1SgPQE`, `/tmp/kl-qa-refactor002-Mj3l8k`.
- stopped DB `kl-qa-refactor002-6405308e3d1f`. 최종 gate clone은 외부 gate-remote.txt로 추적한다.

삭제·이동 정리·guard 우회·배포·실계정 발송은 하지 않았다. UI 변경이 없어 새 시각 검수는 해당 없음이다. live Google/외부 callback/실제 호스트 wake·장시간·최대 부하·백업 복원은 미실행이다. local claim 뒤 crash/read/write 실패 window와 전체 JSONB/row lock 한계는 남는다. 이번 결과는 그 한계나 과거 held 제품 수락을 해제하지 않는다.
