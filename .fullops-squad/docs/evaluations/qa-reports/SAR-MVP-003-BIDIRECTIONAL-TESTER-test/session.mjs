// Shared stdio MCP session and safe HTTP result parsing for the independent trial observer.
import { spawn } from "node:child_process";
import { createServer } from "node:http";
import { sign } from "node:crypto";
import { once } from "node:events";

export const clone = process.env.CLONE;
export const relay = process.env.RELAY_URL || "";
export const checks = [];
export let failed = false;

export function check(name, pass, detail = "") {
  checks.push({ name, pass: Boolean(pass), detail });
  if (!pass) failed = true;
}

export function scrub(text, secrets) {
  let value = String(text ?? "");
  for (const secret of secrets) {
    if (secret && secret.length > 5) value = value.split(secret).join("[redacted]");
  }
  return value.replace(/-----BEGIN [A-Z ]+-----[\s\S]*?-----END [A-Z ]+-----/g, "[redacted-key]");
}

function attach(target, child) {
  target.stderr = "";
  target.pending = new Map();
  target.buffer = "";
  target.next = 1;
  target.child = child;
  child.stderr.on("data", (chunk) => {
    target.stderr += chunk;
  });
  child.stdout.on("data", (chunk) => {
    target.buffer += chunk;
    let index = target.buffer.indexOf("\n");
    while (index >= 0) {
      const line = target.buffer.slice(0, index);
      target.buffer = target.buffer.slice(index + 1);
      if (line.trim()) {
        const message = JSON.parse(line);
        const resolve = target.pending.get(message.id);
        if (resolve) {
          target.pending.delete(message.id);
          resolve(message);
        }
      }
      index = target.buffer.indexOf("\n");
    }
  });
}
export class Mcp {
  constructor(plugin, env) {
    attach(this, spawn(process.execPath, [plugin], { cwd: "/tmp", env, stdio: ["pipe", "pipe", "pipe"] }));
  }
  static fromChild(child) {
    const mcp = Object.create(Mcp.prototype);
    attach(mcp, child);
    return mcp;
  }
  request(method, params) {
    const id = this.next++;
    const payload = JSON.stringify({ jsonrpc: "2.0", id, method, params }) + "\n";
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => reject(new Error(`mcp timeout ${method}`)), 20000);
      this.pending.set(id, (message) => {
        clearTimeout(timer);
        resolve(message);
      });
      this.child.stdin.write(payload);
    });
  }
  notify(method) {
    this.child.stdin.write(JSON.stringify({ jsonrpc: "2.0", method }) + "\n");
  }
  async init() {
    const ready = await this.request("initialize", {
      protocolVersion: "2025-11-25",
      capabilities: {},
      clientInfo: { name: "sar-mvp-003-independent", version: "0" },
    });
    check("mcp initialize", !ready.error && ready.result?.protocolVersion === "2025-11-25", ready.result?.protocolVersion || "error");
    this.notify("notifications/initialized");
  }
  async close() {
    if (this.done) return;
    this.done = true;
    this.child.stdin.end();
    await Promise.race([once(this.child, "exit"), new Promise((resolve) => setTimeout(resolve, 2000))]);
    if (this.child.exitCode === null) this.child.kill();
  }
}

export function baseEnv(extra) {
  return { PATH: "/usr/bin:/bin", ...extra };
}

export async function call(mcp, name, args = {}) {
  const message = await mcp.request("tools/call", { name, arguments: args });
  const text = message.result?.content?.[0]?.text || "";
  return { isError: Boolean(message.result?.isError), body: text ? JSON.parse(text) : {}, rpcError: Boolean(message.error) };
}

export async function listen(handler) {
  const server = createServer(handler);
  server.listen(0, "127.0.0.1");
  await once(server, "listening");
  const address = server.address();
  return { server, url: `http://127.0.0.1:${address.port}` };
}

const envKeys = ["KNOWSLINK_MODE", "RELAY_URL", "AGENT_ID", "AGENT_KID", "AGENT_CREDENTIAL", "AGENT_KEY_FILE", "KNOWSLINK_TEST_PEER", "CF_ACCESS_CLIENT_ID", "CF_ACCESS_CLIENT_SECRET"];
const savedEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));
export function restoreEnv() {
  for (const key of envKeys) {
    if (savedEnv[key] === undefined) delete process.env[key];
    else process.env[key] = savedEnv[key];
  }
}

export function signed(own, peer, fields) {
  const { signingBytes, uuid7 } = fields.core;
  const envelope = {
    v: "relay.v1",
    id: uuid7(),
    from: own.agent,
    to: fields.to || peer,
    intent: fields.intent || "relay.test.message",
    body: fields.body || { text: fields.text },
    deliver: "agent",
    exp: new Date(Date.now() + fields.expMs).toISOString().replace(/\.\d{3}Z$/, "Z"),
    idempotency_key: fields.key,
    sig: { alg: "Ed25519", kid: "key1", value: "" },
  };
  envelope.sig.value = sign(null, signingBytes(envelope), own.pem).toString("base64url");
  if (fields.corrupt) envelope.sig.value = "A".repeat(86);
  return envelope;
}

export async function post(path, token, body, secrets) {
  const response = await fetch(`${relay}${path}`, {
    method: "POST",
    redirect: "error",
    headers: { ...(token ? { authorization: `Bearer ${token}` } : {}), "content-type": "application/json" },
    body: JSON.stringify(body),
  });
  const raw = scrub(await response.text(), secrets);
  let error = "";
  let id = "";
  try {
    const parsed = JSON.parse(raw);
    if (parsed && typeof parsed.error === "string" && /^[a-z0-9_]+$/.test(parsed.error)) error = parsed.error;
    if (parsed && typeof parsed.id === "string" && /^[0-9a-f-]{36}$/.test(parsed.id)) id = parsed.id;
  } catch {
    error = raw ? "opaque" : "";
  }
  return { status: response.status, error, id };
}
