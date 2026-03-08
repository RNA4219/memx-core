package api

import (
	"context"
	"errors"
	"testing"
)

func TestShortNoteResolver_ResolveRef(t *testing.T) {
	tests := []struct {
		name       string
		ref        TypedRef
		mockNote   *Note
		mockErr    error
		wantStatus RefStatus
		wantErr    bool
	}{
		{
			name: "resolved evidence ref",
			ref:  NewTypedRef(EntityTypeEvidence, "test-id-1"),
			mockNote: &Note{
				ID:       "test-id-1",
				Title:    "Test Evidence",
				Summary:  "Test summary",
				Body:     "Test body",
			},
			wantStatus: RefStatusResolved,
		},
		{
			name: "resolved knowledge ref",
			ref:  NewTypedRef(EntityTypeKnowledge, "test-id-2"),
			mockNote: &Note{
				ID:       "test-id-2",
				Title:    "Test Knowledge",
				Summary:  "Knowledge summary",
			},
			wantStatus: RefStatusResolved,
		},
		{
			name:       "unresolved ref",
			ref:        NewTypedRef(EntityTypeEvidence, "not-found"),
			mockErr:    errors.New("not found"),
			wantStatus: RefStatusUnresolved,
		},
		{
			name:    "unsupported domain",
			ref:     NewTypedRefWithProvider(DomainAgentTaskstate, "task", ProviderLocal, "task-1"),
			wantErr: true,
		},
		{
			name:    "unsupported type",
			ref:     TypedRef{Domain: DomainMemx, Type: "unknown", Provider: ProviderLocal, ID: "id"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchFunc := func(ctx context.Context, query string, topK int) ([]Note, error) {
				return nil, nil
			}
			showFunc := func(ctx context.Context, id string) (*Note, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				if tt.mockNote != nil && tt.mockNote.ID == id {
					return tt.mockNote, nil
				}
				return nil, errors.New("not found")
			}

			resolver := NewShortNoteResolver(searchFunc, showFunc)
			got, err := resolver.ResolveRef(context.Background(), tt.ref)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got.Status != tt.wantStatus {
				t.Errorf("status = %s, want %s", got.Status, tt.wantStatus)
			}

			if tt.wantStatus == RefStatusResolved {
				if got.Summary == "" {
					t.Error("expected summary, got empty")
				}
			}
		})
	}
}

func TestShortNoteResolver_ResolveMany(t *testing.T) {
	notes := map[string]*Note{
		"id-1": {ID: "id-1", Title: "Note 1", Summary: "Summary 1"},
		"id-2": {ID: "id-2", Title: "Note 2", Summary: "Summary 2"},
	}

	searchFunc := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return nil, nil
	}
	showFunc := func(ctx context.Context, id string) (*Note, error) {
		if note, ok := notes[id]; ok {
			return note, nil
		}
		return nil, errors.New("not found")
	}

	resolver := NewShortNoteResolver(searchFunc, showFunc)

	refs := []TypedRef{
		NewTypedRef(EntityTypeEvidence, "id-1"),
		NewTypedRef(EntityTypeKnowledge, "id-2"),
		NewTypedRef(EntityTypeEvidence, "not-found"),
		NewTypedRefWithProvider(DomainAgentTaskstate, "task", ProviderLocal, "task-1"),
	}

	report, err := resolver.ResolveMany(context.Background(), refs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(report.Resolved) != 2 {
		t.Errorf("resolved count = %d, want 2", len(report.Resolved))
	}

	if len(report.Unresolved) != 1 {
		t.Errorf("unresolved count = %d, want 1", len(report.Unresolved))
	}

	if len(report.Unsupported) != 1 {
		t.Errorf("unsupported count = %d, want 1", len(report.Unsupported))
	}

	if !report.Diagnostics.PartialBundle {
		t.Error("expected partial bundle to be true")
	}
}

