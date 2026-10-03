---
title: FullOps 0.9.13 업데이트 적용
status: review
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-0.9.13]
summary: 글로벌 gate 수정과 레포 버전 적용 및 사용자 작업 중단 보존을 기록한다
---

# FULLOPS-UPDATE-0.9.13

## 설치와 적용 기준

현재 CLI는 Codex다. 실제 설치 전 plugin.json은 0.9.13이다. 레포 적용 버전은 0.9.12다. marketplace upgrade는 이미 최신, plugin add는 0.9.13 설치 경로를 반환했다. 두 명령은 종료코드 0이다. 새 update-fullops를 직접 다시 읽었다. deps.py --check는 필수 CLI 모두 있음, 종료코드 0이다.

update.py --repo . --from 0.9.13은 더 오래된 레포 기준 0.9.12 이후 릴리스 0.9.13을 출력했다. 이전 0.9.12 적용 기록을 확인했으며 미적용 필수 레포 항목은 없다. 기존 역할 coor/designer/dev/ops/tester, origin/main으로 setup dry-run은 종료코드 0, 생성 예정 0개다. 신규 파일이 없어 setup 실실행은 해당 없음이다.

## 릴리스 판정

0.9.13의 셸 도움말·과제 키 구분·워크트리 selector·완료 보고 수집·hold/resume 재조회·종료 확인 수정은 글로벌 플러그인 설치로 적용한다. 서비스 코드·역할·기본 브랜치 변경은 해당 없음이다. integration status pending은 0건이다. 빈 SHA·과제 키 복구 대상이 없으므로 기존 보고 재수집은 해당 없음이다. 기존 hold·실패·QA 제한과 사용자 중단 지시를 보존한다.

업데이트 요청 직전 Grok 어댑터 준비 커밋 99c0aaf만 작성했다. 실제 worker-start를 실행하지 않았다. 사용자의 작업 중단 요청을 인박스 blocked·PLANS·board에 반영한다. 제품 개발·외부 발송·배포는 진행하지 않는다.

## 검증과 재발 확인

실제 plugin.json 0.9.13과 deps check, setup dry-run, integration status를 확인했다. orchestration check는 현재 터미널이 기존 Run에 bound되지 않아 종료코드 1을 반환했다. 이 오류는 이번 릴리스의 SHA/키 파싱 오차단과 구분하며 성공으로 표시하지 않는다. 정상 새 세션 바인딩 전 Run 변경을 수행하지 않는다.

설치 버전과 필수 적용 확인 후 fullops.json의 plugin_version을 0.9.13으로 갱신한다. 제품 코드 변경이 없어 제품 기능 QA를 재실행하지 않는다. 문서·JSON·Git 공백·산출물 strict·FullOps lint를 실행한다. 수정된 gate 오류의 실제 worker 재발 검증은 새 coordinator 세션에서 다음 명시적 재개 후 확인한다.

## 통합과 보류

운영 변경은 허가된 main 일반 push로 공유한다. 깨끗한 main/coor를 동기화한다. 진행 중 또는 상태 불명 역할의 파일을 덮어쓰지 않는다. 해당 역할은 최신 main 동기화를 예약한다. 개발 재개는 사용자의 중단 지시가 해제되고 새 coordinator 세션을 시작한 뒤 수행한다.

검증 SHA 1f938c3389ea의 FullOps lint는 종료코드 0, product-lint passed, ERROR 0, WARNING 1, 실행 불가 0이다. 경고는 기존 누적 PLANS 503줄의 SIZE-001이며 기록을 삭제하지 않고 보존했다. strict는 13개, 문제·경고 0이다. Git 공백 검사도 통과했다. 이 업데이트의 제품 코드 변경은 없다.
