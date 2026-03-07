---
owner: memx-core
status: active
last_reviewed_at: 2026-03-06
next_review_due: 2026-06-06
priority: high
---

# memx 要求事項（MUST / SHOULD / FUTURE）

> **注意**: 本書は要件の索引・概要を記載します。詳細は分割ファイルを参照してください。

## 仕様文書の正本範囲

- 本書（`requirements.md`）は memx v1.3 の要件正本の索引です。
- API 契約の正本は `contracts/openapi.yaml`、CLI `--json` 契約の正本は `contracts/cli-json.schema.json` とする。
- `error-contract.md` と `quickstart.md` は本書に従属する補助文書です。
- 参照起点は `spec.md` とします。

## 0. 目的とスコープ

設計詳細は `design.md`、I/F 詳細は `interfaces.md` を参照し、契約詳細は正本の `contracts/openapi.yaml` と `contracts/cli-json.schema.json` を参照してください。

## ADR参照運用ルール

- 要件に影響する設計判断を追加/変更する場合、該当節へ ADR リンクを追記してください。
- `design.md` と本書の該当節リンクは同一PRで同期更新してください。
- ADR は `docs/ADR/README.md` 索引更新を必須とし、未更新のままマージしません。

## 1. MUST（v1）

- CLI: `mem in short` / `mem out search` / `mem out show` を提供する。
- API: `POST /v1/notes:ingest` / `POST /v1/notes:search` / `GET /v1/notes/{id}` を提供する。
- CLI `--json` は API レスポンスと同型を維持する。
- 入力不備は 400 系、内部障害は 500 系で返す。
- fail-closed 方針で機密入力を拒否できること。
- ingest/search/show の最小性能目標を満たすこと。

## 2. SHOULD（v1.x）

- `POST /v1/gc:run` は「公開可否」と「実行可否」を分離して扱う。
- `mem.features.gc_short=true` 時のみ `mem gc short` / `POST /v1/gc:run` の実行を有効化する。
- SHOULD 機能は feature flag 既定 OFF で提供し、既定挙動を壊さない。

## 3. FUTURE（v2+）

- Recall/Working/Tag/Meta/Lineage/Distill 系 CLI/API の正式導入。
- 破壊変更を伴う再設計（段階移行前提）。

## 4. ID 定義

主要要件ID:

| 領域 | ID | 概要 |
| --- | --- | --- |
| CLI | `REQ-CLI-001` | v1必須3コマンド |
| API | `REQ-API-001` | v1必須3エンドポイント |
| GC | `REQ-GC-001` | GC dry-run |
| Security | `REQ-SEC-001` | fail-closed |
| Error | `REQ-ERR-001` | エラーモデル |
| NFR | `REQ-NFR-001`〜`REQ-NFR-006` | 非機能要件 |
| typed_ref | `FR-008` | typed_ref 正規化（4セグメント canonical） |
| typed_ref | `AC-006` | typed_ref 一貫性（cross-system 追跡可能） |
| bundle | `FR-006` | 継続用 bundle 保存（監査項目必須） |
| state | `FR-007` | 状態遷移明示化（current/history 分離） |
| context | `FR-003` | Context Rebuild の外部依存制約（tracker optional） |
| conflict | `FR-009` | 競合検出（stale state/bundle 検出） |
| tracker | `FR-010` | Tracker 統合（外部 issue 接続） |
| audit | `NFR-001` | 再現性（bundle 再現可能） |
| audit | `NFR-002` | 可監査性（source_refs/generator_version 追跡） |
| audit | `NFR-004` | 劣化耐性（競合検出優先） |
| audit | `AC-007` | Bundle 監査性（必須監査項目確認） |
| integration | `AC-008` | 段階的導入整合（Phase 1-3） |

---

## 分割ファイル一覧

