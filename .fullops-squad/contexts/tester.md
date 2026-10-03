---
title: tester 컨텍스트
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-SETUP-001-TESTER, SAR-MVP-001-TESTER, SAR-MVP-001-TESTER-FIX, SAR-BETA-001-TESTER, SAR-BETA-001-TESTER-PUBLIC, SAR-BETA-002-TESTER, SAR-MVP-002-DEV-TESTER, SAR-MVP-002-BOT-CATALOG-DEV-TESTER, SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER]
summary: "독립 QA 재사용 경계와 MVP, 베타, Grok Bot 플러그인 검증 결과를 기록한다"
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
