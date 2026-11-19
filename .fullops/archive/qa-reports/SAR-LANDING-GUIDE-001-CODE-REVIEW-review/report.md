---
title: 공개 안내·인증 진입 구현 독립 코드 리뷰
status: review
updated: 2026-10-14
owner: reviewer
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 고정 구현의 인증 경계·공개 안내·필수 검사와 문서 인계를 검토한다
---

# SAR-LANDING-GUIDE-001-CODE-REVIEW

코드와 기술 문서 리뷰는 수락 가능하다. 확정 발견은 critical/high/medium/low 모두 0건이다. 이 결론은 검토자의 판단이며 OCR 자동 판정이 아니다. 독립 tester 동작 QA와 designer 직접 시각 수락은 별도다. native 키보드와 native browser 200% 확대는 DEV에서 unavailable/not_run이며 최종 제품 수락 전 검증해야 한다.

- 검토자: Codex / gpt-6-astra / high, 실제 세션 `01a13aaa-968a-7d81-b8ca-f2d71a16a0d4`.
- 구현자 실제 세션: `01a13a96-02a2-7863-a346-67d01b79d394`. 서로 다른 세션이다.
- base/merge-base: `eae9a8bffae706410f684e4c837e396e7b0bcf26`.
- 최종 head: `fe2387a7676f7997aea76a2606f4cc2dbcd4b56b`. 제품·DB·브라우저 fixture 원 실행 SHA: `4795d619513fde4fbcd4a55b94ec91b056ebbaf0`.
- 관리 snapshot: `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-6599ed38950f468c916585fcd9f1f417`. 공식 review.py snapshot으로 만들고 exact HEAD·detached·clean을 확인했다. 읽기 전용으로 유지했다.
- 정본 결과는 coor의 이 리뷰 폴더에만 작성했다. 제품·snapshot·worker·PLANS·board는 수정하지 않았다. 삭제·이동·원격 실행·운영 배포는 하지 않았다.
- OCR preview는 reviewable 61, excluded 23, 총84개다. reviewed46/skipped38, 직접 검토54.7619%다. 모든 `(path,status)`의 상태·이유는 result.json에 기록했다.

## 적용 기준과 검토

FullOps1.4.2 fullops-review와 open-code-review-delegate를 사용했다. AGENTS.md, FULLOPS.md, fullops-common-0.3.3 README/coding-style/testing/security, project.md, document-writing.md, review/rule.json과 네 OCR rule group을 읽었다. 해당 규칙·lint 설정은 base→head에서 변경되지 않았다. 원 지시와 완료 전문은 handovers/logs/2026-10-14_to_dev.md의 동일 attempt에 있다. 사용자 Astra high 지정이 frontend.json의 Sol 기본값보다 우선한다.

D02 LG-01–08, D04 공개 첫페이지 절과 D05/D10/D11을 코드와 대조했다. 과거 최소 로그인 지시는 공개 안내 범위에서 대체된다. 유효 세션의 `/` 자동 redirect를 없애고 안내200과 `/home` 링크를 제공하는 것은 LG-02의 의도된 변경이다. 회원 이름·이메일·agent·친구·본문·키를 startView에 전달하지 않는다. 이름 미설정은 기존 `/home/name`으로 이어진다.

Service.startView의 모든 호출자를 검색하고 member.go의 이메일 시작·발송 실패·코드 만료·검증 실패·legacy 오류와 googleFailure를 읽었다. Google/Mail 활성 상태를 한 곳에서 결정한다. Google-only/both/email-only/none, 오류 상단 요약과 필드 연결, 만료 안내, authed 폼 비노출을 확인했다. Google start/callback/finishGoogleLogin의 성공 Next는 `/home`, `/connect/...`, `/remote/consent/...`를 유지한다. 재확인·세션 결속·Origin/Fetch-Site·CSP/no-store·html/template escape·기존 POST 경계는 유지된다.

