---
title: 공개 안내 DEV 검증과 격리 QA 인계
status: draft
updated: 2026-10-15
owner: dev
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 인증 상태·공개 정보·실제 제품 검사와 후속 독립 QA를 구분한다
---

# 공개 안내 DEV 검증

2026-10-15 같은 과제의 CSS 두 결함 수정은 [재개 검증](css-resume/report.md)에 기록했다. source `f2b8cb6`의 CTA 대비 7.958:1, 예시 간격 24px, Chrome 38개 검사와 Linux make lint/test 종료0을 확인했다. 아래 기록은 최초 구현의 실행·실패 이력이다. tester5149ac1의 실제 native keyboard/200% PASS를 재사용하므로 아래 당시 native 미검증을 현재 미완료로 해석하지 않는다.

비교 기준은 eae9a8bffae706410f684e4c837e396e7b0bcf26, 착수 HEAD는 b0f6596b6b7bfd4ea4e3b8cfed6d8b1e20a9b93e다. 구현과 기술 판단은 [DEV phase](../../../exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md)에 있다. 이 보고는 DEV 자체 검사이며 독립 QA·직접 시각 수락·코드 리뷰·운영 배포를 대신하지 않는다.

## 검사와 실패 보존

루트에서 `go test ./internal/relay -run 'TestLanding|TestMemberUIFixtures|TestGoogle|TestHuman' -count=1`을 실행했다. 수정 전 새 TestLandingConfiguredLoginAndPublicContent는 nil DB인 Service.Handler의 전역 admission 경로에서 panic했다. 이는 새 테스트 준비 오류다. 그 명령 뒤 orchestration check가 실행되어 셸 최종 코드는0이었지만 Go 결과는 명시적 FAIL이었다. 이 실행을 통과 또는 정확한 Go 종료코드 증거로 쓰지 않는다.

수정용 inline Python/PowerShell 명령은 자동 승인 검토의 `blocked by policy`로 거부됐다. 구체 사유는 없었다. 원 명령을 변형 재실행하지 않았다. coor ask msg_8b9dc2666080과 답변 msg_d75b88015752에 따라 명시적 apply_patch로 테스트 계층을 수정했다. coor가 원문 차단을 자기 외부 blocked-test-edit-evidence.json에 보존했다.

수정 후 같은 관련 단위 회귀는 Go 결과 `ok`와 명령 자신의 종료0을 확인했다. startPage 내용/provider/error 검사는 단위에서, Cache-Control/CSP는 실제 DB integration의 익명·유효세션에서 검증한다. 기존 form golden은 변경하지 않았다. 실제 제품 gate와 DB 결과는 고정 SHA를 연결해 아래에 추가한다.

## 사실·탐색·디자인 근거

README의 네 설치/업데이트 앵커와 installer Command/Arguments·신규 폴더 조건을 대조했다. 공개 예시는 합성 문구이며 실제 회원 데이터가 아니다. 수신·표시·답장 승인·180초·24h의 의미는 D02/D04와 구조/운영 문서를 유지한다.

[source-audit.json](source-audit.json)은 startup packet의61경로에서 원문 hash·제목·관련 구간을 보존한다. 과거 FAIL·미검증·partial/unknown과 잔여 후보를 유지한다. 초기 packet의 큰 D04/member.go는 관련 공개·인증 절을 직접 읽었다. 과거 보고의 PASS를 새 코드 실행으로 재사용하지 않는다.

frontend.py의 codex/gpt-6-astra/high check는 configuration_ready다. 인증·브라우저·MCP는 이 구성 검사로 통과하지 않는다. 주 스킬 frontend-design, 기존 memberStyle과 D04 여섯 토큰을 사용했다. 별도 디자인 lint와 다중테마는 없다.

## 격리 자원과 후속

