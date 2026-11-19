---
title: SAR-SCHEDULE-001-CODE-REVIEW 리뷰
status: review
updated: 2026-10-21
owner: reviewer
summary: cc6f836b2394 병합 전 독립 리뷰 보고서. critical·high·medium 0, low 6(비차단). 수락.
---

# SAR-SCHEDULE-001-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: reviewer 세션 `schedule001-code-review-opus` / Claude Code / claude-opus-5-5. 구현자 세션 `ctx_0a6fa8252664`와 다르다.
- base SHA / head SHA / merge-base: `826ea80be09c09d67fda75f39fd6a7e7fd32f494` / `cc6f836b2394e9350c3b57a6c1078499ab7598f6` / `826ea80be09c09d67fda75f39fd6a7e7fd32f494`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-16a1d653ed824f0cb39a415188ae2844`(detached, HEAD cc6f836, `git status` clean). 리뷰 중 snapshot을 수정하지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.2(55450a9) delegate. 규칙은 coor가 prepare한 `rules.json`(rule_sha256 `f1061481…`). 공통 규칙 fullops-common-0.3.3(security·testing·coding-style), `.fullops-squad/project.md`.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md`(SD-01–06, D010 수치, 확정 제품 규칙, SC-01–SC-09). 지시서에 적힌 `docs/exec-plans/phases/SAR-SCHEDULE-001.md`는 snapshot에 없다. 정본은 위 product-spec이며, DEV 지시서도 이 경로를 정본으로 지정한다. 완료 보고는 `handovers/logs/2026-10-21_to_dev.md` '완료 보고' 절이다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 108 / 52 / 56 / 46 / 62. coverage 100%(pending 0). 지시서 소유 범위 밖 변경은 `adapters/src/core.ts`·`synthetic.ts`다. 결과 전달에 필요하며 완료 보고 (4)에 근거가 있다. 수락한다.
- lint(`lint.json`): coor가 snapshot에서 실행했다. ERROR 0 / WARNING 15 / 실행 불가 2(product-lint·product-test, Windows에 `make` 없음).
- SIZE-002: 지시서는 '예상 규모가 크다'고 적었다. 실제 추가 1,883줄 중 제품 약 650줄, 나머지는 테스트·합성 제공자·fixture·문서다. 연결·조회·공개·철회가 한 상태 모델과 승인 handler를 공유한다. 나누면 '승인됐지만 조회 없음' 중간 상태가 생긴다. 분할하지 않은 이유를 수락한다.
- SIZE-001(13건): 새 `calendar.go` 381줄과 기존 큰 파일의 좁은 델타다. 이 과제에서 분할하지 않는 것을 수락한다.
- SEC-001 `calendar_provider_test.go:64`: `ClientSecret: "local-fixture"`는 합성 제공자 전용 더미다. 오탐으로 판단한다.
- DEP-001: 해당 없음. `go.mod`·`package.json` 변경이 없다. 봉인은 stdlib `crypto/aes`·`crypto/cipher`, 토큰 갱신은 기존 `golang.org/x/oauth2`를 쓴다.
- UI 디자인: 정본은 D04와 SAR-SCHEDULE-001 '화면 문구 원칙'이다. 기존 설정 절·승인 상세 구조와 i18n을 재사용한다. 새 CSS·색 토큰·JS는 없다. 테마 전환은 해당 없음이다. DEV 렌더 13화면×3폭 넘침 0(`render/probe.json`)을 증거로 읽었다. 시각 수락은 designer 몫이다.

## 검토 범위

result.json의 108개 `(path, status)` 모두에 `review_status`와 `reason`을 기록했다.

- reviewed 46: 제품 Go·TS 코드, 테스트, golden, env·compose, 기술·운영 문서, DEV 보고·`remote.py`·`probe.mjs`.
- skipped 62: 원시 로그 14, PNG 39, jev packet·route·find 6, `probe.json`, handover log, `raw-log-sha256.txt`. 하네스가 공통 사유로 미리 채웠다. 제외 증거는 아래 방법으로 확인했다.
  - 로그 14개는 `sha256sum`으로 `raw-log-sha256.txt`와 대조했다. 모두 일치한다. 각 로그의 마지막 EXIT 줄도 읽었다.
  - handover log는 '완료 보고' 절까지 직접 읽었다.
  - PNG·jev는 열지 않았다. DEV 보고와 packet 요약으로 존재만 확인했다. 영향: 시각 검수와 탐색 정확도는 이 리뷰 범위가 아니다.
- 삭제된 `.fullops-squad/handovers/to_dev.md`는 work.py finish가 로그로 옮긴 결과다. 로그에 전문이 있다.

중점 확인 결과:

