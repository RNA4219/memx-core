---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
priority: high
---

# memx 要求事項（MUST / SHOULD / FUTURE）

設計詳細は `design.md`、I/F 詳細は `interfaces.md`、契約詳細は `CONTRACTS.md` を正本とする。

## 1. MUST（v1）
- CLI: `mem in short` / `mem out search` / `mem out show` を提供する。
- API: `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}` を提供する。
- CLI `--json` は API レスポンスと同型を維持する。
- 入力不備は 400 系、内部障害は 500 系で返す。
- fail-closed 方針で機密入力を拒否できること。
- ingest/search/show の最小性能目標を満たすこと。

## 2. SHOULD（v1.x）
- `mem.features.gc_short=true` 時のみ `mem gc short` / `POST /v1/gc:run` を有効化する。
- SHOULD 機能は feature flag 既定 OFF で提供し、既定挙動を壊さない。

## 3. FUTURE（v2+）
- Recall/Working/Tag/Meta/Lineage/Distill 系 CLI/API の正式導入。
- 破壊変更を伴う再設計（段階移行前提）。

## 4. ID 定義

| ID | 区分 | 定義 |
| --- | --- | --- |
| REQ-CLI-001 | CLI | v1 必須 3 コマンドの I/O 互換を維持する。 |
| REQ-API-001 | API | v1 必須 3 エンドポイントの I/O 互換を維持する。 |
| REQ-GC-001 | GC | GC dry-run は DB 非更新で判定結果のみ返す。 |
| REQ-SEC-001 | Security | 機密入力を fail-closed で拒否する。 |
| REQ-RET-001 | Retention | archive 退避/削除の保持要件を満たす。 |
| REQ-SEC-AUD-001 | Security/Audit | `archive_move` の監査固定項目を満たす。 |
| REQ-SEC-AUD-002 | Security/Audit | `archive_purge` の監査固定項目を満たす。 |
| REQ-SEC-GRD-001 | Security/Guardrails | fail-closed 整合チェックを満たす。 |
| REQ-SEC-GRD-001-1 | Security/Guardrails | 判定規則と GUARDRAILS fail-closed 方針が一致する。 |
| REQ-SEC-GRD-001-2 | Security/Guardrails | operation 責任分界（store/output/archive_move/archive_purge）が整合する。 |
| REQ-SEC-GRD-001-3 | Security/Guardrails | 差分時は同一 PR で関連文書を追従更新する。 |
| REQ-ERR-001 | Error | ErrorCode 契約と retryable 規則を維持する。 |
| REQ-NFR-001 | NFR/Performance | ingest/search/show の性能閾値を満たす。 |
| REQ-NFR-002 | NFR/Recovery | `RTO<=30分` かつ `RPO<=5分` を満たす。 |
| REQ-NFR-003 | NFR/Recovery | 検知から暫定復旧まで `15分` 以内とする。 |
| REQ-NFR-004 | NFR/Reliability | 自動再試行は最大 2 回までとする。 |
| REQ-NFR-005 | NFR/Consistency | 障害検知から `30分以内` に整合回復判定へ到達する。 |
| REQ-NFR-006 | NFR/Observability | 復旧・補償フローの収束判定を監査可能にする。 |
