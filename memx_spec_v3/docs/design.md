---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 設計（design）

## 1. レイヤ構成
```
CLI -> API -> Service(Usecase) -> DB / LLM / Gatekeeper
```
- CLI: 入力整形と表示のみを担当。
- API: 安定 JSON I/F を提供。
- Service: ビジネスロジックの唯一入口。
- DB/LLM/Gatekeeper: 副作用を持つインフラ層。

## 2. DB 責務分割
- `short.db`: 一次投入先。短期メモ、GC 対象の起点。
- `chronicle.db`: 時系列ログ（出来事・進捗）。
- `memopedia.db`: 抽象知識（定義・方針）。
- `archive.db`: 退避保管（通常検索対象外）。

共通責務:
- `notes`, `tags`, `note_tags`, `note_embeddings`, `notes_fts`（archive は一部省略可）。

short 固有:
- `short_meta`: GC 判定メタ。
- `lineage`: 蒸留/昇格/退避の系譜。

## 3. 移行戦略
- マイグレーションは `schema/*.sql` を正本として適用する。
- `PRAGMA user_version` を採用し、破壊的/非互換 DDL のみバージョンを進める。
- v1 では後方互換を最優先し、破壊変更は v2+（`FUTURE`）へ隔離する。
- 実験機能は feature flag 既定 OFF で段階導入する。
