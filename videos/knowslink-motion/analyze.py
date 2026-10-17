"""Beat grid analysis with numpy only.

usage:
  python analyze.py scan <mp3>...           # BPM/strength per file
  python analyze.py pick <mp3> [--out beats.json] [--start <sec>]
  python analyze.py peaks <wav-or-mp3>...   # first transient peak per SFX (ms)
"""
import json
import subprocess
import sys

import numpy as np

SR = 22050
HOP = 256
BEATS = 28
TARGET_BPM = 120.0
GRID = [  # brief: SAR-MOTION-INTRO-001-brief.md, beat k at video t = (k-1) * 0.5 s
    "b01 ① 버튼 Google로 계속", "b02 커서 클릭", "b03 ② 로더", "b04 ③ 체크",
    "b05 ④ 이름 입력칸 내 이름", "b06 지우 타이핑 (ㅈ, 지)", "b07 타이핑 계속 (지ㅇ, 지우)", "b08 Enter → ✓",
    "b09 ⑤ 에이전트 카드 Grok·Dots·Claude Code", "b10 Claude Code 클릭", "b11 로더", "b12 연결됨",
    "b13 ⑥ 친구 코드 붙여넣기(마스킹)", "b14 친구 요청", "b15 ⑦ 토글 민지와 서로 수락", "b16 토글이 늘어나며 켜짐",
    "b17 ⑧ 노브가 탭 표시자로 친구·내 에이전트·상대 에이전트", "b18 메모 도우미로 드래그(놓기)", "b19 민지 / 검토 도우미로 드래그(놓기)", "b20 ⑨ 작성창",
    "b21 소개 문구 초안을 확인해 주세요. 타이핑", "b22 보내기", "b23 ⑩ 승인 대화상자", "b24 승인",
    "b25 ⑪ 토스트 민지의 검토 도우미에 도착", "b26 정지 + 줌아웃", "b27 ⑫ 버튼으로 복귀", "b28 정착, 마지막 프레임 = 프레임 0",
]


def decode(path, sr=SR):
    raw = subprocess.run(
        ["ffmpeg", "-v", "error", "-i", path, "-ac", "1", "-ar", str(sr), "-f", "f32le", "-"],
        check=True, capture_output=True).stdout
    return np.frombuffer(raw, dtype=np.float32)


def onset_env(y, lo_hz=0, hi_hz=None):
    n_fft = 1024
    win = np.hanning(n_fft).astype(np.float32)
    frames = np.lib.stride_tricks.sliding_window_view(y, n_fft)[::HOP] * win
    mag = np.abs(np.fft.rfft(frames, axis=1))
    freqs = np.fft.rfftfreq(n_fft, 1 / SR)
    band = (freqs >= lo_hz) & (freqs <= (hi_hz or SR / 2))
    logm = np.log1p(100 * mag[:, band])
    flux = np.maximum(0, np.diff(logm, axis=0)).sum(axis=1)
    flux = np.concatenate([[0], flux])
    flux -= np.convolve(flux, np.ones(16) / 16, mode="same")  # local mean removal
    return np.maximum(flux, 0)


def comb_score(env, bpm, fps):
    """Best phase and score of a beat pulse train at this tempo."""
    period = 60.0 / bpm * fps
    n = int((len(env) - 1) / period)
    best = (-1.0, 0.0)
    for phase in np.arange(0, period, 0.5):
        idx = np.round(phase + period * np.arange(n)).astype(int)
        s = env[idx].mean()
        if s > best[0]:
            best = (s, phase)
    return best


def tempo(env):
    fps = SR / HOP
    ac = np.correlate(env, env, mode="full")[len(env) - 1:]
    lags = np.arange(len(ac))
    bpm_of = lambda lag: 60 * fps / lag
    ok = (lags > 60 * fps / 160) & (lags < 60 * fps / 90)
    lag = lags[ok][np.argmax(ac[ok])]
    coarse = bpm_of(lag)
    grid = np.arange(coarse - 2, coarse + 2, 0.01)
    scores = [comb_score(env, b, fps)[0] for b in grid]
    bpm = float(grid[int(np.argmax(scores))])
    score, phase = comb_score(env, bpm, fps)
    strength = score / (env.mean() + 1e-9)
    return bpm, phase / fps, strength


