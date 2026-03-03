---
intent_id: memx-checklists
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# CHECKLISTS

## リリース前確認項目

### 互換性
- [ ] v1 API の既存 request/response 必須項目を破壊していない。
- [ ] CLI 主要コマンド（`mem in short`, `mem out search`, `mem out show`）が後方互換で動作する。
- [ ] 未分類エラーが `INTERNAL` にフォールバックしている。

### 動作確認
- [ ] `POST /v1/notes:ingest` が成功し note を返す。
- [ ] `POST /v1/notes:search` が期待件数で返る。
- [ ] `GET /v1/notes/{id}` が単一 note を返す。
- [ ] `mem gc short --dry-run` が予定操作を表示し DB を変更しない。

### セーフティ/失敗時
- [ ] Gatekeeper `deny` / `needs_human` で fail-closed になる。
- [ ] retry/timeout（15 秒、最大 2 回）が外部クライアント呼び出しに適用される。
- [ ] ingest 部分失敗時に `notes` 保存が維持される。

### データ/スキーマ
- [ ] DB 4 分割（short/chronicle/memopedia/archive）の前提を壊していない。
- [ ] 非互換 DDL 時のみ `PRAGMA user_version` を +1 している。
- [ ] GC 閾値が `memory_policy.yaml.gc.short` を単一参照している。

### 運用
- [ ] `BLUEPRINT.md`, `RUNBOOK.md`, `GUARDRAILS.md`, `EVALUATION.md`, `CHECKLISTS.md` の front matter が更新されている。
- [ ] `last_reviewed_at`/`next_review_due` が妥当な日付で設定されている。
