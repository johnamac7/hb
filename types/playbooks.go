package types

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Playbooks - Playbook information
type Playbooks struct {
	Playbooks []struct {
		PlayBookName string   `json:"playbook-name" yaml:"playbook-name"`
		Description  string   `json:"description" yaml:"description"`
		Rules        []string `json:"rules"`
		Synopsis     string   `json:"synopsis" yaml:"synopsis"`
	} `json:"playbooks" yaml:"playbooks"`
}

// Parse - tries to parse yaml first, then json into the Playbooks struct
func (c *Playbooks) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	}
	return nil
}

// Dump - outputs Playbooks struct in either 'yaml' or 'json' format
func (c *Playbooks) Dump(format string) string {
	return DumpYAMLOrJSON(format, c)
}
