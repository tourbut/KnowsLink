---
title: SAR-PUBLISH-001-TESTER 독립 QA 보고
status: review
updated: 2026-10-22
owner: tester
tasks: [SAR-PUBLISH-001-TESTER]
summary: "main 17a360b2의 캘린더 제거·일회 정리·/privacy·/terms를 새 격리 fixture에서 독립 확인한 결과. 62개 fixture 확인과 회귀 게이트 모두 PASS, 제품 결함 0"
---

# SAR-PUBLISH-001-TESTER 독립 QA 보고

tester 독립 확인(PC-06의 tester 몫)이다. designer 시각 수락·ops 운영 반영 확인은 이 보고의 범위가 아니다. 정책 문안은 법률 검토(P-05) 전 초안이다. 테스트 레벨 lite.

- 대상 제품: main `17a360b2051d96e39646e1b048c931f6f0ede88b`(워크트리 HEAD `2ad0c6bb`는 `.fullops-squad/` 기록만 더한다: `git diff --stat 17a360b2 HEAD`에 제품 파일 없음).
- 실행 환경: 원격 Linux 격리 clone(`/tmp/knowslink-pubqa-17a360b2051d-*`), 자체 Postgres 컨테이너, 합성 데이터, loopback만. DEV 인스턴스는 재사용하지 않았다. 운영 서버·env·Tunnel·Cloud Console·실계정·실메일은 쓰지 않았다. 종료 때 fixture를 `POST /fixture/stop`(204)으로 끝내고 컨테이너를 내렸다(남은 컨테이너 0).
- 도구: [remote.py](remote.py)(DEV 것을 이 과제용으로 고정 SHA·경로만 바꿈), [explore.py](explore.py), [walk.py](walk.py), [show.py](show.py), [stop.py](stop.py), [measure.py](measure.py). 시나리오: [SAR-PUBLISH-001-TESTER.md](../../scenarios/SAR-PUBLISH-001-TESTER.md).

## 결과 요약

| 항목 | 판정 | 증거 |
|---|---|---|
| 정리: 키 있음 | PASS | `calendar cleanup: connections=2 revoked=2 failed=0`, 연결 0·lookup·disclosure·키 0, 결과 gate `approved→denied`, 로컬 revoke 요청 2 |
| 정리: 키 있음 + `revoke=fail` | PASS | `connections=2 revoked=0 failed=2`, 토큰·연결·미전달 결과 0, gate `denied`, revoke 요청 2(전부 거부됨) |
| 정리: 키 없음 | PASS | `connections=2 revoked=0 failed=2`, 요청 0건, 삭제만 수행, gate `denied`. failed=2는 '철회하지 못한 연결 수'로 읽힌다 |
| 정리 기록에 토큰·이메일 값 없음 | PASS | 세 기록 모두 `날짜 시각 calendar cleanup: connections=N revoked=M failed=K` 형식뿐, `@`와 24자 이상 토큰 문자열 0 |
| 재시작 idempotent | PASS | 세 사례 모두 같은 정리를 두 번째 실행하면 기록 빈 문자열, revoke 요청 증가 0, 상태 불변 |
| 설정·승인 상세 KO/EN 캘린더 문구 없음 | PASS | 정리 뒤 account·pending gate·정리된 gate 6개 화면에서 `calendar`·`캘린더`·`busy time`·`바쁜 시간` 0 |
| `schedule.query` 승인·disclosure 없음 | PASS | 승인 버튼 있음, 승인 POST 받아짐, 승인 뒤 버튼 없음, consume 응답 `disclosure:false`. 정리된 gate consume은 `invalid_gate` |
| `schedule.commit` 기존 동작 | PASS | integration `TestPostgresSafety` PASS(browser fixture에는 commit gate가 없어 화면 왕복은 not_run) |
| `/home/calendar/connect`·`/disconnect` | PASS | KO·EN 세션 POST·GET 모두 404, 비로그인 POST도 404 |
| Google 로그인 scope `openid email` | PASS(코드·단언) | `internal/relay/google.go:50` `Scopes: {openid, email}`, 제품 코드 전체에서 offline·`calendar.freebusy` 옵션 0(grep), 단언 `google_integration_test.go:92`(`scope == "openid email"`, `access_type` 없음) PASS. fixture는 `/auth/google`을 막고 Google 설정이 없어 redirect URL 직접 확인은 not_run |
| `/privacy`·`/terms` KO/EN 비로그인 200 | PASS | 4개 URL 200, 로그인 회원도 200. `TestPolicyPagesPublic` PASS |
| 문안이 정본과 같다 | PASS | 새 fixture에서 받은 렌더 4개와 `SAR-PUBLISH-001.md`의 4개 문안: 단어 단위 일치(523/772/303/442), [parity.txt](parity.txt). 비교 스크립트는 DEV의 `parity.py`를 재사용했다(독립 구현 아님) |
| 영어판 한국어 0·`{ }` 0 | PASS | EN 두 페이지에서 한글 0('한국어' 단어 제외), 4개 페이지 모두 `{`·`}` 0 |
| 첫페이지 하단·설정·상호 링크 | PASS | 첫페이지 하단 두 링크(KO·EN), 설정 두 링크, `/privacy`↔`/terms` 상호 링크, 고유 title·h1 1개 |
| `#start` 정책 안내 한 줄 | PASS(단위) | `landing_pages.go:70`은 Google 또는 이메일 로그인이 켜질 때만 렌더한다. fixture에는 로그인 설정이 없어 한 줄이 없는 것이 정상이고 `TestPolicyLinks` PASS가 맡는다. 운영에서 이 줄이 보이는지는 ops의 live 확인 몫 |
| 공개 ingress | 후보만 확인 | `check_public_ingress.py`는 `make lint`에 포함되어 EXIT=0. live 반영은 ops 몫 |
| 3폭 가로 넘침 | PASS | 정책 4·첫페이지 2·설정 2·승인 상세 2 = 10 페이지 × 1440/390/320에서 `scrollWidth ≤ innerWidth`, BAD `[]`([widths.json](widths.json)). 영어 위탁 표만 320px에서 표 영역 안 스크롤(304>273, 페이지 넘침 아님) |
| 친구·메시지 송수신·승인 회귀 | PASS | integration(`-race -tags=integration`, `TestRemoteConsentBrowserReturn` 제외) EXIT=0, `--- FAIL` 0, `--- PASS` 128 |

