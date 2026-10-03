---
id: D11
title: 사용자설명서
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-BETA-001-OPS]
summary: 본인 전용 합성 베타의 접속·합성 시험·종료 방법을 안내한다
---

# KnowsLink 본인 전용 합성 베타 사용자설명서 (D11)

이 문서는 사용자 본인이 합성 데이터로 베타를 시험하는 방법이다. 서버는 `link.knowslog.com`이다. 접근은 사용자 본인 이메일 한 개로 제한한다. 실제 업무 데이터·실제 벤더·실제 일정은 입력하지 않는다.
공개 연결은 2026-10-03에 적용됐다. 미인증 요청은 모두 Cloudflare Access 로그인으로 이동한다. 마지막 이메일 로그인은 사용자가 직접 확인하는 인간 검사다. 현재 상태는 [D13](transition.md)에서 확인한다.

## 시험 절차

1. 서버 터미널에서 합성 시험 자료를 만든다.
   `/home/shin/deploy/knowslink/deploy/knowslink/beta.sh seed`
   기대: `UI candidate: http://127.0.0.1:8080/owner/gates/<id>`가 출력된다. 이 gate는 생성 후 180초에 만료된다.
2. 같은 터미널에서 로그인 값을 확인한다.
   `/home/shin/deploy/knowslink/deploy/knowslink/beta.sh owner-login`
   기대: `username`과 `basic-auth value`가 출력된다. 합성 값이다. 다른 곳에 붙여 넣지 않는다. 이 명령은 사용자 본인이 자기 터미널에서만 실행한다. agent에게 실행시키지 않는다.
3. 브라우저에서 `https://link.knowslog.com/owner/gates/<id>`를 연다. `<id>`는 1단계 출력에 있다.
4. Cloudflare Access 화면에서 본인 이메일을 입력한다. 이메일로 온 일회용 코드를 입력한다. 기대: 소유자 화면의 HTTP Basic 로그인 창이 나온다.
5. 브라우저 로그인 창의 사용자 이름 칸에 2단계의 `username` 값을 입력한다. 비밀번호 칸에 `basic-auth value` 값을 입력한다. 기대: pending gate 내용이 나온다.
6. Approve 또는 Deny를 누른다. 기대: 결정 결과 화면이 나온다. 페이지를 여는 GET 요청만으로는 결정되지 않는다.
7. 만료된 gate는 1단계를 다시 실행한다.

## 사용하지 않는 것

- 본인 이메일이 아닌 계정은 접근할 수 없다. Access가 차단한다.
- agent의 API 호출은 서버 안의 loopback(`127.0.0.1:8080`)에서만 실행한다. 공개 주소에서 API 자동 호출은 지원하지 않는다.
- 실제 정보 공개·벤더 효과·positive silent done은 꺼져 있다.

## 종료 방법

| 목적 | 명령 | 결과 |
|---|---|---|
| 외부 접근만 중지 | `/home/shin/deploy/knowslink/deploy/knowslink/beta.sh unexpose` | 앱과 데이터는 유지된다. 공개 주소는 오류가 된다 |
| 베타 전체 중지 | `/home/shin/deploy/knowslink/deploy/knowslink/beta.sh stop` | 컨테이너를 멈춘다. 데이터 볼륨은 유지된다 |
| 데이터 백업 | `/home/shin/deploy/knowslink/deploy/knowslink/beta.sh backup` | 0600 dump가 `/home/shin/deploy/knowslink-state/backups/`에 생긴다 |

완전 철회는 [D12](ops-guide.md) 11.3을 따른다. 다른 프로젝트(`myportfolio`)와 기존 `orca` Tunnel은 이 명령으로 바뀌지 않는다.