landing_pages.go의 다섯 구획·단일 인증 영역·두 닫힌 details·합성 text 예시를 D04와 대조했다. README의 네 설치/업데이트 제목과 installer의 Command/Arguments·새 폴더 조건이 일치한다. 수신·화면 알림·매 발송 승인·180초 큐·최대24h 웹 사본을 구분한다. 키 복사·자동 답장·가짜 운영 성공을 안내하지 않는다. 외부 폰트·이미지·script를 추가하지 않았다. 공개 CSS는 `.landing` 범위이며 memberStyle을 변경하지 않아 회원 화면 제목·배치에 전파되지 않는다. D04의 hero/h2/폭/구획 간격 예외를 적용한다.

새 단위·integration·fixture와 기존 member_ui_test/form golden·Google/device/Remote/human 호출 맥락을 읽었다. nil DB 테스트 준비 오류 수정은 startPage 내용/provider/error 검사를 유지하고 실제 DB 익명·authed 두 경로에 Cache-Control/CSP assertion을 남겼다. 원 panic·승인 거부·후속 수정은 DEV QA에 구분돼 있다. 거부된 셸 명령은 재실행하지 않았다.

remote.py/forward.py 전체에서 exact SHA·clean 검사, 고유 clone/DB, known_hosts/RejectPolicy, 비밀값 비출력, loopback 제한, stop 후 보존을 확인했다. opt-in fixture는 실제 handler와 합성 메일·로컬 Google 취소만 사용하며 연결/송신 경로는 allowlist 밖이다. browser.py의 DOM 실제 submit과 native 입력의 차이를 확인했다. 고정 page/포트와 host cookie 공유의 재사용 제약은 인계에 있다.

## 검사·증거와 규모

[verification.json](verification.json)은84개 대상의 존재·hash, JSON 파싱,17개 PNG의 실제 크기·manifest SHA256 일치, 고정 snapshot과 runtime 입력 동일성을 기록한다. 제품/테스트/의존성/빌드/배포/gate 입력은4795→fe에서 diff0이다. 이 확인으로 원 DB/브라우저/build 실행의 적용 범위를 판단했으며 fe에서 새로 실행했다고 바꾸지 않았다.

최종 작성자의 공식 FullOps lint는 정확한 fe/eae 기준이다. 외부 final-lint.json을 [lint-original.json](lint-original.json)과 [lint.json](lint.json)에 byte 그대로 보존했다. reason·상태·종료코드를 수정하지 않았다. 실제 새 Linux clone의 fe-lint.log/fe-test.log는 make lint/test 모두 EXIT=0이다. Go race와 adapter 테스트가 실행됐다. lint의 기존 candidate overwrite/UUID ERROR 출력은 음수 시나리오이고 명령 종료는0이다. ERROR0/WARNING3/unavailable0을 확인했다. 검토자가 같은 제품 검사를 반복 실행하지 않았다.

4795 원 로그의 make build와 DB integration도 EXIT=0이다. TestLandingSessionAndLoginRecovery,4 provider subtests, TestEmailIdentity, TestGoogleHTTP, TestGoogleDeviceHTTP, TestRemoteOAuthHTTP/IsolationAndRefresh, TestHumanPostgres가 실제 PASS했다. 기존 테스트의 기대문을 제거하거나 form golden을 완화하지 않았다. git diff --check base/head도 종료0이다.

SIZE-001 두 건은 기존 PLANS1833(기준1827)와 D04 798(기준796)이다. 과거 기록과 단일 디자인 정본을 유지하고 구현 연결만 추가한 설명을 수락한다. SIZE-002는 lint 선택 범위 추가406>400이다. [작성자 최종 영수증](author-final-handoff.md)은 초기395에서 D05의 공개200 계약과 DEV context 추가로406이 되었고, 단일 공개 template/auth-view와 필수 회귀·fixture의 증거를 다른 과제로 나누지 않았다고 설명한다. 실제 diff가 이 범위와 맞는다. 추가 기능·dependency/schema는 없으며 DEP-001도 없다. 경고 회피를 위한 기록 삭제나 임의 분할은 요구하지 않는다.

