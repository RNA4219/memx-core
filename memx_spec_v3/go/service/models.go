package service

// Note は notes テーブル相当のドメインモデル。
// API 層へそのまま露出しない（必要なら api.Note に変換する）。
//
// 文字列は RFC3339 の UTC を想定。
type Note struct {
	ID             string
	Title          string
	Summary        string
	Body           string
	CreatedAt      string
	UpdatedAt      string
	LastAccessedAt string
	AccessCount    int64

	SourceType  string
	Origin      string
	SourceTrust string
	Sensitivity string
}
