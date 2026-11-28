# SAR-GROK-INSTALL-001-DEV 리뷰 보고

- base: 88742a66e27fe82dcb052adad9d88a7bcbbc3f64
- head: ba339edd34ddd3db0737e4f1d9f3160e466e49f1
- snapshot: C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-38c622de
- 리뷰 초점: 주인 지정 D067(OAuth loopback redirect, scope 생략, `/mcp` orca 형식 송수신)

## 요약

D066 Remote HTTPS 전환 변경 18파일을 읽었다(`dev.r0.md`는 하네스 기록이라 건너뜀). 주인 지정 보안 초점 세 가지에서 결함은 찾지 못했다. 문서·웹 안내에 low 2건, 테스트 nit 1건이 있다.

### 초점별 판정

1. **loopback redirect(DCR)**: `redirectURL`은 scheme `http`, hostname이 정확히 `localhost`·`127.0.0.1`·`::1`, userinfo·fragment 없음만 받는다(`localhost.evil.test`·대문자는 거절). authorize는 등록 목록과 `redirect_uri` 정확 일치를 보고, token은 code의 client·redirect·resource 일치와 PKCE S256 verifier를 다시 본다. 공격자가 등록한 loopback client로 피해자를 동의시켜도 code는 피해자 기기의 loopback으로 간다. 기존에도 임의 HTTPS redirect 등록이 가능했으므로 공격면이 늘지 않는다. open redirect는 동의 화면이 클라이언트·반환 주소를 보이고 Google 재인증·CSRF를 거친 뒤에만 이동하므로 해당 없음. 등록 남용은 기존 1000 client 상한·만료 그대로다. webhook callback은 `remoteURL`(HTTPS 443 전용)을 유지한다.
2. **scope 생략**: 발급 scope는 유일한 `knowslink.text`뿐이고 metadata·등록·token 응답이 모두 그 값을 낸다. 생략을 그 값으로 보는 것은 RFC 6749 3.3 기본 scope와 같다. 다른 scope 값은 계속 거절한다. 권한 상승 없음.
3. **`knowslink_send` kind·body**: kind와 body는 둘 다 있거나 둘 다 없어야 하고, text와 섞이면 거절한다. body는 JCS 정규화 뒤 `validGrantBody`(고정 키, 2 KB)로 검사하고 설치 키로 서명한다. 실제 허용은 `ingestPublicText → grantFor`가 정한다: `self:` scope, 같은 owner, (from,to,kind) 일치, 살아 있는 grant, 생성 시점 key 스냅샷, 시간당 30건 상한. 다른 owner나 grant 없는 발송은 `approval_required`로 막힌다(테스트로 확인). grant 경로는 `copyConversation`을 건너뛰어 본문을 대화 사본에 남기지 않는다. receipt에는 digest만 남는다. 변경 파일에 본문 로그 호출은 없다. receive는 `untrusted:true`를 유지하고 `kinds`만 더한다.

grant 대상: 웹 사전 승인 폼은 `.Own` 전체를 보이므로 Remote 설치 agent도 고를 수 있다. `createGrant`는 같은 owner와 `self:` 권한을 요구한다.

## 발견 사항

| # | 심각도 | 위치 | 내용 |
|---|---|---|---|
| 1 | low | `internal/relay/human_pages.go:46`, `docs/README.md:20` | Grok Bot (원격) '기존 연결 계속'이 로컬 Command bundle 갱신 절차(`#grok-기존-설치-업데이트`)만 링크한다. Remote Grok 사용자에게 맞는 안내는 Dots와 같은 Refresh다. |
| 2 | low | `README.md:78` | '기본 `Grok Bot`' 이름은 DCR `client_name`에 'grok'이 있을 때만 맞는다. Cursor backend가 등록하면 기본값은 'Dots'일 수 있다. 검증되지 않은 주장이다. |
| 3 | nit | `internal/relay/remote_auth.go:262` | scope 생략 authorize→token 경로와 loopback 포트 불일치 거절에 표적 테스트가 없다. |

## 결론

accept. high 0. 보안 초점(D067) 세 가지 모두 결함 없음. low 2건은 안내 문구 정합 문제라 병합을 막지 않으며, 실제 Grok 계정 연결(SAR-ORCA-BOT-001-OPS) 전에 고치길 권한다.
