---
title: 기능 보존 전체 코드 정리 실행 기록
status: draft
updated: 2026-10-13
owner: dev
tasks: [SAR-REFACTOR-001-DEV]
summary: 모듈별 호출 근거로 중복을 정리하고 고정 SHA의 회귀 검증을 기록한다.
---

# SAR-REFACTOR-001-DEV

## 기준과 기술 계획

제품 기준은 `bc0fd031efd252d9755681d8b429eb31c045b7e4`다. 착수 HEAD는 문서 준비 커밋 `0eb7735f0963272b62266992791f9f09152f9266`다. 공통 규칙 `fullops-common-0.3.3`, 현행 HP-01–16, D03, D07을 적용한다. 과거 문서의 후보·미수락 설명은 당시 상태다. 최근 PLANS와 운영 기록의 실계정 왕복·자동 수신 근거를 보존한다.

1. 실제 함수 호출, template 메서드, CLI 진입점, MCP 도구 등록, bundle와 패키징 경로를 확인한다.
2. Go와 adapter에서 명확한 중복만 기존 공통 동작으로 합친다. 외부 계약과 보안 경계는 보존한다.
3. 부모는 cmd/config/scripts/deploy/루트 빌드 연결을 조사하고 하위 결과와 문서를 취합한다.
4. 변경을 커밋한 뒤 고유 Linux clone에서 exact SHA·clean 상태와 `make lint`, `make test`, `make build`를 확인한다.
5. 등록 FullOps lint, packet outcomes, D10, 역할 완료 기록을 마친다. 독립 QA와 고정 SHA 리뷰는 coordinator에게 인계한다.

Subagent level은 standard, Test level은 lite다. 하위 둘은 Codex `gpt-6.1-sol` / high로 실제 기동했다. relay 소유자는 `internal/relay/*.go`와 관련 tests만 수정한다. adapter 소유자는 `adapters/src/*.ts`만 수정한다. 부모만 문서와 Git 커밋을 처리한다.

## 실행 환경과 보존

하위 Run은 `run_aa0857932bbe`다. Go dispatch는 `ctx_07f05ebc8cfa`, adapter dispatch는 `ctx_a614cf6dca4f`다. adapter 최초 시작은 입력만 수락되어 turn 시작이 미확인이었다. 실제 draft를 확인한 뒤 Enter 한 번으로 기존 입력을 제출했다. 중복 worker는 만들지 않았다.

하위 Run 생성 뒤 부모 Run을 명시한 check가 Run binding 불일치로 `consumer_fenced`를 반환했다. coordinator는 기존 부모 dispatch의 owner·incarnation·live·미철회 상태를 독립 확인했다. coordinator 답변에 따라 하위 Run을 유지하고 부모 follow-up은 비소비 inbox로 읽는다. 부모 Run의 coordinator binding은 변경하지 않는다.

파일 삭제·이동, 운영 설정 변경, 배포, 실계정 발송, 로컬 Docker는 실행하지 않는다. coordinator는 고유 검증 clone에서 기존 필수 테스트가 자기 mktemp 자료만 정리하는 기존 trap을 승인 범위로 확인했다. 새로운 정리 명령은 추가하지 않는다. 원본 환경·키·기존 증거는 보존한다.

## 조사·변경 결과

| 범위 | 판단과 근거 |
|---|---|
| relay 화면 선택 | `humanView`의 자기·친구 agent 선택 두 곳을 `humanAgents`로 합쳤다. owner 필터·철회 제외·ID 정렬·키 없는 상태 표시를 유지한다. |
| relay 연결 응답 | `deviceAPI`와 `connectAPI`의 동일 retry·명령 오류·transaction 응답을 `connectionResponse`로 합쳤다. 기존 우선순위와 Retry-After를 유지한다. |
| adapter text | 기존 `validText`를 `core.ts`에서 공유한다. 일반·trial 송신과 수신의 공백·UTF-8 길이·고립 surrogate 검증을 보존한다. |
| MCP 응답 | 기존 `trialResult`의 동일 text content 포장을 `result`에서도 재사용한다. 상태 필드와 isError는 그대로다. |
| CLI·config | `cmd/relay`의 기동·shutdown·health와 `cmd/migrate`의 SQL-only 검사는 서로 다른 실패 경계를 가진다. `config.DatabaseURL`은 두 명령이 호출한다. 더 합치지 않았다. |
| scripts·deploy | Makefile이 검증·schema·패키징 스크립트를 호출한다. `run_trial.py`는 README의 외부 호스트 launcher이며 `mail_sink.py`는 격리 인증 fixture다. `beta.sh`가 `verify.py local`을 호출하고 Access trial 계획은 D12에 참조된다. 별도 CLI·권한·실패 동작을 유지했다. |
| 빌드·패키징 | Makefile·Dockerfile·package.json·HOST_FILES를 확인했다. `index.ts`는 일반 CLI 진입점이며 plugin bundle 진입점은 `mcp.ts`다. `plugin.ts`는 없다. 호스트별 manifest·허용 목록·설치 경로 차이는 필요하다. |
| 유지한 relay 후보 | `leaseMessage`는 legacy/public text tests가 호출한다. `legacyHomePage`는 route에 등록된다. `memberPair.Deadline/Mine/Peer`는 template이 호출한다. MCP tool·Events 등록과 registry/legacy fixture는 외부 계약이다. |
| 유지한 adapter 후보 | trial/public CLI·exports·Inbox·AutoReceiver·wake는 실제 entry 및 도구가 호출한다. 보존·ACK 순서·serial queue·SSE coalescing·retry/uncertain을 건드리지 않았다. |

