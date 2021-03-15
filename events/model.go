package events

import (
	"time"

	"github.com/secureworks/taegis-sdk-go/common"
)

type Event struct {
	ID     string        `json:"id"`
	Values common.Object `json:"values"`
}

// EventQuery defines the overall query status and metadata
type EventQuery struct {
	ID        string              `json:"id"`
	Query     string              `json:"query"`
	Status    string              `json:"status"`
	Reasons   []*EventQueryResult `json:"reasons"`
	Submitted time.Time           `json:"submitted"`
	Completed *time.Time          `json:"completed"`
	Expires   time.Time           `json:"expires"`
	Types     []string            `json:"types"`
	Metadata  common.Object       `json:"metadata"`
}

// EventQueryProgress provides basic metrics about query progress
type EventQueryProgress struct {
	// The total rows available as the result of the search.  For the Athena backend this is always the same as the maxRowsPerQuery, for Arcana it is the total rows that matched the search criteria (but could be capped for other reasons)
	TotalRows *int `json:"TotalRows"`
	// Flag that indicates whether or not the total rows is actually a lower bound (e.g. the actual number of results count be higher).  Indicates that we have *atLeast* TotalRows that match the search query
	TotalRowsIsLowerBound *bool `json:"TotalRowsIsLowerBound"`
	// Flag that indicates whether the results were truncated (there are more results available but we only returned a portion of them)
	ResultsTruncated *bool `json:"ResultsTruncated"`
}

// EventQueryOptions provides ability to override default query behavior
type EventQueryOptions struct {
	// reverses default timestamp order of descending
	TimestampAscending *bool `json:"timestampAscending"`
	// change default page size up to 1K max
	PageSize *int `json:"pageSize"`
	// change default number of rows requested up to 100K max
	MaxRows *int `json:"maxRows"`
}

// EventQueryResult returns query status and if available a page of results for a specific event type
type EventQueryResult struct {
	ID        string              `json:"id"`
	Type      string              `json:"type"`
	Backend   string              `json:"backend"`
	Status    string              `json:"status"`
	Reason    *string             `json:"reason"`
	Submitted time.Time           `json:"submitted"`
	Completed *time.Time          `json:"completed"`
	Expires   *time.Time          `json:"expires"`
	Facets    common.Object       `json:"facets"`
	Rows      []common.Object     `json:"rows"`
	Progress  *EventQueryProgress `json:"progress"`
}

// EventQueryResults contains overall query status and optionally results for a specific event type
type EventQueryResults struct {
	Query  *EventQuery       `json:"query"`
	Result *EventQueryResult `json:"result"`
	// if present points to the next logical page of results across all event types covered by the query
	Next *string `json:"next"`
	// if present points to the prev page of results
	Prev *string `json:"prev"`
}

// Results contains a list of EventQueryResults
type Results []*EventQueryResults

// GetNextPageID runs through a list of Results and extracts the next page pointer if one exists
func (r Results) GetNextPageID() *string {
	for _, result := range r {
		if result.Next != nil {
			return result.Next
		}
	}
	return nil
}
