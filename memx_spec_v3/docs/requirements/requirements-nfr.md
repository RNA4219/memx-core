---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# memx 要求事項 - 非機能要件

> 本書は `requirements.md` から分割された一部です。正本は `requirements.md` を参照してください。

## 5. 非機能要件

### 5-1. 性能目標（v1必須3エンドポイント）

計測条件は以下で固定する。

- データセット条件: short ストア 10,000 件、本文 1 ノートあたり約 500 文字（UTF-8 プレーンテキスト）
- 実行環境: ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）
- ウォームアップ: 各エンドポイント 20 回
- 本計測: 各エンドポイント 200 回（ウォームアップ除外）

| 操作 | API | p50 目標 | p95 目標 |
| --- | --- | --- | --- |
| ingest | `POST /v1/notes:ingest` | `<= 120ms` | `<= 250ms` |
| search | `POST /v1/notes:search` | `<= 80ms` | `<= 180ms` |
| show | `GET /v1/notes/{id}` | `<= 40ms` | `<= 90ms` |
- Requirement ID: `REQ-NFR-001`

#### 計測プロトコル

- データ生成条件（固定）:
  - 対象ストアは `short` 固定。
  - 事前投入データは 10,000 件固定。
  - 各ノート本文は約 500 文字（UTF-8 プレーンテキスト）固定。
- 計測回数（固定）:
  - 各エンドポイントごとにウォームアップ 20 回を先に実行。
  - 本計測は各エンドポイント 200 回で固定し、ウォームアップを集計から除外する。
- 除外条件（固定）:
  - ウォームアップ 20 回は母集団に含めない。
  - `REQ-NFR-001` 判定では `ingest` / `search` / `show` の 3 エンドポイント以外を除外する。
  - ローカル単体（4 vCPU / 16GB RAM / NVMe SSD / Linux x86_64）以外の環境で得た計測値は正式判定から除外する。
- 集計方法（固定）:
  - 本計測 200 回のレイテンシ分布から `p50` と `p95` を算出する。
  - 合否判定値は `p50_ms` / `p95_ms` のみとし、平均値や最大値は参考情報として扱う。
  - 出力フォーマットは `RUNBOOK.md` の `artifacts/perf/perf-result.json` に定義されたキーを正本とする。

### Dependencies

- `BLUEPRINT.md`
- `EVALUATION.md`
- `GUARDRAILS.md`
- `RUNBOOK.md`

- OS：ローカル（Linux / macOS / Windows）で動作する CLI を想定。
- DB：SQLite3（WAL モード、foreign_keys ON）。
- 言語：Go（単一バイナリビルドを前提）。
- 依存：
  - SQLite ドライバ（例：`modernc.org/sqlite` / `github.com/mattn/go-sqlite3`）。
  - CLI は標準 `flag` でもよい（薄いラッパが前提）。
  - HTTP API は標準 `net/http` でよい。
- セキュリティ：
  - APIキーや秘密情報は `memory_policy.yaml` + Gatekeeper により保存前にブロック。
- 拡張性：
  - `journal.db` / `knowledge.db` / `archive.db` のスキーマは、`short.db` と同様の `notes` / `tags` / `note_tags` / `note_embeddings` / `notes_fts` 構造をベースとする。
  - 将来、Working Memory／プロジェクトスコープ／セッションスコープを追加しても、既存の CLI と DB 構造を壊さない。
- トランザクションと一貫性：
  - ATTACH を跨いだ完全な原子的トランザクションは SQLite の仕様上保証できない。
- 設計上「データ喪失より重複を許容する」ポリシーとし、lineage による追跡・再蒸留で整合性を取り戻せるようにする。

### 5-2. 可用性・復旧・整合性回復（運用NFR）

