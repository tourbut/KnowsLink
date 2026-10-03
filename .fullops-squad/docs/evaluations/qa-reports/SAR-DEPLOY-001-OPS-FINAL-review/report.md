---
title: SAR-DEPLOY-001-OPS-FINAL 독립 문서 리뷰
status: review
updated: 2026-10-03
owner: ops
tasks: [SAR-DEPLOY-001-OPS-REVIEW]
summary: 완료 OPS 문서 f5a73a3의 별도 세션 독립 검토 결과와 문서 수락 결론
---

# SAR-DEPLOY-001-OPS-FINAL 리뷰

- 검토자 / CLI / 모델: Claude Code `claude-sonnet-5-5`. Orca task `task_097cae8b7185`, dispatch `ctx_d7459ee5ea60`, 실제 session `6147402e-3bd9-4057-a87c-c98a1db6ec82`.
- base SHA / head SHA / merge-base: `ffca87c9ec501e9f313c50f60124b686d27af728` / `f5a73a3b4b4dc9fe7f3a731fe0952a798f119219` / `ffca87c9ec501e9f313c50f60124b686d27af728`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate. `.fullops-squad/review/rule.json`의 `**/*.md`·`**/*.json` 규칙(rules.json 2그룹). 공통 기준 `fullops-common-0.3.2`(snapshot의 rules/common README·coding-style·testing·security), `project.md`, `docs/agents/document-writing.md`를 snapshot에서 읽었다. 예외 없음.
- 요구사항·완료 기준 원천: snapshot 시점 OPS 지시서(완료 아카이브 `handovers/logs/2026-10-03_to_ops.md`), 현재 리뷰 지시서 `handovers/to_ops.md`. 원래 `SAR-DEPLOY-001-OPS-find/documents-find/context.json` 근거를 재사용했다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 13 / 13 / 0 / 13 / 0.
- lint(`lint.json`): ERROR 0 / WARNING 0 / 실행 불가 0. 아래 검증 절의 첫 실행 실패 이력을 참고한다.

## 독립성

- 구현자: Codex `gpt-6.1-sol`(codex-tui, cwd `fullops-ops`, dispatch `ctx_3c54fe7d2043`) 세션 `01a1010f-21e1-7863-b6f7-83064be0c312`. 세션 메타데이터(`~/.codex/sessions/2026/10/03/rollout-2026-10-03T18-18-51-…jsonl`)에서 확인했다. 이 세션이 f5a73a3(UTC 09:22:50)를 만든 시점과 일치한다. dispatch ID는 세션 ID가 아니므로 추정하지 않았다.
- 최초 D12 초안·서버 조사는 별도 Claude 세션 `f4f63623-13fc-4e08-bfb9-3cdfb3192fbc`다. 이번 검토 세션 `6147402e-…`과 다르다.
- 검토 세션은 구현·초안 세션과 모두 다르다. snapshot `/tmp/knowslink-ops-review-f5a73a3`는 detached head `f5a73a3`, 추적 파일 변경 0이다. ignored `adapters/node_modules`만 설치했다.

## 검토 범위

13개 파일을 모두 reviewed했다. 파일별 근거는 `result.json`에 있다. 제외 파일은 없다.

- **D12(`docs/operations/ops-guide.md`)**: 전문을 읽고 `a6a10c7`의 `compose.yaml`·`Dockerfile`·`.env.example`·D05 인터페이스설계서·`db/migrations/00001_relay.sql`·`cmd/migrate`·`cmd/relay/main.go`·`internal/relay/http.go`와 대조했다. 서비스·네트워크 이름과 `internal`, `127.0.0.1:8080` 게시, Postgres 미게시, migrate `restart: "no"`, restart 정책 부재, `TUNNEL_TOKEN`·`tunnel` profile, `pg_dump -U knowslink`, `relay_state`·goose version table, `/healthz` 200, 미정의 경로 404(기본 ServeMux), 보호 경로 401·CSRF 403이 문서와 일치한다.
- **QA 참조**: `c59537b`의 `SAR-MVP-001-TESTER.md`에서 QA-01–11 실행 항목 통과·held 8·미해결 제품 결함 0을 확인했다. D12가 QA-01/08 기준을 인용한 내용과 일치한다. D12는 QA 통과를 제품 수락으로 확대하지 않는다.
- **C1 high**: D12·실행 기록·contexts가 reviewer `msg_4fbcac80f76c`의 C1 high(agent credential이 `deliver:human`을 pull·persist·ACK·claim)를 근거로 `a6a10c7` 수락·배포를 금지한다. `a6a10c7`의 `pull`/`persist`/`ack`/`claim` 처리 코드에 `deliver:human` 구분이 없음을 읽어 서술과 모순이 없음을 확인했다. 이 검토는 C1을 새로 재현·수정하지 않는다.
- **원래 서버 관찰**: 문서는 원본 명령 로그가 없다고 적는다. 원본 Claude 세션 `f4f63623-…` 기록에서 Docker 29.4.3·Compose 5.1.3, 메모리 13Gi/가용 7.0Gi, cloudflared 2026.8.3·PID 20일 가동, `orca` Tunnel `3e132faa…`, `hermes-9ae8f587` Exited(255), `myportfolio-master-viewer-1` Created, 포트 8765·6768·8766이 존재함을 대조했다. 서버 명령은 이 검토에서 새로 실행하지 않았다(F1).
- **미실행·held**: D12 9장과 실행 기록 모두 DNS·Tunnel·컨테이너 변경, 수락 SHA 배포, 공개 health 검증, DB 복원을 미실행으로 적는다. D11/D13 미작성과 인덱스(D12 draft) 일치. 공개 가입 `POST /v1/owners`가 무인증 합성 가입이므로 공개 정책·신원(DEC-03)이 held인 동안 공개 연결 보류라는 서술이 타당하다.
- **외부 API 근거**: Cloudflare 공식 run-parameters 문서를 조회해 `--token`은 remotely-managed 전용, `--config`는 locally-managed 전용임을 확인했다. D12의 override(`tunnel --no-autoupdate --config … run`)와 `route dns`에 `--overwrite-dns` 미사용 계획이 이 구조와 맞다. 비공개 코드·값은 조회문에 넣지 않았다.
- **문서 링크·비밀값·JSON**: 변경 Markdown의 상대 링크 전수 확인에서 깨진 링크 0. 변경 JSON 6개 모두 유효. 비밀값·토큰·capability·credential 패턴 검색에서 해당 없음. Tunnel UUID 접두·ID는 식별자이며 자격 파일 내용이 아니다. `.env.example`은 jev context에서 refused, `compose.yaml`은 unsent로 처리돼 있다.
- **소유권**: PLANS.md(빈 줄 1개)·board.json·to_dev.md·MVP-DEV route는 coor/designer 계열 변경이며 base 이후 동기화 커밋에서 들어왔다. OPS 소유 변경은 contexts/ops.md·D12·D12 인덱스 행·실행 기록·완료 아카이브·빈 인박스다. D12 원문은 이 검토에서 수정하지 않았다.
- **빈 inbox·아카이브**: snapshot의 `handovers/to_ops.md`는 비어 있다(f5a73a3에서 62줄 삭제). 지시서·완료 보고 전문이 `logs/2026-10-03_to_ops.md`에 보존돼 있다.

