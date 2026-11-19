---
title: SAR-BETA-001 베타 운영 수정 독립 재리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-BETA-001-REVIEW-FINAL]
summary: 고정 SHA f824015의 원래 지적 수정 여부와 새 expose·deploy 경계를 독립 재리뷰한 결과와 공개 held 결론을 기록한다
---

# SAR-BETA-001 베타 운영 수정 독립 재리뷰

- 검토자 / CLI / 모델: Orca dispatch `ctx_b100d949407f`(task `task_df231db29363`, run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 ID는 `941501bd-e2b4-4578-b092-e6c70cdbd306`이다(`~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-dev/`).
- 구현 세션: OPS `2191cc9b-76ef-4522-9fad-d2c9f017bfbc`. 세션 기록의 cwd는 `fullops-ops`이며 `8a7ad36`·`f824015`와 `access.live.json` 작업이 있다. 리뷰 세션과 다르다. 이전 리뷰 세션 `bac29ce7`과도 다르다.
- base SHA / head SHA / merge-base: `437f1432a158670a485413c1aba159debc3759e5` / `f824015314c66bcab42940cfe3db2edabb22e1dd` / `437f143…`
- snapshot: `/tmp/knowslink-beta-review-f824015`. detached(`symbolic-ref` 실패)·HEAD `f824015`·검사 전후 `git status --porcelain` 0줄을 확인했다. Python은 `-B`·`PYTHONDONTWRITEBYTECODE=1`로 실행했다. 실행 검사는 scratch 복제본과 scratch 상태 디렉터리에서 했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate / `.fullops-squad/review/rule.json`(sha256 `ab2116fb…`, 이전 리뷰와 같다)
- 공통 규칙: `fullops-common-0.3.2`, `project.md`, 문서 작성 규칙, ponytail full, `cloudflare-one` 규칙(현재 API 스키마 조회). 예외 없음.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_dev.md`(SAR-BETA-001-REVIEW-FINAL), 이전 리뷰 [SAR-BETA-001-REVIEW-review/report.md](../SAR-BETA-001-REVIEW-review/report.md)(`f625c4e`, 보존·수정하지 않음).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 6 / 6 / 0 / 6 / 0. coverage 100%.
- lint(`lint.json`): ERROR 0 / WARNING 0 / 실행 불가 0, product-lint passed. `f824015`의 scratch clone에서 `npm ci --prefix adapters` 뒤 실행했다. clone 작업 트리는 깨끗했다.

## 검토 범위

result.json의 6개 `(path,status)`를 모두 reviewed로 기록했다. 변경은 `deploy/knowslink/{access_apply.py,beta.sh}`와 OPS 문서 4개다.
제품 diff가 없음을 직접 확인했다. `git diff --name-only 78b1d92 f824015`에는 `.fullops-squad/`와 `deploy/` 밖의 파일이 없다. 제품 전체 QA는 반복하지 않았다.
`8a7ad36..f824015`는 문서 2개(3줄)만 다르다. 따라서 OPS의 `8a7ad36` 재검증 증거는 코드 기준으로 `f824015`에 적용된다.
`verify.py`와 `tunnel/config.yml.tmpl`은 이번 diff에 없다. 이전 리뷰 근거를 재사용하고 `verify.py local/regression`만 재실행했다.

### 원래 지적의 수정 여부

| ID | 판정 | 근거 |
|---|---|---|
| M1 | 해소 | D12 서두·6·9·10장이 역사적 기록으로 표시됐다. `미작성`·`모든 단계가 미실행` 서술이 남지 않았다 |
| M2 | 해소 | D12 11.2·D13·실행 기록에 IdP 읽기 권한을 추가했다. 문서 표기 Write와 토큰 화면 Edit의 관계도 적었다 |
| M3 | 해소(새 N1·N2 주의) | `expose`가 `access_apply.py check`를 호출한다. live 앱 destinations·정책 1개·이메일 단독 allow·IdP 1개·우회 옵션·`aud`==기록==`audTag`·`teamName`을 비교한다. `dig` 실패는 중단한다 |
| M5 | 해소 | `beta.sh deploy <sha>`가 백업·이력·migration 차이 차단·재빌드·`verify.py local`을 수행한다. D12 11.3에 rollback을 적었다 |
| L1 | 부분 해소 | 사후 검증이 넓어졌고 계정 1개를 강제한다. 사전 중복 검사는 여전히 첫 100개 앱의 `domain`만 본다. 상한을 문서에 적었으므로 수락한다 |
| L2 | 해소 | `ss -ltnH` 패턴이 `127.0.0.1`·`0.0.0.0`·`*`·`[::]`를 포함한다 |
| L3 | 해소 | 실패 시 부분 dump를 삭제한다(재현). backup·restore-verify가 `relay_state` 행 수를 출력한다 |
| L4 | 해소 | D12 11.5·D11 2단계가 agent 실행 금지를 명시한다 |
| L5 | 해소 | 실행 기록에 비밀 없는 GET 경로·상태 표가 있다 |
| L6 | 해소 | D11 5단계가 칸별로 나뉘었다 |

### 새 경계 검토

| 항목 | 판정 |
|---|---|
| Access 정책·IdP·단일 이메일 | `only_owner`가 decision allow·include 이메일 하나·exclude/require 비어 있음을 본다. `everyone`·다른 이메일·IdP 2개는 거부됐다(독립 실행) |
| domain·destinations | `domain`은 없음 또는 정확한 host, destinations는 `[public, link.knowslog.com]` 하나만 허용 |
| aud·team 일치 | live `aud`==`access.aud`==config의 64-hex `audTag` 목록 전체. 추가 audTag와 다른 aud는 거부됐다. `teamName`은 config와 `KNOWSLINK_TEAM`(beta.sh `render_config`와 같은 기본값)만 비교한다(N5) |
| stale/missing MCP 증거 | 없음은 `FileNotFoundError`로 exit 1, 10분 초과는 거부. 미래 mtime과 다른 app id는 통과한다(N2) |
| 우회 방지 | 일반 실행에서는 fail-closed다. `PYTHONOPTIMIZE`에서는 모든 assert가 사라져 통과한다(N1) |
| DNS 실패 차단 | `dig` 비정상 종료는 중단한다. SERVFAIL은 exit 0·빈 출력이라 "없음"으로 처리된다. 이어지는 `tunnel route dns`가 기존 레코드를 덮어쓰지 않으므로 결과는 fail-closed다. 별도 지적하지 않는다 |
| deploy 자기 자원 제한 | `dc`는 `-p knowslink`와 `$DEPLOY`·`$STATE`만 쓴다. 다른 Compose 프로젝트·볼륨을 건드리지 않는다. `up`은 relay만 기동한다 |
| deploy 백업·migration 차단 | migration 차이와 unknown commit은 백업·이력·체크아웃 전에 중단한다. pg_dump 실패는 체크아웃 전에 중단하고 dump를 지운다(재현) |
| 자기 교체 안전성 | 중괄호 dispatch와 `exit`로 bash가 실행 전에 전체를 읽는다. checkout 뒤의 `up`은 메모리의 이전 정의, `verify.py`는 새 SHA 파일을 쓴다. 타당하다 |
| 현재 API schema | Cloudflare OpenAPI 검색: `GET /accounts/{account_id}/access/apps/{app_id}` 응답에 `aud`·`options_preflight_bypass`·`allowed_idps`·`destinations`·`custom_pages`·`policies`·`domain`이 있다. policy GET에 `decision`·`include`·`exclude`·`require`가 있다. apps 목록은 `per_page`·`page`를 받는다 |

## 발견 사항

미해결 critical/high는 없다. 아래 항목은 모두 미해결이다.

| ID | 심각도 | 파일·줄 | 내용과 영향 | 권고 |
|---|---|---|---|---|
| N1 | medium | `access_apply.py`:42-53 | `verify_live`의 모든 검사가 `assert`다. `PYTHONOPTIMIZE=1`에서 aud가 다른 live 앱으로 `check`가 exit 0으로 통과했다(재현). expose 게이트가 fail-open된다. 원점 JWT `audTag`가 별도로 막으므로 high가 아니다 | 명시적 `raise SystemExit`로 바꾸거나 `beta.sh`가 `python3 -I`·`PYTHONOPTIMIZE` 해제로 호출 |
| N2 | low | `access_apply.py`:114-126 | MCP 증거 `access.live.json`은 출처 검증 없이 mtime만 본다. mtime이 미래이면 계속 통과하고, 스냅숏 app id가 `access.json`과 달라도 통과한다(둘 다 재현). aud 비교가 앱을 묶으므로 영향은 작다 | `saved` app/policy id 비교, 음수 나이 거부 |
| N3 | low | `transition.md`:19, `ops-guide.md`:172 | 실제 `deploy-history.log`는 `8a7ad36->437f143`, `8a7ad36->8a7ad36` 두 줄이다. `437f143->8a7ad36` 복귀와 현재 HEAD `f824015`로의 이동은 수동 체크아웃이라 기록이 없다. 문서는 로그를 현재 SHA 근거·"이력 기록"으로 적었다 | `rev-parse`를 정본으로 명시하고 수동 이동도 기록 |
| N4 | low | `ops-guide.md`:172 | `8a7ad36` 이전 SHA로 rollback하면 그 SHA의 `beta.sh expose`는 live 검증이 없는 이전 게이트다. 11.3에 이 주의가 없다 | expose는 `access_apply.py check`가 있는 SHA에서만 실행한다고 명시 |
| N5 | low | `access_apply.py`:72 | selftest의 `config.replace("team", "other")`는 `teamName` 키 자체를 바꾼다. 값 불일치만을 격리하지 않는다. 독립 실행으로 값만 다른 경우도 거부됨을 확인했다. `teamName`은 live organization과 비교하지 않는다(불일치 시 원점 거부로 fail-closed) | 값만 바꾸는 사례로 교체 |

N1은 expose 전 수정을 권고한다. N2~N5는 다음 OPS 과제에서 고쳐도 된다. 직접 구현하지 않고 coor에 돌려준다.

## 검증 및 남은 제약

독립 실행 검사(모두 종료코드 직접 확인, 구현자 증거와 별도):

| 검사 | 결과 |
|---|---|
| `bash -n beta.sh` | exit 0 |
| `access_apply.py selftest` | exit 0 |
| `check` 임시 상태(더미 이메일·더미 aud, 0600/0700) | fresh 정상 0 / missing 1 / stale(11분) 1 / 미래 mtime 0(N2) / aud 불일치 1 / 다른 app id 0(N2) / IdP 2개 1 / 다른 이메일 1 / everyone 1 / audTag 추가 1 / teamName 값만 다름 1 / 이메일 파일 0644 1 |
| `PYTHONOPTIMIZE=1 check`(aud 불일치) | exit 0(N1) |
| `beta.sh deploy`(scratch clone, 가짜 docker, scratch 상태) | unknown commit exit 1, `bd9c8dc^` migration 차이 exit 1, `437f143`(migration 같음) pg_dump 실패 exit 1. 세 경우 모두 HEAD `f824015` 불변·dump 0개·이력 0줄 |
| 실제 배포 상태(읽기 전용) | `/home/shin/deploy/knowslink` HEAD `f824015` clean, relay `127.0.0.1:8080` healthy, postgres 비게시, cloudflared 컨테이너 없음, `access.aud`·`access.json` 없음(Access 앱 미생성), backups 0600 |
| `verify.py local`(라이브, 읽기 전용) | exit 0 |
| `verify.py regression`(라이브, 읽기 전용) | exit 0: orca 200·s8 200·mcp 401, PID 506937 불변 |
| 공개 차단 | `dig +short link.knowslog.com A` exit 0·빈 응답. 공개 경로가 없다 |
| Cloudflare OpenAPI 검색 | 위 필드 확인. 계정 API 호출은 하지 않았다 |
| lint | ERROR 0 / WARNING 0 / 실행 불가 0 |
| `deliverables.py --strict`(snapshot·기록 체크아웃) | 둘 다 exit 0, 검사 13 / 문제 0 / 경고 0 |
| `review.py check --from 437f143 --to f824015` | exit 0, reviewed 6 / skipped 0. jev_find 점수 recall 0.5 / precision 0.083(수락 판단에 쓰지 않음) |

실행하지 않은 것과 영향:

- Access 앱 생성·DNS·connector·`verify.py public`·이메일 로그인: 쓰기 권한이 범위 밖이다. live API로 `check`를 실행하지 못했다. 응답 모양은 OpenAPI와 더미 스냅숏으로만 검사했다.
- 실제 Postgres에 대한 `deploy`·`backup`·`restore-verify`: 라이브 DB를 바꾸지 않으려고 실행하지 않았다. 구현자 증거(rollback·forward exit 0, `tables=2 relay_state_rows=1`)는 구현자 자체 증거로만 인용한다.
- 제품 전체 QA: 제품 diff가 없으므로 기존 수락 근거를 재사용했다.

## 검토 결론

`f824015`의 배포 설정과 운영 문서에 미해결 critical/high는 없다. 원래 medium 4건과 low 5건은 해소됐고 L1은 문서화된 상한으로 수락한다. 이 SHA의 리뷰 기록은 통합할 수 있다.
N1은 expose 게이트가 특정 환경에서 fail-open되므로 expose 전에 고치기를 권고한다.
공개 수락은 **held**다. Access 앱·DNS·connector·공개 negative·사용자 이메일 로그인(인간 검사)이 아직 실행되지 않았다. 이 리뷰는 공개 수락 근거가 아니다.
이 결과는 AI 검토자의 판단이며 OCR 자동 판정이 아니다. `review.py check` 통과는 기록 검사다.
