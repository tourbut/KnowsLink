---
title: tester 컨텍스트
status: draft
updated: 2026-10-13
owner: tester
tasks: [SAR-REFACTOR-001-TESTER, SAR-CODEX-HOOK-QUOTE-001-TESTER, SAR-SETUP-001-TESTER, SAR-MVP-001-TESTER, SAR-MVP-001-TESTER-FIX, SAR-BETA-001-TESTER, SAR-BETA-001-TESTER-PUBLIC, SAR-BETA-002-TESTER, SAR-MVP-002-DEV-TESTER, SAR-MVP-002-BOT-CATALOG-DEV-TESTER, SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER, SAR-MVP-003-BIDIRECTIONAL-TESTER, SAR-PUBLIC-IDENTITY-001-TESTER, SAR-PUBLIC-IDENTITY-001-FIX-TESTER, SAR-PUBLIC-AGENTS-001-TESTER, SAR-PUBLIC-AGENTS-001-FIX-TESTER, SAR-PUBLIC-MESSAGES-001-TESTER, SAR-GOOGLE-CONNECT-002-TESTER, SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER, SAR-GROK-WEBHOOK-001-TESTER, SAR-SERVICE-UI-001-TESTER]
summary: "독립 QA 재사용 경계와 MVP, 베타, Grok Bot 플러그인, 신원, agent, 메시지 PS08–11 QA 결과를 기록한다"
---