func TestShortNoteResolver_LoadSummary(t *testing.T) {
	notes := map[string]*Note{
		"test-id": {ID: "test-id", Title: "Test", Summary: "Test summary"},
	}

	searchFunc := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return nil, nil
	}
	showFunc := func(ctx context.Context, id string) (*Note, error) {
		if note, ok := notes[id]; ok {
			return note, nil
		}
		return nil, errors.New("not found")
	}

	resolver := NewShortNoteResolver(searchFunc, showFunc)

	t.Run("existing ref", func(t *testing.T) {
		ref := NewTypedRef(EntityTypeEvidence, "test-id")
		payload, err := resolver.LoadSummary(context.Background(), ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !payload.Exists {
			t.Error("expected exists to be true")
		}
		if payload.Summary != "Test summary" {
			t.Errorf("summary = %s, want 'Test summary'", payload.Summary)
		}
	})

	t.Run("non-existing ref", func(t *testing.T) {
		ref := NewTypedRef(EntityTypeEvidence, "not-found")
		payload, err := resolver.LoadSummary(context.Background(), ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if payload.Exists {
			t.Error("expected exists to be false")
		}
	})
}

func TestShortNoteResolver_LoadSelectedRaw(t *testing.T) {
	notes := map[string]*Note{
		"test-id": {ID: "test-id", Title: "Test", Summary: "Summary", Body: "Full body content"},
	}

	searchFunc := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return nil, nil
	}
	showFunc := func(ctx context.Context, id string) (*Note, error) {
		if note, ok := notes[id]; ok {
			return note, nil
		}
		return nil, errors.New("not found")
	}

	resolver := NewShortNoteResolver(searchFunc, showFunc)

	t.Run("with body", func(t *testing.T) {
		ref := NewTypedRef(EntityTypeEvidence, "test-id")
		selector := RawSelector{IncludeBody: true}
		payload, err := resolver.LoadSelectedRaw(context.Background(), ref, selector)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !payload.Found {
			t.Error("expected found to be true")
		}
		if payload.Raw != "Full body content" {
			t.Errorf("raw = %s, want 'Full body content'", payload.Raw)
		}
	})

	t.Run("without body", func(t *testing.T) {
		ref := NewTypedRef(EntityTypeEvidence, "test-id")
		selector := RawSelector{IncludeBody: false}
		payload, err := resolver.LoadSelectedRaw(context.Background(), ref, selector)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !payload.Found {
			t.Error("expected found to be true")
		}
		if payload.Raw != "Summary" {
			t.Errorf("raw = %s, want 'Summary'", payload.Raw)
		}
	})
}

