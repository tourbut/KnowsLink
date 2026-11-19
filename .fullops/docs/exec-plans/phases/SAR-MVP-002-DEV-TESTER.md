---
title: SAR-MVP-002-DEV-TESTER — 플러그인 독립 QA 실행 기록
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-002-DEV-TESTER]
summary: "고정 후보 552586b의 패키지, stdio MCP, 격리 SQL 독립 검증을 기록한다"
---

# SAR-MVP-002-DEV-TESTER — 플러그인 독립 QA 실행 기록

## 기준

기준 ref는 `dbdd70086971285b790683f362702e5a9ff55acd`다. 판정 후보는 `552586b6e886f95bffa9a000a031ea03070afedb`다.
기록 브랜치는 `fullops/tester`다. 시작 HEAD는 `e8cf6cf531142bc1b440d46d527a7f6720dc317f`다. 실행 clone은 이 체크아웃과 분리된 detached HEAD다.
`fullops-common-0.3.2`, FULLOPS.md, project.md, 문서 작성 규칙, fullops-test를 적용했다. D03, D05, D10과 사용자 설치 README를 판정 기준으로 읽었다.

## 실행

1. clone에서 `make install`을 실행했다. 종료코드는 0이다.
2. `make plugin`으로 ZIP을 만들고 압축 해제 MCP 검사를 포함했다. 종료코드는 0이다. SHA256은 `0e671d1a89c141d896034fff31619b9cd2148b73b567adbc3a97126031989117`이다.
3. 압축 해제본에 `mcp.test.js`와 동시 pull probe를 따로 실행했다. 두 종료코드는 0이다.
4. `make verify-mvp`로 격리 Compose, 실제 SQL, MCP bundle pull, owner approve, denied result를 실행했다. 종료코드는 0이다.
5. `package_plugin.py` 변경 때문에 `make lint`를 clone에서 다시 실행했다. 종료코드는 0이다.

Go race `make test`의 `./cmd/... ./internal/...`는 `c4ebbecd3ef91be10ecbb517fe451e02769358bc` 증거를 재사용했다. 대상 소스 diff 종료코드는 0이다.
Chrome QA `9584aafcbb5fee88dcc6d618caf660884f6a527d`를 재사용했다. Go와 UI 파일 diff 종료코드는 0이다. 이번 기록은 새 브라우저 QA가 아니다.

## 판정과 보류

패키지 허용 목록, held, loopback, redirect, mock 실패 전파, 실제 SQL 권한, gate, 최소 denied result, in-process `busy`와 SQL 동시성 검사는 통과했다.
새 결함과 새 critical/high는 없다. 제품 코드, 공유 fixture, 기존 컨테이너, 운영 설정은 변경하지 않았다.
실제 Grok Bot 계정, marketplace, hosted runtime, 유료 API는 held다. 재개 조건은 계정과 도달 경로의 별도 승인다.
상세 명령과 로그는 [QA 보고서](../../evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md)에 있다. 독립 리뷰는 coor 후속이다.

## 완료 뒤 검사

FullOps lint는 기준 ref에서 깨끗한 기록 커밋으로 실행한다. `deliverables.py --strict`도 그 커밋에서 실행한다. 종료코드는 QA test 디렉터리에 남긴다.
`work.py finish`가 인박스와 완료 보고를 보존한다. 역할 브랜치 일반 push 뒤 worker_done을 보낸다.