새 검증 자료는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-dev/`에 보존한다. [gate/remote.py](gate/remote.py)는 known_hosts·RejectPolicy를 유지하며 현재 HEAD bundle만 새 `/tmp/knowslink-landing001-<SHA>-<suffix>`에 설치한다. `make.cmd`는 등록 명령의 make lint/test를 같은 HEAD Linux clone으로 전달하고 원 종료코드를 반환한다. 레포 lint 설정·필수 검사는 바꾸지 않는다.

합성 Postgres는 과제 고유 컨테이너다. 테스트용 DB는 검사 후 stop하고 보존한다. 브라우저 fixture는 별도 합성 DB·메일 mock·Google 취소 mock을 사용한다. 연결/송신 API는 fixture allowlist로 거부한다. 실제 로그인·실메일·실상대 송신·운영 DB·Tunnel·기존 컨테이너 변경은 없다. [gate/forward.py](gate/forward.py)는 지정 fixture의 loopback 포트만 Windows loopback으로 전달한다.

독립 tester는 실제 handler·POST·redirect·정보 비노출을 확인한다. designer는 같은 후보의 1440/390/320 실제 캡처를 판정한다. 별도 세션 리뷰와 main/origin 통합은 coor 후속이다. 임시 자료의 삭제·이동 정리는 수행하지 않는다.

## 고정 구현 결과

제품·테스트·fixture SHA는 `4795d619513fde4fbcd4a55b94ec91b056ebbaf0`이다. 뒤의 완료 기록 커밋은 문서와 QA 증거만 추가한다. 최종 SHA는 worker_done에 전달하며 최종 필수 gate 영수증은 외부 `landing-guide001-dev/final-lint.json`에 남긴다.

| 실제 실행 | 결과·근거 |
|---|---|
| FullOps lint --from eae9a8b… | ERROR 0, WARNING 2, unavailable 0; [first-lint.json](first-lint.json) |
| 등록 make lint / make test | 실제 Linux 고정 clone, 각각 종료0; logs의 SHA-lint.log/SHA-test.log |
| make build | 종료0; SHA-build.log |
| DB race integration | TestLanding/EmailIdentity/GoogleHTTP/GoogleDeviceHTTP/RemoteOAuth/HumanPostgres, PASS 13.543s; SHA-integration.log |
| deliverables --strict | 13종, error 0/warning 0 |
| git diff --check | 종료0 |

크기 경고는 기존 PLANS 1833줄(기준1827, coor 준비 변경)과 D04 798줄(기준796, 구현 연결2줄)이다. 원문 기록을 삭제해 경고를 숨기지 않았다. D05는 GET `/`의 유효 세션303이200으로 바뀌었으므로 현재 계약을 추가하고 옛 계약을 역사로 명시했다. D10 구현 책임, D11 DEV 후보, D04 구현 연결도 갱신했다. 운영 수락 상태는 변경하지 않았다.

## 브라우저 증거와 한계

[browser.py](browser.py)는 별도 isolated profile `42ea2bbf-1caa-49a7-bbe4-c75037e54a28`, page `7c54ca56-3195-4698-8d36-1d120a8d6997`에서 실행했다. layout16·states9·audit11 assertion을 확인했다. 네 설정별 폼, Google 취소의 위/아래 복구, 만료 안내, 이메일 실제 POST→확인→이름 설정, 로그인한 공개 `/`의 개인정보·폼 비노출을 확인했다. 인증은 DOM의 실제 button.click 및 실제 handler/redirect로 수행했으며 native 마우스/키보드 성공을 뜻하지 않는다. 계정이 없는 이메일 `landing@example.test`만 합성 mail mock으로 인증했다.

CSS viewport1440×900/390×844/320×844에서 hero CTA·전체 overflow·44px min-height·단일h1·lang을 확인했다. 44px의 fractional layout 값43.988964는 브라우저 반올림 허용0.02px로 비교했다. [browser-audit/captures.json](browser-audit/captures.json)은 PNG 크기와 SHA256이다. inherited zoom1.095445 때문에 requested viewport를 보정해 실제 innerWidth를 확인했다. desktop은 네이티브 캡처 반복 타일을 피하기 위해 DPR 약0.712로 캡처하므로 PNG가1440px이라는 뜻이 아니다. 모바일 DPR 약1.095다. 첫 타일 캡처와 viewport/rounding 실패 JSON도 보존한다.

DEV가 hero1440/390, install320, examples390, Google 오류 복구 캡처를 직접 열어 내용·줄바꿈·겹침을 확인했다. 이는 독립 designer 수락이 아니다. [contrast.json](browser-audit/contrast.json)의 실제 CSS 색상 대비는 본문14.70, 보조6.19, 링크/primary7.96, 오류8.01, 선택면 본문12.98, 테두리3.79다. CSS zoom200%에서 overflow 없음도 확인했지만 native browser zoom200%는 지원 명령이 없어 미실행이다.

**native keyboard는 unavailable/not_run**이다. tab switch --focus, click, DOM focus, exec focus, keypress Enter, exec press Enter가 성공을 반환해도 document.hasFocus=false이고 summary가 열리지 않았다. 초기 [keyboard 실패](browser-keyboard/checks.json)와 [입력 한계](browser-keyboard/runtime-limits.json)를 보존한다. coor msg_526751460434는 정책 변경·가짜 DOM key event 없이 이 한계를 기록하고 완료하도록 지시했다. Enter/Space·Tab/Shift+Tab·실제 키보드 form 제출·native200%는 독립 tester의 실제 실행기로 확인할 때까지 미수락이다. DOM summary click과 앵커 focus, focus-visible3px 규칙 검사는 별도 audit에만 기록한다.

states 최초 native click도 실제 이동이 없어 FAIL이었다. DOM 실제 제출로 재실행한 결과는 별도이며 native-click-failure.json을 보존했다. synthetic mail code를 browser fetch로 읽으려던 시도는 제품 CSP connect 제한으로 `Failed to fetch`였다. CSP를 완화하지 않고 로컬 helper가 `/fixture/code`를 읽었다. audit 최초는 동일 host의 포트 간 cookie 공유로 익명 anchor가 없어 실패했으며 읽기 전용 localhost host에서 재검사했다. 새 profile 또는 GET 전 cookie 상태가 필요한 이유다.

## 후속 실행용 fixture

[browser-fixture.json](browser-fixture.json)의 Google/both/email/none URL은 각각 `http://127.0.0.1:18781/`, `http://127.0.0.1:18782/`, `http://127.0.0.1:18783/`, `http://127.0.0.1:18784/`다. 위 page는 DEV 검증용이므로 후속 QA는 새 isolated profile을 사용한다. 동일127.0.0.1의 cookie는 포트 간 공유하므로 로그인 이후 상태 검사를 주의한다. source SHA4795d619…의 실제 handler와 별도 합성 DB다. 현재 foreground fixture와 Hidden forward PID26340을 유지한다. 24시간 자동 종료 한계가 있다.

실행 원천은 `gate/remote.py fixture`와 `gate/forward.py`다. 현재 원격 경로는 `/tmp/knowslink-landing001-4795d619513f-725b85fc`, DB는 `kl-landing001-fixture-4795d619513f-ebfe4439`다. GET 허용: `/`, `/auth/verify`, `/auth/google/callback`, `/home`, `/home/name`, `/home/account`. POST 허용: `/auth/start`, `/auth/verify`, `/auth/google`, `/auth/logout`, `/home/profile`. 합성 helper GET은 `/fixture/code`, `/fixture/google-cancel`, POST는 `/fixture/expire`, `/fixture/stop`이다. 나머지 제품 endpoint는 deny한다. Google은 로컬 취소 mock, 메일은 합성 code만 사용한다.

유지·종료 인계 책임은 coor/tester다. 종료 시 loopback `/fixture/stop` POST로 test 서버를 종료하고 해당 이름의 DB 컨테이너만 docker stop, 해당 Hidden forward PID만 Stop-Process한다. 이는 종료 안내이며 DEV는 아직 실행하지 않았다. DB·clone·bundle·로그·캡처는 삭제하지 않는다. 자격·운영 계정·실메시지 본문은 보고에 포함하지 않는다.