| 중점 | 확인 | 근거 |
|---|---|---|
| refresh token 봉인 | AES-256-GCM, 랜덤 nonce, 키는 env `KNOWSLINK_CALENDAR_KEY`(base64 32B)에서만 읽는다. 잘못된 키는 기동 실패, 없으면 기능이 꺼진다. 평문은 상태·로그·화면에 없다. AAD 없음(L3). | `calendar.go:93-119`, `google.go:50-55` |
| OAuth 증분 동의 | 캘린더 동의만 `openid email calendar.freebusy`·`access_type=offline`·`include_granted_scopes=true`를 쓴다. 로그인 scope는 그대로다. reauth 경로와 `signIn(requiredMember)`로 같은 Google sub만 허용한다. 응답 scope에 freebusy가 없거나 refresh token이 없으면 저장하지 않는다. | `google.go:142-150, 221-242`, `identity.go:225-229` |
| revoke | 해제는 lock 안에서 토큰을 삭제하고, lock 밖에서 provider revoke를 요청한다. revoke 실패는 상세 없이 로그만 남기고 로컬 삭제를 막지 않는다. | `calendar.go:347-427` |
| 승인 1회=조회 1회 | approve transaction이 `Gate.Lookup{Started}` lease를 기록한다. lock 밖에서 조회하고, 두 번째 transaction이 `State==""`일 때만 결과를 1회 기록한다. 중복 제출은 `g.State != pending`으로 409다. 재시작은 lease 1분 만료 뒤 `failed`로 닫고 재조회하지 않는다. | `http.go:668-706`, `calendar.go:199-345` |
| 공개 1회(gate-consume) | lookup이 끝나기 전에는 agent에 `pending`이다. consume은 1회만 허용하고 canonical disclosure를 돌려준다. `relay.result done`은 저장 disclosure와 canonical 비교가 같을 때만 받는다. 수락 뒤 disclosure를 지운다. | `http.go:458-479`, `store.go:502-509, 582-584` |
| 범위 집행 | 단위 15·30·60, 시작 ≥ `now.Truncate(unit)`, 끝 ≤ now+14일, 같은 owner 쌍 24h 5회는 수락 시 검사하고 조회 없이 422로 거부한다. Messages 24h 보존과 revoked agent 24h 보존으로 집계 우회가 없다. 수치는 `scheduleScope` 한 곳에 있다. | `calendar.go:33-40, 121-147`, `store.go:278-299`, `connections.go:20` |
| 실패·철회 | 미연결·reconnect는 승인할 수 없다(거절 폼만). refresh 400/401은 `reconnect`와 `failed`, 일시 오류는 8초 예산 안에서 최대 3회 뒤 `failed`다. 친구 해제·차단은 기존 sweep `revoked`가 결과 기록·consume·disclosure를 막는다. 실패는 `done`으로 보낼 수 없다. | `calendar.go:216-321`, `store.go:333-351`, `http.go:587-600` |
| 동시성 | 모든 상태 변경은 직렬화된 transaction 안에 있다. 외부 호출은 lock 밖이다. 해제와 조회의 경합은 결과 기록 transaction의 `State==""` 검사로 한쪽만 이긴다. | 위 경로 |
| 로그·공개 최소화 | 로그 문자열에 식별자·제공자 상세가 없다. freebusy 응답에서 `busy[].start/end`만 읽는다. 테스트가 상태 JSON·결과·로그의 비밀 필드 부재를 확인한다. | `calendar.go:249-298`, `calendar_integration_test.go:169-267` |
| commit stub | `schedule.commit`은 lookup이 없고 consume `disclosure:false`, `done`은 403이다. | `store.go:502-509`, 테스트 `TestCalendarUnconnectedAndCommitStub` |

## 발견 사항

critical 0, high 0, medium 0, low 6이다. 모두 미해결이며 병합을 막지 않는다. 후속 과제 후보로 PLANS에 남길 것을 권한다.

| # | 심각도 | 파일·줄 | 재현 조건 | 영향 | 수정 제안 |
|---|---|---|---|---|---|
| L1 | low | `internal/relay/store.go:502-518` | 첫 `done` 수락 뒤 같은 key·같은 내용으로 재전송 | 기존 receipt 대신 403 `disclosure_denied`. 동봉 adapter는 매번 새 key를 써서 재현되지 않는다 | done 검사를 idempotency 분기 뒤로 옮긴다 |
| L2 | low | `internal/relay/calendar.go:347-371` | 승인·조회 완료 후 consume 전에 회원이 연결 해제 | 해제 뒤에도 이미 승인한 반올림 결과가 1회 공개될 수 있다. 안내 문구와 어긋날 수 있다 | 해제 때 미전달 disclosure 삭제 여부를 designer가 판단한다 |
| L3 | low | `internal/relay/calendar.go:93-119` | DB 쓰기 권한자가 회원 사이 `Sealed`를 교환 | 다른 회원 캘린더로 조회 | member ID를 GCM AAD로 묶는다 |
| L4 | low | `internal/relay/google.go:221-247` | 캘린더 동의에서 다른 Google 계정 선택(identity_mismatch) | 그 계정에 freebusy grant가 남는다. token은 보관하지 않는다 | 저장하지 않는 refresh token은 revoke한다 |
| L5 | low | `internal/relay/calendar.go:389-394` | revoke 응답이 200이 아님 | 응답 body를 닫지 않는다 | `defer resp.Body.Close()` |
| L6 | low | `internal/relay/http.go:343-346` | owner-revoke로 소유자 비활성화 | 봉인 token이 남고 provider revoke가 없다 | 비활성화 때 연결 해제와 revoke를 함께 실행한다 |

