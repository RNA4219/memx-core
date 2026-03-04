# Design Review Records

## 保存先
- 設計レビュー記録は `memx_spec_v3/docs/reviews/` に保存する。

## 命名規則
- ファイル名は `DESIGN-REVIEW-YYYYMMDD-###.md` とする。
  - `YYYYMMDD`: レビュー実施日（ローカル日付）
  - `###`: 同日内の 001 始まり連番
- 例: `DESIGN-REVIEW-20260304-001.md`

## 更新ルール
- 新規レビューごとに `TEMPLATE.md` を複製して新規ファイルを作成する（既存記録の上書き禁止）。
- 記録作成時は `design-review-spec.md` の必須 6 項目（対象章/関連 REQ-ID/Node IDs/指摘一覧（重大度付き）/再確認結果/判定）を必ず記入する。
- 判定欄には `EVALUATION.md` の該当ルール参照と証跡を必ず記載する。
- `docs/TASKS.md` 連携項目（`Release Note Draft` / `Status` / `Moved-to-CHANGES`）を記録クローズ前に更新する。
