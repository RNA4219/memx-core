package api

// v1.3: API はツール/AI 向けの安定 I/F。

type ErrorCode string

const (
	CodeInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeGatekeepDeny    ErrorCode = "GATEKEEP_DENY"
	CodeInternal        ErrorCode = "INTERNAL"
)

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NoteBase はノートの共通フィールド。
type NoteBase struct {
	ID             string   `json:"id"`
	Ref            TypedRef `json:"ref"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	Body           string   `json:"body"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	LastAccessedAt string   `json:"last_accessed_at"`
	AccessCount    int64    `json:"access_count"`
	SourceType     string   `json:"source_type"`
	Origin         string   `json:"origin"`
	SourceTrust    string   `json:"source_trust"`
	Sensitivity    string   `json:"sensitivity"`
}

// Note は short ストアのノート（API返却モデル）。
type Note NoteBase

// ScopedNote は working_scope を持つノート（journal/knowledge）。
type ScopedNote struct {
	NoteBase
	WorkingScope string `json:"working_scope"`
	IsPinned     bool   `json:"is_pinned"`
}

// JournalNote は journal ストアのノート。
type JournalNote ScopedNote

// KnowledgeNote は knowledge ストアのノート。
type KnowledgeNote ScopedNote

// ArchiveNote は archive ストアのノート。
type ArchiveNote NoteBase

// -------------------- Short Store --------------------

type NotesIngestRequest struct {
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Summary     string   `json:"summary,omitempty"`
	SourceType  string   `json:"source_type,omitempty"`
	Origin      string   `json:"origin,omitempty"`
	SourceTrust string   `json:"source_trust,omitempty"`
	Sensitivity string   `json:"sensitivity,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	NoLLM       bool     `json:"no_llm,omitempty"`
}

type NotesIngestResponse struct {
	Note Note `json:"note"`
}

type NotesSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

type NotesSearchResponse struct {
	Notes []Note `json:"notes"`
}

type GCOptions struct {
	DryRun bool `json:"dry_run,omitempty"`
}

type GCRunRequest struct {
	Target  string    `json:"target"`
	Options GCOptions `json:"options,omitempty"`
}

type GCRunResponse struct {
	Status string `json:"status"`
}

type SummarizeRequest struct {
	ID string `json:"id"`
}

type SummarizeResponse struct {
	Note Note `json:"note"`
}

type SummarizeBatchRequest struct {
	IDs []string `json:"ids"`
}

type SummarizeBatchResponse struct {
	Summary   string `json:"summary"`
	NoteCount int    `json:"note_count"`
}

// -------------------- Journal Store --------------------

type JournalIngestRequest struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Summary      string   `json:"summary,omitempty"`
	SourceType   string   `json:"source_type,omitempty"`
	Origin       string   `json:"origin,omitempty"`
	SourceTrust  string   `json:"source_trust,omitempty"`
	Sensitivity  string   `json:"sensitivity,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	WorkingScope string   `json:"working_scope"`
	IsPinned     bool     `json:"is_pinned,omitempty"`
	NoLLM        bool     `json:"no_llm,omitempty"`
}

type JournalIngestResponse struct {
	Note JournalNote `json:"note"`
}

type JournalSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

type JournalSearchResponse struct {
	Notes []JournalNote `json:"notes"`
}

type JournalListByScopeRequest struct {
	WorkingScope string `json:"working_scope"`
	Limit        int    `json:"limit,omitempty"`
}

type JournalListByScopeResponse struct {
	Notes []JournalNote `json:"notes"`
}

// -------------------- Knowledge Store --------------------

type KnowledgeIngestRequest struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Summary      string   `json:"summary,omitempty"`
	SourceType   string   `json:"source_type,omitempty"`
	Origin       string   `json:"origin,omitempty"`
	SourceTrust  string   `json:"source_trust,omitempty"`
	Sensitivity  string   `json:"sensitivity,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	WorkingScope string   `json:"working_scope"`
	IsPinned     bool     `json:"is_pinned,omitempty"`
	NoLLM        bool     `json:"no_llm,omitempty"`
}

type KnowledgeIngestResponse struct {
	Note KnowledgeNote `json:"note"`
}

type KnowledgeSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k,omitempty"`
}

type KnowledgeSearchResponse struct {
	Notes []KnowledgeNote `json:"notes"`
}

type KnowledgeListByScopeRequest struct {
	WorkingScope string `json:"working_scope"`
	Limit        int    `json:"limit,omitempty"`
}

type KnowledgeListByScopeResponse struct {
	Notes []KnowledgeNote `json:"notes"`
}

