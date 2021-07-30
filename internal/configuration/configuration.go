package configuration

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var AppVersion = "v0.1-beta"

type ConfigStruct struct {
	Version          string  `json:"version,omitempty"`
	ConfigFile       *string `json:"configFile,omitempty"`
	VariablesFile    *string `json:"variablesFile,omitempty"`
	TemplatesFile    *string `json:"templatesFile,omitempty"`
	LabelDirectory   *string `json:"labelDirectory,omitempty"`
	Port             *int    `json:"port,omitempty"`
	InitialSetupDone *bool   `json:"initialSetupDone,omitempty"`

	LastVersionCheckDate *time.Time `json:"versionCheckDate,omitempty"`
	UpdateAvailable      *bool      `json:"updateAvailable,omitempty"`
}

type githubTagsArray struct {
	Name string `json:"name"`
}

func (c *ConfigStruct) IsUpdateAvailable() bool {

	if c.UpdateAvailable != nil && *c.UpdateAvailable && c.Version == AppVersion {
		return true
	}
	if c.LastVersionCheckDate == nil || time.Now().After((*c.LastVersionCheckDate).Add(time.Hour+24)) {

		log.Println("[load] Checking for any updates")
		if resp, err := http.Get("https://api.github.com/repos/decalibrate/overlay-label-manager/tags"); err == nil {
			if body, e2 := io.ReadAll(resp.Body); e2 == nil {
				tags := make([]githubTagsArray, 0)
				json.Unmarshal(body, &tags)

				if len(tags) > 0 && tags[0].Name != AppVersion {
					b := true
					c.UpdateAvailable = &b

					return b
				}
			}
		}
	}

	now := time.Now()
	c.LastVersionCheckDate = &now
	c.Version = AppVersion
	c.UpdateAvailable = nil

	return false
}
