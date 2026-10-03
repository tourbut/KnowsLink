---
title: SAR-BETA-002-TESTER — 루프백 소유자 브라우저 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-002-TESTER]
summary: "배포 28bd1bb의 loopback Chrome에서 소유자 승인, 거절, 만료, 권한 불일치를 기록한다"
---

# SAR-BETA-002-TESTER — 루프백 소유자 브라우저 독립 QA 보고서

## 판정

판정 배포 SHA는 `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`다. 시작 HEAD와 끝 HEAD가 같다. porcelain은 시작과 끝이 0줄이다.
제품 경로와 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 diff는 비어 있다. 제품 코드와 배포 소스는 수정하지 않았다.
판정 실행은 `browser.mjs`의 두 번째 실행이다. 종료코드는 0이다. checks 26개가 통과했다.
headless Google Chrome 148.0.7778.96이 정확한 owner의 dashboard, fresh gate, Approve, Deny, 새로고침을 확인했다.
이전 seed owner와 같은 seed의 다른 owner는 새 gate에서 403 `sender_not_allowed`다. 맞는 owner의 새 컨텍스트는 200이다.
fixture owner의 기존 expired gate 1개는 읽기 전용으로 버튼이 없다. POST는 0개다.
부모 completion은 승인 뒤와 거절 뒤 모두 빈 값이다. 화면의 정책 문장과 비실행 안내 문장은 결정 뒤에도 남는다.
이 판정은 사용자 브라우저 캐시의 관측이 아니다. 이메일 OTP와 공개 로그인 뒤 UI는 실행하지 않았다. 새 critical/high는 없다.
시나리오: [SAR-BETA-002-TESTER.md](../scenarios/SAR-BETA-002-TESTER.md). 로그: [SAR-BETA-002-TESTER-test/](SAR-BETA-002-TESTER-test/).

## 기준

날짜는 2026-10-03이다. 기록 브랜치는 `fullops/tester`다. 공통 기준은 `fullops-common-0.3.2`다.
정본은 `.fullops-squad/project.md`다. 문서 작성 규칙과 `fullops-test`를 적용했다. 예외는 없다.
Jev keep은 `internal/relay/http.go`, `adapters/src/synthetic.ts`, `deploy/knowslink/beta.sh`, 공개 QA, 로컬 QA, 제품 FINAL QA, user-guide, transition, project, 공통 README와 세 규칙, FULLOPS, 현재 인박스다.
충돌 후보와 주의 후보는 없다. find의 absent는 신규 제품 구현 부재 추천이다. 대상 파일은 고정 SHA에 있다.
D01–D13은 읽기만 했다. 갱신하지 않았다.

## 환경

| 항목 | 값 |
|---|---|
| 브라우저 | Google Chrome 148.0.7778.96, headless |
| 조작 도구 | 저장소 밖 `playwright-core` 1.63.0. 경로 `/tmp/sar-beta-002-browser` |
| Node | v22.22.2 |
| 주소 | `http://127.0.0.1:8080` |
| 배포 작업 트리 | clean. HEAD 고정 |
| trace, HAR, 스크린샷 | 만들지 않음 |

저장소에 Playwright 의존성을 추가하지 않았다. DOM과 버튼 상태는 테스트 코드가 판정했다. Jev 웹 조작은 실행하지 않았다.
`beta.sh owner-login`은 실행하지 않았다. credential은 `build/qa-fixture.json`에서 메모리로만 읽었다.

## 브라우저 결과

| 검사 | 결과 |
|---|---|
| browser_launched | 통과. headless |
| seed_approve `01a101ac-33dc-7a5d-8dd8-cb355e30c41d` | 종료코드 0 |
| approve_dashboard, approve_fresh | 200. typed body, intent, 정책, 만료, 안내 문장, 버튼 2개 |
| approve_decided, approve_reload | POST 303. 상태 `approved`. 버튼 0개. 새로고침 유지 |
| approve_no_schedule_effect | completion 빈 값. intent `schedule.query` |
| seed_deny `01a101ac-36f9-778e-928a-0bd5c77d0e84` | 종료코드 0 |
| previous_owner_403 | 403 `sender_not_allowed`. 원문과 버튼 없음 |
| other_owner_403 | 403 `sender_not_allowed`. 원문과 버튼 없음 |
| fresh_context_200 | 200. `상태: pending` |
| deny_dashboard, deny_fresh, deny_decided, deny_reload | 승인 흐름과 같은 표시. 상태는 `denied` |
| deny_no_schedule_effect | completion 빈 값 |
| expiry_existing `01a1019b-eed5-70d4-9e94-6a4245c02485` | 200. `상태: expired`. 버튼 0개. POST 0개 |

