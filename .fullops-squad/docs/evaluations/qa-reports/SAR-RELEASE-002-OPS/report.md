---
title: main 통합 제품 7b4f9bc 운영 배포와 /lang 공개 ingress 추가 결과
status: draft
updated: 2026-10-20
owner: ops
tasks: [SAR-RELEASE-002-OPS]
summary: 제품 7b4f9bc를 기존 운영 relay에 배포한 백업·격리복원·보존·공개 검사 결과와 /lang/(ko|en) 공개 ingress 추가를 기록한다.
---

# SAR-RELEASE-002-OPS 결과

제품 `7b4f9bcf11bb7ab13207d8e635060a6e90e7f8c2`를 기존 link.knowslog.com 운영 relay에 배포했다. 새 private 백업과 격리 복원, 보존 assertion 38개, 공개 smoke를 통과했다. 사용자 승인으로 공개 Tunnel ingress에 `/lang/(ko|en)`만 추가하고 cloudflared를 한 번 재시작했다. 두 경로는 303과 언어 쿠키로 응답하고 이웃 경로는 404다.

## 범위·기준·시간

- 사용자 결정 D006 '배포 진행'(2026-10-20)과 ingress 질문의 '추가·재시작 승인'이다. 준비 HEAD는 `d0fe090`이고 브랜치는 `fullops/ops`다. 준비 HEAD의 제품 파일은 `7b4f9bc`와 같다(`git diff 7b4f9bc HEAD -- . ':!.fullops-squad'` 비어 있음).
- Task `task_ed8188b36480`, Dispatch `ctx_549ff7d89f7d`, worker `term_50c4ac91-32ed-4a58-8ea1-fe42102794fd`, 복귀 coor `term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c`다.
- 하위 에이전트는 쓰지 않았다. 직렬 운영 전환은 직접 수행했다.
- `beta.sh deploy/backup/restore-verify`와 `verify.py regression/public`은 실행하지 않았다. 삭제 포함 절차와 과거 baseline 전제 때문이다. `verify.py local`은 baseline과 전환 직후 두 번 통과했다.

| 항목 | 값 |
|---|---|
| fresh baseline | 2026-10-20T13:49:00Z. 서버 HEAD `1dd6800`·clean. relay 이미지 `sha256:4893d0f5…0098999`·healthy·restart0·시작 2026-10-17T09:38:03Z |
| 전환 구간 | STOP_BEGIN `2026-10-20T13:50:44Z`(22:50:44 KST) → UP_DONE `2026-10-20T13:51:02Z`(22:51:02 KST). 약18초 |
| 최종 relay | 이미지 `sha256:9646febaf4…c85443f66`(manifest list ID)·시작 13:50:56Z·healthy·restart0 |
| Tunnel 재시작 | 2026-10-20T13:52:03Z(22:52:03 KST). 약1.4초 끊김(아래 절) |

## 배포 입력 확인

- bundle SHA256 `43f420b04db965cb7e3c84caba1eb937cd3f63a0b64c816bd59e66149c9cf51f`(81093341bytes)다. 서버에서 hash 대조·`bundle verify`·`FETCH_HEAD`=`7b4f9bc`·`1dd6800`→`7b4f9bc` 조상 관계를 확인했다.
- `db compose.yaml Dockerfile go.mod go.sum Makefile sqlc.yaml`의 diff는 0이다. 새 migration은 없다. `deploy scripts` diff는 `public-ingress.yml`과 `check_public_ingress.py`(직전 배포의 ingress 반영)뿐이다. 제품 diff는 72파일 +5619/−276이고, `internal adapters cmd`는 지시서와 같은 46파일 +4285/−263이다.
- 고정 clean checkout에서 relay/migrate build가 EXIT0이었다. 빌드 중 기존 relay는 healthy였다.
- 새 build 전에 현재 이미지에 `knowslink-rollback:1dd6800-release002-20261020`을 붙였다. inspect 결과 ID가 `sha256:4893d0f5…`로 같다. 기존 tag(`8a93e38-release001-20261017`·`2747fe9-refactor002-20261014`)는 덮지 않았다.
- rollback 이미지는 network none·read-only 컨테이너 `knowslink-release002-rollback-inspect-20261020`(exited0 보존)에서 `/app/relay`·`/app/migrate` 실행 파일과 SHA256을 확인했다. DB 연결이나 앱 기동 검사는 아니다.

