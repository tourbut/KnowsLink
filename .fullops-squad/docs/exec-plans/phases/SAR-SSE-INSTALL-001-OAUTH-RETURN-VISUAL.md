---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL 실행 계획
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL]
summary: tester가 만든 같은 후보 캡처 5장으로 OAuth 반환·오류 화면의 사용성을 직접 판정한 실행 기록
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL 실행 계획

## 목적과 범위

DEV가 OAuth consent 승인 뒤 303 대신 반환 화면을 보여 주도록 고쳤다(제품 SHA `eac603b6a41d2ce002e698600ac66b8fc67ad762`). designer는 tester의 같은 후보 캡처를 직접 보고 반환·재시작 안내가 사용 가능한지 판정한다. 새 기획·스타일·화면은 만들지 않는다. 기능 QA와 고정 SHA 코드 리뷰는 tester가 이미 수행했다.

## 단계와 결과

1. 지시서와 packet을 읽고 `visual-inputs.json`의 5장과 SHA256을 대조했다. 일치.
2. 5장을 이미지 도구로 직접 열었다. 반환 desktop·mobile, 만료 desktop, CSRF mobile, 계정 binding mobile.
3. 화면 문구를 고정 SHA의 `remoteReturn`·`refusal` 정의와 대조했다. 기준 `c399a5e`에는 반환 화면이 없다.
4. 관찰별 PASS와 한계를 [QA 보고](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md)에 기록했다.
5. 제품 바이트가 고정 SHA와 같음을 `git diff`로 확인하고 원 제품 검증을 재사용했다.
6. 등록 lint와 `work.py finish`로 지시서·보고를 보존했다.

## 결정

- 반환·안내 화면의 desktop·mobile 시각 수락: PASS.
- F2(링크 24px): 시각 수락 가능. 실제 클릭 성공과 24px 최소 기준 충족이 근거다.
- F1(medium): 비차단 잔여. DEV 후속에서 원인별 문구를 다룬다.
- V1(low, 모바일 한국어 단어 중간 줄바꿈): 기존 `memberStyle`의 전역 규칙이므로 이번 변경 범위 밖이다.

## 제한과 재개 조건

- 합성 Chromium 캡처 5장만 판정했다. 실제 호스트·계정·운영 복귀는 미확인이다.
- CSRF와 계정 binding 캡처는 같은 파일이며 조건 구분은 tester 보고에 의존한다.
- 재개 조건: 배포 뒤 사용자 계정으로 실제 호스트 복귀를 확인한다. 담당은 부모와 OPS다.
