---
intent_id: memx-runbook
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# RUNBOOK

## 1. API サーバー起動
```bash
mem api serve
```
- 既定のローカル API（例: `http://127.0.0.1:7766`）を起動する。

## 2. ノート投入（ingest）
```bash
mem in short \
  --title "Qwen3.5-27B ローカルメモ" \
  --file ./note.txt \
  --source-type web \
  --origin "https://example.com/article"
```
- `--no-llm`: LLM/Embedding を使わず最小メタで保存する。
- `--tags`: 手動タグをカンマ区切りで付与する。

## 3. キーワード検索（FTS）
```bash
mem out search "Qwen3.5 ベンチ" --store short --limit 10
```
- FTS5 ベースの検索。対象ストア/件数を指定可能。

## 4. 単一ノート表示
```bash
mem out show <note-id>
```
- note id を指定して単一ノートを取得する。

## 5. セマンティック検索（recall）
```bash
mem out recall "Qwen3.5-27B ベンチマーク結果" \
  --scope self \
  --stores short,chronicle,memopedia \
  --top-k 8 \
  --range 3
```
- `top-k` は `1..50`（未指定 8、超過時 50 に丸め）。
- `range` は `0..10`（未指定 3）。
- `Conn.Embed == nil` では既定エラー。明示フラグ時のみ FTS フォールバック。

## 6. GC 実行
```bash
mem gc short
mem gc short --dry-run
```
- `--dry-run` は DB 変更なしで予定操作のみ表示。
- 閾値は `memory_policy.yaml.gc.short` を単一の参照元として使用する。
