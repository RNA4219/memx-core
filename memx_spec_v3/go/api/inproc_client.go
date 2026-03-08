package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"memx/service"
)

type InProcClient struct {
	Svc      *service.Service
	resolver *MultiStoreResolver // P4: 遅延初期化
}

func NewInProcClient(svc *service.Service) *InProcClient {
	return &InProcClient{Svc: svc}
}

// -------------------- Helpers --------------------

func fromServiceNote(n service.Note) Note {
	return Note{
		ID:             n.ID,
		Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
		Title:          n.Title,
		Summary:        n.Summary,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
		LastAccessedAt: n.LastAccessedAt,
		AccessCount:    n.AccessCount,
		SourceType:     n.SourceType,
		Origin:         n.Origin,
		SourceTrust:    n.SourceTrust,
		Sensitivity:    n.Sensitivity,
	}
}

func fromServiceJournalNote(n service.JournalNote) JournalNote {
	return JournalNote{
		NoteBase: NoteBase{
			ID:             n.ID,
			Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
			Title:          n.Title,
			Summary:        n.Summary,
			Body:           n.Body,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:      n.UpdatedAt,
			LastAccessedAt: n.LastAccessedAt,
			AccessCount:    n.AccessCount,
			SourceType:     n.SourceType,
			Origin:         n.Origin,
			SourceTrust:    n.SourceTrust,
			Sensitivity:    n.Sensitivity,
		},
		WorkingScope: n.WorkingScope,
		IsPinned:     n.IsPinned,
	}
}

func fromServiceKnowledgeNote(n service.KnowledgeNote) KnowledgeNote {
	return KnowledgeNote{
		NoteBase: NoteBase{
			ID:             n.ID,
			Ref:            NewTypedRef(EntityTypeKnowledge, n.ID),
			Title:          n.Title,
			Summary:        n.Summary,
			Body:           n.Body,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:      n.UpdatedAt,
			LastAccessedAt: n.LastAccessedAt,
			AccessCount:    n.AccessCount,
			SourceType:     n.SourceType,
			Origin:         n.Origin,
			SourceTrust:    n.SourceTrust,
			Sensitivity:    n.Sensitivity,
		},
		WorkingScope: n.WorkingScope,
		IsPinned:     n.IsPinned,
	}
}

func fromServiceArchiveNote(n service.ArchiveNote) ArchiveNote {
	return ArchiveNote{
		ID:             n.ID,
		Ref:            NewTypedRef(EntityTypeEvidence, n.ID),
		Title:          n.Title,
		Summary:        n.Summary,
		Body:           n.Body,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
		LastAccessedAt: n.LastAccessedAt,
		AccessCount:    n.AccessCount,
		SourceType:     n.SourceType,
		Origin:         n.Origin,
		SourceTrust:    n.SourceTrust,
		Sensitivity:    n.Sensitivity,
	}
}

// -------------------- Short Store --------------------

func (c *InProcClient) NotesIngest(ctx context.Context, req NotesIngestRequest) (NotesIngestResponse, *Error) {
	n, err := c.Svc.IngestShort(ctx, service.IngestNoteRequest{
		Title:       req.Title,
		Body:        req.Body,
		Summary:     req.Summary,
		SourceType:  req.SourceType,
		Origin:      req.Origin,
		SourceTrust: req.SourceTrust,
		Sensitivity: req.Sensitivity,
		Tags:        req.Tags,
		NoLLM:       req.NoLLM,
	})
	if err != nil {
		return NotesIngestResponse{}, mapError(err)
	}
	return NotesIngestResponse{Note: fromServiceNote(n)}, nil
}

func (c *InProcClient) NotesSearch(ctx context.Context, req NotesSearchRequest) (NotesSearchResponse, *Error) {
	ns, err := c.Svc.SearchShort(ctx, req.Query, req.TopK)
	if err != nil {
		return NotesSearchResponse{}, mapError(err)
	}
	return NotesSearchResponse{Notes: mapNotes(ns)}, nil
}

func (c *InProcClient) NotesGet(ctx context.Context, id string) (Note, *Error) {
	n, err := c.Svc.GetShort(ctx, id)
	if err != nil {
		return Note{}, mapError(err)
	}
	return fromServiceNote(n), nil
}

func (c *InProcClient) GCRun(ctx context.Context, req GCRunRequest) (GCRunResponse, *Error) {
	result, err := c.Svc.GCShort(ctx, service.GCRequest{
		Target:  req.Target,
		DryRun:  req.Options.DryRun,
		Enabled: true,
	})
	if err != nil {
		return GCRunResponse{}, mapError(err)
	}
	if result.DryRun && result.DryRunResult != nil {
		return GCRunResponse{Status: result.DryRunResult.ToJSON()}, nil
	}
	return GCRunResponse{Status: result.Status}, nil
}

