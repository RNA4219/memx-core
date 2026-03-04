# Link Integrity Spec

## 1. 目的
- Phase 3 のリンク健全性判定を、再現可能な基準で実施する。
- 判定語彙（重大度・Status 遷移）は `docs/TASKS.md` と同一語彙を使う。

## 2. 検証対象

### 2.1 入力対象
- `memx_spec_v3/docs/*.md`
- `docs/*.md`
- `docs/birdseye/index.json`

### 2.2 リンク分類
入力対象から抽出した参照は、以下 3 分類で検証する。

1. 章内アンカー（Markdown Anchor）
   - 例: `#phase-3`, `./requirements.md#2-2-変更タイプ別チェックリスト`
2. 相対パス（Relative Path）
   - 例: `../memx_spec_v3/docs/design.md`, `./IN-20260303-001.md`
3. 契約 JSON Pointer（Contract JSON Pointer）
   - 例: `memx_spec_v3/docs/contracts/openapi.yaml#/paths/~1memories/get`
   - 例: `memx_spec_v3/docs/contracts/cli-json.schema.json#/properties/output`

## 3. 判定ルール

### 3.1 章内アンカー
- 同一ファイル内またはリンク先 Markdown 内に、対応する見出しアンカーが 1 件以上存在すること。
- 見出し正規化（小文字化・空白/記号のハイフン化）後に一致判定すること。

### 3.2 相対パス
- 参照元基準で正規化した最終パスが実在すること。
- 正規化後パスがリポジトリルート配下に収まること（`..` でルート外へ逸脱しない）。
- Markdown から Markdown/JSON を参照する場合、拡張子が `.md` または `.json` であること。

### 3.3 契約 JSON Pointer
- 参照先ファイルが実在し、JSON/YAML として解釈可能であること。
- `#/<pointer>` で指定されたノードが 1 件に解決できること。
- Pointer が配列/オブジェクト境界を跨ぐ場合も RFC 6901 のエスケープ（`~0`, `~1`）を満たすこと。

## 4. 失敗条件
以下をリンク健全性チェックの失敗とする。

1. リンク先不存在
   - 相対パス正規化後にファイルが存在しない。
2. アンカー未解決
   - 対象 Markdown 内に見出しアンカーが存在しない。
3. 拡張子誤り
   - Markdown/JSON 契約参照なのに `.md`/`.json` 以外を指す。
4. 相対パス階層逸脱
   - `../` 解決後にリポジトリルート外となる。
5. JSON Pointer 未解決
   - 対象契約ファイル内に Pointer 先ノードがない。
6. JSON Pointer 構文誤り
   - RFC 6901 非準拠の Pointer で解決不能。

## 5. 重大度（severity）

### high
- 相対パス階層逸脱
- 契約 JSON Pointer 未解決/構文誤り
- 参照先不存在（契約・運用導線を阻害する主要リンク）

### medium
- アンカー未解決（章間遷移は可能だが、節への直接到達が不能）
- 参照先不存在（補助説明リンク）

### low
- 拡張子誤り（実体は存在し内容把握は可能だが、運用規約違反）
- 将来互換上の軽微な表記ゆれ（同一実体への重複リンク等）

## 6. Status 遷移規定（`docs/TASKS.md` 準拠）
- 許可語彙: `planned` / `active` / `in_progress` / `reviewing` / `blocked` / `done`

遷移規定:
1. `reviewing -> blocked`
   - `severity: high` が 1 件でも検出された場合。
2. `reviewing` 継続
   - `severity: medium` のみ残存し、修正 Task Seed を起票済みの場合。
3. `reviewing -> done`
   - high=0 かつ medium=0（low は許容しない運用とし 0 件化する）。
4. `in_progress -> reviewing`
   - 本仕様に基づくリンク健全性チェック結果を Task Seed の `Requirements` / `Commands` へ反映済みの場合。

## 7. 出力フォーマット
リンク不整合は 1 件につき以下で記録する。

```yaml
check_id: LI-YYYYMMDD-001
severity: low | medium | high
status_transition: reviewing->blocked
source_ref: memx_spec_v3/docs/design.md#L10-L40
target_ref: ../docs/unknown.md#phase-1
failure_type: missing_target | unresolved_anchor | bad_extension | path_escape | pointer_not_found | pointer_syntax_error
action: "リンク先追加または参照先修正"
```

## 8. Phase 3 運用
- `orchestration/memx-design-docs-authoring.md` の Phase 3 では、本仕様をリンク健全性の正規基準として参照する。
- 判定結果は章別 Task Seed の `Requirements` / `Commands` / `Status` 判定に反映する。
