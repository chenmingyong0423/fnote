package domain

import "time"

const (
	ConfigHealthLevelRequired    = "required"
	ConfigHealthLevelRecommended = "recommended"
	ConfigHealthLevelOptional    = "optional"

	ConfigHealthStatusOK      = "ok"
	ConfigHealthStatusMissing = "missing"
	ConfigHealthStatusIgnored = "ignored"
	ConfigHealthStatusSnoozed = "snoozed"

	ConfigCheckStateActive  = "active"
	ConfigCheckStateIgnored = "ignored"
	ConfigCheckStateSnoozed = "snoozed"
)

type ConfigHealth struct {
	Score       int
	Required    ConfigHealthGroup
	Recommended ConfigHealthGroup
	Optional    ConfigHealthGroup
	Items       []ConfigHealthItem
}

type ConfigHealthGroup struct {
	Done      int
	Total     int
	Completed bool
}

type ConfigHealthItem struct {
	Key           string
	Label         string
	Level         string
	Status        string
	Configured    bool
	MissingFields []string
	Description   string
	Href          string
	SnoozedUntil  *time.Time
	IgnoredReason string
}

type ConfigCheckState struct {
	CheckKey      string
	Status        string
	IgnoredReason string
	SnoozedUntil  *time.Time
	UpdatedAt     time.Time
}