type KnowledgeListPinnedRequest struct {
	WorkingScope string `json:"working_scope,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

type KnowledgeListPinnedResponse struct {
	Notes []KnowledgeNote `json:"notes"`
}

type PinResponse struct {
	Success bool `json:"success"`
}

type UnpinResponse struct {
	Success bool `json:"success"`
}

// -------------------- Archive Store --------------------

type ArchiveListRequest struct {
	Limit int `json:"limit,omitempty"`
}

type ArchiveListResponse struct {
	Notes []ArchiveNote `json:"notes"`
}

type ArchiveRestoreResponse struct {
	Note Note `json:"note"`
}

// -------------------- Recall --------------------

type RecallRequest struct {
	Query        string   `json:"query"`
	TopK         int      `json:"top_k,omitempty"`
	MessageRange int      `json:"message_range,omitempty"`
	Stores       []string `json:"stores,omitempty"`
	FallbackFTS  bool     `json:"fallback_fts,omitempty"`
}

type RecallResponse struct {
	Results []NoteWithContext `json:"results"`
}

type NoteWithContext struct {
	Anchor RecallNote   `json:"anchor"`
	Before []RecallNote `json:"before,omitempty"`
	After  []RecallNote `json:"after,omitempty"`
}

type RecallNote struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Summary string  `json:"summary,omitempty"`
	Body    string  `json:"body,omitempty"`
	Store   string  `json:"store"`
	Score   float64 `json:"score"`
}

// -------------------- Resolver API (P4) --------------------

// ResolveRefRequest は単一の typed_ref 解決リクエスト。
type ResolveRefRequest struct {
	Ref TypedRef `json:"ref"`
}

// ResolveRefResponse は単一の typed_ref 解決レスポンス。
type ResolveRefResponse struct {
	Resolved ResolvedRef `json:"resolved"`
}

// ResolveManyRequest は複数の typed_ref 一括解決リクエスト。
type ResolveManyRequest struct {
	Refs []TypedRef `json:"refs"`
}

// ResolveManyResponse は複数の typed_ref 一括解決レスポンス。
type ResolveManyResponse struct {
	Report ResolveReport `json:"report"`
}

// LoadSummaryRequest は要約取得リクエスト（summary-first retrieval）。
type LoadSummaryRequest struct {
	Ref TypedRef `json:"ref"`
}

// LoadSummaryResponse は要約取得レスポンス。
type LoadSummaryResponse struct {
	Payload SummaryPayload `json:"payload"`
}

// LoadSelectedRawRequest は raw データ取得リクエスト。
type LoadSelectedRawRequest struct {
	Ref      TypedRef    `json:"ref"`
	Selector RawSelector `json:"selector"`
}

// LoadSelectedRawResponse は raw データ取得レスポンス。
type LoadSelectedRawResponse struct {
	Payload RawPayload `json:"payload"`
}

// -------------------- Context Bundle (P4) --------------------

// ContextBundle は再開用の文脈バンドル。
// P4仕様: purpose, summary, decision_digest, open_question_digest, artifact_refs, evidence_refs, source_refs, raw_included, resolver_diagnostics
type ContextBundle struct {
	ID                 string             `json:"id"`
	Purpose            string             `json:"purpose"`
	Summary            string             `json:"summary"`
	RebuildLevel       string             `json:"rebuild_level"`
	DecisionDigest     string             `json:"decision_digest,omitempty"`
	OpenQuestionDigest string             `json:"open_question_digest,omitempty"`
	ArtifactRefs       []TypedRef         `json:"artifact_refs,omitempty"`
	EvidenceRefs       []TypedRef         `json:"evidence_refs,omitempty"`
	SourceRefs         []BundleSourceRef  `json:"source_refs,omitempty"`
	TrackerRefs        []TypedRef         `json:"tracker_refs,omitempty"`
	RawIncluded        bool               `json:"raw_included"`
	GeneratorVersion   string             `json:"generator_version"`
	GeneratedAt        string             `json:"generated_at"`
	Diagnostics        BundleDiagnostics  `json:"diagnostics"`
}

// BundleSourceRef はバンドル内のソース参照。
type BundleSourceRef struct {
	Ref          TypedRef `json:"ref"`
	SourceKind   string   `json:"source_kind"`
	SelectedRaw  bool     `json:"selected_raw"`
	MetadataJSON string   `json:"metadata_json,omitempty"`
}

// BundleDiagnostics はバンドル生成時の診断情報。
type BundleDiagnostics struct {
	MissingRefs      []TypedRef `json:"missing_refs,omitempty"`
	UnsupportedRefs  []TypedRef `json:"unsupported_refs,omitempty"`
	ResolverWarnings []string   `json:"resolver_warnings,omitempty"`
	PartialBundle    bool       `json:"partial_bundle"`
}

// BuildBundleRequest は Context Bundle 構築リクエスト。
type BuildBundleRequest struct {
	Purpose      string     `json:"purpose"`
	SourceRefs   []TypedRef `json:"source_refs"`
	IncludeRaw   bool       `json:"include_raw,omitempty"`
	RawSelectors map[string]RawSelector `json:"raw_selectors,omitempty"`
}

// BuildBundleResponse は Context Bundle 構築レスポンス。
type BuildBundleResponse struct {
	Bundle ContextBundle `json:"bundle"`
}
