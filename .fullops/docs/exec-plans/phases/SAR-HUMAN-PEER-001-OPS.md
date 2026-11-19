---
title: SAR-HUMAN-PEER-001-OPS 운영 적용 실행 기록
status: draft
updated: 2026-10-12
owner: ops
tasks: [SAR-HUMAN-PEER-001-OPS]
summary: 수락한 사람 peer·UI 26094e9를 backup·격리 복원·migration 00002와 함께 기존 서버에 적용한 계획·순서·결과와 coor 후속
---

# SAR-HUMAN-PEER-001-OPS

## 범위와 기준

배포 제품 SHA는 `26094e90648aa71e8d07b381e9891acb37438730`이다. 준비·기록 main은 `3425267`이며 제품 파일은 같다. 지시서 attempt는 `f7ff48dfbefc43d2ab71fae18ad90c57`이다. 규약은 `fullops-common-0.3.3`, `project.md`, 문서 작성 규칙, AGENTS의 Windows 파일 작업 안전이다. 테스트 레벨 lite, 하위 위임 standard다. 서버 작업은 순서 의존이 있어 하위 에이전트를 쓰지 않았다. 제품 코드·정책·migration·compose·ingress는 바꾸지 않았다. 증거 정본은 [운영 LIVE 기록](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-OPS/report.md)이다.

| 구분 | 근거 |
|---|---|
| 이전 운영 | `95df83c`, [SERVICE-UI LIVE 기록](../../evaluations/qa-reports/SAR-SERVICE-UI-001-OPS-LIVE.md) |
| 독립 QA·리뷰 | [TESTER 리뷰](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-review/report.md), [UI delta 리뷰](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-UI-DELTA-review/report.md) |
| 시각 수락 | [VISUAL 최종 보고](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md) |
| 기술 인계 | [D07](../../design-docs/database-design.md#sar-human-peer-001-dev-후보), [DEV 기록](SAR-HUMAN-PEER-001-DEV.md) |

## 지시서 전제와의 충돌 처리

과거 OPS의 "새 backup 없음"은 migration이 없던 과제의 제한이었다. 이번에는 migration `00002_human_writer`가 있어 새 backup이 필수다. `beta.sh`의 `deploy`·`backup`·`restore-verify`는 실패 시 `rm`을 포함하고 `deploy`는 서버에 없는 origin을 fetch하므로 호출하지 않았다. 삭제 없는 `pg_dump`(noclobber, 새 고유 파일)와 고유 이름의 자동 소멸 격리 컨테이너로 같은 순서를 수동 실행했다. `beta.sh up`은 내용을 읽고 `up` 항목만 사용했다(build·migrate 완료 후 relay `--wait`, 삭제 없음). `verify.py regression`의 baseline `projects` 결함은 새 실패와 구분해 원 종료코드 1을 보존했다.

## 기술 계획과 복구 계획

migration `00002`의 trigger는 `HumanModel>=1`인 행을 더 낮은 값으로 쓰는 UPDATE를 거부한다. 구 relay는 `HumanModel` 필드가 없는 JSON을 저장하므로 새 relay의 첫 변환 저장 뒤에는 쓰기가 거부된다. 그래서 구 writer로 되돌리지 않는다. Down은 의도적으로 오류를 낸다. 복구는 전진 방식이다: relay 정지→현재 DB 새 backup→현재 schema·친구·차단·세대·철회·키를 보존하는 reader/writer로 진행한다. 과거 backup 복원은 권한을 되살릴 수 있어 금지한다. 위험은 변환 순간(pending pair 만료·`SelfConnect=false` 설정)뿐이다. baseline에서 pending pair가 0건임을 확인했다.

## 순서와 결과

1. 읽기 전용 baseline: 서버 HEAD `95df83c` clean, 상태 해시·권한, relay_state 키별 개수와 기존 필드 md5, 컨테이너·포트·cloudflared·공유 서비스.
2. `95df83c..26094e9` 인프라 diff는 migration 한 파일뿐이다. bundle(SHA256 `7a92999f…6dfe`)을 전달·검증해 `git fetch <bundle> main`으로 받았다. `deploy-history.log`에 한 줄을 추가했다.
3. detached checkout, `compose config --quiet`, relay·migrate 이미지 build(기존 relay 서비스 중).
4. relay 정지 12:43:44Z, epoch 고정 확인 후 `pg_dump -Fc`(0600, 11098 bytes), `pg_restore --list` 종료 0.
5. 격리 복원 리허설: 복원본 Agents md5가 라이브와 일치, migration Up·첫 쓰기 성공·downgrade 쓰기 거부. 첫 시도의 파일 권한 실패(Up 미적용)는 로그로 보존했다.
6. `beta.sh up` 종료 0, goose 2, trigger 1, relay healthy, `verify.py local` 종료 0.
7. 전·후 보존 대조: 기존 필드 md5 동일, 신규 `Name`·`Client`(빈 문자열)·`HumanModel=1`·`FriendCode`·`SelfConnect=false`만 추가.
8. 무인증 공개 smoke: 새 디자인 토큰과 새 h1, 새 `/home/*` 경로의 무세션 303, public deny 404, 공유 서비스 불변.

## 남은 일 (coor)

1. 실제 Google 로그인 탭에서 `https://link.knowslog.com/home`을 연다. 운영 SHA는 `26094e90648aa71e8d07b381e9891acb37438730`이다. 두 기존 회원 모두 이름이 비어 있다. 첫 화면은 `/home/name`의 「내 이름」이다. 안내 문구는 「이 이름은 친구에게 보입니다. 이름은 1–40자로 입력하세요.」이고 버튼은 「저장하고 대화로 이동」이다. 저장 뒤 대화 화면으로 이동한다.
2. 표시 이름 선택, agent 명칭 선택, 친구 코드 등록은 coor 후속이다. 이 과제는 실제 친구에게 요청이나 메시지를 보내지 않았다.
3. 기존 회원의 `SelfConnect=false`다. 자기 에이전트끼리의 연결은 새 UI에서 한 번 동의해야 한다. 기존 pair 1개(회원2)는 「기존 연결」에서 원래 두 agent만 유지한다.
4. 키보드 포커스·200% 확대·모바일의 실계정 확인은 미실행이다.
5. `verify.py` baseline 결함(`projects`)과 `public` 검사의 Access 302 전제는 이번에도 고치지 않았다. 담당은 dev다.
