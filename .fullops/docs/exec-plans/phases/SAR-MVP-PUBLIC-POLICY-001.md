---
title: SAR-MVP-PUBLIC-POLICY-001 — 공개 제품 기준 결정 기록
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-PUBLIC-POLICY-001]
summary: 공개 가입 확정과 DEC-03 제안값 및 승인·기술 근거 대기를 기록한다
---

# SAR-MVP-PUBLIC-POLICY-001 — 공개 제품 기준 결정 기록

## 범위와 적용 기준

designer는 기획 정본과 자기 기록만 수정했다. 작업 브랜치는 `tourbut/fullops-designer-pilot`, 착수 HEAD는 `3228d4e7178eae735aa04e6ab5c27e0db5609910`이다. 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이다. 공통 규칙 `fullops-common-0.3.2`, FULLOPS.md, project.md와 문서 작성 규칙을 적용했다. 받은 지시서의 필수 원천·D02·백로그·DEV/OPS 지시서·designer 컨텍스트를 대조했다. 준비 find/context 결과는 지시서의 기존 근거를 재사용했다.

fullops-work·fullops-deliverables로 문서와 완료 기록을 관리했다. orchestration으로 두 차례 coor에 ask했고 새 preamble의 Run `run_8ca8bc058ab7`, Task `task_9031cdaccb54`, Dispatch `ctx_44f5c365fa4f`로 회신한다. 기술 계획·제품 코드·배포·원천·D03·PLANS·board를 수정하지 않는다.

## 확정 결정과 승인 근거

사용자 배포 승인은 받은 지시서와 OPS 지시서에 기록돼 있다. 현재 서버 Docker·Cloudflare Tunnel 첫 배포와 `link.knowslog.com`을 적용한다. 원천의 `relay.knowslog.com`은 원문 보존 대상이며 현재 hostname과 구분한다.

첫 ask에서 coor는 사용자 결정을 전달했다. “누구나 가입하는 공개 서비스입니다. 초대 전용 파일럿으로 제한하지 마세요.” 가입 초대 제한을 추가하지 않았다. agent pairing 초대·B-human 수락은 기존 규칙대로 유지한다. 첫 단계는 인증된 등록·페어링·합성 안전 요청·human-gate 파일럿이다. 실데이터·실제 일정 효과·실벤더 연결·유료화는 제외한다.

[D02](../../planning/product-specs/SAR-MVP.md)에 공개 제품 기준·거부 의미·사용자 완료 조건을 작성했다. [백로그](../../planning/SAR-MVP-backlog.md)의 DEC-03과 구현 시작 승인 이력을 갱신했다. [D01](../../planning/business-plan.md)은 현재 배포 승인과 D02 연결을 반영했다. 기존 held와 원천 잠금 및 critical/high 차단은 유지한다.

## 제안값과 미정 결정

첫 ask에서 수치 결정을 사용자에게 모두 요구한 접근을 수정했다. coor는 기획자가 권장안·단위·적용 대상·거부 동작과 근거를 먼저 작성하도록 회신했다. 두 번째 ask에 초기 수용량·rate·32KiB 봉투·HTTP/claim 동시성 묶음을 제안했다. coor는 “제안 수치는 확인했으며 아직 승인값이 아닙니다.”라고 답했다. DEV/OPS에 polling/ACK·철회·경합·서버 자원 근거를 요청했다고 전달했다.

상세 수치는 D02 권장안이 정본이다. 수치는 두 agent 합성 흐름을 작은 초기 서비스 용량에서 검증하기 위한 제안이다. 측정치·확정값·성능 보장·Free N이 아니다. 공개 가입이 확정됐어도 실제 공개는 held다.

coor의 보완 요청에 따라 새 등록·초대·enqueue의 자원 상한을 기존 ACK·deny·철회·unpair·receipt replay에 재적용하지 않도록 정했다. rate 제안은 서비스 전체 신규 200회/안전 정리 100회, principal별 신규 40회/정리 20회로 분리했다. 모두 rolling 60s 기준이다. HTTP 동시성은 신규 16개/정리 4개로 분리했다. 신규 작업은 정리 budget을 소진하지 못한다. 정리 경로도 인증·현재 권한·lease·만료·CSRF·크기·남용 방어를 지킨다.

남은 제품 결정은 사용자 수치 묶음 승인/조정이다. DEV/OPS가 정상 polling·ACK·철회·경합·재시작·NAT·source IP·현재 서버 CPU/RAM/디스크·장기 보존 증가의 근거를 제공한다. 부모 요청과 H/R이 queue·receipt 한도에 걸려 안전 종료하지 못하는 문제와 예약량도 기술 근거를 받는다. designer는 근거를 반영해 한도와 거부 의미를 확정하고 coor는 구체적인 묶음만 사용자에게 질문한다.

DEC-01 Free N·가격·slot-unit과 DEC-02 disclosure/result schema는 held다. 무정책 query는 deny이며 optional result/error 공개·positive silent done은 보류다. DEC-04 실제 인터페이스와 DEC-05 기술 운영 설정은 해당 DEV/OPS가 맡는다. 정책 제안만으로 기존 DEV 구현 범위를 확대하지 않는다. 공개 한도 구현은 확정값과 coor 후속 인계 뒤 반영한다.

## 검증과 완료 경계

front matter는 deliverables.py --stamp로 작성한다. deliverables.py --strict, 변경 문서의 로컬 링크·승인/제안 구분·단위·held·파일 소유권과 git diff --check를 검사한다. 실제 결과는 아래에 기록한다. 검사 로그는 `/tmp/SAR-MVP-PUBLIC-POLICY-001-validation/`에 보존한다.

제품 코드가 바뀌지 않아 product-lint 게이트·build/test/runtime·독립 기능 QA·UI 캡처는 이번 문서 기록에 미적용이다. 제품 제한 구현과 독립 QA·직접 UI 검수·별도 세션 fixed-SHA 코드 리뷰·OPS 검증은 후속이다. 문서 준비 완료를 QA PASS·배포 수락·전체 MVP 완료로 표시하지 않는다.

## 문서 검사 결과와 인계

기획 문서 5개에서 로컬 링크 37개를 검사했고 오류 0, 종료코드 0이다. deliverables strict는 검사 13·미작성 10·문제 0·경고 0, 종료코드 0이다. git diff --check는 종료코드 0이다. 착수 HEAD 대비 코드·원천·기술 정본·PLANS·board 보존 검사도 종료코드 0이다. 로그는 위 레포 밖 검증 경로에 명령별로 보존했다.

기획 기록 단계는 완료한다. DEC-03은 공개 범위·거부 원칙의 부분 결정이며 추가 수치는 제안 상태다. DEV/OPS 근거 수신과 사용자 승인/조정 뒤 designer가 확정값을 반영한다. coor는 후속 제한 구현·독립 검증·수락·OPS 공개를 조정한다. 코드 변경이 없어 done-gate의 코드 변경 후 product-lint 검사는 미적용이다.

자기 인박스 완료 보고를 work.py finish로 날짜별 아카이브에 한 번 보존한다. 최종 커밋 후 strict·공백·보존·빈 인박스 검사를 다시 확인하고 고정 SHA를 worker_done에 전달한다.
