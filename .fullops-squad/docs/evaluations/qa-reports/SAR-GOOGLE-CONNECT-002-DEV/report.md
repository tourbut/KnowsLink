---
title: 비용 없는 Google 연결 DEV 검증
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-002-DEV]
summary: 서버 격리 검사·실제 ingress·임시 Bot 설치와 Windows 로그인 회귀 증거 및 수락 한계
---

# 비용 없는 Google 연결 DEV 검증

- 담당: DEV / Codex gpt-6.1-sol medium / 세션 `01a1218c-e248-7962-a471-d3a01468503b`.
- Task `task_73d249b1bcff`, Dispatch `ctx_9dc5e970a7f8`, base `5af28c9fcfc289a845f9730aed34d32308e8b5a6`.
- 제품 코드 검증 후보: `b9e7717`. 최종 문서·archive 커밋은 worker_done의 결과 SHA를 따른다. 최종 gate는 같은 최종 HEAD의 깨끗한 서버 clone에서 실행하며 JSON은 해당 로컬 Git 디렉터리 `fullops-gate/google002-final-lint.json`에 보존한다.

| 위치 / 명령 | 실제 종료코드 | 증거 |
|---|---|---|
| 서버 `make install` | 0 | [server-install.log](server-install.log) |
| 서버 `make lint` | 0 | [server-product-lint.log](server-product-lint.log) |
| 서버 `make test` | 0 | [server-product-test.log](server-product-test.log) |
| 격리 DB `go run ./cmd/migrate` | 0 | [server-migration.log](server-migration.log) |
| `go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP\|TestGoogleHTTP\|TestEmailIdentity' -count=1` | 0 | [server-integration.log](server-integration.log) |
| 실제 cloudflared2026.9.1 validate·15개 ingress rule | 0 | [server-cloudflared.log](server-cloudflared.log) |
| 새 임시 prefix `sh scripts/install_bot_mcp.sh` | 0 | [server-bot-install.log](server-bot-install.log) |
| Windows Node22.22.2 `node adapters/dist/login.test.js` | 0 | [windows-login.log](windows-login.log) |
| Windows Node22.22.2 `node adapters/dist/mcp.test.js` | 0 | [windows-mcp.log](windows-mcp.log) |

서버는 `.env.server`를 비공개 SSH 메모리에서만 읽었다. IP·USER·PW를 argv·출력·Git에 넣지 않았다. Paramiko5.0.0은 로컬 접속 도구이며 제품 의존성이 아니다. 새 서버 호스트 키는 과제 전용 비공개 known_hosts에 고정했다. Node22.22.2/npm10.9.7, module 지정 Go1.27.1, GNU Make 및 기존 pinned Postgres17/cloudflared2026.9.1을 사용했다. 로컬 Docker를 실행하지 않았다.

`/tmp/knowslink-google002-ctx9dc5e970/work`는 Git bundle의 깨끗한 별도 clone이다. 실제 명령 순서는 [server-checks.sh](server-checks.sh)에 보존했다. Postgres는 고유 `knowslink-dev-google002-ctx9dc5e970` 컨테이너·임의 loopback 포트·합성 전용 DB·256MiB/1CPU다. migration·통합검사 뒤 자기 컨테이너만 제거했다. cloudflared binary 복사용 자기 stopped 컨테이너도 제거했다. 전후 docker ps 이름/ID 목록이 cmp exit0으로 같았다. 운영 DB·회원·키·서비스·Tunnel·Access·다른 호스트는 변경하지 않았다. shell set -e와 실제 프로세스 반환값으로 성공을 확인했고 tail 파이프로 명령 종료코드를 가리지 않았다.

make lint에는 public renderer의0600·별도파일·404fallback·기존 후보 덮어쓰기 거부·잘못된 UUID·validator 실패 전파가 포함된다. 로그의 예상 ERROR 두 줄은 음성 검사이며 make lint 자체는0이다. 실제 cloudflared는 member8경로 rule0과 owner/admin/test/미등록7경로 rule1(404)을 확인했다. 실제 edge Access 제거는 수행하지 않았다.

make test는 cap2000번째 허용/다음 거부·member 요청 제외·만료 뒤24h 직전 보존/경계 회복과 기존 다른 회원의 Google callback 바인딩 거부를 포함한다. 후자는 회원·세션 개수와 기존 owner/state가 바뀌지 않는지 확인한다. Windows/Linux login은 상위 폴더가 없는 새 경로를 생성하고 최종 폴더 exclusive 생성·ACL/모드·서명·자동 저장·별도 키·만료를 확인했다. MCP initialize/listTools는 도구9개와 held에서 network/key read 없음·실패 경계를 확인했다. 임시 Bot installer는 pinned Node checksum·압축 해제 bundle의 실제 stdio 검사를 통과했고 사용자 앱 등록을 바꾸지 않았다.

제품/문서 변경은 작은 renderer와 기존 테스트·설치 안내 중심이다. 새 제품 의존성·SQL·auth framework·UI는 없다. 기존 medium 신규 Device 포화 한계는 cap2000/만료 뒤24h 보존 정책과 함께 남는다. 공개 전 coor가 이 한계를 인지해야 한다. 운영 적용·실제 Google 계정·Bot UI 설치·관계 수락·외부 양방향 text·직접 시각 검수·고정 SHA 독립 리뷰는 실행하지 않았고 coor/사용자 후속이다. DEV 합성 통과를 실제 로그인 성공으로 표시하지 않는다.
