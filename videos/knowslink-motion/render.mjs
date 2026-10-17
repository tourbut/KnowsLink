// usage:
//   node render.mjs stills <outDir> <t1,t2,...>   # PNG per time (seconds)
//   node render.mjs beats <outDir>                 # one frame per beat (b01..b28, beat + 0.25 s)
//   node render.mjs selfcheck                      # seek(t) purity assert
//   node render.mjs video <version>                # renders/<version>/knowslink-motion.mp4
import { spawn, spawnSync } from "node:child_process";
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

process.env.PLAYWRIGHT_BROWSERS_PATH ??= "0"; // browsers live in node_modules (installed per README)
const { chromium } = await import("playwright");

const ROOT = dirname(fileURLToPath(import.meta.url));
const SIZE = 1440, FPS = 60, SUB = 4, DUR = 14;
const PY = join(ROOT, ".venv", process.platform === "win32" ? "Scripts/python.exe" : "bin/python");

async function openPages(n) {
  const browser = await chromium.launch({ args: ["--force-color-profile=srgb", "--disable-lcd-text"] });
  const pages = [];
  for (let i = 0; i < n; i++) {
    const page = await browser.newPage({ viewport: { width: SIZE, height: SIZE }, deviceScaleFactor: 1 });
    await page.goto(pathToFileURL(join(ROOT, "index.html")).href);
    await page.evaluate(() => document.fonts.ready);
    const fonts = await page.evaluate(() => ({
      geist: document.fonts.check('600 40px "Geist"', "Google Claude Code Grok Dots"),
      pretendard: document.fonts.check('600 40px "Pretendard"', "민지의 검토 도우미에 도착 ㅈㅇ"),
      loaded: [...document.fonts].filter(f => f.status === "loaded").map(f => f.family),
    }));
    if (!fonts.geist || !fonts.pretendard || fonts.loaded.length < 2) throw new Error(`fonts not loaded: ${JSON.stringify(fonts)}`);
    if (i === 0) console.log("fonts", JSON.stringify(fonts));
    const cdp = await page.context().newCDPSession(page);
    pages.push({ page, cdp });
  }
  return { browser, pages };
}

async function shot({ page, cdp }, t) {
  await page.evaluate(x => window.seek(x), t);
  const { data } = await cdp.send("Page.captureScreenshot", { format: "png", optimizeForSpeed: true });
  return Buffer.from(data, "base64");
}

async function stills(outDir, times, names) {
  mkdirSync(outDir, { recursive: true });
  const { browser, pages } = await openPages(1);
  for (let i = 0; i < times.length; i++) {
    const file = join(outDir, names ? names[i] : `t${times[i].toFixed(3)}.png`);
    writeFileSync(file, await shot(pages[0], times[i]));
    console.log(file);
  }
  await browser.close();
}

async function selfcheck() {
  const { browser, pages } = await openPages(1);
  const p = pages[0], probes = [0.83, 4.62, 8.31, 8.77, 10.21, 12.74];
  let ok = true;
  for (const t of probes) {
    const a = await shot(p, t);
    await shot(p, 13.37 - t); // visit another time in between, then come back
    const b = await shot(p, t);
    const same = a.equals(b);
    ok &&= same;
    console.log(`seek purity t=${t}: ${same ? "same" : "DIFFERENT"}`);
  }
  // loop: the end of the timeline must equal its start
  const first = await shot(p, 0), last = await shot(p, (FPS * DUR - 1) / FPS);
  console.log(`loop frame0 vs frame${FPS * DUR - 1} png bytes equal: ${first.equals(last)}`);
  await browser.close();
  if (!ok) { console.error("selfcheck FAILED"); process.exit(1); }
  console.log("selfcheck passed");
}

function run(cmd, args) {
  const r = spawnSync(cmd, args, { stdio: "inherit", cwd: ROOT });
  if (r.status !== 0) throw new Error(`${cmd} exited ${r.status}`);
}

