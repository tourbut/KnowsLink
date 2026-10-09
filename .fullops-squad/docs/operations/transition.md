---
id: D13
title: 인수인계서
status: draft
updated: 2026-10-10
owner: ops
tasks: [SAR-BETA-001-OPS, SAR-MVP-002-DEV, SAR-MVP-002-BOT-CATALOG-DEV, SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-003-BIDIRECTIONAL-OPS, SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW, SAR-PUBLIC-SERVICE-OPEN-PREP, SAR-GOOGLE-CONNECT-001-DEV, SAR-GOOGLE-CONNECT-002-DEV]
summary: 베타 배포의 현재 상태와 인수 항목 및 남은 일을 기록한다
---

# KnowsLink 본인 전용 합성 베타 인수인계서 (D13)

이 문서는 베타 배포의 인수인계 상태를 기록한다. 배포 구성과 절차 정본은 [D12](ops-guide.md) 11장이다. 사용자 시험 절차는 [D11](user-guide.md)이다.

## 현재 상태 (2026-10-03 보호 연결 적용 후)

| 항목 | 상태 |
|---|---|
| 배포 SHA | 구성 `28bd1bb`(독립 리뷰 `7ba9df0`·QA 통과). 체크아웃 `/home/shin/deploy/knowslink`는 detached다. 제품 코드는 수락 제품 `78b1d92`와 동일하다. 현재 SHA의 정본은 `git -C /home/shin/deploy/knowslink rev-parse HEAD`다. `knowslink-state/deploy-history.log`는 `beta.sh deploy` 이동만 기록한다. 갱신·rollback은 D12 11.3 |
| 로컬 스택 | 기동·검증 완료. relay `127.0.0.1:8080`, Postgres 비게시 |
| Tunnel | `knowslink` Tunnel 연결 4개(icn05/06/07). 컨테이너 `knowslink-cloudflared-1` |
| Access 앱·정책 | 생성됨. self-hosted 앱 `KnowsLink beta (owner-only)`, 대상 `link.knowslog.com` 하나, reusable 정책 `knowslink-beta-owner-only`(사용자 이메일 한 개만 allow), IdP는 One-time PIN 하나, 세션 24h |
| 원점 JWT | Tunnel `originRequest.access.required: true`, teamName `scshin88`, audTag는 앱 aud와 일치 |
| DNS `link.knowslog.com` | proxied CNAME → `knowslink` Tunnel. 기존 레코드 3개 불변 |
| 미인증 공개 검사 | 모든 경로가 302로 `scshin88.cloudflareaccess.com`에 이동한다. 실제 응답 본문 없음 |
| 사용자 이메일 로그인 | **인간 검사. 미실행**. D11 절차로 사용자가 수행한다 |
| 공개 수락 | held — 로그인 확인 전, 그리고 전체 공개 한도·실제 신원·실데이터는 별도 held |

## 인수 항목

- 상태 디렉터리 `/home/shin/deploy/knowslink-state`(0700): `.env`, Tunnel 자격 파일, `access.aud`, 백업. 모두 Git 미추적이다. 자격 파일을 공유 위치로 옮기지 않는다.
- 공유 서비스 baseline은 `shared-baseline.json`에 있다. `verify.py regression`으로 비교한다.
- 첫 DB 백업과 격리 복원 검증은 완료했다.
- 관리 연결: 이번 Access 쓰기는 coordinator가 준비한 Codex file-store OAuth(공식 Cloudflare MCP, scopes에 Access 앱·정책 쓰기 포함)로 수행했다. 원본 credential은 복사·기록하지 않았다. 반복 시 같은 방식이 필요하다.
- 보존 대상: `myportfolio` 프로젝트·볼륨, `orca` Tunnel, 호스트 cloudflared PID 506937, `~/.cloudflared/config.yml`.

## 남은 일

1. 사용자가 D11 절차로 본인 이메일 OTP 로그인과 합성 시험을 수행한다(인간 검사).
2. 종료는 D11의 `unexpose`(외부 접근만) 또는 `stop`(전체)이다. 전체 철회는 D12 11.3이다.
3. held 유지: DEC-03 공개 한도, 실제 신원, 실데이터, 실제 벤더, 무제한 공개.
4. 구성 갱신이 필요하면 `beta.sh deploy <sha>`와 새 독립 리뷰를 거친다. 코드 변경 SHA에서 `expose`를 다시 하기 전 `access_apply.py check`를 통과해야 한다.

## SAR-MVP-002 DEV 플러그인 준비 인계

