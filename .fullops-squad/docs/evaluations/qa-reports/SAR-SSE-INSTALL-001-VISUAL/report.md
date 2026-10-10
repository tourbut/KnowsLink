---
title: 원격 연결 동의 화면의 독립 직접 시각 검수
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-VISUAL]
summary: 고정 febd의 Google 확인과 설치 승인·취소 화면을 desktop1280/mobile390에서 직접 판정한다.
---

# SAR-SSE-INSTALL-001-VISUAL — 직접 시각 검수

## 판정과 범위

고정 `febd1282df2e84a29fa0a53834569e7b9cc028e2`의 합성 시각 검수는 **PASS**다. 최초 Google 계정 확인과 설치 승인·취소 선택 화면을 각각 1280×800 및 390×844에서 직접 확인했다. 일반 반환 URI와 긴 반환 URI를 나누어 총 8개 조건을 관찰했다. 긴 URI의 모바일 하단 2개 캡처를 보완했다. PNG 12개를 모두 직접 열어 판정했다.

제품 수락·실제 Google 동의·실제 호스트 설치·승인 저장·취소 처리는 이 PASS에 포함하지 않는다. 해당 동작 QA와 보안 리뷰는 별도 tester 과제다. 제품 코드와 기획·디자인 원천은 변경하지 않았다. 시각 finding은 0건이다.

## 독립 세션과 원천

| 항목 | 실제 값 |
|---|---|
| 검수 세션 | `01a126c0-2b8d-7332-8de2-e6ca9fb92f95` — 환경의 `CODEX_SESSION_ID`와 `CODEX_THREAD_ID`에서 확인 |
| 구현자 세션 | `01a12678-994d-7680-a15c-83577943b52d` — 지시서의 DEV 구현자와 다름 |
| 실제 terminal | `term_11fc1118-0e42-46ab-a7bb-77dccaa19315` |
| Task / Dispatch | `task_415d9cd925c6` / `ctx_4283f224b494` |
| 기록 체크아웃 시작 HEAD | `cdcc2a9b858fdec73fe487e06c9f874ef3b560b6` |
| UI 고정 source SHA | `febd1282df2e84a29fa0a53834569e7b9cc028e2` |
| `remote_auth.go` Git blob | `a142feb88cabf6892e6dd90d6f914ff5c1081231` |
| `member.go` Git blob | `dd44c5bdf1cc2d0c3144d76e2769e27f608b4e3d` |
| Orca embedded page ID | `c1de83f5-2c16-4f6c-ad74-0e508af70166` |
| Browser runtime ID | `7d495f42-7b90-4591-8523-d39e444844ff` |

`git show <固定SHA>:internal/relay/remote_auth.go`의 `remotePage`와 `member.go:17`의 `memberStyle`을 읽었다. [source-manifest.json](source-manifest.json)에 원천·CSS·fixture의 SHA256을 남겼다. `render_check.go`는 같은 고정 Git blob을 Go `html/template`로 실제 실행한다. 저장된 HTML 4개와 출력 바이트가 모두 같았다. 현재 기록 HEAD와 고정 후보 사이 제품 diff도 0이다.

`build_fixtures.py`는 합성 값과 기존 두 분기만 펼친다. 새 CSS·문구·DOM 구조는 추가하지 않는다. 클라이언트 이름은 `Synthetic Work client`다. 반환 URI는 예약 도메인 `client.example`의 37자 또는 545자 합성 주소다. 숨은 ID와 CSRF 값도 사용 불가능한 합성 값이다. 실제 서버·계정·DB·키를 사용하지 않았다.

기준은 `fullops-common-0.3.3`의 코딩·테스트·보안, `project.md`, 문서 작성 규칙, 역할 책임 규약이다. [D04](../../../design-docs/mockups/SAR-MVP-001-UI.md)의 기존 검수 경계와 `memberStyle`을 적용했다. 별도 디자인·기획 승인은 수행하지 않았다. 하위 위임은 짧은 단일 시각 검수여서 사용하지 않았다. 테스트 레벨은 lite다.

