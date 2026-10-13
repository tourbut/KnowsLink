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

구현 후보 커밋 뒤 고유 Linux clone에서 `make install`, 등록 `make lint`, `make test`, `make build`를 수행한다. 결과는 완료 보고 전에 실제 종료코드로 갱신한다.

helper는 기존 검증된 SSH helper의 과제 사본이다. Git bundle은 추적 HEAD만 전송한다. 접속값은 메모리로만 읽고 known_hosts·RejectPolicy를 유지한다. 로그는 접속값을 가리며 터미널에는 log 파일명과 종료코드만 표시한다. 레포 밖 작업 경로는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor001/`이다. 소스 사본은 [gate/remote.py](gate/remote.py)에 있다. bundle·clone·기존 증거는 삭제하지 않는다.

## 생략과 후속

실계정·운영 배포·SQL 데이터 이행·부하·백업 복원·장시간 무인 호스트 wake는 실행하지 않는다. 관련 코드와 제품 규칙을 변경하지 않았다. UI template/CSS/문구 변경이 없어 직접 시각 검수를 추가하지 않는다. 기존 HTTP·template·MCP·Inbox·SSE·wake 자동 회귀는 등록 테스트가 수행한다. 독립 QA와 읽기 전용 snapshot 리뷰는 coordinator가 최종 후보에서 배정한다. 미해결 critical/high는 계속 수락을 차단한다.