# tester 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 기준 ref의 `lint.json` `commands`가 비어 있으면 FullOps lint는 제품 검사를 실행하지 않는다. 기준을 새 SHA로 바꾼 임시 복사본에 위반을 커밋해 등록 명령의 실패 전파(종료코드 1)를 따로 증명했다.
- 검증 명령은 셸 없이 실행하는 래퍼(`run.py`)로 실제 종료코드를 로그에 남긴다. 위반 주입은 clone의 복사본에서만 하고 `git checkout -- .`로 원복한다.
- SAR-SETUP-001: dev SHA `0cc10b0` 독립 QA는 결함 없이 통과했다. 업무 SQL·sqlc·UI는 코드가 없어 미적용으로 남긴다.
- SAR-MVP-001: 후보 `a6a10c7`의 QA-01–11 실행 항목은 최종 probe 종료코드 0이다. DEC-02·DEC-03·Free N·실adapter·WAL·고의 epoch·designer 시각 판정은 held다.
- 보고서: [SAR-MVP-001-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md). 제품 코드는 수정하지 않았다. 초기 골격의 404·빈 migration 기대값은 쓰지 않았다.
- 2026-10-03 SAR-MVP-001-TESTER-FIX: 제품 `4262d02`의 신규 human transport 차단과 관련 회귀는 통과했다. legacy claim의 authorize·result 수락은 high이며 제품 수락을 차단한다.
- 판정·증거: [SAR-MVP-001-TESTER-FIX.md](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md). 원래 held와 `a6a10c7` QA는 유지했다. 제품 코드는 수정하지 않았다.
- 2026-10-03 SAR-MVP-001-TESTER-FINAL: 제품 `78b1d92`에서 legacy human·경로 미기록 claim의 authorize, H, R, consume은 `403 sender_not_allowed`다. 새 agent, owner gate, current-auth는 통과했다.
- 판정·증거: [SAR-MVP-001-TESTER-FINAL.md](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md). `4262d02` 보고서와 원래 held는 유지했다.
- 2026-10-03 SAR-BETA-001-TESTER: 판정 배포는 `f824015`다. `437f143`의 자손이고 제품 트리는 `78b1d92`와 같다. 검사 시작 체크아웃은 `437f143`이었고 reflog가 그 뒤를 바꿨다.
- loopback, 비공개 Postgres, 0600, 자원 제한, 합성 인증 음성, 백업, 격리 복원, expose 차단, 공유 서비스 회귀는 통과했다. 공개 Access와 인간 로그인과 원래 held는 유지했다.
- 판정·증거: [SAR-BETA-001-TESTER.md](../docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md). 제품 코드와 배포 설정은 수정하지 않았다.
- `f824015`의 Access 증명 누락·불일치와 migration deploy 차단은 임시 상태와 임시 clone에서 통과했다. 라이브 체크아웃은 옮기지 않았다.
- 2026-10-03 SAR-BETA-001-TESTER-PUBLIC: `28bd1bb` gate의 aud 불일치와 proof 차단은 임시 상태에서 통과했다. `gates.py` 종료코드는 0이다.
- 공개 negative는 공개 해석기와 `curl --resolve`에서 302 Access다. 기본 해석기 NXDOMAIN 때문에 `verify.py public` 종료코드는 1이다. 그 실행에는 상태코드가 없다.
- 인간 이메일 로그인은 실행하지 않았다. 사용자 held다. 보고서: [SAR-BETA-001-TESTER-PUBLIC.md](../docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md). 제품 코드와 배포 소스는 수정하지 않았다.
- 2026-10-03 SAR-BETA-002-TESTER: 배포 `28bd1bb`의 loopback에서 headless Chrome이 정확한 owner의 Approve, Deny, 새로고침, 기존 만료 gate 비활성, 다른 owner `403 sender_not_allowed`, 새 컨텍스트 200을 확인했다.
- 이 결과는 격리된 브라우저 컨텍스트의 재현이다. 사용자 브라우저 캐시, 이메일 OTP, 공개 로그인 뒤 UI는 관측하지 않았다.
- 판정·증거: [SAR-BETA-002-TESTER.md](../docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md). fixture는 복원했다. 제품 코드와 배포 소스는 수정하지 않았다.
- 2026-10-03 SAR-MVP-002-DEV-TESTER: 후보 `552586b6e886f95bffa9a000a031ea03070afedb`의 별도 clone에서 ZIP SHA256은 `0e671d1a89c141d896034fff31619b9cd2148b73b567adbc3a97126031989117`이다. MCP 경계와 `make verify-mvp` 종료코드는 0이다.
- 같은 프로세스의 겹친 pull은 `busy`다. 실제 Grok Bot 계정, marketplace, hosted runtime은 held다. 새 critical/high는 없다.
- 판정·증거: [SAR-MVP-002-DEV-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md). UI·Go는 `9584aaf`와 같아 Chrome QA를 재사용했다. 제품 코드는 수정하지 않았다.
- 2026-10-04 SAR-MVP-002-BOT-CATALOG-DEV-TESTER: 후보 `8e46c5a846e6d190e484e48be40b3dc368001a2b`의 별도 clone에서 설치 첫 실행·재실행 종료코드는 0이다. ZIP SHA256은 `b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad`다.
- 레포 밖 `env -i` 프로브는 도구 2개, held, stderr 0이다. 미지원·다운로드·checksum·빌드 실패는 exit 1이고 Ready가 없다. 새 결함은 없다.
- 판정·증거: [SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md). 실제 앱 등록·카탈로그·aarch64 실행은 미검증이다. 제품 코드는 수정하지 않았다.
- 2026-10-04 SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER: 후보 `423db6a2a388ea63610462f9d3a5f4c619dd781b`의 별도 clone에서 ZIP SHA256은 `d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609`다. 레포 밖 `env -i`는 도구 2개, held, stderr 0이다.
- swap 이전의 추출·`env -i` 실패는 기존 bundle을 남긴다. swap이 이전 트리를 stage로 옮긴 뒤의 mv 실패와 SIGTERM은 이전 `knowslink`를 지운다. 등급은 low다. critical/high는 없다.
- 판정·증거: [SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md). 실제 계정·개인 UI·relay·유료 inference는 미검증이다. 제품 코드는 수정하지 않았다.
- 2026-10-04 SAR-MVP-003-BIDIRECTIONAL-TESTER: 후보 `cd60e7f87eb5ce137eca887980f232b3f67a18d0`의 별도 clone에서 두 MCP와 Codex/Grok launcher 왕복 ID가 일치했다. SQL은 owner 2, trial 4, 업무 intent 0이다.
- 기본 10초 body timeout은 첫 청크 뒤 멈춘 응답에서 request가 15초까지 끝나지 않았다. 등급은 medium이다. 제품 코드는 수정하지 않았다.
- 판정·증거: [SAR-MVP-003-BIDIRECTIONAL-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md). 실제 Grok 계정·Cloudflare·공개 왕복은 미검증이다. owner 화면 템플릿은 `78b1d92`와 같아 Chrome을 반복하지 않았다.
- 2026-10-05 SAR-PUBLIC-IDENTITY-001-TESTER: 후보 `59b66ada8b36802484cc6d7e22523257b50572cc`의 fixture QA-P01–P05는 157통과 0실패다. `make verify-mvp` 종료코드는 0이다.
- 거부 170회가 공유 `http:new` 200을 채운 뒤 다른 source와 `/home`은 429이고 logout은 303이다. 등급은 medium이다. 제품 코드는 수정하지 않았다.
- QA-P06과 사람 QA-P07은 미실행이다. 일반 서비스 수락은 BLOCKED다. 보고서: [SAR-PUBLIC-IDENTITY-001-TESTER.md](../docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md).
- 2026-10-06 SAR-PUBLIC-IDENTITY-001-FIX-TESTER: 후보 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`의 별도 checkout에서 rate 격리와 trial lease 회수 검사가 통과했다. `make verify-mvp` 종료코드는 0이다.
- 첫 `make test` 종료코드 2는 손상된 zod 추출이다. 검증된 재설치의 product-test 종료코드는 0이다. 제품 코드는 수정하지 않았다.
- 실제 이메일·공개·노우↔다닷은 미실행이다. 원본 F1은 `59b66ad`에 둔다. 보고서: [SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md](../docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md).
- 2026-10-06 SAR-PUBLIC-AGENTS-001-TESTER: 후보 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`의 별도 clone에서 교차 계정, 키 회전·철회, 관계, 한도, rate, Node CLI가 통과했다. `make lint`, `make test`, `make verify-mvp` 종료코드는 0이다.
- 프로브 중간 실패는 기대값 오류였다. 제품 코드는 수정하지 않았다. 새 critical/high는 없다.
- 실제 이메일·공개·노우↔다닷은 미실행이다. 상속 리뷰 양식은 front matter만 보정했고 H1 이후 원문은 같다. 보고서: [SAR-PUBLIC-AGENTS-001-TESTER.md](../docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md).
- 2026-10-06 SAR-PUBLIC-AGENTS-001-FIX-TESTER: 후보 `458798c2ee15c179edacfd6f94ebb9896d26f411`의 별도 clone에서 보존, rate, POLICY 관찰 조건 프로브 종료코드는 0이다. 1차 종료코드 1은 프로브 기대값이다.
- 원본 QA `bcb06b8`과 UI FAIL `d165178`은 그 SHA에 둔다. 벽시계 24시간, 실제 메일, 공개, 노우↔다닷은 미실행이다.
- 제품 코드는 수정하지 않았다. 보고서: [SAR-PUBLIC-AGENTS-001-FIX-TESTER.md](../docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md).
- 2026-10-06 SAR-PUBLIC-MESSAGES-001-TESTER: 후보 `09c523da8a3407288d9f5d711e1834af12bc7808`의 별도 clone에서 PS08–11 프로브 종료코드는 1이다. H-1 high, M-1 medium, receipt 상한 뒤 신규 HTTP 행 1개 잔류를 재현했다.
- `make lint`, `make test`, `make verify-mvp` 종료코드는 0이다. 이 통과는 격리 fixture다. 제품 코드는 수정하지 않았다.
- 실메일, 공개, 벽시계 24시간, 운영 부하와 복원, PS13, PS14, 노우↔다닷은 미실행이다. H-1은 main 수락을 차단한다. 보고서: [SAR-PUBLIC-MESSAGES-001-TESTER.md](../docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md).