기존 베타·Access·Tunnel·Tailscale·배포 SHA는 변경하지 않았다. 사용자 확정 대상은 xAI 공식 Grok Bot이다. 설치된 Grok Build CLI와 inference API를 대체 대상으로 선택하지 않았다.

`make plugin`으로 `build/knowslink-grok-bot-plugin.zip`을 생성한다. package는 MCP·skill·Cursor manifest와 standalone bundle을 포함하며 기본 held다. [사용자 설치 문서](../../../adapters/README.md)에 marketplace 배포 선행 조건·도구 검색·Node·계정/secret·합성 검사와 실패 의미를 기록했다.

owner/admin은 승인된 marketplace 등록과 실제 Bot 앱의 설치·도구 검색을 담당한다. DEV/OPS는 승인된 relay network·최소 권한·secret 전달을 확인한다. coor는 fixed-SHA 독립 리뷰·TESTER QA와 PLANS/board를 갱신한다. 계정·운영 연결은 이번 과제에서 실행하지 않았다. 문서 완료와 로컬 MCP 성공을 실제 Bot 연결 수락으로 표시하지 않는다. 자세한 지원 근거·검증 SHA·재개 조건은 [실행 기록](../exec-plans/phases/SAR-MVP-002-DEV.md)을 따른다.

## SAR-MVP-002-BOT-CATALOG-DEV 앱 카탈로그 인계

재시험에서 CLI 설치는 성공했지만 앱 카탈로그에 knowslink가 없었다. 앱 계정 등록 단계의 부재는 미확정 가설이다. `scripts/install_bot_mcp.sh`와 README의 Command server 등록 절차를 추가했다. 리뷰 보완([FIX 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md))에서 원인 단정과 `Settings → Plugins` 경로를 정정했다. Command 등록 근거는 Team Bots 문서에 한정된다. 개인 계정 UI·승인 카드·`AddMcpServer` 제공은 미확인이다.

owner는 실제 계정에서 등록과 `knowslink_status` held 호출을 재시험한다. coor는 고정 SHA 독립 리뷰·필요한 QA 뒤 FIX 기록의 [재시험 댓글 초안](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md#재시험-댓글-초안-coor-게시)을 게시한다. 이전 실행 기록의 초안은 게시하지 않는다. 실제 앱 노출은 미검증이다. Remote HTTPS·Marketplace 발행·실제 relay는 별도 승인 전까지 held다.

## 시험 운영 인수 (2026-10-04, SAR-MVP-003-BIDIRECTIONAL-OPS)

현재 배포는 `0911c2c73468f8684260a277d4940a74d26bcf7d`이고 시험 allowlist는 `trial_codex,trial_grok`이다. 공개 `/v1/test/*`는 아직 root owner Access 아래에 있어 외부에서 열리지 않는다. 24h service token·trial 앱 생성은 `access-service-token.write` 권한 부재로 막혀 있다. 사용자가 권한이나 0600 API token 파일을 제공하면 OPS가 재개한다. 절차·rollback·종료는 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md)과 D12 13.5를 따른다. Grok 전용 private 파일은 `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/grok-export/`에 있으며 전달 수단은 미확인이다.

### 재개 결과 (2026-10-04)

위 차단은 해소됐다. 24h machine token 두 개·trial policy·`/v1/test/*` 앱이 적용됐고 원점 Tunnel에 trial AUD rule이 들어갔다. 공개 negative(403·401·302)와 positive(두 local client의 HTTPS 왕복)가 통과했다. actual Grok은 미검증이다. 만료는 2026-10-05T06:42:09Z다. 그 시각 전에 시험을 끝내거나 token을 새로 만든다. 종료 시 token 두 개 revoke·앱과 policy 삭제·원점 config 복원·allowlist 비우기를 한다. Grok 전용 파일은 `grok-export/`에 CF 항목까지 준비됐고 외부 전달은 coor·사용자 안전 채널이 담당한다. 상세는 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md)의 「재개 결과」다.

### 만료 갱신 (2026-10-05, SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW)

service token 두 개가 만료돼 24h로 갱신했다. UUID·secret은 그대로다. 새 만료는 **2026-10-06T08:09:43Z(codex)·08:09:50Z(grok)**이다. 공개 negative·positive와 두 local client 왕복이 다시 통과했다(실제 Grok 아님). Grok 전달용 묶음은 `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/knowslink-grok-trial-20261005.tar.gz`이며 외부 전달은 coor·사용자가 한다. 관리 token은 2026-10-05T23:59:59Z에 먼저 만료되므로 종료 정리를 그 전에 하거나 새 권한을 받는다. 상세는 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md)이다.

