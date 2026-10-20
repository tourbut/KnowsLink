---
title: SAR-UX-ISSUES-002-LANG-DEV DEV 보고
status: review
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-LANG-DEV]
summary: "설정 화면·한국어/English 저장(U3), 헤더 ‘설정’(U7), 지정 화면 영어판의 기술 계획·누락 번역 검사·U3-A 검사·렌더·fixture 기동 근거와 미실행 항목"
---

# SAR-UX-ISSUES-002-LANG-DEV DEV 보고

DEV 검사 기록이다. 독립 리뷰·tester QA·designer 시각 수락이 아니다.

## 기술 계획 (실행함)

1. 문구 관리: 한국어 원문이 키다. 템플릿 본문은 바꾸지 않고 `i18n.go:translateTemplate`가 템플릿 소스를 구간(태그 사이 본문·속성값, 제어 action에서 분리, `{{.X}}` 같은 값 action은 구간 안에 유지)으로 나누어 영어 표로 바꾼 뒤 영어 집합 `memberPagesEN`을 한 번 더 파싱한다. 핸들러 문구는 `render`가 Title·Problem·Notice에, 템플릿 함수 `tr`·`name`·`badge`가 상태·안내·서버 fallback 이름·배지에 `tr(lang, 원문)`을 적용한다. 합성 문구는 `%s`(다시 번역)·`%v`(그대로) 패턴 키로 처리한다. 영어 표는 `i18n_en_member.go`(회원·공개·remote 화면과 그 핸들러)·`i18n_en_human.go`(사람 화면과 그 핸들러)로 파일 소유를 나눠 하위 에이전트 2개가 채우고 부모가 검사로 취합했다. 새 의존성은 없다.
2. 회원 선호 저장: `Member.Locale`(상태 JSON 생략 가능 필드, `ko`/`en`). POST `/home/language`만 바꾼다. 브라우저 쿠키에는 쓰지 않는다.
3. 비로그인 선호 쿠키: `__Host-kl_lang`, 값 `ko`/`en`, Secure·HttpOnly·SameSite=Strict·Path=/(기존 `setCookie` 속성 그대로), 수명 365일. GET `/lang/{code}?back=<경로>`가 쓰고 내부 경로로만 303 복귀한다.
4. 요청별 언어 결정: 회원 저장값 → 쿠키 → 한 응답용 `?lang=` → 한국어(`pageLang`). `browserLang` 미들웨어가 `Content-Language` 헤더에 기록하고 `memberHit`가 회원을 확인하면 회원 값으로 바꾼다. `render`·거부 화면이 그 헤더로 템플릿 집합을 고른다. 손상된 값은 다음 출처로 넘어간다. `?lang=`이 쿠키 없이 적용되면 ‘이 화면에만 적용됩니다’ 안내를 붙인다.
5. 누락 번역 감지: (a) `TestTranslationCoversTemplates` 템플릿 소스의 한글 구간 중 표에 없는 키 0, 값의 `{{…}}` action 동일성, 값에 한글 없음; (b) `TestTranslationCoversGoText` 제품 Go 파일의 한글 문자열 리터럴(메일·humanClock 단어·템플릿 소스 제외)이 키이거나 패턴 키의 일부; (c) `TestEnglishScreensHaveNoKorean` 지정 화면 전부를 ASCII 사용자 데이터로 영어 렌더해 한글 0(`한국어` 언어명만 예외); (d) `TestComposedHandlerTextTranslates` 모든 `agentResume` 상태와 합성 문구의 패턴 번역. 런타임 누락은 한국어 원문으로 남고 내부 키·빈 버튼은 생기지 않는다.

## 변경 요약

- U7: 헤더 ‘내 계정’→‘설정’(회원 head·human-head), 친구 화면 링크, 설정 제목·title.
- U3 설정: 내 이름 → 내 친구 코드 → 언어(radio 한국어/English, ‘언어 저장 / Save language’, 선택만으로 제출 없음) → 로그인과 보안. 저장 성공은 `303 /home/account?n=language#language`와 status 안내, 지원하지 않는 값은 422 alert와 저장값 유지, 저장 실패(503)는 선택한 언어의 거부 화면과 `#language` 복귀.
- 비로그인 진입: 공개 첫페이지 헤더·확인 코드·device 동의(로그인 전)·remote 동의 화면에 ‘한국어 / English’ 링크(`aria-current`로 현재 표시). 공개 페이지 정보 구조·설치 사실·토큰 불변.
- 번역 범위: 공개/로그인/이름/대화/친구/내 에이전트/설정/legacy/receipt/승인/local 연결/remote 동의·복귀의 제목·메뉴·라벨·상태·빈 상태·오류·검증·접근성 이름·도움말. `<html lang>`과 `Content-Language`가 일치. 정확한 기한은 KST 문자열 그대로, 짧은 대화 시각만 단어 변경(`humanClock(lang)`).
- 번역하지 않음: 프로토콜 ID·키·명령·사용자 이름·메시지·검증 원문 JSON, 확인 메일, README, JSON API, `/owner` 작업 화면(지정 화면 아님), 외부 호스트 UI.

