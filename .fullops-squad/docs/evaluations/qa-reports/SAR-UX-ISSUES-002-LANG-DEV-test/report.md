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

(아래 표는 최종 SHA 확정 뒤 채운다)

## 렌더

## fixture

## 하위 에이전트

## 못 한 것 (not_run)

## 산출물