## 백업·격리 복원·재개

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `1573187`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/1dd6800-pre-release002-20261020T135049Z.dump`. 상위 0700·파일 0600·11271bytes |
| 백업 SHA256 | `32b592886b1d753c7874f909c3272c944c5585508403bce297ebe59421bcc8e5` |
| list·복원 | `pg_restore --list` TOC 25줄. network none·host port 없음·postgres 사용자로 복원 |
| 복원 검증 | 복원 `data::text` md5가 frozen 운영 값과 같다. goose2·guard1·rows1·epoch1573187·HumanModel1 |
| 거부 경계 | 격리 DB에서 HumanModel 제거 UPDATE가 거부됐다. live DB에는 거부용 쓰기를 하지 않았다 |
| 재개 | 미리 만든 이미지로 `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |
| 로그 | relay error/panic 0·Google configured=true·전환 뒤 Tunnel 프로세스 error 0 |

복원 fixture `knowslink-release002-rehearse-20261020`은 exited0·network none·port 없음으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. Down·volume 삭제·prune·이미지 삭제는 없다. 서버 `deploy-history.log`에 한 줄을 추가했다.

## 보존

[preservation.json](preservation.json)의 assertion 38개가 통과했다. 비교 단계는 baseline·before·frozen·after·final이다. 값은 read-only 조회로 메모리에서 count와 digest만 만들었다.

- 회원·agent·owner·pair·신원 등 상태 map digest, Remote 설치·구독·Events 메타·Inbox 0이 5단계에서 같다. schema SHA256·`.env`·`tunnel.uuid`·`access.aud`·Tunnel credential·파일 mode가 같다. 기존 백업17개의 hash/크기/mode가 같고 새 백업 1개만 추가됐다.
- 전환 구간에서 바뀐 컨테이너는 relay와 migrate뿐이다. Postgres와 나머지 컨테이너(공유 서비스 포함)는 ID·이미지·시작 시각이 같다. 이후 Tunnel 변경으로 `knowslink-cloudflared-1`이 재시작됐다. 새 컨테이너는 보존한 복원 fixture와 rollback 검사 컨테이너 둘이다.
- 호스트 cloudflared PID, loopback 8080/5432 listener, orca/s8/mcp 응답이 같다.
- **Rates**는 요청 시각 기록과 window 정리로 변하는 volatile map이다. 권한 변경으로 해석하지 않는다. 정확한 원인 요청은 집계로 확정할 수 없다.
- **Friends**는 전환 전(baseline/before/frozen)과 전환 후(after/final)에 각각 digest가 같고, 둘 사이에서 한 번 달라졌다. 원인을 분리했다([friends-diff.log](friends-diff.log), [ct.log](ct.log)). 새 relay가 상태를 저장하면서 기존 친구 1건의 JSON에 `"CodeTail":""`을 더 쓴다. 그 필드는 `omitempty`가 아니다(`internal/relay/human.go`). 의미는 "필드 도입 전 요청"이고 값은 비어 있다. 키 집합은 같고 그 하나 외에 다른 필드가 바뀌지 않았다. `State`·`Exp`가 같다. Friends 밖의 상태 map 변화는 없다. 이 건은 새 omitempty 필드(`Held`·`HeldFrom`·`HeldTail`)가 기존 JSON 없이 읽히는 것과 함께, 기존 상태 JSON이 새 필드 없이 읽힌다는 증거다. relay는 healthy이고 기존 회원·친구·agent 수가 같다.

보존 assertion은 짧은 전환 중의 일치다. 정상 만료·cleanup·사용자 요청에 따른 이후 변경을 금지하는 보장은 아니다.

## 공개 ingress `/lang/(ko|en)`

repo `deploy/knowslink/tunnel/public-ingress.yml`의 `path` 정규식에서 `/home(/.*)?` 뒤에 `/lang/(ko|en)` 대안 하나만 추가했다. `scripts/check_public_ingress.py`에는 허용 `/lang/ko`·`/lang/en`과 차단 `/lang`·`/lang/`·`/lang/fr`·`/lang/ko/x`·`/lang/KO` 기대값을 더했다. 지시서의 소유 목록에는 이 검사 파일이 적혀 있지 않으나 직전 배포와 같은 방식으로 같은 한 줄 변경의 검사로 더했다. 다른 경로는 추가하지 않았다.