func (c *InProcClient) Summarize(ctx context.Context, id string) (SummarizeResponse, *Error) {
	n, err := c.Svc.SummarizeNote(ctx, id)
	if err != nil {
		return SummarizeResponse{}, mapError(err)
	}
	return SummarizeResponse{Note: fromServiceNote(n)}, nil
}

func (c *InProcClient) SummarizeBatch(ctx context.Context, req SummarizeBatchRequest) (SummarizeBatchResponse, *Error) {
	result, err := c.Svc.SummarizeNotes(ctx, req.IDs)
	if err != nil {
		return SummarizeBatchResponse{}, mapError(err)
	}
	return SummarizeBatchResponse{Summary: result.Summary, NoteCount: result.NoteCount}, nil
}

// -------------------- Journal Store --------------------

func (c *InProcClient) JournalIngest(ctx context.Context, req JournalIngestRequest) (JournalIngestResponse, *Error) {
	n, err := c.Svc.IngestJournal(ctx, service.IngestJournalRequest{
		Title:        req.Title,
		Body:         req.Body,
		Summary:      req.Summary,
		SourceType:   req.SourceType,
		Origin:       req.Origin,
		SourceTrust:  req.SourceTrust,
		Sensitivity:  req.Sensitivity,
		Tags:         req.Tags,
		WorkingScope: req.WorkingScope,
		IsPinned:     req.IsPinned,
		NoLLM:        req.NoLLM,
	})
	if err != nil {
		return JournalIngestResponse{}, mapError(err)
	}
	return JournalIngestResponse{Note: fromServiceJournalNote(n)}, nil
}

func (c *InProcClient) JournalSearch(ctx context.Context, req JournalSearchRequest) (JournalSearchResponse, *Error) {
	ns, err := c.Svc.SearchJournal(ctx, req.Query, req.TopK)
	if err != nil {
		return JournalSearchResponse{}, mapError(err)
	}
	return JournalSearchResponse{Notes: mapJournalNotes(ns)}, nil
}

func (c *InProcClient) JournalGet(ctx context.Context, id string) (JournalNote, *Error) {
	n, err := c.Svc.GetJournal(ctx, id)
	if err != nil {
		return JournalNote{}, mapError(err)
	}
	return fromServiceJournalNote(n), nil
}

func (c *InProcClient) JournalListByScope(ctx context.Context, req JournalListByScopeRequest) (JournalListByScopeResponse, *Error) {
	ns, err := c.Svc.ListJournalByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return JournalListByScopeResponse{}, mapError(err)
	}
	return JournalListByScopeResponse{Notes: mapJournalNotes(ns)}, nil
}

// -------------------- Knowledge Store --------------------

func (c *InProcClient) KnowledgeIngest(ctx context.Context, req KnowledgeIngestRequest) (KnowledgeIngestResponse, *Error) {
	n, err := c.Svc.IngestKnowledge(ctx, service.IngestKnowledgeRequest{
		Title:        req.Title,
		Body:         req.Body,
		Summary:      req.Summary,
		SourceType:   req.SourceType,
		Origin:       req.Origin,
		SourceTrust:  req.SourceTrust,
		Sensitivity:  req.Sensitivity,
		Tags:         req.Tags,
		WorkingScope: req.WorkingScope,
		IsPinned:     req.IsPinned,
		NoLLM:        req.NoLLM,
	})
	if err != nil {
		return KnowledgeIngestResponse{}, mapError(err)
	}
	return KnowledgeIngestResponse{Note: fromServiceKnowledgeNote(n)}, nil
}

