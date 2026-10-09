---
title: <과제 키> 리뷰
status: draft
updated: 2026-10-09
owner: coor
tasks: [FULLOPS-UPDATE-1.2.0]
summary: 독립 고정 SHA 리뷰와 탐색 처리·대화 미참조 인계 점검 양식
---

# <과제 키> 리뷰

- 검토자 / CLI / 모델:
- base SHA / head SHA / merge-base:
- OCR 버전 / 적용 규칙:
- 요구사항·완료 기준 원천:
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped:
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유:
- SIZE-002: 지시서의 예상 변경 규모 / 실제 추가 줄 수 / 차이·분할하지 않은 이유:
- DEP-001: 파일별 의존성 변경 여부 / 완료 보고의 필요성·표준 라이브러리 대안:
- UI 디자인: 핸드오버의 정본·공용 컴포넌트·예외 근거 / 디자인 lint의 HEAD·종료코드 / DESIGN 경고 처리 / 테마 전환 검증·해당 없음의 이유:

## 검토 범위

result.json의 모든 (path, status)에 검토 상태를 기록한다. 제외 파일도 검토하거나 구체적인 생략 사유를 남긴다.

## 발견 사항

심각도 / 파일·줄 / 재현 조건 / 영향 / 수정 상태 / 검증 근거.
발견 사항이 없으면 실제 검토 범위와 검증한 내용을 적는다.

## 검증 및 남은 제약

실행한 테스트·결과, 실행하지 못한 검증, skipped 영향과 남은 일을 적는다.
등록된 kind: test 명령의 HEAD·종료코드와 실행 결과를 확인한다. SIZE-002의 예상 규모나 DEP-001의 변경 근거가 누락되면 수정 요청을 남긴다.
UI 작업의 디자인 기준·검증 항목이 비어 있으면 수정 요청을 남긴다. 디자인 경고와 테마 전환의 수락 여부는 검토자가 판단한다. review.py check가 이 문서의 서술을 자동 판정하지는 않는다.

## 검토 결론

수락 가능 여부와 근거를 적는다. check 통과는 AI 검토의 내용·테스트 성공을 자동 보증하지 않는다.

## 독립성과 인계 점검

- 구현자 실제 세션 ID / 검토자 별도 실제 세션 ID:
- 관리된 clean detached snapshot 경로 / 고정 40자리 SHA / 읽기 전용 확인:
- 현재 attempt·packet의 path/category별 처리 근거 / partial·unknown 해소:
- 대화 미참조 점검: 정본 인덱스에서 요구·결정 이유·구조·구현/미완료·실행/검증·운영/복구·다음 작업을 찾은 경로/절과 확인·미확인·해당 없음:
- 누락·오래된 정보·깨진 링크·미지원 anchor / skipped 사유와 영향:
- 테스트 레벨 / 필수 검사와 skipped 사유 / 원래 SHA의 재사용 증거:
- 검토·증거 보존·reviewer release·관리된 snapshot cleanup 결과:

historical check는 보존 증거 확인이며 신규 SHA의 수락 리뷰를 대신하지 않는다.
