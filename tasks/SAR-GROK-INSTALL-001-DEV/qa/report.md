# SAR-GROK-INSTALL-001-DEV QA 보고 (2차, D066 Remote HTTPS)

- 기준: `b3b9cc38` / 브랜치 `fullops/dev` / 테스트 레벨 lite
- 1차 팀 marketplace 안(`98ce0565`, `1cc6b9f5`, `b52f2c17`)은 `3f1c4d40`에서 되돌렸다. 1차 컨테이너 설치 증명 기준은 로컬 설치가 없어져 적용하지 않는다.

## 검증한 것

| 항목 | 방법 | 결과 |
|---|---|---|
| Remote agent가 `self:` agent로 grant를 받음 | `TestRemoteOrcaFormsUseOwnerGrant` | 통과: grant 없음 `approval_required`, task·answer grant 뒤 수락, `Sent` 1건 기록 |
| 형식 검사 | 같은 테스트의 6개 실패 경계(추가 키, text 혼용, body 누락, kind 누락, 비 grant kind, item 0) | 모두 `invalid_schema` |
| 받은 decision의 `kind`·`untrusted` | 같은 테스트 | `kinds[id]=orca.decision.v1`, `untrusted:true` |
| OAuth redirect 규칙 | `TestRemoteRedirectAdmitsHTTPSAndLoopbackOnly` | HTTPS·loopback 허용, 비 loopback http·userinfo·fragment·비 443 거절, webhook은 HTTPS 전용 유지 |
| 웹 안내 ko·en | `TestAutosetupGuideOnlyForClaudeAndCodex`, `TestLandingConfiguredLoginAndPublicContent`, `TestTranslationCoversGoText` | 통과 |
| 전체 | `fullops lint --from b3b9cc38` (product-lint, product-test 212초) | 통과, ERROR 0 / WARNING 3(파일 크기 2, 테스트 시간 1) |

## Grok OAuth 호환 문서 대조

근거와 판정은 [adapters/README.md](../../../adapters/README.md#remote-https-oauth-호환-2026-11-28-문서-대조)에 있다. 공식 문서상 막는 점은 없지만 미확정이다. Grok Bot 문서는 OAuth 세부를 적지 않는다. Cursor MCP 문서를 간접 근거로 써서 loopback redirect 등록을 허용했다.

## 못 한 것

- 실제 Grok Bot 계정 연결(사람 단계), 운영 배포, push.
- MCP 원문 HTTP 왕복은 기존 `remote_integration_test.go`(Docker DB)가 덮으며 새 integration 테스트는 더하지 않았다.
