---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER 독립 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER]
summary: DEV 최종 SHA eac603b의 OAuth 승인 callback 복귀 수정에 대한 고정 SHA 독립 리뷰. 미해결 critical/high 없음, medium 1·low 1 비차단
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER 리뷰

- 검토자 / CLI / 모델: tester / Claude Code / claude-sonnet-5-5. 실제 세션 `d7e2c059-f0af-44ff-82ce-6e0f2f6ea834`. 구현자 세션 `e20cc179-625e-43fa-84fd-b431d432b9c6`와 다르다.
- base SHA / head SHA / merge-base: `c399a5e52c714d4ee61a7cbca03d32ffbfed9bcf` / `eac603b6a41d2ce002e698600ac66b8fc67ad762` / `c399a5e52c714d4ee61a7cbca03d32ffbfed9bcf` (base가 head의 조상). 제품 리뷰 SHA는 DEV 최종 SHA다. coor의 준비 merge HEAD(`fb0a60d`)를 제품 SHA로 쓰지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13 (delegate 모드: `review.py prepare`가 대상·규칙을 준비했고 검토는 호스트 에이전트가 수행했다). `rules.json`은 prepare가 만든 그대로다. 규칙 묶음 `fullops-common-0.3.3`(README·coding-style·security·testing)과 `project.md`를 직접 읽었다. 예외 없음.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_tester.md` 완료 기준 1~6, DEV 실행 기록 `docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md`, 등록 gate `dev-eac603b6-final-gate.json`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 13 / 11 / 2(`browser-runs.log` unsupported_ext, `consent_browser.mjs` default_path) / 11 / 2(DEV 패킷 처리 JSON, DEV 지시서 원문 로그).
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: head `eac603b6a41…`에서 ERROR 0 / WARNING 2(SIZE-001 `remote_auth.go`, SLOP-004 `consent_browser.mjs:13`) / 실행 불가 0. product-lint·product-test passed.
- SIZE-002: 해당 없음(발생하지 않음). 실제 추가 줄 수 292(제품 코드 2 + 테스트 약 200 + 드라이버 80). 분할하지 않은 이유: 한 handler 수정과 그 실제 브라우저 검증이 한 단위다.
- DEP-001: 새 의존성 없음. `package.json`·`go.mod` 변경 없음. 드라이버는 Node 22 내장 WebSocket과 기존 `golang.org/x/oauth2`를 쓴다.
- UI 디자인: 새 CSS·토큰 없음. 기존 `memberStyle`과 `refusal` 화면을 재사용한다. 프로젝트 DESIGN lint가 Go 문자열 CSS를 검사하지 않아 해당 없음이다(`project.md`). 테마 전환은 테마가 없어 해당 없음. 반환·오류 화면의 desktop/mobile 캡처를 `qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/screenshots/`에 남겼고 직접 시각 판정은 designer에게 넘긴다.

## 검토 범위

`result.json`의 13개 `(path,status)`에 reviewed/skipped와 사유를 기록했다. 제품 코드 4개(`remote_auth.go`, 두 integration test, 드라이버)와 설계 문서 3개·컨텍스트·운영 가이드·DEV 실행 기록·브라우저 로그는 직접 읽었다. DEV 패킷 처리 JSON과 지시서 원문 로그는 제품 동작과 무관해 skipped이며 사유와 영향을 적었다. 삭제된 `to_dev.md`는 logs 보존본과 대조했다.

내부 호출자 검색: `remote/consent`·`authorization_failed`·`form-action` 사용처를 snapshot 전체에서 검색했다. consent 응답을 303으로 가정하는 다른 코드는 없다(`remote_integration_test.go`만 갱신됨). `public-ingress.yml`은 `/remote/consent/<43자>`를 GET/POST 구분 없이 허용하며 변경이 필요 없다.

## 발견 사항

| ID | 심각도 | 파일·줄 | 재현 조건 | 영향 | 수정 상태 | 검증 근거 |
|---|---|---|---|---|---|---|
| F1 | medium | `internal/relay/remote_auth.go:366` | 소유 agent가 한도(10건 기록)인 owner가 승인(409), CSRF 값 불일치(403), 로그인 만료·계정 불일치(401) | 원인과 다른 안내("만료되었거나 이미 처리됨; 앱에서 다시 시작")가 표시된다. 한도 초과는 재시작으로 해결되지 않고 `/home`에서 정리가 필요하다. 안전 속성(비밀 비노출, callback 미도달, install 0)은 유지된다 | 미해결(비차단). 후속 DEV 권고 | `qa-run-final.log`의 capacity-desktop·csrf-tamper-mobile·owner-binding-mobile 시나리오와 `screenshots/*-refusal.png` |
| F2 | low | `internal/relay/remote_auth.go:314` | 모든 승인·거절 | 링크 높이 24px, 카드 없는 기본 본문. 자동 이동이 실패한 때의 터치 대상이 작다 | 미해결. designer 확인 | `screenshots/approve-link-*-return.png`, 실제 터치·마우스 클릭으로 callback 도착 |

미해결 critical/high는 없다. 요구의 핵심(승인 POST 뒤 Chromium이 등록 callback에 도착하고 PKCE token 교환이 성공)은 독립 실행으로 확인했다. 기준 handler는 같은 시나리오에서 CSP form-action 위반으로 실패한다.

위험이 큰 변경(인증 흐름)이라 다음을 별도로 확인했다. (1) callback 주소는 attempt에 저장된 등록 redirect에 state·iss·code/error만 추가한다. (2) html/template이 `<meta content>`와 `<a href>`를 escape한다. (3) CSRF·Google 재확인(5분)·owner binding·attempt/code TTL·PKCE·resource·client 바인딩 코드는 diff에 없다. (4) 공통 CSP 문자열은 바뀌지 않았고 테스트가 응답 헤더 전체를 비교한다. (5) code가 HTML 본문에 들어가는 것은 이전 Location header와 같은 노출이며 no-store·no-referrer가 유지된다.

## 검증 및 남은 제약

- 독립 실행(대상 SHA, 기존 Linux의 고유 clone·고유 Postgres 17, Chrome 148): DEV 테스트 재실행 PASS, 기준 handler RED(CSP 위반·callback 미도착), 수정 SHA GREEN, tester 시나리오 8개 그룹 PASS(승인 자동·링크 desktop/mobile·거절·만료·계정 binding·로그인 없음·CSRF·한도, token 경계 7종·만료 code·재사용, authorize 경계 7종). 자세한 명령·결과는 `qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md`다.
- 등록 kind: test 명령 `product-test`(`make test`)는 head `eac603b6a41d…`에서 종료코드 0이다. DEV 게이트 영수증(같은 head, 종료코드 0)을 읽고 대조했다. 이 SHA에서 `lint.py`를 새로 실행해 `lint.json`의 head를 `review.py check`의 `--to`와 맞췄다.
- `review.py check`: `--key SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER --from c399a5e… --to eac603b… --task-key SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV` → "기록 검사 통과: reviewed=11, skipped=2, total=13, lint WARNING=2". tester 과제 패킷은 리뷰 대상 SHA(eac603b)보다 뒤의 coor 커밋에 있어 `--task-key`에 eac603b에 존재하는 DEV 과제 키를 썼다. 이 통과는 기록 검사이며 테스트 성공을 증명하지 않는다.
- 미실행: 실제 Dots·ChatGPT 호스트 callback과 token 교환, 실제 Google 계정, 운영 DB·도메인. Chromium 외 브라우저. 실시간 만료 대기(DB 만료 시각 조작으로 대체).
- 기존 medium(SSE 라우트 보안 헤더 누락)·`verify.py` baseline 결함은 이번 변경 밖이다. 새 결함으로 세지 않았다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `handovers/to_tester.md`, DEV 실행 기록 "원인" 절, `contexts/dev.md` 최근 결정 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | DEV 실행 기록 "수정"·"한계와 후속", `design-docs/module-design.md` Remote 절 | 확인 | 미완료: 실제 Dots 검증(부모·OPS) |
| 실행·검증 방법과 증거 | `qa-reports/…-TESTER-test/report.md` "증거와 재현", `qa/` 소스 | 확인 | 재현에 서버 접속 helper가 필요하다. 비밀은 파일에 없다 |
| 운영·복구(해당 시) | `operations/ops-guide.md` 42~44줄 | 확인 | 운영 변경은 이 과제에서 하지 않았다. 이미 생성된 운영 install 1건 정리는 OPS·부모 판단 |
| 다음 작업·담당·재개 조건 | 완료 보고(worker_done), `PLANS.md` 1466~1477줄 | 확인 | main 수락·push·OPS 배포·실계정 연결은 부모. F1은 후속 DEV 권고 |
| 로컬 링크·절 접근/지원 한계 | report.md 상대 링크 | 확인 | snapshot 경로는 정리 전까지만 유효하다 |
| snapshot 정리 후 정본 접근 | 리뷰 폴더·`-test/` 폴더 | 해당 없음 | 이 과제는 snapshot을 정리하지 않고 보존했다. 정리는 소유자·부모 결정 |

## 검토 결론

수락 가능(조건부 비차단 발견 2건). 요구를 충족한다: 실제 Chromium에서 승인 후 callback 도착과 PKCE token 교환이 성공하고, 기준 handler는 실패한다. 허용 callback·PKCE·resource·issuer·CSRF·계정 binding·만료/사용된 consent의 실패 경계는 안전하게 거부된다. 전역 CSP 확장·CSRF 완화·로그인 생략·TTL 연장은 없다. 미해결 critical/high는 0건이다. 실제 Dots·실계정 검증은 성공으로 간주하지 않는다. check 통과는 AI 검토의 내용·테스트 성공을 자동 보증하지 않는다.
