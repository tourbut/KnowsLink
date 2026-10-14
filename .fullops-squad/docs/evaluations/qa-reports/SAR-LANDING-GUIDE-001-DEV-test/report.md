---
title: 공개 안내 DEV 검증과 격리 QA 인계
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 인증 상태·공개 정보·실제 제품 검사와 후속 독립 QA를 구분한다
---

# 공개 안내 DEV 검증

비교 기준은 eae9a8bffae706410f684e4c837e396e7b0bcf26, 착수 HEAD는 b0f6596b6b7bfd4ea4e3b8cfed6d8b1e20a9b93e다. 구현과 기술 판단은 [DEV phase](../../../exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md)에 있다. 이 보고는 DEV 자체 검사이며 독립 QA·직접 시각 수락·코드 리뷰·운영 배포를 대신하지 않는다.

## 검사와 실패 보존

루트에서 `go test ./internal/relay -run 'TestLanding|TestMemberUIFixtures|TestGoogle|TestHuman' -count=1`을 실행했다. 수정 전 새 TestLandingConfiguredLoginAndPublicContent는 nil DB인 Service.Handler의 전역 admission 경로에서 panic했다. 이는 새 테스트 준비 오류다. 그 명령 뒤 orchestration check가 실행되어 셸 최종 코드는0이었지만 Go 결과는 명시적 FAIL이었다. 이 실행을 통과 또는 정확한 Go 종료코드 증거로 쓰지 않는다.

수정용 inline Python/PowerShell 명령은 자동 승인 검토의 `blocked by policy`로 거부됐다. 구체 사유는 없었다. 원 명령을 변형 재실행하지 않았다. coor ask msg_8b9dc2666080과 답변 msg_d75b88015752에 따라 명시적 apply_patch로 테스트 계층을 수정했다. coor가 원문 차단을 자기 외부 blocked-test-edit-evidence.json에 보존했다.

수정 후 같은 관련 단위 회귀는 Go 결과 `ok`와 명령 자신의 종료0을 확인했다. startPage 내용/provider/error 검사는 단위에서, Cache-Control/CSP는 실제 DB integration의 익명·유효세션에서 검증한다. 기존 form golden은 변경하지 않았다. 실제 제품 gate와 DB 결과는 고정 SHA를 연결해 아래에 추가한다.

## 사실·탐색·디자인 근거

README의 네 설치/업데이트 앵커와 installer Command/Arguments·신규 폴더 조건을 대조했다. 공개 예시는 합성 문구이며 실제 회원 데이터가 아니다. 수신·표시·답장 승인·180초·24h의 의미는 D02/D04와 구조/운영 문서를 유지한다.

[source-audit.json](source-audit.json)은 startup packet의61경로에서 원문 hash·제목·관련 구간을 보존한다. 과거 FAIL·미검증·partial/unknown과 잔여 후보를 유지한다. 초기 packet의 큰 D04/member.go는 관련 공개·인증 절을 직접 읽었다. 과거 보고의 PASS를 새 코드 실행으로 재사용하지 않는다.

frontend.py의 codex/gpt-6-astra/high check는 configuration_ready다. 인증·브라우저·MCP는 이 구성 검사로 통과하지 않는다. 주 스킬 frontend-design, 기존 memberStyle과 D04 여섯 토큰을 사용했다. 별도 디자인 lint와 다중테마는 없다.

## 격리 자원과 후속

새 검증 자료는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-dev/`에 보존한다. [gate/remote.py](gate/remote.py)는 known_hosts·RejectPolicy를 유지하며 현재 HEAD bundle만 새 `/tmp/knowslink-landing001-<SHA>-<suffix>`에 설치한다. `make.cmd`는 등록 명령의 make lint/test를 같은 HEAD Linux clone으로 전달하고 원 종료코드를 반환한다. 레포 lint 설정·필수 검사는 바꾸지 않는다.

합성 Postgres는 과제 고유 컨테이너다. 테스트용 DB는 검사 후 stop하고 보존한다. 브라우저 fixture는 별도 합성 DB·메일 mock·Google 취소 mock을 사용한다. 연결/송신 API는 fixture allowlist로 거부한다. 실제 로그인·실메일·실상대 송신·운영 DB·Tunnel·기존 컨테이너 변경은 없다. [gate/forward.py](gate/forward.py)는 지정 fixture의 loopback 포트만 Windows loopback으로 전달한다.

독립 tester는 실제 handler·POST·redirect·정보 비노출을 확인한다. designer는 같은 후보의 1440/390/320 실제 캡처를 판정한다. 별도 세션 리뷰와 main/origin 통합은 coor 후속이다. 임시 자료의 삭제·이동 정리는 수행하지 않는다.
