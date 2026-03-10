package api

import (
	"fmt"
	"strings"
)

// Domain は typed_ref のシステム領域を表す。
type Domain string

const (
	DomainMemx          Domain = "memx"
	DomainAgentTaskstate Domain = "agent-taskstate"
	DomainTracker       Domain = "tracker"
)

// EntityType は memx 内のエンティティ種別。
type EntityType string

const (
	EntityTypeEvidence      EntityType = "evidence"
	EntityTypeEvidenceChunk EntityType = "evidence_chunk"
	EntityTypeKnowledge     EntityType = "knowledge"
	EntityTypeArtifact      EntityType = "artifact"
	EntityTypeLineage       EntityType = "lineage"
)

// Provider は typed_ref のデータソースを表す。
type Provider string

const (
	ProviderLocal  Provider = "local"
	ProviderJira   Provider = "jira"
	ProviderGitHub Provider = "github"
	ProviderLinear Provider = "linear"
)

// DefaultProvider は provider が省略された場合のデフォルト値。
const DefaultProvider Provider = ProviderLocal

// TypedRef はエンティティへの型付き参照。
// Canonical format: <domain>:<entity_type>:<provider>:<entity_id>
// 例: memx:evidence:local:01HXXXXXXX
//
// 移行期間中は3セグメント形式（memx:<type>:<id>）も受理し、
// provider=local として正規化する。
type TypedRef struct {
	Domain     Domain     `json:"domain"`
	Type       EntityType `json:"type"`
	Provider   Provider   `json:"provider"`
	ID         string     `json:"id"`
}

// String は TypedRef を canonical format（4セグメント）の文字列に変換する。
func (r TypedRef) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", r.Domain, r.Type, r.Provider, r.ID)
}

// MarshalText は encoding.TextMarshaler インターフェースを実装する。
func (r TypedRef) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// UnmarshalText は encoding.TextUnmarshaler インターフェースを実装する。
func (r *TypedRef) UnmarshalText(text []byte) error {
	ref, err := ParseTypedRef(string(text))
	if err != nil {
		return err
	}
	*r = ref
	return nil
}

// ParseTypedRef は typed_ref 文字列を TypedRef に変換する。
// 3セグメント形式（memx:<type>:<id>）と4セグメント形式（memx:<type>:<provider>:<id>）の両方を受理する。
// 3セグメント形式の場合は provider=local として正規化する。
func ParseTypedRef(s string) (TypedRef, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return TypedRef{}, fmt.Errorf("empty ref")
	}

	parts := strings.Split(s, ":")

	// 3セグメント形式（旧形式）の受理
	if len(parts) == 3 {
		return parseThreeSegment(parts)
	}

	// 4セグメント形式（canonical）の受理
	if len(parts) == 4 {
		return parseFourSegment(parts)
	}

	return TypedRef{}, fmt.Errorf("invalid ref format: %s (expected <domain>:<type>:[<provider>:]<id>)", s)
}

// parseThreeSegment は3セグメント形式をパースし、provider=local として正規化する。
func parseThreeSegment(parts []string) (TypedRef, error) {
	domain := Domain(parts[0])
	if domain != DomainMemx {
		return TypedRef{}, fmt.Errorf("invalid domain: %s (expected 'memx' for 3-segment format)", domain)
	}

	entityType := EntityType(parts[1])
	if err := validateEntityType(domain, entityType); err != nil {
		return TypedRef{}, err
	}

	id := parts[2]
	if id == "" {
		return TypedRef{}, fmt.Errorf("empty id in ref")
	}

	return TypedRef{
		Domain:   domain,
		Type:     entityType,
		Provider: DefaultProvider,
		ID:       id,
	}, nil
}

// parseFourSegment は4セグメント形式（canonical）をパースする。
func parseFourSegment(parts []string) (TypedRef, error) {
	domain := Domain(parts[0])
	if err := validateDomain(domain); err != nil {
		return TypedRef{}, err
	}

	entityType := EntityType(parts[1])
	if err := validateEntityType(domain, entityType); err != nil {
		return TypedRef{}, err
	}

	provider := Provider(parts[2])
	if provider == "" {
		return TypedRef{}, fmt.Errorf("empty provider in ref")
	}

	id := parts[3]
	if id == "" {
		return TypedRef{}, fmt.Errorf("empty id in ref")
	}

	return TypedRef{
		Domain:   domain,
		Type:     entityType,
		Provider: provider,
		ID:       id,
	}, nil
}

// validateDomain はドメインが有効かどうかを検証する。
func validateDomain(domain Domain) error {
	switch domain {
	case DomainMemx, DomainAgentTaskstate, DomainTracker:
		return nil
	default:
		return fmt.Errorf("invalid domain: %s (expected memx, agent-taskstate, or tracker)", domain)
	}
}

// validateEntityType はエンティティタイプが有効かどうかを検証する。
// memx ドメインの場合は既知のタイプのみ許可し、それ以外のドメインでは空でないことを確認する。
func validateEntityType(domain Domain, entityType EntityType) error {
	if entityType == "" {
		return fmt.Errorf("empty entity type")
	}

	// memx ドメインの場合は既知のタイプのみ許可
	if domain == DomainMemx {
		switch entityType {
		case EntityTypeEvidence, EntityTypeEvidenceChunk, EntityTypeKnowledge, EntityTypeArtifact, EntityTypeLineage:
			return nil
		default:
			return fmt.Errorf("invalid entity type for memx domain: %s", entityType)
		}
	}

	// 他のドメイン（agent-taskstate, tracker）では空でないことを確認するのみ
	// 実在性確認は別責務とする
	return nil
}

// MustParseTypedRef は ParseTypedRef のパニック版。
func MustParseTypedRef(s string) TypedRef {
	ref, err := ParseTypedRef(s)
	if err != nil {
		panic(err)
	}
	return ref
}

// NewTypedRef は指定された型とIDから TypedRef を作成する（memx ドメイン、local provider）。
func NewTypedRef(entityType EntityType, id string) TypedRef {
	return TypedRef{
		Domain:   DomainMemx,
		Type:     entityType,
		Provider: DefaultProvider,
		ID:       id,
	}
}

// NewTypedRefWithProvider は指定された provider で TypedRef を作成する。
func NewTypedRefWithProvider(domain Domain, entityType EntityType, provider Provider, id string) TypedRef {
	return TypedRef{
		Domain:   domain,
		Type:     entityType,
		Provider: provider,
		ID:       id,
	}
}

// Ref は TypedRef.String() のエイリアス。
func (r TypedRef) Ref() string {
	return r.String()
}

// IsZero は TypedRef がゼロ値かどうかを返す。
func (r TypedRef) IsZero() bool {
	return r.Domain == "" && r.Type == "" && r.Provider == "" && r.ID == ""
}

// IsValid は TypedRef が有効かどうかを返す。
func (r TypedRef) IsValid() bool {
	return r.Domain != "" && r.Type != "" && r.Provider != "" && r.ID != ""
}

// Canonical は TypedRef を canonical format で返す。
// 既に canonical の場合はそのまま返す。
func (r TypedRef) Canonical() string {
	return r.String()
}