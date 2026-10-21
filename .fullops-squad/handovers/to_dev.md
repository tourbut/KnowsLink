---
title: SAR-PUBLISH-001-DEV — 서버 캘린더 연결 제거(D018)와 /privacy·/terms 공개 페이지 구현(D017)
status: draft
updated: 2026-10-21
owner: dev
tasks: [SAR-PUBLISH-001-DEV]
summary: 서버 캘린더 연결 제거(D018)와 /privacy·/terms 공개 페이지 구현(D017)
attempt: 3ae836b54ade423d9e42bda762d3f846
base: dcf882528eb41183071d22a1c086b4c0579db2dc
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-DEV — 서버 캘린더 연결 제거(D018)와 /privacy·/terms 공개 페이지 구현(D017)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-21
- From / To: designer(SAR-PUBLISH-001) → coor → dev
- 상태: done(DEV 완료, 리뷰·tester·designer 수락 대기)
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 서버 캘린더 연결 코드·화면·테스트·fixture·env 요구 제거, 기존 저장 연결의 일회 정리 코드, `/privacy`·`/terms` 구현과 링크, 공개 ingress repo 후보(`public-ingress.yml`·`check_public_ingress.py`), 기술·사용 문서 갱신과 커밋 / 운영 서버·운영 env·live Tunnel config 변경, Cloud Console 변경, 실계정 Google 호출, DB migration, 새 외부 의존성, 정책 문안 변경(designer 확인 필요), force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (coor가 dispatch 때 기록) / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 main `4787e85`이며 서버 캘린더 연결(SAR-SCHEDULE-001)이 들어 있다. 2026-10-21 사용자 결정 D018로 이 기능을 제품에서 제거한다. 캘린더는 각 에이전트가 따로 연결하고 KnowsLink는 메시지만 중계한다. 같은 날 D017로 Google OAuth 앱(knowslink-auth)을 공개 게시한다. 게시에는 `knowslog.com` 아래 홈페이지·개인정보처리방침·약관 URL이 필요하다. 현재 `/privacy`·`/terms`는 없다(404). 결정 원문은 `decision.py --repo . list --all`의 D017·D018이다.

제품 정본(designer SAR-PUBLISH-001):

