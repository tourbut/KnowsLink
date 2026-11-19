---
title: FULLOPS-UPDATE-1.5.0
status: draft
updated: 2026-10-17
owner: coor
tasks: [FULLOPS-UPDATE-1.5.0]
summary: 1.5.0 대규모 변경만 독립 리뷰하는 기준과 snapshot 정리 수정을 적용하고 기존 코드·역할·검사는 보존한다
---

# FULLOPS-UPDATE-1.5.0

## 버전

업데이트 전 실제 설치 버전은 1.4.4(플러그인 캐시와 이전 기록)다. 이번 요청 전에 호스트가 1.5.0 캐시를 이미 준비했으므로 `claude plugin update`는 이미 최신(1.5.0)을 보고했다. 레포 적용 버전은 1.4.4에서 1.5.0으로 갱신한다.

## 설치

`claude plugin marketplace update fullops-squad` 성공. `deps.py --host claude-code` 설치 완료, `--check --host claude-code` 종료0(설치 직후 첫 check 1회는 영수증 갱신 중 실패해 재실행으로 통과). 이번 세션의 hook은 이전 경로일 수 있어 운영 업데이트만 수행하고 제품 worker는 시작하지 않는다.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|1.5.0 적용1 deps|적용. 위 설치 항목.|
|적용2 review_policy|보류 아님·기본값 채택. `fullops.json`에 review_policy를 두지 않아 기본 400줄·15파일 기준이 적용된다. 모든 변경 리뷰가 필요하면 기본 브랜치에 `{"min_code_lines":0,"auto":true}`를 사용자 결정으로 추가한다.|
|적용3 snapshot 정리|적용. `review.py prune` 결과 missing·live 모두 빈 목록. 미추적 snapshot-cleanup.json 3건은 사용자 기존 변경으로 보존.|
|적용4 rule.json exclude|해당 없음. review가 자동 skip하므로 필수 아님. 레포 rule.json 미변경.|
|적용5 .env.example|적용. TYPESAFE_API_KEY 줄 추가. 기존 .env 불변.|
|적용6 lint.json required|해당 없음. 기존 product-lint·product-test 명령 계약 보존, 새 필수 조건 소급 없음.|
|setup 재실행|기존 역할·remote·coor·lite·standard로 dry-run 후 실행. common/testing.md 1줄만 갱신.|
|선택 프런트 도구|해당 없음. 변경 없음.|

## 검증

`update.py --repo . --from 1.4.4`로 릴리스 1.5.0만 확인. setup 후 변경 파일은 testing.md·.env.example·fullops.json·본 기록이다. 제품 코드 변경 없음.
