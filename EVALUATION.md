---
intent_id: memx-evaluation
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# EVALUATION

## Done 定義（受け入れ条件）
- CLI→API の入出力マッピング互換が維持されていること。
- API エラーが方針どおりに返ること（入力不備は 400 系、内部障害は 500 系）。
- v1 必須コマンド/API が動作すること。
  - CLI: `mem in short`, `mem out search`, `mem out show`
  - API: `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}`

## 測定指標
- 性能:
  - `ingest` / `search` / `show` がローカル単体で実用応答時間を満たす。
- 品質:
  - エラー分類が `INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL` へ正規化されている。
  - Gatekeeper deny/needs_human が fail-closed で停止する。
- 互換性:
  - v1 API において後方互換が維持される。
  - 変更が任意フィールド追加中心である。

## 補助評価（運用）
- GC 閾値判定が `memory_policy.yaml.gc.short` のみを参照している。
- recall の入力正規化（stores/top-k/range）が規定範囲で検証される。