def scan(paths):
    for p in paths:
        y = decode(p)
        bpm, phase, strength = tempo(onset_env(y))
        print(f"{p}\t{len(y) / SR:6.1f}s\tbpm={bpm:7.2f}\tphase={phase:.3f}\tstrength={strength:.2f}")


def pick(path, out, start=None):
    y = decode(path)
    env = onset_env(y)
    low = onset_env(y, 0, 150)
    fps = SR / HOP
    bpm, phase, strength = tempo(env)
    period = 60.0 / bpm
    beats = np.arange(phase, len(y) / SR - period, period)
    # downbeat: the bar position (0..3) whose low-band onsets are strongest
    li = np.clip(np.round(beats * fps).astype(int), 0, len(low) - 1)
    bar_pos = int(np.argmax([low[li[k::4]].mean() for k in range(4)]))
    downbeats = beats[bar_pos::4]
    dur = BEATS * period
    energy = np.concatenate([[0], np.cumsum(y.astype(np.float64) ** 2)])
    if start is None:
        # loudest 14 s window that starts on a downbeat and keeps steady beats
        cands = [d for d in downbeats if d > 2.0 and d + dur < len(y) / SR - 1.0]
        def window_score(d):
            a, b = int(d * SR), int((d + dur) * SR)
            idx = np.round((d + period * np.arange(BEATS)) * fps).astype(int)
            rms = np.sqrt((energy[b] - energy[a]) / (b - a))
            return rms * (env[idx].mean() / (env.mean() + 1e-9))
        start = float(max(cands, key=window_score))
    rel = [float(round(b - start, 4)) for b in beats if start - 1e-6 <= b < start + dur - 1e-6]
    # strongest onset within +-46 ms of each grid beat: how well the song really sits on the grid
    w = int(0.046 * fps)
    onsets = []
    for r in rel:
        c = int(round((start + r) * fps))
        onsets.append(round((c - w + int(np.argmax(env[c - w:c + w + 1]))) / fps - start, 4))
    stretch = TARGET_BPM / bpm  # ffmpeg atempo factor so one beat = 0.5 s
    data = {
        "source": path.replace("\\", "/"),
        "measured_bpm": round(bpm, 3),
        "beat_strength": round(float(strength), 3),
        "atempo": round(stretch, 6),
        "segment_start_sec": round(start, 4),
        "segment_source_dur_sec": round(dur, 4),
        "downbeat_bar_position": bar_pos,
        "measured_beats_rel_source_sec": rel,
        "onset_peaks_rel_sec": onsets,
        "onset_offset_ms_max_abs": round(max(abs(o - r) for o, r in zip(onsets, rel)) * 1000, 1),
        "atempo_applied": False,  # |bpm-120| drifts < 1 ms over 14 s, so the song is not stretched
        "video_beat_sec": [round(i * 60 / TARGET_BPM, 3) for i in range(BEATS)],
        "grid": GRID,
    }
    if out:
        with open(out, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
    print(json.dumps({k: v for k, v in data.items() if not k.endswith("_sec") or k == "segment_start_sec"}))
    return data


def peak_index(y):
    """First sample reaching half of the clip's max: the audible transient."""
    a = np.abs(y)
    return int(np.argmax(a >= 0.5 * a.max()))


def peaks(paths):
    for p in paths:
        y = decode(p, 44100)
        print(f"{p}\tpeak_ms={peak_index(y) / 44.1:.1f}\tlen_ms={len(y) / 44.1:.0f}\tmax={np.abs(y).max():.3f}")


MIX_SR = 44100
SFX = {  # kind: file, max length (s), base gain
    "click": ("assets/audio/sfx-click-216.mp3", 0.12, 0.55),
    "key": ("assets/audio/sfx-key-2841.mp3", 0.10, 0.35),
    "switch": ("assets/audio/sfx-switch-2579.mp3", 0.16, 0.6),
    "notify": ("assets/audio/sfx-notify-2357.mp3", 0.38, 0.35),
}
MUSIC_GAIN = 0.7


def mix(beats_path, events_path, out_wav, dur=14.0):
    import wave
    beats = json.load(open(beats_path, encoding="utf-8"))
    events = json.load(open(events_path, encoding="utf-8"))
    n = int(round(dur * MIX_SR))
    music = decode(beats["source"], MIX_SR)
    a = int(round(beats["segment_start_sec"] * MIX_SR))
    out = MUSIC_GAIN * music[a:a + n].astype(np.float64)
    assert len(out) == n, "song too short for the segment"
    stem = np.zeros(n)
    clips = {}
    for kind, (path, max_len, gain) in SFX.items():
        y = decode(path, MIX_SR)[: int(max_len * MIX_SR)].astype(np.float64)
        y *= np.minimum(1, np.linspace(8, 0, len(y)))  # short fade at the trim point
        clips[kind] = (y * gain / np.abs(y).max(), peak_index(y))
    placed = []
    for e in events:
        y, pk = clips[e["kind"]]
        start = int(round(e["t"] * MIX_SR)) - pk  # the transient lands on the event time
        lo, hi = max(0, start), min(n, start + len(y))
        stem[lo:hi] += e.get("gain", 1) * y[lo - start:hi - start]
        placed.append({"t": e["t"], "kind": e["kind"], "start_sample": start, "peak_sample": start + pk})
    out += stem
    fade = int(0.006 * MIX_SR)  # avoid a click where the loop wraps
    out[:fade] *= np.linspace(0, 1, fade)
    out[-fade:] *= np.linspace(1, 0, fade)
    peak = np.abs(out).max()
    if peak > 0.89:  # keep the mix under -1 dBFS
        out *= 0.89 / peak
    for path, y in ((out_wav, out), (out_wav[:-4] + ".sfx.wav", stem / max(1e-9, np.abs(stem).max()))):
        pcm = (np.clip(y, -1, 1) * 32767).astype("<i2")
        with wave.open(path, "wb") as w:
            w.setnchannels(1)
            w.setsampwidth(2)
            w.setframerate(MIX_SR)
            w.writeframes(pcm.tobytes())
    json.dump(placed, open(out_wav + ".events.json", "w"), indent=1)
    print(f"mix {out_wav}: {n / MIX_SR:.3f}s peak_before_limit={peak:.3f} events={len(placed)}")


def verify(final, mix_wav, events_path):
    """UI sound transients vs planned times (on the SFX-only stem), and final-vs-mix delay."""
    stem = decode(mix_wav[:-4] + ".sfx.wav", MIX_SR)
    events = json.load(open(events_path, encoding="utf-8"))
    env = np.abs(stem)
    w = int(0.012 * MIX_SR)
    rows = []
    cs = [int(round(e["t"] * MIX_SR)) for e in events]
    for i, c in enumerate(cs):
        # window ends halfway to the neighbours so a dense typing burst does not bleed in
        lo = max(c - w, (cs[i - 1] + c) // 2 if i else 0)
        hi = min(c + w, (c + cs[i + 1]) // 2 if i + 1 < len(cs) else len(env))
        seg = env[lo:hi]
        onset = lo + int(np.argmax(seg >= 0.5 * seg.max()))  # first half-max crossing near the event
        rows.append(round((onset - c) / MIX_SR * 1000, 2))
    y_mix, y_fin = decode(mix_wav, MIX_SR), decode(final, MIX_SR)
    k = 4 * MIX_SR  # cross-correlate the first 4 s, lags within +-50 ms
    lags = range(-2205, 2206)
    a, b = y_mix[:k], y_fin[:k + 2205 * 2]
    scores = [float(np.dot(a[2205:k - 2205], b[2205 + L:k - 2205 + L])) for L in lags]
    delay = lags[int(np.argmax(scores))] / MIX_SR * 1000
    print(json.dumps({"events": len(rows), "sfx_onset_offset_ms_max_abs": max(map(abs, rows)),
                      "final_vs_mix_delay_ms": round(delay, 2), "final_samples": len(y_fin),
                      "sfx_onset_offset_ms": rows}))


if __name__ == "__main__":
    cmd, args = sys.argv[1], sys.argv[2:]
    if cmd == "scan":
        scan(args)
    elif cmd == "pick":
        out = args[args.index("--out") + 1] if "--out" in args else None
        st = float(args[args.index("--start") + 1]) if "--start" in args else None
        pick(args[0], out, st)
    elif cmd == "peaks":
        peaks(args)
    elif cmd == "mix":
        mix(*args)
    elif cmd == "verify":
        verify(*args)
