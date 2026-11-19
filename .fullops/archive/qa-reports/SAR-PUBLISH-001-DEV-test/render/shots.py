"""Capture each fixture page at 1440/390/320 inside a fixed-width iframe and crop to that width."""
import pathlib, subprocess, sys, tempfile
from PIL import Image
CHROME = r'C:\Program Files\Google\Chrome\Application\chrome.exe'
src = pathlib.Path(sys.argv[1]); out = pathlib.Path(sys.argv[2]); out.mkdir(parents=True, exist_ok=True)
pages = ['policy-privacy-ko', 'policy-privacy-en', 'policy-terms-ko', 'policy-terms-en', 'policy-start-ko', 'policy-start-en', 'human-account', 'en-human-account', 'gate-active', 'en-gate-active']
with tempfile.TemporaryDirectory() as tmp:
    for p in pages:
        for w in (1440, 390, 320):
            hgt = 1400 if w == 1440 else 2400
            wrap = pathlib.Path(tmp) / f'{p}-{w}.html'
            wrap.write_text(f'<!doctype html><body style="margin:0"><iframe style="width:{w}px;height:{hgt}px;border:0;display:block" src="{(src / (p + ".html")).as_uri()}"></iframe></body>', encoding='utf-8')
            shot = out / f'{p}-{w}.png'
            subprocess.run([CHROME, '--headless=new', '--disable-gpu', '--hide-scrollbars', f'--user-data-dir={tmp}/u', f'--window-size={max(w, 500)},{hgt}', '--virtual-time-budget=3000', f'--screenshot={shot}', wrap.as_uri()], capture_output=True, timeout=60)
            Image.open(shot).crop((0, 0, w, hgt)).save(shot)
print(len(list(out.glob('*.png'))))
