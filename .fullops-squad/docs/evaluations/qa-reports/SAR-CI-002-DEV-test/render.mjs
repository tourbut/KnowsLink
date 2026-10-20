// QA render: node render.mjs <repoRoot> <outDir>; browser profile stays under outDir
import { createRequire } from "node:module";
import { join } from "node:path";
import { readFileSync } from "node:fs";
const [root, out] = process.argv.slice(2);
process.env.PLAYWRIGHT_BROWSERS_PATH = "0";
const { chromium } = createRequire(join(root, "videos/knowslink-motion/package.json"))("playwright");
const svg = (n) => "data:image/svg+xml;base64," + readFileSync(join(root, "docs/brand", n)).toString("base64"); // about:blank cannot load file://
const page_html = (items, bg, px) => `<body style="margin:0;background:${bg}">` + items.map(([n, w, h]) =>
  `<img src="${svg(n)}" style="display:inline-block;width:${w}px;height:${h}px;margin:10px">`).join("") + `</body>`;
const ctx = await chromium.launchPersistentContext(join(out, "browser-profile"), { viewport: { width: 760, height: 330 }, deviceScaleFactor: 1, args: ["--force-color-profile=srgb"] });
const page = ctx.pages()[0] ?? await ctx.newPage();
const shots = {
  "favicon-16-32-48-light": [[["knowslink-favicon.svg",16,16],["knowslink-favicon.svg",32,32],["knowslink-favicon.svg",48,48],["knowslink-symbol.svg",16,16],["knowslink-symbol.svg",32,32],["knowslink-symbol.svg",48,48]], "#ffffff", [400,80]],
  "favicon-16-32-48-dark": [[["knowslink-favicon.svg",16,16],["knowslink-favicon.svg",32,32],["knowslink-favicon.svg",48,48]], "#182B31", [400,80]],
  "lockup-light": [[["knowslink-lockup.svg",231,266],["knowslink-lockup-ko.svg",238,272]], "#ffffff", [520,300]],
  "lockup-inverse-dark": [[["knowslink-lockup-inverse.svg",231,266],["knowslink-symbol-inverse.svg",200,200]], "#182B31", [520,300]],
  "lockup-inverse-black": [[["knowslink-lockup-inverse.svg",231,266],["knowslink-symbol-inverse.svg",200,200]], "#000000", [520,300]],
};
for (const [name, [items, bg, [w, h]]] of Object.entries(shots)) {
  await page.emulateMedia({ colorScheme: bg === "#ffffff" ? "light" : "dark" }); // favicon.svg switches colors with the OS scheme
  await page.setViewportSize({ width: w, height: h });
  await page.setContent(page_html(items, bg));
  await page.waitForTimeout(150);
  await page.screenshot({ path: join(out, name + ".png") });
}
await ctx.close();
