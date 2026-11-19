---
title: 회원 웹 polish 후보의 직접 시각 수락
status: review
updated: 2026-10-12
owner: designer
tasks: [SAR-FRONTEND-POLISH-001-VISUAL]
summary: 고정 후보의 합성 PNG를 직접 검수하고 변경 수락과 기존 low·미검증 경계를 보존한다
---

# SAR-FRONTEND-POLISH-001-VISUAL — 직접 시각 검수

**판정: 제공된 변경 화면의 시각 PASS, 비차단 잔여 있음.** 직접 관찰한 화면에서 critical·high는 발견하지 않았다. 제품 전체의 동작·접근성 적합성·운영 수락을 뜻하지 않는다.

## 기준과 증거 동일성

- 제품 후보는 `9597bad351a2fbfbaa726444b320b8cf7000e1e2`다. 검수 시작 HEAD는 `471e28a8f18da409e9abc171b36988f494668503`다. 적용 규칙은 `fullops-common-0.3.3`, `project.md`, D04 `service-member-ui.md` 4–12절이다. 지시서 기준 ref는 `7afaf02c6a652ef77cd3d330afc0b0ecbbceb1bf`다.
- DEV가 기록한 before 실행 코드는 `7afaf02`, after 실행 코드는 `f7b4786`이다. 두 manifest의 PNG 27개·26개를 SHA256와 bytes로 대조했다. 불일치는 0개다. 실제 PNG 크기는 1440×900 또는 390×844다. [동일성 검사](identity.json)에 manifest 해시와 비교 결과를 보존했다.
- `git diff --exit-code f7b4786 9597bad… -- cmd internal adapters db scripts Dockerfile compose.yaml Makefile go.mod go.sum sqlc.yaml`은 exit 0이다. 후보와 검수 시작 HEAD의 같은 비교도 exit 0이다. 제품 변경이 없으므로 after 캡처를 재사용했다. 후보 SHA에서 새 브라우저 실행을 했다고 주장하지 않는다.
- 출처는 [DEV 보고서](../SAR-FRONTEND-POLISH-001-DEV-test/report.md)의 `TestMemberUIFixtures` 합성 HTML과 Orca 내장 브라우저(Chromium 계열)다. 캡처 시각·정확한 Chromium 버전은 manifest에 없어 확인하지 못했다. 실제 계정·키·관계·메시지·운영 서비스는 사용하지 않았다.
- 호스트 tuple은 지시서의 codex / gpt-6.1-sol / medium이다. 주 스킬은 impeccable다. Operate 모드의 critique 판단 항목과 polish의 기존 디자인 보존 원칙을 적용했다. 이번 작업은 고정 캡처 수락이며 별도 전체 `impeccable critique` 명령이나 재설계가 아니다. 하위 위임은 사용하지 않았다.
- `frontend.py --repo . --host codex --model gpt-6.1-sol --effort medium --check`는 exit 0, `status: disabled`다. 프로젝트 설정은 변경하지 않았다. 인증·MCP·새 브라우저 QA는 `not_run`이다. 설치 스킬 원문과 참조는 읽었다. `impeccable context`는 실행을 시도했으나 출력·종료코드를 별도 보존하지 않아 실행 준비의 성공 근거로 사용하지 않는다. 기존 D04와 제품 템플릿을 직접 읽었다.

## 직접 열어 본 화면과 수락 근거

`view_image`로 after manifest의 PNG 26장을 모두 직접 열었다. home desktop 7장, home mobile 10장, home-empty mobile 4장, 연결 수단 mobile 2장, remote-consent·verify·refusal-limited mobile 각 1장이다. 모든 이미지 경로는 [after manifest](../SAR-FRONTEND-POLISH-001-DEV-test/after/manifest.json)에 있다.

before는 home-empty-390-s3of5, remote-consent-390, verify-390, connection-token-390-s1of2, home-390-s2of10의 5장을 직접 열어 비교했다. 다른 before PNG는 해시만 확인했다. 스크롤 조각의 위·아래 잘림과 마지막 조각의 내용 반복은 캡처 구간 중첩이며 제품 overflow로 판정하지 않았다.

