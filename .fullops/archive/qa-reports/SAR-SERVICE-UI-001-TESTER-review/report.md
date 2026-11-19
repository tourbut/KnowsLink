---
title: SAR-SERVICE-UI-001-TESTER 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SERVICE-UI-001-TESTER]
summary: 회원 웹 고정 후보 c7828b2의 독립 snapshot 코드 리뷰 결과와 파일별 검토·발견 사항·인계 점검을 기록한다
---

# SAR-SERVICE-UI-001-TESTER 리뷰

- 검토자 / CLI / 모델: tester 역할 / Claude Code / Claude Sonnet 5.5(`claude-sonnet-5-5`). 검토 세션 `1fe0133c-f609-42f7-a774-bf5b8b5765c2`. 구현자 세션 `5953ccd3-7f4f-4c5e-b63d-c02d36ee35ac`와 다르다. AI가 판단한 결과이며 OCR의 자동 판정이 아니다.
- base SHA / head SHA / merge-base: `8539ef5b6d3e7ed99259762c62d164b8ba40fab0` / `c7828b227b5715f7bdf02bff46bf9cf102a758d9` / `bd24be06b20ff21cb9e8de4680345b29fa9b139c`. 리뷰는 고정 후보를 읽었다. 움직이는 DEV HEAD는 보지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13(fabbdb29), `delegate` 방식(OCR 측 LLM 미사용). `review/rule.json` sha256 `f1061481dcdb66affe2f0f69a6c2354d276babcf23434dedfe6b622171994de8`. 공통 규칙 `fullops-common-0.3.3`(README·coding-style·testing·security), `project.md`, `docs/agents/document-writing.md`, 기준 커밋 tester `030c5bd`(제품 바이트는 `c7828b2`와 같다). 예외 없음.
- 요구사항·완료 기준 원천: D04 `docs/design-docs/mockups/service-member-ui.md`(필수 R-01–09, 권장 R-10–11, A-01–10, I-01–08, 12.2 완료 기준 1–6), 지시서 `handovers/to_tester.md`, DEV 기록 `docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md`. D04 12절의 R-12·다크 모드·실접속은 범위 밖이다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 19 / 18 / 1(`member_forms.golden`, `unsupported_ext`) / reviewed 15 / skipped 4(DEV 컨텍스트·패킷 outcomes JSON·DEV 인박스 비움·DEV 로그 append). 제외 1건은 `git diff`와 재설계 전 golden(`37b85a4`)을 직접 비교했다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 4(SIZE-001 3건, SIZE-002 1건) / 실행 불가 0. Linux clone(`c7828b2` detached, clean)에서 실행했다. `product-lint`(`make lint`)·`product-test`(`make test`) passed.
- SIZE-002: 변경 추가 649줄 > 400줄. DEV 기록은 템플릿 전면 재작성·29 fixture·golden이 원인이며 화면별로 나누면 공유 토큰 검수가 어렵다고 설명한다. 지시서에 별도 예상 규모 수치는 없다. 649줄 중 대부분이 `member_pages.go`(182)·`member_ui_test.go`(194)·문서·packet outcomes(577은 JSON이라 SIZE 계산 밖일 수 있음)다. 근거가 타당하다고 판단했다. SIZE-001 3건은 기준 전부터 300줄을 넘는 파일이다(connections_test 390→392, member.go 357→377, member_agents.go 331→352). view 분리로 `member.go`는 오히려 179줄 줄었다.
- DEP-001: 해당 없음. `go.mod`·`package.json`·lock 변경 0. 새 JS·폰트·이미지·라이브러리 0. Context7·외부 문서 조회는 필요하지 않았다(표준 HTML·CSS·Go `html/template`만 쓴다).
- UI 디자인: 정본 D04 8절(토큰), 공용 정의는 `memberStyle` 한 상수. 별도 디자인 lint·테마 전환은 프로젝트에 구성되어 있지 않다(project.md). 다크 모드는 D04가 후속으로 정했다. 디자인 lint는 해당 없음(Go 문자열 CSS는 FullOps DESIGN이 검사하지 못함). 대신 렌더 검사와 Orca 내장 브라우저 1440·390·320 측정을 실행했다. 시각 수락은 designer 후속이다.

## 검토 범위

result.json의 19개 `(path,status)`에 reviewed/skipped와 이유를 기록했다.

읽은 것: diff 전체, 호출자(`render`·`refused`·`memberAction`·`connectionPage`·`homePage`·`remotePage`·`page`), `member_pages.go` 템플릿 183줄 전문, `member_ui_test.go` 전문, golden, D04 전문, DEV 기록 전문, D06·D10·user-guide·README diff, 공통 규칙 4개. 패킷의 `adapters/src/login.test.ts`·`mcp.ts`·`SKILL.md`·`adapters/README.md`·`scripts/install_bot_mcp.sh`는 UI 문구에 의존하는지 확인했다. 옛 버튼 이름을 쓰지 않는다. 옛 문구가 남은 곳은 D04 표(과거→제안)와 과거 증거뿐이다.

