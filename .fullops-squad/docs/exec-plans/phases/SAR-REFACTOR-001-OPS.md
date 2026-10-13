---
title: 수락 리팩토링 운영 적용 실행 기록
status: draft
updated: 2026-10-13
owner: ops
tasks: [SAR-REFACTOR-001-OPS]
summary: 고정 제품의 무삭제 배포·보존 경계·실패 증거·클라이언트 후속을 연결한다
---

# SAR-REFACTOR-001-OPS 실행 기록

사용자가 승인한 수락 제품 `2747fe98653d5c6946f41c219fc256395333d7fe`를 운영에 적용했다. 기존 제품은 `26094e90648aa71e8d07b381e9891acb37438730`다. 이 과제는 코드 개발이 아니라 승인된 배포·보존 검증·클라이언트 안내·기록이다.

## 기준·기술 판단

공통 `fullops-common-0.3.3`, project.md, AGENTS.md, D07 writer guard와 D12/D13을 적용했다. 테스트 레벨 lite이며 직접 수행했다. 하위 위임은 명시 지시에 따라 사용하지 않았다. 서버 origin 부재·beta deploy/backup의 rm 때문에 기존 SSH helper와 bundle 수동 절차를 재사용했다. known_hosts·RejectPolicy·자격 메모리 보관·비밀 비출력을 유지했다.

26094e9..2747fe9의 실제 제품 diff는 Go/TS 코드·회귀 11개 파일이다. schema·DB·배포 설정·의존성·UI template/CSS·remote MCP schema 변경은 없다. CODE/FINAL 리뷰와 TESTER 고정 제품 cd4dae8의 Linux 전체 게이트 근거를 읽고 재사용했다. 새 전체 제품 QA 실행으로 표시하지 않는다.

## 수행 순서와 결과

1. 서버 HEAD/clean·health·private mode·goose2·guard1·HumanModel1과 보존 집계를 읽었다. 기존 실제 회원/친구·Grok 로컬 설치·Dots remote app/Events 기록을 확인했다.
2. 고정 main bundle hash·verify·FETCH_HEAD·조상 관계·schema/deploy diff0을 확인했다. detached checkout 뒤 relay/migrate를 먼저 build했다.
3. relay만 중지했다. 2초간 epoch912460 고정 뒤 새0600 pg_dump를 만들었다. 백업은 `26094e9-pre-refactor-20261013T132105Z.dump`, 15693bytes, SHA256 `792523eebdb6d007fd6f01cf45ebc75b282942dd96eabaf60424602b74cc95bd`다.
4. pg_restore list25줄·network none 격리 복원·writer downgrade 거부를 확인했다. 새 격리 컨테이너는 정상 종료 후 보존했다. live DB 복원·새 SQL migration·삭제/이동은 없다.
5. beta up·local 검사를 통과했다. 중단부터 healthy 확인까지13:21:02Z~13:21:18Z 약16초다. 수동 deploy-history를 append했다.
6. 전후 schema·회원/agent/친구/키·메시지·Remote/Inbox/Events·설정·기존 백업·공유 컨테이너 digest가 동일했다. 공개 smoke21개·공유200/200/401·약2분 새 relay error/panic0·up 이후 tunnel error0을 확인했다.
7. Grok의 같은 기존 연결 폴더/agent/key/Command·watcher·webhook routine을 유지하는 bundle 갱신·재시작 안내를 작성했다. Dots는 기존 remote app/구독을 유지하며 metadata가 바뀔 때만 refresh한다. 둘 다 직접 업데이트·실계정 발송은 하지 않았다.

## 실패·한계·복구

관측의 old HEAD 조건 오류와 업로드 CRLF 오류는 relay 중단 전에 멈췄다. 실패 자료를 보존하고 고유 LF 사본으로 수정했다. 설명되지 않은 제품 실패는 없다. 전체 과거 문서 전문 재독은 하지 않았다. packet의 unknown/remaining 후보는 존재·hash·관련 검색 구간과 최신 정본·직접 서버 보존 경계로 수동 보완했다.

Events1개는 배포 전후 같은 만료·Verified·attempts/pending/lease metadata로 보존됐다. 새 이벤트·retry·다음 구독 갱신·장시간 host turn은 실행하지 않았다. 실계정 송수신/Google 로그인·클라이언트 직접 갱신·UI/장시간/부하는 not_run이다. 기존 held와 독립 QA/시각 수락은 그대로다.

복구는 현재 DB·친구·키·철회·Remote/Inbox를 보존하고 같은 HumanModel writer를 가진 직전 수락26094e9 코드로만 한다. 과거 DB 덮기·구 writer·Down은 금지한다. rollback은 이번에 실행하지 않았다.

상세 실행·검증·공식 근거·원 자료 보존 경로는 [운영 보고](../../evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md)에 있다. D12/D13에 현재 상태와 후속을 연결했다. coor는 완료 문서 SHA의 독립 snapshot 리뷰·main 통합을 담당한다. 최종 기록 SHA는 worker_done으로 보고하며 운영 제품 SHA와 구분한다.