| 대상 | 직접 관찰과 D04 대조 | 판정 |
|---|---|---|
| home 1440×900·390×844 | 조용한 밝은 배경, 흰 영역, 파란 주 동작, 텍스트와 기호가 있는 상태 배지, 분리된 빨간 위험 구역을 유지한다. 식별자·지문은 영역 안에서 줄바꿈된다. 한국어 설명을 읽을 수 있고 화면 겹침이 없다. | 기존 디자인 유지, 시각 수락 |
| V-05 기한·위험 버튼 | connection-token의 `(최대 10분·1회 사용)`이 하나의 줄로 묶인다. home 위험 버튼은 대상과 결과를 두 줄로 균형 있게 나눈다. | 제공된 390px 사례 resolved |
| V-06 빈 상태 | before의 빈 요청 상태 select·실행 버튼이 after에서 사라진다. 관계에는 에이전트 연결 후 초대하라는 다음 동작이 남는다. | 시각 resolved. 제출 동작은 독립 QA 담당 |
| V-07 활성 수 | home-empty 첫 조각의 `활성 0/5`를 확인했다. | 합성 fixture 시각 resolved |
| V-10 OAuth 동의 | 제품 헤더·흰 영역·전체 폭 파란 승인 버튼과 보조 취소 버튼이 생겼다. 승인 결과 설명과 관리 복귀 링크가 읽힌다. | 동의 화면 390px 시각 resolved |
| V-11 확인 코드 | 입력이 before보다 높고 라벨·오류·남은 시도·재발송·다른 주소 링크의 위계가 유지된다. PNG 입력은 비어 있다. | 입력 크기 시각 resolved. 입력된 숫자 자간은 PNG로 미판정 |
| N-01 영역 끝 여백 | 동의·거절 화면의 마지막 링크 아래 간격은 기존 영역 안쪽 여백과 일관된다. | 제공된 사례 resolved |
| 승인 요청 UUID 링크 | home-1440-s6of7에서 각 UUID는 한 줄이다. mobile에서도 값이 영역을 넘지 않는다. | DEV range 줄 수 경고를 새 결함으로 세지 않음 |

위계의 장점은 상태·다음 동작·위험 동작이 서로 구분되는 점이다. 빈 폼 제거는 선택할 수 없는 동작을 줄인다. OAuth 승인 강조는 사용자에게 설치 연결의 결정 지점을 명확히 보여 준다. 에이전트의 연결 완료가 현재 접속을 뜻하지 않는다는 설명도 유지한다.

## V-02·V-03·V-08·V-12의 제품 범위 판단

| ID | 판단·사유 | 남은 일·담당·재개 조건 |
|---|---|---|
| V-02, low | 받은 초대에 수락·거절·관계 철회가 함께 남는다(home-390-s5of10, home-1440-s4of7). 혼동 가능성은 있으나 승인 주 동작과 위험 외곽선은 구분된다. DEV가 밝힌 pending unpair 계약을 시각 polish로 바꾸지 않는 판단을 수용한다. 해결된 것으로 표시하지 않는다. | designer가 거절·철회의 사용자 의미를 별도 정리하고 DEV가 서버 계약을 검토할 때 재개한다. 현재 동작 유지로 이번 후보를 비차단 수락한다. |
| V-03, low | KST 날짜 전체 표기는 읽을 수 있고 F-UI-04를 유지한다. 짧은 날짜 요구의 미충족은 남는다. 공유 clock의 형식을 이번 과제에서 바꾸지 않는 판단을 수용한다. | 짧은 날짜 표기 과제에서 designer가 표시 규칙을 확정한 뒤 DEV가 공유 영향과 검사를 수행한다. |
| V-08, low·범위 밖 | D04가 gate 문구와 receipt 상태 의미 보존을 요구한다. 이번 PNG 묶음에는 해당 화면이 없어 새 직접 시각 판정을 하지 않는다. 과거 보고의 영어 표현 잔여를 보존한다. | gate·receipt 문구 과제에서 designer가 의미를 정리하고 DEV·tester가 상태별 렌더와 동작을 확인한다. |
| V-12, low | mobile 이동 링크가 확인 필요 띠보다 앞선다. D04 4.2 순서와 차이는 남지만 home-390-s1of10 첫 화면에 띠·두 요청 링크·생성 동작이 보인다. 이번 후보를 비차단 수락한다. | 별도 순서 정리 시 DEV가 DOM·랜드마크·탭 순서를 함께 검토하고 designer가 390·768·1440 배치를 확인한다. |

