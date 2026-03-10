package api

import (
	"testing"
)

func TestTypedRef_String(t *testing.T) {
	tests := []struct {
		ref      TypedRef
		expected string
	}{
		{TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HXXXXXXX"}, "memx:evidence:local:01HXXXXXXX"},
		{TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, Provider: ProviderLocal, ID: "01HYYYYYYY"}, "memx:knowledge:local:01HYYYYYYY"},
		{TypedRef{Domain: DomainMemx, Type: EntityTypeArtifact, Provider: ProviderLocal, ID: "01HZZZZZZZ"}, "memx:artifact:local:01HZZZZZZZ"},
		{TypedRef{Domain: DomainMemx, Type: EntityTypeLineage, Provider: ProviderLocal, ID: "01HAAAAAAA"}, "memx:lineage:local:01HAAAAAAA"},
		{TypedRef{Domain: DomainMemx, Type: EntityTypeEvidenceChunk, Provider: ProviderLocal, ID: "01HCHUNK"}, "memx:evidence_chunk:local:01HCHUNK"},
		{TypedRef{Domain: DomainAgentTaskstate, Type: "task", Provider: ProviderLocal, ID: "task_01J"}, "agent-taskstate:task:local:task_01J"},
		{TypedRef{Domain: DomainTracker, Type: "issue", Provider: ProviderJira, ID: "PROJ-123"}, "tracker:issue:jira:PROJ-123"},
		{TypedRef{Domain: DomainTracker, Type: "issue", Provider: ProviderGitHub, ID: "owner/repo#123"}, "tracker:issue:github:owner/repo#123"},
	}

	for _, tt := range tests {
		if got := tt.ref.String(); got != tt.expected {
			t.Errorf("TypedRef.String() = %s, want %s", got, tt.expected)
		}
	}
}

