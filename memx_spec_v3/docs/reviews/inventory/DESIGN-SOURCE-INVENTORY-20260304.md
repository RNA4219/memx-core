# DESIGN SOURCE INVENTORY (2026-03-04)

- 対象入力は `orchestration/memx-design-docs-authoring.md` Phase 1 Dependencies と一致させる。
- 承認条件: `memx_spec_v3/docs/design-source-inventory-operations-spec.md` に従い、`blocked=0` 行のみで構成される場合に承認可能（`blocked` 行が 1 件でもあれば Phase 1 Done 不可）。

| source_path#section | req_id | contract_ref | node_id | depends_on | owner | reviewed_at | blocked |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `memx_spec_v3/docs/interfaces.md#1.2 \`mem out search\`（IF-CLI-SEARCH-REQ/RES）` | `REQ-CLI-001` | `memx_spec_v3/docs/contracts/cli-json.schema.json#/mem.out.search` | `api` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/interfaces.md#2.1 \`POST /v1/notes:ingest\`（IF-API-INGEST-REQ/RES）` | `REQ-API-001` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1notes:ingest/post` | `api` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/requirements.md#2. SHOULD（v1.x）` | `REQ-GC-001` | `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1v1~1gc:run/post` | `requirements` | `api` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/traceability.md#1. 主要 REQ-ID トレーサビリティ表` | `REQ-SEC-001` | `traceability: REQ-SEC-001 mapping` | `requirements` | `design` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/design.md#2.1 store別設計詳細（short/journal/knowledge/archive）` | `REQ-RET-001` | `design section: retention/archive behavior` | `design` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/design.md#2. DB 責務分割` | `REQ-SEC-AUD-001` | `design section: audit responsibility split` | `design` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `memx_spec_v3/docs/interfaces.md#4.1 ErrorCode × HTTP × retryable × クライアント動作（IF-ERR-MATRIX）` | `REQ-SEC-AUD-002` | `interfaces error/audit linkage` | `api` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `docs/birdseye/index.json#nodes[guardrails]` | `REQ-SEC-GRD-001` | `docs/birdseye/index.json#/nodes[node_id=guardrails]` | `guardrails` | `requirements` | `memx-core` | `2026-03-04` | `0` |
| `docs/birdseye/caps/RUNBOOK.md.json#trace-req-err-001` | `REQ-ERR-001` | `runbook trace id: trace-req-err-001` | `runbook` | `api` | `memx-core` | `2026-03-04` | `0` |
| `docs/birdseye/caps/EVALUATION.md.json#req-nfr-001-passfail` | `REQ-NFR-001` | `evaluation gate: req-nfr-001-passfail` | `evaluation` | `runbook` | `memx-core` | `2026-03-04` | `0` |
