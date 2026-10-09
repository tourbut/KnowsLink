---
title: SAR-GOOGLE-CONNECT-001-DEV 실행 기록
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-001-DEV]
summary: Google 클라이언트 연결 구현·검증과 독립 검토 및 운영 인계
---

# SAR-GOOGLE-CONNECT-001-DEV 실행 기록

## 기술 계획과 원천

현재 기준은 `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`다. 착수 HEAD와 packet HEAD는 `f6daf938e80ed28c9e9500baad6fa3c61ecd502c`로 같다. DEV 인박스와 route/find/context/packet은 전달된 미커밋 원천 스냅샷으로 읽고 보존한다. 공통 규칙 fullops-common-0.3.3과 project.md를 적용한다.

1. 기존 Connections JSON 상태를 재사용한다. 클라이언트가 새 키와 비공개 연결 수단을 만들고 공개키 소유 증명으로 요청을 등록한다. 브라우저에는 연결 수단의 해시로 만든 공개 URL만 전달한다.
2. Google 시도에 연결 ID를 묶는다. Google 인증 후 해당 회원에 요청을 고정한다. 최근 인증과 명시적인 지문 확인·동의 후 별도 agent를 만든다. 기존 completeConnection이 키 소유 증명과 한 번의 자격 발급을 담당한다.
3. Node CLI와 Command MCP에 시작·완료 동작을 연결한다. 개인키와 자격은 해당 클라이언트의 새 비공개 폴더에만 저장한다. 수동 연결·held·시험 모드는 유지한다.
4. 기존 UI/CSS를 재사용한다. 같은 회원의 다른 agent 선택을 관계 초대 UI에 추가하고 수락은 유지한다. 공개 서비스 ingress/Access의 좁은 변경과 복귀 절차는 OPS에 인계한다.
5. 핵심 성공·만료·replay·타 회원/키 분리 검사를 추가한다. Windows의 make·Node 22·race compiler를 로컬 도구 경로에 복구하고 등록 필수 검사를 실행한다.

기존 호출자는 connect.ts의 prepare/complete, text.ts의 memberTransport, mcp.ts의 public-node 모드다. relay는 member_agents.go·google.go·connections.go와 공통 HTTP origin 보호를 공유한다. SQL schema와 새 dependency는 필요하지 않다. 운영 자격·DB·Google/Cloudflare 계정은 변경하지 않는다.

Context7 Node.js resolve는 monthly quota 오류였다. 파일 권한의 Windows 차이는 Node 공식 fs 문서와 로컬 API 동작으로 확인한다. 기존 OAuth API는 SAR-GOOGLE-LOGIN-001-DEV의 고정 의존성 근거를 재사용하며 새 OAuth API를 추가하지 않는다.

## 구현 결과와 검증 인계

Device 요청·Google 회원 결합·최근 인증과 명시적 동의·서명 조회·기존 complete를 연결했다. CLI/MCP가 자기 키와 credential을 자동 저장한다. 같은 회원의 두 연결은 독립 agent이며 자동 pairing은 없다. 기존 수동 confirm이 Device를 승인하지 못하게 전용 동의 경계를 강제한다. D03/D05/D06/D09~D13을 갱신했고 D07/D08 SQL은 변경하지 않았다.

Windows 검사 환경은 사용자 로컬 도구 폴더에 Node 22.22.2/npm 10과 w64devkit 2.10.0(make/gcc/sh)을 설치해 복구했다. Go module 지정 toolchain 1.27.1을 사용했다. 기존 검사 명령·timeout·실패 기준은 낮추지 않았다. Node 의존성 추가는 없다. LF 정규화로 Windows gofmt와 동결 registry digest 실패를 해결했다. import한 CLI가 bundle에서 실행되던 문제는 entry 파일명 검사로 고쳤다. WindowsPS의 Set-Acl module 충돌은 .NET ACL API로 고쳤다. 테스트의 실제 최근 인증 거부 계약은 422여서 새 기대값의 403 오기를 수정했다.

최종 측정은 같은 변경 후보에서 `make test`, `make lint`, `go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP|TestGoogleHTTP|TestEmailIdentity' -count=1`의 실제 종료코드를 사용한다. 결과와 로그 위치는 QA 실행 보고를 따른다. 마무리 commit 뒤 기준 ref f9f7675로 등록 lint.py를 실행하며 최종 SHA와 결과는 authentic worker_done에서 전달한다.

로컬 Node HTTP 검사와 격리 Postgres 17/합성 RSA OAuth 검사는 실제 Google 계정 성공이 아니다. 직접 Chrome 화면 검수는 실제 Service handler에 합성 회원·세션·연결을 넣은 일회성 localhost fixture였다. desktop과 mobile 390×844에서 지문/동의 버튼을 읽었고 실제 document scrollWidth=clientWidth=375였다. 동의 후 approved와 잘못된 요청 거부 화면을 확인했다. 캡처 파일은 저장하지 않았으며 부모의 독립 시각 QA는 남는다.

새 런타임 의존성은 없다. 변경 규모는 trust-boundary 검증과 Windows ACL 및 회귀 검사·운영 문서 때문이다. 기존 큰 module을 별도 재구성하지 않았다. Context7 Node/Cloudflare 조회는 monthly quota 초과였고 공식 Node fs 및 Cloudflare 경로 우선순위/Tunnel 문서를 확인했다. 운영 계정·DB·실사용자 환경과 비공개 .env 링크를 보존했다. 부모가 고정 SHA 리뷰·좁은 QA·main 통합·OPS 적용을 맡는다.

## 운영 서버 검증으로 전환

최신 coor 메시지 msg_641fb3898ab1/msg_d9bb45f4a345에 따라 로컬 Docker 사용을 중단하고 .env.server로 비공개 SSH 접속했다. 구현 b8fce7e의 별도 임시 디렉터리에서 등록 make lint=0, make test=0을 확보했다. 별도 합성 Postgres17/임의 loopback 포트의 migration=0, GoogleDeviceHTTP/GoogleHTTP/EmailIdentity race integration=0이다. 기존 운영 DB·회원·키·Tunnel·다른 서비스를 보존하고 임시 컨테이너만 제거했다. 서버 로그·도구 버전·제한은 QA 보고의 최신 절을 따른다. 최종 문서 commit의 clean clone에서 등록 게이트를 실행한다.