- 제거 규칙: `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` '철회 — D018' 절(SR-01–06, 기존 연결 회원 처리 1–6, 제거 완료 조건 RC-01–05). 그 아래 본문은 철회 전 기록이다.
- 정책 페이지: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(화면 기준, 링크 위치, KO/EN 문안 전문, 게시 전 확인 항목 P-01–05, 완료 조건 PC-01–06).
- 화면: D04 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`의 '서버 캘린더 연결 제거'·'개인정보처리방침·서비스 약관 페이지' 절과 '공개 첫페이지 정본'의 하단 문장.

designer가 확인한 코드 위치(기준 `dcf8825`):

- 캘린더: `internal/relay/calendar*.go`(구현·테스트·browser fixture), `google.go`(`calendarScope`·증분 동의·`ConfigureGoogle`의 키 인자), `store.go`(`State.Calendars`·`CalendarReady`·`Gate.Lookup`·sweep), `http.go`(gate lookup·disclosure), `protocol.go`(`disclosure`), `gate_summary.go`('공개되는 내용'), `human_http.go`·`human_pages.go`·`member.go`·`member_pages.go`(설정 절·`/home/calendar/*`), `i18n_en_*.go`, `cmd/relay/main.go`, `.env.example`, `compose.yaml`의 `KNOWSLINK_CALENDAR_KEY`.
- 첫페이지: `landing_pages.go`. 정책 문안의 수집·보관 사실 근거표는 SAR-PUBLISH-001 '문안의 사실 근거'에 있다.

## 적용 기준과 예외

- 규칙: fullops-common-0.3.4(`.fullops-squad/rules/common/README.md`와 연결 규칙, 특히 security: 토큰을 로그·화면에 남기지 않음), `.fullops-squad/project.md`, 문서 규칙.
- 기준 커밋: main `dcf882528eb41183071d22a1c086b4c0579db2dc`. coor가 designer 결과 병합 뒤 그 main SHA로 갱신할 수 있다.
- ID 접두어: `SR-`·`RC-` → SAR-SCHEDULE-001.md '철회 — D018'. `PC-`·`P-` → SAR-PUBLISH-001.md. `U5` → `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md`. `MVP-` → `product-specs/SAR-MVP.md`. 산출물 `D01~D13` → `docs/deliverables/README.md`. 결정 `D009·D011·D017·D018`은 decision.py 기록이다.
- Google 토큰 철회 API는 Context7 `resolve-library-id` → `query-docs`(golang.org/x/oauth2, Google OAuth revoke)로 확인한다. 실패하면 공식 문서로 확인한다. 기존 `calendar.go`의 revoke 경로를 재사용할 수 있으면 새 조회는 필요 없다.
- Windows에는 make가 없다. 이전 단위와 같이 격리 원격 Linux clone에서 make lint/test/integration을 실행한다. lint.py는 깨끗한 임시 detached worktree에서 실행하고 force 없이 정리한다.
- 예외 없음.

## 먼저 읽을 문서

1. SAR-SCHEDULE-001.md '철회 — D018' 절(현행 정본).
2. SAR-PUBLISH-001.md 전체.
3. D04 service-member-ui.md의 두 새 절과 '공개 첫페이지 정본' 하단 문장.
4. SAR-UX-ISSUES-002.md U5(되돌아갈 승인 상세 기준)와 SAR-MVP.md MVP-08.
5. D10 `module-design.md` '캘린더 연동 모듈'과 D03 `architecture.md` 'Google Calendar 빈 시간 조회' 절(지울 범위 확인).
6. `.fullops-squad/contexts/dev.md`.

필요 시 확인: `data-model.md`·`interface-design.md`·`crud-design.md`의 캘린더 항목, `docs/development.md`·`docs/operations.md`·`README.md`의 캘린더 문장, SAR-SCHEDULE-001 DEV·QA 보고(이전 구현 근거).

## 해야 할 일과 파일 소유권

- [x] 짧은 기술 계획(완료 보고): 제거 범위, 기존 상태 JSON(`Calendars`·`Gate.Lookup`·`GoogleAttempt.Calendar`)의 일회 정리 방식과 시점, 키 유무에 따른 철회 요청 처리, `schedule.query` gate를 U5 기준으로 되돌리는 방식, rollback 영향.
- [x] RC-01–05 구현: 설정 캘린더 절, `/home/calendar/*`, 증분 동의, freebusy 호출, 토큰 봉인, disclosure, D010 범위·횟수 검사, 캘린더 문구(KO/EN)를 제거한다. Google 로그인 scope는 `openid email`만 남긴다.
- [x] 기존 연결 정리(SAR-SCHEDULE-001 '기존 연결 회원 처리' 1–3·5·6): 각 연결의 Google refresh token 철회를 한 번 요청한다. 결과와 무관하게 토큰·연결·미전달 조회 결과를 삭제한다. consume되지 않은 해당 gate는 `denied`다. 키가 없으면 철회 요청 없이 삭제한다. 기록은 연결 수·성공 수·실패 수만 남긴다. 운영 실행에 `KNOWSLINK_CALENDAR_KEY`가 필수가 아니다.
- [x] `/privacy`·`/terms` 구현(SAR-PUBLISH-001 화면 기준과 KO/EN 문안 전문): 첫페이지 템플릿 스타일 재사용, 고유 title, `aria-current`, 시행일 단일 상수. 첫페이지 하단·`#start` 한 줄·회원 설정 링크를 추가한다. `{ }` 자리 표시는 P-01–04의 값 또는 그 표에 적힌 대체 문구로 채운다. 시행일(P-01)은 2026-10-22(coor 지정, 운영 배포 예정일. 배포일이 바뀌면 ops가 상수만 고친다).
- [x] 공개 ingress repo 후보: `deploy/knowslink/tunnel/public-ingress.yml`과 `scripts/check_public_ingress.py`에 `GET /privacy`·`GET /terms` 정확 경로를 추가한다. `/home/calendar/*`가 더 필요 없음을 반영한다. live config는 바꾸지 않는다(ops 몫).
- [x] 테스트: 캘린더 테스트·fixture를 삭제하거나 대체한다. 실패 경계 자동 검사: 정리 코드(키 있음/없음, 철회 실패에도 토큰 0, 미전달 결과 미공개), `schedule.query` 승인 뒤 disclosure 없음과 U5 문구, 정책 두 페이지 200·두 언어·영어판 한국어 없음·링크. 기존 친구·메시지·승인 회귀를 확인한다.
- [x] tester·designer용 격리 fixture: 정리 전 연결 상태가 있는 합성 데이터에서 시작하는 경로와 정책 페이지 KO/EN을 볼 수 있는 기동 방법을 DEV 보고에 적는다.
- [x] 문서 갱신(실제 변경만): D03 architecture, D05 interface, D06 data-model, D09 crud, D10 module-design, D11 user-guide, D12 ops-guide, D13 transition(제거 배포 인수 항목: 운영 키 제거 시점, Cloud Console scope·Calendar API 정리, 공개 ingress 두 경로), `docs/development.md`·`docs/operations.md`·`README.md`의 캘린더 문장, `docs/deliverables/README.md` D02 줄(SAR-SCHEDULE-001 철회, SAR-PUBLISH-001 추가).
- [x] 원격 Linux make lint/test/integration, FullOps lint.py --from (지시서 base), DEV 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`, packet-outcomes, finish, 커밋, worker_done.

소유 파일: `internal/relay/`의 캘린더 제거 경로·정책 페이지·첫페이지·설정 링크와 테스트, `cmd/relay/main.go`, `.env.example`, `compose.yaml`, 공개 ingress repo 후보 두 파일, 위 기술 문서, 자기 기록·DEV 보고.

예상 규모: 삭제가 크다(캘린더 구현·테스트 약 1,800줄 이상). 추가는 정책 템플릿·문구 약 300–400줄이다. SIZE 경고는 제거 단위라 분할하지 않은 이유와 함께 보고한다.

worker는 항목을 끝낼 때마다 이 목록에 체크한다. 대화가 요약돼도 이 목록이 진행 상태의 기준이다.

## 완료 기준과 검증

- RC-01–05와 PC-01–05가 자동 검사와 실제 기동 fixture로 확인된다. PC-06(designer 시각 검수·tester 확인)은 수락 절차다.
- 원격 Linux make lint/test/integration EXIT=0.
- FullOps lint `python3 <runtime>/scripts/lint.py --repo <dev 워크트리> --from <base>` ERROR 0. 명령 자신의 종료코드를 `| tail` 등으로 가리지 않는다.
- 검증 증거: `.fullops-squad/` 아래 `.md`는 front matter가 필요하다(DOC-003). 원시 로그는 `.txt`/`.log`로 둔다.
- worker 완료는 제품 수락이 아니다. 리뷰·tester·designer 시각 검수·운영 배포·Cloud Console 정리·앱 게시는 coor가 배정한다.

tester 검증 지시서 기준(coor가 dev 완료 뒤 작성):

- 격리 fixture에서 정리 전 연결 데이터(키 있음/없음)로 시작한다. 정리 뒤 토큰·연결·미전달 결과가 0이고 해당 gate가 `denied`다. 정리 기록에 비밀 값이 없다.
- 설정·승인 상세 KO/EN에 캘린더 문구가 없다. `schedule.query`는 승인할 수 있고 승인 뒤 disclosure가 없다. `schedule.commit`은 기존 동작이다.
- `/home/calendar/connect`·`/home/calendar/disconnect` 요청이 동작하지 않는다. Google 로그인 요청 scope가 `openid email`이다.
- `/privacy`·`/terms`가 비로그인 200이다. KO/EN 문안이 정본과 같다. 영어판에 한국어가 없다. `{ }`가 없다. 첫페이지·`#start`·설정 링크가 있다. 1440·390·320px에서 가로 넘침이 없다.
- 친구·메시지 송수신·승인 회귀가 없다. 운영·실계정·실메일은 쓰지 않는다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04 `service-member-ui.md`, SAR-PUBLISH-001.md '화면 기준'. 토큰은 `memberStyle`·`body.landing`이다.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: `landing_pages.go`의 헤더·하단·`lang-nav`, 기존 표·`muted` 스타일, 기존 i18n 구조. 새 테마·JS·외부 리소스는 없다. 예외: 정책 본문 표는 320px에서 표 영역만 가로 스크롤한다.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py. DESIGN 경고는 완료 보고에 처리 이유를 적는다.
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환은 없다. KO/EN 전환 뒤 `/privacy`·`/terms`·첫페이지·설정·승인 상세를 1440·390·320px로 캡처해 DEV 보고에 연결한다.
- 피할 것: 크림 배경, 알약 버튼, 모노스페이스 라벨, 장식 그라디언트, 이모지 아이콘, 접는 상자, 영어판의 한국어 잔존.

## 갱신할 산출물

D03·D05·D06·D09·D10·D11·D12·D13 중 실제 변경(위 문서 갱신 항목). D02·D04는 designer가 SAR-PUBLISH-001에서 갱신했다. 문안 변경이 필요하면 designer에게 묻는다.

## 기대 산출물

제거·정책 페이지 코드와 테스트, 공개 ingress repo 후보, 갱신한 기술 문서, `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`, packet-outcomes, 완료 보고.

## 제약·협업·후속

- 하위 작업은 [위임 기준](../rules/delegation.md)을 따른다(standard, 최대 2). 나누면 제거 경로와 정책 페이지로 파일 소유를 분리한다.
- 정책 문안·링크 위치·제거 규칙의 변경은 제품 판단이다. coor를 거쳐 designer에게 묻는다.
- 운영 env·live Tunnel·Cloud Console·실계정은 범위 밖이다. dev 완료 뒤 coor가 tester 검증과 ops 배포를 배정한다. 배포 뒤 coor가 사용자와 Cloud Console 브랜딩 입력, `calendar.freebusy` scope·Calendar API 정리, 앱 게시를 한다.
- 법률 검토(P-05)는 사용자 몫이다. DEV는 문안을 법적 최종본으로 표시하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-DEV-packet.json` / SHA `4d6c3a3232dc7ca08cbd933f69e166dcceae3c81` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 9, 9, 18, 18, 18, 23, 24, 43, 47, 47, 129, 217, 241, 321, 335, 359, 366, 372, 394, 394, 403, 403 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 14, 18, 40, 47, 54, 69, 235 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 274 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 96 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 328 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 17, 51 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/system-design.md` (document_read) · 줄 75, 78, 113 · inferred
- `internal/relay/calendar.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/google.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/http.go` (direct_edit) · 줄 전체/미확인 · unknown
나머지 참조 후보 6건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 6건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; internal/relay/calendar.go: missing, oversized, sensitive or invalid; manual inspection; internal/relay/google.go: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

- 브랜치 / SHA: `fullops/dev` / 제품 최종 `e20c199a4ec2d81c434962954614a8f0feba2b85`(원격 게이트 기준). 이후 커밋은 보고·지시서·packet-outcomes 기록뿐이다.
- 변경 이유: D018 서버 캘린더 연결 제거와 기존 연결 일회 정리, D017 `/privacy`·`/terms` 공개 페이지와 링크, 공개 ingress repo 후보 두 경로.
- 기술 계획: SAR-SCHEDULE-001 이전 `67284e3` 판으로 캘린더 파일을 되돌렸다(`internal`·`cmd`·adapter·env). 정리는 HTTP 시작 전 `RetireCalendars` 한 곳이며 legacy 필드는 그 정리만 읽는다. 상세는 [DEV 보고](../docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md) '기술 계획'.
- 검토 필요: `calendar_retire.go`(토큰 처리·로그 비노출), `policy_pages.go`·`policy/*.html`, `landing_pages.go`의 style·footer 분리. 정책 페이지 하단에 '홈'을 넣은 위치(서비스 약관 뒤)는 designer 확인 대상이다. P-04 이메일 업체 행은 운영 SMTP 미설정이라 뺐다.
- 검증한 것: 원격 Linux make lint/test/integration EXIT=0(`e20c199`), 격리 fixture 기동·probe(정리 전 상태 → 키 있음·철회 거부·키 없음 모두 토큰·연결·결과 0, gate denied, 정책 200, 영어판 한국어 0, `/home/calendar/*` 404), 정본 문안 4개 일치, 1440·390·320 가로 넘침 없음. FullOps lint `681bacf` ERROR 0·WARNING 17(처리 이유는 DEV 보고)·실행 불가 2(Windows make 없음, 원격 Linux에서 실행).
- 못 한 것: 실제 Google revoke 호출, 운영 상태 정리, live ingress, 키보드·200% 실행, 독립 리뷰·tester·designer 수락.
- 산출물·로그: [DEV 보고](../docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md), [render](../docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/render/), [packet-outcomes](../docs/evaluations/jev/SAR-PUBLISH-001-DEV-packet-outcomes.json).
- 의존성: 변경 없음(go.mod·package.json 미변경). revoke는 표준 `net/http` POST다.
- SIZE: 제거 단위라 분할하지 않았다(제품 코드 +232/−673, 테스트 포함 +898/−2,052).
- 남은 일: coor가 고정 SHA 독립 코드 리뷰를 배정한다. tester가 격리 fixture로 독립 QA를 한다. designer가 KO/EN 3폭 직접 시각 검수(PC-06)를 한다. ops가 키를 둔 채 제거 배포·정리 로그 확인·키 제거·live ingress 두 경로를 한다. coor·user가 Cloud Console scope·Calendar API 정리와 앱 게시, 법률 검토(P-05)를 한다.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
