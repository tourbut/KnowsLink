---
title: Google 클라이언트 연결 DEV 검사
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-001-DEV]
summary: Windows Google 연결 핵심 검사와 운영 수락 한계를 기록한다
---

# Google 클라이언트 연결 DEV 검사

검사 담당은 DEV이며 테스트 레벨은 lite다. 기준 ref는 `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`이고 착수 HEAD는 `f6daf938e80ed28c9e9500baad6fa3c61ecd502c`다. commit 전 변경 후보의 검사이며 최종 고정 SHA와 commit 후 등록 게이트 결과는 worker_done 및 Git 공용 디렉터리의 `fullops-gate/google-connect-lint.json`을 따른다. 독립 리뷰·QA와 제품 운영 수락은 부모 담당이다.

| 명령 | 실제 종료코드 | 증거 |
|---|---|---|
| `make lint` | 0 | [product-lint.log](product-lint.log) |
| `make test` | 0 | [product-test.log](product-test.log) |
| `go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP\|TestGoogleHTTP\|TestEmailIdentity' -count=1` | 0 | [integration.log](integration.log) |

product-lint는 gofmt/vet/module 검증, Prettier/ESLint/TypeScript, Compose, shell/Python 및 공개 ingress 경계 self-check를 수행했다. product-test는 Go race와 MCP initialize/discovery/held, 시험 경계, 수동 connect, Google login 로컬 서명·자동 저장·별도 키·권한·만료 검사를 수행했다. 통합 검사는 격리 Postgres 17과 합성 RSA OAuth/PKCE 서버로 같은 Google 회원의 두 연결·별도 agent/credential·명시적 동의·최근 인증·Origin·수동 confirm 우회 거부·replay 및 기존 Google/이메일 회귀를 검사했다.

Windows의 기존 Node24/make 부재를 로컬 Node22.22.2/npm10 및 w64devkit2.10.0으로 복구했다. 등록 명령과 실패 기준은 변경하지 않았다. 초기 registry digest/format 실패는 LF checkout으로 해결했다. 초기 bundle CLI 부작용은 entry 파일명 검사로 해결했다. WindowsPS Set-Acl module 충돌은 .NET ACL로 해결했다. 새 최근 인증 테스트의 403 기대값은 기존 422 계약에 맞게 고쳤다. 이 중간 실패들을 최종 PASS로 소급하지 않는다.

직접 UI 검수는 기존 Service handler를 쓰는 합성 localhost fixture에서 Chrome으로 수행했다. desktop과 mobile 390×844에서 키 지문과 동의 버튼을 확인했고 document scrollWidth/clientWidth는 모두 375였다. 명시적 동의 후 approved 화면과 잘못된 요청 거부 화면을 확인했다. 화면 캡처 파일은 저장하지 않았다. 임시 프로세스와 격리 DB는 검사 후 정리한다.

실제 Google 계정·Cloudflare Access/Tunnel 적용·외부 Bot 설치·실제 왕복은 실행하지 않았다. 기존 개인키·비공개 .env 링크와 전달 WIP를 보존했다. 실제 운영 검증·복구는 D12, 사용자 절차는 D11, 추적성·기술 선택은 실행 기록과 packet-outcomes를 따른다. SQL schema·dependency 변경은 없으며 기존 held/fail은 유지한다.

## 최신 지시: 운영 서버 격리 검사

coor의 `msg_641fb3898ab1`과 `msg_d9bb45f4a345`를 수신한 뒤 로컬 Docker 작업을 중단했다. 사용자 제공 `.env.server`를 비공개 SSH 클라이언트로 읽었다. 접속 IP·USER·PW는 인자·출력·Git에 넣지 않았다. 기존 known_hosts를 검사하고 처음 보는 서버 키는 이 과제의 비공개 임시 known_hosts에 고정했다. Paramiko 5.0.0은 로컬 접속 도구이며 제품 의존성이 아니다. Context7 quota 오류 후 [공식 SSHClient 문서](https://docs.paramiko.org/en/stable/api/client.html)를 확인했다.

구현 SHA `b8fce7e99fb7cc02ae7f0cc11bf9bb7407002bf6`의 Git archive를 운영 서버의 별도 `/tmp/knowslink-google-connect-b8fce7e-ctx0b917bf8`에 업로드했다. Node22.22.2/npm10.9.7/Go1.27.1/GNU Make4.3과 기존 Postgres17 이미지를 사용했다. 기존 relay/cloudflared/postgres와 다른 서비스는 정지·재시작·설정 변경하지 않았다. 운영 환경 파일을 소스 작업 공간으로 복사하지 않았다.

| 운영 서버 명령 | 종료코드 | 증거 |
|---|---|---|
| `make lint` | 0 | [server-product-lint.log](server-product-lint.log) |
| `make test` | 0 | [server-product-test.log](server-product-test.log) |
| `go run ./cmd/migrate` | 0 | [server-exits.txt](server-exits.txt) |
| 좁은 Google/이메일 integration race 검사 | 0 | [server-integration.log](server-integration.log) |

DB 검사는 별도 `knowslink-dev-google-ctx0b917bf8` 컨테이너에만 적용했다. 임의 loopback 포트·합성 전용 계정·256MiB/1CPU 제한을 사용했고 EXIT trap으로 컨테이너를 제거했다. 운영 DB·회원·키·Tunnel과 공유 서비스 데이터는 읽거나 변경하지 않았다. 운영 서버에서 실행한 합성 검사를 실제 Google/외부 Bot 운영 수락이라고 부르지 않는다.

구현 commit 후 Windows lint.py는 동일 기준 ref에서 ERROR0/WARNING8/실행 불가0, product-lint/product-test 모두 exit0였다. SEC-001은 로컬 합성 OAuth의 `ClientSecret:local-fixture`이며 실자격이 아니다. DEP-001은 package.json 테스트 명령에 login 검사를 추가한 것이며 의존성 변경은 없다. SLOP-004는 ingress self-check의 사용자 결과 출력이다. SIZE 경고는 기존 큰 문서·member 모듈과 새 MCP 도구·보안/ACL 검사·문서 추가 때문이다. 별도 기능이나 구조 개편으로 분할하지 않았으며 부모가 고정 SHA 리뷰에서 확인한다. 최종 문서 후속 SHA의 게이트는 운영 서버의 별도 clean clone에서 실행하고 JSON을 worker_done에 연결한다.
