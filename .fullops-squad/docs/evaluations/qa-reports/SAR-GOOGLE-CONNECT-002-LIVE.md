---
title: 비용 없는 Google 연결 운영 적용
status: draft
updated: 2026-10-10
owner: coor
tasks: [SAR-GOOGLE-CONNECT-002-DEV, SAR-GOOGLE-CONNECT-002-TESTER]
summary: 독립 수락 후보 배포와 Tunnel 차단 및 Access 앱 제거의 실제 결과
---

# 운영 적용 결과

DEV dbbe2f17f353282430f28b6e3ed98e72d4450e69와 tester 98f08a96e474de24558c8f21e4b0e7168418ce9b를 main ba3754d2c41d3ba8258eae76fa0cf0cfad5b27c4에 통합하고 origin/main에 일반 push했다. 두 결과는 로컬·원격 main의 조상이다. 고정 snapshot의 review check는 reviewed28/skipped9, critical/high0, lint WARNING0으로 통과했다. 별도 구현자·검토자 실제 세션은 리뷰 정본에 있다.

운영 서버에 Git bundle로 main을 전달했다. 기존 origin 경로가 없어 fetch origin 대신 bundle을 사용했다. DB migration diff는 없다. 기존 배포04a65b2의 DB backup은 `04a65b20-20261009T171623Z.dump`다. 크기8432bytes, TOC23, relay_state1행을 확인했다. main 배포 뒤 migrate0, relay healthy, loopback health200·owner/API401·비공개PG·자원 제한·0600/0700 검사가 통과했다. 로컬 Docker는 실행하지 않았다.

원본 Tunnel config SHA256은 `40ce65ea1450329ac73c9b3188f12966b99ea562fae31c334e06dd92b06cf283`이다. 운영 상태 디렉터리에 `tunnel/config.pre-google002.yml`0600으로 보존했다. 기존 pinned cloudflared2026.9.1 바이너리를 별도 private tools 경로에 준비하고 render-public-config를 실행했다. 실제15경로 검사를 통과한 뒤 config를 적용했다. Compose 환경 변수 누락으로 첫 restart는 실행 전에 실패했다. 기존 state/UID/GID를 지정한 재시도는 성공했다. 실행 중인 컨테이너의 mounted config validate와 owner/test 경로의404 선택을 다시 확인했다. 적용 SHA256은 `3b6f569e5603489496cda27802626043ef5716dae9954c0e457c95da4c7c8cd3`이다.

두 KnowsLink Access 앱의 full app/policies는 Git 밖 비공개 백업에 보존했다. 백업 SHA256은 `0b0ee7179c586cc3432bd89e0f68dd392ab0eeddfcd3aa92363d5207c36d5019`이다. 관리자 차단 적용 뒤 member 앱 fc81b205-d4d1-445e-a2bd-384a2ed82f62와 owner 앱 bd210310-fd5e-4e6e-8cb4-d36d618cbebd만 삭제했다. 로그인된 Orca dashboard의 인증 요청은 각각202/success였고 목록 GET200의 apps는 빈 배열이다. 첫 DELETE는 dashboard 보안 헤더 누락으로 HTML 응답이어서 목록으로 미변경을 확인한 뒤 정상 요청 헤더로 실행했다. 정책·IdP·DNS·다른 Tunnel·구독은 변경하지 않았다. Access 활성화·결제·초과 자동과금 동의는 하지 않았다.

실제 Orca 브라우저에서 익명 root가 KnowsLink Google 로그인 화면으로 열린다. 외부 HTTPS 검사에서 root200, 세션 없는 home303→자체 `/`, 잘못된 poll401 JSON이다. owner/owners/정확한 keys/authorize/test/healthz/unknown은 모두404다. 서비스 CSP 때문에 브라우저 eval fetch는 실패하여 음성 HTTP 검사는 서버의 외부 HTTPS 요청으로 확인했다. DOM 화면과 실제 로그인 링크는 Orca 브라우저로 확인했다.

공유 서비스 응답200/200/401, 별도 host connector PID, myportfolio 컨테이너 상태와 다른 Compose project가 같다. 기존 verify.py regression은 baseline의 knowslink만 제외하지 않는 비교 결함으로 실패했다. 양쪽 projects에서 이 과제의 knowslink를 제외한 직접 비교는 통과했다. 제품 검사기 수정은 이 운영 적용에서 하지 않았다.

## 남은 수락과 복구

현재 Codex 세션의 실제 stdio MCP로 connect를 실행했다. 별도 로컬 private 폴더에 요청이 생성됐고 Orca 브라우저 Google 로그인 화면에서 사용자 로그인을 기다린다. Grok Bot에는 수락된 main 갱신·실제 connect를 요청했다. 실제 Google 등록·두 agent 관계 수락·양방향 대화는 아직 수락하지 않았다. 키·credential·로그인 비밀은 이 문서에 저장하지 않는다.

medium 익명 신규연결 cap2000/만료 뒤24h 포화 한계는 남는다. 기존 회원·키의 메시지는 이 cap에 묶이지 않는다. 장애 시 D12대로 먼저 KnowsLink cloudflared를 중지한다. 원본 config와 Access 앱을 복구할 때 새 aud를 대조한다. 과거 DB를 덮어쓰지 않는다.

reviewer와 DEV 신규 터미널은 정상 release했다. snapshot cleanup은 관리 경로 정체성 및 Orca release 확인 실패로 보류했다. 강제 삭제하지 않았다. 담당 coor, 재개 조건은 정본 기록 checkout과 release 조회의 정체성 확인이다. 기존 user_owned 세션은 유지한다. 역할 동기화는 쉬는 상태·clean 확인 뒤 수행하며 상태 불명 역할은 최신main ba3754d 기준으로 다음 dispatch 전에 동기화한다.
