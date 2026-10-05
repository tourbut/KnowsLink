---
title: SAR-PUBLIC-IDENTITY-001-TESTER — 고정 신원 QA
status: draft
updated: 2026-10-05
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-TESTER]
summary: "후보 59b66ad의 fixture 통과, F1 medium 관측, 공개·사람 확인 BLOCKED를 기록한다"
---

# SAR-PUBLIC-IDENTITY-001-TESTER — 고정 신원 QA

## 판정

판정 후보는 `59b66ada8b36802484cc6d7e22523257b50572cc`다. 제품 코드 SHA는 `a446d89ff288c4243ad6d7f8780a778517154584`다. 기준은 `94533b207b456c0560800fe30a7c90b2b5887c6e`다.
기록 준비 HEAD는 `746ecd9193e9283369267d51e110ccdb87927af4`다. 실행 위치는 `/tmp/knowslink-public-identity-qa-59b66ad`다. 시작과 끝의 추적 트리는 깨끗하다.
fixture QA-P01–P05는 통과다. `probe.py`는 157통과, 0실패, 2건너뜀이다. 건너뜀은 QA-P06과 QA-P07 사람 확인이다.
`npm ci --prefix adapters` 종료코드는 0이다. `make verify-mvp` 종료코드는 0이다. 격리 프로젝트는 `knowslink-mvp-c29927e9ca`다.
F1 medium을 이 Run에서 관측했다. 한 익명 source의 거부 요청이 공유 `http:new` 200/60초를 채운다. 다른 source의 첫 요청과 기존 세션의 `/home`은 429다. logout은 303이다.
일반 서비스 수락은 BLOCKED다. QA-P06 공개 경계와 QA-P07의 designer 검수·실제 사용자 이메일은 미실행이다. 제품 코드, Cloudflare, 배포, PLANS, board는 수정하지 않았다.
리뷰 기록은 `25b110fb694d`의 [SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md](SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md)다. 이 보고의 F1 수치는 그 기록을 복사한 값이 아니다.

## 기준

날짜는 2026-10-05다. 공통 기준은 `fullops-common-0.3.2`다. `fullops-test`와 문서 작성 규칙을 적용했다.
지시서의 QA-P01–P05 필수 실패와 미해결 critical/high는 없다. F1은 medium이다. medium은 그 문장의 identity 차단 조건에 들어가지 않는다.
QA-P06과 QA-P07 사람 확인이 미완료이므로 일반 서비스 수락은 미완료다. 로컬 synthetic 403은 `p04-synthetic-closed`이며 공개 통과가 아니다.

## 환경

| 항목 | 값 |
|---|---|
| QA checkout | `/tmp/knowslink-public-identity-qa-59b66ad`, detached `59b66ada8b36802484cc6d7e22523257b50572cc` |
| 프로브 | `SAR-PUBLIC-IDENTITY-001-TESTER-test/probe.py`, 증거 `result.json` |
| F1 | `f1_observe.py`, 증거 `f1-result.json` |
| F1 Postgres | `postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24`, 컨테이너 `knowslink-pi-f1-pg`, 실행 뒤 삭제 |
| 메일 | loopback `scripts/mail_sink.py`, 디렉터리 0700, 파일 0600 |
| 회귀 | QA checkout에서 `npm ci --prefix adapters` 뒤 `make verify-mvp` |

프로브용 임시 Postgres와 `/tmp/knowslink-pi-qa-run`도 실행 뒤 삭제했다. `127.0.0.1:5432`와 `127.0.0.1:8080`의 기존 listener는 변경하지 않았다.

## fixture QA-P01–P05

