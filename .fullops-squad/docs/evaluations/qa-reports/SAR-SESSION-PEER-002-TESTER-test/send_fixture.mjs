// Fixture sender (no host approval gate): signs and sends one text as the synthetic agent in <folder>. Receiver-side QA only.
// usage: node send_fixture.mjs <dist-dir> <folder> <peer> <idempotency-key> <text-file>
import { pathToFileURL } from "node:url";
import { join } from "node:path";
import { readFileSync } from "node:fs";
const [dist, folder, peer, key, textFile] = process.argv.slice(2);
const { memberTransport } = await import(pathToFileURL(join(dist, "text.js")).href);
const t = await memberTransport(folder);
const text = readFileSync(textFile, "utf8");
const r = await t.sendText(peer, text, key, true);
console.log(JSON.stringify({ ...r, bytes: Buffer.byteLength(text) }));
