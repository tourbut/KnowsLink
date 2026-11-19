---
title: 최종 메시지 후보의 독립 QA 중단 기록
status: draft
updated: 2026-10-07
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER]
summary: 후보 d089의 lint·unit/race·build는 통과했으나 runtime 재빌드와 QA 컴파일이 실패하여 PS08–11은 미검증으로 마감한다
---

# SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER — 사용자 중단 지시로 QA 미완료

## 판정

고정 후보는 `d08903a55c3638128827010400e66e9d45b61d7c`다. 독립 PS08–11 기능 QA는 **미완료**다. 새 제품 적합성 PASS를 선언하지 않는다.
coordinator 메시지 `msg_0aa33027c858`는 사용자 요청에 따라 새 테스트·전체 재실행·탐색 확대를 중단하고 현재 결과를 기록하도록 지시했다. 해당 지시 수신 후 추가 기능 실행을 하지 않았다.
후보의 `make lint`, `make test`, `make build`는 각각 exit 0이다. 이 세 결과는 제품 기본 검사다. 전체 integration·교차 owner 왕복·신규 포화 중 자기 정리의 독립 판정을 대신하지 않는다.

## 독립성과 원본

실제 `CODEX_THREAD_ID`는 `01a116a7-1869-7503-9c40-e5c45acfae97`다. 세션 metadata의 모델은 `gpt-6.1-sol`이고 effort는 `high`다.
구현 세션 `aa85544d-18c3-43d3-95d0-b729aa9e9e8c` 및 최종 reviewer `01a116a3-2a1d-7961-92fa-33417bc9c46c`와 다르다.
관리 snapshot은 `/home/shin/orca/workspaces/KnowsLink/.fullops-review-a72ebe9890c04c478171375d18be5239`다. 시작과 종료 모두 위 SHA·detached·clean이다. 읽기 전용으로 사용했고 coordinator의 후속 회수를 위해 남겼다.
실행 scratch는 `/tmp/sar-messages-fix3-tester-88aqgkr5`다. 후보 detached clone에서 실행했다. 종료 뒤 자기 scratch를 회수했다. 제품 소스는 수정하지 않았다.
원 TESTER09c·실패 로그와 FIX2 Sol 진행분·원 FIX-REVIEW의 78개 파일은 SHA256과 byte 수로 보존했다. `original-hashes.json`과 `original-preserved.json`을 연결한다.
원 Grok 실패·중간 프로브 오류·미실행·OPS dfc H2 high·UI09c 실패·FIX2 pending·플랫폼 차단은 소급 변경하지 않는다.

## 실제 실행

모든 로그·명령 exit는 [SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER-test/](SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER-test/)에 있다. 파이프의 마지막 명령으로 exit를 대체하지 않았다.

| 명령 | 실행 위치·대상 | 실제 결과 | 근거 |
|---|---|---|---|
| `make install` | 후보 scratch | exit 0 | install.log / install.exit |
| `make lint` | 후보 scratch | exit 0 | lint.log / lint.exit |
| `make test` | 후보 scratch | exit 0, Go race·adapter unit | test.log / test.exit |
| `make build` | 후보 scratch | exit 0 | build.log / build.exit |
| `make verify-mvp` | 후보 scratch | exit 2, integration 시작 전 Node 재빌드 종료 139 | verify-mvp.log / verify-mvp.exit |
| 독립 `run_probe.py` | 후보 scratch·자기 PostgreSQL | exit 1, Go QA 컴파일 실패 | probe-runner.log / probe-runner.exit |
| fixture up / migrate build / migrate / down | 고유 `knowslink-messagesqa-*` | 각 exit 0 | fixture-logs/의 log와 exit |

`make verify-mvp`의 결정적 줄은 `Segmentation fault (core dumped)`와 `make: *** [Makefile:33: build] Error 139`다. 재빌드 전의 별도 build 성공을 runtime 성공으로 바꾸지 않는다. 원인은 이번 중단 범위에서 분석하지 않았다. 제품 결함으로 단정하지 않는다.
독립 QA 컴파일의 결정적 줄은 `internal/relay/qa_fixed_normal_test.go:9:2: "net/http" imported and not used`다. tester가 만든 시험 소스의 오류다. 제품 실패가 아니다. 기능 assertion은 한 개도 실행되지 않았다. 재실행·오류 수정은 사용자 중단 지시로 하지 않았다.

## 요구값과 현재 관측