`probe.py` 종료코드는 0이다. 기대와 관측이 다른 check는 없다.
QA-P01은 확인 전 회원 부재, 마스킹, 오답 잔여, 재사용 거부, 5회 오답 뒤 만료, 만료 코드, cross-site 403, Origin 불일치 403, 발송 실패 503, mail 미설정 503, 비-loopback 평문 SMTP 거부를 통과했다.
QA-P02는 재로그인 동일 회원, 다른 발급자 비병합, 동시 첫 가입 회원 하나, 동시 확인 상태 303 두 건, 점 별칭과 plus 별칭의 분리를 통과했다.
QA-P03은 재시작 후 세션 유지, 현재·전체 로그아웃, 옛 세션 거부, 5분 밖 전체 로그아웃 403, 재확인 뒤 전체 로그아웃 복구, 유휴·절대 만료, agent의 별도 수명을 통과했다. cookie는 `__Host-` 이름과 Secure, HttpOnly, SameSite=Strict를 가졌다.
QA-P04는 synthetic 403, 세션 없는 agent API 401, 회원 cookie의 `/owner` 401, 다른 gate 조회 403, 다른 gate 결정 409, cross-site 결정 403과 미승인을 통과했다.
QA-P05는 즉시 재발송 429, 61초 뒤 303, 이메일·IP·전역 발송 한도의 재시작 유지, 익명 30과 다른 source 허용, 회원 40, 정리 budget의 logout 303, 전역 `http:new` 200, 회원 100에서 기존 로그인 허용과 신규 503을 통과했다. `X-Forwarded-For`로 기동한 relay 종료코드는 1이다.
`p05-anonymous-30`은 한 source의 31회 근처만 본다. 그 검사는 다른 source를 아직 허용한다. F1의 증거가 아니다.

## F1 medium

위치는 `internal/relay/identity.go`의 `hit`과 `anonymousRate`, `memberRate`다. `http:new` 200/60초가 per-IP 30/60초와 per-member 40/60초보다 먼저 증가한다.
1단계 관측이다. source A의 잘못된 형식 30회는 모두 422다. 이어서 170회는 모두 429다. 그 직후 `http:new` 길이는 200이고 A의 IP 버킷은 31이며 B의 IP 버킷은 0이다. B의 첫 요청은 429다. 그 거절 뒤 `http:new`는 201이고 B의 IP 버킷은 그대로 0이다.
2단계 관측이다. fixture 로그인의 start와 verify는 303이다. 그 뒤 `http:new` 길이는 2다. 다른 source의 30회는 422이고 168회는 429다. `/home` 직전 `http:new` 길이는 200이다. 기존 세션의 `/home`은 429다. 새 source의 첫 start도 429다. logout은 303이고 Location은 `/?n=logout`이다. cleanup 버킷 길이는 1이다. `http:member` 키는 0개다.
relay 로그에 확인 코드는 없다. 제품 동작은 바꾸지 않았다. 임시 컨테이너와 작업 디렉터리는 삭제했다.

## 회귀

`make verify-mvp` 안의 `go test -tags=integration -race -count=1 -v ./internal/relay` 종료코드는 0이다. `TestEmailIdentity`와 PostgresSafety, Gate, Trial이 통과했다.
`node adapters/dist/synthetic.js` 두 번과 `node adapters/dist/trial-check.js`의 종료코드는 0이다. `docker compose down --volumes` 종료코드는 0이다.
이 회귀는 DEV 로그의 재진술이 아니다. QA checkout에서 이 Run이 실행했다.

## 미실행

QA-P06은 미실행이다. 운영 공개 후보와 현재 hostname의 owner-only 경계가 이 Run에 없다.
QA-P07의 HTML 문구 세 검사는 통과했다. designer의 직접 화면 검수와 실제 사용자 이메일 확인은 미실행이다.
운영 SMTP와 Cloudflare edge rate limit은 미실행이다.

## 증거

증거 디렉터리는 [SAR-PUBLIC-IDENTITY-001-TESTER-test/](SAR-PUBLIC-IDENTITY-001-TESTER-test/)다. `result.json`과 `f1-result.json`에는 `@`와 DB URL이 없다.