## 실행한 명령과 종료코드

| 명령 | 결과 | 로그 |
|---|---|---|
| `python remote.py prepare`(clone·checkout `17a360b2`·`make install`) | EXIT=0 | 원격 |
| `python remote.py fixture` | EXIT=0 | 원격 |
| `python explore.py walk.py`(fixture 확인 62개) | EXIT=0, 62/62 PASS | [walk-result.json](walk-result.json) |
| `python remote.py make lint` | EXIT=0 | [lint](17a360b2051d96e39646e1b048c931f6f0ede88b-lint.log) |
| `python remote.py make test` | EXIT=0 | [test](17a360b2051d96e39646e1b048c931f6f0ede88b-test.log) |
| `python remote.py integration` | EXIT=0, FAIL 0 | [integration](17a360b2051d96e39646e1b048c931f6f0ede88b-integration.log) |
| `python parity.py` / `python measure.py` | EXIT=0 / EXIT=0 | [parity.txt](parity.txt), [widths.json](widths.json) |

## 관찰과 검증하지 못한 것

- 시행 방식 기록: 첫 `walk.py` 실행들(탐색·스크립트 오류 수정)은 같은 fixture를 이미 정리·승인한 뒤 돌려 판정이 아니다. 판정은 새로 띄운 fixture에서 한 번 돌린 62/62다. 처음에는 키 없음 사례의 기대값을 `failed=0`으로 잘못 적었다(검사 스크립트 오류). 제품은 `failed=2`로 기록하고 요청은 0건이다. 이 동작은 DEV 보고의 '요청 없이 삭제만'과 같다. 기록 문구의 의미는 coor·ops가 ops-guide의 설명과 맞는지 확인하면 된다.
- `gate_pending`(lookup이 없는 pending gate)은 정리 뒤에도 `pending`이다. 지시(lookup이 있고 consume되지 않은 gate만 `denied`)와 U5 기준(승인 가능)에 맞다. 결함이 아니다.
- 리뷰 low 'commit~revoke 사이 crash 시 정리 기록 줄 없음': fixture에 그 사이를 끊는 지점이 없어 재현하지 않았다. 코드 구조상(삭제 commit 뒤 revoke 루프·기록)이라는 리뷰 설명만 있고 이번 QA의 관찰은 없다. ops-guide 절차 보강은 리뷰 권고 그대로 남는다.
- 실제 Google revoke 끝점·운영 상태 정리·live 반영·Cloud Console은 검증하지 않았다(ops·coor 몫).
- 실제 브라우저의 키보드·200% 확대·시각 판정은 하지 않았다(designer, SAR-PUBLISH-001-VISUAL).
- 정책 문안의 법적 적합성은 보지 않았다(P-05).
- 제품 결함: 없음. 제품 코드는 수정하지 않았다.
