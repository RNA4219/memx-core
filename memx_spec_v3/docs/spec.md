---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 仕様（spec）

## 1. 対象ユースケース
- ローカル環境での個人メモ投入・検索・参照。
- LLM/Agent からの機械呼び出し（CLI/API）での短期記憶運用。
- 将来の蒸留・昇格・GC を見据えた 4 ストア構成（short/chronicle/memopedia/archive）。

## 2. スコープ境界
### In Scope（v1）
- CLI: `mem in short` / `mem out search` / `mem out show`。
- API: `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}`。
- ローカル SQLite を前提にした単体運用。
- CLI `--json` と API レスポンスの同型維持。

### Out of Scope（v1）
- Web UI。
- マルチユーザー運用、認証・認可・監査基盤の本格提供。
- 常駐必須プロセス設計。
- 完全自律エージェントランタイム。

## 3. 非ゴール
- GUI での操作体験最適化。
- クラウド前提の水平分散や外部ベクターDB必須化。
- v1 内での破壊的 API/CLI 変更。

## 4. 受け入れ観点
- 互換性: v1 必須 I/F の後方互換を維持する。
- エラー: 入力不備は 400 系、内部障害は 500 系で返す。
- 品質: ingest/search/show がローカル単体で実用応答時間を満たす。
- 安全性: fail-closed 方針に従い、機密入力は保存拒否できる。
