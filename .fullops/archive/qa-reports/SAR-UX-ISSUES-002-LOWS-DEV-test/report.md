---
title: SAR-UX-ISSUES-002-LOWS-DEV DEV 보고
status: review
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-LOWS-DEV]
summary: "잔여 UI low 5건(F-R1·F-R2·F-R3·O-LG-05·내부 언어 표식 헤더)의 수정, 수정 전 실패 확인, 원격 Linux make lint/test/integration과 FullOps lint 근거"
---

# SAR-UX-ISSUES-002-LOWS-DEV DEV 보고

DEV 검사 기록이다. 독립 리뷰·tester QA·designer 시각 수락이 아니다.

- 기준 / 제품 커밋: `c3b51926dbae01bdb325a9d3be01ed1ad09275c0`(지시서 base, main `4a4cdc0`과 `internal/`·`cmd/`·`adapters/` 차이 없음) / `e9da442e9b3d111c2073b9edd0c9d6f38e54713b`.
- 적용 기준: fullops-common-0.3.3, `.fullops-squad/project.md`, [SAR-UX-ISSUES-002.md](../../../design-docs/SAR-UX-ISSUES-002.md) U4(재요청 결정·차단 비공개)·U3.
- 원천: [CONNECT-FIX-REVIEW](../SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW-review/report.md) F-R1–F-R3, [LANG-TESTER](../SAR-UX-ISSUES-002-LANG-TESTER-test/report.md) O-LG-05·O-LG-06 끝부분, LANG-DEV 보고 ‘남은 한계’.

## 수정

| ID | 원인 | 수정 | 검사 |
|---|---|---|---|
| F-R1 | 수신 한도 보류 때 `Inviter`가 새 요청자로 바뀌어, 예전 발신자였던 받은 사람에게 행이 받은 쪽 기록으로 보였다 | `Friend.HeldFrom`·`HeldTail`(생략 가능)에 이전 요청자·끝 4자를 보존한다. `humanView`는 받은 사람이 `HeldFrom`이면 보낸 기록과 그 끝 4자로 보인다. 두 필드가 없는 기존 행은 이전과 같다 | `TestFriendRerequestReusesOneRow/receive_limit_reversed`: 받은 사람의 렌더 화면이 재요청 전·후·요청자 취소 뒤 같다 |
| F-R2 | block 판정이 `State == "pending"`만 대기로 봐서 보류 행이 `removed`가 되었다 | 받은 사람의 기한 전 `blocked-pending`도 대기처럼 `blocked-pending`으로 남긴다 | `TestFriendBlockOnHeldRequestKeepsRequesterPending`: 요청자 화면이 대기 |
| F-R3 | 쌍 행이 있으면 열린 대리 행 검사를 건너뛰었다 | 대리 행 검사를 쌍 행 유무와 무관하게 한다. 대리 행이 대기 중이면 무동작이다. 대리 행 이전(migrate)은 계속 쌍 행이 없을 때만 한다 | `TestFriendRepeatWithOpenLegacyPlaceholderIsNoop` |
| O-LG-05 | 오류 문장이 일반 텍스트로만 렌더됐다 | `intro`의 Problem을 `langPart`로 렌더한다. 영어판의 ‘한국어’에 `lang="ko"`, 한국어판의 ‘English’에 `lang="en"` span을 붙인다. 나머지는 기존처럼 이스케이프한다. 레이아웃·문구 변경 없음 | `TestLanguagePartsAndMarkers` |
| 내부 표식 헤더 | `render`만 표식을 지워 303·JSON 응답에 `X-KnowsLink-Lang-Saved: 1`이 남았다 | `boundedHTTP`가 응답을 `langMarks`로 감싸 헤더 전송 직전(`WriteHeader`·첫 `Write`)에 `Lang-Saved`·`Lang-Unsaved`를 지운다. `Unwrap`으로 `ResponseController` deadline·flush를 유지하고 본문 한도(`MaxBytesReader`)는 원래 writer를 쓴다. `/v1/text/events`는 이 경로 밖이며 표식을 쓰지 않는다 | `TestLanguagePreferencePostgres`: 저장값 회원의 `POST /home/language`·`/auth/logout` 303 송신 헤더에 표식 없음. 단위 검사도 redirect를 본다 |

