---
title: SAR-PUBLISH-001-DEV-REVIEW 리뷰
status: review
updated: 2026-10-22
owner: reviewer
tasks: [SAR-PUBLISH-001-DEV, SAR-PUBLISH-001-DEV-REVIEW]
summary: "722feed9a89b 병합 전 독립 리뷰 보고서. 서버 캘린더 제거(D018)·일회 정리·/privacy·/terms(D017)를 수락한다. critical/high 0, low 4"
---

# SAR-PUBLISH-001-DEV-REVIEW 리뷰

- 검토자 / CLI / 모델: tester worktree의 별도 reviewer 세션 `8c2de73c-2985-47a5-9f5e-a7eb008be2d4` / Claude Code / `claude-opus-5-5`. 구현자 세션 `a5727aa6-dee3-44de-8040-810796a5d9ca`와 다르다. coor가 dispatch로 배정했다(task `task_bc65c78ecea2`, dispatch `ctx_297da3025001`).
- base SHA / head SHA / merge-base: `05d0800e825768f78ca6750209d74c74aa1174c9` / `722feed9a89b5e0f415ed90f8ede64174f7b6437` / `05d0800e825768f78ca6750209d74c74aa1174c9`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-8d05ab5817d2480285051e41ecf15231`(detached, 읽기 전용으로 사용, owner coor). 리뷰는 이 snapshot과 git 객체만 읽었다. dev 워크트리는 gate 디렉터리의 lint 결과 유무만 확인했다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.2(55450a9). `.fullops-squad/review/rule.json`(rule_sha256 `f1061481…`), FullOps 플러그인 1.8.0 `review.py`.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC-01–06, P-01–05), `SAR-SCHEDULE-001.md` '철회 — D018'(기존 연결 회원 처리 1–6, RC-01–05), 기준 `handovers/to_dev.md`(05d0800e).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 119 / 64 / 55 / 68 / 51. skipped 51개는 prepare가 미리 채운 `qa-reports/`·`evaluations/jev/`·`handovers/logs/` 하네스 자동 생성 기록이다. 지시서 범위(캘린더 제거·정책 페이지·ingress 후보·문서)와 실제 diff는 일치한다. 범위 밖 제품 변경은 없다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 16 / 2. 출처는 아래 '검증 및 남은 제약'에 적었다.
- SIZE-002: 지시서는 줄 수 예상치를 따로 두지 않았다. 지시서는 '캘린더 제거 + 정책 4문안'을 한 과제로 정했다. 실제 추가 1,062줄 중 약 214줄이 정책 HTML, 약 420줄이 새 테스트·fixture, 나머지가 문서다. 제품 코드(테스트·정적 문안 제외)는 +232/−673이다. 제거 단위라 분할하지 않은 이유는 타당하다.
- DEP-001: 경고 없음. go.mod·package.json 변경이 없다. revoke는 표준 `net/http`·`crypto/aes`·`crypto/cipher`만 쓴다.
- UI 디자인: 정본은 SAR-PUBLISH-001 '화면 기준'이다. `landing-style`·`landing-footer`·`memberStyle`을 재사용하고 새 색·글꼴·JS가 없다(`policy_pages.go`의 72ch 폭 규칙과 표 스크롤만 추가). 디자인 lint는 이 레포에 등록되지 않았다. 3폭 가로 넘침은 DEV 렌더 측정(`render/widths.json`)만 있고 테마 전환 대상은 없다(단일 테마). designer 직접 시각 수락(PC-06)과 tester 독립 QA는 아직 없다. 이것은 PLANS 후속에 있으며 코드 리뷰 수락의 전제가 아니다.

## 검토 범위

result.json의 68개 대상 파일을 모두 `reviewed`로 기록했다. 코드는 diff와 연결 코드를 직접 읽었다. 확인한 연결 코드는 다음과 같다.

- 기준 `internal/relay/calendar.go`의 `sealToken`·`openToken`·`revoke`. 새 `revokeSealed`의 nonce||ciphertext·AAD(member ID) 형식이 같다.
- `member.go` `startPage`. 정책 페이지의 세션 언어 처리와 같은 패턴이다.
- 템플릿의 `.Calendar`·`.CanDeny` 참조. 테스트 밖 Go·TS·YAML·Python 전체 grep에서 남은 참조가 0이다.

되돌림 범위도 검증했다. `git diff 67284e3 722feed9 -- internal cmd adapters/src .env.example compose.yaml`의 제품 코드 차이는 다음 세 가지뿐이다.

- legacy 필드 3줄
- `calendar_retire.go`·정책 페이지와 그 연결
- `main.go` 정리 호출

`git log 67284e3..05d0800e`의 11개 커밋이 모두 SAR-SCHEDULE-001 커밋임을 확인했다. 비캘린더로 보이는 `4a708c15`도 캘린더 테스트 파일만 바꾼다.

skipped 51개(`SAR-PUBLISH-001-DEV-test/`의 로그·PNG·JSON, jev packet-outcomes, handovers log)는 파일마다 열지 않았다. 대신 다음 방법으로 존재와 무결성을 확인했다.

- `report.md` 본문과 원격 로그 헤더·꼬리(`EXIT=0`, FAIL 0, `TestRetireCalendars`·`TestCalendarRoutesRetired` PASS)를 읽었다.
- `raw-log-sha256.txt`를 대조했다. git blob(LF)과 CRLF 사본의 sha256이 각각 기록된 두 값과 일치한다(예: lint `5c6f9d73…`/`323c9e96…`, test `8a9bc7e5…`/`07356ca0…`).

## 발견 사항

critical·high는 없다. 아래 low 4건은 모두 `resolved: false`로 남긴다. 수락을 막지 않는다.

| 심각도 | 파일·줄 | 재현 조건 | 영향 | 권고 |
|---|---|---|---|---|
| low | `internal/relay/calendar_retire.go:31-34` | 정리 transaction이 DB 오류나 상태 JSON 해석 오류로 실패한다 | 기동 실패 메시지가 "calendar cleanup failed"뿐이라 원인을 알 수 없다 | 토큰이 없는 `err`를 감싸 반환하거나 로그에 남긴다 |
| low | `internal/relay/calendar_retire.go:39-51` | 삭제 commit 뒤 revoke 루프·로그 줄 전에 프로세스가 종료된다(OOM·SIGKILL·배포 중단) | 토큰은 revoke 없이 지워지고 정리 줄이 없다. ops-guide의 "줄이 없으면 지울 연결이 없었다"는 해석이 이때 틀린다 | ops 절차에 첫 기동의 정상 진행 확인을 더한다. 이상 종료 시 테스트 사용자 전원에게 수동 접근 삭제 안내(SAR-SCHEDULE-001 처리 4)를 적용한다 |
| low | `internal/relay/store.go:497-499` | 캘린더 빌드에서 이미 받은 `relay.result` `done`을 같은 키로 재전송한다 | 기존 receipt 대신 403 `disclosure_denied`를 받는다. 메시지 24시간 보관이므로 배포 직후 24시간만 해당한다 | 67284e3 동작과 같아 수용한다. 기록만 남긴다 |
| low | `internal/relay/policy/privacy.ko.html:6-38`(EN 동일) | 운영에서 SMTP를 설정해 이메일 코드 로그인을 켠다 | 처리 위탁 표에 이메일 발송 업체가 없어 고지가 빠진다. P-04는 PLANS에서 ops 확인 대기인데 DEV가 SMTP 미설정 근거로 행을 뺐다 | ops의 P-04 확인과 'SMTP 활성화 전 문안 갱신' 조건을 PLANS 후속에 남긴다 |

중점 위험 확인 결과:

- 일회 정리: 키 유무·revoke 결과와 관계없이 토큰을 지운다. 삭제가 먼저 commit된다. 재기동 시 대상이 0이므로 요청도 0이다. 로그는 개수만 남기고 토큰·이메일·회원 ID는 남기지 않는다. 정리 실패 시 relay는 기동하지 않는다(fail closed). revoke는 토큰당 5초 timeout으로 순차 실행된다. 테스트 사용자 규모에서는 기동 지연이 무시할 만하다.
- 로그인 scope는 `openid email`만이다(`google.go:52`, integration 단언). `/home/calendar/*` 라우트가 없다. gate 조회는 `g.State`를 그대로 쓴다. gate-consume은 항상 `disclosure:false`다. `schedule.query`는 U5 기준으로 승인과 거절이 모두 가능하다. done은 `disclosure_denied`다. adapter도 항상 `denied`로 답한다.
- 정책 문안: KO/EN 4페이지가 SAR-PUBLISH-001 문안과 절 단위로 같다. P-01 시행일 2026-10-22, P-02·P-03은 표의 대체 문구, P-04는 위 low 1건이다. EN 본문의 한글 문자는 0이다(grep, `TestPolicyPages`의 `koreanLeft`). 수집 항목 근거 코드는 이 과제가 바꾸지 않았다. 캘린더 데이터 경로가 사라져 '캘린더를 수집하지 않는다'가 제거 뒤 코드와 맞다.
- ingress 후보: `^(...|/privacy|/terms|...)$`로 두 경로만 정확히 연다. 검사 스크립트가 `/privacy/`·`/privacy/x`·`/privacyx`·`/terms/`·`/terms.html` 차단을 단언한다.
- DB: SQL migration·스키마 삭제는 없다. `relay_state` JSON 한 행에서 다음 값이 지워진다.
  - `Calendars`(봉인 토큰·연결 상태)
  - 모든 `Gate.Lookup`(조회 시각·결과·미전달 disclosure)
  - calendar 표시 Google attempt

  이 삭제는 SAR-SCHEDULE-001 처리 2·3의 의도된 손실이다. consume된 gate의 조회 감사 기록도 함께 지워진다. 이것은 회수가 아니라 기록 삭제다. gate 자체도 24시간 뒤 sweep 대상이다. 정리 뒤 이전 빌드로 rollback하면 연결 0에서 재연결이 필요하다(transition.md에 기록).

## 검증 및 남은 제약

리뷰어는 lint·테스트를 다시 실행하지 않았다. 사용한 증거는 다음과 같다.

- `lint.json`(sha256 `1727b38c…`): coor가 깨끗한 임시 detached worktree에서 `722feed9`(base `05d0800e`)로 FullOps `lint.py`를 실행한 결과다.
  - 구현자 dev gate 디렉터리(`D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/`)에는 `722feed9`의 `lint-result.json`이 없었다. `lint-attempt.json`만 있었고, 구현자 gate 기록은 `722feed9`에서 blocked였다. 그래서 첫 `prepare`는 lint 없이 폴더를 만들었다.
  - coor의 첫 결과(`lint-722feed9.json`, sha256 `ae73fecf…`)는 실행 불가 2건의 `reason`이 비어 `check_lint`가 거부한다. coor가 같은 실행 결과에 reason만 보충한 `lint-722feed9-reasoned.json`을 주었다. 두 파일의 차이가 두 `reason` 값뿐임을 확인했다. 원본은 coor scratchpad에 남아 있다.
  - prepare는 기존 리뷰 폴더를 덮어쓰지 않는다. 그래서 보충 파일을 바이트 그대로 `lint.json`에 복사했다. prepare의 lint 처리는 `check_lint` 뒤 원본 바이트 기록뿐이고 `review.py check`가 같은 검사를 다시 수행한다.
- 실행 불가 2건(product-lint `make lint`, product-test `make test`) 판단: Windows 호스트에 make가 없다. DEV가 원격 Linux 격리 clone에서 제품 최종 SHA `e20c199`로 make lint·make test·integration(`-race -tags=integration`)을 실행했고 모두 EXIT=0이다. `git diff --name-only e20c199 722feed9`는 `.fullops-squad/` 기록 파일만 바꾼다. 그래서 같은 제품 코드의 증거로 수용한다. 로그 hash는 위 검토 범위에서 대조했다. 남은 제약: `722feed9` 자체에서 make를 직접 실행한 증거는 없다.
- WARNING 16건 판단:
  - SIZE-001 4건(`PLANS.md`, `ops-guide.md`, `i18n_en_member.go`, `member.go`): 각 파일은 1–13줄만 늘었다. 분할은 범위 밖이다.
  - ANTI-005 6건: 캘린더 테스트 5개 파일과 `human_pages_test.go` 1사례 삭제. 삭제된 기능의 테스트이며 정리 동작은 `calendar_retire_integration_test.go`가 대체한다. 검증 무력화가 아니다.
  - DOC-001 4건(`policy/*.html`): 방문자에게 그대로 보내는 정적 본문이다. 설명은 `policy_pages.go` 헤더에 있다. 수용한다.
  - SLOP-004 1건: 검사 결과 출력 줄이며 오탐이다.
  - SIZE-002 1건: 위 머리말의 근거로 수용한다.
- 미검증: 실제 Google revoke 끝점, 운영 상태의 실제 정리, live ingress 반영, 실브라우저 키보드·확대 검사, designer 시각 수락(PC-06), tester 독립 QA. 모두 PLANS 후속(coor·ops·tester·designer)에 있다.

## 대화 미참조 인계 점검

이전 대화 없이 `.fullops-squad/docs/deliverables/README.md`에서 시작해 snapshot에서 확인했다.

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `planning/product-specs/SAR-PUBLISH-001.md`, `SAR-SCHEDULE-001.md` '철회 — D018' | 확인 | P-02·P-05(user), P-03·P-04(ops) 열림 |
| 구조와 구현/미완료 상태 | `design-docs/architecture.md`·`module-design.md`·`data-model.md` 'SAR-PUBLISH-001-DEV' 절 | 확인 | legacy 필드 제거 조건이 data-model과 코드 `shortcut:` 주석에 있다 |
| 실행·검증 방법과 증거 | `evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md` | 확인 | 증거 SHA는 e20c199. 722feed9와 제품 코드가 같다 |
| 운영·복구(해당 시) | `operations/ops-guide.md` 18–20행, `operations/transition.md` | 확인 | 이상 종료 시 정리 줄 해석은 low 발견 |
| 다음 작업·담당·재개 조건 | `.fullops-squad/PLANS.md` SAR-PUBLISH-001-DEV 후속 2건 | 확인 | PLANS의 리뷰 대상 표기는 `e20c199`다. 실제 리뷰는 같은 제품 코드의 `722feed9`다 |
| 로컬 링크·절 접근/지원 한계 | `deliverables/README.md` → product-specs, SAR-SCHEDULE-001 '#철회--d018-2026-10-21' anchor | 확인 | anchor 해석은 GitHub 방식을 전제로 한다 |
| snapshot 정리 후 정본 접근 | 이 리뷰 폴더(`result.json`·`lint.json`·`report.md`), main의 정본 문서 | 확인 | snapshot 정리는 coor 몫 |

## 검토 결론

수락한다. 캘린더 제거는 캘린더 과제 전 판(67284e3)으로 깨끗하게 되돌렸다. 일회 정리는 SAR-SCHEDULE-001 처리 1–3·5·6을 지키고, 실패해도 토큰을 남기지 않으며, 재실행에도 안전하다. 정책 페이지는 정본 문안·화면 기준과 맞고, ingress 후보는 두 경로만 정확히 연다. lint ERROR는 0이다. low 4건은 후속으로 남긴다. `review.py check` 통과는 기록 검사이며 이 AI 검토의 판단이나 테스트 성공을 자동으로 보증하지 않는다.