## 2026-10-07 최종 코드 수락과 공개 준비

사용자 요청으로 최종 main 수락과 공개 준비를 진행한다. [공개 준비 실행 기록](../exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md)의 고정 SHA·실제 백업/격리 복원·현재 보호상태·배포/복귀 절차를 따른다. 기존 배포0911c2c는유지됐고새서비스/실메일은아직미검증이다. 과거시험의미완료/완료는원래시점으로보존한다. 이 연결은 운영 공개 수락이 아니다.

## Google 로그인 우선 출시 — 2026-10-07

최신 사용자 결정은 SMTP 대신 Google 로그인으로 핵심 기능을 먼저 출시하는 것이다. 과거 이메일 전용 공개 준비·실메일 미검증은 원래 기록으로 유지한다. 새 프로젝트 knowslink-auth의 Web OAuth client를 등록했고 지정 계정은 testing 사용자다. 결제 설정을 활성화하지 않았다.

운영 env 정본은 `/home/shin/deploy/knowslink-state/.env`0600이다. `KNOWSLINK_GOOGLE_CLIENT_ID`, `KNOWSLINK_GOOGLE_CLIENT_SECRET`, `KNOWSLINK_GOOGLE_REDIRECT_URL`을 Compose가 relay에 전달한다. 값은 Git/대화/로그에 남기지 않는다. callback은 `https://link.knowslog.com/auth/google/callback`이다. Google discovery/token/JWK HTTPS egress가 필요하다. 기존 Tunnel·Access와 Postgres volume은 유지한다.

제품 후보017bf456의 독립 인증 delta 리뷰와 기존 lint/test를 확인했다. main 통합 뒤 `beta.sh deploy <수락 SHA>`로 사전 DB backup·migration diff 거부·기동 확인을 수행한다. 이번 변경은 DB migration이 없다. 이전 제품 배포d08903a55c3638128827010400e66e9d45b61d7c로 코드 복귀할 때는 새 Google env 세 개를 비공개 환경에서 비활성화한다. 실제 사용자 로그인은 배포 뒤 `/`의 Google로 계속 → Google 화면의 선택/동의 → 자기 홈 링크로 확인한다. 일반 공개·Grok Bot 실연결 완료는 별도 기록한다.

## Google 연결 DEV 인계 — SAR-GOOGLE-CONNECT-001-DEV

Google 회원 인증과 클라이언트 로컬 키를 연결하고 명시적 동의 뒤 자동 저장하는 구현을 전달한다. 같은 Google 회원의 두 클라이언트는 별도 agent·키·credential을 사용하고 자동 pairing은 없다. D03/D05/D06/D09/D10/D11/D12와 실행 보고를 함께 검토한다. D07/D08 SQL 변경은 없다.

부모는 worker_done의 고정 SHA로 독립 리뷰와 좁은 QA를 수행하고 main 통합 및 D12 운영 변경을 담당한다. 로컬 합성 OAuth/Postgres와 화면 검증은 실제 Google·외부 Bot 성공 증거가 아니다. 기존 held/fail 제품 판정은 유지한다. 운영 자격·개인키·기존 WIP는 보존했다. 실제 두 클라이언트 가입·관계 수락·왕복 전달 및 관리자 보호/공유 호스트 회귀를 운영 수락 조건으로 남긴다.
## 비용 없는 Google 연결 인계 — SAR-GOOGLE-CONNECT-002-DEV

DEV는 별도 public config 후보 생성, 설치 안내, Google 연결의 parent 폴더 생성 및 회귀를 보완한다. D03 D10 D11 D12를 함께 갱신한다. 기존 DB·Google 자격·키·trial 코드·기본 held·비명시 발송 금지는 보존한다. 운영 config·Access·DNS·Bot UI는 DEV가 변경하지 않는다.

coor는 고정 결과 SHA 독립 리뷰·좁은 QA 뒤 [D12의 새 적용 절차](ops-guide.md#비용-없는-tunnel-적용--sar-google-connect-002-dev)로 allowlist+404를 먼저 적용하고 백업한 두 KnowsLink 앱만 제거한다. 실제 Google·독립 두 클라이언트·관계 수락·명시적 양방향 text와 공유 호스트 회귀를 확인한다. 익명 신규 Device 포화의 medium 한계는 남는다. 복구 시 먼저 Tunnel을 닫고 config/앱/aud를 재검증한다. 기존 DB 백업을 덮지 않는다.
