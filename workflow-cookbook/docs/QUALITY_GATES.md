# workflow-cookbook 品質ゲート参照

適用対象: 本書は `workflow-cookbook/` 配下の Task/実装変更に適用する。

- Python: `ruff check .` / `mypy .` / `pytest -q`
- 補助確認（必要に応じて）: `python tools/ci/check_governance_gate.py`

※ memx 本体（`memx_spec_v3/`）の品質ゲートは `../QUALITY_GATES.md` ではなく、リポジトリルートの `docs/QUALITY_GATES.md` を参照する。
