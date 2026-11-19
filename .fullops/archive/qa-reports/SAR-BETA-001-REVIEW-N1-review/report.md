---
title: SAR-BETA-001 Access 게이트 수정 독립 검토
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-BETA-001-REVIEW-N1]
summary: 고정 SHA 28bd1bb의 N1-N5 수정과 Access 게이트 negative의 일반·최적화 실행을 독립 검토한 결과와 공개 held 결론을 기록한다
---

# SAR-BETA-001 Access 게이트 수정 독립 검토

- 검토자 / CLI / 모델: Orca dispatch `ctx_9779a5971c1d`(task `task_be8f037e8453`, run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 ID는 `fdfab4f3-8e18-4129-94b1-ac6f85ce216f`이다(`~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-dev/`).
- 구현 세션: OPS `2191cc9b-76ef-4522-9fad-d2c9f017bfbc`. 리뷰 세션과 다르다. 이전 리뷰 세션 `941501bd`(REVIEW-FINAL)·`bac29ce7`(REVIEW)과도 다르다.
- base SHA / head SHA / merge-base: `f824015314c66bcab42940cfe3db2edabb22e1dd` / `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2` / `f824015…`. 사이 커밋은 `28bd1bb` 하나다.
- snapshot: `/tmp/knowslink-beta-review-28bd1bb`. detached(`symbolic-ref` 실패)·HEAD `28bd1bb`·검사 전후 `git status --porcelain` 0줄을 확인했다. snapshot 파일은 수정하지 않았다. 실행 검사는 `python3 -B`·`PYTHONDONTWRITEBYTECODE=1`과 scratch 상태 디렉터리를 사용했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate / `.fullops-squad/review/rule.json`(sha256 `ab2116fb…`, 이전 리뷰와 같다)
- 공통 규칙: `fullops-common-0.3.2`, `project.md`, 문서 작성 규칙, ponytail full. 예외 없음.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_dev.md`(SAR-BETA-001-REVIEW-N1), 이전 리뷰 [SAR-BETA-001-REVIEW-FINAL-review/report.md](../SAR-BETA-001-REVIEW-FINAL-review/report.md)(`12a88b2`, 보존·수정하지 않음).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 4 / 4 / 0 / 4 / 0. coverage 100%.
- lint(`lint.json`): ERROR 0 / WARNING 0 / 실행 불가 0, product-lint passed. `28bd1bb`의 scratch clone에서 `npm ci --prefix adapters` 뒤 실행했다. clone 작업 트리는 깨끗했다.

## 검토 범위

result.json의 4개 `(path,status)`를 모두 reviewed로 기록했다. 변경은 `deploy/knowslink/access_apply.py`와 OPS 문서 3개다.
제품 diff가 없음을 직접 확인했다. `git diff --name-only 78b1d92 28bd1bb`에는 `.fullops-squad/`와 `deploy/` 밖의 파일이 없다.
`beta.sh`와 `verify.py`는 `f824015..28bd1bb`에서 바뀌지 않았다. 이전 리뷰의 제품·설정·deploy 검사는 반복하지 않았다.
이번 검토는 게이트 negative의 일반·최적화 실행과 MCP 증거(`access.live.json`)의 결합 경계에 한정했다.

### 이전 지적 N1~N5의 수정 여부

| ID | 판정 | 근거 |
|---|---|---|
| N1 | 해소 | `verify_live`·`apply`·`__main__`의 모든 게이트 `assert`가 `need()`(`SystemExit`)로 바뀌었다. 남은 `assert`는 `selftest` 63~65행뿐이다. 아래 표의 negative 10개가 네 가지 실행 모드 모두에서 exit 1이다 |
| N2 | 해소 | `0 <= age <= 600`으로 미래 mtime을 거부한다. 스냅숏의 app·policy id가 `access.json`과 다르면 거부한다 |
| N3 | 해소(잔여 R2) | D12 11.3과 D13이 `rev-parse`를 정본으로, 로그를 `beta.sh deploy` 이동 기록으로 한정한다. 실제 live HEAD는 `28bd1bb`이고 로그에는 그 이동이 없다. 문서 설명과 맞다 |
| N4 | 해소 | D12 11.3이 `check` 없는 SHA에서 `expose` 금지를 적었다. 이전 SHA의 `beta.sh`는 바꿀 수 없으므로 문서 통제가 맞는 수단이다 |
| N5 | 해소 | `selftest`가 `teamName: team`→`teamName: other`로 값만 바꾼다. `selftest`는 일반·`-O` 모두 exit 0이다 |

### 게이트 negative 독립 실행

snapshot의 `access_apply.py check`를 직접 실행했다. 상태는 scratch 디렉터리(0700)에 더미 이메일·더미 64-hex aud·더미 id로 만들었다. `CF_API_TOKEN_FILE`은 설정하지 않아 MCP 스냅숏 경로만 탔다. 표의 값은 종료코드다.

| 사례 | 일반 | `-O` | `-OO` | `PYTHONOPTIMIZE=1` | 거부 메시지 |
|---|---|---|---|---|---|
| 정상(fresh) | 0 | 0 | 0 | 0 | — |
| live aud 불일치 | 1 | 1 | 1 | 1 | app aud differs |
| `access.live.json` 없음 | 1 | 1 | 1 | 1 | `FileNotFoundError` |
| 11분 경과 | 1 | 1 | 1 | 1 | last 10 minutes |
| 미래 mtime(+5분) | 1 | 1 | 1 | 1 | last 10 minutes |
| 스냅숏 app id 다름 | 1 | 1 | 1 | 1 | snapshot app/policy ids differ |
| 스냅숏 policy id 다름 | 1 | 1 | 1 | 1 | snapshot app/policy ids differ |
| 다른 이메일 | 1 | 1 | 1 | 1 | only the owner email |
| IdP 2개 | 1 | 1 | 1 | 1 | exactly one identity provider |
| `teamName` 값만 다름 | 1 | 1 | 1 | 1 | tunnel teamName |
| config에 audTag 추가 | 1 | 1 | 1 | 1 | audTag must be exactly the app aud |

대조군으로 같은 상태에서 `f824015`의 `access_apply.py`를 실행했다. aud 불일치는 일반 1·`PYTHONOPTIMIZE=1` 0, 미래 mtime과 다른 app id는 두 모드 모두 0이었다. 따라서 이 행렬은 N1·N2를 검출하며, `28bd1bb`에서 결과가 바뀐 것은 수정 때문이다.

`beta.sh expose`는 `python3 … access_apply.py check || die`로 호출한다. 환경의 `PYTHONOPTIMIZE`가 상속돼도 위 결과대로 실패가 전달된다. `expose` 자체는 실행하지 않았다(정적 확인).

### proof 경계

- MCP 스냅숏은 이제 `access.json`의 app·policy id와 10분 이내 mtime에 묶인다. 출처 서명은 없다. 상태 디렉터리는 0700 소유자 전용이므로 위조에는 같은 사용자 권한이 필요하다. 이 수준을 수락한다.
- 토큰 경로는 `access.json`의 id로 직접 GET하므로 결합이 구조적으로 보장된다. 이번 diff는 이 경로를 바꾸지 않았다.
- `verify.py`의 사후 증명은 여전히 `assert`다(R1).

## 발견 사항

미해결 critical/high는 없다. 새 지적 2건은 모두 low이며 미해결이다.

| ID | 심각도 | 파일·줄 | 내용과 영향 | 권고 |
|---|---|---|---|---|
| R1 | low | `verify.py`:45-74(diff 밖, 이번 SHA 불변) | `regression`·`local`·`public`의 검사가 모두 `assert`다. `-O`·`PYTHONOPTIMIZE` 실행에서 검사가 사라지고 성공 메시지를 출력한다. `beta.sh deploy`의 사후 `local` 증명과 공개 negative 증명이 거짓 통과할 수 있다. 게이트가 아니고 Access·원점 JWT가 별도로 막으므로 low다. 언어 의미에 따른 정적 판단이며 실행하지 않았다 | `need()`와 같은 명시적 검사로 바꾸거나, tester·OPS 실행 기록에 `PYTHONOPTIMIZE` 미설정을 남긴다 |
| R2 | low | `ops-guide.md`:172, `transition.md`:19 | 11.3 끝의 "(둘 다 exit 0, 백업 생성, 이력 기록)"은 앞으로의 이동도 기록된 것처럼 읽힌다. 실제 로그는 `8a7ad36->437f143`과 수동 체크아웃 뒤의 no-op `8a7ad36->8a7ad36`뿐이다. D13 배포 SHA 칸은 "기록한다"와 "갱신·rollback" 사이에 마침표가 없다. 동작 영향은 없다 | 다음 OPS 문서 갱신에서 표현을 고친다 |

R1·R2는 직접 고치지 않고 coor에 돌려준다. 둘 다 `expose` 전 필수 수정이 아니다.

## 검증 및 남은 제약

독립 실행 검사(모두 종료코드 직접 확인, OPS 자체 증거와 별도):

| 검사 | 결과 |
|---|---|
| snapshot 상태 | detached, HEAD `28bd1bb`, porcelain 0줄(검사 전후) |
| 4파일 diff 커버리지 | `git diff --name-status f824015 28bd1bb` 4개 = result.json 4개, 누락 0 |
| 제품·설정 불변 | `78b1d92..28bd1bb` 제품 파일 0개, `beta.sh`·`verify.py` `f824015..28bd1bb` diff 없음 |
| `access_apply.py selftest` | 일반 exit 0, `-O` exit 0 |
| `check` negative 행렬 | 위 표: 정상 0, negative 10개 × 4개 모드 모두 1 |
| `f824015` 대조군 | `PYTHONOPTIMIZE=1` aud 불일치 0, 미래 mtime 0, 다른 app id 0(N1·N2 재현) |
| 실제 배포 상태(읽기 전용) | `/home/shin/deploy/knowslink` HEAD `28bd1bb` clean, `deploy-history.log` 2줄, `access.*` 파일 없음(Access 앱 미생성) |
| lint | ERROR 0 / WARNING 0 / 실행 불가 0 |
| `deliverables.py --strict`(snapshot) | exit 0, 검사 13 / 문제 0 / 경고 0 |

`review.py check`와 기록 체크아웃의 `--strict`·lint 결과는 완료 보고에 적는다.

실행하지 않은 것과 영향:

- Access 앱 생성·DNS·connector·`beta.sh expose`·`verify.py public`·이메일 로그인: 쓰기·노출이 범위 밖이다. live API로 `check`를 실행하지 못했다. 응답 모양은 이전 리뷰의 OpenAPI 근거와 더미 스냅숏으로만 검사했다.
- 제품 전체 QA와 deploy·backup 재실행: 이번 diff에 없으므로 이전 리뷰 근거를 재사용했다.
- R1은 실행하지 않은 정적 판단이다.

## 검토 결론

`28bd1bb`의 `access_apply.py`와 OPS 문서에 미해결 critical/high는 없다. 이전 N1~N5는 모두 해소됐다. N1의 최적화 실행 우회가 닫혔음을 네 가지 모드로 독립 확인했다. 이 SHA의 리뷰 기록은 통합할 수 있다.
새 low R1·R2는 미해결이며 다음 OPS 과제에서 처리해도 된다.
공개 수락은 계속 **held**다. Access 앱·DNS·connector·공개 negative·사용자 이메일 로그인(인간 검사)이 아직 실행되지 않았다. 이 리뷰는 공개 수락 근거가 아니다.
이 결과는 AI 검토자의 판단이며 OCR 자동 판정이 아니다. `review.py check` 통과는 기록 검사다.