브라우저 layout16/states9/audit11의 원 assertion과 실패 기록을 읽었다. CSS viewport1440/390/320, 전체 overflow,44px min-height, 대비, DOM auth POST/redirect는 DEV 검사다. CSS body zoom=2는 native browser zoom이 아니다. native click/keyboard 실패, hasFocus=false, viewport/DPR 보정과 소수 반올림 실패를 보존했다. CSS/DOM 성공을 키보드·native200% PASS로 승격하지 않았다. 별도 디자인 lint와 다중테마는 없어 해당 없음이며 일반 DESIGN 검사로 Go 문자열 CSS의 시각 수락을 주장하지 않는다.

skipped38개는 PNG17, 생성 탐색/history7, 과거 설계 리뷰 JSON9, 생성 DOM snapshot4, source-audit1이다. 파싱·hash·요약/manifest는 확인했지만 과거 발췌 전체의 의미나 캡처 시각 판정을 재감사하지 않았다. 핵심 제품9개·helper3개·변경 정본과 새 검증/실패 요약·원 로그는 직접 읽었다. packet partial=true/unknown22와 producer 불일치 이력을 유지했다. outcomes의 현재 packet input hash·attempt 결속 및 핵심 read/update/direct-edit 후보를 실제 변경과 대조했다. 자동 find 교집합 점수는 보조 지표이며 호출자 읽기나 완전성 증명이 아니다.

## 대화 미참조 인계 점검

| 항목 | 정본 경로/절 | 결과 |
|---|---|---|
| 현재 요구·이유 | deliverables/README → D02 SAR-HUMAN-PEER 공개 소개 → D04 공개 첫페이지 | 확인. 최신 공개 안내 요구·최소 로그인 대체·LG-01–08이 연결된다. |
| 구조·구현/미완료 | D05 공개 진입, D10 공개 안내, D11 새 사용 흐름 | 확인.200 계약·공유 view·개인정보 경계·DEV 후보와 후속 수락을 구분한다. |
| 실행·검증 | DEV phase → DEV-test/report → 원 로그/fixture/browser.py | 확인. 실제4795와 최종fe 필수 gate·실패·native 미실행을 구분한다. |
| 운영·복구 | README 신규/업데이트, DEV QA 후속 fixture | 확인. 기존 키 보존·합성자원 유지/stop 담당이 있다. 운영 배포/실계정 시험은 해당 없음이다. |
| 다음 작업·담당 | archive 완료 보고와 DEV QA 후속 | 확인. tester 입력/확대·독립기능, designer 시각, coor 통합/자원 유지가 남는다. |
| 링크·정본 접근 | D04/D05/D10/D11와 README 네 설치/업데이트 제목 | 확인. 현재 파일과 대상 절을 읽었다. 원격 사이트 가용성·모든 renderer anchor는 미실행이다. |
| snapshot 정리 뒤 보존 | coor 리뷰 폴더와 고정 commit의 원문·QA | 확인. 결과는 snapshot 밖에 있으며 정리는 coor 후속이다. |

review.py check는 exact base/head와 task-key로 실행해 종료0이다. 결과는 reviewed46/skipped38/total84, lint WARNING3이다. 실제 명령·종료코드는 [check-execution.json](check-execution.json)에 보존했다. find 교집합은 code recall1.0/precision0.25, documents recall0.25/precision0.25다. check 통과는 기록 검사이며 독립 tester/시각 PASS·제품 최종 수락·배포 승인을 의미하지 않는다.

실제 브라우저 조작·설치·실계정 로그인·실메일·메시지 발송·서버/배포는 실행하지 않았다. native 키보드와 native200% 및 직접 시각 수락이 남아 있다.