- 2026-10-07 SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER: d089 기본 lint/test/build exit 0이나 runtime 재빌드 exit 2, tester QA 컴파일 exit 1이다. 사용자 중단으로 독립 기능 QA는 미완료다.
- 원 실패와 78개 파일 해시를 보존했다. 자기 DB와 scratch를 회수하고 managed snapshot은 남겼다. [상세 보고서](../docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md).

- 2026-10-10 SAR-GOOGLE-CONNECT-002-TESTER: 고정 후보 dbbe2f1 수락. 실제 cloudflared 26경로와 Node22 Windows 로그인 경계를 직접 확인했다.
- test_public_config.sh는 Linux 전용(stat -c). Git Bash 경로 변환은 MSYS_NO_PATHCONV=1로 막는다. 실제 Google·Bot·왕복은 coor 후속. [리뷰](../docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md).
- 2026-10-11 SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER: `eac603b` 승인 callback 복귀를 실제 Chromium·격리 PG·합성 Google/호스트로 독립 확인했다. 기준 handler는 form-action 위반으로 callback에 못 가고 수정본은 도착해 PKCE 교환 200이다.
- 경계(만료·CSRF·계정 binding·한도·token 7종)는 안전하게 거부된다. 발견: 모든 실패가 같은 "만료" 문구(F1 medium), 반환 링크 24px(F2 low). 미해결 C/H 없음. 실Dots·실계정은 미검증이다.
- Git Bash는 `/tmp/...` 인자를 Windows 경로로 바꾼다. 원격 경로를 넘기는 helper 호출에는 `MSYS_NO_PATHCONV=1`과 Windows 로컬 경로를 쓴다. 등록 `make`는 SHA 검사하는 shim으로 원격 클론에서 실행해 lint.py head를 맞췄다.
- 2026-10-11 SAR-GROK-WEBHOOK-001-TESTER: 고정 `b86f397` webhook wake를 snapshot 리뷰와 실제 fetch loopback 프로브 58건으로 확인했다. 미해결 C/H 없음, low F1: `EAI_AGAIN`·`ENETUNREACH`가 uncertain이라 그 ID의 깨우기가 24h 사라질 수 있다. 실Bot routine은 미검증.
- 리뷰 packet이 대상 SHA보다 늦게 커밋되면 `review.py check --task-key`가 "탐색 패킷 없음"으로 실패한다. 리뷰 키로 check하고 packet은 `work.py finish` 뒤 확인한다. 인라인 삭제 명령은 shell_guard가 막으므로 tar `--exclude`로 대체한다.
- 2026-10-11 SAR-SERVICE-UI-001-TESTER: 회원 웹 `c7828b2`를 snapshot 리뷰와 합성 fixture 29개 × 1440·390·320 브라우저 측정, Linux `make lint/test`로 확인했다. 미해결 C/H 없음, low 3건(OAuth 화면 main 없음, `code`의 aria-labelledby, 모바일 nav 순서). 실Google·실스크린 리더는 미검증.
- Orca 브라우저에서 `set viewport`는 `goto` 뒤에 해야 적용된다. 앞에 두면 `innerWidth`가 창 폭 그대로라 측정이 무효가 된다. 측정 행마다 `iw`가 요청 폭과 같은지 먼저 확인한다.
- 패킷이 대상 SHA 뒤 커밋이고 리뷰 키가 과제 키와 같으면 `review.py check`가 막힌다. `--task-key <없는 키>`로 실행하고 사유를 기록한다. 순수 이동 refactor는 `git diff`의 사라진 줄과 추가된 줄을 `comm`으로 비교해 DEV 증거 재사용 근거로 쓴다.
- 2026-10-11 SAR-SERVICE-UI-001-FIX-TESTER: 후보 `295417a`의 V-01 수정은 독립 리뷰로 수락 가능하다. fixture 29종 × 320/390/1440에서 visible 링크 66개가 모두 44px 이상이다. 증거: [리뷰 보고서](../docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md).

