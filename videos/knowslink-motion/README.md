# KnowsLink 단일 도형 UI 모션 (SAR-MOTION-INTRO-001)

1440×1440, 60fps, 14.0초(840프레임), H.264 yuv420p + AAC. 120 BPM 7마디 = 28비트 × 0.5초. 브리프: `.fullops-squad/docs/design-docs/motion/SAR-MOTION-INTRO-001-brief.md`.

## 파일

| 파일 | 역할 |
|---|---|
| `index.html` | 영상 전체. `window.seek(t)`가 모든 스타일을 시간의 순수 함수로 계산한다(전환·타이머·프레임 간 상태 없음). 스프링은 closed-form step response의 합. |
| `render.mjs` | Playwright 렌더: 프레임당 4 subframe(180° 셔터) → ffmpeg `tmix` 모션 블러 → 60fps, 음원 믹스 mux. 비트 프레임·순수성 검사도 같은 진입점. |
| `analyze.py` | numpy 비트 분석(`scan`, `pick`), 효과음 피크(`peaks`), 믹스(`mix`), 최종 오디오 정렬 측정(`verify`). |
| `beats.json` | 측정 BPM, 사용 구간 시작점(다운비트), 비트 시각, 28비트 매핑. |
| `assets/` | 글꼴(커밋), 음원 출처·해시(`SOURCES.md`, `SHA256SUMS`). mp3는 라이선스상 커밋하지 않는다. |
| `renders/vNN/` | 렌더 결과. 수정 렌더는 새 버전 폴더(이전 파일을 덮지 않음, `render.mjs`가 거부한다). |

## 설치 (이 폴더 안에만)

```sh
npm ci                                              # playwright 1.64.1 (package-lock.json)
PLAYWRIGHT_BROWSERS_PATH=0 npx playwright install chromium   # node_modules 안에 Chromium
python -m venv .venv && .venv/Scripts/python -m pip install numpy==2.5.4   # Linux/macOS: .venv/bin/python
sh fetch-audio.sh                                   # Mixkit mp3 받기 + sha256 검증
```

PATH에 `ffmpeg`/`ffprobe`가 필요하다(검증: 8.0.1 full build). 검증 환경: Windows 11, Node v24.13.0, Python 3.13.12.

Windows에서 폴더 경로가 길면 `node_modules` 안 Chromium 압축 해제가 MAX_PATH(260자)를 넘어 `chrome-headless-shell.exe ENOENT`로 실패한다. 그때는 짧은 경로에 받고 렌더할 때도 같은 값을 준다: `PLAYWRIGHT_BROWSERS_PATH=C:/pw npx playwright install chromium`, `PLAYWRIGHT_BROWSERS_PATH=C:/pw node render.mjs ...`.

## 렌더

```sh
node render.mjs selfcheck              # seek 순수성 assert + 프레임0 = 프레임839
node render.mjs beats <dir>            # 비트당 1프레임(b01..b28, 비트 + 0.25초)
node render.mjs stills <dir> 0,6.25    # 임의 시각 PNG
node render.mjs video v02              # renders/v02/knowslink-motion.mp4
```

최종 후보는 `renders/v03/`(H.264 고정 QP 18, 마지막 프레임 I-frame). v01·v02는 수정 근거로 보존한다. 검수 기록: `.fullops-squad/docs/evaluations/qa-reports/SAR-MOTION-INTRO-001-DEV-test/a2/report.md`.

`video`는 먼저 페이지의 `window.EVENTS`(클릭·키·토글·알림 시각)를 `tmp/events.json`으로 내보내고 `analyze.py mix`로 효과음의 측정 피크가 그 시각에 오도록 믹스한 뒤 프레임을 렌더한다.

## 음원 바꾸기

선택: House Vibes(측정 119.995 BPM, 구간 64.084초 = 32마디째 다운비트). 대안은 Dry Gin(120.005), Shame(120.005). 바꾸려면:

```sh
.venv/Scripts/python analyze.py pick assets/audio/mixkit-dry-gin-868.mp3 --out beats.json   # 구간 자동 선택
.venv/Scripts/python analyze.py pick assets/audio/mixkit-dry-gin-868.mp3 --out beats.json --start 63.9674   # 오프셋 지정
node render.mjs video v02
```

영상 타임라인은 0.5초 그리드에 고정이라 음원 쪽 `segment_start_sec`만 바뀐다.

## 루프

모든 스프링은 13.9초 전에 0.02px 이하로 정착한다. `seek`은 13.9초 이후를 정착 상태(t=0)로 평가하므로 마지막 프레임(839)이 프레임 0과 같다(커서 위치·속도 포함).