추적 제품 소스의 `bytes.splitlines()` 기준으로 relay 6792→6786, adapter 2793→2787, 부모 범위 1792→1792다. 합계 11377→11365로 12줄 감소했다. 새 경계 회귀는 별도로 추가했다. 줄 수를 줄이기 위해 기존 테스트·문서·호환 경로를 삭제하지 않았다. [집계](../../evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/line-counts.json)를 따른다.

Go 자체 focused 회귀·relay suite·vet가 통과했다. Windows race는 cgo 비활성으로 실행 전 실패했으며 통과로 기록하지 않는다. Linux 고정 SHA의 등록 검사와 빌드 결과는 QA 보고에 기록한다.

## 패킷 수동 보완과 문서 영향

188개 경로와 193개 path/category 쌍의 존재·정본·검색 근거를 확인했다. [검토 지도](../../evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/packet-review.json)는 각 파일의 hash·줄 수·원 packet 근거 줄·현재 계약 검색 줄·제목을 보존한다. 과거 로그는 해당 구간과 이력을 확인하며 전체 전문 재독이나 새 SHA 재실행을 주장하지 않는다.

partial/unknown 17건은 producer의 bounded search·큰 문서·없는 `plugin.ts` 때문이다. 큰 PLANS/운영/핸드오버는 현행 사람 모델·writer guard·검증 명령·기존 결과 절을 로컬 검색으로 보완했다. `remote_events.go`는 실재하며 callback 검증·서명·구독·retry 등록을 유지했다. documents-find 잔여 후보와 optional/context 경로는 검토 지도에서 모두 추적한다. semantic language-server의 완전한 호출 그래프를 주장하지 않는다. 파일 삭제 후보는 확정하지 않았다.

D10에는 공유 처리와 회귀를 추가한다. D03/D05/D06/D07/D08/D11/D12/D13의 외부 계약·schema·운영·설치·사용자 동작은 바뀌지 않아 기존 원천을 보존한다. 과거 후보 설명과 최근 운영 수락의 차이는 시점 차이다. UI template·CSS·문구는 변경하지 않았다. 직접 시각 재검수는 해당 없음이며 기존 렌더 회귀는 필수 테스트에 포함된다.

## 고정 후보 검증과 완료 재개

제품 후보는 `cd4dae810ae262759bcff29895165f1a64f8c51a`다. Linux 고유 clone `/tmp/knowslink-refactor001-cd4dae810ae2-45997998/repo`에서 install/lint/test/build가 모두 종료코드 0이다. 매 make 호출은 HEAD 일치와 clean 상태를 먼저 확인했다. 설치 FullOps 1.4.1의 등록 product-lint/product-test도 같은 후보에서 통과했다. ERROR 0·WARNING 3·unavailable 0이며 경고는 기존 대형 파일 PLANS/core/human_test의 크기다. 증거 원본은 외부 refactor001 폴더에 보존하고 가려진 로그 사본과 gate JSON을 [QA 보고](../../evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md)에 연결했다. 플러그인 1.4.0 경로 소실의 최초 실패는 별도 보존했다.

원 구현 부모는 세션 `01a13509-e1a3-7b62-bf9d-5a48405d706a` / dispatch `ctx_0ff7da5ff7ef`다. Go/TS 자식 세션은 `01a1350c-1383-7323-9d75-748304b33792` / `01a1350c-7f6b-7f03-822b-e3a3ec516f30`다. 11:54 UTC operator terminal close 뒤 원 dispatch가 철회됐다. 현재 세션 `01a1351f-c007-7f22-ac7c-1837dca51f0c` / dispatch `ctx_2ee36546ed1f`는 원 후보와 실제 통과 결과를 이어받아 문서·finish를 완료한다. 최신 후속 확인 전에 시작한 host 하위 둘은 편집 없이 종료했다. 제품을 다시 구현하거나 원 실패를 성공으로 바꾸지 않았다.

원 packet outcomes의 193개 비선택 path/category와 partial 수동 보완을 유지했다. D10은 후보에서 이미 갱신됐으며 D13의 새 제품·운영 계약 변경은 없어 원문을 보존한다. 독립 tester QA와 작성자와 다른 고정 SHA snapshot 리뷰, main 통합·일반 push는 coordinator의 후속이다. 배포·실계정 발송·파일 삭제·이동은 수행하지 않았다.
