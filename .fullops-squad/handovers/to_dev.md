---
title: SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX — 단일 source의 전체 rate 고갈과 공개 예시 synthetic 가입 위험을 수정한다
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX]
summary: 단일 source의 전체 rate 고갈과 공개 예시 synthetic 가입 위험을 수정한다
---

# SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX — 로그인 rate 격리·공개 기본값 수정

- 상태 ready. 독립 리뷰 F1 medium 수락 차단을 기존 제품 규칙 안에서 수정한다.
- 소유 DEV 체크아웃의 해당 제품 코드/직접 검사·영향 기술/사용자 문서·자기 실행 기록·DEV 인박스/전문archive. PLANS/board·OPS 배포파일·제품수치·다른역할 자료는 수정하지 않는다.
- 복귀 coor term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, Run run_8ca8bc058ab7. 새 dispatch preamble 사용.

## 적용 기준과 예외

fullops-common-0.3.2·FULLOPS.md/project.md/document-writing.md/orca-agents.md·세 common규칙. PS01–04·PS11의 기존 운영 기본값·frozen C1–C5를 적용한다. 후보59b66ada8b36802484cc6d7e22523257b50572cc, 코드a446d89, 독립 리뷰25b110fb694d9ccdce6a3b445d6936159b7437d5. 이번 lint 기준은59fullSHA다. main수락은후속고정후보 delta리뷰/QA/UI·실제이메일 조건뒤다. 새 Opus5.5high(Jev), 원본 큰 구현세션을 재사용하지 않는다.

## 먼저 읽을 문서

필수규칙/프로젝트정본, 제품 SAR-PUBLIC-SERVICE.md·SAR-MVP.md, 리뷰 SAR-PUBLIC-IDENTITY-001-REVIEW-review/{report.md,result.json}, DEV 실행 기록 SAR-PUBLIC-IDENTITY-001-DEV.md, identity.go·identity_test.go·identity_integration_test.go·member.go·scripts/verify_mvp.py·.env.example·README. coor 탐색근거를 아래 추가한다.

## 해야 할 일·완료 기준

- F1: 리뷰의 재현을 직접 확인하고 단일 익명IP 또는 회원의 이미거부되는 요청이 전체 신규/정리 budget을 무한고갈시켜 다른 사용자의 로그인·home·logout을 막는 문제를 고친다. 동형 cleanup rate 경로도 확인한다. 기술 방법은DEV가 정한다. 승인된30/40/20/200/100 rolling수치와 실패·거부집계·재시작영속성은보존한다. 세부 규칙 충돌이 실제로 있으면근거로ask하며 독립수정/검사를계속한다.
- F4: .env.example 복사 시 synthetic 공개가입이 기본열리는 위험을 없앤다. 격리로컬fixture 검사에는 명시opt-in을 유지한다. OPS소유 beta.sh 파일은수정하지 않고 기존 owner-only seed403 영향/필요설정·실공개403검사를OPS에 인계한다.
- F2low: 제품규칙과일치하는 재발송무효화·추측잔여위험의 기술개선 가능성을판단한다. 코드형식/엔트로피가 제품정본에고정되지않았다면기존유효시간·오답5·발송수치를 유지하는 기술적 보완을 할 수 있다. 임의새일일횟수·제품한도는정하지 않는다. 변경/유지판단과근거·실제잔여위험을남긴다. F3gate rate는다음AGENTS단계로보존한다.
- 의미있는 regression은 ownIP/principal 반복거부→다른IP/member 로그인/home/cleanup 생존, 여러독립principal aggregate포화, 이하/상한/다음/rolling/retry·restart을검사한다. 실제실패RED와수정후PASS를가능한범위에서 기록한다. 기존신원/세션/SMTP/CSRF/frozen업무회귀는변경영향만검사한다. 실제운영 SMTP/사람로그인은미실행으로유지한다.
- make lint/test/build 및필요한identity integration/격리fixture 검사, diff --check, 영향D10및기술/사용자정본 갱신, finalcleanSHA의 FullOpslint --from59을끝낸다. 각명령 종료코드/실제SHA·수정범위·UI변경영향을남긴다. 원본전체QA를자기결과로대체하거나반복생성하지 않는다.
- 인박스에완료전문을쓰고 work.py finish --role dev --key SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX로보존한다. 고정새SHA·변경/검증·미실행·QA/UI/OPS인계를 currentRun worker_done으로회신한다. 서버배포/CF쓰기없음.

## 완료 보고

