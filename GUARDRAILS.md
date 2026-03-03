---
intent_id: memx-guardrails
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# GUARDRAILS

## 絶対遵守ルール

### 1) 破壊的変更の禁止
- v1 系（`/v1/*`）の既存クライアントを壊す変更を禁止する。
- 既存必須フィールドの削除/意味変更を禁止する。
- DB スキーマの非互換変更は `user_version` の +1 と移行手順を必須化する。

### 2) 後方互換の維持
- v1 内で許可されるのは任意フィールド追加などの後方互換変更のみ。
- CLI オプションは API request フィールドへの 1:1 マッピングを維持する。
- エラー互換は `INVALID_ARGUMENT` / `NOT_FOUND` / `INTERNAL` 最小集合を保証し、未分類は `INTERNAL` にフォールバックする。

### 3) 失敗時挙動（fail-safe/fail-closed）
- Gatekeeper が `deny` または `needs_human` を返した場合は fail-closed で処理中断する。
- ingest の部分失敗は `notes` 保存成功を優先し、タグ/埋め込み失敗は後追い再実行対象として扱う。
- クライアント呼び出しは 15 秒 timeout、再試行規約（最大 2 回）に従う。

### 4) 変更適用ルール
- レイヤ境界（CLI→API→Service→Infra）を越えた直接依存を追加しない。
- 設定の単一情報源（例: GC 閾値は `memory_policy.yaml.gc.short`）を崩さない。
- 秘密情報は保存前に policy + Gatekeeper でブロックする。
