package sdk

import (
	"time"
)

type SwaThreeFieldApplication struct {
	Credentials   *SchemeApplicationCredentials     `json:"credentials,omitempty"`
	Embedded      interface{}                       `json:"_embedded,omitempty"`
	Links         interface{}                       `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility         `json:"accessibility,omitempty"`
	Created       *time.Time                        `json:"created,omitempty"`
	Features      []string                          `json:"features,omitempty"`
	Id            string                            `json:"id,omitempty"`
	Label         string                            `json:"label,omitempty"`
	LastUpdated   *time.Time                        `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing             `json:"licensing,omitempty"`
	Name          string                            `json:"name,omitempty"`
	Profile       interface{}                       `json:"profile,omitempty"`
	Settings      *SwaThreeFieldApplicationSettings `json:"settings,omitempty"`
	SignOnMode    string                            `json:"signOnMode,omitempty"`
	Status        string                            `json:"status,omitempty"`
	Visibility    *ApplicationVisibility            `json:"visibility,omitempty"`
}

func NewSwaThreeFieldApplication() *SwaThreeFieldApplication {
	return &SwaThreeFieldApplication{
		Name:       "template_swa3field",
		SignOnMode: "BROWSER_PLUGIN",
	}
}

func (a *SwaThreeFieldApplication) IsApplicationInstance() bool {
	return true
}
