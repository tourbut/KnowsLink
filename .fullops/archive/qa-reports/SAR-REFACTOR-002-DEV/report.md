---
title: SAR-REFACTOR-002 구현·회귀·성능 결과
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-REFACTOR-002-DEV]
summary: Inbox 본문 유실·Remote 도착 지연 수정과 보관 이력 비용 절감 및 검증 한계를 기록한다.
---

# 구현과 검증

기준은 `8f9370bb62c9133c73e10ea2307b4cd696e10f7e`이며 제품 후보는 `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`다. 이후 기술 문서·검증 helper·증거만 갱신한다. 최종 인계 SHA와 정확한 HEAD의 등록 gate 결과는 worker_done 및 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002/final-lint.json`에 고정한다.

기존 계약 안에서 세 제품 파일을 고쳤다. Inbox는 housekeeping 실패 시 미표시 본문을 보존한다. Remote는 마지막 lease 소진 뒤 새 pending을 즉시 재확인하며 같은 호출의 agent 검색을 재사용한다. sweep은 만료된 이력의 불필요한 현재 권한 순회만 생략한다. 원인·재현·호출자·잔여 한계는 [감사](audit.md), [adapter](adapter-report.md), [remote](remote-report.md)에 있다.

## 동일 조건 측정

Windows Go1.27.1/AMD Ryzen 5 5600, `-benchmem -benchtime=200ms -count=3`의 중앙값이다. 100명/200 agent와 만료 delivered 이력 개수만 바꿨다. fixture·테스트 소스·컴파일러는 동일하며 제품 수정 전후 각각 실행했다. 원 로그는 [before](benchmarks/before-windows.txt)와 [after](benchmarks/after-windows.txt)다.

| 보관 메시지 | sweep 이전/이후 ns/op | transaction CPU 이전/이후 ns/op | transaction CPU 이전/이후 B/op | 이전/이후 allocs/op |
|---:|---:|---:|---:|---:|
| 100 | 420339 / 3776 | 2222973 / 1559679 | 510904 / 505597 | 3439 / 3339 |
| 1000 | 4248390 / 27181 | 13161200 / 7546163 | 4276403 / 4131659 | 12585 / 11580 |
| 10000 | 38099700 / 290768 | 120844250 / 69121733 | 58206720 / 44143589 | 103478 / 93461 |

10000개에서 sweep 약 99.2%, transaction CPU 약 42.8% 감소다. sweep 할당은 모든 크기에서 메시지당 16B/1회에서 0으로 줄었다. transaction CPU는 JSON decode·event hash·sweep·encode를 포함하며 SQL·네트워크·row lock 대기는 제외한다. 메모리 값은 GC·buffer 재사용에 따라 변동한다. 짧은 비격리 반복이므로 통계적 신뢰구간이나 운영 처리량을 주장하지 않는다.

같은 설치의 4개 구독/10000개 기록은 Remote claim 중앙값 590438→162572ns/op, 64B/1alloc 유지다. 1개 구독에서는 148226→156203ns/op로 증가했다. 모든 경로가 빨라졌다고 주장하지 않는다. 상세 조건은 remote 보고를 따른다.

## 검증 증거

| 검사 | 대상·결과 |
|---|---|
| 만료 경계·Windows Go unit | 미커밋 제품 후보, native exit0. 만료 정각/직전·철회·delivered 유지·claim 제거·24h metadata 제거 |
| Remote RED/GREEN | 새 도착 wake 회귀 exit1→0. 기존 Remote/TextEvent focused exit0 |
| Inbox RED/GREEN | Node22 실제 파일 stat 오류 주입 exit1→0. unread 본문·ENOENT·중복 보존, owner wire 회귀 exit0 |
| exact Linux install/build | `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`, 각각 native exit0 |
| 첫 SQL integration | native exit1. adapter build 완료 전 `connect.js` 부재. [원 로그](gate/integration-before-build-failed.log)를 보존하며 제품 결함 통과로 재분류하지 않음 |
| 재실행 SQL·DB 벤치 | 동일 제품 후보의 고유 Postgres17에서 전체 `go test -race -tags=integration ./internal/relay -count=1` native0, 111.338초. 전후 벤치도 native0이며 [전체 로그](gate/kl-refactor002-bcb1dd6629b5-0a7414c0-integration.log)의 set-e 체인으로 실패를 전파 |
| 최종 FullOps 등록 gate | 최종 archive 커밋 뒤 `lint.py --from 8f9370bb62c9133c73e10ea2307b4cd696e10f7e`. 정확한 최종 SHA·native exit·ERROR/WARNING은 외부 final-lint.json과 worker_done에 기록 |

제품 후보 FullOps 첫 시도는 로컬 문서 dirty 때문에 실행 전 native2로 거부됐다([원문](gate/bcb-dirty-gate-failed.out)). 제품 후보 검증은 외부 `refactor002/bcb-gate`의 깨끗한 exact bcb clone에서 다시 실행한다. Windows 긴 경로로 누락된 과거 증거 파일은 clone에만 `git -c core.longpaths=true restore --worktree .`로 복원했고 clean을 확인했다. 사용자 작업 트리 변경을 버리거나 파일 정리를 한 것이 아니다.

clean bcb의 등록 product-lint/product-test는 native0이다([FullOps 결과](gate/bcb-clean-fullops-lint.json)). FullOps 전체는 `inbox.pending()`를 테스트 보류 함수 `pending()`으로 오인한 ANTI-004 하나로 native1이었다. unread 파일의 실제 본문을 바로 다음 readFile assertion이 확인하므로 중복 목록 assertion 한 줄만 제거했다. 테스트 skip·설정 억제·검사기 변경은 없다. 제품 구현은 bcb와 동일하며 최종 SHA에서 전체 gate를 다시 실행한다. DEP-001은 의존성 아닌 회귀 등록이고 SIZE-001의 remote_events는 기존 463줄 파일에 같은 상태 전이 8줄을 추가한 것으로, 독립 추상화 분할보다 한 경로 유지가 적절하다.

Linux 기본 Go launcher는 1.22.2이고 module의 auto toolchain은 Go1.27.1이다. Node22.22.2/npm10.9.7을 사용한다. helper는 추적 HEAD bundle만 새 `/tmp/knowslink-refactor002-*`로 전송하고 매 make 전에 HEAD/clean을 검사한다. 고유 loopback Postgres17 컨테이너에만 synthetic credentials·DB를 만들고 migration Up·tagged race를 실행한다. 운영 DB·계정·Tunnel은 사용하지 않는다.

등록 `make lint/test`는 외부 task-local `make.cmd`가 같은 SHA의 Linux clone에 전달한다. 실제 원 명령의 종료코드가 FullOps에 돌아오며 설정 정본·gate를 완화하지 않는다. build·SQL·벤치의 고정 제품 후보와 문서만 추가한 최종 SHA를 구분한다.

### 실제 DB 전후 측정

같은 Linux AMD Ryzen 3 3200G, Postgres17 digest 고정, DB CPU 1개·메모리 512MiB 제한, loopback 연결을 사용했다. 같은 DB에서 이전 코드 `9255852fa0e46cd40b7be1021eaa1f6c7f7d083b`를 먼저 실행하고 후보를 실행했다. baseline에는 후보와 동일한 두 retention 테스트 파일만 복사했다. schema·제품 코드는 baseline 그대로다. 각 benchmark는 자기 합성 상태를 다시 seed하며 `-benchtime=5x -count=3`을 사용했다.

| 보관 이력 | DB transaction 이전/이후 ns/op 중앙값 | 이전/이후 B/op | 이전/이후 allocs/op |
|---:|---:|---:|---:|
| 100 | 17018286 / 14334587 | 1198078 / 1387009 | 3359 / 3284 |
| 1000 | 82354063 / 72626928 | 8707016 / 8485273 | 12540 / 11549 |
| 10000 | 588663216 / 542726660 | 80896920 / 79062267 | 103326 / 93329 |

10000개에서 실제 transaction은 약 7.8% 감소했다. 같은 Linux CPU-only sweep은 57910169→915579ns/op다. SQL·JSONB 변환·저장 비용은 남으며 한 요청의 입장/업무/퇴장을 합친 HTTP 지연이 아니다. 여러 client의 동시 처리량·최대 본문·장시간 포화 검사는 하지 않았다. 5회씩 세 번인 짧은 표본이며 실행 순서·공유 호스트·DB cache·GC 효과를 완전히 격리하지 않았다. B/op는 Go client 할당이며 Postgres 메모리와 peak RSS를 뜻하지 않는다. 원본 [before](benchmarks/before-linux.txt)·[after](benchmarks/after-linux.txt)와 [중앙값](benchmarks/linux-medians.json)을 보존한다.

## 기록·소유권·잔여

D03/D05/D07/D09/D10와 docs/system-design.md를 실제 구현에 맞춰 갱신했다. D02는 제품 규칙 불변, D04는 외형 불변, D06/D08은 저장 필드·SQL/migration 불변으로 no_change다. 의존성 추가는 없고 package.json은 새 회귀 등록만 바꿨다. 제품 수정 3파일·테스트 4파일·테스트 등록 1파일로 범위를 좁혔다. 새 레이어·정규화·상시 cache·설정 옵션을 만들지 않았다.

packet의 비선택 179 path/category와 partial/unknown 보완을 outcomes에 기록하고 정식 completion 검사로 확인했다. 제목·날짜·검색 구간의 과거 high/model 일치를 신규 계약으로 읽지 않았다. 현재 제품 규약과 날짜 있는 운영 관측을 분리했다. 원본 packet은 정식 history에 보존했다.

원 helper 생성은 shell guard가 부재 검사 문자열을 삭제 코드로 오인해 차단했다. [원문](guard-blocked-command.txt)을 보존했다. coor의 `msg_0eb2c7322744`가 직접 검토 후 비파괴 파일 편집만 허용했다. 거절된 명령을 재실행하거나 hook을 비활성화하지 않았다. helper의 compile·RejectPolicy·HEAD/clean·민감값 가림을 확인했다. 이 판단은 원 명령의 자동 guard PASS가 아니다.

하위 둘은 requested/effective codex/gpt-6-astra/high였다. adapter task `task_9e4b0f1970a2`/dispatch `ctx_fd1e99cea542`는 완료 뒤 release했다. remote task `task_516a9e21d9b4`/dispatch `ctx_89868aaae53b`는 변경 동결 뒤 부모 통합 gate를 기다린다. 부모가 최종 결과를 전달하고 모든 하위를 정착시킨 뒤 worker_done 한다. 별도 독립 리뷰는 coor가 제품 후보의 다른 실제 Astra high 세션에 배정했다.

임시 자료는 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002/`, 두 adapter 합성 fixture, 원격 고유 clone·DB에 보존한다. 성공 DB는 정확한 이름으로 stop만 하며 실패 DB는 보존한다. 임의 파일/디렉터리 삭제·이동 정리는 실행하지 않았다.

제품 후보의 clone은 `/tmp/knowslink-refactor002-bcb1dd6629b5-16beef99/repo`, baseline은 같은 상위의 `baseline`이다. 성공 DB `kl-refactor002-bcb1dd6629b5-0a7414c0`은 stop 상태로 보존했다. 첫 실패 DB `kl-refactor002-bcb1dd6629b5-5a295562`는 남아 있으며 coordinator 정리 대상이다. 최종 gate clone은 외부 최종 SHA의 `.base` 파일로 추적한다.

UI 외형 변화가 없어 새 캡처·직접 시각 판정은 해당 없음이다. 실계정 발송·live OAuth/callback·실제 host 화면·운영 배포·장시간 수신·최대 부하·백업 복원은 수행하지 않았다. singleton JSONB·행 lock과 local claim 이후 crash window는 남는다. 독립 QA·최종 리뷰 수락·main/origin 통합은 coordinator 후속이다.
