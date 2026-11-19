---
title: SAR-HUMAN-PEER-001 고정 SHA 독립 코드 리뷰
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-HUMAN-PEER-001-TESTER]
summary: 254ff6a 사람 친구 후보의 권한·발신·보존·이행을 별도 세션에서 검토하고 low 문구 3건을 남긴다
---

# SAR-HUMAN-PEER-001-TESTER 리뷰

- 검토자 / CLI / 모델: 세션 `01a12fd5-0d90-7803-85ed-d7546b4f7189`. 호스트는 Grok 4.7이다. dispatch가 high를 지정했다. 작성 세션은 `01a12f82-ca28-7de1-8f74-acbfa1a6c823`이다.
- base SHA / head SHA / merge-base: `f849215bdefd3984dbf385c9f0973a74525be545` / `254ff6a98b12b7e51f0835491300f7695425d264` / 같은 base.
- OCR 버전 / 적용 규칙: open-code-review v1.12.12. 호스트가 리뷰했다. OCR LLM 모드는 사용하지 않았다. 규칙 파일은 `.fullops-squad/review/rule.json`이다. 공통 규칙은 `fullops-common-0.3.3`의 README, coding-style, testing, security다. 기준 커밋은 지시서의 `f20b33a11addeb3debcb23e2e835be31c7dd1454`다.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` HP-01–16, 현행 D04 `service-member-ui.md`, D03 architecture, D11 user-guide.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 243 / reviewed 88 / skipped 155. pending은 0이다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0, WARNING 17, 실행 불가 2. Windows 호스트에 `make`가 없다. 사유를 기록했다. 이 상태를 passed로 바꾸지 않았다.
- SIZE-002: 추가 2577줄. 상한은 400줄이다. 사람 친구 권한, 웹 owner 발신, 보존, legacy 이행이 한 후보에 있다. 기능을 나누면 인가 경계가 갈라진다. 분할하지 않은 이유를 수락한다. 파일 길이 SIZE-001은 경고이며 동작 결함이 아니다.
- DEP-001: `adapters/package.json`의 test 스크립트에 `node dist/owner-text.test.js`가 추가됐다. 새 패키지 의존성은 없다. 표준 라이브러리 대안은 해당 없다.
- UI 디자인: 정본은 현행 D04다. 공용 시작점은 기존 Go template와 memberStyle이다. 다크 테마는 범위 밖이다. Go 문자열 CSS를 디자인 lint가 판정하지 못하는 한계를 유지한다. 390/1440 시각 판정은 designer 몫이다. 이 리뷰는 화면 수락을 대체하지 않는다.

snapshot은 `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-cba0ca432fa64d08903c989aa30d9464`이다. detached, clean, HEAD `254ff6a98b12b7e51f0835491300f7695425d264`, read_only다. 결과 파일은 이 기록 체크아웃에만 썼다. snapshot worktree는 삭제하지 않았다.

## 검토 범위

변경 243개 경로에 reviewed 또는 skipped와 사유가 있다. 권한 호출, 친구와 자기 연결, legacy 정확 pair, generation, 해제와 차단, 웹 owner와 agent 서명 발신, 부모 TTL과 분리된 새 발송, 24h 참여 사본, 멱등, quota, CSRF, writer rollback, adapter, SSE, remote, gate를 snapshot에서 읽었다.

skipped 155개는 과거 QA 로그, 생성 탐색, 핸드오버 이력이다. 각 파일의 사유는 그 경로를 이번 SHA의 수락 증거로 재명명하지 않는다는 뜻이다. 존재는 같은 폴더의 요약과 이 보고서가 구분한다.

## 발견 사항

미해결 critical과 high는 없다. 아래 3건은 문구 공백이며 동작 한계와 일치한다.

1. low, 미해결. `scripts/package_plugin.py` 146행. 패키징된 문구가 `at most24 hours`다. 보존은 24시간이다.
2. low, 미해결. `adapters/skills/knowslink/SKILL.md` 42–48행. `at most180 seconds`와 `until24 hours`다. 큐 TTL은 180초이고 보존은 24시간이다.
3. low, 미해결. `.fullops-squad/docs/operations/user-guide.md` 29행. 이번 diff가 추가한 문장이 `뺀1–40자`다. 서버 한계는 1–40 Unicode rune이다. 재현은 `git show 254ff6a98b12b7e51f0835491300f7695425d264:.fullops-squad/docs/operations/user-guide.md`의 29행이다.

HTML `maxlength=80`은 UTF-16과 rune의 차이다. 서버가 40 rune을 거부한다. 결함으로 올리지 않았다. 대화 헤더가 정렬 전 첫 원소를 쓰는 점은 같은 쌍의 두 방향에서 이름이 바뀐다. 라이브에서 잘못된 상대 이름이 고정되는 재현은 없다. 결함으로 올리지 않았다.

## 검증 및 남은 제약

`review.py check --task-key SAR-HUMAN-PEER-001-DEV`는 문구 3건을 반영한 뒤에도 통과했다. reviewed 88, skipped 155, WARNING 17이다. check는 기록 검사다. 테스트 성공을 자동 보증하지 않는다.

Windows `lint.py`의 product-lint와 product-test는 unavailable이다. Linux clone `254ff6a`에서 `make lint` EXIT 0, `make test` EXIT 0이다. 로그는 `C:\Users\shin\AppData\Local\KnowsLinkDevTools\humanpeer001-runtime-ctx9ded\1791806105370100800-254ff6a98b12b7e51f0835491300f7695425d264-make_remote`다. `test_public_config.sh`의 부정 경로 문장은 스위트 실패가 아니다.

격리 DB `humanpeer_integration_f46538af714c`에서 `TestHumanPostgres`와 `TestHumanLegacyPostgres`가 통과했다. 소요 3.622초다. 이 필터 실행은 패키지 전체가 아니다. 패키지 전체는 `make test` EXIT 0을 사용한다.

독립 HTTP QA는 QA 보고서가 정본이다. DEV `runtime/report.md`와 이전 PNG SHA `4eaefa6`, `59d419d`, `7f83fd6`, `f567abf`는 과거 캡처다. 이번 리뷰 통과의 대체 증거가 아니다.

jev_find 변경 파일 겹침은 code recall 0.333 / precision 0.556, documents recall 0.111 / precision 0.167이다. 수락 판단에 쓰지 않는다.

SEC-001과 SLOP-004는 테스트 fixture와 테스트 출력이다. 제품 비밀값 하드코딩으로 보지 않았다. SLOP-002는 `text.ts`의 조건부 스프레드이며 owner-text 파싱 분기다. 동작 결함으로 올리지 않았다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | SAR-HUMAN-PEER.md HP-01–16 | 확인 | 과거 agent별 수락은 현행 HP가 대체한다 |
| 구조와 구현/미완료 상태 | architecture.md, DEV phase, SHA 254ff6a | 확인 | 운영 적용은 후속이다 |
| 실행·검증 방법과 증거 | final-report.md와 이 검수 QA 보고서 | 확인 | runtime/report.md는 과거 이력이다 |
| 운영·복구 | database-design.md writer 차단 | 확인 | 라이브 goose Down은 실행하지 않았다 |
| 다음 작업·담당·재개 조건 | designer 시각 판정, low 문구 3건 | 확인 | 제품 수정은 DEV가 한다 |
| 로컬 링크·절 접근 | user-guide의 DEV phase 링크 | 확인 | 깨진 anchor는 이번 수락 차단 사유가 아니다 |
| snapshot 정리 후 정본 접근 | snapshot은 아직 active | 확인 | 삭제·이동 금지로 cleanup을 실행하지 않았다 |

## 검토 결론

코드 리뷰 기록은 수락할 수 있다. 미해결 critical/high가 없다. 권한, owner와 서명 발신, 멱등, CSRF, 이행 fail-closed, 24h 참여 열람이 정본과 맞다. low 문구 3건은 남긴다. 이 결론은 제품 전체 수락이 아니다. designer의 직접 시각 판정이 남아 있다.
