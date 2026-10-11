---
title: SAR-SSE-INSTALL-001-REVIEW 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-REVIEW]
summary: febd128 고정 SHA의 독립 리뷰와 lite QA 결론
---

# SAR-SSE-INSTALL-001-REVIEW 리뷰

- 검토자 / CLI / 모델: session `01a126ba-6211-7961-b7ec-d72b4ccb6575` / Grok 4.7. 구현자 session은 `01a12678-994d-7680-a15c-83577943b52d`다.
- base SHA / head SHA / merge-base: `4786a00a1274479ea5c3d80decccc0ea556c52e8` / `febd1282df2e84a29fa0a53834569e7b9cc028e2` / 같은 base.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13. `review/rule.json` sha256 `f1061481dcdb66affe2f0f69a6c2354d276babcf23434dedfe6b622171994de8`. 호스트가 판정했다. OCR 자동 판정이 아니다.
- 요구사항·완료 기준 원천: `handovers/to_tester.md`, DEV/INSTALL/REMOTE 계획, D03/D05/D10/D11/D12/D13.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 65 / 51 / 14 / 51 / 14.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0, WARNING 18, 실행 불가 0. product-lint와 product-test는 passed, 종료코드 0, head febd128이다.
- SIZE-002: 추가 3352줄이 상한 400줄을 넘는다. OAuth, SSE, webhook, 설치 패키지가 한 relay_state 계약이라 나누지 않았다. 경고로 남기고 수정 요청으로 올리지 않는다.
- DEP-001: `adapters/package.json`은 scripts만 바뀌었다. dependencies 추가는 없다. SDK 1.32.0과 pgx v5.10.0을 유지한다. OAuth와 HMAC은 Go 표준 라이브러리다.
- UI 디자인: consent는 `html/template`와 기존 `memberStyle`이다. 디자인 lint와 테마 전환은 이 저장소 화면에 해당하지 않는다. 최종 시각 검수는 designer dispatch `ctx_4283f224b494` 소유이며 tester는 실행하지 않았다.

적용 기준은 fullops-common-0.3.3의 README, coding-style, testing, security와 `project.md`, `docs/agents/document-writing.md`, `rules/delegation.md`다. snapshot은 `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-757ce58618f543ad94da80b02d0f320a`다. detached HEAD와 clean 상태를 유지했다. 제품 코드는 고치지 않았다.

## 검토 범위

51개 리뷰 대상은 snapshot diff와 대응 구현을 읽었다. 14개 `*.log`는 skipped다. 짝이 되는 gate JSON과 evidence.json으로 존재와 과거 종료코드만 확인했다. febd128의 새 실행으로 세지 않는다.

`e9ab92d13d5836ed01bbccd4dab4f3ee52a25822`와 head 사이 제품 트리 diff는 비어 있다. 그 구간의 제품 바이트는 같다. evidence.json의 `security_regression_sha`는 e9ab92d다. 최종 gate는 dev git-dir `fullops-gate/pass.json`의 head febd128이다. attempt는 `0b3a8ed545a443e18fff072b126aadd5`다.

## 발견 사항

medium, 미해결, `internal/relay/http.go` 202–213. `GET /v1/text/events`가 보안 헤더 middleware 밖에 있다. 인증 실패 `respond()`와 스트림 응답에 Cache-Control no-store, CSP, Referrer-Policy, nosniff가 없다. D05 35행은 모든 응답이 no-store라고 적는다. 본문은 빈 `ready`/`arrival`과 오류 코드다. 자격 증명과 메시지 원문은 나가지 않는다. critical/high가 아니므로 이 기록의 수락을 막지 않는다.

low, 미해결, `internal/relay/remote_events.go` 197–236. `64:ff9b::/96`은 거절 목록에 없다. `2002::/16`과 `2001::/23`은 거절한다. 호스트가 NAT64 prefix를 내부 IPv4로 라우팅하고 DNS가 그 주소를 주면 pinning이 통과할 수 있다. 콜백 전제는 443과 공인 호스트다.