| Requirement ID | 区分 | 要件 |
| --- | --- | --- |
| `REQ-NFR-002` | 可用性/復旧目標 | 障害時の目標復旧時間は `RTO <= 30分`、目標復旧時点は `RPO <= 5分` を満たすこと。 |
| `REQ-NFR-003` | 検知〜暫定復旧 | 障害検知から暫定復旧（サービス縮退を含む）完了までの上限時間を `15分` とする。 |
| `REQ-NFR-004` | 再処理上限 | 再試行/再処理は 1 リクエスト（または 1 ノート）あたり `最大 2 回` までとし、3 回目以降は自動再試行を禁止して運用エスカレーションする。 |
| `REQ-NFR-005` | 整合性回復時間 | `short→archive` 補償フローの整合性回復は、障害検知から `30分以内` に収束判定（または `docs/IN-*.md` 起票）へ到達すること。 |

補足:
- `RPO` は障害復旧後に再投入が必要なデータ欠損許容時間を指す。
- `RTO` は障害検知時点から正常系または暫定系サービス復帰までの許容時間を指す。
- 再処理方針は `at-least-once` を採用し、重複は許容するが欠損は許容しない（重複は `REQ-NFR-005` の収束条件で解消する）。

#### 5-2-1. RTO/RPO 判定の固定ルール

- `RTO` の起点は `detected_at`、終点は `mitigated_at`（暫定復旧）または `resolved_at`（恒久復旧）の先着時刻とする。
- `RPO` は障害復旧後に再投入が必要だった最古データ時刻と `detected_at` の差分で算出する。
- `REQ-NFR-002` の合否は、同一インシデントに対して `rto_minutes <= 30` かつ `rpo_minutes <= 5` の同時成立を必須とする。
- 判定証跡は `artifacts/ops/incident-summary.json` を正本とし、手計算値は参考情報扱いとする。

#### 5-2-2. 再試行方針（運用固定）

- 自動再試行の対象は一時障害（DB lock、LLM timeout、HTTP 429/502/503/504）のみとする。
- 再試行回数は 1 リクエスト（または 1 ノート）あたり最大 2 回、待機は指数バックオフ（`1s -> 2s`、ジッタ許容）を推奨値とする。
- `INVALID_ARGUMENT` / `NOT_FOUND` / `GATEKEEP_DENY` / 恒久障害の `INTERNAL` は再試行禁止とし、運用エスカレーションへ遷移する。
- 再試行打ち切り後は `rollback` または `replan` を必須化し、`docs/IN-*.md` に再試行回数と打ち切り理由を記録する。

### 5-3. 整合性回復要件（Archive 補償フロー）

- Requirement ID: `REQ-NFR-005`
- `mem gc short` の再実行で、`lineage` と `archive.notes` の不整合を `30分以内` に収束させる。
- 整合性回復の 1 サイクルで実施する再処理回数は `最大 2 回`（`REQ-NFR-004` 準拠）とし、未収束時は `docs/IN-*.md` 起票を必須とする。
- 収束条件は「`archive 実在 + 対応 lineage 実在` を満たし、short 側 Delete 判定が再開可能」であること。

収束判定（数値/状態）:
1. 対象バッチの `pending_compensation_count == 0`（未補償 0 件）。
2. 対象 `src_note_id` ごとに `archive 実在 + archived_from lineage 実在` が `1 組以上` 成立。
3. 同一 `src_note_id` の重複 archive は `dup_archive_count <= 1`（許容重複 1 件）まで削減済み、または削減不能理由を `docs/IN-*.md` に記録済み。
4. `short_delete_ready_ratio == 1.0`（Delete 判定対象の全件が削除可能状態）。
5. 上記 1〜4 を障害検知から 30 分以内に満たせない場合は「未収束」とし、`docs/IN-*.md` 起票と再計画チケット発行を必須とする。

状態定義（重複許容後の収束条件）:

| 状態ID | 名称 | 判定条件 | 次状態 |
| --- | --- | --- | --- |
| `S0` | 検知直後 | `pending_compensation_count > 0` | `S1` |
| `S1` | 補償実行中 | `retry_count <= 2` かつ `archive+lineage` の片系不足が存在 | `S2` or `S3` |
| `S2` | 重複許容安定 | `dup_archive_count >= 1` を許容しつつ `欠損=0` を維持 | `S4`（`pending_compensation_count==0` かつ `short_delete_ready_ratio==1.0` かつ `dup_archive_count<=1`）or `S3` |
| `S3` | 未収束 | `retry_count > 2` または 30分超過で `pending_compensation_count > 0` | `S5` |
| `S4` | 収束完了 | `pending_compensation_count == 0` かつ `short_delete_ready_ratio == 1.0` かつ `dup_archive_count <= 1` | 終端 |
| `S5` | 要起票終端 | `docs/IN-*.md` 起票 + 再計画チケット発行済み | 終端 |

