---
title: SAR-BETA-001 베타 배포 설정 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-BETA-001-REVIEW]
summary: 고정 SHA 437f143의 본인 전용 합성 베타 배포 설정과 D11/D12/D13 독립 리뷰 결과와 공개 held 결론을 기록한다
---

# SAR-BETA-001 베타 배포 설정 리뷰

- 검토자 / CLI / 모델: Orca dispatch `ctx_fa48cdc4605a`(task `task_470cc6941994`, run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 ID는 `bac29ce7-11bc-46ce-8331-c354bf56ade2`이다.
- 구현 세션: OPS `2191cc9b-76ef-4522-9fad-d2c9f017bfbc`. `~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-ops/`의 세션 기록에서 cwd(`fullops-ops`), 브랜치 `fullops/ops`, task `task_0cc034d2aaa9`·dispatch `ctx_e0f114306c6c`, 커밋 `66d1057`~`437f143` 작업을 직접 확인했다. 리뷰 세션과 다르다.
- base SHA / head SHA / merge-base: `1314e7fcd02d0c4b13dd86ba2be72fa0f153129c` / `437f1432a158670a485413c1aba159debc3759e5` / `1314e7f…`
- snapshot: `/tmp/knowslink-beta-review-437f143`. detached·clean·HEAD `437f143`을 확인했다. 검사 전후 `git status --porcelain` 0줄이다. Python은 `-B`·`PYTHONDONTWRITEBYTECODE=1`로 실행했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate / `.fullops-squad/review/rule.json`(sha256 `ab2116fb…`)
- 공통 규칙: `fullops-common-0.3.2`(`.fullops-squad/rules/common/README.md`와 코딩·테스트·보안 규칙), `project.md`, 문서 작성 규칙, ponytail full. 예외 없음.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_dev.md`, OPS 실행 기록 `docs/exec-plans/phases/SAR-BETA-001-OPS.md`, D11/D12/D13.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 11 / 10 / 1(`tunnel/config.yml.tmpl`, unsupported_ext) / 11 / 0. coverage 100%.
- lint(`lint.json`): ERROR 0 / WARNING 0 / 실행 불가 0, product-lint passed. snapshot을 바꾸지 않으려고 `437f143`의 별도 scratch clone에서 `npm ci --prefix adapters` 뒤 실행했다. clone 작업 트리는 깨끗했다.

## 검토 범위

result.json의 11개 `(path,status)`를 모두 reviewed로 기록했다. OCR이 제외한 `tunnel/config.yml.tmpl`은 직접 검토했다.
제품 코드는 기존 수락(`557ebc3`/`78b1d92`)과 같으므로 제품 전체 QA를 반복하지 않았다. 제품 근거는 기존 QA `659f4b0`·최종 리뷰 `311381f`를 그대로 연결한다. 이번에는 `adapters/src/synthetic.ts`의 fixture 권한(0600)·seed 출력·gate 만료(180초)만 대조했다.

검토 항목과 판정은 다음과 같다.

| 항목 | 판정 |
|---|---|
| Access 정책 내용 | email selector 하나만 include, exclude·require 비어 있음, decision allow. 기존 reusable `knowslog-bot - Production`은 재사용하지 않는다 |
| 정확한 hostname | 앱 `destinations=[{type: public, uri: link.knowslog.com}]`, Tunnel ingress `link.knowslog.com`→`relay:8080`, 나머지 404. 경로 제한 없음(전체 경로 보호) |
| 사용자 한 명·reusable | 새 reusable policy 하나만 앱에 연결, 사후 GET으로 정책 ID 하나만 확인. 이메일은 0600 파일에서만 읽고 출력하지 않는다 |
| 바이패스/다른 Allow | 앱 본문에 bypass·service auth·`overrides`·`options_preflight_bypass` 없음. 사후 검증 범위가 좁다(L1) |
| team/aud | 렌더링은 같은 `access.aud`를 쓴다. expose 게이트가 일치를 재확인하지 않는다(M3) |
| origin JWT fail-closed | `originRequest.access.required: true`·`teamName`·`audTag`가 공식 Origin parameters 문서와 일치한다. Access 앱이 사라지거나 aud가 달라도 cloudflared가 요청을 거부한다 |
| 비밀 취급 | `.env` 0600·상태 디렉터리 0700·자격 파일 0600·dump 0600·qa-fixture 0600 확인. access_apply는 토큰·이메일·응답 본문을 출력하지 않는다. Basic 값 출력은 L4 |
| 신규 환경값 | prepare는 기존 `.env`가 있으면 거부하고 `openssl rand`로 새 값을 만든다. 재사용 없음 |
| 공유 서비스·DB 비노출 | Postgres는 `internal` 네트워크·게시 포트 없음, relay는 `127.0.0.1:8080`만 게시. cloudflared `TUNNEL_TOKEN` 환경은 `!reset`으로 제거된다 |
| backup/restore | pg_dump -Fc 0600, 네트워크 없는 임시 컨테이너 복원. OPS 세션 기록에서 `437f143` 코드로 3회 `tables=2 restore=0` 확인 |
| rollback | 노출 중지·전체 중지·철회는 있다. 배포 SHA 갱신·코드 rollback·라이브 복구 절차는 없다(M5) |
| restart/자원 제한 | unless-stopped, 메모리·CPU·pids 제한, cap_drop ALL·no-new-privileges·read_only(postgres 제외), 로그 10m×3 |
| TLS | 엣지 TLS는 Cloudflare가 종료한다. Tunnel 원점은 Docker 내부 네트워크의 HTTP다 |
| 현재 API schema | Cloudflare OpenAPI(`/accounts/{account_id}/access/apps` POST): `destinations`·`policies[{id,precedence}]`·`allowed_idps` 필드 일치. `domain`과 `destinations` 동시 전송 허용. policies·apps·IdP 목록은 페이지네이션이 있다 |
| 읽기 전용 호출 기록 | 요약 표만 있고 원시 기록은 없다(L5) |

## 발견 사항

미해결 critical/high는 없다. 아래 항목은 모두 미해결이다.

| ID | 심각도 | 파일·줄 | 내용과 영향 | 권고 |
|---|---|---|---|---|
| M1 | medium | `ops-guide.md`:14-15, 92, 133-136, 143 | 서두와 9·10장이 "계획, 컨테이너·Tunnel 미변경, D11/D13 미작성"으로 남아 11장·D13과 모순된다. 운영자가 현재 상태를 잘못 판단할 수 있다 | 역사적 기록으로 표시하거나 현재 상태로 갱신 |
| M2 | medium | `transition.md`:36 | 최소 권한을 `Access: Apps and Policies Edit` 하나로 적었다. `access_apply.py`는 IdP 목록을 읽으므로 `Access: Organizations, Identity Providers, and Groups Read`(또는 `Access: Identity Providers Read`)도 필요하다. 실행 기록은 `Write`로 표기했다. 이 권한만 주면 apply가 HTTP 403으로 멈춘다(fail-closed) | 권한 목록 정정 |
| M3 | medium | `beta.sh`:193-200 | expose 게이트가 `access.aud` 존재와 `required: true` 문자열만 본다. 라이브 앱·정책, `audTag`==현재 aud, teamName을 확인하지 않는다. 빈 `audTag` 설정도 `ingress validate`를 통과함을 재현했다. `dig` 실패는 "레코드 없음"으로 처리된다 | expose 전 읽기 전용 API 재확인과 aud·team 비교 추가 |
| M5 | medium | `ops-guide.md`:171-177, `transition.md`:19 | 배포 SHA 갱신·코드 rollback·라이브 복구 절차가 없다. prepare는 기존 체크아웃을 거부한다. D13이 가리키는 실행 기록에 배포 SHA가 없다(읽기 전용 확인 결과 `437f143`) | 갱신·rollback 명령과 배포 SHA 기록 |
| L1 | low | `access_apply.py`:61-77 | 사전 검사는 apps 첫 페이지의 `domain`만 본다. 사후 검증은 정책 ID만 보고 `destinations`·`overrides`(behavior public은 우회)·`allowed_idps`를 보지 않는다. 첫 계정을 쓴다 | 사후 검증 확장 |
| L2 | low | `beta.sh`:119-121 | 8080 점유 검사가 `[::]:8080`·`*:8080`을 놓친다. `verify.py local`이 보완한다 | 패턴 확장 |
| L3 | low | `beta.sh`:140-159 | backup 실패 시 부분 dump가 남는다. restore-verify는 테이블 수만 센다 | 실패 시 삭제, 행 수 비교 |
| L4 | low | `beta.sh`:167-174 | owner-login이 Basic 값을 터미널에 출력한다. agent가 실행하면 세션 기록에 남는다 | D12에 agent 실행 금지 명시 |
| L5 | low | `SAR-BETA-001-OPS.md`:23-32 | 읽기 전용 Cloudflare 재확인의 비밀 없는 원시 기록이 없다 | 경로·상태코드 기록 보존 |
| L6 | low | `user-guide.md`:65 | 5단계가 두 값을 모두 비밀번호 칸에 넣는 것으로 읽힌다 | 칸별로 분리 |

M3은 원점 JWT가 불일치 시 거부하므로 high로 분류하지 않았다. Access 앱과 원점 검사가 동시에 빠져야 공개된다.

## 검증 및 남은 제약

실행한 최소 검사(모두 종료코드 직접 확인):

| 검사 | 결과 |
|---|---|
| `bash -n beta.sh` | exit 0 |
| `access_apply.py selftest` | exit 0 |
| `verify.py local`(라이브, 읽기 전용) | exit 0: health 200, 미지 경로 404, owner/API 401, loopback 전용, Postgres 비게시, 제한, 0600/0700 |
| `verify.py regression`(라이브, 읽기 전용) | exit 0: orca 200·s8 200·mcp 401, 호스트 cloudflared PID 506937 불변 |
| 더미 env `docker compose … --profile tunnel config` | exit 0. 병합 결과를 위 표처럼 확인 |
| 더미 값 렌더링 `cloudflared:2026.9.1 ingress validate`·`ingress rule` | exit 0, `link.knowslog.com`→`relay:8080`. 빈 audTag도 exit 0(M3) |
| 공개 차단(토큰 없음) | `dig +short link.knowslog.com`(로컬·1.1.1.1) 빈 응답, `curl https://link.knowslog.com/healthz` exit 6(해석 불가). 현재 공개 경로가 없다 |
| 배포 상태(읽기 전용) | `/home/shin/deploy/knowslink` HEAD `437f143` clean, relay `127.0.0.1:8080` healthy, postgres 비게시, cloudflared 컨테이너 없음, `.env`·자격 파일·dump 0600 |
| lint | ERROR 0 / WARNING 0 / 실행 불가 0 |
| `deliverables.py --strict`(snapshot·기록 체크아웃) | 둘 다 exit 0, 검사 13 / 문제 0 / 경고 0 |
| `review.py check --from 1314e7f --to 437f143` | exit 0, reviewed 11 / skipped 0 |

실행하지 않은 것과 영향:

- Access 앱 생성·DNS·connector·`verify.py public`·이메일 로그인: 쓰기 권한이 없고 이번 범위가 아니다. 공개 Access 동작과 원점 JWT 거부는 관찰하지 못했다.
- Cloudflare API 쓰기 호출: 수행하지 않았다. 스키마는 OpenAPI 검색과 공식 문서로만 대조했다.
- 제품 전체 QA: 기존 수락 근거를 재사용했다.

## 검토 결론

`437f143`의 배포 설정과 운영 문서에 미해결 critical/high는 없다. 이 SHA의 설정 리뷰는 기록 통합이 가능하다.
M2·M3은 expose 전에 해결하기를 권고한다. M1·M5는 운영 문서 정확성을 위해 다음 OPS 과제에서 고친다.
공개 수락은 **held**다. Access 앱·DNS·connector·공개 negative·사용자 이메일 로그인(인간 검사)이 아직 실행되지 않았다. 이 리뷰는 공개 수락 근거가 아니다.
이 결과는 AI 검토자의 판단이며 OCR 자동 판정이 아니다. `review.py check` 통과는 기록 검사다.
