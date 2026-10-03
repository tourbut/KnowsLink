// Development adapter entry point; no polling, vendor API, or tool execution.
export const adapterStatus = {
  state: "unimplemented",
  transport: "pull",
  webhook: false,
  evidenceFetch: false,
} as const;

console.info(JSON.stringify(adapterStatus));