잔여 low: mobile home의 에이전트 영역 머리말에서 `새 에이전트 만들기`가 전체 폭이 아니다(home-390-s1of10). D04 8.4의 mobile 주 버튼 전체 폭 요구와 차이가 있다(N-02, 이번 직접 관찰). 버튼 이름은 잘리지 않고 주 동작으로 구분된다. 이번 수락은 비차단이며 기존 디자인 요구를 변경하지 않는다. 다음 레이아웃 과제에서 DEV·designer가 다른 생성 상태와 함께 확인한다.

기존 명령 줄바꿈 부채도 유지한다. home-390-s8of10에서 `--confirmed`의 `--`와 `confirmed`가 서로 다른 줄에 나온다. 문자열 복사 값이 같다는 DEV 설명은 재검증하지 않았다. 설명용 코드가 읽히지만 수동 입력 혼동은 남는다. 다음 코드 표시 개선에서 확인한다.

## 재사용 검사와 미검증

- [measurements.json](../SAR-FRONTEND-POLISH-001-DEV-test/after/measurements.json)은 29개 fixture × 1440·390·320의 87개 결과다. overflow·44px 미만·중복 버튼 이름·펼침 표시 단독 줄이 없다는 DEV 결과를 재사용한다. PNG 26장이 29개 화면 전체의 새 직접 검수를 뜻하지 않는다.
- [keyboard.json](../SAR-FRONTEND-POLISH-001-DEV-test/after/keyboard.json)의 home·remote-consent 포커스 `solid 3px`, Enter 펼침, failures 0을 재사용한다. PNG에는 포커스 상태가 없어 직접 입력 검증으로 표시하지 않는다. 건너뛰기 링크 키보드 이동은 DEV도 미검증이다.
- V-04의 320px 개선, V-09의 접근성 이름, OAuth 반환 화면·로그인 화면 desktop/mobile 전체, 기기 승인·gate·receipt 상태별 전체는 이번 PNG로 직접 확인하지 않았다. 이 항목의 DEV 측정·기존 수락을 새 직접 관찰로 바꾸지 않는다.
- 실제 Google 왕복·실계정·로그아웃·운영 반영·200% 확대·스크린 리더·다른 엔진·정확한 렌더 대비 비율·문자 간격 확대·입력 완료 숫자 상태는 미검증이다. 독립 코드 리뷰와 동작 QA는 tester, 운영 통합은 coor/OPS가 담당한다.
- 새 detector·제품 테스트·제품 lint는 실행하지 않았다. 제품 코드·토큰·테마가 바뀌지 않았다. DEV detector의 side-tab·em-dash는 D04 8.5의 명시된 디자인에 해당하며 재사용한다. 문서 변경의 공백·링크·산출물 메타데이터를 검사한다.
- Questions skipped: 지시서가 고정 후보의 수락과 범위 판단을 위임했다. 새로운 제품 동작이나 사실 문구를 제안·변경하지 않았으므로 추가 질문은 필요하지 않았다.

## 인계

D04에는 이번 시각 수락 근거만 추가했다. 이전 실패·보류와 미검증 기록은 유지했다. 제품 파일은 수정하지 않았다. 삭제·이동·외부 발송·운영 조작은 수행하지 않았다. 본 보고와 identity.json, packet-outcomes, 완료 아카이브를 함께 커밋해 coor에 보낸다.

문서 검증: D04 strict는 exit 0(문제 0·경고 0), `git diff --check`는 exit 0, 새 보고서·phase의 로컬 링크 누락은 0개다. 명령 자체 종료코드는 [validation.json](validation.json)에 있다. 제품 코드 변경이 없어 이번 세션의 product-lint·product-test done-gate는 적용하지 않았다.
