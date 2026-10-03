---
id: D13
title: 인수인계서
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-BETA-001-OPS, SAR-MVP-002-DEV, SAR-MVP-002-BOT-CATALOG-DEV]
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

재시험에서 CLI 설치는 성공했지만 앱 카탈로그에 knowslink가 없었다. 원인은 앱 계정 등록 단계의 부재다. CLI plugin은 앱에 등록되지 않는다. `scripts/install_bot_mcp.sh`와 README의 Command server 등록 절차를 추가했다.

owner는 실제 계정에서 등록과 `knowslink_status` held 호출을 재시험한다. coor는 고정 SHA 독립 리뷰·TESTER QA 뒤 [재시험 댓글 초안](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md#재시험-댓글-초안-coor-게시)을 게시한다. 실제 앱 노출은 미검증이다. Remote HTTPS·Marketplace 발행·실제 relay는 별도 승인 전까지 held다.
