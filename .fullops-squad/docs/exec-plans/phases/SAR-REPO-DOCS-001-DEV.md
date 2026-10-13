---
title: README와 루트 docs 현행화 실행 기록
status: draft
updated: 2026-10-13
owner: dev
tasks: [SAR-REPO-DOCS-001-DEV]
summary: README 설치·업데이트와 루트 docs 작성·검증·재개 근거
---

# README와 루트 docs 현행화 실행 기록

## 기술 계획

1. 운영 제품 `2747fe98653d5c6946f41c219fc256395333d7fe`와 실제 친구 대화·수신 알림 기록을 현재 상태 근거로 사용한다.
2. README에 사용자 흐름, Grok 로컬 신규 설치/기존 업데이트, Dots 원격 신규 연결/기존 업데이트를 분리한다.
3. 루트 docs에 독자별 목차·구조·개발·운영 설명을 작성한다. 설치 본문은 README로 연결한다.
4. adapter 안내의 현재와 충돌하는 절만 정정한다. 합성·legacy·다른 호스트 검사 한계는 보존한다.
5. 실제 소스의 명령·환경·저장·writer 계약과 문서를 대조한다. 링크·fragment·코드펜스·공백·산출물 strict를 검사한다.
6. 탐색 원추천과 수동 보완을 outcomes에 기록한다. 인박스 의미는 변경하지 않고 finish·커밋·표준 회신한다.

## 적용 범위와 경계

기준 `473da85ad6edc56d3833fd42cf540e8d7887920a`, 착수 HEAD `a2961818ee33f3786aa6f949c5b411bd4b31aa1a`, 공통 규칙 `fullops-common-0.3.3`을 적용한다. 문서와 완료 증거만 변경한다. 제품 코드·의존성·운영 설정·실계정 발송·배포·삭제·이동은 수행하지 않는다. 하위 위임은 없다. 작성자 세션과 검증 결과는 QA 보고서에 기록한다.

기존 Context7 월 quota 실패와 OPS의 2026-10-13 공식 근거를 재사용한다. OpenAI Docs 스킬로 Add MCP/Refresh 및 MCP Events 공식 원문도 열어 확인했다. 외부 계정별 UI·요금제·권한은 확정하지 않는다.

## 결과와 후속

검증 결과와 한계는 [QA 보고서](../../evaluations/qa-reports/SAR-REPO-DOCS-001-DEV/report.md)에 기록한다. coor가 작성자와 다른 실제 세션에서 고정 SHA 문서 리뷰·main 통합/push를 수행한다. 미해결 critical/high 차단을 유지한다.

## 재개와 완료

실제 작성자 세션은 `01a13599-947e-7253-b385-f596fc4b0e2d`다. 이전 세션 `01a1358b-dde4-7e80-97de-8c6d335de28a`의 phase·README 작업 사본을 보존했다. coor의 정상 stop/release 뒤 같은 인박스·기준·Sol high에서 이어 작성했다.

README와 루트 docs4개, adapter의 충돌 절을 작성했다. 로컬 링크73개·기존 inbound12개·exact 소스항목20개를 검사했다. 구 heading alias를 보존했고 제품 파일은 바꾸지 않았다. packet189경로192쌍·remaining155의 관련 검색/원문 확인과 unknown19/producer4와 직접원문보완4 구분·보완은 QA audit/outcomes에 남겼다. 인박스는 의미변경 없이 완료 결과만 추가한다. 새 설치/배포/계정/삭제/이동은 실행하지 않았다. work finish·커밋·고정 SHA 회신 뒤 coor의 독립 문서 리뷰/통합을 기다린다.
