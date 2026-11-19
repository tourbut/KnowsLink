const task = await taskSpace("connect-simple-qa-20261027");
console.log("space", task.spaceId);
const page = task.page("p1");
const out = "C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/";
const base = "http://127.0.0.1:18951";
const link = "http://127.0.0.1:18951/connect/73AQUOOR8435riQLutueJub6RJViR5jGPnSGrztY4to";
await page.goto(base + "/");
// the fixture's session endpoint stands in for a returning signed-in member D (cookie set by the relay itself)
await page.evaluate(() => { const f = document.createElement("form"); f.method = "POST"; f.action = "/fixture/as/d"; document.body.appendChild(f); f.submit(); });
await page.waitForURL("**/fixture/as/d");
console.log("fixture login ok", await page.url());
await page.goto(base + "/home/agents");
await page.waitForSelector("#unnamed");
await page.screenshot({ path: out + "shot-01-agents-unnamed-section.png" });
const unnamed = await page.evaluate(() => {
  const s = document.querySelector("#unnamed");
  return { title: s.querySelector("h2").textContent, inDetails: !!s.closest("details"), hasInput: !!s.querySelector("input[name=name]"), button: s.querySelector("button").textContent, firstSection: [...document.querySelectorAll("main section")].map((x) => x.id) };
});
console.log(JSON.stringify(unnamed));
// the app's link opened by the person
await page.goto(link);
await page.screenshot({ path: out + "shot-02-connect-link-signed-in-before-google.png" });
console.log(JSON.stringify(await page.evaluate(() => document.body.innerText.slice(0, 400))));
