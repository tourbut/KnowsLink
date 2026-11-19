---
title: SAR-UX-ISSUES-002-LANG-CODE-REVIEW 리뷰
status: draft
updated: 2026-10-20
owner: reviewer
summary: 7704e6d31288 병합 전 독립 리뷰 보고서
---

# SAR-UX-ISSUES-002-LANG-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: reviewer_session `lang002-code-review-opus` / Claude Code / `claude-opus-5-5`. 실제 host 세션 ID는 `93aee436-47d5-4be1-9995-709e1628359a`다. 구현자 세션은 `ctx_9fec89dfb1ec`이며 이 세션과 다르다.
- base SHA / head SHA / merge-base: `131544873193795e77a2886dcdd81e8a2658c0f6` / `7704e6d3128831468bb474330210bc2a5099bce7` / `131544873193795e77a2886dcdd81e8a2658c0f6`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-c18675c1cf8f476fb89c9ef2245dec75` (detached, HEAD 7704e6d, `git status` clean). snapshot은 읽기만 했고 결과는 coor 기록 폴더에만 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.2 (55450a9), `.fullops-squad/review/rule.json`(rule_sha256 `f1061481…`). 공통 규칙 `fullops-common-0.3.3`(README·coding-style·testing·security), `.fullops-squad/project.md`, 문서 규칙 `docs/agents/document-writing.md`. 범위 안에서 규칙 파일 변경은 없다.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U3(설정 순서·선호 우선순위 표·비로그인 진입·번역 범위·lang/KST·누락 fallback·U3-A)와 U7 '설정'. 지시서와 완료 보고는 `.fullops-squad/handovers/logs/2026-10-20_to_dev.md` 723행 이후의 SAR-UX-ISSUES-002-LANG-DEV 절이다. DEV 보고는 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/report.md`다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 101 / 39 / 62 / 31 / 70. 지시서의 소유 범위(`internal/relay/`·기술 문서·own 기록)와 실제 diff가 일치한다. `PLANS.md`·`board.json` 변경은 coor 준비 커밋 `a3a8f443`이며 DEV 커밋이 아니다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0, WARNING 7(SIZE-001 6, SIZE-002 1), 실행 불가 2(product-lint·product-test). 사유는 Windows 호스트에 make가 없어서다. coor가 같은 SHA 7704e6d를 격리 Linux clone에서 실행한 로그로 보완했다(아래 '검증' 참고).
- SIZE-002: 지시서는 "예상 규모가 크다(전 템플릿)"로 예고했다. 실제 추가는 1931줄이다. 번역 표 2개 파일(약 750줄)과 검사·fixture가 대부분이다. 표는 키 하나에 한 줄이므로 분할해도 줄 수가 줄지 않는다는 DEV 설명을 수용한다.
- SIZE-001: `human_http.go`·`member.go`·`remote_auth.go`는 언어 저장·render·세션 헤더 몇 줄만 늘었다. `i18n_en_member.go`·`i18n_test.go`는 표와 검사 파일의 특성이다. `PLANS.md`는 이 과제 변경이 아니다. 모두 수용한다.
- DEP-001: 해당 없음. `go.mod`·`go.sum`·`adapters/package.json`·`Makefile`은 바뀌지 않았다. 번역은 표준 라이브러리와 기존 `html/template`만 쓴다.
- UI 디자인: 지시서의 정본은 D04와 SAR-UX-ISSUES-002 '화면 배치와 자기 비평'이다. 새 테마·폰트·JS·색 토큰은 없다(diff 확인). 테마 전환은 해당 없음(테마가 없다). DEV가 두 언어 17화면을 1440·390·320에서 렌더해 가로 넘침 0을 보고했다(`render/probe.json` 51건 모두 overflow false). designer 시각 수락은 별도다.

## 검토 범위

result.json의 모든 (path, status)에 검토 상태를 기록했다.