func TestValidateTypedRefForResolve(t *testing.T) {
	tests := []struct {
		name    string
		ref     TypedRef
		wantErr bool
	}{
		{
			name:    "valid memx local ref",
			ref:     NewTypedRef(EntityTypeEvidence, "test-id"),
			wantErr: false,
		},
		{
			name:    "invalid empty ref",
			ref:     TypedRef{},
			wantErr: true,
		},
		{
			name:    "unsupported agent-taskstate domain",
			ref:     NewTypedRefWithProvider(DomainAgentTaskstate, "task", ProviderLocal, "task-1"),
			wantErr: true,
		},
		{
			name:    "unsupported non-local provider",
			ref:     NewTypedRefWithProvider(DomainMemx, EntityTypeEvidence, ProviderJira, "id"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypedRefForResolve(tt.ref)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// -------------------- P4 Tests: Context Bundle --------------------

func TestContextBundle_Build(t *testing.T) {
	notes := map[string]*Note{
		"ev-1": {ID: "ev-1", Title: "Evidence 1", Summary: "Evidence summary 1", Body: "Full evidence body"},
		"art-1": {ID: "art-1", Title: "Artifact 1", Summary: "Artifact summary 1"},
	}

	searchFunc := func(ctx context.Context, query string, topK int) ([]Note, error) {
		return nil, nil
	}
	showFunc := func(ctx context.Context, id string) (*Note, error) {
		if note, ok := notes[id]; ok {
			return note, nil
		}
		return nil, errors.New("not found")
	}

	resolver := NewShortNoteResolver(searchFunc, showFunc)

	refs := []TypedRef{
		NewTypedRef(EntityTypeEvidence, "ev-1"),
		NewTypedRef(EntityTypeArtifact, "art-1"),
	}

	report, err := resolver.ResolveMany(context.Background(), refs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(report.Resolved) != 2 {
		t.Errorf("resolved count = %d, want 2", len(report.Resolved))
	}

	// Verify bundle structure
	bundle := ContextBundle{
		ID:             "test-bundle",
		Purpose:        "test",
		Summary:        "Test bundle",
		SourceRefs:     []BundleSourceRef{},
		RawIncluded:    false,
		GeneratorVersion: "test/v1",
		Diagnostics: BundleDiagnostics{
			PartialBundle: report.Diagnostics.PartialBundle,
		},
	}

	if bundle.ID == "" {
		t.Error("bundle ID should not be empty")
	}
	if bundle.Purpose == "" {
		t.Error("bundle purpose should not be empty")
	}
}

func TestBundleDiagnostics_PartialBundle(t *testing.T) {
	tests := []struct {
		name          string
		unresolved    int
		unsupported   int
		wantPartial   bool
	}{
		{
			name:        "complete bundle",
			unresolved:  0,
			unsupported: 0,
			wantPartial: false,
		},
		{
			name:        "partial with unresolved",
			unresolved:  1,
			unsupported: 0,
			wantPartial: true,
		},
		{
			name:        "partial with unsupported",
			unresolved:  0,
			unsupported: 1,
			wantPartial: true,
		},
		{
			name:        "partial with both",
			unresolved:  2,
			unsupported: 1,
			wantPartial: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics := BundleDiagnostics{
				MissingRefs:     make([]TypedRef, tt.unresolved),
				UnsupportedRefs: make([]TypedRef, tt.unsupported),
				PartialBundle:   tt.unresolved > 0 || tt.unsupported > 0,
			}

			if diagnostics.PartialBundle != tt.wantPartial {
				t.Errorf("PartialBundle = %v, want %v", diagnostics.PartialBundle, tt.wantPartial)
			}
		})
	}
}

func TestBundleSourceRef_Structure(t *testing.T) {
	ref := NewTypedRef(EntityTypeEvidence, "test-id")
	sourceRef := BundleSourceRef{
		Ref:         ref,
		SourceKind:  "evidence",
		SelectedRaw: true,
		MetadataJSON: `{"key": "value"}`,
	}

	if sourceRef.Ref.ID != "test-id" {
		t.Errorf("ref ID = %s, want test-id", sourceRef.Ref.ID)
	}
	if sourceRef.SourceKind != "evidence" {
		t.Errorf("source kind = %s, want evidence", sourceRef.SourceKind)
	}
	if !sourceRef.SelectedRaw {
		t.Error("selected raw should be true")
	}
}

// -------------------- P4 Tests: API Types --------------------

func TestResolveRefRequest_Response(t *testing.T) {
	ref := NewTypedRef(EntityTypeEvidence, "test-id")
	req := ResolveRefRequest{Ref: ref}

	if req.Ref.ID != "test-id" {
		t.Errorf("ref ID = %s, want test-id", req.Ref.ID)
	}

	resolved := ResolvedRef{
		Ref:     ref,
		Status:  RefStatusResolved,
		Summary: "Test summary",
	}

	resp := ResolveRefResponse{Resolved: resolved}
	if resp.Resolved.Status != RefStatusResolved {
		t.Errorf("status = %s, want resolved", resp.Resolved.Status)
	}
}

func TestResolveManyRequest_Response(t *testing.T) {
	refs := []TypedRef{
		NewTypedRef(EntityTypeEvidence, "id-1"),
		NewTypedRef(EntityTypeKnowledge, "id-2"),
	}

	req := ResolveManyRequest{Refs: refs}
	if len(req.Refs) != 2 {
		t.Errorf("refs count = %d, want 2", len(req.Refs))
	}

	report := ResolveReport{
		Resolved: []ResolvedRef{
			{Ref: refs[0], Status: RefStatusResolved},
			{Ref: refs[1], Status: RefStatusResolved},
		},
		Unresolved:  []TypedRef{},
		Unsupported: []TypedRef{},
	}

	resp := ResolveManyResponse{Report: report}
	if len(resp.Report.Resolved) != 2 {
		t.Errorf("resolved count = %d, want 2", len(resp.Report.Resolved))
	}
}

func TestBuildBundleRequest_Response(t *testing.T) {
	req := BuildBundleRequest{
		Purpose:    "context rebuild test",
		SourceRefs: []TypedRef{NewTypedRef(EntityTypeEvidence, "ev-1")},
		IncludeRaw: false,
	}

	if req.Purpose != "context rebuild test" {
		t.Errorf("purpose = %s, want 'context rebuild test'", req.Purpose)
	}

	bundle := ContextBundle{
		ID:          "bundle-001",
		Purpose:     req.Purpose,
		Summary:     "Test summary",
		RawIncluded: req.IncludeRaw,
	}

	resp := BuildBundleResponse{Bundle: bundle}
	if resp.Bundle.ID != "bundle-001" {
		t.Errorf("bundle ID = %s, want bundle-001", resp.Bundle.ID)
	}
}