## 직접 관찰

| 조건 | 캡처 | 관찰과 판정 |
|---|---|---|
| desktop 최초 확인 | [desktop-login.png](desktop-login.png) | 제목·클라이언트·반환 주소·서버 키 보관 안내를 읽을 수 있다. `Google 계정 확인`과 관리 링크가 구분된다. PASS. |
| desktop 승인·취소 | [desktop-consent.png](desktop-consent.png) | `이 설치 연결 승인`과 `취소`가 한 줄에 표시된다. 글자와 버튼 경계가 잘리지 않는다. PASS. |
| desktop 긴 URI 최초 확인 | [desktop-login-long.png](desktop-login-long.png) | 545자 URI가 내용 폭 안에서 줄바꿈된다. 안내와 확인 버튼이 800px 높이 안에 보인다. PASS. |
| desktop 긴 URI 승인·취소 | [desktop-consent-long.png](desktop-consent-long.png) | URI가 버튼을 가리지 않는다. 승인·취소 두 버튼과 관리 링크가 보인다. PASS. |
| mobile 최초 확인 | [mobile-login.png](mobile-login.png) | 제목·본문 18px·확인 버튼을 읽을 수 있다. 20px 좌우 여백과 URI 줄바꿈을 유지한다. PASS. |
| mobile 승인·취소 | [mobile-consent.png](mobile-consent.png) | 두 버튼이 같은 줄에 표시된다. 두 레이블이 명확하며 가로 잘림이 없다. PASS. |
| mobile 긴 URI 최초 확인 | [첫 화면](mobile-login-long.png), [하단](mobile-login-long-bottom.png), [전체 보조](mobile-login-long-full.png) | URI와 안내가 세로로 길어진다. 하단으로 스크롤하면 확인 버튼과 관리 링크를 온전히 읽을 수 있다. PASS. |
| mobile 긴 URI 승인·취소 | [첫 화면](mobile-consent-long.png), [하단](mobile-consent-long-bottom.png), [전체 보조](mobile-consent-long-full.png) | 세로 스크롤 뒤 승인·취소 버튼이 모두 보인다. 가로 이동은 필요 없다. PASS. |

본문은 `18px system-ui`, 전경은 `rgb(23, 35, 43)`, 배경은 `rgb(245, 247, 248)`이다. 고정 원천의 기존 CSS와 같다. 제목·문단·기본 버튼·파란 관리 링크의 기존 형태를 유지한다. 상태나 선택을 색만으로 구분하지 않는다. 승인·취소 버튼은 같은 기본 스타일이며 문구와 개별 테두리로 구분된다. 두 버튼 사이 여백은 0px인 기존 template 출력이다. 문구가 겹치거나 버튼 영역이 잘리지 않아 새 디자인 변경 조건으로 올리지 않았다.

모든 버튼 높이는 52px다. 승인 버튼 폭은 약 172.98px, 취소는 64px다. 모바일에서 두 버튼의 오른쪽 끝은 약 256.98px이며 내용 폭 안에 있다. 긴 URI 모바일의 버튼 위쪽은 문서 기준 약 842.22px다. 첫 화면 하단에 걸리지만 142px 세로 스크롤 뒤 완전히 보인다. 이 조건을 처음부터 모든 버튼이 보이는 화면으로 기록하지 않는다.

## 실제 viewport와 캡처 한계

| 조건 | 실제 innerWidth×innerHeight | clientWidth / scrollWidth | 문서 높이 |
|---|---|---|---|
| desktop 4개 | 1280×800 | 1280 / 1280 | 800 |
| mobile 일반 2개 | 390×844 | 390 / 390 | 844 |
| mobile 긴 URI 2개 | 390×844 | 375 / 375 | 986 |

Orca의 `set viewport` 응답만으로 크기를 판정하지 않았다. 탐색 뒤 크기가 초기화되는 것을 확인했다. 각 탐색 뒤 다시 설정하고 실제 `innerWidth/innerHeight`와 PNG IHDR를 확인했다. 기본 캡처 8개와 하단 캡처 2개는 모두 지정 viewport와 정확히 같다. 긴 모바일의 15px 차이는 세로 scrollbar다. `scrollWidth == clientWidth`이므로 가로 overflow는 없다.

