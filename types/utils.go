package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration - structures that get loaded from files
type Configuration interface {
	Parse(data []byte) error
	Dump(format string) string
}

// LoadConfiguration - populates a configuration type with data from file
func LoadConfiguration(filelocation string, configuration Configuration) error {
	data, err := ioutil.ReadFile(filelocation)
	if err != nil {
		return err
	}
	if err := configuration.Parse(data); err != nil {
		return err
	}
	return nil
}

// DumpYAMLOrJSON - For a configuration, output json or yaml
func DumpYAMLOrJSON(format string, configuration Configuration) string {
	var data []byte
	switch format {
	case "yaml":
		_ = yaml.Unmarshal(data, configuration)
	default:
		_ = json.Unmarshal(data, configuration)
	}
	return fmt.Sprintf("%v", data)
}