## 발견 사항

미해결 critical/high는 없다. 모두 low이며 수락을 막지 않는다.

| 심각도 | 파일·줄 | 내용 | 상태 |
|---|---|---|---|
| low (F1) | `docs/operations/ops-guide.md:14` | 원본 명령 로그가 없다고 적지만 원본 Claude 세션 `f4f63623-…`에 출력이 남아 있고 이 검토가 대조했다. 재개 때 이 세션을 근거 경로로 연결하면 추적성이 좋아진다. 서술은 사실과 모순되지 않는다(신규 PASS를 주장하지 않음). | 열림, 후속 |
| low (F2) | `board/board.json:70` | 운영 배포 phase note가 C1 배포 금지·사용자 중지 지시를 담지 않는다. coor 소유다. | 열림, coor |
| low (F3) | `handovers/to_dev.md:115` | "수락 뒤 … 기존 승인 근거로 Tunnel 연결" 문구가 이후 중지 지시·D12 1장과 충돌할 수 있다. 실행 권한은 D12를 따른다. | 열림, coor |

기타 관찰: D12 4장의 자격 파일 UID(65532) 읽기 권한, 백업 보관·암호화 정책, 쓰기 권한은 문서가 이미 후속 확인·held로 표시한다. 결함이 아니다.

## 검증 및 남은 제약

- `python3 <0.9.12 deliverables.py> --repo /tmp/knowslink-ops-review-f5a73a3 --strict`: 종료코드 0, 검사 13·미작성 9·문제 0·경고 0.
- `lint.py --repo /tmp/knowslink-ops-review-f5a73a3 --from ffca87c… --out …/lint.json`: 첫 실행은 `adapters/node_modules` 부재로 `prettier: not found`, 종료코드 1이었다. 의존성 부재일 뿐 제품·문서 결함이 아니다. snapshot의 ignored 경로에 `npm ci --prefix adapters`(종료코드 0)를 설치한 뒤 재실행해 `product-lint` passed, ERROR 0·WARNING 0·실행 불가 0, 종료코드 0이다. `lint.json`은 통과 결과만 보존한다. snapshot 작업 트리는 이후에도 깨끗하다.
- `review.py check`는 아래 결론 절 직후 실행하고 종료코드를 worker_done에 기록한다.
- 미실행: 서버·Docker·DNS·Tunnel·DB 복원·공개 health/auth 재검증, C1 재현. 이 검토는 문서 리뷰이며 해당 실행은 범위 밖이다. 제품 코드는 base 대비 변경이 없어 `make test`·`make build`·runtime은 적용하지 않았다.
- 직접 시각 검수·독립 QA는 이 문서 변경의 대상이 아니다.

## 검토 결론

**문서 수락: 가능.** 완료 OPS 문서(D12 draft 계획·실행 기록·contexts·아카이브·인덱스)는 고정 SHA `f5a73a3` 기준으로 실제 제품 참조·공식 문서·원본 관찰과 일치하고, 미실행·held·C1 high 배포 금지·비밀값 보호·rollback 경계를 정확히 서술한다. 미해결 critical/high 없음, lint ERROR 0이다. 이 수락은 coor의 main 통합 대상이 되는 문서 기록에 대한 것이다.

**제품·배포 수락: 아님.** `a6a10c7`은 C1 high로 수락·배포가 금지된다. DEV 수정, 새 고정 SHA의 독립 QA·리뷰·직접 시각 검수·coor 수락, 사용자 재개 지시와 후속 Dispatch 전에는 DNS·Tunnel·컨테이너 변경을 하지 않는다. D11/D13은 미작성이며 held 항목(DEC-02/03 등)은 유지한다. check 통과는 기록 검사이며 내용·테스트 성공의 자동 보증이 아니다. 이 결과는 AI가 판단한 것이다.
