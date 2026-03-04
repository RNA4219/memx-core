# memx

Local-first personal memory & knowledge store for LLM agents.

`memx` は、ローカルLLM／エージェント向けの 4 ストア構成（short / chronicle / memopedia / archive）のメモリ基盤です。

## 設計書作成開始入口（IA）
- 参照入口（責務境界・更新優先順位・Trigger別必須更新先）: [memx_spec_v3/docs/design-doc-ia-spec.md](./memx_spec_v3/docs/design-doc-ia-spec.md)

## 正本ドキュメント（memx_spec_v3/docs）
- 参照起点・役割分担（正本/補助）: [memx_spec_v3/docs/spec.md](./memx_spec_v3/docs/spec.md)
- 要求事項の正本（MUST/SHOULD/FUTURE, ID 定義）: [memx_spec_v3/docs/requirements.md](./memx_spec_v3/docs/requirements.md)
- API 契約の正本: [memx_spec_v3/docs/contracts/openapi.yaml](./memx_spec_v3/docs/contracts/openapi.yaml)
- CLI `--json` 契約の正本: [memx_spec_v3/docs/contracts/cli-json.schema.json](./memx_spec_v3/docs/contracts/cli-json.schema.json)

## 補助・参照ドキュメント（memx_spec_v3/docs）
- 設計詳細: [memx_spec_v3/docs/design.md](./memx_spec_v3/docs/design.md)
- I/F 詳細: [memx_spec_v3/docs/interfaces.md](./memx_spec_v3/docs/interfaces.md)
- 契約一覧（参照）: [memx_spec_v3/docs/CONTRACTS.md](./memx_spec_v3/docs/CONTRACTS.md)

> README では重複記述を避け、詳細仕様は `spec.md` の役割分担に従って参照してください。

- 品質ゲート参照先と適用範囲: memx 本体（`memx_spec_v3/`）は [`docs/QUALITY_GATES.md`](./docs/QUALITY_GATES.md)、`workflow-cookbook/` は [`workflow-cookbook/docs/QUALITY_GATES.md`](./workflow-cookbook/docs/QUALITY_GATES.md) を参照。

## Governance Docs
- [BLUEPRINT.md](./BLUEPRINT.md)
- [RUNBOOK.md](./RUNBOOK.md)
- [GUARDRAILS.md](./GUARDRAILS.md)
- [EVALUATION.md](./EVALUATION.md)
- [CHECKLISTS.md](./CHECKLISTS.md)

## License
Apache-2.0
