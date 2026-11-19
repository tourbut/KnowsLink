---
title: SAR-CI-002-DEV 렌더·대비 확인
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-CI-002-DEV]
summary: CI 파생 로고 렌더 확인과 WCAG 대비 수치
---

# SAR-CI-002-DEV 대비 수치 (WCAG 2.x 상대 휘도 비율, 비텍스트 기준 3:1)

| 전경 | 배경 | 비율 | 3:1 |
|---|---|---|---|
| `#F2F0EA` 반전 선·글자 | `#182B31` | 12.90 | 통과 |
| `#2FA3AD` 반전 청록 | `#182B31` | 4.87 | 통과 |
| `#F2F0EA` 반전 선·글자 | `#000000` | 18.43 | 통과 |
| `#2FA3AD` 반전 청록 | `#000000` | 6.96 | 통과 |
| `#182B31` 원본 잉크(참고) | `#FFFFFF` | 14.70 | 통과 |
| `#005A62` 원본 청록(참고) | `#FFFFFF` | 7.96 | 통과 |
| `#005A62` 원본 청록을 어두운 배경에 그대로(탈락 근거) | `#182B31` | 1.85 | 미달 |
| `#2FA3AD` 반전 청록 vs 반전 선(참고) | `#F2F0EA` | 2.65 | 미달 |

## 렌더 확인

- `favicon-16-32-48-light.png`·`-x4.png`: 파비콘(왼쪽 3개)과 원본 심볼(오른쪽 3개)을 16·32·48px 실제 크기로 렌더(`-x4`는 최근접 4배 확대). 원본은 16·32px에서 선 틈·점이 흐려지고 파비콘은 두 칸과 청록 겹침이 읽힌다.
- `favicon-16-32-48-dark.png`·`-x4.png`: 어두운 모드(`prefers-color-scheme: dark`)에서 파비콘이 반전 색으로 바뀐 모습.
- `lockup-light.png`: 영문·한글 lockup, 흰 배경.
- `lockup-inverse-dark.png`, `lockup-inverse-black.png`: 반전 lockup·심볼, `#182B31`과 `#000000` 배경.
- 재현: `node render.mjs <repo 절대경로> <출력 폴더 절대경로>`(Playwright는 `videos/knowslink-motion` 설치본). 브라우저 프로필은 실행 뒤 삭제했다.
- 한글 글자는 Pretendard 500(`videos/knowslink-motion/assets/fonts/PretendardVariable.woff2`, SIL OFL 1.1)에서 윤곽선으로 뽑았다. 렌더에서 폰트 의존 없음.
- 시각 수락은 하지 않았다. 사용자 또는 designer 확인이 필요하다.