기존 expired gate가 fixture owner에게 1개 있었다. 180초 대기 분기는 실행하지 않았다. 시계와 라이브 DB 시각은 바꾸지 않았다.
승인 화면과 거절 화면은 `정보 공개나 일정 실행을 허용하지 않습니다.`와 `deny even after gate approve`를 유지했다.
adapter의 completion `denied` 경로는 다시 실행하지 않았다. 그 경로는 기존 합성 검증의 범위다. 이번 클릭은 completion을 만들지 않았다.
403 재현은 새 Playwright 컨텍스트에 이전 owner 자격 또는 같은 seed의 다른 owner 자격을 넣은 것이다. 사용자 Chrome 프로필은 열지 않았다. `user_cache_observed`는 false다.

## 보존

| 항목 | 실행 전 | 실행 후 |
|---|---|---|
| fixture SHA-256 | `d5c6e55306a5be50c961018400315ce1be2a56fc3dac07e83e7cf276eed7d3db` | 같음 |
| fixture 모드 | 600 | 600 |
| `relay_state` 행 | 1 | 1 |
| `relay_state` 바이트 | 37510 | 48856 |
| 배포 HEAD | `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2` | 같음 |
| 배포 porcelain | 0줄 | 0줄 |

바이트 증가는 합성 owner, agent, gate 추가다. 기존 행을 지우거나 fixture를 교체한 채 남기지 않았다.
첫 실행은 브라우저 기동 전에 종료코드 1이었다. Postgres가 `count(*)`와 `octet_length(data)`를 한 SELECT에서 거부했다. seed는 실행되지 않았다. fixture 해시와 배포 트리는 그대로였다. 그 로그는 판정 실행이 덮어썼다. 판정은 두 번째 실행이다.

## 재사용한 증거

| 증거 | SHA | 이번 실행 |
|---|---|---|
| 공개 Access negative | `c993d599efab0bfdc5741bc9cb02053ca9afd856` | 재사용. 다시 실행하지 않음 |
| 로컬 런타임 | `1762b430bed1c0584fecd163ae81567a4a5d04a9` | 재사용. 판정 배포는 `f824015` |
| 제품 QA | `659f4b06dfda1bd59997f8f2b06cba23b48026a6` | 재사용. 제품 `78b1d92` |
| 배포 리뷰 | `7ba9df046611c67109b53ad2b43812543b937f95` | 재사용 |
| 제품 리뷰 | `311381feb0f18203f60329605d762257bfff8421` | 재사용 |

coor의 loopback HTTP 승인·거절 기록은 서버 직접 호출이다. 이번 브라우저 QA의 통과로 기록하지 않았다.
인간 이메일 OTP는 실행하지 않았다. 담당은 사용자이고 상태는 held다.
모바일 실기기와 공개 주소 로그인 뒤 UI는 실행하지 않았다. loopback headless Chrome 결과와 구분한다.

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `node .fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER-test/browser.mjs` | 0 | [browser.log](SAR-BETA-002-TESTER-test/browser.log) |
| 같은 명령의 stderr | 빈 출력 | [browser.err](SAR-BETA-002-TESTER-test/browser.err) |

`make test`, `make verify-mvp`, 백업, 복원, `verify.py public`은 실행하지 않았다. 제품 트리와 해당 증거의 실행 조건이 같기 때문이다.
healthz 200은 브라우저 검사의 전제다. UI 통과의 대체 증거가 아니다.

## 한계

headless Chrome은 사용자 데스크톱의 보이는 창과 사용자 프로필이 아니다. 버튼과 DOM 판정은 이 엔진에서 수행했다.
합성 seed가 만든 owner와 gate는 라이브 DB에 추가된 채 남는다. 자격 파일은 이전 바이트로 복원했다.
사용자에게 보인 `sender_not_allowed`의 원인으로 이 재현을 단정하지 않는다. 재현은 이전 owner 자격이 새 gate에서 403이 되고, 맞는 새 컨텍스트가 200이 된다는 것이다.
