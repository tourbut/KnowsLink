const task = await taskSpace("connect-simple-qa-20261027");
const page = task.page("p1");
const out = "C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/";
await page.goto("http://127.0.0.1:18951/connect/73AQUOOR8435riQLutueJub6RJViR5jGPnSGrztY4to");
console.log((await page.evaluate(() => document.body.innerText)).slice(0, 300).replace(/\n/g, " "));
await page.screenshot({ path: out + "shot-06-link-revisit-completed.png" });
// another browser state: sign out, same link
await page.evaluate(() => { const f = document.createElement("form"); f.method = "POST"; f.action = "/auth/logout"; document.body.appendChild(f); f.submit(); });
await page.waitForLoadState();
await page.goto("http://127.0.0.1:18951/connect/73AQUOOR8435riQLutueJub6RJViR5jGPnSGrztY4to");
console.log("logged-out:", (await page.evaluate(() => document.body.innerText)).slice(0, 300).replace(/\n/g, " "));
await page.screenshot({ path: out + "shot-07-link-revisit-logged-out.png" });
