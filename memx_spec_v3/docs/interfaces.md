---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx インターフェース仕様（interfaces）

## 0. 文書の位置づけ
- 本書は人間可読の I/O 説明と互換方針を扱う。
- 契約正本は `docs/contracts/openapi.yaml` と `docs/contracts/cli-json.schema.json`。
- `CONTRACTS.md` は正本スキーマへの索引・抜粋のみを扱う。

## 1. CLI I/O（v1 必須）

### 1.1 `mem in short`（IF-CLI-INGEST-REQ/RES）

#### Input（IF-CLI-INGEST-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `title` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `body` (`--stdin`/`--file`) | string | yes | none | non-empty UTF-8 text | 必須削除禁止・意味変更禁止 |
| `summary` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `source_type` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `origin` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `source_trust` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `sensitivity` | string | no | `""` 相当（未指定可） | string | 追加のみ許可（既存意味維持） |
| `tags` | array<string> | no | `[]` 相当（未指定可） | 各要素 string | 追加のみ許可（既存意味維持） |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-INGEST-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `note` | object (`Note`) | yes | none | `Note` スキーマ準拠 | トップレベル構造変更禁止 |
| `note.id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `note.title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `note.access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `note.source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `note.sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

### 1.2 `mem out search`（IF-CLI-SEARCH-REQ/RES）

#### Input（IF-CLI-SEARCH-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `query` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `top_k` | integer | no | implementation default | integer | 追加のみ許可（既存意味維持） |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-SEARCH-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `notes` | array<`Note`> | yes | none | 各要素 `Note` 準拠 | トップレベル構造変更禁止 |

### 1.3 `mem out show`（IF-CLI-SHOW-REQ/RES）

#### Input（IF-CLI-SHOW-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `--json` | boolean(flag) | no | false | flag | true 時は API response と同型を維持 |

#### Output（IF-CLI-SHOW-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

## 2. API I/O（v1 必須）

### 2.1 `POST /v1/notes:ingest`（IF-API-INGEST-REQ/RES）

#### Request（IF-API-INGEST-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `title` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `summary` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `source_type` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `origin` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `source_trust` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `sensitivity` | string | no | omitted | string | 追加のみ許可（既存意味維持） |
| `tags` | array<string> | no | omitted | 各要素 string | 追加のみ許可（既存意味維持） |

#### Response 200（IF-API-INGEST-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `note` | object (`Note`) | yes | none | `Note` 準拠 | トップレベル構造変更禁止 |

### 2.2 `POST /v1/notes:search`（IF-API-SEARCH-REQ/RES）

#### Request（IF-API-SEARCH-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `query` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `top_k` | integer | no | omitted | integer | 追加のみ許可（既存意味維持） |

#### Response 200（IF-API-SEARCH-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `notes` | array<`Note`> | yes | none | 各要素 `Note` 準拠 | トップレベル構造変更禁止 |

### 2.3 `GET /v1/notes/{id}`（IF-API-GET-REQ/RES）

#### Request（IF-API-GET-REQ）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` (path) | string | yes | none | `minLength: 1` | 必須削除禁止・意味変更禁止 |

#### Response 200（IF-API-GET-RES）
| name | type | required | default | validation | compatibility rule |
| --- | --- | --- | --- | --- | --- |
| `id` | string | yes | none | non-empty | 必須削除禁止・意味変更禁止 |
| `title` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `summary` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `body` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `created_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `updated_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `last_accessed_at` | string | yes | none | timestamp string | 必須削除禁止・意味変更禁止 |
| `access_count` | integer | yes | none | integer | 必須削除禁止・意味変更禁止 |
| `source_type` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `origin` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `source_trust` | string | yes | none | string | 必須削除禁止・意味変更禁止 |
| `sensitivity` | string | yes | none | string | 必須削除禁止・意味変更禁止 |

## 3. 互換ルール
- 必須フィールド削除禁止。
- 既存フィールド意味変更禁止。
- 成功レスポンストップレベル構造変更禁止。
- 破壊変更は v2+ で段階移行（互換フラグまたは新バージョン導入）。

## 4. エラー面

### 4.1 ErrorCode × HTTP × retryable × クライアント動作（IF-ERR-MATRIX）
| ErrorCode | HTTP | retryable | クライアント動作 |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | 400 | false | 入力修正後に再実行（自動再試行しない） |
| `NOT_FOUND` | 404 | false | 対象 ID/条件を見直し、必要なら ingest 後に再実行 |
| `CONFLICT` | 409 | false | 競合解消（重複・状態整合）後に再実行 |
| `GATEKEEP_DENY` | 403 | false | ポリシー変更または権限見直しまで停止（実装済み: `service.ErrPolicyDenied`) |
| `NEEDS_HUMAN` | 403 | false | 人間承認フロー完了まで停止（実装済み: `service.ErrNeedsHuman`） |
| `FEATURE_DISABLED` | 409 | false | feature flag 有効化まで停止 |
| `INTERNAL` | 500 | conditional | 一時障害判定時のみ指数バックオフ再試行（最大2回） |

