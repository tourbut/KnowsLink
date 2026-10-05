---
title: 일반 이메일 신원 후보의 직접 UI 검수 결과
status: review
updated: 2026-10-05
owner: designer
tasks: [SAR-PUBLIC-IDENTITY-001-UI]
summary: 고정 후보의 로컬 fixture UX01–03 시각 PASS와 Orca 캡처 장애 및 실제 확인 미실행을 구분한다
---

# SAR-PUBLIC-IDENTITY-001-UI — 직접 시각 검수

## 판정과 수락 경계

고정 후보 `59b66ada8b36802484cc6d7e22523257b50572cc`의 **로컬 fixture UX-01–03 직접 시각 검수는 PASS**다. 새로운 UI critical/high 발견은 없다. 기존 회원 gate의 시각 표기에 low 1건을 남긴다.

Orca 전용 페이지에서 실제 입력·SMTP 코드 확인·회원 홈·재확인·로그아웃을 조작했다. Orca 캡처가 불안정해 필요한 시각 항목은 coordinator가 허용한 별도 실제 Chromium에서 보완했다. `pw-*.png`는 별도 Chromium 캡처이며 Orca PNG가 아니다. 저장된 PNG를 검수자가 직접 열어 확인했다. 자동 DOM 판정만으로 시각 PASS를 만들지 않았다.

**실제 사람 이메일, 운영 SMTP, 공개 배포, QA-P06 공개 확인, QA-P07 사람 확인, 동일 이메일 노우↔다닷 최종 시험은 미실행**이다. 이번 결과는 로컬 신원 기능 단계의 UX 근거다. 전체 일반 서비스 PS-13 수락이나 독립 보안 리뷰 수락을 대신하지 않는다. coor가 전달한 별도 보안 리뷰의 F1/F4/F2 수정 상태를 이번 UI 검사로 해소하지 않는다.

검증된 코드 단계와 운영 최종 수락은 구분한다. 독립 리뷰·QA·이 직접 UI 결과를 연결해 코드 통합 여부를 판단할 수 있다. OPS는 통합 후보로 운영 SMTP·배포·외부 확인을 수행한다. 실제 이메일과 노우↔다닷 완료 조건은 그대로 남는다. 운영 공개가 완료돼야만 로컬 코드 단계를 검수할 수 있다는 순환 조건으로 해석하지 않는다. 새 제품 수치나 기술 승인은 추가하지 않았다.

## 기준과 실행 환경

| 항목 | 실제 값 |
|---|---|
| 제품 후보 | `59b66ada8b36802484cc6d7e22523257b50572cc` |
| 구현 코드 | `a446d89ff288c4243ad6d7f8780a778517154584` |
| 지시서 준비 | `1207bdf4542fac98d228f86de79aad9f4e126ce8` |
| 지시서 기준 | `94533b207b456c0560800fe30a7c90b2b5887c6e` |
| 제품 실행 | `/tmp/knowslink-public-identity-ui-59b66ad`의 깨끗한 detached checkout |
| 기록 checkout | `/home/shin/orca/workspaces/KnowsLink/fullops-coor` |
| 실행 기간 | 2026-10-05 KST; 최초 relay 기동 22:23:55, 최종 보완 캡처 22:40–41 |
| 독립 DB | `knowslink-identity-ui-59b66ad-pg`, PostgreSQL 17, loopback port `47201` |
| relay | 후보에서 `go build`한 독립 바이너리, `http://localhost:52731` |
| SMTP | 후보 `scripts/mail_sink.py`, loopback; 초기 port `50045`, 복구 뒤 새 임시 port |
| 설정 | 합성 가입·trial allowlist·trusted IP header 미설정; 운영 인증값 미사용 |
| Orca 페이지 | `56050624-50c1-4700-b313-a3ed78e71059`; 모든 페이지 조작에 `--page` 명시 |
| 보완 브라우저 | Chromium `147.0.7727.15`, Playwright core `1.63.0`, headless, 새 비영속 context |
| 화면 폭 | `390×844`, `1280×900`; 모바일 기기 실물 검사는 아님 |
| Task / Dispatch | `task_1a9cb470bd74` / `ctx_a3efa466a50f` |
| 복귀 Run | `run_8ca8bc058ab7` |

임시 브라우저는 설치된 Chromium과 Playwright core를 읽기 전용으로 사용했다. 다른 QA의 브라우저·DB·프로필·파일을 변경하지 않았다. 두 브라우저 모두 같은 고정 제품과 같은 자기 DB·SMTP에 접속했다. 두 context는 독립 쿠키를 썼다.