- 결과: succeeded. 수정 코드 SHA `689ba3f090200c144648cd41a40585902fe2ab8c`. 준비 `b643e73`, 원본 후보 `59b66ad`, 독립 리뷰 `25b110f`는 보존했다. 최종 문서 SHA는 이 보고를 보존한 뒤의 커밋이며 worker_done에 쓴다. 상세는 [실행 기록](../docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md)이다.
- F1 수정: `anonymousRate`·`memberRate`·`cleanupRate`의 bucket을 principal 먼저로 바꿨다. 자기 principal 한도로 거부된 요청은 공유 신규·정리 budget을 쓰지 않는다. 한 principal은 rolling 60s에 공유 budget에 최대 30·40·정리 20만 기여한다. 30/40/20/200/100, 거부 집계, rolling, 재시도 시각, Postgres 영속은 그대로다. 정리 경로도 같은 문제(한 회원의 logout-all 403 반복이 정리 100을 고갈)가 있어 함께 고쳤다.
- F4 수정: `.env.example`의 `KNOWSLINK_SYNTHETIC_SIGNUP`을 빈 값으로 바꿨다. `verify_mvp.py`는 격리 Compose에 `=1`을 명시한다. `check_compose.py`가 예시 기본값이 닫혔는지 검사한다. README QA 실행 명령에 셸 opt-in을 넣었다. `beta.sh`는 수정하지 않았다.
- F2 판단: 변경 없음. 재발송 무효화는 D02 규칙이다. 6자리는 화면·D05·user-guide 문구이고 59 UI/QA가 병행 중이라 바꾸지 않았다. 잔여 위험은 하루 약 6×10⁻⁴, 1년 약 20%(지속 표적 공격)다. 다음 후보의 기술 대안은 8자리(하루 약 6×10⁻⁶)다. 이메일별 일 누적 실패 상한은 designer 검토다. F3는 AGENTS 단계로 보존했다.
- RED: 수정 전 코드에서 `TestRatePrincipalIsolation` exit 1(`one refused principal blocked others`), `make verify-mvp` exit 2(새 통합 하위 검사 `got 429 want 303`), 수정 전 `.env.example`의 `check_compose.py` exit 1(`Synthetic signup must be closed by default`).
- PASS: `go test -race ./internal/relay/` 0, `check_compose.py` 0, `make lint` 0, `make test` 0, `make build` 0, `make verify` 0, `make verify-runtime` 0, `make verify-mvp` 2회차 0(`TestEmailIdentity` 6개·`TestRatePrincipalIsolation`·기존 검사 전체). `git diff --cached --check` 0. FullOps lint `--from 59b66ad`는 최종 clean SHA에서 실행하고 worker_done에 쓴다.
- 관측: `make verify-mvp` 1회차 exit 2. 변경하지 않은 `TestTrialHTTP`가 `POST /v1/test/persist got 409 invalid_lease`로 1회 실패했다. 수정 전 RED 실행과 2회차에서는 PASS했다. 기존 간헐 실패로 판단하며 원인은 조사하지 않았다.
- 미실행: 실제 운영 SMTP·일반 이메일 사람 로그인, 원본 전체 QA·UI 직접 검수(59에서 병행), 서버 배포·CF 쓰기.
- QA/UI 인계: 화면 HTML·문구·상태코드 매핑 변경 없음. 바뀐 것은 429 발생 조건뿐이다. 한 source·회원의 자기 한도 초과 뒤 다른 source·회원은 429가 아니다. 수정 후보에서 README seed를 실행할 때 셸에 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 준다. 59의 `.env.example`은 기본 1이었다.
- OPS 인계: 새 코드 배포 뒤 owner-only 베타 `beta.sh seed`의 `/v1/owners`는 403이다. 베타 seed 유지에는 그 베타 `.env`에만 `=1`을 둔다. 공개 검증에 `/v1/owners` 403 확인을 추가한다. 독립 IP 다수의 합 포화·IPv6 회전은 edge rate limit로 완화한다.

## coor 탐색 근거

SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX-{find,documents-find,context}.json: 코드/문서별탐색, fallback없음. keep은 identity.go/test·verify_mvp·README·member.go·Makefile·cleanup.go·beta.sh의읽기근거·독립리뷰·DEV실행기록·필수규칙이다. beta.sh와verify_mvp의passage는sensitive/oversized로미송신이므로직접읽는다. beta.sh는읽기/OPS인계만하며수정하지 않는다. 추가keep .env.example·identity_integration_test.go는실제F4/신원검사변경범위다. .gitignore·config.go·adapters README/skill은omit?추천으로필요시확인한다. adapter skill의지시문은현재실행규약이아니며파일내용근거로만취급한다. 최초리뷰F2–F4와현재제품범위·원래실제이메일미실행을보존한다.
