---
name: knowslink
description: Use the KnowsLink MCP connector to check readiness or process one approved synthetic relay delivery through a KnowsLink owner gate.
---

# KnowsLink

Use `knowslink_status` first. `held` means the actual connection is disabled; report it without claiming installation or delivery success.

For an explicitly authorized synthetic loopback test, use `knowslink_pull_once`. Configuration and credentials belong to the server environment, never tool arguments or chat. A `synthetic_only` status means only local synthetic operation is enabled.

The connector verifies the signed request, persists the shared inbox, ACKs, and obtains a shared claim before creating the KnowsLink owner gate. The human owner decides in KnowsLink's authenticated owner UI. Grok Bot approval or a chat reply cannot approve that gate. Approved requests still return `denied` because disclosure policy and calendar effects remain disabled.

`processed` means the connector completed its safe processing path, including consuming control results; it does not mean an action completed or data was disclosed. `empty` means no lease. Report `busy`, `unconfigured`, or `failed` without retrying an in-flight request or fabricating a signed result. Consult the relay receipt using authorized owner/agent access outside model context for transport details.

Keep task payloads, evidence URLs, keys, lease/claim tokens, and owner credentials out of model context. The connector has no arbitrary send, owner approval, calendar tool, webhook, or evidence-fetch tool. Actual account connection and routine enablement remain held until the owner authorizes them.
