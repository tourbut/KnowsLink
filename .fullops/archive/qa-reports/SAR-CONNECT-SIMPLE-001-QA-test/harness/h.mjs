// Drives product modules (dist/login.js) directly: `start <base> <folder>` begins one attempt and returns; `poll <folder>` runs finishLogin once.
const [, , cmd, a, b] = process.argv;
const dist = "file:///C:/Users/shin/orca/workspaces/KnowsLink/tester/adapters/dist/";
await import(dist + "defaults.js");
const login = await import(dist + "login.js");
try {
  if (cmd === "start") console.log(JSON.stringify(await login.beginLogin(a, b)));
  else if (cmd === "poll") console.log(JSON.stringify({ state: await login.finishLogin(a) }));
} catch (e) { console.log(JSON.stringify({ error: String(e.message ?? e) })); process.exitCode = 1; }
