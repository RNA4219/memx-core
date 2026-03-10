---
owner: memx-core
status: active
last_reviewed_at: 2026-03-06
next_review_due: 2026-06-06
---

# memx 要求事項 - API

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 6. API 要件（v1.3 追加）

- Requirement ID: `REQ-API-001`
- **実装状況: ✅ 完了（2026-03-06）**

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `memx_spec_v3/docs/quickstart.md`

### 6-1. 目的

- ツール/AI から呼びやすい **安定 JSON I/F** を提供する。
- CLI は API の薄いラッパとして実装し、ビジネスロジックを持たない。

### 6-2. 提供形態

- **HTTP（ローカル）**：`mem api serve` で起動。
  - 例：`http://127.0.0.1:7766`
  - 将来的に unix socket 対応も想定。
- **in-proc**：CLI や別バイナリが同一プロセスで呼ぶ。

### 6-3. エンドポイント（v1）
- ADR: [ADR-0002: v1必須3エンドポイント固定方針](../../docs/ADR/ADR-0002-v1-required-endpoints.md)

- `GET /healthz` → `ok` ✅
- `POST /v1/notes:ingest` ✅
  - request: `{title, body, summary?, source_type?, origin?, source_trust?, sensitivity?, tags?}`
  - response: `{note: Note}`
- `POST /v1/notes:search` ✅
  - request: `{query, top_k?}`
  - response: `{notes: Note[]}`
- `GET /v1/notes/{id}` → `Note` ✅
- `POST /v1/gc:run` ✅（SHOULD (v1.x): 実装完了）
  - request: `{target, options: {dry_run?}}`
  - response: `{status}`

**`POST /v1/gc:run` 実装状況**:
- Phase 0: トリガ判定 ✅
- Phase 3: Archive退避 ✅
- dry_run オプション ✅
- Feature flag: CLI側で `--enable-gc` または `--dry-run` で制御

### 6-3-1. v1必須3エンドポイント契約（`requirements.md` × `go/api/types.go` 照合）

#### POST `/v1/notes:ingest`

