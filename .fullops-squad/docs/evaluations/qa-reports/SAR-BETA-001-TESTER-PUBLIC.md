---
title: SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 베타 연결의 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-001-TESTER-PUBLIC]
summary: 28bd1bb gate와 보호된 공개 연결의 독립 QA 결과와 인간 로그인 held를 기록한다
---

# SAR-BETA-001-TESTER-PUBLIC — 최신 gate와 보호된 공개 베타 연결의 독립 QA 보고서

## 판정

판정 대상 배포 SHA는 `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`다. 읽기 전용 snapshot과 배포 체크아웃의 HEAD가 같다. 둘 다 detached이고 porcelain은 0줄이다.
`gates.py` 종료코드는 0이다. aud 불일치는 일반 실행, `-O`, `-OO`, `PYTHONOPTIMIZE=1`에서 종료코드 1이다. proof ID 불일치, missing proof, stale proof, future proof도 종료코드 1이다.
이 검사는 임시 상태 디렉터리만 썼다. snapshot, 배포 체크아웃, `/home/shin/deploy/knowslink-state`의 메타데이터 다이제스트는 검사 전과 같다.
coor가 `apply_complete`를 회신한 뒤에 외부 HTTP와 공유 서비스 회귀와 로컬 증명을 실행했다. tester는 Access, DNS, connector, 제품 소스, 배포 소스를 변경하지 않았다.
공개 해석기 경로의 미인증 HTTP는 종료코드 0이다. `/`, `/owner`, `/v1/registry`, `/healthz`, `/_probe`와 가짜 JWT, service token, Bearer는 상태 302다. Location 호스트는 `scshin88.cloudflareaccess.com`이다. 응답 본문에 `relay`가 없다.
HTTP 80의 `/`는 상태 301이다. Location scheme은 `https`이고 호스트는 `link.knowslog.com`이다.
`verify.py public` 종료코드는 1이다. 예외는 `socket.gaierror`다. 이 서버의 기본 해석기는 `link.knowslog.com`에 대해 NXDOMAIN을 반환한다. `dig @1.1.1.1`과 `dig @8.8.8.8`은 같은 Cloudflare A 두 개를 반환한다. 공개 판정은 `public_http.py`다.
`verify.py regression` 종료코드는 0이다. `orca.knowslog.com`은 200, `s8.knowslog.com`은 200, `mcp.knowslog.com`은 401이다. 호스트 cloudflared PID는 `506937`이다.
`knowslink-cloudflared-1`은 별도 컨테이너로 실행 중이다. 호스트 포트 바인딩은 없다. 루트 파일 시스템은 읽기 전용이다. 명령은 상태 디렉터리의 tunnel 설정을 사용한다.
보존된 `access.live.json`은 0600이고 검사 시점 나이는 332초다. 앱 ID, 정책 ID, aud, teamName, `required: true`가 실제 `tunnel/config.yml` 및 `access.json`, `access.aud`와 일치한다. 정책은 reusable allow 1개이고 include 이메일은 소유자 파일과 같다. 이메일 값은 기록하지 않았다.
`allowed_idps`는 ID 1개다. 보존된 GET에 IdP 타입 필드는 없다. 나이가 10분 이하이고 앱·정책·aud·team이 일치하므로 추가 관리 GET은 요청하지 않았다.
인간 이메일 OTP 로그인은 실행하지 않았다. 담당은 사용자이고 상태는 held다.
제품 경로와 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 diff는 비어 있다. `verify.py`, `beta.sh`, `compose.ops.yaml`, tunnel 템플릿 blob은 `f824015314c66bcab42940cfe3db2edabb22e1dd`와 같다. 배포 트리에서 바뀐 파일은 `access_apply.py` 하나다.

시나리오: [SAR-BETA-001-TESTER-PUBLIC.md](../scenarios/SAR-BETA-001-TESTER-PUBLIC.md). 로그: [SAR-BETA-001-TESTER-PUBLIC-test/](SAR-BETA-001-TESTER-PUBLIC-test/).

## 기준

날짜는 2026-10-03이다. 기록 브랜치는 `fullops/tester`다. 공통 기준은 `fullops-common-0.3.2`다.
정본은 `.fullops-squad/project.md`다. 문서 작성 규칙과 `fullops-test`를 적용했다. 예외는 없다.
snapshot 경로는 `/tmp/knowslink-beta-review-28bd1bb`다. 실행 파일은 그 트리의 `deploy/knowslink/access_apply.py`다. `PYTHONDONTWRITEBYTECODE=1`과 `python3 -B`를 사용했다. 게이트 실행의 `PYTHONOPTIMIZE`는 사례별로만 지정했다.
`verify.py` 실행 시점의 `sys.flags.optimize`는 0이다. 환경 변수 `PYTHONOPTIMIZE`는 없었다.
Jev keep 문서를 읽었다. 충돌 후보와 주의 후보는 없다. D11, D12, D13은 읽기만 했다.

## 게이트

