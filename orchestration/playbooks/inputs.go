package playbooks

import "github.com/secureworks/tdr-sdk-go/common"

// PlaybookInput defines the mutable fields of a playbook. Modifications the version will create a new version
type PlaybookInput struct {
	common.ObjectMetaInput
	Head       *common.ID            `json:"head"`
	Version    *PlaybookVersionInput `json:"version"`
	Categories common.IDs            `json:"categories" yaml:"categories"`
	Title      *string               `json:"title,omitempty"`
}

// PlaybookVersionInput defines the mutable fields of a playbook version
type PlaybookVersionInput struct {
	Inputs   common.Object `json:"inputs" yaml:"inputs" gorm:"type:JSONB"`
	Outputs  common.Object `json:"outputs" yaml:"outputs" gorm:"type:JSONB"`
	Requires common.IDs    `json:"requires" yaml:"requires" gorm:"type:uuid[]"`
	DSL      common.Object `json:"dsl" yaml:"dsl" gorm:"type:JSONB"`
}

// PlaybookInstanceInput defines the mutable fields of a playbook instance
type PlaybookInstanceInput struct {
	common.ObjectMetaInput
	Trigger     PlaybookTriggerInput `json:"trigger"`
	Enabled     bool                 `json:"enabled"`
	Inputs      common.Object        `json:"inputs"`
	Connections common.IDs           `json:"connections"`
}

type PlaybookTriggerInput struct {
	common.ObjectMetaInput
	TypeID common.ID     `json:"typeId"`
	Config common.Object `json:"config"`
}

type ClonePlaybookInput struct {
	PlaybookID common.ID `json:"playbookId"`
	Name       string    `json:"name"`
}