func TestParseTypedRef(t *testing.T) {
	tests := []struct {
		input    string
		expected TypedRef
		hasError bool
	}{
		// 4-segment format (canonical)
		{"memx:evidence:local:01HXXXXXXX", TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HXXXXXXX"}, false},
		{"memx:knowledge:local:01HYYYYYYY", TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, Provider: ProviderLocal, ID: "01HYYYYYYY"}, false},
		{"memx:artifact:local:01HZZZZZZZ", TypedRef{Domain: DomainMemx, Type: EntityTypeArtifact, Provider: ProviderLocal, ID: "01HZZZZZZZ"}, false},
		{"memx:lineage:local:01HAAAAAAA", TypedRef{Domain: DomainMemx, Type: EntityTypeLineage, Provider: ProviderLocal, ID: "01HAAAAAAA"}, false},
		{"agent-taskstate:task:local:task_01J", TypedRef{Domain: DomainAgentTaskstate, Type: "task", Provider: ProviderLocal, ID: "task_01J"}, false},
		{"tracker:issue:jira:PROJ-123", TypedRef{Domain: DomainTracker, Type: "issue", Provider: ProviderJira, ID: "PROJ-123"}, false},
		{"tracker:issue:github:owner/repo#123", TypedRef{Domain: DomainTracker, Type: "issue", Provider: ProviderGitHub, ID: "owner/repo#123"}, false},

		// 3-segment format (legacy, normalized to provider=local)
		{"memx:evidence:01HXXXXXXX", TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HXXXXXXX"}, false},
		{"memx:knowledge:01HYYYYYYY", TypedRef{Domain: DomainMemx, Type: EntityTypeKnowledge, Provider: ProviderLocal, ID: "01HYYYYYYY"}, false},
		{"memx:artifact:01HZZZZZZZ", TypedRef{Domain: DomainMemx, Type: EntityTypeArtifact, Provider: ProviderLocal, ID: "01HZZZZZZZ"}, false},
		{"memx:lineage:01HAAAAAAA", TypedRef{Domain: DomainMemx, Type: EntityTypeLineage, Provider: ProviderLocal, ID: "01HAAAAAAA"}, false},

		// Error cases
		{"", TypedRef{}, true},                          // empty
		{"invalid", TypedRef{}, true},                   // wrong format
		{"foo:evidence:local:01HXXX", TypedRef{}, true}, // wrong domain
		{"memx:unknown:local:01HXXX", TypedRef{}, true}, // unknown type for memx
		{"memx:evidence:local:", TypedRef{}, true},      // empty id
		{"agent-taskstate:task:01J", TypedRef{}, true},            // 3-segment only memx allowed
	}

	for _, tt := range tests {
		got, err := ParseTypedRef(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseTypedRef(%s) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseTypedRef(%s) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("ParseTypedRef(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		}
	}
}

func TestTypedRef_RoundTrip(t *testing.T) {
	refs := []TypedRef{
		{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HXXXXXXX"},
		{Domain: DomainMemx, Type: EntityTypeKnowledge, Provider: ProviderLocal, ID: "01HYYYYYYY"},
		{Domain: DomainMemx, Type: EntityTypeArtifact, Provider: ProviderLocal, ID: "01HZZZZZZZ"},
		{Domain: DomainMemx, Type: EntityTypeLineage, Provider: ProviderLocal, ID: "01HAAAAAAA"},
		{Domain: DomainAgentTaskstate, Type: "task", Provider: ProviderLocal, ID: "task_01J"},
		{Domain: DomainTracker, Type: "issue", Provider: ProviderJira, ID: "PROJ-123"},
	}

	for _, original := range refs {
		s := original.String()
		parsed, err := ParseTypedRef(s)
		if err != nil {
			t.Errorf("ParseTypedRef(%s) error: %v", s, err)
			continue
		}
		if parsed != original {
			t.Errorf("RoundTrip: got %v, want %v", parsed, original)
		}
	}
}

func TestTypedRef_ThreeToFourSegmentMigration(t *testing.T) {
	// 3-segment input should be normalized to 4-segment with provider=local
	input := "memx:evidence:01HXXX"
	parsed, err := ParseTypedRef(input)
	if err != nil {
		t.Fatalf("ParseTypedRef(%s) error: %v", input, err)
	}

	// Should be normalized
	expected := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HXXX"}
	if parsed != expected {
		t.Errorf("ParseTypedRef(%s) = %v, want %v", input, parsed, expected)
	}

	// Output should be 4-segment
	output := parsed.String()
	expectedOutput := "memx:evidence:local:01HXXX"
	if output != expectedOutput {
		t.Errorf("String() = %s, want %s", output, expectedOutput)
	}
}

func TestNewTypedRef(t *testing.T) {
	ref := NewTypedRef(EntityTypeEvidence, "01HTEST")
	if ref.Type != EntityTypeEvidence {
		t.Errorf("NewTypedRef type = %s, want %s", ref.Type, EntityTypeEvidence)
	}
	if ref.ID != "01HTEST" {
		t.Errorf("NewTypedRef id = %s, want 01HTEST", ref.ID)
	}
	if ref.Domain != DomainMemx {
		t.Errorf("NewTypedRef domain = %s, want %s", ref.Domain, DomainMemx)
	}
	if ref.Provider != ProviderLocal {
		t.Errorf("NewTypedRef provider = %s, want %s", ref.Provider, ProviderLocal)
	}
}

func TestNewTypedRefWithProvider(t *testing.T) {
	ref := NewTypedRefWithProvider(DomainTracker, "issue", ProviderJira, "PROJ-123")
	if ref.Domain != DomainTracker {
		t.Errorf("domain = %s, want %s", ref.Domain, DomainTracker)
	}
	if ref.Type != "issue" {
		t.Errorf("type = %s, want issue", ref.Type)
	}
	if ref.Provider != ProviderJira {
		t.Errorf("provider = %s, want %s", ref.Provider, ProviderJira)
	}
	if ref.ID != "PROJ-123" {
		t.Errorf("id = %s, want PROJ-123", ref.ID)
	}
}

func TestTypedRef_IsZero_IsValid(t *testing.T) {
	zero := TypedRef{}
	if !zero.IsZero() {
		t.Error("empty TypedRef should be zero")
	}
	if zero.IsValid() {
		t.Error("empty TypedRef should not be valid")
	}

	valid := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HTEST"}
	if valid.IsZero() {
		t.Error("non-empty TypedRef should not be zero")
	}
	if !valid.IsValid() {
		t.Error("non-empty TypedRef should be valid")
	}

	noID := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: ""}
	if noID.IsValid() {
		t.Error("TypedRef without ID should not be valid")
	}

	noType := TypedRef{Domain: DomainMemx, Type: "", Provider: ProviderLocal, ID: "01HTEST"}
	if noType.IsValid() {
		t.Error("TypedRef without Type should not be valid")
	}

	noProvider := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: "", ID: "01HTEST"}
	if noProvider.IsValid() {
		t.Error("TypedRef without Provider should not be valid")
	}

	noDomain := TypedRef{Domain: "", Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HTEST"}
	if noDomain.IsValid() {
		t.Error("TypedRef without Domain should not be valid")
	}
}

func TestTypedRef_MarshalUnmarshal(t *testing.T) {
	original := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HTEST"}

	data, err := original.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if string(data) != "memx:evidence:local:01HTEST" {
		t.Errorf("MarshalText = %s, want memx:evidence:local:01HTEST", string(data))
	}

	var parsed TypedRef
	if err := parsed.UnmarshalText(data); err != nil {
		t.Fatalf("UnmarshalText error: %v", err)
	}
	if parsed != original {
		t.Errorf("UnmarshalText = %v, want %v", parsed, original)
	}
}

func TestTypedRef_MarshalUnmarshal_LegacyInput(t *testing.T) {
	// Unmarshaling 3-segment should produce canonical 4-segment
	legacy := []byte("memx:evidence:01HTEST")

	var parsed TypedRef
	if err := parsed.UnmarshalText(legacy); err != nil {
		t.Fatalf("UnmarshalText error: %v", err)
	}

	// Should be normalized to 4-segment
	expected := TypedRef{Domain: DomainMemx, Type: EntityTypeEvidence, Provider: ProviderLocal, ID: "01HTEST"}
	if parsed != expected {
		t.Errorf("UnmarshalText = %v, want %v", parsed, expected)
	}

	// Re-marshaling should produce 4-segment
	data, err := parsed.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if string(data) != "memx:evidence:local:01HTEST" {
		t.Errorf("MarshalText = %s, want memx:evidence:local:01HTEST", string(data))
	}
}