func (c *InProcClient) KnowledgeSearch(ctx context.Context, req KnowledgeSearchRequest) (KnowledgeSearchResponse, *Error) {
	ns, err := c.Svc.SearchKnowledge(ctx, req.Query, req.TopK)
	if err != nil {
		return KnowledgeSearchResponse{}, mapError(err)
	}
	return KnowledgeSearchResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgeGet(ctx context.Context, id string) (KnowledgeNote, *Error) {
	n, err := c.Svc.GetKnowledge(ctx, id)
	if err != nil {
		return KnowledgeNote{}, mapError(err)
	}
	return fromServiceKnowledgeNote(n), nil
}

func (c *InProcClient) KnowledgeListByScope(ctx context.Context, req KnowledgeListByScopeRequest) (KnowledgeListByScopeResponse, *Error) {
	ns, err := c.Svc.ListKnowledgeByScope(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return KnowledgeListByScopeResponse{}, mapError(err)
	}
	return KnowledgeListByScopeResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgeListPinned(ctx context.Context, req KnowledgeListPinnedRequest) (KnowledgeListPinnedResponse, *Error) {
	ns, err := c.Svc.ListPinnedKnowledge(ctx, req.WorkingScope, req.Limit)
	if err != nil {
		return KnowledgeListPinnedResponse{}, mapError(err)
	}
	return KnowledgeListPinnedResponse{Notes: mapKnowledgeNotes(ns)}, nil
}

func (c *InProcClient) KnowledgePin(ctx context.Context, id string) (PinResponse, *Error) {
	if err := c.Svc.PinKnowledge(ctx, id); err != nil {
		return PinResponse{}, mapError(err)
	}
	return PinResponse{Success: true}, nil
}

func (c *InProcClient) KnowledgeUnpin(ctx context.Context, id string) (UnpinResponse, *Error) {
	if err := c.Svc.UnpinKnowledge(ctx, id); err != nil {
		return UnpinResponse{}, mapError(err)
	}
	return UnpinResponse{Success: true}, nil
}

// -------------------- Archive Store --------------------

func (c *InProcClient) ArchiveGet(ctx context.Context, id string) (ArchiveNote, *Error) {
	n, err := c.Svc.GetArchive(ctx, id)
	if err != nil {
		return ArchiveNote{}, mapError(err)
	}
	return fromServiceArchiveNote(n), nil
}

func (c *InProcClient) ArchiveList(ctx context.Context, req ArchiveListRequest) (ArchiveListResponse, *Error) {
	ns, err := c.Svc.ListArchive(ctx, req.Limit)
	if err != nil {
		return ArchiveListResponse{}, mapError(err)
	}
	return ArchiveListResponse{Notes: mapArchiveNotes(ns)}, nil
}

func (c *InProcClient) ArchiveRestore(ctx context.Context, id string) (ArchiveRestoreResponse, *Error) {
	n, err := c.Svc.RestoreFromArchive(ctx, id)
	if err != nil {
		return ArchiveRestoreResponse{}, mapError(err)
	}
	return ArchiveRestoreResponse{Note: fromServiceNote(n)}, nil
}

// -------------------- Resolver API (P4) --------------------

// ResolveRef は単一の typed_ref を解決する。
func (c *InProcClient) ResolveRef(ctx context.Context, req ResolveRefRequest) (ResolveRefResponse, *Error) {
	resolver := c.getResolver()
	if resolver == nil {
		return ResolveRefResponse{}, &Error{Code: CodeInternal, Message: "resolver not initialized"}
	}

	resolved, err := resolver.ResolveRef(ctx, req.Ref)
	if err != nil {
		return ResolveRefResponse{}, mapResolverError(err)
	}

	return ResolveRefResponse{Resolved: resolved}, nil
}

// ResolveMany は複数の typed_ref を一括解決する。
func (c *InProcClient) ResolveMany(ctx context.Context, req ResolveManyRequest) (ResolveManyResponse, *Error) {
	resolver := c.getResolver()
	if resolver == nil {
		return ResolveManyResponse{}, &Error{Code: CodeInternal, Message: "resolver not initialized"}
	}

	report, err := resolver.ResolveMany(ctx, req.Refs)
	if err != nil {
		return ResolveManyResponse{}, mapResolverError(err)
	}

	return ResolveManyResponse{Report: *report}, nil
}

// LoadSummary は要約を取得する（summary-first retrieval）。
func (c *InProcClient) LoadSummary(ctx context.Context, req LoadSummaryRequest) (LoadSummaryResponse, *Error) {
	resolver := c.getResolver()
	if resolver == nil {
		return LoadSummaryResponse{}, &Error{Code: CodeInternal, Message: "resolver not initialized"}
	}

	payload, err := resolver.LoadSummary(ctx, req.Ref)
	if err != nil {
		return LoadSummaryResponse{}, mapResolverError(err)
	}

	return LoadSummaryResponse{Payload: *payload}, nil
}

// LoadSelectedRaw は必要時のみ raw データを取得する。
func (c *InProcClient) LoadSelectedRaw(ctx context.Context, req LoadSelectedRawRequest) (LoadSelectedRawResponse, *Error) {
	resolver := c.getResolver()
	if resolver == nil {
		return LoadSelectedRawResponse{}, &Error{Code: CodeInternal, Message: "resolver not initialized"}
	}

	payload, err := resolver.LoadSelectedRaw(ctx, req.Ref, req.Selector)
	if err != nil {
		return LoadSelectedRawResponse{}, mapResolverError(err)
	}

	return LoadSelectedRawResponse{Payload: *payload}, nil
}

// getResolver は Service から MultiStoreResolver を取得する。
// 遅延初期化でキャッシュする。
func (c *InProcClient) getResolver() *MultiStoreResolver {
	if c.resolver != nil {
		return c.resolver
	}

	// Service の各ストア機能を使って MultiStoreResolver を構築
	c.resolver = NewMultiStoreResolver(
		c.svcSearchShort(),
		c.svcShowShort(),
		c.svcSearchJournal(),
		c.svcShowJournal(),
		c.svcSearchKnowledge(),
		c.svcShowKnowledge(),
		c.svcShowArchive(),
	)
	return c.resolver
}

// svcSearchShort は Service.Short.Search を adapter 関数に変換する。
func (c *InProcClient) svcSearchShort() func(ctx context.Context, query string, topK int) ([]Note, error) {
	return func(ctx context.Context, query string, topK int) ([]Note, error) {
		notes, err := c.Svc.SearchShort(ctx, query, topK)
		if err != nil {
			return nil, err
		}
		return mapNotes(notes), nil
	}
}

// svcShowShort は Service.Short.Show を adapter 関数に変換する。
func (c *InProcClient) svcShowShort() func(ctx context.Context, id string) (*Note, error) {
	return func(ctx context.Context, id string) (*Note, error) {
		n, err := c.Svc.GetShort(ctx, id)
		if err != nil {
			return nil, err
		}
		note := fromServiceNote(n)
		return &note, nil
	}
}

// svcSearchJournal は Service.Journal.Search を adapter 関数に変換する。
func (c *InProcClient) svcSearchJournal() func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
	return func(ctx context.Context, query string, topK int) ([]JournalNote, error) {
		notes, err := c.Svc.SearchJournal(ctx, query, topK)
		if err != nil {
			return nil, err
		}
		return mapJournalNotes(notes), nil
	}
}

// svcShowJournal は Service.Journal.Show を adapter 関数に変換する。
func (c *InProcClient) svcShowJournal() func(ctx context.Context, id string) (*JournalNote, error) {
	return func(ctx context.Context, id string) (*JournalNote, error) {
		n, err := c.Svc.GetJournal(ctx, id)
		if err != nil {
			return nil, err
		}
		note := fromServiceJournalNote(n)
		return &note, nil
	}
}

// svcSearchKnowledge は Service.Knowledge.Search を adapter 関数に変換する。
func (c *InProcClient) svcSearchKnowledge() func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
	return func(ctx context.Context, query string, topK int) ([]KnowledgeNote, error) {
		notes, err := c.Svc.SearchKnowledge(ctx, query, topK)
		if err != nil {
			return nil, err
		}
		return mapKnowledgeNotes(notes), nil
	}
}