O-LG-06 (1)–(3)과 O-CT-03·V-AP-01·O-LG-01·02·04는 지시대로 바꾸지 않았다.

## 검사

| 검사 | 대상 / 위치 | 결과 | 근거 |
|---|---|---|---|
| 수정 전 실패(친구 3건) | Windows 로컬, 새 검사 + 수정 전 코드 | FAIL 3 (EXIT=1) | [prefix-friend-red.log](prefix-friend-red.log) |
| 수정 전 실패(O-LG-05) | Windows 로컬, `member_pages.go`의 한 줄만 일시 복원 | FAIL (EXIT=1) | [prefix-langpart-red.log](prefix-langpart-red.log) |
| 수정 전 실패(전체) | 원격 Linux, `e9da442` clone 사본에서 제품 5파일과 `i18n_test.go`만 `c3b5192`로 되돌림, integration 태그 | 4건 모두 FAIL(`GO_TEST_EXIT=1`): F-R1·F-R2·F-R3·303 표식 | [red 로그](e9da442e9b3d111c2073b9edd0c9d6f38e54713b-red.log) |
| `make install` | 원격 Linux 고유 `/tmp` clone, HEAD `e9da442` 확인 | EXIT=0 | [로그](e9da442e9b3d111c2073b9edd0c9d6f38e54713b-install.log) |
| `make lint` | 같은 clone, clean | EXIT=0 | [로그](e9da442e9b3d111c2073b9edd0c9d6f38e54713b-lint.log). ‘invalid tunnel UUID’는 기존 음성 검사 출력이다(이전 단위 로그에도 있음) |
| `make test` | 같은 clone | EXIT=0 | [로그](e9da442e9b3d111c2073b9edd0c9d6f38e54713b-test.log) |
| integration(`-race -tags=integration`, `TestRemoteConsentBrowserReturn` 제외) | 같은 clone, 고유 Postgres 17 컨테이너 | `ok … 137.047s`, EXIT=0 | [로그](e9da442e9b3d111c2073b9edd0c9d6f38e54713b-integration.log) |
| FullOps lint.py `--from c3b5192` | `e9da442` 깨끗한 임시 detached worktree(`git worktree remove`, force 없이 제거) | ERROR 0, WARNING 4, 실행 불가 2 | [로그](e9da442-fullops-lint.log) |

- `TestRemoteConsentBrowserReturn` 제외는 이전 단위와 같다(실제 Chromium 필요).
- lint WARNING은 SIZE-001 4건이다. `PLANS.md`(coor 변경, 이 과제 아님)와 이미 상한을 넘은 `human.go`(+3)·`human_http.go`(+4)·`i18n_test.go`(+19)다. 기존 함수 안 분기라 분할하지 않았다.
- 실행 불가 2(product-lint·product-test)는 Windows에 make가 없어서이며 위 원격 Linux EXIT=0으로 대신한다.
- DEP-001 없음. `go.mod`·`go.sum`·`package.json` 변경이 없다. 새 import는 표준 라이브러리 `fmt`(i18n.go)뿐이다.
- 원격 사본: 원격 clone과 `red.*` 사본, 정지한 컨테이너 두 개(`kl-lows002-integration-…`, `kl-lows002-red-…`)는 이전 단위처럼 보존했다. private 도구 디렉터리는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/lows-ux002-dev`이며 실행 스크립트는 [remote.py](remote.py)다.

## 영향과 호환

- 상태 JSON: `Friend.HeldFrom`·`HeldTail` 추가(생략 가능). migration 없음. 되돌리기 시 이전 코드는 두 필드를 무시해 F-R1 이전 표시로 돌아간다.
- 문서: [data-model](../../../design-docs/data-model.md)·[interface-design](../../../design-docs/interface-design.md)·[module-design](../../../design-docs/module-design.md)·[D13](../../../operations/transition.md).
- UI: 시각 변화 없음. 오류 문장에 `lang` span만 붙는다. 테마·토큰 변경 없음.

## 검증하지 못한 것

- 스크린리더 실제 낭독(O-LG-05)과 실제 브라우저 응답 헤더 관찰. tester 재확인이 필요하면 coor가 조정한다.
- 독립 코드 리뷰(고정 SHA)는 coor 조정 대상이다.
