---
title: SAR-PUBLISH-001-P03 DEV 검증 보고
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-PUBLISH-001-P03]
summary: 개인정보처리방침 백업 문장을 14일 보관으로 교체하고 원격 lint·test를 통과시킴
---

# SAR-PUBLISH-001-P03 DEV 검증 보고

코드 SHA f15fceac(기준 d5695ff1). 변경은 문장 교체뿐이다. 화면 구조·스타일·시행일(P-01)은 그대로다.

| 검사 | 명령·위치 | SHA | 결과 |
|---|---|---|---|
| 원격 lint | `remote.py make lint`(Linux clone) | f15fceac | EXIT=0 |
| 원격 test | `remote.py make test`(Linux clone) | f15fceac | EXIT=0 |
| 영어판 한글 | `grep -cP` 한글 범위, `privacy.en.html` | f15fceac | 0건 |

정책 문안 일치 테스트는 이 문장을 직접 비교하지 않아 고칠 테스트가 없었다(grep 확인).

검증하지 못한 것: 운영 `/privacy` 실서버 화면(coor가 ops 정리와 함께 배포 후 확인). 브라우저 렌더 캡처는 문장 교체라 만들지 않았다.
