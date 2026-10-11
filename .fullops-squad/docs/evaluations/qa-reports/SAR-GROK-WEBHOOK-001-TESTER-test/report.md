---
title: SAR-GROK-WEBHOOK-001-TESTER 경계 QA 증거
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-GROK-WEBHOOK-001-TESTER]
summary: b86f397 webhook wake의 독립 lite 경계 프로브 실행 명령과 결과
---

# SAR-GROK-WEBHOOK-001-TESTER 경계 QA

- 대상 SHA `b86f39711c2f5240d37332f44a9930d4dc2653fa`. 실행 위치: Linux `/tmp/knowslink-webhook-tester-b86f397-s254e09d8/repo`(clean, detached), Node v22.22.2.
- 명령: `node probe.mjs <repo>/adapters/dist` (`webhook-boundary-probe.mjs`). `make test`가 만든 `dist/`를 사용한다. 실제 `fetch`와 loopback HTTP 서버만 쓰며 외부로 나가는 요청은 DNS `.invalid` 조회뿐이다.
- 결과: 종료코드 0, 58 PASS / 0 FAIL, 정보 2건. 전문은 `probe-and-lint.log`.
- 범위: 2xx·3xx(미추적)·4xx·5xx·timeout·reset·거절·NXDOMAIN 분류, 헤더·body 키, ID 필터·상한, 불신 config 16종과 파일 권한·symlink·크기, key 회전, env 배타·legacy, AutoReceiver marker(중복·uncertain·retry 해제), 비밀 비노출, `wake-check` CLI.
- 정보: `[2001:db8::1]`와 `192.0.2.1`에서 네트워크 도달 불가가 `uncertain/webhook_no_response`로 분류됨(리뷰 F1, low).
- 비밀값: 프로브 key는 합성 값이다. 실제 key·URL·메시지 본문은 쓰지 않았다.
- 미실행: 실제 Bot routine·api2.cursor.sh 호출, Windows ACL, 부하·장시간(lite 범위 밖).
- `lint.py` 출력(종료코드 0, ERROR 0, WARNING 1, 실행 불가 0)도 `probe-and-lint.log` 끝에 있다.
