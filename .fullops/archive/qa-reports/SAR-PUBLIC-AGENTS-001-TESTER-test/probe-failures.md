---
title: SAR-PUBLIC-AGENTS-001-TESTER 프로브 중간 실패
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-TESTER]
summary: 독립 프로브의 중간 실패 원문과 재실행 이유를 보존한다
---

# 독립 프로브의 중간 실패 원문

최종 `probe-go.log`는 통과 실행이다. 아래 원문은 그 전에 같은 후보에서 프로브 기대값이 어긋나 실패한 실행에서 캡처한 fatal이다. 제품 소스는 바꾸지 않았고, 재실행 이유는 프로브만 고친 것이다. 중간 로그 파일은 다음 실행이 같은 경로에 덮어썼다.

## 게시된 Postgres에서의 첫 실패

- `cross_owner`: 연결 확인 GET에 `node adapters/dist/connect.js prepare`가 없다. 그 문장은 발급 POST의 연결 대기 화면에 있다.
- `relationships`: `invalid_schema`. `wire`의 idempotency key가 16자보다 짧았다.
- `capacity`: `got 409 want 303`. 전역 200이 찬 상태에서 owner 5를 만들려고 했다.
- `malformed`: `got 422 want 429`. 회원 예산은 40인데 31회째를 429로 기대했다.
- 같은 실행에서 키 회전과 Node CLI 서브테스트는 PASS였다.

## 기대값을 맞춘 뒤의 실패

```
qa_agents_probe_test.go:239: /v1/connect/info got 200 want 401: {"agent":"[redacted]","client":"node-local","exp":"[redacted]","mode":"register","owner":"[redacted]","state":"prepared"}
qa_agents_probe_test.go:421: /v1/send got 403 want 200: {"error":"[redacted]"}
qa_agents_probe_test.go:553: got 303 want 409: /home
```

재실행 이유: 같은 owner의 agent로 저장 상태를 바꾸면 info가 200이다. 수신 agent에 키가 없으면 send는 403이다. 비활성 owner의 pending은 전역 200에 들어가지 않아 invite가 303이다.

```
qa_agents_probe_test.go:197: /v1/connect/prepare got 200 want 422: {"state":"prepared"}
qa_agents_probe_test.go:439: got 403 want 303: 문제: 처리할 수 없습니다. 대상과 현재 권한·상태를 확인하세요. 새 권한은 생성되지 않았습니다.
```

재실행 이유: `browserAgent`가 새로 만든 id가 아니라 화면에서 사전순으로 마지막 `<code>agent_`를 반환해 치환과 초대가 같은 id에 걸렸다.

```
qa_agents_probe_test.go:456: got 403 want 303: 문제: 처리할 수 없습니다. 대상과 현재 권한·상태를 확인하세요. 새 권한은 생성되지 않았습니다.
```

재실행 이유: 관계가 생긴 뒤 홈의 마지막 agent 코드는 상대 수신 id다. 같은 owner 초대가 그 id를 사용했다. 이후에는 `<h3><code>`의 자기 agent만 새로 고른다.

## 최종

`PROBE_EXIT:0`. `TestQAAgentsIndependent`와 하위 6개 서브테스트가 PASS다.