- `S2` は「重複許容の暫定安定状態」とし、欠損ゼロを維持したまま `S4` へ収束させる中間状態とする。
- `S2` のまま 30 分を超過した場合は `S3`（未収束）へ遷移し、`S5` へ進める。
- `S2 -> S4` 遷移は「重複許容後の最終収束条件（未補償ゼロ・Delete 再開可・重複許容上限内）」を同時充足した場合に限定する。

### 5-4. インシデント記録（`docs/IN-*.md`）最小監査項目

- Requirement ID: `REQ-NFR-006`
- 受入対象の実運用インシデント記録は `docs/IN-<実日付>-<連番>.md` 形式のみとする。
- `docs/IN-*.md` には、以下の監査項目を必須記録する。
  1. 事象識別子: `インシデントID` / `発生日` / `起票日` / `重大度` / `ステータス`
  2. 要件トレーサビリティ: `関連要件ID` / `要件違反有無` / `違反した要件IDまたは節`
  3. 時間監査: `検知日時` / `暫定復旧完了日時` / `恒久復旧完了日時`
  4. 復旧行動監査: `再試行回数` / `ロールバック実施有無` / `再計画チケットID`
  5. 影響監査: `影響対象` / `影響期間` / `影響規模` / `CIA影響`
  6. 証跡: `関連ログ・メトリクス・判定結果ファイル` の保存先パス

#### 5-4-1. waiver 時の必須記録（`docs/IN-*.md` 運用連動）

- waiver を許容する場合でも、記録媒体は必ず `docs/IN-<実日付>-<連番>.md` とする（口頭/チャットのみは不可）。
- 必須項目の記載フォーマットは `docs/IN-BASELINE.md` / `docs/IN-YYYYMMDD-001.md` / `docs/IN-202603xx-001.md` の waiver セクションに準拠する。これら3文書はテンプレート資料であり、要件根拠としては参照禁止とする。
- `EVALUATION.md` の運用NFR合否判定で機械的に参照できること。
- 必須記録項目:
  1. waiver対象要件ID
  2. waiver理由（技術的制約/外部依存/緊急運用の別）
  3. 期限（UTC、失効日時）
  4. 暫定リスク受容者（承認者）
  5. 代替統制（監視強化・手動運用手順・追加検証）
  6. 解除条件（どの証跡が揃えば waiver を解消するか）
  7. 関連証跡パス（`artifacts/ops/incident-summary.json`、`artifacts/ops/recovery-log.ndjson` など）

### 5-5. typed_ref 一貫性（AC-006）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-1-typed_ref-canonical-format-の固定`

- Requirement ID: `AC-006`

`agent-taskstate` / `memx-core` / `tracker-bridge` をまたぐ主要参照が、同一 canonical `typed_ref` 形式で保存・再利用・追跡可能であること。

#### 受入条件

1. **出力整合**: 新規生成される全ての typed_ref が canonical format（4セグメント）であること
2. **読込互換**: 移行期間中、旧形式（3セグメント）の typed_ref を正規化して読み込めること
3. **追跡可能性**: `source_refs`、`entity_link`、`lineage`、`context_bundle_source` で同じ ref 形式を共有できること
4. **検証一元化**: 3システムが同じ validation rule を持つこと

#### 検証コマンド

```bash
# memx-core での typed_ref 形式検証
go test ./memx_spec_v3/go/api/... -run TestTypedRef

# canonical format の確認
rg -n "memx:[a-z_]+:[a-z]+:" memx_spec_v3/go/ --type go
```

#### 関連要件

- `FR-008` typed_ref 正規化（requirements-api.md#6-6-typed_ref-正規化fr-008）