- 제품 코드 18개 파일과 기술 문서 4개를 diff와 연결 코드로 읽었다. 연결 코드는 `capacity.go`의 `boundedHTTP`·`capacityResponse`, `receipt_context.go`의 `receiptFailure`, `gate_summary.go`, `protocol.go`의 timezone 검증, `http.go`의 `crossOrigin` 보호다.
- 지시서와 완료 보고(handover 로그)는 직접 읽고 reviewed로 기록했다.
- 번역 표 2개(754줄)는 표본으로 확인했다. 설정·언어·친구·에이전트·철회·연결 안내·공개 첫페이지·로그인·세션·오류·패턴 키에서 약 110항목을 읽었다. 값에 태그나 속성 따옴표가 새로 들어가지 않았는지 전체를 grep으로 확인했다. 전체 문장 교정은 하지 않았다.
- skipped: `docs/evaluations/jev/` 자동 생성 기록 7개, DEV 원시 로그 10개, 해시 목록 1개, 렌더 PNG 51개, `probe.json` 1개다. 원시 로그는 `raw-log-sha256.txt`로 `sha256sum -c`를 실행해 10개 모두 OK였다. PNG는 `probe.json`과 DEV 보고로 존재만 확인했다.

## 발견 사항

### F1 — medium, 미해결: admission 거부 화면이 항상 한국어이고 DEV 보고 서술과 다르다

- 위치: `internal/relay/http.go:212-218`(`browserLang` 호출 위치), 연결 코드 `internal/relay/capacity.go`의 `boundedHTTP`·`capacityResponse`.
- 재현 조건: 쿠키 `__Host-kl_lang=en`인 비로그인 브라우저, 또는 저장값이 `en`인 회원이 동시 처리 한도나 요청 budget을 넘는다. 이때 `capacityResponse`('요청 대기' 429), admission의 '요청 제한' 429, `/home/receipts`의 `receiptFailure`(429/503)가 화면을 그린다.
- 원인: `browserLang`은 `bounded`의 내부 핸들러 안에서만 실행된다. admission 거부는 그 전에 `render`를 부른다. `Content-Language`가 비어 있으므로 `responseLang`은 항상 `ko`다.
- 영향: 영어를 고른 사용자가 혼잡 시 한국어 거부 화면을 본다. `lang="ko"`와 내용은 일치하므로 깨진 표시는 아니다. DEV 보고 59행과 완료 보고의 "브라우저 쿠키 언어를 따른다(회원 저장값 미적용)"는 코드와 다르다.
- 수정 요청: `browserLang`을 `boundedHTTP` 앞(root 또는 bounded 진입부)으로 옮긴다. 또는 보고를 "항상 한국어"로 정정하고 알려진 한계로 남긴다.

### F2 — medium, 미해결: 영어 화면의 한국어 외부 문서 링크에 언어 표시가 없다

- 위치: `internal/relay/i18n_en_human.go:135-153`, `internal/relay/i18n_en_member.go:254-303`.
- 재현 조건: 영어로 공개 첫페이지나 내 에이전트의 연결 안내를 연다.
- 원인: README·adapters/README 링크 문구가 'Grok new installation guide', 'Install and usage guide'처럼 번역된다. 한국어 문서라는 표시나 `hreflang="ko"`가 없다. `README.md`는 한국어 문서다.
- 영향: U3의 "외부 한국어 문서 링크는 언어를 표시한다"를 지키지 않는다.
- 수정 요청: 영어 값에 '(Korean)' 같은 보이는 표시를 붙인다. 또는 템플릿 링크에 `hreflang="ko"`와 보이는 표시를 함께 둔다.

### F3 — low, 미해결: 공개 첫페이지의 언어 링크가 로그인 회원에게 반응하지 않는다

- 위치: `internal/relay/landing_pages.go:13`.
- 재현 조건: 저장값이 `ko`인 회원이 로그인 상태로 `/`에서 English를 누른다.
- 영향: 브라우저 쿠키만 `en`이 되고 회원 저장값이 우선하므로 화면은 바뀌지 않는다. 안내도 없다. 로그아웃 뒤 그 쿠키가 같은 브라우저의 비로그인 화면과 저장값 없는 회원 화면에 적용된다. 명시적 브라우저 선택이므로 계정 선호 복사는 아니다.
- 수정 요청: `Authenticated`이면 lang-nav 대신 설정 링크를 둔다. 또는 회원 저장값이 우선한다고 안내한다.

### 문제를 찾지 못한 중점 항목