// svcShowKnowledge は Service.Knowledge.Show を adapter 関数に変換する。
func (c *InProcClient) svcShowKnowledge() func(ctx context.Context, id string) (*KnowledgeNote, error) {
	return func(ctx context.Context, id string) (*KnowledgeNote, error) {
		n, err := c.Svc.GetKnowledge(ctx, id)
		if err != nil {
			return nil, err
		}
		note := fromServiceKnowledgeNote(n)
		return &note, nil
	}
}

// svcShowArchive は Service.Archive.Show を adapter 関数に変換する。
func (c *InProcClient) svcShowArchive() func(ctx context.Context, id string) (*ArchiveNote, error) {
	return func(ctx context.Context, id string) (*ArchiveNote, error) {
		n, err := c.Svc.GetArchive(ctx, id)
		if err != nil {
			return nil, err
		}
		note := fromServiceArchiveNote(n)
		return &note, nil
	}
}

// mapResolverError は Resolver 関連のエラーを API Error に変換する。
func mapResolverError(err error) *Error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*ErrUnsupportedRef); ok {
		return &Error{Code: CodeInvalidArgument, Message: err.Error()}
	}
	if _, ok := err.(*ErrUnresolvedRef); ok {
		return &Error{Code: CodeNotFound, Message: err.Error()}
	}
	return &Error{Code: CodeInternal, Message: err.Error()}
}

