---
intent_id: INT-001
owner: your-handle
status: active   # draft|active|deprecated
last_reviewed_at: 2025-10-21
next_review_due: 2025-11-21
---

# Runbook

## Environments

- Local / CI / Prod の差分（キー名だけ列挙）

## Execute

- 準備 → 実行 → 確認（最短手順）
- 例）
  - 準備: データ投入 / キャッシュ初期化
  - 実行: コマンド/ジョブ名
  - 確認: 出力の存在・件数・整合

## Observability

- ログ/メトリクスの確認点、失敗時の兆候（[ADR-021: メトリクスと可観測性の統合](docs/ADR/ADR-021-metrics-observability.md) を参照）
- KPI の目的・閾値は [EVALUATION.md#KPIs](EVALUATION.md#kpis) を参照し、収集手順と解釈を同期する。
- インシデント発生時は docs/IN-YYYYMMDD-XXX.md に記録し、最新サンプル（[IN-20250115-001](docs/IN-20250115-001.md)）を参照して検知し、ログ・メトリクスの抜粋を添付
- QA メトリクス収集と確認（`tools/perf/` 共通テンプレート準拠）
  - `python -m tools.perf.collect_metrics --suite qa --metrics-url <Prometheus URL> --log-path <StructuredLogger 等の運用ログ>`
    を実行する。`--suite qa` は `.ga/qa-metrics.json` への書き出しを既定とし、Prometheus
    （`workflow_review_latency_*` は `review_latency` に正規化され、旧 `legacy_review_latency_*`
    も互換処理で継続利用可能。`task_seed_cycle_time_*`／`birdseye_refresh_delay_*` も同様に平均化）と
    運用ログ（`checklist_compliance_rate` の比率計算も含む）から統合メトリクスを取得する。
    出力先を変更したい場合は `--output <JSON パス>` を追加指定する。
    `semantic_retention` を取得するには `tools/perf/context_trimmer.trim_messages` へ
    `semantic_options`（例: `{"embedder": <callable>}`）を渡せるよう、ログエミッタ側で埋め込み関数を設定しておく。
    指標と検証手順の詳細は [docs/addenda/D_Context_Trimming.md](docs/addenda/D_Context_Trimming.md) を参照する。
    埋め込み関数はテキストを `Sequence[float]` へ変換できる必要があり、トリミング後の意味保持率は
    このベクトル間のコサイン類似度として集計される。
    `--metrics-url` または `--log-path` のどちらか片方しか利用できない場合は、利用可能な入力のみ指定する。
  - 対話 UI や運用ツールからメトリクスを出力する場合は `tools.perf.structured_logger.StructuredLogger`
    を利用する。例: `from tools.perf.structured_logger import StructuredLogger` →
    `StructuredLogger(name="ops-ui", path="~/.logs/ops/metrics.log")`
    `.inference(metrics={"semantic_retention": 0.9})`。
    こうして生成された JSON ログ行は `collect_metrics --log-path ~/.logs/ops/metrics.log` で取り込まれ、
    `metrics` キー配下の辞書がそのまま運用ログ由来メトリクスとして集計される。
  - FastAPI などの Web サービスに組み込む場合は `tools.perf.metrics_registry.MetricsRegistry` を共有シングルトン
    として初期化し、トリミング完了時に `observe_trim` を呼び出す。`compress_ratio=` を直接指定する新 API と、
    既存の `original_tokens=` / `trimmed_tokens=` を渡す後方互換 API のどちらでも動作し、`semantic_retention`
    の有無も任意。`@app.get("/metrics")` エンドポイントで `return PlainTextResponse(registry.export_prometheus())`
    を返すと Prometheus が取得可能となる。収集 CLI は公開 API として `compress_ratio` / `semantic_retention`
    を参照しつつ、Prometheus 上では `trim_compress_ratio_*` / `trim_semantic_retention_*` を優先的に解釈する。
    再オープン率は `workflow_reopen_rate_*` → `docops_reopen_rate` → `reopen_rate` を順に確認し、スペック充足率は
    `workflow_spec_completeness_*` → `spec_completeness_*` → `spec_completeness` をフォールバックとして参照する。
  - 実行後に `.ga/qa-metrics.json` がリポジトリルート配下へ生成されていることを確認する。生成されない場合は
    `--output` に明示したパスと標準出力を突き合わせ、異常がないか確認する。
  - `python - <<'PY'` を実行し、以下を評価して各メトリクスの値を抽出する:

    ```python
    import json

    with open('.ga/qa-metrics.json', encoding='utf-8') as fh:
        data = json.load(fh)

    keys = (
        'checklist_compliance_rate',
        'compress_ratio',
        'semantic_retention',
        'task_seed_cycle_time_minutes',
        'birdseye_refresh_delay_minutes',
        'review_latency',
        'reopen_rate',
        'spec_completeness',
    )
    print({k: data[k] for k in keys})
    ```

    閾値は最新サンプルと突き合わせ、外れた場合は直近成功値との差分と再現条件を記録して共有する。
  - FastAPI 等へ常駐組み込みする際は `tools.perf.metrics_registry.MetricsRegistry` を介し、トリミング結果を逐次記録する:

      ```python
    from fastapi import FastAPI, Response

    from tools.perf.metrics_registry import MetricsRegistry

    registry = MetricsRegistry(default_labels={"service": "workflow"})
    app = FastAPI()

    @app.post("/trim")
    async def record_trim(payload: dict[str, float]) -> dict[str, str]:
        registry.observe_trim(
            compress_ratio=payload["compress_ratio"],
            semantic_retention=payload["semantic_retention"],
            labels={"model": payload.get("model", "unknown")},
        )
        # 旧 API を利用する場合の例（compress_ratio が未計算なときなど）:
        # registry.observe_trim(
        #     original_tokens=payload["original_tokens"],
        #     trimmed_tokens=payload["trimmed_tokens"],
        #     semantic_retention=payload.get("semantic_retention"),
        #     labels={"model": payload.get("model", "unknown")},
        # )
        return {"status": "ok"}

    @app.get("/metrics")
    async def metrics() -> Response:
        return Response(registry.export_prometheus(), media_type="text/plain")
      ```

  - `snapshot()` で `{"trim_compress_ratio": [{"labels": {...}, "count": 2, ...}]}` 形式の統計を確認できる。
    Prometheus エクスポートでは `trim_compress_ratio_{count,sum,avg,min,max}` および
    `trim_semantic_retention_{count,sum,avg,min,max}` を同一ラベル集合ごとに出力する。
    例:

    ```text
    # HELP trim_compress_ratio_count Compression ratio observed after trimming. (count).
    trim_compress_ratio_count{model="gpt-5",service="workflow"} 2
    trim_compress_ratio_avg{model="gpt-5",service="workflow"} 0.45
    ```

  - 公開メトリクス名: `trim_compress_ratio` / `trim_semantic_retention`
    （各 `_count`、`_sum`、`_avg`、`_min`、`_max` を同時出力）に加え、後方互換 Gauge
    `compress_ratio` / `semantic_retention` も平均値として公開する。

- 失敗兆候と一次対応
  - `.ga/qa-metrics.json` が生成されない / 壊れている: `python -m tools.perf.collect_metrics --help` でオプションを再確認し、一時ファイルやログ出力設定を洗い直してから再実行。
  - `checklist_compliance_rate` が 95% を下回る: 実行時のチェックリスト完了ログを抽出し、どの項目が未完了かを Birdseye や Git 履歴で確認する。改善作業が必要な場合は Task Seed を追加投入する。
  - `task_seed_cycle_time_minutes` が 1440 分を超過: 受付から着手までの待機要因（担当者アサイン、依頼内容不備など）を振り返り、対応 SLA を再共有する。
  - `birdseye_refresh_delay_minutes` が 60 分を超過: Birdseye 更新ジョブの実行ログとスケジューラ状態を確認し、必要に応じて手動更新を実施。

## HUB 更新トリガー時の最短運用手順

- 再生成: `python tools/codemap/update.py --targets docs/birdseye/index.json,docs/birdseye/hot.json --emit index+caps` を実行し、`index.json` と `hot.json` を同時更新する。
- 差分確認: `git diff -- docs/birdseye/index.json docs/birdseye/hot.json` で想定外差分の有無を確認する。
- 鮮度確認: `docs/birdseye/index.json.generated_at` と `docs/birdseye/hot.json.generated_at` の更新、ならびに `docs/birdseye/hot.json.refresh_command` と `docs/BIRDSEYE.md` 記載コマンドの一致をレビューゲートとして確認する。

## Outbound Request Approval

- 申請項目: [アウトバウンド通信申請テンプレート](tickets/outbound-request.md#申請テンプレート) を利用し、依頼者・宛先ドメイン/ポート・用途・送信データ分類・想定期間・
  フォールバック手順・再試行上限を記入し、最新の `network/allowlist.yaml` 差分を添付する。
- 承認者: 当番 SRE（一次）とセキュリティ/プライバシー担当（`docs/addenda/G_Security_Privacy.md#4-通信制御とツール実行`）が双方承認して初めて通信開始可とする。
- 記録: 承認完了後に `audit/outbound-requests.log` へ記録し、関連チケットへ決裁ログと実行結果リンクを残す。失敗時はインシデントテンプレ（`docs/IN-YYYYMMDD-XXX.md`）へ転記する。
- 再試行条件: 承認済み通信で 5xx / Timeout が発生した場合のみ、指数バックオフ（初回 2 分、最大 3 回）で自動再試行を許可する。
  4xx や `network/allowlist.yaml` 未反映による失敗は再申請を行い、承認が完了するまで再試行を禁止する。

## Confirm

- Execute 結果を主要メトリクス・アウトプットと突き合わせ、`CHECKLISTS.md` の [Hygiene](CHECKLISTS.md#hygiene) で整合性と未完了項目を再確認
- インシデント記録を [docs/INCIDENT_TEMPLATE.md](docs/INCIDENT_TEMPLATE.md) に沿って初動報告→確定記録まで更新し、関連 PR / チケットへリンクを共有
- `Observability` で検知したアラート・兆候の解消を運用チャネルへ報告し、残るフォローアップを RUNBOOK / docs/IN-YYYYMMDD-XXX.md に追記

## Rollback / Retry

- どこまで戻すか、再実行条件
- インシデントサマリを更新後、該当PRの説明欄と本RUNBOOKの該当セクションにリンクを追加する