| 범위 | 새 후보의 요구값 | 이번 실제 관측·판정 |
|---|---|---|
| PS08 로컬 일반회원 왕복 | 이메일 fixture의 서로 다른 owner·agent, 명시 송신·수신·persist/ACK·관련 reply·receipt ID 일치 | 새 교차 owner 시나리오를 작성했으나 컴파일 전에 중단됐다. 미검증 |
| PS09/10 TTL·멱등·권한·CSRF·gate | TTL 180s, UTF-8 4096 bytes, 동시 같은 key 한 메시지, 타 회원·원문 부재·기한·정책 차단 | 원09c 프로브를 새 후보용으로 준비했으나 미실행. unit/race suite만 exit 0 |
| PS11 공유 경계 | 신규 16·정리 4, queue 100·claim 4·receipt 20000, rate와 재시작 회수 | 기본 suite exit 0. DB integration·독립 상한 판정은 미검증 |
| 원 H1 | 익명 slow body가 정리·신규 admission을 미리 점유하지 않음 | 기존 실제 TCP regression의 소스를 확인했다. 이번 integration은 미실행. 해소 선언하지 않음 |
| 원 M1 | 유효 owner·CSRF의 form deny는 신규 포화 중 정리 budget으로 성공 | form deny 기대를 준비했다. 실행 없음. 원09c 실패 보존 |
| HTTP 행 잔류 | 큰 receipt 거절·정리 뒤 admission 행이 남지 않음 | 원 기대 0을 유지했다. 새 실행 없음. 이전 medium 해소 미검증 |
| FIX3 동일 DB 시각 | 늦은 낮은 epoch가 최신 snapshot을 덮지 않음, 더 늦은 복원 commit은 갱신 | 기존 unit 포함 make test exit 0. 독립 값/epoch 관측 시험은 컴파일 전 미실행 |
| FIX3 타 instance·restart | 신규 포화 중 현재 자기 기록을 읽고 유효 정리를 허용 | 기존 source와 요구를 대조했다. 실제 DB 관측 미실행 |
| FIX3 owner 제어·agent ACK | agent ACK가 owner 철회를 막지 않음, 공정성 단위마다 1개 | DEV 기존 자동 PASS는 구현자 근거다. tester integration 미실행 |
| 정상 cleanup caller | ACK·cancel·invite deny·gate deny·unpair·key/agent revoke·logout/all은 권한·유효 상태를 지킴 | 실제 상태 assertion을 준비했으나 새 기능 판정은 미검증 |

`probe-adaptation.diff`는 원09c와 이번 시험 준비의 차이다. `admission` 필드 추가로 positional literal을 named fields로 바꿨다. 원래 shared cleanup 시험의 빈 ACK에는 유효 lease가 없다. 새 계약은 이 요청을 신규로 분류하므로 해당 부분은 기존 verified-record integration에 연결하려고 분리했다. 익명 slow helper의 슬롯 점유 전제도 새 계약과 달라 기존 실제 TCP regression으로 연결했다. 이 변경을 실행 PASS나 원 실패의 소급 해소로 기록하지 않는다.

## 회수·한계·인계

fixture `down --volumes`는 exit 0이다. Docker after 목록에 자기 project가 없다. QA용 Go source는 scratch에서 제거됐다. scratch 전체도 회수했다. Node/MCP는 QA 컴파일 실패 때문에 새 fixture에서 시작하지 않았다. 공유 서버·Tunnel·다른 역할 자원은 중지하지 않았다.
실메일·공개 서버·실24h·운영 부하·CPU/상태 크기·백업/복원·PS13·PS14·Grok Bot/다닷 계정·최종 노우↔다닷·designer 직접 시각 판정은 미검증이다. 기존 L1·LA·LB와 운영/vendor 제한을 유지한다.
coor는 사용자 지시의 배포와 독립 기능 검수 미완료를 구분한다. 새 기능 QA를 재개하려면 별도 지시가 필요하다. 준비한 시험의 unused import와 실제 DB 검사를 먼저 확인해야 한다. 이번 tester는 제품 코드를 고치지 않았다.

## 기록 마감

`fullops-common-0.3.3`, project, document-writing, PS08–11/PS04·06·07, UX06·07을 기준으로 기록했다. 패킷의 partial·unknown·미전송 verify_mvp.py는 허용된 소스 확인과 중단 경계를 packet-outcomes에 기록했다. 추가 후보 읽기와 수정은 중단했다.
D10은 이 중단 보고서·시나리오 연결만 추가했다. D12/D13·공개 수락·제품 기준은 바꾸지 않았다.
`work.py finish`로 지시서와 완료 보고 전문을 보존하고 인박스를 비운다. 최종 커밋 SHA의 필수 FullOps lint/test·strict·diff·역할 push는 레포 밖 기록과 실제 worker_done으로 전달한다. 이 절 자체는 최종 SHA 검사 통과의 주장으로 사용하지 않는다.