- POST `/home/language`: 기존 `humanAction` 경로를 그대로 쓴다. `memberHit`가 세션과 rate를 확인하고, 익명은 `invalid_auth`로 `expired` 처리된다. CSRF는 `http.CrossOriginProtection`(Sec-Fetch-Site/Origin)과 `SameSite=Strict` 세션 쿠키로 다른 `/home` POST와 같은 수준이다. 본문은 `MaxBytesReader` 16 KiB이고 값은 `validLang`(ko/en) 화이트리스트다. 저장은 `Member.Locale`만 바꾸며 쿠키를 쓰지 않는다.
- GET `/lang/{code}`: 알 수 없는 code는 404이고 쿠키를 쓰지 않는다. `safeBack`은 스킴·호스트·`//`·백슬래시·CR/LF·상대 경로를 거부한다. 제어 문자는 `url.Parse`가 거부한다. 이후 `http.Redirect`가 `path.Clean`을 적용하므로 `/.//evil`도 `/evil`이 된다. `/%2F%2Fevil`은 percent 그대로 경로로 남는다. 열린 리다이렉트를 찾지 못했다. 쿠키는 `__Host-` 접두사, Secure, HttpOnly, SameSite=Strict, Path=/, Domain 없음, 365일이다.
- 상태 변경 GET: `/lang/{code}`는 cross-site 링크로 비로그인 언어 쿠키를 바꿀 수 있다. 영향은 표시 언어뿐이다. 회원 저장값이 우선하고 다른 상태는 바뀌지 않는다. 승인된 설계(쿠키 1개, 링크 진입)의 범위로 수용한다.
- 계정 교대: 회원 저장값은 쿠키로 복사되지 않는다. 쿠키는 `/lang/` 링크에서만 쓴다. 통합 검사가 B의 A 선호 미상속과 명시 저장 뒤 쿠키 불변을 단언한다.
- 손상 locale: `validLang`이 아니면 다음 출처로 넘어간다. 값은 화면에 나오지 않는다. 우선순위 10조합 단위 검사와 통합 검사(`xx`)가 있다.
- 템플릿 escape·구조: `translateTemplate`는 템플릿 소스의 텍스트 구간만 표 값으로 바꾼다. 그 뒤 `html/template`로 다시 parse하므로 모든 action 출력은 문맥 escape를 유지한다. 표는 정적 Go 상수이고 외부 입력이 아니다. 값의 action 집합이 키와 같은지 검사가 확인한다. 값에 새 태그나 속성 따옴표가 없음을 grep으로 확인했다. 큰따옴표가 있는 값 1건은 본문 텍스트이고, `&lt;` 2건은 키와 같다. `tr`은 문자열을 돌려주므로 escape된다. `badge`는 `HTMLEscapeString`을 쓴다. lang-nav의 `back` 값은 URL query 문맥에서 escape된다.
- 원문·ID·서명값·KST: 메시지 본문·ID·Fingerprint·Kid·명령은 `tr`을 거치지 않는다. `tr`이 적용되는 값은 서버 상태 라벨, client 라벨, gate 요약 라벨이다. gate의 timezone은 `time.LoadLocation`으로 검증된 ASCII라 번역 대상이 아니다. `clock`은 두 언어 모두 `2006-01-02 15:04:05 KST`로 같다. `humanClock`은 같은 KST instant에서 단어만 바꾼다(단위 검사 3건).
- 공개 첫페이지: `<html lang>`과 헤더 lang-nav 하나만 추가됐다. 섹션·설치 사실·명령·링크 구조는 같다. 390px에서 링크가 다음 줄로 내려가는 것은 DEV가 보고했다.
- fixture: `lang_browser_fixture_test.go`는 `//go:build integration && langfixture`인 `_test.go`다. 제품 바이너리에 들어가지 않는다. 환경 변수가 없으면 실행하지 않고 manifest는 0600이다.

## 검증 및 남은 제약

