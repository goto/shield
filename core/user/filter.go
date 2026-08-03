package user

type Filter struct {
	Limit                     int32
	Page                      int32
	Keyword                   string
	ProjectID                 string
	ServiceDataKeyResourceIds []string

	// MetadataPaths and the fields below filter users by their servicedata
	// metadata, mirroring Guardian's appeal_details* generic JSONB filter. Each
	// path is dot-delimited: the first segment is the metadata (servicedata) key
	// and the remaining segments navigate into that key's JSON value, e.g.
	// "functional_roles.roles". When empty, no metadata filter is applied.
	MetadataPaths []string
	// Positive filters (OR-ed together across paths). A matched value that is a
	// JSON array matches when it contains the value.
	Metadatas          []string // equals / IN
	MetadataStartsWith string
	MetadataEndsWith   string
	MetadataContains   string
	// Negative filters (AND-ed together; user is excluded on a match).
	NotMetadatas          []string // not equals / NOT IN
	MetadataNotStartsWith string
	MetadataNotEndsWith   string
	MetadataNotContains   string
}
