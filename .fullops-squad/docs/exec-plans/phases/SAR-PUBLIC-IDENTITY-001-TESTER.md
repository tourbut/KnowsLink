---
title: SAR-PUBLIC-IDENTITY-001-TESTER — 독립 QA 실행 기록
status: draft
updated: 2026-10-05
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-TESTER]
summary: 고정 후보의 독립 QA 순서와 F1 관측 및 남은 공개 검증을 기록한다
---

# SAR-PUBLIC-IDENTITY-001-TESTER — 독립 QA 실행 기록

## 입력과 범위

지시서는 `.fullops-squad/handovers/to_tester.md`다. 과제 키는 `SAR-PUBLIC-IDENTITY-001-TESTER`다. Run은 `run_8ca8bc058ab7`다.
고정 후보는 `59b66ada8b36802484cc6d7e22523257b50572cc`다. 기록 체크아웃 준비 커밋은 `746ecd9193e9283369267d51e110ccdb87927af4`다.
제품 검증 트리는 `/tmp/knowslink-public-identity-qa-59b66ad`다. 제품 코드, Cloudflare, 배포, PLANS, board는 수정하지 않았다.
coordinator 상태 메시지 `msg_47cb347d1391`은 원본 QA를 끝내고 F1 공유 예산을 관측해 기록하라고 했다. 수락 기준과 소스 checkout은 유지했다.

## 수행 순서

1. 지시서의 고정 후보 절, 신원 구현, sink 절차, 리뷰 보고의 F1 절을 읽었다.
2. 임시 PostgreSQL, loopback sink, loopback relay로 `probe.py`를 실행했다. 종료코드는 0이다. 157통과, 0실패, 2건너뜀이다.
3. QA checkout에서 `npm ci --prefix adapters`를 실행했다. 종료코드는 0이다. 이어서 `make verify-mvp`를 실행했다. 종료코드는 0이다.
4. 프로브가 끝난 뒤 `f1_observe.py`로 F1만 한 번 관측했다. 종료코드는 0이다. 157개 check를 다시 돌리지 않았다.
5. 시나리오, QA 보고, 이 실행 기록, 인박스 완료 보고를 작성했다. 기록 커밋 lint 종료코드는 0이다. `work.py finish`가 인박스를 `handovers/logs/2026-10-05_to_tester.md`에 보존하고 인박스를 비웠다.

## 결과

fixture QA-P01–P05는 통과다. 상세는 [QA 보고서](../../evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md)다.
F1 medium은 재현됐다. source A의 거부 170회 뒤 `http:new` 길이는 200이고 source B의 버킷은 0인데 B의 첫 요청은 429다. 포화 뒤 기존 `/home`은 429이고 logout은 303이다.
일반 서비스 수락은 BLOCKED다. QA-P06과 QA-P07 사람 확인은 미실행이다.

## 정리

프로브 컨테이너와 F1 컨테이너 `knowslink-pi-f1-pg`는 삭제했다. F1 작업 디렉터리 `/tmp/knowslink-pi-f1`은 삭제했다.
QA checkout의 추적 파일 diff는 없다. `127.0.0.1:5432`와 `127.0.0.1:8080`은 변경하지 않았다.

## 남은 일

QA-P06은 운영 공개 후보가 생긴 뒤 현재 hostname에서 확인한다. QA-P07은 designer 화면 검수와 실제 사용자 이메일이 필요하다.
F1은 DEV의 기존 규칙 안 수정과 좁은 회귀 대상이다. 후속 고정 SHA가 오면 변경 영향만 추가 검증한다.
