package service

import (
	"testing"

	"memx/db"
)

func TestNormalizeRecallStores(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []db.StoreKind
	}{
		{
			name:     "empty defaults to short",
			input:    []string{},
			expected: []db.StoreKind{db.StoreShort},
		},
		{
			name:     "nil defaults to short",
			input:    nil,
			expected: []db.StoreKind{db.StoreShort},
		},
		{
			name:     "single valid store",
			input:    []string{"journal"},
			expected: []db.StoreKind{db.StoreJournal},
		},
		{
			name:     "multiple valid stores",
			input:    []string{"short", "journal", "knowledge"},
			expected: []db.StoreKind{db.StoreShort, db.StoreJournal, db.StoreKnowledge},
		},
		{
			name:     "invalid store filtered",
			input:    []string{"short", "invalid", "knowledge"},
			expected: []db.StoreKind{db.StoreShort, db.StoreKnowledge},
		},
		{
			name:     "case insensitive",
			input:    []string{"SHORT", "Journal", "KNOWLEDGE"},
			expected: []db.StoreKind{db.StoreShort, db.StoreJournal, db.StoreKnowledge},
		},
		{
			name:     "whitespace trimmed",
			input:    []string{" short ", " journal "},
			expected: []db.StoreKind{db.StoreShort, db.StoreJournal},
		},
		{
			name:     "all invalid defaults to short",
			input:    []string{"invalid", "unknown"},
			expected: []db.StoreKind{db.StoreShort},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeRecallStores(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("normalizeRecallStores() = %v, want %v", result, tt.expected)
				return
			}
			for i, s := range result {
				if s != tt.expected[i] {
					t.Errorf("normalizeRecallStores()[%d] = %v, want %v", i, s, tt.expected[i])
				}
			}
		})
	}
}