## 검사

제품 코드는 `2ff9e44`에서 확정됐다(`9028971` 첫 구현 → `824b4dd` 주석 → `2ff9e44` 번역표 gofmt). 이후 `db94002`·`326b82c`·`cc016b5`는 `lang_integration_test.go`만 바꿨다(`git diff --stat db94002 cc016b5`: 그 파일 1개). 원격 Linux 격리 clone(`remote.py`, 자기 Postgres 컨테이너)에서 실행했고 원시 로그는 이 폴더에 있다([raw-log-sha256.txt](raw-log-sha256.txt)).

| 검사 | SHA | 결과 | 근거 |
|---|---|---|---|
| 원격 `make lint` | db94002 | EXIT=0 | [lint](db94002bd88561224d07531296df66f625059df1-lint.log). `824b4dd`에서는 번역표 두 파일의 gofmt 위반으로 EXIT=2였다(Windows PATH의 옛 gofmt가 통과시킨 정렬 차이). 모듈 toolchain gofmt로 고쳐 `2ff9e44`부터 통과 |
| 원격 `make test` | db94002 | EXIT=0 | [test](db94002bd88561224d07531296df66f625059df1-test.log) |
| 원격 integration 전체(`-race`, `-skip TestRemoteConsentBrowserReturn`) | db94002 | EXIT=1, 실패 1건 = 새 `TestLanguagePreferencePostgres`뿐(116 PASS) | [integration](db94002bd88561224d07531296df66f625059df1-integration.log). 실패 원인은 테스트 자체(저장 뒤 `?n=language` 없이 조회, 같은 주소 재로그인의 메일 60초 간격 제한, 링크 markup 순서 기대)였고 제품 결함이 아니다. `2ff9e44` 전체 실행도 같은 1건만 실패([로그](2ff9e44a38b2cc6cc5543617c74f5f94f74f37e9-integration.log)) |
| 원격 `TestLanguagePreferencePostgres`(수정 뒤) | cc016b5 | EXIT=0 PASS | [로그](cc016b5a7ec68baee29ba77c77e82ae1f950f58d-integration-TestLanguagePreferencePostgres.log); 직전 `326b82c`는 markup 순서 기대로 FAIL([로그](326b82c892a8fd3c9b3f5fb0b4f8e30e6c7ec1e8-integration-TestLanguagePreferencePostgres.log)) |
| 로컬 `go test ./internal/relay`(unit, 번역 커버리지 포함) | cc016b5 | ok | 아래 자동 검사 |
| fixture 컴파일·실제 기동·probe | cc016b5 | EXIT=0 | fixture 절 |
| FullOps lint.py `--from 1315448` | 2ff9e44(깨끗한 detached worktree) · 최종 HEAD(dev, 사용자 미추적 폴더를 그 프로세스에만 `core.excludesFile`로 제외) | ERROR 0 | 완료 보고의 WARNING 목록 |

U3-A 자동 검사: `TestLanguagePreferencePostgres`(Postgres)는 기본 한국어와 설정 순서, English 저장 → 같은 설정 위치 영어·status 안내·radio 유지·`Content-Language`·쿠키 미기록, 다른 화면과 거부(422) 영어, 다른 브라우저 같은 계정 영어, 로그아웃 뒤 공개 페이지 한국어, 같은 브라우저 계정 교대(B에 A 선호 미복사), 비로그인 `/lang/en` 쿠키·복귀·다음 방문 영어, 비로그인 선호를 가진 B 첫 로그인은 화면 적용만(저장값 없음), B 명시 저장 ko가 쿠키 en보다 우선·로그아웃 뒤 쿠키 복귀, 잘못된 값 422·저장값 유지, 손상 locale `xx` 한국어 fallback·비노출, 익명 저장 차단(`303 /?n=expired`), 쿠키 없는 `?lang=` 한 응답 적용과 안내·마커 비누출, 외부 `back` 거부, 언어 동작 전후 Friends/Agents/Gates/Messages/Connections/Pairs 수 동일, 같은 기한 문자열이 두 언어에 동일함을 단언한다. 단위: `TestLanguageChoiceOrderAndLink`(우선순위 표 10조합, `safeBack`, 쿠키 속성·수명, 404, 한 응답 안내), `TestHumanClockKeepsInstantAcrossLanguages`, `TestComposedHandlerTextTranslates`, `TestTranslationCoversTemplates`·`TestTranslationCoversGoText`·`TestEnglishScreensHaveNoKorean`(지정 화면 영어 렌더 한글 0). 문구 복제 테스트는 만들지 않았다.

