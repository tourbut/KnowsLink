---
title: SAR-HOST-VERIFY-001-LIVE 중간 보고 — Claude Code 작업 세션 C1–C4
status: draft
updated: 2026-10-21
owner: tester
tasks: [SAR-HOST-VERIFY-001-LIVE]
summary: 운영 7b4f9bc 대상 Claude Code 작업 세션 신규 설치 C1–C4를 검증했다. 사용자 중지로 C5 이후와 Grok G1–G9는 not_run이다
---

# SAR-HOST-VERIFY-001-LIVE 중간 보고

대상은 운영 `https://link.knowslog.com`의 제품 `7b4f9bcf11bb7ab13207d8e635060a6e90e7f8c2`다. 절차는 [시나리오](../../scenarios/SAR-HOST-VERIFY-001-LIVE.md)를 따랐다. 단계별 값은 [evidence.json](evidence.json)에 있다.

사용자는 2026-10-21 C5 Google 로그인 전에 '그만'으로 실행을 중지했다. 사용자 결정 D015에 따라 C1–C4를 중간 보고로 커밋하고 종료한다. 두 호스트 모두 '신규 설치 검증 완료'가 아니다.

## 판정 — Claude Code 작업 세션

환경: Claude Code 2.1.304, Windows 11 Pro 10.0.26300, Node v22.22.2, npm 10.9.7, Python 3.13.12, Git 2.52.0.windows.1. 설치 루트는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/host-verify-001`이다.

| ID | 단계 | 판정 | 근거 |
|---|---|---|---|
| C1 | 준비물 기록 | PASS | evidence.json `prerequisites` |
| C2 | 기본 `git clone` | FAIL (exit 128) | [c2-clone-default-fail.log](c2-clone-default-fail.log) |
| C2 | `git clone -c core.longpaths=true` 뒤 `7b4f9bc` detach | PASS (exit 0, `git status --short` 0줄) | [c2-clone-longpaths.log](c2-clone-longpaths.log), [c2-checkout.log](c2-checkout.log) |
| C2 | `npm ci` | PASS (exit 0) | [c2-npm-ci.tail.log](c2-npm-ci.tail.log) |
| C2 | build | PASS (exit 0) | [c2-npm-build.log](c2-npm-build.log) |
| C2 | `package_plugin.py --host claude --connect --verify` | PASS (exit 0) | [c2-package-claude.log](c2-package-claude.log) |
| C3 | `claude plugin validate` | PASS (exit 0) | [c3-plugin-validate.log](c3-plugin-validate.log) |
| C3 | `claude plugin test` | PASS (11 pass, 0 fail) | [c3-plugin-test.log](c3-plugin-test.log) |
| C4 | 새 작업 폴더 `claude --plugin-dir` 도구 검색 | PASS (9개 도구, `/knowslink-session status` 미연결 문구) | [c4-mcp-tools.txt](c4-mcp-tools.txt) |
| C5 | `knowslink_connect` 첫 연결 요청 | 만료 | 아래 발견 3 |
| C5–C11 | 인증·이름 등록·세션 on·친구·발송·수신 toast·재개 | not_run | 사용자 중지(D015) |

## 판정 — Grok Bot

| ID | 단계 | 판정 | 근거 |
|---|---|---|---|
| — | README 'Grok 신규 설치' 소스 고정 SHA | FAIL | 아래 발견 2 |
| G1–G9 | 앱 기록·설치·Command MCP·환경·인증·친구·수신·알림·재개 | not_run | 사용자 중지(D015). Grok 앱은 조작하지 않았다 |

## 발견

1. **Windows 기본 `git clone` 실패.** exit 128, `error: unable to create file .fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/.../...-integration.log: Filename too long`, `fatal: unable to checkout working tree`. 공개 안내대로 clone하면 Windows 기본 설정에서 체크아웃이 깨진다. 재현: Windows에서 `core.longpaths` 없이 `git clone https://github.com/tourbut/KnowsLink.git src`. 담당 제안: DEV·ops — 커밋된 runtime 로그 경로를 줄이거나 제거하고, 아니면 안내에 `core.longpaths`를 적는다.
2. **README Grok 설치 SHA 불일치.** README 'Grok 신규 설치'는 `2747fe98653d5c6946f41c219fc256395333d7fe`를 고정한다. 운영 제품은 `7b4f9bc`다. 사용자 결정에 따라 공개 안내 불일치 FAIL로 기록한다. 담당 제안: DEV(README)·ops(배포 SHA 동기화).
3. **C5 첫 연결 요청 만료.** `knowslink_connect {"confirmed":true}` 뒤 웹 페이지 지문이 세션 값과 일치했다. 사용자 조작이 없어 기한 2026-10-21 00:17:11 KST(00:17 KST)에 만료됐다. 제품 결함으로 판정하지 않는다.

## 범위와 남은 자원

운영 데이터 변경·로그인·친구 요청·발송·Grok 조작은 하지 않았다. 비밀값·토큰·private 연결 폴더 내용은 증거에 넣지 않았다. 아래 자원은 건드리지 않고 존재와 경로만 적는다.

- Orca 터미널 `hv-claude`: 검증용 Claude Code 세션이 열려 있다.
- 설치 루트 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/host-verify-001`: `inst-claude-01`, `inst-claude-01-private`, `logs`, `src`, `src-lp`, `work-claude`, `t.sh`.
- private 연결 폴더 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/host-verify-001/inst-claude-01-private`.

판정에 쓰지 않은 터미널 구동 스크립트는 `raw/`(git 제외)로 옮겼다. 재개하려면 C5부터 사용자 로그인 승인을 다시 받는다. 설계 U8 호스트 표 갱신은 이 보고로 바뀌지 않는다.
