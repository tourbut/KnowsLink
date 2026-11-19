---
title: SAR-LANDING-GUIDE-001-TESTER 조작 테스트 — 통과 (passed)
status: draft
updated: 2026-10-14
owner: tester
tasks: [SAR-LANDING-GUIDE-001-TESTER]
summary: "Read the public introduction, steps, installation summaries, fictional examples and login section. Do not log in or follow external links."
---

# SAR-LANDING-GUIDE-001-TESTER 조작 테스트 — 통과 (passed)

- 시나리오: `.fullops-squad/docs/evaluations/scenarios/web/landing-public-tester.json`
- covers: LG-01 public structure, LG-03 order, LG-04 install summaries, LG-05 fictional examples, LG-06 delivery limits, LG-07 both providers
- 목표: Read the public introduction, steps, installation summaries, fictional examples and login section. Do not log in or follow external links.
- 스텝 0 · Jev 0회 · 비용 0 · Jev 시간 0초

## 검증 조건

| 조건 | 읽은 값 | 결과 |
|---|---|---|
| js `['main','how','install','examples','start'].every(id=>!!document.getElementById(id))` == True | True | 통과 |
| js `document.querySelectorAll('h1').length===1 && document.documentElement.lang==='ko'` == True | True | 통과 |
| js `document.querySelectorAll('details').length===2 && [...document.querySelectorAll('details')].every(e=>!e.open)` == True | True | 통과 |
| 문구 "기능을 설명하기 위한 가상 예시입니다."가 보여야 함 | True | 통과 |
| 문구 "최대 180초"가 보여야 함 | True | 통과 |
| 문구 "최대 24시간"가 보여야 함 | True | 통과 |
| js `document.querySelectorAll('#start form').length===2 && document.querySelectorAll('form').length===2` == True | True | 통과 |

실행기: `orca` / `1.4.224` · 범위: `assertions_only`
관측: `{"total": 38, "offset": 0, "observed": 38, "remaining": 0, "unsupported": [], "complete": true}`
중단 사유: `passed` · 독립 assertions: `True`

## 진행

| 스텝 | 동작 | 대상 | 입력값 키 | 확신 | 결과 |
|---|---|---|---|---|---|

원본: `result.json`(판정), `events.jsonl`(스텝별 선택·확률·비용).
