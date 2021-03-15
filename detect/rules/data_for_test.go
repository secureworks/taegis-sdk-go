package rules

import (
	"time"

	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/testutils"
)

// This file is separate from client_test.go because that file is code generated
// and we don't want any custom code there.

func referenceInputFromReferences(refs []RuleReference) []RuleReferenceInput {
	result := []RuleReferenceInput{}

	for _, ref := range refs {
		result = append(result, RuleReferenceInput{
			Description: ref.Description,
			URL:         ref.URL,
		})
	}

	return result
}

func ruleInputFromRule(r *Rule) RuleInput {
	return RuleInput{
		ID:               &r.ID,
		EventType:        &r.EventType,
		Name:             &r.Name,
		Description:      &r.Description,
		Visibility:       &r.Visibility,
		ResultVisibility: &r.ResultVisibility,
		Severity:         &r.Severity,
		Confidence:       &r.Confidence,
		CreateAlert:      &r.CreateAlert,
		Tags:             r.Tags,
		AttackCategories: r.AttackCategories,
		EndpointPlatform: r.EndpointPlatform,
		References:       referenceInputFromReferences(r.References),
	}
}

func convertFilters(filters []RuleFilterInput) []interface{} {
	result := []interface{}{}

	for _, f := range filters {
		result = append(result, testutils.ToGenericMap(f))
	}

	return result
}

var (
	id = "123"

	count = common.IntP(10)
	page  = common.IntP(1)

	ruleType    = RuleTypeRegex
	ruleTypePtr = &ruleType

	timestampStr = time.Now().Format(time.RFC3339)
	// Build the time from the string to ensure it is the same as a JSON time
	timestamp, _ = time.Parse(time.RFC3339, timestampStr)

	filter = RuleFilter{
		Key:     "process",
		Pattern: "cmd.exe",
	}

	eventType    = RuleEventTypeProcess
	eventTypePtr = &eventType

	rule = &Rule{
		ID:               id,
		EventType:        eventType,
		Name:             "Test Rule",
		Description:      "A test rule",
		Visibility:       RuleVisibilityVisible,
		ResultVisibility: RuleVisibilityVisible,
		Severity:         0.5,
		Confidence:       0.5,
		Enabled:          true,
		CreateAlert:      true,
		Tags:             []string{"tag1", "tag2"},
		EndpointPlatform: []RuleEndpointPlatform{RuleEndpointPlatformPlatformWindows},
		References: []RuleReference{{
			Description: "A reference",
			URL:         "https://www.secureworks.com",
		}},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Filters:   []RuleFilter{filter},
	}

	getRulesResponse = struct {
		Out []*Rule `json:"rules"`
	}{
		Out: []*Rule{rule},
	}

	getDeletedRulesResponse = struct {
		Out []*Rule `json:"deletedRules"`
	}{
		Out: []*Rule{rule},
	}

	getRulesCountResponse = struct {
		Out int `json:"rulesCount"`
	}{
		Out: 10,
	}

	getRulesForEventResponse = struct {
		Out []*Rule `json:"rulesForEvent"`
	}{
		Out: []*Rule{rule},
	}

	getRulesForEventCountResponse = struct {
		Out int `json:"rulesForEventCount"`
	}{
		Out: 5,
	}

	getRuleResponse = struct {
		Out *Rule `json:"rule"`
	}{
		Out: rule,
	}

	getFilterKeysResponse = struct {
		Out []string `json:"filterKeys"`
	}{
		Out: []string{"a", "b", "c"},
	}

	getChangesSinceResponse = struct {
		Out []*Rule `json:"changesSince"`
	}{
		Out: []*Rule{rule},
	}

	ruleInput = ruleInputFromRule(rule)
	filters   = []RuleFilterInput{
		{
			Key:     filter.Key,
			Pattern: filter.Pattern,
		},
	}
	genericFilters = convertFilters(filters)

	createRuleResponse = struct {
		Out Rule `json:"createRule"`
	}{
		Out: *rule,
	}

	ruleID      = rule.ID
	filterInput = filters[0]

	addFilterToRuleResponse = struct {
		Out RuleFilter `json:"addFilterToRule"`
	}{
		Out: filter,
	}

	updateRuleResponse = struct {
		Out Rule `json:"updateRule"`
	}{
		Out: *rule,
	}

	deleteRuleResponse = struct {
		Out Rule `json:"deleteRule"`
	}{
		Out: *rule,
	}

	filterID = "456"

	updateFilterResponse = struct {
		Out RuleFilter `json:"updateFilter"`
	}{
		Out: filter,
	}

	deleteFilterResponse = struct {
		Out RuleFilter `json:"deleteFilter"`
	}{
		Out: filter,
	}

	redQLFilter = RuleRedQLFilter{
		ID:    id,
		Query: "from process commandline = explorer.exe",
	}
	redQLFilterInput = RuleRedQLFilterInput{
		Query: redQLFilter.Query,
	}

	createRedQLRuleResponse = struct {
		Out Rule `json:"createRedQLRule"`
	}{
		Out: *rule,
	}

	updateRedQLFilterResponse = struct {
		Out RuleRedQLFilter `json:"updateRedQLFilter"`
	}{
		Out: redQLFilter,
	}

	disableRuleResponse = struct {
		Out Rule `json:"disableRule"`
	}{
		Out: *rule,
	}

	enableRuleResponse = struct {
		Out Rule `json:"enableRule"`
	}{
		Out: *rule,
	}
)
