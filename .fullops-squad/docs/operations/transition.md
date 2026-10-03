---
id: D13
title: 인수인계서
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-BETA-001-OPS]
summary: 베타 배포의 현재 상태와 인수 항목 및 남은 일을 기록한다
---

# KnowsLink 본인 전용 합성 베타 인수인계서 (D13)

이 문서는 베타 배포의 인수인계 상태를 기록한다. 배포 구성과 절차 정본은 [D12](ops-guide.md) 11장이다. 사용자 시험 절차는 [D11](user-guide.md)이다.

## 현재 상태 (이 문서 작성 시점)

| 항목 | 상태 |
|---|---|
| 배포 SHA | 체크아웃 `/home/shin/deploy/knowslink`는 detached다. 제품 코드는 수락 제품 `78b1d92`와 동일하다. 현재 SHA는 `git -C /home/shin/deploy/knowslink rev-parse HEAD`와 `knowslink-state/deploy-history.log`로 확인한다. 갱신·rollback은 D12 11.3 |
| 로컬 스택 | 기동·검증 완료. relay `127.0.0.1:8080`, Postgres 비게시 |
| Tunnel | `knowslink` Tunnel 생성 완료. connector 미기동 |
| Access 앱·정책 | **미생성** — Access 쓰기 권한이 필요하다 |
| DNS `link.knowslog.com` | **없음** — 보호 확인 전 연결하지 않는다 |
| 공개 노출 | **held** |
| 사용자 이메일 로그인 확인 | 인간 검사. 미실행 |

## 인수 항목

- 상태 디렉터리 `/home/shin/deploy/knowslink-state`(0700): `.env`, Tunnel 자격 파일, `access.aud`, 백업. 모두 Git 미추적이다. 자격 파일을 공유 위치로 옮기지 않는다.
- 공유 서비스 baseline은 `shared-baseline.json`에 있다. `verify.py regression`으로 비교한다.
- 첫 DB 백업과 격리 복원 검증은 완료했다.
- 보존 대상: `myportfolio` 프로젝트·볼륨, `orca` Tunnel, 호스트 cloudflared PID 506937, `~/.cloudflared/config.yml`.

## 남은 일

1. Access 앱 생성 권한을 연결한다. 최소 권한은 Account 범위의 `Access: Apps and Policies Edit`와 `Access: Organizations, Identity Providers, and Groups Read` 두 개다. 둘 중 하나가 빠지면 `apply`가 HTTP 403으로 멈춘다(fail-closed).
2. 독립 리뷰와 coordinator main 통합 후 D12 11.2 순서로 노출한다.
3. 사용자가 본인 이메일로 로그인해 시험한다.
4. held 유지: DEC-03 공개 한도, 실제 신원, 실데이터, 실제 벤더, 무제한 공개.
