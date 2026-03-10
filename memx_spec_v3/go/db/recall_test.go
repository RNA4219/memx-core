package db

import (
	"context"
	"testing"
)

func TestRecallQuery_Validation(t *testing.T) {
	tests := []struct {
		name    string
		q       RecallQuery
		wantErr bool
	}{
		{
			name: "valid query",
			q:    RecallQuery{Text: "test query", Stores: []StoreKind{StoreShort}},
			wantErr: false,
		},
		{
			name: "empty query",
			q:    RecallQuery{Text: "", Stores: []StoreKind{StoreShort}},
			wantErr: true,
		},
		{
			name: "whitespace only query",
			q:    RecallQuery{Text: "   ", Stores: []StoreKind{StoreShort}},
			wantErr: true,
		},
		{
			name: "query too long",
			q:    RecallQuery{Text: string(make([]byte, 1001)), Stores: []StoreKind{StoreShort}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecallQuery(tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRecallQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeStores(t *testing.T) {
	tests := []struct {
		name     string
		stores   []StoreKind
		expected []StoreKind
	}{
		{
			name:     "empty stores defaults to short",
			stores:   []StoreKind{},
			expected: []StoreKind{StoreShort},
		},
		{
			name:     "nil stores defaults to short",
			stores:   nil,
			expected: []StoreKind{StoreShort},
		},
		{
			name:     "single valid store",
			stores:   []StoreKind{StoreJournal},
			expected: []StoreKind{StoreJournal},
		},
		{
			name:     "multiple valid stores",
			stores:   []StoreKind{StoreShort, StoreJournal, StoreKnowledge},
			expected: []StoreKind{StoreShort, StoreJournal, StoreKnowledge},
		},
		{
			name:     "archive is filtered out",
			stores:   []StoreKind{StoreShort, StoreArchive},
			expected: []StoreKind{StoreShort},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeStores(tt.stores)
			if len(result) != len(tt.expected) {
				t.Errorf("normalizeStores() = %v, want %v", result, tt.expected)
				return
			}
			for i, s := range result {
				if s != tt.expected[i] {
					t.Errorf("normalizeStores()[%d] = %v, want %v", i, s, tt.expected[i])
				}
			}
		})
	}
}

func TestEmbeddingToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected string
	}{
		{
			name:     "empty embedding",
			input:    []float64{},
			expected: "[]",
		},
		{
			name:     "single value",
			input:    []float64{1.0},
			expected: "[1.000000]",
		},
		{
			name:     "multiple values",
			input:    []float64{1.0, 2.5, 0.123456},
			expected: "[1.000000,2.500000,0.123456]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := embeddingToJSON(tt.input)
			if result != tt.expected {
				t.Errorf("embeddingToJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortByScore(t *testing.T) {
	notes := []RecallNote{
		{ID: "1", Score: 0.5},
		{ID: "2", Score: 0.9},
		{ID: "3", Score: 0.3},
		{ID: "4", Score: 0.7},
	}

	sortByScore(notes)

	expected := []float64{0.9, 0.7, 0.5, 0.3}
	for i, n := range notes {
		if n.Score != expected[i] {
			t.Errorf("sortByScore()[%d].Score = %v, want %v", i, n.Score, expected[i])
		}
	}
}

func TestFloat32To64(t *testing.T) {
	input := []float32{1.0, 2.5, 3.14}
	result := float32To64(input)

	if len(result) != len(input) {
		t.Errorf("float32To64() length = %d, want %d", len(result), len(input))
		return
	}

	for i, v := range result {
		if float32(v) != input[i] {
			t.Errorf("float32To64()[%d] = %v, want %v", i, v, input[i])
		}
	}
}

// MockEmbeddingClient はテスト用のモック埋め込みクライアント。
type MockEmbeddingClient struct {
	Embeddings [][]float32
	Err        error
}

func (m *MockEmbeddingClient) EmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if len(m.Embeddings) > 0 {
		return m.Embeddings, nil
	}
	// デフォルトのダミー埋め込みを返す
	result := make([][]float32, len(texts))
	for i := range texts {
		result[i] = make([]float32, 384)
		for j := range result[i] {
			result[i][j] = 0.1
		}
	}
	return result, nil
}