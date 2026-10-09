---
title: FULLOPS-UPDATE-1.2.0
status: draft
updated: 2026-10-09
owner: coor
tasks: [FULLOPS-UPDATE-1.2.0]
summary: 설치 1.2.0과 레포 0.9.14의 차이 적용·검증·보류 기록
---

# FULLOPS-UPDATE-1.2.0

## 목표와 적용 기준

현재 Codex의 FullOps 설치와 레포의 누락된 운영 변경을 적용한다. 업데이트 전 실제 설치 버전은 1.2.0이다. 레포 적용 버전은 0.9.14이며 대상은 1.2.0이다. 기준 ref는 `cdcb3c12e2456136feebf6bf06aedc0e5d3de95d`다. 역할·브랜치·모델·원격·사용자 설정·진행 중 인박스·제품 실패와 보류를 보존한다.

## 설치와 사전 확인

`codex plugin marketplace upgrade fullops-squad --json`은 errors 없이 종료코드 0이다. `codex plugin add fullops-squad@fullops-squad --json`은 version 1.2.0과 실제 설치 경로를 반환했다. 새 패키지의 update-fullops와 INSTALL.md를 직접 확인했다. 출처 전환과 이전 캐시 삭제는 수행하지 않았다.

Python은 3.12.3이다. 초기 `deps.py --check --host codex`는 종료코드 1이며 host dependency installation incomplete였다. `deps.py --host codex`로 pinned OCR 1.12.12·Context7 4.1.1·외부 스킬과 Codex ponytail 구성을 설치했다. 설치는 종료코드 0이다. 재검사는 호스트 의존성 검증 완료, 종료코드 0이다.

`update.py --repo . --from 1.2.0`은 더 오래된 레포 기준 0.9.14 이후 릴리스 다섯 개를 반환했다. 기존 다섯 역할·origin/main의 setup dry-run과 실행은 생성/작성 0개, 종료코드 0이다. 기존 파일은 자동 교체되지 않았다. 관리 블록은 제공 setup의 operating_block과 동일한 내용으로 기존 문서에 통합했다.

main과 등록 다섯 역할은 시작 시 clean이다. Orca 조회에서 하위 네 역할과 main은 inactive·터미널 0개다. coor는 현재 업데이트 세션만 working이다. 미통합 결과는 integration status에서 pending 0개다. 과거 리뷰 공간·실패 dispatch·사용자 소유 터미널은 삭제하지 않는다.

## 릴리스별 판정

| 릴리스 | 판정과 근거 |
|---|---|
| 1.0.0 필수 | 적용. Python 3.10 이상과 Codex 호스트 의존성을 확인했다. 세션 선택·고정 기준 SHA·Task key/Purpose·reopen/attempt·route bind-inbox·탐색 패킷·전송 receipt·직렬 lint·독립 snapshot 수명·대화 미참조 인계 점검을 운영 문서와 양식에 반영했다. 진행 중 인박스를 일괄 변환하지 않는다. 다음 신규 인계부터 새 계약을 적용한다. |
| 1.0.0 조건부 | 공유 D06/D07/D09 원천은 해당 없음. 인덱스에 서로 다른 세 파일이 등록돼 복수 ID 변환이 필요 없다. Unity bridge는 Go/TypeScript 서비스이므로 해당 없음이다. 실제 Jev 호출은 이번 운영 변경의 검증에 필요 없어 실행하지 않는다. 새 리뷰 snapshot은 관리 도구로 만든다. 기존 상설·임시 사용자 공간은 보존한다. |
| 1.1.0 | 적용. 이슈 작업은 기본 OFF이며 사용자 요청 때만 켠다. 자동 이슈 과제는 draft PR·integration hold까지이며 main 병합은 사용자 판단이다. 활성화 요청이 없어 설정·enable·GitHub 댓글을 실행하지 않는다. |
| 1.1.1 | 해당 없음. 비공개 token-file은 이슈 모드 활성화 때 새 스킬을 따른다. 현재 lease token을 만들거나 전역 권한을 변경하지 않는다. |
| 1.1.2 | 해당 없음. 영상 요청이 없어 제작 도구·의존성·역할을 추가하지 않는다. 요청 시 fullops-motion을 사용한다는 포인터만 등록했다. |
| 1.2.0 | 적용. 기존 coor/standard를 유지하며 FULLOPS.md·orca-agents.md·testing.md의 관리 블록을 보강했다. mode 전환·테스트 축소·기존 lint 명령의 level 변경은 하지 않는다. 과거 pending 리뷰 supersede는 이번 운영 업데이트와 관계없어 기존 원본을 보존한다. |

## 변경 범위와 완료 조건

운영 문서·양식·보드·적용 버전과 업데이트 기록을 변경한다. 예상 검사 대상 추가 줄은 200줄 이내이며 제공 보드 HTML은 별도 생성물이다. 제품 코드·의존성 선언·검사 명령·역할 모델·제품 규칙은 변경하지 않는다. UI·테마 검수는 제품 화면 변경이 없어 해당 없음이다.

문서 strict·JSON·공백·관리 블록과 기존 설정 보존 검증 뒤 plugin_version을 1.2.0으로 갱신한다. 고정 커밋의 product-lint·product-test와 독립 리뷰를 통과한 뒤 main에 통합하고 원격 push한다. clean·inactive 역할만 동기화한다. 검증 실패는 적용 완료로 보고하지 않는다.

## 검증과 보류

제품 서비스 공개·Google 정책·외부 에이전트 왕복·기존 QA/UI 보류를 이번 업데이트로 수락하지 않는다. 이슈 실환경 wake·Unity Player·Windows/macOS·실제 영상·새 세션 hook은 이 Linux 업데이트에서 검증하지 않는다. 다음 개발은 새 coordinator 세션에서 진행한다.

필수 적용의 보류가 생기면 coor가 사유·담당·재개 조건을 PLANS.md와 이 기록에 남긴다. 역할 동기화는 레포 적용과 구분하며 실행 중·dirty·상태 불명 역할은 덮어쓰지 않는다.
