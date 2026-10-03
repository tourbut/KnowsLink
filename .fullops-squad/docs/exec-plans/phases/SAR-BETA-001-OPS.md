---
title: SAR-BETA-001-OPS 실행 기록
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-BETA-001-OPS]
summary: 본인 전용 합성 베타 배포 준비의 근거와 검증 증거 및 held 항목을 기록한다
---

# SAR-BETA-001-OPS 실행 기록

## 기준·범위

기준 ref는 `1314e7f`다. 수락 후보는 main `557ebc3`, 제품 `78b1d92`, QA `659f4b0`, 최종 리뷰 `311381f`, 직접 UI `e238777`이다. OPS 체크아웃의 제품 코드(`.fullops-squad` 밖)는 `78b1d92`와 diff가 없음을 `git diff --stat`로 확인했다. Task는 task_0cc034d2aaa9, Dispatch는 ctx_e0f114306c6c, Run은 run_8ca8bc058ab7이다.
이전 OPS의 `a6a10c7` 배포 금지·C1 high는 역사적 기록이다. D12 서두에 새 수락 SHA와 연결해 갱신했다. 전체 공개 한도·신원·실데이터 held는 유지한다.

## 적용한 외부 근거

- Cloudflare 스킬 `cloudflare-one`: 현재 문서·API 스키마를 조회한 뒤 설정한다는 규칙, Access 기본 deny, 새 공개 정책을 먼저 좁게 둔다는 규칙을 적용했다.
- 공식 docs MCP(`search_cloudflare_documentation`)와 Context7 `/cloudflare/cloudflare-docs`: Tunnel 원점 `originRequest.access.required/teamName/audTag`, reusable policy와 email selector, OTP는 이메일 제한과 함께만 쓴다는 경고, `tunnel route dns`는 DNS만 만들고 connector가 꺼져 있으면 노출하지 않는다는 설명, Access 쓰기 권한 이름을 확인했다. 문서의 API 표기는 `Access: Apps and Policies Write`이고 토큰 화면의 표기는 `Edit`이다. 같은 권한이다. 비공개 코드·이메일·비밀값은 조회문에 넣지 않았다.
- Cloudflare API MCP 스키마 검색으로 self-hosted app 생성 필드(`destinations`·`policies`·`allowed_idps` 등)를 확인했다.

## 읽기 전용 재확인 (2026-10-03)

| 항목 | 결과 |
|---|---|
| 계정·zone | Cloudflare API MCP `GET /accounts` 200, zone `knowslog.com` active. MCP 권한 목록에 edit/write 없음 — **읽기 전용** |
| Access | 앱 0개. reusable policy 1개(`knowslog-bot - Production`, 다른 이메일 허용 — 재사용하지 않음). IdP는 One-time PIN 1개. team `scshin88` |
| DNS | `mcp`·`orca`·`s8` proxied CNAME 3개. `link` 없음. `dig A link.knowslog.com` 빈 응답, HTTPS 000 |
| 서버 | 8080 비어 있음. 호스트 cloudflared PID 506937. `myportfolio` 5개 Up |
| 공유 baseline | `orca`·`s8` 200, `mcp` 401. `verify.py baseline`으로 보존했다 |
| 토큰 환경변수 | `CLOUDFLARE*`·`CF_*` 없음 |

## 실행·검증 증거 (모든 종료코드는 직접 확인)

| 검사 | 결과 |
|---|---|
| `lint.py --repo . --from 1314e7f` | 첫 실행 exit 1: `prettier: not found`(`npm ci` 미실행). `npm ci --prefix adapters` 뒤 exit 0, ERROR 0 / WARNING 0 / 실행 불가 0, product-lint passed |
| Compose 병합 검사 | `docker compose … config --format json`: cloudflared `environment` 없음, command·user·마운트 적용, relay `127.0.0.1:8080` 단일 게시 |
| `beta.sh prepare` | exit 0. 체크아웃 detached, 상태 디렉터리 0700, `.env` 0600. 새 비밀값은 `openssl rand`로 생성했고 기록하지 않는다 |
| `beta.sh up` | exit 0. postgres healthy, migrate Exited(0), relay healthy. 구성 SHA 변경 후 재실행에서 migrate·relay가 재생성됐고 migration은 멱등이다 |
| `verify.py local` | exit 0. health 200, `/_probe` 404, `/owner`·`/v1/contacts`·`/v1/receipts/x` 인증 없음 401, relay loopback 전용, Postgres 비게시, restart·메모리 제한, `.env`·tunnel 권한 |
| 합성 e2e(`synthetic.js`) | exit 0: TS verify/persist/ACK/claim, 정책 deny, owner approve, 실효과 없음 |
| owner/CSRF 음성 | 인증 없음 401, agent credential으로 owner 경로 401, 잘못된 CSRF 403 |
| `beta.sh backup` | exit 0. dump 0600, `pg_restore --list` 23행 |
| `beta.sh restore-verify` | exit 0. 네트워크 없는 임시 컨테이너 복원, public 테이블 2개, 컨테이너 정리 확인 |
| Tunnel config | 임시 UUID·AUD로 렌더링한 설정을 `cloudflare/cloudflared:2026.9.1` 이미지에서 `ingress validate` exit 0, `ingress rule`이 `link.knowslog.com`→`relay:8080`로 매칭 |
| `beta.sh tunnel-create` | exit 0. 새 Tunnel `knowslink`. `cert.pem`의 Tunnel 쓰기 권한 확인. 자격 파일 0600을 상태 디렉터리에 두었다. `orca` Tunnel·호스트 cloudflared 불변 |
| `access_apply.py selftest` | 통과(본문 모양·이메일 단독 allow 판정) |
| 공유 서비스 회귀 | `verify.py regression` 여러 번 통과: 코드 200/200/401, PID 506937, `myportfolio` 상태 불변 |

