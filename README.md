# memx

Local-first personal memory & knowledge store for LLM agents.

`memx` は、ローカルLLM／エージェントから使うことを前提にした  
**4層構成（short / chronicle / memopedia / archive）** のメモリ／知識ストアと CLI です。

ログをただ溜めるのではなく、

- 短期メモ（短期記憶）
- 時系列の出来事（エピソード記憶）
- 抽象化された知識（意味記憶）
- 退避された履歴（長期保管）

に分けて管理し、「蒸留（distill）」と GC（Observer / Reflector）を通じて  
ログを「参照しやすい知識」に育てていくことを目標としています。

> ※まだ **pre-alpha** 段階です。仕様が固まりつつあるところで、実装はこれから。

---

## Goals

- ローカルで完結する「個人用・長期メモリ＆知識ストア」を提供する
- LLM／エージェントから使いやすい **シンプルな CLI と API** を用意する
- 「短期メモ → 観測ノート → 知識ページ」への**一方向のフロー**で、記憶の整理を自動化する
- Semantic Recall（意味ベース検索）＋タグ＋FTS で **検索性能を重視** する
- Web UI や常時稼働エージェントではなく、**バッチ／都度呼び出し前提の設計**にする

---

## Architecture Overview

### 4 Stores

物理的に 4つの SQLite DB を持ちます（ATTACH 前提）。

- `short.db`
  - すべてのノートが最初に入る「短期メモ」ストア
  - 断片的なログ・生のメモなど
- `chronicle.db`
  - 日記・旅程・プロジェクト進捗など、「時間軸で意味を持つログ」
- `memopedia.db`
  - 用語定義・設計・ポリシーなど、「時間軸から独立した知識ベース」
- `archive.db`
  - 古い／優先度の低いノートを退避するストア
  - 通常検索からは外すが、バックトラック用に保持

### Storage / Indexing

各ストアは基本的に同じ構造を持ちます（`short.db` は superset）:

- `notes` … ノート本体
- `tags` / `note_tags` … タグとノートの多対多
- `note_embeddings` … 意味検索用のベクター（埋め込み）
- `notes_fts` … FTS5 による全文検索インデックス（archive は optional）

`short.db` にのみ、追加で:

- `short_meta` … GC 用メタ情報（note_count / token_sum / last_gc_at など）
- `lineage` … 蒸留・昇格・退避の系譜（どのノートがどこへ統合されたか）

### LLM Roles

LLM は役割ごとに分けて扱います：

- **EmbeddingClient**
  - テキスト → ベクター（埋め込み）
  - Semantic Recall で使用
- **MiniLLMClient**
  - タグ生成
  - `relevance / quality / novelty / importance` の初期スコア推定
  - 機密度（`sensitivity`）推定
  - 軽量モデル（1B〜3B）想定
- **ReflectLLMClient**
  - 観測ノートクラスタの要約（Observer）
  - Memopedia ページの更新（Reflector）
  - 7B〜30B クラス想定

これらは Go の interface（`go/llm_client.go`）として定義され、  
`db.Conn` に注入して使う構成になっています。

---

## CLI Design (draft)

コマンド名は `mem` を想定しています。

### 基本コマンド

- `mem in short`
  - 短期ストア（short.db）へのノート投入
- `mem out search`
  - FTS5 によるキーワード検索
- `mem out recall`
  - ベクター検索＋前後文脈の Semantic Recall
- `mem gc short`
  - short.db の GC（Observer / Reflector を含む）

今後追加予定のもの：

- `mem distill`
  - 手動での蒸留／統合
- `mem working`
  - Memopedia の Working Memory（常時ピン留めノート）の操作
- `mem lineage`
  - あるノートがどのログから来たかを辿る

### 入力例（mem in short）

```bash
mem in short \
  --title "Qwen3.5-27B ローカルメモ" \
  --file ./note.txt \
  --source-type web \
  --origin "https://example.com/article"