### 5-6. 再現性・可監査性（NFR-001 / NFR-002）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-3-context-bundle-の必須監査項目の明確化`

#### NFR-001 再現性

- Requirement ID: `NFR-001`

同一状態集合と同一 `generator_version` から再構成される bundle は、意味的に再現可能でなければならない。

**判定条件:**
- 同一入力に対して同一 `generator_version` で生成された bundle は、内容が意味的に等価であること
- `source_refs` の順序違いは許容するが、欠落・過剰は許容しない

#### NFR-002 可監査性

- Requirement ID: `NFR-002`

bundle 本体だけでなく、bundle 生成時に使用された `source_refs`、`generator_version`、`diagnostics` を後から監査可能でなければならない。

**監査可能項目:**
- `source_refs`: bundle の参照元を一意に特定可能
- `generator_version`: 生成時のロジックを追跡可能
- `diagnostics`: 生成時の警告・欠損を確認可能
- `generated_at`: 生成時刻を記録

### 5-7. Bundle 監査性（AC-007）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-3-context-bundle-の必須監査項目の明確化`

- Requirement ID: `AC-007`

任意の bundle について、以下の監査項目を確認できること。

#### 必須監査項目

| 項目 | 確認内容 |
|------|---------|
| `purpose` | 生成目的 |
| `source_refs` | 参照元の typed_ref リスト |
| `raw_included_flag` | raw データ含有有無 |
| `generator_version` | 生成器バージョン |
| `diagnostics` | 診断情報（missing/unsupported refs 等） |

#### 検証コマンド

```bash
# bundle 監査項目の確認
rg -n "purpose|source_refs|raw_included|generator_version|diagnostics" <bundle-path>
```

#### 関連要件

- `FR-006` 継続用 bundle 保存（requirements-api.md#6-7）
- `FR-007` 状態遷移明示化（requirements-api.md#6-8）
- `NFR-001` 再現性（本節）
- `NFR-002` 可監査性（本節）

### 5-8. 劣化耐性（NFR-004 追記）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-6-stale-state--stale-bundle-への競合制御`

- Requirement ID: `NFR-004`

再開時に古い state や bundle が使われた場合でも、誤更新より競合検出を優先しなければならない。

#### 優先順位

1. **競合検出**: 古いデータに基づく更新を防ぐ
2. **データ保護**: 意図しない上書きを防ぐ
3. **再開継続**: 可能な限り再開を支援する（競合解決後）

#### 検証コマンド

```bash
# 競合検出機構の確認
rg -n "state_revision|task_version|bundle_generated_at" <state-path>
```

#### 関連要件

- `FR-009` 競合検出（requirements-api.md#6-10）

### 5-9. 段階的導入整合（AC-008）

> Source: `docs/kv-priority-roadmap/kv-cache-independence-amendments.md#追記案-4-memx-core-依存の段階化`

- Requirement ID: `AC-008`

memx-core 依存の導入は段階的に評価してよい。

#### 導入フェーズ

| Phase | 内容 | 必須度 |
|-------|------|--------|
| Phase 1 | task 継続に必要な構造状態を外部化 | 必須 |
| Phase 2 | summary-first で evidence/knowledge/artifact を再構成に利用 | 必須 |
| Phase 3 | lineage/journal/distilled knowledge により根拠追跡を強化 | 推奨 |

#### 受入条件

1. 初期段階では `work state` のみで再開可能であってもよい
2. 後続段階では summary-first の memory integration を通じて evidence/knowledge/artifact を bundle に取り込めること
3. 最終的には raw evidence と distilled knowledge の双方へ辿れる構造を持たなければならない

#### 記憶基盤連携原則

- summary-first を原則とする
- raw データは必要時のみ selected inclusion とする
- 段階的に機能を導入し、各段階で動作確認を行う

#### 検証コマンド

```bash
# Phase 2 以降の確認
rg -n "evidence|knowledge|artifact" <bundle-path>
```

#### 関連要件

- `FR-003` Context Rebuild の外部依存制約（requirements-api.md#6-9）
- `FR-006` 継続用 bundle 保存（requirements-api.md#6-7）