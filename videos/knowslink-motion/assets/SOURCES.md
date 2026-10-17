# 자산 출처

받은 날: 2026-10-17. 모든 파일은 공식 배포처에서 받았고 수정하지 않았다(음원은 렌더 때 구간만 잘라 쓴다).

## 글꼴 (커밋함)

| 파일 | 출처 | 버전 | 라이선스 |
|---|---|---|---|
| `fonts/Geist-Variable.woff2` | npm `geist` (Vercel) `https://registry.npmjs.org/geist/-/geist-1.8.1.tgz` → `dist/fonts/geist-sans/Geist-Variable.woff2` | 1.8.1 | SIL OFL 1.1, `fonts/Geist-OFL.txt` |
| `fonts/PretendardVariable.woff2` | npm `pretendard` (orioncactus/pretendard) `https://registry.npmjs.org/pretendard/-/pretendard-1.3.9.tgz` → `dist/web/variable/woff2/PretendardVariable.woff2` | 1.3.9 | SIL OFL 1.1 (Reserved Font Name Pretendard), `fonts/Pretendard-OFL.txt` |

OFL은 수정하지 않은 글꼴을 라이선스 파일과 함께 재배포하는 것을 허용한다.

## 음원·효과음 (커밋하지 않음)

Mixkit 라이선스(`https://mixkit.co/license/`)
- Stock Music Free License: 상업·비상업 영상과 웹·소셜 게시 허용. CD/DVD·게임·TV/라디오 방송·리믹스(음악 단독 트랙)·권리 등록 금지.
- Sound Effects Free License: 더 큰 결과물에 포함해 사용 허용. **항목 단독·소스 파일과 함께 재배포 금지.**

그래서 mp3는 저장소에 넣지 않고 `../fetch-audio.sh`로 다시 받는다(`SHA256SUMS`로 검증). 귀속 표기 의무는 없다.

| 파일 | 제목 / 작가 | 찾은 목록 페이지 | 원본 URL | 측정 BPM |
|---|---|---|---|---|
| `audio/mixkit-house-vibes-129.mp3` **사용** | House Vibes / Alejandro Magaña (A. M.) | `https://mixkit.co/free-stock-music/house/` | `https://assets.mixkit.co/music/129/129.mp3` | 119.995 |
| `audio/mixkit-dry-gin-868.mp3` 대안 | Dry Gin / Michael Ramir C. | `https://mixkit.co/free-stock-music/funk/` | `https://assets.mixkit.co/music/868/868.mp3` | 120.005 |
| `audio/mixkit-shame-553.mp3` 대안 | Shame / Eugenio Mininni | `https://mixkit.co/free-stock-music/pop/` | `https://assets.mixkit.co/music/553/553.mp3` | 120.005 |
| `audio/sfx-click-216.mp3` | Clear mouse clicks | `https://mixkit.co/free-sound-effects/click/` | `https://assets.mixkit.co/active_storage/sfx/216/216-preview.mp3` | 피크 40.3 ms |
| `audio/sfx-key-2841.mp3` | Single key type | `https://mixkit.co/free-sound-effects/keyboard/` | `https://assets.mixkit.co/active_storage/sfx/2841/2841-preview.mp3` | 피크 8.3 ms |
| `audio/sfx-switch-2579.mp3` | On or off light switch tap | `https://mixkit.co/free-sound-effects/switch/` | `https://assets.mixkit.co/active_storage/sfx/2579/2579-preview.mp3` | 피크 20.5 ms |
| `audio/sfx-notify-2357.mp3` | Positive notification | `https://mixkit.co/free-sound-effects/notification/` | `https://assets.mixkit.co/active_storage/sfx/2357/2357-preview.mp3` | 피크 39.4 ms |

Mixkit 목록 페이지는 곡별 상세 페이지 링크를 내보내지 않아 목록 페이지와 원본 파일 URL을 기록했다. 제목·작가·길이는 목록 페이지의 JSON-LD(`MusicRecording`)에서 읽었다.