## 읽기 전용 Cloudflare 호출 원시 기록 (비밀 없음)

모두 Cloudflare API MCP `execute`의 GET이며 응답 `success: true`였다.

| 경로 | 상태 | 요지 |
|---|---|---|
| `/accounts` | 200 | 계정 1개 |
| `/accounts/{id}/access/organizations` | 200 | auth_domain `scshin88.cloudflareaccess.com` |
| `/accounts/{id}/access/apps` | 200 | 0개 |
| `/accounts/{id}/access/policies` | 200 | reusable 1개 |
| `/accounts/{id}/access/identity_providers` | 200 | One-time PIN 1개 |
| `/accounts/{id}/cfd_tunnel?is_deleted=false` | 200 | `orca` 1개(healthy, 연결 4) |
| `/zones?name=knowslog.com` | 200 | active, 권한 목록에 edit/write 없음 |
| `/zones/{id}/dns_records?per_page=100` | 200 | proxied CNAME 3개(mcp·orca·s8) |

## 독립 리뷰 대응 (dev 리뷰 f625c4e, 대상 437f143)

미해결 critical/high는 없었다. medium·low는 다음처럼 처리했다.

| ID | 처리 |
|---|---|
| M1 | D12 서두·6장·9장·10장을 역사적 기록으로 표시하고 현재 상태를 11장·D13으로 연결했다 |
| M2 | 최소 권한에 IdP 읽기를 추가했다(D12·D13·본 문서) |
| M3 | `access_apply.py check`를 추가하고 `beta.sh expose`가 호출한다. live 앱·정책·`aud`·`teamName`·Tunnel `audTag`를 비교한다. `dig` 실패는 중단한다. `selftest`가 깨진 앱·설정을 거부하는지 검사한다. 임시 상태 디렉터리와 더미 이메일로 `check`의 통과·거부(정책 2개)를 직접 실행해 확인했다 |
| M5 | `beta.sh deploy <sha>`(백업·이력·마이그레이션 차이 시 중단·재빌드·검증)와 rollback 절차를 D12 11.3에 추가했다 |
| L1 | 사후 검증이 destinations·IdP 수·우회 옵션·`aud`를 본다. 앱 목록은 `per_page=100`이다(상한: 100개 초과 계정은 미지원). 토큰이 계정 하나만 볼 때만 진행한다 |
| L2 | 8080 점유 검사가 `0.0.0.0`·`*`·`[::]`를 포함한다 |
| L3 | backup 실패 시 부분 dump를 삭제하고, backup·restore-verify가 `relay_state` 행 수를 출력해 비교한다 |
| L4 | D12 11.5와 D11이 owner-login을 사용자 본인만 실행하도록 명시했다 |
| L5 | 위 원시 기록 표를 추가했다 |
| L6 | D11 5단계를 칸별로 분리했다 |

재검증(수정 SHA에서): `beta.sh deploy 437f143`(rollback, exit 0)·`beta.sh deploy <수정 SHA>`(no-op 이동, exit 0)·`verify.py local` 통과·`restore-verify` `tables=2 relay_state_rows=1` exit 0·`verify.py regression` 불변. rollback 대상에 `deploy` 명령이 없어 앞으로 이동은 수동 체크아웃이 필요했다(D12 11.3에 한계로 기록).

## 재리뷰 대응 (SAR-BETA-001-REVIEW-FINAL, 대상 f824015)

critical/high는 없었다. 신규 N1~N5를 처리했다.

| ID | 처리와 증거 |
|---|---|
| N1 | `access_apply.py`의 모든 게이트 검사를 `assert`에서 명시적 `need()`(SystemExit)로 바꿨다. `PYTHONOPTIMIZE=1`에서도 `selftest` exit 0이고 aud 불일치 `check`가 exit 1이다(임시 상태 디렉터리·더미 값으로 직접 실행). 정상 `check`는 exit 0이다 |
| N2 | 스냅숏의 app·policy id가 `access.json`과 다르면 거부한다. 미래 mtime(음수 나이)도 거부한다. 직접 실행: 다른 app id exit 1, 미래 mtime exit 1, 11분 경과 exit 1 |
| N3·N4 | D12 11.3과 D13에 `rev-parse` 정본, 수동 이동 미기록, `check` 없는 SHA에서 `expose` 금지를 명시했다 |
| N5 | `selftest`가 `teamName` 값만 바꾼 설정을 거부하는지 검사한다 |

## 아직 실행하지 않은 것 (held)

- Access 앱·reusable policy 생성: MCP와 `cert.pem`에 Access 쓰기 권한이 없다. 필요한 최소 권한을 coordinator에 ask했다. 최소 권한은 `Access: Apps and Policies Edit`와 `Access: Organizations, Identity Providers, and Groups Read` 두 개다(리뷰 M2로 정정).
- DNS `link.knowslog.com`·connector 기동·공개 negative 검사·사용자 이메일 로그인(인간 검사): Access 보호 확인 전에는 하지 않는다.
- 원점 JWT 거부의 동작 검증: Access가 앞단에서 먼저 차단하므로 공개 경로로는 원점 거부를 독립 관찰할 수 없다. 설정 존재와 `ingress validate`만 확인했고, 동작 확인은 인증된 요청 성공과 cloudflared 로그로만 가능하다. 이 한계를 숨기지 않는다.

## 산출물

[D11](../../operations/user-guide.md), [D12](../../operations/ops-guide.md) 11장, [D13](../../operations/transition.md), 구성 `deploy/knowslink/`.
