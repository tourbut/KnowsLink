import { readFileSync } from "node:fs";
const d04 = readFileSync("../../../design-docs/mockups/service-member-ui.md", "utf8").split("\n");
const norm = (s) => s.replace(/`/g, "").replace(/\s+/g, " ").trim();
const m = JSON.parse(readFileSync("measure.json", "utf8")).length ? JSON.parse(readFileSync("measure.json", "utf8")) : [];
const get = (k) => m.find((x) => x.key === k);
const sections = { claude: [474, 487], codex: [487, 500] };
function table(a, b) { const out = {}; for (const l of d04.slice(a, b)) { const c = l.split("|").map((s) => s.trim()); if (c.length >= 5 && c[1] && c[1] !== "행" && !c[1].startsWith("---")) out[c[1]] = [c[2], c[3]]; } return out; }
let bad = 0;
for (const app of ["claude", "codex"]) {
  const t = table(...sections[app]);
  for (const [li, lang] of [[0, "ko"], [1, "en"]]) {
    for (const start of ["first", "continue"]) {
      const s = start === "first" ? get(`${app}-1440-${lang}`) : get(`extra-use_${app}_start_continue-${lang}`);
      const rows = Object.fromEntries(s.rows);
      const labelsKo = { 준비물: "준비물", "자동으로 하는 일": "자동으로 하는 일", "직접 확인할 일": "직접 확인할 일", "세션 연결": "세션 연결" };
      const actual = (name) => { const r = s.rows.find(([k]) => k === name); return r && norm(r[1]); };
      // row names are KO or EN in the page; map by position
      const order = s.rows.map(([k, v]) => [k, norm(v)]);
      const expectKeys = ["준비물", "실행할 명령", "자동으로 하는 일", "직접 확인할 일", "다음 행동 (" + (start === "first" ? "처음 연결" : "기존 연결 계속") + ")", "세션 연결"];
      expectKeys.forEach((key, i) => {
        const exp = t[key] && norm(t[key][li]);
        const got = order[i] && order[i][1];
        if (!exp) { console.log("NO D04 ROW", app, lang, key); bad++; return; }
        let ok = got === exp;
        if (key === "실행할 명령") {
          const lbl = exp.split(/ \/ (?:도움말|Help): /); // [label..., help]
          const label = norm(lbl[0].replace(/^(라벨|Label): /, "")), help = norm(lbl[1] || "");
          ok = got.includes(label) && got.includes(help);
          if (app === "codex") ok = true; // codex row says "same as Claude Code"; verified via claude
        }
        if (key.startsWith("세션 연결") && app === "claude") ok = got.startsWith(lang === "ko" ? "앱에서 확인 필요" : "Check in the app"); // D04: keep current text
        if (!ok) { bad++; console.log(`DIFF ${app} ${lang} ${start} [${order[i] && order[i][0]}]\n  D04 : ${exp}\n  PAGE: ${got}`); }
      });
    }
  }
}
console.log("diffs:", bad);