다른 경계에서는 미해결 critical/high를 찾지 못했다. refresh 재사용은 검증된 installation/client/resource MAC만 설치 Exp를 지금으로 두고 `sweepRemote`한다. `remoteGrantRefusal`은 error가 nil이라 transaction이 commit한다. 위조 token은 철회하지 않는다. webhook은 Proxy 없음, TLS 1.2 이상, redirect 거부, 포트 443, 모든 DNS 답 검사, 첫 공개 IP pinning이다. 구독 비밀과 설치 개인키는 relay_state JSON에 있다. 문서가 DB 접근과 백업 비공개를 경계로 적는다.

## 검증 및 남은 제약

Linux 고유 클론 `/tmp/knowslink-sse001-tester-ctx70e952fbc5c8`에서 head febd128, dirty 0, `make install` 종료코드 0, `make lint` 종료코드 0, `make test` 종료코드 0이다. `go test -race`의 `internal/relay`와 adapter `events.test.js`가 PASS다. integration tag 검사는 make test에 없다. 로컬 Docker는 쓰지 않았다.

`test_public_config.sh` stderr에 `cannot overwrite existing file`과 `invalid tunnel UUID`가 있다. SSH helper가 stdout 뒤에 stderr를 붙인다. make lint 종료코드는 0이다. 그 문구는 라이브 설정 덮어쓰기 거부다.

SEC-001은 `internal/relay/remote_test.go` 210행의 합성 `whsec_`다. 테스트 음성 대조이며 운영 비밀값이 아니다. SIZE-001은 이미 긴 파일과 새 remote 파일의 300행 상한 경고다. 프로토콜 검사를 인라인으로 둔 이유이며 이번 수락을 막지 않는다. SLOP-004는 `package_plugin.py`와 `verify_grok_plugin.py`의 사용자 PASS 출력이다. 디버그 잔재로 수정 요청하지 않는다.

실행하지 않은 항목은 PASS가 아니다. 실제 Google 로그인, 외부 앱 등록, Claude Desktop MCPB 런타임, Dots 계정 Events callback, 운영 도메인 SSE, tester의 동의 화면 시각 검수다. 시각 검수 담당은 designer `ctx_4283f224b494`다.

Windows에서 잘못된 bash 인자가 `rm -rf \`로 풀린 프로세스를 강제 종료했다. 캡처는 `Program Files\Go`의 권한 거부였다. tester 체크아웃과 snapshot은 그 뒤에도 남아 있다. 삭제가 조용히 더 있었는지는 이 로그로 전부 증명하지 못한다. 같은 형태의 로컬 `rm`은 다시 실행하지 않았다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | D03 architecture.md SSE 절, REMOTE.md 구현 계약 | 확인 | 실제 계정 작업은 수락 뒤다. |
| 구조와 구현/미완료 상태 | D05, D10, `text_events.go`, `remote_auth.go`, `remote_events.go` | 확인 | SSE 보안 헤더 medium이 문서의 no-store보다 좁다. |
| 실행·검증 방법과 증거 | DEV.md, evidence.json, dev `pass.json`, tester QA report | 확인 | evidence.json의 과거 SHA는 새 febd128 실행이 아니다. |
| 운영·복구 | D12 SSE 절, D13 rollback 문장 | 확인 | tester는 운영 명령을 재실행하지 않는다. |
| 다음 작업·담당·재개 조건 | D13 후속 순서, to_tester.md 제약 | 확인 | 병합과 운영 적용은 coor다. 시각 검수는 designer다. |
| 로컬 링크·절 접근 | D11에서 D13, D12에서 transition.md | 확인 | 인덱스 밖의 과거 운영 초안은 새 통과가 아니다. |
| snapshot 정리 후 정본 접근 | 리뷰 폴더는 tester 체크아웃에 있다 | 확인 | snapshot cleanup은 하지 않았다. coor가 release 뒤에 정리한다. |

## 검토 결론

수락 가능이다. 미해결 critical/high는 없다. medium 헤더 누락과 low NAT64는 남는다. exact SHA의 make lint와 make test는 종료코드 0이다. 실제 앱 계정과 시각 검수를 이 결론의 통과로 세지 않는다. 병합은 coor가 최신 SHA와 이 기록을 보고 판단한다.