- repo 검사: 로컬 `check_public_ingress.py` 통과. 서버 Linux에서 `git archive 7b4f9bc`의 deploy·scripts 위에 두 파일을 덮고 `check_public_ingress.py`·`test_public_config.sh`를 실행해 EXIT0이다([ingress-test.log](ingress-test.log)). 로그의 ERROR 두 줄은 `test_public_config.sh`가 일부러 만든 거부 사례의 출력이다. Windows에서는 NTFS mode 한계로 `test_public_config.sh`를 쓰지 않았다.
- live 원본 백업: `release002-ops-20261020/config.yml.pre-release002`, 0600, SHA256 `f2d8e6eed9dade1fe30f3aa1cf25c8313aa0cb480ac1d9385d570ed224ca2eaa`(적용 직전 live와 같음).
- 후보: live 규칙에서 한 대안만 더한 파일이며 repo 규칙과 문자열이 같다(`config.candidate.yml`). 변경 줄 2개. `cloudflared tunnel ingress validate` OK. 경로별 `ingress rule`에서 `/lang/ko`·`/lang/en`·`/`는 규칙 #0, `/lang`·`/lang/`·`/lang/fr`·`/lang/ko/x`와 이웃·보호 경로는 404 fallback #1이다([ingress-apply.log](ingress-apply.log)).
- 적용: 2026-10-20T13:52:03Z에 live config를 제자리 갱신하고(mode 0600 유지) `docker restart knowslink-cloudflared-1`을 한 번 실행했다. 서버에서 0.5초 간격으로 공개 `/`를 관찰했다([restart-poll.txt](restart-poll.txt)). 502(13:52:03.655) → 530(13:52:04.139) → 200(13:52:05.07)이다. 끊김은 약1.4초다. 새 프로세스는 연결 4개를 등록했다. 로그의 ERR 5줄은 종료 기록이다.
- 동일성: 적용 뒤 live의 `path` 규칙이 repo 규칙(SHA256 `99789dba9991f021f60260d23a1c38282930b4caf3f06919d2cb56b9c58ced10`인 `new-public-ingress.yml`)과 같다. 추가된 대안은 `/lang/(ko|en)` 하나다([final.log](final.log)).
- 되돌림은 필요하지 않았다. 필요하면 백업 config를 복사하고 cloudflared를 재시작한다.
- 서버 checkout은 배포 SHA `7b4f9bc`다. 이 ingress 수정 커밋은 다음 main 반영 때 서버 checkout에 도착한다.

## 공개 smoke

모두 통과했다. 실제 로그인·송신·pull·ACK·token 소비는 없다. 응답 본문과 cookie 값은 보존하지 않았다.

[public-smoke.json](public-smoke.json)(직전 smoke 재실행, `all_pass`)
- 기존 안전 검사 21/21: `/` 200과 no-store/CSP/no-referrer/nosniff, 무세션 `/home*` 303과 expired redirect, 보호 경로 404, SSE·MCP 무인증 401, OAuth well-known 200, consent 401, connect 405, 공유 orca/s8/mcp 200/200/401.
- 한국어 첫페이지 필수 문구·금지 문자열 검사 통과, Google 폼 1개·이메일 폼 0개·`<details>` 2개·`Set-Cookie` 없음.
- 새 API 무인증 경계(`/v1/contacts/{agent}`·`/v1/text/seen` 401)와 이웃 404 6개.

[public-smoke-lang.json](public-smoke-lang.json)(이번 변경 검사, `all_pass`, 22/22)
- `GET /lang/ko`·`/lang/en`·`?back=%2F`·`?back=/auth/verify`: 303, `Location`이 `?lang=`을 붙인 같은 사이트 경로, `no-store`, 쿠키 `__Host-kl_lang`(Path=/·HttpOnly·Secure·SameSite=Strict·Max-Age=31536000) 하나. `back=https://evil.invalid/`·`//evil.invalid`는 `/?lang=en`으로 정리된다.
- 이웃 경로 `/lang/fr`·`/lang/KO`·`/lang/ko/x`·`/lang/`·`/lang`·`/lang/ko%2Fx`·`/language`는 404이고 쿠키가 없다. `GET /home/language`는 404/405다.
- 첫페이지 다섯 경우(기본·`__Host-kl_lang=en`·`?lang=en`·`=ko`·잘못된 값 `fr`): `Content-Language`와 `<html lang>`이 선택 언어와 같고 잘못된 값은 한국어다. 영어는 'Connect with friends, converse through agents'다. 한국어/영어 링크와 `aria-current`가 맞고 Google 폼 1개·이메일 폼 0개, 보안 헤더(no-store·CSP·no-referrer·nosniff)가 있고 `Set-Cookie`가 없다.
- 무인증 회원 경계: `GET /home`·`/home/account`·`POST /home/language`가 303 `/?n=expired`이며 세션 쿠키 삭제(`Max-Age=0`) 외의 쿠키는 없다.
- 관찰: `/lang/*` 303 응답의 `Content-Language`는 `ko`다. 리다이렉트는 렌더하지 않는 응답이므로 이 값은 기본값이고, 이동한 페이지가 선택 언어를 따른다. 제품 결함으로 판정하지 않는다.
- 관찰: 공개 응답에는 직전과 같이 Cloudflare edge가 붙인 Web Analytics `<script>` 하나가 있다. origin 응답에는 없다. 브라우저에서 CSP가 막는지는 확인하지 않았다.