| 검사 | 모드 | 종료코드 | 결과 |
|---|---|---|---|
| selftest | 일반, `-O` | 0 | 통과 |
| fresh proof | 일반, `-O` | 0 | 통과 |
| aud 불일치 | 일반, `-O`, `-OO`, `PYTHONOPTIMIZE=1` | 1 | 통과 |
| proof app ID 불일치 | 일반, `PYTHONOPTIMIZE=1` | 1 | 통과 |
| proof policy ID 불일치 | 일반 | 1 | 통과 |
| `access.live.json` 없음 | 일반 | 1 | 통과 |
| `access.json` 없음 | 일반 | 1 | 통과 |
| stale 660초 | 일반, `PYTHONOPTIMIZE=1` | 1 | 통과 |
| future 300초 | 일반, `PYTHONOPTIMIZE=1` | 1 | 통과 |
| snapshot·배포·상태 메타데이터 | 검사 후 | 불변 | 통과 |

더미 이메일과 더미 aud만 임시 디렉터리에 썼다. 로그에 이메일 값은 없다. 전체 종료코드는 [gates.log](SAR-BETA-001-TESTER-PUBLIC-test/gates.log)의 `EXIT:0`이다.

## 공개 HTTP와 증명

| 검사 | 종료코드 | 관찰 |
|---|---|---|
| `public_http.py` | 0 | 443 경로 8개 상태 302. 본문 `relay` 없음. HTTP 80 상태 301 |
| `verify.py public` | 1 | `socket.gaierror`. 상태코드 없음 |
| `verify.py regression` | 0 | 200, 200, 401. PID `506937` |
| `proof_check.py` | 0 | 나이 332초. 불일치 항목 없음 |
| connector inspect | 0 | `knowslink-cloudflared-1` running. `ports={}` |

probe 실행에서 1.1.1.1과 8.8.8.8의 A 목록은 같았고 개수는 2다. curl은 1.1.1.1의 첫 A에 `--resolve`로 연결했다. 이후 [dns.log](SAR-BETA-001-TESTER-PUBLIC-test/dns.log)의 A는 `104.21.54.184`와 `172.67.141.9`다. 시스템 `getaddrinfo`는 실패했다.
증명 파일과 tunnel 설정, `.env`, 자격 json 1개는 0600이다. 상태 디렉터리와 tunnel 디렉터리는 0700이다. 소유자 이메일 파일은 0600이다.
정책 이름은 `knowslink-beta-owner-only`다. decision은 allow다. exclude와 require는 비어 있다. `options_preflight_bypass`는 false다. custom page는 없다. launcher는 숨김이다. `originRequest.access.required`는 true다. audTag는 64-hex 1개이고 기록된 aud와 같다. 값은 보고서에 적지 않았다.

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `gates.py` | 0 | [gates.log](SAR-BETA-001-TESTER-PUBLIC-test/gates.log) |
| `proof_check.py` | 0 | [proof.json](SAR-BETA-001-TESTER-PUBLIC-test/proof.json) |
| `public_http.py` | 0 | [public.log](SAR-BETA-001-TESTER-PUBLIC-test/public.log) |
| `verify.py public` | 1 | [verify-public.log](SAR-BETA-001-TESTER-PUBLIC-test/verify-public.log) |
| `verify.py regression` | 0 | [verify-regression.log](SAR-BETA-001-TESTER-PUBLIC-test/verify-regression.log) |
| `docker inspect knowslink-cloudflared-1` | 0 | [connector.log](SAR-BETA-001-TESTER-PUBLIC-test/connector.log) |

`make test`, `make verify-mvp`, 백업, 복원, `beta.sh expose`는 실행하지 않았다. 제품 트리와 `verify.py` blob이 기존 QA 조건과 같기 때문이다. 캡처와 영상은 만들지 않았다.

## 재사용

로컬 런타임 QA의 보고 커밋은 `1762b430bed1c0584fecd163ae81567a4a5d04a9`다. 실행 배포는 `f824015314c66bcab42940cfe3db2edabb22e1dd`다. loopback, 백업, 격리 복원, expose 차단, 합성 인증 음성은 그 조건으로 재사용한다.
제품 QA 보고 커밋은 `659f4b06dfda1bd59997f8f2b06cba23b48026a6`다. 제품 SHA는 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`다. 리뷰 `311381feb0f18203f60329605d762257bfff8421`도 원래 SHA와 조건으로 재사용한다.
이번 새 실행은 `28bd1bb` gate와 공개 해석기 HTTP와 regression과 로컬 증명 대조다. 이전 held를 이번 통과로 바꾸지 않는다.

## held

인간 이메일 OTP 로그인은 held다. tester가 대신 실행하지 않았다.
DEC-02 실데이터, DEC-03 공개 한도와 신원, 실adapter, 실제 벤더, WAL 또는 backup 삭제, 장기 운영은 held다.
보존된 GET은 IdP 타입을 포함하지 않는다. OTP 타입 문자열은 이 QA가 관찰한 값이 아니다.
`verify.py`의 검사는 `assert`다. 이번 실행은 최적화 플래그 0으로만 했다. 리뷰 R1은 그대로다.
기본 해석기의 NXDOMAIN 캐시는 공개 DNS 실패와 구분한다. 공개 해석기 조회는 성공했다.

## 정적 검사

lint, strict, 공백 검사의 종료코드는 증거 커밋 뒤 깨끗한 트리에서 확인하고 이 절과 완료 보고에 갱신한다.