// BuildBundle は Context Bundle を構築する。
func (c *InProcClient) BuildBundle(ctx context.Context, req BuildBundleRequest) (BuildBundleResponse, *Error) {
	resolver := c.getResolver()
	if resolver == nil {
		return BuildBundleResponse{}, &Error{Code: CodeInternal, Message: "resolver not initialized"}
	}

	// Bundle ID を生成
	bundleID := generateBundleID()

	// Source refs を解決
	report, err := resolver.ResolveMany(ctx, req.SourceRefs)
	if err != nil {
		return BuildBundleResponse{}, mapResolverError(err)
	}

	// Bundle source refs を構築
	bundleSourceRefs := make([]BundleSourceRef, 0, len(report.Resolved))
	for _, resolved := range report.Resolved {
		bundleSourceRefs = append(bundleSourceRefs, BundleSourceRef{
			Ref:         resolved.Ref,
			SourceKind:  string(resolved.Ref.Type),
			SelectedRaw: false,
		})
	}

	// Summary を構築
	var summaryBuilder strings.Builder
	for _, resolved := range report.Resolved {
		if resolved.Summary != "" {
			summaryBuilder.WriteString(resolved.Summary)
			summaryBuilder.WriteString("\n")
		}
	}

	// Evidence refs と Artifact refs を分類
	var evidenceRefs, artifactRefs []TypedRef
	for _, ref := range req.SourceRefs {
		switch ref.Type {
		case EntityTypeEvidence, EntityTypeEvidenceChunk:
			evidenceRefs = append(evidenceRefs, ref)
		case EntityTypeArtifact:
			artifactRefs = append(artifactRefs, ref)
		}
	}

	// Diagnostics を構築
	diagnostics := BundleDiagnostics{
		MissingRefs:      report.Diagnostics.MissingRefs,
		UnsupportedRefs:  report.Diagnostics.UnsupportedRefs,
		ResolverWarnings: report.Diagnostics.ResolverWarnings,
		PartialBundle:    report.Diagnostics.PartialBundle,
	}

	bundle := ContextBundle{
		ID:               bundleID,
		Purpose:          req.Purpose,
		Summary:          summaryBuilder.String(),
		RebuildLevel:     "summary",
		ArtifactRefs:     artifactRefs,
		EvidenceRefs:     evidenceRefs,
		SourceRefs:       bundleSourceRefs,
		RawIncluded:      req.IncludeRaw,
		GeneratorVersion: "memx-core/v1",
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		Diagnostics:      diagnostics,
	}

	return BuildBundleResponse{Bundle: bundle}, nil
}

// generateBundleID は Bundle ID を生成する。
func generateBundleID() string {
	return fmt.Sprintf("bundle_%s", time.Now().UTC().Format("20060102_150405"))
}

// -------------------- Recall --------------------

func (c *InProcClient) Recall(ctx context.Context, req RecallRequest) (RecallResponse, *Error) {
	results, err := c.Svc.Recall(ctx, service.RecallRequest{
		Query:        req.Query,
		TopK:         req.TopK,
		MessageRange: req.MessageRange,
		Stores:       req.Stores,
		FallbackFTS:  req.FallbackFTS,
	})
	if err != nil {
		return RecallResponse{}, mapError(err)
	}
	return RecallResponse{Results: mapNotesWithContext(results)}, nil
}

func mapNotesWithContext(nwc []service.NoteWithContext) []NoteWithContext {
	out := make([]NoteWithContext, 0, len(nwc))
	for _, n := range nwc {
		out = append(out, NoteWithContext{
			Anchor: RecallNote{
				ID:      n.Anchor.ID,
				Title:   n.Anchor.Title,
				Summary: n.Anchor.Summary,
				Body:    n.Anchor.Body,
				Store:   n.Anchor.Store,
				Score:   n.Anchor.Score,
			},
			Before: mapRecallNotes(n.Before),
			After:  mapRecallNotes(n.After),
		})
	}
	return out
}

func mapRecallNotes(notes []service.RecallNote) []RecallNote {
	out := make([]RecallNote, 0, len(notes))
	for _, n := range notes {
		out = append(out, RecallNote{
			ID:      n.ID,
			Title:   n.Title,
			Summary: n.Summary,
			Body:    n.Body,
			Store:   n.Store,
			Score:   n.Score,
		})
	}
	return out
}

// -------------------- Slice Mappers --------------------

func mapNotes(ns []service.Note) []Note {
	out := make([]Note, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceNote(n))
	}
	return out
}

func mapJournalNotes(ns []service.JournalNote) []JournalNote {
	out := make([]JournalNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceJournalNote(n))
	}
	return out
}

func mapKnowledgeNotes(ns []service.KnowledgeNote) []KnowledgeNote {
	out := make([]KnowledgeNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceKnowledgeNote(n))
	}
	return out
}

func mapArchiveNotes(ns []service.ArchiveNote) []ArchiveNote {
	out := make([]ArchiveNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, fromServiceArchiveNote(n))
	}
	return out
}