async function video(version) {
  const out = join(ROOT, "renders", version);
  if (existsSync(join(out, "knowslink-motion.mp4"))) throw new Error(`${out} exists; pick a new version`);
  mkdirSync(out, { recursive: true });
  mkdirSync(join(ROOT, "tmp"), { recursive: true });
  const workers = 4;
  const { browser, pages } = await openPages(workers);
  // audio first: UI sounds come from the same timeline the frames use
  const events = await pages[0].page.evaluate(() => window.EVENTS);
  writeFileSync(join(ROOT, "tmp", "events.json"), JSON.stringify(events, null, 1));
  const wav = join(ROOT, "tmp", `mix-${version}.wav`);
  run(PY, ["analyze.py", "mix", "beats.json", "tmp/events.json", wav]);

  const total = FPS * DUR * SUB;
  const ff = spawn("ffmpeg", ["-y", "-v", "error",
    "-f", "image2pipe", "-framerate", String(FPS * SUB), "-c:v", "png", "-i", "-",
    "-i", wav,
    "-filter_complex", `[0:v]tmix=frames=${SUB},select='eq(mod(n\\,${SUB})\\,${SUB - 1})',setpts=N/(${FPS}*TB),format=yuv420p[v]`,
    "-map", "[v]", "-map", "1:a", "-r", String(FPS), "-fps_mode", "cfr",
    // constant QP and an I-frame at the end keep the encoded loop seam identical to frame 0
    "-c:v", "libx264", "-preset", "slow", "-qp", "18", "-pix_fmt", "yuv420p", "-movflags", "+faststart",
    "-force_key_frames", `expr:eq(n,0)+eq(n,${FPS * DUR - 1})`,
    "-c:a", "aac", "-b:a", "192k", "-t", String(DUR),
    join(out, "knowslink-motion.mp4")], { stdio: ["pipe", "inherit", "inherit"] });
  const done = new Promise((res, rej) => ff.on("close", c => (c === 0 ? res() : rej(new Error(`ffmpeg exited ${c}`)))));

  // subframe k sits at frame time + centred 180-degree shutter offset
  const timeOf = k => Math.floor(k / SUB) / FPS + ((k % SUB) - (SUB - 1) / 2) / (FPS * SUB * 2);
  const ready = new Map();
  let next = 0, written = 0;
  const t0 = Date.now();
  const flush = async () => {
    while (ready.has(written)) {
      const buf = ready.get(written); ready.delete(written);
      if (!ff.stdin.write(buf)) await new Promise(r => ff.stdin.once("drain", r));
      written++;
      if (written % 480 === 0) console.log(`subframes ${written}/${total} (${((Date.now() - t0) / 1000).toFixed(0)} s)`);
    }
  };
  await Promise.all(pages.map(async p => {
    while (next < total) {
      const k = next++;
      while (k - written > 64) await new Promise(r => setTimeout(r, 5)); // bound memory
      ready.set(k, await shot(p, timeOf(k)));
      await flush();
    }
  }));
  await flush();
  ff.stdin.end();
  await done;
  await browser.close();
  console.log(`rendered ${join(out, "knowslink-motion.mp4")} in ${((Date.now() - t0) / 1000).toFixed(0)} s`);
}

const [mode, ...args] = process.argv.slice(2);
if (mode === "stills") await stills(resolve(args[0]), args[1].split(",").map(Number));
else if (mode === "beats") {
  const times = Array.from({ length: 28 }, (_, i) => i * 0.5 + 0.25);
  await stills(resolve(args[0]), times, times.map((_, i) => `b${String(i + 1).padStart(2, "0")}.png`));
} else if (mode === "selfcheck") await selfcheck();
else if (mode === "video") await video(args[0] ?? "v01");
else { console.error(readFileSync(fileURLToPath(import.meta.url), "utf8").split("\n").slice(0, 5).join("\n")); process.exit(2); }