## 429 관찰(반복)

직전 배포의 교훈(재시작 직후 반복 호출 금지)을 지키지 못했다. 재시작 동안 서버 측 관찰 루프가 `/`를 0.5초 간격으로 90번 호출했다. 13:52:15Z경부터 relay를 거치는 공개 경로가 429가 됐다([restart-poll.txt](restart-poll.txt)에서 마지막 429는 13:52:54Z). 같은 구간의 `/lang/ko`·`/lang/en`·`/v1/contacts/…` 첫 확인도 429였고 `/lang/fr`·`/lang/ko/x`는 404였다. 이후 간격을 두고 다시 실행한 smoke가 모두 통과했다. 429의 정확한 종료 시각은 재지 않았다. 요청 제한의 단위(IP별인지 공유인지)는 확인하지 않았다. 다른 사용자가 영향을 받았는지는 알 수 없다. 다음 Tunnel 재시작 검사는 호출 간격을 1초 이상으로 한다.

## gate 재사용·미실행·한계

- 제품 코드 gate는 각 단위의 수락 증거(고정 SHA 독립 리뷰·tester QA·designer 시각 검수, main 통합 기록)를 재사용한다. 이번 운영 검사와 구분한다. 이 작업에서 제품 코드·테스트·의존성·migration은 바꾸지 않았다.
- 미실행: 실계정 로그인·메시지 왕복, 로그인 상태의 언어 저장(`POST /home/language` 성공 경로), 인증된 새 화면, pull/ACK/token 소비, 실제 rollback, Access·DNS 변경, 부하·장시간 검사. 이 항목은 후속 실설치·실계정 검증 과제에서 다룬다.
- 한계: 브라우저 시각 확인은 하지 않았다(UI 코드 변경 없음, 지시서 기준). 공개 페이지는 HTTPS 내용·헤더만 확인했다. rollback 이미지는 파일 존재·hash까지 확인했고 앱 기동은 하지 않았다. 새 relay는 manifest list ID로 보인다. rollback 시 새 `Member.Locale`·`Friend.CodeTail` 등 필드는 이전 코드가 무시하며, `Held` 계열 보류 행은 이전 코드에서 받은 사람에게 보일 수 있다(D13 DEV 인계 기록).
- Windows `make`가 없어 `lint.py`의 `make lint`·`make test`는 실행 불가로 남는다. 같은 한계는 직전 배포와 같다. 아래 절의 결과를 따른다.
- 후속(DEV): 새 공개 경로를 추가하는 변경은 `public-ingress.yml`·`check_public_ingress.py`를 함께 바꿔야 한다는 직전 후속이 이번에도 해당했다. 경로를 추가한 DEV 단위(`GET /lang/{code}`)가 배포 전 단위 gate로 잡히지 않았다. Grok 로컬 bundle과 Claude Code hook 패키지는 서버 배포로 갱신되지 않는다.

## 문서 HEAD의 FullOps 검사

work.py finish 전에 기록한다.

## 보존 위치

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/release002-ops/`(bundle·스크립트·원 로그).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/release002-ops-20261020/`(bundle·관측·backup.list·config 백업·후보·ingress-check 사본).
- 서버 백업: `/home/shin/deploy/knowslink-state/backups/1dd6800-pre-release002-20261020T135049Z.dump`(0600).
- 보존 컨테이너: `knowslink-release002-rehearse-20261020`(복원 fixture), `knowslink-release002-rollback-inspect-20261020`(rollback artifact 검사). 이미지 tag `knowslink-rollback:1dd6800-release002-20261020`.
- Git 증거: 이 폴더의 redacted JSON/log와 [commands.txt](commands.txt)(실행 스크립트 전문). 원 credential·token·cookie 값·실메시지·DB row는 포함하지 않는다. `?token=x`는 무인증 거부를 보는 합성 경로다.
