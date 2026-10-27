---
title: SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV install 마운트 repo 이전과 adapters npm audit 처리
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV]
summary: "compose.ops.yaml relay에 install 변수·RO 마운트를 넣어 서버 override 없이 같은 구성이 됨을 확인하고, adapters npm audit 4건(전부 런타임)을 semver 호환 패치로 0건으로 만든 결과"
---

# SAR-CONNECT-SIMPLE-001-FOLLOWUP-DEV 결과

기준 main `540df637`(인박스 `49dc5999`). 운영 서버에는 접속하지 않았다.

## 1. install 마운트 → `deploy/knowslink/compose.ops.yaml`

relay에 `KNOWSLINK_INSTALL_DIR: /install`과 `${KNOWSLINK_STATE_DIR:?…}/install:/install:ro`를 더했다. 다른 서비스·제한값은 그대로다.

로컬 Docker 29.9.2, `KNOWSLINK_STATE_DIR=/home/shin/deploy/knowslink-state`, 더미 `.env`(비밀 아님), 서버 override와 같은 내용의 임시 `compose.install.yaml`로 확인했다.

| 명령 | EXIT | 결과 |
|---|---|---|
| `docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file <dummy> config relay` | 0 | `KNOWSLINK_INSTALL_DIR: /install`, bind `source: /home/shin/deploy/knowslink-state/install` `target: /install` `read_only: true` |
| 위 + `-f compose.install.yaml` | 0 | 중복 마운트 오류 없음. 두 출력 `diff` 결과 동일 |
| `KNOWSLINK_STATE_DIR=` 로 위 첫 명령 | 1 | `required variable KNOWSLINK_STATE_DIR is missing a value`(relay·cloudflared 모두 거절, 의도된 동작) |

`beta.sh`의 `dc()`는 `-f compose.yaml -f deploy/knowslink/compose.ops.yaml`을 쓰므로 수정 없이 install 마운트를 포함한다. D12의 compose override·호출 규칙·override 제거 절차(전환 순서·확인 명령·되돌리는 법)를 갱신했다.

## 2. adapters npm audit 4건

`npm audit --omit=dev`와 `npm audit` 결과가 같았다(4건 모두 런타임 경로, EXIT 1). 경로는 `npm ls` 기준이다.

| 패키지 | 심각도 | 권고 | 경로 | 런타임 | 수정 |
|---|---|---|---|---|---|
| body-parser 2.3.0 | high | GHSA-27r3-jxmh-p828 | `@modelcontextprotocol/sdk` → `express@5.2.1` | 예 | 2.4.0 (lockfile, minor·범위 내; 의존 `qs` 6.16.0→6.17.0 동반) |
| fast-uri 3.1.8 | moderate | GHSA-w74x-v6f9-j3xg 외 3 | sdk → `ajv@8.20.0`(·`ajv-formats`) | 예 | 3.1.9 (lockfile) |
| ip-address 10.7.3 | moderate | GHSA-rxp7-6hxx-x8q5 | sdk → `express-rate-limit@8.7.0` | 예 | 10.7.5 (lockfile) |
| zod 4.6.5 | high | GHSA-mrhx-c7xc-m3w5 | 직접 의존 + sdk·`zod-to-json-schema` | 예 | 4.6.7 (patch). `package.json`이 `4.6.5`로 정확 고정이라 `"zod": "4.6.7"`로 바꿨다. major·새 의존성 없음 |

처리: Node 22.22.2/npm 10.9.7(engine-strict, main checkout `.devtools/node-v22.22.2-win-x64`)로 `npm install` 뒤 `npm audit fix`(비강제). 이후 `npm audit`·`npm audit --omit=dev` 모두 `found 0 vulnerabilities` EXIT 0.

| 검증 | EXIT |
|---|---|
| `npm run check` | 0 |
| `npm test`(build + 13개 테스트 파일 PASS) | 0 |

## 3. 남은 일

- ops: 이 커밋으로 배포한 뒤 서버 `knowslink-state/compose.install.yaml`을 D12 절차대로 `.retired`로 바꾼다(서버 미접속, 미실행).
- `make lint`·`make test`: Windows에 `make`가 없어 실행하지 않았다. 대신 FullOps `lint.py --from 540df637`과 adapters `check`·`test`를 돌렸다.
