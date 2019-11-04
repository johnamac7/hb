package types

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// PlaybookInstances - wrapper type for Device Groups, with only the Playbook relevant information described
type PlaybookInstances struct {
	DeviceGroup []struct {
		DeviceGroupName string    `json:"device-group-name" yaml:"device-group-name"`
		Devices         *[]string `json:"devices,omitempty" yaml:"devices,omitempty"`
		Playbooks       []string  `json:"playbooks,omitempty" yaml:"playbooks,omitempty"`
		Variable        []struct {
			InstanceID    string `json:"instance-id" yaml:"instance-id"`
			Playbook      string `json:"playbook"`
			Rule          string `json:"rule"`
			VariableValue []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"variable-value,omitempty" yaml:"variable-value,omitempty"`
		} `json:"variable"`
	} `json:"device-group" yaml:"device-group"`
}

// Parse - tries to parse yaml first, then json into the PlaybookInstances struct
func (c *PlaybookInstances) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	}
	return nil
}

// Dump - outputs PlaybookInstances struct in either 'yaml' or 'json' format
func (c *PlaybookInstances) Dump(format string) string {
	return DumpYAMLOrJSON(format, c)
}
