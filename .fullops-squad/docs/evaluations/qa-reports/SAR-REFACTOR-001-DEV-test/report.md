---
title: 기능 보존 리팩터링 DEV 검증 보고
status: draft
updated: 2026-10-13
owner: dev
tasks: [SAR-REFACTOR-001-DEV]
summary: 공통 처리 네 곳과 핵심 경계 회귀 및 exact SHA 제품 게이트를 추적한다.
---

# SAR-REFACTOR-001-DEV 검증

## 범위와 기준

기준은 `bc0fd031efd252d9755681d8b429eb31c045b7e4`와 공통 `fullops-common-0.3.3`이다. Test level은 lite다. 이는 DEV 자체 검사이며 독립 tester QA·fixed-SHA 리뷰·제품 수락을 대신하지 않는다.

화면 agent 선택·연결 응답·일반/trial text 검증·MCP 응답 포장을 공유했다. 제품 코드 12줄 감소와 별도 경계 회귀를 기록했다. 상세 보존 이유는 [실행 기록](../../../exec-plans/phases/SAR-REFACTOR-001-DEV.md)을 따른다. API·서명·auth·UI·운영설정·SQL·의존성은 변경하지 않았다.

## 자체 검사

| 검사 | 결과 |
|---|---|
| Go focused 회귀 | child가 종료코드 0 확인. 응답 우선순위·retry header·agent owner/철회/차단/정렬/키 없는 표시 |
| `go test ./internal/relay -count=1` | child가 종료코드 0 확인 |
| `go vet ./internal/relay` | child가 종료코드 0 확인 |
| Windows Go race | 실행 전 실패: `go: -race requires cgo; enable cgo by setting CGO_ENABLED=1`. Linux 필수 회귀로 검증 예정 |
| adapter check/build·MCP·owner-text focused 회귀 | child 종료코드 0. Node24/npm11에서 수행한 보조 근거이며 최종 지원 Node22/npm10 회귀는 Linux에서 수행 |
| task-local SSH helper | Python parse·RejectPolicy·clean SHA 검사·출력 비공개 자기 검사 통과 |

## 고정 SHA 검사

후보 `cd4dae810ae262759bcff29895165f1a64f8c51a`의 고유 Linux clone `/tmp/knowslink-refactor001-cd4dae810ae2-45997998/repo`에서 아래 검사를 완료했다. helper는 매 make 호출 전에 exact HEAD와 `git status --porcelain`이 빈 상태인지 검사했다. 기존 결과를 재개 세션에서 직접 읽어 확인했으며 새 실행으로 표시하지 않는다.

| 검사 | 종료코드·근거 |
|---|---|
| `make install` | 0 · [install.log](gate/install.log), 기존 lockfile의 npm ci |
| `make lint` | 0 · [lint.log](gate/lint.log), Go format/vet/modules·adapter format/ESLint/types·config |
| `make test` | 0 · [test.log](gate/test.log), Go race·TS 전체 등록 suite·shell guard 4 tests |
| `make build` | 0 · [build.log](gate/build.log), Go relay/migrate·TS·plugin bundle |
| 설치 FullOps 1.4.1 lint, `--from bc0fd031efd252d9755681d8b429eb31c045b7e4` | 0 · [lint-candidate-141.json](gate/lint-candidate-141.json), exact cd4dae8, product-lint/product-test passed·ERROR 0·WARNING 3·unavailable 0 |

크기 경고 3건은 PLANS.md 1736줄, core.ts 334줄, human_test.go 517줄이다. 기록·회귀·공통 경계를 보존하며 줄 수 경고를 없애기 위한 분리는 하지 않았다. lint.log의 `cannot overwrite existing file`, `ERROR: candidate exists or cannot be written; live config unchanged`, `ERROR: invalid tunnel UUID`는 기존 부정 테스트의 예상 거절이다. 해당 명령 전체의 `EXIT=0`과 정상 테스트 요약을 함께 확인했다.

플러그인 갱신 중 기존 1.4.0 경로가 없어 최초 lint가 종료코드 2로 실패했다. [원 실패 출력](gate/lint-candidate.out)을 보존하고 설치된 1.4.1의 실제 통과 결과와 구분한다.

helper는 기존 검증된 SSH helper의 과제 사본이다. Git bundle은 추적 HEAD만 전송한다. 접속값은 메모리로만 읽고 known_hosts·RejectPolicy를 유지한다. 로그는 접속값을 가리며 터미널에는 log 파일명과 종료코드만 표시한다. 레포 밖 작업 경로는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor001/`이다. 소스 사본은 [gate/remote.py](gate/remote.py)에 있다. bundle·clone·기존 증거는 삭제하지 않는다.

## 중단 뒤 재개와 작성자 식별

원 구현 부모 세션은 `01a13509-e1a3-7b62-bf9d-5a48405d706a`, terminal `term_5ad12e64-a5e7-40af-9d78-ad9cba83b1c5`, dispatch `ctx_0ff7da5ff7ef`다. 실제 Go/TS 구현 자식 세션은 각각 `01a1350c-1383-7323-9d75-748304b33792`, `01a1350c-7f6b-7f03-822b-e3a3ec516f30`다. Git 작업 트리의 fullops-gate/flow 기록으로 대조했다. 원 세션의 candidate cd4dae8과 결과를 보존한다.

현재 완료 기록 담당 세션은 `01a1351f-c007-7f22-ac7c-1837dca51f0c`, terminal `term_188de7b2-366d-4aa7-9eb5-df99c7693447`, dispatch `ctx_2ee36546ed1f`다. coordinator의 재개 지시 `msg_92a1be5e891e`, 기존 gate 확인 `msg_6004ea668992`, 원 작성자 후속 인계 `msg_269c8b556bb9`를 적용했다. 최초 지시서의 병렬 구현 요구로 시작한 host 하위 둘은 최신 후속을 읽은 즉시 중단 지시를 받았다. 둘 모두 읽기 전용 조사만 했고 편집·검사·Git 작업은 하지 않았으며 종료했다. 현재 세션은 제품 코드를 변경하지 않고 증거·문서·역할 완료 기록만 보완한다.

## 생략과 후속

실계정·운영 배포·SQL 데이터 이행·부하·백업 복원·장시간 무인 호스트 wake는 실행하지 않는다. 관련 코드와 제품 규칙을 변경하지 않았다. UI template/CSS/문구 변경이 없어 직접 시각 검수를 추가하지 않는다. 기존 HTTP·template·MCP·Inbox·SSE·wake 자동 회귀는 등록 테스트가 수행한다. 독립 QA와 읽기 전용 snapshot 리뷰는 coordinator가 최종 후보에서 배정한다. 미해결 critical/high는 계속 수락을 차단한다.
