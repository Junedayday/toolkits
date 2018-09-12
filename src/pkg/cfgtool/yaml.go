package cfgtool

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadCfgFromYamlFile(path string, cfg interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, cfg)
}