검증한 직접 확인 항목:

- CSRF: gate의 `csrf` hidden·`approve`/`deny` value·`formaction`이 그대로다. 회원 홈 form은 쿠키 `SameSite=Strict`와 `crossOrigin.Handler`(`http.go:189`)에 의존하며 이 코드는 diff에 없다. golden이 모든 form의 method·action·hidden·필드·button 라우팅을 고정한다.
- CSP: `http.go:193`·`capacity.go:169`의 헤더 문자열이 diff에 없다. 템플릿에 `<script`·`src=`·`<link`·`url(`·`@import`·`style=` 속성이 없다(테스트가 검사하고 브라우저에서도 0건).
- owner 분리: `homePage`의 `agents`·`pairs`·`gates`는 owner 비교로 걸러진다(변경 없음). 새로 추가한 건수(`ActiveAgents`·`IncomingInvites`·`PendingGates`)는 걸러진 목록만 센다. 다른 owner의 값이 새 화면에 닿는 경로를 찾지 못했다.
- 5분 재인증: 서버 검사 코드 변경 없음. 화면은 이미 계산된 `Recent`만 쓴다. `Recent`가 아니면 안내와 재확인 버튼이 먼저 보이고, 서버 거부는 그대로 남는다.
- 파괴 확인: 새 모달·확인 화면 없음. 위험 구역 `details`를 펼쳐야 단일 `POST` 버튼에 닿는다. 키보드로 확인했다(닫힌 상태에서는 Tab으로 닿지 않음).
- form action/method/hidden: golden을 재설계 전 템플릿(`37b85a4`)의 golden과 `diff`했다. 차이는 홈 `post /auth/logout` 2행(헤더 로그아웃, D04 4.2의 의도)과 원격 동의 fixture 2행(처음 pre에 없던 화면 추가)뿐이다.
- 출력 이스케이프: `badge`는 고정 맵과 `template.HTMLEscapeString`을 쓴다. 나머지는 `html/template` 문맥 이스케이프다. 연결 수단 token은 `<code>` 텍스트 노드로 출력된다. 사용자 값이 `template.HTML`로 승격되는 경로는 `badge`의 알 수 없는 상태 대체뿐이며 이스케이프된다.
- 파괴·데이터 손실: 상태 변경 로직·저장소·migration 변경 없음.

## 발견 사항

미해결 critical/high는 없다. 아래 3건은 모두 low이고 수락을 막지 않는다.

| # | 심각도 | 파일·줄 | 재현 조건 | 영향 | 수정 상태 | 검증 근거 |
|---|---|---|---|---|---|---|
| F-1 | low | `internal/relay/remote_auth.go:312` | fixture `remote-consent`·`remote-login` 로드, `document.querySelectorAll('main').length` | OAuth 동의·반환 화면에 `main` 랜드마크·건너뛰기 링크가 없다(A-02). 영문 "agent"도 남는다. DEV가 미변경을 기록했고 R-01 목록 밖이다 | 미해결(후속) | 87건 측정에서 두 화면 모두 overflow·44px 위반 0. `main`=0만 해당 |
| F-2 | low | `internal/relay/member_pages.go:114` | 연결 수단 `<code aria-labelledby="token-l">` | ARIA 1.2에서 `code` 역할은 이름을 받지 못해 일부 보조기기가 라벨을 읽지 않을 수 있다(A-07). 값 바로 위에 라벨 문장이 있어 정보 손실은 낮다 | 미해결(후속) | 코드 읽기. 실제 스크린 리더 미검증 |
| F-3 | low | `internal/relay/member_pages.go:141` | `home` 390px, nav y=89, 확인 필요 띠 y=199 | 960px 미만에서 섹션 이동이 확인 필요 띠 위에 놓인다. D04 4.2 순서와 다르다. 기능 영향은 낮다 | 미해결(designer 판단) | 브라우저 위치 측정 |

기록만 하는 관찰(발견 아님): 받은 초대의 수락·거절 옆에 "관계 철회"가 함께 보인다(form 계약 보존, DEV 공개). 홈 헤더·계정 섹션 로그아웃이 둘이다(공개). 서버 안내문에 영문 "agent"가 남는다(D04 허용, 정책 검사가 문장을 고정).

## 검증 및 남은 제약

