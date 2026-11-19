import { readFileSync } from "node:fs";
const F = JSON.parse(readFileSync(process.argv[2], "utf8"));
const cred = (k) => JSON.parse(readFileSync(F[k].folder + "/agent.json", "utf8")).credential;
for (const [who, target] of [["b", "a"], ["a", "b"]]) {
  const r = await fetch(`${F.relay}/v1/contacts/${F[target].agent}`, { headers: { Authorization: `Bearer ${cred(who)}` } });
  console.log(`contacts ${who} -> ${target}: ${r.status} ${(await r.text()).trim()}`);
}
