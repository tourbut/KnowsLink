---
title: SAR-MOTION-INTRO-001-DEV a2 검수 기록 — 단일 도형 UI 모션
status: draft
updated: 2026-10-18
owner: dev
tasks: [SAR-MOTION-INTRO-001-DEV]
summary: Playwright + ffmpeg 직접 렌더 v01–v03의 규격·루프·비트·시각·오디오 검수와 깨끗한 폴더 재렌더 결과.
---

# SAR-MOTION-INTRO-001-DEV a2 검수 기록

- attempt: a767a2ea1444498ab7ef1eebe2c7aa97 / base 9fb42da / 실행 HEAD 2244558(작업 미커밋 상태, dirty) → 결과 커밋은 완료 보고 참조
- 도구: Windows 11 Pro 10.0.26300, Node v24.13.0, Playwright 1.64.1(Chrome Headless Shell 156.0.8078.4, `node_modules` 안), FFmpeg 8.0.1 full build(gyan.dev), Python 3.13.12 venv + numpy 2.5.4
- 엔진: 사용자 프롬프트 지정대로 Playwright + ffmpeg 직접 렌더. HyperFrames·Remotion 공식 스킬 검사 항목은 해당 없음.
- 최종 후보: `videos/knowslink-motion/renders/v03/knowslink-motion.mp4` (2,315,100 B, sha256 `ca95b79ecc11239d8b676b0ceefcda19db1965c9631d43f0118df16e79e926f2`). v01·v02는 수정 근거로 보존.

## 버전 이력

| 버전 | 변경 | 판정 근거 |
|---|---|---|
| 비트 r1 (미렌더) | 첫 28장: 토글 노브가 트랙 뒤에 가려짐, 토스트 문구와 노브 겹침, 화면 채움 부족(폭 1100px), 체크 그리기 느림 | `videos/knowslink-motion/tmp/`(커밋 안 함) |
| v01 | 위 수정 + 13.9초 이후 정착 상태 평가로 원본 루프 동일. CRF 18 | 프레임 0/839 인코딩 차이 PSNR 52.1 dB, 최대 61/255(`loop-diff-v01.txt`). 7.94초 장거리 커서 이동이 4 subframe 잔상으로 갈라짐 |
| v02 | 마지막 프레임 I-frame 강제, 커서 이동 7.8초 시작 | PSNR 50.8 dB, 최대 63 — CRF/mbtree가 첫·끝 I-frame QP를 다르게 줌(`loop-diff-v02.txt`) |
| v03 | 고정 QP 18, 커서 장거리 이동을 느린 스프링(7.72초 시작) | 아래 표 |

## 검수 표 (v03)

