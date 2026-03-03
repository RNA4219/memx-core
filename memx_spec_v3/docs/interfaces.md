---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx インターフェース仕様（interfaces）

## 1. CLI I/O（v1 必須）
- `mem in short`
  - Input: title/body/source 等。
  - Output: 生成 note ID と store、作成時刻。
- `mem out search`
  - Input: query/store/limit 等。
  - Output: 検索結果 items + total。
- `mem out show`
  - Input: note ID。
  - Output: ノート詳細（id/store/title/body/created_at）。

### JSON 出力規則
- `--json` は API レスポンスと同型（同一キー体系・同一意味）を維持する。

## 2. API I/O（v1 必須）
- `POST /v1/notes:ingest`
  - request: `store`, `title`, `body`。
  - response: `id`, `store`, `created_at`。
- `POST /v1/notes:search`
  - request: `store`, `query`, `limit`。
  - response: `items[]`, `total`。
- `GET /v1/notes/{id}`
  - response: `id`, `store`, `title`, `body`, `created_at`。

## 3. 互換ルール
- 必須フィールド削除禁止。
- 既存フィールド意味変更禁止。
- 成功レスポンストップレベル構造変更禁止。
- 破壊変更は v2+ で段階移行（互換フラグまたは新バージョン導入）。

## 4. エラー面
- 入力不正: HTTP 400 + `INVALID_ARGUMENT` 系。
- ポリシー拒否: HTTP 403 + `POLICY_DENIED`（fail-closed）。
- 内部障害: HTTP 500 + `INTERNAL`。
- `retryable` は error code と整合させる。