참고(수정 요청 아님):

- 정본의 '결과를 받은 쪽은 바쁨과 정보 없음을 구분한다' 규칙은 요청자 agent 표현 쪽이다. 현재 adapter·SKILL에는 요청자 쪽 표현 안내가 없다. 화면은 이 과제의 SC-07 대상이 아니다. designer·tester가 확인한다.
- 시작 시각 허용은 `now.Truncate(unit)`이다. 따라서 최대 한 단위 과거 slot까지 조회될 수 있다. DEV 문서에 명시돼 있어 수락한다.
- 완료 보고의 검토 필요 (1) 미연결 회원의 거절 허용, (2) 이메일 회원 안내, (3) DEV 작성 운영 문구는 제품·designer 확인 사항이다. 코드 결함은 아니다.

## 검증 및 남은 제약

- 재실행 없음. 지시대로 lint·테스트를 다시 실행하지 않았다. 임시 테스트를 만들지 않았고 구현자 워크트리도 수정하지 않았다.
- 제품 검사 증거: DEV 보고의 원격 Linux 격리 clone 실행이다. 제품 SHA `dfd5500`에서 make lint·make test·integration(`-skip TestRemoteConsentBrowserReturn`)·fixture-compile·fixture 기동이 모두 EXIT=0이다. 로그 SHA256이 기록과 일치한다.
- SHA 동일성 근거: `git diff --name-only dfd5500..cc6f836`은 `.fullops-squad/` 아래 파일만 바꾼다. 따라서 제품 코드는 `dfd5500`과 head가 같다.
- 환경 근거: 원격 clone은 자기 Postgres 컨테이너를 썼다(`remote.py`). 의존성 변경이 없다.
- 7ad11be·c0eb0bc의 integration EXIT=1은 이후 보정 커밋으로 해소됐다. 보정 내용(JSONB 재정규화, 거절 폼 유지, 테스트 기대 보정)을 diff에서 확인했다.
- 기존 테스트 변경: `queryBody` 현재 시각화, 일부 seed의 commit 전환, 회원 캘린더 연결 seed다. 새 규칙(범위 집행·미연결 승인 불가) 반영이며 검증 완화가 아니다.
- lint.json 실행 불가 2건: Windows `make` 부재다. 같은 명령의 원격 Linux EXIT=0으로 대체 근거를 확인했다. lint.json의 command `reason` 필드는 비어 있다. check가 이를 문제로 보면 coor가 사유를 보강한다.
- 미검증: 실제 Google 계정·Cloud Console·운영 env(D011 열림), 키 교체 절차, 브라우저 접근성, 다중 프로세스 실제 경합. 합성 제공자 통과는 실계정 수락이 아니다(SC-09).

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `docs/planning/product-specs/SAR-SCHEDULE-001.md` 결정 기록·SD-01–06·확정 수치 | 확인 | 지시서가 지정한 `exec-plans/phases/SAR-SCHEDULE-001.md`는 없다(경로 오기) |
| 구조와 구현/미완료 상태 | `design-docs/architecture.md`·`module-design.md`·`data-model.md`·`crud-design.md`·`interface-design.md`의 SAR-SCHEDULE-001-DEV 절 | 확인 | 코드와 일치한다 |
| 실행·검증 방법과 증거 | `qa-reports/SAR-SCHEDULE-001-DEV-test/report.md`, `remote.py`, `raw-log-sha256.txt` | 확인 | 로그 hash 일치 |
| 운영·복구(해당 시) | `operations/transition.md` 2026-10-21 절, `docs/development.md`, `.env.example` | 확인 | 키 교체·분실 절차는 미정으로 명시됨 |
| 다음 작업·담당·재개 조건 | handover log 완료 보고 '남은 일', transition.md 운영 수락 조건 | 확인 | D011 답 뒤 실계정 1경로를 coor가 확인한다 |
| 로컬 링크·절 접근/지원 한계 | DEV 보고의 상대 링크(로그·render·remote.py) | 확인 | 깨진 링크 없음 |
| snapshot 정리 후 정본 접근 | 결과는 coor 체크아웃 리뷰 폴더에 있다 | 확인 | snapshot 정리는 coor 담당 |

## 검토 결론

수락한다. 근거는 다음과 같다.

- SC-01–SC-08의 핵심 경계가 코드와 자동 검사에 모두 있다. 봉인·증분 동의·revoke, 승인 1회=조회 1회, 결과 1회 기록, consume 1회 공개, echo 검증, D010 범위 집행, 실패·철회 분류, commit stub이 해당한다.
- 원격 Linux 증거가 head와 같은 제품 코드에서 EXIT=0이다. lint ERROR는 0이다.
- low 6건은 데이터 손실이나 무승인 공개를 만들지 않는다. 후속 과제로 처리할 수 있다.

check 통과는 이 문서의 서술과 테스트 성공을 자동 보증하지 않는다. SC-09(독립 QA·designer 시각 검수·실계정 1경로)는 남아 있다.