**request (`NotesIngestRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `title` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `body` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。v1 では削除・意味変更禁止。 |
| `summary` | `string` | 任意 | `""` | 追加バリデーションなし | 任意フィールドのため、未指定互換を維持。 |
| `source_type` | `string` | 任意 | `"manual"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定（クライアント未送信を維持）。 |
| `origin` | `string` | 任意 | `""` | trim のみ | 任意フィールドとして後方互換追加可。 |
| `source_trust` | `string` | 任意 | `"user_input"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `sensitivity` | `string` | 任意 | `"internal"` | trim 後、空ならデフォルト補完 | 既定値補完ルールを固定。 |
| `tags` | `[]string` | 任意 | `[]` 扱い（未指定時は処理なし） | 各要素 trim、空要素は無視 | 任意配列として未指定・空配列を同等扱い。 |

**response (`NotesIngestResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `note` | `Note` | 必須 | なし | 保存成功時に返却 | v1 では wrapper 構造（`{note: ...}`）を維持。 |
| `note.id` | `string` | 必須 | なし | 32 hex 形式の ID を生成 | 型・キー名固定。 |
| `note.title` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.summary` | `string` | 必須 | `""` 許容 | request 値 | 必須キーだが空文字を許容。 |
| `note.body` | `string` | 必須 | なし | request 値（trim 後） | 型固定。 |
| `note.created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.last_accessed_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `note.access_count` | `int64` | 必須 | `0` | ingest 直後は `0` | 数値型固定。 |
| `note.source_type` | `string` | 必須 | `"manual"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.origin` | `string` | 必須 | `""` | request 値 | 必須キーとして維持。 |
| `note.source_trust` | `string` | 必須 | `"user_input"` 補完あり | request/補完値 | 補完込みで常に返却。 |
| `note.sensitivity` | `string` | 必須 | `"internal"` 補完あり | request/補完値 | 補完込みで常に返却。 |

#### POST `/v1/notes:search`

**request (`NotesSearchRequest`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `query` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | 必須を維持。 |
| `top_k` | `int` | 任意 | `20`（`<=0` 時） | `<=0` は service で `20` に補正 | 任意数値として未指定互換を維持。 |

**response (`NotesSearchResponse`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `notes` | `[]Note` | 必須 | `[]` | 一致結果を最大 `top_k` 件返却 | v1 では配列 wrapper 構造を維持。 |
| `notes[].*` | `Note` 各フィールド | 必須 | `Note` 契約準拠 | `Note` と同一 | `Note` フィールドの型・キー名を維持。 |

#### GET `/v1/notes/{id}`

**request（path parameter）**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | trim 後に空文字不可（空の場合 `INVALID_ARGUMENT`） | path パラメータ必須を維持。 |

**response (`Note`)**

| field | type | required | default | validation | backward-compatibility note |
| --- | --- | --- | --- | --- | --- |
| `id` | `string` | 必須 | なし | ノート存在時に返却、非存在は `NOT_FOUND` | 型・キー名固定。 |
| `title` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `summary` | `string` | 必須 | `""` 許容 | 保存済み値 | 必須キーとして維持。 |
| `body` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `created_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `updated_at` | `string` | 必須 | なし | UTC RFC3339Nano | 日時文字列フォーマットを維持。 |
| `last_accessed_at` | `string` | 必須 | なし | 取得時刻で更新 | 取得時更新の挙動を維持。 |
| `access_count` | `int64` | 必須 | なし | 取得ごとに +1 | カウンタ型・加算挙動を維持。 |
| `source_type` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `origin` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `source_trust` | `string` | 必須 | なし | 保存済み値 | 型固定。 |
| `sensitivity` | `string` | 必須 | なし | 保存済み値 | 型固定。 |

### 6-4. エラーモデル
- ADR: [ADR-0003: ErrorCodeとretryable設計（再試行可/不可の境界）](../../docs/ADR/ADR-0003-errorcode-retryable-boundary.md)

- Requirement ID: `REQ-ERR-001`

本節を ErrorCode 契約の正本とし、`error-contract.md` は本節の運用向け要約として同期する。

共通：

```json
{"code":"INVALID_ARGUMENT","message":"...","details":{}}
```

- `INVALID_ARGUMENT` → 400
- `NOT_FOUND` → 404
- `CONFLICT` → 409
- `GATEKEEP_DENY` → 403
- `INTERNAL` → 500

`go/service/errors` と `go/api/errors.go` の現行方針（`ErrInvalidArgument` / `ErrNotFound` を明示マッピングし、それ以外は `INTERNAL` へフォールバック）を前提に、ErrorCode を 2 段で再定義する。

#### ErrorCode 区分（v1）

| 区分 | ErrorCode | HTTP | 契約レベル | 条件 |
| --- | --- | --- | --- | --- |
| v1必須保証 | `INVALID_ARGUMENT` | 400 | MUST | 常時有効。 |
| v1必須保証 | `NOT_FOUND` | 404 | MUST | 常時有効。 |
| v1必須保証 | `INTERNAL` | 500 | MUST | 常時有効。未分類エラーのフォールバック先。 |
| v1.x拡張（feature/sentinel依存） | `CONFLICT` | 409 | SHOULD | service sentinel（例: `ErrConflict`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |
| v1.x拡張（feature/sentinel依存） | `GATEKEEP_DENY` | 403 | SHOULD | gatekeeper deny 系 sentinel（例: `ErrGatekeepDeny`）を実装し `mapError` へ明示追加した場合のみ返却。未実装時は `INTERNAL`。 |

#### 現行実装との差分注記（`go/api/types.go` / `go/api/http_server.go`）

- `go/api/types.go` には `CONFLICT` / `GATEKEEP_DENY` 定数が定義済みだが、返却は service sentinel と `mapError` 実装に依存する。
- `go/api/http_server.go` は `writeErr` で `CONFLICT=409` / `GATEKEEP_DENY=403` を処理可能だが、上流が当該コードを返さない限り到達しない。
- sentinel 未実装時は `mapError` 方針に従って `INTERNAL`（500）へフォールバックする。

| 代表事象 | service 層の分類 | API `code` | HTTP | 再試行可否 | 備考 |
| --- | --- | --- | --- | --- | --- |
| 入力不備（必須欠落・形式不正・空文字） | `ErrInvalidArgument` | `INVALID_ARGUMENT` | 400 | 不可 | クライアント入力修正が必要。 |
| DB ロック（`database is locked` / busy timeout 超過） | 現行は汎用エラー（将来 `ErrConflict` 候補） | `INTERNAL`（将来 `CONFLICT` 検討可） | 500（将来 409 検討可） | 条件付き可 | 短時間で解消しうるため指数バックオフ再試行対象。長時間継続時は運用アラート。 |
| LLM タイムアウト（上流 timeout / 502/503/504） | 汎用エラー | `INTERNAL` | 500 | 条件付き可 | 一時障害として再試行対象。最大回数超過で失敗扱い。 |
| Gatekeeper deny / needs_human(v1.3 では deny 扱い) | 将来 sentinel 化（`ErrGatekeepDeny`） | `GATEKEEP_DENY`（未実装時は `INTERNAL`） | 403（未実装時は 500） | 不可 | ポリシー判定のため再試行では解消しない。入力/運用判断が必要。 |

#### エラーコード別 再試行可否表

| API `code` | HTTP | 再試行可否 | ルール |
| --- | --- | --- | --- |
| `INVALID_ARGUMENT` | 400 | 不可 | 入力を修正して再実行。 |
| `NOT_FOUND` | 404 | 不可 | 対象 ID・クエリを見直す。 |
| `CONFLICT` | 409 | 条件付き可 | DB ロック等の一時競合のみ再試行。恒久競合は不可。 |
| `GATEKEEP_DENY` | 403 | 不可 | ポリシー deny は即時失敗。再試行禁止。 |
| `INTERNAL` | 500 | 条件付き可 | 原因が一時障害（DB lock / LLM timeout / upstream 502/503/504）の場合のみ指数バックオフ。恒久障害は不可。 |

### 6-5. 互換性ポリシー（v1）

- v1 系（`/v1/*`）では、既存クライアントを壊さない後方互換を維持する。

許容する変更（v1 内）：

- request/response への**任意フィールド追加**（既存必須フィールドは維持）。
- 新規エンドポイント追加（既存エンドポイントの契約は維持）。
- 既存フィールドのバリデーション強化のうち、既存の正常系入力を拒否しない変更。

禁止する変更（v1 内）：

- 既存の必須フィールド削除。
- 既存フィールドの意味変更（同じキー名で別意味にする変更）。
- 既存レスポンスの型変更（例：`string` → `object`）や、既存成功レスポンスの構造破壊。

破壊変更が必要な場合の移行手順：

1. まずは **新エンドポイント**（例：`/v1/notes:search2`）で並行提供し、既存挙動を維持する。
2. 並行提供で吸収できない場合は **`/v2` を新設**し、v1 を非推奨化する。
3. `CHANGES.md` に移行期限・差分・移行例を記載し、少なくとも 1 リリースは移行猶予を置く。

CLI 出力との互換責任範囲：

- 人間向けの通常表示（非 `--json`）は、可読性改善のための文言・並び変更を許容する。
- `--json` 出力は API の互換ポリシーと同等に扱い、v1 では破壊的変更を禁止する。
- 機械連携は API または CLI `--json` を正とし、互換責任はこの 2 系統に対して負う。

### 6-6. typed_ref 正規化（FR-008）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-1-typed_ref-canonical-format-の固定`

- Requirement ID: `FR-008`

システムは、外部状態・bundle・根拠参照で用いる `typed_ref` を canonical format で保存・返却しなければならない。

#### canonical format

```txt
<domain>:<entity_type>:<provider>:<entity_id>
```

#### 例

| domain | entity_type | provider | entity_id | canonical format |
|--------|-------------|----------|-----------|------------------|
| `memx` | `evidence` | `local` | `01HXXX` | `memx:evidence:local:01HXXX` |
| `memx` | `artifact` | `local` | `01HYYY` | `memx:artifact:local:01HYYY` |
| `memx` | `knowledge` | `local` | `01HZZZ` | `memx:knowledge:local:01HZZZ` |
| `agent-taskstate` | `task` | `local` | `task_01J...` | `agent-taskstate:task:local:task_01J...` |
| `tracker` | `issue` | `jira` | `PROJ-123` | `tracker:issue:jira:PROJ-123` |
| `tracker` | `issue` | `github` | `owner/repo#123` | `tracker:issue:github:owner/repo#123` |

#### 必須条件

1. 新規生成される ref は canonical format であること
2. migration 期間中は旧 ref（3セグメント形式）を canonical format へ正規化可能であること
3. 実在性確認と形式妥当性確認を分離すること

#### 移行期の read-both / write-one

- **パーサー**: 一時的に 3セグメント（`memx:<type>:<id>`）も受理し、`provider=local` として正規化する
- **フォーマッタ**: 常に 4セグメント canonical format を出力する
- **DB**: 既存の 3セグメント文字列を直ちに破壊しない（移行完了後に一括更新）

#### 検証ルール

- 4セグメントに split 可能であること
- `domain`, `entity_type`, `provider`, `entity_id` がいずれも空でないこと
- `domain` は既知 namespace（`memx` / `agent-taskstate` / `tracker`）のいずれかであること
- 実在性確認は別責務とする（形式検証とは分離）

### 6-7. 継続用 bundle 保存（FR-006）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-3-context-bundle-の必須監査項目の明確化`

- Requirement ID: `FR-006`

継続用 context bundle は、再開・監査・再現性のために以下の必須項目を持たなければならない。

#### bundle 必須項目

| 項目 | 説明 | 必須 |
|------|------|------|
| `purpose` | bundle 生成目的 | 必須 |
| `rebuild_level` | 再構成レベル（summary/raw/full） | 必須 |
| `summary` | bundle 内容の要約 | 必須 |
| `state_snapshot` | 現在状態のスナップショット | 必須 |
| `decision_digest` | 判断の要約 | 必須 |
| `open_question_digest` | 未解決質問の要約 | 必須 |
| `source_refs` | 参照元の typed_ref リスト | 必須 |
| `raw_included_flag` | raw データ含有フラグ | 必須 |
| `generator_version` | bundle 生成器のバージョン | 必須 |
| `generated_at` | 生成タイムスタンプ | 必須 |
| `diagnostics` | 診断情報 | 必須 |

#### diagnostics 必須項目

| 項目 | 説明 |
|------|------|
| `missing_refs` | 解決できなかった ref リスト |
| `unsupported_refs` | 未対応の ref リスト |
| `partial_bundle_flag` | 部分的な bundle かどうか |
| `resolver_warnings` | resolver からの警告 |

### 6-8. 状態遷移明示化（FR-007）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-2-current-dashboard-と-state-history-の分離`

- Requirement ID: `FR-007`

システムは、task の進行状態を以下の 2 層で保持しなければならない。

#### 2 層構造

1. **materialized current state**: 最新状態を示すダッシュボード
2. **state transition history**: append-only の状態遷移履歴

#### 状態遷移ルール

- 状態変更は暗黙更新ではなく、履歴付き遷移として記録されなければならない
- `task.status` の変更は専用の状態遷移処理を経由し、直接更新してはならない

#### 状態遷移履歴の必須項目

| 項目 | 説明 |
|------|------|
| `from_state` | 遷移元状態 |
| `to_state` | 遷移先状態 |
| `reason` | 遷移理由 |
| `actor` | 遷移実行者 |
| `related_run_ref` | 関連 run の typed_ref |
| `changed_at` | 遷移タイムスタンプ |

#### 失敗条件

- current state は復元できるが state history が失われている
- `task.status` が直接更新され、履歴付き遷移として辿れない

### 6-9. Context Rebuild の外部依存制約（FR-003）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-5-tracker-情報は-optional-input-であり必須依存ではないことの明記`

- Requirement ID: `FR-003`

tracker issue snapshot や外部 issue 情報は optional input として利用してよいが、task 再開の必須条件としてはならない。

#### 許可される依存

| 入力 | 必須/任意 | 備考 |
|------|----------|------|
| 内部状態（task/state/decision） | 必須 | KV Cache Independence の核心 |
| memx-core（evidence/knowledge/artifact） | 必須（段階的導入） | Phase 2 以降で完全必須化 |
| tracker issue snapshot | 任意 | 利用可能だが欠けても再開可能 |
| 外部 issue 情報 | 任意 | 同期されていなくても継続可能 |

#### 失敗条件

- 外部 tracker 情報がないと task を再開できない
- 外部 issue snapshot が欠けると current state を復元できない

#### 非対象

外部 tracker を内部作業状態の正本として扱うことは本要求の対象外とする。

### 6-10. 競合検出（FR-009）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-6-stale-state--stale-bundle-への競合制御`

- Requirement ID: `FR-009`

システムは、stale state または stale bundle に基づく更新を検出可能でなければならない。

#### 競合検出機構

最低限、以下のいずれかを持たなければならない。

| 機構 | 説明 |
|------|------|
| `state_revision` | 状態の版番号 |
| `task_version` | task のバージョン |
| `expected_current_state` | 期待される現在状態のハッシュ |
| `bundle_generated_at` | bundle 生成時刻 |
| `source_snapshot_version` | 参照元スナップショットの版 |

#### 競合時の挙動

- 不一致時は **暗黙 merge を行わず、競合として扱わなければならない**
- 競合検出時は operator へ通知し、明示的な解決を求める

#### 失敗条件

- stale bundle からの更新を検出できない
- 古い current state に基づいて superseded decision や古い next action を採用してしまう

### 6-11. Tracker 統合（FR-010）

> Source: `docs/kv-priority-roadmap/04-tracker-bridge-minimum-integration.md`

- Requirement ID: `FR-010`

外部 issue を「外界の窓口」として agent-taskstate / memx-core に接続する。ただし tracker を正本にしない。

#### 境界

**tracker-bridge が持つもの:**
- `tracker_connection`: 外部 tracker 接続情報
- `issue_cache`: 取得した issue のキャッシュ
- `entity_link`: issue と agent-taskstate task のリンク
- `sync_event`: 同期イベント履歴

**tracker-bridge が持たないもの:**
- task state の正本
- decision の正本
- evidence / knowledge の正本
- context build policy

#### MVP スコープ

| 機能 | 内容 |
|------|------|
| **inbound** | issue fetch / normalize / cache 更新 / sync_event 記録 |
| **linking** | `tracker:issue:*:*` と `agent-taskstate:task:local:*` のリンク（role: primary/related/duplicate/blocks） |
| **outbound** | status 反映 / short comment 反映 |
| **snapshot export** | context build 用の最小 snapshot 提供 |

#### snapshot 必須項目

| 項目 | 説明 |
|------|------|
| `remote_key` | 外部 issue の一意キー |
| `title` | issue タイトル |
| `status` | issue 状態 |
| `assignee` | 担当者 |
| `updated_at` | 更新日時 |
| `last_sync_result` | 最終同期結果 |

#### 受入条件

1. 外部 issue を 1 件 fetch して `issue_cache` に保存できる
2. issue と `agent-taskstate task` を entity_link できる
3. inbound / outbound の sync_event を追跡できる
4. `agent-taskstate context build` が tracker snapshot を optional input として使える
5. tracker を見失っても `agent-taskstate` と `memx-core` だけで task 継続は可能

#### typed_ref 形式

```
tracker:issue:jira:PROJ-123
tracker:issue:github:owner/repo#123
```

#### 関連要件

- `FR-003` Context Rebuild の外部依存制約（requirements-api.md#6-9）
- `FR-008` typed_ref 正規化（requirements-api.md#6-6）