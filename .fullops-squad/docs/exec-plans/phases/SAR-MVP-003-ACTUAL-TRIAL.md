---
title: 실제 Grok 시험 메시지 송신과 회신 수신 기록
status: draft
updated: 2026-10-05
owner: coor
tasks: [SAR-MVP-003-BIDIRECTIONAL]
summary: Grok 파일 설치와 준비 완료 회신 후 실제 relay에서 Codex 송신과 trial_grok 회신 수신을 대조한 결과를 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL — 실제 시험 round

사용자는 Grok 파일 설치 댓글 확인을 요청했다. 기존 실제 메시지 시험 승인을 유지한다. 제품은 검수된0911c2c이며 제품 변경은 없다. 대상은 기존https://link.knowslog.com, 시험text만 사용했다. 업무pull·자동wake·dots는 수행하지 않았다.

## Grok 준비 근거

[준비 완료 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993472191): private폴더0700 box, key/environment0600 box, HEAD0911c2c, receive exit0/received_or_empty/message:null. 설치와 파일전달은완료됐다고Bot이회신했다. 비밀내용은공개하지않았다.

## 직접 실행한 송신과 수신

coor가 `scripts/run_trial.py`와trial_codex의privateconfig로send/receive를실행했다. trial_grok config로로컬대역송신을하지않았다. 기존180초TTL과수동pull을사용했다.

| 항목 | 관측 |
|---|---|
| Codex send key | sar-mvp-003-codex-actual-20261005-1127 |
| Codex send ID | 01a10bd1-0aa4-7f2a-88dc-2af44cdd66d6 |
| from/to | trial_codex → trial_grok |
| accepted_at | 2026-10-05T11:26:51.677202Z |
| exp | 2026-10-05T11:29:51Z |
| send 판정 | queued. 이 결과만으로Grok수신성공이라고판정하지않았다. |
| Codex가 수신한 reply ID | 01a10bd1-164b-79ae-9c7e-337807e0c1dc |
| reply from/to | trial_grok → trial_codex |
| reply text | Grok trial reply to 01a10bd1-0aa4-7f2a-88dc-2af44cdd66d6 |
| reply exp | 2026-10-05T11:29:54Z |
| reply 처리 | untrusted:true 데이터, 명령실행없음 |

원래송신ID가회신text에일치한다. from/to도기대값과일치한다. nonce는송신문에는있지만회신은기존안내형식으로ID만참조하므로회신nonce일치라고기록하지않는다. Codex관측으로relay에서인증된trial_grok의관련회신수신까지완료했다. Grok측실제receive/send댓글의ID대조는추가보고대기다.

[송신ID 안내](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993505386)·[Codex 수신확인](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993526732)을게시했다. 같은key추가회신을요청하지않았다. actual왕복확정근거를보존하고Grok측보고의from/to·수신ID·replyID를추가대조한다.

## 후속과 운영 범위

시험service자격은Codex2026-10-06T08:09:43Z/Grok08:09:50Z까지다. 관리API token은2026-10-05T23:59:59Z에먼저만료된다. 자동wake/MCP부모도구4개갱신/dots연결은미검증이며이번manualCLIround와분리한다. 시험종료시우리trialtoken·앱·정책정리와원점owner복귀/allowlist비움은D12의종료절차를따른다. userprivate파일은자동삭제하지않았다.

## Grok측 최종 대조 — 성공 확정

[실제 결과 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993585072)을 확인했다. Grok수신ID가Codex sendID와같고, GrokreplyID가여기서받은ID와같다. 양쪽from/to도같다. Grok수신nonce가원래값과같다. 회신문구는nonce없이ID만참조하며고정문구사용을Bot이보고했다. receive/send각exit0, reply accepted_at2026-10-05T11:26:55.048742Z, key sar-mvp-003-grok-renew-round-1 한 번이다. 10초간격수동CLIreceive루프는종료됐다고보고했다.

실제Codex→Grok Bot→Codex 수동CLI왕복을성공으로확정했다. 이전추가보고대기는해소됐다. [성공 판정 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993750758)을게시하고Grok측private시험파일제거회신을요청했다. 기존승인된D12시험종료절차의Accesstrial자원/원점/allowlist정리는OPS에게인계한다. 부모MCP도구목록과자동wake/dots는별도후속이며이슈를열어둔다.

## 종료 확인

Grok은 [정리 완료 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993825445)로 private 시험 키·설정과 받은 압축파일 제거를 보고했다. OPS 고정 SHA `2182401680806f4938cf8ad46a9f04c644e60bf5`에서 서버 시험 인증 삭제·owner 보호 복원·공유 서비스 회귀를 통과했다. coor가 [종료 정리 결과](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993872965)를 게시하고 원문을 대조했다. 실제 수동 CLI 왕복과 이번 시험 종료는 완료다. MCP 부모 도구 목록·자동 wake·dots는 후속으로 남긴다.