- 리뷰어는 테스트·lint를 다시 실행하지 않았다. 기존 증거만 읽었다.
- coor 보완 로그: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev/`에 있다. `.base`는 `/tmp/knowslink-lang002-7704e6d31288-…`이고 install 로그는 `7704e6d3128831468bb474330210bc2a5099bce7` detached checkout을 기록한다.
  - `7704e6d…-lint.log`: `go vet`·`go mod verify`·adapters check·prettier·compose·public ingress·public config 검사 후 `EXIT=0`이다. 로그의 `ERROR: candidate exists`·`invalid tunnel UUID` 두 줄은 `test_public_config.sh`가 의도한 실패 경계의 출력이다. 같은 검사가 "validation failures propagated"로 통과를 보고한다.
  - `7704e6d…-test.log`: `go test -race ./cmd/... ./internal/...` 전체 ok, adapters 테스트 PASS, `test_shell_guard.py` 4건 OK, `EXIT=0`이다.
  - `7704e6d…-integration.log`: `internal/relay` ok(137.289s)이고 `EXIT=0`이다. `TestLanguagePreferencePostgres`·`TestLanguageChoiceOrderAndLink`·`TestHumanClockKeepsInstantAcrossLanguages`가 PASS다.
- DEV 원시 로그 10개는 `raw-log-sha256.txt`와 해시가 같다. 제품 코드는 `2ff9e44` 이후 바뀌지 않았다. `2ff9e44..7704e6d`의 코드 변경은 `lang_integration_test.go`뿐이다.
- 테스트 감소·기대값 완화: `member_ui_test.go`·`connect_guide_test.go`·`gate_summary_integration_test.go`는 템플릿 변수 통합과 `badge` 서명 변경만 반영했다. 단언은 줄지 않았다.
- 미실행 범위(DEV not_run 그대로): 저장 실패 503 자동 검사, 실제 키보드·200%·스크린리더·보조 브라우저, fixture 브라우저 왕복, `TestRemoteConsentBrowserReturn`, 실제 브라우저의 쿠키 차단, 영어 전체 교정. tester·designer가 확인할 범위다. admission 거부 화면 언어(F1)는 자동 검사가 없다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `docs/design-docs/SAR-UX-ISSUES-002.md` U3·U7, 지시서 '적용 기준과 예외' | 확인 | 외부 한국어 문서 표시 요구가 미구현이다(F2). |
| 구조와 구현/미완료 상태 | `docs/design-docs/module-design.md` '설정·언어 모듈 (5단위)', `interface-design.md` '설정·언어 (5단위)' | 확인 | 코드와 일치한다. admission 한계 서술은 없다. |
| 실행·검증 방법과 증거 | DEV `report.md`, `remote.py`, `raw-log-sha256.txt`, coor 보완 로그 경로 | 확인 | DEV 보고 59행 서술이 코드와 다르다(F1). |
| 운영·복구(해당 시) | `docs/operations/transition.md` 'DEV 후보 인계' | 확인 | migration 없음, 운영 미반영이다. `Member.Locale`은 생략 가능 필드라 되돌려도 상태 JSON이 깨지지 않는다. |
| 다음 작업·담당·재개 조건 | 완료 보고 '남은 일 (coor)', `PLANS.md` 5단위 항목 | 확인 | tester 독립 QA, designer 시각 검수, main 통합. 이 리뷰의 F1–F3 후속을 PLANS에 더해야 한다. |
| 로컬 링크·절 접근/지원 한계 | 완료 보고의 DEV 보고·packet-outcomes·remote.py 상대 링크 | 확인 | 링크 대상 파일이 snapshot에 있다. README 앵커 percent-encoding은 브라우저 클릭으로 확인하지 않았다. |
| snapshot 정리 후 정본 접근 | coor 기록 폴더 `SAR-UX-ISSUES-002-LANG-CODE-REVIEW-review/` | 확인 | 결과·보고·lint는 snapshot 밖에 있다. coor 보완 로그는 레포 밖 private 도구 디렉터리에 있어 정리되면 사라질 수 있다. |

## 검토 결론

조건부 수락한다. critical/high 발견 사항은 없다. 중점 보안 경계(CSRF·인증·입력 검증·열린 리다이렉트·쿠키 속성·계정 교대 미복사·손상 locale fallback·템플릿 escape)에서 결함을 찾지 못했다. 원문·ID·KST·공개 첫페이지 구조·fixture 격리도 지시서와 정본을 지킨다.

medium 2건(F1, F2)과 low 1건(F3)은 미해결이다. 병합 전 수정하거나 PLANS.md에 후속 과제로 남기는 판단은 coor가 한다. 이 결론은 AI 검토 판단이며 OCR 자동 판정이 아니다. `review.py check` 통과는 기록 형식 검사이며 내용과 테스트 성공을 보증하지 않는다.
