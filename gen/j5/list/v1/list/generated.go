package list

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/j5/list/v1/list

import ()

// Sort Proto: Sort
type Sort struct {
	Field      string `json:"field,omitempty"`
	Descending bool   `json:"descending,omitempty"`
}

// FieldType Proto: FieldType
type FieldType struct {
	Value string `json:"value,omitempty"`
	Range *Range `json:"range,omitempty"`
}

// And Proto: And
type And struct {
	Filters []*Filter `json:"filters,omitempty"`
}

// Search Proto: Search
type Search struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

// Or Proto: Or
type Or struct {
	Filters []*Filter `json:"filters,omitempty"`
}

// PageResponse Proto: PageResponse
type PageResponse struct {
	NextToken *string `json:"nextToken,omitempty"`
}

// QueryRequest Proto: QueryRequest
type QueryRequest struct {
	Searches []*Search `json:"searches,omitempty"`
	Sorts    []*Sort   `json:"sorts,omitempty"`
	Filters  []*Filter `json:"filters,omitempty"`
}

// PageRequest Proto: PageRequest
type PageRequest struct {
	Token    *string `json:"token,omitempty"`
	PageSize *int64  `json:"pageSize,omitempty,string"`
}

// Range Proto: Range
type Range struct {
	Min string `json:"min,omitempty"`
	Max string `json:"max,omitempty"`
}

// Field Proto: Field
type Field struct {
	Name string     `json:"name,omitempty"`
	Type *FieldType `json:"type"`
}

// Filter Proto Oneof: j5.list.v1.Filter
type Filter struct {
	J5TypeKey string `json:"!type,omitempty"`
	Field     *Field `json:"field,omitempty"`
	And       *And   `json:"and,omitempty"`
	Or        *Or    `json:"or,omitempty"`
}

func (s Filter) OneofKey() string {
	if s.Field != nil {
		return "field"
	}
	if s.And != nil {
		return "and"
	}
	if s.Or != nil {
		return "or"
	}
	return ""
}

func (s Filter) Type() interface{} {
	if s.Field != nil {
		return s.Field
	}
	if s.And != nil {
		return s.And
	}
	if s.Or != nil {
		return s.Or
	}
	return nil
}