| ファイル | 内容 | 元セクション |
| --- | --- | --- |
| [requirements/release.md](./requirements/requirements-release.md) | Release Scope Matrix / バージョニング / トレーサビリティ | 0-1, 0-2, 0-3 |
| [requirements/architecture.md](./requirements/requirements-architecture.md) | 全体アーキテクチャ / ストア構成 / Gatekeeper | 1-0 〜 1-5 |
| [requirements/data-model.md](./requirements/requirements-data-model.md) | データモデル / Security & Retention | 2-1 〜 2-7 |
| [requirements/cli.md](./requirements/requirements-cli.md) | CLI 要件 | 3-1 〜 3-6 |
| [requirements/llm.md](./requirements/requirements-llm.md) | LLM 戦略 | 4-1 〜 4-4 |
| [requirements/nfr.md](./requirements/requirements-nfr.md) | 非機能要件 | 5-1 〜 5-4 |
| [requirements/api.md](./requirements/requirements-api.md) | API 要件 | 6-1 〜 6-5 |
| [requirements/incident.md](./requirements/requirements-incident.md) | インシデント対応要件 | 11 |

---

## クイックリファレンス

### 主要要件ID一覧

| ID | 説明 | 詳細 |
| --- | --- | --- |
| `REQ-CLI-001` | CLI v1必須3コマンド | [requirements/cli.md](./requirements/requirements-cli.md) |
| `REQ-API-001` | API v1必須3エンドポイント | [requirements/api.md](./requirements/requirements-api.md) |
| `REQ-GC-001` | GC dry-run | [requirements/cli.md#3-5](./requirements/requirements-cli.md#3-5-mem-gc-shortobserver--reflector) |
| `REQ-SEC-001` | fail-closed security | [requirements/data-model.md#2-7](./requirements/requirements-data-model.md#2-7-security--retention-requirements) |
| `REQ-ERR-001` | エラーモデル | [requirements/api.md#6-4](./requirements/requirements-api.md#6-4-エラーモデル) |
| `REQ-NFR-001` | 性能目標 | [requirements/nfr.md#5-1](./requirements/requirements-nfr.md#5-1-性能目標v1必須3エンドポイント) |
| `FR-003` | Context Rebuild 外部依存制約 | [requirements/api.md#6-9](./requirements/requirements-api.md#6-9-context-rebuild-の外部依存制約fr-003) |
| `FR-006` | 継続用 bundle 保存 | [requirements/api.md#6-7](./requirements/requirements-api.md#6-7-継続用-bundle-保存fr-006) |
| `FR-007` | 状態遷移明示化 | [requirements/api.md#6-8](./requirements/requirements-api.md#6-8-状態遷移明示化fr-007) |
| `FR-008` | typed_ref 正規化 | [requirements/api.md#6-6](./requirements/requirements-api.md#6-6-typed_ref-正規化fr-008) |
| `FR-009` | 競合検出 | [requirements/api.md#6-10](./requirements/requirements-api.md#6-10-競合検出fr-009) |
| `NFR-001` | 再現性 | [requirements/nfr.md#5-6](./requirements/requirements-nfr.md#5-6-再現性可監査性nfr-001--nfr-002) |
| `NFR-002` | 可監査性 | [requirements/nfr.md#5-6](./requirements/requirements-nfr.md#5-6-再現性可監査性nfr-001--nfr-002) |
| `NFR-004` | 劣化耐性 | [requirements/nfr.md#5-8](./requirements/requirements-nfr.md#5-8-劣化耐性nfr-004-追記) |
| `AC-006` | typed_ref 一貫性 | [requirements/nfr.md#5-5](./requirements/requirements-nfr.md#5-5-typed_ref-一貫性ac-006) |
| `AC-007` | Bundle 監査性 | [requirements/nfr.md#5-7](./requirements/requirements-nfr.md#5-7-bundle-監査性ac-007) |
| `AC-008` | 段階的導入整合 | [requirements/nfr.md#5-9](./requirements/requirements-nfr.md#5-9-段階的導入整合ac-008) |