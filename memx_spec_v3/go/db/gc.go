package db

import "context"

// GCOptions は GC 実行時のオプション。
type GCOptions struct {
	DryRun bool
	// 将来: MaxDuration, MaxNotes などを追加可能
}

// GCShort は short.db に対する GC（Observer / Reflector を含む）を実行する。
// Phase0: トリガ判定
// Phase1: Observer（クラスタリング & 관찰ノ트生成）
// Phase2: Reflector（Knowledge ページ更新）
// Phase3: Short -> Archive 退避
// 実際の LLM 呼び出しには Conn.Mini / Conn.Reflect を利用する想定。
func (c *Conn) GCShort(ctx context.Context, opt GCOptions) error {
	// TODO: 実装
	return nil
}

// RunObserver は Observer フェーズのみを実行する。
func (c *Conn) RunObserver(ctx context.Context, opt GCOptions) error {
	// TODO: 実装
	return nil
}

// RunReflector は Reflector フェーズのみを実行する。
func (c *Conn) RunReflector(ctx context.Context, opt GCOptions) error {
	// TODO: 実装
	return nil
}
