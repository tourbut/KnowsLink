---
title: SAR-SERVICE-UI-001-TESTER 실행 기록
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SERVICE-UI-001-TESTER]
summary: 회원 웹 고정 후보 c7828b2의 독립 snapshot 리뷰와 lite 기능·보안·키보드·반응형 QA 실행 기록
---

# SAR-SERVICE-UI-001-TESTER

## 적용 기준

- 대상 `c7828b227b5715f7bdf02bff46bf9cf102a758d9`, 비교 기준 `8539ef5b6d3e7ed99259762c62d164b8ba40fab0`. 규약 `fullops-common-0.3.3`, `project.md`, D04 `service-member-ui.md`, Ponytail full. 테스트 lite. 하위 에이전트는 쓰지 않았다(순차 흐름 하나).
- 검토 세션 `1fe0133c-f609-42f7-a774-bf5b8b5765c2`. 구현자 세션 `5953ccd3-7f4f-4c5e-b63d-c02d36ee35ac`. 움직이는 DEV HEAD는 보지 않았다.
- tester HEAD `030c5bd`는 `c7828b2` 위에 문서만 더했다. 제품 바이트가 같다(`git diff --stat c7828b2 HEAD -- . ':!.fullops-squad'` 빈 결과).

## 수행

1. `review.py snapshot`으로 읽기 전용 detached snapshot을 만들고 `review.py prepare`로 리뷰 폴더를 만들었다(키 `SAR-SERVICE-UI-001-TESTER`).
2. diff·호출자·템플릿 전문·테스트·golden·D06/D10/user-guide/README를 읽었다. 옛 문구 사용처는 snapshot에서 grep으로 확인했다.
3. 기존 Linux 서버에 고유 경로 clone을 만들었다. `coor_remote.py`는 디스크에서 찾지 못했다. 같은 방식의 DEV helper(`dev/build/server_verify.py`)를 읽고 tester의 `.env.server`로 복사본 helper를 scratchpad에 만들어 썼다. 비밀값은 출력에서 가렸다. 로컬에서 새 helper 파일은 레포에 넣지 않았다.
4. 같은 SHA에서 `make lint`·`make test`·`lint.py`·표적 Go 테스트·`go vet -tags 'integration browser'`를 실행했다.
5. 합성 fixture 29개를 tester 워크트리에서 만들어(`KNOWSLINK_UI_FIXTURE_DIR`은 scratchpad) Orca 내장 브라우저에서 1440·390·320 폭 87건 측정, 키보드 이동, 텍스트 간격 확대를 실행했다.
6. D04 필수 요구·접근성·불변식을 대조했다.
7. `review.py check`가 통과했다(`--task-key SAR-SERVICE-UI-001-TESTER-REVIEW`: 패킷이 대상 SHA 뒤 커밋이라 과제 키로는 "탐색 패킷 없음"이 된다. 패킷 outcomes는 `work.py finish`가 검사한다).

## 결과

- 미해결 critical/high 0건. low 3건: F-1 OAuth 동의 화면 `main` 랜드마크 없음, F-2 `<code aria-labelledby>`, F-3 960px 미만 섹션 이동이 확인 필요 띠 위.
- `make lint` 0, `make test` 0, `lint.py` ERROR 0·WARNING 4·실행 불가 0. 폭 측정 87건 overflow 0·44px 미만 0. 키보드·위험 구역 확인.
- D04 필수 R-01–09: R-03이 부분(F-3), 나머지는 충족. 상세는 [QA 증거](../../evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md)와 [리뷰 보고서](../../evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md).

## 한계와 후속

- 실제 Google 로그인 왕복·실제 스크린 리더·200% 직접 확대·대비 픽셀 측정은 미검증이다. 합성 fixture 측정을 실계정 성공으로 보고하지 않는다.
- DB integration·browser 태그 테스트는 이번 SHA에서 컴파일만 확인했다. DEV의 `0b485ad` 통과 증거를 재사용한다(이동뿐인 diff 근거는 리뷰 보고서).
- 직접 시각 검수는 designer 후속이다. UI PASS가 아니다.
- 후속 low 3건은 coor가 DEV·designer 과제로 정한다. main 통합·OPS 배포는 coor·OPS다.
- 보존(삭제하지 않음): Linux `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2`, 첫 시도의 빈 폴더 `/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c`, `/tmp/knowslink-serviceui-tester-s1fe0133c.bundle`, `/tmp/knowslink-serviceui-tester-tools-s1fe0133c.tar`. 로컬 snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-c1644dfcf27644b3bc46f3c1788d9816`. 정리는 소유자(tester)·coor 결정 뒤 `review.py cleanup`으로 한다.