기준 문서는 `fullops-common-0.3.2`, FULLOPS.md, project.md, document-writing.md, orca-agents.md, 제품 PS-01–04·신원 PS-11, UX-01–03, frozen SAR-MVP C1–C5다. designer 인박스 전문과 coor의 코드/문서 find·context 근거를 읽었다. DEV 실행 기록, README의 일반 이메일 절차, 사용자 가이드, member.go·identity.go·http.go·mail_sink.py의 실제 흐름을 확인했다.

## 관측 결과

아래 PNG는 [증거 폴더](SAR-PUBLIC-IDENTITY-001-UI-evidence/)에 있다. 같은 이름의 `.txt`는 민감값을 가린 실제 화면 텍스트다.

| 항목 | 직접 관측과 다음 동작 | 판정·대표 증거 |
|---|---|---|
| UX-01 시작 | 이메일 레이블·가입/로그인 통합 안내·확인 전 미로그인 문구가 보인다. 관리자·SSH·별도 초대를 요구하지 않는다. | PASS; `pw-01-start-mobile.png`, `pw-01-start-desktop.png` |
| UX-02 확인 대기 | 실제 sink로 발송한 뒤 코드 입력 화면이 열린다. 발송과 신원 확인을 구분한다. 마스킹한 주소·기한·5회 오답·재발송 60초·다른 주소 동작을 보인다. | PASS; `pw-02-pending-mobile.png` |
| UX-02 오답 | 오답 제출 뒤 남은 시도 4회와 빨간 오류 표식을 보인다. 코드 입력과 재시도 동작이 남는다. | PASS; `pw-03-wrong-mobile.png` |
| UX-02 제한 | 즉시 재발송을 누르면 429 안내와 구체적 재시도 시각을 보인다. 성공/로그인 완료를 표시하지 않는다. | PASS; `pw-04-resend-limited-mobile.png` |
| UX-02 만료 | 자기 DB의 challenge 기한을 과거로 둔 뒤 실제 경로를 연다. 만료 안내와 새 코드 요청 화면을 보인다. | PASS; `pw-05-code-expired-mobile.png` |
| UX-02 발송 실패 | 자기 SMTP sink만 중지한 뒤 브라우저에서 새 요청을 보낸다. 미발송·미로그인·1분 뒤 재시도를 명시한다. | PASS; `pw-13-mail-failed-mobile.png` |
| UX-03 자기 홈 | 실제 sink 코드를 확인한 세션으로 자기 신원·회원 ID·세션 기한을 보인다. agent·관계·승인 요청은 빈 목록이다. 관리 메뉴·타 회원 목록이 없다. | PASS; `pw-06-home-mobile.png`, `pw-06-home-desktop.png` |
| UX-03 준비 중 | agent 연결과 계정 비활성화에 준비 중 상태와 현재 실행 불가를 명시한다. 미구현 성공 버튼은 없다. | PASS; `pw-06-home-mobile.png` |
| UX-03 재확인 | 자기 fixture 세션의 Verified를 6분 전으로 바꾼다. 전체 로그아웃 버튼 대신 이메일 재확인 안내를 보인다. 실제 새 sink 코드 확인 뒤 같은 회원으로 돌아온다. | PASS; `pw-07-reauth-home-desktop.png`, `pw-08-reauth-pending-desktop.png` |
| UX-03 전체 로그아웃 | 최근 확인 뒤 전체 로그아웃을 누른다. 브라우저 세션 종료와 agent 자격 유지 안내가 보인다. `/home` 재접속은 세션 종료 안내로 돌아온다. | PASS; `pw-09-logout-all-desktop.png`, `pw-10-old-session-desktop.png` |
| UX-03 재로그인·현재 종료 | 같은 fixture 이메일을 다시 확인하면 같은 회원 ID와 같은 빈 목록을 보인다. 현재 로그아웃을 누르면 시작 화면과 종료 안내가 보인다. | PASS; `pw-12-current-logout-desktop.png`; continuity JSON |
| 회원 gate 변경 영향 | 자기 owner에 결속한 만료·원문 부재 fixture gate를 실제 회원 경로로 연다. 활성 승인/거절 버튼이 없고 `/home` 복귀 링크를 보인다. | PASS, low 표기 문제 별도; `pw-11-expired-member-gate-desktop.png` |
| 레이블·키보드 | 이메일·코드 입력에 연결된 레이블이 있다. 실제 Tab으로 이메일 입력에서 발송 버튼으로 이동한다. Orca의 화면 밖 클릭은 포커스와 Enter로 보완했다. | PASS; Orca DOM·명령 로그, Chromium keyboard 관측 |
| 가독성·폭 | 제목·섹션·오류·다음 동작을 읽을 수 있다. 성공/실패를 색만으로 구분하지 않는다. 지정 폭에서 버튼·문구가 잘리지 않고 수평 overflow가 없다. 긴 회원 ID는 줄바꿈된다. | PASS; 직접 PNG, `playwright-observations.json`의 width/scrollWidth |

