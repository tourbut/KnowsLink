---
title: 완료 detached 리뷰 워크트리 정리
status: review
updated: 2026-10-14
owner: coor
tasks: [SAR-REFACTOR-002-DEV]
summary: 사용자 승인으로 완료 리뷰 공간 8개를 제거하고 커밋·증거·상설 역할을 보존했다.
---

# 정리 결과

사용자는 detached HEAD 정리를 명시적으로 요청했다. coordinator가 완료된 리뷰용 Git worktree 8개만 Orca `worktree rm`으로 제거했다. force와 archive hook은 사용하지 않았다. 모든 대상은 KnowsLink workspace 바로 아래의 확인된 절대경로이며 일반 디렉터리다.

[사전 검사](preflight.json)는 정확한 HEAD·생성 기록·정본 결과의 동일성, main 조상, detached 상태, 변경/미추적/ignored 파일 0, 터미널 0, 대상 경로를 사용하는 process 0과 증거 SHA256을 기록한다. [첫 응답](removal.json)과 [나머지 응답](removal-rest.json)은 비동기 `removing` 영수증이다. 첫 응답 직후 디렉터리가 남아 있어 즉시 완료로 처리하지 않았으며, 이후 실제 파일과 Git/Orca 목록을 다시 확인했다. [최종 결과](final.json)는 8개 경로의 제거 완료를 기록한다.

정리 뒤 Git와 Orca에는 main/coor/designer/dev/ops/tester 상설 체크아웃 6개만 남았다. detached는 0개다. 커밋과 기존 QA/리뷰 정본은 삭제하지 않았다. 사용자 소유 터미널·미확인 역할 공간·레포 밖 임시 검증 자료는 건드리지 않았다.

# 자동 정리와 구분

FullOps `review.py cleanup`은 해당 snapshot에서 실행된 reviewer dispatch의 정상 release를 요구한다. 이번 과거 root 리뷰에는 그 dispatch가 없다. 사용자 요청과 coordinator의 직접 비사용 확인으로 Orca 운영 정리를 실행했으며, 가짜 dispatch/release를 만들지 않았다. 공용 snapshot 기록에는 `removed_user_authorized`와 이 근거를 남겼다. 자동 cleanup 또는 historical check 성공은 주장하지 않는다. 과거 검토 결과의 판정은 바꾸지 않았다. 재검수가 필요하면 보존 SHA로 새 snapshot을 만든다.

# 후속

새 Astra 구현과 별도 세션 검수는 새로운 과제다. 새 리뷰 snapshot은 완료 뒤 reviewer 종료·비사용과 증거 보존을 확인해 별도로 정리한다. 운영 서비스·계정·DB에는 이번 정리로 변경한 것이 없다.
