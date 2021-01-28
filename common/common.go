package common

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	AuthorizationHeader  = "Authorization"
	XTenantContextHeader = "X-Tenant-Context"
)

type ID string

type IDs []ID

// Object is a generic json container
type Object map[string]interface{}

type ObjectMetaInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Tags        Tags    `json:"tags,omitempty"`
}

type Tags []string

func (o Object) Keys() []string {
	keys := make([]string, len(o))
	i := 0
	for k := range o {
		keys[i] = k
		i++
	}
	return keys
}

func (o *Object) UnmarshalGQL(v interface{}) error {
	switch d := v.(type) {
	case []byte:
		if err := json.Unmarshal(d, &o); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(d), &o); err != nil {
			return err
		}
	case map[string]interface{}:
		*o = d
		return nil

	default:
		return fmt.Errorf("UnmarshalGQL: invalid type %T", v)
	}
	return nil
}

func (o Object) MarshalGQL(w io.Writer) {
	_ = json.NewEncoder(w).Encode(o)
}

func NewPaginationOptions(page, perPage int) Pagination {
	return Pagination{
		Page:    &page,
		PerPage: &perPage,
	}
}

type Pagination struct {
	Page    *int `json:"page"`
	PerPage *int `json:"perPage"`
}

// IntP is a helper function to return a pointer to an int which is useful for
// optional int parameters in APIs.
func IntP(i int) *int {
	return &i
}

// StringP is a helper function to return a pointer to a string which is useful
// for optional string parameters in APIs.
func StringP(s string) *string {
	return &s
}
