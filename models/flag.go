package models

// Flags represents feature flags.
type Flags struct {
	Theme           string `json:"theme"`
	Welcome         bool   `json:"welcome"`
	Sidebar         bool   `json:"sidebar"`
	SearchAutofocus bool   `json:"search_autofocus"`
	ViewRedesign    bool   `json:"view_redesign"`
}