| 검사 | 결과 | 근거 |
|---|---|---|
| 실행 | passed | `node render.mjs video v03` exit 0, 76초(`render-v03.log`) |
| 규격 | passed | `ffprobe-v03.json`: h264 High 1440×1440 yuv420p, 60/1, `nb_read_frames` 840, 14.000000초 / aac LC 44.1 kHz mono 14.000000초 |
| seek(t) 순수성 | passed | `selfcheck-v03.txt`: 6개 시각에서 다른 시각을 거친 뒤 같은 t로 돌아와 PNG 바이트 동일 |
| 루프 동일성(원본) | passed | 같은 파일: seek(0)과 seek(839/60) PNG 바이트 동일. 프레임 0의 subframe(±3.1 ms)도 스타일 차이는 숨겨진 스피너 회전뿐 |
| 루프 동일성(인코딩 MP4) | passed (근사) | `frames-v03/frame0000.png` vs `frame0839.png`: PSNR 62.68 dB, 최대 차 8/255, 바이트 동일 아님(`loop-diff-v03.txt`). 원인은 H.264 양자화이며 원본은 동일 |
| 비트 그리드(28장) | passed | `beat-frames-v03/b01..b28.png`(비트 + 0.25초), `beat-frames-v03-sheet.png`. 각 비트에 브리프 그리드 상태가 하나씩 있음. 매핑은 `beats.json`의 `grid` |
| 시각 대표 프레임 | passed | `frames-v03/`: 0, 45, 80, 110, 220, 225, 255, 300, 345, 400, 450, 459, 471, 490, 540, 565, 590, 640, 690, 745, 760, 800, 830, 839. 전환 직후(220·490·565·690·760·800)는 의도된 blur·이동 중 |
| 타이포/내용 | passed | 화면 문구는 브리프 그리드 문자열만: Google로 계속, 내 이름, 지우(ㅈ·지·지ㅇ 조합 단계 포함), Grok, Dots, Claude Code, 연결됨, 친구 요청, 민지와 서로 수락, 친구, 내 에이전트, 상대 에이전트, 메모 도우미, 민지 / 검토 도우미, 소개 문구 초안을 확인해 주세요., 보내기, 승인, 민지의 검토 도우미에 도착. 친구 코드는 `•••• •••• ••••` 마스킹. 잘림·겹침·한글 glyph 누락 없음(직접 확인). `document.fonts`에서 Geist·Pretendard 둘 다 loaded |
| 금지 요소 | passed (자체 판정) | 바운스 이징 없음(모든 스프링 ζ ≥ 0.8, 오버슛 ≤ 1.5 %), 파티클·글로우·UI 크롬 그라디언트 없음, 아이콘 선 굵기 하나(5 world px), 카메라 스케일 요소에 `will-change` 없음, accent `#005A62` 하나 |
| 음원 비트 | passed | `beats.json`: House Vibes 측정 119.995 BPM, 강도 12.69, 구간 시작 64.0838초(32마디째 다운비트). 28비트의 실제 onset 피크가 그리드에서 최대 10.1 ms(onset 해상도 11.6 ms 이내). 늘이기 미적용(14초 누적 오차 < 1 ms) |
| 효과음 정렬 | passed | `audio-align-v03.json`, `sfx-placement-v03.json`: 37개 이벤트. 연속 타이핑 구간 밖 23개는 측정 피크가 계획 시각에서 최대 0.5 ms. 10.0–10.43초 24 ms 간격 타이핑 14개는 앞 키의 꼬리와 겹쳐 개별 onset 측정 불가(배치는 구성상 피크 = 시각) |
| 최종 믹스 | passed | 최종 AAC ↔ 믹스 WAV 지연 0.0 ms(교차상관). Peak −1.1 dBFS, RMS −14.6 dB(`audio-stats-v03.txt`), 리미터 ×0.71(합산 피크 1.256 → 0.89). 루프 경계 6 ms 페이드 |
| 청취 | not_run | 사용자 청취 수락 몫. 음량 균형·효과음 톤·곡 선택 미청취 |
| 비용/권리 | passed | 비용 0(로컬 렌더). 글꼴 OFL 1.1(커밋), Mixkit Music/SFX Free License(재배포 금지로 mp3 미커밋, `fetch-audio.sh` + `SHA256SUMS`). `videos/knowslink-motion/assets/SOURCES.md` |
| 재렌더 | 아래 절 | 깨끗한 폴더 |

## 깨끗한 폴더 재렌더

| 항목 | 결과 |
|---|---|
| 폴더 | 세션 scratchpad `clean-rerender/`(저장소 밖). 소스·`package-lock.json`·`beats.json`·글꼴만 복사, `node_modules`·`.venv`·mp3 없음 |
| 설치 | `npm ci` 0 vulnerabilities, `pip install numpy==2.5.4`, `sh fetch-audio.sh` 7개 sha256 OK |
| 1차 실행 | failed: `PLAYWRIGHT_BROWSERS_PATH=0`으로 받은 Chromium이 긴 경로(Windows MAX_PATH 260자) 때문에 풀리지 않아 `chrome-headless-shell.exe ENOENT`(`clean-rerender/rerender.log`) |
| 2차 실행 | passed: 같은 Playwright 1.64.1 브라우저 빌드 1248을 `PLAYWRIGHT_BROWSERS_PATH`로 원 작업 폴더에서 지정(소스·lock·자산은 깨끗한 폴더 것). selfcheck 통과, `video v01` exit 0, 79초(`clean-rerender/rerender2.log`) |
| 규격 | passed: h264 1440×1440 yuv420p 60/1 840프레임 14.000초 + aac 44.1 kHz 14.000초(`clean-rerender/ffprobe-clean.json`), 2,319,384 B, sha256 `2aa69d20…`(바이트 동일성 요구 없음) |
| 시각 대조 | passed: v03과 840프레임 전체 PSNR 평균 63.35 dB, 최소 49.45 dB(`clean-rerender/psnr-vs-v03.txt`). 대표 프레임 0·400·640·839(`clean-rerender/frames-0-400-640-839.png`) 직접 확인 동일 |
| 오디오 대조 | passed: 디코딩 PCM 바이트 동일 |

## 미확인·사용자 판단

- 청취: 음악·효과음 음량과 톤, 곡 선택(대안 Dry Gin 120.005 BPM·Shame 120.005 BPM은 `beats.json`만 바꿔 재렌더 가능).
- 움직임 체감: 4 subframe 블러에서 빠른 커서는 여전히 약한 잔상이 보일 수 있다(프롬프트 지정 subframe 수 유지).
- 오디오는 mono 믹스다.