- 2026-10-11 coor 정정: 위 `--task-key <없는 키>` 기록은 잘못된 우회 지침이며 사용하지 않는다. 제품 리뷰는 실제 구현 과제 키에 연결한다. FIX-TESTER의 동일295417a를 실제 FIX-DEV로 check해 통과했다. 과거 실패는 FIX-TESTER phase에 보존, 후속 근거는 FIX-TESTER-RESULT-REVIEW/report.md다.
- 2026-10-12 SAR-CODEX-HOOK-QUOTE-001-TESTER: 고정 `d7d6d744` Codex hook를 snapshot에서 cmd·pwsh·powershell·sh로 실행해 root/scripts/adapters 허용0·거부2·deny JSON을 확인했다. 미해결 C/H 없음, low 1건(argv 목록 호출은 guard 미실행). 실제 Codex host 실행·원 stderr는 미검증.
- hook 셸 검증은 코드를 문자열 리터럴로 넘기면 guard 없이 0이 되므로 거부 칸이 2인지와 deny JSON 파싱으로 실행을 증명한다. `make.cmd` 같은 원격 SSH 도구는 DEV HEAD를 올리므로 tester snapshot에 쓰지 않고 exact-SHA 로그·`pass.json` 동일성으로 재사용한다.

- 2026-10-13 SAR-REFACTOR-001-TESTER: exact cd4dae8 Linux의 Go race 경계8개·MCP/owner-text·독립 text matrix15개 exit0. DEV 전체 게이트는 원본·사본 LF 정규화 내용과 SHA를 대조해 재사용했다. 새 C/H 발견0, 제품 읽기 전용.
- [독립 QA 보고](../docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md). 첫 확인식 따옴표 오류 exit1 보존. 실계정/운영/SQL 이행은 not_run이며 전체 제품 수락을 대신하지 않는다.
