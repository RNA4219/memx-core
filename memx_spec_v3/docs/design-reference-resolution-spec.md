---
owner: memx-core
status: active
last_reviewed_at: 2026-03-04
next_review_due: 2026-06-04
---

# design reference resolution spec

`orchestration/memx-design-docs-authoring.md` で使う入力参照名を、Task Seed 作成時に一意な正規パスへ解決するための仕様。

## 0. 適用スコープ

本仕様の参照解決ルールは、以下の成果物すべてに適用する。

- Task Seed
- Phase 1 抽出表（Design Source Inventory）
- 章ドラフト
- レビュー記録の `Source` 欄

## 1. 対象入力参照名と正規パスマッピング

| 入力参照名（非正規） | 正規パス（必須） |
| --- | --- |
| `requirements.md` | `memx_spec_v3/docs/requirements.md` |
| `design.md` | `memx_spec_v3/docs/design.md` |
| `interfaces.md` | `memx_spec_v3/docs/interfaces.md` |
| `traceability.md` | `memx_spec_v3/docs/traceability.md` |
| `contracts.md` | `memx_spec_v3/docs/CONTRACTS.md` |
| `CONTRACTS.md` | `memx_spec_v3/docs/CONTRACTS.md` |
| `EVALUATION.md` | `docs/birdseye/caps/EVALUATION.md.json` |
| `RUNBOOK.md` | `docs/birdseye/caps/RUNBOOK.md.json` |
| `docs/birdseye/index.json` | `docs/birdseye/index.json` |

### 1-2. EVALUATION/RUNBOOK 正本パス規約

- 正本（canonical source）はリポジトリルートの **`EVALUATION.md`** / **`RUNBOOK.md`** に固定する。
- `memx_spec_v3/docs/` 配下に `EVALUATION.md` / `RUNBOOK.md` を正本として配置・参照してはならない。
- 設計仕様での `Source` / `Dependencies` / 入力成果物記述は、評価・運用参照を必ず以下へ正規化する。
  - `docs/birdseye/caps/EVALUATION.md.json`
  - `docs/birdseye/caps/RUNBOOK.md.json`
- `memx_spec_v3/docs/design-phase-gate-spec.md` と `orchestration/memx-design-docs-authoring.md` の入力成果物記述は、上記正規表記に統一する。

### 1-1. `contracts.md` / `CONTRACTS.md` の扱い（固定）

- 入力参照名 `contracts.md` と `CONTRACTS.md` は、表記ゆれとして同一扱いにする。
- 正本（canonical source）は **`memx_spec_v3/docs/CONTRACTS.md`** に固定し、他パスへの解決を禁止する。
- `Source` / `Dependencies` / レビュー記録で `memx_spec_v3/docs/contracts.md`（小文字）を検出した場合は誤参照として fail 扱いにする。

## 2. `Source` 正規化ルール（`path#Section` 統一）

Task Seed / 章ドラフトの `Source` は、以下をすべて満たすこと。

1. `path#Section` 形式を必須とする（`#Section` が不要な場合も末尾に `#` ではなく章名を明示する）。
2. 相対名（例: `requirements.md#...`）を禁止し、上表の正規パスへ解決した絶対的なリポジトリ相対パスのみ許可する。
3. 曖昧名（例: `EVALUATION.md#...`、`RUNBOOK.md#...`）を禁止する。
4. 複数候補へ解決される参照名は自動補完せず fail とし、Task Seed を `reviewing` のまま差し戻す。

### 2-1. `docs/birdseye/caps/*.json` 前提の `path#section` 解決

- `EVALUATION.md#...` / `RUNBOOK.md#...` の入力は、`path#section` へ解決する際に必ず以下へ変換する。
  - `docs/birdseye/caps/EVALUATION.md.json#<section_key>`
  - `docs/birdseye/caps/RUNBOOK.md.json#<section_key>`
- `section` は Markdown 見出し名ではなく、caps JSON 側のキー（snake_case）を使用する。
- `memx_spec_v3/docs/EVALUATION.md#...` / `memx_spec_v3/docs/RUNBOOK.md#...` は不正な `path#section` として fail 扱いにする。

## 3. Task Seed 作成時の必須チェック

Task Seed 起票時に、以下チェックをすべて通過しなければならない。

- `Source` の全行が本仕様の正規パスマッピングに一致している。
- `Source` の全行が `path#Section` 形式で、`#Section` が空でない。
- 相対名・曖昧名・複数候補解決のいずれも 0 件。

## 4. 運用例（誤/正）

- 誤: `requirements.md`
  - 正: `memx_spec_v3/docs/requirements.md#6-4. エラーモデル`
- 誤: `design.md#API`
  - 正: `memx_spec_v3/docs/design.md#3. API設計`
- 誤: `EVALUATION.md#pass-fail`
  - 正: `docs/birdseye/caps/EVALUATION.md.json#pass_fail_rules`
- 誤: `RUNBOOK.md#rollback`
  - 正: `docs/birdseye/caps/RUNBOOK.md.json#rollback`
- 誤: `contracts.md#4. CLI JSON`
  - 正: `memx_spec_v3/docs/CONTRACTS.md#4. CLI JSON`

## 5. 参照文字列の棚卸し（正規化対象一覧）

以下の文書を対象に、参照文字列の正規化対象を棚卸しした。

- `docs/TASKS.md`
  - `requirements.md` / `design.md` / `interfaces.md` / `traceability.md` / `EVALUATION.md` / `RUNBOOK.md` / `docs/birdseye/index.json` の入力参照名が記載されているため、本仕様の正規パスマッピング適用対象。
- `orchestration/memx-design-docs-authoring.md`
  - 入力成果物として `requirements.md` / `traceability.md` / `design.md` / `interfaces.md` / `EVALUATION.md` / `RUNBOOK.md` / `docs/birdseye/index.json` が繰り返し記載されているため、本仕様での解決対象。
- `memx_spec_v3/docs/design-review-spec.md`
  - `requirements.md` / `traceability.md` / `EVALUATION.md` などの参照文字列が含まれるため、レビュー記録の `Source` 記述時は本仕様の正規パスへ正規化する。

本一覧で特定した参照名のうち、`contracts.md` / `CONTRACTS.md` は 1-1 節の固定ルールを優先適用する。

## 6. 完了条件（EVALUATION/RUNBOOK 参照の収束）

- 完了条件として、`memx_spec_v3/docs` 配下の設計仕様（`design-*.md`）から `memx_spec_v3/docs/EVALUATION.md` / `memx_spec_v3/docs/RUNBOOK.md` 参照が 0 件であることを必須とする。
