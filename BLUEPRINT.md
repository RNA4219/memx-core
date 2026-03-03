---
intent_id: memx-blueprint
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# BLUEPRINT

## 目的
- 個人用・ローカル運用の「知識＋記憶」管理基盤を提供する。
- 短期メモ（short）/時系列ログ（chronicle）/知識ベース（memopedia）/保管庫（archive）を分離して管理する。
- CLI + API + SQLite を前提に、実行環境依存の小さい知識OS下層を構成する。

## スコープ/非ゴール
- 非ゴール: Web UI、マルチユーザー運用前提の認証/権限/監査、常駐必須設計、完全自動運転エージェント。
- v1 必須:
  - CLI: `mem in short`, `mem out search`, `mem out show`
  - API: `POST /v1/notes:ingest`, `POST /v1/notes:search`, `GET /v1/notes/{id}`
- v1.1 以降: GC/recall/working/tag/meta/lineage の拡張。

## 制約
- レイヤリング: `CLI -> API -> Service -> DB/LLM/Gatekeeper`。
- CLI は入出力整形のみを担当し、DB へ直接アクセスしない。
- API は安定 JSON I/F を提供し、CLI と 1:1 対応を維持する。
- DB は 4 分割（`short.db`, `chronicle.db`, `memopedia.db`, `archive.db`）を維持する。
- schema 移行は SQL ファイル適用 + `PRAGMA user_version` ルールに従う。

## 非機能要件
- OS: Linux/macOS/Windows 上でローカル実行。
- DB: SQLite3（WAL, foreign_keys=ON）。
- 言語/実装前提: Go 単一バイナリ、HTTP は `net/http` 想定。
- セキュリティ: 秘密情報は保存前に policy + Gatekeeper で遮断。
- 拡張性: 将来機能追加時も CLI/API/DB の後方互換を維持する。
- 一貫性方針: ATTACH 跨ぎの完全原子性は前提にせず、「データ喪失より重複許容」。

## 性能/信頼性方針
- v1 の最小性能目標として `ingest/search/show` のローカル実用応答を維持。
- 外部クライアント呼び出し契約:
  - タイムアウト 15 秒/回
  - 最大 2 回リトライ（指数バックオフ）
  - 429/502/503/504 と接続系は再試行可
  - 400/401/403/スキーマ不整合は再試行不可
- ingest 部分失敗時は `notes` 保存成功を最小コミット境界とし、Gatekeeper deny 時のみ fail-closed。