- 등록 `kind: test` 명령 `make test`: Linux clone `c7828b2`에서 종료코드 0(race·adapters 포함). `make lint` 0. 같은 clone의 `lint.py`가 `lint.json`에 head `c7828b2`·base `8539ef5…`로 기록했다. 자기 실행이다. DEV 실행(`678736d`·`0b485ad`)과 구분한다.
- 표적 테스트 `TestMemberUIFixtures`·`TestMemberPagesShowNextSteps` PASS. 29개 합성 HTML 생성·브라우저 87건 측정(overflow 0, 44px 미만 0).
- integration·browser 태그 테스트(DB 필요)는 이번에 실행하지 않았다. `go vet`·`go test -c`로 컴파일만 확인했다. DEV의 `0b485ad` 통과 증거(`ok … 114.875s`)를 재사용한다. 근거: `0b485ad`→`c7828b2`의 제품 diff는 `member.go`의 템플릿 블록·`memberStyle`을 `member_pages.go`로 옮긴 것(사라진 줄 중 `html/template` 임포트 한 줄만 이동 대상이 아니다)과 `README.md` 2줄, `member_ui_test.go` 주석 1줄이다. 같은 바이트의 `c7828b2` 실행 결과가 아니므로 자기 실행으로 기록하지 않는다. 남은 위험은 낮다.
- 미검증: 실제 Google 로그인 왕복(실행하지 않음), 실제 스크린 리더, 200% 확대 직접 조작, 대비 픽셀 측정. 직접 시각 검수는 designer 후속이며 이 리뷰가 UI PASS를 선언하지 않는다.
- 보존: snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-c1644dfcf27644b3bc46f3c1788d9816`(읽기 전용, 정리는 `review.py cleanup` 대상). Linux clone `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2`. 삭제·이동 없음.
- 상세 QA 증거: [report.md](../SAR-SERVICE-UI-001-TESTER-test/report.md).

## 대화 미참조 인계 점검

이전 대화 없이 `.fullops-squad/docs/` 산출물 인덱스(`docs/deliverables/README.md`)와 phase·정본에서 시작해 확인했다. 경로는 `.fullops-squad/` 기준이다.

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `docs/design-docs/mockups/service-member-ui.md` 12.1 요구표, 6절 문구 표, 11절 불변식 | 확인 | D04 front matter `status: review`. 요구와 이유(O-01–12)가 문서 안에서 닫힌다. 이 문서의 12절이 R-12를 후속으로 둔다 |
| 구조와 구현/미완료 상태 | `docs/design-docs/module-design.md` 「회원 웹 상용 UI」 절, `docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` 변경·미검증 절 | 확인 | 파일 표가 실제 코드와 일치한다. 미완료(후속) 목록이 phase에 있다. 지침의 `SAR-PUBLIC-SERVICE-UX.md` 22행 "계정 비활성화"는 designer 소유로 갱신 필요(DEV 기록) |
| 실행·검증 방법과 증거 | DEV phase 「검증」·「시각 검수 환경」, 이 리뷰의 `SAR-SERVICE-UI-001-TESTER-test/report.md` | 확인 | 방법대로 29개 fixture를 새로 만들어 재현했다. 명령·위치·종료코드가 기록되어 있다 |
| 운영·복구(해당 시) | — | 해당 없음 | 서버 상태·배포·migration 변경이 없다. 운영 배포는 OPS 후속 |
| 다음 작업·담당·재개 조건 | DEV phase 「미검증과 후속」, 이 phase의 후속 | 확인 | designer 1440·390 직접 시각 검수 → main 통합·OPS 배포. 후속 low 3건은 coor가 DEV 또는 designer 과제로 정한다 |
| 로컬 링크·절 접근/지원 한계 | D04·D10의 상대 링크, phase의 `../../design-docs/…` | 확인 | 링크가 열린다. `lint.py`의 deliverables 검사가 ERROR를 내지 않았다. anchor 링크는 쓰지 않았다 |
| snapshot 정리 후 정본 접근 | 이 폴더(`result.json`·`report.md`·`lint.json`·`preview.json`·`rules.json`) | 확인 | 증거가 snapshot 밖 기록 체크아웃에 있다. snapshot 정리 뒤에도 `review.py check --historical`로 동일성 확인 가능하다 |

## 검토 결론

수락 가능(조건: 후속 low 3건과 designer 직접 시각 검수 기록). 미해결 critical/high는 없다. CSRF·CSP·owner 분리·5분 재인증·파괴 확인·form action/method/hidden 보존을 코드와 브라우저에서 직접 확인했다. 한계는 위 「검증 및 남은 제약」과 같다. DB integration은 DEV 증거 재사용이며 이 리뷰의 자기 실행이 아니다. check 통과는 기록 검사이며 AI 검토의 내용·테스트 성공을 자동 보증하지 않는다.