**実装状況（2026-03-06 更新）**:
- `service.ErrPolicyDenied`: Gatekeeper が `deny` 判定時に返却
- `service.ErrNeedsHuman`: Gatekeeper が `needs_human` 判定時に返却（v1.3 ではエラー扱い）
- `api/errors.go`: `mapError` で `GATEKEEP_DENY` にマッピング

### 4.2 整合元
- エラーコード区分・運用差分は `error-contract.md` を正本運用要約として参照。
- API 契約正本は `docs/contracts/openapi.yaml`。


### 4.3 GC route 無効時契約（`POST /v1/gc:run`）
- route 非公開時は `NOT_FOUND` / `404`。
- route 公開かつ `mem.features.gc_short=false` 時は `INTERNAL` / `500`（v1 正本契約の固定値）。
- `FEATURE_DISABLED` / `409` は将来移行候補であり、v1 では採用しない。

## 5. 契約変更手順（更新順序固定）
1. `memx_spec_v3/docs/requirements.md` を更新する。
2. 正本スキーマ（`memx_spec_v3/docs/contracts/openapi.yaml` / `memx_spec_v3/docs/contracts/cli-json.schema.json`）を更新する。
3. `memx_spec_v3/docs/interfaces.md` と `memx_spec_v3/docs/CONTRACTS.md` を更新する。
4. `EVALUATION*` / `memx_spec_v3/docs/operations-spec.md`（RUNBOOK 相当）を更新する。
5. 契約差分チェック手順（`operations-spec.md`）を実施し、レビュー可能状態にする。


## 6. 付録: RUNBOOK連携 I/F ID（v1運用）
| I/F 項目ID | 対象 |
| --- | --- |
| `IF-GC-SHORT-REQ` | `mem gc short` の入力（`--dry-run` 等の実行パラメータ） |
| `IF-GC-SHORT-RES` | `mem gc short` の出力（dry-run 結果 JSON / 実行結果） |

## 7. typed_ref I/F 仕様（FR-008 / AC-006）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-1-typed_ref-canonical-format-の固定`

### 7.1 canonical format

```txt
<domain>:<entity_type>:<provider>:<entity_id>
```

| セグメント | 説明 | 許容値 |
|-----------|------|--------|
| `domain` | システム領域 | `memx`, `agent-taskstate`, `tracker` |
| `entity_type` | エンティティ種別 | `evidence`, `artifact`, `knowledge`, `lineage`, `task`, `decision`, `issue` 等 |
| `provider` | データソース | `local`, `jira`, `github`, `linear` 等 |
| `entity_id` | 一意識別子 | システム固有のID形式 |

### 7.2 memx-core での使用

| I/F 項目ID | フォーマット | 例 |
|-----------|-------------|-----|
| `IF-TYPEDREF-EVIDENCE` | `memx:evidence:local:<id>` | `memx:evidence:local:01HXXX` |
| `IF-TYPEDREF-ARTIFACT` | `memx:artifact:local:<id>` | `memx:artifact:local:01HYYY` |
| `IF-TYPEDREF-KNOWLEDGE` | `memx:knowledge:local:<id>` | `memx:knowledge:local:01HZZZ` |
| `IF-TYPEDREF-LINEAGE` | `memx:lineage:local:<id>` | `memx:lineage:local:01HAAA` |

### 7.3 移行期の互換（read-both / write-one）

#### 入力（パーサー）

| 入力形式 | 正規化結果 | 備考 |
|---------|-----------|------|
| `memx:evidence:01HXXX` | `memx:evidence:local:01HXXX` | 3セグメント → `provider=local` 補完 |
| `memx:evidence:local:01HXXX` | `memx:evidence:local:01HXXX` | canonical（そのまま） |

#### 出力（フォーマッタ）

常に canonical format（4セグメント）を出力。

### 7.4 検証ルール

```go
// 検証条件
// 1. 4セグメントに split 可能
// 2. 全セグメントが空でない
// 3. domain が既知 namespace（memx, agent-taskstate, tracker）
// 4. 実在性確認は別責務（形式検証とは分離）
```

### 7.5 関連要件

- `FR-008` typed_ref 正規化（requirements-api.md#6-6）
- `AC-006` typed_ref 一貫性（requirements-nfr.md#5-5）