이 검사는 상태별 UI 표시를 확인한다. challenge 만료와 최근 재확인은 자기 DB fixture의 시각을 바꿨다. 반복 캡처 때 자기 send/HTTP bucket을 초기화했다. 따라서 10분 실시간 대기, 5회 전체 오답, 12시간/60분 실시간 만료, 다중 브라우저 전체 무효화, 시간당 모든 한도 집행의 독립 기능 QA를 주장하지 않는다. 기존 QA 담당의 결과를 대체하지 않는다.

## Finding

**F-UI-01 / low / 기존 회원 gate 만료 시각의 가독성.** 새 신원 화면은 `KST` 시각을 보여 주지만 `/home/gates/ui-expired-gate`는 Go 기본 UTC 시각과 중복된 `+0000`을 보여 준다. 상태도 기존 `expired` 표기를 유지한다. 만료·원문 부재 문구와 비활성 승인 상태는 명확하므로 이번 신원 UX 수락을 차단하지 않는다.

재현: 후보의 정상 회원 세션에서 자기 owner에 결속된 만료 gate를 연다. 만료 항목의 UTC 기본 문자열을 확인한다. 근거는 `pw-11-expired-member-gate-desktop.png`다. 후속 담당은 DEV이며 gate 표시를 수정하는 과제에서 검토한다. 이번 검수자는 제품 코드를 수정하지 않았다.

## Orca 캡처 장애와 보완 근거

Orca screenshot/full-screenshot은 정상 DOM과 레이아웃 좌표가 있는 페이지를 배경만으로 캡처했다. 새 전용 탭 선택·390/1280 폭·키보드 포커스·networkidle·스크롤·RAF·무수정 화면·exec screenshot을 확인했다. 일부 캡처만 복구돼 안정된 방법으로 판정하지 않았다.

실제 Orca 직접 관측은 `04-start-focused.png`의 시작 화면과 `18-exclusive-orca.png`의 모바일 현재 로그아웃 화면이다. `orca-capture-blank.png`는 장애 증거다. 배경만 있는 PNG를 UI PASS 근거로 쓰지 않는다. coor가 같은 브라우저에서 병행 조작했다는 회신은 가능한 경합 근거이며 원인 확정은 아니다.

이 Dispatch의 coordinator `ask` 응답은 전용 페이지 독점 복구 후에도 blank가 남으면 별도 실제 Chromium에서 같은 fixed59 localhost fixture를 검수하도록 허용했다. 지시·질문·답변은 Run에 보존된다. 허용 뒤 실제 새 browser/context/page로 로그인·오답·제한·만료·재확인·종료 흐름을 실행했다. HTML 재현물이나 합성 렌더링으로 대체하지 않았다.

캡처 전에 이메일 텍스트만 `[이메일 가림]`으로 치환했다. 이메일·코드 입력과 홈 이메일 값은 회색 마스크로 가렸다. 인증값·cookie·token·메일 전문은 증거에 저장하지 않았다. 핵심 대기·미로그인 문구까지 가린 첫 보완 캡처는 증거 결함으로 보고 해당 캡처를 교체했다.

## 실행·정리·미실행

후보 Go 빌드 exit 0, 최종 실제 Chromium 흐름 명령 exit 0이다. `visual.exit`와 `visual.log`는 명령 자신의 종료코드를 보존한다. PNG 직접 관측과 자동 상태/overflow 확인은 서로 다른 근거다. 새 제품 코드를 변경하지 않았으므로 product-lint의 변경 검사 대상은 없다. 전체 제품 테스트는 별도 QA 담당 범위여서 중복 실행하지 않았다.

자기 relay·SMTP 프로세스, PostgreSQL 컨테이너와 자기 anonymous volume을 정리했다. 자기 메일 파일을 제거했다. 세 loopback 포트의 종료를 확인했고 전용 Orca 페이지도 닫았다. detached checkout은 코드 변경 없이 보존했다. [정리 근거](SAR-PUBLIC-IDENTITY-001-UI-evidence/cleanup.json)를 확인한다.

실제 사람 이메일/운영 공개의 재개 담당은 OPS·coor·사용자다. 운영 SMTP와 수락된 공개 후보가 준비되면 QA-P06·QA-P07을 수행한다. 최종 노우↔다닷은 일반 서비스 수락 뒤 같은 일반 이메일로 별도 agent·키·관계를 확인한다. 이번 결과는 새 보안 수정 후보의 실행 근거가 아니다. 해당 후보가 오면 UI 변경 영향만 별도 검수한다.