## 렌더

`TestEnglishScreensHaveNoKorean`(영어, `en-*.html`)·`TestHumanPagesRenderAndEscape`·`TestMemberUIFixtures`가 `KNOWSLINK_UI_FIXTURE_DIR`로 만든 합성 HTML을 headless Chrome(비모바일 viewport) 1440·390·320에서 [render/probe.mjs](render/probe.mjs)로 렌더했다. 대상 17개(영어 15: 설정·공개 첫페이지·확인 코드·친구·내 에이전트(Claude 선택)·대화·legacy·조회·승인·연결 수단·device 로그인·remote 동의·복귀·거부·Google 완료, 한국어 2: 설정·공개 첫페이지), 51장 모두 페이지 가로 넘침 0([probe.json](render/probe.json)). 390px 공개 첫페이지에서 ‘한국어 / English’ 링크는 헤더 다음 줄로 내려간다. DEV 확인이며 designer 시각 수락이 아니다.

## 하위 에이전트

Purpose: implementation 2건(standard 예산 2, lite). 파일 소유를 `i18n_en_member.go`(회원·공개·remote 화면과 그 핸들러)와 `i18n_en_human.go`(사람 화면과 그 핸들러)로 나누고 같은 용어집·패턴 키 목록을 전달했다. 부모가 세 커버리지 검사로 취합했고, 하위가 지적한 템플릿 결함 2건(`{{.Client}}`·`{{.Scope}}` 미번역 출력)과 배지 라벨의 버튼 문구 충돌(‘거절’=Decline/Denied)을 부모가 고쳤다(배지는 `badges`의 영어 열로 분리). README 앵커의 한글 fragment는 값에 한글을 둘 수 없어 percent-encoding했다(같은 대상). 독립 코드 리뷰는 수행하지 않았으며 coor의 고정 SHA 리뷰 대상이다.

## 못 한 것 (not_run)

- 저장 실패(503) 경로의 자동 검사: DB 장애 주입 수단이 없어 코드 경로만 둔다(선택 언어의 거부 화면·`#language` 복귀).
- 실제 키보드·200% 확대·스크린리더·보조 브라우저(U7-A 재검증), fixture 브라우저 왕복, `TestRemoteConsentBrowserReturn`(기존 skip 조건 그대로).
- 쿠키를 막은 실제 브라우저: `?lang=` 한 응답 적용과 안내는 HTTP 수준 검사(`TestLanguageChoiceOrderAndLink`·통합 검사)로만 확인했다.
- 공용 admission의 429/503 거부 화면은 회원 조회 전이라 브라우저 쿠키 언어를 따른다(회원 저장값 미적용). 알려진 한계로 둔다.
- 영어 문장의 사람 교정: 하위 에이전트 번역을 부모가 핵심 화면·설정·헤더·오류만 표본 확인했다. 전체 문체 검수는 designer/tester 몫이다.

## 산출물

- 제품: 새 `internal/relay/i18n.go`·`i18n_en_member.go`·`i18n_en_human.go`, 수정 `member_pages.go`·`human_pages.go`·`landing_pages.go`·`remote_auth.go`·`member.go`·`human_http.go`·`http.go`·`identity.go`·`device.go`·`google.go`·`member_agents.go`·`member_receipt.go`.
- 검사: 새 `i18n_test.go`·`lang_integration_test.go`·`lang_browser_fixture_test.go`, 수정 `member_ui_test.go`·`connect_guide_test.go`·`gate_summary_integration_test.go`.
- 문서: D05·D10·D11·D13, [packet-outcomes](../../jev/SAR-UX-ISSUES-002-LANG-DEV-packet-outcomes.json), 이 보고와 원시 로그([remote.py](remote.py)).