Orca `full-screenshot`은 세로 scrollbar를 제외해 375×986 PNG를 만든다. 이 2개는 전체 내용 확인용 보조 자료다. 최초 캡처 도구는 full 이미지 폭을 390px로 가정하여 exit1을 냈다. 실제 `clientWidth=375`를 확인하고 해당 검사를 바로잡았다. 제품 실패나 보호 차단은 아니다. 정확한 390×844 하단 캡처를 추가하여 버튼 판정을 보완했다. 원래 오류와 재실행 결과는 [verification.json](verification.json)에 남겼다.

좁은 viewport는 desktop engine의 `mobile:false`, DPR 1이다. 실제 휴대폰·터치·모바일 UA·스크린리더 검증이 아니다. 실제 로그인, 승인·취소 제출, 관리 링크 이동은 수행하지 않았다. 저장 HTML 렌더는 서버의 인증·CSRF·redirect·Google 동의 처리나 저장 결과를 증명하지 않는다.

## 탐색 패킷 처리

[packet-audit.json](packet-audit.json)에 93개 경로의 역할·제목·직접 UI 참조를 기록했다. [packet outcomes](../../jev/SAR-SSE-INSTALL-001-VISUAL-packet-outcomes.json)는 비선택 path/category 93쌍을 모두 처리한다. `direct_edit`는 비어 있다. `PLANS.md`의 update 추천은 coor 소유와 이번 제한에 따라 no_change다.

partial의 route/find/documents-find 신원 불일치는 검색 신뢰 한계로 유지했다. `member.go` oversized 원문 대신 실제 스타일 정의를 직접 읽었다. `member_auth.go`는 후보 오기였으며 원천으로 사용하지 않았다. 실제 `remotePage`와 `memberStyle`이 존재함을 고정 Git blob으로 확인했다. 과거 archive의 unknown 경로는 현재 존재와 자료 역할을 확인했다. 역사적 QA·운영·설치 자료는 현재 합성 HTML에 입력되지 않는 이유로 no_change다. 전체 재독이나 예전 PASS의 새 실행 전환은 하지 않았다.

## 검증·정리·인계

[verification.json](verification.json)에 실제 명령 종료코드와 로컬 검증을 보존했다. fixture 생성, Go template 바이트 대조, 제품 diff0, 공백 검사는 각각 exit0이다. viewport·버튼 경계·overflow·PNG와 fixture 해시도 확인했다. 직접 시각 관찰은 자동 검사와 분리했다.

등록 `product-lint`/`product-test`와 FullOps 최종 게이트 결과는 아래 완료 기록에 추가한다. Windows 기본 PATH에는 `make`가 없으며 임의 설치·전역 PATH 변경·운영 서버 사용은 하지 않는다. 별도 theme/design lint는 프로젝트에 미구성이며 Go 문자열 CSS의 시각 수락을 대신하지 않는다.

검수에는 Orca embedded browser의 file preview만 사용했다. foreground/background helper, 웹 OS computer-use, 외부 브라우저, 실제 계정, 운영 DB, 새 유료 API, local Docker를 사용하지 않았다. 자동승인 거부를 재시도하거나 보호를 완화하지 않았다. 임시 브라우저 탭은 종료 시 정리한다. 보존된 HTML은 재현용 합성 QA fixture다.

coor는 이 좁은 합성 시각 PASS와 별도 tester의 고정 SHA 동작QA·독립 리뷰를 함께 확인한다. 실제 Google 동의·Dots 등록 ID·public callback·앱 실행·호스트 wake는 coor/사용자의 후속이다. 재개 조건은 후보 수락과 실제 앱 권한·등록 ID 및 승인된 실제 설치 환경이다. 제품 미완료 조건을 시각 PASS로 대체하지 않는다.
