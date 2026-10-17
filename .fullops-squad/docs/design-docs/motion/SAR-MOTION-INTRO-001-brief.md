---
title: SAR-MOTION-INTRO-001 단일 도형 UI 모션 브리프
status: draft
updated: 2026-10-17
owner: coor
tasks: [SAR-MOTION-INTRO-001-DEV]
summary: 사용자가 준 모션 프롬프트 원문, 확정 입력, 승인된 28비트 그리드.
---

# SAR-MOTION-INTRO-001 단일 도형 UI 모션 브리프

## 확정 입력 (사용자 2026-10-17)

- 상태: KnowsLink 흐름 12개(아래 그리드). 원문 예시의 음악 플레이어·차트·⌘K는 쓰지 않는다.
- 색: 흑백 + accent `#005A62` 하나. 캔버스는 밝은 warm-gray.
- 음악: Mixkit(상업 무료 라이선스)에서 약 120 BPM 곡을 받는다.
- 규격: 1440×1440 정사각, 60fps, 14초(120 BPM 7마디 = 28비트 × 0.5초), H.264 MP4 + AAC.
- 글꼴: Geist(라틴) + Pretendard(한글). 둘 다 OFL, 공식 배포처에서 받는다(승인됨).
- 그리드 승인: 사용자 "진행해"(2026-10-17).

## 비트 그리드 (비트당 0.5초, b1 = 첫 다운비트)

| 마디 | b1 | b2 | b3 | b4 |
|---|---|---|---|---|
| 1 | ① 버튼 `Google로 계속` | 커서 클릭 | ② 로더 | ③ 체크 |
| 2 | ④ 이름 입력칸 `내 이름` | `지우` 타이핑 | (타이핑 계속) | Enter → ✓ |
| 3 | ⑤ 에이전트 카드, 칩 `Grok` `Dots` `Claude Code` | `Claude Code` 클릭 | 로더 | `연결됨` |
| 4 | ⑥ 친구 코드 붙여넣기(마스킹) | `친구 요청` | ⑦ 토글 `민지와 서로 수락` | 토글이 늘어나며 켜짐 |
| 5 | ⑧ 노브가 탭 표시자로: `친구` `내 에이전트` `상대 에이전트` | `메모 도우미`로 드래그 | `민지 / 검토 도우미`로 드래그 | ⑨ 작성창 |
| 6 | `소개 문구 초안을 확인해 주세요.` 타이핑 | `보내기` | ⑩ 승인 대화상자 | `승인` |
| 7 | ⑪ 토스트 `민지의 검토 도우미에 도착` | 정지 + 줌아웃 | ⑫ 버튼으로 복귀 | 마지막 프레임 = 프레임 0 |

문구는 이 표와 첫페이지 정본(`docs/design-docs/mockups/service-member-ui.md`)의 예시 이름(민지·메모 도우미·검토 도우미)만 쓴다. 실제 계정·코드·키는 쓰지 않는다(친구 코드는 마스킹 표시).

## 사용자 프롬프트 원문

<inputs>
Ask me for: 8 to 12 UI states I want the shape to become (e.g. button, loader, player, slider, toggle, tabs, chart, command palette, toast), pure black and white or one accent color, and a royalty-free song around 120 BPM (e.g. Mixkit, free for commercial use).
</inputs>

<direction>
Dribbble-level UI motion. One shape, never cut: every state is the same element morphing its size, radius and color while its content swaps with a short blur. A cursor drives every change with real clicks and drags. Light warm-gray canvas, black and white components, one clean UI font (Geist). Springs everywhere, a tiny overshoot at most. The camera zooms so each state fills the frame. The last frame is the first frame, so it loops.
Banned: bouncy easing, particle bursts, glows, gradients on UI chrome, mismatched icon strokes, dead time, anything that looks like a template.
</direction>

<structure>
120 BPM, 7 bars, something happens on every beat.
Button → loader → check → dynamic island → music player with a play/pause morph → scrub the progress bar → it becomes a volume slider that stretches when dragged past max → a toggle flips on the beat → the knob becomes a liquid tab indicator → the tabs open into a chart that draws itself, with a tooltip on hover → it collapses into ⌘K → type to filter → enter → toast → back to the button.
</structure>

<build>
1. One HTML file, square 1440x1440. Every style is computed from time inside seek(t): no CSS transitions, no timers, no state carried between frames.
2. Springs are closed-form step responses. A value that changes target many times is the sum of one spring per change, so it stays a pure function of time.
3. The tab indicator's two edges ride different springs, so the leading edge stretches ahead of the trailing one. Same trick for the toggle knob.
4. Drags are direct manipulation: while the cursor is held, the value is computed from its position. On release it springs back from wherever it was.
5. Analyze the song with numpy for the beat grid and start on a downbeat. Place every UI sound by its measured peak.
6. Render with Playwright: 4 subframes per frame, blended with ffmpeg tmix for motion blur at 60fps.
7. Render one frame per beat before the full render. Fix anything off the grid, cramped or hard to read.
</build>

<gotchas>
Never put will-change on anything the camera scales or the text renders blurry. Text that swaps inside a morphing container needs its own enter and exit timing or it overlaps. Make the last frame identical to the first, cursor position and speed included, or the loop stutters.
</gotchas>
